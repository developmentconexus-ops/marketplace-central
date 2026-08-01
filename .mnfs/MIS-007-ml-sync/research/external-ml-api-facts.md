# Research — Fatos externos da API ML que bindam o planning MIS-007

Fonte primária: [MIS-007-ML-SYNC-DESIGN.md §10](../../../docs/design/MIS-007-ML-SYNC-DESIGN.md)
(brainstorm ratificado 2026-07-31; doc oficial via context7 `/websites/developers_mercadolivre_br` + web).
Evidência de medição live própria (2026-07-20, conta real):
[docs/design/evidence/ml-api/](../../../docs/design/evidence/ml-api/) — dumps JSON T0–T12, F1–F5.
ATENÇÃO: dumps untracked com PII scrub pendente (pendência herdada da MIS-006) — não commitar sem scrub.

| # | Fato | Status | Evidência |
| --- | --- | --- | --- |
| 1 | `order_items[].sale_fee` é POR UNIDADE (multiplicar por quantity) | verified (medição live 2026-07-20; doc NÃO declara) | `evidence/ml-api/T2-order-sale-fee.json` |
| 2 | Multiget de shipments NÃO existe — GETs individuais paralelos | verified (live 2026-07-20) | `T5-shipments-multiget.json`, `T5b-shipments-multiget-21.json` |
| 3 | Cotação de frete por dimensão/peso sem `item_id` NÃO existe | verified (ausência na doc, pesquisa 2026-07-31) | design §4/§10 |
| 4 | `listing_prices` exige `logistic_type`/`shipping_mode`/`billable_weight` p/ bater com cobrado | verified (doc comissão-por-vender, context7 2026-07-31) | design §10 |
| 5 | Comissão dentro de `/items/{id}/sale_price` | NÃO confirmado — usar `listing_prices` | design §10 |
| 6 | Webhook: payload só `resource`+`topic`+`user_id`+`attempts`; dado se rebusca; UTC | verified (doc oficial context7 2026-07-31) | design §5 |
| 7 | Entrega: 8 tentativas/~1h; perdidas em `GET /missed_feeds` por 2 dias | verified (doc oficial context7 2026-07-31) | design §5 |
| 8 | Topic recomendado p/ vendas: `orders_v2` | verified (doc oficial) | design §5 |
| 9 | IPs de origem oficiais: 54.88.218.97, 18.215.140.160, 18.213.114.129, 18.206.34.84 | verified (doc oficial; lista mutável → uso log-only, gate P1) | design §5 |
| 10 | "200 em 500ms" como exigência dura | NÃO confirmado (exemplo oficial 498ms aceito); design independe: 200 imediato | design §5 |
| 11 | Rate limit ~1500 req/min por seller | assumed (fonte parcial); doc oficial recomenda backoff exp + jitter | design §10; `T12-429-behavior.json` |
| 12 | Orders search: cursor `date_last_updated` + `sort=date_desc` (bug com `date_asc` provado) | verified (live T7, D-120) | design §6; INTEGRATION-FINDINGS-D120 |
| 13 | Items scan: 1000/batch com `scroll_id`; multiget 20 | verified (live D-120 scan-paging + doc) | `T11-scan-timing.json`; design §6 |
| 14 | `/shipments/{id}/costs` com header `X-Costs-New`; shipment com `x-format-new` | verified (live) | `T6b-shipment-costs.json`, `F3-shipment-single-xformatnew.json` |
| 15 | SLA por `/shipments/{id}/sla` | verified (live) | `F5-shipment-sla-*.json` |
| 16 | `billing_info` por pedido | verified (live) | `T4-billing-info-*.json` |
| 17 | pack_id agrupa pedidos (carrinho) | verified (live) | `T3-pack.json` |
| 18 | Anúncio de terceiro: 403 em `/items` | verified (live) | `F2*-third-party-*.json` |
| 19 | `POST /items/catalog_listings` (optin catálogo) existe — WRITE, fora da missão | verified (doc) | design §10 |
| 20 | Predição categoria `GET /sites/MLB/domain_discovery/search?q=` | verified (doc) — MIS-008+ | design §10 |

Contagem de anúncios da conta (F1, 2026-07-20): active/paused/closed/all nos dumps
`F1-listing-count-*.json` — insumo p/ dimensionar backfill e fixtures (ler valores no P5/P6;
fixtures multi-página continuam obrigatórias independente do volume real).

Nada aqui é `verify-at-install`: tudo verificado por doc oficial datada ou medição live
própria, exceto #5, #10, #11 marcados explicitamente.
