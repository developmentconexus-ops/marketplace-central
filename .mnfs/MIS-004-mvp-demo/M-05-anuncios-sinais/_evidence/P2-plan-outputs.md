# P2 PLAN OUTPUTS — M-05-anuncios-sinais (CORE §4.1)

```yaml
milestone: M-05-anuncios-sinais
mission: MIS-004-mvp-demo
base_sha: 8b6c4b3093f9465cd3b91209b054af4fa702171a
planner: P2 BATCH PLANNER (cold Opus, plan-only)
features: F-01 (5 slices), F-02 (6 slices)
open_questions_total: 0
blocked: none
```

## 1. WRITE-SET SUMMARY (write-DAG; disjoint from forbidden + concurrent tracks)

### F-01 union write-set
- `contracts/api/marketplace-central.openapi.yaml` — `/listings*` section ONLY (additive).
- `packages/sdk-runtime/src/listings.ts` — NEW.
- `packages/sdk-runtime/src/listings.test.ts` — NEW.
- `packages/sdk-runtime/src/index.ts` — BARREL-M05 line ONLY.
- `internal/modules/listings/domain/{signal.go(NEW), read_model.go, filter.go}`.
- `internal/modules/listings/domain/signal_test.go` — NEW.
- `internal/modules/listings/ports/` — NEW listings-side evidence consumer port.
- `internal/modules/listings/application/{read_service.go, read_service_test.go}`.
- `internal/modules/listings/adapters/postgres/repository.go`.
- `internal/modules/listings/transport/http_handler.go`.
- `internal/composition/{root.go (ROOT-M05 region), market_adapters.go}`.

### F-02 union write-set
- `apps/web/src/pages/{AnunciosPage.tsx, AnunciosTable.tsx, ListingDetailPanel.tsx, ListingsSummary.tsx,
  ListingsRefreshControl.tsx, anunciosQueryState.ts, anunciosQueries.ts}` (+ their `*.test.tsx`).
- `apps/web/src/routes/anuncios.tsx` (if touched).
- `apps/web/src/app/AppRouter.test.tsx` — APPTEST-M05 `/anuncios` case ONLY (iff invalidated).

### Disjointness
- All paths ⊆ OWNERSHIP grant. NONE in FORBIDDEN (`modules/market/**`, `modules/product_links/**`,
  `modules/pricing/**`, `modules/orders/**`, `apps/web/src/app/**` except APPTEST-M05, `packages/ui/**`).
- Concurrent tracks M-07/M-08 touch DISJOINT regions of `index.ts` (their own barrel lines), `root.go` (their own
  wiring regions), `AppRouter.test.tsx` (their own route cases). Textual merge is the hub's job, not this chip's.
- ZERO migration (computed join read over `ListingReadModel.Link.ProductID`). No `db/migrations/**` in any write-set.

## 2. CONTRACT-SATISFIABILITY (diffed against current contract state)

Current state quoted from base 8b6c4b30:

- `filter.exception` enum — CURRENT `[sync_error, stale, unlinked, below_margin]`
  (openapi.yaml:556 `/listings`, :654 `/listings/by-product`). DELTA: append `sem_vinculo, abaixo_custo, sem_evidencia`.
  Additive (existing 4 values retained; no consumer sending the old values breaks). SATISFIABLE.
- `ListingReadModel` — `required: [listing_id, installation_id, provider, provider_listing_id, title, listing_type,
  status, link, price, published_quantity, sync_state, sync_error, quality_score, pending_issue, sales_30d, cost,
  below_margin_worst_case, icms_worst_case_by_uf, fetched_at]` (openapi.yaml:2897). DELTA: add `market_signal` (object,
  nullable) + `signal_status` (enum) as NEW OPTIONAL properties — NOT added to `required`. W1 clients that ignore
  unknown properties keep parsing. SATISFIABLE, strictly additive.
- `ListingSummaryExceptions` — `required: [sync_error, stale, unlinked, below_margin_worst_case, margin_unknown]`
  (openapi.yaml:3038). DELTA: add `sem_vinculo`/`abaixo_custo`/`sem_evidencia` integer counters as NEW OPTIONAL props —
  NOT added to `required`. SATISFIABLE, strictly additive.
- SDK `packages/sdk-runtime/src/index.ts` — listings types + client INLINE (`ListingReadModel:309`,
  `ListingSummaryExceptions:363`, methods `1503-1520`). NO edit except the single BARREL-M05 export line. New surface
  lives entirely in NEW `listings.ts` (imports base types, declares extended `ListingWithSignal` + `withListingSignals`
  wrapper). No identifier collision (barrel already `export * from "./market"`; listings.ts symbols are listings-prefixed
  / distinct). SATISFIABLE — index.ts inline client stays intocado.
