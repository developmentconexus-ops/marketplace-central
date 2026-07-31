package transport

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/dashboard/domain"
	"marketplace-central/apps/server_core/internal/platform/apierror"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type SummaryService interface {
	Summary(context.Context, string) (domain.Summary, error)
}

type Handler struct {
	service SummaryService
}

func NewHandler(service SummaryService) Handler {
	return Handler{service: service}
}

func (h Handler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("GET /dashboard/summary", h.HandleSummary)
}

func (h Handler) HandleSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeDashboardError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", "")
		return
	}

	installationID := strings.TrimSpace(r.URL.Query().Get("installation_id"))
	if installationID == "" {
		writeDashboardError(w, http.StatusBadRequest, "installation_required", "installation_id é obrigatório", "installation_id")
		return
	}

	summary, err := h.service.Summary(r.Context(), installationID)
	if err != nil {
		if errors.Is(err, domain.ErrInstallationNotFound) {
			writeDashboardError(w, http.StatusNotFound, "installation_not_found", "installation not found", "")
			return
		}
		writeDashboardError(w, http.StatusInternalServerError, "internal_error", "internal error", "")
		return
	}
	if summary.Degraded == nil {
		summary.Degraded = []string{}
	}

	httpx.WriteJSON(w, http.StatusOK, summary)
}

func writeDashboardError(w http.ResponseWriter, status int, code, message, key string) {
	details := map[string]any{}
	if key != "" {
		details["key"] = key
	}
	apierror.Write(w, status, code, message, details)
}
