# CHIP-M02-LIVE-3 — Evidence Pack (live-drive reprova of FINDING-M02-LIVE-2)

- **Branch:** `chip/m02-live3`  **base:** `d69c0ddbda09ab6665326b746de03e4116eb6018` (main)  **tip:** `6bef1d4df6968c12431618aa86cf734fa6d33d9c`
- **Ledger:** FINDING-M02-LIVE-2, D-86 (demo T-1, 2026-07-20). Market collection returned NO competitor price live; two live-contract defects survived the mock-based P6.
- **Constraint honored:** ZERO ML writes (all GET). No server boot / no `.env` read by the chip (hub owns the dev stack). `GOCACHE` absolute, no `GOFLAGS`. No push.
- **Commits:**
  - `25e146cd` feat(market): surface raw provider body + route on 4xx and decode-fail  *(instrument-first, per hub order)*
  - `75eee00a` fix(market): decode numeric seller_id + wrap-proof pricing failure log  *(DEFEITO A)*
  - `6bef1d4d` fix(market): normalize no-variation listing ref; stop fabricating provider 400  *(DEFEITO B)*

Hub-mandated corrective sequence was followed exactly: **instrument first → reproduce live → capture ground truth → fix by evidence (no guessing) → re-drive**. Hub owns the dev stack; the hub captured every live body and ran every re-drive.

---

## DEFEITO A — competitor offers never decoded (seller_id shape drift)

### Ground truth (hub live capture @25e146cd)
The instrumentation surfaced the raw provider route + body on the decode-fail path:
```
WARN market collection provider failure operation=ListCatalogOffers status=400
  error_code=CONNECTORS_PROVIDER_PAYLOAD_INVALID
  provider_body="GET /products/MLB22624877/items?offset=0 -> HTTP 200:
    {"paging":{"total":16,"offset":0,"limit":100},"results":[{"item_id":"MLB4735328201",
     "site_id":"MLB","seller_id":691607102,"accepts_mercadopago":true,...}]}"
```
Route hit `/products/MLB22624877/items` directly with **HTTP 200** and `paging.total=16` — 16 real competitor offers were waiting behind a failing decode.

### Root cause
`seller_id` arrives as a bare **NUMBER** (`691607102`), but `mlCatalogOffer.SellerID` was a `string`, so `json.Unmarshal` failed on the 2xx body and the whole page was discarded as `PROVIDER_PAYLOAD_INVALID`. (NOT a `children_ids` shape / leaf-fanout issue as first hypothesized — the external-researcher confirmed `children_ids` is `[]string`; the instrumentation proved the real cause.)