- `market.EvidenceReader` port signature (`Signals/Aggregates/Verdicts`) is FROZEN and consumed READ-ONLY. No change.
- `market_signal` object shape claimed by F-01 = superset projection over IC-03 `CompetitiveSignal` + `MarketAggregate`
  + `Verdict` (all fields sourced, none invented): `status` (listings-derived composite), `position`, `price_to_win`
  (= IC-03 `target_price`), `delta_pct` (listings-computed our_price vs price_to_win), `match_status` (IC-01 via Verdict),
  `n_offers`/`n_sellers` (MarketAggregate), `evidence.{source, fetched_at, freshness}` (MarketAggregate.Source +
  CompetitiveSignal.FetchedAt + derived freshness). No field requires a market contract change.

No collision, no already-occupied path. All additive.

## 3. VERIFICATION MAP (every criterion -> command/QA + file)

| Criterion | Verification | Carrying file(s) |
| --- | --- | --- |
| **C01** listings enriched via internal port, no HTTP self-call, unlinked ⇒ reason not zero | F01-S3 read_service_test (fake EvidenceReader) + grep-zero HTTP-to-/market in listings + F01-S5 booted List shows `signal_status` | `listings/application/read_service_test.go`, `internal/composition/*`; QA §port |
| **C02** each signal carries source/fetched_at/n_offers/n_sellers/match_status; negative states named, never null/0 | F01-S3 cases (strong / NO_PRICE_EVIDENCE / INSUFFICIENT_MARKET) + F01-S4 summary transcript | `listings/application/read_service_test.go`, `listings/transport/http_handler.go`; QA §evidencia |
| **C03** AnunciosPage honest signals + FreshnessIndicator + refresh + light/dark | F02-S2 table test (4 states) + F02-S5 drawer test + F02-S3 chips + F02-S6 retheme; P7 QA live-drive `/anuncios` (light+dark, refresh renews age) | `pages/AnunciosTable.test.tsx`, `pages/ListingDetailPanel.test.tsx`, `pages/AnunciosPage.test.tsx`; QA §tela |
| **C04** ownership clean, zero migration, lanes green | P5 lane ladder L0-L2 + diff-vs-matrix (all write-sets in §1) + zero `db/migrations` diff | this doc §1 + §5; QA §seams |
| F-01 VE: GET /listings 4-status transcript + summary counts + W1-subset regression + `?exception=abaixo_custo` | F01-S3 + F01-S4 tests | `listings/application/read_service_test.go` |
| F-02 VE: 4-state table screenshot + chip⇒URL filter + C08/C09 regression + group toggle | F02-S2/S3/S4/S6 tests + P7 screenshots | `pages/Anuncios*.test.tsx`, `AppRouter.test.tsx` |

No unmapped criterion.

## 4. PRE-IDENTIFIED ADDITIVE LOCKS

- **BARREL-M05** — `packages/sdk-runtime/src/index.ts`: exactly ONE new line `export * from "./listings";` (mirrors
  existing `export * from "./market"` at :3-4). Additive: no existing export changed; index.ts inline client untouched.
  Region: barrel export block only.
- **ROOT-M05** — `internal/composition/root.go`: additive imports + LISTINGS wiring region (:586-597) — swap
  `NewReadService(...)` for `NewReadServiceWithEvidence(..., newListingsEvidenceAdapter(marketEvidenceSvc))`, plus new
  adapter in `market_adapters.go`. Additive: `marketEvidenceSvc` already constructed (:580); zero edits to any other
  module's wiring; no permanent nil/stub on the live path.
- **APPTEST-M05** — `apps/web/src/app/AppRouter.test.tsx`: the `/anuncios` route case ONLY (:119-123), edited ONLY if
  F-02 markup invalidates the "Todos" tab assertion. Single case; no other route touched.

## 5. SEAM-CLOSURE CHECKLIST

### F-01
- Composition-root wiring (ROOT-M05): F01-S5 write-set. COVERED.
- SDK client surface (methods, not just types): `withListingSignals` wrapper in listings.ts re-typing
  `listListings`/`listListingsByProduct`/`getListingsSummary`. F01-S1. COVERED.
- API spec: openapi.yaml `/listings*` additive delta. F01-S1. COVERED.
- Governance/registry: listings module ALREADY in `modules.json` — no new module dir, NO `GOV_MODULE_COVERAGE` hit;
  registry entry lands via chip merge (not pre-merge). CONFIRMED no action needed.
