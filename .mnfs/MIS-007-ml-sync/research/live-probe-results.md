# Live ML API probe — MIS-007 round 3 (M1–M10)

Conta: instalação `inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0`, seller_id `691607102`.
Ferramenta: `apps/server_core/cmd/mlprobe` (`-round3` / `-round3-m9`), executada via
`docker exec marketplace-central-backend-1 sh -c 'cd /workspace/apps/server_core && GOCACHE=/tmp/gocache go run ./cmd/mlprobe -round3'`.
Evidência crua (uma linha por chamada, PII mascarada na captura) em
`docs/design/evidence/ml-api/M1-*.json` .. `M10-*.json` dentro do container backend.
Todas as chamadas são GET. Zero escrita no ML.

DB dev usado para descobrir IDs reais: `marketplace-central-postgres-1` / `marketplace_central`
(34 listings, todos `paused` ou `under_review` — nenhum `active`; 38 orders,
`orders_marketplace_orders` + `order_shipments`; 7 orders `cancelled`).

**CORREÇÃO (remeasure round 4):** a rodada 3 original de M3 e M5 foi capturada
com a conta em **modo férias** — por isso os 34 listings liam 0 `active` no
Postgres (estado velho) e ao vivo. O operador reativou os anúncios depois. M3
e M5 abaixo foram **remedidas ao vivo** (`-round4`), confirmando primeiro via
`GET /items?ids=` (nunca confiando no Postgres) quais IDs estão `active`
agora. As seções originais foram substituídas pelas versões remedidas; o
resto do documento (M1, M2, M4, M6–M10) é da rodada 3 e permanece válido —
não dependia de nenhum item estar `active`.

---

## M1 — offset ceiling em /orders/search + search_type=scan [BLOQUEIA M-06]

Comando: `GET /orders/search?seller=691607102&offset=N&limit=1` para N em
0, 500, 1000, 1500, 5000, 10000.

| offset | status | resultado |
|---|---|---|
| 0 | 200 | paging.total=38 |
| 500 | 200 | paging.total=38, results=[] |
| 1000 | 200 | paging.total=38, results=[] |
| 1500 | 200 | paging.total=38, results=[] |
| 5000 | 200 | paging.total=38, results=[] |
| 10000 | **400** | `{"error":"limit.maximum_exceeded","message":"Limit must be a lower or equal than 10000"}` |

O erro dispara exatamente em offset=10000 (limit=1), ou seja o teto real de
`offset+limit` para paginação normal (sem scan) em `/orders/search` é **10000**,
igual ao documentado para `/items/search`. Como esta conta só tem 38 pedidos,
o teto nunca é tocado hoje — mas para um backfill de 12 meses em conta grande
(>10000 pedidos) a paginação `offset` comum PARA nesse limite.

Segundo comando: `GET /orders/search?seller=691607102&search_type=scan&limit=1`
→ **status 200**, corpo com `display:"partial"` e `paging.scroll_id` presente
(base64, formato igual ao usado em `/users/{id}/items/search?search_type=scan`).
Resultado real veio populado (1 pedido cancelado completo, com `order_items`,
`payments`, `shipping.id`, `cancel_detail`).

Veredito: **MEDIDO**.
Decisão: **dá para fazer backfill de 12 meses acima de 1000/10000 pedidos, SIM**,
pelo caminho `search_type=scan` — não documentado para `/orders/search` na doc
pública, mas a API aceita e devolve `scroll_id` funcional. O backfill por
`offset` comum precisa parar de avançar ao cruzar `offset+limit > 10000` e
trocar para `scan` a partir daí (ou usar `scan` desde o início para contas
grandes).

---

## M2 — listing_prices: logistic_type muda sale_fee_details.fixed_fee? [BLOQUEIA M-05/M-07]

Comando: `GET /sites/MLB/listing_prices?price=P&category_id=MLB270310&listing_type_id=gold_special[&logistic_type=LT]`
para LT em `""` (sem param), `drop_off`, `xd_drop_off`, `fulfillment`, e P em 10, 13, 100.
category_id=MLB270310 é a categoria mais frequente dos 34 listings reais (23/34).

12 combinações, todas status 200. Tabela (`sale_fee_details`):

