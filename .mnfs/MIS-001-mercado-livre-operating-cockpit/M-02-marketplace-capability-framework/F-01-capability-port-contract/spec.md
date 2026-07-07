# Feature Spec

```yaml
id: F-01
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-01
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-capability-port-contract

## Problem

MPC needs business-facing capability contracts for listing, stock, and order operations before Mercado Livre-specific adapter work begins. The current codebase has provider metadata and integration capability state, but it does not yet expose small business-facing ports or normalized snapshots that later modules can depend on without importing provider payloads.

## Requirements

- Requirement: Add normalized marketplace capability types for listing, stock, stock write, and order operations under the `connectors` module.
  Acceptance evidence: Go unit tests compile and assert normalized fields and structured errors through fake implementations.
- Requirement: Expose small ports for listing read, stock read/write, and order read operations with provider ids kept as strings.
  Acceptance evidence: Unit tests compile against fake ports and do not import Mercado Livre adapter packages.
- Requirement: Unsupported capability behavior is explicit and stable for later features to reuse.
  Acceptance evidence: Tests verify unsupported capability returns `CONNECTORS_PROVIDER_UNSUPPORTED_SHAPE` instead of nil success or panic.

## Non-Goals

- Implement Mercado Livre HTTP requests or payload mapping.
- Change OpenAPI or frontend provider catalog capability fields.
- Implement Stock Seguro policy, product links, or order profitability logic.

## Design

Create a small domain surface inside `apps/server_core/internal/modules/connectors/domain` with:

- provider/account/listing/order reference value types;
- normalized listing, stock, stock-write, and order snapshots/results;
- connector error codes and a typed error wrapper for structured provider failures.

Create small operation-specific ports in `apps/server_core/internal/modules/connectors/ports`:

- `ListingReader`
- `StockReader`
- `StockWriter`
- `OrderReader`

Create a small application-level registry/service in `apps/server_core/internal/modules/connectors/application` that lets business-facing callers depend on the ports and resolves explicit unsupported behavior when an operation is absent for a provider.

## Edge Cases

- Provider code or required ids are blank: return a structured invalid/unsupported error rather than nil behavior.
- Provider supports listing reads but not stock writes: stock write lookup must return `CONNECTORS_PROVIDER_UNSUPPORTED_SHAPE`.
- Business-facing caller uses fake implementations in tests: no provider adapter import is required.

## Acceptance Criteria

- Criterion: Business-facing listing, stock, and order capabilities are represented as small ports with provider ids as strings and no Mercado Livre endpoint payloads in business modules.
  Traces to milestone criterion ID: M-02-C01
  Proven by (verification command or QA step): `go test ./internal/modules/connectors/...`

- Criterion: Unsupported provider operations return explicit unsupported behavior instead of nil success.
  Traces to milestone criterion ID: M-02-C02
  Proven by (verification command or QA step): `go test ./internal/modules/connectors/...`

- Criterion: No business-facing application code imports Mercado Livre adapter packages for the capability surface.
  Traces to milestone criterion ID: M-02-C01
  Proven by (verification command or QA step): `rg -n "modules/connectors/adapters/mercado_livre|modules/integrations/adapters/mercadolivre" apps/server_core/internal/modules/connectors/application apps/server_core/internal/modules/connectors/ports apps/server_core/internal/modules/connectors/domain`

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Write `plan.md` and implement the scoped feature.
- Required files/evidence: feature brief, spec, milestone contract, validation expectations
- Blockers or open decisions: None.
