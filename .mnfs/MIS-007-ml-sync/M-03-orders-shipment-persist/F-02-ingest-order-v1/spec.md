# F-02-ingest-order-v1 — spec

Feature 2 of 3 in the serial DAG (F-01→F-02→F-03) for MIS-007/M-03-orders-shipment-persist.
Depends on F-01 (`connectors/adapters/mercado_livre/order_ingest_reader.go`,
`shipment_ingest_reader.go` — CLOSED @cf2c09e, not modified here).

## Goal

`IngestOrder(ctx, installationID, providerOrderID) error` — the single write path (ADR-04)
for one order: fetch OrderDetail (F-01), conditionally fetch ShipmentDetail (F-01), fetch
BuyerFiscal (existing `ports.BuyerFiscalReader`), resolve SKU links (existing
`ports.LinkReader`), read the current `faturado_at` watermark, derive the bucket via
`domain.DeriveOrderBucket` verbatim, and persist header + items + payments + 0089 columns +
an optional `order_shipments` row in ONE transaction.

```go
// orders/ports/order_ingestor.go
func (s *IngestService) IngestOrder(ctx context.Context, installationID, providerOrderID string) error
```

`ImportService.Import` is refactored to enumerate provider order ids from `ports.OrderSource`
(unchanged) and delegate each id to `ports.OrderIngestor.IngestOrder`, aggregating
imported/skipped counts. `normalizeOrders` (the old parallel write path) is deleted in the
same diff, per ADR-04 ("one write path only").

## New files

- `orders/domain/order_shipment.go` — `OrderShipment` (0088 shape, module-owned, not a reuse
  of `connectorsdomain.ShipmentDetail`).
- `orders/domain/order_ingest_errors.go` — `ErrOrderUnavailable` (403/404 order-level
  classification for IngestOrder callers).
- `orders/ports/order_detail_reader.go`, `orders/ports/shipment_detail_reader.go`,
  `orders/ports/order_ingestor.go` — new ports.
- `orders/application/ingest_service.go` — `IngestService`, the new writer.
- `orders/application/ingest_service_test.go` — DB-independent unit tests (fakes).
- `orders/adapters/postgres/ingest_repo_test.go` — `//go:build integration` DB tests for the
  repo-level persistence/atomicity/idempotency scenarios feature.md's Validation Expectations
  ask for.
- `composition/orders_ingest_adapters.go` — composition adapters wrapping
  `CapabilityAdapter.GetOrderDetail`/`GetShipmentDetail`, mirroring
  `orders_adapters.go`'s `accountRefForInstallation` pattern (kept in a separate file per the
  dispatch brief, `orders_adapters.go` itself untouched).

## Edited files

- `orders/domain/order.go` — `MarketplaceOrder` extended with 13 new 0089 fields (`PackID`,
  `ProviderShipmentID`, `Bucket`, `DateLastUpdatedML`, 9 `Buyer*` fiscal fields).
  `net_amount`/`margin_pct`/`decomposition` (0089, M-06 scope) are deliberately absent.
- `orders/ports/order_store.go` — new `OrderIngestStore` interface (`GetFaturadoAt` +
  `IngestOrder`).
- `orders/adapters/postgres/order_repo.go` — `upsertOrder`'s SQL extended to the 13 new
  columns (COALESCE-never-erase for 12 of them, `bucket` always fresh); new
  `upsertOrderShipment` (0088 upsert, `fetched_at`-watermarked); new `GetFaturadoAt`
  (`pgx.ErrNoRows` → `nil, nil`, honest-absence for first ingest); new `IngestOrder` (one tx:
  `upsertOrder` → items/payments or backfill → optional shipment → commit).
- `orders/application/import_service.go` — rewritten: `ImportServiceConfig{Source, Ingestor,
  Now}`, enumerate-and-delegate `Import`. `normalizeOrders`, `collectIdentities`,
  `trimNonEmpty` deleted (repo-wide grep confirmed no other callers).
