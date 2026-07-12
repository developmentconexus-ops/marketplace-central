# F-05 Docker live-Oracle runner validation

```yaml
id: F-05
type: feature-validation
status: blocked
owner: Feature Implementer
parent: F-05
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-05-docker-live-oracle-runner

## Summary

This correction resolves the Milestone Orchestrator's rejection of the original
F-05 return: its Docker `ProcessStartInfo` populated keys without first clearing
the inherited process environment. The runner now clears the Docker child,
restores only an explicit Windows execution allowlist, and adds the five exact
container runtime keys only for `docker run`. The actual live Oracle build/run
remains blocked because the available Docker CLI cannot reach its daemon. No
build, container, or Oracle connection was attempted, and no live success is
claimed.

## Quick Validation Result

- Result: Blocked
- Result owner: Feature Implementer
- Decision date: 2026-07-12
- Final feature state for handoff: `blocked`

## Evidence Honesty

| Check | Evidence type | Result | Actual result / artifact |
| --- | --- | --- | --- |
| Context compile | ran | Pass | `scripts/.runs/a522db0e0be649f5b9cdf39e918754ec/context-pack.json`; status passed. |
| Context validate | ran | Pass | Same context pack; status passed with current base SHA. |
| Direct inspection of the rejected implementation | ran | Fail reproduced | `Invoke-LiveOracleDockerCommand` added values to `ProcessStartInfo.Environment` without `Environment.Clear()`, so the Docker child inherited unrelated caller variables. |
| `Invoke-Pester -Path scripts/tests/live-oracle-docker-runner.tests.ps1` | ran | Pass | 5 passed, 0 failed. The new deterministic fixture constructs the Docker child start info after adding ambient MPC/database/Oracle/provider/legacy-alias values; it proves they are absent and the `MPC_*` runtime set is exactly the five canonical keys. It would fail without the cleared allowlist construction. Existing fixtures still prove the backend Dockerfile, exact smoke-test target, read-only mount, key-only Docker forwarding, no secret argv values, canonical precedence, and missing-key stop. |
| `go test ./internal/modules/internal_read/adapters/oracle -run '^TestLoadConfigFromEnvDoesNotLeakSecretValues$' -count=1` with workspace-local absolute `GOCACHE` | ran | Pass | Oracle config unit test passed. An initial relative `GOCACHE=.gocache` attempt was rejected by Windows Go because the cache path must be absolute; the test was then rerun using the equivalent absolute `apps/server_core/.gocache` path. |
| `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/run-live-oracle-docker.ps1 -PreflightOnly` | could-not-run | Blocked | Docker CLI was present, but Docker could not load its user config and could not connect to its named-pipe daemon (exit 1). The runner stopped before credential handling, build, or run. |
| Live Docker build and `TestOracleLiveSmoke` execution | could-not-run | Blocked | Not attempted because Docker availability preflight failed. No simulated result. |
| Manual command/profile inspection | ran | Pass | `docker/live-oracle/profile.json` and the generated plan use only `docker build` with `docker/dev/backend.Dockerfile`, `docker run --rm`, a read-only checkout mount, the Oracle package, and `^TestOracleLiveSmoke$`; no Compose, `.env`, migration, or app command. |

## Correction Applied

- `New-LiveOracleDockerProcessStartInfo` now sets `UseShellExecute=false`, calls
  `Environment.Clear()`, restores only `SystemRoot`, `WINDIR`, `ComSpec`,
  `PATH`, `PATHEXT`, `TEMP`, `TMP`, `USERPROFILE`, `APPDATA`, and
  `LOCALAPPDATA` when present, and then adds the supplied runtime dictionary.
- The build child receives only the execution allowlist. The run child receives
  that allowlist plus exactly `MPC_ORACLE_USERNAME`, `MPC_ORACLE_PASSWORD`,
  `MPC_ORACLE_CONNECT_STRING`, `MPC_ORACLE_LIVE_TEST`, and
  `MPC_ORACLE_LIB_DIR`.
- Docker CLI arguments remain key-only `--env` references; no credential value
  is placed in argv or emitted by runner status output.

## Spec Adherence

- `F05-AC01`: deterministic Pester coverage passed; the real target remains
  blocked solely by Docker daemon availability.
- `F05-AC02`: deterministic Pester coverage passed; actual preflight truthfully
  stopped before a Docker build/run.
- Scope correction: `docs/research/evidence/2026-07-12-ml-price-monitoring-manifest.md`
  was observed as an out-of-scope untracked artifact. This feature did not
  create or modify it, so it was preserved for its owner.

## Quick Validation State

- fixup_attempts: 1
- max_fixup_attempts: 1
- Original defect resolution: Yes — verified by the new deterministic child
  environment test (5 Pester tests passed).
- Correction fixup: replaced inherited `ProcessStartInfo.Environment` use with
  a clear-and-allowlist start-info constructor; added ambient-environment
  regression coverage and operator documentation.
- Full quick-validation plan rerun after correction: deterministic Pester tests
  and targeted Oracle config unit test passed; Docker preflight remained
  blocked.

## Changed Paths

Current correction paths:

- `scripts/run-live-oracle-docker.ps1`
- `scripts/tests/live-oracle-docker-runner.tests.ps1`
- `docs/operations/live-oracle-docker.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-05-docker-live-oracle-runner/feature.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-05-docker-live-oracle-runner/validation.md`

The amended F-05 commit also retains the original scoped feature paths:

- `docker/live-oracle/profile.json`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-05-docker-live-oracle-runner/spec.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-05-docker-live-oracle-runner/plan.md`

