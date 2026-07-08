# Marketplace Integration Foundation Audit

Date: 2026-07-08

## Confirmed Current State

- `AuthFlowService.HandleCallback` rotates credential and updates auth session/status, but does not atomically project the full installation connection snapshot.
- `IntegrationInstallation` is flattened and exposes credential-oriented fields.
- Mercado Livre connector read paths exist but are not instantiated in composition as runtime marketplace capabilities.
- Mercado Livre stock write code exists but is outside slice 1 runtime exposure.
- `pricing_fee_sync` must be verified as live provider-backed or demoted from runtime-live semantics.

## Slice 1 Runtime Capabilities

- `account_probe`
- `listing_read`
- `order_read`
- `fee_quote_read` only after live provider-backed implementation
- `stock_read` only after composed and live validated

## Explicitly Modeled But Not Runtime Available

- `stock_write`
- `message_reply`
- `shipment_write`
- `webhook_receive`
- `listing_write`

## Stop Conditions

- Stop if OpenAPI, SDK, backend, or UI cannot agree on one installation shape.
- Stop if implementing live fee sync requires unsafe writes or speculative provider semantics.
- Stop if a test seam is the only evidence for a live provider claim.
