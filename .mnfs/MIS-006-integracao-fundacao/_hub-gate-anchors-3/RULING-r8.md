# RULING — CHIP-ANCHORS-3 rodada 8 · **BLOQUEADO**

`tip julgado: 1542839a` · gate `GATE-P6-r8.md` + `HUB-LANE-r8.md` · assentos A (Opus) e B (Sol
medium), independentes · **os dois BLOCKED**

Terceira rodada seguida em que os dois lados reprovam. O que muda desta vez é ONDE: as rodadas
anteriores morreram no pack. Esta morre no **código e no write-set**.

---

## A-1 · o comentário do SQL é falso, e é refutado por dois contraexemplos independentes — BLOQUEANTE

`query_repository.go:155`:

> `-- importados and enfileirados count import ROWS; vinculados counts linked internal PRODUCTS.`
> `-- They differ only when one import names the same product twice ('101' and '00101'), and`
> `-- then they differ truthfully.`

Os dois assentos derrubaram a mesma linha **com casos diferentes**, sem contato:

| assento | contraexemplo | por que refuta |
|---|---|---|
| A | fixture `#660-E` **do próprio chip** (`patch:182-206`): codprods `'00101'`,`'ABC'`,`'00ABC'`, um link, dois pendentes — resultado asserido `3 · 1 · 2` (`patch:214`) | três números diferentes, nenhum produto nomeado duas vezes, **asserido pelo próprio commit** |
| B | linhas `101` e `102`, link só para `101` | `importados=2, vinculados=1` sem duplicata |

Um comentário total refutado por uma asserção do mesmo patch. O mesmo commit ainda afirma o
oposto numa superfície de consumidor (`patch:1108`) — P e ¬P numa revisão. E a prosa se repete no
teste de integração (`chain_query_repository_integration_test.go:331`), então o teste documenta a
mesma falsidade que ele mesmo contradiz.

Custo: um mantenedor lê `55 · 0 · 55`, acredita que é artefato de codprod duplicado, e "conserta"
o join que está correto.

## A-2 · garantia de terceiro APAGADA fora do grant — BLOQUEANTE (F8-5)

Achado do assento A. **Medido pelo hub contra a `main`**, que é a base que o F8-5 manda usar:

```
git grep -c "exact EAN ERP seller SKU empty"  main       -> 1
git grep -c "exact EAN ERP seller SKU empty"  1542839a   -> (nada)

git grep -c "id-with-path"                    main       -> 3
git grep -c "id-with-path"                    1542839a   -> (nada)
```

O grant é `ImportChainPanel.test.tsx:173,:186`, `ImportacaoSection.tsx:151`,
`ImportacaoSection.test.tsx:89`. Nada disso está nele.

1. **Caso de tabela inteiro deletado** — `TestNamedMissingAnchorSitesAreIncomparableWithCorrectSide`
   garantia na `main` que ausência de `seller_sku` é atribuída a `side=erp`. A justificativa da
   remoção se apoia em `findProducts:277-283`, que está fora do diff. Se aquele filtro for
   relaxado, nada acende: o caso que pegaria foi removido neste commit.
2. **Garantia trocada** — `TestHandlerGetImportChainMapsServiceErrorsAndUTCResponse` fixava que o
   handler encaminha o valor cru do path (`id-with-path`) ao service. Depois do merge fixa outra
   coisa (encaminha UUID). A garantia velha é deletada em silêncio.
3. **Asserção pré-existente invertida** (`patch:936-951`) — `TestCase3EANAloneYieldsMediaConfirm`
   passa `seller_sku` de INCOMPARABLE para UNAVAILABLE e para de atribuir side. Mudança de
   comportamento em superfície embarcada, declarada só em comentário de teste.

Esta é a classe do `diff-vs-base-de-despacho-esconde-REVERT`, agora em teste em vez de feature.
Write-set não é sugestão: é o que permite dois chips na mesma onda.

## A-3 · o B-01 foi realocado, não morto — BLOQUEANTE

