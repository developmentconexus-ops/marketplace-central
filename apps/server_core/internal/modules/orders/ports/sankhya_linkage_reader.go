package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type SankhyaLinkageReader interface {
	ValidateConfiguration(ctx context.Context) error
	ConfigurationRevision() string
	FindCandidates(ctx context.Context, externalOrderKey string) ([]domain.AssistedSankhyaCandidate, error)
	ListDescendants(ctx context.Context, origin domain.InternalDocumentLineIdentity, expectedOriginQuantity *float64) (domain.AssistedSankhyaLineage, error)
}
