# D7-C — Durable Work & External Effects

> **Status:** OPERATOR-RATIFIED  
> **Parent:** `D7-RUNTIME-JOBS-TRANSACTIONS.md`  
> **Accepted prerequisites:** D7-A Runtime Envelope + D7-B PostgreSQL Isolation/Transactions — OPERATOR-RATIFIED  
> **Parent authorities:** accepted D3 communication/failure semantics, accepted D4 external-effect contracts, accepted D7-A/B runtime and transaction boundaries  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Derived:** 2026-08-21  
> **Ratified:** 2026-08-22

## 1. Purpose

D7-C selects the smallest durable-work realization that makes accepted D3 recoverable propagation and D4 external-effect ambiguity/reconciliation semantics executable without creating a generic workflow/event business authority.

D7-C owns only technical mechanics for transactionally durable post-commit work, admitted producer-fact → consumer reaction delivery, external-effect dispatch and crash recovery, retry/backoff classification, authoritative reconciliation work, and bounded schedules/polling/recovery sweeps already required by accepted D4 consumers.

D7-C does not create new Product operations, business events, intents, owners, workflow states, provider semantics, broker topology or Product implementation.

## 2. Target invariant

> **If accepted business state requires later work, the owner commit and the durable handoff are atomic; work may execute more than once without duplicating business meaning; external writes are never repeated after possible acceptance merely because a worker crashed or timed out; and queue/scheduler state never becomes business truth or required historical authority.**

## 3. Accepted inputs

### From D3

- Q/C/E/P meaning is semantic and transport-neutral;
- event delivery may duplicate, arrive late/out of order, fail or replay;
- consumer-owned semantic idempotency, not transport dedupe, prevents duplicate business effects;
- no global delivery order or exactly-once claim is required;
- consequential E propagation must be detectable/recoverable;
- progression consumers re-query/revalidate current owner truth when currentness matters;
- material evidence occurrences remain recoverable from the smallest sufficient durable authority, not from an infinite transport log;
- replay never authorizes a fresh external effect;
- required accepted consumers must not be hidden behind one best-effort fan-out;
- projections remain rebuildable/read-only and transport is never their sole history authority.

### From D4

- provider notification/polling is acquisition evidence, not MPC business truth;
- current material provider meaning comes from authoritative reread;
- consequential external-effect contracts distinguish rejection, accepted submission, synchronous effect, pending and ambiguity;
- transport `2xx` does not prove convergence;
- possible acceptance must not be blindly retried;
- every external write retains its Installation/SourceInstance-qualified target, owner intent/correlation anchor and authoritative reread/reconciliation surface;
- Sankhya remains API Gateway only; Direct Oracle is never a recovery path.

### From D7-A/B

- one Go process per replica hosts HTTP + workers over the same PostgreSQL;
- owner consequential intake commits owner state + required technical correctness records in one owner-local transaction;
- external network effects occur only after that transaction commits;
- Organization scope is explicit and structurally enforced with D7-B transaction-local scope/RLS;
- `pgx/v5` + `pgxpool` are accepted PostgreSQL primitives.

## 4. Durable-work mechanism

### 4.1 River — ACCEPTED

Select current River (`github.com/riverqueue/river`) with the pgx v5 driver as the bounded durable-work engine.

Why it is the smallest fit:

- PostgreSQL is already the canonical state store, so no Redis/broker service is required;
- River supports transaction-bound `InsertTx`, allowing owner state and required job handoff to commit/rollback atomically;
- a job inserted in the owner transaction is not work-visible until commit and disappears if the transaction rolls back;
- workers execute in-process under the D7-A topology;
- River provides retries/backoff, scheduled/periodic jobs, unique jobs, graceful stop, stuck-job rescue and transactional completion without MPC building a generic queue framework;
- River can later run insert-only or with different worker roles if D7-A's process-split reopen triggers become real.

Exact River version is an implementation-manifest concern; D7-C depends on these documented behavioral properties, not on a speculative version pin.

### 4.2 No second generic MPC outbox

For work executed by River, the transactionally inserted River job row **is the durable technical handoff**.

Do not add a second generic `outbox_events`/`messages` MPC table merely to copy owner facts into River later. That would add another relay, failure boundary and retention surface while `InsertTx` already provides atomic owner-state → durable-work handoff inside the same PostgreSQL transaction.

A separate outbox may be reconsidered only if a later accepted consumer requires a transport that cannot participate in the owner PostgreSQL transaction, or if River itself fails a required property. Such a reopen remains D7 mechanism scope, not business authority.

