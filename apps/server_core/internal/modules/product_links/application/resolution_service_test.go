package application

import (
	"context"
	"errors"
	"testing"
	"time"

	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type stubWorkflowStore struct {
	linkFound      bool
	currentLink    productlinksdomain.ProductLink
	links          []productlinksdomain.ProductLink
	audits         []productlinksdomain.ProductLinkAuditEntry
	applied        []productlinksdomain.ProductLinkTransition
	auditIDCounter int
	batches        []productlinksdomain.ProductLinkBatch
	// linksByIdentity tracks every applied link by identity (rather than
	// just the single pre-seeded currentLink), so batch-apply tests with
	// multiple distinct identities can assert each one individually
	// (e.g. re-evaluating a just-applied candidate now reports
	// ALREADY_RESOLVED, without clobbering other identities' state).
	linksByIdentity map[string]productlinksdomain.ProductLink
}

func (s *stubWorkflowStore) GetProductLink(_ context.Context, identity productlinksdomain.ListingIdentity) (productlinksdomain.ProductLink, bool, error) {
	key := identityKey(identity.InstallationID, identity.ProviderItemID, identity.ProviderVariationID)
	if link, ok := s.linksByIdentity[key]; ok {
		return link, true, nil
	}
	if s.linkFound &&
		s.currentLink.InstallationID == identity.InstallationID &&
		s.currentLink.ProviderItemID == identity.ProviderItemID &&
		s.currentLink.ProviderVariationID == identity.ProviderVariationID {
		return s.currentLink, true, nil
	}
	return productlinksdomain.ProductLink{}, false, nil
}

func (s *stubWorkflowStore) ApplyProductLinkTransition(_ context.Context, transition productlinksdomain.ProductLinkTransition) error {
	s.applied = append(s.applied, transition)
	s.currentLink = transition.Link
	s.linkFound = true
	s.links = []productlinksdomain.ProductLink{transition.Link}
	s.audits = append([]productlinksdomain.ProductLinkAuditEntry{transition.Audit}, s.audits...)
	if s.linksByIdentity == nil {
		s.linksByIdentity = map[string]productlinksdomain.ProductLink{}
	}
	key := identityKey(transition.Link.InstallationID, transition.Link.ProviderItemID, transition.Link.ProviderVariationID)
	s.linksByIdentity[key] = transition.Link
	return nil
}

func (s *stubWorkflowStore) ListProductLinks(_ context.Context, installationID string, _ int) ([]productlinksdomain.ProductLink, error) {
	items := make([]productlinksdomain.ProductLink, 0, len(s.links))
	for _, item := range s.links {
		if item.InstallationID == installationID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (s *stubWorkflowStore) ListProductLinkAuditEntries(_ context.Context, installationID string, _ int) ([]productlinksdomain.ProductLinkAuditEntry, error) {
	items := make([]productlinksdomain.ProductLinkAuditEntry, 0, len(s.audits))
	for _, item := range s.audits {
		if item.InstallationID == installationID {
			items = append(items, item)
		}
	}
	return items, nil
}

// InsertBatch is the S3 batch-audit write. Kept on the shared stub (rather
// than only in batch_service_test.go) because widening
// ports.ProductLinkWorkflowStore with this method means every stub
// implementing the interface — including this resolution_service_test.go
// fixture used by NewResolutionService — must keep compiling.
func (s *stubWorkflowStore) InsertBatch(_ context.Context, batch productlinksdomain.ProductLinkBatch) error {
	s.batches = append(s.batches, batch)
	return nil
}

// GetAuditEntry / LatestAuditForIdentity / ListAuditByBatch (S4 undo reads)
// are additive on ports.ProductLinkWorkflowStore. This stub's s.audits slice
// is already newest-first (ApplyProductLinkTransition prepends), so a
// front-to-back scan naturally yields the latest match for any filter.

func (s *stubWorkflowStore) GetAuditEntry(_ context.Context, auditID string) (productlinksdomain.ProductLinkAuditEntry, bool, error) {
	for _, audit := range s.audits {
		if audit.AuditID == auditID {
			return audit, true, nil
		}
	}
	return productlinksdomain.ProductLinkAuditEntry{}, false, nil
}

func (s *stubWorkflowStore) LatestAuditForIdentity(_ context.Context, identity productlinksdomain.ListingIdentity) (productlinksdomain.ProductLinkAuditEntry, bool, error) {
	for _, audit := range s.audits {
		if audit.InstallationID == identity.InstallationID &&
			audit.ProviderItemID == identity.ProviderItemID &&
			audit.ProviderVariationID == identity.ProviderVariationID {
			return audit, true, nil
		}
	}
	return productlinksdomain.ProductLinkAuditEntry{}, false, nil
}

func (s *stubWorkflowStore) ListAuditByBatch(_ context.Context, batchID string) ([]productlinksdomain.ProductLinkAuditEntry, error) {
	items := make([]productlinksdomain.ProductLinkAuditEntry, 0)
	for _, audit := range s.audits {
		if audit.BatchID == batchID {
			items = append(items, audit)
		}
	}
	return items, nil
}

func TestApproveCandidateCreatesResolvedLinkAndAudit(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 21, 0, 0, 0, time.UTC)
	productID := 101
	candidateStore := &stubCandidateStore{
		candidates: []productlinksdomain.LinkCandidate{{
			CandidateID:           "cand-1",
			InstallationID:        "inst-1",
			ProviderCode:          "mercado_livre",
			ProviderItemID:        "MLB123",
			InternalProductID:     &productID,
			InternalProductName:   "Produto Interno",
			InternalReferenceCode: "SKU-101",
			State:                 productlinksdomain.LinkCandidateStateExactEAN,
			MatchInput:            productlinksdomain.LinkCandidateMatchInputEAN,
			MatchValue:            "789",
			CreatedAt:             now.Add(-time.Hour),
			UpdatedAt:             now.Add(-time.Hour),
		}},
	}
	workflowStore := &stubWorkflowStore{}
	svc := NewResolutionService(ResolutionServiceConfig{
		Candidates: candidateStore,
		Workflows:  workflowStore,
		Now:        func() time.Time { return now },
		NewAuditID: func() string { return "audit-1" },
	})

	result, err := svc.ApproveCandidate(context.Background(), ApproveCandidateInput{
		CandidateID: "cand-1",
		Actor: productlinksdomain.ActorMetadata{
			ActorType: "operator",
			ActorID:   "user-1",
			ActorName: "Ana",
		},
		Reason: "EAN confirmed by operator",
	})
	if err != nil {
		t.Fatalf("ApproveCandidate() error = %v", err)
	}

	if len(workflowStore.applied) != 1 {
		t.Fatalf("applied transitions = %d, want 1", len(workflowStore.applied))
	}
	if result.Link.State != productlinksdomain.ProductLinkStateResolved {
		t.Fatalf("link state = %s, want resolved", result.Link.State)
	}
	if result.Link.SourceCandidateID != "cand-1" {
		t.Fatalf("source candidate = %q, want cand-1", result.Link.SourceCandidateID)
	}
	if result.Link.InternalProductID == nil || *result.Link.InternalProductID != 101 {
		t.Fatalf("internal product id = %#v, want 101", result.Link.InternalProductID)
	}
	if result.Audit.Action != productlinksdomain.ProductLinkActionApproveCandidate {
		t.Fatalf("audit action = %s, want approve_candidate", result.Audit.Action)
	}
	if result.Audit.PreviousState != productlinksdomain.ProductLinkStateNone || result.Audit.NextState != productlinksdomain.ProductLinkStateResolved {
		t.Fatalf("audit states = %s -> %s", result.Audit.PreviousState, result.Audit.NextState)
	}
}

func TestRejectListingReplacesExistingLinkWithRejectedStateAndAudit(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 21, 5, 0, 0, time.UTC)
	productID := 202
	workflowStore := &stubWorkflowStore{
		linkFound: true,
		currentLink: productlinksdomain.ProductLink{
			InstallationID:      "inst-1",
			ProviderCode:        "mercado_livre",
			ProviderItemID:      "MLB999",
			State:               productlinksdomain.ProductLinkStateResolved,
			InternalProductID:   &productID,
			InternalProductName: "Produto Atual",
			CreatedAt:           now.Add(-2 * time.Hour),
			UpdatedAt:           now.Add(-time.Hour),
		},
	}
	svc := NewResolutionService(ResolutionServiceConfig{
		Candidates: &stubCandidateStore{},
		Workflows:  workflowStore,
		Now:        func() time.Time { return now },
		NewAuditID: func() string { return "audit-2" },
	})

	result, err := svc.RejectListing(context.Background(), RejectListingInput{
		InstallationID: "inst-1",
		ProviderCode:   "mercado_livre",
		ProviderItemID: "MLB999",
		Actor: productlinksdomain.ActorMetadata{
			ActorType: "operator",
			ActorID:   "user-2",
		},
		Reason: "Listing does not belong to our catalog",
	})
	if err != nil {
		t.Fatalf("RejectListing() error = %v", err)
	}

	if result.Link.State != productlinksdomain.ProductLinkStateRejected {
		t.Fatalf("link state = %s, want rejected", result.Link.State)
	}
	if result.Link.InternalProductID != nil {
		t.Fatalf("internal product id = %#v, want nil", result.Link.InternalProductID)
	}
	if result.Audit.PreviousState != productlinksdomain.ProductLinkStateResolved || result.Audit.NextState != productlinksdomain.ProductLinkStateRejected {
		t.Fatalf("audit states = %s -> %s", result.Audit.PreviousState, result.Audit.NextState)
	}
	if result.Audit.PreviousInternalProductID == nil || *result.Audit.PreviousInternalProductID != 202 {
		t.Fatalf("previous internal product id = %#v, want 202", result.Audit.PreviousInternalProductID)
	}
}

func TestManualResolveCreatesResolvedLinkAndAuditWithoutCandidateSource(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 21, 10, 0, 0, time.UTC)
	svc := NewResolutionService(ResolutionServiceConfig{
		Candidates: &stubCandidateStore{},
		Workflows:  &stubWorkflowStore{},
		Now:        func() time.Time { return now },
		NewAuditID: func() string { return "audit-3" },
	})

	result, err := svc.ManualResolve(context.Background(), ManualResolveInput{
		InstallationID:        "inst-1",
		ProviderCode:          "mercado_livre",
		ProviderItemID:        "MLB777",
		ProviderVariationID:   "VAR-1",
		InternalProductID:     303,
		InternalProductName:   "Produto Manual",
		InternalReferenceCode: "REF-303",
		Actor: productlinksdomain.ActorMetadata{
			ActorType: "operator",
			ActorID:   "user-3",
			ActorName: "Bia",
		},
		Reason: "Variation confirmed manually",
	})
	if err != nil {
		t.Fatalf("ManualResolve() error = %v", err)
	}

	if result.Link.State != productlinksdomain.ProductLinkStateResolved {
		t.Fatalf("link state = %s, want resolved", result.Link.State)
	}
	if result.Link.SourceCandidateID != "" {
		t.Fatalf("source candidate = %q, want empty", result.Link.SourceCandidateID)
	}
	if result.Link.InternalProductID == nil || *result.Link.InternalProductID != 303 {
		t.Fatalf("internal product id = %#v, want 303", result.Link.InternalProductID)
	}
	if result.Audit.Action != productlinksdomain.ProductLinkActionManualResolve {
		t.Fatalf("audit action = %s, want manual_resolve", result.Audit.Action)
	}
}

