package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

const (
	serviceLineOne domain.MPCLineID = "mpl_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	serviceLineTwo domain.MPCLineID = "mpl_bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
)

type fakeAssistedOrderLookup struct {
	order domain.MarketplaceOrder
	found bool
	err   error
	scope domain.LinkageScope
}

func (f *fakeAssistedOrderLookup) FindExactOrder(_ context.Context, scope domain.LinkageScope) (domain.MarketplaceOrder, bool, error) {
	f.scope = scope
	return f.order, f.found, f.err
}

type descendantCall struct {
	origin   domain.InternalDocumentLineIdentity
	expected *float64
}

type fakeAssistedReader struct {
	configurationErr      error
	configurationRevision string
	candidates            []domain.AssistedSankhyaCandidate
	candidateErr          error
	keys                  []string
	validateCalls         int
	descendants           map[domain.InternalDocumentLineIdentity]domain.AssistedSankhyaLineage
	descendantErrors      map[domain.InternalDocumentLineIdentity]error
	descendantCalls       []descendantCall
}

func (f *fakeAssistedReader) ValidateConfiguration(context.Context) error {
	f.validateCalls++
	return f.configurationErr
}

func (f *fakeAssistedReader) ConfigurationRevision() string { return f.configurationRevision }

func (f *fakeAssistedReader) FindCandidates(_ context.Context, key string) ([]domain.AssistedSankhyaCandidate, error) {
	f.keys = append(f.keys, key)
	return f.candidates, f.candidateErr
}

func (f *fakeAssistedReader) ListDescendants(_ context.Context, origin domain.InternalDocumentLineIdentity, expected *float64) (domain.AssistedSankhyaLineage, error) {
	f.descendantCalls = append(f.descendantCalls, descendantCall{origin: origin, expected: expected})
	if err := f.descendantErrors[origin]; err != nil {
		return domain.AssistedSankhyaLineage{}, err
	}
	if lineage, ok := f.descendants[origin]; ok {
		return lineage, nil
	}
	return domain.AssistedSankhyaLineage{Origin: origin, State: domain.AssistedSankhyaLineageNone}, nil
}

type fakeAssistedLinkageRepository struct {
	appended []domain.SankhyaLinkage
	result   domain.SankhyaLinkage
	err      error
}

func (f *fakeAssistedLinkageRepository) LoadCurrent(context.Context, domain.LinkageScope) (domain.SankhyaLinkage, bool, error) {
	return domain.SankhyaLinkage{}, false, nil
}

func (f *fakeAssistedLinkageRepository) AppendConfirmation(_ context.Context, linkage domain.SankhyaLinkage) (domain.SankhyaLinkage, error) {
	f.appended = append(f.appended, linkage)
	if f.err != nil {
		return domain.SankhyaLinkage{}, f.err
	}
	if f.result.Header.DocumentID > 0 {
		return f.result, nil
	}
	return linkage, nil
}

func TestAssistedSankhyaListCandidatesUsesOnlyPersistedExactMLKey(t *testing.T) {
	service, orders, reader, repo, _, _ := assistedServiceFixture()

	got, err := service.ListCandidates(context.Background(), ListAssistedSankhyaCandidatesInput{
		TenantID: " tenant-1 ", InstallationID: " install-1 ", ProviderOrderID: " order-1 ",
	})
	if err != nil {
		t.Fatalf("ListCandidates() error = %v", err)
	}
	if orders.scope != (domain.LinkageScope{TenantID: "tenant-1", InstallationID: "install-1", ProviderOrderID: "order-1"}) {
		t.Fatalf("lookup scope = %#v", orders.scope)
	}
	if len(reader.keys) != 1 || reader.keys[0] != "ml:v1:install-1:order-1" || got.ExternalOrderKey != reader.keys[0] {
		t.Fatalf("exact keys = %#v result=%q", reader.keys, got.ExternalOrderKey)
	}
	if reader.validateCalls != 1 || len(repo.appended) != 0 || len(got.Candidates) != 1 {
		t.Fatalf("validation=%d appends=%d candidates=%d", reader.validateCalls, len(repo.appended), len(got.Candidates))
	}
}

