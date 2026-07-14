# F-01-tanstack-adoption

```yaml
id: F-01
type: feature-plan
status: planned
owner: Feature Implementer
parent: M-05
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: feature
split_decision: single
split_reason: The shared cache/no-cache seam and three view migrations must be implemented and tested together for one intentional commit; the packet explicitly requires this coherent unit.
```

## Feature ID

F-01-tanstack-adoption

## Steps

1. Add the shared web-query package and install/lock TanStack Query v5 in apps/web and the three importing feature packages.
2. Wire the root QueryClientProvider and apps/web-owned no-cache fetch wrapper into ClientContext.
3. Replace the legacy ProductsPage with CatalogPage, catalog infinite/search hooks, quality-safe fact rendering, refresh, route, and navigation.
4. Migrate StockSeguroPage's risk read to the inventory query hook with IC-01 key/staleTime/as-of/refresh, preserving the manual action and refetch seam exactly.
5. Migrate OrdersPage's profitability-input read to the profitability query hook with IC-01 key/staleTime/as-of/refresh, preserving the import write seam.
6. Add/update component coverage, run the complete quick-validation commands, and write validation.md with exact evidence.

## Files Expected To Change

- `.mnfs/.../F-01-tanstack-adoption/spec.md`, `plan.md`, `validation.md`: feature execution artifacts.
- `packages/web-query/package.json`, `packages/web-query/src/index.ts`: shared query keys, stale times, QueryClient defaults, freshness helper, no-cache controller.
- `apps/web/package.json`, `package-lock.json`, `apps/web/src/main.tsx`, `apps/web/src/app/ClientContext.tsx`, `apps/web/src/app/AppRouter.tsx`, `apps/web/src/app/Layout.tsx`: dependency/root wiring and catalog route/nav.
- `packages/feature-products/package.json`, `src/CatalogPage.tsx`, `src/catalogQueries.ts`, `src/index.ts`: catalog implementation; delete legacy `ProductsPage.tsx` and test.
- `packages/feature-inventory/package.json`, `src/StockSeguroPage.tsx`, `src/StockSeguroPage.test.tsx`: stock query migration and coverage.
- `packages/feature-orders/package.json`, `src/OrdersPage.tsx`, `src/OrdersPage.test.tsx`: profitability query migration and coverage.
- `apps/web/src/app/TanStackQueries.test.tsx` (or equivalent): cross-cutting component coverage for pagination, staleTime, as-of, refresh, and no-cache header.

## Verification Commands

- Command: `npm view @tanstack/react-query version`
  - Satisfies criterion ID: M-05-C02 and install requirement
  - Expected result: a v5.x version is returned and recorded with the v5 API note.
- Command: `npm test` (from `apps/web`)
  - Satisfies criterion ID: M-05-C01, M-05-C02, M-05-C03
  - Expected result: all web/component tests pass, including cursor append/null-stop, staleTime remount, as-of, refresh header, and indicator update assertions.
- Command: `npm run test --workspace @marketplace-central/feature-products`
  - Satisfies criterion ID: M-05-C01, M-05-C02, M-05-C03
  - Expected result: catalog feature tests pass.
- Command: `npm run test --workspace @marketplace-central/feature-inventory`
  - Satisfies criterion ID: M-05-C02, M-05-C03
  - Expected result: stock-risk feature tests pass.
- Command: `npm run test --workspace @marketplace-central/feature-orders`
  - Satisfies criterion ID: M-05-C02, M-05-C03
  - Expected result: profitability feature tests pass.
- Command: `npm run build` (from `apps/web`)
  - Satisfies criterion ID: M-05-C05
  - Expected result: Vite/esbuild production build exits 0.
- Command: `rg -n "fetch\\(|listCatalogProducts|listInventoryStockRisks|listProfitabilityMarginInputs" packages/feature-products/src/CatalogPage.tsx packages/feature-inventory/src/StockSeguroPage.tsx packages/feature-orders/src/OrdersPage.tsx packages/web-query/src apps/web/src
  - Satisfies criterion ID: M-05-C05
  - Expected result: no direct fetch or deprecated/manual read in migrated views; SDK calls occur in query hooks and only the apps/web fetch wrapper owns fetch customization.

## QA Steps

- Step: Render catalog with two mocked envelope pages, trigger the next-page control, and verify both pages append while a null cursor disables continuation.
  - Expected result: one call per page, no first-page repeat, terminal `Fim da lista` state.
- Step: Render a migrated view with an envelope, click its refresh control, and inspect the captured SDK fetch init.
  - Expected result: refresh is disabled while pending, GET carries `Cache-Control: no-cache`, and the new `as_of` is displayed.
- Step: Render a null/quality-flag catalog fact.
  - Expected result: unknown value and quality flags are visible; no zero/default is introduced.

## Rollback/Risk Notes

- Keep `packages/sdk-runtime/src/**`, Go, OpenAPI, and out-of-scope pages untouched. If no-cache cannot be implemented in apps/web's `fetchImpl`, stop as blocked rather than modifying SDK runtime.
- Existing feature tests must receive a QueryClientProvider; use isolated QueryClients per test to prevent cache leakage.
- Preserve F-02 seams: `applyInventoryManualStockAction` plus its post-action manual refetch and `importProfitabilityMarginInputs` write plus its existing refresh behavior remain intact.
- A single intentional commit is required; do not merge or modify main.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Implement and record validation evidence.
- Required files/evidence: spec, changed paths, exact verification results, validation.md, commit SHA.
- Blockers or open decisions: None.
