# Research Note

```yaml
id: R-001
type: research
status: draft
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: support
```

## Topic

Evidence for a Mercado Livre-first marketplace cockpit that still uses a future multi-marketplace hub architecture.

## Sources Checked

- Source: Mercado Livre Developers via Context7, `/websites/developers_mercadolivre_br_pt_br`.
  Why it matters: Official source for item stock, variations, orders, and notifications.
- Source: https://github.com/mercadolibre/golang-sdk.
  Why it matters: Official Go SDK repository status affects adapter dependency decisions.
- Source: https://github.com/topics/marketplace-integration.
  Why it matters: Public marketplace integration projects describe common product scope and orchestration vocabulary.
- Source: https://github.com/timothyryner/marketplace-hub.
  Why it matters: User-cited example of small-business inventory/POS/listing/intelligence hub.
- Source: `C:\Users\leandro.theodoro\Documents\MNOS\semantic\governance\*.yml` and `semantic\views\*.sql`.
  Why it matters: Local evidence for Sankhya stock, price, cost, tax, and sales semantics.

## Findings

- Finding: Mercado Livre stock updates are item/variation operations over existing listings.
  Evidence: Context7 official docs show `PUT /items/{item_id}` with `available_quantity`, and variation updates by sending `variations` with `id` and `available_quantity`.
- Finding: Setting available quantity to zero can pause/out-of-stock an item.
  Evidence: Context7 official docs note quantity `0` changes item status/sub-state in supported cases.
- Finding: Orders contain the fields needed for initial margin quality.
  Evidence: Context7 official docs show `GET /orders/{ORDER_ID}` with `order_items`, `quantity`, `unit_price`, `sale_fee`, `payments`, `shipping`, status, tags, and cancellation details.
- Finding: Mercado Livre notifications are resource pointers, not full business payloads.
  Evidence: Context7 official docs show notification payloads containing `resource`, `user_id`, `topic`, `application_id`, `attempts`, and timestamps; the app must fetch the resource.
- Finding: The official Mercado Livre Go SDK is not a safe dependency.
  Evidence: The GitHub repository is archived and states it is no longer maintained, not functional, and recommends the official documentation.
- Finding: Marketplace hub examples converge on inventory, listing management, orders, POS/source-of-truth bridge, and intelligence.
  Evidence: `timothyryner/marketplace-hub` describes unified inventory, marketplace listings, POS, and AI intelligence; GitHub marketplace-integration topic includes API-first orchestration platforms for products, inventory, and orders.
- Finding: MNOS already contains validated internal business semantics.
  Evidence: `VW_ESTOQUE_SALDO` defines stock bucket grain and `ESTOQUE`/`RESERVADO`; `VW_PRECO_TABELA` defines current/as-of price; `TGFCUS` defines `CUSSEMICM` cost as-of; `VW_FAT_VENDA_ITEM` defines signed item sales.

## Recommendation

Build MPC as a capability-based marketplace hub. Business modules own ProductLink, StockPolicy, StockAction, MarketplaceOrder, and ProfitSnapshot. Provider adapters implement capabilities such as `ListingReader`, `StockReader`, `StockWriter`, and `OrderReader`. Mercado Livre is the first adapter, and direct HTTP is preferred over the deprecated Go SDK.

## Impact On Mission

- Architecture uses capability ports instead of Mercado Livre-specific service logic.
- M-02 exists before M-04/M-05 to prevent provider hardcoding.
- M-03 imports MNOS/Sankhya semantics as contracts before stock math or margin math.
- Webhooks are deferred because resource-pointer notifications require the same idempotent fetch/upsert paths as polling.

## Handoff

- Current status: Evidence summarized for mission planning.
- Next owner: Mission reviewer, then feature implementers.
- Next action: Use official docs again during feature execution for endpoint-specific details.
- Required files/evidence: Feature `validation.md` files must cite exact docs and tests used.
- Blockers or open decisions: None for architecture; product exclusion lists require operator examples during Stock Policy work.
