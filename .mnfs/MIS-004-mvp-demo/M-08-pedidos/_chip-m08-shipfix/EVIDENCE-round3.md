# CHIP-M08-SHIPFIX — Round 3 (FINDING-P7-SHIPDECODE-3)

Correction round 3. Round-2 shipped @9e3256d5 (merged); the unmarshal now succeeds
but the hub P7 live-drive still logged, on every order:

> `WARN orders: shipment lookup failed ... error="shipment_gross_amount.currency must be a non-empty upper-case string"`

- **Branch:** `chip/m08-shipcosts3`
- **Base SHA (governance anchor):** `9e3256d5709609d11ea1520e78b6710a73274734`

## Root cause

1. **The real new-format `GET /shipments/{id}/costs` body carries NO currency
   field** (ML docs `/websites/developers_mercadolivre_br_pt_br`,
   gerenciamento-de-envios: the cost body is `gross_amount` + `receiver` +
   `senders` only). Rounds 1–2 fixtures **fabricated** `currency_id`, masking this.
   In production `mlShipmentCostsResponse.Currency` decoded to `""` → every
   `providerMoney` call fed `domain.ValidateMoney` an empty currency → it rejects
   with the exact WARN above → the whole shipment read sank.
2. **`mapShipmentCosts` failure was FATAL** — it ran *after* the costs-decode
   degrade branch and returned `domain.ShipmentInfo{}, err`, so any mapping error
   sank the read and lost status/substatus/rastreio. This violated the round-1
   invariant *"costs failure degrades, status survives."*

## Fix (Go adapter only — `mercado_livre/shipping_reader.go`)

- **(a) Degrade invariant (structural fix).** `getShipmentInfo` now computes
  `result := mapShipmentInfo(shipment, a.now())` **before** costs mapping; on any
  `mapShipmentCosts` error it returns `result, nil` (status/substatus/UF survive,
  `Costs` nil) instead of sinking. The costs-DECODE degrade (404 /
  ProviderPayloadInvalid) is unchanged; **5xx / auth / rate-limit / transient
  still sink** (those are provider-error codes routed through
  `mapPricingReaderError`, never reachable by the new mapping-error degrade).
- **(b) Currency resolved at the adapter from the site.** New `mlSiteCurrency`
  resolver (ML-documented site→currency table: MLB→BRL, MLA→ARS, MLM→MXN, MLC→CLP,
  MCO→COP, MLU→UYU, MPE→PEN), doc-cited. Threaded as
  `mlSiteCurrency(firstNonEmpty(shipment.SiteID, a.siteID))` — prefer the
  shipment's own `site_id`, fall back to the **account-configured** site before the
  BRL default, so a non-BR tenant whose payload omits `site_id` is not silently
  relabelled BRL (unknown ≠ default; same `firstNonEmpty(input, a.siteID)` pattern
  the listing-price path already uses). The *amount* is still the real provider
  fact; only the currency label — a structural invariant of the ML site — is
  derived from a known fact.
- **(c) Fixtures.** Removed the fabricated `Currency`/`currency_id` field from
  `mlShipmentCostsResponse`; rewrote every `/costs` fixture to the real
  no-currency shape (`gross_amount` + `receiver{user_id,cost}` +
  `senders[]{user_id,cost}`). Assertions still expect BRL, now via the resolver.

No FE / OpenAPI / SDK / migration / modules.json delta. Scope stayed inside
`connectors/adapters/mercado_livre` (no ESCALATION needed).

## Failing-test-first (RED → GREEN)

RED on merged code (`go test ./.../mercado_livre/`):

```
--- FAIL: TestGetShipmentInfoMapsShipmentAndCosts
    error = shipment_gross_amount.currency must be a non-empty upper-case string
--- FAIL: TestGetShipmentInfoDecodesNewFormatNumericID   (same currency WARN)
--- FAIL: TestGetShipmentInfoAbsentOptionalsStayNil       (same currency WARN)
--- FAIL: TestGetShipmentInfoRequestsSendXFormatNewHeader (same currency WARN)
--- FAIL: TestGetShipmentInfoCostsMappingFailureDegradesToStatusOnly
    error = shipment_gross_amount.amount must be a non-empty decimal string, want degrade
```

The first four reproduce the **exact production WARN** once the fabricated
`currency_id` is removed from the fixtures. GREEN after the fix.

New tests:
- **`TestGetShipmentInfoCostsMappingFailureDegradesToStatusOnly`** — `/costs` body
  `{"gross_amount":2.4e1}` decodes cleanly into `*json.Number` (literal text kept)
  but exponent notation fails `decimalPattern`, so `mapShipmentCosts` errors
  *after* decode. Asserts status/substatus survive + `Costs` nil + nil error. This
  is the degrade invariant, and is distinct from `...LegacyShapeDegrades...`
  (which is a *decode* failure).
- Currency-resolution fidelity is covered by the rewritten happy-path fixtures
  (real no-currency `/costs` → costs mapped with BRL from the site).

## Ladder (GOCACHE absolute, no GOFLAGS / workspace mode) — GREEN

| Lane | Result |
|------|--------|
| `go build ./...` | exit 0 |
| `go vet ./internal/modules/connectors/...` | exit 0 |
| `go test ./...` | exit 0 |

## Governance lane — PASSED

`harness:governance -BaseSha 9e3256d5709609d11ea1520e78b6710a73274734` → `status=passed`
(Go-only additive delta, no contract/registry drift).

## P6 dual gate — PASS / PASS (agreement)

- **Cold correctness (Opus):** PASS. Verified degrade returns non-zero
  ShipmentInfo + nil Costs + nil error; 5xx/transport still sink (confirmed via
  status→code mapper, `TestGetShipmentInfoCostsServerErrorFails`); `mlSiteCurrency`
  panic-free + uppercase; SHIP-8 exponent genuinely reaches mapping (non-theater);
  `mapShipmentCosts` single caller, free-shipping path untouched. One soft note on
  default-on-unknown → actioned (see below).
- **Adversarial (sonnet):** one 🟡 — `mlSiteCurrency(shipment.SiteID)` skipped the
  account-configured `a.siteID` fallback, so a non-BR tenant with a blank payload
  `site_id` would be silently relabelled BRL (AGENTS.md "unknown ≠ zero/default").
  **VALID + convergent with the Opus soft note. ACTIONED:** threaded
  `firstNonEmpty(shipment.SiteID, a.siteID)`. Re-ran build/vet/test → GREEN. No
  other findings; degrade-vs-sink boundary, exponent-reaches-mapping, and
  fixture-fidelity all confirmed non-refuting.

## HAND-BACK

Hub owns merge (onto main @9e3256d5, currently 1 ledger commit ahead at fe23ba08)
+ the live P7 browser re-drive — confirm the `WARN ... shipment_gross_amount.currency`
line is gone and Rastreio/Destino/substatus/costs render on shipped orders. Chip
cannot boot the server / has zero ML writes. Chip did NOT push.
