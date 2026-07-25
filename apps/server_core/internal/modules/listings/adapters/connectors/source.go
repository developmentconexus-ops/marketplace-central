package connectors

import (
	"context"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
)

// SnapshotObserver receives every provider snapshot page this source reads,
// BEFORE it is mapped down to canonical listing rows. The canonical row keeps
// only what the listings table stores, so provider identity fields the linker
// needs (EAN, seller SKU) are dropped at that mapping; an observer sees them.
// Implemented by the product_links import service — one provider fetch feeds
// both the listings catalog and the link matcher.
type SnapshotObserver interface {
	AbsorbProviderSnapshots(ctx context.Context, installationID string, snapshots []connectorsdomain.ListingSnapshot) error
}

type Source struct {
	capabilities *connectorsapp.MarketplaceCapabilityService
	observer     SnapshotObserver
}

func NewSource(capabilities *connectorsapp.MarketplaceCapabilityService) Source {
	return Source{capabilities: capabilities}
}

// NewSourceWithObserver is NewSource plus the snapshot observer seam. A nil
// observer is the plain source.
func NewSourceWithObserver(capabilities *connectorsapp.MarketplaceCapabilityService, observer SnapshotObserver) Source {
	return Source{capabilities: capabilities, observer: observer}
}

// observe hands the page to the observer. Its failure fails the read: a pull
// that silently skipped the linker would report a complete catalog while the
// matcher stayed blind to part of it — exactly the state this seam exists to
// end. Callers retry the pull; the upsert is idempotent.
func (s Source) observe(ctx context.Context, account listingsports.InstallationAccount, snapshots []connectorsdomain.ListingSnapshot) error {
	if s.observer == nil || len(snapshots) == 0 {
		return nil
	}
	return s.observer.AbsorbProviderSnapshots(ctx, account.InstallationID, snapshots)
}

// scanPageReader is the optional deep-paging capability: providers whose
// offset paging is depth-capped (ML rejects offset+limit past 1000) expose a
// cursor-based scan; the returned cursor feeds the next ReadPage call.
type scanPageReader interface {
	ListListingsScanPage(ctx context.Context, input connectorsdomain.ListListingsInput) ([]connectorsdomain.ListingSnapshot, string, error)
}

func (s Source) ReadPage(ctx context.Context, account listingsports.InstallationAccount, cursor string, limit int) (listingsports.ListingPage, error) {
	reader, err := s.capabilities.ListingReader(account.ProviderCode)
	if err != nil {
		return listingsports.ListingPage{}, err
	}
	ref := connectorsdomain.ProviderAccountRef{TenantID: account.TenantID, InstallationID: account.InstallationID, ProviderCode: account.ProviderCode, ProviderAccountID: account.ProviderAccountID}
	input := connectorsdomain.ListListingsInput{AccountRef: ref, Cursor: cursor, Limit: limit}

	if scanner, ok := reader.(scanPageReader); ok {
		snapshots, next, err := scanner.ListListingsScanPage(ctx, input)
		if err != nil {
			return listingsports.ListingPage{}, err
		}
		if err := s.observe(ctx, account, snapshots); err != nil {
			return listingsports.ListingPage{}, err
		}
		rows, err := MapListingSnapshotToCanonicalRows(ref, snapshots)
		if err != nil {
			return listingsports.ListingPage{}, err
		}
		return listingsports.ListingPage{ProviderItemCount: len(snapshots), Rows: rows, NextCursor: next}, nil
	}

	snapshots, err := reader.ListListings(ctx, input)
	if err != nil {
		return listingsports.ListingPage{}, err
	}
	if err := s.observe(ctx, account, snapshots); err != nil {
		return listingsports.ListingPage{}, err
	}
	rows, err := MapListingSnapshotToCanonicalRows(ref, snapshots)
	if err != nil {
		return listingsports.ListingPage{}, err
	}
	return listingsports.ListingPage{ProviderItemCount: len(snapshots), Rows: rows}, nil
}
