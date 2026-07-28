# CHIP-VINC-NEUTRO — `/vinculos`: vocabulário neutro, `INCOMPARABLE` e badge de auto-aprovado

```yaml
chip: CHIP-VINC-NEUTRO
branch: chip/vinc-neutro
base_sha: 5441fe18f64171ef61cb03b51b5bf66e2922e4eb
wave: 2 (paralelo a CHIP-IMPORT-CHAIN e CHIP-ANCHORS-3 — write-sets disjuntos)
milestone: M-06 (F-04 + F-05, fatia por TELA conforme D-E)
design_ref: .mnfs/MIS-004-mvp-demo/DESIGN-REFERENCE.md (/vinculos)
```

## Autoridade

- `.mnfs/MIS-006-integracao-fundacao/DECISOES-D122-anchors-telas.md` — **leia inteiro**. D-B
  (`INCOMPARABLE` + `side`), D-C (ordem estável, `≡`) e D-E (M-06 fatiado por TELA) governam esta
  tela. O bloco de contrato congelado no fim é normativo.
- `.mnfs/MIS-006-integracao-fundacao/_hub-gate-anchors-2/p6-reconciliation-r1.md` — de onde vem a
  restrição B-02 abaixo.
- `.mnfs/MIS-006-integracao-fundacao/M-06-telas-sdk/milestone.md` — F-04. **Com duas correções
  obrigatórias, abaixo.** O resto do M-06 (F-01/F-02/F-03) é do CHIP-IMPORT-CHAIN.

Aponte para esses arquivos. Não os recopie (R-14).

## Ponto de partida honesto

`main` carrega **`tsc` VERMELHO em `apps/web`** desde `dbdcdfb1`, de propósito. 15 erros: 3 são deste
chip, 12 são baseline. Os 3 seus:

```
apps/web/src/pages/vinculos/QueueRow.tsx(34,7): TS2741: Property 'INCOMPARABLE' is missing … Record<ProductLinkReasonDirection, string>
apps/web/src/pages/vinculos/QueueRow.tsx(75,7): TS2741: Property 'INCOMPARABLE' is missing … Record<ProductLinkReasonDirection, string>
apps/web/src/pages/vinculos/VinculoDrawer.tsx(118,17): TS7053: Element implicitly has an 'any' type … 'ProductLinkReasonDirection' can't be used to index …
```

Os 12 baseline (`anunciosQueries.ts`, `anunciosQueryState.test.ts` ×2, `AnunciosTable.test.tsx`,
`ListingsRefreshControl.test.tsx`, `MutationPreviewModal.tsx`, `MutationResultSummary.tsx`,
`ProdutoPage.partialFailure.test.tsx` ×4, `ProdutoPage.test.tsx`) **não são seus**. Não os conserte:
são superfície de outro dono e consertá-los infla seu diff. Declare-os pré-existentes com a lista.

## F-05 — vocabulário neutro e `INCOMPARABLE` renderizado

### A restrição que decide se este chip presta: `QueueRow.tsx:159`

```tsx
const shown = [...byDirection("AGAINST"), ...byDirection("FOR"), ...byDirection("UNAVAILABLE")].slice(
```

Isto enumera direções por **string literal**. É type-correct. O compilador **não vai te avisar**.

Consequência para uma linha cujos motivos são todos `INCOMPARABLE` — que é a linha exata que D-B
criou: `shown` fica **vazio**, `hidden = reasons.length - 0 > 0`, e a célula renderiza **um botão
`+2` sozinho, com zero chips de motivo**. Isso viola a invariante ADR-17 que o próprio arquivo
documenta em `:154-156`:

```
// Rank AGAINST → FOR → UNAVAILABLE and slice. Ranking (never filtering) is
// what keeps at least one motivo on screen even for a row whose only signals
// are UNAVAILABLE ones (ADR-17 — motivo sempre visível).
```

O comentário diz "ranking, nunca filtragem". Com a quarta direção existindo, a expressão **virou uma
filtragem** que descarta tudo.

**Consertar os 3 erros que o `tsc` reporta deixa este defeito calado.** O mecanismo `Record<Direction,
…>` do D-B torna esquecer o caso novo impossível **onde o sistema de tipos alcança** — e aqui ele não
alcança. Seu gate falha se você fechar o `tsc` e não fechar isto.

Os mapas `directionClasses` (`:34`) e `directionGlyphs` (`:75`) também alimentam
`compactChipLabel`/`compactChip` por indexação; sem a chave nova eles rendem o literal `undefined` no
glyph e uma `className` `undefined`.

### `refforn` saiu do vocabulário

`anchorShortLabels` (`QueueRow.tsx:65-71`) ainda mapeia `refforn: "Ref. forn."`. O F-01 do
CHIP-ANCHORS-2 tirou `refforn` do vocabulário cross-side no backend — nenhum motivo com essa âncora é
emitido hoje. Entrada morta.

Decida e **declare** o que fez: remover a entrada (o fallthrough verbatim de `:81-83` já cobre
qualquer âncora desconhecida, então nada some da tela), ou mantê-la com o motivo escrito. Remover em
silêncio e não dizer nada é FAIL de evidência.

### Vocabulário neutro de provider

