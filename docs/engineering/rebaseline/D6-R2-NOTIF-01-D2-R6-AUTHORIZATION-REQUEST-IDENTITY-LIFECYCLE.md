# NOTIF-01 D2-R6 — Authorization Request Identity & Lifecycle

> **Status:** OPEN — EXACT D2-R6 CANDIDATE / OPERATOR ADJUDICATION REQUIRED
> **Operator-approved premise:** a pre-decision authorization episode is a distinct canonical Governance identity named `AuthorizationRequest`, separate from `AuthorizationDecision`.
> **Parent identity authority:** [D2 Identity / Tenant / Data Ownership](D2-IDENTITY-TENANT-DATA-OWNERSHIP.md)
> **Boundary revalidation:** [D1-R2 Authorization-Request Boundary Revalidation](D6-R2-NOTIF-01-D1-R2-AUTHORIZATION-REQUEST-BOUNDARY-REVALIDATION.md) — CURRENT STRUCTURE CONFIRMED
> **Trigger / supersession:** [P9-F1 Supersession](D6-R2-NOTIF-01-D6-R-P9-F1-SUPERSESSION-RATIFICATION.md)
> **Canonical Product wire:** unchanged — 104 Product operations · 31 ordinary Permissions · H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Why D2 reopens

D2 currently gives stable MPC-owned identity to the terminal `AuthorizationDecision` occurrence but not to the materially real pre-decision authorization case.

P9 proved this gap operationally. The deeper architecture review then proved that `target + target revision` cannot substitute for request identity because materially distinct authorization episodes may reference the same target revision.

Example:

```text
PriceIntent PI-100 / revision 7
→ authorization episode A
→ authorized
→ later material governing-context drift makes that authorization insufficient for execution

PriceIntent PI-100 / revision 7
→ authorization episode B
→ new human authorization required
```

Episode A and B are different historical authorization meanings even though the action-owning target/revision may be unchanged.

D3's accepted reopen law says that if correctness requires a materially new business identity rather than an occurrence discriminator, perform a targeted D2 identity review. That condition is now met.

## 2. Root cause

> **The accepted model represents the business intent and the terminal authorization decision, but lacks canonical identity for the Governance-owned interval in which one concrete authorization episode exists, may be actionable by changing human Principals, may terminate without a decision, and must remain historically distinguishable from later reauthorization of the same target.**

An API query over current targets would expose data but would not fix this missing identity/lifecycle.

## 3. Target invariant

For every material authorization episode that enters Governance:

```text
one stable AuthorizationRequest identity
→ one immutable reviewed authorization basis/context
→ current decision eligibility remains Governance-owned and revocable
→ material current validity is revalidated semantically, not inferred from generic resource drift
→ zero or one terminal AuthorizationDecision
→ terminal history is never rewritten into another episode
→ later reauthorization creates a new AuthorizationRequest identity
```

The request never becomes the action-owning Business Intent or execution authority.

## 4. Canonical name

Candidate comparison:

- `ApprovalRequest` — rejected: too human-UI-specific and conflates approval with the broader authorization boundary.
- `AuthorizationCase` — rejected: unnecessarily case-management/generic-workflow flavored.
- `PendingAuthorization` — rejected: bakes one lifecycle state into identity.
- **`AuthorizationRequest` — selected.** It states exactly that an action-owning domain has asked Governance for authorization under a concrete bounded context.

`Request` here is business/Governance semantics, not HTTP transport.

## 5. Canonical `AuthorizationRequestID`

`AuthorizationRequest` is MPC-owned canonical state under **Controlled Action Governance** with one stable opaque `AuthorizationRequestID`.

Identity laws:

- opaque, non-business-semantic and non-reusable;
- does not encode Organization, target kind/ID, Principal, Permission, outcome, timestamp or revision;
- is distinct from the governed Business Intent identity;
- is distinct from `AuthorizationDecisionID`;
- remains stable while decision eligibility changes;
- once terminal, is never reopened/reused for a later authorization episode;
- a new authorization episode for the same target/revision gets a new ID.

No generic platform `RequestID`, `CaseID`, workflow identity or universal entity identity is admitted.

## 6. Closed governed-target binding

Launch-V1 `AuthorizationRequest` targets remain exactly the already accepted Governance target union:

```text
ListingIntent
PriceIntent
BusinessOrderIntent
InvoicingIntent
```

The request preserves a closed target identity and enough immutable source-owner-derived provenance to identify the reviewed business action/scope for this authorization episode.

A source-owner resource revision/ETag may be retained as **evidence of what was observed**, but it is not `AuthorizationRequest` identity and is not automatically the semantic validity predicate for authorization.

Why:

```text
irrelevant source-resource field changes
→ may change a whole-resource revision/ETag
→ must not automatically force reapproval
```

D0 explicitly requires materiality-based authorization validity. Therefore target-resource revision and authorization-basis validity are separate meanings.

No arbitrary `{entity_type, entity_id}` graph is introduced.

