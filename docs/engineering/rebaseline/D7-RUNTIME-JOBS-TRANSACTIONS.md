# D7 — Runtime / Jobs / Transactions

> **Status:** OPEN / ACTIVE — D7-A OPERATOR-RATIFIED / D7-B NEXT
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

Decide PostgreSQL connection/transaction ownership, structural Organization isolation/RLS, owner command transactions, ETag/revision realization, idempotency lifecycle and only justified read projections.

### D7-C — Durable work and external effects

Decide post-commit durable work, schedules/polling, admitted E reactions, external-effect dispatch/reconciliation, retry classification and atomic handoff between owner state and required work.

### D7-D — Authentication/session request-trust realization

Realize ApplicationSession, OIDC exchange, CSRF, Keycloak/provider topology and A/S token validation under the already-decided D5-R1 carrier split.

### D7-E — Operability boundary

Decide only the structured observability, secrets/configuration, migrations, deployment/readiness/backup and real-dependency proof mechanics needed by the resulting runtime.

## 5. Candidate set is evidence, not selection

Accepted architecture has already identified candidates worth testing, including:

```text
modular-monolith class
pgx / pgxpool
PostgreSQL structural isolation / RLS
River-first durable-work falsification
OpenTelemetry / OTLP / slog
sqlc
tern
Keycloak
real PostgreSQL / real-dependency test tooling
```

None is selected merely by appearing here. D7 must test each against an actual accepted consumer/property, current primary evidence and the smallest sufficient alternative. Exact versions are implementation-manifest concerns unless a version-specific property is architectural.

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
2. **D7-B Persistence:** isolation + transaction invariants before database/schema mechanics — **NEXT**.
3. **D7-C Durable work:** atomic handoff/retry/reconciliation before worker selection.
4. **D7-D Authentication:** session/OIDC/CSRF/token mechanics under the accepted carrier split.
5. **D7-E Operability:** only required deployment/observability/secret/migration mechanics.
6. One combined D7 coherence/proof review before closeout.

## 9. Current primary evidence for D7-A

Revalidated on 2026-08-21:

- Go `net/http.Server.Shutdown` provides graceful HTTP shutdown without interrupting active connections: <https://pkg.go.dev/net/http#Server.Shutdown>.
- Go `ServeMux` supports method, host and path routing patterns, so route expressibility alone does not require a framework: <https://pkg.go.dev/net/http#ServeMux>.
- Go `signal.NotifyContext` provides signal-driven context cancellation for process lifecycle: <https://pkg.go.dev/os/signal#NotifyContext>.
- Current River documentation shows a PostgreSQL-backed client can start worker loops in background goroutines inside the application process, can be insert-only when queues are omitted, supports graceful worker stop, and supports transaction-bound `InsertTx`: <https://pkg.go.dev/github.com/riverqueue/river> and <https://riverqueue.com/>.

River remains only a D7-C candidate here. Its evidence matters to D7-A because it proves that durable PostgreSQL work does **not** force a separate worker service/process and that a future split remains possible without a business-authority rewrite.

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
  ├─ scheduler coordinator seam (mechanism D7-C)
  └─ PostgreSQL connection pool
```

Binding laws:

- **one process does not mean one business module**: accepted owner/context boundaries remain explicit inside the modular monolith;
- Product API and Technical Ingress use separate route/middleware composition boundaries even though they share listener/process infrastructure;
- no internal owner-to-owner HTTP merely for symmetry: Q/C/P inside the monolith use explicit owner-owned application ports/interfaces; E remains durable only where D3 admits an independent reaction;
- browser/static + API remain one public origin. A separate CDN/edge/runtime is not baseline and requires an operational consumer;
- no separate screen-shaped BFF appears: server-side OIDC/session mediation protects the same Product API;
- correctness must not depend on exactly one replica. The process may later be replicated if D7-C coordination and D7-D session persistence are safe across replicas;
- no HTTP router/framework is selected by D7-A. Current Go routing capability means a dependency must prove a concrete missing property before admission.

### 10.3 Owner transaction ownership

For a consequential owner command:

```text
request / admitted internal C
  -> authenticate + scope/Permission/precondition gates
  -> begin one owner command transaction
       -> claim/validate idempotency identity when required
       -> mutate only that owner's canonical business state
       -> write required audit/intake/outcome metadata
       -> atomically persist required durable handoff/job/outbox record
  -> COMMIT
  -> return accepted/pending/domain outcome
  -> background worker performs external effect after commit
  -> authoritative reread/reconciliation determines convergence
```

Laws:

- one transaction may include platform records required for correctness (idempotency, audit, durable handoff) but **must not mutate another semantic owner's private business state**;
- cross-owner business changes do not use a distributed transaction. They cross explicit Q/C/E contracts and each owner commits its own meaning;
- external network writes are never performed while holding the owner PostgreSQL transaction open;
- potentially accepted external effects are never blindly replayed after transport ambiguity;
- external reads/preflight evidence, when required, occur outside the write transaction; local revision/scope/precondition checks are revalidated before commit as needed;
- ordinary Q reads do not acquire a write transaction. Snapshot/read-only transaction semantics remain D7-B only where a concrete multi-read consistency property requires them;
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

Split API and worker process roles only if a real proof establishes at least one of:

- background work causes unacceptable API resource starvation or shutdown coupling;
- API and worker need materially different scaling envelopes;
- a security/exposure boundary requires separate processes;
- D8 composed-flow proof cannot satisfy graceful drain/recovery with the combined process;
- deployment constraints make combined lifecycle unsafe or operationally unmanageable.

Preference for process separation, microservices or framework convention is not a reopen trigger.

## 11. Exact next D7 work

Proceed to **D7-B — PostgreSQL Isolation & Transaction Realization**: derive the structural Organization/RLS model, transaction context propagation, idempotency storage boundary and ETag/revision enforcement before choosing pgx/sqlc/tern details.

Do not begin D8, D9 or Product implementation.
