# GATE P6 — CHIP-ANCHORS-3, round 5 (forma nova: os assentos leem o DIFF, não o pack)

`alvo` = `main` @ `678a6d51` · `chip` = `chip/anchors-3` @ `f8363839`
`entrada dos assentos` = `code-diff.patch` (957 linhas, gerado por
`git diff main chip/anchors-3 -- apps contracts packages`) + os critérios abaixo + a saída crua da
lane do assento executor do hub.

Ratificado em `f19793d7` / `678a6d51`: o pack do chip **não é entrada de revisor**. Quatro rounds
anteriores gastaram dois assentos cada lendo 9.730 linhas de pack para produzir, no round 4, três
achados bloqueantes sobre um diff que o próprio assento leitor mediu como *"40 insertions, 12
deletions, all `//` text; no behavior changed"*. Esta rodada corrige a MIRA, não o rigor.

## Régua desta rodada (`3f8560b1`)

- **BLOQUEIA** só quando o achado nomeia um **observável errado**: comportamento, segurança, dado,
  ou contrato publicado.
- **REPORT** para prosa, contagem, metadado, citação, formatação, higiene. Registra, conserta, não
  segura merge e não abre rodada.
- Discriminador, respondido ANTES de atribuir a classe: *deixa o achado exatamente como está, sobe
  o código — o que um usuário, um operador, um chamador ou uma linha gravada faz de diferente?*
  Sem resposta = REPORT.
- Reclassificar é do hub. O assento reporta tudo que achar, na severidade que acreditar.

## Verificação de custódia, feita pelo hub antes de abrir a rodada

O branch foi trazido pra frente (`git merge main`, merge `65fbfe7b`) depois que o hub mediu que ele
revertia a tela `/importacoes` inteira. Re-medido pelo hub, independente do chip:

```
git diff --numstat main chip/anchors-3 -- apps/web            → 0 linhas (intocado)
git diff --numstat main chip/anchors-3 -- .../cmd/mlprobe     → 0 linhas (intocado)
git branch --contains main                                    → chip/anchors-3 listado
```

Nenhuma linha do delta tem `0` na coluna de inserções. O revert está morto. **Isto não é critério
desta rodada** — está aqui para o assento não gastar tempo re-derivando.

## Critérios — verbatim do `chip.md`, e são a única coisa que o chip prometeu

### CORR-1 — BLOQUEANTE. `seller_sku` estava lendo `refforn`

`identityAnchorValues`, case `"seller_sku"`, decidia se o lado ERP tem a âncora lendo
`product.ReferenceCode`, que é **`refforn`** (referência do fornecedor), enquanto o CODPROD canônico
está no mesmo struct em `InternalProductID` e não era usado.

Duas saídas erradas que o defeito produzia:
1. produto ERP com `refforn` vazio → `INCOMPARABLE` + `side=erp` com detail *"sem CODPROD para
   corroborar o EAN"* — o motivo afirma que o produto ERP não tem CODPROD, impossível, porque
   CODPROD é parte da chave primária de `erp_import_products`;
2. produto ERP **com** `refforn` + anúncio com `seller_sku` → classificador conclui "presente dos
   dois lados", devolve `emit=false`, e a âncora **some de `reasons[]`** — o sumiço silencioso que
   o F-02 existe para matar.

**Exigido:** o valor do lado ERP para `seller_sku` sai do CODPROD canônico
(`InternalProductID` / `ProductID`), nunca de `ReferenceCode`. A alternativa (manter `ReferenceCode`
e reescrever os details) está **negada por escrito** no brief.

**Armadilha nomeada no brief, e é onde o assento deve olhar:** o teste que existia
(`generation_service_test.go:417`) usava `&internalreaddomain.ProductCandidate{}` — candidato SEM
`InternalProductID`, que `generation_service.go:279` filtra antes de qualquer candidato existir.
Ele asseria `side == erp` sobre um fixture impossível em produção: **codificava o defeito em vez de
guardar contra ele**. O fixture novo tem que carregar id canônico, senão a asserção continua
inalcançável e trocou-se um teste vazio por outro. *Verifique isso no diff.*

### CORR-2 — BLOQUEANTE. `vinculados` perdia CODPROD com zero à esquerda

`ON links.internal_product_id::text = products.codprod`. Os dois lados vêm de pipelines diferentes:
`erp_import_products.codprod` guarda a string crua; `product_links.internal_product_id` é
`strconv.ParseInt` numa coluna `integer`. `IsValidCodprod` **aceita** zero à esquerda. Então
`'00101'` vira link `101` e a junção compara `'101' = '00101'` → falso: o endpoint reporta um número
**menor que a verdade, sem erro nenhum**.

**Exigido:** junção compara a mesma representação canônica dos dois lados, e existe fixture de
regressão — `codprod = '00101'`, link `internal_product_id = 101`, esperado DENTRO de `vinculados`.

