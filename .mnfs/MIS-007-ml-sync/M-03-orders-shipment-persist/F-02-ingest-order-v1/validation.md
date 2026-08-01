# F-02-ingest-order-v1 — validation

All Go commands run from `apps/server_core` with `GOCACHE`/`GOMODCACHE` exported as absolute
paths, each as a separate bash invocation, on Windows/Git Bash.

## go build ./...

```
$ export GOCACHE="$(pwd)/.gocache" GOMODCACHE="$(pwd)/.gomodcache" && go build ./...
(no output — clean)
```

## go vet ./... (plain and integration-tagged)

First pass failed:

```
vet.exe: internal\modules\orders\application\ingest_service_test.go:35:6: fakeBuyerFiscalReader redeclared in this block
```

Cause: `fakeBuyerFiscalReader` collided with an identically-named type already defined in
F-03's `enrich_service_test.go` (same package, `orders/application`). Fixed by renaming my
type to `fakeIngestBuyerFiscalReader` throughout `ingest_service_test.go` — F-03's file was
not touched.

Second pass failed differently:

```
vet.exe: internal\modules\orders\application\enrich_link_refresh_test.go:33:30: undefined: floatPtr
```

Cause: the pre-existing `import_service_test.go` owned three test helpers
(`floatPtr`, `intPtr`, `stubLinkReader`) that F-03's `enrich_service_test.go`/
`enrich_link_refresh_test.go` also call. My rewrite of `import_service_test.go` (needed for
the new `ImportServiceConfig{Source, Ingestor, Now}` shape) dropped them along with the old
`normalizeOrders`-era tests. Fixed by restoring all three in the same file, unmodified.

After both fixes:

```
$ export GOCACHE="$(pwd)/.gocache" GOMODCACHE="$(pwd)/.gomodcache" && go vet ./...
(no output — clean)
$ export GOCACHE="$(pwd)/.gocache" GOMODCACHE="$(pwd)/.gomodcache" && go vet -tags=integration ./...
(no output — clean)
```

## go test ./internal/modules/orders/... ./internal/composition/... -count=1

```
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/integrations	3.607s
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/internalread	3.679s
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/postgres	4.187s
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/pricingtax	3.745s
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/productlinks	3.704s
ok  	marketplace-central/apps/server_core/internal/modules/orders/application	4.022s
ok  	marketplace-central/apps/server_core/internal/modules/orders/domain	3.664s
ok  	marketplace-central/apps/server_core/internal/modules/orders/ports	3.796s
ok  	marketplace-central/apps/server_core/internal/modules/orders/transport	3.914s
ok  	marketplace-central/apps/server_core/internal/composition	4.763s
```

All packages pass. `orders/application` includes 12 new `ingest_service_test.go` tests
(`IngestService.IngestOrder`, fakes, no DB) and 6 rewritten `import_service_test.go` tests
(`ImportService.Import`, stub ingestor). Coverage against the brief's required scenarios:

- (a) truth-table bucket arms: `TestIngestOrderDerivesFaturarBucketWhenPaidAndNotFaturado`,
  `TestIngestOrderDerivesEnviarBucketWhenPaidAndFaturado`,
  `TestIngestOrderDerivesEnviadoBucketFromDeliveredTag`,
  `TestIngestOrderPersistsShipmentWhenFound` (shipped-status arm).
- (b) replay idempotency — DB-level, see the integration section below.
- (c) must-fail atomicity injection — DB-level, see the integration section below.
- (d) shipment-404-but-order-persists: `TestIngestOrderPersistsWithoutShipmentOnHonestAbsence`.
- (e) order-403-no-writes-at-all: `TestIngestOrderUnavailableOnOrderAuthErrorWritesNothing`,
  `TestIngestOrderUnavailableOnOrderNotFoundWritesNothing`.

Plus: `TestIngestOrderPropagatesShipmentRealErrorWithoutWriting` (a genuine, non-honest-absence
shipment error also aborts before any write), `TestIngestOrderPropagatesFaturadoLookupError`,
`TestIngestOrderAlwaysFetchesBuyerFiscal`, `TestIngestOrderRejectsBlankIdentifiers`, and
`TestImportServiceCountsAnyIngestFailureAsSkip` (the IC-06-driven design decision named in
spec.md).

## gofmt

This repo checks out with CRLF (`core.autocrlf=true`, `.gitattributes`) but `gofmt` normalizes
to LF, so a raw `gofmt -l` flags every touched file as needing formatting even when the
committed (LF) content is already clean — confirmed by running it against untouched files
too (`order_bucket.go`, `buyer_fiscal_reader.go`), which show the same false positive. Real
gofmt-cleanliness was checked by stripping `\r` first:

```
$ for f in <15 new/edited .go files>; do tr -d '\r' < "$f" > /tmp/lf_test.go; gofmt -l /tmp/lf_test.go; done
```

