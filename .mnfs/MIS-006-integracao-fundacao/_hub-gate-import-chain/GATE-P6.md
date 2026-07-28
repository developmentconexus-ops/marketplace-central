# P6 dual gate — CHIP-IMPORT-CHAIN — rodado pelo HUB

`2026-07-28` · diff sob revisão: `5441fe18` → `42d7c2d1`, 15 arquivos, +454/−19, tudo `apps/web`
Código congelado verificado por número: `git diff --name-only aceff011 c2229f26 -- . ':(exclude).mnfs'` = **0 caminhos**;
`git diff --stat 42d7c2d1 c2229f26 -- apps/web` = **vazio**. Tudo depois de `aceff011` é evidência.

## VEREDITO DO HUB: ACEITO COM CONDIÇÕES

Os dois lados voltaram **APROVA**, independentes e cegos um pro outro. Nenhum achou bloqueante.
As condições abaixo não são desacordo com eles — são o que os dois marcaram como ressalva e o
hub decidiu cobrar antes do merge.

| Lado | Assento | Veredito | Artefato |
|---|---|---|---|
| Opus | cold gate-reviewer, read-only físico (sem Bash, sem git) | **APROVA**, 2 ressalvas + 3 notas | `p6-opus-r1.md` |
| GPT-5.6 Sol medium | `codex exec --sandbox read-only`, OS-process | **APROVA**, 2 ressalvas + 1 nota | `p6-sol-r1-raw.md` (stdout verbatim, só ANSI removido) |

Prompt único, congelado antes dos dois disparos, apontando os dois para o meio que ninguém tinha
revisto com olho fresco. Nenhum dos dois viu o veredito do outro.

## A convergência que vale — os dois acharam o MESMO buraco, sozinhos

`ImportChainPanel.tsx`, o `queue_read_at`. **Os dois lados, cegos, escreveram a mesma ressalva** e
os dois disseram explicitamente que a review fria anterior tinha afirmado que esse campo estava
`guarded correctly` — e que a afirmação é falsa.

O `formatDateTime` (`packages/web-query/src/index.ts`) guarda `!value` e `Number.isNaN(getTime())`.
Não tem `typeof === "string"`. O Sol executou a coerção em vez de argumentar sobre ela:

```
["2026-07-18T12:00:00Z"] → 2026-07-18T12:00:00.000Z
true                     → 1970-01-01T00:00:00.001Z
1                        → 1970-01-01T00:00:00.001Z
```

Um `queue_read_at` que chegue como número ou array não vira `—`: vira uma data que parece boa. E
esse é exatamente o mecanismo que o chip JÁ endureceu em `aceff011` para o `protocol`
(`typeof value === "string" && value.trim().length > 0`), e que o `renderCounter` já aplica aos
contadores (`typeof value === "number"`). Dois dos três mecanismos checam tipo. O terceiro não.

Um sítio consertado e o irmão deixado é a forma que a regra da terceira rodada existe para pegar.
Aqui foi pego na rodada 1, por dois leitores independentes, e por isso é condição e não REPORT.

## As condições

1. **`queue_read_at` recebe o guard de tipo**, do mesmo formato do `protocol`, com teste. O caso
   discriminante é payload com o campo em tipo errado → `—`, não data de 1970.
2. **Trava unitária do zero conhecido.** Nenhum dos 8 testes do painel assere `0`. O
   `ImportacaoDetailPage.test.tsx` passa `enfileirados: 0` e não assere nada sobre ele. O
   comportamento está certo por leitura e foi observado ao vivo no L2 — mas uma "simplificação"
   futura de `renderCounter` para `value ? value : <UnknownValue/>` quebra a metade do ADR-17 que
   o operador mais lê, e os 521 continuam verdes. Um teste: `enfileirados: 0` →
   `toHaveTextContent("0")` e `not.toHaveTextContent("—")`.
3. **404 deixa de ser atribuído por status.** `candidate.status === 404 || candidate.error ===
   "import_not_found"` curto-circuita no status, então QUALQUER 404 — rota não montada, prefixo de
   `baseUrl` errado, proxy — vira *"Importação não encontrada."*, afirmando um fato sobre o DADO
   quando a verdade é que a ROTA não está lá. É a classe do catalog-503 do M-02 falando com
   confiança errada. O discriminador certo já está escrito e é o segundo termo do `||`: o servidor
   emite corpo achatado (`writeError(w, 404, "import_not_found", "")` →
   `map[string]string{"error": code}` → o SDK lança `{status, error}`), verificado pelo lado Opus
   ponta a ponta. Basta cair o `status === 404 ||`.

