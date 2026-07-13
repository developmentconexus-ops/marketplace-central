package application

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

type SankhyaLinkageService struct {
	reader ports.SankhyaLinkageReader
}

func NewSankhyaLinkageService(reader ports.SankhyaLinkageReader) SankhyaLinkageService {
	return SankhyaLinkageService{reader: reader}
}

func (s SankhyaLinkageService) ValidateConfiguration(ctx context.Context) error {
	if s.reader == nil {
		return domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage reader is not configured", nil)
	}
	return s.reader.ValidateConfiguration(ctx)
}

func (s SankhyaLinkageService) FindCandidates(ctx context.Context, input ports.SankhyaCandidateInput) ([]domain.SankhyaDocumentCandidate, error) {
	if s.reader == nil {
		return nil, domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage reader is not configured", nil)
	}
	return s.reader.FindCandidates(ctx, input)
}

func (s SankhyaLinkageService) ListDescendants(ctx context.Context, input ports.SankhyaDescendantInput) (domain.SankhyaLineage, error) {
	if s.reader == nil {
		return domain.SankhyaLineage{}, domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage reader is not configured", nil)
	}
	return s.reader.ListDescendants(ctx, input)
}
