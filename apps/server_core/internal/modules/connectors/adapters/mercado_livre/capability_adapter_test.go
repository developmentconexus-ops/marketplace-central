package mercadolivre

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/connectors/ports"
)

func TestProviderCapabilitySetDeclaresMercadoLivreIdentityAnchors(t *testing.T) {
	t.Parallel()

	got := NewCapabilityAdapter(CapabilityAdapterConfig{}).ProviderCapabilitySet().IdentityAnchors
	want := []ports.IdentityAnchor{
		ports.IdentityAnchorSellerSKU,
		ports.IdentityAnchorEAN,
		ports.IdentityAnchorTitle,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("IdentityAnchors = %#v, want %#v", got, want)
	}
}

func TestCapabilityAdapterReadListingPublishesCanonicalFactsWithoutProviderPayload(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/items/MLB-DECIMAL" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{
			"id":"MLB-DECIMAL",
			"title":"Produto decimal",
			"status":"paused",
			"available_quantity":23,
			"price":89.00,
			"currency_id":"BRL",
			"listing_type_id":"gold_pro",
			"variations":[
				{"id":101,"available_quantity":7},
				{"id":102,"available_quantity":11}
			]
		}`)
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-decimal", nil
		},
	})

	snapshot, err := adapter.ReadListing(context.Background(), listingRef("MLB-DECIMAL", ""))
	if err != nil {
		t.Fatalf("ReadListing() error = %v", err)
	}
	if snapshot.PriceAmount == nil || *snapshot.PriceAmount != "89.00" {
		t.Fatalf("price amount = %#v, want exact decimal 89.00", snapshot.PriceAmount)
	}
	if snapshot.PriceCurrency != "BRL" || snapshot.ListingTypeCode != "gold_pro" {
		t.Fatalf("published facts = %#v", snapshot)
	}
	if snapshot.ProviderStatus != "paused" || snapshot.AvailableQuantity == nil || *snapshot.AvailableQuantity != 23 {
		t.Fatalf("item facts = %#v", snapshot)
	}
	if len(snapshot.Variations) != 2 || snapshot.Variations[0].AvailableQuantity == nil || *snapshot.Variations[0].AvailableQuantity != 7 || snapshot.Variations[1].AvailableQuantity == nil || *snapshot.Variations[1].AvailableQuantity != 11 {
		t.Fatalf("variation facts = %#v", snapshot.Variations)
	}
	if raw, ok := snapshot.RawProviderRef.(string); !ok || raw != "/items/MLB-DECIMAL" {
		t.Fatalf("raw provider ref = %#v, want sanitized path string", snapshot.RawProviderRef)
	}
}

func TestCapabilityAdapterClassifiesProviderAuthFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		resolver   AccessTokenResolver
		wantCalls  int
	}{
		{name: "http 401", statusCode: http.StatusUnauthorized, wantCalls: 1},
		{name: "http 403", statusCode: http.StatusForbidden, wantCalls: 1},
		{name: "credential unavailable", resolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "", errors.New("credential unavailable")
		}, wantCalls: 0},
		{name: "empty credential", resolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "", nil
		}, wantCalls: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			requestCalls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestCalls++
				if r.URL.Path != "/items/MLB-AUTH" {
					t.Fatalf("unexpected path %q", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			resolver := tt.resolver
			if resolver == nil {
				resolver = func(context.Context, domain.ProviderAccountRef) (string, error) {
					return "token-auth", nil
				}
			}
			adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
				BaseURL:             server.URL,
				HTTPClient:          server.Client(),
				AccessTokenResolver: resolver,
			})

			_, err := adapter.ReadListing(context.Background(), listingRef("MLB-AUTH", ""))
			if got := domain.ErrorCodeOf(err); got != domain.ErrCodeProviderAuth {
				t.Fatalf("error code = %q, err = %v, want provider auth", got, err)
			}
			if requestCalls != tt.wantCalls {
				t.Fatalf("provider request calls = %d, want %d", requestCalls, tt.wantCalls)
			}
		})
	}
}

func TestCapabilityAdapterReadListingMapsItemWithoutVariation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/items/MLB1" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token-1" {
			t.Fatalf("authorization = %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":                  "MLB1",
			"title":               "Martelo",
			"status":              "active",
			"available_quantity":  7,
			"last_updated":        "2026-07-06T10:00:00Z",
			"seller_custom_field": "SKU-1",
			"attributes": []map[string]any{
				{"id": "EAN", "value_name": "789000000001"},
			},
		})
	}))
	defer server.Close()

	frozen := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		Now:        func() time.Time { return frozen },
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-1", nil
		},
	})

	snapshot, err := adapter.ReadListing(context.Background(), listingRef("MLB1", ""))
	if err != nil {
		t.Fatalf("ReadListing() error = %v", err)
	}

	if snapshot.ProviderItemID != "MLB1" || snapshot.Title != "Martelo" || snapshot.SellerSKU != "SKU-1" {
		t.Fatalf("unexpected snapshot %#v", snapshot)
	}
	if snapshot.AvailableQuantity == nil || *snapshot.AvailableQuantity != 7 {
		t.Fatalf("available quantity = %#v", snapshot.AvailableQuantity)
	}
	if snapshot.EAN != "789000000001" {
		t.Fatalf("EAN = %q", snapshot.EAN)
	}
	if !snapshot.FetchedAt.Equal(frozen) {
		t.Fatalf("FetchedAt = %s", snapshot.FetchedAt)
	}
}

func TestCapabilityAdapterProbeAccountReadsUsersMe(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/me" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token-1" {
			t.Fatalf("authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":691607102,"nickname":"METALNOBREACABAMENTOS","site_id":"MLB","status":{"site_status":"active"}}`))
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-1", nil
		},
		Now: func() time.Time { return time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC) },
	})

	snapshot, err := adapter.ProbeAccount(context.Background(), domain.ProviderAccountRef{
		TenantID:          "tenant_default",
		InstallationID:    "inst-1",
		ProviderCode:      "mercado_livre",
		ProviderAccountID: "691607102",
	})
	if err != nil {
		t.Fatalf("ProbeAccount() error = %v", err)
	}
	if snapshot.ProviderAccountID != "691607102" {
		t.Fatalf("account id = %q", snapshot.ProviderAccountID)
	}
	if snapshot.ProviderAccountName != "METALNOBREACABAMENTOS" {
		t.Fatalf("account name = %q", snapshot.ProviderAccountName)
	}
	if snapshot.Status != "active" {
		t.Fatalf("status = %q", snapshot.Status)
	}
}

