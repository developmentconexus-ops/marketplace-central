package application_test

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/application"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

type memStore struct {
	current map[string]application.CurrentListing
	saves   int
	// stored is the recorded-payload side of the store, in the order the
	// postgres adapter orders it: channel, account, listing id.
	stored []contracts.StoredObservation
	// storedTenant is the tenant the last StoredObservations call asked for,
	// so a test can prove the scope reached the store instead of being
	// dropped somewhere above it.
	storedTenant tenant.ID
	// lastObservation is the full observation handed to the most recent
	// SaveVersion call, kept so tests can inspect what was actually stored
	// (e.g. a specific fact) without that inspection depending on the shape
	// of CurrentListing itself.
	lastObservation contracts.ListingObservation
}

func newMemStore() *memStore {
	return &memStore{current: map[string]application.CurrentListing{}}
}

func memKey(k contracts.SourceListingKey) string {
	return k.Tenant().String() + "|" + k.Account().Channel().String() + "|" + k.Account().External() + "|" + k.ListingID()
}

func (m *memStore) Current(_ context.Context, k contracts.SourceListingKey) (application.CurrentListing, bool, error) {
	c, ok := m.current[memKey(k)]
	return c, ok, nil
}

// StoredObservations mirrors the adapter's contract: rows strictly after the
// cursor, in (channel, account, listing) order, at most limit of them.
func (m *memStore) StoredObservations(_ context.Context, t tenant.ID, after contracts.StoredCursor, limit int) ([]contracts.StoredObservation, error) {
	m.storedTenant = t
	var out []contracts.StoredObservation
	for _, o := range m.stored {
		if o.Key.Tenant() != t {
			continue
		}
		if !after.IsStart() && !afterCursor(o, after) {
			continue
		}
		out = append(out, o)
		if len(out) == limit {
			break
		}
	}
	return out, nil
}

// afterCursor is the row-value comparison the SQL does, spelled out.
func afterCursor(o contracts.StoredObservation, after contracts.StoredCursor) bool {
	row := []string{o.Key.Account().Channel().String(), o.Key.Account().External(), o.Key.ListingID()}
	cur := []string{after.Channel(), after.Account(), after.ListingID()}
	for i := range row {
		if row[i] != cur[i] {
			return row[i] > cur[i]
		}
	}
	return false
}

func (m *memStore) SaveVersion(_ context.Context, o contracts.ListingObservation, version int) error {
	m.saves++
	m.lastObservation = o
	m.current[memKey(o.Key)] = application.CurrentListing{Version: version, State: o.State}
	return nil
}
