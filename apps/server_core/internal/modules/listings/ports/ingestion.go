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

type CompletedPullStore interface {
	ApplyCompletedPull(ctx context.Context, installationID string, rows []domain.Listing, completedAt time.Time) error
}
