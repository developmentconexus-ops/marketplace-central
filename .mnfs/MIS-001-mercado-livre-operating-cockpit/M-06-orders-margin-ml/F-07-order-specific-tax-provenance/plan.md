# F-07 order-specific tax provenance plan

```yaml
id: F-07
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-07
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-1
lifecycle_scope: feature
split_decision: single
split_reason: Three ordered implementation steps across one bounded internal-read/profitability seam with no unresolved design choice.
```

## Feature ID

F-07-order-specific-tax-provenance

## Steps

1. Compile and validate the feature context against the accepted base SHA and
   exclusive path scope.
2. Add the exact Oracle tax source identity to the internal-read domain/port,
   fake adapter, and Oracle adapter; prove exact predicates, returned
   provenance, missing identity, missing rows, and partial tax.
3. Carry optional verified source identity through profitability; do not query
   tax without it, preserve nil/missing inputs and incomplete snapshots, and
   run the focused proof plus diff-scope check.

## Files Expected To Change

- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-07-order-specific-tax-provenance/**`
- `apps/server_core/internal/modules/internal_read/domain/internal_tax.go`
- `apps/server_core/internal/modules/internal_read/ports/reader.go`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/reader_test.go`
- `apps/server_core/internal/modules/internal_read/adapters/fake/reader.go`
- `apps/server_core/internal/modules/internal_read/adapters/fake/reader_test.go`
- `apps/server_core/internal/modules/profitability/ports/internal_read.go`
- `apps/server_core/internal/modules/profitability/adapters/internalread/fact_reader.go`
- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`

## Verification Commands

- Command ID `go-profitability-internal-read`: set `GOCACHE=.gocache`; run
  `go test ./internal/modules/profitability/...` from `apps/server_core`.
  Satisfies: `F07-AC02`, `F07-AC03`.
- Command ID `go-tax-provenance-focused`: set `GOCACHE=.gocache`; run
  `go test ./internal/modules/internal_read/adapters/oracle ./internal/modules/internal_read/adapters/fake` from `apps/server_core`.
  Satisfies: `F07-AC01`.
- Command ID `git-diff-check`: run `git diff --check` and inspect changed paths.
  Satisfies: `F07-AC01`, `F07-AC02`, `F07-AC03` scope/format proof.

## QA Steps

- Inspect the Oracle SQL to confirm exact document/line/product/incidence
  predicates and absence of product/date aggregation.
- Inspect the profitability call boundary to confirm a resolved product alone
  cannot produce Oracle tax and that missing provenance remains explicit.

## Rollback/Risk Notes

The contract is internal and carries optional provenance. Reverting the scoped
commit restores the prior behavior. The main operational risk is that upstream
orders do not yet supply the Oracle identity; this intentionally yields
incomplete tax rather than a false margin.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-07",
  "required_sources": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-contract.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/research/mnos-sankhya-read-interface-contract.md"
  ],
  "knowledge_route_ids": ["portfolio-core", "orders-margin"],
  "allowed_paths": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-07-order-specific-tax-provenance/**",
    "apps/server_core/internal/modules/internal_read/domain/internal_tax.go",
    "apps/server_core/internal/modules/internal_read/ports/reader.go",
    "apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go",
    "apps/server_core/internal/modules/internal_read/adapters/oracle/reader_test.go",
    "apps/server_core/internal/modules/internal_read/adapters/fake/reader.go",
    "apps/server_core/internal/modules/internal_read/adapters/fake/reader_test.go",
    "apps/server_core/internal/modules/profitability/ports/internal_read.go",
    "apps/server_core/internal/modules/profitability/adapters/internalread/fact_reader.go",
    "apps/server_core/internal/modules/profitability/application/service.go",
    "apps/server_core/internal/modules/profitability/application/service_test.go"
  ],
  "forbidden_paths": ["docs/research/**", "packages/sdk-runtime/**", "apps/server_core/openapi/**"],
  "side_effects": {"allowed": ["repository-write", "isolated-cache-write"], "forbidden": ["database-mutation", "external-network", "provider-write"]},
  "commands": [
    {"id": "go-profitability-internal-read", "command_id": "go-profitability-internal-read", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "go-tax-provenance-focused", "command_id": "go-tax-provenance-focused", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "git-diff-check", "command_id": "git-diff-check", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F07-AC01", "milestone_criterion_id": "M-06-C02", "command_ids": ["go-tax-provenance-focused", "git-diff-check"]},
    {"id": "F07-AC02", "milestone_criterion_id": "M-06-C02", "command_ids": ["go-profitability-internal-read", "git-diff-check"]},
    {"id": "F07-AC03", "milestone_criterion_id": "M-06-C02", "command_ids": ["go-profitability-internal-read", "git-diff-check"]}
  ],
  "stop_conditions": [
    {"code": "contract-conflict", "condition": "Exact Oracle identity requires an owner-approved mapping not present in repository truth."},
    {"code": "scope-conflict", "condition": "Implementation requires a change outside the allowed paths."},
    {"code": "validation-blocked", "condition": "Focused Go proof cannot execute."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "commit", "changed-paths", "evidence", "blockers", "next"]
}
```

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Compile/validate context, implement, and record validation.
- Required files/evidence: spec, context pack, focused Go tests, diff check.
- Blockers or open decisions: None; absent real linkage is represented as
  missing and remains an operational handoff gap.
