package transport

import (
	"context"
	"errors"
	"net/http"
	"time"

	"marketplace-central/apps/server_core/internal/modules/market/application"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/platform/apierror"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type ReadService interface {
	ListObservations(context.Context, []string) (application.ObservationReadResult, error)
	ListReferences(context.Context, []string) (application.ReferenceReadResult, error)
}

// CollectionService is the write side of on-demand price-intel collection
// (F-04-S5). Single-product synchronous by design (D-F4-p): no batch POST.
type CollectionService interface {
	Collect(ctx context.Context, codprod string) (application.CollectionSummary, error)
}

// EvidenceReader is the read side of assembled price-intel evidence
// (F-04-S5), matching ports.EvidenceReader's frozen signature.
type EvidenceReader interface {
	Signals(ctx context.Context, listingIDs []string) ([]domain.CompetitiveSignal, error)
	Aggregates(ctx context.Context, codprods []string) ([]domain.MarketAggregate, error)
	Verdicts(ctx context.Context, codprods []string) ([]domain.Verdict, error)
}

type Handler struct {
	service     ReadService
	collections CollectionService
	evidence    EvidenceReader
}

func NewHandler(service ReadService) Handler {
	return Handler{service: service}
}

// NewHandlerWithCollections extends NewHandler with the F-04-S5 collection
// and evidence capabilities. Additive: NewHandler's single-argument
// signature is unchanged so the existing composition/root.go call site keeps
// compiling untouched.
func NewHandlerWithCollections(service ReadService, collections CollectionService, evidence EvidenceReader) Handler {
	return Handler{service: service, collections: collections, evidence: evidence}
}

func (h Handler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("GET /market/observations", h.handleObservations)
	mux.HandleFunc("GET /market/references", h.handleReferences)
	if h.collections != nil {
		mux.HandleFunc("POST /market/collections", h.handleCollect)
	}
	if h.evidence != nil {
		mux.HandleFunc("GET /market/signals", h.handleSignals)
		mux.HandleFunc("GET /market/aggregates", h.handleAggregates)
		mux.HandleFunc("GET /market/verdicts", h.handleVerdicts)
	}
}

type marketObservationsEnvelope struct {
	Items []domain.MarketObservation `json:"items"`
	AsOf  time.Time                  `json:"as_of"`
}

type marketReferencesEnvelope struct {
	Items []domain.MarketReference `json:"items"`
	AsOf  time.Time                `json:"as_of"`
}

func (h Handler) handleObservations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		apierror.Write(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", map[string]any{})
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
		apierror.Write(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", map[string]any{})
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
		apierror.Write(w, http.StatusBadRequest, invalidFilter.Code(), invalidFilter.Error(), invalidFilter.Details())
	case errors.As(err, &installationRequired):
		apierror.Write(w, http.StatusBadRequest, installationRequired.Code(), "installation_id é obrigatório", installationRequired.Details())
	default:
		apierror.Write(w, http.StatusBadRequest, "invalid_filter", "filtro inválido", map[string]any{})
	}
}

func writeMarketServiceError(w http.ResponseWriter, err error) {
	var tooMany *application.TooManyIDsError
	if errors.As(err, &tooMany) || errors.Is(err, application.ErrTooManyIDs) {
		apierror.Write(w, http.StatusUnprocessableEntity, application.ErrTooManyIDs.Error(), "too many ids", map[string]any{})
		return
	}
	apierror.Write(w, http.StatusInternalServerError, "internal_error", "internal error", map[string]any{})
}

// writeMarketError forwards to apierror.Write. Kept (not inlined at every
// call site) because verdict_handler.go in this package calls it too and is
// outside this slice's write_set — deleting it would break that file's
// build. See slice report for the alternatives-considered note.
func writeMarketError(w http.ResponseWriter, status int, code, message string, details map[string]any) {
	apierror.Write(w, status, code, message, details)
}
