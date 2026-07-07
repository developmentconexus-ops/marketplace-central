# Module: Inventory

Layer: operational control
Path: `apps/server_core/internal/modules/inventory/` (planned)

## Main Question It Answers

"Is Mercado Livre announcing stock that is safe relative to our internal stock truth?"

## What This Module Owns

| Entity | Purpose |
|---|---|
| `StockSnapshot` | Internal and Mercado Livre stock values with timestamps |
| `StockPolicy` | Safety buffer, max advertised quantity, and write rules |
| `StockAction` | Proposed or applied provider stock change with audit trail |

## Rules

- Internal stock comes from Sankhya/MetalShopping through ports, not from connector state.
- Mercado Livre writes require a resolved product link, current provider state, safety policy, and audit record.
- Unknown stock, stale stock, ambiguous link, or provider read failure must block automatic writes.
- `available_quantity` updates must be idempotent and record before/after values.
- `internal_read` returning `missing_stock` means stock quantity stays nil and the workflow is blocked rather than defaulted to `0`.

## Initial Mercado Livre Scope

- Compare internal available stock to Mercado Livre announced quantity.
- Show divergence and risk: oversell, undersell, stale, unresolved, healthy.
- Start with assisted/manual apply; automatic write is a later policy decision.
