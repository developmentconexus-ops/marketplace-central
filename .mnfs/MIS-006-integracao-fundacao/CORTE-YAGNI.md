# CORTE — decisão do operador, 2026-07-28

Ordem: YAGNI. Fechar a missão com o que muda o funcionamento do software. Tudo que não muda,
sai.

## O que estava errado

A R-24 ("alegação total ou vira report") pegou coisa real no começo — o B-01, o guard parcial sob
frase total. Depois virou gate de **todo comentário do repo**. Os critérios que o hub despachou
(V9-1..V9-7, F8-1..F8-6) são quase todos sobre verdade de prosa, então os assentos devolvem
prosa. Custo: 8 rodadas no ANCHORS-3, 9 no VINC-NEUTRO, com uma rodada inteira gasta em
injetividade de uma função cuja única entrada viva é `ProviderCode: "mercado_livre"` hardcoded em
`connectors/adapters/mercado_livre/capability_adapter.go:81`.

Erro de dosagem do hub: comentário tratado com o mesmo peso de tela.

## R-24, novo alcance (vigente já)

Dispara em:
- **string que o operador lê na tela**;
- **contrato publicado** (OpenAPI, tipos do SDK).

Não dispara em comentário de código. Comentário errado vira revisão normal ou task, nunca rodada
de gate.

## Estado real da missão

`M-01`..`M-05` **mergeados**: sync state + scheduler, mirror/active-source, adapter xlsx, adapter
Sankhya, auto-vínculo. A integração funciona. `M-06` é tela, `M-07` é descoberta.

## Vereditos revistos

### CHIP-VINC-NEUTRO r9 — **ACEITO**

Todos os achados da rodada 9 são prosa; o próprio ruling do hub registra que nenhum é visível ao
operador. Sob o novo alcance, nenhum bloqueia. Merge como está. Duas correções de 2 linhas
(`wireFixtures.ts:79` contradizendo `:455`; dois títulos de teste alegando universal aberto)
entram no mesmo commit se o chip quiser; não são condição.

O `550` vs `551` e o `tsc 12`: irrelevantes, não voltam.

### CHIP-ANCHORS-3 r8 — **ACEITO COM 3 CONDIÇÕES**, sem novo gate

Só o que chega ao operador ou apaga guarda. Verificado por medição do hub, não por alegação:

1. **String falsa na tela.** `"sem CODPROD para corroborar o EAN: …"` é emitida em caso onde o
   produto TEM codprod. `reason.detail` renderiza em `QueueRow.tsx:59`, `:95` e
   `VinculoDrawer.tsx:117`. É o B-01 realocado de `Side` para `Detail`.
2. **Cópia falsa na tela.** `ImportChainPanel.tsx:47` — `Três medidas independentes` — quando o
   SQL garante `enfileirados ≤ importados` e `vinculados ≤ importados`.
3. **Testes de terceiro apagados fora do grant.** Medido contra a `main`:
   `"exact EAN ERP seller SKU empty"` main=1 / chip=0; `"id-with-path"` main=3 / chip=0. Restaurar
   ou justificar por comportamento, não por comentário.

Fora disso: nada. Sem critério de prosa, sem assento, sem rodada de gate. O hub verifica os três
por string e mergeia.

## Achados rebaixados a task (não bloqueiam nada)

Cadeia de protótipo em lookup de objeto literal; "every string a transform can produce…";
"render every unmapped code verbatim (injective)"; "Only two shapes actually close it"; "not a
subset of importados"; "never data loss"; o comentário `differ only when…` do
`query_repository.go`; os 5 universais falsos pré-existentes. Nenhum tem caminho de input
não-confiável — `provider_code`, `match_status`, `direction`, `side` e `confidence_band` são
todos constante de Go.

## Escopo da missão, redefinido pelo operador

**MIS-006 é integração ERP + planilha + vínculo. Só.**

`M-07` (quebra do chicken-egg, descoberta EAN→`catalog_product_id`) **sai do escopo**. Descoberta
de catálogo é caminho de MERCADO, não de vínculo — vai para a missão mercado junto com a coleta
(MC-11), que já estava fora. O gate live T13-T16 **não roda nesta missão**; a pergunta "a API do
ML resolve EAN→catálogo?" fica aberta e vira pré-requisito da missão mercado. `cmd/mlprobe` segue
untracked.

## Ordem até o fim

ANCHORS-3 (3 condições) → merge → VINC-NEUTRO → merge → **fechar M-06** → um chip com os 4 bugs
observáveis (deadline em `POST /erp/imports`; candidato STALE fora do cap; duplicado entre
âncoras em `buildCollisionCandidates`; painel de cadeia travando em `Carregando…` no 5xx) →
**encerrar MIS-006**.

**Executado.** VINC-NEUTRO `3847fb4f` · ANCHORS-3 `28ac8ac5` · M-06 fechado `9c828044` ·
CHIP-FIM `312adc2d` · missão fechada `975ac82d`.

## Últimos dois cortes (operador, 2026-07-28)

`A2-R1` (forma do ramo AGAINST) e `G4` (índice para `(tenant_id, state,
internal_product_id)`) **saem sem sucessor**. Mesmo critério que matou o gate de prosa: nenhum
dos dois é defeito que o operador alcança. O `G4` não é sequer mensurável fora de produção — os
dois assentos que o levantaram disseram isso na cara, que faltava `EXPLAIN (ANALYZE, BUFFERS)`
em escala real. Registrar um item que ninguém pode medir é registrar trabalho para o gate, não
para o software.

Volta se e quando existir sintoma: fila lenta medida, ou um operador lendo um AGENTE errado na
tela.
