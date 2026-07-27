# CHIP-ANCHORS-2 — contrato de validação

Cada critério é verificável por leitura ou por execução. Um critério que só pode ser atestado por
opinião não está aqui, de propósito.

## C1 — `refforn` fora do vocabulário cross-side (F-01)

`IdentityAnchorRefforn` não aparece mais em `knownIdentityAnchors`. Provado por **string**:
`git grep -n "IdentityAnchorRefforn"` na tip, com a saída colada. As ocorrências que
sobrarem — se sobrarem — são nomeadas uma a uma e cada uma tem de estar fora da lista de âncoras
declaradas ao gerador.

`refforn` continua existindo do lado ERP. Um grep em `erp_import_products` / `products_mirror`
mostrando o campo intacto é parte da prova: a remoção é da lista, não do dado.

## C2 — âncora declarada cujo VALOR falta de algum lado não some calada (F-02)

> Reescrito pelo hub em A2-R1/A2-R2. O título anterior era *"nenhuma âncora declarada some calada"*,
> e alegava mais do que a tabela abaixo verifica: o caso "declarada e presente dos dois lados"
> continua não emitindo nada, porque o classificador compara PRESENÇA e nunca estabeleceu que os
> dois lados concordam. Estreitar a alegação é o remédio de R-24; alargar o código para emitir um
> `FOR` não verificado seria o defeito que R-24 existe para impedir.

Teste de unidade sobre a geração, com um provider que **declara** uma âncora, cobrindo os quatro
casos da tabela do pack:

| Caso | Esperado |
|---|---|
| provider não declara (`Supplied == false`) | motivo `UNAVAILABLE`, **sem** `side` |
| declara, valor do anúncio vazio | `INCOMPARABLE` + `side = provider` |
| declara e tem valor, produto ERP sem o dado | `INCOMPARABLE` + `side = erp` |
| declara, faltando dos dois lados | `INCOMPARABLE` + `side = both` |

O teste assere que a âncora **está presente** em `reasons[]` nos quatro casos. O defeito de origem
era ausência silenciosa, então um teste que só verifica o valor de um motivo que existe não prova
nada sobre ele.

## C3 — `side` só existe em `INCOMPARABLE`

Serialização JSON de um motivo `FOR`, um `AGAINST` e um `UNAVAILABLE`: a chave `side` **não aparece**
(`omitempty`), verificada no JSON produzido, não no struct. Um `INCOMPARABLE` no mesmo teste mostra
a chave presente.

## C4 — política de auto-aprovação intacta (D-121)

O critério mais importante deste chip, porque é o que ele pode quebrar sem querer.

Teste que prova que `INCOMPARABLE` **não é evidência contra**: um par que hoje auto-aprova por
CODPROD+EAN concordantes continua auto-aprovando depois de ganhar um ou mais motivos `INCOMPARABLE`
de outras âncoras. E um par de âncora única continua indo para CONFIRMAÇÃO, não sendo rebaixado a
REVIEW por causa de um `INCOMPARABLE`.

Prova adicional por **string**: nenhum somatório de confiança/score lê `INCOMPARABLE`. Cole o grep.

## C5 — ordenação estável e display de equivalência (F-03)

Duas partes, ambas obrigatórias.

**(a) Determinismo.** Guard rodando `-count=10`, com a saída colada. Um guard de ordem que roda uma
vez não distingue estável de sortudo.

**(b) Load-bearing provado.** Reverta `SortStableFunc` para `SortFunc`, rode o mesmo guard em
`-count=10`, mostre-o **falhando**, restaure. As duas saídas vão no EVIDENCE. Um guard que passa nas
duas versões do código não é guard.

**(c) Display.** `50cm` × `500MM` produz `50cm ≡ 500MM`. `50cm` × `50cm` produz `50cm`. Os dois
casos no teste — sem o segundo, a mudança gera `50cm ≡ 50cm` em produção e ninguém percebe.

## C6 — chain-read responde os três números (F-04)

Teste de integração com dados semeados: um protocolo com N produtos, M deles com `product_links` em
`state = 'resolved'`, K deles presentes em `sync_state.cursor -> 'pending'` para `entity = 'market'`.
A resposta traz exatamente `N`, `M`, `K`.

