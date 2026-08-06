# ADR-028: Auto-link only on concordant anchors

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed the MIS-006-integracao-fundacao mission's
product-links auto-linking policy (M-05-auto-vinculo), amended 2026-07-25 (D-121-2,
ratified by the operator), but no document was ever written. It is reconstructed here
from the two live-code citations of the two-digit token `ADR-05` (a numbering collision
with an unrelated pre-existing `005-mercado-livre-first-control-plane.md`; see
`docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md`), harvested at
`docs/architecture/decisions/_citations/adr-05-twodigit-citations.md`. Every clause below
is traceable to code and tests that already assert it. Nothing here is new policy.

## Context

When a Mercado Livre listing is matched against the internal product catalog, the matcher
can be handed two independent anchors: the seller SKU (compared against `CODPROD`) and the
EAN. Either anchor alone can resolve to a product, and a resolved anchor is not proof of
identity by itself — a single matching field is exactly the kind of coincidence a human
should confirm, not the kind of fact a system should act on unattended. Before this policy
was ratified, the original wording auto-approved on "EAN exact and unique" alone; that was
amended (D-121-2) after operator review because a single anchor, however unique, is still
one field agreeing, not two independent fields agreeing with each other.

## Decision

**The only automatic linking path is CODPROD and EAN both resolving to the same product,
with no hard-negative signal between the listing and that product. Every other outcome —
one anchor only, several anchors disagreeing, several products colliding on one anchor, or
a title contradiction overriding agreement — goes to a human, never to auto-approval.**

### The clauses

**§1 — Auto-approval requires both anchors to resolve to the same product.** The single
automatic path is stated in code as CODPROD and EAN concordance, and nowhere in
`AutoApproveCandidate` is a single-anchor match accepted.
> `apps/server_core/internal/modules/product_links/application/resolution_service.go:217`
> — "AutoApproveCandidate applies the single automatic path (D-121-2, ADR-05 amended):
> CODPROD and EAN resolving the same product, with no hard negative."
> `apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go:14`
> — "The auto-linking policy ratified as ADR-05 amended (D-121-2): the ONLY automatic path
> is CODPROD and EAN resolving the same product."

**§2 — A hard negative blocks the automatic path even with both anchors agreeing.** A
title contradiction (kit-vs-unit, color, voltage) between the listing and the resolved
product overrides concordant anchors and keeps the candidate out of auto-approval.
> `apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go:292`
> (`TestHardNegativeBlocksTheAutomaticPathEvenWithBothAnchors`) — "a title contradiction
> beats agreeing anchors. Must-fail proof: dropping the detectHardNegative check in
> buildConcordantCandidate makes this listing auto-approve."
> `apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go:314`
> (`TestHardNegativeKindsBlockConcordantSKUAndEAN`) exercises kit/combo, color, and voltage
> contradictions and asserts each candidate resolves to REJECT with nothing auto-approved.

**§3 — A single resolved anchor goes to confirmation, never auto-approval.** When only
CODPROD or only EAN resolves, the candidate's match status is CONFIRM (a one-click,
human-in-the-loop state), and the missing anchor is named in the reason rather than left
generic.
> `apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go:167`
> (`TestSingleAnchorGoesToConfirmationNeverAutoApproved`) — "a single anchor resolves one
> product, but nothing corroborates it — so it waits for a human, and the warning names the
> anchor that is missing (AC-11 forbids a generic one)."

**§4 — An anchor that resolves to more than one product goes to review, approving
nothing.** An EAN matching several distinct products offers every one of them for human
review; none is auto-selected.
> `apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go:225`
> (`TestCollidingAnchorGoesToReviewAndApprovesNothing`) — "an EAN matching four products...
> identified none of them. Every product is offered for review and nothing is approved."

**§5 — Anchors pointing at different products go to review with no winner elected.** When
CODPROD and EAN each independently resolve, but to different products, both candidates are
surfaced as REVIEW and neither is discarded or auto-selected.
> `apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go:257`
> (`TestConflictingAnchorsApproveNothingAndElectNoWinner`) — "the anchors point at different
> products. Neither wins — the operator decides, and nothing is written automatically."

