package transport

import (
	"context"
	"errors"
	"net/http"
	"time"

	"marketplace-central/apps/server_core/internal/modules/integrations/ports"
	"marketplace-central/apps/server_core/internal/platform/apierror"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type RunReadService interface {
	ListRuns(context.Context, ports.RunListQuery) (ports.RunPage, error)
}

type RunReadHandler struct {
	service RunReadService
}

type runRouteClassRegistrar interface {
	RegisterRouteClass(string, httpx.RouteClass)
}

func NewRunReadHandler(service RunReadService) RunReadHandler {
	return RunReadHandler{service: service}
}

func (h RunReadHandler) Register(mux httpx.RouteRegistrar) {
	if registrar, ok := mux.(runRouteClassRegistrar); ok {
		registrar.RegisterRouteClass("/sync/runs", httpx.InteractiveRouteClass)
	}
	mux.HandleFunc("GET /sync/runs", h.HandleList)
}

type runPageEnvelope struct {
	Items      []runReadItem `json:"items"`
	NextCursor *string       `json:"next_cursor"`
	PageSize   int           `json:"page_size"`
	AsOf       time.Time     `json:"as_of"`
}

type runReadItem struct {
	OperationRunID      string     `json:"operation_run_id"`
	InstallationID      string     `json:"installation_id"`
	Module              string     `json:"module"`
	Status              string     `json:"status"`
	ResultCode          string     `json:"result_code"`
	FailureCode         string     `json:"failure_code"`
	TranslatedErrorCode string     `json:"translated_error_code"`
	AttemptCount        int        `json:"attempt_count"`
	DurationMs          int64      `json:"duration_ms"`
	StartedAt           *time.Time `json:"started_at"`
	FinishedAt          *time.Time `json:"finished_at"`
}

func (h RunReadHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeRunReadError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", "")
		return
	}

	query, err := ParseRunQuery(r.URL.Query())
	if err != nil {
		writeRunQueryError(w, err)
		return
	}
	page, err := h.service.ListRuns(r.Context(), query)
	if err != nil {
		if errors.Is(err, ports.ErrInvalidRunCursor) {
			writeRunReadError(w, http.StatusBadRequest, "invalid_cursor", "cursor inválido", "cursor")
			return
		}
		if errors.Is(err, ports.ErrInvalidRunListQuery) {
			writeRunReadError(w, http.StatusBadRequest, "invalid_filter", "filtro inválido", "")
			return
		}
		writeRunReadError(w, http.StatusInternalServerError, "internal_error", "internal error", "")
		return
	}

	items := make([]runReadItem, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, runReadItem{
			OperationRunID: item.OperationRunID, InstallationID: item.InstallationID, Module: item.Module,
			Status: string(item.Status), ResultCode: item.ResultCode, FailureCode: item.FailureCode,
			TranslatedErrorCode: item.TranslatedErrorCode, AttemptCount: item.AttemptCount, DurationMs: item.DurationMs,
			StartedAt: item.StartedAt, FinishedAt: item.FinishedAt,
		})
	}
	var nextCursor *string
	if page.NextCursor != nil {
		encoded, err := page.NextCursor.Encode()
		if err != nil {
			writeRunReadError(w, http.StatusInternalServerError, "internal_error", "internal error", "")
			return
		}
		nextCursor = &encoded
	}

	httpx.WriteJSON(w, http.StatusOK, runPageEnvelope{
		Items: items, NextCursor: nextCursor, PageSize: len(items), AsOf: time.Now().UTC(),
	})
}

func writeRunQueryError(w http.ResponseWriter, err error) {
	var invalidFilter *InvalidFilterError
	var invalidCursor *InvalidCursorError
	var installationRequired *InstallationRequiredError
	switch {
	case errors.As(err, &invalidFilter):
		writeRunReadError(w, http.StatusBadRequest, invalidFilter.Code(), invalidFilter.Error(), invalidFilter.Key)
	case errors.As(err, &invalidCursor):
		writeRunReadError(w, http.StatusBadRequest, invalidCursor.Code(), "cursor inválido", "cursor")
	case errors.As(err, &installationRequired):
		writeRunReadError(w, http.StatusBadRequest, installationRequired.Code(), "installation_id é obrigatório", "installation_id")
	default:
		writeRunReadError(w, http.StatusBadRequest, "invalid_filter", "filtro inválido", "")
	}
}

func writeRunReadError(w http.ResponseWriter, status int, code, message, key string) {
	details := map[string]any{}
	if key != "" {
		details["key"] = key
	}
	apierror.Write(w, status, code, message, details)
}
