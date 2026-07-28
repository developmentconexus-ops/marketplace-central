# CORR-3 — a classe é maior que a minha CORR-2, e a saída de escape do docstring é falsa

`tip: 7cc17731` (inalterado, gate cortado) · medido pelo hub por EXECUÇÃO
· instrumento: `lane-r9/hub-06-branch-enumeration.mjs`, saída em `hub-06-branch-enumeration.log`

O chip devolveu, sob freeze e sem tocar na árvore, que a CORR-2 continua estreita. Confirmei por
execução, e a confirmação produziu um **segundo** achado que nenhum dos dois lados tinha.

## O erro da CORR-2 é o mesmo que ela corrigia, um nível acima

A CORR-2 alargou a classe de "colapso por CAIXA" para "código não-mapeado cuja string literal
iguala a saída TIPOGRAFADA de outro código". Ainda estreita: ela assume que **um dos lados da
colisão saiu da transformação**.

Meu instrumento da rodada 9 enumerou sobre uma lista escolhida a mão e agrupou por saída. Bom
agrupamento, população errada. O chip enumerou por **RAMO**, que é a população derivada do
mecanismo, e o instrumento novo faz o mesmo — para cada chave mapeada, o LITERAL de saída dela
também é um `provider_code` candidato, porque `registry.go:100-114` dedupa por igualdade exata de
string e nada proíbe esse código.

Medido:

```
COLLIDE "Amazon"             <- Amazon [VERBATIM]             | amazon [TYPESET]
COLLIDE "Amazon Marketplace" <- Amazon Marketplace [VERBATIM] | amazon_marketplace [TYPESET]
COLLIDE "Mercado Livre"      <- Mercado Livre [VERBATIM]      | mercado_livre [MAPPED]
COLLIDE "Shopee"             <- Shopee [VERBATIM]             | shopee [TYPESET]
classes em colisao: 4
```

A terceira linha é **MAPEADO × VERBATIM**, e o membro mapeado é `mercado_livre` — a única
declaração que existe hoje. Um adapter cujo `provider_code` fosse literalmente `"Mercado Livre"`
veste o nome do provider que já está em produção. Nenhuma das duas redações anteriores cobre esse
par. Falta ainda MAPEADO × MAPEADO, que ninguém pode instanciar com uma entrada só no `Record`.

## Redação que fecha — CORR-3

> A função tem N ramos e UM contradomínio. Injetividade é propriedade do CONTRADOMÍNIO, e a
> classe é todo par de códigos distintos cujas saídas coincidem, em QUALQUER combinação de ramos
> — inclusive dois ramos que nunca transformam nada. Estreitar domínio, mapear literal e cair em
> verbatim são três formas de PRODUZIR string; nenhuma é uma forma de RESERVAR string.

Redação do chip, ratificada verbatim. Alcance e gatilho inalterados: `mercado_livre` é a única
declaração e nenhuma linha viva alcança colisão. **R-1 continua REPORT.**

## D-3 · a saída de escape do docstring é chamada de injetiva e NÃO É — BLOQUEANTE por R-24

Isto é do hub, não do chip, e não estava em nenhuma das duas mensagens. `QueueRow.tsx`, no tip
cortado, fecha o parágrafo do R-1 assim:

> Only two shapes actually close it — render every unmapped code verbatim (**injective**,
> uglier), or require a display name in the registry beside the code (injective, and a contract
> change).

Rodei a opção (a) — `f(x) = map[x]` se mapeado, senão `x` — sobre a mesma população:

```
--- OPCAO (a): "render every unmapped code verbatim (injective, uglier)" ---
COLLIDE "Mercado Livre" <- Mercado Livre [VERBATIM] | mercado_livre [MAPPED]
classes em colisao sob a OPCAO (a): 1
opcao (a) NAO e injetiva
```

**A opção (a) não fecha nada.** Ela remove o ramo TYPESET e deixa MAPEADO × VERBATIM de pé, que é
justamente o par que a CORR-3 acabou de nomear. O docstring afirma `injective` em forma total
sobre um remédio que não tem a propriedade.

E o desenho da falha é o do resto da onda: o parágrafo foi escrito para ser HONESTO sobre um
defeito aberto, e a frase que vende o remédio alega mais do que o remédio entrega. É a quarta vez
nesta missão que um chip caçando universal falso escreve o seu — e desta vez dentro da própria
declaração de lacuna.

A opção (b) — nome de display no registry ao lado do código — é injetiva só se o registry
IMPEDIR dois códigos com o mesmo nome. Nada no texto diz isso, e o `Record` de hoje não impede.
Mesma frase, mesma forma, um grau mais fraco.

**Severidade: BLOQUEANTE por R-24**, independentemente de alcance: alegação total dispara em
código. Não muda o veredito sobre B-1 nem sobre B-2, que continuam fechados.

## O que isto NÃO muda

Os assentos estão rodando sobre este tip e **não foram informados** deste achado — a
independência deles é o instrumento e contaminá-la custaria mais do que ganharia. O achado fica
registrado aqui, datado, e entra no ruling junto com o que eles devolverem.
