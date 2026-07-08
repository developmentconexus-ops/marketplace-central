# Marketplace Integration Foundation Design

Date: 2026-07-08
Status: approved design for planning
Repository: `C:\Users\leandro.theodoro\Documents\marketplace-central`

## Purpose

Design the long-term marketplace integration foundation for Marketplace Central while keeping the first implementation slice safe, real, and verifiable.

The platform must stop treating Mercado Livre as the architectural center. Mercado Livre is the first real adapter used for validation, but the core concepts must belong to Marketplace Central: installations, credentials, connection health, runtime capabilities, operation runs, business policies, and provider adapters.

This design intentionally separates two layers:

- Target model: the complete architecture for reads, writes, messages, shipments, webhooks, and future providers.
- Implementation slice 1: the subset that becomes operational now: auth, connection snapshot, executable read/probe capabilities, operation runs, API/SDK/UI alignment, and live Mercado Livre validation.

## Global Maximum Rationale

This is the global maximum because it preserves the reason each module exists:

- `integrations` owns connection lifecycle, OAuth/auth sessions, credentials, runtime capability state, provider health, and operation audit.
- `connectors` owns external API calls, provider payload mapping, provider error translation, and provider-specific protocol details.
- Business modules own business decisions: pricing, inventory, orders, messages, shipments, safety policies, idempotency, and internal state changes.
- Provider adapters are replaceable implementations of generic platform ports, not places where platform policy leaks.

This avoids the two most dangerous legacy paths:

- A god-module where `integrations` starts deciding stock, fee, order, shipment, and message business rules.
- A provider-first platform where every business module calls Mercado Livre directly and spreads tokens, payloads, and error semantics across the codebase.

The design also prevents "capability theater": runtime APIs and UI must show only capabilities that are implemented, registered, composed, and executable. Future support may be modeled and documented, but not presented as available operation.

## Target Architecture

### Module Roles

```text
integrations
  Installation aggregate
  ConnectionSnapshot projection
  OAuth/auth session lifecycle
  Credential storage and rotation
  Runtime capability registry
  OperationRun audit
  Provider health and reauth state

connectors
  Generic capability ports
  Marketplace adapter implementations
  Provider API clients
  Provider payload mapping
  Error translation
  No tenant business state

business modules
  pricing, inventory, orders, messages, shipments
  Business rules and policies
  Internal persistence and projections
  Idempotency and before/after audit for business effects

frontend/sdk
  Thin operational surfaces
  API contract driven through sdk-runtime
  No direct provider calls
  No credential/token exposure
```

### Core Concepts

`Installation` is the user-facing and API-facing unit of integration. It means: a tenant connected one external account for one provider.

`ConnectionSnapshot` is the safe operational projection stored on the installation. It contains:

- provider code
- installation id
- external account id
- external account display name
- connection state
- health status
- auth strategy
- last verified timestamp
- credential expiry timestamp when known
- required next action
- reauth reason when applicable

`Credential` and raw tokens are internal to `integrations`. Controllers, React, business modules, and generic connector consumers should not receive raw provider tokens.

`RuntimeCapability` is an operational capability that is actually available. A capability is runtime-available only when:

- the provider adapter implements the port
- composition registers the adapter
- the installation is connected and healthy enough for that capability
- the application exposes a safe execution path
- local contract validation exists

For capabilities that claim real provider behavior, live validation evidence is also required. A local-only capability may exist during development, but it must be labeled as not live-validated and must not be used as proof of provider integration.

`ProviderModel` describes what the platform can conceptually support over time. It can include stock writes, messages, shipments, webhooks, listing changes, and future provider features. It must not be confused with `RuntimeCapability`.

`OperationRun` records technical execution against an external provider. It captures the installation, capability, start/end time, status, provider evidence, translated error, duration, and correlation metadata. It is not the business audit record for stock, price, order, or message decisions.

## Capability Model

Capabilities are platform concepts, not Mercado Livre feature names.

Target capability ports:

- `AccountProbe`: read external account identity and health.
- `ListingReader`: read listings/items from provider.
- `OrderReader`: read orders from provider.
- `StockReader`: read provider-visible stock.
- `StockWriter`: write stock safely after internal policy approval.
- `FeeQuoteReader`: read fee quote for product/category/listing context.
- `FeeSyncReader`: synchronize fee tables or fee rules when provider supports it.
- `MessageReader`: read buyer questions/messages.
- `MessageReplier`: send replies after policy and audit checks.
- `ShipmentReader`: read shipment/tracking data.
- `WebhookReceiver`: receive, authenticate, normalize, and dispatch provider events.
- `ListingWriter`: future capability for listing changes, not in slice 1.

