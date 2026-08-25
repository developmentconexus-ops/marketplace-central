# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed current owners.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Main baseline | **`181f606ceaf5fadd7b25aab2008d0256ed6ad7de` — PR #71 repository-health rebaseline integrated; required CI PASS** |
| Method profile | **local [`engineering-method.md`](development/engineering-method.md) + [`frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md) v2.3** |
| Current prerequisite | **PR #70 — Human-Operable Read Projection & Wire Conformance — REANCHORED / GLOBAL MAXIMUM DESIGN APPROVED / IMPLEMENTATION PLAN REVIEW REQUIRED** |
| Written design | **[`docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md`](superpowers/specs/2026-08-25-human-operable-read-projection-design.md) — operator-approved 2026-08-25; reanchored on current main** |
| Implementation plan | **[`docs/superpowers/plans/2026-08-25-human-operable-read-projection-implementation-plan.md`](superpowers/plans/2026-08-25-human-operable-read-projection-implementation-plan.md) — FINAL REANCHORED CANDIDATE / OPERATOR APPROVAL REQUIRED** |
| B20 increment | **PR #69 — PAUSED / NO P8; cannot resume before prerequisite acceptance/integration and bounded B10 correspondence revalidation** |
| B10 status | **Global Maximum + main structure/operator LOCK preserved; correspondence region REOPEN REQUIRED by upstream contract falsifier; P9 must rerun after repaired wire** |
| Upstream finding | **Canonical refs/keys were reused as human read projections; Readiness dynamic vocabulary/source/correspondence candidate projection is insufficient; directly implicated W2→OAD Offering read drift also exists** |
| Current decision | **`RESTRUCTURE NOW` APPROVED — targeted D4/D4-R1 + D5 W2/W3 + OAD read-projection/conformance repair; D0/D1/D2/D3/W1/W4 confirmed unless execution finds a new material falsifier** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only — expected unchanged by prerequisite** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; universal checks always run; Product proof is diff-aware and fails safe when reliable changed-surface detection is unavailable** |
| Exact next action | **Operator reviews the reanchored implementation plan. Do NOT modify D4/D5/OAD or B10/B20 HTML until the plan is explicitly approved. After approval, execute Tasks 1–6 on PR #70; B20 remains paused.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; later Personal Notifications meaning consolidated into current owner by PR #71 |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; Personal Notifications + AuthorizationRequest boundaries consolidated by PR #71 |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; Notification + AuthorizationRequest identity/lifecycle consolidated by PR #71 |
| D3 — Communication / Events | ACCEPTED / CLOSED; bounded prerequisite execution must STOP if it proves a genuinely new semantic dependency |
| D4 — External Integrations | ACCEPTED / CLOSED; **bounded presentation-evidence reopen approved by current prerequisite, execution not yet authorized** |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL; **bounded human-operability reopen approved, execution not yet authorized** |
| D5 — API | ACCEPTED / CLOSED; current owners consolidated to 106/31/H-A-S; canonical OAD unchanged by reanchor |
| D5-R2 — Operational Read Projection Repair | ACCEPTED / CANONICAL / current consumer evidence retargeted to P5 |
| D6 — Frontend | ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110/B10 operator-LOCK evidence preserved |
| D7 — Runtime / Jobs / Transactions | ACCEPTED / CLOSED; runtime baseline remains NONE |
| D8 — Golden Flows | ACCEPTED / CLOSED / INTEGRATED |
| D8-R2 — GF-02 Operational Read Revalidation | ACCEPTED / PASS |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — PR #70 implementation-plan review; PR #69 / B20 PAUSED** |
| Pre-D9 readiness | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Current result

- PR #71 is integrated at `181f606ceaf5fadd7b25aab2008d0256ed6ad7de`; current authority owners, proportional CI and Git-as-history lifecycle are active on `main`.
- PR #70 was reanchored without rewriting history. Its design/plan reference no path retired by #71, and all required D4/D5/P5/OAD owners remain current.
- Reanchor analysis added the new focused verifier to the diff-aware Product-proof trigger set required by #71; a verifier-only change must not skip the proof it protects.
- PR #64 is integrated; B10's simplified Global Maximum remains accepted: marketplace requirements + source values + downstream ListingIntent authoring/provider validation.
- B20 planning exposed that Offering `MarketplaceListing` collection/detail reads do not provide sufficient human presentation for an operator, while Performance already carries a local non-authoritative Listing label.
- The operator approved the Global Maximum design: **Canonical Ref ≠ Current Read Projection ≠ Purpose/Historical Snapshot**, with owner-specific typed projections only for proven human recognition/selection/explanation jobs.
- The implementation plan is finalized and reanchored, but is **not execution authority** until the operator explicitly approves it.
- Expected Product inventory remains **106 operations / 31 Permissions / H-A-S**; a count change is a new finding and must stop/reopen explicitly.
- PR #69 / B20 remains paused. No B20 P8 HTML is authorized against the deficient Listing read contract.

```text
PR #64 integrated
→ B20 planning exposes human-operability gap
→ #69 paused / no P8
→ Global Maximum design approved
→ #71 repository-health rebaseline integrated
→ #70 reanchored on current main
→ operator implementation-plan review   ← CURRENT
→ only after approval: Tasks 1–6 prerequisite implementation/proof
→ prerequisite integration
→ B10 correspondence-region revalidation + P9 + operator re-LOCK
→ resume B20
```

One coherent prerequisite lands before B20 resumes. Return to [`index.md`](index.md) for current-owner routing.