func TestManualResolveRejectsNonPositiveInternalProductID(t *testing.T) {
	svc := NewResolutionService(ResolutionServiceConfig{Workflows: &stubWorkflowStore{}})
	for _, id := range []int{-1, 0} {
		_, err := svc.ManualResolve(context.Background(), ManualResolveInput{InternalProductID: id})
		if !errors.Is(err, productlinksdomain.ErrInvalidInternalProductID) {
			t.Fatalf("ManualResolve(%d) error = %v", id, err)
		}
	}
}

func TestListLinkWorkflowsMergesCandidatesWithCurrentLinkAndAudit(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 8, 21, 15, 0, 0, time.UTC)
	productID := 404
	candidateStore := &stubCandidateStore{
		candidates: []productlinksdomain.LinkCandidate{{
			CandidateID:         "cand-4",
			InstallationID:      "inst-1",
			ProviderCode:        "mercado_livre",
			ProviderItemID:      "MLB404",
			InternalProductID:   &productID,
			InternalProductName: "Produto Workflow",
			State:               productlinksdomain.LinkCandidateStateConflict,
			MatchInput:          productlinksdomain.LinkCandidateMatchInputSellerSKU,
			MatchValue:          "SKU-404",
			CreatedAt:           now.Add(-time.Hour),
			UpdatedAt:           now.Add(-time.Hour),
		}},
	}
	workflowStore := &stubWorkflowStore{
		links: []productlinksdomain.ProductLink{{
			InstallationID:    "inst-1",
			ProviderCode:      "mercado_livre",
			ProviderItemID:    "MLB404",
			State:             productlinksdomain.ProductLinkStateResolved,
			InternalProductID: &productID,
			CreatedAt:         now.Add(-30 * time.Minute),
			UpdatedAt:         now.Add(-30 * time.Minute),
		}},
		audits: []productlinksdomain.ProductLinkAuditEntry{{
			AuditID:               "audit-4",
			InstallationID:        "inst-1",
			ProviderCode:          "mercado_livre",
			ProviderItemID:        "MLB404",
			Action:                productlinksdomain.ProductLinkActionApproveCandidate,
			PreviousState:         productlinksdomain.ProductLinkStateConflict,
			NextState:             productlinksdomain.ProductLinkStateResolved,
			NextInternalProductID: &productID,
			CreatedAt:             now.Add(-20 * time.Minute),
		}},
	}
	svc := NewResolutionService(ResolutionServiceConfig{
		Candidates: candidateStore,
		Workflows:  workflowStore,
		Now:        func() time.Time { return now },
		NewAuditID: func() string { return "unused" },
	})

	items, err := svc.ListLinkWorkflows(context.Background(), ListLinkWorkflowsInput{
		InstallationID: "inst-1",
		Limit:          10,
	})
	if err != nil {
		t.Fatalf("ListLinkWorkflows() error = %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("items = %d, want 1", len(items))
	}
	if items[0].CurrentLink == nil || items[0].CurrentLink.State != productlinksdomain.ProductLinkStateResolved {
		t.Fatalf("current link = %#v", items[0].CurrentLink)
	}
	if len(items[0].Candidates) != 1 || items[0].Candidates[0].CandidateID != "cand-4" {
		t.Fatalf("candidates = %#v", items[0].Candidates)
	}
	if len(items[0].Audit) != 1 || items[0].Audit[0].AuditID != "audit-4" {
		t.Fatalf("audit = %#v", items[0].Audit)
	}
}