River job retention is not business history. Owner/domain state and accepted external evidence remain the durable authorities required by D2/D3.

### 4.3 Job arguments

Job arguments are bounded technical routing data, not provider/business DTOs or a universal business event envelope.

They carry only what a worker needs to re-enter an accepted owner boundary, proportionately:

```text
organization_id
job kind / technical reaction kind
stable semantic target or owner occurrence discriminator
owner intent / resolution / decision / work reference when applicable
source-qualified target reference when required for technical routing
```

Rules:

- no bearer/session credentials;
- no raw provider payload by convenience;
- no arbitrary PII merely to avoid an owner/source reread;
- no generic SagaID/WorkflowID/EventID is invented;
- business actor attribution references accepted canonical lineage; the worker process identity never replaces the actual human/automation/source cause.

## 5. Producer commit → consumer reaction

For an accepted D3 **E** whose consumer reaction is required:

```text
BEGIN producer owner transaction + exact Organization scope
  -> commit producer-owned fact/state
  -> InsertTx one durable reaction job for each explicit accepted consumer
COMMIT
```

Each consumer gets an independently recoverable job/reaction lane. One consumer failing cannot silently starve another consumer behind a shared in-memory fan-out.

The River job identifies the producer occurrence using the smallest accepted semantic discriminator. It does not copy mutable producer state merely to avoid a Q.

Consumer worker:

```text
claim technical job
  -> enter exact Organization scope
  -> call consumer owner application boundary
  -> consumer applies its semantic duplicate predicate
  -> Q/revalidate producer current meaning when D3 says currentness matters
  -> commit only consumer-owned state
  -> complete job
```

When completion can be transactionally coupled to the consumer's local PostgreSQL mutation, River transactional completion is preferred so consumer state and job completion do not create an avoidable acknowledgement gap. Duplicate delivery remains safe even without this optimization.

## 6. Delivery semantics

### 6.1 At-least-once / repeat-safe baseline

D7-C makes **no exactly-once execution claim**.

River may retry failed jobs and may rescue jobs believed to be stuck. A rescued job can overlap/repeat if the original execution was still alive. Therefore:

> **River retry/uniqueness is operational optimization; accepted owner semantic idempotency is correctness.**

A job must be safe if delivered again after process death, timeout, rescue or acknowledgement loss.

### 6.2 Unique jobs

River unique-job constraints may suppress redundant queued work when a stable technical key is available, but uniqueness is never the only duplicate-business-effect control.

The design must still be correct with uniqueness disabled or with two distinct River jobs representing the same producer occurrence. Consumer semantic identity/discriminator + Organization scope remains authoritative.

### 6.3 Ordering

No queue order becomes business order.

Late/out-of-order progression work revalidates current owner truth before consequential progression. Evidence work uses the occurrence's source/domain time and discriminator. River priority/scheduling may optimize latency but never defines business chronology.

## 7. External-effect execution protocol

A consequential external write uses an explicit owner effect/intake state; the River job is only the executor trigger.

### 7.1 Pre-dispatch gate

Before any network write, the worker enters one owner transaction and:

1. establishes exact Organization scope;
2. resolves the owner intent/effect anchor;
3. resolves exact prior idempotency/intake result;
4. revalidates current owner business disposition and execution-time validity;
5. revalidates required Governance authorization/currentness where applicable;
6. revalidates material provider/source prerequisites and target/blast-radius evidence required by D4;
7. establishes one durable **dispatch-attempt marker** proving that this semantic attempt is about to leave MPC;
8. commits.

If the gate fails, no provider write occurs.

The dispatch marker is execution-safety mechanism/history. It does not become a generic Mutation business owner and does not decide business permissibility itself.

### 7.2 Network write outside the database transaction

Only after the pre-dispatch transaction commits does the adapter perform the external write.

The adapter returns the smallest D4 semantic effect evidence/classification, never raw provider payload as business truth.

### 7.3 Post-dispatch persistence

After a response/failure is observed, the worker enters a new owner transaction and persists the owner-visible effect outcome proportionately: definitive rejection; accepted/pending external work; synchronous confirmed effect where the D4 contract truly proves it; ambiguous/unknown possible acceptance; or reconciliation-required state/evidence.

Convergence remains controlled by the accepted authoritative reread surface, not by queue/job success.

## 8. Crash and ambiguity rule

The persisted pre-dispatch marker creates the key fail-safe distinction:

```text
no dispatch marker
  -> worker may safely try the pre-dispatch path

dispatch marker + definitive persisted outcome
  -> resolve/replay existing outcome

dispatch marker + no definitive outcome
  -> POSSIBLE ACCEPTANCE
  -> DO NOT RE-DISPATCH
  -> schedule/perform authoritative reconciliation
```

Therefore a process crash after bytes may have left MPC but before the response/outcome was persisted cannot become a blind retry.

A River stuck-job rescue/redelivery encountering an unresolved dispatch marker must transition into reconciliation behavior rather than call the external write again.

## 9. Retry classification

Retry is classified by **effect safety**, not generic HTTP status or queue attempt count.

### Class R1 — safe technical pre-dispatch failure

Examples: database serialization/deadlock retry before dispatch marker commits, local dependency unavailable before any external request, rate limiter cannot grant dispatch, adapter cannot build/validate request before send.

- automatic retry/backoff is allowed;
- no external acceptance could have occurred.

### Class R2 — read/reconciliation/acquisition failure

Authoritative reads, polling and reconciliation reads are repeatable subject to source rate limits/freshness contracts.

- automatic retry/backoff is allowed;
- repeated read must preserve D4 partial/unavailable semantics.

### Class R3 — definitive external rejection/precondition failure

Provider/business-system contract proves the write was not accepted.

- do not blindly auto-resubmit the same effect;
- persist rejection/conflict evidence;
- owner decides whether a new/rebased attempt is admissible.

### Class R4 — possible acceptance / ambiguous write

Timeout, connection loss, process death or protocol behavior where acceptance may have survived.

- **never auto-resubmit**;
- persist/retain ambiguous attempt state;
- enqueue reconciliation/read work only;
- owner resolves accepted/rejected/pending/converged state from authoritative evidence.

### Class R5 — unknown/unclassified write failure

If the adapter cannot prove non-acceptance, fail toward R4 ambiguity rather than R1 retry.

A library's default retry behavior never overrides this classification. External-effect workers must return/mark job outcomes so River cannot automatically redeliver an ambiguous write as a fresh dispatch.

## 10. Reconciliation

Reconciliation is a technical execution pattern under owner semantics, not a generic Reconciliation business domain.

For each external-effect contract, D4 already names the authoritative reread surface. D7-C supplies durable scheduling/execution only.

Reconciliation work:

```text
exact Organization + owner effect anchor
  -> authoritative source/provider reread
  -> adapter translates bounded current evidence
  -> owner interprets convergence/ambiguity under its accepted semantics
  -> persist owner state
  -> if still unresolved, schedule another bounded read according to owner/source policy
```

A reconciliation job never issues the original write unless the owner establishes a **new** admitted semantic attempt under the normal intake/pre-dispatch path.

Persistent unresolved actionable conditions may become Operational Work under accepted source-owner + Work semantics. A queue failure alone does not automatically create Work.

## 11. Scheduling and polling

River scheduled/periodic jobs are admitted as technical wake-up mechanisms for already-accepted needs such as D4 acquisition/polling where notification/coverage alone is insufficient, reconciliation rereads, bounded recovery sweeps for missed required reactions, and source-token/config maintenance only when later D7-D/E authority admits it.

Laws:

- periodic scheduler state is not business truth or source completeness;
- an exact cron tick is not a domain occurrence identity;
- duplicate periodic wake-ups must be harmless;
- a missed periodic tick must not create an unrecoverable lifecycle stall; the next sweep/startup/recovery path can discover remaining due work from durable owner/source state;
- correctness-critical deadlines remain owner/source semantics; a generic scheduler does not own them;
- no global `backfill|incremental|sweep` business vocabulary is introduced.

Current River client leadership/maintenance scheduling is therefore an efficiency mechanism, not the sole correctness proof for a due business obligation.

## 12. Queue topology and backpressure

Exact queue names, worker counts and priorities are deployment/implementation details unless D8 proves a correctness need.

D7-C requires only that a noisy work class cannot force another correctness-critical class to violate a proved service/deadline property. River supports multiple queues/concurrency limits if D8/D7-E proves such isolation is required.

Do not pre-create per-owner queues or a generic priority taxonomy by symmetry.

## 13. Failed/discarded work

A durable job that exhausts safe technical retries remains inspectable; it does not disappear silently.

Operational handling must distinguish job transport/executor failure, owner semantic rejection/block, external-effect ambiguity, source/provider unavailability, and unresolved actionable business condition.

River discarded/failed state may support technical inspection, but canonical business meaning remains in owner state. D7-E decides the smallest operational alert/UI/log surface; D7-C does not make River UI a Product requirement.

