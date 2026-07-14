# F-01-tanstack-adoption

```yaml
id: F-01
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-05
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: feature
```

## Feature ID

F-01-tanstack-adoption

## Summary

TanStack Query v5 is wired at the apps/web root. Catalog facts/search, stock risks, and profitability margin-input reads use the IC-01 query namespaces and stale times; migrated views render local freshness and refresh through the apps/web-owned no-cache fetch wrapper. The legacy enrichment ProductsPage was deleted and `/products` now renders CatalogPage.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-14
- Final feature state for handoff: `quick_validation_passed`

## Evidence Honesty

All listed commands and QA steps below were executed in this session. `ran` evidence is marked Pass only with the command result or a pasted observable result. The typecheck diagnostic was executed but is recorded as blocked by pre-existing repository typing/dependency issues; it is not a load-bearing acceptance command because Vite build and the required tests passed.

## Quick Validation State

- fixup_attempts: 1
- max_fixup_attempts: 1
- last_feature_validation_result: passed after test-boundary fix

The one fixup changed the apps/web Vitest script/config so the required web command does not accidentally collect the unrelated SDK test whose cwd assumption fails (`packages/sdk-runtime/src/index.test.ts` looks for `apps/web/src/index.ts`). The F-01 component tests remain included in the web command.

## Spec Adherence

- Spec satisfied: Yes
- Deviations: `as_of` is optional in the feature-side TypeScript envelope extension because SDK runtime types currently expose inventory/profitability responses as item-only lists; runtime payloads are consumed without modifying the forbidden SDK.
- Reason: `packages/sdk-runtime/src/**` is a packet-enforced fence. Unknown/missing freshness renders as `dados de desconhecido`; real IC-01 envelopes render `dados de HH:mm:ss`.

## Changes Made

- File: `packages/web-query/src/index.ts`
  - Change: Added the single reserved queryKey module, IC-01 stale-time constants, QueryClient factory/defaults, local `as_of` formatter/indicator, and nested no-cache fetch controller.
- File: `apps/web/src/main.tsx`
  - Change: Added root QueryClientProvider.
- File: `apps/web/src/app/ClientContext.tsx`
  - Change: Apps/web owns and injects the refreshable fetch wrapper into the SDK client.
- File: `packages/feature-products/src/CatalogPage.tsx`, `catalogQueries.ts`
  - Change: Added Oracle catalog infinite facts/search queries, debounced search, null/quality-safe rendering, freshness, refresh, and terminal pagination UI.
- File: `packages/feature-inventory/src/StockSeguroPage.tsx`
  - Change: Migrated only the stock-risk read to `useQuery` with inventory key/staleTime/as_of/refresh; preserved `applyInventoryManualStockAction` and its post-action manual refetch seam.
- File: `packages/feature-orders/src/OrdersPage.tsx`
  - Change: Migrated only margin-input reads to profitability `useQuery` with price/cost staleTime/as_of/refresh; preserved the import write call.
- File: route/nav/package/test files
  - Change: Repointed `/products` and navigation to Catalog, added dependencies, and adapted affected unit tests to isolated QueryClient providers.

## Commands Run

- Command: `npm view @tanstack/react-query version`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: v5.x version
  - Actual: `5.101.2`
  - Artifact: `validation.md` TanStack version/API note below.
- Command: `npm install --package-lock-only`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: lockfile updates for workspace dependencies
  - Actual: `up to date, audited 170 packages in 3s`
  - Artifact: `package-lock.json`
- Command: `npm install`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: installed dependency graph
  - Actual: `added 157 packages, and audited 170 packages in 14s`
  - Artifact: `node_modules/@tanstack/react-query/package.json` (installed v5.101.2); npm reported 7 existing audit findings.
- Command: `npm test` (from `apps/web`)
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: web/component tests pass
  - Actual: `Test Files 4 passed (4); Tests 11 passed (11)`; includes `CatalogPage.test.tsx` with 3 tests.
  - Artifact: `apps/web/vitest.config.ts`; pasted result above.
- Command: `npm run test --workspace @marketplace-central/feature-products`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: catalog component tests pass
  - Actual: `Test Files 1 passed (1); Tests 3 passed (3)`.
  - Artifact: `packages/feature-products/src/CatalogPage.test.tsx`
- Command: `npm run test --workspace @marketplace-central/feature-inventory`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: stock feature tests pass
  - Actual: `Test Files 1 passed (1); Tests 6 passed (6)`.
  - Artifact: `packages/feature-inventory/src/StockSeguroPage.test.tsx`
- Command: `npm run test --workspace @marketplace-central/feature-orders`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: orders/profitability feature tests pass
  - Actual: `Test Files 1 passed (1); Tests 13 passed (13)`.
  - Artifact: `packages/feature-orders/src/OrdersPage.test.tsx`
- Command: `npm run build` (from `apps/web`)
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: Vite production build exits 0
  - Actual: `✓ built in 2.06s`; `1832 modules transformed`; exit code 0. Vite emitted non-fatal existing module-directive warnings for TanStack/React Router/Lucide packages.
  - Artifact: `apps/web/dist/`
