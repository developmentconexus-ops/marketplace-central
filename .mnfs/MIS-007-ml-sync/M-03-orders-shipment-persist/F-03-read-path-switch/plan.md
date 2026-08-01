# F-03-read-path-switch — plan

## Sequence executed

1. New Postgres readers: `orders/adapters/postgres/shipment_reader.go`
   (`ShipmentReader.GetShipment`, backs `ports.ShipmentReader`, reads `order_shipments`/0088)
   and `buyer_fiscal_reader.go` (`BuyerFiscalReader.GetBuyerFiscal`, backs
   `ports.BuyerFiscalReader`, reads `orders_marketplace_orders` buyer_* columns/0089). Both
   split their column-mapping logic into pure functions (`scanShipmentRow`,
   `buildBuyerFiscalInfo`) taking already-scanned `pgtype.*` values, so DB-free unit tests
   exercise every branch (honest-absence, partial data, cost-without-currency, etc.) without a
   database.
2. `ports.ErrShipmentNotFound` sentinel added to `orders/ports/shipment_reader.go`; wired into
   `GetShipment`'s no-row branch and consumed by `EnrichService.fetchShipment` via
   `errors.Is`.
3. root.go rewiring: construct `orderspostgres.NewShipmentReader`/`NewBuyerFiscalReader`,
   pass them into `NewEnrichServiceWithReaders`, delete the
   `ordersShipmentReader := newOrdersShipmentReaderAdapter(...)` line. Left
   `ordersBuyerFiscalReader := newOrdersBuyerFiscalReaderAdapter(...)` untouched — F-02's
   `ordersIngestSvc.BuyerFiscal` batch path still needs it.
4. Delete site A (`ordersShipmentReaderAdapter` struct + constructor + `GetShipment` method)
   from `composition/orders_adapters.go`. Left `ordersBuyerFiscalReaderAdapter`/
   `ordersCostReaderAdapter` untouched.
5. archguard: shrink `mlAllowlist` from 4 to 2 entries (drop the two orders sites), add 3
   entries to `mlExcludedSymbols` (order-detail reader, shipment-detail reader — both
   batch-only per F-02 — and the reclassified buyer-fiscal reader, now batch-only-only since
   its interactive use is gone). Update the doc comment above `mlAllowlist` to describe sites
   A (retired) and B (reclassified) accurately. Fix the now-broken
   `testdata/three_sites/root.go` fixture (hardcoded the two retired/reclassified symbols) to
   mirror only the current real wiring.
6. `GetOrderBucketCounts`: two-tier read (persisted `bucket` column, fallback to
   `DeriveOrderBucket` re-derivation). New integration test
   `order_repo_bucket_counts_test.go` seeds 3 rows via `UpsertOrders` proving a genuine
   before/after delta on one row (persisted "enviado" vs. what re-derivation with the old
   hardcoded empty shipment status would have produced, "faturar").
7. `sumSaleFee` bugfix in `enrich_service.go` (missing `×Quantity`), doc comment corrected,
   existing unit test coverage extended.
8. Golden `comprador_fiscal` proof split across the two packages that own the two unexported
   mapping functions: `postgres.TestBuildBuyerFiscalInfo_MatchesOldLiveAdapterShapeForEquivalentData`
   and `transport.TestMapCompradorFiscalGoldenOldAdapterShapeVsNewDBReaderShape`. Both hand-
   mirror the OLD live adapter's exact mapping (read from
   `connectors/adapters/mercado_livre/buyer_fiscal_reader.go`) for equivalent input data and
   assert JSON-byte-identical output against the NEW path.
9. Zero-live-call proof (`orders/adapters/postgres/zero_live_ml_call_test.go`): an import-scan
   test proving neither new reader file imports the mercado_livre connector adapter package
   (structural, compile-enforced — not just "doesn't currently call it"), plus a nil-pool
   degrade test and a pure-function smoke test. Combined with the already-passing
   `TestRealRepoInteractiveMLSites_MatchesAllowlist` (proves root.go's interactive wiring no
   longer references the retired symbol) as the composition-wiring half of the same claim.
10. `composition/orders_adapters_test.go` cleanup: removed the now-dangling
    `TestOrdersShipmentReaderAdapter` (drove the deleted constructor) and its single-use
    `erroringInstallationRepo` fixture; dropped the two imports that became unused as a
    result.

## Verification ladder

Run from `apps/server_core` with `GOCACHE`/`GOMODCACHE` exported as absolute paths:

```
go build ./...
go vet ./...
go vet -tags=integration ./...
go test ./internal/modules/orders/... ./internal/composition/... ./internal/platform/archguard/...
```

Plus the harness integration lane (`npm run harness:pg:up` then `npm run harness:integration`
from the repo root) for `order_repo_bucket_counts_test.go` (real Postgres, `//go:build
integration`). Plus a live must-fail/must-pass round-trip on the archguard test (temporarily
reintroduce a call shaped like the deleted site A, confirm RED naming the exact symbol/line,
revert, confirm GREEN). See validation.md for full output of every step.
