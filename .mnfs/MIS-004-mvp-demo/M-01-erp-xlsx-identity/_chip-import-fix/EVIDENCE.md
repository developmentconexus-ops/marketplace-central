# CHIP-IMPORT-FIX — Evidence Pack

Branch: `claude/zealous-mendel-d7699e` (UNPUSHED)
Base SHA: `9e0beb41ad47f3145e3002f384b4cf7df6bdb1ee`
Scope: (A) real customer catalog cannot be imported (parser); (B) `/integracoes`
"Configuração da plataforma" screen is a stub; (C) imported `catalogo_cliente`
products are invisible because the `erp_source` active-source toggle was never
wired. All three fixed for the 2026-07-20 demo.

---

## GOAL A — backend parser (`erp_import/adapters/xlsx/parser.go`)

Three defects fixed against the raw Sankhya "Produto" export shape:

1. **Multi-sheet union.** `parseWithRequired` iterates ALL sheets in
   `GetSheetList()` order and UNIONs the data rows (was: `GetSheetList()[0]`
   only). Sheet name is carried as the product category (`DescrGrupo`) on the
   lenient path when the sheet has no explicit group column.
2. **Header-row detection past preamble.** `detectHeader` finds the first row
   carrying both identity columns (CODPROD + DESCRPROD, EN or PT alias); the
   title/emission preamble rows above it are skipped. A sheet with no detectable
   product header is skipped with a **WARNING file-issue**, never a hard-fail of
   the whole file (ADR-17).
3. **PT header aliases.** `canonicalHeaderKey` folds `Código`→CODPROD,
   `Descrição`→DESCRPROD, `Código de Barra(s)`→EAN,
   `Referência do Fornecedor`→REFFORN, `Marca`→MARCA (two Marca columns, first
   non-empty wins), `NCM`→NCM. English strict keys retained. `Estoque Mínimo` is
   NOT physical stock (unmapped). Absent CUSTO/ESTOQUE_FISICO → honest-unknown
   (empty), never fabricated 0.

**Strict path preserved.** `Parser.Parse` still requires
`[CODPROD, DESCRPROD, CUSTO, ESTOQUE_FISICO]`; a detected header missing a
required column returns `MISSING_REQUIRED_COLUMN` FileError (hard-fail). Only the
lenient `ParseLenient` path relaxes CUSTO/ESTOQUE and injects sheet-name category.

**Signature change (threaded through 4 owned files):** `Parse`/`ParseLenient`
now return `([]NormalizedRow, []Issue, error)` — the middle slice is file-level
issues (skipped sheets). `ports.Parser`, `application/import_service.go`
(`issues = append(fileIssues, rowIssues...)`), and both test files updated.

**Migration avoidance:** skipped-sheet warning reuses the existing
`MISSING_REQUIRED_COLUMN` code with `Kind=Warning` so persistence stays inside
the `erp_import_issues.code` CHECK constraint (0047+0073). A dedicated
`SHEET_SKIPPED` code would need a hub-owned migration → flagged as follow-up.

### Go tests — GOCACHE/GOMODCACHE absolute, cwd = worktree apps/server_core

```
go test ./internal/modules/erp_import/...
ok  .../erp_import/adapters/internalread   2.732s
ok  .../erp_import/adapters/postgres       2.226s
ok  .../erp_import/adapters/xlsx           2.226s
ok  .../erp_import/application             1.697s
ok  .../erp_import/domain                  1.923s
ok  .../erp_import/transport               2.670s
EXIT: 0
```

New / evolved xlsx tests (all PASS):
- `TestParserRawExportGolden2012` — raw 3-sheet export (FACAS 138 + FERRAMENTAS
  1837 + MAQUINAS 37 = **2012**) → 2012 accepted, 0 rejected, 0 file-issues,
  ZERO fabricated custo/stock, category counts exact, 2012 MISSING_CUSTO + 2012
  MISSING_ESTOQUE honest-unknown warnings.
- `TestParserRawExportHeaderDetectionSkipsPreamble` — 2 preamble rows skipped.
- `TestParserRawExportPTAliases` — every PT alias + Marca first-non-empty across
  the twin Marca columns.
- `TestParserRawExportSkipsHeaderlessSheetWarnNotFail` — headerless RESUMO sheet
  → 1 WARNING, good sheet survives, file never hard-fails.
