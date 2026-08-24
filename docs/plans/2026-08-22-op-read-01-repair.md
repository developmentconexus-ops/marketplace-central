# OP-READ-01 Repair Implementation Plan

> **For agentic workers:** execute this plan task-by-task on `stage/d6-r2-frontend-realization`; Product implementation remains blocked until accepted D9.

**Goal:** close OP-READ-01 by enriching only existing owner-local operational collections so the frontend can perform honest operational triage without N+1 detail fan-out, synthetic workflow state, or a screen-shaped API.

**Architecture:** preserve the 99-operation / 30-Permission / H-A-S Product surface. Add only owner-native fields already present on point resources to the relevant ListItem schemas, plus typed owner-local filters required by the now-evidenced operational consumer. No new semantic owner, Product operation, Permission, cross-owner lifecycle, total count, priority score, or `/operational-dashboard` endpoint.

**Tech Stack:** OpenAPI 3.1.2 YAML, Redocly bundle/lint, openapi-typescript, oapi-codegen, Node proof scripts, repository PowerShell gate.

**Spec:** `docs/engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md`

## Global Constraints

- Product surface remains exactly **99 Product operations / 30 ordinary Permissions / Principal kinds H / A / S only**.
- Preserve existing semantic owners: MarketplaceSales, BusinessSystemMaterialization, Fulfillment, PostSaleResolution, OperationalWork.
- No Product implementation or runtime code.
- No `OperationalWorkflow`, `operational_stage`, `next_action`, `priority`, `total_count`, generic filter DSL, or screen-shaped aggregate endpoint.
- ListItem additions must be semantic subsets of the same owner point resource.
- Filters must be typed owner-native predicates with AND composition under W3.
- After wire changes, revalidate the affected GF-02 read/composition properties before frontend progression.

---

### Task 1: Add executable OP-READ-01 contract proof

**Files:**
- Create: `scripts/verify-operational-read-contract.mjs`
- Modify: `package.json`
- Modify: `scripts/gate.ps1`

**Produces:** a bundled-OAD proof that asserts the exact approved list-item fields/filters and rejects synthetic operational workflow fields.

- [ ] Add a proof script that bundles `contracts/api/product/openapi.yaml` and asserts:
  - `BusinessOrderIntentListItem` requires/exposes `convergence`.
  - `ListBusinessOrderIntents` admits only the new operational filters `external_effect_state?` and `convergence?` in addition to its existing query grammar.
  - `InvoicingIntent` and `InvoicingIntentListItem` expose source-qualified `sale`; the ListItem also exposes `fulfillment_execution_id?` and required `convergence`.
  - `ListInvoicingIntents` admits `external_effect_state?` and `convergence?` in addition to existing filters.
  - `FulfillmentExecutionListItem` exposes/requires the four checkpoint objects and exposes optional `provider_dispatch_deadline`.
  - `ListFulfillmentExecutions` admits `physical_readiness?`, four checkpoint-state filters (`pending|recorded`), and `provider_dispatch_deadline_before?`.
  - `ShipmentListItem` exposes optional `sale` and `dispatch_deadline`; `ListShipments` admits `state?`.
  - forbidden Product properties/parameters remain absent: `operational_stage`, `next_action`, `priority`, `urgency_score`, `total_count`, `kanban_column`.
  - Product operations remain 99 and ordinary Permissions remain 30.
- [ ] Add negative controls that mutate one required operational field away and inject one forbidden synthetic field; both must fail validation.
- [ ] Add the script to both `npm run gate` and `npm run gate:full` after the existing Product/performance proofs.
- [ ] Add the proof script to the gate required-file list.

### Task 2: Enrich W2 owner-local list projections

**Files:**
- Modify: `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md`
- Modify: `contracts/api/product/components.yaml`

**Produces:** the minimal owner-native list projection required for operational triage.

- [ ] Business-System Materialization:
  - `BusinessOrderIntentListItem`: add required `convergence` with the same enum/meaning as `BusinessOrderIntent.convergence`.
  - `InvoicingIntent`: add required source-qualified `sale` reference.
  - `InvoicingIntentListItem`: add required `sale` and `convergence`; add optional `fulfillment_execution_id`.
- [ ] Fulfillment:
  - `FulfillmentExecutionListItem`: add required `separation`, `physical_conference`, `packing`, `dispatch_handoff`; add optional `provider_dispatch_deadline`.
  - Reuse the existing `FulfillmentCheckpoint` and `Instant` schemas; do not create `stage` or `next_action`.
- [ ] Shipment:
  - `ShipmentListItem`: add optional source-qualified `sale` and optional `dispatch_deadline`, reusing the point resource meaning.
- [ ] Update W2 prose to state these fields are owner-native triage projections, not a new operational lifecycle.

### Task 3: Enrich W3 owner-local query grammar

**Files:**
- Modify: `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md`
- Modify: `contracts/api/product/paths-economics-governance-sales-materialization.yaml`
- Modify: `contracts/api/product/paths-fulfillment-postsale-work.yaml`

**Produces:** typed server-side narrowing for real operational queues without creating new list operations.

- [ ] `ListBusinessOrderIntents`: add optional `external_effect_state` and `convergence` filters using the existing owner enums.
- [ ] `ListInvoicingIntents`: add optional `external_effect_state` and `convergence` filters.
- [ ] `ListFulfillmentExecutions`: add optional `physical_readiness`, `separation_state`, `physical_conference_state`, `packing_state`, `dispatch_handoff_state`, `provider_dispatch_deadline_before`.
- [ ] Each checkpoint-state filter is exactly `pending|recorded`; no result/severity/priority filter.
- [ ] `ListShipments`: add optional `state` using the existing Shipment state enum.
- [ ] Preserve existing deterministic ordering, pagination/cursor law, source qualification, and AND filter composition.
- [ ] Update W3 matrix/prose; operation count remains unchanged.

### Task 4: Close OP-READ-01 authority and revalidate GF-02

**Files:**
- Modify: `docs/engineering/rebaseline/D6-R2-P4-R1-GLOBAL-IA-OPERATIONAL-MASS-REOPEN.md`
- Modify: `docs/engineering/rebaseline/D8-GOLDEN-FLOWS.md`
- Modify: `docs/roadmap.md`
- Modify: `docs/engineering/rebaseline/D6-R2-P8-BLOCK-LEDGER.md`

**Produces:** durable record that the bounded D5 repair is accepted as the current contract and that GF-02 remains coherent.

- [ ] Mark OP-READ-01 `RESOLVED` only after exact-head CI passes.
- [ ] Record that the repair changes read projection/filter expressibility only: no owner, operation, Permission, business lifecycle, choreography, write/effect, or runtime semantics changed.
- [ ] Add a GF-02 revalidation note proving the new collection reads preserve owner boundaries and do not create cross-owner workflow authority.
- [ ] Advance the D6-R2 next action to re-render corrected B00 global IA; B10 remains suspended until that B00 is operator-adjudicated.

### Task 5: Verify the exact candidate head

**Verification:**

- [ ] Run PR CI full gate on the final HEAD.
- [ ] Require PASS for Product OAD baseline non-regression, canonical OAD 99/99, Permissions 30/30, H/A/S, generated TypeScript/Go projections, auth proof, performance knowledge proof, new operational-read proof, repository gate, and legacy runtime population 0.
- [ ] Inspect changed files and PR metadata; keep PR #61 draft and unmerged.
- [ ] Only after the final exact HEAD is green, re-render B00 using the already operator-approved IA masses and stop for operator review.
