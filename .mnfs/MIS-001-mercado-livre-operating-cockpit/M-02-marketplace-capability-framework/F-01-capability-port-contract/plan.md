# Feature Plan

```yaml
id: F-01
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-01
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-capability-port-contract

## Steps

1. Add failing tests for normalized capability types, explicit unsupported behavior, and business-facing fake implementations.
2. Add connector domain types, structured connector errors, and operation-specific ports.
3. Add a small application service/registry that resolves provider capability implementations and explicit unsupported behavior.
4. Run targeted connector tests and import-boundary grep checks, then record validation evidence.

## Files Expected To Change

- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/F-01-capability-port-contract/spec.md`
  Reason: Feature specification required by MNFS feature execution.
- Path: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/F-01-capability-port-contract/plan.md`
  Reason: Feature plan required by MNFS feature execution.
- Path: `apps/server_core/internal/modules/connectors/domain/*`
  Reason: New normalized capability references, snapshots, results, and structured errors.
- Path: `apps/server_core/internal/modules/connectors/ports/*`
  Reason: New listing, stock, and order business-facing ports.
- Path: `apps/server_core/internal/modules/connectors/application/*`
  Reason: Small provider capability registry/service and tests proving fake implementations work.

## Verification Commands

- Command: `go test ./internal/modules/connectors/...`
  Satisfies criterion ID: M-02-C01
  Expected result: Connector capability tests pass and business-facing callers compile against ports without provider payload coupling.

- Command: `go test ./internal/modules/connectors/...`
  Satisfies criterion ID: M-02-C02
  Expected result: Unsupported provider operations return `CONNECTORS_PROVIDER_UNSUPPORTED_SHAPE` and tests pass.

- Command: `rg -n "modules/connectors/adapters/mercado_livre|modules/integrations/adapters/mercadolivre" apps/server_core/internal/modules/connectors/application apps/server_core/internal/modules/connectors/ports apps/server_core/internal/modules/connectors/domain`
  Satisfies criterion ID: M-02-C01
  Expected result: No matches in the business-facing capability surface.

## QA Steps

- Step: Review changed connector domain and port files for provider-agnostic naming and string provider ids.
  Expected result: Names and fields stay normalized and do not embed Mercado Livre payload structs.

## Rollback/Risk Notes

- Risk: Capability abstractions grow too broad before real adapters exist.
  Recovery: Keep the first surface limited to IC-001 listing, stock, stock write, and order operations only.
- Risk: Unsupported behavior becomes inconsistent across operations.
  Recovery: Centralize unsupported behavior in one application helper and test it directly.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Implement and record validation evidence.
- Required files/evidence: spec, changed paths, verification commands, QA steps
- Blockers or open decisions: None.
