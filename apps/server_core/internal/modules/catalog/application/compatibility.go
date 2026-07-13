package application

import (
	"strconv"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/catalog/domain"
)

type CompatibilityStatus string

const (
	CompatibilityMapped           CompatibilityStatus = "mapped"
	CompatibilityNotFound         CompatibilityStatus = "not_found"
	CompatibilityIdentityConflict CompatibilityStatus = "identity_conflict"
)

// ResolveLegacyProductReference accepts only exact positive decimal CODPROD
// equality. EAN, references, titles, and provider SKU are intentionally not
// compatibility keys.
func ResolveLegacyProductReference(legacy string, candidates []domain.CanonicalProduct) (*domain.InternalProductID, CompatibilityStatus) {
	n, err := strconv.Atoi(strings.TrimSpace(legacy))
	if err != nil || n <= 0 {
		return nil, CompatibilityNotFound
	}
	var matched []domain.InternalProductID
	for _, candidate := range candidates {
		if int(candidate.InternalProductID) == n {
			matched = append(matched, candidate.InternalProductID)
		}
	}
	if len(matched) == 0 {
		return nil, CompatibilityNotFound
	}
	if len(matched) != 1 {
		return nil, CompatibilityIdentityConflict
	}
	return &matched[0], CompatibilityMapped
}
