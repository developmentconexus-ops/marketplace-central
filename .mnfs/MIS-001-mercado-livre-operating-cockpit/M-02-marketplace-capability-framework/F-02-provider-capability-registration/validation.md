# Feature Validation

```yaml
id: F-02
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: F-02
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-02-provider-capability-registration

## Summary

Quick validation passed. Marketplace capability profiles, integration provider declarations, OpenAPI, and SDK runtime now use the M-02 business capability vocabulary for `listing_read`, `stock_read`, `stock_write`, and `order_read`, with conservative support declared only for Mercado Livre.

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
- last_feature_validation_result: planned checks passed on first run after contract and fixture alignment.

## Spec Adherence

- Spec satisfied: Yes
- Deviations: None.
- Reason: The feature updated the public/internal capability vocabulary and conservative provider declarations without touching runtime adapter behavior.

## Changes Made

- File: `apps/server_core/internal/modules/marketplaces/domain/marketplace_def.go`
  Change: Replaced legacy marketplace capability names and status vocabulary with business capability names and IC-001-aligned status values.
- File: `apps/server_core/internal/modules/marketplaces/registry/*.go`
  Change: Updated marketplace public capability profiles, with Mercado Livre marked supported for the four business capabilities and future providers kept blocked/unsupported.
- File: `apps/server_core/internal/modules/marketplaces/transport/http_handler_test.go`
  Change: Updated public capability profile assertions.
- File: `apps/server_core/internal/modules/integrations/adapters/*/auth_adapter.go`
  Change: Updated provider `DeclaredCapabilities` so only Mercado Livre declares the new operational capabilities and future providers stay conservative.
- File: `apps/server_core/internal/modules/integrations/application/capability_service_test.go`
  Change: Updated resolver tests to the new capability names.
- File: `apps/server_core/internal/modules/integrations/adapters/providers/registry_test.go`
  Change: Added capability expectations for Mercado Livre and conservative declarations for future providers.
- File: `apps/server_core/internal/modules/integrations/adapters/postgres/provider_definition_repo_test.go`
  Change: Updated stored capability fixture vocabulary.
- File: `contracts/api/marketplace-central.openapi.yaml`
  Change: Updated `CapabilityProfile` schema field names and enums.
- File: `packages/sdk-runtime/src/index.ts`
  Change: Updated SDK capability profile types and enums.
- File: `packages/sdk-runtime/src/index.test.ts`
  Change: Updated provider declaration fixtures to the conservative vocabulary.

## Commands Run

- Command: `go test ./internal/modules/marketplaces/... ./internal/modules/integrations/...`
  Status: Pass
  Evidence type: ran
  Owner: Feature Implementer
  Expected: marketplace and integrations capability registration/resolution tests pass with the new vocabulary.
  Actual: Command exited 0 from `apps/server_core` using `cmd.exe` and `GOCACHE=C:\temp\gocache-f02`.
  Artifact: changed Go files listed above
  Blocking condition: capability registration or resolver tests fail

- Command: `npm run test --workspace @marketplace-central/sdk-runtime -- --run src/index.test.ts`
  Status: Pass
  Evidence type: ran
  Owner: Feature Implementer
  Expected: SDK runtime tests pass with the updated capability vocabulary.
  Actual: `src/index.test.ts` passed with `16` tests green in `24.39s`.
  Artifact: `packages/sdk-runtime/src/index.test.ts`
  Blocking condition: SDK contract tests fail

## Manual QA

- QA level: QA-0
- Flow or step: Review provider capability declarations for conservative non-Mercado Livre support.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: only Mercado Livre advertises the new operational business capabilities.
- Actual: Mercado Livre declares `listing_read`, `stock_read`, `stock_write`, and `order_read`; Amazon, Leroy Merlin, and MadeiraMadeira declare none; Magalu and Shopee keep only `pricing_fee_sync`.
- Blocking condition: a future provider appears operational for unimplemented capability work

## Evidence

- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/F-02-provider-capability-registration/spec.md`
  Status: Pass
  Evidence type: ran
  Owner: Feature Implementer
  Blocking condition: spec missing or lacks milestone-traced criteria

- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/F-02-provider-capability-registration/plan.md`
  Status: Pass
  Evidence type: ran
  Owner: Feature Implementer
  Blocking condition: plan missing expected changed paths or verification mapping

- Artifact: `contracts/api/marketplace-central.openapi.yaml`
  Status: Pass
  Evidence type: ran
  Owner: Feature Implementer
  Blocking condition: public contract still exposes legacy capability fields

- Artifact: `packages/sdk-runtime/src/index.ts`
  Status: Pass
  Evidence type: ran
  Owner: Feature Implementer
  Blocking condition: SDK type surface still exposes legacy capability fields

## Risks

- The UI currently uses `declared_capabilities` mostly as a preview string, so the semantic value of the new vocabulary becomes fully useful only after F-03 plugs real operational capability behavior behind it.
- Non-Mercado Livre provider capability declarations were intentionally reduced to conservative values; if later runtime support already exists for a provider, that should be raised in a dedicated feature rather than assumed here.

## Milestone Acceptance Review

- Decision: accepted
- Decision owner: Milestone Orchestrator
- Scope check: changed paths stayed inside marketplace/integration contract surfaces, SDK, and feature artifacts.
- Evidence check: all load-bearing criteria rely on `ran` evidence only.
- Notes: this feature changes contract and registry surfaces but does not cross auth/PII/security handling in a way that needs independent QA Validator review before milestone integration.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Feature Implementer
- Next action: Start F-03 using the accepted capability contracts and provider declarations from F-01 and F-02.
- Required files/evidence: feature brief, spec, plan, changed paths, validation evidence
- Blockers or open decisions: None for F-02.
