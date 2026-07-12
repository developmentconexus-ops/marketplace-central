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
2. Add a narrow ignored-local-`.env` parser which accepts exactly the five
   `MPC_SANKHYA_ORACLE_*` connection inputs, permits but ignores the local
   schema assignment, ignores unrelated non-reserved `.env` keys, rejects
   unknown reserved keys and credential aliases, and gives matching
   caller-process values documented precedence; construct the governed
   `host:port/service` connect string, then map resolved values only to the pre-existing governed
   `MPC_ORACLE_*` container credential names while retaining the existing Docker
   profile and read-only smoke-test target.
3. Add Pester fixtures proving service-name construction, the `.env` whitelist,
   caller-process precedence, schema isolation, unrelated-key isolation,
   generic `MPC_ORACLE_*` and ambient-alias rejection, exact mapping to the
   governed container names, and secret-safe Docker forwarding.
4. Document the narrow local `.env` contract and precedence; run deterministic
   contract tests before any Docker or Oracle activity and record evidence.

## Files Expected To Change

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
  environment forwarding, narrow `.env` allowlist, caller-process precedence,
  generic and ambient-alias rejection, and prerequisite stops without live
  credentials.

## QA Steps

- Inspect the runner's generated invocation to confirm the only credential names
  accepted from local `.env` or caller process are `MPC_SANKHYA_ORACLE_*`,
  caller process precedence is explicit, forwarding uses only the governed
  container names `MPC_ORACLE_USERNAME`, `MPC_ORACLE_PASSWORD`, and
  `MPC_ORACLE_CONNECT_STRING`, the only test target is `TestOracleLiveSmoke`,
  there is no `docker compose`, Compose-wide `.env` inheritance, migration, or
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
