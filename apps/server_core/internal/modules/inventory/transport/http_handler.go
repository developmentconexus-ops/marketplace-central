package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/inventory/application"
	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type StockRiskLister interface {
	List(ctx context.Context, input application.ListStockRisksInput) ([]domain.StockRiskListItem, error)
}

type ManualActionApplier interface {
	Apply(ctx context.Context, input application.ApplyManualActionInput) (application.ApplyManualActionResult, error)
}

type Handler struct {
	risks   StockRiskLister
	actions ManualActionApplier
}

func NewHandler(risks StockRiskLister, actions ManualActionApplier) Handler {
	return Handler{risks: risks, actions: actions}
}

func (h Handler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("/inventory/stock-risks", h.handleStockRisks)
	mux.HandleFunc("/inventory/stock-actions/manual-apply", h.handleManualApply)
}

type apiError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

func writeInventoryError(w http.ResponseWriter, status int, code, message string) {
	httpx.WriteJSON(w, status, map[string]any{
		"error": apiError{
			Code:    code,
			Message: message,
			Details: map[string]any{},
		},
	})
}

func (h Handler) handleStockRisks(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeInventoryError(w, http.StatusMethodNotAllowed, "INVENTORY_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	items, err := h.risks.List(r.Context(), application.ListStockRisksInput{
		InstallationID: r.URL.Query().Get("installation_id"),
		State:          r.URL.Query().Get("state"),
		LinkState:      r.URL.Query().Get("link_state"),
		Actionability:  r.URL.Query().Get("actionability"),
		Limit:          parsePositiveInt(r.URL.Query().Get("limit"), 50),
	})
	if err != nil {
		status, code := mapInventoryError(err)
		slog.Error("inventory.stock_risks", "action", "list", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeInventoryError(w, status, code, err.Error())
		return
	}
	slog.Info("inventory.stock_risks", "action", "list", "result", "200", "count", len(items), "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h Handler) handleManualApply(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeInventoryError(w, http.StatusMethodNotAllowed, "INVENTORY_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	var req struct {
		StockActionID       string `json:"stock_action_id"`
		InstallationID      string `json:"installation_id"`
		ProviderItemID      string `json:"provider_item_id"`
		ProviderVariationID string `json:"provider_variation_id"`
		RequestedQuantity   int    `json:"requested_quantity"`
		Reason              string `json:"reason"`
		Approval            struct {
			Approved   bool      `json:"approved"`
			ApprovedAt time.Time `json:"approved_at"`
			Operator   struct {
				ActorType string `json:"actor_type"`
				ActorID   string `json:"actor_id"`
				ActorName string `json:"actor_name"`
			} `json:"operator"`
		} `json:"approval"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeInventoryError(w, http.StatusBadRequest, "INVENTORY_INVALID_REQUEST", "malformed request body")
		return
	}
	result, err := h.actions.Apply(r.Context(), application.ApplyManualActionInput{
		ActionID:            req.StockActionID,
		InstallationID:      req.InstallationID,
		ProviderItemID:      req.ProviderItemID,
		ProviderVariationID: req.ProviderVariationID,
		RequestedQuantity:   req.RequestedQuantity,
		Reason:              req.Reason,
		Approval: domain.ManualApproval{
			Approved:   req.Approval.Approved,
			ApprovedAt: req.Approval.ApprovedAt,
			Operator: domain.StockActionOperator{
				ActorType: req.Approval.Operator.ActorType,
				ActorID:   req.Approval.Operator.ActorID,
				ActorName: req.Approval.Operator.ActorName,
			},
		},
	})
	if err != nil {
		status, code := mapInventoryError(err)
		slog.Error("inventory.stock_actions.manual_apply", "action", "apply", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeInventoryError(w, status, code, err.Error())
		return
	}
	slog.Info("inventory.stock_actions.manual_apply", "action", "apply", "result", "200", "stock_action_id", result.Action.ActionID, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, result)
}

func mapInventoryError(err error) (int, string) {
	msg := strings.TrimSpace(err.Error())
	switch {
	case strings.HasSuffix(msg, "_NOT_FOUND"):
		return http.StatusNotFound, msg
	case strings.HasPrefix(msg, "INVENTORY_"), strings.HasPrefix(msg, "INTEGRATIONS_"), strings.HasPrefix(msg, "PRODUCT_LINKS_"):
		return http.StatusBadRequest, msg
	default:
		return http.StatusInternalServerError, "INVENTORY_INTERNAL_ERROR"
	}
}

func parsePositiveInt(raw string, fallback int) int {
	value := strings.TrimSpace(raw)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
