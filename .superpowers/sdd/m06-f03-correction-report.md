# M-06 F-03 Correction Report

## Status

`IN_PROGRESS`

## Design And Plan

- `docs/superpowers/specs/2026-07-09-m06-f03-order-realization-design.md`
- `docs/superpowers/plans/2026-07-09-m06-f03-order-realization.md`

## Task Evidence

Pending Task 1 TDD cycle.

## Task 1 - Profitability-Owned Order Facts

- RED: focused `TestOrderReader|TestImportMarginInputs` run failed as expected because the profitability-owned order fact types and constants did not exist.
- GREEN: `go test ./internal/modules/profitability/adapters/orders ./internal/modules/profitability/application -count=1` passed both packages with `GOCACHE=.gocache`.
- Boundary: exact `rg` over profitability domain, ports, and application returned no orders-module imports (exit `1`, empty output).
- Paths: `profitability/domain/order_fact.go`, `profitability/ports/order_reader.go`, `profitability/adapters/orders/order_reader.go`, `profitability/adapters/orders/order_reader_test.go`, `profitability/application/service.go`, and `profitability/application/service_test.go`.
- No commit was created.

## Task 2 - Realization-Aware Snapshot Calculation

- RED: focused `TestCalculateSnapshots` failed as expected because `ProfitSnapshot.RealizationState`, `not_realized`, `order_cancelled`, and `order_state_unknown` did not yet exist.
- GREEN: focused application test passed; fresh `go test ./internal/modules/profitability/... ./internal/composition -count=1` passed all profitability packages and composition with `GOCACHE=.gocache`.
- Semantics: realized behavior is retained; cancelled and unknown snapshots preserve known inputs but always have nil contribution/margin; missing order IDs normalize to unknown; nullable monetary facts remain nullable.
- Boundary/format: exact orders-import `rg` over profitability domain/ports/application returned no matches (exit `1`), and final `gofmt -d` was empty (exit `0`).
- Paths: `profitability/domain/input.go`, `profitability/application/service.go`, `profitability/application/service_test.go`, and `.superpowers/sdd/m06-f03-task-2-report.md`.
- Evidence is deterministic unit/contract evidence, not live provider/database validation. No commit was created.

## Task 3 - PostgreSQL, OpenAPI, And SDK Realization Contract

- RED: the database-gated repository test compiled and skipped because `MC_DATABASE_URL` was absent; SDK tests then failed on the two expected missing `realization_state` fixture assertions (`33 passed, 2 failed`).
- GREEN: `go test ./internal/modules/profitability/... ./internal/composition -count=1` passed with the PostgreSQL integration test still explicitly skipped; SDK passed `35/35` tests.
- Contract: migration `0031`, repository INSERT/SELECT/Scan, OpenAPI, and SDK now share `realized | not_realized | unknown`; quality adds `not_realized`; flags add `order_cancelled` and `order_state_unknown`; cancelled/unknown monetary nils remain nil.
- Database evidence level: compile/skip only, not real PostgreSQL validation. Controller must apply migrations `0028`-`0031` to PostgreSQL 16 and rerun the focused test with `MC_DATABASE_URL`.
- Paths and complete evidence: `.superpowers/sdd/m06-f03-task-3-report.md`. No commit was created.

## Task 4 - Orders UI Realization Semantics

- RED: exact focused feature-orders run failed semantically with `3` failed and `6` passed tests because `Data quality`, `not_realized` filter/summary/labels, operational state text, and em-dash null rendering were absent.
- GREEN: exact feature-orders run passed `9/9`; exact web regression passed `8/8`; production web build transformed `1783` modules and passed in `2m 3s`.
- Semantics: cancelled is `Not realized` / `Order not realized` with `Order cancelled`; unknown remains `Incomplete` with `Order state unknown`; operational flags are separated from missing-input `Data quality`; cancelled/unknown contribution and margin render `—`; cancelled does not show negative margin.
- Contract/scope: local manual-adjustment actor is required; no React business math and no backend/OpenAPI/SDK/migration changes; `ClientContext.tsx` required no Task 4 edit.
- Evidence level: unit/regression/build only. Built-in-browser desktop/mobile QA remains pending; dependency-level Node/Vite warnings were observed.
- Paths and complete evidence: `.superpowers/sdd/m06-f03-task-4-report.md`. No commit was created.

## Final Whole-Feature Review Correction

- RED: low-limit snapshot calculation persisted/count-reported only `1/4` snapshots; actor UI lacked a disabled absent state and replaced supplied provenance with synthesized `Leandro` identity.
- GREEN: complete snapshots are deterministically sorted and atomically replaced before response limiting; calculated count is complete; stale realized cancelled data is replaced by `not_realized`. Immutable optional operator provenance is forwarded exactly, while omitted/blank required identity disables submission with explicit text.
- Go application/PostgreSQL packages passed; focused actor UI passed. PostgreSQL integration tests skipped explicitly because `MC_DATABASE_URL` is absent, so no real DB claim is made.
- Full feature-orders/web/build reruns are controller-pending because this worker's approval escalation was rejected after focused Vitest verification. Known dependency warnings must remain visible if reproduced.
- Boundary `rg` was empty, `gofmt -d` was empty, scoped `git diff --check` had no whitespace errors, and no commit was created.
- Complete evidence: `.superpowers/sdd/m06-f03-final-review-fix-report.md`.
