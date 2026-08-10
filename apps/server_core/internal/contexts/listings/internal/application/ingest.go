package application

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
)

// Service is the one use case this leg ships: fold a channel observation in.
type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

// Current exposes what the store knows about a listing, for callers that need
// to read without folding (the module's public CurrentState wraps this).
func (s *Service) Current(ctx context.Context, key contracts.SourceListingKey) (CurrentListing, bool, error) {
	return s.store.Current(ctx, key)
}

// Ingest decides created/changed/idempotent by comparing the facts this
// observation would store against the facts already stored — not by the
// provider's payload hash. The row we keep is a function of TWO inputs, the
// payload AND the mapping that turns it into facts, and a fold keyed on the
// hash alone only sees the first one move. That is what let a mapper defect
// hide: fixing the mapper changes the facts a listing maps to without
// changing the channel's bytes, so a hash-keyed fold reports idempotent and
// the corrected facts are never written — exactly what happened when a live
// re-run healed only 14 of 34 rows, the other 20 having an unchanged payload.
// Comparing State instead still keeps the property that made the hash
// attractive in the first place: re-polling an unchanged channel must be
// free, and it still is, because an unchanged payload through an unchanged
// mapper produces identical facts, and ListingState.Equal says so.
func (s *Service) Ingest(ctx context.Context, o contracts.ListingObservation) (contracts.IngestResult, error) {
	if err := o.Validate(); err != nil {
		return contracts.IngestResult{}, err
	}
	current, exists, err := s.store.Current(ctx, o.Key)
	if err != nil {
		return contracts.IngestResult{}, fmt.Errorf("listings: read current: %w", err)
	}
	switch {
	case !exists:
		if err := s.store.SaveVersion(ctx, o, 1); err != nil {
			return contracts.IngestResult{}, fmt.Errorf("listings: save v1: %w", err)
		}
		return contracts.IngestResult{Disposition: contracts.DispositionCreated, Version: 1}, nil
	case current.State.Equal(o.State):
		return contracts.IngestResult{Disposition: contracts.DispositionIdempotent, Version: current.Version}, nil
	default:
		next := current.Version + 1
		if err := s.store.SaveVersion(ctx, o, next); err != nil {
			return contracts.IngestResult{}, fmt.Errorf("listings: save v%d: %w", next, err)
		}
		return contracts.IngestResult{Disposition: contracts.DispositionChanged, Version: next}, nil
	}
}
