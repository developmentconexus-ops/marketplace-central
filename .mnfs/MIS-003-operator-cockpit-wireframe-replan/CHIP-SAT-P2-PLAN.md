# CHIP-SAT P2 batch plan — MIS-003

Scope is backend-only. No `apps/web/**`, `packages/feature-*`, `packages/ui`, or `packages/web-query` files enter any slice. The F-01 Vite proxy instruction is overridden by the mission Parallel Execution Plan.

Implementation order is strictly:

1. M-05 F-01 `aggregate-sync-endpoints`
2. M-06 F-02 `market-contract-module`

One pre-implementation contract issue requires hub resolution: the current orders schema cannot supply all required canonical fields or fulfillment filtering without reading `raw_provider_ref`, which is forbidden.

---

## 1. M-05 F-01 — aggregate-sync-endpoints

### Proposed HTTP shapes

- `GET /dashboard/summary?installation_id=…`
  - Nullable counters: `sync_errors`, `pending_links`, `below_margin`, `missing_gtin`, `orders_today`, `orders_7d`.
  - `last_sync_at` values nullable by module.
  - `degraded[]` source vocabulary: `listings`, `linkage`, `orders`, `sync`.
  - `as_of`.
- Existing `GET /orders`
  - Preserve `installation_id` and `limit`.
  - Add `cursor`, `filter.status`, `filter.fulfillment`, `filter.date_from`, `filter.date_to`, `q`.
  - Preserve `operationId: listMarketplaceOrders`.
  - Add optional `next_cursor`, `page_size`, `as_of` alongside existing `items`.
- `GET /orders/{provider_order_id}?installation_id=…`
- `GET /sync/runs`
  - Required `installation_id`; optional `cursor`, `limit`, `filter.module`, `filter.status`.
  - Fixed 90-day window.
  - Output maps stored `completed_at` to canonical `finished_at`.
  - Sort/keyset: `started_at DESC, operation_run_id DESC`.

### Pre-implementation gate — REQUEST required

`orders_marketplace_orders` (`0027`) stores provider status/timestamps, shipping ID, tags, and `raw_provider_ref`; payments store numeric totals but no currency. Neither `0027` nor `0033` stores:

- buyer nickname;
- canonical fulfillment state;
- order currency;
- a defined canonical order-status projection.

Therefore `buyer nickname`, `{total,currency}`, canonical `fulfillment`, and meaningful `filter.fulfillment` cannot be implemented from approved columns. Extracting them from `raw_provider_ref` in the read transport would leak provider payload semantics outside the adapter.

The hub must choose one:

1. Grant a migration block and authorize additive canonical columns populated during adapter ingestion; or
2. Ratify these facts as nullable/unavailable this mission and remove or narrow unsupported filters; or
3. Supply an existing canonical source not present in the inspected repository.

Until then, slices O2–O4 below are contract-ready but not implementable to the full brief.

### Slice cards

#### F01-O1 — Orders query and cursor grammar

- Goal: establish strict additive query parsing and an opaque versioned keyset cursor without changing current endpoint behavior.
- Write-set:
  - `apps/server_core/internal/modules/orders/domain/read_model.go` — new
  - `apps/server_core/internal/modules/orders/ports/read_query.go` — new
  - `apps/server_core/internal/modules/orders/ports/cursor.go` — new
  - `apps/server_core/internal/modules/orders/ports/cursor_test.go` — new
  - `apps/server_core/internal/modules/orders/transport/query.go` — new
  - `apps/server_core/internal/modules/orders/transport/query_test.go` — new
- Failing test first:
  - `TestDecodeOrderCursorRejectsMalformedWrongVersionAndTrailingJSON`
  - `TestParseOrderQueryRejectsMalformedDateAsInvalidFilter`
  - `TestParseOrderQueryRejectsUnknownOrRepeatedFilter`
  - `TestParseOrderQueryPreservesExistingLimitDefault`
  - Date range requires RFC3339/date grammar selected by the hub contract decision; `date_from > date_to` returns `invalid_filter`.
- Done:
  - Existing `installation_id` and `limit` remain valid.
  - Unknown filters are rejected, not ignored.
  - Cursor encodes the newest-first tie-break fields and returns `invalid_cursor`.
