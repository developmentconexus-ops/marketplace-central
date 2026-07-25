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
		// A single UNIQUE anchor is exactly what generation puts in the
		// confirmation queue; the rule may name that anchor only because the
		// generator found it unique.
		MatchStatus: productlinksdomain.LinkCandidateMatchStatusConfirm,
	}
}

// ambiguousCandidate is what the operator sees when NO anchor resolved the
// listing — a colliding EAN or a CODPROD≠EAN conflict. It still carries the
// anchor it was built from, which is precisely why the trail must not read
// that anchor as the rule that decided it.
func ambiguousCandidate(id string, anchor productlinksdomain.LinkCandidateMatchInput, productID *int) productlinksdomain.LinkCandidate {
	candidate := singleAnchorCandidate(id, anchor, productID)
	candidate.State = productlinksdomain.LinkCandidateStateConflict
	candidate.MatchStatus = productlinksdomain.LinkCandidateMatchStatusReview
	return candidate
}

// concordantCandidate is what the generator produces when CODPROD and EAN both
// resolve the same product — the only shape the automatic path accepts.
func concordantCandidate(id string, productID *int) productlinksdomain.LinkCandidate {
	candidate := singleAnchorCandidate(id, productlinksdomain.LinkCandidateMatchInputSellerSKU, productID)
	candidate.MatchStatus = productlinksdomain.LinkCandidateMatchStatusAccept
	return candidate
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

// AC-08/ADR-17: approving a candidate that reached the operator BECAUSE no
// anchor resolved it must not record that anchor as the rule. `exact_ean_unique`
// on the operator's real 4-way collision would assert a uniqueness nobody
// established and would read as though an anchor had won.
func TestApprovingAnAmbiguousCandidateIsRecordedAsAManualCall(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	productID := 100001
	for name, anchor := range map[string]productlinksdomain.LinkCandidateMatchInput{
		"colliding EAN":   productlinksdomain.LinkCandidateMatchInputEAN,
		"conflicting SKU": productlinksdomain.LinkCandidateMatchInputSellerSKU,
	} {
		t.Run(name, func(t *testing.T) {
			svc, store := decisionServiceWith(t, []productlinksdomain.LinkCandidate{ambiguousCandidate("cand-1", anchor, &productID)}, now)
			if _, err := svc.ApproveCandidate(context.Background(), ApproveCandidateInput{CandidateID: "cand-1", Actor: operator("user-1")}); err != nil {
				t.Fatalf("ApproveCandidate() error = %v", err)
			}
			decisions, _ := store.ListDecisionsForLink(context.Background(), productlinksdomain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"})
			if len(decisions) != 1 {
				t.Fatalf("decisions = %d, want 1", len(decisions))
			}
			if decisions[0].RuleMatched != productlinksdomain.DecisionRuleManual {
				t.Errorf("rule_matched = %q, want manual: no anchor resolved this listing", decisions[0].RuleMatched)
			}
			if decisions[0].Actor != productlinksdomain.DecisionActorOperator {
				t.Errorf("actor = %q, want operator", decisions[0].Actor)
			}
		})
	}
}

// An undo is a decision too: it says the link that stood should not. Left out
// of the trail, the decision it reverts keeps reading as in force — and the
// automatic path stays blocked by a row nobody stands behind, with nothing
// saying why.
func TestUndoSupersedesTheDecisionItReverts(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	productID := 100001
	svc, store := decisionServiceWith(t, []productlinksdomain.LinkCandidate{singleAnchorCandidate("cand-1", productlinksdomain.LinkCandidateMatchInputEAN, &productID)}, now)
	result, err := svc.ApproveCandidate(context.Background(), ApproveCandidateInput{CandidateID: "cand-1", Actor: operator("user-1")})
	if err != nil {
		t.Fatalf("ApproveCandidate() error = %v", err)
	}
	if _, err := svc.UndoResolution(context.Background(), UndoResolutionInput{
		AuditID: result.Audit.AuditID,
		Reason:  "vínculo errado",
		Actor:   operator("user-1"),
	}); err != nil {
		t.Fatalf("UndoResolution() error = %v", err)
	}

	decisions, _ := store.ListDecisionsForLink(context.Background(), productlinksdomain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"})
	if len(decisions) != 2 {
		t.Fatalf("decisions = %d, want the undo recorded alongside what it reverted", len(decisions))
	}
	if decisions[0].SupersededBy != decisions[1].DecisionID {
		t.Errorf("the reverted decision still reads as in force: superseded_by = %q", decisions[0].SupersededBy)
	}
	if decisions[1].RuleMatched != productlinksdomain.DecisionRuleManual {
		t.Errorf("undo rule = %q, want manual", decisions[1].RuleMatched)
	}
	if decisions[1].Actor != productlinksdomain.DecisionActorOperator {
		t.Errorf("undo actor = %q, want operator", decisions[1].Actor)
	}
	if decisions[1].CollisionsAtDecision != nil {
		t.Errorf("collisions_at_decision = %d, want unknown: an undo reads no collision count", *decisions[1].CollisionsAtDecision)
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

// A rejection names no anchor, but it settles the link — so it must supersede
// the decision it overrules. Without its own row the automatic approval stays
// the one in force, and the trail reports the system as still standing behind
// a link the operator killed.
func TestRejectingAListingSupersedesTheDecisionItOverrules(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	productID := 100002
	svc, store := decisionServiceWith(t, []productlinksdomain.LinkCandidate{concordantCandidate("cand-1", &productID)}, now)
	collisions := 1
	if _, err := svc.AutoApproveCandidate(context.Background(), AutoApproveCandidateInput{
		Candidate:            concordantCandidate("cand-1", &productID),
		CollisionsAtDecision: &collisions,
	}); err != nil {
		t.Fatalf("AutoApproveCandidate() error = %v", err)
	}
	identity := productlinksdomain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"}
	if _, err := svc.RejectListing(context.Background(), RejectListingInput{
		InstallationID: "inst-1",
		ProviderCode:   "mercado_livre",
		ProviderItemID: "MLB123",
		Actor:          operator("user-1"),
	}); err != nil {
		t.Fatalf("RejectListing() error = %v", err)
	}

	decisions, _ := store.ListDecisionsForLink(context.Background(), identity)
	if len(decisions) != 2 {
		t.Fatalf("decisions = %+v, want the rejection recorded alongside what it overruled", decisions)
	}
	var live []productlinksdomain.ProductLinkDecision
	for _, decision := range decisions {
		if decision.SupersededBy == "" {
			live = append(live, decision)
		}
	}
	if len(live) != 1 {
		t.Fatalf("live decisions = %+v, want exactly one in force", live)
	}
	if live[0].Actor != productlinksdomain.DecisionActorOperator || live[0].RuleMatched != productlinksdomain.DecisionRuleManual {
		t.Errorf("decision in force = %+v, want the operator's manual rejection, not the system's approval", live[0])
	}
	if live[0].CollisionsAtDecision != nil {
		t.Errorf("a rejection read no anchor, so it has no collision count to record")
	}
}

// A corroborated candidate approved BY HAND was still carried by both anchors.
// The automatic path is not the only route an ACCEPT candidate takes — a
// listing the operator rejected earlier is locked out of it, batch-approve
// comes through here, and the AutoApprover is optional wiring — so the trail
// must record what actually decided it, not who pressed the button.
func TestOperatorApprovingACorroboratedCandidateRecordsCorroboration(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	productID := 100002
	svc, store := decisionServiceWith(t, []productlinksdomain.LinkCandidate{concordantCandidate("cand-1", &productID)}, now)
	if _, err := svc.ApproveCandidate(context.Background(), ApproveCandidateInput{CandidateID: "cand-1", Actor: operator("user-1")}); err != nil {
		t.Fatalf("ApproveCandidate() error = %v", err)
	}
	decisions, _ := store.ListDecisionsForLink(context.Background(), productlinksdomain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"})
	if len(decisions) != 1 {
		t.Fatalf("decisions = %+v, want one", decisions)
	}
	if decisions[0].RuleMatched != productlinksdomain.DecisionRuleConcordantCodprodEAN {
		t.Errorf("rule_matched = %q, want the corroborated rule: both anchors named this product", decisions[0].RuleMatched)
	}
	if decisions[0].Actor != productlinksdomain.DecisionActorOperator {
		t.Errorf("actor = %q, want operator: a human took this one", decisions[0].Actor)
	}
}

// AC-03: a caller that read no collision count must be able to say so. The
// count is the generator's reading, and an approval taken without one stores
// NULL — never 0, which would read as "the anchor matched nothing".
func TestAnAutomaticApprovalWithoutACountRecordsNoCount(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	productID := 100002
	svc, store := decisionServiceWith(t, nil, now)
	if _, err := svc.AutoApproveCandidate(context.Background(), AutoApproveCandidateInput{
		Candidate: concordantCandidate("cand-1", &productID),
	}); err != nil {
		t.Fatalf("AutoApproveCandidate() error = %v", err)
	}
	decisions, _ := store.ListDecisionsForLink(context.Background(), productlinksdomain.ListingIdentity{InstallationID: "inst-1", ProviderItemID: "MLB123"})
	if len(decisions) != 1 {
		t.Fatalf("decisions = %+v, want one", decisions)
	}
	if decisions[0].CollisionsAtDecision != nil {
		t.Fatalf("collisions_at_decision = %d, want NULL: nobody read a count here", *decisions[0].CollisionsAtDecision)
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
