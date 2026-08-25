# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed current owners.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Main baseline | **`bdbbef43ed3a5e9d912e67ddac5173024352eaa3` — PR #64 / B10 integrated** |
| Method profile | **local [`engineering-method.md`](development/engineering-method.md) + [`frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md) v2.3** |
| Current prerequisite | **PR #71 — Repository Governance & Context Health Rebaseline — EXECUTION ACTIVE** |
| Health finding | **accepted NOTIF-01 / AuthorizationRequest meaning reached the 106/31 OAD but was not fully consolidated back into canonical D0/D1/D2/D3/D5/D6/D7/D8 owners; historical routing/proof residue remained active** |
| Health execution | **Task 1 retirement audit PASS / NO DELETE; Task 2 selective bootstrap/index routing in progress** |
| B20 increment | **PR #69 — PAUSED / NO P8** |
| Human-operable read-projection prerequisite | **PR #70 — PAUSED at implementation-plan gate until health integration/reanchor** |
| B10 status | **main structure/operator LOCK preserved; correspondence region remains the bounded upstream-repair reopen after health work** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; health work will make heavy Product proof diff-aware without weakening Product-affecting proof** |
| Exact next action | **Finish PR #71 canonical rehome/retirement + proportional-CI proof. Do not render B20 or execute PR #70 Product/OAD changes during health work. After PR #71 is proved, obtain explicit operator merge authorization; then reanchor #70 and #69 from cleaned `main`.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; **NOTIF accepted meaning pending canonical consolidation by PR #71** |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; **NOTIF / AuthorizationRequest later amendments pending canonical consolidation by PR #71** |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; **D2-R6 later amendments pending canonical consolidation by PR #71** |
| D3 — Communication / Events | ACCEPTED / CLOSED; **D3-R3 later amendments pending canonical consolidation by PR #71** |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; **current OAD 106/31; canonical prose owners pending post-amendment consolidation by PR #71** |
| D5-R2 — Operational Read Projection Repair | ACCEPTED / CANONICAL |
| D6 — Frontend | ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110/B10 operator-LOCK evidence preserved |
| D7 — Runtime / Jobs / Transactions | ACCEPTED / CLOSED baseline; later AuthorizationRequest runtime repair pending canonical consolidation by PR #71 |
| D8 — Golden Flows | ACCEPTED / CLOSED / INTEGRATED; later AuthorizationRequest revalidation pending canonical consolidation by PR #71 |
| D8-R2 — GF-02 Operational Read Revalidation | ACCEPTED / PASS |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — repository-health prerequisite before B10 bounded revalidation and B20 continuation** |
| Pre-D9 readiness | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Current result

- PR #64 is integrated in `main`; its B10 operator-LOCK evidence remains current evidence.
- PR #69 exposed the human-operability/read-projection prerequisite and remains paused before P8.
- PR #70 contains the approved Global Maximum design/implementation-plan candidate for that upstream repair and remains paused while repository governance is corrected.
- PR #71 health audit proved that the repository accumulated accepted intermediate amendments without the final **canonical rehome → retirement** step. Current Product wire is newer than several canonical prose parents.
- Task 1 classified every targeted D6-R2/NOTIF/AuthorizationRequest/ADR/plan/proof artifact before deletion. Current LOCKED frontend HTML and current B10/B110/Notifications P8/P9 evidence are protected.
- Normal repository navigation is being restored to `AGENTS → roadmap → applicable method → index → smallest current owner`; Git history remains the archive.
- Product semantics, OAD bytes, 106/31/H-A-S, runtime NONE and D9 implementation block are unchanged by health work.

```text
PR #64 integrated
→ B20 planning exposes upstream gap
→ #69 paused
→ #70 design/plan approved but paused
→ #71 repository-health rebaseline
   → classify
   → canonical rehome
   → retire absorbed intermediates
   → proportional CI
   → prove/review
→ operator-authorized #71 integration
→ reanchor #70
→ bounded B10 correspondence revalidation/P9/re-LOCK
→ resume #69 / B20
```

Return to [`index.md`](index.md) for current-owner routing.
