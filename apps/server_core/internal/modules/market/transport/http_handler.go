package transport

import (
	"context"
	"errors"
	"net/http"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/application"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type ReadService interface {
	ListObservations(context.Context, []string) (application.ObservationReadResult, error)
	ListReferences(context.Context, []string) (application.ReferenceReadResult, error)
}

type Handler struct {
	service ReadService
}

func NewHandler(service ReadService) Handler {
	return Handler{service: service}
}

func (h Handler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("GET /market/observations", h.handleObservations)
	mux.HandleFunc("GET /market/references", h.handleReferences)
}

type marketObservationsEnvelope struct {
	Items []domain.MarketObservation `json:"items"`
	AsOf  time.Time                  `json:"as_of"`
}

type marketReferencesEnvelope struct {
	Items []domain.MarketReference `json:"items"`
	AsOf  time.Time                `json:"as_of"`
}

type apiError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

type apiErrorResponse struct {
	Error apiError `json:"error"`
}

func (h Handler) handleObservations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeMarketError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", map[string]any{})
		return
	}
	query, err := ParseObservationQuery(r.URL.Query())
	if err != nil {
		writeMarketQueryError(w, err)
		return
	}
	result, err := h.service.ListObservations(r.Context(), query.ListingIDs)
	if err != nil {
		writeMarketServiceError(w, err)
		return
	}
	if result.Items == nil {
		result.Items = []domain.MarketObservation{}
	}
	httpx.WriteJSON(w, http.StatusOK, marketObservationsEnvelope{Items: result.Items, AsOf: result.AsOf.UTC()})
}

func (h Handler) handleReferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeMarketError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", map[string]any{})
		return
	}
	query, err := ParseReferenceQuery(r.URL.Query())
	if err != nil {
		writeMarketQueryError(w, err)
		return
	}
	result, err := h.service.ListReferences(r.Context(), query.ProductIDs)
	if err != nil {
		writeMarketServiceError(w, err)
		return
	}
	if result.Items == nil {
		result.Items = []domain.MarketReference{}
	}
	httpx.WriteJSON(w, http.StatusOK, marketReferencesEnvelope{Items: result.Items, AsOf: result.AsOf.UTC()})
}

func writeMarketQueryError(w http.ResponseWriter, err error) {
	var invalidFilter *InvalidFilterError
	var installationRequired *InstallationRequiredError
	switch {
	case errors.As(err, &invalidFilter):
		writeMarketError(w, http.StatusBadRequest, invalidFilter.Code(), invalidFilter.Error(), invalidFilter.Details())
	case errors.As(err, &installationRequired):
		writeMarketError(w, http.StatusBadRequest, installationRequired.Code(), "installation_id é obrigatório", installationRequired.Details())
	default:
		writeMarketError(w, http.StatusBadRequest, "invalid_filter", "filtro inválido", map[string]any{})
	}
}

func writeMarketServiceError(w http.ResponseWriter, err error) {
	var tooMany *application.TooManyIDsError
	if errors.As(err, &tooMany) || errors.Is(err, application.ErrTooManyIDs) {
		writeMarketError(w, http.StatusUnprocessableEntity, application.ErrTooManyIDs.Error(), "too many ids", map[string]any{})
		return
	}
	writeMarketError(w, http.StatusInternalServerError, "internal_error", "internal error", map[string]any{})
}

func writeMarketError(w http.ResponseWriter, status int, code, message string, details map[string]any) {
	httpx.WriteJSON(w, status, apiErrorResponse{Error: apiError{Code: code, Message: message, Details: details}})
}
