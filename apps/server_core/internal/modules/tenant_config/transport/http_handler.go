// Package transport exposes the tenant active-source config over HTTP:
// GET reads the tenant's current selection (fail-closed 400 when unset),
// PUT sets it. source_kind is always server-derived (tenant_config.DefaultKind)
// and never accepted from the request body (ADR-17).
package transport

import (
	"context"
	"encoding/json"
	"errors"
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
	}
	mux.HandleFunc("GET /config/active-source", h.handleGet)
	mux.HandleFunc("PUT /config/active-source", h.handlePut)
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
