package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// CatalogMatchReader resolves a seller product against the provider catalog,
// buy-box, and category predictions. Read-only (GET); no provider writes.
type CatalogMatchReader interface {
	ReadCatalogMatch(ctx context.Context, input domain.CatalogMatchInput) (domain.CatalogMatchSnapshot, error)
}
