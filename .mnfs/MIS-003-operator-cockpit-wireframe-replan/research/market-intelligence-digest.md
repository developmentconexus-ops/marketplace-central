# Research Note

```yaml
id: R-04
type: research
status: draft
owner: Codebase Investigator
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: support
```

## Topic

Decision-dense digest of the commercial-intelligence research corpus, grounding the market-data contract (IC-04) and the deferral of Mercado UI surfaces.

## Sources Checked

- Source: `docs/research/2026-07-12-mercado-livre-competitor-price-monitoring.md` (defines gates G1–G7 in its section 12), `docs/research/2026-07-12-mercado-livre-prelisting-batch50.md`, `docs/research/2026-07-13-mercado-livre-provider-homologation.md` (vendor-homologation ladder), superseded briefs M-07/M-13/M-14 and `mission.md` of MIS-001.
- Why it matters: market-data contracts must encode what is actually obtainable and forbidden; UI deferral decision rests on this.

## Findings

- **Post-listing (our items), official API proven live**: `GET /items/{id}` + `/sale_price` (10/10), `GET /items/{id}/price_to_win?version=v2` (10/10) → `status` (winning/competing/not_listed), `winner.item_id`, `winner.price`, `price_to_win` target, boosts. Winner ≠ lowest price (frete, Full, installments, reputation participate).
- **Blocked with current app credentials**: `/items/{competitor}` 403, `/sites/MLB/search` 403 — plausibly missing DevCenter functional permissions, not proven forbidden (G2).
- `GET /products/{catalog_id}/items` works live but is **shutdown-flagged** in en-US docs → transitory: only behind replaceable adapter/feature flag, never a foundation. `buy_box_winner` absent in 9/9 — not usable.
- **Price signal schema (never collapse into one field)**: `our_sale_price`, `winner_price`, `competitive_target`, `catalog_offer_price` (transitory), `catalog_min/median/max` (derived), `visual_search_card_price` (manual only). Forbidden: installment-as-price, freight-in-price, zero-for-unknown, silent stale, "lowest offer = winner", price_to_win as victory promise.
- **Pre-listing coverage (batch50, frozen)**: GTIN-exact catalog candidate 22/50; conservative adjudication → acceptable match 12/50 (24%), with price 11/50, ≥5 sellers 7/50 (14%); executable gate: zero false accepts, 16% auto-coverage. Two proven FALSE GTIN matches → semantic gate must beat GTIN. Fallbacks added ~0 official coverage; external engines 8% safe accepts.
- **Statistics rules**: per-seller dedup mandatory (lowest valid offer per seller); median-per-seller with n≥5 is the defensible stat; `INSUFFICIENT_MARKET` when <5 sellers; `NO_PRICE_EVIDENCE` when nothing validated; condition=new, BRL, positive price, known timestamp only.
- **Gates G1–G7 before any live collection**: G1 OAuth operational readiness (**currently FAILED**: `StartAuthorize` clobbers `connected → pending_connection`; canonical expiry = `integration_auth_sessions.access_token_expires_at`; refresh_token single-use, atomic persist, per-installation mutex); G2 permissions/durability audit; G3 written usage authorization (ToS clause prohibits robots/scraping; polling/retention/aggregation need written confirmation); G4 visual parity 30/30 cent-exact; G5 matching CI95 ≥0.99 per auto class; G6 quota/429 measurement (start 10 listings 1×/day); G7 security (rotate Oracle credential exposed in M-04 F-02 artifact; never log tokens/buyer PII).
- **Vendor ladder (frozen)**: 1 official API + own matcher, 2 DataForSEO (canary prepared, not executed), 3 JoomPulse, 4 Precifica, 5 Cargoos; scraping vendors rejected regardless of intermediary (ML clause 7.6). Minimum vendor response contract preserves `synthetic_id, requested_gtin, source, catalog_product_id, item_id, url, observed_gtin, title, brand, manufacturer_reference, variant_attributes, condition, price, original_price, currency, seller_key, captured_at, match_method, provider_confidence`.
- **Identity**: future `seller_sku` = `CODPROD`; live listing observed carrying legacy `REFFORN` as SELLER_SKU. Oracle coverage of 10,519 active products: 100% unique CODPROD; 68.10% valid+unique GTIN; 91 GTIN collisions.
- **Still-binding MIS-001 ADRs**: ADR-006 MVP evidence boundary; ADR-007 object-centered workspaces + deep links; ADR-008 product≠listing, link states unresolved|conflict|resolved|rejected; ADR-009 proportional security (writes off by default, no PII/secrets in evidence); ADR-010 mocks never claim integration; ADR-011 M-06 failed gate preserved at SHA `1eb8831f…`.
- **Carry-forwards from superseded briefs**: M-07 "no automatic price/listing writes; recommendations always show evidence + quality"; M-13 route contracts survive reload, composite listing segment, IC-003 quality vocabulary (current/stale/unknown/conflict), single-writer AppRouter/Layout/SDK-context seams; M-14 values-minimized evidence discipline, missing real sample = `externally_blocked`, never invented.

## Recommendation

MIS-003 pins the market-data contract (IC-04: schema, nullability, evidence semantics, collector port) and implements NO collector and NO market UI. Live collection is a successor-mission concern strictly sequenced behind G1–G7.

## Impact On Mission

Justifies deferral of 2c/Concorrência/@mercado UI; shapes IC-04; adds risk rows (G1 OAuth defect, `/products/{id}/items` shutdown) to mission risk register.

## Handoff

- Current status: complete.
- Next owner: Mission Strategist.
- Next action: none.
- Required files/evidence: the three research docs + MIS-001 artifacts cited above.
- Blockers or open decisions: business actions pending in corpus (ML support ticket, DataForSEO canary authorization) are explicitly NOT MIS-003 blockers.
