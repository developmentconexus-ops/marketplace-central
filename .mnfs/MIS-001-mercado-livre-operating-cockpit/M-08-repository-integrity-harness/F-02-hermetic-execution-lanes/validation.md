# F-02 Hermetic Execution Lanes — Quick Validation

```yaml
id: F-02
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-08
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: feature
```

## Summary

The versioned PowerShell harness and root npm aliases now expose isolated unit,
integration, live, browser, and provider-write lanes. Unit execution does not
load `.env`; live output contains target/key names only; provider-write rejects
missing actor or idempotency before any network path. No provider write, real
database, Oracle, or browser action was invoked in this feature validation.

## RED/GREEN Evidence

- RED — `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/hermetic-lanes.tests.ps1`
  - Target: fake.
  - Actual: exit 1 before implementation because `scripts/harness.ps1` was absent.
  - Evidence type: ran; no values or external state consumed.
- GREEN — same command after implementation.
  - Target: fake.
  - Actual: exit 0, `PASS hermetic lane tests`.
  - Coverage: unit inherited-config rejection, live allowlist/redaction and
    missing-key behavior, provider actor/idempotency guard, and root aliases.

## Commands Run

- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command unit -PreflightOnly`
  - Target: fake.
  - Actual: exit 0; `target=fake`, `env=ignored`, PostgreSQL/Oracle/provider
    network/migrations all disabled. A repository `.env` remained present and
    was not read.
- `npm run harness:unit -- -PreflightOnly`
  - Target: fake.
  - Actual: exit 0 through the root alias; same disabled-target inventory.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command unit`
  - Target: fake.
  - Actual: exit 0; Go unit package passed (`23.126s`) and web Vitest passed
    (`20` files, `185` tests). Existing React `act(...)` warnings are non-fatal
    baseline output and were not masked.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command live -PreflightOnly`
  - Target: live-oracle (preflight only).
  - Actual: exit 0; emitted `target=live-oracle` and only allowlisted canonical
    or legacy Oracle key names. No secret values were printed or persisted to
    committed artifacts.
- Focused GREEN test provider case.
  - Target: live-provider (guard only).
  - Actual: missing actor and idempotency rejected before network; no provider
    adapter or HTTP request was invoked.

Raw summaries are generated under ignored `scripts/.runs/<run-id>` directories;
committed evidence contains no environment values, secrets, or buyer PII.

## Scope and Limitations

- Integration requires an explicitly supplied `MPC_TEST_DATABASE_URL` whose
  database name matches `mpc_test_*`; F-03 owns lifecycle, migrations, and
  cleanup, so no real PostgreSQL evidence is claimed here.
- Browser lane is an explicit preflight surface and reports that browser
  automation must be invoked separately; no browser target was contacted.
- Provider-write remains intentionally outside F-02 business behavior; even
  `-Execute` stops with a no-network boundary message.

## Changed Paths

- `.gitignore`
- `package.json`
- `scripts/harness.ps1`
- `scripts/tests/hermetic-lanes.tests.ps1`
- `F-02-hermetic-execution-lanes/spec.md`
- `F-02-hermetic-execution-lanes/plan.md`
- `F-02-hermetic-execution-lanes/validation.md`

## Handoff

- Current status: `quick_validation_passed`
- Next owner: final F-02 SPEC+QUALITY reviewers
- Next action: review exact allowlists, redaction behavior, target labels, and
  command aliases; do not infer real integration/provider readiness from fake
  or preflight evidence.
- Blockers: none for F-02 scoped harness behavior; live and ephemeral database
  validation remain later-lane responsibilities.

## Correction Batch — Final Review Findings

### RED/GREEN

- RED — focused test after adding final-review cases:
  `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/hermetic-lanes.tests.ps1`
  - Target: fake.
  - Actual: exit 1; inherited `MPC_TEST_DATABASE_URL`, `MPC_PROVIDER_*`,
    `MPC_ORACLE_*`, legacy Oracle, proxy, and tunnel cases were not all
    rejected and the prior catch path mutated key identifiers.
- GREEN — same focused command after correction:
  - Target: fake.
  - Actual: exit 0, `PASS hermetic lane tests`.
  - Coverage: anchored/prefix unit guard with key-only reporting; exact
    required/missing/forbidden names; process and EnvFile legacy
    `SANKHYA_ORACLE_*`; integration ambient-config rejection and F-03
    fail-closed boundary; subprocess redaction; npm alias execution; and
    provider-write no-network guards with and without `-Execute`.

### Required commands

- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command unit -PreflightOnly`
  - Target: fake; exit 0; `.env` ignored and PostgreSQL, Oracle, provider
    network, and migrations disabled.
- `npm run harness:unit -- -PreflightOnly`
  - Target: fake; exit 0 through the root alias; same disabled-target output.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command unit`
  - Target: fake; exit 0; Go unit package passed with `GOCACHE=.gocache`, and
    web Vitest passed (20 files, 185 tests). Existing React `act(...)` and
    Node deprecation warnings remain non-fatal baseline output.
- Integration preflight with no explicit `-DatabaseUrl` exits 1 and reports
  only the required key name; ambient Oracle/provider/OAuth/proxy/database
  values are ignored. With an explicit `mpc_test_*` URL it emits
  `target=ephemeral-postgres` and `key=MPC_TEST_DATABASE_URL`, then exits 1
  with the F-03 lifecycle blocker before any DB/network attempt.
- Provider-write with valid actor/idempotency exits 1 before network without
  `-Execute`; `-Execute` is also rejected as outside F-02. No provider adapter
  or HTTP request was invoked.

### Scope limitations

- This correction batch proves deterministic fake/preflight behavior only.
  It does not validate a real PostgreSQL, Oracle, OAuth, provider, browser,
  proxy/tunnel, migration, or network target. F-03 owns the ephemeral
  PostgreSQL lifecycle and migration evidence.
- Subprocess output is captured and structurally redacted before forwarding;
  key identifiers in required/missing/forbidden messages remain exact.
- Run summaries are created only when a lane reaches `Write-Summary`; blocked
  lanes do not create new empty `scripts/.runs` directories.
