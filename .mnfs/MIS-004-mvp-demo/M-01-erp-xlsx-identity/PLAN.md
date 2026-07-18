# M-01 implementation plan — ERP XLSX identity

Planning base: `59d0e62fdbf15db068542432ef5d5731b6fa9f83`. This plan covers `F-01 ∥ F-02 → F-03`; it is not an implementation record. The binding semantics come from IC-01/IC-02 and ADR-02/03/12/17. In particular, current code and tests that treat `TGFPRO.REFERENCIA` as manufacturer reference are defects: ADR-03 requires checksum-valid `REFERENCIA` to become EAN and IC-01 requires `REFFORN` alone to become manufacturer reference (`.mnfs/MIS-004-mvp-demo/mission.md:86-87`, `.mnfs/MIS-004-mvp-demo/research/identity-matching-interface-contract.md:18-26`).

Command convention below: run PowerShell from repository root. Every Go command sets a checkout-local cache and changes to `apps/server_core`; it does not use a shared cache. Each behavioral slice starts by adding the named failing test, observing the intended failure, and only then adding production code.

## A. Per-feature slice cards

### F-01 — identity semantics fix

#### CARD F01-S1

- `id`: `F01-S1`
- `title`: Publish the canonical pure GTIN validator
- `feature`: `F-01`
- `depends_on`: `[]`
- `complexity`: `standard`
- `validation_kind`: `unit`
- `write_set`:
  - `apps/server_core/internal/modules/catalog/domain/gtin.go` (new)
  - `apps/server_core/internal/modules/catalog/domain/gtin_test.go` (new; failing first)
- `commands`:
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/modules/catalog/domain -run TestIsValidGTIN; Pop-Location`
- `expected_artifacts`:
  - Table test with valid GTIN-8/UPC-A/GTIN-13/GTIN-14 and invalid checksum, invalid 12-digit value, empty/whitespace, wrong length, and non-digits.
  - Published compile-time seam: `catalog/domain.IsValidGTIN`.
- `open_questions`: `[]`

#### CARD F01-S2

- `id`: `F01-S2`
- `title`: Correct Oracle identity lookup and direct-reader projection
- `feature`: `F-01`
- `depends_on`: `[F01-S1]`
- `complexity`: `complex`
- `validation_kind`: `contract`
- `write_set`:
  - `apps/server_core/internal/modules/internal_read/domain/internal_product.go`
  - `apps/server_core/internal/modules/internal_read/domain/internal_tax.go`
  - `apps/server_core/internal/modules/internal_read/domain/quality_flag.go`
  - `apps/server_core/internal/modules/internal_read/domain/quality_flag_test.go`
  - `apps/server_core/internal/modules/internal_read/domain/contract_test.go`
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go`
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/reader_test.go`
- `commands`:
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/modules/internal_read/domain ./internal/modules/internal_read/adapters/oracle -run 'Test.*(Identity|Product|GTIN|Collision|EAN|Reference|NCM)'; Pop-Location`
- `expected_artifacts`:
  - Failing-first Oracle fixtures for valid/invalid/blank `REFERENCIA`, independent present/absent `REFFORN`, `MARCA`, `NCM`, EAN-only lookup, and two CODPRODs sharing a valid EAN.
  - `ProductCandidate` gains nullable `NCM`; `ReferenceCode` remains the compatibility Go name but is populated only from `REFFORN`. `TaxInputs` gains nullable `NCM` without changing any Reader method signature.
  - New identity flags `invalid_ean` and `ean_collision`; collision detection must query the full active-product population (not merely count rows in a filtered one-product result).
  - Oracle SELECT includes `REFERENCIA`, `REFFORN`, `MARCA`, `CODMARCA`, and `NCM`; EAN lookup matches `REFERENCIA`, then validates it in Go. Invalid/blank `REFERENCIA` is discarded, never copied to `ReferenceCode`.
- `open_questions`: `[]`

#### CARD F01-S3

- `id`: `F01-S3`
- `title`: Carry canonical identity through Oracle catalog pages and catalog HTTP DTOs
- `feature`: `F-01`
- `depends_on`: `[F01-S2]`
- `complexity`: `complex`
- `validation_kind`: `contract`
- `write_set`:
  - `apps/server_core/internal/modules/internal_read/ports/catalog_page.go`
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/catalog_page.go`
  - `apps/server_core/internal/modules/internal_read/adapters/oracle/catalog_page_test.go`
  - `apps/server_core/internal/modules/catalog/domain/canonical_product.go`
  - `apps/server_core/internal/modules/catalog/domain/canonical_product_test.go`
  - `apps/server_core/internal/modules/catalog/adapters/internalread/reader.go`
  - `apps/server_core/internal/modules/catalog/adapters/internalread/reader_test.go`
  - `apps/server_core/internal/modules/catalog/transport/http_handler.go`
  - `apps/server_core/internal/modules/catalog/transport/http_handler_test.go`
- `commands`:
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/modules/internal_read/adapters/oracle ./internal/modules/catalog/... -run 'Test.*(Catalog|Identity|EAN|Reference|NCM|Collision)'; Pop-Location`
- `expected_artifacts`:
  - Contract tests flip the current `CAST(NULL ... AS EAN)` behavior at `internal_read/adapters/oracle/catalog_page.go:124-133` and assert null-not-empty semantics.
  - `CatalogProductFact` and `CanonicalProduct` carry `EAN`, manufacturer reference from `REFFORN`, brand from `MARCA`, `NCM`, and identity quality flags. The transport retains the existing `reference` field as a compatibility alias of `REFFORN` and adds the governed names; no raw `REFERENCIA` escapes the Oracle adapter.
  - Page and direct-product responses expose `invalid_ean`/`ean_collision`, preserve both colliding CODPRODs, and never synthesize empty strings.
