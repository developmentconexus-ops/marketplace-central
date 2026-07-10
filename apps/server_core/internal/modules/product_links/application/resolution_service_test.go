package application

import (
	"context"
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
}

func (s *stubWorkflowStore) GetProductLink(_ context.Context, identity productlinksdomain.ListingIdentity) (productlinksdomain.ProductLink, bool, error) {
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
			AuditID:           "audit-4",
			InstallationID:    "inst-1",
			ProviderCode:      "mercado_livre",
			ProviderItemID:    "MLB404",
			Action:            productlinksdomain.ProductLinkActionApproveCandidate,
			PreviousState:     productlinksdomain.ProductLinkStateConflict,
			NextState:         productlinksdomain.ProductLinkStateResolved,
			NextInternalProductID: &productID,
			CreatedAt:         now.Add(-20 * time.Minute),
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
