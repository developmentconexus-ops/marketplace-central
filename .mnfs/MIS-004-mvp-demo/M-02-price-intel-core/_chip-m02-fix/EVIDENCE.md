# CHIP-M02-FIX — FINDING-M02-COLLECT-4XX evidence pack

- **Mission / milestone:** MIS-004-mvp-demo / M-02-price-intel-core
- **Finding / ledger:** FINDING-M02-COLLECT-4XX (D-86) — demo-critical
- **Branch:** `chip/m02-fix-siteid`
- **Base SHA:** `f6bcb55e5d91974d03d5b34fab28d55235007d61` (main, confirmed = merge-base)
- **Tip SHA:** `c03fce0020cddadeccbcf048a39e34078cc60dbd`
- **Owned seam:** `connectors/adapters/mercado_livre/catalog_identity_reader.go` + market
  collection error mapping / observability (market application + transport) +
  composition market identity ACL. Disjoint from T2-MIN (catalog_match_reader,
  pricing/**), BUYER (shipping/orders/billing), GRUPO-IMPORT (erp_import source).
- **Untouched (verified via `git diff --stat`):** catalog_match_reader.go, pricing/**,
  erp_import source, OpenAPI/SDK.

## Commits (per green slice)

| SHA | Slice | Item |
|-----|-------|------|
| `720fa5bb` | `fix(market): send site_id on ML catalog EAN search` | 1 (fix) + 2 (regression test) |
| `973d7e63` | `fix(market): map ERP not-found to sentinel so unknown codprod is 404` | 4 |
| `c03fce00` | `feat(market): populate causas[].detail with honest provider observability` | 3 |

## Item reconciliation

| Item | Scope | Root cause | Fix | Verified |
|------|-------|-----------|-----|----------|
| 1 | site_id on `/products/search` | `searchCatalogByEAN` omitted `site_id`; ML 4XXs without it (match_reader honors it, D-83 live) | `query.Set("site_id", a.siteID)` in `catalog_identity_reader.go` | green |
| 2 | regression test (5th field/param-absent instance) | — | `TestSearchCatalogByEANSendsSiteIDAndProductIdentifier` asserts site_id + product_identifier on the wire | RED without fix (`site_id = "", want MLB`) → GREEN |
| 3 | `causas[].detail` always null | domain carried no per-cause detail; `newCollectionResponse` hardcoded `Detail: nil` | `CollectionSummary.Causas` → `[]CollectionCauseDetail{Cause, Detail}`; `providerCauseDetail(operation, err)` emits honest MPC-native note (operation + HTTP status; no payload/URL/PII/token); handler forwards `c.Detail` | RED without fix (5xx/429/timeout detail = `<nil>`) → GREEN |
| 4 | 500 → 404 for unknown codprod | `GetLocalIdentity` forwarded ERP reader's `*erp_import.ERPProductNotFoundError` raw; handler's `errors.Is(err, ErrProductNotFound)` missed → 500 | `errors.As` maps it to `marketports.ErrProductNotFound` at the composition ACL (erp_import source untouched) | RED without fix (`ERP product 90001 not found`) → GREEN |
| 5 | OpenAPI + SDK | wire `detail` field already exists (`MarketCollectionCause.detail`, nullable string) | **no change** — additive population only | n/a |

## GETs / verification ladder (see `verify.log`)

- `go build ./...` (server_core) — exit 0
- `go vet` market + composition + connectors ML — exit 0
- `go test ./internal/modules/market/...` — ok (application, domain, transport)
- `go test ./internal/composition/` (identity/mapping/accountref) — PASS
- `go test ./internal/modules/connectors/adapters/mercado_livre/` — PASS
- `GOCACHE=<repo>/.gocache` (absolute), no GOFLAGS. No server booted, zero ML writes.

## Honest-error / observability notes

- `detail` strings are MPC-native operation labels + provider HTTP status class only
  (e.g. `catalog offers fetch: provider responded HTTP 503`,
  `own listing pricing: provider request timed out`). No raw provider payload, URL,
  PII, or token crosses the adapter boundary. Non-provider causes (`NO_IDENTITY`)
  keep `detail: null`.
- Live-drive of collection against the ML sandbox is the hub's post-merge step
  (mocks prove contract behavior here per ADR-17 / no-live-in-chip).
