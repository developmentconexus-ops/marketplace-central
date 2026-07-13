# F-13 Assisted Linkage API and Runtime Wiring Plan

```yaml
id: F-13
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-13
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-0
lifecycle_scope: feature
split_decision: single
```

## Feature ID

F-13-assisted-linkage-api-runtime

## Split Decision

`single`: although the expected path count exceeds the split heuristic, the
dispatch requires one plan/execution session and the runtime, HTTP, OpenAPI,
SDK, and governance changes form one atomic shared seam.

## Owner, Ports, Consumers, and Legacy Decision

- Owner: orders application owns assisted-linkage behavior; composition owns
  runtime selection; transport owns safe DTO/error mapping.
- Ports: extend the existing F-12 ledger-facing service contract only as needed
  for `GetCurrent`; retain the F-11 bridge as the only orders-to-internal-read
  adapter.
- Consumers: orders HTTP routes and `sdk-runtime` callers.
- Legacy decision: keep all existing endpoints and unrelated Oracle composition
  unchanged. Invalid assisted-linkage configuration disables only these routes.
- Explicit unknowns: nullable facts remain nil/null; invalid configuration is
  unavailable; lineage conflict remains conflict; no zero/default substitute.

## Steps

1. Compile/validate the bounded context after this plan and read only emitted
   selectors; record any known route-budget failure and use the smallest
   truthful reduced context pack if necessary.
2. Inspect the exact F-11/F-12/F-09 and composition/transport/API/SDK seams named
   by context; define strict runtime parsing with explicit bounded settings and
   governance entries.
3. Add `GetCurrent` and server-owned evidence reference propagation only where
   needed, with focused application/adapter tests.
4. Add fail-closed composition plus strict current/candidates/confirm handlers,
   DTO conversion, stable safe errors, and transport/config/composition tests.
5. Update OpenAPI and SDK runtime atomically, including parity/transport tests
   that prove forbidden confirm controls are absent.
6. Run all registered quick-validation commands, inspect the scoped diff, write
   `validation.md`, stage only owned paths, and create one intentional commit.

## Files Expected To Change

- Feature `spec.md`, `plan.md`, `context.json`, `validation.md`: durable scope,
  context, evidence, and handoff.
- `apps/server_core/internal/modules/internal_read/adapters/oracle/sankhya_linkage_config.go`
  and test: explicit runtime configuration parsing/validation if adapter-owned.
- Orders domain/ports/application/internalread files and tests listed by the
  dispatch: current read and evidence/revision propagation.
- `apps/server_core/internal/modules/orders/transport/sankhya_linkage_handler.go`
  and test: strict routes, DTOs, and error mapping.
- `apps/server_core/internal/composition/root.go`: fail-closed runtime wiring.
- `contracts/api/marketplace-central.openapi.yaml`: public contract.
- `packages/sdk-runtime/src/index.ts` and test: matching SDK transport/runtime.
- `contracts/governance/runtime-config.json`: all new environment settings.

No other path is authorized.

## Verification Commands

- Command ID: `go-assisted-linkage-http`
  - Satisfies: F13-AC01, F13-AC02, F13-AC03, F13-AC05
  - Expected: registered focused Go tests/builds exit 0 using
    `GOCACHE=.gocache`; invalid config is unavailable, exact scope and strict
    requests are enforced, and no live adapter is invoked.
- Command ID: `sdk-assisted-linkage`
  - Satisfies: F13-AC02, F13-AC04, F13-AC05
  - Expected: registered SDK tests/build/parity checks exit 0; confirm exposes
    no caller-controlled server facts and nullable/lineage DTOs match OpenAPI.
- Command ID: `governance-contracts`
  - Satisfies: F13-AC01, F13-AC05
  - Expected: runtime registry/schema/drift checks exit 0.
- Command ID: `git-diff-check`
  - Satisfies: F13-AC04, F13-AC05
  - Expected: `git diff --check` exits 0 and status/diff contains only owned
    paths with no forbidden or live-system artifacts.

## QA Steps

- Inspect OpenAPI confirm request and SDK confirm input together.
  - Expected: neither contains tenant, event ID, configuration revision,
    evidence reference, actor type, external key, or unknown-field passthrough.
- Inspect runtime parser and composition unavailable path.
  - Expected: every field is explicit, fixed TOP values are not environment
    defaults, startup continues, and only assisted-linkage endpoints are 503.
