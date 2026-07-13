# F-13 Assisted Linkage API and Runtime Wiring Validation

```yaml
id: F-13
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: F-13
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-13-assisted-linkage-api-runtime

## Summary

Quick validation passed. The assisted Sankhya workflow is always registered at
three exact-scope orders routes, but only composes when all eight explicit
runtime settings and the shared Oracle reader are available. Confirmation uses
only operator selection/intent facts from the caller; event identity, default
tenant, external key, fixed actor type, configuration revision, and opaque
evidence reference remain server-owned. OpenAPI and SDK changed together.

## Quick Validation Result

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-13
- Final feature state for handoff: `quick_validation_passed`

## Quick Validation State

- fixup_attempts: 1
- max_fixup_attempts: 1
- last_feature_validation_result: Pass after full rerun

The initial focused Go run reproduced one compile-only defect:
`sankhya_linkage_handler.go:260:2: undefined: SankhyaLinkageResponse`. The
confirmation DTO embedded the wrong type name. The single fixup changed it to
the intended private `sankhyaLinkageResponse`; the complete focused Go command
then passed, followed by repository-wide Go tests.

## Context Evidence

- The first planning compile failed `CTX_FEATURE_INVALID` at
  `machine-work-contract-count` because the newly written plan did not yet
  contain the required machine contract; the next failed `mnfs-headings`
  because the supplied brief used `Outcome` rather than compiler headings.
- The plan now preserves every dispatch path, side effect, proof command,
  criterion, stop condition, and route. `feature.md` carries equivalent `Brief`
  and `Expected Output` headings.
- The canonical full L2 context compiled successfully with explicit overflow
  reasons for every source; the known L1 route-budget defect did not recur, so
  no reduced pack was used or claimed.
- Final `context.json` was recompiled after API/SDK/governance changes and
  passed `context-validate -RequireCurrentBase` at accepted base
  `90e771ae589d057928fe3f810d02f43e45967256`.
- Evidence type: ran
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-13-assisted-linkage-api-runtime/context.json`

## Spec Adherence

- Spec satisfied: Yes
- Deviations: None.
- Reason: Runtime, service, transport, composition, OpenAPI, SDK, governance,
  and fake/unit proof implement F13-AC01 through F13-AC05 without live access.

## Changes Made

- `apps/server_core/internal/modules/internal_read/adapters/oracle/sankhya_linkage_config.go`
  and test: explicit enabled/schema/header-field/revision/attestation/three-limit
  loading, bounds, fixed TOP 313/306, and opaque deterministic evidence digest.
- Orders reader port/internal-read bridge/application service and tests: one
  evidence reference alongside the revision, exact `GetCurrent`, and audited
  confirmation propagation.
- `apps/server_core/internal/modules/orders/transport/sankhya_linkage_handler.go`
  and test: strict current/candidate/confirm routes, default tenant, exact query
  and path scope, bounded strict JSON, narrow nullable DTOs, and safe errors.
- `apps/server_core/internal/composition/root.go`: always register the workflow;
  wire the Oracle reader, bridge, service, and F-09 repository only when the
  explicit runtime and shared Oracle source are available.
- OpenAPI, SDK runtime, SDK tests, and governance registry: atomic public
  contract plus every environment key.
- Feature artifacts: brief/compiler headings, spec, plan, final context, and
  validation evidence.

## Commands Run

### context-compile / context-validate

