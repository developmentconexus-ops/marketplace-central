package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

type fakeListService struct {
	page  ports.ListingRowPage
	calls int
}

func (f *fakeListService) List(context.Context, ports.ListingQuery) (ports.ListingRowPage, error) {
	f.calls++
	return f.page, nil
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
