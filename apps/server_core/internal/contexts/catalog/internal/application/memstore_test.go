package application_test

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/application"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/domain"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// memStore is an in-memory Store. It is a test double for the decision logic
// only: it proves nothing about persistence, atomicity or isolation, and the
// integration lane in Task 10 is what covers those.
type memStore struct {
	byID  map[string]domain.Product
	order []string
}

func newMemStore() *memStore {
	return &memStore{byID: map[string]domain.Product{}}
}

func (m *memStore) BySourceKey(_ context.Context, k contracts.SourceProductKey) (domain.Product, bool, error) {
	for _, id := range m.order {
		p := m.byID[id]
		for _, have := range p.SourceKeys() {
			if have.String() == k.String() {
				return p, true, nil
			}
		}
	}
	return domain.Product{}, false, nil
}

func (m *memStore) ByIdentifier(_ context.Context, t tenant.ID, id contracts.Identifier) ([]domain.Product, error) {
	var out []domain.Product
	for _, key := range m.order {
		p := m.byID[key]
		if p.Tenant() != t {
			continue
		}
		for _, have := range p.Identifiers() {
			if have.String() == id.String() {
				out = append(out, p)
				break
			}
		}
	}
	return out, nil
}

func (m *memStore) Insert(_ context.Context, p domain.Product) error {
	if _, exists := m.byID[p.ID().String()]; exists {
		return fmt.Errorf("memstore: duplicate insert of %s", p.ID())
	}
	m.byID[p.ID().String()] = p
	m.order = append(m.order, p.ID().String())
	return nil
}

func (m *memStore) Update(_ context.Context, p domain.Product) error {
	if _, exists := m.byID[p.ID().String()]; !exists {
		return fmt.Errorf("memstore: update of unknown %s", p.ID())
	}
	m.byID[p.ID().String()] = p
	return nil
}

// seqIDs hands out predictable canonical ids so a test can name them.
type seqIDs struct{ n int }

func (s *seqIDs) NewProductID() (domain.ProductID, error) {
	s.n++
	return domain.NewProductID(fmt.Sprintf("%s%032d", domain.ProductIDPrefix, s.n))
}

var _ application.Store = (*memStore)(nil)
var _ application.IDFactory = (*seqIDs)(nil)
