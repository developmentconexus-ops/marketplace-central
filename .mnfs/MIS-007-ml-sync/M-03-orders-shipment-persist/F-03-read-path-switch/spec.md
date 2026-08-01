# F-03-read-path-switch — spec

## Scope (per feature.md brief)

Replacement-before-deletion (ADR-05): `GET /orders/{id}`'s enrich path stops calling the
live Mercado Livre shipment/buyer-fiscal readers and reads Postgres rows written by F-02's
`IngestOrder` instead. Sites A (shipment) and B (fiscal) are retired from composition wiring
in the same diff that adds their Postgres replacements. The archguard interactive-ML
allowlist (M-02 F-04) shrinks to match. `comprador_fiscal` stays byte-identical in shape;
only its data source changes. `GetOrderBucketCounts` reads the persisted `bucket` column
with re-derivation fallback for legacy rows. `sumSaleFee` gets a bugfix (missing
`×Quantity`) discovered while re-reading the enrich path.

## Design decisions

1. **`ports.ErrShipmentNotFound` sentinel** — a shipment row genuinely absent (order has no
   `order_shipments` row, e.g. never shipped or pre-F-02) is honest-unknown, not an error to
   warn about. `EnrichService.fetchShipment` checks `errors.Is(err, ports.ErrShipmentNotFound)`
   first (silent degrade to nil) before the generic error path (which still warns once) — this
   preserves the existing warn-on-real-failure behavior for actual DB/query errors while
   removing false-positive warnings for the common "not shipped yet" case.

2. **Money `::text` cast in SQL** — `order_shipments.cost_gross`/`cost_seller` are
   `numeric(14,2)`. Scanning numeric columns into Go float64 and reformatting risks producing
   `"4.25e+07"`-style scientific notation, which fails
   `connectorsdomain.ValidateMoney`'s decimal regex
   (`^[+-]?(?:\d+(?:\.\d+)?|\.\d+)$`, plus uppercase-currency requirement). Casting to
   `::text` in the SELECT (`cost_gross::text`, `cost_seller::text`) hands back the driver's
   canonical fixed-point representation verbatim (e.g. `"42.50"`), which always satisfies the
   regex. Verified directly in `shipment_reader_test.go`'s
   `TestScanShipmentRow_MoneyAmountSatisfiesValidateMoneyRegex`.

3. **`StateName` always nil (BuyerFiscalReader)** — migration 0089 has no `uf_nome`/
   `state_name` column (only `buyer_address_state_code`). VERIFIED safe, not merely assumed:
   `apps/web/src/pages/pedidos/PedidoDrawer.tsx`'s `formatEndereco` (around line 341-347)
   reads only `uf_codigo`. This is a genuine shape difference from the OLD live adapter
   (`connectors/adapters/mercado_livre/buyer_fiscal_reader.go`'s `mapBuyerFiscalAddress` DOES
   populate `StateName` from the ML payload's `state.Name`), so it is explicitly excluded from
   both sides of the golden-shape comparison (see `TestBuildBuyerFiscalInfo_StateNameAlwaysNil`
   and the golden tests' doc comments) rather than silently treated as "no difference."

4. **Accepted, documented, out-of-scope UX gap** — `ShipmentEnrichment.CarrierName`,
   `TrackingURL`, and `Costs.ReceiverCost` become permanently nil after the switch: migration
   0088 has no backing columns for them. This regresses `PedidoDrawer.tsx`'s rastreio/
   frete_real display fields. Per the brief, this is captured here as an accepted gap, not
   fixed — re-adding a live call to backfill these fields would reintroduce the exact live-read
   dependency this feature removes. See `TestScanShipmentRow_AlwaysNilFieldsPer0088Gap`.

5. **`GetOrderBucketCounts` two-tier read** — persisted `bucket` column (non-null, non-empty)
   wins; empty-string or SQL-NULL falls back to `DeriveOrderBucket` re-derivation with an
   empty shipment status string, byte-identical to the pre-F-03 behavior. Empty-string and
   NULL are treated identically because `UpsertOrders` (unlike `IngestOrder`) writes
   `string(order.Bucket)` unconditionally — a legacy/never-re-ingested row's Go zero value
   `""` round-trips as an empty non-NULL string, not SQL NULL.

6. **`sumSaleFee` bugfix** — `domain.MarketplaceOrderItem.SaleFeeAmount` is documented as
   PER-UNIT (verbatim from the provider's `OrderDetailItem.SaleFeeUnit`, never
   pre-multiplied). The pre-existing `sumSaleFee` summed `SaleFeeAmount` alone, under-counting
   commission (and therefore `retorno_liquido`/`margem_pct`) for every multi-quantity line.
   Fixed to `sum += *item.SaleFeeAmount * float64(item.Quantity)`.

## Ownership boundaries respected

Owned: `orders/application/enrich_service.go`, `orders/adapters/postgres/` (new reader files
+ `order_repo.go`'s `GetOrderBucketCounts`), `orders/transport/http_handler_test.go`
(additive golden test only — `http_handler.go` itself untouched, no DTO/OpenAPI change),
`composition/root.go` orders region, `composition/orders_adapters.go` (site A deletion),
`platform/archguard/archguard_test.go` (+ its `testdata/three_sites/root.go` fixture, a
necessary minimal consequence of the allowlist shrink). `composition/orders_adapters_test.go`
also required a matching edit: it held a now-dangling test (`TestOrdersShipmentReaderAdapter`)
against the deleted `newOrdersShipmentReaderAdapter` constructor — removed along with its
single-use `erroringInstallationRepo` fixture, same necessary-consequence reasoning as the
archguard testdata fixture.

Forbidden and untouched: `modules/connectors/`, `modules/listings/`, `modules/pricing/`,
F-02's `ingest_service.go`/`ingest_repo_test.go`/`order_ingest_errors.go`/`order_shipment.go`.

## Not done (explicitly out of scope per brief §8)

No `pack_id` or other new fields added to `enrichedOrderDTO`/OpenAPI/SDK. `comprador_fiscal`
DTO shape is unchanged — verified via the golden tests, not merely asserted.
