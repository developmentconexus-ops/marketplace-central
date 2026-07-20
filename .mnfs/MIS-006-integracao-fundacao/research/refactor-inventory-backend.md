# Refactor Inventory — Backend (MIS-006-integracao-fundacao)

Base: `main` @138aac3d. Repo root: `apps/server_core`. Read-only audit — verified against
real files (paths verified 2026-07-20), extends `main-code-audit-d120` memory. Target-state
per `docs/design/SYSTEM-BLUEPRINT.md`, `STORAGE-SCHEMA.md`, `INTEGRATION-DATA-CONTRACT.md`,
`evidence/INTEGRATION-FINDINGS-D120.md`.

---

## 1. products_mirror canonical table

Confirmed absent: no `products_mirror` string anywhere in `apps/server_core` (grep 0 hits),
no migration creates it (migrations 0001–0074 enumerated, none named/contain it).

| path:line | current state | target | action | why |
|---|---|---|---|---|
| (none — table does not exist) | today "product truth" is read live per-request from 2 divergent readers below (no canonical row store) | new `products_mirror` (+ `products_mirror_stock_locations`) table, tenant+codigo_produto scoped, upsert-merge keep-absent (F-XLSX-1), `stale_since` col | CREATE | STORAGE-SCHEMA.md §products_mirror; is the join target for vínculos/pricing/mercado/pedidos |
| `internal/modules/erp_import/adapters/internalread/reader.go:84-107` (`Reader.snapshot`) | reads straight off `erp_import_products` via `LatestCompletedSnapshot` per-request — no materialized "current state" row, everything recomputed from the whole accepted-rows snapshot on every call (`FindProductsForLinking`, `GetSellableStock`, `GetCostAsOf`, `catalogPage` all re-scan `snapshot.AcceptedRows`) | becomes the xlsx-side writer INTO `products_mirror` (upsert-merge on import completion), reads move to a thin mirror-reader | REFACTOR | today's O(n) snapshot rescanning per query is the "read = compute" pattern the blueprint explicitly forbids (SYSTEM-BLUEPRINT.md line 8-9: enrichment at ingest, not read) |
| `internal/modules/internal_read/adapters/oracle/reader.go:22-116` (`Reader` / `FindProductsForLinking`) | live read-through query straight to Oracle (TGF* tables) per request, no cache table | becomes SankhyaAdapter's read-through path feeding `products_mirror` on a sync cadence (per D7, snapshot into mirror is optional but strongly implied for consumers to share one join target with xlsx) | REFACTOR | today two live-query adapters (xlsx-snapshot-rescan, oracle-live) implement the SAME `ports.Reader`/`FindProductsForLinking` interface shape already — good sign the port boundary exists informally, just not materialized into a shared table |
| `internal/modules/erp_import/domain/import.go:10-22` (`NormalizedRow`) | only 6 of 10 E2 fields modeled: Codprod, Descrprod, Custo, StockPhysical/Reserved, EAN, Refforn, Marca, NCM (no `local`/depósito, no `grupo` used downstream, no `preco venda`) | extend to full E2 10-field contract (add `local`, `preco`, formalize `grupo`) | REFACTOR | INTEGRATION-DATA-CONTRACT.md §E2 fields 4 (local), 7 (preco) "novo — não existe no modelo atual" |

---

## 2. ProductSourceAdapter port

Confirmed absent: grep `ProductSourceAdapter` across `apps/server_core` = 0 hits. No shared
port interface unifies xlsx and Sankhya today — they implement parallel, independently-typed
reader interfaces that happen to be structurally similar.

