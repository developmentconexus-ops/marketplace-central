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

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/application"
	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
	"marketplace-central/apps/server_core/internal/platform/apierror"
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

// OrderFaturador marks an order as invoiced ("Foi faturado", goal B).
// Our-DB only — never a Mercado Livre write.
type OrderFaturador interface {
	MarkFaturado(ctx context.Context, installationID, providerOrderID string) error
}

var ErrInstallationRequired = errors.New("installation_required")

type Handler struct {
	importer   OrderImporter
	lister     OrderLister
	readLister OrderReadLister
	enricher   *application.EnrichService
	summary    application.SummaryService
	faturador  OrderFaturador
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

// WithFaturador attaches the POST /orders/{id}/faturado writer. Additive: nil
// faturador simply omits the route.
func (h Handler) WithFaturador(f OrderFaturador) Handler { h.faturador = f; return h }

func (h Handler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("/orders/import", h.handleImport)
	mux.HandleFunc("/orders", h.handleList)
	mux.HandleFunc("/orders/summary", h.handleSummary)
	if h.readLister != nil {
		mux.HandleFunc("/orders/{provider_order_id}", h.handleGet)
		if h.faturador != nil {
			mux.HandleFunc("/orders/{provider_order_id}/faturado", h.handleFaturado)
		}
	}
}

func (h Handler) handleImport(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		apierror.Write(w, http.StatusMethodNotAllowed, "ORDERS_METHOD_NOT_ALLOWED", "method not allowed", nil)
		return
	}
	var req struct {
		InstallationID string `json:"installation_id"`
		Limit          int    `json:"limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierror.Write(w, http.StatusBadRequest, "ORDERS_INVALID_REQUEST", "malformed request body", nil)
		return
	}
	result, err := h.importer.Import(r.Context(), application.ImportOrdersInput{
		InstallationID: req.InstallationID,
		Limit:          req.Limit,
	})
	if err != nil {
		status, code := mapOrdersError(err)
		slog.Error("orders.import", "action", "import", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		apierror.Write(w, status, code, err.Error(), nil)
		return
	}
	slog.Info("orders.import", "action", "import", "result", "200", "count", result.ImportedCount, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h Handler) handleList(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		apierror.Write(w, http.StatusMethodNotAllowed, "ORDERS_METHOD_NOT_ALLOWED", "method not allowed", nil)
		return
	}
	if h.readLister != nil {
		h.handleReadList(w, r)
		return
	}
	if h.lister == nil {
		apierror.Write(w, http.StatusInternalServerError, "ORDERS_INTERNAL_ERROR", "orders list is not configured", nil)
		return
	}
	items, err := h.lister.List(r.Context(), application.ListOrdersInput{
		InstallationID: r.URL.Query().Get("installation_id"),
		Limit:          parsePositiveInt(r.URL.Query().Get("limit"), 20),
	})
	if err != nil {
		status, code := mapOrdersError(err)
		slog.Error("orders.list", "action", "list", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		apierror.Write(w, status, code, err.Error(), nil)
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
		apierror.Write(w, http.StatusMethodNotAllowed, "ORDERS_METHOD_NOT_ALLOWED", "method not allowed", nil)
		return
	}
	if h.readLister == nil {
		apierror.Write(w, http.StatusInternalServerError, "ORDERS_INTERNAL_ERROR", "orders read is not configured", nil)
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
		enriched := h.enricher.EnrichOne(r.Context(), installationID, model)
		httpx.WriteJSON(w, http.StatusOK, mapEnrichedOrder(enriched))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, model)
}

// handleFaturado implements POST /orders/{provider_order_id}/faturado, the
// operator's "Foi faturado" action (goal B). OUR-DB ONLY — this never issues
// a Mercado Livre write; it only records our own faturado_at fact so
// DeriveOrderBucket can move the order from "A FATURAR" to "A ENVIAR".
func (h Handler) handleFaturado(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		apierror.Write(w, http.StatusMethodNotAllowed, "ORDERS_METHOD_NOT_ALLOWED", "method not allowed", nil)
		return
	}
	providerOrderID := strings.TrimSpace(r.PathValue("provider_order_id"))
	var req struct {
		InstallationID string `json:"installation_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierror.Write(w, http.StatusBadRequest, "ORDERS_INVALID_REQUEST", "malformed request body", nil)
		return
	}
	if err := h.faturador.MarkFaturado(r.Context(), req.InstallationID, providerOrderID); err != nil {
		status, code := mapFaturadoError(err)
		slog.Error("orders.faturado", "action", "faturado", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		apierror.Write(w, status, code, err.Error(), nil)
		return
	}
	slog.Info("orders.faturado", "action", "faturado", "result", "204", "provider_order_id", providerOrderID, "duration_ms", time.Since(start).Milliseconds())
	w.WriteHeader(http.StatusNoContent)
}

// mapFaturadoError maps FaturadoService/OrderFaturadoStore errors to an HTTP
// status + error code. ports.ErrOrderNotFound must map to 404: mapOrdersError
// cannot be reused here because its "_NOT_FOUND" suffix match is
// case-sensitive and ports.OrderNotFoundError.Error() renders lowercase
// ("order_not_found"), so it falls through to 500 there.
func mapFaturadoError(err error) (int, string) {
	switch {
	case errors.Is(err, ports.ErrOrderNotFound):
		return http.StatusNotFound, "order_not_found"
	case errors.Is(err, application.ErrFaturadoInvalidInput):
		return http.StatusBadRequest, "ORDERS_INVALID_REQUEST"
	default:
		return http.StatusInternalServerError, "ORDERS_INTERNAL_ERROR"
	}
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
		apierror.Write(w, http.StatusMethodNotAllowed, "ORDERS_METHOD_NOT_ALLOWED", "method not allowed", nil)
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
		apierror.Write(w, http.StatusBadRequest, "unsupported_summary_dimension", "unsupported summary dimension", map[string]any{"by": by})
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
		apierror.Write(w, http.StatusBadRequest, "installation_required", "installation_id é obrigatório", map[string]any{"key": "installation_id"})
	case errors.Is(err, application.ErrSummaryStoreNotConfigured):
		apierror.Write(w, http.StatusServiceUnavailable, "ORDERS_SUMMARY_STORE_NOT_CONFIGURED", "orders summary is not configured", nil)
	default:
		apierror.Write(w, http.StatusInternalServerError, "ORDERS_INTERNAL_ERROR", "orders summary failed", nil)
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
	// Transportadora/UrlRastreio come from the shipment /carrier sub-resource.
	// Additive/omitempty: absent (carrier 404 / not yet dispatched) is honest
	// absence, never a fabricated carrier (ADR-17).
	Transportadora string `json:"transportadora,omitempty"`
	UrlRastreio    string `json:"url_rastreio,omitempty"`
}

// freteRealDTO surfaces the real shipment freight actuals (from the /costs
// sub-resource, mapped by the connectors adapter) — distinct from the modeled
// decomposicao.frete. Every amount is pointer/omitempty: nil means the cost
// could not be honestly sourced (ADR-17), never a fabricated zero.
type freteRealDTO struct {
	Bruto    *float64 `json:"bruto,omitempty"`
	Receiver *float64 `json:"receiver,omitempty"`
	Sender   *float64 `json:"sender,omitempty"`
}

// enrichedItemDTO carries the existing item JSON fields (embedded, so no
// duplication/drift risk) plus the read-time-only per-unit cost.
// CustoUnitario is omitted (never 0) when the cost could not be honestly
// resolved (ADR-17).
type enrichedItemDTO struct {
	domain.MarketplaceOrderItem
	CustoUnitario *float64 `json:"custo_unitario,omitempty"`
	// CustoAproximado marks a cost the ERP source observed at another instant
	// than the order date (single-snapshot source, no cost history); when true
	// CustoObservadoEm carries that instant so the UI can say WHEN, instead of
	// passing the amount off as an exact as-of answer.
	CustoAproximado  bool       `json:"custo_aproximado,omitempty"`
	CustoObservadoEm *time.Time `json:"custo_observado_em,omitempty"`
}

// decomposicaoDTO is the response shape for the order-level cost/fee
// decomposition (F01-C1). Every amount is a pointer/omitempty: nil (absent
// key) means the component could not be honestly sourced (ADR-17).
// ComponentesDesconhecidos is never omitted — it is what explains the "—"s.
type decomposicaoDTO struct {
	Comissao *float64 `json:"comissao,omitempty"`
	TaxaFixa *float64 `json:"taxa_fixa,omitempty"`
	Frete    *float64 `json:"frete,omitempty"`
	Imposto  *float64 `json:"imposto,omitempty"`
	Difal    *float64 `json:"difal,omitempty"`
	// ICMSSaida/PisCofins/RestituicaoST are the P2.b T5 per-item D-41 tax
	// components, replacing Imposto's contribution to the sum end to end
	// (Imposto above stays in the payload — always omitted now, since
	// BuildProfitability no longer populates it — for any consumer still
	// reading the old key; contract-level docs land in Task 6).
	ICMSSaida                *float64 `json:"icms_saida,omitempty"`
	PisCofins                *float64 `json:"pis_cofins,omitempty"`
	TarifaFull               *float64 `json:"tarifa_full,omitempty"`
	Custo                    *float64 `json:"custo,omitempty"`
	RestituicaoST            *float64 `json:"restituicao_st,omitempty"`
	MargemValor              *float64 `json:"margem_valor,omitempty"`
	MargemPct                *float64 `json:"margem_pct,omitempty"`
	ComponentesDesconhecidos []string `json:"componentes_desconhecidos"`
}

// compradorFiscalEnderecoDTO is the buyer's billing address for ERP
// registration. Every field is pointer/omitempty: nil (absent key) means the
// provider did not supply it (ADR-17). UfCodigo/UfNome are the provider's own
// state code + label, rendered verbatim (no UF derivation — that belongs to the
// shipment destination fact, not the fiscal address).
type compradorFiscalEnderecoDTO struct {
	Logradouro *string `json:"logradouro,omitempty"`
	Numero     *string `json:"numero,omitempty"`
	Cidade     *string `json:"cidade,omitempty"`
	UfCodigo   *string `json:"uf_codigo,omitempty"`
	UfNome     *string `json:"uf_nome,omitempty"`
	Cep        *string `json:"cep,omitempty"`
	Pais       *string `json:"pais,omitempty"`
}

// compradorFiscalDTO is the buyer fiscal identity surfaced for ERP
// registration: legal name, an opaque document (type + number) and a billing
// address. DocTipo is the provider's document-type string rendered verbatim —
// never mapped to a CPF/CNPJ enum (ADR-17). Every field is pointer/omitempty:
// a masked/absent provider value is honestly absent, never a fabricated blank.
// The whole block is omitted (nil) when the buyer carries no billing data.
type compradorFiscalDTO struct {
	Nome      *string                     `json:"nome,omitempty"`
	DocTipo   *string                     `json:"doc_tipo,omitempty"`
	DocNumero *string                     `json:"doc_numero,omitempty"`
	Endereco  *compradorFiscalEnderecoDTO `json:"endereco,omitempty"`
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
	DestinoCep               *string              `json:"destino_cep,omitempty"`
	Destinatario             *string              `json:"destinatario,omitempty"`
	FreteReal                *freteRealDTO        `json:"frete_real,omitempty"`
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
	// CompradorFiscal is the buyer's fiscal identity for ERP registration
	// (name, opaque document, billing address). Additive/omitempty: nil (buyer
	// without billing data, or a lookup that could not run) omits the block
	// entirely — never a fabricated identity (ADR-17). The document number is
	// rendered by the caller only; it is never logged (LGPD).
	CompradorFiscal *compradorFiscalDTO `json:"comprador_fiscal,omitempty"`
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
			dto.CustoAproximado = e.ItemCosts[i].Approximate
			dto.CustoObservadoEm = e.ItemCosts[i].ObservedAt
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
		// shipment status + the persisted faturado_at fact. The delivered tag
		// (order_repo read-through) is the robust signal that survives a
		// shipment-lookup failure; the live e.Shipment status refines
		// shipped/delivered when the tag is absent; both win over the faturado
		// gate. faturado_at (not shipping_id presence) now gates the
		// paid/confirmed faturar<->enviar split (goal B) — the operator's
		// explicit "Foi faturado" mark, our-DB only. The SQL summary path shares
		// the same signals, and the KPI cards derive their counts from this same
		// per-order bucket (FE), so KPI == Lista by construction.
		Bucket:         domain.DeriveOrderBucket(e.Order.Status, shipmentStatusOf(e.Shipment), e.Order.Tags, e.Order.FaturadoAt != nil),
		RetornoLiquido: e.Profitability.RetornoLiquido,
		MargemPct:      e.Profitability.MargemPct,
		Decomposicao:   mapDecomposicao(e.Profitability.Decomposition),
		Difal:          difalDTO{Amount: e.Profitability.Difal.Amount, UFRoute: e.Profitability.Difal.UFRoute, DueDate: e.Profitability.Difal.DueDate, Paid: e.Profitability.Difal.Paid},
	}
	if e.Shipment != nil {
		dto.SLA = &enrichedSLADTO{Due: e.Shipment.SLADue, Atrasado: e.Shipment.Delayed}
		dto.DestinoUF = e.Shipment.DestinationUF
		dto.DestinoCep = e.Shipment.DestinationZip
		dto.Destinatario = e.Shipment.ReceiverName
		dto.FreteReal = mapFreteReal(e.Shipment.Costs)
		dto.Rastreio = &enrichedRastreioDTO{
			ShipmentID:     e.Shipment.ShipmentID,
			Status:         e.Shipment.Status,
			Substatus:      e.Shipment.Substatus,
			Transportadora: derefString(e.Shipment.CarrierName),
			UrlRastreio:    derefString(e.Shipment.TrackingURL),
		}
	}
	dto.CompradorFiscal = mapCompradorFiscal(e.BuyerFiscal)
	return dto
}

