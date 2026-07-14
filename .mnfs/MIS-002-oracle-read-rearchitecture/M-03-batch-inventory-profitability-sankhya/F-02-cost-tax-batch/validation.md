# F-02 Validation Evidence

Date: 2026-07-14
Base: `a714acf7aa87f4e93d1c44102cc51504a532c7dc`

PowerShell uses an absolute `GOCACHE` equivalent for the requested
`GOCACHE=.gocache` commands; all commands below ran from `apps/server_core`.

## M-03-C02 — profitability batch adoption

Command:

```text
GOCACHE=.gocache go test ./internal/modules/profitability/... -run 'Batch|QueryCount' -v
```

Key output:

```text
=== RUN   TestImportMarginInputsBatchQueryCountN1200
--- PASS: TestImportMarginInputsBatchQueryCountN1200 (0.01s)
=== RUN   TestBatchMissingFactsMakeProductIncomputable
--- PASS: TestBatchMissingFactsMakeProductIncomputable (0.00s)
PASS
ok   marketplace-central/apps/server_core/internal/modules/profitability/application 1.738s
```

The counting batch fake received 1200 deduplicated product IDs and observed
3 cost + 3 tax chunks. Nil cost/tax facts produced nil amounts and the item
snapshot had `missing_cost`, `missing_tax`, incomplete quality, and no
contribution amount.

Command:

```text
GOCACHE=.gocache go test ./internal/modules/internal_read/adapters/oracle -run 'BatchReader' -v
```

Key output:

```text
=== RUN   TestBatchReaderUses500ChunksAndNoPartialMapOnFailure
--- PASS: TestBatchReaderUses500ChunksAndNoPartialMapOnFailure (0.00s)
=== RUN   TestBatchReaderSalesHistoryPeeksAt5001AndMarksTruncated
--- PASS: TestBatchReaderSalesHistoryPeeksAt5001AndMarksTruncated (0.01s)
PASS
ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle 1.668s
```

The `database/sql/driver` fixture returned 5001 sales rows; the adapter
returned 5000 entries, `Truncated=true`, and emitted structured
`slow_query=true truncated=true row_cap=5000`. A forced second-chunk Oracle
failure returned typed `source_unavailable` with a nil result map.

## M-03-C06 — import ceiling

Command:

```text
GOCACHE=.gocache go test ./internal/modules/profitability/... -run 'Limit|Incomputable' -v
```

Key output:

```text
=== RUN   TestImportMarginInputsLimitExceededBeforeAnyPortCall
--- PASS: TestImportMarginInputsLimitExceededBeforeAnyPortCall (0.00s)
=== RUN   TestHandleImportLimitExceededReturns422AndNamesCap
--- PASS: TestHandleImportLimitExceededReturns422AndNamesCap (0.01s)
PASS
```

`Limit=200` is accepted by the batch-count test. `Limit=201` returns 422 with
`{"error":"limit_exceeded","limit":200,...}` before order, cost, tax, or
Oracle calls; the fake call counts remain zero.

## Full repository proof

Command:

```text
GOCACHE=.gocache go test ./...
```

Key output:

```text
ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle (cached)
ok   marketplace-central/apps/server_core/internal/modules/profitability/application (cached)
ok   marketplace-central/apps/server_core/internal/composition (cached)
ok   marketplace-central/apps/server_core/tests/integration 2.606s
ok   marketplace-central/apps/server_core/tests/unit 3.476s
```

No OpenAPI or `sdk-runtime` files were changed. The local
`internal_read/adapters/oracle/oraclebatch/stub.go` exists only to compile this
parallel worktree and is intentionally untracked/excluded from the commit;
F-01 supplies the canonical package.
