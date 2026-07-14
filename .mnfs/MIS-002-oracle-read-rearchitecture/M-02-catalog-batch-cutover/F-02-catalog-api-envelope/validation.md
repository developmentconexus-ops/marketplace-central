# F-02 catalog-api-envelope validation

```yaml
id: F-02
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-02
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: feature
```

## Summary

- Result: Pass
- Result owner: Feature Implementer
- Decision date: 2026-07-14
- Final feature state for handoff: `quick_validation_passed`

The IC-01 list/search routes, error matrix, route deadline registration,
no-cache freshness mapping, OpenAPI, and typed SDK are implemented. The
production route uses `RouteClassMux`; the explicitly deprecated legacy
registration remains only for direct plain-mux harnesses without the F-01 page
reader, preserving existing compatibility tests while the orchestrator wires
the page reader at composition acceptance.

## Evidence Honesty

Every result below is `ran` with output captured in this artifact unless noted
otherwise. No assumed evidence is used.

## Quick Validation Result

- Result: Pass
- Feature state: `quick_validation_passed`
- Fixup attempts: 1 implementation fixup after full-suite compatibility tests
  exposed the legacy plain-mux harness; focused and full suites were rerun.

## Spec Adherence

- Spec satisfied: Yes
- Deviations: None from IC-01. The supplied F-01 page-port signature has no
  freshness parameter, so `Cache-Control: no-cache` is carried as typed
  `internal_read/domain.FreshnessPolicy{MaxAge: 0}` in request context only.
- Reason: No port or adapter edits are allowed in F-02 and no cache layer is
  added.

## Changed Paths

- `apps/server_core/internal/modules/catalog/transport/http_handler.go`
- `apps/server_core/internal/modules/catalog/transport/http_handler_test.go`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`
- `packages/sdk-runtime/src/index.test.ts`
- `packages/sdk-runtime/package.json`
- `packages/sdk-runtime/tsconfig.json`
- `packages/sdk-runtime/vitest.config.ts`
- `.mnfs/.../F-02-catalog-api-envelope/spec.md`
- `.mnfs/.../F-02-catalog-api-envelope/plan.md`
- `.mnfs/.../F-02-catalog-api-envelope/validation.md`

## Commands Run

### Focused httptest transcript

- Command: `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test -v ./internal/modules/catalog/transport -run 'TestCatalogPageRoutesFollowThreePageCursorChain|TestCatalogPageRoutesValidateBeforePortCall|TestCatalogSearchPageEnvelopeAndNoCachePolicy|TestCatalogPageRoutesMapSourceAndDeadlineErrors'`
- Status: Pass
- Evidence type: `ran`
- Expected: three-page envelope chain, validation before port, search envelope,
  no-cache mapping, source redaction, and deadline code.
- Actual transcript:

```text
GET /catalog/products?limit=1 -> 200 {"as_of":"2026-07-14T12:00:00Z","items":[{"active":true,"cost":{"amount":null,"currency":"BRL","quality":["missing_cost"]},"current_price":{"amount":"12.90","currency":"BRL","quality":[]},"description":"PARAFUSO SEXTAVADO M8","ean":null,"internal_product_id":1,"reference":"ABC-123","sellable_stock":{"quality":[],"quantity":41}}],"next_cursor":"MQ==","page_size":1}
GET /catalog/products?limit=1&cursor=MQ== -> 200 {"as_of":"2026-07-14T12:01:00Z","items":[{"active":true,"cost":{"amount":null,"currency":"BRL","quality":["missing_cost"]},"current_price":{"amount":"12.90","currency":"BRL","quality":[]},"description":"PARAFUSO SEXTAVADO M8","ean":null,"internal_product_id":2,"reference":"ABC-123","sellable_stock":{"quality":[],"quantity":41}}],"next_cursor":"Mg==","page_size":1}
GET /catalog/products?limit=1&cursor=Mg== -> 200 {"as_of":"2026-07-14T12:02:00Z","items":[{"active":true,"cost":{"amount":null,"currency":"BRL","quality":["missing_cost"]},"current_price":{"amount":"12.90","currency":"BRL","quality":[]},"description":"PARAFUSO SEXTAVADO M8","ean":null,"internal_product_id":3,"reference":"ABC-123","sellable_stock":{"quality":[],"quantity":41}}],"next_cursor":null,"page_size":1}
GET /catalog/products?cursor=%25%25%25garbage -> 400 {"error":"invalid_cursor"}
GET /catalog/products?limit=0 -> 400 {"allowed_range":"1..100","error":"invalid_limit"}
GET /catalog/products?limit=101 -> 400 {"allowed_range":"1..100","error":"invalid_limit"}
GET /catalog/products/search?q=PARAFUSO&limit=51 -> 400 {"allowed_range":"1..50","error":"invalid_limit"}
GET /catalog/products/search?q=PARAFUSO&limit=50 -> 200 {"as_of":"2026-07-14T13:00:00Z","items":[{"active":true,"cost":{"amount":null,"currency":"BRL","quality":["missing_cost"]},"current_price":{"amount":"12.90","currency":"BRL","quality":[]},"description":"PARAFUSO SEXTAVADO M8","ean":null,"internal_product_id":10,"reference":"ABC-123","sellable_stock":{"quality":[],"quantity":41}}],"next_cursor":null,"page_size":1}
GET /catalog/products -> 503 {"error":"source_unavailable"}
GET /catalog/products -> 504 {"error":"deadline_exceeded"}
PASS
ok marketplace-central/apps/server_core/internal/modules/catalog/transport
```

### Required build/test commands

- Command: `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go build ./...`
  from `apps/server_core`.
  - Status: Pass
  - Evidence type: `ran`
  - Actual: exit code 0; server build completed with no stdout.
- Command: `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./...`
  from `apps/server_core`.
  - Status: Pass
  - Evidence type: `ran`
  - Actual: exit code 0; all packages passed, including
    `internal/modules/catalog/transport`, `internal/platform/httpx`, and
    `tests/unit`.
- Command: `npm run build` from `packages/sdk-runtime`.
  - Status: Pass
  - Evidence type: `ran`
  - Actual transcript:

```text
> @marketplace-central/sdk-runtime@0.1.0 build
> tsc --noEmit

