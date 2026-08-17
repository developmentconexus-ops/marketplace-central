# D3 Batch B1 — Communication Topology & Edge Matrix — Independent Review Candidate

> **STATUS: REVIEW CANDIDATE — NOT ARCHITECTURE AUTHORITY**  
> **Stage:** D3 — Communication / Events  
> **Operator posture:** B1 direction approved; independent Fable challenge required before canonical filing  
> **Disposable:** delete after adjudication; durable meaning belongs only in `D3-COMMUNICATION-EVENTS.md` after explicit operator batch acceptance  
> **Date:** 2026-08-16

## Reviewer bootstrap

Read the current authority path independently before reviewing this candidate:

1. `AGENTS.md`
2. `docs/engineering/rebaseline/README.md`
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/architecture/decisions/README.md`
6. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
7. `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
8. `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
9. `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`
10. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
11. only the specific legacy ADRs/current code needed as evidence.

Do **not** treat this file, chat history, prior GPT summaries or `AI-DIALOG.md` as architecture authority.

After review, append the independent Fable result to `AI-DIALOG.md` as the next append-only Fable round. Do not modify canonical D3 authority directly.

---

## 1. Problem being solved

D1 established legal semantic dependencies between business authorities. D2 established canonical/external identity, Organization isolation, one-write-authority persistence semantics, domain-local material Intent identity, Principal attribution, provenance and historical constraints.

B1 must now answer:

> **How is each already-accepted dependency realized semantically without turning communication into a second authority, without converting business cycles into private-code cycles, and without choosing D4/D5/D6/D7 mechanisms early?**

The candidate deliberately rejects a system-wide communication ideology. The proposed Global Maximum is a **semantic hybrid**: choose the minimum communication form required by the meaning crossing the boundary.

---

## 2. Candidate B1.1 — Four semantic communication forms

B1 proposes four legitimate forms.

### Q — synchronous owner query

Use when a consumer needs producer-owned **current meaning** to complete the consumer's current decision.

Examples:

- Availability needs current Readiness conclusion;
- Offering needs current economic conclusion before forming a price intent;
- a consequential transition needs current producer state after being awakened by an older event.

`Synchronous` is a semantic requirement for request/response at decision time. B1 does **not** choose HTTP, gRPC, Go interface/package layout, process topology or network transport.

### C — explicit owner capability request

Use when one context needs to ask another owner to perform/accept work that belongs to the **callee's own authority**.

The caller supplies an explicit semantic request/context; the callee decides and mutates only callee-owned state.

A capability request may return `accepted`, `rejected`, `pending` or an owner-owned identity/reference. `Synchronous` does not imply keeping a transport request open until a human/long-running workflow finishes.

### E — committed domain event

Use when a producer-owned business fact has already been committed and another accepted consumer needs to react independently.

An event is a statement of producer-owned past/current committed meaning, not a disguised command such as `PleaseDoX`.

Provider notifications/webhooks/poll results are **not automatically MPC domain events**. D4 owns concrete provider acquisition/translation; an MPC domain event exists only after the owning MPC domain has established/committed the MPC meaning it owns.

### P — projection / read model

Use when multiple authorities need to be composed for reading, operator attention, UX or analytics without creating a new write authority.

A projection may be incrementally maintained but remains rebuildable read state. It never becomes correctness/write authority merely because it is convenient or fast.

---

## 3. Candidate B1.2 — Selection rule

For each allowed D1/D2 dependency:

```text
Need producer-owned current meaning to decide now?
  → Q

Need to ask the producer/owner to perform or accept producer-owned work?
  → C

Producer-owned fact is already committed and an independent consumer must react?
  → E

Only composing multiple authorities for reading/attention/analytics?
  → P

None applies?
  → do not create communication.

Required semantic dependency absent from D1/D2?
  → STOP / targeted D1 reopen; never hide it in mechanism.
```

One accepted semantic edge may legitimately use more than one form for different moments. B1 therefore rejects classifying an entire domain pair as globally `sync` or globally `event-driven`.

---

## 4. Candidate B1.3 — Event is a wake-up/fact, not automatic current truth

When a consumer's consequential decision requires **current producer state**, receiving event `E` may trigger/reawaken the consumer, but the event does not automatically substitute for a current owner query.

Candidate pattern:

```text
producer commits fact
    ↓ E
consumer becomes eligible to react
    ↓ Q when current producer truth is material
