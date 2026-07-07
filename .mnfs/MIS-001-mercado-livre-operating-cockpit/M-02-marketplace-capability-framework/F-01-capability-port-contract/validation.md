# Feature Validation

```yaml
id: F-01
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: F-01
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-capability-port-contract

## Summary

Quick validation passed. The `connectors` module now exposes normalized capability DTOs, operation-specific business-facing ports, and an explicit unsupported-behavior resolver without importing Mercado Livre adapter packages into the business-facing capability surface.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-06
- Final feature state for handoff: quick_validation_passed

Feature quick validation is recorded by the Feature Implementer. Milestone Orchestrator accepts, rejects, or blocks the feature output later during milestone integration review. QA Validator owns a feature-level verdict only when Milestone Orchestrator explicitly invokes formal feature validation review.

## Evidence Honesty

- Tag every command, QA step, and artifact with an `Evidence type`: `ran` (executed this session, real output captured), `assumed` (expected but not executed), or `could-not-run` (attempted but blocked; name the reason).
- Record `Pass` only when `Evidence type: ran` and an artifact path or pasted output is present. Never record `Pass` on `assumed` or `could-not-run`.
- A load-bearing check that is `assumed` or `could-not-run` makes the result `Blocked`, not `Pass`.

## Quick Validation State

- fixup_attempts: 0
- max_fixup_attempts: 1
- last_feature_validation_result: all planned checks passed after implementation; PowerShell-based `go test` attempts stalled without output, so final quick validation used `cmd.exe` with a short external `GOCACHE` path.

## Spec Adherence

- Spec satisfied: Yes
- Deviations: None.
- Reason: The feature added normalized connector capability types, operation-specific ports, and explicit unsupported behavior only.

## Changes Made

- File: `apps/server_core/internal/modules/connectors/domain/capability.go`
  Change: Added normalized references, snapshots, stock write types, capability enums, and structured connector errors.
- File: `apps/server_core/internal/modules/connectors/ports/marketplace_capability.go`
  Change: Added business-facing listing, stock read, stock write, and order read ports plus named capability constants.
- File: `apps/server_core/internal/modules/connectors/application/marketplace_capability_service.go`
  Change: Added provider capability registry/service with explicit unsupported operation behavior.
- File: `apps/server_core/internal/modules/connectors/application/marketplace_capability_service_test.go`
  Change: Added fake-implementation tests proving business-facing callers can use the ports and receive explicit unsupported errors.

## Commands Run

- Command: `gofmt -w apps/server_core/internal/modules/connectors/domain/capability.go apps/server_core/internal/modules/connectors/ports/marketplace_capability.go apps/server_core/internal/modules/connectors/application/marketplace_capability_service.go apps/server_core/internal/modules/connectors/application/marketplace_capability_service_test.go`
  Status: Pass
  Evidence type: ran
  Owner: Feature Implementer
  Expected: Connector capability files are formatted successfully.
  Actual: Command exited 0.
  Artifact: workspace files listed above
  Blocking condition: formatting command fails

- Command: `go test ./...`
  Status: Pass
  Evidence type: ran
  Owner: Feature Implementer
  Expected: Connector module tests pass, proving the capability surface compiles and fake implementations work.
  Actual: Command exited 0 from `apps/server_core/internal/modules/connectors` using `cmd.exe` and `GOCACHE=C:\temp\gocache-f01full`.
  Artifact: `apps/server_core/internal/modules/connectors/application/marketplace_capability_service_test.go`
  Blocking condition: connector tests fail or panic

- Command: `rg -n "modules/connectors/adapters/mercado_livre|modules/integrations/adapters/mercadolivre" apps/server_core/internal/modules/connectors/application apps/server_core/internal/modules/connectors/ports apps/server_core/internal/modules/connectors/domain`
  Status: Pass
  Evidence type: ran
  Owner: Feature Implementer
  Expected: No Mercado Livre adapter imports in the business-facing capability surface.
  Actual: Command exited 1 with no matches.
  Artifact: `apps/server_core/internal/modules/connectors/application/marketplace_capability_service.go`
  Blocking condition: adapter package appears in application/domain/ports surface

## Manual QA

- QA level: QA-0
- Flow or step: Source review for provider-agnostic naming and string provider ids.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: Names and value objects stay normalized and do not embed provider payload structs.
- Actual: Added capability types use normalized names such as `ListingSnapshot`, `ProviderListingRef`, `StockWriteRequest`, and `OrderSnapshot`; no Mercado Livre HTTP types were introduced.
- Blocking condition: provider payload types or business policy leak into connector capability contracts

## Evidence

- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/F-01-capability-port-contract/spec.md`
  Status: Pass
  Evidence type: ran
  Owner: Feature Implementer
  Blocking condition: spec missing or lacks milestone-traced criteria

- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/F-01-capability-port-contract/plan.md`
  Status: Pass
  Evidence type: ran
  Owner: Feature Implementer
  Blocking condition: plan missing expected changed paths or verification mapping

- Artifact: `apps/server_core/internal/modules/connectors/application/marketplace_capability_service_test.go`
  Status: Pass
  Evidence type: ran
  Owner: Feature Implementer
  Blocking condition: fake consumer tests missing

## Risks

- The capability surface is intentionally small; pagination and provider-specific metadata remain for later adapter work in F-03.
- `PowerShell`-based `go test` invocations in this session stalled without returning output, so validation standardized on `cmd.exe` for Go commands.

## Milestone Acceptance Review

- Decision: accepted
- Decision owner: Milestone Orchestrator
- Scope check: changed paths stayed inside the `connectors` module and the feature artifact folder.
- Evidence check: all load-bearing criteria rely on `ran` evidence only.
- Notes: this feature does not cross auth/PII/secret or multi-role runtime boundaries, so independent QA Validator review is not required before milestone integration.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Feature Implementer
- Next action: Start F-02 using the accepted connector capability contracts from F-01.
- Required files/evidence: feature brief, spec, plan, changed paths, validation evidence
- Blockers or open decisions: None for F-01.
