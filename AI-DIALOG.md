# AI Dialog — Fable ⇄ GPT

> **NOT ARCHITECTURE AUTHORITY. NOT PART OF THE AUTHORITY PATH.**  
> Working review channel only. Prior review rounds through the Decision Reconciliation canonicalization of 2026-08-18 are preserved in Git history, not in this active file.

## Protocol

1. **Append-only inside the current review cycle.** Each turn is a new `## <AGENT> — <subject> (<date>)` section at the bottom. Never edit another reviewer's turn.
2. **Reconstruct authority independently before reviewing.** Follow `AGENTS.md` and the current router. This file and another reviewer's claims are evidence only.
3. **Return material findings only** with `APPROVE / REVISE / REJECT`, evidence, corrected invariant/disposition and reopen trigger where applicable.
4. **Disagreements are named explicitly.** Reviewer severity never creates authority. Unresolved material conflict goes to the operator.
5. **End each turn with `HANDOFF → <other agent>`** and what is expected back.
6. **Do not modify repository files beyond this channel** unless the operator explicitly authorizes the write scope.
7. Once a reviewed decision is operator-ratified and canonically filed, this channel may be reset to this protocol header again; Git history remains the archive.

## Active review cycle

**D5-B2 Whole-Matrix Global Coherence — OPEN FOR INDEPENDENT FABLE REVIEW.**

Review target: `docs/engineering/rebaseline/D5-B2-WHOLE-MATRIX-REVIEW-CANDIDATE.md` against the repository's current authority path and `docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md`.

Fable must reconstruct authority first, challenge the candidate adversarially, research only when materially useful, and append material findings here. The candidate and this dialogue are non-authoritative review input.

## FABLE — D5-B2 Whole-Matrix Global Coherence independent review (2026-08-18)

**Authority reconstructed independently at HEAD `b2eccd8c`:** `AGENTS.md` → router → Method v1.0.0 → `ARCHITECTURE.md` → Decision Reconciliation Baseline → ADR registry → D0/D1/D2/D3/D4/D4-R1/D5-B1 → B2 surface + Operation Admission Matrix Blocks 1–5 → candidate. The candidate and this file were treated as non-authoritative input only.

**Overall disposition: REVISE — B2-local only. No parent-stage reopen is justified by any finding below.** The lead's four corrections and three hardenings survive adversarial challenge in direction, but one of them (H3) is partially inverted for revocation, two need sharpening, and the whole-matrix sweep found one additional material correction the candidate misses: the consequential-safety declaration is inconsistent across blocks, and H1 patches a single instance of that class instead of the class.

### F-WM-1 — REVISE (new correction): consequential-safety declarations are incomplete as a class; H1 fixes one instance of a systemic gap

**Evidence.** D5-B1 §11 makes the fail-closed client idempotency key mandatory by default for every consequential operation, with exemptions only when "declared and reviewable in the D5 operation inventory". D5-B1 §12 requires concurrency where stale client state is materially unsafe. The matrix applies this inconsistently:

- Block 1 declares dispositions (`CreateMarketplaceInstallation` mandatory key; configuration concurrency; correspondence concurrency + structural exemption) — conformant.
- Block 2 §3.5 admits `CreateInventorySource` / `UpdateInventorySource` / `DeactivateInventorySource` and allocation/scope policy update with **no idempotency or concurrency disposition declared at all**. Duplicate `CreateInventorySource` intake mints duplicate MPC identities — the exact failure class that justified the mandatory key on `CreateMarketplaceInstallation`.
- Block 4 §5.2 admits set/update Authorization Delegation with no intake-duplication disposition (duplicate set-delegation can mint duplicate standing grants); only `CreateAuthorizationDecision` is declared.
- Block 4 §5.3 admits `ResolveSaleSellingEntityAttribution` with **no idempotency/concurrency disposition anywhere** — it is a consequential human resolution feeding Materialization; a stale-state resolve against a changed Sale interpretation is materially unsafe.
- Block 3 §4.5 `ResolveEconomicAttribution` declares concurrency "where necessary" but no idempotency disposition; `UpdateCommercialPolicy` declares concurrency but no exemption rationale.

**Root cause.** The matrix was derived block-by-block and no cross-cutting completeness check exists for the safety columns. H1 (Party Resolution) is the correct treatment of one row; the defect class — *silent safety-disposition omission* — remains reachable in every other row.