- Complexity: **standard** → Luna-high.

#### F01-O2 — Tenant-scoped orders keyset repository

- Goal: add filtered page/detail reads over `orders_marketplace_orders`, items/payments, and `orders_sankhya_linkage_events`.
- Write-set:
  - `apps/server_core/internal/modules/orders/ports/order_store.go`
  - `apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go`
  - `apps/server_core/internal/modules/orders/adapters/postgres/order_read_integration_test.go` — new
- Failing test first:
  - `TestListOrderPageUsesTenantInstallationAndKeyset`
  - `TestListOrderPageNewestFirstWithoutDuplicates`
  - `TestOrderFiltersStatusDateQueryAndFulfillment`
  - `TestGetOrderCannotCrossTenant`
  - `TestGetUnknownOrderReturnsNotFound`
  - Assert every query includes both `tenant_id` and `installation_id`.
- Done:
  - Stable cursor walk with no gaps/duplicates.
  - Date predicates use the ratified canonical order timestamp; null timestamps are not defaulted.
  - NF state is `linked` only when an exact `orders_sankhya_linkage_events` row exists; absence stays null, never “pending”.
  - Repository never exposes `raw_provider_ref`.
- Complexity: **complex** → Sol-low, because of keyset/filter SQL and child-row projection.

#### F01-O3 — Orders application and HTTP evolution

- Goal: evolve `/orders` in place and add order detail without mutation behavior.
- Write-set:
  - `apps/server_core/internal/modules/orders/application/list_service.go`
  - `apps/server_core/internal/modules/orders/application/list_service_test.go` — new
  - `apps/server_core/internal/modules/orders/transport/http_handler.go`
  - `apps/server_core/internal/modules/orders/transport/http_handler_test.go`
- Failing test first:
  - `TestListOrdersReturnsCursorEnvelopeAndCanonicalProjection`
  - `TestListOrdersMalformedCursorReturns400InvalidCursor`
  - `TestListOrdersMalformedDateReturns400InvalidFilter`
  - `TestGetUnknownOrderReturns404OrderNotFound`
  - `TestOrderTransportNeverSerializesRawProviderRef`
- Done:
  - Existing `/orders?installation_id=…&limit=…` still works.
  - New filters/cursor work alongside existing params.
  - Detail requires installation scope and returns `order_not_found`.
  - No order mutation/faturar/NF-write route is added.
- Complexity: **standard** → Luna-high.

#### F01-S1 — Sync-runs cursor and 90-day repository

- Goal: provide a tenant-scoped filtered keyset read over `integration_operation_runs`.
- Write-set:
  - `apps/server_core/internal/modules/integrations/ports/run_read.go` — new
  - `apps/server_core/internal/modules/integrations/ports/run_cursor.go` — new
  - `apps/server_core/internal/modules/integrations/ports/run_cursor_test.go` — new
  - `apps/server_core/internal/modules/integrations/ports/operation_run_store.go`
  - `apps/server_core/internal/modules/integrations/adapters/postgres/operation_run_repo.go`
  - `apps/server_core/internal/modules/integrations/adapters/postgres/operation_run_read_integration_test.go` — new
- Failing test first:
  - `TestListRunsTenantScopedNewestFirstWithCursor`
  - `TestListRunsExcludesStartedBeforeNinetyDays`
  - `TestListRunsFiltersStatusAndModule`
  - `TestRunningRunPreservesNullFinishedAt`
  - `TestMalformedRunCursorRejected`
- Done:
  - SQL contains `tenant_id`, `installation_id`, and the fixed cutoff.
  - `finished_at` remains null for running rows.
  - No timestamps are synthesized from zero values.
  - Proposed `filter.module` values are exact stored operation types such as `listings_refresh`, `pricing_fee_sync`, `listing_read`, and `order_read`; no guessed prefix mapping.
- Complexity: **complex** → Sol-low.

#### F01-S2 — Sync-runs application and transport