func TestCapabilityAdapterReadFeeQuoteReadsListingPrices(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sites/MLB/listing_prices" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		query := r.URL.Query()
		if query.Get("price") != "100" || query.Get("listing_type_id") != "gold_special" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"listing_type_id":"gold_special","sale_fee_details":{"percentage_fee":11,"fixed_fee":6}}]`))
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		SiteID:     "MLB",
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-1", nil
		},
		Now: func() time.Time { return time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC) },
	})

	quote, err := adapter.ReadFeeQuote(context.Background(), domain.FeeQuoteInput{
		AccountRef: domain.ProviderAccountRef{
			TenantID:          "tenant_default",
			InstallationID:    "inst-1",
			ProviderCode:      "mercado_livre",
			ProviderAccountID: "691607102",
		},
		SiteID:        "MLB",
		ListingTypeID: "gold_special",
		PriceAmount:   100,
		CurrencyID:    "BRL",
	})
	if err != nil {
		t.Fatalf("ReadFeeQuote() error = %v", err)
	}
	if quote.CommissionPercent == nil || *quote.CommissionPercent != 11 {
		t.Fatalf("commission percent = %v", quote.CommissionPercent)
	}
	if quote.FixedFeeAmount == nil || *quote.FixedFeeAmount != 6 {
		t.Fatalf("fixed fee = %v", quote.FixedFeeAmount)
	}
}

func TestCapabilityAdapterReadFeeQuoteAcceptsSingleObjectPayload(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sites/MLB/listing_prices" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		query := r.URL.Query()
		if query.Get("price") != "199.9" || query.Get("listing_type_id") != "gold_special" || query.Get("category_id") != "MLB1000" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"currency_id":"BRL","listing_type_id":"gold_special","sale_fee_amount":25.99,"sale_fee_details":{"percentage_fee":13,"fixed_fee":0}}`))
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		SiteID:     "MLB",
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-1", nil
		},
		Now: func() time.Time { return time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC) },
	})

	quote, err := adapter.ReadFeeQuote(context.Background(), domain.FeeQuoteInput{
		AccountRef: domain.ProviderAccountRef{
			TenantID:          "tenant_default",
			InstallationID:    "inst-1",
			ProviderCode:      "mercado_livre",
			ProviderAccountID: "691607102",
		},
		SiteID:        "MLB",
		CategoryID:    "MLB1000",
		ListingTypeID: "gold_special",
		PriceAmount:   199.9,
		CurrencyID:    "BRL",
	})
	if err != nil {
		t.Fatalf("ReadFeeQuote() error = %v", err)
	}
	if quote.CommissionPercent == nil || *quote.CommissionPercent != 13 {
		t.Fatalf("commission percent = %v", quote.CommissionPercent)
	}
	if quote.FixedFeeAmount == nil || *quote.FixedFeeAmount != 0 {
		t.Fatalf("fixed fee = %v", quote.FixedFeeAmount)
	}
}