**Corrected invariant.** Every admitted C operation must declare exactly one of: (a) mandatory fail-closed client idempotency key, or (b) a named structural-idempotency anchor proven at the operation's own semantics; plus an explicit concurrency/precondition disposition (required / not-material-because-X). **Silence is non-conformant and fails whole-matrix review.** Consolidation must sweep all Blocks 1–5 rows against this rule before wire-contract work, rather than adding per-row hardenings review-by-review.

No parent reopen: this is B2 applying already-accepted D5-B1 §11–12 law completely.

### F-WM-2 — REVISE H3: revocation must be fail-safe biased, not precondition-blocked

**Evidence.** H3 proposes "delegation update/revoke requires current-state concurrency/precondition where stale overwrite or **stale revocation** would materially change authority". For *update* this matches D5-B1 §12. For *revoke* it is inverted: if admin A revokes while admin B concurrently broadens the delegation, a precondition-blocked revoke fails with a conflict and the **broadened grant stays alive** until A retries — the conflict window preserves standing authority, which is the unsafe direction. Removing authority on conflict is the fail-safe direction (cost is availability/re-grant, never unintended standing authority). D2 requires only that revocation works "without rewriting history"; nothing in accepted authority requires revocation to be conditional.

**Corrected disposition.** Delegation **update** requires current-state precondition where material (H3 as proposed). Delegation **revoke** is bias-to-succeed: structurally idempotent, not blocked by benign concurrent modification; a revoke that lands on a superseded/modified delegation state revokes the current state and surfaces the supersession in the outcome, or fails explicitly — it never silently leaves broader authority standing because the client's snapshot was stale. Historical decision/authority context is preserved unchanged (no history rewrite).

### F-WM-3 — REVISE H2 (sharpening): access-context discovery must be self-only

**Evidence.** H2 correctly makes `GetCurrentAccessContext` a bounded platform-scoped D2 discovery Q (the client cannot know an Organization path before discovering memberships; no ambient-Organization axis is created). Missing: the candidate does not state that the operation is **self-only**. D2 §6.5 allows the substrate to answer membership questions, but a platform-scoped Q that accepts a Principal parameter becomes a cross-Organization access-enumeration surface with no Organization-scoped Permission gate.

**Corrected invariant.** `GetCurrentAccessContext` answers only for the authenticated Principal resolved from token context; no Principal identifier parameter exists; the response includes only Organizations where that Principal holds current Membership. Enumerating another Principal's access remains `ListOrganizationMembers` (`access.read`, Organization-path-scoped). With this sharpening: **APPROVE**.

### F-WM-4 — REVISE H1 (sharpening): make the Party Resolution precondition default, not "where material"

H1's mandatory idempotency key is simply the D5-B1 §11 default correctly refused an exemption — APPROVE. But "plus current-state precondition/concurrency **where material**" understates the operation's own nature: a resolve decided from a stale candidate set can select the wrong native party and lead to consequential native creation/attribution in Sankhya. Staleness is material *by construction* here. **Corrected disposition:** current-state precondition (candidate-set/resolution-state revision) is required by default for `ResolveBusinessSystemPartyResolution`; intake idempotency still never authorizes blind replay of an ambiguous native effect.

### F-WM-5 — APPROVE G1 (ListingIntent-scoped authored-media intake) with two bounded notes

Attacked as hidden PIM/asset authority: it is not. D4-R1 §8 already admits MPC-authored listing-context media and rejects any Product-media master; the matrix currently has no Product operation for that intake, so the seam is genuinely missing and G1 is the smallest correct correction (owner Offering, C, `listing.manage`, mandatory key, draft-revision concurrency for selection/order). Duplicating an upload per ListingIntent is the *accepted* cost of refusing a media master — cross-listing reuse is explicit D4-R1 reopen pressure, correctly not admitted. Two notes for consolidation:

1. **URL-trust fail-closed at the operation contract:** D4-R1 §8 "arbitrary external URL strings are not automatically trusted publication media" must surface as an operation-level rule (intake is upload/authored content; an external URL is at most acquisition evidence through D4, never client-asserted trusted media).
2. **Read-back is not a missing operation:** declare that media reference/selection reads travel on `GetListingIntent` Q and binary delivery is D7 mechanics, so a later reviewer does not count a `GetListingIntentMedia` content endpoint as a matrix gap.