- `open_questions`: `[]`

#### CARD F01-S4

- `id`: `F01-S4`
- `title`: Atomically align catalog OpenAPI and the manual SDK types
- `feature`: `F-01`
- `depends_on`: `[F01-S3]`
- `complexity`: `standard`
- `validation_kind`: `openapi-sdk`
- `write_set`:
  - `contracts/api/marketplace-central.openapi.yaml` (only catalog schemas named in §C)
  - `packages/sdk-runtime/src/index.ts` (only additive catalog-type grant at current lines 164-204)
  - `packages/sdk-runtime/src/index.test.ts`
- `commands`:
  - `npm run build --workspace @marketplace-central/sdk-runtime`
  - `npm run test --workspace @marketplace-central/sdk-runtime`
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/modules/catalog/transport; Pop-Location`
- `expected_artifacts`:
  - One atomic OpenAPI+SDK commit with parity tests for nullable `ean`, `manufacturer_reference`, `brand_name`, `ncm`, and identity flags.
  - No generated output: `packages/sdk-runtime/package.json:7-10` has only manual TypeScript build/test scripts and ADR-12 explicitly declares the SDK manual (`.mnfs/MIS-004-mvp-demo/mission.md:96`).
- `open_questions`: `[]`

### F-02 — ERP import module

#### CARD F02-S1

- `id`: `F02-S1`
- `title`: Define import domain, parser boundary, exact validation codes, and normalized rows
- `feature`: `F-02`
- `depends_on`: `[F01-S1]`
- `complexity`: `standard`
- `validation_kind`: `unit`
- `write_set`:
  - `apps/server_core/internal/modules/erp_import/domain/import.go` (new)
  - `apps/server_core/internal/modules/erp_import/domain/validation.go` (new)
  - `apps/server_core/internal/modules/erp_import/domain/validation_test.go` (new; failing first)
  - `apps/server_core/internal/modules/erp_import/ports/parser.go` (new)
  - `apps/server_core/internal/modules/erp_import/ports/repository.go` (new)
- `commands`:
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/modules/erp_import/domain; Pop-Location`
- `expected_artifacts`:
  - Normalized domain row uses pointers for every optional field and decimal text/validated numeric representation that cannot turn unknown into zero.
  - Exact row errors `EMPTY_CODPROD`, `DUPLICATE_CODPROD`, `EMPTY_DESCRPROD`, `INVALID_CUSTO`, `INVALID_ESTOQUE`; warnings `INVALID_EAN`, `INVALID_NCM`; file errors are typed separately.
  - EAN validation calls the fixed `catalogdomain.IsValidGTIN(string) bool`; workbook/excelize types never cross `ports.Parser`.
- `open_questions`: `[]`

#### CARD F02-S2

- `id`: `F02-S2`
- `title`: Create tenant-scoped snapshot schema in the granted migration block
- `feature`: `F-02`
- `depends_on`: `[]`
- `complexity`: `standard`
- `validation_kind`: `migration`
- `write_set`:
  - `apps/server_core/migrations/0045_create_erp_import_protocols.sql` (new)
  - `apps/server_core/migrations/0046_create_erp_import_products.sql` (new)
  - `apps/server_core/migrations/0047_create_erp_import_issues.sql` (new; includes latest-snapshot, tenant/hash uniqueness, and reader indexes)
  - `apps/server_core/internal/platform/migrate/runner_test.go`
- `commands`:
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/platform/migrate; Pop-Location`
  - `Get-ChildItem apps/server_core/migrations/004[5-7]*.sql | Select-Object -ExpandProperty Name`
- `expected_artifacts`:
  - `erp_import_protocols`: UUID/id, `tenant_id`, unique tenant/file SHA-256, sanitized filename, stored `#NNN-E` protocol, source `xlsx`, UTC `imported_at`, `COMPLETED|REJECTED`, accepted/rejected/warning counts.
  - `erp_import_products`: `(tenant_id, protocol_id, codprod)` key plus descrprod, exact custo, physical/reserved stock, nullable ean/refforn/marca/ncm.
  - `erp_import_issues`: tenant/protocol/row/nullable column/kind (`REJECTION|WARNING`)/exact code/detail/nullable offending value. No workbook blob or raw row is persisted.
  - Runner fixture moves atomically from 41 to 44 at both assertions currently at `runner_test.go:25-26,64-65`.
- `open_questions`: `[]`

#### CARD F02-S3

