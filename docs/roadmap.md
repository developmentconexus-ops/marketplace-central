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
| Method profile | **`developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535`** |
| Current repository-governance increment | **PR #65 — proportional CI verification — independent of Product/B10** |
| Parallel Product candidate | **PR #64 — B10 Preparação — UNMERGED; not integrated authority on this branch** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Prove PR #65 Draft quick lane, transition the exact candidate to Ready so the full aggregate lane runs, complete required review, then STOP for explicit merge authorization. Do not reinterpret or advance B10/Product in this increment.** |
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
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B10 continues separately; PR #65 changes verification machinery only** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- Methodology adoption is integrated at `689dab34b0a756cbd7c790a6c5277d887ced0b4c`.
- PR #65 removes accidental CI ceremony: Draft PRs use quick/change-aware verification; Ready candidates and `main` use the full aggregate gate; the stable `required` check remains.
- Quick verification preserves repository/bootstrap/methodology/diff hygiene and invokes only directly affected lightweight verifiers; it does not claim whole-Product non-regression.
- Full verification remains the integration proof and retains Product historical/current, D7-R, D8-R, W1 and OAD source-reachability coverage.
- PR #64/B10 is an explicit independent parallel candidate. Nothing in #65 treats it as accepted, merged, or semantically changed.
- Product **106/31/H-A-S**, accepted OAD, runtime baseline NONE, existing LOCKs and implementation block remain unchanged.

```text
PR #65 CI proportionality → integrate only after full proof/review → continue PR #64 B10 independently → Pre-D9 readiness → D9 → implementation only after accepted D9
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
