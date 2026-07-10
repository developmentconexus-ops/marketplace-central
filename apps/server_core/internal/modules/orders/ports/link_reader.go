package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type LinkReader interface {
	ResolveLinks(ctx context.Context, installationID string, identities []domain.ListingIdentity) (map[string]domain.ItemLink, error)
}
