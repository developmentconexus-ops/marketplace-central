# D6-R2 — AuthorizationRequest Global-Maximum Independent Fable Challenge

> Review branch only: `review/d6r2-authorization-request-fable`
> Candidate branch: `stage/d6-r2-frontend-realization`
> Candidate HEAD expected: `5ea52e223b2c4f077a15872300a54c20e89f4c6e`
> Candidate PR: #61 — `docs(d6-r2): open complete frontend realization closure`
> Base `main` expected: `d409e643126a1d58cb1047be121da1ee4e7b92f0`
> Canonical Product candidate: **106 Product operations / 31 ordinary Permissions / Principal kinds H-A-S only**
> Active runtime baseline: **NONE**
> D7-R / D8-R: **BLOCKED**
> Product implementation: **BLOCKED UNTIL accepted D9**

## Purpose

Run the required isolated **independent Fable adversarial challenge** of the complete `AuthorizationRequest` redesign before it may be declared Global-Maximum closed or D7-R may open.

This is not a preference, naming or style review. Attempt to falsify the current architecture as one composed Product:

```text
D1-R2 boundary revalidation
→ D2-R6 AuthorizationRequest identity/lifecycle
→ D3-R3 communication/recovery
→ D5-R6 canonical 106/31 Product wire
→ B00-R2/B11/B12 Notifications experience
→ B110 Approvals experience
→ final P9 bidirectional Screen Contracts
```

The redesign was intentionally allowed to reopen previously accepted planning before implementation. Therefore do **not** preserve it merely because the operator/GPT already approved it. If a broader targeted reconstruction is materially better under current Product needs, say so and prove why. Conversely, do not recommend broader generic workflow/BPM/case-management machinery without a concrete current consumer/falsifier.

Reviewer output is **Evidence, not authority**. Do not edit the candidate branch or PR #61. Write only below `## Fable response` in this file on this review branch.

## Mandatory revalidation before analysis

Independently record and verify:

1. remote `main` HEAD;
2. `stage/d6-r2-frontend-realization` HEAD;
3. PR #61 state/draft/base/head/mergeability and changed-file count;
4. exact GitHub Actions status on candidate HEAD — expected CI #601 SUCCESS and pr-title #671 SUCCESS;
5. this review branch ancestry from the exact candidate HEAD;
6. `candidate...review` differs by **only** `docs/work/current/ai-dialog.md`;
7. Product OAD proof currently reports 106/31 while historical 95/29 + 99/30 remain green.

If the candidate HEAD differs from `5ea52e223b2c4f077a15872300a54c20e89f4c6e`, **STOP** and report `STALE_REVIEW`. Do not review a moved candidate.

## Strict reading discipline

Start exactly:

1. `AGENTS.md`
2. `docs/index.md`
3. `docs/roadmap.md`
4. `docs/engineering/rebaseline/D6-R2-P9-AUTHORIZATION-REQUEST-BIDIRECTIONAL-SCREEN-CONTRACTS.md`
5. `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-R6-AUTHORIZATION-REQUEST-OAD-WIRE-PROOF.md`

Then read this bounded redesign pack:

6. `docs/engineering/rebaseline/D6-R2-NOTIF-01-D1-R2-GOVERNANCE-BOUNDARY-REVALIDATION.md`
7. `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R6-AUTHORIZATION-REQUEST-IDENTITY-CANDIDATE.md`
8. `docs/engineering/rebaseline/D6-R2-NOTIF-01-D2-R6-RATIFICATION.md`
9. `docs/engineering/rebaseline/D6-R2-NOTIF-01-D3-R3-AUTHORIZATION-REQUEST-COMMUNICATION-CANDIDATE.md`
10. `docs/engineering/rebaseline/D6-R2-NOTIF-01-D3-R3-RATIFICATION.md`
11. `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-R5-AUTHORIZATION-REQUEST-PRODUCT-SURFACE-RATIFICATION.md`
12. `docs/engineering/rebaseline/D6-R2-NOTIF-01-D6-R-P9-F1-SUPERSESSION-RATIFICATION.md`
13. `docs/engineering/rebaseline/D6-R2-NOTIF-01-D6-R-P8-RATIFICATION.md`
14. `docs/engineering/rebaseline/D6-R2-P8-B110-APPROVALS-RATIFICATION.md`

Use exact wire when material:

- `contracts/api/product/openapi.yaml`
- `contracts/api/product/paths-authorization-requests.yaml`
- `contracts/api/product/paths-notifications-authorization-request.yaml`
- `contracts/api/product/paths-notifications.yaml`
- `contracts/api/product/paths-work-authorization-request.yaml`

Use rendered P8 evidence only when challenging frontend operability:

- `qualification/d6-r2-wireframes/b00-r2-notifications.html`
- `qualification/d6-r2-wireframes/b11-notifications-inbox.html`
- `qualification/d6-r2-wireframes/b12-notification-routing-settings.html`
- `qualification/d6-r2-wireframes/b110-approvals.html`

Read broader accepted D0–D6 authority only when one concrete falsifier cannot be answered from this bounded pack. Repository current authority > this review prompt.

