// Package listings is the context's façade. Everything past this file is
// under internal/; the composition root names a pool and nothing else.
package listings

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/application"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/postgres"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

type Module struct{ service *application.Service }

func New(pool *pgxpool.Pool) *Module {
	return &Module{service: application.NewService(postgres.NewRepository(pool))}
}

// IngestListing folds one channel observation into Listings.
func (m *Module) IngestListing(ctx context.Context, o contracts.ListingObservation) (contracts.IngestResult, error) {
	return m.service.Ingest(ctx, o)
}

// There is deliberately no read method here. This façade carries what a
// caller outside the context actually needs, and nothing is reading listings
// back yet -- an exported operation whose only caller is a test is the orphan
// surface this context's plan rules out by name. The property such a method
// would have proven (what SaveVersion writes is what Current reads back) is
// proven better without it, in tests/integration/listings_ingest_test.go: a
// second fold of an identical observation must report idempotent, which is
// that round-trip taken through the door production uses, and the columns are
// then checked in SQL that no reader of ours touches. When a real consumer
// arrives, it brings its own port and its own reason.

// StoredObservations reads one page of the payloads this context has recorded,
// so the root can hand them back to the channel adapter that produced them and
// fold the re-derived facts in through IngestListing. The bytes leave here as
// opaquely as they arrived.
func (m *Module) StoredObservations(ctx context.Context, t tenant.ID, after contracts.StoredCursor, limit int) (contracts.StoredPage, error) {
	return m.service.StoredObservations(ctx, t, after, limit)
}