- Goal: expose `GET /sync/runs` with strict filters and error envelopes.
- Write-set:
  - `apps/server_core/internal/modules/integrations/application/operation_service.go`
  - `apps/server_core/internal/modules/integrations/application/operation_service_test.go`
  - `apps/server_core/internal/modules/integrations/transport/run_query.go` — new
  - `apps/server_core/internal/modules/integrations/transport/run_query_test.go` — new
  - `apps/server_core/internal/modules/integrations/transport/run_read_handler.go` — new
  - `apps/server_core/internal/modules/integrations/transport/run_read_handler_test.go` — new
- Failing test first:
  - `TestListSyncRunsRunningHasNullFinishedAt`
  - `TestListSyncRunsInvalidCursor`
  - `TestListSyncRunsInvalidStatusOrModule`
  - `TestListSyncRunsRequiresInstallation`
- Done:
  - Page envelope is `items/next_cursor/page_size/as_of`.
  - Existing `/integrations/installations/{id}/operations` remains unchanged.
- Complexity: **standard** → Luna-high.

#### F01-D1 — Orders dashboard counters

- Goal: publish today/7-day order counts as an orders application service.
- Write-set:
  - `apps/server_core/internal/modules/orders/ports/order_store.go`
  - `apps/server_core/internal/modules/orders/application/summary_service.go` — new
  - `apps/server_core/internal/modules/orders/application/summary_service_test.go` — new
  - `apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go`
  - `apps/server_core/internal/modules/orders/adapters/postgres/order_summary_integration_test.go` — new
- Failing test first:
  - `TestOrderSummaryCountsTodayAndSevenDaysTenantScoped`
  - `TestOrderSummaryDoesNotTreatUnknownTimestampAsToday`
- Done:
  - One tenant/installation-scoped aggregate query.
  - Counts derive from the hub-ratified canonical order timestamp.
- Complexity: **standard** → Luna-high.

#### F01-D2 — Product-link pending and missing-GTIN counters

- Goal: expose linkage counters without dashboard SQL.
- Write-set:
  - `apps/server_core/internal/modules/product_links/ports/summary_reader.go` — new
  - `apps/server_core/internal/modules/product_links/application/summary_service.go` — new
  - `apps/server_core/internal/modules/product_links/application/summary_service_test.go` — new
  - `apps/server_core/internal/modules/product_links/adapters/postgres/summary_reader.go` — new
  - `apps/server_core/internal/modules/product_links/adapters/postgres/summary_reader_integration_test.go` — new
- Failing test first:
  - `TestLinkageSummaryCountsUnresolvedAndEmptyEAN`
  - `TestLinkageSummaryScopesTenantAndInstallation`
- Done:
  - `pending_links` comes from unresolved/conflict linkage state.
  - `missing_gtin` counts blank `product_link_listing_snapshots.ean`.
  - Dashboard never queries product-link tables itself.
- Complexity: **standard** → Luna-high.

#### F01-D3 — Last-sync application projection

- Goal: expose latest successful/attempted sync timestamp per operation module for dashboard composition.
- Write-set:
  - `apps/server_core/internal/modules/integrations/ports/operation_run_store.go`
  - `apps/server_core/internal/modules/integrations/application/operation_service.go`
  - `apps/server_core/internal/modules/integrations/application/operation_service_test.go`
  - `apps/server_core/internal/modules/integrations/adapters/postgres/operation_run_repo.go`
  - `apps/server_core/internal/modules/integrations/adapters/postgres/operation_run_repo_integration_test.go`
- Failing test first:
  - `TestLatestRunsByModuleTenantScoped`
  - `TestLatestRunsKeepsUnknownTimestampNull`
- Done:
  - Repository returns latest stored facts; no “now” fallback.
  - Dashboard consumes the application method only.
- Complexity: **standard** → Luna-high.

#### F01-D4 — Honest dashboard composition

- Goal: compose installation, listings, linkage, orders, and sync application services with ADR-17 degradation.
- Write-set:
  - `apps/server_core/internal/modules/dashboard/domain/summary.go` — new
  - `apps/server_core/internal/modules/dashboard/ports/sources.go` — new
  - `apps/server_core/internal/modules/dashboard/application/service.go` — new
  - `apps/server_core/internal/modules/dashboard/application/service_test.go` — new
- Failing test first:
  - `TestSummaryOneFailedLinkageSourceReturnsPendingLinksNull`
  - `TestSummaryAllSourcesFailReturnsAllNullAndFullDegradedList`
  - `TestSummaryUnknownInstallationReturnsInstallationNotFound`
  - `TestSummaryDoesNotConvertUnavailableBelowMarginToZero`
