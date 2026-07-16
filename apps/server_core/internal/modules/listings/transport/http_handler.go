package transport

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/listings/application"
	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type ListService interface {
	List(context.Context, ports.ListingQuery) (ports.ListingRowPage, error)
	ByProduct(context.Context, ports.ListingGroupQuery) (ports.ListingGroupRowPage, error)
	Get(context.Context, domain.ListingID) (domain.ListingReadModel, []domain.TimelineEvent, error)
	Summary(context.Context, ports.SummaryQuery) (ports.ListingSummaryRow, error)
}
type listingSummaryEnvelope struct {
	Total      int                      `json:"total"`
	Active     int                      `json:"active"`
	Paused     int                      `json:"paused"`
	Exceptions listingSummaryExceptions `json:"exceptions"`
	AsOf       any                      `json:"as_of"`
}
type listingSummaryExceptions struct {
	SyncError            int  `json:"sync_error"`
	Stale                int  `json:"stale"`
	Unlinked             int  `json:"unlinked"`
	BelowMarginWorstCase *int `json:"below_margin_worst_case"`
	MarginUnknown        *int `json:"margin_unknown"`
}

func (h ReadHandler) HandleSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeListError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", "")
		return
	}
	installation, ok := single(r.URL.Query(), "installation_id")
	if !ok || strings.TrimSpace(installation) == "" {
		writeListError(w, http.StatusBadRequest, "installation_required", "installation_id é obrigatório", "installation_id")
		return
	}
	row, err := h.service.Summary(r.Context(), ports.SummaryQuery{InstallationID: installation})
	if err != nil {
		if errors.Is(err, application.ErrInstallationIDRequired) {
			writeListError(w, http.StatusBadRequest, "installation_required", "installation_id é obrigatório", "installation_id")
			return
		}
		if errors.Is(err, application.ErrInstallationNotFound) {
			writeListError(w, http.StatusNotFound, "installation_not_found", "instalação não encontrada", "installation_id")
			return
		}
		writeListError(w, http.StatusServiceUnavailable, "source_unavailable", "fonte de dados indisponível", "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, listingSummaryEnvelope{
		Total: row.Total, Active: row.Active, Paused: row.Paused, AsOf: row.AsOf.UTC(),
		Exceptions: listingSummaryExceptions{SyncError: row.SyncError, Stale: row.Stale, Unlinked: row.Unlinked, BelowMarginWorstCase: row.BelowMarginWorstCase, MarginUnknown: row.MarginUnknown},
	})
}

type ReadHandler struct{ service ListService }

func NewReadHandler(service ListService) ReadHandler { return ReadHandler{service: service} }

type routeClassRegistrar interface {
	RegisterRouteClass(string, httpx.RouteClass)
}

func (h ReadHandler) Register(mux httpx.RouteRegistrar) {
	if registrar, ok := mux.(routeClassRegistrar); ok {
		for _, pattern := range []string{"/listings", "/listings/by-product", "/listings/summary", "/listings/{id}"} {
			registrar.RegisterRouteClass(pattern, httpx.InteractiveRouteClass)
		}
	}
	// POST /listings/refresh belongs here in the later write slice.
	mux.HandleFunc("GET /listings", h.HandleList)
	mux.HandleFunc("GET /listings/by-product", h.HandleByProduct)
	mux.HandleFunc("GET /listings/summary", h.HandleSummary)
	mux.HandleFunc("GET /listings/{id}", h.HandleGet)
}

type listError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}
type listErrorEnvelope struct {
	Error listError `json:"error"`
}
type listingPageEnvelope struct {
	Items      any     `json:"items"`
	NextCursor *string `json:"next_cursor"`
	PageSize   int     `json:"page_size"`
	AsOf       any     `json:"as_of"`
}
type listingGroupPageEnvelope struct {
	Groups     []domain.ListingGroup `json:"groups"`
	NextCursor *string               `json:"next_cursor"`
	PageSize   int                   `json:"page_size"`
	AsOf       any                   `json:"as_of"`
}

type listingDetailEnvelope struct {
	domain.ListingReadModel
	Timeline []domain.TimelineEvent `json:"timeline"`
}

