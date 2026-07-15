package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

type fakeListService struct {
	page     ports.ListingRowPage
	calls    int
	model    domain.ListingReadModel
	timeline []domain.TimelineEvent
	getErr   error
}

func (f *fakeListService) List(context.Context, ports.ListingQuery) (ports.ListingRowPage, error) {
	f.calls++
	return f.page, nil
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
	for _, key := range []string{"listing_type", "price", "published_quantity", "sync_error", "quality_score", "pending_issue", "sales_30d", "cost", "below_margin_worst_case", "icms_worst_case_by_uf", "fetched_at"} {
		if _, ok := item[key]; !ok {
			t.Errorf("missing %s: %s", key, w.Body.String())
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
			var body listErrorEnvelope
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
