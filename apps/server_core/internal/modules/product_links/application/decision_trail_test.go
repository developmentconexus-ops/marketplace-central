package application

import (
	"context"
	"fmt"
	"testing"
	"time"

	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

func decisionServiceWith(t *testing.T, candidates []productlinksdomain.LinkCandidate, now time.Time) (*ResolutionService, *stubWorkflowStore) {
	t.Helper()
	workflowStore := &stubWorkflowStore{}
	issued := 0
	return NewResolutionService(ResolutionServiceConfig{
		Candidates: &stubCandidateStore{candidates: candidates},
		Workflows:  workflowStore,
		Now:        func() time.Time { return now },
		NewAuditID: func() string { issued++; return fmt.Sprintf("audit-%d", issued) },
		NewDecisionID: func() string {
			return fmt.Sprintf("decision-%d", issued)
		},
	}), workflowStore
}

func singleAnchorCandidate(id string, anchor productlinksdomain.LinkCandidateMatchInput, productID *int) productlinksdomain.LinkCandidate {
	state := productlinksdomain.LinkCandidateStateExactEAN
	if anchor == productlinksdomain.LinkCandidateMatchInputSellerSKU {
		state = productlinksdomain.LinkCandidateStateExactSKU
	}
	return productlinksdomain.LinkCandidate{
		CandidateID:         id,
		InstallationID:      "inst-1",
		ProviderCode:        "mercado_livre",
		ProviderItemID:      "MLB123",
		InternalProductID:   productID,
		InternalProductName: "PUXADOR DHARMA PONTO 10MM CR",
		State:               state,
		MatchInput:          anchor,
		MatchValue:          "100001",
	}
}

func operator(id string) productlinksdomain.ActorMetadata {
	return productlinksdomain.ActorMetadata{ActorType: "operator", ActorID: id, ActorName: "Ana"}
}

// M05-C22: the operator confirming a single-anchor candidate is a human
// decision on one anchor, and the trail has to keep saying which anchor —
// never actor=system, never the corroborated rule.
func TestOperatorConfirmationRecordsTheAnchorItWasTakenOn(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	productID := 100001
	for name, tc := range map[string]struct {
		anchor productlinksdomain.LinkCandidateMatchInput
		want   productlinksdomain.ProductLinkDecisionRule
	}{
		"CODPROD alone": {productlinksdomain.LinkCandidateMatchInputSellerSKU, productlinksdomain.DecisionRuleExactCodprodUnique},
		"EAN alone":     {productlinksdomain.LinkCandidateMatchInputEAN, productlinksdomain.DecisionRuleExactEANUnique},
	} {
		t.Run(name, func(t *testing.T) {
			svc, store := decisionServiceWith(t, []productlinksdomain.LinkCandidate{singleAnchorCandidate("cand-1", tc.anchor, &productID)}, now)
			if _, err := svc.ApproveCandidate(context.Background(), ApproveCandidateInput{CandidateID: "cand-1", Actor: operator("user-1")}); err != nil {
				t.Fatalf("ApproveCandidate() error = %v", err)
			}
			decisions, err := store.ListDecisionsForLink(context.Background(), productlinksdomain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"})
			if err != nil {
				t.Fatalf("ListDecisionsForLink() error = %v", err)
			}
			if len(decisions) != 1 {
				t.Fatalf("decisions = %d, want 1", len(decisions))
			}
			decision := decisions[0]
			if decision.RuleMatched != tc.want {
				t.Errorf("rule_matched = %q, want %q", decision.RuleMatched, tc.want)
			}
			if decision.Actor != productlinksdomain.DecisionActorOperator {
				t.Errorf("actor = %q, want operator", decision.Actor)
			}
			if decision.CollisionsAtDecision != nil {
				t.Errorf("collisions_at_decision = %d, want unknown: no count was read at confirmation time", *decision.CollisionsAtDecision)
			}
			if decision.LinkID != "inst-1|MLB123|" {
				t.Errorf("link_id = %q", decision.LinkID)
			}
		})
	}
}

// The decision row and the link state are one write. Whatever the store is
// handed, it is handed together — a link resolved without its decision would
// leave E10 unable to answer why.
func TestApprovalHandsTheStoreLinkAuditAndDecisionInOneTransition(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	productID := 100001
	svc, store := decisionServiceWith(t, []productlinksdomain.LinkCandidate{singleAnchorCandidate("cand-1", productlinksdomain.LinkCandidateMatchInputEAN, &productID)}, now)
	if _, err := svc.ApproveCandidate(context.Background(), ApproveCandidateInput{CandidateID: "cand-1", Actor: operator("user-1")}); err != nil {
		t.Fatalf("ApproveCandidate() error = %v", err)
	}
	if len(store.applied) != 1 {
		t.Fatalf("applied transitions = %d, want 1", len(store.applied))
	}
	transition := store.applied[0]
	if transition.Decision == nil {
		t.Fatalf("approving transition carried no decision")
	}
	if transition.Decision.CreatedAt != transition.Audit.CreatedAt {
		t.Errorf("decision and audit disagree on when the decision was taken: %v vs %v", transition.Decision.CreatedAt, transition.Audit.CreatedAt)
	}
}

func TestManualResolveIsRecordedAsAManualDecision(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	svc, store := decisionServiceWith(t, nil, now)
	if _, err := svc.ManualResolve(context.Background(), ManualResolveInput{
		InstallationID:    "inst-1",
		ProviderCode:      "mercado_livre",
		ProviderItemID:    "MLB123",
		InternalProductID: 100001,
		Actor:             operator("user-1"),
	}); err != nil {
		t.Fatalf("ManualResolve() error = %v", err)
	}
	decisions, _ := store.ListDecisionsForLink(context.Background(), productlinksdomain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"})
	if len(decisions) != 1 || decisions[0].RuleMatched != productlinksdomain.DecisionRuleManual {
		t.Fatalf("decisions = %+v, want a single manual decision", decisions)
	}
	if decisions[0].CollisionsAtDecision != nil {
		t.Errorf("a manual resolution decided no anchor, so no collision count exists to record")
	}
}

// Rejecting a listing approves nothing. E10 records decisions about which
// product a link names; a rejection names none.
func TestRejectingAListingWritesNoDecision(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	svc, store := decisionServiceWith(t, nil, now)
	if _, err := svc.RejectListing(context.Background(), RejectListingInput{
		InstallationID: "inst-1",
		ProviderCode:   "mercado_livre",
		ProviderItemID: "MLB123",
		Actor:          operator("user-1"),
	}); err != nil {
		t.Fatalf("RejectListing() error = %v", err)
	}
	decisions, _ := store.ListDecisionsForLink(context.Background(), productlinksdomain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"})
	if len(decisions) != 0 {
		t.Fatalf("decisions = %+v, want none", decisions)
	}
}

// M05-C10: a later decision stamps the one it replaced instead of rewriting
// it. Two decisions on the same link leave two rows — the superseded one still
// readable, which is the whole point of an append-only trail.
func TestALaterDecisionSupersedesTheEarlierOneWithoutErasingIt(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	productID := 100001
	svc, store := decisionServiceWith(t, []productlinksdomain.LinkCandidate{singleAnchorCandidate("cand-1", productlinksdomain.LinkCandidateMatchInputEAN, &productID)}, now)
	if _, err := svc.ApproveCandidate(context.Background(), ApproveCandidateInput{CandidateID: "cand-1", Actor: operator("user-1")}); err != nil {
		t.Fatalf("ApproveCandidate() error = %v", err)
	}
	if _, err := svc.ManualResolve(context.Background(), ManualResolveInput{
		InstallationID:    "inst-1",
		ProviderCode:      "mercado_livre",
		ProviderItemID:    "MLB123",
		InternalProductID: 100002,
		Actor:             operator("user-2"),
	}); err != nil {
		t.Fatalf("ManualResolve() error = %v", err)
	}
	decisions, _ := store.ListDecisionsForLink(context.Background(), productlinksdomain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"})
	if len(decisions) != 2 {
		t.Fatalf("decisions = %d, want 2 (history is append-only)", len(decisions))
	}
	if decisions[0].SupersededBy != decisions[1].DecisionID {
		t.Errorf("first decision superseded_by = %q, want %q", decisions[0].SupersededBy, decisions[1].DecisionID)
	}
	if decisions[0].RuleMatched != productlinksdomain.DecisionRuleExactEANUnique {
		t.Errorf("the superseded decision was rewritten: rule = %q", decisions[0].RuleMatched)
	}
	if decisions[1].SupersededBy != "" {
		t.Errorf("the decision in force must not be superseded: %q", decisions[1].SupersededBy)
	}
}
