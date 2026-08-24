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
