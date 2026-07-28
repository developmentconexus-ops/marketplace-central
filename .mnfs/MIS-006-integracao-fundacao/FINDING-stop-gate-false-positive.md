# FINDING — `stop-gate.sh` bloqueia turnos que não reivindicaram CLOSED

`status: upstream-candidate` · `provenance: 2026-07-28 · D-122 · relatado por CHIP-ANCHORS-3, e disparado no HUB no mesmo dia`
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

## O que NÃO é a causa — duas hipóteses testadas e refutadas

1. **"O `**` do glob exige ≥2 níveis"** (hipótese do CHIP-ANCHORS-3). Refutada: a linha 39 é
   `find "$CWD/.mnfs" -name EVIDENCE.md -path '*_chip*'`, agnóstica a profundidade. O
   `.mnfs/**/_chip-*/EVIDENCE.md` existe só na string de mensagem, não no mecanismo.
2. **"`find` falha com path Windows"** (hipótese do hub). Refutada por execução: `find` com
   `C:\Users\…\.mnfs` resolve e lista normalmente no Git Bash desta máquina.

## O que continua SEM explicação, e é honesto dizer

O bloqueio do CHIP-ANCHORS-3 **não** é explicado pelo defeito 2. Rodando a linha 39 verbatim contra
o worktree dele agora: **35** packs. A checagem de existência passa ali, então ele não deveria ter
recebido a mensagem de "no evidence pack".

Duas leituras sobrevivem e nenhuma está medida: (a) o `$CWD` do payload do hook não era o worktree
dele no momento do disparo; (b) o `.mnfs` dele estava em outro estado naquele instante. O dado que
decide é uma coisa só — o `cwd` que o hook recebeu naquele turno. Enquanto não existir, isto fica
**não fechado**, e o defeito 1 (certo) não deve ser usado para dar o caso por explicado.