## Candidate invariants to attack

### 1. Missing pre-decision identity — was `AuthorizationRequest` really the right repair?

Challenge whether a canonical Governance-owned pre-decision episode is genuinely required and correctly separated from:

```text
Business Intent
AuthorizationDecision
Operational Work
Notification
```

Attempt to prove one of these instead:

- pending authorization is not a material identity/lifecycle and should remain a derived projection;
- `AuthorizationDecision` could safely model pending + terminal state without semantic collapse;
- `AuthorizationRequest` duplicates Business Intent identity/authority;
- the chosen identity is too weak to distinguish materially separate reauthorization episodes;
- the chosen identity is too strong and creates unnecessary case/workflow state.

Specifically attack the case where the same target revision may require two materially different authorization episodes because governing authority/policy/evidence changed.

### 2. Lifecycle minimality and historical explainability

Falsify the lifecycle:

```text
PENDING → DECIDED
PENDING → INVALIDATED
```

Challenge whether another current state is required or whether one of these is unnecessary. Distinguish:

- human reject;
- invalidation without human decision;
- reauthorization;
- no eligible decider;
- source drift;
- stale concurrent decision.

Challenge `1 AuthorizationRequest → 0..1 AuthorizationDecision` and bounded predecessor lineage for reauthorization.

### 3. Review-basis snapshot design

Attack the claim that Governance needs an immutable **typed + bounded** decision-purpose snapshot while the source owner remains authoritative.

Challenge all four families:

```text
listing_intent
price_intent
business_order_intent
invoicing_intent
```

Attempt to find:

- missing information that makes a legitimate approver unable to understand what is being authorized without broad source read;
- copied data that transfers business ownership into Governance;
- mutable/current source data incorrectly presented as historical decision basis;
- generic payload/metadata pressure hiding a missing canonical schema;
- privacy/disclosure excess in snapshots;
- historical reconstruction depending on current mutable owner truth.

### 4. Material validity vs concurrency

Challenge the separation:

```text
AuthorizationRequest ETag/revision → request-local concurrency
semantic owner revalidation       → material authorization-basis validity
```

Construct cases where:

- irrelevant target drift should not force reapproval;
- material governing-context drift must invalidate even if target ETag is unchanged;
- two approvers race;
- authority/delegation changes after read but before decision;
- source owner is unavailable during validity revalidation.

Attempt to prove target ETag should be used differently, or that the Q/revalidation protocol is insufficient/overbuilt.

### 5. Duplicate-safe intake and authorization episode identity

Attack `authorization_episode_key` semantics from D3-R3:

- retry after ambiguous owner→Governance acceptance must find the same Request;
- materially new reauthorization must create a new Request;
- retries must not collapse different authorization episodes;
- retries must not duplicate the same episode.

Look for a missing stable owner-side identity or an accidental hidden distributed transaction requirement.

### 6. D5 Product surface minimality

Challenge the selected public surface:

```text
ListMyActionableAuthorizationRequests
GetMyActionableAuthorizationRequest
CreateAuthorizationDecision  (reanchored on request)
ListAuthorizationDecisions
GetAuthorizationDecision
```

Attempt to prove:

- one of the two actionable reads is redundant;
- an additional public Create/Invalidate/Reauthorize request operation is genuinely required;
- internal D3 communication was incorrectly exposed or hidden;
- actionable and historical reads should share/should not share authority differently;
- the resulting 106/31 census hides another necessary Product operation;
- a screen-shaped API exists only because of frontend convenience.

Do not optimize for a particular operation count. Count is consequence, never target.

### 7. Least privilege and Permission separation

Falsify:

```text
governance.decide != governance.read
governance.decide != source-owner read
governance.read   != governance.decide
```

Attempt to prove the bounded actionable projection leaks too much or too little. Challenge whether a human can legitimately decide with `governance.decide` but without `offering.read` / `materialization.read`, and whether the review basis makes that safe and understandable.

Possession of `AuthorizationRequestID`, Notification, URL or candidate-list presence must never grant authority.

### 8. F13 action-required Notification

Challenge:

```text
AUTHORIZATION_ACTION_REQUIRED
→ AuthorizationRequestRef
→ GetMyActionableAuthorizationRequest
→ current server eligibility
```

Attack stale/delayed materialization, decision by another approver before delivery, request invalidation, lost eligibility, Organization switch and historical Notification retention.

Attempt to prove F13 should point somewhere else or carry additional authority/state. Reject any design where Notification becomes a capability token or current Governance truth.

### 9. F14 decision-result Notification

Challenge the intentional asymmetry:

```text
AUTHORIZATION_DECISION_RESULT
→ target-oriented continuation
```

It must not require `governance.read` merely to return the requester to work and must not default to `AuthorizationDecisionID` by symmetry.

Attack target resolution for all four target families, especially BusinessOrder/Invoicing where the eventual human destination is the current Sale composition. Determine whether existing owner reads provide enough correlation under current source authorization without a hidden new API.

### 10. Zero eligible human deciders

Attack the chosen disposition:

```text
PENDING approval-required request
+ known zero current human deciders
→ explicit Operational Work condition
```

