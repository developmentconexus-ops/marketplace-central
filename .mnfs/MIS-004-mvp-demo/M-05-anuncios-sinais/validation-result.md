# M-05-anuncios-sinais — Validation Result

```yaml
id: M-05
type: milestone-validation-result
status: chip-closed (P6 dual-gate PASS; P7 live-drive = HUB post-merge per D-64)
parent: MIS-004
review_sha: bc3f843   # code frozen; HEAD 0e2d7d1 = same tree + P6 ledger
base_sha: 8b6c4b3093f9465cd3b91209b054af4fa702171a
gate: dual (cold Opus + adversarial-REFUTE, contingency lane D-23)
```

## Summary

F-01 (listings-signals-api) + F-02 (anuncios-ui-sinais) COMPLETE, all slices reviewed-green.
P5 verify ladder ALL GREEN. P6 dual gate PASS. C01–C04 code-verified at `bc3f843`; the
browser live-drive (§tela screenshots + C03 refresh-signal drive) is deferred to the HUB
post-merge P7 persona per HUB ledger D-64 (single dev-stack serializes wave B; M-07 merges
first). No stack booted by the chip (chips never boot the stack).

## §port — C01 (Listings enriquecidos via port interno)  → PASS (code-verified)

- Market data read ONLY through the Go port `ports.EvidenceReader` — `application/read_service.go:599-601`
  (`s.evidence.Signals/Aggregates/Verdicts`). Port is consumer-owned by listings
  (`ports/evidence.go:43-47`); the module never imports `modules/market`.
- Anti-corruption mapping at the composition root: `composition/market_adapters.go:479-536`
  (`listingsEvidenceAdapter`); `net/http` there is used only for status-code constants in
  error mapping, not a client call.
- ZERO internal self-HTTP call to `/market`: no `http.Client` anywhere in `modules/listings`.
- Additive: legacy `/listings` response preserved; new fields nullable. Unlinked listing ⇒
  `SEM_VINCULO`, no number, never zero (`read_service.go:580-583`).

## §evidencia — C02 (Evidência IC-03 em cada sinal)  → PASS (code-verified)

- `MarketSignal` carries `MatchStatus`, `NOffers`, `NSellers`, `Evidence{Source, FetchedAt,
  Freshness}` (`domain/signal.go:45-54`), populated at `read_service.go:640-648`.
- Negative states are named enums, never null/0: `DeriveSignalStatus` returns distinct
  `SEM_VINCULO` / `NO_PRICE_EVIDENCE` / `STALE` / `OK` (`signal.go:59-70`).
- Tests: `domain/signal_test.go:10-55`, `application/read_service_test.go:893-971`.

## §tela — C03 (AnunciosPage com sinais honestos)  → PASS (code-verified); live-drive = HUB P7

- PREÇO cell folds valor R$ + sign-colored %-chip; STALE → âmbar + FreshnessIndicator
  (`apps/web/src/pages/AnunciosTable.tsx:89-131`). SEM_VINCULO → value + link, no number;
  NO_PRICE_EVIDENCE → value only, no chip. NUNCA dupla sublinha (single-branch render).
- Drawer evidence panel renders source / fetched_at / n_offers / n_sellers + FreshnessIndicator,
  re-skinned to the ratified 2×2 card idiom + LINHA DO TEMPO timeline (real `detail.timeline`)
  (`apps/web/src/pages/ListingDetailPanel.tsx:129-191`, `198-286`). Margem honest boolean-derived
  (`:60-67`); QUAL renders `quality_score` as INT % directly, never ×100 (`:259`).
- Refresh control present + tokenized (`ListingsRefreshControl.tsx`).
- FE honest-state tests assert the ABSENCE of a number for SEM_VINCULO/NO_PRICE_EVIDENCE and
  92→"92%" not "9200%" (`AnunciosTable.test.tsx:73-75,142-224`, `ListingDetailPanel.test.tsx:163-192`).
- **DEFERRED to HUB post-merge P7 (D-64):** live-drive `/anuncios` on the fixture (wave-A closed
  + M-04 vínculos, ≥1 vinculado-com-evidência + ≥1 sem-evidência), light+dark screenshots, and
  the refresh-signal → renewed-age step per the C03 Drive block. Persona also confirms anuncios
  a11y is not WORSE than the app baseline (see FINDING-WCAG disposition).

## §seams — C04 (Ownership limpo)  → PASS

- Write-set = 28 files, ALL ⊆ owned surfaces: `modules/listings/**`, `contracts/api/*.openapi.yaml`
  (/listings additive), `packages/sdk-runtime/src/index.ts`(+test), Anuncios*/Listing*/ListingsSummary/
  ListingsRefreshControl(+tests)+`anunciosQueryState.ts`(+test), `composition/root.go`+`market_adapters.go`
  (ROOT-M05 grant), `.mnfs/**` evidence.
