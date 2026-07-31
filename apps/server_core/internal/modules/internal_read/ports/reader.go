package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type FindProductsInput struct {
	ProductID       *int
	EAN             *string
	SellerSKU       *string
	Title           *string
	IncludeInactive bool
}

type SellableStockInput struct {
	ProductID int
	Policy    domain.SellableStockPolicy
	Freshness domain.FreshnessPolicy
}

type CurrentPriceInput struct {
	ProductID int
	Policy    domain.CurrentPricePolicy
	Freshness domain.FreshnessPolicy
}

type CostAsOfInput struct {
	ProductID int
	Policy    domain.CostAsOfPolicy
	Freshness domain.FreshnessPolicy
}

type SalesHistoryInput struct {
	ProductID      *int
	ProductGroupID *int
	Window         domain.SalesHistoryWindow
	Freshness      domain.FreshnessPolicy
}

type TaxInput struct {
	ProductID int
	Policy    domain.TaxPolicy
	Freshness domain.FreshnessPolicy
}

type Reader interface {
	FindProductsForLinking(ctx context.Context, input FindProductsInput) ([]domain.ProductCandidate, error)
	GetSellableStock(ctx context.Context, input SellableStockInput) (domain.SellableStock, error)
	GetCurrentPrice(ctx context.Context, input CurrentPriceInput) (domain.CurrentPrice, error)
	GetCostAsOf(ctx context.Context, input CostAsOfInput) (domain.CostAsOf, error)
	GetSalesHistory(ctx context.Context, input SalesHistoryInput) (domain.SalesHistory, error)
	GetTaxInputs(ctx context.Context, input TaxInput) (domain.TaxInputs, error)
}

// Source is a complete internal_read source: it answers the per-product reads
// AND pages the catalog. Every decorator in the live chain — cache, timing,
// application.Service — wraps one, so each takes this single parameter type
// rather than a Reader it has to assert its way back out of.
//
// That assertion is what this type exists to delete. A decorator declaring a
// plain Reader field could still forward the paging methods, by recovering the
// capability at call time with `next.(CatalogPageReader)`; when the recovery
// failed it either panicked mid-request or reported the source as unavailable,
// so a mis-wiring surfaced as an outage rather than as a build error. Requiring
// both halves in the parameter type moves that whole class to compile time.
type Source interface {
	Reader
	CatalogPageReader
}
