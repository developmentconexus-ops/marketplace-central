# Mercado Livre Operating Cockpit Research

Date: 2026-07-06

## Context

ADR-005 pivots Marketplace Central away from VTEX and toward a Mercado Livre-first internal cockpit backed by Sankhya/MetalShopping data.

## Mercado Livre Capability Priority

| Priority | Capability | Useful API concepts | Practical value | Main risk |
|---|---|---|---|---|
| 1 | Items, stock, and price | Items/listings, variations, `available_quantity`, listing types/exposures | Align announced stock and price with internal truth | Variation rules and provider-side ownership can change write semantics |
| 2 | Orders and packs | Order read model, order items, payments, status, cancellation details | Feed stock risk, fulfillment, and profitability | Orders, packs, payments, and shipments are separate concepts |
| 3 | Shipments and labels | Shipment read model, labels, shipment items, tracking | Operational dispatch visibility | Logistics mode changes behavior (`me1`, `me2`, fulfillment, drop-off, custom) |
| 4 | Notifications | Event topics and resource re-fetch pattern | Reduce polling delay and stale stock/order views | Requires deduplication, retries, and idempotent processing |
| 5 | Questions and post-sale messages | Questions API, post-sale messages | Centralize seller response workflow | Buyer data and message scope are restricted by context |
| 6 | Pricing/fees strategy | Listing types, exposures, listing prices/net proceeds where available | Support margin-aware pricing decisions | Fee data is category/site/listing-type dependent and not a single universal source |

## Business Opportunity Priority

| Priority | Opportunity | Data needed | Success metric | Recommended phase |
|---|---|---|---|---|
| 1 | Net margin radar by SKU | cost, price, taxes, ML fees, freight, historical sales | fewer negative-margin SKUs, higher visible margin coverage | Orders + Margin ML |
| 2 | Margin-band repricing | current price, cost, stock, fee, sales history | improved margin while preserving volume | Pricing + Strategy ML |
| 3 | Rupture and coverage alerts | available stock, sales velocity, lead time, open orders | fewer stockouts and canceled orders | Stock Seguro ML |
| 4 | Stock aging with action suggestions | stock age, turnover, cost, margin | lower aged inventory and capital tied up | Pricing + Strategy ML |
| 5 | Cost/price/tax divergence detection | cost history, price table, taxes, orders | fewer bad-margin surprises | Orders + Margin ML |
| 6 | Kit/bundle candidates | component stock/cost, sales history, margin | higher ticket and long-tail sell-through | Pricing + Strategy ML |

## Recommended Mission Sequence

1. Stock Seguro ML: product links, stock divergence, assisted stock actions.
2. Orders + Margin ML: order ingestion, fees/freight/cost reconciliation, margin quality flags.
3. Shipping + Notifications: shipment visibility and event-driven refresh.
4. Pricing + Strategy ML: margin bands, aging stock, promotions, kits.
5. Questions + Messages: seller response cockpit after core operations are stable.

## Implementation Implications

- Product links are prerequisite state for stock writes, order margin, and repricing.
- Provider writes must be audited, idempotent, and guarded by safety policies.
- Unknown cost, freight, fees, taxes, or links must be explicit data-quality states.
- Frontend must render recommendations and actions from backend application services; no margin or stock safety math in React.

## Source Notes

Research used official Mercado Livre/Mercado Libre developer documentation via Context7 and delegated review:

- Mercado Livre Developers `/websites/developers_mercadolivre_br_pt_br`
- Global Selling Mercado Libre devsite references for items, orders, shipments, notifications, messages, listing types, and pricing