| path:line | current state | target | action | why |
|---|---|---|---|---|
| `internal/modules/internal_read/ports/*` (interface consumed by both `erp_import/adapters/internalread.Reader` and `internal_read/adapters/oracle.Reader`) | both readers independently satisfy `readports.Reader`/`FindProductsForLinking`/`GetSellableStock`/`GetCostAsOf`/`GetTaxInputs` — this IS today's de-facto port, just scoped to "read" not "source→canonical write" | formalize into `ProductSourceAdapter` port with an explicit `Sync(ctx) (ImportedRows, error)` or equivalent write-side method feeding `products_mirror`; keep the existing read-shape for consumers | REFACTOR (rename/extend, not rewrite) | the read-side abstraction already exists and works — SYSTEM-BLUEPRINT.md §2 says XlsxAdapter/SankhyaAdapter are "two adapters of one port," which is directionally what's here today, minus the write/sync half and minus the shared name |
| `internal/composition/root.go:772` (`func erpSource(getenv func(string) string) (string, error)`) | boot-time env resolution (`MC_ERP_SOURCE`) picks ONE of {oracle, xlsx} for the WHOLE backend at startup — wires either `internalreadoracle.Reader` or `erp_import/adapters/internalread.Reader` into `internalReadSvc`, never both simultaneously | replaced by per-tenant DB config selecting active adapter; both adapters can coexist per tenant | REMOVE (the boot-time selection) / REFACTOR (composition wiring) | INTEGRATION-DATA-CONTRACT.md §1b "Fonte ativa = configuração de tenant em banco... `MC_ERP_SOURCE` morre"; D120-POSTMORTEM I1 names this exact line as "root cause confusão" |

---

## 3. XlsxAdapter over the lenient xlsx parser

