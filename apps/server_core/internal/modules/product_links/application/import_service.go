package application

import (
	"context"
	"errors"
	"strings"
	"time"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/product_links/domain"
	"marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

type ListingSource interface {
	ListListings(ctx context.Context, installationID string, limit int) ([]connectorsdomain.ListingSnapshot, error)
}

type ImportService struct {
	source ListingSource
	store  ports.ListingSnapshotStore
	now    func() time.Time
}

type ImportServiceConfig struct {
	Source ListingSource
	Store  ports.ListingSnapshotStore
	Now    func() time.Time
}

type ImportListingSnapshotsInput struct {
	InstallationID string
	Limit          int
}

func NewImportService(cfg ImportServiceConfig) *ImportService {
	now := cfg.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &ImportService{
		source: cfg.Source,
		store:  cfg.Store,
		now:    now,
	}
}

func (s *ImportService) ImportListingSnapshots(ctx context.Context, input ImportListingSnapshotsInput) (domain.ListingSnapshotImportResult, error) {
	if s.source == nil || s.store == nil {
		return domain.ListingSnapshotImportResult{}, errors.New("PRODUCT_LINKS_IMPORT_NOT_CONFIGURED")
	}

	installationID := strings.TrimSpace(input.InstallationID)
	if installationID == "" {
		return domain.ListingSnapshotImportResult{}, errors.New("PRODUCT_LINKS_INSTALLATION_REQUIRED")
	}

	listings, err := s.source.ListListings(ctx, installationID, input.Limit)
	if err != nil {
		return domain.ListingSnapshotImportResult{}, err
	}

	snapshots := normalizeListingSnapshots(installationID, listings, s.now())
	if err := s.store.UpsertListingSnapshots(ctx, snapshots); err != nil {
		return domain.ListingSnapshotImportResult{}, err
	}

	return domain.ListingSnapshotImportResult{
		InstallationID: installationID,
		ImportedCount:  len(snapshots),
		Items:          snapshots,
	}, nil
}

// AbsorbProviderSnapshots upserts snapshots the CALLER already fetched from the
// provider, instead of fetching them again. The listings sync already walks the
// seller's whole catalog (ML scan paging, no depth cap) and holds the very same
// ListingSnapshot values — including EAN and SellerSKU, which the canonical
// listings table has no column for and therefore discarded. Without this seam
// the only way a listing ever got a product_link snapshot was the manual POST
// path below, whose HTTP default of limit=20 left the linker seeing a fraction
// of the catalog: the operator then saw thousands of anúncios "sem vínculo"
// that had simply never been offered to the matcher. One provider fetch, two
// consumers; the linker's coverage is now the sync's coverage.
func (s *ImportService) AbsorbProviderSnapshots(ctx context.Context, installationID string, listings []connectorsdomain.ListingSnapshot) error {
	if s.store == nil {
		return errors.New("PRODUCT_LINKS_IMPORT_NOT_CONFIGURED")
	}
	installationID = strings.TrimSpace(installationID)
	if installationID == "" {
		return errors.New("PRODUCT_LINKS_INSTALLATION_REQUIRED")
	}
	if len(listings) == 0 {
		return nil
	}
	return s.store.UpsertListingSnapshots(ctx, normalizeListingSnapshots(installationID, listings, s.now()))
}

func normalizeListingSnapshots(installationID string, listings []connectorsdomain.ListingSnapshot, now time.Time) []domain.ListingSnapshot {
	snapshots := make([]domain.ListingSnapshot, 0, len(listings))
	for _, listing := range listings {
		if len(listing.Variations) == 0 {
			snapshots = append(snapshots, domain.ListingSnapshot{
				InstallationID:    installationID,
				ProviderCode:      listing.ProviderCode,
				ProviderItemID:    strings.TrimSpace(listing.ProviderItemID),
				ProviderStatus:    strings.TrimSpace(listing.ProviderStatus),
				SellerSKU:         strings.TrimSpace(listing.SellerSKU),
				EAN:               strings.TrimSpace(listing.EAN),
				Title:             strings.TrimSpace(listing.Title),
				AvailableQuantity: listing.AvailableQuantity,
				SourceUpdatedAt:   listing.SourceUpdatedAt,
				FetchedAt:         listing.FetchedAt.UTC(),
				CreatedAt:         now,
				UpdatedAt:         now,
			})
			continue
		}

		for _, variation := range listing.Variations {
			snapshots = append(snapshots, domain.ListingSnapshot{
				InstallationID:      installationID,
				ProviderCode:        listing.ProviderCode,
				ProviderItemID:      strings.TrimSpace(listing.ProviderItemID),
				ProviderVariationID: strings.TrimSpace(variation.ProviderVariationID),
				ProviderStatus:      strings.TrimSpace(listing.ProviderStatus),
				SellerSKU:           strings.TrimSpace(variation.SellerSKU),
				EAN:                 strings.TrimSpace(variation.EAN),
				Title:               strings.TrimSpace(listing.Title),
				AvailableQuantity:   variation.AvailableQuantity,
				SourceUpdatedAt:     listing.SourceUpdatedAt,
				FetchedAt:           listing.FetchedAt.UTC(),
				CreatedAt:           now,
				UpdatedAt:           now,
			})
		}
	}
	return snapshots
}