**§6 — The automatic path is closed by construction, not merely by convention.**
`AutoApproveCandidate` refuses a candidate whose `MatchStatus` is not `Accept` — i.e. one
the generator itself did not mark corroborated — even if called directly.
> `apps/server_core/internal/modules/product_links/application/resolution_service.go:230-233`
> — `if candidate.MatchStatus != domain.LinkCandidateMatchStatusAccept { return false,
> errors.New("PRODUCT_LINKS_AUTO_APPROVE_NOT_CORROBORATED") }`.

**§7 — Auto-approval never overrides a decision already in force.** A listing the operator
already resolved or rejected is left untouched by a later automatic run; the link's own
state, not just the decision trail, is checked so links resolved before the audit trail
existed are not silently reopened.
> `apps/server_core/internal/modules/product_links/application/resolution_service.go:249-254`
> — "A listing the operator already settled is settled. The link's own state is what proves
> it... links resolved before E10 existed carry no decision at all, so the trail alone would
> report those as undecided and the automatic path would reopen them."
> `apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go:426`
> (`TestAutoApprovalNeverOverridesAnOperatorDecision`) and `:457`
> (`TestAutoApprovalNeverReopensAListingTheOperatorRejected`) exercise both sides.

**§8 — Only a system actor may take the automatic path; every operator-facing entry point
refuses one.** `ApproveCandidate`, `ManualResolve`, `RejectListing`, and `UndoResolution`
all reject a system actor with `PRODUCT_LINKS_SYSTEM_ACTOR_NOT_PERMITTED` before validating
or writing anything; `AutoApproveCandidate` is the sole path a machine may take.
> `apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go:576-591`
> — "Every operator-facing entry point refuses a machine. AutoApproveCandidate is the sole
> exception... a machine may corroborate, and nothing else (D-121-2 follow-up, operator
> ruling)."

## Rationale

A single matched field is a coincidence a catalog of any size will produce regularly —
SKUs and EANs are reused, mistyped, or shared across near-identical variants. Requiring two
independent fields to agree, and additionally requiring that no textual signal contradicts
that agreement, is the minimum bar at which an automatic decision is more likely correct
than a coin flip informed by one fact. Below that bar the system's only honest move is to
hand the decision to someone who can look at both listings side by side. The original
"EAN exact and unique" rule fell below that bar because uniqueness of a match is not the
same claim as corroboration by a second, independent anchor — the amendment (D-121-2)
exists because the operator caught that distinction in review.

## Consequences

- A product legitimately sold under many listings is not penalized: two listings resolving
  the same CODPROD both auto-link, because concordance is evaluated per listing, not
  per product (`TestOneProductAutoLinksToManyListings`).
- Re-running generation over the same snapshot does not duplicate the link or the decision
  row; the automatic path is idempotent per listing (`TestRerunningGenerationDoesNotDuplicateTheAutomaticLink`).
- A failed approval in a batch does not abandon the rest of the batch — every corroborated
  candidate is still offered, and failures are reported by name
  (`TestAFailedApprovalDoesNotAbandonTheRestOfTheBatch`).
- The collision count backing a decision is carried from the moment the generator judged
  the candidate, never re-derived at approval time, so the audit row reflects what was true
  when the decision was made, not what is true now.

## Alternatives Considered

**EAN exact-and-unique alone (the original, pre-amendment wording).** Superseded by
D-121-2: a unique match on one field is not corroboration, and the operator ruled that a
second independent anchor must agree before the system decides unattended.
> `apps/server_core/internal/modules/product_links/application/auto_link_policy_test.go:14`
> names this as the policy being amended, though the pre-amendment wording itself is not
> preserved in code — see Unverified claims.

## Unverified claims

- The exact original pre-amendment wording ("EAN exact-and-unique") is stated in the
  MIS-006 mission record per the harvest, but the two live-code anchors verified here
  (`resolution_service.go:217`, `auto_link_policy_test.go:14`) only name the amended
  (D-121-2) rule as current; the prior rule's precise scope is not independently
  reconstructed from code in this document.
- `root.go:1002` ("ADR-05: flag defaults OFF", the `mlCatalogOffersEnabled` read-gate) was
  investigated separately and is not part of this decision — it names a different rule
  under the same colliding token and is out of scope for ADR-028 (see accompanying reply).
