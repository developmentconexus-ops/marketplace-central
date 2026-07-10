# F-01 Specification - Order Ingestion

```yaml
id: M-06-F-01
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-06
created: 2026-07-09
updated: 2026-07-09
```

## Goal

Create the first `orders` slice as a real business module that ingests Mercado Livre orders idempotently, persists normalized tenant-scoped snapshots, and records link-quality truth without coupling business state to `integrations` or provider payloads.

## Architectural Position

- `integrations` keeps installation health, credentials, runtime capability gating, and provider operation audit.
- `connectors` keeps Mercado Livre API calls, payload mapping, and provider error semantics.
- `orders` owns normalized order state, ingestion idempotency, link-quality projection, and future order lifecycle/business projections.
- `profitability` is deferred; this feature must not calculate margin or infer unknown money values as zero.

## Scope

This feature includes:

- a new `orders` module under `apps/server_core/internal/modules/orders`
- normalized `MarketplaceOrder` and `MarketplaceOrderItem` snapshots
- payment, shipping, tags, and cancellation snapshot fields required by the brief
- tenant-scoped Postgres persistence
- idempotent re-ingestion by provider order id plus provider update timestamp
- ingestion application service that reads through provider capability ports
- explicit link-quality state per order item based on existing `product_links` truth

This feature does not include:

- profitability or margin calculation
- freight/cost/tax enrichment
- shipment-detail reads beyond `shipping.id`
- buyer PII beyond what is strictly required for operational references
- operator UI

## Source Contracts

### Provider Input

Provider reads continue to come from the already wired Mercado Livre order capability path:

- `integrations/application/provider_operation_service.go`
- `connectors/ports/marketplace_capability.go`
- `connectors/adapters/mercado_livre/capability_adapter.go`

The new `orders` module may depend on an orders-owned source port implemented by an adapter that delegates to `ProviderOperationService`, but it must not import Mercado Livre payload types.

### Product Link Truth

Link resolution truth comes from `product_links` state:

- `resolved` and `rejected` come from persisted product links
- `conflict` and unresolved candidate states come from workflow/candidate state
- missing link is quality state, not ingestion blocker

The `orders` module should consume an adapter-facing port that exposes only the listing identity and resolved internal product identity/state needed for ingestion quality.

## Normalized Domain Shape

### Order Identity

Each order is uniquely identified by:

- `tenant_id`
- `installation_id`
- `provider_order_id`

### Idempotency Rule

Reprocessing the same provider order:

- inserts once when absent
- updates the existing row when the provider snapshot changes
- never creates duplicate order rows
- never creates duplicate item rows
- never creates duplicate payment rows

Provider freshness must be stored explicitly:

- `provider_updated_at` from `last_updated` when present
- fallback to `date_last_updated` when that is the authoritative provider field
- `fetched_at` for local observation time

The upsert must only replace the current snapshot when the incoming provider update timestamp is newer or equal to the stored snapshot. Unknown provider update timestamp is allowed but must not silently erase a known newer timestamp.

### Order Fields To Persist

Persist at minimum:

- installation/provider identity
- `provider_order_id`
- `provider_status`
- `provider_status_detail`
- `provider_created_at`
- `provider_closed_at`
- `provider_updated_at`
- `shipping_id`
- `cancellation_detail`
- `tags`
- `raw_provider_ref` limited to safe provider references, not raw full payload dumps
- `fetched_at`
- `created_at`
- `updated_at`

### Item Fields To Persist

Persist per item:

- `provider_item_id`
- `provider_variation_id`
- `seller_sku`
- `title`
- `quantity`
- `unit_price`
- `sale_fee_amount`
- item-level link quality state
- resolved `internal_product_id` when available

### Payment Fields To Persist

Persist per payment:

- `provider_payment_id`
- `provider_status`
- `transaction_amount`
- `total_paid_amount`

## Link Quality Model

Each ingested order item must carry explicit link quality:

- `resolved`
- `rejected`
- `conflict`
- `unresolved`
- `missing`

Rules:

- `resolved` means an unambiguous internal product link exists for the listing identity
- `rejected` means operators explicitly rejected the listing link
- `conflict` means candidate/workflow state is ambiguous
- `unresolved` means candidate/workflow exists but is not resolved
- `missing` means no persisted link/candidate truth exists yet

Link quality never blocks order ingestion.

## Provider Mapping Rules

Based on official Mercado Livre docs, the ingestion path may safely normalize:

- `GET /orders/{ORDER_ID}`: `id`, `date_created`, `date_closed`, `last_updated`, `status`, `status_detail`, `order_items[].item.id`, `order_items[].item.title`, `order_items[].item.variation_id`, `order_items[].item.seller_sku`, `quantity`, `unit_price`, `sale_fee`, `payments[].id`, `payments[].status`, `payments[].transaction_amount`, `payments[].total_paid_amount`, `shipping.id`, `tags`
- `GET /orders/search`: required `seller`, optional `order.status`, `order.date_last_updated.from`, `order.date_last_updated.to`, with `results[]` usable to enumerate order ids and summary timestamps

Source: Mercado Livre Developers docs retrieved on 2026-07-09 via Context7 from `/websites/developers_mercadolivre_br_pt_br`, pages "gerenciamento-de-vendas", "gestao-packs", and "pedidos-e-opinioes".

## Storage Design

Add first-class order tables instead of JSON blobs as canonical business state.

Planned tables:

- `orders_marketplace_orders`
- `orders_marketplace_order_items`
- `orders_marketplace_order_payments`

All tables must:

- carry `tenant_id`
- scope by `installation_id`
- use forward-only migrations
- support deterministic upsert/update

Tags may be stored as JSON/text array. Cancellation detail may start as a normalized string summary in F-01 if the provider shape is inconsistent, but the schema must leave room for richer structured fields later.

## Application Flow

Recommended F-01 flow:

1. `orders` application requests provider order snapshots for an installation and cursor/filter.
2. Provider source adapter delegates to `integrations` provider operation service.
3. `orders` normalizes provider snapshots into domain aggregates.
4. `orders` resolves per-item link quality through a module-owned link reader port.
5. `orders` upserts orders, then replaces or upserts items/payments inside one transaction.
6. Result returns imported/updated counts and the high-level cursor context.

## API Surface For This Slice

The minimal operational surface for F-01 should prefer a manual import trigger over an operator UI.

Recommended first surface:

- `POST /orders/import`
- `GET /orders?installation_id=...`

If the first implementation needs to stay even smaller, the trigger endpoint may land first and the list endpoint can be added in the same slice only if it helps validation materially.

## Validation Requirements

Tests must cover:

- paid order ingestion
- canceled order ingestion
- order with variation id
- order with missing product link
- idempotent reprocessing without duplicate rows

Runtime evidence for this feature may use real Mercado Livre order reads as source proof, but must clearly separate:

- provider read validation
- orders-module persistence validation

## Non-Negotiables

- no buyer PII leakage into contracts or persistence beyond necessary operational references
- no margin or freight math in this slice
- no unknown-to-zero fallback
- no direct `connectors` or provider payload types crossing into transport/UI contracts
- no cross-module SQL into `product_links`; use ports/adapters
