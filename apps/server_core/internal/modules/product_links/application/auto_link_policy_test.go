package application

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

// The auto-linking policy ratified as ADR-05 amended (D-121-2): the ONLY
// automatic path is CODPROD and EAN resolving the same product. One anchor
// alone is a confirmation the operator gives; ambiguity, conflict and
// hard-negatives are review. These tests drive generation with the resolution
// service actually wired in, so what is asserted is the link and the E10 trail
// the pair really wrote — not a mock's opinion of them.

const (
	// Golden data from the dev-stack mirror (source=sankhya).
	goldenSKUConcordant = "100002"
	goldenEANConcordant = "7909251304727"
	goldenNameDharma192 = "PUXADOR DHARMA 192MM PINT TIT"
	goldenSKUOnly       = "100001"
	goldenEANOnly       = "7909251260214"
	goldenEANColliding  = "7896902180697"
)

func autoLinkingServices(t *testing.T, snapshots []productlinksdomain.ListingSnapshot, matcher *stubProductMatcher, now time.Time) (*GenerationService, *stubWorkflowStore, *stubCandidateStore) {
	t.Helper()
	candidateStore := &stubCandidateStore{}
	workflowStore := &stubWorkflowStore{}
	issued := 0
	resolution := NewResolutionService(ResolutionServiceConfig{
		Candidates:    candidateStore,
		Workflows:     workflowStore,
		Now:           func() time.Time { return now },
		NewAuditID:    func() string { issued++; return fmt.Sprintf("audit-%d", issued) },
		NewDecisionID: func() string { return fmt.Sprintf("decision-%d", issued) },
	})
	generation := NewGenerationService(GenerationServiceConfig{
		Snapshots:    &stubSnapshotReader{snapshots: snapshots},
		Matcher:      matcher,
		Store:        candidateStore,
		AutoApprover: resolution,
		Now:          func() time.Time { return now },
	})
	return generation, workflowStore, candidateStore
}

func concordantSnapshot(itemID string, now time.Time) productlinksdomain.ListingSnapshot {
	return productlinksdomain.ListingSnapshot{
		InstallationID: "inst-m05", ProviderCode: "mercado_livre", ProviderItemID: itemID,
		SellerSKU: goldenSKUConcordant, EAN: goldenEANConcordant, Title: goldenNameDharma192, FetchedAt: now,
	}
}

func concordantMatcher() *stubProductMatcher {
	product := internalreaddomain.ProductCandidate{InternalProductID: canonicalIDPtr(100002), Name: goldenNameDharma192}
	return &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:" + goldenSKUConcordant: {product},
		"ean:" + goldenEANConcordant: {product},
	}}
}

func generateAll(t *testing.T, svc *GenerationService, installationID string) []productlinksdomain.LinkCandidate {
	t.Helper()
	result, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: installationID})
	if err != nil {
		t.Fatalf("GenerateLinkCandidates() error = %v", err)
	}
	return result.Items
}

func decisionsFor(t *testing.T, store *stubWorkflowStore, itemID string) []productlinksdomain.ProductLinkDecision {
	t.Helper()
	decisions, err := store.ListDecisionsForLink(context.Background(), productlinksdomain.ListingIdentity{
		InstallationID: "inst-m05", ProviderItemID: itemID,
	})
	if err != nil {
		t.Fatalf("ListDecisionsForLink() error = %v", err)
	}
	return decisions
}