func TestAssistedSankhyaConfirmAppendsExactAuditThenReadsEveryPersistedOrigin(t *testing.T) {
	service, _, reader, repo, input, recordedAt := assistedServiceFixture()
	originOne := reader.candidates[0].Lines[0].Identity
	originTwo := reader.candidates[0].Lines[1].Identity
	persisted := expectedLinkage(input, reader.candidates[0], recordedAt, "generated-event-1", "cfg-runtime-1")
	persisted.Audit.EventID = "persisted-original-event"
	persisted.Audit.RecordedAt = recordedAt.Add(-time.Minute)
	persisted.Lines[0], persisted.Lines[1] = persisted.Lines[1], persisted.Lines[0]
	repo.result = persisted
	reader.descendantErrors[originTwo] = &domain.AssistedSankhyaReadError{Kind: domain.AssistedSankhyaReadUnavailable}
	attended := 1.0
	reader.descendants[originOne] = domain.AssistedSankhyaLineage{
		Origin: originOne, State: domain.AssistedSankhyaLineageComplete,
		Descendants: []domain.AssistedSankhyaDescendant{{Identity: domain.InternalDocumentLineIdentity{DocumentID: 30601, LineNumber: 1}, AttendedQuantity: &attended}},
	}

	got, err := service.Confirm(context.Background(), input)
	if err != nil {
		t.Fatalf("Confirm() error = %v", err)
	}
	if got.Linkage.Audit.EventID != "persisted-original-event" || len(repo.appended) != 1 {
		t.Fatalf("result linkage/appends = %#v/%d, want persisted idempotent aggregate", got.Linkage.Audit, len(repo.appended))
	}
	appended := repo.appended[0]
	if appended.ExternalOrderKey != "ml:v1:install-1:order-1" || appended.Audit.EvidenceState != domain.LinkageEvidenceExact ||
		appended.Audit.ActorType != domain.AssistedSankhyaActorOperatorSuppliedUnverified || appended.Audit.ActorID != "operator-7" ||
		appended.Audit.EventID != "generated-event-1" || appended.Audit.ConfigurationRevision != "cfg-runtime-1" {
		t.Fatalf("appended exact audit = %#v", appended)
	}
	if len(got.Lines) != 2 || len(reader.descendantCalls) != 2 {
		t.Fatalf("lineage results/calls = %d/%d, want every persisted line", len(got.Lines), len(reader.descendantCalls))
	}
	if reader.descendantCalls[0].origin != persisted.Lines[0].Origin || reader.descendantCalls[1].origin != persisted.Lines[1].Origin {
		t.Fatalf("descendant calls = %#v, want persisted mapping order", reader.descendantCalls)
	}
	states := map[domain.InternalDocumentLineIdentity]domain.AssistedSankhyaLineageState{}
	for _, line := range got.Lines {
		states[line.Origin] = line.State
	}
	if states[originOne] != domain.AssistedSankhyaLineageComplete || states[originTwo] != domain.AssistedSankhyaLineageUnavailable {
		t.Fatalf("states = %#v, want complete plus unavailable", states)
	}
	if reader.descendantCalls[0].expected == nil && reader.descendantCalls[1].expected == nil {
		t.Fatal("expected candidate quantities were not propagated")
	}
}

func TestAssistedSankhyaConfirmRejectsInvalidSelectionBeforeAppend(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*fakeAssistedOrderLookup, *fakeAssistedReader, *ConfirmAssistedSankhyaLinkageInput)
	}{
		{name: "non candidate", mutate: func(_ *fakeAssistedOrderLookup, _ *fakeAssistedReader, input *ConfirmAssistedSankhyaLinkageInput) {
			input.SelectedDocumentID = 999
		}},
		{name: "wrong TOP", mutate: func(_ *fakeAssistedOrderLookup, reader *fakeAssistedReader, _ *ConfirmAssistedSankhyaLinkageInput) {
			reader.candidates[0].OperationCode = 306
		}},
		{name: "missing line", mutate: func(_ *fakeAssistedOrderLookup, _ *fakeAssistedReader, input *ConfirmAssistedSankhyaLinkageInput) {
			input.Selections = input.Selections[:1]
		}},
		{name: "legacy line", mutate: func(orders *fakeAssistedOrderLookup, _ *fakeAssistedReader, _ *ConfirmAssistedSankhyaLinkageInput) {
			orders.order.Items[0].ReconciliationState = domain.LineReconciliationLegacyUnresolved
		}},
		{name: "different document line", mutate: func(_ *fakeAssistedOrderLookup, _ *fakeAssistedReader, input *ConfirmAssistedSankhyaLinkageInput) {
			input.Selections[0].CandidateLine.DocumentID = 777
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, orders, reader, repo, input, _ := assistedServiceFixture()
			test.mutate(orders, reader, &input)
			if _, err := service.Confirm(context.Background(), input); !errors.Is(err, domain.ErrInvalidAssistedSankhyaLinkage) {
				t.Fatalf("Confirm() error = %v, want invalid", err)
			}
			if len(repo.appended) != 0 {
				t.Fatalf("appends = %d, want none before full validation", len(repo.appended))
			}
		})
	}
}

