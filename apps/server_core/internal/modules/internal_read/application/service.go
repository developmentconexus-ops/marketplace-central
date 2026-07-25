package application

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

type Service struct {
	reader ports.Reader
}

func NewService(reader ports.Reader) Service {
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

func (s Service) ListCatalogProductFacts(ctx context.Context, cursor ports.Cursor, limit int) (ports.CatalogFactPage, error) {
	reader, ok := s.reader.(ports.CatalogPageReader)
	if !ok {
		return ports.CatalogFactPage{}, domain.NewReadError(domain.ReadErrorSourceUnavailable, "oracle catalog page reader is unavailable", nil)
	}
	return reader.ListCatalogProductFacts(ctx, cursor, limit)
}

func (s Service) SearchCatalogProductFacts(ctx context.Context, q string, limit int) (ports.CatalogFactPage, error) {
	reader, ok := s.reader.(ports.CatalogPageReader)
	if !ok {
		return ports.CatalogFactPage{}, domain.NewReadError(domain.ReadErrorSourceUnavailable, "oracle catalog page reader is unavailable", nil)
	}
	return reader.SearchCatalogProductFacts(ctx, q, limit)
}

func (s Service) CatalogProductFactsByIDs(ctx context.Context, ids []int64) (ports.CatalogFactPage, error) {
	reader, ok := s.reader.(ports.CatalogPageReader)
	if !ok {
		return ports.CatalogFactPage{}, domain.NewReadError(domain.ReadErrorSourceUnavailable, "oracle catalog page reader is unavailable", nil)
	}
	return reader.CatalogProductFactsByIDs(ctx, ids)
}

var _ ports.CatalogPageReader = Service{}