Challenge whether this is the correct existing owner/boundary, whether it duplicates request state, whether Work can reconcile safely, and whether the frontend can surface it through existing `ListWork` / `GetWork` without new API.

Fail the design if Work assignment/escalation can accidentally become authorization or if any implicit fallback approver (admin, owner, requester, `governance.decide` holder, Notification recipient) is created.

### 11. Decision retry, idempotency and stale concurrency

Challenge the human write:

```text
POST .../authorization-requests/{id}:decide
If-Match: AuthorizationRequest ETag
Idempotency-Key: same semantic attempt
body: outcome only
```

Construct:

- two simultaneous approvers;
- 201 response lost after commit;
- 503 known no-effect validity-unavailable response;
- 412 stale precondition;
- 409 no-longer-admitted/current-state conflict;
- browser refresh during confirmation;
- same request with a changed outcome;
- same outcome after request revision changed.

Check that no blind consequential auto-retry or idempotency-key misuse is required by P9.

### 12. P8/P9 frontend coherence

Challenge all locked affected surfaces:

```text
G00-E bell
U01 preview
R128 Inbox
R129 Notification routing
R110 Approvals queue/history
R111 request/history detail
R100/R101 Work disposition
```

Attempt to find:

- material UI with no backend owner;
- Product operation with no realistic human home;
- URL state becoming business authority;
- router becoming a second server-state cache;
- inline local draft becoming durable source truth;
- Organization switch retaining incompatible scoped state;
- a deep link that depends on remembered ambient authority;
- hidden Permission implication;
- screen-shaped parallel domain/client authority.

### 13. Notifications routing semantics

Retest the existing Notifications repair against the new Governance flow:

- self Inbox remains H-only without `notifications.read`;
- `notifications.manage` remains routing-admin only;
- candidate projection remains only `principal_id + display_name`;
- candidate presence remains discovery, not authorization;
- `CONFIGURED([])` remains impossible;
- historical ineligible recipient remains explicit without opaque-ID leakage;
- opening source remains separate from read/archive awareness mutation;
- no count/bulk/search platform is implied by symmetry.

### 14. Global-Maximum / YAGNI challenge

This review must answer both directions:

**Under-modeling:** does the redesign still leave a root cause that implementation would have to invent locally?

**Over-modeling:** did the repair introduce a workflow/case platform, generalized approval framework, extra state, API or abstraction without a current Product 1.0 consumer?

Explicitly compare the accepted design against at least:

1. query-only actionable projection over target/revision;
2. `AuthorizationDecision` with a pending state;
3. `AuthorizationRequest` + `AuthorizationDecision` (current design);
4. a broader generic approval/workflow engine.

State which is the Global Maximum under current constraints and why.

### 15. Proof quality

Independently inspect executable proof, not just prose. At minimum challenge whether current green gates genuinely establish:

- historical Product 95/29 non-regression;
- historical Product 99/30 non-regression;
- current Product 106/31;
- four typed review-basis families;
- AuthorizationRequest negative controls 10/10;
- Notification negative controls 8/8;
- B00-R2/B11/B12/B110 structural proofs;
- P9 six surfaces / ten human operations / 12 negative controls;
- generated TS/Go projection determinism where claimed;
- no active legacy runtime.

Look for tautological string checks that fail to prove a material semantic relationship, blind spots where prose and OAD could drift together, or proof that claims runtime behavior before D7.

### 16. D7 leakage

Search the redesign for premature commitment to:

- database/schema/RLS realization;
- queue/outbox/job implementation;
- exact HTTP router/framework;
- session/deployment topology;
- exact transaction semantics beyond D3 obligations;
- realtime transport;
- worker scheduling;
- persistent client state implementation.

Classify genuine future obligations as `D7_OBLIGATION`, not redesign blockers, unless current D1–D6 correctness already depends on an unchosen mechanism.

## Mandatory targeted falsifiers

Record PASS / FAIL / INSUFFICIENT-EVIDENCE with reasoning for at least these:

1. `AuthorizationRequest` is a materially distinct business identity, not a convenience DTO.
2. Same target revision can safely have multiple reauthorization episodes without identity collision.
3. `AuthorizationDecision` remains terminal history rather than pending lifecycle.
4. INVALIDATED cannot be confused with human reject.
5. Four review-basis families are sufficient and bounded.
6. Governance does not steal source-owner business truth.
7. Historical explainability does not reconstruct past decisions from current mutable truth.
8. Request ETag solves concurrent decision without becoming business-validity oracle.
9. Material validity can change independently of target ETag.
10. Irrelevant target drift does not force reapproval.
11. Intake retry cannot duplicate one authorization episode.
12. Reauthorization cannot collapse into a prior episode.
13. `governance.decide` does not imply `governance.read`.
14. `governance.decide` does not imply source read.
15. Actionable request IDs/URLs/Notifications cannot bypass current eligibility.
16. F13 cannot become a capability token.
17. Delayed F13 does not create stale actionable awareness.
18. F14 does not need AuthorizationDecisionID by symmetry.
19. F14 target continuation is viable for all four target types using current owner authority/correlation.
20. Zero-decider Work cannot become fallback authorization.
21. No new Work Product operation is required.
22. Two simultaneous approvers cannot create two Decisions.
23. Ambiguous successful decision retry reuses the same Idempotency-Key.
24. Known 503 no-effect does not trigger blind auto-retry.
25. 412 stays distinct from business rejection/invalidity.
26. Client decision body does not need target ETag or copied target authority.
27. R128 source-open cannot implicitly mark read.
28. R129 candidate presence cannot authorize route save.
29. `CONFIGURED([])` remains unrepresentable.
30. P9 URL/navigation state cannot become a second server/business-state authority.
31. Every one of the ten redesigned human Product operations has a coherent frontend home.
32. No hidden 107th Product operation is required by the composed UX.
33. No generic approval/workflow engine has a proven current consumer.
34. No D7 mechanism is prematurely selected.
35. Product implementation remains blocked until D9.
36. The accepted design is a true Global Maximum rather than the best patch inside a flawed earlier structure.

