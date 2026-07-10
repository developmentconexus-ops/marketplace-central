# Feature Spec

```yaml
id: F-03
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-03
created: 2026-07-09
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-03-manual-stock-action-audit

## Problem

Stock Seguro must never write Mercado Livre stock from an unsafe row. A manual apply flow needs an application service that revalidates link truth, source freshness, product eligibility, provider stock shape, and policy recommendation immediately before using a provider stock writer. Every attempt must leave auditable evidence, including skipped/blocked cases where no provider call is made.

## Requirements

- Model manual stock action states: `proposed`, `blocked`, `approved`, `applied`, `failed`, and `skipped`.
- Apply stock only through an inventory-owned `StockWriter` port.
- Reclassify the supplied stock risk input immediately before any provider write.
- Block unresolved, rejected, conflict, stale, ineligible, unsupported, and no-recommendation rows before calling the writer.
- Require explicit manual operator approval before calling the writer.
- Use `stock_action_id` as the provider idempotency key.
- Repeating the same `stock_action_id` must not create a second provider write.
- Audit every result with before quantity, requested quantity, operator, trigger, policy id, source timestamps, provider response summary, idempotency key, and blocking/failure reason when present.
- Persist failed provider responses as failed audit evidence rather than losing the attempt.

## Non-Goals

- No background or automatic stock apply.
- No live Mercado Livre write validation without explicit operator approval.
- No HTTP API, SDK, UI, or Postgres repository in this feature slice.
- No direct dependency on `connectors` domain types inside `inventory/domain`.

## Design

Add inventory-owned action domain types plus an application service. The service owns the orchestration:

1. Check idempotency by reading an existing action from an `ActionStore`.
2. Re-run `ClassifyStockRisk`.
3. Persist a blocked/skipped action when safety gates fail.
4. Persist an approved action before provider intent.
5. Call the inventory `StockWriter` port with `IdempotencyKey == StockActionID`.
6. Persist applied or failed action with provider response summary.

The stock writer port is intentionally inventory-owned. A later adapter can map it to `connectors/ports.StockWriter`.

## Edge Cases

- Missing approval returns `skipped` and does not write.
- Blocked risk state returns `blocked` and does not write.
- Existing action id returns the stored action and does not write again.
- Provider rejection returns `failed` with provider result summary.
- Writer error returns `failed` with error code/message and idempotency key preserved.
- Applied result stores before quantity, requested quantity, timestamps, policy id, operator, and provider summary.

## Acceptance Criteria

- M-05-C02: Tests prove every blocked condition skips the writer.
- M-05-C03: Tests prove applied and failed actions store before quantity, requested quantity, manual operator trigger, policy id, source timestamps, provider response summary, and idempotency key.
- M-05-C03: Tests prove duplicate action id does not create duplicate provider intent.
- M-05-C03: Tests prove unapproved manual action is skipped with audit evidence.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: implement the planned application slice
- Required files/evidence: `plan.md`, implementation tests, `validation.md`
- Blockers or open decisions: live provider write validation still requires explicit operator approval
