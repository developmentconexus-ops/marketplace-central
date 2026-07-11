# F-03 Ephemeral PostgreSQL and Canonical Migrations — Validation

```yaml
id: F-03
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Feature ID

F-03-ephemeral-postgres-migrations

## Summary

Quick validation passed. F-03 owns a disposable PostgreSQL 16 lifecycle with
embedded canonical migrations, integration-tagged database suites, typed
FK-complete fixtures, bounded cleanup, hermetic child processes, and explicit
fake/ephemeral/dev-observer evidence boundaries.

This is feature evidence for Milestone Orchestrator review. It does not mark
M-08 passed and does not prove Oracle or marketplace-provider readiness.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-11
- Final feature state for handoff: `quick_validation_passed`

## Evidence Honesty

- Fake contract evidence validates deterministic lifecycle ordering, failure
  classification, redaction, and cleanup policy. It is not database evidence.
- Real integration evidence used locally available `postgres:16-bookworm` with
  `--pull=never`, random loopback ports, generated databases, and tmpfs. It is
  real ephemeral PostgreSQL evidence, not a live provider or Oracle target.
- Dev invariance evidence is a SELECT-only observer against the known healthy
  Compose PostgreSQL. Sessions enforced
  `default_transaction_read_only=on`; only a digest and table count were
  emitted. It is observer evidence, not an integration test target.
- No provider write, provider read, Oracle read, external network request, or
  image pull was executed or claimed.

## Quick Validation State

- fixup_attempts: Consolidated L3 correction batch completed; final isolated
  regression-child finding corrected and revalidated.
- max_fixup_attempts: One consolidated correction batch plus fixed-finding
  revalidation under Milestone Orchestrator direction.
- last_feature_validation_result: Pass on accepted HEAD `7ea94672`.

## Spec Adherence

- Spec satisfied: Yes
- Deviations: The nullable marketplace installation relation required a scoped
  adapter correction: empty optional domain IDs are persisted as SQL `NULL`
  with `NULLIF`, while readback remains the domain empty value. Explicit
  missing-parent fixtures still prove SQLSTATE `23503`.
- Reason: The schema relation is optional, but persisting an empty string
  incorrectly exercised the FK as a real identifier. No API/domain expansion
  or FK weakening was introduced.

## Changes Made

- Commits `218cfd6d`, `9edd3717`, `dac207c6`, and `7ea94672`: added the
  PostgreSQL lifecycle, embedded migration dispatcher, typed test target and
  fixtures, build-tagged real suites, deterministic product-links PostgreSQL
  coverage, governance cutover, read-only dev observer, structured diagnostics,
  confirmed held-connection cleanup, and hermetic regression children.
- `apps/server_core/internal/testsupport/postgres/`: central fail-closed target,
  pool creation, marketplace/provider/installation fixtures, and cleanup
  helpers.
- Database integration suites under `apps/server_core/tests/integration/` and
  the orders, profitability, and product-links packages: use only the generated
  test target and no ambient database fallback or skip.
- `apps/server_core/internal/modules/marketplaces/adapters/postgres/repository.go`:
  maps an absent optional installation ID to SQL `NULL`; true missing parents
  remain rejected.
- `scripts/harness.ps1`, `scripts/harness/Postgres.psm1`, and F-03 test scripts:
  own the disposable container/database lifecycle, safe diagnostics, cleanup,
  dev invariance, and complete regression/current-context gates.
- `docker-compose.yml`, `docker/dev/backend-entrypoint.sh`, and
  `contracts/governance/runtime-config.json`: retired migration-directory and
  product-links temporary runtime readers after source cutover.

## Commands Run

### Fake contract and Go evidence

- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/postgres-contract.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: reject caller targets, enforce generated identifiers/no pull/no
    volume/dynamic loopback, SELECT-only observer source, and safe diagnostics.
  - Actual: `PASS postgres contract tests target=fake`.
  - Artifact: console transcript; committed test contract.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/postgres-lifecycle.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: success/failure ordering, primary exit preservation, bounded
    readiness, cleanup, structured diagnostic fallback, and held-connection
    confirmation contract.
  - Actual: `PASS postgres lifecycle tests target=fake`.
  - Artifact: console transcript; committed fake probe.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/postgres-go.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: embedded sorted migration inventory and fail-closed typed target.
  - Actual: migrate and testsupport packages passed.
  - Artifact: console transcript.

### Real ephemeral PostgreSQL evidence

- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command integration`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: migrations `32/0`, all four selected packages, dynamic port, and
    zero labelled resources.
  - Actual: run `20e9d7c3659643f1821fcefaef2194db`, port `64707`, migrations
    `32/0`, resource count `0`.
  - Artifact: `scripts/.runs/20e9d7c3659643f1821fcefaef2194db/summary.txt`.
