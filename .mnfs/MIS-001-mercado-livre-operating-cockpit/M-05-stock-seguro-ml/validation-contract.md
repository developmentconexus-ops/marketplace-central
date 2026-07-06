# Milestone Validation Contract

```yaml
id: M-05
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: milestone
```

## Milestone ID

M-05

## QA Level

QA-0 planning contract.

## Required Outcome

Stock Seguro classifies linked listing stock risk and applies only manual audited stock actions.

## Criteria

## Criterion: Safe Quantity Formula
ID: M-05-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: inventory policy tests
- Expected: Default recommended Mercado Livre quantity is `max(0, SUM(ESTOQUE - RESERVADO) - 1)` for `CODEMP IN (1,2)` and `CODLOCAL=10101`.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/validation-result.md`
Blocking failure: Recommendation includes showroom/reserved stock or omits default buffer.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Blocked Unsafe Actions
ID: M-05-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: stock action service tests
- Expected: Unresolved link, conflict link, stale internal source, stale provider source, ineligible product, and unsupported provider shape produce blocked actions and no provider write.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/validation-result.md`
Blocking failure: Provider write is attempted for a blocked condition.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Manual Audit Evidence
ID: M-05-C03
Level: Milestone
Type: Observability
Required: Yes
Status: Pending
Evidence:
- Command: stock action repository/API tests
- Expected: Applied/failed action stores before quantity, requested quantity, operator/manual trigger, policy id, source timestamps, provider response summary, and idempotency key.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/validation-result.md`
Blocking failure: Manual write can be applied without complete audit evidence.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Stock Dashboard States
ID: M-05-C04
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: frontend tests and browser validation
- Expected: Dashboard shows healthy, oversell, undersell, stale, unresolved, conflict, ineligible, loading, error, and empty states.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/validation-result.md`
Blocking failure: Operator cannot distinguish blocked from actionable stock rows.
Blocking failure observed: No
Drive (UI — agent-browser; UI criteria only; omit for non-UI):
- Fixture: seeded stock risks for healthy, oversell, stale, unresolved, and conflict rows.
- Steps:
  - open `/inventory/stock-seguro`
  - assert text "Stock Seguro"
  - assert text "Oversell"
  - assert text "Conflito"
  - assert text "Sem vinculo"
- Expected: Risk rows and blocked states are visible in the first viewport or filterable table.
Owner: QA Validator

## Evidence Requirements

- API/OpenAPI/SDK evidence for stock risk/action endpoints.
- UI browser evidence for main workflow.
- Provider write tests use fake adapter unless operator approves live listing write.

## Blocking Failures

- Automatic stock writes.
- Stock action without resolved link.
- Missing source timestamp/policy in audit.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned.
- Next owner: QA Validator after execution.
- Next action: Validate Stock Seguro end to end.
- Required files/evidence: F-*/validation.md.
- Blockers or open decisions: Live write validation requires explicit operator approval.
