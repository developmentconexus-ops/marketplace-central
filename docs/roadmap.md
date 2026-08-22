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
| Exact next action | **Operator adjudicate/ratify D7-B; do not open D7-C before that decision.** |
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
| D5 — API | ACCEPTED / CLOSED |
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **OPEN / ACTIVE — D7-A ratified; D7-B candidate** |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## D7 accepted inputs

```text
Product operations     99
ordinary Permissions   30
Principal kinds        H / A / S only
stable origin          https://conexus.fun
Go backend             canonical business execution
PostgreSQL             canonical MPC-owned state store
active runtime         NONE
```

Human browser authentication remains server-side OIDC → Secure HttpOnly application session + CSRF; A/S remain Client Credentials bearer clients. D7 realizes mechanics only and may not move Product/frontend/business authority.

## D7-A accepted

One Go application process per replica hosts same-origin delivery, Product API, Technical Ingress, human session mediation and in-process durable-worker runner over PostgreSQL. Owner consequential intake commits only owner state plus required technical correctness records; external writes happen after commit. No distributed business transaction spans semantic owners.

## D7-B candidate

[D7-B](engineering/rebaseline/D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md) proposes:

```text
organization_id on org-owned state/evidence
+ composite Organization FKs
+ transaction-local scope
+ PostgreSQL ENABLE/FORCE RLS
+ non-owner / NOSUPERUSER / NOBYPASSRLS runtime role
+ READ COMMITTED + explicit row locking
+ opaque random owner revision token
+ org+operation-scoped idempotency record
+ pgx/v5 + pgxpool
```

Principal-self and technical-routing modes are bounded bootstrap exceptions, never generic cross-tenant business access. Real PostgreSQL negative proof is required for RLS, scope reset, FK isolation, worker scope, revision concurrency and idempotency. ORM/repository abstraction, database-per-Organization, role-per-Organization, global `SERIALIZABLE` and database portability are rejected absent a consumer; `sqlc` and `tern` remain deferred.

D8–D9 remain blocked. Product implementation remains blocked until D9. Reopen D0–D6 only for a material falsifier at the smallest owner.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).
