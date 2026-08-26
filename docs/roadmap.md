# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed current owners.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Main baseline | **PR #75 — B10 correspondence revalidation integrated at `ad06e70cb31c1037b5ffcebc116a57749e4728d4` over the PR #70 read-projection prerequisite; required CI PASS** |
| Method profile | **local [`engineering-method.md`](development/engineering-method.md) + [`frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md) v2.3** |
| Current acceptance increment | **B20 — Publicações core / R20–R21 — OPEN; P8 structural design OPERATOR-APPROVED 2026-08-26; browser-operable candidate rendered; operator walkthrough NEXT** |
| B20 accepted input | **R20 Marketplace Listing collection + exact Installation context; R21 one source-qualified Listing detail with owner-separated material regions; no screen-shaped Product capability** |
| B20 read contract | **Repaired by PR #70 and revalidated through PR #75: owner-specific typed read projections (Canonical Ref ≠ Current Read Projection ≠ Purpose/Historical Snapshot) are the only admitted human presentation basis** |
| Prior increment | **B10 — fully LOCKED (correspondence region RE-LOCKED 2026-08-26); P9 PASS / CLOSED; integrated through PR #64 + PR #75** |
| B10 Global Maximum | **requirements + source values + downstream authoring/provider validation; `source_sufficiency` REJECTED; NO NEW UPSTREAM WIRE FIELD** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; universal checks always run; Product proof is diff-aware and fails safe when reliable changed-surface detection is unavailable** |
| Exact next action | **Operate the browser-operable B20 R20/R21 candidate (`qualification/d6-r2-wireframes/b20-publications.html`) and obtain operator disposition `LOCK / REVISE / UPSTREAM FINDING`. Do not run B20 P9, begin B23, Pre-D9/D9 or Product implementation before explicit LOCK.** |
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
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B10 fully LOCKED / P9 CLOSED / integrated; B20 RESUMED — P8 design adjudication next** |
| Pre-D9 readiness | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Current result

- PR #75 is integrated at `ad06e70cb31c1037b5ffcebc116a57749e4728d4`: the operator re-LOCKED the B10 correspondence region on 2026-08-26 and the bounded P9 rerun passed against the integrated OAD. **B10 is fully LOCKED and P9 CLOSED.**
- PR #70 remains the read-projection authority: owner-specific typed projections for proven human recognition/selection/explanation jobs, with **Canonical Ref ≠ Current Read Projection ≠ Purpose/Historical Snapshot**.
- Product remains **106 operations / 31 ordinary Permissions / H-A-S**, runtime NONE; no Product/OAD change was introduced by the B10 closure.
- The B20 pause condition is fully discharged: the deficient Listing read contract that paused PR #69 was repaired by PR #70 and revalidated through PR #75. **B20 is RESUMED** on the reanchored PR #69.
- P5 places **B20 Publications core** immediately after B10 and explicitly marks P6/P7 **NOT TRIGGERED** because collection/detail is conventional and authority separation is already explicit.
- B20 covers **R20 `/publicacoes`** and **R21 source-qualified Listing detail**. The later ListingIntent editor remains B23/R22–R23 and is not pulled into B20.
- The B20 P8 design must preserve exact Marketplace Installation context, source-qualified Listing identity, honest Product knowledge states, owner separation in composed detail evidence, navigation-handoff ≠ mutation, and must consume only the repaired typed read projections — never canonical refs/keys as human presentation.
- No bulk-selection framework, saved-view platform, provider-direct write, screen-shaped API, universal normalized entity store or speculative shared component authority is admitted by default.

```text
PR #70 read-projection prerequisite integrated
→ PR #75 B10 correspondence re-LOCK + P9 rerun PASS integrated
→ B20 acceptance increment RESUMED on reanchored PR #69
→ bounded P8 structural design adjudicated: operator approved (2026-08-26)
→ browser-operable R20/R21 candidate rendered + proof
→ operator walkthrough: LOCK / REVISE / UPSTREAM FINDING   ← CURRENT
→ P9 bidirectional Screen Contract
→ P10 pattern consolidation
→ required CI
→ operator-authorized integration
→ only then dependent B23 increment
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
