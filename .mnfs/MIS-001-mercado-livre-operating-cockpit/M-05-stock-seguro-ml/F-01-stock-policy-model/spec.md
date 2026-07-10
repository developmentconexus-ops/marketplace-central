# Feature Spec

```yaml
id: F-01
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-01
created: 2026-07-09
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-stock-policy-model

## Problem

M-05 needs Stock Seguro decisions that are owned by the inventory module, not by Oracle query code or UI logic. The current internal read layer can query sellable stock using mission defaults, but it does not model the downstream safety policy that classifies source freshness, buffers, eligibility, and the non-negative recommended Mercado Livre quantity.

## Requirements

- Create `inventory` domain value objects for stock scope, safety buffer, freshness, eligibility, and policy evidence.
- Preserve the mission default scope: `CODEMP IN (1,2)`, `CODLOCAL=10101`, excluded showroom location `10108`, and formula `SUM(ESTOQUE - RESERVADO)`.
- Compute recommended Mercado Livre quantity as `max(0, internal_sellable_quantity - buffer)`.
- Default buffer is `1`.
- Freshness policy must mark missing timestamps and sources older than the threshold as stale.
- Product exclusion rules must be explicit, visible, and extensible for future group, weight, size, and margin rules.
- Provide an explicit mapping from inventory policy to the existing `internal_read` sellable stock input policy so the business rule remains owned by `inventory`.

## Non-Goals

- No stock risk list, API, SDK, UI, or dashboard in F-01.
- No provider read/write or Mercado Livre `available_quantity` mutation in F-01.
- No automatic stock action policy in this mission.
- No final product group, weight, size, or margin exclusion examples until owner examples are available.

## Design

The `inventory/domain` package owns `StockPolicy`, `StockScope`, `StockBuffer`, `FreshnessPolicy`, `EligibilityPolicy`, and policy evidence types. The policy is pure domain code and can be tested without database or provider dependencies.

The `inventory/adapters/internalread` package translates the domain policy to `internal_read/domain.SellableStockPolicy` and `internal_read/domain.FreshnessPolicy`. This keeps Oracle query shape at the adapter boundary while allowing F-02 and F-03 to consume inventory-owned concepts.

## Edge Cases

- Negative or zero internal stock recommends `0`.
- Buffer larger than internal stock recommends `0`.
- Zero buffer is allowed but does not increase the recommendation above internal sellable stock.
- Missing source timestamp is stale.
- Exact threshold age is still fresh; older than threshold is stale.
- Eligibility defaults to eligible when no explicit exclusion rules match.

## Acceptance Criteria

- M-05-C01: Unit tests prove the default policy preserves company/location/showroom/formula defaults and computes `max(0, quantity - 1)`.
- M-05-C01: Unit tests prove recommendation never becomes negative.
- M-05-C01: Unit tests prove the inventory policy maps to `internal_read` sellable stock policy without losing mission defaults.
- M-05-C02: Unit tests prove stale source detection is explicit for missing and old timestamps.
- M-05-C02: Unit tests prove product ineligibility rules are explicit and produce visible reasons.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: implement the planned domain slice
- Required files/evidence: `plan.md`, implementation tests, `validation.md`
- Blockers or open decisions: concrete group/weight/size/margin exclusion examples remain deferred to later owner refinement