O lado da FILA é imune (o enqueuer empurra `row.Codprod` cru), mas as duas colunas são lidas na
mesma tela como decomposição de uma população — se só um lado for canonicalizado, o operador lê a
diferença entre dois números como fila travada.

### CORR-3 — `{id}` malformado respondia 500

`GET /erp/imports/not-a-uuid/chain` → `PathValue("id")` virava `ImportID` sem validação → SQL contra
coluna `uuid` → Postgres levanta `22P02`, que não é `pgx.ErrNoRows` → `internal_error`. O chamador
não distingue "você mandou lixo" de "nós quebramos".

**Exigido:** valida o id antes da query, **nas duas rotas** (`{id}` e `{id}/chain`), num lugar só. Se
mudou código de status documentado, OpenAPI + `packages/sdk-runtime` no MESMO commit (profile §7).

### CORR-4 — comentário de doutrina afirmava um universal que o repo refuta

`marketplace_capability.go:33-34` dizia que `refforn` *"answers `no` for every provider present and
future"*, enquanto `identity_resolver.go:90-92` anexa uma âncora `"refforn"` comparando o `RefForn`
do ERP contra atributo `MODEL` do marketplace, e esse resolver está ligado.

**Exigido:** a remoção da lista continua certa; a **justificativa** é que estava falsa. R-25 — a
metade falsa é deletada ou estreitada, nunca anotada. Estreitar ao vocabulário que o arquivo governa
ou nomear a exceção do módulo `market`.

### CORR-6 — dois cercos baratos

- `jsonb_array_elements_text` sobre um `pending` não-array levanta erro em query-time e derruba o
  endpoint inteiro para o tenant; `COALESCE` só defende de `NULL`. Guardar com
  `jsonb_typeof(state.cursor -> 'pending') = 'array'`.
- `product := *comparison.product` — deref incondicional num ponteiro que as duas funções irmãs
  checam contra nil.

## Fora de escopo — não é achado bloqueante se cair aqui

- **CORR-5**: frase do EVIDENCE do CHIP-ANCHORS-2. É do hub, já feita em `bcab8269`.
- **B-02** (`QueueRow.tsx`) e tudo em `apps/web`: dono é o CHIP-VINC-NEUTRO. O delta deste chip não
  toca `apps/web` — se o assento achar `apps/web` no diff, isso é o achado.
- Os **5 universais falsos pré-existentes** fora do delta: do hub, na fila.
- `generation_integration_test.go` está atrás de `//go:build integration` — excluído em COMPILAÇÃO,
  não pulado em runtime. Lane Postgres real é do assento executor do hub, não do assento leitor.

## O que o assento leitor NÃO deve fazer

- Não pedir o pack. Se um achado de CÓDIGO depender da derivação escrita do chip, peça uma seção
  **nomeada e limitada**; nunca a árvore.
- Não emitir achado sobre prosa, contagem, citação ou consistência de documento. Não há documento
  nesta entrada — só o diff.
- Não citar `file:line` de arquivo que não esteja no diff. Um assento já FABRICOU uma citação para
  arquivo que não existe em árvore nenhuma (`_hub-gate-import-chain`, registrado no §11). Se o
  arquivo não está no patch, ele não está nesta rodada.

## Assento EXECUTOR do hub — saída crua, medida fora do chip

Checkout: `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\hub-exec-anchors3`
Tip medido: `f8363839` (detached). `GOCACHE` absoluto, `GOMODCACHE` absoluto.

```
BUILD_EXIT=0
VET_EXIT=0
go test ./internal/modules/product_links/... ./internal/modules/erp_import/... -count=1 -v
  FUNCOES (linhas '=== RUN' sem '/')  = 199
  TOTAL com subtestes ('=== RUN')     = 304
  --- PASS = 304   --- FAIL = 0   --- SKIP = 0
  13 pacotes `ok`, 3 `[no test files]`, 0 `FAIL`
```

Reproduz o número que o chip reportou, medido por instrumento do hub em checkout do hub.
Contagem gerada por `grep -c`, não lida em tail.

## Tip do chip avançou para `617fd41f` DEPOIS do despacho dos assentos — a rodada continua válida

```
git diff --numstat f8363839 617fd41f -- apps contracts packages   → vazio
git diff --numstat f8363839 617fd41f                              → 56/12 em EVIDENCE.md, só
sha256 code-diff.patch (gerado em f8363839)                       = 559c15f5…4da546f3
sha256 git diff main 617fd41f -- apps contracts packages          = 559c15f5…4da546f3
```

Patch byte-idêntico nos dois tips. É a propriedade que a forma nova compra: churn de pack não
invalida rodada de gate, porque pack não é entrada de assento.