// mapCompradorFiscal maps the neutral BuyerFiscalInfo onto the transport
// comprador_fiscal block. Returns nil (block omitted) when the enrichment is
// nil — honest absence, never a fabricated identity (ADR-17). doc_tipo is
// carried verbatim; the document number is passed straight through for the
// caller to render and is NEVER logged here (LGPD).
func mapCompradorFiscal(info *connectorsdomain.BuyerFiscalInfo) *compradorFiscalDTO {
	if info == nil {
		return nil
	}
	return &compradorFiscalDTO{
		Nome:      info.Name,
		DocTipo:   info.DocType,
		DocNumero: info.DocNumber,
		Endereco:  mapCompradorFiscalEndereco(info.Address),
	}
}

// mapCompradorFiscalEndereco maps the neutral billing address onto the
// transport endereco block. Returns nil (endereco omitted) when the address is
// nil (the connectors adapter already collapses an all-blank address to nil).
func mapCompradorFiscalEndereco(addr *connectorsdomain.BuyerFiscalAddress) *compradorFiscalEnderecoDTO {
	if addr == nil {
		return nil
	}
	return &compradorFiscalEnderecoDTO{
		Logradouro: addr.StreetName,
		Numero:     addr.StreetNumber,
		Cidade:     addr.City,
		UfCodigo:   addr.StateCode,
		UfNome:     addr.StateName,
		Cep:        addr.ZipCode,
		Pais:       addr.CountryID,
	}
}

