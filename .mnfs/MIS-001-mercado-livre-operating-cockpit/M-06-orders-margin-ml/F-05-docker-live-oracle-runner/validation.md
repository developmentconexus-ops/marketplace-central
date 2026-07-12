# F-05 Docker live-Oracle runner validation

```yaml
id: F-05
type: feature-validation
status: quick_validation_passed
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

This correction makes the canonical runner usable with the ignored local
secure `.env` without Compose-wide inheritance. DBeaver connection evidence
establishes that the supplied `MPC_SANKHYA_ORACLE_CONNECT_STRING` is a service
name. Its narrow parser accepts the five connection inputs
`MPC_SANKHYA_ORACLE_USERNAME`, `MPC_SANKHYA_ORACLE_PASSWORD`,
`MPC_SANKHYA_ORACLE_HOST`, `MPC_SANKHYA_ORACLE_PORT`, and
`MPC_SANKHYA_ORACLE_CONNECT_STRING`; `MPC_SANKHYA_ORACLE_SCHEMA` is permitted
but ignored as unrelated local configuration. Unrelated non-reserved
assignments are ignored without reading their values; unknown reserved keys,
generic/ambient aliases, duplicate reserved keys, or malformed lines stop
before Docker. A nonempty caller-process value with the same exact name
overrides its corresponding local value. Generic
`MPC_ORACLE_*` credentials and ambient aliases cannot satisfy preflight. After
resolution, the runner constructs the governed `host:port/service` connect
string and maps only it plus the resolved username and password to the
pre-existing governed container names `MPC_ORACLE_USERNAME`,
`MPC_ORACLE_PASSWORD`, and `MPC_ORACLE_CONNECT_STRING` consumed by the Oracle
configuration. Docker arguments remain key-only, and the existing read-only
smoke-test target and fail-closed child environment are unchanged.

The correction changes no runtime governance. It preserves the permitted
runner boundary mapping to the existing governed container names and does not
read, log, stage, or persist credential values.

No Docker, Oracle, application, migration, or provider command was run. The
dispatch requires a separate read-only live-validation session after review.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-12
- Final feature state for handoff: `quick_validation_passed`

## Evidence Honesty

| Check | Evidence type | Result | Actual result / artifact |
| --- | --- | --- | --- |
| Spec and plan correction | ran | Pass | `spec.md` and `plan.md` define the exact local-file whitelist, documented caller-process precedence, existing governed container mapping, and deterministic proof. |
| `Invoke-Pester scripts/tests/live-oracle-docker-runner.tests.ps1` | ran | Pass | 10 passed, 0 failed. Covers service-name DSN construction, ignored unrelated-local and schema isolation, unknown reserved-key and alias rejection before Docker, exact caller-process precedence, one-way governed mapping, key-only Docker argv, isolated child environment, and fixed smoke-test target. |
| Static connection-boundary inspection | ran | Pass | The runner reads only the five explicit caller connection names; the parser permits but does not retain schema or unrelated non-reserved assignments, rejects unknown reserved keys and aliases, and Docker receives the five governed runtime keys by reference only. |
| `git diff --check` | ran | Pass | No whitespace errors in the scoped correction. |
| Docker preflight or live smoke test | assumed | Not run | Deliberately deferred: this bounded correction prohibits Docker/Oracle activity before contract tests and delegates live validation to a separate read-only session after acceptance. |

## Spec Adherence

- F05-AC01: Pass. Deterministic coverage preserves the existing Dockerfile,
  read-only mount, exact `TestOracleLiveSmoke` target, key-only forwarding, and
  the exact governed container keys consumed by Oracle configuration.
- F05-AC02: Pass. Deterministic coverage proves service-name construction, the
  reserved local-file allowlist, unrelated-key and schema isolation,
  documented caller-process precedence, generic/ambient rejection, and that
  the one-way mapping does not place connection values in Docker argv/output.
- Scope: only the assigned runner, its deterministic tests, operator document,
  and F-05 correction artifacts changed. No runner profile change was needed.

## Changed Paths

- `scripts/run-live-oracle-docker.ps1`
- `scripts/tests/live-oracle-docker-runner.tests.ps1`
- `docs/operations/live-oracle-docker.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-05-docker-live-oracle-runner/spec.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-05-docker-live-oracle-runner/plan.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-05-docker-live-oracle-runner/validation.md`

## Risks

- The live Docker/Oracle lane remains unexecuted by design. Its evidence must
  be produced only by the separate read-only validation session using the
  secure local file or the same exact caller-process namespace.
- This correction relies on the pre-existing governed container names named by
  the independent review. Any future request to alter those canonical names
  requires an explicit runtime-governance stop and a separately scoped change.
- This record is Feature Implementer quick-validation evidence, not milestone
  acceptance or QA verdict.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: fixed-SHA review, proportional QA, then dispatch the separate
  read-only live-validation session if accepted.
- Required files/evidence: this validation record, the Pester result, and the
  future read-only live validation result.
- Blockers or open decisions: None for this bounded correction; live evidence
  is intentionally deferred.

## Correction M06-F05-preflight-order-retry-1

- Retry: 1 of 2.
- Original fixed-SHA review failure trace preserved: at
  `19022d6ee3153f325778ede1b85f0d7f281e5213`,
  `Invoke-LiveOracleDockerRunner` called
  `Test-LiveOracleDockerAvailable` (`docker version`) before
  `New-LiveOracleDockerPlan` parsed the local `.env`. This did not prove that
  invalid, non-whitelisted, duplicate, or malformed local-file input fails
  before any Docker request.
- Smallest correction: resolve the runner's credentials plan before its Docker
  availability check. The runner now accepts an optional test-only local-file
  path parameter while retaining the canonical ignored `.env` default.
- Targeted proof: the runner-level Pester case provides a non-whitelisted,
  opaque fixture key, mocks Docker availability, asserts the narrow parser's
  rejection, and asserts the Docker availability command was invoked zero
  times. No Docker, Oracle, provider, migration, or application command is
  run by this correction.
- Command evidence (ran): `Invoke-Pester -Path
  scripts/tests/live-oracle-docker-runner.tests.ps1` passed 8, failed 0,
  skipped 0, pending 0, inconclusive 0. This includes the runner-order
  regression test.
- Command evidence (ran): `git diff --check` passed with no whitespace errors.

## Correction M06-F05-structured-service-contract

- Authority: owner-supplied DBeaver evidence selects `ORCL` as **Service Name**
  for the supplied host and port.
- Runtime route gap resolved: the governed Oracle configuration receives
  `MPC_ORACLE_CONNECT_STRING`, and the adapter test establishes the expected
  `host:1521/service` value form. No governance, adapter, or profile change was
  required.
- Smallest correction: require host, port, and service name together with the
  existing username and password; construct `host:port/service` only inside the
  child environment. Permit `MPC_SANKHYA_ORACLE_SCHEMA` only as ignored local
  configuration. It is neither retained by the parser nor forwarded, logged, or
  used for routing.
- Targeted proof: deterministic Pester tests pass 9/9, including service-name
  construction, exact caller precedence, schema ignored/non-forwarded,
  unknown/alias rejection, key-only Docker forwarding, and isolated child
  environment. `git diff --check` passed. No Docker, Oracle, migration,
  application, or provider command was run.

## Correction M06-F05-unrelated-env-policy-retry-1

- Retry: 1 of 2.
- Original failure trace preserved: host validation at
  `755c9cf6bbf5b54521aea6dd403211b0b20087c9` stopped before Docker because
  the strict parser rejected the unrelated local assignment category
  `MC_DATABASE_URL`. No values were observed or output, and the canonical
  runner was not invoked.
- Assigned failure addressed: a local ignored `.env` was incorrectly treated
  as entirely owned by the runner rather than as a source for its exact
  reserved inputs.
- Smallest correction: ignore non-reserved assignments without reading or
  retaining their values. Within `MPC_SANKHYA_ORACLE_*`, permit only the five
  connection inputs and local-only `SCHEMA`; reject unknown reserved keys.
  Reject generic and ambient credential aliases, including `MPC_ORACLE_*`,
  before Docker. Exact caller-process precedence and the one-way governed
  container mapping are unchanged.
- Files changed: `scripts/run-live-oracle-docker.ps1`,
  `scripts/tests/live-oracle-docker-runner.tests.ps1`,
  `docs/operations/live-oracle-docker.md`, and this F-05 `spec.md`, `plan.md`,
  and `validation.md`.
- Targeted validation evidence (ran): `Invoke-Pester -Path
  scripts/tests/live-oracle-docker-runner.tests.ps1` passed 10, failed 0,
  skipped 0, pending 0, inconclusive 0. Fixtures prove unrelated
  `MC_DATABASE_URL` isolation plus unknown-reserved and generic-alias
  rejection with Docker availability mocked and asserted at zero calls.
- Targeted validation evidence (ran): `git diff --check` passed with no
  whitespace errors in the scoped correction.
- Runtime: not run by correction scope. No Docker, Oracle, application,
  migration, or provider command was invoked.
- Original blocking failure resolved: Yes, deterministically. A separate
  read-only host Docker validation remains required for live evidence.
