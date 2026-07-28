# VEREDITO — CHIP-VINC-NEUTRO, rodada 8 · **BLOQUEADA, com AGREEMENT**

`tip julgado: bfc1d9bb` · `assento A: Opus frio → BLOCKED` · `assento B: GPT-5.6 Sol medium → BLOCKED`
Verbatim dos dois colados no retorno: `SEAT-A-opus-round8-VERBATIM.md`, `SEAT-B-sol-round8-VERBATIM.md`.

**P6-DUAL-GATE: AGREEMENT — os dois REPROVAM.** A rodada não fecha. Pela primeira vez nesta missão
o agreement é sobre BLOQUEAR, não sobre liberar, e os dois assentos chegaram lá por caminhos
diferentes: um achou defeito fora dos critérios, o outro dentro.

Nenhum bloqueante é sobre o fix central. **V8-1, V8-2 (lista 1), V8-3 e V8-5 passaram nos dois
assentos.** O que bloqueia é prosa que o código contradiz — duas vezes, nos dois casos em superfície
escrita pelo próprio chip nesta onda.

---

## B-1 · BLOQUEANTE — universal falso no docstring de `providerDisplayName`

Origem: assento A, fora de V8-1..V8-5. O assento declarou que estava fora e recusou-se a amolecer
por isso. Achado fora do critério só vira ordem se o HUB confirmar por observável próprio.
**Confirmei por EXECUÇÃO** (`HUB-VERIFY-round8.md`): função extraída do tip e rodada em `node`.

```
"amazon-marketplace"     -> "Amazon-marketplace"
```

O bloco diz, em forma total:

> Anything else — hyphens, mixed case, embedded spaces, anything that would lose a character —
> renders verbatim

**Hífen não renderiza verbatim.** Mecânica: `typesetSlug` faz `split("_")`, então um código com
hífen e sem underscore é um token só, é capitalizado, e o round-trip `painted.toLowerCase()` volta a
bater — o guard APROVA a transformação que a frase diz que ele recusa.

Bloqueia por R-24, independentemente de alcance: **alegação total dispara em CÓDIGO.** É a terceira
vez nesta missão que um chip caçando universal falso escreve o seu.

## B-2 · BLOQUEANTE — o buraco declarado é falso numa célula nomeada

Origem: **os DOIS assentos**, com severidades diferentes. Sol bloqueou; Opus rebaixou a REPORT
argumentando que a lista 1 enuncia a regra, então as duas listas JUNTAS ficam corretas.

Verificado por string, e o quadro é mais preciso do que qualquer um dos dois disse:

```
wireFixtures.ts:294   if (reason.side !== undefined && reason.direction !== "INCOMPARABLE") fail(...)
wireFixtures.ts:450   "does NOT check ... the direction/side of an absence on a SUPPLIED anchor"
```

A frase é uma conjunção e ela é falsa em **uma** célula, não em todas:

| caso | frase diz | código faz | |
|---|---|---|---|
| `side` numa ausência UNAVAILABLE | passa | **LANÇA** (`:294`) | **FALSO** |
| `side` numa ausência INCOMPARABLE | passa | passa | verdadeiro |
| a DIREÇÃO da ausência em âncora suprida | não checa | não checa | verdadeiro |

**Ruling: BLOQUEANTE, e eu fico com o Sol.** A razão contra o steelman do Opus, dita inteira: este
docstring é a declaração AUTORITATIVA de escopo — o chip o escreveu exatamente para que ninguém
tivesse de fazer engenharia reversa do guard. Quem consulta a lista do que NÃO é checado para saber
o que passa recebe resposta falsa. A lista 1 enunciar a regra não repara isso; significa que **o
docstring contradiz a si mesmo**, o que é pior do que estar incompleto, porque o leitor não tem
regra para decidir qual lista vence.

Reparo é uma cláusula: dizer "o VALOR de `side` numa ausência INCOMPARABLE, e a DIREÇÃO de uma
ausência em âncora suprida". O crédito da forma escopada fica de pé — este bloco é o primeiro desta
missão a passar no R-24 por construção, e o escorregão é num item, não no desenho.

---

## R-1 · REPORT — colapso por CAIXA, produzível e NÃO alcançável

`"amazon"` e `"Amazon"` renderizam os dois `"Amazon"` (medido). É exatamente o dano que o docstring
nomeia — *two providers wearing one name is wrong information* — e o guard não pega, porque
`restored = painted.toLowerCase()` destrói a caixa antes de comparar.

