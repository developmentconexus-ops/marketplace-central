# CHIP-VINC-NEUTRO — contrato de validação

PASS/FAIL/NOT-PROVEN por critério. Um critério que você não consegue fazer **totalmente** vira
REPORT e o chip fecha sem ele (R-24). Verifique por **string**, nunca por linha.

---

## V1 — `INCOMPARABLE` renderiza chip com glyph e classe próprios

Os três mapas indexados por direção (`directionClasses`, `directionGlyphs`, e o do
`VinculoDrawer.tsx:118`) têm a chave `INCOMPARABLE` com valor de token de design — nunca literal
Tailwind, nunca `undefined`.

Cite o glyph escolhido e o par de tokens. O contrato congelado do D-122 governa a semântica: um
`INCOMPARABLE` **não é** um `AGAINST` (não bloqueia) e **não é** um `UNAVAILABLE` (o provider FORNECE
a âncora; o outro lado é que não permite comparar), então reusar as cores de qualquer um dos dois é
FAIL — a tela passaria a mentir a distinção que D-B foi criado para mostrar.

## V2 — linha 100% `INCOMPARABLE` mostra ao menos um motivo

**O critério que decide este chip.** Teste com uma linha cujos `reasons` são TODOS `INCOMPARABLE`.

Esperado: pelo menos um chip de motivo visível no modo compacto. Proibido: célula com só o botão
`+N` e zero chips.

Asserção sobre o **DOM renderizado**, não sobre o array `shown`. Uma asserção sobre a variável
interna não prova a tela.

## V3 — must-fail do V2

Reverta `QueueRow.tsx:159` para a enumeração de três direções, mostre o V2 **falhando**, com a
assinatura citada (quantos chips apareceram, o que o botão dizia). Restaure, mostre verde.

Sem must-fail, o V2 está inspecionado, não provado — foi o F-9 do CHIP-ANCHORS-2.

## V4 — `side` do `INCOMPARABLE` chega ao operador

O contrato congelado do D-122: `side` só existe em `INCOMPARABLE` e diz **de que lado** falta o dado.
Um `INCOMPARABLE` renderizado sem essa informação em lugar nenhum (chip, tooltip, ou drawer) perde
exatamente o que D-B acrescentou.

Se `side` não vier no wire para algum caminho, é REPORT com o caminho nomeado — não invente o lado.

## V5 — `tsc` do write-set fecha, baseline declarada

`tsc -p apps/web --noEmit` a partir da raiz do repo principal. Esperado: os **3** erros de
`/vinculos` sumiram; os 12 baseline listados um a um, por caminho, declarados pré-existentes e fora
do write-set.

**`tsc` verde não é o critério deste chip.** V2 é. Um diff que zera os 3 erros e falha o V2 é FAIL —
essa é a lição B-02 inteira.

Se você tocar em qualquer um dos 12 baseline, o chip está fora de escopo.

## V6 — "Identificado por" ou está correta, ou não existe

Uma das duas, e diga qual:

- **Implementada:** prove que a fonte são as âncoras que DECIDIRAM, e mostre um caso `title FOR`
  (que já existe na base sob `TitleMatch`) **não** aparecendo na coluna. Sem esse caso negativo, o
  critério é NOT-PROVEN — é exatamente o modo de falha que a A2-R2 nomeou.
- **Não implementada:** REPORT dizendo qual campo falta no wire e por que a derivação disponível
  seria falsa. NOT-PROVEN honesto vale mais que uma coluna errada.

## V7 — badge de auto-aprovado dispara em `actor_type === "system"`

Três casos, os três com asserção sobre o DOM:

1. vínculo resolvido com entrada de auditoria `actor.actor_type === "system"` → **badge presente**;
2. vínculo resolvido por operador → **badge ausente**;
3. vínculo antigo sem entrada de auditoria resolvedora → **badge ausente**, sem crash, sem
   `undefined` na tela.

O caso 3 é o que o `milestone.md:213` chama de registro pré-M-05. Ausência de badge = "não foi
automático", que é verdade. Nunca um badge fabricado (ADR-17).

## V8 — a correção do brief F-04 está registrada

O EVIDENCE declara, por escrito: que o predicado do `milestone.md:208` não foi implementado como
escrito; que `rule_matched` não está no wire (cite o grep vazio em `contracts/` e
`packages/sdk-runtime/`); que `0082_product_link_decisions.sql:54` proíbe `actor=system` com
`exact_ean_unique`; e qual predicado você usou.

Isto não é burocracia. Sem esse registro, o próximo leitor abre o `milestone.md` e acredita nele.

## V9 — `refforn` no vocabulário: decisão declarada

O que aconteceu com `anchorShortLabels.refforn`, e por quê. Se removeu, mostre que
`anchorShortLabel()` continua devolvendo a âncora verbatim para entrada desconhecida — nada some da
tela, que é a promessa do comentário em `QueueRow.tsx:62-64`.

## V10 — vocabulário neutro não apagou dado de provider

Prove as duas metades: um rótulo estrutural que ficou neutro, **e** um valor de dado que continua
dizendo de qual provider o anúncio é. Neutralizar o valor é FAIL.

E: nenhum slug cru na tela. `grep` por `"mercado_livre"` no seu write-set renderizado retorna zero.
Essa é a armadilha que pegou o CHIP-PED-FILA em 4 superfícies, a quarta achada só pelo refuter do
hub.

## V11 — sem dano colateral

- `vitest` com a contagem citada antes e depois. Baseline vermelha continua vermelha e é declarada;
  verde que virou vermelha é FAIL.
- `VinculosDesign.golden.test.tsx` verde. Se você afrouxou alguma asserção, diga qual e por quê.
- `git diff --name-only 5441fe18f64171ef61cb03b51b5bf66e2922e4eb HEAD` — **zero** caminhos fora de
  `apps/web/src/pages/vinculos/` (fora o `.mnfs/` do pack), e **zero** ocorrências de
  `VinculosPage.tsx` ou `ImportacaoSection.tsx`.

`base_sha` é PISO, não ponto fixo. `main` se move enquanto você trabalha; rode o comando acima contra
o `main` de verdade antes do fecho e declare o que apareceu que não é seu.

---

## O que este contrato NÃO cobre

- **Nenhum L2.** O drive ao vivo de `/vinculos` é do hub, no P7.
- **Os 12 erros de `tsc` baseline** não são deste chip.
- **`/importacoes`, `/integracoes` e o `ImportacaoSection`** são do CHIP-IMPORT-CHAIN.
- **Qualquer mudança de wire.** Se um critério aqui exigir campo novo, ele vira REPORT + REQUEST ao
  hub, nunca uma edição de contrato por este chip.