## Risks

- Docker user-config access and Docker daemon/API access must be restored before
  the canonical live lane can be executed. Oracle connectivity is still
  unverified.
- Docker command output is intentionally suppressed by the runner to prevent
  accidental credential disclosure; operators should retain only safe status
  evidence.

## Handoff

- Current status: `blocked`
- Next owner: Milestone Orchestrator
- Next action: review the scoped implementation; when Docker daemon access is
  available, rerun preflight and then the canonical runner with credentials
  supplied only through the calling process or explicit secure parameters.
- Required files/evidence: this validation record, deterministic Pester result,
  targeted Oracle config test result, and a future live execution result.
- Blockers or open decisions: Docker daemon/API unavailable.

## Milestone Acceptance Review — 2026-07-12

- Decision: **accepted for integration**.
- Scope: the implementation changed only the assigned Docker profile, runner,
  deterministic checks, operator documentation, and F-05 artifacts. The M-06
  feature table now records F-05 as validation support.
- Evidence: context compile/validate, the five deterministic runner checks,
  and the targeted Oracle configuration test were all rerun with `ran` pass
  evidence. The correction proves the Docker child clears inherited runtime
  configuration before it adds its allowlist.
- Constraint: the actual Docker build and `TestOracleLiveSmoke` execution are
  still blocked by Docker daemon/API access. This does not reduce the runner's
  deterministic implementation evidence to a live pass; it remains pending
  operational evidence for later QA.
- Next owner: independent fixed-SHA reviewer, then proportional QA.

## Fixed-SHA Review Correction — 2026-07-12

- Review finding reproduced: `Test-LiveOracleDockerAvailable` used
  `Start-Process` for `docker version --format {{.Server.Version}}`. That
  bypassed `New-LiveOracleDockerProcessStartInfo`, so the preflight child could
  inherit ambient MPC, database, Oracle, provider, legacy, or secret values.
- Correction: `New-LiveOracleDockerPreflightPlan` now exposes the exact
  preflight arguments and `ProcessStartInfo`. It creates `docker version
  --format {{.Server.Version}}` with
  `New-LiveOracleDockerProcessStartInfo` and an empty runtime dictionary.
  `Test-LiveOracleDockerAvailable` executes that start info through the same
  redacted process helper used by build and run; `Start-Process` is no longer
  used by the runner.
- Deterministic reproduction and proof: the added Pester fixture injects
  ambient MPC, canonical runtime, database, Oracle, provider, and legacy alias
  values. It verifies the preflight start info contains exactly the available
  OS execution allowlist, no `MPC_*` runtime key, and none of the injected
  values. The existing run-child fixture still verifies the exact five-key
  runtime set.
- Validation rerun: `Invoke-Pester -Path
  scripts/tests/live-oracle-docker-runner.tests.ps1` passed with 6 passed and
  0 failed. `go test ./internal/modules/internal_read/adapters/oracle -run
  '^TestLoadConfigFromEnvDoesNotLeakSecretValues$' -count=1` passed with the
  workspace-local absolute `GOCACHE`.
- Docker preflight/live run: not rerun by correction scope; the previously
  recorded Docker daemon/API blocker remains unchanged. No Docker build,
  container, Oracle connection, provider operation, or database operation was
  attempted in this correction.
- Original finding resolved: **Yes**. This is deterministic start-info evidence
  only; it is not a QA pass and does not supersede the daemon blocker.