## 7. Authorization review/basis snapshot

A human must be able to understand what exact material action/context is being authorized, and later history must explain what was reviewed without reconstructing mutable current state.

Therefore each request retains one **immutable, bounded, typed authorization review/basis snapshot**.

Laws:

1. the action owner remains authority for source business meaning used to author the snapshot;
2. Governance owns retention of the authorization-purpose snapshot after accepting the request;
3. the snapshot is typed by the closed governed-target/action family; no generic payload/metadata/JSON bag;
4. it contains only information materially needed to understand/authorize the request under accepted action semantics;
5. it may preserve policy/evidence/scope provenance when material to why authorization was required, without making Governance owner of those underlying policies/facts;
6. it is not a general target read projection and grants no source-owner read capability;
7. it never becomes current source truth after creation;
8. if materially governing context changes enough that the reviewed basis is no longer valid, the request is invalidated or a new authorization episode is required rather than silently rewriting the snapshot;
9. non-material presentation/resource drift does not by itself require a new request.

Exact field schemas are deliberately not D2 authority. D3/D5/D6 must derive the minimum typed per-target review contracts after D2-R6 closes.

## 8. Request concurrency vs material authorization validity

This distinction is binding.

### 8.1 Request concurrency

`AuthorizationRequest` has an owner-local revision for stale-write/concurrent-decision protection.

Later D5 may encode this with standard conditional semantics such as an ETag/`If-Match` on the **AuthorizationRequest representation itself**.

This protects cases such as:

```text
João and Maria can both decide AR-501
→ both read request revision R3
→ João decides first, request becomes DECIDED / R4
→ Maria's stale R3 decision cannot create a second decision
```

D2 does not select the HTTP carrier.

### 8.2 Material authorization validity

A request's business target may change without every change being material to the authorization basis. Conversely, authority/policy/evidence context may change materially even when the target's own coarse revision does not.

Therefore:

> **A generic target-resource ETag is not the authorization validity oracle.**

Before a decision is committed, Governance must still revalidate:

```text
request is PENDING
+ current Principal is still authorized to decide
+ the preserved authorization episode is still materially valid enough to decide
```

When current source/business meaning is needed for that last predicate, Governance uses the already-accepted action-owner semantic boundary; the action owner evaluates material source/business validity under its own authority.

The exact Q/C/E communication and recovery contract is D3 feed-forward after D2-R6. Delayed/lost event delivery may not be the sole authority for current validity.

If material validity fails, the request does not silently retarget/rewrite itself; it becomes invalid or a later new authorization episode is created.

This preserves D0's two simultaneous laws:

```text
irrelevant drift → no needless reapproval
material governing drift → no stale authorization
```

## 9. Requester / initiator lineage

The request preserves exact accountable Principal lineage supplied by the action owner when materially available:

```text
requester_or_initiator_principal_id?
```

This is historical/correlation identity, not transfer of the Business Intent or actor policy to Governance.

It may be H/A/S according to the originating action. Human-only Notification result semantics may consume it only when the accepted Notification family requires an exact human recipient.

Email, username, role name and Permission are never requester identity.

## 10. Lifecycle — candidate

D2-R6 admits the smallest lifecycle required by known consumers.

### 10.1 `PENDING`

The request has been accepted by Governance and no terminal `AuthorizationDecision` exists.

`PENDING` does **not** mean:

- a currently eligible human necessarily exists;
- the request is visible to every holder of `governance.decide`;
- the target is executable;
- a Notification has been delivered;
- the request can be decided without current Governance + material-validity revalidation.

Current eligible decision Principals are dynamic Governance-owned meaning and may change while request identity remains stable.

### 10.2 `DECIDED`

Exactly one `AuthorizationDecision` has been committed for the request.

```text
AuthorizationRequest 1
→ 0..1 AuthorizationDecision
```

A baseline request never accumulates multiple votes/decisions. Multi-stage/quorum/voting semantics are not admitted.

The historical request remains `DECIDED` even if the resulting authorization later becomes insufficient for execution because materially governing context drifted. That later execution-validity fact does not rewrite the past decision.

If new human authorization is required, create a new `AuthorizationRequest`.

### 10.3 `INVALIDATED`

A pending request may terminate without a human decision when its authorization episode can no longer truthfully be decided under the preserved review/basis context.

Known invalidation classes are bounded to:

```text
target withdrawn / no longer seeks authorization
material source/governing-context invalidation
authorization-context invalidation
```

Routine changes in which individual Principals are currently eligible to decide do **not** automatically invalidate the request; Governance may recompute the current decision-Principal set.

Invalidation preserves enough typed provenance to explain why the request became terminal without fabricating a `reject` decision. Exact wire/event spelling belongs later.

An invalidated request is never reopened. A legitimate later authorization episode gets a new ID.

### 10.4 No baseline expiry state

No generic `expired` lifecycle is admitted merely by symmetry. Time-based authorization expiry may be added only if a governing product rule/consumer proves it.

