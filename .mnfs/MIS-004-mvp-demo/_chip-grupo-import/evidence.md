# CHIP-GRUPO-IMPORT — evidence pack

Mission: MIS-004-mvp-demo · Scope: DESIGN-TARIFAS-ML §11 item 3 (ratified D-84)
Base SHA (hub-confirmed): `f6bcb55e5d91974d03d5b34fab28d55235007d61`
Branch: `chip/grupo-import`
Migration grant (hub): `0069` (numbering ruling D-8x: next-free at grant, no future reservations)

## Scope delivered
Capture ERP xlsx `grupo`/`descrição grupo` columns the importer previously discarded — additive,
nullable, honest-absent (ADR-17). NO de-para grupo→categoria ML (CHIP-CATMAP pós-demo), NO
category tree, NO new UI, NO new endpoint. Data stored now; catmap consumes later.

## Edit surface (8 files, zero scope creep)
1. `internal/modules/erp_import/domain/import.go` — NormalizedRow +`Grupo`/`DescrGrupo *string` (nullable)
2. `adapters/xlsx/parser.go` — `optionalCell(...,"GRUPO")` / `optionalCell(...,"DESCRGRUPO")`
   (same header-fold path as EAN/Marca/NCM: lowercase + accent-fold, `optionalCell` empty→nil)
3. `adapters/postgres/import_repository.go` — INSERT +`grupo,descrgrupo` ($12,$13)
4. `adapters/postgres/query_repository.go` — SELECT + Scan +`grupo,descrgrupo`
5. `migrations/0069_erp_import_products_grupo.sql` — `ADD COLUMN IF NOT EXISTS grupo/descrgrupo text` (nullable, NO default)
6. `internal/platform/migrate/runner_test.go` — 2 hardcoded counts 55→56 (+ hub reconcile comment)
7. `adapters/xlsx/parser_test.go` — grupo present / absent-column→nil / accent+case fold
8. `adapters/postgres/query_repository_test.go` — grupo write→read round-trip (integration)

## OpenAPI/SDK: N/A (verified)
No transport/query_service/OpenAPI response exposes NormalizedRow product fields (ean/marca/ncm
today are internal-snapshot only). grupo follows suit. No endpoint invented (YAGNI, dispatch item 4).

## Gate results
- `go build ./...` = exit 0
- `go vet ./internal/modules/erp_import/...` + `-tags integration` = exit 0
- `go test ./internal/modules/erp_import/...` = ok (xlsx/domain/application/transport/internalread/postgres)
- `go test ./internal/platform/migrate/...` = ok (both count assertions 56)
- Hermetic integration lane (`npm run harness:integration`, ephemeral postgres, warm modcache):
  `status=passed`, `migrations_first=56` (0069 applied clean over base — no false HPG -1),
  `migrations_second=0` (idempotent), run_id=1ecba0f8828147bdb2e878c46b03f764
- Fixture "com e sem colunas grupo": parser_test builds both in-memory (xlsxBytes) — legacy file
  WITHOUT grupo columns → import OK, fields nil (TestParserValidAndInvalidRows absent assertion);
  file WITH columns → captured (TestParserCapturesGrupoColumns); empty grupo cell → nil.

## Open note flagged to hub (header string)
Header keys chosen: `GRUPO` / `DESCRGRUPO` (parser convention UPPERCASE_UNDERSCORE + design §11
shorthand "grupo/descrgrupo"). Exact real Sankhya-export header for descrição-grupo unconfirmed
(hub owns client file). If real header differs (e.g. "DESCRICAO GRUPO" space-form), column is
honestly-nil (never fake) and fix is a 1-line header string. HUB: verify at post-merge reimport
of the real xlsx; column captured non-null = confirmation.

## Numbering / merge note (hub duty)
Forked from f6bcb55e (has 0067, NOT T1's 0068 — T1 unmerged). Local tree = 56 migrations
(base 55 + 0069). T1's 0068 lands in parallel; hub reconciles runner_test count to 57 at the
second merge (comment left on both count lines).
