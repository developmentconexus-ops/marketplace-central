# M-06 F-03 Task 2 Review Package

## Baseline

- Branch: `main`
- HEAD unchanged at `8dba7db`; Task 2 created no commit or staging changes.
- The profitability module is untracked in the shared dirty worktree, so the three scoped files below are full-file review snapshots rather than conventional Git hunks.

## Scoped Code Files (read each once in full)

1. `apps/server_core/internal/modules/profitability/domain/input.go`
2. `apps/server_core/internal/modules/profitability/application/service.go`
3. `apps/server_core/internal/modules/profitability/application/service_test.go`

## Evidence

- Requirements: `.superpowers/sdd/m06-f03-task-2-brief.md`
- Implementer report: `.superpowers/sdd/m06-f03-task-2-report.md`
- Aggregate correction evidence: `.superpowers/sdd/m06-f03-correction-report.md`

## Controller Checks

- Focused RED failed on the required missing realization field/quality/flags.
- Focused GREEN and full profitability/composition regression passed.
- `gofmt -d` was empty and boundary `rg` returned no forbidden orders imports.
- Reviewer must verify exact realized/not-realized/unknown semantics and the absence of Task 3-4 scope.