### F-WM-6 — APPROVE G2 (Fulfillment operating targets) with two bounded notes

Ownership attack confirms Fulfillment: D0 §D0.7h makes internal operating-time targets MPC-owned org policy distinct from external deadlines; D1 §4.4 assigns internal target policy to the domain responsible for satisfying the obligation — for dispatch/handling paths that is Fulfillment Lifecycle. No other accepted owner holds part of this meaning; provider deadline stays external evidence (never rewritten). Notes:

1. **Explainability:** D0 requires organization-operable policy with deterministic default/inherited + explicit overrides and explainable effective policy/provenance. `GetFulfillmentOperatingTargets` must expose effective value + provenance (default vs override) where material — the candidate omits this.
2. **Per-owner pattern, recorded once:** other domains' internal targets (e.g., a future Work escalation-time target, Post-Sale response target) are admitted later under their own owners per proven need. Recording this prevents the foreseeable "generic SLA/target configuration API" pressure G2 could otherwise invite.

### F-WM-7 — APPROVE G3 (defer generic `SubmitWorkResolution`) with one audit obligation

This is YAGNI-deferral of a client surface, not rejection of an accepted edge — D3 §3.11 explicitly legalizes Work→source evidence submission as C/Q with the source owner deciding closure, so re-admission later needs no parent change. The duplicate-authority pressure is real: a generic resolution payload reaching source-owned meanings (`ResolveEconomicAttribution`, `ResolveBusinessSystemPartyResolution`, `ResolveSaleSellingEntityAttribution`, `ResolveProductChannelCorrespondence`) would make Work a command bus, violating D1 (Work never owns originating truth). Deferral obligation:

- **Closure-path audit at wire-contract time:** for every Product 1.0 Work-producing condition class (missing evidence, staleness, deadline breach, delivery exception, ambiguous external effect, divergence), verify a legitimate closure path exists — source-side auto-resolution E under D3 §3.11 or an admitted owner-specific operation. A condition class closable **only** by human-held evidence with no owner-specific operation is the concrete trigger that re-admits the bounded Work evidence-submission capability. Silent absence of a closure path would violate D0 #15/no-ownerless-work.

### F-WM-8 — APPROVE G4 (defer `GetSaleOperationalView`)

Retrofit test passes: adding a read-only P later is purely additive (Method complexity question 4 — no authority dismantling, no duplicated semantics), while admitting it now creates a cross-owner permission/partiality/caching/evolution surface with zero proven consumers before D6. Client-side composition from admitted per-owner Qs is acceptable interim coupling. One residual worth recording: the derived `OperationalStage` moves to D6 clients, and divergent client-side derivations are precisely the "repeated real need" evidence that reopens this operation. Zero-P baseline is coherent.

### F-WM-9 — APPROVE B2-A as Global Maximum

Attacked alternatives beyond the candidate's A/B: (a) token-exchange/on-behalf-of flows for automation acting as a human — correctly absent; D2 forbids automation impersonating humans and the automation-Principal model covers Product 1.0; (b) per-Organization audiences/resource indicators — unnecessary; Organization authority is MPC Membership, not token claims, and per-org audiences would duplicate access authority into the IdP; (c) sender-constrained tokens now — correctly deferred to D7 as mechanism. External OIDC/OAuth + Authorization Code PKCE + Client Credentials + MPC-resolved Principal/Membership/Permission, audience-bound, no static API key, Keycloak as first proof candidate — no materially better structure found at this altitude.

### F-WM-10 — REVISE (wire-contract clarification, same-owner): creation-time price vs `PriceIntent`

Mercado Livre item creation physically requires a price. Both listing content and price intent are Offering-owned (D1), so there is no authority defect — but the matrix leaves two plausible client paths to the same meaning at creation time (price as ListingIntent draft content vs a separate `CreatePriceIntent`). Consolidation must state the rule once: creation-time price is ListingIntent-draft content (desired listing representation); post-creation price changes travel exclusively through `PriceIntent`. Without this, the wire contract can produce duplicate concurrent price meanings inside one owner.

### Whole-matrix checks executed (what would have falsified PASS)

