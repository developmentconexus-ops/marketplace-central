# F-03-read-path-switch — validation

All Go commands run from `apps/server_core` with `GOCACHE`/`GOMODCACHE` exported as absolute
paths, on Windows/Git Bash.

## go build / go vet / go vet -tags=integration

```
$ export GOCACHE="$(pwd)/.gocache" GOMODCACHE="$(pwd)/.gocache/mod"
$ go build ./...
(no output — clean)
$ go vet ./...
(no output — clean)
$ go vet -tags=integration ./...
(no output — clean)
```

First pass of `go vet ./...` failed with
`internal\composition\orders_adapters_test.go:119:14: undefined: newOrdersShipmentReaderAdapter`
— this test file held `TestOrdersShipmentReaderAdapter`, coverage for the site-A adapter that
this feature deletes from `orders_adapters.go`. Fixed by removing that test and its
single-use `erroringInstallationRepo` fixture (used nowhere else in the package — confirmed by
grep before deleting), then removing the two now-unused imports
(`integrationsapp`, `integrationsdomain`). This is a necessary, minimal consequence of the
site-A deletion the brief specifies (item 4), not scope creep.

## go test ./internal/modules/orders/... ./internal/composition/... ./internal/platform/archguard/...

```
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/integrations	2.7s
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/internalread	2.7s
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/postgres	3.8s
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/pricingtax	3.0s
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/productlinks	1.7s
ok  	marketplace-central/apps/server_core/internal/modules/orders/application	3.2s
ok  	marketplace-central/apps/server_core/internal/modules/orders/domain	3.0s
ok  	marketplace-central/apps/server_core/internal/modules/orders/ports	3.1s
ok  	marketplace-central/apps/server_core/internal/modules/orders/transport	3.8s
ok  	marketplace-central/apps/server_core/internal/composition	4.6s
ok  	marketplace-central/apps/server_core/internal/platform/archguard	0.7s
```

All packages pass, `-count=1` (no stale cache).

## gofmt

**FINDING (tooling gotcha):** `gofmt -d` intermittently produces an empty diff on this
Windows/Git Bash environment for a file that genuinely differs from canonical formatting —
confirmed reproducibly on `shipment_reader.go` (5 consecutive `gofmt -l` runs against a
freshly CRLF-stripped copy all flagged it as needing formatting, while `gofmt -d` against the
same copy printed nothing). Root cause not fully isolated (suspect `gofmt -d` shelling out to
an external `diff` binary that behaves oddly on this host, since `-l`, which does the
comparison internally, was consistent). Reliable method used instead: copy the CRLF-stripped
file, run `gofmt -w` on the copy, then `diff -u` the two copies directly — this caught a real
misalignment in `shipment_reader.go` (struct-field column alignment inside `scanShipmentRow`,
introduced by a hand-edit) that `gofmt -d` was silently hiding. Do not trust a clean `gofmt -d`
alone on this host; corroborate with the copy-and-diff method above.

This repo checks out with CRLF (`core.autocrlf=true`) but `gofmt` normalizes to LF, so a raw
`gofmt -l`/`-w` against the checked-out file always reports a "diff" from line endings alone —
every check below strips `\r` first (`tr -d '\r' < file > tmp.go`) before invoking gofmt.

16 touched/created files checked with the copy-and-diff method:

```
order_repo_bucket_counts_test.go, shipment_reader_test.go, buyer_fiscal_reader_test.go,
shipment_reader.go, buyer_fiscal_reader.go, zero_live_ml_call_test.go,
ports/shipment_reader.go, enrich_service.go, enrich_service_test.go, order_repo.go,
root.go, orders_adapters.go, orders_adapters_test.go, archguard_test.go,
testdata/three_sites/root.go, transport/http_handler_test.go
```