| path:line | current state | target | action | why |
|---|---|---|---|---|
| `internal/modules/erp_import/adapters/xlsx/parser.go` | lenient multi-sheet parser: PT aliases, preamble skip, union across sheets, honest-unknown for missing custo/estoque — functional, exercised live (#004-E 2012 produtos) | KEEP the parsing logic as-is; it is the "protocolo de import → snapshot" step (SYSTEM-BLUEPRINT.md §2 step 1-2) | KEEP | F-XLSX-1 + memory `lenient-xlsx-import-two-source`: parser is proven live, not the broken part |
| `internal/modules/erp_import/application/import_service.go` | orchestrates parse → validate → persist `erp_import_products` snapshot per protocol; nothing downstream of "snapshot completed" | ADD post-completion hook: upsert-merge into `products_mirror` + trigger link-candidate generation + enqueue market collection (INTEGRATION-DATA-CONTRACT.md §1c "4 elos quebrados") | REFACTOR | this is exactly the "import não é fim, é começo de cadeia" gap (contract §1b) — `import_service.go` today stops at "snapshot persisted," never calls product_links generation or market collection |
| `internal/modules/erp_import/adapters/internalread/reader.go:28-47` (`WithActiveSource`/`activeSourceFromContext`) | works correctly (fixed post D-119 IMPORT-FIX per memory `chip-import-fix-closed`) — ctx-carried source toggle between `xlsx` and `catalogo_cliente`, defaults to xlsx absent | becomes the per-tenant-config lookup instead of ctx-default; same mechanism, different config source | REFACTOR | already correct pattern (ctx-scoped, not global env) — just needs its selection source moved from request-param/default to tenant DB config |
| `internal/modules/erp_import/domain/import.go:64-67` (`ImportSource` = `xlsx` \| `catalogo_cliente`) | only 2 values, no `sankhya` | ADD `sankhya` (or equivalent) as a third `ImportSource`/active-source value once SankhyaAdapter exists | REFACTOR | needed once Sankhya becomes selectable per-tenant alongside xlsx |
| STORAGE-SCHEMA.md's old "rebuild mirror" language | — | upsert-merge keep-absent + `stale_since` (F-XLSX-1 corrected the design doc itself) | N/A (doc already flags this) | any XlsxAdapter CREATE work must implement merge, never wipe — `erp_import_products` is history, `products_mirror` is current-state and must never lose `product_links` by cascading a wipe |

---

## 4. SankhyaAdapter

There is no product-catalog Sankhya *sync* adapter (the kind that would feed `products_mirror`).
What exists is (a) a live Oracle read-through for product-linking/pricing/tax, and (b) an
unrelated "assisted Sankhya linkage" feature for orders (invoice/document reconciliation) —
different concern, do not conflate.

| path:line | current state | target | action | why |
|---|---|---|---|---|
| `internal/modules/internal_read/adapters/oracle/reader.go:1-32` (`Reader`, `NewReader(db queryer)`) | live-query Oracle reader: `FindProductsForLinking`, `GetSellableStock` (continues past line 120), `GetCostAsOf`, `GetTaxInputs`, backed by hand-built SQL against TGFPRO/TGFEST/etc (see `buildFindProductsQuery`) | becomes the query-execution core of `SankhyaAdapter`; add a sync/refresh entrypoint that writes into `products_mirror` on a cadence instead of (or in addition to) live read-through | REFACTOR (reuse, don't rewrite — this is mature working Oracle plumbing) | D1 ratified: "Sankhya read-through (modelo atual do adapter Oracle) — a fonte é viva, não duplicamos" — but products_mirror as shared join target still needs SOME write path per D7 |
| `internal/modules/internal_read/adapters/oracle/sankhya_linkage_reader.go`, `sankhya_linkage_config.go`, `internal_read/domain/sankhya_linkage.go`, `internal_read/application/sankhya_linkage_service.go` | separate feature: "assisted Sankhya linkage" for `orders` module (invoice reconciliation), NOT product catalog sync — wired only when `source == "oracle"` (`composition/root.go:520`) | out of scope for products_mirror; KEEP as-is, unrelated subsystem | KEEP | do not merge this with SankhyaAdapter — different bounded context (order↔fiscal linkage, not product master data) |
| `[TESTAR-SKW]` column mapping (custo: gerencial/reposição/média? estoque: reservado subtracted where? which NUTAB for preço?) | undecided — INTEGRATION-DATA-CONTRACT.md §E2 marks every Sankhya column hypothesis `[TESTAR-SKW]` | ratify with Oracle/Sankhya specialist session before coding sync | CREATE (design work, not code) | blocks any SankhyaAdapter CREATE work — memory `db-specialist-session` names the consult channel |

---

## 5. CHICKEN-EGG — Collect(codprod) via LinkedListings

Exact call chain today (confirmed at current line numbers, matches memory verdict):

```
market/application/collection_pipeline_service.go:134  Collect(ctx, codprod)
  → :140  s.identity.GetLocalIdentity(ctx, codprod)          -- local ERP identity (EAN etc)
  → :163  s.collector.SearchCatalogByEAN / GetCatalogProduct  -- resolves ML catalog identity
  → :184  if !run.stopped: s.collectCompetitiveSignals(...)
      → :269  listingIDs, err := s.identity.LinkedListings(ctx, codprod)
      → :275  for _, listingID := range listingIDs { ... GetOwnItemPricing / GetPriceToWin ... }
```

`collectCompetitiveSignals` (own-listing pricing signal, `market_price_snapshots` +
`competitive_signals`) is **keyed off `LinkedListings(codprod)`** — a codprod with zero
`product_links` rows produces `listingIDs == []`, the loop body never executes, and
`signals` stays empty (`:340-343` early-returns nil, appends nothing). Separately (not this
call chain but the actual Oportunidades gap): `collectCatalogEvidence` (`:193-244`, reached
via `SearchCatalogByEAN`→`ResolveIdentity`) does NOT require a link — it only needs EAN. So
`market_aggregates`/`competitor_offers` CAN populate for an unlinked product. The vínculo
gate is specifically on the **own-price competitive-signal** half, not the aggregate half —
but per memory + contract §1c, product never reaches Oportunidades anyway because nothing
ever calls `Collect(codprod)` for an unlinked/unimported product in the first place (no
trigger exists — see §6/S4 below); the *practical* chicken-egg is "nothing invokes Collect at
all," compounded by "even if invoked, competitive-signal half is link-gated."

**Minimal change to give every ERP product a market path:**
1. Trigger `Collect(codprod)` automatically post-import (xlsx) / post-sync (Sankhya) for every
   `products_mirror` row with a valid EAN — currently NOTHING calls this pipeline except a
   manual per-product UI action (`market/transport/collection_handler.go`, 1 codprod/POST,
   sync, no batch — D120-POSTMORTEM S5).
2. `collectCompetitiveSignals` (`:268-344`) should not be a hard prerequisite for Oportunidades
   visibility — Oportunidades needs `products_mirror ⋈ market_aggregates` (catalog evidence,
   EAN-only), not `product_links`. The design doc's own join (STORAGE-SCHEMA.md line 126:
   `products_mirror pm JOIN market_aggregates ma LEFT JOIN product_links`) already treats
   `product_links` as a LEFT JOIN (i.e., correctly NOT required for Oportunidades) — so no
   code fix is needed in `Collect` itself for that path; the fix is purely "make something
   call `Collect`/the catalog-evidence half in bulk," per item 1.

---

## 6. Auto-link chain (mirror × listings by EAN, auto-approve unique-exact-EAN)

| path:line | current state | target | action | why |
|---|---|---|---|---|
| `product_links/application/generation_service.go:60-78` (`GenerateLinkCandidates`) | generates candidates from listing snapshots × ERP matcher (`ProductMatcher.FindProductsForLinking`) — works, but nothing ever calls it automatically; only reachable via orphan endpoint | ADD trigger: call after backfill of listings AND after ERP import/sync completes (INTEGRATION-DATA-CONTRACT.md §E4 "Trigger automático... hoje endpoints órfãos") | REFACTOR (wiring only, logic works) | D120-POSTMORTEM S4: `POST /product-links/link-candidates/generations` (`product_links/transport/http_handler.go:90`) has zero FE/composition callers |
| `product_links/transport/http_handler.go:89-90` | `/product-links/listing-snapshots/imports` and `/product-links/link-candidates/generations` registered but curl-only, no caller | KEEP handlers; ADD internal (non-HTTP) invocation from import/backfill completion instead of/in addition to the HTTP surface | REFACTOR | orphan since D120-POSTMORTEM S4 audit; confirmed still registered, still uncalled |
| `product_links/application/resolution_service.go:129-149` (`ApproveCandidate`) | `if candidate.InternalProductID == nil { return ... "PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE" }` — approval is 100% manual (`ApproveCandidate` called only from operator-driven transport), no code path auto-approves | ADD auto-approve rule: EAN exact-match with `collisions[ean] == 1` (unique) → auto-call `ApproveCandidate`-equivalent transition at generation time, with audit trail | REFACTOR (extend generation_service, reuse resolution_service's transition/audit machinery) | INTEGRATION-DATA-CONTRACT.md D4 "RATIFICADO: auto-aprovar com audit trail"; operator already bulk-approved 31 manually in D-95 (memory `demo-data-provisioned`) as a stopgap |
| `erp_import/adapters/internalread/reader.go:344-355` (`validEANCounts`) + `356-366` (`identityQuality`) | uniqueness-collision detection ALREADY EXISTS for xlsx-side candidates (`QualityEANCollision` flag when `collisions[ean] >= 2`) | REUSE this exact signal as the auto-approve gate (`collisions[ean] == 1` ⇒ safe to auto-approve) | KEEP + REFACTOR (wire into generation/resolution) | the "is this EAN unique" computation the auto-approve policy needs is already implemented and tested — do not reimplement, just consume `QualityFlags` at generation time |
| `internal_read/adapters/oracle/reader.go:70-76` | same collision pattern (`activeEANCount.Int64 >= 2`) mirrored for the Oracle reader | same as above — SankhyaAdapter side is already collision-aware too | KEEP | confirms the signal exists symmetrically on both today's readers, good building block for the unified auto-approve rule regardless of which adapter is active |

---

## 7. sync_state + scheduler skeleton

Confirmed: no `sync_state` table (grep across migrations = 0; grep across code = only
comments/test doubles in unrelated modules, no real table/repo). Only 2 real tickers exist.

| path:line | current state | target | action | why |
|---|---|---|---|---|
| (none) | no `sync_state` table anywhere | new table: `tenant_id, installation_id, entity, cursor JSONB, last_full_sync_at, last_incremental_at, last_error JSONB, consecutive_failures` | CREATE | STORAGE-SCHEMA.md §sync_state; SYSTEM-BLUEPRINT.md Fase 0 names this the FIRST thing to build ("nada muda p/ usuário, tudo passa a ser rastreável") |
| `internal/composition/root.go:575` | `go integrationsbg.NewRefreshTicker(authSessionRepo, authFlowSvc, 5*time.Minute).Start(...)` — OAuth token refresh only | KEEP as-is (unrelated to product/listing/order sync) | KEEP | this ticker is fine, not in scope of the sync engine build |
| `internal/composition/root.go:577` | `go integrationsbg.NewFeeSyncScheduler(installationSvc, providerSvc, feeSyncSvc, 15*time.Minute).Start(...)` — runs every 15min but calls `FeeSyncer.Sync` which just re-seeds the SAME static 16%/22% constants (see §8) | becomes real live-tariff sweep once fee_sync.go is fixed; scheduler mechanism itself (interval-ticker pattern) is reusable as a TEMPLATE for the other 4 scheduled jobs (pedidos 5min, anúncios diário, mercado diário, tarifas semanal) | REFACTOR (the ticker harness is reusable; the job it calls is a stub — see §8) | only existing "scheduler" precedent in the codebase; SYSTEM-BLUEPRINT.md §4 wants 4 more of these keyed to `sync_state` cursors instead of blind interval reseeding |
| listings refresh: `internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:244-254` (sequential per-item fan-out, no goroutines) | manual-trigger-only, no ticker at all | become the body of the "anúncios: diário" scheduled job (SYSTEM-BLUEPRINT.md §4), parallelized (goroutines/errgroup) + bounded workers | REFACTOR | D120-POSTMORTEM S2: zero goroutine/errgroup usage confirmed repo-wide for these fan-outs |
| orders import: `internal/modules/orders/*` (no ticker; FE never calls `/orders/import` — P1 finding) | manual/out-of-band only | become "pedidos: 5min" scheduled job, cursor = `date_last_updated`, `sort=date_desc` | CREATE (job) + REFACTOR (existing import code becomes job body) | SYSTEM-BLUEPRINT.md §4/§5; P1/P2 in D120-POSTMORTEM |
| market collection: `market/transport/collection_handler.go:13-14` (comment "D-F4-p: no batch POST") | 1-codprod-per-POST, synchronous, no scheduler | become "mercado: diário" batch job iterating `products_mirror` (or linked set) | REFACTOR | D120-POSTMORTEM S5; also unblocks §5 chicken-egg fix (bulk-trigger mechanism) |

---

## 8. fee_sync.go:29 static hole + DTOs missing Raw json.RawMessage

| path:line | current state | target | action | why |
|---|---|---|---|---|
| `internal/modules/connectors/adapters/mercado_livre/fee_sync.go:27-49` (`FeeSyncer.Sync`) | hardcodes `CommissionPercent: 0.16` (Clássico) / `0.22` (Premium), `Source: "seeded"`, comment at line 13 admits "seeds... with static defaults only... not live provider-backed" | call `GET /sites/MLB/listing_prices?price=&category_id=&listing_type_id=` live, per category, write to new `ml_tariffs` history table (never UPDATE, close `valid_to` + insert) | REFACTOR (replace body; the `FeeScheduleSyncer` port/registration pattern at `:17-23` stays) | F-ADAPTER-1 hole #2; INTEGRATION-DATA-CONTRACT.md §E7 "sync agendado REAL... matar seed hardcoded 16%/22%" |
| `internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:971` (`mlItemResponse`), `:1022` (`mlOrderResponse`), `:1045,1051,1055,1059,1066` (order sub-DTOs) | typed-fields-only structs, no `Raw json.RawMessage` capturing original body; only `flexString` (tolerant string/number) is defensive, not the whole struct | add `Raw json.RawMessage` field to each top-level DTO (item, order, shipment), populate from the original response body, persist alongside mapped columns | REFACTOR (additive field + one assignment per unmarshal site, not a rewrite) | F-ADAPTER-1 hole #3: "DTOs rígidos SEM Raw... campo que o ML adiciona/renomeia = perda silenciosa" — direct answer to operator's stated fear |
| `capability_adapter.go:607-609` (429 recognized) / `:653` (`ErrCodeProviderRateLimited`, gives up) / price_writer `:79`, listing_writer `:105` | no retry, no backoff, no Retry-After honoring, `http.Client{Timeout: 15s}` (`:62`) flat timeout | add backoff/retry middleware honoring `Retry-After` before any batch/scheduled sync goes live | REFACTOR | F-ADAPTER-1 hole #1 — MUST land before Fase 2+ backfill (thousands of items) per the finding's own ordering note |
| everything else in `capability_adapter.go` (Bearer header `:716`, `context.Context` `:711`, multi-tenant resolver `:26,:721`, mapper layer `mapListing:753`/`mapOrder:798`, pagination `:233` offset + `:279` scroll_id, redaction `price_writer:19`, LimitReader `:744`, `x-format-new`/`flexString` `shipping_reader.go:56,155`, idempotency `:724`) | mature, correct, already better than the official golang-sdk (F-SDK-1/F-ADAPTER-1 verdict) | — | KEEP | explicit "REUSAR, não reescrever" finding — do not touch these in the refactor; only the 3 holes above are real |

---

## 9. Active-source config per tenant

| path:line | current state | target | action | why |
|---|---|---|---|---|
| `internal/composition/root.go:772` (`erpSource(getenv)`) + its call site wiring `internalReadSvc` to exactly one of oracle/xlsx reader at boot | boot-time, whole-backend, env-var (`MC_ERP_SOURCE`) — confirmed still the mechanism in main | replace with per-tenant DB-stored config, both adapters constructible simultaneously, selection resolved per-request/per-tenant | REFACTOR (kill the env branch, keep both reader constructions available in composition) | I1 in D120-POSTMORTEM: this is the FIRST of the "two confused source axes"; contract §1b "MC_ERP_SOURCE morre" is explicit |
| `erp_import/adapters/internalread/reader.go:24-70` (`activeSourceKey`, `WithActiveSource`, `ParseActiveSource`) | SECOND source axis — request/ctx-scoped toggle between `xlsx` and `catalogo_cliente` (both xlsx-shaped imports, different upload streams), defaults to xlsx absent selection; this is the one D-119 IMPORT-FIX fixed (`WithActiveSource` is now actually called, per memory `chip-import-fix-closed`) | this ctx mechanism is the right SHAPE for the target (request-scoped, explicit, no silent fallback) — extend its value domain to include a real `sankhya` source once that adapter exists, and source the values from tenant config instead of transport param default | KEEP shape, REFACTOR value domain | already correctly built request-scoped (not a global), already has `ErrUnknownActiveSource` fail-closed (`:53`) instead of silent fallback (:59-69) — good pattern to extend, not replace |
| `internal_read/adapters/cache/cache.go:306` area (`canonicalKey()`) | cache key composition — per D120-POSTMORTEM I3 "sem tenant nem source" at the time of that finding; CHIP-IMPORT-FIX (memory `chip-import-fix-closed`) says this was FIXED via `activeSourceKey(ctx)` fix + must-fail-proven test | verify still holds post-138aac3d (memory says merged); if a THIRD source (sankhya) is added, the same key composition must extend, not regress | REFACTOR (extend when sankhya lands) / VERIFY (confirm current fix still in place) | memory lesson explicitly: "ctx-carried scope MUST be in downstream cache key" — this is the exact class of bug a SankhyaAdapter addition could reintroduce if the cache key isn't extended alongside the source domain |
| `erp_import/domain/import.go:64-67` (`ImportSource` enum) | 2 values only, `xlsx` / `catalogo_cliente` — note neither of these is literally "sankhya"; `catalogo_cliente` is a SECOND xlsx-shaped source (prospect catalog), not the Oracle live source | Oracle/Sankhya live source is a THIRD, structurally different kind of source (not upload-based) — the per-tenant config model must represent "live read-through" vs "uploaded snapshot" as different source KINDS, not just add a 3rd enum value naively | CREATE (new source-kind modeling) | today's `ImportSource` enum models only upload-shaped sources; Sankhya is read-through, doesn't have a "protocol"/"snapshot" the same way — naive enum extension would misfit the domain model |
| `erp_import/adapters/postgres/query_repository.go:66` (`LatestCompletedSnapshot`) | reads `erp_import_protocols`/`erp_import_products` filtered by `source` column | works correctly for the 2 upload-shaped sources; N/A for Sankhya read-through (no snapshot to be "latest completed" of, per D1) | KEEP (for xlsx-shaped sources) | confirms scope: this repo method is xlsx/catalogo_cliente-only by design, Sankhya bypasses it entirely via the Oracle reader |

---

## Collision surfaces (parallel-planning)

Modules/tables/migrations more than one refactor stream will touch — sequence or lock these:

- **`internal/modules/erp_import/`** (domain, adapters/internalread, adapters/postgres, adapters/xlsx, application) — touched by §1 (mirror write), §3 (post-import trigger), §6 (link-gen trigger), §9 (source domain extension). Highest collision density in this inventory.
- **`internal/composition/root.go`** — touched by §2 (kill `erpSource`/`MC_ERP_SOURCE` boot branch), §7 (register new schedulers), §9 (per-tenant source wiring), and indirectly §8 (fee_sync ticker interval/body). Sequence: land `products_mirror` + tenant-config plumbing BEFORE touching the scheduler wiring in the same file to avoid repeated merge conflicts on this single large file.
- **New migrations**: `products_mirror` (+`_stock_locations`), `sync_state`, `ml_tariffs`, per-tenant active-source config table — all net-new, no existing migration to conflict with, but ALL 9xx-series migrations should be sequenced (mirror before sync_state per SYSTEM-BLUEPRINT.md §8 "Fase 0... Fase 1"; contract says Fase 0 = sync_state+scheduler skeleton+config THEN Fase 1 = mirror — note STORAGE-SCHEMA.md §"Ordem de migração" and SYSTEM-BLUEPRINT.md §8 actually disagree on which comes first, sync_state-first vs mirror-first — resolve before creating migration files).
- **`internal/modules/product_links/`** (application/generation_service.go, resolution_service.go, transport/http_handler.go) — touched by §6 (auto-approve, trigger wiring). Depends on §1 (mirror must exist as candidate source) and §3 (import completion hook is the trigger origin) landing first.
- **`internal/modules/market/`** (application/collection_pipeline_service.go, transport/collection_handler.go) — touched by §5 (bulk-trigger) and §7 (mercado scheduler). Depends on §1 (mirror as the Oportunidades join source) and §6 (auto-linked EANs feeding the competitive-signal half) landing first.
- **`internal/modules/connectors/adapters/mercado_livre/`** — touched by §8 (fee_sync.go, DTO Raw field, backoff) independent of everything else in this inventory; safe to parallelize against §1-§7,§9 (no shared table/file).
- **`internal/modules/internal_read/`** (ports, adapters/oracle, adapters/cache) — touched by §1 (SankhyaAdapter reuse), §4, §9 (cache key extension when 3rd source lands). Cache key change (§9) must land in the SAME change as any new source-kind addition, never split across two chips (memory lesson `chip-import-fix-closed`: this exact split caused the stale cross-source cache pollution bug once already).
