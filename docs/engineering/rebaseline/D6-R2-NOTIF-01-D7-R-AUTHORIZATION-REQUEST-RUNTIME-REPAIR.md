# NOTIF-01 D7-R — AuthorizationRequest Runtime / Jobs / Transactions Repair

> **Status:** OPERATOR-RATIFIED / ACCEPTED
> **Operator ratification:** 2026-08-24 — `Aprovado`
> **Trigger:** [AuthorizationRequest Fable Closure Ratification](D6-R2-AUTHORIZATION-REQUEST-FABLE-RATIFICATION.md)
> **Runtime owner:** [D7 Runtime / Jobs / Transactions](D7-RUNTIME-JOBS-TRANSACTIONS.md)
> **Bounded mechanisms:** [D7-B PostgreSQL](D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md) + [D7-C Durable Work](D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md) + [D7-E HTTP / Proof](D7-E-OPERABILITY-DEPLOYMENT-PROOF.md)
> **Canonical Product:** 106 Product operations · 31 ordinary Permissions · H/A/S
> **Active runtime:** NONE
> **Implementation:** BLOCKED UNTIL accepted D9

```text
D7R_GLOBAL_MAXIMUM:ACCEPTED
D7R_PROOF_CLASS:STRUCTURAL_AUTHORITY_ONLY_NOT_RUNTIME_PROOF
```

## 1. Purpose and boundary

D7-R realizes the accepted `AuthorizationRequest` semantics against the already-ratified D7 runtime architecture. It does not redesign the Product, create a second runtime stack or begin Product implementation.

The Fable closure carries eight runtime obligations into D7-R:

1. committed same-key replay resolves before stale `If-Match` evaluation;
2. typed semantic validity-unavailable 503 remains distinguishable from ambiguous/non-semantic 503;
3. eligible-decider derivation stays current without event/cache authority;
4. material-validity Q timeout/unavailability realizes the typed 503 contract;
5. zero-decider Work is materialized and reconciled recoverably;
6. F13/F14 delivery revalidates currentness and remains duplicate-safe;
7. request invalidation cannot be permanently lost when its wake-up event is missed;
8. `Idempotency-Key` storage/scope/retention preserves exact replay semantics.

D7-R adds **no Product operation, Permission, Principal kind, business owner, NotificationKind, Work capability, generic workflow engine, external broker, Redis, second outbox, new database, new server framework or new source of business truth**.

## 2. D7 revalidation result

The accepted D7 mechanisms remain sufficient:

```text
HTTP / strict wire       Chi v5 + oapi-codegen strict server + OAD middleware
owner transactions       pgx/v5 + PostgreSQL READ COMMITTED + explicit locks
Organization isolation   transaction-local scope + FORCE RLS + composite FKs
generic idempotency      D7-B owner-local technical records
post-commit reactions    River InsertTx
recovery                 owner Q/revalidation + reconciliation/sweeps
```

D7-B already requires exact prior idempotent intake to be resolved before stale revision re-evaluation. D7-C already requires transactionally durable consumer jobs, semantic duplicate safety, no exactly-once claim, and recovery sweeps that do not make scheduler/event delivery business authority. D7-E already requires Product Problems to be emitted from canonical Product wire semantics rather than middleware error leakage.

Therefore this is a **bounded D7 composition repair**, not D7 reconstruction.

## 3. Alternatives challenged and adjudicated

### A — reuse accepted D7 primitives with AuthorizationRequest-specific sequencing

**SELECTED / OPERATOR-RATIFIED — GLOBAL MAXIMUM.** It closes every current runtime obligation while adding no infrastructure or business authority.

### B — generic approval/workflow engine

**REJECTED.** No current quorum, stage, SLA, voting, arbitrary case or workflow consumer exists. It would recreate the broader structure already rejected by D2/D5/Fable.

### C — synchronous-only propagation with no durable work

**REJECTED.** A committed Decision/invalidation/eligibility occurrence could lose required action-owner, F13/F14 or Work propagation on process/network failure.

### D — new event broker / generic outbox beside River

**REJECTED.** Accepted PostgreSQL + River `InsertTx` already gives atomic owner-state → durable-reaction handoff. Another broker/outbox adds a relay and failure boundary without a consumer.

## 4. `CreateAuthorizationDecision` exact runtime order

The decision command has one authoritative sequence.

### 4.1 Transport/access envelope

Before the owner transaction:

```text
H session/AuthN
→ OAD request validation
→ exact Organization Membership / current ordinary governance.decide access gate
→ canonical request decoding
```

These gates are not idempotency/business replay. Current Product access may still fail closed even when an old idempotent result exists.

### 4.2 Principal-scoped idempotency namespace

