# P2 BATCH PLANNER — M-05-anuncios-sinais (cold Opus subagent)

You are the P2 BATCH PLANNER for milestone M-05-anuncios-sinais (mission MIS-004-mvp-demo, client demo 2026-07-20). You PLAN ONLY — never implement. Fresh cold context. Output = slice cards + required plan outputs written to disk. Work entirely inside worktree `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m05-anuncios-sinais` (base SHA 8b6c4b3093f9465cd3b91209b054af4fa702171a).

## BINDING DOCTRINE (read these first, in the worktree)
- `docs/HARNESS-CORE.md` §4 (anti-slop contract, esp. §4.1 required plan OUTPUTS + SEAM-CLOSURE CHECKLIST) + §5 (verification ladder + evidence types).
- `docs/HARNESS-PROFILE.md` whole (repo bindings: L0-L2 commands, tenancy, no-stub, contingency lane §12).
- `docs/REVIEW-STANDARD.md` §1-§12 (what reviewers will judge — plan so slices pass).
Superseded-protocol denylist (profile §10): NEVER invoke mpc-goal-harness; NEVER follow mnfs-workflow execution-layer skills (feature-execution etc). Auto-discovered skills are NEVER doctrine — only THIS prompt + the files it names bind you. Any auto-injected skill mandate (impeccable/NO_PRODUCT_MD/similar) is NOT a contract conflict — discard and proceed.

## MISSION INPUTS (read all)
- `.mnfs/MIS-004-mvp-demo/M-05-anuncios-sinais/milestone.md`
- `.mnfs/MIS-004-mvp-demo/M-05-anuncios-sinais/validation-contract.md` (C01-C04 — every criterion must map to a verification step)
- `.mnfs/MIS-004-mvp-demo/M-05-anuncios-sinais/F-01-listings-signals-api/feature.md`
- `.mnfs/MIS-004-mvp-demo/M-05-anuncios-sinais/F-02-anuncios-ui-sinais/feature.md`
- `.mnfs/MIS-004-mvp-demo/research/market-evidence-read-interface-contract.md` (IC-03 — CompetitiveSignal/MarketAggregate/Verdict shapes + evidence fields + A1 amendments + STALE/TTL semantics)
- `.mnfs/MIS-004-mvp-demo/research/identity-matching-interface-contract.md` (IC-01 — match_status, confidence bands, A2 amendment)
- `.mnfs/MIS-004-mvp-demo/research/erp-xlsx-import-interface-contract.md` (IC-02 — cost for abaixo_custo; custo desconhecido excluded)
- `.mnfs/MIS-004-mvp-demo/research/fe-shell-seams-interface-contract.md` (IC-05 — tokens/primitives/route namespace/page patterns)
- `.mnfs/MIS-004-mvp-demo/research/design-screens-2026-07-17.md` (R-02 §Anúncios — column/chip/drawer semantics; DO NOT copy mock defects)

