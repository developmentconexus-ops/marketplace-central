# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed current owners.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Main baseline | **`bdbbef43ed3a5e9d912e67ddac5173024352eaa3` — PR #64 / B10 integrated** |
| Method profile | **local [`engineering-method.md`](development/engineering-method.md) + [`frontend-product-experience-planning-method.md`](development/frontend-product-experience-planning-method.md) v2.3** |
| Repository-health increment | **PR #71 — PROVED / INTEGRATION CANDIDATE; merge requires explicit operator authorization** |
| Health result | **current-authority consolidation COMPLETE; absorbed D6-R2/NOTIF/AuthorizationRequest intermediates retired; satisfied ADR/plan residue retired; temporary health working artifacts retired; Git history remains archive** |
| Independent review | **Codex full review + review-slice challenge completed; all material findings adjudicated and fixed (P12 assumptions rehomed; D5-R2 evidence link retargeted; no-base CI fail-safe/implementation block corrected)** |
| Health proof | **Product/proof-input diff → full 95/99/106 Product proof PASS; isolated docs-only PR #72 → `product_oad_proof: SKIPPED_NOT_AFFECTED` + gate PASS; Product OAD and operator-locked P8 HTML unchanged** |
| B20 increment | **PR #69 — PAUSED / NO P8** |
| Human-operable read-projection prerequisite | **PR #70 — PAUSED pending reanchor after health integration** |
| B10 status | **main structure/operator LOCK preserved; correspondence region remains the bounded upstream-repair revalidation after health integration** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; universal checks always run; heavy Product proof runs for Product/proof/toolchain inputs and fails safe to full proof when reliable diff detection is unavailable; implementation block also fails safe without a diff base** |
| Exact next action | **Operator reviews PR #71 and explicitly authorizes or declines merge. Do not merge implicitly. After integration, revalidate `main`, audit remote branches, reanchor PR #70, complete the bounded B10 correspondence/P9 revalidation required by that prerequisite, then resume PR #69 / B20.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; later Personal Notifications meaning consolidated into current owner by PR #71 |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; Personal Notifications supporting owner + AuthorizationRequest boundaries consolidated by PR #71 |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; Notification + AuthorizationRequest identity/lifecycle consolidated by PR #71 |
| D3 — Communication / Events | ACCEPTED / CLOSED; Notification + AuthorizationRequest propagation/recovery consolidated by PR #71 |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; current owners consolidated to 106/31/H-A-S; canonical OAD unchanged |
| D5-R2 — Operational Read Projection Repair | ACCEPTED / CANONICAL / current consumer evidence retargeted to P5 |
| D6 — Frontend | ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110/B10 operator-LOCK evidence preserved |
| D7 — Runtime / Jobs / Transactions | ACCEPTED / CLOSED; AuthorizationRequest runtime composition consolidated by PR #71 |
| D8 — Golden Flows | ACCEPTED / CLOSED / INTEGRATED; AuthorizationRequest revalidation consolidated by PR #71 |
| D8-R2 — GF-02 Operational Read Revalidation | ACCEPTED / PASS |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — repository-health prerequisite proved; waiting only for operator decision on PR #71 integration** |
| Pre-D9 readiness | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Current result

- PR #64 is integrated in `main`; B10 operator-LOCK evidence remains current.
- PR #69 remains paused before B20 P8 after exposing the human-operability/read-projection prerequisite.
- PR #70 contains the approved Global Maximum prerequisite plan and remains paused until this health increment lands and is reanchored.
- PR #71 corrected the structural repository defect: accepted later amendments now live in current D0–D8/W1–W4/P5/P8/P9 owners rather than requiring historical ratification/repair chains.
- The remaining P12 assumptions A01–A05 are explicitly current in P5 and may not silently disappear with historical closure material.
- Git history is the archive. The active tree retains current owners, current named evidence, and only ADR residues with a still-live retirement condition.
- Historical 95/99 Product non-regression expectations live in explicit technical fixtures under `scripts/fixtures/`, not stale semantic Markdown.
- `scripts/gate.ps1` preserves the full Product proof for Product/proof/toolchain changes, skips it for unrelated docs/frontend planning changes, and fails safe when a trustworthy base is unavailable.
- A disposable PR #72 proved the docs-only skip path and was closed without merge.
- Independent Codex reviews found material issues in the cleanup/proof boundary; every finding was fixed and its review thread resolved.
- Product semantics, canonical OAD bytes, 106/31/H-A-S, active runtime NONE and D9 implementation block remain unchanged.

```text
PR #64 integrated
→ B20 planning exposes human-operability gap
→ #69 paused
→ #70 approved prerequisite plan paused
→ #71 repository-health rebaseline
   → classify
   → canonical rehome
   → retire absorbed intermediates
   → isolate historical proof fixtures
   → proportional CI proof
   → independent review + adjudication
   → retire health working artifacts
   → PROVED / INTEGRATION CANDIDATE   ← CURRENT
→ explicit operator merge authorization only
→ reanchor #70
→ bounded B10 correspondence revalidation/P9 as required
→ resume #69 / B20
```

Return to [`index.md`](index.md) for current-owner routing.
