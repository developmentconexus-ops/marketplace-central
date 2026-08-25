# D3 — Communication / Events

> **Status:** CLOSED / ACCEPTED AS A WHOLE — current consolidated authority  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authorities:** `D0-PRODUCT-SYSTEM-DEFINITION.md`, `D1-DOMAINS-BOUNDARIES.md`, `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Original acceptance:** 2026-08-16  
> **Later accepted amendments consolidated:** Personal Notifications + AuthorizationRequest communication/recovery

## 1. Purpose and imported invariants

D3 defines how accepted owners communicate **without moving or duplicating authority**. D3 does not choose provider DTO/auth/paging (D4), Product HTTP wire (D5), frontend topology (D6), or queue/outbox/worker/transaction/process realization (D7).

Binding parent invariants:

- communication is legal only inside accepted D1 semantic edges;
- Organization scope is explicit and never inferred from provider/source identity;
- actor attribution distinguishes human, automation/system and external-source provenance;
- domain-owned Intents/Requests/Decisions retain their own identities; no generic Action/Workflow/Correlation authority appears;
- historical refs/snapshots/projections never become current producer truth;
- provider webhook/poll payload is D4 acquisition evidence, not an MPC domain event by itself;
- unknown/absence/freshness/coverage remain fail-honest;
- **Mechanism != Authority**.

---

## 2. Communication grammar — Q / C / E / P

Marketplace Central uses the smallest communication form coherent with the meaning crossing the boundary:

- **Q — owner query:** consumer needs producer-owned **current** meaning to decide now;
- **C — owner capability request:** caller asks the callee to accept/perform work whose meaning/state belongs to the callee;
- **E — committed producer fact:** producer already committed a fact it owns and an independent accepted consumer must react;
- **P — projection/read model:** compose multiple authorities for reading/attention/analytics without creating write authority.

Selection:

```text
need current producer truth now?                   → Q
need callee to accept/perform callee-owned work?   → C
already-true producer fact + independent reaction? → E
read-only multi-authority composition?             → P
none?                                              → no communication
required semantic edge absent from D1?             → STOP / targeted D1 reopen
```

One pair of owners may legitimately use different forms at different moments. D3 does not impose one global sync/event ideology.

### 2.1 Event-worthiness

An `E` is justified only when:

1. the producer statement is already true/committed independent of the consumer; and
2. the consumer reaction is consumer-owned semantics applied to that fact.

Imperative meaning disguised as a message remains `C`. Event-per-CRUD/state-change is rejected.

### 2.2 Progression vs evidence

**Progression:** an event may wake the consumer, but the consumer Q/revalidates current producer truth when consequential currentness matters.

**Evidence:** when an individual historical occurrence is materially required for attribution/reconciliation/closure, that occurrence must remain recoverable from the smallest sufficient durable authority; latest mutable state cannot erase it.

No universal event history/event store/event sourcing is required.

---

## 3. Core accepted owner flows

### 3.1 Current-meaning feed-forward edges

Baseline `Q` remains appropriate for current meaning such as:

```text
Portfolio → marketplace-facing owners
Readiness → Offering / Availability
Offering → Availability / Market / Economics
Market → Economics
Economics → Offering
identity/access substrate → current access checks
```

A future proven autonomous reaction may add `E` inside an existing legal edge without moving authority.

### 3.2 Marketplace Sales fan-out

After Marketplace Sales commits canonical sale interpretation/attribution:

```text
Sales → Materialization       E
Sales → Fulfillment           E
Sales → Commercial Economics E
```

Consumers use `Q` for current Sale truth when material. Sales→Post-Sale eventing is only for committed post-sale-relevant Sale facts; ordinary Sale existence does not itself create a PostSaleResolution.

### 3.3 Materialization ⇄ Fulfillment

Business cycle, not shared write cycle:

```text
Fulfillment commits physical-readiness/conference fact
→ E
→ Materialization Q/revalidates as needed
→ Materialization alone owns InvoicingIntent progression

