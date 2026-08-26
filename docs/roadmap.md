# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed current owners.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Main baseline | **PR #77 — B24 pricing workbench + MKT-01 Market projection integrated at `b53f34ad460d33f07e490ca84e0c054bb2689ad1`; required CI PASS** |
| Method profile | **local [`engineering-method.md`](development/engineering-method.md) + [`frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md) v2.3** |
| Current acceptance increment | **B30 — Disponibilidade / R30 — NOT OPENED; P8 structural design adjudication is the entry gate** |
| B30 accepted input | **R30 `/disponibilidade` — sellable availability current population + provenance/knowledge/config continuation (D6-R2 P5 inventory); Availability owns the meaning, Offering owns any listing write** |
| Prior increment | **B24 — P8 LOCKED 2026-08-26 (revision 5) + P9 PASS / CLOSED (MKT-01 Market position projection), integrated by PR #77** |
| Read contract law | **Owner-specific typed read projections (Canonical Ref ≠ Current Read Projection ≠ Purpose/Historical Snapshot) are the only admitted human presentation basis (PR #70; extended by the PR #69 `source_product_link` repair)** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **108 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only — EvaluateCompetitivePositionScenario admitted 2026-08-26 (MKT-01)** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; universal checks always run; Product proof is diff-aware and fails safe when reliable changed-surface detection is unavailable** |
| Exact next action | **Open the B30 — Disponibilidade acceptance increment: adjudicate its bounded P8 structural design with the operator before rendering any candidate. No B30 HTML, further blocks, Pre-D9/D9 or implementation before that adjudication.** |
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
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00/B01/B00-R2/B11/B12/B110/B10/B20/B23/B24 LOCKED / integrated; B30 next** |
| Pre-D9 readiness | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Current result

- PR #77 is integrated at `b53f34ad460d33f07e490ca84e0c054bb2689ad1`: B24 Preços is **P8 OPERATOR-LOCKED (2026-08-26, revision 5) and P9 PASS / CLOSED**. The locked shape is the compact pricing workbench — owner facts per row (current margin; delivered range and our rank in the observed comparable population), a debounced owner evaluation of the typed price, the cascade behind a row-level disclosure, and one explicit supersede-only write per row.
- **MKT-01** projected already-accepted Market meaning into D5 (D1 competitive-position ownership; D4-B4 §6.2 provider evidence lane): `MarketDeliveredPriceRange`, `MarketRank` with a closed `observed_comparable_population` basis, and `EvaluateCompetitivePositionScenario`. Surface **108 / 31 / H-A-S**; historical non-regression preserved.
- Ten D6-R2 blocks are LOCKED. The remaining D6-R2 surface is R30 Disponibilidade, R40+ Performance, R50 Mercado, R60+ Economia, R70+ Vendas and the rest of the P5 inventory.
- **B30 — Disponibilidade (R30)** is the next acceptance increment and is **not yet opened**: its bounded P8 structural design must be adjudicated with the operator first.

```text
PR #77 B24 integrated
→ B30 acceptance increment opened
→ bounded P8 structural design adjudicated with the operator   ← NEXT
→ candidate rendered + deterministic proof + browser operation
→ operator walkthrough: LOCK / REVISE / UPSTREAM FINDING
→ P9 bidirectional Screen Contract
→ P10 pattern consolidation
→ required CI
→ operator-authorized integration
→ only then the next D6-R2 block
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
