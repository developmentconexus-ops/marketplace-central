# CHIP-M07-XLSXCOST — Evidence Pack

**Chip:** CHIP-M07-XLSXCOST (M-07 correction, xlsx cost/price 500)
**Branch:** `chip/m07-xlsxcost` · **base:** `96879eed` (main at dispatch) · **final SHA:** `525b884`
**Lane:** Claude-only contingency (D-23, codex quota dead til 2026-07-25)

## Defect (verified by hub, re-confirmed here)
`MC_ERP_SOURCE=xlsx`: POST /pricing/decompose|solve → 500 PRICING_INTERNAL_ERROR for every
`product_id`. Chain: pricing `products.Exists` → catalog `ListProductsByIDs` →
`GetCanonicalProduct` → `products()` → `product()` called `GetCurrentPrice`, which the xlsx
erp_import reader hardcodes to `ReadError{unsupported_query}`. `product()` returned that fatally →
existence check 500'd an existing product.

## Fix (catalog adapter only)
`apps/server_core/internal/modules/catalog/adapters/internalread/reader.go`
- New helper `isDegradableReadError` (reader.go:188-196): true iff error code is
  `ReadErrorUnsupportedQuery` OR `ReadErrorSourceUnavailable` (via `readdomain.IsReadErrorCode`,
  mirrors `composition/market_adapters.go:64-69`).
- `product()` stock/price/cost exits (reader.go:~200-230): on a degradable read error, reset the
  fact to its zero-value struct (`readdomain.SellableStock{}` / `CurrentPrice{}` / `CostAsOf{}` —
  nil Amount/Quantity) so `fact()` emits `QualityMissing*` (honest "—"). Any OTHER error still
  propagates unchanged.

### Before / after (each of the 3 exits)
```
BEFORE:  if err != nil { return CanonicalProduct{}, err }            // fatal → 500
AFTER:   if err != nil {
             if !isDegradableReadError(err) { return CanonicalProduct{}, err }  // genuine infra err still 500s
             <fact> = readdomain.<ZeroStruct>{}                       // degrade → QualityMissing*
         }
```

Principle: ADR-17 (profile §7) — unknown ≠ error; unavailable/unsupported fact degrades to honest
missing, never aborts existence/enrichment. `unknown ≠ zero` preserved: degraded value is nil,
Quality=Unknown (domain constructor `NewNumericSourceFact` hard-rejects Quality=Unknown with a
non-nil value — no fabricated zero can leak).

**erp_import reader, pricing, composition, module boundaries: UNTOUCHED. FE (apps/web): UNTOUCHED.**
Working `/catalog/products` LIST path (`ListCatalogProductFacts`/`catalogPage`) never calls
per-product `GetCurrentPrice` — unaffected (verified by both gates).

## Tests (reader_test.go, all new, package `internalread`)
- `TestCanonicalProductDegradesUnavailableFactsToMissing` — subtests `unsupported_query` +
  `source_unavailable`: all 3 facts error degradable → product resolves, each fact missing, no error.
- `TestListProductsByIDsDegradesUnavailableFacts` — same through the by-ID entrypoint pricing uses.
- `TestCanonicalProductDegradesOnlyTheFailingFact` — price unsupported while stock/cost return
  real values → price missing, stock+cost survive as Current with real values (scoped-degrade proof;
  added per cold-Opus gate suggestion).
- `TestCanonicalProductPropagatesUnexpectedFactError` — subtests stock/price/cost: a non-classified
  error (`errors.New`) still propagates → regression / over-degrade guard.

### Failing-test-first evidence (`ran`)
Temporarily reverted the price exit to old fatal behavior → `TestCanonicalProductDegrades*` +
`TestListProductsByIDsDegrades*` FAILED with `got error: not available in this source mode` /
`xlsx snapshot does not contain current price` (the 500 cause). Restored fix → all green.

## Ladder (worktree, GOCACHE=absolute .gocache, GOMODCACHE warmed)
| Level | Command | Exit |
|---|---|---|
| L0 build | `go build ./...` | 0 |
| L0 vet | `go vet ./...` | 0 |
| L1 test | `go test ./...` (full unit sweep) | 0 |

Governance lane: NOT run — change is logic-only on an already-registered catalog adapter, zero
`modules.json`/registry/contract delta → `GOV_MODULE_COVERAGE` cannot fire; hub re-runs governance
on integrated master at acceptance (profile §2). Integration lane: not applicable (no DB shape /
migration / platform change).

## P6 dual gate (contingency lane: cold Opus + adversarial sonnet)
Both dispatched sync, bounded input (diff + pointed files), fixed SHA `525b884` predecessor
`8e86d80`; the scoped-degrade test was the only gate-driven change, re-verified green.

- **Cold Opus (general-purpose, model=opus):** **PASS.** Degrade resolves 500 without fabricating
  values; classification correct & tighter than blanket swallow; LIST path untouched; no AI-slop;
  regression guard present. Findings: 1 suggestion (partial-degrade coverage — **ADOPTED**, added
  `TestCanonicalProductDegradesOnlyTheFailingFact`), 1 question (source_unavailable outage semantics
  — confirm intended), 1 nit (spacing, matches file style).
- **Adversarial sonnet refuter (general-purpose, model=sonnet):** **CANNOT-REFUTE (=PASS).** Zero-value
  flow verified nil→Unknown, constructor blocks fabricated value; no other caller broken
  (`go test ./internal/modules/catalog/...` green); tests exercise real code not stubs. One PLAUSIBLE
  caveat (see reconciliation).

### Reconciliation
Both gates independently flagged ONE shared item: in **oracle** mode, `wrapOracleError`
(oracle/reader.go:529-534) classifies every non-ErrNoRows Oracle failure as `source_unavailable`, so
a genuine Oracle outage on an existing product now degrades to 200-with-missing instead of 500 — in
oracle mode too, not just xlsx.
**Ruling — not a defect, gates clear:** (a) the chip spec explicitly mandated `source_unavailable`
degrade on all three exits per ADR-17 ("a product that exists must resolve to a CanonicalProduct with
missing facts, never a 500"); (b) pre-existing precedented convention — identical classify-and-degrade
in `composition/market_adapters.go:65` and `product_links/application/generation_service.go:562`; this
diff extends an accepted convention to a new call site, invents no new risk; (c) fail-honest: the
degraded fact is Quality=Unknown / value=nil → UI shows "—", the honest signal, not a wrong number.
A total-outage-all-facts-missing product is a visible, honest state per ADR-17, not a masked error.
