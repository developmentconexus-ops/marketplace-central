package mercadolivre

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

func TestSearchCatalogByEANMapsCandidatesInProviderOrder(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 17, 12, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/products/search" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("product_identifier"); got != "789 123/456" {
			t.Fatalf("product_identifier = %q", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("authorization = %q", got)
		}
		_, _ = io.WriteString(w, `{"results":[{"id":" MLB-CATALOG-2 ","name":" Second product ","attributes":[{"id":"GTIN","name":" GTIN ","value_name":" 789123456 "},{"id":"BRAND","name":"Brand","value_name":null}]},{"id":"MLB-CATALOG-1","name":"First product","attributes":[{"id":"MODEL","name":"Model"}]}]}`)
	}))
	defer server.Close()

	result, err := catalogIdentityTestAdapter(server.URL, fetchedAt).SearchCatalogByEAN(context.Background(), pricingAccountRef(), "789 123/456")
	if err != nil {
		t.Fatalf("SearchCatalogByEAN() error = %v", err)
	}
	if !result.FetchedAt.Equal(fetchedAt) {
		t.Fatalf("FetchedAt = %v, want %v", result.FetchedAt, fetchedAt)
	}
	if len(result.Products) != 2 {
		t.Fatalf("products = %#v", result.Products)
	}
	if got := result.Products[0]; got.CatalogProductID != "MLB-CATALOG-2" || got.Title != "Second product" || len(got.Attrs) != 2 {
		t.Fatalf("products[0] = %#v", got)
	}
	if got := result.Products[0].Attrs[0]; got.ID != "GTIN" || got.Name != "GTIN" || got.ValueName == nil || *got.ValueName != "789123456" {
		t.Fatalf("products[0].attrs[0] = %#v", got)
	}
	if got := result.Products[0].Attrs[1]; got.ValueName != nil {
		t.Fatalf("products[0].attrs[1].ValueName = %q, want nil", *got.ValueName)
	}
	if got := result.Products[1]; got.CatalogProductID != "MLB-CATALOG-1" || got.Title != "First product" || len(got.Attrs) != 1 || got.Attrs[0].ValueName != nil {
		t.Fatalf("products[1] = %#v", got)
	}
}

// TestSearchCatalogByEANSendsSiteIDAndProductIdentifier is the regression guard
// for FINDING-M02-COLLECT-4XX (D-86): the ML /products/search endpoint 4XXs
// when site_id is absent, so every catalog EAN search MUST carry site_id
// alongside product_identifier — the same contract catalog_match_reader's
// searchCatalog already honors and that homologation proved live (D-83).
func TestSearchCatalogByEANSendsSiteIDAndProductIdentifier(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/products/search" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("site_id"); got != "MLB" {
			t.Fatalf("site_id = %q, want MLB", got)
		}
		if got := r.URL.Query().Get("product_identifier"); got != "7891234560019" {
			t.Fatalf("product_identifier = %q, want 7891234560019", got)
		}
		_, _ = io.WriteString(w, `{"results":[]}`)
	}))
	defer server.Close()

	if _, err := catalogIdentityTestAdapter(server.URL, time.Now().UTC()).SearchCatalogByEAN(context.Background(), pricingAccountRef(), "7891234560019"); err != nil {
		t.Fatalf("SearchCatalogByEAN() error = %v", err)
	}
}