## Output contract

Append below `## Fable response` using this structure.

### 1. Verdict

Choose exactly one:

- `ACCEPT`
- `ACCEPT WITH BOUNDED FIXES`
- `REOPEN SMALLEST AUTHORITY`
- `REJECT / RECONSTRUCT`

### 2. Revalidation record

Record exact main SHA, candidate SHA, PR #61 state/base/head, changed-file count, exact CI results and review-isolation result.

### 3. Global-Maximum assessment

Answer explicitly whether the current structure is the best coherent design under real current Product constraints, including comparison against query-only, pending-Decision and generic-workflow alternatives.

### 4. Material findings

Number highest severity first. For each include:

- **classification:** `D1_FIX`, `D2_FIX`, `D3_FIX`, `D5_FIX`, `D6_FIX`, `D7_OBLIGATION`, `REPOSITORY_FIX`, `LATER_NON_BLOCKING`, `AUTHORITY_CONTRADICTION`, or `REVIEW_FALSE_POSITIVE`;
- **severity:** Critical / Important / Minor;
- exact candidate location;
- governing repository authority;
- concrete counterexample/failure;
- smallest correction;
- why it belongs at that exact stage/owner.

If there are no material findings, say so explicitly. Do not manufacture style findings.

### 5. Identity / lifecycle assessment

Adjudicate AuthorizationRequest identity, lifecycle, reauthorization lineage and separation from Decision/Work/Notification.

### 6. Communication / recovery assessment

Adjudicate episode-key retry, material-validity Q, invalidation, decision propagation, F13 delayed materialization and zero-decider handling.

### 7. Product API / least-privilege assessment

Adjudicate 106/31, actionable vs historical operations, Permission separation, public-vs-internal surface minimality and hidden-operation risk.

### 8. Notification / continuation assessment

Adjudicate F13 request-ref, F14 target-oriented continuation, PersonalNotifications awareness semantics and routing administration.

### 9. Frontend / P9 assessment

Adjudicate routes, URL/server/local state split, all six affected surfaces, Work disposition, stale/503 recovery and bidirectional operation homes.

### 10. Retry / concurrency assessment

Adjudicate If-Match, Idempotency-Key, ambiguous retry, 409/412/503 handling and multi-approver races.

### 11. Snapshot / historical explainability assessment

Adjudicate four typed review-basis families, disclosure bounds, source ownership and immutable historical explanation.

### 12. Proof assessment

State exactly what current executable proof establishes, what it does not, and any material proof blind spot.

### 13. YAGNI / reconstruction decision

Answer explicitly: is any broader reconstruction or new generic approval/workflow machinery justified **before implementation**?

### 14. D7 obligations

Separate current architecture defects from valid future runtime implementation obligations.

### 15. Continuation recommendation

State the smallest exact action after GPT/operator adjudication. Review output itself authorizes **nothing**: no merge, no D7-R, no D8-R, no B10 resumption, no Product implementation.

---

## Interaction rule

Fable writes **only** to this file on `review/d6r2-authorization-request-fable`. Do not edit PR #61, `stage/d6-r2-frontend-realization`, `main` or any other review branch.

GPT/operator will independently adjudicate every material finding against current repository authority and executable evidence. A second Fable round is justified only if a surviving material contradiction requires a bounded fix and fresh independent revalidation.

---

## Fable response

<!-- Fable: append the independent AuthorizationRequest Global-Maximum review here. -->

### 1. Verdict

**ACCEPT WITH BOUNDED FIXES**

No material contradiction survives against the selected architecture (model C: MPC-owned `AuthorizationRequest` + immutable `AuthorizationDecision`). The redesign is the Global Maximum under present requirements. The bounded fixes are two executable-proof blind spots, one wire-field disposition adjudication, one 503-discriminator adjudication, and one explicitly recorded D7 obligation. None reopens D1–D3 or the ratified D5 surface.

### 2. Revalidation record

Independently measured 2026-08-24 against remote after `git fetch --prune`:

```text
remote main HEAD                          d409e643126a1d58cb1047be121da1ee4e7b92f0  (= expected PR base)
candidate stage/d6-r2-frontend-realization 5ea52e223b2c4f077a15872300a54c20e89f4c6e  (= expected)
review review/d6r2-authorization-request-fable e473d269cae125b859ffbec87aae633adf9503fa (= expected)
PR #61                                    OPEN / DRAFT / mergedAt=null / base=main@d409e643 / head=5ea52e22
CI on 5ea52e22                            ci #601 SUCCESS · pr-title #671 SUCCESS (later pr-title #672 also SUCCESS)
review isolation                          candidate...review = 0 behind / 1 ahead; merge-base = 5ea52e22
only differing file                       docs/work/current/ai-dialog.md
Product OAD census                        106 operationIds wired at openapi.yaml root; historical 95/29 + 99/30
                                          verifiers present and green per CI #601 full gate
```

`STALE_REVIEW` condition: **not met**. The exact expected candidate was reviewed.

### 3. Global-Maximum assessment

**The current structure (C) is the Global Maximum under present Product requirements.** Each alternative is refuted by a present falsifier, not by preference:

**A — query-only current-target/current-revision projection: REFUTED by two independent present requirements.**
First, `CreateAuthorizationDecision` needs a concurrency anchor. Under A the only available anchor is the governed target ETag, which D0's materiality law forbids in both directions: irrelevant target drift would 412-force reapproval, and material governing-context drift (authority/policy/evidence) that leaves the coarse target revision unchanged would silently pass. Second, episode identity: the PI-100/revision-7 counterexample (episode A invalidated by governing-context drift; episode B required for the same target+revision) means `target+revision` is not injective over authorization episodes — F13 could not name which episode needs deciding, duplicate-safe intake has no anchor, and predecessor lineage is unrepresentable. A also converges toward C anyway: the Decision must still carry an immutable reviewed basis, so A saves no schema, only the identity that makes the rest coherent. A is not a smaller correct design; it is the root cause of P9-F1 preserved.

**B — AuthorizationDecision with pending/pre-decision lifecycle: REFUTED primarily by Permission-regime collision.**
Actionable pre-decision reads are H-only + `governance.decide` + exact-current-eligibility; Decision history is H/A/S + `governance.read`. Under B one resource must serve both disjoint authority regimes, so either every history read filters lifecycle (and a filter bug leaks actionable data to `governance.read` automation) or pending state leaks into history by construction. C makes the separation structural — the strongest form of `P9_PERMISSION_LAW:GOVERNANCE_DECIDE_READ_INDEPENDENT`. Secondary refutations: a pending Decision has no outcome/decider/decided_at, so B degrades the immutable-history invariant to "immutable eventually"; INVALIDATED then requires either a fabricated third outcome (violating "no fake reject") or record deletion (rewriting history); and the `1 request → 0..1 decision` law becomes a self-referential mutation rule instead of a relation.

**D — generic approval/workflow/request engine: REFUTED by absence of any present consumer.**
Present requirements are: single-stage, single-decider-commit, four closed target families, three lifecycle states, no quorum/stages/SLA/comments. Every generic capability would be dead weight with attack surface. D1-R2's rejection stands; the D2-R6 §16 exclusion list plus D5-R6 §8 absent-operations list plus executable negative control #10 (public `CreateAuthorizationRequest` leak fails the gate) keep D structurally out.

**C minimality check (over-modeling direction):** two identities (one new), three lifecycle states (each with a proven consumer: PENDING = actionability, DECIDED = history, INVALIDATED = truthful no-decision termination), one bounded optional predecessor link (known reauthorization consumer), one owner-local episode anchor (required by D3 §4.11 ambiguous-acceptance law), +2 Product operations, 0 new Permissions, 0 new Notification kinds, 0 new Work operations. I attempted to remove each element and each removal breaks a stated present law. Nothing found that exists only for UI symmetry: the two reads are respectively the queue source and the ETag-bearing decision anchor; a collection-as-detail merge would leave no resource whose ETag is the If-Match validator.

### 4. Material findings

No Critical findings. No finding reopens D1, D2-R6, D3-R3 or the ratified D5-R5 surface.

**F-1 — AuthorizationDecision target/review-basis pairing has no executable guard.**
- **classification:** `REPOSITORY_FIX` · **severity:** Important (proof), not semantic
- **location:** `scripts/verify-authorization-request-oad.mjs` lines ~187–192; schema `contracts/api/product/paths-authorization-requests.yaml` lines 283–287.
- **authority:** D5-R6 §5 — "The decision representation structurally preserves the same target/review-basis pairing used by the request"; repository verification rule "material guards require a deterministic falsifier/negative control".
- **counterexample:** delete the four-branch pairing `oneOf` from `AuthorizationDecision` (keeping the plain `target` + `review_basis` unions). `validate()` still passes — it compares the two unions individually and `validateActionableView()` checks pairing only for the actionable view. No negative control mutates decision pairing. All gates stay green while a `listing_intent` target paired with a `price_intent` basis becomes representable in Governance history, falsifying the D5-R6 §5 claim silently.
- **smallest correction:** assert the pairing `oneOf` on `AuthorizationDecision` (4 branches, target kind == review-basis kind per branch, mirroring `validateActionableView`) and add one negative control that strips it (census 10→11, with the ratcheted count updated in the same change).
- **why this owner/stage:** pure proof-harness gap; wire and semantics are already correct at HEAD.

