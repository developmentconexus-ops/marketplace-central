package composition

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/catalog"
	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/port"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// IngestReport is what one full walk of a feed did. Every disposition is counted
// separately because "nothing changed" and "nothing was read" are different
// outcomes and a single total cannot tell them apart.
type IngestReport struct {
	Pages      int
	Observed   int
	Created    int
	Changed    int
	Idempotent int
	Conflicts  []contracts.Identifier
}

// RunCatalogIngest walks a feed to exhaustion and folds every observation into
// Catalog. It is the production path: the composition root owns it, the context
// owns the decision, the adapter owns the SQL.
func RunCatalogIngest(ctx context.Context, module *catalog.Module, feed port.ProductFeed, t tenant.ID, pageSize int) (IngestReport, error) {
	if pageSize <= 0 {
		return IngestReport{}, fmt.Errorf("composition: page size must be positive, got %d", pageSize)
	}
	var report IngestReport
	cursor := port.Cursor{}
	for {
		page, err := feed.NextPage(ctx, t, cursor, pageSize)
		if err != nil {
			return report, fmt.Errorf("composition: read catalog feed page %d: %w", report.Pages+1, err)
		}
		report.Pages++
		for _, obs := range page.Observations {
			result, err := module.IngestProduct(ctx, obs)
			if err != nil {
				return report, fmt.Errorf("composition: ingest %s: %w", obs.Key, err)
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
			report.Conflicts = append(report.Conflicts, result.DuplicateIdentifiers...)
		}
		if page.Done {
			return report, nil
		}
		if page.Next.IsStart() {
			// A source that is not done must advance. Returning the start cursor
			// again would loop forever, and a job that never ends looks exactly
			// like a job that is working.
			return report, fmt.Errorf("composition: feed reported more pages but did not advance the cursor")
		}
		cursor = page.Next
	}
}
