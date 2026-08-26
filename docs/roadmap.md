# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed current owners.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Main baseline | **PR #76 — B23 listing-intent authoring integrated at `6ce3902bbd8e7593249c7f4a45658c9d0027bb96`; required CI PASS** |
| Method profile | **local [`engineering-method.md`](development/engineering-method.md) + [`frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md) v2.3** |
| Current acceptance increment | **B24 — Preços / R24 — OPEN; P8 OPERATOR-LOCKED 2026-08-26 (revision 5); P9 run raised one blocking UPSTREAM FINDING (Market position projection); operator adjudication NEXT** |
| B24 accepted input | **R24 `/publicacoes/precos` — PriceIntent collection/detail/create with Market/Economics handoff; Offering owns price writes (`price.manage`, H/A, Idempotency-Key); supersede model, no in-place price edit** |
| Prior increment | **B23 — P8 LOCKED 2026-08-26 + P9 PASS / CLOSED (variations + category discovery repairs), integrated by PR #76** |
| Read contract law | **Owner-specific typed read projections (Canonical Ref ≠ Current Read Projection ≠ Purpose/Historical Snapshot) are the only admitted human presentation basis (PR #70; extended by the PR #69 `source_product_link` repair)** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **107 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only — SearchPublicationContexts admitted 2026-08-26** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; universal checks always run; Product proof is diff-aware and fails safe when reliable changed-surface detection is unavailable** |
| Exact next action | **Adjudicate the B24 P9 upstream finding — the Market owner carries no delivered range, rank or candidate-price position ([P9 contract §6](engineering/rebaseline/D6-R2-P9-B24-PRICE-INTENTS-SCREEN-CONTRACT.md)) — as RATIFIED / REJECTED / DEFERRED. No B24 P9 PASS, further blocks, Pre-D9/D9 or implementation before that adjudication.** |
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
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B00/B01/B00-R2/B11/B12/B110/B10/B20/B23/B24 P8 LOCKED; B24 P9 blocked on one upstream finding** |
| Pre-D9 readiness | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Current result

- PR #76 is integrated at `6ce3902bbd8e7593249c7f4a45658c9d0027bb96`: B23 ListingIntent authoring is **P8 OPERATOR-LOCKED (2026-08-26, revision 3) and P9 PASS / CLOSED**, carrying the variations model and the category discovery/equivalence repair; the Product surface is **107 operations / 31 Permissions / H-A-S** with historical non-regression preserved.
- Nine D6-R2 blocks are LOCKED; **B24 — Preços (R24)** is the next dependent increment: the price destination the B23 editor and B20 detail point at.
- The B24 design must preserve: Offering as sole price-write owner (`price.manage`, H/A, Idempotency-Key); typed targets (existing listing / pre-creation) with typed target presentation; honest state truth (pending/applied/rejected/ambiguous/superseded) and convergence; supersede-instead-of-edit; ambiguous never blindly retried; Market/Economics evidence as read-only owner-separated context; no automatic repricing engine.
- No bulk repricing, rule engine, screen-shaped API or speculative shared component authority is admitted by default.

```text
PR #76 B23 integrated
→ B24 acceptance increment opened
→ bounded P8 structural design adjudicated: operator approved (2026-08-26)
→ first R24 candidate rendered; operator alignment → listing-centric pricing workbench approved
→ workbench revision rendered + proof (13/13)
→ second alignment: owner-evaluated waterfall/indicators + single create home
→ revision 3 rendered + proof (20/20)
→ third alignment: live debounced evaluation, market range, collapsible waterfall
→ revision 4 rendered + proof (24/24)
→ fourth alignment: compact row, labeled margin, range + competitive rank live in the Mercado column
→ revision 5 rendered + proof (28/28)
→ operator walkthrough: P8 LOCKED (2026-08-26, revision 5)
→ P9 bidirectional Screen Contract run
→ UPSTREAM FINDING: Market carries no delivered range / rank / candidate position   ← CURRENT
→ operator adjudication: RATIFY / REJECT / DEFER
→ bounded repair + proof (if ratified)
→ P9 PASS
→ P10 pattern consolidation
→ required CI
→ operator-authorized integration
→ only then the next D6-R2 block
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
