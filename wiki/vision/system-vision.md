# System Vision: Marketplace Central

## What MPC Is

Marketplace Central is an intelligence and control surface for sellers operating in multiple marketplaces through VTEX.

MPC centralizes:
- pricing strategy by channel
- integration auth and health status
- fee schedule visibility for simulation
- operational traceability (sync/auth runs)

## Scope of This Wiki

This wiki focuses on architecture and operations for Marketplace + Integrations.

Planning note:
- Task sequencing and progress tracking are maintained in Nexus Brain.
- This wiki is not the execution roadmap.

## Product Direction

- Keep business policy in Marketplaces module.
- Keep credential/auth lifecycle in Integrations module.
- Keep coupling explicit and minimal.
- Keep frontend data-driven from backend definitions.

## MetalShopping Compatibility

The current structure remains merge-friendly:
- modular boundaries match `server_core/internal/modules/*` expectations
- platform helpers remain reusable
- database design stays tenant-ready and forward-only