// M05-C15 + M05-C3: the corroborated pair is the one automatic path, and it
// leaves an E10 row saying so — rule, actor, and the collision count the
// generator itself read.
func TestConcordantAnchorsAreTheOnlyAutomaticPath(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	svc, workflows, _ := autoLinkingServices(t, []productlinksdomain.ListingSnapshot{concordantSnapshot("MLB-C", now)}, concordantMatcher(), now)

	candidates := generateAll(t, svc, "inst-m05")
	if len(candidates) != 1 || candidates[0].MatchStatus != productlinksdomain.LinkCandidateMatchStatusAccept {
		t.Fatalf("candidates=%#v, want a single ACCEPT", candidates)
	}
	if len(workflows.applied) != 1 {
		t.Fatalf("applied transitions = %d, want the listing auto-approved", len(workflows.applied))
	}
	link := workflows.applied[0].Link
	if link.State != productlinksdomain.ProductLinkStateResolved || link.InternalProductID == nil || *link.InternalProductID != 100002 {
		t.Fatalf("link = %#v, want resolved onto codprod 100002", link)
	}
	decisions := decisionsFor(t, workflows, "MLB-C")
	if len(decisions) != 1 {
		t.Fatalf("decisions = %#v, want exactly one", decisions)
	}
	decision := decisions[0]
	if decision.RuleMatched != productlinksdomain.DecisionRuleConcordantCodprodEAN {
		t.Errorf("rule_matched = %q, want concordant_codprod_ean", decision.RuleMatched)
	}
	if decision.Actor != productlinksdomain.DecisionActorSystem {
		t.Errorf("actor = %q, want system", decision.Actor)
	}
	if decision.CollisionsAtDecision == nil || *decision.CollisionsAtDecision != 1 {
		t.Errorf("collisions_at_decision = %v, want the generator's count of 1", decision.CollisionsAtDecision)
	}
}

// M05-C2 and M05-C14: a single anchor resolves one product, but nothing
// corroborates it — so it waits for a human, and the warning names the anchor
// that is missing (AC-11 forbids a generic one).
func TestSingleAnchorGoesToConfirmationNeverAutoApproved(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	for name, tc := range map[string]struct {
		snapshot productlinksdomain.ListingSnapshot
		matcher  *stubProductMatcher
		anchor   string
		warning  string
	}{
		"CODPROD sem EAN": {
			snapshot: productlinksdomain.ListingSnapshot{
				InstallationID: "inst-m05", ProviderCode: "mercado_livre", ProviderItemID: "MLB-B",
				SellerSKU: goldenSKUOnly, Title: "PUXADOR DHARMA PONTO CR", FetchedAt: now,
			},
			matcher: &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
				"sku:" + goldenSKUOnly: {{InternalProductID: canonicalIDPtr(100001), Name: "PUXADOR DHARMA PONTO CR"}},
			}},
			anchor:  "ean",
			warning: "sem EAN para corroborar o CODPROD",
		},
		"EAN sem CODPROD": {
			snapshot: productlinksdomain.ListingSnapshot{
				InstallationID: "inst-m05", ProviderCode: "mercado_livre", ProviderItemID: "MLB-A",
				EAN: goldenEANOnly, Title: "PUXADOR DHARMA PONTO CR", FetchedAt: now,
			},
			matcher: &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
				"ean:" + goldenEANOnly: {{InternalProductID: canonicalIDPtr(100000), Name: "PUXADOR DHARMA PONTO CR"}},
			}},
			anchor:  "seller_sku",
			warning: "sem CODPROD para corroborar o EAN",
		},
	} {
		t.Run(name, func(t *testing.T) {
			svc, workflows, _ := autoLinkingServices(t, []productlinksdomain.ListingSnapshot{tc.snapshot}, tc.matcher, now)
			candidates := generateAll(t, svc, "inst-m05")
			if len(candidates) != 1 {
				t.Fatalf("candidates=%#v, want one", candidates)
			}
			if candidates[0].MatchStatus != productlinksdomain.LinkCandidateMatchStatusConfirm {
				t.Fatalf("match_status = %q, want CONFIRM", candidates[0].MatchStatus)
			}
			reason, ok := findReason(candidates[0].Reasons, tc.anchor, productlinksdomain.LinkCandidateReasonDirectionUnavailable)
			if !ok || reason.Detail != tc.warning {
				t.Fatalf("reasons=%#v, want %q on the %s anchor", candidates[0].Reasons, tc.warning, tc.anchor)
			}
			if len(workflows.applied) != 0 {
				t.Fatalf("applied = %#v, want no link approved without a human", workflows.applied)
			}
			if got := decisionsFor(t, workflows, tc.snapshot.ProviderItemID); len(got) != 0 {
				t.Fatalf("decisions = %#v, want none until the operator confirms", got)
			}
		})
	}
}

