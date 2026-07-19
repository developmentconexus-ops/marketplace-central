# CHIP-M08-SHIPFIX — Round 2 (FINDING-P7-SHIPDECODE-2)

Correction round 2. Round-1 shipped @2c13f6ad (merged); hub P7 live-drive found
the KPI/bucket defect fixed but the **live shipment read still failed on every
order** (`CONNECTORS_PROVIDER_PAYLOAD_INVALID`) — bucket was carried only by the
order `delivered`-tag fallback; Rastreio/Destino/substatus rendered "—".

- **Branch:** `chip/m08-shipdecode2` (forked fresh off merged main)
- **Base SHA (governance anchor):** `2c13f6adf82ea9b461f62daba663e6f9b50a4a2c`
- **Fix commit:** `5a5a305`

## Root cause (ML official docs, Context7 `/websites/developers_mercadolivre_br_pt_br`)

1. **Primary `GET /shipments/{id}` also requires `x-format-new: true`** — verbatim
   (gerenciamento-de-envios): *"ao fazer uma requisição GET, é necessário enviar o
   header 'x-format-new: true'."* Round-1 scope (hub) targeted only `/costs`. Without
   the header ML returns the LEGACY shape → `mlShipmentResponse` fails to unmarshal
   → `doJSON` ProviderPayloadInvalid → the whole read sinks **before** the round-1
   costs degrade is ever reached.
2. **New-format `id` is a JSON number** (docs status-de-pedidos-rastreamento:
   `"id": 28264263908`). The old `ID string` field cannot unmarshal a bare number,
   so even with the header the real payload sinks. Round-1 tests only used quoted
   string ids, so this never surfaced.

## Fix

- Primary GET now uses `doJSONWithHeaders(..., {"x-format-new":"true"}, &shipment)`.
  Primary decode failure stays a **hard error (honest)** — unlike costs there is no
  partial shipment to degrade to.
- `mlShipmentResponse.ID`: `string` → `flexString`, a custom type whose
  `UnmarshalJSON` accepts a JSON string OR a bare number (keeps the literal text);
  `null`/empty → "". Shape drift on `id` degrades the field, never the read (ADR-17).
- Costs call + degrade path **unchanged** from round 1.

## Failing-test-first (RED → GREEN)

- `TestGetShipmentInfoDecodesNewFormatNumericID` — primary body with numeric
  `"id":28264263908` + extra top-level fields. **RED on merged code**:
  `CONNECTORS_PROVIDER_PAYLOAD_INVALID` (the exact hub WARN). Green after flexString.
- `TestGetShipmentInfoRequestsSendXFormatNewHeader` — asserts `x-format-new:true` on
  **both** the primary and the costs request. **RED on merged code**: primary header "".

## Ladder (GOCACHE absolute, no GOFLAGS / workspace mode) — GREEN

| Lane | Result |
|------|--------|
| `go build ./...` | exit 0 |
| `go vet ./internal/modules/connectors/...` | clean |
| `go test ./...` | exit 0 |

No FE / OpenAPI / SDK / migration / modules.json delta this round (Go adapter only).

## Governance lane — PASSED

`harness:governance -BaseSha 2c13f6adf82ea9b461f62daba663e6f9b50a4a2c` → `status=passed`
(no contract/registry delta).

## P6 dual gate — PASS / PASS (agreement, no findings)

- **Cold correctness (Opus):** PASS. Verified flexString covers string/int/float/
  null/empty (bounds-safe, no nil-deref); header applied per-request via fresh
  `http.NewRequestWithContext`, zero leakage; primary-decode-failure-is-error
  invariant preserved; both new tests fail if the fix is reverted; no regression to
  SLA/DestinationUF/Delayed/substatus mapping.
- **Adversarial (sonnet):** PASS, no issues. Confirmed header isolation, hard-error
  on primary, costs degrade scoped to 404/payload-invalid, tests non-theater, ADR-17.

## HAND-BACK

Hub owns merge (onto main @2c13f6ad) + the live P7 browser re-drive — confirm the
`WARN orders: shipment lookup failed` line is gone and Rastreio/Destino/substatus
now render on shipped orders. Chip cannot boot the server / has zero ML writes.
