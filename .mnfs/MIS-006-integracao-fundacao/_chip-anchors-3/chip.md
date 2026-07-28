# CHIP-ANCHORS-3 — corretivos do gate P6 do CHIP-ANCHORS-2

```yaml
chip: CHIP-ANCHORS-3
branch: chip/anchors-3
base_sha: 5441fe18f64171ef61cb03b51b5bf66e2922e4eb
wave: 1.5 (paralelo aos dois chips de FE — write-sets disjuntos, backend vs apps/web)
authority: .mnfs/MIS-006-integracao-fundacao/_hub-gate-anchors-2/p6-reconciliation-r1.md
```

## De onde este chip vem

O CHIP-ANCHORS-2 foi mergeado em `main @dbdcdfb1` e **depois** passou pelo dual gate real que o
operador exigiu — Opus e `gpt-5.6-sol` medium, adversariais, sobre input congelado. **Os dois lados
reprovaram**, por evidências diferentes. A reconciliação está em
`_hub-gate-anchors-2/p6-reconciliation-r1.md` e é a **autoridade** deste chip; os verdicts brutos
estão ao lado (`p6-sol-gate-r1.md`, `p6-opus-gate-r1.md`). Leia a reconciliação inteira antes de
tocar em código — ela contém as rulings de severidade que já foram decididas, e reabrir uma delas é
retrabalho.

O merge fica registrado como **mergeado e NÃO aprovado por gate**. Este chip é o que fecha isso.

Não recopie os achados para o seu EVIDENCE. Aponte para eles (R-14).

## Escopo — cinco corretivos

### CORR-1 — BLOQUEANTE. `seller_sku` está lendo `refforn`, e emite motivo falso

`identityAnchorValues`, case `"seller_sku"`
([generation_service.go:734](../../../apps/server_core/internal/modules/product_links/application/generation_service.go:734)),
decide se o lado ERP tem a âncora lendo `product.ReferenceCode`. Esse campo é **`refforn`** — a
referência do fornecedor:

- `ReferenceCode: copyTrimmed(row.Referencia)` ([reader.go:418](../../../apps/server_core/internal/modules/erp_import/adapters/internalread/reader.go:418))
- `MirrorProduct{… Referencia: row.Refforn …}` ([source_contract_test.go:349](../../../apps/server_core/internal/modules/erp_import/adapters/internalread/source_contract_test.go:349))

O CODPROD canônico está no MESMO struct, em `InternalProductID`, e não é usado. Em todo o resto do
sistema o contraparte ERP de um `seller_sku` de provider é o CODPROD — o próprio matcher compara
`row.CodigoProduto == *sku` ([reader.go:456](../../../apps/server_core/internal/modules/erp_import/adapters/internalread/reader.go:456)).

Duas saídas erradas, as duas novas neste merge:

1. Produto ERP com `refforn` vazio → `INCOMPARABLE` + `side=erp` colado no detail
   `"sem CODPROD para corroborar o EAN"`. O motivo afirma que o produto ERP **não tem CODPROD**.
   Impossível: CODPROD é parte da chave primária de `erp_import_products`.
2. Produto ERP **com** `refforn` + anúncio com `seller_sku` → o classificador conclui "presente dos
   dois lados", devolve `emit=false`, e **a âncora some de `reasons[]`**. É o sumiço silencioso que
   o F-02 existe para matar, reintroduzido por comparar um SKU de marketplace contra uma referência
   de fornecedor.

**Correção:** o valor do lado ERP para `seller_sku` sai do CODPROD canônico do produto
(`InternalProductID` / `ProductID`), nunca de `ReferenceCode`.

**A alternativa está NEGADA, e por escrito.** Manter `ReferenceCode` e reescrever os details que
dizem "CODPROD" recolocaria `refforn` como o lado ERP de uma comparação cross-side — na mesma missão
cujo F-01 o tirou de lá (D-A). O par como está no código não pode ser verdadeiro dos dois lados.

**O teste que existe hoje não guarda nada.** `generation_service_test.go:417` usa
`product: &internalreaddomain.ProductCandidate{}` — candidato sem `InternalProductID`, que
`generation_service.go:279` filtra antes de qualquer candidato existir. Ele assere `side == erp`
sobre um fixture impossível em produção: **codifica o defeito em vez de guardar contra ele**. O
fixture novo carrega um id canônico, senão a asserção continua inalcançável e você trocou um teste
vazio por outro.

### CORR-2 — BLOQUEANTE. `vinculados` perde CODPROD com zero à esquerda

[query_repository.go:89](../../../apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go:89):

```sql
ON links.internal_product_id::text = products.codprod
```

Os dois lados vêm de pipelines diferentes. `erp_import_products.codprod` guarda a string crua;
`product_links.internal_product_id` é `strconv.ParseInt(...)` ([reader.go:162](../../../apps/server_core/internal/modules/erp_import/adapters/internalread/reader.go:162))
numa coluna `integer`. `IsValidCodprod` **aceita** zero à esquerda
([seller_sku.go:20-34](../../../apps/server_core/internal/modules/internal_read/domain/seller_sku.go:20)).
Então `'00101'` vira link `101`, e a junção compara `'101' = '00101'` → falso. O endpoint reporta um
número **menor que a verdade, sem erro nenhum**.

O hub decidiu isto como BLOQUEANTE contra a leitura de um dos reviewers: o valor é digitado por
humano na planilha do upload xlsx, e ausência da forma no dado de HOJE não prova que o caminho está
fechado.