// mapFreteReal maps the connectors ShipmentCosts (real freight actuals) onto the
// transport frete_real block. Returns nil when no costs were resolved so the
// field is honestly absent (ADR-17); each amount is nil when its Money is nil.
func mapFreteReal(c *connectorsdomain.ShipmentCosts) *freteRealDTO {
	if c == nil {
		return nil
	}
	bruto := moneyAmountFloat(c.GrossAmount)
	receiver := moneyAmountFloat(c.ReceiverCost)
	sender := moneyAmountFloat(c.SenderCost)
	// All three amounts unresolvable -> omit the block entirely rather than
	// emit an empty {}; honest absence, not a fabricated-empty object (ADR-17).
	if bruto == nil && receiver == nil && sender == nil {
		return nil
	}
	return &freteRealDTO{Bruto: bruto, Receiver: receiver, Sender: sender}
}

// moneyAmountFloat converts a connectors domain.Money (a validated decimal
// string amount) to the *float64 the order DTO uses for every money field. It
// returns nil for a nil Money or an unparseable amount — honest absence, never
// a fabricated zero (ADR-17).
func moneyAmountFloat(m *connectorsdomain.Money) *float64 {
	if m == nil {
		return nil
	}
	v, err := strconv.ParseFloat(m.Amount, 64)
	if err != nil {
		return nil
	}
	return &v
}

