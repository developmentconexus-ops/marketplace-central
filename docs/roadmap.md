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
| NOTIF-01 | D0-R + D1-R + D2-R/R2 + [D2-R3](engineering/rebaseline/D6-R2-NOTIF-01-D2-R3-RATIFICATION.md) + [D2-R4](engineering/rebaseline/D6-R2-NOTIF-01-D2-R4-RATIFICATION.md) + [D3-R](engineering/rebaseline/D6-R2-NOTIF-01-D3-R-RATIFICATION.md) **ACCEPTED** · [D3-R1](engineering/rebaseline/D6-R2-NOTIF-01-D3-R1-PRESENTATION-FEED-FORWARD-REVALIDATION.md) **PASS** · [D5-F3](engineering/rebaseline/D6-R2-NOTIF-01-D5-F3-GLOBAL-MAXIMUM-OPERATION-PERMISSION-REVIEW.md) findings retained except four-operation count · [D5-F4](engineering/rebaseline/D6-R2-NOTIF-01-D5-F4-RECIPIENT-DISCOVERY-GLOBAL-MAXIMUM.md) **ACCEPTED / OPERATOR-RATIFIED** · [D2-R5](engineering/rebaseline/D6-R2-NOTIF-01-D2-R5-TYPED-RESULT-CONTINUATION.md) **CANDIDATE** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicates only NOTIF-01 D2-R5 typed-result/requester-continuation repair. D5-R3 final five-operation table/OAD remain blocked.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; **NOTIF-01 D0-R ACCEPTED** |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; **NOTIF-01 D1-R ACCEPTED** |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; **NOTIF-01 D2-R/R2/R3/R4 ACCEPTED; D2-R5 CANDIDATE** |
| D3 — Communication / Events | ACCEPTED / CLOSED; **NOTIF-01 D3-R ACCEPTED / D3-R1 PASS; D3-R2 blocked by D2-R5** |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; **NOTIF-01 D5-F4 RECIPIENT-DISCOVERY GLOBAL MAXIMUM ACCEPTED; final D5-R3 table blocked by D2-R5/D3-R2** |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00/B01 LOCKED; B10 SUSPENDED; NOTIF-01 D2-R5 gate** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- Counts are consequences, never constraints. D5-F4 supersedes only D5-F3's four-operation conclusion: the final NOTIF-01 direction now requires **five** public operations, including the purpose-bounded `ListNotificationRouteRecipientCandidates`; the future wire consequence is **104 Product operations / 31 ordinary Permissions** if the final table is ratified after the remaining gates.
- `ListNotificationRouteRecipientCandidates` is H-only, Organization-scoped, owned by IdentityAccess, gated by `notifications.manage`, and exposes only current eligible human `principal_id` + `display_name`; it does not imply `access.read` or disclose `role_keys`/Permissions.
- `ListMyNotifications` must admit bounded filtering by one-or-more Product-defined `NotificationKind` values; no generic filter DSL/source-family identity.
- ORG_ROUTED route writes use explicit recipients plus current source-read Permission eligibility (`portfolio.read`, `availability.read`, `economics.read`, `sales.read`, `materialization.read`, `fulfillment.read`, `post_sale.read` as applicable); Permission never selects recipients and candidate discovery never authorizes the write.
- D2-R5 candidate adds only typed immutable F02/F14 result atoms and changes F14 requester continuation from `AuthorizationDecisionRef` to `AuthorizationTargetRef`, avoiding broad `governance.read` grants.
- OAD remains 99/30 until D2-R5 + bounded D3 feed-forward + final D5-R3 five-operation table ratification. No Product implementation.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.