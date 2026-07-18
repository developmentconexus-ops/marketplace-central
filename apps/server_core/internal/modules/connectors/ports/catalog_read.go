package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

type CatalogOffersReader interface {
	ListCatalogOffers(ctx context.Context, accountRef domain.ProviderAccountRef, catalogProductID string) ([]domain.CatalogOffer, error)
}

type CatalogReader interface {
	SearchCatalogByEAN(ctx context.Context, accountRef domain.ProviderAccountRef, ean string) (domain.CatalogSearchResult, error)
	GetCatalogProduct(ctx context.Context, accountRef domain.ProviderAccountRef, catalogProductID string) (domain.CatalogProduct, error)
}
