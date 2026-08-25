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
| Method profile | **local [`engineering-method.md`](development/engineering-method.md) + [`frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md) v2.3** |
| Current acceptance increment | **PR #64 / B10 — P8 REOPENED / CANDIDATE; A01 `ACCEPT_FOR_LOCK_WITH_LATER_PROBE`; P9 PAUSED** |
| B10 Global Maximum | **requirements + source values + downstream authoring/provider validation; `source_sufficiency` REJECTED; NO NEW UPSTREAM WIRE FIELD** |
| B10 evidence | **[P6 / Global Maximum revalidation](engineering/rebaseline/D6-R2-P6-B10-PREPARATION-REFERENCE-STUDY.md) + [P8 reopen record](engineering/rebaseline/D6-R2-P8-B10-PREPARATION-RATIFICATION.md) + [paused P9](engineering/rebaseline/D6-R2-P9-B10-PREPARATION-SCREEN-CONTRACT.md)** |
| Resolved upstream finding | **PublicationRequirements wire gap remains RESOLVED by integrated PR #68** |
| LOCK impact | **B00 / B01 / B00-R2 / B11 / B12 / B110 UNAFFECTED; prior B10 LOCK preserved as historical evidence but current B10 is reopened** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; run #758 PASS at `2bb957dc447f432a16f3bb58382c5183810c22cd`; Product OAD 106/106 + auth proof PASS** |
| Exact next action | **Run a fresh operator walkthrough of the simplified B10 P8 candidate. Only the operator may LOCK / REVISE / raise an UPSTREAM FINDING. Do not rerun/close P9, begin P10/P11, Pre-D9/D9 or Product implementation before a fresh B10 LOCK.** |
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
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110 LOCKED; B10 boundedly REOPENED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; D7-R ACCEPTED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; D8-R ACCEPTED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B10 P8 REOPENED / CANDIDATE; aggregate proof GREEN; operator walkthrough NEXT** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- PR #68 remains integrated; Product **106/31/H-A-S**, runtime NONE and unrelated LOCKs are unchanged.
- P9 finding `F-P9-B10-01` triggered a wider operator-requested Global Maximum evaluation rather than a cosmetic wording patch.
- Fresh AnyMarket, Hub2b, Magis5 and Mercado Livre evidence supports the simpler hub model: marketplace requirements + available source values + downstream listing authoring/provider validation.
- A generic `source_sufficiency` layer is rejected as accidental complexity. No Product operation, Permission, Principal kind or new Product wire field is required.
- The prior B10 P8 LOCK is preserved as historical evidence, but the operator authorized a bounded reopen because the visible requirement model changes materially. A01 remains accepted and unaffected.
- CI was simplified to one durable aggregate `required` check. Planning prose/P8/P9/ratification string verifiers and CI-self-testing were retired; the canonical Product OAD proof and structural repository rails remain protected. Run #758 passed with current Product **106/106** and authentication proof green.

```text
P9 F-P9-B10-01
→ GLOBAL MAXIMUM REVALIDATED
→ simplified B10 P8 candidate
→ required aggregate proof PASS
→ operator walkthrough
→ LOCK / REVISE / UPSTREAM FINDING
→ P9 rerun only after LOCK
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
