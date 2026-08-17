# D3 Batch B2 — Communication Contract & Failure Semantics — Independent Review Candidate

> **STATUS: REVIEW CANDIDATE — NOT ARCHITECTURE AUTHORITY**  
> **Stage:** D3 — Communication / Events  
> **Operator posture:** B2 direction approved; independent Fable challenge required before canonical filing  
> **Parent B1:** `D3-COMMUNICATION-EVENTS.md` §3 is ACCEPTED / CANONICAL  
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

Accepted D3-B1 decides **which semantic communication forms exist and why**:

- **Q** — current owner query;
- **C** — explicit owner capability request;
- **E** — committed producer-owned fact for independent reaction;
- **P** — read-only multi-authority projection.

B1 also establishes:

- progression events may wake a consumer while current truth still comes from the owner;
- evidence consumers may require individual material occurrences rather than latest mutable state;
- consequential event propagation must be recoverable, never permanently silent;
- transport does not become authority;
- provider notifications are not automatically MPC domain events;
- no cross-owner distributed transaction or shared mutable authority is required.

B2 must now answer:

> **What minimum semantic contract keeps Q/C/E/P correct when communication duplicates, arrives late or out of order, disappears, is replayed, times out after possible acceptance, or feeds a rebuildable projection — without choosing D4/D7 transport/runtime mechanisms or introducing universal event sourcing?**

The target is not perfect delivery. The target is **business correctness despite imperfect delivery**.

---

## 2. Evidence / reachable failure classes

This candidate is justified by current evidence, not hypothetical distributed-systems fashion.

### 2.1 Duplicate and out-of-order notifications are reachable

The Evidence Register records current Bling webhook behavior where:

- the same webhook can be sent more than once;
- generation order is not guaranteed;
- retries can continue for days;
- receipt is not proof of completeness/current authoritative state.

Therefore duplicate/out-of-order handling is a D3 correctness question, not speculative hardening.

### 2.2 Submission/acceptance is not final convergence

Current provider evidence includes:

- Amazon bulk operations queued asynchronously with later item issues and possible partial work even when a feed ends `FATAL`;
- Mirakl imports returning tracking identity while price changes may remain pending before becoming effective;
- Magalu fiscal submission where HTTP success is not accepted-state proof.

Therefore `accepted`, `completed` and `converged` are materially different meanings.

### 2.3 Transport dedupe is not semantic idempotency

Legacy ADR-013 proved a concrete path where a notification lacking transport `_id` could be inserted repeatedly while domain correctness still required zero duplicate business effect. The durable lesson is:

> **transport/message deduplication may reduce work but cannot be the sole proof of semantic idempotency.**

### 2.4 Blind write replay is unsafe

Legacy ADR-029 plus accepted D0/ARCHITECTURE constraints show that a write whose acceptance may have occurred cannot be blindly retried unless replay safety is established.

B2 must preserve that same semantic distinction for capability requests and later D7 execution paths without importing HTTP-method rules as architecture.

### 2.5 Trigger arrival order is not source truth

Legacy ADR-024 proved the defect class where multiple observation triggers can see the same resource at different moments and last-scheduled commit must not defeat a newer observation merely because it arrived later.

B2 must carry the anti-regression principle without inheriting old `backfill` / `incremental` / `webhook` implementation vocabulary as target architecture.

---

## 3. Root cause and target invariant

### Root cause

B1 defines semantic communication topology but, without a B2 contract, implementations can still accidentally make correctness depend on:

- exactly-once delivery;
- transport-level dedupe;
- arrival order;
- latest mutable state where occurrence history is material;
- timeout interpreted as rejection;
- blind retry after ambiguous acceptance;
- a broker/event log as hidden history authority;
- projection update time as source freshness;
- one opaque batch outcome instead of granular member outcomes.

Those are mechanism accidents becoming business semantics.

### Target invariant

> **Communication may duplicate, arrive late/out of order, fail or be replayed without changing business truth. Current truth comes from its owner; material historical occurrence comes from the smallest sufficient durable authority; ambiguous acceptance remains explicit until reconciled; transport remains mechanism, never authority.**

---

## 4. Candidate B2.1 — Explicit Organization scope