First pass found 3 real misalignments (interface-assertion block in `order_repo.go`, a struct
tag column in `order.go`, a struct field column in `ingest_service_test.go`) — all gofmt
artifacts of hand-edited column alignment, not logic. Fixed with targeted `Edit` calls.
Second pass: all 15 files clean (`ALL_CLEAN=1`).

## Integration tests (`//go:build integration`, real Postgres)

```
$ npm run harness:pg:up
container=mpc-pg-session-e7e297a5
port=64121
status=ready
$ npm run harness:integration   # first run
...
failure_token=package=marketplace-central/apps/server_core/tests/integration
failure_token=test=TestListingsReadContractEndToEnd
status=blocked
postgres lifecycle failed reasons=HPG_TEST_FAILED exit_code=1
```

`TestListingsReadContractEndToEnd` is a catalog/listings contract test in
`apps/server_core/tests/integration`, unrelated to this feature (orders module only) and not
touched by this diff. No `orders`-package failure token appeared. Re-ran to check for a flake:

```
$ npm run harness:integration   # second run
...
target=ephemeral-postgres
status=passed
```

Clean pass. The harness auto-discovers every package with `//go:build integration` in its
first 5 lines (`scripts/harness/Postgres.psm1`'s `Get-HarnessIntegrationTestPackages`), so
this run included `orders/adapters/postgres/ingest_repo_test.go` (new, 4 top-level tests):

- `TestOrderRepositoryIngestOrderPersistsHeaderItemsPaymentsAndShipment` — scenario (a):
  order+items+payments+order_shipments row all persist in one call; asserts the caller-supplied
  `Bucket` (`BucketEnviado`) and `ProviderShipmentID` round-trip through Postgres untouched.
  (Bucket DERIVATION correctness for the two truth-table arms is unit-tested DB-independently
  in `ingest_service_test.go` above — this test only proves the persisted VALUE survives the
  round-trip.)
- `TestOrderRepositoryIngestOrderPersistsOrderWithoutShipmentRow` — scenario (d) at the repo
  level: `shipment=nil` still persists the header with `provider_shipment_id` set, zero
  `order_shipments` rows.
- `TestOrderRepositoryIngestOrderIsIdempotentOnReplay` — scenario (b): identical order+shipment
  ingested twice; asserts `orders`/`items`/`payments`/`order_shipments` row counts are
  unchanged after the replay.
- `TestOrderRepositoryIngestOrderRollsBackWholeTransactionWhenShipmentWriteFails` — scenario
  (c), the must-fail atomicity test: ingests a baseline order, then a second `IngestOrder`
  call that mutates the header (`status: paid → shipped`, new bucket) AND supplies a shipment
  with a 4-character `currency` value (`order_shipments.currency` is `char(3)`, 0088) to force
  a Postgres data-length violation on `upsertOrderShipment`, which runs AFTER `upsertOrder` in
  the same transaction. Asserts: the call returns an error naming the violation
  (`value too long`), the order row's `status`/`bucket` are still the BASELINE values (not the
  mutated ones), and zero `order_shipments` rows exist. Names the property under test in its
  doc comment: if `IngestOrder` ever split the header write and the shipment write into two
  transactions, this test starts failing because the header would show `shipped` even though
  the shipment write failed.

```
$ npm run harness:pg:down
target=pg-session
status=stopped
```

## Scope discipline (`git status`/`git diff --stat`, from the worktree root)

```
 M apps/server_core/internal/composition/root.go
 M apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go
 M apps/server_core/internal/modules/orders/application/import_service.go
 M apps/server_core/internal/modules/orders/application/import_service_test.go
 M apps/server_core/internal/modules/orders/domain/order.go
 M apps/server_core/internal/modules/orders/ports/order_store.go
?? apps/server_core/internal/composition/orders_ingest_adapters.go
?? apps/server_core/internal/modules/orders/adapters/postgres/ingest_repo_test.go
?? apps/server_core/internal/modules/orders/application/ingest_service.go
?? apps/server_core/internal/modules/orders/application/ingest_service_test.go
?? apps/server_core/internal/modules/orders/domain/order_ingest_errors.go
?? apps/server_core/internal/modules/orders/domain/order_shipment.go
?? apps/server_core/internal/modules/orders/ports/order_detail_reader.go
?? apps/server_core/internal/modules/orders/ports/order_ingestor.go
?? apps/server_core/internal/modules/orders/ports/shipment_detail_reader.go
```

Only owned paths changed: `internal/composition/` (2 files), `orders/adapters/postgres/`
(2 files), `orders/application/` (4 files), `orders/domain/` (3 files), `orders/ports/`
(4 files). No connectors/, transport/, scheduler/ files touched.

`git diff -- apps/server_core/internal/composition/root.go` shows both hunks fall strictly
inside the original 578-603 line range: the moved `ordersBuyerFiscalReader` construction plus
the new `ordersOrderDetailReader`/`ordersShipmentDetailReader`/`ordersIngestSvc`
block, and the updated `ordersImportSvc` config. Every other line in the file (imports, the
enrich/tax/faturado wiring below, everything outside this func) is unchanged.

No git commit, push, or stage was run — working tree left for the orchestrator.