The generic D7-B scope is specialized for `CreateAuthorizationDecision` because the immutable result carries the server-derived deciding Principal and a raw client key must never become cross-Principal replay/disclosure authority.

The accepted uniqueness namespace is:

```text
organization_id
+ effective PrincipalID
+ CreateAuthorizationDecision operation identity
+ digest(Idempotency-Key)
```

```text
D7R_IDEMPOTENCY_SCOPE:ORGANIZATION_PRINCIPAL_OPERATION_KEY
```

The effective Principal is derived from trusted authentication/current Product access context. It is **not** a client field, credential digest or caller-authored authority.

Within that namespace, the semantic fingerprint includes at least:

```text
authorization_request_id
+ supplied If-Match
+ outcome
```

```text
D7R_IDEMPOTENCY_FINGERPRINT:REQUEST_IFMATCH_OUTCOME
```

A same raw Idempotency-Key under a different effective Principal is a different idempotency namespace. It neither replays nor conflicts with the first Principal's intake merely because the raw key string is equal. The second Principal proceeds as a genuinely new decision attempt, where current request state/revision/eligibility decides admissibility.

> **Binding law:** same raw Idempotency-Key under a different effective Principal is a different idempotency namespace.

This specialization does not change the Product OAD or add a public identity parameter. It is a bounded D7 realization of server-attributed actor semantics.

### 4.3 Idempotency before request concurrency

Inside one exact-Organization Governance transaction:

```text
lookup/claim Principal-scoped Idempotency-Key
  ├─ existing same fingerprint + committed terminal result
  │    → replay the established result
  │    → BEFORE request lock / If-Match / current eligibility/material-validity re-evaluation
  ├─ existing different fingerprint
  │    → accepted reused-key failure
  ├─ existing genuinely in-progress intake
  │    → accepted in-progress conflict
  └─ new intake
       → continue to request/currentness evaluation
```

```text
D7R_REPLAY_ORDER:IDEMPOTENCY_BEFORE_IF_MATCH
```

An exact replay is recovery of an already committed historical command, not a new decision attempt. It still passes current AuthN/ordinary Product access gates, but it does not fabricate a second current eligibility/material-validity decision.

### 4.4 New decision attempt

Only for a genuinely new intake:

```text
lock exact AuthorizationRequest under Organization scope
→ request must still be PENDING
→ compare request-local If-Match
→ revalidate exact current human decision eligibility
→ Q action owner for material authorization-basis validity
→ commit one terminal meaning or no-effect response
```

`If-Match` is request concurrency only. It is never material-validity authority.

## 5. Current decision eligibility

Actionability list/detail and every **new** decision attempt derive eligibility from current Governance authority/delegation plus current IdentityAccess truth.

```text
no durable eligibility cache as authority
no Notification-derived eligibility
no River-event-derived eligibility
events/wakeups may optimize reevaluation only
```

```text
D7R_ELIGIBILITY:AUTHORITATIVE_Q_NOT_EVENT_CACHE
```

For the decision command, the authoritative eligibility Q occurs after the request has been locked/current-checked and before Decision creation. A revocation/change committed before that Q must be observed. A change after the successful Q is later authority drift; it does not rewrite a historical Decision, and action-owner execution-time revalidation remains binding.

D7-R does not acquire private cross-owner IdentityAccess locks or create a distributed transaction merely to freeze future access state.

## 6. Material-validity Q and transaction safety

Governance never calls a marketplace/provider/business-system network endpoint while its owner transaction/request lock is open.

The decision-time material-validity Q is an in-process action-owner query over that owner's current MPC truth/evidence. If sufficient current evidence is not available without an external acquisition/reconciliation step, the action owner answers `UNKNOWN_OR_UNAVAILABLE` rather than performing hidden network I/O or guessing `VALID`.

```text
D7R_VALIDITY_Q:IN_PROCESS_NO_EXTERNAL_NETWORK
```

The three outcomes are realized as follows.

### `VALID`

In the Governance transaction:

```text
create immutable AuthorizationDecision
+ set request DECIDED / rotate request revision
+ persist terminal Principal-scoped idempotency result reference
+ InsertTx independent required reaction jobs
COMMIT
```

Required Decision consumers get independent durable jobs, including action-owner propagation and F14 when exact accepted requester lineage permits it.

### `INVALID`

Governance atomically terminates the still-current request as `INVALIDATED`, rotates its revision, persists the terminal Principal-scoped idempotent command outcome, and inserts required owner/Work reconciliation reactions. It creates **no AuthorizationDecision and no F14**.

### `UNKNOWN_OR_UNAVAILABLE`

