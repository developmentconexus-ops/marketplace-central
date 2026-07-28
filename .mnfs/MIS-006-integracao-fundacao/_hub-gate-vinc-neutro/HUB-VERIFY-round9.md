# Verificação do HUB — rodada 9 · **o chip refutou uma ordem minha, e ele está certo**

`tip: 7cc17731` · medido pelo hub por EXECUÇÃO, antes de aceitar a mensagem do chip.

A rodada 9 devolveu uma correção à ordem da rodada 8. O chip não a aceitou de véspera: mediu
antes. A regra da casa vale nos dois sentidos — achado do chip contra o hub também só vira fato
se o HUB confirmar por observável próprio. Confirmei, e o meu instrumento vai além do dele.

## CORR-1 · A minha frase era falsa

`VERDICTS-round8.md:119` dizia:

> A primeira fecha B-1 e R-1 juntos.

**É falso.** Função extraída de `7cc17731:apps/web/src/pages/vinculos/QueueRow.tsx` e rodada em
`node`, com o conjunto de entradas do hub ampliado:

```
"mercado_livre"          -> "Mercado Livre"
"amazon"                 -> "Amazon"
"Amazon"                 -> "Amazon"          <-- R-1 SOBREVIVE
"amazon-marketplace"     -> "amazon-marketplace"   <-- B-1 FECHADO
"amazon_marketplace"     -> "Amazon Marketplace"
"Amazon_Marketplace"     -> "Amazon_Marketplace"
"amazon__market"         -> "amazon__market"
"_amazon"                -> "_amazon"
"amazon_"                -> "amazon_"
"shopee"                 -> "Shopee"
"Shopee"                 -> "Shopee"
"AMAZON"                 -> "AMAZON"
"amazon marketplace"     -> "amazon marketplace"
"amazon2"                -> "Amazon2"
"2amazon"                -> "2amazon"
"Amazon Marketplace"     -> "Amazon Marketplace"
```

Classes em colisão, computadas sobre a saída e não julgadas a olho:

```
"Amazon"              <- ["amazon", "Amazon"]
"Amazon Marketplace"  <- ["amazon_marketplace", "Amazon Marketplace"]
"Shopee"              <- ["shopee", "Shopee"]
total de classes em colisao: 3
```

## Por que eu errei, dito com precisão

Eu tratei "cair em verbatim" como se verbatim fosse um destino SEGURO. Não é: verbatim é a
IDENTIDADE, e a identidade divide contradomínio com a transformação. Estreitar o domínio não
separa `amazon` de `Amazon` — só move `Amazon` para o ramo onde ele continua batendo com a saída
que `amazon` produz. A colisão não foi eliminada; **mudou de lugar**, de dentro da transformação
para a costura entre a transformação e a própria escotilha dela.

## O que a minha medição acrescenta à do chip

O chip mediu a colisão por CAIXA. A terceira classe acima mostra que **caixa é só uma instância**:
`"Amazon Marketplace"` (com espaço, fora do domínio, identidade) colide com `"amazon_marketplace"`
(dentro do domínio, transformado). A forma geral é a que o chip enunciou — *toda string que a
transformação produz é também uma string que o fallback produz* — e ela tem pelo menos duas
instâncias distintas, não uma.

**Consequência para o R-1:** o registro dele na rodada 8 diz "colapso por CAIXA", e isso é
estreito demais. A classe correta é *código não-mapeado cuja string literal iguala a saída
tipografada de outro código*. Reescrito abaixo.

## R-1 · reescrito, continua REPORT, gatilho inalterado

Colisão entre o ramo identidade e a imagem da transformação. Duas instâncias medidas (caixa;
substituição de separador por espaço). Continua **não alcançável**: `mercado_livre` é a única
capacidade declarada e está no mapa literal, então nenhuma linha viva chega a qualquer um dos
dois ramos. **Gatilho: o segundo adapter registrado.** As duas saídas que fecham de verdade
(verbatim universal; display name no registry ao lado do código) estão fora do write-set deste
chip — a segunda é mudança de contrato, e contrato é do hub.

## Custódia — dois pontos, e um deles precisa de atribuição

O chip verificou com `main...HEAD` (base de merge). O gate de MERGE mede contra o **tip do alvo**,
e aí aparece um caminho a mais:

```
git diff --name-only main 7cc17731 -- ':!.mnfs'   ->  + docs/HARNESS-PROFILE.md
```

Atribuído antes de virar achado:

```
git diff 4f4377d1 7cc17731 -- docs/HARNESS-PROFILE.md   ->  VAZIO   (o chip não tocou)
git diff 4f4377d1 main     -- docs/HARNESS-PROFILE.md   ->  +96 -1  (a MAIN avançou)
git diff 4f4377d1 main     -- apps/web/src/pages/vinculos/  ->  VAZIO
```

Então: nenhum revert. O arquivo é avanço da própria main, e o merge 3-way não o desfaz porque o
chip não tem commit nele. E como a main não tocou `pages/vinculos` desde a base, **dois pontos e
três pontos são idênticos na superfície revisada** — a diff do pack serve para o gate de merge sem
ressalva.

## O que o hub NÃO mediu nesta rodada

A atribuição do MUST-FAIL 5 (mutação `:294 -> if (false && …)`, 1 vermelho isolado) e as lanes
estão no pack do chip como log commitado, não como prosa. O assento lê o log; o hub não o
reproduziu. Declarado, não escondido.
