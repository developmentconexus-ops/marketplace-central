# Feature Plan

```yaml
id: F-02
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-02
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-02-oracle-adapter-implementation

## Steps

1. Add failing tests for config secrecy, adapter-boundary ownership, and deterministic service behavior.
2. Implement or rewrite `application.Service` only as a thin consumer-facing layer over `ports.Reader`.
3. Implement deterministic fake adapter behavior that mirrors the new contract semantics.
4. Implement real Oracle config, connection, query organization, mapping, and error translation.
5. Run focused tests and record what still requires real-environment validation.

## Files Expected To Change

- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-02-read-adapter-seam/spec.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-02-read-adapter-seam/plan.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-02-read-adapter-seam/validation.md`
- `apps/server_core/internal/modules/internal_read/application/service.go`
- `apps/server_core/internal/modules/internal_read/application/service_test.go`
- `apps/server_core/internal/modules/internal_read/adapters/fake/reader.go`
- `apps/server_core/internal/modules/internal_read/adapters/fake/reader_test.go`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/*.go`

## Verification Commands

- Command: `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/...`
  - Satisfies criterion ID: M-03-C03
  - Expected result: Pass. Service, fake adapter, and Oracle adapter compile and pass scoped tests.
- Command: `rg -n "INSERT|UPDATE|DELETE|MERGE" apps/server_core/internal/modules/internal_read/adapters/oracle`
  - Satisfies criterion ID: M-03-C04
  - Expected result: Pass. No ERP write-path SQL exists in the adapter boundary.

## QA Steps

- Step: inspect adapter imports from downstream modules.
  - Expected result: only the adapter package imports Oracle driver/query code; business modules remain on `ports.Reader`.
- Step: record exact gap between package tests and live Oracle validation.
  - Expected result: validation artifact clearly separates fake/unit proof from real Oracle proof.

## Rollback/Risk Notes

- Risk: pushing too much source-specific logic into application/domain code.
- Recovery: keep every Oracle-specific concern inside `adapters/oracle` or helper files owned by that boundary.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: implement and record scoped test evidence plus live-validation gaps
- Required files/evidence: spec, changed paths, verification commands, QA steps
- Blockers or open decisions: real Oracle credentials and target validation queries must be operator-approved
