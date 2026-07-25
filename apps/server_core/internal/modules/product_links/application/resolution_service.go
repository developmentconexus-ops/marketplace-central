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
	// BatchID tags the resulting audit entry with the batch it was
	// approved under (S3 ApplyBatch). Empty for a standalone approval —
	// additive field, existing callers default to "".
	BatchID string
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

// UndoResolutionInput is the S4 single-resolution undo request: the audit
// row identifying the transition to reverse, plus the acting operator
// (threaded into the reversal's own audit entry — the reversal is a new
// audited action, never an anonymous rewrite).
type UndoResolutionInput struct {
	AuditID string
	Actor   domain.ActorMetadata
	Reason  string
}

// UndoBatchItem is the itemized per-item outcome of a batch-undo run. An
// item that is already SUPERSEDED/ALREADY_UNDONE is itemized as FAILED —
// it is never fatal to the rest of the batch (mirrors ApplyBatch's
// partial-failure semantics).
type UndoBatchItem struct {
	AuditID     string `json:"audit_id"`
	CandidateID string `json:"candidate_id,omitempty"`
	Status      string `json:"status"`
	Cause       string `json:"cause,omitempty"`
}

// UndoBatchResult is the itemized outcome of a batch-undo run.
type UndoBatchResult struct {
	BatchID  string          `json:"batch_id"`
	Reverted []UndoBatchItem `json:"reverted"`
	Failed   []UndoBatchItem `json:"failed"`
}

// ErrProductLinkAuditNotFound signals UndoResolution/UndoBatch found no
// audit row for the given id/batch (module sentinel style, mapped to 404 by
// the existing "_NOT_FOUND" suffix convention).
var ErrProductLinkAuditNotFound = errors.New("PRODUCT_LINKS_AUDIT_NOT_FOUND")

// ErrProductLinkBatchNotFound signals UndoBatch found no audit rows tagged
// with the given batch_id.
var ErrProductLinkBatchNotFound = errors.New("PRODUCT_LINKS_BATCH_NOT_FOUND")

// ErrProductLinkSuperseded (409) — the audit entry targeted by
// UndoResolution is no longer the latest action for its listing identity: a
// newer resolution has since replaced it, so reversing it would not be a
// safe/coherent undo of the CURRENT state.
var ErrProductLinkSuperseded = errors.New("SUPERSEDED")

