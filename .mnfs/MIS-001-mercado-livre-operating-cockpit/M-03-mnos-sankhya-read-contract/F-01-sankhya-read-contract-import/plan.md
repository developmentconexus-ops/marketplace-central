# Feature Plan

```yaml
id: F-01
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-01
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-oracle-read-contract-redesign

## Steps

1. Add failing contract tests that expose the superseded assumptions and assert MPC-owned policy/value-object boundaries.
2. Rewrite `internal_read/domain` types so they model source facts, quality states, and policy inputs explicitly.
3. Rewrite `internal_read/ports.Reader` and its typed inputs/outputs to match the Oracle-first contract.
4. Run focused `go test` for `./internal/modules/internal_read/...` and record the rewritten contract evidence.

## Files Expected To Change

- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/spec.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/plan.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/validation.md`
- `apps/server_core/internal/modules/internal_read/domain/*.go`
- `apps/server_core/internal/modules/internal_read/domain/*_test.go`
- `apps/server_core/internal/modules/internal_read/ports/reader.go`
- `apps/server_core/internal/modules/internal_read/ports/reader_test.go`

## Verification Commands

- Command: `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/...`
  - Satisfies criterion ID: M-03-C01
  - Expected result: Pass. Domain and port tests prove the contract is MPC-owned and adapter-agnostic.
- Command: `rg -n "MS_DATABASE_URL|MS_TENANT_ID|SANKHYA_" apps/server_core/internal/modules/internal_read`
  - Satisfies criterion ID: M-03-C01
  - Expected result: Pass with only intentional historical compatibility references, not active contract dependence.

## QA Steps

- Step: Review changed paths for boundary discipline.
  - Expected result: only feature artifacts, `internal_read/domain`, and `internal_read/ports` change; no live Oracle query implementation or downstream-module business logic is introduced.

## Rollback/Risk Notes

- Risk: preserving old field names and assumptions can freeze the wrong contract again.
- Recovery: prefer MPC-owned terminology and typed policy objects wherever the old seam was too coupled to the superseded source model.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: implement and record rewritten contract evidence
- Required files/evidence: spec, changed paths, verification commands, QA steps
- Blockers or open decisions: exact real-source policy defaults still need confirmation from Oracle evidence
