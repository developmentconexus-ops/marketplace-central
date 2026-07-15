# Interface Contract

```yaml
id: IC-04
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

Market/competitor data contract: `market` module schema + collector port + read API. **This mission pins the contract only — no collector implementation, no market UI.** Consumers (tela Mercado 2c, aba Concorrência 2b, colunas @mercado 2d) arrive in a successor mission after gates G1–G7 pass (see [[market-intelligence-digest]] / R-04).

## Why This Contract Exists

Three future UI surfaces + the simulator would otherwise each invent price-signal semantics. The research corpus proved the failure modes (winner ≠ lowest price; GTIN lies; installment-as-price) — encoding them now as schema is cheap; retrofitting is rework across every consumer. Deferral decision: operator, P3 review 2026-07-14.

## Resources Or Entities

**MarketObservation** (post-listing; one per own listing per capture):

| Field | Type | Null | Notes |
| --- | --- | --- | --- |
| `listing_id` | string | no | IC-02 canonical id |
| `our_sale_price` | money | yes | from `/items/{id}/sale_price` |
| `competitive_status` | enum | yes | `winning \| competing \| not_listed` (raw provider value preserved) |
| `winner_price` | money | yes | `price_to_win.winner.price` |
| `competitive_target` | money | yes | `price_to_win.price_to_win` — a target, NEVER a victory promise |
| `catalog_offer_price` | money | yes | TRANSITORY source (`/products/{id}/items` is shutdown-flagged); only via replaceable adapter |
| `catalog_stats` | object | yes | `{min, p25, median, trimmed_mean, p75, max, n_offers, n_sellers}` — per-seller deduped; null unless `n_sellers ≥ 5` |
| `evidence_state` | enum | no | `observed \| insufficient_market \| no_price_evidence` |
| `captured_at` | string | yes | RFC3339; NOT NULL on stored observation rows (write path always has a capture time); null ONLY in synthetic `no_price_evidence` read items materialized for ids with no stored row |
| `source` | enum | yes | `official_api \| vendor \| manual`; NOT NULL on stored rows; null ONLY in synthetic `no_price_evidence` read items (no stored row → no provenance to claim) |

**MarketReference** (pre-listing; one per ERP product per capture): `product_id` (CODPROD), `catalog_product_id` (nullable), `match_state` (`accept | review | reject | no_candidate`), `match_method`, `catalog_stats` (same shape/rules), `evidence_state`, `captured_at`, `source`. Identity gate outcomes are OURS, not the vendor's (homologation protocol §4); `accept` requires the semantic gate, never GTIN alone.

`money` = `{amount: string decimal, currency: "BRL"}`.

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| `listMarketObservations` `GET /market/observations` | future Mercado UI / simulator | `installation_id`, `listing_ids[]` (max 200) | `{items[], as_of}` | latest observation per listing; missing → item with `evidence_state=no_price_evidence`; sort: input order |
| `listMarketReferences` `GET /market/references` | future oportunidades UI | `product_ids[]` (max 200) | `{items[], as_of}` | same semantics |

Both endpoints are REAL in MIS-003 (OpenAPI + SDK + handler) but serve only what storage contains — which is nothing until a collector exists. Empty is honest: `no_price_evidence`.

**CollectorPort** (Go, module `market/ports`): `CollectObservations(ctx, installationID, listingIDs) ([]MarketObservation, error)` + `CollectReferences(ctx, productIDs) ([]MarketReference, error)`. MIS-003 ships NO production adapter. A test-double adapter exists ONLY in `_test` packages (mocks never claim integration, ADR-010).

## Enums And Statuses

- `evidence_state`: `observed | insufficient_market | no_price_evidence`. UI copy (future): "sem evidência de preço" / "mercado insuficiente (<5 vendedores)".
- `source`: `official_api | vendor | manual`.
- Forbidden semantics (schema-enforced where possible): no zero-for-unknown; no installment-derived price; no freight folded into price; `catalog_stats` null when `n_sellers < 5`.

## Error Cases

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| >200 ids | 422 | `too_many_ids` | request is well-formed but exceeds the batch cap (consistent with IC-03 `selection_too_large` 422) |
| missing installation (observations) | 400 | `installation_required` | |
| unknown ids | 200 | — | returned as `no_price_evidence` items, not errors |

## Persistence Expectations

Tables `market_observations`, `market_references` (tenant-scoped, append-only per capture; read = latest per key). Migrations land in MIS-003; tables stay empty in production.

## Canonical Examples

Honest-empty response:

```json
{"items": [{"listing_id": "inst_1~MLB3456790~-", "evidence_state": "no_price_evidence", "our_sale_price": null, "competitive_status": null, "winner_price": null, "competitive_target": null, "catalog_offer_price": null, "catalog_stats": null, "captured_at": null, "source": null}], "as_of": "2026-07-14T18:00:00Z"}
```

## Database Shape

- `market_observations`: PK `(tenant_id, listing_id, captured_at)`; money as NUMERIC + currency; `evidence_state` check constraint.
- `market_references`: PK `(tenant_id, product_id, captured_at)`.

## Seed Data

- None in production ever (fabricated market facts forbidden). Test fixtures live in `_test` only.

## Timestamp And ID Semantics

- `captured_at` = provider capture time, RFC3339 UTC; `as_of` = serve time.
- `captured_at` is part of both PKs, therefore NOT NULL in storage; the nullable `captured_at` in the read shape exists only for synthetic `no_price_evidence` items (no stored row to cite).

## Compatibility Rules

- Vendor sources must map to the homologation minimum response contract (protocol §3) at the adapter; canonical shape here never grows vendor-specific fields.
- Future collector milestone extends by implementing CollectorPort + a scheduler; API and schema unchanged.

## Route Namespace

- Server: `/market` prefix, `market` module (M-06 owns).
- Client: none this mission.

## Transport And Integration

- Vite dev proxy: add `/market` (even unused, keeps convention).

## Must Preserve

- Never invent price; "no evidence" is a first-class value.
- The 6-signal separation (`our_sale_price` / `winner_price` / `competitive_target` / `catalog_offer_price` / `catalog_stats` / manual) — never collapsed.
- Live collection sequenced strictly behind G1–G7 (G1 currently FAILED — OAuth defect).

## Must Not Decide In Feature Execution

- Schema, evidence-state semantics, stats rules (n≥5, per-seller dedup), port signature.

## Validation Impact

M-06 criteria: OpenAPI+SDK parity, honest-empty responses proven, contract tests on evidence-state rules, zero production adapter present.
