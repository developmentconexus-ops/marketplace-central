# NOTIF-01 D3-R3 — Authorization Request Communication & Recovery Feed-Forward

> **Status:** OPEN — EXACT D3-R3 CANDIDATE / OPERATOR ADJUDICATION REQUIRED
> **Parent communication authority:** [D3 Communication / Events](D3-COMMUNICATION-EVENTS.md)
> **Accepted identity authority:** [D2-R6 Ratification](D6-R2-NOTIF-01-D2-R6-RATIFICATION.md)
> **Boundary authority:** [D1-R2](D6-R2-NOTIF-01-D1-R2-AUTHORIZATION-REQUEST-BOUNDARY-REVALIDATION.md) — CURRENT STRUCTURE CONFIRMED
> **Trigger:** D6 P9 falsifier / AuthorizationRequest Global-Maximum repair
> **Canonical Product wire:** unchanged — 104 Product operations · 31 ordinary Permissions · H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Purpose

D3-R3 answers only:

> **How does one accepted `AuthorizationRequest` episode move between the action owner, Controlled Action Governance, Personal Notifications and Operational Work without duplicate request creation, stale decision authority, silent stalls or transferred ownership?**

D3-R3 does not select HTTP routes, OpenAPI schemas, operation counts, queues, transactions, database tables, router/server implementation or frontend components.

## 2. Imported invariants

The following are already fixed and are not reopened here:

- the action-owning domain owns Business Intent, intended target, effective action disposition and execution-time validity;
- Controlled Action Governance owns `AuthorizationRequest`, Grant/Delegation/current decision eligibility and `AuthorizationDecision`;
- Personal Notifications owns personal awareness only;
- Operational Work owns responsibility/assignment/escalation/work lifecycle only;
- IdentityAccess owns current Principal access eligibility and Organization Membership/ordinary access facts;
- one request episode is not target identity and not decision identity;
- request concurrency is not generic target-resource validity;
- `PENDING | DECIDED | INVALIDATED` is the exact baseline request lifecycle;
- a request has at most one terminal decision;
- later reauthorization creates a new `AuthorizationRequestID`.

No new D1 semantic edge is required. D3-R3 uses only already-accepted action-owner ⇄ Governance, domain → Work, Governance → Personal Notifications and current identity/access Q boundaries.

## 3. Global communication graph

```text
ACTION OWNER
  owns Business Intent + approval-required disposition
      |
      | C  request one authorization episode
      v
GOVERNANCE
  owns AuthorizationRequest(PENDING)
      |
      | E  current human decision-eligibility episode
      v
PERSONAL NOTIFICATIONS
  owns F13 personal awareness only

HUMAN DECIDER
      |
      | C  decide request using request revision
      v
GOVERNANCE
      |
      | Q  material authorization-basis validity when current source meaning matters
      v
ACTION OWNER
      |
      | valid / invalid / unknown-or-unavailable
      v
GOVERNANCE
  DECIDED + AuthorizationDecision
  or INVALIDATED
      |
      +---- E DecisionCommitted ------> ACTION OWNER
      |                                  revalidates execution-time validity
      |
      +---- E DecisionCommitted ------> PERSONAL NOTIFICATIONS
      |                                  optional F14 to exact human requester
      |
      +---- E RequestInvalidated -----> ACTION OWNER

PENDING request + zero eligible human decision principals
      |
      | E material Governance actionable condition
      v
OPERATIONAL WORK
  owns explicit responsibility/work lifecycle; no fallback approver
```

Transport technology and transaction materialization remain D7.

## 4. Creating one `AuthorizationRequest` — C + semantic anchor

When an action owner has a durable Business Intent whose current owner disposition requires one new authorization episode, the owner asks Governance to accept that episode through **C**.

The action owner supplies the closed governed target plus the immutable bounded review/basis material required by D2-R6 and one stable owner-local **authorization episode key**.

Semantic acceptance anchor:

```text
AuthorizationRequestAnchor
= Organization
+ governed target identity
+ action-owner-local authorization_episode_key
```

The `authorization_episode_key`:

- is stable for retry/recovery of the **same** semantic episode;
- is newly minted by the action owner when that owner concludes a **new** authorization episode is required;
- is scoped to the owning action/target lineage and is not a universal platform RequestID/CommandID/EventID;
- is only an occurrence discriminator/acceptance anchor; canonical request identity remains Governance-owned `AuthorizationRequestID`;
- cannot be reused with materially different request basis/scope. Same anchor + materially different request is an explicit conflict, not silent mutation.