// --- S4 undo -----------------------------------------------------------

func TestUndoResolutionRevertsCandidateToQueue(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 18, 16, 0, 0, 0, time.UTC)
	productID := 701
	candidateStore := &stubCandidateStore{
		candidates: []productlinksdomain.LinkCandidate{{
			CandidateID:        "cand-undo-1",
			InstallationID:     "inst-1",
			ProviderCode:       "mercado_livre",
			ProviderItemID:     "MLB701",
			InternalProductID:  &productID,
			State:              productlinksdomain.LinkCandidateStateUnresolved,
			CreatedAt:          now.Add(-time.Hour),
			UpdatedAt:          now.Add(-time.Hour),
		}},
	}
	workflowStore := &stubWorkflowStore{}
	auditSeq := 0
	svc := NewResolutionService(ResolutionServiceConfig{
		Candidates: candidateStore,
		Workflows:  workflowStore,
		Now:        func() time.Time { return now },
		NewAuditID: func() string {
			auditSeq++
			return "audit-undo-" + string(rune('0'+auditSeq))
		},
	})

	approveResult, err := svc.ApproveCandidate(context.Background(), ApproveCandidateInput{
		CandidateID: "cand-undo-1",
		Actor:       productlinksdomain.ActorMetadata{ActorType: "operator", ActorID: "user-1"},
		Reason:      "approved for undo test",
	})
	if err != nil {
		t.Fatalf("ApproveCandidate() error = %v", err)
	}
	if approveResult.Link.State != productlinksdomain.ProductLinkStateResolved {
		t.Fatalf("precondition: link state = %s, want resolved", approveResult.Link.State)
	}

	undoResult, err := svc.UndoResolution(context.Background(), UndoResolutionInput{
		AuditID: approveResult.Audit.AuditID,
		Actor:   productlinksdomain.ActorMetadata{ActorType: "operator", ActorID: "user-1"},
		Reason:  "operator changed their mind",
	})
	if err != nil {
		t.Fatalf("UndoResolution() error = %v", err)
	}

	// The candidate is back in the queue: state reverted to its
	// pre-approval value (unresolved), internal product cleared.
	if undoResult.Link.State != productlinksdomain.ProductLinkStateUnresolved {
		t.Fatalf("link state after undo = %s, want unresolved", undoResult.Link.State)
	}
	if undoResult.Link.InternalProductID != nil {
		t.Fatalf("internal product id after undo = %#v, want nil", undoResult.Link.InternalProductID)
	}

	// Reversal is a NEW audit row (action=undo), the original row is
	// untouched — full history preserved, never a delete/rewrite.
	if undoResult.Audit.Action != productlinksdomain.ProductLinkActionUndo {
		t.Fatalf("undo audit action = %s, want undo", undoResult.Audit.Action)
	}
	if undoResult.Audit.PreviousState != productlinksdomain.ProductLinkStateResolved {
		t.Fatalf("undo audit previous state = %s, want resolved", undoResult.Audit.PreviousState)
	}
	if undoResult.Audit.NextState != productlinksdomain.ProductLinkStateUnresolved {
		t.Fatalf("undo audit next state = %s, want unresolved", undoResult.Audit.NextState)
	}
	if len(workflowStore.applied) != 2 {
		t.Fatalf("applied transitions = %d, want 2 (approve + undo, no rewrite)", len(workflowStore.applied))
	}
	originalStillPresent := false
	for _, audit := range workflowStore.audits {
		if audit.AuditID == approveResult.Audit.AuditID && audit.Action == productlinksdomain.ProductLinkActionApproveCandidate {
			originalStillPresent = true
		}
	}
	if !originalStillPresent {
		t.Fatalf("original approve audit row missing after undo — history must be preserved, not deleted")
	}
}

