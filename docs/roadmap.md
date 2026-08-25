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
| Current acceptance increment | **PR #64 / B10 — P8 OPERATOR-RATIFIED / LOCKED; A01 `ACCEPT_FOR_LOCK_WITH_LATER_PROBE`; P9 DERIVED / BLOCKED — P8 REOPEN REQUIRED** |
| B10 P8 evidence | **[P8 ratification](engineering/rebaseline/D6-R2-P8-B10-PREPARATION-RATIFICATION.md) — human-first candidate operated/LOCKED** |
| B10 P9 | **[Screen contract](engineering/rebaseline/D6-R2-P9-B10-PREPARATION-SCREEN-CONTRACT.md) — `F-P9-B10-01`: `known` source evidence does not equal requirement satisfied; UPSTREAM FINDING NONE** |
| B10 proof | **Draft quick CI #747 PASS at `2300172ef96f2d34d943deef760ebe5f6fcfbb57`; P9 falsifiers 5/5; P8 ratification falsifiers 3/3; B10 P8 falsifiers 12/12** |
| Resolved upstream finding | **PublicationRequirements wire gap RESOLVED by integrated PR #68** |
| LOCK impact | **B00 / B01 / B00-R2 / B11 / B12 / B110 UNAFFECTED** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **STOP. Operator adjudicates bounded B10 P8 wording reopen: `Atendido` → `Informação disponível`; `requisitos atendidos` → `com informação disponível`. Do not change the LOCKED screen or begin P10 before approval, corrected P8 proof/re-ratification and P9 closure.** |
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
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110 LOCKED; B10 LOCK challenged only by P9 F-P9-B10-01** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; D7-R ACCEPTED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; D8-R ACCEPTED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B10 P9 BLOCKED pending bounded P8 wording adjudication** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- PR #68 remains integrated; Product **106/31/H-A-S**, runtime NONE and existing LOCKs are unchanged.
- B10 P8 was operator-ratified with A01 accepted for later probe. P9 successfully bound route/state/read/write/failure authority and found no upstream Product gap.
- `F-P9-B10-01` proves one frontend overclaim: source evidence `known` says information exists, not that Readiness issued a per-requirement satisfied/met fact. The only proposed repair is human wording; no layout/flow/Product change.
- Exact current-head quick proof is GREEN; the branch remains Draft because the operator must adjudicate reopening the LOCKED wording before P9 can close.

```text
B10 P8 LOCKED → P9 F-P9-B10-01 → OPERATOR ADJUDICATION → bounded P8 correction/re-LOCK → P9 close → P10
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
