package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/application"
	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

// errorEnvelopeShape decodes the apierror.Write envelope for assertions that
// need the code/details fields (json.Unmarshal only needs field names to
// match; no import of the unexported apierror types is required).
type errorEnvelopeShape struct {
	Error struct {
		Code    string         `json:"code"`
		Message string         `json:"message"`
		Details map[string]any `json:"details"`
	} `json:"error"`
}

func trimJSON(body string) string {
	var value any
	if err := json.Unmarshal([]byte(body), &value); err != nil {
		return body
	}
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

func TestListingsErrorEnvelopeWholeBody(t *testing.T) {
	h := NewReadHandler(&fakeListService{})
	r := httptest.NewRequest(http.MethodGet, "/listings?installation_id=i&cursor=bad", nil)
	w := httptest.NewRecorder()
	h.HandleList(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	want := `{"error":{"code":"invalid_cursor","message":"cursor inválido","details":{"key":"cursor"}}}`
	if trimJSON(w.Body.String()) != trimJSON(want) {
		t.Fatalf("body = %s, want %s", trimJSON(w.Body.String()), trimJSON(want))
	}
}

type refreshServiceStub struct {
	id  string
	err error
}

func (s refreshServiceStub) Start(context.Context, string) (string, error) { return s.id, s.err }

func TestRefreshHandlerContract(t *testing.T) {
	tests := []struct {
		name, body string
		service    refreshServiceStub
		wantStatus int
		wantBody   string
	}{
		{"missing", `{}`, refreshServiceStub{err: application.ErrInstallationIDRequired}, 400, `"code":"installation_required"`},
		{"blank", `{"installation_id":" "}`, refreshServiceStub{err: application.ErrInstallationIDRequired}, 400, `"code":"installation_required"`},
		{"unknown", `{"installation_id":"missing"}`, refreshServiceStub{err: application.ErrInstallationNotFound}, 404, `"code":"installation_not_found"`},
		{"accepted", `{"installation_id":"inst"}`, refreshServiceStub{id: "op_new"}, 202, `{"operation_run_id":"op_new"}`},
		{"active", `{"installation_id":"inst"}`, refreshServiceStub{err: &application.RefreshInProgressError{OperationRunID: "op_active"}}, 409, `"operation_run_id":"op_active"`},
		{"malformed", `{"installation_id":`, refreshServiceStub{}, 400, `"code":"invalid_request"`},
		{"trailing", `{"installation_id":"inst"}{}`, refreshServiceStub{}, 400, `"code":"invalid_request"`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/listings/refresh", bytes.NewBufferString(tc.body))
			NewRefreshHandler(tc.service).HandleRefresh(rec, req)
			if rec.Code != tc.wantStatus || !strings.Contains(strings.TrimSpace(rec.Body.String()), tc.wantBody) {
				t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestRefreshHandlerRegistersPostOnly(t *testing.T) {
	mux := http.NewServeMux()
	NewRefreshHandler(refreshServiceStub{id: "op"}).Register(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/listings/refresh", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d", rec.Code)
	}
}

type fakeListService struct {
	page       ports.ListingRowPage
	groupPage  ports.ListingGroupRowPage
	calls      int
	model      domain.ListingReadModel
	timeline   []domain.TimelineEvent
	getErr     error
	summary    ports.ListingSummaryRow
	summaryErr error
}

func TestReadHandlerRegisterMapsReadRoutes(t *testing.T) {
	mux := httpx.NewRouteClassMux()
	NewReadHandler(&fakeListService{}).Register(mux)

	for _, tc := range []struct {
		path   string
		status int
	}{
		{path: "/listings", status: http.StatusBadRequest},
		{path: "/listings/by-product", status: http.StatusBadRequest},
		{path: "/listings/summary", status: http.StatusBadRequest},
		{path: "/listings/not-a-listing-id", status: http.StatusNotFound},
	} {
		t.Run(tc.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if recorder.Code != tc.status {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, tc.status, recorder.Body.String())
			}
		})
	}
}

func (f *fakeListService) Summary(context.Context, ports.SummaryQuery) (ports.ListingSummaryRow, error) {
	return f.summary, f.summaryErr
}

func (f *fakeListService) List(context.Context, ports.ListingQuery) (ports.ListingRowPage, error) {
	f.calls++
	return f.page, nil
}
func (f *fakeListService) ByProduct(context.Context, ports.ListingGroupQuery) (ports.ListingGroupRowPage, error) {
	return f.groupPage, nil
}
func (f *fakeListService) Get(context.Context, domain.ListingID) (domain.ListingReadModel, []domain.TimelineEvent, error) {
	return f.model, f.timeline, f.getErr
}
func TestListHandlerValidationErrors(t *testing.T) {
	h := NewReadHandler(&fakeListService{})
	for _, tc := range []struct{ url, key, code string }{{"/listings", "installation_id", "installation_required"}, {"/listings?installation_id=i&filter.nope=x", "nope", "invalid_filter"}, {"/listings?installation_id=i&cursor=bad", "cursor", "invalid_cursor"}} {
		r := httptest.NewRequest(http.MethodGet, tc.url, nil)
		w := httptest.NewRecorder()
		h.HandleList(w, r)
		if w.Code != 400 {
			t.Fatalf("%s status=%d", tc.url, w.Code)
		}
		var body struct {
			Error struct {
				Code    string         `json:"code"`
				Details map[string]any `json:"details"`
			} `json:"error"`
		}
		json.Unmarshal(w.Body.Bytes(), &body)
		if body.Error.Code != tc.code || body.Error.Details["key"] != tc.key {
			t.Fatalf("%s body=%s", tc.url, w.Body.String())
		}
	}
}
func TestListHandlerIC02Envelope(t *testing.T) {
	at := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)
	svc := &fakeListService{page: ports.ListingRowPage{Items: []domain.ListingReadModel{{ListingID: "i~x~-", InstallationID: "i", Provider: "mercadolivre", ProviderListingID: "x", Title: "X", Status: domain.ListingStatusActive, Link: domain.ListingLink{State: domain.LinkStateUnresolved}, SyncState: domain.ListingSyncStateSynced}}, AsOf: at}}
	w := httptest.NewRecorder()
	NewReadHandler(svc).HandleList(w, httptest.NewRequest(http.MethodGet, "/listings?installation_id=i", nil))
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	var body map[string]any
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["page_size"] != float64(1) || body["next_cursor"] != nil || body["as_of"] != "2026-07-15T12:00:00Z" {
		t.Fatalf("body=%s", w.Body.String())
	}
	item := body["items"].([]any)[0].(map[string]any)
	for _, key := range []string{"listing_type", "price", "published_quantity", "sync_error", "pending_issue", "cost", "below_margin_worst_case", "icms_worst_case_by_uf", "fetched_at"} {
		if _, ok := item[key]; !ok {
			t.Errorf("missing %s: %s", key, w.Body.String())
		}
	}
	// quality_score/sales_30d left the HTTP read contract in Task 8 (ADR-C3):
	// no producer wires them into the response anymore.
	for _, key := range []string{"quality_score", "sales_30d"} {
		if _, ok := item[key]; ok {
			t.Errorf("stale %s still on the wire: %s", key, w.Body.String())
		}
	}
}

func TestGetHandlerMalformedAndUnknownIDsReturnNested404(t *testing.T) {
	for _, raw := range []string{"", "i~x", "i~~-", "i~x~-~extra"} {
		t.Run(raw, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/listings/"+raw, nil)
			r.SetPathValue("listing_id", raw)
			w := httptest.NewRecorder()
			NewReadHandler(&fakeListService{getErr: &domain.ListingNotFoundError{}}).HandleGet(w, r)
			if w.Code != http.StatusNotFound {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			var body errorEnvelopeShape
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Error.Code != "listing_not_found" || body.Error.Details == nil {
				t.Fatalf("body=%s", w.Body.String())
			}
		})
	}
}

func TestGetHandlerSuccessIncludesTopLevelModelAndNonNilTimeline(t *testing.T) {
	svc := &fakeListService{
		model:    domain.ListingReadModel{ListingID: "i~x~-", InstallationID: "i", ProviderListingID: "x", Title: "X", ICMSWorstCaseByUF: &[]domain.ICMWorstCaseByUF{}},
		timeline: []domain.TimelineEvent{},
	}
	r := httptest.NewRequest(http.MethodGet, "/listings/i~x~-", nil)
	r.SetPathValue("listing_id", "i~x~-")
	w := httptest.NewRecorder()
	NewReadHandler(svc).HandleGet(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["listing_id"] != "i~x~-" || body["title"] != "X" {
		t.Fatalf("body=%s", w.Body.String())
	}
	if timeline, ok := body["timeline"].([]any); !ok || timeline == nil {
		t.Fatalf("timeline=%#v body=%s", body["timeline"], w.Body.String())
	}
}

func TestGetHandlerServiceNotFoundAndWrongMethod(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/listings/i~x~-", nil)
	r.SetPathValue("listing_id", "i~x~-")
	w := httptest.NewRecorder()
	NewReadHandler(&fakeListService{getErr: errors.New("listing_not_found")}).HandleGet(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("untyped error status=%d body=%s", w.Code, w.Body.String())
	}
	r = httptest.NewRequest(http.MethodPost, "/listings/i~x~-", nil)
	r.SetPathValue("listing_id", "i~x~-")
	w = httptest.NewRecorder()
	NewReadHandler(&fakeListService{}).HandleGet(w, r)
	if w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") != http.MethodGet {
		t.Fatalf("status=%d allow=%q", w.Code, w.Header().Get("Allow"))
	}
}

func TestByProductHandlerValidationAndIC02Envelope(t *testing.T) {
	h := NewReadHandler(&fakeListService{})
	for _, url := range []string{"/listings/by-product", "/listings/by-product?installation_id=i&filter.nope=x", "/listings/by-product?installation_id=i&cursor=bad"} {
		w := httptest.NewRecorder()
		h.HandleByProduct(w, httptest.NewRequest(http.MethodGet, url, nil))
		if w.Code != 400 {
			t.Fatalf("%s status=%d body=%s", url, w.Code, w.Body.String())
		}
	}
	title := "sem produto"
	svc := &fakeListService{groupPage: ports.ListingGroupRowPage{Groups: []domain.ListingGroup{{ProductTitle: &title, ListingCount: 1, GroupState: domain.GroupStateAttention, Listings: []domain.ListingReadModel{{ListingID: "i~x~-"}}}}, AsOf: time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)}}
	w := httptest.NewRecorder()
	NewReadHandler(svc).HandleByProduct(w, httptest.NewRequest(http.MethodGet, "/listings/by-product?installation_id=i", nil))
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	groups := body["groups"].([]any)
	group := groups[0].(map[string]any)
	if w.Code != 200 || body["page_size"] != float64(1) || group["product_id"] != nil || group["product_title"] != title || group["listing_count"] != float64(1) || group["group_state"] != "attention" {
		t.Fatalf("body=%s", w.Body.String())
	}
}

func TestSummaryHandlerEnvelopeAndErrors(t *testing.T) {
	one := 1
	at := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)
	svc := &fakeListService{summary: ports.ListingSummaryRow{Total: 3, Active: 1, Paused: 1, SyncError: 1, Stale: 1, Unlinked: 1, BelowMarginWorstCase: &one, MarginUnknown: nil, SemVinculo: 1, AbaixoCusto: 2, SemEvidencia: 3, AsOf: at}}
	w := httptest.NewRecorder()
	NewReadHandler(svc).HandleSummary(w, httptest.NewRequest(http.MethodGet, "/listings/summary?installation_id=i", nil))
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	exceptions := body["exceptions"].(map[string]any)
	if w.Code != 200 || body["total"] != float64(3) || body["as_of"] != "2026-07-15T12:00:00Z" || exceptions["below_margin_worst_case"] != float64(1) || exceptions["margin_unknown"] != nil {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	// F01-S4 additive counters: strictly additional keys alongside the 4
	// original ones (W1's response remains a valid subset of this envelope).
	if exceptions["sem_vinculo"] != float64(1) || exceptions["abaixo_custo"] != float64(2) || exceptions["sem_evidencia"] != float64(3) {
		t.Fatalf("additive exception counters missing/wrong: %+v", exceptions)
	}
	for _, tc := range []struct {
		name, url, method string
		err               error
		status            int
		code              string
	}{
		{"required", "/listings/summary", http.MethodGet, nil, 400, "installation_required"},
		{"not found", "/listings/summary?installation_id=i", http.MethodGet, application.ErrInstallationNotFound, 404, "installation_not_found"},
		{"source", "/listings/summary?installation_id=i", http.MethodGet, errors.New("down"), 503, "source_unavailable"},
		{"method", "/listings/summary?installation_id=i", http.MethodPost, nil, 405, "method_not_allowed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			NewReadHandler(&fakeListService{summaryErr: tc.err}).HandleSummary(w, httptest.NewRequest(tc.method, tc.url, nil))
			var got errorEnvelopeShape
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if w.Code != tc.status || got.Error.Code != tc.code {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if tc.status == 405 && w.Header().Get("Allow") != http.MethodGet {
				t.Fatalf("allow=%q", w.Header().Get("Allow"))
			}
		})
	}
}
