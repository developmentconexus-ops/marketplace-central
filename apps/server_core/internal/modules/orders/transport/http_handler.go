package transport

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/application"
	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type OrderImporter interface {
	Import(ctx context.Context, input application.ImportOrdersInput) (domain.ImportResult, error)
}

type OrderLister interface {
	List(ctx context.Context, input application.ListOrdersInput) ([]domain.MarketplaceOrder, error)
}

type OrderReadLister interface {
	List(ctx context.Context, query ports.OrderListQuery) (ports.OrderPage, error)
	Get(ctx context.Context, installationID, providerOrderID string) (domain.OrderReadModel, error)
}

var ErrInstallationRequired = errors.New("installation_required")

type Handler struct {
	importer   OrderImporter
	lister     OrderLister
	readLister OrderReadLister
	enricher   *application.EnrichService
	summary    application.SummaryService
}

func NewHandler(importer OrderImporter, lister OrderLister) Handler {
	return Handler{importer: importer, lister: lister}
}

func NewHandlerWithReader(importer OrderImporter, reader OrderReadLister) Handler {
	return Handler{importer: importer, readLister: reader}
}

// NewHandlerWithEnricher wires the /orders read transport with a read-time
// EnrichService: handleReadList/handleGet map results through enricher and
// emit the enriched response DTO instead of the raw domain.OrderReadModel.
// enricher may be nil (defensive) — the handler then falls back to the raw
// marshal exactly like NewHandlerWithReader, so no caller regresses.
func NewHandlerWithEnricher(importer OrderImporter, reader OrderReadLister, enricher *application.EnrichService) Handler {
	return Handler{importer: importer, readLister: reader, enricher: enricher}
}

// NewHandlerWithSummary extends NewHandlerWithEnricher with the counts
// summary service backing GET /orders/summary (F01-S5). Additive: the
// enricher-only constructor keeps working for its existing callers, and
// summary is a value type whose zero value already answers honestly (see
// application.SummaryService.Summary, store-nil branch), so passing a fresh
// Handler through this constructor never regresses NewHandlerWithEnricher
// call sites.
func NewHandlerWithSummary(importer OrderImporter, reader OrderReadLister, enricher *application.EnrichService, summary application.SummaryService) Handler {
	return Handler{importer: importer, readLister: reader, enricher: enricher, summary: summary}
}

func (h Handler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("/orders/import", h.handleImport)
	mux.HandleFunc("/orders", h.handleList)
	mux.HandleFunc("/orders/summary", h.handleSummary)
	if h.readLister != nil {
		mux.HandleFunc("/orders/{provider_order_id}", h.handleGet)
	}
}

