package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/product_links/domain"
	"marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

type ResolutionService struct {
	candidates ports.LinkCandidateStore
	workflows  ports.ProductLinkWorkflowStore
	now        func() time.Time
	newAuditID func() string
}

type ResolutionServiceConfig struct {
	Candidates ports.LinkCandidateStore
	Workflows  ports.ProductLinkWorkflowStore
	Now        func() time.Time
	NewAuditID func() string
}

type ApproveCandidateInput struct {
	CandidateID string
	Actor       domain.ActorMetadata
	Reason      string
}

type RejectListingInput struct {
	InstallationID      string
	ProviderCode        string
	ProviderItemID      string
	ProviderVariationID string
	Actor               domain.ActorMetadata
	Reason              string
}

type ManualResolveInput struct {
	InstallationID        string
	ProviderCode          string
	ProviderItemID        string
	ProviderVariationID   string
	InternalProductID     int
	InternalProductName   string
	InternalReferenceCode string
	Actor                 domain.ActorMetadata
	Reason                string
}

type ListLinkWorkflowsInput struct {
	InstallationID string
	Limit          int
}

func NewResolutionService(cfg ResolutionServiceConfig) *ResolutionService {
	now := cfg.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	newAuditID := cfg.NewAuditID
	if newAuditID == nil {
		newAuditID = func() string { return fmt.Sprintf("pla_%d", now().UTC().UnixNano()) }
	}
	return &ResolutionService{
		candidates: cfg.Candidates,
		workflows:  cfg.Workflows,
		now:        now,
		newAuditID: newAuditID,
	}
}

func (s *ResolutionService) ApproveCandidate(ctx context.Context, input ApproveCandidateInput) (domain.ProductLinkResolutionResult, error) {
	if s.candidates == nil || s.workflows == nil {
		return domain.ProductLinkResolutionResult{}, errors.New("PRODUCT_LINKS_RESOLUTION_NOT_CONFIGURED")
	}
	candidateID := strings.TrimSpace(input.CandidateID)
	if candidateID == "" {
		return domain.ProductLinkResolutionResult{}, errors.New("PRODUCT_LINKS_CANDIDATE_REQUIRED")
	}
	candidate, found, err := s.candidates.GetLinkCandidate(ctx, candidateID)
	if err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	if !found {
		return domain.ProductLinkResolutionResult{}, errors.New("PRODUCT_LINKS_CANDIDATE_NOT_FOUND")
	}
	if candidate.InternalProductID == nil {
		return domain.ProductLinkResolutionResult{}, errors.New("PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE")
	}
	identity := domain.ListingIdentity{
		InstallationID:      candidate.InstallationID,
		ProviderItemID:      candidate.ProviderItemID,
		ProviderVariationID: candidate.ProviderVariationID,
	}
	current, found, err := s.workflows.GetProductLink(ctx, identity)
	if err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	fallbackState := domain.ProductLinkStateNone
	if !found {
		fallbackState = candidateStateToProductLinkState(candidate.State)
	}
	result := s.buildTransition(current, found, fallbackState, buildResolvedLinkInput{
		InstallationID:        candidate.InstallationID,
		ProviderCode:          candidate.ProviderCode,
		ProviderItemID:        candidate.ProviderItemID,
		ProviderVariationID:   candidate.ProviderVariationID,
		SourceCandidateID:     candidate.CandidateID,
		InternalProductID:     candidate.InternalProductID,
		InternalProductName:   candidate.InternalProductName,
		InternalReferenceCode: candidate.InternalReferenceCode,
		Action:                domain.ProductLinkActionApproveCandidate,
		Reason:                input.Reason,
		Actor:                 input.Actor,
	})
	if err := s.workflows.ApplyProductLinkTransition(ctx, result); err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	return domain.ProductLinkResolutionResult{Link: result.Link, Audit: result.Audit}, nil
}

func (s *ResolutionService) RejectListing(ctx context.Context, input RejectListingInput) (domain.ProductLinkResolutionResult, error) {
	if s.workflows == nil {
		return domain.ProductLinkResolutionResult{}, errors.New("PRODUCT_LINKS_RESOLUTION_NOT_CONFIGURED")
	}
	identity, providerCode, err := normalizeIdentity(input.InstallationID, input.ProviderCode, input.ProviderItemID, input.ProviderVariationID)
	if err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	current, found, err := s.workflows.GetProductLink(ctx, identity)
	if err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	result := s.buildTransition(current, found, domain.ProductLinkStateNone, buildRejectedLinkInput{
		InstallationID:      identity.InstallationID,
		ProviderCode:        providerCode,
		ProviderItemID:      identity.ProviderItemID,
		ProviderVariationID: identity.ProviderVariationID,
		Action:              domain.ProductLinkActionRejectListing,
		Reason:              input.Reason,
		Actor:               input.Actor,
	})
	if err := s.workflows.ApplyProductLinkTransition(ctx, result); err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	return domain.ProductLinkResolutionResult{Link: result.Link, Audit: result.Audit}, nil
}

