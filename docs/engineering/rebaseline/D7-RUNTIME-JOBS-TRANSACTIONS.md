# D7 — Runtime / Jobs / Transactions

> **Status:** OPEN / ACTIVE — D7-A + D7-B OPERATOR-RATIFIED / D7-C CANDIDATE PENDING RATIFICATION
> **Program:** Architecture Rebaseline / Technical System Design
> **Opened:** 2026-08-21
> **Parent authorities:** accepted D0–D6 semantics, `ARCHITECTURE.md`, canonical Product OAD, and bounded owner authority routed by `docs/index.md`
> **Method:** DevelopmentConexus Engineering Method v1.0.0

## 1. Purpose and boundary

D7 defines the smallest target runtime realization capable of executing the already-accepted Marketplace Central Product and frontend contracts without creating new business authority.

D7 owns technical realization for:

- server serving/process topology;
- PostgreSQL persistence, structural Organization isolation and transaction boundaries;
- durable jobs, scheduling, outbox/effect handoff and retry/reconciliation mechanics;
- session/CSRF/OIDC and machine-token runtime mechanics already required by D5-R1;
- secrets, observability and deployment topology required to operate the admitted transports;
- cursor/runtime mechanics only where needed to realize accepted Product behavior.

D7 does **not** reopen Product operations, ordinary Permissions, Principal kinds, semantic ownership, frontend interaction meaning or provider/business-system semantics by implementation convenience. D8 golden-flow choreography, D9 adversarial review and Product implementation remain blocked.

## 2. Imported invariants

1. **Go is the canonical backend execution language.** Concrete source realization remains blocked until D9.
2. **PostgreSQL is the canonical MPC-owned state store.** External systems remain sources/dependencies, never alternate writable application stores.
3. **Organization isolation must be structural.** Correctness may not depend only on developers remembering tenant predicates.
4. **Product OpenAPI remains the single machine-readable Product wire authority.** Runtime conforms to it; runtime does not create a second contract.
5. **Q/C/E/P meaning remains owner-defined.** Runtime transport cannot move authority or invent hidden cross-owner calls.
6. **Consequential writes preserve idempotency, ambiguity, auditability and reconciliation.** Potentially accepted effects are never blindly replayed.
7. **Partial/unknown/unavailable remain honest.** Runtime/persistence cannot turn missing evidence into plausible known state.
8. **Provider protocol stays in adapters.** Mercado Livre/Sankhya DTOs and native choreography do not become MPC business ontology.
9. **Sankhya uses the sanctioned API Gateway.** Direct Oracle/database access is not an admitted fallback.
10. **Human browser credentials stay server-side.** H uses same-origin application session + CSRF; browser JavaScript does not own OIDC access/refresh tokens.
11. **Machine clients remain A/S bearer clients.** Client Credentials does not become a Principal kind or authorization shortcut.
12. **No Product code before D9.** D7 produces accepted realization authority and proof contracts, not the application implementation.

## 3. D7 target invariant

> **Every admitted Product or Technical-Ingress interaction can be realized through explicit serving, persistence, transaction and durable-work boundaries that preserve Organization isolation, authority ownership, authentication carrier separation, idempotency/concurrency semantics and honest external-effect outcomes, while platform mechanisms remain business-policy-free.**

## 4. Decision surface

D7 must resolve only the mechanics needed by accepted D0–D6 behavior.

### D7-A — Runtime envelope and process topology

Decide the minimum number and responsibility of runtime processes/executables and the serving boundary between browser/static delivery, Product API, Technical Ingress and background work.

### D7-B — Persistence, isolation and transactions

Owning accepted authority: [D7-B PostgreSQL Isolation & Transactions](D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md).

Decide PostgreSQL connection/transaction ownership, structural Organization isolation/RLS, owner command transactions, ETag/revision realization, idempotency lifecycle and only justified read projections.

### D7-C — Durable work and external effects

Owning candidate: [D7-C Durable Work & External Effects](D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md).

Decide post-commit durable work, schedules/polling, admitted E reactions, external-effect dispatch/reconciliation, retry classification and atomic handoff between owner state and required work.

### D7-D — Authentication/session request-trust realization

Realize ApplicationSession, OIDC exchange, CSRF, Keycloak/provider topology and A/S token validation under the already-decided D5-R1 carrier split.

### D7-E — Operability boundary

Decide only the structured observability, secrets/configuration, migrations, deployment/readiness/backup and real-dependency proof mechanics needed by the resulting runtime.

## 5. Candidate set is evidence, not selection

