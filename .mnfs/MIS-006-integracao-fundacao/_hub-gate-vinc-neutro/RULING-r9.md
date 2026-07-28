# RULING — CHIP-VINC-NEUTRO rodada 9 · **BLOQUEADO**

`tip julgado: 7cc17731` · gate `GATE-P6-r9.md` · assentos: A (Opus, `harness:gate-reviewer`) e
B (Sol medium, `codex:codex-rescue`), independentes, sem contato · **os dois BLOCKED**

O hub não encaminha veredito de assento. Cada achado abaixo que decide alguma coisa foi
**re-medido por execução na casa** antes de entrar. Instrumento: `lane-r9/hub-06-proto.mjs`.

---

## D-1 · lookup em objeto literal aceita membro herdado — BLOQUEANTE (V9-1, V9-6)

Os dois assentos acharam isto **por caminhos separados**: A partiu da prosa de totalidade e
perguntou o que a viola; B varreu a forma `X[chaveDoFio]`. Convergência independente em cima do
mesmo mecanismo é o sinal mais forte que este gate produz.

`providerDisplayNames` é objeto literal, então herda `Object.prototype`. O `if (mapped)` aceita
qualquer membro herdado ANTES de o domínio ser testado. Medido:

```
toString               typeof: function   verbatim? false
constructor            typeof: function   verbatim? false
__proto__              typeof: object     verbatim? false
valueOf / hasOwnProperty / isPrototypeOf / toLocaleString / propertyIsEnumerable   idem
```

**O tipo de retorno declarado é `string` e é violado em runtime.** A prosa em `QueueRow.tsx:103`
promete que qualquer coisa fora do domínio "renders verbatim"; nenhuma das oito renderiza
verbatim.

A classe não é de uma função. Assento A enumerou os sítios e o consequente de cada um; o pior é
`statusDecidingAnchors[candidate.match_status]` (`:229`), porque é justamente o reader cujo
comentário promete que um status desconhecido derruba **uma célula e não a tabela**: com
`match_status: "toString"` o `if (!rule)` é contornado, `rule(...)` devolve `"[object Undefined]"`,
e `decided.join(" + ")` estoura **durante o render** — exatamente o desfecho que o comentário diz
impedir. Mesma classe em `wireFixtures.ts:270`, onde a rejeição documentada vira `TypeError` sem
atribuição.

### Alcance — MEDIDO pelo hub, e ele corta este achado ao meio

O hub mediu de onde `provider_code` sai:

```
connectors/adapters/mercado_livre/capability_adapter.go:81   ProviderCode: "mercado_livre",
```

Literal Go. Constante → coluna `product_links.provider_code` → API → FE. **Não há caminho de
input não-confiável**: nem operador, nem usuário, nem a API do ML alimentam esse campo. Idem
`match_status`, `direction`, `side`, `confidence_band` — todos enum de Go.

Para `"toString"` chegar ali, um dev precisa escrever `ProviderCode: "toString"` num adapter
novo. Isso é **exatamente** o mesmo gatilho do R-1, que este mesmo ruling arquiva como REPORT.
Tratar um como BLOQUEANTE e o outro como REPORT é dois pesos para um balde. Corrigido:

- **BLOQUEANTE, e independe de alcance**: a frase `renders verbatim` (`:103`) é falsa sobre o
  domínio que ela mesma declara, e o tipo de retorno `: string` devolve `function`. R-24 dispara
  pela frase e pela assinatura. Conserto é `Object.hasOwn` — uma linha por sítio.
- **REPORT com gatilho, rebaixado**: a queda da tabela em `statusDecidingAnchors` e os outros ~9
  sítios. Mecanismo confirmado por execução, gatilho idêntico ao do R-1 (o segundo adapter), e
  arquivado ao lado dele.

## D-2 · a justificativa da lacuna é falsa — BLOQUEANTE (V9-1)

`QueueRow.tsx:139`:

> the fallback is identity, so **every** string a transform can produce is also a string the
> fallback can produce.

Achado do assento B. Medido:

```
2amazon    TYPESET -> "2amazon"    verbatim consegue emitir essa saida? false
0          TYPESET -> "0"          verbatim consegue emitir essa saida? false
9x         TYPESET -> "9x"         verbatim consegue emitir essa saida? false
```

O verbatim só emite string que **falha** no regex. Saída do typeset que **passa** no regex é
inalcançável por ele. A frase é total e tem contraexemplo.

## D-3 · a saída de escape do próprio parágrafo não é injetiva — BLOQUEANTE (V9-1)

**Achado do hub. Nenhum dos dois assentos chegou nele.** Mesmo parágrafo, `QueueRow.tsx:141`:

> render every unmapped code verbatim (**injective**, uglier)

Rodei a opção (a) — `f(x) = map[x]` se mapeado, senão `x`:

```
COLLIDE "Mercado Livre" <- Mercado Livre [VERBATIM] | mercado_livre [MAPPED]
classes em colisao sob a OPCAO (a): 1
```

Não é injetiva. Sobra MAPEADO × VERBATIM, que é o par que a CORR-3 nomeia. O remédio anunciado
não fecha o defeito que o parágrafo declara aberto.

