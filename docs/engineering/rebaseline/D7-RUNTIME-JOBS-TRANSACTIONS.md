# D7 — Runtime / Jobs / Transactions

> **Status:** OPEN / ACTIVE — D7-A + D7-B + D7-C OPERATOR-RATIFIED / D7-D CANDIDATE PENDING RATIFICATION
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

### D7-C — Durable work and external effects

Owning accepted authority: [D7-C Durable Work & External Effects](D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md).

### D7-D — Authentication/session request-trust realization

Owning candidate: [D7-D Authentication / Session / CSRF / Machine Token Realization](D7-D-AUTHENTICATION-SESSION-CSRF.md).

### D7-E — Operability boundary

Decide only the structured observability, secrets/configuration, migrations, deployment/readiness/backup and real-dependency proof mechanics needed by the resulting runtime.

## 5. Candidate set is evidence, not selection

Accepted architecture identified candidates worth testing, including modular-monolith class, pgx/pgxpool, PostgreSQL RLS, River, OpenTelemetry/OTLP/slog, sqlc, tern, Keycloak and real-dependency test tooling. D7-A selected the modular-monolith process envelope; D7-B selected pgx/pgxpool + PostgreSQL structural isolation; D7-C selected River. Keycloak and auth/session libraries are under D7-D adjudication. Remaining candidates stay unselected until their owning slice proves a current consumer/property.

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
3. **D7-C Durable work:** atomic handoff/retry/reconciliation before worker selection — **OPERATOR-RATIFIED**.
4. **D7-D Authentication:** session/OIDC/CSRF/token mechanics under the accepted carrier split — **CANDIDATE / OPERATOR RATIFICATION PENDING**.
5. **D7-E Operability:** only required deployment/observability/secret/migration mechanics.
6. One combined D7 coherence/proof review before closeout.

## 9. D7-A accepted — Runtime Envelope & Transaction Ownership Boundary

> **OPERATOR-RATIFIED:** one Go application process per replica, one public application origin, Product API + Technical Ingress + human session mediation + durable worker runner in the same process; owner business transactions stay local to one semantic owner and all consequential external effects happen outside the database transaction after an atomic durable handoff.

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

No internal owner-to-owner HTTP is introduced merely for symmetry; Q/C/P enter explicit owner application boundaries and E is durable only where D3 admits an independent reaction. Browser/static + API remain one public origin. A separate API/worker process split reopens only for measured resource, failure, scaling, security or deployment need.

Owner consequential intake uses one owner-local transaction for canonical owner state plus required idempotency/audit/durable handoff. External network writes happen only after commit. No distributed business transaction spans semantic owners.

Graceful shutdown marks not-ready, stops new HTTP and new job claiming, drains bounded active work, leaves durable work recoverable, then closes telemetry/database resources. Provider/Sankhya unavailability alone is not application unreadiness.

## 10. D7-B accepted — PostgreSQL Isolation & Transaction Realization

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

## 11. D7-C accepted — Durable Work & External Effects

[D7-C](D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md) is **OPERATOR-RATIFIED**.

Accepted baseline:

```text
River over pgx/PostgreSQL
+ InsertTx atomic owner-state -> durable-work handoff
+ no second generic MPC outbox for River work
+ repeat-safe delivery / no exactly-once claim
+ owner semantic idempotency = correctness
+ persisted pre-dispatch marker before consequential external writes
+ possible acceptance => no redispatch
+ authoritative reconciliation
+ scheduled/periodic River work = wake-up mechanism only
```

River completion, uniqueness, retry, rescue and scheduler state are technical optimization/state only. They never become business truth or history authority. Crash/timeout after possible external dispatch moves to reconciliation, never generic retry.

Real PostgreSQL + River crash/restart, transactional-enqueue, duplicate-delivery and cross-Organization proof remains mandatory before D7 closeout.

## 12. Current D7-D gate

The bounded [D7-D Authentication / Session / CSRF / Machine Token Realization](D7-D-AUTHENTICATION-SESSION-CSRF.md) candidate is pending operator ratification.

It selects Keycloak as the first provider while retaining a standards-based OIDC boundary; server-side Authorization Code + required PKCE S256 through `go-oidc` + `x/oauth2`; one-time PostgreSQL login transaction; human OIDC token-set discard after verified `(issuer, sub)` binding; an opaque hashed server-side MPC session with bounded idle/absolute lifetime; session-bound synchronizer `X-CSRF-Token` plus Go `CrossOriginProtection`; and strict audience/issuer/JWKS machine-token validation through `jwx/v3` with explicit client→A/S Principal binding.

No browser bearer, IdP role→Permission mapping, Redis session store, persistent human refresh-token cache, realm-per-Organization or generic IAM engine is admitted.

Do not begin D7-E until D7-D is operator-ratified.

Do not begin D8, D9 or Product implementation.