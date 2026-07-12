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

The Docker runner credential contract now exclusively uses the explicit caller
process namespace `MPC_SANKHYA_ORACLE_USERNAME`,
`MPC_SANKHYA_ORACLE_PASSWORD`, and
`MPC_SANKHYA_ORACLE_CONNECT_STRING`. Generic `MPC_ORACLE_*` credentials and
ambient aliases cannot satisfy preflight. Docker arguments remain key-only, and
the existing read-only smoke-test target and fail-closed child environment are
unchanged.

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
| Spec and plan correction | ran | Pass | `spec.md` and `plan.md` name the explicit Sankhya namespace and map F05-AC01/F05-AC02 to deterministic runner tests. |
| `Invoke-Pester -Path scripts/tests/live-oracle-docker-runner.tests.ps1` | ran | Pass | 6 passed, 0 failed. Includes generic `MPC_ORACLE_*` and ambient-alias rejection, exact five-key forwarding, no credential values in Docker argv, isolated child environment, and fixed smoke-test target. |
| Static credential-name audit of `scripts/run-live-oracle-docker.ps1` | ran | Pass | The runner contains only `MPC_SANKHYA_ORACLE_USERNAME`, `MPC_SANKHYA_ORACLE_PASSWORD`, and `MPC_SANKHYA_ORACLE_CONNECT_STRING` as credential keys; the remaining `MPC_ORACLE_*` names are non-credential live-test/library settings. |
| `git diff --check` | ran | Pass | No whitespace errors in the scoped correction. |
| Docker preflight or live smoke test | assumed | Not run | Deliberately deferred: this bounded correction prohibits Docker/Oracle activity before contract tests and delegates live validation to a separate read-only session after acceptance. |

## Spec Adherence

- F05-AC01: Pass. Deterministic coverage preserves the existing Dockerfile,
  read-only mount, exact `TestOracleLiveSmoke` target, and key-only forwarding.
- F05-AC02: Pass. Deterministic coverage proves the explicit Sankhya keys are
  required and generic/ambient alternatives are rejected without printing any
  credential value.
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