func (h Handler) handleImport(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeOrdersError(w, http.StatusMethodNotAllowed, "ORDERS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	var req struct {
		InstallationID string `json:"installation_id"`
		Limit          int    `json:"limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeOrdersError(w, http.StatusBadRequest, "ORDERS_INVALID_REQUEST", "malformed request body")
		return
	}
	result, err := h.importer.Import(r.Context(), application.ImportOrdersInput{
		InstallationID: req.InstallationID,
		Limit:          req.Limit,
	})
	if err != nil {
		status, code := mapOrdersError(err)
		slog.Error("orders.import", "action", "import", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeOrdersError(w, status, code, err.Error())
		return
	}
	slog.Info("orders.import", "action", "import", "result", "200", "count", result.ImportedCount, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h Handler) handleList(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeOrdersError(w, http.StatusMethodNotAllowed, "ORDERS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	if h.readLister != nil {
		h.handleReadList(w, r)
		return
	}
	if h.lister == nil {
		writeOrdersError(w, http.StatusInternalServerError, "ORDERS_INTERNAL_ERROR", "orders list is not configured")
		return
	}
	items, err := h.lister.List(r.Context(), application.ListOrdersInput{
		InstallationID: r.URL.Query().Get("installation_id"),
		Limit:          parsePositiveInt(r.URL.Query().Get("limit"), 20),
	})
	if err != nil {
		status, code := mapOrdersError(err)
		slog.Error("orders.list", "action", "list", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeOrdersError(w, status, code, err.Error())
		return
	}
	slog.Info("orders.list", "action", "list", "result", "200", "count", len(items), "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

type orderPageEnvelope struct {
	Items      []domain.OrderReadModel `json:"items"`
	NextCursor *string                 `json:"next_cursor"`
}

type enrichedOrderPageEnvelope struct {
	Items      []enrichedOrderDTO `json:"items"`
	NextCursor *string            `json:"next_cursor"`
}

func (h Handler) handleReadList(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	installationID, err := requiredInstallation(values)
	if err != nil {
		writeOrderReadError(w, err)
		return
	}
	query, err := ParseOrderQuery(values)
	if err != nil {
		writeOrderReadError(w, err)
		return
	}
	page, err := h.readLister.List(r.Context(), query)
	if err != nil {
		writeOrderReadError(w, err)
		return
	}
	if page.Items == nil {
		page.Items = []domain.OrderReadModel{}
	}
	var next *string
	if page.NextCursor != nil {
		encoded, encodeErr := page.NextCursor.Encode()
		if encodeErr != nil {
			writeOrderReadError(w, encodeErr)
			return
		}
		next = &encoded
	}
	if h.enricher != nil {
		enriched := h.enricher.Enrich(r.Context(), installationID, page.Items)
		dtos := make([]enrichedOrderDTO, len(enriched))
		for i, e := range enriched {
			dtos[i] = mapEnrichedOrder(e)
		}
		httpx.WriteJSON(w, http.StatusOK, enrichedOrderPageEnvelope{Items: dtos, NextCursor: next})
		return
	}
	httpx.WriteJSON(w, http.StatusOK, orderPageEnvelope{Items: page.Items, NextCursor: next})
}

func (h Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeOrdersError(w, http.StatusMethodNotAllowed, "ORDERS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	if h.readLister == nil {
		writeOrdersError(w, http.StatusInternalServerError, "ORDERS_INTERNAL_ERROR", "orders read is not configured")
		return
	}
	values := r.URL.Query()
	installationID, err := requiredInstallation(values)
	if err != nil {
		writeOrderReadError(w, err)
		return
	}
	for key := range values {
		if key != "installation_id" {
			writeOrderReadError(w, &InvalidFilterError{Key: key})
			return
		}
	}
	providerOrderID := strings.TrimSpace(r.PathValue("provider_order_id"))
	if providerOrderID == "" {
		writeOrderReadError(w, &ports.OrderNotFoundError{InstallationID: installationID})
		return
	}
	model, err := h.readLister.Get(r.Context(), installationID, providerOrderID)
	if err != nil {
		writeOrderReadError(w, err)
		return
	}
	if h.enricher != nil {
		enriched := h.enricher.Enrich(r.Context(), installationID, []domain.OrderReadModel{model})
		httpx.WriteJSON(w, http.StatusOK, mapEnrichedOrder(enriched[0]))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, model)
}

// orderSummaryDTO is the response shape for GET /orders/summary. Keys are
// locked byte-identical to the OpenAPI OrderSummary schema and the SDK
// OrderSummary interface (F01-S5 CONSISTENCY requirement).
type orderSummaryDTO struct {
	Today     int64 `json:"today"`
	SevenDays int64 `json:"seven_days"`
}

// orderBucketCountsDTO is the response shape for GET /orders/summary?by=status
// (F01-A). The four keys are locked byte-identical to domain.OrderBucket's
// novo/faturar/enviar/enviado constants, the OpenAPI by_status schema, and
// the SDK OrderSummary.by_status interface (CONSISTENCY requirement).
type orderBucketCountsDTO struct {
	Novo    int64 `json:"novo"`
	Faturar int64 `json:"faturar"`
	Enviar  int64 `json:"enviar"`
	Enviado int64 `json:"enviado"`
}

// handleSummary resolves installationID the same way handleGet/handleReadList
// do (requiredInstallation over the query string) and delegates counting to
// application.SummaryService. On any Summary() error it writes an honest
// error status — never a fabricated {"today":0,"seven_days":0} body — since
// an unresolved count is unknown, not zero.
//
// The optional "by" query param selects the response dimension (F01-A):
// absent/empty keeps the existing {today, seven_days} shape; "status"
// returns {by_status: {...}} bucket counts; any other value is an honest 400
// (never a silent fallback to the default shape).
func (h Handler) handleSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeOrdersError(w, http.StatusMethodNotAllowed, "ORDERS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	installationID, err := requiredInstallation(r.URL.Query())
	if err != nil {
		writeOrderReadError(w, err)
		return
	}
	switch by := r.URL.Query().Get("by"); by {
	case "":
		h.handleSummaryDefault(w, r, installationID)
	case "status":
		h.handleSummaryByStatus(w, r, installationID)
	default:
		writeOrdersErrorDetails(w, http.StatusBadRequest, "unsupported_summary_dimension", "unsupported summary dimension", map[string]any{"by": by})
	}
}

func (h Handler) handleSummaryDefault(w http.ResponseWriter, r *http.Request, installationID string) {
	summary, err := h.summary.Summary(r.Context(), installationID, time.Now())
	if err != nil {
		writeSummaryError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, orderSummaryDTO{Today: summary.Today, SevenDays: summary.SevenDays})
}

func (h Handler) handleSummaryByStatus(w http.ResponseWriter, r *http.Request, installationID string) {
	counts, err := h.summary.BucketCounts(r.Context(), installationID)
	if err != nil {
		writeSummaryError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"by_status": orderBucketCountsDTO{Novo: counts.Novo, Faturar: counts.Faturar, Enviar: counts.Enviar, Enviado: counts.Enviado},
	})
}

func writeSummaryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, application.ErrSummaryInstallationRequired):
		writeOrdersErrorDetails(w, http.StatusBadRequest, "installation_required", "installation_id é obrigatório", map[string]any{"key": "installation_id"})
	case errors.Is(err, application.ErrSummaryStoreNotConfigured):
		writeOrdersError(w, http.StatusServiceUnavailable, "ORDERS_SUMMARY_STORE_NOT_CONFIGURED", "orders summary is not configured")
	default:
		writeOrdersError(w, http.StatusInternalServerError, "ORDERS_INTERNAL_ERROR", "orders summary failed")
	}
}

// enrichedBuyerDTO is the ONLY buyer identification allowed onto the
// transport: masked Display plus honest, possibly-unknown City/UF. Never the
// raw provider nickname or any other PII (C06/LGPD).
type enrichedBuyerDTO struct {
	Display string  `json:"display"`
	City    *string `json:"city,omitempty"`
	UF      *string `json:"uf,omitempty"`
}

type enrichedSLADTO struct {
	Due      *time.Time `json:"due,omitempty"`
	Atrasado *bool      `json:"atrasado,omitempty"`
}

type enrichedRastreioDTO struct {
	ShipmentID string `json:"shipment_id"`
	Status     string `json:"status"`
	// Substatus is the ML shipment sub-status (e.g. "out_for_delivery",
	// "receiver_absent"). Additive/omitempty: absent when the provider reports
	// no sub-status, never a fabricated value (ADR-17).
	Substatus string `json:"substatus,omitempty"`
}

// enrichedItemDTO carries the existing item JSON fields (embedded, so no
// duplication/drift risk) plus the read-time-only per-unit cost.
// CustoUnitario is omitted (never 0) when the cost could not be honestly
// resolved (ADR-17).
type enrichedItemDTO struct {
	domain.MarketplaceOrderItem
	CustoUnitario *float64 `json:"custo_unitario,omitempty"`
}

// decomposicaoDTO is the response shape for the order-level cost/fee
// decomposition (F01-C1). Every amount is a pointer/omitempty: nil (absent
// key) means the component could not be honestly sourced (ADR-17).
// ComponentesDesconhecidos is never omitted — it is what explains the "—"s.
type decomposicaoDTO struct {
	Comissao                 *float64 `json:"comissao,omitempty"`
	TaxaFixa                 *float64 `json:"taxa_fixa,omitempty"`
	Frete                    *float64 `json:"frete,omitempty"`
	Imposto                  *float64 `json:"imposto,omitempty"`
	Difal                    *float64 `json:"difal,omitempty"`
	TarifaFull               *float64 `json:"tarifa_full,omitempty"`
	Custo                    *float64 `json:"custo,omitempty"`
	MargemValor              *float64 `json:"margem_valor,omitempty"`
	MargemPct                *float64 `json:"margem_pct,omitempty"`
	ComponentesDesconhecidos []string `json:"componentes_desconhecidos"`
}

// difalDTO is the response shape for the order-level DIFAL fact (F01-C1).
// All fields pointer/omitempty: nil (absent key) means unknown (ADR-17).
type difalDTO struct {
	Amount  *float64   `json:"amount,omitempty"`
	UFRoute *string    `json:"uf_route,omitempty"`
	DueDate *time.Time `json:"due_date,omitempty"`
	Paid    *bool      `json:"paid,omitempty"`
}

// enrichedOrderDTO embeds the canonical OrderReadModel (reusing its existing
// JSON fields, including the always-nil buyer_nickname the additive SDK lock
// still requires) and overrides Items with the per-item enriched shape, then
// adds the S1-S3 order-level enrichment facts. Every enrichment field is a
// pointer/omitempty slice so an unresolved fact is honestly absent, never a
// fabricated zero value (ADR-17).
type enrichedOrderDTO struct {
	domain.OrderReadModel
	Items                    []enrichedItemDTO    `json:"items"`
	VinculoStatus            string               `json:"vinculo_status"`
	Buyer                    enrichedBuyerDTO     `json:"buyer"`
	ComponentesDesconhecidos []string             `json:"componentes_desconhecidos,omitempty"`
	SLA                      *enrichedSLADTO      `json:"sla,omitempty"`
	DestinoUF                *string              `json:"destino_uf,omitempty"`
	Rastreio                 *enrichedRastreioDTO `json:"rastreio,omitempty"`
	// Bucket is the workflow bucket derived by domain.DeriveOrderBucket. It is
	// ALWAYS derivable (never empty), so unlike the fields above it is
	// required/non-omitempty (ADR-17).
	Bucket domain.OrderBucket `json:"bucket"`
	// RetornoLiquido/MargemPct/Decomposicao/Difal (F01-C1) are the
	// decomposição/DIFAL/retorno seam. Decomposicao and Difal are always
	// present (non-omitempty) because they carry the honest-unknown
	// explanation (componentes_desconhecidos); RetornoLiquido/MargemPct are
	// pointer/omitempty like every other ADR-17 fact.
	RetornoLiquido *float64        `json:"retorno_liquido,omitempty"`
	MargemPct      *float64        `json:"margem_pct,omitempty"`
	Decomposicao   decomposicaoDTO `json:"decomposicao"`
	Difal          difalDTO        `json:"difal"`
}

// mapEnrichedOrder maps an application.EnrichedOrder onto the transport DTO.
// ItemCosts is built by EnrichService in the same iteration order as
// Order.Items (application/enrich_service.go resolveItemCosts, S2), so the
// zip below is positional and order-safe by construction.
func mapEnrichedOrder(e application.EnrichedOrder) enrichedOrderDTO {
	items := make([]enrichedItemDTO, len(e.Order.Items))
	for i, item := range e.Order.Items {
		dto := enrichedItemDTO{MarketplaceOrderItem: item}
		if i < len(e.ItemCosts) {
			dto.CustoUnitario = e.ItemCosts[i].UnitCost
		}
		items[i] = dto
	}

	dto := enrichedOrderDTO{
		OrderReadModel:           e.Order,
		Items:                    items,
		VinculoStatus:            string(e.VinculoStatus),
		Buyer:                    enrichedBuyerDTO{Display: e.Buyer.Display, City: e.Buyer.City, UF: e.Buyer.UF},
		ComponentesDesconhecidos: e.ComponentesDesconhecidos,
		// Bucket derives from provider_status + the order "delivered" tag + the live
		// shipment status. The delivered tag (order_repo read-through) is the robust
		// signal that survives a shipment-lookup failure; the live e.Shipment status
		// refines shipped/delivered vs ready_to_ship when the tag is absent. hasShipment
		// (shipping_id column) stays the pre-shipment faturar/enviar proxy. The SQL
		// summary path shares the tag signal, and the KPI cards derive their counts from
		// this same per-order bucket (FE), so KPI == Lista by construction.
		Bucket:                   domain.DeriveOrderBucket(e.Order.Status, shipmentStatusOf(e.Shipment), e.Order.Tags, e.Order.ShippingID != ""),
		RetornoLiquido:           e.Profitability.RetornoLiquido,
		MargemPct:                e.Profitability.MargemPct,
		Decomposicao:             mapDecomposicao(e.Profitability.Decomposition),
		Difal:                    difalDTO{Amount: e.Profitability.Difal.Amount, UFRoute: e.Profitability.Difal.UFRoute, DueDate: e.Profitability.Difal.DueDate, Paid: e.Profitability.Difal.Paid},
	}
	if e.Shipment != nil {
		dto.SLA = &enrichedSLADTO{Due: e.Shipment.SLADue, Atrasado: e.Shipment.Delayed}
		dto.DestinoUF = e.Shipment.DestinationUF
		dto.Rastreio = &enrichedRastreioDTO{ShipmentID: e.Shipment.ShipmentID, Status: e.Shipment.Status, Substatus: e.Shipment.Substatus}
	}
	return dto
}

// shipmentStatusOf returns the live shipment status, or "" when the shipment
// could not be read (nil enrichment). "" makes DeriveOrderBucket fall back to
// the order tag / provider_status signals — never a fabricated status (ADR-17).
func shipmentStatusOf(s *application.ShipmentEnrichment) string {
	if s == nil {
		return ""
	}
	return s.Status
}

// mapDecomposicao maps an application-layer domain.OrderDecomposition onto
// its transport DTO field-by-field. Flows via mapEnrichedOrder to BOTH the
// /orders list and /orders/{id} detail responses (shared mapper).
func mapDecomposicao(d domain.OrderDecomposition) decomposicaoDTO {
	return decomposicaoDTO{
		Comissao:                 d.Comissao,
		TaxaFixa:                 d.TaxaFixa,
		Frete:                    d.Frete,
		Imposto:                  d.Imposto,
		Difal:                    d.Difal,
		TarifaFull:               d.TarifaFull,
		Custo:                    d.Custo,
		MargemValor:              d.MargemValor,
		MargemPct:                d.MargemPct,
		ComponentesDesconhecidos: d.ComponentesDesconhecidos,
	}
}

func requiredInstallation(values url.Values) (string, error) {
	items, present := values["installation_id"]
	if !present || len(items) == 0 || (len(items) == 1 && strings.TrimSpace(items[0]) == "") {
		return "", ErrInstallationRequired
	}
	if len(items) != 1 {
		return "", &InvalidFilterError{Key: "installation_id"}
	}
	return items[0], nil
}

func writeOrdersError(w http.ResponseWriter, status int, code, message string) {
	writeOrdersErrorDetails(w, status, code, message, map[string]any{})
}

func writeOrderReadError(w http.ResponseWriter, err error) {
	var invalidFilter *InvalidFilterError
	var invalidCursor *InvalidCursorError
	switch {
	case errors.As(err, &invalidFilter):
		writeOrdersErrorDetails(w, http.StatusBadRequest, invalidFilter.Code(), invalidFilter.Error(), invalidFilter.Details())
	case errors.As(err, &invalidCursor), errors.Is(err, ports.ErrInvalidCursor):
		writeOrdersErrorDetails(w, http.StatusBadRequest, "invalid_cursor", "cursor inválido", map[string]any{})
	case errors.Is(err, ports.ErrOrderNotFound):
		writeOrdersErrorDetails(w, http.StatusNotFound, "order_not_found", "order not found", map[string]any{})
	case errors.Is(err, ErrInstallationRequired):
		writeOrdersErrorDetails(w, http.StatusBadRequest, "installation_required", "installation_id é obrigatório", map[string]any{"key": "installation_id"})
	default:
		writeOrdersErrorDetails(w, http.StatusInternalServerError, "ORDERS_INTERNAL_ERROR", "orders read failed", map[string]any{})
	}
}

func writeOrdersErrorDetails(w http.ResponseWriter, status int, code, message string, details map[string]any) {
	httpx.WriteJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
			"details": details,
		},
	})
}

func mapOrdersError(err error) (int, string) {
	msg := strings.TrimSpace(err.Error())
	switch {
	case strings.HasSuffix(msg, "_NOT_FOUND"):
		return http.StatusNotFound, msg
	case strings.HasPrefix(msg, "ORDERS_"), strings.HasPrefix(msg, "INTEGRATIONS_"), strings.HasPrefix(msg, "PRODUCT_LINKS_"):
		return http.StatusBadRequest, msg
	default:
		return http.StatusInternalServerError, "ORDERS_INTERNAL_ERROR"
	}
}

func parsePositiveInt(raw string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
