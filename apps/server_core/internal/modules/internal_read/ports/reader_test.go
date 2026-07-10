package ports

import (
	"context"
	"reflect"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type stubReader struct{}

func TestReaderContractCompiles(t *testing.T) {
	var _ Reader = (*stubReader)(nil)
}

func TestCurrentPriceInputUsesTypedPolicy(t *testing.T) {
	now := time.Date(2026, 7, 7, 10, 0, 0, 0, time.UTC)
	input := CurrentPriceInput{
		ProductID: 42664,
		Policy:    domain.DefaultCurrentPricePolicy(now),
	}

	if input.Policy.PriceTableID != domain.DefaultPriceTableID {
		t.Fatalf("expected default price table %d, got %d", domain.DefaultPriceTableID, input.Policy.PriceTableID)
	}
	if input.Policy.LocationID != domain.DefaultPriceLocalID {
		t.Fatalf("expected default price local %d, got %d", domain.DefaultPriceLocalID, input.Policy.LocationID)
	}

	typ := reflect.TypeOf(CurrentPriceInput{})
	if field, ok := typ.FieldByName("Policy"); !ok || field.Type != reflect.TypeOf(domain.CurrentPricePolicy{}) {
		t.Fatalf("expected Policy field of type CurrentPricePolicy, got %v", field.Type)
	}
}

func TestSalesHistoryInputSupportsProductOrGroup(t *testing.T) {
	start := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)
	productOnly := SalesHistoryInput{ProductID: intPtr(42664), Window: domain.SalesHistoryWindow{Start: start, End: end}}
	groupOnly := SalesHistoryInput{ProductGroupID: intPtr(12), Window: domain.SalesHistoryWindow{Start: start, End: end}}

	if productOnly.ProductID == nil || productOnly.ProductGroupID != nil {
		t.Fatalf("expected product-only history input to keep group nil, got %+v", productOnly)
	}
	if groupOnly.ProductGroupID == nil || groupOnly.ProductID != nil {
		t.Fatalf("expected group-only history input to keep product nil, got %+v", groupOnly)
	}
}

func (stubReader) FindProductsForLinking(context.Context, FindProductsInput) ([]domain.ProductCandidate, error) {
	return nil, nil
}

func (stubReader) GetSellableStock(context.Context, SellableStockInput) (domain.SellableStock, error) {
	return domain.SellableStock{}, nil
}

func (stubReader) GetCurrentPrice(context.Context, CurrentPriceInput) (domain.CurrentPrice, error) {
	return domain.CurrentPrice{}, nil
}

func (stubReader) GetCostAsOf(context.Context, CostAsOfInput) (domain.CostAsOf, error) {
	return domain.CostAsOf{}, nil
}

func (stubReader) GetSalesHistory(context.Context, SalesHistoryInput) (domain.SalesHistory, error) {
	return domain.SalesHistory{}, nil
}

func (stubReader) GetTaxInputs(context.Context, TaxInput) (domain.TaxInputs, error) {
	return domain.TaxInputs{}, nil
}

func intPtr(v int) *int {
	return &v
}