**F-2 — P9 verifier proves prose self-consistency, not prose↔wire consistency.**
- **classification:** `REPOSITORY_FIX` · **severity:** Minor
- **location:** `scripts/verify-d6-r-p9-screen-contracts.mjs` (all 90 lines are string-marker asserts over the P9 markdown).
- **authority:** review mandate §15 — "blind spots where prose and OAD could drift together"; repository rule "presence is not execution".
- **counterexample:** a coordinated operationId rename (OAD + `verify-authorization-request-oad.mjs` `byId` literals updated together) leaves the P9 doc's `P9_OP_HOME:` markers naming a dead operation while every gate stays green. The ten-operation bidirectional trace — P9's core claim — is never checked against the bundled document.
- **smallest correction:** in the P9 verifier, bundle (or read the cached bundle of) the canonical OAD and assert each of the ten `P9_OP_HOME` names is a real `operationId`.
- **why this owner/stage:** proof harness only; the trace itself was manually confirmed correct at HEAD (all ten names exist in the wired OAD).

**F-3 — Two actionable-view wire fields have no consumer and one has an undischargeable purpose.**
- **classification:** `D5_FIX` (bounded) or explicit adjudication as non-UI wire disposition · **severity:** Minor
- **location:** `contracts/api/product/paths-authorization-requests.yaml` — `requester_or_initiator_principal_id` and `predecessor_authorization_request_id` on `ActionableAuthorizationRequestListItem` and all four `Actionable*AuthorizationRequest` variants.
- **authority:** P9 §1 (every material element has an owner/home or explicit non-UI disposition); D2-R1 (humans see labels, not opaque IDs); D2-R6 §11 (predecessor link exists so "history can distinguish" reauthorization).
- **counterexample:** the locked B110 wireframe renders neither field (verified: no requester/predecessor rendering in `b110-approvals.html`). For the H-only `governance.decide` consumer both are dead opaque IDs: the requester ID is unresolvable to a label without access-administration reads the approver need not hold, and the predecessor ID is undereferenceable — terminal requests have no Product read, `GetMyActionableAuthorizationRequest` fails on non-actionable requests, and `ListAuthorizationDecisions` has no `authorization_request_id` filter. Further, `AuthorizationDecision` does not carry predecessor lineage, so once the successor request is decided, the D2-R6 §11 distinguish-reauthorization purpose has **no reader anywhere on the wire** — the lineage survives only as internal Governance state.
- **smallest correction:** one adjudication line in the D5/P9 authority declaring both fields wire-level correlation data with explicitly deferred UI/history exposure (reopen-triggered), **or** drop them from the actionable views until a consumer exists. Do **not** add a history read or list filter now — no present consumer justifies it, and inventing one would violate the same YAGNI law the design correctly applies elsewhere.
- **why this owner/stage:** D5 projection choice; D2-R6 identity semantics are correct and unaffected — internal retention of lineage is proper. This is exposure-without-consumer, the mirror image of the usual missing-operation defect, and costs one adjudication sentence.

**F-4 — The 503 no-effect guarantee is carried only by response-body shape.**
- **classification:** `D5_FIX` (one-line adjudication) with a `D7_OBLIGATION` tail · **severity:** Minor
- **location:** `paths-authorization-requests.yaml` `AuthorizationValidityUnavailable` / `AboutBlank503Problem` (`type: about:blank`, status 503).
- **authority:** P9 §8.1 — "503 means … no Decision was recorded"; MPC safety rail — ambiguous potentially-accepted writes are not blindly retried.
- **counterexample:** an infrastructure 503 (proxy/gateway, before the decision commit point — or after it) is byte-wise indistinguishable in status and problem `type` from the semantic no-effect 503. P9's "known 503 no-effect … is a new human attempt (new key)" recovery is only safe when the client can prove the server evaluated the request. If a client treats a bodyless transport 503 as known-no-effect and mints a new Idempotency-Key, it converts an ambiguous possibly-accepted write into a second semantic attempt — exactly the class the safety rail forbids.
- **smallest correction:** adjudicate explicitly (P9 or D5-R6 errata line): only a parseable `AuthorizationValidityUnavailable` body is the known-no-effect discriminator; any 503 without that body is ambiguous acceptance and takes the same-Idempotency-Key recovery path. No wire change required; a distinct problem `type` URI would be stronger but is not mandated by current repo problem-type convention.
- **why this owner/stage:** the distinction is client-contract semantics (D5/P9), not runtime mechanics; D7 then implements it.