First pass found 3 real misalignments: `shipment_reader_test.go` (struct field + literal
column alignment), `shipment_reader.go` (struct field column alignment inside
`scanShipmentRow`), `archguard_test.go` (one map-literal value column in `mlExcludedSymbols`
after a key got a 2-character-longer name). All three are gofmt artifacts of hand-edited
column alignment, not logic changes. Fixed with targeted `Edit` calls. Second pass:
`ALL CLEAN (diff-verified)`.

## archguard before/after

**Before** (baseline at the start of this feature — F-01/F-02 already landed, deleting site A
from root.go's wiring was still pending):

```go
var mlAllowlist = []Site{
	{Name: "orders.list.shipment", File: "internal/composition/root.go", Symbol: "newOrdersShipmentReaderAdapter"},
	{Name: "orders.detail.buyer_fiscal", File: "internal/composition/root.go", Symbol: "newOrdersBuyerFiscalReaderAdapter"},
	{Name: "pricing.decompose.category_resolver", File: "internal/composition/root.go", Symbol: "newPricingCategoryResolverAdapter"},
	{Name: "pricing.solve.commission_quoter", File: "internal/composition/root.go", Symbol: "newPricingCommissionQuoterAdapter"},
}
var mlExcludedSymbols = map[string]string{
	"newMarketPriceIntelCollectorAdapter": "market/* competitor-price collection is MIS-008 scope, not MIS-007 F-04; tracked separately",
}
```

**After** (current):

```go
var mlAllowlist = []Site{
	{Name: "pricing.decompose.category_resolver", File: "internal/composition/root.go", Symbol: "newPricingCategoryResolverAdapter"},
	{Name: "pricing.solve.commission_quoter", File: "internal/composition/root.go", Symbol: "newPricingCommissionQuoterAdapter"},
}
var mlExcludedSymbols = map[string]string{
	"newMarketPriceIntelCollectorAdapter":  "market/* competitor-price collection is MIS-008 scope, not MIS-007 F-04; tracked separately",
	"newOrdersOrderDetailReaderAdapter":    "batch-only: feeds ordersIngestSvc.OrderDetail for POST /orders/import (F-02), not an interactive request-time path",
	"newOrdersShipmentDetailReaderAdapter": "batch-only: feeds ordersIngestSvc.ShipmentDetail for POST /orders/import (F-02), not an interactive request-time path",
	"newOrdersBuyerFiscalReaderAdapter":    "RECLASSIFIED by F-03 (was allowlist entry orders.detail.buyer_fiscal): the interactive GET /orders/{id} path no longer calls this constructor's result -- it now reads Postgres via orders/adapters/postgres.BuyerFiscalReader. The remaining root.go call site feeds ordersIngestSvc.BuyerFiscal for POST /orders/import batch ingest only (F-02), same batch-only reasoning as the two entries above",
}
```

`mlAllowlist`: 4 -> 2 entries (site A fully retired — deleted, not merely excluded; site B
reclassified batch-only). `mlExcludedSymbols`: 1 -> 4 entries. Doc comment above
`mlAllowlist` rewritten to describe sites A/B's retirement/reclassification with file:line
citations.

### Baseline test output (site A present, before deletion)

```
raw detector found 7 call site(s) in ../../composition/root.go, expected exactly 4 (allowlist) + 1 (documented exclusions) = 5; a new, undocumented site appeared -- raw hits: [{Symbol:newOrdersBuyerFiscalReaderAdapter File:../../composition/root.go Line:583} {Symbol:newOrdersOrderDetailReaderAdapter File:../../composition/root.go Line:586} {Symbol:newOrdersShipmentDetailReaderAdapter File:../../composition/root.go Line:587} {Symbol:newOrdersShipmentReaderAdapter File:../../composition/root.go Line:604} {Symbol:newMarketPriceIntelCollectorAdapter File:../../composition/root.go Line:696} {Symbol:newPricingCategoryResolverAdapter File:../../composition/root.go Line:859} {Symbol:newPricingCommissionQuoterAdapter File:../../composition/root.go Line:860}]
--- FAIL: TestRealRepoInteractiveMLSites_MatchesAllowlist (0.00s)
FAIL
FAIL	marketplace-central/apps/server_core/internal/platform/archguard	1.266s
FAIL
```

### Current test output (clean)

```
=== RUN   TestRealRepoInteractiveMLSites_MatchesAllowlist
    archguard_test.go:435: baseline PASS: 2 interactive Mercado Livre site(s) in ../../composition/root.go match mlAllowlist exactly (4 documented exclusion(s) filtered):
    archguard_test.go:437:   newPricingCategoryResolverAdapter        ../../composition/root.go:867
    archguard_test.go:437:   newPricingCommissionQuoterAdapter        ../../composition/root.go:868
--- PASS: TestRealRepoInteractiveMLSites_MatchesAllowlist (0.00s)
=== RUN   TestRealRepoTransportAndApplication_NeverImportMLDirectly
    archguard_test.go:457: ../../modules/orders: 0 direct Mercado Livre imports (as expected)
    archguard_test.go:457: ../../modules/pricing: 0 direct Mercado Livre imports (as expected)
--- PASS: TestRealRepoTransportAndApplication_NeverImportMLDirectly (0.02s)
=== RUN   TestFixture_FifthSiteIsDetectedAndNamed
--- PASS: TestFixture_FifthSiteIsDetectedAndNamed (0.00s)
=== RUN   TestFixture_AliasedSiteIsDetectedAndNamed
--- PASS: TestFixture_AliasedSiteIsDetectedAndNamed (0.00s)
=== RUN   TestFixture_ShrunkAllowlistStillPasses
    archguard_test.go:569: must-fail PASS (state 3): shrunk 1-entry allowlist matches the 1-site fixture exactly (guard shrinks structurally, not hardcoded to 4)
--- PASS: TestFixture_ShrunkAllowlistStillPasses (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/platform/archguard	1.073s
```

`TestFixture_ShrunkAllowlistStillPasses` required rewriting `testdata/three_sites/root.go`
(hardcoded the now fully-retired site A and reclassified site B, which no longer belong in
the shrunk fixture) to mirror only the current real wiring — a necessary, minimal consequence
of the mlAllowlist shrink (item 5 of the brief), not scope creep.

### Must-fail proof (fresh, re-run in this session, not reused from an older capture)

Temporarily re-inserted, immediately before the `ordersEnrichShipmentReader :=` line in
root.go: `_ = newOrdersShipmentReaderAdapter(mercadoLivreCapabilities, installationSvc,
cfg.DefaultTenantID)` — a call shaped exactly like the deleted site A.

RED:

```
$ go test ./internal/platform/archguard/... -run TestRealRepoInteractiveMLSites_MatchesAllowlist -v
=== RUN   TestRealRepoInteractiveMLSites_MatchesAllowlist
    archguard_test.go:426: raw detector found 7 call site(s) in ../../composition/root.go, expected exactly 2 (allowlist) + 4 (documented exclusions) = 6; a new, undocumented site appeared -- raw hits: [{Symbol:newOrdersBuyerFiscalReaderAdapter File:../../composition/root.go Line:583} {Symbol:newOrdersOrderDetailReaderAdapter File:../../composition/root.go Line:586} {Symbol:newOrdersShipmentDetailReaderAdapter File:../../composition/root.go Line:587} {Symbol:newOrdersShipmentReaderAdapter File:../../composition/root.go Line:611} {Symbol:newMarketPriceIntelCollectorAdapter File:../../composition/root.go Line:705} {Symbol:newPricingCategoryResolverAdapter File:../../composition/root.go Line:868} {Symbol:newPricingCommissionQuoterAdapter File:../../composition/root.go Line:869}]
--- FAIL: TestRealRepoInteractiveMLSites_MatchesAllowlist (0.00s)
FAIL
FAIL	marketplace-central/apps/server_core/internal/platform/archguard	1.218s
FAIL
```

The RED names the exact injected symbol (`newOrdersShipmentReaderAdapter`) and its exact
line (`root.go:611`, where it was inserted).

Reverted (removed the inserted line). GREEN confirmed:

```
=== RUN   TestRealRepoInteractiveMLSites_MatchesAllowlist
    archguard_test.go:435: baseline PASS: 2 interactive Mercado Livre site(s) in ../../composition/root.go match mlAllowlist exactly (4 documented exclusion(s) filtered):
--- PASS: TestRealRepoInteractiveMLSites_MatchesAllowlist (0.00s)
```

Also confirmed `go build ./...` clean both before the must-fail proof (with the extra line
present — note: the raw AST scan does not require the referenced symbol to exist, so this
would NOT have compiled had it been left in; it was never left in a build-checked state) and
after the revert.

## GetOrderBucketCounts — before/after (per-order attribution)

Fixture (`order_repo_bucket_counts_test.go`,
`TestGetOrderBucketCountsUsesPersistedBucketWithLegacyFallback`, run against real Postgres via
`npm run harness:integration`):

| order | provider_status | tags | persisted `bucket` column | OLD behavior (always re-derive, shipmentStatus="") | NEW behavior (two-tier read) | delta |
|---|---|---|---|---|---|---|
| persisted-enviado | paid | none | `"enviado"` | **Faturar** (paid, not delivered-tagged, faturado_at NULL) | **Enviado** (reads persisted column) | **CORRECTION** — OLD could never see the real "shipped" shipment status F-02 recorded, because it always re-derived with a hardcoded empty shipment-status string |
| legacy-faturar | paid | none | `""` (empty/NULL — simulates a pre-F-02 row) | Faturar | Faturar (falls back to re-derivation, byte-identical to OLD) | none |
| persisted-novo | created | none | `"novo"` | Novo | Novo (persisted column agrees with re-derivation) | none |

Aggregate counts: `want := ports.OrderBucketCounts{Novo: 1, Faturar: 1, Enviar: 0, Enviado: 1}`
— asserted equal to `GetOrderBucketCounts`'s actual return.

The test also asserts the persisted bucket columns really hold what the fixture claims
(`assertPersistedBucketColumn`, guards against the test accidentally not exercising the
two-tier branch it claims to), and cross-checks that re-deriving persisted-enviado with the
OLD hardcoded `shipmentStatus=""` genuinely produces `BucketFaturar` (proves the delta is
real and attributable, not incidental — if this assertion ever failed, the test's own premise
would be broken and it says so in its failure message).

Confirmed via `npm run harness:pg:up` (session-reused Postgres, port 50486) then
`npm run harness:integration` (self-discovering, globs every `//go:build integration`
package). Final confirming run (this session, after all edits including the must-fail
round-trip on root.go):

```
target=ephemeral-postgres
status=passed
```

Clean pass, no `failure_token`. (An earlier run in this feature's history hit the documented
intermittent flake `TestListingsReadContractEndToEnd` — docs/HARNESS-PROFILE.md lines 61-64,
unrelated catalog/listings test, not touched by this diff — as the sole failure; re-ran once
per that entry's instruction and got a clean pass, confirming the flake diagnosis.)

## sumSaleFee bugfix

Before (pre-existing, under-counting):

```go
func sumSaleFee(items []domain.MarketplaceOrderItem) *float64 {
	...
	for _, item := range items {
		...
		sum += *item.SaleFeeAmount   // missing quantity multiplier
	}
	...
}
```

After (`enrich_service.go`):

```go
// sumSaleFee sums the order's per-line SaleFeeAmount × Quantity into the order
// comissão. SaleFeeAmount is PER-UNIT (domain.MarketplaceOrderItem.SaleFeeAmount
// doc comment, corrected during F-02's review: verbatim from the provider's
// OrderDetailItem.SaleFeeUnit, never pre-multiplied by quantity) — this sum is
// exactly the "consumer needing a line/order total must multiply by Quantity
// itself" case that doc comment calls out. Before this fix, sumSaleFee summed
// SaleFeeAmount alone (no ×Quantity), under-counting commission — and
// therefore retorno_liquido/margem_pct — for every multi-quantity line. Any
// nil-fee line or an empty item list makes the whole sum honest-unknown
// (ADR-17: never a partial fabricated total).
func sumSaleFee(items []domain.MarketplaceOrderItem) *float64 {
	if len(items) == 0 {
		return nil
	}
	var sum float64
	for _, item := range items {
		if item.SaleFeeAmount == nil {
			return nil
		}
		sum += *item.SaleFeeAmount * float64(item.Quantity)
	}
	...
}
```

Unit test coverage in `enrich_service_test.go` exercises a quantity>1 line, asserting the
corrected (multiplied) total.

## Golden comprador_fiscal proof (both halves)

Split across the two packages that own the two unexported mapping functions being compared
(hexagonal-boundary-respecting — no illegal cross-package unexported access):

```
$ go test ./internal/modules/orders/adapters/postgres/... ./internal/modules/orders/transport/... -v -run 'Golden|MatchesOldLiveAdapter'
=== RUN   TestBuildBuyerFiscalInfo_MatchesOldLiveAdapterShapeForEquivalentData
--- PASS: TestBuildBuyerFiscalInfo_MatchesOldLiveAdapterShapeForEquivalentData (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/postgres	1.6s
=== RUN   TestMapCompradorFiscalGoldenOldAdapterShapeVsNewDBReaderShape
--- PASS: TestMapCompradorFiscalGoldenOldAdapterShapeVsNewDBReaderShape (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/orders/transport	2.1s
```

`postgres.TestBuildBuyerFiscalInfo_MatchesOldLiveAdapterShapeForEquivalentData`: hand-mirrors
the OLD live adapter's exact mapping (read from
`connectors/adapters/mercado_livre/buyer_fiscal_reader.go`'s `mapBuyerFiscalInfo`/
`mapBuyerFiscalAddress`) for a billing payload equivalent to a seeded row, asserts JSON-byte-
identical to `buildBuyerFiscalInfo`'s real output (StateName excluded from both sides by
construction — see spec.md decision 3; FetchedAt excluded — both sides zero-valued, the field
is verified dead downstream).

