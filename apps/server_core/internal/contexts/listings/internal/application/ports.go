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
type Store interface {
	Current(ctx context.Context, key contracts.SourceListingKey) (CurrentListing, bool, error)
	SaveVersion(ctx context.Context, o contracts.ListingObservation, version int) error
}
