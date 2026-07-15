# F-02-catalogo-workspace

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-04
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contracts: IC-05 (route `/catalogo`, catalog ns keys, states), R-01 screen 1f, R-02 (current Products page facts). ADR-15/17.

## Milestone

M-04 read-workspaces. Depends on F-01 (product detail exists as row target).

## Brief

Rebuild Catálogo (wireframe 1f) at `/catalogo`, replacing the legacy `/products` page: product table (CODPROD, título, GTIN state, custo, estoque disponível, completude, nº anúncios, ⚑ ERP provenance), filters per 1f (search, completude band, GTIN missing, sem anúncio), row click → `/catalogo/produtos/:productId`. Server state via existing catalog endpoints through TanStack Query (`catalog` ns, staleTime 300s) — this page migrates OFF direct fetch per its migration brief. Freshness indicator + manual refresh (createRefreshableFetch path).

EARS:
- While the operator filters/searches, when state changes, URL query shall encode it and reload restore it.
- While a product has no cost, when the row renders, custo shall show UnknownValue "custo?" and completude shall count the gap — never "R$ 0,00".
- While catalog data is stale (>300s), when the page shows, FreshnessIndicator shall show fetch time; manual refresh forces no-cache fetch and updates indicator.
- While legacy `/products?…` is visited, when router resolves, redirect shall land `/catalogo?…` (M-02 redirect, re-proven here after page swap).

## Inputs

- R-01 §1f columns/filters, R-02 current Products implementation (parity checklist source), existing catalog list endpoint + sdk-runtime, IC-05 keys/components, F-01 detail route.

## Expected Output

- `/catalogo` page, TanStack Query only; legacy Products component deleted (not orphaned).
- Parity checklist in spec.md from current Products features; each item kept or explicitly dropped-with-reason.
- Component tests: URL round-trip, unknown-cost render, refresh flow, row navigation.

## Constraints

- No new server endpoints; existing catalog reads only.
- No enrichment editing here (detail page owns it).
- Deletion of legacy page includes its direct-fetch helpers if unused elsewhere (grep proof).

## Negative Scenarios

- Catalog API error → full-page ErrorState with retry.
- Empty catalog → EmptyState with hint "Verifique a sincronização com o ERP em Integrações".
- Filter yielding zero → EmptyState + clear-filters affordance.

## Validation Expectations

- Vitest output: round-trip, unknown-render, navigation tests green.
- Browser proof: 1f screenshot; F5 fidelity with filters; redirect from `/products`.
- Grep proof: no `useEffect`-fetch in the new page; legacy component gone.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (after F-01 accepted).
- Next action: compile context pack; read R-01 §1f + R-02 Products section + IC-05.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
