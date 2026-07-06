# Module: Orders

Layer: operational monitoring
Path: `apps/server_core/internal/modules/orders/` (planned)

## Main Question It Answers

"What happened with each Mercado Livre order, and what operational action is needed?"

## What This Module Owns

| Entity | Purpose |
|---|---|
| `MarketplaceOrder` | Provider order header and lifecycle state |
| `MarketplaceOrderItem` | Provider item/variation quantity, price, fee, and internal link reference |
| `ShipmentState` | Shipment reference/status needed for dispatch guardrails |
| `CancellationDetail` | Provider cancellation reason, requester, and date |

## Rules

- Orders read from Mercado Livre via connectors; internal product/cost resolution happens through ports.
- Provider raw payloads may be snapshotted for audit, but domain logic uses typed fields.
- Missing product links must not block order ingestion; they produce data-quality flags for downstream modules.
- Order ingestion must be idempotent by provider order ID and update timestamp.

## Initial Mercado Livre Scope

- Import paid/canceled orders, items, payments, sale fees, shipment IDs, status, tags, and cancellation details.
- Feed profitability and inventory exception views.
