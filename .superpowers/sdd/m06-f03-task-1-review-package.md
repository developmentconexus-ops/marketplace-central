# M-06 F-03 Task 1 Review Package

## Baseline

- Branch: `main`
- HEAD before and after Task 1: `8dba7db docs(handoff): pause M-06 at F-03 RED`
- No commit, staging, reset, revert, or clean operation was performed.
- The profitability module is untracked in the shared worktree, so Git has no tracked pre-Task-1 blob from which to produce a conventional hunk diff. For this task, each scoped code file below is a full-file change snapshot and is the review surface.

## Scoped Code Files (read each once in full)

1. `apps/server_core/internal/modules/profitability/domain/order_fact.go`
2. `apps/server_core/internal/modules/profitability/ports/order_reader.go`
3. `apps/server_core/internal/modules/profitability/adapters/orders/order_reader.go`
4. `apps/server_core/internal/modules/profitability/adapters/orders/order_reader_test.go`
5. `apps/server_core/internal/modules/profitability/application/service.go`
6. `apps/server_core/internal/modules/profitability/application/service_test.go`

## Evidence Files

- `.superpowers/sdd/m06-f03-task-1-report.md`
- `.superpowers/sdd/m06-f03-correction-report.md`

## Scoped Status

All eight Task 1 code/evidence paths are untracked (`??`) in the existing heavily dirty worktree. No unrelated path is part of this review package.

## Controller Checks

- Implementer report contains expected RED compile failures before production edits.
- Fresh scoped GREEN reports both required Go packages passing.
- `gofmt -d` was empty after formatting.
- Boundary `rg` over profitability domain, ports, and application returned no matches.
- `git diff --check` reported no whitespace errors for tracked hunks; full untracked files are reviewed directly.