// M05-C4: an EAN matching four products (the operator's real 7896902180697)
// identified none of them. Every product is offered for review and nothing is
// approved.
func TestCollidingAnchorGoesToReviewAndApprovesNothing(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	colliding := make([]internalreaddomain.ProductCandidate, 0, 4)
	for _, id := range []int{42535, 42536, 42537, 42538} {
		colliding = append(colliding, internalreaddomain.ProductCandidate{InternalProductID: canonicalIDPtr(id), Name: fmt.Sprintf("PARAFUSO %d", id)})
	}
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-m05", ProviderCode: "mercado_livre", ProviderItemID: "MLB-D",
		EAN: goldenEANColliding, Title: "PARAFUSO", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"ean:" + goldenEANColliding: colliding,
	}}
	svc, workflows, _ := autoLinkingServices(t, []productlinksdomain.ListingSnapshot{snapshot}, matcher, now)

	candidates := generateAll(t, svc, "inst-m05")
	if len(candidates) != 4 {
		t.Fatalf("candidates = %d, want all four colliding products offered", len(candidates))
	}
	for _, candidate := range candidates {
		if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusReview {
			t.Fatalf("candidate %#v, want REVIEW", candidate)
		}
	}
	if len(workflows.applied) != 0 {
		t.Fatalf("applied = %#v, want nothing approved off an ambiguous anchor", workflows.applied)
	}
}

// M05-C16 / AC-08: the anchors point at different products. Neither wins — the
// operator decides, and nothing is written automatically.
func TestConflictingAnchorsApproveNothingAndElectNoWinner(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	snapshot := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-m05", ProviderCode: "mercado_livre", ProviderItemID: "MLB-E",
		SellerSKU: "100000", EAN: "7899656858195", Title: "PUXADOR", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:100000":        {{InternalProductID: canonicalIDPtr(100000), Name: "PUXADOR DHARMA PONTO CR"}},
		"ean:7899656858195": {{InternalProductID: canonicalIDPtr(100003), Name: "PUXADOR OUTRO"}},
	}}
	svc, workflows, _ := autoLinkingServices(t, []productlinksdomain.ListingSnapshot{snapshot}, matcher, now)

	candidates := generateAll(t, svc, "inst-m05")
	if len(candidates) != 2 {
		t.Fatalf("candidates = %#v, want both sides of the conflict", candidates)
	}
	seen := map[int]bool{}
	for _, candidate := range candidates {
		if candidate.MatchStatus != productlinksdomain.LinkCandidateMatchStatusReview {
			t.Fatalf("candidate %#v, want REVIEW", candidate)
		}
		seen[*candidate.InternalProductID] = true
	}
	if !seen[100000] || !seen[100003] {
		t.Fatalf("candidates = %#v, want both products present (no anchor discarded)", candidates)
	}
	if len(workflows.applied) != 0 {
		t.Fatalf("applied = %#v, want no automatic resolution of a conflict", workflows.applied)
	}
}