Every communication concerning Organization-owned state or persisted external evidence is evaluated under an explicit Organization isolation scope.

Accepted candidate rules:

1. Organization scope must be explicit at the semantic boundary.
2. Physical realization may carry it as a typed contract field or trusted execution context, but it must be available to both producer and consumer correctness checks.
3. Organization must never be inferred from Marketplace Installation, Selling Entity, provider account, source key, Principal's last-used organization, request-global default or process-global state.
4. Cross-Organization replay/delivery is invalid even when external resource identifiers collide.
5. Platform-scoped definitions remain platform-scoped only where D2 already permits that meaning.

B2 does **not** choose row-level-security, schema, transaction-context or message-envelope implementation; D7 owns enforcement mechanics.

---

## 5. Candidate B2.2 — Minimum event meaning and payload

A D3 event remains an immutable notification of a producer-owned committed fact.

Candidate rules:

1. Event meaning is producer-owned and past/committed, never an imperative request.
2. Event payload contains only the minimum stable identities and immutable occurrence facts consumers materially need.
3. Event payload does not mirror a whole mutable aggregate merely to save a query.
4. If a consumer needs current mutable producer truth, it uses **Q** per accepted B1.
5. Provider webhook/callback/poll payload remains external acquisition evidence, not an MPC domain event until an owning domain commits its own meaning.
6. PII is minimized; credentials/secrets never become communication payload merely for convenience.

B2 does **not** require universal event-schema inheritance, a shared business-event base class or a schema registry.

---

## 6. Candidate B2.3 — Event/message occurrence identity

B2 does **not** introduce a universal business `EventID`, `EvidenceID` or generic entity graph.

Candidate direction:

- a concrete runtime may assign a technical unique occurrence/message identifier for delivery dedupe, tracing or replay bookkeeping;
- that technical identifier carries **no business authority**;
- semantic idempotency is anchored in identities/subjects the consuming domain actually understands: Sale, domain Business Intent, Post-Sale Resolution, Authorization Decision, Work/source condition, external source-qualified movement/result or another typed owner-defined subject;
- when an owner genuinely has a stable domain occurrence identity, consumers may use it, but B2 does not manufacture one globally.

Transport occurrence identity helps answer “have I processed this delivery?”; it does not automatically answer “has this business consequence already been represented?”

---

## 7. Candidate B2.4 — Principal / actor attribution

Communication carries Principal/actor attribution only when actor identity is materially part of the committed meaning, decision lineage, authorization or audit explanation.

Candidate rules:

1. Human, delegated automation and system action remain distinguishable per D2.
2. A provider-originated or externally observed fact is attributed to source provenance; it does not invent an MPC human/system Principal merely because the system processed it.
3. A technical worker/service that delivers a message is not automatically the business actor that caused the underlying decision.
4. Where an Authorization Decision or controlled Business Intent already carries canonical actor/authority context, events reference that owner-owned lineage rather than duplicating a second actor authority.

Exact authentication/session propagation belongs to D5/D7 where applicable.

---

## 8. Candidate B2.5 — Correlation and causation

B2 preserves business correlation without creating a universal workflow identity.

Candidate direction:

- use existing typed semantic identities when they express the relationship: Sale, domain Business Intent, Resolution, Authorization Decision, Work, owner-defined occurrence/result and external source-qualified result;
- correlation across a multi-step flow may carry more than one typed reference where materially required;
- technical trace/span/correlation identifiers may exist for observability but never become canonical business identity or reconstruct business meaning by themselves;
- no universal `CorrelationID` business aggregate, generic SagaID or WorkflowID is introduced merely because flows cross domains.

Causation is preserved to the smallest depth needed for material explanation/recovery. B2 does not require an infinite causation graph.

---

## 9. Candidate B2.6 — Provenance and time

Communication must preserve the D2 distinction between materially different time/provenance meanings.

Candidate rules:

1. Source/effective/event time, MPC observation/acquisition time, MPC record/commit time, decision time, external deadline and transport delivery time are not collapsed into one universal `occurred_at`.
2. Event transport may carry technical recorded/published/delivered timestamps, but those do not replace domain/source time semantics.
3. Re-delivery time does not rewrite when the business fact occurred.
4. Currentness/freshness conclusions remain with the consuming owner using the relevant provenance/time; message age alone is not a universal freshness rule.
5. Unknown source time remains unknown rather than fabricated from delivery time.

