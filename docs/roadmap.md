# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D7 — Runtime / Jobs / Transactions — OPEN / ACTIVE — CLOSEOUT RATIFIED / INTEGRATION PENDING** |
| Accepted baseline | **D0–D6 ACCEPTED / CLOSED** |
| D7 authority | [D7 Runtime / Jobs / Transactions](engineering/rebaseline/D7-RUNTIME-JOBS-TRANSACTIONS.md) — **ACCEPTED / CLOSED authority pending integration** |
| D7-A→D7-E | **OPERATOR-RATIFIED** |
| D7-R1 | **OPERATOR-RATIFIED / ACCEPTED** — [authority](engineering/rebaseline/D7-R1-WHOLE-STAGE-COHERENCE.md) |
| Whole-D7 review | **COMPLETE / CONVERGED** |
| Operator closeout | **APPROVED / RATIFIED — 2026-08-22** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Integrate accepted D7 through PR #58 only after explicit merge authorization; then revalidate `main` and only then transition D8 to NEXT / NOT STARTED.** |
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
| D7 — Runtime / Jobs / Transactions | **OPEN / ACTIVE — CLOSEOUT RATIFIED / PENDING INTEGRATION INTO `main`** |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Accepted D7 authority

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
recovery       PITR/restore + affirmative out-of-rollback continuity witness + fail-closed recovery fence
```

Whole-stage internal review and isolated Fable challenge converged without reopening D0–D6. The sole blocking Fable finding — automatic recovery-fence arming after database rollback — was incorporated into D7-R1 and passed the post-fix repository gate. Fable's `chi-server` proof suggestion and bootstrap-budget note remain non-blocking.

## Integration boundary

Operator ratification accepts D7 as target runtime authority and proof contract, but does **not** by itself authorize Git merge. PR #58 remains the sole D7 integration vehicle.

Until D7 lands in `main`:

- keep the mutable program stage at D7 integration;
- keep D8 and D9 blocked;
- keep Product implementation blocked;
- do not stack D8 work on the unmerged D7 branch.

After an authorized merge, revalidate `main`, PRs and branches first. Only then may the roadmap record D7 as integrated `ACCEPTED / CLOSED` and D8 as `NEXT / NOT STARTED`.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).
