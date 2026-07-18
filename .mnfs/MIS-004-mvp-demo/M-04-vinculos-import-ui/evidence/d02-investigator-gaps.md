# D-02 Gap-Closer Investigation — M-04 vinculos-import-ui

## A. SDK-runtime structure

### A.1 Generated vs hand-authored
**Verdict: hand-authored.**
- `packages/sdk-runtime/package.json:7-10` — scripts are only `"build": "tsc --noEmit"` and `"test": "vitest run --config vitest.config.ts"`. No codegen/generate script.
- Root `package.json` has no openapi-typescript/orval/codegen script referencing sdk-runtime (grep returned nothing).
- No "DO NOT EDIT"/"GENERATED" header in `src/index.ts` line 1 (`export * from "./erpImport";`) or anywhere in the file — the only hit for "generated" is the field name `generated_count` (index.ts:885), not a file marker.
- Every domain method + type is committed by hand in feature commits (see A.4 history below), each paired with the OpenAPI spec update in the same commit (e.g. commit a2c0698 "IC-01 identity fields in OpenAPI + manual SDK").

### A.2 index.ts size and assembly
- `packages/sdk-runtime/src/index.ts` total: **1804 lines**. `erpImport.ts`: 57 lines.
- `createMarketplaceCentralClient` defined at **index.ts:1322-1325**:
  ```ts
  export function createMarketplaceCentralClient(options: {
    baseUrl: string;
    fetchImpl?: typeof fetch;
  }) {
  ```
  Body (1326-1431) defines local helpers `getJson/putJson/postJson/deleteJson` + query-string builders (`catalogQuery`, `listingQuery`, `syncRunQuery`, `orderQuery`, `mutationQuery`), then `return { ... }` (starts index.ts:1433) is one giant object literal where **every domain method's implementation is written inline** — there is no per-domain file split for method implementations (only `erpImport.ts` holds pure types, no functions).
