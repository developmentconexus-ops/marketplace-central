# M-06 F-03 Final Review Fix Brief

## Goal

Resolve the whole-feature review's complete Critical/Important findings before real-environment validation.

## Finding 1 - Critical: Snapshot Replacement Is Truncated And Nondeterministic

Evidence: `CalculateSnapshots` builds snapshots from Go maps, truncates the resulting slice before deriving replacement order IDs, and sends only that subset to `ReplaceSnapshots`. A low limit can omit a cancelled order and leave its previous realized row untouched.

### Required behavior

- Write a failing application regression before production edits.
- Build and deterministically sort the complete bounded snapshot set.
- Derive replacement order IDs from the complete set.
- Persist the complete set so every calculated order is replaced atomically.
- Apply the request `limit` only to the returned `Items`, after persistence.
- `CalculatedCount` reports the complete number calculated/persisted, not only returned items.
- Stable sort key: provider order ID, then scope with order before item, then provider item ID, then provider variation ID.
- Test a low response limit with multiple orders/snapshots and a previously realized cancelled order; prove the store receives/replaces all relevant order IDs and the cancelled snapshot is `not_realized`.
- Extend the PostgreSQL-gated regression if useful to prove stale rows for selected order IDs are removed/replaced. Label skip versus real evidence accurately.
- Preserve nullable monetary values and all realization semantics.

## Finding 2 - Important: Editable/Default Actor Is Not Trustworthy

Evidence: `OrdersPage` defaults the form to `Leandro` and derives `actor_id` from editable display text. The app currently has no authenticated/session principal.

### Required safe behavior

- Write failing UI/type tests before production edits.
- Remove editable/defaulted operator identity from the adjustment form.
- Accept an immutable optional `operator: ProfitabilityActor` through `OrdersPageProps`.
- When a valid operator is supplied, send its exact `actor_type`, `actor_id`, and optional `actor_name`; never derive ID from display text.
- When trustworthy operator identity is absent/blank, show a clear non-color-only explanation and disable manual-adjustment submission. Do not invent an environment identity or claim authentication.
- Existing `AppRouter` may omit the operator until a real authenticated application context exists; this must result in safely disabled adjustment creation.
- Keep the reason/amount/scope/category/item controls and adjustment listing intact.
- Do not change backend actor invariants or SDK/OpenAPI contracts.

## Root Cause And TDD

- Use `superpowers:systematic-debugging` Phase 1-4 and `superpowers:test-driven-development`.
- Record explicit hypotheses and reproduce each failing behavior with focused RED tests before any production edit.
- Implement the smallest root-cause fixes, not symptom patches.

## Verification

```powershell
cd apps/server_core
$env:GOCACHE="$PWD\.gocache"
go test ./internal/modules/profitability/application ./internal/modules/profitability/adapters/postgres -count=1
cd ../../..
npm run test --workspace @marketplace-central/feature-orders -- OrdersPage.test.tsx
npm run test --workspace @marketplace-central/web -- AppRouter.test.tsx ClientContext.test.tsx viteProxy.test.ts
npm run build --workspace @marketplace-central/web
```

- Run `gofmt` for changed Go files.
- Run the profitability orders-import boundary `rg`.
- Output should be pristine except accurately identified pre-existing dependency warnings.

## Files

Expected:

- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`
- optionally extend `apps/server_core/internal/modules/profitability/adapters/postgres/profit_snapshot_integration_test.go`
- `packages/feature-orders/src/OrdersPage.tsx`
- `packages/feature-orders/src/OrdersPage.test.tsx`
- evidence reports only

Do not modify unrelated files, contracts, migrations, backend actor validation, or provider behavior.

## Report

Write `.superpowers/sdd/m06-f03-final-review-fix-report.md` with root-cause evidence, RED/GREEN output, exact paths, database evidence level, actor safety behavior, self-review, and remaining real gates. Append a concise correction section to `.superpowers/sdd/m06-f03-correction-report.md`.

Do not stage or commit.