- ZERO new DB migration (`git diff --name-only base..bc3f843 -- migrations` = empty) — join/projection only.
- OpenAPI `/listings*` STRICTLY ADDITIVE (only appended enum values + new nullable props; nothing
  removed/renamed/retyped/required) and sdk-runtime types landed SAME commit (`5318b29b`; `07cd4688`
  SDK inline reframe per SDK-LISTINGS-M05).
- Grant usage: ROOT-M05 used (root.go swap to `NewReadServiceWithEvidence` + market_adapters append);
  SDK-LISTINGS-M05 used (inline into index.ts, not a standalone listings.ts); **APPTEST-M05 UNUSED**
  (AppRouter.test.tsx not in diff — write-set tighter than granted).
- Lanes L0–L2 GREEN (see §ladder).

## §ladder — P5 self-verification (clean state, chip-run)

- L0-GO `go build ./...` exit 0. L1-GO `go test ./internal/modules/listings/...` all `ok`, 0 FAIL.
- L1-WEB `npx vitest run` 266/266 (35 files). L0-WEB `tsc --noEmit` 288 == baseline, TS2688 == 0
  (0 new; FINDING-P2 false-alarm signature stable — report below).
- L2 zero-migration HELD; ownership 28⊆owned; OpenAPI+SDK coherent.

## §dual-gate — P6 (contingency lane D-23)

- **P6a cold Opus milestone gate → MILESTONE PASS.** C01-C04 all PASS (anchored); domain invariants
  1-6 HELD (ADR-17 three states; `GET /listings` no-500 degrade to NO_PRICE_EVIDENCE
  `read_service.go:602-614`; port-only/no self-HTTP; additive+zero-migration; zero ML writes; tenant
  scoping `l.tenant_id=$1 AND l.installation_id=$2` intact `repository.go:79,254,325`); design CONFORMS
  1:1. 0 blocking. 2 suggestion non-blockers (sync-pill map dup — rule-of-three not tripped; fan-out
  scaling caveat OOS). No abstentions.
- **P6b adversarial-REFUTE → all 7 surfaces SURVIVE, 0 REFUTED.** No fabricated/dishonest signal,
  no 500/panic path, QUAL scale correct, OpenAPI additive-only, ownership clean, tests genuine
  (real negative assertions, not theater), design 1:1 (no VS-MERCADO column, no dupla sublinha,
  no raw non-token literals).

## Dispositions (open items → closed/carried)

- **OPEN-Q1 — RESOLVED.** abaixo_custo target = `price_to_win` (IC-03 alias; `read_service.go:163,728,746-750`,
  `signal.go:547-552`). NO winner_price fallback — target-nil ⇒ `BelowCost` false (honest degrade). C03
  does not require a fallback; adding one would fabricate a target = ADR-17 violation. Port is contract-correct.
- **FINDING-WCAG — CARRIED (non-blocking, backlog).** Dark-theme `text-white` on solid `bg-accent`
  (~2.5:1, below WCAG AA 3:1) on 3 controls. Pre-existing app-wide idiom (also `vinculos/*`), NOT
  M-05-introduced. Hub D-64: non-blocking for M-05 → app-wide a11y backlog ticket. P7 persona confirms
  anuncios is not WORSE than baseline (not that it fixes the app-wide issue).
- **FINDING-P2 — report for profile §2 ratification.** Repo-wide `apps/web` raw `tsc --noEmit` carries
  PRE-EXISTING unrelated REDs (jest-dom matcher TS2339 across *.test.tsx; unrelated precos.tsx) beyond the
  documented TS2688 @types/node. Verify web slices via stash-diff / count-vs-baseline, not raw repo-wide tsc.
- **FINDING-G1 — RESOLVED by hub** (pre-existing base RCFG_UNDECLARED_READ from D-21; hub declared on main
  `a5d9478`; close-gate governance runs post-merge on fixed main).

## Handoff

- Current status: chip work COMPLETE through P6 (dual-gate PASS). Code frozen @ `bc3f843`.
- Next owner: HUB — wave-B merge (M-05 second, after M-07) → post-merge ladder → rebuild main stack →
  dispatch fresh cold QA persona to live-drive C01-C04 (§tela) → milestone flips closed on P7 PASS.
- Required files/evidence: this file + `_evidence/dispatch-ledger.md` (full P1-P6 record).
- Blockers or open decisions: none. Only the hub-owned post-merge P7 live-drive remains.
