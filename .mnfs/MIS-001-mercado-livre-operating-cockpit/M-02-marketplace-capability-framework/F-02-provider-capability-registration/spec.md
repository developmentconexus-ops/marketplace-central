# Feature Spec

```yaml
id: F-02
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-02
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-02-provider-capability-registration

## Problem

Marketplace Central still exposes legacy marketplace capability names such as `publish`, `stock_sync`, and `orders`, while M-02 and IC-001 require business capability names for listing read, stock read, stock write, and order read. The provider catalogs must align so UI and future business services can reason about Mercado Livre support without inferring unsupported providers as operational.

## Requirements

- Requirement: Replace the public marketplace capability profile fields and enums with business capability names and IC-001 status vocabulary.
  Acceptance evidence: Go tests, TypeScript tests, and contract inspection show `listing_read`, `stock_read`, `stock_write`, and `order_read` with `supported|unsupported|degraded|blocked`.
- Requirement: Update provider definitions so Mercado Livre declares the four business capabilities and future providers do not imply unsupported operations.
  Acceptance evidence: Provider registry tests assert Mercado Livre support and conservative capability declarations for other providers.
- Requirement: Capability state resolution tests use the new business capability names.
  Acceptance evidence: Go capability service tests pass with the new names.

## Non-Goals

- Implement operational reads/writes against Mercado Livre APIs.
- Add new frontend screens or change auth lifecycle behavior.
- Model shipment/question operational capabilities beyond conservative placeholders.

## Design

Update the public marketplace profile shape in `marketplaces/domain`, registries, OpenAPI, and `packages/sdk-runtime` to the business capability vocabulary required by M-02. Keep `messages`, `questions`, `freight_quotes`, `webhooks`, and `sandbox` as secondary fields, but move their status enum to the same `supported|unsupported|degraded|blocked` vocabulary.

Update integration provider definitions to declare only capabilities that are truly executable today. Mercado Livre declares `listing_read`, `stock_read`, `stock_write`, and `order_read`. Other providers keep only capabilities already grounded in existing runtime behavior such as `pricing_fee_sync`, and do not advertise the new business operations.

## Edge Cases

- Public and internal provider catalogs must change together; otherwise SDK/UI drift breaks capability reads.
- Existing capability state resolution for persisted values must accept the new names without migration.
- Future providers may remain active in auth/catalog while operational capabilities stay absent or blocked.

## Acceptance Criteria

- Criterion: The public and internal provider contracts expose business capability names for the M-02 operations without implying unsupported providers are operational.
  Traces to milestone criterion ID: M-02-C01
  Proven by (verification command or QA step): `go test ./internal/modules/marketplaces/... ./internal/modules/integrations/...`

- Criterion: Unsupported capability support remains explicit in provider definitions and resolver tests.
  Traces to milestone criterion ID: M-02-C02
  Proven by (verification command or QA step): `go test ./internal/modules/integrations/...`

- Criterion: SDK and OpenAPI reflect the new marketplace capability profile vocabulary.
  Traces to milestone criterion ID: M-02-C01
  Proven by (verification command or QA step): `npm test -- --run packages/sdk-runtime/src/index.test.ts`

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Write `plan.md` and implement the scoped feature.
- Required files/evidence: feature brief, spec, milestone contract, validation expectations
- Blockers or open decisions: None.
