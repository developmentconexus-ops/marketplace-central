# Amazon SP-API: Orders and Fulfillment

Last verified: 2026-04-27

## What this section covers

- Orders retrieval and filtering
- Fulfillment status transitions
- Shipping label and shipment workflows
- PII access paths tied to orders

## Representative docs

- [Orders API](https://developer-docs.amazon.com/sp-api/docs/orders-api)
- [Orders API Rate Limits](https://developer-docs.amazon.com/sp-api/docs/orders-api-rate-limits)
- [Get Orders with Filtering Criteria](https://developer-docs.amazon.com/sp-api/docs/get-orders-with-filtering-criteria)
- [Fulfill Orders](https://developer-docs.amazon.com/sp-api/docs/fulfill-orders)
- [Track a Partially Fulfilled Order](https://developer-docs.amazon.com/sp-api/docs/track-partially-fulfilled-order)
- [Merchant Fulfillment API](https://developer-docs.amazon.com/sp-api/docs/merchant-fulfillment-api)
- [Retrieve Shipment Labels](https://developer-docs.amazon.com/sp-api/docs/retrieve-shipment-labels)
- [Get authorization to access PII for order items](https://developer-docs.amazon.com/sp-api/docs/get-authorization-to-access-pii-for-order-items)

## MPC notes

- Keep explicit status maps for order and package state machines.
- Isolate restricted data calls behind RDT-aware connector paths.

