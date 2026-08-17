# D3 — Communication / Events

> **Status:** CLOSURE CANDIDATE — D3-B1+B2 accepted and consolidated; final Global Coherence completed; pending explicit operator ratification of D3 as a whole  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authorities:** `D0-PRODUCT-SYSTEM-DEFINITION.md`, `D1-DOMAINS-BOUNDARIES.md`, `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-16

## 1. Purpose and boundary

D3 answers:

> **How do the owners and boundaries already accepted in D1/D2 communicate without transferring, duplicating or hiding authority?**

D3 decides, proportionately:

- when an accepted semantic dependency requires a synchronous owner query/capability interaction;
- when a committed producer-owned fact justifies an event;
- when multi-authority composition belongs in a projection/read model;
- how accepted D1 workflow cycles can be realized without private-code or shared-write-authority cycles;
- the minimum event/communication semantics needed for delivery, ordering, duplication, replay and recovery contracts;
- how Organization scope, Principal/actor attribution, domain-local Business Intent identity, provenance and material historical correctness survive communication.

D3 does **not** choose:

- Mercado Livre/Sankhya/payment concrete contracts, DTOs, auth, pagination, provider webhook/polling semantics or source completeness — **D4**;
- HTTP/OpenAPI endpoint/operation/error contracts — **D5**;
- frontend screens, work inboxes, approval UX or concrete UI projection topology — **D6**;
- workers, queues, brokers, outbox implementation, transaction boundaries, locks, retries, RLS enforcement, deployment/process topology or execution-safety runtime realization — **D7**;
- product implementation.

Implementation remains blocked until D9 is accepted.

## 2. Imported parent invariants

These are already authoritative from D0–D2; D3 does not reopen them by choosing communication:

1. **Every material business meaning has one semantic authority.** Communication may transport or reference meaning; it does not create a second write authority.
2. **D3 may communicate only inside the accepted D1 semantic edge set.** If a necessary semantic dependency is missing, stop and reopen only the implicated D1 decision instead of hiding the dependency in an event, API, queue, database or projection.
3. **A consumer uses producer-owned meaning only through an explicit public semantic boundary owned by the producer.** No private implementation import, cross-context SQL/private-table access or shared mutable business entity.
4. **Organization scope is explicit.** Communication must not reconstruct tenant scope from provider account, Installation, Selling Entity, source key or process-global defaults.
5. **Actor attribution remains semantically honest.** Human, delegated automation and system action remain distinguishable where material.
6. **Material domain-owned Business Intents retain domain-local identity.** D3 must not introduce a generic Action/Mutation/Command/BusinessIntent owner merely to route communication.
7. **Historical snapshots/references/projections never become current producer authority.** A read model may compose authorities; it may not become write authority.
8. **External provider/business-system facts remain external/source-qualified.** Provider notifications or payloads are not automatically MPC domain events; D4 owns concrete acquisition/translation contracts.
9. **Unknown/absence, provenance and material time meanings remain fail-honest.** Communication must not turn stale, partial, duplicated or delayed information into plausible current truth.
10. **Mechanism != Authority.** Shared transport/runtime mechanics may later centralize accidental complexity without acquiring business meaning.

---

## 3. D3-B1 — Communication Topology & Edge Matrix — ACCEPTED

**Outcome:** `CURRENT STRUCTURE CONFIRMED` with bounded review corrections. No D0/D1/D2 reopen is required.

The operator explicitly accepted the converged B1 batch after independent Fable challenge and GPT adjudication. Reviewer findings are evidence; the authoritative result is the operator-approved structure recorded here.

### 3.1 Governing communication rule

Marketplace Central does **not** adopt one global communication ideology.

> **Each dependency uses the smallest communication form coherent with the meaning crossing the boundary. Communication never changes who owns that meaning.**

Four semantic forms are accepted:

- **Q — synchronous owner query:** consumer needs producer-owned current meaning to complete the consumer's current decision.
- **C — explicit owner capability request:** caller asks the callee to perform/accept work whose meaning and state belong to the callee's own authority.
- **E — committed domain event:** producer has already committed a producer-owned fact and an accepted consumer must react independently.
- **P — projection/read model:** multiple authorities are composed for reading/attention/UX/analytics without creating new write authority.

`Synchronous` describes semantic request/response at decision time; it does not choose HTTP, gRPC, Go package layout or process topology.

### 3.2 Selection rule

For each accepted D1/D2 dependency:

```text
Need producer-owned current meaning to decide now?
  -> Q

Need the owner to perform/accept owner-owned work?
  -> C

Producer-owned fact is already committed and an independent consumer must react?
  -> E

Only composing authorities for reading/attention/analytics?
  -> P

None applies?
  -> no communication.

Required semantic dependency absent from D1/D2?
  -> STOP / targeted D1 reopen.
```

One semantic edge may use more than one form at different moments. A domain pair is therefore not globally classified as `sync` or `event-driven`.

### 3.3 Event predicate and event-worthiness

A committed fact is a legitimate **E** only when both are true:

1. the producer's statement is true and complete independent of any consumer reaction; and
2. each consumer's reaction is determined by consumer-owned semantics applied to that fact.

A communication whose semantic meaning is effectively:

> "callee, produce outcome X that you own, selected by the caller"

is **C**, regardless of whether the physical transport later looks like a message.

B1 does **not** adopt event-per-state-change or event-per-CRUD. A baseline event is justified only when:

- the producer has committed meaning it owns;
- an accepted consumer has a real independent reaction;
- delayed delivery does not become authority;
- the event removes a real coupling/fan-out/workflow problem rather than creating speculative event surface.

Therefore baseline events are not required merely because Portfolio configuration, Readiness, Offering state, Market Intelligence evidence or Economics conclusions changed. If future evidence proves an autonomous reaction is necessary inside an already accepted semantic edge, `E` may be added without moving authority.

### 3.4 Progression edges vs evidence edges

`E -> Q` must not collapse two different correctness needs.

#### Progression semantics

When the consumer's next consequential decision depends on **current producer truth**, an event may wake the consumer, but the consumer queries/revalidates the producer's current public meaning before deciding when currentness is material.

```text
producer commits fact
    ↓ E
