# F-03 Sankhya linkage lean — validation evidence

Validation ran from `apps/server_core`. On Windows, this Go toolchain requires
an absolute `GOCACHE`; each command below points at the workspace-local
`.gocache` directory.

## M-03-C03 — linkage reader lean

Command:

```powershell
$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/adapters/oracle/... -run 'Linkage|Redact|SafeCause' -v
```

Result: PASS. `TestSankhyaLinkageCandidateReadsChunkLinesWithoutRequestValidation`
passed for 1 and 700 candidates; the fake queryer observed 2 total queries for
1 candidate and 3 total queries for 700 candidates (one candidate query plus
one or two 500-ID IN-list line queries). It observed `pings=0`, and the test
rejected metadata/uniqueness queries on the request path. Existing typed
overflow, nullable descendant, identity, and ordering-related tests passed.

Command:

```powershell
$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/product_links/... -run 'Linkage' -v
```

Result: PASS. Product-links packages compiled and passed; the filter reported
no tests to run in those packages because the linkage fake-queryer proof is in
the internal-read Oracle adapter package.

Command:

```powershell
rg -n "ValidateConfiguration\(" 'internal/modules/internal_read/adapters/oracle/sankhya_linkage_reader.go' 'internal/composition/root.go'
```

Key output:

```text
internal/composition/root.go:396: if err := source.ValidateConfiguration(context.Background()); err != nil {
internal/modules/internal_read/adapters/oracle/sankhya_linkage_reader.go:51: func (r *SankhyaLinkageReader) ValidateConfiguration(ctx context.Context) error {
```

The request methods contain no validation call; enabled linkage validation is
performed once during composition-root startup and startup failure is returned
as an error.

Command:

```powershell
if (rg -n "FreshnessPolicy|cache|Cache" 'internal/modules/internal_read/adapters/oracle/sankhya_linkage_reader.go' 'internal/modules/product_links') { Write-Error 'unexpected linkage cache path'; exit 1 } else { Write-Output 'no linkage FreshnessPolicy/cache references' }
```

Key output: `no linkage FreshnessPolicy/cache references`.

## M-03-C04 — Oracle cause redaction

Command:

```powershell
$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./... -run 'Redact|SafeCause' -v
```

Result: PASS. `TestSankhyaLinkageOracleErrorsRedactDSNCauseAndLogs` passed. A
DSN-shaped fake driver error exposed neither username, host, service,
password, secret, nor raw `ORA-00942` text in the returned error or logs; the
safe cause was `oracle error code=942` and logs contained `oracle_code=942`.

## Full regression

Command:

```powershell
$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./...
```

Result: PASS across all server packages, including `internal/composition`,
internal-read Oracle, product-links, orders, inventory, profitability, and
integration/unit test packages.
