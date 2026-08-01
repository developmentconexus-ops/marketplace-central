# Mapa de capacidade da API do Mercado Livre — macro e micro

```yaml
id: MIS-007-ML-CAPABILITY-MAP
type: research
status: draft
owner: hub
parent: MIS-007
created: 2026-08-01
scope: site MLB, conta de vendedor, OAuth de vendedor
```

## 0. Como ler

Fontes, em ordem de autoridade:

1. **Payload real** capturado da conta do operador (`docs/design/evidence/ml-api/`). Ganha sempre.
2. **Doc oficial** (developers.mercadolivre.com.br), acessada via context7.
3. **Inferência** — marcada como tal, nunca disfarçada de fato.

Onde 1 e 2 discordam, a divergência está registrada explicitamente. Onde nenhuma das duas
responde, está escrito `NÃO VERIFICADO` — não foi preenchido com suposição.

Nota de método: as ferramentas MCP do context7 não carregaram dentro dos subagentes; o
corpus foi consultado pelo CLI `ctx7` (mesma base). Registrado para que ninguém leia
"context7" como algo que não foi feito.

---

## 1. Macro — o que a API é

O Mercado Livre não tem "a API do anúncio" nem "a API do pedido". Tem **um recurso raso e
muitos sub-recursos caros**. Essa é a característica que deve governar todo o desenho:

- `GET /items` traz muito num multiget barato, mas **qualidade, visitas, competição de
  catálogo e infrações são endpoints separados, um por item**.
- `GET /orders/{id}` traz o pedido inteiro **menos o que importa para margem**: o custo real
  de frete está em `/shipments/{id}/costs`, e **não existe multiget de shipments** (404
  medido, e a doc não documenta nenhum). Cada pedido custa no mínimo 2 chamadas, 4 se
  quisermos SLA e fiscal.
- **Notificação nunca carrega dado.** Todo webhook manda só `{resource, topic, user_id, …}`
  e obriga um `GET`. Webhook reduz latência, não reduz chamadas.

Consequência de desenho: **o custo de sincronização é dominado por sub-recurso 1-a-1**, não
pela listagem. Qualquer plano que trate "buscar pedidos" como uma chamada está errado por
uma ordem de grandeza.

### 1.1 Limites que restringem o desenho

| limite | valor | fonte |
|---|---|---|
| multiget de itens | 20 **ou** 50 por chamada — **a doc se contradiz entre duas páginas** | doc |
| multiget de shipments | **não existe** (404 com 2 e com 21 ids) | dump + ausência na doc |
| visitas por janela | **1 item por chamada**, janela ≤ 150 dias | doc |
| varredura > 1000 anúncios | `search_type=scan`, `scroll_id` expira em **5 min**, não combina com offset | doc + dump |
| varredura > 1000 pedidos | **`NÃO VERIFICADO`** — scan/scroll não é documentado para `/orders/search` | — |
| webhook | responder **200 em ≤ 500 ms**; 8 tentativas em 1 h; `missed_feeds` guarda só 2 dias; falha repetida **desativa o tópico** | doc |
| rate limit | tabela genérica (100–200 req/min); **100 chamadas seguidas deram 0× 429** no teste real | doc fraca + dump |

---

## 2. Micro — anúncio

### 2.1 O que vem barato, no multiget de `/items`

`id · title · price · base_price · original_price · currency_id · listing_type_id ·
available_quantity · initial_quantity · sold_quantity · status · sub_status[] · tags[] ·
category_id · domain_id · catalog_listing · catalog_product_id · user_product_id ·
inventory_id · attributes[] (inclui SELLER_SKU e GTIN) · sale_terms[] · variations[] ·
permalink · thumbnail · pictures[] · shipping{mode, logistic_type, free_shipping,
local_pick_up, store_pick_up, dimensions, methods[], tags[]} · seller_address · deal_ids[]`

Três leituras que mudam o produto:

- **`sold_quantity` é vitalício, não janela.** Usar como "vendas do mês" é erro de semântica.
- **`initial_quantity` não é `available + sold`** — reposição incrementa disponível sem
  tocar o inicial. Não são dois campos do mesmo número.
- **`health` vem `null`** em 10/10 itens reais. Ver 2.2.

### 2.2 Qualidade — estávamos lendo o campo errado

A doc oficial declara: *`/health` será descontinuada e substituída pela API `/performance`*.

- **`GET /item/{id}/performance`** (singular `item`) → `score`, `level`, `level_wording`,
  `calculated_at`, `buckets[].variables[].rules[]` com `status`, `progress`, `mode`
  (`OPPORTUNITY`) e `wordings.{title,label,link}` — o link aponta direto para a ação de
  correção.
