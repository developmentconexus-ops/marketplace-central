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
| D7-C | **Durable Work & External Effects — OPERATOR-RATIFIED** — [authority](engineering/rebaseline/D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md) |
| D7-D | **Authentication / Session / CSRF / Machine Token Realization — CANDIDATE / OPERATOR RATIFICATION PENDING** — [candidate](engineering/rebaseline/D7-D-AUTHENTICATION-SESSION-CSRF.md) |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicate/ratify D7-D; do not open D7-E before that decision.** |
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
| D7 — Runtime / Jobs / Transactions | **OPEN / ACTIVE — D7-A + D7-B + D7-C ratified; D7-D candidate** |
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

**D7-A:** one Go process per replica hosts same-origin delivery, Product API, Technical Ingress, human session mediation and in-process durable workers over PostgreSQL. Owner transactions remain local; external writes happen only after commit.

**D7-B:** explicit Organization scope + composite Organization FKs + transaction-local PostgreSQL `ENABLE/FORCE RLS` under a non-owner/non-`BYPASSRLS` runtime role; `READ COMMITTED` + explicit locking; opaque revision tokens; org/operation-scoped idempotency; `pgx/v5` + `pgxpool`. Real PostgreSQL negative proof remains required.

**D7-C:** River is the bounded PostgreSQL durable-work engine. `InsertTx` is the owner-state → durable-work atomic handoff, without a second generic MPC outbox. Delivery is repeatable; semantic idempotency is correctness. Possible external acceptance after crash/timeout never redispatches and instead reconciles against authoritative source truth. Scheduler/job state is not business truth.

## D7-D candidate

[D7-D](engineering/rebaseline/D7-D-AUTHENTICATION-SESSION-CSRF.md) proposes:

```text
Keycloak first OIDC/OAuth provider
human: confidential Authorization Code + PKCE S256
       go-oidc + x/oauth2
       verified (issuer, sub) -> fresh opaque MPC session
       discard human OIDC token set after callback
       30m idle / 8h absolute server-side session
       X-CSRF-Token synchronizer + net/http CrossOriginProtection
machine: Keycloak Client Credentials/service account
         audience includes https://conexus.fun
         jwx/v3 trusted-JWKS verification
         explicit machine-client -> A/S Principal binding
```

No browser bearer, human refresh-token cache, Redis session store, IdP role→Permission mapping, realm-per-Organization, wildcard CORS or generic IAM engine is admitted.

D8–D9 remain blocked. Product implementation remains blocked until D9. Reopen D0–D6 only for a material falsifier at the smallest owner.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).