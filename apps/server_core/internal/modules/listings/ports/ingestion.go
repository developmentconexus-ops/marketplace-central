package ports

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/domain"
)

type InstallationAccount struct {
	TenantID          string
	InstallationID    string
	ProviderCode      string
	ProviderAccountID string
}

type ListingPage struct {
	ProviderItemCount int
	Rows              []domain.Listing
	// NextCursor is the provider-issued cursor for the next page (e.g. the ML
	// scan scroll_id). Empty means the source pages by numeric offset — the
	// ingestion loop then derives the next cursor from the running offset.
	NextCursor string
}

type PageSource interface {
	ReadPage(ctx context.Context, account InstallationAccount, cursor string, limit int) (ListingPage, error)
}

// CompletedPullStore persists a pull, split into two calls so a resumable
// multi-tick run (F-03 backfill/sweep) can upsert progress tick-by-tick
// without a legitimate empty terminal tick (scroll exhausted) disabling
// keep-absent for the whole run — see UpsertPulledRows and MarkRunComplete
// for the exact contract each half owns.
//
// A single-shot caller (e.g. Ingestion.Pull) calls UpsertPulledRows once per
// page loop and then MarkRunComplete once at the end, which reproduces the
// old single-call ApplyCompletedPull(rows, runStartedAt, complete=true)
// behavior exactly. A partial/aborted run (complete=false in the old
// signature) simply never calls MarkRunComplete: it upserts whatever rows it
// received via UpsertPulledRows and stops, which is what "must never mark
// any row absent" now means structurally instead of via a bool flag.
type CompletedPullStore interface {
	// UpsertPulledRows persists one tick's worth of rows. seenAt is the RUN's
	// own reference time (the same value across every tick of a run): it is
	// stamped onto last_seen_at/updated_at for every upserted row, which is
	// what bounds "not seen since" for MarkRunComplete's keep-absent step.
	UpsertPulledRows(ctx context.Context, installationID string, rows []domain.Listing, seenAt time.Time) error

	// MarkRunComplete runs the keep-absent step (ADR-06/IC-06) for a run
	// whose enumeration + hydration fully drained. runStartedAt must be the
	// exact same timestamp passed as seenAt to every UpsertPulledRows call
	// in this run — it both bounds the keep-absent UPDATE and backs the
	// run-scoped safety pin: a run that never upserted anything under this
	// timestamp is a no-op here, never a mass-closure.
	MarkRunComplete(ctx context.Context, installationID string, runStartedAt time.Time) error
}

// BatchIDPage is one enumerator tick (F-03/IC-06): a page of provider item
// ids plus the next enumeration cursor. NextCursor == "" means the scan is
// exhausted, REGARDLESS of len(IDs) — a page can legitimately report zero
// new ids and still be the terminal page (e.g. the catalog count is an exact
// multiple of the page size), and a caller must still treat that as run
// completion, not as "nothing happened this tick".
type BatchIDPage struct {
	IDs        []string
	NextCursor string
}

// IDEnumerator produces ids ONLY — it never hydrates (IC-06 "enumerador
// nunca hidrata"). Implemented by connectors.BackfillSource, wrapping the
// mercado_livre ids-only scan (ListListingsScanIDs, Passo 1).
type IDEnumerator interface {
	ReadIDPage(ctx context.Context, account InstallationAccount, cursor string, limit int) (BatchIDPage, error)
}

// BatchHydrator hydrates one batch of ids into canonical rows — it never
// enumerates (IC-06 "hidratação nunca enumera"). Implemented by
// connectors.MultigetHydrator, wrapping GetItemsMultiget (M-01/F-02) + the
// Passo 2 mapper, and feeding the product-links SnapshotObserver (ADR-13) so
// the new path keeps the re-vínculo matcher fed the same way the old
// single-item ReadListing path does. A batch-level error (e.g. a 429 series
// exceeding M-01's retry budget) fails the whole call — a per-item error
// (one bad id in the batch) never does; the implementation records those
// internally and returns the rows that DID hydrate.
type BatchHydrator interface {
	HydrateBatch(ctx context.Context, account InstallationAccount, ids []string) ([]domain.Listing, error)
}

// ListingIngestor is IC-06's Ingest{Listing} operation: a resource-addressed,
// idempotent, single-item upsert, callable by every producer (refresh,
// webhook worker, ad-hoc resync) through the SAME path the backfill job
// uses. Bound to one installation account at construction — same idiom as
// mutations/adapters/listings.NewResyncWriter's account parameter — so the
// method itself only needs the resource id.
type ListingIngestor interface {
	IngestListing(ctx context.Context, providerListingID string) error
}
