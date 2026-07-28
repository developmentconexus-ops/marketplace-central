# GATE P6 — CHIP-ANCHORS-3, rodada 7 (F-5)

`tip: 57666417` · `base do diff: tip ATUAL da main` · `diff: r7-code-diff.patch (1194 linhas, .mnfs e HARNESS-PROFILE.md excluídos)`

Rodadas 1–6 cobriram CORR-1..CORR-6 (`GATE-P6.md`). **Esta rodada julga o F-5 e só o F-5.** As seis
anteriores estão descarregadas; reabrir uma delas exige observável novo, não releitura.

## Custódia, verificada pelo hub antes de abrir

- Nenhum arquivo de código apagado contra o tip atual (`git diff --diff-filter=DR` vazio) — a
  armadilha que custou quatro rodadas nesta missão.
- O chip **nunca tocou** `docs/HARNESS-PROFILE.md` (`git log <merge-base>..57666417 --` vazio). As
  97 linhas que o diff bruto mostra são adições do HUB ausentes no branch; o merge não as remove.
- Pack do chip existe e é versionado no tip dele.

## O que o F-5 promete

O painel de `/importacoes/<id>` afirmava um funil (`importados → vinculados → enfileirados`) que o
servidor não apura: são três medidas de **populações independentes** em **duas unidades** (linhas de
arquivo × produtos internos), então `55 · 0 · 55` é estado normal e lia como sistema quebrado.

Entregas: painel renomeado com unidade em cada rótulo e a seta removida; prosa falsa do OpenAPI
DELETADA (R-25) com resíduo escrito para os identificadores que sobrevivem por custo de quebra;
TSDoc do SDK espelhando as `description`s; duas strings de teste sob micro-grant.

## Critérios desta rodada

**F5-1 · BLOQUEANTE.** Nenhuma superfície entregue afirma sequência ou subconjunto entre as três
medidas. Verificar por STRING no tip, não por leitura: `→` no painel, a palavra `cadeia` na cópia,
`processing chain` no contrato. A alegação é TOTAL, então uma única ocorrência sobrevivente a refuta.

**F5-2 · BLOQUEANTE.** Todo rótulo nomeia sua UNIDADE. Este é o defeito real da rodada e ele não é
o desenho: três números anônimos lado a lado continuam convidando à comparação mesmo sem seta.
Mesmo defeito da rodada 5 (`ltrim` fundindo chaves distintas) em superfície diferente — a falsidade
mora na unidade não nomeada.

**F5-3 · BLOQUEANTE.** A dívida de identificador está DECLARADA, não silenciada. `chain` sobrevive
em quatro sítios que quebram cliente — path `:3263`, `operationId: getErpImportChain` `:3266`,
`$ref` `:3282`, schema `:8090`. Silêncio sobre eles seria falsidade nova. O `operationId` é o nome
do método do SDK; **a tabela original do hub o omitia e o chip corrigiu.**

**F5-4.** A mudança no SDK é comment-only, PROVADA e não afirmada: `numstat`, contagem de linhas
adicionadas não-comentário, e sha256 do arquivo com comentários removidos, idênticos entre revisões.

**F5-5.** Nenhuma asserção de terceiro quebrada fora do micro-grant (`ImportChainPanel.test.tsx`
`:173`, `:186`).

## Entrada dos assentos (§11, incluindo o item 4)

1. `r7-code-diff.patch` — o diff contra o tip ATUAL da main;
2. estes critérios verbatim;
3. saídas cruas de lane do chip;
4. **`DRIVE-EVIDENCE-f5.md` — observável MEDIDO PELO HUB.** Três telas dirigidas em navegador real
   contra a stack. Nenhum assento tem navegador; sem isto, F5-1 e F5-2 seriam decididos lendo o
   fonte, que decide o que o código PODE renderizar e nunca o que a tela renderiza.

O artefato declara o próprio limite e o assento deve respeitá-lo: **o drive é PARCIAL** — FE e SDK
do tip, backend da `main`. Todo veredito que dependa de resposta de servidor vale contra a `main`.

## O que o assento NÃO deve fazer

- **Não cobrar o defeito do 5xx deste diff.** O painel travar em `Carregando…` para sempre num 5xx
  está medido no artefato e é **pré-existente da `main`**: `ImportChainPanel.tsx:46-50` e
  `useErpImportChain.ts:23-28` são idênticos nas duas revisões e o diff não toca nenhum dos dois.
- **Não cobrar o `500` do id malformado da `main`.** O `400 invalid_import_id` está NESTE diff; a
  discrepância existe justamente porque ele ainda não mergeou.
- **Não abrir achado de prosa do pack.** O assento recebe o diff, não o pack (§11). Um achado que só
  se sustenta lendo a justificativa do chip é REPORT, nunca BLOQUEANTE.
- **Não tratar o drive como cobertura.** Três telas, um provider, um dia.

## Régua

BLOQUEIA só com observável: file:line, ou string que existe/não existe, ou saída de comando. Achado
de gosto, de escopo especulativo, ou de "poderia ser melhor" é REPORT. Um achado sem observável
nomeado não é achado.