func TestCapabilityAdapterDefaultClockAdvancesBetweenCalls(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/me" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":691607102,"nickname":"METALNOBREACABAMENTOS","site_id":"MLB","status":{"site_status":"active"}}`))
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-1", nil
		},
	})

	first, err := adapter.ProbeAccount(context.Background(), domain.ProviderAccountRef{
		TenantID:          "tenant_default",
		InstallationID:    "inst-1",
		ProviderCode:      "mercado_livre",
		ProviderAccountID: "691607102",
	})
	if err != nil {
		t.Fatalf("first ProbeAccount() error = %v", err)
	}

	time.Sleep(5 * time.Millisecond)

	second, err := adapter.ProbeAccount(context.Background(), domain.ProviderAccountRef{
		TenantID:          "tenant_default",
		InstallationID:    "inst-1",
		ProviderCode:      "mercado_livre",
		ProviderAccountID: "691607102",
	})
	if err != nil {
		t.Fatalf("second ProbeAccount() error = %v", err)
	}

	if !second.FetchedAt.After(first.FetchedAt) {
		t.Fatalf("fetched_at did not advance: first=%s second=%s", first.FetchedAt, second.FetchedAt)
	}
}

func TestCapabilityAdapterReadStockMapsVariation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/items/MLB2" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":           "MLB2",
			"title":        "Camiseta",
			"status":       "active",
			"last_updated": "2026-07-06T10:00:00Z",
			"variations": []map[string]any{
				{
					"id":                  101,
					"available_quantity":  4,
					"seller_custom_field": "SKU-VAR-101",
					"attributes": []map[string]any{
						{"id": "GTIN", "value_name": "789000000101"},
					},
				},
			},
		})
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-2", nil
		},
	})

	snapshot, err := adapter.ReadStock(context.Background(), listingRef("MLB2", "101"))
	if err != nil {
		t.Fatalf("ReadStock() error = %v", err)
	}

	if snapshot.Scope != domain.StockScopeVariation {
		t.Fatalf("scope = %q", snapshot.Scope)
	}
	if snapshot.ProviderVariationID != "101" || snapshot.AvailableQuantity != 4 {
		t.Fatalf("unexpected stock snapshot %#v", snapshot)
	}
	if snapshot.SellerSKU != "SKU-VAR-101" || snapshot.EAN != "789000000101" {
		t.Fatalf("unexpected variation mapping %#v", snapshot)
	}
}

func TestCapabilityAdapterReadOrderMapsSaleFee(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/orders/2001" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":            2001,
			"status":        "paid",
			"status_detail": "accredited",
			"last_updated":  "2026-07-06T11:00:00Z",
			"shipping": map[string]any{
				"id": 9001,
			},
			"order_items": []map[string]any{
				{
					"quantity":   2,
					"unit_price": 62.96,
					"sale_fee":   11.07,
					"item": map[string]any{
						"id":           "MLB3737188005",
						"title":        "Espacador",
						"variation_id": 12345,
						"seller_sku":   "SKU-ORDER-1",
						"variation_attributes": []map[string]any{
							{"id": "EAN", "value_name": "789000000201"},
						},
					},
				},
			},
			"payments": []map[string]any{
				{
					"id":                91776699099,
					"status":            "approved",
					"total_paid_amount": 125.92,
				},
			},
		})
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-3", nil
		},
	})

	snapshot, err := adapter.ReadOrder(context.Background(), orderRef("2001"))
	if err != nil {
		t.Fatalf("ReadOrder() error = %v", err)
	}

	if snapshot.ProviderStatus != "paid" || snapshot.ShippingID != "9001" {
		t.Fatalf("unexpected order snapshot %#v", snapshot)
	}
	if snapshot.SaleFeeAmount == nil || *snapshot.SaleFeeAmount != 11.07 {
		t.Fatalf("sale fee amount = %#v", snapshot.SaleFeeAmount)
	}
	if len(snapshot.Items) != 1 || snapshot.Items[0].ProviderVariationID != "12345" {
		t.Fatalf("unexpected items %#v", snapshot.Items)
	}
	if len(snapshot.Payments) != 1 || snapshot.Payments[0].Amount == nil || *snapshot.Payments[0].Amount != 125.92 {
		t.Fatalf("unexpected payments %#v", snapshot.Payments)
	}
}

func TestCapabilityAdapterReadOrderOmitsNilVariationID(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/orders/2002" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":           2002,
			"status":       "cancelled",
			"last_updated": "2026-07-06T10:00:00Z",
			"shipping": map[string]any{
				"id": 9002,
			},
			"order_items": []map[string]any{
				{
					"quantity":   1,
					"unit_price": 19.9,
					"item": map[string]any{
						"id":           "MLB4834373620",
						"variation_id": nil,
						"seller_sku":   "S-18",
						"title":        "Placa",
					},
				},
			},
			"payments": []map[string]any{
				{
					"id":                91776699100,
					"status":            "refunded",
					"total_paid_amount": 19.9,
				},
			},
		})
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-7", nil
		},
	})

	snapshot, err := adapter.ReadOrder(context.Background(), orderRef("2002"))
	if err != nil {
		t.Fatalf("ReadOrder() error = %v", err)
	}
	if len(snapshot.Items) != 1 {
		t.Fatalf("items = %#v", snapshot.Items)
	}
	if snapshot.Items[0].ProviderVariationID != "" {
		t.Fatalf("provider_variation_id = %q, want empty string for nil upstream value", snapshot.Items[0].ProviderVariationID)
	}
}

func TestCapabilityAdapterUpdateAvailableQuantityReturnsRejectedOnProviderValidation(t *testing.T) {
	t.Parallel()

	var putBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/items/MLB3":
			if r.Method == http.MethodGet {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id":                 "MLB3",
					"title":              "Item 3",
					"status":             "active",
					"available_quantity": 2,
				})
				return
			}
			raw, _ := io.ReadAll(r.Body)
			putBody = string(raw)
			http.Error(w, `{"message":"quantity not allowed"}`, http.StatusBadRequest)
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-4", nil
		},
	})

	result, err := adapter.UpdateAvailableQuantity(context.Background(), stockWriteRequest("MLB3", "", "action-1", 0))
	if err != nil {
		t.Fatalf("UpdateAvailableQuantity() error = %v", err)
	}
	if result.Result != domain.StockWriteResultRejected {
		t.Fatalf("result = %#v", result)
	}
	if !strings.Contains(putBody, `"available_quantity":0`) {
		t.Fatalf("unexpected put body %q", putBody)
	}
}

func TestCapabilityAdapterUpdateAvailableQuantityMapsRateLimit(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":                 "MLB4",
				"available_quantity": 1,
			})
		case http.MethodPut:
			w.WriteHeader(http.StatusTooManyRequests)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-5", nil
		},
	})

	_, err := adapter.UpdateAvailableQuantity(context.Background(), stockWriteRequest("MLB4", "", "action-2", 5))
	if got := domain.ErrorCodeOf(err); got != domain.ErrCodeProviderRateLimited {
		t.Fatalf("error code = %q, err = %v", got, err)
	}
}

func TestCapabilityAdapterUpdateAvailableQuantityRejectsUnsupportedVariationShape(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "MLB5",
			"variations": []map[string]any{
				{"id": 987, "available_quantity": 1},
			},
		})
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-6", nil
		},
	})

	_, err := adapter.UpdateAvailableQuantity(context.Background(), stockWriteRequest("MLB5", "", "action-3", 3))
	if got := domain.ErrorCodeOf(err); got != domain.ErrCodeProviderUnsupportedShape {
		t.Fatalf("error code = %q, err = %v", got, err)
	}
}

func listingRef(itemID string, variationID string) domain.ProviderListingRef {
	return domain.ProviderListingRef{
		AccountRef: domain.ProviderAccountRef{
			TenantID:          "tenant-1",
			InstallationID:    "inst-1",
			ProviderCode:      "mercado_livre",
			ProviderAccountID: "seller-1",
		},
		ProviderItemID:      itemID,
		ProviderVariationID: variationID,
	}
}

func orderRef(orderID string) domain.ProviderOrderRef {
	return domain.ProviderOrderRef{
		AccountRef: domain.ProviderAccountRef{
			TenantID:          "tenant-1",
			InstallationID:    "inst-1",
			ProviderCode:      "mercado_livre",
			ProviderAccountID: "seller-1",
		},
		ProviderOrderID: orderID,
	}
}

var testAccountRef = domain.ProviderAccountRef{
	TenantID:          "tenant-1",
	InstallationID:    "inst-1",
	ProviderCode:      "mercado_livre",
	ProviderAccountID: "seller-1",
}

func newTestAdapter(t *testing.T, baseURL string) *CapabilityAdapter {
	t.Helper()
	return NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-test", nil
		},
	})
}

// TestListOrderRefsDoesNotHydrateEachOrder is the negative control for F-00 Task 3: before this
// change, ListOrders (the only enumeration path) called ReadOrder — one /orders/{id} request —
// per hit off /orders/search. The batch consumer (orders.ImportService) only ever used the id, so
// that hydration was a wasted provider call per order per cycle. ListOrderRefs must enumerate
// without ever reaching /orders/{id}.
func TestListOrderRefsDoesNotHydrateEachOrder(t *testing.T) {
	t.Parallel()

	var readOrderCalls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/orders/search"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"results":[
				{"id":2000000000001,"date_last_updated":"2026-08-02T10:00:00.000-03:00"},
				{"id":2000000000002,"date_last_updated":"2026-08-02T11:00:00.000-03:00"}]}`)
		default:
			atomic.AddInt32(&readOrderCalls, 1)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	adapter := newTestAdapter(t, srv.URL)
	refs, err := adapter.ListOrderRefs(context.Background(), domain.ListOrdersInput{
		AccountRef: testAccountRef, Limit: 50,
	})
	if err != nil {
		t.Fatalf("list order refs: %v", err)
	}
	if got := atomic.LoadInt32(&readOrderCalls); got != 0 {
		t.Fatalf("enumeracao id-only chamou /orders/{id} %d vez(es); quero 0", got)
	}
	if len(refs) != 2 {
		t.Fatalf("refs: quero 2, recebi %d", len(refs))
	}
	if refs[0].ProviderOrderID != "2000000000001" {
		t.Fatalf("id: quero \"2000000000001\", recebi %q", refs[0].ProviderOrderID)
	}
	if refs[0].ProviderUpdatedAt == nil {
		t.Fatalf("date_last_updated presente no payload tem que chegar no ref")
	}
	if !refs[0].ProviderUpdatedAt.Equal(time.Date(2026, 8, 2, 13, 0, 0, 0, time.UTC)) {
		t.Fatalf("date_last_updated: quero 13:00Z (10:00-03:00), recebi %s", refs[0].ProviderUpdatedAt)
	}
}

