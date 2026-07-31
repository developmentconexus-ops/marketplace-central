package application

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

type Service struct {
	reader ports.Source
}

func NewService(reader ports.Source) Service {
	return Service{reader: reader}
}

func (s Service) FindProductsForLinking(ctx context.Context, input ports.FindProductsInput) ([]domain.ProductCandidate, error) {
	return s.reader.FindProductsForLinking(ctx, input)
}

func (s Service) GetSellableStock(ctx context.Context, input ports.SellableStockInput) (domain.SellableStock, error) {
	return s.reader.GetSellableStock(ctx, input)
}

func (s Service) GetCurrentPrice(ctx context.Context, input ports.CurrentPriceInput) (domain.CurrentPrice, error) {
	return s.reader.GetCurrentPrice(ctx, input)
}

func (s Service) GetCostAsOf(ctx context.Context, input ports.CostAsOfInput) (domain.CostAsOf, error) {
	return s.reader.GetCostAsOf(ctx, input)
}

func (s Service) GetSalesHistory(ctx context.Context, input ports.SalesHistoryInput) (domain.SalesHistory, error) {
	return s.reader.GetSalesHistory(ctx, input)
}

func (s Service) GetTaxInputs(ctx context.Context, input ports.TaxInput) (domain.TaxInputs, error) {
	return s.reader.GetTaxInputs(ctx, input)
}

func (s Service) ListCatalogProductFacts(ctx context.Context, cursor ports.Cursor, limit int, policy *ports.SellableAssortmentPolicy) (ports.CatalogFactPage, error) {
	return s.reader.ListCatalogProductFacts(ctx, cursor, limit, policy)
}

func (s Service) SearchCatalogProductFacts(ctx context.Context, q string, cursor ports.Cursor, limit int, policy *ports.SellableAssortmentPolicy) (ports.CatalogFactPage, error) {
	return s.reader.SearchCatalogProductFacts(ctx, q, cursor, limit, policy)
}

func (s Service) GetCatalogAssortmentCounts(ctx context.Context, policy *ports.SellableAssortmentPolicy) (ports.CatalogAssortmentCounts, error) {
	return s.reader.GetCatalogAssortmentCounts(ctx, policy)
}

func (s Service) CatalogProductFactsByIDs(ctx context.Context, ids []int64) (ports.CatalogFactPage, error) {
	return s.reader.CatalogProductFactsByIDs(ctx, ids)
}

var _ ports.CatalogPageReader = Service{}
