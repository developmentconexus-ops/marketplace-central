# M-07-simulador — Validation Result

Milestone: M-07-simulador · Mission: MIS-004-mvp-demo · Chip: CHIP-M07 · branch `chip/m07-simulador`
Evidence tree HEAD at authoring: post P6-remediation S7–S9 + S10 adversarial-fix + S11 race-polish; web-vitest 188/188 + build GREEN. P6 dual-gate GREEN (see below).
Full dispatch/verification trail: `_evidence/dispatch-ledger.md` (D0–D21 + findings log).

Status legend: ✅ verified (evidence in-tree) · ⏳ Pending-P7 (live-drive on the hub dev-stack; chip never boots a server) · 📋 hub ruling pending.

Dual gate (P6) precedes the live-drive (contract §Evidence Requirements: "Dual gate antes do live-drive"). Round-1 P6 = FAIL (D17); remediated S7–S9; round-2 re-gate below.

---

## §golden — IC-04 decomposition golden vectors  ✅
- `modules/pricing/domain/decompose.go` + golden test: **12 golden vectors GREEN** (D4, commit `d288cd12`). Per-component breakdown comissão/taxa fixa/frete/imposto/DIFAL/tarifa Full/custo/retorno; decimal strings end-to-end; unknown components → null (ADR-17), never fabricated 0.
- Whole pricing domain tree GREEN; `go vet`+`fmt` clean.

## §difal — DIFAL seed + override, disclaimer, verify-at-execution  ✅
- Seed padrão 2026: 27 UFs seeded (migration 0057), hermetically applied (throwaway pg, D6 `63da61cf`): efetivo SP 6.0 / BA 13.5 / MA 16.0 / MG 6.0 = domain-exact; `tenant_id NOT NULL`; paired-override CHECK rejects lone override.
- Repo round-trip (docker-exec psql, D7 `6eb9e6f5`): DIFAL override 19.05@2dp, list ordered BA<SP, override clears.
- Disclaimer string EXACT: "seed padrão 2026 — não é orientação fiscal" — surfaced on the DIFAL drawer (`DifalDrawer.tsx`, asserted `PricingPage.test.tsx` "não é orientação fiscal"). `PricingDifalListResponse.disclaimer` (SDK).
- R-04 verify-at-execution: override persists only on confirmed divergence Δ>0,049pp, audited (application `CalcService`, D8); no silent apply.
- FE DIFAL-off honesty (S9, D20): when `profile.difal_enabled=false` the decomposition shows "Margem sem DIFAL — não use para decisão de venda interestadual." (`DecompositionPanel.tsx:57-61`; asserted `DecompositionPanel.test.tsx`).

## §solver — bidirectional margem-alvo → preço + UNREACHABLE_TARGET  ✅
- Backend `modules/pricing/domain/solve.go` (D5 `139231db`): 5 golden tests GREEN incl. exact 15,00%; high-segment discontinuity crossing; 60%@comissão16% ⇒ UNREACHABLE ceiling 58.75; custo-nil blocks; per-segment monotonic.
- IC-04 HTTP: `POST /pricing/solve` → `PricingSolveResponse{reached,preco,ceiling_pct,desconhecidos,blocking_state,code}` (D8/D9); unreachable = HTTP 200 + `UNREACHABLE_TARGET` (not an error).
- **FE (S7, D18 `d1ef8ce`):** `SolverPanel.tsx` drives `pricingSolveTarget` from the UI (margem-alvo input); `reached=true` ⇒ suggested price; `reached=false` ⇒ achievable `ceiling_pct` surfaced, **price node not rendered — no fabricated price (ADR-17)**; SEM_CUSTO ⇒ named blocking banner; error ⇒ ErrorState. 6/6 tests GREEN incl. "surfaces the achievable ceiling without fabricating a price".

## §ports — IC-04 ports single source of truth; scenarios CRUD  ✅
- `modules/pricing/ports/calc_ports.go` (`CalcEngine`, `DifalForUF`) — decimal, published/COMMITTED to hub at D4; single decompose formula (no 2nd decompose in modules/orders — grep-confirmed D4). M-08 F-01 consumes this port (gate intra-wave B).
- `CalcRepository` composes `DifalReader`; `CostReader` nil=absent (ADR-17); scenarios CRUD (`domain/scenario.go`, D7): newest-first ordering, upsert-no-dup, delete — hermetically validated (D7).
- **FE scenarios (S8, D19 `006aa0a`):** `ScenariosPanel.tsx` drives `listPricingScenarios`/`createPricingScenario`/`deletePricingScenario` from the UI; save snapshots the working sim payload, reload re-applies to page state, delete + list-invalidate. 6/6 tests GREEN.
- Contract-note (C07): the SDK pricing surface landed ADDITIVELY in `packages/sdk-runtime/src/index.ts` (PRICING region, SDK-PRICING-M07 grant), NOT a separate `sdk-runtime/src/pricing.ts` as the C07 wording sketched — a ratified deviation (D9); OpenAPI `/pricing/*` + SDK same commit `b41583b8`.

