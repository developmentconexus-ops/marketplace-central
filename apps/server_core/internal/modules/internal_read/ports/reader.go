package ports

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type FindProductsInput struct {
	EAN       *string
	SellerSKU *string
	Title     *string
}

type SellableStockInput struct {
	Codprod   int
	Companies []int
	Locations []int
	Now       time.Time
}

type CurrentPriceInput struct {
	Codprod  int
	Codtab   int
	Codlocal int
	Now      time.Time
}

type CostAsOfInput struct {
	Codprod  int
	Codemp   int
	SaleDate string
}

type SalesHistoryInput struct {
	Codprod   *int
	StartDate string
	EndDate   string
}

type TaxInput struct {
	Codprod  int
	SaleDate string
}

type Reader interface {
	FindProductsForLinking(ctx context.Context, input FindProductsInput) ([]domain.ProductCandidate, error)
	GetSellableStock(ctx context.Context, input SellableStockInput) (domain.SellableStock, error)
	GetCurrentPrice(ctx context.Context, input CurrentPriceInput) (domain.CurrentPrice, error)
	GetCostAsOf(ctx context.Context, input CostAsOfInput) (domain.CostAsOf, error)
	GetSalesHistory(ctx context.Context, input SalesHistoryInput) (domain.SalesHistory, error)
	GetTaxInputs(ctx context.Context, input TaxInput) (domain.TaxInputs, error)
}