## 11. Reauthorization lineage

Material reauthorization is a known consumer, so D2-R6 preserves a bounded optional predecessor link:

```text
predecessor_authorization_request_id?
```

It may be set only when the new request is explicitly a successor authorization episode for the same governed business action/target lineage.

It is not a generic related-entity graph and does not authorize traversal to arbitrary requests.

This allows history to distinguish:

```text
initial authorization request
→ decision
→ later material invalidation for execution
→ explicit reauthorization request
```

without inferring the relationship from timestamps/titles.

## 12. Duplicate intake / same-vs-different law

Retry/recovery of the **same semantic authorization episode** must resolve to the same `AuthorizationRequest`; it must not mint duplicates because transport was retried.

A **new** request is required when Governance/action-owner semantics establish a genuinely new authorization episode, including reauthorization after the prior episode no longer suffices.

Exact idempotency carrier, request anchor, persistence and ambiguous-intake recovery are D3/D5/D7 work. D2 freezes only the semantic same-vs-different property.

## 13. `AuthorizationDecision` feed-forward

Existing `AuthorizationDecision` remains canonical and distinct.

D2-R6 changes its lineage requirement conceptually to:

```text
AuthorizationDecision
→ belongs to exactly one AuthorizationRequest
→ preserves deciding Principal
→ preserves authorize | reject outcome
→ preserves exact authorization/authority context required for historical explanation
→ preserves the immutable reviewed target/scope/basis from that request episode
```

A Decision does not become a mutable `pending` object.

The client should not need to reconstruct a current target-resource ETag merely to identify what was reviewed; later D5 derives the exact command contract from `AuthorizationRequest` identity/revision and the D3 current-validity revalidation contract.

Later revocation/reapproval/execution invalidation never edits the historical decision outcome.

## 14. Notification feed-forward

This D2 repair changes only source identity semantics needed for later D3/D5/D6 derivation.

Candidate consequence:

```text
AUTHORIZATION_ACTION_REQUIRED
→ AuthorizationRequestRef
```

because the human is being told that **this exact pending authorization episode** needs a decision.

`AUTHORIZATION_DECISION_RESULT` does not automatically move to `AuthorizationRequestRef`; its accepted requester continuation may remain target-oriented so a requester is not forced to gain Governance-history access.

No Notification stores the request review snapshot, target current state, concurrency validator, Permission set or authority context.

## 15. Authority fences

`AuthorizationRequest` does not own or grant:

```text
business action disposition
Business Intent lifecycle
target-owner current state
ordinary Product Permission
source-owner read access
automatic execution
execution-time validity
Work assignment
Notification routing
provider/business-system protocol
```

Possession of `AuthorizationRequestID` grants nothing.

## 16. YAGNI / excluded structure

D2-R6 does **not** admit:

```text
workflow engine
BPMN
arbitrary approval stages
quorum/voting
parallel approval rules
custom request types
request form designer
generic case management
universal governed entity graph
free-form payload/metadata
request-level comments/chat
approval SLA/expiry
```

If a future named consumer proves one of those meanings, reopen the smallest responsible authority then.

## 17. Adversarial checks

The candidate fails if:

- two materially distinct authorization episodes for the same target/revision collapse into one identity;
- a terminal request is reopened instead of creating a new request;
- `AuthorizationDecision` is repurposed as both pending request and terminal decision;
- current approver eligibility becomes immutable recipient identity on the request;
- request snapshot becomes source-owner current truth or generic target mirror;
- request identity grants `governance.read`, target-owner read/write or execution authority;
- every source-resource ETag change forces reapproval regardless of materiality;
- material governing drift is ignored merely because a target-resource revision did not change;
- every delegation/member change invalidates the request by default;
- a human `reject` is fabricated merely because the request became invalid before decision;
- reauthorization lineage requires a generic relationship graph;
- D2 selects HTTP paths, operation count, database schema or runtime topology.

## 18. Candidate decision

**Recommendation:** ACCEPT D2-R6 with the model above.

The repair introduces exactly one missing canonical business identity under the already-correct Governance boundary and one bounded predecessor relation justified by the known reauthorization case. It also separates request concurrency from material authorization validity so we preserve both no-stale-authorization and no-needless-reapproval laws.

This is essential authorization complexity, not a generic workflow abstraction.

## 19. Gate

```text
D1-R2 boundary revalidation       PASS / CURRENT STRUCTURE CONFIRMED
D2-R6 exact identity/lifecycle    OPERATOR ADJUDICATION REQUIRED
D3 feed-forward                   BLOCKED BY D2-R6
D5/OAD                            104/31 UNCHANGED / BLOCKED
D6 P9                             PAUSED / BLOCKED BY D2-R6→D3→D5
D7-R / D8-R                       BLOCKED
Product implementation            BLOCKED UNTIL D9
```

**Exact next action:** operator adjudicates this D2-R6 contract. No D3/D5/OpenAPI modification before that adjudication.