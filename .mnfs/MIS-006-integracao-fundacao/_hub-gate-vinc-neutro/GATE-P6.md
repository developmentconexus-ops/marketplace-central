# GATE P6 — CHIP-VINC-NEUTRO, round 7 (forma nova: os assentos leem o DIFF, não o pack)

`alvo` = `main` @ `94374920` · `chip` = `chip/vinc-neutro` @ `293c1485` (merge de `main` em `43794cdf`)
`entrada dos assentos` = `code-diff.patch` (2.662 linhas, `git diff main chip/vinc-neutro -- apps
contracts packages`, sha256 `511a9b04…e3e6bb3a`) + os critérios abaixo + a lane crua do chip.

Ratificado em `f19793d7` / `678a6d51`: o pack do chip **não é entrada de revisor**. Seis rounds
anteriores neste chip gastaram dois assentos cada lendo pack. Esta rodada corrige a MIRA, não o rigor.

## Régua desta rodada (`3f8560b1`)

- **BLOQUEIA** só quando o achado nomeia um **observável errado**: comportamento, segurança, dado,
  ou contrato publicado.
- **REPORT** para prosa, contagem, metadado, citação, formatação, higiene. Registra, conserta, não
  segura merge e não abre rodada.
- Discriminador, respondido ANTES de atribuir a classe: *deixa o achado exatamente como está, sobe
  o código — o que um usuário, um operador, um chamador ou uma linha gravada faz de diferente?*
  Sem resposta = REPORT.
- Reclassificar é do hub. O assento reporta tudo que achar, na severidade que acreditar.

## Custódia, medida pelo hub antes de abrir a rodada

```
git branch --contains main                                     → chip NÃO listado (main avançou depois do merge do chip)
git diff --numstat main chip/vinc-neutro -- apps contracts packages | awk '$1==0'  → vazio
delta = 11 arquivos, TODOS em apps/web/src/pages/vinculos/     → nada fora do write-set concedido
```

Nenhum revert. **Isto não é critério desta rodada** — está aqui para o assento não re-derivar.

## Critérios — verbatim do `chip.md`

### V-1 — `QueueRow.tsx` enumerava direções por string literal (o critério que decide o chip)

O código original era

```tsx
const shown = [...byDirection("AGAINST"), ...byDirection("FOR"), ...byDirection("UNAVAILABLE")].slice(
```

Type-correct, e o compilador não avisa. Numa linha cujos motivos são TODOS `INCOMPARABLE` — a linha
que a mudança de backend criou — `shown` fica vazio, `hidden = reasons.length - 0 > 0`, e a célula
renderiza **um botão `+2` sozinho, com zero chips**. Viola a invariante que o próprio arquivo
documenta: *"Ranking (never filtering) is what keeps at least one motivo on screen"*. Com a quarta
direção existindo, a expressão VIROU filtragem.

**Explicitamente no contrato:** *"Consertar os 3 erros que o `tsc` reporta deixa este defeito
calado. Seu gate falha se você fechar o `tsc` e não fechar isto."* `tsc` verde NÃO é o critério.

Irmãos na mesma armadilha, nomeados no contrato: `directionClasses` e `directionGlyphs` são
indexados; sem a chave nova rendem o literal `undefined` no glyph e `className` `undefined`.

### V-2 — `refforn` é entrada morta no vocabulário

`anchorShortLabels` mapeava `refforn: "Ref. forn."`; o F-01 do CHIP-ANCHORS-2 tirou `refforn` do
vocabulário cross-side no backend. Exigido: **decidir e DECLARAR** — remover (o fallthrough verbatim
já cobre âncora desconhecida) ou manter com motivo escrito. Remover em silêncio é FAIL de evidência.

### V-3 — vocabulário neutro de provider, sem apagar o dado

