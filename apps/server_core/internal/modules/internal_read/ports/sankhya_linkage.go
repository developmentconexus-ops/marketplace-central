package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type SankhyaCandidateInput struct {
	ExternalOrderKey string
}

type SankhyaDescendantInput struct {
	Origin                 domain.InternalDocumentLine
	ExpectedOriginQuantity *float64
}

// SankhyaLinkageReader is deliberately separate from the broad Reader port.
// Implementations must fail closed on configuration validation before reads.
type SankhyaLinkageReader interface {
	ValidateConfiguration(ctx context.Context) error
	FindCandidates(ctx context.Context, input SankhyaCandidateInput) ([]domain.SankhyaDocumentCandidate, error)
	ListDescendants(ctx context.Context, input SankhyaDescendantInput) (domain.SankhyaLineage, error)
}
