# M-06 F-03 Final Review Fix Report

## Status

`DONE_WITH_CONCERNS`

Both whole-feature findings are fixed with focused RED/GREEN evidence. The controller must still run the complete feature-orders/web/build verification because this worker's required full npm run was rejected by the approval system after the focused Vitest runs succeeded.

## Root Cause And Hypotheses

- Snapshot truncation: `CalculateSnapshots` sliced the map-derived snapshot set before deriving replacement order IDs and before `ReplaceSnapshots`. Hypothesis confirmed by a focused regression: `limit=1` produced calculated/persisted/returned counts `1/1/1` instead of `4/4/1`, so response pagination controlled canonical persistence and map iteration controlled which row survived.
- Actor provenance: `OrdersPage` stored a default editable operator name and synthesized `actor_id` from that display text. Hypothesis confirmed by focused UI regressions: no absent-identity explanation/disable existed, and a supplied `{operator,user-42,Trusted Operator}` was ignored in favor of `{operator,leandro,Leandro}`.
- No architecture ambiguity remained: the existing snapshot store already accepts a complete atomic replacement set, and the current app safely permits `OrdersPage` to omit an operator until authenticated context exists.

## RED Evidence

- Go: `go test ./internal/modules/profitability/application -run TestCalculateSnapshotsLimitsOnlyResponseAndReplacesCompleteDeterministicSet -count=1` failed with `calculated:1 persisted:1 returned:1, want 4, 4, 1`.
- UI: focused Vitest ran 2 tests and both failed for the expected semantics: absent operator had no explanation/disabled submit, and supplied immutable actor was replaced by the editable/default Leandro identity.
- An initial sandboxed UI attempt failed before Vitest because root config access was denied; it is not counted as RED evidence.

## Fixes

- Complete bounded snapshots are sorted by provider order ID, scope (`order` before `item`), provider item ID, then provider variation ID.
- Replacement order IDs are deterministically derived from the complete sorted set; `ReplaceSnapshots` receives all snapshots before response limiting.
- `CalculatedCount` reports the complete persisted count; only returned `Items` are sliced by `limit`.
- Regression seeds a stale realized cancelled order, uses a low response limit, and proves all four snapshots/two order IDs are replaced and the cancelled order snapshot becomes `not_realized`.
- `OrdersPageProps` now accepts immutable optional `operator: Readonly<ProfitabilityActor>`; exact actor fields are forwarded without ID derivation.
- Editable/default operator state and input are removed. Omitted or blank required actor identity shows explicit text and disables submission. `AppRouter` may continue omitting the prop safely.
- Backend, SDK, OpenAPI, migrations, actor invariants, and provider behavior are unchanged.

## GREEN And Verification Evidence

- Focused Go regression: passed after the stale-row seed was added.
- Required Go packages: `go test ./internal/modules/profitability/application ./internal/modules/profitability/adapters/postgres -count=1` passed both packages.
- Focused actor UI: 2/2 tests passed before the additional blank-identity coverage was added; controller rerun of the complete file remains required.
- `gofmt -d` on changed Go files was empty.
- Profitability orders-import boundary `rg` returned exit `1` with empty output, proving no forbidden orders-module imports.
- `git diff --check` on scoped paths returned no whitespace errors; the only message was Git's pre-existing LF-to-CRLF working-copy warning.

## Database Evidence Level

PostgreSQL integration evidence is compile/skip only, not real database validation. Verbose repository tests explicitly skipped `TestManualAdjustmentsAppendOnlyReadbackAndConstraints` and `TestProfitSnapshotRealizationPersistence` because `MC_DATABASE_URL` is not set. Apply required migrations to a real PostgreSQL 16 target and rerun with `MC_DATABASE_URL` for integration evidence.

## Remaining Verification And Warnings

- Worker could not rerun the complete feature-orders test, exact web tests, or production web build: escalation was rejected due approval/usage limit after a sandbox run had already shown root Vitest config access denial. Controller owns these exact commands before acceptance.
- Preserve and report known dependency-level warnings if reproduced: Node `DEP0205 module.register()` and React Router/Lucide `use client` build warnings. They are not fixed or hidden here.
- No stage, commit, reset, revert, clean, contract change, or unrelated-file edit was performed.

## Exact Changed Paths

- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`
- `packages/feature-orders/src/OrdersPage.tsx`
- `packages/feature-orders/src/OrdersPage.test.tsx`
- `.superpowers/sdd/m06-f03-final-review-fix-report.md`
- `.superpowers/sdd/m06-f03-correction-report.md`

## Self-Review

The fix is confined to the two root causes. Snapshot monetary pointers and realization calculation are untouched; sorting only establishes persistence/response order. Actor validation checks only required type/ID presence, preserves optional name, and forwards supplied values exactly. Existing amount, reason, scope, category, item controls, adjustment listing, and refresh flow remain intact.
