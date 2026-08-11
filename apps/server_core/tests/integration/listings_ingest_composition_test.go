//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre"
	"marketplace-central/apps/server_core/internal/composition"
	"marketplace-central/apps/server_core/internal/contexts/listings"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// listingsCompositionAccountID is the fake seller id the fixture below serves
// under /users/{id}/items/search — it has to agree with the Config.UserID
// passed to mercadolivre.New, exactly like feed_test.go's fixture does.
const listingsCompositionAccountID = "179571326"

// listingsCompositionFixtureServer is DUPLICATED, not extracted, from
// apps/server_core/internal/adapters/marketplace/mercadolivre/listings/feed_test.go's
// fixtureServer (plan Task 8 step 1 explicitly allows either choice for a
// _test.go helper shared across distinct packages). It serves a deterministic
// 2-page scan (page1: MLB1+MLB2, page2: MLB3) followed by an empty/Done page,
// and deterministic multiget bodies — no timestamps or counters — so a
// second walk over the same feed produces byte-identical raw payloads and
// therefore Idempotent dispositions, not Changed ones.
func listingsCompositionFixtureServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/users/"+listingsCompositionAccountID+"/items/search", func(w http.ResponseWriter, r *http.Request) {
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
		for _, id := range listingsCompositionSplitIDs(ids) {
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

func listingsCompositionSplitIDs(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// TestListingsIngestCompositionWritesRows is the production-path acceptance
// for the listings slice: a real mercadolivre.Bundle feed, driven by
// composition.RunListingsIngest, against the real Postgres lane, must land
// real rows — not just a report the code makes about itself. A second walk
// over the same (stateless) fixture must be free.
func TestListingsIngestCompositionWritesRows(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_listings_composition")
	tid, err := tenant.Parse(cfg.DefaultTenantID)
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}

	srv := listingsCompositionFixtureServer(t)
	defer srv.Close()

	bundle, err := mercadolivre.New(mercadolivre.Config{
		BaseURL:   srv.URL,
		UserID:    listingsCompositionAccountID,
		Channel:   "mercado_livre",
		AccountID: listingsCompositionAccountID,
		Token:     func(context.Context) (string, error) { return "tok-test", nil },
	})
	if err != nil {
		t.Fatalf("mercadolivre.New: %v", err)
	}

	module := listings.New(pool)
	ctx := context.Background()

	report, err := composition.RunListingsIngest(ctx, module, bundle.ListingFeed, tid, 50)
	if err != nil {
		t.Fatalf("RunListingsIngest (first walk): %v", err)
	}
	if report.Pages != 3 || report.Observed != 3 || report.Created != 3 {
		t.Fatalf("first walk report = %+v, want Pages=3 Observed=3 Created=3", report)
	}

	second, err := composition.RunListingsIngest(ctx, module, bundle.ListingFeed, tid, 50)
	if err != nil {
		t.Fatalf("RunListingsIngest (second walk): %v", err)
	}
	if second.Idempotent != 3 {
		t.Fatalf("second walk report = %+v, want Idempotent=3", second)
	}

	// A aceitação é o BANCO, não o retorno (aceite = observável). Todas as três
	// tabelas: uma listagem gravada sem a sua variação, ou sem a observação de
	// origem que a justifica, é uma escrita parcial que o count da tabela-mãe
	// sozinho não vê. source_observations continua em 3 depois do segundo walk
	// porque o hash é o mesmo — idempotente não grava versão nova.
	for _, want := range []struct {
		table string
		query string
		count int
	}{
		{"listings.listings", `SELECT count(*) FROM listings.listings WHERE tenant_id=$1`, 3},
		{"listings.listing_variations", `SELECT count(*) FROM listings.listing_variations WHERE tenant_id=$1`, 3},
		{"listings.source_observations", `SELECT count(*) FROM listings.source_observations WHERE tenant_id=$1`, 3},
	} {
		var got int
		if err := scanOne(pool, tid, &got, want.query); err != nil {
			t.Fatalf("count %s: %v", want.table, err)
		}
		if got != want.count {
			t.Fatalf("%s count = %d, want %d", want.table, got, want.count)
		}
	}
}