- `TestParserHeaderFoldingAndUnionsAllSheets` — union assertion (evolved from the
  old FirstSheetOnly test that encoded defect #1; behaviour change is mandated).

### Guards proven load-bearing (revert → must-fail, then restored)
- Drop PT aliases → INVALID_FILE (no identity header detected).
- Break union (first sheet only) → golden returns 138 not 2012.
- Force header index 0 → MISSING_REQUIRED_COLUMN (preamble row read as header).

---

## GOAL B — `/integracoes` "Configuração da plataforma" screen

New page module (single-line AppRouter change):
- `apps/web/src/pages/integracoes/IntegracoesPage.tsx` — header, **Importar
  catálogo** card (source selector `Catálogo do cliente`→`catalogo_cliente`
  lenient [default] vs `Sankhya (custo + estoque)`→`xlsx` strict; file
  picker + drag-drop `.xlsx`; Importar button; live result summary
  protocolo/aceitos/rejeitados/avisos + issue lists; PT error/duplicate/
  missing-column/loading states), an **inert** ML connect card
  (`disponível em breve`, disabled, no network), and the reused read-only
  `ImportacaoSection` history below.
- `apps/web/src/pages/integracoes/useErpImportUpload.ts` — mutation hook: POSTs
  multipart via `client.createErpImport`, invalidates the shared import-history
  list on success, normalizes failures to a typed `ErpImportUploadError`
  (invalid_file / missing_required_column / duplicate_file / import_in_progress
  / internal_error / network_error).
- `apps/web/src/routes/integracoes.tsx` — route wrapper (pedidos pattern).
- `AppRouter.tsx` — one route element swapped (`WorkspacePlaceholder` →
  `IntegracoesRoute`) + one import line. `/mercado` and `/marketplaces` (CHIP-
  MERCADO's seam) untouched.

### SDK + contract lockstep
- `packages/sdk-runtime/src/index.ts` — new `createErpImport(file, source,
  fileName?)` building `FormData` (`file` + `source`), posted via a new
  `postMultipart` helper (no manual Content-Type; surfaces the parsed error body
  so callers branch on 409/422).
- `packages/sdk-runtime/src/erpImport.ts` — `ErpImportSourceInput` union
  (`xlsx | catalogo_cliente`); `ErpImportSummary.source` widened to the union;
  issue codes extended with `MISSING_CUSTO`, `MISSING_ESTOQUE`,
  `MISSING_REQUIRED_COLUMN` (backend truth).
- `contracts/api/marketplace-central.openapi.yaml` — multipart `source` enum
  field documented; `ErpImportSummary.source` enum `[xlsx, catalogo_cliente]`;
  `ErpImportIssueCode` enum extended to match the SDK.

### apps/web vitest (throwaway chip config: abs setupFiles + fs.allow for the
node_modules junction; config DELETED pre-commit)
```
IntegracoesPage.test.tsx  7 tests PASS
AppRouter.test.tsx       19 tests PASS  (incl. rewritten /integracoes assertion)
Full apps/web suite: 366 passed | 1 failed (367)
```
The single red is `PricingMatrix.test.tsx` (precos) — a file this chip does NOT
touch (`git diff --stat HEAD -- apps/web/src/pages/precos/` = empty); it is the
known pre-existing baseline red. My in-scope surface is 100% green.

New IntegracoesPage tests assert: heading renders; upload POSTs
`createErpImport(file, "catalogo_cliente", name)` by default and renders the
result summary from the fetched detail; Sankhya selection POSTs
`(file, "xlsx", name)`; 409 duplicate shows the prior protocol; 422 shows the
missing column name; a non-.xlsx file is rejected locally with NO network call;
the inert ML affordance is disabled and never calls the network.

### tsc --noEmit
- `apps/web` — 0 new errors. The reported errors are the documented baseline
  (anuncios*, ListingsRefresh, Mutation*, Produto* — all pre-existing, none in
  the chip write-set: integracoes/, useErpImportUpload, AppRouter.tsx,
  sdk-runtime absent).
- `packages/sdk-runtime` — **0 errors** (clean).

### LIVE-VERIFIED
Throwaway vite on free port **5199** (NOT 5174/8080), backend absent.
`LIVE-VERIFIED: /integracoes renders — header "Configuração da plataforma",
Importar-catálogo card with working source selector (Catálogo do cliente ⇄
Sankhya toggles), .xlsx dropzone + Selecionar arquivo, Importar disabled with no
file, inert ML "Conectar — disponível em breve" disabled, Importação history
section. Dark theme confirmed; light parity by semantic-token discipline
(bg-surface/text-ink/border-border/accent-soft — no hardcoded colors).`

The real end-to-end 2012-row import through the live UI is **HUB-VERIFIED** (needs
the backend dev stack at :8080) — NOT claimed here.

---

---

## GOAL C — active-source toggle (`erp_source`)

Operator-ratified addition: imported `catalogo_cliente` products were invisible
because `internalread.WithActiveSource(ctx, source)` existed but was never called.
This goal wires the toggle end-to-end (transport → reader snapshot selection → FE).

### Backend — validated `erp_source` at three read seams
- `erp_import/adapters/internalread/reader.go` — new exported
  `ParseActiveSource(raw) (ImportSource, bool, error)`: `""`→(_,false,nil) [absent,
  default], `xlsx`/`catalogo_cliente`→(source,true,nil), anything else→
  `ErrUnknownActiveSource`. Case-sensitive; NO silent fallback. Exported
  `ActiveSourceFromContext` (present=true only when set); private
  `activeSourceFromContext` routes through it. `LatestCompletedSnapshot` /
  `Reader` unchanged — ctx stays a pure passthrough.
- `catalog/transport/http_handler.go` — `requestContext(r)` now returns
  `(ctx, error)`, parses `erp_source`; unknown → `catalogPageQueryError{code:
  "invalid_erp_source"}` → `writeCatalogPageError` 400 `{"error":
  "invalid_erp_source"}`; present → `WithActiveSource`. Both `handleProducts` and
  `handleSearch` hoist `ctx, err := requestContext(r)` and bail on error before
  any port call. Absent param → no `WithActiveSource` → byte-stable default.
- `market/transport/collection_handler.go` — `handleCollect` parses
  `erp_source` after `ParseCollectionRequest`; unknown → `writeMarketQueryError(
  &InvalidFilterError{Key:"erp_source"})`; present → `ctx =
  WithActiveSource(ctx, source)` BEFORE `Collect(ctx, req.Codprod)`. ctx chain
  confirmed to `GetLocalIdentity`→`FindProductsForLinking`→`snapshot(ctx)`→
  `activeSourceFromContext`, so identity resolves from the selected snapshot.

### FE — visible, persisted "Fonte ativa" selector threaded into /catalogo
- `packages/web-query/src/activeErpSource.ts` (NEW) — `ActiveErpSource` union,
  `DEFAULT_ERP_SOURCE="xlsx"`, localStorage key `mc.active-erp-source`,
  get/set (set dispatches a CustomEvent + writes storage), `subscribe`
  (CustomEvent + cross-tab `storage`), `useActiveErpSource()` via
  `useSyncExternalStore`. Unrecognized stored value → default (never trusts junk).
- `packages/feature-products/src/catalogQueries.ts` — hooks take `erpSource`,
  put it in the `queryKey` (so switching source refetches, no stale cross-source
  cache), and pass `erp_source: sourceParam(erpSource)` to the client.
  `sourceParam` returns `undefined` for the default (xlsx) → param omitted →
  URL byte-stable with pre-toggle clients.
- `CatalogPage.tsx` — `const [activeSource] = useActiveErpSource()`, threaded
  into both query hooks and both refresh refetch keys.
- `IntegracoesPage.tsx` — `ActiveSourceCard` radio group (`ERP Sankhya` /
  `Catálogo do cliente`, testids `active-source-selector` /
  `active-source-{xlsx,catalogo_cliente}`) rendered above the upload card; hint
  states catalogo_cliente shows custo/estoque as "—" (ADR-17 honest-unknown).
- `packages/web-query/src/index.ts` — `catalogQueryKeys.search` appends the
  params object ONLY when params are supplied, so the pre-toggle two-arg key
  shape `["catalog","search",q]` stays byte-stable (guard test retained).

### SDK + OpenAPI lockstep (additive grant, same commit)
- `sdk-runtime/src/index.ts` — `erp_source?: ErpImportSourceInput` on
  `CatalogPageOptions` + `CatalogSearchPageOptions`, threaded into the list/search
  query builders; `catalogQuery` drops undefined keys so absent stays absent.
- `sdk-runtime/src/market.ts` — `collectMarketPriceIntel(codprod, erpSource?)`
  appends `?erp_source=` only when supplied (file kept independent of index.ts).
- `contracts/api/marketplace-central.openapi.yaml` — additive `erp_source`
  query param (enum `[xlsx, catalogo_cliente]`) on `/catalog/products`,
  `/catalog/products/search`, `/market/collections` (+ a 400 on collections).

### Golden matrix — regression tests (all PASS)
- `catalog/transport/http_handler_test.go` —
  `TestCatalogPageThreadsActiveSourceToggle` (absent→no WithActiveSource;
  catalogo_cliente/xlsx→threaded; search too) +
  `TestCatalogPageRejectsUnknownActiveSource` (400, port NOT called).
- `market/transport/collection_handler_test.go` —
  `TestCollectThreadsActiveSourceIntoIdentityResolution`,
  `TestCollectAbsentActiveSourceIsByteStable`,
  `TestCollectRejectsUnknownActiveSource` (400 `key:erp_source`, service not
  called).
- `internalread/reader_test.go` — `TestParseActiveSource` (table: absent,
  both valid, case-sensitivity, unknown) + `TestReaderSelectsSnapshotByActiveSource`
  (absent→xlsx snapshot; catalogo_cliente→prospect snapshot with honest-unknown
  `missing_cost`/`missing_stock`, never zero).
- `packages/web-query/src/activeErpSource.test.ts` (NEW, 4) — default xlsx;
  ignores junk stored value; persists + reads back; multi-hook sync on change.
- `packages/feature-products/src/CatalogPage.test.tsx` — new: default view omits
  `erp_source` from the fetch URL; `catalogo_cliente` active → URL carries
  `erp_source=catalogo_cliente`.
- `apps/web/.../IntegracoesPage.test.tsx` — new: Fonte ativa defaults to Sankhya,
  switching to Catálogo do cliente persists (`getActiveErpSource()`).

### Verification snapshot (A+B+C, pre-gate)
```
Go write-set  (GOCACHE/GOMODCACHE abs, cwd apps/server_core):
  go build ./...                                        BUILD_OK
  go vet catalog/transport market/transport internalread VET_OK
  gofmt -l (3 changed .go)                              clean
  go test erp_import/... catalog/transport market/transport   ALL ok

FE (real package configs — chip vitest config used for apps/web only, DELETED):
  packages/sdk-runtime  (own config)     71/71 pass
  packages/web-query index+activeErpSource  pass
  packages/feature-products CatalogPage     pass
  apps/web IntegracoesPage + AppRouter      pass

tsc --noEmit:
  write-set files (sdk-runtime, web-query, feature-products, integracoes) 0 errors
  apps/web total 12 errors — ALL pre-existing baseline in NON-write-set surfaces
  (anuncios*, mutations*, produto*, ListingsRefresh); none cite chip files.
```

### Harness artifacts (NOT source regressions — proven)
- `PricingMatrix.test.tsx (a)` fails ONLY under the throwaway apps/web chip
  vitest config; it fails IDENTICALLY (9/10) with packages resolved to MAIN's
  unmodified source (no worktree alias), so it is independent of this chip's
  edits — a config-timing race (market read renders a skeleton row before flush).
  Green on main under the real config; /precos is not this chip's seam.
- The 4 sdk-runtime openapi-parity tests fail under the chip config only because
  it runs from the worktree ROOT and those tests read the spec via
  `process.cwd()`-relative `../../contracts/...`; under each package's own config
  (correct cwd) they are 71/71 green.

## Hub follow-ups flagged
1. Dedicated `SHEET_SKIPPED` issue code (needs a migration to the
   `erp_import_issues.code` CHECK) — currently reusing `MISSING_REQUIRED_COLUMN`
   as a Warning to stay migration-free.
2. `FINDING-raw-sankhya-produto-export.md` did not exist at base SHA; parser was
   built to the EXEMPLO-IO shape in the dispatch prompt. Confirm the real
   export's exact header labels against a live file during hub live-drive.

## Notes / discipline
- Base-relative, minimal AppRouter edit; CHIP-MERCADO seam lines untouched.
- No server booted by the chip; no deps installed; node_modules is a junction to
  main (gitignored, not committed); chip vitest config deleted pre-commit.
- P6 dual gate (cold reviewer + adversarial refuter) run over this whole surface
  — verdicts recorded below.
