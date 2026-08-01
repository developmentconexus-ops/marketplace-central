# F-01-ml-ingest-readers — validation

All commands run from `apps/server_core` with `GOCACHE`/`GOMODCACHE` exported as absolute
paths, each as a separate bash invocation (per harness process rule), on Windows/Git Bash.

## go build ./...

```
$ export GOCACHE="$(pwd)/.gocache" GOMODCACHE="$(pwd)/.gomodcache" && go build ./...
(no output — clean)
```

## go vet ./...

```
$ export GOCACHE="$(pwd)/.gocache" GOMODCACHE="$(pwd)/.gomodcache" && go vet ./...
(no output — clean)
```

First `go vet` pass caught one real defect before `go test` even ran: `order_ingest_reader_test.go`
compared `detail != (domain.OrderDetail{})`, which does not compile because `OrderDetail`
contains a `[]string` field (`Tags`) and Go structs with slice fields are not comparable with
`!=`. Fixed by comparing specific fields (`ProviderOrderID`, `Items`, `RawOrder`) instead.
Second `go vet ./...` was clean.

## go test ./internal/modules/connectors/...

First run surfaced a real production bug via the fixture-based decode test (exactly the
scenario the dispatch prompt's "VALIDATION EXPECTED" section asked this test to catch):

```
--- FAIL: TestGetOrderDetailDecodesAllTargetFields (0.00s)
    order_ingest_reader_test.go:85: GetOrderDetail() error =
    CONNECTORS_PROVIDER_PAYLOAD_INVALID: GET /orders/2000003508372839 -> HTTP 200: { "id":
    2000003508372839, ... "pack_id": 2000003508372800, ... }
```

Root cause: `mlIngestOrderResponse.PackID` was declared `*string`, but the fixture (built to
match ML's documented numeric-id convention, the same as `order.id`/`shipping.id`) sends
`pack_id` as a bare JSON number. `json.Unmarshal` cannot decode a JSON number into `*string`.
Fixed by decoding `PackID any` and normalizing through the existing `normalizeAnyID` helper in
`mapOrderDetail` (`apps/server_core/internal/modules/connectors/adapters/mercado_livre/order_ingest_reader.go`),
matching the pattern already used for every other provider id in the file.

After the fix, full clean run:

```
$ export GOCACHE="$(pwd)/.gocache" GOMODCACHE="$(pwd)/.gomodcache" && go test ./internal/modules/connectors/...
ok  	marketplace-central/apps/server_core/internal/modules/connectors/adapters/melhorenvio	3.945s
ok  	marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre	13.210s
ok  	marketplace-central/apps/server_core/internal/modules/connectors/application	2.347s
ok  	marketplace-central/apps/server_core/internal/modules/connectors/domain	1.832s
?   	marketplace-central/apps/server_core/internal/modules/connectors/events	[no test files]
?   	marketplace-central/apps/server_core/internal/modules/connectors/ports	[no test files]
?   	marketplace-central/apps/server_core/internal/modules/connectors/readmodel	[no test files]
ok  	marketplace-central/apps/server_core/internal/modules/connectors/transport	2.619s
```

All packages in the connectors module pass. No pre-existing test in the module regressed.

## Correction round: SLALimitAt (adversarial review finding)

An adversarial review of this feature's first pass found a real defect: `ShipmentDetail.SLALimitAt`
had been left nil with the stated reason "source is the unconfirmed `/sla` sub-resource." That
reasoning was wrong. `sla_limit_at` does NOT come from `GET /shipments/{id}/sla` — it comes from
`lead_time.estimated_delivery_limit.date` on the SAME primary `GET /shipments/{id}` call
already made (x-format-new: true header). Proof, both already present in this worktree before
this feature started:

- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/shipping_reader.go`
  (existing, live, untouched by this feature) already decodes this exact field today:
  `mlShipmentResponse.LeadTime` (line 19) → `mapShipmentInfo`'s `SLADue` (lines 236-237), via
  the pre-existing private types `mlShipmentLeadTime`/`mlEstimatedDeliveryLimit` (lines 77-83).
- `docs/design/handoff-2026-07/API-MAP.md:18` documents the SLA due-date as coming from
  `GET /shipments/$ID` (`estimated_delivery_*`), not a separate `/sla` call. Per this repo's
  truth order (`docs/` outranks `.mnfs/research/`), that doc wins over IC-03's `/sla` claim on
  this specific point — IC-03 is wrong here.

Fix applied (no new endpoint call, no new types, `capability_adapter.go`/`shipping_reader.go`
untouched):

1. `apps/server_core/internal/modules/connectors/adapters/mercado_livre/shipment_ingest_reader.go`
   — added a `LeadTime *mlShipmentLeadTime` field (`json:"lead_time"`) to
   `mlIngestShipmentResponse`, reusing the existing `mlShipmentLeadTime`/`mlEstimatedDeliveryLimit`
   types from `shipping_reader.go` (same package). `mapShipmentDetail` now populates
   `detail.SLALimitAt` from `shipment.LeadTime.EstimatedDeliveryLimit.Date` via `parseTimePtr`,
   with the same nil-guarded traversal `mapShipmentInfo` already uses (absent at any level stays
   nil, never fabricated).
2. `apps/server_core/internal/modules/connectors/domain/shipment_detail.go` — split the combined
   `SLAStatus`/`SLALimitAt` doc comment: `SLAStatus` keeps its original honest-unknown rationale
   (the `/sla` sub-resource is genuinely not called, shape unconfirmed); `SLALimitAt` now documents
   its real source (the primary payload's `lead_time`) and points to the shipping_reader.go
   precedent.
3. Test coverage added/extended in
   `apps/server_core/internal/modules/connectors/adapters/mercado_livre/shipment_ingest_reader_test.go`:
   - `TestGetShipmentDetailDecodesAllTargetFields` — fixture now includes
     `"lead_time":{"estimated_delivery_limit":{"date":"2026-07-25T23:59:59.000-04:00"}}`; asserts
     `SLALimitAt` equals the parsed UTC instant (`2026-07-26T03:59:59Z`). `SLAStatus` assertion
     kept separate and still asserts nil (that field's honest-unknown status is unchanged).
   - `TestGetShipmentDetailSLALimitAtAbsentStaysNil` (new) — payload with no `lead_time` block
     at all → `SLALimitAt == nil`, proving the gap still degrades honestly rather than erroring
     or defaulting when the source field truly is absent.

`SLAStatus`, `LogisticType`, `TrackingMethod` were left exactly as-is (nil) — no repo evidence
contradicts those three; only the `SLALimitAt` reasoning was wrong.

Verification re-run after the fix:

```
$ export GOCACHE="$(pwd)/.gocache" GOMODCACHE="$(pwd)/.gomodcache" && go build ./internal/modules/connectors/...
(no output — clean)
$ export GOCACHE="$(pwd)/.gocache" GOMODCACHE="$(pwd)/.gomodcache" && go vet ./internal/modules/connectors/...
(no output — clean)
$ export GOCACHE="$(pwd)/.gocache" GOMODCACHE="$(pwd)/.gomodcache" && go test ./internal/modules/connectors/adapters/mercado_livre/... -run 'TestGetShipmentDetail' -v
--- PASS: TestGetShipmentDetailDecodesAllTargetFields (0.07s)
--- PASS: TestGetShipmentDetailSLALimitAtAbsentStaysNil (0.07s)
--- PASS: TestGetShipmentDetailNotFoundIsHonestAbsence (0.01s)
    --- PASS: TestGetShipmentDetailNotFoundIsHonestAbsence/Not_Found (0.00s)
    --- PASS: TestGetShipmentDetailNotFoundIsHonestAbsence/Gone (0.00s)
--- PASS: TestGetShipmentDetailCostsDegradeOn404 (0.21s)
--- PASS: TestGetShipmentDetailCostsRealErrorPropagates (0.07s)
--- PASS: TestGetShipmentDetailPrimaryErrorPropagates (0.00s)
--- PASS: TestGetShipmentDetailEmptyShipmentIDIsInvalidReference (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre	3.185s
$ export GOCACHE="$(pwd)/.gocache" GOMODCACHE="$(pwd)/.gomodcache" && gofmt -l internal/modules/connectors/domain/shipment_detail.go internal/modules/connectors/adapters/mercado_livre/shipment_ingest_reader.go internal/modules/connectors/adapters/mercado_livre/shipment_ingest_reader_test.go
(no output — all 3 clean)
$ export GOCACHE="$(pwd)/.gocache" GOMODCACHE="$(pwd)/.gomodcache" && go test ./internal/modules/connectors/...
ok  	marketplace-central/apps/server_core/internal/modules/connectors/adapters/melhorenvio	(cached)
ok  	marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre	11.158s
ok  	marketplace-central/apps/server_core/internal/modules/connectors/application	(cached)
ok  	marketplace-central/apps/server_core/internal/modules/connectors/domain	(cached)
ok  	marketplace-central/apps/server_core/internal/modules/connectors/transport	(cached)
```

Full connectors module regression suite still green after the fix — 11 top-level test functions
now in the two new files (19 including subtests).

### F-01's own new tests (isolated re-run, verbose, pre-correction baseline)

10 top-level test functions, all PASS (17 including subtests) — this count predates the
`TestGetShipmentDetailSLALimitAtAbsentStaysNil` addition above:

- `TestGetOrderDetailDecodesAllTargetFields` — PASS. Decodes every IC-03/0089 field
  (status/status_detail/timestamps/pack_id/shipping.id/buyer.nickname/tags/order_items
  qty=3+qty=1/sale_fee per-unit/payments), asserts `RawOrder` never contains the substring
  `"billing_info"` nor the synthetic PII marker, and asserts non-fiscal buyer data + order_items
  content DO survive the redaction (byte-exact, no precision loss).
- `TestGetOrderDetailAbsentShippingIDStaysNil` — PASS. Payload without `shipping.id` →
  `ShippingID == nil`, `PackID == nil`, `BuyerNickname == nil`; raw still billing_info-free.
- `TestGetOrderDetailMapsProviderErrors` — PASS (5 subtests: 403, 401, 404, 429, 500). 403
  (third-party order) and 401 both map to `domain.ErrUnauthorized`, typed and non-retryable,
  detail stays zero-value.
- `TestGetOrderDetailEmptyProviderOrderIDIsInvalidReference` — PASS.
- `TestGetShipmentDetailDecodesAllTargetFields` — PASS (pre-correction: also asserted
  `SLALimitAt` nil; post-correction it asserts `SLALimitAt` populated from `lead_time` — see
  "Correction round" above). Decodes every 0088 column this reader fills
  (status/substatus/tracking_number/date_created/destination block/costs/sla_limit_at), and
  asserts the remaining honest-gap fields (`LogisticType`, `TrackingMethod`, `SLAStatus`) stay
  nil.
- `TestGetShipmentDetailNotFoundIsHonestAbsence` — PASS (2 subtests: 404, 410). Both →
  `domain.ShipmentDetail{}`, `Found == false`, `err == nil`.
- `TestGetShipmentDetailCostsDegradeOn404` — PASS. Primary succeeds, costs 404s → shipment
  still returned with `Found=true`, cost fields nil, rest of shipment intact.
- `TestGetShipmentDetailCostsRealErrorPropagates` — PASS. Costs 500 → `domain.ErrProviderUnavailable`
  propagates, detail zero-value.
- `TestGetShipmentDetailPrimaryErrorPropagates` — PASS. Primary 403 → `domain.ErrUnauthorized`.
- `TestGetShipmentDetailEmptyShipmentIDIsInvalidReference` — PASS.

## gofmt

```
$ gofmt -l internal/modules/connectors/domain/order_detail.go internal/modules/connectors/domain/shipment_detail.go internal/modules/connectors/adapters/mercado_livre/order_ingest_reader.go internal/modules/connectors/adapters/mercado_livre/shipment_ingest_reader.go internal/modules/connectors/adapters/mercado_livre/order_ingest_reader_test.go internal/modules/connectors/adapters/mercado_livre/shipment_ingest_reader_test.go
(no output — all 6 files clean)
```

## Honest gaps left nil (not blocking, per feature.md's own escape hatch)

`SLALimitAt` is NO LONGER in this table (see "Correction round" above) — it is now populated
from the primary payload's `lead_time.estimated_delivery_limit.date`, the same field the
existing `mapShipmentInfo` already decodes.

| Field | Why |
|---|---|
| `ShipmentDetail.LogisticType` | Top-level JSON key for the v2 `logistic_type` vocabulary (fulfillment/cross_docking/self_service/drop_off) could not be confirmed against any fixture/doc present in this worktree. A legacy `mode` field IS present in an in-repo fixture but is a different, older concept — aliasing it would be a guess. |
| `ShipmentDetail.TrackingMethod` | Same: no confirmed JSON key in this worktree. |
| `ShipmentDetail.SLAStatus` | IC-03 names `GET /shipments/{id}/sla` as the source; that endpoint's existence is live-verified (fact #15) but its response shape is not — the referenced evidence file is absent (PII-scrub debt, see repo memory `harness-debts-file`). The endpoint is not called at all rather than called and guessed. |

None of these gaps block this feature — all are structural absences (nil pointers), which is
the ADR-17-compliant outcome for an unconfirmed fact, not an error condition.

## Not touched (verified)

- `capability_adapter.go` — 0 edits.
- `shipping_reader.go` — 0 edits; its `getShipmentInfo`/`mlShipmentResponse` remain untouched
  and still serve their existing call site.
- `internal/modules/orders/**`, `root.go` — 0 edits (F-02/F-03 scope).
- No git commands were run (no commit, no push, no stage) — working tree left for the
  orchestrator to review.
