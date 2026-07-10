# Feature Plan

```yaml
id: F-01
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-01
created: 2026-07-09
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: feature
split_decision: single
split_reason: pure domain plus one adapter mapping, within four expected changed implementation paths
```

## Steps

1. Add failing domain tests for default policy evidence, non-negative recommendations, freshness decisions, and explicit eligibility reasons.
2. Implement `inventory/domain` value objects and pure policy behavior.
3. Add failing adapter tests for mapping inventory policy to `internal_read` stock/freshness policy.
4. Implement `inventory/adapters/internalread` mapping.
5. Run focused Go validation and record evidence in `validation.md`.

## Expected Changed Paths

- `apps/server_core/internal/modules/inventory/domain/policy.go`
- `apps/server_core/internal/modules/inventory/domain/policy_test.go`
- `apps/server_core/internal/modules/inventory/adapters/internalread/policy_mapper.go`
- `apps/server_core/internal/modules/inventory/adapters/internalread/policy_mapper_test.go`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-01-stock-policy-model/spec.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-01-stock-policy-model/plan.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-01-stock-policy-model/validation.md`

## Verification Commands

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/inventory/... -count=1`
  - Satisfies criterion IDs: M-05-C01, M-05-C02
  - Expected result: Pass. Inventory policy domain and internal-read mapping are covered.
- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/internal_read/... ./apps/server_core/internal/modules/inventory/... -count=1`
  - Satisfies criterion IDs: M-05-C01, M-05-C02
  - Expected result: Pass. Existing internal read stock contract remains compatible with the inventory policy mapping.

## Manual QA

No browser QA is required for F-01 because this feature has no UI surface. Browser validation starts in F-04.

## Rollback/Risk Notes

- Risk: duplicating the internal-read stock policy can create drift.
- Mitigation: keep a single explicit adapter mapping test that proves the inventory-owned default maps to the current `internal_read` query contract.
- Risk: product eligibility rules are under-specified.
- Mitigation: implement a generic explicit reason model now and defer concrete group/weight/size/margin examples until owner examples exist.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: run the TDD implementation steps
- Required files/evidence: focused Go test output and `validation.md`
- Blockers or open decisions: none for the F-01 foundation slice