consumer becomes eligible to react
    ↓ Q when current producer truth matters
consumer makes its own decision
```

An older event is never automatically current authority.

#### Evidence-accumulation semantics

Some consumers require individual material occurrences, not merely the producer's latest mutable state. Examples include attributable fiscal/economic movements, reversals and consequence evidence used by Commercial Economics or Post-Sale Resolution.

For those edges:

> **The evidence-consuming domain determines which occurrences are material to its correctness claim. Each material occurrence must remain recoverable from the smallest sufficient durable authority. Latest mutable state must not erase an occurrence whose existence/sequence is material to attribution, reconciliation, closure or historical explanation.**

The sufficient durable authority may be canonical MPC state/history or preserved/re-observable authoritative external observation/evidence as allowed by D2. This does **not** require a universal producer history API, universal event history, event store or event sourcing.

If no accepted durable authority can recover a genuinely material occurrence class, surface the gap; do not substitute latest state and do not create a universal event store by mechanism.

### 3.5 Consequential event propagation is recoverable

A committed fact whose consumer reaction is required for an accepted Product 1.0 lifecycle progression has **recoverable propagation semantics**:

> **loss is detectable and recoverable, never a silent permanent stall.**

This applies across consequential event edges, including sale fan-out, Materialization/Fulfillment checkpoints, later Governance decisions and material actionable conditions.

B1 defines the semantic obligation. **B2** defines per-edge duplicate/order/recovery/replay contracts. **D7** chooses outbox/queue/worker/transaction/poller or other concrete mechanism.

### 3.6 Feed-forward D1 edge matrix

For these accepted D1 feed-forward dependencies, baseline communication is **Q**:

| Accepted D1 edge | Accepted B1 realization | Rationale |
|---|---|---|
| Marketplace Portfolio -> marketplace-facing domains | **Q** | Consumers ask Portfolio for current installation participation/configuration/posture and eligible Selling Entity participation when needed. |
| Readiness -> Offering | **Q** | Offering consumes Readiness-owned correspondence/readiness and does not recompute it. |
| Readiness -> Availability | **Q** | Availability consumes current Readiness meaning for its own decision. |
| Offering -> Availability | **Q** | Availability obtains current marketplace representation/target while retaining Sellable Availability authority. |
| Offering -> Market Intelligence | **Q** | Market Intelligence consumes the organization's own offer representation while retaining comparability authority. |
| Offering -> Commercial Economics | **Q** | Economics consumes offer/listing/current commercial context needed for offer-specific economics. |
| Market Intelligence -> Commercial Economics | **Q** | Economics consumes Market Intelligence-owned comparable-market meaning rather than reinterpreting competitor payloads. |
| Commercial Economics -> Offering | **Q** | Offering consumes economic conclusions/implications and alone owns resulting Price Intent/listing action. |

No speculative baseline event is required on these edges. A later proven autonomous reaction may add `E` without moving authority.

### 3.7 Marketplace Sales fan-out

Once Marketplace Sales establishes and commits canonical sale interpretation/context/correlation and transaction-specific Selling Entity attribution, independent downstream authorities may react without making Sale existence depend on their synchronous availability.

Accepted baseline:

```text
Marketplace Sales
  -> Materialization      E
  -> Fulfillment          E
  -> Commercial Economics E
```

Downstream consumers use current Sales meaning by **Q** when currentness is material and must never independently reinterpret provider transaction semantics or Selling Entity attribution as their own authority.

**Sales -> Post-Sale** is narrower: `E` is baseline-worthy only for committed **post-sale-relevant** Sales facts. An ordinary committed Sale does not itself create a Post-Sale Resolution; D2 permits `0..N` resolutions per Sale. Post-Sale may query Sales by Q when it needs additional context.

Provider order notifications are not `SaleCommitted`; D4 owns provider acquisition/refetch/translation before Sales commits MPC meaning.

### 3.8 Materialization <-> Fulfillment

This accepted business workflow cycle is realized as two separate semantic flows, never bilateral mutation or shared mutable workflow authority.

#### Fulfillment -> Materialization

```text
Fulfillment commits physical-readiness/conference checkpoint
    ↓ E
Materialization becomes eligible to progress
    ↓ Q when current physical state is consequential
Materialization alone creates/blocks/advances Invoicing Intent
```

Fulfillment never creates or mutates Invoicing Intent.

#### Materialization -> Fulfillment

```text
Materialization commits material business/fiscal result
    ↓ E
Fulfillment becomes eligible to progress
    ↓ Q when current materialization meaning is consequential
Fulfillment alone owns provider-readiness/packing/dispatch progression
```

Materialization never mutates Fulfillment state.

> **The business cycle remains bidirectional; write authority does not. Each side commits only its own meaning, observes the other's public meaning, and owns its own next transition.**

No distributed transaction or shared cross-owner workflow entity is required.

### 3.9 Materialization -> Commercial Economics

Material attributable business/fiscal occurrences that can independently change economic attribution/reconciliation use **E**.

Commercial Economics processes the material occurrences it needs for its own evidence chain and may use **Q** for current/public Materialization meaning where currentness matters. Latest fiscal state cannot erase a material occurrence required for economic attribution/history.

Economics never reads Materialization private tables or becomes fiscal authority.

### 3.10 Controlled Action Governance <-> action-owning domains

#### Action owner -> Governance

Use **C + Q** as needed.

The action-owning domain supplies:

- domain-owned Business Intent;
- intended target scope;
- effective action disposition;
- authorization-relevant context.

Governance applies only its authorization-specific Grant/Delegation/Authorization Decision semantics. A request may return `pending`; semantic synchrony does not require an interactive transport request to remain open until a human decides.

#### Governance -> action owner

When a material Authorization Decision is committed after the original request, especially after pending human review, use **E** to wake the owner and **Q/revalidation** where current authorization/context matters.

An approval event never executes a provider action and never waives execution-time domain validity, freshness, readiness, policy, correspondence or mandatory safety invariants.

> **Governance decides authorization; the action owner retains Business Intent, business disposition and execution-time validity.**

### 3.11 Operational Work <-> originating domains

#### Source domain -> Work

The source domain may emit **E** asserting only its own committed material actionable condition.

It does **not** assert that a particular Work object must exist. Operational Work owns whether/how that condition is represented as Work, including obligation scoping, deduplication against an existing obligation and Work-local lifecycle.

However, D0's no-ownerless-work invariant remains binding:

> **A source-committed material actionable condition ends represented in Work state or explicitly reconciled as already covered/superseded. Silent disappearance is a propagation failure, not a legitimate Work decision.**

Duplicate delivery is therefore a B2/D7 idempotency concern, not a reason to move Work authority back to the source domain.

#### Work -> source domain

Use **C/Q** when Work submits or points to resolution evidence that the source domain must evaluate under the source's own closure semantics.

```text
Work submits/points to resolution evidence
    ↓ C/Q
