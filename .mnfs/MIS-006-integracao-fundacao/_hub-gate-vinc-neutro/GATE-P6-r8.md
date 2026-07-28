# GATE P6 — CHIP-VINC-NEUTRO, rodada 8

`alvo` = `main` @ `da3248b7` · `chip` = `chip/vinc-neutro` @ `bfc1d9bb`
`entrada dos assentos` = `code-diff-r8.patch` (3.016 linhas, `git diff main bfc1d9bb -- apps contracts packages`,
sha256 `20f375ab…c8571e6`) + estes critérios verbatim + a lane crua do chip + `DRIVE-EVIDENCE.md`.

Assentos CEGOS entre si. **Nenhum recebe o pack do chip** (§11, `f19793d7`). Sete rodadas anteriores
gastaram assento lendo pack; esta corrige a MIRA, não o rigor.

## Custódia, verificada pelo hub antes de abrir

- Nenhum arquivo de código apagado ou renomeado contra o tip ATUAL da main
  (`git diff --diff-filter=DR main bfc1d9bb -- apps contracts packages` vazio). Esta é a armadilha
  que custou quatro rodadas nesta missão: chip velho apaga feature da main no merge e o gate,
  medindo contra a base de despacho, não pode ver.
- O chip **nunca tocou** `docs/HARNESS-PROFILE.md` (`git log <merge-base>..bfc1d9bb --` vazio). As
  58 linhas que o diff bruto mostraria são do HUB, chegadas pelo merge de `main` em `0b1d0f36`.
  Excluídas da entrada do assento.
- Pack do chip existe e é versionado no tip dele.

## O que mudou desde a rodada 7 (o que esta rodada julga)

`8271991b..bfc1d9bb`, três commits. Duas coisas, e as duas nasceram de ordem da rodada 7:

1. **`QueueRow.tsx`** — a justificativa do fix foi REESCRITA. A antiga dizia que o caso manchete
   (todos os motivos INCOMPARABLE ⇒ `shown` vazio ⇒ `+N` sobre zero chips) era o dano. O chip
   corrigiu contra si: aquilo é **PRODUZÍVEL mas INALCANÇÁVEL** nesta árvore, e a metade alcançável
   é menor — com `COMPACT_CHIP_LIMIT === 2`, linhas que gastam os dois slots num FOR e no
   `marca UNAVAILABLE` (ausência permanente, acionável por ninguém) e não conseguem exibir o
   `ean INCOMPARABLE`, que o operador PODE acionar.
2. **`wireFixtures.guard.test.ts`** — dois braços MUST-FAIL novos (tupla produzível com conjunto de
   sinais que sítio nenhum emite; uma âncora carregando sinal e ausência ao mesmo tempo) e **um
   braço MUST-PASS** (`title` FOR + `title` AGAINST no mesmo candidato, que o gerador emite de
   verdade em `applySingleAnchorScore`).

## Critérios desta rodada

**V8-1 · BLOQUEANTE — a alegação de ALCANCE é sobre realidade em execução, e o assento tem a medida.**
O comentário de `QueueRow.tsx` afirma, em forma existencial: *"Rows of exactly that shape are on the
screen today."* Isso não se decide lendo o gerador — ler o gerador decide o que é PRODUZÍVEL e é
cego para o que é ALCANÇÁVEL, que é a distinção inteira desta rodada. Decide-se contra
`DRIVE-EVIDENCE.md`, que é observável MEDIDO PELO HUB em navegador real (§11 item 4). Nenhum assento
tem navegador. Se a forma alegada não estiver no dump de nove linhas, a frase é falsa e cai por R-25
(prosa falsa em contrato se DELETA, não se ameniza).

**V8-2 · BLOQUEANTE — o bloco de escopo de `wireCandidate` é uma alegação em DUAS direções.**
O docstring (`wireFixtures.ts`, bloco de `A candidate the backend can actually emit`) diz "throws
unless ALL of" seguido de três marcadores, e depois diz "What it does NOT check" seguido de três
buracos declarados. As duas listas são verificáveis por string e **as duas podem ser falsas**:

- um marcador da primeira lista que o código não impõe é garantia vendida e não entregue;
- um item da segunda lista que o código NA VERDADE checa é buraco declarado onde não há — falsidade
  igual, apontando para o outro lado, e é a que ninguém procura.

Verificar marcador a marcador contra o corpo de `assertProducibleReasons` / `assertProducibleCandidate`.

**V8-3 · BLOQUEANTE — o braço MUST-PASS precisa de poder discriminante PROVADO.**
Uma suíte só de must-fail não enxerga guard LARGO DEMAIS: um guard mais estrito deixa todo must-fail
mais verde. O braço `title FOR + title AGAINST` existe exatamente para essa direção. Então ele só
vale se a lane mostrar o par: **VERMELHO sob a regra larga** (a que rejeitava âncora repetida) e
**VERDE sob a regra atual**, com saída crua nomeando o braço. Braço must-pass sem o par é decoração.

**V8-4.** `wireFixtures.ts:56` diz que `DECLARED_PROVIDER_CAPABILITIES` *"is the whole set of
declarations that exist"* — alegação TOTAL sobre a ÁRVORE, não sobre o arquivo. Verificar por string
que hoje há um único adapter declarando capacidade. Se houver dois, a frase é falsa; se a frase for
verdadeira hoje mas o guard não puder detectar um adapter novo, isso é REPORT com o gatilho de
expiração nomeado, não bloqueante.

**V8-5.** Nenhuma asserção de terceiro quebrada fora do write-set do chip.

## Entrada dos assentos (§11, incluindo o item 4)

1. `code-diff-r8.patch` — diff contra o tip ATUAL da main;
2. estes critérios verbatim;
3. saídas cruas de lane do chip;
4. **`DRIVE-EVIDENCE.md` — observável MEDIDO PELO HUB.** Dump de wire de nove linhas + a tabela
   aritmética que mapeia as nove pelo `QueueRow.tsx` até os chips renderizados, conferida contra o
   navegador.

**O artefato declara o próprio limite e o assento deve respeitá-lo:** o drive testemunha **3 dos 16
sítios PRODUZÍVEIS** — nove linhas, um provider, uma instalação, um dia. Tratá-lo como cobertura do
wire é super-leitura. Ele é decisivo para ALCANCE (uma forma que está na tela está na tela) e mudo
sobre os treze sítios não testemunhados.

## O que o assento NÃO deve fazer

- **Não cobrar deste chip a invariante `Side` só em INCOMPARABLE no produtor Go.** Backend Go está
  fora do write-set (V-6, verbatim: "Backend Go: nenhum"). Já é dívida do HUB, registrada.
- **Não abrir achado de prosa do pack.** O assento recebe o diff, não o pack. Achado que só se
  sustenta lendo a justificativa do chip é REPORT, nunca BLOQUEANTE.
- **Não reabrir as rodadas 1–7** sem observável NOVO. Releitura não reabre.
- **Não tratar o drive como cobertura.**

## Régua (`3f8560b1`)

BLOQUEIA só quando o achado nomeia um **observável errado**: comportamento, segurança, dado, ou uma
alegação que o código contradiz — com `file:line`, ou string que existe/não existe, ou saída de
comando. Achado de gosto, de escopo especulativo, ou de "poderia ser melhor" é REPORT. Achado sem
observável nomeado não é achado.