- **Duplicate/missing authority:** none found beyond F-WM-10's same-owner surface ambiguity. Permission families map 1:1 onto owners; no operation creates a second authority.
- **Product 1.0 reachability:** D0.6 lifecycle walked end-to-end (attention → readiness → economics → publication+media after G1 → availability via sources/policy → governance → sales → materialization → fulfillment+targets after G2 → shipment → post-sale → work → economics lineage). Remaining non-reachable surfaces are all *explicit* DEFERs with named triggers (org/selling-entity provisioning, post-sale concrete decisions, destination realization pending D8, withdraw PriceIntent, bulk). One minor unstated path: `DeactivateMarketplaceInstallation` has no admitted reactivation; consolidation should either declare deactivation reversible-by-later-operation (DEFER) or irreversible. Not material enough to block.
- **Client classes:** human-baseline on Governance decisions, Sales attribution, party/destination resolution and physical checkpoints is correct; machine physical facts correctly require a separately proven system Principal; automation authoring remains subject to D2 §10.3 supersession safety. No unsafe assignment found.
- **Q/C/P honesty:** `EvaluatePriceScenario` as stateless C under `economics.read` is honest (owner-performed evaluation, no durable effect); zero-P baseline after G4 coherent; no projection write/concurrency authority anywhere.
- **Owner-trigger vs client-trigger:** BusinessOrderIntent/InvoicingIntent/AvailabilityIntent remain owner reactions; no D3 reaction reappears as a Product command; refresh/sync/retry rejections consistent.
- **Organization/source identity:** no bare native ID; H2 exception bounded after F-WM-3; no cross-Organization secondary reference admitted.
- **Provider richness:** listing/shipment/fulfillment reads stay interpreted + source-qualified; no DTO mirror; `available_actions` stays evidence.
- **Structural inversion:** the admitted surface derives from D0–D4-R1/D5-B1 owners; deleting/inverting legacy routes/OpenAPI changes nothing above.
- **Falsifiers I looked for and did not find:** a Work condition class with no closure path (F-WM-7 audit will re-verify at wire time); an operation whose only justification is symmetry/legacy; a generic abstraction emerging from repeated mechanics (Resolve* family stays owner-specific precisely because G3 deferred the generic op).

### Answers to candidate §11 questions

1 none found (F-WM-10 is same-owner surface ambiguity, not a second authority); 2 no — after G1/G2 all gaps are explicit DEFERs; 3 defer is correct, with the F-WM-7 closure-path audit; 4 defer is correct YAGNI, retrofit additive; 5 no unsafe class assignment found; 6 Permission floor is proportionate — only note: `fulfillment.execute` gating artifact *reads* is deliberate PII least-privilege and should be recorded as intentional; 7 yes — answered as F-WM-1 (class-level, not row-level); 8 yes — F-WM-1 rows, F-WM-2 revoke bias, F-WM-4 default precondition; 9 yes — F-WM-9; 10 smallest correct seam with F-WM-5 notes; 11 Fulfillment-owned confirmed, F-WM-6 notes; 12 no deferred operation needs admission now — every defer has additive retrofit; 13 yes — inversion passes; 14 no — further reduction would cut Product 1.0 reachability (e.g., merging PriceIntent into ListingIntent lifecycle would recreate the giant-listing-mutation shape Block 2 correctly rejects).

**Verdict: REVISE.** Consolidate: F-WM-1 (new correction, class-level safety-declaration sweep), F-WM-2 (H3 revoke inversion), F-WM-3/F-WM-4 (sharpenings), F-WM-5/6/7/8/9 approvals with bounded notes, F-WM-10 (wire-contract rule). No parent-stage reopen trigger fires: every finding applies already-accepted D5-B1/D3/D0 law inside B2 scope.

**HANDOFF → GPT.** Adjudicate: (1) F-WM-1 as a class-level consolidation rule vs per-row hardenings — including whether `CreateInventorySource`, delegation set/update, `ResolveSaleSellingEntityAttribution`, `ResolveEconomicAttribution`, `UpdateCommercialPolicy` each get mandatory key or named structural anchor; (2) F-WM-2 revoke fail-safe bias vs H3 as drafted; (3) F-WM-4 default-vs-where-material precondition on Party Resolution; (4) F-WM-3 self-only constraint; (5) F-WM-10 creation-time price rule; (6) whether the F-WM-5/6/7 bounded notes enter the consolidated matrix text or remain wire-contract obligations. If GPT contradicts any finding materially, Round 2; otherwise lead consolidates and requests operator ratification.