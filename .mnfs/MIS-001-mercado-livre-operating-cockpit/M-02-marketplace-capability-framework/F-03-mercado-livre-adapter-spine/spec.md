# F-03 Spec

```yaml
id: F-03
type: feature-spec
status: in_progress
owner: Feature Implementer
parent: M-02
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Goal

Provide a Mercado Livre capability adapter spine in `connectors/adapters/mercado_livre` that implements the normalized marketplace capability ports with direct HTTP seams and no dependency on the archived official Go SDK.

## Scope

- Implement `ListingReader`, `StockReader`, `StockWriter`, and `OrderReader` for Mercado Livre.
- Read listing/item shapes from `GET /items/{ITEM_ID}` and seller item enumeration via `GET /users/{USER_ID}/items/search`.
- Read orders from `GET /orders/{ORDER_ID}` and seller order enumeration via `GET /orders/search?seller=...`.
- Update stock via `PUT /items/{ITEM_ID}` for item-level or variation-level quantity writes.
- Map rate limit, validation, unsupported-shape, transient, and invalid-reference cases to connector capability errors/results.

## Non-Goals

- No credential repository wiring, background sync orchestration, or transport endpoints in this feature.
- No live provider write validation without operator-controlled credentials/listings.
- No stock policy or business reconciliation rules in the adapter.

## Design

- Adapter accepts an HTTP client, base URL, optional site ID, clock seam, and access-token resolver function.
- `ListListings` and `ListOrders` enumerate IDs/search results, then hydrate normalized snapshots through `ReadListing` and `ReadOrder`.
- `ReadStock` uses item shape inspection to distinguish item-level versus variation-level stock.
- `UpdateAvailableQuantity` pre-reads the item so unsupported variation shapes are rejected before write attempts.

## Error Contract

- `429` -> `CONNECTORS_PROVIDER_RATE_LIMITED`
- `5xx` / network failures / decode failures -> `CONNECTORS_PROVIDER_TRANSIENT` or `CONNECTORS_PROVIDER_PAYLOAD_INVALID`
- missing provider refs / empty account context -> `CONNECTORS_PROVIDER_INVALID_REFERENCE`
- ambiguous item-vs-variation stock writes -> `CONNECTORS_PROVIDER_UNSUPPORTED_SHAPE`
- provider write rejection (`4xx` except `429`) -> `StockWriteResult{Result: rejected}`

## Validation Plan

- Unit tests cover:
  - item with no variation
  - item with variation
  - order with `sale_fee`
  - provider rejection
  - `429` mapping
  - unsupported variation shape
