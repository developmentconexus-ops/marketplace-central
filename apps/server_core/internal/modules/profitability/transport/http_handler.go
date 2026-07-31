package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	profitabilityapp "marketplace-central/apps/server_core/internal/modules/profitability/application"
	profitabilitydomain "marketplace-central/apps/server_core/internal/modules/profitability/domain"
	"marketplace-central/apps/server_core/internal/platform/apierror"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type MarginInputImporter interface {
	ImportMarginInputs(ctx context.Context, input profitabilityapp.ImportMarginInputsInput) (profitabilitydomain.ImportInputsResult, error)
	ListInputs(ctx context.Context, installationID string, limit int) ([]profitabilitydomain.MarginInput, error)
	CreateManualAdjustment(ctx context.Context, input profitabilityapp.CreateManualAdjustmentInput) (profitabilitydomain.ManualAdjustment, error)
	ListAdjustments(ctx context.Context, installationID string, limit int) ([]profitabilitydomain.ManualAdjustment, error)
	CalculateSnapshots(ctx context.Context, installationID string, limit int) (profitabilitydomain.CalculateSnapshotsResult, error)
	ListSnapshots(ctx context.Context, installationID string, limit int) ([]profitabilitydomain.ProfitSnapshot, error)
}

type Handler struct{ service MarginInputImporter }

func NewHandler(service MarginInputImporter) Handler { return Handler{service: service} }

func (h Handler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("/profitability/margin-inputs/import", h.handleImport)
	mux.HandleFunc("/profitability/margin-inputs", h.handleListInputs)
	mux.HandleFunc("/profitability/manual-adjustments", h.handleAdjustments)
	mux.HandleFunc("/profitability/profit-snapshots/calculate", h.handleCalculateSnapshots)
	mux.HandleFunc("/profitability/profit-snapshots", h.handleListSnapshots)
}

func (h Handler) handleImport(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "PROFITABILITY_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	var req struct {
		InstallationID string `json:"installation_id"`
		Limit          int    `json:"limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "PROFITABILITY_INVALID_REQUEST", "malformed request body")
		return
	}
	result, err := h.service.ImportMarginInputs(r.Context(), profitabilityapp.ImportMarginInputsInput{InstallationID: req.InstallationID, Limit: req.Limit})
	if err != nil {
		status, code := mapError(err)
		slog.Error("profitability.margin_inputs.import", "action", "import", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeError(w, status, code, err.Error())
		return
	}
	slog.Info("profitability.margin_inputs.import", "action", "import", "result", "200", "count", result.ImportedCount, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h Handler) handleListInputs(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "PROFITABILITY_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	items, err := h.service.ListInputs(r.Context(), r.URL.Query().Get("installation_id"), parsePositiveInt(r.URL.Query().Get("limit"), 50))
	if err != nil {
		status, code := mapError(err)
		slog.Error("profitability.margin_inputs.list", "action", "list", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeError(w, status, code, err.Error())
		return
	}
	slog.Info("profitability.margin_inputs.list", "action", "list", "result", "200", "count", len(items), "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h Handler) handleAdjustments(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	switch r.Method {
	case http.MethodGet:
		items, err := h.service.ListAdjustments(r.Context(), r.URL.Query().Get("installation_id"), parsePositiveInt(r.URL.Query().Get("limit"), 50))
		if err != nil {
			status, code := mapError(err)
			slog.Error("profitability.manual_adjustments.list", "action", "list", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			writeError(w, status, code, err.Error())
			return
		}
		slog.Info("profitability.manual_adjustments.list", "action", "list", "result", "200", "count", len(items), "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
	case http.MethodPost:
		var req struct {
			InstallationID      string                                       `json:"installation_id"`
			IdempotencyKey      string                                       `json:"idempotency_key"`
			ProviderOrderID     string                                       `json:"provider_order_id"`
			ProviderItemID      string                                       `json:"provider_item_id"`
			ProviderVariationID string                                       `json:"provider_variation_id"`
			Scope               profitabilitydomain.InputScope               `json:"scope"`
			Category            profitabilitydomain.ManualAdjustmentCategory `json:"category"`
			Amount              float64                                      `json:"amount"`
			Currency            string                                       `json:"currency"`
			Reason              string                                       `json:"reason"`
			Actor               profitabilitydomain.ActorMetadata            `json:"actor"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "PROFITABILITY_INVALID_REQUEST", "malformed request body")
			return
		}
		result, err := h.service.CreateManualAdjustment(r.Context(), profitabilityapp.CreateManualAdjustmentInput{
			InstallationID: req.InstallationID, IdempotencyKey: req.IdempotencyKey, ProviderOrderID: req.ProviderOrderID, ProviderItemID: req.ProviderItemID,
			ProviderVariationID: req.ProviderVariationID, Scope: req.Scope, Category: req.Category, Amount: req.Amount,
			Currency: req.Currency, Reason: req.Reason, Actor: req.Actor,
		})
		if err != nil {
			status, code := mapError(err)
			slog.Error("profitability.manual_adjustments.create", "action", "create", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
			writeError(w, status, code, err.Error())
			return
		}
		slog.Info("profitability.manual_adjustments.create", "action", "create", "result", "200", "adjustment_id", result.AdjustmentID, "duration_ms", time.Since(start).Milliseconds())
		httpx.WriteJSON(w, http.StatusOK, result)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "PROFITABILITY_METHOD_NOT_ALLOWED", "method not allowed")
	}
}

func (h Handler) handleCalculateSnapshots(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "PROFITABILITY_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	var req struct {
		InstallationID string `json:"installation_id"`
		Limit          int    `json:"limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "PROFITABILITY_INVALID_REQUEST", "malformed request body")
		return
	}
	result, err := h.service.CalculateSnapshots(r.Context(), req.InstallationID, req.Limit)
	if err != nil {
		status, code := mapError(err)
		slog.Error("profitability.profit_snapshots.calculate", "action", "calculate", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeError(w, status, code, err.Error())
		return
	}
	slog.Info("profitability.profit_snapshots.calculate", "action", "calculate", "result", "200", "count", result.CalculatedCount, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h Handler) handleListSnapshots(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "PROFITABILITY_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	items, err := h.service.ListSnapshots(r.Context(), r.URL.Query().Get("installation_id"), parsePositiveInt(r.URL.Query().Get("limit"), 50))
	if err != nil {
		status, code := mapError(err)
		slog.Error("profitability.profit_snapshots.list", "action", "list", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeError(w, status, code, err.Error())
		return
	}
	slog.Info("profitability.profit_snapshots.list", "action", "list", "result", "200", "count", len(items), "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

// writeError keeps its limit_exceeded branch: that call carries a real value
// (the enforced snapshot-calculation limit) that must land in details, not be
// dropped. Every other status/code pair passes through with empty details.
func writeError(w http.ResponseWriter, status int, code, message string) {
	var details map[string]any
	if status == http.StatusUnprocessableEntity && code == "limit_exceeded" {
		details = map[string]any{"limit": 200}
	}
	apierror.Write(w, status, code, message, details)
}

func mapError(err error) (int, string) {
	msg := strings.TrimSpace(err.Error())
	if strings.HasPrefix(msg, "limit_exceeded") {
		return http.StatusUnprocessableEntity, "limit_exceeded"
	}
	if strings.HasPrefix(msg, "PROFITABILITY_") {
		return http.StatusBadRequest, msg
	}
	return http.StatusInternalServerError, "PROFITABILITY_INTERNAL_ERROR"
}

func parsePositiveInt(raw string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
