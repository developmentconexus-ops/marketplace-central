# GATE P6 — CHIP-ANCHORS-3, rodada 8 (F-5)

`tip: c3acf62b` · `base do diff: tip ATUAL da main` · `diff: r8-code-diff.patch (1.227 linhas, sha256 362fa10c…ef957fd1)`

Rodadas 1–6 cobriram CORR-1..CORR-6 (`GATE-P6.md`) e estão descarregadas. A rodada 7 fechou
**SEM AGREEMENT** (`VERDICTS-round7.md`): um bloqueante do chip e **dois critérios meus mal
escritos**, que foram a causa da divergência dos assentos. Esta rodada julga o F-5 com F5-3, F5-4 e
F5-5 **reescritos pelo hub** — renumerados F8-3, F8-4, F8-5.

## Custódia, verificada pelo hub antes de abrir

- Nenhum arquivo de código apagado ou renomeado contra o tip ATUAL da main
  (`git diff --diff-filter=DR main c3acf62b -- apps contracts packages` vazio).
- O chip **nunca tocou** `docs/HARNESS-PROFILE.md` (`git log <merge-base>..c3acf62b --` vazio).
- 15 arquivos tocados, todos dentro do write-set declarado do chip.

## O que mudou desde a rodada 7

`57666417..c3acf62b`, um commit. O bloqueante da rodada 7: `Ver cadeia` em
`ImportacaoSection.tsx:151` — a porta de entrada prometendo a sequência que o destino nega, medida
pelo hub em navegador. Reparo: `Ver cadeia` → **`Ver estado`** em `:151` e no nome acessível de
`:89`. Mais a declaração da dívida no contrato e no TSDoc, e a tabela de delta por classe.

## Critérios desta rodada

**F8-1 · BLOQUEANTE — alegação TOTAL, e a superfície agora está nomeada.**
Nenhuma superfície em **`apps/web/src` inteiro** que leve o operador a `/importacoes/<id>` ou que
renderize o painel nomeia o destino como cadeia, nem afirma sequência ou subconjunto entre as três
medidas. A alegação é TOTAL sobre a JORNADA, não sobre o diretório do chip: **uma única ocorrência
sobrevivente a refuta.** Verificar por STRING, não por leitura.

Segunda metade, e ela é do reparo: **a troca não pode comprar falsidade com ambiguidade.** Em
`:145` existe um controle vizinho `Ver detalhes`/`Ocultar detalhes` que EXPANDE inline; `:151`
NAVEGA. Checar papel e nome acessível dos dois — dois controles adjacentes com o mesmo nome
acessível, um expandindo e outro navegando, seria defeito novo desta rodada.

**F8-2 · BLOQUEANTE — inalterado.** Todo rótulo nomeia sua UNIDADE. A falsidade desta onda não mora
no desenho, mora na unidade não nomeada: três números anônimos lado a lado continuam convidando à
comparação mesmo sem seta.

**F8-3 · BLOQUEANTE — REESCRITO. A dívida vai declarada onde o CONSUMIDOR lê.**
Na rodada 7 escrevi "a dívida está DECLARADA" sem nomear a superfície, e os assentos divergiram por
isso: um recusou supor que estava no pack (correto), o outro passou lendo a asserção do gate como se
fosse a declaração (circular). **Pack não conta e gate não conta.** A dívida tem que estar em
superfície que quem consome a API lê: `description` no OpenAPI **e** TSDoc no `sdk-runtime`.

E a declaração tem que ser VERDADEIRA e COMPLETA sobre quais sítios sobrevivem. Dado verificável
para o assento decidir sozinho: o `$ref` `#/components/schemas/ErpImportChain` é **referência ao
nome do schema** — renomear o schema renomeia o `$ref` —, então não é sítio independente dos outros
três. Se o assento julgar que é, tem que dizer por qual string.

**F8-4 · BLOQUEANTE — REESCRITO. O par de revisões está nomeado, e a alegação total morreu.**
Na rodada 7 pedi "comment-only" sem dizer identidade **entre o quê**, e os dois assentos acertaram
perguntas diferentes: um mediu o commit de documentação (comment-only VERDADEIRO), o outro mediu
contra a `main` (mudança de tipo). A pergunta certa é contra a **`main`**, porque é isso que mergeia.

O delta contra a `main` vai declarado **POR CLASSE** — linhas que afetam TIPO, linhas de comentário,
linhas de prosa — com **cada linha de tipo NOMEADA**. Uma alegação total de "comment-only" sobre o
chip é falsa e não pode aparecer. R-24 aplicado ao hub tanto quanto ao chip: alegação total que o
código não sustenta vira enumeração.

**F8-5 · REESCRITO. A unidade é o MERGE, então a base é a `main` e o grant é a UNIÃO.**
Na rodada 7 pedi delta declarado e entreguei diff cumulativo — mesmo defeito de unidade. Nenhuma
asserção de terceiro alterada fora do grant, medido **contra a `main`**, com o grant sendo a união
de todas as rodadas: `ImportChainPanel.test.tsx:173`, `:186`, `ImportacaoSection.tsx:151`,
`ImportacaoSection.test.tsx:89`.

**F8-6 · BLOQUEANTE — o braço must-fail do rótulo.**
Asserção de rótulo não se certifica sozinha. Prova exigida na lane crua: rótulo revertido no fonte
com o teste intacto → VERMELHO **pelo NOME ACESSÍVEL**, não por `data-testid`; restaurado → VERDE.
Um teste que falha por `testid` sobreviveria à falsidade e por isso não descarrega nada aqui.

## Entrada dos assentos (§11, incluindo o item 4)

1. `r8-code-diff.patch` — o diff contra o tip ATUAL da main;
2. estes critérios verbatim;
3. saídas cruas de lane do chip;
4. **`DRIVE-EVIDENCE-f5.md`** — três telas dirigidas pelo hub em navegador real, e a captura de
   `/importacoes` que produziu o bloqueante da rodada 7.

**Limite declarado, e ele corta contra o hub:** o drive é PARCIAL (FE e SDK do tip, backend da
`main`), e **o reparo desta rodada NÃO foi dirigido em navegador**. `Ver estado` está verificado por
STRING no tip e nada mais. O drive do reparo acontece no QA de merge, que é onde a doutrina o põe —
só QA passa milestone. Nenhum assento deve tratar o reparo como medido em tela.

## O que o assento NÃO deve fazer

- **Não cobrar o defeito do 5xx deste diff.** O painel travar em `Carregando…` para sempre num 5xx é
  **pré-existente da `main`**: `ImportChainPanel.tsx:46-50` e `useErpImportChain.ts:23-28` são
  idênticos nas duas revisões e o diff não toca nenhum dos dois.
- **Não cobrar o `500` do id malformado da `main`.** O `400 invalid_import_id` está NESTE diff; a
  discrepância existe justamente porque ele ainda não mergeou.
- **Não abrir achado de prosa do pack.** O assento recebe o diff, não o pack (§11). Achado que só se
  sustenta lendo a justificativa do chip é REPORT, nunca BLOQUEANTE.
- **Não reabrir CORR-1..CORR-6** sem observável NOVO. Releitura não reabre.
- **Não tratar o drive como cobertura.** Três telas, um provider, um dia.

## Régua

BLOQUEIA só com observável: `file:line`, ou string que existe/não existe, ou saída de comando.
Achado de gosto, de escopo especulativo, ou de "poderia ser melhor" é REPORT. Achado sem observável
nomeado não é achado.