- Command: absolute `scripts/harness.ps1 -Command integration` invoked after
  changing the caller CWD to the system temporary directory.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: identical embedded-migration behavior independent of caller CWD.
  - Actual: run `a9797f707cd84d46b3f4abe13420803a`, port `61049`, migrations
    `32/0`, resource count `0`.
  - Artifact: `scripts/.runs/a9797f707cd84d46b3f4abe13420803a/summary.txt`.
- Command: nested integration executed by
  `scripts/tests/postgres-dev-invariance.integration.tests.ps1`.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: a distinct real ephemeral target inside the observer wrapper.
  - Actual: nested target `ephemeral-postgres` passed; its generated target URL
    and password were not emitted.
  - Artifact: dev-invariance console transcript.
- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/postgres-docker.integration.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: confirm a live held connection, force child exit `17`, preserve
    that exit, force-drop the database, remove the container, and leave zero
    resources.
  - Actual: final run `5fe5476b2c044888bfa94eb2c8205eaa`, port `62061`,
    `held_connection=confirmed-and-force-dropped`, child exit `17`, resource
    count `0`.
  - Artifact: console transcript.

### Read-only dev observer evidence

- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/postgres-dev-invariance.integration.tests.ps1`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: read-only session proof, SELECT-only per-table counts, equal
    before/after digest, and a passing nested ephemeral integration.
  - Actual: digest
    `54007d0aff84ebad9f123e1aa92f5299ef65de28ef6dc23a0a2858c53e87b029`,
    table count `34`, before/after equal, nested target passed.
  - Artifact: console transcript; committed observer source.

### Governance, regression, and current context

- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/f03-regression.tests.ps1 -BaseSha 7ea94672de4a8fad43511a87abad31c52e301c50`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: every fake/Go/F-08 regression marker, governance contracts and
    drift, context compiler/current validation, and typed children isolated
    from contaminated parent DB/Oracle/provider values.
  - Actual: external exit `0`; `contamination_count=0`; postgres contract,
    lifecycle, Go, environment, execution, hermetic, aliases, governance,
    context compiler, context current, and aggregate markers all passed.
  - Artifact:
    `scripts/.runs/99ba9a36fcea4401851db71129e420e0/context-pack.json`.
- Command: post-correction context compile/validate and governance drift at
  accepted HEAD `7ea94672`.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: source hashes/base/current governance agree with accepted HEAD.
  - Actual: context compile, `RequireCurrentBase`, governance validate, and
    governance drift passed.
  - Artifact:
    `scripts/.runs/99ba9a36fcea4401851db71129e420e0/context-pack.json` and
    `contracts/governance/`.
- Command: boundary searches for retired runtime keys and final Docker
  inventory by exact run label/name.
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Expected: zero retired runtime references in executable/governance paths and
    zero `mpc-pg-*` or run-labelled resources.
  - Actual: retired reference count `0`; label resource count `0`; name resource
    count `0`.
  - Artifact: console transcript; governance zero-reference assertion.

## Manual QA

- QA level: QA-3
- Flow or step: Independent fixed-commit SPEC/SAFETY and QUALITY review,
  consolidated correction, and fixed-finding revalidation.
- Status: Pass
- Evidence type: ran
- Owner: Milestone Orchestrator / independent reviewers
- Expected: no unresolved P0/P1/P2 finding on accepted HEAD.
- Actual: review findings were consolidated, corrected in `dac207c6` and
  `7ea94672`, and the final isolation finding was revalidated on accepted HEAD.
- Blocking condition: None for F-03 quick-validation handoff.

## Evidence

- Artifact: commits `218cfd6d`, `9edd3717`, `dac207c6`, and `7ea94672`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: None.
- Artifact: real ephemeral run summaries and console transcripts listed above
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: None.
- Artifact: current context pack
  `scripts/.runs/99ba9a36fcea4401851db71129e420e0/context-pack.json`
  - Status: Pass
  - Evidence type: ran
  - Owner: Feature Implementer
  - Blocking condition: None.

## Risks

- F-04 still owns cold provisioning and cold-gate proof. F-03 used the local
  image and locally available Go dependencies with network disabled; this does
  not prove a cold machine can provision them.
- F-04 still owns real Oracle/provider read evidence. F-03 deliberately removed
  the mixed Oracle/product-links test and makes no live integration claim.
- The dev digest proves invariance only across the observed F-03 window; normal
  application activity can legitimately change persistent dev counts outside
  that bounded observation.
- No unresolved F-03 functional or cleanup risk remains in quick validation.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: Review F-03 spec, plan, changed paths, validation evidence, and
  accepted fixed-commit review result; accept or reject the feature for M-08
  integration.
- Required files/evidence: `feature.md`, `spec.md`, `plan.md`, this
  `validation.md`, commits through `7ea94672`, real run summaries, dev digest,
  governance/current context evidence, and independent review result.
- Blockers or open decisions: None for F-03 acceptance. M-08 remains in
  progress and must not be marked passed from this feature alone.
