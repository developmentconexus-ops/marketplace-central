package connectors

import (
	"context"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
)

type Source struct {
	capabilities *connectorsapp.MarketplaceCapabilityService
}

func NewSource(capabilities *connectorsapp.MarketplaceCapabilityService) Source {
	return Source{capabilities: capabilities}
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
	rows, err := MapListingSnapshotToCanonicalRows(ref, snapshots)
	if err != nil {
		return listingsports.ListingPage{}, err
	}
	return listingsports.ListingPage{ProviderItemCount: len(snapshots), Rows: rows}, nil
}
