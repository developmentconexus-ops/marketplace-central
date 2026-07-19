package mercadolivre

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

func catalogMatchAccountRef() domain.ProviderAccountRef {
	return domain.ProviderAccountRef{
		TenantID:          "tenant_default",
		InstallationID:    "inst-1",
		ProviderCode:      "mercado_livre",
		ProviderAccountID: "691607102",
	}
}

func newCatalogMatchAdapter(baseURL string, client *http.Client) *CapabilityAdapter {
	return NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    baseURL,
		SiteID:     "MLB",
		HTTPClient: client,
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-1", nil
		},
		Now: func() time.Time { return time.Date(2026, 7, 19, 12, 0, 0, 0, time.UTC) },
	})
}

func TestCapabilityAdapterReadCatalogMatchResolvesBuyBoxWinner(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/products/search":
			query := r.URL.Query()
			if query.Get("site_id") != "MLB" || query.Get("product_identifier") != "7891234567890" {
				t.Fatalf("search query = %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"results":[{"id":"MLB20270041","name":"Torneira Gourmet","status":"active","domain_id":"MLB-BATHROOM_FAUCETS_AND_MIXERS"}]}`))
		case "/products/MLB20270041":
			_, _ = w.Write([]byte(`{"name":"Torneira Gourmet","buy_box_winner":{"category_id":"MLB1276","price":199.9,"listing_type_id":"gold_special"}}`))
		case "/sites/MLB/domain_discovery/search":
			if r.URL.Query().Get("q") != "Torneira Cozinha" || r.URL.Query().Get("limit") != "3" {
				t.Fatalf("discovery query = %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`[{"domain_id":"MLB-BATHROOM_FAUCETS_AND_MIXERS","category_id":"MLB1276","category_name":"Torneiras"}]`))
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	adapter := newCatalogMatchAdapter(server.URL, server.Client())

	snapshot, err := adapter.ReadCatalogMatch(context.Background(), domain.CatalogMatchInput{
		AccountRef: catalogMatchAccountRef(),
		EAN:        "7891234567890",
		Query:      "Torneira Cozinha",
	})
	if err != nil {
		t.Fatalf("ReadCatalogMatch() error = %v", err)
	}

	if len(snapshot.CatalogHits) != 1 {
		t.Fatalf("catalog hits = %d", len(snapshot.CatalogHits))
	}
	hit := snapshot.CatalogHits[0]
	if hit.ProductID != "MLB20270041" || hit.DomainID != "MLB-BATHROOM_FAUCETS_AND_MIXERS" || hit.Status != "active" {
		t.Fatalf("hit = %+v", hit)
	}
	if snapshot.BuyBox == nil {
		t.Fatal("buy box = nil, want winner")
	}
	if snapshot.BuyBox.CategoryID != "MLB1276" || snapshot.BuyBox.ListingType != "gold_special" {
		t.Fatalf("buy box = %+v", snapshot.BuyBox)
	}
	if snapshot.BuyBox.Price == nil || *snapshot.BuyBox.Price != 199.9 {
		t.Fatalf("buy box price = %v", snapshot.BuyBox.Price)
	}
	if len(snapshot.DomainDiscovery) != 1 || snapshot.DomainDiscovery[0].CategoryID != "MLB1276" {
		t.Fatalf("domain discovery = %+v", snapshot.DomainDiscovery)
	}
}

func TestCapabilityAdapterReadCatalogMatchNilBuyBoxFallsBackToProductName(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/products/search":
			_, _ = w.Write([]byte(`{"results":[{"id":"MLB22490763","name":"Ducha Higienica","status":"active","domain_id":"MLB-SHOWER_HEADS"}]}`))
		case "/products/MLB22490763":
			// Observed live 2026-07-13: catalog product with no buy_box_winner.
			_, _ = w.Write([]byte(`{"name":"Ducha Higienica","buy_box_winner":null}`))
		case "/sites/MLB/domain_discovery/search":
			if r.URL.Query().Get("q") != "Ducha Higienica" {
				t.Fatalf("discovery should fall back to product name, got q = %s", r.URL.Query().Get("q"))
			}
			_, _ = w.Write([]byte(`[{"domain_id":"MLB-SHOWER_HEADS","category_id":"MLB5551","category_name":"Duchas"}]`))
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	adapter := newCatalogMatchAdapter(server.URL, server.Client())

	snapshot, err := adapter.ReadCatalogMatch(context.Background(), domain.CatalogMatchInput{
		AccountRef: catalogMatchAccountRef(),
		EAN:        "7899999999999",
	})
	if err != nil {
		t.Fatalf("ReadCatalogMatch() error = %v", err)
	}
	if snapshot.BuyBox != nil {
		t.Fatalf("buy box = %+v, want nil", snapshot.BuyBox)
	}
	if len(snapshot.DomainDiscovery) != 1 || snapshot.DomainDiscovery[0].CategoryID != "MLB5551" {
		t.Fatalf("domain discovery = %+v", snapshot.DomainDiscovery)
	}
}

func TestCapabilityAdapterReadCatalogMatchRequiresIdentifierOrQuery(t *testing.T) {
	t.Parallel()

	adapter := newCatalogMatchAdapter("http://example.invalid", http.DefaultClient)

	_, err := adapter.ReadCatalogMatch(context.Background(), domain.CatalogMatchInput{
		AccountRef: catalogMatchAccountRef(),
	})
	if domain.ErrorCodeOf(err) != domain.ErrCodeProviderValidation {
		t.Fatalf("error code = %v, want %v", domain.ErrorCodeOf(err), domain.ErrCodeProviderValidation)
	}
}
