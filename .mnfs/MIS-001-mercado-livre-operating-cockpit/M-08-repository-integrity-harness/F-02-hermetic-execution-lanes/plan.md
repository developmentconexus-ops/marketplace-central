# F-02 Hermetic Execution Lanes — Plan

```yaml
id: F-02
type: feature-plan
status: planned
owner: Feature Implementer
parent: M-08
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: feature
split_decision: single-pass
```

## Steps

1. Add a focused PowerShell test covering missing harness, unit environment
   rejection, live redaction/missing-key behavior, provider-write guards, and
   root npm aliases; run it RED before implementation.
2. Implement the Windows-first `scripts/harness.ps1` dispatcher with isolated
   unit, integration, live, browser, and provider-write lanes.
3. Add root npm aliases and ignored run-ID summary output; keep values out of
   stdout and committed artifacts.
4. Run the focused test GREEN, then run unit preflight and the available unit
   command without enabling database, Oracle, migration, or provider targets.
5. Record target-labelled evidence in `validation.md`; stage only F-02 paths
   and commit one intentional feature commit.

## Verification Commands

- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/hermetic-lanes.tests.ps1`
  - Target: fake; expected RED before harness and GREEN after implementation.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command unit -PreflightOnly`
  - Target: fake; expected exit 0 with `.env` ignored and external targets off.
- `npm run harness:unit -- -PreflightOnly`
  - Target: fake; expected same unit preflight through the root alias.
- `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command live -PreflightOnly -EnvFile <fixture>`
  - Target: live-oracle/provider; expected key names only and no values.

## Scope Boundary

Changes are limited to the harness script, focused harness test, root aliases,
ignored run directory rule, and F-02 execution artifacts. F-03 owns ephemeral
PostgreSQL lifecycle/migration behavior; no provider network action is added.
