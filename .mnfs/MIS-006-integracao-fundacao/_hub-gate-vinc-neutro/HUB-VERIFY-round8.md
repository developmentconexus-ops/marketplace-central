# Verificação do HUB — achado bloqueante do assento A, rodada 8

O assento A bloqueou **fora de V8-1..V8-5**, declarou que estava fora, e disse que não amoleceria
por isso. Regra da casa: achado fora do critério só vira ordem se o HUB o confirmar por observável
próprio. Confirmei — **por EXECUÇÃO, não por leitura.**

## Instrumento

Extraí a função do tip do chip (`bfc1d9bb:apps/web/src/pages/vinculos/QueueRow.tsx`) e a rodei em
`node` sobre as entradas em disputa. Não é releitura do código do assento: é a função rodando.

```
"mercado_livre"          -> "Mercado Livre"
"amazon"                 -> "Amazon"
"Amazon"                 -> "Amazon"          <-- COLAPSO
"amazon-marketplace"     -> "Amazon-marketplace"
"amazon_marketplace"     -> "Amazon Marketplace"
"Amazon_Marketplace"     -> "Amazon_Marketplace"
"amazon__market"         -> "amazon__market"   <-- guard de whitespace FUNCIONA
"amazon_market"          -> "Amazon Market"
```

## São DOIS defeitos, e eles têm classes diferentes

### D-1 · UNIVERSAL FALSO no docstring — BLOQUEANTE por R-24

O bloco acima da função diz, em forma total:

> Anything else — hyphens, mixed case, embedded spaces, anything that would lose a character —
> renders verbatim

**Hífen não renderiza verbatim.** `amazon-marketplace` → `Amazon-marketplace`, medido acima. A causa
é mecânica: `typesetSlug` faz `split("_")`, então um código com hífen e sem underscore vira um único
token, é capitalizado, e o round-trip `painted.toLowerCase()` volta a bater com o código de entrada —
o guard aprova a transformação que a frase diz que ele recusa.

Isto é BLOQUEANTE independentemente de alcance, porque R-24 é sobre a alegação, não sobre o
usuário: **alegação total dispara em CÓDIGO.** E é a mesma forma que esta onda inteira existe para
matar — o chip que caça universal falso escreveu o seu, terceira vez nesta missão.

### D-2 · Colapso por CASE — REPORT hoje, e a razão é a desta rodada

`"amazon"` e `"Amazon"` renderizam os dois como `"Amazon"`. É exatamente o dano que o próprio
docstring nomeia ("two providers wearing one name is wrong information") e o guard não pega, porque
`restored = painted.toLowerCase()` destrói a informação de caixa antes de comparar.

**Não é alcançável nesta árvore.** Um único adapter declara capacidade (`mercado_livre`), e ele está
no mapa literal, então nem passa pelo round-trip. Dois códigos diferindo só por caixa exigem dois
adapters registrados. É **PRODUZÍVEL e não ALCANÇÁVEL** — a distinção que esta rodada inteira
ratificou, aplicada agora contra um achado que favorecia bloquear.

O assento não errou ao chamá-lo de bloqueante: ele **não tinha** como medir alcance (sem Bash, sem
browser, Go fora do read-set) e disse isso na linha 5 dos limites dele. Quem tinha o instrumento era
o hub.

## Consequência

D-1 e D-2 têm o mesmo conserto barato, então a separação não muda o trabalho — muda o que a ordem
pode exigir. O guard erra porque testa round-trip em vez de testar o DOMÍNIO em que a transformação
é injetiva. Restringir a aplicação a slugs que casem `^[a-z0-9]+(_[a-z0-9]+)*$` faz caixa e hífen
caírem em verbatim, que é o que a frase promete, e aí a frase vira verdadeira sem ser reescrita.

Alternativa honesta se o chip preferir: manter o mecanismo e **DELETAR a frase falsa** (R-25),
declarando o escopo real. As duas saídas servem; escolher é do chip, porque o dono do arquivo é ele.

## O que o teste do chip não cobre, e por quê importa

`patch:1399-1421` — `does not let two provider codes collapse onto one name through whitespace` —
passa e continuaria passando com D-2 intacto: ele alimenta variante de whitespace, nunca variante de
caixa. Guard parcial sob frase total, de novo. O teste está certo sobre o que testa; a frase é que
alega mais do que ele sustenta.