A tela fala "ML"/"Mercado Livre" em rótulos estruturais onde o dado é `provider_code`. Neutralize o
que é **rótulo de coluna/seção**; **não** neutralize o valor de dado — o anúncio É de um provider e
esconder isso é mentira, não neutralidade.

Cuidado com a armadilha do `provider_code` que já queimou o CHIP-PED-FILA em 4 superfícies: o valor
do wire é o slug `"mercado_livre"`, e mostrá-lo cru na tela é bug. Se você renderizar o provider,
renderize o nome de exibição.

### "Identificado por" — **NÃO é derivável de `direction == FOR`**

Restrição nomeada pela ruling A2-R2. `title FOR` **já existe na base** sob `TitleMatch`, antes deste
chip. Uma coluna "Identificado por" alimentada por `direction == FOR` mostraria "Título" como âncora
identificadora — e o D-121 é explícito que title-only vai para REVIEW, nunca identifica.

O contrato congelado do D-122 diz: "Identificado por" mostra as âncoras que **DECIDIRAM**, juntas por
` + ` (ex.: `CODPROD + EAN`); uma âncora `UNAVAILABLE` ou `INCOMPARABLE` aparece em **Motivo** e
**não** em Identificado por.

Se o wire não carrega qual âncora decidiu, **não invente a coluna**. Reporte como lacuna de contrato
e pare — uma coluna derivada errada é pior que uma coluna ausente. Não é sua a chamada de estender o
wire; é REQUEST ao hub.

## F-04 — badge de auto-aprovado. **O brief do M-06 está errado em dois pontos**

O `milestone.md:208` diz que o badge dispara em `rule_matched=exact_ean_unique` **E** `actor=system`.
As duas metades estão furadas, e você não deve implementar como está escrito:

1. **`rule_matched` não existe no wire.** `grep` em `contracts/marketplace-central.openapi.yaml` e
   `packages/sdk-runtime/src/` retorna **zero**. O campo existe só no DB (`migrations/0082…sql:38`) e
   num repo read per-link (`link_candidate_repo.go:411 ListDecisionsForLink`) que rota nenhuma expõe.
2. **`exact_ean_unique` com `actor=system` é proibido por constraint.** `0082…sql:54`:
   `CHECK (actor <> 'system' OR rule_matched = 'concordant_codprod_ean')`. Ator system só pode ter
   `concordant_codprod_ean`. O brief é da era D-120; **D-121 estreitou a política** (só CODPROD+EAN
   concordantes auto-aprovam) e o brief não acompanhou.

**O caminho que existe hoje, e que não custa backend nenhum:** a auto-aprovação grava a entrada de
auditoria com `ActorType: "system", ActorID: "auto_linker"`
([resolution_service.go:280](../../../apps/server_core/internal/modules/product_links/application/resolution_service.go:280)),
e `item.audit[].actor.actor_type` **já está no wire** — `useVinculosResolved.ts:18-23` já lê
`item.audit` e filtra por `next_state`.

Então o predicado do badge é: **a entrada de auditoria que resolveu o vínculo tem
`actor.actor_type === "system"`**. Mesma entrada que `resolutionAuditId` já localiza. Nenhum campo
novo, nenhum OpenAPI, nenhum SDK, nenhuma migration.

Registro obrigatório no EVIDENCE: que o brief foi corrigido, os dois motivos, e o predicado que você
usou. Não é liberdade de escopo — é o brief estando errado sobre o repo, e o pack corrigindo por
escrito em vez de você adivinhar em silêncio.

Vínculo antigo sem entrada `system` → **sem badge**, nunca um badge fabricado (ADR-17). Ausência de
badge significa "não foi automático", que é a verdade para todo vínculo manual.

## Matriz de propriedade

| Eixo | CHIP-VINC-NEUTRO |
|---|---|
| Migração | **nenhuma** |
| Backend Go | **nenhum**. Se você precisar de backend, é REQUEST ao hub, não edição |
| Contrato/SDK | **nenhum**. `contracts/` e `packages/sdk-runtime/` fora do write-set |
| FE — seu | `apps/web/src/pages/vinculos/`: `QueueRow.tsx`, `VinculoDrawer.tsx`, `QueueTab.tsx`, `ResolvidosTab.tsx`, `useVinculosResolved.ts` + os testes desses arquivos |
| FE — **NÃO seu** | **`VinculosPage.tsx`** e **`ImportacaoSection.tsx`** são do CHIP-IMPORT-CHAIN (ele move o `ImportacaoSection` para fora e ajusta as 2 linhas do `VinculosPage`). Tocar neles = colisão |
| `tsc` | fecha os **3** erros de `/vinculos`. Os 12 baseline continuam vermelhos e são declarados |

`VinculosDesign.golden.test.tsx` é seu, mas trate-o como gate de design, não como arquivo de
conveniência: se você o afrouxar para passar, diga qual asserção afrouxou e por quê.

## Ladder

- L0: `tsc` no write-set + lint.
- L1: `vitest` — baseline citada antes e depois.
- L2: **do hub**. Este chip não sobe servidor, não liga em `:8080`, não lê `.env*`.

O worktree do chip não tem `node_modules`. Junction para o do repo principal + config de vitest do
chip (`fs.allow` + `setupFiles` absoluto), e **delete a config antes de commitar**. `npx --no-install
tsc` no worktree passa vazio — é pass vacuoso; rode o `tsc` a partir da raiz principal apontando para
o worktree.