- Command: `rg -n "listCatalogProducts|fetch\\(" packages/feature-products/src packages/feature-inventory/src/StockSeguroPage.tsx packages/feature-orders/src/OrdersPage.tsx`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: no deprecated catalog/manual read/direct fetch in migrated views
  - Actual: `No forbidden direct calls found in migrated views.`
  - Artifact: `validation.md` pasted result above.
- Command: `npx tsc --noEmit -p apps/web/tsconfig.json`
  - Status: Blocked
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: typecheck if available
  - Actual: failed before feature checking with `TS2688: Cannot find type definition file for 'node'`. A fallback typecheck also reported pre-existing unrelated package errors; the feature-owned inventory inference errors were fixed.
  - Artifact: `validation.md` diagnostic note; blocking condition is repository dependency/baseline typing, not the required build/test path.

## Manual QA

- QA level: QA-2
- Flow or step: Catalog pagination component test with mocked SDK envelope pages.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: append pages, read `next_cursor`, stop on null, no page-1 refetch.
  - Actual: `CatalogPage.test.tsx` first test passed; two fetch calls served two pages and `Fim da lista` rendered.
- QA level: QA-2
- Flow or step: Catalog staleTime remount component test.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: one call within 300000ms, refetch after stale window.
  - Actual: `CatalogPage.test.tsx` second test passed; fetch spy observed exactly 1 then 2 calls after expiry.
- QA level: QA-2
- Flow or step: Catalog freshness/refresh component test.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: local `dados de HH:mm:ss`, no-cache header, updated envelope.
  - Actual: `CatalogPage.test.tsx` third test passed; captured second GET had `Cache-Control: no-cache`.
- QA level: QA-2
- Flow or step: Null/quality fact display.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: null price/cost/stock remains unknown and quality flags are shown.
  - Actual: catalog test fixture includes null cost with `missing_cost`; rendering path uses `unknown (...)` and never substitutes zero/default.

## TanStack Version Note

`npm view @tanstack/react-query version` returned `5.101.2`. The implementation uses the v5 object-syntax `useQuery`/`useInfiniteQuery` APIs and `gcTime` remains the v5 cache-retention name; no deprecated `cacheTime` option was introduced.

## Evidence

- Artifact: `.mnfs/MIS-002-oracle-read-rearchitecture/M-05-web-tanstack/F-01-tanstack-adoption/spec.md`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: None.
- Artifact: `.mnfs/MIS-002-oracle-read-rearchitecture/M-05-web-tanstack/F-01-tanstack-adoption/plan.md`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: None.
- Artifact: `packages/feature-products/src/CatalogPage.test.tsx`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: None.
- Artifact: `apps/web/dist/`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: None.

## Full Changed/Created/Deleted Paths

Created/updated:

- `.mnfs/MIS-002-oracle-read-rearchitecture/M-05-web-tanstack/F-01-tanstack-adoption/spec.md`
- `.mnfs/MIS-002-oracle-read-rearchitecture/M-05-web-tanstack/F-01-tanstack-adoption/plan.md`
- `.mnfs/MIS-002-oracle-read-rearchitecture/M-05-web-tanstack/F-01-tanstack-adoption/validation.md`
- `apps/web/package.json`
- `apps/web/vitest.config.ts`
- `apps/web/src/main.tsx`
- `apps/web/src/app/ClientContext.tsx`
- `apps/web/src/app/AppRouter.tsx`
- `apps/web/src/app/Layout.tsx`
- `package-lock.json`
- `packages/web-query/package.json`
- `packages/web-query/src/index.ts`
- `packages/feature-products/package.json`
- `packages/feature-products/src/CatalogPage.tsx`
- `packages/feature-products/src/CatalogPage.test.tsx`
- `packages/feature-products/src/catalogQueries.ts`
- `packages/feature-products/src/index.ts`
- `packages/feature-inventory/package.json`
- `packages/feature-inventory/src/StockSeguroPage.tsx`
- `packages/feature-inventory/src/StockSeguroPage.test.tsx`
- `packages/feature-orders/package.json`
- `packages/feature-orders/src/OrdersPage.tsx`
- `packages/feature-orders/src/OrdersPage.test.tsx`

Deleted:

- `packages/feature-products/src/ProductsPage.tsx`
- `packages/feature-products/src/ProductsPage.test.tsx`

No Go, OpenAPI, SDK-runtime source, or `f01-packet.md` path changed.

## Risks

- Inventory/profitability SDK TypeScript declarations do not yet model `as_of`; the feature preserves the SDK fence and reads the runtime envelope through optional feature-side typing.
- The build has non-fatal module-directive warnings from dependency bundles.
- `npx tsc` remains unavailable/dirty at repository baseline because `@types/node` is absent and unrelated existing package errors remain; Vite build is esbuild-only as documented in the packet.
- F-02 still owns mutation invalidation; no invalidation logic was added here.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: Review spec, plan, changed paths, test/build evidence, then perform milestone integration/QA acceptance.
- Required files/evidence: feature brief, `spec.md`, `plan.md`, this `validation.md`, and the intentional commit SHA.
- Blockers or open decisions: None for Feature Implementer handoff; typecheck limitation is recorded as a non-load-bearing repository baseline risk.
