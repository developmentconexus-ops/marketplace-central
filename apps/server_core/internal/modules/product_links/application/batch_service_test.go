package application

import (
	"context"
	"errors"
	"testing"
	"time"

	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

func TestPreviewBatchReturnsItemizedResultsWithoutPersisting(t *testing.T) {
	t.Parallel()

	productA := 501
	productB := 502
	productC := 503
	now := time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	candidateStore := &stubCandidateStore{
		candidates: []productlinksdomain.LinkCandidate{
			{
				CandidateID:       "cand-ok-1",
				InstallationID:    "inst-1",
				ProviderCode:      "mercado_livre",
				ProviderItemID:    "MLB501",
				InternalProductID: &productA,
			},
			{
				CandidateID:       "cand-ok-2",
				InstallationID:    "inst-1",
				ProviderCode:      "mercado_livre",
				ProviderItemID:    "MLB502",
				InternalProductID: &productB,
			},
			{
				CandidateID:       "cand-resolved",
				InstallationID:    "inst-1",
				ProviderCode:      "mercado_livre",
				ProviderItemID:    "MLB503",
				InternalProductID: &productC,
			},
		},
	}
	workflowStore := &stubWorkflowStore{
		linkFound: true,
		currentLink: productlinksdomain.ProductLink{
			InstallationID:    "inst-1",
			ProviderCode:      "mercado_livre",
			ProviderItemID:    "MLB503",
			State:             productlinksdomain.ProductLinkStateResolved,
			InternalProductID: &productC,
			CreatedAt:         now.Add(-time.Hour),
			UpdatedAt:         now.Add(-time.Hour),
		},
	}

	svc := NewBatchService(BatchServiceConfig{Candidates: candidateStore, Workflows: workflowStore})

	result, err := svc.PreviewBatch(context.Background(), PreviewBatchInput{Approvals: []BatchApprovalInput{
		{CandidateID: "cand-ok-1"},
		{CandidateID: "cand-ok-2"},
		{CandidateID: "cand-resolved"},
	}})
	if err != nil {
		t.Fatalf("PreviewBatch() error = %v", err)
	}
	if len(result.Items) != 3 {
		t.Fatalf("items = %d, want 3", len(result.Items))
	}

	byID := map[string]BatchPreviewItem{}
	for _, item := range result.Items {
		byID[item.CandidateID] = item
	}
	okCount, failedCount := 0, 0
	for _, item := range result.Items {
		switch item.Status {
		case BatchItemStatusOK:
			okCount++
		case BatchItemStatusFailed:
			failedCount++
		default:
			t.Fatalf("item %q has unexpected status %q", item.CandidateID, item.Status)
		}
	}
	if okCount != 2 || failedCount != 1 {
		t.Fatalf("ok/failed = %d/%d, want 2/1", okCount, failedCount)
	}
	if byID["cand-ok-1"].Status != BatchItemStatusOK {
		t.Fatalf("cand-ok-1 status = %s, want OK", byID["cand-ok-1"].Status)
	}
	if byID["cand-ok-2"].Status != BatchItemStatusOK {
		t.Fatalf("cand-ok-2 status = %s, want OK", byID["cand-ok-2"].Status)
	}
	if byID["cand-resolved"].Status != BatchItemStatusFailed {
		t.Fatalf("cand-resolved status = %s, want FAILED", byID["cand-resolved"].Status)
	}
	if byID["cand-resolved"].Cause != "ALREADY_RESOLVED" {
		t.Fatalf("cand-resolved cause = %q, want ALREADY_RESOLVED", byID["cand-resolved"].Cause)
	}

	// Dry-run: zero-persist proof. Preview must not write to either store.
	if len(workflowStore.applied) != 0 {
		t.Fatalf("workflow store applied = %d transitions, want 0 (dry-run must not persist)", len(workflowStore.applied))
	}
	if candidateStore.identities != nil || candidateStore.installationID != "" {
		t.Fatalf("candidate store recorded a write, want none (dry-run must not persist)")
	}
}

func TestPreviewBatchRejectsEmptyApprovals(t *testing.T) {
	t.Parallel()

	svc := NewBatchService(BatchServiceConfig{Candidates: &stubCandidateStore{}, Workflows: &stubWorkflowStore{}})

	_, err := svc.PreviewBatch(context.Background(), PreviewBatchInput{})
	if !errors.Is(err, ErrBatchApprovalsRequired) {
		t.Fatalf("PreviewBatch() error = %v, want ErrBatchApprovalsRequired", err)
	}
	_, err = svc.PreviewBatch(context.Background(), PreviewBatchInput{Approvals: []BatchApprovalInput{}})
	if !errors.Is(err, ErrBatchApprovalsRequired) {
		t.Fatalf("PreviewBatch() with empty slice error = %v, want ErrBatchApprovalsRequired", err)
	}
}

func TestEvaluateApprovalCases(t *testing.T) {
	t.Parallel()

	productID := 900
	resolvedProductID := 901
	now := time.Date(2026, 7, 18, 13, 0, 0, 0, time.UTC)

	candidateStore := &stubCandidateStore{
		candidates: []productlinksdomain.LinkCandidate{
			{
				CandidateID:       "cand-valid",
				InstallationID:    "inst-1",
				ProviderCode:      "mercado_livre",
				ProviderItemID:    "MLB900",
				InternalProductID: &productID,
			},
			{
				CandidateID:    "cand-not-resolvable",
				InstallationID: "inst-1",
				ProviderCode:   "mercado_livre",
				ProviderItemID: "MLB910",
				// InternalProductID intentionally nil.
			},
			{
				CandidateID:       "cand-already-resolved",
				InstallationID:    "inst-1",
				ProviderCode:      "mercado_livre",
				ProviderItemID:    "MLB920",
				InternalProductID: &resolvedProductID,
			},
		},
	}
	workflowStore := &stubWorkflowStore{
		linkFound: true,
		currentLink: productlinksdomain.ProductLink{
			InstallationID:    "inst-1",
			ProviderCode:      "mercado_livre",
			ProviderItemID:    "MLB920",
			State:             productlinksdomain.ProductLinkStateResolved,
			InternalProductID: &resolvedProductID,
			CreatedAt:         now.Add(-time.Hour),
			UpdatedAt:         now.Add(-time.Hour),
		},
	}
	svc := NewBatchService(BatchServiceConfig{Candidates: candidateStore, Workflows: workflowStore})

	cases := []struct {
		name        string
		candidateID string
		wantOK      bool
		wantCause   string
	}{
		{"valid resolvable candidate", "cand-valid", true, ""},
		{"nonexistent candidate", "cand-missing", false, "PRODUCT_LINKS_CANDIDATE_NOT_FOUND"},
		{"candidate without internal product", "cand-not-resolvable", false, "PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE"},
		{"candidate already backing a resolved link", "cand-already-resolved", false, "ALREADY_RESOLVED"},
		{"empty candidate id", "", false, "PRODUCT_LINKS_CANDIDATE_REQUIRED"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ok, cause := svc.evaluateApproval(context.Background(), tc.candidateID)
			if ok != tc.wantOK {
				t.Fatalf("evaluateApproval(%q) ok = %v, want %v", tc.candidateID, ok, tc.wantOK)
			}
			if cause != tc.wantCause {
				t.Fatalf("evaluateApproval(%q) cause = %q, want %q", tc.candidateID, cause, tc.wantCause)
			}
		})
	}

	if len(workflowStore.applied) != 0 {
		t.Fatalf("workflow store applied = %d transitions, want 0 (evaluateApproval must not persist)", len(workflowStore.applied))
	}
}

