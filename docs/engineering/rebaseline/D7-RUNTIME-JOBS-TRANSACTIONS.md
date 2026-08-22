# D7 — Runtime / Jobs / Transactions

> **Status:** OPEN / ACTIVE — D7-A RUNTIME ENVELOPE NEXT
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

Questions:

- one deployable modular-monolith class or more than one process for a proved operational reason;
- whether HTTP serving and background execution share one binary/process or only one codebase/deployment unit;
- where composition/wiring lives without becoming business authority;
- shutdown/readiness/health boundaries needed for correctness;
- same-origin topology required by D5-R1/D6 without creating a screen-shaped BFF.

### D7-B — Persistence, isolation and transactions

Decide:

- PostgreSQL connection/transaction ownership;
- structural Organization isolation / RLS model and fail-closed proof;
- transaction boundary for owner command intake;
- concurrency/precondition realization for accepted ETag/revision semantics;
- idempotency record ownership/lifecycle;
- read-projection persistence only where a named consumer requires it.

No schema/table census is frozen before these invariants are established.

### D7-C — Durable work and external effects

Decide the smallest mechanism for:

- post-commit durable work;
- schedules/polling where accepted acquisition requires them;
- owner fact → independent consumer reaction where D3 admits E;
- external-effect dispatch, ambiguous outcomes, authoritative reread and reconciliation;
- retry classification/backoff without changing business meaning;
- atomic handoff between owner state commit and required durable work.

No generic workflow engine, event bus or connector platform is admitted by default.

### D7-D — Authentication/session request-trust realization

Realize the already-decided D5-R1 profile:

- ApplicationSession representation/persistence/expiry/rotation/revocation;
- OIDC Authorization Code exchange and server-side token handling;
- CSRF bootstrap/generation/carriage/rotation;
- Keycloak/provider client/realm topology only to the extent required now;
- A/S audience-bound bearer validation and cache mechanics;
- same-origin browser/API topology.

D7 may choose mechanisms; it may not return the human SPA to bearer-token ownership.

### D7-E — Operability boundary

Decide only what is needed for an executable production target:

- structured logs/traces/metrics and correlation boundaries;
- secrets/configuration ownership and injection;
- migration ownership;
- deployment/readiness/backup expectations required by accepted state and auth;
- real-dependency test/proof seams needed before D8.

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

Absent a concrete falsifier/consumer, D7 does not introduce:

- microservices/service mesh;
- Kubernetes/operator platform work;
- Kafka/NATS or a generic event-stream platform;
- a generic workflow/orchestration engine;
- CQRS/event sourcing as architecture;
- multi-region active-active;
- generic cache layer/Redis;
- generic repository/ORM abstraction;
- second database or analytics warehouse;
- separate BFF business API;
- provider plugin framework;
- real-time/WebSocket infrastructure;
- AI/MCP runtime authority.

## 8. Decision order

D7 proceeds dependency-last:

1. **D7-A Runtime envelope:** fix process/deploy/serving responsibilities and transaction ownership boundary at architectural level.
2. **D7-B Persistence:** establish isolation + transaction invariants before database libraries/schema mechanics.
3. **D7-C Durable work:** establish atomic handoff/retry/reconciliation properties before queue/worker selection.
4. **D7-D Authentication:** select session/OIDC/CSRF/token mechanics under the already-accepted carrier split.
5. **D7-E Operability:** add only deployment/observability/secret/migration mechanics required by the resulting runtime.
6. Run one combined D7 coherence/proof review before closeout.

A later step may force a bounded revisit of an earlier D7 mechanism, but accepted D0–D6 authority reopens only for a material falsifier at the smallest owning stage.

## 9. Exact next D7 work

**D7-A — Runtime Envelope & Transaction Ownership Boundary.**

Derive the smallest serving/process model capable of hosting Product API, Technical Ingress, human same-origin session mediation and durable background work while keeping owner transactions explicit. Compare the modular-monolith/process alternatives against current primary evidence and accepted MPC constraints before selecting any concrete HTTP/runtime/job library.

Do not begin D8, D9 or Product implementation.