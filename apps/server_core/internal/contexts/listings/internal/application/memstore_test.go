package application_test

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/application"
)

type memStore struct {
	current map[string]application.CurrentListing
	saves   int
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

func (m *memStore) SaveVersion(_ context.Context, o contracts.ListingObservation, version int) error {
	m.saves++
	m.current[memKey(o.Key)] = application.CurrentListing{Version: version, PayloadHash: o.Evidence.PayloadHash()}
	return nil
}