consumer makes its own domain decision
```

This is especially important when producer state can materially change between occurrence/delivery/processing and the consumer's action.

B2 will define the minimum ordering/duplication/replay/failure semantics. B1 decides only where current owner revalidation is semantically required rather than trusting an old message as current authority.

---

## 5. Candidate B1.4 — Event-worthiness threshold

B1 does **not** adopt `event for every state change`.

A committed domain event is baseline-worthy only when:

1. the producer has committed a meaning it actually owns;
2. a current accepted D1/D2 consumer must react independently rather than merely being the caller of the operation;
3. eventual/delayed delivery does not become business authority by itself;
4. the event eliminates a real coupling/fan-out/workflow need rather than creating speculative event surface.

Therefore B1 does **not** require baseline events merely because Portfolio configuration, Readiness, Offering state, Market Intelligence evidence or Economics conclusions changed. Current owner query remains sufficient where the consumer only needs the meaning when making a later decision.

If future evidence proves autonomous reaction is required inside an already accepted semantic edge, an event may be added without moving authority.

---

## 6. Candidate B1.5 — Feed-forward D1 edge matrix

Legend:

- `Q` = consumer queries current producer-owned meaning;
- `C` = consumer explicitly requests producer-owned capability/work;
- `E→Q` = committed producer fact awakens independent consumer; consumer re-queries producer current state where current truth is material;
- `P` = read-only composition only.

The D1 arrow remains the semantic-authority flow from producer meaning toward consumer dependency; the text below names who consumes whom to avoid arrow ambiguity.

| Accepted D1 edge | Candidate realization | B1 rationale |
|---|---|---|
| Marketplace Portfolio → marketplace-facing domains | **Q** | Consumers obtain current installation participation/configuration/posture and eligible Selling Entity participation from Portfolio when needed; no broadcast baseline is currently required. |
| Readiness → Offering | **Q** | Offering consumes Readiness-owned correspondence/readiness; it must not recompute or copy Readiness authority. |
| Readiness → Availability | **Q** | Availability consumes Readiness-owned correspondence/readiness for its own availability decision. |
| Offering → Availability | **Q** | Availability obtains the marketplace representation/target it must synchronize; Offering still does not own Sellable Availability. |
| Offering → Market Intelligence | **Q** | Market Intelligence consumes own-offer representation for comparison while retaining comparability authority. |
| Offering → Commercial Economics | **Q** | Economics consumes offer/listing representation/current commercial context needed for offer-specific economic interpretation. |
| Market Intelligence → Commercial Economics | **Q** | Economics consumes Market Intelligence-owned comparable-market meaning rather than reinterpreting provider competitor payloads. |
| Commercial Economics → Offering | **Q** | Offering consumes current economic conclusion/candidate implications and then owns any resulting Price Intent. Economics never requests/executes a price write by authority. |

No baseline event is required for these feed-forward edges unless independent reaction becomes a proven current requirement.

---

## 7. Candidate B1.6 — Marketplace Sales fan-out

For the accepted edge:

```text
Marketplace Sales
  → Materialization
  → Fulfillment
  → Commercial Economics
  → Post-Sale Resolution
```

candidate realization is **E→Q**.

Reason: once Marketplace Sales has established and committed canonical MPC sale interpretation/context/correlation and transaction-specific Selling Entity attribution, several downstream authorities may need to progress independently. The Sale's existence/meaning must not depend on every downstream consumer being synchronously available.

Candidate invariant:

> **One Marketplace Sales-owned committed interpretation feeds all downstream owners. Downstream owners may react independently but never reinterpret provider transaction semantics or transaction-specific Selling Entity attribution as their own authority.**

A provider order notification itself is not `SaleCommitted`; concrete provider acquisition/verification/translation remains D4.

---

## 8. Candidate B1.7 — Materialization ⇄ Fulfillment cycle

This is the highest-risk accepted workflow cycle.

B1 proposes realizing the normal cross-owner progression through **committed checkpoints plus owner queries**, not bilateral mutation.

### Fulfillment → Materialization

```text
Fulfillment commits a material physical-readiness/conference checkpoint
    ↓ E
Materialization becomes eligible to progress
    ↓ Q when current physical state is consequential
Materialization owns its decision to create/block/advance Invoicing Intent
```

Fulfillment never creates or mutates Invoicing Intent.

### Materialization → Fulfillment

```text
Materialization commits a material business/fiscal/native-materialization result
    ↓ E
Fulfillment becomes eligible to progress
    ↓ Q when current materialization state/result is consequential