- Example inline method entries (index.ts:1434-1457 region): `listCatalogProductFacts` (1434), `searchCatalogProductFacts` (1438), `listListings` (1442), `getDashboardSummary` (1454), `listSyncRuns` (1456).
- Existing product-links client methods, all inline in the returned object literal:
  - `importProductLinkListingSnapshots` — index.ts:1508-1509
  - `generateProductLinkCandidates` — index.ts:1510-1511
  - `listProductLinkCandidates` — index.ts:1512-1515
  - `listProductLinkWorkflows` — index.ts:1516-1519
  - `approveProductLinkCandidate` — index.ts:1520-1521
  - `rejectProductLinkListing` — index.ts:1522-1529
  - `manualResolveProductLink` — index.ts:1530-1545 (includes inline validation: throws `Error("invalid_identity: ...")` if `internal_product_id` isn't a positive integer, before delegating to `postJson`)

  Associated request/response types are declared earlier in the same file at index.ts:850-949 (`ProductLinkListingSnapshotImportResult` 850, `ProductLinkCandidateState`/`ProductLinkCandidateMatchInput` 856/864, `ProductLinkCandidateItem` 866, `ProductLinkCandidateGenerationResult` 883, `ProductLinkState`/`ProductLinkAction` 889/891, `ProductLinkActor` 893, `ProductLinkListingIdentity` 899, `ProductLink` 905, `ProductLinkAuditEntry` 919, `ProductLinkWorkflowItem` 936, `ProductLinkResolutionResult` 943).

### A.3 erpImport.ts
- Exports **types only** — `ErpImportStatus`, `ErpImportIssueCode`, `ErpImportIssue`, `ErpImportSummary`, `ErpImportDetail`, `ErpImportList`, `ErpImportCreated`, `ErpImportError`, `ErpImportConflict` (erpImport.ts:1-58). No functions, no client methods.
- index.ts consumes it via a single barrel re-export line: `index.ts:1` → `export * from "./erpImport";`. There is no `import ... from "./erpImport"` anywhere else in index.ts (grep for `ErpImport` inside index.ts returned no matches beyond that) — i.e. the erp-import client methods themselves (if any exist) are NOT in erpImport.ts; only its types are surfaced through the barrel. (No `erpImport*` methods were found inline in the returned client object during this pass — out of scope of this query but worth flagging if F-02 assumed otherwise.)
- Git history: `erpImport.ts` created whole-cloth in commit 894ece73 "feat(erp-import): OpenAPI /erp/imports* + manual SDK types (F02-S7)"; the barrel line was added one commit later in 332c592 "feat(sdk): BARREL-01 — re-export erpImport types from package root" (explicitly noted as a **hub-owned** edit granted via HUB-LEDGER D-13).

### A.4 Recipe to add a new domain client (pattern = product-links / IC-01, most recently added domains)
1. Add/extend request+response **types** directly in `index.ts` (near other domain types) — see commit a2c0698 (IC-01, 25 lines added straight into index.ts) and the product-links type block at 850-949.
2. Add the **method implementation inline** inside the `return { ... }` object literal of `createMarketplaceCentralClient` (index.ts:1433+), following the `getJson<T>/postJson<T>` helper pattern — see product-links methods at 1508-1545.
3. Update **OpenAPI spec in the same commit** (ADR-12 discipline; confirmed by a2c0698 commit message: "OpenAPI spec and the manual sdk-runtime types land in the same commit").
4. No regeneration step exists — `packages/sdk-runtime` has no build/codegen script beyond `tsc --noEmit`; changes are hand-typed and type-checked, not generated.
5. Exception pattern (erpImport.ts): a **standalone types-only file** was created once (F02-S7) and then barrel-re-exported in a *separate, hub-owned* commit (BARREL-01) — this is not the default path; it required explicit hub grant (HUB-LEDGER D-13). Default/expected path for a new domain (e.g. product-links batch resolutions) is inline-in-index.ts per steps 1-2, not a new per-domain file.

### A.5 All barrel re-export lines in index.ts
- Only one: `index.ts:1` → `export * from "./erpImport";`
- No other `export * from` lines exist in the file (grep confirmed single match).

---

## B. Matcher / anchor data

### B.1 ProductMatcher / FindProductsForLinking signatures
- Port interface (canonical): `apps/server_core/internal/modules/internal_read/ports/reader.go:48-49`
  ```go
  type Reader interface {
      FindProductsForLinking(ctx context.Context, input FindProductsInput) ([]domain.ProductCandidate, error)
      ...
  }
  ```
  `FindProductsInput` struct: `reader.go:9-15` — `ProductID *int`, `EAN *string`, `SellerSKU *string`, `Title *string`, `IncludeInactive bool`.
- Narrower local interface used by product-links application layer: `apps/server_core/internal/modules/product_links/application/generation_service.go:17-19`
  ```go
  type ProductMatcher interface {
      FindProductsForLinking(ctx context.Context, input internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error)
  }
  ```
- Return struct `domain.ProductCandidate`: `apps/server_core/internal/modules/internal_read/domain/internal_product.go:3-21`
  ```go
  type ProductCandidate struct {
      InternalProductID *InternalProductID
      ProductID         int
      Name              string
      EAN               *string
      ReferenceCode     *string
      NCM               *string
      ProductGroupID    *int
      ProductGroupName  *string
      BrandID           *int
      BrandName         *string
      IsActive          bool
      UsageType         *string
      Source            SourceMetadata
      QualityFlags      []QualityFlag
  }
  ```
  **Yes** — brand (`BrandName`/`BrandID`), manufacturer reference (`ReferenceCode`), `NCM`, and `EAN` are all exposed per matched product.
- Wiring in composition root: `apps/server_core/internal/composition/root.go:398-457`. `productMatcher` starts as `unavailableProductMatcher` (line 398), then gets assigned either the xlsx-backed `internalReadSvc` (line 419, when `source == "xlsx"`) or the Oracle-backed `internalReadSvc` (line 438, on successful Oracle connection). `productLinkGenerationSvc` is built at line 457 with `Matcher: productMatcher`.

### B.2 Adapter columns (Oracle)
- `apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go:399-419` — SQL SELECT list in `buildFindProductsQuery`:
  ```sql
  SELECT
      p.CODPROD, p.DESCRPROD, p.REFERENCIA, p.REFFORN, p.NCM,
      p.CODGRUPOPROD, p.MARCA, p.CODMARCA, p.ATIVO, p.USOPROD,
      ( SELECT COUNT(DISTINCT collision.CODPROD) ... ) AS EAN_ACTIVE_COUNT
  FROM METALPRD.TGFPRO p
  ```
  Scan targets at reader.go:52-64 map these into `productID, name, eanValue(REFERENCIA-derived), referenceValue(REFFORN), ncmValue, productGroupID, brandName(MARCA), brandID(CODMARCA), activeValue, usageType, activeEANCount`.
- **Confirmed**: `p.MARCA`/`p.CODMARCA` (brand) reach the matcher via `BrandName`/`BrandID` (reader.go:59-60, 85-86); `p.REFFORN` (manufacturer_reference) reaches the matcher via `ReferenceCode` (reader.go:56, 83).
- Caveat: field naming is non-obvious — `REFERENCIA` (not `REFFORN`) is scanned into `eanValue`/`ean` (reader.go:55, 69), and `REFFORN` is scanned into `referenceValue`/`ReferenceCode` (reader.go:56, 83). There is an explicit code comment at reader.go:431-432: `// EAN is intentionally not matched against TGFPRO.REFERENCIA. Until a governed barcode source is added, EAN-only linking remains unproved.` — i.e. today's "EAN" field is really `REFERENCIA` data repurposed, and true EAN-based WHERE-clause matching is deliberately unsupported (reader.go:441-443 returns a `ReadErrorUnsupportedQuery` if only EAN is supplied).

### B.3 Listing-snapshot columns
- Confirmed in one line: migration `0022_product_links_listing_snapshots.sql:1-17` and `product_links/domain/listing_snapshot.go:5-19` both define only `provider_status, seller_sku, ean, title, available_quantity` (+ identity/timestamps) — **no brand or manufacturer_reference column exists on the listing-snapshot side**.

### B.4 GenerateLinkCandidates decision ladder (generation_service.go:58-160, actually spans 58-146 for the core ladder)
- Ladder, in `generateForSnapshot` (generation_service.go:120-145): try exact SKU match (`findProducts(..., SellerSKU)`, line 123) and exact EAN match (line 127) in parallel → `buildExactCandidates` (132) which returns `exact_sku`/`exact_ean` state if exactly one product matches on either axis, or `conflict` state if either axis yields >1 product OR the two axes disagree on product identity (buildExactCandidates, lines 192-213; conflict logic at 216-232). If no exact/conflict result, falls through to title LIKE-match (line 137) → `title_match` state (142). If nothing matches, `unresolved` state with `ProductCandidate{}` empty (145).
- **Data in hand at decision time**: only the snapshot's `SellerSKU`, `EAN`, `Title` strings (each queried independently against Oracle/xlsx via `FindProductsInput{SellerSKU|EAN|Title}` — reader.go:433-440 builds an OR'd WHERE across whichever single field is populated per call) and, once matched, the full `ProductCandidate` (including brand/NCM/reference) is available in the closure at `newCandidate` (line 242) but **only `InternalProductName` and `InternalReferenceCode` are persisted into `domain.LinkCandidate`** (lines 260-262) — brand/NCM/quality-flags are read but dropped, not stored on the candidate.
- **No per-candidate multi-field comparison exists today.** Each candidate carries exactly one `match_input` anchor (`ProductLinkCandidateMatchInput`: "manual"|"seller_sku"|"ean"|"title"|"none", index.ts:864) and one `match_value` — there is no scored/weighted comparison of brand+reference+ncm+ean together; conflict detection (buildConflictCandidates, 216-232) only checks whether SKU-axis and EAN-axis point to different product IDs, not a confidence blend across identity fields.