- O `health` embutido no item está `null` porque está sendo desligado, não porque a conta
  não tem qualidade medida.

**Isso reverte o veredito anterior de "`quality_score` sem fonte".** Tem fonte; é outro
endpoint, custa 1 chamada por item, e entrega muito mais do que um número: entrega a lista
de o-que-fazer.

### 2.3 Competição de catálogo — o dado que faltava para preço

`GET /items/{id}/price_to_win?version=v2` → `current_price`, `price_to_win`, `status`
(`winning`/`losing`/`listed`), `visit_share`, `competitors_sharing_first_place`,
`winner.{item_id, price, boosts[]}`, e **`boosts[]`** (`fulfillment`, `free_shipping`,
`free_installments`, `shipping_collect`, `same_day_shipping`) cada um com
`opportunity`/`boosted`.

`boosts[]` é a resposta para "por que estou perdendo mesmo com preço igual". Nenhuma tela
nossa mostra isso hoje.

`GET /products/{id}/items` complementa com as ofertas concorrentes (preço, frete, seller) —
e é o **único** caminho legítimo: item de terceiro dá **403 sempre**, com auth, sem auth,
com `attributes=` restrito. Medido em 5 variantes.

### 2.4 Por que um anúncio está parado

- `sub_status[]` no próprio item: `paused_by_seller`, `moderation_penalty`,
  `waiting_for_patch`, `forbidden`, `held`, `pending_documentation`, `suspended`,
  `suspended_for_prevention`, `picture_downloading_pending`, `poor_quality_thumbnail`,
  `deleted`.
- `GET /moderations/infractions/{user_id}` → `reason` e `remedy` em texto — a explicação e a
  correção. Também existe `/items/{id}/infractions`.
- Usuário suspenso: `GET /users/{id}` → `status.list.allow == false`.

### 2.5 Tarifa

`GET /sites/MLB/listing_prices` aceita `category_id`, `price`, `currency_id`,
`listing_type_id`, **`logistic_type`**, `shipping_mode`, `billable_weight`, `tags`,
`quantity`. Devolve `sale_fee_amount`, `sale_fee_details{fixed_fee, gross_amount,
percentage_fee, financing_add_on_fee}`, `listing_fee_amount`, `listing_fee_details{…}`.

**A doc chama `logistic_type` de crucial para o `fixed_fee` sair correto.** Simulação que
não passa `logistic_type` devolve tarifa errada — e isso atinge M-05 e M-07 diretamente.

Medido nos dumps (categoria MLB420116, `gold_special`): `percentage_fee` fixo em 11,5 %;
`fixed_fee` cai de 6,24 para **0 entre R$ 12,50 e R$ 13,00**. `listing_fee_amount` = 0.

**Não existe** endpoint de "tarifa deste anúncio publicado" — `listing_prices` é paramétrico.
Tarifa por anúncio é composição nossa, não leitura direta.

### 2.6 Estoque Full

`inventory_id` no item aponta o estoque Full; o saldo se lê em
`GET /user-products/{id}/stock` → `locations[]{type, quantity}` com
`selling_address` / `meli_facility` / `seller_warehouse`.

**Esta conta não usa Full** (`inventory_id` null em 10/10). Não é prioridade, mas o desenho
não pode assumir que estoque mora só no item.

---

## 3. Micro — pedido, envio, dinheiro

### 3.1 O pedido

`GET /orders/{id}`: `status`, `status_detail`, `cancel_detail`, `total_amount`,
`paid_amount`, `currency_id`, `taxes`, `tags[]`, `static_tags[]`, `fulfilled`, `pack_id`,
`buying_mode`, `context{channel, site, flows}`, `feedback`, `mediations[]`, `order_items[]`,
`payments[]`, `shipping{id}`.

- **`status` tem vocabulário fechado documentado**: `confirmed, payment_required,
  payment_in_process, partially_paid, paid, partially_refunded, pending_cancel, cancelled,
  invalid`.
- **`status_detail` NÃO tem vocabulário publicado.** `NÃO VERIFICADO`.
- **`cancel_detail` é um objeto** — `group`, `code`, `description`, `requested_by`
  (buyer/seller/Mercado Livre), `date`. É a fonte real do motivo de cancelamento.
- `tags[]` observado: `pack_order`, `catalog`, `paid`, `not_delivered`. Vocabulário completo
  não publicado.

### 3.2 Comissão — a armadilha

| campo | onde | valor no pedido real medido |
|---|---|---|
| `order_items[].sale_fee` | por item | **120,43** |
| `payments[].marketplace_fee` | por pagamento | **0** |

