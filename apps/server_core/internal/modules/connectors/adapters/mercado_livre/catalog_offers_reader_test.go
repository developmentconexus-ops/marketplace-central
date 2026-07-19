package mercadolivre

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

func TestListCatalogOffersPaginatesAndPreservesProviderOrder(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("authorization = %q", got)
		}
		if r.URL.Path == "/products/MLB-CATALOG" {
			_, _ = io.WriteString(w, `{"children_ids":[]}`)
			return
		}
		calls.Add(1)
		if r.Method != http.MethodGet || r.URL.Path != "/products/MLB-CATALOG/items" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		switch r.URL.Query().Get("offset") {
		case "0":
			_, _ = io.WriteString(w, `{"paging":{"total":3,"offset":0,"limit":2},"results":[{"seller_id":" seller-2 ","price":120.00,"currency_id":"BRL","condition":" new ","shipping":{"mode":" me2 "}},{"seller_id":"seller-1","price":99.90,"currency_id":"BRL","condition":"used","shipping":{"mode":"custom"}}]}`)
		case "2":
			_, _ = io.WriteString(w, `{"paging":{"total":3,"offset":2,"limit":2},"results":[{"seller_id":"seller-2","price":110.50,"currency_id":"BRL","condition":"new","shipping":{"mode":"me1"}}]}`)
		default:
			t.Fatalf("offset = %q", r.URL.Query().Get("offset"))
		}
	}))
	defer server.Close()

	offers, err := catalogOffersTestAdapter(server.URL, true).ListCatalogOffers(context.Background(), pricingAccountRef(), "MLB-CATALOG")
	if err != nil {
		t.Fatalf("ListCatalogOffers() error = %v", err)
	}
	if calls.Load() != 2 || len(offers) != 3 {
		t.Fatalf("calls = %d, offers = %#v", calls.Load(), offers)
	}
	wants := []struct{ seller, amount, condition, shipping string }{{"seller-2", "120.00", "new", "me2"}, {"seller-1", "99.90", "used", "custom"}, {"seller-2", "110.50", "new", "me1"}}
	for i, want := range wants {
		got := offers[i]
		if got.SellerID != want.seller || got.Price == nil || got.Price.Amount != want.amount || got.Price.Currency != "BRL" || got.Condition != want.condition || got.ShippingMode != want.shipping {
			t.Fatalf("offers[%d] = %#v, want %#v", i, got, want)
		}
	}
}

func TestListCatalogOffersFlagOffAvoidsTokenAndHTTP(t *testing.T) {
	t.Parallel()

	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { requests.Add(1) }))
	defer server.Close()
	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL: server.URL,
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			t.Fatal("token resolver called while flag disabled")
			return "", nil
		},
	})

	offers, err := adapter.ListCatalogOffers(context.Background(), pricingAccountRef(), "MLB-CATALOG")
	if !errors.Is(err, domain.ErrCatalogOffersUnavailable) || offers != nil || requests.Load() != 0 {
		t.Fatalf("offers = %#v, error = %v, requests = %d", offers, err, requests.Load())
	}
}

func TestListCatalogOffersMidPaginationFailureReturnsNoPartialResults(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/products/MLB-CATALOG" {
			_, _ = io.WriteString(w, `{"children_ids":[]}`)
			return
		}
		if r.URL.Query().Get("offset") == "0" {
			_, _ = io.WriteString(w, `{"paging":{"total":2,"offset":0,"limit":1},"results":[{"seller_id":"seller-1","price":90,"currency_id":"BRL"}]}`)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// A mid-pagination 5xx is a provider fault, NOT the capability being disabled.
	// Disabled is reserved for the flag-off path (FINDING-M02-LIVE-2 D1, D-86).
	offers, err := catalogOffersTestAdapter(server.URL, true).ListCatalogOffers(context.Background(), pricingAccountRef(), "MLB-CATALOG")
	if !errors.Is(err, domain.ErrProviderUnavailable) || errors.Is(err, domain.ErrCatalogOffersUnavailable) || offers != nil {
		t.Fatalf("offers = %#v, error = %v", offers, err)
	}
}

