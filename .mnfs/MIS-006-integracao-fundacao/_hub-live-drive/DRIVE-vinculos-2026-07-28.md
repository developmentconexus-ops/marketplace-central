# Live drive — /vinculos na `main` de hoje, dado real

Ordem do operador: *"Porque nçao deixou subir servidor? Regra é não ter fixture e sempre teste
real lembra"*. Stack subida pelo hub, frontend re-apontado para o checkout primário, `/vinculos`
dirigida com dado real. Sem fixture, sem mock, sem `vi.mock`.

Base: `main` @ `4f4377d1`. Nenhum chip da onda 2 mergeado. Zero erro de console.

```
GET /product-links/link-candidates?installation_id=inst-mercado_livre-7e0d2125-…&limit=200
candidatos: 9   linhas sem reason: 0   linhas 100% INCOMPARABLE/UNAVAILABLE: 0
direções: FOR 5 · AGAINST 6 · INCOMPARABLE 3 · UNAVAILABLE 9
side por direção: (INCOMPARABLE,both) 2 · (INCOMPARABLE,provider) 1
```

## O que o drive comprou

### V-1 não é alcançável com dado real

O defeito manchete de 7 rodadas do CHIP-VINC-NEUTRO — a linha sem chip nenhum quando todo motivo
é INCOMPARABLE/UNAVAILABLE — **não ocorre**. Toda linha carrega `marca UNAVAILABLE`, que ESTÁ em
`shown`. Sete rodadas de gate sobre um caso que o produtor não emite.

Isso não retira o valor do chip (ver F-3), mas mede o custo: o alvo foi escolhido pela fixture,
não pelo dado.

### F-1 — duplicado entre âncoras (BACKEND, não nomeado em nenhuma das ~18 rodadas)

`generation_service.go:365-377`, `buildCollisionCandidates`. O laço percorre as DUAS âncoras
(seller_sku, ean) e emite um candidato por produto por âncora. `uniqueProducts` (:367) deduplica
**dentro** de uma âncora; **nada deduplica entre âncoras**.

Consequência: o produto que aparece nas duas âncoras é emitido DUAS VEZES para o mesmo anúncio.
E esse é justamente o caso em que as âncoras CONCORDAM — o sinal mais forte disponível.

Vivo, agora:

| anúncio | codprod | confiança | match_input | score aplicado |
|---|---|---|---|---|
| MLB4735378521 | 22467 | 40% | seller_sku | `applyAmbiguousCorroborationScore` (:373) |
| MLB4735378521 | 44975 | 20% | ean | `applyCollisionScore` (:371) |
| MLB4735378521 | **22467** | **20%** | **ean** | `applyCollisionScore` (:371) |
| MLB6896001442 | 42519 | 40% | seller_sku | `applyAmbiguousCorroborationScore` |
| MLB6896001442 | 43534 | 20% | ean | `applyCollisionScore` |
| MLB6896001442 | **42519** | **20%** | **ean** | `applyCollisionScore` |

A fila mostra 9 linhas onde há 7 pares distintos.

Não corrompe dado: `product_links` tem `PRIMARY KEY (tenant_id, installation_id,
provider_item_id, provider_variation_id)` — o link é por ANÚNCIO. Aprovar qualquer linha liquida
as três. É defeito de **superfície de decisão**, não de persistência.

O comentário de `:347-352` declara a intenção ("todo produto do conjunto ambíguo é oferecido; a
âncora que resolveu sozinha também") e ela está certa. O que ninguém considerou foi a
INTERSEÇÃO dos dois conjuntos.

### F-2 — GTIN contraditório no mesmo par (consequência de F-1; FE ABSOLVIDO)

Na tela, duas linhas adjacentes do MESMO anúncio e do MESMO produto 22467:

```
MLB4735378521  MIST.LAV.B.ALTA POLO CROMADO  22467  GTIN —        BAIXA 40%  ✓ SKU
MLB4735378521  MIST.LAV.B.ALTA POLO CROMADO  22467  GTIN ✓ igual  BAIXA 20%  – Marca
```

Verdictos de GTIN opostos para o mesmo par. Mas `QueueRow.tsx:191-194` está **correto e
documentado**:

```tsx
// GTIN "✓ igual" is only honest when the listing matched the product on EAN
// (match_input === "ean"); otherwise the GTIN relationship is UNKNOWN → "—",
const gtinEqual = candidate.match_input === "ean" && Boolean(candidate.match_value);
```

`—` significa DESCONHECIDO nesta linha, não "sem GTIN". A regra é honesta por linha. É F-1 que
põe duas linhas de `match_input` diferente do mesmo par lado a lado e converte um UNKNOWN
honesto em contradição visível. **Consertar F-1 apaga F-2.** Nada a fazer no FE.

Mesma mecânica na coluna Motivo: `não corrobora` numa linha e `âncora ambígua` na outra, sobre
o mesmo par.

### F-3 — `side` colapsa na tela; é exatamente o que o VINC-NEUTRO conserta

Três motivos `ean INCOMPARABLE` vivos, dois `side=both` e um `side=provider`, e os três com
detail **byte-idêntico**:

```
side=both      detail=sem EAN para corroborar o CODPROD
side=both      detail=sem EAN para corroborar o CODPROD
side=provider  detail=sem EAN para corroborar o CODPROD
```

`both` = nenhum dos lados tem EAN. `provider` = o ERP tem, o anúncio não. Fatos diferentes, texto
igual. `reasonSideLabel` do CHIP-VINC-NEUTRO é o consumidor que separa os dois, e não está na
`main`. **O chip compra valor real** — só que não o valor que ele passou 7 rodadas defendendo.

### F-4 — dado do cliente, não nosso

`22467` e `44975` têm nome idêntico (`MIST.LAV.B.ALTA POLO CROMADO`) e refforn diferindo por um
ponto final: `1877.C33.` vs `1877.C33`. Vem assim do wire. Quase certamente o mesmo produto
físico duplicado no ERP. A coluna "Produto sugerido" — a coluna de decisão — não discrimina
exatamente onde discriminar é o trabalho do operador. Não é bug nosso; é achado para o operador
levar ao cliente, e é argumento a favor de exibir refforn na fila.

## O que isto diz sobre as ~18 rodadas

Os gates leram o WIRE e o RENDER. Nenhum leu o CONJUNTO de linhas de um anúncio — que é como o
operador lê a tela. F-1 e F-2 são visíveis em cinco segundos de tela e sobreviveram a todas as
rodadas; V-1 consumiu sete e não existe no dado.

O padrão não é "os gates falharam". É que a população que eles examinavam foi fixada pela
fixture, e fixture é escrita por quem já tem a hipótese. Instância direta do §11
("a população é dada pelo FATO, nunca pela pegada da edição") aplicada a validação em vez de
varredura.
