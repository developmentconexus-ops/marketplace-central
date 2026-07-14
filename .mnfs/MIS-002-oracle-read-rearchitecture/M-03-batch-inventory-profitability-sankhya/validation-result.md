# M-03 validation result — fixed-SHA review + proportional QA

Date: 2026-07-14  
Frozen SHA: `379305857b8e00fd288ea9a9429fafeedc82ac2a`  
Accepted base SHA: `a714acf7aa87f4e93d1c44102cc51504a532c7dc`  
Mode: `fixed_sha_review THEN proportional_qa`

## Fixed-SHA review — PASS

Review was read-only and restricted to the named M-03 seam. `HEAD` matched the
frozen SHA and `git diff --check` passed. M-02 catalog/OpenAPI/SDK files were
ignored.

- C01: **PASS** — stock and shared batch chunking is 500; missing facts remain
  nil with `missing_stock`, and no zero substitution is introduced.
- C02: **PASS** — cost/tax batch reads chunk at 500; sales peeks at 5001 and
  preserves `truncated`; missing facts remain incomputable.
- C03: **PASS** — the only executable `ValidateConfiguration` call is startup
  validation in `composition/root.go:404` using `context.Background()`. The
  orders wrapper is an in-memory guard and does not delegate. Request methods
  use candidate/descendant reads, preserving typed configuration/unavailable
  behavior. Candidate lines are read in 500-ID IN-list chunks; linkage has no
  cache path.
- C04: **PASS** — new Oracle error paths use `wrapOracleError` or
  `sankhyaOracleError` with `safeOracleCause`.
- C05: **PASS** — one shared semaphore is configured from
  `MPC_ORACLE_BATCH_PERMITS`, default 4, and is batch-only.
- C06: **PASS** — the 200 import cap is enforced before Oracle/port calls and
  maps to `limit_exceeded`.

Additional checks passed: canonical-identity behavior from `97fd4b58` remains
preserved; SQL/driver types remain in internal-read Oracle adapters; unknown
operational facts do not become zero/default values.

Static call-site check excerpt:

```text
apps/server_core/internal/composition/root.go:404: if err := source.ValidateConfiguration(context.Background()); err != nil {
apps/server_core/internal/modules/orders/adapters/internalread/sankhya_linkage_reader.go:39: func (r *SankhyaLinkageReader) ValidateConfiguration(_ context.Context) error {
```

## Proportional QA

All commands ran from `apps/server_core` with an absolute local GOCACHE:
`$env:GOCACHE = (Join-Path (Get-Location) '.gocache')`.

### C01 — PASS

```text
go test ./internal/modules/inventory/... ./internal/modules/internal_read/... -run 'StockBatch' -v
```

Relevant output: `TestStockBatchQueriesExpectedChunkCount` passed for
`one`, `exact_chunk`, `one_over_chunk`, and `three_chunks`; observed query
counts were `1/1/2/3` for `N=1/500/501/1200`. Empty input made zero queries;
chunk failure returned no partial map. `TestStockBatchSemaphoreBoundsConcurrency`
and `TestStockBatchSemaphoreDoesNotBlockInteractiveClass` both passed.

### C02 — PASS

```text
go test ./internal/modules/profitability/... -run 'Batch|QueryCount' -v
go test ./internal/modules/internal_read/adapters/oracle/... -run 'BatchReader' -v
```

Relevant output: `TestImportMarginInputsBatchQueryCountN1200` and
`TestBatchMissingFactsMakeProductIncomputable` passed; the fake observed 3
cost and 3 tax chunks. `TestBatchReaderUses500ChunksAndNoPartialMapOnFailure`
and `TestBatchReaderSalesHistoryPeeksAt5001AndMarksTruncated` passed; sales
returned 5000 rows with `truncated=true` after the 5001-row peek.

### C03 — PASS

```text
go test ./internal/modules/internal_read/adapters/oracle/... -run 'Linkage' -v
go test ./internal/modules/orders/application/... ./internal/modules/orders/adapters/internalread/... -run 'Linkage|Sankhya' -v
```

Relevant output: `TestSankhyaLinkageCandidateReadsChunkLinesWithoutRequestValidation`
passed for 1 and 700 candidates; the linkage suite passed the redaction,
no-cache, exact-identity, and typed-state tests. Orders application and
internal-read adapter linkage suites passed, including unavailable/configuration
state coverage.

Observation: the contract names `./internal/modules/product_links/... -run
'Linkage'`, but that location contains no linkage tests. The packet-directed
Oracle and orders commands above are the substantive C03 evidence.

### C04 — PASS

```text
go test ./... -run 'Redact|SafeCause' -v
```

Relevant output: `TestSankhyaLinkageOracleErrorsRedactDSNCauseAndLogs` passed;
the forced DSN-shaped error exposed only numeric Oracle code and no username,
host, service, or raw driver cause.

One initial invocation from the repository root failed before test setup because
the command was outside the packet-required module directory. The exact command
was rerun from `apps/server_core` and passed as recorded above.

### C05 — PASS

Covered by the C01 semaphore tests: maximum in-flight was 4 and interactive
work completed while batch permits were exhausted.

### C06 — PASS

```text
go test ./internal/modules/profitability/... -run 'Limit' -v
```

Relevant output: `TestImportMarginInputsLimitExceededBeforeAnyPortCall` and
`TestHandleImportLimitExceededReturns422AndNamesCap` passed. `Limit=201`
returned HTTP 422 with `limit_exceeded` and the 200 cap; `Limit=200` remained
accepted.

### Full — PASS

```text
go test ./...
```

Exit code `0`; all server packages, including composition, internal-read
Oracle/oraclebatch, inventory, orders, profitability, integration, and unit
tests passed.

## Overall verdict

**PASS.** C01–C06 passed, and the full registered repository gate passed.

## Blockers and next

Blockers: none. Next: accept M-03 at frozen SHA
`379305857b8e00fd288ea9a9429fafeedc82ac2a`.
