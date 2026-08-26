# D6-R2 P8 — B20 Publicações Candidate

> **Status:** CANDIDATE / OPERATOR WALKTHROUGH REQUIRED — no LOCK claimed
> **Block:** B20 — Publicações core / R20–R21
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Candidate evidence:** `qualification/d6-r2-wireframes/b20-publications.html`
> **Wire prerequisite:** PR #70 read projections + PR #75 B10 closure — ACCEPTED / INTEGRATED
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Design adjudication

The bounded B20 P8 structural design was adjudicated in chat on 2026-08-26 and **approved by the operator** (`concordo`):

- **R20 `/publicacoes`** — Marketplace Listing collection under exact Organization + Marketplace Installation context;
- **R21** — one source-qualified Listing detail composed of owner-separated read-only regions;
- human presentation comes only from the typed `presentation` projection (`known + display_name` / `unknown` / `unavailable`) — never canonical refs/keys;
- honest population states: known-populated, known-empty, unknown and unavailable never collapse;
- lifecycle rendered as observed fact (`active / paused / closed / unknown`), with no synthetic status, score or per-field sufficiency;
- no mutation home in R20/R21; the ListingIntent editor remains B23 and continuation is navigation-only;
- no bulk-selection framework, saved-view platform, provider-direct write or screen-shaped API.

P6/P7 remain NOT TRIGGERED per P5: conventional collection/detail with authority separation already explicit.

## 2. Authority consumed

| Region | Operation | Semantic owner | Permission | Principal kinds |
| --- | --- | --- | --- | --- |
| R20 collection | `ListMarketplaceListings` | Offering | `offering.read` | H/A/S |
| R21 observation | `GetMarketplaceListing` | Offering | `offering.read` | H/A/S |
| R21 Disponibilidade | `GetSellableAvailability` | Availability | availability read | H/A/S |
| R21 Performance | `GetMarketplaceListingPerformance` | Performance | performance read | H/A/S |
| R21 Mercado | `GetCompetitivePosition` | Market | market read | H/A/S |
| R21 Economia | `GetExpectedEconomics` | Economics | economics read | H/A/S |

`MarketplaceListingListItem` / `MarketplaceListing` carry `listing` (canonical ref), typed `presentation`, `lifecycle`, observed fields/media/price, `observed_at` and provenance. The candidate binds `presentation.display_name` as the only human label and keeps `native listing key` in secondary technical disclosure.

## 3. Candidate interaction

- exact account selection with no hidden default (B00/B10 law);
- collection rows show display label, observed lifecycle tag and `observed_at`; a row with unknown presentation states so honestly instead of showing a key;
- list-level known-empty ≠ unknown ≠ unavailable, each with distinct copy; unavailable explicitly does not imply absence;
- detail composes Offering observation plus Availability/Performance/Market/Economics/Work regions, each carrying `data-region-owner` and its own knowledge state; degraded scenario proves unknown/unavailable regions do not block or fake data;
- displaying regions together grants no region mutation authority; every "Ver em…" is navigation to the owner surface;
- "Ir para intenções de anúncio" reveals only a navigation boundary: no draft is created by navigating;
- cursor-based collection; no bulk selection, no saved views.

## 4. Proof

- deterministic verifier `scripts/verify-d6-r-b20-publications-wireframe.mjs`: **8/8 negative controls**, wired into `npm run gate` for the B20 surface;
- browser-operated flow (Chromium): account gate → 4-row list with one honest unknown-presentation row → unavailable/known-empty list states distinct → detail with 5 owner regions → degraded scenario yields 4 unknown/unavailable regions → navigation-only boundary → technical disclosure;
- 390px viewport: no horizontal document overflow; mobile drawer navigation law active;
- browser console warnings/errors: **0**.

## 5. Revision 2 — operator walkthrough round 1 (`REVISE`, 2026-08-26)

The operator's first walkthrough returned `REVISE` with three findings; all three were adjudicated and the revised design approved in chat:

1. **Source-product identity in R20 (upstream repair executed).** The Listing reads carried no source-product identity, so an operator could not recognize their own SKU/reference. This reused the accepted PR #70 pattern: the OAD gained a required typed `source_product_link` on `MarketplaceListingListItem` and `MarketplaceListing` — `resolved` (SourceProductRef + `SourceProductPresentation` with display name/SKU/GTIN/source instance), `unresolved`, `unknown`, `unavailable`. Product surface remains **106/31/H-A-S**; no operation, path or Permission changed; the repair is protected by `verify-human-operable-read-projection.mjs`. The R20 table now carries a `Seu produto` column (name + SKU + source), with honest unresolved/unknown/unavailable states, and unresolved routes the human to Preparação.
2. **Richer observed detail (no wire change needed).** The wire already carried more than the first candidate rendered. The detail now shows publication context (category/product type), the full observed-field census including `not_applicable`, an observed-media gallery with per-item availability, and provenance with evidence custody.
3. **Faster in-context inspection (read-only peek).** Region continuations felt navigation-heavy. Following the common platform side-peek pattern, each owner region gained `Ver rápido`: a read-only slide-over summary (`data-peek-kind="read-only-summary"`, `data-peek-writes="none"`) that closes on Esc/backdrop and contains the `Abrir na área responsável` continuation. The peek renders only already-consulted owner reads, respects degraded region states, and grants no mutation authority — it is LOCAL_EPHEMERAL disclosure state only.

Revised proof: verifier **10/10 negative controls** (adds source-link collapse and peek-write widening controls); browser operation re-passed end to end (SKU filtering, subject link line, 9-row observed census, media gallery, peek open/Esc/degraded/mobile) with **0 console issues** and no 390px overflow.

## 6. Revision 3 — operator walkthrough round 2 (`REVISE`, 2026-08-26)

The operator liked the round-2 candidate but found the peek indirection unnecessary: the summary facts fit the region cards themselves. Applied as the simpler Global Maximum:

- each owner region card now renders its read-only facts inline (e.g. Disponibilidade: vendável agora / origem do estoque / última observação) — same already-consulted owner reads, zero extra clicks;
- the side-peek panel was removed entirely (redundant mechanism once facts are inline);
- degraded regions hide the facts table and state unknown/unavailable honestly — facts are never faked;
- the read-only law moved to the cards: `data-region-facts-kind="inline-read-only"`, with negative controls against widening cards into a write surface and against fake facts in degraded states.

Revised proof: verifier **11/11 negative controls**; browser operation re-passed (5 inline fact tables, degraded hides facts, Work region keeps facts) with **0 console issues** and no 390px overflow.

## 7. P8 operator gate

Current disposition:

```text
CANDIDATE
operator walkthrough required
LOCK / REVISE / UPSTREAM FINDING
```

The assistant, verifier and reviewer cannot set `LOCKED`. P9 must not run for B20 until the operator explicitly LOCKs this candidate.

## 8. Reopen triggers

Reopen only if operation of this candidate proves that:

- `presentation.display_name` is insufficient for real Listing recognition;
- an operator job requires a write that would need presentation rather than canonical identity;
- the owner-separated composition hides a materially required cross-owner fact;
- honest population/knowledge states cannot preserve the monitoring job;
- responsive or accessible operation cannot preserve the collection/detail job.
