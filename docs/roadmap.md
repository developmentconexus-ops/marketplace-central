# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D7 — Runtime / Jobs / Transactions — OPEN / ACTIVE** |
| Accepted baseline | **D0–D6 ACCEPTED / CLOSED** |
| D7 authority | [D7 Runtime / Jobs / Transactions](engineering/rebaseline/D7-RUNTIME-JOBS-TRANSACTIONS.md) |
| D7-A | **Runtime Envelope & Transaction Ownership Boundary — OPERATOR-RATIFIED** |
| D7-B | **PostgreSQL Isolation & Transaction Realization — CANDIDATE / OPERATOR RATIFICATION PENDING** — [candidate](engineering/rebaseline/D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md) |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicate/ratify D7-B before opening D7-C durable-work/effect realization.** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED |
| D3 — Communication / Events | ACCEPTED / CLOSED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED — bounded D5-R1 browser-auth carrier correction integrated and proved through D6 |
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **OPEN / ACTIVE — D7-A operator-ratified; D7-B candidate pending operator ratification** |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Accepted input to D7

D7 inherits, without reopening by preference:

```text
Product operations     99
ordinary Permissions   30
Principal kinds        H / A / S only
stable origin          https://conexus.fun
Go backend             canonical business execution
PostgreSQL             canonical MPC-owned state store
active runtime         NONE
```

Human/machine authentication remains:

```text
H browser  -> server-side OIDC login -> Secure HttpOnly application session + CSRF on unsafe requests
A / S      -> Client Credentials -> audience-bound bearer
```

Frontend realization remains React + TypeScript strict, TanStack Query, TanStack Router, `openapi-typescript` and `openapi-fetch`; D7 realizes the server/runtime side without moving Product or frontend authority.

## D7 boundary

D7 owns technical realization for serving/process topology, PostgreSQL structural Organization isolation and transaction boundaries, durable jobs/scheduling/outbox/effect handoff, session/CSRF/OIDC and machine-token mechanics, secrets, observability, migrations and deployment topology required by the admitted transports.

D7 must preserve Organization isolation, owner-defined Q/C/E/P meaning, Product idempotency/ETag/Problem grammar, no-blind-retry/reconciliation, provider protocol isolation, sanctioned Sankhya API Gateway, human session + CSRF, A/S bearer separation and honest unknown/partial/unavailable state.

Pre-vetted candidates such as modular-monolith class, pgx/pgxpool, PostgreSQL RLS, River, OpenTelemetry/OTLP/slog, sqlc, tern, Keycloak and real-dependency test tooling remain candidates until their exact D7 slice establishes a current consumer/property.

## D7-A accepted authority

The operator ratified the minimum runtime envelope:

```text
one Go application process per replica
  -> same-origin frontend/static delivery
  -> Product API boundary
  -> Technical Non-Product Ingress boundary
  -> H session/CSRF/OIDC mediation
  -> in-process durable worker runner
  -> PostgreSQL
```

Owner consequential command intake uses one owner-local PostgreSQL transaction for canonical owner state plus required idempotency/audit/durable-handoff records. External writes happen only after commit and convergence remains authoritative-reread/reconciliation-driven. Cross-owner business state is never updated through a distributed/shared transaction.

Separate API/worker processes are deferred until measured resource/failure/scaling/security/deployment evidence requires them; microservices remain rejected absent an independent consumer. River remains a D7-C candidate, not a D7-A dependency selection.

## D7-B candidate

The current D7-B candidate combines complementary structural controls rather than relying on one tenant mechanism:

```text
explicit organization_id on organization-owned rows
+ composite (organization_id, id) referential integrity
+ transaction-local scope
+ PostgreSQL ENABLE/FORCE RLS
+ runtime role: non-owner / NOSUPERUSER / NOBYPASSRLS
```

All scoped database work runs inside an explicit transaction so Organization/Principal/technical-routing scope is transaction-local and cannot leak through pooled connection reuse. Organization business/evidence tables accept only exact Organization scope; principal-self and technical-routing are bounded platform modes for `/access-context` and technical bootstrap only.

The candidate uses `READ COMMITTED` by default with explicit row locking/constraints for protected owner mutations; global `SERIALIZABLE` is rejected without a concrete invariant requiring its retry cost. Strong Product ETags map to persisted opaque random owner revision tokens and never expose `xmin`, counters or timestamps. Required idempotency uses an organization + operation scoped technical record with a semantic fingerprint and exact replay/result correlation in the owner transaction.

`pgx/v5` + `pgxpool` is the candidate direct PostgreSQL primitive. Generic ORM/repository abstraction and database portability are rejected; `sqlc` remains deferred until concrete queries exist and `tern` remains D7-E migration mechanism territory.

D7-B requires real PostgreSQL negative proof for RLS, scope reset, composite FK isolation, worker scope, concurrency, ETag and idempotency before ratification/closeout claims. Mocks are insufficient for those properties.

D8–D9 remain blocked. Product implementation remains blocked until D9. Reopen accepted D0–D6 authority only for a material falsifier at the smallest owning stage.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).
