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
| D7-B | **PostgreSQL Isolation & Transaction Realization — OPERATOR-RATIFIED** — [authority](engineering/rebaseline/D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md) |
| D7-C | **Durable Work & External Effects — CANDIDATE / OPERATOR RATIFICATION PENDING** — [candidate](engineering/rebaseline/D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md) |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicate/ratify D7-C; do not open D7-D before that decision.** |
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
| D7 — Runtime / Jobs / Transactions | **OPEN / ACTIVE — D7-A + D7-B ratified; D7-C candidate** |
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

## D7 accepted decisions

**D7-A:** one Go application process per replica hosts same-origin delivery, Product API, Technical Ingress, human session mediation and in-process durable-worker runner over PostgreSQL. Owner consequential intake commits only owner state plus required technical correctness records; external writes happen after commit. No distributed business transaction spans semantic owners.

**D7-B:** organization-owned state/evidence carries explicit Organization scope; composite Organization FKs + transaction-local scope + PostgreSQL `ENABLE/FORCE RLS` under a non-owner/non-`BYPASSRLS` runtime role provide structural isolation. Baseline is `READ COMMITTED` + explicit row locking, opaque random owner revision tokens, organization/operation-scoped idempotency and `pgx/v5` + `pgxpool`. Real PostgreSQL negative proof remains required before D7 closeout.

## D7-C candidate

[D7-C](engineering/rebaseline/D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md) proposes River as the one PostgreSQL-backed durable-work engine over the accepted pgx stack. Owner state + required jobs use River `InsertTx` atomically, so no second generic MPC outbox is baseline. Delivery remains repeatable; River uniqueness/retry is optimization, never business correctness.

External-effect workers persist a pre-dispatch marker before any network write. If possible acceptance exists after timeout/crash, the write is **not** redispatched; durable authoritative reconciliation runs instead. River scheduling is only a wake-up mechanism over durable owner/source state, never business truth/history.

D8–D9 remain blocked. Product implementation remains blocked until D9. Reopen D0–D6 only for a material falsifier at the smallest owner.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).
