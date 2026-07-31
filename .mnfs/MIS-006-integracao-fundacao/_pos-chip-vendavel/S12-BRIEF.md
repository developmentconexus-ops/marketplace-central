# S12-CATALOGO-FE — brief corrigido (medido, não alegado)

Base: `9c139f45` (tree clean). Worktree `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel`, branch `worktree-chip-vendavel`.

O slice card (BATCH-PLAN.md:801-830) foi escrito antes de S9/S11 existirem. Este brief é
a versão medida. Onde os dois divergirem, **este manda**; o `failing_test_first` e as
strings de tela do card ficam verbatim.

## Escopo

FE puro. **Nada de Go, nada de OpenAPI, nada de `packages/sdk-runtime/`.** Medido:
`CatalogPageOptions` (sdk-runtime/src/index.ts:245-250) e `CatalogSearchPageOptions`
(:258-265) **já declaram `include_all?: boolean`**, e o handler já o interpreta
(`catalogPolicy`, catalog/transport/http_handler.go:254-271, chamado nos dois assentos:
:131 lista, :234 busca). O contrato está pronto; falta só a tela usar.

## Write set (fechado)

- `packages/feature-products/src/catalogQueries.ts`
- `packages/feature-products/src/CatalogPage.tsx`
- `packages/feature-products/src/CatalogPage.test.tsx`
- `apps/web/src/app/AppRouter.test.tsx`

Fora do write set, e o card estava errado ao listá-los: `packages/web-query/src/index.ts`
e `index.test.ts` — os construtores de chave já aceitam params arbitrários
(`catalogQueryKeys.facts(params)` index.ts:30, `search(q, params)` :35-38), então
`include_all` entra pelo objeto de params sem tocar no seam compartilhado.
`apps/web/src/app/AppRouter.tsx` também sai: :23-29 **já** monta o componente real com
`client` + `erpSource`; só o teste está faltando.

**Proibido criar arquivo de teste novo em `packages/feature-products/`.**
`apps/web/vitest.config.ts` inclui esse pacote por **nome de arquivo exato**
(`"../../packages/feature-products/src/CatalogPage.test.tsx"`), não por glob — um
segundo arquivo de teste ali roda em lane nenhuma e o verde é fantasma. Todo teste de
S12 vai dentro de `CatalogPage.test.tsx`.

## Teste que falha primeiro (verbatim do card)

```
opens with Vendáveis 2 de 4 and ver todos never mutates tenant config
```
em `packages/feature-products/src/CatalogPage.test.tsx`.

### A armadilha desse teste — leia antes de escrevê-lo

`CatalogQueriesClient` (catalogQueries.ts:9-16) hoje declara só
`listCatalogProductFacts` e `searchCatalogProductFacts`. Se você montar a página com um
client que **não tem** `setSellableAssortment`, a asserção "never mutates tenant config"
é estruturalmente incapaz de falhar — é teatro. O client passado no teste **tem que
carregar um spy `setSellableAssortment`** (e `getCatalogAssortmentCounts`), para que a
asserção `expect(setSpy).not.toHaveBeenCalled()` tenha um mundo em que reprova.
Prove isso: no must-fail, chame `client.setSellableAssortment` no handler do "Ver todos"
e mostre o teste VERMELHO nomeando essa asserção.

## Mudanças, por arquivo

### `catalogQueries.ts`

1. `CatalogQueriesClient` ganha `include_all?: boolean` nos options dos dois métodos, e
   ganha `getCatalogAssortmentCounts: () => Promise<CatalogAssortmentCounts>` (a página
   precisa satisfazer o `SellableAssortmentClient` do hook de contagem).
2. `useCatalogFactsQuery` e `useCatalogSearchQuery` recebem `includeAll: boolean`.
   - **Chave:** `include_all` entra no objeto de params das duas chaves. Sem isso,
     filtrado e "ver todos" colidem na mesma entrada de cache e a volta ao filtrado
     serve as linhas erradas sem nenhuma requisição.
   - **Wire:** quando `includeAll` é false, **não envie o parâmetro** — nem
     `include_all=false`. O handler trata ausente e `false` como idênticos
     (http_handler.go:256-268), e o teste existente em `CatalogPage.test.tsx:81`
     asserta `toContain("/catalog/products/search?q=CUBA")`, que quebra se outro
     parâmetro entrar antes do `q`.

### `CatalogPage.tsx`

1. Estado `includeAll` (default `false`), botão "Ver todos" ligando, e o caminho de
   volta ao filtrado. O botão **nunca** chama `setSellableAssortment` — é escape de
   tela, não mudança de regra do tenant.
