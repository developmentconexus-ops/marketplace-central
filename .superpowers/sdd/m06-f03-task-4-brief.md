# M-06 F-03 Task 4 Brief - Orders UI Semantics And Regression

Source plan: `docs/superpowers/plans/2026-07-09-m06-f03-order-realization.md`, Task 4.

## Goal

Present realized, cancelled, and unknown order profitability semantics accurately in the existing orders cockpit without doing business calculations in React.

## Product Context

Read `packages/feature-orders/PRODUCT.md`. Preserve the established product UI and component vocabulary. Separate operational state from data quality, use state text in addition to color, and keep desktop/mobile behavior intact. This is a semantic correction, not a visual redesign.

## Files

- Modify `packages/feature-orders/src/OrdersPage.tsx`
- Modify `packages/feature-orders/src/OrdersPage.test.tsx`
- Modify only if required by type integration: `apps/web/src/app/ClientContext.tsx`
- Append `.superpowers/sdd/m06-f03-correction-report.md`

## Contract

- Consume SDK `realization_state`, snapshot quality, and flags from Task 3.
- Keep the F-02 manual-adjustment actor required in `OrdersClient`; change the local actor field from optional to required.
- Do not calculate contribution, margin, or realization in React.

## TDD

1. Add cancelled and unknown fixtures before production edits.
2. Assert cancelled orders render `Not realized`, `Order cancelled`, and an em dash (`—`) for contribution and margin.
3. Assert unknown orders render `Incomplete` and `Order state unknown`.
4. Assert cancelled orders do not render a negative-margin label.
5. Make the local `OrdersClient` actor field required.
6. Run RED:

```powershell
npm run test --workspace @marketplace-central/feature-orders -- OrdersPage.test.tsx
```

## Implementation

- Add `not_realized` to quality filter options, summary semantics, `qualityLabel`, and `qualityTone`.
- Render realization flags as operational state, not as missing inputs.
- Use `Data quality` for missing-input flags.
- Use `Order not realized` for cancellation.
- Preserve known input display and render nullable contribution/margin through the existing em-dash formatter.
- Keep loading, error, empty, selected-order, item, input, and adjustment behavior intact.
- Avoid new decorative patterns, motion, or layout churn.

## GREEN

```powershell
npm run test --workspace @marketplace-central/feature-orders -- OrdersPage.test.tsx
npm run test --workspace @marketplace-central/web -- AppRouter.test.tsx ClientContext.test.tsx viteProxy.test.ts
npm run build --workspace @marketplace-central/web
```

## Limits

- No backend, migration, OpenAPI, SDK, or business-math changes.
- Preserve the required manual-adjustment actor contract.
- Do not commit, stage, reset, revert, or clean.
- Unit/build evidence is not built-in-browser desktop/mobile QA; browser evidence is a later controller gate.
- Preserve all unrelated dirty-worktree changes.