func TestUndoResolutionSecondUndoIsAlreadyUndone(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 18, 16, 5, 0, 0, time.UTC)
	productID := 702
	candidateStore := &stubCandidateStore{
		candidates: []productlinksdomain.LinkCandidate{{
			CandidateID:       "cand-undo-2",
			InstallationID:    "inst-1",
			ProviderCode:      "mercado_livre",
			ProviderItemID:    "MLB702",
			InternalProductID: &productID,
			State:             productlinksdomain.LinkCandidateStateUnresolved,
		}},
	}
	workflowStore := &stubWorkflowStore{}
	auditSeq := 0
	svc := NewResolutionService(ResolutionServiceConfig{
		Candidates: candidateStore,
		Workflows:  workflowStore,
		Now:        func() time.Time { return now },
		NewAuditID: func() string {
			auditSeq++
			return "audit-undo2-" + string(rune('0'+auditSeq))
		},
	})

	approveResult, err := svc.ApproveCandidate(context.Background(), ApproveCandidateInput{CandidateID: "cand-undo-2"})
	if err != nil {
		t.Fatalf("ApproveCandidate() error = %v", err)
	}

	if _, err := svc.UndoResolution(context.Background(), UndoResolutionInput{AuditID: approveResult.Audit.AuditID}); err != nil {
		t.Fatalf("first UndoResolution() error = %v", err)
	}

	_, err = svc.UndoResolution(context.Background(), UndoResolutionInput{AuditID: approveResult.Audit.AuditID})
	if !errors.Is(err, ErrProductLinkAlreadyUndone) {
		t.Fatalf("second UndoResolution() error = %v, want ErrProductLinkAlreadyUndone", err)
	}
}

