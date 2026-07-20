# Catálogo Canônico de Consultas ML — DRAFT v0 (D-120)

Regra: **uma função canônica por necessidade de dado, sem redundância**. Toda função aqui é
A forma de buscar aquele dado; nenhum outro código chama a API ML fora deste catálogo.
Status: `DOC` = confirmado na documentação oficial (fonte citada); `LIVE` = precisa prova
live antes de ratificar. Fonte-base: developers.mercadolivre.com.br (lido 2026-07-20).

## Convenções globais (DOC)
- Multiget: `GET /{resource}?ids=a,b,c` — **máx 20 ids**; resposta = array `[{code, body}]`
  por id (200/404 individuais). Suporta `?attributes=` p/ field selection.
- Rate limit: ~1500 req/min por seller (não confirmado verbatim) → 429 corpo vazio, **sem
  Retry-After confirmado**. Cliente canônico: token bucket + backoff exponencial c/ jitter.
- OAuth: access 6h; refresh **single-use** (guardar o novo a cada refresh, atomicamente);
  refresh expira com 6 meses sem uso. Adicionar refresh-retry-once no request path (fix I8).
- Test users: `POST /users/test_user` (máx 10, 60d inatividade some, só transacionam entre si).
  Servem p/ fluxo de venda; concorrência/catálogo exigem conta real (read-only GETs).

---

## F1. Sincronização de anúncios

### F1.1 `listar_ids_anuncios(seller)` — DOC
`GET /users/{seller_id}/items/search?search_type=scan&scroll_id=...`
- Retorna **só IDs** (results = strings) — hidratação é sempre F1.2.
- Paginação normal: default 50, máx 100, teto ~1000 resultados → catálogo >1000 usa SCAN.
- **SCAN: `scroll_id` expira em 5 MINUTOS** — o loop scan→multiget não pode intercalar
  trabalho lento; coletar todos os ids primeiro, hidratar depois. Lotes de 1000/chamada,
  mesmo scroll_id reutilizado, fim = resultado vazio.
- Filtros úteis: `status`, `sku`, `seller_sku`, `orders=last_updated_desc`,
  `missing_product_identifiers`, `reputation_health_gauge` (BR).
- Aviso da própria doc: item search "não substitui notificações" — é backfill, não sync.

### F1.2 `hidratar_anuncios(ids[])` — DOC + LIVE
`GET /items?ids=<até 20>&attributes=<lista>`
- Envelope `{code, body}` por id; tratar 404 individual sem abortar lote.
- `attributes=` p/ payload enxuto no sync (id,title,status,sub_status,price,original_price,
  available_quantity,sold_quantity,category_id,condition,permalink,thumbnail,date_created,
  last_updated,listing_type_id,catalog_product_id,catalog_listing,health,tags,shipping,
  official_store_id,seller_custom_field,attributes,variations).
- **LIVE**: multiget devolve `variations[]` completas? `sold_quantity`/`tags`/`catalog_*`
  presentes? (exemplo da doc é item de teste mínimo — não prova).
- ⚠️ DOC: `price`/`base_price`/`original_price` do `/items` serão DESCONTINUADOS → preço
  canônico virá de F1.3.
- 2006 anúncios = 3 chamadas scan + 101 multigets ≈ **104 requests** (vs 2027 hoje).

### F1.3 `preco_atual(item)` — DOC
`GET /items/{id}/sale_price?context=channel_marketplace` (400 sem context — já conhecido).
Retorna amount, regular_amount, metadata.promotion_id/type. Fonte canônica de preço vigente.

### F1.4 `visitas(items)` — DOC
- Totais (2 anos): `GET /visits/items?ids=...` → `{"MLB...": 552}`.
- Janela por data: `GET /items/visits?ids=UM_SO&date_from&date_to` — **máx 1 id**, janela ≤150d.
- Série diária: `GET /items/{id}/visits/time_window?last=N&unit=day`.

---

## F2. Sincronização de pedidos

### F2.1 `buscar_pedidos(seller, desde)` — DOC (parcial) + LIVE
`GET /orders/search?seller={id}&sort=date_desc&order.date_last_updated.from=<cursor>`
- ⚠️ DOC: **sort default = `date_asc`** (mais antigo primeiro!) — bug atual dos "pedidos de
  hoje" confirmado pela doc. Sorts: date_asc/desc, updated_asc/desc, closed_asc/desc.
