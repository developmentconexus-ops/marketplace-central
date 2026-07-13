# F-15 Partial Tax Snapshot Quality Plan

## Fixed scope

- Owner: bounded F-15 Feature Implementer.
- Application seam: profitability snapshot input application and rollup.
- Ports/interfaces: unchanged; existing tax-lineage and exact-tax inputs only.
- Consumers: existing `ImportOrder` item and order snapshot construction.
- Legacy decision: replace the amount-only completeness assumption in place;
  add no fallback, estimate, or compatibility branch.
- Explicit unknowns: nil amounts remain nil; non-complete qualities remain
  incomplete even when amounts are known; no operational fact is inferred.

## Paths

Implementation and proof are limited to:

- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`
- this F-15 Feature directory
- F-14 `validation.md` for the evidence correction
- M-06 `corrections/correction-task.md` for the append-only attempt-2 entry

The parent-owned modified milestone review and all review/QA, research, probe,
output, auth/manual-adjustment, lineage, Oracle, API/SDK, migration, UI, live
system, and dependency paths remain untouched.

## Execution

1. Compile and validate the scoped `context.json` against frozen SHA
   `ef4b08c78d30a5e2269e79b051a432c9dc12b58d`.
2. Inspect only the compiled selectors in the two profitability application
   files and the two authorized evidence artifacts.
3. Make the smallest application change so a tax component is marked missing
   when its amount is nil or its quality is not complete, without dropping a
   known amount.
4. Add import-to-snapshot regressions for partial lineage/all-known components,
   complete lineage/incomplete source quality, and fully complete control.
5. Amend F-14 validation truthfully and append correction attempt 2 while
   retaining the explicit C03 failure/deferment.

## Registered proof

- `go-profitability-partial-tax`: from `apps/server_core`, with repository-local
  `GOCACHE=.gocache`, run the focused profitability application regression.
- `go-repository`: from `apps/server_core`, with repository-local
  `GOCACHE=.gocache`, run `go test -count=1 ./...`.
- `git-diff-check`: run working-tree and staged `git diff --check` checks and
  verify only allowed paths are staged.

Record exact commands/results in `validation.md`, validate current context,
create exactly one intentional commit, and return the compact handoff.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-15",
  "required_sources": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-14-profitability-sankhya-lineage/spec.md"
  ],
  "allowed_paths": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-15-partial-tax-snapshot-quality/**",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-14-profitability-sankhya-lineage/validation.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/corrections/correction-task.md",
    "apps/server_core/internal/modules/profitability/application/service.go",
    "apps/server_core/internal/modules/profitability/application/service_test.go"
  ],
  "forbidden_paths": [
    "apps/server_core/migrations/**",
    "apps/server_core/internal/modules/orders/**",
    "apps/server_core/internal/modules/internal_read/**",
    "contracts/**",
    "packages/**",
    "apps/web/**",
    "scripts/**",
    "docs/research/**",
    "docker/**"
  ],
  "side_effects": {
    "allowed": ["repository-write", "isolated-cache-write"],
    "forbidden": ["database-mutation", "external-network", "provider-write"]
  },
  "commands": [
    {"id": "go-profitability-partial-tax", "command_id": "go-profitability-partial-tax", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "go-repository", "command_id": "go-repository", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "git-diff-check", "command_id": "git-diff-check", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F15-AC01", "milestone_criterion_id": "M-06-C02", "command_ids": ["go-profitability-partial-tax"]},
    {"id": "F15-AC02", "milestone_criterion_id": "M-06-C02", "command_ids": ["go-profitability-partial-tax"]},
    {"id": "F15-AC03", "milestone_criterion_id": "M-06-C02", "command_ids": ["go-profitability-partial-tax", "go-repository"]},
    {"id": "F15-AC04", "milestone_criterion_id": "M-06-C02", "command_ids": ["git-diff-check"]}
  ],
  "stop_conditions": [
    {"code": "discard-known-tax", "condition": "A known partial tax amount would be discarded instead of retained."},
    {"code": "incomplete-margin", "condition": "A non-complete tax quality could still produce contribution or margin."},
    {"code": "c03-scope", "condition": "The correction would expand into C03 authentication, authorization, or manual adjustments."},
    {"code": "forbidden-scope", "condition": "A lineage, Oracle, API, SDK, migration, UI, live-system, dependency, secret, or PII change would be required."}
  ],
  "retry_budget": {"max_correction_attempts": 2},
  "handoff_fields": ["status", "commit", "changed-paths", "evidence", "blockers", "next"]
}
```
