# Schema de Armazenamento Alvo — DRAFT v0 (D-120)

Princípio: **estender o que existe via migração, não recriar**. Schema atual está ~certo
(listings, orders_*, product_links, market_*); buracos = colunas faltantes + tabelas de
sync/espelho. Tudo `tenant_id`-scoped. Joins canônicos no fim.

## Visão (entidades e joins)

```
products_mirror (ERP espelho) ──┐
                                ├─ product_links ──┐
listings (anúncios ML) ─────────┘                  │
    │  └─ listing_variations                       │
    │                                              │
    ├── market: competitor_offers ─ sellers_cache  │
    │          market_aggregates                   │
    │                                              │
orders ── order_items ─────────────────────────────┘ (item.id+variation → link → codprod)
    ├── order_payments
    └── order_shipments (novo: enrichment persistido)
ml_tariffs (histórico comissão)      sync_state (cursores por conta/entidade)
```

## Tabelas

### `products_mirror` — NOVA (decisão D7 revista: SIM)
Espelho canônico de produto p/ JOIN SQL, alimentado por qualquer adapter (xlsx grava do
snapshot; Sankhya sync grava do read-through). `erp_import_products` continua como histórico
por protocolo (auditoria); mirror = estado corrente.
```
tenant_id, source (sankhya|xlsx), codigo_produto (PK c/ tenant), descricao, referencia,
ean, marca, grupo_codigo, grupo_descricao, custo NUMERIC NULL, preco_venda NUMERIC NULL,
estoque_total NUMERIC NULL, updated_at, protocol_id (origem)
+ products_mirror_stock_locations(tenant_id, codigo_produto, local_codigo, local_descricao, quantidade)
```
NULL = honesto-desconhecido (nunca 0 fake). Ingest = **upsert-merge keep-absent** (ADR-031 /
F-XLSX-1): linha do arquivo/query novo faz UPSERT por `codigo_produto`; linha ausente do snapshot
novo NÃO é deletada — marca `absent_in_last_snapshot`/`stale_since` (preserva `product_links`).
Nunca wipe, nunca delete físico. Coluna: `stale_since timestamptz NULL`.

### `listings` — ESTENDER (existente, migração 0036)
Adicionar colunas (contrato E3): `sub_status text[]`, `sold_quantity`, `initial_quantity`,
`category_id`, `condition`, `permalink`, `thumbnail_url`, `date_created_ml` (**data real do
ML; `fetched_at` continua sendo nosso**), `start_time`, `stop_time`, `catalog_product_id`,
`catalog_listing bool`, `health`, `tags text[]`, `official_store_id`, `shipping_mode`,
`shipping_free bool`, `logistic_type`, `original_price`, `brand`, `model`,
`commission_amount NUMERIC NULL`, `commission_pct NUMERIC NULL` (F4.1 no ingest),
`free_shipping_cost NUMERIC NULL` (F4.2 no ingest), `visits_30d INT NULL` (F1.4 diário).

### `listing_variations` — NOVA
Hoje variação vive na PK de listings (variation_id) com campos do pai. Formalizar:
```
tenant_id, installation_id, provider_listing_id, variation_id (PK), price,
available_quantity, sold_quantity, seller_sku, attribute_combinations JSONB, picture_ids text[]
```

### `orders_marketplace_orders` — ESTENDER (0027 + 0074)
Adicionar: `pack_id BIGINT NULL`, `channel` (context.channel), `total_amount`, `paid_amount`,
`currency_id`, `taxes_amount NUMERIC NULL`, `coupon_amount NUMERIC NULL` (F2.2 /discounts),
`decomposicao JSONB NULL` + `retorno_liquido NUMERIC NULL` + `margem_pct NUMERIC NULL`
(**decomposição PERSISTIDA no ingest** — motor deixa de ser nil), `raw_tags text[]` (expor).
`order_items`: + `currency_id`, `sale_fee_total` (semântica ratificada pós-T2), `full_unit_price`.