- `id`: `F02-S3`
- `title`: Parse first-sheet XLSX behind the parser adapter
- `feature`: `F-02`
- `depends_on`: `[F02-S1]`
- `complexity`: `standard`
- `validation_kind`: `integration`
- `write_set`:
  - `apps/server_core/go.mod`
  - `apps/server_core/go.sum`
  - `apps/server_core/internal/modules/erp_import/adapters/xlsx/parser.go` (new)
  - `apps/server_core/internal/modules/erp_import/adapters/xlsx/parser_test.go` (new; failing first)
  - `apps/server_core/internal/modules/erp_import/adapters/xlsx/testdata/valid-and-invalid.xlsx` (new, non-sensitive)
  - `apps/server_core/internal/modules/erp_import/adapters/xlsx/testdata/missing-custo.xlsx` (new, non-sensitive)
  - `apps/server_core/internal/modules/erp_import/adapters/xlsx/testdata/empty.xlsx` (new, non-sensitive)
- `commands`:
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/modules/erp_import/adapters/xlsx; Pop-Location`
- `expected_artifacts`:
  - Header matching is case-insensitive, trimmed, and accent-insensitive; only first sheet is read.
  - Corrupt/empty/non-XLSX becomes `INVALID_FILE`; a missing required header names the column and becomes `MISSING_REQUIRED_COLUMN`.
  - The adapter closes workbook/file resources and returns normalized rows only; no raw cell values are logged.
- `open_questions`:
  - `BLOCKED — hub REQUEST dep-grant`: approve `github.com/xuri/excelize/v2` and the resulting `go.mod`/`go.sum` edits. Resolution: hub records the dependency grant (or names an already-approved equivalent); then this card becomes dispatch-ready. The dependency is currently absent at `apps/server_core/go.mod:5-22`.

#### CARD F02-S4

- `id`: `F02-S4`
- `title`: Implement atomic protocol, duplicate/concurrency, snapshot, and report persistence
- `feature`: `F-02`
- `depends_on`: `[F02-S1, F02-S2]`
- `complexity`: `complex`
- `validation_kind`: `integration`
- `write_set`:
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/import_repository.go` (new)
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/import_repository_test.go` (new; failing first)
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go` (new)
  - `apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository_test.go` (new; failing first)
- `commands`:
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test -tags=integration ./internal/modules/erp_import/adapters/postgres; Pop-Location`
- `expected_artifacts`:
  - A transaction/session API uses `pg_try_advisory_xact_lock` keyed by tenant for fail-fast `IMPORT_IN_PROGRESS`, rechecks `(tenant_id,file_hash)` under the lock, and returns the original id/protocol for `DUPLICATE_FILE`.
  - Protocol, products, issues, and final status commit atomically; all-rejected persists protocol/report but no products. Failed parsing/header validation creates no completed snapshot.
  - Every SELECT/INSERT/UPDATE/DELETE includes `tenant_id`; integration fixtures create a second tenant and prove no cross-tenant duplicate, list, detail, or snapshot leakage.
  - List is newest-first; detail returns exact row/column/code/detail/offending value; latest Reader query ignores `REJECTED` snapshots.
- `open_questions`: `[]`

#### CARD F02-S5

- `id`: `F02-S5`
- `title`: Orchestrate synchronous full-snapshot imports and honest partial rejection
- `feature`: `F-02`
- `depends_on`: `[F02-S3, F02-S4]`
- `complexity`: `complex`
- `validation_kind`: `unit`
- `write_set`:
  - `apps/server_core/internal/modules/erp_import/application/import_service.go` (new)
  - `apps/server_core/internal/modules/erp_import/application/import_service_test.go` (new; failing first)
  - `apps/server_core/internal/modules/erp_import/application/query_service.go` (new)
  - `apps/server_core/internal/modules/erp_import/application/query_service_test.go` (new; failing first)
- `commands`:
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/modules/erp_import/application; Pop-Location`
- `expected_artifacts`:
  - SHA-256 is computed while bounded upload content is spooled; filename is sanitized; raw workbook bytes and row payloads are never logged or exposed beyond the parser adapter.
  - Valid rows plus per-row rejections/warnings produce `COMPLETED`; zero valid rows produce `REJECTED`; missing header/invalid file persist no products; duplicate and concurrent errors retain typed details.
  - Protocol formatting is deterministic `^#[0-9]{3,}-E$`; source/import time are explicit, and no blanket recover/fallback converts parse or integrity failures into success.
- `open_questions`: `[]`

#### CARD F02-S6

- `id`: `F02-S6`
- `title`: Expose multipart POST and tenant-scoped list/detail handlers
- `feature`: `F-02`
- `depends_on`: `[F02-S5]`
- `complexity`: `standard`
- `validation_kind`: `contract`
- `write_set`:
  - `apps/server_core/internal/modules/erp_import/transport/http_handler.go` (new)
  - `apps/server_core/internal/modules/erp_import/transport/http_handler_test.go` (new; failing first)
