## 1. `product_links` module map

- `domain/`: `link_candidate.go`, `product_link.go`, `listing_snapshot.go`, `internal_product_id.go`.
  - `LinkCandidate`: `CandidateID`, installation/provider/item/variation IDs, `InternalProductID *int`, name/reference, state, match input/value, snapshot/created/updated timestamps (`domain/link_candidate.go:31-46`).
  - `ProductLink`: identity, state, source candidate, internal product ID/name/reference, timestamps (`domain/product_link.go:28-40`).
  - `ProductLinkAuditEntry`, `ProductLinkTransition`, `ProductLinkWorkflowItem`, `ProductLinkResolutionResult` (`domain/product_link.go:42-74`).
  - `ValidateInternalProductID(id int) error` (`domain/internal_product_id.go:7`).

- `application/`:
  - `ImportService.ImportListingSnapshots(ctx, ImportListingSnapshotsInput) (ListingSnapshotImportResult, error)` (`import_service.go:47`).
  - `GenerationService.GenerateLinkCandidates(ctx, GenerateLinkCandidatesInput) (LinkCandidateGenerationResult, error)` (`generation_service.go:58`).
  - `GenerationService.ListLinkCandidates(ctx, ListLinkCandidatesInput) ([]LinkCandidate, error)` (`generation_service.go:105`).
  - `ResolutionService.ApproveCandidate(...)`, `RejectListing(...)`, `ManualResolve(...)`, each returning `ProductLinkResolutionResult, error` (`resolution_service.go:77,130,157`).
  - `ResolutionService.ListLinkWorkflows(...) ([]ProductLinkWorkflowItem, error)` (`resolution_service.go:191`).
  - `SummaryService.Summary(ctx, installationID) (ports.LinkageSummary, error)` (`summary_service.go:22`).
  - Generation order: seller SKU → EAN → title → unresolved (`generation_service.go:120-145`).

- `ports/`: `link_candidate_store.go`, `listing_snapshot_store.go`, `summary_reader.go`, `workflow_store.go`.
  - Candidate store: replace/list/get.
  - Workflow store: get link, apply transition, list links, list audit entries (`ports/workflow_store.go:10-17`).

- `adapters/postgres/`: `link_candidate_repo.go`, `listing_snapshot_repo.go`, `summary_reader.go`; integration test.
  - `ReplaceLinkCandidates`, list/get candidates, get/list links, `ApplyProductLinkTransition`, list audit entries (`link_candidate_repo.go:26-389`).
  - Transition persists `product_links` and `product_link_audit_entries` (`link_candidate_repo.go:243-321`).

- `transport/`: `http_handler.go`; handler interfaces and route registration (`http_handler.go:17-65`).

Resolution routes/application calls:

- `POST /product-links/link-resolutions/approve-candidate` → `ApproveCandidate` (`http_handler.go:62,223-252`).
- `POST /product-links/link-resolutions/reject-listing` → `RejectListing` (`http_handler.go:63,255-290`).
- `POST /product-links/link-resolutions/manual-resolve` → `ManualResolve` (`http_handler.go:64,292-336`).
- Generation: `POST /product-links/link-candidates/generations` → `GenerateLinkCandidates` (`http_handler.go:59,143-175`).
- Workflow listing: `GET /product-links/link-workflows` → `ListLinkWorkflows` (`http_handler.go:61,201-220`).

## 2. Product-link DB tables

Migrations:

- `0022_product_links_listing_snapshots.sql`
- `0023_product_link_candidates.sql`
- `0025_product_link_workflows.sql`

`0022_product_links_listing_snapshots.sql:1-15`

```sql
CREATE TABLE IF NOT EXISTS product_link_listing_snapshots (
  tenant_id text NOT NULL,
  installation_id text NOT NULL,
  provider_code text NOT NULL,
  provider_item_id text NOT NULL,
  provider_variation_id text NOT NULL DEFAULT '',
  provider_status text NOT NULL DEFAULT '',
  seller_sku text NOT NULL DEFAULT '',
  ean text NOT NULL DEFAULT '',
  title text NOT NULL DEFAULT '',
  available_quantity integer,
  source_updated_at timestamptz,
  fetched_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, installation_id, provider_item_id, provider_variation_id)
);
```

`0023_product_link_candidates.sql:1-17`

```sql
CREATE TABLE IF NOT EXISTS product_link_candidates (
  tenant_id text NOT NULL,
  candidate_id text NOT NULL,
  installation_id text NOT NULL,
  provider_code text NOT NULL,
  provider_item_id text NOT NULL,
  provider_variation_id text NOT NULL DEFAULT '',
  internal_product_id integer,
  internal_product_name text NOT NULL DEFAULT '',
  internal_reference_code text NOT NULL DEFAULT '',
  state text NOT NULL,
  match_input text NOT NULL,
  match_value text NOT NULL DEFAULT '',
  source_snapshot_fetched_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, candidate_id)
);
```

`0025_product_link_workflows.sql:1-39`