### `order_shipments` — NOVA (enrichment persistido; mata as 3 calls no read)
```
tenant_id, order_id, shipment_id (PK), status, substatus, logistic_type,
receiver_address JSONB (rua/nº/cidade/UF/CEP — NF), sla_expected_date timestamptz,
sla_status text (on_time|delayed — live F5), estimated_delivery_date timestamptz,
tracking_number, tracking_url, carrier_name, cost_seller NUMERIC NULL,
cost_receiver NUMERIC NULL, gross_amount NUMERIC NULL, fetched_at
```
LIVE T5: multiget /shipments NÃO existe → alimentada por GETs individuais paralelizados
(`/shipments/{id}` x-format-new + `/shipments/{id}/sla` + `/costs`) no ingest do pedido.
Bucket da Fila derivado daqui + faturado_at, persistido em `orders.bucket` (indexável).

### `competitor_offers` — SUBSTITUI `market_validated_offers` (estender)
Por oferta/coleta: + `competitor_item_id` (**hoje dropado — chave do rastreio**),
`listing_type_id`, `official_store_id`, `shipping_free`, `logistic_type`, `collected_at`.
LIVE T8/F2: /items de terceiro = 403 (todas variantes) → **sold_quantity de concorrente
não existe via API**; fonte única = `/products/{id}/items` (price/seller/shipping/listing_type).
"Quantos venderam" = honesto-desconhecido; proxy = transactions.total do seller (T9).

### `sellers_cache` — NOVA
```
seller_id (PK), nickname, permalink, reputation_level, power_seller_status,
transactions_total (histórico TOTAL do seller, NÃO vendas do item — T9), user_type,
city, state, fetched_at (TTL ~7d)
```
Sem tenant (dado público global). Sem multiget — GET individual cacheado.

### `market_aggregates` — MANTER (0053/0070)
+ `collected_at` exposto na UI (fix M4). `buy_box_winner` unificado (fix M6):
+ colunas winner_item_id, winner_seller_id, winner_price, winner_listing_type.

### `ml_tariffs` — NOVA (design ml-tariff-design-pending, ratificar junto)
```
site, category_id, listing_type_id, logistic_type, price_from, price_to,
percentage_fee, fixed_fee, financing_add_on_fee, valid_from, valid_to NULL,
source (VENDA|COTACAO), fetched_at
```
Histórico: nunca UPDATE, sempre nova linha + fecha valid_to. Mata seed 16%/22%.
LIVE T10/F4: percentage 11,5% na categoria testada; fixed_fee = 50% do preço p/ itens
≤ R$12,50, zero acima — limiar é POR categoria (sweep no sync), não constante global.

### `sync_state` — NOVA (coração do motor de sync)
```
tenant_id, installation_id, entity (orders|listings|market|tariffs|products),
cursor JSONB (ex.: {date_last_updated_from}, {scroll_checkpoint}), last_full_sync_at,
last_incremental_at, last_error JSONB NULL, consecutive_failures INT
```
Um scheduler lê daqui: pedidos ~5min (F2.1 incremental), anúncios diário (F1.1+F1.2),
mercado diário (F3.*, fila com rate budget), tarifas semanal (F4.1).

### Mantidas sem mudança estrutural
`product_links` + candidates + audit (só + trigger automático de geração),
`erp_import_protocols/products` (histórico), `integration_operation_runs` (+ LIMIT na query — L5),
credenciais/installations.

## Joins canônicos (as perguntas do negócio)
```sql
-- anúncio → produto → custo (margem por anúncio)
listings l JOIN product_links pl USING (tenant_id, installation_id, provider_listing_id, variation_id)
           JOIN products_mirror pm ON pm.codigo_produto = pl.internal_product_id
-- pedido → item → produto (custo congelado + margem real)
order_items oi JOIN product_links pl ON pl.provider_item_id = oi.item_id ...
               JOIN products_mirror pm ...
-- pedido → envio (Fila/SLA): orders o JOIN order_shipments os USING (order_id)
-- anúncio → concorrência: l JOIN market_aggregates ma ON ma.codprod = pl.internal_product_id
-- oportunidades: products_mirror pm JOIN market_aggregates ma LEFT JOIN product_links (excluir vinculados — fix M1)
```

## Ordem de migração proposta
1. `sync_state` + scheduler esqueleto (nada muda p/ usuário, tudo passa a ser rastreável)
2. `products_mirror` (+ fix source por tenant)
3. `listings` estendida + `listing_variations` (junto do novo backfill scan+multiget)
4. `order_shipments` + orders estendida (junto do ingest incremental + decomposição)
5. `competitor_offers`/`sellers_cache`/aggregates (junto da coleta em lote)
6. `ml_tariffs`