Slice 1 should expose only safe read/probe capabilities that have real paths:

- `account_probe`
- `listing_read`
- `order_read`
- `fee_quote_read` or a narrowly scoped `pricing_fee_sync` only if it calls the live provider API rather than static seed data
- `stock_read` only if the adapter path is composed and validated against the live provider or explicitly marked local-only until live evidence exists

Writes must not appear as runtime capabilities in slice 1.

## Auth And Credential Flow

The OAuth callback must end in a transactional application operation named conceptually as `ApplyAuthResult`.

`ApplyAuthResult` responsibilities:

- validate callback/session state
- exchange provider code for credential payload
- rotate or create active credential
- persist or update auth session
- project external account id/name onto installation
- project `ConnectionSnapshot` onto installation
- set installation status and health
- update last verified timestamp
- activate only runtime capabilities that are truly executable
- record enough evidence for troubleshooting without exposing secrets

This fixes the current failure mode where Mercado Livre is connected for real, but the installation still displays empty credential/account fields.

Credential resolution should use an internal `CredentialResolver` port owned by `integrations`. Callers pass execution context such as tenant id, installation id, and provider code. Raw tokens should not be passed through controllers, React, or business modules.

## Operation Execution Flow

Recommended runtime flow:

```text
business module
  -> integrations validates installation and runtime capability
  -> integrations opens OperationRun
  -> connector executes provider capability using CredentialResolver
  -> connector returns normalized result or typed provider error
  -> integrations closes OperationRun with evidence
  -> business module applies business policy, if any
```

This keeps provider communication, connection governance, and business decisions separate.

For read/probe flows, the business module may simply consume normalized data.

For write flows, the business module must own the policy gate before `connectors` executes the write. Example for future stock write:

- internal product link is resolved and unambiguous
- source stock and timestamp are visible
- safety policy is explicit
- write is idempotent or duplicate-safe
- before/after values are audited
- provider response is stored

`OperationRun` records the external operation. Business modules record business effect audit.

## API, SDK, And UI Shape

The OpenAPI contract should expose installation state as a safe operational model:

```text
IntegrationInstallation
  id
  tenant_id
  provider_code
  display_name
  status
  connection
  runtime_capabilities
  last_operation_summary
  created_at
  updated_at

ConnectionSnapshot
  state
  health
  external_account_id
  external_account_name
  auth_strategy
  last_verified_at
  expires_at
  next_action
  reauth_reason
```

The frontend must render the installation, not raw credential internals.

Recommended UI zones:

- Connection: provider, external account, connection state, health, last verified, expires/reauth, next action.
- Capabilities: only runtime capabilities with states like `available`, `unavailable`, `needs_auth`, `degraded`, `not_configured`.
- Operations: recent probes/syncs with status, translated errors, provider evidence, duration, and next action.

The UI should not show "Credential Not connected" when the installation is actually connected. Credential id is an internal audit detail, not the operator's primary state.

## Error And Quality States

Provider and connection failures should become typed platform states:

- `needs_auth`
- `unauthorized`
- `forbidden`
- `rate_limited`
- `provider_unavailable`
- `provider_contract_changed`
- `unsupported_capability`
- `not_configured`
- `ambiguous_account`
- `missing_external_account`
- `degraded`

Absence must not be converted into a default value. Unknown fees, stock, freight, tax, product linkage, or account identity are explicit quality states.

For writes, unsupported or unsafe states must block execution rather than degrade into best-effort behavior.

## Implementation Slice 1

Slice 1 implements the foundation without pretending future writes are ready.

### Backend Domain And Application

Work in `apps/server_core/internal/modules/integrations`:

- Introduce or refine domain types for `Installation`, `ConnectionSnapshot`, `RuntimeCapability`, `OperationRun`, `AuthSession`, and `CredentialRef`.
- Add typed enums for provider code, connection state, health, capability state, operation status, and next action.
- Add `ApplyAuthResult` transaction semantics.
- Add capability validation before operation execution.
- Add operation run begin/finish application service.
- Ensure auth status and installation list read from the same connection snapshot truth.