Materialization commits material business/fiscal fact
→ E
→ Fulfillment Q/revalidates as needed
→ Fulfillment alone owns provider-readiness/packing/dispatch progression
```

No shared workflow entity or cross-owner atomic mutation is required.

### 3.4 Materialization → Economics

Material attributable fiscal/business occurrences use `E`; Economics may `Q` current public meaning. Economics never becomes fiscal authority and latest state cannot erase a material occurrence required by economic history.

### 3.5 Governance ⇄ action owner

Action owner→Governance uses `C` plus `Q` as needed for one AuthorizationRequest episode and current validation. Governance owns AuthorizationRequest/Decision; action owner retains Business Intent, business disposition and execution-time validity.

When Governance commits a Decision/result after pending human review, `E` wakes the relevant owner/requester flow. An authorization event never executes the external action and never waives current policy/readiness/evidence/safety revalidation.

### 3.6 Operational Work ⇄ source owner

A source owner may emit `E` asserting only its own committed material actionable condition. Work owns whether/how that condition becomes or maps to one Work obligation.

Work may use `C/Q` to submit resolution evidence for source-owner evaluation. Closing Work alone never changes source truth. Source resolution may emit `E` so Work reconciles its lifecycle.

### 3.7 Post-Sale coordination

Post-Sale requests owner-specific consequences by `C`; consequence owners emit committed outcomes/checkpoints by `E` and expose current truth by `Q`. Post-Sale owns only resolution coordination/closure, not the underlying Sale/Fulfillment/Materialization/Economics semantics.

---

## 4. Personal Notifications communication

Personal Notifications consumes committed source-owner awareness meaning; it never becomes source truth or authorization.

Accepted producer edges are the D1 set:

```text
Portfolio
Offering
Availability
Economics
Sales
Materialization
Fulfillment
Post-Sale
Operational Work
Controlled Action Governance
  → Personal Notifications
