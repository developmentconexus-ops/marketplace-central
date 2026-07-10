# M-06 F-02 Correction Review Package

## Requirements

- `.superpowers/sdd/m06-f02-correction-brief.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-02-margin-input-model/spec.md`

## Evidence

- `.superpowers/sdd/m06-f02-correction-report.md`

## Current-State Files

- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`
- `apps/server_core/internal/modules/profitability/adapters/postgres/store.go`
- `apps/server_core/internal/modules/profitability/adapters/postgres/manual_adjustment_integration_test.go`
- `apps/server_core/migrations/0028_profitability_margin_inputs.sql`
- `apps/server_core/migrations/0030_profitability_manual_adjustment_invariants.sql`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`
- `packages/sdk-runtime/src/index.test.ts`

## Binding Constraints

- Manual adjustments are append-only audited events.
- Actor type/id and reason are required and persisted.
- Scope/category and item identity cannot be ambiguous.
- Freight and commission both work and retain signed amounts.
- IDs remain unique when timestamps are equal.
- OpenAPI and SDK agree that actor is required.
- Unknown margin inputs remain nil; F-03 calculation behavior is out of scope.
