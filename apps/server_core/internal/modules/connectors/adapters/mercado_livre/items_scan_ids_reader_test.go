package mercadolivre

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// TestListListingsScanIDsReturnsIDsWithoutHydrating proves Passo 1's central
// contract (IC-06 "enumerador nunca hidrata"): the ids-only scan must return
// exactly the ids from `results` plus the next scroll_id, and it must NEVER
// call the per-item hydration endpoint (/items/{id}) the way
// ListListingsScanPage's N+1 loop does — this test's server 404s any request
// to that path so a regression would be caught immediately, not silently.
func TestListListingsScanIDsReturnsIDsWithoutHydrating(t *testing.T) {
	t.Parallel()

	hydrationCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/users/seller-1/items/search" {
			if got := r.URL.Query().Get("search_type"); got != "scan" {
				t.Fatalf("search_type = %q, want scan", got)
			}
			if got := r.URL.Query().Get("limit"); got != "50" {
				t.Fatalf("limit = %q, want 50 (default)", got)
			}
			_, _ = io.WriteString(w, `{"scroll_id":"scroll-abc","results":["MLB1","MLB2","MLB3"]}`)
			return
		}
		hydrationCalls++
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	ids, nextCursor, err := pricingTestAdapter(server.URL, fixedNow()).ListListingsScanIDs(context.Background(), connectorsdomain.ListListingsInput{AccountRef: pricingAccountRef()})
	if err != nil {
		t.Fatalf("ListListingsScanIDs() error = %v", err)
	}
	if hydrationCalls != 0 {
		t.Fatalf("hydration calls = %d, want 0 (enumerator must never hydrate)", hydrationCalls)
	}
	wantIDs := []string{"MLB1", "MLB2", "MLB3"}
	if len(ids) != len(wantIDs) {
		t.Fatalf("ids = %#v, want %#v", ids, wantIDs)
	}
	for i, id := range ids {
		if id != wantIDs[i] {
			t.Fatalf("ids[%d] = %q, want %q", i, id, wantIDs[i])
		}
	}
	if nextCursor != "scroll-abc" {
		t.Fatalf("nextCursor = %q, want scroll-abc", nextCursor)
	}
}

// TestListListingsScanIDsChainsAcrossMultiplePages is the required test (2),
// R-3: a fixture with MORE THAN ONE enumeration page. The server hands out
// scroll_id "page-2" on the first call (no incoming scroll_id) and an empty
// scroll_id (scan exhausted) on the second call (incoming scroll_id =
// "page-2") — the caller must chain cursor -> next scroll_id across calls and
// must never resend a scroll_id it wasn't just given.
func TestListListingsScanIDsChainsAcrossMultiplePages(t *testing.T) {
	t.Parallel()

	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		scrollID := r.URL.Query().Get("scroll_id")
		switch calls {
		case 1:
			if scrollID != "" {
				t.Fatalf("call 1 scroll_id = %q, want empty (first page)", scrollID)
			}
			_, _ = io.WriteString(w, `{"scroll_id":"page-2","results":["MLB-P1-A","MLB-P1-B"]}`)
		case 2:
			if scrollID != "page-2" {
				t.Fatalf("call 2 scroll_id = %q, want page-2", scrollID)
			}
			_, _ = io.WriteString(w, `{"scroll_id":"","results":["MLB-P2-A"]}`)
		default:
			t.Fatalf("unexpected call %d", calls)
		}
	}))
	defer server.Close()

	adapter := pricingTestAdapter(server.URL, fixedNow())
	account := pricingAccountRef()

	page1IDs, cursor1, err := adapter.ListListingsScanIDs(context.Background(), connectorsdomain.ListListingsInput{AccountRef: account})
	if err != nil {
		t.Fatalf("page 1: %v", err)
	}
	if len(page1IDs) != 2 || cursor1 != "page-2" {
		t.Fatalf("page 1 = %#v cursor=%q, want [MLB-P1-A MLB-P1-B] page-2", page1IDs, cursor1)
	}

	page2IDs, cursor2, err := adapter.ListListingsScanIDs(context.Background(), connectorsdomain.ListListingsInput{AccountRef: account, Cursor: cursor1})
	if err != nil {
		t.Fatalf("page 2: %v", err)
	}
	if len(page2IDs) != 1 || page2IDs[0] != "MLB-P2-A" {
		t.Fatalf("page 2 = %#v, want [MLB-P2-A]", page2IDs)
	}
	if cursor2 != "" {
		t.Fatalf("cursor2 = %q, want empty (scan exhausted)", cursor2)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
}

// TestListListingsScanIDsSkipsBlankResultEntries defends against a provider
// response containing whitespace-only ids (defensive, mirrors the trimming
// ListListingsScanPage's caller relies on for scroll_id).
func TestListListingsScanIDsSkipsBlankResultEntries(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"scroll_id":"  ","results":["MLB1","  ","MLB2",""]}`)
	}))
	defer server.Close()

	ids, nextCursor, err := pricingTestAdapter(server.URL, fixedNow()).ListListingsScanIDs(context.Background(), connectorsdomain.ListListingsInput{AccountRef: pricingAccountRef()})
	if err != nil {
		t.Fatalf("ListListingsScanIDs() error = %v", err)
	}
	if len(ids) != 2 || ids[0] != "MLB1" || ids[1] != "MLB2" {
		t.Fatalf("ids = %#v, want [MLB1 MLB2]", ids)
	}
	if nextCursor != "" {
		t.Fatalf("nextCursor = %q, want empty (trimmed)", nextCursor)
	}
}

// TestListListingsScanIDsRespectsExplicitLimit ensures the batch-size
// (limit=1000 for enumeration per IC-06) is forwarded verbatim, not silently
// defaulted, whenever the caller supplies one.
func TestListListingsScanIDsRespectsExplicitLimit(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("limit"); got != "1000" {
			t.Fatalf("limit = %q, want 1000", got)
		}
		ids := make([]string, 0, 3)
		for i := 0; i < 3; i++ {
			ids = append(ids, "MLB"+strconv.Itoa(i))
		}
		fmt.Fprintf(w, `{"scroll_id":"","results":[%q,%q,%q]}`, ids[0], ids[1], ids[2])
	}))
	defer server.Close()

	_, _, err := pricingTestAdapter(server.URL, fixedNow()).ListListingsScanIDs(context.Background(), connectorsdomain.ListListingsInput{AccountRef: pricingAccountRef(), Limit: 1000})
	if err != nil {
		t.Fatalf("ListListingsScanIDs() error = %v", err)
	}
}