| price | sem param | drop_off | xd_drop_off | fulfillment |
|---|---|---|---|---|
| 10 | fixed_fee=4.99, gross=6.19 | fixed_fee=4.99, gross=6.19 | fixed_fee=4.99, gross=6.19 | fixed_fee=4.99, gross=6.19 |
| 13 | fixed_fee=0, gross=1.56 | fixed_fee=0, gross=1.56 | fixed_fee=0, gross=1.56 | fixed_fee=0, gross=1.56 |
| 100 | fixed_fee=0, gross=12 | fixed_fee=0, gross=12 | fixed_fee=0, gross=12 | fixed_fee=0, gross=12 |

`percentage_fee=12` em todas. Os 4 valores de `logistic_type` (incluindo
ausência do parâmetro) produziram **bit-a-bit o mesmo resultado** nos 3 preços
testados, para esta categoria/listing_type.

Veredito: **MEDIDO**.
Decisão: nesta categoria (MLB270310, gold_special) `logistic_type` **NÃO** muda
`sale_fee_details.fixed_fee` nem `sale_fee_amount`. O cálculo de fee pode
ignorar `logistic_type` como dimensão de tarifa — mas isso foi medido só numa
categoria; se M-05/M-07 dependerem de outras categorias, vale reconfirmar lá
antes de generalizar a regra para todo o catálogo.

---

## M3 — GET /item/{id}/performance (singular) [REMEDIDO round 4]

**A rodada anterior mediu contra uma conta em modo férias (0 listings active
de verdade). Esta seção substitui a original.**

Passo 0 (ao vivo, não via Postgres): `GET /items?ids=<34 ids>&attributes=id,status,sub_status,catalog_product_id`
em 2 lotes de ≤20 (limite confirmado em M8). Resultado: **9/34 listings estão
`active` AGORA** (eram 0 na rodada anterior). Dos 9 ativos, **8 têm
`catalog_product_id`** e **1 é não-catálogo**: `MLB4834219830`.

Comando: `GET /item/MLB4834219830/performance` → **status 200** (primeira vez
que o endpoint devolve sucesso nesta missão).

Estrutura real confirmada (chaves top-level):
`calculated_at, buckets[], level_wording, entity_type, entity_id, score, level`.

- `score` = 49 (float64)
- `level` = `"bad"`
- `level_wording` = `"Básica"`
- `calculated_at` = `"2026-08-01T22:39:11.162Z"`
- `entity_type` = `"USER_PRODUCT"`, `entity_id` = `"MLBU783470824"` (id do
  produto do usuário, diferente do `item_id`/`MLB...` do anúncio)
- `buckets[]`: 2 buckets reais, `USER_PRODUCT` ("Dados do produto", score
  70.59) e `ITEM` ("Condições de venda", score 18.42) — cada bucket tem
  `{key, type, title, status, score, calculated_at, variables[]}`
- `buckets[].variables[]`: cada variável tem `{key, title, status, score,
  calculated_at, rules[]}` — ex.: `UP_GTIN`, `UP_PICTURES`,
  `UP_TECHNICAL_SPECIFICATIONS_MAIN`, `UP_SHORTS`, `UP_TITLE`,
  `UP_STOCK_AVAILABILITY_TIME`, `UP_STOCK_DEPOSITO`, `UP_FREE_SHIPPING`,
  `UP_PRICE`, `UP_FINANCING`
- `buckets[].variables[].rules[]`: cada regra tem `{key, status, progress,
  mode, calculated_at, wordings}`. Exemplo real:
  `status="COMPLETED", progress=1, mode="OPPORTUNITY", wordings={label,
  link, title}` — `mode` observado sempre `"OPPORTUNITY"` nesta amostra;
  `status` variou entre `"COMPLETED"` e `"PENDING"`; `progress` variou 0 a 1
  (fracionário em pelo menos um caso: 0.33333334).

Veredito: **MEDIDO** (para o único item live-active + não-catálogo
disponível na conta — 1/1 amostra). O shape de sucesso é real, não é o
promessa da doc: as chaves batem com o pedido (`score`, `level`,
`level_wording`, `calculated_at`, `buckets[].variables[].rules[]` com
`status`/`progress`/`mode`/`wordings`), mas a amostra é n=1 porque só existe
1 listing live-active não-catálogo nesta conta agora.
Decisão: qualquer leitor do motor pode consumir esse shape com confiança
estrutural, mas — dado n=1 — vale reconfirmar em mais itens quando a conta
tiver mais ativos não-catálogo, antes de travar o parser em produção. Os
outros 8 ativos são catalog-linked e não foram testados neste endpoint (ver
M5 — todos ainda `not_listed` mesmo ativos, o que sugere que catálogo é uma
trilha separada de `/item/{id}/performance`).

