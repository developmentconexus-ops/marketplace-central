# CHIP-M08-SHIPFIX — Evidence Pack

Milestone-correction chip against M-08 pedidos (merged @e2a5570). Fixes the
`/pedidos` "A ENVIAR 24 / ENVIADOS 0" defect: shipped/delivered orders never
leaving the *enviar* bucket.

- **Branch:** `chip/m08-shipfix`
- **Base SHA (governance anchor):** `8b9e0ec01cae7b64f68001054464b3cb7f534160`
- **Scope:** A (connector shipment decode) · B (bucket semantics) · C (FE-derived KPI)
- **Operator-surfaced scope addition:** order `tags` (`delivered`) as a bucket
  signal + shipment `substatus` threading — see HAND-BACK / ledger. Flagged to hub.

## Root causes → fixes

| # | Root cause | Fix | Files |
|---|------------|-----|-------|
| 1 | ML `/shipments/{id}/costs` called without required `x-format-new: true` → legacy body → `ProviderPayloadInvalid` | Send `x-format-new: true` on the **costs** request only (primary GET works header-less; not touched) via new per-request header path in the capability adapter | `capability_adapter.go`, `shipping_reader.go` |
| 2 | `getShipmentInfo` treated costs decode-invalid as FATAL (only 404 degraded) → whole shipment read (incl. status) sank | Degrade on `ProviderInvalidReference` **and** `ProviderPayloadInvalid` → return shipment-with-status, `Costs=nil`. 5xx still sinks (transient). | `shipping_reader.go` |
| 3 | `DeriveOrderBucket` read only `provider_status`, which ML pins at `paid` for the whole post-payment lifecycle → `enviado` unreachable | New signature consumes `shipmentStatus` + order `tags`: `delivered` tag → enviado (survives shipment-lookup failure); `shipped`/`delivered` status → enviado; `ready_to_ship`/`handling`/`pending` → enviar. `provider_status` keeps sole authority over novo/faturar/cancelado. | `order_bucket.go`, `read_model.go`, `order_repo.go`, `enrich_service.go`, `http_handler.go` |
| 4 | KPI cards fed by shipment-status-blind `GetOrderBucketCounts` SQL → diverged from Lista | KPI derived on FE from the already-fully-loaded enriched list via the same `bucketTabCount` helper the Lista tabs use → **KPI == Lista by construction**. No migration, no backend fan-out. | `PedidosPage.tsx`, `pedidosTabs.ts` |

Substatus is captured, threaded, and exposed on the `OrderRastreio` DTO for
display, but is deliberately **not** a `DeriveOrderBucket` parameter — every
`shipped` substatus maps to `enviado`, so it would be a dead bucketing param
(anti-slop). Additive-only contract change (optional `substatus`).

## Failing-test-first

- `TestGetShipmentInfoCostsLegacyShapeDegradesToStatusOnly` — `/costs` returns a
  legacy-shape body (`senders` object, not the new-format array) → decode fails →
  read degrades to shipment-with-status (`shipped` + substatus `receiver_absent`),
  `Costs=nil`. Reproduces the pre-fix ProviderPayloadInvalid sink.
- `TestGetShipmentInfoCostsRequestSendsXFormatNewHeader` — asserts `x-format-new:
  true` rides the `/costs` request and is **absent** on the primary GET.
- `order_bucket_test.go` — exhaustive truth table over
  {providerStatus, shipmentStatus, tags, hasShipment}; adds delivered-tag and
  shipment-status cases.
- `order_repo_scan_test.go` — tags_json scan (populated / NULL→nil, ADR-17).
- `PedidosPage.test.tsx` — KPI counts derived from the list == Lista tab counts;
  list-load failure → honest "—", never "0".

## Verification ladder (L0/L1) — GREEN

Run in worktree, `GOCACHE=C:/Users/leandro.theodoro/Documents/marketplace-central/.gocache` (absolute).

| Lane | Command | Result |
|------|---------|--------|
| Go build | `go build ./...` | clean |
| Go vet | `go vet ./internal/modules/connectors/... ./internal/modules/orders/...` | clean |
| Go test (touched) | `go test ./internal/modules/connectors/... ./internal/modules/orders/...` | ok (all pkgs) |
| Go test (full) | `go test ./...` | EXIT=0 |
| Web vitest (full) | `vitest run` | 45 files / 312 tests passed |
| Web vitest (pedidos) | `vitest run src/pages/pedidos` | 22 tests passed |
| SDK typecheck | `tsc --noEmit` (packages/sdk-runtime) | EXIT=0 |
| Web build | `vite build` | ✓ built, EXIT=0 |

**tsc raw-invocation note (false alarm):** `tsc --noEmit -p apps/web/tsconfig.json`
emits `TS2339` jest-dom matcher errors (`toBeInTheDocument`, `toBeDisabled`,
`toHaveTextContent`) across `*.test.tsx` — a repo-wide setup-type gap that also
hits AppRouter/Header test files on base, unrelated to this change. Zero errors
in touched production `.ts/.tsx`. Authoritative FE gate = `vite build` (green).

## Governance lane — PASSED

`npm run harness:governance -- -BaseSha 8b9e0ec01cae7b64f68001054464b3cb7f534160`
(clean worktree, full 40-hex anchor) → `status=passed`. The `OrderRastreio.substatus`
addition is optional/additive → no semantic drift. All emitted `baseline_exception`
lines are the pre-existing repo baseline (direct-reader secrets, module edges,
migration-prefix, production-panic allowlist) — none introduced by this chip.

## P6 dual gate — PASS / PASS (agreement)

Two independent read-only reviewers of committed diff `8b9e0ec0..aa4bb11`.

**Cold correctness (Opus, did not write the code): PASS.** Verified: x-format-new
on /costs only; degrade on ProviderPayloadInvalid|InvalidReference, 5xx still
sinks; bucket precedence coherent (cancel > delivered-tag > shipment-status >
provider-status fallback); both DeriveOrderBucket call sites consistent; tags
genuinely populated end-to-end (migration 0027 `tags_json` jsonb → scanReadModel
→ EnrichedOrder.Order.Tags → handler); ADR-17 honest-unknown throughout;
contract/SDK additive+consistent; tests exercise the claimed behavior (not theater).

**Adversarial slop (sonnet): PASS, no BLOCKER/MAJOR.** Confirmed error
classification, bucket priority, contract/SDK sync, nil-safe FE chaining, and
that the KPI-from-list rewrite fixes the defect.

### Adjudicated review notes (non-blocking)

1. **`not_delivered` top-level shipment status → enviar** (Opus). A failed-delivery
   shipment falls through to the provider fallback and buckets as *enviar*. Rare,
   genuinely ambiguous (no live *devolucao* bucket yet — placeholder), orthogonal
   to the shipped/delivered defect. **Design judgment, left as-is; flagged to hub
   as a follow-up if `not_delivered` appears in real data.**
2. **Silent `json.Unmarshal` on `tags_json`** (sonnet, both call sites).
   **Intentional per ADR-17** — corrupt/absent tags degrade to nil (honest
   unknown), never fabricated. Reviewer acknowledged as intentional. Left as-is.
3. **Happy-path substatus not asserted** (sonnet). VALID — actioned: added
   `Substatus == "ready_to_print"` assertion to
   `TestGetShipmentInfoMapsShipmentAndCosts` (removes a dangling mock field,
   guards happy-path substatus mapping). ML connector package re-run green.

## HAND-BACK

Hub owns the live P7 browser re-drive of `/pedidos` against the demo dataset
(KPI + Lista + drawer substatus display) — this chip cannot boot the server.
