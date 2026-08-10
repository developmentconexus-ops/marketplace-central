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

// CurrentState reads back what Listings has stored for a listing, translated
// to the context's public vocabulary.
//
// Its caller is the round-trip assertion in tests/integration
// (listings_ingest_test.go): that test lives outside internal/contexts/listings
// and so cannot name the repository, which is exactly the boundary this façade
// states. Without this method the only way to prove that what SaveVersion
// wrote is what Current reads back — the property the ingest fold now depends
// on — would be to reach past the façade or to assert on raw columns, and
// neither proves the state a fold actually compares.
func (m *Module) CurrentState(ctx context.Context, key contracts.SourceListingKey) (contracts.ListingState, int, bool, error) {
	current, found, err := m.service.Current(ctx, key)
	if err != nil || !found {
		return contracts.ListingState{}, 0, found, err
	}
	return current.State, current.Version, true, nil
}

// StoredObservations reads one page of the payloads this context has recorded,
// so the root can hand them back to the channel adapter that produced them and
// fold the re-derived facts in through IngestListing. The bytes leave here as
// opaquely as they arrived.
func (m *Module) StoredObservations(ctx context.Context, t tenant.ID, after contracts.StoredCursor, limit int) (contracts.StoredPage, error) {
	return m.service.StoredObservations(ctx, t, after, limit)
}