Neutralizar **rótulo de coluna/seção**; **não** neutralizar o VALOR — o anúncio é de um provider e
esconder isso é mentira, não neutralidade. Armadilha que já queimou o CHIP-PED-FILA em 4 superfícies:
o wire carrega o slug `"mercado_livre"`; renderizar o slug cru na tela é bug. Se renderizar provider,
renderizar nome de exibição.

### V-4 — "Identificado por" NÃO é derivável de `direction == FOR`

`title FOR` já existe na base sob `TitleMatch`. Uma coluna alimentada por `direction == FOR` mostraria
"Título" como âncora identificadora, e o D-121 é explícito que title-only vai para REVIEW, nunca
identifica. Contrato congelado: "Identificado por" mostra as âncoras que **DECIDIRAM**, juntas por
` + `; âncora `UNAVAILABLE` ou `INCOMPARABLE` aparece em **Motivo** e não em Identificado por.
Se o wire não carrega qual âncora decidiu, **não inventar a coluna** — reportar lacuna e parar.

### V-5 — badge de auto-aprovado, com o brief do M-06 corrigido

O brief pedia `rule_matched=exact_ean_unique` E `actor=system`. As duas metades furadas:
`rule_matched` não está no wire (grep zero em `contracts/` e `packages/sdk-runtime/`), e
`CHECK (actor <> 'system' OR rule_matched = 'concordant_codprod_ean')` proíbe o par pedido.
Predicado correto, sem custo de backend: **a entrada de auditoria que resolveu o vínculo tem
`actor.actor_type === "system"`** — já no wire. Vínculo antigo sem entrada `system` → **sem badge**,
nunca badge fabricado.

### V-6 — propriedade

Migração: nenhuma. Backend Go: nenhum. Contrato/SDK: nenhum — `contracts/` e `packages/sdk-runtime/`
FORA do write-set. `VinculosPage.tsx` e `ImportacaoSection.tsx` são do CHIP-IMPORT-CHAIN: tocar =
colisão. `VinculosDesign.golden.test.tsx` é gate de design — afrouxar asserção sem dizer qual é FAIL.

## Lane crua do chip (round 6), para o assento não re-medir

```
tsc     → 12 erros, 0 em pages/vinculos/  (baseline inalterado)
vitest  → 64 arquivos / 534 testes, exit 0  (eram 531; +3 são os must-fail)
must-fail 1 — ACCEPT exists only as (exact_sku, seller_sku)                     RED antes / GREEN depois
must-fail 2 — a provider with no capability declaration is not producible       RED antes / GREEN depois
must-fail 3 — an INCOMPARABLE `side` outside the SDK union renders VERBATIM     RED antes / GREEN depois
sentinela (rename por sufixo de knownIdentityAnchors, só no port) → 2 failed | 3 passed (5), sentinela ✗
```

## Fora de escopo — não é achado bloqueante se cair aqui

- Qualquer coisa fora de `apps/web/src/pages/vinculos/`. Se o assento achar arquivo Go, `contracts/`
  ou `packages/sdk-runtime/` no diff, **isso** é o achado.
- `BatchPreviewModal.tsx` — união por literal em 3 sítios; já é REPORT registrado, dono é o hub
  (CHIP-B09 na fila). Não é defeito hoje.
- Os 12 erros de `tsc` do baseline fora de `pages/vinculos/`: do hub.

## O que o assento leitor NÃO deve fazer

- Não pedir o pack. Se um achado de CÓDIGO depender da derivação escrita do chip, peça uma seção
  **nomeada e limitada**; nunca a árvore.
- Não emitir achado sobre prosa, contagem, citação ou consistência de documento. Não há documento
  nesta entrada — só o diff.
- Não citar `file:line` de arquivo que não esteja no diff. Um assento já FABRICOU uma citação para
  arquivo que não existe em árvore nenhuma. Se o arquivo não está no patch, ele não está na rodada.
- Não achar que `tsc` verde ou `vitest` verde descarrega V-1. O contrato diz o contrário por escrito.
