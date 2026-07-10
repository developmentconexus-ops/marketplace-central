package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type ListingSnapshotStore interface {
	UpsertListingSnapshots(ctx context.Context, snapshots []domain.ListingSnapshot) error
}
