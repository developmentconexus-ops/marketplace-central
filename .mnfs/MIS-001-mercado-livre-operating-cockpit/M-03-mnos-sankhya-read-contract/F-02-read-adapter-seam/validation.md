# Feature Validation

```yaml
id: F-02
type: feature-validation
status: passed
owner: Feature Implementer
parent: F-02
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-02-oracle-adapter-implementation

## Summary

The MPC-owned Oracle adapter boundary is implemented in code and now validated against a real Oracle runtime on this host. Config loading, DSN assembly, query builders, row mapping, error translation, and `cgo`/`!cgo` runtime split live inside `internal_read/adapters/oracle`, with live smoke proof for product, stock, price, cost, sales, and tax reads.

## Current Validation State

- Result: Passed for live Oracle validation
- Result owner: Feature Implementer
- Decision date: 2026-07-07
- Final feature state for handoff: ready_for_milestone_acceptance

## Evidence

- Command: `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/...`
  - Result: Pass
- Command: `rg -n "INSERT|UPDATE|DELETE|MERGE" apps/server_core/internal/modules/internal_read/adapters/oracle`
  - Result: Pass (no matches)
- Command: temporary `go run` diagnostic using `godror` `PingContext` with normalized `.env`, portable `gcc`, and Oracle Instant Client
  - Result: Pass
  - Output summary: `ping-ok`
- Command: `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/adapters/oracle -run TestOracleLiveSmoke -v`
  - Result: Pass
  - Output summary:
    - `PASS: TestOracleLiveSmoke/product_lookup`
    - `PASS: TestOracleLiveSmoke/sellable_stock`
    - `PASS: TestOracleLiveSmoke/current_price`
    - `PASS: TestOracleLiveSmoke/cost_as-of`
    - `PASS: TestOracleLiveSmoke/sales_history`
    - `PASS: TestOracleLiveSmoke/tax_inputs`

## Observed

- The live-capable adapter uses Oracle views/tables anchored in the current research evidence: `METALPRD.TGFPRO`, `METALPRD.TGFEST`, `METALPRD.TGFCUS`, `LEANDRO.VW_PRECO_TABELA`, `LEANDRO.VW_FAT_VENDA_ITEM`, and `LEANDRO.VW_IMPOSTO_ITEM`.
- `OpenDB` was proven live on Windows using `godror` with `cgo`, a user-scoped portable `gcc`, and Oracle Instant Client.
- Query/mapping/error logic remains isolated to `adapters/oracle`; downstream code stays on `ports.Reader`.
- Real Oracle semantics discovered during validation were folded back into code:
  - Oracle empty-string semantics required `LENGTH(TRIM(REFERENCIA)) > 0` in live discovery.
  - current-price reads require duplicated bind arguments for repeated effective-date predicates under `database/sql`/`godror`.
  - tax reads require joining `LEANDRO.VW_IMPOSTO_ITEM` to `METALPRD.TGFITE` for `CODPROD`.

## Scope Declaration

- contract_validated: Yes
- integration_validated: Yes
- live_validation_target: real Oracle via `godror` + Oracle Instant Client on Windows

- Current status: `passed`
- Next owner: Milestone Orchestrator
- Next action: accept F-02 into milestone closeout and keep the live test harness as the runtime proof path
- Required files/evidence: retained command output for Oracle ping and `TestOracleLiveSmoke`
- Blockers or open decisions: none at feature scope
