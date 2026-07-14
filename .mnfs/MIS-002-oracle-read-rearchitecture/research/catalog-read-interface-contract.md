# Interface Contract — Catalog Paginated Read, Cache & Batch Semantics

```yaml
id: IC-01
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-002
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: support
```

## Boundary

HTTP API (`apps/server_core` transport) ↔ TS SDK (`packages/sdk-runtime`) ↔ web (`apps/web`); plus the internal cache/batch semantics every Oracle-backed port shares. Multiple workers (M-02, M-03, M-04, M-05) touch this seam.

## Why This Contract Exists

Pagination envelope, cursor semantics, `as_of`, TTL classes, route classes, IN-list chunking, and error codes would otherwise be invented independently by 4 milestones.

## Resources Or Entities

`CatalogProductFact` (one row per product page item):

```json
{
  "internal_product_id": 12345,
  "reference": "ABC-123",
  "description": "PARAFUSO SEXTAVADO M8",
  "ean": null,
  "active": true,
  "sellable_stock": { "quantity": 41.0, "quality": [] },
  "current_price": { "amount": "12.90", "currency": "BRL", "quality": [] },
  "cost": { "amount": null, "currency": "BRL", "quality": ["missing_cost"] }
}
```

Nullability: unknown facts are `null` + a `quality` flag entry. NEVER `0`. `amount` is a decimal string. `ean` stays null until a governed barcode source exists (M-09 rule preserved).

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| `GET /catalog/products` | operator opens list / next page | `?cursor=<opaque>&limit=<1..100, default 50>` | `200` page envelope | sort: `internal_product_id` ascending (keyset), stable |
| `GET /catalog/products/search` | debounced search | `?q=<text>&limit=<1..50, default 50>` | `200` page envelope, `next_cursor` always `null` | bounded FETCH FIRST; sort `internal_product_id` asc |
| force refresh (any Oracle-backed GET) | operator refresh control | header `Cache-Control: no-cache` | fresh read, new `as_of` | maps to `FreshnessPolicy.MaxAge=0`; bypasses L2 |

### Page envelope (Required Outputs)

```json
{
  "items": [ CatalogProductFact ],
  "next_cursor": "MTIzNDU=",
  "page_size": 50,
  "as_of": "2026-07-13T18:22:05Z"
}
```

- `next_cursor`: opaque base64 of the last item's keyset key (`internal_product_id`). `null` = last page. Clients MUST treat it as opaque.
- `as_of`: RFC3339 UTC — moment facts were read from Oracle (fresh) or originally read (cache hit). Present on EVERY Oracle-backed response body in this mission (catalog, stock-risks, linkage candidates).

### Required Inputs

`cursor` optional (absent = first page); `limit` optional. Invalid cursor (not base64 / not a positive int) → 400 `invalid_cursor`.

## Enums And Statuses

`quality` flag values reuse existing domain set (`missing_stock`, `missing_price`, `missing_cost`, `ambiguous_price`, ...). No new enum values introduced without adding here first.

## Error Cases / Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| invalid/malformed cursor | 400 | `invalid_cursor` | body `{"error":"invalid_cursor"}` |
| limit out of range | 400 | `invalid_limit` | states allowed range |
| Oracle unavailable | 503 | `source_unavailable` | existing typed error; no fallback data; no raw driver detail in body |
| interactive deadline exceeded (15s) | 504 | `deadline_exceeded` | context timeout at transport |
| batch cap exceeded (`ImportMarginInputs.Limit`>200) | 422 | `limit_exceeded` | body names the cap; NEVER silent truncation |
| ambiguous price (multi-row) | 200 item-level | quality flag `ambiguous_price` + null amount | per-item, does not fail the page |

No downstream feature returns an error case absent from this matrix; new case = new row here first.

## Cache Semantics (L2 server + L1 TanStack)

| Data class | L2 TTL (env-tunable default) | TanStack staleTime | Force-refresh |
| --- | --- | --- | --- |
| catalog page / taxonomy / search | 5min | 5min | yes |
| sellable stock | 45s | 45s | yes |
| price / cost | 2min | 2min | yes |
| Sankhya linkage (candidates/descendants/config data) | NEVER cached (config validation result MAY be cached 5min) | staleTime 0 | n/a |

