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

func (h Handler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("/orders/import", h.handleImport)
	mux.HandleFunc("/orders", h.handleList)
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
}

// enrichedItemDTO carries the existing item JSON fields (embedded, so no
// duplication/drift risk) plus the read-time-only per-unit cost.
// CustoUnitario is omitted (never 0) when the cost could not be honestly
// resolved (ADR-17).
type enrichedItemDTO struct {
	domain.MarketplaceOrderItem
	CustoUnitario *float64 `json:"custo_unitario,omitempty"`
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
	}
	if e.Shipment != nil {
		dto.SLA = &enrichedSLADTO{Due: e.Shipment.SLADue, Atrasado: e.Shipment.Delayed}
		dto.DestinoUF = e.Shipment.DestinationUF
		dto.Rastreio = &enrichedRastreioDTO{ShipmentID: e.Shipment.ShipmentID, Status: e.Shipment.Status}
	}
	return dto
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
