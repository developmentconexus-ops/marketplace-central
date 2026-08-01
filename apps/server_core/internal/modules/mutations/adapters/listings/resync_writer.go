package listings

import (
	"context"
	"errors"
	"strings"

	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	mutationsdomain "marketplace-central/apps/server_core/internal/modules/mutations/domain"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

// ListingIngestor is the per-item hydrate+persist operation
// (listings/ports.ListingIngestor's method set, mirrored locally per this
// package's existing idiom of package-local interfaces shaped after
// listingsapp's real types). A resync WriteItem names exactly one listing;
// this writer must call this exactly once per Apply, never fan out into a
// whole-catalog pull (listingsapp.BackfillRunner) per queued item.
type ListingIngestor interface {
	IngestListing(ctx context.Context, providerListingID string) error
}

// ResyncWriter applies an admin-triggered listing resync by delegating to a
// shared ListingIngestor, scoped to one installation account, for exactly
// the one listing named by item.ListingID. It no longer owns any
// hydrate/persist logic itself (that lives in the ingestor) — its own job is
// strictly the guard clauses below, the installation-scope check, and
// extracting the provider listing id from the canonical ListingID.
type ResyncWriter struct {
	ingestor ListingIngestor
	account  listingsports.InstallationAccount
}

func NewResyncWriter(ingestor ListingIngestor, account listingsports.InstallationAccount) *ResyncWriter {
	return &ResyncWriter{ingestor: ingestor, account: account}
}

func (w *ResyncWriter) Apply(ctx context.Context, item mutationsports.WriteItem) (mutationsports.WriteOutcome, error) {
	if w == nil || w.ingestor == nil {
		return mutationsports.WriteOutcome{}, errors.New("listing resync writer is not configured")
	}
	if item.ProtocolType != mutationsdomain.ProtocolTypeListingResync {
		return mutationsports.WriteOutcome{}, errors.New("listing resync writer received unsupported protocol type")
	}
	if strings.TrimSpace(w.account.TenantID) == "" || strings.TrimSpace(w.account.InstallationID) == "" || strings.TrimSpace(item.InstallationID) != strings.TrimSpace(w.account.InstallationID) {
		return mutationsports.WriteOutcome{}, errors.New("listing resync installation scope is invalid")
	}
	parsed, err := listingsdomain.ParseListingID(item.ListingID)
	if err != nil {
		return mutationsports.WriteOutcome{}, err
	}
	if err := w.ingestor.IngestListing(ctx, parsed.ProviderListingID); err != nil {
		return mutationsports.WriteOutcome{}, err
	}
	return mutationsports.WriteOutcome{}, nil
}