source domain decides source condition:
  resolved / unresolved / unknown-or-pending
```

Closing Work alone never changes source truth.

If the source condition resolves independently, the source may emit a committed resolution fact **E** so Work reconciles/closes its own lifecycle.

### 3.12 Post-Sale Resolution coordination

Post-Sale remains coordinator/correlator of a material Resolution; it does not absorb Sales, Materialization, Fulfillment or Economics authority.

- **Sales -> Post-Sale:** `E` only for committed post-sale-relevant Sales facts, with `Q` for additional/current context where needed.
- **Post-Sale -> Materialization/Fulfillment/Economics:** **C** when Post-Sale requests a consequence whose semantics belong to that owner.
- **Consequence owner -> Post-Sale:** **E** for committed consequence outcomes/checkpoints; Post-Sale may use `Q` where current owner meaning is material.

This deliberately rejects imperative-looking choreography such as `RefundNeeded` or `CancelEverything` as disguised commands.

Each consequence owner accepts/rejects/pends and owns its domain-local consequence intent/state. Post-Sale decides only whether **its Resolution** has sufficient evidence to close.

### 3.13 Commercial Economics <-> Marketplace Offering Operations

The accepted semantic cycle does not require bilateral event/mutation choreography.

Normal path:

```text
Commercial Economics owns economic conclusion
    ↑ Q by Offering when needed
Marketplace Offering Operations owns resulting Price Intent / listing action
```

Economics never directly requests/executes marketplace price writes by authority. No baseline `PriceRecommended` event is required merely to create event-driven shape.

If future evidence proves a real independent attention/re-evaluation reaction, `E` may awaken Offering while Offering remains price-intent authority.

### 3.14 D2 identity/access substrate

Correctness-critical ordinary identity/access checks use **Q** against the D2 identity/access substrate for current membership/RoleAssignment/Permission meaning.

Revocation correctness cannot depend solely on eventual `RoleChanged`/`MembershipChanged` delivery. Events/caches may later optimize access checks, but delayed propagation cannot become the sole authority for current ordinary access.

The substrate still cannot answer substantive marketplace action permissibility, consequential authorization or execution validity; those remain with action owners/Governance.

### 3.15 Projection/read-model semantics

A **P** combines authorities for reading/attention/UX/analytics without creating a new business authority.

Expected uses include, where later D6/D7 justify them:

- portfolio attention;
- normalized OperationalStage;
- Work + originating condition;
- Authorization Decision + Business Intent;
- material lifecycle/history views.

Rules:

1. A projection may consume owner public queries and committed events for incremental maintenance.
2. A projection never commands/mutates canonical state by authority.
3. A consequential write cannot use a projection as sole correctness authority when the owning domain must be consulted.
4. Rebuild consumes public owner current state plus only the material historical state/evidence the particular projection genuinely requires.
5. Event transport is an incremental-maintenance optimization, never the sole rebuild authority or system of record.
6. If required history is not durably available from any accepted authority, the projection must shrink its claimed content honestly rather than promote transport-log retention into business/historical authority.
7. D3 does not require universal event sourcing, infinite event retention or event-log replay as the rebuild model.
8. Exact projection schema/topology remains D6/D7.

### 3.16 Provider / D4 fence

Provider webhook/callback/poll result is **not** automatically a D3 domain event.

```text
provider notification / polling evidence
    ↓ D4 acquisition/refetch/translation
owning domain establishes its MPC meaning
    ↓ commit
