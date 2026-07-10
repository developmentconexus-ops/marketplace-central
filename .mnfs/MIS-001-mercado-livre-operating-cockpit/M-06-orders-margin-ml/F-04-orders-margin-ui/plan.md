# F-04 Orders Margin UI Plan

## Task 1 - Route and package skeleton

- Create `packages/feature-orders`.
- Add `/orders` route to `apps/web`.
- Add navigation entry and CSS source registration.
- Add router coverage for the new route.

## Task 2 - Orders workspace page

- Build `OrdersPage` with installation selection and quality filter.
- Load orders, snapshots, inputs, and adjustments from `sdk-runtime`.
- Render list/detail layout with explicit loading, error, and empty states.

## Task 3 - Manual adjustment workflow

- Add operator adjustment form for order/item scope.
- Persist via `createProfitabilityManualAdjustment`.
- Refresh adjustments and snapshots after save.
- Keep calculations server-owned by calling snapshot recalculation instead of local math.

## Task 4 - Verification

- Add page tests for the required UI states.
- Run targeted frontend tests for the new package and router.
- Run the app and validate in the in-app browser against the live local backend.
- Record validation evidence in `validation.md`.
