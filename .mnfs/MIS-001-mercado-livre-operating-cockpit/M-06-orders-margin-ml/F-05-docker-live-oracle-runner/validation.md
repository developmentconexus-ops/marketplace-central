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

Correction attempt 1 resolves the fixed-SHA review finding. The runner accepts
only the explicit caller-process namespace `MPC_SANKHYA_ORACLE_USERNAME`,
`MPC_SANKHYA_ORACLE_PASSWORD`, and
`MPC_SANKHYA_ORACLE_CONNECT_STRING`; generic `MPC_ORACLE_*` credentials and
ambient aliases cannot satisfy preflight. After resolution, the runner maps
only those values to the pre-existing governed container names
`MPC_ORACLE_USERNAME`, `MPC_ORACLE_PASSWORD`, and
`MPC_ORACLE_CONNECT_STRING` consumed by the Oracle configuration. Docker
arguments remain key-only, and the existing read-only smoke-test target and
fail-closed child environment are unchanged.

Original fixed-SHA review trace: the prior runner forwarded
`MPC_SANKHYA_ORACLE_*` into the container while the Oracle adapter consumes
governed `MPC_ORACLE_*`; no Dockerfile mapper existed, so a real smoke test
would fail (P0). The reviewer also noted that runtime configuration already
contains the governed container names and required a stop if the canonical
runtime contract changed (P1). This correction changes no runtime governance:
it repairs only the permitted runner boundary mapping to those existing names.

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
| Spec and plan correction | ran | Pass | `spec.md` and `plan.md` distinguish exclusive Sankhya caller inputs from the existing governed Oracle container keys and map F05-AC01/F05-AC02 to deterministic runner tests. |
| `Invoke-Pester -Path scripts/tests/live-oracle-docker-runner.tests.ps1` | ran | Pass | 6 passed, 0 failed. Includes generic `MPC_ORACLE_*` and ambient-alias input rejection, exact mapping from resolved Sankhya values to the three governed container names, five-key forwarding, no credential values in Docker argv, isolated child environment, and fixed smoke-test target. |
| Static credential-name audit of `scripts/run-live-oracle-docker.ps1` | ran | Pass | `static-name-audit=passed`: all three exclusive caller names and all three governed container names are present, and no governed container key is read as a caller-process input. |
| `git diff --check` | ran | Pass | No whitespace errors in the scoped correction. |
| Docker preflight or live smoke test | assumed | Not run | Deliberately deferred: this bounded correction prohibits Docker/Oracle activity before contract tests and delegates live validation to a separate read-only session after acceptance. |

## Spec Adherence

- F05-AC01: Pass. Deterministic coverage preserves the existing Dockerfile,
  read-only mount, exact `TestOracleLiveSmoke` target, key-only forwarding, and
  the exact governed container keys consumed by Oracle configuration.
- F05-AC02: Pass. Deterministic coverage proves the explicit Sankhya keys are
  required, generic/ambient alternatives are rejected, and the one-way mapping
  does not print any credential value.
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
  explicit caller-process namespace.
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
