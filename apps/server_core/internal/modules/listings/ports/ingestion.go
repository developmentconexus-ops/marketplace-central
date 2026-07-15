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
}

type PageSource interface {
	ReadPage(ctx context.Context, account InstallationAccount, cursor string, limit int) (ListingPage, error)
}

type CompletedPullStore interface {
	ApplyCompletedPull(ctx context.Context, installationID string, rows []domain.Listing, completedAt time.Time) error
}