- Done:
  - `degraded` is deterministic and deduplicated.
  - Linkage failure specifically yields `"pending_links": null` and includes `"linkage"`.
  - All-source failure returns 200 with all-null counters.
  - Installation lookup failure remains 404, not degraded.
  - No dashboard package imports Postgres or table names.
- Complexity: **standard** → Luna-high.

#### F01-D5 — Dashboard HTTP transport

- Goal: expose `GET /dashboard/summary`.
- Write-set:
  - `apps/server_core/internal/modules/dashboard/transport/http_handler.go` — new
  - `apps/server_core/internal/modules/dashboard/transport/http_handler_test.go` — new
- Failing test first:
  - `TestDashboardSummaryPartialFailureIs200`
  - `TestDashboardSummaryUnknownInstallationIs404`
  - `TestDashboardSummaryRequiresInstallation`
- Done:
  - JSON retains explicit null keys.
  - Errors use `installation_required` and `installation_not_found`.
- Complexity: **standard** → Luna-high.

#### F01-C1 — Orders OpenAPI and SDK compatibility lock

- Goal: evolve `/orders` in place and add detail without breaking the existing operation or SDK call.
- Write-set:
  - `contracts/api/marketplace-central.openapi.yaml`
  - `packages/sdk-runtime/src/index.ts`
  - `packages/sdk-runtime/src/index.test.ts`
- Failing test first:
  - SDK test proves legacy `listMarketplaceOrders("inst-1", 5)` still emits the same URL.
  - New `listOrders(options)` emits cursor/filter params.
  - New `getOrder(installationId, providerOrderId)` encodes both identifiers.
- Done:
  - `/orders` remains a single YAML path.
  - `operationId: listMarketplaceOrders` remains unchanged.
  - Existing parameters and `items: MarketplaceOrder[]` remain.
  - Additive page fields are optional in the schema but emitted by the new handler.
  - New SDK `listOrders` is an additive alias/new method, not a rename.
- Complexity: **standard** → Luna-high.

#### F01-C2 — Dashboard/sync OpenAPI and SDK lock

- Goal: add `getDashboardSummary` and `listSyncRuns`.
- Write-set:
  - `contracts/api/marketplace-central.openapi.yaml`
  - `packages/sdk-runtime/src/index.ts`
  - `packages/sdk-runtime/src/index.test.ts`
- Failing test first:
  - SDK URL encoding and nullable response type tests.
  - Contract test asserts sync `finished_at: null` and nullable dashboard counters.
- Done:
  - Only dashboard-summary and sync-runs paths/schemas change.
  - No mutation/protocolo, market, or shared preamble edits.
- Complexity: **standard** → Luna-high.

#### F01-R1 — Composition-root registration

- Goal: add only the pre-granted registrations and source wiring.
- Write-set:
  - `apps/server_core/internal/composition/root.go`
  - `apps/server_core/internal/composition/root_test.go`
- Failing test first:
  - Root route smoke proves `/dashboard/summary` and `/sync/runs` are mounted.
  - Existing `/orders` remains mounted once.
- Done:
  - No existing tables or unrelated module wiring changed.
- Complexity: **standard** → Luna-high.

#### F01-I1 — End-to-end integration contract

- Goal: prove the complete backend contract over real Postgres repositories.
- Write-set:
  - `apps/server_core/tests/integration/aggregate_sync_read_test.go` — new
- Failing test first:
  - `TestAggregateSyncReadContract`
  - Seeds both tenants, orders, link snapshots, listings, and runs.
  - Walks order and run cursors.
  - Exercises status/date filters, unknown order, 90-day exclusion, running/null-finished, and summary real counts.
- Done:
  - Tenant B never appears.
  - No provider raw payload appears.
  - Degraded behavior is covered by D4 unit tests; real-repository counts by this integration test.
- Complexity: **standard** → Luna-high.

### F-01 write-DAG

```text
Hub orders-shape decision
  └─ O1 → O2(complex) → O3 → C1
                     └──────→ D1 ─┐
S1(complex) → S2 → C2            │
          └──────→ D3 ───────────┤
D2 ──────────────────────────────┤
existing listings Summary ───────┤
                                 └─ D4 → D5 → R1 → I1
```