Fulfillment owns packing/dispatch/provider-readiness progression
```

Materialization never mutates Fulfillment state.

Candidate invariant:

> **The business cycle remains bidirectional; write authority does not. Each side commits only its own meaning, observes the other's public meaning, and owns its own next transition.**

B1 does not require a distributed transaction or shared mutable workflow object across both owners.

---

## 9. Candidate B1.8 — Materialization → Commercial Economics

Candidate realization: **E→Q** for material attributable business/fiscal results that can independently change economic interpretation/reconciliation.

Materialization commits its result; Commercial Economics interprets/attributes it under Economics authority and may query Materialization current/public meaning when required.

Economics does not read private fiscal tables or convert a copied materialization row into a second fiscal authority.

---

## 10. Candidate B1.9 — Controlled Action Governance ⇄ action-owning domains

### Action owner → Governance

Candidate: **C + Q** as needed.

The action-owning domain supplies its domain-owned Business Intent, intended target scope, effective action disposition and authorization-relevant context. Governance applies only authorization-specific Grant/Delegation/Decision semantics.

A capability request may create/accept an authorization case and return `pending`; it need not block until a human decision occurs.

### Governance → action owner

Candidate: **E→Q** when a material Authorization Decision is committed after the original request, especially human/pending approval.

The event awakens the action owner. It does **not** by itself waive execution-time business validity, freshness/readiness, correspondence, policy or other domain-owned safety checks.

Candidate invariant:

> **Governance decides authorization; the action-owning domain retains business disposition, domain Intent and execution-time validity. Governance never becomes provider executor.**

---

## 11. Candidate B1.10 — Operational Work ⇄ originating domains

### Source domain → Work

Candidate: **E** when the source domain commits that a material actionable condition exists and a durable Work obligation should exist.

Work then owns responsibility/assignment/escalation/work-state lifecycle without becoming owner of the originating business truth.

D3 semantically requires recoverable propagation of a material Work-creating fact so the obligation cannot disappear silently after source commit. **D7** decides outbox/worker/transaction/queue implementation.

### Work → source domain

Candidate: **C/Q** when resolution/closure evidence must be evaluated by the source domain's semantic closure requirement.

```text
Work submits/points to resolution evidence
    ↓ C/Q
source domain decides originating condition:
  resolved / unresolved / unknown-or-pending
```

Closing Work alone never changes source truth.

If the originating condition resolves independently, the source domain may emit committed resolution fact `E` so Work can reconcile/close its own lifecycle.

---

## 12. Candidate B1.11 — Post-Sale Resolution coordination

Post-Sale remains coordinator/correlator of a material Resolution; it does not absorb Sales, Materialization, Fulfillment or Economics authority.

### Sales → Post-Sale

Candidate: **E→Q** for committed sale/post-sale-relevant facts owned by Sales.

### Post-Sale → consequence-owning domains

Candidate: **C** when a Resolution deliberately requests a consequence whose semantics belong to Materialization, Fulfillment or Economics.

The semantic request carries explicit Resolution/scope/correlation context; the callee decides/executes only callee-owned consequence semantics.

This is preferred over publishing an imperative-looking `RefundNeeded`/`CancelEverything` event and making consumers infer commands from choreography.

### Consequence owners → Post-Sale

Candidate: **E→Q** for their committed consequence outcomes/checkpoints. Post-Sale uses those facts to decide whether **its Resolution** is sufficiently evidenced/closed.

Post-Sale never rewrites the other domains' facts to force closure.

---

## 13. Candidate B1.12 — Offering ⇄ Commercial Economics cycle

The accepted semantic cycle does **not** justify a bilateral mutation/event choreography baseline.

Candidate normal path:

```text
Commercial Economics owns economic conclusion
    ↓ Q by Offering when decision is needed
Marketplace Offering Operations owns resulting Price Intent / listing action
```

Economics does not publish an imperative `PriceRecommended` event that automatically mutates price. If future evidence proves an independent attention/re-evaluation reaction is needed, an event may awaken Offering while Offering still owns intent and current validity.

---

## 14. Candidate B1.13 — D2 identity/access substrate

The D2 identity/access substrate is a cross-cutting non-domain authority, not a 13th D1 business domain.

Candidate communication for correctness-critical ordinary access is **Q**:

```text
Principal + Organization + ordinary Permission question
    ↓ Q
identity/access substrate
    ↓
