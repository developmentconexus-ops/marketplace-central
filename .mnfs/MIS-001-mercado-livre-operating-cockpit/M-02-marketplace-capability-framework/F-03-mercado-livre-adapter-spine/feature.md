# F-03-mercado-livre-adapter-spine

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-02
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-02 Marketplace capability framework.

## Brief

Create the Mercado Livre adapter spine using direct HTTP seams/stubs and official docs for item, variation stock, and order shapes. Do not depend on deprecated official Go SDK.

## Inputs

- IC-001.
- Context7 Mercado Livre docs.
- Existing Mercado Livre OAuth installation support.

## Expected Output

- Adapter maps documented listing/stock/order fields into normalized snapshots.
- Stock write operation accepts an idempotent action id and returns structured result.

## Constraints

- No real provider write tests without explicit operator-controlled credentials/listings.
- Business stock policy remains outside adapter.

## Inputs/Outputs

- Input: provider listing/order refs and credential context.
- Output: normalized listing, stock, order, and stock write results from IC-001.

## Negative Scenarios

- Unsupported variation/distributed stock shape returns `CONNECTORS_PROVIDER_UNSUPPORTED_SHAPE`.
- Provider 429 maps to `CONNECTORS_PROVIDER_RATE_LIMITED`.

## Validation Expectations

- Unit tests cover item with no variation, item with variation, order with `sale_fee`, and provider rejection.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-03/validation.md.
- Blockers or open decisions: Real live write validation needs operator approval.