Parallel-eligible/file-disjoint:

- O1/O2, S1, and D2 are mutually file-disjoint.
- After O2, D1 may run alongside S2/D3 and D2.
- C1 and C2 are logically disjoint sections but edit the same OpenAPI/SDK files, so they are ordered.
- D4 waits for all four source services.
- R1 and I1 are terminal ordered slices.
- Nothing in F-01 may overlap F-02 because both features eventually write OpenAPI, SDK, and `root.go`; chip execution is sequential.

### F-01 contract satisfiability

| Path/section | Current state | Collision check | Verdict |
|---|---|---|---|
| `GET /dashboard/summary` | Free; no `/dashboard` path | Disjoint from CHIP-M03 and F-02 | **PASS — additive** |
| `GET /orders` | Occupied by `listMarketplaceOrders` | Same path as this feature; must be evolved, never duplicated | **PASS WITH LOCK — preserve path, operationId, params, response** |
| `GET /orders/{provider_order_id}` | Free | Existing deeper linkage paths use the same identifier but no direct GET at this key | **PASS — additive** |
| `GET /sync/runs` | Free; no `/sync` path | Disjoint from CHIP-M03 and F-02 | **PASS — additive** |
| Orders canonical schemas | Existing `MarketplaceOrder` and `ListMarketplaceOrdersResponse` occupied | Add nullable canonical projection/page fields; do not remove existing properties | **CONDITIONAL — storage/source decision required** |
| CHIP-M03 `/mutations...` | Free on this base, sibling-owned | No path/schema overlap | **PASS — do not touch** |

---

## 2. M-06 F-02 — market-contract-module

IC-04 is applied verbatim. “Manual” is represented by `source: manual`; it does not create a synthetic `manual_price` or aggregate “market price”. The value signals remain distinct fields:

`our_sale_price`, `winner_price`, `competitive_target`, `catalog_offer_price`, `catalog_stats`, plus manual provenance.

### Migration allocation

- `0043_market_observations.sql`
- `0044_market_references.sql`
- `0045` remains reserved and unused unless a concrete IC-04 requirement cannot fit the first two migrations.
- Do not create an empty `0045` migration.

Stored observation columns:

- `tenant_id`, `listing_id`
- four independent amount/currency pairs
- `competitive_status`
- `catalog_stats`
- `evidence_state`
- `captured_at NOT NULL`
- `source NOT NULL`
- lifecycle `created_at`
- PK `(tenant_id, listing_id, captured_at)`

Stored reference columns:

- `tenant_id`, `product_id`, `catalog_product_id`
- `match_state`, `match_method`
- `catalog_stats`, `evidence_state`
- `captured_at NOT NULL`, `source NOT NULL`
- lifecycle `created_at`
- PK `(tenant_id, product_id, captured_at)`

### Slice cards

#### F02-M1 — IC-04 migration contract

- Goal: create tenant-scoped append-only storage with mandatory provenance.
- Write-set:
  - `apps/server_core/migrations/0043_market_observations.sql` — new
  - `apps/server_core/migrations/0044_market_references.sql` — new
  - `apps/server_core/migrations/market_test.go` — new
- Failing test first:
  - `TestMarketMigrationsMatchIC04`
  - Assert exact columns/PKs, money-pair checks, enum checks, `source`/`captured_at NOT NULL`, no seed inserts, and no generic `market_price`.
- Done:
  - `evidence_state`: `observed|insufficient_market|no_price_evidence`.
  - `source`: `official_api|vendor|manual`.
  - Stored rows cannot have null source/captured time.
  - Synthetic no-evidence rows are never persisted.
- Complexity: **standard** → Luna-high.

#### F02-M2 — Domain and ports contract

- Goal: encode IC-04 entities, money/stats invariants, CollectorPort, and storage ports.
- Write-set:
  - `apps/server_core/internal/modules/market/domain/market.go` — new
  - `apps/server_core/internal/modules/market/domain/market_test.go` — new
  - `apps/server_core/internal/modules/market/ports/collector.go` — new
  - `apps/server_core/internal/modules/market/ports/store.go` — new