```sql
CREATE TABLE IF NOT EXISTS product_links (
  tenant_id text NOT NULL,
  installation_id text NOT NULL,
  provider_code text NOT NULL,
  provider_item_id text NOT NULL,
  provider_variation_id text NOT NULL DEFAULT '',
  state text NOT NULL,
  source_candidate_id text NOT NULL DEFAULT '',
  internal_product_id integer,
  internal_product_name text NOT NULL DEFAULT '',
  internal_reference_code text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, installation_id, provider_item_id, provider_variation_id)
);

CREATE TABLE IF NOT EXISTS product_link_audit_entries (
  tenant_id text NOT NULL,
  audit_id text NOT NULL,
  installation_id text NOT NULL,
  provider_code text NOT NULL,
  provider_item_id text NOT NULL,
  provider_variation_id text NOT NULL DEFAULT '',
  action text NOT NULL,
  reason text NOT NULL DEFAULT '',
  source_candidate_id text NOT NULL DEFAULT '',
  actor_type text NOT NULL DEFAULT '',
  actor_id text NOT NULL DEFAULT '',
  actor_name text NOT NULL DEFAULT '',
  previous_state text NOT NULL,
  next_state text NOT NULL,
  previous_internal_product_id integer,
  next_internal_product_id integer,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, audit_id)
);
```

Current top migration: `0047_create_erp_import_issues.sql`.

Hardcoded migration count: `44` at `apps/server_core/internal/platform/migrate/runner_test.go:25` and `:64`.

## 3. OpenAPI paths

`contracts/api/marketplace-central.openapi.yaml`:

- `/product-links/listing-snapshots/imports` POST, `importProductLinkListingSnapshots`; request `ImportProductLinkListingSnapshotsRequest` (`installation_id`, `limit`), response `ProductLinkListingSnapshotImportResult` (`installation_id`, `imported_count`, `items`) (`:1364-1380`, schemas `:5029-5055`).
- `/product-links/link-candidates/generations` POST, `generateProductLinkCandidates`; request `GenerateProductLinkCandidatesRequest` (`installation_id`, `limit`), response `ProductLinkCandidateGenerationResult` (`installation_id`, `generated_count`, `items`) (`:1381-1397`, schemas `:5105-5131`).
- `/product-links/link-candidates` GET, `listProductLinkCandidates`; response `ListProductLinkCandidatesResponse` (`items: ProductLinkCandidate`) (`:1398-1421`, schemas `:5057-5103`, `:5133-5141`).
- `/product-links/link-workflows` GET, `listProductLinkWorkflows`; response `ListProductLinkWorkflowsResponse` (`items: ProductLinkWorkflowItem`) (`:1422-1445`, schemas `:5258-5286`).
- `/product-links/link-resolutions/approve-candidate` POST, `approveProductLinkCandidate`; request `ApproveProductLinkCandidateRequest` (`candidate_id`, `reason`, `actor`), response `ProductLinkResolutionResult` (`link`, `audit`) (`:1446-1462`, schemas `:5565-5576`, `:5288-5297`).
- `/product-links/link-resolutions/reject-listing` POST, `rejectProductLinkListing`; request `RejectProductLinkListingRequest` (installation/provider/item/variation, reason, actor), response `ProductLinkResolutionResult` (`:1463-1479`, schemas `:5578-5598`).
- `/product-links/link-resolutions/manual-resolve` POST, `manualResolveProductLink`; request `ManualResolveProductLinkRequest` (identity, internal product ID/name/reference, reason, actor), response `ProductLinkResolutionResult` (`:1480-1496`, schemas `:5599-5630`).

ERP paths:

- `/erp/imports` POST, `createErpImport`; response `ErpImportCreated` (`import_id`, `protocol`, `status`) (`:2499-2521`).
- `/erp/imports` GET, `listErpImports`; response `ErpImportList` (`items: ErpImportSummary`) (`:2546-2557`).
- `/erp/imports/{id}` GET, `getErpImport`; response `ErpImportDetail` (`ErpImportSummary` plus `rejected_rows`, `warnings`) (`:2565-2584`, schemas `:6198-6269`).

Protocol shape: string matching `^#[0-9]{3,}-E$` in migration `0045...:15`; statuses `COMPLETED|REJECTED` (`openapi:6183-6185`); counts `accepted_count`, `rejected_count`, `warning_count`; rejection/warning issue fields `row`, `code`, `detail`, nullable `column`, nullable `offending_value` (`openapi:6229-6244`).

## 4. `packages/sdk-runtime`

`src/*.ts`:

- `erpImport.ts`
- `index.ts`

`src/productLinks.ts` does not exist.

Pattern: `createMarketplaceCentralClient({baseUrl, fetchImpl?})` returns an object of arrow-function clients (`index.ts:1322-1326`). Generic `getJson<T>` performs fetch, parses JSON, and throws `{status, error}` on non-2xx (`index.ts:1328-1335`). `erpImport.ts` contains exported type/interface declarations only (`erpImport.ts:1-57`).

Barrel example: `export * from "./erpImport";` (`index.ts:1`).

## 5. Web seams

`apps/web/src/routes/vinculos.tsx`:

```tsx
import { WorkspacePlaceholder } from "../pages/WorkspacePlaceholder";

export function VinculosRoute() {
  return <WorkspacePlaceholder />;
}
```

