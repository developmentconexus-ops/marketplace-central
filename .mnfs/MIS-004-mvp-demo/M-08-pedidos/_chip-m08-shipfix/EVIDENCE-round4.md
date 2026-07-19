# CHIP-M08-SHIPFIX — Round 4 (FINDING-P7-SHIPDECODE-4)

Correction round 4. Round 3 shipped @b40335f2 (merged); the shipment read now
succeeds and costs map cleanly. The hub doc-grounded re-inspection (deepmap
`ml-api-shipments-orders-deepmap.md`) found the read still surfaces **no
destino / carrier / real-freight** facts, and the destination it *does* decode
targets a field name that does not exist in the documented new-format schema.

- **Branch:** `chip/m08-destino4`
- **Base SHA (governance anchor):** `62208dfef8c9474aea7af33f32cea0112b501395`

## Root cause (three defects)

1. **Destino field-name bug (4th instance of the fabricated-field class).** The
   adapter decoded `receiver_address.state.id` — a field **absent** from the
   documented new-format `GET /shipments/{id}`. The real payload carries
   `destination.shipping_address.state{id,name}` + `city{id,name}` + `zip_code`,
   and `destination.receiver_name`. So `DestinationUF` was silently never
   populated from the real body, and city / zip / receiver were not surfaced at
   all — the ERP needs them.
2. **Carrier / tracking never fetched.** Carrier + tracking URL live on a
   dedicated sub-resource `GET /shipments/{id}/carrier` (`x-format-new:true`)
   → `{url,name}`. The adapter never called it, so `rastreio` carried no
   transportadora / URL.
3. **Costs never surfaced.** `mapShipmentCosts` output (round 3) reached
   `domain.ShipmentInfo.Costs` but **stopped there** — nothing mapped it onto
   the order DTO, so the real freight actuals were invisible to the client.

## Fix

### (1) Destination decode — adapter + neutral domain
`mercado_livre/shipping_reader.go`: replaced the fabricated `receiver_address`
structs with `mlDestination{ receiver_name, shipping_address{ state{id,name},
city{id,name}, zip_code } }` copied verbatim from the doc shape. `mapShipmentInfo`
populates the **marketplace-neutral** domain fields:
- `DestinationUF` — `state.id` with the `BR-` prefix trimmed (`"BR-RJ"→"RJ"`,
  bare `"SP"` passes through), guarded so an empty/masked value degrades to nil.
- `DestinationCity` ← `city.name`, `DestinationZip` ← `zip_code`,
  `ReceiverName` ← `receiver_name`, each via `trimmedPtr` (nil or trims-to-empty
  → nil). **ML obfuscates the buyer address until payment is confirmed** — a
  masked value is honest absence (`—`), never a WARN/error (ADR-17).

`domain.ShipmentInfo` gained neutral pointer fields only (`DestinationCity/Zip`,
`ReceiverName`, `CarrierName`, `TrackingURL`). **No provider shape name
(`destination`/`shipping_address`/`x-format-new`) crosses the adapter boundary.**

### (2) Carrier sub-resource — degrade-only
New `fillShipmentCarrier` GETs `…/carrier` with `x-format-new:true` →
`{url,name}` into `CarrierName`/`TrackingURL`. It runs **before** the costs GET,
so a costs degrade cannot drop the carrier. **Every** failure (404, decode error,
transport/auth) degrades to nil carrier while the shipment status/costs survive —
a carrier 404 is the documented *normal* "none" (doc precedent: `/delays` 404 =
no delay), never a WARN. Documented divergence from the costs sink-on-5xx: carrier
is peripheral; status/bucket survival is the paramount M-08 invariant.

### (3) Costs + destino surfacing onto the order DTO
`orders/application/enrich_service.go`: `ShipmentEnrichment` gained
`DestinationCity/Zip`, `ReceiverName`, `CarrierName`, `TrackingURL`, and
`Costs *connectorsdomain.ShipmentCosts`. `resolveShipment` also sets
`buyer.City` from `DestinationCity` (blessed source = shipment, per the
`MaskedBuyer` doc comment).
`orders/transport/http_handler.go`: additive DTO fields
`destino_cep`, `destinatario`, `frete_real{bruto,receiver,sender}`, plus
`rastreio.transportadora` + `rastreio.url_rastreio`. `mapFreteReal` converts each
`Money` decimal-string amount → `*float64` (nil Money / unparseable → nil, never
a fabricated 0); when **all three** amounts are nil the whole `frete_real` block
is omitted (honest absence, not an empty `{}`).

### Contract (same commit)
`contracts/api/marketplace-central.openapi.yaml`: `OrderRead` += `destino_cep`,
`destinatario`, `frete_real` (new `OrderFreteReal` schema, nullable doubles);
`OrderRastreio` += `transportadora`, `url_rastreio`.
`packages/sdk-runtime/src/index.ts`: 1:1 mirror (`OrderFreteReal` interface,
optional/`number|null` fields matching the Go json tags).