func (s *ResolutionService) ManualResolve(ctx context.Context, input ManualResolveInput) (domain.ProductLinkResolutionResult, error) {
	if s.workflows == nil {
		return domain.ProductLinkResolutionResult{}, errors.New("PRODUCT_LINKS_RESOLUTION_NOT_CONFIGURED")
	}
	if input.InternalProductID == 0 {
		return domain.ProductLinkResolutionResult{}, errors.New("PRODUCT_LINKS_INTERNAL_PRODUCT_REQUIRED")
	}
	identity, providerCode, err := normalizeIdentity(input.InstallationID, input.ProviderCode, input.ProviderItemID, input.ProviderVariationID)
	if err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	current, found, err := s.workflows.GetProductLink(ctx, identity)
	if err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	internalProductID := input.InternalProductID
	result := s.buildTransition(current, found, domain.ProductLinkStateNone, buildResolvedLinkInput{
		InstallationID:        identity.InstallationID,
		ProviderCode:          providerCode,
		ProviderItemID:        identity.ProviderItemID,
		ProviderVariationID:   identity.ProviderVariationID,
		InternalProductID:     &internalProductID,
		InternalProductName:   strings.TrimSpace(input.InternalProductName),
		InternalReferenceCode: strings.TrimSpace(input.InternalReferenceCode),
		Action:                domain.ProductLinkActionManualResolve,
		Reason:                input.Reason,
		Actor:                 input.Actor,
	})
	if err := s.workflows.ApplyProductLinkTransition(ctx, result); err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	return domain.ProductLinkResolutionResult{Link: result.Link, Audit: result.Audit}, nil
}

func (s *ResolutionService) ListLinkWorkflows(ctx context.Context, input ListLinkWorkflowsInput) ([]domain.ProductLinkWorkflowItem, error) {
	if s.candidates == nil || s.workflows == nil {
		return nil, errors.New("PRODUCT_LINKS_RESOLUTION_NOT_CONFIGURED")
	}
	installationID := strings.TrimSpace(input.InstallationID)
	if installationID == "" {
		return nil, errors.New("PRODUCT_LINKS_INSTALLATION_REQUIRED")
	}
	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}
	candidates, err := s.candidates.ListLinkCandidates(ctx, installationID, limit)
	if err != nil {
		return nil, err
	}
	links, err := s.workflows.ListProductLinks(ctx, installationID, limit)
	if err != nil {
		return nil, err
	}
	audits, err := s.workflows.ListProductLinkAuditEntries(ctx, installationID, limit*5)
	if err != nil {
		return nil, err
	}
	itemsByKey := map[string]*domain.ProductLinkWorkflowItem{}
	order := make([]string, 0, len(candidates)+len(links))
	addItem := func(identity domain.ListingIdentity) *domain.ProductLinkWorkflowItem {
		key := identityKey(identity.InstallationID, identity.ProviderItemID, identity.ProviderVariationID)
		item, ok := itemsByKey[key]
		if ok {
			return item
		}
		item = &domain.ProductLinkWorkflowItem{
			Identity:   identity,
			Candidates: []domain.LinkCandidate{},
			Audit:      []domain.ProductLinkAuditEntry{},
		}
		itemsByKey[key] = item
		order = append(order, key)
		return item
	}
	for _, candidate := range candidates {
		item := addItem(domain.ListingIdentity{
			InstallationID:      candidate.InstallationID,
			ProviderItemID:      candidate.ProviderItemID,
			ProviderVariationID: candidate.ProviderVariationID,
		})
		item.Candidates = append(item.Candidates, candidate)
	}
	for _, link := range links {
		item := addItem(domain.ListingIdentity{
			InstallationID:      link.InstallationID,
			ProviderItemID:      link.ProviderItemID,
			ProviderVariationID: link.ProviderVariationID,
		})
		linkCopy := link
		item.CurrentLink = &linkCopy
	}
	for _, audit := range audits {
		item := addItem(domain.ListingIdentity{
			InstallationID:      audit.InstallationID,
			ProviderItemID:      audit.ProviderItemID,
			ProviderVariationID: audit.ProviderVariationID,
		})
		item.Audit = append(item.Audit, audit)
	}
	items := make([]domain.ProductLinkWorkflowItem, 0, len(order))
	for _, key := range order {
		items = append(items, *itemsByKey[key])
	}
	return items, nil
}

type buildResolvedLinkInput struct {
	InstallationID        string
	ProviderCode          string
	ProviderItemID        string
	ProviderVariationID   string
	SourceCandidateID     string
	InternalProductID     *int
	InternalProductName   string
	InternalReferenceCode string
	Action                domain.ProductLinkAction
	Reason                string
	Actor                 domain.ActorMetadata
}

