# D6-R2 — Independent Fable Adversarial Review: Method v2.3 + Whole-Repository Coherence

> **Reviewer:** Fable (independent adversarial review)
> **Date:** 2026-08-24
> **Scope:** Frontend Product Experience Planning Method v2.3 (supplied file) + Marketplace Central as method consumer + whole-repository architectural coherence
> **This document is review EVIDENCE only. No finding becomes authority without operator/GPT adjudication.**

---

# 1. Review identity

```text
main                        d409e643126a1d58cb1047be121da1ee4e7b92f0        (revalidated)
candidate branch            stage/d6-r2-frontend-realization
candidate HEAD              60c954d9caeb6869bad7d7c3b8b6428cdc37ef7a        (revalidated)
PR                          #61 — docs(d6-r2): open complete frontend realization closure
                            OPEN / DRAFT / MERGEABLE / NOT MERGED           (revalidated)
CI on candidate HEAD        ci #622 SUCCESS · pr-title #697 (+#698) SUCCESS (revalidated)
methodology reviewed        Frontend Product Experience Planning Method
                            v2.3 — OPERATOR-RATIFIED 2026-08-23
                            (supplied file: Downloads/functional-html-wireframe-method.md;
                             NOT present in the repository — see §6)
review branch               review/d6r2-methodology-whole-repo-fable
review isolation            branched from exact candidate HEAD 60c954d9;
                            merge-base(review, candidate) == 60c954d9;
                            review adds only docs/work/current/ai-dialog.md;
                            candidate untouched; no force pushes
roadmap state consumed      D7-R + D8-R OPERATOR-RATIFIED/ACCEPTED; AUTH GLOBAL-MAXIMUM CLOSED;
                            B00/B01/B00-R2/B11/B12/B110 LOCKED; B10 P8 RESUMED (next action)
```

The candidate branch moved from `5ea52e22` to `60c954d9` between the prior AuthorizationRequest review and this one; this review consumed the **current** HEAD `60c954d9` and its roadmap (D7-R/D8-R accepted, B10 resumed). All frozen `Exact next action` prose inside routed stage documents was treated as historical snapshot per repository law, not as drift.

# 2. Verdict

**ACCEPT WITH BOUNDED FIXES.**

- Method v2.3 is architecturally sound and is the strongest version of this method to date; findings against it are amendments, not invalidation.
- Marketplace Central genuinely follows the method rather than merely citing it; compliance findings are hygiene-level except the version-adoption drift (C-1).
- The whole-repository pass surfaced **one MATERIAL cross-authority contradiction (R-1: If-Match on the `:decide` custom-method URI vs the ratified W1 carrier law)** and one IMPORTANT hygiene trend (R-2: OAD source orphan growth). Neither invalidates the AuthorizationRequest Global-Maximum architecture (model C stands), the 106/31 census, or any operator LOCK.

# 3. Methodology findings (Objective A)

## M-1 — LOCK revalidation after upstream authority change is under-specified — IMPORTANT

- **Claim attacked:** §5.3 "preserve valid LOCKED blocks unless new evidence falsifies them" is a sufficient lock-safety law.
- **Evidence:** §5.3 names the preservation rule but no obligation to *enumerate and disposition every existing LOCKED block* when accepted authority materially changes. Who decides "unless falsifies", and against which checklist?
- **Counterexample/falsifier:** an authority change whose impact on an already-LOCKED block is real but unnoticed — nothing in the method forces the sweep, so a stale LOCK survives structurally until P11/P12.
- **Why the method survives partially:** MPC practice invented the missing discipline on its own: NOTIF-01 triggered an explicit bounded **B00-R2** shell reopen; D2-R6/D5-R6 triggered an explicit targeted **B110 P8 reopen** ("Method v2.1 therefore requires a targeted P8 reopen before final P9"); D8-R produced a per-flow impact matrix. The method got the right behavior from a disciplined consumer, not from its own text.
- **Global-Maximum alternatives:** (A) keep as-is, rely on consumer discipline; (B) add a small law: on material authority change, every currently-LOCKED block receives an explicit disposition `UNAFFECTED / REVALIDATE / REOPEN` recorded with the rebaseline. (B) costs one table per rebaseline and removes a silent-staleness class.
- **Recommended disposition:** adopt (B) as a v2.3.x amendment.
- **Smallest owning section:** §5.3.