Therefore:

```text
same episode / repeated C / timeout recovery
→ one AuthorizationRequest

same target, even same coarse target revision, but new reauthorization episode
→ new authorization_episode_key
→ new AuthorizationRequestID
```

This directly applies accepted D3 §4.11 ambiguous-acceptance law without inventing a generic command identity.

### 4.1 Capability outcome

Governance acceptance means Governance has durably accepted/recovered one canonical request episode. The resulting request may be `PENDING`; `PENDING` is request lifecycle, not evidence that transport acceptance was uncertain.

If caller observation of acceptance is ambiguous, the action owner reconciles with Governance by the semantic anchor before retrying. Retry of the same anchor converges on the existing request.

D5 later decides whether any Product operation exposes this capability; D3 only fixes owner-to-owner semantics.

## 5. Current decision eligibility is dynamic Governance truth

For a `PENDING` request, Governance resolves the exact current set of human Principals eligible to decide from:

```text
request authorization class/scope
+ current Governance Grant/Delegation semantics
+ current IdentityAccess Principal eligibility / Organization Membership as required
```

Binding laws:

- current eligibility is Q/current-state meaning, not a historical recipient list on the request;
- `governance.decide` alone is not the recipient selector;
- a PENDING request may have `0..N` current eligible humans;
- ordinary membership/Principal or delegation changes may change the set without changing request identity;
- merely adding/removing another eligible decider does **not** stale unrelated deciders by changing request revision;
- every decision attempt revalidates the exact current human's eligibility; stale UI/cached eligibility is never authority.

IdentityAccess current truth may be cached/awakened later by D7 mechanisms, but delayed access-change events can never be the sole revocation authority.

## 6. F13 `AUTHORIZATION_ACTION_REQUIRED` — Governance E → Notifications with currentness revalidation

When Governance establishes a **material current decision-responsibility episode** for exact human Principal `P` on `AuthorizationRequest R`, Governance owns the producer fact:

> `P` is currently an eligible decision Principal for pending request `R`.

That committed Governance fact may produce **E** to Personal Notifications.

The event contract carries only what Personal Notifications needs to react:

```text
explicit Organization
AuthorizationRequestRef
effective human Principal P
stable Governance-owned decision-eligibility episode discriminator
source occurrence time
```

No request review snapshot, Permission set, role data, target current state or authority internals are copied into Notification propagation.

### 6.1 Same vs distinct F13 occurrence

Governance owns a bounded `decision_eligibility_episode_key` or equivalent stable occurrence discriminator for the request+Principal responsibility episode.

```text
duplicate delivery of same eligibility episode
→ one semantic F13 Notification

P loses eligibility, then later gains a genuinely new decision-responsibility episode while R is still pending
→ distinct occurrence discriminator
→ a new F13 awareness occurrence may be created
```

Configuration flapping must not create arbitrary duplicate awareness; Governance owns the semantic episode boundary, not transport delivery count.

### 6.2 Late delivery must not create stale actionable awareness

F13 is progression/current-actionability awareness, not historical evidence accumulation.

Before Personal Notifications materializes a delayed/recovered F13, currentness is revalidated against Governance:

```text
request still PENDING
+ P still currently eligible to decide
```

If either predicate is false:

```text
late/replayed eligibility E
→ no new F13 Notification
```

Thus an event emitted before another decider commits a decision cannot later create a misleading “action required” notification after the request is `DECIDED`.

While `R` remains PENDING and `P` remains eligible, a missing required F13 reaction is recoverable from current Governance state plus Personal Notifications state. Event transport is not the sole recovery authority.

### 6.3 Existing historical F13 after source becomes non-actionable

A F13 Notification already created for P is historical personal awareness. Later `DECIDED`, `INVALIDATED` or loss of P's eligibility does **not** rewrite/delete/retarget that Notification.

Opening it later rechecks Governance current request/actionability and may truthfully say the request is no longer actionable. Personal Notifications still owns read/archive state only.

D2/D5 feed-forward changes F13 source identity to `AuthorizationRequestRef`; exact wire waits for D5.

## 7. Zero eligible human deciders — explicit condition, never fallback

A known-empty current decision Principal set is a materially different state from unavailable/unknown resolution.

```text
PENDING AuthorizationRequest
+ current eligible human set = known empty
→ authorization cannot progress through the required human path
```

