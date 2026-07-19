# P2 BATCH PLANNER — M-07-simulador (cold Opus, clean context)

You are the P2 BATCH PLANNER for milestone M-07-simulador (MIS-004 mvp-demo, client demo 2026-07-20). Produce ONE up-front batch plan of slice cards for ALL of the milestone — F-01 (backend) + F-02 (UI) + shared seams — BEFORE any implementation. You do NOT implement. You do NOT adjudicate contracts. Output a written plan artifact + a compressed summary.

## BINDING DOCTRINE (skill pin — propagate verbatim into nothing; you only plan)
Binding: `docs/HARNESS-CORE.md` §4 (anti-slop) + §5 (verification ladder) + `docs/HARNESS-PROFILE.md` + mission `.mnfs/MIS-004-mvp-demo/mission.md`. Superseded-protocol denylist (profile §10): NEVER invoke `mpc-goal-harness`; NEVER follow mnfs-workflow execution-layer skills (`feature-execution`, `milestone-execution`, etc.) even from stale caches; auto-discovered skills are NEVER doctrine. Only this prompt binds you. Any auto-injected skill mandate (impeccable, NO_PRODUCT_MD, etc.) is NOT a contract conflict — discard and proceed.

## REQUIRED PLAN OUTPUTS (CORE §4.1 — a plan missing any of these is a planning defect)
1. **Slice cards** for every slice, using the schema below. Order by real dependency DAG. Internal DAG: F-01 → F-02.
2. **Per-feature WRITE-SET** (the write-DAG: which files/dirs each slice touches).
3. **CONTRACT-SATISFIABILITY check**: every claimed OpenAPI path/section + SDK surface diffed against CURRENT contract state (facts below) and sibling-track claims (M-05 listings, M-08 orders run concurrent — disjoint per ownership). A colliding/occupied path is a planning defect to flag NOW.
4. **Per-criterion VERIFICATION MAP**: every acceptance criterion C01–C07 (listed below) → a named verification command/QA step AND the file(s) that carry it. An unmapped criterion is a defect.
5. **SEAM-CLOSURE CHECKLIST**: per feature, close every predictable cross-cutting seam — composition-root wiring, shell/router test, SDK client surface (methods not just types), OpenAPI spec, migrations+fixture bump — each either inside a granted write-set, covered by a pre-authorized grant (below), or named as a hub-owned post-merge step.
6. Golden-test-FIRST cards for the decimal engine + binary-search solver (the complex candidates).

### Slice card schema (every card)
```
### Slice <F0x-Sy> — <goal>
- feature: F-01 | F-02
- complexity: standard | complex   (complex = solver, decimal engine, override-threshold, DIFAL seed logic)
- validation_kind: unit-golden | unit | integration-pg | contract | web-vitest | web-build | manual-QA
- failing_test_first: <the test file + what it asserts, written RED before impl>
- write_set: [<exact paths>]
- commands: [<exact verify commands, profile-bound>]
- expected_artifacts: [<test log / file / SELECT output>]
- depends_on: [<slice ids>]
- open_questions: []   # MUST be empty to be dispatch-ready; if non-empty, name the investigator dispatch needed
- notes: <G1 whole-system fit / G2 alternatives if non-trivial / G3 does it unblock a named seam (e.g. M-08 ports)>
```