D3 domain event, only if event-worthy
```

A duplicate/out-of-order provider notification therefore does not automatically become duplicate/out-of-order MPC business truth. D4 and B2/D7 later close concrete acquisition, ordering and recovery mechanics.

### 3.17 Cross-owner atomicity

B1 requires no atomic mutation spanning multiple business authorities.

Cross-owner workflows are correlated/convergent. Partial outcomes remain explicit. D7 may exploit a local database transaction where safe and useful, but semantic correctness must not depend on every current owner remaining colocated in one process/database.

No cross-provider or cross-owner atomicity is invented.

### 3.18 Public semantic boundary and code-cycle rule

Inter-domain communication depends only on the producer's public semantic contract.

A consumer must not import:

- producer repository/store;
- producer private application types;
- producer private tables;
- provider DTOs leaked through another business context.

A bidirectional semantic edge is implemented as two owner-specific semantic flows, not a shared mutable business object.

The producer owns the meaning of its public contract. D3 does not freeze exact Go interface/package placement and does not create `shared/domain`, universal `shared/contracts` or `shared/business-events` as informal authority containers.

### 3.19 External-effect safety mechanism fence

B1 preserves the reusable safety lesson behind legacy ADR-018 without preserving its generic Mutation business owner/table/poller shape.

Every path that can reach an external side effect must cross structurally enforced execution-safety mechanics/proofs appropriate to the action, including as material:

- intent/attempt correlation;
- actor attribution capture;
- idempotency/duplicate protection;
- ambiguity handling;
- audit/attempt/outcome capture;
- fail-closed proof that required owner-issued disposition/validity and Governance authorization are present/current enough.

The shared mechanism verifies **proofs**; it does not own the **answers**:

- action disposition/business policy remains with the action-owning domain;
- consequential authorization remains with Controlled Action Governance;
- execution-time validity remains with the action owner under D0.7n;
- provider protocol remains D4.

Exact runtime mechanism remains D7.

### 3.20 Legacy ADR disposition at B1

B1 adjudicates the D3 semantic portions of ADR-018/019/024/026 as follows.

#### ADR-018 — generic mutation envelope

The generic `mutation_protocols`/`mutation_items` business shape, `/mutations` ownership and in-process poller are **not target business architecture**. Domain Business Intents remain with action owners; reusable external-effect safety is a mechanism under §3.19.

Pending approved/unexecuted domain intents are already canonical durable state under D2. Poller, claim strategy and `FOR UPDATE SKIP LOCKED` remain D7 evidence.

**D3 portion adjudicated; D7 residue remains open.**

#### ADR-019 — hidden secondary consumer starvation

The durable semantic lesson is generalized:

- accepted consumers of a committed fact must be explicit enough that a producer-path rewrite cannot silently starve one consumer while another remains healthy;
- consequential propagation failure must be visible/recoverable rather than partial-silent;
- honest content translation/parity remains required where D4 maps external evidence.

Legacy one-row-per-item/PK/sentinel schema mechanics do not constrain the clean target database.

B2 still must define concrete missed/duplicate-delivery/recovery semantics before the D3 portion is fully retireable.

#### ADR-024 — single writer for order ingest

The target preserves:

1. one Marketplace Sales semantic interpretation/write authority for provider sale meaning; and
2. trigger convergence/anti-regression: multiple acquisition triggers must converge on the owner's one interpretation path, and a late older observation cannot regress a newer committed interpretation merely because of scheduling order.

The first principle is B1 authority topology; the second becomes a B2 ordering/duplication contract. Legacy import/backfill/webhook worker names and current code shape are evidence only.

#### ADR-026 — scheduler phase vocabulary

No global `backfill | incremental | sweep` D3 vocabulary is carried forward.

Its useful semantic kernel is already D0 authority: full/terminal and incremental observations make different coverage claims; conflating them can corrupt completeness/freshness conclusions. Cursor/scheduler phase mechanics remain D4/D7 evidence.

**D3 portion adjudicated; D7 residue remains open.**

### 3.21 YAGNI / explicit non-decisions

B1 does **not** create or choose:

- generic Event Bus business abstraction;
- generic Command Bus;
- generic Action/Mutation business owner;
- generic Workflow/Saga engine;
- universal CQRS;
- event sourcing;
- universal producer/domain fact history;
- event-per-CRUD;
- projection-per-domain;
- shared mutable business model;
- cross-context SQL;
- distributed transaction;
- broker technology such as Kafka/RabbitMQ;
- concrete outbox/queue/worker/poller topology;
- retry/lock/lease framework;
- exactly-once delivery claim;
- microservice decomposition.

B1 prepares the semantic seams D3 actually needs without designing D7 runtime in advance.

### 3.22 Proof / strongest counterexamples

B1 must remain true under these checks:

- Economics cannot directly mutate listing/price.
- Availability cannot copy/recalculate Readiness authority.
- Fulfillment conference cannot create Invoicing Intent.
- Materialization result cannot mutate packing/dispatch state.
- Governance approval cannot execute provider action.
- Work closure cannot change originating source truth.
- Source actionable condition cannot dictate Work-owned obligation representation and cannot disappear silently.
- Post-Sale consequence remains owned by the consequence domain.
- Provider webhook duplicate does not automatically duplicate domain truth.
- Delayed progression event cannot masquerade as current producer truth.
- Evidence edge cannot lose a material historical occurrence merely because latest state changed.
- Projection rebuild cannot require the event transport as sole history/system of record.
- Projection cannot mutate canonical state.
- A new semantic consumer outside D1 edges triggers targeted D1 reopen rather than hidden communication.
- Future process separation may change transport but must not change semantic ownership/contracts.

### 3.23 B1 reopen / stop triggers

Revisit B1 only for material evidence such as:

1. Product 1.0 requires a semantic dependency absent from D1 -> targeted D1 reopen.
2. Materialization <-> Fulfillment cannot close without new owner/meaning -> implicated authority review.
3. Post-Sale must own consequence semantics currently owned elsewhere -> D1 reopen.
4. Operational Work must decide originating business truth itself -> D1 review.
5. Identity/access substrate must decide business action permissibility -> D1/D2 boundary conflict.
6. Correctness truly requires atomic cross-owner mutation -> reopen the implicated semantic decision; do not smuggle in a distributed transaction.
7. A projection must become write authority -> stop; re-evaluate ownership.
8. A required event cannot be stated as an already-true producer-owned fact -> it is likely C or an authority problem.
9. A feed-forward Q edge gains a proven autonomous-reaction requirement -> add E inside the accepted edge; no D1 reopen unless meaning changes.
10. B2 cannot give a consequential E edge recoverable propagation without moving semantic authority -> return to B1.
11. A genuinely material evidence occurrence cannot be recovered from any sufficient durable authority -> surface the D2 lineage gap; never silently substitute latest state or universal event sourcing.
12. D7 attempts to make shared execution-safety mechanics decide business disposition/authorization -> stop; Mechanism != Authority.

Framework preference, event-driven fashion or future microservice topology are not reopen evidence.

---

## 4. D3-B2 — Communication Contract & Failure Semantics — ACCEPTED

**Outcome:** `CURRENT STRUCTURE CONFIRMED` with bounded review corrections. No D0/D1/D2/B1 reopen and no B3 are required.

The operator explicitly accepted the converged B2 batch after independent Fable challenge, GPT adjudication and reviewer convergence with no remaining dispute.

### 4.1 Governing failure contract

> **Communication may duplicate, arrive late/out of order, fail or be replayed without changing business truth. Current truth comes from its owner; material historical occurrence comes from the smallest sufficient durable authority; ambiguous acceptance remains explicit until reconciled; transport remains mechanism, never authority.**

D3 does not require perfect delivery. It requires business correctness despite imperfect delivery.

### 4.2 Explicit Organization scope

Every communication concerning Organization-owned state or persisted external evidence is evaluated inside an explicit Organization isolation scope.

- Q/C that live entirely inside one trusted execution context may obtain the explicit scope from that context.
- Any durable communication/recovery state that can outlive its producing execution context must preserve Organization scope explicitly in its durable representation/container so later processing does not reconstruct it.
- Installation, SourceInstance, Selling Entity, provider account, external resource ID, Principal last-used Organization, request-global default or process-global state never substitute for or determine Organization.
- Installation/SourceInstance/native identifiers may participate as namespace qualifiers **inside** the explicit Organization scope.
- Duplicate predicates, reconciliation anchors, replay and recovery are always evaluated inside the explicit Organization scope. A bare external identifier cannot collapse identical IDs from different Organizations.

Exact schema/RLS/transaction/message-envelope enforcement remains D7.

### 4.3 Event payload and occurrence discrimination

An event remains an immutable notification of a producer-owned committed fact.

- payload contains only stable identities and immutable occurrence facts consumers materially need;
- it does not mirror a whole mutable aggregate merely to avoid a Q;
- provider webhook/callback/poll payload is still acquisition evidence, not an MPC event until the owner commits MPC meaning;
- PII is minimized and credentials/secrets are never propagated for convenience.

D3 introduces no universal business `EventID`, `EvidenceID` or occurrence aggregate.

However:

> **When a consumer's correctness requires distinguishing “the same material occurrence delivered again” from “two distinct material occurrences”, the producer/source public contract exposes a stable occurrence discriminator sufficient for that semantic distinction.**

The discriminator may be:

- source-qualified external movement/result identity;
- an existing canonical occurrence such as Authorization Decision;
- a bounded domain-local occurrence key/identity where that meaning genuinely exists;
- another immutable owner/source-defined key sufficient for the contract.

The evidence-consuming domain decides when that distinction is materially required. A discriminator answers **same vs different**; it does not establish before/after ordering and does not justify a universal occurrence identity class.

If correctness requires inventing a materially new business identity not supported by D2, reopen only the implicated D2 identity/lineage decision.

### 4.4 Actor, correlation, provenance and time

Communication preserves actor/provenance honestly:

- human, automation and system Principals remain distinguishable where material;
- provider-originated evidence uses source provenance rather than inventing an MPC Principal;
- a technical worker delivering a communication is not automatically the business actor that caused the underlying action/decision;
- where canonical Intent/Authorization/Work lineage already contains actor/authority context, communication references that lineage rather than creating a duplicate actor authority.

Business correlation uses typed semantic identities that already express the relationship — Sale/source-qualified order, domain Business Intent, Resolution, Authorization Decision, Work, owner-defined occurrence/result and similar accepted identities. Technical trace/span IDs may exist for observability but never become canonical business meaning. No universal SagaID/WorkflowID/Correlation aggregate is introduced.

Material time/provenance meanings remain distinct:

- source/effective/event time;
- observation/acquisition time;
- MPC record/commit time;
- decision time;
- external deadline/window;
- transport publish/delivery/redelivery time.

Delivery/redelivery time never rewrites when a business fact occurred. Unknown source time remains unknown. Freshness remains consumer/use-sensitive; message age alone is not a universal freshness rule.

### 4.5 Duplicate delivery and semantic idempotency

Duplicate delivery is permitted and must be safe.

> **Transport dedupe may reduce repeated work; consumer-owned semantic idempotency prevents duplicate business effect.**

Examples:

- repeated delivery of the same Sale obligation does not create a second Business Order Intent;
- repeated delivery of the same actionable condition does not create a second Work obligation while genuinely distinct conditions do not collapse;
- repeated Authorization Decision notification does not create another Decision or bypass current revalidation;
- projection update redelivery has no side effect outside projection state.

The consuming owner defines its semantic duplicate predicate using the explicit Organization scope plus the smallest stable semantic anchor/discriminator required by the case. A universal `{event_id -> processed}` table may be a runtime optimization but is never sufficient proof of domain idempotency.

### 4.6 Ordering, late delivery and anti-regression

D3 assumes **no global delivery order**.

> **Arrival order never defines business order or authority.**

Progression edges:

- late events wake/reconcile current progression;
- current owner truth is re-queried/revalidated when consequential currentness matters;
- an older delivery cannot roll current meaning backward merely because it arrived last.

Evidence edges:

- a true material occurrence remains processable when delivered late;
- source/domain time and provenance determine historical meaning;
- occurrence discrimination does not imply ordering.

If a specific producer genuinely requires a monotonic revision to distinguish superseding current-state observations, it may define one locally. D3 does not define global sequence, total order, universal aggregate version, vector/Lamport clock or queue-partition order as business truth.

This fully rehomes ADR-024's D3 anti-regression meaning.

### 4.7 Recoverable propagation and missed-reaction ownership

After a producer commits a fact whose consumer reaction is required for an accepted Product 1.0 lifecycle:

> **sufficient durable authority/evidence must remain for a missing required reaction to be detected and recovered.**

Producer commit does not become uncommitted because propagation fails. `commit -> best-effort in-memory publish -> forget forever` is not sufficient correctness by itself.

The owner of the progression/convergence that should have happened owns the semantic conclusion that its required reaction is missing. A generic reconciliation runtime does not acquire that authority.

Recovery may use, as appropriate:

- producer public/current state;
- consumer canonical/pending state;
- domain Business Intent;
- preserved material occurrence evidence;
- honest re-observation of authoritative external evidence;
- another accepted durable owner state capable of proving the gap.

D7 may realize this with outbox, worker, queue, reconciliation sweep, poller, checkpoint or another smaller mechanism. D3 does not require durable message records or make the event transport log the recovery authority.

A successful automatic recovery discharges the miss without automatically creating Work. Only an unresolved condition that is materially actionable becomes Work/attention under the owning domain + Operational Work semantics. Persistent/recurring recovery degradation may itself become a material condition when the responsible domain judges it so.

### 4.8 Progression recovery vs evidence recovery

For a progression trigger where current owner state is sufficient, recovery may compare public producer state with consumer state and determine whether required progression is missing. Reconstructing every historical event is unnecessary.

For an evidence edge where correctness requires individual occurrences, latest mutable state is insufficient. Each material occurrence must remain recoverable from the smallest sufficient durable authority, for example owner canonical history/state, preserved external observation/evidence, or an authoritative source that can still be re-observed honestly.

If no accepted authority can recover a genuinely required occurrence class, surface a targeted D2 lineage gap. Never substitute latest state and never introduce universal event sourcing as a transport workaround.

### 4.9 Replay / redelivery

Replay is not permission to create new business history or repeat an external effect.

- replay/redelivery of one communication represents the same producer occurrence;
- progression replay re-evaluates current owner validity/readiness where material;
- evidence replay applies the same occurrence idempotently to the consumer's interpretation/history;
- replay associated with an external effect never means blind external re-execution;
- projection replay/rebuild has no business side effects;
- replay does not rerun current policy and claim that today's policy was the historical reason for an old decision/action.

D3 does not require infinite event retention or global replayability.

### 4.10 Query result semantics

Q preserves honest knowledge state. The minimum semantic distinctions are:

- **known value**;
- **known empty/absent**, only when the owner can legitimately prove absence for the asked scope;
- **unknown / insufficiently known**;
- **unavailable / error** because the owner could not answer.

Failure to reach/query an owner never silently becomes `false`, `0`, empty, absent, ready or permitted.

Freshness is orthogonal to those four states. When freshness-for-use is material, the Q result supplies or references enough owner-controlled provenance/observation time for the consuming domain to judge freshness. A `known value` may still be insufficiently fresh for a particular use.

Exact wire/status/error encoding remains D5/D7.

### 4.11 Capability outcomes, ambiguity and retry safety

C distinguishes request acceptance from completion/convergence.

Where applicable, semantic outcomes are:

- **accepted** — callee accepted/created/continued owner-owned work;
- **rejected** — callee definitively refused under its own semantics;
- **pending** — decision/work remains unresolved;
- **ambiguous / unknown acceptance** — only when the caller cannot know whether acceptance occurred and acceptance could have survived the caller's timeout/failure.

`accepted != completed != externally applied != converged`.

Timeout/disconnect is not automatically `rejected` where acceptance may already have occurred. Ambiguity is not universalized: when a realization proves atomic non-acceptance, failure remains definitive.

For capability classes subject to ambiguous acceptance:

> **the caller supplies a stable Organization-scoped semantic anchor; the callee's public semantic contract supports acceptance reconciliation by that anchor and callee-owned semantic idempotency so retry of the same semantic request converges on already-accepted work rather than creating a duplicate.**

Examples may use Resolution + consequence scope, domain Business Intent, source condition/subject or another accepted owner-specific identity. D3 does not introduce a generic CommandID/Request aggregate solely for retry safety.

If the callee cannot expose acceptance reconciliation without leaking private implementation, revisit that edge's public contract; never bypass the semantic boundary.

### 4.12 Projection failure/rebuild contract

A projection remains read-only derived state.

- incremental update tolerates duplicate/out-of-order communication;
- arrival order does not become business order;
- rebuild uses public owner current state plus only the historical state/evidence the projection actually requires;
- event transport is never the sole rebuild/system-of-record authority;
- rebuild/replay has no business side effect;
- `projection.updated_at` does not prove source freshness or population completeness;
- if required history is not durably available from an accepted authority, the projection shrinks its claimed content honestly rather than inventing history or promoting transport retention into authority.

Exact projection storage/update topology remains D6/D7.

### 4.13 Multi-target / partial outcomes

Where communication carries a material multi-target operation/result, it preserves the D0 distinction between:

- intended target scope owned by the action domain;
- authorized scope snapshot owned by Governance;
- attempted/outcome scope from execution evidence;
- member-level `confirmed`, `rejected`, `ambiguous`, `not-executed` or equivalent distinctions when material.

A batch-level correlation/result cannot erase member-level correctness or make whole-batch blind replay safe.

### 4.14 Contract ownership, cutover and evolution

The producer owns the meaning of its public semantic communication contract. Consumers may not fork producer meaning by convenience.

D3 intentionally rejects a baseline schema registry, universal event-version hierarchy, upcaster framework or multi-version consumer machinery.

However:

> **an incompatible communication-contract cutover must preserve every required reaction that remains pending/recoverable under the old contract.**

Depending on D7 realization, cutover may drain retained pending records, translate them, or safely regenerate/reconcile required reactions from owner authority. Silently orphaning a required reaction is the same propagation failure §4.7 forbids.

This rule does not assume durable message records. True simultaneous multi-version consumer support is deferred until evidence proves it is required.

### 4.15 External-effect safety remains mechanism, not authority

Shared execution-safety runtime may structurally verify proofs and capture attempts/results, including intent/attempt correlation, duplicate/replay safety, ambiguity handling, actor/audit evidence and presence/currentness of owner-issued disposition/authorization/validity proofs.

It does not evaluate or grant business disposition, policy, readiness or authorization. Those answers remain with action owners/Governance. Provider protocol remains D4. Exact runtime realization remains D7.

### 4.16 Legacy ADR disposition after B2

#### ADR-019

Its remaining D3 meaning is fully rehomed by:

- explicit accepted-consumer semantics from B1;
- recoverable propagation/missed-reaction ownership from §4.7;
- duplicate discrimination/semantic idempotency from §§4.3/4.5;
- honest external translation parity remaining a D0/D4 concern.

Legacy listing observer/table/PK implementation does not constrain target architecture.

**D3 adjudication complete: ADR-019 becomes historical.**

#### ADR-024

Its target meaning is fully rehomed by:

- one Marketplace Sales interpretation/write authority under D1/B1;
- owner convergence and public boundary rules;
- arrival-order independence and anti-regression in §4.6.

Legacy trigger taxonomy/import/backfill/webhook implementation is D4/D7 evidence only.

**D3 adjudication complete: ADR-024 becomes historical.**

#### ADR-018

No additional D3 semantic residue remains. D1/B1/B2 already rehome domain intent ownership, external-effect safety proofs, ambiguity/replay safety and restart/recovery semantics.

Poller/table/claim/`FOR UPDATE SKIP LOCKED` and concrete execution topology remain D7 evidence.

**ADR-018 remains reopened — D7 only.**

#### ADR-026

Its D3 semantic portion remains adjudicated; no global phase vocabulary is carried forward. Scheduler/cursor mechanics remain D7 evidence.

**ADR-026 remains reopened — D7 only.**

### 4.17 YAGNI / explicit non-decisions

B2 does **not** create or require:

- universal business EventID/EvidenceID/Occurrence aggregate;
- generic CommandID/Request aggregate;
- global ordering/sequence;
- universal aggregate version;
- exactly-once delivery;
- universal `{event_id -> processed}` domain correctness table;
- universal event history/event store/event sourcing;
- infinite event retention;
- universal Saga/Workflow/Correlation business identity;
- schema registry/upcaster/versioning framework;
- multi-version consumer support;
- generic reconciliation business domain/component;
- broker, outbox table, queue, worker, poller or scheduler topology;
- cross-owner transaction or process topology.

B2 prepares only semantic seams justified by reachable failure classes.

### 4.18 Proof / strongest counterexamples

B2 must remain true under these cases:

1. same Sale occurrence delivered twice -> no duplicate Business Order Intent;
2. same material source condition delivered twice -> no duplicate Work obligation, while distinct conditions remain distinct;
3. old Fulfillment progression event arrives late -> current owner truth prevents regression;
4. invoice + reversal/adjustment arrive late/out of order -> each material occurrence remains attributable without treating arrival order as truth;
5. process dies after producer commit and before propagation -> required downstream reaction remains detectable/recoverable;
6. consumer commits own state and dies before acknowledgement -> redelivery converges idempotently;
7. Authorization Decision arrives after material governing drift -> action owner revalidates and stale approval does not execute;
8. capability timeout after possible acceptance -> caller reconciles by semantic anchor; callee does not create duplicate work;
9. Q owner unavailable -> does not become false/zero/empty/ready/permitted;
10. known Q value is materially stale -> consumer can judge freshness from provenance/time rather than trusting `known` alone;
11. projection receives duplicate/out-of-order updates -> remains rebuildable and side-effect free;
12. replay associated with external effect -> never blind re-execution;
13. multi-target partial outcome -> member-level distinctions remain intact;
14. same external/native identifier exists in two Organizations -> explicit Organization scope prevents cross-tenant dedupe/correlation collapse;
15. provider webhook duplicates/out-of-order -> transport dedupe/order never becomes domain correctness;
16. evidence-edge duplicate cannot be distinguished from distinct same-valued occurrences -> contract must provide stable bounded occurrence discriminator or surface a targeted lineage/identity gap;
17. technical worker delivers an event -> worker identity cannot overwrite the actual human/automation/source cause;
18. incompatible communication contract deploys while a required reaction remains recoverable -> cutover drains/translates/regenerates/reconciles rather than silently losing the reaction;
19. projection history is not reconstructable from owner authority -> projection reduces its claim rather than requiring a hidden event store;
20. future process separation changes failure mode -> semantic ambiguity/recovery rules survive without moving ownership.

### 4.19 B2 reopen / stop triggers

Revisit B2 or its parent decision only for material evidence such as:

1. a consequential edge cannot be made recoverable without moving semantic authority -> B1 review;
2. a genuinely material evidence occurrence cannot be recovered from any accepted durable authority -> targeted D2 lineage review;
3. evidence correctness requires a materially new business identity rather than a bounded occurrence discriminator -> targeted D2 identity review;
4. a callee subject to ambiguous acceptance cannot expose anchor-based reconciliation without violating its public boundary -> revisit that edge contract;
5. a real edge proves a baseline monotonic revision/order contract is required -> add only owner-local ordering semantics justified by that evidence;
6. D7 cannot preserve explicit Organization isolation structurally/fail-closed -> surface conflict against D2/`ARCHITECTURE.md`, never infer scope from another identity;
7. incompatible contract cutover silently drops a still-required recoverable reaction -> stop-the-line propagation failure;
8. shared safety/reconciliation mechanism begins deciding business truth/disposition/authorization -> stop; Mechanism != Authority;
9. a projection must become write authority or event transport must become sole history -> stop and re-evaluate ownership/lineage.

Framework preference, desire for Kafka/event sourcing, process separation or generic distributed-systems patterns are not reopen evidence.

---

## 5. Final D3 Global Coherence + YAGNI / Overengineering / Future-Cost review — COMPLETED

**Outcome: CURRENT STRUCTURE CONFIRMED with no material correction. No B3, D0/D1/D2 reopen or B1/B2 reopen is required.**

The review evaluates accepted B1+B2 as one communication system against D0–D2, `ARCHITECTURE.md` and the DevelopmentConexus Engineering Method.

### 5.1 Duplicate / missing authority

**PASS.** Q/C/E/P transport meaning but do not create a new business owner. Events remain producer facts; capability mutations remain with the callee; projections remain read-only; recovery/missed-reaction conclusions remain with the domain whose progression/convergence is at stake. No generic event/reconciliation/workflow authority appears.

### 5.2 Semantic-edge completeness

**PASS.** B1 realizes every accepted D1 edge plus the D2 identity/access dependency without adding a dependency outside D1. B2 strengthens failure semantics without introducing a new semantic edge. New future semantic dependencies still trigger targeted D1 reopen.

### 5.3 Business cycles / deadlock / authority cycles

**PASS.** Materialization <-> Fulfillment is a business cycle realized by two owner-specific committed-fact/query flows, not mutual writes. Governance/action-owner and Work/source cycles similarly preserve separate authorities. No shared mutable workflow object or circular write authority is introduced.

### 5.4 Current truth vs historical occurrence

**PASS.** Progression uses current owner revalidation when currentness matters; evidence consumers preserve/recover material occurrences when latest state is insufficient. The distinction avoids both stale-event authority and universal event-history/event-sourcing requirements.

### 5.5 Failure/recovery coherence

**PASS.** Duplicate delivery is semantically idempotent; arrival order is not truth; missed required reactions remain recoverable; replay does not recreate history or external effects; ambiguous capability acceptance is reconciled by owner/domain anchors. None of these properties require exactly-once delivery, global ordering or one transport technology.

### 5.6 Tenant / identity coherence

**PASS.** Organization remains explicit and is never inferred from Installation/SourceInstance/provider key. Durable communication/recovery state preserves scope. Occurrence discrimination uses existing source/domain semantics only when materially needed and does not create a universal EventID identity graph.

### 5.7 Trust / actor / authorization layering

**PASS.** Technical delivery workers do not become business actors. Governance still owns authorization only; action owners retain disposition and execution-time validity. Approval events cannot execute provider actions or waive revalidation. Shared execution-safety mechanisms verify proofs but never own policy/authorization answers.

### 5.8 External-authority / D4 fence

**PASS.** Provider webhooks/callbacks/poll results remain acquisition evidence, not MPC events by default. D3 defines only the semantics after an owning domain commits meaning. Concrete provider capabilities, authoritative rereads, source completeness and protocol remain D4.

### 5.9 Projection coherence

**PASS.** Projections never become write authority, never make update time equivalent to source freshness/completeness, and rebuild from owner state/evidence. Missing historical authority shrinks projection claims rather than turning event transport into system of record.

### 5.10 Multi-target / partial outcome coherence

**PASS.** Communication preserves intended, authorized and attempted/outcome scope separation and member-level partial/ambiguous states. Batch-level communication cannot manufacture cross-target atomicity or make whole-batch retry safe.

### 5.11 YAGNI / overengineering

**PASS.** D3 explicitly refuses unsupported generic capability:

- no event-per-state/CRUD;
- no generic Event/Command Bus as business authority;
- no Workflow/Saga engine;
- no universal CQRS/event sourcing/event store;
- no universal EventID/CommandID/SagaID business identity;
- no global sequence/version/vector clock;
- no exactly-once promise;
- no schema registry/upcaster framework;
- no generic reconciliation domain;
- no broker/outbox/queue/worker topology;
- no distributed transaction;
- no microservice split.

Every retained abstraction has a current semantic consumer or protects an accepted safety/isolation/history invariant.

### 5.12 Future-cost / seam review

**PASS.** D3 prepares only justified seams:

- public semantic contracts permit later process separation without changing authority;
- Q/C/E/P allow the simplest current realization while preserving async fan-out where required;
- explicit Organization scope survives future tenancy/process changes;
- domain/source occurrence discrimination is added only where evidence correctness needs it;
- anchor-based capability reconciliation avoids a generic Command identity while remaining safe under later remote failure modes;
- recoverable propagation allows D7 to choose the smallest runtime mechanism rather than forcing a broker/outbox now;
- contract cutover protects pending/recoverable reactions without imposing permanent multi-version support.

No irreversible structural dead end was found.

### 5.13 Later-stage leakage

**PASS.** D3 does not select provider DTOs/contracts/auth (D4), HTTP/OpenAPI/error encoding (D5), concrete projection/UI topology (D6), or worker/queue/outbox/transaction/lock/retry/RLS/process/deployment topology (D7). B2 contract-cutover and recoverability requirements constrain correctness only, not mechanism.

### 5.14 Legacy ADR coherence

**PASS.** ADR-019 and ADR-024 have no remaining target role after B1+B2 rehome their durable meaning and may be historical. ADR-018 and ADR-026 retain only D7 mechanism/runtime residue. ADR-013/029 remain carried external-evidence/write-safety constraints through D0/`ARCHITECTURE.md`; D3 does not inherit their old implementation shapes.

### 5.15 Strongest counterexamples checked

- duplicate Sale delivery cannot duplicate Business Order Intent;
- duplicate Work-condition delivery cannot create duplicate material obligation;
- same provider ID in two Organizations cannot collapse tenant scope;
- late old progression event cannot regress current owner meaning;
- same-valued distinct evidence occurrences can remain distinguishable where material;
- invoice/reversal sequence cannot be erased by latest state;
- producer crash before notification cannot leave a permanent silent stall;
- consumer crash after own commit cannot make redelivery duplicate semantic effect;
- ambiguous capability timeout cannot create a second callee-owned intent by blind retry;
- approval arriving after governing drift cannot bypass execution-time validity;
- owner-query failure cannot become a plausible business value;
- projection cannot require broker history to become truthful;
- incompatible contract cutover cannot silently orphan a required recoverable reaction;
- runtime safety/reconciliation mechanics cannot become business authority.

**Conclusion:** no material contradiction, missing authority, hidden God component, speculative framework or later-stage mechanism leak remains in D3.

---

## 6. Review and authority protocol

D3 used the same accelerated protocol proven in D2:

1. GPT prepared coherent B1/B2 candidate batches from repository authority/evidence.
2. The operator approved each candidate direction for independent challenge.
3. Disposable `D3-B<n>-REVIEW-CANDIDATE.md` files were explicitly non-authority.
4. The operator invoked Fable separately; Fable reconstructed the authority path independently and appended material findings to `AI-DIALOG.md`.
5. Reviewer findings remained evidence; GPT independently adjudicated material findings against repository authority/evidence.
6. Material disagreements received another reviewer round until convergence.
7. The operator explicitly ratified each converged batch before canonical consolidation.
8. Review candidates are disposable after consolidation.
9. D3 received the final Global Coherence + YAGNI / Overengineering / Future-Cost review in §5.
10. D3 closes only after explicit operator ratification as a whole.

`AI-DIALOG.md`, review candidates, chat summaries and reviewer statements are not architecture authority.

---

## 7. Current D3 state / exact next action

D3 is a **CLOSURE CANDIDATE**.

- **D3-B1 — Communication Topology & Edge Matrix: ACCEPTED / CANONICAL.**
- **D3-B2 — Communication Contract & Failure Semantics: ACCEPTED / CANONICAL.**
- **Final Global Coherence + YAGNI / Overengineering / Future-Cost review: COMPLETED / PASS.**
- **B3: NOT REQUIRED.**

Exact next action: **explicit operator ratification of D3 as a whole**.

If ratified:

1. mark D3 `CLOSED / ACCEPTED`;
2. update the rebaseline router so **D4 — External Integrations** becomes the exact next stage;
3. do not begin product implementation; implementation remains blocked until D9.

If a material issue is found, reopen only the implicated D3 decision rather than re-running the whole stage.
