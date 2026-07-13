# M-13-integrated-operator-workspaces

```yaml
id: M-13
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-001
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-3
lifecycle_scope: milestone
```

## Mission

MIS-001 Mercado Livre Operating Cockpit MVP replan.

## Outcome

The browser becomes one operator cockpit with Overview, Products/Product 360,
Listings, Sales/Margin, and Operations workspaces connected by stable deep links,
shared installation context, common quality states, and a non-executing stock preview.

## Why This Milestone Exists

The repository already contains substantial product, link, stock, order, margin, and
integration behavior, but the UI mirrors feature packages and forces the operator to
reconstruct one business journey across unrelated routes.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Operator shell and attention queue | Establish route shell, shared installation context, Overview, and filtered deep links |
| F-02 | Catalog and Product 360 | Join canonical product facts, listings, stock, margin, sales, and evidence around CODPROD |
| F-03 | Listings workspace and stock preview | Move Product Links/Stock Seguro into one listing journey with simulation-only correction |
| F-04 | Sales and margin workspace | Make routable sale detail connect order, listing, product, Sankhya linkage, inputs, and margin |
| F-05 | Operations convergence and UX consistency | Merge channel setup/health, retire duplicate navigation, and align states/copy/responsiveness |

## Dependencies

- M-09 canonical product foundation passed.
- M-02 through M-06 implemented capability/API/SDK evidence.
- IC-003 route, identity, state, and simulation semantics.

## Risks

- Workspace composition could duplicate domain calculations in React.
- Concurrent features share AppRouter, Layout, SDK context, and UI primitives.
- Legacy real-write controls could remain reachable or appear to execute.
- Product/listing/sale fetches could lose installation context across deep links.

## Done Means

- Main navigation contains Visão geral, Produtos, Anúncios, Vendas, and Operações.
- `/products/:productId`, `/listings/:listingRef`, and `/sales/:providerOrderId`
  survive reload and preserve installation/entity context.
- Attention items open exact filtered objects and share IC-003 quality vocabulary.
- Product 360 and sale detail cross-link product, listing, stock, inputs, and margin.
- Stock action UI displays `Simulação`, current/proposed values, evidence, preview
  payload, and `executed=false`; no provider mutation is reachable.
- Legacy routes redirect with context, Portuguese copy is consistent, and applicable
  loading/error/empty/stale/unknown/conflict states are visible.

## Handoff

- Current status: Planned.
- Next owner: Milestone Orchestrator after M-09 passes.
- Next action: Order shared-seam ownership so only one feature writes router/layout at a time.
- Required files/evidence: M-13 feature validations, fixed-SHA review, `validation-result.md`.
- Blockers or open decisions: Stop if a required aggregate needs a client-side business calculation.

## Correction Handoff

- QA failure summary: Not applicable during planning.
- Correction scope: None until QA reports a named failed criterion.
- Attempts used/remaining: 0/2.
- Next artifact: `F-01-operator-shell-attention/feature.md`.
- Revalidation evidence required: M-13 contract plus browser evidence.
