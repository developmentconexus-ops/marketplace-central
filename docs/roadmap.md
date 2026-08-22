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
| D7-A | **OPERATOR-RATIFIED** |
| D7-B | **OPERATOR-RATIFIED** — [authority](engineering/rebaseline/D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md) |
| D7-C | **OPERATOR-RATIFIED** — [authority](engineering/rebaseline/D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md) |
| D7-D | **OPERATOR-RATIFIED** — [authority](engineering/rebaseline/D7-D-AUTHENTICATION-SESSION-CSRF.md) |
| D7-E | **OPERATOR-RATIFIED** — [authority](engineering/rebaseline/D7-E-OPERABILITY-DEPLOYMENT-PROOF.md) |
| Whole-D7 review | **OPEN / ACTIVE — coherence + executable-proof/adversarial review** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Run whole-D7 coherence + executable-proof/adversarial review over D7-A→D7-E; then independent Fable challenge and GPT adjudication. Do not open D8.** |
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
| D7 — Runtime / Jobs / Transactions | **OPEN / ACTIVE — D7-A→D7-E ratified; whole-D7 review open** |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Accepted D7 baseline

```text
Product            99 operations · 30 Permissions · H/A/S
backend            Go modular-monolith process per replica
state              PostgreSQL + pgx/v5 + pgxpool
isolation          explicit Organization + composite FKs + ENABLE/FORCE RLS
transactions       owner-local · READ COMMITTED + explicit locking
work/effects       River InsertTx · repeat-safe · possible acceptance => reconcile, never redispatch
auth H             Keycloak OIDC code+PKCE -> opaque PostgreSQL MPC session + CSRF
auth A/S           Client Credentials -> audience-bound bearer -> explicit A/S Principal binding
HTTP               Chi v5 + oapi-codegen strict server + OAD runtime validation
bytes              private S3-compatible custody + authenticated Go delivery
migrations         tern/v2 with separate migration owner
observability      JSON slog + OTel traces/metrics over OTLP/HTTP
deploy             one immutable OCI app image behind trusted TLS edge
recovery           PostgreSQL PITR + Keycloak subject continuity + binary integrity restore proof
```

Detailed laws, exclusions and falsifiers remain in the D7 owner documents; this roadmap does not duplicate them.

## Whole-D7 review boundary

The review may challenge only material contradictions, missing seams, duplicated authority, impossible proof obligations or hidden changes to accepted D0–D6 meaning. It may produce bounded D7 fixes. It may not begin D8, D9 or Product implementation.

Required sequence:

```text
internal whole-D7 coherence/adversarial review
-> executable-proof-plan review
-> isolated independent Fable challenge
-> GPT adjudication
-> bounded fixes + fresh exact-head gate if required
-> operator D7 closeout decision
```

D8 remains blocked until D7 closes. Reopen D0–D6 only for a material falsifier at the smallest owning authority.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).