## M-2 — OPEN assumptions can ride into LOCK without lock-time disposition — IMPORTANT

- **Claim attacked:** assumption tracking (§4, §21) is strong enough to prevent OPEN assumptions from contaminating LOCKED work.
- **Evidence:** §14.6 exit requires "no blocking finding remains" — but a registered assumption is not a finding, so a block may LOCK while an assumption *material to that block's structure* is OPEN, with the falsification deferred to P12. Worst case is a safety-critical block: MPC's A02 (device/work-floor evidence) is material to B80 Fulfillment structure; if B80 ever locked with A02 OPEN and P12 falsified it, the most safety-critical block would reopen at the end of the program.
- **Why the method survives partially:** the P14 `material assumptions OPEN = 0` closure law bounds the damage window, and MPC practice again self-corrected: P5 explicitly conditions B80 rendering on "P6 B80 + A02 evidence path". But that routing was consumer judgment, not method law.
- **Global-Maximum alternatives:** (A) status quo (P12/P14 backstop only); (B) require each block's LOCK record to list the OPEN assumptions its structure depends on, with an explicit operator risk-acceptance for each; (C) forbid LOCK while any structurally-material assumption is OPEN (too strong — would have blocked B01 on A04 and stalled the program for evidence that only operation can produce).
- **Recommended disposition:** (B). YAGNI-safe: it is one line per assumption at a decision that already exists.
- **Smallest owning sections:** §14.6 + §21.

## M-3 — P11 composition mechanics are unspecified — IMPORTANT

- **Claim attacked:** "P11 assembles already-LOCKED block prototypes" is executable as written.
- **Evidence:** P8 artifacts are, by law (§3.9), technically disposable single-file HTML with local fixtures. §17 does not say how N disposable files become one integrated product with cross-block navigation, deep links and shared shell behavior without either (a) rewriting the locked evidence (mutating what the operator adjudicated) or (b) accreting a mini-framework that drifts toward the production-authority §3.9 forbids. There is also no fidelity law: nothing requires the assembled P11 to be verified as structurally equivalent to each locked P8.
- **Counterexample/falsifier:** a P11 assembly that quietly diverges from a locked block (different reading order, dropped state) would pass every written P11 exit criterion because the criteria are all cross-block, none per-block-fidelity.
- **Why this matters now:** MPC will hit P11 with ~15 blocks. MPC already has a working pattern that generalizes: per-block executable structural verifiers (`verify-d6-r-b*.mjs`). Requiring assembled-P11 to pass the same structural invariants per block closes the fidelity hole.
- **Recommended disposition:** amend §17 with (1) an explicit statement that P11 is a *new* assembled artifact, locked P8 files remain immutable evidence; (2) a fidelity obligation: assembled P11 must preserve each locked block's structural invariants, demonstrably.
- **Smallest owning section:** §17.

## M-4 — v2.3 compression dropped two useful v2.1 instruments — OPTIONAL

- v2.1 §7 (Proportionality / "Small single-block delta" path) and v2.1 §16.4 (operator walkthrough question set) have no v2.3 equivalent. §3.5/§12/§23 carry part of the proportionality intent, but the explicit cheap path for trivial deltas is gone, and the walkthrough scaffold that structures what "operator uses it" means is gone. Risk: ceremony inflation on small changes; unstructured, disposition-only operator reviews. Note the MPC lock records: every lock shows `material changes requested: NONE` — the loop demonstrably has teeth (IA-01 was a real operator falsification mid-program), but the *recorded* evidence of what was operated is thin, and v2.1's walkthrough questions would have fixed that for free.
- **Recommended disposition:** restore both in condensed form. No authority impact.
- **Smallest owning sections:** §3.5 (proportionality) + §14.5 (walkthrough record).

## M-5 — P1 discovery can be satisfied entirely by authority-derived synthesis — OPTIONAL

- §7 (P1) accepts needs derived from accepted actor classes with zero direct user evidence; frequency/urgency capture is "recommended", and nothing forces real-operator evidence before block design. MPC's N01–N16 were derived from D0 actor classes; A01 (frequency), A02 (devices), A03 (terminology) have been OPEN since P1 across every LOCK. This is the strongest residue of backend anchoring the method still permits — not backend-*shaped UX* (the IA/flows are job-shaped), but backend-*sourced personas*.
- **Why it survives:** for an operator-run product where the operator *is* a primary user class, the LOCK loop itself is user evidence; and the P12/P14 assumption laws force eventual validation. Largely subsumed by M-2's lock-time assumption disposition.
- **Recommended disposition:** no structural change; fold into M-2's amendment.