func TestAssistedSankhyaConfirmPreservesRepositoryConflict(t *testing.T) {
	service, _, _, repo, input, _ := assistedServiceFixture()
	repo.err = domain.ErrSankhyaLinkageConflict
	if _, err := service.Confirm(context.Background(), input); !errors.Is(err, domain.ErrSankhyaLinkageConflict) {
		t.Fatalf("Confirm() error = %v, want linkage conflict", err)
	}
}

func TestAssistedSankhyaConfirmFailsBeforeAppendWithoutRuntimeRevision(t *testing.T) {
	service, _, reader, repo, input, _ := assistedServiceFixture()
	reader.configurationRevision = ""
	_, err := service.Confirm(context.Background(), input)
	var readErr *domain.AssistedSankhyaReadError
	if !errors.As(err, &readErr) || readErr.Kind != domain.AssistedSankhyaReadConfigurationInvalid {
		t.Fatalf("Confirm() error = %v, want configuration invalid", err)
	}
	if len(repo.appended) != 0 || len(reader.keys) != 0 {
		t.Fatalf("candidate reads/appends = %d/%d, want none without runtime revision", len(reader.keys), len(repo.appended))
	}
}

func TestAssistedSankhyaConfirmFailsBeforeAppendWhenEventIDGenerationFails(t *testing.T) {
	service, _, _, repo, input, _ := assistedServiceFixture()
	generationErr := errors.New("entropy unavailable")
	service.newEventID = func() (string, error) { return "", generationErr }
	if _, err := service.Confirm(context.Background(), input); !errors.Is(err, generationErr) {
		t.Fatalf("Confirm() error = %v, want event ID generation error", err)
	}
	if len(repo.appended) != 0 {
		t.Fatalf("appends = %d, want none after event ID failure", len(repo.appended))
	}
}

func TestDefaultAssistedSankhyaEventIDsAreServerGeneratedAndDistinct(t *testing.T) {
	service := NewAssistedSankhyaLinkageService(AssistedSankhyaLinkageServiceConfig{})
	first, err := service.newEventID()
	if err != nil {
		t.Fatalf("first event ID: %v", err)
	}
	second, err := service.newEventID()
	if err != nil {
		t.Fatalf("second event ID: %v", err)
	}
	if len(first) != 36 || len(second) != 36 || first[:4] != "asl_" || second[:4] != "asl_" || first == second {
		t.Fatalf("generated event IDs = %q, %q", first, second)
	}
}

func TestAssistedSankhyaPostAppendLineageFailuresStayExact(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*fakeAssistedReader, domain.InternalDocumentLineIdentity)
		want      domain.AssistedSankhyaLineageState
	}{
		{name: "typed lineage conflict", configure: func(reader *fakeAssistedReader, origin domain.InternalDocumentLineIdentity) {
			reader.descendantErrors[origin] = &domain.AssistedSankhyaReadError{Kind: domain.AssistedSankhyaReadLineageConflict}
		}, want: domain.AssistedSankhyaLineageConflict},
		{name: "unexpected error fails conflict", configure: func(reader *fakeAssistedReader, origin domain.InternalDocumentLineIdentity) {
			reader.descendantErrors[origin] = errors.New("unexpected read failure")
		}, want: domain.AssistedSankhyaLineageConflict},
		{name: "malformed returned origin", configure: func(reader *fakeAssistedReader, origin domain.InternalDocumentLineIdentity) {
			reader.descendants[origin] = domain.AssistedSankhyaLineage{Origin: domain.InternalDocumentLineIdentity{DocumentID: 999, LineNumber: 1}, State: domain.AssistedSankhyaLineageNone}
		}, want: domain.AssistedSankhyaLineageConflict},
		{name: "source unavailable", configure: func(reader *fakeAssistedReader, origin domain.InternalDocumentLineIdentity) {
			reader.descendantErrors[origin] = &domain.AssistedSankhyaReadError{Kind: domain.AssistedSankhyaReadUnavailable}
		}, want: domain.AssistedSankhyaLineageUnavailable},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, _, reader, _, input, _ := assistedServiceFixture()
			origin := reader.candidates[0].Lines[0].Identity
			test.configure(reader, origin)
			got, err := service.Confirm(context.Background(), input)
			if err != nil {
				t.Fatalf("Confirm() error = %v", err)
			}
			states := map[domain.InternalDocumentLineIdentity]domain.AssistedSankhyaLineageState{}
			for _, line := range got.Lines {
				states[line.Origin] = line.State
			}
			if states[origin] != test.want {
				t.Fatalf("lineage state = %q, want %q", states[origin], test.want)
			}
		})
	}
}

