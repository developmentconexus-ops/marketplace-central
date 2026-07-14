# F-02 catalog-api-envelope

```yaml
id: F-02
type: feature-plan
status: planned
owner: Feature Implementer
parent: M-02
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: feature
split_decision: single
```

## Steps

1. Add transport DTOs, query validation, page-port wiring, typed error
   mapping, freshness-policy context mapping, and interactive route
   registration.
2. Add transport-local fake-port httptest coverage for the cursor chain,
   search, invalid inputs, unavailable source, deadline, and no-cache mapping.
3. Update OpenAPI paths/schemas and the SDK page types and methods in the same
   change.
4. Run focused tests, then the required Go build/test, SDK build, and commit
   stat evidence.

## Files Expected To Change

- `apps/server_core/internal/modules/catalog/transport/http_handler.go`:
  paginated list/search transport and IC-01 mapping.
- `apps/server_core/internal/modules/catalog/transport/http_handler_test.go`:
  fake-port httptest evidence.
- `contracts/api/marketplace-central.openapi.yaml`: IC-01 catalog paths and
  schemas.
- `packages/sdk-runtime/src/index.ts`: typed catalog page models and methods.
- `packages/sdk-runtime/src/index.test.ts`: SDK request/type coverage.
- `.mnfs/.../F-02-catalog-api-envelope/spec.md`: feature specification.
- `.mnfs/.../F-02-catalog-api-envelope/plan.md`: this implementation plan.
- `.mnfs/.../F-02-catalog-api-envelope/validation.md`: command transcripts.

## Verification Commands

- Command: `GOCACHE=.gocache go test ./internal/modules/catalog/transport`
  from `apps/server_core`.
  - Satisfies criterion ID: M-02-C01, M-02-C02
  - Expected result: Pass; all transport transcripts and deadline wiring pass.
- Command: `GOCACHE=.gocache go build ./...` from `apps/server_core`.
  - Satisfies criterion ID: M-02-C01, M-02-C02
  - Expected result: Pass; server compiles against F-01 interface.
- Command: `GOCACHE=.gocache go test ./...` from `apps/server_core`.
  - Satisfies criterion ID: M-02-C01, M-02-C02
  - Expected result: Pass; all server tests pass.
- Command: `npm run build` from `packages/sdk-runtime`.
  - Satisfies criterion ID: M-02-C03
  - Expected result: Pass; typed SDK compiles.
- Command: `npm test -- --run` from `packages/sdk-runtime`.
  - Satisfies criterion ID: M-02-C03
  - Expected result: Pass; SDK request and OpenAPI lockstep tests pass.
- Command: `git show --stat HEAD` after the intentional commit.
  - Satisfies criterion ID: M-02-C03
  - Expected result: One commit lists OpenAPI, SDK, routes, tests, and evidence.

## QA Steps

- Exercise the three-page list cursor chain and assert the final cursor is
  JSON `null`, with one fake-port call per page.
- Exercise garbage cursor, limits `0`, `101`, and search limit `51`; assert
  exact IC-01 error code bodies and zero port calls for validation failures.
- Exercise source unavailable and a deadline-bound route; assert `503` and
  `504` redacted bodies.
- Exercise `Cache-Control: no-cache`; assert the fake observes
  `FreshnessPolicy{MaxAge: 0}` through request context.

## Rollback/Risk Notes

- The old unpaginated `/catalog/products` flow is intentionally replaced by
  the page route; non-list catalog operations remain on the existing handler.
- F-01's page interface has no freshness parameter. The no-cache policy is
  carried in request context as a typed seam only; no port package edit or
  cache behavior is introduced.
- If the supplied F-01 interface changes, stop and return to the spec rather
  than widening scope into `internal_read`.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Implement and record validation evidence.
- Required files/evidence: spec, changed paths, exact command output
- Blockers or open decisions: None.
