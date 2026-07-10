package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
)

type ListingSnapshotReader interface {
	ListListingSnapshots(ctx context.Context, installationID string, limit int) ([]domain.ListingSnapshot, error)
}

type ListingLinkReader interface {
	ListLinks(ctx context.Context, installationID string, limit int) ([]domain.LinkRecord, error)
}

type InstallationReader interface {
	GetInstallation(ctx context.Context, installationID string) (domain.ProviderStockRef, bool, error)
}

type InternalStockReader interface {
	GetSellableStock(ctx context.Context, productID int, policy domain.StockPolicy) (domain.InternalStockEvidence, domain.ProductEvidence, error)
}

type ListingSnapshotStore interface {
	SaveListingSnapshot(ctx context.Context, snapshot domain.ListingSnapshot) error
}