**Não é alcançável nesta árvore.** Dois códigos diferindo só por caixa exigem dois adapters
registrados; existe um (`mercado_livre`), e ele está no mapa literal, então nem chega ao round-trip.
Produzível e não alcançável — a distinção que ESTA rodada ratificou, aplicada agora contra um achado
que era mais cômodo deixar bloqueante.

O assento não errou ao chamá-lo de bloqueante: sem shell, sem browser, Go fora do read-set, ele não
tinha instrumento para medir alcance, e disse isso na linha 5 dos próprios limites. Quem tinha era o
hub.

**Gatilho de expiração:** vira BLOQUEANTE no dia em que um segundo adapter registrar capacidade.

> **CORR-2 · 2026-07-28 · a CLASSE aqui está estreita demais.** "Colapso por caixa" é uma
> instância. A medição do hub na rodada 9 achou uma segunda (`"Amazon Marketplace"` colide com
> `"amazon_marketplace"`), então a classe é *código não-mapeado cuja string literal iguala a saída
> tipografada de outro código*. Alcance e gatilho não mudam. Reescrito em `HUB-VERIFY-round9.md`.

## R-2 · REPORT — o guard não descobre adapter novo

**Os dois assentos**, mesma coisa: `DECLARED_PROVIDER_CAPABILITIES` / `GO_SEAM` leem só o adapter do
`mercado_livre`, então um `shopee/capability_adapter.go` novo deixa os oito testes verdes com a
tabela obsoleta. **Sem ordem:** o gatilho de expiração já está escrito na árvore, que é a forma que
o V8-4 pede. Guard que não se auto-descobre, com a condição de expiração declarada, é dívida
honesta, não falsidade.

---

## Limite de assento que o HUB descarregou

Sol declarou: *"The TypeScript lane has 12 errors outside `pages/vinculos`; their claimed
pre-existence was not independently reproducible under the read-only boundary."* Correto, e ele não
tinha como. Eu tinha. Medido na `main` de hoje:

```
npx tsc -p apps/web/tsconfig.json --noEmit   →  15 erros
                            em pages/vinculos →   3
```

Os três de `pages/vinculos` são `QueueRow.tsx:34`, `:75` e `VinculoDrawer.tsx:118` — todos
`INCOMPARABLE` faltando num `Record<ProductLinkReasonDirection, …>`, que é **exatamente o que este
chip conserta**. 15 − 3 = **12**, e 12 é o número da lane dele. A alegação de pré-existência está
corroborada por medição independente, e o número bate na casa.

## Ordem para a rodada 9

1. **B-1** — a frase sobre hífen/caixa. Duas saídas, escolha do chip porque o arquivo é dele:
   restringir a aplicação ao domínio onde a transformação É injetiva (`^[a-z0-9]+(_[a-z0-9]+)*$`),
   o que faz caixa e hífen caírem em verbatim e torna a frase verdadeira sem reescrevê-la; ou
   DELETAR a frase falsa (R-25) e declarar o escopo real. A primeira fecha B-1 e R-1 juntos.
   > **CORR-1 · 2026-07-28 · a última frase é FALSA.** O chip mediu antes de aceitar e o hub
   > confirmou por execução em `HUB-VERIFY-round9.md`: a restrição de domínio fecha B-1 e **não**
   > fecha R-1. Verbatim é a IDENTIDADE, e a identidade divide contradomínio com a transformação —
   > `Amazon` sai do domínio e continua pintando `Amazon`, que é o que `amazon` produz. Estreitar o
   > domínio move a entrada para o ramo onde o parceiro da colisão mora; não separa nada. A frase
   > fica registrada porque foi ela que a rodada 9 executou; o veredito dela não fica.
2. **B-2** — a cláusula do buraco declarado, com as três células acima decididas por string.
3. **R-1 e R-2** não são ordem. Se a saída escolhida em B-1 fechar R-1, diga; se não, R-1 fica
   registrado com o gatilho.
4. **O teste de colapso** (`patch:1399-1421`) alimenta variante de whitespace e nunca de caixa.
   Ele está certo sobre o que testa; se B-1 fechar por domínio, o braço de caixa vira barato e a
   frase passa a ser sustentada pelo teste, não só pelo mecanismo.