// M05-C17: a title contradiction beats agreeing anchors. Must-fail proof:
// dropping the detectHardNegative check in buildConcordantCandidate makes this
// listing auto-approve.
func TestHardNegativeBlocksTheAutomaticPathEvenWithBothAnchors(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	snapshot := concordantSnapshot("MLB-KIT", now)
	snapshot.Title = "KIT 2 UN PUXADOR DHARMA 192MM PINT TIT"
	svc, workflows, _ := autoLinkingServices(t, []productlinksdomain.ListingSnapshot{snapshot}, concordantMatcher(), now)

	candidates := generateAll(t, svc, "inst-m05")
	if len(candidates) != 1 {
		t.Fatalf("candidates = %#v, want one", candidates)
	}
	if candidates[0].MatchStatus == productlinksdomain.LinkCandidateMatchStatusAccept {
		t.Fatalf("match_status = ACCEPT, want the kit/unit contradiction to block the automatic path")
	}
	if _, ok := findReason(candidates[0].Reasons, "title", productlinksdomain.LinkCandidateReasonDirectionAgainst); !ok {
		t.Fatalf("reasons = %#v, want the contradiction stated", candidates[0].Reasons)
	}
	if len(workflows.applied) != 0 {
		t.Fatalf("applied = %#v, want nothing approved over a hard negative", workflows.applied)
	}
}

// M05-C21: confirmation and review are separate queues. A count by status
// returns two groups, and the confirmation warning survives the split —
// mapping CONFIRM onto REVIEW collapses this to one group and fails here.
func TestConfirmationAndReviewAreCountedSeparately(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	confirmable := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-m05", ProviderCode: "mercado_livre", ProviderItemID: "MLB-B",
		SellerSKU: goldenSKUOnly, Title: "PUXADOR DHARMA PONTO CR", FetchedAt: now,
	}
	reviewable := productlinksdomain.ListingSnapshot{
		InstallationID: "inst-m05", ProviderCode: "mercado_livre", ProviderItemID: "MLB-D",
		EAN: goldenEANColliding, Title: "PARAFUSO", FetchedAt: now,
	}
	matcher := &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{
		"sku:" + goldenSKUOnly: {{InternalProductID: canonicalIDPtr(100001), Name: "PUXADOR DHARMA PONTO CR"}},
		"ean:" + goldenEANColliding: {
			{InternalProductID: canonicalIDPtr(42535), Name: "PARAFUSO 42535"},
			{InternalProductID: canonicalIDPtr(42536), Name: "PARAFUSO 42536"},
		},
	}}
	svc, _, _ := autoLinkingServices(t, []productlinksdomain.ListingSnapshot{confirmable, reviewable}, matcher, now)

	byStatus := map[productlinksdomain.LinkCandidateMatchStatus][]productlinksdomain.LinkCandidate{}
	for _, candidate := range generateAll(t, svc, "inst-m05") {
		byStatus[candidate.MatchStatus] = append(byStatus[candidate.MatchStatus], candidate)
	}
	if len(byStatus) != 2 {
		t.Fatalf("statuses = %#v, want confirmation and review as separate groups", byStatus)
	}
	confirmations := byStatus[productlinksdomain.LinkCandidateMatchStatusConfirm]
	reviews := byStatus[productlinksdomain.LinkCandidateMatchStatusReview]
	if len(confirmations) != 1 || len(reviews) != 2 {
		t.Fatalf("confirmations=%d reviews=%d, want 1 and 2", len(confirmations), len(reviews))
	}
	reason, ok := findReason(confirmations[0].Reasons, "ean", productlinksdomain.LinkCandidateReasonDirectionUnavailable)
	if !ok || reason.Detail != "sem EAN para corroborar o CODPROD" {
		t.Fatalf("confirmation reasons=%#v, want the missing-anchor warning preserved", confirmations[0].Reasons)
	}
}

// M05-C20 / AC-09: the same product legitimately sells under several listings.
// Two listings resolving codprod 100002 produce two links, neither flagged.
func TestOneProductAutoLinksToManyListings(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	first := concordantSnapshot("MLB-C1", now)
	second := concordantSnapshot("MLB-C2", now)
	svc, workflows, _ := autoLinkingServices(t, []productlinksdomain.ListingSnapshot{first, second}, concordantMatcher(), now)

	generateAll(t, svc, "inst-m05")
	if len(workflows.applied) != 2 {
		t.Fatalf("applied = %d, want one link per listing", len(workflows.applied))
	}
	for _, transition := range workflows.applied {
		if transition.Link.State != productlinksdomain.ProductLinkStateResolved {
			t.Fatalf("link %#v, want resolved — sharing a product is not a defect", transition.Link)
		}
	}
}