- Inspect response DTO conversion.
  - Expected: only generic identifiers/safe display/audit facts are present;
    unknown values remain nullable and no Oracle field/SQL/provider/PII leaks.

## Rollback/Risk Notes

- API/SDK drift is controlled by a single commit and parity tests; revert the
  whole Feature commit if contract integration fails.
- Composition changes are isolated behind an unavailable implementation so a
  malformed deployment does not affect unrelated routes.
- No migration or external write exists to roll back; confirmed ledger history
  remains append-only and is never deleted.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-13",
  "required_sources": [
    "AGENTS.md",
    "ARCHITECTURE.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-contract.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/sankhya-admin-spec.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-11-sankhya-linkage-read-contract/spec.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-12-assisted-linkage-service/spec.md",
    "contracts/governance/knowledge-routes.json",
    "contracts/governance/runtime-config.json",
    "contracts/governance/shared-seams.json",
    "wiki/modules/orders.md"
  ],
  "knowledge_route_ids": ["portfolio-core", "orders-margin"],
  "allowed_paths": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-13-assisted-linkage-api-runtime/**",
    "apps/server_core/internal/modules/internal_read/adapters/oracle/sankhya_linkage_config.go",
    "apps/server_core/internal/modules/internal_read/adapters/oracle/sankhya_linkage_config_test.go",
    "apps/server_core/internal/modules/orders/domain/assisted_sankhya_linkage.go",
    "apps/server_core/internal/modules/orders/ports/sankhya_linkage_reader.go",
    "apps/server_core/internal/modules/orders/application/assisted_sankhya_linkage_service.go",
    "apps/server_core/internal/modules/orders/application/assisted_sankhya_linkage_service_test.go",
    "apps/server_core/internal/modules/orders/adapters/internalread/sankhya_linkage_reader.go",
    "apps/server_core/internal/modules/orders/adapters/internalread/sankhya_linkage_reader_test.go",
    "apps/server_core/internal/modules/orders/transport/sankhya_linkage_handler.go",
    "apps/server_core/internal/modules/orders/transport/sankhya_linkage_handler_test.go",
    "apps/server_core/internal/composition/root.go",
    "contracts/api/marketplace-central.openapi.yaml",
    "packages/sdk-runtime/src/index.ts",
    "packages/sdk-runtime/src/index.test.ts",
    "contracts/governance/runtime-config.json"
  ],
  "forbidden_paths": [
    "apps/server_core/migrations/**",
    "apps/server_core/internal/modules/profitability/**",
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
    {"id": "go-assisted-linkage-http", "command_id": "go-assisted-linkage-http", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "sdk-assisted-linkage", "command_id": "sdk-assisted-linkage", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "governance-contracts", "command_id": "governance-contracts", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "git-diff-check", "command_id": "git-diff-check", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F13-AC01", "milestone_criterion_id": "M-06-C02", "command_ids": ["go-assisted-linkage-http", "governance-contracts"]},
    {"id": "F13-AC02", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-assisted-linkage-http", "sdk-assisted-linkage"]},
    {"id": "F13-AC03", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-assisted-linkage-http"]},
    {"id": "F13-AC04", "milestone_criterion_id": "M-06-C02", "command_ids": ["sdk-assisted-linkage", "git-diff-check"]},
    {"id": "F13-AC05", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-assisted-linkage-http", "sdk-assisted-linkage", "governance-contracts", "git-diff-check"]}
  ],
  "stop_conditions": [
    {"code": "invented-runtime-default", "condition": "A runtime field, schema, revision, attestation, evidence reference, or limit default would be invented."},
    {"code": "caller-audit-control", "condition": "Caller-controlled tenant, event, configuration, evidence, actor-type, or external-key data would reach the ledger."},
    {"code": "api-sdk-drift", "condition": "OpenAPI and sdk-runtime cannot change atomically."},
    {"code": "authentication-claim", "condition": "Production authentication or manual-adjustment authorization would be claimed."},
    {"code": "forbidden-scope", "condition": "A migration, profitability, UI, Docker, live database/provider, dependency, secret, or PII operation would be required."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "commit", "changed-paths", "evidence", "blockers", "next"]
}
```

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Compile/validate context, implement the six steps, run all four
  command IDs, write `validation.md`, and commit once.
- Required files/evidence: spec, plan, dispatch, context, focused command
  outputs, scoped diff
- Blockers or open decisions: None.