func TestApplyBatchPartialFailureAppliesValidItemsAndItemizesFailures(t *testing.T) {
	t.Parallel()

	productA := 601
	productB := 602
	now := time.Date(2026, 7, 18, 14, 0, 0, 0, time.UTC)
	candidateStore := &stubCandidateStore{
		candidates: []productlinksdomain.LinkCandidate{
			{
				CandidateID:       "cand-apply-1",
				InstallationID:    "inst-1",
				ProviderCode:      "mercado_livre",
				ProviderItemID:    "MLB601",
				InternalProductID: &productA,
			},
			{
				CandidateID:       "cand-apply-2",
				InstallationID:    "inst-1",
				ProviderCode:      "mercado_livre",
				ProviderItemID:    "MLB602",
				InternalProductID: &productB,
			},
			{
				CandidateID:    "cand-apply-invalid",
				InstallationID: "inst-1",
				ProviderCode:   "mercado_livre",
				ProviderItemID: "MLB603",
				// InternalProductID intentionally nil => NOT_RESOLVABLE.
			},
		},
	}
	workflowStore := &stubWorkflowStore{}
	auditCounter := 0
	resolver := NewResolutionService(ResolutionServiceConfig{
		Candidates: candidateStore,
		Workflows:  workflowStore,
		Now:        func() time.Time { return now },
		NewAuditID: func() string {
			auditCounter++
			return "audit-apply-" + string(rune('0'+auditCounter))
		},
	})

	batchIDCounter := 0
	svc := NewBatchService(BatchServiceConfig{
		Candidates: candidateStore,
		Workflows:  workflowStore,
		Approver:   resolver,
		Now:        func() time.Time { return now },
		NewBatchID: func() string {
			batchIDCounter++
			return "batch-apply-1"
		},
	})

	actor := productlinksdomain.ActorMetadata{ActorType: "operator", ActorID: "user-1", ActorName: "Ana"}
	result, err := svc.ApplyBatch(context.Background(), ApplyBatchInput{
		Approvals: []ApplyApprovalInput{
			{CandidateID: "cand-apply-1"},
			{CandidateID: "cand-apply-2"},
			{CandidateID: "cand-apply-invalid"},
		},
		Actor: actor,
	})
	if err != nil {
		t.Fatalf("ApplyBatch() error = %v", err)
	}
	if result.BatchID != "batch-apply-1" {
		t.Fatalf("batch id = %q, want batch-apply-1", result.BatchID)
	}

	// Partial failure: 2 valid approvals applied, the 1 invalid failed —
	// and the invalid item did NOT prevent the 2 valid ones from applying.
	if len(result.Applied) != 2 {
		t.Fatalf("applied = %d, want 2", len(result.Applied))
	}
	if len(result.Failed) != 1 {
		t.Fatalf("failed = %d, want 1", len(result.Failed))
	}
	appliedIDs := map[string]bool{}
	for _, item := range result.Applied {
		appliedIDs[item.CandidateID] = true
	}
	if !appliedIDs["cand-apply-1"] || !appliedIDs["cand-apply-2"] {
		t.Fatalf("applied items = %#v, want cand-apply-1 and cand-apply-2", result.Applied)
	}
	if result.Failed[0].CandidateID != "cand-apply-invalid" {
		t.Fatalf("failed candidate = %q, want cand-apply-invalid", result.Failed[0].CandidateID)
	}
	if result.Failed[0].Cause != "PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE" {
		t.Fatalf("failed cause = %q, want PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE", result.Failed[0].Cause)
	}

	// The 2 valid candidates were actually resolved: re-evaluating them now
	// reports ALREADY_RESOLVED (queue no longer offers them for approval).
	if ok, cause := svc.evaluateApproval(context.Background(), "cand-apply-1"); ok || cause != "ALREADY_RESOLVED" {
		t.Fatalf("cand-apply-1 post-apply evaluateApproval = (%v, %q), want (false, ALREADY_RESOLVED)", ok, cause)
	}
	if ok, cause := svc.evaluateApproval(context.Background(), "cand-apply-2"); ok || cause != "ALREADY_RESOLVED" {
		t.Fatalf("cand-apply-2 post-apply evaluateApproval = (%v, %q), want (false, ALREADY_RESOLVED)", ok, cause)
	}

	// Exactly one product_link_batches row, with requested/applied/failed
	// counts 3/2/1 and a status reflecting the partial outcome.
	if len(workflowStore.batches) != 1 {
		t.Fatalf("batches persisted = %d, want 1", len(workflowStore.batches))
	}
	batch := workflowStore.batches[0]
	if batch.BatchID != "batch-apply-1" {
		t.Fatalf("batch.BatchID = %q, want batch-apply-1", batch.BatchID)
	}
	if batch.RequestedCount != 3 || batch.AppliedCount != 2 || batch.FailedCount != 1 {
		t.Fatalf("batch counts = %d/%d/%d, want 3/2/1", batch.RequestedCount, batch.AppliedCount, batch.FailedCount)
	}
	if batch.Status != productlinksdomain.ProductLinkBatchStatusPartial {
		t.Fatalf("batch.Status = %q, want PARTIAL", batch.Status)
	}
	if batch.Actor.ActorType != "operator" {
		t.Fatalf("batch.Actor.ActorType = %q, want operator", batch.Actor.ActorType)
	}
	if batch.InstallationID != "inst-1" {
		t.Fatalf("batch.InstallationID = %q, want inst-1 (never empty, ADR-17)", batch.InstallationID)
	}

	// Each applied item's audit entry carries the batch_id.
	if len(workflowStore.applied) != 2 {
		t.Fatalf("applied transitions = %d, want 2", len(workflowStore.applied))
	}
	for _, transition := range workflowStore.applied {
		if transition.Audit.BatchID != "batch-apply-1" {
			t.Fatalf("audit entry for %s has BatchID = %q, want batch-apply-1", transition.Audit.SourceCandidateID, transition.Audit.BatchID)
		}
	}
}

