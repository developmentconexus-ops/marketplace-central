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
| Current repository-governance increment | **PR #65 — proportional CI verification — R1 adjudicated; F1/F2 corrected; R2 confirmation pending** |
| Parallel Product candidate | **PR #64 — B10 Preparação — UNMERGED; not integrated authority on this branch** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Run full aggregate verification on the exact corrected PR #65 head under ruleset-required context `required`, then one isolated R2 confirmation of R1-F1/F2. If converged, STOP for explicit merge authorization. Do not advance B10/Product here.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; NOTIF-01 D0-R ACCEPTED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; D1-R2 PASS |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; D2-R6 ACCEPTED |
| D3 — Communication / Events | ACCEPTED / CLOSED; D3-R3 PASS |
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
- PR #65 removes accidental CI ceremony: Draft runs publish advisory context `quick`; Ready candidates and `main` publish ruleset-required context `required` and run the full aggregate gate.
- Quick verification preserves repository/bootstrap/methodology/diff hygiene and runs only lightweight verifiers explicitly mapped to touched surfaces. It is an early-signal under-approximation, **not** a completeness or whole-Product non-regression proof.
- Full verification remains the sole integration completeness boundary and retains Product historical/current, D7-R, D8-R, W1 and OAD source-reachability coverage.
- R1 found no MATERIAL defect. R1-F1/F2 were accepted and corrected with exact step-association falsifiers plus dynamic `quick`/`required` check naming; R1-F3 was clarified as intentional under-selection; R1-F4 is deferred until real map growth/failure; R1-F5 requires no change while `conventional-title` remains advisory.
- PR #64/B10 is an explicit independent parallel candidate. Nothing in #65 treats it as accepted, merged, or semantically changed.
- Product **106/31/H-A-S**, accepted OAD, runtime baseline NONE, existing LOCKs and implementation block remain unchanged.

```text
PR #65 corrected candidate → full required proof → R2 confirmation → merge only after explicit authorization → continue PR #64 B10 independently
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