func TestAssistedSankhyaRequiresExactPersistedMercadoLivreOrder(t *testing.T) {
	service, orders, reader, repo, input, _ := assistedServiceFixture()
	orders.order.ProviderCode = "other"
	if _, err := service.Confirm(context.Background(), input); !errors.Is(err, domain.ErrAssistedSankhyaOrderNotFound) {
		t.Fatalf("Confirm() error = %v, want exact ML order not found", err)
	}
	if reader.validateCalls != 0 || len(repo.appended) != 0 {
		t.Fatalf("reader/appends ran before exact persisted order: %d/%d", reader.validateCalls, len(repo.appended))
	}
}

func assistedServiceFixture() (*AssistedSankhyaLinkageService, *fakeAssistedOrderLookup, *fakeAssistedReader, *fakeAssistedLinkageRepository, ConfirmAssistedSankhyaLinkageInput, time.Time) {
	quantityOne, quantityTwo := 1.0, 2.0
	candidate := domain.AssistedSankhyaCandidate{
		Header: domain.InternalDocumentIdentity{DocumentID: 31301}, OperationCode: 313,
		Lines: []domain.AssistedSankhyaCandidateLine{
			{Identity: domain.InternalDocumentLineIdentity{DocumentID: 31301, LineNumber: 1}, Quantity: &quantityOne},
			{Identity: domain.InternalDocumentLineIdentity{DocumentID: 31301, LineNumber: 2}, Quantity: &quantityTwo},
		},
	}
	orders := &fakeAssistedOrderLookup{found: true, order: domain.MarketplaceOrder{
		InstallationID: "install-1", ProviderCode: "mercado_livre", ProviderOrderID: "order-1",
		Items: []domain.MarketplaceOrderItem{
			{MPCLineID: serviceLineOne, ReconciliationState: domain.LineReconciliationStable},
			{MPCLineID: serviceLineTwo, ReconciliationState: domain.LineReconciliationStable},
		},
	}}
	reader := &fakeAssistedReader{
		configurationRevision: "cfg-runtime-1",
		candidates:            []domain.AssistedSankhyaCandidate{candidate},
		descendants:           make(map[domain.InternalDocumentLineIdentity]domain.AssistedSankhyaLineage),
		descendantErrors:      make(map[domain.InternalDocumentLineIdentity]error),
	}
	repo := &fakeAssistedLinkageRepository{}
	recordedAt := time.Date(2026, 7, 13, 15, 0, 0, 0, time.UTC)
	service := NewAssistedSankhyaLinkageService(AssistedSankhyaLinkageServiceConfig{
		Orders: orders, Reader: reader, Linkages: repo, Now: func() time.Time { return recordedAt },
		NewEventID: func() (string, error) { return "generated-event-1", nil },
	})
	input := ConfirmAssistedSankhyaLinkageInput{
		TenantID: "tenant-1", InstallationID: "install-1", ProviderOrderID: "order-1", SelectedDocumentID: 31301,
		Selections: []domain.AssistedSankhyaLineSelection{
			{MPCLineID: serviceLineOne, CandidateLine: candidate.Lines[0].Identity},
			{MPCLineID: serviceLineTwo, CandidateLine: candidate.Lines[1].Identity},
		},
		ActorID: " operator-7 ", Reason: " exact operator selection ", IdempotencyKey: "idem-1",
		SourceAt: time.Date(2026, 7, 13, 14, 0, 0, 0, time.UTC),
	}
	return service, orders, reader, repo, input, recordedAt
}

func expectedLinkage(input ConfirmAssistedSankhyaLinkageInput, candidate domain.AssistedSankhyaCandidate, recordedAt time.Time, eventID, configurationRevision string) domain.SankhyaLinkage {
	return domain.SankhyaLinkage{
		Scope:            domain.LinkageScope{TenantID: input.TenantID, InstallationID: input.InstallationID, ProviderOrderID: input.ProviderOrderID},
		ExternalOrderKey: "ml:v1:install-1:order-1", Header: candidate.Header,
		Lines: []domain.SankhyaLineMapping{
			{MPCLineID: serviceLineOne, Origin: candidate.Lines[0].Identity},
			{MPCLineID: serviceLineTwo, Origin: candidate.Lines[1].Identity},
		},
		Audit: domain.LinkageAudit{
			EventID: eventID, EventType: domain.LinkageEventConfirmed,
			ActorType: domain.AssistedSankhyaActorOperatorSuppliedUnverified, ActorID: "operator-7", Reason: "exact operator selection",
			IdempotencyKey: input.IdempotencyKey, SourceAt: input.SourceAt, RecordedAt: recordedAt,
			ConfigurationRevision: configurationRevision, EvidenceState: domain.LinkageEvidenceExact,
		},
	}
}
