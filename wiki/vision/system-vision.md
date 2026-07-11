# System Vision: Marketplace Central

## What MPC Is

Marketplace Central is an internal operations and intelligence cockpit for selling through Mercado Livre, backed by Sankhya/MetalShopping data.

MPC centralizes:
- product and listing linkage between internal products and Mercado Livre announcements
- stock reconciliation between internal stock and Mercado Livre available quantity
- pricing and profitability strategy using cost, fees, freight, and tax inputs
- order monitoring, cancellation analysis, and operational traceability
- integration auth, health status, and sync/action audit

## Scope of This Wiki

This wiki focuses on architecture and operations for Marketplace + Integrations.

Planning note:
- Task sequencing and progress tracking are maintained in `.mnfs/`.
- This wiki is not the execution roadmap.

## Product Direction

- Mercado Livre first; other marketplaces remain deferred catalog entries until the Mercado Livre operating loop is reliable.
- Keep business policy in Marketplaces and future operational modules.
- Keep credential/auth lifecycle in Integrations.
- Keep Sankhya/MetalShopping reads behind explicit ports and adapters.
- Keep frontend data-driven from backend definitions and SDK methods.
- Treat unknown cost, freight, fees, taxes, and product links as explicit data-quality states, never optimistic defaults.

## MetalShopping Compatibility

The current structure remains merge-friendly:
- modular boundaries match `server_core/internal/modules/*` expectations
- platform helpers remain reusable
- database design stays tenant-ready and forward-only
