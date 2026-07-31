# CHIP-VENDAVEL — estado congelado por falha de ferramenta (shell morto)

**Data:** 2026-07-31
**Chip:** CHIP-VENDAVEL (MIS-006, pós-missão)
**Worktree:** `.claude/worktrees/chip-vendavel`, branch `worktree-chip-vendavel`
**BASE-SHA:** `554788d576d04b21719f4a17e4702dd0f0aff4e1` · tip da main (base do gate): `774d75aa`

## 1. O bloqueio

Toda invocação do shell falha **antes** de rodar o comando, com assinatura constante:

```
/usr/bin/bash: -c: line 71: syntax error: unexpected end of file from `{' command on line 70
```

Reproduzido em: `pwd`, `true`, `echo probe`, `git --version`; com `run_in_background`;
com sandbox desligado; via tool Monitor; e **dentro de um subagente** (processo novo,
mesmo erro byte-idêntico). Logo: não é o meu comando, não é o sandbox, não é a sessão.

Sem shell não há `go build`, `go test`, `go vet`, `git`, nem lanes. **Nada abaixo do
último run bem-sucedido está verificado por execução.**

## 2. Diagnóstico

Classe: bloco não fechado em arquivo de inicialização carregado pelo bash (mesma classe
do thread Arch #168953 — lá era um `if` sem `fi`).

Verificado e **descartado**:

- `~/.claude/shell-snapshots/snapshot-bash-*.sh` — 100 arquivos, todos com ≤44 linhas
  (20 ou 44). O mais novo (`snapshot-bash-1785506322265-vzq686.sh`) foi lido inteiro:
  `rg` abre em 7 e fecha em 18, `fi` em 19, `pkill` abre em 22 e fecha em 43, `export PATH`
  em 44. **Balanceado.** Nenhum snapshot chega à linha 70.
- `~/.bashrc` — não existe.
- `<scratchpad>/testdb-env.sh` — 8 linhas.

Conclusão: o script de 71 linhas é **composto em tempo de chamada** pelo harness
(snapshot de 44 linhas + ~27 linhas de wrapper). O `{` órfão está na parte gerada, que
**não existe em disco** — não há arquivo para eu corrigir com Edit. O conserto é do
operador: reiniciar o app/sessão do Claude Code (regenera o wrapper).

Nada em `apps/` foi alterado por esse diagnóstico.

## 3. O que está COMMITADO (verde na época, execução verificada)

- `0d3004d` F1 — `erp_import/adapters/postgres/chain_query_repository_integration_test.go`
- `b9f858d` F2 — `product_links/application/generation_service.go` + `_test.go` + fixtures FE
- `219fb39` F3 — `validation-contract.md`
- `7659bac` — `BATCH-PLAN.md:950`
- evidências: `evidence/F1-s13-clock.txt`, `evidence/F2-unproved-copy.txt`

## 4. O que está NÃO-COMMITADO e NÃO VERIFICADO

Último `go test ./...` bem-sucedido terminou com exatamente 4 falhas, todas de rota legada
em `tests/unit/catalog_handler_test.go` (verbatim):

```
--- FAIL: TestCatalogHandlerGetReturnsProducts (0.00s)
    catalog_handler_test.go:129: expected 200, got 503: {"error":"source_unavailable"}
--- FAIL: TestCatalogHandlerRejectsPost (0.00s)
    catalog_handler_test.go:154: expected JSON error body on 405: invalid character 'M' looking for beginning of value
--- FAIL: TestCatalogSearchEndpoint (0.00s)
    catalog_handler_test.go:187: expected 200, got 503: {"error":"source_unavailable"}
--- FAIL: TestCatalogSearchEndpointRequiresQuery (0.00s)
    catalog_handler_test.go:208: expected 400, got 503
```

As 4 foram tratadas por edição **manual** depois disso, e o shell morreu antes de compilar.
Estado atual da árvore, por item:

### F4 — dobra da porta (compilou verde antes do wedge)
Porta única `internalreadports.Source = Reader + CatalogPageReader`. `CatalogAssortmentReader`
e seu comentário falso deletados. 10 pacotes nomeados pelo compilador ajustados, incluindo
`internal/composition/root.go:440` e `internal/platform/migrate/runner_test.go` (71→72, dois sítios).
Evidência escrita: `evidence/F4-port-fold.txt`, `evidence/F4-compiler-census.txt`.

### F4b — decorator único no cache (NÃO verificado)
`internal_read/adapters/cache/cache.go`: `CatalogPageReader`/`NewCatalogPageReader`
desexportados para `catalogPageReader`/`newCatalogPageReader`; asserções colapsadas em
`var _ internalreadports.Source = Reader{}`.
`internal_read/observability/timing.go`: asserção passa a `ports.Source`.

### F5 — deleção da rota legada não-cortada (NÃO verificado)
Deletados: o fork em `catalog/transport/http_handler.go:Register`, `handleLegacyProducts`,
`handleLegacySearch`; `ListProducts`/`SearchProducts` de `catalog/application/service.go` e
`canonical_service.go`; os métodos correspondentes de `catalog/ports/canonical_reader.go` e
`catalog/ports/repository.go`; 4 walks + 4 métodos de `UnavailableReader` em
`catalog/adapters/internalread/reader.go`.

### F6 — guard de contrato `q` obrigatório (NÃO verificado, e é mudança de comportamento)
`handleSearch` passa a rejeitar `q` em branco com 400 `invalid_q`. OpenAPI declara
`q: required: true` (`contracts/api/marketplace-central.openapi.yaml:462-472`); o handler
paginado nunca cobrou — o check morava no handler legado deletado no F5 — então `q` vazio
descia como padrão vazio e voltava o catálogo inteiro vestido de resultado de busca.

### Testes migrados à mão (NÃO compilados)
- `tests/unit/cache_composed_test.go` — reescrito sobre a cadeia real
  (cache → TimingReader → routing.Reader → Service → RouteClassMux). Novo
  `TestComposedCatalogCutTravelsWireToSource`. Must-fail JÁ FOI feito e reproduzido verbatim:
  `cache_composed_test.go:263: stored-rule read returned [101 202], want the sellable cut [101]`.
- `tests/unit/catalog_handler_test.go` — novo `catalogHandlerPageReader`; 405 agora assere
  `Allow: GET` (o 405 é do router, não do handler); blank-q agora assere `error == invalid_q`.
- `tests/unit/catalog_service_test.go`, `tests/unit/enrichment_precedence_test.go` — chamadas
  reapontadas para `ListProductsByIDs`; testes das rotas mortas deletados com comentário.
- `catalog/adapters/internalread/reader_test.go` — `TestCanonicalWalksAskForTheTenantRule`
  deletado, com comentário apontando para o teste composto.

## 5. Pendências conhecidas, ainda não feitas

1. `go vet ./...` + `go test ./...` completos (nada acima está verificado por execução).
2. Must-fail do guard `invalid_q` (remover, mostrar a falha nomeada, restaurar).
3. `internal/composition/catalog_routes_test.go` — o cabeçalho do arquivo e as mensagens
   descrevem o fork `if h.PageReader == nil && !hasRouteClasses` que o F5 **deletou**.
   É prosa falsa sobre o repo (R-25) e a última checagem (`CATALOG_METHOD_NOT_ALLOWED`)
   ficou vácua: a string só sobrevive em `handleTaxonomy`
   (`catalog/transport/http_handler.go:477`), que nada tem com essa rota. Compila, mas
   precisa ser reescrito para o mundo de uma rota só — **não editei** para não empilhar
   mudança não verificável.
4. Commit do 71→72 em `runner_test.go`.
5. `EVIDENCE.md` escrito DURANTE o run do `go test`, antes de qualquer evento CLOSED.
6. Escada P5 completa, incluindo o run explícito de
   `./internal/modules/tenant_config/... ./internal/composition/...` que nenhuma lane
   descobre (dívida B-6), e `harness:governance -- -BaseSha 774d75aa`.
7. Gate duplo P6 em árvore congelada (Opus frio + GPT-5.6 Sol medium), ambos lendo o diff
   contra `774d75aa`.
8. `pg-session-down` no fecho (A-25 passo 5). Container da sessão:
   `mpc-pg-session-1e6bac12`, porta 55841.

## 6. Alargamento de escopo a declarar ao hub

O contrato despachado cobria F1–F3. F4, F4b, F5 e F6 são alargamento, sob a ordem
permanente do operador de eliminar redundância e local maximum em um batch. F6 em
particular **muda comportamento de API** (400 onde antes vinha 200 com o catálogo inteiro)
e conserta uma violação real do OpenAPI. Nenhum deles pode ser fechado sem a escada.

## 7. Achados de harness desta sessão

- **B-8 (novo):** o shell do harness pode morrer para a sessão inteira, subagentes
  incluídos, com erro na própria preâmbulo gerada. Sem shell não há degradação graciosa:
  toda a escada de verificação cai junto. O wrapper não está em disco, então não há
  conserto local — só reiniciar o app. Vale um probe de shell no bootstrap do chip.
- **B-6** (lane de integração varre só as 5 primeiras linhas dos `_test.go`) e **B-7**
  (`vitest.config.ts` inclui feature-products por nome exato) permanecem como ratificados.
- Classe já registrada: critério pode ser **vácuo contra o TIPO** — "nunca chama X" não
  reprova se o tipo do colaborador não declara X.
