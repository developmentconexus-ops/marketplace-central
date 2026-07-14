# F-01 stock-batch validation

Date: 2026-07-14
Base SHA: `a714acf7aa87f4e93d1c44102cc51504a532c7dc`

## M-03-C01 — stock batch chunks correctly

Requested command:

```text
GOCACHE=.gocache go test ./internal/modules/inventory/... ./internal/modules/internal_read/... -run 'StockBatch' -v
```

PowerShell with the repository Go toolchain rejected the relative cache path
before running tests:

```text
build cache is required, but could not be located: GOCACHE is not an absolute path
```

Passing equivalent command, from `apps/server_core`:

```text
$env:GOCACHE=(Resolve-Path .gocache).Path; go test ./internal/modules/inventory/... ./internal/modules/internal_read/... -run 'StockBatch' -v
```

Key output:

```text
--- PASS: TestStockBatchQueriesExpectedChunkCount/one
--- PASS: TestStockBatchQueriesExpectedChunkCount/exact_chunk
--- PASS: TestStockBatchQueriesExpectedChunkCount/one_over_chunk
--- PASS: TestStockBatchQueriesExpectedChunkCount/three_chunks
--- PASS: TestStockBatchDeduplicatesAndEmptyInputDoesNotQuery
--- PASS: TestStockBatchChunkFailureReturnsNoPartialMap
--- PASS: TestStockBatchAdapterSemaphoreBoundsInstrumentedFake
ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle
```

The instrumented fake observed exactly 1, 1, 2, and 3 queries for 1, 500,
501, and 1200 IDs. Empty input made zero queries; duplicate IDs made one
query; the second-chunk failure returned a typed source-unavailable error and
nil map. Stock-risk service proof made one batch call, preserved nil internal
quantity, and emitted `missing_stock` for an absent fact.

## M-03-C05 — batch semaphore bounds concurrency

Requested command:

```text
GOCACHE=.gocache go test ./internal/modules/inventory/... ./internal/modules/internal_read/... -run 'StockBatch|Semaphore' -v
```

Passing equivalent command used for the same absolute-cache reason:

```text
$env:GOCACHE=(Resolve-Path .gocache).Path; go test ./internal/modules/inventory/... ./internal/modules/internal_read/... -run 'StockBatch|Semaphore' -v
```

Key output:

```text
--- PASS: TestStockBatchAdapterSemaphoreBoundsInstrumentedFake
--- PASS: TestStockBatchSemaphoreBoundsConcurrency
--- PASS: TestStockBatchSemaphoreDoesNotBlockInteractiveClass
ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle/oraclebatch
```

Eight concurrent adapter calls against the instrumented fake observed
`max in-flight = 4`. The interactive-class test completed while a batch permit
was exhausted, proving it does not use the batch semaphore.

## Full repository gate

Requested command:

```text
GOCACHE=.gocache go test ./...
```

Passing equivalent command:

```text
$env:GOCACHE=(Resolve-Path .gocache).Path; go test ./...
```

Key output:

```text
ok   marketplace-central/apps/server_core/internal/composition
ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle
ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle/oraclebatch
ok   marketplace-central/apps/server_core/internal/modules/inventory/application
ok   marketplace-central/apps/server_core/tests/integration
```

Exit code: `0`.
