package ports

import (
	"context"
	"marketplace-central/apps/server_core/internal/modules/catalog/domain"
)

// CanonicalProductReader is the Oracle/internal_read catalog boundary.
type CanonicalProductReader interface {
	ListCanonicalProducts(context.Context) ([]domain.CanonicalProduct, error)
	GetCanonicalProduct(context.Context, domain.InternalProductID) (domain.CanonicalProduct, error)
	SearchCanonicalProducts(context.Context, string) ([]domain.CanonicalProduct, error)
}