- `commands`:
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/modules/erp_import/transport; Pop-Location`
- `expected_artifacts`:
  - `POST /erp/imports` uses `http.MaxBytesReader`, `ParseMultipartForm`, and `FormFile("file")`, processes synchronously, and returns 201 `{import_id,protocol,status}`.
  - `GET /erp/imports` and `GET /erp/imports/{id}` return list/detail; exact 400/409/422 mappings include duplicate original id/protocol.
  - No new auth middleware. The present root only wraps routes with `httpx.CORSMiddleware` (`internal/composition/root.go:659`; `platform/httpx/router.go:16-28`) and injects `cfg.DefaultTenantID` into repositories (for example `root.go:459`); F03 performs the same constructor injection, while F02 services never accept caller-supplied tenant IDs.
- `open_questions`: `[]`

#### CARD F02-S7

- `id`: `F02-S7`
- `title`: Atomically add ERP import OpenAPI paths and manual SDK client
- `feature`: `F-02`
- `depends_on`: `[F01-S4, F02-S6]`
- `complexity`: `standard`
- `validation_kind`: `openapi-sdk`
- `write_set`:
  - `contracts/api/marketplace-central.openapi.yaml` (only `/erp/imports*` and `ErpImport*` components in §C)
  - `packages/sdk-runtime/src/erpImport.ts` (new)
  - `packages/sdk-runtime/src/erpImport.test.ts` (new)
- `commands`:
  - `npm run build --workspace @marketplace-central/sdk-runtime`
  - `npm run test --workspace @marketplace-central/sdk-runtime`
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/modules/erp_import/transport; Pop-Location`
- `expected_artifacts`:
  - One atomic OpenAPI+SDK commit made on top of F01-S4; typed client sends `FormData` and models all success/error/report fields.
  - The one-line `export * from "./erpImport"` is intentionally absent from this card and is handed to the hub seam lock.
- `open_questions`: `[]`

### F-03 — reader adapter and source selection

#### CARD F03-S1

- `id`: `F03-S1`
- `title`: Implement the full Reader plus catalog-page adapter over the latest completed snapshot
- `feature`: `F-03`
- `depends_on`: `[F01-S3, F02-S4, F02-S5]`
- `complexity`: `complex`
- `validation_kind`: `integration`
- `write_set`:
  - `apps/server_core/internal/modules/erp_import/adapters/internalread/reader.go` (new)
  - `apps/server_core/internal/modules/erp_import/adapters/internalread/reader_test.go` (new; failing first)
  - `apps/server_core/internal/modules/erp_import/adapters/internalread/reader_integration_test.go` (new; failing first)
- `commands`:
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/modules/erp_import/adapters/internalread; go test -tags=integration ./internal/modules/erp_import/adapters/internalread; Pop-Location`
- `expected_artifacts`:
  - Compile assertions for `internal_read/ports.Reader` (all six methods at `ports/reader.go:48-55`) and `CatalogPageReader` so imported products are visible through `/catalog`.
  - Exported package sentinel `ErrNoErpSnapshot = errors.New("no_erp_snapshot")`; absence/as-of mismatch returns `domain.NewReadError(domain.ReadErrorSourceUnavailable, "no ERP snapshot", ErrNoErpSnapshot)`. This supports both `errors.Is` and existing `IsReadErrorCode` without changing the frozen Reader interface.
  - Exported typed `ERPProductNotFoundError` for a missing CODPROD; no zero-value success.
  - Every query includes tenant and `status='COMPLETED'`; a newer `REJECTED` protocol is ignored.
  - Exact per-method behavior is frozen in §C.
- `open_questions`: `[]`

#### CARD F03-S2

- `id`: `F03-S2`
- `title`: Select xlsx/oracle behind cache/timing decorators, register routes, and register governance
- `feature`: `F-03`
- `depends_on`: `[F02-S6, F03-S1]`
- `complexity`: `standard`
- `validation_kind`: `contract`
- `write_set`:
  - `apps/server_core/internal/composition/root.go` (sole writer; additive imports/config/registration/source override only)
  - `apps/server_core/internal/composition/root_test.go`
  - `contracts/governance/modules.json` (one `erp_import` entry only)
- `commands`:
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/composition; Pop-Location`
  - `npm run harness:governance` — REVIEW-BINDING: the governance lane false-fails inside a chip worktree (it scans `.claude/worktrees`). Run it from a CLEAN detached worktree at the full 40-hex BaseSha; the `modules.json` `erp_import` entry is validated at merge, never asserted green naively in-chip. In-chip, only lint the entry shape.
- `expected_artifacts`:
  - Pure `erpSource(getenv func(string) string) (string, error)` accepts empty/`oracle`/`xlsx`, defaults empty to `oracle`, and rejects all other values at startup; it reads `MC_ERP_SOURCE` with the existing `os.Getenv` injection idiom (`root.go:253-258,662-676`).
  - For xlsx, construct the repository/import handlers and snapshot Reader, then wrap it `cache.NewReader` followed by `observability.NewTimingReader`, matching current order at `root.go:401-410`; set the same service used by product linking and catalog pages. For oracle/default, retain the current construction unchanged.
  - Add `/erp/imports` to batch route classification and register the handler. Repository receives `cfg.DefaultTenantID`; no request body/query can select a tenant.
  - Governance entry: `{"id":"erp_import","root":"apps/server_core/internal/modules/erp_import","code_owner_path":"apps/server_core/internal/modules/erp_import","composition_required":true,"openapi_prefixes":["/erp/imports"],"dependencies":["catalog","internal_read"]}` in sorted module order.
- `open_questions`: `[]`

#### CARD F03-S3

