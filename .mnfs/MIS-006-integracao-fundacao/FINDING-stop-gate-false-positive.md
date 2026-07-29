# FINDING — `stop-gate.sh` bloqueia turnos que não reivindicaram CLOSED

`status: upstream-candidate · CAUSA ATRIBUÍDA nos dois casos (2026-07-28, medido)`
`provenance: 2026-07-28 · D-122 · relatado por CHIP-ANCHORS-3, e disparado no HUB no mesmo dia`
`locus: mnfs-harness/harness/0.4.0/hooks/stop-gate.sh:32-43`
`destino: HARNESS-CORE upstream (Documents\mnfs-harness). NÃO editar o cache do plugin — sync é gate do operador.`

Duas sessões independentes bloqueadas no mesmo dia sem nenhuma reivindicação de close. Não é
ergonomia: um gate determinístico que dispara sem o evento que ele guarda ensina a ignorá-lo, e um
gate ignorado não guarda nada.

## Defeito 1 — o gatilho é substring crua, não evento (CERTO, os dois casos)

```bash
TAIL="$(tail -c 30000 "$TRANSCRIPT" 2>/dev/null)"
case "$TAIL" in
  *CLOSED*) ;;
  *) exit 0 ;;
esac
```

`*CLOSED*` casa com **qualquer** ocorrência do token em 30 KB de transcript cru: prosa de doutrina,
mensagem do hub citada, o pack lido no turno — e, decisivamente, **o próprio prompt de dispatch, que
ENUMERA o vocabulário de eventos** (`CLOSED`/`BLOCKED`/`ESCALATION`/`REQUEST`). Todo chip nasce
armado, e permanece armado enquanto o prompt estiver dentro da janela.

CHIP-ANCHORS-3 verificou o inverso por string no próprio pack: `status: NAO FECHADO PELO CHIP` na
linha 37, `P6-DUAL-GATE:` em três lugares, os três dizendo que **a linha é do hub**, e `AGREEMENT`
só em frases afirmando que nunca houve. Nenhum veredito escrito pelo chip, em rodada nenhuma.

**Forma do conserto:** casar o EVENTO emitido (a linha de evento do chip), não o token em prosa.

## Defeito 2 — o guard de escopo transforma o HUB em chip (mede o caso do hub inteiro)

O guard classifica sessão de chip por `cwd` conter `/.claude/worktrees/`. Nesta topologia o **hub**
tem cwd de worktree: o shell reseta para `.claude/worktrees/epic-lehmann-4ffbad`, um worktree stale
que o hub não usa para trabalhar. Medido hoje:

```
.claude/worktrees/epic-lehmann-4ffbad     .mnfs existe · packs de chip encontrados = 0   → BLOQUEIA
.claude/worktrees/happy-montalcini-b010c0 .mnfs existe · packs de chip encontrados = 35  → passa
```

O hub é bloqueado por um gate de chip, contra uma árvore onde ele nunca escreveu, por uma reivindicação
que ele não faz — o hub não emite `CLOSED`, ele o **recebe**.

**Forma do conserto:** classificar por papel declarado, não por forma de caminho.

## Defeito 3 — `$CWD/.mnfs` é resolvido do cwd do turno, não da raiz do repo (EXPLICA o caso do chip)

Medido pelo CHIP-ANCHORS-3: com `cwd` mais FUNDO que a raiz do repo (`apps/web`,
`packages/sdk-runtime`, `apps/server_core`), a linha 39 procura `$CWD/.mnfs`, que não existe →
**packs = 0** → guard dispara, com o pack intacto na raiz.

**O `cwd` do disparo, extraído pelo HUB (o chip não leu o próprio transcript):**

```
21:45:36.830Z  type=hook_blocking_error  cwd=…\happy-montalcini-b010c0\.mnfs\MIS-006-…\_chip-anchors-3
21:45:36.831Z  type=system               cwd=…\happy-montalcini-b010c0
21:46:09.922Z  type=message              cwd=…\happy-montalcini-b010c0
21:46:27.976Z  type=message              cwd=…\happy-montalcini-b010c0
```

Há **um** disparo, não cinco: só o registro `hook_blocking_error` é o hook. Os outros três são o
chip citando o texto depois, já de volta na raiz — é daí que vinha a impressão de que o `cwd` era a
raiz. **O `cwd` do disparo era o diretório do próprio pack**, então a linha 39 procurou
`…/_chip-anchors-3/.mnfs` e achou zero.

