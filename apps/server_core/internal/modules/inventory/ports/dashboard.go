package ports

import (
	"context"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
)

type StockFact = internalreaddomain.StockFact

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

type InternalStockBatchReader interface {
	GetStockFactsByIDs(ctx context.Context, ids []int64) (map[int64]*StockFact, error)
}

type ListingSnapshotStore interface {
	SaveListingSnapshot(ctx context.Context, snapshot domain.ListingSnapshot) error
}