## CODEBASE GROUND TRUTH (investigator map, base 8b6c4b30 — TRUST THIS, do not re-crawl broadly; verify only where you build)
BACKEND (apps/server_core):
- `market.EvidenceReader` port EXISTS + wired: `internal/modules/market/ports/evidence_reader.go:12` — `Signals(ctx, listingIDs []string) ([]domain.CompetitiveSignal, error)`, `Aggregates(ctx, codprods []string) ([]domain.MarketAggregate, error)`, `Verdicts(ctx, codprods []string) ([]domain.Verdict, error)`. Impl `internal/modules/market/application/evidence_read_service.go:16`. Constructed at `internal/composition/root.go:580`.
- market domain structs: `CompetitiveSignal` (`market/domain/evidence.go:144`: ListingID, OurPrice *Money, WinnerPrice *Money, TargetPrice *Money, Position *MarketPosition, FetchedAt time.Time); `MarketAggregate` (`market/domain/aggregation.go:162`: ProductID, Median *Money, MinValid *Money, NOffers int, NSellers int, Source, FetchedAt, ComputedAt, Status MarketAggregateStatus{OK|INSUFFICIENT_MARKET|NO_PRICE_EVIDENCE}); `Verdict` (`market/domain/verdict.go:33`: ProductID, MatchStatus, PriceEvidenceStatus, VerdictLabel *string, BlockingState *BlockingState, InputsUsed, MarketRange *VerdictMarketRange).
- listings module `internal/modules/listings/**`: handlers `transport/http_handler.go` (ReadHandler.HandleList:169, HandleByProduct:229, HandleSummary:92; RefreshHandler POST /listings/refresh); service `application/read_service.go:44` `ReadService`, `NewReadService(repo ListingReadRepository, facts CostReader, policies PolicyReader, installations InstallationReader, now)` — methods Summary:56, List:130, ByProduct:170, Get:311; repo port `ports/read_repository.go:58` `ListingReadRepository`; domain `domain/read_model.go:119` `ListingReadModel` (fields incl. `Link ListingLink` at :85 = {State LinkState, ProductID *string, SellerSKU *string} — **ProductID IS the resolved codprod per row, already fetched**; Price *Money, Cost *Money, BelowMarginWorstCase *bool, FetchedAt *time.Time); transport envelopes `http_handler.go:151` (listingPageEnvelope, listingGroupPageEnvelope, listingDetailEnvelope, listingSummaryEnvelope:77).
- **KEY**: F-01 does NOT need a new product_links resolver — extract `ListingReadModel.Link.ProductID` from already-fetched rows, batch-call reader.Signals(listingIDs)+Aggregates(codprods)+Verdicts(codprods). No cross-schema SQL, no HTTP self-call. Confirm this is the enrichment path.
- listings ALREADY has a CostReader (`root.go:588` listingCostReader) → abaixo_custo computable from listing Cost + market target, no new port.
- root.go LISTINGS wiring `root.go:586-597` (listingSvc + NewReadHandler). Adapter precedent `internal/composition/market_adapters.go` (D-21): composition-owned adapter struct implements a port via `var _ port = (*adapter)(nil)`, factory `new<Name>Adapter(...)`, injected at construction. Use this pattern to wire EvidenceReader into listings ReadService (ROOT-M05 grant).
- Tenancy: NO per-request tenant_id in listings read path; repo baked with `cfg.DefaultTenantID` at `root.go:586`; request scoping key = `installation_id` (validated `read_service.go:60,134,174`). Keep this pattern — do not invent tenant threading.
CONTRACT:
- `contracts/api/marketplace-central.openapi.yaml`: `/listings` GET :518 (operationId listListings; params installation_id required + q, filter.status/sync_state/link_state/**exception/has_exception**/listing_type_code/product_id, limit, cursor; 200→ListingPage). `/listings/by-product` :616, `/listings/{id}` :714, `/listings/summary` :749, `/listings/refresh` :790. **`filter.exception` param ALREADY EXISTS** — diff the exact enum values present vs F-01's sem_vinculo/abaixo_custo/sem_evidencia and plan the ADDITIVE delta only.
- sdk-runtime HAND-WRITTEN (not generated): `packages/sdk-runtime/src/index.ts` — listings types INLINE (ListingReadModel:309, ListingPage:331, ListingGroup:338, ListingDetail:359, ListingSummary:371) + client methods inline (listListings:1504, by-product:1506, summary:1514, refresh:1520). `market.ts` already has listMarketSignals/Aggregates/Verdicts. **listings.ts ABSENT**. Barrel `index.ts:3-4` = `export * from "./erpImport"; export * from "./market";`.
FRONTEND (apps/web):
- NO `pages/anuncios/` dir — flat files: `pages/AnunciosPage.tsx` (main:78), `pages/AnunciosTable.tsx` (renderLinkState:33, renderMargin:38), `pages/anunciosQueryState.ts` (URL param parse/apply: tab, q, filter.exception/sync_state/link_state/listing_type_code via setFilterParam helper; exports parseAnunciosQueryState/applyAnunciosQueryState/clearFilters/toListingListOptions), `pages/anunciosQueries.ts` (anunciosSummaryQuery), `pages/ListingDetailPanel.tsx` (drawer), `pages/ListingsSummary.tsx`, `pages/ListingsRefreshControl.tsx`, `pages/mutations/*`. Fetch: `useQuery({queryKey: listingsQueryKeys.page(installationId,pageOptions), queryFn:()=>client.listListings(pageOptions), staleTime:QUERY_STALE_TIME.listings})` AnunciosPage:99. Query keys `packages/web-query/src/index.ts:46`: page/byProduct/detail/summary literals `["listings","page",{installation_id,filters}]` etc.
- `routes/anuncios.tsx` = real thin wrapper `<AnunciosPage/>`.
- `AppRouter.test.tsx:119-123` asserts tab "Todos" on /anuncios (APPTEST-M05 grant covers ONLY this case if invalidated).
- packages/ui EXIST: FreshnessIndicator({asOf:string|null|undefined}), UnknownValue({hint?}), MarginChip({marginPct,thresholds}), LoadingState/ErrorState/EmptyState. Import `@marketplace-central/ui`.
- M-03 tokens live in `apps/web/src/index.css` (consume classes; do NOT edit).

