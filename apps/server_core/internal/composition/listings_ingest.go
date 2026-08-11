package composition

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/listings"
	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/port"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// ListingsIngestReport counts every disposition separately: "nothing changed"
// and "nothing was read" are different outcomes (catalog_ingest.go:13-15).
type ListingsIngestReport struct {
	Pages      int
	Observed   int
	Created    int
	Changed    int
	Idempotent int
}

// RunListingsIngest walks a listing feed to exhaustion and folds every
// observation into Listings. Production path: the root owns it, the context
// decides, the adapter speaks wire.
func RunListingsIngest(ctx context.Context, module *listings.Module, feed port.ListingFeed, t tenant.ID, pageSize int) (ListingsIngestReport, error) {
	if pageSize <= 0 {
		return ListingsIngestReport{}, fmt.Errorf("composition: page size must be positive, got %d", pageSize)
	}
	var report ListingsIngestReport
	cursor := port.Cursor{}
	for {
		page, err := feed.NextPage(ctx, t, cursor, pageSize)
		if err != nil {
			return report, fmt.Errorf("composition: read listings feed page %d: %w", report.Pages+1, err)
		}
		report.Pages++
		for _, obs := range page.Observations {
			result, err := module.IngestListing(ctx, obs)
			if err != nil {
				return report, fmt.Errorf("composition: ingest %s: %w", obs.Key.ListingID(), err)
			}
			report.Observed++
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
			return report, fmt.Errorf("composition: feed reported more pages but did not advance the cursor")
		}
		cursor = page.Next
	}
}
