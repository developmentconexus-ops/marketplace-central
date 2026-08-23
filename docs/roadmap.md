# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| D6-R2 | [Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) + [P4-R1](engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md) + [P8](engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md) — **B00/B01 LOCKED; B10 SUSPENDED** |
| NOTIF-01 | D0-R + D1-R + D2-R + D2-R2 + [D3-R](engineering/rebaseline/D6-R2-NOTIF-01-D3-R-RATIFICATION.md) **ACCEPTED** · [D5-F1](engineering/rebaseline/D6-R2-NOTIF-01-D5-F1-OPERATION-SURFACE-ROUTE-REVERSIBILITY.md) proved route reversibility gap · [D2-R3](engineering/rebaseline/D6-R2-NOTIF-01-D2-R3-ROUTE-UNCONFIGURE-CUTOVER.md) **CANDIDATE** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicates only NOTIF-01 D2-R3 explicit route-unconfigure cutover. D5-R3 table/OAD remain blocked.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; **NOTIF-01 D0-R ACCEPTED** |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; **NOTIF-01 D1-R ACCEPTED** |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; **NOTIF-01 D2-R + D2-R2 ACCEPTED / D2-R3 CANDIDATE** |
| D3 — Communication / Events | ACCEPTED / CLOSED; **NOTIF-01 D3-R ACCEPTED** |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; **NOTIF-01 D5-R3 analysis OPEN / table not ratified** |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00/B01 LOCKED; B10 SUSPENDED; NOTIF-01 D2-R3 gate** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- H3, census, D0-R, D1-R, 14 family contracts, D2-R, D2-R2 and D3-R are accepted.
- D5-R3 derives four Product operations only: `ListMyNotifications`, `UpdateMyNotificationAwarenessState`, `ListNotificationRoutes`, `SetNotificationRoute`; candidate new ordinary Permission: `notifications.manage` only.
- Self Inbox remains human/exact-recipient + current Organization membership, not a new `notifications.read` axis; no cross-Principal Inbox.
- D5-F1 proves routing configuration must be reversible; D2-R3 candidate adds explicit `UNCONFIGURED` temporal revisions without configured-empty, opt-outs, subscriptions or routing DSL.
- Canonical OAD remains unchanged pending D2-R3 + operation-table ratification.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.