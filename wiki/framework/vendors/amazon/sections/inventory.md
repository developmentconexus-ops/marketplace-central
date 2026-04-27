# Amazon SP-API: Inventory

Last verified: 2026-04-27

## What this section covers

- Inventory summaries and replenishment paths
- Inbound shipment workflows
- Supply source and multi-location inventory concerns

## Representative docs

- [FBA Inventory API](https://developer-docs.amazon.com/sp-api/docs/fba-inventory-api)
- [Get FBA Inventory Summaries](https://developer-docs.amazon.com/sp-api/docs/get-fba-inventory-summaries)
- [Create an Inbound Shipment](https://developer-docs.amazon.com/sp-api/docs/create-an-inbound-shipment)
- [Create an Inbound Order](https://developer-docs.amazon.com/sp-api/docs/create-an-inbound-order)
- [Track Inbound Shipment with SKU Details](https://developer-docs.amazon.com/sp-api/docs/track-an-inbound-shipment-with-sku-details)
- [Supply Sources API Rate Limits](https://developer-docs.amazon.com/sp-api/docs/supply-sources-api-rate-limits)
- [Multi-Location Inventory Integration Guide](https://developer-docs.amazon.com/sp-api/docs/mli-integration-guide)
- [Retrieve Inventory Eligible for Blank Box](https://developer-docs.amazon.com/sp-api/docs/retrieve-inventory-eligible-for-blank-box)

## MPC notes

- Keep inventory mutation endpoints capability-gated by fulfillment mode.
- Use reconciliation jobs for eventual consistency across inbound/stock states.