Routes:

- `anuncios.tsx`
- `dashboard.tsx`
- `pedidos.tsx`
- `precos.tsx`
- `produto.tsx`
- `vinculos.tsx`

`apps/web/src/pages/vinculos/` is absent. Existing page seam is flat files; `AnunciosPage.tsx` uses `useQuery`, `useClient`, `listingsQueryKeys`, `QUERY_STALE_TIME.listings`, and `LoadingState`/`ErrorState`/`EmptyState` (`apps/web/src/pages/AnunciosPage.tsx:1-7,99-104,274-286`).

`packages/web-query/src/index.ts` exports:

- `QUERY_STALE_TIME`
- `queryKeyNamespaces`
- `catalogQueryKeys`, `inventoryQueryKeys`, `linkageQueryKeys`, `profitabilityQueryKeys`, `listingsQueryKeys`, `mutationsQueryKeys`, `ordersQueryKeys`, `syncQueryKeys`, `installationsQueryKeys`
- `createWebQueryClient`
- `FreshnessIndicator`
- `invalidateAfterMutation` and invalidation types (`:4-78,137-142`).

`packages/ui/src/index.ts` exports `LoadingState`, `ErrorState`, `EmptyState`, `UnknownValue`, `FreshnessIndicator`, `ConflictTag`, `MarginChip`, `Button`, `SurfaceCard`, `Badge`, `StatCard`, `ProductPicker`, `PaginatedTable`, `DetailPanel`, `DataTable`, `DetailDrawer` (`:1-22`). No tabs/modal/toast primitive files are present.

## 6. Catalog identity

`GET /catalog/products` returns `CatalogProductFactPage`; `/catalog/products/{id}` returns `CanonicalCatalogProduct` (`openapi:338-404`).

Identity fields:

- `ean` exists.
- `ncm` exists.
- `refforn` and `marca` do not exist under those names in the catalog response; equivalents are `manufacturer_reference` and `brand_name` (`openapi:2904-2935`, `2942-2973`).
- Go response struct: `catalogProductFactResponse` has `ManufacturerReference`, `EAN`, `BrandName`, `NCM` JSON fields (`apps/server_core/internal/modules/catalog/transport/http_handler.go:213-226`).
- Canonical domain struct has the same identity fields (`apps/server_core/internal/modules/catalog/domain/canonical_product.go:69-82`).

Candidate generation does not directly read the catalog handler/schema. It calls injected `ProductMatcher.FindProductsForLinking` for seller SKU, EAN, and title (`product_links/application/generation_service.go:17-19,120-145`). Composition wires this to internal-read service / Oracle fallback (`internal/composition/root.go:398-455`).

## 7. Governance/composition

`contracts/governance/modules.json:19-21`:

```json
{
  "id": "product_links",
  "root": "apps/server_core/internal/modules/product_links",
  "code_owner_path": "apps/server_core/internal/modules/product_links",
  "composition_required": true,
  "openapi_prefixes": ["/product-links"],
  "dependencies": ["connectors", "internal_read"]
}
```

Composition root: `apps/server_core/internal/composition/root.go`.

- Imports product-links adapters/application/transport (`:101-103`).
- Constructs snapshot/candidate repositories (`:391-396`).
- Constructs generation and resolution services (`:457-465`).
- Registers handler: `productlinkstransport.NewHandler(...).Register(mux)` (`:467`).

## 8. Existing audit/protocol pattern

ERP import protocol pattern:

- Tables: `erp_import_protocols`, `erp_import_products`, `erp_import_issues`.
- Protocol columns: `id`, `tenant_id`, `file_sha256`, `protocol`, `source`, `imported_at`, `status`, accepted/rejected/warning counts (`migrations/0045_create_erp_import_protocols.sql:1-18`).
- Issues columns: `tenant_id`, `protocol_id`, `row_number`, `column_name`, `kind`, `code`, `detail`, `offending_value` (`migrations/0047_create_erp_import_issues.sql:1-31`).
- Store pattern: `PersistSnapshotAtomically(ctx, tenantID, snapshot)` opens transaction, tenant advisory lock, duplicate-file lookup, inserts protocol/products/issues, commits (`erp_import/adapters/postgres/import_repository.go:18-66`).
- Query pattern scopes all reads by tenant and protocol ID (`erp_import/adapters/postgres/query_repository.go:24-99`).

Product-links already mirrors this with `product_link_audit_entries`, storing previous/next state, actor, reason, and product IDs (`0025_product_link_workflows.sql:20-39`).

## 9. `tenant_id` scoping

Repositories capture tenant at construction:

```go
func NewLinkCandidateRepository(pool *pgxpool.Pool, tenantID string) *LinkCandidateRepository
```

`apps/server_core/internal/modules/product_links/adapters/postgres/link_candidate_repo.go:17-24`.

Example query:

```sql
SELECT ...
FROM product_link_candidates
WHERE tenant_id = $1
  AND installation_id = $2
```

`link_candidate_repo.go:97-111`.

Product-link lookup similarly requires `tenant_id = $1` plus installation/item/variation identity (`link_candidate_repo.go:200-214`).