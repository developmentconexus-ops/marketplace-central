package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

type ListAssistedSankhyaCandidatesInput struct {
	TenantID        string
	InstallationID  string
	ProviderOrderID string
}

type ConfirmAssistedSankhyaLinkageInput struct {
	TenantID           string
	InstallationID     string
	ProviderOrderID    string
	SelectedDocumentID int64
	Selections         []domain.AssistedSankhyaLineSelection
	ActorID            string
	Reason             string
	IdempotencyKey     string
	SourceAt           time.Time
}

type AssistedSankhyaEventIDGenerator func() (string, error)

type AssistedSankhyaLinkageService struct {
	orders     ports.OrderLookup
	reader     ports.SankhyaLinkageReader
	linkages   ports.SankhyaLinkageRepository
	now        func() time.Time
	newEventID AssistedSankhyaEventIDGenerator
}

type AssistedSankhyaLinkageServiceConfig struct {
	Orders     ports.OrderLookup
	Reader     ports.SankhyaLinkageReader
	Linkages   ports.SankhyaLinkageRepository
	Now        func() time.Time
	NewEventID AssistedSankhyaEventIDGenerator
}

func NewAssistedSankhyaLinkageService(config AssistedSankhyaLinkageServiceConfig) *AssistedSankhyaLinkageService {
	now := config.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	newEventID := config.NewEventID
	if newEventID == nil {
		newEventID = newAssistedSankhyaEventID
	}
	return &AssistedSankhyaLinkageService{orders: config.Orders, reader: config.Reader, linkages: config.Linkages, now: now, newEventID: newEventID}
}

func (s *AssistedSankhyaLinkageService) ListCandidates(ctx context.Context, input ListAssistedSankhyaCandidatesInput) (domain.AssistedSankhyaCandidateResult, error) {
	scope := normalizedLinkageScope(input.TenantID, input.InstallationID, input.ProviderOrderID)
	if _, err := s.loadExactOrder(ctx, scope); err != nil {
		return domain.AssistedSankhyaCandidateResult{}, err
	}
	if s.reader == nil {
		return domain.AssistedSankhyaCandidateResult{}, &domain.AssistedSankhyaReadError{Kind: domain.AssistedSankhyaReadConfigurationInvalid}
	}
	if err := s.reader.ValidateConfiguration(ctx); err != nil {
		return domain.AssistedSankhyaCandidateResult{}, err
	}
	externalKey := domain.ExternalOrderKeyFor(scope)
	candidates, err := s.reader.FindCandidates(ctx, externalKey)
	if err != nil {
		return domain.AssistedSankhyaCandidateResult{}, err
	}
	return domain.AssistedSankhyaCandidateResult{Scope: scope, ExternalOrderKey: externalKey, Candidates: candidates}, nil
}

func (s *AssistedSankhyaLinkageService) Confirm(ctx context.Context, input ConfirmAssistedSankhyaLinkageInput) (domain.AssistedSankhyaConfirmationResult, error) {
	scope := normalizedLinkageScope(input.TenantID, input.InstallationID, input.ProviderOrderID)
	order, err := s.loadExactOrder(ctx, scope)
	if err != nil {
		return domain.AssistedSankhyaConfirmationResult{}, err
	}
	if s.reader == nil || s.linkages == nil {
		return domain.AssistedSankhyaConfirmationResult{}, &domain.AssistedSankhyaReadError{Kind: domain.AssistedSankhyaReadConfigurationInvalid}
	}
	if input.SelectedDocumentID <= 0 {
		return domain.AssistedSankhyaConfirmationResult{}, assistedInvalid("selected_document_id")
	}
	if err := s.reader.ValidateConfiguration(ctx); err != nil {
		return domain.AssistedSankhyaConfirmationResult{}, err
	}
	configurationRevision := strings.TrimSpace(s.reader.ConfigurationRevision())
	if configurationRevision == "" {
		return domain.AssistedSankhyaConfirmationResult{}, &domain.AssistedSankhyaReadError{Kind: domain.AssistedSankhyaReadConfigurationInvalid}
	}
	externalKey := domain.ExternalOrderKeyFor(scope)
	candidates, err := s.reader.FindCandidates(ctx, externalKey)
	if err != nil {
		return domain.AssistedSankhyaConfirmationResult{}, err
	}
	selected, err := selectExactCandidate(candidates, input.SelectedDocumentID)
	if err != nil {
		return domain.AssistedSankhyaConfirmationResult{}, err
	}
	mappings, err := domain.ValidateAssistedSankhyaSelection(order, selected, input.Selections)
	if err != nil {
		return domain.AssistedSankhyaConfirmationResult{}, err
	}
	eventID, err := s.newEventID()
	if err != nil {
		return domain.AssistedSankhyaConfirmationResult{}, fmt.Errorf("generate assisted sankhya event id: %w", err)
	}
	eventID = strings.TrimSpace(eventID)
	if eventID == "" {
		return domain.AssistedSankhyaConfirmationResult{}, assistedInvalid("event_id")
	}

	linkage := domain.SankhyaLinkage{
		Scope:            scope,
		ExternalOrderKey: externalKey,
		Header:           selected.Header,
		Lines:            mappings,
		Audit: domain.LinkageAudit{
			EventID:               eventID,
			EventType:             domain.LinkageEventConfirmed,
			ActorType:             domain.AssistedSankhyaActorOperatorSuppliedUnverified,
			ActorID:               strings.TrimSpace(input.ActorID),
			Reason:                strings.TrimSpace(input.Reason),
			IdempotencyKey:        strings.TrimSpace(input.IdempotencyKey),
			SourceAt:              input.SourceAt.UTC(),
			RecordedAt:            s.now().UTC(),
			ConfigurationRevision: configurationRevision,
			EvidenceState:         domain.LinkageEvidenceExact,
		},
	}
	if err := linkage.Validate(); err != nil {
		return domain.AssistedSankhyaConfirmationResult{}, err
	}
	persisted, err := s.linkages.AppendConfirmation(ctx, linkage)
	if err != nil {
		return domain.AssistedSankhyaConfirmationResult{}, err
	}

	result := domain.AssistedSankhyaConfirmationResult{Linkage: persisted}
	expectedQuantity := candidateQuantityByOrigin(selected)
	for _, mapping := range persisted.Lines {
		lineage, readErr := s.reader.ListDescendants(ctx, mapping.Origin, expectedQuantity[mapping.Origin])
		if readErr != nil {
			result.Lines = append(result.Lines, domain.AssistedSankhyaLineage{
				MPCLineID: mapping.MPCLineID,
				Origin:    mapping.Origin,
				State:     lineageErrorState(readErr),
			})
			continue
		}
		if !validLineageForOrigin(lineage, mapping.Origin) {
			result.Lines = append(result.Lines, domain.AssistedSankhyaLineage{
				MPCLineID: mapping.MPCLineID,
				Origin:    mapping.Origin,
				State:     domain.AssistedSankhyaLineageConflict,
			})
			continue
		}
		lineage.MPCLineID = mapping.MPCLineID
		result.Lines = append(result.Lines, lineage)
	}
	return result, nil
}

