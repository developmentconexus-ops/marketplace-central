
# F-01 listings-module-ingestion implementation plan

## Codebase Anchors

1. **Closest module analog — `catalog`**

   - `apps/server_core/internal/modules/catalog/domain/canonical_product.go:1` — domain package location and canonical entity ownership.
   - `apps/server_core/internal/modules/catalog/application/service.go:6-23` — application depends on domain/ports and receives concrete dependencies through constructors.
   - `apps/server_core/internal/modules/catalog/ports/repository.go:9-23` — module-owned ports describe named application consumers.
   - `apps/server_core/internal/modules/catalog/adapters/postgres/repository.go:15-24` — PostgreSQL adapter asserts its port implementation and receives pool plus tenant at composition time.
   - `apps/server_core/internal/modules/catalog/adapters/postgres/repository.go:32-36` and `apps/server_core/internal/modules/catalog/adapters/postgres/repository.go:85-89` — every read query includes `tenant_id = $1`.
   - `apps/server_core/internal/modules/catalog/transport/http_handler.go:20-24` — transport receives application dependencies rather than constructing infrastructure.
   - `apps/server_core/internal/modules/catalog/transport/http_handler.go:57-60` — standard nested `{error:{code,message,details}}` response.
   - `apps/server_core/internal/modules/catalog/transport/http_handler.go:63-78` — handler-owned route registration on `httpx.RouteRegistrar`, including method-aware patterns.
   - `apps/server_core/internal/composition/root.go:13-17` and `apps/server_core/internal/composition/root.go:353-360` — composition root imports each layer, constructs repository/service/handler, and calls `Register`.
   - Mirror this structure with `listings/domain`, `listings/application`, `listings/ports`, `listings/adapters/connectors`, `listings/adapters/integrations`, `listings/adapters/postgres`, and `listings/transport`. Do not create unused `events` or `readmodel` packages in F-01.

2. **Published connectors `ListListings` capability**

   - `apps/server_core/internal/modules/connectors/ports/marketplace_capability.go:24-27` — published signature is `ListListings(context.Context, domain.ListListingsInput) ([]domain.ListingSnapshot, error)` plus `ReadListing`; F-01 must consume this port, never the Mercado Livre adapter’s private HTTP types.
   - `apps/server_core/internal/modules/connectors/domain/capability.go:85-90` — account reference carries tenant, installation, provider, and provider account identity.
   - `apps/server_core/internal/modules/connectors/domain/capability.go:103-108` — pagination input is `{AccountRef, Cursor string, Status string, Limit int}`.
   - `apps/server_core/internal/modules/connectors/domain/capability.go:118-130` — current returned listing snapshot exposes provider/item/variation/status/title/nullable quantity/timestamps plus variations, but not price, currency, or listing modality.
   - `apps/server_core/internal/modules/connectors/domain/capability.go:143-148` — variation snapshot carries variation identity and nullable quantity.
   - `apps/server_core/internal/modules/connectors/domain/capability.go:38-45` and `apps/server_core/internal/modules/connectors/domain/capability.go:47-82` — typed capability errors currently cover rate limit, validation, unsupported shape, transient, invalid payload, and invalid reference; `ErrorCodeOf` extracts the translated code.
   - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:107-141` — ML implementation treats `Cursor` as an offset, requests one provider page, reads each returned item, and returns only the slice; there is no provider total or next-cursor result.
   - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:117-125` — page size defaults to 50 and cursor becomes `offset`.
   - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:129-141` — any item-detail failure aborts that page; partial page results are not returned.
   - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:538-566` — private ML payload becomes the published snapshot and variations are expanded without leaking the private response struct.
   - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:742-752` — private item response currently omits `price`, `currency_id`, and `listing_type_id`; extend this adapter and the published canonical snapshot rather than reading ML private types from `listings`.
   - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:472-495` — provider HTTP errors are translated before leaving the adapter.
   - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:449-460` — credential-resolution errors are currently collapsed into transient/invalid-reference errors.
   - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:477-487` — 401/403 currently fall into provider validation; this cannot honestly record an auth-class refresh failure and requires a contract-level connector error classification in F-01.

3. **`integration_operation_runs` lifecycle**

   - `apps/server_core/internal/modules/integrations/domain/lifecycle.go:51-55` — statuses are fixed as `queued`, `running`, `succeeded`, `failed`, and `cancelled`.
   - `apps/server_core/internal/modules/integrations/domain/operation_run.go:5-22` — run fields include operation type, result/failure/translated error codes, actor, provider evidence, duration, and lifecycle timestamps.
   - `apps/server_core/internal/modules/integrations/ports/operation_run_store.go:9-12` — current published store can save/upsert a run and list installation history.
   - `apps/server_core/internal/modules/integrations/application/operation_service.go:13-28` — `RecordOperationInput` is the application boundary for all run states and evidence.
   - `apps/server_core/internal/modules/integrations/application/operation_service.go:39-65` — recording builds a tenant-scoped run and saves it through the store.
   - `apps/server_core/internal/modules/integrations/application/fee_sync_service.go:119-190` — closest async pattern: validate an existing installation, record `queued`, dispatch background execution, and return the run ID.
   - `apps/server_core/internal/modules/integrations/application/fee_sync_service.go:193-266` — execution records `running`, performs work, then records terminal success/failure.
   - `apps/server_core/internal/modules/integrations/application/fee_sync_service.go:269-306` — current concurrency check scans run history for queued/running states, but is not an atomic database guard; F-01 needs a transactionally serialized begin operation for M01-C02.
   - `apps/server_core/internal/modules/integrations/application/provider_operation_service.go:233-261` — connector failures become failed operation runs and `connectorsdomain.ErrorCodeOf` is copied to `translated_error_code`.
   - `apps/server_core/internal/modules/integrations/adapters/postgres/operation_run_repo.go:24-55` — one run ID is upserted through queued/running/terminal states.
   - `apps/server_core/internal/modules/integrations/adapters/postgres/operation_run_repo.go:59-70` — installation history query is tenant- and installation-scoped.
   - `apps/server_core/migrations/0016_integrations_foundation.sql:196-224` — table definition, status check, installation foreign key, and tenant-scoped indexes.
   - `apps/server_core/migrations/0021_integration_operation_run_evidence.sql:1-4` — translated error, provider evidence, and duration extensions.
   - Use operation type `listings_refresh` consistently for the exclusive guard and lifecycle. Do not call `StartAuthorize`, `StartReauth`, or any provider-write capability.

