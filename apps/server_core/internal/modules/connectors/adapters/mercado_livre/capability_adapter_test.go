package mercadolivre

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

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
