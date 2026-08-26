# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed current owners.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Main baseline | **PR #69 — B20 publications core integrated at `675ca74e09e46ff8abe5b3d77635e25d314913ea` over the PR #75/#70 baseline; required CI PASS** |
| Method profile | **local [`engineering-method.md`](development/engineering-method.md) + [`frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md) v2.3** |
| Current acceptance increment | **B23 — ListingIntent authoring / R22–R23 — P8 OPERATOR-LOCKED 2026-08-26 (revision 3, after variations + category upstream repairs); P9 PASS / CLOSED; awaiting operator-authorized integration** |
| Resolved upstream finding | **Listing variations RATIFIED 2026-08-26 and REPAIRED: census scope + variation_axes, ListingIntentDesired.variations (coordinate-keyed, SKU-level source refs), observed_variations; 106/31/H-A-S preserved; projection proof 16/16** |
| B23 accepted input | **R22 `/publicacoes/intencoes` ListingIntent collection/work context; R23 revision-aware ListingIntent authoring/editor/media/evidence/submit-discard outcome; Offering owns all ListingIntent writes** |
| Read contract law | **Owner-specific typed read projections (Canonical Ref ≠ Current Read Projection ≠ Purpose/Historical Snapshot) are the only admitted human presentation basis (PR #70; extended by the PR #69 `source_product_link` repair)** |
| Prior increments | **B10 fully LOCKED + P9 CLOSED (PR #64/#75); B20 P8 LOCKED 2026-08-26 + P9 PASS / CLOSED, integrated by PR #69** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **107 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only — SearchPublicationContexts admitted 2026-08-26** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; universal checks always run; Product proof is diff-aware and fails safe when reliable changed-surface detection is unavailable** |
| Exact next action | **Obtain explicit operator authorization to integrate the B23 increment (variations + SearchPublicationContexts repairs, LOCKED P8, P9 PASS) into `main` via a pull request. The next D6-R2 block opens only after integration.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; later Personal Notifications meaning consolidated into current owner by PR #71 |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; Personal Notifications + AuthorizationRequest boundaries consolidated by PR #71 |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; Notification + AuthorizationRequest identity/lifecycle consolidated by PR #71 |
| D3 — Communication / Events | ACCEPTED / CLOSED; bounded increment execution must STOP if it proves a genuinely new semantic dependency |
| D4 — External Integrations | ACCEPTED / CLOSED; **bounded presentation-evidence repair integrated through PR #70** |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL; **bounded human-operability repair integrated through PR #70** |
| D5 — API | ACCEPTED / CLOSED; **PR #70 OAD repair integrated at unchanged 106/31/H-A-S** |
| D5-R2 — Operational Read Projection Repair | ACCEPTED / CANONICAL / current consumer evidence retargeted to P5 |
| D6 — Frontend | ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110/B10 operator-LOCK evidence preserved |
| D7 — Runtime / Jobs / Transactions | ACCEPTED / CLOSED; runtime baseline remains NONE |
| D8 — Golden Flows | ACCEPTED / CLOSED / INTEGRATED |
| D8-R2 — GF-02 Operational Read Revalidation | ACCEPTED / PASS |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00/B01/B00-R2/B11/B12/B110/B10/B20/B23 LOCKED; B23 increment awaits operator-authorized integration** |
| Pre-D9 readiness | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Current result

- PR #69 is integrated at `675ca74e09e46ff8abe5b3d77635e25d314913ea`: B20 Publicações core is **P8 OPERATOR-LOCKED (2026-08-26, revision 3) and P9 PASS / CLOSED**, including the bounded `source_product_link` OAD repair at unchanged **106/31/H-A-S** and the deterministic B20 wireframe verifier (11/11) in the gate.
- PR #70/#75 remain the read-projection authority: owner-specific typed projections for proven human recognition/selection/explanation jobs; canonical refs never render as human labels.
- **B23 — ListingIntent authoring (R22–R23)** is the next dependent increment: the `Ir para intenções de anúncio` boundary locked in B20 lands here.
- The B23 design must preserve: Offering as sole ListingIntent write owner (`listing.manage`, H/A, Idempotency-Key on creates); draft/submitted/discarded lifecycle truth; `follow_source` vs `explicit_override` requirement resolutions keyed by canonical keys only; honest dispatchability (dispatchable/blocked/unknown/unavailable) with explicit blockers; external-effect honesty (not_attempted/pending/accepted/rejected/ambiguous, no blind retry after ambiguity); revision-aware authoring against `requirements_revision`; no provider-direct write and no live marketplace effect without explicit operator authorization.
- No generic form engine, provider field bag, screen-shaped API or speculative shared component authority is admitted by default.

```text
PR #69 B20 integrated
→ B23 acceptance increment opened
→ bounded P8 structural design adjudicated: operator approved (2026-08-26)
→ first browser-operable R22/R23 candidate rendered + proof
→ operator walkthrough exposed authoring simplicity → P6 reference study run
→ variations upstream finding RATIFIED; B23 P8 PAUSED
→ variations Global Maximum design adjudicated + repair applied with proof
→ B23 P8 RESUMED with full P6 dispositions + Variações region
→ operator raised category/taxonomy question; verification + evaluation increment opened
→ category discovery/equivalence adjudicated: C1+C2 RATIFIED, C3 DEFERRED; repair integrated (107/31)
→ B23 context region revised with real discovery/suggestion
→ operator walkthrough → LOCK (2026-08-26, revision 3)
→ P9 bidirectional Screen Contract PASS + P10 no new shared pattern
→ required CI
→ operator-authorized integration   ← CURRENT
→ only then the next D6-R2 block
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
