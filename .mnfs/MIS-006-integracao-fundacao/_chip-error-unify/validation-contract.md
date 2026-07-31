# CHIP-ERROR-UNIFY — Validation Contract

Critérios estáveis. Todo criterio exige evidência medida (comando + saída) em
`evidence/` — alegação sem observável = não cumprido (R-24).

## VC-1 — Um único produtor de erro no backend
Grep-censo pós: zero produtor local de erro em `internal/modules/*/transport` e adapters
HTTP; todo sítio de erro chama `platform/apierror`. **A POPULAÇÃO é a medida do chip, não
os 20 do pack** (emenda 2026-07-31 por ESCALATION do chip: censo real = 34 funções — 20
terminais + 14 wrappers delegantes — + 1 sítio inline sem helper nomeado em
catalog/transport/http_handler.go:440 + o fallback plain-text de
platform/httpx/json.go:12). Evidência: `evidence/EU-census-backend.txt` reconciliando
POPULAÇÃO vs EXTRAÇÃO com os DOIS números impressos (WriteJSON call-sites ≥400,
`http.Error(`, `w.WriteHeader(4|5)`), e contagem 0 de produtor fora de
`platform/apierror` + testes. O "20" do pack fica visivelmente como subconjunto
histórico, nunca como o todo.

## VC-2 — Shape universal no fio
Para CADA um dos 12 módulos transport: 1 teste de unidade pinando o JSON COMPLETO do
envelope de um erro representativo (não substring — lição ORA-937). Campos ad hoc
migrados p/ `details` (`allowed_range`, `limit`, `column`, `import_id`, `protocol`)
cobertos por teste onde existiam.

## VC-3 — Recover central
Middleware no mux raiz: teste que faz handler entrar em panic e recebe
`500 {"error":{"code":"internal_error",...}}` sem vazar mensagem do panic no body; log
registra. Controle negativo: rota sã não passa pelo shape de panic.

## VC-4 — OpenAPI coerente
`CatalogPageErrorResponse` flat removido; especializações por domínio via base
`ErrorResponse` (padrão MutationErrorResponse); TODAS as respostas de erro das rotas
referenciam schema envelope; `erp_source` declarado no counts (enum
`xlsx|catalogo_cliente`). Evidência: diff do yaml + grep zero `type: string` como shape
de campo `error` raiz em schemas de erro.

## VC-5 — SDK: throw path único tipado
`MarketplaceCentralClientError` classe Error real (`instanceof` funciona); `code` uniões
por domínio + global; `isApiError`/`hasCode` exportados com narrowing (teste de tipo
compile-time + runtime); parse com validação runtime — teste alimenta shape desconhecido
e recebe `internal_error` com `details.raw` preservado (nunca `undefined`). Zero
`as ErrorResponse` cego no caminho de erro (grep).

## VC-6 — FE: padrão único de match
Censo pós: zero comparação de literal de código de erro fora de `hasCode()`/união tipada
em TODOS os pacotes FE (`apps/web`, `packages/feature-*`, `packages/sdk-runtime`
consumidores) — varredura total, não só arquivos tocados (lição guard-parcial).
Evidência: `evidence/EU-census-fe.txt` com o grep e 0 achados fora do padrão.

## VC-7 — Must-fail duplo nomeado
(a) Vermelho que NOMEIA teste falhando quando um transport emite shape velho (injetar
defeito temporário, capturar `failure_token=test=`, reverter); (b) grep-gate do VC-1/VC-6
demonstrado sustentando (introduzir helper local temporário → gate acusa → reverter).

## VC-8 — Escada completa verde
go build/vet/test (107+ packages) · integração hermética `status=passed` + run_id ·
tsc = teto 12 pré-existentes EXATOS (lista comparada) · vitest lanes explícitas web +
sdk-runtime + feature-* tocados · governança set-diff zero violação nova vs baseline 54.
Emenda 2026-07-31 (achado do chip): lane nova `tsconfig.test.json` do sdk-runtime LIGADA
(antes NENHUM teste do SDK era type-checked — `@ts-expect-error` provando união era
checado por nada). Baseline da lane nova = os 7 erros pré-existentes NOMEADOS
(activeSource.test.ts diretiva morta; index.test.ts OrderRead ×2 +
frete_desconhecido; listings-signals.test.ts median/min_valid/max_valid ×2) — zero erro
NOVO além deles; correção dos 7 = fila do hub, fora do escopo do chip.

## VC-9 — QA live-drive (HUB, pós-merge — só isto passa o chip)
3 telas com erro REAL forçado (ex.: `erp_source` inválido via URL, panic route de teste
removida antes do merge, 404 de recurso): UI mostra estado de erro correto (toast/mensagem),
nunca tela branca/undefined; console zero erro não-relacionado; network mostra envelope
novo. Before/after se defeito achado (padrão A-26: banco velho, sem limpeza prévia).
