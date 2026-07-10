# Implementation Report

```yaml
id: F-04
type: implementation-report
status: ready_for_review
owner: Feature Implementer
parent: F-04
created: 2026-07-09
updated: 2026-07-09
```

## Claimed Scope

Implemented the Stock Seguro dashboard vertical slice end-to-end: inventory risk listing API, manual stock apply API, Postgres action persistence, SDK runtime client methods, route wiring, and operator UI backed by persisted Mercado Livre listing snapshots and internal sellable stock evidence.

## Changed Paths

- `apps/server_core/internal/composition/root.go`
- `apps/server_core/internal/modules/inventory/**`
- `apps/server_core/migrations/0026_inventory_stock_actions.sql`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`
- `packages/sdk-runtime/src/index.test.ts`
- `packages/feature-inventory/**`
- `apps/web/src/app/AppRouter.tsx`
- `apps/web/src/app/Layout.tsx`
- `apps/web/src/app/AppRouter.test.tsx`
- `apps/web/src/app/viteProxy.test.ts`
- `apps/web/vite.config.ts`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-04-stock-seguro-dashboard/spec.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-04-stock-seguro-dashboard/plan.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-04-stock-seguro-dashboard/validation.md`

## TDD And Runtime Evidence

GREEN:

- Command: `go test ./internal/modules/inventory/... -count=1`
- Result: Passed.

- Command: `npx vitest run apps/web/src/app/viteProxy.test.ts apps/web/src/app/AppRouter.test.tsx`
- Result: Passed.

- Command: `npm run test --workspace @marketplace-central/feature-inventory`
- Result: Passed.

RUNTIME FIXUP 1:

- Reproduction: in-app browser navigation to `http://localhost:5174/inventory/stock-seguro` rendered `404 page not found` because Vite proxied the SPA route to backend.
- Change: narrowed Vite proxy keys from broad `/inventory` and `/product-links` prefixes to API-only subpaths.
- Result: browser route opened correctly after frontend restart.

RUNTIME FIXUP 2:

- Reproduction: real `GET /inventory/stock-risks` returned `blocking_reason.Code` and `blocking_reason.Message` instead of the JSON contract shape.
- Change: added JSON tags to inventory blocking/action domain structs and restarted backend.
- Result: runtime response now returns `blocking_reason.code` and `blocking_reason.message`.

## Concerns

- Desktop browser validation succeeded with real backend data.
- Mobile browser validation is only partial in this session because the in-app browser viewport capability reported support but did not change `window.innerWidth`; responsive shell hardening was still applied in `Layout.tsx`, but true narrow-width browser proof remains tool-limited.
- Existing unrelated web workspace failure remains in `packages/feature-product-links/src/ProductLinksPage.test.tsx` when running the whole `@marketplace-central/web` suite.
