# CHIP-M02-LIVE-2 — FINDING-M02-LIVE-2 evidence pack

- **Mission / milestone:** MIS-004-mvp-demo / M-02-price-intel-core
- **Finding / ledger:** FINDING-M02-LIVE-2 (D-86) — demo-critical (competitor-price path)
- **Branch:** `chip/m02-live2`
- **Base SHA:** `910d819688b5db64a870ccc26d464b14d32ffd0a` (main tip, hub-granted)
- **Tip SHA:** `7fde1998ebc4f6ae38fbcde8945a7df24d1947a8`
- **Owned seam:** `connectors/adapters/mercado_livre/{catalog_offers_reader,pricing_reader}.go`
  + the market-collection error/observability boundary in
  `composition/market_adapters.go` (the `marketPriceIntelCollectorAdapter`
  collector methods + ACL mapping). Disjoint from T2-MIN (catalog_match_reader),
  BUYER (shipping/orders/billing), GRUPO-IMPORT (erp_import source).
- **Untouched (verified via `git diff --stat`):** catalog_identity_reader.go,
  catalog_match_reader.go, root.go, OpenAPI/SDK, erp_import source. No shape
  change → no contract/SDK edit (OpenAPI intentionally untouched).
- **Textual-overlap note:** hub flagged possible overlap with CHIP-T2-MIN in
  `composition/market_adapters.go`/`root.go`. This chip touched only the
  collector methods + added helpers in `market_adapters.go`; `root.go` untouched.
  Hub resolves at merge.

## Commits (per green slice, red-then-green)

| SHA | Slice | Defect |
|-----|-------|--------|
| `6c444a5b` | `fix(market): send context=channel_marketplace on ML own-item sale_price` | D3 |
| `2cb1ba5d` | `fix(market): reserve ErrCatalogOffersUnavailable for the flag-off path only` | D1 |
| `e0eedac4` | `fix(market): resolve catalog offers via parent product then child leaves` | D2 |
| `7fde1998` | `feat(market): warn-once on provider-body faults at the collection boundary` | Observability |

## Defect reconciliation

| Defect | Root cause | Fix | Red → Green |
|--------|-----------|-----|-------------|
| **D1** false FLAG_DISABLED | `catalog_offers_reader.go` wrapped any non-401/404/429 provider error (and pagination-integrity fault) into `ErrCatalogOffersUnavailable`; composition maps that to `ErrCatalogOffersDisabled`, so causes read "capability disabled" while the flag was ON | provider fault → `mapPricingReaderError(err)` surfaces as itself; pagination-integrity → `ErrProviderUnavailable`; `ErrCatalogOffersUnavailable` now emitted only on the flag-off path (`capability_adapter.go:171`) | `TestListCatalogOffersMidPaginationFailure...` / `...EmptyPageBeforeTotal...` asserted `ErrProviderUnavailable` && NOT `ErrCatalogOffersUnavailable` (RED: got `catalog offers unavailable: ...`) → GREEN. `...FlagOffAvoidsTokenAndHTTP` stays `ErrCatalogOffersUnavailable`. |
| **D2** offers 4xx on parent | `MLB22624877` (codprod 90008) is a PARENT catalog product; `/products/{id}/items` exists only on LEAF children, so the parent's `/items` 4xxs by design → zero live evidence | `listCatalogOffers` GETs `/products/{id}` first; if `children_ids` non-empty, fan out to each child leaf's `/items` and aggregate (`listLeafCatalogOffers`); leaf keeps the single-id paging path | `TestListCatalogOffersFansOutParentToChildLeaves` (RED: `unexpected request = GET /products/MLB-PARENT/items`) → GREEN; existing leaf tests updated to serve the `/products/{id}` metadata GET |
| **D3** own listing pricing 400 | `pricing_reader.go` GET `/items/{id}/sale_price` omitted `?context=channel_marketplace`; ML 400s without it → own price silently failed | add `context=channel_marketplace` query param | `TestOwnItemPricing` asserts `context=channel_marketplace` on the wire (RED: `GET /items/MLB-PRICE/sale_price?` empty query) → GREEN |
| **Observability** | provider validation body behind a 400 died silently at the collection boundary | 5 collector methods route provider errors through `mapAndObserve`, which warns once per `operation+status` with status + error code + `sanitizeProviderBody` (whitespace-collapsed, bearer-token-redacted, rune-bounded 256); timeouts/bare sentinels carry no body → not logged | `TestObserveProviderFailureWarnsOnceSanitized` (dedup=1, token not leaked, distinct status warns again) + `TestObserveProviderFailureIgnoresBodilessErrors` (RED: compile-fail undefined method) → GREEN |

## ML API facts (official docs, FINDING-M02-LIVE-2)

- `GET /products/{id}/items` — offers only on LEAF products; PARENT (`children_ids`
  non-empty) `/items` 4xxs by design. `/sites/MLB/search?catalog_product_id` unsupported.
- `GET /items/{id}/sale_price?context=channel_marketplace` — `context` REQUIRED (400 without).
- `price_to_win` already uses `?version=v2` (query-param precedent).

## Verification ladder (see `verify.log`)

- `go build ./...` — exit 0
- `go vet` connectors + market + composition — exit 0
- `go test ./internal/modules/connectors/... ./internal/modules/market/... ./internal/composition/` — all ok
- `GOCACHE=<repo>/.gocache` (absolute), no GOFLAGS. Zero ML writes, no server boot, no `.env`.

## Honest-error / no-leak notes

- Provider body only reaches the log via `sanitizeProviderBody`: whitespace
  collapsed, `Bearer <token>` redacted, rune-bounded to 256. No raw token/PII crosses.
- Live-drive of the corrected collection path against the ML sandbox is the hub's
  post-merge step (mocks prove contract here per ADR-17 / no-live-in-chip).
