# CHIP-ERROR-UNIFY — Context Pack (hub, 2026-07-31)

Decisão do operador 2026-07-29 (ratificada 2026-07-31 em brainstorm): **padrão único de erro
HTTP, solução global máxima, FE com padrão de match tipado — nunca por string ad hoc**.

## Identidade

- Chip: CHIP-ERROR-UNIFY. Branch: `worktree-chip-error-unify`, worktree
  `.claude/worktrees/chip-error-unify` (isolado; `npm ci` na raiz do worktree, junction
  banido — profile §3).
- Base de despacho: main tip no momento do spawn (registrar SHA no primeiro commit).
- Hub: sessão "MIS-006 hub" neste repo — eventos via ccd `send_message` (procure sessão
  por título contendo "hub" no cwd marketplace-central; fallback: escreva evento em
  `.mnfs/MIS-006-integracao-fundacao/_chip-error-unify/events.md` e commite no branch).
- Chips falam com hub SÓ por eventos: `CLOSED`/`BLOCKED`/`ESCALATION`/`REQUEST`/
  `SPLIT-REQUEST`/`COMMITTED`/`ACK`.

## Doença (censo medido pelo hub, 2026-07-31 — não é alegação, é grep)

**4 famílias de shape, 20 helpers, 12 módulos transport.** Todos passam por
`platform/httpx.WriteJSON` (json.go:9) mas montam payload independente.

Famílias:
- A — envelope `{"error":{"code","message","details"}}`: maioria (14 helpers).
- B — flat `{"error":"<code>"}`: catalog (`writeCatalogPageError` http_handler.go:398,
  + `allowed_range` opcional); erp_import (http_handler.go:195) e tenant_config
  (http_handler.go:165) usam `{"error":"<code>","detail":"..."}` (`detail`, não `message`).
- C — sem `details`: connectors (`writeConnectorsError` http_handler.go:42).
- D — ad hoc: erp_import `writeImportError` (http_handler.go:162) emite
  `{"error":"missing_required_column","column":...}` e
  `{"error":"duplicate_file","import_id":...,"protocol":...}`; profitability
  (http_handler.go:176) caso especial 422 `{"error":"limit_exceeded","limit":200,"message":...}`.

Helpers (file:line — todos morrem, substituídos pelo pacote único):
catalog http_handler.go:75,398,467 · erp_import http_handler.go:162,195 ·
tenant_config http_handler.go:165 · profitability http_handler.go:176 ·
classifications http_handler.go:33 · inventory http_handler.go:45 ·
listings http_handler.go:333,340 · integrations http_handler.go:48 +
run_read_handler.go:107,123 · connectors http_handler.go:42 +
adapters/melhorenvio/oauth.go:249 · market http_handler.go:128,141,150 +
collection_handler.go:49 · pricing http_handler.go:36 + calc_handler.go:62 ·
product_links http_handler.go:112 · orders http_handler.go:359,688,692,709 +
sankhya_linkage_handler.go:174 · dashboard http_handler.go:64 ·
marketplaces http_handler.go:42 · mutations errors.go:37,44 + query_handler.go:326.

SDK (`packages/sdk-runtime/src/index.ts`): throw path único faz
`throw { status, error: (data as ErrorResponse).error }` (getJson :1747, putJson :1818,
deleteJson :1831, postJson :1839, postVoid :1852) — cast cego. **Mistyping confirmado**:
rotas catalog (:1887-1900) e erp_import lançam `error` = string em runtime, tipo diz
objeto → `err.error.code` = `undefined` silencioso.

OpenAPI: `ErrorResponse` :7553 (107 refs), `CatalogPageErrorResponse` :3984 (flat, 13 refs),
`ConnectorErrorResponse` :7584, `MutationErrorResponse` :8129.

**Zero recover central**: nenhum `recover()` em transport/platform — panic em handler vaza
sem shape.

Gap conhecido: `GET /catalog/products/counts` lê `erp_source`
(catalog/transport/http_handler.go:53 via requestContext) mas o bloco OpenAPI :429-461 não
declara o parâmetro (irmão `/catalog/products` declara :370-380).

## Design ratificado (operador aprovou 2026-07-31)

1. **Fio**: UM shape universal `{"error":{"code","message","details"}}`.
   `CatalogPageErrorResponse` flat MORRE; vira especialização do `ErrorResponse` (padrão
   `MutationErrorResponse`: base + enum de códigos por domínio). Campos ad hoc migram p/
   `details`: `allowed_range`, `limit`, `column`, `import_id`, `protocol`. Códigos
   normalizados onde inconsistentes (rename seguro: match vira tipado, tsc pega).
2. **Backend**: pacote novo `apps/server_core/internal/platform/apierror`:
   `Write(w, status, code, message string, details map[string]any)` (+ variação p/ erro
   tipado se natural). Os 20 helpers morrem; 12 transports chamam o único. Middleware
   recover central no mux raiz: panic → 500 `{"error":{"code":"internal_error",...}}` + log
   (nunca vazar mensagem do panic no body).
