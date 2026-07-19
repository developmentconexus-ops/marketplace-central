# CHIP-M08-LIST-TIMEOUT — Evidence Pack

Reopen of CHIP-M08-BUYER. **FINDING-M08-LIST-TIMEOUT** (demo-critical, T-1).

- **Mission:** MIS-004-mvp-demo · **Milestone:** M-08-pedidos
- **Branch:** chip/m08-list-timeout
- **BaseSha:** `5ab5714cf2b114ebb156896c213739e7e43662a0` (main tip)
- **HEAD:** `73a86c4e`

## Regression

The merged buyer-fiscal enrichment resolved the two-step ML billing-info lookup **per order** inside `EnrichService.Enrich`. Both handlers used `Enrich`:
- `handleList` (GET /orders) — 24-order real page
- `handleGet` (GET /orders/{id})

So the list did 24 × (+2 sequential ML calls) → `context deadline exceeded` @ ~15s → **504**, cascading into the shipment lookups that previously passed. `/pedidos` (demo screen) went dead ("Carregando…", counters "—"). Evidence: `curl` list = 504 @15.003s; pre-merge P7 round-4 PASSED.

## Ruling (hub)

Buyer fiscal is DRAWER (detail) data — the list must NEVER pay billing-info. Resolve fiscal only on the detail path.

## Fix (form: EnrichOne variant)

- `Enrich` (list path) drops `resolveBuyerFiscal` entirely → every `EnrichedOrder` carries nil `BuyerFiscal`. Restores pre-merge round-4 list cost.
- New `EnrichOne` (detail path): runs the same base enrichment for a single order, then resolves buyer fiscal for that one order only — bounding the two-step lookup to the drawer the user opened.
- `handleGet` → `EnrichOne`; `handleList` → `Enrich` (unchanged).

`EnrichOne` reuses `Enrich(ctx,id,[]{order})[0]` (always length-1, `[0]` panic-safe) and mutates only the `BuyerFiscal` field on the value copy — all base fields (masked buyer + shipment UF, item costs, profitability) survive. Degrade semantics (`resolveBuyerFiscal`: honest absence → nil, real error → warn-once → nil) untouched. comprador_fiscal DTO contract intact (detail emits via omitempty, list omits).

## Out of scope (documented, NOT touched)

The list's per-order **shipment** lookup is still sequential. It was marginal pre-merge and passed P7 round-4; the 504 is purely the buyer-fiscal addition. Bounding shipment with a limited errgroup is a separate optimization not required to clear this regression — left out to keep the T-1 demo fix minimal and reversible. Flagged for a future optimization chip if the list still reads slow under load.

## Tests (red-then-green)

- **RED:** `EnrichOne` undefined → package build-fails (watched).
- `TestEnrichServiceEnrich_ListNeverResolvesBuyerFiscal`: reader wired + loaded with data, `Enrich` over N orders → `reader.calls == 0` AND every element's `BuyerFiscal == nil` (not a weak calls-only assertion).
- `TestEnrichServiceEnrichOne_BuyerFiscal`: present / no-data / reader-error / nil-reader / empty-provider-id / composite(mask+cost+shipment) — detail path resolves exactly once, degrade behavior unchanged.

## Gate evidence

| Gate | Result |
|------|--------|
| `go build ./apps/server_core/...` | OK |
| `go vet` (orders/...) | OK |
| `go test` (orders application + transport + adapters) | all `ok`, no FAIL/panic |
| `gofmt -l` (3 touched files) | clean |
| Governance lane (`harness.ps1 governance -BaseSha 5ab5714c…`) | `status=passed` (only pre-existing baseline exceptions) |

## P6 dual gate (re-gate of delta)

| Reviewer | Verdict |
|----------|---------|
| Cold Opus (independent) | **PASS** — Enrich drops resolveBuyerFiscal; EnrichOne single-order, [0] panic-safe, all base fields preserved; handlers wired; DTO/degrade intact; zero writes/migration |
| Adversarial sonnet | **PASS** — grep confirms only EnrichOne reaches the reader; list test asserts nil per element; no shared-state mutation; gofmt/build/vet/test clean; FE list has no comprador_fiscal reference |

**Both PASS → agreement.**

## Files (delta, BaseSha..HEAD)

```
orders/application/enrich_service.go        (Enrich drops fiscal; add EnrichOne)
orders/application/enrich_service_test.go   (list-never test + EnrichOne suite)
orders/transport/http_handler.go            (handleGet → EnrichOne)
```

Full unified diff: `fix.diff` (this directory).

## Hub owns

- Merge to main.
- P7 re-drive: GET /orders (24-order list) loads without 504; drawer still shows comprador_fiscal + destinatario/CEP/frete_real/transportadora — the M-08 re-close gate (D-73 reversal).
