# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Integrated checkpoints | **PR #61 → `b54d17bfe6d794645d198a9160f4a2a1c63647e8`; PR #62 methodology → `689dab34b0a756cbd7c790a6c5277d887ced0b4c`; PR #65 proportional CI → `a2aeac19816c90ee30bf373cef0448d52a486c7e`; PR #68 PublicationRequirements wire → `ed3d164b0574b7950c2c7467d150c89576bba1ec`** |
| Method profile | **`developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535` → `METHOD + FRONTEND-METHOD`** |
| Current acceptance increment | **PR #64 / B10 — Preparação — B10 P8 OPERATOR-RATIFIED / LOCKED; A01 `ACCEPT_FOR_LOCK_WITH_LATER_PROBE`; P9 NEXT** |
| B10 P8 evidence | **Human-first candidate operated and approved; pre-lock Draft quick proof PASS; ratification now carries operator LOCK while the low-fi remains candidate evidence** |
| Resolved upstream finding | **PublicationRequirements wire gap RESOLVED by integrated PR #68** |
| LOCK impact | **B00 / B01 / B00-R2 / B11 / B12 / B110 UNAFFECTED** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Run P9 — B10 Screen Contract + bidirectional backend trace against the LOCKED human-first P8. Do not redesign B10, begin P10/P11, Pre-D9/D9, or Product implementation before the ordered method gates.** |
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
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110/B10 LOCKED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; D7-R ACCEPTED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; D8-R ACCEPTED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B10 P8 LOCKED; P9 NEXT** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- PR #68 is integrated and post-merge `main` full CI passed.
- B10 keeps the accepted search/list → exact-subject detail structure and complete provider/context-specific requirement truth.
- Operator REVISE removed backend-shaped UX; the locked experience leads with **Resumo da preparação** and **Requisito do marketplace / Exigência / Situação / Informação atual / O que fazer**, with technical identities/keys behind secondary disclosure.
- The operator operated and approved the human-first candidate, then explicitly accepted A01 as `ACCEPT_FOR_LOCK_WITH_LATER_PROBE`; B10 P8 is therefore **OPERATOR-RATIFIED / LOCKED**.
- PR #68 wire distinctions remain encoded/testable: class vs applicability, six source-evidence states, seven value-spec families, candidate identity and source-media separation. Existing LOCKED blocks remain unaffected; Product **106/31/H-A-S**, runtime NONE and Pre-D9/D9/implementation blocks are unchanged.

```text
PR #68 INTEGRATED / MAIN GREEN
→ B10 HUMAN-FIRST P8 OPERATOR-RATIFIED / LOCKED
→ P9 SCREEN CONTRACT + BIDIRECTIONAL TRACE
→ P10 only after P9 closure
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