`transport.TestMapCompradorFiscalGoldenOldAdapterShapeVsNewDBReaderShape`: constructs one
`BuyerFiscalInfo` the OLD-adapter way and one the NEW-DB-reader way, asserts
`mapCompradorFiscal(...)` produces byte-identical JSON for both, AND asserts the resulting
JSON matches a hardcoded expected string
(`{"nome":"Fulano de Tal","doc_tipo":"CPF","doc_numero":"12345678900","endereco":{"logradouro":"Rua das Flores","numero":"100","cidade":"Sao Paulo","uf_codigo":"SP","cep":"01000-000","pais":"BR"}}`).

Together these two tests chain to prove the whole DB-row -> comprador_fiscal JSON path is
unchanged for equivalent data.

## Zero-live-call proof (EARS requirement)

`orders/adapters/postgres/zero_live_ml_call_test.go`, 3 tests:

```
$ go test ./internal/modules/orders/adapters/postgres/... -v -run 'ZeroLiveML|NeverImport|NilPoolDegrades|NeverProduce'
=== RUN   TestShipmentAndBuyerFiscalReadersNeverImportMercadoLivreConnector
--- PASS: TestShipmentAndBuyerFiscalReadersNeverImportMercadoLivreConnector (0.00s)
=== RUN   TestScanShipmentRowAndBuildBuyerFiscalInfoNeverProduceAnHTTPCall
--- PASS: TestScanShipmentRowAndBuildBuyerFiscalInfoNeverProduceAnHTTPCall (0.00s)
=== RUN   TestShipmentReaderGetShipment_NilPoolDegradesWithoutPanicOrNetworkCall
--- PASS: TestShipmentReaderGetShipment_NilPoolDegradesWithoutPanicOrNetworkCall (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/postgres	1.3s
```

