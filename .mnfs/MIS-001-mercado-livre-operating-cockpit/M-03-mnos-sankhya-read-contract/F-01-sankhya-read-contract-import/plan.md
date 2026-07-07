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

F-01-sankhya-read-contract-import

## Steps

1. Write failing contract tests for the `internal_read` domain defaults and `Reader` port compilation.
2. Add the minimal `internal_read/domain` value objects and quality-flag constants required by IC-002.
3. Add the `internal_read/ports.Reader` interface and typed inputs aligned with the domain contract.
4. Run focused `go test` for `./internal/modules/internal_read/...` and record the evidence in `validation.md`.

## Files Expected To Change

- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/spec.md`
  - Reason: record the scoped feature contract and MNOS evidence.
- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/plan.md`
  - Reason: record ordered implementation and verification steps.
- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/validation.md`
  - Reason: capture focused validation evidence and handoff state.
- Path: `apps/server_core/internal/modules/internal_read/domain/quality_flag.go`
  - Reason: define shared quality flags for missing and stale inputs.
- Path: `apps/server_core/internal/modules/internal_read/domain/internal_product.go`
  - Reason: define product-link candidate read models.
- Path: `apps/server_core/internal/modules/internal_read/domain/internal_stock.go`
  - Reason: define stock scope defaults and sellable stock model.
- Path: `apps/server_core/internal/modules/internal_read/domain/internal_price.go`
  - Reason: define current price contract with nullable source value.
- Path: `apps/server_core/internal/modules/internal_read/domain/internal_cost.go`
  - Reason: define `CUSSEMICM` cost-as-of contract.
- Path: `apps/server_core/internal/modules/internal_read/domain/internal_tax.go`
  - Reason: define nullable tax inputs and quality flags.
- Path: `apps/server_core/internal/modules/internal_read/domain/internal_sales.go`
  - Reason: define sales history contract with source timestamp.
- Path: `apps/server_core/internal/modules/internal_read/domain/contract_test.go`
  - Reason: assert default scope and required flags.
- Path: `apps/server_core/internal/modules/internal_read/ports/reader.go`
  - Reason: define the application-facing read interface and typed inputs.
- Path: `apps/server_core/internal/modules/internal_read/ports/reader_test.go`
  - Reason: assert the port surface compiles from the consumer boundary.

## Verification Commands

- Command: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/modules/internal_read/...`
  - Satisfies criterion ID: M-03-C01
  - Expected result: Pass. Domain tests confirm formula `SUM(ESTOQUE - RESERVADO)`, companies `[1, 2]`, locations `[10101]`, and scope code `revenda`.
- Command: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/modules/internal_read/...`
  - Satisfies criterion ID: M-03-C02
  - Expected result: Pass. Contract tests confirm the explicit quality flags and the nullable `CUSSEMICM` cost basis.

## QA Steps

- Step: Review the changed paths for scope discipline.
  - Expected result: only feature artifacts, `internal_read/domain`, and `internal_read/ports` change; no Sankhya write path or SQL adapter is introduced.

## Rollback/Risk Notes

- Risk: Over-modeling the future adapter surface would expand Task 1 beyond the contract seam.
- Recovery: keep the module limited to pure types and interface signatures; defer adapters, SQL, and application services to later features.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Implement and record validation evidence.
- Required files/evidence: spec, changed paths, verification commands, QA steps
- Blockers or open decisions: None.