## 14. Rejected/deferred alternatives

### Hand-rolled PostgreSQL queue / `FOR UPDATE SKIP LOCKED`

**REJECT baseline.** Reimplements claiming, retries, scheduling, rescue, inspection and shutdown behavior that River already provides over the accepted PostgreSQL/pgx stack. Reopen only if real proof shows River cannot preserve a required MPC property.

### External broker (Kafka/NATS/RabbitMQ/Redis queue)

**REJECT baseline.** No accepted throughput, fan-out, isolation or cross-database consumer requires another infrastructure service. It would also require a second atomic-handoff mechanism/outbox to bridge the owner PostgreSQL transaction.

### Generic workflow/saga engine

**REJECT.** Business progression remains distributed across accepted semantic owners using Q/C/E, not one technical workflow authority.

### Universal outbox/event log

**REJECT.** River transactional jobs handle current durable reaction/effect work; owner state/evidence owns history. No consumer requires an infinite event transport/history surface.

### Exactly-once execution framework

**REJECT.** It is unnecessary and misleading under process/network failure. Repeat-safe owner semantics are required instead.

## 15. Falsifiable proof contract

D7-C cannot close D7 without executable proof capable of falsifying at least:

1. owner transaction commits state but loses required River handoff;
2. owner transaction rolls back but leaves a runnable River job;
3. one explicit consumer failure silently prevents another accepted consumer reaction from existing;
4. duplicate/redelivered reaction creates duplicate consumer business state;
5. correctness breaks when River uniqueness suppression is disabled;
6. late/out-of-order progression work regresses current owner truth;
7. process dies after dispatch marker commit but before external outcome persistence, then restarted worker re-dispatches the write;
8. timeout after possible external acceptance causes automatic write retry;
9. definitive pre-dispatch failure is never retried even though no effect was possible;
10. reconciliation job can mutate via the original write path without a new owner attempt;
11. job Organization tampering/cross-tenant processing escapes D7-B RLS/runtime-role constraints;
12. raw provider payload/credential/PII appears in River args/loggable metadata by convenience;
13. queue/job retention becomes the sole history required for material evidence or projection rebuild;
14. missed/duplicate periodic wake-up permanently loses a required recoverable reaction;
15. exhausted safe retries make required work disappear without inspectable owner/technical state;
16. target Sankhya recovery path falls back to Direct Oracle.

Proof requiring crash/restart, transactional enqueue, repeat execution and cross-Organization worker behavior must use real PostgreSQL/River integration, not mocks alone.

## 16. Current primary evidence

Revalidated on 2026-08-21 from current River primary sources:

- River repository/current docs: <https://github.com/riverqueue/river>
- River documentation: <https://riverqueue.com/>
- River Go API: <https://pkg.go.dev/github.com/riverqueue/river>

Relevant current properties:

- transaction-bound job insertion commits/rolls back with the caller transaction and jobs do not run before commit;
- clients may work jobs in-process or be insert-only;
- retries/backoff, scheduled/periodic jobs, unique jobs, graceful stop, stuck-job rescue and transactional completion are built in;
- River explicitly documents that stuck-job rescue can cause repeated/duplicate execution if a job was still alive, reinforcing the need for semantic idempotency;
- River uses client leadership for maintenance facilities including periodic work across a shared schema, so scheduler leadership is treated only as a wake-up optimization and not as durable business authority;
- current River releases remain on the v0 line; exact pin belongs to the implementation manifest.

Current review also found open 2026 River issues around concurrent stuck-job rescue and periodic-scheduler state advancement. Those observations do not become MPC authority; they reinforce the accepted rule that River completion/scheduling state is never business truth and that durable owner/source state must make duplicate rescue and missed wake-ups recoverable.

## 17. Adjudication

**OPERATOR-RATIFIED:** River is the accepted bounded durable-work engine over `pgx/v5` PostgreSQL. `InsertTx` is the baseline owner-state → durable-work atomic handoff and no second generic MPC outbox is admitted for River-executed work. Delivery is repeatable; River uniqueness/retry is optimization only. Consequential external writes persist a pre-dispatch marker; possible acceptance/crash uncertainty routes to authoritative reconciliation instead of redispatch. Scheduled/periodic River work is only a wake-up mechanism over durable owner/source state.

No Product operation, Permission, Principal kind, semantic owner, frontend decision or provider contract changes.

Next is **D7-D — Authentication / Session / CSRF / Machine Token Realization**.

Do not begin D8, D9 or Product implementation.