### O parágrafo inteiro

D-2, D-3 e o REPORT do assento A ("**Only two shapes** actually close it" — existe uma terceira,
tornar os contradomínios disjuntos) estão nas MESMAS dez linhas. Três totais falsos no parágrafo
que foi escrito para ser honesto sobre um defeito aberto. Não é coincidência de redação: quem
declara lacuna escreve com a mesma confiança de quem declara garantia, e a lacuna não é auditada
como garantia. Vira emenda de perfil.

## D-4 · o reparo consertou o parágrafo errado — BLOQUEANTE (V9-2)

Os dois assentos. Confirmado por leitura dos dois sítios no `7cc17731`:

`wireFixtures.ts:79`, lista **WHAT IT DOES NOT CHECK** do ARQUIVO:
> - the DIRECTION and `side` of an absence on a SUPPLIED anchor.

`wireFixtures.ts:455`, parágrafo do `wireCandidate`:
> The `side` of an absence is NOT in that list… `side` **is checked**… a `side` on an UNAVAILABLE
> absence THROWS

Mesmo arquivo, um diz que não checa, o outro diz que checa, sem regra de precedência. E
`wireFixtures.guard.test.ts` (MUST-FAIL 5), do mesmo commit, **afirma o throw**. O diff refuta o
próprio docstring.

Isto é a regra do perfil literal: **guard parcial sob frase total é pior que guard nenhum.** O
chip corrigiu a instância que a rodada 8 citou e deixou a mesma falsidade uma docstring acima.

## D-5 · dois títulos de teste alegam o universal que o R-1 declara aberto — BLOQUEANTE (V9-5)

Os dois assentos. Confirmado por string:

```
589:  it("typesets an unmapped provider without ever collapsing two distinct codes onto one name"
1046: it("does not let two provider codes collapse onto one name"
```

A exceção existe só em comentário. Nenhuma lane imprime comentário; toda lane imprime título. Um
leitor do vitest verde lê garantia de não-colapso que a suíte não sustenta — e o hub mediu 4
classes em colisão. Título de teste é superfície publicada.

## CORR-3 · redação da classe R-1, ratificada

> A função tem N ramos e UM contradomínio. Injetividade é propriedade do CONTRADOMÍNIO, e a
> classe é todo par de códigos distintos cujas saídas coincidem, em QUALQUER combinação de ramos
> — inclusive dois ramos que nunca transformam nada. Estreitar domínio, mapear literal e cair em
> verbatim são três formas de PRODUZIR string; nenhuma é uma forma de RESERVAR string.

Redação do chip, verbatim. Detalhe em `HUB-VERIFY-round9-CORR3.md`. **R-1 continua REPORT** —
alcance e gatilho inalterados.

## B-1 e B-2

Continuam **fechados**. D-1..D-5 não os reabrem. V9-3 PASS pelos dois assentos: MUST-FAIL 5 é
isolado, o regex bate em um único `fail()` do arquivo, e a mutação deu 1 vermelho de 42.

## O `550`

Pergunta do chip respondida: **emenda no r10.** A árvore vai se mover de qualquer jeito, então o
argumento de "mover o tip debaixo de gate cortado" morreu. Corrige os dois arquivos de pack no
mesmo commit do reparo e declara o delta; a mensagem imutável de `7cc17731` fica coberta pela
correção datada. **São dois números, não um** — o `tsc 12` tem o mesmo vício e sobreviveu por
sorte. Ambos re-medidos como última coisa antes do `git add`.

## Peso desta rodada, dito sem inflar

Depois do corte de alcance acima, o que resta nesta rodada é **quase todo prosa**, e o hub diz
isso em voz alta em vez de deixar a contagem de BLOQUEANTES sugerir outra coisa:

| achado | efeito com os dados que já estão no banco |
|---|---|
| D-1 (frase + assinatura) | nenhum — some quando um dev escreve um provider_code esquisito |
| D-2, D-3 | nenhum — propriedade de função sobre entrada que não ocorre |
| D-4 | real para o próximo dev: duas docstrings do mesmo arquivo se contradizem |
| D-5 | real para quem lê a lane: título verde alega garantia que a suíte não tem |
| R-1, e o resto da classe protótipo | REPORT com gatilho no segundo adapter |

Nenhum defeito desta rodada é visível ao operador hoje. Isso **não** desfaz o BLOCKED — R-24 é a
regra da missão e frase falsa em código dispara por si — mas o r10 é uma rodada curta de redação
e um `Object.hasOwn`, não um redesenho. Dimensione assim.

## Não-evidenciados que NÃO viram achado

Registrados como limite do instrumento, não como aprovação: "só `mercado_livre` é declarado"
(nenhum assento nem lane enumera adapters registrados); `BatchPreviewModal.tsx` sob
`pages/vinculos/` e fora do patch pode ainda ter tabela privada de tokens; a alegação de importe
em `VinculosPage.test.tsx:41` é estática sobre árvore fora do diff. **Os três viram critério do
r10**, medidos por string, não declarados.
