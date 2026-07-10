package application

import (
	"context"
	"testing"
	"time"

	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type stubOrderSource struct {
	items []ordersdomain.OrderIngestionSnapshot
}

func (s stubOrderSource) ListOrders(context.Context, string, int) ([]ordersdomain.OrderIngestionSnapshot, error) {
	return s.items, nil
}

type stubLinkReader struct {
	links map[string]ordersdomain.ItemLink
}

func (s stubLinkReader) ResolveLinks(context.Context, string, []ordersdomain.ListingIdentity) (map[string]ordersdomain.ItemLink, error) {
	return s.links, nil
}

type stubOrderStore struct {
	orders []ordersdomain.MarketplaceOrder
}

func (s *stubOrderStore) UpsertOrders(_ context.Context, orders []ordersdomain.MarketplaceOrder) (int, int, error) {
	s.orders = orders
	return len(orders), 0, nil
}

func (s *stubOrderStore) ListOrders(context.Context, string, int) ([]ordersdomain.MarketplaceOrder, error) {
	return s.orders, nil
}

func TestImportServiceNormalizesLinkQualityAndPayments(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	store := &stubOrderStore{}
	service := NewImportService(ImportServiceConfig{
		Source: stubOrderSource{items: []ordersdomain.OrderIngestionSnapshot{{
			ProviderCode:         "mercado_livre",
			ProviderOrderID:      "2001",
			ProviderStatus:       "paid",
			ProviderStatusDetail: "accredited",
			ProviderCreatedAt:    &now,
			ProviderUpdatedAt:    &now,
			FetchedAt:            now,
			ShippingID:           "ship-1",
			Tags:                 []string{"paid"},
			Items: []ordersdomain.OrderIngestionItem{{
				ProviderItemID:      "MLB123",
				ProviderVariationID: "VAR1",
				SellerSKU:           "SKU-1",
				Title:               "Produto",
				Quantity:            2,
				SaleFeeAmount:       floatPtr(11.07),
			}},
			Payments: []ordersdomain.OrderIngestionPayment{{
				PaymentID:         "pay-1",
				Status:            "approved",
				TransactionAmount: floatPtr(125.92),
				TotalPaidAmount:   floatPtr(125.92),
			}},
		}}},
		Links: stubLinkReader{links: map[string]ordersdomain.ItemLink{
			"inst-1|MLB123|VAR1": {
				Quality:           ordersdomain.LinkQualityResolved,
				InternalProductID: intPtr(321),
			},
		}},
		Store: store,
		Now:   func() time.Time { return now },
	})

	result, err := service.Import(context.Background(), ImportOrdersInput{
		InstallationID: "inst-1",
		Limit:          1,
	})
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if result.ImportedCount != 1 {
		t.Fatalf("imported count = %d", result.ImportedCount)
	}
	if len(store.orders) != 1 {
		t.Fatalf("stored orders = %d", len(store.orders))
	}
	item := store.orders[0].Items[0]
	if item.LinkQuality != ordersdomain.LinkQualityResolved {
		t.Fatalf("link quality = %s", item.LinkQuality)
	}
	if item.InternalProductID == nil || *item.InternalProductID != 321 {
		t.Fatalf("internal product id = %#v", item.InternalProductID)
	}
	payment := store.orders[0].Payments[0]
	if payment.TotalPaidAmount == nil || *payment.TotalPaidAmount != 125.92 {
		t.Fatalf("total paid amount = %#v", payment.TotalPaidAmount)
	}
}

func TestImportServiceDefaultsMissingLinkQuality(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	closedAt := now.Add(-time.Hour)
	store := &stubOrderStore{}
	service := NewImportService(ImportServiceConfig{
		Source: stubOrderSource{items: []ordersdomain.OrderIngestionSnapshot{{
			ProviderCode:         "mercado_livre",
			ProviderOrderID:      "2002",
			ProviderStatus:       "cancelled",
			ProviderStatusDetail: "buyer_cancelled",
			ProviderClosedAt:     &closedAt,
			ProviderUpdatedAt:    &now,
			FetchedAt:            now,
			CancellationDetail:   "cancelled_by_buyer",
			Items: []ordersdomain.OrderIngestionItem{{
				ProviderItemID: "MLB999",
				Quantity:       1,
			}},
		}}},
		Links: stubLinkReader{links: map[string]ordersdomain.ItemLink{}},
		Store: store,
		Now:   func() time.Time { return now },
	})

	if _, err := service.Import(context.Background(), ImportOrdersInput{InstallationID: "inst-1", Limit: 1}); err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if got := store.orders[0].Items[0].LinkQuality; got != ordersdomain.LinkQualityMissing {
		t.Fatalf("link quality = %s", got)
	}
	stored := store.orders[0]
	if stored.ProviderStatus != "cancelled" || stored.ProviderStatusDetail != "buyer_cancelled" {
		t.Fatalf("cancelled status = %q/%q", stored.ProviderStatus, stored.ProviderStatusDetail)
	}
	if stored.ProviderClosedAt == nil || !stored.ProviderClosedAt.Equal(closedAt) {
		t.Fatalf("provider closed at = %v, want %v", stored.ProviderClosedAt, closedAt)
	}
	if stored.CancellationDetail != "cancelled_by_buyer" {
		t.Fatalf("cancellation detail = %q", stored.CancellationDetail)
	}
}

func TestImportServiceDerivesSafeProviderReference(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	store := &stubOrderStore{}
	service := NewImportService(ImportServiceConfig{
		Source: stubOrderSource{items: []ordersdomain.OrderIngestionSnapshot{{
			ProviderCode:    "mercado_livre",
			ProviderOrderID: "2003",
			FetchedAt:       now,
		}}},
		Links: stubLinkReader{links: map[string]ordersdomain.ItemLink{}},
		Store: store,
		Now:   func() time.Time { return now },
	})

	if _, err := service.Import(context.Background(), ImportOrdersInput{InstallationID: "inst-1", Limit: 1}); err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if got := store.orders[0].RawProviderRef; got != "/orders/2003" {
		t.Fatalf("raw provider reference = %#v, want exact safe reference", got)
	}
}

func floatPtr(v float64) *float64 { return &v }
func intPtr(v int) *int           { return &v }