func TestUndoResolutionSupersededByNewerResolution(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 7, 18, 16, 10, 0, 0, time.UTC)
	productID := 703
	candidateStore := &stubCandidateStore{
		candidates: []productlinksdomain.LinkCandidate{{
			CandidateID:       "cand-undo-3",
			InstallationID:    "inst-1",
			ProviderCode:      "mercado_livre",
			ProviderItemID:    "MLB703",
			InternalProductID: &productID,
			State:             productlinksdomain.LinkCandidateStateUnresolved,
		}},
	}
	workflowStore := &stubWorkflowStore{}
	auditSeq := 0
	svc := NewResolutionService(ResolutionServiceConfig{
		Candidates: candidateStore,
		Workflows:  workflowStore,
		Now:        func() time.Time { return now },
		NewAuditID: func() string {
			auditSeq++
			return "audit-undo3-" + string(rune('0'+auditSeq))
		},
	})

	firstApprove, err := svc.ApproveCandidate(context.Background(), ApproveCandidateInput{CandidateID: "cand-undo-3"})
	if err != nil {
		t.Fatalf("first ApproveCandidate() error = %v", err)
	}

	// A newer resolution supersedes the first: a manual resolve on the same
	// identity, e.g. the operator re-resolved the listing to a different
	// product.
	if _, err := svc.ManualResolve(context.Background(), ManualResolveInput{
		InstallationID:      "inst-1",
		ProviderCode:        "mercado_livre",
		ProviderItemID:      "MLB703",
		InternalProductID:   999,
		InternalProductName: "Produto Substituto",
	}); err != nil {
		t.Fatalf("ManualResolve() error = %v", err)
	}

	_, err = svc.UndoResolution(context.Background(), UndoResolutionInput{AuditID: firstApprove.Audit.AuditID})
	if !errors.Is(err, ErrProductLinkSuperseded) {
		t.Fatalf("UndoResolution() of superseded entry error = %v, want ErrProductLinkSuperseded", err)
	}
}