---

## C. FE patterns for F-02

### C.1 DetailDrawer / DetailPanel
- `DetailDrawer` props: `packages/ui/src/DetailDrawer.tsx:4-13` — `{ open, onClose, title, subtitle?, children, actions?, closeLabel?, width? }`. It's a thin wrapper that renders `DetailPanel` with `footer={actions}` (DetailDrawer.tsx:26-36) and default `width ?? 360`.
- `DetailPanel` props: `packages/ui/src/DetailPanel.tsx:4-13` — `{ open, onClose, title, closeLabel="Close panel", subtitle?, children, footer?, width=380 }`. Implements Escape-to-close (useEffect, lines 25-32), fixed right-side panel (`role="complementary"`, lines 36-42), header with close `X` button, scrollable body, optional footer.
- **DetailDrawer itself has no live callsite in `apps/web/src`** (only its own test file references it) — grep for `DetailDrawer` under apps/web/src returned no matches. **DetailPanel is the one actually used**, via the page-local `ListingDetailPanel` wrapper:
  - `apps/web/src/pages/ListingDetailPanel.tsx:10,182-191` imports `DetailPanel` from `@marketplace-central/ui` and wraps it with listing-specific content/props (`ListingDetailPanelProps` at line 20, component at line 167).
  - Usage from a page: `apps/web/src/pages/AnunciosPage.tsx:18,327`:
    ```tsx
    import { ListingDetailPanel } from "./ListingDetailPanel";
    ...
    <ListingDetailPanel listingId={openListingId} onClose={() => setOpenListingId(null)} />
    ```

