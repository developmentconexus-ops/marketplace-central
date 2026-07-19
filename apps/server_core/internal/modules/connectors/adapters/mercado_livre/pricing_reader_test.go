package mercadolivre

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/connectors/ports"
)

var _ ports.MarketReader = (*CapabilityAdapter)(nil)

func TestOwnItemPricing(t *testing.T) {
	t.Parallel()

	frozen := time.Date(2026, 7, 17, 12, 30, 0, 123456000, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ML /items/{id}/sale_price 400s without context; channel_marketplace is
		// the required marketplace-pricing context (FINDING-M02-LIVE-2, official
		// docs precio-por-cantidad/kits-virtuales).
		if r.Method != http.MethodGet || r.URL.Path != "/items/MLB-PRICE/sale_price" || r.URL.Query().Get("context") != "channel_marketplace" {
			t.Fatalf("request = %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("authorization = %q", got)
		}
		_, _ = io.WriteString(w, `{"amount":89.90,"currency_id":"BRL","regular_amount":120.00}`)
	}))
	defer server.Close()

	adapter := pricingTestAdapter(server.URL, frozen)
	pricing, err := adapter.GetOwnItemPricing(context.Background(), pricingAccountRef(), "MLB-PRICE")
	if err != nil {
		t.Fatalf("GetOwnItemPricing() error = %v", err)
	}
	if pricing.ItemID != "MLB-PRICE" || pricing.Amount == nil || *pricing.Amount != "89.90" || pricing.Currency != "BRL" || pricing.RegularAmount == nil || *pricing.RegularAmount != "120.00" || !pricing.FetchedAt.Equal(frozen) {
		t.Fatalf("pricing = %#v", pricing)
	}
}

func TestOwnItemPricingKeepsAbsentPricesUnknown(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"currency_id":"BRL","amount":null,"regular_amount":null}`)
	}))
	defer server.Close()

	pricing, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetOwnItemPricing(context.Background(), pricingAccountRef(), "MLB-PRICE")
	if err != nil {
		t.Fatalf("GetOwnItemPricing() error = %v", err)
	}
	if pricing.Amount != nil || pricing.RegularAmount != nil {
		t.Fatalf("pricing = %#v, want nil amounts", pricing)
	}
}

func TestPriceToWin(t *testing.T) {
	t.Parallel()

	frozen := time.Date(2026, 7, 17, 12, 30, 0, 123456000, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/items/MLB-PRICE/price_to_win" || r.URL.Query().Get("version") != "v2" {
			t.Fatalf("request = %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
		_, _ = io.WriteString(w, `{"status":"winning","current_price":120.00,"price_to_win":99.50,"currency_id":"BRL","position":{"rank":2,"total":5},"min_price":80.00}`)
	}))
	defer server.Close()

	result, err := pricingTestAdapter(server.URL, frozen).GetPriceToWin(context.Background(), pricingAccountRef(), "MLB-PRICE")
	if err != nil {
		t.Fatalf("GetPriceToWin() error = %v", err)
	}
	if result.ItemID != "MLB-PRICE" || result.Status != "winning" || result.CurrentPrice == nil || result.CurrentPrice.Amount != "120.00" || result.TargetPrice == nil || result.TargetPrice.Amount != "99.50" || result.TargetPrice.Amount == "80.00" || result.Position == nil || result.Position.Rank != 2 || result.Position.Total != 5 || !result.FetchedAt.Equal(frozen) {
		t.Fatalf("price to win = %#v", result)
	}
}

func TestPriceToWinKeepsAbsentValuesUnknown(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"status":"unknown","currency_id":"BRL"}`)
	}))
	defer server.Close()

	result, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetPriceToWin(context.Background(), pricingAccountRef(), "MLB-PRICE")
	if err != nil {
		t.Fatalf("GetPriceToWin() error = %v", err)
	}
	if result.CurrentPrice != nil || result.TargetPrice != nil || result.Position != nil {
		t.Fatalf("price to win = %#v, want nil values", result)
	}
}

func TestPriceToWinErrorMappingWithoutRetry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status int
		want   error
	}{
		{name: "unauthorized", status: http.StatusUnauthorized, want: connectorsdomain.ErrUnauthorized},
		{name: "forbidden", status: http.StatusForbidden, want: connectorsdomain.ErrUnauthorized},
		{name: "not found", status: http.StatusNotFound, want: connectorsdomain.ErrNotFound},
		{name: "rate limited", status: http.StatusTooManyRequests, want: connectorsdomain.ErrRateLimited},
		{name: "provider unavailable", status: http.StatusBadGateway, want: connectorsdomain.ErrProviderUnavailable},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.WriteHeader(test.status)
			}))
			defer server.Close()

			_, err := pricingTestAdapter(server.URL, time.Now().UTC()).GetPriceToWin(context.Background(), pricingAccountRef(), "MLB-PRICE")
			if !errors.Is(err, test.want) {
				t.Fatalf("error = %v, errors.Is(_, %v) = false", err, test.want)
			}
			if calls != 1 {
				t.Fatalf("provider calls = %d, want 1", calls)
			}
		})
	}
}

func TestOwnItemPricingTimeoutMapsToProviderUnavailable(t *testing.T) {
	t.Parallel()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    "http://provider.invalid",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, context.DeadlineExceeded })},
		AccessTokenResolver: func(context.Context, connectorsdomain.ProviderAccountRef) (string, error) {
			return "test-token", nil
		},
	})
	_, err := adapter.GetOwnItemPricing(context.Background(), pricingAccountRef(), "MLB-PRICE")
	if !errors.Is(err, connectorsdomain.ErrProviderUnavailable) {
		t.Fatalf("error = %v, want provider unavailable", err)
	}
}

func pricingTestAdapter(baseURL string, now time.Time) *CapabilityAdapter {
	return NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: time.Second},
		AccessTokenResolver: func(context.Context, connectorsdomain.ProviderAccountRef) (string, error) {
			return "test-token", nil
		},
		Now: func() time.Time { return now },
	})
}

func pricingAccountRef() connectorsdomain.ProviderAccountRef {
	return connectorsdomain.ProviderAccountRef{
		TenantID: "tenant-1", InstallationID: "installation-1", ProviderCode: "mercado_livre", ProviderAccountID: "seller-1",
	}
}