**F-5 — Idempotent-replay resolution must precede If-Match evaluation.**
- **classification:** `D7_OBLIGATION` · **severity:** Important as an obligation; not a D0–D6 defect
- **location:** `CreateAuthorizationDecision` semantics (D5-R6 §3, P9 §8.1).
- **authority:** D5-R5 §2.3 — Idempotency-Key exists for "ambiguous acceptance / exact retry recovery".
- **counterexample:** decision commits, 201 is lost, client retries with the same Idempotency-Key and same If-Match (per P9 law). The commit advanced the request revision, so a precondition-first runtime returns 412. The decider — holding `governance.decide` only — then has **no wire path to learn their own outcome**: the actionable GET now 404s and Decision history requires `governance.read`. The stated purpose of the Idempotency-Key is defeated. Correct behavior: same-key replay of a committed identical attempt returns the original 201/Decision before precondition evaluation; 412 is reserved for genuinely different attempts.
- **smallest correction:** record this evaluation-order law in the D7-R obligations list so it cannot be silently lost in runtime realization. The D0–D6 semantics already imply it; nothing reopens.
- **why this owner/stage:** evaluation ordering is runtime realization; flagging it now prevents a plausible D7 implementation that is green on every current proof yet breaks P9's retry promise.

### 5. Identity / lifecycle assessment

Sound. `AuthorizationRequest` is a materially distinct business identity (episode ≠ target+revision, proven by the governing-context-drift counterexample), not a convenience DTO — it carries its own lifecycle, concurrency revision, snapshot ownership and lineage that no existing identity can host. The three-state lifecycle is minimal-complete: I attempted to add states (expired, suspended, awaiting-validity) and each is correctly representable as PENDING plus explicit at-decide-time outcomes; I attempted to remove INVALIDATED and it forces either eternal PENDING or fabricated rejects. `1 → 0..1` Decision with never-reopen and new-episode-new-ID cleanly separates reject (human decision, F14 fires) from invalidation (no decision, no F14, owner E). Predecessor lineage is bounded to one optional field — see F-3 for its wire-exposure wrinkle, which does not touch the D2 semantics.

### 6. Communication / recovery assessment

Sound. The episode anchor (`Organization + target + authorization_episode_key`) gives exactly-once business meaning without exactly-once transport claims; same-anchor-different-basis is an explicit conflict, killing silent mutation. Material validity Q (`VALID | INVALID | UNKNOWN_OR_UNAVAILABLE`, never defaulting valid) correctly separates from request ETag concurrency, satisfying both D0 laws (no needless reapproval, no stale authorization). F13 delayed-materialization revalidation (still PENDING + still eligible) prevents stale actionable awareness; historical F13 is honestly retained and rechecks on open. Zero-decider is a known-empty vs unavailable distinction feeding the existing Work boundary with no fallback approver — the D3 §12 failure matrix covers every failure class I constructed, including the invalidation-E-lost case (Governance must reconcile; silent stale-PENDING forbidden). No event transport is ever sole authority.

### 7. Product API / least-privilege assessment

106/31 is sufficient; no hidden 107th operation is required. I specifically attacked: (a) decider outcome confirmation — discharged by the 201 body returning the created Decision, no `governance.read` needed; (b) requester pending-state visibility — owner intent reads, owner-owned disposition; (c) predecessor dereference — no present consumer (F-3); (d) org-wide pending oversight — deliberately absent (`ListAllAuthorizationRequests`), stalls surface owner-side and through zero-decider Work; acceptable for v1 with reopen triggers. Neither actionable read is redundant (queue source vs ETag-bearing anchor). Permission separation is structural and negative-controlled (decide widening to read, automation admission, search leak, target-in-body, ETag return, generic payload, F13 fallback, Work origin loss, public create — all fail the gate). The internal D3 intake/invalidation surface is correctly non-public. Nothing exists only for screen shape.

### 8. Notification / continuation assessment

F13 = `AuthorizationRequestRef`, verified in wire including the negative control against target-ref fallback; awareness-not-capability holds because the actionable GET re-derives eligibility. F14 target-oriented asymmetry is correct and verified: `AuthorizationDecisionID` would force `governance.read` on requesters or create capability leakage — the requester's need is continuation, not Governance history. Continuation viability checked against the wired OAD for all four families: `ListingIntent`, `PriceIntent`, `BusinessOrderIntent`, `InvoicingIntent` all have owner detail reads (`/listing-intents/{id}`, `/price-intents/{id}`, `/business-order-intents/{id}`, `/invoicing-intents/{id}`), and the BusinessOrder/Invoicing owner representations carry the sale correlation needed to reach the Sale composition under existing owner authority. No hidden continuation API is required. Notifications routing semantics (ten ORG_ROUTED slots, minimal candidate projection, `CONFIGURED([])` unrepresentable, ineligible-historical-recipient labeling) are unchanged from the locked repair and unbroken by the Governance flow; F13 is exact-principal-driven, not route-configured, so no routing fallback can fabricate an approver.

### 9. Frontend / P9 assessment

The bidirectional trace is genuinely closed at operation granularity: ten human operations, each with exactly one home; six surfaces, each with one owner/wire per material element; internal D3 semantics correctly have no screen. The four state classes prevent every second-authority failure mode I constructed (URL as business state, router cache, durable local drafts, cross-Org state retention — all named and forbidden with the Org-switch invalidation law). Stale/412, 409, 503 and no-longer-actionable each have distinct recoverable UX bound to rereading Governance truth. F-3 is the only trace-granularity gap (field-level, not operation-level). The zero-decider disposition through existing R100/R101 respects Work-not-authority; a `work.read` consumer sees the obligation exists but cannot read the request — correct fail-closed asymmetry.

