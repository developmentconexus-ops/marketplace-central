# F-02 listings-read-api implementation plan

## Codebase Anchors

1. **Catalog cursor-pagination precedent to mirror**

   - `apps/server_core/internal/modules/internal_read/ports/catalog_page.go:12-24` defines the typed `invalid_cursor` sentinel/error, and `apps/server_core/internal/modules/internal_read/ports/catalog_page.go:26-55` defines the zero-value first-page cursor plus strict standard-base64 encode/decode.
   - `apps/server_core/internal/modules/internal_read/adapters/oracle/catalog_page.go:20-32` validates the cursor and requests `limit+1`; `apps/server_core/internal/modules/internal_read/adapters/oracle/catalog_page.go:78-82` trims the look-ahead row and emits the last returned key; `apps/server_core/internal/modules/internal_read/adapters/oracle/catalog_page.go:86-153` applies the keyset predicate before a matching `ORDER BY` and bounded fetch.
   - `apps/server_core/internal/modules/catalog/transport/http_handler.go:174-187` rejects duplicate/undecodable cursor query values before the port call, `apps/server_core/internal/modules/catalog/transport/http_handler.go:206-210` fixes the `items/next_cursor/page_size/as_of` envelope, and `apps/server_core/internal/modules/catalog/transport/http_handler.go:235-266` encodes only a real next cursor.
   - `apps/server_core/internal/modules/catalog/transport/http_handler_test.go:52-90` proves a complete three-page walk with a null terminal cursor; `apps/server_core/internal/modules/internal_read/ports/catalog_page_test.go:8-31` proves round-trip and rejects empty, malformed, zero, negative, and non-integer cursors.
   - F-02 mirrors those mechanics with a listings-owned cursor containing the complete ordered key (`title`, canonical `listing_id`), and a separate group cursor containing the complete group key. It does not reuse the single-integer catalog cursor.
   - **Contract reconciliation:** catalog page errors at `apps/server_core/internal/modules/catalog/transport/http_handler.go:276-293` use legacy flat `{"error":"invalid_cursor"}` JSON. F-02 mirrors catalog cursor/keyset and envelope semantics, not that body; IC-02 requires the nested envelope.

2. **`product_links` read model and exact link identity**

   - `apps/server_core/internal/modules/product_links/domain/product_link.go:8-14` contains link states. IC-02 consumes only `unresolved|conflict|resolved|rejected`, so internal `none` or an absent row must normalize to `unresolved`, never leak.
   - `apps/server_core/internal/modules/product_links/domain/product_link.go:28-39` carries `InternalProductID`, `InternalProductName`, and `InternalReferenceCode`; only a resolved link may expose `product_id`.
   - `apps/server_core/internal/modules/product_links/ports/workflow_store.go:9-13` is the published workflow read surface. The exact single-link lookup is `apps/server_core/internal/modules/product_links/adapters/postgres/link_candidate_repo.go:200-240`; its SQL at lines 204-214 binds tenant, installation, provider item, and provider variation.
   - `apps/server_core/migrations/0025_product_link_workflows.sql:1-18` fixes the physical key and tenant/installation index.
   - Seller SKU is not in `product_links`: `apps/server_core/internal/modules/product_links/domain/listing_snapshot.go:5-18` and `apps/server_core/migrations/0022_product_links_listing_snapshots.sql:1-16` place it in the product-links-owned `product_link_listing_snapshots` projection. Its reader is tenant/installation scoped at `apps/server_core/internal/modules/product_links/adapters/postgres/listing_snapshot_repo.go:75-130`.
   - F-02 therefore uses tenant-scoped `LEFT JOIN` projections to both tables: `product_links` for state/product/title and `product_link_listing_snapshots` for seller SKU. Every `ON` clause repeats tenant, installation, provider item, and variation.
   - **Physical-shape mismatch:** F-01 listings use literal `-` for no variation, while product-links tables use `''` (`0022_product_links_listing_snapshots.sql:6` and `0025_product_link_workflows.sql:6`). Normalize the variation explicitly in the join; never broaden to item-only matching.