- Failing test first:
  - `TestObservationKeepsAllSignalsSeparate`
  - `TestCatalogStatsRequireAtLeastFiveSellers`
  - `TestStoredObservationRequiresSourceAndCapturedAt`
  - Compile-time CollectorPort signature assertion.
- Done:
  - CollectorPort signatures exactly match IC-04.
  - Unknown money remains nil.
  - No derived/winning promise field is introduced.
- Complexity: **standard** → Luna-high.

#### F02-M3 — Observation latest-row repository

- Goal: persist append-only observations and read the latest requested rows in input order.
- Write-set:
  - `apps/server_core/internal/modules/market/adapters/postgres/repository.go` — new
  - `apps/server_core/internal/modules/market/adapters/postgres/observation_repository.go` — new
  - `apps/server_core/internal/modules/market/adapters/postgres/observation_repository_integration_test.go` — new
- Failing test first:
  - `TestObservationLatestPerListingPreservesInputOrder`
  - `TestObservationRepositoryCannotCrossTenant`
  - `TestObservationWithoutSourceRejectedByDatabase`
  - `TestObservationSignalsRoundTripSeparately`
- Done:
  - Writes are INSERT-only, never upsert/update.
  - Every SQL query contains `tenant_id`.
  - Missing IDs are absent from repository results for application synthesis.
- Complexity: **complex** → Sol-low, due to latest-per-key/input-order SQL.

#### F02-M4 — Reference latest-row repository

- Goal: implement the corresponding reference storage/read contract.
- Write-set:
  - `apps/server_core/internal/modules/market/adapters/postgres/reference_repository.go` — new
  - `apps/server_core/internal/modules/market/adapters/postgres/reference_repository_integration_test.go` — new
- Failing test first:
  - `TestReferenceLatestPerProductPreservesInputOrder`
  - `TestReferenceRepositoryCannotCrossTenant`
  - `TestReferenceProvenanceIsMandatory`
- Done:
  - Latest capture per `(tenant_id, product_id)`.
  - Match-state and catalog-stats invariants survive round-trip.
- Complexity: **complex** → Sol-low.

#### F02-M5 — Test-only collection round-trip

- Goal: prove CollectorPort can feed storage without shipping a production implementation.
- Write-set:
  - `apps/server_core/internal/modules/market/application/collection_service.go` — new
  - `apps/server_core/internal/modules/market/application/collection_service_test.go` — new, package `application_test`
- Failing test first:
  - `TestCollectorPortRoundTripPersistsMandatoryProvenance`
  - `TestCollectorSignalsAreNeverMerged`
- Done:
  - Fake CollectorPort type exists only in `_test.go`.
  - Production service depends on the port but composition root does not instantiate a collector or collection scheduler.
- Complexity: **standard** → Luna-high.

#### F02-M6 — Honest-empty read application

- Goal: synthesize `no_price_evidence` only for requested IDs lacking stored rows.
- Write-set:
  - `apps/server_core/internal/modules/market/application/read_service.go` — new
  - `apps/server_core/internal/modules/market/application/read_service_test.go` — new
- Failing test first:
  - `TestVirginObservationReadReturnsNoPriceEvidenceWithNullSignals`
  - `TestMalformedListingIDReturnsItemLevelNoEvidence`
  - `TestUnknownProductIDReturnsNoEvidence`
  - `TestStoredSignalsRoundTripWithoutMerge`
  - `TestTwoHundredOneIDsReturnsTooManyIDs`
- Done:
  - Unknown/malformed IDs return 200 items in input order.
  - Synthetic items have null `captured_at`, null `source`, and all signal fields null.
  - Exactly 200 IDs are allowed; 201 returns `too_many_ids`.
- Complexity: **standard** → Luna-high.

#### F02-M7 — `/market` transport

- Goal: expose both IC-04 read endpoints.
- Write-set:
  - `apps/server_core/internal/modules/market/transport/query.go` — new
  - `apps/server_core/internal/modules/market/transport/query_test.go` — new
  - `apps/server_core/internal/modules/market/transport/http_handler.go` — new
  - `apps/server_core/internal/modules/market/transport/http_handler_test.go` — new