Process exited with code 0
```

- Command: `npm test -- --run` from `packages/sdk-runtime`.
  - Status: Pass
  - Evidence type: `ran`
  - Actual: `src/index.test.ts` — 41 tests passed; 1 test file passed.
- Command: `python -c "import yaml; doc=yaml.safe_load(open('contracts/api/marketplace-central.openapi.yaml', encoding='utf-8')); print('openapi', doc['openapi'], 'catalog-paths', all(p in doc['paths'] for p in ['/catalog/products','/catalog/products/search']))"`
  - Status: Pass
  - Evidence type: `ran`
  - Actual: `openapi 3.1.0 catalog-paths True`.
- Command: `git diff --check`
  - Status: Pass
  - Evidence type: `ran`
  - Actual: no whitespace errors.

## Fixup History

- Reproduction: the first focused transport assertion compared JSON strings
  with different object-key order for `invalid_limit`.
  - Change: compare normalized JSON values; reran focused tests.
- Reproduction: the first SDK build attempt found no package `build` script,
  and the added build then found option-shape and test-type configuration
  issues.
  - Change: added the package-local typecheck build, explicit query mapping,
  package tsconfig, and package-local Node Vitest config; reran build/tests.
- Reproduction: final full Go test initially failed only in legacy unit tests
  that use plain `http.ServeMux` with no page reader.
  - Change: isolated an explicitly deprecated plain-mux compatibility
  registration; RouteClassMux and any page-reader wiring remain IC-01.
  - Result: focused and full Go suites pass.

## Evidence Artifacts

- `spec.md`: ran; acceptance criteria and scope recorded.
- `plan.md`: ran; every acceptance criterion maps to verification.
- `validation.md`: ran; this transcript.
- `apps/server_core/internal/modules/catalog/transport/http_handler_test.go`:
  ran; fake-port httptest evidence.
- `packages/sdk-runtime/src/index.test.ts`: ran; typed client and OpenAPI
  lockstep evidence.

## Risks

- Composition root wiring of `PageReader` remains intentionally outside F-02
  and must be performed by the Milestone Orchestrator at integration/accept.
- The old plain-mux compatibility path must not be used as production
  composition; production route-class registration is explicitly IC-01.

## Handoff

- Current status: `quick_validation_passed`
- Next owner: Milestone Orchestrator
- Next action: Review this evidence, integrate F-01 adapter and composition
  wiring, then perform milestone acceptance/QA.
- Required files/evidence: feature brief, spec, plan, changed paths,
  validation transcript, post-commit stat.
- Blockers or open decisions: None.

## Post-commit stat transcript

- Command: `git show --stat HEAD`
- Status: Pass
- Evidence type: `ran`
- Actual transcript before evidence amend:

```text
commit 1853a79111152123b53001169dcb2c08e6a7a0e4
Author: leandrotcawork <leandrotca.work@gmail.com>
Date:   Tue Jul 14 10:55:50 2026 -0300

    feat(m-02/f-02): catalog IC-01 envelope routes + openapi + sdk

 .../F-02-catalog-api-envelope/plan.md              |  91 ++++++++
 .../F-02-catalog-api-envelope/spec.md              |  89 ++++++++
 .../F-02-catalog-api-envelope/validation.md        | 170 +++++++++++++++
 .../modules/catalog/transport/http_handler.go      | 242 ++++++++++++++++++++-
 .../modules/catalog/transport/http_handler_test.go | 202 +++++++++++++++++
 contracts/api/marketplace-central.openapi.yaml     | 159 +++++++++++++-
 packages/sdk-runtime/package.json                 |   3 +-
 packages/sdk-runtime/src/index.test.ts            |  28 ++-
 packages/sdk-runtime/src/index.ts                 |  73 ++++++-
 packages/sdk-runtime/tsconfig.json                |   8 +
 packages/sdk-runtime/vitest.config.ts             |   9 +
 11 files changed, 1053 insertions(+), 21 deletions(-)
```

The evidence append changes only this feature artifact and is amended into the
same intentional feature commit; the final SHA is reported in the handoff.
