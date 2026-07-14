# F-02 catalog-api-envelope

```yaml
id: F-02
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-02
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-2
lifecycle_scope: feature
```

## Problem

The catalog transport still exposes an unpaginated composition-oriented list
response. F-01 now provides a page port, so the existing `/catalog/*`
namespace needs a stable IC-01 envelope for listing and search, with transport
validation and bounded error mapping.

## Requirements

- `GET /catalog/products` accepts optional `cursor` and `limit` (`1..100`,
  default `50`) and calls `CatalogPageReader.ListCatalogProductFacts`.
- `GET /catalog/products/search` accepts `q` and optional `limit` (`1..50`,
  default `50`) and calls `CatalogPageReader.SearchCatalogProductFacts`.
- Both success responses use `{items, next_cursor, page_size, as_of}`. Search
  always returns `next_cursor: null`.
- Invalid cursor and limit are rejected before the port call. Port source
  failures map to `503 source_unavailable`; route-class deadline expiry maps to
  `504 deadline_exceeded`, without driver detail.
- `Cache-Control: no-cache` maps to an internal-read `FreshnessPolicy` with
  `MaxAge=0` in the request context. No cache is added.
- OpenAPI and `sdk-runtime` describe and call the same typed page envelope.

## Non-Goals

- No edits to F-01 ports, adapters, application services, or composition root.
- No caching, frontend changes, Oracle logic, or business logic in handlers.
- No changes to unrelated catalog operations.

## Design

The catalog handler receives an optional `internal_read/ports.CatalogPageReader`
and maps query parameters into the F-01 cursor and limit inputs. It maps the
port page and typed read errors into transport DTOs. The existing list path is
repurposed as the paginated route, removing the old unpaginated flow. Both
routes are explicitly registered as `interactive`; `RouteClassMux` supplies
the 15-second context deadline.

The SDK adds `CatalogProductFact`, `CatalogProductFactPage`, and typed list and
search methods. OpenAPI adds matching schemas, parameters, paths, and IC-01
error responses in the same change.

## Edge Cases

- Missing cursor means the first page; an explicitly empty or malformed cursor
  is `400 {"error":"invalid_cursor"}` and does not call the port.
- Missing limit uses the operation default; non-integer and out-of-range
  values are `400 invalid_limit` with the operation's allowed range.
- Search ignores a returned cursor and always emits JSON `null`.
- Missing facts remain JSON `null`; quality arrays remain present.
- A no-cache header is mapped without introducing a cache or changing the
  page-port interface.

## Acceptance Criteria

- Criterion: Catalog list and search expose the IC-01 page envelope and route
  namespace, with cursor/limit validation and no business logic in transport.
  - Traces to milestone criterion ID: M-02-C01
  - Proven by: `GOCACHE=.gocache go test ./...` and transport httptest transcript.
- Criterion: The full IC-01 transport error matrix is redacted and the
  interactive route deadline produces `deadline_exceeded`.
  - Traces to milestone criterion ID: M-02-C02
  - Proven by: transport httptest transcript and route-deadline tests in
    `GOCACHE=.gocache go test ./...`.
- Criterion: OpenAPI and `sdk-runtime` remain in lockstep with typed catalog
  page methods.
  - Traces to milestone criterion ID: M-02-C03
  - Proven by: `npm run build` in `packages/sdk-runtime` and OpenAPI/SDK tests.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Write `plan.md`, implement, and record evidence.
- Required files/evidence: plan, changed paths, validation transcript
- Blockers or open decisions: None.