### C.2 Tab / Modal / Toast
- **Tab**: no shared Tab component in `packages/ui` (barrel `packages/ui/src/index.ts` has no Tab export). Pattern is page-local: `apps/web/src/pages/AnunciosPage.tsx:24` defines a local `tabs: Array<{ value: AnunciosTab; label: string }>` array, rendered at lines 183-196 with `role="tab"`, `aria-selected={state.tab === tab.value}`, and `onClick={() => updateState({ ...state, tab: tab.value })}` — plain local state, no abstraction.
- **Modal/Dialog**: no shared Modal/Dialog component in `packages/ui` either. Pattern is page-local: `apps/web/src/pages/mutations/MutationPreviewModal.tsx:127` — `<section role="dialog" aria-modal="true" aria-labelledby="mutation-modal-title" className="...rounded-xl bg-white shadow-lg">`, hand-built with no shared primitive.
- **Toast/snackbar/notification**: **none found** — grep across `apps/web/src` and `packages/ui/src` for toast/snackbar/notification (case-insensitive, .tsx) returned zero files.
- Conclusion: for tabs and modals, F-02 should follow the existing page-local pattern (plain state + `role="tab"`/`role="dialog"` + Tailwind classes) rather than expect a shared component; there is no toast mechanism at all, so any success/error feedback needs a new pattern or reuse of `ErrorState`/`EmptyState`/`FreshnessIndicator` style inline banners (packages/ui/src barrel lines 1-10).

### C.3 web-query: linkageQueryKeys, invalidateAfterMutation, QUERY_STALE_TIME
- `linkageQueryKeys`: `packages/web-query/src/index.ts:38-40`
  ```ts
  export const linkageQueryKeys = {
    workflows: (installation_id: string) => ["linkage", { installation_id }] as const,
  };
  ```
- `QUERY_STALE_TIME`: `packages/web-query/src/index.ts:4-13` — keys: `catalog: 300_000, stock: 45_000, pricecost: 120_000, listings: 45_000, mutations: 5_000, orders: 120_000, sync: 30_000, market: 300_000`. **No `linkage` key exists yet** — F-02 will need to add one if it wants a distinct linkage stale-time (currently nothing maps to `linkageQueryKeys`'s namespace for staleTime purposes).
- `invalidateAfterMutation` signature: `packages/web-query/src/invalidation.ts:36-39`
  ```ts
  export async function invalidateAfterMutation(
    queryClient: QueryClient,
    type: MutationInvalidationType,
  ): Promise<void>
  ```
  `MutationInvalidationType` (invalidation.ts:5-12) already includes `"link_apply"` mapped to namespaces `["listings", "inventory"... ]` — actually `link_apply: ["listings", "linkage", "catalog", "mutations"]` (invalidation.ts:26) — so the `linkage` namespace is already wired into the invalidation map, just not into `QUERY_STALE_TIME`.
- Example invocation from a page/hook: `apps/web/src/pages/mutations/useMutationProtocol.ts:43`:
  ```ts
  void invalidateAfterMutation(queryClient, protocol.type as MutationInvalidationType);
  ```
  called inside a `useEffect` (lines 33-44) gated on mutation terminal state and a `WeakMap<QueryClient, Set<string>>` de-dupe guard (lines 13, 36-42) to invalidate exactly once per protocol.

### C.4 Page-local component pattern confirmation
- Confirmed: `apps/web/src/pages/mutations/` is a page-local component folder owned by `ProtocoloPage.tsx`, containing: `MutationBulkActions.tsx`, `MutationIntentForm.tsx`, `MutationItemsTable.tsx`, `mutationPresentation.ts`, `MutationPreviewModal.tsx` (+ tests), `MutationResultSummary.tsx`, `useMutationProtocol.ts` (+ test). This is the precedent for F-02 to place a `pages/vinculos/` (or similar) folder with its own sub-components rather than routing everything through `packages/ui`.
- Simpler precedent: `apps/web/src/pages/ListingDetailPanel.tsx` is a single page-local file (not a folder) sitting directly under `pages/`, imported by `AnunciosPage.tsx:18` — a lighter-weight pattern for a single local sub-component.
