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

// Ingest decides created/changed/idempotent by payload hash. Re-polling an
// unchanged channel must be free (same rule the catalog leg ratified).
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
	case current.PayloadHash == o.Evidence.PayloadHash():
		return contracts.IngestResult{Disposition: contracts.DispositionIdempotent, Version: current.Version}, nil
	default:
		next := current.Version + 1
		if err := s.store.SaveVersion(ctx, o, next); err != nil {
			return contracts.IngestResult{}, fmt.Errorf("listings: save v%d: %w", next, err)
		}
		return contracts.IngestResult{Disposition: contracts.DispositionChanged, Version: next}, nil
	}
}