## §tela — /precos live-drive (preço→margem, margem→preço, destino≠SP, DIFAL off, sem custo)  ⏳ HUB POST-MERGE P7
- **Deferred to HUB post-merge P7 live-drive** (hub ledger D-63): NO worktree re-point — this chip carries migrations 0055-0059 + new decompose/solve endpoints not on main; re-pointing the shared dev DB to the worktree would drift it. Hub runs the fresh-QA-persona P7 on the proper main stack after merge (M-01/M-04 pattern). Chip never boots a server.
- Planned steps per contract C05 Drive: open `/precos?params=1` → assert "Parâmetros"; product w/ custo; preço 129,90; destino BA; assert "DIFAL" + "seed padrão 2026"; expect decomposição com DIFAL BA (efetivo interna−7%). Plus: margem→preço solver drive, destino≠SP, DIFAL-off recalc, produto sem custo ⇒ SEM_CUSTO (nunca 0). Screenshots light+dark.
- Static precondition met: all C05 UI surfaces present + wired (region-decomposicao/solver/comparacao/aplicar/cenarios), web-vitest 186/186 + build GREEN (D21).

## §mutacao — aplicar preço, teto previewed, zero ML write  ⏳ HUB POST-MERGE P7
- **Deferred to HUB post-merge P7 live-drive** (hub ledger D-63): live-drive on the main stack — trigger "aplicar preço" on a simulation, `SELECT` the created protocol, inspect adapter log for the window.
- Expected: `price_update` via /mutations with preview+protocol; final status `previewed`; UI offers NO approve (teto); dispatcher provider OFF; ZERO ML write requests in the log window.
- Static precondition met: `ApplyPriceAction` caps at `previewed` (createMutation→previewMutation, approve NEVER invoked; surface exposes no approve control — asserted `ApplyPrice.test.tsx`, D14). Demo listing 90001 ↔ MLB3758134295 (F-listing-1).

## §seams — migrations block + ownership diff  ✅
- Migrations `pricing_calc_profiles`/`pricing_difal_rates`/`pricing_scenarios` in block 0055–0058 (0059 reserved, no file); runner fixture 51→55; runner_test GREEN@55 (D6).
- Ownership diff: F-01 backend within `modules/pricing/**` + additive composition wiring (ROOT-M07) + OpenAPI `/pricing/*` + SDK index.ts PRICING region; F-02 within `apps/web/src/pages/precos/**` + `routes/precos.tsx` + `packages/feature-simulator/src/**` + AppRouter.test :93-97 (APPTEST-M07). **S7–S9 remediation diff = 8 files, ALL under `apps/web/src/pages/precos/**`** (`git diff --name-only 2ca0af6~3 2ca0af6` — zero forbidden-path touches). Reconciliation: `F-02-simulador-ui/changed-path-reconciliation.md`.

---

## P6 dual gate (round 2)
- Cold milestone-reviewer (opus, mnfs-workflow:milestone-reviewer) @ `5033305`: **PASS** on the code re-gate (D22) — all 3 D17 gaps MET; per-simulation DIFAL confirmed a feature-brief-vs-validation-contract divergence, NOT a milestone acceptance criterion (C05 requires only "toggle off recalcula sem DIFAL", satisfied by global toggle + warning); out of the F-02 FE grant → hub scope ruling, not chip code. 7/7 must-meet. Carries two non-code CLOSE preconditions: F-difal-1 hub ruling + P7 live-drive evidence (§tela/§mutacao).
- Adversarial-REFUTE (sonnet) @ `5033305`: **REFUTED** (D23) — found the S8 scenario-reload silent wrong-product desync + a reload integration-coverage gap. Everything else clean (ADR-17, wiring, DIFAL condition, forbidden-path, money precision). → fixed in **S10** (`113203d`, D24): unmatched selectedId ⇒ honest-empty not products[0]; absent-product reload shows a notice; integration tests added.
- Adversarial re-verify @ `113203d`: **COULD NOT REFUTE** (D25) — S10 fully resolves the desync; only a low-severity loading-race caveat (spurious notice over correct data while the catalog is still loading, not a refutation). Eliminated structurally in **S11** (derives `productMissing` reactively with a `products.length > 0` gate instead of storing notice state — net-simpler, race-free). web-vitest 188/188 + build GREEN. → **P6 dual-gate GREEN** (cold PASS + adversarial could-not-refute).

## CLOSE preconditions
1. ✅ **F-difal-1 — RESOLVED** (hub ledger D-60, ruling A). Global-only DIFAL toggle + D20 "Margem sem DIFAL — não use p/ decisão" warning ships for the demo; DIFAL is computed correctly from the profile (honest, ADR-17 clean). Per-simulation DIFAL toggle (option B) REJECTED pre-CLOSE — it would reopen the FROZEN/COMMITTED IC-04 PricingCalcInput + OpenAPI + backend override on a live M-08 consumer, two days from demo → moved to **POST-DEMO BACKLOG** (IC-04 `PricingCalcInput.difal_enabled/difal_destino_uf` + OpenAPI + backend). Not a correctness/honesty defect — a scoping refinement.
2. ⏳ **P7 live-drive** — §tela + §mutacao evidence. **HUB-COORDINATED** (hard rule: chips never boot the stack): REQUEST → hub re-points dev-stack FE to this worktree for browser QA (D-17 precedent) OR runs P7 live-drive hub-side post-merge (M-01/M-04 pattern). C05/C06 are live-drive criteria.