### Connector Ports And Mercado Livre Adapter

Work in `apps/server_core/internal/modules/connectors`:

- Define generic capability ports for account probe, listings, orders, fees, and optional stock read.
- Keep provider payload structs inside provider adapter packages.
- Keep SQL, provider SDK/client logic, HTTP details, and token usage outside business modules.
- Translate provider errors into typed connector errors.
- Wire Mercado Livre adapter through composition rather than leaving capability code orphaned.

For Mercado Livre specifically:

- Use real account probe equivalent to `/users/me`.
- Use live reads for listings/orders/fees where safe.
- Replace or demote static `pricing_fee_sync` seed behavior so runtime metadata does not claim live API sync unless it actually performs one.

### Contracts, SDK, And Frontend

- Update `contracts/api/marketplace-central.openapi.yaml` for installation connection snapshot, runtime capabilities, and operation summaries.
- Regenerate or update `packages/sdk-runtime` consistently with OpenAPI.
- Update integrations UI to consume SDK types only.
- Remove credential-centric display from operator-facing state.
- Show clear loading, error, empty, connected, degraded, and needs-auth states.

### Validation

Local validation:

- Unit tests for `ApplyAuthResult` projection.
- Unit tests for capability declaration gating.
- Unit tests for credential resolver usage without token leakage.
- Unit tests for operation run status/error transitions.
- UI tests or component tests for connection/capability states if the repo already has a pattern.

Live validation:

- Mercado Livre `/users/me` account probe using active credential.
- Listings read using active credential.
- Orders read using active credential.
- Fee quote/read using active credential.
- Document exact command, endpoint, timestamp, installation id, provider account id, and response summary.

Validation boundaries:

- Fake/mock/seam tests prove only local contract and deterministic behavior.
- Live Mercado Livre evidence is required before declaring real provider integration validated.
- No write capability is validated or exposed until a later slice with explicit safety gates.

## Deferred But Modeled Slices

### Slice 2: Stock Safety And Writes

- Product link resolution.
- Stock read reconciliation.
- Stock write policy.
- Idempotency key strategy.
- Before/after audit.
- Provider response persistence.
- Live validation only against intentionally selected safe product/link.

### Slice 3: Messages

- Message/question read model.
- Reply policy and templates.
- Human approval or automation mode.
- Message reply idempotency.
- Customer communication audit.

### Slice 4: Shipments

- Shipment reader.
- Tracking normalization.
- Order/shipment association.
- Provider-specific shipment states mapped to platform states.

### Slice 5: Webhooks

- Provider webhook registration model.
- Signature/auth verification.
- Event normalization.
- Idempotent event ingestion.
- Dispatch into business modules through events.
- Replay and dead-letter strategy.

### Slice 6: Listing Writes

- Listing linkage.
- Field-level write policies.
- Draft vs publish mode.
- Before/after diff.
- Safe rollback or compensating action where provider supports it.

## Execution Plan For Clean Session

The implementation should be run as MNFS-style execution with feature boundaries. Use subagents where tasks are independent, but keep architectural decisions and final acceptance in the main session.

### Phase 0: Rehydrate Context

Read:

- `C:\Users\leandro.theodoro\Documents\marketplace-central\ARCHITECTURE.md`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\wiki\README.md`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.brain\system-pulse.md`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.brain\roadmap.json`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\AGENTS.md`
- this spec

Then inspect current integration code, connector code, OpenAPI, SDK, and frontend integrations page.

### Phase 1: Audit And Reconcile Truth

Use subagents for parallel read-only audits:

- Subagent A: `integrations` domain/application/transport/persistence/auth flow.
- Subagent B: `connectors` ports/adapters/Mercado Livre capability code and composition.
- Subagent C: OpenAPI, SDK runtime, frontend integrations UI.
- Subagent D: tests, validation artifacts, live Mercado Livre evidence paths.

Main session reconciles findings into a scoped implementation plan. Stop if architecture, contract, runtime, or ownership contradictions appear.

### Phase 2: Domain And Application Foundation

Implement `integrations` domain/application changes first:

- `ConnectionSnapshot`
- runtime capability model
- operation run model
- `ApplyAuthResult`
- capability validation
- credential resolver boundary

Validation:

