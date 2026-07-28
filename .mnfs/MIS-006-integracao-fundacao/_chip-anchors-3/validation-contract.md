# CHIP-ANCHORS-3 — contrato de validação

Autoridade: `_hub-gate-anchors-2/p6-reconciliation-r1.md`. Cada critério é PASS/FAIL/NOT-PROVEN.
Um critério que você não consegue fazer **totalmente** vira REPORT, e o chip fecha sem ele (R-24) —
não invente uma alegação parcial que soa completa.

Verifique por **string**, nunca por linha. Números de linha neste pack são de `5441fe18`; eles vão
deslizar assim que você editar.

---

## A1 — o lado ERP do `seller_sku` é o CODPROD canônico

`identityAnchorValues`, case `"seller_sku"`, obtém `productValue` do id canônico do produto. `grep`
por `ReferenceCode` dentro dessa função retorna **zero** ocorrências.

PASS exige as duas metades: o grep vazio E o novo caminho lendo o campo canônico.

## A2 — o teste de fixação é alcançável em produção

O caso de teste que assere `side` para `seller_sku` usa um `ProductCandidate` **com**
`InternalProductID` preenchido — o mesmo formato que `generation_service.go:279` deixa passar.

Prove a alcançabilidade, não afirme. O caminho aceito: um teste que arranca do gerador (não da
função interna) com um candidato do formato de produção e observa o motivo emitido.

**Nomeie explicitamente** no EVIDENCE o que aconteceu com
`TestNamedMissingAnchorSitesAreIncomparableWithCorrectSide`: corrigido, substituído, ou deletado. Ele
está verde hoje sobre um fixture impossível. Deixá-lo verde e sem tocar é FAIL deste critério, porque
o próximo leitor vai contá-lo como cobertura.

## A3 — must-fail do A1

Reverta só o `productValue` do case `"seller_sku"` para `product.ReferenceCode` e mostre o teste do
A2 **falhando**, com a assinatura de falha citada (qual valor apareceu no lugar de qual). Depois
restaure e mostre `git diff HEAD` **vazio** no arquivo.

Sem must-fail, A1 está inspecionado, não provado — foi exatamente o F-9 do CHIP-ANCHORS-2.

## A4 — âncora `seller_sku` presente dos dois lados não some mais por comparando errado

Cenário: anúncio com `seller_sku` que resolve o produto ERP, produto ERP com CODPROD (sempre tem).
Antes do fix, `refforn` não-vazio disparava o ramo both-present do A2-R2 e a âncora sumia de
`reasons[]`. Depois do fix, a decisão passa a ser sobre CODPROD contra CODPROD.

Declare qual direção sai agora e por quê. Se a resposta honesta for "continua não emitindo, mas agora
por comparar as coisas certas", diga isso — A2-R2 nega o `FOR` e continua valendo. O que este
critério cobre é o comparando, não a emissão.

## A5 — `vinculados` conta CODPROD com zero à esquerda

Fixture: `erp_import_products.codprod = '00101'`, link resolvido com `internal_product_id = 101`.
Esperado: o produto **DENTRO** de `vinculados`.

## A6 — must-fail do A5

Reverta a junção para a forma antiga, mostre o A5 falhando com o número exato (`Vinculados: 0`, não
"falhou"), restaure, mostre verde.

## A7 — o `DISTINCT` do `vinculados` continua vivo

O CHIP-ANCHORS-2 pagou por este guard (codprod resolvido duas vezes → ainda conta 1). Sua mudança de
junção passa por cima dele. Re-rode e cite. Se a mudança de junção o quebrou, é FAIL — não é
"comportamento aceitável novo".

## A8 — `{id}` malformado responde 4xx nas DUAS rotas

`GET /erp/imports/not-a-uuid` e `GET /erp/imports/not-a-uuid/chain`. Teste de nível de handler serve;
não é preciso servidor.

Se o código de status mudar em relação ao que o OpenAPI declara hoje, o OpenAPI e
`packages/sdk-runtime` mudam **no mesmo commit** (profile §7).

## A9 — o comentário do `marketplace_capability.go` não afirma mais um universal falso

O texto resultante não pode ser refutado por `market/domain/identity_resolver.go:90-92`. Cite o texto
novo inteiro no EVIDENCE — é uma frase, e o critério é sobre a frase.

## A10 — `pending` não-array não derruba o endpoint

Fixture com `cursor -> 'pending'` sendo um objeto ou um escalar. Esperado: resposta válida, não erro.

## A11 — `*comparison.product` não deref nil

Ou o teste chega no sítio com `product == nil`, ou você declara **NOT-PROVEN** e explica por que o
sítio é inalcançável — com o caminho, não com "parece seguro". As duas irmãs checam nil; se este não
precisa, o motivo é uma informação de valor.

## A12 — nenhum dano colateral

- `go build ./...` = 0, `go vet ./...` = 0.
- `go test ./...` verde, com a contagem citada.
- Os dois guards de ordem em `-count=10`, 10/10.
- A suíte de política D-121 verde, com a contagem. **E** com a nota do A2 ao lado: uma das que
  contavam era o fixture impossível.
- `git diff --name-only 5441fe18f64171ef61cb03b51b5bf66e2922e4eb HEAD -- apps contracts packages` —
  **zero** caminhos sob `apps/web/`, **zero** migrations, **zero** `platform/httpx`.

`base_sha` é PISO, não ponto fixo. `main` se move enquanto você trabalha. Antes do fecho, rode o
comando acima contra o `main` de verdade e declare o que apareceu que não é seu.

## A13 — o EVIDENCE aponta, não recopia

O pack cita `p6-reconciliation-r1.md` por caminho e nome de achado (B-01, G1, …). Recopiar os achados
para dentro do EVIDENCE cria uma segunda cópia que vai divergir. R-14.

---

## O que este contrato NÃO cobre

Dito assim porque um marcador que omite sua cobertura é a mesma alegação inflada que R-24 proíbe.

- **Nenhum L2.** Este chip não sobe servidor. Drive ao vivo é do hub.
- **B-08** (deadline da rota) e **G4** (índice) estão fora, com dono nomeado no `chip.md`.
- **O ramo `AGAINST` de A2-R1** continua decisão do operador.
- **Os 15 erros de `tsc` do `apps/web`** não são seus, nem os 3 do CHIP-ANCHORS-2 nem os 12 baseline.
