# F-10 Plan — Postgres Harness Migration Inventory

## Execution

1. Compile and validate `context.json` from this plan, the dispatch, and the
   accepted base SHA.
2. Add a fail-closed canonical SQL inventory resolver to `Postgres.psm1`, carry
   the resolved count on the run spec, and compare lifecycle output to it.
3. Make the fake probe count SQL files from its Go-process working directory;
   make fake and Docker assertions consume the derived run-spec count.
4. Add focused missing/empty inventory cases and run `postgres-harness-fake`.
5. Run the registered `go-orders-linkage-postgres` integration lifecycle,
   verify F-09 repository test execution and cleanup, then run
   `git-diff-check`.
6. Record `validation.md`, inspect scoped status/diff, and create one intentional
   commit.

## Fixed decisions

- Owner/seam: this worker exclusively owns the dispatched `harness-control`
  paths for F-10.
- Interface: `New-HarnessPostgresRunSpec` remains the repository-root boundary
  and adds `ExpectedMigrationCount`; lifecycle callers do not accept an
  override.
- Consumers: lifecycle gate and fake/Docker assertions use that property; the
  fake probe independently reads canonical files from `apps/server_core` CWD.
- Legacy: no 32/33 fallback or compatibility branch remains.
- Unknown: absent, unreadable, or empty canonical inventory is explicitly
  `HPG_MIGRATION_INVENTORY_INVALID` and stops before external side effects.
- Allowed paths: Feature artifacts plus the four dispatched harness/test files.
- Forbidden paths: application, migration, Oracle/provider, production, API,
  SDK, governance, dependencies, and `.agents`.

## Commands and proof

- Context: `pwsh -NoProfile -File scripts/harness.ps1 -Command context-compile
  -FeaturePath <feature-dir> -BaseSha
  a4b54e935cb525bb8b0470ad1bb1fe53065cde59 -AllowedPath <each dispatched
  path>`; save the generated values-safe pack as `context.json`, then run
  `context-validate -ContextPath <context.json> -RequireCurrentBase`.
- `postgres-harness-fake`: `pwsh -NoProfile -File
  scripts/tests/postgres-lifecycle.tests.ps1`.
- `go-orders-linkage-postgres`: `pwsh -NoProfile -File scripts/harness.ps1
  -Command integration`.
- `git-diff-check`: `git diff --check` plus scoped status and diff inspection.

## Stop conditions

Stop and return exact safe tokens if work requires a forbidden path, the
registered integration fails outside harness-count/F-09 linkage scope, owned
Docker cleanup cannot be proven, or any Oracle/provider/production/secret/PII
operation would occur.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-10",
  "required_sources": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-contract.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-09-order-line-identity-linkage-ledger/validation.md",
    "contracts/governance/execution-lanes.json",
    "contracts/governance/shared-seams.json"
  ],
  "knowledge_route_ids": ["harness-control-plane", "orders-margin"],
  "allowed_paths": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-10-postgres-harness-migration-inventory/**",
    "scripts/harness/Postgres.psm1",
    "scripts/tests/postgres-lifecycle.tests.ps1",
    "scripts/tests/postgres-docker.integration.tests.ps1",
    "scripts/tests/fixtures/harness/postgres-docker-probe.mjs"
  ],
  "forbidden_paths": [
    "apps/**",
    "contracts/api/**",
    "packages/**",
    "docs/research/**",
    "contracts/governance/**",
    ".agents/**"
  ],
  "side_effects": {
    "allowed": ["repository-write", "database-write", "isolated-cache-write"],
    "forbidden": ["external-network", "provider-write"]
  },
  "commands": [
    {"id": "postgres-harness-fake", "command_id": "postgres-harness-fake", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "go-orders-linkage-postgres", "command_id": "go-orders-linkage-postgres", "lane_id": "integration", "expected_exit_code": 0},
    {"id": "git-diff-check", "command_id": "git-diff-check", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F10-AC01", "milestone_criterion_id": "M-06-C01", "command_ids": ["postgres-harness-fake", "go-orders-linkage-postgres"]},
    {"id": "F10-AC02", "milestone_criterion_id": "M-06-C01", "command_ids": ["postgres-harness-fake"]},
    {"id": "F10-AC03", "milestone_criterion_id": "M-06-C01", "command_ids": ["postgres-harness-fake"]},
    {"id": "F10-AC04", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-orders-linkage-postgres", "git-diff-check"]}
  ],
  "stop_conditions": [
    {"code": "forbidden-path", "condition": "The change requires application, migration, provider, Oracle, production, API, SDK, governance, dependency, or agent files."},
    {"code": "unrelated-integration-failure", "condition": "The registered integration fails outside harness-count or F-09 linkage scope."},
    {"code": "cleanup-unproved", "condition": "Owned Docker resource cleanup cannot be proven by the registered lifecycle."},
    {"code": "forbidden-side-effect", "condition": "Any command would access Oracle, a provider, production data, secrets, or PII."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "commit", "changed-paths", "evidence", "blockers", "next"]
}
```
