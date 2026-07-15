# F-01-listings-module-ingestion

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-01
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contract: IC-02 `../../research/listings-read-interface-contract.md` (entity, enums, DB shape, seed data). ADR-12 (read-only module, ingestion via connectors only), ADR-17 (nullable unknowns).

## Milestone

M-01 listings-read-spine.

## Brief

Create the `listings` module (`apps/server_core/internal/modules/listings/` with domain/application/ports/adapters/transport layout per governance `modules.json` conventions), the `listings` table migration exactly per IC-02 Database Shape, the ML→canonical mapping in an adapter that consumes the existing connectors `ListListings` capability, and `POST /listings/refresh` which runs a full-pull ingestion recorded as an integration operation run.

EARS:
- While an installation is connected, when `POST /listings/refresh` is accepted, the listings module shall pull all listing pages via the connectors capability, upsert canonical rows keyed `(tenant_id, installation_id, provider_listing_id, variation_id)`, and mark rows absent from the completed pull as `status=closed`.
- While a refresh for the same installation is already running, when a second refresh is requested, the module shall return 409 `refresh_in_progress` with the active `operation_run_id`.
- While a provider value cannot be mapped to a canonical enum, when ingesting, the adapter shall store `unknown` (status) or null (facts) and never a guessed value.
- While the provider is unreachable, when a refresh runs, the module shall fail the operation run honestly and leave existing rows untouched.

## Inputs

- IC-02 (all shapes/enums/keys), R-03 gap matrix, connectors capability signatures (`apps/server_core/internal/modules/connectors/domain/capability.go`, `adapters/mercado_livre/capability_adapter.go`), `integration_operation_runs` pattern in integrations module, migration numbering (next after 0033), composition root wiring pattern.

## Expected Output

- Migration `00xx_listings.sql` per IC-02 DB shape (all fact columns nullable; enum checks).
- Module code compiling, wired in composition root.
- `POST /listings/refresh` per IC-02 (202 envelope; 400/404/409 per Error Matrix).
- Unit tests: mapping (incl. unmappable → unknown/null), upsert keying, closed-marking.
- Integration test (ephemeral-postgres): refresh against a stubbed capability seeds the IC-02 seed set (MLBTEST0001..0006 shapes).
- One intentional commit; OpenAPI + sdk-runtime updated in the SAME commit for the refresh endpoint.

## Constraints

- Read-only module: no provider write imports, no mutation of provider state.
- Provider payloads (raw ML JSON) never leave the adapter.
- Link fields NOT stored (joined from product_links at read — F-02 concern).
- Do not touch product_links snapshots or its module.
- `GOCACHE` absolute for Go tests; governance lanes must stay green.

## Negative Scenarios

- Unknown installation → 404 `installation_not_found`.
- Missing installation_id → 400 `installation_required`.
- Concurrent refresh → 409 `refresh_in_progress` (proof: two requests, second returns 409 + active run id).
- Capability error mid-pull → operation run `failed`; row count unchanged from pre-refresh (proof: count assert).

## Validation Expectations

- `go test` output showing the mapping + upsert + closed-marking tests green.
- Integration transcript: `POST /listings/refresh` → 202 `{operation_run_id}`; second concurrent → 409 body with `refresh_in_progress`.
- SQL assert: seeded pull produces exactly the seed rows with nulls where IC-02 says nullable (e.g. `sales_30d IS NULL`).

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (creates `spec.md`).
- Next action: compile feature context pack; read only IC-02 + named code paths.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
