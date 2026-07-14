# M-02 proportional QA validation result

- frozen_sha: `648adbdf4e4621700515278f44487c906babf876`
- qa_verdict: `PASS`

## Criteria

### M-02-C01 — One Oracle query per page: PASS

Command (from `apps/server_core`, absolute `GOCACHE`):

```text
go test ./internal/modules/internal_read/... ./internal/modules/catalog/... -run 'CatalogPage' -v
```

Key observed output (verbatim):

```text
=== RUN   TestCatalogPageUsesOneQueryForEveryListSize
=== RUN   TestCatalogPageUsesOneQueryForEveryListSize/limit-1
=== RUN   TestCatalogPageUsesOneQueryForEveryListSize/limit-50
=== RUN   TestCatalogPageUsesOneQueryForEveryListSize/limit-100
--- PASS: TestCatalogPageUsesOneQueryForEveryListSize (0.00s)
--- PASS: TestCatalogPageCursorChainIsGaplessAndNonOverlapping
```

The test passed for page sizes 1, 50, and 100 with one fake query independent of item count.

### M-02-C02 — Envelope and pagination conform to IC-01: PASS

Command:

```text
go test ./internal/modules/catalog/... -run 'Envelope|Cursor|Page' -v
```

Key observed output (verbatim):

```text
GET /catalog/products?limit=1 -> 200 {"as_of":"2026-07-14T12:00:00Z","items":[{"active":true,"cost":{"amount":null,"currency":"BRL","quality":["missing_cost"]},"current_price":{"amount":"12.90","currency":"BRL","quality":[]},"description":"PARAFUSO SEXTAVADO M8","ean":null,"internal_product_id":1,"reference":"ABC-123","sellable_stock":{"quality":[],"quantity":41}}],"next_cursor":"MQ==","page_size":1}
GET /catalog/products?limit=1&cursor=MQ== -> 200 {"as_of":"2026-07-14T12:01:00Z","items":[{"active":true,"cost":{"amount":null,"currency":"BRL","quality":["missing_cost"]},"current_price":{"amount":"12.90","currency":"BRL","quality":[]},"description":"PARAFUSO SEXTAVADO M8","ean":null,"internal_product_id":2,"reference":"ABC-123","sellable_stock":{"quality":[],"quantity":41}}],"next_cursor":"Mg==","page_size":1}
GET /catalog/products?limit=1&cursor=Mg== -> 200 {"as_of":"2026-07-14T12:02:00Z","items":[{"active":true,"cost":{"amount":null,"currency":"BRL","quality":["missing_cost"]},"current_price":{"amount":"12.90","currency":"BRL","quality":[]},"description":"PARAFUSO SEXTAVADO M8","ean":null,"internal_product_id":3,"reference":"ABC-123","sellable_stock":{"quality":[],"quantity":41}}],"next_cursor":null,"page_size":1}
--- PASS: TestCatalogPageRoutesFollowThreePageCursorChain (0.01s)
```

Envelope keys, decimal-string amount, RFC3339 `as_of`, null final cursor, and gapless/non-overlapping three-page chain passed.

### M-02-C03 — Error matrix enforced: PASS

Registered command:

```text
go test ./internal/modules/catalog/... -run 'Error|InvalidCursor|InvalidLimit|SourceUnavailable' -v
```

The selector did not match the validation test’s full name, so the allowed widened selector was also run:

```text
go test ./internal/modules/catalog/... -run 'TestCatalogPageRoutesValidateBeforePortCall|TestCatalogPageRoutesMapSourceAndDeadlineErrors' -v
```

Key observed output (verbatim):

```text
GET /catalog/products?cursor=%25%25%25garbage -> 400 {"error":"invalid_cursor"}
GET /catalog/products?limit=0 -> 400 {"allowed_range":"1..100","error":"invalid_limit"}
GET /catalog/products?limit=101 -> 400 {"allowed_range":"1..100","error":"invalid_limit"}
GET /catalog/products -> 503 {"error":"source_unavailable"}
--- PASS: TestCatalogPageRoutesValidateBeforePortCall (0.00s)
--- PASS: TestCatalogPageRoutesMapSourceAndDeadlineErrors (0.00s)
```

The source-unavailable body contains no driver detail.

### M-02-C04 — Unknown facts stay null with quality flags: PASS

Command:

```text
go test ./internal/modules/internal_read/... -run 'Nullable|Quality|Ambiguous' -v
```

Key observed output (verbatim):

```text
--- PASS: TestFakeReaderMissingStockStaysNilWithQualityFlag (0.00s)
--- PASS: TestCatalogPageMapsNullableFactsAndAmbiguousPrice (0.00s)
--- PASS: TestRequiredQualityFlagsRemainExplicit (0.00s)
--- PASS: TestMissingCostStaysNilWithQualityFlag (0.00s)
--- PASS: TestQualityFlagsAreStable (0.00s)
```

The focused tests cover missing stock/price/cost and duplicate active price behavior: null values plus quality flags, `ambiguous_price`, page 200, and no zero substitution.

### M-02-C05 — OpenAPI and sdk-runtime synchronized: PASS

Commands:

```text
git show --stat cbd4d009
cd packages/sdk-runtime && npm run build
```

Key `git show --stat cbd4d009` output (verbatim):

```text
contracts/api/marketplace-central.openapi.yaml     | 159 +++++++++++++-
packages/sdk-runtime/src/index.ts                  |  73 ++++++-
```

The same commit also contains the catalog handler and tests. SDK build output:

```text
> @marketplace-central/sdk-runtime@0.1.0 build
> tsc --noEmit
```

Exit code was 0. F-02 evidence confirms the old plain-mux listing registration is explicitly deprecated and production composition uses the paginated page reader route.

### M-02-C06 — Search route bounded with null cursor: PASS

Command:

```text
go test ./internal/modules/catalog/... -run 'Search' -v
```

Key observed output (verbatim):

```text
GET /catalog/products/search?q=PARAFUSO&limit=50 -> 200 {"as_of":"2026-07-14T13:00:00Z","items":[{"active":true,"cost":{"amount":null,"currency":"BRL","quality":["missing_cost"]},"current_price":{"amount":"12.90","currency":"BRL","quality":[]},"description":"PARAFUSO SEXTAVADO M8","ean":null,"internal_product_id":10,"reference":"ABC-123","sellable_stock":{"quality":[],"quantity":41}}],"next_cursor":null,"page_size":1}
--- PASS: TestCatalogSearchPageEnvelopeAndNoCachePolicy (0.01s)
```

The test passed its assertions for ≤50 results, ascending `internal_product_id`, exactly one query, and limit 51 returning `400 invalid_limit` (also observed in C03).

## Whole-suite gate

Command (from `apps/server_core`, absolute `GOCACHE`):

```text
go build ./... && go test ./...
```

Observed result: exit code 0; build completed and every Go package passed, including catalog transport, internal-read packages, integration, and unit suites.

Handoff: QA passed M-02 at frozen SHA `648adbdf4e4621700515278f44487c906babf876`; Milestone Orchestrator may close the milestone using this validation result.
