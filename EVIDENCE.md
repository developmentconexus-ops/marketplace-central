# CHIP-VENDAVEL — EVIDENCE (P5 ladder, hub takeover per operator order 2026-07-31)

Tree: `worktree-chip-vendavel`, tip at ladder time `a2f96fa1` (F6 @86fc6400, cond.6 @83c91556,
F4/F4b/F5 @80f0319c, custody @c8f7d90e). BASE-SHA de despacho `554788d5`; gate/merge base =
main tip `4ad36272` (re-lido na execução).

## Pré-condições

- Árvore limpa antes das lanes; `git status --porcelain` vazio; zero `.js` emitido em packages/ (find vazio).
- pg-session recriado no worktree (A-25: endpoint morto se RECRIA): container `mpc-pg-session-3eee515d`, porta 50265, run_id `d1856cbaa2f34800bdd3dee8065d9d59`.
- `.gomodcache` aquecido (147M) antes da lane de integração.

## Migrate (A-25)

DB `mpc_test_d1856cbaa2f34800bdd3dee8065d9d59` (casa com `^mpc_test_[0-9a-f]{32}$`).
`go run ./cmd/testdb migrate` → `applied 72 migration(s)`; segunda rodada → `applied 0 migration(s)`.

## Lanes

| Lane | Comando | Resultado |
|---|---|---|
| No-lane packages (B-6) | `go test ./internal/modules/tenant_config/... ./internal/composition/... -v -count=1` com MPC_TEST_DATABASE_URL da sessão | exit=0, **PASS=115 SKIP=0 FAIL=0** (sem DB: 4 SKIPs — os 4 são DB-backed; rodados contra a sessão) |
| go vet | `go vet ./...` | limpo |
| Go full | `go test ./... -count=1` | exit=0, **107 pacotes ok, 0 FAIL** |
| tsc | `../../node_modules/.bin/tsc --noEmit` em apps/web | exit=2, **12 erros = teto pré-existente, ZERO em arquivo tocado** (lista conferida: anuncios/mutations/produto — nenhum do chip) |
| vitest web | `npm run test --workspace @marketplace-central/web` | exit=0, **67 files / 566 tests passed**, zero "No test files" |
| vitest sdk-runtime | `npm run test --workspace @marketplace-central/sdk-runtime` | exit=0, **5 files / 78 tests passed** |
| vitest feature-products | `npm run test --workspace @marketplace-central/feature-products` | exit=0, **1 file / 13 tests passed** |
| integração | `npm run harness:integration` | **status=passed**, migrations_first=72, migrations_second=0, container=session-reuse porta 50265 |
| governança | `pwsh scripts/harness.ps1 -Command governance -BaseSha 4ad36272...` | ver abaixo |

## Governança — lane VERMELHA pré-existente no main (REPORT)

Baseline medido em worktree LIMPO destacado no main tip `4ad36272`: **status=failed, 51
violações** (19 RCFG_UNAPPROVED_READER, 13 GOV_MODULE_DEPENDENCY, 10 RCFG_READER_MISSING,
6 GOV_MODULE_LAYER, 5 RCFG_UNDECLARED_READ, 2 GOV_MODULE_COVERAGE) — nada disso é do chip.

Chip tree: **50 violações = subconjunto estrito do main**. Diff exato (parse por tripla
code/id/path): ONLY-MAIN = `GOV_MODULE_COVERAGE tenant_config` (resolvida pela entry nova em
modules.json); ONLY-CHIP = vazio. `GOV_API_SDK_SPLIT` não dispara — OpenAPI + sdk-runtime no
mesmo commit (F6 @86fc6400). Dois consertos no caminho: `removal_owner` fora do pattern do
schema (GOV_SCHEMA_INVALID, corrigido p/ M-10) e leitor não aprovado de MPC_TEST_DATABASE_URL
em root_test.go (trocado pelo boundary testpostgres.SkipWithoutTarget) — ambos @a2f96fa1.

REPORTs p/ o hub: (1) lane governance vermelha no main com 51 violações pré-existentes —
verde absoluto é inalcançável p/ qualquer chip até saldar; critério usado = zero violação
NOVA vs baseline do main tip, medido por diff de conjunto. (2) lane no checkout do hub
trava >20min varrendo dumps untracked de `docs/design/evidence/ml-api/` (não excluídos
pelo filtro da Policy) — caso novo da classe "governance lane clean worktree".

## Must-fail (F6)

`evidence/F6-must-fail.txt` — guard invalid_q removido: vermelho NOMEIA invalid_q em 2
arquivos; allowed_range re-gateado: vermelho NOMEIA o campo ausente (invalid_ids +
invalid_erp_source). Restaurado, verde. Must-fail do corte (F4): reproduzido pré-wedge em
cache_composed_test.go (`stored-rule read returned [101 202], want the sellable cut [101]`),
registrado em STATE-SHELL-BLOCKED.md.

## FE caller (cond.4)

`evidence/F6-fe-caller-guard.txt` — caminho único até searchCatalogProductFacts com guard
duplo (`enabled: Boolean(query)` + trim/debounce + roteamento p/ lista quando vazio); risco
residual anotado (métodos SDK sem guard, zero consumidores).
