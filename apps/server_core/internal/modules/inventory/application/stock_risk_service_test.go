package application

import (
	"context"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
)

type stubSnapshotReader struct {
	items []domain.ListingSnapshot
}

func (s stubSnapshotReader) ListListingSnapshots(context.Context, string, int) ([]domain.ListingSnapshot, error) {
	return s.items, nil
}

type stubLinkReader struct {
	items []domain.LinkRecord
}

func (s stubLinkReader) ListLinks(context.Context, string, int) ([]domain.LinkRecord, error) {
	return s.items, nil
}

type stubInternalStockReader struct {
	stock   domain.InternalStockEvidence
	product domain.ProductEvidence
}

func (s stubInternalStockReader) GetSellableStock(context.Context, int, domain.StockPolicy) (domain.InternalStockEvidence, domain.ProductEvidence, error) {
	return s.stock, s.product, nil
}

func TestStockRiskServiceClassifiesOversellAndFilters(t *testing.T) {
	now := time.Now().UTC()
	service := NewStockRiskService(
		stubSnapshotReader{items: []domain.ListingSnapshot{{
			Identity: domain.ListingIdentity{
				InstallationID: "inst-1",
				ProviderItemID: "MLB123",
			},
			ProviderCode:      "mercado_livre",
			Title:             "Produto teste",
			AvailableQuantity: testIntPtr(9),
			ObservedAt:        &now,
		}}},
		stubLinkReader{items: []domain.LinkRecord{{
			Identity: domain.ListingIdentity{
				InstallationID: "inst-1",
				ProviderItemID: "MLB123",
			},
			State:             domain.LinkStateResolved,
			InternalProductID: testIntPtr(101),
		}}},
		stubInternalStockReader{
			stock: domain.InternalStockEvidence{
				Quantity: testIntPtr(8),
				Source:   domain.SourceEvidence{ObservedAt: &now},
			},
			product: domain.ProductEvidence{ProductID: 101},
		},
	)

	items, err := service.List(context.Background(), ListStockRisksInput{
		InstallationID: "inst-1",
		State:          "oversell",
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items)=%d, want 1", len(items))
	}
	if items[0].State != domain.StockRiskOversell {
		t.Fatalf("state=%q, want oversell", items[0].State)
	}
	if !items[0].Actionable {
		t.Fatalf("actionable=false, want true")
	}
	if items[0].RecommendedQuantity == nil || *items[0].RecommendedQuantity != 7 {
		t.Fatalf("recommended=%v, want 7", items[0].RecommendedQuantity)
	}
}

func testIntPtr(value int) *int {
	return &value
}