- L2 key: `(port method, canonical params)`. No tenant component: Oracle read ports carry no tenant dimension — ERP facts are installation-global (recorded in mission `Accepted assumptions`; revisit if per-tenant Oracle scoping is ever introduced). Concurrent identical misses collapse via singleflight.
- Evict-on-mutation: confirming a linkage, applying a stock action, or any write that invalidates a fact class evicts the matching L2 keys AND the client invalidates matching queryKeys, per the fixed crosswalk below (features reference it, never invent labels).

### Invalidation Crosswalk (server fact class ↔ client queryKey namespace)

| Mutation | Server `InvalidateClass` | Client `invalidateQueries` |
| --- | --- | --- |
| linkage confirm | `catalog` | `['linkage']` + `['catalog']` |
| product edit (DORMANT) | `catalog` | `['catalog']` |
| stock-affecting action | `inventory` | `['inventory']` |
| margin-input import | `pricecost` | `['catalog']` + `['profitability']` |

Server fact classes: `catalog`, `inventory`, `pricecost` (linkage has no class — never cached).

DORMANT row (product edit): M-05 removed the legacy ProductsPage — the only product-edit write surface — so this row currently has no implementation target. It is a forward obligation, not dead: any reintroduced product-edit surface MUST invalidate server class `catalog` and client `['catalog']`.
- `MaxAge=0` (from `Cache-Control: no-cache`) bypasses L2 and repopulates it.

## Batch & Route-Class Rules (server-internal, cross-milestone)

- Route classes declared at registration: `interactive` → 15s context deadline; `batch` → 120s. Batch routes: `/profitability/margin-inputs/import`, `/profitability/profit-snapshots/calculate`, `/orders/import`, `/product-links/*/imports|generations`, `/pricing/simulations/batch`, fee syncs.
- Batch Oracle work acquires one of 4 semaphore permits before touching the pool (pool stays 12).
- IN-list chunk size 500 (ORA-01795 guard). Chunk results merge in the adapter; partial chunk failure fails the whole call (typed error), never partial-silent results.
- `ImportMarginInputs.Limit` ceiling 200 → exceeding = 422 `limit_exceeded`. `GetSalesHistory` row cap 5000 → result carries an explicit `truncated=true` marker (200; report marked partial) — NEVER silent truncation; a 422 here would block whole profitability reports the operator cannot shrink.

## Persistence Expectations

None new — Oracle read-only; Postgres untouched by this contract.

## Timestamp And ID Semantics

- `internal_product_id`: positive int64 (CODPROD); positivity enforced (existing canonical-identity rule).
- All timestamps RFC3339 UTC.
- Sankhya linkage candidates list: ordered by match score descending, tie-break `internal_product_id` ascending (existing reader semantics preserved by M-03 F-03).

## Compatibility Rules

- Envelope is the base shape; later features may ADD fields to items or envelope, never rename/reshape `items`/`next_cursor`/`page_size`/`as_of`.
- OpenAPI (`contracts/api/marketplace-central.openapi.yaml`) and `packages/sdk-runtime` change in the SAME commit as the handler (repo rule).

## Route Namespace

- Server: `/catalog/*` owned by catalog transport; no new prefixes introduced by this mission.
- Client: existing feature-package pages; TanStack queryKey namespace: `['catalog', ...]`, `['inventory', ...]`, `['linkage', ...]`, `['profitability', ...]` — reserved per feature package; no other namespaces.

## Must Preserve

- nil + quality flags (never zero/default); typed `source_unavailable`; credential/DSN redaction (`safeOracleCause` discipline) in all new code paths.

## Must Not Decide In Feature Execution

- Envelope field names, cursor encoding, TTL classes, route-class deadlines, chunk size, error matrix rows, queryKey namespaces.

## Validation Impact

Criteria in mission/milestone contracts reference this contract's exact shapes (M-02-C01..C03, M-03-C01..C03, M-04-C01..C02, M-05-C01..C02).