Binding laws:

- never fall back to all admins, all `governance.decide` holders, Organization owner, request initiator or routing configuration;
- never fabricate a F13 recipient;
- unavailable eligibility resolution is not known-empty;
- a persistently blocking known-empty condition is a Governance-owned material actionable condition, not a silent queue state.

D0 requires material actionable conditions not to become ownerless. Governance therefore feeds the already-accepted domain → Operational Work boundary:

```text
Governance commits:
AuthorizationRequest R has no current eligible human decision Principal
    ↓ E
Operational Work
    decides whether/how to represent/deduplicate the explicit obligation
```

The Work subject is the sufficiently explicit Governance condition anchored by `AuthorizationRequestID`; this does not create another business identity or new Notification kind.

When the source condition resolves because at least one eligible human exists or the request becomes terminal, Governance commits the corresponding source-resolution fact so Work reconciles its own obligation. Work closure never mutates the request.

Exact Work projection/source-ref wire is a D5 feed-forward question; no new Product operation is assumed here.

## 8. Deciding a request — C with request-local concurrency + current semantic validity

A human decision asks Governance to decide exactly one `AuthorizationRequest`.

Semantic command input later exposed by D5 must be derivable from:

```text
AuthorizationRequestID
+ request owner-local revision / precondition
+ outcome = authorize | reject
```

It does **not** require the client to fetch/reconstruct a whole-target ETag merely to identify the reviewed request.

Before committing a decision, Governance revalidates in this order proportionately:

```text
1. request exists in exact Organization
2. request is PENDING
3. supplied request revision/precondition is current
4. caller is exact current H decision Principal
5. current IdentityAccess eligibility/Membership gates still hold
6. preserved authorization episode is still materially valid enough to decide
```

Step 6 uses **Q to the action owner** when current source/business meaning is required.

### 8.1 Action-owner validity Q

The action owner answers only the material predicate it owns:

```text
VALID
  preserved request basis is still materially valid enough for a decision

INVALID
  preserved episode can no longer truthfully be decided; current owner meaning requires invalidation/new episode/abandonment

UNKNOWN_OR_UNAVAILABLE
  owner cannot currently establish sufficient truth
```

Rules:

- `VALID` does not mean executable now or guarantee later execution;
- irrelevant target/resource drift may still be `VALID`;
- material business/policy/evidence/readiness drift may be `INVALID` even if a coarse target ETag is unchanged;
- `UNKNOWN_OR_UNAVAILABLE` never becomes VALID by default.

Outcome handling:

```text
VALID
→ Governance may commit exactly one AuthorizationDecision

INVALID
→ no human Decision is fabricated
→ Governance transitions request to INVALIDATED if still pending/current

UNKNOWN_OR_UNAVAILABLE
→ no Decision
→ request remains PENDING
→ explicit unavailable/unknown outcome to caller
```

Exact problem/status vocabulary is D5.

### 8.2 Concurrent deciders

Two humans may legitimately read the same PENDING request revision.

```text
João reads R3
Maria reads R3
João commits decision → request DECIDED / next revision
Maria attempts with R3 → stale / no second decision
```

The request-local revision protects one-request/one-decision. Changes to unrelated current eligible-decider membership do not by themselves create needless request revision churn; current eligibility is separately revalidated.

## 9. Decision propagation — Governance E → action owner and Personal Notifications

When Governance commits an `AuthorizationDecision`, that is an already-true Governance-owned occurrence and therefore **E** is correct for independent reactions.

The canonical `AuthorizationDecisionID` is sufficient same-vs-distinct occurrence identity.

### 9.1 Governance → action owner

Decision propagation to the action owner is **recoverable**.

Event meaning is bounded to:

```text
Organization
AuthorizationRequestRef
AuthorizationDecisionRef
authorize | reject
committed decision time
```

The event does not execute or mutate the Business Intent.

The action owner reacts under its own semantics:

- duplicate decision E is idempotent by canonical Decision identity;
- current Governance decision/context may be Q-revalidated when material;
- `authorize` only supplies authorization evidence/context;
- before any consequential execution, the action owner still performs its D0 execution-time validity revalidation;
- `reject` informs owner-local next state/disposition; Governance does not mutate that state.

Loss of the Decision E cannot permanently stall the action owner. Missing required reaction is recoverable from canonical Governance request/decision state plus action-owner pending state.

### 9.2 Governance → Personal Notifications F14

