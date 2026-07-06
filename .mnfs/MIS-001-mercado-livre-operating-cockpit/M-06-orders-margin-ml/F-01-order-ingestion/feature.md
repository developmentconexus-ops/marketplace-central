# F-01-order-ingestion

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-06
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Brief

Ingest Mercado Livre orders, items, payments, shipment ids, tags, and cancellation details through capability ports.

## Inputs

- IC-001 ListOrders/ReadOrder.
- Mercado Livre official docs for orders.

## Expected Output

- Tenant-scoped `MarketplaceOrder` and item snapshots.
- Missing product link creates quality flag but does not block ingestion.

## Constraints

- Do not calculate margin in the adapter.
- Do not expose unnecessary buyer PII.

## Inputs/Outputs

- Input: provider account, date/status cursor, provider order id.
- Output: normalized order snapshot with items, quantities, unit price, sale fee, payments, shipping id, cancellation detail.

## Negative Scenarios

- Duplicate order update is idempotent.
- Unknown provider field shape stores raw reference only where allowed and maps known fields safely.

## Validation Expectations

- Tests cover paid order, canceled order, order with variation, and order with missing product link.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-01/validation.md.
- Blockers or open decisions: None.
