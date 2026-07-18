# F-02-simulador-ui — Changed-Path Reconciliation

Core §4 ob.7 reconciliation of DECLARED slice write_sets vs ACTUAL changed paths.
Range: `d9d53d5^..3265353` (F02-S1..S6). Hub ruling: F-plan-2 ACCEPTED — additive
per-slice in-page wiring to `PricingPage.tsx` is legal (inside F-02 milestone
ownership `apps/web/src/pages/precos/**`, additive within granted scope). This file
is the required reconciliation, not a scope-change request.

## Declared — changed as planned

| Path | Slice | Note |
|------|-------|------|
| `apps/web/src/pages/precos/PricingPage.tsx` | S1 | Shell/page (declared S1). |
| `apps/web/src/pages/precos/PricingPage.test.tsx` | S1 | Declared S1. |
| `apps/web/src/pages/precos/index.ts` | S1 | Barrel (declared S1). |
| `apps/web/src/routes/precos.tsx` | S1 | Route rewire (declared S1, F-02 owned path). |
| `apps/web/src/app/AppRouter.test.tsx` | S1 | APPTEST-M07 grant, :93-97 /precos case only. |
| `apps/web/src/pages/precos/DecompositionPanel.tsx` (+`.test.tsx`) | S2 | Declared S2. |
| `apps/web/src/pages/precos/ParamsDrawer.tsx` (+`.test.tsx`) | S3 | Declared S3. |
| `apps/web/src/pages/precos/DifalDrawer.tsx` (+`.test.tsx`) | S3 | Declared S3. |
| `apps/web/src/pages/precos/MarketComparison.tsx` (+`.test.tsx`) | S4 | Declared S4. |
| `apps/web/src/pages/precos/ApplyPriceAction.tsx` / `ApplyPrice.test.tsx` | S5 | Declared S5. |
| `packages/feature-simulator/src/PricingSimulatorPage.tsx` | S6 | Declared S6 (retheme). |

## Changed-undeclared — additive, justified (F-plan-2)

| Path | Slices | Justification |
|------|--------|---------------|
| `apps/web/src/pages/precos/PricingPage.tsx` | S2, S3, S4, S5 | Per-slice ADDITIVE in-page wiring (F-plan-2): each component slice mounts its panel/drawer into the region shell S1 established (region-decomposicao, params/difal triggers+drawers, region-comparacao, region-aplicar) so /precos stays demoable at every slice boundary. No rewrite of prior slices' regions — each edit adds imports/state/JSX only. Inside F-02 owned path; sole-committer serialized. The slice cards' write_sets omitted this file (plan hygiene), which C05/C06 (page must SHOW the surfaces) require. |
| `apps/web/src/pages/precos/PricingPage.test.tsx` | S3, S4, S5 | Extended alongside each wiring: added mutation/query mocks (putPricingProfile, listPricingDifal, putPricingDifalOverride, listIntegrationInstallations, listListingsByProduct) + child stubs (MarketComparison, ApplyPriceAction) so the page test stays hermetic. Same owned path, additive. |

## Declared-but-unchanged

None. Every declared slice target was touched. No planned path was dropped.

## Boundary compliance

- Forbidden paths (feature.md): `apps/web/src/app/**` (except APPTEST-M07 :93-97), `packages/ui/**`, other routes/pages, `sdk-runtime/**`, backend — **all untouched** by F-02. Verified: F-02 diff hits only `apps/web/src/{pages/precos/**, routes/precos.tsx, app/AppRouter.test.tsx}` + `packages/feature-simulator/src/PricingSimulatorPage.tsx`.
- MarketComparison consumes IC-03 via a LOCAL base-URL resolver (not imported from `app/ClientContext`) precisely to avoid coupling to a forbidden path — see F-plan-3.
- Single-writer: orchestrator sole committer (F-boot-2); one clean commit per GREEN slice.