No universal timeline framework is introduced.

---

## 10. Candidate B2.7 — Duplicate delivery and semantic idempotency

Duplicate delivery is permitted and must be safe.

Candidate invariant:

> **Transport dedupe may reduce repeated work; consumer-owned semantic idempotency prevents duplicate business effect.**

Examples:

- two deliveries of the same Sale occurrence must not create two Business Order Intents for the same accepted semantic obligation;
- two deliveries of one source actionable condition must not create two Work obligations for the same material condition, while genuinely distinct conditions must not be collapsed by an over-broad dedupe key;
- repeated notification of one Authorization Decision must not create a second Decision or waive current revalidation;
- projection updaters must tolerate duplicate deliveries without side effects outside projection state.

The consuming owner chooses its semantic duplicate predicate. B2 forbids a universal `{event_id -> processed}` table from being treated as sufficient domain correctness.

---

## 11. Candidate B2.8 — Ordering, late delivery and anti-regression

B2 assumes **no global delivery order**.

Candidate invariant:

> **Arrival order never defines business order or authority.**

### Progression edges

When a late event merely wakes current progression, the consumer re-queries/revalidates current owner truth before a consequential transition when currentness is material.

A delayed old checkpoint therefore cannot roll current state backward merely because it arrived after a newer one.

### Evidence edges

A material occurrence remains processable even when delivered late. Its true source/domain time/provenance determines historical meaning; arrival order does not erase it.

### Bounded local revision

If a specific producer contract genuinely requires a monotonic revision/version to distinguish superseding current-state observations, that producer may define one locally.

B2 does **not** create:

- global sequence;
- total order across domains;
- universal aggregate version;
- Lamport/vector clock framework;
- ordering by queue partition as business truth.

This carries ADR-024's anti-regression principle without inheriting its old trigger taxonomy.

---

## 12. Candidate B2.9 — Recoverable propagation / missed delivery

Accepted B1 requires every consequential committed event reaction to be recoverable rather than permanently silent.

B2 makes that semantic requirement explicit:

> **After a producer commits a fact whose consumer reaction is required for an accepted Product 1.0 lifecycle, sufficient durable evidence/state must remain for the missed reaction to be detected and recovered.**

Therefore this is not a sufficient correctness design by itself:

```text
commit canonical state
  -> best-effort in-memory publish
  -> forget forever
```

Candidate rules:

1. producer commit does not become uncommitted merely because propagation fails;
2. propagation failure must become detectable/reconcilable rather than silent permanent loss;
3. recovery may use owner public state, durable pending work/intent, preserved occurrence evidence, reconciliation sweep or replay where retained — exact choice depends on the edge;
4. the event transport log is not required to be the durable recovery authority;
5. D7 chooses outbox/transaction/worker/queue/poller/checkpoint implementation.

B2 defines the semantic obligation, not the mechanism.

---

## 13. Candidate B2.10 — Progression recovery vs evidence recovery

B1 already distinguishes two classes; B2 defines their recovery contract.

### Progression recovery

When the missed message merely represents a progression trigger and current owner state is sufficient, recovery may reconcile:

```text
producer current/public state
+ consumer current state
-> is required progression missing?
```

Reconstructing every historical event is unnecessary.

### Evidence recovery

When a consuming domain's correctness requires individual material occurrences — e.g. invoice, reversal, adjustment, settlement movement — latest state is insufficient.

Those occurrences must be recoverable from the **smallest sufficient durable authority**, which may be:

- owner canonical history/state where material;
- preserved external observation/evidence;
- authoritative external source that can still be re-observed honestly.

If no accepted authority can recover a genuinely required occurrence class, surface a D2 lineage gap. Do **not** silently substitute latest state and do **not** introduce universal event sourcing as a transport workaround.

---

## 14. Candidate B2.11 — Replay / redelivery semantics

Replay is not permission to recreate business history or blindly repeat external effects.

Candidate rules:

1. redelivering/replaying one communication does not create a new producer business fact;
2. progression replay re-evaluates current owner validity/readiness before a consequential transition when currentness matters;
3. evidence replay re-applies the same material occurrence idempotently to the consuming domain's interpretation/history;
4. replay of communication associated with an external effect never implies blind re-execution of that external effect;
5. projection replay/rebuild is side-effect free outside projection state;
6. replay does not rerun current policy and claim it was the historical reason for an old action.

B2 does not require that every event be retained forever or globally replayable.

---

## 15. Candidate B2.12 — Query result semantics

Q must preserve honest knowledge state.

Candidate distinctions:

- **known value**;
- **known empty/absent** where the owner can legitimately prove absence for the asked scope;
- **unknown/insufficiently known** where the owner has that semantic state;
- **unavailable/error** where the owner could not answer.

These meanings are not interchangeable.

In particular:

> **failure to reach/query an owner must never silently become `false`, `0`, empty, absent, ready or permitted.**

Exact API/error encoding belongs to D5. B2 fixes only the semantic distinction.

---

## 16. Candidate B2.13 — Capability outcomes

A capability request asks another owner to perform/accept owner-owned work. Its outcome vocabulary must distinguish request acceptance from final completion/convergence.

Candidate semantic outcomes where applicable:

- **accepted** — callee accepted/created/continued the owner-owned work or intent;
- **rejected** — callee definitively refused under its own semantics;
- **pending** — decision/work remains unresolved, for example pending human authorization or long-running owner process;
- **ambiguous / unknown acceptance** — only where the caller cannot know whether acceptance occurred and acceptance could have survived the caller's failure/timeout.

Rules:

1. `accepted` does not mean completed, externally applied or converged;
2. timeout/disconnect is not automatically `rejected` when acceptance may already have occurred;
3. ambiguity is not universalized to every local function call — if the realization can prove atomic non-acceptance, a definitive failure remains definitive;
4. later completion/result belongs to the callee and may surface by owner query or committed fact according to B1.

Exact wire status/error codes belong to D5/D7.

---

## 17. Candidate B2.14 — Capability duplicate/retry safety

A caller must not create duplicate callee-owned intent/state merely because the capability response was lost.

Candidate invariant:

> **When acceptance may have occurred, reconcile against a stable semantic business anchor before creating a second consequence.**

Use the smallest owner-understood anchor available, for example:

- domain Business Intent identity;
- Post-Sale Resolution + consequence scope;
- existing Work/source condition subject;
- Authorization Decision / authorization case;
- Sale + materialization obligation;
- other typed owner-defined scope.

B2 does **not** introduce a universal CommandID/ActionID just to make every C request idempotent.

If a specific provider later supports an external idempotency key, D4/D7 may use it; it does not replace owner-domain semantic duplicate safety.

---

## 18. Candidate B2.15 — Controlled external-effect safety fence

Accepted B1/ADR-018 disposition allows shared execution-safety mechanics only without business authority.

B2 candidate invariant:

> **Every path capable of reaching a consequential external side effect must structurally prove the required safety context, but the shared mechanism verifies owner-issued proofs; it does not decide business policy, disposition, authorization or validity meaning.**

Shared mechanics may centralize:

- attempt/correlation bookkeeping;
- duplicate/idempotency protection;
- ambiguous-outcome handling;
- actor attribution capture;
- audit/attempt/outcome capture;
- fail-closed verification that required domain disposition, Governance authorization and execution-time owner validity evidence are present/current enough according to their owners.

The answers remain owned:

- business disposition -> action-owning domain;
- consequential authorization -> Controlled Action Governance;
- execution-time validity -> action owner under D0.7n;
- provider protocol -> D4 adapter;
- runtime transaction/retry implementation -> D7.

Mechanism != Authority.

---

## 19. Candidate B2.16 — Projection update and rebuild contract

A projection remains disposable derived read state.

Candidate rules:

1. incremental updater tolerates duplicate and out-of-order communication;
2. arrival order cannot regress the projection's claim about owner state/history;
3. projection update has no business side effects;
4. rebuild reads public owner current state plus only material historical state/evidence genuinely required by that projection;
5. retained event transport is never the sole rebuild system of record;
6. `projection.updated_at` / rebuild completion time does not prove source freshness, evidence coverage or completeness;
7. if required owner history does not exist, the projection reduces its claimed historical/content coverage honestly rather than fabricating it;
8. a projection never becomes the sole correctness source for consequential writes.

