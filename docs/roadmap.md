# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Integrated checkpoints | **PR #61 → `b54d17bfe6d794645d198a9160f4a2a1c63647e8`; PR #62 methodology → `689dab34b0a756cbd7c790a6c5277d887ced0b4c`; PR #65 proportional CI → `a2aeac19816c90ee30bf373cef0448d52a486c7e`** |
| Method profile | **`developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535`** |
| Current prerequisite increment | **PR #68 — PublicationRequirements wire truth repair — REVIEW-CONVERGED; final bootstrap guard revalidation pending** |
| Independent review | **R1 findings absorbed; R2 confirmed the corrected candidate with no new major issue** |
| B10 Product candidate | **PR #64 — Preparação — PAUSED at P8 / NOT LOCKED pending prerequisite integration** |
| B10 finding | **UPSTREAM FINDING — accepted publication-requirement semantics exceed the current Product OAD wire realization** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Restore bootstrap `<= 20 KiB` and obtain fresh `required`/full PASS on PR #68; then STOP for explicit squash-merge authorization. Do not resume PR #64/B10 P8, P9, Pre-D9/D9, or Product implementation first.** |
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
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110 LOCKED; B110 REVALIDATED / UNAFFECTED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; D7-R ACCEPTED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; D8-R ACCEPTED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B10 P8 paused until PR #68 is integrated** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- PR #65 is integrated. Draft PRs use advisory `quick`; Ready candidates and `main` use ruleset-required `required` with the full aggregate gate.
- B10 P8 exposed an upstream wire gap. PR #68 preserves provider/context-qualified publication requirements, bounded value constraints, and source knowledge without adding Product/PIM master, new operation, Permission, Principal kind, runtime, or ListingIntent editor.
- R1 findings were corrected and R2 confirmed them. Full proof at technical candidate `36e7c13f6712a9bbb88078931918f09f3ce69e64` passed Product **106/31/H-A-S**, OAD/source-value and negative-path controls. The later roadmap closeout made only the bootstrap size guard red (`20923 > 20480`).
- PR #64 remains paused at B10 P8; existing LOCKs and Pre-D9/D9/implementation blocks are unchanged.

```text
PR #68 bootstrap compact/revalidate
→ required/full GREEN
→ explicit operator squash-merge authorization
→ revalidate main + rebase/revalidate PR #64 B10 P8
→ operator adjudication / LOCK before P9
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
