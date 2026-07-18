package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

type CatalogOffersReader interface {
	ListCatalogOffers(ctx context.Context, accountRef domain.ProviderAccountRef, catalogProductID string) ([]domain.CatalogOffer, error)
}
