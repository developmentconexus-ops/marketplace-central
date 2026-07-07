# Feature Validation

```yaml
id: F-03
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: F-03
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-03-data-quality-rules

## Quick Validation

- Command: `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/...`
  - Evidence type: ran
  - Expected: fake seam proves missing facts stay flagged and missing numerics stay nil.
  - Actual: `ok .../adapters/fake`; `ok .../adapters/oracle`; `ok .../application`; `ok .../domain`; `ok .../ports`

- Command: `rg -n "CUSVARIAVEL|INSERT|UPDATE|DELETE" apps/server_core/internal/modules/internal_read`
  - Evidence type: ran
  - Expected: no `CUSVARIAVEL`; no write-path SQL in the internal read seam.
  - Actual: no matches

- Command: `rg -n "SANKHYA_DSN|SANKHYA_PASSWORD|password=" .mnfs apps/server_core/internal/modules/internal_read`
  - Evidence type: ran
  - Expected: no secret values leaked in artifacts or module errors.
  - Actual: only environment-variable names in `adapters/oracle/config.go` and `config_test.go`; no secret values leaked.

## Summary

Quick validation passed for reusable quality-flag semantics. The fake seam now keeps missing product, stock, cost, and tax as explicit quality states, and missing numerics remain `nil` rather than `0`.
