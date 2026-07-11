# F-02 Hermetic Execution Lanes — Specification

```yaml
id: F-02
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-08
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: feature
```

## Problem

The repository needs named Windows-first execution lanes whose environment
boundaries are visible and testable. Unit tests must not inherit `.env` or live
targets; live preflight must never print values; provider writes must be an
explicit, actor- and idempotency-gated action.

## Requirements

- Expose `harness:unit`, `harness:integration`, `harness:live`,
  `harness:browser`, and `harness:provider-write` npm aliases.
- Unit preflight rejects inherited database, migration, Oracle, provider,
  OAuth, proxy, and live-origin variables and never loads `.env`.
- Integration requires an explicitly named `mpc_test_*` PostgreSQL target.
- Live preflight consumes only target-specific allowlisted keys and emits key
  names/target type, never values.
- Provider-write is a separate command requiring provider, actor, idempotency,
  and explicit execution before any network adapter could be invoked.
- Raw run summaries are written below ignored `scripts/.runs/<run-id>` paths.

## Non-goals

No framework migration, provider business rules, migration implementation, or
masking of existing inventory/contract failures. F-02 does not perform a
provider write or claim live integration evidence.

## Acceptance Criteria

1. `harness.ps1 -Command unit -PreflightOnly` exits 0 with `.env` present and
   reports `target=fake`, all external targets disabled, and `env=ignored`.
2. The same command exits non-zero when a forbidden live/database variable is
   inherited, before test execution.
3. Live preflight reports `target=live-*` and allowlisted `key=` names only;
   missing configuration reports names without values.
4. Provider-write rejects missing actor/idempotency before network and cannot
   run implicitly from unit or live commands.
5. Focused hermetic lane tests pass and root aliases resolve to the versioned
   PowerShell harness.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Required evidence: RED/GREEN lane tests, target-labelled command output,
  and secret/PII redaction checks.