4. **Migration layout and allocation**

   - `ARCHITECTURE.md:251-252` — migrations are forward-only sequential `NNNN_description.sql` files.
   - `apps/server_core/migrations/source.go:10-16` — all SQL files are embedded and identified by full filename.
   - `apps/server_core/migrations/0033_orders_sankhya_linkage.sql:1`, `apps/server_core/migrations/0034_catalog_legacy_product_reference_evidence.sql:1`, and `apps/server_core/migrations/0035_catalog_codprod_compatibility.sql:1` — the current checkout already contains committed migrations using every number allocated to M-01.
   - `contracts/governance/invariants.json:14` — `GOV_MIGRATION_PREFIX` checks migration-prefix collisions.
   - `contracts/governance/invariants.json:26-33` — only the historical 0021 duplicate has an exact exception; a new duplicate 0033 would fail governance.
   - **Blocking contradiction:** the task requires the next listings migration to be `0033` and forbids exceeding `0035`, while this checkout already owns 0033–0035 for unrelated committed migrations. Do not overwrite, rename, or duplicate them. The milestone/hub must reconcile the accepted base or allocate a valid number before Slice 1 implementation.

5. **Governance module and layering rules**

   - `contracts/governance/modules.json:5-15` — every module declares root, code-owner path, composition requirement, OpenAPI prefixes, and explicit module dependencies.
   - Add `listings` with root/code owner `apps/server_core/internal/modules/listings`, `composition_required: true`, OpenAPI prefix `["/listings"]`, and dependencies `["connectors","integrations"]`; do not declare `product_links` in F-01 because it is not consumed until F-02.
   - `contracts/governance/invariants.json:5-12` — module coverage, declared dependency, target-layer, composition import, application import, PostgreSQL driver, and API/SDK atomicity checks apply.
   - `scripts/harness/Policy.psm1:313-325` — cross-module imports must be declared; imports of another module’s adapters/transport/registry are rejected.
   - `scripts/harness/Policy.psm1:328-334` — every composition-required module must appear in the composition root.
   - `contracts/governance/invariants.json:12` — changing OpenAPI without `packages/sdk-runtime`, or vice versa, fails `GOV_API_SDK_SPLIT`.
   - `contracts/governance/shared-seams.json:4-7` — API/SDK, migration sequence, and composition root are exclusive seams.
   - `contracts/governance/shared-seams.json:12-20` — the connectors capability contract is also an exclusive seam; Slice 2 must have explicit ownership before editing it.
   - `contracts/governance/execution-lanes.json:15-28` — integration proof uses localhost ephemeral PostgreSQL with no external side effects.
   - `scripts/harness/Postgres.psm1:119` — `./tests/integration` already runs in the registered integration lane, so the F-01 end-to-end test belongs there and does not require a harness edit.

6. **Tenant scoping invariant**

   - `ARCHITECTURE.md:37` — every business table carries `tenant_id`.
   - `apps/server_core/internal/modules/catalog/adapters/postgres/repository.go:17-24` — repository instances receive the active tenant at construction.
   - `apps/server_core/internal/modules/catalog/adapters/postgres/repository.go:32-36` and `apps/server_core/internal/modules/catalog/adapters/postgres/repository.go:85-89` — reads bind tenant as the first predicate.
   - `apps/server_core/internal/modules/product_links/adapters/postgres/link_candidate_repo.go:101-111`, `apps/server_core/internal/modules/product_links/adapters/postgres/link_candidate_repo.go:157-165`, and `apps/server_core/internal/modules/product_links/adapters/postgres/link_candidate_repo.go:204-214` — list/get/join identities all carry tenant and installation predicates.
   - `apps/server_core/internal/modules/integrations/adapters/postgres/operation_run_repo.go:66-70` — operation-run reads also use tenant plus installation.
   - Every F-01 `SELECT`, `INSERT ... ON CONFLICT`, close update, active-run lookup, and cleanup assertion must bind `tenant_id`; integration tests must prove tenant B survives a tenant A refresh.