- Sem filtro nenhum a busca "não faz nada" — sempre passar seller + filtro.
- Offset máx 9999; >10k registros → paginação `from_id`. Retenção: **12 meses** (teto
  natural do backfill D5 ✓).
- 206 = dados parciais (tratar, não explodir).
- **LIVE**: nomes canônicos dos filtros (`order.date_last_updated.from`, `order.status`) —
  snippet não confirmou verbatim.
- NÃO existe multiget de orders — `/orders/search` + `GET /orders/{id}` é o caminho.

### F2.2 `detalhar_pedido(id)` — DOC
`GET /orders/{id}` — campos confirmados: total_amount, currency_id (root + por item + por
payment), order_items[].{item.id, variation_id, quantity, unit_price, sale_fee}, payments[],
buyer, shipping.id, tags[], context.{channel: marketplace|proximity|mp-channel, site, flows},
taxes, feedback. `pack_id` no root (LIVE confirmar shape).
- **LIVE crítico: `sale_fee` é POR UNIDADE ou total da linha?** Doc não resolve. (gross_price
  confirmado = (unit_price+discounts)×qty, ou seja, linha.) Provar com pedido real qty>1.
- Cupom/desconto NÃO é inline: `GET /orders/{id}/discounts` (endpoint próprio).

### F2.3 `pack(pack_id)` — DOC
`GET /packs/{pack_id}` → pedidos do carrinho/pacote. Shape exato LIVE.

### F2.4 `dados_nf(order)` — DOC + LIVE
`GET /orders/billing-info/{site}/{billing_info_id}` com header **`x-version: 2`**
(v1 em depreciação — nosso código atual usa qual? migrar). Shape PF (CPF) vs PJ
(CNPJ + razão social + IE) = LIVE.

### F2.5 `envio(shipment_ids[])` — DOC + LIVE
- **MULTIGET DE SHIPMENTS EXISTE** (`?ids=`, limite 20 ou 50 — doc conflita, LIVE decidir).
  Mata o 3-calls-por-pedido sequencial: enrichment em lote no ingest.
- `GET /shipments/{id}` header `x-format-new: true` → status, substatus, logistic_type,
  destination.shipping_address (endereço completo p/ NF).
- **SLA: `GET /shipments/{id}/lead_time`** (sub-recurso próprio, não vem no root):
  `estimated_handling_limit` = prazo de DESPACHO (coluna da Fila); `estimated_delivery_limit`.
- `GET /shipments/{id}/costs` header **`X-Costs-New: true`** (header próprio, diferente!):
  gross_amount, receiver.receiver_cost, senders[] (custo do vendedor), discounts[].
- `GET /shipments/{id}/carrier` → transportadora + URL rastreio.

---

## F3. Concorrência / catálogo

### F3.1 `identidade_catalogo(ean)` — DOC
`GET /products/search?site_id=MLB&status=active&product_identifier={EAN}` → id, domain_id,
attributes (BRAND/MODEL/GTIN), pictures, parent_id, children_ids. Fallback `q=` keywords.
Preferir SEMPRE `catalog_product_id` que já vem no anúncio (F1.2) — busca por EAN só p/
produto sem anúncio nosso.

### F3.2 `produto_catalogo(product_id)` — DOC
`GET /products/{id}` → status, name, domain_id, permalink, pictures, attributes,
parent_id/children_ids (ofertas SÓ em leaf — fan-out children), buy_box_winner{item_id,
seller_id, price, shipping, listing_type_id, official_store_id, seller.reputation_level_id},
buy_box_winner_price_range{min,max}. **Um único DTO** (fix M6).

### F3.3 `ofertas_concorrentes(product_id)` — DOC
`GET /products/{id}/items` (paging offset/limit, default 100) → por oferta: **item_id** (persistir!),
seller_id, price, currency, condition, listing_type_id, official_store_id, tags,
shipping{free_shipping, mode, logistic_type}. Filtros: price=, shipping_cost=free, seller_id=.
- ⚠️ **`sold_quantity` NÃO está no response documentado.** "Quantos venderam" vem de F3.5.

### F3.4 `quem_e_o_seller(seller_ids[])` — DOC
`GET /users?ids=<até 20>` (multiget!) → nickname, permalink, seller_reputation{level_id,
power_seller_status, transactions{period, total}}. Público só garante transactions.total
(completed = LIVE). Cachear forte (muda pouco).

