package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
)

// ObservationStore appends stored observations and reads the latest stored
// row per requested listing. Appends are insert-only: an existing capture is
// never updated or replaced. LatestObservations omits listings without rows
// and preserves the input order for listings that are present. Tenant scope is
// carried by the repository/adapter instance, not by this port.
type ObservationStore interface {
	AppendObservations(ctx context.Context, observations []domain.MarketObservation) error
	LatestObservations(ctx context.Context, listingIDs []string) ([]domain.MarketObservation, error)
}

// ReferenceStore appends stored references and reads the latest stored row per
// requested product. Appends are insert-only: an existing capture is never
// updated or replaced. LatestReferences omits products without rows and
// preserves the input order for products that are present. Tenant scope is
// carried by the repository/adapter instance, not by this port.
type ReferenceStore interface {
	AppendReferences(ctx context.Context, references []domain.MarketReference) error
	LatestReferences(ctx context.Context, productIDs []string) ([]domain.MarketReference, error)
}