// M05-C8: re-running generation over the same snapshot (a re-import or the
// next sync) must not pile up links or decisions.
func TestRerunningGenerationDoesNotDuplicateTheAutomaticLink(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	svc, workflows, _ := autoLinkingServices(t, []productlinksdomain.ListingSnapshot{concordantSnapshot("MLB-C", now)}, concordantMatcher(), now)

	generateAll(t, svc, "inst-m05")
	generateAll(t, svc, "inst-m05")

	if len(workflows.applied) != 1 {
		t.Fatalf("applied = %d, want the second run to be a no-op", len(workflows.applied))
	}
	if got := decisionsFor(t, workflows, "MLB-C"); len(got) != 1 {
		t.Fatalf("decisions = %#v, want one — a re-sync decides nothing new", got)
	}
}

// M05-C10: the operator's answer is final. A later automatic run finds a
// decision in force and leaves it alone instead of reverting to actor=system.
func TestAutoApprovalNeverOverridesAnOperatorDecision(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	svc, workflows, _ := autoLinkingServices(t, []productlinksdomain.ListingSnapshot{concordantSnapshot("MLB-C", now)}, concordantMatcher(), now)

	// The operator resolved this listing onto a different product first.
	workflows.decisions = append(workflows.decisions, productlinksdomain.ProductLinkDecision{
		DecisionID:     "decision-operator",
		InstallationID: "inst-m05",
		ProviderItemID: "MLB-C",
		LinkID:         productlinksdomain.LinkID("inst-m05", "MLB-C", ""),
		RuleMatched:    productlinksdomain.DecisionRuleManual,
		Actor:          productlinksdomain.DecisionActorOperator,
		CreatedAt:      now.Add(-time.Hour),
	})

	generateAll(t, svc, "inst-m05")

	if len(workflows.applied) != 0 {
		t.Fatalf("applied = %#v, want the operator's decision untouched", workflows.applied)
	}
	decisions := decisionsFor(t, workflows, "MLB-C")
	if len(decisions) != 1 || decisions[0].Actor != productlinksdomain.DecisionActorOperator {
		t.Fatalf("decisions = %#v, want only the operator's, unsuperseded", decisions)
	}
}

// M05-C10, second half: "não é meu anúncio" is an operator answer too, and
// rejection writes no E10 row — there is no rule to record for it. The link's
// own terminal state is what keeps the automatic path from quietly reopening a
// listing the operator closed on the next sync.
func TestAutoApprovalNeverReopensAListingTheOperatorRejected(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	svc, workflows, _ := autoLinkingServices(t, []productlinksdomain.ListingSnapshot{concordantSnapshot("MLB-C", now)}, concordantMatcher(), now)

	rejected := productlinksdomain.ProductLink{
		InstallationID: "inst-m05",
		ProviderCode:   "mercado_livre",
		ProviderItemID: "MLB-C",
		State:          productlinksdomain.ProductLinkStateRejected,
		CreatedAt:      now.Add(-time.Hour),
		UpdatedAt:      now.Add(-time.Hour),
	}
	workflows.links = []productlinksdomain.ProductLink{rejected}
	workflows.currentLink = rejected
	workflows.linkFound = true

	generateAll(t, svc, "inst-m05")

	if len(workflows.applied) != 0 {
		t.Fatalf("applied = %#v, want the operator's rejection left standing", workflows.applied)
	}
	if got := decisionsFor(t, workflows, "MLB-C"); len(got) != 0 {
		t.Fatalf("decisions = %#v, want none — the automatic path decided nothing", got)
	}
	if workflows.links[0].State != productlinksdomain.ProductLinkStateRejected {
		t.Fatalf("state = %q, want the listing still rejected", workflows.links[0].State)
	}
}

