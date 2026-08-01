# F-01-ml-ingest-readers — plan

## Files created (all new — no existing file was edited)

1. `apps/server_core/internal/modules/connectors/domain/order_detail.go`
   — `OrderDetail`, `OrderDetailItem`, `OrderDetailPayment` DTOs.
2. `apps/server_core/internal/modules/connectors/domain/shipment_detail.go`
   — `ShipmentDetail` DTO (1:1 with `order_shipments`, 0088).
3. `apps/server_core/internal/modules/connectors/adapters/mercado_livre/order_ingest_reader.go`
   — `GetOrderDetail`, `doJSONRaw` helper, `redactBillingInfo`, `mapOrderDetail`, and the
   `mlIngest*` decode structs for `/orders/{id}`.
4. `apps/server_core/internal/modules/connectors/adapters/mercado_livre/shipment_ingest_reader.go`
   — `GetShipmentDetail`, `getShipmentDetail`, `doJSONAllowGone` helper, `mapShipmentDetail`,
   `applyShipmentCosts`, and `mlIngestShipmentResponse`.
5. `apps/server_core/internal/modules/connectors/adapters/mercado_livre/order_ingest_reader_test.go`
6. `apps/server_core/internal/modules/connectors/adapters/mercado_livre/shipment_ingest_reader_test.go`

## Design decisions

- **New DTOs, not reused ones.** `OrderSnapshot`/`ShipmentInfo` stay scoped to their existing
  read-path call sites; `OrderDetail`/`ShipmentDetail` are a deliberately separate, richer
  superset for the ingest-write path (0088/0089 column-complete).
- **New decode structs (`mlIngest*`), not reused `mlOrderResponse`/`mlShipmentResponse`.**
  Avoids any risk of drifting the closed file's types and keeps the two call sites (list-read
  vs. ingest-write) independently evolvable, per the dispatch's own framing.
- **`doJSONRaw`** (new, in `order_ingest_reader.go`): the existing `doJSONWithHeaders`
  (capability_adapter.go, closed) decodes into a target but never returns the raw response
  bytes to its caller. Since a raw capture with billing_info stripped is required from the
  SAME GET (not a second round trip — that would double rate-limit spend and violate the "1
  GET" requirement), a small new helper wraps the shared `doRaw` choke point and replicates
  the exact same status-code→error mapping as `doJSONWithHeaders`, then returns both the
  decoded target and the raw bytes.
- **`doJSONAllowGone`** (new, in `shipment_ingest_reader.go`): `doJSONWithHeaders` only
  special-cases 404. IC-03/feature.md require 410 (Gone) to ALSO degrade to honest-absence for
  shipments. Rather than edit the closed file, a local wrapper adds the 410 branch and is used
  only by the shipment ingest path.
- **`redactBillingInfo`** operates on `map[string]json.RawMessage`, not
  `map[string]any`/re-marshal, specifically to avoid float64 precision loss on ML's large
  numeric order/buyer IDs — every field except `buyer.billing_info` is re-emitted
  byte-for-byte from its original `json.RawMessage`.
- **No `/carrier` or `/sla` calls.** See spec.md "Honest gaps" — neither serves a persisted
  0088 column with a confirmed shape in this worktree; calling them anyway would be
  motion without a wired outcome.
- **`x-format-new: true` header** (not the research doc's `X-Costs-New` name) — mirrors the
  existing, proven-working `getShipmentInfo`/costs call idiom already live in
  `shipping_reader.go`, rather than trusting a possibly-inaccurate research-doc header name
  against live traffic.

## Sequencing

1. Read `capability_adapter.go`, `shipping_reader.go`, `pricing_reader.go`,
   `buyer_fiscal_reader.go`, `domain/capability.go`, `domain/shipping_read.go`, IC-03, the
   ML-facts ledger, and 0088/0089 DDL — build the exact target field set before writing any
   Go.
2. Write `domain/order_detail.go` and `domain/shipment_detail.go` (DTOs first, since both
   adapter files depend on them).
3. Write `order_ingest_reader.go`, then `shipment_ingest_reader.go`.
4. `go build ./...` / `go vet ./...` — both clean on first pass after one self-caught typo fix
   (`jsonUnmarshal` → `json.Unmarshal`, missing `encoding/json` import) during authoring.
5. `gofmt -w` the 4 new files (manual editing had misaligned struct tags on 3 of them);
   re-ran `gofmt -l` clean.
6. Write both `_test.go` files against the established `httptest.NewServer` +
   `pricingTestAdapter`/`pricingAccountRef` helper pattern already used by
   `shipping_reader_test.go`/`pricing_reader_test.go`.
7. `go vet ./...` caught a real bug: `detail != (domain.OrderDetail{})` doesn't compile
   because the struct contains a `[]string` field (not comparable) — fixed the test to compare
   specific fields instead.
8. First `go test` run surfaced a real production bug (not a test bug): `pack_id` is decoded
   as `*string` but ML emits it as a bare JSON number — `json.Unmarshal` failed with
   `CONNECTORS_PROVIDER_PAYLOAD_INVALID` against a realistic fixture. Fixed by decoding
   `PackID` as `any` (same pattern already used for `order.id`/`shipping.id`) and normalizing
   via the existing `normalizeAnyID` helper in `mapOrderDetail`. This is exactly the kind of
   defect the "fixture-based decode test" requirement in the dispatch prompt exists to catch.
9. Re-ran `gofmt -l`, `go build ./...`, `go vet ./...`, `go test ./internal/modules/connectors/...`
   — all clean/green. See validation.md for verbatim output.
