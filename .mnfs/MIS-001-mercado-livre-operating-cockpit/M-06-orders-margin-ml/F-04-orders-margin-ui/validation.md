# F-04 Validation

```yaml
id: M-06-F-04
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-06
created: 2026-07-09
updated: 2026-07-09
```

## Scope

Validate the orders and margin UI against the real local backend and the persisted M-06 order/profitability state.

## Local Contract Validation

- Command:
  - `npm run test --workspace @marketplace-central/feature-orders -- OrdersPage.test.tsx`
- Result:
  - Passed with `6` tests.
  - Covers loading, error, empty, incomplete margin, complete margin, negative margin, and manual adjustment refresh behavior.
- Command:
  - `npm run test --workspace @marketplace-central/web -- AppRouter.test.tsx`
- Result:
  - Passed with `4` tests.
  - Confirms `/orders` route wiring through the app shell.
- Command:
  - `npm run test --workspace @marketplace-central/web -- ClientContext.test.tsx`
- Result:
  - Passed with `3` tests.
  - Confirms dev-time API base URL resolution behavior.
- Command:
  - `npm run test --workspace @marketplace-central/web -- viteProxy.test.ts`
- Result:
  - Passed with `1` test.
  - Confirms Vite proxy still covers inventory/product-links APIs without shadowing SPA routes such as `/orders`.
- Command:
  - `npm run build --workspace @marketplace-central/web`
- Result:
  - Passed.
  - Production bundle built successfully after adding the new feature package and route.

## Runtime Validation

- Environment:
  - Backend: `http://localhost:8080`
  - Web QA server: `http://localhost:5175`
  - Browser surface: Codex in-app browser
- Validation steps:
  - Opened `/orders?installation=inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98`
  - Confirmed the page rendered the persisted order list and selected-order detail
  - Confirmed list/detail surfaced:
    - order `2000012659424976`
    - order status `cancelled`
    - margin quality `Incomplete`
    - explicit flags `missing link`, `missing commission`
    - item profitability row for `MLB4834373620`
    - persisted manual freight adjustment `R$12.50`
  - Added a new manual adjustment through the UI:
    - category `commission`
    - amount `R$4.00`
    - reason `QA commission validation`
  - Confirmed the page refreshed after the write and now showed:
    - toast `Manual adjustment saved for order 2000012659424976.`
    - top summary moved from `Incomplete=1 / Negative=0` to `Incomplete=0 / Negative=1`
    - order margin quality changed to `Negative margin`
    - contribution updated to `-R$6.23`
    - margin updated to `-31.31%`
    - new manual adjustment entry `commission / R$4.00 / QA commission validation`

## Infrastructure Findings Fixed During Validation

- The app shell originally had no frontend package route for orders UI.
- Local Vite dev originally shadowed the SPA `/orders` route when `/orders` was proxied directly to the backend.
- Local dev also inherited a wrong `VITE_API_BASE_URL` (`http://localhost:8082`) in this environment.
- Final fix:
  - `feature-orders` package + `/orders` route added
  - dev API base URL now resolves safely to `http://localhost:8080` when no valid override is supplied
  - Vite proxy no longer shadows `/orders` or `/profitability`

## Live Provider Boundary

- Live validation covers:
  - UI readback of the real Mercado Livre order imported in F-01
  - UI readback of profitability inputs/snapshots created in F-02/F-03
  - real UI write of a persisted manual commission adjustment
  - real UI-triggered profitability snapshot recalculation and re-read
- Live validation does not cover:
  - multiple-order installation datasets in this session
  - mobile viewport/browser proof for F-04

## Open Blockers

- None for F-04.