- Failing test first:
  - `TestObservationsRequiresInstallation`
  - `TestObservationsTwoHundredOneIDsReturns422`
  - `TestMalformedListingIDRemains200`
  - `TestReferencesTwoHundredOneIDsReturns422`
- Done:
  - `listing_ids`/`product_ids` parse as comma-separated ordered lists per validation examples.
  - Missing observation installation returns 400 `installation_required`.
  - 201 IDs return 422 `too_many_ids`.
- Complexity: **standard** → Luna-high.

#### F02-C1 — Market OpenAPI and SDK lock

- Goal: add the two real market operations and verbatim IC-04 schemas.
- Write-set:
  - `contracts/api/marketplace-central.openapi.yaml`
  - `packages/sdk-runtime/src/index.ts`
  - `packages/sdk-runtime/src/index.test.ts`
- Failing test first:
  - `listMarketObservations` and `listMarketReferences` URL/order encoding.
  - Honest-empty nullable type fixture.
  - Six-signal separation type fixture.
- Done:
  - No synthetic `market_price`.
  - Stored provenance is required by domain/storage; response provenance remains nullable for synthetic items.
  - No mutation/protocolo or F-01 section changes.
- Complexity: **standard** → Luna-high.

#### F02-C2 — Category-attribute contract reservation

- Goal: add the pre-assigned category-attribute OpenAPI/SDK contract for later M-06 F-01 without implementing its handler here.
- Write-set:
  - `contracts/api/marketplace-central.openapi.yaml`
  - `packages/sdk-runtime/src/index.ts`
  - `packages/sdk-runtime/src/index.test.ts`
- Failing test first:
  - SDK `getCategoryAttributes(categoryId)` URL encoding.
  - Type fixture pins provider-order-preserving `attributes[]`, nullable constraints, and fixed error codes.
- Done:
  - Adds `GET /listings/categories/{category_id}/attributes`, `operationId: getCategoryAttributes`.
  - Documents 404 `category_not_found` and 502 `provider_unavailable`.
  - No connector capability, cache, handler, or production adapter is added; those remain M-06 F-01.
- Complexity: **standard** → Luna-high.

#### F02-R1 — Market composition registration

- Goal: register repository-backed market reads only.
- Write-set:
  - `apps/server_core/internal/composition/root.go`
  - `apps/server_core/internal/composition/root_test.go`
- Failing test first:
  - Root smoke proves both `/market` routes are mounted.
  - Test confirms no collector/scheduler is wired.
- Done:
  - No CollectorPort implementation in production packages.
  - No background job or seed call.
- Complexity: **standard** → Luna-high.

#### F02-I1 — Market end-to-end and absence proof

- Goal: prove virgin-DB honesty, tenant isolation, provenance, limits, and forbidden-path absence.
- Write-set:
  - `apps/server_core/tests/integration/market_contract_test.go` — new
- Failing test first:
  - `TestMarketContractVirginDatabase`
  - `TestMarketContractRoundTripAndTenantIsolation`
  - `TestMarketContractTwoHundredOneIDs`
  - `TestMarketContractMalformedListingID`
- Done:
  - Virgin observations return all requested items as `no_price_evidence`.
  - Missing source insertion is rejected.
  - Grep evidence finds no production CollectorPort implementation, market seed, scraping dependency, or ML market/search call.
- Complexity: **standard** → Luna-high.

### F-02 write-DAG

```text
M1 → M2
      ├─ M3(complex) ─┐
      ├─ M4(complex) ─┼─ M6 → M7 → R1 → I1
      └─ M5 ──────────┘

M6/M7 → C1 → C2
```

Parallel-eligible/file-disjoint:

- After M2, M3 and M5 are file-disjoint.
- M4 is source-file-disjoint from M3 except for the shared repository constructor created by M3; it may start only after that common file is green.
- C1 and C2 edit the same OpenAPI/SDK files and are ordered.
- R1 waits for M7.
- I1 waits for repository, transport, contract, and root slices.
- F-02 starts only after all F-01 slices using OpenAPI/SDK/root have completed.

### F-02 contract satisfiability