- targeted Go unit tests for domain/application behavior
- no live claims yet

### Phase 3: Connector Wiring And Live Read Capabilities

Implement/wire generic connector ports and Mercado Livre adapter execution:

- account probe
- listing read
- order read
- fee read/quote or honest demotion of current static seed sync
- typed provider error translation

Validation:

- local adapter contract tests with test doubles
- live read probes using active Mercado Livre credential
- no writes

### Phase 4: Contract, SDK, UI

Update the external platform surface:

- OpenAPI installation model
- SDK runtime types/client
- integrations page rendering connection/capabilities/operations

Validation:

- frontend tests/build as available
- browser validation through ngrok/local app
- verify UI no longer shows credential-disconnected state for connected installation

### Phase 5: Evidence And Closeout

Update evidence artifacts:

- MNFS feature validations
- milestone validation result
- any wiki/architecture notes if implementation reveals new truth

Final claim must distinguish:

- local contract validation
- live Mercado Livre read validation
- deferred write validation

## Handoff Prompt For Clean Session

Use this prompt to start a clean implementation session:

```text
Context:
Repo: C:\Users\leandro.theodoro\Documents\marketplace-central
We approved the marketplace integration foundation design in:
C:\Users\leandro.theodoro\Documents\marketplace-central\docs\superpowers\specs\2026-07-08-marketplace-integration-foundation-design.md

Objective:
Create the implementation plan and execute slice 1 of the marketplace integration foundation.

Global maximum:
Model the complete marketplace platform, but operationalize only safe, real, verifiable read/probe capabilities now.

Approved architecture:
- integrations owns Installation, ConnectionSnapshot, auth sessions, credentials, runtime capabilities, operation runs, health, and credential resolution.
- connectors owns provider API calls, payload mapping, provider errors, and marketplace adapters.
- business modules own pricing/inventory/orders/messages/shipments policies and internal side effects.
- Mercado Livre is first real adapter, not the architecture center.
- Runtime capabilities appear only when implemented, wired, composed, and executable.
- Future writes/messages/shipments/webhooks are modeled but not exposed as operational capabilities in slice 1.

Must read first:
- AGENTS.md
- ARCHITECTURE.md
- wiki/README.md
- .brain/system-pulse.md
- .brain/roadmap.json
- docs/superpowers/specs/2026-07-08-marketplace-integration-foundation-design.md

Required workflow:
- Use writing-plans before coding.
- Use subagents for read-only audits of integrations, connectors, API/SDK/UI, and validation.
- Do not write code until the plan is explicit.
- Update OpenAPI and sdk-runtime together for API changes.
- Never claim live integration validation from mocks/fakes/seams.
- No provider writes in this slice.
- Do not expose raw tokens outside integrations.
- Do not show future capabilities as runtime available.

Known audit facts:
- Mercado Livre auth is real and live account reads previously succeeded.
- UI currently can show misleading credential/account state because installation projection is incomplete.
- Some Mercado Livre connector capability code exists but is not fully composed/exposed.
- Current pricing_fee_sync behavior may be static seed; do not claim live fee sync unless converted to real provider API behavior or demoted honestly.

Target deliverable:
Implemented slice 1 with tests/builds/live read evidence, updated validation artifacts, and a clear closeout separating local, integration, and live validation.
```

## Acceptance Criteria

- The design remains provider-neutral while Mercado Livre works as first real adapter.
- Installation is the primary operator-facing unit.
- Connection state is projected through `ConnectionSnapshot`.
- Runtime capabilities are honest and executable.
- Credential handling stays inside `integrations`.
- Provider payloads stay inside provider adapters.
- Business modules own business rules and side effects.
- Operation runs provide technical external-operation audit.
- OpenAPI, SDK, backend, and UI agree on the same model.
- Validation artifacts explicitly separate fake/local evidence from live provider evidence.
- Writes remain modeled but unavailable until a later safety milestone.

## Open Non-Blocking Decisions

- Exact naming of `FeeQuoteReader` versus `FeeSyncReader` can be finalized after auditing current pricing needs.
- Exact persistence shape for `ConnectionSnapshot` can be column-based or JSON projection, but API shape must be typed and stable.
- Whether operation execution is centralized in an `integrations` application service or exposed as small application ports can be chosen during implementation planning, as long as ownership boundaries remain intact.
