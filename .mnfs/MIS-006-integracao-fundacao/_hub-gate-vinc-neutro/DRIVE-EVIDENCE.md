# DRIVE-EVIDENCE — dado vivo, medido pelo HUB, para os assentos do gate

**Este artefato existe porque o assento não tem navegador.** O chip não pode subir servidor (§2 L2)
e o assento leitor recebe o DIFF, não o pack. Sem isto, o assento teria que INFERIR alcançabilidade
da mesma fonte de onde o chip inferiu — o gerador — e chegaria à mesma conclusão errada por
construção. Foi exatamente esse buraco que deixou uma premissa falsa sobreviver sete rodadas dos
DOIS lados do gate.

Medido pelo hub na `main` em `4f4377d1`, stack real, dado real de ERP e de Mercado Livre.
Nenhuma fixture, nenhum mock. Instalação `inst-mercado_livre-7e0d2125-…`.

```
GET /product-links/link-candidates?installation_id=…&limit=200

candidatos: 9
direções:   FOR 5 · AGAINST 6 · INCOMPARABLE 3 · UNAVAILABLE 9
side:       INCOMPARABLE/both 2 · INCOMPARABLE/provider 1
linhas sem reason: 0        linhas 100% INCOMPARABLE/UNAVAILABLE: 0
```

## As nove linhas, verbatim do wire

```
MLB3758161505  codprod=10741   70% MEDIA  seller_sku  [(seller_sku,FOR,) (ean,INCOMPARABLE,both)     (marca,UNAVAILABLE,)]
MLB4735304125  codprod=31704   70% MEDIA  seller_sku  [(seller_sku,FOR,) (ean,INCOMPARABLE,both)     (marca,UNAVAILABLE,)]
MLB4735326915  codprod=33698   70% MEDIA  seller_sku  [(seller_sku,FOR,) (ean,INCOMPARABLE,provider) (marca,UNAVAILABLE,)]
MLB4735378521  codprod=22467   40% BAIXA  seller_sku  [(seller_sku,FOR,) (ean,AGAINST,)              (marca,UNAVAILABLE,)]
MLB6896001442  codprod=42519   40% BAIXA  seller_sku  [(seller_sku,FOR,) (ean,AGAINST,)              (marca,UNAVAILABLE,)]
MLB4735378521  codprod=44975   20% BAIXA  ean         [(ean,AGAINST,)                                (marca,UNAVAILABLE,)]
MLB4735378521  codprod=22467   20% BAIXA  ean         [(ean,AGAINST,)                                (marca,UNAVAILABLE,)]
MLB6896001442  codprod=43534   20% BAIXA  ean         [(ean,AGAINST,)                                (marca,UNAVAILABLE,)]
MLB6896001442  codprod=42519   20% BAIXA  ean         [(ean,AGAINST,)                                (marca,UNAVAILABLE,)]
```

## O que isto DECIDE, e o assento não precisa acreditar em ninguém

### (a) A linha 100% INCOMPARABLE é INALCANÇÁVEL — o critério antigo era falso

`linhas 100% INCOMPARABLE/UNAVAILABLE: 0`, e mais forte: **toda** linha carrega
`marca UNAVAILABLE`, direção que a enumeração VELHA já cobria. A causa está no gerador e é
estrutural, não amostral: `mercado_livre` é a única declaração de capability da árvore e não supre
`marca`, então `resolveIdentityAnchors:149-169` põe `marca` UNAVAILABLE em TODO candidato. Célula
vazia exigiria um provider que suprisse as quatro âncoras. Não existe.

Producibilidade é fechada sobre os SÍTIOS do gerador. Alcançabilidade é fechada sobre as
DECLARAÇÕES que existem. Ler o gerador decide a primeira e é cego para a segunda.

### (b) O defeito REAL, e ele é aritmética — o assento pode conferir sozinho

Algoritmo da `main`, `QueueRow.tsx`:

```tsx
const shown  = [...byDirection("AGAINST"), ...byDirection("FOR"), ...byDirection("UNAVAILABLE")]
                 .slice(0, COMPACT_CHIP_LIMIT);   // COMPACT_CHIP_LIMIT = 2
const hidden = reasons.length - shown.length;
```

`INCOMPARABLE` **não está na enumeração**. Não é rankeado por último — não pode ser rankeado. É
inexibível a QUALQUER limite; subir `COMPACT_CHIP_LIMIT` não o traz. E `hidden` o CONTA.

Aplicado às nove linhas acima, e conferido contra o que o navegador renderizou:

| linha | `shown` calculado | `hidden` | tela renderizou | bate? |
|---|---|---|---|---|
| 70% (×3) | `[FOR, UNAVAILABLE]` | 1 | `✓ SKU · – Marca · +1` | sim |
| 40% (×2) | `[AGAINST, FOR]` | 1 | `✕ EAN … · ✓ SKU · +1` | sim |
| 20% (×4) | `[AGAINST, UNAVAILABLE]` | 0 | `✕ EAN … · – Marca` | sim |

Nas três linhas de 70%, os dois slots vão para um FOR e para `marca UNAVAILABLE` — ausência
**permanente, acionável por ninguém**, porque o provider nunca supre a âncora. O `ean INCOMPARABLE`,
sobre o qual o operador PODE agir (cadastrar o EAN no ERP, ou corrigir o anúncio), **não aparece**,
e o `+1` o conta sem que chip nenhum possa nomeá-lo. Um botão que promete o que a célula é incapaz
de exibir.

O comentário imediatamente acima da expressão afirma *"Ranking (never filtering) is what keeps at
least one motivo on screen"*. Com a quarta direção existindo, essa frase é **falsa na `main`**.

### (c) `side` colapsa: dois fatos, um texto

Os três `ean INCOMPARABLE` vivos, com o `detail` do wire:

```
side=both      detail=sem EAN para corroborar o CODPROD
side=both      detail=sem EAN para corroborar o CODPROD
side=provider  detail=sem EAN para corroborar o CODPROD
```

`both` = nenhum dos dois lados tem EAN → a ação é cadastrar no ERP.
`provider` = o ERP tem, o anúncio não → a ação é corrigir o anúncio.

Ações operacionais diferentes, **texto byte-idêntico**. Ler o campo `side` é a única coisa que
separa; parsear o português não separa, porque o português é o mesmo.

### (d) Defeito de BACKEND, fora do write-set deste chip

`MLB4735378521` aparece com `codprod=22467` **duas vezes** (40% via `seller_sku`, 20% via `ean`);
idem `MLB6896001442`/`42519`. Nove linhas para sete pares distintos.
`generation_service.go:365-377` percorre as duas âncoras e `uniqueProducts` deduplica só DENTRO de
uma. Consequência visível: duas linhas adjacentes do MESMO par renderizam GTIN `—` e GTIN
`✓ igual`, porque `QueueRow.tsx:191-194` deriva GTIN de `match_input` e está correto por linha.

**Não é deste chip.** Registrado aqui para que nenhum assento o atribua a este diff.

## Regra de uso

Este artefato é MEDIÇÃO, não alegação. Toda linha acima é reproduzível pelo hub e nenhuma depende
de prosa de chip. Um assento pode contradizê-lo — mas tem que contradizer o dado, nomeando qual
linha e qual campo, não o argumento.