// derefString returns the pointed-to string, or "" for a nil pointer (so the
// omitempty DTO field is absent). Used for the additive carrier fields.
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
		ICMSSaida:                d.ICMSSaida,
		PisCofins:                d.PisCofins,
		TarifaFull:               d.TarifaFull,
		Custo:                    d.Custo,
		RestituicaoST:            d.RestituicaoST,
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

func writeOrderReadError(w http.ResponseWriter, err error) {
	var invalidFilter *InvalidFilterError
	var invalidCursor *InvalidCursorError
	switch {
	case errors.As(err, &invalidFilter):
		apierror.Write(w, http.StatusBadRequest, invalidFilter.Code(), invalidFilter.Error(), invalidFilter.Details())
	case errors.As(err, &invalidCursor), errors.Is(err, ports.ErrInvalidCursor):
		apierror.Write(w, http.StatusBadRequest, "invalid_cursor", "cursor inválido", nil)
	case errors.Is(err, ports.ErrOrderNotFound):
		apierror.Write(w, http.StatusNotFound, "order_not_found", "order not found", nil)
	case errors.Is(err, ErrInstallationRequired):
		apierror.Write(w, http.StatusBadRequest, "installation_required", "installation_id é obrigatório", map[string]any{"key": "installation_id"})
	default:
		apierror.Write(w, http.StatusInternalServerError, "ORDERS_INTERNAL_ERROR", "orders read failed", nil)
	}
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