| Path/section | Current state | Collision check | Verdict |
|---|---|---|---|
| `GET /market/observations` | Free | Disjoint from CHIP-M03 and F-01 | **PASS — additive** |
| `GET /market/references` | Free | Disjoint from CHIP-M03 and F-01 | **PASS — additive** |
| `GET /listings/categories/{category_id}/attributes` | Free | Under existing `/listings` prefix, but explicitly granted to M-06; no current route occupies it | **PASS — additive reservation for F-01** |
| Market schemas/operationIds | Absent | No existing `MarketObservation`, `MarketReference`, `listMarketObservations`, or `listMarketReferences` | **PASS** |
| CHIP-M03 `/mutations...` | Free on this base, sibling-owned | No overlap | **PASS — do not touch** |
| F-01 dashboard/orders/sync | Dashboard/sync free; orders occupied as above | Shared YAML/SDK files but disjoint pre-assigned sections; serialized inside this chip | **PASS** |

Migration check:

- Current maximum migration is `0037`.
- No tracked or untracked `0043*`, `0044*`, or `0045*` files exist.
- **Verdict: `0043–0045` are free and reserved for F-02.**

---

## Expected additive `root.go` contract-lock diff

Exact anticipated imports:

```go
dashboardapp "marketplace-central/apps/server_core/internal/modules/dashboard/application"
dashboardtransport "marketplace-central/apps/server_core/internal/modules/dashboard/transport"
marketpostgres "marketplace-central/apps/server_core/internal/modules/market/adapters/postgres"
marketapp "marketplace-central/apps/server_core/internal/modules/market/application"
markettransport "marketplace-central/apps/server_core/internal/modules/market/transport"
```

F-01 additive wiring/registration:

```go
productLinkSummarySvc := productlinksapp.NewSummaryService(productLinkSnapshotRepo, productLinkCandidateRepo)
ordersSummarySvc := ordersapp.NewSummaryService(ordersRepo)

integrationstransport.NewRunReadHandler(operationSvc).Register(mux)

dashboardSvc := dashboardapp.NewService(
	installationSvc,
	listingSvc,
	productLinkSummarySvc,
	ordersSummarySvc,
	operationSvc,
	time.Now,
)
dashboardtransport.NewHandler(dashboardSvc).Register(mux)
```

The existing line remains once and is not duplicated:

```go
orderstransport.NewHandler(ordersImportSvc, ordersListSvc).Register(mux)
```

F-02 additive wiring/registration:

```go
marketRepo := marketpostgres.NewRepository(pool, cfg.DefaultTenantID)
marketReadSvc := marketapp.NewReadService(marketRepo, time.Now)
markettransport.NewHandler(marketReadSvc).Register(mux)
```

There must be no line resembling:

```go
marketapp.NewCollectionService(...)
marketbackground.NewScheduler(...)
```

---

## Open questions and REQUEST-worthy items

1. **Blocking — orders canonical source.** Request a hub ruling on buyer, currency/total, canonical status, fulfillment, and fulfillment filtering. Current tables cannot satisfy them without forbidden raw-provider decoding. A migration need must receive a new hub-assigned block; F-01 may not use `0043–0045`.

2. **Contract-lock ratification — dashboard field names.** The brief pins semantics but not the exact JSON nesting. Ratify the proposed nullable counter names and `degraded` vocabulary before F01-C2.

3. **Sync module grammar.** The plan uses exact stored `operation_type` values for `filter.module`; if the UI expects broader names such as `listings` or `orders`, the mapping must be pinned rather than inferred.

4. **Hub deferral — Vite proxy rows.** `/dashboard` and `/sync` proxy additions are explicitly excluded from CHIP-SAT. Report them as a `REQUEST`/deferred hub item for the frontend-seam owner. `/market` is likewise assumed to be supplied by M-02; this chip does not verify by editing Vite.

5. **Category-attribute runtime gap is intentional.** F-02 adds only its pre-assigned OpenAPI/SDK reservation. M-06 F-01 later owns connector mapping, category cache, handler, 404/502 behavior, and live route implementation.

6. **Migration `0045`.** Leave it unused unless implementation reveals a concrete third schema migration. Any requirement beyond `0045` is a hub `REQUEST`.

7. **No dependencies.** Any proposed scraping/market client library, collector adapter, seed path, or scheduler is a scope defect, not a dependency request.