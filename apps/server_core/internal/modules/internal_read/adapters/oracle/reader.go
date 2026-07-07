package oracle

import (
	"context"
	"errors"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

var errReadUnavailable = errors.New("SANKHYA_READ_UNAVAILABLE")

type Reader struct {
	cfg Config
}

func NewReader(cfg Config) *Reader {
	return &Reader{cfg: cfg}
}

func (r *Reader) FindProductsForLinking(context.Context, ports.FindProductsInput) ([]domain.ProductCandidate, error) {
	return nil, errReadUnavailable
}

func (r *Reader) GetSellableStock(context.Context, ports.SellableStockInput) (domain.SellableStock, error) {
	return domain.SellableStock{}, errReadUnavailable
}

func (r *Reader) GetCurrentPrice(context.Context, ports.CurrentPriceInput) (domain.CurrentPrice, error) {
	return domain.CurrentPrice{}, errReadUnavailable
}

func (r *Reader) GetCostAsOf(context.Context, ports.CostAsOfInput) (domain.CostAsOf, error) {
	return domain.CostAsOf{}, errReadUnavailable
}

func (r *Reader) GetSalesHistory(context.Context, ports.SalesHistoryInput) (domain.SalesHistory, error) {
	return domain.SalesHistory{}, errReadUnavailable
}

func (r *Reader) GetTaxInputs(context.Context, ports.TaxInput) (domain.TaxInputs, error) {
	return domain.TaxInputs{}, errReadUnavailable
}