- `TestShipmentAndBuyerFiscalReadersNeverImportMercadoLivreConnector`: parses
  `shipment_reader.go`/`buyer_fiscal_reader.go`'s import declarations (syntax-only,
  `parser.ImportsOnly`) and asserts neither imports
  `marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre`.
  Since Go cannot call a package's exported identifiers without importing it, this is a
  structural, compiler-enforced guarantee — not "does not currently call it" but "cannot call
  it without this test going red first."
- `TestScanShipmentRowAndBuildBuyerFiscalInfoNeverProduceAnHTTPCall`: calls the pure mapping
  functions directly, confirming they take no capability/client argument at all (their only
  inputs are already-scanned column values).
- `TestShipmentReaderGetShipment_NilPoolDegradesWithoutPanicOrNetworkCall`: an unconfigured
  reader (`nil` pool) returns an error immediately, never dialing out.

Combined with the already-passing `TestRealRepoInteractiveMLSites_MatchesAllowlist` (proves
root.go's interactive composition wiring no longer references
`newOrdersShipmentReaderAdapter`, and the must-fail proof above shows this specific claim is
falsifiable) this closes all three legs of the brief's zero-live-call requirement: (1)
composition wiring doesn't reference the retired constructor, (2) the new readers structurally
cannot import the ML connector, (3) their pure mapping logic has no path to a network call.