No mesmo pedido. **`marketplace_fee` veio zero havendo comissão real cobrada.** Quem calcular
margem por `marketplace_fee` calcula margem inflada.

E `sale_fee` é **por unidade**: no pedido medido, `quantity=2`, `unit_price=729,90`,
`sale_fee=120,43` → 16,5 % do preço unitário. Se fosse total da linha daria 8,25 %, fora de
qualquer faixa de comissão do ML. A doc **não** diz qual das duas é — todos os exemplos
oficiais usam `quantity=1`, o que torna a doc estruturalmente incapaz de desambiguar.

**Custo de comissão da linha = `sale_fee × quantity`.** Fonte: medição, não documentação.

`gross_price = (unit_price + discounts.full) × quantity` — fórmula documentada e confirmada
no payload real.

### 3.3 Frete — a segunda armadilha

| campo | significa |
|---|---|
| `lead_time.cost` / `cost_type: "free"` | custo **do comprador**. Zero aqui = frete grátis para ele |
| `/shipments/{id}/costs` → `senders[].cost` | **custo do vendedor**, já líquido de descontos |
| `/shipments/{id}/costs` → `receiver.cost` | o que o comprador paga |

No envio real medido: comprador pagou **0**, vendedor pagou **19,85** (com desconto
obrigatório de 50 %). Ler `lead_time.cost=0` como "frete não custou nada" erra a margem em
R$ 19,85 por pedido.

### 3.4 Envio

`GET /shipments/{id}` (com `x-format-new`): `status`, `substatus`, `logistic{mode, type,
direction}`, `tracking_number`, `tracking_method`, `lead_time{…}`, `destination`, `origin`,
`dimensions`, `declared_value`, `tags`, `date_created`, `last_updated`.

- 8 status primários: `to_be_agreed, pending, handling, ready_to_ship, shipped, delivered,
  not_delivered, cancelled`. **Mais de 100 substatus**, e a doc **não** publica uma tabela
  única — está fragmentada por página de produto. O substatus real desta conta
  (`invoice_pending`) **não aparece em nenhuma lista oficial consultada**.
- Sub-recursos: `/sla` (`status` ∈ `on_time, delayed, early, insuficient_info`;
  `expected_date`), `/costs`, `/lead_time`, `/items`, `/history` (transições
  `{status, substatus, date}` — a linha do tempo real do envio), `/invoice_data` (nota
  fiscal: `invoice_number`, `invoice_amount`, `status`, `fiscal_key`).

**Conclusão operacional: vocabulário de substatus é aberto.** Qualquer CHECK, enum ou
`switch` exaustivo sobre substatus vai quebrar em produção. Confirma IC-07/ADR-06 —
verbatim, sempre.

### 3.5 Sincronização incremental

`GET /orders/search?seller=…&order.date_last_updated.from=<cursor>&sort=date_asc` —
confirmado por doc **e** por chamada real. É o mecanismo correto.

Ressalvas medidas:
- `order.date_created` tem granularidade só até a **hora**; não confirmado se vale para
  `date_last_updated`.
- **Offset máximo e suporte a scan/scroll em `/orders/search` = `NÃO VERIFICADO`.** O
  scan/scroll é documentado para `/items/search` e `/questions/search`, **não** para orders.
  Backfill de 12 meses acima de ~1000 pedidos é risco aberto.

### 3.6 Fiscal

- `GET /orders/{id}/billing_info` (header `x-version: 2`) — funciona na conta. Dados fiscais
  do comprador. **A doc documenta outra rota** (`/orders/billing-info/{site}/{id}`) — as duas
  podem coexistir; qual é a canônica é `NÃO VERIFICADO`.
- `invoice_type` foi **removido** das respostas; mapear tipo de nota virou responsabilidade
  do integrador.
- `GET /shipments/{id}/invoice_data` — nota fiscal associada ao envio.
- `GET /users/{id}/invoices/{id}` — XML/estruturado da nota.
- Relatórios de faturamento do vendedor: `/billing/…/monthly/periods` + `/documents` +
  `/summary` + `/details`. **A própria doc grafa o prefixo de dois jeitos diferentes** —
  validar ao vivo antes de implementar.

### 3.7 Pós-venda — domínio inteiro ausente do nosso produto

`GET /post-purchase/v1/claims/search` e `/claims/{id}` (+ `/detail`, `/actions-history`,
`/status-history`, `/affects-reputation`, `/reasons/{id}`).

