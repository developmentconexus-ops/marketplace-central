package ports

import (
	"context"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type stubReader struct{}

func TestReaderContractCompiles(t *testing.T) {
	var _ Reader = (*stubReader)(nil)
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
