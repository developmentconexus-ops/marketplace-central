# F-02-ingest-order-v1 — plan

## Sequence executed

1. Domain layer: `orders/domain/order_shipment.go` (new `OrderShipment`),
   `orders/domain/order_ingest_errors.go` (new `ErrOrderUnavailable`),
   `orders/domain/order.go` extended with the 13 new 0089 fields.
2. Ports layer: `orders/ports/order_detail_reader.go`, `shipment_detail_reader.go`,
   `order_ingestor.go` (new); `orders/ports/order_store.go` extended with
   `OrderIngestStore`.
3. Composition adapters: `internal/composition/orders_ingest_adapters.go` (new file, per
   dispatch brief — `orders_adapters.go` itself untouched), wrapping
   `CapabilityAdapter.GetOrderDetail`/`GetShipmentDetail` behind the same
   `accountRefForInstallation` resolution the existing shipment/buyer-fiscal adapters use.
4. Repository layer: `orders/adapters/postgres/order_repo.go` — extended `upsertOrder`'s SQL
   (17 → 30 columns), added `upsertOrderShipment`, `GetFaturadoAt`, `IngestOrder` (single-tx
   orchestration reusing `upsertOrder`/`replaceItems`/`replacePayments`/
   `backfillMissingLines` verbatim for the header/items/payments branch, so a concurrent
   ingest of the same order from two triggers still can't regress a newer snapshot).
5. Application layer: `orders/application/ingest_service.go` (new `IngestService`) —
   fetch-then-write ordering (OrderDetail → ShipmentDetail if `ShippingID≠nil` → BuyerFiscal →
   links → `GetFaturadoAt` → `DeriveOrderBucket` → single `store.IngestOrder` call).
6. Refactor `orders/application/import_service.go`: `ImportServiceConfig{Source, Ingestor,
   Now}`, `Import` now enumerates via `ports.OrderSource.ListOrders` and calls
   `ports.OrderIngestor.IngestOrder` per id. Deleted `normalizeOrders`, `collectIdentities`,
   `trimNonEmpty` after a repo-wide grep confirmed zero external callers. `safeOrderProviderReference`/
   `keyOf` moved (not deleted) into `ingest_service.go`, which still needs them.
7. Wiring: `root.go` lines 578-603 — added `ordersOrderDetailReader`,
   `ordersShipmentDetailReader`, `ordersIngestSvc`, updated `ordersImportSvc`'s config;
   moved `ordersBuyerFiscalReader`'s construction earlier (still used unchanged by
   `ordersEnrichSvc` below it — F-03's adapter, not touched).
8. Verification: `go build`/`go vet` (plain and `-tags=integration`), package unit tests, new
   `//go:build integration` repo tests run against the harness's ephemeral Postgres session,
   `gofmt -l` (LF-normalized, see validation.md's CRLF note), `git status`/`git diff --stat`
   scope check.
9. This artifact set (`spec.md`/`plan.md`/`validation.md`).

## Correction made mid-implementation

`go vet` (step 8) failed after the `import_service.go`/`import_service_test.go` rewrite:
`enrich_link_refresh_test.go:33: undefined: floatPtr`. Root cause: the pre-existing
`import_service_test.go` owned three package-wide test helpers
(`floatPtr`, `intPtr`, `stubLinkReader`) that F-03's `enrich_service_test.go`/
`enrich_link_refresh_test.go` also depend on; the rewrite dropped them along with the
`normalizeOrders`-specific tests they used to support. Fixed by re-adding all three to the
rewritten `import_service_test.go` (same file, since it already owned them) without touching
either F-03 test file. Re-ran `go build`/`go vet`/`go test` clean after the fix — see
validation.md.

## Verification commands run (see validation.md for full output)

- `go build ./...`
- `go vet ./...` and `go vet -tags=integration ./...`
- `go test ./internal/modules/orders/... ./internal/composition/... -count=1`
- `gofmt -l` on all 15 new/edited `.go` files (CRLF-normalized first — see validation.md)
- `npm run harness:pg:up` then `npm run harness:integration` (twice — first run hit an
  unrelated pre-existing flake in `tests/integration`'s `TestListingsReadContractEndToEnd`,
  a catalog/listings contract test outside this feature's scope; second run passed clean,
  which includes the new `orders/adapters/postgres/ingest_repo_test.go` scenarios since the
  harness auto-discovers every `//go:build integration` package)
- `npm run harness:pg:down`
- `git status --porcelain` / `git diff --stat` / `git diff -- root.go` (scope discipline)

No git commit/push/stage was run — working tree left for the orchestrator.