No Decision and no request lifecycle mutation are committed. A newly claimed idempotency intake is not retained as a false effect. The strict Product responder emits only the canonical semantic Problem:

```text
https://conexus.fun/marketplace-central/problems/product/authorization-validity-unavailable
```

The request stays PENDING.

```text
D7R_503:EXACT_TYPED_KNOWN_NO_EFFECT
```

Any bodyless, unparsable, proxy/infrastructure or otherwise non-matching 503 is not this semantic result and remains **ambiguous potentially accepted** to the client. Runtime middleware/proxy errors may not counterfeit the typed Product Problem.

## 7. Idempotency result storage and retention

For committed `CreateAuthorizationDecision` terminal effects, the idempotency record stores only the minimum durable replay correlation, proportionately:

```text
effective Principal namespace identity
+ key digest
+ semantic fingerprint digest
+ stable result/status reference
```

It does not duplicate the immutable review basis or raw request body merely for replay.

For a committed Decision, replay can resolve through the immutable `AuthorizationDecisionID`. For an attempt that atomically caused request invalidation, the terminal request/outcome correlation remains durable enough to replay the established result semantics.

Because Product authority defines no shorter safe ambiguous-retry horizon, the committed terminal mapping for this command is retained **no shorter than the corresponding AuthorizationRequest/AuthorizationDecision historical authority**. This is operation-bounded; it does not impose perpetual retention on every MPC idempotency record.

```text
D7R_TERMINAL_REPLAY_RETENTION:AUTHORIZATION_HISTORY_LIFETIME
```

Transient rolled-back/no-effect attempts do not gain historical idempotency state merely by arrival.

## 8. F13 action-required materialization

A current eligible-human episode may awaken PersonalNotifications through one durable River reaction. Job args are bounded to technical routing/correlation such as Organization, AuthorizationRequest reference, exact human Principal, occurrence discriminator and source time; no review basis, Permission set, role state or provider payload is copied into River.

Before PersonalNotifications materializes F13 it re-Qs Governance/current eligibility:

```text
request still PENDING?
+ exact human still eligible now?
```

If either is false, the delayed/replayed job completes without creating new actionable awareness.

Duplicate delivery of the **same eligibility occurrence** is semantically idempotent. If a human becomes ineligible and later legitimately eligible again while the request remains PENDING, Governance may produce a new eligibility occurrence/discriminator and therefore a new bounded F13 occurrence.

```text
D7R_F13:REVALIDATE_PENDING_AND_ELIGIBLE
```

Existing historical Notifications are never rewritten.

## 9. F14 decision-result materialization

A committed Decision inserts the F14 consumer reaction independently from action-owner propagation when exact accepted human requester lineage permits F14.

The semantic duplicate key is anchored by the immutable Decision occurrence + exact recipient; duplicate/rescued jobs cannot create duplicate result Notifications.

F14 remains target-oriented. River carries no source-read authority and PersonalNotifications does not require `governance.read` to notify the requester.

```text
D7R_F14:DECISION_OCCURRENCE_IDEMPOTENT
```

## 10. Zero-decider Work lifecycle

Known-empty current eligible-human set is detected at request creation/current evaluation and is never filled by a fallback admin/Permission holder.

Governance feeds the existing Operational Work owner through durable reaction semantics. Work's semantic obligation is anchored by the AuthorizationRequest origin + the zero-decider condition so duplicate reevaluations do not create duplicate Work meaning.

```text
PENDING + known-empty deciders
→ ensure zero-decider Work obligation
```

When a valid decider exists again **or the request becomes terminal**, Governance/Work reconciliation closes the no-longer-applicable Work obligation. Work closure never mutates or declares the AuthorizationRequest truth by itself.

Assignment/escalation of that Work never grants `governance.decide` or current Governance eligibility.

```text
D7R_ZERO_DECIDER:WORK_MATERIALIZE_AND_RECONCILE
```

## 11. Reevaluation wakeups and recovery sweep

Correctness cannot depend on every eligibility/invalidation wake-up being delivered.

D7-R uses two complementary lanes:

```text
fast lane
  accepted local state changes / owner events may enqueue reevaluation work

recovery lane
  bounded River scheduled/recovery sweep scans durable PENDING Governance truth
  → re-Q current eligibility
  → re-Q material validity where due/required
  → reconcile F13 / zero-decider Work / invalidation obligations
```

The sweep cursor/tick is technical state only. Exact cadence/queue names/concurrency are not frozen without evidence.

A missed event or scheduler tick cannot permanently strand a correctness-required reaction; later sweep/startup recovery rediscovers it from durable owner state.

```text
D7R_RECOVERY:EVENT_PLUS_DURABLE_PENDING_SWEEP
```

## 12. Invalidation recovery