- Command: `pwsh -NoProfile -File scripts/harness.ps1 -Command context-compile -FeaturePath .mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-13-assisted-linkage-api-runtime -BaseSha 90e771ae589d057928fe3f810d02f43e45967256 -AllowedPath <dispatch allowed_paths>`; then `pwsh -NoProfile -File scripts/harness.ps1 -Command context-validate -ContextPath .mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-13-assisted-linkage-api-runtime/context.json -RequireCurrentBase`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: canonical bounded L2 context compiles and is current at accepted SHA.
- Actual: final compile and current-base validation exited 0.
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-13-assisted-linkage-api-runtime/context.json`
- Blocking condition: None.

### go-assisted-linkage-http

- Command: from `apps/server_core`, set `GOCACHE=.gocache`; `go test ./internal/modules/internal_read/adapters/oracle ./internal/modules/orders/application ./internal/modules/orders/adapters/internalread ./internal/modules/orders/transport ./internal/composition`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: config, current/evidence service, bridge, strict transport, safe
  errors, fail-closed route registration, and composition build pass without a
  database or network.
- Actual: after the one recorded compile-name fixup, all four test packages
  passed and composition built (`[no test files]`).
- Artifact: `apps/server_core/internal/modules/orders/transport/sankhya_linkage_handler_test.go`
- Blocking condition: None.

### repository-wide Go build/test

- Command: from `apps/server_core`, set `GOCACHE=.gocache`; `go test ./...`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: the server commands and all Go packages compile/test with no live
  side effect.
- Actual: exited 0; command packages and composition built, unit packages passed.
- Artifact: `apps/server_core`
- Blocking condition: None.

### sdk-assisted-linkage

- Command: `npm test --workspace @marketplace-central/sdk-runtime`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: SDK route, encoding, response, nullable lineage, and confirm-body
  tests pass.
- Actual: Vitest exited 0.
- Artifact: `packages/sdk-runtime/src/index.test.ts`
- Blocking condition: None.

- Command: `npm exec -- tsc --noEmit --target ES2022 --module NodeNext --moduleResolution NodeNext --lib ES2022,DOM packages/sdk-runtime/src/index.ts`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: standalone SDK public types and methods compile.
- Actual: exited 0 with no diagnostics.
- Artifact: `packages/sdk-runtime/src/index.ts`
- Blocking condition: None.

- Command: parse OpenAPI with the available Python YAML runtime; assert all
  three paths, `additionalProperties: false`, exact confirm required/property
  parity, and absence of tenant/event/config/evidence/actor-type/external-key;
  inspect the matching SDK interface for the same forbidden fields.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: public contract parses and caller cannot submit server-owned facts.
- Actual: `PASS assisted linkage OpenAPI contract`; `PASS assisted linkage SDK input parity`.
- Artifact: `contracts/api/marketplace-central.openapi.yaml`
- Blocking condition: None.

### governance-contracts

- Command: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/governance-contracts.tests.ps1`; then `pwsh -NoProfile -File scripts/harness.ps1 -Command governance -BaseSha 90e771ae589d057928fe3f810d02f43e45967256`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: runtime registry schema/policy and API-SDK atomic seam pass.
- Actual: governance contract tests printed `PASS`; governance command exited 0.
- Artifact: `contracts/governance/runtime-config.json`
- Blocking condition: None.

### git-diff-check

- Command: `git diff --check`; inspect tracked/untracked paths against
  `dispatch.json.allowed_paths`, excluding pre-existing user-owned
  `docs/research/**` and `output/**` state from this Feature.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: no whitespace errors and only dispatched owned paths enter the
  Feature commit.
- Actual: no whitespace errors; all Feature-owned changes are within dispatch.
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-13-assisted-linkage-api-runtime/dispatch.json`
- Blocking condition: None. Git emitted only existing LF/CRLF conversion warnings.

## Manual QA

- QA level: QA-0
- Flow or step: Compare OpenAPI confirm request and SDK confirm input.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: only selected document/lines, actor ID, reason, source time, and
  idempotency key are writable.
- Actual: automated parity assertions and SDK body test passed; all six
  forbidden server controls are absent.
- Blocking condition: None.

- QA level: QA-0
- Flow or step: Inspect runtime/composition unavailable behavior.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: absent/disabled/malformed/out-of-range values never default or stop
  unrelated startup; routes remain registered and return stable 503.
- Actual: table-driven loader tests and nil-runtime transport tests passed;
  root composition retains unrelated Oracle behavior.
- Blocking condition: None.

- QA level: QA-0
- Flow or step: Inspect response mapping for nullable facts, provenance, and
  lineage states.
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Expected: no Oracle field/SQL/provider/PII exposure and no unknown-to-zero.
- Actual: transport and SDK tests preserve null facts, generic identifiers,
  server audit fields, exact descendants, and explicit states.
- Blocking condition: None.

## Evidence

- Artifact: `apps/server_core/internal/modules/internal_read/adapters/oracle/sankhya_linkage_config_test.go`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Blocking condition: None.
- Artifact: `apps/server_core/internal/modules/orders/application/assisted_sankhya_linkage_service_test.go`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Blocking condition: None.
- Artifact: `apps/server_core/internal/modules/orders/transport/sankhya_linkage_handler_test.go`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Blocking condition: None.
- Artifact: `packages/sdk-runtime/src/index.test.ts`
- Status: Pass
- Evidence type: ran
- Owner: Feature Implementer
- Blocking condition: None.

## Risks

- Production authentication/manual-adjustment authorization remains explicitly
  deferred; actor provenance is `operator_supplied_unverified` by contract.
- No live Oracle metadata/uniqueness/lineage validation was run. A syntactically
  valid deployment still fails closed through F-11 typed runtime validation if
  metadata, uniqueness, or source availability is not proved.
- This is Feature quick validation only; Milestone review and proportional QA
  still own acceptance/pass.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: Review the F-13 commit, spec/plan/context, changed paths, and
  validation evidence; integrate before fixed-SHA review and QA.
- Required files/evidence: `feature.md`, `spec.md`, `plan.md`, `context.json`,
  `validation.md`, committed diff
- Blockers or open decisions: None.
