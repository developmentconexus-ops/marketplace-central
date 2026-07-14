package application

import (
	"context"
	"testing"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	profitabilitydomain "marketplace-central/apps/server_core/internal/modules/profitability/domain"
)

type countingBatchReader struct {
	stubInternalFactReader
	costIDs        []int64
	taxIDs         []int64
	costQueryCount int
	taxQueryCount  int
}

func (r *countingBatchReader) GetCostFactsByIDs(_ context.Context, ids []int64) (map[int64]*internalreaddomain.CostAsOf, error) {
	r.costIDs = append([]int64(nil), ids...)
	r.costQueryCount = (len(ids) + 499) / 500
	result := make(map[int64]*internalreaddomain.CostAsOf, len(ids))
	for _, id := range ids {
		result[id] = nil
	}
	return result, nil
}

func (r *countingBatchReader) GetTaxFactsByIDs(_ context.Context, ids []int64) (map[int64]*internalreaddomain.TaxInputs, error) {
	r.taxIDs = append([]int64(nil), ids...)
	r.taxQueryCount = (len(ids) + 499) / 500
	result := make(map[int64]*internalreaddomain.TaxInputs, len(ids))
	for _, id := range ids {
		result[id] = nil
	}
	return result, nil
}

func TestImportMarginInputsBatchQueryCountN1200(t *testing.T) {
	now := time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)
	items := make([]profitabilitydomain.OrderItemFact, 1200)
	for index := range items {
		productID := index + 1
		items[index] = profitabilitydomain.OrderItemFact{
			ProviderItemID:    "item",
			LinkQuality:       profitabilitydomain.OrderLinkResolved,
			InternalProductID: &productID,
			Quantity:          1,
		}
	}
	reader := &countingBatchReader{}
	service := NewService(ServiceConfig{
		Orders:   stubOrderReader{orders: []profitabilitydomain.OrderFact{{InstallationID: "inst-1", ProviderOrderID: "order-1", FetchedAt: now, Items: items}}},
		Internal: reader,
		Inputs:   &stubInputStore{},
		Now:      func() time.Time { return now },
	})
	if _, err := service.ImportMarginInputs(context.Background(), ImportMarginInputsInput{InstallationID: "inst-1", Limit: 200}); err != nil {
		t.Fatalf("ImportMarginInputs() error = %v", err)
	}
	if reader.costQueryCount != 3 || reader.taxQueryCount != 3 {
		t.Fatalf("batch query counts cost=%d tax=%d, want 3/3", reader.costQueryCount, reader.taxQueryCount)
	}
	if len(reader.costIDs) != 1200 || len(reader.taxIDs) != 1200 {
		t.Fatalf("batch IDs cost=%d tax=%d, want 1200/1200", len(reader.costIDs), len(reader.taxIDs))
	}
}

func TestBatchMissingFactsMakeProductIncomputable(t *testing.T) {
	now := time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)
	reader := &countingBatchReader{}
	price, fee := 100.0, 10.0
	productID := 42
	service := NewService(ServiceConfig{
		Orders:   stubOrderReader{orders: []profitabilitydomain.OrderFact{{InstallationID: "inst-1", ProviderOrderID: "order-1", FetchedAt: now, RealizationState: profitabilitydomain.OrderRealizationRealized, Items: []profitabilitydomain.OrderItemFact{{ProviderItemID: "item-1", UnitPrice: &price, SaleFeeAmount: &fee, Quantity: 1, LinkQuality: profitabilitydomain.OrderLinkResolved, InternalProductID: &productID}}}}},
		Internal: reader,
		Inputs:   &stubInputStore{},
		Now:      func() time.Time { return now },
	})
	result, err := service.ImportMarginInputs(context.Background(), ImportMarginInputsInput{InstallationID: "inst-1", Limit: 1})
	if err != nil {
		t.Fatalf("ImportMarginInputs() error = %v", err)
	}
	snapshots := buildProfitSnapshots(result.Items, nil, map[string]profitabilitydomain.OrderRealizationState{"order-1": profitabilitydomain.OrderRealizationRealized}, now)
	for _, snapshot := range snapshots {
		if snapshot.Scope != profitabilitydomain.InputScopeItem {
			continue
		}
		if snapshot.ContributionAmount != nil || snapshot.Quality != profitabilitydomain.ProfitSnapshotIncomplete {
			t.Fatalf("snapshot = %+v, want incomplete with no contribution", snapshot)
		}
		if !hasFlag(snapshot.Flags, profitabilitydomain.ProfitFlagMissingCost) || !hasFlag(snapshot.Flags, profitabilitydomain.ProfitFlagMissingTax) {
			t.Fatalf("snapshot flags = %v, want missing_cost and missing_tax", snapshot.Flags)
		}
		return
	}
	t.Fatal("expected item snapshot")
}

type countingOrderReader struct{ calls int }

func (r *countingOrderReader) ListOrders(context.Context, string, int) ([]profitabilitydomain.OrderFact, error) {
	r.calls++
	return nil, nil
}

func TestImportMarginInputsLimitExceededBeforeAnyPortCall(t *testing.T) {
	orders := &countingOrderReader{}
	reader := &countingBatchReader{}
	service := NewService(ServiceConfig{Orders: orders, Internal: reader, Inputs: &stubInputStore{}})
	_, err := service.ImportMarginInputs(context.Background(), ImportMarginInputsInput{InstallationID: "inst-1", Limit: 201})
	if err == nil || err.Error() != "limit_exceeded: maximum limit is 200" {
		t.Fatalf("error = %v, want limit_exceeded", err)
	}
	if orders.calls != 0 || len(reader.costIDs) != 0 || len(reader.taxIDs) != 0 {
		t.Fatalf("calls after rejected limit: orders=%d cost=%d tax=%d, want zero", orders.calls, len(reader.costIDs), len(reader.taxIDs))
	}
}

func hasFlag(flags []profitabilitydomain.ProfitSnapshotFlag, want profitabilitydomain.ProfitSnapshotFlag) bool {
	for _, flag := range flags {
		if flag == want {
			return true
		}
	}
	return false
}