func TestUndoResolutionMissingAuditIsNotFound(t *testing.T) {
	t.Parallel()

	svc := NewResolutionService(ResolutionServiceConfig{
		Candidates: &stubCandidateStore{},
		Workflows:  &stubWorkflowStore{},
	})

	_, err := svc.UndoResolution(context.Background(), UndoResolutionInput{AuditID: "does-not-exist"})
	if !errors.Is(err, ErrProductLinkAuditNotFound) {
		t.Fatalf("UndoResolution() error = %v, want ErrProductLinkAuditNotFound", err)
	}
}

func TestUndoBatchRevertsAllItemsAndItemizesAlreadyUndone(t *testing.T) {
	t.Parallel()

	productA := 801
	productB := 802
	now := time.Date(2026, 7, 18, 16, 15, 0, 0, time.UTC)
	candidateStore := &stubCandidateStore{
		candidates: []productlinksdomain.LinkCandidate{
			{
				CandidateID:       "cand-batch-undo-1",
				InstallationID:    "inst-1",
				ProviderCode:      "mercado_livre",
				ProviderItemID:    "MLB801",
				InternalProductID: &productA,
				State:             productlinksdomain.LinkCandidateStateUnresolved,
			},
			{
				CandidateID:       "cand-batch-undo-2",
				InstallationID:    "inst-1",
				ProviderCode:      "mercado_livre",
				ProviderItemID:    "MLB802",
				InternalProductID: &productB,
				State:             productlinksdomain.LinkCandidateStateUnresolved,
			},
		},
	}
	workflowStore := &stubWorkflowStore{}
	auditSeq := 0
	resolver := NewResolutionService(ResolutionServiceConfig{
		Candidates: candidateStore,
		Workflows:  workflowStore,
		Now:        func() time.Time { return now },
		NewAuditID: func() string {
			auditSeq++
			return "audit-batch-undo-" + string(rune('0'+auditSeq))
		},
	})
	batchSvc := NewBatchService(BatchServiceConfig{
		Candidates: candidateStore,
		Workflows:  workflowStore,
		Approver:   resolver,
		Now:        func() time.Time { return now },
		NewBatchID: func() string { return "batch-undo-1" },
	})

	applyResult, err := batchSvc.ApplyBatch(context.Background(), ApplyBatchInput{
		Approvals: []ApplyApprovalInput{
			{CandidateID: "cand-batch-undo-1"},
			{CandidateID: "cand-batch-undo-2"},
		},
		Actor: productlinksdomain.ActorMetadata{ActorType: "operator", ActorID: "user-1"},
	})
	if err != nil {
		t.Fatalf("ApplyBatch() error = %v", err)
	}
	if len(applyResult.Applied) != 2 {
		t.Fatalf("applied = %d, want 2", len(applyResult.Applied))
	}

	// Find each item's own audit entry (tagged with the batch id) so we can
	// pre-undo one of them individually before the batch-undo runs.
	var undoneAheadOfTime productlinksdomain.ProductLinkAuditEntry
	found := false
	for _, audit := range workflowStore.audits {
		if audit.BatchID == applyResult.BatchID && audit.SourceCandidateID == "cand-batch-undo-1" {
			undoneAheadOfTime = audit
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("could not find batch audit row for cand-batch-undo-1")
	}
	if _, err := resolver.UndoResolution(context.Background(), UndoResolutionInput{AuditID: undoneAheadOfTime.AuditID}); err != nil {
		t.Fatalf("pre-undo of cand-batch-undo-1 error = %v", err)
	}

	result, err := resolver.UndoBatch(context.Background(), applyResult.BatchID)
	if err != nil {
		t.Fatalf("UndoBatch() error = %v", err)
	}
	if result.BatchID != applyResult.BatchID {
		t.Fatalf("UndoBatch() batch id = %q, want %q", result.BatchID, applyResult.BatchID)
	}
	// cand-batch-undo-1 was already undone ahead of time: itemized as
	// failed, not fatal to the rest of the batch-undo.
	if len(result.Failed) != 1 || result.Failed[0].CandidateID != "cand-batch-undo-1" {
		t.Fatalf("failed = %#v, want 1 item for cand-batch-undo-1", result.Failed)
	}
	if result.Failed[0].Cause != "ALREADY_UNDONE" {
		t.Fatalf("failed cause = %q, want ALREADY_UNDONE", result.Failed[0].Cause)
	}
	// cand-batch-undo-2 was reverted cleanly by the batch-undo call.
	if len(result.Reverted) != 1 || result.Reverted[0].CandidateID != "cand-batch-undo-2" {
		t.Fatalf("reverted = %#v, want 1 item for cand-batch-undo-2", result.Reverted)
	}

	link2, found, err := workflowStore.GetProductLink(context.Background(), productlinksdomain.ListingIdentity{
		InstallationID: "inst-1",
		ProviderItemID: "MLB802",
	})
	if err != nil || !found {
		t.Fatalf("GetProductLink(MLB802) error = %v, found = %v", err, found)
	}
	if link2.State != productlinksdomain.ProductLinkStateUnresolved {
		t.Fatalf("cand-batch-undo-2 link state after batch-undo = %s, want unresolved", link2.State)
	}
}
