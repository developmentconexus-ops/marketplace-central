# F-05 Docker live-Oracle runner plan

```yaml
id: F-05
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-05
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-0
lifecycle_scope: feature
split_decision: single
split_reason: Four scoped implementation paths and no unresolved ownership decision.
```

## Feature ID

F-05-docker-live-oracle-runner

## Steps

1. Compile and validate this feature's context pack against the base SHA and
   exclusive write scope.
2. Add the PowerShell runner and a narrow Docker profile/entrypoint which use
   the existing backend Dockerfile and execute only `TestOracleLiveSmoke`.
3. Add Pester fixtures for the generated Docker invocation, preflight, and
   secret-safe environment forwarding.
4. Document the canonical invocation, prerequisites, and truthful blocked/live
   result handling; run deterministic validation and record evidence.

## Files Expected To Change

- `docker/live-oracle/**`: isolated image execution profile.
- `scripts/run-live-oracle-docker.ps1`: canonical Docker preflight and launch.
- `scripts/tests/live-oracle-docker-runner.tests.ps1`: deterministic runner
  contract tests.
- `docs/operations/live-oracle-docker.md`: operator procedure and boundaries.
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-05-docker-live-oracle-runner/**`:
  feature specification, plan, context, and validation evidence.

## Verification Commands

- Command: `Invoke-Pester scripts/tests/live-oracle-docker-runner.tests.ps1`
  Satisfies criterion IDs: `F05-AC01`, `F05-AC02`.
  Expected result: fixtures prove the exact build/run argv, key-only Docker
  environment forwarding, and prerequisite stops without live credentials.
- Command: `scripts/run-live-oracle-docker.ps1 -PreflightOnly`
  Satisfies criterion ID: `F05-AC02`.
  Expected result: checks Docker and canonical process/parameter inputs without
  building or running a container; outcome is recorded truthfully.
- Command: `scripts/run-live-oracle-docker.ps1`
  Satisfies criterion IDs: `F05-AC01`, `F05-AC02`.
  Expected result: executed only after safe preflight; otherwise recorded as
  could-not-run, never as a pass.

## QA Steps

- Inspect the runner's generated invocation to confirm the only test target is
  `TestOracleLiveSmoke`, there is no `docker compose`, `.env`, migration, or
  application command, and no secret appears in argv/output.

## Rollback/Risk Notes

The runner makes no persistent configuration changes. If the new image profile
is defective, remove only the newly added scoped files. A credential, Docker,
or Oracle availability failure blocks only the optional live execution and
must remain visible in `validation.md`.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-05",
  "required_sources": [
    "docker/dev/backend.Dockerfile",
    "apps/server_core/internal/modules/internal_read/adapters/oracle/reader_live_test.go",
    "contracts/governance/execution-lanes.json",
    "contracts/governance/runtime-config.json",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-contract.md"
  ],
  "knowledge_route_ids": ["portfolio-core", "orders-margin"],
  "allowed_paths": [
    "docker/live-oracle/**",
    "scripts/run-live-oracle-docker.ps1",
    "scripts/tests/live-oracle-docker-runner.tests.ps1",
    "docs/operations/live-oracle-docker.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-05-docker-live-oracle-runner/**"
  ],
  "forbidden_paths": [
    "docker/dev/backend.Dockerfile",
    "docker-compose.yml",
    ".env",
    "scripts/harness/**",
    "scripts/harness.ps1",
    "contracts/governance/**",
    "apps/**",
    "packages/**",
    "contracts/api/**",
    "migrations/**",
    "package.json",
    "Makefile"
  ],
  "side_effects": {
    "allowed": ["repository-write"],
    "forbidden": ["database-mutation", "provider-write"]
  },
  "commands": [
    {"id": "docker-live-oracle-runner-tests", "command_id": "docker-live-oracle-runner-tests", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "docker-live-oracle-preflight", "command_id": "docker-live-oracle-preflight", "lane_id": "live-oracle", "expected_exit_code": 0},
    {"id": "docker-live-oracle-run", "command_id": "docker-live-oracle-run", "lane_id": "live-oracle", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F05-AC01", "milestone_criterion_id": "M-06-C02", "command_ids": ["docker-live-oracle-runner-tests", "docker-live-oracle-run"]},
    {"id": "F05-AC02", "milestone_criterion_id": "M-06-C02", "command_ids": ["docker-live-oracle-runner-tests", "docker-live-oracle-preflight"]}
  ],
  "stop_conditions": [
    {"code": "runtime-contract-conflict", "condition": "Canonical runtime keys or the live-oracle lane require a governance change."},
    {"code": "scope-conflict", "condition": "The runner requires a Docker, Compose, harness, Oracle test, or application change outside the allowed paths."},
    {"code": "live-prerequisite-missing", "condition": "Docker or a nonempty canonical credential is unavailable for a requested live execution."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "commit", "changed-paths", "evidence", "blockers", "next"]
}
```

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Compile and validate the context pack before implementation.
- Required files/evidence: spec, context pack, runner tests, validation record.
- Blockers or open decisions: None.
