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
| NOTIF-01 | [D0-R](engineering/rebaseline/D6-R2-NOTIF-01-D0-R-TRIGGER-SCOPE-CORRECTION.md) **ACCEPTED** · [D1-R](engineering/rebaseline/D6-R2-NOTIF-01-D1-R-PRODUCER-ROUTING-BOUNDARY-CORRECTION.md) + [Semantic Closure](engineering/rebaseline/D6-R2-NOTIF-01-D1-R-SEMANTIC-CLOSURE.md) **FINAL CANDIDATE** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicates final NOTIF-01 D1-R: 10 explicit producer edges, 14 approved family contracts, three audience strategies, bounded A3 routing, Work-replacement suppression and Sale→Post-Sale precedence. Do not begin D2-R/D3/OAD/B00/B10 first.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; **NOTIF-01 D0-R ACCEPTED** |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; **NOTIF-01 D1-R FINAL CANDIDATE** |
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
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00/B01 LOCKED; B10 SUSPENDED; NOTIF-01 D1-R final gate** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- H3, Trigger/Audience Census, D0-R and all 14 family semantic contracts are operator-approved.
- D1-R final candidate: **10 explicit producer edges → 14 awareness families → one Personal Notifications owner**.
- DIRECT_SOURCE and OWNER_DERIVED recipient meaning stays with producers; ORG_ROUTED exact-human configuration belongs to Personal Notifications.
- Two bounded per-recipient suppression laws only: source alert→assigned Work; Sale attention→richer Post-Sale continuation for the same consequence/recipient.
- No `AnyDomain` fan-out, Permission-implied responsibility, hidden subtype routing, generic subscription/dedup DSL, broker or Product implementation.

```text
D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent gate lands before the next. Return to [`index.md`](index.md) for task routing.