### Fix — file:line
- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/catalog_offers_reader.go`
  - `mlCatalogOffer.SellerID` `string` → **`flexString`** (the same tolerant decoder `shipping_reader.go` already uses for the x-format-new numeric shipment id; decodes number OR string to literal text — contract drift degrades the field, never the read; ADR-17).
  - offer construction: `SellerID: strings.TrimSpace(string(result.SellerID))`.

### Red→green
`TestListCatalogOffersDecodesNumericSellerID` (catalog_offers_reader_test.go): numeric `seller_id` body → RED on the old `string` field (`PROVIDER_PAYLOAD_INVALID`), GREEN on `flexString` (offer decodes, `SellerID == "691607102"`, price `169.99`).

### Live proof (hub re-drive @75eee00a → confirmed @6bef1d4d)
`INFO catalog offers read route=/products/{id}/items status=ok page_count=1 offer_count=16` — **16 competitor offers collected and aggregated.**

---

## DEFEITO B — "own listing pricing HTTP 400" was a LOCAL fabrication, not an ML 400

### Ground truth (hub inspection of the worktree @75eee00a)
The re-drive still failed own-pricing with "provider HTTP 400", yet **no WARN fired** from the wrap-proof reader logger — because `doJSON` never ran. The chain (hub-verified, postgres + source):
1. `product_links.provider_variation_id` = **empty string** for both resolved links (90008→MLB4735328201, 90001→MLB3758134295).
2. `LinkedListings` (composition/market_adapters.go) built `ListingID{VariationID: ""}` → `String()` = `"<inst>~MLB4735328201~"` — **empty third segment**.
3. `ParseListingID` (listings/domain/read_model.go) requires `len(parts)==3 && all non-empty` → **rejects its own `String()` output**. Round-trip `String()→Parse` is broken for a listing without variations.
4. `accountRefForListing` (market_adapters.go) treated the parse failure by **fabricating `ProviderStatusError{StatusCode: 400}`** → pipeline classified `PROVIDER_4XX` → "own listing pricing: provider responded HTTP 400". A local failure wearing an invented provider status = double ADR-17 violation.

`GetOwnItemPricing`/`GetPriceToWin` never reached the adapter, so `logProviderReadFailure` (correct, from @75eee00a) was unreachable — no WARN, consistent with everything observed. **`sale_price`/`?context=channel_marketplace` was never the problem** (Prices-API research not needed).

### Convention that dictated the fix form
The codebase already has one no-variation representation: `NoVariationID = "-"` (listings/domain/listing.go:9). `NewListingKey` normalizes empty→`"-"` (listing.go:76-78); the mapper and the postgres read path store/serialize `"-"` and bridge the product_links `""` representation via `CASE WHEN l.variation_id='-' THEN '' ELSE l.variation_id END`. `ParseListingID` requiring all segments non-empty is **intentional** — `"-"` is the sentinel, empty is malformed by design.

### Fix — file:line (two honest parts)
- **(i) demo-core:** `composition/market_adapters.go` — `LinkedListings` now calls **`listingVariationOrSentinel(link.ProviderVariationID)`** (empty/blank → `listingsdomain.NoVariationID`), matching the rest of the system. Own-item pricing now resolves for variation-less listings. One file, zero blast radius across the ~12 `ParseListingID` consumers.
- **(ii) ADR-17:** parse-fail no longer fabricates `ProviderStatusError{400}`.
  - `market/ports/price_intel_collector.go` — new `ErrInvalidListingRef` sentinel.
  - `composition/market_adapters.go` `accountRefForListing` — returns `fmt.Errorf("%w: %v", marketports.ErrInvalidListingRef, err)`.
  - `market/application/collection_pipeline_service.go` — new `CollectionCauseInvalidListingRef`; `classifyProviderFailure` maps the sentinel → local cause (non-fatal, `stop=false`); `providerCauseDetail` renders "…: unresolvable listing reference (local)".

### FINDING reported to hub (accepted)
Hub-diagnosed option (i) "relax `ParseListingID` to accept empty" would contradict `NewListingKey`/`mapper.go` (which deliberately convert `""`→`"-"`) and create a second no-variation encoding rippling across all consumers. **Normalizing at the leak (`LinkedListings`) is the convention-consistent form.** Hub accepted the divergence.

### Red→green
- `TestListingVariationOrSentinelRoundTripsThroughParse` (composition): empty/blank variation → id round-trips through `ParseListingID` (parsed `VariationID == NoVariationID`); a real variation passes through untouched.
- `TestClassifyInvalidListingRefIsLocalNotProvider4xx` (market/application): `ErrInvalidListingRef` (direct and wrapped) → `(INVALID_LISTING_REF, false)`, never `PROVIDER_4XX`.

### Live proof (hub re-drive @6bef1d4d)
- **90008 → COMPLETED** (first time). match ACCEPT, `price_evidence_status=OK`, `blocking_state=null`, causes EMPTY. `GET sale_price` fired for real and **succeeded — HTTP 200, 1.94s**, no fabricated 400.
- **90001 → PARTIAL / NO_CANDIDATE** — honest (no ML catalog match); single cause `PRICE_UNAVAILABLE` "price-to-win: provider omitted the target price" = honest absence, not an error. The 400 is gone.
- **No `GetOwnItemPricing`/`GetPriceToWin` WARN since restart** — no provider failure occurred; wrap-proof logger silent as designed.
- **DB proof (`market_competitive_signals`):**
  ```
  inst-…~MLB3758134295~- | our=19.90  | winner=null | target=null | 19:22:56
  inst-…~MLB4735328201~- | our=169.99 | winner=null | target=169  | 19:22:55
  ```
  `listing_id` carries the `-` NoVariationID sentinel; own prices are REAL provider values (169.99 vs price_to_win target 169 = a genuine competitive signal for the demo). Also populated: market_price_snapshots=14, validated_offers=24, aggregates=6.

---

## Verification ladder
Full output: `verify.log` (this dir). Summary — all GREEN, run from the worktree `apps/server_core` with `GOCACHE=<repo>/.gocache`:
- `go build ./...` → exit 0
- `go vet` (composition, market, listings, connectors) → exit 0
- `go test` (composition, market, listings, connectors suites) → all `ok`
- New guards `-v`: all three PASS.

## Follow-up seam (OUT of this chip's scope — hub to decide post-close)
Our own listing (MLB4735328201, seller 691607102) appears inside its OWN catalog offers (the 16 aggregated). Own-vs-competitor seller filtering in the market composition is a distinct seam; hub confirmed out of scope for this chip.

## Non-committed field notes (checkout-local)
`docker/dev/*.sh` carry CRLF in this worktree; `set -euo pipefail\r` breaks bash. Hub worked around it with a `tr -d '\r'` command override. NOT fixed here (checkout-local artifact, not a source defect).
