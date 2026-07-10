# Feature Validation

```yaml
id: F-01
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: F-01
created: 2026-07-09
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-stock-policy-model

## Quick Validation Result

Passed.

## Spec Adherence

- M-05-C01: Default policy preserves `CODEMP IN (1,2)`, `CODLOCAL=10101`, showroom exclusion `10108`, formula `SUM(ESTOQUE - RESERVADO)`, and buffer `1`.
- M-05-C01: Recommended provider quantity is `max(0, internal_sellable_quantity - buffer)` and never goes negative.
- M-05-C01: Inventory-owned policy maps to the existing `internal_read` sellable stock policy without losing mission defaults.
- M-05-C02: Freshness evaluation marks missing and older-than-threshold source timestamps as stale.
- M-05-C02: Eligibility rules return explicit ineligible reasons and support future rule types through a generic rule model.

## Changed Paths

- `apps/server_core/internal/modules/inventory/domain/policy.go`
- `apps/server_core/internal/modules/inventory/domain/policy_test.go`
- `apps/server_core/internal/modules/inventory/adapters/internalread/policy_mapper.go`
- `apps/server_core/internal/modules/inventory/adapters/internalread/policy_mapper_test.go`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-01-stock-policy-model/spec.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-01-stock-policy-model/plan.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-01-stock-policy-model/implementation-report.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-01-stock-policy-model/validation.md`

## Commands Run

### RED

- Evidence type: `ran`
- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/inventory/... -count=1`
- Expected: Fail because the new inventory policy API did not exist yet.
- Actual: Failed with undefined `DefaultStockPolicy`, `FreshnessPolicy`, `SourceEvidence`, `EligibilityPolicy`, and `MapStockPolicy`.
- Result: Pass as RED evidence.

### REVIEW FIXUP RED

- Evidence type: `ran`
- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/inventory/... -count=1`
- Expected: Fail because the generic eligibility API did not exist yet.
- Actual: Failed with undefined `EligibilityRule`, `GroupExclusion`, `EligibilityRuleTypeWeight`, `WeightGrams`, and `RuleType`.
- Result: Pass as RED evidence.

### GREEN

- Evidence type: `ran`
- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/inventory/... -count=1`
- Expected: Pass.
- Actual: Passed.
- Result: Pass.

### REGRESSION

- Evidence type: `ran`
- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/internal_read/... ./apps/server_core/internal/modules/inventory/... -count=1`
- Expected: Pass.
- Actual: Passed.
- Result: Pass.

## Manual QA

- Evidence type: `could-not-run`
- Step: Browser validation for Stock Seguro dashboard.
- Result: Not applicable to F-01 because this feature has no UI surface. Browser QA is required for F-04.

## Review Evidence

- Evidence type: `ran`
- Reviewer: Copernicus
- Result: Initial review found Important issues for default freshness, eligibility extensibility, and negative buffer handling. All were fixed with test-first changes and revalidated.

## Same-Session Fixups

- Attempt 1:
  - Reproduction: reviewer findings were converted into failing tests.
  - Change: added default freshness, generic eligibility rules, and negative-buffer protection.
  - Result: focused and regression validation passed.

## Risks

- Concrete product group, weight, size, and margin exclusion examples remain unavailable. The model supports these rule types, but final values require owner input in later features.
- This feature proves policy behavior only. It does not prove risk list generation, provider stock writes, API/SDK contracts, or UI behavior.

## Handoff

- Current status: `quick_validation_passed`
- Handoff target: Milestone Orchestrator
- Handoff reason: F-01 policy foundation is ready for feature acceptance review and F-02 risk engine planning.
- Required next inputs: source freshness threshold confirmation can be refined if operators want a value other than the current default `30m`.