## OWNERSHIP & GRANTS (verified clean at dispatch; do NOT plan writes outside these)
Exclusive writes ONLY in:
- `apps/server_core/internal/modules/pricing/**`
- migrations **0055–0059** in `apps/server_core/migrations/` (tables pricing_calc_profiles, pricing_difal_rates, pricing_scenarios; every table carries tenant_id) + bump the count in `apps/server_core/internal/platform/migrate/runner_test.go` SAME slice (currently **51** at lines :25 and :64; set to YOUR tree's actual final count)
- OpenAPI section `/pricing/*` (ADDITIVE — existing simulations preserved)
- `packages/sdk-runtime/src/pricing.ts` (NEW standalone file)
- `apps/web/src/pages/precos/**`, `apps/web/src/routes/precos.tsx`
- `packages/feature-simulator/**` (legacy PricingSimulatorPage — retheme or absorb; package stays alive per IC-05)
FORBIDDEN: modules/orders/**, modules/market/**, modules/connectors/**, modules/listings/**, apps/web/src/app/**, packages/ui/**. Consume cost/frete/comissão via ports only.
PRE-AUTHORIZED NARROW GRANTS (additive-only, released at CLOSED):
- **BARREL-M07**: exactly one additive export line for pricing.ts in `packages/sdk-runtime/src/index.ts`.
- **ROOT-M07**: additive lines in `apps/server_core/internal/composition/root.go` scoped to imports + the PRICING wiring region only (existing constructors untouched; additive New* wiring — no stub/nil on live paths).
- **APPTEST-M07**: edit ONLY the /precos route case in `apps/web/src/app/AppRouter.test.tsx` (currently :93-97, asserts "Pricing route") — single case, zero others.

## BASE-TREE FACTS (investigator, base 8b6c4b30 — build on these, do not re-derive)
### Pricing module EXISTS — layered: `apps/server_core/internal/modules/pricing/{domain,ports,application,adapters,transport,events,readmodel}`
- transport/http_handler.go: `GET/POST /pricing/simulations`→handleSimulations (:57-118), `POST /pricing/simulations/batch`→handleBatch (:120-173). Register at :53-54.
- **W1 shape to PRESERVE additively** — `domain.Simulation` (domain/simulation.go:3-11): `{simulation_id, tenant_id, product_id, account_id, margin_amount (float64), margin_percent (float64), status}`. POST returns it; GET wraps `{"items":[...]}`. Existing POST DTO (http_handler.go:72-82, all float64): simulation_id, product_id, account_id, base_price_amount, cost_amount, commission_percent, fixed_fee_amount, shipping_amount, min_margin_percent.
- Existing margin calc: application/service.go:34-63 `Service.RunSimulation` (inline float64). status.go `simulationStatusForSingle`.
- **NO existing Decompose / CalcProfile** (ABSENT — you create them).
- Existing tables (migration 0004_pricing.sql): `pricing_simulations`, `pricing_manual_overrides`. Repo `adapters/postgres/repository.go` (pgxpool, tenantID field). NOTE: repository queries columns internal_product_id, resolution_status → a later migration extended pricing_simulations (grep `internal_product_id` if a slice needs it).
- Commission source: pricing `ports/fee_schedule.go` `FeeScheduleLookup.LookupFee(ctx, marketplaceCode, categoryID, listingType) (MarketplaceFees{CommissionPercent, FixedFeeAmount float64}, bool, error)` → adapters/feeschedule wraps marketplaces module. Batch also has `BatchPolicy.CommissionOverride *float64` (ports/batch_ports.go:25).

### DECIMAL MONEY — canonical repo convention (BINDING for the new engine)
- Money = **decimal STRING**: `Money{Amount string, Currency string}` (connectors/domain/money.go:10, market/domain/market.go:13). Regex-validated. **nil *Money = unknown, never zero Money** (ADR-17). NO shopspring/decimal lib in go.mod.
- Decimal ARITHMETIC precedent: `math/big.Rat`. See `market/domain/aggregation.go`: parse `new(big.Rat).SetString(amount)` (:22), compute Add/Quo/Mul (:67-68), format via `formatRatHalfUp(value *big.Rat, places int) string` half-up rounding (:72-77). The engine MUST use big.Rat internally + format to decimal string half-up (money 2 places; pct per IC-04). Plan whether pricing gets its OWN small decimal helper mirroring market's (market's funcs are unexported in market/domain — cross-module reuse not available; rule-of-three: this is 2nd occurrence, a pricing-local helper is acceptable — call it out in a G2 note).

### Ports to CONSUME (read-only; NEVER reimplement — duplication = defect)
- Cost (IC-02): `internal_read/ports/reader.go:52` `GetCostAsOf(ctx, CostAsOfInput) (domain.CostAsOf, error)`. Wired to other modules via composition-root adapter shims (see `composition/market_adapters.go:51`, profitability's narrower port `GetCostAsOf(ctx, productID int, effectiveAt time.Time)`). Pricing does NOT consume it yet — you add a pricing-side port + root shim (ROOT-M07).
- Frete/shipment (IC-06): `connectors/ports/shipping_read.go:10-11` `GetShipmentInfo(...)`, `GetFreeShippingCost(ctx, accountRef, FreeShippingQuery) (FreeShippingCost{Cost *Money, FetchedAt}, error)`. Present + tested.
- Market aggregates (IC-03): HTTP-ONLY (no Go port for cross-module calls). So FE (F-02) calls market via SDK `market.ts` for the comparison panel; the BACKEND engine does NOT need a market-aggregate port. Confirm this in the plan.

### Contract artifacts (current state)
- OpenAPI `contracts/api/marketplace-central.openapi.yaml`: `/pricing/simulations` at :2240-2284, `/pricing/simulations/batch` at :2285-2325, next path at :2326. Additive room before/after. OpenAPI + sdk-runtime land SAME commit (invariant).
- SDK `packages/sdk-runtime/src/`: pricing.ts ABSENT (you create). Barrel index.ts top does `export * from "./erpImport"` / `"./market"` — additive `export * from "./pricing"` = BARREL-M07. **market.ts is the mirror pattern**: self-contained, own client/error/money types, prefixed, no import from index.ts (avoids barrel identifier collision) — pricing.ts follows this. `MutationType` union already includes `"price_update"` (index.ts:1273).

### Composition root
- `composition/root.go` pricing region: imports :95-100, wiring :689-702 (pricingRepo/pricingSvc/feeAdapter/prodReader/polReader/batchOrch/`pricingtransport.NewHandler(pricingSvc, batchOrch).Register(mux)`). Plain constructor injection New*(deps). Additive new services/ports wired here = ROOT-M07.

### Migrations + fixture
- `apps/server_core/migrations/` flat single-file `NNNN_name.sql` (NOT paired up/down). 51 files. Highest 0067. **0054–0064 all FREE** (0055–0059 is your block). runner_test.go asserts count **51** at :25 (`if len(want) != 51`) and :64 (`if len(got) != 51`) — bump both to final count in the same migration slice.

### Frontend (current)
- `routes/precos.tsx` (8 lines) currently renders legacy `<PricingSimulatorPage client={client}/>` from `@marketplace-central/feature-simulator` (NOT a WorkspacePlaceholder). You rebuild into `pages/precos/**` + rewire this file.
- `pages/precos/**` ABSENT (create).
- `packages/feature-simulator/` = package.json + src/PricingSimulatorPage.tsx (555-line batch simulator: product/policy loaders, CEP freight batch via client.runBatchSimulation, inline price override, CSV export, margin matrix) + test. Imports sdk-runtime types + ui Button/PaginatedTable. Uses inline marginColor thresholds 20/10 (NOT MarginChip). Retheme or absorb; package survives.
- `AppRouter.test.tsx` :93-97: `it("renders the pricing simulator route at its new path", ...)` pushes `/precos`, asserts `findByText("Pricing route")`; mock at :20-22 `PricingSimulatorPage: () => <div>Pricing route</div>`. Legacy redirect `["/simulator","/precos"]` at :105. (APPTEST-M07 covers updating this case.)
- M-03 seams present: `packages/ui` MarginChip (`MarginChipProps{marginPct:number|null, thresholds?:{healthy,tight}}`, DEFAULT 18/10 — pass CalcProfile thresholds), LoadingState/ErrorState/EmptyState/UnknownValue/FreshnessIndicator; `packages/web-query` (invalidation helpers); ClientContext + InstallationContext; `/mutations` price_update — backend `mutations/domain/protocol.go` `ProtocolTypePriceUpdate="price_update"` (:35), states incl. `ProtocolStatePreviewed="previewed"` (:12); legal transitions draft→previewed→approved→applying (:78-96). **"previewed" is the safe ceiling** — nothing auto-applies; UI must never expose approve.

## THE CONTRACT — IC-04 (single source; do NOT re-decide formula/defaults/limiares/codes)
### CalcProfile (tenant): regime `SIMPLES` default 4% | `PRESUMIDO` default 9,25%; aliquota_pct editable; limiar_verde_pct default 18; limiar_amarelo_pct default 10; tarifa_full editable NULLABLE (null=unknown); difal_enabled bool; difal_destino_uf (UF|null). Initial state ⇒ regime default SIMPLES 4% APPLIED with origem `default` explicit in response.
### DifalRate: uf(27), interna_pct, interestadual_pct (12% for MG,PR,RJ,RS,SC,SP; 7% rest; origem fixed SC), efetivo_pct = max(interna−interestadual, 0), origem_versao `padrao-2026`, override nullable {interna_pct, updated_at}. Override persists ONLY if Δ>0,049pp, audited. Mandatory `disclaimer` field on EVERY DIFAL surface: "seed padrão 2026 — não é orientação fiscal".
### Decomposition (SINGLE formula, Simulador AND Pedidos):
`retorno = preço − comissão(pct_modalidade×preço) − taxa_fixa(preço<79 ⇒ 6,50; senão 0) − frete(preço≥79 ⇒ frete_produto; senão 0) − imposto(aliquota×preço) − difal(efetivo_uf×preço, se enabled e destino conhecido) − tarifa_full(modalidade `full` ⇒ CalcProfile.tarifa_full, null ⇒ componente UNKNOWN propaga; demais modalidades ⇒ 0 explícito) − custo_erp`; `margem_pct = retorno/preço`.
- Port `Decompose(input) → {preco, comissao, taxa_fixa, frete, imposto, difal, tarifa_full, custo, margem_valor, margem_pct, componentes_desconhecidos[]}`. tarifa_full is an EXPLICIT nullable component. ANY unknown component ⇒ margem null + componentes_desconhecidos naming the missing (e.g. ["custo_erp"]). NEVER 0.
- Bidirectional `SolveTargetPrice(margem_alvo)` binary search, converges across the 79 step-discontinuity; margem inatingível ⇒ **200 `UNREACHABLE_TARGET`** citing the attainable ceiling.
- Port `DifalForUF(uf) → {efetivo_pct, versao}` — signature FROZEN IC-04 (M-08 consumes read-only; reflects active override).
### Endpoints: `GET/PUT /pricing/profile` (PUT validates aliquota 0–35%), `GET /pricing/difal` (DifalRate[27] ordered by UF), `PUT /pricing/difal/{uf}` (interna 0–35; Δ≤0,049pp ⇒ 200 without persisting override), `POST /pricing/simulations` extended additively (accepts full decomposition + margem→preço direction), `GET/POST/DELETE /pricing/scenarios` (newest-first).
### Error matrix (EXACT codes): aliquota fora 0–35 ⇒ **422 INVALID_RATE**; UF inválida ⇒ **404 UF_NOT_FOUND**; sem custo ⇒ 200 + decomposição unknown + blocking_state (NÃO erro); margem inatingível ⇒ 200 UNREACHABLE_TARGET; preço ≤ 0 ⇒ **422 INVALID_PRICE**; codprod inexistente ⇒ **404 ITEM_NOT_FOUND**; cenário inexistente ⇒ **404 SCENARIO_NOT_FOUND**.

## SEED DATA — R-04 (interna_pct modal per UF; 27 UFs)
AC 19,0 · AL 20,5 (post-01/04/2026 Lei 9.776/2025; pre-change 20% is an operator-override window — seed uses 20,5) · AM 20,0 · AP 18,0 · BA 20,5 · CE 20,0 · DF 20,0 · ES 17,0 · GO 19,0 · MA 23,0 · MG 18,0 · **MS 17,0 (DISPUTED — verify-at-execution; a source says 19%; seed 17,0 + operator-override expected)** · MT 17,0 · PA 19,0 · PB 20,0 · PE 20,5 · PI 22,5 · PR 19,5 · RJ 20,0 · RN 20,0 · RO 19,5 · RR 20,0 · RS 17,0 · SC 17,0 (origem) · SE 19,0 · SP 18,0 · TO 20,0. interestadual: 12% for MG,PR,RJ,RS,SC,SP; 7% all others. efetivo = max(interna−interestadual, 0). FCP excluded from interna_pct.

## ACCEPTANCE CRITERIA to map (validation-contract C01–C07)
- **C01 golden**: ≥10 cases (SIMPLES/PRESUMIDO × preço</≥79 × UF 12%/7% × custo conhecido/desconhecido × modalidade classico/premium/full) — every decimal component EXACT; soma fecha (preço = Σ componentes + margem); taxa fixa 6,50 se <79 senão 0; frete produto se ≥79 senão 0; tarifa_full só em full (0 explícito demais; null em full ⇒ unknown propaga); unknown ⇒ margem null + componentes_desconhecidos; decimal not float64.
- **C02 DIFAL seed+override**: GET /pricing/difal = 27 UFs ordered, interna=R-04, interestadual 12/7, efetivo=max(interna−inter,0), origem_versao padrao-2026, disclaimer present; PUT Δ>0,049pp ⇒ persisted+audited; Δ≤0,049pp ⇒ 200 no persist; aliquota 40 ⇒ 422 INVALID_RATE; UF XX ⇒ 404 UF_NOT_FOUND.
- **C03 solver**: margem-alvo 15% ⇒ preço whose re-sim returns 15,00% exact; 60% w/ comissão 16% ⇒ 200 UNREACHABLE_TARGET citing ceiling; old W1 shape preserved (regression).
- **C04 ports frozen**: Go contract test of Decompose + DifalForUF EXACT IC-04 sigs; DifalForUF reflects active override; grep for a 2nd decomposition impl in modules/orders = zero.
- **C05 tela** (QA/browser, F-02): decomposição per component + MarginChip (verde≥18/âmbar≥10/vermelho, CalcProfile thresholds); market comparison shows source/fetched_at/n_offers/n_sellers/match_status (IC-03); destino real muda DIFAL; toggle off recalcula sem DIFAL; sem custo ⇒ UnknownValue + blocking SEM_CUSTO (never 0); deep-link /precos?params=1 opens Parâmetros drawer; DIFAL drawer lists 27 UFs + disclaimer. Live comissão/frete via M-02 ports (hardcoded comissão = BLOCKING).
- **C06 aplicar** (QA/security, F-02): "aplicar preço" ⇒ price_update intent via /mutations, final status `previewed` (UI no approve), provider dispatcher OFF, ZERO ML write requests in log window.
- **C07 migrations+seams**: 0055–0059 hold the 3 tables, fixture = real count, diff only in owned paths.

## HOW TO WORK
Bulk-read the actual base files you need to size slices precisely (you have Read/Grep/Glob) — e.g. read the existing pricing transport/service/domain, market/domain/aggregation.go (decimal pattern), mutations/domain/protocol.go, MarginChip.tsx, feature-simulator/src/PricingSimulatorPage.tsx, the OpenAPI /pricing block. Keep slices ≤~300 changed lines (mechanical/generated exempt). Sequence F-01 so the **ic04-ports publication slice lands EARLY** (Decompose + DifalForUF committed + reviewable → unblocks M-08). Then F-02.

## DELIVERABLE
Write the full plan to `.mnfs/MIS-004-mvp-demo/M-07-simulador/_evidence/M07-plan.md` (all slice cards + the 6 required outputs, clearly sectioned). Then RETURN a compressed summary: slice count per feature, the DAG order, the ic04-ports slice id, any contract-satisfiability flags, and any residual open_questions (there should be none — if you find one, name the investigator dispatch needed instead of leaving it hanging).
