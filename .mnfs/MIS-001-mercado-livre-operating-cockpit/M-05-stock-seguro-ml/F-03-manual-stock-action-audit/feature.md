# F-03-manual-stock-action-audit

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Brief

Create proposed manual stock actions and apply them through the capability `StockWriter` only after re-validating link, source freshness, policy, and provider state.

## Inputs

- Stock risk rows.
- Resolved product link.
- Mercado Livre `StockWriter`.
- Operator manual approval.

## Expected Output

- Action states: `proposed`, `blocked`, `approved`, `applied`, `failed`, `skipped`.
- Complete audit for every apply attempt.

## Constraints

- No background/automatic apply.
- Revalidate before provider write.
- Idempotency key is the stock action id.

## Negative Scenarios

- Stale source blocks action.
- Provider rejection persists failed audit with response summary.
- Repeating same action id does not create duplicate provider intent.

## Validation Expectations

- Tests prove every blocked condition skips provider adapter.
- Tests prove applied/failed action audit contains required fields.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-03/validation.md.
- Blockers or open decisions: live provider write needs explicit operator approval.
