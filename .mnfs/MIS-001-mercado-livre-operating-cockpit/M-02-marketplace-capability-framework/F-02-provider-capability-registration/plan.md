# Feature Plan

```yaml
id: F-02
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-02
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-02-provider-capability-registration

## Steps

1. Update the marketplace capability profile domain, registries, OpenAPI schema, and SDK type definitions to the business capability vocabulary.
2. Update integration provider `DeclaredCapabilities` so only Mercado Livre advertises the new operational capabilities.
3. Update Go and TypeScript tests that assert capability names and profile shapes.
4. Run targeted Go and SDK test suites and record validation evidence.

## Files Expected To Change

- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/F-02-provider-capability-registration/spec.md`
  Reason: Feature specification required by MNFS feature execution.
- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/F-02-provider-capability-registration/plan.md`
  Reason: Feature plan required by MNFS feature execution.
- Path: `apps/server_core/internal/modules/marketplaces/domain/marketplace_def.go`
  Reason: Public capability profile shape and status vocabulary.
- Path: `apps/server_core/internal/modules/marketplaces/registry/*.go`
  Reason: Marketplace public capability registration.
- Path: `apps/server_core/internal/modules/marketplaces/transport/http_handler_test.go`
  Reason: Public capability profile assertions.
- Path: `apps/server_core/internal/modules/integrations/adapters/*/auth_adapter.go`
  Reason: Provider `DeclaredCapabilities` alignment.
- Path: `apps/server_core/internal/modules/integrations/adapters/providers/registry_test.go`
  Reason: Registry capability expectations.
- Path: `apps/server_core/internal/modules/integrations/application/capability_service_test.go`
  Reason: Resolver tests with the new capability names.
- Path: `contracts/api/marketplace-central.openapi.yaml`
  Reason: Public contract alignment.
- Path: `packages/sdk-runtime/src/index.ts`
  Reason: SDK type alignment.
- Path: `packages/sdk-runtime/src/index.test.ts`
  Reason: SDK fixture expectations.

## Verification Commands

- Command: `go test ./internal/modules/marketplaces/... ./internal/modules/integrations/...`
  Satisfies criterion ID: M-02-C01
  Expected result: marketplace and integrations capability registration tests pass with the new vocabulary.

- Command: `go test ./internal/modules/integrations/...`
  Satisfies criterion ID: M-02-C02
  Expected result: capability resolver tests pass and unsupported operations remain explicit.

- Command: `npm test -- --run packages/sdk-runtime/src/index.test.ts`
  Satisfies criterion ID: M-02-C01
  Expected result: SDK contract tests pass with the updated capability profile shape and declared capabilities.

## QA Steps

- Step: Review marketplace registry values to confirm only Mercado Livre advertises `listing_read`, `stock_read`, `stock_write`, and `order_read`.
  Expected result: future providers remain conservative and do not imply runtime support for unimplemented operations.

## Rollback/Risk Notes

- Risk: updating only one catalog surface causes public/internal drift.
  Recovery: keep OpenAPI, SDK, marketplaces registry, and integration provider definitions in the same change.
- Risk: capability status values drift from IC-001.
  Recovery: normalize all public capability statuses to `supported|unsupported|degraded|blocked`.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Implement and record validation evidence.
- Required files/evidence: spec, changed paths, verification commands, QA steps
- Blockers or open decisions: None.