`status` (`opened`/`closed`) · `type` (`mediations, return, fulfillment, ml_case,
cancel_sale, cancel_purchase, change, service`) · `stage` (`claim, dispute, recontact,
stale, none`) · `reason_id` (`PNR` não recebido, `PDD` defeituoso, `CS` cancelado) ·
`players[].available_actions` (`send_message, refund, open_dispute,
add_shipping_evidence, allow_return, allow_partial_refund, send_tracking_number`) ·
`affects_reputation`.

Reclamação aberta é a coisa mais urgente que existe num pedido e **nossa tela não tem nem
o conceito**.

### 3.8 Webhooks — lista real

`orders_v2` · `orders_feedback` · `payments` · `shipments` · `flex-handshakes` ·
`fbm_stock_operations` · `post_purchase` (`claims`, `claims_actions`) · `messages`
(`created`, `read`) · `items` · `questions` · `stock-location` · `items_prices` ·
`price_suggestion` · `catalog_item_competition_status` · `catalog_suggestions` ·
`user_products_families` · `public_offers` · `public_candidates` · `invoices` ·
`leads_credits` · `vis_leads`.

Payload: `{_id, resource, user_id, topic, application_id, attempts, sent, received}`; com
subtópico, `actions[]`.

**Correções a fazer no brief do M-08:** a política real é **8 tentativas em janela de 1 h**
(o brief diz 5), o receptor precisa responder **200 em ≤ 500 ms**, `missed_feeds` só guarda
**2 dias**, e falha repetida **desativa o tópico** para a aplicação — isso é perda
silenciosa de sincronismo, e merece alarme próprio.

IPs de origem documentados: `54.88.218.97`, `18.215.140.160`, `18.213.114.129`,
`18.206.34.84`, `35.236.253.169`, `35.245.91.34`, `35.245.20.104`, `35.186.182.146`.
Vieram de resumo processado, **não** de leitura literal da tabela — confirmar antes de virar
allowlist de firewall.

---

## 4. Divergências entre doc e realidade (todas medidas)

| # | doc diz | conta real diz | quem ganha |
|---|---|---|---|
| 1 | `health` populado 0–1 | `null` em 10/10 | real — campo em descontinuação |
| 2 | multiget de itens até 20 / até 50 (duas páginas) | não testado além de 10 | doc se contradiz |
| 3 | `marketplace_fee` ≈ comissão | `marketplace_fee = 0` com `sale_fee = 120,43` | real |
| 4 | `sale_fee` sem unidade declarada | 16,5 % do preço unitário ⇒ por unidade | real |
| 5 | order tem objeto `coupon` | ausente do payload | real |
| 6 | substatus listados por página de produto | `invoice_pending` não consta em nenhuma | real — vocabulário aberto |
| 7 | `billing_info` em `/orders/billing-info/{site}/{id}` | `/orders/{id}/billing_info` responde 200 | as duas coexistem; canônica indefinida |
| 8 | item de terceiro sem restrição declarada | **403 em 5 variantes** | real |
| 9 | soma por status = total | `9 + 18 + 0 = 27 ≠ 34` | real — 7 em `under_review` |
| 10 | rate limit 100–200 req/min | 100 chamadas seguidas, **zero 429** | limite real não alcançado |

---

## 5. O que a API **não** dá

- **Vendas por janela de tempo por anúncio.** Não existe endpoint. `sold_quantity` é
  vitalício. Só derivando de `/orders/search` por item e data. → `sales_30d` é **derivado**,
  nunca lido.
- **Multiget de shipments.** Cada envio custa uma chamada; com `/costs` e `/sla`, três.
- **Tarifa de um anúncio publicado.** Só a função paramétrica `listing_prices`.
- **Item de terceiro por id.** Só via catálogo.
- **Vocabulário fechado** de: `status_detail` de pedido, `tags`/`static_tags`,
  `payments[].status`, substatus de envio, `logistic.mode/type/direction`, `cust_type`.
  Todos abertos ⇒ **verbatim obrigatório, enum proibido**.

---

## 6. Lacunas que ficam abertas (nenhuma foi preenchida com suposição)

1. Limite real do multiget de itens (20 ou 50).
2. `/orders/search` suporta scan/scroll? Offset máximo? — **bloqueia o backfill de 12 meses
   do M-06 se a conta passar de ~1000 pedidos.**
3. Escopo OAuth exigido por `billing_info`; restrição de retenção/LGPD explícita.
4. Rota canônica de `billing_info`.
5. Prefixo correto dos relatórios de faturamento.
6. Vocabulário de `payments[].status` no domínio ML (distinto do Mercado Pago).
7. API de devoluções (Returns) — existe, não foi detalhada.
8. Custo de Full para o vendedor — conta não usa, sem payload real.
9. Limiar real de 429 e header `Retry-After`.
10. `settings` completo de `GET /categories/{id}`.
