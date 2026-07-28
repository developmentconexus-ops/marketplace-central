# HUB GATE — CHIP-ANCHORS-2 (`chip/anchors-2`)

**Pack do hub.** O chip não carimbou o marcador abaixo e fez certo: a linha `P6-DUAL-GATE:` é
autoridade do hub. O fecho do chip é `AGREEMENT — P6 discharged`, com o ledger de discharge ao lado
(`_chip-anchors-2/EVIDENCE.md`). Este arquivo é o registro independente do hub e é o que o hook de
merge lê, porque o hook varre a working tree de `main` antes de o branch do chip chegar.

- Branch: `chip/anchors-2`
- Code tip: `ae6c6525a3db3ad4c4b278e587040442ec153d22`
- Pack tip: `2d42dd358d806c18afa5055ab1127143853cf5ff`
- Rebasado em `e98d8193` — `merge-base` confirmado pelo hub, o branch é `main` mais 9 commits
- Escopo: F-01 `refforn` · F-02 `INCOMPARABLE`+`side` · F-03 ordem estável + `≡` · F-04 chain-read
- Rulings: `_chip-anchors-2/hub-rulings.md` (A2-R1, A2-R2)

---

P6-DUAL-GATE: REFUTED — round 1 do gate real, ver `p6-reconciliation-r1.md`

> **Este marcador foi corrigido em 2026-07-28, depois do merge.** Ele dizia `AGREEMENT` com uma
> declaração honesta de que faltava o verdict do lado GPT. O operador determinou que declarar a
> lacuna não é o mesmo que fechá-la, e exigiu o dual gate que a harness prescreve — Opus e Sol
> medium, adversarial, específico sobre o que foi implementado e o que impacta. O gate rodou e
> **os dois lados REPROVARAM, independentemente, por caminhos diferentes**, sobre o mesmo input
> congelado (`p6-input-r1.patch`, sha256 `0762d05f…dd32`).
>
> O merge `dbdcdfb1` fica registrado como **mergeado e NÃO aprovado por gate**. Os corretivos
> CORR-1..CORR-6 estão nomeados na reconciliação e caem em `main`. O bloqueante que ninguém
> tinha visto — nem os dois gates anteriores do chip, nem este hub na leitura pré-merge — é o
> B-01: o classificador decide se o lado ERP tem a âncora `seller_sku` lendo `refforn`, e emite
> um motivo que afirma que o produto ERP não tem CODPROD.
>
> A seção "O que este AGREEMENT cobre" abaixo fica como está, porque descreve com precisão o que
> aconteceu antes desta correção. Ela está SUPERSEDIDA como veredito e preservada como registro.

LIVE-WAIVED-BY-OPERATOR: waiver dado pelo operador em 2026-07-28, com a alternativa e o risco
declarados antes da escolha. F-04 é o único item com superfície live e não tem tela até a onda 2; o
compose declara `env_file: .env` relativo ao diretório do compose, então dirigir a partir do worktree
do chip exigiria copiar credenciais Sankhya/ML para um segundo diretório. O operador optou por
dirigir o endpoint no P7 do CHIP-IMPORT-CHAIN, que é a tela `/importacoes` que o consome — prova ao
vivo onde um usuário de verdade a vê. Risco aceito, e nomeado antes do aceite: um defeito de FIAÇÃO
(composition root, decorator perdido) fica escondido até a onda 2. É a classe exata do catalog-503 do
M-02, que teste nenhum pegou.

---

## O que este AGREEMENT cobre, e o que ele NÃO cobre

Declarado assim porque um marcador que omite sua cobertura é a mesma alegação inflada que R-24
proíbe.

**Não houve verdict do lado GPT sobre o tip final.** O ledger de dispatch do chip mostra `gpt-5.6-sol`
e `gpt-5.6-luna` como **IMPLEMENTADORES** (D2, D4, D6, D7, D8, D10, D11, D12), não como gate. Os três
passes adversariais — D3, D9, D13 — foram subagentes Claude sonnet, independentes do implementador. O
implementador não é gate de si mesmo, então o lado GPT não gateou.

**Os artefatos dos três reviewers estão 0 bytes em disco.** O chip declarou isso por conta própria:
os transcripts foram streamados, não persistidos, e os verdicts estão transcritos das mensagens de
conclusão. O hub não consegue reler os verdicts originais. Registrado, não desculpado — é uma lacuna
de proveniência, e a lacuna é do harness, não da honestidade do chip.

**O segundo pass model-side foi feito pelo HUB, lendo o código.** O que o hub verificou por si, no
tip, citado por string:

