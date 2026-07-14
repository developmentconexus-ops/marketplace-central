# F-02 Cost/Tax Batch Specification

## Contract

- Cost and tax reads expose map-keyed batch ports keyed by internal product ID.
- IDs are deduplicated, split into chunks of at most 500, executed sequentially,
  and merged only after every chunk succeeds. A failed chunk returns a typed,
  safe Oracle error and no partial map.
- Every batch call acquires the shared Oracle batch semaphore for the complete
  batch operation. Missing facts remain nil and carry `missing_cost` or
  `missing_tax`; no numeric zero is substituted.
- Profitability gathers resolved product IDs once per import and performs one
  cost batch and one tax batch. Existing lineage-aware per-line tax behavior is
  retained for non-batch implementations and test doubles; the production
  batch implementation supplies product-keyed tax facts without raw provider
  payloads crossing the adapter boundary.
- `ImportMarginInputs.Limit > 200` fails before order or Oracle port calls with
  a 422 `limit_exceeded` response naming the 200 cap. `Limit == 200` remains
  valid.
- Sales history uses a 5001-row peek with a 5000-row result cap. The domain
  result carries `Truncated=true` when the peek row is present, and the batch
  adapter emits a structured warning. No OpenAPI or SDK response change is
  made because this baseline has no documented profitability sales-report
  response to extend.

## Acceptance proof

- Batch query counts are ceil(N/500), including N=1200, with duplicate IDs not
  increasing the count.
- A missing cost or tax fact produces nil amount and the matching quality flag;
  snapshot finalization marks the product incomplete/incomputable.
- 5001 sales fixture rows produce 5000 returned entries and `Truncated=true`.
- Import limits 200 and 201 prove acceptance and zero port calls respectively.
