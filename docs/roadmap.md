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
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **D7-A — derive the smallest runtime/process topology and transaction-ownership boundary before selecting concrete HTTP/job/database mechanisms.** |
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
| D7 — Runtime / Jobs / Transactions | **OPEN / ACTIVE — D7-A Runtime Envelope next** |
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

D7 must preserve:

- Organization isolation as a structural invariant;
- owner-defined Q/C/E/P meaning;
- accepted Product idempotency/ETag/Problem grammar;
- no-blind-retry and authoritative reconciliation for ambiguous external effects;
- provider DTO/protocol isolation;
- Sankhya sanctioned API Gateway with no Direct Oracle fallback;
- human session + CSRF with no browser OIDC bearer ownership;
- A/S bearer resolution independent from Membership/Permission/owner/Governance gates;
- honest unknown/partial/unavailable state.

Pre-vetted candidates such as modular-monolith class, pgx/pgxpool, PostgreSQL RLS, River, OpenTelemetry/OTLP/slog, sqlc, tern, Keycloak and real-dependency test tooling remain **candidates only** until D7 proves a current consumer/property and adjudicates the smallest sufficient mechanism.

## D7-A next gate

First derive the minimum runtime envelope before dependency selection:

1. number/responsibility of server and background processes;
2. Product API vs Technical Ingress serving boundary;
3. same-origin human browser/session mediation without a screen-shaped BFF;
4. owner command transaction ownership and post-commit durable handoff boundary;
5. shutdown/readiness/health responsibilities needed for correctness.

Use current primary technology evidence only when a concrete D7-A mechanism comparison requires it. Do not infer target runtime from removed code or historical ADR implementation detail.

D8–D9 remain blocked. Product implementation remains blocked until D9. Reopen accepted D0–D6 authority only for a material falsifier at the smallest owning stage.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).
