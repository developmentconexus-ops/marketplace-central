# Feature Validation

```yaml
id: F-01
type: feature-validation
status: passed
owner: Feature Implementer
parent: F-01
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-oracle-read-contract-redesign

## Summary

The Oracle-first contract rewrite is implemented and locally verified. `internal_read/domain` and `internal_read/ports` now expose MPC-owned policy/value objects, typed quality states, and typed read errors without any active `MS_DATABASE_URL` or `SANKHYA_*` contract dependence.

## Current Validation State

- Result: Passed for contract ownership and local proof
- Result owner: Feature Implementer
- Decision date: 2026-07-07
- Final feature state for handoff: ready_for_oracle_adapter_consumers

## Evidence

- Command: `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/...`
  - Result: Pass
- Command: `rg -n "MS_DATABASE_URL|MS_TENANT_ID|SANKHYA_" apps/server_core/internal/modules/internal_read`
  - Result: Pass (no matches)

## Observed

- `SellableStockPolicy` now owns the default `CODEMP IN (1,2)`, `CODLOCAL=10101`, `CODLOCAL=10108` exclusion, and `SUM(ESTOQUE - RESERVADO)` semantics.
- `CurrentPricePolicy`, `CostAsOfPolicy`, `TaxPolicy`, and `SalesHistoryWindow` move policy/time semantics into MPC-owned types.
- `ReadError` now distinguishes `source_unavailable` from `unsupported_query`.
- Missing price/cost/tax/stock remain nil-valued outputs with explicit quality states, never zero-filled defaults.

## Scope Declaration

- contract_validated: Yes
- integration_validated: No
- blocked_for_real_validation: F-01 does not claim live Oracle behavior on its own

## Handoff

- Current status: `passed`
- Next owner: Feature Implementer
- Next action: keep F-02/F-03 aligned with the rewritten contract during live Oracle validation
- Required files/evidence: F-02 and milestone live-validation evidence
- Blockers or open decisions: none for the contract layer beyond live-consumer exercise
