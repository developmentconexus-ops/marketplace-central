# M-06 F-03 Task 4 Report - Orders UI Realization Semantics

## Status

`DONE_WITH_CONCERNS`

Task 4 is implemented and verified by focused frontend tests, web regression tests, and a production web build. Built-in-browser desktop/mobile QA remains pending as required by the controller gate.

## Scope

- Implemented realization-aware presentation in the existing orders cockpit.
- Preserved SDK-provided values; no contribution, margin, or realization calculation was added to React.
- Kept the manual-adjustment actor required in the local `OrdersClient` request contract.
- Did not modify backend, migrations, OpenAPI, SDK, or provider/database behavior.
- Did not modify `apps/web/src/app/ClientContext.tsx`; the existing client already satisfied the required actor contract.

## TDD Evidence

### RED

Before production edits, cancelled and unknown order fixtures, realization fields for existing snapshots, the required-actor type assertion, and the Task 4 UI assertions were added to `OrdersPage.test.tsx`.

Command:

```powershell
npm run test --workspace @marketplace-central/feature-orders -- OrdersPage.test.tsx
```

Semantic RED result: `1` test file failed; `3` tests failed and `6` passed.

- Incomplete detail rendered `Missing inputs` instead of required `Data quality`.
- The filter had no `Not realized` option/summary semantics.
- Cancelled quality fell through to `All qualities`; operational flag text was lowercase and nullable contribution used `-`.
- Unknown realization rendered lowercase flag text and nullable contribution used `-`.

An earlier sandboxed attempt was blocked before Vitest started because workspace config resolution was denied; it was not counted as RED evidence.

### GREEN

Command:

```powershell
npm run test --workspace @marketplace-central/feature-orders -- OrdersPage.test.tsx
```

Result: `1` test file passed; `9/9` tests passed with no React test warning.

Coverage includes loading, error, empty, complete, incomplete, negative-margin, cancelled, unknown, quality filtering, and manual-adjustment refresh behavior.

## Regression And Build Evidence

Command:

```powershell
npm run test --workspace @marketplace-central/web -- AppRouter.test.tsx ClientContext.test.tsx viteProxy.test.ts
```

Result: `3` test files passed; `8/8` tests passed. Node emitted the existing `DEP0205 module.register()` deprecation warning.

Command:

```powershell
npm run build --workspace @marketplace-central/web
```

Result: production Vite build passed; `1783` modules transformed and output emitted in `2m 3s`. The build emitted dependency-level Node `DEP0205` and React Router/Lucide `use client` directive warnings; no Task 4 type or bundle error occurred.

## UI Semantics

- `not_realized` is represented in filter options, filter behavior, summary count, quality label, and neutral quality tone.
- Cancelled orders show `Not realized`, `Order cancelled`, and selected-detail text `Order not realized`.
- Unknown realization remains profitability quality `Incomplete` while separately showing operational state `Order state unknown`.
- Operational realization flags are separated from missing-input flags in selected detail.
- Missing-input flags appear under `Data quality`.
- Known revenue/inputs remain visible; nullable contribution and margin render an em dash (`—`).
- Cancelled orders do not receive a negative-margin label.
- State meaning is conveyed with text in addition to color.

## Self-Review

- Scoped production changes are limited to `OrdersPage.tsx`; no layout redesign, decorative pattern, motion, or new business math was introduced.
- Existing labelled selects, order buttons, loading/error/empty states, selected-order flow, item detail, inputs, and manual-adjustment workflow remain intact.
- Responsive structure remains breakpoint-based; the summary grid only expands from four to five semantic counters.
- Operational and data-quality callouts use explicit text and readable existing palette classes, so meaning is not color-only.
- Required actor request data is still constructed for every manual adjustment.

## Exact Changed Paths

- `packages/feature-orders/src/OrdersPage.tsx`
- `packages/feature-orders/src/OrdersPage.test.tsx`
- `.superpowers/sdd/m06-f03-task-4-report.md`
- `.superpowers/sdd/m06-f03-correction-report.md`

## Concerns And Remaining Evidence

- Browser evidence level is unit/regression/build only.
- Built-in-browser desktop/mobile QA is still pending and must verify the five-counter summary, filter behavior, selected cancelled/unknown detail, responsive layout, and manual-adjustment interaction visually.
- Dependency deprecation/directive warnings remain outside Task 4 scope.
- No commit, stage, reset, revert, or clean operation was performed.

## Fix After Independent Review

### Changed Path

- `packages/feature-orders/src/OrdersPage.test.tsx`
- `.superpowers/sdd/m06-f03-task-4-report.md`

No production, backend, SDK, contract, or migration file was changed for this review fix.

### Fixes

- Replaced global em-dash counts with assertions scoped to the selected-order `SurfaceCard`.
- Bound the `Contribution` and `Margin` labels to their direct value siblings and asserted each value is exactly `—` after selecting cancelled and unknown orders.
- Replaced the ambiguous `Not realized` text count with a direct assertion that the `Not realized` StatCard label's value sibling is exactly `1`.
- Reused existing semantic `section`, heading, label, and value structure; no test-only production hook was added.

### Verification

Exact command:

```powershell
npm run test --workspace @marketplace-central/feature-orders -- OrdersPage.test.tsx
```

Result: exit `0`; `1` test file passed; `9/9` tests passed in `1.87s`. Output was pristine: no stderr, React `act(...)` warning, or other test warning.

### Self-Review

- Cancelled and unknown assertions now prove selected-detail contribution and margin independently rather than succeeding from unrelated order-list rows.
- The summary assertion now proves the `Not realized` counter value, not merely repeated visible text.
- Query scope follows user-visible structure and remains resilient to unrelated em dashes elsewhere on the page.
- The fix is tests/report only and preserves all production behavior and shared dirty-worktree changes.
