# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Integrated checkpoints | **PR #61 Product/frontend checkpoint → `b54d17bfe6d794645d198a9160f4a2a1c63647e8`; PR #62 methodology adoption → `689dab34b0a756cbd7c790a6c5277d887ced0b4c`** |
| Method profile | **`developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535` → `METHOD + FRONTEND-METHOD`** |
| Current acceptance increment | **B10 — Preparação — bounded post-methodology revalidation / P8 CANDIDATE; NOT LOCKED** |
| LOCKED frontend blocks | **B00 · B01 · B00-R2 · B11 · B12 · B110** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Revalidate B10 P6/P7 against the pinned frontend method, prove the corrected browser-operable P8 HTML, then operator operates it and chooses REVISE / UPSTREAM FINDING / explicit LOCK. Do not begin B10 P9 before LOCK.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; NOTIF-01 D0-R ACCEPTED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; D1-R2 PASS |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; D2-R6 ACCEPTED |
| D3 — Communication / Events | ACCEPTED / CLOSED; D3-R3 ACCEPTED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; D5-R6 106/31 PROVED; D5-R7 W1 REPAIR ACCEPTED |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110 LOCKED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; D7-R ACCEPTED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; D8-R ACCEPTED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B10 P8 CANDIDATE / OPERATOR ADJUDICATION NEXT AFTER GREEN** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- Methodology adoption is integrated; accepted Product/frontend authority remains intact.
- B10 is the only active Product-experience increment. Its existing search-first structure is being falsified against the pinned `FRONTEND-METHOD.md`; methodology adoption alone does not reopen Product authority.
- Product **106/31/H-A-S**, AuthorizationRequest/W1, D7-R, D8-R, OAD source hygiene and existing operator LOCKs remain unchanged.
- B10 may advance to P9 only after explicit operator LOCK of its operated P8 candidate.

```text
B10 P6/P7 revalidation → corrected P8 HTML → operator adjudication → LOCK only if explicit → B10 P9/P10 → integrate increment → next Bxx
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
