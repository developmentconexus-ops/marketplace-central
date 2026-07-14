# F-01 stock-batch specification

## Contract

`oraclebatch.Chunks` preserves input order and splits into chunks no larger
than the requested maximum. `oraclebatch.Semaphore` bounds only batch-class
calls; the composition root creates one instance from
`MPC_ORACLE_BATCH_PERMITS`, defaulting to four permits.

The internal-read stock batch port is:

```go
GetStockFactsByIDs(ctx context.Context, ids []int64) (map[int64]*StockFact, error)
```

The Oracle adapter deduplicates IDs, returns immediately for an empty set,
chunks at 500 IDs, and runs those chunks sequentially while holding one batch
permit. Results are keyed by product ID. A product absent from Oracle is
absent from the map; a present row with an unknown quantity keeps a nil
quantity and `missing_stock`. Any query, scan, or iteration failure returns a
wrapped typed Oracle error and no partial map.

Stock-risk listing collects all resolved internal product IDs, performs one
batch-port call per run, and classifies every snapshot from the returned map.
Missing IDs produce nil internal quantity and the `missing_stock` quality
flag, so unknown stock never becomes zero.

## Acceptance evidence

- Query-count tests prove 1, 1, 2, and 3 Oracle calls for 1, 500, 501, and
  1200 IDs respectively.
- Empty input performs zero Oracle calls; duplicate IDs are queried once;
  chunk failure returns no partial map.
- Eight concurrent batch calls observe at most four in-flight operations, and
  an interactive operation does not acquire the batch semaphore.
- Stock-risk tests prove one batch call per run and preserve nil stock plus
  `missing_stock` for absent facts.
- Targeted and full Go test commands are recorded in `validation.md`.
