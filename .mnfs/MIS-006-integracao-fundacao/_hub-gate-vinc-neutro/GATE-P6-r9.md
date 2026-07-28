# GATE P6 — CHIP-VINC-NEUTRO, rodada 9 · contrato de validação

`tip julgado: 7cc17731` · `base: main @ cf018d81 (TIP do alvo, não merge-base)`
`patch: r9-code-diff.patch` — 3122 linhas, 11 arquivos, sha256
`90dfc51aeb5e113c8f57e3fe3dda1e584614fa58e7ee6fd539ac9f7b2616198a`

## Suas entradas, e são exaustivas

1. `r9-code-diff.patch` — `git diff main 7cc17731 -- apps/web/src/pages/vinculos/`.
2. Este contrato, verbatim.
3. `lane-r9/` — saídas CRUAS medidas pelo HUB no tip do chip, não pelo chip.

Você **não** recebe o pack do chip e não deve pedi-lo. Achado sobre a prosa do pack é
inatingível daqui por construção, e isso é intencional. Se um achado de CÓDIGO depender do
raciocínio do chip, peça uma seção NOMEADA e limitada — nunca a árvore.

## Estado: esta é uma rodada de FECHAMENTO após BLOQUEIO

A rodada 8 foi bloqueada pelos dois assentos, por prosa que o código contradizia. Nada do que
bloqueou tocava o fix central. As duas ordens estão respondidas nesta rodada; um terceiro item
(R-1) foi expressamente deixado ABERTO com gatilho, e a honestidade dessa declaração é critério.

## Fatos que o HUB mediu — verbatim, e você pode refutá-los nomeando linha e valor

Estes não são alegações do chip. Foram medidos pelo hub no `7cc17731`, e o instrumento viaja
junto para você poder contestá-lo (`lane-r9/hub-05-…-instrument.mjs`).

`providerDisplayName` executada em `node` (`lane-r9/hub-05-…-execution.log`):

```
"amazon-marketplace" -> "amazon-marketplace"     "amazon" -> "Amazon"
"Amazon"             -> "Amazon"                 "amazon_marketplace" -> "Amazon Marketplace"
"Amazon Marketplace" -> "Amazon Marketplace"     "amazon__market" -> "amazon__market"

classes em colisão: 3
  "Amazon"             <- ["amazon", "Amazon"]
  "Amazon Marketplace" <- ["amazon_marketplace", "Amazon Marketplace"]
  "Shopee"             <- ["shopee", "Shopee"]
```

Lanes medidas pelo hub no tip do chip, com o EXIT do PROCESSO (não do `sed` do pipe):

```
vitest src/pages/vinculos    6 files, 42 tests, 0 failed          EXIT(vitest)=0
vitest apps/web (completo)  67 files, 551 tests, 0 failed         EXIT(vitest)=0
tsc -p tsconfig --noEmit    12 errors, 0 em pages/vinculos        EXIT(tsc)=2
```

Atribuição do MUST-FAIL 5, mutação isolada `:294 -> if (false && …)`, refeita pelo hub
(`lane-r9/hub-04-mutation.log`): **1 failed | 41 passed (42)**, e o vermelho é
`MUST-FAIL 5 — an UNAVAILABLE absence cannot carry a \`side\``. Restauração conferida por md5
(`bdf56153…` dos dois lados) e `git diff --quiet HEAD` = 0.

Custódia: contra o TIP da main o único caminho fora de `pages/vinculos` é
`docs/HARNESS-PROFILE.md`, e ele é avanço da PRÓPRIA main — o chip não tem commit nele, e a main
não tocou `pages/vinculos` desde a base. Não é revert. Não está no seu patch.

---

## Critérios

### V9-1 · O universal do `providerDisplayName` é verdadeiro do código

O bloco acima da função faz afirmação TOTAL sobre o que renderiza verbatim. Decida se ela é
verdadeira **sobre todo o domínio de entrada da função**, não só sobre o contraexemplo que a
rodada 8 usou. O guard agora rejeita ANTES de transformar (`INJECTIVE_PROVIDER_SLUG`), e o chip
declara explicitamente que não sobrou round-trip atrás dele.

Falha se: alguma frase total do bloco for falsa em alguma entrada; ou se a regex admitir uma
entrada que a frase diz que cai em verbatim; ou se a frase declarar uma lacuna que não existe
(**buraco declarado falso é tão falso quanto garantia não cumprida, e aponta para o outro lado**).

