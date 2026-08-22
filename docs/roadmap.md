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
| D7-A→D7-E | **OPERATOR-RATIFIED** |
| D7-R1 | **Whole-stage coherence corrections — FABLE REVIEWED / GPT ADJUDICATED / BOUNDED FIX APPLIED** — [authority candidate](engineering/rebaseline/D7-R1-WHOLE-STAGE-COHERENCE.md) |
| Whole-D7 review | **CONVERGED — EXACT-HEAD PROOF PENDING / OPERATOR CLOSEOUT PENDING** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Run fresh exact-head gate after the accepted Fable F-1 recovery-continuity amendment; if green, present D7 for explicit operator closeout. Do not open D8.** |
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
| D7 — Runtime / Jobs / Transactions | **OPEN / ACTIVE — closeout proof + operator decision remaining** |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Accepted D7 baseline

```text
runtime        one Go process/replica · same-origin · in-process River workers
state          PostgreSQL + pgx/v5/pgxpool · structural Organization RLS + composite FKs
transactions   owner-local · READ COMMITTED + explicit locking · opaque revisions · scoped idempotency
effects        River InsertTx · repeat-safe · possible acceptance => authoritative reconciliation
auth H         Keycloak OIDC code+PKCE -> opaque PostgreSQL MPC session + CSRF
auth A/S       Client Credentials -> audience-bound bearer -> explicit A/S Principal binding
HTTP           Chi v5 + oapi-codegen strict server + OAD runtime validation
bytes          private S3-compatible custody + authenticated Go delivery
migrations     tern/v2 for MPC schema + version-matched River migration tool for River schema
observability  JSON slog + OTel traces/metrics over OTLP/HTTP
deploy         immutable OCI image behind trusted TLS edge
recovery       PITR/restore proof + affirmative continuity witness + fail-closed recovery fence
```

Detailed laws, exclusions and proof contracts remain in the D7 owners.

## Whole-D7 review result

Internal review produced five bounded D7-R1 seams. Independent Fable review returned **ACCEPT WITH BOUNDED FIXES** and found one Important gap: the PITR recovery fence was correct once engaged but lacked deterministic automatic arming after rollback. GPT accepted only that blocking finding and D7-R1 now requires timeline continuity to be affirmatively proved using an out-of-rollback-domain witness; absence/mismatch automatically fences external dispatch without a manual restore flag.

Fable's other findings are non-blocking: optional future `chi-server` generation proof and bootstrap-budget hygiene. No D0–D6 reopen or D7 reconstruction is required. Product remains 99/30/H-A-S; Sankhya remains API-Gateway-only; Product implementation remains blocked.

## Closeout boundary

Required sequence:

```text
fresh exact-head gate
-> operator D7 closeout decision
-> only after closeout may roadmap open D8
```

PR #58 remains Draft and unmerged. D8/D9 remain blocked until authority explicitly changes here.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).