Scope stayed inside `connectors/adapters/mercado_livre` +
`connectors/domain/shipping_read.go` + orders enrich/DTO + contract/SDK. No
migration, no `modules.json` delta. No ESCALATION needed.

## Failing-test-first (RED → GREEN)

RED: `TestGetShipmentInfoDecodesNewFormatNumericID` (and the primary
`…MapsShipmentAndCosts`) failed `DestinationUF = nil, want SP` the moment the
fixtures were rewritten to the real `destination` shape — reproducing the
production defect (the old `receiver_address` decode path never saw the field).
GREEN after decoding `destination`.

New adapter tests (fixtures verbatim from the deepmap doc):
- `TestGetShipmentInfoMaskedDestinationDegradesToNil` (SHIP-9) — blanked
  state/city/zip/receiver → all four nil, `err == nil` (masked ≠ error).
- `TestGetShipmentInfoCarrierNotFoundLeavesCarrierNil` (SHIP-10) — carrier 404 →
  nil carrier, status + costs survive, no WARN.
- `TestGetShipmentInfoCarrierSurvivesCostsDegrade` (SHIP-11) — undecodable costs
  → `Costs` nil but `CarrierName` populated (ordering proof).
- Carrier happy-path body (`{"url":"http://tracking.totalexpress.com.br/…",
  "name":"Total Express"}`) copied verbatim; `…RequestsSendXFormatNewHeader` now
  also asserts the `/carrier` call sends `x-format-new:true`.

New surfacing tests:
- `TestMapEnrichedOrderSurfacesShipmentDestinoCarrierFrete` — asserts
  `destino_cep`/`destinatario`/`destino_uf`, `rastreio.transportadora`/
  `url_rastreio`, and `frete_real.{bruto,receiver,sender}` all reach the JSON.
- `TestMapEnrichedOrderOmitsShipmentFactsWhenAbsent` — a bare shipment omits
  every destino/carrier/frete field (honest absence).
- `TestMapEnrichedOrderOmitsFreteRealWhenAllAmountsNil` — non-nil Costs with all
  amounts nil omits `frete_real` entirely.
- `enrich_service_test.go` subtest — city/zip/receiver + carrier + costs
  propagate to `ShipmentEnrichment` and `buyer.City`.

## Ladder (GOCACHE absolute, no GOFLAGS / workspace mode) — GREEN

| Lane | Result |
|------|--------|
| `go build ./...` | exit 0 |
| `go vet ./...` | exit 0 |
| `go test ./...` | exit 0 |
| `gofmt -l` (LF-normalized, touched files) | clean |
| `tsc --noEmit` (sdk-runtime) | exit 0 |
| OpenAPI YAML parse + schema presence | OK |

## Governance lane — PASSED

`harness:governance -BaseSha 62208dfef8c9474aea7af33f32cea0112b501395` (run from
the clean chip worktree, post-commit @af4ff9e) → `status=passed`. Only
pre-existing `baseline_exception=*` entries listed (direct-reader secrets,
migration-0021 dup, module-edge adapters, production-panic pricing) — none
touched by this delta. No new drift: the contract change is purely additive
(OpenAPI + SDK in the same commit), no `modules.json` / migration delta.

## P6 dual gate — PASS / PASS (agreement)

- **Cold correctness (Opus):** PASS. Verified destination decode matches the doc
  shape; `BR-` trim + masked→nil; carrier fills before costs and swallows all
  failures while status/costs survive; `frete_real` money→*float64 with nil-safe
  parse; no provider-shape leakage; all three required fixtures present;
  contract+SDK additive and consistent; no scope creep / test theater. One
  non-blocking note (empty `frete_real:{}` when Costs present but all amounts nil).
- **Adversarial (sonnet):** PASS. Attacked ordering, masked-destination, UF trim
  edge (`BR-`), 404-as-WARN, money parse, provider-shape leakage grep, fixture
  fidelity, contract drift, test theater, scope — could not break it. Confirmed
  reverting the field name to `receiver_address` breaks the real-`destination`
  fixture assertions (mutation-sensitive).
- **Action:** the Opus non-blocking note was hardened — `mapFreteReal` now
  returns nil (omits the block) when all three amounts are nil, covered by
  `TestMapEnrichedOrderOmitsFreteRealWhenAllAmountsNil`. Full ladder re-run GREEN.

## HAND-BACK

Hub owns merge (onto main tip; base anchor `62208dfe`) + the live P7 browser
re-drive — confirm shipped orders render Destino (UF/cidade/CEP/destinatário),
Rastreio (transportadora + URL), and Frete real. Regenerate the SDK from the
updated OpenAPI as part of the merge if the build pipeline requires it.
Heads-up for CHIP-M08-BUYER: the DTO additions are tight/additive (three new
`OrderRead` fields + two `OrderRastreio` fields + `OrderFreteReal`) so the
billing-info follow-up rebases cleanly. Chip cannot boot the server / has zero
ML writes. Chip did NOT push.
