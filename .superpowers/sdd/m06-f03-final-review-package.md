# M-06 F-03 Whole-Feature Review Package

## Baseline And Review Method

- Branch: `main`; HEAD remains `8dba7db` because the approved plan forbids task commits in this shared dirty worktree.
- Profitability/orders feature files are largely untracked and OpenAPI/SDK had pre-existing dirty hunks. A clean commit-range diff is unavailable.
- Review the scoped surfaces below as the complete F-03 correction surface. Do not broaden into unrelated dirty work.

## Requirements

- Design: `docs/superpowers/specs/2026-07-09-m06-f03-order-realization-design.md`
- Plan: `docs/superpowers/plans/2026-07-09-m06-f03-order-realization.md`
- Task briefs: `.superpowers/sdd/m06-f03-task-{1,2,3,4}-brief.md`
- Aggregate report: `.superpowers/sdd/m06-f03-correction-report.md`

## Task 1 Surface

- `apps/server_core/internal/modules/profitability/domain/order_fact.go`
- `apps/server_core/internal/modules/profitability/ports/order_reader.go`
- `apps/server_core/internal/modules/profitability/adapters/orders/order_reader.go`
- `apps/server_core/internal/modules/profitability/adapters/orders/order_reader_test.go`

## Task 2 Surface

- `apps/server_core/internal/modules/profitability/domain/input.go`
- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`

## Task 3 Surface

- `apps/server_core/migrations/0031_profitability_order_realization.sql`
- `Store.ReplaceSnapshots` and `Store.ListSnapshots` in `apps/server_core/internal/modules/profitability/adapters/postgres/store.go`
- `apps/server_core/internal/modules/profitability/adapters/postgres/profit_snapshot_integration_test.go`
- OpenAPI realization/quality/flag/snapshot schemas around lines 2666-2705 in `contracts/api/marketplace-central.openapi.yaml`
- SDK realization/quality/flag/snapshot types around lines 380-405 in `packages/sdk-runtime/src/index.ts`
- SDK profitability fixtures/assertions around lines 1100-1190 in `packages/sdk-runtime/src/index.test.ts`

## Task 4 Surface

- `packages/feature-orders/src/OrdersPage.tsx`
- `packages/feature-orders/src/OrdersPage.test.tsx`
- Product context: `packages/feature-orders/PRODUCT.md`

## Completed Task Gates

- Task 1: independent SPEC compliant + QUALITY approved, no findings.
- Task 2: independent SPEC compliant + QUALITY approved, no findings.
- Task 3: independent SPEC compliant + QUALITY approved; one report-count Minor corrected; real PostgreSQL test still pending.
- Task 4: initial assertion findings fixed; independent re-review SPEC compliant + QUALITY approved, no findings; browser QA pending.

## Validation Boundary

- Deterministic Go/SDK/UI/build evidence is recorded in per-task reports.
- Do not infer real integration readiness from mocks, skipped tests, compilation, or build.
- PostgreSQL 16 migration/round-trip, live Mercado Livre paid/cancelled orders, real Oracle cost/tax, and built-in-browser desktop/mobile QA remain controller gates after this review.
