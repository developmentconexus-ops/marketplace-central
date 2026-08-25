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
| Current prerequisite | **PR #71 — Repository Governance & Context Health Rebaseline — PROOF COMPLETE / INDEPENDENT REVIEW ACTIVE** |
| Health result | **canonical rehome COMPLETE; absorbed D6-R2/NOTIF/AuthorizationRequest history retired; satisfied ADR/plan residue retired; Product proof decoupled from stale Markdown fixtures and made diff-aware** |
| Health proof | **Product/proof-input diff → full 95/99/106 Product proof PASS; isolated docs-only PR #72 → `product_oad_proof: SKIPPED_NOT_AFFECTED` + gate PASS; Product OAD and locked P8 HTML unchanged** |
| B20 increment | **PR #69 — PAUSED / NO P8** |
| Human-operable read-projection prerequisite | **PR #70 — PAUSED pending reanchor after health integration** |
| B10 status | **main structure/operator LOCK preserved; correspondence region remains the bounded upstream-repair revalidation after health integration** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Active runtime baseline | **NONE** |
| Aggregate CI | **one required check / one `npm run gate`; universal checks always run; heavy Product proof runs only for Product/proof/toolchain inputs or fail-safe when diff base is unavailable** |
| Exact next action | **Complete independent review of PR #71, adjudicate any valid Critical/Important findings, rerun required CI, then retire temporary health spec/plan/audit. Stop for explicit operator merge authorization. Do not render B20 or execute PR #70 Product/OAD changes first.** |
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
| D5-R2 — Operational Read Projection Repair | ACCEPTED / CANONICAL |
| D6 — Frontend | ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110/B10 operator-LOCK evidence preserved |
| D7 — Runtime / Jobs / Transactions | ACCEPTED / CLOSED; AuthorizationRequest runtime composition consolidated by PR #71 |
| D8 — Golden Flows | ACCEPTED / CLOSED / INTEGRATED; AuthorizationRequest revalidation consolidated by PR #71 |
| D8-R2 — GF-02 Operational Read Revalidation | ACCEPTED / PASS |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — health review/cleanup prerequisite before upstream repair and B20 continuation** |
| Pre-D9 readiness | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Current result

- PR #64 is integrated in `main`; B10 operator-LOCK evidence remains current.
- PR #69 remains paused before B20 P8 after exposing the human-operability/read-projection prerequisite.
- PR #70 contains the approved Global Maximum prerequisite design/implementation plan and remains paused until this health increment lands and is reanchored.
- PR #71 corrected the structural repository defect: accepted later amendments are now consolidated into their current D0–D8/W1–W4/P5/P8/P9 owners instead of requiring historical ratification/repair chains.
- Git history is the archive. The active tree retains current owners, current named evidence and only the ADR residues with an unresolved current retirement condition.
- Historical 95/99 Product non-regression expectations now live in explicit technical fixtures under `scripts/fixtures/`, not in stale current semantic Markdown.
- `scripts/gate.ps1` preserves the exact heavy Product proof for Product/proof/toolchain changes and skips it for unrelated docs/frontend planning changes while keeping universal repository checks.
- A disposable PR #72 proved the docs-only route and was closed without merge.
- Product semantics, canonical OAD bytes, 106/31/H-A-S, active runtime NONE and D9 implementation block are unchanged.

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
   → independent review        ← CURRENT
   → retire health working artifacts
→ explicit operator merge authorization
→ reanchor #70
→ bounded B10 correspondence revalidation/P9 as required
→ resume #69 / B20
```

Return to [`index.md`](index.md) for current-owner routing.