type buildRejectedLinkInput struct {
	InstallationID      string
	ProviderCode        string
	ProviderItemID      string
	ProviderVariationID string
	Action              domain.ProductLinkAction
	Reason              string
	Actor               domain.ActorMetadata
}

func (s *ResolutionService) buildTransition(current domain.ProductLink, found bool, fallbackPreviousState domain.ProductLinkState, input any) domain.ProductLinkTransition {
	now := s.now().UTC()
	prevState := fallbackPreviousState
	var prevProductID *int
	createdAt := now
	if found {
		prevState = current.State
		prevProductID = current.InternalProductID
		createdAt = current.CreatedAt
	}
	switch typed := input.(type) {
	case buildResolvedLinkInput:
		link := domain.ProductLink{
			InstallationID:        typed.InstallationID,
			ProviderCode:          typed.ProviderCode,
			ProviderItemID:        typed.ProviderItemID,
			ProviderVariationID:   typed.ProviderVariationID,
			State:                 domain.ProductLinkStateResolved,
			SourceCandidateID:     typed.SourceCandidateID,
			InternalProductID:     typed.InternalProductID,
			InternalProductName:   typed.InternalProductName,
			InternalReferenceCode: typed.InternalReferenceCode,
			CreatedAt:             createdAt,
			UpdatedAt:             now,
		}
		return domain.ProductLinkTransition{
			Link: link,
			Audit: domain.ProductLinkAuditEntry{
				AuditID:                s.newAuditID(),
				InstallationID:         typed.InstallationID,
				ProviderCode:           typed.ProviderCode,
				ProviderItemID:         typed.ProviderItemID,
				ProviderVariationID:    typed.ProviderVariationID,
				Action:                 typed.Action,
				Reason:                 strings.TrimSpace(typed.Reason),
				SourceCandidateID:      typed.SourceCandidateID,
				Actor:                  typed.Actor,
				PreviousState:          prevState,
				NextState:              domain.ProductLinkStateResolved,
				PreviousInternalProductID: prevProductID,
				NextInternalProductID:  typed.InternalProductID,
				CreatedAt:              now,
			},
		}
	case buildRejectedLinkInput:
		link := domain.ProductLink{
			InstallationID:      typed.InstallationID,
			ProviderCode:        typed.ProviderCode,
			ProviderItemID:      typed.ProviderItemID,
			ProviderVariationID: typed.ProviderVariationID,
			State:               domain.ProductLinkStateRejected,
			CreatedAt:           createdAt,
			UpdatedAt:           now,
		}
		return domain.ProductLinkTransition{
			Link: link,
			Audit: domain.ProductLinkAuditEntry{
				AuditID:                s.newAuditID(),
				InstallationID:         typed.InstallationID,
				ProviderCode:           typed.ProviderCode,
				ProviderItemID:         typed.ProviderItemID,
				ProviderVariationID:    typed.ProviderVariationID,
				Action:                 typed.Action,
				Reason:                 strings.TrimSpace(typed.Reason),
				Actor:                  typed.Actor,
				PreviousState:          prevState,
				NextState:              domain.ProductLinkStateRejected,
				PreviousInternalProductID: prevProductID,
				CreatedAt:              now,
			},
		}
	default:
		panic("unsupported product link transition input")
	}
}

func candidateStateToProductLinkState(state domain.LinkCandidateState) domain.ProductLinkState {
	switch state {
	case domain.LinkCandidateStateConflict:
		return domain.ProductLinkStateConflict
	case domain.LinkCandidateStateUnresolved, domain.LinkCandidateStateTitleMatch, domain.LinkCandidateStateManual:
		return domain.ProductLinkStateUnresolved
	default:
		return domain.ProductLinkStateNone
	}
}

func normalizeIdentity(installationID, providerCode, providerItemID, providerVariationID string) (domain.ListingIdentity, string, error) {
	installationID = strings.TrimSpace(installationID)
	providerCode = strings.TrimSpace(providerCode)
	providerItemID = strings.TrimSpace(providerItemID)
	providerVariationID = strings.TrimSpace(providerVariationID)
	if installationID == "" {
		return domain.ListingIdentity{}, "", errors.New("PRODUCT_LINKS_INSTALLATION_REQUIRED")
	}
	if providerCode == "" || providerItemID == "" {
		return domain.ListingIdentity{}, "", errors.New("PRODUCT_LINKS_LISTING_IDENTITY_REQUIRED")
	}
	return domain.ListingIdentity{
		InstallationID:      installationID,
		ProviderItemID:      providerItemID,
		ProviderVariationID: providerVariationID,
	}, providerCode, nil
}

func identityKey(installationID, providerItemID, providerVariationID string) string {
	return installationID + "|" + providerItemID + "|" + providerVariationID
}
