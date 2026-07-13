package application

import (
	"marketplace-central/apps/server_core/internal/modules/catalog/domain"
	"testing"
)

func TestResolveLegacyProductReference(t *testing.T) {
	product := func(id int) domain.CanonicalProduct {
		return domain.CanonicalProduct{InternalProductID: domain.InternalProductID(id)}
	}
	cases := []struct {
		name, legacy string
		candidates   []domain.CanonicalProduct
		status       CompatibilityStatus
		id           *domain.InternalProductID
	}{
		{"mapped", "42664", []domain.CanonicalProduct{product(42664)}, CompatibilityMapped, ptrID(42664)},
		{"unmapped string", "REF-42664", []domain.CanonicalProduct{product(42664)}, CompatibilityNotFound, nil},
		{"unmapped numeric", "42664", []domain.CanonicalProduct{product(9)}, CompatibilityNotFound, nil},
		{"ambiguous", "42664", []domain.CanonicalProduct{product(42664), product(42664)}, CompatibilityIdentityConflict, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id, status := ResolveLegacyProductReference(tc.legacy, tc.candidates)
			if status != tc.status {
				t.Fatalf("status=%s", status)
			}
			if (id == nil) != (tc.id == nil) || id != nil && *id != *tc.id {
				t.Fatalf("id=%v", id)
			}
		})
	}
}
func ptrID(v int) *domain.InternalProductID { id := domain.InternalProductID(v); return &id }