## O que o hub REJEITA das ressalvas

**`Number.isFinite` aceitar `-1` ou `1.5` (Sol).** Não vira condição. Apertar para inteiro não
negativo inventa uma regra de domínio que o contrato não tem, e faz um servidor que mandou um
número real — surpreendente, mas real — virar `—`. Isso é a OUTRA mentira do ADR-17: dizer "não
sei" quando o servidor disse. Os contadores saem de `count(*)`; se um dia saírem negativos o
defeito é do backend e o FE deve mostrar o que chegou, não maquiar.

## I9 — descarregado pelo HUB, não pelo chip, não pelos reviewers

Os DOIS lados deixaram o I9 NOT-PROVEN pelo mesmo motivo, e o motivo não é falha deles: um assento
de review não roda nada. Então o hub rodou, no worktree do chip fixado em `c2229f26`:

```
L0  npx --no-install tsc -p apps/web --noEmit
    → 15 erros. Zero em pages/importacoes/.
      5 pages/produto · 5 pages · 3 pages/vinculos · 2 pages/mutations
      (os 3 de pages/vinculos/ são INCOMPARABLE ausente — pré-existentes, do VINC-NEUTRO)

L1  npx --no-install vitest run
    → Test Files 65 passed (65) · Tests 521 passed (521) · 39.84s
```

Bate com a declaração do chip nos dois números. **A declaração era verdadeira; agora ela é
verificada.** É a primeira vez nesta missão que a metade executável do gate teve verificador
independente do implementador — ver o FINDING do CHIP-ANCHORS-3, que nomeou a lacuna, e a decisão
do hub sobre o assento executor.

## Erros dos próprios reviewers, registrados para ninguém caçar fantasma

**Sol leu a árvore errada, e a culpa é do prompt do hub.** Meu prompt deu `WORKSPACE:` o checkout
principal. O `EVIDENCE.md` do chip vive no worktree do chip e não está commitado no main. Então o
`I7 — FAIL (NOT-EVIDENCED)` e o `I8 — NOT-EVIDENCED` do Sol repousam sobre *"there is no
EVIDENCE.md in `_chip-import-chain/`"*, que é artefato de árvore errada, não achado. O lado Opus,
que eu apontei explicitamente para o worktree do chip, leu o EVIDENCE e deu PASS nos dois. **Erro
meu de instrumento, não do chip e não do modelo.**

**Sol citou um arquivo que não existe.** Registrou como hazard vivo que
`apps/web/vitest.chip.config.ts` *"exists on disk"* e citou o conteúdo dele em `:4` e `:16`.
Verificado: não existe no checkout principal, não existe no worktree do chip, não é tracked em
nenhuma árvore, e `git ls-tree` em `c2229f26` não tem o caminho. A menção provável de origem é
`chip.md:113-114`, que manda deletar o arquivo antes do commit — o modelo leu a ORDEM e reportou
como OBSERVAÇÃO. Uma citação file:line fabricada num gate é séria: não invalida os outros achados
dele (o da coerção de data ele PROVOU executando), mas o peso de qualquer alegação de existência
vinda desse lado cai. Anotado no profile como assinatura, não como caso isolado.

## O que continua não provado, e fica declarado

- **`vinculados` não-zero.** Ninguém exerceu. O L2 provou que o `0` da tela é VERDADEIRO (conferido
  contra SQL direto, com a subcontagem de zero à esquerda descartada por duas medições), mas o
  caminho com contagem real não foi percorrido. Herdado como limitação declarada.
- **Montagem com literalmente zero instalações** (nota do Sol, justa). O L2 dirigiu num tenant COM
  instalação ML e usou como observável do gate os `href` limpos das rotas não-gated. Isso mostra a
  decisão de roteamento, não o estado sem marketplace. A leitura estática dos dois lados diz que
  funciona — `InstallationProvider` sempre renderiza filhos, só o `InstallationGate` bloqueia, e
  `/importacoes` está fora do bloco. Fica como I1 provado por leitura + observável indireto.
- **Payload malformado contra servidor conforme** é inalcançável por construção: OpenAPI declara os
  cinco campos `required` e o `protocol` tem `CHECK (protocol ~ '^#[0-9]{3,}-E$')` `NOT NULL` no
  banco. O teste unitário é a única prova possível dessa metade, e foi aceito como tal.