E o pack estava lá: `EVIDENCE.md` com mtime `21:43:34Z`, **dois minutos antes** do bloqueio.
"Pack ausente" está refutado por dois instrumentos independentes.

**Nenhum defeito sozinho explica o caso.** O defeito 1 armou o guard (o chip nunca emitiu `CLOSED`);
o defeito 3 fez a checagem de existência falhar. Precisa dos dois.

**Limite do instrumento, dito:** o `cwd` lido é o campo do REGISTRO que o hook produziu, não o
stdin do hook. É o cwd de sessão carimbado no próprio registro do bloqueio — o mais apertado que
este instrumento chega, e não é identidade.

**Defeitos 2 e 3 são o mesmo erro de raiz:** o hook trata `cwd` como se fosse a raiz do repo. Um
conserto cobre os dois — resolver a raiz (`git rev-parse --show-toplevel`) e classificar por papel
declarado, nunca por forma de caminho nem por cwd do turno.

## O que NÃO é a causa — três hipóteses testadas e refutadas

1. **"O `**` do glob exige ≥2 níveis"** (hipótese do CHIP-ANCHORS-3). Refutada: a linha 39 é
   `find "$CWD/.mnfs" -name EVIDENCE.md -path '*_chip*'`, agnóstica a profundidade. O
   `.mnfs/**/_chip-*/EVIDENCE.md` existe só na string de mensagem, não no mecanismo.
2. **"`find` falha com path Windows"** (hipótese do hub). Refutada por execução: `find` com
   `C:\Users\…\.mnfs` resolve e lista normalmente no Git Bash desta máquina.
3. **"Todo Bash meu carrega `cd <absoluto>`, então o cwd do disparo era a raiz"** — contra-indício
   que o próprio CHIP-ANCHORS-3 ofereceu CONTRA o mecanismo dele. Refutado pela medição acima: o
   cwd do registro do disparo é o diretório do pack. O chip estava certo no mecanismo e errado no
   contra-indício que ele mesmo levantou contra si.

## Reincidência 2026-07-29 (D-123) — dois disparos no HUB, defeitos 1+2 confirmados sem medição nova

O hub foi bloqueado duas vezes no dia do despacho do CHIP-VENDAVEL, e o mecanismo é exatamente o já
diagnosticado — nada de novo a medir:

1. **No turno de despacho.** O hub tinha acabado de emitir o prompt do chip, que ENUMERA
   `CLOSED`/`BLOCKED`/`ESCALATION`/`REQUEST` na cláusula (e) do contrato de comunicação. Defeito 1:
   o prompt de dispatch arma o guard contra quem o escreve, não só contra quem o recebe.
2. **No turno do RULING DR-2.** A janela de 30 KB ainda continha `CLOSED` da doutrina de aceitação
   que o hub cita ao chip ("`CLOSED` só com: branch + SHA, vereditos…").

Nos dois casos o hub não reivindicou fecho de coisa nenhuma — **o hub nunca emite `CLOSED`, ele o
RECEBE** — e o cwd continua o worktree stale `epic-lehmann-4ffbad` (defeito 2), onde o hub nunca
escreveu e nunca escreverá.

A partir daqui o hub PARA de contar disparo a disparo: enquanto o hub estiver a arbitrar um chip
vivo, a janela de 30 KB contém sempre o vocabulário de eventos — o guard passou de intermitente a
**permanentemente armado**, e registar cada ocorrência seria ruído, não evidência. O disparo do
turno do RULING DR-3 é o quinto e fica como o último contado.

**Custo acumulado, que é o argumento upstream:** o guard já disparou em 2 sessões distintas em
2026-07-28 e mais 3 vezes em 2026-07-29, **zero verdadeiros positivos**. Um gate cujo histórico
inteiro é falso-positivo não guarda o invariante que nomeia; ensina a passar por cima dele, e o dia
em que alguém reivindicar `CLOSED` sem pack ele será ignorado como os outros quatro. A reincidência
não muda o conserto proposto (casar EVENTO emitido; resolver raiz por `git rev-parse
--show-toplevel`; classificar por papel declarado) — muda a prioridade dele.

## Sobre a regra de não ler transcript

O chip pediu liberação da regra para este campo, ou o finding ABERTO. **Nenhuma das duas.** A regra
existe contra AUTO-certificação: o chip lendo o próprio transcript para provar a própria alegação.
Quem extraiu foi o hub, para quem a regra nunca foi sobre isso. A regra fica intacta e o dado
existe.
