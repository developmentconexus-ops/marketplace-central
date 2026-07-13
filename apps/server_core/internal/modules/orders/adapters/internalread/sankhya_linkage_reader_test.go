package internalread

import (
	"context"
	"errors"
	"testing"
	"time"

	internaldomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type fakeSankhyaLinkageSource struct {
	candidateInput  internalports.SankhyaCandidateInput
	descendantInput internalports.SankhyaDescendantInput
	candidates      []internaldomain.SankhyaDocumentCandidate
	lineage         internaldomain.SankhyaLineage
	err             error
}

func (f *fakeSankhyaLinkageSource) ValidateConfiguration(context.Context) error { return f.err }
func (f *fakeSankhyaLinkageSource) FindCandidates(_ context.Context, input internalports.SankhyaCandidateInput) ([]internaldomain.SankhyaDocumentCandidate, error) {
	f.candidateInput = input
	return f.candidates, f.err
}
func (f *fakeSankhyaLinkageSource) ListDescendants(_ context.Context, input internalports.SankhyaDescendantInput) (internaldomain.SankhyaLineage, error) {
	f.descendantInput = input
	return f.lineage, f.err
}

func TestSankhyaLinkageReaderConvertsGenericCandidatesAndNullableFacts(t *testing.T) {
	documentNumber, productID := int64(44), int64(99)
	quantity, amount := 2.0, 10.5
	negotiatedAt := time.Date(2026, 7, 13, 10, 0, 0, 0, time.UTC)
	source := &fakeSankhyaLinkageSource{candidates: []internaldomain.SankhyaDocumentCandidate{{
		DocumentID: 31301, DocumentNumber: &documentNumber, OperationCode: 313, NegotiatedAt: &negotiatedAt, Amount: nil,
		Lines: []internaldomain.SankhyaLineEvidence{{
			Identity: internaldomain.InternalDocumentLine{DocumentID: 31301, LineNumber: 7}, ProductID: &productID, Quantity: &quantity, Amount: &amount,
		}},
	}}}
	reader, err := NewSankhyaLinkageReader(source, "cfg-runtime-1")
	if err != nil {
		t.Fatalf("NewSankhyaLinkageReader() error = %v", err)
	}

	got, err := reader.FindCandidates(context.Background(), "ml:v1:install:order")
	if err != nil {
		t.Fatalf("FindCandidates() error = %v", err)
	}
	if source.candidateInput.ExternalOrderKey != "ml:v1:install:order" || len(got) != 1 || got[0].Header.DocumentID != 31301 || got[0].Amount != nil {
		t.Fatalf("converted candidates = %#v input=%#v", got, source.candidateInput)
	}
	if len(got[0].Lines) != 1 || got[0].Lines[0].Identity.LineNumber != 7 || got[0].Lines[0].ProductID == nil {
		t.Fatalf("converted candidate lines = %#v", got[0].Lines)
	}
}

func TestSankhyaLinkageReaderConvertsExactLineageStateAndInput(t *testing.T) {
	expected, attended := 3.0, 2.0
	source := &fakeSankhyaLinkageSource{lineage: internaldomain.SankhyaLineage{
		Origin: internaldomain.InternalDocumentLine{DocumentID: 31301, LineNumber: 2}, State: internaldomain.SankhyaLineagePartial,
		Descendants: []internaldomain.SankhyaDescendant{{Identity: internaldomain.InternalDocumentLine{DocumentID: 30601, LineNumber: 4}, AttendedQuantity: &attended}},
	}}
	reader, err := NewSankhyaLinkageReader(source, "cfg-runtime-1")
	if err != nil {
		t.Fatalf("NewSankhyaLinkageReader() error = %v", err)
	}
	origin := ordersdomain.InternalDocumentLineIdentity{DocumentID: 31301, LineNumber: 2}

	got, err := reader.ListDescendants(context.Background(), origin, &expected)
	if err != nil {
		t.Fatalf("ListDescendants() error = %v", err)
	}
	if source.descendantInput.Origin.DocumentID != origin.DocumentID || source.descendantInput.ExpectedOriginQuantity == nil || *source.descendantInput.ExpectedOriginQuantity != expected {
		t.Fatalf("descendant input = %#v", source.descendantInput)
	}
	if got.State != ordersdomain.AssistedSankhyaLineagePartial || len(got.Descendants) != 1 || got.Descendants[0].Identity.DocumentID != 30601 {
		t.Fatalf("converted lineage = %#v", got)
	}
}

func TestSankhyaLinkageReaderTranslatesInternalErrors(t *testing.T) {
	source := &fakeSankhyaLinkageSource{err: internaldomain.NewReadError(internaldomain.ReadErrorCandidateAmbiguous, "internal detail", errors.New("driver detail"))}
	reader, constructorErr := NewSankhyaLinkageReader(source, "cfg-runtime-1")
	if constructorErr != nil {
		t.Fatalf("NewSankhyaLinkageReader() error = %v", constructorErr)
	}
	_, err := reader.FindCandidates(context.Background(), "ml:v1:i:o")
	var readErr *ordersdomain.AssistedSankhyaReadError
	if !errors.As(err, &readErr) || readErr.Kind != ordersdomain.AssistedSankhyaReadCandidateConflict {
		t.Fatalf("translated error = %v, want orders candidate conflict", err)
	}
	var internalErr *internaldomain.ReadError
	if errors.As(err, &internalErr) {
		t.Fatalf("internal read error leaked through adapter: %T", err)
	}
}

func TestSankhyaLinkageReaderRequiresAndExposesRuntimeRevision(t *testing.T) {
	if _, err := NewSankhyaLinkageReader(&fakeSankhyaLinkageSource{}, "  "); err == nil {
		t.Fatal("NewSankhyaLinkageReader() error = nil, want empty revision rejected")
	}
	reader, err := NewSankhyaLinkageReader(&fakeSankhyaLinkageSource{}, " cfg-runtime-7 ")
	if err != nil {
		t.Fatalf("NewSankhyaLinkageReader() error = %v", err)
	}
	if got := reader.ConfigurationRevision(); got != "cfg-runtime-7" {
		t.Fatalf("ConfigurationRevision() = %q", got)
	}
}