```

### 4.1 General propagation law

For an accepted awareness occurrence:

```text
source owner commits source meaning + stable owner-local occurrence discriminator
→ recoverable E / equivalent durable producer contract
→ Personal Notifications materializes exact personal awareness
→ source continuation later Q/revalidates current source access/truth
```

Delivery may duplicate, arrive late/out-of-order or be replayed without changing source truth. Required awareness propagation must be recoverable; a crash between source commit and materialization cannot create a silent permanent gap.

Provider webhook/polling does not bypass this rule: D4 acquisition/refetch/translation occurs before an MPC owner commits awareness-worthy meaning.

### 4.2 Payload/content fence

The source-owner communication carries only the bounded immutable atoms required by D2:

- Organization scope;
- typed source/reference identity;
- stable source occurrence discriminator;
- source occurrence/commit time where material;
- exact recipient or owner-derived audience basis when applicable;
- safe source-owned `subject_display_label` material;
- F02/F14 typed result atom where applicable;
- bounded Work/Post-Sale replacement correlation where applicable.

No provider DTO, arbitrary payload/metadata/template bag, copied mutable source state, credentials or unnecessary PII is admitted.

### 4.3 Audience and route-time semantics

- `DIRECT_SOURCE`: source owner supplies exact historical human lineage.
- `OWNER_DERIVED`: source owner resolves exact current eligible human set under its own authority.
- `ORG_ROUTED`: Personal Notifications resolves the D2 route revision applicable at `source_committed_at`.

Delayed/replayed delivery must not use a later current route. Unconfigured-at-commit stays no-recipient; later configuration does not backfill. Recipient access/Membership continuity is revalidated against the D2 eligibility epoch semantics.

### 4.4 Bounded supersession under arbitrary ordering

The two accepted overlap cases must converge regardless of delivery order:

```text
source attention ↔ richer causally-related WORK_ASSIGNMENT
SALE_ATTENTION ↔ richer same-resolution POST_SALE_ATTENTION
```

Exact typed occurrence/replacement correlation is preserved in communication so Personal Notifications may deduplicate/supersede without title/time guessing or a generic causal graph.

---

## 5. AuthorizationRequest communication/recovery

### 5.1 Request intake

An action owner asks Governance to accept one semantic authorization episode by `C`, anchored to the domain-owned target/action context. Same semantic intake retry converges to the same AuthorizationRequest; a genuinely new reauthorization episode creates a distinct Request.

Ambiguous transport acceptance must be reconcilable by the owner/request semantic anchor rather than creating duplicate Governance state.

### 5.2 Current decision eligibility and material validity

Before Governance commits a decision it does not rely on a stale notification/event or target ETag as authority. It revalidates by current `Q`/owner contract as needed:

```text
Request still PENDING
+ Principal still eligible to decide
+ authorization basis still materially valid enough
```

Irrelevant target drift does not force reapproval. Material governing drift may invalidate the pending request or require a later new episode.

### 5.3 Action-required awareness

When a pending AuthorizationRequest becomes actionable for one or more exact human Principals, Governance emits/communicates committed actionable meaning sufficient for `AUTHORIZATION_ACTION_REQUIRED` awareness keyed to `AuthorizationRequestRef`.

Eligibility changes do not create Notification authority. Personal Notifications materializes awareness; Governance remains authority for whether the Request is currently decidable by a Principal.

### 5.4 Decision result continuation

When Governance commits a terminal Decision, result communication carries:

```text
AuthorizationTargetRef
+ authorize | reject
+ exact requester/initiator human lineage when that family applies
+ Governance-owned occurrence discriminator
```

Requester awareness remains target-oriented; it does not force `governance.read` and does not claim the authorized action executed/converged.

### 5.5 Zero-decider condition

If a pending Request has no currently eligible human decider, Governance/Work communication may create or reconcile explicit Work/attention. Work does not become approver or decision authority. Recovery continues to be based on current Governance truth, not on whether one notification was delivered.

---

## 6. Failure, duplicate, ordering and recovery semantics

### 6.1 Explicit Organization scope

Any durable communication/recovery state that can outlive the producing call preserves Organization scope explicitly. Bare Installation/SourceInstance/native key/Principal last-used Organization never supplies tenant scope.

### 6.2 Occurrence discrimination

Where correctness must distinguish duplicate delivery from two distinct same-valued occurrences, the producer exposes the smallest stable source/domain occurrence discriminator. It is owner-local and does not create a universal EventID.

If this requires a materially new business identity, reopen D2 instead of inventing a transport ID.

### 6.3 Semantic idempotency

Transport dedupe may optimize, but consumer-owned semantic idempotency prevents duplicate business effect. A generic processed-event table is never sufficient domain correctness by itself.

Examples:

- duplicate Sale delivery → no duplicate BusinessOrderIntent;
- duplicate actionable condition → no duplicate Work obligation while distinct conditions remain distinct;
- duplicate Notification source occurrence → no duplicate personal awareness;
- duplicate authorization intake/result → no duplicate Request/Decision.

### 6.4 Arrival order is not business truth

No global delivery order is assumed. Late progression triggers current owner reread/reconciliation; older delivery cannot regress newer committed meaning. Historical evidence remains interpretable from source/domain time/provenance, not arrival order.

Owner-local monotonic revision is introduced only when that owner actually needs it.

### 6.5 Recoverable propagation

For a committed fact whose reaction is required for Product progression/awareness, sufficient durable owner/evidence state must remain to detect and recover a missed reaction.

D7 may implement outbox, queue, worker, poller, checkpoint or sweep; none becomes business authority. Automatic recovery need not create Work unless unresolved degradation itself becomes actionable.

### 6.6 Replay/redelivery

Replay means the same semantic occurrence:

- no new business history;
- no blind external re-execution;
- current progression revalidates current truth;
- evidence re-applies the same occurrence idempotently;
- projection rebuild remains side-effect free;
- current policy is never substituted as the historical reason for an old decision.

### 6.7 Q result honesty

A Q distinguishes at least:

```text
known value
known empty/absent   // only when owner can prove it
unknown/insufficient
unavailable/error
```

Failure never becomes false/0/empty/ready/permitted. Freshness is a separate owner/consumer concern.

### 6.8 C acceptance and ambiguity

Where relevant:

```text
accepted
rejected
pending
ambiguous / unknown acceptance
```

`accepted != completed != externally applied != converged`.

A potentially accepted timeout is not automatically rejection. The caller supplies a stable Organization-scoped semantic anchor; the owner supports acceptance reconciliation/idempotent retry for the same semantic request when the class can be ambiguous.

---

## 7. Projections and multi-target communication

A `P` is rebuildable read-only composition. It may use owner queries and committed events but cannot become write authority or sole correctness truth for consequential action.

Projection update time does not prove source freshness/coverage. If required historical authority is unavailable, shrink the projection's claim rather than promote event transport into system of record.

Material multi-target communication preserves:

```text
intended target scope       // action owner
authorized scope snapshot   // Governance
attempted/outcome scope     // execution evidence
member-level partial/ambiguous outcomes when material
```

Batch correlation never manufactures cross-target atomicity or makes whole-batch blind retry safe.

---

## 8. External-effect safety fence

Every path that can reach an external effect structurally preserves as material:

- domain Intent/attempt correlation;
- actor attribution;
- duplicate/idempotency protection;
- ambiguity/no-blind-retry handling;
- attempt/outcome evidence;
- proof that required owner disposition/validity and Governance authorization are current enough.

Shared D7 machinery may verify these proofs; it never supplies the business answers. Provider protocol remains D4.

---

## 9. Contract evolution and cutover

Producer owns public semantic contract meaning. Incompatible communication-contract cutover must preserve every still-required recoverable reaction under the old contract by draining, translating, regenerating or reconciling from owner authority as the chosen D7 mechanism permits.

D3 does not require a schema registry, upcaster framework, universal event version hierarchy or simultaneous multi-version consumer platform.

---

## 10. Explicit non-decisions

D3 does not create or require:

- generic Event/Command Bus as business authority;
- generic Action/Mutation/Workflow/Saga owner;
- event-per-CRUD/state-change;
- universal CQRS/event sourcing/event store;
- universal EventID/CommandID/SagaID/Correlation identity;
- global order/sequence/vector clock;
- exactly-once delivery;
- generic reconciliation business domain;
- broker/outbox/queue/worker/lock/process topology;
- distributed transaction or microservice decomposition.

Exact D7 mechanics remain free to choose the smallest realization that satisfies this semantic contract.

---

## 11. Proof / reopen triggers

A conforming realization must survive, among others:

- duplicate Sale/actionable/Notification/authorization delivery without duplicate semantic effect;
- late progression event without state regression;
- late individual economic/fiscal occurrence without historical loss;
- producer crash after commit before propagation without silent permanent stall;
- consumer crash after own commit before acknowledgement with idempotent redelivery;
- ambiguous capability timeout without duplicate callee-owned work;
- stale approval after material governing drift cannot execute;
- owner Q unavailable cannot become plausible business value;
- same native ID in different Organizations cannot collide in recovery/dedupe;
- provider webhook duplicates/out-of-order cannot become MPC truth/order;
- old Notification route revision still governs a delayed occurrence committed under it;
- revoke→re-enable cannot silently reactivate old routed responsibility;
- requester Notification never grants Governance/source capability;
- zero-decider recovery never converts Work into authorization authority;
- projection replay cannot mutate business state or require broker history as source of truth.

Reopen only the implicated D3/parent decision if a required edge cannot be made recoverable without moving authority, a real occurrence requires new D2 identity, D7 cannot preserve Organization/temporal semantics, or a shared mechanism begins deciding business truth.

Framework preference or desire for a particular broker/event architecture is not reopen evidence.