- `orders/application/import_service_test.go` — rewritten around the new shape;
  `floatPtr`/`intPtr`/`stubLinkReader` (package-wide test helpers this file already owned,
  also consumed by F-03's `enrich_service_test.go`/`enrich_link_refresh_test.go`) kept intact.
- `internal/composition/root.go` lines 578-603 — wires the two new readers + `IngestService`
  + the updated `ImportServiceConfig`; every line outside that range is untouched (confirmed
  via `git diff`).

## Ground truth consulted

- `.mnfs/MIS-007-ml-sync/research/orders-persistence-interface-contract.md` (IC-03) — 0088/0089
  column set, "typed only, no raw" rule for shipment PII.
- `.mnfs/MIS-007-ml-sync/research/sync-ingest-ports-interface-contract.md` (IC-06) —
  `IngestOrder` signature, "erro por item no ingest: registra, segue o batch" (not scoped to
  403/404 alone — see Design decisions below).
- `migrations/0088_order_shipments.sql`, `migrations/0089_orders_marketplace_orders_sync_fields.sql`
  — exact column list/types (`currency char(3)`, no FK, additive-only).
- `orders/domain/order_bucket.go` — `DeriveOrderBucket` signature and bucket vocabulary,
  called verbatim, never reimplemented.
- `orders/adapters/postgres/order_repo.go` (pre-edit, 952 lines) — `upsertOrder`,
  `replaceItems`/`replacePayments`/`backfillMissingLines`, `nullableTime`/`nullableFloat8`
  helpers, `*string` pgx parameter convention (confirmed against
  `market/adapters/postgres/observation_repository.go`).
- `connectors/adapters/mercado_livre/shipment_ingest_reader.go` — confirmed 404/410 on the
  shipment call degrades internally to `ShipmentDetail{Found: false}, nil` (never a Go error),
  which is what let `IngestOrder` treat every OTHER error uniformly (abort-before-write, no
  per-site special-casing needed).

## Design decisions (brief left room for judgment)

1. **ImportService's skip/continue scope.** The brief's own example is 403/404-only, but
   IC-06's line "erro por item no ingest: registra, segue o batch (run não aborta)" is not
   qualified to that case. `ImportService.Import` treats ANY error from `IngestOrder` as
   skip+log+continue — never aborting the whole run — while `IngestOrder` itself still wraps
   the order-level 403/404 case in the typed `domain.ErrOrderUnavailable` for diagnostic
   log distinction only, not different control flow at the `ImportService` level.
   (`import_service_test.go`'s `TestImportServiceCountsAnyIngestFailureAsSkip` names this.)
2. **SaleFeeAmount mapping.** `OrderDetailItem.SaleFeeUnit` is passed straight through into
   `MarketplaceOrderItem.SaleFeeAmount`, unchanged (no `×Quantity`). The `sale_fee_amount`
   column has no documented semantic (migration 0027 has no comment), the pre-refactor path
   did the same pass-through, and `SaleFeeUnit`'s own doc comment defers any total-multiply to
   a downstream/read-side decision. Flagged as an honest gap, not silently resolved.
3. **Buyer fiscal is always fetched.** Unlike `EnrichService`'s list-path fiscal-skip
   (FINDING-M08-LIST-TIMEOUT, a different feature), `IngestOrder` runs off the interactive
   request path (import/backfill/sync trigger, never a page load), so there is no timeout
   budget forcing a skip — IC-03's Operations table lists fiscal as part of ingest, not an
   optional enrichment.
4. **Atomicity is sequencing, not special-casing.** All three fetches (OrderDetail,
   ShipmentDetail, BuyerFiscal) happen before any write. Combined with the F-01 reader's
   documented honest-absence degrade on shipment 404/410, this means ANY error that DOES
   propagate from a fetch is inherently "real" (auth/transient/validation) and safe to treat
   uniformly: abort before writing, no per-fetch-site branching required.

## Not touched (forbidden paths, verified)

- `capability_adapter.go`, `shipping_reader.go` — 0 edits (M-01-owned).
- `connectors/**` — 0 edits (F-01 closed).
- `orders/application/enrich_service.go`, `orders/application/enrich_service_test.go`,
  `orders/application/enrich_link_refresh_test.go`, `orders/transport/archguard_test.go` —
  0 edits (F-03's files). Two package-level test helpers (`floatPtr`, `intPtr`,
  `stubLinkReader`) that those F-03 test files depend on live in
  `import_service_test.go` (this feature's file) and were preserved across the rewrite —
  see validation.md's `go vet` section for how this was caught.
- `root.go` — only lines 578-603 (confirmed via `git diff`, see validation.md).