// The automatic path is closed by construction: a caller handing
// AutoApproveCandidate anything the generator did not mark corroborated is
// refused rather than quietly approved.
func TestAutoApproveRefusesACandidateItDidNotCorroborate(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	workflows := &stubWorkflowStore{}
	svc := NewResolutionService(ResolutionServiceConfig{
		Candidates: &stubCandidateStore{},
		Workflows:  workflows,
		Now:        func() time.Time { return now },
	})
	productID := 100001
	oneCollision := 1
	approved, err := svc.AutoApproveCandidate(context.Background(), AutoApproveCandidateInput{
		Candidate: productlinksdomain.LinkCandidate{
			InstallationID: "inst-m05", ProviderCode: "mercado_livre", ProviderItemID: "MLB-B",
			InternalProductID: &productID,
			State:             productlinksdomain.LinkCandidateStateExactSKU,
			MatchInput:        productlinksdomain.LinkCandidateMatchInputSellerSKU,
			MatchStatus:       productlinksdomain.LinkCandidateMatchStatusConfirm,
		},
		CollisionsAtDecision: &oneCollision,
	})
	if err == nil || approved {
		t.Fatalf("AutoApproveCandidate() = %v, %v; want a refusal for a confirmation-queue candidate", approved, err)
	}
	if len(workflows.applied) != 0 {
		t.Fatalf("applied = %#v, want nothing written", workflows.applied)
	}
}

// failingApprover records every listing it was offered and fails the ones it
// was told to, so a batch that abandons its remainder is visible as a listing
// that was never offered at all.
type failingApprover struct {
	failOn  map[string]bool
	offered []string
}

func (f *failingApprover) AutoApproveCandidate(_ context.Context, input AutoApproveCandidateInput) (bool, error) {
	f.offered = append(f.offered, input.Candidate.ProviderItemID)
	if f.failOn[input.Candidate.ProviderItemID] {
		return false, fmt.Errorf("approval refused for %s", input.Candidate.ProviderItemID)
	}
	return true, nil
}

// One listing that cannot be approved is not a reason to leave the rest of the
// run unlinked. Before the batch continued, the first failure returned and
// every corroborated candidate behind it stayed pending with nothing saying
// why — on a 10k-listing sync that is the whole batch lost to one bad row.
func TestAFailedApprovalDoesNotAbandonTheRestOfTheBatch(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	approver := &failingApprover{failOn: map[string]bool{"MLB-FAIL": true, "MLB-ALSO-FAIL": true}}
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: &stubSnapshotReader{snapshots: []productlinksdomain.ListingSnapshot{
			concordantSnapshot("MLB-FAIL", now),
			concordantSnapshot("MLB-ALSO-FAIL", now),
			concordantSnapshot("MLB-OK", now),
		}},
		Matcher:      concordantMatcher(),
		Store:        &stubCandidateStore{},
		AutoApprover: approver,
		Now:          func() time.Time { return now },
	})

	_, err := svc.GenerateLinkCandidates(context.Background(), GenerateLinkCandidatesInput{InstallationID: "inst-m05"})
	if err == nil {
		t.Fatal("GenerateLinkCandidates() error = nil, want the refused approval reported")
	}
	// Two failures, so the join is load-bearing: keeping only the last error
	// would still report a failure and still continue the batch, and a single
	// failing listing could not tell the two apart.
	for _, want := range []string{"MLB-FAIL", "MLB-ALSO-FAIL"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error = %v, want it to name %s", err, want)
		}
	}
	if len(approver.offered) != 3 {
		t.Fatalf("offered = %v, want every listing attempted", approver.offered)
	}
	if approver.offered[2] != "MLB-OK" {
		t.Errorf("offered = %v, want MLB-OK attempted after both failures", approver.offered)
	}
}

