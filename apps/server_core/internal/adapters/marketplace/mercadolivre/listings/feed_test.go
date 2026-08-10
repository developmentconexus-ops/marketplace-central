package listings_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre"
	"marketplace-central/apps/server_core/internal/contexts/listings/port"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// fixture: scan devolve 2 ids na página 1 (scroll_id avança) e 1 id na página
// 2; multiget devolve bodies com title/price/status/sku/GTIN e uma variação.
func fixtureServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/users/179571326/items/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok-test" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.URL.Query().Get("search_type") != "scan" {
			t.Errorf("search_type = %q, want scan", r.URL.Query().Get("search_type"))
		}
		switch r.URL.Query().Get("scroll_id") {
		case "":
			_ = json.NewEncoder(w).Encode(map[string]any{"scroll_id": "s2", "results": []string{"MLB1", "MLB2"}})
		case "s2":
			_ = json.NewEncoder(w).Encode(map[string]any{"scroll_id": "s3", "results": []string{"MLB3"}})
		default:
			_ = json.NewEncoder(w).Encode(map[string]any{"scroll_id": "", "results": []string{}})
		}
	})
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query().Get("ids")
		var out []map[string]any
		for _, id := range splitIDs(ids) {
			body := map[string]any{
				"id": id, "title": "Produto " + id, "status": "active",
				"price": json.Number("199.90"), "currency_id": "BRL",
				"listing_type_id": "gold_special", "available_quantity": 5,
				"seller_sku": "SKU-" + id,
				"attributes": []map[string]any{{"id": "GTIN", "value_name": "7891234567890"}},
				"variations": []map[string]any{{
					"id": 111, "price": json.Number("199.90"), "available_quantity": 3,
					"seller_sku": "VSKU-" + id,
					"attributes": []map[string]any{{"id": "GTIN", "value_name": "7899999999999"}},
				}},
			}
			out = append(out, map[string]any{"code": 200, "body": body})
		}
		_ = json.NewEncoder(w).Encode(out)
	})
	return httptest.NewServer(mux)
}

func splitIDs(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func newFeed(t *testing.T, baseURL string) port.ListingFeed {
	t.Helper()
	bundle, err := mercadolivre.New(mercadolivre.Config{
		BaseURL:   baseURL,
		UserID:    "179571326",
		Channel:   "mercado_livre",
		AccountID: "179571326",
		Token:     func(context.Context) (string, error) { return "tok-test", nil },
	})
	if err != nil {
		t.Fatalf("bundle: %v", err)
	}
	return bundle.ListingFeed
}

func TestFeedWalksEveryScanPageAndMapsFacts(t *testing.T) {
	srv := fixtureServer(t)
	defer srv.Close()
	feed := newFeed(t, srv.URL)
	tid, _ := tenant.Parse("tenant_default")

	page1, err := feed.NextPage(context.Background(), tid, port.Cursor{}, 50)
	if err != nil {
		t.Fatalf("page 1: %v", err)
	}
	if len(page1.Observations) != 2 || page1.Done {
		t.Fatalf("page 1: %d obs done=%v, want 2 obs not done", len(page1.Observations), page1.Done)
	}
	obs := page1.Observations[0]
	if title, ok := obs.Title.Value(); !ok || title != "Produto MLB1" {
		t.Fatalf("title = %v known=%v", title, ok)
	}
	if price, ok := obs.Price.Value(); !ok || price.Amount().StringFixed(2) != "199.90" || price.Currency().String() != "BRL" {
		t.Fatalf("price fact wrong: %v known=%v", price, ok)
	}
	if gtin, ok := obs.GTIN.Value(); !ok || gtin != "7891234567890" {
		t.Fatalf("gtin = %q known=%v", gtin, ok)
	}
	if len(obs.Variations) != 1 || obs.Variations[0].VariationID != "111" {
		t.Fatalf("variations = %+v", obs.Variations)
	}
	if len(obs.RawPayload) == 0 {
		t.Fatal("raw payload not retained")
	}

	page2, err := feed.NextPage(context.Background(), tid, page1.Next, 50)
	if err != nil {
		t.Fatalf("page 2: %v", err)
	}
	if len(page2.Observations) != 1 {
		t.Fatalf("page 2: %d obs, want 1 (page-1 truncation blindness)", len(page2.Observations))
	}

	page3, err := feed.NextPage(context.Background(), tid, page2.Next, 50)
	if err != nil {
		t.Fatalf("page 3: %v", err)
	}
	if !page3.Done {
		t.Fatal("empty scan page must report Done")
	}
}

func TestFeedFailsLoudOnPerItemError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/179571326/items/search", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"scroll_id": "s2", "results": []string{"MLB9"}})
	})
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]any{{"code": 404, "body": map[string]any{}}})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	feed := newFeed(t, srv.URL)
	tid, _ := tenant.Parse("tenant_default")
	_, err := feed.NextPage(context.Background(), tid, port.Cursor{}, 50)
	if err == nil || !strings.Contains(err.Error(), "MLB9") {
		t.Fatalf("per-item failure must fail the page naming the id, got: %v", err)
	}
}

func TestFeedRejectsNonemptyResultsWithEmptyScrollID(t *testing.T) {
	// ML returning results with an empty scroll_id is malformed: port.NewCursor("")
	// is the FEED START cursor, and handing that back for a nonempty page would let
	// a future caller loop forever re-reading page 1. The adapter is the only code
	// that knows what an empty scroll_id from ML means, so it must fail loud here
	// rather than let the port contract carry an ambiguous cursor.
	mux := http.NewServeMux()
	mux.HandleFunc("/users/179571326/items/search", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"scroll_id": "", "results": []string{"MLB1"}})
	})
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]any{{"code": 200, "body": map[string]any{"id": "MLB1"}}})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	feed := newFeed(t, srv.URL)
	tid, _ := tenant.Parse("tenant_default")
	_, err := feed.NextPage(context.Background(), tid, port.Cursor{}, 50)
	if err == nil || !strings.Contains(err.Error(), "scroll_id") {
		t.Fatalf("nonempty results with empty scroll_id must fail naming scroll_id, got: %v", err)
	}
}

func TestFeedMapsAbsentFieldsToUnknownNeverZero(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/179571326/items/search", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"scroll_id": "s2", "results": []string{"MLB7"}})
	})
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		// body só com id: todo o resto ausente.
		_ = json.NewEncoder(w).Encode([]map[string]any{{"code": 200, "body": map[string]any{"id": "MLB7"}}})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	feed := newFeed(t, srv.URL)
	tid, _ := tenant.Parse("tenant_default")
	page, err := feed.NextPage(context.Background(), tid, port.Cursor{}, 50)
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	o := page.Observations[0]
	if _, ok := o.Price.Value(); ok {
		t.Fatal("absent price came back Known — the zero-fabrication this kernel exists to end")
	}
	if o.Price.State().String() != "unknown" {
		t.Fatalf("price state = %s, want unknown", o.Price.State())
	}
}