// TestListOrderRefsLeavesUpdatedAtNilWhenProviderOmitsIt is the unknown/absent guard (ADR-17):
// absence of date_last_updated in the payload is a DIFFERENT state from "very old" or "zero time".
// A zero time.Time would read as January 1st, year 1, entering a cursor watermark as a fabricated
// fact instead of an honest unknown.
func TestListOrderRefsLeavesUpdatedAtNilWhenProviderOmitsIt(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/orders/search"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"results":[{"id":2000000000003}]}`)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	adapter := newTestAdapter(t, srv.URL)
	refs, err := adapter.ListOrderRefs(context.Background(), domain.ListOrdersInput{
		AccountRef: testAccountRef, Limit: 50,
	})
	if err != nil {
		t.Fatalf("list order refs: %v", err)
	}
	if len(refs) != 1 {
		t.Fatalf("refs: quero 1, recebi %d", len(refs))
	}
	if refs[0].ProviderOrderID != "2000000000003" {
		t.Fatalf("id: quero \"2000000000003\", recebi %q", refs[0].ProviderOrderID)
	}
	if refs[0].ProviderUpdatedAt != nil {
		t.Fatalf("date_last_updated ausente no payload; ProviderUpdatedAt tinha que ser nil, veio %v", refs[0].ProviderUpdatedAt)
	}
}

func stockWriteRequest(itemID string, variationID string, actionID string, quantity int) domain.StockWriteRequest {
	return domain.StockWriteRequest{
		TenantID:            "tenant-1",
		InstallationID:      "inst-1",
		ProviderCode:        "mercado_livre",
		ProviderAccountID:   "seller-1",
		ProviderItemID:      itemID,
		ProviderVariationID: variationID,
		IdempotencyKey:      actionID,
		RequestedQuantity:   quantity,
		Reason:              "test",
	}
}
