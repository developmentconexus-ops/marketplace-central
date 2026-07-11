# F-03 Ephemeral PostgreSQL and Canonical Migrations — Plan

```yaml
id: F-03
type: feature-plan
status: planned
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Scope and Ownership

F-03 is an L3 change over the accepted harness runtime, migration source,
database-test boundary, and migration shared seam. One serial writer owns these
paths until acceptance. The dev Compose PostgreSQL remains an observer target
only during final validation.

## Machine Work Contract

```json
{
  "schema_version":"1.0",
  "feature_id":"F-03",
  "required_sources":[
    "scripts/tests/hermetic-lanes.tests.ps1",
    "scripts/tests/harness-aliases.tests.ps1",
    "contracts/governance/execution-lanes.json",
    "contracts/governance/runtime-config.json",
    "scripts/harness.ps1",
    "apps/server_core/internal/platform/migrate/runner.go",
    "apps/server_core/tests/integration/migrate_runner_test.go"
  ],
  "allowed_paths":[
    "scripts/harness.ps1",
    "scripts/harness/Postgres.psm1",
    "scripts/harness/Environment.psm1",
    "scripts/tests/postgres-contract.tests.ps1",
    "scripts/tests/postgres-lifecycle.tests.ps1",
    "scripts/tests/postgres-go.tests.ps1",
    "scripts/tests/postgres-docker.integration.tests.ps1",
    "scripts/tests/postgres-dev-invariance.integration.tests.ps1",
    "scripts/tests/f03-regression.tests.ps1",
    "scripts/tests/hermetic-lanes.tests.ps1",
    "scripts/tests/harness-aliases.tests.ps1",
    "scripts/tests/fixtures/harness/postgres-docker-probe.mjs",
    "apps/server_core/migrations/source.go",
    "apps/server_core/internal/platform/migrate/runner.go",
    "apps/server_core/internal/platform/migrate/runner_test.go",
    "apps/server_core/internal/testsupport/postgres/**",
    "apps/server_core/cmd/migrate/main.go",
    "apps/server_core/cmd/migrate/main_test.go",
    "apps/server_core/cmd/testdb/**",
    "apps/server_core/tests/integration/phase1_smoke_test.go",
    "apps/server_core/tests/integration/migrate_runner_test.go",
    "apps/server_core/tests/integration/marketplaces_repository_test.go",
    "apps/server_core/tests/integration/integrations_credential_repo_test.go",
    "apps/server_core/internal/modules/orders/adapters/postgres/order_repo_test.go",
    "apps/server_core/internal/modules/profitability/adapters/postgres/profit_snapshot_integration_test.go",
    "apps/server_core/internal/modules/profitability/adapters/postgres/manual_adjustment_integration_test.go",
    "apps/server_core/internal/modules/product_links/application/generation_live_test.go",
    "apps/server_core/internal/modules/product_links/application/generation_integration_test.go",
    "contracts/governance/runtime-config.json",
    "contracts/governance/execution-lanes.json",
    "contracts/governance/schemas/execution-lanes.schema.json",
    "contracts/governance/schemas/feature-work-contract.schema.json",
    "contracts/governance/schemas/context-pack.schema.json",
    "scripts/tests/governance-contracts.tests.ps1",
    "scripts/tests/governance-drift.tests.ps1",
    "docker-compose.yml",
    "docker/dev/backend-entrypoint.sh",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-03-ephemeral-postgres-migrations/**"
  ],
  "forbidden_paths":[
    "contracts/api/**",
    "packages/**",
    "apps/web/**",
    "apps/server_core/internal/modules/integrations/adapters/mercadolivre/**",
    ".env"
  ],
  "side_effects":{"allowed":["repository-write","database-write"],"forbidden":["external-network","provider-write"]},
  "commands":[
    {"id":"postgres-contract","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/postgres-contract.tests.ps1","lane_id":"unit","expected_exit_code":0},
    {"id":"postgres-lifecycle","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/postgres-lifecycle.tests.ps1","lane_id":"unit","expected_exit_code":0},
    {"id":"postgres-go","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/postgres-go.tests.ps1","lane_id":"unit","expected_exit_code":0},
    {"id":"docker-real","command_template":"npm run harness:integration","lane_id":"integration","expected_exit_code":0},
    {"id":"docker-failure-cleanup","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/postgres-docker.integration.tests.ps1","lane_id":"integration","expected_exit_code":0},
    {"id":"dev-invariance","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/postgres-dev-invariance.integration.tests.ps1","lane_id":"dev-invariance","expected_exit_code":0},
    {"id":"regression-gate","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/f03-regression.tests.ps1 -BaseSha {base_sha}","lane_id":"unit","expected_exit_code":0}
  ],
  "criteria":[
    {"id":"F03-AC01","milestone_criterion_id":"M-08-C03","command_ids":["postgres-contract","postgres-lifecycle","docker-real"]},
    {"id":"F03-AC02","milestone_criterion_id":"M-08-C03","command_ids":["postgres-go","docker-real"]},
    {"id":"F03-AC03","milestone_criterion_id":"M-08-C03","command_ids":["postgres-go","docker-real"]},
    {"id":"F03-AC04","milestone_criterion_id":"M-08-C03","command_ids":["postgres-lifecycle","docker-failure-cleanup"]},
    {"id":"F03-AC05","milestone_criterion_id":"M-08-C03","command_ids":["dev-invariance"]},
    {"id":"F03-AC06","milestone_criterion_id":"M-08-C04","command_ids":["postgres-contract","postgres-lifecycle","docker-failure-cleanup"]},
    {"id":"F03-AC07","milestone_criterion_id":"M-08-C09","command_ids":["regression-gate"]}
  ],
  "stop_conditions":[
    {"code":"resource-leak","condition":"A generated database, container, network, mount, port, or credential artifact survives cleanup."},
    {"code":"dev-state-change","condition":"The validation observer reports a changed persistent dev table-count digest."},
    {"code":"external-side-effect","condition":"Execution requires image pull, external network, provider write, or non-generated database target."},
    {"code":"scope-conflict","condition":"An unrelated dirty path or competing shared-seam owner appears."}
  ],
  "retry_budget":{"max_correction_attempts":1},
  "handoff_fields":["commit-sha","changed-paths","red-evidence","green-evidence","target-labels","resource-inventory","dev-digest","remaining-risks"]
}
```

## Phase 1 — Adversarial RED Contract

Create fake Docker/process fixtures and Go tests before production changes.

- `postgres-contract.tests.ps1`: caller URL rejection, exact generated
  identifiers, no side effect/run dir before validation, static no-volume/no-
  fixed-port contract, explicit CWD/environment, and redaction.
- `postgres-lifecycle.tests.ps1`: ordered start/ready/port/create/migrate twice/
  tests/drop/remove calls; exit `17`, readiness/create/drop/remove failures;
  primary+cleanup failure preservation; active connection; zero label inventory.
- Migration tests: embedded sorted canonical filenames, expected current count,
  per-file transaction failure, second run zero, and no CWD dependency.
- Test-target tests: malformed/remote/dev/wrong-name rejection before pgxpool,
  explicit missing-target failure, and no ambient `MC_DATABASE_URL` fallback.

Run all new fake/unit suites. RED must be the missing lifecycle/embed/helper,
not syntax or unavailable Docker. Commit tests only:

```text
test(harness): define ephemeral postgres contract
```

## Phase 2 — Lifecycle, Embedded Migrations, and Dispatcher GREEN

- Create `scripts/harness/Postgres.psm1` using typed process requests and a
  fakeable executable/argument-prefix seam.
- Embed canonical migrations and refactor runner/command to `fs.FS`; remove
  migration-directory configuration and loop-defer debt.
- Create the central typed test target/helper.
- Route normal integration through generated Docker lifecycle in `try/finally`.
  Keep preflight contact-free; reject any caller `-DatabaseUrl`.
- Update inherited F-08 integration preflight fixtures to `ready`/exit `0` and
  migrate the runner integration test directly to the embedded `fs.FS` API;
  do not retain a string/`any` compatibility bridge.
- Extend the integration child with repo-local offline Go caches and generated
  explicit `MPC_TEST_DATABASE_URL`; do not read reserved targets ambiently.

Run Phase 1 suites plus existing F-08 suites and migration/helper unit tests.
Commit:

```text
feat(harness): add ephemeral postgres lifecycle
```

Use one fixed-commit migration/lifecycle contract gate. Final independent
SPEC/SAFETY and QUALITY reviews remain Phase 3 exit gates.

## Phase 3 — Real DB Suites, Fixtures, and Governance Cutover

- Add integration build tags and central target use to real database tests.
- Add only required typed marketplace/provider/product fixtures; keep FK
  negatives. Split/remove mixed product-links live behavior without claiming
  Oracle evidence.
- Remove every repaired F-03 runtime direct-reader exception and retire
  `MC_MIGRATIONS_DIR`/`MPC_PRODUCT_LINKS_*` keys only when no reader remains.
- Update Compose/entrypoint to use embedded migrations, not a path override.
- Run the real Docker integration lane twice, then forced-failure/active-
  connection cleanup and complete governance/context regressions.

Commit:

```text
refactor(testdb): centralize postgres integration fixtures
```

Run independent SPEC/SAFETY and QUALITY reviews concurrently against one fixed
commit, consolidate findings into at most one correction batch, then perform
focused revalidation. A new architecture contradiction returns to replan
rather than opening review loops.

## Phase 4 — Real Evidence and Handoff

From a clean final HEAD:

1. Capture the read-only dev table-count digest through the existing healthy
   Compose PostgreSQL; emit only digest/count inventory.
2. Run fake contract/lifecycle suites.
3. Run normal real `harness:integration` twice from repo and foreign CWD.
4. Run forced-failure/held-connection Docker integration suite.
5. Verify exact migration set, 32/0 counts, selected Go suites, distinct
   identifiers/ports, zero post-run labelled resources, unchanged dev digest,
   and zero secret/absolute-path leakage.
6. Run F-08 regressions, governance/context gates, and recompile/current-
   validate the final F-03 context pack.

Create `validation.md`, update `feature.md` to
`quick_validation_passed`, and commit:

```text
docs(m08): record ephemeral postgres evidence
```

Evidence must distinguish fake contract runs, real ephemeral PostgreSQL, and
read-only dev observer proof. It cannot claim live provider/Oracle readiness or
M-08 completion.

## Handoff

- Current status: `planned`.
- Next owner: Fresh Phase 1 Feature Implementer.
- Next action: Commit RED fixtures only.
- Review policy: L3, one Phase 2 contract gate and one concurrent final
  SPEC/SAFETY + QUALITY review pair on the fixed Phase 3 commit.
- Blockers: None; do not start a container before Phase 2 GREEN implementation.
