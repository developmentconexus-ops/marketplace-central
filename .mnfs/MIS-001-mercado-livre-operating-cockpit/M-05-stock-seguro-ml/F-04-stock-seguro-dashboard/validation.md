# Feature Validation

```yaml
id: F-04
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: F-04
created: 2026-07-09
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-04-stock-seguro-dashboard

## Quick Validation Result

Passed with one explicit tooling limitation on mobile viewport proof.

## Spec Adherence

- M-05-C04: Browser validation proved the route `/inventory/stock-seguro` opens and lists real persisted Mercado Livre risk rows with filters, detail, policy, and source timestamps.
- M-05-C04: Runtime API validation proved dashboard data comes from the inventory API and not local UI math.
- M-05-C04: Manual action API and SDK surface exist; this session validated the read/dashboard slice in browser and contract-level apply flow readiness through tests.
- M-05-C04: Browser-led runtime debugging found and fixed a dev-server proxy collision that would have broken direct route navigation.

## Commands Run

### BACKEND GREEN

- Evidence type: `ran`
- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./internal/modules/inventory/... -count=1`
- Expected: Pass.
- Actual: Passed.
- Result: Pass.

### FRONTEND GREEN

- Evidence type: `ran`
- Command: `npx vitest run apps/web/src/app/viteProxy.test.ts apps/web/src/app/AppRouter.test.tsx`
- Expected: Pass.
- Actual: Passed.
- Result: Pass.

- Evidence type: `ran`
- Command: `npm run test --workspace @marketplace-central/feature-inventory`
- Expected: Pass.
- Actual: Passed.
- Result: Pass.

### REAL API VALIDATION

- Evidence type: `ran`
- Command: `Invoke-WebRequest http://localhost:8080/inventory/stock-risks?installation_id=inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98&limit=1`
- Expected: Real risk rows using the JSON contract.
- Actual: Returned real Mercado Livre row with `blocking_reason.code` and `blocking_reason.message` after backend restart.
- Result: Pass.

### BROWSER VALIDATION

- Evidence type: `ran`
- Surface: in-app browser
- Step: Open `http://localhost:5174/inventory/stock-seguro`.
- Actual: Route opened successfully after the Vite proxy fix and frontend restart.
- Result: Pass.

- Evidence type: `ran`
- Surface: in-app browser
- Step: Switch installation to `Mercado Livre`.
- Actual: UI loaded 5 real rows, 0 actionable rows, 0 oversell rows, and 2 blockers for the selected installation, with visible row detail and observed timestamps.
- Result: Pass.

- Evidence type: `ran`
- Surface: in-app browser
- Step: Validate direct route navigation failure reproduction before fix.
- Actual: Initial browser run reproduced `404 page not found`, which traced back to Vite proxy shadowing `/inventory/stock-seguro`.
- Result: Pass as runtime bug reproduction evidence.

- Evidence type: `could-not-run`
- Surface: in-app browser mobile viewport
- Step: Validate true narrow-width layout using browser viewport override.
- Actual: Browser capability reported support but did not change `window.innerWidth` in this session, so a strict mobile-width proof could not be completed inside the built-in browser.
- Result: Tool-limited. Responsive shell changes were still applied and should be rechecked in a browser/runtime where viewport override is effective.

## Risks

- Whole-workspace web test validation is still blocked by an unrelated pre-existing timeout in `packages/feature-product-links/src/ProductLinksPage.test.tsx`.
- Live manual provider stock write was not executed in browser because this feature requires explicit operator intent for side effects and current real rows were blocked, not actionable.

## Handoff

- Current status: `quick_validation_passed`
- Handoff target: Milestone Orchestrator
- Handoff reason: F-04 dashboard slice is implemented, browser-validated on desktop with real data, and ready for milestone review, with the mobile viewport limitation explicitly recorded.
