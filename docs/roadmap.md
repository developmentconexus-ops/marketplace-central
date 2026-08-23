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
| NOTIF-01 evidence | [Design](engineering/rebaseline/D6-R2-NOTIF-01-NOTIFICATION-ARCHITECTURE-DESIGN.md) + [Reference Study / Census](engineering/rebaseline/D6-R2-NOTIF-01-REFERENCE-STUDY.md) — **H3 + Trigger/Audience Census OPERATOR-APPROVED** |
| NOTIF-01 D0-R | [Trigger-Scope Correction](engineering/rebaseline/D6-R2-NOTIF-01-D0-R-TRIGGER-SCOPE-CORRECTION.md) — **OPERATOR-APPROVED / ACCEPTED** |
| NOTIF-01 D1-R | [Producer & Routing Boundary Correction](engineering/rebaseline/D6-R2-NOTIF-01-D1-R-PRODUCER-ROUTING-BOUNDARY-CORRECTION.md) — **CANDIDATE; RATIFICATION DEFERRED PENDING FAMILY CONTRACTS** |
| NOTIF-01 family contracts | [14 Notification Family Semantic Contracts](engineering/rebaseline/D6-R2-NOTIF-01-NOTIFICATION-FAMILY-SEMANTIC-CONTRACTS.md) — **CANDIDATE / CURRENT OPERATOR GATE** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicates the 14 NOTIF-01 family semantic contracts plus two coherence findings: MATERIALIZATION_ATTENTION uses one materialization/backoffice ORG_ROUTED audience, and same-consequence SALE_ATTENTION→POST_SALE_ATTENTION gets bounded per-recipient precedence. Do not ratify D1-R or begin D2-R/D3/OAD/B00/B10 before this gate closes.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; **NOTIF-01 D0-R ACCEPTED** |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; **NOTIF-01 D1-R CANDIDATE / FAMILY CONTRACT REVIEW OPEN** |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; **NOTIF-01 prior candidate SUSPENDED / D2-R BLOCKED** |
| D3 — Communication / Events | ACCEPTED / CLOSED; NOTIF-01 reopen BLOCKED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00/B01 LOCKED; B10 SUSPENDED; NOTIF-01 family-semantic gate** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- Global IA remains locked and OP-READ-01 remains resolved.
- NOTIF-01 H3 + Trigger/Audience Census + D0-R corrected Product scope are operator-approved.
- D1-R remains candidate: **10 explicit producer-owner edges → 14 awareness families → one Personal Notifications supporting owner**.
- The 14 family contracts now define human meaning, birth/not-this boundaries, audience, deep link, repeat and suppression semantics before D1-R ratification.
- Global coherence review finds no fifteenth family, but proposes two D1-R clarifications: one stable materialization/backoffice audience for `MATERIALIZATION_ATTENTION`, and bounded same-consequence `POST_SALE_ATTENTION` precedence over `SALE_ATTENTION` per recipient.
- Source owners retain attention-occurrence truth; Personal Notifications owns bounded awareness/routing/dedup semantics only.
- D2-R/D3/OAD/bell/Inbox/settings remain blocked until family contracts and then D1-R are operator-ratified.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.