## Attacks run against the method that FAILED (method survives) — see also §8

P8-LOCK-before-P9 ordering; static-approval risk; all-at-once generation; screen-shaped API law strength; YAGNI/DEFERRED-REJECTED strength; reference prestige import; block-vs-whole-product risk; premature IA lock; accessibility/responsive gating; pattern consolidation timing; status-model layering. Details in §8.

# 4. Marketplace Central methodology-compliance findings (Objective B)

## C-1 — v2.1 → v2.3 adoption drift; ratified method text exists only outside the repository — IMPORTANT

- **Repo evidence:** `docs/development/frontend-product-experience-planning-method.md` is **Version 2.1** (integrated at main HEAD `d409e643`, PR #60, 2026-08-22). `docs/index.md:34` routes "Frontend Product Experience Planning Method v2.1". Every D6-R2 program document cites v2.1. The operator-ratified **v2.3 (2026-08-23)** exists only as a loose file outside the repository (`Downloads/functional-html-wireframe-method.md`). No v2.2/v2.3 text, citation or ratification record exists anywhere in the repo.
- **Classification: adoption drift, not silent fork and not executed-decision damage.** I checked every operative v2.1→v2.3 delta against what the program actually did:
  - v2.1 permits **static** P8 lock media; v2.3 mandates browser-operable functional HTML. All seven `qualification/d6-r2-wireframes/*.html` artifacts are functional (state machines, scenario switches, onclick/addEventListener interaction — verified file-by-file), and the P8 ledger §1 already declares the executable-HTML law locally. Practice = v2.3.
  - v2.3 adds §3.10A (no backend-shaped UX), the P7 blocking law and the disposition vocabulary; the program already practiced all three (OP-READ-01 smallest owner-local repair; NOTIF-01 P7 feasibility lines with PRESENT-IN-AUTHORITY dispositions; P9-F1 → AuthorizationRequest upstream reopen instead of UX degradation). Practice = v2.3.
  - Nothing executed under v2.1 is forbidden by v2.3, and nothing v2.3 requires was skipped in executed work. **No LOCK is invalidated.**
- **Forward risk (why IMPORTANT, not OPTIONAL):** (1) B10 resumes P8 *next*; under the routed v2.1 text a future block could legally LOCK on a static wireframe — the exact regression v2.2/v2.3 exist to prevent; (2) the operator-ratified method text has no durable home — a single unversioned file in a Downloads folder is not survivable authority.
- **Smallest correction:** replace the repo method file content with v2.3 verbatim + update the `docs/index.md` route row to v2.3 (one file + one line). Historical stage documents citing "v2.1" remain frozen snapshots — do **not** rewrite them. This belongs on the candidate branch (or a bounded docs commit), not on this review branch.

## C-2 — P8 Block Ledger is stale relative to lock state — OPTIONAL

- **Repo evidence:** ledger header and §2.4 say B12 `LOCK: NO / PENDING OPERATOR` and carry no B110-LOCKED row; actual locks live in `D6-R2-NOTIF-01-D6-R-P8-RATIFICATION.md` (B12) and `D6-R2-P8-B110-APPROVALS-RATIFICATION.md` (B110), each with an explicit supersession sentence, and the roadmap (sole status authority) states the current truth. So this is *reachable* and authority-clean — but the method (§21/§23) designates the block ledger as the per-block record, and a fresh actor who opens the ledger for block state reads a wrong lock state unless they already know the supersession chain.
- **Smallest correction:** one supersession pointer line in the two affected ledger rows at the next legitimate docs commit. No authority change.

## C-3 — B10 P6 study predates the v2.3 disposition vocabulary — OPTIONAL

- The B10 Preparation reference study dispositions every surfaced reference capability in prose/matrix form ("NO baseline", "no invented score", rejected with Product-ownership reasoning — notably **never** with "current API lacks it", so the §12 anti-law is respected in substance). It does not use the explicit v2.3 disposition classes. When B10 resumes P8, restating the study's capability dispositions in the five-class vocabulary is a ten-minute conformance touch, not a redo.

## C-4 — Verifier token typo — OPTIONAL / cosmetic

- `scripts/verify-d6-r-b110-approvals-wireframe.mjs:128` emits `d6_r_b110_notification=AWAWARENESS_NOT_CAPABILITY` (sic), propagated into the B110 ratification record. Harmless (log token only), fix opportunistically.

## Compliance audit — everything else checked and PASSING

- **P0–P5 executed as global foundation:** authority pack bounded; N01–N16 job-shaped (not endpoint-shaped); UF01–UF16 complete with failure branches; 99/99 coverage with zero orphans and zero invented capability; P4 ran as *falsification* of accepted IA (not redesign) and correctly exited CANDIDATE, with global-IA LOCK earned only through the rendered B00 cycle — exactly the v2.1/v2.3 law.
- **Blocks are really enumerated** (P5 §9: B00…B120, 15 blocks, sequenced with per-block P6/P7 trigger dispositions) — no "B04+ NOT OPEN" opacity.
- **Bounded rebaselines for 99→104→106 exist and are reachable:** NOTIF-01 D6-R feed-forward carries the P1/P3/P5/P7 deltas for the five Notification operations; P5-B110 supersession + final P9 §12 carry the AuthorizationRequest deltas; every one of the ten new human operations has an exact frontend home; internal D3 request-creation correctly has **no** frontend home. The base P0–P4 record staying frozen at 99/30 is legitimate append-only history, not drift.
- **LOCKs are genuinely operator-made on operated functional HTML.** The reviewed HTML deliberately self-reports `CANDIDATE` while LOCK authority lives in ratification records — this is the *correct* resolution of "assistant artifact must not act as operator LOCK", and I verified no artifact self-claims lock authority.
- **B10 resume gating is correct:** suspended through the D6-R→D7-R→D8-R chain, resumed only after D8-R ratification, next action scoped against current 106/31; the 7-operation delta (Notifications/Governance) does not touch Preparation semantics, so the existing B10 candidate remains valid input.
- **IA-01 and OP-READ-01 are model examples of method §26 order-of-classification** — IA failure fixed as IA (B00 re-lock, routes preserved); missing read truth fixed as the *smallest owner-local* Product repair (D5-R2) with negative controls against dashboard/workflow/priority inventions, revalidated by D8-R2.
- **Responsive/accessibility are structurally evidenced** (locked responsive laws per block; keyboard/focus/non-color laws in the feed-forward §14; mobile drawer behavior working in the artifacts) — not merely mentioned.
- **Assumptions are carried, not laundered:** A01–A05 plus the two NOTIF-01 assumptions remain OPEN with named P12 probes; none was silently promoted to fact; A05 (bulk) correctly produced *absence* of bulk UX rather than a feature.

# 5. Whole-repository architecture findings (Objective C)

## R-1 — `CreateAuthorizationDecision` If-Match on the `:decide` custom-method URI contradicts the ratified W1 carrier law — MATERIAL

- **Exact current authority in conflict:**
  - `docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md` (operator-ratified, F-GPT-1 convergence): §14.2 — "A custom method such as `/resource/{id}:verb` does **not** use the base resource's ETag as `If-Match`, because the custom-method URI is a different HTTP request target… the custom request carries the acted-on resource's same opaque validator as **typed technical request data**"; forbidden-list item 9 — "`If-Match` carrying a base-resource ETag on a different `:verb` request target by private convention"; closing law — "HTTP `If-Match` is reserved for a request whose literal target is the protected resource."
  - Canonical OAD at candidate HEAD: `openapi.yaml:153` binds `POST /organizations/{organization_id}/authorization-requests/{authorization_request_id}:decide`; `paths-authorization-requests.yaml:42-64` gives `CreateAuthorizationDecision` an **`If-Match` header** (+`Idempotency-Key`), body `{outcome}` only, with `412`/`428` responses. The ETag it consumes is minted by `GetMyActionableAuthorizationRequest` on the **base** request URI — precisely the forbidden pattern.
  - Internal inconsistency proof inside the same OAD: `ResolveEconomicAttribution` (`…:resolve`, same C-class custom-verb shape) carries the validator as a **typed body field** — `ResolveEconomicAttributionRequest` requires `etag` (`components.yaml:649-653`) and has 409/422, no If-Match/412/428. Two carrier conventions for the same class in one contract.
- **How it survived:** D5-R5 §"If-Match = current AuthorizationRequest validator" and D5-R6 chose the header carrier without citing, testing or amending W1; the prior independent AuthorizationRequest review (F-1…F-5), GPT adjudication, operator ratification, final P9 (`P9_DECISION_CONCURRENCY:IF_MATCH_AUTHORIZATION_REQUEST`), D7-R (fingerprint includes "supplied If-Match") and D8-R all built on top of it. This is a silent contradiction between two operator-ratified authorities — exactly the class AGENTS.md orders surfaced rather than silently resolved.
- **What it is NOT:** not an architecture defect. Model C (AuthorizationRequest + immutable Decision) stands; request-local concurrency semantics are identical under either carrier; census 106/31 unchanged; no operation/Permission change is implied. It is a wire-carrier coherence defect that makes the contract non-implementable-as-written: an implementer reading W1 builds typed-etag/422/409; an implementer reading D5-R6/P9/D7-R builds If-Match/428/412.
- **Global-Maximum alternatives** (see §7 matrix): (A) bounded wire correction of `:decide` to typed `etag` body field (W1 stays universal; D5-R6/P9/D7-R/D8-R texts get a mechanical carrier substitution; semantics untouched); (B) ratified W1 amendment creating explicit alias semantics for request-local decide (preserves current OAD bytes; weakens the converged universal law and leaves `:resolve` vs `:decide` asymmetric); (C) leave both, document the fork (rejected — perpetuates the contradiction).
- **Recommended disposition:** operator adjudication required; recommend (A) as Global Maximum — one universal carrier grammar, smallest total authority change, zero semantic change.
- **Smallest owning authority to reopen:** the D5-R5/D5-R6 wire-carrier choice (bounded), never W1's universal law by convenience and never the AuthorizationRequest architecture.

## R-2 — OAD source orphan accumulation is growing without a control — IMPORTANT

- **Evidence (measured at candidate HEAD):** 16 orphaned `pathItems` in source files that the root OAD no longer references — `paths-economics-governance-sales-materialization.yaml` (superseded `AuthorizationDecisions`/`AuthorizationDecision`), `paths-fulfillment-postsale-work.yaml` (all 7 superseded Work items), `paths-identity-portfolio-readiness.yaml` (5 superseded access items), `paths-notifications.yaml` (2 superseded MyNotification items) — plus chain-orphaned superseded schemas (e.g. the old target-shaped `AuthorizationDecision`/`CreateAuthorizationDecisionRequest` in `components.yaml:665-684`, shadowed by the root remap at `openapi.yaml:252`). The whole-D8 review counted 5 orphans; the supersession pattern has tripled the count in one repair cycle. No verifier bounds it (no orphan/unreferenced check exists in `scripts/`).
- **Risk class (already proven in this repo's history):** source-text checks over a superseded copy go vacuous once the bundle prunes it, and a future editor can edit the dead copy believing it canonical — the old target-shaped `CreateAuthorizationDecisionRequest` sitting next to the live outcome-only one is a live instance of that trap.
- **Smallest correction:** either prune superseded pathItems/schemas at each supersession commit, or add one gate verifier asserting `unreferenced source pathItems == 0` (with an explicit allowlist if any orphan is deliberately retained). Docs/tooling only; no Product semantics.

## Architecture invariants actively re-challenged and CONFIRMED COHERENT

Verified directly at candidate HEAD (not assumed from prior reviews):

- **Census/fences:** 106 operations / 31 ordinary Permissions / H/A/S enforced by the gate verifier chain (`package.json` gate scripts chain all 12 verifiers; CI #622 green; local gate run below). New Notification ops are `authenticated`+H (no ordinary permission for self-Inbox); Governance ops carry `governance.decide`/`governance.read` independently; `notifications.manage` gates routing only.
- **F-3 shrink real:** zero occurrences of `requester_or_initiator_principal_id` / `predecessor_authorization_request_id` anywhere in `contracts/api/product/`.
- **F-4 real:** the semantic 503 Problem `type` is a `const` of the exact `authorization-validity-unavailable` URI on the `:decide` response; only that operation carries it.
- **Replay-before-If-Match (F-5):** D7-R §4.3 fixes the order inside the owner transaction, Principal-scoped namespace (`org + effective Principal + operation + key digest`), with falsifier list items 1–3 covering lost-201 replay, cross-Principal same-raw-key isolation and reused-key mismatch.
- **AuthorizationDecision ≠ target execution; Notification awareness ≠ capability; Work ≠ authority; zero-decider fallback forbidden** — each carried consistently across D5-R6 wire, P9 negative controls (`P9_FORBID:*`), D7-R runtime order and D8-R falsifiers; no layer quietly re-derives capability from possession.
- **Idempotency-state ≠ PITR continuity oracle:** D8-R §12 explicitly composes D7-R replay retention with the D7-R1 recovery fence; restored key/record presence is not timeline proof; fence arms automatically when lineage cannot be positively established.
- **Golden-flow set remains smallest falsifiable:** the four review-basis kinds land exactly inside GF-01/GF-02 where the governed actions live; a GF-04 would duplicate defect classes without a new invariant — the D8-R rejection is correct, and no flow was added merely because Governance exists.
- **No screen-shaped API and no backend-shaped UX in the realized blocks:** every forbidden-convenience list (no unread count, no mark-all-read, no search platform, no candidate-search API, no approval bulk/filters) is an explicit adjudicated DEFERRED/REJECTED with a named reopen trigger or P12 probe — not a silent API-absence surrender. Conversely, where a real gap existed the program reopened upstream (OP-READ-01, AuthorizationRequest) instead of degrading UX.
- **Frontend composition ≠ cross-owner workflow authority:** R71/R21 region-composition rules, operational-cockpit negative controls (`no operational_stage/next_action/priority/total_count`, no cross-owner Kanban) and the D5-R2/D8-R2 pair remain mutually consistent.
- **Probed but dismissed (no finding):** requester-side rejection-rationale absence (AuthorizationDecision has no reason field; F14 carries outcome atom only; a requester without `governance.read` learns outcome + target continuation, not rationale). Under §3.10A discipline: no evidence yet proves the material human need, a free-text reason field adds disclosure surface, and the target owner's current state answers "what now". Correct disposition today is a **named P12 probe**, not an upstream finding — recorded here so it cannot silently disappear.

# 6. Version/adoption drift determination

```text
supplied to review          v2.3 — OPERATOR-RATIFIED 2026-08-23 (outside repo only)
integrated in repository    v2.1 (main d409e643, PR #60, 2026-08-22)
operationally practiced     v2.3 semantics (functional P8, blocking law, dispositions, §3.10A)
v2.2                        intermediate draft (2026-08-22, outside repo), never integrated
```

**Classification: v2.1→v2.3 adoption drift (C-1, IMPORTANT).** Not a silent fork (no diverging local edits — the repo file is a faithful v2.1), not a justified consumer deviation (no deviation record exists), and not merely historical (the index routes v2.1 as the *current* execution method while the operator has ratified v2.3). No executed decision or LOCK is invalidated by the delta — verified operative-clause by operative-clause in §4/C-1. Smallest correction: land v2.3 verbatim as the repo method file + update the single `docs/index.md` route row; freeze historical v2.1 citations untouched.

# 7. Global-Maximum challenge matrix

| Area | Current design | Alternative A | Alternative B | Alternative C | YAGNI pressure | Whole-system tradeoff | Selected recommendation |
|---|---|---|---|---|---|---|---|
| **R-1 decision-wire carrier** (MATERIAL) | `If-Match` header on `:decide` URI (+428/412), contradicting ratified W1 | Typed `etag` body field on `:decide` (+422/409), W1 stays universal | Ratified W1 amendment: explicit alias semantics for request-local decide; keep header | Document the fork, keep both conventions | No new capability either way; pure coherence | A: mechanical text updates to D5-R6/P9/D7-R/D8-R + OAD, zero semantic change, one grammar. B: zero OAD change but two grammars forever (`:resolve` ≠ `:decide`) and weakens a converged universal law. C: perpetuates contradiction | **A** — operator adjudication required |
| **C-1 method version** (IMPORTANT) | Repo routes v2.1; ratified v2.3 homeless | Land v2.3 file + index row | Keep v2.1, record deviation | — | Zero — one file | A restores single-source ratified authority before B10 P8 resumes | **A** |
| **M-1 lock revalidation** (IMPORTANT) | §5.3 preservation prose only | Add UNAFFECTED/REVALIDATE/REOPEN sweep per rebaseline | Status quo + consumer discipline | — | One table per rebaseline | MPC already pays this cost voluntarily; law makes it portable | **A** (method amendment) |
| **M-2 lock-time assumptions** (IMPORTANT) | Assumptions linked, not dispositioned at LOCK | LOCK record lists depended-on OPEN assumptions + operator risk-acceptance | Forbid LOCK on OPEN material assumption | Status quo (P12/P14 backstop) | A: one line per assumption | B over-blocks (A04 would have stalled B01); C defers falsification to the most expensive moment | **A** (method amendment) |
| **M-3 P11 composition** (IMPORTANT) | Unspecified assembly | New assembled artifact + per-block structural-fidelity verification (reuse `verify-d6-r-b*` pattern) | Ad-hoc assembly, review-only fidelity | — | Verifiers already exist as a pattern | A prevents silent lock divergence at ~15-block scale | **A** (method amendment) |
| **R-2 OAD orphans** (IMPORTANT) | 16 orphaned pathItems, growing, uncontrolled | Prune at each supersession OR orphan==0 gate verifier | Leave; bundle is canonical anyway | — | Verifier is ~30 lines | Bundle correctness is safe today; source-drift trap is the proven risk class | **A** |

# 8. Confirmed non-findings (attacked and survived)

1. **P8 LOCK before final P9 (the handoff's sequencing challenge — proved safe, not assumed).** The load-bearing gate is the P7 feasibility line + v2.3 blocking law (P8 MUST NOT begin with an unresolved blocking upstream finding), not P9. The one real historical failure of this class — P9-F1, where P9 exposed that a `governance.decide` human lacked a purpose-bounded review object — occurred under v2.1 (which had no blocking law), and the method's escape path (bounded reopen → targeted B110 P8 re-lock → final P9) repaired it without collateral restart. v2.3's blocking law + disposition vocabulary is precisely the codification of that lesson. Moving full P9 rigor before LOCK would bind exact contracts to unlocked structure and churn; current ordering is the Global Maximum given the P7 gate exists.
2. **Operator-LOCK semantics.** Attacked as both too weak (rubber-stamp risk: every lock shows zero revision rounds) and too strong (bottleneck). Survives: IA-01 proves the operator loop falsifies materially mid-program; B12's operator-confirmed-meaning record proves comprehension, not ceremony; and CANDIDATE-in-artifact / LOCK-in-ledger separation cleanly prevents assistant artifacts from acting as authority. Residual improvement is only M-4's walkthrough record.
3. **Functional-HTML P8 requirement (over-prescription attack).** For this product class the requirement earns its cost: every one of the seven artifacts encodes knowledge-state switching (unknown≠empty≠unavailable) that static frames cannot falsify, and fixture discipline (deterministic local fixtures + CANDIDATE self-report + disposable-code law) blocked fake Product authority in practice. §3.5 still allows static exploration where it belongs (P7).
4. **Block-by-block vs whole-product.** IA-01 (global grouping falsified during B10 review, after B00 lock) is exactly the feared failure — and it was caught and repaired *mid-program* by the operator loop + smallest-reopen law, long before P12. P4/P5 + continuous operation are adequate safeguards; P12 is a backstop, not the first line.
5. **YAGNI resistance.** Reference study (B10, Notifications) imported task patterns and rejected prestige capabilities with Product reasons; the repo's forbidden-convenience ledgers (counts, search, bulk, saved views, candidate-search API, GF-04, workflow engine, second broker) are uniformly explicit rejections with reopen triggers — speculative-feature pressure is demonstrably contained.
6. **No screen-shaped backend / no backend-shaped UX (both directions).** OP-READ-01 is the canonical positive case: a real frontend-proven need repaired by the smallest owner-local read enrichment, with negative controls forbidding the dashboard-endpoint shortcut, and D8-R2 revalidation. P9-F1→AuthorizationRequest is the deep case: the frontend need reopened D1/D2/D3/D5 rather than being weakened. Nothing in the realized blocks navigates backend nouns.
7. **Status/program layering.** Roadmap-as-sole-mutable-authority + frozen snapshots survived a deliberately hostile read (the stage doc's own status header contradicts current state and is *correctly* superseded); only the ledger hygiene note (C-2) remains.
8. **Notification/Governance permission fences at the wire.** Independently re-derived from the YAML, not the docs: self-Inbox needs no ordinary Permission and no `notifications.manage`; decide/read independence is real; candidate discovery exposes `principal_id + display_name` only.
9. **D7-R/D8-R consumption of the eight Fable runtime obligations.** All eight traceable into D7-R sections + falsifier lists; D8-R adds the PITR composition seam (idempotency ≠ continuity oracle) and keeps the probe ledger closed with P2 preserved.

# 9. Progression impact

**B. B10 may resume only after bounded fixes.**

Precisely: nothing in B10 Preparation itself is defective — its candidate, P6 study and placement all survive review. But (1) R-1 is a MATERIAL contradiction between two ratified authorities and must be operator-adjudicated before further authority-consuming progression (the fix is bounded either way and touches no B10 content); (2) C-1 should land before B10 P8 resumes, because resuming P8 under routed v2.1 while the operator has ratified v2.3 would deepen the exact adoption drift this review was asked to catch. Both fixes are small; once adjudicated/landed, B10 resumes unchanged. No D6-R2 stop and no upstream architecture reopen (option C) is warranted — the AuthorizationRequest Global-Maximum closure stands.

# 10. Exact next action

Operator adjudicates R-1 (carrier class: typed-etag correction vs ratified W1 alias amendment — recommendation: typed-etag) and authorizes the C-1 v2.3 method integration; then B10 P8 resumes per the roadmap.

---

# Special review question — grounded answer

**"Does Method v2.3 provide a sufficiently strong bidirectional mechanism to discover and repair missing Product/backend capabilities without (a) preserving a backend-shaped local maximum or (b) letting frontend convenience destabilize the architecture?"**

**Yes — with the caveat that part of the proof predates the text that now mandates it.**

**Direction (a) — the mechanism demonstrably breaks backend-shaped local maxima.** Three repairs of ascending depth, all triggered by frontend realization and all resolved upstream rather than by degrading UX:
- **IA-01:** the accepted `OPERAÇÕES` grouping was falsified by operator mental-model evidence during rendered-block review → smallest IA reopen, re-lock. (Method category 2 failure, repaired as IA.)
- **OP-READ-01:** the evidenced operational cockpit exposed owner-local read contracts too poor for honest triage → the need was *not* dropped to fit the API, and the API was *not* given a screen-shaped `/operational-dashboard`; instead the smallest owner-local projection/filter enrichment (D5-R2) landed with explicit negative controls, and D8-R2 proved no choreography drift. This is §3.10A executing perfectly in both directions at once.
- **P9-F1 → AuthorizationRequest:** the deepest case — frontend contract-binding exposed that the accepted Governance model could not give a legitimate decider a purpose-bounded review object. The program did not weaken the locked B110 need and did not bolt on a query-only convenience (that remedy was explicitly SUPERSEDED); it reopened the smallest owning chain D1-R2→D2-R6→D3-R3→D5-R6, re-locked B110, re-traced P9, survived independent adversarial review, and fed genuine runtime obligations into D7-R/D8-R. Frontend evidence falsified backend planning four documents deep, and the architecture came out *more* coherent.

**Direction (b) — convenience is demonstrably contained.** The same chain shows the brakes working: the F-3 shrink removed wire fields with no legitimate frontend consumer (frontend presence pressure did not keep them); every "nice" capability (counts, search, bulk, candidate-search API, GF-04, workflow engine) died an explicit REJECTED/DEFERRED death with recorded Product reasoning and named reopen triggers; and the smallest-owner law kept every repair bounded (D5-R2 did not become a dashboard platform; the AuthorizationRequest repair did not become a workflow engine — alternative D was constructed and rejected twice).

**Where the honest separation lies:** those repairs were executed under v2.1's text, which lacked §3.10A, the blocking law and the disposition vocabulary. The *practice* was ahead of the *method*; v2.3 is the codification of what MPC's execution discovered (P9-F1 is visibly the ancestor of the P7 blocking law). So: the repairs prove the v2.3 mechanism works — because v2.3 *is* those repairs, generalized. The residual method weaknesses this review found are exactly the places where MPC's discipline is still ahead of the text: lock-revalidation sweeps (M-1), lock-time assumption disposition (M-2), and P11 assembly fidelity (M-3). Amend those, land v2.3 in the repo (C-1), and the method text finally equals the practice that validated it. The one thing the whole apparatus missed — R-1's carrier contradiction — was missed by *every* layer including the prior independent review, which is an argument for the method's §27 independent-review recursion, not against the method: it was still caught before a single line of Product code exists.

---

*End of adversarial review. No fixes applied; no authority modified; PR #61 untouched; B10 untouched.*