func TestListCatalogOffersEmptyPageBeforeTotalReturnsNoPartialResults(t *testing.T) {
	t.Parallel()

	// Provider lies: Total=2 but the second page is genuinely empty (not an HTTP
	// error). The non-advancing guard must hard-stop with a typed provider error
	// (integrity fault, NOT capability-disabled), never return the one
	// already-accumulated offer nor loop forever (FINDING-M02-LIVE-2 D1).
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/products/MLB-CATALOG" {
			_, _ = io.WriteString(w, `{"children_ids":[]}`)
			return
		}
		if r.URL.Query().Get("offset") == "0" {
			_, _ = io.WriteString(w, `{"paging":{"total":2,"offset":0,"limit":1},"results":[{"seller_id":"seller-1","price":90,"currency_id":"BRL"}]}`)
			return
		}
		_, _ = io.WriteString(w, `{"paging":{"total":2,"offset":1,"limit":1},"results":[]}`)
	}))
	defer server.Close()

	offers, err := catalogOffersTestAdapter(server.URL, true).ListCatalogOffers(context.Background(), pricingAccountRef(), "MLB-CATALOG")
	if !errors.Is(err, domain.ErrProviderUnavailable) || errors.Is(err, domain.ErrCatalogOffersUnavailable) || offers != nil {
		t.Fatalf("offers = %#v, error = %v", offers, err)
	}
}

func TestListCatalogOffersPreservesAbsentPrice(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"paging":{"total":1,"offset":0,"limit":50},"results":[{"seller_id":"seller-1","condition":"new","shipping":{"mode":"me2"}}]}`)
	}))
	defer server.Close()

	offers, err := catalogOffersTestAdapter(server.URL, true).ListCatalogOffers(context.Background(), pricingAccountRef(), "MLB-CATALOG")
	if err != nil || len(offers) != 1 || offers[0].Price != nil {
		t.Fatalf("offers = %#v, error = %v", offers, err)
	}
}

func TestListCatalogOffersMapsProviderErrors(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		status int
		want   error
	}{{http.StatusUnauthorized, domain.ErrUnauthorized}, {http.StatusNotFound, domain.ErrNotFound}, {http.StatusTooManyRequests, domain.ErrRateLimited}} {
		t.Run(fmt.Sprint(test.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(test.status) }))
			defer server.Close()
			offers, err := catalogOffersTestAdapter(server.URL, true).ListCatalogOffers(context.Background(), pricingAccountRef(), "MLB-CATALOG")
			if !errors.Is(err, test.want) || offers != nil {
				t.Fatalf("offers = %#v, error = %v, want %v", offers, err, test.want)
			}
		})
	}
}

func TestListCatalogOffersFansOutParentToChildLeaves(t *testing.T) {
	t.Parallel()

	// MLB22624877 (codprod 90008) is a PARENT catalog product: /items exists only
	// on its LEAF children. The reader must read /products/{id}, and when
	// children_ids is non-empty fan out to each child's /items and aggregate
	// (FINDING-M02-LIVE-2 D2, D-86; official docs competencia-en-catalogo).
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("authorization = %q", got)
		}
		switch r.URL.Path {
		case "/products/MLB-PARENT":
			_, _ = io.WriteString(w, `{"children_ids":["MLB-CHILD-1","MLB-CHILD-2"]}`)
		case "/products/MLB-CHILD-1/items":
			_, _ = io.WriteString(w, `{"paging":{"total":1,"offset":0,"limit":50},"results":[{"seller_id":"seller-1","price":100.00,"currency_id":"BRL","condition":"new","shipping":{"mode":"me2"}}]}`)
		case "/products/MLB-CHILD-2/items":
			_, _ = io.WriteString(w, `{"paging":{"total":1,"offset":0,"limit":50},"results":[{"seller_id":"seller-2","price":95.50,"currency_id":"BRL","condition":"new","shipping":{"mode":"me1"}}]}`)
		default:
			t.Fatalf("unexpected request = %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	offers, err := catalogOffersTestAdapter(server.URL, true).ListCatalogOffers(context.Background(), pricingAccountRef(), "MLB-PARENT")
	if err != nil {
		t.Fatalf("ListCatalogOffers() error = %v", err)
	}
	if len(offers) != 2 || offers[0].SellerID != "seller-1" || offers[0].Price == nil || offers[0].Price.Amount != "100.00" || offers[1].SellerID != "seller-2" || offers[1].Price == nil || offers[1].Price.Amount != "95.50" {
		t.Fatalf("offers = %#v", offers)
	}
}

func catalogOffersTestAdapter(baseURL string, enabled bool) *CapabilityAdapter {
	return NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:              baseURL,
		HTTPClient:           &http.Client{Timeout: time.Second},
		CatalogOffersEnabled: enabled,
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "test-token", nil
		},
	})
}
