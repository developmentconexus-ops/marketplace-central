package application

import (
	"context"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

type fakeSankhyaLinkageReader struct {
	validated   int
	candidates  int
	descendants int
}

func (f *fakeSankhyaLinkageReader) ValidateConfiguration(context.Context) error {
	f.validated++
	return nil
}

func (f *fakeSankhyaLinkageReader) FindCandidates(context.Context, ports.SankhyaCandidateInput) ([]domain.SankhyaDocumentCandidate, error) {
	f.candidates++
	return []domain.SankhyaDocumentCandidate{{DocumentID: 31301, OperationCode: 313}}, nil
}

func (f *fakeSankhyaLinkageReader) ListDescendants(context.Context, ports.SankhyaDescendantInput) (domain.SankhyaLineage, error) {
	f.descendants++
	return domain.SankhyaLineage{State: domain.SankhyaLineageNone}, nil
}

func TestSankhyaLinkageServiceUsesSeparateReaderContract(t *testing.T) {
	reader := &fakeSankhyaLinkageReader{}
	service := NewSankhyaLinkageService(reader)
	ctx := context.Background()

	if err := service.ValidateConfiguration(ctx); err != nil {
		t.Fatalf("ValidateConfiguration() error = %v", err)
	}
	if _, err := service.FindCandidates(ctx, ports.SankhyaCandidateInput{ExternalOrderKey: "ml:v1:account:order"}); err != nil {
		t.Fatalf("FindCandidates() error = %v", err)
	}
	if _, err := service.ListDescendants(ctx, ports.SankhyaDescendantInput{Origin: domain.InternalDocumentLine{DocumentID: 31301, LineNumber: 1}}); err != nil {
		t.Fatalf("ListDescendants() error = %v", err)
	}
	if reader.validated != 1 || reader.candidates != 1 || reader.descendants != 1 {
		t.Fatalf("reader calls = validate:%d candidates:%d descendants:%d", reader.validated, reader.candidates, reader.descendants)
	}
}

func TestSankhyaLinkageServiceFailsClosedWithoutReader(t *testing.T) {
	service := NewSankhyaLinkageService(nil)
	_, err := service.FindCandidates(context.Background(), ports.SankhyaCandidateInput{ExternalOrderKey: "ml:v1:account:order"})
	if !domain.IsReadErrorCode(err, domain.ReadErrorConfigurationInvalid) {
		t.Fatalf("error = %v, want configuration_invalid", err)
	}
}