// ErrProductLinkAlreadyUndone (409) — the latest action for the target
// audit's listing identity is itself already an undo, so there is nothing
// further to reverse.
var ErrProductLinkAlreadyUndone = errors.New("ALREADY_UNDONE")

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
	if err := domain.ValidateInternalProductID(*candidate.InternalProductID); err != nil {
		return domain.ProductLinkResolutionResult{}, err
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
		BatchID:               input.BatchID,
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
	if err := domain.ValidateInternalProductID(input.InternalProductID); err != nil {
		return domain.ProductLinkResolutionResult{}, err
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

// UndoResolution reverses a single prior resolution transition (S4). The
// decision between proceeding, SUPERSEDED, and ALREADY_UNDONE is derived
// entirely from the audit chain (never from time/random):
//
//  1. Load the target audit entry by AuditID. Missing ⇒ ErrProductLinkAuditNotFound.
//  2. Load the latest audit entry for that same listing identity.
//  3. If the latest action is itself "undo" ⇒ ErrProductLinkAlreadyUndone —
//     the identity's current state has already been reverted (this also
//     covers re-undoing the same target: its own reversal becomes the new
//     latest, so a second UndoResolution(sameAuditID) sees latest.Action ==
//     undo before it ever compares audit ids).
//  4. Else if the latest audit id differs from the target's ⇒
//     ErrProductLinkSuperseded — a newer resolution has replaced the target
//     since it was written; reversing the target would not reflect the
//     identity's current state.
//  5. Else (target IS the latest, non-undo action) ⇒ write a REVERSAL
//     transition: state <- target.PreviousState, internal product <-
//     target.PreviousInternalProductID, action = undo. This is a NEW audit
//     row (history preserved), never a delete/rewrite of the target.
func (s *ResolutionService) UndoResolution(ctx context.Context, input UndoResolutionInput) (domain.ProductLinkResolutionResult, error) {
	if s.workflows == nil {
		return domain.ProductLinkResolutionResult{}, errors.New("PRODUCT_LINKS_RESOLUTION_NOT_CONFIGURED")
	}
	auditID := strings.TrimSpace(input.AuditID)
	if auditID == "" {
		return domain.ProductLinkResolutionResult{}, errors.New("PRODUCT_LINKS_AUDIT_REQUIRED")
	}
	target, found, err := s.workflows.GetAuditEntry(ctx, auditID)
	if err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	if !found {
		return domain.ProductLinkResolutionResult{}, ErrProductLinkAuditNotFound
	}
	return s.undoAuditEntry(ctx, target, input.Actor, input.Reason)
}

// UndoBatch fans out UndoResolution over every original (non-undo) audit
// row tagged with batchID (S4). Each item is undone independently — an item
// that is already SUPERSEDED/ALREADY_UNDONE is itemized as failed, never
// fatal to the rest of the batch.
func (s *ResolutionService) UndoBatch(ctx context.Context, batchID string) (UndoBatchResult, error) {
	if s.workflows == nil {
		return UndoBatchResult{}, errors.New("PRODUCT_LINKS_RESOLUTION_NOT_CONFIGURED")
	}
	batchID = strings.TrimSpace(batchID)
	if batchID == "" {
		return UndoBatchResult{}, errors.New("PRODUCT_LINKS_BATCH_REQUIRED")
	}
	entries, err := s.workflows.ListAuditByBatch(ctx, batchID)
	if err != nil {
		return UndoBatchResult{}, err
	}
	// Only original (non-undo) rows are undo targets. Reversal rows this
	// same method writes are tagged with batchID too (traceability) but
	// must never be re-selected as targets.
	originals := make([]domain.ProductLinkAuditEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.Action != domain.ProductLinkActionUndo {
			originals = append(originals, entry)
		}
	}
	if len(originals) == 0 {
		return UndoBatchResult{}, ErrProductLinkBatchNotFound
	}
	result := UndoBatchResult{BatchID: batchID}
	for _, target := range originals {
		item := UndoBatchItem{AuditID: target.AuditID, CandidateID: target.SourceCandidateID}
		if _, err := s.undoAuditEntry(ctx, target, domain.ActorMetadata{ActorType: "operator"}, fmt.Sprintf("batch undo %s", batchID)); err != nil {
			item.Status = BatchItemStatusFailed
			item.Cause = err.Error()
			result.Failed = append(result.Failed, item)
			continue
		}
		item.Status = BatchItemStatusOK
		result.Reverted = append(result.Reverted, item)
	}
	return result, nil
}

// undoAuditEntry applies the SUPERSEDED/ALREADY_UNDONE ordering check
// against the latest audit for target's identity, then (if clear) writes
// the reversal transition. Shared by UndoResolution and UndoBatch so both
// paths agree on the exact same ordering decision.
func (s *ResolutionService) undoAuditEntry(ctx context.Context, target domain.ProductLinkAuditEntry, actor domain.ActorMetadata, reason string) (domain.ProductLinkResolutionResult, error) {
	identity := domain.ListingIdentity{
		InstallationID:      target.InstallationID,
		ProviderItemID:      target.ProviderItemID,
		ProviderVariationID: target.ProviderVariationID,
	}
	latest, found, err := s.workflows.LatestAuditForIdentity(ctx, identity)
	if err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	if !found {
		// Cannot happen if target itself was persisted, but guard honestly
		// rather than assume.
		return domain.ProductLinkResolutionResult{}, ErrProductLinkAuditNotFound
	}
	if latest.Action == domain.ProductLinkActionUndo {
		return domain.ProductLinkResolutionResult{}, ErrProductLinkAlreadyUndone
	}
	if latest.AuditID != target.AuditID {
		return domain.ProductLinkResolutionResult{}, ErrProductLinkSuperseded
	}

	current, foundLink, err := s.workflows.GetProductLink(ctx, identity)
	if err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	now := s.now().UTC()
	createdAt := now
	prevStateForAudit := domain.ProductLinkStateNone
	var prevProductForAudit *int
	if foundLink {
		createdAt = current.CreatedAt
		prevStateForAudit = current.State
		prevProductForAudit = current.InternalProductID
	}

	// The audit chain carries the previous product ID but not its descriptive
	// fields (name / reference code / source candidate). When the reversal
	// lands back on the SAME internal product the link already points at,
	// those descriptors are still true, so carry them over instead of
	// blanking the row (the operator saw "—" where a product name belonged).
	// When the reversal points at a DIFFERENT product — or at none — we have
	// no honest source for its name here, so it stays empty (ADR-17) rather
	// than inheriting the previous product's description.
	link := domain.ProductLink{
		InstallationID:      target.InstallationID,
		ProviderCode:        target.ProviderCode,
		ProviderItemID:      target.ProviderItemID,
		ProviderVariationID: target.ProviderVariationID,
		State:               target.PreviousState,
		InternalProductID:   target.PreviousInternalProductID,
		CreatedAt:           createdAt,
		UpdatedAt:           now,
	}
	if foundLink && sameInternalProduct(current.InternalProductID, target.PreviousInternalProductID) {
		link.InternalProductName = current.InternalProductName
		link.InternalReferenceCode = current.InternalReferenceCode
		link.SourceCandidateID = current.SourceCandidateID
	}
	audit := domain.ProductLinkAuditEntry{
		AuditID:                   s.newAuditID(),
		InstallationID:            target.InstallationID,
		ProviderCode:              target.ProviderCode,
		ProviderItemID:            target.ProviderItemID,
		ProviderVariationID:       target.ProviderVariationID,
		Action:                    domain.ProductLinkActionUndo,
		Reason:                    strings.TrimSpace(reason),
		Actor:                     actor,
		PreviousState:             prevStateForAudit,
		NextState:                 target.PreviousState,
		PreviousInternalProductID: prevProductForAudit,
		NextInternalProductID:     target.PreviousInternalProductID,
		BatchID:                   target.BatchID,
		CreatedAt:                 now,
	}
	transition := domain.ProductLinkTransition{Link: link, Audit: audit}
	if err := s.workflows.ApplyProductLinkTransition(ctx, transition); err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	return domain.ProductLinkResolutionResult{Link: link, Audit: audit}, nil
}

// sameInternalProduct compares two optional internal product IDs by VALUE.
// Two absent IDs are "the same" only in the sense that neither carries a
// product, and that case never reaches the descriptor carry-over above.
func sameInternalProduct(a, b *int) bool {
	if a == nil || b == nil {
		return false
	}
	return *a == *b
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
	BatchID               string
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
				AuditID:                   s.newAuditID(),
				InstallationID:            typed.InstallationID,
				ProviderCode:              typed.ProviderCode,
				ProviderItemID:            typed.ProviderItemID,
				ProviderVariationID:       typed.ProviderVariationID,
				Action:                    typed.Action,
				Reason:                    strings.TrimSpace(typed.Reason),
				SourceCandidateID:         typed.SourceCandidateID,
				Actor:                     typed.Actor,
				PreviousState:             prevState,
				NextState:                 domain.ProductLinkStateResolved,
				PreviousInternalProductID: prevProductID,
				NextInternalProductID:     typed.InternalProductID,
				BatchID:                   typed.BatchID,
				CreatedAt:                 now,
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
				AuditID:                   s.newAuditID(),
				InstallationID:            typed.InstallationID,
				ProviderCode:              typed.ProviderCode,
				ProviderItemID:            typed.ProviderItemID,
				ProviderVariationID:       typed.ProviderVariationID,
				Action:                    typed.Action,
				Reason:                    strings.TrimSpace(typed.Reason),
				Actor:                     typed.Actor,
				PreviousState:             prevState,
				NextState:                 domain.ProductLinkStateRejected,
				PreviousInternalProductID: prevProductID,
				CreatedAt:                 now,
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