7. **`product_links` ownership**

   - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/listings-read-interface-contract.md:134-136` — listing persistence excludes link, cost, and below-margin fields; link truth is joined from `product_links` only at read time.
   - `apps/server_core/migrations/0022_product_links_listing_snapshots.sql:1-16` — the old listing snapshot table is a distinct product-link-owned projection.
   - `apps/server_core/internal/modules/product_links/adapters/postgres/link_candidate_repo.go:204-214` — canonical link truth is read by tenant, installation, provider item, and variation.
   - `apps/server_core/internal/modules/product_links/adapters/postgres/listing_snapshot_repo.go:31-67` — the existing snapshot writer is owned by `product_links` and must not be reused or modified.
   - F-01 must not edit any `product_links` file, create link columns, seed link rows, or import the module. The IC-02 “resolved/unlinked” seed characteristics become join assertions in F-02; F-01 seeds only the listings-owned subset under the same fixed IDs.

8. **OpenAPI and SDK lockstep**

   - `contracts/api/marketplace-central.openapi.yaml:506-528` — nearby async fee-sync operation shows a `202` response and shared `ErrorResponse` references.
   - `contracts/api/marketplace-central.openapi.yaml:2504-2517` — accepted-operation schema pattern.
   - `contracts/api/marketplace-central.openapi.yaml:4292-4321` — shared error envelope requires `code` and `message` and permits arbitrary `details`.
   - `packages/sdk-runtime/src/index.ts:981-992` — matching SDK error and client-error types.
   - `packages/sdk-runtime/src/index.ts:1041-1052` — shared JSON POST helper preserves typed errors.
   - `packages/sdk-runtime/src/index.ts:1085-1086` — nearby async method pattern.
   - `packages/sdk-runtime/src/index.test.ts:170-195` — SDK request tests capture URL/method and assert the typed response.
   - `packages/sdk-runtime/package.json:7-10` — SDK verification is TypeScript build plus Vitest.
   - There is no OpenAPI code-generation command in the repository. Follow the established contract-first pattern: change OpenAPI first, then hand-maintained SDK types/method/tests in the same slice and same commit; do not claim generated evidence.

9. **Architecture/contract reconciliation**

   - `ARCHITECTURE.md:39` and `ARCHITECTURE.md:243` require all provider access through connector ports and keep connector-owned code free of business state; the planned connector adapter plus listings-owned persistence satisfies this.
   - `ARCHITECTURE.md:245` forbids synchronous web requests from depending on connector availability; returning `202` after durable queueing and running the pull asynchronously satisfies this.
   - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:111` fixes ADR-12: read-only listings module, connector-only ingestion, composite identity, and manual refresh.
   - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:116` fixes ADR-17: nullable unknown facts from the first migration.
   - **Documentation contradiction:** `ARCHITECTURE.md:7` says accepted decisions are indexed under `docs/architecture/decisions`, but `docs/architecture/decisions/README.md:5-10` lists only ADR-004 through ADR-007. ADR-12 and ADR-17 exist only as decided rows in the MIS-003 mission. Their required behavior is nevertheless unambiguous and is followed here; the missing formal ADR records should be repaired by the architecture owner, not by F-01.
   - Apart from the migration-number collision and missing formal ADR records, no contradiction was found between IC-02, the F-01 brief, M01-C01/C02/C03, and the module architecture.

## Slices

### Slice 1: Establish the canonical listings schema, domain vocabulary, and governed module registration

- **Complexity:** standard

- **Files:**

  - Create `apps/server_core/migrations/0033_listings.sql` only after the migration-allocation blocker is resolved; this is the filename required by the supplied allocation.
  - Create `apps/server_core/migrations/listings_test.go`.
  - Create `apps/server_core/internal/modules/listings/domain/listing.go`.
  - Create `apps/server_core/internal/modules/listings/domain/listing_test.go`.
  - Edit `contracts/governance/modules.json`.

- **Failing test first:**

  - Add `TestListingsMigrationMatchesIC02` in `apps/server_core/migrations/listings_test.go`, RED before the migration exists. Read the embedded migration by exact filename and assert:
    - the table is exactly `listings`;
    - the exact columns are present;
    - `variation_id` is text and uses `'-'` as the no-variation value at the application boundary;
    - the primary key is `(tenant_id, installation_id, provider_listing_id, variation_id)`;
    - `listing_type_code`, `price_amount`, `price_currency`, `published_quantity`, `sync_error`, `quality_score`, `sales_30d`, and `fetched_at` permit SQL NULL under ADR-17;
    - status allows only `active|paused|closed|unknown`;
    - sync state allows only `synced|error|stale|queued|syncing|paused_sync`;
    - quality score is constrained to 0–100 when non-null;
    - no link, product, seller SKU, cost, below-margin, or pending-issue columns exist.
  - Add table-driven domain tests in `apps/server_core/internal/modules/listings/domain/listing_test.go` for the fixed status/sync-state values and `ListingKey` normalization:
    - no provider variation becomes literal `-`;
    - two variations under the same provider item remain distinct keys;
    - empty tenant, installation, provider listing ID, or title is rejected rather than defaulted;
    - nullable facts remain nil, including an explicit cross-check that nil quantity is not converted to zero.

- **Implementation:**

  - Define only the F-01-owned canonical persistence entity, value types, fixed enums, composite `ListingKey`, and validation errors. Do not add F-02 filters, cursor types, link DTOs, summary DTOs, or GET transport.
  - Migration creates the IC-02 table with:
    - non-null tenant/installation/provider/provider-listing/variation/title/status/sync-state/created/updated identity and lifecycle columns;
    - nullable listing modality and all provider facts that may be unknown;
    - exact PK and checks from IC-02;
    - UTC-capable `TIMESTAMPTZ` timestamps;
    - an FK to `(tenant_id, installation_id)` only if the milestone resolves the FK/delete-policy question below.
  - Do not use zero/default values for nullable facts. Defaults may be used only for database-generated `created_at`/`updated_at`; ingestion still supplies one explicit capture time.
  - Add one F-02-ready index on `(tenant_id, installation_id, title, provider_listing_id, variation_id)` only if accepted as part of M01-C09 preparation; do not add speculative filter-specific indexes in F-01.
  - Register `listings` in `modules.json` with declared dependencies on `connectors` and `integrations`, composition required, and `/listings` prefix.
  - Do not create unused interfaces. Ports appear in the slices where their application consumers are introduced.

- **Done:**

  - Migration/domain tests are green.
  - `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command governance` is green with no migration-prefix exception added for F-01.
  - From `apps/server_core`, resolve `.gocache` to an absolute path and run `go test ./migrations ./internal/modules/listings/domain -count=1`.
  - From `apps/server_core`, run `go build ./...` with the same absolute `GOCACHE`.
  - Schema review confirms no `product_links` mutation or forbidden columns.
  - Advances **M01-C01** by establishing the exact upsert key and closed status, and **M01-C03** by making unknown facts nullable from the first migration.

### Slice 2: Extend the published listing capability and map Mercado Livre facts into canonical nullable listings

- **Complexity:** standard

- **Files:**

  - Edit `apps/server_core/internal/modules/connectors/domain/capability.go`.
  - Edit `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`.
  - Edit `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter_test.go`.
  - Create `apps/server_core/internal/modules/listings/adapters/connectors/mapper.go`.
  - Create `apps/server_core/internal/modules/listings/adapters/connectors/mapper_test.go`.

- **Failing test first:**

  - Extend `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter_test.go` with an ML item fixture containing exact decimal price, `currency_id`, `listing_type_id`, item status, item quantity, and two variations. Assert the published `connectors/domain.ListingSnapshot` carries normalized facts and contains no raw JSON object.
  - Add `TestMapListingSnapshotToCanonicalRows` in `apps/server_core/internal/modules/listings/adapters/connectors/mapper_test.go` and make it RED for:
    - item without variations → one `variation_id="-"` row;
    - item with variations → one row per variation and no synthetic item row;
    - item facts such as title, price, modality, status, and fetched time propagate to each variation while variation quantity remains variation-specific;
    - `active`, `paused`, and `closed` map exactly;
    - an unrecognized provider status maps to canonical `unknown`;
    - an unrecognized or empty modality maps to nil `listing_type_code`;
    - nil price, quantity, quality, and sales stay nil rather than zero;
    - decimal price remains a base-10 string suitable for PostgreSQL `NUMERIC`, without a float round trip;
    - missing required identity/title/fetched facts fails the page mapping honestly rather than dropping or fabricating a row.
  - Add a negative adapter test proving HTTP 401/403 and credential-unavailable failures receive the ratified provider-auth connector error classification rather than provider validation/transient, without calling authorization APIs.

- **Implementation:**

  - Extend the published connector listing snapshot with only the provider-agnostic facts required by IC-02: decimal price amount, price currency, and provider listing modality code. Do not add ML response types or JSON blobs to the published contract.
  - Decode ML price without `float64` so `89.00` is not rounded or rewritten through binary floating point. Keep the raw response struct and JSON decoding inside `adapters/mercado_livre`.
  - Extend the private ML item response with `price`, `currency_id`, and `listing_type_id`; map those into the published connector snapshot.
  - Keep `RawProviderRef` as the existing sanitized path reference only; the listings mapper ignores it, and no raw response payload is persisted or returned.
  - Implement a provider-aware mapper under `listings/adapters/connectors`:
    - validates required published facts;
    - applies the fixed status mapping;
    - applies the milestone-ratified modality allowlist;
    - expands variations using `-` only for items with no variations;
    - sets successful ingestion `sync_state=synced`, `sync_error=nil`, `quality_score=nil`, and `sales_30d=nil`;
    - preserves nullable price/quantity facts.
  - Add a typed connector auth error classification once its exact code is ratified. `ErrorCodeOf` must carry it into operation-run `translated_error_code`; do not trigger reauthorization.
  - Update all existing connector snapshot literals/tests affected by the additive fields without changing old behavior.
  - Because `capability.go` is an exclusive provider-capability seam, this slice starts only after the milestone confirms sole ownership.

- **Done:**

  - Mapping tests cover known, unmappable, null, variation, malformed-identity, and auth-class cases.
  - Existing connector adapter and capability-service tests remain green.
  - From `apps/server_core`, run `go test ./internal/modules/connectors/... ./internal/modules/listings/adapters/connectors -count=1` with absolute `GOCACHE`.
  - Run `go build ./...` and the governance lane.
  - No listing mapper imports `connectors/adapters/mercado_livre`; only the published connector domain/port surface is consumed.
  - Directly satisfies **M01-C03** and prepares provider-error evidence for **M01-C01**.

### Slice 3: Implement full-page ingestion with one atomic tenant-scoped upsert-and-close transaction

- **Complexity:** complex

- **Files:**

  - Create `apps/server_core/internal/modules/listings/ports/ingestion.go`.
  - Create `apps/server_core/internal/modules/listings/application/ingestion.go`.
  - Create `apps/server_core/internal/modules/listings/application/ingestion_test.go`.
  - Create `apps/server_core/internal/modules/listings/adapters/connectors/source.go`.
  - Create `apps/server_core/internal/modules/listings/adapters/connectors/source_test.go`.
  - Create `apps/server_core/internal/modules/listings/adapters/postgres/repository.go`.
  - Create `apps/server_core/internal/modules/listings/adapters/postgres/repository_integration_test.go`.

- **Failing test first:**

  - Add `TestIngestionWalksEveryConnectorPageBeforePersistence` in `application/ingestion_test.go`. A source returns two full pages then one short page; assert cursors advance by the number of provider items returned, all variation-expanded keys reach one completed-store call, and no persistence occurs before pagination finishes.
  - Add `TestIngestionCapabilityErrorMidPullDoesNotApplySnapshot`. Page 1 succeeds and page 2 returns a typed connector error; assert the store receives zero calls and the original error remains classifiable through `connectorsdomain.ErrorCodeOf`.
  - Add `TestIngestionRejectsDuplicateCanonicalKey`. Feed duplicate item/variation identities across pages and assert the refresh fails before persistence instead of allowing last-write-wins ambiguity.
  - Add real PostgreSQL tests under the integration build tag in `repository_integration_test.go`:
    - seed two rows for tenant A/installation A and an identical provider identity for tenant B;
    - apply a completed pull that updates one row, adds a distinct variation, and omits one old tenant-A row;
    - assert the PK keeps variations distinct, omitted tenant-A row becomes `closed`, returned rows retain mapped statuses/facts, existing `created_at` is preserved, returned rows get refreshed `updated_at/fetched_at`, and tenant B is byte-for-byte unchanged;
    - force an insert/check failure within the transaction and assert all earlier close/upsert effects roll back;
    - apply a completed empty pull and assert all rows for only the target tenant/installation become closed.

- **Implementation:**

  - Define two named ports consumed by `application.Ingestion`:
    - a page source that accepts the already-resolved installation account, cursor, and limit and returns published connector snapshots;
    - a completed-pull store that atomically applies canonical rows for one tenant/installation.
  - The connector source selects the published `ListingReader` through `MarketplaceCapabilityService`, constructs `ProviderAccountRef` from the already-connected installation, calls `ListListings`, and invokes the Slice 2 mapper. It does not use `ProviderOperationService.ListListings`, because that method has no cursor and would create a separate operation run per page.
  - The application loops from cursor `""`, uses a bounded configured page size, advances the numeric offset by provider items returned, and completes only after a short page. It validates duplicate canonical keys across the whole pull.
  - A capability or mapping error at any page returns before the repository is called, guaranteeing the provider-unreachable row-count invariant.
  - The PostgreSQL adapter opens one transaction and:
    - marks all current rows for the exact tenant/installation `status='closed'`;
    - upserts every returned row on `(tenant_id, installation_id, provider_listing_id, variation_id)`;
    - restores each returned row’s mapped status and facts;
    - updates `updated_at` on every upsert;
    - sets `fetched_at` only from a successful provider snapshot;
    - preserves original `created_at` on conflict;
    - commits only after every row succeeds.
  - Closing first and re-upserting inside the same transaction avoids temporary externally visible closure and avoids constructing unsafe dynamic `NOT IN` SQL. Empty completed pulls close all rows; failed or incomplete pulls never enter this transaction.
  - Every statement includes both tenant and installation predicates. Repository constructors receive tenant ID from composition, matching the catalog/integrations idiom.
  - Do not add read-query ports or GET endpoint code.

- **Done:**

  - Application/source tests prove full pagination, distinct variation keys, duplicate rejection, and zero store calls on a mid-pull failure.
  - PostgreSQL integration tests prove real upsert/close/rollback/cross-tenant behavior rather than asserting SQL strings or mocks.
  - Run targeted unit tests with absolute `GOCACHE`.
  - Run the registered ephemeral-PostgreSQL integration lane and confirm the tagged repository test executes.
  - Run `go build ./...`, touched-package `go test ./...`, and governance.
  - Directly advances **M01-C01** and the unchanged-row failure half of **M01-C03**.

### Slice 4: Add the asynchronous refresh endpoint, atomic concurrent-run guard, operation lifecycle, composition wiring, and OpenAPI/SDK contract in one commit

- **Complexity:** complex

- **Files:**

  - Create `apps/server_core/internal/modules/listings/ports/integrations.go`.
  - Create `apps/server_core/internal/modules/listings/application/refresh_service.go`.
  - Create `apps/server_core/internal/modules/listings/application/refresh_service_test.go`.
  - Create `apps/server_core/internal/modules/listings/adapters/integrations/gateway.go`.
  - Create `apps/server_core/internal/modules/listings/adapters/integrations/gateway_test.go`.
  - Create `apps/server_core/internal/modules/listings/transport/http_handler.go`.
  - Create `apps/server_core/internal/modules/listings/transport/http_handler_test.go`.
  - Edit `apps/server_core/internal/modules/integrations/ports/operation_run_store.go`.
  - Edit `apps/server_core/internal/modules/integrations/application/operation_service.go`.
  - Edit `apps/server_core/internal/modules/integrations/application/operation_service_test.go`.
  - Edit `apps/server_core/internal/modules/integrations/application/credential_service_test.go`.
  - Edit `apps/server_core/internal/modules/integrations/application/provider_operation_service_test.go`.
  - Edit `apps/server_core/internal/modules/integrations/adapters/postgres/operation_run_repo.go`.
  - Create `apps/server_core/internal/modules/integrations/adapters/postgres/operation_run_repo_integration_test.go` if the exclusive-begin proof is not kept entirely in the final endpoint integration test.
  - Edit `apps/server_core/internal/composition/root.go`.
  - Edit `apps/server_core/internal/composition/root_test.go`.
  - Edit `contracts/api/marketplace-central.openapi.yaml`.
  - Edit `packages/sdk-runtime/src/index.ts`.
  - Edit `packages/sdk-runtime/src/index.test.ts`.

- **Failing test first:**

  - In `refresh_service_test.go`, add:
    - missing installation ID returns typed `installation_required` before any integration lookup;
    - unknown installation returns typed `installation_not_found`;
    - a connected installation records `queued`, returns its operation ID, then records `running` and `succeeded` around ingestion;
    - an ingestion connector error records `failed`, preserves the connector translated error code, contains no raw provider body, and does not call `StartAuthorize`/`StartReauth` because those dependencies do not exist;
    - the async runner is injected so tests deterministically hold and release execution;
    - if terminal run persistence fails, the error is surfaced to logs/test hooks and is not represented as success.
  - In `operation_service_test.go` and the PostgreSQL operation repository integration test, issue two simultaneous exclusive starts for the same tenant/installation/`listings_refresh`; assert exactly one queued run is created and the loser receives the active run object. Also assert different installations and different tenants do not block each other.
  - In `http_handler_test.go`, add exact transport cases:
    - `{}` or blank `installation_id` → 400 `installation_required`;
    - unknown installation → 404 `installation_not_found`;
    - accepted connected installation → 202 exactly `{"operation_run_id":"..."}`;
    - active refresh → 409 `refresh_in_progress` with the same active operation ID in the ratified error location;
    - malformed JSON and trailing JSON are rejected as 400 rather than partially accepted;
    - GET on `/listings/refresh` is not accepted.
  - In `root_test.go`, assert `POST /listings/refresh` is registered by the production composition root and is not a 404.
  - In `packages/sdk-runtime/src/index.test.ts`, assert `refreshListings({installation_id:"inst_test"})` sends `POST /listings/refresh`, serializes exactly that body, returns the operation ID, and preserves a 409 typed error including the active run ID.
  - Add a contract parity test/assertion that OpenAPI contains operation ID `refreshListings`, request schema, exact `202/400/404/409` responses, and the same SDK method/type names.

- **Implementation:**

  - Add an integrations-owned exclusive operation-start primitive rather than letting `listings` query the integrations table:
    - serialize by tenant, installation, and operation type inside one PostgreSQL transaction using a database advisory transaction lock or an equivalently atomic PostgreSQL primitive;
    - query queued/running `listings_refresh` runs with tenant and installation predicates;
    - return the active run when found;
    - otherwise insert the new queued run before releasing the lock.
  - Extend the published integrations store/application boundary and update its existing test stubs. Do not import `integrations/adapters/postgres` from `listings`.
  - `listings/adapters/integrations` adapts the published installation and operation application services into narrow listings-owned ports. It exposes only the fields required to resolve the already-connected installation and record lifecycle changes.
  - Validate installation existence and the contract-approved runnable status before starting a run. Never call authorization or mutate installation connection state.
  - `RefreshService.Start`:
    - trims and validates installation ID;
    - resolves the installation;
    - creates a cryptographically collision-resistant operation run ID;
    - atomically starts `listings_refresh` as queued or returns the active-run conflict;
    - dispatches a background function and immediately returns the queued run ID.
  - Background execution:
    - uses a fresh background context, not the cancelled request context;
    - records `running` with start time;
    - invokes Slice 3 ingestion;
    - records `succeeded` with completed time and sanitized `page_count`/`listing_count` evidence, or `failed` with safe fixed failure code and connector `translated_error_code`;
    - never includes raw provider JSON or credentials in failure/evidence fields.
  - Inject the async runner and clock for tests; composition supplies the production goroutine runner. Do not add blanket panic recovery or a new scheduler.
  - Transport performs strict JSON decoding, maps only the IC-02 error matrix, and uses the standard nested error envelope. The active run ID goes in the contract-ratified location; if the milestone adopts the existing generic schema convention, use `error.details.operation_run_id`.
  - Register only `POST /listings/refresh`; no GET listing routes are introduced.
  - Composition creates:
    - tenant-scoped listings repository;
    - connector page source backed by the already-built marketplace capability service;
    - integrations gateway backed by installation and operation services;
    - ingestion and refresh services;
    - listings handler registration.
  - Update OpenAPI first with:
    - `POST /listings/refresh`;
    - request `{installation_id}` required;
    - `202` response containing only required `operation_run_id`;
    - 400/404/409 shared error responses and exact documented codes;
    - no F-02 GET paths or schemas.
  - In the same slice and commit, add SDK request/response types and `refreshListings` client method plus tests. The repository has no generator, so record this as contract-first hand-maintained parity, not generated output.
  - Do not split endpoint and OpenAPI/SDK work into separate commits; doing so would fail M01-C08/GOV_API_SDK_SPLIT.

- **Done:**

  - Unit and PostgreSQL concurrency tests prove exactly one active run under real simultaneous starts.
  - Handler tests prove 202/400/404/409 status and error-code parity, including the active run ID.
  - Failed ingestion produces a terminal failed run with honest classification and no row mutation.
  - Structural dependency review confirms no authorization method is reachable from refresh.
  - `go build ./...`, touched Go tests, SDK `npm run build --workspace @marketplace-central/sdk-runtime`, SDK tests, and governance all pass.
  - `git show --stat` for this slice’s one green commit contains the handler/composition changes, OpenAPI, and `packages/sdk-runtime`.
  - Directly satisfies **M01-C02**, completes operation-run evidence for **M01-C01**, and advances **M01-C08**.

### Slice 5: Prove the complete refresh contract against ephemeral PostgreSQL and a stubbed published connector capability

- **Complexity:** complex

- **Files:**

  - Create `apps/server_core/tests/integration/listings_refresh_test.go`.
  - Edit `apps/server_core/internal/testsupport/postgres/fixtures.go` only if a reusable connected-installation fixture is required; otherwise keep all F-01 fixture setup local to the new integration test.

- **Failing test first:**

  - Add an integration-build-tagged suite using `testsupport/postgres.OpenPool`, real migrations, real listings and operation-run repositories, real application/transport services, and a stub implementation of the published `connectors/ports.ListingReader`.
  - `TestListingsRefreshSeedsIC02RowsAndClosesMissing`:
    - seed connected `inst_test` for tenant A;
    - stub paginated connector results producing exactly `MLBTEST0001..MLBTEST0006`, including item/variation identities, active/paused/closed/unknown-capable statuses, synced/error/stale/paused-sync states as applicable to listings-owned data, nullable facts, and at least one variation;
    - POST refresh and assert 202 with an operation ID;
    - wait deterministically for the run to reach `succeeded`;
    - SQL assert exactly six tenant-A rows with the composite PK, `variation_id='-'` where appropriate, and required nullable fields such as `sales_30d IS NULL`;
    - remove one provider listing from the completed stub pull, refresh again, and assert only that row becomes closed while the other five keep their mapped facts/statuses;
    - do not seed or assert `product_links`; resolved/unlinked join state belongs to F-02.
  - `TestListingsRefreshRejectsConcurrentRunWithActiveID`:
    - block the first source call after its queued/running run exists;
    - first POST returns 202;
    - second POST for the same installation returns 409 `refresh_in_progress` carrying the first operation ID;
    - SQL assert only one queued/running run exists;
    - release the source and assert the same run succeeds.
  - `TestListingsRefreshCapabilityErrorMidPullLeavesRowsUnchanged`:
    - preseed known listing rows;
    - use a small test page size so page 1 succeeds and page 2 returns a typed capability error;
    - POST returns 202;
    - wait for `failed`;
    - SQL compare pre/post row count and persisted values, not only the count;
    - assert `translated_error_code` is the ratified provider-auth code for an auth error, or the exact connector transient code for an unreachable provider;
    - assert provider evidence contains no raw response body.
  - `TestListingsRefreshIsTenantIsolated`:
    - seed the same installation/provider listing identity under tenant B;
    - refresh/close tenant A;
    - assert tenant B listings and operation runs are unaffected.
  - `TestListingsRefreshUnknownInstallation`:
    - POST a nonexistent installation and assert 404 with no operation run and no listing changes.

- **Implementation:**

  - Construct the stub through `connectors/application.MarketplaceCapabilityService` and a `ProviderCapabilitySet` so the test exercises the published connector selection boundary, not a listings-only fake that bypasses it.
  - Use a controllable page stub supporting cursor/limit, deliberate blocking, and a configured failure page.
  - Build the real handler on an in-memory `http.ServeMux` or `httpx.RouteClassMux`, POST through `httptest`, and poll the real operation-run table with a bounded deadline.
  - Seed an already-connected installation with external account ID and no authorization call path. Do not invoke live ML, OAuth, `StartAuthorize`, product writes, or `product_links`.
  - Query all assertions with tenant and installation predicates.
  - Clean up both tenants’ listings and operation runs through exact predicates; rely on the ephemeral database lifecycle rather than touching any developer database.
  - Capture the integration transcript and SQL assertions for `F-01-listings-module-ingestion/validation.md` during implementation, without adding evidence-writing behavior to production code.

- **Done:**

  - Registered integration lane runs the new tagged test and reports ephemeral PostgreSQL/migrations green.
  - First pull yields six rows; second completed pull closes exactly the removed row.
  - Concurrent requests yield 202 then 409 with the same run ID and only one active run.
  - Mid-pull failure leaves all prior rows unchanged and records a terminal failed operation honestly.
  - Unmappable status/modality and nil facts are stored as `unknown`/SQL NULL.
  - Cross-tenant assertions pass.
  - Final L0/L1 ladder:
    - from `apps/server_core`, absolute `GOCACHE` pointing at `.gocache`, then `go build ./...`;
    - governance validation and drift lanes;
    - `go test ./... -count=1` for non-tagged touched packages;
    - registered ephemeral-PostgreSQL integration lane for tagged tests;
    - SDK build and Vitest.
  - Supplies complete feature evidence for **M01-C01**, **M01-C02**, and **M01-C03**.

## Open Risks / Decisions for the milestone session

1. **Blocking migration allocation conflict**

   The current accepted checkout already contains committed 0033, 0034, and 0035 migrations, while the supplied task reserves exactly those numbers and requires listings to use 0033. `GOV_MIGRATION_PREFIX` has no applicable exception. The hub must reconcile the base or issue a new valid allocation before Slice 1. F-01 must not overwrite, rename, duplicate, or exceed the allocated range autonomously.

2. **Formal ADR records are missing**

   ADR-12 and ADR-17 are fixed in `.mnfs/.../mission.md`, but absent from the architecture ADR directory/index. Their behavior is clear enough to plan and implement, but the architecture owner should create/ratify the formal records outside this feature’s owned paths.

3. **Connector listing contract lacks required facts**

   The published snapshot currently lacks price, currency, and modality, while IC-02 requires them. Extending the published connector contract and ML adapter is necessary and is an exclusive shared seam. The milestone must confirm Slice 2 owns that seam and no sibling work edits it concurrently.

4. **Provider-auth error vocabulary is not defined**

   The brief requires auth-class operation-run failures, but connectors currently define no auth error code; token resolution becomes transient and HTTP 401/403 becomes validation. The milestone/contract owner must ratify the exact connector `ErrorCode` string before Slice 2. Once fixed, refresh records it in `translated_error_code`; no authorization flow is started.

5. **Modality allowlist and labels are incomplete**

   IC-02 pins nullable `listing_type_code` and gives `gold_pro`/“Premium” as an example, but does not enumerate every recognized ML modality or its label. F-01 needs a ratified code allowlist to distinguish recognized values from unmappable values. Do not pass arbitrary provider strings through or invent labels. F-02 will need the same registry to form `{code,label}`.

6. **Active run ID placement in the 409 body**

   IC-02 requires the body to carry the active `operation_run_id` but does not state whether it is `error.operation_run_id`, `error.details.operation_run_id`, or top-level. The existing shared `ErrorResponse` supports `error.details`; using `error.details.operation_run_id` is the least disruptive repository convention, but the milestone must ratify it before OpenAPI/SDK locking.

7. **Known but non-connected installation status**

   The error matrix fixes missing/unknown/concurrent cases but does not define the response for an existing installation in draft, disconnected, suspended, failed, degraded, or requires-reauth state. RK-06 forbids reauthorization. The milestone must define the HTTP status/code or explicitly declare that only connected installations reach this endpoint in M-01 tests. Do not map these states to `installation_not_found` or start auth without contract authority.

8. **`price_currency` nullability**

   IC-02 marks the price object nullable and ADR-17 says unknown facts stay nullable; the database list explicitly marks `price_amount` nullable but does not annotate `price_currency`. This plan treats both amount and currency as nullable when the provider price fact is absent, avoiding a fabricated currency. Confirm this interpretation before freezing migration 0033.

9. **Listings-to-installation FK/delete policy**

   The table identity includes installation and unknown-installation handling depends on integrations, but IC-02 does not specify an FK or `ON DELETE` behavior. Existing integration-owned tables use tenant-scoped FKs. The milestone must decide whether listings cascade on installation deletion or remain as retained operational history; do not invent the deletion policy in SQL.

10. **Pagination completion semantics**

    The connector accepts an offset cursor but returns only a slice, with no total/next cursor. The planned full pull advances by returned provider items and terminates on a short page. Confirm that every supported `ListingReader` obeys this offset/short-page contract; otherwise the published interface must gain explicit completion metadata before F-01 can safely mark absent rows closed.

11. **Duplicate keys across provider pages**

    IC-02 fixes the canonical key but does not specify duplicate-provider-page behavior. The safe plan fails the operation before persistence when two snapshots map to one canonical key, because last-write-wins could hide provider inconsistency and make closed marking unsafe. Ratify this fail-honest behavior.

12. **Stale queued/running run recovery**

    The contract defines a concurrent guard but not crash recovery for an operation left queued/running after process termination. F-01 should not invent a timeout that could permit overlapping work. The milestone should either accept manual cleanup for M-01 or define a stale-run recovery rule in a later feature.

13. **IC-02 seed link states versus F-01 scope**

    The canonical seed description includes resolved and unlinked examples, but F-01 is explicitly forbidden from touching `product_links`. F-01 can seed the six fixed listing IDs and listings-owned facts only; F-02 must add the link fixtures and joined assertions. This is a scope boundary, not authorization to duplicate link fields.

14. **Verification command wording**

    The F-01 brief requires an absolute `GOCACHE`, while the supplied ladder writes `GOCACHE=.gocache`. On Windows, resolve the repository’s `.gocache` directory to an absolute path and assign that value before Go commands; this preserves the intended cache location while satisfying Go’s absolute-path requirement.

15. **No SDK generator exists**

    `packages/sdk-runtime` has build/test scripts but no OpenAPI generation command. The endpoint slice must update OpenAPI first and then the hand-maintained runtime client in the same commit. Validation must say “contract-first SDK parity,” not claim regeneration that did not occur.