3. **F-01 listings seam**

   - No listings package exists in this checkout yet; F-01 is represented by its accepted plan.
   - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-01-listings-read-spine/F-01-listings-module-ingestion/plan.md:121-157` creates migration 0036 per D-1, the canonical domain entity, and an F-02 keyset index while expressly excluding F-02 filters/cursors/read DTOs.
   - The same plan at lines 227-265 creates the PostgreSQL repository only for the completed-pull writer and expressly excludes read-query ports.
   - The same plan at lines 276-360 creates the listings handler, registers only `POST /listings/refresh`, composes the module, and locks only refresh OpenAPI/SDK.
   - D-5 requires one modality registry in the listings module, consumed by F-01 recognition and F-02 label assembly.
   - The planned F-01 repository has no F-02 read surface or timeline source. The next section is therefore a mandatory shared-seam addition for F-01, not optional abstraction.

4. **Cost fact and pricing-policy sources**

   - `apps/server_core/internal/modules/internal_read/domain/internal_cost.go:13-27` models nullable latest/as-of cost and its quality/source.
   - `apps/server_core/internal/modules/internal_read/ports/batch_reader.go:9-13` publishes bounded `GetCostFactsByIDs`; a missing map value is explicitly unknown.
   - `apps/server_core/internal/modules/internal_read/adapters/oracle/batch_reader.go:40-95` performs the latest-cost Oracle read in chunks, preserves null, and labels missing cost rather than defaulting. Composition already wraps it in the price/cost cache at `apps/server_core/internal/composition/root.go:441-448`.
   - `apps/server_core/migrations/0003_marketplaces.sql:13-25` stores `min_margin_percent`; `apps/server_core/internal/modules/marketplaces/domain/policy.go:3-16` is the policy shape. Tenant-scoped policy reads exist only as all/by-ID lists at `apps/server_core/internal/modules/marketplaces/ports/repository.go:9-15` and `apps/server_core/internal/modules/marketplaces/adapters/postgres/repository.go:91-167`.
   - Installation-to-account linkage is physical at `apps/server_core/migrations/0016_integrations_foundation.sql:225-235`, but no published `GetPolicyForInstallation` exists. F-02 needs a narrow installation-scoped read, not list-all-and-guess.
   - Existing pricing calculation at `apps/server_core/internal/modules/pricing/application/service.go:34-51` deducts commission, fixed fee, and shipping before comparing margin. IC-02/F-02 only say listing price + latest cost + min-margin policy and do not settle whether those deductions apply.
   - **Blocking contract/verification gap:** costs are in Oracle while listing/link/policy rows are PostgreSQL, and no local latest-cost table exists. The endpoint cannot both read authoritative latest cost and truthfully claim the complete summary executed as one SQL query without a ruling or owned read projection.

5. **Transport error envelope and route registrar**

   - `apps/server_core/internal/modules/product_links/transport/http_handler.go:67-84` and `apps/server_core/internal/modules/integrations/transport/http_handler.go:38-55` are the exact nested `{error:{code,message,details}}` idiom.
   - Existing 400/404 mapping is at `apps/server_core/internal/modules/product_links/transport/http_handler.go:87-101`. F-02 uses typed listings errors rather than string-prefix inference while preserving the envelope.
   - `apps/server_core/internal/platform/httpx/route_deadline.go:13-18` is the registrar; method-aware GET registration is demonstrated at `apps/server_core/internal/modules/catalog/transport/http_handler.go:63-89`.
   - The composition root constructs/registers adjacent modules at `apps/server_core/internal/composition/root.go:349-371`. F-01 adds one listings handler there; F-02 extends that same handler's `Register` and must not mount another prefix handler.
   - Register static `GET /listings/by-product` and `GET /listings/summary`, then `GET /listings/{listing_id}` and exact `GET /listings`; retain F-01 `POST /listings/refresh` in the same registrar.

6. **OpenAPI and hand-maintained SDK**

   - `contracts/api/marketplace-central.openapi.yaml:6-55` shows GET list cursor/limit parameters and a typed page response; lines 1655-1672 fix the page envelope schema.
   - `contracts/api/marketplace-central.openapi.yaml:4292-4321` is the reusable nested error schema with open details.
   - `packages/sdk-runtime/src/index.ts:177-214` shows page/item/options types, lines 994-1018 show typed GET/error and query construction, and lines 1054-1067 show typed list methods.
   - `packages/sdk-runtime/src/index.test.ts:42-66` captures exact GET URL/method/page behavior; lines 26-40 demonstrate OpenAPI/SDK source parity assertions.
   - D-15 is binding: OpenAPI first, then hand-maintained `index.ts` and tests in the same endpoint-exposing commit; never claim generation.

7. **Tenant scoping**

   - `ARCHITECTURE.md:37` requires tenant ownership on business data.
   - Repository instances receive tenant at `apps/server_core/internal/modules/catalog/adapters/postgres/repository.go:17-24`; reads bind it first at lines 32-36 and 85-89.
   - Link reads bind tenant + installation + full listing identity at `apps/server_core/internal/modules/product_links/adapters/postgres/link_candidate_repo.go:204-214`; seller-SKU reads bind tenant + installation at `apps/server_core/internal/modules/product_links/adapters/postgres/listing_snapshot_repo.go:79-89`.
   - Every F-02 query, CTE, subquery, join, group, getter, timeline, summary, cost-ID selection, and EXPLAIN fixture binds repository tenant. A predicate on `listings` alone is insufficient: both product-links joins bind their own tenant and installation.
   - Integration tests seed the same installation/provider identities for tenant B and prove tenant A list/get/group/summary/q/filter requests neither return nor count tenant B.

## Shared Seam Requirements for F-01

These must be folded into F-01 before F-02 starts. They have one named consumer, F-02 `application.ReadService`, so they are not speculative one-implementation interfaces.

1. **F-01 must expose a listings-owned read repository port because F-02 needs contract-shaped tenant-scoped reads without importing PostgreSQL infrastructure.**

   - Create `apps/server_core/internal/modules/listings/ports/read_repository.go`.
   - Expose `ListingReadRepository` with:
     - `ListListingRows(ctx, ListingQuery) (ListingRowPage, error)` — keyset walk in title/listing-ID order, `limit+1` look-ahead, all non-cost filters and q in SQL, joined link/SKU facts and nullable listing facts.
     - `ListListingGroupRows(ctx, ListingGroupQuery) (ListingGroupRowPage, error)` — cursor over group keys and complete children for only selected groups, without N+1.
     - `GetListingRow(ctx, ListingKey) (ListingReadRow, bool, error)` — exact tenant/composite lookup with the same projection as list.
     - `GetListingsSummary(ctx, SummaryQuery) (ListingSummaryRow, error)` — exactly one PostgreSQL aggregate once the cost/policy ruling supplies queryable inputs.
     - `ListListingTimeline(ctx, ListingKey, limit int) ([]TimelineEventRow, error)` — full tenant/composite identity, newest first with deterministic tie-break, maximum 10.
     - `ListResolvedProductIDs(ctx, installationID) ([]int64, error)` only if the accepted cost design requires a bulk pre-read; omit it if costs become an owned SQL projection.
   - Inputs contain only validated values: installation, typed filters, trimmed q, decoded cursor, limit, and ratified cost/policy input. Raw URL values never enter the repository.
   - Row types preserve pointers for modality, paired price fields, quantity, sync-error members, quality, sales, fetched time, SKU, product ID/title, cost, policy, and below-margin. No `COALESCE` may fabricate 0, false, BRL, or an empty contract object.

2. **F-01 must make its single PostgreSQL repository satisfy both completed-pull writer and read port because F-02 needs one tenant-bound adapter.**

   - Keep `apps/server_core/internal/modules/listings/adapters/postgres/repository.go` as the concrete repository with compile-time assertions for each named port.
   - Centralize one joined projection/scanner used by list, group children, and get so shapes cannot drift.
   - Construct listing ID from stored installation/provider-listing/variation; use the identical expression/order for response and keyset comparison.
   - Join both product-links tables on tenant, installation, provider item, and normalized variation. Never item-only join.
   - Normalize absent/none to unresolved; product ID/title only for resolved; seller SKU only from snapshot and null when empty/absent.

3. **F-01 must create the exact keyset index because F-02 needs M01-C09 without offset scans.**

   - Migration `apps/server_core/migrations/0036_listings.sql` creates a named B-tree on `(tenant_id, installation_id, title, provider_listing_id, variation_id)`.
   - SQL ordering/cursor predicate use the same collation and tuple. With installation fixed, provider-listing + variation are the stored listing-ID suffix.
   - No speculative filter indexes. Add one only after 2,000-row EXPLAIN names a real bottleneck/consumer.

4. **F-01 must provide a durable listing-scoped timeline because F-02 cannot derive the last ten events from current row state or installation-wide operation runs.**

   - Before migration 0036 freezes, ratify event kinds/messages and whether insert, refresh, close, and failure create events.
   - Then create a listings-owned event source (for example a separate `listing_sync_events` table, not new listing columns) keyed by tenant + installation + provider listing + variation + event ID, with `at/kind/message_pt` and index `(tenant_id, installation_id, provider_listing_id, variation_id, at DESC, event_id DESC)`.
   - F-01 writes events in the same completed-pull transaction as affected listing state. F-02 only reads them. An always-empty fallback or installation-wide history is forbidden.

5. **F-01 must publish the D-5 modality lookup because F-02 needs the same registry as ingestion.**

   - Keep one lookup under `apps/server_core/internal/modules/listings/domain` returning label + recognized for the seven D-5 codes.
   - F-02 maps nullable stored code through it; null/unrecognized produces JSON null. No transport/SDK label map.

6. **F-01 must leave one handler/composition mount because duplicate `/listings` registrars conflict.**

   - F-01 registers one handler in composition and initially only `POST /listings/refresh`.
   - Its constructor/config permits F-02 `ReadService` injection; F-02 edits the same `Register` to add four method-aware GETs.

7. **F-01 must declare known F-02 module dependencies because F-02 consumes product-links, internal-read, and policy facts.**

   - The final listings entry in `contracts/governance/modules.json` contains F-01 dependencies plus ratified F-02 dependencies: `product_links`, `internal_read`, and `marketplaces` or `pricing` according to the policy ruling.
   - Listings adapters import only published domain/application/port surfaces. The IC-02 SQL read join is limited to projection and never mutates product-links tables.

8. **F-01/milestone must resolve the cost/policy query seam because F-02 needs cost-aware filters and honest single-query summary.**

   - Add a marketplace-owned method such as `GetPricingPolicyForInstallation(ctx, installationID) (Policy, bool, error)` using tenant + `marketplace_accounts.integration_installation_id`. Multiple matches fail honestly unless a selection rule is ratified.
   - Use `internal_read/ports.BatchReader.GetCostFactsByIDs` only if C09 is clarified to allow that Oracle read plus one PostgreSQL aggregate. A missing reader is not an empty known cost set.
   - If C09 literally means one database query for the endpoint, architecture must first provide a tenant-scoped local latest-cost projection. F-02 does not invent a cache/table or use profitability history.

## Slices

### Slice 1: Lock filter grammar, read model, and opaque cursor foundations

- **Complexity:** complex (Sol-low — filter parser and composite cursor)
- **Goal:** Encode IC-02 grammar, nullable domain shape, typed errors, and separate listing/group cursor contracts before repository work.
- **Files:**
  - Create `apps/server_core/internal/modules/listings/domain/read_model.go`.
  - Create `apps/server_core/internal/modules/listings/domain/read_model_test.go`.
  - Create or align the F-01 shared seam at `apps/server_core/internal/modules/listings/ports/read_repository.go`.
  - Create `apps/server_core/internal/modules/listings/ports/cursor.go`.
  - Create `apps/server_core/internal/modules/listings/ports/cursor_test.go`.
  - Create `apps/server_core/internal/modules/listings/transport/query.go`.
  - Create `apps/server_core/internal/modules/listings/transport/query_test.go`.
- **Failing test first:**
  - Accept exactly filter keys `status,sync_state,link_state,exception,has_exception,listing_type_code,product_id` and all fixed values. q is separate/trimmed; default limit 50; range 1..200.
  - Reject bogus/array/operator keys, invalid enum values, malformed boolean, duplicate scalar/filter values, and invalid limits before a port call. `invalid_filter` names the key in message and `error.details.key`.
  - Round-trip Unicode title + canonical ID through strict standard base64; reject empty/garbage/malformed payload, missing/extra fields, wrong cursor kind/version, and noncanonical IDs. Zero cursor means first page but is never emitted.
  - Prove nullable fields remain nil; D-5 labels are exact; malformed composite IDs take typed not-found path.
  - Unit-test below-margin tri-state: null cost/price/policy/source => nil. Add known true/false cases only after formula ruling.
  - Add pending-issue/group-severity unit matrices only after Risk 3 is ratified.
- **Implementation:**
  - Filter parsing stays in transport and returns typed `domain.ListingFilter`; repository receives no raw URL strings.
  - Repeated filter values are invalid because grammar is scalar, not arrays.
  - Versioned listing cursor contains last title + listing ID. Versioned group cursor contains null-last discriminator + product title + product ID tie-break.
  - Responses use pointer fields without `omitempty` where JSON null is contractual.
  - Typed application errors map only in transport. No generic cross-module filter/cursor framework.
- **Done:** targeted domain/ports/transport tests green with absolute GOCACHE; no open-risk semantic guessed. Advances **M01-C04/C06/C07**.

### Slice 2: Implement tenant-scoped `GET /listings` keyset behavior

- **Complexity:** complex (Sol-low — joined keyset SQL, filters/q, tri-state enrichment)
- **Goal:** Implement/test list repository, application, and handler behavior, but leave production GET registration to Slice 6.
- **Files:**
  - Edit `apps/server_core/internal/modules/listings/adapters/postgres/repository.go`.
  - Create `apps/server_core/internal/modules/listings/adapters/postgres/read_repository_integration_test.go`.
  - Create `apps/server_core/internal/modules/listings/application/read_service.go`.
  - Create `apps/server_core/internal/modules/listings/application/read_service_test.go`.
  - Create `apps/server_core/internal/modules/listings/adapters/internalread/cost_reader.go` if the ratified design uses the batch cost port.
  - Create `apps/server_core/internal/modules/listings/adapters/marketplaces/policy_reader.go` if marketplaces owns the ratified installation-policy boundary.
  - Edit `apps/server_core/internal/modules/listings/transport/http_handler.go`.
  - Edit `apps/server_core/internal/modules/listings/transport/http_handler_test.go`.
  - Edit `apps/server_core/internal/composition/root.go` only to inject read dependencies into the existing F-01 handler; do not register GET yet.
- **Failing test first:**
  - Real PostgreSQL page walk at limit 2 yields all seed rows once in title/listing-ID order and null terminal cursor; tenant B mirrors never leak.
  - Each filter key independently/in combination; absent link => unresolved; product ID only resolved; exception/has-exception predicates exact.
  - q matches case-insensitive title, provider listing ID, and joined snapshot seller SKU; percent/underscore are literal text, not SQL wildcards.
  - Cost/policy cases prove paired currencies, explicit cost/below-margin nulls, unavailable source not converted to empty facts, and known formula filtering after ruling.
  - Handler proves 400 nested installation/filter/cursor errors and field-for-field IC-02 response including required nulls, page size, UTC as-of, no provider-private data.
- **Implementation:**
  - Validate installation before reads; resolve ratified policy/source through published boundaries; batch cost without N+1; stamp one injected UTC serve time.
  - SQL selects `limit+1`, matches cursor tuple/order/collation, and derives next cursor from last returned row.
  - Direct listing and normalized link filters stay in SQL; q uses bound escaped `ILIKE`.
  - D-5 modality and D-8 price pair are exact. All nullable facts stay nullable.
  - Joined read projection supports link/q filtering before truncation; no item-by-item link reads.
  - No response cache, persisted link/cost/below-margin, or fallback on source failure.
- **Done:** repository/application/handler tests green; EXPLAIN uses keyset index at 2,000 rows; no offset/N+1. Advances **M01-C04/C06/C07** and list **C09**.

### Slice 3: Implement composite get with last ten timeline events

- **Complexity:** standard (Luna-high)
- **Goal:** Implement exact composite lookup and listing-scoped timeline, leaving route exposure to Slice 6.
- **Files:**
  - Edit `apps/server_core/internal/modules/listings/adapters/postgres/repository.go`.
  - Edit `apps/server_core/internal/modules/listings/adapters/postgres/read_repository_integration_test.go`.
  - Edit `apps/server_core/internal/modules/listings/application/read_service.go`.
  - Edit `apps/server_core/internal/modules/listings/application/read_service_test.go`.
  - Edit `apps/server_core/internal/modules/listings/transport/http_handler.go`.
  - Edit `apps/server_core/internal/modules/listings/transport/http_handler_test.go`.
- **Failing test first:**
  - Valid no-variation and real variation IDs pass; missing/extra/empty components, malformed decoded path, and unknown well-formed ID all return 404 `listing_not_found`, never 400.
  - Seed 12 tenant-A events with tied times plus tenant-B mirrors; return exactly newest 10 by at DESC + event-ID tie-break with no leak.
  - Get uses identical listing mapper as list. Successful body has top-level Listing fields plus non-null `timeline[]`; true empty history is `[]`.
- **Implementation:**
  - Parse composite server-side into F-01 `ListingKey`; literal `-` stays stored identity.
  - Getter binds tenant and all components, reuses projection/scanner, then reads timeline only after existence.
  - Cost/policy enrichment and failure semantics match list. Never substitute operation-run history.
- **Done:** malformed/absent 404, deterministic limit/order, exact shape, tenant isolation green. Advances **M01-C04/C06**.

### Slice 4: Implement by-product grouping with null group last

- **Complexity:** complex (Sol-low — group keyset SQL and worst-child reduction)
- **Goal:** Cursor over complete product groups in product-title order; retain every unlinked listing.
- **Files:**
  - Edit `apps/server_core/internal/modules/listings/adapters/postgres/repository.go`.
  - Edit `apps/server_core/internal/modules/listings/adapters/postgres/read_repository_integration_test.go`.
  - Edit `apps/server_core/internal/modules/listings/application/read_service.go`.
  - Edit `apps/server_core/internal/modules/listings/application/read_service_test.go`.
  - Edit `apps/server_core/internal/modules/listings/transport/http_handler.go`.
  - Edit `apps/server_core/internal/modules/listings/transport/http_handler_test.go`.
- **Failing test first:**
  - Seed two resolved products with equal titles/distinct IDs, multiple children, all non-resolved/absent link states, and tenant B mirrors.
  - Walk groups limit 1: each group once; product ID tie-break; complete children; `listing_count == len(listings)`.
  - Only resolved links form product groups. All others appear in one `{product_id:null,product_title:"sem produto"}` group, last globally and at most once.
  - Filters/q apply before grouping. Worst-child mapping covers every ratified state and is order-independent. Tenant B cannot affect state/count/order.
- **Implementation:**
  - SQL first computes normalized filtered children, selects `limit+1` group keys using explicit null-last discriminator/title/ID, then set-fetches all children for returned groups. Prefer one CTE statement; at most one bounded keys query + one set child query, never N+1.
  - Cursor encodes complete group key/kind/version; default NULL order is not relied on.
  - Reuse full list item mapping; compute deterministic group state only from ratified severity table.
- **Done:** group walk, equal-title stability, filters-before-grouping, complete children, null-last, cross-tenant green. Satisfies **M01-C05**, advances **C04**.

### Slice 5: Implement honest one-aggregate summary

- **Complexity:** complex (Sol-low — conditional aggregate and nullable margin semantics)
- **Goal:** Return exact summary from the ratified single-query design, excluding unknown margin rows.
- **Files:**
  - Edit `apps/server_core/internal/modules/listings/adapters/postgres/repository.go`.
  - Create `apps/server_core/internal/modules/listings/adapters/postgres/summary_query_test.go`, or place the real-PostgreSQL query assertions in `apps/server_core/internal/modules/listings/adapters/postgres/read_repository_integration_test.go` if the query is not factored.
  - Edit `apps/server_core/internal/modules/listings/application/read_service.go`.
  - Edit `apps/server_core/internal/modules/listings/application/read_service_test.go`.
  - Edit `apps/server_core/internal/modules/listings/transport/http_handler.go`.
  - Edit `apps/server_core/internal/modules/listings/transport/http_handler_test.go`.
- **Failing test first:**
  - Seed active/paused/closed, sync-error/stale, resolved/unlinked, below-margin true/false/null, and tenant B mirror. Assert exact tenant-A counters.
  - Null cost/price/policy yields listing null and contributes neither true nor false to below-margin; unavailable source follows ratified nullable counter behavior, not 0.
  - Missing installation returns 400 before dependencies; unknown installation follows 404 `installation_not_found`.
  - Query-count hook proves exactly one PostgreSQL aggregate; capture `EXPLAIN (ANALYZE, BUFFERS)` at 2,000 rows.
- **Implementation:**
  - One conditional-aggregate statement over tenant/installation joined projection; reuse exact list predicates for sync-error/stale/unlinked/below-margin.
  - Counters whose integrity is unknown are pointers; never blanket-`COALESCE` them to 0. Known empty total/active/paused may be numeric zero.
  - One injected UTC as-of.
  - Slice waits for Risks 1-2. If ruling permits Oracle bulk read plus one PostgreSQL aggregate, validation states that exact interpretation.
- **Done:** summary/null/error/tenant tests and query/EXPLAIN proof green. Satisfies **M01-C07** and summary **C09**, advances **C06**.

### Slice 6: Expose all GET routes contract-first with SDK in one commit

- **Complexity:** standard (Luna-high)
- **Goal:** Publicly register the tested paths only when handler, OpenAPI, SDK, and parity tests land together.
- **Files:**
  - Edit `apps/server_core/internal/modules/listings/transport/http_handler.go`.
  - Edit `apps/server_core/internal/modules/listings/transport/http_handler_test.go`.
  - Edit `apps/server_core/internal/composition/root.go`.
  - Edit `apps/server_core/internal/composition/root_test.go`.
  - Edit `contracts/api/marketplace-central.openapi.yaml` first.
  - Edit `packages/sdk-runtime/src/index.ts` second.
  - Edit `packages/sdk-runtime/src/index.test.ts`.
  - Edit `contracts/governance/modules.json` only if F-01 lacks the ratified F-02 dependencies.
- **Failing test first:**
  - Production root exposes exact four GETs and retains F-01 POST refresh; unsupported methods rejected.
  - OpenAPI parity asserts operation IDs `listListings/listListingsByProduct/getListing/getListingsSummary`, parameters, limit, all filters, nullable fields, response/error schemas/codes.
  - SDK URL tests cover dotted filters, q, composite path encoding, no undefined values, typed nulls/errors. Preserve `refreshListings`.
  - Source parity requires all four operation IDs/methods and matching required Listing fields.
- **Implementation:**
  - OpenAPI defines four paths and reusable Listing/link/money/error/pending/timeline/group/page/summary schemas. Nullable facts remain required properties with nullable values.
  - SDK adds string unions, nullable types, typed options with nested filter serialized to dotted keys, four typed GET methods, and shared typed errors.
  - Extend the one F-01 handler registrar; static routes are not swallowed by path variable. Composition constructs/injects once.
  - Slices 2-5 remain internal/unexposed. The first commit registering any GET path includes OpenAPI and SDK. Plan commit boundaries up front; do not rely on later squash/rewrite.
- **Done:** Go/root, SDK build/Vitest, parity and governance including `GOV_API_SDK_SPLIT` green; route-exposing `git show --stat` contains transport/composition + both contracts; no `apps/web` changes. Satisfies **M01-C08** and exposure for **C04-C07**.

### Slice 7: Prove IC-02 end-to-end and 2,000-row performance

- **Complexity:** complex (Sol-low — cross-surface integration and performance evidence)
- **Goal:** Exercise production-composed routes over real migrations/repositories and capture the complete ladder.
- **Files:**
  - Create `apps/server_core/tests/integration/listings_read_test.go`.
  - Create `apps/server_core/tests/integration/listings_read_perf_test.go` if separating the performance lane keeps the contract suite deterministic.
  - Edit `apps/server_core/internal/testsupport/postgres/fixtures.go` only if a reusable tenant-safe fixture has a second named consumer; otherwise keep fixtures local.
  - During execution, create `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-01-listings-read-spine/F-02-listings-read-api/validation.md` and named evidence files.
- **Failing test first:**
  - Use migration 0036, real listings repository/handler/routes/product-links tables, explicit deterministic stub only at external Oracle cost port, and real policy repository if ratified. Stub proves contract behavior, not live Oracle.
  - Seed exact `MLBTEST0001..0006` cases and F-02 product-links/snapshot fixtures per D-13, including resolved/unlinked and `-` to empty variation.
  - Full small-page walk returns six once, exact order/envelope/required nulls.
  - Every filter key; exception sync-error only seeded error; q finds title/provider ID/SKU.
  - By-product group walk proves counts/state/equal-title tie/null group last.
  - Get/timeline proves variation identities, last-ten order, malformed/absent 404.
  - Error matrix asserts status **and** `error.code`: 400 installation-required, 400 invalid-filter key/value, 400 invalid-cursor, 404 listing-not-found, and accepted F-01 concurrent 409 refresh-in-progress. Do not alter/re-plan refresh.
  - Null-cost JSON is exactly cost null + below-margin null; summary excludes it; all nullable facts remain null.
  - Tenant B identical identities but distinct data never affect tenant A list/q/filter/get/group/summary/cursor.
  - Performance seeds 2,000 rows and matching fixtures, starts direct dev-local server configuration, warms up, records 100 sequential list limit-50 calls, computes nearest-rank p95, fails at >=500ms.
  - Capture summary query count and EXPLAIN for 2,000 rows; assert ratified one-query rule.
- **Implementation:**
  - Bound SQL fixtures and exact tenant cleanup in ephemeral PostgreSQL only; no developer DB/live provider/OAuth/provider write/browser.
  - Exercise actual production constructors/registrar. External fakes are explicit and not live evidence.
  - Compare JSON keys/nulls, not decoded Go zero values. Walk returned cursors; never manufacture later cursors.
  - Performance uses monotonic timer, sequential requests, environment/SHA, and no concurrent refresh; retain samples/p95/query count/EXPLAIN/commands in validation evidence.
- **Done:**
  - Registered ephemeral-postgres lane runs all tests green.
  - L0: resolve repo `.gocache` absolute, set `$env:GOCACHE`, run `go build ./...` from `apps/server_core`; run `governance-validate`, `governance-drift`, and `governance` with zero findings.
  - L1: same GOCACHE for touched packages/integration lane; SDK workspace build and test green.
  - Evidence proves page walk, filters/q, grouping, errors, null honesty, tenant isolation, p95 <500ms, and ratified summary-query rule. Completes **M01-C04-C09**.

## Open Risks / Decisions

1. **Blocking — authoritative cost conflicts with literal C09 single-query wording.** Latest cost is Oracle while listings/link/policy are PostgreSQL and no local latest-cost projection exists. Ratify either one PostgreSQL aggregate after a separate authoritative bulk read, or an owned local projection/freshness contract. No ad hoc cache, profitability snapshot, zero, or production mock.

2. **Blocking — below-margin formula and policy selection are incomplete.** IC-02 says latest cost + min-margin and brief says cost + price, while current pricing logic also deducts commission/fixed fee/shipping. No installation-scoped policy getter or zero/multiple-policy rule exists. Ratify formula, percentage units, required inputs, and unique selection before known-result tests/SQL.

3. **Blocking for exact items/groups — pending-issue precedence and group severity are absent.** Contract does not say which issue wins when error/stale/unlinked/below-margin overlap, how attribute-required derives, localized messages, or which child maps to attention/error. Ratify one table reused by filters, has-exception, item mapping, grouping, and summary.

4. **Blocking for get timeline — no listing-scoped event source/vocabulary exists.** F-01 current row and installation-wide operation runs cannot supply last ten events. F-01 needs the shared source and ratified kinds/messages; permanent empty timeline is forbidden fallback.

5. **Contract-strength contradiction — mutable-title keyset cannot guarantee snapshot-stable walking through concurrent refresh with a single mutable row.** Atomic refresh prevents partial transactions, but a title update between pages can move across cursor. Encoding as-of + filtering updated-at also drops in-place-updated rows without versions. Clarify whether stability means deterministic keyset/no offset, or authorize snapshot/version/session storage.

6. **Non-blocking physical seam — no-variation sentinel differs.** Listings `-` vs product-links empty string is resolved by exact normalized full-identity joins and mandatory integration coverage.

7. **Non-blocking API detail — invalid-filter key location is unfixed.** Plan uses existing `error.details` as `details.key` and names key in message. Lock another location before Slice 6 if desired; status/code stay fixed.

8. **Do not reopen resolved decisions:** D-5 modality, D-8 paired nullable price, D-9 no cross-module FK, D-13 F-02 link fixtures, D-14 absolute GOCACHE, and D-15 contract-first SDK are incorporated. D-6/D-7/D-10/D-11 are F-01 refresh semantics and are not re-planned.