- `id`: `F03-S3`
- `title`: Prove xlsx/Oracle substitutability and complete milestone regression lanes
- `feature`: `F-03`
- `depends_on`: `[F01-S4, F02-S7, F03-S2]`
- `complexity`: `standard`
- `validation_kind`: `contract`
- `write_set`:
  - `apps/server_core/internal/modules/erp_import/adapters/internalread/source_contract_test.go` (new; failing first)
  - `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/fixtures/example-erp.xlsx` (new, operator-approved non-sensitive fixture)
  - `.mnfs/MIS-004-mvp-demo/M-01-erp-xlsx-identity/fixtures/identity-rejections.xlsx` (new, non-sensitive fixture)
- `commands`:
  - `Push-Location apps/server_core; $env:GOCACHE="$PWD\.gocache"; go test ./internal/modules/erp_import/adapters/internalread ./internal/modules/internal_read/adapters/oracle ./internal/modules/catalog/...; go build ./cmd/server; Pop-Location`
  - REVIEW-BINDING (chip non-negotiable): `npm run harness:live` is **NOT** an implementer command — the chip never loads `.env*`/Oracle creds nor boots a server. The live-Oracle lane is QA-only at milestone close (P7). The implementer proves substitutability against the hermetic Oracle-shaped fake + `go build ./cmd/server` only.
- `expected_artifacts`:
  - One contract scenario exercised against a hermetic Oracle-shaped fake and the Postgres snapshot adapter with identical supported shapes and typed unsupported/unavailable outcomes.
  - Live-Oracle transcript is recorded later by QA in `validation-result.md §oracle`; mocks only prove shape, never live integration.
  - Versioned ≥50-row import fixture and the exact invalid-row fixture required by M-01-C01/C02.
- `open_questions`: `[]`

## B. Write-DAG and serialization

```text
F01-S1 ─┬─> F01-S2 -> F01-S3 -> F01-S4 ───────────────┐
         └─> F02-S1 -> F02-S3(dep grant) ─┐            │
F02-S2 ───────────────> F02-S4 ───────────┼-> F02-S5 -> F02-S6 -> F02-S7
                                          │                         │
F01-S3 + F02-S4 + F02-S5 ───────────────> F03-S1 -> F03-S2 ───────┴-> F03-S3
```

- First landing is **F01-S1** because F02 compiles against its exact signature. Once it lands, F01-S2, F02-S1, and F02-S2 may be in flight simultaneously; they have disjoint write sets.
- F02-S3 may start only after F02-S1 and the hub dependency grant. F02-S4 can run in parallel with F02-S3. F02-S5 joins them.
- `contracts/api/marketplace-central.openapi.yaml` is serialized: F01-S4 lands first; only then does F02-S7 edit a disjoint region and manually align the SDK on top. Neither card may split OpenAPI from SDK.
- `packages/sdk-runtime/src/index.ts` has one F01 writer only. The F02 export is a hub merge action after F02-S7, never a chip card.
- All three migrations and both runner-count assertions are one F02-S2 atomic write-DAG unit (41→44), so the fixture is never observed at 42/43 with a three-file migration set or vice versa. No 0048/0049 migration is planned: 0047 includes the reader/latest/hash indexes. If implementation proves another index necessary, it requires a revised card and monotonic 44→45 bump.
- `apps/server_core/internal/composition/root.go` and `contracts/governance/modules.json` have one writer, F03-S2, after F02 exposes constructors. F02 never wires itself.
- F03 starts only after the F01 identity projection and F02 persistence/service artifacts exist. F03-S3 is the join/final regression gate.

## C. Contract-satisfiability check

### F-01 catalog schemas and SDK

The current catalog read DTOs are:

- `CanonicalCatalogProduct`, `contracts/api/marketplace-central.openapi.yaml:2786-2815`, returned by `GET /catalog/products/{id}`. It already has nullable `ean`, `manufacturer_reference`, and `brand_name` (`:2795-2806`). Add `ncm` and `quality_flags` to `required`; add `ncm: {type: string, nullable: true, pattern: '^[0-9]{8}$'}` and `quality_flags` as a non-null array whose items enum is `[invalid_ean, ean_collision]`.
- `CatalogProductFact`, `contracts/api/marketplace-central.openapi.yaml:2816-2839`, used by `CatalogProductFactPage` at `:2864-2871` for list/search. It has nullable `ean` but only the legacy `reference` field. Retain `reference` to avoid a rename; define it explicitly as a deprecated compatibility alias of `manufacturer_reference`. Add required nullable `manufacturer_reference`, `brand_name`, `ncm`, plus required non-null `quality_flags` with the same enum. Thus IC-01 `refforn` maps to API `manufacturer_reference`; `marca` maps to `brand_name`; `ean` and `ncm` keep their names.
- Do **not** edit `CatalogProduct` (`:2892-2935`): it is the legacy enrichment DTO, not the canonical GET/list read contract. Other `ean` occurrences belong to listing/integration/inventory schemas and are outside F-01.
- Apply the same additive fields to the manual TypeScript interfaces at `packages/sdk-runtime/src/index.ts:164-204`. There is no generator/config or generate script; F01-S4 manually maintains parity.

The identity flags are necessary to satisfy IC-01's observable invalid-EAN warning and collision signal; without them, the current HTTP DTO at `catalog/transport/http_handler.go:213-222` drops `ProductCandidate.QualityFlags` and the EARS contract is unsatisfiable.

