package mercadolivre

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/connectors/ports"
)

func TestCapabilityAdapterUpdatePriceCapturesAbsoluteMLWrite(t *testing.T) {
	t.Parallel()

	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != http.MethodPut || r.URL.Path != "/items/MLB-PRICE-1" {
			t.Errorf("request = %s %s, want PUT /items/MLB-PRICE-1", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("authorization = %q, want bearer token", got)
		}
		if got := r.Header.Get("X-Idempotency-Key"); got != "MP-000001:MLB-PRICE-1" {
			t.Errorf("idempotency key = %q", got)
		}
		var body struct {
			Price      json.Number `json:"price"`
			CurrencyID string      `json:"currency_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		if body.Price.String() != "49.90" || body.CurrencyID != "BRL" {
			t.Errorf("request body = %+v, want price 49.90 and currency BRL", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"MLB-PRICE-1","price":49.90,"currency_id":"BRL"}`)
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "test-token", nil
		},
	})

	result, err := adapter.UpdatePrice(context.Background(), domain.PriceWriteRequest{
		TenantID:       "tenant-1",
		InstallationID: "installation-1",
		ListingID:      "MLB-PRICE-1",
		IdempotencyKey: "MP-000001:MLB-PRICE-1",
		Price:          domain.Price{Amount: "49.90", Currency: "BRL"},
	})
	if err != nil {
		t.Fatalf("UpdatePrice() error = %v", err)
	}
	if calls != 1 {
		t.Fatalf("provider calls = %d, want 1", calls)
	}
	if result.Result != domain.WriteResultApplied || result.ListingID != "MLB-PRICE-1" || result.IdempotencyKey != "MP-000001:MLB-PRICE-1" || result.Price.Amount != "49.90" || result.Price.Currency != "BRL" {
		t.Fatalf("result = %+v", result)
	}
}

func TestCapabilityAdapterUpdatePriceRejectsMalformedCanonicalInputBeforeNetwork(t *testing.T) {
	t.Parallel()

	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls++ }))
	defer server.Close()
	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "must-not-be-used", nil
		},
	})

	tests := []domain.PriceWriteRequest{
		{TenantID: "tenant", InstallationID: "installation", ListingID: "MLB-1", IdempotencyKey: "key", Price: domain.Price{Amount: "0", Currency: "BRL"}},
		{TenantID: "tenant", InstallationID: "installation", ListingID: "MLB-1", IdempotencyKey: "key", Price: domain.Price{Amount: "49,90", Currency: "BRL"}},
		{TenantID: "tenant", InstallationID: "installation", ListingID: "MLB-1", IdempotencyKey: "key", Price: domain.Price{Amount: "49.90", Currency: ""}},
		{TenantID: "tenant", InstallationID: "installation", ListingID: "", IdempotencyKey: "key", Price: domain.Price{Amount: "49.90", Currency: "BRL"}},
		{TenantID: "tenant", InstallationID: "installation", ListingID: "MLB-1", IdempotencyKey: "", Price: domain.Price{Amount: "49.90", Currency: "BRL"}},
	}
	for _, request := range tests {
		if _, err := adapter.UpdatePrice(context.Background(), request); domain.ErrorCodeOf(err) != domain.ErrCodeProviderValidation {
			t.Errorf("request %+v error = %v, want provider validation", request, err)
		}
	}
	if calls != 0 {
		t.Fatalf("provider calls = %d, want 0", calls)
	}
}

func TestCapabilityAdapterUpdatePriceMapsProviderFailuresWithSanitizedMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code int
		want domain.ErrorCode
	}{
		{name: "auth 401", code: http.StatusUnauthorized, want: domain.ErrCodeProviderAuth},
		{name: "auth 403", code: http.StatusForbidden, want: domain.ErrCodeProviderAuth},
		{name: "rate limit", code: http.StatusTooManyRequests, want: domain.ErrCodeProviderRateLimited},
		{name: "validation", code: http.StatusBadRequest, want: domain.ErrCodeProviderValidation},
		{name: "unavailable", code: http.StatusBadGateway, want: domain.ErrCodeProviderTransient},
		{name: "unknown", code: http.StatusFound, want: domain.ErrorCode("CONNECTORS_INTERNAL")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.WriteHeader(tt.code)
				_, _ = io.WriteString(w, `{"message":"provider rejected price; token=super-secret","error":"provider_error"}`)
			}))
			defer server.Close()
			adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
				BaseURL:    server.URL,
				HTTPClient: server.Client(),
				AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
					return "request-token", nil
				},
			})

			_, err := adapter.UpdatePrice(context.Background(), validPriceRequest())
			if got := domain.ErrorCodeOf(err); got != tt.want {
				t.Fatalf("error code = %q, err = %v, want %q", got, err, tt.want)
			}
			if calls != 1 {
				t.Fatalf("provider calls = %d, want 1", calls)
			}
			if strings.Contains(err.Error(), "super-secret") || strings.Contains(err.Error(), "request-token") || strings.Contains(err.Error(), `{"message"`) {
				t.Fatalf("provider error was not sanitized: %v", err)
			}
			if !strings.Contains(err.Error(), "provider rejected price") {
				t.Fatalf("sanitized provider message missing: %v", err)
			}
		})
	}
}

func TestCapabilityAdapterUpdatePriceMapsTimeoutAsUnavailable(t *testing.T) {
	t.Parallel()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL: "http://provider.invalid",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, context.DeadlineExceeded
		})},
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "request-token", nil
		},
	})

	_, err := adapter.UpdatePrice(context.Background(), validPriceRequest())
	if got := domain.ErrorCodeOf(err); got != domain.ErrCodeProviderTransient {
		t.Fatalf("error code = %q, err = %v, want provider transient", got, err)
	}
	if strings.Contains(err.Error(), "request-token") || strings.Contains(err.Error(), "DeadlineExceeded") {
		t.Fatalf("timeout error was not sanitized: %v", err)
	}
}

func TestCapabilityAdapterExposesPriceWriterCapability(t *testing.T) {
	t.Parallel()
	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{})
	capabilities := adapter.ProviderCapabilitySet()
	if capabilities.PriceWrites != adapter {
		t.Fatalf("PriceWrites = %T, want adapter", capabilities.PriceWrites)
	}
	var _ ports.PriceWriter = adapter
}

func validPriceRequest() domain.PriceWriteRequest {
	return domain.PriceWriteRequest{
		TenantID:       "tenant-1",
		InstallationID: "installation-1",
		ListingID:      "MLB-PRICE-1",
		IdempotencyKey: "MP-000001:MLB-PRICE-1",
		Price:          domain.Price{Amount: "49.90", Currency: "BRL"},
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