### 10. Retry / concurrency assessment

If-Match on the request-local revision solves the two-approver race (first commit advances revision; second gets 412) without becoming a business-validity oracle — validity is separately revalidated per D3 §8. Eligibility-set churn not bumping the revision is correct: it avoids needless 412 storms while per-attempt eligibility revalidation closes the gap. Idempotency-Key is required-on-wire, reuse bounded to the same semantic attempt, auto-retry forbidden. The 409/412/503/428 vocabulary is distinct and each maps to a distinct recovery. Two residual sharp edges are F-4 (503 discriminator) and F-5 (replay-before-precondition ordering); both are recoverable with one adjudication line plus one recorded D7 obligation, and neither invalidates the design.

### 11. Snapshot / historical explainability assessment

The four typed families are closed, structurally paired with target kinds, and generic-field-proofed (`payload/metadata/data/attributes/entity_type/entity_id/fields/raw` all executable-forbidden). Disclosure is purpose-bounded: I checked the widest snapshot — `MarketplaceSale` in business-order/invoicing bases — and it carries sale ref, occurrence time, selling-entity attribution and lines only; no buyer PII beyond what party-resolution review materially requires, consistent with the PII-minimization rail. Optional price/economics evidence is correctly marked historical provenance, not current Economics authority. Historical explanation of DECIDED episodes reconstructs entirely from the immutable Decision (basis + target + outcome + decider + time) with zero dependence on current mutable owner truth — falsifier discharged. INVALIDATED episodes retain explanation internally with no wire reader; that is YAGNI-correct today (see F-3 note) and reopen-triggered.

### 12. Proof assessment

What current proof establishes: 106/31 census with uniqueness; historical 95/29 and 99/30 non-regression; the three Governance operations' exact method/path/owner/Permission/principal contracts; list-query closure (`limit,cursor` only); detail-ETag presence; decision body closure (`outcome` only) with required If-Match + Idempotency-Key; four closed review-basis families with per-family required fields and generic-field prohibition; actionable-view pairing; target refs ETag-free; F13 `AuthorizationRequestRef` (with anti-fallback control); F14 four-family target union; Work origin admission; six forbidden public operations absent; ten mutation-based negative controls that all fire. This is a genuinely strong bundle-level semantic proof, not string theater.

What it does not establish: the Decision-pairing invariant (F-1); prose↔wire operation-name consistency in P9 (F-2); and — inherently, correctly — any runtime behavior (eligibility derivation, validity Q, replay ordering, revalidation sequencing), which remains D7+ with `active runtime baseline = NONE`. The B110/P9 structural verifiers are marker checks over reviewed artifacts; acceptable as lock-preservation guards given the OAD verifiers carry the semantic weight, with F-2 as the one cheap cross-check worth adding.

### 13. YAGNI / reconstruction decision

**No broader reconstruction and no generic approval/workflow machinery is justified before implementation.** Every generic capability (stages, quorum, expiry, comments, request platform, entity graph) lacks a present consumer and has an explicit reopen trigger. Conversely no further shrinkage survives: each element of C is load-bearing against a stated present law. The redesign solved the root cause of P9-F1 (missing pre-decision identity) rather than relocating it — the superseded query-only remedy would have preserved it inside the decision wire, and its rejection was correct.

### 14. D7 obligations

Valid future obligations, not architecture defects: (1) F-5 replay-before-precondition evaluation order; (2) F-4 client handling of bodyless 503 as ambiguous acceptance; (3) eligibility-derivation freshness mechanics (cache/awaken) under the "delayed events never sole revocation authority" law; (4) validity-Q timeout/unavailability realization mapping to the 503 semantics; (5) zero-decider Work materialization/reconciliation jobs; (6) F13/F14 delivery + currentness-revalidation mechanics; (7) invalidation-reconciliation sweep preventing silent stale-PENDING; (8) Idempotency-Key storage/scope/retention. None of these hides an unmade D0–D6 semantic decision — I checked each against the accepted authorities and the deciding semantics exist for all eight.

### 15. Continuation recommendation

Smallest exact action after GPT/operator adjudication of F-1..F-5:

1. adjudicate F-3 and F-4 (each is one authority sentence or a bounded field removal — operator chooses exposure-deferral vs removal);
2. apply F-1 and F-2 as repository proof fixes on the candidate branch under the normal RED→GREEN discipline (F-1: pairing assertion + 11th negative control; F-2: op-home names asserted against the bundle);
3. record F-5 plus the §14 list in the D7-R obligation register;
4. only then declare the AuthorizationRequest redesign Global-Maximum closed and open D7-R per roadmap.

This review authorizes nothing: no merge of PR #61, no D7-R/D8-R opening, no B10 resumption, no Product implementation. Findings are evidence pending adjudication against current repository authority.