### F3.5 `vendas_dos_concorrentes(item_ids[])` — LIVE
`GET /items?ids=<até 20>&attributes=id,sold_quantity,price,available_quantity,seller_id`
sobre os item_id de F3.3. **LIVE**: sold_quantity vem em item de OUTRO seller? (público no
site; provável, mas provar). ⚠️ available_quantity público é BUCKETIZADO (faixas: 1, 50,
100... — nunca tratar como exato).

### F3.6 `nossa_posicao(item)` — DOC
`GET /items/{id}/price_to_win?version=v2` → price_to_win, status (winning/competing/
sharing_first_place/listed), reason[] (enum documentado: reputation_below_threshold,
item_paused...), boosts[] (fulfillment/free_shipping/... × boosted/opportunity),
visit_share, competitors_sharing_first_place, winner{item_id, price}. Não lista concorrentes
— lista é F3.3.

---

## F4. Tarifas

### F4.1 `comissao(price, category, listing_type)` — DOC
`GET /sites/MLB/listing_prices?price=&category_id=&listing_type_id=[&logistic_type=]`
→ sale_fee_amount + sale_fee_details{percentage_fee, meli_percentage_fee, fixed_fee,
financing_add_on_fee, gross_amount}.
- ⚠️ **Regra 2026 da taxa fixa: depende do TIPO LOGÍSTICO** (Flex cobra abaixo do TH;
  ME1/custom sempre abaixo do TH; ≥TH ninguém cobra). Valor do TH não publicado na doc
  (FAQ) — **LIVE por faixa de preço**. R$79 hardcode continua banido; agora nem o threshold
  é fixo por site, é por logística.
- Passar `logistic_type` na consulta quando conhecido (temos via F1.2 shipping).

### F4.2 `frete_gratis_custo(seller, item)` — DOC (atual)
`GET /users/{seller}/shipping_options/free?item_id=` (mantém; sem batch documentado).

---

## F5. Webhooks (fase 2 — pós URL pública; scheduler-first ratificado)
- Tópicos: `orders_v2`, `items`, `shipments`, `price_suggestion`, `catalog_item_competition`.
- Config no painel da aplicação (Minhas Aplicações → Tópicos + callback URL HTTPS), não por API.
- Callback deve responder **200 em ≤500ms** → ack imediato + fila assíncrona, nunca processar inline.
- Retry ~1h / 8 tentativas; falha repetida DESATIVA o tópico. Recuperação: `GET /missed_feeds`
  (guarda 2 dias; período desativado é irrecuperável).
- Payload: {resource:"/orders/123", user_id, topic, attempts, sent, received}.

---

## Rodada LIVE consolidada (bloqueia ratificação de DTO)
| # | Prova | Função |
|---|-------|--------|
| T1 | multiget /items devolve variations/tags/catalog_*/sold_quantity completos | F1.2 |
| T2 | sale_fee por unidade vs linha (pedido qty>1) | F2.2 |
| T3 | pack_id no root do pedido + shape /packs/{id} | F2.2/F2.3 |
| T4 | billing-info x-version:2 — shape PF vs PJ (IE p/ PJ?) | F2.4 |
| T5 | multiget shipments: existe? limite 20 ou 50? campos c/ x-format-new | F2.5 |
| T6 | lead_time: estimated_handling_limit presente em pedido real | F2.5 |
| T7 | filtros canônicos /orders/search (date_last_updated.from etc.) + sort=date_desc | F2.1 |
| T8 | sold_quantity de item de terceiro via multiget | F3.5 |
| T9 | transactions.completed público em /users/{id} | F3.4 |
| T10 | listing_prices: varrer faixa de preço p/ achar TH por logistic_type | F4.1 |
| T11 | scan: scroll na nossa conta 2006 itens, medir tempo total | F1.1 |
| T12 | 429: comportamento real (header? corpo?) num burst controlado | global |

Tudo read-only GET na conta real conectada; test users só p/ fluxo de pedido sintético.
Evidência: respostas cruas salvas em `docs/design/evidence/ml-api/` (1 arquivo por prova).

---

## VEREDICTOS LIVE — executado 2026-07-20 (rounds T0–T12 + follow-ups F1–F5)

Conta: METALNOBREACABAMENTOS (seller 691607102, PJ/CNPJ). Ferramenta: `apps/server_core/cmd/mlprobe`
(roda no container backend, token via credencial ativa do banco, nunca exposto). Evidência crua
(token redigido, PII mascarado) em `docs/design/evidence/ml-api/`.

