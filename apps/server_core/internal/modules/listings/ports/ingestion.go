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
