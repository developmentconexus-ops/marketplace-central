# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Integrated checkpoints | **PR #61 → `b54d17bfe6d794645d198a9160f4a2a1c63647e8`; PR #62 methodology → `689dab34b0a756cbd7c790a6c5277d887ced0b4c`; PR #65 proportional CI → `a2aeac19816c90ee30bf373cef0448d52a486c7e`; PR #68 PublicationRequirements wire → `ed3d164b0574b7950c2c7467d150c89576bba1ec`; PR #64 B10 → `bdbbef43ed3a5e9d912e67ddac5173024352eaa3`** |
| Method profile | **local [`engineering-method.md`](development/engineering-method.md) + [`frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md) v2.3** |
| Current prerequisite | **Human-Operable Read Projection & Wire Conformance — GLOBAL MAXIMUM DESIGN APPROVED / IMPLEMENTATION PLAN REVIEW REQUIRED** |
| Written design | **[`docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md`](superpowers/specs/2026-08-25-human-operable-read-projection-design.md) — operator-approved 2026-08-25** |
| Implementation plan | **[`docs/superpowers/plans/2026-08-25-human-operable-read-projection-implementation-plan.md`](superpowers/plans/2026-08-25-human-operable-read-projection-implementation-plan.md) — CANDIDATE / OPERATOR REVIEW REQUIRED** |
| B20 increment | **PR #69 — PAUSED / NO P8; cannot resume before prerequisite acceptance/integration and bounded B10 correspondence revalidation** |
| B10 status | **Global Maximum + main structure preserved; correspondence region REOPEN REQUIRED by upstream contract falsifier; P9 must rerun after repaired wire** |
| Upstream finding | **Canonical refs/keys were reused as human read projections; Readiness dynamic vocabulary/source/correspondence candidate projection is insufficient; directly implicated W2→OAD Offering read drift also exists** |
| Current decision | **`RESTRUCTURE NOW` APPROVED — targeted D4/D4-R1 + D5 W2/W3 + OAD read-projection/conformance repair; D0/D1/D2/D3/W1/W4 confirmed unless execution finds a new material falsifier** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only — expected unchanged by prerequisite** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; proportional targeted proof only when the claim needs it** |
| Exact next action | **Operator reviews the implementation plan. Do NOT modify D4/D5/OAD or B10/B20 HTML until the implementation plan is explicitly approved. After plan approval, execute the prerequisite task-by-task; B20 remains paused.** |
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
| D4 — External Integrations | ACCEPTED / CLOSED; **bounded presentation-evidence reopen approved by current prerequisite, execution not yet authorized** |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL; **bounded human-operability reopen approved, execution not yet authorized** |
| D5 — API | ACCEPTED / CLOSED; D5-R6 106/31 PROVED; D5-R7 W1 REPAIR ACCEPTED |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110 LOCKED; B10 main structure re-LOCKED in D6-R2** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; D7-R ACCEPTED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; D8-R ACCEPTED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — human-operable read-projection prerequisite at implementation-plan review; B20 PAUSED** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- PR #64 is integrated in `main`; B10's simplified Global Maximum remains accepted: marketplace requirements + source values + downstream ListingIntent authoring/provider validation.
- B20 planning exposed that Offering `MarketplaceListing` collection/detail reads do not provide sufficient human presentation for an operator, while Performance already carries a local non-authoritative Listing `display_name`.
- Global Coherence found the broader defect class: canonical refs/keys are correct identity/write carriers but are inconsistently reused as human read projections.
- B10 exposed stronger Readiness falsifiers: multi-source/exact-subject presentation is incomplete; dynamic requirement/option/unit/context and FOLLOW_SOURCE candidate keys lack sufficient human presentation; correspondence may ask a human to Resolve by `candidate_key` without returning a selectable human-recognizable candidate population.
- Canonical W2 already requires directly implicated MarketplaceListing/ListingIntent read axes that are absent or incomplete in the current OAD, including MarketplaceListing publication context/media/provenance and ListingIntent authored-media presentation separation.
- Performance's local Listing `display_name` and Governance's immutable `subject_display_label` prove that current presentation and purpose/historical snapshots are already distinct Product meanings in different places; the prerequisite consolidates the rule without creating a generic presentation authority.
- The operator approved the Global Maximum design: **Canonical Ref ≠ Current Read Projection ≠ Purpose/Historical Snapshot**, with owner-specific typed projections only for proven human recognition/selection/explanation jobs.
- A concrete implementation plan now exists and has been self-reviewed for scope, file paths, type-name consistency, TDD/negative proof, one-gate CI, D3 stop conditions, and independent review. It is **not yet execution authority**.
- Expected Product inventory remains **106 operations / 31 Permissions / H-A-S**; a count change is a new finding and must stop/reopen explicitly.
- PR #69 / B20 remains paused. No B20 P8 HTML is authorized against the deficient Listing read contract.

```text
PR #64 integrated
→ B20 planning
→ human-operability falsifier
→ Global Coherence Review
→ RESTRUCTURE NOW design
→ operator-approved written spec
→ implementation plan candidate
→ operator plan review
→ only after approval: bounded prerequisite implementation/proof
→ prerequisite integration
→ B10 correspondence-region revalidation + P9 + operator re-LOCK
→ resume B20
```

One coherent prerequisite lands before B20 resumes. Return to [`index.md`](index.md) for task routing.
