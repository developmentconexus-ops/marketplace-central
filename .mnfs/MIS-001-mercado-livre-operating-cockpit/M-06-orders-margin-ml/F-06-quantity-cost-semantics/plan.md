# F-06 Quantity Cost Semantics — Plan

```yaml
id: F-06
type: feature-plan
status: planned
owner: Feature Implementer
parent: M-06
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-0
lifecycle_scope: feature
split_decision: single
split_reason: Two scoped implementation files; the three MNFS artifacts are required protocol evidence, not an additional implementation seam.
```

## Steps

1. Add durable amount-scope source-reference constants for per-unit CUSSEMICM
   and item-line sale-fee/tax inputs.
2. Extend only a non-nil CUSSEMICM amount by the Mercado Livre item quantity
   while mapping the cost input; preserve all existing quality behavior.
3. Add import tests for quantities 1, 2, and 7 that assert extended cost and
   unchanged sale-fee/tax line totals, plus a nil-cost preservation test.
4. Run the registered profitability/internal-read command and the focused
   verbose test command; record exact results in `validation.md`.

## Files Expected To Change

- `apps/server_core/internal/modules/profitability/application/service.go`:
  compose per-unit CUSSEMICM into item-line cost and label scopes.
- `apps/server_core/internal/modules/profitability/application/service_test.go`:
  prove quantity semantics and unknown-cost preservation.
- `F-06-quantity-cost-semantics/spec.md`, `plan.md`, `validation.md`:
  feature specification, plan, and executed evidence.

## Verification Commands

- Command: `cd apps/server_core; $env:GOCACHE="$PWD\\.gocache"; go test ./internal/modules/profitability/application ./internal/modules/internal_read/... -count=1`
  - Satisfies criterion ID: `M-06-C02`
  - Expected result: profitability and internal-read application/unit packages pass.
- Command: `cd apps/server_core; $env:GOCACHE="$PWD\\.gocache"; go test ./internal/modules/profitability/application -run "TestImportMarginInputsExtendsOnlyUnitCostByQuantity|TestImportMarginInputsPreservesUnknownCost" -count=1 -v`
  - Satisfies criterion ID: `M-06-C02`
  - Expected result: quantities 1/2/7 extend only cost; nil cost remains nil
    with missing quality.

## QA Steps

- Review the focused test output for source references proving cost is extended
  from CUSSEMICM per unit while sale fee and tax components stay line totals.

## Rollback/Risk Notes

The only numeric behavior change is cost extension during input import. A
regression can be reverted by restoring the cost mapper call/extension; no
provider, Oracle, or runtime data mutation occurs in this feature.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Implement and record validation evidence.
- Blockers or open decisions: None.
