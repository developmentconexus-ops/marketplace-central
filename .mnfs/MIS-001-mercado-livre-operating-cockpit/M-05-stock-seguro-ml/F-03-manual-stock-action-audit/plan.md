# Feature Plan

```yaml
id: F-03
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-03
created: 2026-07-09
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: feature
split_decision: single
split_reason: application service plus ports and focused fake tests, no transport or persistence layer in this slice
```

## Steps

1. Add failing application tests for blocked states, unapproved skip, applied audit, provider failure audit, and duplicate idempotency.
2. Add inventory domain action evidence types.
3. Add inventory-owned ports for stock writer and action store.
4. Implement manual stock action application service.
5. Run focused inventory validation and regression validation.
6. Record evidence in `validation.md`.

## Expected Changed Paths

- `apps/server_core/internal/modules/inventory/domain/action.go`
- `apps/server_core/internal/modules/inventory/ports/stock_action.go`
- `apps/server_core/internal/modules/inventory/application/stock_action_service.go`
- `apps/server_core/internal/modules/inventory/application/stock_action_service_test.go`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-03-manual-stock-action-audit/spec.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-03-manual-stock-action-audit/plan.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-03-manual-stock-action-audit/validation.md`

## Verification Commands

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/inventory/... -count=1`
  - Satisfies criterion IDs: M-05-C02, M-05-C03
  - Expected result: Pass. Action safety gates, idempotency, and audit evidence are covered.
- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/internal_read/... ./apps/server_core/internal/modules/inventory/... -count=1`
  - Satisfies criterion IDs: M-05-C02, M-05-C03
  - Expected result: Pass. Existing internal-read compatibility remains intact.

## Manual QA

- No browser QA is required for F-03 because this slice has no UI surface.
- Live provider write QA is intentionally not run in this feature without explicit operator approval.

## Rollback/Risk Notes

- Risk: provider responses become raw payload dependencies.
- Mitigation: store only inventory-owned provider response summary fields in the action audit.
- Risk: idempotency is only as strong as the store implementation.
- Mitigation: service treats existing action id as terminal for provider intent; Postgres uniqueness is planned for the persistence slice.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: run the TDD implementation steps
- Required files/evidence: focused Go test output and `validation.md`
- Blockers or open decisions: live provider write validation requires explicit operator approval
