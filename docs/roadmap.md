# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed current owners.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Main baseline | **PR #70 — human-operable read-projection prerequisite integrated over `181f606ceaf5fadd7b25aab2008d0256ed6ad7de`; required CI PASS** |
| Method profile | **local [`engineering-method.md`](development/engineering-method.md) + [`frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md) v2.3** |
| Current prerequisite | **PR #70 — Human-Operable Read Projection & Wire Conformance — ACCEPTED / INTEGRATED / PROOF PASS** |
| Current increment | **B10 correspondence-region revalidation — functional P8 CANDIDATE / OPERATOR WALKTHROUGH REQUIRED; main B10 structure LOCK preserved** |
| Written design | **[`docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md`](superpowers/specs/2026-08-25-human-operable-read-projection-design.md) — DESIGN ACCEPTED / IMPLEMENTATION INTEGRATED through PR #70** |
| Implementation plan | **[`docs/superpowers/plans/2026-08-25-human-operable-read-projection-implementation-plan.md`](superpowers/plans/2026-08-25-human-operable-read-projection-implementation-plan.md) — FINAL REANCHORED / OPERATOR-APPROVED 2026-08-25** |
| B20 increment | **PR #69 — PAUSED / NO P8; cannot resume before prerequisite acceptance/integration and bounded B10 correspondence revalidation** |
| B10 status | **Global Maximum + main structure/operator LOCK preserved; correspondence region functional candidate active; P9 rerun awaits operator re-LOCK** |
| Resolved upstream finding | **PR #70 repaired Readiness dynamic vocabulary/source/correspondence candidate presentation and directly implicated W2→OAD Offering read drift** |
| Current decision | **`CURRENT STRUCTURE CONFIRMED` for B10 — consume the integrated PR #70 projection through a bounded correspondence-region repair; no broader B10 or Product reopen** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only — expected unchanged by prerequisite** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; universal checks always run; Product proof is diff-aware and fails safe when reliable changed-surface detection is unavailable** |
| Exact next action | **Operate the bounded B10 correspondence candidate and obtain operator disposition `LOCK / REVISE / UPSTREAM FINDING`; do not rerun P9 or resume B20 before explicit re-LOCK.** |
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
| D4 — External Integrations | ACCEPTED / CLOSED; **bounded presentation-evidence repair integrated through PR #70** |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL; **bounded human-operability repair integrated through PR #70** |
| D5 — API | ACCEPTED / CLOSED; **PR #70 OAD repair integrated at unchanged 106/31/H-A-S** |
| D5-R2 — Operational Read Projection Repair | ACCEPTED / CANONICAL / current consumer evidence retargeted to P5 |
| D6 — Frontend | ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110/B10 operator-LOCK evidence preserved |
| D7 — Runtime / Jobs / Transactions | ACCEPTED / CLOSED; runtime baseline remains NONE |
| D8 — Golden Flows | ACCEPTED / CLOSED / INTEGRATED |
| D8-R2 — GF-02 Operational Read Revalidation | ACCEPTED / PASS |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B10 correspondence P8 candidate awaits operator walkthrough; PR #69 / B20 PAUSED** |
| Pre-D9 readiness | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Current result

- PR #70 is integrated over the PR #71 baseline; repaired D4/D4-R1/W2/W3/OAD authority, proportional CI and Git-as-history lifecycle are active on `main`.
- PR #70 was reanchored without rewriting history before execution. Its design/plan reference no path retired by #71, and all required D4/D5/P5/OAD owners remain current.
- Reanchor analysis added the new focused verifier to the diff-aware Product-proof trigger set required by #71; a verifier-only change must not skip the proof it protects.
- PR #64 is integrated; B10's simplified Global Maximum remains accepted: marketplace requirements + source values + downstream ListingIntent authoring/provider validation.
- B20 planning exposed that Offering `MarketplaceListing` collection/detail reads do not provide sufficient human presentation for an operator, while Performance already carries a local non-authoritative Listing label.
- The operator approved the Global Maximum design: **Canonical Ref ≠ Current Read Projection ≠ Purpose/Historical Snapshot**, with owner-specific typed projections only for proven human recognition/selection/explanation jobs.
- The implementation plan is finalized, reanchored and **operator-approved**; Tasks 1–6 and the human-operable read-projection repair are integrated through PR #70.
- The focused proof passes **12/12 negative controls**, the publication source guard passes **13/13**, `npm run gate` passes, and required CI is green.
- Independent full review was adjudicated: two correctness findings were accepted and repaired; one stale authorization/status finding was rejected against current repository authority. All review threads are resolved.
- Product inventory remains **106 operations / 31 Permissions / H-A-S** with zero new Product paths, operations, Permissions, Principal kinds, generic presentation services or metadata bags.
- The fresh B10 increment preserves the locked main structure and renders only the repaired correspondence candidate population: selectable known candidates plus known-empty, unknown and unavailable states.
- The bounded candidate submits only `candidate_key`, requires explicit selection, blocks after consequential effect until reread, and passes **8/8 deterministic negative controls** plus desktop/mobile browser operation.
- No B10 re-LOCK is claimed and P9 has not been rerun; both remain operator-gated.
- PR #69 / B20 remains paused. No B20 P8 HTML is authorized before the bounded B10 re-LOCK and P9 rerun.

```text
PR #64 integrated
→ B20 planning exposes human-operability gap
→ #69 paused / no P8
→ Global Maximum design approved
→ #71 repository-health rebaseline integrated
→ #70 reanchored on current main
→ operator implementation-plan approval
→ Tasks 1–6 prerequisite implementation/proof complete
→ explicit operator authorization to integrate PR #70
→ prerequisite integration complete
→ bounded B10 correspondence-region functional candidate + proof
→ operator walkthrough: LOCK / REVISE / UPSTREAM FINDING   ← CURRENT
→ after re-LOCK: bounded P9 rerun
→ resume B20
```

The bounded B10 re-LOCK and P9 rerun must land before B20 resumes. Return to [`index.md`](index.md) for current-owner routing.
