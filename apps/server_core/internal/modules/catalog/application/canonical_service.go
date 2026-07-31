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
func (s CanonicalService) GetProduct(ctx context.Context, id domain.InternalProductID) (domain.CanonicalProduct, error) {
	return s.reader.GetCanonicalProduct(ctx, id)
}
