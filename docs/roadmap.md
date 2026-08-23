# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| D6-R2 authority | [Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) + [P4-R1](engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md) + [P8](engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md) — **B00/B01 LOCKED; B10 SUSPENDED** |
| NOTIF-01 | [Design](engineering/rebaseline/D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md) + [Reference Study](engineering/rebaseline/D6-R2-NOTIF-01-REFERENCE-STUDY.md) + [Authority Amendment](engineering/rebaseline/D6-R2-NOTIF-01-AUTHORITY-AMENDMENT.md) — **multi-owner trigger/audience correction under operator review; D2 suspended** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicates H3: curated Product-defined notification kinds from explicit source owners + bounded Organization routing to exact human Principals. Do not continue D2/D3/OAD/B00/B10 before this decision.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; **NOTIF-01 trigger-scope correction OPEN** |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; **NOTIF-01 producer/routing correction OPEN** |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; **NOTIF-01 candidate SUSPENDED** |
| D3 — Communication / Events | ACCEPTED / CLOSED; NOTIF-01 reopen BLOCKED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00/B01 LOCKED; B10 SUSPENDED; NOTIF-01 reference-study gate** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- Global IA remains locked and OP-READ-01 remains resolved.
- Operator evidence falsified Work-only notification origin; D2's Work-only source model is suspended.
- Reference study converges on curated multi-category alerts, explicit audience resolution, and source reread rather than notification-as-truth.
- Recommended H3: explicit producer-owner occurrences → Product-defined NotificationKind → direct or bounded Organization routing → exact human Principal Notifications.
- No generic `AnyDomain` fan-out, Permission-implied responsibility, external broker, or Product implementation.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.