// TestApplyBatchInstallationIDFromFirstValidItemNotFirstLookup is the ADR-17
// guard: the FIRST approval is invalid (candidate not found — a raw lookup
// would yield ""), but a later approval is valid. The persisted batch row's
// installation_id must come from the VALID item's resolution, never "" and
// never the failed lookup.
func TestApplyBatchInstallationIDFromFirstValidItemNotFirstLookup(t *testing.T) {
	t.Parallel()

	productValid := 611
	now := time.Date(2026, 7, 18, 15, 0, 0, 0, time.UTC)
	candidateStore := &stubCandidateStore{
		candidates: []productlinksdomain.LinkCandidate{
			{
				CandidateID:       "cand-valid-later",
				InstallationID:    "inst-real",
				ProviderCode:      "mercado_livre",
				ProviderItemID:    "MLB611",
				InternalProductID: &productValid,
			},
			// "cand-missing-first" intentionally absent from the store.
		},
	}
	workflowStore := &stubWorkflowStore{}
	resolver := NewResolutionService(ResolutionServiceConfig{
		Candidates: candidateStore,
		Workflows:  workflowStore,
		Now:        func() time.Time { return now },
		NewAuditID: func() string { return "audit-adr17-1" },
	})
	svc := NewBatchService(BatchServiceConfig{
		Candidates: candidateStore,
		Workflows:  workflowStore,
		Approver:   resolver,
		Now:        func() time.Time { return now },
		NewBatchID: func() string { return "batch-adr17-1" },
	})

	result, err := svc.ApplyBatch(context.Background(), ApplyBatchInput{
		Approvals: []ApplyApprovalInput{
			{CandidateID: "cand-missing-first"}, // invalid, looked up first
			{CandidateID: "cand-valid-later"},   // valid, provides installation_id
		},
		Actor: productlinksdomain.ActorMetadata{ActorType: "operator", ActorID: "user-1"},
	})
	if err != nil {
		t.Fatalf("ApplyBatch() error = %v", err)
	}
	if len(result.Applied) != 1 || len(result.Failed) != 1 {
		t.Fatalf("applied/failed = %d/%d, want 1/1", len(result.Applied), len(result.Failed))
	}
	if len(workflowStore.batches) != 1 {
		t.Fatalf("batches persisted = %d, want 1", len(workflowStore.batches))
	}
	if got := workflowStore.batches[0].InstallationID; got != "inst-real" {
		t.Fatalf("batch.InstallationID = %q, want inst-real (from the valid item, never empty)", got)
	}
}