func (s *AssistedSankhyaLinkageService) loadExactOrder(ctx context.Context, scope domain.LinkageScope) (domain.MarketplaceOrder, error) {
	if s.orders == nil || strings.TrimSpace(scope.TenantID) == "" || strings.TrimSpace(scope.InstallationID) == "" || strings.TrimSpace(scope.ProviderOrderID) == "" {
		return domain.MarketplaceOrder{}, assistedInvalid("scope")
	}
	order, found, err := s.orders.FindExactOrder(ctx, scope)
	if err != nil {
		return domain.MarketplaceOrder{}, err
	}
	if !found {
		return domain.MarketplaceOrder{}, domain.ErrAssistedSankhyaOrderNotFound
	}
	if order.InstallationID != scope.InstallationID || order.ProviderOrderID != scope.ProviderOrderID || order.ProviderCode != "mercado_livre" {
		return domain.MarketplaceOrder{}, domain.ErrAssistedSankhyaOrderNotFound
	}
	return order, nil
}

func normalizedLinkageScope(tenantID, installationID, providerOrderID string) domain.LinkageScope {
	return domain.LinkageScope{
		TenantID:        strings.TrimSpace(tenantID),
		InstallationID:  strings.TrimSpace(installationID),
		ProviderOrderID: strings.TrimSpace(providerOrderID),
	}
}

func selectExactCandidate(candidates []domain.AssistedSankhyaCandidate, documentID int64) (domain.AssistedSankhyaCandidate, error) {
	var selected domain.AssistedSankhyaCandidate
	matches := 0
	for _, candidate := range candidates {
		if candidate.Header.DocumentID == documentID {
			selected = candidate
			matches++
		}
	}
	if matches != 1 || selected.OperationCode != 313 {
		return domain.AssistedSankhyaCandidate{}, assistedInvalid("selected_candidate")
	}
	return selected, nil
}

func candidateQuantityByOrigin(candidate domain.AssistedSankhyaCandidate) map[domain.InternalDocumentLineIdentity]*float64 {
	result := make(map[domain.InternalDocumentLineIdentity]*float64, len(candidate.Lines))
	for _, line := range candidate.Lines {
		result[line.Identity] = line.Quantity
	}
	return result
}

func lineageErrorState(err error) domain.AssistedSankhyaLineageState {
	var readErr *domain.AssistedSankhyaReadError
	if errors.As(err, &readErr) && (readErr.Kind == domain.AssistedSankhyaReadConfigurationInvalid || readErr.Kind == domain.AssistedSankhyaReadUnavailable) {
		return domain.AssistedSankhyaLineageUnavailable
	}
	return domain.AssistedSankhyaLineageConflict
}

func validLineageForOrigin(lineage domain.AssistedSankhyaLineage, origin domain.InternalDocumentLineIdentity) bool {
	if lineage.Origin != origin {
		return false
	}
	switch lineage.State {
	case domain.AssistedSankhyaLineageNone:
		return len(lineage.Descendants) == 0
	case domain.AssistedSankhyaLineagePartial:
		if len(lineage.Descendants) == 0 {
			return false
		}
	case domain.AssistedSankhyaLineageComplete:
		if len(lineage.Descendants) == 0 {
			return false
		}
	case domain.AssistedSankhyaLineageConflict:
		return true
	default:
		return false
	}
	seen := make(map[domain.InternalDocumentLineIdentity]struct{}, len(lineage.Descendants))
	for _, descendant := range lineage.Descendants {
		if descendant.Identity.DocumentID <= 0 || descendant.Identity.LineNumber <= 0 {
			return false
		}
		if _, duplicate := seen[descendant.Identity]; duplicate {
			return false
		}
		seen[descendant.Identity] = struct{}{}
		if descendant.AttendedQuantity == nil {
			if lineage.State == domain.AssistedSankhyaLineageComplete {
				return false
			}
			continue
		}
		quantity := *descendant.AttendedQuantity
		if quantity < 0 || math.IsNaN(quantity) || math.IsInf(quantity, 0) {
			return false
		}
	}
	return true
}

func newAssistedSankhyaEventID() (string, error) {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}
	return "asl_" + hex.EncodeToString(value[:]), nil
}

func assistedInvalid(field string) error {
	return fmt.Errorf("%w: %s", domain.ErrInvalidAssistedSankhyaLinkage, field)
}