O fixture precisa de **pelo menos duas instalações** com filas diferentes e um `codprod` presente nas
duas, para provar que `enfileirados` agrega com `DISTINCT` em vez de somar duplicado. Uma amostra de
uma instalação só não distingue as duas implementações.

E precisa de pelo menos um `codprod` do protocolo que **não** esteja na fila, senão `K = N` e o
filtro nunca é exercido.

## C7 — semântica de `enfileirados` declarada no contrato

O `description` do campo no OpenAPI diz que é a fila **atual**, e que o número cai conforme o
scheduler consome. `queue_read_at` está na resposta e no schema. Provado por string no
`marketplace-central.openapi.yaml`.

Isto é critério e não detalhe de redação: o campo cai sozinho com o tempo, e um consumidor que o
leia como histórico de importação vai reportar perda de dado que não houve.

## C8 — casos negativos do endpoint

- protocolo inexistente → `404`, não `200` com zeros;
- protocolo real sem nenhuma linha em `sync_state` → `200` com `enfileirados: 0`, **não** `null`.
  A pergunta tem resposta conhecida; ADR-17 é para o que não se sabe, e aqui se sabe.

## C9 — OpenAPI + SDK no mesmo commit do comportamento

`git show --stat` do commit de F-02 e do commit de F-04, mostrando o arquivo Go, o
`marketplace-central.openapi.yaml` e o `packages/sdk-runtime/src/*` juntos (profile §7). Dois
commits separados, um com o comportamento e outro com o contrato, reprovam este critério mesmo que a
tip final esteja correta.

## C10 — write-set disjunto, `apps/web` intocado

`git diff --name-only <base_sha> <tip>` não contém **nenhum** caminho sob `apps/web/`, e nenhuma
migration nova. Saída colada.

O `tsc` vermelho de `apps/web` causado por F-02 é **declarado** no EVIDENCE, com a mensagem do
compilador e o nome do chip da onda 2 que fecha (CHIP-VINC-NEUTRO). Declarar não é consertar: um
diff que toca `QueueRow.tsx` reprova este critério mesmo que o conserto esteja certo.

## C11 — `UNAVAILABLE` volta a ter um sentido só (A2-R1)

Acrescentado pelo hub ao conceder o grant de A2-R1. Leia a ruling antes deste critério.

**(a) Reclassificação.** Nos sítios nomeados em A2-R1 — localizados por STRING — um motivo cujo
valor falta de algum lado sai como `INCOMPARABLE` com o `side` da tabela da ruling, não como
`UNAVAILABLE`. O `detail` de cada sítio permanece.

**(b) Ausência provada por string.** Depois da reclassificação, nenhum `UNAVAILABLE` sobra em
`generation_service.go` fora do caminho `!anchor.Supplied`. Cole o grep. Se sobrar algum, ele é
nomeado um a um com a razão — um grep com exceções explicadas é evidência; um grep com exceções
silenciosas é o defeito que R-24 descreve.

**(c) Ramo excluído intocado.** Teste que prova que o ramo de A2-R1 — anúncio com valor e produto
ERP com valor **não-vazio e diferente** — continua produzindo exatamente o motivo de hoje. Sem este
teste a exclusão é intenção, não fato. O ramo vai ao EVIDENCE como FINDING, com evidência por
string, para o hub levar ao operador; **não** vira `INCOMPARABLE` (seria trocar uma afirmação errada
por outra) e **não** vira `AGAINST` (mexeria em D-121, que é decisão do operador).

**(d) C4 revalidado.** A política de auto-aprovação D-121 continua intacta DEPOIS da
reclassificação, não só antes. Rode C4 de novo no tip final e cole a saída — este grant toca os
caminhos `ExactSKU` e `ExactEAN`, que são a fila de confirmação.

---

## Condições de merge

1. L0 verde: `go build ./...`, `go vet ./...`, lane de governança.
2. L1 verde: `go test ./...` mais o guard de F-03 em `-count=10`, com a prova de must-fail (C5b).
3. C1–C11 com evidência citada por **string**, nunca por linha.
4. Fecho com `AGREEMENT — P6 discharged` e o ledger de discharge ao lado. A linha
   `P6-DUAL-GATE:` é do hub — o chip não a escreve.
5. Nenhum push. O merge é do hub.

O live-drive é do hub: este chip não sobe servidor, não liga em `:8080` e não lê `.env*`.
