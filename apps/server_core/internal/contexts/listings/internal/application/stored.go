package application

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// StoredObservations reads one page of the payloads Listings has recorded, so
// a caller can hand them back to the adapter that produced them and fold the
// re-derived facts in through Ingest.
//
// This is a read, not a second write path: nothing here decides anything about
// a listing. The decision stays in Ingest, which is why a reprocess and a live
// poll converge on the same rule instead of two.
func (s *Service) StoredObservations(ctx context.Context, t tenant.ID, after contracts.StoredCursor, limit int) (contracts.StoredPage, error) {
	if limit <= 0 {
		return contracts.StoredPage{}, fmt.Errorf("listings: stored observations limit must be positive, got %d", limit)
	}
	rows, err := s.store.StoredObservations(ctx, t, after, limit)
	if err != nil {
		return contracts.StoredPage{}, fmt.Errorf("listings: read stored observations: %w", err)
	}
	if len(rows) == 0 {
		return contracts.StoredPage{Done: true}, nil
	}
	// A short page ends the walk; a full one cannot, because "exactly limit
	// rows left" and "more rows left" are indistinguishable from here. The
	// cost of being right is one extra empty page per walk.
	return contracts.StoredPage{
		Observations: rows,
		Next:         contracts.StoredCursorAfter(rows[len(rows)-1].Key),
		Done:         len(rows) < limit,
	}, nil
}
