# F-02-listings-read-api

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-01
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contract: IC-02 `../../research/listings-read-interface-contract.md` (operations, filter grammar, cursor, error matrix, canonical examples). ADR-17 (unknowns render null, never zero).

## Milestone

M-01 listings-read-spine. Depends on F-01 (table + module exist).

## Brief

Implement the four read operations over the F-01 table: `GET /listings`, `GET /listings/by-product`, `GET /listings/{listingId}`, `GET /listings/summary` — exactly per IC-02, including the flat filter grammar, cursor pagination, link-state join from product_links, and read-time `below_margin` computation. Update OpenAPI + `packages/sdk-runtime` in the same commit. (Dev-proxy row `/listings` is added by M-02 F-02 per IC-05 single-writer rule — NOT here; tests hit server_core directly.)

EARS:
- While listings exist, when `GET /listings?installation_id=X` is called, the API shall return cursor-paginated items sorted title ASC then listing_id ASC, each item shaped per IC-02 Canonical Examples (nullable facts as JSON null).
- While a filter key outside `{status, sync_state, link_state, exception, has_exception, listing_type_code, product_id}` or a non-enum value is supplied, when listing, the API shall return 400 `invalid_filter` naming the offending key.
- While a listing has no cost in ERP, when computing `below_margin`, the API shall return `below_margin: null` (unknown), never `false`.
- While listings share a product link, when `GET /listings/by-product` is called, the API shall group by product with cursor over products and place the synthetic null-product group last.
- While a listing_id is malformed or absent, when `GET /listings/{listingId}` is called, the API shall return 404 `listing_not_found` (composite id `installation~provider_listing_id~variation`, literal `-` for null variation).

## Inputs

- IC-02 (operations, DTOs, filter grammar, cursor rules, error matrix, seed set), F-01 module code, product_links read model (link_state source), existing cursor-pagination precedent in catalog module, `contracts/api/marketplace-central.openapi.yaml` + `packages/sdk-runtime` update pattern (`GOV_API_SDK_SPLIT`: same commit).

## Expected Output

- Four endpoints per IC-02, wired through transport layer.
- `q` free-text search (title + provider_listing_id + seller_sku, the latter joined from product_links at read) separate from `filter.*` — all three IC-02 fields.
- Summary endpoint returning IC-02 counter shape (unknown counters null, not 0).
- OpenAPI paths + schemas; sdk-runtime typed methods (`listListings`, `listListingsByProduct`, `getListing`, `getListingsSummary`, `refreshListings`) — same commit.
- Integration tests against IC-02 seed: pagination, each filter key, grouping order, error matrix rows.

## Constraints

- No new columns; link fields joined at read from product_links, never persisted in `listings`.
- `below_margin` computed at read from cost+price; null inputs → null output.
- Cursor opaque (base64 keyset), stable under concurrent refresh; `invalid_cursor` → 400.
- No UI changes. No write endpoints.
- `GOCACHE` absolute; governance lanes green.

## Negative Scenarios

- `filter.bogus=x` → 400 `invalid_filter`.
- `filter.status=banana` → 400 `invalid_filter`.
- Garbage cursor → 400 `invalid_cursor`.
- Unknown composite id → 404 `listing_not_found`.
- Missing `installation_id` on list/summary → 400 `installation_required`.

## Validation Expectations

- Integration transcript over seed MLBTEST0001..0006: page walk covering full set exactly once; `filter.exception=sync_error` returns only the seeded error row; by-product shows null-product group last.
- JSON proof: seeded no-cost listing returns `"below_margin": null` and `"cost": null`.
- Error-matrix table test output: all five negative rows asserted status+code.
- Diff proof: `contracts/api/marketplace-central.openapi.yaml` + sdk-runtime in the single feature commit (no `apps/web/vite.config.ts` change — proxy row owned by M-02 F-02).

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (after F-01 accepted).
- Next action: compile context pack; read IC-02 + F-01 module only.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
