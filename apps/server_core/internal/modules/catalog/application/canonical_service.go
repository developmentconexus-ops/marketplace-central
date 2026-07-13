package application

import (
	"context"
	"marketplace-central/apps/server_core/internal/modules/catalog/domain"
	"marketplace-central/apps/server_core/internal/modules/catalog/ports"
)

type CanonicalService struct{ reader ports.CanonicalProductReader }

func NewCanonicalService(reader ports.CanonicalProductReader) CanonicalService {
	return CanonicalService{reader: reader}
}
func (s CanonicalService) ListProducts(ctx context.Context) ([]domain.CanonicalProduct, error) {
	return s.reader.ListCanonicalProducts(ctx)
}
func (s CanonicalService) GetProduct(ctx context.Context, id domain.InternalProductID) (domain.CanonicalProduct, error) {
	return s.reader.GetCanonicalProduct(ctx, id)
}
func (s CanonicalService) SearchProducts(ctx context.Context, q string) ([]domain.CanonicalProduct, error) {
	return s.reader.SearchCanonicalProducts(ctx, q)
}