Assento A, e é o pior achado da rodada porque é **o defeito que este chip existe para matar**.

`generation_service.go` (`patch:620-624`) diz que `"produto ERP sem CODPROD"` é alegação que
aquele lado **nunca pode fazer com verdade**, e que ler `refforn` ali cunhou exatamente essa
falsidade. O teste NOVO do mesmo chip (`patch:747`, `:756`) asserta a string idêntica como saída
esperada:

> `wantDetail: "sem CODPROD para corroborar o EAN: o seller_sku do anúncio não casa nenhum produto"`

com entrada onde o produto ERP **tem** CODPROD 101 e o anúncio **tem** seller_sku — o próprio chip
escreve, em `patch:937-939`, que os dois lados carregam valor e apenas discordam.

Só `Side` e `Direction` foram reparados. A frase falsa foi mantida e **recém-certificada por um
teste verde**. B-01 saiu do campo `Side` e entrou no campo `Detail`.

## A-4 · "not a subset of importados" — BLOQUEANTE (F8-3, lacuna declarada que não existe)

`openapi.yaml` (`patch:1124-1128`) e `erpImport.ts` (`patch:1202`). Estruturalmente
`vinculados <= importados` sempre: todo `links.internal_product_id` contado precisa de pelo menos
uma linha de import cujo codprod resolva para ele, e uma linha resolve para exatamente um id. As
três fixtures do próprio chip obedecem o limite (1≤1, 1≤2, 1≤3).

Consumidor avisado de que os dois não têm relação aceita `vinculados=7, importados=3` como leitura
normal e não constrói checagem nenhuma. Estado inalcançável hoje — então uma regressão futura que
o produzisse é lida como "unidades diferentes" em vez de bug.

## A-5 · "Três medidas independentes" — BLOQUEANTE (F8-1, o espelho da falsidade)

`ImportChainPanel.tsx` (`patch:992-995`), cópia visível ao operador. "Não são etapas de um funil"
é verdade; "independentes" não é. `queued_products` faz join de pending contra `import_products`,
logo `enfileirados <= importados` sempre; e `vinculados <= importados` pelo A-4.

O F8-1 proíbe afirmar sequência/subconjunto. Afirmar **ausência de qualquer relação** é a
falsidade espelhada, e o remédio da rodada passou de um lado direto para o outro.

## B-1 (Sol) · "never data loss" — BLOQUEANTE (F8-3)

`openapi.yaml:8123`: *"No consumer is wired yet, so today it only grows — a later drop is
drainage, never data loss."* O mesmo diff trata escrita malformada/parcial do cursor `pending`
como fila vazia: `{"pending":["501"]}` → `{"pending":"501"}` leva `enfileirados` de 1 para 0 sem
dreno nenhum. O contrato publicado chama essa perda de drenagem.

## O que PASSOU

Não é rodada perdida. F8-6 passou no **desenho** pelos dois assentos: o arm mutado falha pelo NOME
ACESSÍVEL (`role link` / name `Ver estado`), não por `data-testid`, 1 vermelho de 29, restauração
conferida por md5. O assento A ainda fez uma verificação que o hub não pediu e que vale a pena
ratificar: o `index 794c21d9..33380892` do log da mutação é o **reverso exato** do
`index 33380892..794c21d9` do patch, provando que a lane rodou contra a versão do chip daquele
arquivo. F8-2 passou nos três rótulos do diff. A metade anti-ambiguidade do F8-1 passou.

## F8-1 é indischargeável como a rodada está empacotada — não é achado, é defeito do GATE

Os dois assentos devolveram a mesma coisa: o F8-1 exige varredura por STRING sobre
`apps/web/src` inteiro, nenhum assento tem shell, e **nenhuma lane do hub contém essa varredura**.
Isso é meu erro, não deles — é a mesma forma do achado da rodada 1 (assento read-only não descarrega
critério de execução). A varredura vira lane do hub no r9, não critério de assento.