The same committed Decision may independently wake Personal Notifications for `AUTHORIZATION_DECISION_RESULT` when accepted requester/initiator lineage resolves an exact currently eligible human recipient.

Binding semantics remain:

```text
source occurrence identity = AuthorizationDecisionID
presentation outcome = authorize | reject
source_ref for requester continuation = AuthorizationTargetRef unless later D5/D6 proof requires a smaller correction
```

The requester is not forced to hold `governance.read` merely to navigate back to the original work. A/S requester lineage does not create a human Inbox by symmetry.

Duplicate Decision E does not duplicate F14.

## 10. Request invalidation

`INVALIDATED` terminates a pending request without inventing a human decision.

Two authority sources can prove invalidation.

### 10.1 Action-owner material invalidation — owner E → Governance

When the action owner commits a fact it owns such as:

```text
target/intent withdrawn
owner no longer seeks authorization
material reviewed business basis is no longer valid
```

that producer fact is already true independent of Governance reaction. Therefore **E** is the baseline cross-owner form.

Governance consumes it and Q-revalidates current owner meaning where necessary before transitioning the matching PENDING request to INVALIDATED.

The action owner does not directly write Governance state.

Propagation is recoverable: if the E is lost/delayed, Governance must be able to discover the still-material mismatch from current owner public meaning while the request remains pending; silent permanent stale-PENDING state is not acceptable.

### 10.2 Governance-owned authorization-context invalidation

Governance may independently know that its own preserved authorization episode/context has become invalid. This does **not** include routine changes in which individual humans are currently eligible; those normally recompute actionability without terminating the request.

A true authorization-context invalidation transitions the request inside Governance authority.

### 10.3 Governance invalidation → action owner

A committed `AuthorizationRequest INVALIDATED` occurrence is **E** to the action owner when the owner has a pending lifecycle consequence.

The action owner then decides under its own semantics whether to:

```text
abandon / remain blocked
recompute disposition/context
create a new authorization episode
```

No F14 decision-result Notification is created because no AuthorizationDecision exists.

Existing F13 awareness is not rewritten; current request read shows terminal non-actionable state.

## 11. Reauthorization

A new authorization episode is always an action-owner decision.

Typical path:

```text
prior request DECIDED
→ action owner later determines prior authorization is insufficient under material execution-time validity
→ owner establishes new authorization basis
→ new authorization_episode_key
→ C Governance
→ new AuthorizationRequestID
→ optional predecessor_authorization_request_id = prior request
```

A new episode never reopens or mutates the prior request/decision. Same target and same coarse target revision do not collapse the two requests.

Retry of the new episode reuses its new anchor; it does not fall back to the predecessor request.

## 12. Failure / replay / recovery matrix

| Failure class | Required semantic outcome |
| --- | --- |
| Action-owner C response lost after Governance accepted | Reconcile by Organization + target + owner-local authorization episode key; recover same request, no duplicate. |
| Same request C delivered twice | One `AuthorizationRequest`; same anchor + materially different basis fails explicitly. |
| F13 eligibility E duplicated | One Notification for same Governance eligibility episode. |
| F13 eligibility E delayed until request terminal | No new action-required Notification after current Governance revalidation. |
| P loses eligibility before F13 materialization | No new F13 for P; prior Notification history, if any, remains. |
| Missing F13 while R still pending/P eligible | Recover from current Governance actionability + Notification state. |
| Eligibility temporarily resolves unavailable | Never interpret as no eligible Principal and never invent recipient. |
| Eligible set known empty | Explicit Governance actionable condition → Operational Work; no admin fallback. |
| Two deciders race | First valid request-revision decision wins; second cannot create another Decision. |
| Whole target ETag changed for irrelevant reason | Does not automatically invalidate request. |
| Material authorization basis invalid despite unchanged target ETag | Action-owner validity Q returns INVALID → request INVALIDATED, no Decision. |
| Action-owner validity Q unavailable/unknown | No Decision; request remains PENDING with explicit uncertainty. |
| Decision E duplicated | Action-owner and F14 reactions idempotent by AuthorizationDecisionID. |
| Decision E lost | Recover required owner reaction from canonical request/decision + owner pending state. |
| Owner invalidation E delayed/lost | Governance reconciles current owner meaning; stale PENDING cannot persist silently. |
| Request INVALIDATED | No fake reject Decision; no F14; action owner decides next step. |
| Reauthorization of same target/revision | New authorization episode key + new AuthorizationRequestID; optional predecessor link. |