func (h ReadHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeListError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", "")
		return
	}
	values := r.URL.Query()
	installation, ok := single(values, "installation_id")
	if !ok || strings.TrimSpace(installation) == "" {
		writeListError(w, http.StatusBadRequest, "installation_required", "installation_id é obrigatório", "installation_id")
		return
	}
	cursor := ports.ListingCursor{}
	if raw, present := values["cursor"]; present {
		if len(raw) != 1 {
			writeListError(w, 400, "invalid_cursor", "cursor inválido", "cursor")
			return
		}
		decoded, err := ports.DecodeListingCursor(raw[0])
		if err != nil {
			writeListError(w, 400, "invalid_cursor", "cursor inválido", "cursor")
			return
		}
		cursor = decoded
	}
	parsed, err := ParseListingQuery(values)
	if err != nil {
		var invalid *InvalidFilterError
		if errors.As(err, &invalid) {
			writeListError(w, 400, invalid.Code(), invalid.Error(), invalid.Key)
			return
		}
		writeListError(w, 400, "invalid_filter", err.Error(), "")
		return
	}
	page, err := h.service.List(r.Context(), ports.ListingQuery{InstallationID: installation, Filter: parsed.Filter, Q: parsed.Q, Cursor: cursor, Limit: parsed.Limit})
	if err != nil {
		if errors.Is(err, application.ErrInstallationIDRequired) {
			writeListError(w, 400, "installation_required", "installation_id é obrigatório", "installation_id")
			return
		}
		if errors.Is(err, application.ErrInstallationNotFound) {
			writeListError(w, http.StatusNotFound, "installation_not_found", "instalação não encontrada", "installation_id")
			return
		}
		writeListError(w, http.StatusServiceUnavailable, "source_unavailable", "fonte de dados indisponível", "")
		return
	}
	var next *string
	if page.NextCursor != nil {
		encoded, e := page.NextCursor.Encode()
		if e != nil {
			writeListError(w, 500, "internal_error", "internal error", "")
			return
		}
		next = &encoded
	}
	httpx.WriteJSON(w, http.StatusOK, listingPageEnvelope{Items: page.Items, NextCursor: next, PageSize: len(page.Items), AsOf: page.AsOf.UTC()})
}

func (h ReadHandler) HandleByProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeListError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", "")
		return
	}
	values := r.URL.Query()
	installation, ok := single(values, "installation_id")
	if !ok || strings.TrimSpace(installation) == "" {
		writeListError(w, 400, "installation_required", "installation_id é obrigatório", "installation_id")
		return
	}
	cursor := ports.GroupCursor{}
	if raw, present := values["cursor"]; present {
		if len(raw) != 1 {
			writeListError(w, 400, "invalid_cursor", "cursor inválido", "cursor")
			return
		}
		decoded, err := ports.DecodeGroupCursor(raw[0])
		if err != nil {
			writeListError(w, 400, "invalid_cursor", "cursor inválido", "cursor")
			return
		}
		cursor = decoded
	}
	parsed, err := ParseListingQuery(values)
	if err != nil {
		var invalid *InvalidFilterError
		if errors.As(err, &invalid) {
			writeListError(w, 400, invalid.Code(), invalid.Error(), invalid.Key)
		} else {
			writeListError(w, 400, "invalid_filter", err.Error(), "")
		}
		return
	}
	page, err := h.service.ByProduct(r.Context(), ports.ListingGroupQuery{InstallationID: installation, Filter: parsed.Filter, Q: parsed.Q, Cursor: cursor, Limit: parsed.Limit})
	if err != nil {
		if errors.Is(err, application.ErrInstallationIDRequired) {
			writeListError(w, 400, "installation_required", "installation_id é obrigatório", "installation_id")
			return
		}
		if errors.Is(err, application.ErrInstallationNotFound) {
			writeListError(w, 404, "installation_not_found", "instalação não encontrada", "installation_id")
			return
		}
		writeListError(w, 503, "source_unavailable", "fonte de dados indisponível", "")
		return
	}
	var next *string
	if page.NextCursor != nil {
		encoded, err := page.NextCursor.Encode()
		if err != nil {
			writeListError(w, 500, "internal_error", "internal error", "")
			return
		}
		next = &encoded
	}
	httpx.WriteJSON(w, 200, listingGroupPageEnvelope{Groups: page.Groups, NextCursor: next, PageSize: len(page.Groups), AsOf: page.AsOf.UTC()})
}

func (h ReadHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeListError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed", "")
		return
	}
	rawID := r.PathValue("id")
	if rawID == "" {
		rawID = r.PathValue("listing_id")
	}
	id, err := domain.ParseListingID(rawID)
	if err != nil {
		writeListError(w, http.StatusNotFound, "listing_not_found", "anúncio não encontrado", "listing_id")
		return
	}
	model, timeline, err := h.service.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrListingNotFound) {
			writeListError(w, http.StatusNotFound, "listing_not_found", "anúncio não encontrado", "listing_id")
			return
		}
		writeListError(w, http.StatusInternalServerError, "internal_error", "internal error", "")
		return
	}
	if timeline == nil {
		timeline = []domain.TimelineEvent{}
	}
	httpx.WriteJSON(w, http.StatusOK, listingDetailEnvelope{ListingReadModel: model, Timeline: timeline})
}
func single(v map[string][]string, key string) (string, bool) {
	values, ok := v[key]
	return func() (string, bool) {
		if !ok || len(values) != 1 {
			return "", false
		}
		return values[0], true
	}()
}
func writeListError(w http.ResponseWriter, status int, code, message, key string) {
	details := map[string]any{}
	if key != "" {
		details["key"] = key
	}
	httpx.WriteJSON(w, status, listErrorEnvelope{Error: listError{Code: code, Message: message, Details: details}})
}