3. **SDK**: `MarketplaceCentralClientError` vira classe Error real
   (`{status, code, message, details}`); `code` = uniões tipadas por domínio + união global
   geradas do OpenAPI. Type guards exportados: `isApiError(e)` e `hasCode(e, "x")` com
   narrowing genérico. Parse com validação runtime: shape desconhecido → erro
   `internal_error` com raw preservado em `details.raw` — NUNCA `undefined` silencioso.
   Path de throw único; zero `as`-cast cego.
4. **FE** (`apps/web` + packages feature-*): padrão ÚNICO de match: `hasCode()`/switch
   exaustivo na união. Censo de todo match ad hoc (comparação string, `.error.code`
   direto, `err.error ===`) → migrar TUDO. Grep-gate no censo final: zero literal de
   código de erro comparado fora do padrão.
5. **Brinde**: declarar `erp_source` como query param no OpenAPI do counts (enum
   `xlsx|catalogo_cliente`, igual ao irmão :370-380).

## Regras que MORDEM neste chip (aprendidas em sangue)

- **R-24**: alegação TOTAL no pack/censo só com prova por string/grep; senão vira REPORT.
- **Substring-pin não prova statement** (ORA-937, VENDAVEL): teste que pina pedaço de
  payload não prova shape completo — pine o JSON inteiro do envelope por rota amostrada.
- **Must-fail NOMEIA o teste**: verde de integração só é evidência depois do vermelho
  nomear o teste que falha (`failure_token=test=`); pulado-vs-verde é byte-idêntico.
- **OpenAPI + sdk-runtime juntos** no mesmo commit (AGENTS.md). FE junto quando toca fio.
- **Guard parcial sob frase total é pior que guard nenhum** (VINC-NEUTRO): se o grep-gate
  do FE promete "zero match por string", o gate tem que varrer TODOS os pacotes FE, não só
  os tocados.
- **tsc**: teto = 12 erros pré-existentes EXATOS (compare lista, não contagem). `npx
  --no-install` só depois de `npm ci` real na raiz do worktree (pass vacuoso documentado).
- **Unknown operational fact nunca vira zero/default** — parse de erro desconhecido
  preserva raw, não inventa shape.

## Escada de verificação (worktree, sempre `cd apps/server_core` p/ Go)

- Go: `GOCACHE=.gocache GOFLAGS=-modcacherw go build ./... && go vet ./... && go test ./...`
  (gomodcache: aqueça `.gomodcache` senão HPG_MIGRATION_FAILED falso na integração;
  grep-check gomodcache=0 pré-commit).
- Integração hermética: `pwsh -NoProfile -File scripts/harness.ps1 integration-test`
  (retry CREATE DATABASE no primeiro boot; `failure_token=test=` p/ must-fail).
- FE: `npm ci` na raiz → `npx tsc --noEmit` por pacote tocado + vitest lanes EXPLÍCITAS
  (web, sdk-runtime, feature-* tocados).
- Governança: worktree limpo detached, BaseSha 40-hex do main tip; critério = set-diff
  por tripla vs baseline, zero violação nova (baseline atual 54).
- pg-session p/ testes que precisam de Postgres:
  `pwsh -NoProfile -File scripts/harness.ps1 pg-session-up` (derrubar no fim).

## Interdições

- Chip NUNCA sobe servidor/:8080/.env/dev stack — precisa = `REQUEST` ao hub.
- Nunca `git push`. Nunca reset/revert/stash/clean/-D. Dep nova = `REQUEST`.
- Nunca ler/criar/printar `.env*`; nunca printar valor de env var (nome só).
- Skills de execução mnfs legadas = denylist (profile §10). Codex morto até Aug 5 — não
  despachar contra a parede; gate é do hub (waiver §12).
- Evidência em `.mnfs/MIS-006-integracao-fundacao/_chip-error-unify/` — não escrito = não
  aconteceu. Prosa com forma de veredito no pack = violação (R-24).

## Gate e fechamento

P6 = dual gate do hub (waiver §12: Opus cold + sonnet adversarial; Sol retro pós-Aug-5).
Só QA (live-drive do hub pós-merge) passa o chip. Evento `CLOSED` com: SHA final, escada
verde medida, censo pós (0 helpers locais, 0 match string FE), must-fail nomeado, diff do
OpenAPI/SDK/FE coerente.

## ENCERRAMENTO (hub, 2026-07-31)

Chip MERGED @6bc30c4c; QA VC-9 PASS @9f0fad33. Gate P6: round 0 Opus REJECT (artefato
VC-4) + sonnet APPROVE; remediação 42808712/bce555a8 (incluiu conserto de falso negativo
do próprio gate: 400 sem `content:` no OpenAPI); round 1 Opus APPROVE. Escada pós-merge
verde nas duas lanes; dev stack rebuildado; worktree/branch removidos.

**FILA DE CHIPS AVULSOS ENCERRADA** (decisão do operador 2026-07-31, "finalizar tudo"):
MIS-006 closed desde 28/07; VENDAVEL e ERROR-UNIFY foram os dois últimos chips avulsos.
Trabalho restante (ProdutoPage 404, router fallback, união SDK completa, inline→allOf,
7 tsc sdk test-lane, PII scrub ml-api, candidates, D-7..D-10) = insumo de MISSÃO NOVA
com escopo e critério de encerramento próprios. Nenhum despacho novo deste hub.
