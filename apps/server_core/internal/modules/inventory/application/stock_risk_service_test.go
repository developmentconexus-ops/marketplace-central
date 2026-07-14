package application

import (
	"context"
	"testing"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
	"marketplace-central/apps/server_core/internal/modules/inventory/ports"
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
	facts     map[int64]*ports.StockFact
	callCount int
	requested []int64
}

func (s *stubInternalStockReader) GetStockFactsByIDs(_ context.Context, ids []int64) (map[int64]*ports.StockFact, error) {
	s.callCount++
	s.requested = append(s.requested, ids...)
	return s.facts, nil
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
		&stubInternalStockReader{facts: map[int64]*ports.StockFact{101: {
			ProductID: 101,
			Quantity:  testFloatPtr(8),
			Source:    internalreaddomain.SourceMetadata{FetchedAt: now},
		}}},
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

func TestStockRiskServiceUsesOneBatchCallAndFlagsMissingStock(t *testing.T) {
	now := time.Now().UTC()
	stockReader := &stubInternalStockReader{facts: map[int64]*ports.StockFact{}}
	service := NewStockRiskService(
		stubSnapshotReader{items: []domain.ListingSnapshot{{
			Identity:          domain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"},
			AvailableQuantity: testIntPtr(2),
			ObservedAt:        &now,
		}}},
		stubLinkReader{items: []domain.LinkRecord{{
			Identity:          domain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"},
			State:             domain.LinkStateResolved,
			InternalProductID: testIntPtr(101),
		}}},
		stockReader,
	)

	items, err := service.List(context.Background(), ListStockRisksInput{InstallationID: "inst-1"})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if stockReader.callCount != 1 || len(stockReader.requested) != 1 || stockReader.requested[0] != 101 {
		t.Fatalf("batch call = %d with ids %v, want one call for [101]", stockReader.callCount, stockReader.requested)
	}
	if len(items) != 1 || items[0].InternalQuantity != nil {
		t.Fatalf("items = %+v, want one item with nil internal quantity", items)
	}
	if len(items[0].QualityFlags) != 1 || items[0].QualityFlags[0] != "missing_stock" {
		t.Fatalf("quality flags = %v, want [missing_stock]", items[0].QualityFlags)
	}
}

func testIntPtr(value int) *int {
	return &value
}

func testFloatPtr(value float64) *float64 {
	return &value
}