Action-owner material invalidation may awaken Governance via accepted E, but that E is never sole authority.

```text
exact Organization + request
→ re-Q current Governance request
→ if terminal: no-op/reconcile dependents
→ if PENDING: Q action-owner material validity
   ├─ VALID                 → leave PENDING
   ├─ INVALID               → Governance INVALIDATED transaction + required reactions
   └─ UNKNOWN/UNAVAILABLE   → leave PENDING + safe read/recovery work
```

An INVALIDATED transaction creates no Decision/F14 and reconciles zero-decider Work. Duplicate/out-of-order invalidation jobs are harmless because request state is current owner truth.

```text
D7R_INVALIDATION:OWNER_EVENT_NOT_SOLE_AUTHORITY
```

## 13. Product HTTP / Problem realization

D7-E's generated strict server + Product validator remains binding.

Only an explicit application result representing material validity `UNKNOWN_OR_UNAVAILABLE` may construct the exact `authorization-validity-unavailable` Problem body. Generic middleware, dependency outage, panic/error recovery, reverse proxy or transport failure does not map to that semantic type by status-code symmetry.

```text
exact typed Problem 503  → known no Decision / request remains PENDING
any other uncertain 503  → ambiguous potentially accepted
```

No new Product response/status or operation is required.

## 14. Required real-dependency falsifiers

D7-R extends the accepted D7 proof contract. Later implemented proof must be capable of falsifying at least:

1. committed Decision + lost 201 + exact same Principal/key/command returns 412 instead of replaying the original Decision;
2. same raw key under a different effective Principal collides with or replays the first Principal's idempotency record instead of using an independent namespace;
3. same Principal + same key + changed request/`If-Match`/outcome avoids the accepted reused-key failure;
4. two new different-key deciders consume the same request revision and both create Decisions;
5. exact-current eligibility is replaced by stale cache/Notification/event state;
6. a provider/business-system network call occurs while the Governance decision transaction is holding its request lock;
7. exact typed semantic 503 leaves a Decision/request mutation or durable false-success intake behind;
8. bodyless/non-semantic 503 is emitted/treated as the typed known-no-effect Product Problem;
9. delayed F13 creates new awareness for a terminal request or no-longer-eligible human;
10. duplicate same-occurrence F13 or same-Decision F14 produces duplicate semantic Notifications;
11. known-empty deciders fail to materialize Work, duplicate reevaluation creates duplicate Work, or restored eligibility leaves stale blocking Work indefinitely;
12. a missed invalidation wake-up permanently leaves a request PENDING even though durable current owner evidence proves invalidity;
13. replay correlation is evicted while its AuthorizationRequest/Decision historical authority remains within the admitted replay lifetime;
14. River job args copy review basis, credentials, arbitrary PII or provider payload by convenience;
15. any D7-R path bypasses exact Organization scope/RLS or creates cross-owner writes inside the Governance owner transaction.

Real PostgreSQL/River/runtime execution is required when implementation proof becomes authorized; document/mock-only tests cannot claim those runtime properties.

## 15. Explicit non-decisions / reopen triggers

D7-R does not freeze:

```text
queue names
worker counts
sweep cadence
SQL/table names
exact random-token encoding
exact Go package layout
metrics dashboards
SSE/WebSocket/realtime push
```

Reopen only if measured/runtime evidence proves one of these is correctness-relevant.

Actionable-list scale may later justify a bounded eligibility candidate projection/cache, but only as an optimization with authoritative revalidation; it may never become decision authority by existence.

D8-R is now **NEXT / NOT STARTED**. B10 remains suspended. Product implementation remains blocked until accepted D9.

## 16. Accepted result

```text
D7 stack reconstruction                 NO
new Product operation                   0
new ordinary Permission                 0
new NotificationKind                    0
new Work operation                      0
new infrastructure service              0
idempotency namespace                   Organization + Principal + operation + key
idempotent replay before If-Match       REQUIRED
eligibility authority                   current Q / no event-cache authority
validity external network in owner tx   FORBIDDEN
semantic 503 discriminator              exact Product Problem type
F13/F14                                 independent River reactions + revalidation/idempotency
zero-decider                            Work materialize + reconcile
invalidation                            event wakeup + durable PENDING recovery sweep
architecture proof class                STRUCTURAL AUTHORITY ONLY
runtime real proof                      REQUIRED AFTER IMPLEMENTATION GATE
```

**Operator-ratified outcome:** D7-R is accepted as the bounded Global-Maximum runtime realization for the current AuthorizationRequest authority. D8-R becomes NEXT / NOT STARTED; B10 remains suspended; PR #61 remains unmerged; Product implementation remains blocked until D9.
