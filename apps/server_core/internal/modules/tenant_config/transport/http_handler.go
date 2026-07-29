// Package transport exposes the tenant configuration over HTTP.
//
// /config/active-source: GET reads the tenant's current selection (fail-closed
// 400 when unset), PUT sets it. source_kind is always server-derived
// (tenant_config.DefaultKind) and never accepted from the request body (ADR-17).
//
// /config/sellable-assortment: GET reads the three stored toggles, PUT writes
// them. A body missing any of the three is a 400 — an absent boolean would
// otherwise decode to false and be written as a rule no operator authored.
// A tenant with no config row answers 500 here, unlike active-source's 400:
// the published contract offers no other status on this surface, and the
// honest status for that state is an open question with the hub.
package transport

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"marketplace-central/apps/server_core/internal/modules/tenant_config"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

// Store is the narrow port this handler depends on. *tenant_config.Repository
// satisfies it (compile-checked below); http_handler_test.go supplies a fake
// for pure-unit coverage — two consumers justify the port.
type Store interface {
	Get(ctx context.Context, tenantID string) (tenant_config.Config, error)
	Set(ctx context.Context, cfg tenant_config.Config) error
	SetSellableAssortment(ctx context.Context, tenantID string, a tenant_config.SellableAssortment) error
}

var _ Store = (*tenant_config.Repository)(nil)

type Handler struct {
	store    Store
	tenantID string
}

func NewHandler(store Store, tenantID string) Handler {
	return Handler{store: store, tenantID: tenantID}
}

type routeClassRegistrar interface {
	RegisterRouteClass(string, httpx.RouteClass)
}

func (h Handler) Register(mux httpx.RouteRegistrar) {
	if registrar, ok := mux.(routeClassRegistrar); ok {
		registrar.RegisterRouteClass("/config/active-source", httpx.InteractiveRouteClass)
		registrar.RegisterRouteClass("/config/sellable-assortment", httpx.InteractiveRouteClass)
	}
	mux.HandleFunc("GET /config/active-source", h.handleGet)
	mux.HandleFunc("PUT /config/active-source", h.handlePut)
	mux.HandleFunc("GET /config/sellable-assortment", h.handleGetSellableAssortment)
	mux.HandleFunc("PUT /config/sellable-assortment", h.handlePutSellableAssortment)
}

func (h Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.store.Get(r.Context(), h.tenantID)
	if err != nil {
		if errors.Is(err, tenant_config.ErrUnknownActiveSource) {
			writeError(w, http.StatusBadRequest, "unknown_erp_source", "")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, newActiveSourceResponse(cfg))
}

type putActiveSourceRequest struct {
	ActiveSource string  `json:"active_source"`
	SetBy        *string `json:"set_by"`
}

func (h Handler) handlePut(w http.ResponseWriter, r *http.Request) {
	var body putActiveSourceRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_body", "")
		return
	}
	setBy := ""
	if body.SetBy != nil {
		setBy = *body.SetBy
	}
	cfg := tenant_config.Config{
		TenantID: h.tenantID,
		Source:   tenant_config.ActiveSource(body.ActiveSource),
		SetBy:    setBy,
	}
	if err := h.store.Set(r.Context(), cfg); err != nil {
		if errors.Is(err, tenant_config.ErrInvalidActiveSource) {
			writeError(w, http.StatusBadRequest, "invalid_active_source", "")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "")
		return
	}
	stored, err := h.store.Get(r.Context(), h.tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, newActiveSourceResponse(stored))
}

func (h Handler) handleGetSellableAssortment(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.store.Get(r.Context(), h.tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, newSellableAssortmentResponse(cfg.SellableAssortment))
}

type putSellableAssortmentRequest struct {
	OnlyRevenda           *bool `json:"only_revenda"`
	OnlyEmEstoque         *bool `json:"only_em_estoque"`
	OnlyEcommerceEligible *bool `json:"only_ecommerce_eligible"`
}

func (h Handler) handlePutSellableAssortment(w http.ResponseWriter, r *http.Request) {
	var body putSellableAssortmentRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil ||
		decoder.Decode(&struct{}{}) != io.EOF ||
		body.OnlyRevenda == nil || body.OnlyEmEstoque == nil || body.OnlyEcommerceEligible == nil {
		writeError(w, http.StatusBadRequest, "invalid_body", "")
		return
	}
	assortment := tenant_config.SellableAssortment{
		OnlyRevenda:           *body.OnlyRevenda,
		OnlyEmEstoque:         *body.OnlyEmEstoque,
		OnlyEcommerceEligible: *body.OnlyEcommerceEligible,
	}
	if err := h.store.SetSellableAssortment(r.Context(), h.tenantID, assortment); err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "")
		return
	}
	stored, err := h.store.Get(r.Context(), h.tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, newSellableAssortmentResponse(stored.SellableAssortment))
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	payload := map[string]string{"error": code}
	if message != "" {
		payload["detail"] = message
	}
	httpx.WriteJSON(w, status, payload)
}

type activeSourceResponse struct {
	ActiveSource string `json:"active_source"`
	SourceKind   string `json:"source_kind"`
	SetAt        string `json:"set_at"`
	SetBy        string `json:"set_by"`
}

func newActiveSourceResponse(cfg tenant_config.Config) activeSourceResponse {
	return activeSourceResponse{
		ActiveSource: string(cfg.Source),
		SourceKind:   string(cfg.Kind),
		SetAt:        cfg.SetAt.UTC().Format(time.RFC3339),
		SetBy:        cfg.SetBy,
	}
}

type sellableAssortmentResponse struct {
	OnlyRevenda           bool `json:"only_revenda"`
	OnlyEmEstoque         bool `json:"only_em_estoque"`
	OnlyEcommerceEligible bool `json:"only_ecommerce_eligible"`
}

func newSellableAssortmentResponse(a tenant_config.SellableAssortment) sellableAssortmentResponse {
	return sellableAssortmentResponse{
		OnlyRevenda:           a.OnlyRevenda,
		OnlyEmEstoque:         a.OnlyEmEstoque,
		OnlyEcommerceEligible: a.OnlyEcommerceEligible,
	}
}
