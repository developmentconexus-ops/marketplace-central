# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Integrated checkpoints | **PR #61 → `b54d17bfe6d794645d198a9160f4a2a1c63647e8`; PR #62 methodology → `689dab34b0a756cbd7c790a6c5277d887ced0b4c`; PR #65 proportional CI → `a2aeac19816c90ee30bf373cef0448d52a486c7e`; PR #68 PublicationRequirements wire → `ed3d164b0574b7950c2c7467d150c89576bba1ec`; PR #64 B10 Preparation → `bdbbef43ed3a5e9d912e67ddac5173024352eaa3`** |
| Method profile | **local [`engineering-method.md`](development/engineering-method.md) + [`frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md) v2.3** |
| Current acceptance increment | **B20 — Publicações core / R20–R21 — OPEN; P6/P7 NOT TRIGGERED; P8 structural design approval NEXT** |
| B20 accepted input | **R20 Marketplace Listing collection + exact Installation context; R21 one source-qualified Listing detail with owner-separated material regions; no screen-shaped Product capability** |
| Prior increment | **B10 — OPERATOR-RATIFIED / LOCKED; P9 PASS / CLOSED; P10 COMPLETE / NO NEW SHARED PATTERN; integrated by PR #64** |
| B10 Global Maximum | **requirements + source values + downstream authoring/provider validation; `source_sufficiency` REJECTED; NO NEW UPSTREAM WIRE FIELD** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; main run #771 PASS at `bdbbef43ed3a5e9d912e67ddac5173024352eaa3`** |
| Exact next action | **Adjudicate the bounded B20 P8 structural design in chat. After explicit approval, render the browser-operable R20/R21 candidate and run the operator walkthrough. Do not begin B23, Pre-D9/D9 or Product implementation before the B20 acceptance increment closes and integrates.** |
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
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B10 integrated; B20 OPEN / P8 design adjudication next** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- PR #64 is integrated into `main` at `bdbbef43ed3a5e9d912e67ddac5173024352eaa3`; its post-merge `required` run #771 passed.
- B10 is now accepted input, not active work: operator-LOCKED P8, P9 PASS, P10 no-new-shared-pattern.
- Product remains **106 operations / 31 ordinary Permissions / H-A-S**, runtime NONE; no Product/OAD change was introduced by B10.
- P5 places **B20 Publications core** immediately after B10 and explicitly marks P6/P7 **NOT TRIGGERED** because collection/detail is conventional and authority separation is already explicit.
- B20 covers **R20 `/publicacoes`** and **R21 source-qualified Listing detail**. The later ListingIntent editor remains B23/R22–R23 and is not pulled into B20.
- The B20 P8 design must preserve exact Marketplace Installation context, source-qualified Listing identity, honest Product knowledge states, owner separation in composed detail evidence, and navigation-handoff ≠ mutation.
- No bulk-selection framework, saved-view platform, provider-direct write, screen-shaped API, universal normalized entity store or speculative shared component authority is admitted by default.

```text
PR #64 B10 integrated
→ B20 acceptance increment opened
→ bounded P8 structural design adjudication
→ render browser-operable R20/R21 candidate only after approval
→ operator LOCK / REVISE / UPSTREAM FINDING
→ P9 bidirectional Screen Contract
→ P10 pattern consolidation
→ required CI
→ operator-authorized integration
→ only then dependent B23 increment
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