---

## M4 — GET /moderations/infractions/{user_id}

Comando: `GET /moderations/infractions/691607102` → **status 200**.

Corpo: `{"infractions":[...], "paging":{"limit":20,"offset":0,"total":55}, "sorting_type":"date_created_desc"}`.
Cada infração: `{id, element_id, element_type, related_item_id, filter_subgroup, reason, remedy (HTML), site_id, user_id, date_created}`.

Os 7 listings `under_review` da base (MLB4834245368, MLB4834408368,
MLB4834419602, MLB4834395600, MLB4834206566, MLB4834232268, MLB4834408366)
**aparecem todos** na lista de 55 infrações, cada um com `reason` e `remedy`
reais, ex.:
- MLB4834206566 → reason: "A foto de capa não tem fundo branco digitalizado."
- MLB4834395600 → reason: "A foto de capa tem logotipos"
- MLB4834408366 → reason: "Não tem vendas há muito tempo." (mesmo texto que os `paused`)

Nota: `total=55` > 34 listings — a lista cobre infrações históricas mesmo de
itens não mais ativos, não é 1:1 com o status atual do listing.

Veredito: **MEDIDO**.
Decisão: os 7 `under_review` aparecem com `reason`/`remedy` reais e utilizáveis
para exibir causa/remediação na Fila de Atenção; o endpoint precisa de
paginação (`paging.total=55` > `limit=20` default) para cobrir tudo.

---

## M5 — GET /items/{id}/price_to_win?version=v2 [REMEDIDO round 4]

**A rodada anterior mediu contra uma conta em modo férias. Esta seção
substitui a original.**

DB: 10 dos 34 listings têm `catalog_product_id` não-nulo. Confirmado ao vivo
(mesmo passo 0 do M3): **8 desses 10 estão `active` AGORA** (os outros 2
continuam `paused`/`under_review` mesmo após a reativação do operador — não é
100% dos catalog-linked). Rodado `GET /items/{id}/price_to_win?version=v2`
nos 8 ativos, todos **status 200**:

| item | catalog_product_id (DB) | live status | resposta |
|---|---|---|---|
| MLB4735364085 | MLB28262741 | active | not_listed / item_not_opted_in / boosts=null |
| MLB4735304125 | MLB27562650 | active | not_listed / item_not_opted_in / boosts=null |
| MLB6896001442 | MLB23427630 | active | not_listed / item_not_opted_in / boosts=null |
| MLB6896003832 | MLB32390517 | active | not_listed / item_not_opted_in / boosts=null |
| MLB6896003262 | MLB19858075 | active | not_listed / item_not_opted_in / boosts=null |
| MLB4735328201 | MLB22624877 | active | not_listed / item_not_opted_in / boosts=null |
| MLB4735324525 | MLB35928565 | active | not_listed / item_not_opted_in / boosts=null |
| MLB4735326915 | MLB60275101 | active | not_listed / item_not_opted_in / boosts=null |

**Resultado idêntico ao da rodada anterior (conta em férias): mesmo com os 8
itens confirmadamente `active` agora, todos continuam `status=not_listed,
reason=[item_not_opted_in]`, `catalog_product_id: null` na resposta, e
`boosts: null` nos 8/8.** Nenhum trouxe `price_to_win`, `current_price`,
`winner` ou `visit_share` preenchidos — os 4 vieram `null` nos 8.