current identity/membership/assignment/ordinary-access answer
```

B1 does not require an eventual `RoleChanged`/`MembershipChanged` event as the sole correctness mechanism for revocation. Events/caches may later be optimization, but current ordinary access authority cannot depend only on delayed delivery.

The substrate still cannot answer substantive marketplace action permissibility/approval/execution validity; those remain with action owner/Governance per D2.

---

## 15. Candidate B1.14 — Projection/read-model semantics

Candidate uses `P` only for read composition such as:

- portfolio/operator attention across Portfolio, Readiness, Offering, Availability, Economics, Work and lifecycle state;
- normalized `OperationalStage`;
- Work + originating condition/evidence;
- authorization + domain Business Intent read history;
- material lifecycle/history views.

Projection rules:

1. projection may consume public owner queries and committed domain events for incremental maintenance;
2. projection may combine multiple authorities;
3. projection preserves enough owner/provenance semantics to avoid pretending it created the facts;
4. projection never accepts a business command by authority;
5. projection never mutates canonical owner state;
6. projection is not the correctness source for a consequential write;
7. rebuildability must not require universal event sourcing, because D2 explicitly rejected universal event-sourced persistence.

Exact physical projection/schema/frontend topology remains D6/D7.

---

## 16. Candidate B1.15 — Public semantic boundaries and code-cycle fence

For every inter-domain dependency:

> **Consumer code may depend on the producer's public semantic contract, never producer private implementation/repository/table model.**

For an accepted bidirectional semantic relationship `A ⇄ B`, B1 does not require a shared mutable object or `A implementation ↔ B implementation` import cycle. Each direction is a distinct public relationship against the relevant owner boundary.

The target must not solve import cycles by creating a generic `shared/business`, universal `contracts`, `Action`, workflow or event domain that becomes an informal owner of everybody's semantics.

Exact Go package/interface placement is not frozen by D3; only the semantic dependency constraint is.

---

## 17. Candidate B1.16 — Cross-owner atomicity

B1 does **not** make correctness depend on one transaction atomically mutating multiple D1 business authorities.

Examples such as:

```text
Sale committed
→ Materialization responsibility
→ Fulfillment responsibility
```

or:

```text
Post-Sale Resolution
→ Materialization consequence
→ Fulfillment consequence
```

are modeled as explicit correlated/convergent cross-owner communication, not as a hidden distributed/shared transaction.

D7 may use local transactional optimizations where safe, but the semantic contract should remain valid if accepted boundaries later run in different processes.

This does **not** pre-decide retries/outbox/queues/transactions; those are D7 mechanisms.

---

## 18. Candidate B1.17 — What B1 deliberately does not create

B1 does **not** create or select:

- universal Event Bus business abstraction;
- universal Command Bus;
- generic Action/Mutation owner;
- generic Workflow/Saga engine;
- universal CQRS;
- universal event sourcing;
- event for every state mutation;
- projection for every domain;
- shared mutable business model;
- cross-domain SQL;
- distributed transaction framework;
- Kafka/RabbitMQ/broker/queue;
- outbox table/implementation;
- worker/scheduler topology;
- retry timing/backoff framework;
- lock/lease strategy;
- exactly-once delivery claim;
- microservice split.

The proposed B1 Global Maximum is **semantic hybrid communication** with only the seams justified by current D1/D2 requirements.

---

## 19. Candidate legacy-ADR treatment for B1

D2's ADR transition policy requires D3 to adjudicate its reopened legacy ADRs before those records can be retired. For B1, attack these as historical/current-state evidence, not target authority:

- **ADR-018 — generic mutation envelope/poller:** preserve any real shared external-write safety/correlation property, but do not resurrect `mutations` as a universal business authority. Domain Intent stays with action owner; Governance stays authorization owner; D4 owns provider protocol; D7 owns runtime realization.
- **ADR-019 — snapshot observer/callback:** preserve the defect-class lesson that a producer path must not silently update one consumer while another legitimate consumer becomes blind. Target answer should be explicit public semantic communication, not necessarily the legacy observer callback.
- **ADR-024 — single-writer order ingest:** preserve one semantic interpretation/write authority for Sale/provider-order meaning; do not inherit the legacy ingest service shape by default.
- **ADR-026 — scheduler phase vocabulary:** preserve only any legitimate domain-local completeness/freshness/coverage meaning; `backfill`/`incremental`/`sweep`, cursor and scheduler mechanics are not global D3 business semantics and remain D4/D7 unless independently justified.

Reviewer must identify any still-valid semantic invariant in these ADRs that this candidate would accidentally lose.

---

## 20. Proof / falsification expectations

B1 should be rejected or revised if any of these cannot be made true without adding hidden authority:

1. Economics cannot directly mutate marketplace price; Offering owns Price Intent.
2. Availability can use Readiness meaning without recomputing/copying Readiness authority.
3. Fulfillment can establish physical readiness without creating Invoicing Intent.
4. Materialization can establish fiscal/native result without mutating packing/dispatch state.
5. Governance can approve/reject without executing provider action.
6. Work can close its lifecycle without declaring the originating business condition resolved by itself.
7. Post-Sale can coordinate consequences without becoming owner of fiscal, fulfillment or economic semantics.
8. A provider notification cannot become duplicate/corrupt MPC business truth merely by being delivered twice.
9. A delayed event is not automatically treated as current producer authority when current state is material.
10. A projection can be rebuilt without mutating canonical owner state or requiring universal event history.
11. A new consumer needing meaning outside D1's accepted edge set forces explicit D1 reopen.
12. Moving accepted boundaries between in-process and separated runtime topology later does not change semantic ownership.

---

## 21. Requested independent challenge

Return **material findings only**. Attack especially:

1. **Global Maximum:** Is semantic hybrid (`Q/C/E/P`) actually the smallest sustainable topology for MPC, or is a more uniform model globally better under accepted D0–D2 constraints?
2. **Selection correctness:** Are any feed-forward `Q` edges incorrectly deprived of an event required for Product 1.0 correctness, or is event baseline there speculative?
3. **Event-as-current-truth:** Does `E→Q` correctly protect consequential decisions from stale events, or does it create redundant coupling/hidden failure modes?
4. **Sales fan-out:** Is committed Sale interpretation truly an event-worthy fact for all four downstream owners, or should any relation use a different form?
5. **Materialization ⇄ Fulfillment:** Can the proposed checkpoint/event/query realization converge without hidden command cycles, private-code cycles or a new semantic owner? Identify the strongest counterexample.
6. **Governance:** Does `C/Q` request + `E→Q` later decision preserve pending human approval, stale authorization, intended/authorized/attempted scope separation and action-owner validity correctly?
7. **Operational Work:** Does source-domain `E` → Work plus Work `C/Q` → source preserve durable responsibility without Work becoming truth authority? Is the requirement for recoverable propagation correctly D3-semantic / D7-mechanism?
8. **Post-Sale:** Are explicit capability requests from Resolution coordinator to consequence owners preferable to choreography events without moving authority? Identify any hidden imperative/event mix-up.
9. **Offering ⇄ Economics:** Is Q-only baseline correct, or is an independent committed economic event required now?
10. **Identity/access:** Is correctness-critical access appropriately query-based without requiring an eventual access-event model? Does this stay within D2's non-domain substrate fence?
11. **Projection:** Can projections be incrementally event-fed yet rebuildable without universal event sourcing or becoming correctness authority?
12. **Cross-owner atomicity:** Does refusing correctness dependence on a multi-owner atomic transaction lose any Product 1.0 invariant? If yes, name the exact flow and authority conflict.
13. **Legacy ADRs:** Re-adjudicate ADR-018/019/024/026 and name every semantic invariant B1 must preserve before those legacy records can later be retired.
14. **D1 reopen test:** Does any candidate flow actually require a semantic dependency absent from the accepted D1 edge set? If yes, identify the exact missing edge; do not solve it with mechanism.
15. **Stage leakage:** Identify anything in B1 that is really D4 provider contract, D5 HTTP/API, D6 frontend or D7 runtime/transaction/queue/outbox design and should be narrowed/deferred.
16. **B2 boundary:** Identify any delivery/ordering/duplication/replay/recovery detail that B1 is prematurely deciding versus correctly leaving for D3-B2.
17. **YAGNI/future cost:** Is any event, query, capability or projection here unsupported by a current Product 1.0 consumer? Conversely, is any low-cost seam missing whose later absence would force semantic migration?

## 22. Expected response shape

Append the next Fable round to `AI-DIALOG.md`:

- **Subject:** independent challenge of `D3-B1-REVIEW-CANDIDATE.md`;
- **Head reviewed:** exact Git SHA;
- confirmation that the authority chain was read independently;
- **VERDICT:** APPROVE / REVISE / REJECT;
- material findings with IDs only;
- corrected invariant(s) where needed;
- explicit D1 reopen requirement, if any;
- explicit legacy ADR-018/019/024/026 disposition corrections, if any;
- reopen triggers;
- final statement: `READY FOR D3-B1 OPERATOR ADJUDICATION: YES/NO`;
- `HANDOFF → GPT / OPERATOR` with the exact expected next action.

Fable findings remain evidence. Canonical B1 architecture changes only after GPT adjudication, reviewer convergence where needed, and explicit operator batch acceptance.