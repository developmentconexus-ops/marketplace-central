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

func (s Source) ReadPage(ctx context.Context, account listingsports.InstallationAccount, cursor string, limit int) (listingsports.ListingPage, error) {
	reader, err := s.capabilities.ListingReader(account.ProviderCode)
	if err != nil {
		return listingsports.ListingPage{}, err
	}
	ref := connectorsdomain.ProviderAccountRef{TenantID: account.TenantID, InstallationID: account.InstallationID, ProviderCode: account.ProviderCode, ProviderAccountID: account.ProviderAccountID}
	snapshots, err := reader.ListListings(ctx, connectorsdomain.ListListingsInput{AccountRef: ref, Cursor: cursor, Limit: limit})
	if err != nil {
		return listingsports.ListingPage{}, err
	}
	rows, err := MapListingSnapshotToCanonicalRows(ref, snapshots)
	if err != nil {
		return listingsports.ListingPage{}, err
	}
	return listingsports.ListingPage{ProviderItemCount: len(snapshots), Rows: rows}, nil
}
