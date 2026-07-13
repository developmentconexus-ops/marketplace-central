# Current UI Journey Inventory

```yaml
id: R-004
type: research
status: verified
owner: Codebase Investigator
parent: MIS-001
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-0
lifecycle_scope: support
```

## Topic

Current Marketplace Central browser routes, implemented capabilities, and UX fragmentation.

## Sources Checked

- `apps/web/src/app/AppRouter.tsx` and `Layout.tsx`.
- `packages/feature-products`, `feature-classifications`, `feature-marketplaces`,
  `feature-integrations`, `feature-product-links`, `feature-inventory`,
  `feature-orders`, and `feature-simulator` source.
- `packages/sdk-runtime/src/index.ts` for already-exposed operations.

## Findings

- Current routes mirror technical packages: `/products`, `/classifications`,
  `/marketplaces`, `/integrations`, `/product-links`, `/inventory/stock-seguro`,
  `/orders`, and `/simulator`.
- There is no routable Product 360 or Listing detail. Order detail is selected inside
  `/orders` with query state rather than a stable `/sales/:providerOrderId` route.
- Products, Classifications, and Simulator repeat catalog search/filter/selection.
- Marketplaces and Integrations overlap installation lifecycle and health actions.
- Product Links, Stock Seguro, and Orders do not provide reciprocal deep links.
- Orders already composes order, input, adjustment, snapshot, and margin data; its
  adjustment is disabled because the route wrapper supplies no trusted operator.
- SDK runtime already exposes listing imports, link candidate generation, assisted
  Sankhya linkage, stock risks, orders, inputs, adjustments, and snapshots.
- Stock Seguro currently exposes a real manual apply path; the approved MVP requires
  preview semantics and default-disabled execution instead.

## Recommendation

Reuse existing SDK operations in five operator workspaces. Add shared installation
context, stable entity deep links, a cross-domain attention vocabulary, and consistent
loading/error/empty/stale/conflict states before adding backend breadth.

## Impact On Mission

M-13 is a UX and composition outcome over existing capabilities. It is not a request
to merge backend modules or create an endpoint per screen.

## Handoff

- Current status: Verified from repository source; no browser or tests ran.
- Next owner: M-13 Milestone Orchestrator.
- Next action: Use IC-003 and feature briefs to consolidate journeys.
- Required files/evidence: this note and M-13 feature validation artifacts.
- Blockers or open decisions: None - live visual state remains M-14 evidence.
