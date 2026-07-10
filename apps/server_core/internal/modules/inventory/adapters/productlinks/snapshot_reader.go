package productlinks

import (
	"context"
	"time"

	inventorydomain "marketplace-central/apps/server_core/internal/modules/inventory/domain"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
	productlinksports "marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

type SnapshotReader struct {
	reader productlinksports.ListingSnapshotReader
}

func NewSnapshotReader(reader productlinksports.ListingSnapshotReader) SnapshotReader {
	return SnapshotReader{reader: reader}
}

func (r SnapshotReader) ListListingSnapshots(ctx context.Context, installationID string, limit int) ([]inventorydomain.ListingSnapshot, error) {
	items, err := r.reader.ListListingSnapshots(ctx, installationID, limit)
	if err != nil {
		return nil, err
	}
	result := make([]inventorydomain.ListingSnapshot, 0, len(items))
	for _, item := range items {
		result = append(result, mapListingSnapshot(item))
	}
	return result, nil
}

type SnapshotStore struct {
	store productlinksports.ListingSnapshotStore
}

func NewSnapshotStore(store productlinksports.ListingSnapshotStore) SnapshotStore {
	return SnapshotStore{store: store}
}

func (s SnapshotStore) SaveListingSnapshot(ctx context.Context, snapshot inventorydomain.ListingSnapshot) error {
	now := time.Now().UTC()
	observedAt := snapshot.ObservedAt
	if observedAt == nil {
		observedAt = &now
	}
	return s.store.UpsertListingSnapshots(ctx, []productlinksdomain.ListingSnapshot{{
		InstallationID:      snapshot.Identity.InstallationID,
		ProviderCode:        snapshot.ProviderCode,
		ProviderItemID:      snapshot.Identity.ProviderItemID,
		ProviderVariationID: snapshot.Identity.ProviderVariationID,
		ProviderStatus:      snapshot.ProviderStatus,
		SellerSKU:           snapshot.SellerSKU,
		EAN:                 snapshot.EAN,
		Title:               snapshot.Title,
		AvailableQuantity:   snapshot.AvailableQuantity,
		SourceUpdatedAt:     observedAt,
		FetchedAt:           now,
		CreatedAt:           now,
		UpdatedAt:           now,
	}})
}

func mapListingSnapshot(item productlinksdomain.ListingSnapshot) inventorydomain.ListingSnapshot {
	observedAt := item.SourceUpdatedAt
	if observedAt == nil {
		observedAt = &item.FetchedAt
	}
	return inventorydomain.ListingSnapshot{
		Identity: inventorydomain.ListingIdentity{
			InstallationID:      item.InstallationID,
			ProviderItemID:      item.ProviderItemID,
			ProviderVariationID: item.ProviderVariationID,
		},
		ProviderCode:      item.ProviderCode,
		ProviderStatus:    item.ProviderStatus,
		SellerSKU:         item.SellerSKU,
		EAN:               item.EAN,
		Title:             item.Title,
		AvailableQuantity: item.AvailableQuantity,
		ObservedAt:        observedAt,
	}
}
