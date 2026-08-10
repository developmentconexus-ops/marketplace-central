// Package listings is the context's façade. Everything past this file is
// under internal/; the composition root names a pool and nothing else.
package listings

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/application"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/postgres"
)

type Module struct{ service *application.Service }

func New(pool *pgxpool.Pool) *Module {
	return &Module{service: application.NewService(postgres.NewRepository(pool))}
}

// IngestListing folds one channel observation into Listings.
func (m *Module) IngestListing(ctx context.Context, o contracts.ListingObservation) (contracts.IngestResult, error) {
	return m.service.Ingest(ctx, o)
}
