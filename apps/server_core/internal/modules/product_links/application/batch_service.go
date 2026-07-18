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

// ErrBatchApprovalsRequired signals an empty batch-preview/apply request
// (F-01 empty-batch => 422).
var ErrBatchApprovalsRequired = errors.New("PRODUCT_LINKS_BATCH_APPROVALS_REQUIRED")

// ErrBatchResolutionNotConfigured signals ApplyBatch was invoked without an
// injected CandidateApprover or workflow store (mirrors
// ResolutionService's own "not configured" sentinel).
var ErrBatchResolutionNotConfigured = errors.New("PRODUCT_LINKS_RESOLUTION_NOT_CONFIGURED")

const (
	BatchItemStatusOK     = "OK"
	BatchItemStatusFailed = "FAILED"
)

type BatchApprovalInput struct {
	CandidateID string
}

type PreviewBatchInput struct {
	Approvals []BatchApprovalInput
}

// BatchPreviewItem is the itemized outcome for one approval in a batch
// dry-run: OK when the candidate is resolvable as-is, FAILED (with Cause)
// otherwise. A per-item FAILED is never surfaced as an HTTP error.
type BatchPreviewItem struct {
	CandidateID string `json:"candidate_id"`
	Status      string `json:"status"`
	Cause       string `json:"cause,omitempty"`
}

type PreviewBatchResult struct {
	Items []BatchPreviewItem `json:"items"`
}

// CandidateApprover executes the SAME resolution transition as
// ResolutionService.ApproveCandidate for a single candidate. ApplyBatch (S3)
// reuses it per item instead of reimplementing the approval business rules
// — *ResolutionService satisfies this interface as-is.
type CandidateApprover interface {
	ApproveCandidate(ctx context.Context, input ApproveCandidateInput) (domain.ProductLinkResolutionResult, error)
}

// BatchService implements the batch-preview (S2) and batch-apply (S3)
// surface for product-links resolutions. It shares evaluateApproval across
// both so the dry-run and the apply path agree on what counts as a
// resolvable candidate.
type BatchService struct {
	candidates ports.LinkCandidateStore
	workflows  ports.ProductLinkWorkflowStore
	approver   CandidateApprover
	now        func() time.Time
	newBatchID func() string
}

type BatchServiceConfig struct {
	Candidates ports.LinkCandidateStore
	Workflows  ports.ProductLinkWorkflowStore
	// Approver is the per-item resolver ApplyBatch reuses (S3). Optional
	// for PreviewBatch-only wiring (S2); required for ApplyBatch.
	Approver CandidateApprover
	Now      func() time.Time
	// NewBatchID generates the batch_id for ApplyBatch. Injected (like
	// ResolutionService.NewAuditID) so the pure batch-assembly logic never
	// reaches for time/random itself.
	NewBatchID func() string
}

func NewBatchService(cfg BatchServiceConfig) *BatchService {
	now := cfg.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	newBatchID := cfg.NewBatchID
	if newBatchID == nil {
		newBatchID = func() string { return fmt.Sprintf("plb_%d", now().UnixNano()) }
	}
	return &BatchService{
		candidates: cfg.Candidates,
		workflows:  cfg.Workflows,
		approver:   cfg.Approver,
		now:        now,
		newBatchID: newBatchID,
	}
}

// PreviewBatch evaluates every approval in the batch and returns an itemized
// OK/FAILED result. It is READ-ONLY by construction: it only ever calls
// GetLinkCandidate / GetProductLink (reads) and never ApplyProductLinkTransition
// or any other store write.
func (s *BatchService) PreviewBatch(ctx context.Context, input PreviewBatchInput) (PreviewBatchResult, error) {
	if len(input.Approvals) == 0 {
		return PreviewBatchResult{}, ErrBatchApprovalsRequired
	}
	items := make([]BatchPreviewItem, 0, len(input.Approvals))
	for _, approval := range input.Approvals {
		candidateID := strings.TrimSpace(approval.CandidateID)
		item := BatchPreviewItem{CandidateID: candidateID}
		if ok, cause := s.evaluateApproval(ctx, candidateID); ok {
			item.Status = BatchItemStatusOK
		} else {
			item.Status = BatchItemStatusFailed
			item.Cause = cause
		}
		items = append(items, item)
	}
	return PreviewBatchResult{Items: items}, nil
}

// ApplyApprovalInput is one item of an ApplyBatch request.
type ApplyApprovalInput struct {
	CandidateID string
}

// ApplyBatchInput is the batch-apply (S3) request: the approvals to resolve
// plus the acting operator (threaded into each item's audit entry).
type ApplyBatchInput struct {
	Approvals []ApplyApprovalInput
	Actor     domain.ActorMetadata
}

// AppliedBatchItem identifies a candidate that was successfully resolved as
// part of the batch.
type AppliedBatchItem struct {
	CandidateID string `json:"candidate_id"`
}

// FailedBatchItem identifies a candidate that could not be resolved, with an
// honest cause (ADR-17) — a per-item failure that never aborts the batch.
type FailedBatchItem struct {
	CandidateID string `json:"candidate_id"`
	Cause       string `json:"cause"`
}

// ApplyBatchResult is the itemized outcome of a batch-apply run, plus the
// batch_id of the persisted product_link_batches row.
type ApplyBatchResult struct {
	BatchID string             `json:"batch_id"`
	Applied []AppliedBatchItem `json:"applied"`
	Failed  []FailedBatchItem  `json:"failed"`
}