Exact projection storage/topology belongs to D6/D7.

---

## 20. Candidate B2.17 — Multi-target / partial outcome communication

Communication carrying material multi-target actions or outcomes preserves the D0/D1 scope distinctions.

Candidate rules:

- intended target scope remains owned by the action-owning domain;
- authorized target scope remains Governance decision context;
- attempted/outcome scope records actual execution evidence;
- communication cannot collapse member-level `confirmed`, `rejected`, `ambiguous`, `not-executed` distinctions into one opaque batch boolean when those distinctions are material;
- a failed/ambiguous member does not make already-confirmed members safe to replay;
- cross-provider atomicity is never invented by a batch message.

Exact bulk API/wire shape remains D5/D7.

---

## 21. Candidate B2.18 — Contract ownership and evolution

The producer owns the meaning of its public semantic communication contract.

B2 intentionally adopts a minimal evolution rule:

1. incompatible semantic changes require explicit impact analysis across material consumers;
2. consumers must not infer private producer fields/types as stable contract;
3. hard-cutover is allowed while current deployment/process constraints permit it;
4. no mandatory universal schema version, registry, compatibility layer, upcaster framework or multi-version consumer support is created now;
5. if D7 later introduces independently deployed consumers or durable messages that can outlive incompatible code changes, that new evidence may justify bounded contract versioning then.

Prepare the seam, not the framework.

---

## 22. ADR-019 / ADR-024 B2 disposition candidate

### ADR-019

D3-B1 already rehomes the semantic dependency/fan-out lesson. B2 proposes to rehome its remaining D3 correctness residue through:

- accepted consumer reaction is explicit;
- consequential propagation is recoverable;
- propagation failure cannot silently starve one accepted consumer while another remains fed;
- duplicate handling is semantic, not only transport-level;
- provider/content translation remains D4.

If accepted, ADR-019 has no remaining target architecture role and may become `historical` after canonical B2 consolidation.

### ADR-024

D3-B1 already rehomes one semantic interpretation/write owner. B2 proposes to rehome the remaining D3 residue:

- all trigger/acquisition classes converge on the same owner semantic interpretation path;
- arrival/scheduling order never defines truth;
- an older observation must not regress a newer committed interpretation merely because it arrives later;
- exact freshness/revision predicate remains owner/source-specific rather than a universal scheduler rule.

If accepted, ADR-024 has no remaining target architecture role and may become `historical` after canonical B2 consolidation.

Backfill/incremental/webhook/manual trigger vocabulary remains D4/D7 evidence only.

---

## 23. Explicit YAGNI / rejected accidental complexity

B2 deliberately does **not** create:

- universal EventID as business identity;
- universal CommandID / ActionID;
- universal CorrelationID business aggregate;
- Event Bus business domain;
- Command Bus business domain;
- Workflow/Saga engine;
- universal exactly-once guarantee;
- global ordering / total sequence;
- universal aggregate version;
- event sourcing;
- universal event history;
- mandatory CQRS;
- schema registry;
- mandatory version field/upcasters for every event;
- universal retry framework;
- broker choice;
- Kafka / RabbitMQ choice;
- outbox table/schema;
- worker/poller/scheduler topology;
- lock/lease/transaction design;
- microservice/process split.

D7 chooses the smallest runtime mechanism that can prove the accepted B2 semantics.

---

## 24. Decision Loop / Global Maximum

### Evidence

Known reachable failure classes include duplicate/out-of-order webhooks, asynchronous provider processing, partial/later outcomes, ambiguous external-write acceptance and hidden-secondary-consumer starvation in legacy paths.

### Root cause

Without an explicit communication contract, delivery/runtime accidents become business semantics.

### Target invariant

Communication failure/replay/order cannot change authority or business truth; recovery anchors in accepted owner state/evidence; ambiguity stays explicit.

### Alternatives

#### A. Best effort + latest-state reread everywhere

Advantages:

- minimal runtime complexity;
- easy local implementation.

Fails because:

- producer crash can permanently starve required reaction;
- latest state loses material historical occurrences;
- capability timeout may duplicate accepted work;
- no explicit duplicate/order semantics.