Veredito: **MEDIDO** (8/8 dos catalog-linked ativos testados).
Decisão: **reativar o anúncio não basta** — confirmado experimentalmente, não
suposição. `item_not_opted_in` persiste com o item `active`, então falta uma
etapa de opt-in/inscrição em catálogo do lado do ML que é independente do
status do anúncio (provavelmente uma ação manual ou automática de "vincular
ao catálogo" que nunca foi disparada para estes 10 itens). Isso é
**informação de produto**: `catalog_product_id` persistido no nosso banco é
só um candidato/link histórico, nunca virou participação real — e a
reativação do anúncio, sozinha, não resolve isso. `boosts[]` não pôde ser
observado preenchido (sempre `null` quando `status=not_listed`) — a estrutura
de `boosted`/`opportunity` fica **NÃO MEDIDA**: só aparece quando o item está
de fato listado no catálogo, o que não é o caso de nenhum dos 10 itens desta
conta hoje.

---

## M6 — GET /orders/{id} nos 7 pedidos cancelados: cancel_detail vs status_detail

Os 7 cancelled: 2000012659424976, 2000012747964244, 2000016865576676,
2000016891167232, 2000016959645380, 2000016971758980, 2000017258505630.

Comando por pedido: `GET /orders/{id}`. Todos **status 200**. Em todos os 7:

- `cancel_detail` **PRESENTE**, com chaves `{application_id, code, date,
  description, group, requested_by}`.
- `status_detail` = **null** nos 7 — nunca preenchido para pedido cancelado.
- `requested_by` variou: `buyer` (3x, code=`buyer_cancel_express`),
  `meli` (3x, code=`mediations`), `seller` (1x, code=`feedback_unavailable_product`).

Exemplo (order 2000012659424976):
```
"cancel_detail": {"code":"buyer_cancel_express","group":"buyer","requested_by":"buyer",
  "description":"There is a mediation with status cancel_purchase", "date":"2025-08-28T17:09:47.000-04:00", ...}
"status": "cancelled", "status_detail": null
```

Veredito: **MEDIDO** (6/7 possíveis, mas na prática 7/7 chamados e 7/7 com
o mesmo padrão).
Decisão (DEF-16): `cancel_detail` é a **única** fonte confiável de "quem
cancelou" e "por quê" — `status_detail` fica **sempre null** em pedido
cancelado nesta conta, não é um valor diferente-mas-preenchido, é
sistematicamente ausente. Qualquer leitura que dependa de `status_detail`
para motivo de cancelamento vai ler `null` sempre; o motor tem que ler
`cancel_detail.requested_by` / `.code` / `.description`.

---

## M7 — GET /shipments/{id}/costs em 5 envios reais + comparação com lead_time.cost

5 shipments delivered amostrados do banco: 47247564687 (RJ), 47255960984
(SP), 47257895403 (SP), 47259138757 (ES), 47262140888 (SP).

Comando por shipment: `GET /shipments/{id}/costs` (header `x-format-new: true`)
e `GET /shipments/{id}/lead_time` (mesmo header). Todos status 200.

| shipment | senders[0].cost | senders[0].discounts | receiver.cost | lead_time.cost |
|---|---|---|---|---|
| 47247564687 | 73.95 | promoted_amount=73.95, rate=0.5, type=mandatory | 0 | **0** |
| 47255960984 | 138.95 | promoted_amount=138.95, rate=0.5, type=mandatory | 0 | **0** |
| 47257895403 | 26.25 | promoted_amount=26.25, rate=0.5, type=mandatory | 0 | **0** |
| 47259138757 | 26.25 | promoted_amount=26.25, rate=0.5, type=mandatory | 0 | **0** |
| 47262140888 | 23.65 | promoted_amount=23.65, rate=0.5, type=mandatory | 0 | **0** |

Padrão idêntico nos 5: `senders[0]` (o vendedor, id 691607102) paga o `cost`
integral do frete, com um desconto `mandatory` de 50% já refletido em `save`;
`receiver.cost` é sempre 0 (comprador não paga); mas também tem um desconto
`ratio` em `receiver.discounts` no shipment 1 (`promoted_amount=43.9`) que não
afeta `receiver.cost` (permanece 0) — é informativo, não um custo cobrado.

Veredito: **MEDIDO** e **CONFIRMADO**: `senders[].cost` ≠ `lead_time.cost`.
`lead_time.cost` veio **0 nos 5 casos**, enquanto `senders[0].cost` variou de
23.65 a 138.95 — são campos com semânticas diferentes; `lead_time.cost` NÃO é
uma fonte válida para o custo de frete pago pelo vendedor.
Decisão: o cálculo de margem tem que ler `senders[0].cost` (ou o campo
equivalente) de `/shipments/{id}/costs`, nunca `lead_time.cost` — este último
parece refletir "custo cobrado no momento do lead_time" (0 = frete grátis
nessa etapa), não o custo real do envio.

---

## M8 — multiget /items?ids=...: limite real

Comando: `GET /items?ids=id1,id2,...` com 20, 21, 50, 51 ids (ids reais dos
34 listings, repetidos ciclicamente para completar a contagem quando
necessário).

| n solicitado | status | resultado |
|---|---|---|
| 20 | 200 | array de 20 respostas, cada uma com `code`+`body` |
| 21 | **400** | `{"error":"bad_request","message":"The parameter 'ids' only allows 20 elements. You exceeded this value by sending 21 elements"}` |
| 50 | **400** | mesma mensagem, mas reporta "sending 34 elements" (o servidor **deduplica** os ids repetidos antes de contar — só temos 34 ids distintos na conta) |
| 51 | **400** | mesma coisa |

Veredito: **MEDIDO**.
Decisão: o limite real e documentado bate: **20 ids por chamada**, erro claro
400 `bad_request` acima disso. Nenhuma ambiguidade encontrada com 21/50/51 —
a doc não se contradiz na prática, o teto é 20. Bônus não pedido mas
observado: o servidor deduplica ids antes de aplicar o limite (a mensagem de
erro em n=50/51 reporta a contagem de ids ÚNICOS, não a contagem bruta
enviada).

---

## M9 — UF de destino: shipment vs billing_info; campo de contribuinte ICMS [ALIMENTA DIFAL]

Comando por pedido (38 pedidos com `provider_shipment_id`): `GET
/shipments/{id}` (`x-format-new: true`) lendo
`destination.shipping_address.state.id` (ex.: `BR-RJ`) e `GET
/orders/{id}/billing_info` (`x-version: 2`) lendo
`buyer.billing_info.address.state.code` (ex.: `BR-RJ`).

Resultado: **38/38 pares com status 200 nos dois lados**, e **38/38 UFs
coincidentes** (nenhum mismatch). Amostra (só UF, sem PII):

| order_id | shipment UF | billing UF | match |
|---|---|---|---|
| 2000012659424976 | RJ | RJ | true |
| 2000012747964244 | SP | SP | true |
| 2000016852695592 | ES | ES | true |
| 2000016878163148 | ES | ES | true |
| ...(34 outros) | ... | ... | true |

Campo de contribuinte ICMS: em `buyer.billing_info.taxes`, chave
**`taxpayer_type.description`** aparece com valores `"Contribuinte"` /
`"Não contribuinte"` quando presente (e `{}` vazio quando não aplicável); e
`taxes.inscriptions.state_registration` traz a IE quando o comprador é PJ
contribuinte (observado em pelo menos 2 dos 38 pedidos). Também
`billing_info.attributes.cust_type` variou entre `BU` (pessoa física/buyer) e
`CO` (empresa/company) — sinal correlato mas distinto do de contribuinte.

Veredito: **MEDIDO** (38/38 pedidos).
Decisão: as duas UFs **coincidem sempre** nesta conta — não há caso
observado de divergência shipment vs billing_info; o motor de DIFAL pode usar
qualquer uma das duas fontes para a UF de destino nesta amostra (mas com
apenas 38 pedidos de uma conta a "sempre coincide" é evidência local, não
prova estrutural — ainda vale manter as duas leituras e logar divergência se
aparecer, em vez de assumir a garantia como regra da API). Existe sim um
campo de contribuinte de ICMS utilizável: `billing_info.address` state para
UF, `billing_info.taxes.taxpayer_type.description` para
contribuinte/não-contribuinte, `taxes.inscriptions.state_registration` para a
IE quando houver.

---

## M10 — GET /shipments/{id}/sla

Comando: `GET /shipments/47247564687/sla` (`x-format-new: true`) → **status 200**.

Corpo: `{"expected_date":"2026-06-09T23:59:59-03:00","last_updated":"2026-06-09T19:24:02Z","service":"","status":"on_time"}`.

Veredito: **MEDIDO**.
Decisão: o endpoint existe e responde 4 campos simples — `expected_date`
(prazo esperado), `status` (`on_time` observado; outros valores não
amostrados), `service` (vazio nesta amostra), `last_updated`. Hoje nunca é
chamado pelo sync; é uma fonte mais enxuta que `/shipments/{id}/lead_time`
para status de prazo, mas `service` vazio nesta amostra não permite confirmar
se o campo é útil — vale reconfirmar em mais de 1 shipment antes de adotar.

---

## Resumo (linha por medição)

1. **M1** — MEDIDO: offset trava em 10000 (`limit.maximum_exceeded`), mas `search_type=scan` FUNCIONA em `/orders/search` (200, scroll_id real) → backfill >10k pedidos É possível via scan.
2. **M2** — MEDIDO: `logistic_type` (drop_off/xd_drop_off/fulfillment/ausente) NÃO altera `sale_fee_details` em nenhuma das 12 combinações testadas (categoria MLB270310) → fee não depende dessa dimensão, nesta categoria.
3. **M3 [REMEDIDO]** — MEDIDO (1/1 item live-active não-catálogo): shape de sucesso real confirmado — `score=49, level="bad", level_wording="Básica", calculated_at, buckets[].variables[].rules[]` com `status/progress/mode/wordings` batendo com o pedido; `mode` sempre `"OPPORTUNITY"` na amostra. n=1 porque só 1 dos 9 listings live-active é não-catálogo.
4. **M4** — MEDIDO: os 7 `under_review` aparecem nas 55 infrações com `reason`/`remedy` reais e utilizáveis na Fila de Atenção.
5. **M5 [REMEDIDO]** — MEDIDO (8/8 catalog-linked confirmados `active` ao vivo): resultado idêntico à rodada anterior — `status=not_listed, reason=item_not_opted_in, boosts=null` nos 8/8, mesmo com o anúncio ativo → **reativar não resolve**, falta opt-in de catálogo (ação separada do status do anúncio); estrutura `boosted`/`opportunity` de `boosts[]` continua NÃO MEDIDA (nunca aparece populada nesta conta).
6. **M6** — MEDIDO (7/7): `cancel_detail` sempre presente com `requested_by`(buyer/meli/seller); `status_detail` sempre `null` em cancelado → só `cancel_detail` serve para motivo/autor do cancelamento (DEF-16).
7. **M7** — MEDIDO (5/5) e CONFIRMADO: `senders[0].cost` (23.65–138.95) ≠ `lead_time.cost` (sempre 0) → margem tem que usar `costs.senders[].cost`, nunca `lead_time.cost`.
8. **M8** — MEDIDO: limite real de `/items?ids=` é exatamente **20**; acima disso 400 `bad_request` com mensagem clara; servidor deduplica ids antes de contar.
9. **M9** — MEDIDO (38/38): UF de shipment e de billing_info **sempre coincidem** nesta conta; campo de contribuinte ICMS existe em `billing_info.taxes.taxpayer_type.description` (+ `inscriptions.state_registration` para IE).
10. **M10** — MEDIDO: `/shipments/{id}/sla` existe, responde `{expected_date, status, service, last_updated}`; amostra única, `service` vazio não confirmado como útil.
11. **BÔNUS** — MEDIDO: dos 34 listings, **9 estão `active` ao vivo agora** (eram 0 na rodada anterior — conta em férias); `sub_status[]` vem **vazio (`[]`) em 9/9** listings ativos observados.

Achado lateral da rodada 3 (agora corrigido): a rodada anterior mediu contra
**conta em modo férias** (0/34 listings `active`, tanto no Postgres quanto ao
vivo) — isso limitou M3 e produziu um M5 que não distinguia "não opt-in por
estar pausado" de "não opt-in mesmo ativo". A rodada 4 confirmou ao vivo (via
`/items?ids=`, nunca via Postgres) que 9/34 estão `active` agora, e que o
resultado de M5 **não muda** com a reativação — é uma decisão de produto
real (falta opt-in em catálogo), não um efeito colateral do modo férias.

**Correção de segurança aplicada nesta rodada:** `piiKeys` em
`apps/server_core/cmd/mlprobe/main.go` não cobria as chaves `identification`
e `state_tax_id` — por isso `docs/design/evidence/ml-api/T0-identity.json`
(rodada anterior à desta sessão) guardava o CNPJ e a inscrição estadual do
vendedor em claro. Adicionadas as duas chaves ao `piiKeys` (mascaramento
agora cobre esses campos em qualquer captura futura) e o arquivo com o dado
em claro foi apagado. Nenhum commit foi feito.