### F-02 paths and components

The current `paths` mapping occupies OpenAPI lines 5-2498 and `components` begins at line 2499. No `/erp/imports` path or `ErpImport*` component exists, so F02-S7 adds a new, non-colliding path block immediately before current line 2499:

- `/erp/imports`: `post` multipart/form-data with required binary `file`; responses 201 `ErpImportCreated`, 400 `ErpImportError`, 409 `ErpImportConflict`, 422 `ErpImportError`. `get` returns `ErpImportList` newest-first.
- `/erp/imports/{id}`: `get` with required string/UUID path parameter and 200 `ErpImportDetail`, 404 `ErpImportError`.
- New schemas placed together at the end of `components.schemas`: `ErpImportStatus` (`COMPLETED|REJECTED`), `ErpImportCreated` (`import_id`, `protocol`, `status`), `ErpImportSummary` (id/protocol/status/sanitized filename/imported_at/accepted_count/rejected_count/warning_count), `ErpImportIssueKind` (`REJECTION|WARNING`), `ErpImportIssue` (row, nullable column, kind, code enum, detail, nullable offending_value), `ErpImportDetail` (summary plus rejection/warning arrays), `ErpImportList` (items), `ErpImportError` (standard code/message/details), and `ErpImportConflict` (code plus nullable existing `import_id`/`protocol`; populated for `DUPLICATE_FILE`).
- Exact error-code enums cover `INVALID_FILE`, `MISSING_REQUIRED_COLUMN`, `IMPORT_IN_PROGRESS`, `DUPLICATE_FILE`, and `ERP_IMPORT_NOT_FOUND`; issue code enums cover the seven exact IC-02 rejection/warning codes.

This region is disjoint from F01's product schemas. F02 manually authors `packages/sdk-runtime/src/erpImport.ts` because the package is manual; it does not regenerate or overwrite F01 catalog types.

### Fixed GTIN seam

```go
package domain // marketplace-central/apps/server_core/internal/modules/catalog/domain

func IsValidGTIN(value string) bool
```

The function trims outer whitespace, then accepts only exactly 8, 12, 13, or 14 ASCII digits whose GS1 mod-10 check digit is correct. It does not normalize punctuation, strip internal whitespace, or return a cleaned value. F02 persists the trimmed original only when this returns true.

### F-03 Reader conformance

The adapter implements all six methods from `internal_read/ports.Reader` (`reader.go:48-55`) and the optional catalog page interface:

| Method | XLSX behavior |
| --- | --- |
| `FindProductsForLinking` | Resolve latest completed tenant snapshot; filter by CODPROD/ProductID, valid EAN, SellerSKU-as-CODPROD only, normalized title, and active=true (snapshot rows are active). Return identity from stored normalized fields, including NCM and stored warnings/collision flag. No snapshot → wrapped `ErrNoErpSnapshot`; requested missing product → `ERPProductNotFoundError`; never return a fabricated empty candidate. |
| `GetSellableStock` | Load product from latest completed snapshot. Reserved present (including explicit zero) → `Quantity = physical - reserved`, source observed/fetched at `imported_at`. Reserved null → `Quantity=nil` plus `QualityMissingStock`; physical remains persisted/queryable as physical and is never substituted as sellable. Missing snapshot/as-of → `ErrNoErpSnapshot`; missing CODPROD → typed not-found. |
| `GetCurrentPrice` | Always return `ReadErrorUnsupportedQuery` with message `xlsx snapshot does not contain current price`; never return a zeroed success, even when a snapshot exists. |
| `GetCostAsOf` | Choose newest completed snapshot with `imported_at <= Policy.EffectiveAt`; exact positive cost, per-unit scope, `ObservedAt=imported_at`. No qualifying snapshot → `ErrNoErpSnapshot`; missing CODPROD → typed not-found. |
| `GetSalesHistory` | Always return `ReadErrorUnsupportedQuery` with message `xlsx snapshot does not contain sales history`; never return an empty success as if known. |
| `GetTaxInputs` | Resolve snapshot/product as above. NCM present → `TaxInputs.NCM` with source time `imported_at`, monetary tax fields nil, because the workbook contains classification rather than calculated tax amounts. NCM absent/invalid → NCM nil plus `QualityMissingTax`. No snapshot → `ErrNoErpSnapshot`; missing CODPROD → typed not-found. |

`ListCatalogProductFacts`/`SearchCatalogProductFacts` project the same identity fields and honest stock/cost facts, with current price unknown; this is required for M-01-C03's GET `/catalog` proof. The method signatures remain unchanged; adding nullable NCM to existing return structs is additive and lands under F-01's `internal_read/**` ownership.

### Contract tensions resolved by this plan

