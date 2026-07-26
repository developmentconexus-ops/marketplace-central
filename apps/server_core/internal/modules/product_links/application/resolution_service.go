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
	candidates    ports.LinkCandidateStore
	workflows     ports.ProductLinkWorkflowStore
	now           func() time.Time
	newAuditID    func() string
	newDecisionID func() string
}

type ResolutionServiceConfig struct {
	Candidates    ports.LinkCandidateStore
	Workflows     ports.ProductLinkWorkflowStore
	Now           func() time.Time
	NewAuditID    func() string
	NewDecisionID func() string
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
	LinkLimit      int
	AuditLimit     int
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
	newDecisionID := cfg.NewDecisionID
	if newDecisionID == nil {
		newDecisionID = func() string { return fmt.Sprintf("pld_%d", now().UTC().UnixNano()) }
	}
	return &ResolutionService{
		candidates:    cfg.Candidates,
		workflows:     cfg.Workflows,
		now:           now,
		newAuditID:    newAuditID,
		newDecisionID: newDecisionID,
	}
}

func (s *ResolutionService) ApproveCandidate(ctx context.Context, input ApproveCandidateInput) (domain.ProductLinkResolutionResult, error) {
	if s.candidates == nil || s.workflows == nil {
		return domain.ProductLinkResolutionResult{}, errors.New("PRODUCT_LINKS_RESOLUTION_NOT_CONFIGURED")
	}
	// AutoApproveCandidate is the only door a machine may come through, and it
	// does not come through this one: it states rule=concordant_codprod_ean for
	// itself. Anything reaching here with actor_type=system is an automation
	// pressing the operator's button, and this path would file it as a human.
	if err := errSystemActorNotPermitted(input.Actor); err != nil {
		return domain.ProductLinkResolutionResult{}, err
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
		// A human is approving, whatever the anchor was. The collision count
		// is not read here: it was the generator's reading at candidate time,
		// and re-deriving it now would record a different moment's fact.
		DecisionRule:  decisionRuleForCandidate(candidate),
		DecisionActor: domain.DecisionActorOperator,
	})
	if err := s.workflows.ApplyProductLinkTransition(ctx, result); err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
	return domain.ProductLinkResolutionResult{Link: result.Link, Audit: result.Audit}, nil
}

// AutoApproveCandidateInput is the corroborated candidate the generator judged
// automatic, plus the collision count IT read when it judged so. The count is
// carried, never re-derived: re-reading the ERP now would record a different
// moment's fact in a row that claims to describe the decision.
// The count is a POINTER because a caller that did not read one must be able
// to say so. A plain int forces the absent case to 0, and 0 reads as "the
// anchor matched nothing" — a fact nobody established, and one the E10 CHECK
// rejects outright (ADR-17 / AC-03).
type AutoApproveCandidateInput struct {
	Candidate            domain.LinkCandidate
	CollisionsAtDecision *int
}

// AutoApproveCandidate applies the single automatic path (D-121-2, ADR-05
// amended): CODPROD and EAN resolving the same product, with no hard negative.
// It runs the SAME transition machine as ApproveCandidate — same link, same
// audit row, same E10 write in the same transaction — and differs only in who
// decided: actor=system, rule=concordant_codprod_ean.
//
// It never overrides a decision already in force. A link the operator resolved
// keeps their answer (M05-C10), and a re-run over an already auto-approved link
// writes nothing rather than a fresh row per sync. Reports whether it approved.
func (s *ResolutionService) AutoApproveCandidate(ctx context.Context, input AutoApproveCandidateInput) (bool, error) {
	if s.workflows == nil {
		return false, errors.New("PRODUCT_LINKS_RESOLUTION_NOT_CONFIGURED")
	}
	candidate := input.Candidate
	if candidate.MatchStatus != domain.LinkCandidateMatchStatusAccept {
		return false, errors.New("PRODUCT_LINKS_AUTO_APPROVE_NOT_CORROBORATED")
	}
	if candidate.InternalProductID == nil {
		return false, errors.New("PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE")
	}
	if err := domain.ValidateInternalProductID(*candidate.InternalProductID); err != nil {
		return false, err
	}
	identity := domain.ListingIdentity{
		InstallationID:      candidate.InstallationID,
		ProviderItemID:      candidate.ProviderItemID,
		ProviderVariationID: candidate.ProviderVariationID,
	}
	current, found, err := s.workflows.GetProductLink(ctx, identity)
	if err != nil {
		return false, err
	}
	// A listing the operator already settled is settled. The link's own state is
	// what proves it: a rejection does write an E10 row, but links resolved
	// before E10 existed carry no decision at all, so the trail alone would
	// report those as undecided and the automatic path would reopen them.
	if found && (current.State == domain.ProductLinkStateRejected || current.State == domain.ProductLinkStateResolved) {
		return false, nil
	}
	decisions, err := s.workflows.ListDecisionsForLink(ctx, identity)
	if err != nil {
		return false, err
	}
	for _, decision := range decisions {
		if decision.SupersededBy == "" {
			return false, nil
		}
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
		Reason:                "auto-vínculo: CODPROD e EAN concordantes",
		Actor:                 domain.ActorMetadata{ActorType: "system", ActorID: "auto_linker"},
		DecisionRule:          domain.DecisionRuleConcordantCodprodEAN,
		DecisionActor:         domain.DecisionActorSystem,
		CollisionsAtDecision:  input.CollisionsAtDecision,
	})
	if err := s.workflows.ApplyProductLinkTransition(ctx, result); err != nil {
		return false, err
	}
	return true, nil
}