## Scope discipline (git status --porcelain, from repo root)

`git rev-parse --show-toplevel` confirmed:
`C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/sleepy-perlman-d0d325`

```
 M apps/server_core/internal/composition/orders_adapters.go
 M apps/server_core/internal/composition/orders_adapters_test.go
 M apps/server_core/internal/composition/root.go
 M apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go
 M apps/server_core/internal/modules/orders/application/enrich_service.go
 M apps/server_core/internal/modules/orders/application/enrich_service_test.go
 M apps/server_core/internal/modules/orders/ports/shipment_reader.go
 M apps/server_core/internal/modules/orders/transport/http_handler_test.go
 M apps/server_core/internal/platform/archguard/archguard_test.go
 M apps/server_core/internal/platform/archguard/testdata/three_sites/root.go
?? apps/server_core/internal/modules/orders/adapters/postgres/buyer_fiscal_reader.go
?? apps/server_core/internal/modules/orders/adapters/postgres/buyer_fiscal_reader_test.go
?? apps/server_core/internal/modules/orders/adapters/postgres/order_repo_bucket_counts_test.go
?? apps/server_core/internal/modules/orders/adapters/postgres/shipment_reader.go
?? apps/server_core/internal/modules/orders/adapters/postgres/shipment_reader_test.go
?? apps/server_core/internal/modules/orders/adapters/postgres/zero_live_ml_call_test.go
```

F-01/F-02 files (`ingest_service.go`, `order_ingest_errors.go`, `order_shipment.go`,
`orders_ingest_adapters.go`, etc.) do not appear — they were already committed
(`cf2c09e1`, prior `65e31c0`) before this feature started and were not touched. Every changed
or new file falls inside this feature's owned paths: `composition/` (3 files, all within the
brief's orders-region/allowlist/site-A scope), `orders/adapters/postgres/` (6 files: 2 new
reader implementations + 4 test files, all owned), `orders/application/` (2 files, the
`sumSaleFee` fix + its test), `orders/ports/` (1 file, the sentinel), `orders/transport/` (1
file, additive golden test only — `http_handler.go` itself untouched), `platform/archguard/`
(2 files, allowlist + its testdata fixture).

No `modules/connectors/`, `modules/listings/`, or `modules/pricing/` files touched. No
`ingest_repo_test.go`/`ingest_service.go` (F-02-owned) files touched.

**No commit, no `git add`, no push was run** — working tree left for the orchestrator's
independent review.