// Every operator-facing entry point refuses a machine. AutoApproveCandidate is
// the sole exception and is covered by its own tests — a machine may
// corroborate, and nothing else (D-121-2 follow-up, operator ruling).
//
// The two failure modes being closed are opposites, which is why one table
// covers them rather than two suites. Reject and undo write
// rule_matched='manual', a pair the E10 CHECK turns down INSIDE the
// transition's transaction: the call rolled back with a 500 naming nothing.
// Approve and manual-resolve hardcode actor='operator', so a machine there was
// never refused — it was FILED AS A HUMAN, and the trail asserted a person
// decided when nobody had.
//
// Each case is given input that would ALSO fail a later validation (a blank
// candidate id, a zero product id, a blank item id). A run that reports the
// later error instead proves the guard is not first, which is the whole claim:
// nothing is validated, read or written on behalf of a caller that may not act.
func TestEveryOperatorPathRefusesASystemActor(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 25, 10, 0, 0, 0, time.UTC)
	system := productlinksdomain.ActorMetadata{ActorType: productlinksdomain.DecisionActorSystem, ActorID: "nightly_sweep"}

	newService := func(workflows *stubWorkflowStore) *ResolutionService {
		return NewResolutionService(ResolutionServiceConfig{
			Candidates: &stubCandidateStore{},
			Workflows:  workflows,
			Now:        func() time.Time { return now },
			NewAuditID: func() string { return "audit-1" },
		})
	}

	cases := []struct {
		name      string
		workflows *stubWorkflowStore
		call      func(*ResolutionService) error
	}{
		{
			name:      "approve candidate",
			workflows: &stubWorkflowStore{},
			call: func(s *ResolutionService) error {
				_, err := s.ApproveCandidate(context.Background(), ApproveCandidateInput{CandidateID: "", Actor: system})
				return err
			},
		},
		{
			name:      "manual resolve",
			workflows: &stubWorkflowStore{},
			call: func(s *ResolutionService) error {
				_, err := s.ManualResolve(context.Background(), ManualResolveInput{
					InstallationID:    "inst-m05",
					ProviderCode:      "mercado_livre",
					ProviderItemID:    "MLB-SYS-MANUAL",
					InternalProductID: 0,
					Actor:             system,
				})
				return err
			},
		},
		{
			name:      "reject listing",
			workflows: &stubWorkflowStore{},
			call: func(s *ResolutionService) error {
				_, err := s.RejectListing(context.Background(), RejectListingInput{
					InstallationID: "inst-m05",
					ProviderCode:   "mercado_livre",
					ProviderItemID: "",
					Actor:          system,
					Reason:         "scripted sweep",
				})
				return err
			},
		},
		{
			// Undo's guard sits inside undoAuditEntry, past the audit lookup, so
			// this case needs a real target to reach it at all.
			name: "undo resolution",
			workflows: &stubWorkflowStore{audits: []productlinksdomain.ProductLinkAuditEntry{{
				AuditID:        "audit-target",
				InstallationID: "inst-m05",
				ProviderCode:   "mercado_livre",
				ProviderItemID: "MLB-SYS-UNDO",
				Action:         productlinksdomain.ProductLinkActionApproveCandidate,
				NextState:      productlinksdomain.ProductLinkStateResolved,
				CreatedAt:      now,
			}}},
			call: func(s *ResolutionService) error {
				_, err := s.UndoResolution(context.Background(), UndoResolutionInput{AuditID: "audit-target", Actor: system})
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.call(newService(tc.workflows))
			if err == nil || err.Error() != "PRODUCT_LINKS_SYSTEM_ACTOR_NOT_PERMITTED" {
				t.Fatalf("error = %v, want PRODUCT_LINKS_SYSTEM_ACTOR_NOT_PERMITTED", err)
			}
			if len(tc.workflows.applied) != 0 {
				t.Errorf("applied = %#v, want nothing written", tc.workflows.applied)
			}
		})
	}
}
