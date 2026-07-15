# Interface Contract

```yaml
id: IC-02
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

Canonical marketplace-agnostic listing read model: new `listings` Go module (read-only) → HTTP transport → sdk-runtime → web UI. Serves 2a Anúncios, 2b aba "Anúncios vinculados", 1e counters, and the selection input of IC-03 mutations.

## Why This Contract Exists

No listing entity exists today (only thin snapshots inside product_links). Six screens plus the mutation envelope consume listing rows; without one contract, workers would invent incompatible identity, filter, and unknown semantics. Bulk selection (IC-03) references THIS contract's filter expression — decided here so M-01/M-02/M-03 cannot diverge (review finding 9a).

## Resources Or Entities

**Listing** (canonical, provider-agnostic):

| Field | Type | Null | Notes |
| --- | --- | --- | --- |
| `listing_id` | string | no | canonical composite `installation~provider_listing_id~variation`, literal `-` for null variation (carried from MIS-001 M-13) |
| `installation_id` | string | no | integration installation |
| `provider` | string | no | e.g. `mercado_livre`; from provider registry |
| `provider_listing_id` | string | no | e.g. `MLB3456790` |
| `title` | string | no | provider title |
| `listing_type` | object | yes | `{code, label}` provider modality, e.g. `{"code":"gold_pro","label":"Premium"}`; canonical field, provider-mapped at adapter |
| `status` | enum | no | `active \| paused \| closed \| unknown` |
| `link` | object | no | `{state, product_id, seller_sku}`; `state ∈ unresolved\|conflict\|resolved\|rejected` (ADR-008); `product_id` = CODPROD as string, null unless resolved |
| `price` | object | yes | `{amount: string decimal, currency: "BRL"}`; null = unknown, never 0 |
| `published_quantity` | integer | yes | null = unknown |
| `sync_state` | enum | no | see Enums |
| `sync_error` | object | yes | `{code, message_pt, message_provider}` when `sync_state=error` |
| `quality_score` | integer | yes | 0–100; null = not computed |
| `pending_issue` | object | yes | `{kind, message_pt}`; `kind ∈ sync_error\|stale\|unlinked\|below_margin\|attribute_required` |
| `sales_30d` | integer | yes | null = unknown (no provider metric yet) |
| `cost` | object | yes | `{amount: string decimal, currency: "BRL"}`; latest cost fact joined at read; null = unknown ("custo?") |
| `below_margin` | boolean | yes | computed at read (pricing policy min-margin vs cost); null when cost or policy unknown — never `false`-by-default |
| `fetched_at` | string | yes | RFC3339 UTC; provider-fact capture time; null = never fetched |

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| `listListings` `GET /listings` | page load, filter change | `installation_id` (required), `cursor`, `limit` (default 50, max 200), `filter` (see below), `q` | `ListingPage {items[], next_cursor, page_size, as_of}` | sort: `title ASC, listing_id ASC` (stable); same envelope shape as `/catalog/products` |
| `listListingsByProduct` `GET /listings/by-product` | agrupar-por-produto view | same as above | `ListingGroupPage {groups[], next_cursor, page_size, as_of}`; group = `{product_id, product_title, listing_count, group_state, listings[]}` + one synthetic group `product_id=null` ("sem produto") ordered last | cursor over groups; sort `product_title ASC`; `group_state ∈ ok\|attention\|error` = worst child |
| `getListing` `GET /listings/{listing_id}` | drawer / deep link | path id | `Listing` + `timeline[]` (last 10 sync events `{at, kind, message_pt}`, sorted `at DESC` — newest first) | 404 when unknown |
| `getListingsSummary` `GET /listings/summary` | faixa de exceções, 1e counters | `installation_id` | `{total, active, paused, exceptions: {sync_error, stale, unlinked, below_margin}, as_of}` | counters computed server-side, one query |
| `refreshListings` `POST /listings/refresh` | "Atualizar" / first connect | `{installation_id}` | `202 {operation_run_id}` | triggers ingestion via connectors `ListListings`; concurrent refresh for same installation → 409 |
| `getCategoryAttributes` `GET /listings/categories/{category_id}/attributes` | corrigir-atributo flow (M-06) | path `category_id` | `{category_id, attributes[]}`; attribute = `{id, name_pt, type ∈ enum\|number\|boolean\|string, required, allowed_values[], constraints}`; `attributes[]` preserves provider order (required first per ML payload — never re-sorted) | **row owned by M-06 F-01** (extension of this M-01-owned prefix, granted here); ML payload mapped at connectors adapter; cached via L2 class `category_meta` (24h) |

### Filter expression (shared with IC-03 bulk selection)

`filter` is a flat JSON object, URL-encoded as `filter.<key>=<value>` query params:
`{status?, sync_state?, link_state?, exception?, has_exception?, listing_type_code?, product_id?}` — values are single enum strings (no arrays, no wildcards, no operators this mission; extension = add keys, never change existing key semantics). `has_exception ∈ true|false` matches any/no `pending_issue` (powers the "Com pendência" tab). `q` (free text over title/provider_listing_id/seller_sku) is separate from `filter`.

**Hub-blessed clarification (Option 2 ruling).** Lists filtered by `exception=below_margin` or
`has_exception=true|false` use a bounded iterative keyset scan because the worst-case margin fact
depends on request-local Oracle reads. These filters may therefore return a short non-final page.
One request scans at most 50 candidate pages; when that cap is reached before `limit` matches are
found, the response contains the accumulated matches and `next_cursor` identifies the last scanned
row. Passing that cursor resumes without silently skipping candidates. The response shape is
unchanged. Lists without a below-margin-dependent predicate retain the single-query `limit+1` path.

**Grouping (`GET /listings/by-product`).** Below-margin-dependent filters use the same bounded
scan mapped to group keys: a group appears iff at least one child survives, and the cursor advances
over every scanned group key, including dropped zero-survivor groups. The cap is 50 group-key pages;
on a cap hit the cursor is the last scanned group key. `listing_count` is the surviving-child count
(equal to `len(listings)`); every emitted child is cost-evaluated, with no sampling. Below margin is
defined only for linked listings, so `exception=below_margin` excludes "sem produto",
`has_exception=true` includes it via `unlinked`, and `has_exception=false` excludes it. The response
shape is unchanged.

## Enums And Statuses

- `sync_state`: `synced | error | stale | queued | syncing | paused_sync`. pt-BR labels fixed: sincronizado / com erro / desatualizado / na fila / sincronizando / pausado. ("sem vínculo" is `link.state`, NOT a sync_state.)
- `exception` filter values: `sync_error | stale | unlinked | below_margin`.
- `status`: `active | paused | closed | unknown` — provider status mapped at adapter; unmappable → `unknown`, never guessed.

## Error Cases

See Error Matrix.

## Persistence Expectations

New table `listings` (module `listings`), tenant-scoped, upserted by ingestion keyed on `(tenant_id, installation_id, provider_listing_id, variation_id)`. Ingestion is full-page pull via connectors capability; rows absent from a completed pull are marked `status=closed` (never deleted). `below_margin` exception computed at read time from pricing policy min-margin vs latest cost fact (null cost → NOT below_margin; it is `unknown`, surfaced via product completeness, not a false alarm).

## Canonical Examples

Success — `GET /listings?installation_id=inst_1&filter.exception=sync_error&limit=2`:

```json
{
  "items": [
    {
      "listing_id": "inst_1~MLB3456790~-",
      "installation_id": "inst_1",
      "provider": "mercado_livre",
      "provider_listing_id": "MLB3456790",
      "title": "Eletrodo 6013 2,5mm Pacote 5kg",
      "listing_type": {"code": "gold_pro", "label": "Premium"},
      "status": "active",
      "link": {"state": "resolved", "product_id": "2210", "seller_sku": "2210"},
      "price": {"amount": "89.00", "currency": "BRL"},
      "published_quantity": 0,
      "sync_state": "error",
      "sync_error": {"code": "provider_validation", "message_pt": "A marca é obrigatória nesta categoria.", "message_provider": "Attribute [BRAND] is required"},
      "quality_score": 61,
      "pending_issue": {"kind": "attribute_required", "message_pt": "Marca obrigatória"},
      "sales_30d": null,
      "cost": null,
      "below_margin": null,
      "fetched_at": "2026-07-14T17:02:11Z"
    }
  ],
  "next_cursor": null,
  "page_size": 2,
  "as_of": "2026-07-14T17:05:00Z"
}
```

Rejection — `GET /listings` without `installation_id`:

```json
{"error": {"code": "installation_required", "message": "installation_id é obrigatório"}}
```

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| missing `installation_id` | 400 | `installation_required` | all list/summary/refresh ops |
| unknown installation | 404 | `installation_not_found` | |
| invalid cursor | 400 | `invalid_cursor` | same semantics as catalog page contract |
| invalid filter key/value | 400 | `invalid_filter` | unknown keys rejected, not ignored |
| unknown listing_id | 404 | `listing_not_found` | getListing |
| refresh already running | 409 | `refresh_in_progress` | body carries active `operation_run_id` |
| unknown category_id | 404 | `category_not_found` | getCategoryAttributes (M-06 row) |
| provider unreachable on attributes fetch (cache empty) | 502 | `provider_unavailable` | getCategoryAttributes; no stale-cache fabrication |
| provider unreachable during refresh | 202 then failed run | — | failure lands on the operation run, never invents rows |

## Database Shape

- Table: `listings` — `tenant_id`, `installation_id`, `provider`, `provider_listing_id`, `variation_id` (text, `'-'` sentinel), `title`, `listing_type_code`, `status`, `price_amount NUMERIC NULL`, `price_currency`, `published_quantity INT NULL`, `sync_state`, `sync_error JSONB NULL`, `quality_score INT NULL`, `sales_30d INT NULL`, `fetched_at TIMESTAMPTZ NULL`, `created_at`, `updated_at`. PK `(tenant_id, installation_id, provider_listing_id, variation_id)`.
- Link fields are NOT duplicated: resolved at read time by join to `product_links` (single source of link truth, ADR-008).
- `cost` and `below_margin` are NOT columns: joined/computed at read time (latest cost fact + pricing policy); null inputs → null outputs.
- Check constraints: `sync_state` and `status` enum checks; `quality_score BETWEEN 0 AND 100`.
- Timestamps: TIMESTAMPTZ UTC.

## Seed Data

- Integration tests seed 5 fixed listings covering: synced+resolved, error+attribute_required, stale, unlinked, paused; plus 1 closed. IDs `MLBTEST0001..0006` under `installation inst_test`.
- Reset: ephemeral-postgres lane per governance; no live-provider data in unit/integration lanes.

## Timestamp And ID Semantics

- `listing_id` composite uses `~` separators, `-` for null variation; opaque to clients (never parsed client-side).
- All timestamps RFC3339 UTC. `as_of` = read-model serve time; `fetched_at` = provider capture time. Distinct, both surfaced.
- `updated_at` moves on every upsert; `fetched_at` only on successful provider fetch.

## Compatibility Rules

- New providers extend by new `provider` value + adapter mapping; canonical fields never gain provider-specific names.
- `filter` extends by new keys only. `Listing` extends by new nullable fields only.
- `sales_30d` stays null until a provider metrics source exists — UI must render null as "—", never 0.

## Route Namespace

- Server: `/listings` prefix mounted by the `listings` module (M-01 owns). Single granted extension: the `getCategoryAttributes` operation row is owned by M-06 F-01 (no other milestone adds routes under this prefix).
- Client pages consuming it: `/anuncios` (M-02 owns per IC-05).

## Transport And Integration

- Same-origin via Vite dev proxy: the `/listings` row in `apps/web/vite.config.ts` is added by M-02 F-02 (IC-05 owns the proxy list; single writer). M-01 tests hit server_core directly and do not touch the proxy file.
- Auth/session: unchanged platform behavior (no new cookies).
- CORS/credentials: no new CORS surface — apps/web and server_core are served same-origin in every environment (dev via Vite proxy; any future prod topology mirrors it with a same-origin reverse proxy). `sdk-runtime`/`createRefreshableFetch` keeps fetch credentials mode `same-origin`, unchanged, for all new prefixes (`/listings`, `/mutations`, `/market`, `/dashboard`, `/sync`). A cross-origin deployment would require a new ADR naming allowed origins — out of scope this mission.
- Freshness: responses carry `as_of`; `Cache-Control: no-cache` honored for forced refresh (same convention as catalog page reads).

## Must Preserve

- Unknown never becomes zero/default (price, quantity, quality, sales all nullable).
- Provider payloads stay at the ML adapter; transport carries canonical shape only.
- Cursor envelope shape identical to `/catalog/products` (`items/next_cursor/page_size/as_of`).
- OpenAPI + sdk-runtime updated in the same commit as any change here.

## Must Not Decide In Feature Execution

- Identity format, enum values, filter grammar, sort orders, group envelope shape, error codes — all fixed here.

## Validation Impact

Criteria in M-01 validation contract reference this file; the seed set above is the fixture for endpoint proofs.
