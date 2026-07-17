package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

// CollectorPort is the contract for a future market-data collector. This
// mission defines the interface only; no production implementation is wired.
type CollectorPort interface {
	CollectObservations(ctx context.Context, installationID string, listingIDs []string) ([]domain.MarketObservation, error)
	CollectReferences(ctx context.Context, productIDs []string) ([]domain.MarketReference, error)
}