// TestApplyBatchAllInvalidPersistsNoEmptyInstallationRow is the ADR-17
// all-failed guard: when NO item resolves, no batch row is persisted (which
// would otherwise carry an empty installation_id). The response is still a
// coherent itemized result (empty batch_id, all failures listed).
func TestApplyBatchAllInvalidPersistsNoEmptyInstallationRow(t *testing.T) {
	t.Parallel()

	candidateStore := &stubCandidateStore{} // no candidates => every lookup fails
	workflowStore := &stubWorkflowStore{}
	resolver := NewResolutionService(ResolutionServiceConfig{Candidates: candidateStore, Workflows: workflowStore})
	svc := NewBatchService(BatchServiceConfig{
		Candidates: candidateStore,
		Workflows:  workflowStore,
		Approver:   resolver,
		NewBatchID: func() string { return "batch-allfail-1" },
	})

	result, err := svc.ApplyBatch(context.Background(), ApplyBatchInput{
		Approvals: []ApplyApprovalInput{
			{CandidateID: "cand-x"},
			{CandidateID: "cand-y"},
		},
		Actor: productlinksdomain.ActorMetadata{ActorType: "operator", ActorID: "user-1"},
	})
	if err != nil {
		t.Fatalf("ApplyBatch() error = %v, want nil (all-failed is a coherent 200)", err)
	}
	if len(result.Applied) != 0 || len(result.Failed) != 2 {
		t.Fatalf("applied/failed = %d/%d, want 0/2", len(result.Applied), len(result.Failed))
	}
	if result.BatchID != "" {
		t.Fatalf("result.BatchID = %q, want empty (no batch row persisted)", result.BatchID)
	}
	if len(workflowStore.batches) != 0 {
		t.Fatalf("batches persisted = %d, want 0 (must not write an empty-installation_id row, ADR-17)", len(workflowStore.batches))
	}
	if len(workflowStore.applied) != 0 {
		t.Fatalf("applied transitions = %d, want 0", len(workflowStore.applied))
	}
}

func TestApplyBatchRejectsEmptyApprovals(t *testing.T) {
	t.Parallel()

	workflowStore := &stubWorkflowStore{}
	candidateStore := &stubCandidateStore{}
	resolver := NewResolutionService(ResolutionServiceConfig{Candidates: candidateStore, Workflows: workflowStore})
	svc := NewBatchService(BatchServiceConfig{Candidates: candidateStore, Workflows: workflowStore, Approver: resolver})

	_, err := svc.ApplyBatch(context.Background(), ApplyBatchInput{})
	if !errors.Is(err, ErrBatchApprovalsRequired) {
		t.Fatalf("ApplyBatch() error = %v, want ErrBatchApprovalsRequired", err)
	}
}
