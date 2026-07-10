# Feature Spec

```yaml
id: F-02
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-02
created: 2026-07-09
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-02-stock-risk-engine

## Problem

M-05 needs a deterministic stock risk classifier before any API, dashboard, or stock action can be safe. The classifier must consume already-read link truth, internal sellable stock, provider stock snapshot, and inventory policy evidence. It must not hide unresolved links, stale sources, or missing quantities behind optimistic recommendations.

## Requirements

- Classify stock risk states: `healthy`, `oversell`, `undersell`, `stale`, `unresolved`, `conflict`, `ineligible`, and `unsupported`.
- Consume resolved product-link truth without importing product-links internals into inventory domain.
- Block unresolved, rejected, missing, and conflict links before quantity recommendation.
- Block stale or missing internal/provider source evidence before quantity recommendation.
- Block ineligible products before quantity recommendation.
- Compare internal recommended quantity against provider announced quantity only after link, source, and eligibility checks pass.
- Produce recommended quantity, policy id, source timestamps, and blocking reason in the risk row.
- Do not call provider synchronously from this feature.

## Non-Goals

- No HTTP API, SDK, UI, persistence, provider calls, or stock writes in F-02.
- No manual action proposal or audit persistence; that belongs to F-03.
- No direct import of `product_links` domain from `inventory/domain`.

## Design

Add a pure `inventory/domain` risk classifier. The classifier accepts a `RiskInput` composed of:

- `LinkEvidence`: inventory-owned representation of upstream link truth.
- `InternalStockEvidence`: internal sellable quantity and source evidence.
- `ProviderStockEvidence`: provider announced quantity and source evidence.
- `StockPolicy`: F-01 policy.
- `ProductEvidence`: eligibility facts already modeled in F-01.

The result is `StockRiskRow`, a future API/readmodel-friendly domain value that carries state, quantities, recommendation, source timestamps, policy id, and blocking reason.

## Edge Cases

- Unresolved or rejected link returns `unresolved` with no recommendation.
- Conflict link returns `conflict` with no recommendation.
- Unsupported provider stock shape returns `unsupported` with no recommendation.
- Missing internal quantity returns `stale` with no recommendation.
- Missing or stale source timestamps return `stale` with no recommendation.
- Ineligible product returns `ineligible` with no recommendation.
- Provider announced quantity greater than recommendation returns `oversell`.
- Provider announced quantity lower than recommendation returns `undersell`.
- Provider announced quantity equal to recommendation returns `healthy`.

## Acceptance Criteria

- M-05-C01: Risk engine uses the F-01 stock policy recommendation formula.
- M-05-C02: Unit tests cover unresolved, conflict, stale, ineligible, unsupported, healthy, oversell, and undersell states.
- M-05-C02: Unit tests prove blocked states do not produce actionable recommendations.
- M-05-C02: Unit tests prove source timestamps and blocking reasons are visible in the risk row.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: implement the planned domain risk classifier
- Required files/evidence: `plan.md`, implementation tests, `validation.md`
- Blockers or open decisions: none for the pure domain classifier
