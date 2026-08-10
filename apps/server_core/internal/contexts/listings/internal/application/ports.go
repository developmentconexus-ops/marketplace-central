// Package application owns the ingest decision. Its Store port is internal:
// only the postgres adapter inside this context implements it.
package application

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
)

// CurrentListing is what the store knows about a listing before an ingest.
type CurrentListing struct {
	Version     int
	PayloadHash string
}

// Store persists listing versions. SaveVersion must be atomic: the listing
// row, its variations and the observation land together or not at all.
//
// KNOWN GAP — the pair is not atomic, only each half is. Ingest reads Current
// and then calls SaveVersion in a second transaction, so two ingests of the
// same listing running at once can both read version N and both write N+1: the
// later write wins the row and one version silently ceases to exist. Nothing
// today can reach that state — cmd/listingsingest is a manual, single-process
// command and there is no scheduler — which is why this is recorded rather
// than designed around. It becomes real the moment a second writer exists
// (a scheduler, a webhook consumer, or two operators running the command at
// once), and the fix is a decision about this port, not a patch inside the
// adapter: either SaveVersion takes the expected current version and fails on
// mismatch, or the two calls collapse into one store operation that reads and
// writes under the same transaction.
type Store interface {
	Current(ctx context.Context, key contracts.SourceListingKey) (CurrentListing, bool, error)
	SaveVersion(ctx context.Context, o contracts.ListingObservation, version int) error
}
