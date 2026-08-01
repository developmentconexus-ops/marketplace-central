package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
)

// listingsBackfillPageLimit is the enumeration page size (IC-06: "scan/
// scroll_id coleta ids (1000/batch)"). Distinct from the multiget hydration
// batch size (20, M-01/F-02) — GetItemsMultiget partitions internally, so a
// single ingestBatch call for a whole enumeration page already fans out into
// the right number of 20-sized HTTP calls without this layer chunking again.
const listingsBackfillPageLimit = 1000

// ingestBatch is the shared plumbing behind ListingIngestor.IngestListing
// (one id) and NewListingsJob (one enumeration page's worth of ids per
// tick): hydrate, then persist via CompletedPullStore.UpsertPulledRows. It
// NEVER calls MarkRunComplete — only the resumable job decides a run is
// complete; a single-item ingest has no notion of "the whole catalog was
// scanned" and must never trigger the keep-absent step.
func ingestBatch(ctx context.Context, hydrator ports.BatchHydrator, store ports.CompletedPullStore, account ports.InstallationAccount, ids []string, seenAt time.Time) ([]domain.Listing, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := hydrator.HydrateBatch(ctx, account, ids)
	if err != nil {
		return nil, fmt.Errorf("hydrate batch: %w", err)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	if err := store.UpsertPulledRows(ctx, account.InstallationID, rows, seenAt); err != nil {
		return nil, fmt.Errorf("upsert batch: %w", err)
	}
	return rows, nil
}

type listingIngestor struct {
	hydrator ports.BatchHydrator
	store    ports.CompletedPullStore
	account  ports.InstallationAccount
	now      func() time.Time
}

// NewListingIngestor builds IC-06's IngestListing operation, bound to one
// installation account (same idiom as mutations/adapters/listings.
// NewResyncWriter). Every producer — refresh, webhook worker, ad-hoc resync
// — that needs a single-item upsert goes through this, sharing ingestBatch
// with the backfill job so both paths hydrate and persist identically.
func NewListingIngestor(hydrator ports.BatchHydrator, store ports.CompletedPullStore, account ports.InstallationAccount, now func() time.Time) ports.ListingIngestor {
	if now == nil {
		now = time.Now
	}
	return &listingIngestor{hydrator: hydrator, store: store, account: account, now: now}
}

func (i *listingIngestor) IngestListing(ctx context.Context, providerListingID string) error {
	id := strings.TrimSpace(providerListingID)
	if id == "" {
		return fmt.Errorf("provider listing id is required")
	}
	_, err := ingestBatch(ctx, i.hydrator, i.store, i.account, []string{id}, i.now().UTC())
	return err
}

// backfillCursor is IC-06's listings cursor shape while phase=="backfill":
// {"phase":"backfill","scroll_id":"eyJ...","ids_collected":2400,
// "run_started_at":"2026-07-31T03:00:00Z"}. RunStartedAt is the run's own
// reference time — the SAME value passed as seenAt to every ingestBatch
// call in this run and as runStartedAt to the terminal MarkRunComplete, so
// the keep-absent step's "not seen since" bound covers every tick of the
// run, not just the last one.
type backfillCursor struct {
	Phase        string    `json:"phase"`
	ScrollID     string    `json:"scroll_id,omitempty"`
	IDsCollected int       `json:"ids_collected,omitempty"`
	RunStartedAt time.Time `json:"run_started_at"`
}

// sweepCursor is IC-06's terminal listings cursor:
// {"phase":"sweep","last_full_sweep_at":"2026-07-31T03:41:10Z"}. Kept as its
// own type (rather than reusing backfillCursor with omitted fields) so a
// completed run never round-trips a stale scroll_id/run_started_at — the
// NEXT tick's tolerant parse (parseListingsCursor) treats phase=="sweep" as
// "start a fresh backfill run", which is exactly how a daily-cadence
// scheduler re-walks the whole catalog each cycle (ML has no reliable
// incremental listing search — ADR-07's ratified vocabulary only resolves
// "incremental" to true, never "sweep").
type sweepCursor struct {
	Phase           string    `json:"phase"`
	LastFullSweepAt time.Time `json:"last_full_sweep_at"`
}

// parseListingsCursor is IC-06's tolerant parse: an absent, malformed, or
// non-backfill (e.g. a prior terminal "sweep") cursor all fall through to a
// FRESH backfill run — never an error, matching the scheduler's own
// tolerant-parse convention for unknown phases (scheduler.go's
// inferIncremental). A cursor that parses as phase=="backfill" but is
// missing run_started_at is ALSO treated as fresh: resuming a backfill
// without knowing its own run's reference time would corrupt the
// keep-absent bound, so this is deliberately as strict as a missing cursor.
func parseListingsCursor(raw json.RawMessage, now func() time.Time) backfillCursor {
	if len(raw) > 0 {
		var cur backfillCursor
		if err := json.Unmarshal(raw, &cur); err == nil && cur.Phase == "backfill" && !cur.RunStartedAt.IsZero() {
			return cur
		}
	}
	return backfillCursor{Phase: "backfill", RunStartedAt: now().UTC()}
}

// NewListingsJob builds F-03's resumable backfill/sweep job body (IC-06's
// JobFunc for the listings entity), modeled on
// sync/composition/products_job.go's NewProductsJob shape but state-machine
// driven by the cursor instead of a stateless per-cycle sync.
//
// Each invocation processes exactly ONE enumeration page (bounded batch in
// memory — IC-06 "batch em memória bounded, tamanho de batch de hidratação,
// não catálogo" — Ingestion.Pull's whole-catalog accumulation is NOT
// inherited here): read one id page, hydrate+upsert it via ingestBatch, and
// either advance the cursor to the next scroll_id (enumeration continues) or
// — when the enumerator reports exhaustion (NextCursor == "", regardless of
// how many ids that final page had) — call MarkRunComplete for the run and
// return the terminal sweep cursor. A page with zero ids but NextCursor=="")
// still triggers MarkRunComplete: this is precisely the "legitimate empty
// terminal tick" the F-02/F-03 split (UpsertPulledRows vs MarkRunComplete)
// exists for.
//
// An enumerator error (e.g. an expired scroll_id, or a 429 series beyond
// M-01's budget) returns the error: the scheduler records last_error and
// keeps the STORED cursor untouched (JobFunc's error path discards the
// returned cursor) — the run stays incomplete, nothing gets marked absent,
// no row upserted so far is lost (IC-06 error matrix).
func NewListingsJob(enumerator ports.IDEnumerator, hydrator ports.BatchHydrator, store ports.CompletedPullStore, account ports.InstallationAccount, now func() time.Time) syncapp.JobFunc {
	if now == nil {
		now = time.Now
	}
	return func(ctx context.Context, rawCursor json.RawMessage) (json.RawMessage, error) {
		cursor := parseListingsCursor(rawCursor, now)

		page, err := enumerator.ReadIDPage(ctx, account, cursor.ScrollID, listingsBackfillPageLimit)
		if err != nil {
			return rawCursor, fmt.Errorf("listings backfill: enumerate page: %w", err)
		}

		if len(page.IDs) > 0 {
			if _, err := ingestBatch(ctx, hydrator, store, account, page.IDs, cursor.RunStartedAt); err != nil {
				return rawCursor, fmt.Errorf("listings backfill: %w", err)
			}
			cursor.IDsCollected += len(page.IDs)
		}

		if strings.TrimSpace(page.NextCursor) == "" {
			if err := store.MarkRunComplete(ctx, account.InstallationID, cursor.RunStartedAt); err != nil {
				return rawCursor, fmt.Errorf("listings backfill: mark run complete: %w", err)
			}
			return json.Marshal(sweepCursor{Phase: "sweep", LastFullSweepAt: now().UTC()})
		}

		cursor.Phase = "backfill"
		cursor.ScrollID = page.NextCursor
		return json.Marshal(cursor)
	}
}