## 13. Projection/read implications

D3-R3 proves that the human approvals queue can later read **Governance-owned `AuthorizationRequest` current state + immutable typed review snapshot** without requiring broad target-owner read Permissions merely to understand what is being authorized.

That does not make Governance a general target read authority. The snapshot is authorization-purpose historical/current request context accepted at request creation, while material current validity remains checked through the action-owner boundary when a decision is attempted.

Completed-decision history remains `AuthorizationDecision` meaning; pending-actionable work remains `AuthorizationRequest` meaning.

Exact Product Q/P surface is D5.

## 14. No new Product census yet

D3-R3 creates **no Product operation or Permission** by itself.

Canonical wire remains:

```text
104 Product operations
31 ordinary Permissions
H / A / S
```

D5 must later derive, from the accepted D2/D3 semantics rather than from screen shape:

- the minimum Request read/list surface for human actionability and exact-request continuation;
- the exact decision command shape bound to `AuthorizationRequest` revision;
- any bounded internal/owner request-intake exposure actually required by a Product consumer;
- Notification F13 `AuthorizationRequestRef` source union correction;
- any Work representation/source-ref feed-forward required for the no-eligible-decider condition;
- historical/non-regression proof and the resulting operation census as a consequence.

No endpoint name/count is frozen in D3.

## 15. Negative controls / falsifiers

D3-R3 fails if any realization:

1. retries one semantic episode into multiple AuthorizationRequests;
2. uses `AuthorizationRequestID` as a generic client-generated CommandID;
3. treats target ETag equality as sufficient authorization validity or any target ETag drift as automatic invalidation;
4. lets action owner directly write Governance request lifecycle;
5. lets Governance mutate Business Intent or execute the governed action;
6. lets Notification contain request review snapshot/current source truth or act as authorization token;
7. materializes F13 from a late event after request/recipient actionability ended;
8. infers F13 recipients from Permission/admin/default role fallback;
9. treats known-empty eligible set as success/healthy/no-action;
10. leaves a materially blocking zero-decider request silently ownerless instead of explicit Work;
11. allows a second Decision for the same request;
12. fabricates human reject when the request becomes invalid without a decision;
13. treats action-owner validity Q unavailable/unknown as valid;
14. executes automatically merely because an authorize Decision exists;
15. rewrites a prior request/decision for later reauthorization;
16. creates a generic workflow/saga/event-bus authority or universal entity/correlation graph;
17. makes event transport the sole recovery authority for request, actionability, decision or invalidation progression.

## 16. Decision

**Recommendation:** ACCEPT D3-R3.

The selected topology is the smallest sustainable communication model that makes D2-R6 operationally complete:

- C with a domain-local semantic episode anchor for exactly-once **business meaning without exactly-once transport claims**;
- Q for current decision eligibility/material owner validity where currentness matters;
- E for committed Governance/action-owner occurrences that require independent reaction;
- recoverable current-state progression rather than event-log authority;
- Operational Work only for the genuinely material no-eligible-decider stall;
- no new domain, broker, workflow engine, generic command identity or Product operation assumed in advance.

## 17. Closure-review requirement

D3-R3 is **not** the final Global-Maximum closure of this redesign.

After D3-R3 is ratified and D5 wire + final D6 P9 trace are derived/proved, the full AuthorizationRequest redesign must undergo an **independent Fable review** before D7-R opens. The review must challenge authority, lifecycle, reauthorization, stale/current validity, concurrency, zero-decider Work handling, Notifications, least privilege, D5 surface and P9 operability. Findings must be adjudicated; review existence alone is not proof.

## 18. Gate

```text
D1-R2                         PASS / CURRENT STRUCTURE CONFIRMED
D2-R6                         OPERATOR-RATIFIED / ACCEPTED
D3-R3                         OPERATOR ADJUDICATION REQUIRED
D5/OAD                        104/31 UNCHANGED / BLOCKED BY D3-R3
D6 P9                         PAUSED / BLOCKED BY D3-R3→D5
Independent Fable review      REQUIRED AFTER D5+P9, BEFORE GLOBAL-MAXIMUM CLOSURE / D7-R
D7-R / D8-R                   BLOCKED
Product implementation        BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates this D3-R3 contract. Do not modify D5/OpenAPI, resume P9/B10, open D7-R/D8-R or implement Product code before this gate closes.