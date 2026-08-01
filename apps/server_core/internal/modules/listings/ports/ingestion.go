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

// CompletedPullStore persists one pull's worth of rows. runStartedAt is the
// run's own reference time: it bounds "not seen since" for the keep-absent
// step and is stamped onto last_seen_at/updated_at for every upserted row —
// a row upserted in this run therefore never satisfies last_seen_at <
// runStartedAt, which is what naturally excludes it from being marked absent.
// complete tells the store whether the caller's enumeration + hydration
// fully drained (IC-06); when false (aborted/rate-limited/killed run), the
// store must upsert whatever rows it did receive but must never mark any
// row absent — a partial run persists partial progress without punishing
// what it didn't get to see.
type CompletedPullStore interface {
	ApplyCompletedPull(ctx context.Context, installationID string, rows []domain.Listing, runStartedAt time.Time, complete bool) error
}
