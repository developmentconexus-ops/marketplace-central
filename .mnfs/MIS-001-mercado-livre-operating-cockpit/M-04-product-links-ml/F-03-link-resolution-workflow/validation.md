# Feature Validation

```yaml
id: F-03
type: feature-validation
status: passed
owner: Feature Implementer
parent: F-03
created: 2026-07-08
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-03-link-resolution-workflow

## Summary

Backend workflow, contract parity, dedicated UI route, and real operator-resolution persistence are validated for the operator workflow that turns generated product link candidates into audited resolved or rejected link truth.

## Current Validation State

- Result: Local and live validation passed
- Result owner: Feature Implementer
- Decision date: 2026-07-08
- Final feature state for handoff: ready_for_milestone_review

## Planned Evidence

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/product_links/... ./apps/server_core/internal/composition -count=1`
  - Result: Pass on 2026-07-08
  - Evidence: `application` and `transport` packages passed with workflow persistence, audit, and transport coverage green together
- Command: `npm run test --workspace @marketplace-central/sdk-runtime -- src/index.test.ts`
  - Result: Pass on 2026-07-08
  - Evidence: 26 tests passed including workflow list and approve-candidate request coverage
- Command: `npm run test --workspace @marketplace-central/feature-product-links`
  - Result: Pass on 2026-07-08
  - Evidence: 6 tests passed covering loading, error, empty, conflict, unresolved, resolved, rejected, approve, reject, and manual resolve flows
- Command: `npm run test --workspace @marketplace-central/web -- src/app/AppRouter.test.tsx`
  - Result: Pass on 2026-07-08
  - Evidence: app router resolves both `/integrations` and `/product-links`
- Command: `npm run build --workspace @marketplace-central/web`
  - Result: Pass on 2026-07-08
  - Evidence: production build completed with the new feature package linked into the app
- Command: `go run ./cmd/migrate` against `postgres://marketplace:marketplace@127.0.0.1:5435/marketplace_central?sslmode=disable` with `MC_MIGRATIONS_DIR=migrations`
  - Result: Pass on 2026-07-08
  - Evidence: applied 1 migration (`0025_product_link_workflows.sql`) to the live Docker Postgres used by the running stack
- Command: `POST http://localhost:8090/product-links/link-resolutions/approve-candidate`
  - Result: Pass on 2026-07-08
  - Evidence: real installation `inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98` resolved `MLB4807275656` to internal product `20307` with audit action `approve_candidate`
- Command: `POST http://localhost:8090/product-links/link-resolutions/reject-listing`
  - Result: Pass on 2026-07-08
  - Evidence: real installation rejected `MLB4834408384` with audit action `reject_listing`
- Command: `POST http://localhost:8090/product-links/link-resolutions/manual-resolve`
  - Result: Pass on 2026-07-08
  - Evidence: real installation manually resolved `MLB4834419602` to internal product `20312` with audit action `manual_resolve`
- Command: `GET http://localhost:8090/product-links/link-workflows?installation_id=inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98&limit=20`
  - Result: Pass on 2026-07-08
  - Evidence: workflow response shows `current_link` plus `audit` entries for the three live actions while unresolved candidates remain visible for untouched items
- Command: direct Postgres proof against `product_links` and `product_link_audit_entries`
  - Result: Pass on 2026-07-08
  - Evidence:
    - `MLB4807275656` persisted as `resolved` with `internal_product_id=20307` and source candidate preserved
    - `MLB4834408384` persisted as `rejected`
    - `MLB4834419602` persisted as `resolved` with `internal_product_id=20312`
    - audit rows store action, previous/next state, actor name, and reason for all three actions
- Command: `docker compose up --build -d backend frontend`
  - Result: Pass on 2026-07-08
  - Evidence: official `marketplace-central-backend-1` returned to `healthy`; official `http://localhost:8080/product-links/link-workflows?...` now serves the new workflow endpoint from the running stack

## Scope Declaration

- contract_validated: Yes
- local_business_validation: Yes
- live_validation: Yes
- blocked_for_real_validation: No

## Notes

- OpenAPI now exposes workflow list plus approve/reject/manual resolve commands.
- `packages/sdk-runtime` now mirrors the workflow contract and typed resolution actions.
- `packages/feature-product-links` owns the dedicated operator UI surface instead of coupling the workflow into the integrations hub.
- Real validation used the Docker Postgres and Mercado Livre installation already populated from live listing reads.
- The first live backend attempt failed because the Docker database had not yet received `0025_product_link_workflows.sql`; after applying the migration, the workflow endpoints and persistence behaved correctly.

## Handoff

- Current status: `live_validated`
- Next owner: Milestone Orchestrator
- Next action: milestone review can proceed; if desired, add cleanup or seed-reset guidance for the three live validation rows before broader operator use
- Required files/evidence: kept in this validation artifact plus the persisted rows already present in Docker Postgres
- Blockers or open decisions: none for feature acceptance