- Current Oracle code explicitly says `REFERENCIA` is manufacturer reference and sets EAN nil (`internal_read/adapters/oracle/reader.go:69-75`); current catalog SQL selects `REFERENCIA` and a hardcoded null EAN (`catalog_page.go:124-133`). ADR-03 + IC-01 are higher truth, so F01-S2/S3 deliberately reverse those assertions and tests.
- IC-02 asks `GetTaxInputs` to carry NCM, but current `TaxInputs` has no NCM (`internal_read/domain/internal_tax.go:29-40`). F01-S2 adds nullable NCM to the return struct without changing the Reader interface; monetary tax values remain unknown.
- IC-02 says physical stock remains consultable while sellable becomes unknown when reserved is absent, but current `SellableStock` exposes only sellable `Quantity` (`internal_read/domain/internal_stock.go:28-34`). This plan preserves physical stock in the snapshot/DB and returns unknown sellable stock; it does not mislabel physical as sellable or expand the frozen port.
- The F-02 brief says reuse the existing HTTP guard, but no auth/tenant request guard exists in the current root. The existing baseline is composition-time `DefaultTenantID`; F03-S2 reuses that exact pattern and adds no auth surface.

## D. Pre-identified additive locks

| Lock | Scope | Holder / serialization | Release condition |
| --- | --- | --- | --- |
| Catalog SDK types | Only additive `ncm`, governed manufacturer/brand fields where absent, and identity flags in `packages/sdk-runtime/src/index.ts:164-204`; no unrelated formatting/client edits | F01-S4 | Released after the atomic F01 OpenAPI+SDK commit lands and its build/test pass. |
| OpenAPI catalog schema | `CanonicalCatalogProduct` and `CatalogProductFact` blocks only | F01-S4 before F02-S7 | Released when F01-S4 lands; F02 rebases/manual-parity checks before editing `/erp/imports`. |
| OpenAPI ERP import | `/erp/imports*` and `ErpImport*` components only | F02-S7 | Released after its atomic OpenAPI+SDK client commit. |
| SDK barrel | One line exporting `erpImport.ts`; F01 catalog edits do not grant any other barrel change | Hub | Hub adds it after F02-S7 is merged and both SDK checks pass. |
| Composition root | Additive erp-import registration and source selection in `internal/composition/root.go` | F03-S2 sole writer | Released after F03-S2 composition tests and governance check. |
| Governance registry | One sorted `erp_import` entry | F03-S2 sole writer via chip merge | Released only after clean-worktree governance validation. |
| Migration inventory fixture | Three files 0045-0047 plus both count assertions 41→44 | F02-S2 single atomic slice | Released after migration test sees exact canonical count 44. |

## E. Per-criterion verification map

| Criterion | Satisfying slice(s) | Named command / QA step | Proof carrier |
| --- | --- | --- | --- |
| `M-01-C01` complete synchronous XLSX, protocol, duplicate/no raw logs | F02-S3, F02-S4, F02-S5, F02-S6, F03-S2, F03-S3 | F02 unit/integration commands; QA `POST /erp/imports` twice with `example-erp.xlsx`, inspect `SELECT` and grep server logs for fixture cost/description (zero hits) | xlsx/parser tests; postgres repo tests; import service/handler tests; `.mnfs/.../validation-result.md §import` |
| `M-01-C02` honest per-row report/warnings/all rejected | F02-S1, F02-S3, F02-S5, F02-S6, F03-S3 | `go test ./internal/modules/erp_import/...`; QA POST `identity-rejections.xlsx` | domain/parser/service/handler tests; fixture; `validation-result.md §rejeicao` |
| `M-01-C03` IC-01 identity Reader/API/collision/null | F01-S1..S4, F03-S1, F03-S3 | F01 contract commands; F03 adapter command; QA GET imported product | GTIN, Oracle reader/page, catalog handler, source contract tests; `validation-result.md §identidade` |
| `M-01-C04` snapshot cost/stock and reserved-unknown | F02-S2, F02-S4, F03-S1 | `go test -tags=integration ./internal/modules/erp_import/adapters/postgres ./internal/modules/erp_import/adapters/internalread` | migration/repository/reader integration tests; `validation-result.md §reader` |
| `M-01-C05` Oracle intact with new semantics | F01-S2, F01-S3, F03-S2, F03-S3 | Oracle unit contracts + `npm run harness:live` + `go build ./cmd/server` | Oracle tests, source contract test, `validation-result.md §oracle` live transcript |
| `M-01-C06` migrations/seams/atomic OpenAPI-SDK | F02-S2, F01-S4, F02-S7, F03-S2 | migration command, SDK build/test after each OpenAPI card, `npm run harness:governance`, diff/commit-boundary review | runner test, OpenAPI, SDK tests, modules registry, `validation-result.md §seams` |
| F01 EARS: valid REFERENCIA→EAN and REFFORN independent | F01-S2, F01-S3 | F01 Oracle/catalog contract commands | `reader_test.go`, `catalog_page_test.go`, catalog adapter/handler tests |
| F01 EARS: invalid/blank REFERENCIA→null+warning, never REFFORN | F01-S1..S3 | GTIN + Oracle/catalog contract commands | GTIN and identity fixture tests |
| F01 EARS: duplicate EAN preserves both + collision flag | F01-S2, F01-S3 | `go test ... -run 'Test.*Collision'` | Oracle reader/page and catalog HTTP tests |
| F02 EARS: valid synchronous full snapshot→201 COMPLETED/counts | F02-S3..S6 | parser/repo/service/handler commands | corresponding tests and `validation-result.md §import` |
| F02 EARS: partial invalid rows only rejected | F02-S1, F02-S3, F02-S5 | domain/parser/service commands | validation/parser/service tests |
| F02 EARS: all rows rejected→201 REJECTED/report/no products | F02-S4..S6 | repo/service/handler commands | integration and handler tests |
| F02 EARS: required header absent→422/no products | F02-S3, F02-S5, F02-S6 | parser/service/handler commands | missing-custo fixture and tests |
| F02 negatives: invalid file, duplicate, concurrent, exact cost/stock semantics | F02-S1, F02-S3, F02-S4, F02-S5 | all F02 commands including integration tag | parser/domain/repository/service tests |
| F03 EARS: cost with imported_at source time | F03-S1 | F03 integration command | `reader_integration_test.go` |
| F03 EARS: no qualifying snapshot→typed ErrNoErpSnapshot | F03-S1 | F03 unit/integration command | `reader_test.go`, `reader_integration_test.go` |
| F03 EARS: oracle config preserves current path | F03-S2, F03-S3 | composition test + live lane | `root_test.go`, source contract, `validation-result.md §oracle` |
| F03 negatives: missing CODPROD, as-of before import, newer REJECTED ignored | F03-S1 | F03 integration command | `reader_integration_test.go` |
| ADR-17 all methods: unsupported/unknown never zero-default | F01-S2, F02-S1, F02-S4, F03-S1 | domain/repository/reader commands | negative tests for nil/typed errors and integrity rollback |