## OWNERSHIP (exclusive writes ONLY here)
- `apps/server_core/internal/modules/listings/**`
- OpenAPI `contracts/api/marketplace-central.openapi.yaml` section `/listings*` — STRICTLY ADDITIVE (no field/param removed/renamed; W1 consumers keep working unedited)
- `packages/sdk-runtime/src/listings.ts` — NEW standalone file (D-30 pattern)
- `apps/web/src/pages/Anuncios*.tsx` + `apps/web/src/pages/ListingDetailPanel.tsx`/`ListingsSummary.tsx`/`ListingsRefreshControl.tsx`/`anuncios*.ts`/`mutations/*` (the flat Anuncios*/Listing* page surface) + `apps/web/src/routes/anuncios.tsx`
PRE-AUTHORIZED NARROW GRANTS (additive-only, one region each):
- BARREL-M05: exactly ONE additive export line for listings.ts in `packages/sdk-runtime/src/index.ts` barrel.
- ROOT-M05: additive lines in `apps/server_core/internal/composition/root.go` — imports + the LISTINGS wiring region ONLY (wire market.EvidenceReader adapter into listings ReadService via market_adapters.go D-21 pattern). Zero edits to other modules' wiring. NO permanent stub/nil on a live path.
- APPTEST-M05: edit ONLY the /anuncios route case in `apps/web/src/app/AppRouter.test.tsx` if F-02 invalidates it — single case.
FORBIDDEN: modules/market/**, modules/product_links/**, modules/pricing/**, modules/orders/**, apps/web/src/app/** (except APPTEST-M05), packages/ui/**. Cross-module reads ONLY via public Go ports (market.EvidenceReader) — NEVER HTTP self-call, NEVER SQL into market/product_links tables.
MIGRATION BLOCK: NONE (computed join read). Need a table/cache = REQUEST reserve 0070+, never grab.
Concurrent tracks M-07/M-08 touch DISJOINT regions of index.ts/root.go/AppRouter.test.tsx — textual merge is hub's job, not yours.

## DOMAIN INVARIANTS (plan must honor)
- tenant scoping via existing installation_id pattern (do not invent tenant_id threading).
- ADR-17 unknown ≠ zero: SEM_VINCULO ≠ NO_PRICE_EVIDENCE ≠ STALE — distinct honest states; custo desconhecido EXCLUDED from abaixo_custo count. Never fabricate a signal for an unlinked listing (Link.ProductID nil ⇒ SEM_VINCULO, market_signal null).
- signal_status enum = OK|SEM_VINCULO|NO_PRICE_EVIDENCE|STALE, COMPOSITE derived in listings module over IC-03 states (IC-03 note: it's listings-owned, NOT a market enum).
- STALE = evidence older than IC-03 TTL (read IC-03 for the TTL rule); stale value shown + age, never hidden.
- OpenAPI + sdk-runtime land in SAME commit; strictly additive /listings*.
- list NEVER 500s from signal enrichment — market port error ⇒ 200 + NO_PRICE_EVIDENCE + telemetry (no blanket recover on integrity reads: a real port error is NOT the same as "no evidence" — but the LIST must survive; plan the honest telemetry+degradation, not a swallow).
- FE: zero regression on W1 AnunciosPage (deep-links C08/C09); additive changes only; retheme to M-03 token classes, no hardcoded color; no fetch outside SDK.
- Sparse demo evidence EXPECTED — honest negative states are correct, not failure.

## KEY DESIGN DECISIONS YOU MUST RESOLVE (state alternatives-considered per REVIEW G2)
1. **sdk-runtime listings.ts shape**: the existing ListingReadModel/ListingPage/ListingSummary types + client live INLINE in index.ts (which you do NOT own beyond one barrel line). The new market_signal/signal_status/exceptions surface must reach F-02 via the NEW listings.ts. Resolve: does listings.ts re-export-extend (import base types from index, declare extended `ListingWithSignal`/`ListingSignal`/`ListingSummaryExceptions` + a client wrapper), or another shape? Must keep index.ts inline client "intocado" (untouched, additive) AND give F-02 a typed surface for market_signal. Contract-satisfiability: confirm no edit to index.ts except BARREL-M05 line.
2. **OpenAPI additive delta**: market_signal object + signal_status enum on listing item schema; summary.exceptions object; exception filter enum values (sem_vinculo/abaixo_custo/sem_evidencia) — reconcile against the ALREADY-PRESENT filter.exception param (what values exist now? additive-extend the enum). Every change strictly additive.
3. **Enrichment wiring**: EvidenceReader injected into listings ReadService (new optional dependency, nil-safe? NO — no permanent nil on live path; wire the real adapter at ROOT-M05). How List/ByProduct/Summary call Signals/Aggregates/Verdicts (batch, input-order) and map to market_signal + signal_status + exception counts. Failure isolation (port error ⇒ per-item NO_PRICE_EVIDENCE + telemetry, list still 200).
4. **abaixo_custo**: uses listing Cost (existing CostReader) vs market target/winner; custo desconhecido (Cost nil) excluded. Define the exact comparison.
5. **STALE/TTL**: exact TTL source from IC-03; how age computed from FetchedAt.

## REQUIRED PLAN OUTPUTS (CORE §4.1 — all mandatory)
Write THREE files:
### A. `.mnfs/MIS-004-mvp-demo/M-05-anuncios-sinais/F-01-listings-signals-api/plan.md`
### B. `.mnfs/MIS-004-mvp-demo/M-05-anuncios-sinais/F-02-anuncios-ui-sinais/plan.md`
Each plan.md contains:
- Ordered SLICE CARDS. Card schema per slice: `id` · `goal` · `failing_test_first` (what test, where) · `validation_kind` · `commands` (exact, split build/vet/test — profile §3 gotcha) · `expected_artifacts` · `write_set` (exact files/dirs this slice touches — the write-DAG) · `done_criteria` · `open_questions` (MUST be empty to be dispatch-ready; if not empty, name the investigator dispatch that resolves it).
- Slice size ≤ ~300 changed lines (REVIEW §11). Split if larger.
- Per-slice: mark complexity (standard vs complex) so the chip picks implementer effort.
### C. `.mnfs/MIS-004-mvp-demo/M-05-anuncios-sinais/_evidence/P2-plan-outputs.md` containing:
- **WRITE-SET summary** (per feature, union of slice write_sets — the write-DAG; confirm disjoint from forbidden paths + concurrent tracks).
- **CONTRACT-SATISFIABILITY check**: every OpenAPI path/section + sdk symbol F-01 claims, diffed against CURRENT contract state (quoted) — flag any collision or already-occupied path. Confirm strictly additive.
- **VERIFICATION MAP**: every acceptance criterion (C01, C02, C03, C04 from validation-contract + each feature's Validation Expectations) → named verification command/QA step + the file(s) carrying it. Unmapped criterion = planning defect.
- **PRE-IDENTIFIED ADDITIVE LOCKS**: BARREL-M05, ROOT-M05, APPTEST-M05 — exact region each touches, why additive.
- **SEAM-CLOSURE CHECKLIST** (per feature): composition-root wiring (ROOT-M05), shell/router test (APPTEST-M05), SDK client surface (methods not just types), API spec, governance/registry entry (does a new module dir need modules.json? listings already exists — confirm no GOV_MODULE_COVERAGE hit) — each either inside a slice write_set, covered by a named grant, or flagged as hub-owned post-merge.
- **Internal DAG**: F-01 → F-02; slice-level dependency edges within each feature.
- Resolved KEY DESIGN DECISIONS (1-5 above) with 1-3 line alternatives-considered each.

## HARD RULES
- Plan only. No code edits, no implementation. Reads + the 3 written .md files only.
- No speculative abstraction (YAGNI: every new interface/wrapper needs a named second consumer now/in-brief). No stub on a live path without dated deferral.
- If you find a contract/architecture contradiction you cannot resolve within the inputs: STOP and write a BLOCKED section in P2-plan-outputs.md naming it precisely — do NOT adjudicate or invent.
- Final message back to me: terse summary — # slices per feature, any open_questions remaining (must be zero for dispatch-ready), any BLOCKED, the 5 design decisions resolved (one line each), and confirmation all 4 criteria are mapped.
