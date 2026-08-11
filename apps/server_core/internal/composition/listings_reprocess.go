package composition

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/port"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// StoredListings is the part of the listings module a reprocess uses: read
// back what was stored, fold what it maps to. *listings.Module satisfies it.
//
// Declared here, where it is consumed, and narrower than the module on
// purpose. A reprocess is the operation that must not be able to call Mercado
// Livre, and this interface is the shape of that promise: nothing reachable
// through it fetches anything.
type StoredListings interface {
	StoredObservations(ctx context.Context, t tenant.ID, after contracts.StoredCursor, limit int) (contracts.StoredPage, error)
	IngestListing(ctx context.Context, o contracts.ListingObservation) (contracts.IngestResult, error)
}

// ListingsReprocessReport counts each disposition separately, like
// ListingsIngestReport. Created is counted even though a reprocess maps rows
// that already exist: if one ever appears it means the listing was deleted
// between the read and the fold, and a silent zero would hide that.
type ListingsReprocessReport struct {
	Pages      int
	Read       int
	Created    int
	Changed    int
	Idempotent int
}

// RunListingsReprocess re-derives every stored listing's facts from the bytes
// already recorded and folds them back in, without calling the channel.
//
// It exists because the stored row is a function of two inputs — the payload
// AND the mapping that turns it into facts — and only one of them arrives from
// the channel. When a mapper is corrected, a live poll heals only the listings
// the channel still serves and only for as long as it keeps serving them; a
// listing deleted at the source keeps the old mapping's answer forever, and
// nothing about that is visible in a green run. This walks what we hold.
//
// The decision stays in Ingest: this loop reads, maps, and hands over. That is
// what makes a reprocess and a live poll converge instead of drifting, and it
// is why an unchanged mapping over unchanged bytes reports idempotent here for
// exactly the same reason it does there.
func RunListingsReprocess(ctx context.Context, module StoredListings, mapper port.ListingMapper, t tenant.ID, pageSize int) (ListingsReprocessReport, error) {
	if pageSize <= 0 {
		return ListingsReprocessReport{}, fmt.Errorf("composition: page size must be positive, got %d", pageSize)
	}
	var report ListingsReprocessReport
	cursor := contracts.StoredCursor{}
	for {
		page, err := module.StoredObservations(ctx, t, cursor, pageSize)
		if err != nil {
			return report, fmt.Errorf("composition: read stored listings page %d: %w", report.Pages+1, err)
		}
		report.Pages++
		for _, stored := range page.Observations {
			observation, err := mapper.MapStored(stored)
			if err != nil {
				return report, fmt.Errorf("composition: map stored %s: %w", stored.Key.ListingID(), err)
			}
			result, err := module.IngestListing(ctx, observation)
			if err != nil {
				return report, fmt.Errorf("composition: ingest %s: %w", stored.Key.ListingID(), err)
			}
			report.Read++
			switch result.Disposition {
			case contracts.DispositionCreated:
				report.Created++
			case contracts.DispositionChanged:
				report.Changed++
			case contracts.DispositionIdempotent:
				report.Idempotent++
			}
		}
		if page.Done {
			return report, nil
		}
		if page.Next.IsStart() {
			return report, fmt.Errorf("composition: stored listings page %d reported more pages but did not advance the cursor", report.Pages)
		}
		cursor = page.Next
	}
}
