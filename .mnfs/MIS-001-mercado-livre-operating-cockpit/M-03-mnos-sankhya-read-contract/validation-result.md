# Milestone Validation Result

```yaml
id: M-03
type: milestone-validation-result
status: passed
owner: QA Validator
parent: M-03
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: milestone
```

## Milestone ID

M-03-mnos-sankhya-read-contract

## Verdict

Passed

## Validation Scope Declaration

- contract_validated: Yes
- integration_validated: Yes
- live_validation_target: real Oracle on Windows using `godror`, Oracle Instant Client, and a user-scoped portable `gcc`

This pass reflects the Oracle-first rewrite now present in code and the real Oracle validation completed in this session. It is execution truth for both local contract proof and live adapter behavior.

## Scope Reviewed

- `F-01-sankhya-read-contract-import`
- `F-02-read-adapter-seam`
- `F-03-data-quality-rules`

## Accepted Feature Evidence

- F-01
  - Artifacts present: `feature.md`, `spec.md`, `plan.md`, `validation.md`
  - Evidence: focused `go test ./internal/modules/internal_read/...` passed with Oracle-first contract/policy objects and no legacy env-key references
  - Changed paths: `internal_read/domain`, `internal_read/ports`, F-01 validation
- F-02
  - Artifacts present: `feature.md`, `spec.md`, `plan.md`, `validation.md`
  - Evidence: focused `go test ./internal/modules/internal_read/...` passed for application service, fake adapter, Oracle query/mapping/config package, and explicit no-`cgo` failure path
  - Changed paths: `internal_read/application`, `internal_read/adapters/fake`, `internal_read/adapters/oracle`, `apps/server_core/go.mod`, `apps/server_core/go.sum`, F-02 validation
- F-03
  - Artifacts present: `feature.md`, `spec.md`, `plan.md`, `validation.md`
  - Evidence: focused `go test ./internal/modules/internal_read/...` passed with explicit missing-value tests, typed read errors, and nil-preserving guardrails
  - Changed paths: quality flags, nil-preserving tests, fake adapter semantics, F-03 validation, wiki env page

## Criterion Results

### M-03-C01: MPC Owns The Oracle Read Contract

- Verdict: Pass
- Evidence:
  - `apps/server_core/internal/modules/internal_read/domain/contract_test.go`
  - `apps/server_core/internal/modules/internal_read/ports/reader_test.go`
  - Command: `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/...`
- Observed:
  - domain and port types are MPC-owned and Oracle-first
  - no active `MS_DATABASE_URL`, `MS_TENANT_ID`, or `SANKHYA_*` dependency remains inside `internal_read`
  - source semantics and operator policy are separated into typed value objects

### M-03-C02: Sellable Stock And Margin Inputs Are Explicit And Auditable

- Verdict: Pass
- Evidence:
  - `apps/server_core/internal/modules/internal_read/adapters/fake/reader_test.go`
  - `apps/server_core/internal/modules/internal_read/domain/quality_flag_test.go`
  - Command: `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/...`
- Observed:
  - missing product returns unresolved candidate flagged `missing_product`
  - missing stock returns `missing_stock` with nil quantity
  - missing price returns `missing_price` with nil amount
  - missing cost returns `missing_cost` with nil `CUSSEMICM` basis amount
  - missing tax returns `missing_tax` with nil tax numeric fields
  - no missing numeric value is converted to `0`

### M-03-C03: Real Oracle Adapter Works For Claimed Runtime Paths

- Verdict: Pass
- Evidence:
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go`
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/open_cgo.go`
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/open_nocgo.go`
  - temporary `go run` diagnostic using `godror` `PingContext` with normalized `.env`, portable `gcc`, and Oracle Instant Client
  - `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/adapters/oracle -run TestOracleLiveSmoke -v`
- Observed:
  - adapter code now targets Oracle tables/views directly through `database/sql`
  - Oracle ping succeeded against the configured service
  - live smoke passed for product, stock, price, cost, sales, and tax runtime paths
  - real-source semantics found during validation were applied back into code:
    - empty Oracle strings behave as `NULL`
    - repeated effective-date predicates require duplicated bind args with `godror`
    - tax product resolution requires `TGFDIN -> TGFITE` join for `CODPROD`

### M-03-C04: Security And Boundary Discipline Hold

- Verdict: Pass
- Evidence:
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/config_test.go`
  - `rg -n "INSERT|UPDATE|DELETE|MERGE" apps/server_core/internal/modules/internal_read/adapters/oracle`
  - `rg -n "MS_DATABASE_URL|MS_TENANT_ID|SANKHYA_" apps/server_core/internal/modules/internal_read`
- Observed:
  - no ERP write-path SQL exists in the Oracle adapter boundary
  - config validation does not leak secret values in error messages
  - Oracle specifics stay isolated to `adapters/oracle`

## Blocking Failure Check

- Oracle write path
  - Command: `rg -n "INSERT|UPDATE|DELETE|MERGE" apps/server_core/internal/modules/internal_read/adapters/oracle`
  - Result: no matches
- Legacy framing
  - Command: `rg -n "MS_DATABASE_URL|MS_TENANT_ID|SANKHYA_" apps/server_core/internal/modules/internal_read`
  - Result: no matches
- ERP mirror
  - Result: no persistence tables or snapshot schema introduced in M-03

## Commands Run

- `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/...`
  - Result: Pass
  - Output summary:
    - `ok .../adapters/fake`
    - `ok .../adapters/oracle`
    - `ok .../application`
    - `ok .../domain`
    - `ok .../ports`
- `rg -n "MS_DATABASE_URL|MS_TENANT_ID|SANKHYA_" apps/server_core/internal/modules/internal_read`
  - Result: Pass (no matches)
- `rg -n "INSERT|UPDATE|DELETE|MERGE" apps/server_core/internal/modules/internal_read/adapters/oracle`
  - Result: Pass (no matches)
- temporary `go run` diagnostic using `godror` `PingContext` with normalized `.env`, portable `gcc`, and Oracle Instant Client
  - Result: Pass
  - Output summary:
    - `ping-ok`
- `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/adapters/oracle -run TestOracleLiveSmoke -v`
  - Result: Pass
  - Output summary:
    - `PASS: TestOracleLiveSmoke/product_lookup`
    - `PASS: TestOracleLiveSmoke/sellable_stock`
    - `PASS: TestOracleLiveSmoke/current_price`
    - `PASS: TestOracleLiveSmoke/cost_as-of`
    - `PASS: TestOracleLiveSmoke/sales_history`
    - `PASS: TestOracleLiveSmoke/tax_inputs`

## Risks

- The live validation currently depends on a user-scoped portable `gcc` path and normalized `.env` loading in the session shell; this should be codified in operator setup if repeated manually.
- `reader_live_test.go` is a smoke harness, not exhaustive coverage for all Oracle business edge cases.

## Recommendation

Advance M-03 as passed. Keep the Oracle-first implementation, preserve the live smoke harness, and codify the validated Windows setup path (`gcc` + Instant Client + normalized env mapping) in operator guidance if this validation will be repeated outside Codex.