// ApplyBatch resolves every approval in the batch. Per-item failure
// (evaluateApproval says not-ok, or the underlying ApproveCandidate call
// errors) is itemized in Failed and never rolls back or aborts the other
// items — only each item's own resolution transition is atomic, not the
// batch as a whole. When at least one item resolves, ONE product_link_batches
// row is persisted with the aggregate counts/status; every applied item's
// audit entry is tagged with the same batch_id via CandidateApprover
// (ResolutionService reuse — no duplicated business rules, no provider/ML
// write: apply calls only the local resolver).
//
// ADR-17 (unknown facts never become silent zero/default): the batch's
// installation_id is taken from the resolution result of the FIRST item that
// actually resolves — never from a raw pre-check lookup of a possibly
// failed/not-found candidate. If NO item resolves (all failed), no batch row
// is persisted (BatchID stays empty in the response) rather than writing a
// row with an empty installation_id. The response is still a coherent 200
// with itemized failures.
func (s *BatchService) ApplyBatch(ctx context.Context, input ApplyBatchInput) (ApplyBatchResult, error) {
	if len(input.Approvals) == 0 {
		return ApplyBatchResult{}, ErrBatchApprovalsRequired
	}
	if s.approver == nil || s.workflows == nil || s.candidates == nil {
		return ApplyBatchResult{}, ErrBatchResolutionNotConfigured
	}
	batchID := s.newBatchID()
	applied := make([]AppliedBatchItem, 0, len(input.Approvals))
	failed := make([]FailedBatchItem, 0)
	installationID := ""
	for _, approval := range input.Approvals {
		candidateID := strings.TrimSpace(approval.CandidateID)
		if ok, cause := s.evaluateApproval(ctx, candidateID); !ok {
			failed = append(failed, FailedBatchItem{CandidateID: candidateID, Cause: cause})
			continue
		}
		result, err := s.approver.ApproveCandidate(ctx, ApproveCandidateInput{
			CandidateID: candidateID,
			Actor:       input.Actor,
			BatchID:     batchID,
		})
		if err != nil {
			failed = append(failed, FailedBatchItem{CandidateID: candidateID, Cause: err.Error()})
			continue
		}
		// installation_id from the FIRST definitively-resolved item's link
		// (ADR-17: an honest, verified installation context — never derived
		// from a failed/not-found candidate lookup).
		if installationID == "" {
			installationID = result.Link.InstallationID
		}
		applied = append(applied, AppliedBatchItem{CandidateID: candidateID})
	}

	// All items failed: do NOT persist a batch row with an empty
	// installation_id (ADR-17). Return the itemized failures with no
	// batch_id; handleBatchApply still responds 200.
	if len(applied) == 0 {
		return ApplyBatchResult{Applied: applied, Failed: failed}, nil
	}

	status := domain.ProductLinkBatchStatusCompleted
	if len(failed) > 0 {
		status = domain.ProductLinkBatchStatusPartial
	}
	batch := domain.ProductLinkBatch{
		BatchID:        batchID,
		InstallationID: installationID,
		Actor: domain.ActorMetadata{
			ActorType: "operator",
			ActorID:   input.Actor.ActorID,
			ActorName: input.Actor.ActorName,
		},
		RequestedCount: len(input.Approvals),
		AppliedCount:   len(applied),
		FailedCount:    len(failed),
		Status:         status,
		CreatedAt:      s.now(),
	}
	if err := s.workflows.InsertBatch(ctx, batch); err != nil {
		return ApplyBatchResult{}, err
	}
	return ApplyBatchResult{BatchID: batchID, Applied: applied, Failed: failed}, nil
}

// evaluateApproval is the single per-item eval helper shared by batch-preview
// (S2) and batch-apply (S3): it checks that the candidate exists, carries a
// resolvable internal product, and does not already back a resolved product
// link. It performs reads only (GetLinkCandidate, GetProductLink) — no writes
// — so S3 can call it verbatim as its pre-apply check before mutating.
func (s *BatchService) evaluateApproval(ctx context.Context, candidateID string) (ok bool, cause string) {
	if s.candidates == nil || s.workflows == nil {
		return false, "PRODUCT_LINKS_RESOLUTION_NOT_CONFIGURED"
	}
	if candidateID == "" {
		return false, "PRODUCT_LINKS_CANDIDATE_REQUIRED"
	}
	candidate, found, err := s.candidates.GetLinkCandidate(ctx, candidateID)
	if err != nil {
		return false, err.Error()
	}
	if !found {
		return false, "PRODUCT_LINKS_CANDIDATE_NOT_FOUND"
	}
	if candidate.InternalProductID == nil {
		return false, "PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE"
	}
	if err := domain.ValidateInternalProductID(*candidate.InternalProductID); err != nil {
		return false, err.Error()
	}
	identity := domain.ListingIdentity{
		InstallationID:      candidate.InstallationID,
		ProviderItemID:      candidate.ProviderItemID,
		ProviderVariationID: candidate.ProviderVariationID,
	}
	current, found, err := s.workflows.GetProductLink(ctx, identity)
	if err != nil {
		return false, err.Error()
	}
	if found && current.State == domain.ProductLinkStateResolved {
		return false, "ALREADY_RESOLVED"
	}
	return true, ""
}
