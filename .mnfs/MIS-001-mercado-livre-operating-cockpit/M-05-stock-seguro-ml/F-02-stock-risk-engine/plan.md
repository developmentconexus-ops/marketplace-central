# Feature Plan

```yaml
id: F-02
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-02
created: 2026-07-09
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: feature
split_decision: single
split_reason: pure domain classifier with one expected implementation path and no API/persistence/UI work
```

## Steps

1. Add failing domain tests for every required risk state and blocked recommendation behavior.
2. Implement `inventory/domain` risk types and classifier.
3. Run focused inventory validation.
4. Run inventory plus internal-read regression validation.
5. Record evidence in `validation.md`.

## Expected Changed Paths

- `apps/server_core/internal/modules/inventory/domain/risk.go`
- `apps/server_core/internal/modules/inventory/domain/risk_test.go`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-02-stock-risk-engine/spec.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-02-stock-risk-engine/plan.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-02-stock-risk-engine/validation.md`

## Verification Commands

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/inventory/... -count=1`
  - Satisfies criterion IDs: M-05-C01, M-05-C02
  - Expected result: Pass. Risk classifier and F-01 policy remain green together.
- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/internal_read/... ./apps/server_core/internal/modules/inventory/... -count=1`
  - Satisfies criterion IDs: M-05-C01, M-05-C02
  - Expected result: Pass. Existing internal-read contract remains compatible with inventory policy and risk domain.

## Manual QA

No browser QA is required for F-02 because this feature has no UI surface. Browser validation starts in F-04.

## Rollback/Risk Notes

- Risk: importing `product_links` domain directly would couple inventory to another module's internals.
- Mitigation: use inventory-owned link evidence types and map from product links later in application code.
- Risk: blocked rows accidentally carry recommendations.
- Mitigation: tests assert no recommendation for all blocked states.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: run the TDD implementation steps
- Required files/evidence: focused Go test output and `validation.md`
- Blockers or open decisions: none for the F-02 domain slice
