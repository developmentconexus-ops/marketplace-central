# F-03-listings-stock-preview

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-13
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-3
lifecycle_scope: feature
```

## Mission

MIS-001 MVP operator journey.

## Milestone

M-13 Integrated Operator Workspaces.

## Brief

Unify listing import/link resolution and Stock Seguro into `/listings` plus listing
detail, with contextual Product/Sale links and a non-executing stock preview.

## Inputs

- Product Links and Inventory SDK methods/read models.
- IC-003 ListingRef, attention/quality states, route encoding, and StockSimulation.

## Inputs/Outputs

- Input: installation, link/risk/attention filters, listing composite identity.
- Output: listing status, product link/evidence, stock facts/risk/recommendation,
  recent sale refs, audit, and StockSimulation with `executed=false`.

## Interaction Model

- Aggregate list supports batch filtering; detail handles link resolution and preview.
- Approve/reject/manual-resolve refetches the listing and attention queue.
- Review simulation opens current/proposed/evidence/payload; it has Close, not Execute.

## State Model

Link states `unresolved|conflict|resolved|rejected`; source quality
`current|stale|unknown|conflict`; simulation `draft|reviewed`. Conflict/unresolved
blocks simulation review.

## Negative Scenarios

- Invalid listing ref: `invalid_identity`.
- Ambiguous link: `identity_conflict`, no stock preview.
- Missing recommendation/evidence: `simulation_unavailable`, no payload.
- Source unavailable: a prior trustworthy listing fact stays stale; without one it
  is null/unknown + nonblank reason.

## Expected Output

The operator understands one listing's product, stock risk, evidence, and proposed
correction without entering separate Product Links and Stock Seguro applications.

## Constraints

- Remove/disable the M-13 route's call to manual provider apply; do not delete historical audit.
- No authentication or provider-write implementation.
- Owned paths: feature-product-links/inventory convergence into listings, shared listing
  UI, and this feature root. AppRouter/Layout/context/legacy redirects are read-only
  because F-01 wires stable outlets before releasing the seam.
- `packages/ui` is read-only; listing-specific components stay in the listing-owning
  feature package.
- Forbidden paths: provider adapter mutation, M-06 evidence, Product/Sales internals.

## Criteria IDs

- M-13-C01 Navigation and attention closure.
- M-13-C02 Product/listing/sale cross-links.
- M-13-C04 Proportional security and simulation.
- M-13-C05 State and responsive consistency.
- M-13-C06 SDK-only thin client.

## Validation Expectations

- Fixture conflict blocks simulation with visible reason.
- Resolved fixture shows current 8, proposed 6, `executed=false`, and no execute control.
- Legacy Product Links/Stock Seguro URLs redirect with filters preserved.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: Briefed.
- Next owner: Feature Implementer after shell seam release.
- Next action: Create spec/plan and converge existing components without backend breadth.
- Required files/evidence: IC-003, M-04/M-05 validation evidence, and this feature's `validation.md`.
- Blockers or open decisions: Stop if disabling execution would violate a public contract.
