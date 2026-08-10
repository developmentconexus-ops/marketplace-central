package listings_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre"
	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/port"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// fixture: scan devolve 2 ids na página 1 (scroll_id avança) e 1 id na página
// 2; multiget devolve bodies com title/price/status/sku/GTIN e uma variação.
//
// O SKU vem onde o Mercado Livre realmente o põe: no atributo SELLER_SKU e no
// seller_custom_field. NÃO existe campo seller_sku de topo — nem na resposta
// completa documentada (developers.mercadolivre.com.br/pt_br/variacoes), nem
// em nenhum dos 34 payloads reais gravados no primeiro live drive. A fixture
// anterior inventava esse campo, e por isso provava o mapper contra um canal
// que não existe.
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
				"seller_custom_field": "SCF-" + id,
				"attributes": []map[string]any{
					{"id": "GTIN", "value_name": "7891234567890"},
					{"id": "SELLER_SKU", "value_name": "SKU-" + id},
				},
				"variations": []map[string]any{{
					"id": 111, "price": json.Number("199.90"), "available_quantity": 3,
					"seller_custom_field": "VSCF-" + id,
					"attributes": []map[string]any{
						{"id": "GTIN", "value_name": "7899999999999"},
						{"id": "SELLER_SKU", "value_name": "VSKU-" + id},
					},
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
	if title, ok := obs.State.Title.Value(); !ok || title != "Produto MLB1" {
		t.Fatalf("title = %v known=%v", title, ok)
	}
	if price, ok := obs.State.Price.Value(); !ok || price.Amount().StringFixed(2) != "199.90" || price.Currency().String() != "BRL" {
		t.Fatalf("price fact wrong: %v known=%v", price, ok)
	}
	if gtin, ok := obs.State.GTIN.Value(); !ok || gtin != "7891234567890" {
		t.Fatalf("gtin = %q known=%v", gtin, ok)
	}
	if sku, ok := obs.State.SellerSKU.Value(); !ok || sku != "SKU-MLB1" {
		t.Fatalf("seller_sku = %q known=%v state=%s reason=%q; ML sent it in the SELLER_SKU attribute and in seller_custom_field",
			sku, ok, obs.State.SellerSKU.State(), obs.State.SellerSKU.Reason())
	}
	if len(obs.State.Variations) != 1 || obs.State.Variations[0].VariationID != "111" {
		t.Fatalf("variations = %+v", obs.State.Variations)
	}
	if sku, ok := obs.State.Variations[0].SellerSKU.Value(); !ok || sku != "VSKU-MLB1" {
		t.Fatalf("variation seller_sku = %q known=%v state=%s reason=%q",
			sku, ok, obs.State.Variations[0].SellerSKU.State(), obs.State.Variations[0].SellerSKU.Reason())
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

// observeOne serves a single multiget body and returns the mapped observation,
// so a case only has to state the SKU-bearing fields it wants to exercise.
func observeOne(t *testing.T, body map[string]any) contracts.ListingObservation {
	t.Helper()
	body["id"] = "MLB1"
	mux := http.NewServeMux()
	mux.HandleFunc("/users/179571326/items/search", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"scroll_id": "s2", "results": []string{"MLB1"}})
	})
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]any{{"code": 200, "body": body}})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	tid, _ := tenant.Parse("tenant_default")
	page, err := newFeed(t, srv.URL).NextPage(context.Background(), tid, port.Cursor{}, 50)
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	return page.Observations[0]
}

// TestFeedResolvesSellerSKUByMercadoLivresDocumentedPriority pins the order ML
// itself publishes (developers.mercadolivre.com.br/pt_br/mercado-envios-2):
// SELLER_SKU attribute of the variation, then the variation's
// seller_custom_field, then the item's SELLER_SKU attribute, then the item's
// seller_custom_field. Attribute before custom field is not a preference — ML
// documents SELLER_SKU as the place a SKU belongs and seller_custom_field as
// the legacy field it superseded, so reading the custom field first would let a
// stale value shadow the one the seller maintains.
func TestFeedResolvesSellerSKUByMercadoLivresDocumentedPriority(t *testing.T) {
	attr := func(id, value string) map[string]any {
		return map[string]any{"id": id, "value_name": value}
	}

	t.Run("item attribute beats item custom field", func(t *testing.T) {
		obs := observeOne(t, map[string]any{
			"seller_custom_field": "legacy-value",
			"attributes":          []map[string]any{attr("SELLER_SKU", "attribute-value")},
		})
		if sku, ok := obs.State.SellerSKU.Value(); !ok || sku != "attribute-value" {
			t.Fatalf("seller_sku = %q known=%v, want the SELLER_SKU attribute", sku, ok)
		}
	})

	t.Run("item custom field is used when no attribute", func(t *testing.T) {
		obs := observeOne(t, map[string]any{"seller_custom_field": "legacy-value"})
		if sku, ok := obs.State.SellerSKU.Value(); !ok || sku != "legacy-value" {
			t.Fatalf("seller_sku = %q known=%v, want seller_custom_field", sku, ok)
		}
	})

	t.Run("neither source is unknown with a reason", func(t *testing.T) {
		obs := observeOne(t, map[string]any{})
		if _, ok := obs.State.SellerSKU.Value(); ok {
			t.Fatal("a listing with no SKU anywhere came back Known")
		}
		if obs.State.SellerSKU.State().String() != "unknown" || obs.State.SellerSKU.Reason() == "" {
			t.Fatalf("state=%s reason=%q, want unknown with a reason", obs.State.SellerSKU.State(), obs.State.SellerSKU.Reason())
		}
	})

	t.Run("variation attribute beats every other source", func(t *testing.T) {
		obs := observeOne(t, map[string]any{
			"seller_custom_field": "item-legacy",
			"attributes":          []map[string]any{attr("SELLER_SKU", "item-attribute")},
			"variations": []map[string]any{{
				"id": 111, "seller_custom_field": "variation-legacy",
				"attributes": []map[string]any{attr("SELLER_SKU", "variation-attribute")},
			}},
		})
		if sku, ok := obs.State.Variations[0].SellerSKU.Value(); !ok || sku != "variation-attribute" {
			t.Fatalf("variation seller_sku = %q known=%v", sku, ok)
		}
	})

	t.Run("variation custom field beats the item", func(t *testing.T) {
		obs := observeOne(t, map[string]any{
			"attributes": []map[string]any{attr("SELLER_SKU", "item-attribute")},
			"variations": []map[string]any{{"id": 111, "seller_custom_field": "variation-legacy"}},
		})
		if sku, ok := obs.State.Variations[0].SellerSKU.Value(); !ok || sku != "variation-legacy" {
			t.Fatalf("variation seller_sku = %q known=%v", sku, ok)
		}
	})

	t.Run("variation with no SKU of its own falls back to the item", func(t *testing.T) {
		obs := observeOne(t, map[string]any{
			"seller_custom_field": "item-legacy",
			"variations":          []map[string]any{{"id": 111}},
		})
		if sku, ok := obs.State.Variations[0].SellerSKU.Value(); !ok || sku != "item-legacy" {
			t.Fatalf("variation seller_sku = %q known=%v, want the item's value", sku, ok)
		}
	})
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
	if _, ok := o.State.Price.Value(); ok {
		t.Fatal("absent price came back Known — the zero-fabrication this kernel exists to end")
	}
	if o.State.Price.State().String() != "unknown" {
		t.Fatalf("price state = %s, want unknown", o.State.Price.State())
	}
}