func (s *ResolutionService) RejectListing(ctx context.Context, input RejectListingInput) (domain.ProductLinkResolutionResult, error) {
	if s.workflows == nil {
		return domain.ProductLinkResolutionResult{}, errors.New("PRODUCT_LINKS_RESOLUTION_NOT_CONFIGURED")
	}
	if err := errSystemActorNotPermitted(input.Actor); err != nil {
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
	if err := errSystemActorNotPermitted(input.Actor); err != nil {
		return domain.ProductLinkResolutionResult{}, err
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
		// The operator named the product themselves — no anchor decided this,
		// so no collision count was read.
		DecisionRule:  domain.DecisionRuleManual,
		DecisionActor: domain.DecisionActorOperator,
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
	candidateLimit := input.Limit
	if candidateLimit <= 0 {
		candidateLimit = 20
	}
	linkLimit := input.LinkLimit
	if linkLimit <= 0 {
		linkLimit = 2000
	}
	auditLimit := input.AuditLimit
	if auditLimit <= 0 {
		auditLimit = 10000
	}
	candidates, err := s.candidates.ListLinkCandidates(ctx, installationID, candidateLimit)
	if err != nil {
		return nil, err
	}
	links, err := s.workflows.ListProductLinks(ctx, installationID, linkLimit)
	if err != nil {
		return nil, err
	}
	audits, err := s.workflows.ListProductLinkAuditEntries(ctx, installationID, auditLimit)
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
	if err := errSystemActorNotPermitted(actor); err != nil {
		return domain.ProductLinkResolutionResult{}, err
	}
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
	// An undo is a decision too: it says the link that stood should not. Left
	// out of the trail, the decision it reverts keeps reading as in force, and
	// the automatic path would be blocked by a row nobody stands behind — with
	// nothing anywhere saying why. Recording it supersedes that row and names
	// who took it back.
	// errSystemActorNotPermitted turned any system caller away above.
	decision := s.newDecisionRow(identity, domain.DecisionRuleManual, domain.DecisionActorOperator, nil, now)
	transition := domain.ProductLinkTransition{Link: link, Audit: audit, Decision: decision}
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
	// DecisionRule/DecisionActor/CollisionsAtDecision carry the E10 row this
	// approval writes alongside the transition. CollisionsAtDecision stays nil
	// unless a collision count was actually read for this decision (ADR-17).
	DecisionRule         domain.ProductLinkDecisionRule
	DecisionActor        string
	CollisionsAtDecision *int
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
			Link:     link,
			Decision: s.buildDecision(typed, now),
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
			// A rejection is a decision too. It names no anchor — there is no
			// rule behind "this anúncio is not ours" — but leaving it out of
			// the trail lets the decision it overrules keep reading as the one
			// in force, so the live row would report the system as still
			// standing behind a link the operator killed. `manual` is the
			// honest name for a call a human made on no anchor.
			// errSystemActorNotPermitted turned any system caller away above.
			Decision: s.newDecisionRow(
				domain.ListingIdentity{
					InstallationID:      typed.InstallationID,
					ProviderItemID:      typed.ProviderItemID,
					ProviderVariationID: typed.ProviderVariationID,
				},
				domain.DecisionRuleManual,
				domain.DecisionActorOperator,
				nil,
				now,
			),
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

// newDecisionRow is the ONE place an E10 row is assembled. Approve, undo and
// reject all reach it, so a path added later inherits the row's shape instead
// of restating it — the three defects this milestone's gates found were each a
// different hand-built literal disagreeing with the others.
//
// It does NOT decide the (rule, actor) pair. The caller states it, and
// migration 0082's CHECK stays the single authority on which pairs are legal:
// re-deriving that policy here would put a second copy of it in Go, free to
// drift from the one the database enforces.
func (s *ResolutionService) newDecisionRow(identity domain.ListingIdentity, rule domain.ProductLinkDecisionRule, actor string, collisions *int, now time.Time) *domain.ProductLinkDecision {
	return &domain.ProductLinkDecision{
		DecisionID:           s.newDecisionID(),
		InstallationID:       identity.InstallationID,
		ProviderItemID:       identity.ProviderItemID,
		ProviderVariationID:  identity.ProviderVariationID,
		LinkID:               domain.LinkID(identity.InstallationID, identity.ProviderItemID, identity.ProviderVariationID),
		RuleMatched:          rule,
		Actor:                actor,
		CollisionsAtDecision: collisions,
		CreatedAt:            now,
	}
}

// buildDecision produces the E10 row for an approving transition. A transition
// that names no rule writes no decision: the trail records decisions actually
// taken, and inventing a rule for a caller that did not state one would be the
// fabrication E10 exists to prevent.
//
// An unstated actor is NOT defaulted. It used to become 'operator', which is
// the one wrong answer available — a caller that forgot to say who decided got
// a row asserting a human did. Passed through empty, 0082's actor CHECK turns
// the row down and the caller learns; there is no honest guess to make here.
func (s *ResolutionService) buildDecision(input buildResolvedLinkInput, now time.Time) *domain.ProductLinkDecision {
	if input.DecisionRule == "" {
		return nil
	}
	identity := domain.ListingIdentity{
		InstallationID:      input.InstallationID,
		ProviderItemID:      input.ProviderItemID,
		ProviderVariationID: input.ProviderVariationID,
	}
	return s.newDecisionRow(identity, input.DecisionRule, input.DecisionActor, input.CollisionsAtDecision, now)
}

// decisionRuleForCandidate reads the rule off the anchor the candidate was
// built from. A human approving a single-anchor candidate is exactly what the
// confirmation queue asks for, and the trail must keep saying which anchor
// carried it — a candidate approved on no anchor at all is a manual call.
//
// Only an anchor the generator found UNIQUE may name itself in the trail. A
// collision (one EAN, four produtos) or a conflict (CODPROD and EAN
// disagreeing) reaches the operator precisely because no anchor resolved it;
// recording `exact_ean_unique` for a decision taken over a four-way collision
// would assert a uniqueness nobody established (ADR-17) and would read as
// though an anchor had won (AC-08). The operator's judgement is what carried
// those, and `manual` is the honest name for it.
func decisionRuleForCandidate(candidate domain.LinkCandidate) domain.ProductLinkDecisionRule {
	switch candidate.MatchStatus {
	case domain.LinkCandidateMatchStatusAccept:
		// Corroborated: CODPROD and EAN both named this product. Who pressed
		// the button does not change what carried the decision, and the
		// automatic path is not the only way an ACCEPT candidate gets
		// approved — a listing the operator rejected earlier is blocked from
		// it, batch-approve goes through here, and the AutoApprover is an
		// optional wiring. Filing those as a single-anchor rule would
		// under-claim the evidence exactly as naming an anchor over a
		// collision over-claims it.
		return domain.DecisionRuleConcordantCodprodEAN
	case domain.LinkCandidateMatchStatusConfirm:
	default:
		return domain.DecisionRuleManual
	}
	switch candidate.MatchInput {
	case domain.LinkCandidateMatchInputSellerSKU:
		return domain.DecisionRuleExactCodprodUnique
	case domain.LinkCandidateMatchInputEAN:
		return domain.DecisionRuleExactEANUnique
	default:
		return domain.DecisionRuleManual
	}
}

// errSystemActorNotPermitted refuses a system actor on every operator path:
// approve, manual resolve, reject and undo. AutoApproveCandidate is the one
// door a machine may use, and it does not pass through here.
//
// Two different failures are being closed. Reject and undo write
// rule_matched='manual', and the E10 CHECK admits exactly one rule a system
// actor may take — concordant_codprod_ean — so the row is one the schema turns
// down INSIDE the transition's transaction: the whole call rolled back and the
// caller got a 500 naming nothing. Approve and manual resolve had the opposite
// problem: they hardcode actor='operator', so a machine driving them was not
// refused at all, it was RECORDED AS A HUMAN. The trail asserted a person
// decided when nobody had.
//
// Refusing rather than widening the schema is the operator's ruling (D-121-2
// follow-up): a machine may corroborate, and nothing else. Automation that
// needs to reject or undo goes through a person.
//
// Not covered, and it cannot be from here: UndoBatch hardcodes an operator
// actor of its own (see its call to undoAuditEntry), so this guard is
// unreachable on that path and every batch reversal is still signed by an
// operator who was never named. Fixing it means giving a published operation a
// request body — a seam decision, reported to the hub.
func errSystemActorNotPermitted(actor domain.ActorMetadata) error {
	if actor.ActorType == domain.DecisionActorSystem {
		return errors.New("PRODUCT_LINKS_SYSTEM_ACTOR_NOT_PERMITTED")
	}
	return nil
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