- Shell/router test: N/A for backend.

### F-02
- SDK client surface: consumed via `withListingSignals(useClient())` (F02-S1) — no new SDK surface authored here.
- Shell/router test (APPTEST-M05): F02-S6 guards `/anuncios` case. COVERED by named grant.
- No composition-root or API-spec seam (FE-only). Query keys reuse existing `listingsQueryKeys.*` (no new key needed;
  `summary`/`byProduct`/`page` already exist in web-query — hub/M-03-owned, consumed not edited).

## 6. INTERNAL DAG

```
F-01:  S1 (contract+SDK, independent, unblocks F-02)
       S2 (domain) -> S3 (enrichment) -> S4 (summary+filter) -> S5 (ROOT-M05 wiring)
F-02  (starts only after F-01 S1 merged — listings.ts + OpenAPI fields present):
       S1 (SDK wire + query-state) -> S2 (VS MERCADO) -> S3 (chips)
                                    \-> S4 (group toggle, after S2)
                                    \-> S5 (drawer, after S2)
       S6 (retheme + APPTEST) after S2-S5
Cross-feature edge: F-02.* depends on F-01.S1 (typed surface) and F-01.S3+ (runtime fields) for live QA.
```

## 7. RESOLVED KEY DESIGN DECISIONS (alternatives considered)

1. **sdk-runtime listings.ts shape** — RESOLVED: standalone `listings.ts` imports base types from `./index` and
   declares `ListingWithSignal extends ListingReadModel` + a `withListingSignals(client)` wrapper that re-types the base
   client's responses (no second fetch). *Alt considered:* duplicate the full listings types + a fresh fetch client in
   listings.ts (rejected — duplicate fetch + drift risk); *alt:* edit inline index.ts types (rejected — index.ts not
   owned beyond the barrel line). Chosen keeps index.ts intocado, DRY, one barrel line.
2. **OpenAPI additive delta** — RESOLVED: `market_signal`+`signal_status` as NEW optional props on `ListingReadModel`
   (not in `required`); 3 new optional counters on `ListingSummaryExceptions` (not in `required`); extend the
   already-present `filter.exception` enum from 4 to 7 values. *Alt:* a separate `/listings/signals` endpoint (rejected —
   feature.md mandates per-item enrichment on the existing list; extra endpoint = non-additive surface + N+1 for FE).
3. **Enrichment wiring** — RESOLVED: listings owns a small consumer port; a D-21 composition adapter bridges it to
   `market.EvidenceReader`; injected via additive `NewReadServiceWithEvidence` so the ~40 existing `NewReadService`
   callers (many out-of-ownership) compile unchanged; List/ByProduct/Get batch `Signals(listingIDs)` +
   `Aggregates/Verdicts(codprods)` input-order; port error ⇒ per-item `NO_PRICE_EVIDENCE` + telemetry, list still 200.
   *Alt:* change `NewReadService` signature (rejected — breaks ~40 out-of-ownership callsites); *alt:* nil-safe optional
   dep on live path (rejected — no permanent nil/stub on production; real adapter wired at ROOT-M05).
4. **abaixo_custo** — RESOLVED: `abaixo_custo` when listing `Cost` is known (CostReader, IC-02) AND market target
   (`price_to_win`, fallback `winner_price`) < `Cost`; `Cost == nil` (custo desconhecido) EXCLUDED from the counter
   (ADR-17: unknown ≠ violation). *Alt:* `our_price < cost` (rejected — that is a margin check, not "vs market"; feature
   brief specifies market target/winner comparison).
5. **STALE/TTL** — RESOLVED: `signal_status = STALE` when `now - signal.FetchedAt > signalStaleTTL`, a listings-owned
   const `= time.Hour` mirroring market `snapshotTTL` (`market/application/collection_pipeline_service.go:19`, which sets
   snapshot `ExpiresAt = fetchedAt + 1h` — the IC-03 EXPIRED horizon). Age is derived from `FetchedAt` because the frozen
   reader port exposes no `expires_at`; STALE shows the value + age, never hidden. *Alt:* import market's `snapshotTTL`
   const (rejected — unexported + cross-module dependency violation); *alt:* invent an arbitrary TTL (rejected — anchor to
   the market horizon for honesty). IC-03 line 29 confirms the composite is listings-owned, so defining the const here is
   in-scope, not a contract change.

## 8. BLOCKED

None. No contract/architecture contradiction found within the inputs. All 11 slices dispatch-ready (zero open_questions).
