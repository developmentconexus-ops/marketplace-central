package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/application"
	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type OrderImporter interface {
	Import(ctx context.Context, input application.ImportOrdersInput) (domain.ImportResult, error)
}

type OrderLister interface {
	List(ctx context.Context, input application.ListOrdersInput) ([]domain.MarketplaceOrder, error)
}

type Handler struct {
	importer OrderImporter
	lister   OrderLister
}

func NewHandler(importer OrderImporter, lister OrderLister) Handler {
	return Handler{importer: importer, lister: lister}
}

func (h Handler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("/orders/import", h.handleImport)
	mux.HandleFunc("/orders", h.handleList)
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

func writeOrdersError(w http.ResponseWriter, status int, code, message string) {
	httpx.WriteJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
			"details": map[string]any{},
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