QA owns creation of `validation-result.md`; implementation tests are necessary but do not substitute for the fresh local-stack HTTP/DB/live-Oracle evidence named by the validation contract.

## F. Open questions and hub REQUESTs

1. **REQUEST DEP-GRANT-01 (blocking F02-S3 and therefore F02-S5 onward):** approve `github.com/xuri/excelize/v2` for isolated use in `erp_import/adapters/xlsx`, including `go.mod` and `go.sum`. Resolution is a recorded hub grant or a hub-selected approved alternative. Until then F02-S3 is not dispatch-ready; cards that depend on it cannot start.
2. **REQUEST BARREL-01 (merge seam, non-blocking implementation):** after F02-S7, the hub adds the one-line `packages/sdk-runtime/src/index.ts` export for `erpImport.ts`, then reruns SDK build/test. F02 does not own this line.
3. **Operator fixture approval (blocks only final QA, not implementation):** approve a non-sensitive ≥50-row example workbook for `.mnfs/.../fixtures/example-erp.xlsx`. If the real client workbook contains sensitive data, it must not be committed; derive a sanitized contract-equivalent fixture and use the real file only in the authorized local QA run.

There is no unresolved ADR-vs-ADR contradiction. The REFERENCIA remap is explicitly aligned between ADR-03 and IC-01; only current code/tests contradict it. All cards except F02-S3 are individually dispatch-ready subject to their `depends_on`; F02-S3 is explicitly blocked by DEP-GRANT-01.

## G. Orchestrator review verdict (CHIP-M01, 2026-07-17)

**ACCEPTED with 4 binding corrections applied in-place:**
1. `erp_import/adapter/internalread` → `erp_import/adapters/internalread` (plural, matches repo `adapters/` convention + F02's own xlsx/postgres dirs). Applied to F03-S1/S2/S3.
2. F03-S3: `npm run harness:live` struck from implementer commands — chip cannot load `.env*`/Oracle creds/boot server; live-Oracle is QA-only at P7. Implementer uses hermetic Oracle-shaped fake + `go build`.
3. F03-S2: `npm run harness:governance` annotated clean-worktree/40-hex-BaseSha (false-fails in-chip; entry validated at merge).
4. Environment prerequisite (not a plan defect): `.gomodcache` is warm; **`npm ci` at repo root is required BEFORE the SDK slices** (F01-S4, F02-S7) — node_modules not yet installed. Run once in bootstrap.

Cross-checked accurate against independently-gathered evidence: GTIN seam-first ordering; REFERENCIA→GTIN→ean reversal with full-active-population collision query; Reader 6-method conformance (`ports/reader.go:48-55`) with per-method honest-unavailable (current_price/sales_history → unsupported_query; reserved-null → Quantity nil + QualityMissingStock, never zero); `ErrNoErpSnapshot` via `NewReadError(source_unavailable,…)` dual-matchable; migrations 0045-0047 + `runner_test.go` 41→44 atomic; OpenAPI DTO scoping (`CanonicalCatalogProduct`/`CatalogProductFact` edited, legacy `CatalogProduct` excluded) consistent with observed `ean` line anchors; one-writer seams (OpenAPI serialized F01-S4→F02-S7, barrel=hub, root.go+modules.json=F03-S2).

**Execution order (respecting DAG + dep-grant block):**
Wave 1 (parallel, disjoint writes): **F01-S1** (validator, lands first) → then F01-S2, F02-S1, F02-S2 concurrent.
Wave 2: F01-S3 · F02-S3 (gated on DEP-GRANT-01) ∥ F02-S4.
Wave 3: F01-S4 · F02-S5 → F02-S6 → F02-S7 (after F01-S4).
Wave 4: F03-S1 → F03-S2 → F03-S3 (join/final).
Every green slice: independent Claude-sonnet reviewer before any dependent slice starts. Dual gate (cold Opus subagent + Sol-medium) + QA live-drive at close.