func TestSearchCatalogByEANZeroResultsReturnsEmptyNonNil(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"results":[]}`)
	}))
	defer server.Close()

	result, err := catalogIdentityTestAdapter(server.URL, time.Now().UTC()).SearchCatalogByEAN(context.Background(), pricingAccountRef(), "789")
	if err != nil {
		t.Fatalf("SearchCatalogByEAN() error = %v", err)
	}
	if result.Products == nil || len(result.Products) != 0 {
		t.Fatalf("products = %#v, want empty non-nil slice", result.Products)
	}
}

func TestGetCatalogProductMapsBuyBoxAndAttributes(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Date(2026, 7, 17, 12, 1, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/products/MLB-CATALOG-1" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"id":" MLB-CATALOG-1 ","name":" Catalog product ","buy_box_winner":{"seller_id":" seller-7 ","price":123.45,"currency_id":"BRL","condition":"new"},"attributes":[{"id":"BRAND","name":" Brand ","value_name":" Maker "}]}`)
	}))
	defer server.Close()

	product, err := catalogIdentityTestAdapter(server.URL, fetchedAt).GetCatalogProduct(context.Background(), pricingAccountRef(), "MLB-CATALOG-1")
	if err != nil {
		t.Fatalf("GetCatalogProduct() error = %v", err)
	}
	if product.ID != "MLB-CATALOG-1" || product.Title != "Catalog product" || !product.FetchedAt.Equal(fetchedAt) {
		t.Fatalf("product = %#v", product)
	}
	if product.BuyBoxWinner == nil || product.BuyBoxWinner.SellerID != "seller-7" || product.BuyBoxWinner.Price == nil || product.BuyBoxWinner.Price.Amount != "123.45" || product.BuyBoxWinner.Price.Currency != "BRL" {
		t.Fatalf("buy_box_winner = %#v", product.BuyBoxWinner)
	}
	if len(product.Attrs) != 1 || product.Attrs[0].ID != "BRAND" || product.Attrs[0].Name != "Brand" || product.Attrs[0].ValueName == nil || *product.Attrs[0].ValueName != "Maker" {
		t.Fatalf("attrs = %#v", product.Attrs)
	}
}

func TestGetCatalogProductNullBuyBoxWinnerAndAbsentValueNameStayNil(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"id":"MLB-CATALOG-NULL","name":"No winner","buy_box_winner":null,"attributes":[{"id":"MODEL","name":"Model","value_name":null},{"id":"GTIN","name":"GTIN"}]}`)
	}))
	defer server.Close()

	product, err := catalogIdentityTestAdapter(server.URL, time.Now().UTC()).GetCatalogProduct(context.Background(), pricingAccountRef(), "MLB-CATALOG-NULL")
	if err != nil {
		t.Fatalf("GetCatalogProduct() error = %v", err)
	}
	if product.BuyBoxWinner != nil {
		t.Fatalf("BuyBoxWinner = %#v, want nil", product.BuyBoxWinner)
	}
	if len(product.Attrs) != 2 || product.Attrs[0].ValueName != nil || product.Attrs[1].ValueName != nil {
		t.Fatalf("attrs = %#v", product.Attrs)
	}
}

func TestCatalogIdentityReaderMapsProviderErrors(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		status int
		want   error
	}{{http.StatusUnauthorized, domain.ErrUnauthorized}, {http.StatusNotFound, domain.ErrNotFound}, {http.StatusTooManyRequests, domain.ErrRateLimited}, {http.StatusInternalServerError, domain.ErrProviderUnavailable}} {
		t.Run(fmt.Sprint(test.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(test.status) }))
			defer server.Close()

			// Both ports route provider failures through the same mapPricingReaderError; assert each in isolation.
			if _, err := catalogIdentityTestAdapter(server.URL, time.Now().UTC()).SearchCatalogByEAN(context.Background(), pricingAccountRef(), "789"); !errors.Is(err, test.want) {
				t.Fatalf("SearchCatalogByEAN error = %v, want %v", err, test.want)
			}
			if _, err := catalogIdentityTestAdapter(server.URL, time.Now().UTC()).GetCatalogProduct(context.Background(), pricingAccountRef(), "MLB-CATALOG-1"); !errors.Is(err, test.want) {
				t.Fatalf("GetCatalogProduct error = %v, want %v", err, test.want)
			}
		})
	}
}

func catalogIdentityTestAdapter(baseURL string, fetchedAt time.Time) *CapabilityAdapter {
	return NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: time.Second},
		Now:        func() time.Time { return fetchedAt },
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "test-token", nil
		},
	})
}