**Correção:** a junção compara a mesma representação canônica dos dois lados. Fixture de regressão
obrigatório: `codprod = '00101'` no protocolo, link resolvido `internal_product_id = 101`, esperado
DENTRO de `vinculados`.

Nota que economiza uma investigação: o lado da FILA é imune — o enqueuer empurra `row.Codprod` cru,
então `queued_products` já é string contra string.

### CORR-3 — `{id}` malformado responde 500

`GET /erp/imports/not-a-uuid/chain` → `PathValue("id")` vira `ImportID` sem validação → o SQL liga
numa coluna `uuid` → Postgres levanta `22P02`, que não é `pgx.ErrNoRows` → cai no `internal_error`.
O chamador não distingue "você mandou lixo" de "nós quebramos".

Não é regressão deste chip: `GET /erp/imports/{id}`
([http_handler.go:111](../../../apps/server_core/internal/modules/erp_import/transport/http_handler.go:111))
tem a forma idêntica e é anterior. `type ImportID string` não tem validador e não existe guarda de
uuid inválido em `internal/` inteiro.

**Correção:** valide o id antes da query, **nas duas rotas**, num lugar só. Isto descarrega o F-8 do
CHIP-ANCHORS-2, que estava arquivado como "não verificado por ninguém".

### CORR-4 — comentário de doutrina afirma um universal que o próprio repo refuta

[marketplace_capability.go:33-34](../../../apps/server_core/internal/modules/connectors/ports/marketplace_capability.go:33)
diz que `refforn` "answers `no` for every provider present and future". Mas
[identity_resolver.go:90-92](../../../apps/server_core/internal/modules/market/domain/identity_resolver.go:90)
anexa uma âncora `"refforn"` comparando o `RefForn` do ERP contra um atributo `MODEL` vindo do
marketplace, e esse resolver está ligado.

A remoção da lista continua certa; a **justificativa** é que está falsa. R-25: a metade falsa é
deletada ou estreitada, nunca anotada. Estreite para o vocabulário que o arquivo governa, ou nomeie
a exceção do módulo `market`.

### CORR-6 — dois cercos baratos

- `jsonb_array_elements_text` sobre um `pending` não-array levanta erro em query-time e derruba o
  endpoint inteiro para o tenant. `COALESCE` só defende de `NULL`. Guarde com
  `jsonb_typeof(state.cursor -> 'pending') = 'array'`.
- `product := *comparison.product` em [generation_service.go:490](../../../apps/server_core/internal/modules/product_links/application/generation_service.go:490)
  deref incondicional, num ponteiro que as duas funções irmãs checam contra nil.

## Fora de escopo, nomeado

- **CORR-5** (a frase do EVIDENCE do CHIP-ANCHORS-2 sobre o ramo default) é **do hub**. Não edite
  pack de chip anterior.
- **B-02** — `QueueRow.tsx:159` derruba todo motivo `INCOMPARABLE`. É `apps/web`, dono é o
  CHIP-VINC-NEUTRO.
- **B-08** — `RouteClassMux` indexa a classe pelo padrão nu e o `Handle` recebe o padrão com método,
  então o upload de xlsx declara 120s e recebe 15s. Real, pré-existente, **e não é seu**: mexer em
  `platform/httpx` colide com toda rota do repo. Vai virar chip próprio. Se seu diff tocar
  `route_deadline.go`, é escopo vazando.
- **G4** — índice para `(tenant_id, state, internal_product_id)`. Precisa de `EXPLAIN` em escala de
  produção, que nenhum reviewer pôde medir e nenhum alegou ter medido.
- **O ramo `AGAINST` de A2-R1** continua decisão do operador e continua intocado.
- **Migrations: nenhuma.** Se seu diff criar uma, é escopo vazando.

## Matriz de propriedade

| Eixo | CHIP-ANCHORS-3 |
|---|---|
| Migração | **nenhuma** |
| DB shape | nenhuma |
| Módulo Go | `product_links/application`, `erp_import/adapters/postgres`, `erp_import/transport`, `connectors/ports` — **não** `platform/httpx`, **não** `market/` |
| `root.go` | nenhum |
| Contrato/SDK | só se o CORR-3 mudar código de status documentado; aí OpenAPI + `packages/sdk-runtime` no MESMO commit (profile §7) |
| FE | **nenhum**. `apps/web/` está fora de escopo, inclusive os 3 erros de `tsc` — dono é o CHIP-VINC-NEUTRO |

## Ladder

- L0: `go build ./...`, `go vet ./...`, lane de governança.
- L1: `go test ./...`, mais o guard de ordem em `-count=10`.
- L2: **do hub**. Este chip não sobe servidor, não liga em `:8080`, não lê `.env*`.

`cd apps/server_core` antes de qualquer comando `go` — rodar da raiz polui `git add -A` com
`.gomodcache`. Não passe `-mod=mod`: o repo está em workspace mode e o comando falha.

## Regras de fecho

- Fecho com `AGREEMENT — P6 discharged` e o ledger de discharge ao lado. A linha `P6-DUAL-GATE:` é
  **do hub** — o chip não a escreve.
- Nenhum push. O merge é do hub.
- **Reviewer adversarial é obrigatório e tem que persistir o artefato.** Os três reviewers do
  CHIP-ANCHORS-2 deixaram `.output` de 0 bytes e por isso os verdicts originais são irrecuperáveis.
  Streamar não é persistir. Se o artefato do reviewer sair vazio, o pass não aconteceu.
