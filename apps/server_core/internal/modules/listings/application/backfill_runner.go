package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

// maxBackfillPages is BackfillRunner's honest ceiling on a single synchronous
// Pull call, mirroring the old Ingestion.Pull's maxIngestionPages guard: an
// enumerator that never reports NextCursor=="" (bug or malice) would
// otherwise loop and grow the run without bound. Sized identically to the
// old constant because it bounds the SAME unit — listingsBackfillPageLimit
// (1000 ids) per page — so 10_000 pages is the same ~10M-id catalog ceiling
// as before; a real catalog that legitimately exceeds it must be resynced in
// smaller runs rather than trusted to loop forever.
const maxBackfillPages = 10_000

// BackfillRunner is the whole-catalog, synchronous Pull operation both
// RefreshService (manual/API-triggered refresh) and ResyncWriter
// (mutations/adapters/listings, an admin per-item resync that still walks
// the account's full catalog) need: walk enumerator pages to exhaustion,
// hydrating+persisting each one through the SAME ingestBatch helper
// NewListingsJob's per-tick scheduler path uses (backfill.go), so every
// producer of a full-catalog walk shares identical hydrate+persist behavior
// (IC-06 "MESMO IngestListing").
//
// Unlike NewListingsJob (one enumeration page per scheduler tick, cursor
// persisted between ticks so the process can restart mid-run), Pull walks to
// completion inside a single call — the shape callers that need "block until
// the whole catalog is refreshed" require. This mirrors the OLD Ingestion's
// external contract exactly (same Pull(ctx, account) error signature), but
// NOT its internal atomicity: Ingestion accumulated every row in memory and
// upserted once at the very end (a page-N failure discarded pages 1..N-1
// too), whereas Pull persists page-by-page via ingestBatch, same as
// NewListingsJob. A later page's enumerator/hydrator failure returns the
// error but does NOT undo earlier pages' upserts — this is the F-02/F-03
// CompletedPullStore contract's own behavior (UpsertPulledRows/
// MarkRunComplete split), which Pull simply reuses rather than working
// around; it never calls MarkRunComplete on that path, so no row is marked
// absent by an incomplete run (see backfill.go's ingestBatch doc).
type BackfillRunner struct {
	enumerator ports.IDEnumerator
	hydrator   ports.BatchHydrator
	store      ports.CompletedPullStore
	now        func() time.Time
}

// NewBackfillRunner builds the runner. now may be nil (defaults to time.Now).
func NewBackfillRunner(enumerator ports.IDEnumerator, hydrator ports.BatchHydrator, store ports.CompletedPullStore, now func() time.Time) *BackfillRunner {
	if now == nil {
		now = time.Now
	}
	return &BackfillRunner{enumerator: enumerator, hydrator: hydrator, store: store, now: now}
}

// Pull walks the account's whole catalog to completion: enumerate a page of
// ids, hydrate+persist it, and either advance to the next page or (once the
// enumerator reports exhaustion) call MarkRunComplete and return.
//
// runStarted is captured ONCE at the start of the call and held constant
// across every page — the same invariant NewListingsJob's cursor.
// RunStartedAt enforces across ticks (backfill.go) — so MarkRunComplete's
// keep-absent bound covers every id seen this run, not just the last page.
func (r *BackfillRunner) Pull(ctx context.Context, account ports.InstallationAccount) error {
	if r == nil || r.enumerator == nil || r.hydrator == nil || r.store == nil || r.now == nil {
		return errors.New("listing backfill runner is not configured")
	}

	runStarted := r.now().UTC()

	cursor := ""
	for pages := 0; ; pages++ {
		if pages >= maxBackfillPages {
			return fmt.Errorf("listing backfill exceeded %d pages without a final short page", maxBackfillPages)
		}
		page, err := r.enumerator.ReadIDPage(ctx, account, cursor, listingsBackfillPageLimit)
		if err != nil {
			return fmt.Errorf("listing backfill: enumerate page: %w", err)
		}
		if len(page.IDs) > 0 {
			if _, err := ingestBatch(ctx, r.hydrator, r.store, account, page.IDs, runStarted); err != nil {
				return fmt.Errorf("listing backfill: %w", err)
			}
		}
		if strings.TrimSpace(page.NextCursor) == "" {
			return r.store.MarkRunComplete(ctx, account.InstallationID, runStarted)
		}
		cursor = page.NextCursor
	}
}