### V9-2 · O parágrafo "what it does NOT check" é verdadeiro nas TRÊS células

`wireFixtures.ts` — o texto foi corrigido nesta rodada. As três células que a rodada 8 decidiu
por string:

| caso | veredito da rodada 8 |
|---|---|
| `side` numa ausência UNAVAILABLE | o guard LANÇA (`:294`) |
| `side` numa ausência INCOMPARABLE | passa |
| a DIREÇÃO da ausência em âncora suprida | não é checada |

Falha se o texto novo continuar falso em qualquer célula, ou se o texto ficar auto-contraditório
(era esse o veredito: docstring que se contradiz é pior que incompleto, porque o leitor não tem
regra para saber qual lista vence).

### V9-3 · O MUST-FAIL 5 sustenta a frase, e sustenta só o que isola

Arm novo em `wireFixtures.guard.test.ts`. Falha se ele passar por razão diferente da regra que a
frase nomeia (asserção por mensagem genérica, fixture que já violaria outra regra antes, etc.).
A atribuição por mutação está medida acima pelo hub — você julga o DESENHO do arm, não precisa
reproduzir a execução.

### V9-4 · R-1 declarado como lacuna HONESTA, não fechado por prosa

O chip declara que a colisão entre o ramo transformado e o ramo identidade **não** está fechada,
nomeia por que estreitar o domínio não pode fechá-la, nomeia as duas formas que fechariam e por
que estão fora do write-set, e nomeia o gatilho.

Falha se: o texto alegar fechamento; ou se o alcance declarado (`mercado_livre` é a única
capacidade declarada, mapeada literal, nenhuma linha viva chega a qualquer ramo) for falso no
diff; ou se o gatilho não for verificável.

### V9-5 · O teste não grava o defeito como requisito

Achado do próprio chip: a asserção V10 afirmava `getByText("Amazon-marketplace")` — a saída
ERRADA do guard antigo — e ficava verde por isso, que foi como o caso do hífen sobreviveu. Foi
corrigida. O chip também declara ter OMITIDO deliberadamente a asserção de caixa, para não
gravar o defeito R-1 sobrevivente como requisito, e declara a omissão no próprio teste.

Falha se: alguma asserção nova codificar a saída atual de um defeito declarado aberto; ou se a
omissão declarada não estiver no teste com a razão; ou se o teste alegar cobertura que não tem.

### V9-6 · Varredura de CLASSE sobre todo o diff — não pare nos sítios que eu nomeei

Um brief que NOMEIA o sítio ensina o assento a parar nele. Varra os 11 arquivos do patch pela
CLASSE: prosa de totalidade (`never`/`always`/`only`/`every`/`cannot`/`no longer`/`verbatim`/
`renders`/`unreachable`) e lacuna declarada. Um veredito sem seção de SWEEP está incompleto na
cara. Declare o escopo que varreu — aprovação corrobora só dentro do escopo declarado.

### V9-7 · Nada regrediu no fix central

O fix é `INCOMPARABLE` ganhando par de token próprio, ranking sem filtrar, e as leituras
tolerantes a drift do wire. Falha se o diff regredir qualquer um, ou tocar caminho fora de
`apps/web/src/pages/vinculos/`.

---

## Regras do assento

- **Reporte tudo.** Amolecer no assento é a falsidade que o R-25 proíbe. A reclassificação
  BLOQUEANTE/REPORT é do HUB, com verificação por string.
- **BLOQUEANTE** é reservado a OBSERVÁVEL errado — comportamento, segurança, dado, contrato
  publicado. Discriminador: suba o código com o achado intocado e nomeie o que um usuário,
  operador, chamador ou linha gravada faz de diferente. Sem resposta = REPORT. Universal falso em
  código é BLOQUEANTE por R-24 independentemente de alcance — a alegação é o dano.
- Achado **fora** destes critérios: traga, e diga que está fora. Não amoleça por isso. Ele só
  vira ordem se o hub confirmar por observável próprio, e foi assim que a rodada 8 bloqueou.
- Cite `arquivo:linha` e **cole o excerto**. Cenário de falha CONCRETO, com entrada e saída.
- Declare seus LIMITES: o que não conseguiu medir com o instrumento que tem.
- Formato: veredito (`BLOCKED` / `CLEAR`), tabela `severidade | critério | file:line | defeito |
  cenário concreto`, seção SWEEP com escopo declarado, seção LIMITES.
