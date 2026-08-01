# F-01-ml-ingest-readers — spec

Feature 1 of 3 in the serial DAG (F-01→F-02→F-03) for MIS-007/M-03-orders-shipment-persist.

## Goal

Add two new reader methods on `*mercadolivre.CapabilityAdapter`, in new files, so F-02
(orders module, out of scope here) can persist a full order + its shipment without touching
`capability_adapter.go` (CLOSED, M-01-owned) or the existing narrower `getShipmentInfo`
reader (shipping_reader.go, a different, still-LIVE call site).

```go
func (a *CapabilityAdapter) GetOrderDetail(ctx context.Context, accountRef domain.ProviderAccountRef, providerOrderID string) (domain.OrderDetail, error)
func (a *CapabilityAdapter) GetShipmentDetail(ctx context.Context, accountRef domain.ProviderAccountRef, shipmentID string) (domain.ShipmentDetail, error)
```

## Ground truth consulted

- `migrations/0088_order_shipments.sql`, `migrations/0089_orders_marketplace_orders_sync_fields.sql`
  — exact DB column set both DTOs must be able to fill.
- `.mnfs/MIS-007-ml-sync/research/orders-persistence-interface-contract.md` (IC-03) — binding
  field/source contract.
- `.mnfs/MIS-007-ml-sync/research/external-ml-api-facts.md` — verified/unverified fact ledger
  (facts #14 sale_fee-per-unit, #15 /shipments/{id}/sla exists but shape unconfirmed, #17
  pack_id groups a cart).
- `internal/modules/orders/application/import_service.go` `normalizeOrders` (108-162) — baseline
  field set already relied on elsewhere in the codebase.
- `docs/design/handoff/API-MAP.md` — doc-reviewed confirmation that `tracking_number` is a
  direct top-level field of `GET /shipments/{id}`.

## Constraints honored

- ADR-17: unknown provider facts are nil pointers, never fabricated zero/blank.
- ADR-03: raw `/orders/{id}` payload held in memory only, `buyer.billing_info` stripped;
  never persisted (F-02's job to not persist it — this DTO carries it in-memory only).
  Shipment/costs payload is never captured raw at all — `order_shipments` (0088) has no raw
  column.
- `sale_fee` is per-unit (fact T2/#14) — DTO field is `SaleFeeUnit`, never `SaleFeeTotal`.
- 404/410 on shipment → honest-absence (`domain.ShipmentDetail{}, nil` with `Found=false`).
- 403 on order (third-party order, a documented live fact) → typed non-retryable error via the
  same status-code mapping as 401.
- No `shipping.id` in the order payload → `OrderDetail.ShippingID` stays nil (caller must not
  call `GetShipmentDetail`).
- `capability_adapter.go` and `shipping_reader.go` were read for reusable private
  types/helpers but never edited.
- `orders` module and `root.go` (F-02/F-03 scope) were not touched.

## Honest gaps (documented, not blocking)

Left nil because their JSON key/endpoint shape could not be confirmed against any
fixture/doc present in this worktree (see `domain/shipment_detail.go` doc comments for the
per-field reasoning):

- `ShipmentDetail.LogisticType`
- `ShipmentDetail.TrackingMethod`
- `ShipmentDetail.SLAStatus` (the `/shipments/{id}/sla` sub-resource is not called at all — its
  shape is unconfirmed and calling it with no field mapping would spend a rate-limit token for
  zero benefit)

`ShipmentDetail.SLALimitAt` is NOT a gap: it is populated from the primary
`GET /shipments/{id}` payload's `lead_time.estimated_delivery_limit.date` — the same field the
existing, live `mapShipmentInfo` (shipping_reader.go) already decodes. An earlier pass of this
feature incorrectly filed it under "unconfirmed `/sla` endpoint"; corrected after adversarial
review (see validation.md "Correction round").

The `/carrier` sub-resource (mentioned in the original dispatch prompt's summary of "three
calls: primary, costs, carrier") is deliberately NOT called: `order_shipments` (0088) has no
`carrier_name`/`carrier_id` column, so there is nothing to persist from it. This is a
concrete-DDL-over-prompt-paraphrase decision.