#### B. Exactly-once + global ordered event bus

Advantages:

- superficially simple consumer model.

Fails because:

- exactly-once end-to-end business effects are not established by broker delivery guarantees;
- total ordering is unnecessary for most edges;
- introduces heavy mechanism before D7 evidence;
- risks transport becoming authority.

#### C. Universal event sourcing

Advantages:

- replay/history appear uniform.

Fails because:

- D0/D2 explicitly refuse universal event sourcing;
- forces all business persistence through one accidental model;
- creates history/compatibility/versioning cost without a proven consumer.

#### D. Semantic failure contract over Q/C/E/P

- duplicate safe by owner semantics;
- no assumed global order;
- missed consequential reaction recoverable;
- current truth revalidated from owner;
- material historical occurrence recovered from smallest sufficient durable authority;
- capability ambiguity explicit only where real;
- projection rebuild independent of transport history;
- runtime mechanism deferred to D7.

### Local vs Global Maximum

**Global Maximum candidate: D.**

It preserves accepted correctness with the least new structure and remains valid in a simple modular monolith or later process separation.

### Essential complexity retained

- duplicate safety;
- ordering/late-delivery honesty;
- recoverability;
- material evidence occurrence retention/re-observation;
- ambiguity after possible acceptance;
- Organization scope;
- historical provenance/time distinctions;
- granular partial outcomes.

### Accidental complexity rejected

- exactly-once infrastructure ideology;
- global order;
- universal event history;
- mandatory broker/schema/versioning framework;
- universal IDs that duplicate business identities.

### Outcome candidate

`CURRENT STRUCTURE CONFIRMED` for D0–D2 + accepted D3-B1, with B2 adding bounded failure-contract semantics. No D1/D2/B1 reopen is expected unless adversarial review finds a real authority/lineage contradiction.

---

## 25. Proof / counterexample set

The candidate must survive at least these cases:

1. **Sale occurrence delivered twice** -> no duplicate Business Order Intent/consequence for the same accepted semantic obligation.
2. **Actionable condition delivered twice** -> no duplicate Work obligation for the same material condition; distinct material conditions still remain distinct.
3. **Old Fulfillment checkpoint delivered late** -> Materialization does not regress; current owner truth governs consequential progression.
4. **Invoice + reversal delivered out of order to Economics** -> both material occurrences remain attributable; arrival order erases neither.
5. **Producer dies after Sale commit and before delivery** -> downstream required reaction remains detectable/recoverable.
6. **Consumer dies after committing own state and before transport acknowledgement** -> redelivery does not duplicate semantic effect.
7. **Authorization Decision delivered after material governing drift** -> action owner revalidates and does not execute from stale approval.
8. **Capability response lost after possible acceptance** -> caller reconciles against stable business anchor before creating another intent/consequence.
9. **Owner query unavailable** -> result does not become `false`, `0`, known absence, ready or permitted.
10. **Projection receives A, duplicate A, C, then late B** -> remains side-effect free/rebuildable and does not let arrival order invent business truth.
11. **Communication replay references an external effect** -> no blind re-execution.
12. **Multi-target action has partial outcome** -> member-level confirmed/rejected/ambiguous/not-executed remain distinguishable where material.
13. **Provider webhook duplicated with no useful transport dedupe key** -> no duplicate domain effect solely because delivery duplicated.
14. **External resource IDs collide across Organizations** -> no cross-tenant consumption/correlation.
15. **Projection update timestamp is fresh but source coverage is partial** -> projection does not claim source completeness/currentness from its own update time.
16. **Historic material occurrence required but unavailable from any authority** -> explicit lineage gap; no fabricated latest-state answer/event-store shortcut.
17. **Technical worker delivered a human-approved event** -> worker identity does not replace the human/authorization lineage.
18. **Contract change during hard-cutover-capable deployment** -> no speculative compatibility framework required unless real durable/incompatible consumers exist.

---

## 26. Reopen / stop triggers

Reopen only the implicated prior decision when evidence changes semantics/authority:

1. A consequential E edge cannot be made recoverable without creating a new semantic dependency outside D1 -> targeted D1/B1 reopen.
2. A material evidence occurrence required by accepted correctness cannot be retained/re-observed by any sufficient durable authority -> targeted D2 lineage reopen; never universal event sourcing by stealth.
3. Consumer correctness truly requires producer-global ordering/versioning beyond a bounded owner-local contract -> revisit that edge/producer contract.
4. Capability duplicate safety cannot be expressed using existing owner/domain identities without inventing a materially new Business Intent/identity -> targeted D2 ownership/identity revisit.
5. Shared execution-safety mechanism must itself decide business disposition/authorization for correctness -> stop; authority design is wrong.
6. Projection must become write authority to satisfy a real flow -> stop/revisit ownership.
7. Cross-Organization correctness cannot be enforced while Organization remains explicit tenant root -> surface D2/D7 conflict, never infer scope.

Do **not** reopen for Kafka/RabbitMQ/outbox preference, worker topology, transaction strategy, framework fashion or future microservice decomposition.

---

## 27. Independent review challenges requested

Fable should attack at least these questions and add any stronger counterexample found from authority/evidence:

1. Does B2 accidentally require universal durable event history despite rejecting event sourcing?
2. Can every accepted consequential E edge recover after `producer commit -> process death -> no live delivery` without B2 choosing D7 mechanism prematurely?
3. Does `smallest sufficient durable authority` leave a materially required Economics/Post-Sale occurrence ownerless or unrecoverable?
4. Is the no-universal-EventID decision safe for transport dedupe/replay, or is a technical occurrence identity semantically required at D3 altitude?
5. Can duplicate Work-condition delivery preserve exactly one material Work obligation without making the source domain own Work creation?
6. Does the capability `ambiguous` outcome overgeneralize distributed uncertainty to local calls, or is the bounded predicate sufficient?
7. Can capability retry be safe without a generic CommandID when caller/callee use domain Business Intent / Resolution / subject identities?
8. Is `accepted != completed != converged` sufficient for Governance, Post-Sale consequence requests and long-running owner work?
9. Can a late evidence occurrence legitimately update Economics/Post-Sale history without regressing current producer authority?
10. Does any accepted edge actually require per-producer monotonic revision now, making B2 too YAGNI by refusing a baseline version?
11. Does Organization scope need to be physically present in every message payload, or is explicit trusted boundary context semantically sufficient without weakening isolation?
12. Could technical actor/worker attribution overwrite the real human/automation cause in any accepted flow?
13. Does projection rebuild remain possible without replaying transport history for every proposed projection class?
14. Are `known empty`, `unknown`, `unavailable/error` distinctions sufficient, or does a real owner query require another semantic state before D5?
15. Does B2 preserve D0 historical-decision correctness so replay never reruns current policy and rewrites historical explanation?
16. Does multi-target communication preserve intended/authorized/attempted scope separation under partial/ambiguous outcomes?
17. Does ADR-019 have any still-valid D3 invariant not rehomed by B1+B2?
18. Does ADR-024 have any still-valid D3 invariant not rehomed by owner convergence + anti-regression?
19. Does ADR-018 require any additional D3 semantic obligation before its remaining residue can safely stay D7-only?
20. Is any candidate rule actually D4, D5, D6 or D7 mechanism leakage rather than D3 semantics?
21. Does the full B1+B2 set create a hidden Message/Workflow God Component despite rejecting one explicitly?
22. What is the strongest reachable counterexample that would force a B3 rather than proceeding directly to Global Coherence?

---

## 28. Expected reviewer response

Append one new `## FABLE — Round N (2026-08-16)` section to `AI-DIALOG.md`.

Required response shape:

- **Subject** and HEAD independently reviewed;
- **VERDICT: APPROVE / REVISE / REJECT**;
- material findings only, each with a stable finding ID such as `B2-F1`;
- evidence/authority citation and exact violated/at-risk invariant;
- corrected invariant or specific disposition;
- whether D0/D1/D2/B1 reopen is required;
- ADR-019 / ADR-024 / ADR-018 disposition challenge;
- additional reopen triggers if material;
- explicit answer: **`READY FOR D3 GLOBAL COHERENCE: YES/NO`** only if no B3 is needed, otherwise name the required B3 surface;
- end with `HANDOFF → GPT / OPERATOR` and the expected next response.

Do not edit canonical D3 authority or this candidate during the independent review.
