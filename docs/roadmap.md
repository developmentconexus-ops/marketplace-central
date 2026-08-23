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
| NOTIF-01 | D0-R + D1-R + D2-R/R2 + [D2-R3](engineering/rebaseline/D6-R2-NOTIF-01-D2-R3-RATIFICATION.md) + [D2-R4](engineering/rebaseline/D6-R2-NOTIF-01-D2-R4-RATIFICATION.md) + [D2-R5](engineering/rebaseline/D6-R2-NOTIF-01-D2-R5-RATIFICATION.md) **ACCEPTED** · [D3-R](engineering/rebaseline/D6-R2-NOTIF-01-D3-R-RATIFICATION.md) **ACCEPTED** · [D3-R1](engineering/rebaseline/D6-R2-NOTIF-01-D3-R1-PRESENTATION-FEED-FORWARD-REVALIDATION.md) + [D3-R2](engineering/rebaseline/D6-R2-NOTIF-01-D3-R2-TYPED-RESULT-CONTINUATION-FEED-FORWARD.md) **PASS** · [D5-F3](engineering/rebaseline/D6-R2-NOTIF-01-D5-F3-GLOBAL-MAXIMUM-OPERATION-PERMISSION-REVIEW.md) findings retained except four-operation count · [D5-F4](engineering/rebaseline/D6-R2-NOTIF-01-D5-F4-RECIPIENT-DISCOVERY-GLOBAL-MAXIMUM.md) **ACCEPTED / OPERATOR-RATIFIED** · [final D5-R3](engineering/rebaseline/D6-R2-NOTIF-01-D5-R3-OPERATION-ADMISSION-TABLE.md) **CANDIDATE / READY FOR OPERATOR REVIEW** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicates only the final NOTIF-01 D5-R3 five-operation admission table. Do not edit the canonical Product OpenAPI first.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; **NOTIF-01 D0-R ACCEPTED** |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; **NOTIF-01 D1-R ACCEPTED** |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; **NOTIF-01 D2-R/R2/R3/R4/R5 ACCEPTED / OPERATOR-RATIFIED as applicable** |
| D3 — Communication / Events | ACCEPTED / CLOSED; **NOTIF-01 D3-R ACCEPTED / D3-R1 + D3-R2 PASS** |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; **NOTIF-01 final D5-R3 five-operation table CANDIDATE / READY FOR OPERATOR REVIEW** |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00/B01 LOCKED; B10 SUSPENDED; NOTIF-01 final D5-R3 gate** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- Counts are consequences, never constraints. The final NOTIF-01 candidate requires **five** public operations and one new `notifications.manage` Permission; if ratified and later authored/proved in the canonical OAD, the derived wire consequence is **104 Product operations / 31 ordinary Permissions**.
- `ListNotificationRouteRecipientCandidates` is H-only, Organization-scoped, owned by IdentityAccess, gated by `notifications.manage`, and exposes only current eligible human `principal_id` + `display_name`; it does not imply `access.read`, disclose `role_keys`/Permissions or authorize the route write.
- `ListMyNotifications` admits bounded filtering by one-or-more Product-defined `NotificationKind` values; no generic filter DSL/source-family identity.
- D2-R5 is accepted: F02/F14 retain only their typed immutable result atoms, and F14 continues through `AuthorizationTargetRef` rather than granting broad `governance.read` merely for Notification navigation.
- D3-R2 PASS proves those D2-R5 atoms/continuation fit the existing source-owner committed-fact `E` topology with no new communication form, owner, broker or generic payload.
- ORG_ROUTED route writes use explicit recipients plus current source-read Permission eligibility (`portfolio.read`, `availability.read`, `economics.read`, `sales.read`, `materialization.read`, `fulfillment.read`, `post_sale.read` as applicable); Permission never selects recipients and candidate discovery never authorizes the write.
- OAD remains **99/30** until the final D5-R3 table is operator-ratified and the separate D5 wire-authoring/proof step executes. No Product implementation.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.