| Verificação | Resultado |
|---|---|
| `merge-base HEAD 2d42dd35` | `e98d8193` — branch é `main` + 9 commits, sem enxerto |
| write-set fora de `.mnfs/` | 20 arquivos, **0** sob `apps/web/`, **0** migrations (C10) |
| `IdentityAnchorRefforn` em `apps/server_core` | nenhuma ocorrência (C1) |
| `DirectionUnavailable` em `generation_service.go` | 3 sítios: `:704` é o caminho `!anchor.Supplied`; `:658` é comparação, não semeadura; `:642` é o ramo EXCLUÍDO por A2-R1 (C11b) |
| ramo excluído de A2-R1 | `missingMatchedAnchorReason`, `default:` só alcançável com valor não-vazio dos DOIS lados — exatamente o que a ruling mandou não tocar, e nada além disso |
| SQL do chain-read | as três CTEs escopadas por `tenant_id`, `DISTINCT` em `resolved_products` e `queued_products`, `entity = 'market'`, `COALESCE(cursor -> 'pending', '[]')`, **sem `ORDER BY`** — a armadilha do `::text` sem alias não é alcançável |

O hub é Opus e leu no tip; o reviewer adversarial é sonnet e leu nos tips de feature; o implementador
é GPT. Três contextos, dois modelos do lado de revisão. É o que existe, e é o que o marcador afirma —
não "dual gate completo conforme a doutrina", que seria falso.

## Achados que saem deste merge (nenhum bloqueante)

Os três primeiros o chip levantou; os dois últimos são observação do hub.

- **F-8, do chip** — `{id}` malformado na rota nova provavelmente responde 500 e não 404
  (`PathValue` vai direto para um parâmetro tipado `uuid`). **Não verificado por ninguém**, e dito
  assim. É o padrão pré-existente de `handleGetImport`, logo não é regressão deste chip, e nenhum
  critério o cobre.
- **F-1, do chip** — lacuna do C2 por desenho: âncora presente dos dois lados não emite motivo.
  Fechado como A2-R2; a onda 2 sabe.
- **F-1b, do chip** — o ramo `UNAVAILABLE` excluído por A2-R1. A verdade ali é `AGAINST`, e virar
  `AGAINST` mexe em D-121. **Decisão do operador**, levada por este merge.
- **Observação do hub — o join do `vinculados` é igualdade de string depois de castar um inteiro.**
  `links.internal_product_id::text = products.codprod`: um `codprod` com zero à esquerda (`007`)
  nunca casaria o inteiro `7`. Os dados de hoje não têm essa forma (`33698`, `90008`), então não é
  defeito agora — é uma suposição sobre a forma do dado que ninguém escreveu em lugar nenhum, e que
  quebraria em silêncio, contando a menos sem erro nenhum.
- **Observação do hub — `apps/web` `tsc` fica VERMELHO em `main` a partir deste merge.** 3 dos 15
  erros são deste chip (`QueueRow.tsx` ×2, `VinculoDrawer.tsx`), causados de propósito por
  `INCOMPARABLE` entrar em `ProductLinkReasonDirection`; os outros 12 são baseline. É o mecanismo
  desenhado em D-B — o `Record<Direction, …>` torna esquecer o caso novo impossível em vez de
  improvável — e o dono é o CHIP-VINC-NEUTRO da onda 2. Enquanto ele não fechar, `main` carrega FE
  `tsc` vermelho. Consequência aceita em D-F, não descoberta agora.

## Erros de dono que o chip declarou em vez de esconder

Registrados porque a prática é o que se quer repetir, não só o resultado.

- **F-3** — um fixture que não executa não é evidência (protocolo estendido do pgx, um comando por
  `Exec`). Duas rodadas cegas falharam antes de a correção sair na lane do dono.
- **F-6** — F-02 foi despachado como feature inteira, contra o core §4. Do D5 em diante, um dispatch
  por slice.
- **F-7** — `git checkout --` usado para reverter uma edição de must-fail descartou trabalho de slice
  não commitado. Pego por grep, reautorado, reverificado.
- **F-9** — um `DISTINCT` que fixture nenhum consegue quebrar está inspecionado, não provado. O
  reviewer do F-04 pegou; virou guard de verdade (codprod `101` resolvido duas vezes, must-fail
  `Vinculados:3 FAIL`).
- No mesmo round, o OpenAPI prometia uma drenagem do scheduler que código nenhum executa — só o job
  `products` está registrado, e a fila de mercado só recebe append. **Deletado sob R-25**, não
  anotado, com o comentário Go corrigido no mesmo commit.