2. Chip de contagem via o hook que **já existe** (S11):
   `useCatalogAssortmentCountsQuery(client, erpSource)` de
   `@marketplace-central/web-query` (activeSource.ts:89-99, chave
   `catalogQueryKeys.counts(erpSource)`, `enabled: Boolean(erpSource)`).
   **Não escreva um segundo hook de contagem.**
   - Texto exato quando há dado: `Vendáveis 2 de 4` (i.e. `Vendáveis {sellable_count} de {total_count}`).
   - **Sem dado (fonte não resolvida, ou a query falhou): não renderize o chip.**
     `Vendáveis 0 de 0` é fato operacional desconhecido virando zero — proibido pelo
     profile §7 / ADR-17. Ausência é a única leitura honesta aqui.
3. Badge `Sem estoque` na coluna Estoque:
   - mostra **sse** `quantity != null && quantity <= 0`;
   - `quantity == null` continua `— (motivo)` via `factQuantity` (:30-32) e **nunca**
     vira `Sem estoque` — desconhecido ≠ zero.
   - A string já foi prometida ao operador em
     `apps/web/src/pages/integracoes/IntegracoesPage.tsx:394` ("produtos sem saldo
     continuam no sortimento com o aviso Sem estoque"); tem que bater byte a byte.
4. **`refresh()` (:65-83) reconstrói as chaves à mão.** Se você não puser `include_all`
   ali, o botão "Atualizar" em modo "ver todos" refetcha uma chave morta e a tela não
   se move — falha silenciosa que nenhum teste de render pega. Corrija e teste.

### `AppRouter.test.tsx`

O critério do card ("prova que /catalogo monta o componente real com o client e a
partição por fonte") está **não atendido**: :29 faz
`vi.mock("@marketplace-central/feature-products", () => ({ CatalogPage: () => <div>Catalog route</div> }))`,
um mock cego que não vê prop nenhuma. Não desmocke (puxaria a página inteira e a rede
para a lane do router, e as outras rotas continuam mockadas — assimetria). Em vez disso
faça o mock **observar as props**: renderize o `erpSource` recebido e a presença do
client, e asserte no teste da rota. Isso prova a fiação da partição de cache sem montar
a página real.

## Testes exigidos (todos em `CatalogPage.test.tsx`, salvo o do router)

1. `opens with Vendáveis 2 de 4 and ver todos never mutates tenant config` (verbatim).
2. Modo filtrado **não** manda `include_all` na URL (lista e busca).
3. "Ver todos" refetcha com `include_all=true`.
4. Voltar ao filtrado usa chave distinta — prove pelo número de fetches, não por
   inspeção de chave: filtrado → todos → filtrado com o mesmo `QueryClient` não pode
   servir a página errada, e a volta ao filtrado não dispara requisição nova (cache
   próprio) enquanto o dado está fresco.
5. Produto com `quantity: 0` retornado em modo "ver todos" continua na tabela com badge
   `Sem estoque`; produto com `quantity: null` mostra `— (motivo)` e **não** o badge.
6. Sem `erpSource`: chip ausente, `getCatalogAssortmentCounts` nunca chamado.
7. "Atualizar" em modo "ver todos" dispara fetch (chave viva).
8. Router: `/catalogo` passa `erpSource` resolvido ao componente.

Os testes existentes (6, linhas 49-154) continuam passando **sem edição**. Se algum
exigir edição, pare e reporte — é sinal de que a mudança quebrou contrato de tela.

## Must-fails obrigatórios (RED verbatim no relatório)

- MF-1: "Ver todos" chama `setSellableAssortment` → teste 1 vermelho nomeando a asserção.
- MF-2: `include_all` fora da chave de query → teste 4 vermelho.
- MF-3: badge disparando em `quantity == null` → teste 5 vermelho.
- MF-4: `refresh()` sem `include_all` na chave → teste 7 vermelho.

Restaure o código byte a byte depois de cada must-fail.

## Lanes (forma verificada neste chip — NÃO use a do card)

```powershell
Set-Location 'C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-vendavel'
npm run test --workspace @marketplace-central/web -- --run ../../packages/feature-products/src/CatalogPage.test.tsx src/app/AppRouter.test.tsx
```

`npx --no-install vitest` (forma do card) já produziu verde vacuoso neste repo. Não use.

Reporte a contagem por teste: conte as linhas `--- ` da saída. **SKIP tem que ser 0** —
um pulado silencioso é o modo de falha que este chip já pegou duas vezes.

Rode também o tsc da raiz: `npm run typecheck` (ou o alvo equivalente da raiz), porque
alargar `CatalogQueriesClient` pode quebrar outro chamador.

## Regras duras

- **Não commite. Não `git add`. Não `git push`.** O orquestrador revisa e commita.
- Não suba servidor, não toque `:8080`, não leia nem crie `.env*`.
- Não instale dependência. Precisa = pare e reporte.
- Copy em pt-BR; nenhum marcador de dev na tela.
- Sem abstração especulativa, sem comentário narrando o óbvio; comentário só onde
  registra POR QUE (o padrão dos arquivos vizinhos).

## Relatório final

- Diff por arquivo.
- Saída RED e GREEN verbatim de cada must-fail e da lane final, com contagem
  passed/failed/skipped.
- Qualquer divergência entre este brief e o que o código exigiu, nomeada.