| # | Veredicto | Impacto no design |
|---|-----------|-------------------|
| T1 | ✅ multiget `/items?ids=` devolve TUDO dos NOSSOS itens: sold_quantity, date_created, catalog_product_id, sub_status, tags, variations, health, permalink, shipping | F1.2 ratificada; hidratação de anúncio = 1 call/20 itens |
| T2 | ✅ **sale_fee é POR UNIDADE**: qty=2, unit_price=729.90, sale_fee=120.43 = 16,5% de 729.90 | Comissão total = `sale_fee × quantity`. Decomposição usa isso |
| T2b | `/orders/{id}/discounts` existe; 404 `discount_not_found` quando pedido sem cupom | 404 = "sem desconto", não erro |
| T3 | ✅ pack_id no root do pedido; `/packs/{id}` → {orders[], shipment.id, status} | Carrinho: agrupar pedidos por pack, 1 envio por pack |
| T4 | ✅ billing_info c/ `x-version:2` OK: identification{type CPF/CNPJ, number}, endereço completo, taxes.inscriptions (IE viria aqui p/ PJ; amostras eram CPF) | Dados NF ratificados. PJ real ainda não amostrado |
| T5 | ❌ **multiget /shipments?ids= NÃO EXISTE** (404 com 2 e com 21 ids) | F2.5 revista: GET individual por shipment, paralelizado (goroutines + budget), persistido no ingest |
| T6 | `/lead_time` existe mas SEM estimated_handling_limit; tem estimated_delivery_time{handling:24h, shipping:72h}, pay_before, list_cost | prazo de entrega ok, prazo de despacho NÃO vem daqui |
| F5 | ✅ **`/shipments/{id}/sla` = fonte do SLA da Fila**: {expected_date, status: on_time/delayed} | Coluna SLA da Fila = expected_date + status deste endpoint, persistido |
| F3 | ✅ `/shipments/{id}` (x-format-new) devolve tudo: destination+endereço, tracking, logistic, priority_class, lead_time embutido | 1 call cobre envio+endereço+tracking; /sla à parte |
| T6b | ✅ `/costs` (X-Costs-New) OK: senders[].cost=19.85 c/ desconto 50% mandatory, gross 63 | frete_vendedor real p/ decomposição |
| T7 | ✅ `date_last_updated.from` + `sort=date_desc` funcionam (paging.total=35 últimos 30d) | Cursor incremental F2.1 ratificado |
| T8/F2 | ❌ **/items de TERCEIRO = 403 access_denied em TODAS as variantes** (multiget, single, com/sem auth — PolicyAgent) | sold_quantity/detalhe de concorrente INVIÁVEL via /items. Lane concorrente = `/products/{id}/items` (T8a ✅: price, seller_id, item_id, listing_type, shipping, official_store — SEM sold_quantity) + perfil público do seller |
| T9 | ✅ `/users/{seller_id}` público: nickname, reputation{level_id, power_seller_status, transactions.total=416}, user_type, cidade/UF | "quem é" e proxy de volume do concorrente. transactions.total = TODAS transações históricas, não vendas do item |
| T10/F4 | ✅ tarifa: percentage_fee 11,5% (cat. teste, gold_special). **fixed_fee = 50% do preço para itens ≤ R$12,50; ZERO acima** (13–200 testado). Regra "R$79/R$6" NÃO observada nesta categoria | ml_tariffs guarda percentage+fixed por categoria/listing_type/faixa; limiar por categoria via sweep, não constante |
| T11/F1 | Conta tem **34 anúncios** (9 active, 18 paused, 0 closed): scan 2 batches/799ms. Os 2006 eram do CLIENTE PROSPECT (outra conta, catálogo xlsx #004-E) — não desta conta ML | Backfill desta conta é trivial; arquitetura dimensionada p/ contas grandes continua (scan+multiget) |
| T12 | 100 GETs concorrentes (10 workers) → 100×200, **zero 429** | Limite bem acima disso; backoff defensivo mantido, Retry-After segue não provado |

### Pendências pós-live
- PJ billing_info real (amostrar pedido com comprador CNPJ) — shape de IE.
- Limiar fixed_fee em OUTRAS categorias (sweep por categoria no sync de tarifas).
- 429 real nunca observado — instrumentar log quando acontecer em produção.
- `sale_fee` documentado por unidade AQUI (gold_pro, catálogo); revalidar em pedido gold_special.
