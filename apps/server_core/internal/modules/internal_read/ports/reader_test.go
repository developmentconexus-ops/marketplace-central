package ports

import (
	"context"
	"reflect"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type stubReader struct{}

func TestReaderContractCompiles(t *testing.T) {
	var _ Reader = (*stubReader)(nil)
}

func TestCurrentPriceInputAllowsImplicitDefaults(t *testing.T) {
	input := CurrentPriceInput{Codprod: 42664}

	if input.Codtab != nil {
		t.Fatalf("expected CODTAB to be optional, got %v", *input.Codtab)
	}
	if input.Codlocal != nil {
		t.Fatalf("expected CODLOCAL to be optional, got %v", *input.Codlocal)
	}

	typ := reflect.TypeOf(CurrentPriceInput{})
	if field, ok := typ.FieldByName("Codtab"); !ok || field.Type.Kind() != reflect.Ptr {
		t.Fatalf("expected Codtab to be a pointer field, got %v", field.Type)
	}
	if field, ok := typ.FieldByName("Codlocal"); !ok || field.Type.Kind() != reflect.Ptr {
		t.Fatalf("expected Codlocal to be a pointer field, got %v", field.Type)
	}
}

func TestSalesHistoryInputSupportsProductOrGroup(t *testing.T) {
	productOnly := SalesHistoryInput{Codprod: intPtr(42664), StartDate: "2026-07-01", EndDate: "2026-07-31"}
	groupOnly := SalesHistoryInput{CodgrupoProd: intPtr(12), StartDate: "2026-07-01", EndDate: "2026-07-31"}

	if productOnly.Codprod == nil || productOnly.CodgrupoProd != nil {
		t.Fatalf("expected product-only history input to keep group nil, got %+v", productOnly)
	}
	if groupOnly.CodgrupoProd == nil || groupOnly.Codprod != nil {
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