Accepted architecture identified candidates worth testing, including modular-monolith class, pgx/pgxpool, PostgreSQL RLS, River, OpenTelemetry/OTLP/slog, sqlc, tern, Keycloak and real-dependency test tooling. D7-A selected the modular-monolith process envelope; D7-B selected pgx/pgxpool + PostgreSQL structural isolation. River is under D7-C adjudication. Remaining candidates stay unselected until their owning slice proves a current consumer/property.

## 6. Proof obligations

Before D7 can close, the target design must make at least these properties falsifiable:

- cross-Organization access/bypass fails structurally, including privileged/service paths;
- owner state plus required durable handoff cannot commit one without the other where atomicity is required;
- same idempotency key + materially different request fails closed;
- ambiguous potentially accepted provider/business-system write is not blindly retried;
- stale revision/precondition remains distinct from business/provider rejection;
- human bearer use cannot silently replace the accepted session carrier;
- unsafe H request without valid CSRF trust fails before owner effect;
- A/S bearer cannot resolve to H and token audience/issuer/client binding fails closed;
- provider callback/acquisition evidence cannot impersonate authoritative current source truth;
- worker replay/duplication/out-of-order delivery cannot mutate committed business truth incorrectly;
- failed job/effect remains inspectable/reconcilable rather than disappearing;
- target Sankhya wiring cannot fall back to Direct Oracle;
- migrations/boot fail on an incompatible persistence baseline rather than silently continuing;
- observability does not log secrets/provider PII by convenience.

D8 will exercise composed golden flows. D7 must provide the runtime contracts and proof seams D8 needs; it does not pre-run D8 here.

## 7. YAGNI exclusions

Absent a concrete falsifier/consumer, D7 does not introduce microservices/service mesh, Kubernetes/operator work, Kafka/NATS, a generic workflow engine, CQRS/event sourcing, multi-region active-active, Redis/cache infrastructure, generic ORM/repository layers, a second database, separate BFF business API, provider plugin framework, realtime/WebSocket infrastructure or AI/MCP runtime authority.

## 8. Decision order

D7 proceeds dependency-last:

1. **D7-A Runtime envelope:** process/serving responsibilities + transaction ownership boundary — **OPERATOR-RATIFIED**.
2. **D7-B Persistence:** isolation + transaction invariants before database/schema mechanics — **OPERATOR-RATIFIED**.
3. **D7-C Durable work:** atomic handoff/retry/reconciliation before worker selection — **CANDIDATE / OPERATOR RATIFICATION PENDING**.
4. **D7-D Authentication:** session/OIDC/CSRF/token mechanics under the accepted carrier split.
5. **D7-E Operability:** only required deployment/observability/secret/migration mechanics.
6. One combined D7 coherence/proof review before closeout.

## 9. Current primary evidence for D7-A

Revalidated on 2026-08-21:

- Go `net/http.Server.Shutdown` provides graceful HTTP shutdown without interrupting active connections: <https://pkg.go.dev/net/http#Server.Shutdown>.
- Go `ServeMux` supports method, host and path routing patterns, so route expressibility alone does not require a framework: <https://pkg.go.dev/net/http#ServeMux>.
- Go `signal.NotifyContext` provides signal-driven context cancellation for process lifecycle: <https://pkg.go.dev/os/signal#NotifyContext>.
- Current River documentation shows a PostgreSQL-backed client can start worker loops in background goroutines inside the application process, can be insert-only when queues are omitted, supports graceful worker stop, and supports transaction-bound `InsertTx`: <https://pkg.go.dev/github.com/riverqueue/river> and <https://riverqueue.com/>.

## 10. D7-A accepted — Runtime Envelope & Transaction Ownership Boundary

> **OPERATOR-RATIFIED:** one Go application process per replica, one public application origin, Product API + Technical Ingress + human session mediation + durable worker runner in the same process; owner business transactions stay local to one semantic owner and all consequential external effects happen outside the database transaction after an atomic durable handoff.

### 10.1 Alternatives adjudicated

| Alternative | Disposition | Reason |
| --- | --- | --- |
| one Go process per replica: HTTP + workers | **ACCEPTED** | smallest topology; current Go lifecycle and PostgreSQL-backed durable-work capabilities satisfy the required envelope without another service boundary |
| same codebase but separate API and worker processes | **DEFER / REOPEN TRIGGER** | useful only with proved resource/failure/scaling isolation; current job candidates allow this split later without changing business ownership |
| microservices / per-owner services | **REJECT** | no current independent deployment/security/scaling consumer; adds network/distributed-transaction/operability failure modes without Product value |

### 10.2 Process and serving shape

Baseline:

```text
https://conexus.fun
        |
        v
one Go application process / replica
  ├─ same-origin frontend/static delivery
  ├─ Product API handler boundary
  ├─ Technical Non-Product Ingress handler boundary
  ├─ H session / CSRF / OIDC mediation boundary
  ├─ in-process durable worker runner
  ├─ scheduler coordinator seam
  └─ PostgreSQL connection pool
```

Binding laws:

- **one process does not mean one business module**: accepted owner/context boundaries remain explicit inside the modular monolith;
- Product API and Technical Ingress use separate route/middleware composition boundaries even though they share listener/process infrastructure;
- no internal owner-to-owner HTTP merely for symmetry: Q/C/P inside the monolith use explicit owner-owned application ports/interfaces; E remains durable only where D3 admits an independent reaction;
- browser/static + API remain one public origin. A separate CDN/edge/runtime is not baseline and requires an operational consumer;
- no separate screen-shaped BFF appears: server-side OIDC/session mediation protects the same Product API;
- correctness must not depend on exactly one replica;
- no HTTP router/framework is selected by D7-A. A dependency must prove a concrete missing property before admission.

### 10.3 Owner transaction ownership

For a consequential owner command:

```text
request / admitted internal C
  -> authenticate + scope/Permission/precondition gates
  -> begin one owner command transaction
       -> claim/validate idempotency identity when required
       -> mutate only that owner's canonical business state
       -> write required audit/intake/outcome metadata
       -> atomically persist required durable handoff/job
  -> COMMIT
  -> return accepted/pending/domain outcome
  -> background worker performs external effect after commit
  -> authoritative reread/reconciliation determines convergence
```

Laws:

- one transaction may include platform records required for correctness but **must not mutate another semantic owner's private business state**;
- cross-owner business changes do not use a distributed transaction;
- external network writes are never performed while holding the owner PostgreSQL transaction open;
- potentially accepted external effects are never blindly replayed after transport ambiguity;
- external reads/preflight evidence, when required, occur outside the write transaction;
- ordinary Q reads do not acquire a write transaction;
- Product HTTP handlers, workers and Technical Ingress adapters enter the same owner application boundary rather than duplicating business logic.

### 10.4 Lifecycle and health boundary

Startup must fail closed on missing critical configuration, incompatible database/migration baseline or inability to establish required local dependencies.

Graceful shutdown order:

```text
1. mark application not-ready;
2. stop accepting new HTTP work;
3. stop fetching new durable jobs/scheduled work;
4. allow active HTTP requests and jobs a bounded drain window;
5. cancel remaining work after the bound while leaving durable work recoverable;
6. flush telemetry and close database/resources.
```

Provider/Sankhya unavailability is **not** application unreadiness by itself; accepted Product knowledge/outcome semantics represent those dependencies as unavailable/partial/pending where appropriate.

### 10.5 Reopen triggers

Split API and worker process roles only if a real proof establishes at least one of background resource starvation/shutdown coupling, materially different scaling envelopes, a required security/exposure boundary, D8 drain/recovery failure or an unsafe deployment lifecycle. Preference for process separation or microservices is not evidence.

## 11. D7-B accepted — PostgreSQL Isolation & Transaction Realization

[D7-B](D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md) is **OPERATOR-RATIFIED**.

Accepted baseline:

```text
organization_id on organization-owned state/evidence
+ composite Organization foreign keys
+ transaction-local scope
+ PostgreSQL ENABLE + FORCE RLS
+ runtime role = non-owner / NOSUPERUSER / NOBYPASSRLS
+ READ COMMITTED baseline + explicit row locking
+ opaque random owner revision tokens
+ organization/operation-scoped idempotency records
+ pgx/v5 + pgxpool
```

Real PostgreSQL negative proof remains mandatory before D7 closeout. No generic ORM/repository abstraction, database/schema/role-per-Organization or global `SERIALIZABLE` baseline is admitted.

## 12. Current D7-C gate

The bounded [D7-C Durable Work & External Effects](D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md) candidate is pending operator ratification.

It selects River over the accepted pgx/PostgreSQL stack as the one durable-work engine, uses transaction-bound `InsertTx` instead of a second generic MPC outbox for River-executed handoffs, treats delivery as repeatable rather than exactly-once, persists a pre-dispatch marker before consequential external writes, and routes possible acceptance/crash uncertainty to authoritative reconciliation rather than redispatch.

River scheduling is a wake-up mechanism over durable owner/source state, never business truth/history. River uniqueness/retry may optimize work but cannot replace accepted semantic idempotency or D4 ambiguity rules.

Do not begin D7-D until D7-C is operator-ratified.

Do not begin D8, D9 or Product implementation.