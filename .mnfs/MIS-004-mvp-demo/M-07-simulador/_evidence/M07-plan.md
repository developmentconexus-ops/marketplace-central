# M-07-simulador — P2 BATCH PLAN (cold Opus planner)

Base tree: `8b6c4b30`. Mission: MIS-004-mvp-demo (client demo 2026-07-20).
Binding: HARNESS-CORE §4/§5 + HARNESS-PROFILE + mission.md + IC-04 contract + R-04 seed.
Slice budget ≤~300 changed lines (mechanical/generated exempt). Internal DAG: F-01 → F-02.
**ic04-ports publication slice = F01-S3** (Decompose + DifalForUF committed + reviewable → unblocks sibling M-08). It lands 3rd, after only the two pure prerequisites (decimal helper, difal domain+seed).

Slice count: **F-01 = 8, F-02 = 6, total = 14.** Zero residual open_questions.

---

## SECTION 0 — ORCHESTRATOR P2 REVIEW CORRECTIONS (BINDING — override the cards below on conflict)

Two corrections applied by CHIP-M07 at P2 acceptance (planner artifact accepted-with-corrections):

- **C-A · NO 4th DIFAL table.** The grant (chip prompt g) names EXACTLY three tables: `pricing_calc_profiles`, `pricing_difal_rates`, `pricing_scenarios`. DIFAL override + audit are **nullable COLUMNS on `pricing_difal_rates`** per IC-04 (`override {interna_pct, updated_at}`; add `override_interna_pct`, `override_updated_at`, `override_actor` to satisfy "audited"). Migration block uses **0055–0058** (0055 profiles · 0056 difal_rates · 0057 difal seed · 0058 scenarios); **0059 stays RESERVE** (matches C07 "0059 reserva do bloco"). Fixture bump = 51 → **actual final count (55)**. A genuine audit-trail TABLE, if later judged necessary, is a REQUEST to the hub — never self-granted. This supersedes F01-S5's 5-file/0059-audit-table version.
- **C-B · npm, not pnpm.** Repo is an npm-workspaces monorepo (only `package-lock.json`; root `package.json` workspaces `apps/web`+`packages/*`). Every `pnpm --filter <pkg> <script>` in the cards below is REPLACED at dispatch time by the npm equivalent:
  - web tests: `npm run test --workspace @marketplace-central/web -- <fileOrPattern>` (script = `vitest run --config vitest.config.ts`)
  - sdk-runtime / feature-simulator tests: `npm run test --workspace @marketplace-central/<pkg>`
  - web build: `npm run build --workspace @marketplace-central/web` (or root `npm run web:build`)
  - typecheck: `npx tsc --noEmit -p <pkgDir>`
  Go lanes carry `GOCACHE=$(pwd)/.gocache` (+ hermetic `GOMODCACHE`) per PROFILE §2/§3 — bindings-block concern, enforced in each implementer dispatch.

---

## SECTION A — SLICE CARDS (dependency-DAG order)

### Slice F01-S1 — pricing-local decimal money helper + decimal value types
- feature: F-01
- complexity: complex (decimal engine foundation)
- validation_kind: unit-golden
- failing_test_first: `apps/server_core/internal/modules/pricing/domain/decimal_test.go` — asserts `parseRat("79.00")`, `formatRatHalfUp` half-up at money 2 places AND pct places (IC-04), rejects non-decimal strings, `nil *Money` = unknown (ADR-17, never zero). Golden vectors for rounding boundary (…x.xx5 up).
- write_set: [`apps/server_core/internal/modules/pricing/domain/decimal.go`, `apps/server_core/internal/modules/pricing/domain/decimal_test.go`]
- commands: [`go test ./apps/server_core/internal/modules/pricing/domain/...`]
- expected_artifacts: [go test PASS log for pricing/domain decimal]
- depends_on: []
- open_questions: []
- notes: G2 — market/domain has `formatRatHalfUp` but funcs are UNEXPORTED and cross-module reuse of another domain is a layering violation; rule-of-three says this is the 2nd occurrence, so a pricing-LOCAL helper mirroring `market/domain/aggregation.go:72-84` is the sanctioned choice (call it out; do NOT import market/domain). G1 — every downstream money computation routes through this one helper so C01 "decimal not float64" holds engine-wide.

### Slice F01-S2 — CalcProfile + DifalRate domain + R-04 seed + DifalForUF pure logic
- feature: F-01
- complexity: complex (DIFAL seed logic)
- validation_kind: unit-golden
- failing_test_first: `apps/server_core/internal/modules/pricing/domain/difal_test.go` — 27-UF seed table equals R-04 exactly; interestadual = 12% for {MG,PR,RJ,RS,SC,SP} else 7%; `efetivo = max(interna−interestadual, 0)`; origem_versao `padrao-2026`; MS seeded 17,0 (DISPUTED marker in a code comment, verify-at-execution). `calcprofile_test.go` — regime SIMPLES default 4%, PRESUMIDO default 9,25%, limiar_verde 18 / limiar_amarelo 10, tarifa_full nullable, initial state ⇒ SIMPLES 4% with origem `default` explicit.
- write_set: [`apps/server_core/internal/modules/pricing/domain/calcprofile.go`, `apps/server_core/internal/modules/pricing/domain/difal.go`, `apps/server_core/internal/modules/pricing/domain/difal_seed.go`, `apps/server_core/internal/modules/pricing/domain/difal_test.go`, `apps/server_core/internal/modules/pricing/domain/calcprofile_test.go`]
- commands: [`go test ./apps/server_core/internal/modules/pricing/domain/...`]
- expected_artifacts: [go test PASS with 27-UF golden + regime-default assertions]
- depends_on: [F01-S1]
- open_questions: []
- notes: G3 — DifalForUF pure computation lands here so the port slice (S3) only wires the interface. Disclaimer string constant `"seed padrão 2026 — não é orientação fiscal"` defined here for reuse by every DIFAL surface (C02/C05).

### Slice F01-S3 — **IC-04 PORTS PUBLICATION**: Decompose engine + DifalForUF port (EARLY, unblocks M-08)
- feature: F-01
- complexity: complex (decimal engine — the single decomposition formula)
- validation_kind: unit-golden + contract
- failing_test_first: (RED first) `apps/server_core/internal/modules/pricing/domain/decompose_golden_test.go` — **C01 ≥10 cases**: SIMPLES/PRESUMIDO × preço </≥79 × UF 12%/7% × custo conhecido/desconhecido × modalidade classico/premium/full; every component EXACT decimal; soma fecha (`preço = Σ componentes + margem_valor`); taxa_fixa 6,50 se preço<79 senão 0; frete = frete_produto se preço≥79 senão 0; tarifa_full só em modalidade `full` (0 explícito nas demais; `full` + tarifa_full null ⇒ componente UNKNOWN propaga); qualquer unknown ⇒ `margem_pct` null + `componentes_desconhecidos` nomeando o faltante (ex ["custo_erp"]); NEVER 0-as-unknown. PLUS `apps/server_core/internal/modules/pricing/ports/calc_ports_contract_test.go` — **C04**: freezes exact IC-04 signatures `Decompose(input) → {...}` and `DifalForUF(uf) → {efetivo_pct, versao}`; DifalForUF reflects active override.
- write_set: [`apps/server_core/internal/modules/pricing/ports/calc_ports.go` (Decompose + DifalForUF interfaces + input/output structs, FROZEN sigs), `apps/server_core/internal/modules/pricing/domain/decompose.go` (formula engine, big.Rat internal → decimal string), `apps/server_core/internal/modules/pricing/domain/decompose_golden_test.go`, `apps/server_core/internal/modules/pricing/ports/calc_ports_contract_test.go`]
- commands: [`go test ./apps/server_core/internal/modules/pricing/...`, `grep -rn "Decompose\|margem_valor" apps/server_core/internal/modules/orders` (must be zero — C04 no 2nd impl)]
- expected_artifacts: [golden test PASS log; contract test PASS log; empty grep output for orders decomposition]
- depends_on: [F01-S1, F01-S2]
- open_questions: []
- notes: **G3 — THIS is the M-08 unblock. Publish + hub-review this slice before proceeding.** Ports interface + engine impl land together because M-08 consumes the *single* formula (Simulador AND Pedidos), not just a type. Formula verbatim from IC-04 §Decomposition. `tarifa_full` is an explicit nullable component. G1 — engine is pure (no I/O), cost/comissão/frete arrive as already-resolved decimal inputs from ports (S6), keeping the golden hermetic.

### Slice F01-S4 — SolveTargetPrice binary search (bidirectional) + UNREACHABLE_TARGET
- feature: F-01
- complexity: complex (solver)
- validation_kind: unit-golden
- failing_test_first: `apps/server_core/internal/modules/pricing/domain/solve_golden_test.go` — **C03**: margem-alvo 15% ⇒ preço whose re-`Decompose` returns margem_pct 15,00% EXACT; converges ACROSS the preço=79 step-discontinuity (taxa_fixa + frete flip); margem-alvo 60% w/ comissão 16% ⇒ `UNREACHABLE_TARGET` citing the attainable ceiling; monotonic bracket + tolerance documented.
- write_set: [`apps/server_core/internal/modules/pricing/domain/solve.go`, `apps/server_core/internal/modules/pricing/domain/solve_golden_test.go`]
- commands: [`go test ./apps/server_core/internal/modules/pricing/domain/...`]
- expected_artifacts: [solver golden PASS incl. discontinuity-crossing + unreachable-ceiling cases]
- depends_on: [F01-S3]
- open_questions: []
- notes: G2 — binary search chosen over closed-form because the 79-threshold makes retorno(preço) piecewise (taxa_fixa+frete discontinuity); solver must bracket both segments and report the reachable ceiling on failure, mapped to HTTP 200 `UNREACHABLE_TARGET` in S7.

### Slice F01-S5 — migrations 0055–0058 (3 tenant tables; 0059 reserve) + fixture bump 51→55  ⟵ SECTION 0 · C-A
- feature: F-01
- complexity: standard
- validation_kind: integration-pg
- failing_test_first: `apps/server_core/internal/platform/migrate/runner_test.go` bumped to the tree's ACTUAL final count (51 current + 4 new = **55**) at :25 and :64 (RED until the 4 files exist) — asserts inventory count and lexicographic order. Implementer sets both lines to the REAL count if 51 is stale.
- write_set: [`apps/server_core/migrations/0055_pricing_calc_profiles.sql`, `apps/server_core/migrations/0056_pricing_difal_rates.sql`, `apps/server_core/migrations/0057_pricing_difal_seed.sql`, `apps/server_core/migrations/0058_pricing_scenarios.sql`, `apps/server_core/internal/platform/migrate/runner_test.go`]
- commands: [`go test ./apps/server_core/internal/platform/migrate/...`, hermetic-pg apply (warm .gomodcache first per memory), `SELECT count(*) FROM pricing_difal_rates;` ⇒ 27]
- expected_artifacts: [migrate runner PASS at the real count (55); psql `\d pricing_calc_profiles|pricing_difal_rates|pricing_scenarios` showing tenant_id NOT NULL on every table; `\d pricing_difal_rates` showing the override columns; SELECT 27 rows]
- depends_on: [F01-S2]
- notes: EXACTLY 3 granted tables = pricing_calc_profiles, pricing_difal_rates, pricing_scenarios — **every table carries tenant_id** (grant condition). **DIFAL override + audit = NULLABLE COLUMNS on `pricing_difal_rates`**, NOT a 4th table (grant names only 3; a 4th ⇒ REQUEST to hub). IC-04 override shape `{interna_pct, updated_at}` → columns `override_interna_pct numeric NULL`, `override_updated_at timestamptz NULL`, plus `override_actor text NULL` for the "audited" requirement (Δ>0,049pp persists, audited per IC-04). 0056 defines the table WITH these columns; the 27-row seed (0057) leaves them NULL. **0059 stays RESERVE** (C07 "0059 reserva do bloco"). 0054/0060–0064 left free. **C07** anchor. Bump BOTH runner_test lines to the REAL count in THIS slice (invariant). Follow hermetic-lane CREATE DATABASE retry-loop gotcha (memory).

### Slice F01-S6 — postgres repositories (profile/difal/scenarios) + cost port + ROOT shim
- feature: F-01
- complexity: standard
- validation_kind: integration-pg
- failing_test_first: `apps/server_core/internal/modules/pricing/adapters/postgres/calc_repository_test.go` — profile upsert+read (initial ⇒ SIMPLES 4% origem `default`); difal read 27 ordered by UF + override reflected; scenarios CRUD newest-first; tenant-scoped. PLUS `apps/server_core/internal/modules/pricing/adapters/costread/reader_test.go` — pricing cost port maps internal_read `GetCostAsOf` → decimal Money, `nil` when absent (never zero).
- write_set: [`apps/server_core/internal/modules/pricing/adapters/postgres/calc_repository.go`, `apps/server_core/internal/modules/pricing/adapters/postgres/calc_repository_test.go`, `apps/server_core/internal/modules/pricing/ports/cost_read.go` (pricing-side cost port), `apps/server_core/internal/modules/pricing/adapters/costread/reader.go`, `apps/server_core/internal/modules/pricing/adapters/costread/reader_test.go`]
- commands: [`go test ./apps/server_core/internal/modules/pricing/adapters/...`]
- expected_artifacts: [pg repo test PASS; cost-adapter test PASS incl. nil-on-absent case]
- depends_on: [F01-S5, F01-S2]
- notes: Cost consumed via IC-02 `internal_read/ports/reader.go:52 GetCostAsOf` — pricing declares its OWN narrow port (mirror profitability's pattern) + root shim wired in S8 under ROOT-M07. Frete via IC-06 `connectors/ports/shipping_read.go GetFreeShippingCost` (present/tested) — reuse, do NOT reimplement. G1 — engine (S3) stays pure; these adapters resolve inputs into it.

### Slice F01-S7 — HTTP transport: profile/difal/scenarios endpoints + extended simulations + error matrix
- feature: F-01
- complexity: standard
- validation_kind: unit + contract
- failing_test_first: `apps/server_core/internal/modules/pricing/transport/calc_handler_test.go` — routes + IC-04 error matrix EXACT: aliquota fora 0–35 ⇒ 422 INVALID_RATE; UF inválida ⇒ 404 UF_NOT_FOUND; preço ≤0 ⇒ 422 INVALID_PRICE; codprod inexistente ⇒ 404 ITEM_NOT_FOUND; cenário inexistente ⇒ 404 SCENARIO_NOT_FOUND; sem custo ⇒ 200 + decomposição unknown + blocking_state (NÃO erro); margem inatingível ⇒ 200 UNREACHABLE_TARGET; DIFAL PUT Δ≤0,049pp ⇒ 200 no-persist. Existing W1 POST DTO + GET `{items:[...]}` shape preserved (regression, C03).
- write_set: [`apps/server_core/internal/modules/pricing/transport/calc_handler.go` (new handlers: GET/PUT /pricing/profile, GET /pricing/difal, PUT /pricing/difal/{uf}, GET/POST/DELETE /pricing/scenarios, extended POST /pricing/simulations decomposition+margem→preço), `apps/server_core/internal/modules/pricing/transport/http_handler.go` (ADDITIVE Register lines only), `apps/server_core/internal/modules/pricing/transport/calc_handler_test.go`, `apps/server_core/internal/modules/pricing/application/calc_service.go` (orchestrates engine+repos+ports), `apps/server_core/internal/modules/pricing/application/calc_service_test.go`]
- commands: [`go test ./apps/server_core/internal/modules/pricing/transport/... ./apps/server_core/internal/modules/pricing/application/...`]
- expected_artifacts: [transport test PASS covering full error matrix; W1 regression PASS]
- depends_on: [F01-S3, F01-S4, F01-S6]
- notes: POST /pricing/simulations extended ADDITIVELY — existing float64 DTO fields (http_handler.go:72-82) preserved; new decomposition + margem-alvo direction are additive. Register additions at transport :53-54 region. May approach 300 lines — if so, split scenarios CRUD into F01-S7b (still same DAG position, depends identical).

### Slice F01-S8 — OpenAPI /pricing/* additive + pricing.ts SDK (same commit) + composition-root wiring
- feature: F-01
- complexity: standard
- validation_kind: contract + web-build
- failing_test_first: `packages/sdk-runtime/src/pricing.test.ts` — pricing client methods (getProfile/putProfile/listDifal/putDifal/listScenarios/create/delete/runSimulation) hit correct paths + decode decomposition/difal/error shapes; mirrors market.ts self-contained pattern. OpenAPI validated by existing contract lane.
- write_set: [`contracts/api/marketplace-central.openapi.yaml` (ADDITIVE new paths after :2325 + new schemas — /pricing/profile, /pricing/difal, /pricing/difal/{uf}, /pricing/scenarios; existing /simulations preserved), `packages/sdk-runtime/src/pricing.ts` (NEW standalone, Pricing*-prefixed, own client/error/money types, no import from index.ts), `packages/sdk-runtime/src/pricing.test.ts`, `packages/sdk-runtime/src/index.ts` (BARREL-M07: exactly one `export * from "./pricing"`), `apps/server_core/internal/composition/root.go` (ROOT-M07: additive imports :95-100 region + PRICING wiring :689-702 region — new calc service/repos/cost shim, existing constructors untouched, no stub/nil on live paths)]
- commands: [`npx tsc --noEmit -p packages/sdk-runtime`, `npm run test --workspace @marketplace-central/sdk-runtime`, `go build ./apps/server_core/...`, OpenAPI lint/contract lane]
- expected_artifacts: [tsc clean; sdk-runtime vitest PASS; go build PASS; OpenAPI validates]
- depends_on: [F01-S7]
- notes: OpenAPI + sdk-runtime land SAME commit (invariant). BARREL-M07 = one additive line only. ROOT-M07 = additive wiring only. market.ts is the exact mirror pattern for pricing.ts (own createPricingClient, getJson/postJson, throws `{status, error}`). **This slice publishes the SDK surface F-02 consumes — F-02 starts after this.**

---

### Slice F02-S1 — /precos page scaffold + route rewire + AppRouter test (APPTEST-M07)
- feature: F-02
- complexity: standard
- validation_kind: web-vitest
- failing_test_first: `apps/web/src/pages/precos/PricingPage.test.tsx` — page mounts, loads profile+product via pricing.ts client (mocked), renders shell regions (decomposição, parâmetros trigger, comparação, aplicar). `apps/web/src/app/AppRouter.test.tsx` :93-97 updated (APPTEST-M07, that single /precos case only) to assert new page marker.
- write_set: [`apps/web/src/pages/precos/PricingPage.tsx`, `apps/web/src/pages/precos/PricingPage.test.tsx`, `apps/web/src/pages/precos/index.ts`, `apps/web/src/routes/precos.tsx` (rewire from legacy PricingSimulatorPage to new page), `apps/web/src/app/AppRouter.test.tsx` (APPTEST-M07: /precos case only)]
- commands: [`npm run test --workspace @marketplace-central/web -- PricingPage AppRouter`]
- expected_artifacts: [vitest PASS for PricingPage + AppRouter /precos case]
- depends_on: [F01-S8]
- notes: APPTEST-M07 grant = edit ONLY the :93-97 /precos case, zero other cases; legacy redirect ["/simulator","/precos"] at :105 untouched. Consumes pricing.ts from F01-S8. C05 shell anchor.

### Slice F02-S2 — decomposition panel (per-component) + MarginChip + UnknownValue/blocking SEM_CUSTO
- feature: F-02
- complexity: standard
- validation_kind: web-vitest
- failing_test_first: `apps/web/src/pages/precos/DecompositionPanel.test.tsx` — renders each component (comissão/taxa_fixa/frete/imposto/difal/tarifa_full/custo/margem); MarginChip fed CalcProfile thresholds (verde≥18/âmbar≥10/vermelho); `componentes_desconhecidos` ⇒ UnknownValue (—) + blocking banner `SEM_CUSTO`, NEVER renders 0 for unknown; tarifa_full unknown propagates.
- write_set: [`apps/web/src/pages/precos/DecompositionPanel.tsx`, `apps/web/src/pages/precos/DecompositionPanel.test.tsx`]
- commands: [`npm run test --workspace @marketplace-central/web -- DecompositionPanel`]
- expected_artifacts: [vitest PASS incl. unknown-not-zero + threshold-band cases]
- depends_on: [F02-S1]
- notes: Reuse `packages/ui` MarginChip (pass CalcProfile {healthy: limiar_verde, tight: limiar_amarelo}) + UnknownValue/ErrorState/EmptyState — do NOT fork ui. C05 anchor. DO NOT copy legacy inline marginColor 20/10 thresholds.

### Slice F02-S3 — Parâmetros drawer (CalcProfile edit) + DIFAL drawer (27 UFs + disclaimer) + deep-link + destino recalc
- feature: F-02
- complexity: standard
- validation_kind: web-vitest
- failing_test_first: `apps/web/src/pages/precos/ParamsDrawer.test.tsx` + `DifalDrawer.test.tsx` — Parâmetros drawer edits regime/aliquota/limiares/tarifa_full/difal_enabled/destino_uf; deep-link `/precos?params=1` opens Parâmetros drawer on mount; changing destino UF re-decomposes with that UF's efetivo; toggle difal_enabled OFF recomputes without DIFAL component; DIFAL drawer lists 27 UFs ordered + shows disclaimer on the surface.
- write_set: [`apps/web/src/pages/precos/ParamsDrawer.tsx`, `apps/web/src/pages/precos/DifalDrawer.tsx`, `apps/web/src/pages/precos/ParamsDrawer.test.tsx`, `apps/web/src/pages/precos/DifalDrawer.test.tsx`]
- commands: [`npm run test --workspace @marketplace-central/web -- ParamsDrawer DifalDrawer`]
- expected_artifacts: [vitest PASS incl. deep-link mount + destino-recalc + toggle-off + 27-UF-with-disclaimer]
- depends_on: [F02-S2]
- notes: Disclaimer text `"seed padrão 2026 — não é orientação fiscal"` mandatory on every DIFAL surface (C02/C05). aliquota client-validates 0–35 before PUT (server re-validates → 422 INVALID_RATE).

### Slice F02-S4 — market comparison panel via market.ts SDK (IC-03 HTTP)
- feature: F-02
- complexity: standard
- validation_kind: web-vitest
- failing_test_first: `apps/web/src/pages/precos/MarketComparison.test.tsx` — calls existing market.ts `listMarketAggregates`; renders source / fetched_at / n_offers / n_sellers / match_status; NO_PRICE_EVIDENCE / INSUFFICIENT_MARKET states shown, never fabricated.
- write_set: [`apps/web/src/pages/precos/MarketComparison.tsx`, `apps/web/src/pages/precos/MarketComparison.test.tsx`]
- commands: [`npm run test --workspace @marketplace-central/web -- MarketComparison`]
- expected_artifacts: [vitest PASS with all IC-03 evidence fields + blocking states]
- depends_on: [F02-S2]
- notes: IC-03 is HTTP-ONLY (confirmed: no Go cross-module market port) — FE consumes market aggregates via existing `packages/sdk-runtime/src/market.ts` (createMarketPriceIntelClient); the BACKEND engine correctly needs NO market-aggregate port. verdict_label is always null from M-02 (margin labels are M-07-owned).

### Slice F02-S5 — "aplicar preço" → price_update intent via /mutations, previewed ceiling (no approve)
- feature: F-02
- complexity: complex (override-threshold / safety ceiling)
- validation_kind: web-vitest + manual-QA
- failing_test_first: `apps/web/src/pages/precos/ApplyPrice.test.tsx` — **C06**: "aplicar preço" creates a `price_update` intent via /mutations SDK; resulting/asserted final status `previewed`; UI exposes NO approve/apply control (previewed is the safe ceiling); no auto-transition past previewed.
- write_set: [`apps/web/src/pages/precos/ApplyPriceAction.tsx`, `apps/web/src/pages/precos/ApplyPrice.test.tsx`]
- commands: [`npm run test --workspace @marketplace-central/web -- ApplyPrice`]
- expected_artifacts: [vitest PASS asserting previewed-only + no approve control]
- depends_on: [F02-S2]
- notes: Backend `mutations/domain/protocol.go` — price_update enabled; legal draft→previewed→approved; "previewed" is the safe ceiling (UI must never expose approve/apply). Manual-QA (P7) confirms provider dispatcher OFF + ZERO ML write requests in log window. Consumes /mutations read/create via SDK — no write into mutations module (FORBIDDEN path respected).

### Slice F02-S6 — legacy feature-simulator retheme/absorb (IC-05 survives) + final build + QA drive
- feature: F-02
- complexity: standard
- validation_kind: web-build + manual-QA
- failing_test_first: `packages/feature-simulator/src/PricingSimulatorPage.test.tsx` updated — package still builds/exports (IC-05: package stays alive); retheme aligns thresholds to MarginChip 18/10 (drop inline 20/10). Full `npm run build` green.
- write_set: [`packages/feature-simulator/src/PricingSimulatorPage.tsx` (retheme — swap inline marginColor/marginBg 20/10 for MarginChip; absorb batch flows the new page keeps), `packages/feature-simulator/src/PricingSimulatorPage.test.tsx`]
- commands: [`npm run build --workspace @marketplace-central/feature-simulator && npm run test --workspace @marketplace-central/feature-simulator`, `npm run build --workspace @marketplace-central/web`]
- expected_artifacts: [feature-simulator build+test PASS; web build PASS; P7 browser QA log vs C05/C06]
- depends_on: [F02-S3, F02-S4, F02-S5]
- notes: IC-05 — feature-simulator package must survive (do NOT delete). Retheme or absorb into new /precos flow; keep exports. Final web-build gate + hand to P7 QA live-drive.

---

## SECTION B — PER-FEATURE WRITE-SET (write-DAG)

### F-01 (backend) — all within granted exclusive paths
| Slice | Write targets (dir/file) |
|---|---|
| F01-S1 | pricing/domain/decimal.go(+test) |
| F01-S2 | pricing/domain/{calcprofile,difal,difal_seed}.go(+tests) |
| F01-S3 | pricing/ports/calc_ports.go, pricing/domain/decompose.go(+tests) |
| F01-S4 | pricing/domain/solve.go(+test) |
| F01-S5 | migrations/0055–0058*.sql (0059 reserve), platform/migrate/runner_test.go (fixture bump) |
| F01-S6 | pricing/adapters/postgres/calc_repository.go(+test), pricing/ports/cost_read.go, pricing/adapters/costread/reader.go(+test) |
| F01-S7 | pricing/transport/{calc_handler,http_handler}.go(+test), pricing/application/calc_service.go(+test) |
| F01-S8 | contracts/api/…openapi.yaml (additive), sdk-runtime/src/pricing.ts(+test), sdk-runtime/src/index.ts (BARREL-M07), composition/root.go (ROOT-M07) |

### F-02 (frontend) — all within granted paths
| Slice | Write targets |
|---|---|
| F02-S1 | apps/web/src/pages/precos/{PricingPage,index}, routes/precos.tsx, app/AppRouter.test.tsx (APPTEST-M07) |
| F02-S2 | apps/web/src/pages/precos/DecompositionPanel.tsx(+test) |
| F02-S3 | apps/web/src/pages/precos/{ParamsDrawer,DifalDrawer}.tsx(+tests) |
| F02-S4 | apps/web/src/pages/precos/MarketComparison.tsx(+test) |
| F02-S5 | apps/web/src/pages/precos/ApplyPriceAction.tsx(+test) |
| F02-S6 | packages/feature-simulator/src/PricingSimulatorPage.tsx(+test) |

Write-DAG serialization point: `apps/server_core/internal/modules/pricing/transport/http_handler.go` (S7 Register) and `sdk-runtime/src/index.ts` (S8 barrel) — single-writer, no cross-slice contention. `composition/root.go` written once (S8). No two slices write the same file except the deliberate additive touches (http_handler.go S1..none/S7; runner_test.go S5 only).

---

## SECTION C — CONTRACT-SATISFIABILITY CHECK (vs current contract + sibling tracks)

Current OpenAPI pricing block: `/pricing/simulations` (:2240-2284), `/pricing/simulations/batch` (:2285-2325); next path `/connectors/...` at :2326. Current SDK: pricing.ts ABSENT; index.ts exports erpImport + market.

| Claimed surface | Current state | Verdict |
|---|---|---|
| `GET/PUT /pricing/profile` | ABSENT | ADD — clean, additive after :2325 |
| `GET /pricing/difal`, `PUT /pricing/difal/{uf}` | ABSENT | ADD — clean |
| `GET/POST/DELETE /pricing/scenarios` | ABSENT | ADD — clean |
| `POST /pricing/simulations` (extended) | EXISTS :2251 | EXTEND ADDITIVELY — W1 float64 DTO + `{items}` GET preserved (C03 regression guards) |
| `packages/sdk-runtime/src/pricing.ts` | ABSENT | CREATE — standalone, market.ts mirror |
| `index.ts` barrel `export * from "./pricing"` | ABSENT | ADD 1 line — BARREL-M07 |
| Decompose / DifalForUF Go ports | ABSENT (no existing Decompose/CalcProfile) | CREATE — no collision |

Sibling-track disjointness: M-05 (listings) owns modules/listings + listing paths; M-08 (orders) owns modules/orders + order paths; M-07 owns modules/pricing + /pricing/*. **No path/schema/SDK-identifier collision.** M-08 is a read-only CONSUMER of DifalForUF (frozen sig, F01-S3) — a dependency, not a write collision.

**FLAGS: none.** No colliding/occupied path. No ownership collision. All F-01/F-02 write targets fall inside the granted exclusive paths or the three narrow grants (BARREL-M07 / ROOT-M07 / APPTEST-M07).

---

## SECTION D — PER-CRITERION VERIFICATION MAP (C01–C07)

| Crit | Verification command / QA step | Carrying file(s) | Slice |
|---|---|---|---|
| **C01** golden ≥10 (every decimal component exact, soma fecha, taxa/frete 79-rule, tarifa_full full-only, unknown⇒null not 0, decimal not float64) | `go test ./apps/server_core/internal/modules/pricing/domain/... -run Decompose` | decompose_golden_test.go; decompose.go; decimal.go | F01-S3 (+S1) |
| **C02** DIFAL seed 27/ordered, interna=R-04, inter 12/7, efetivo=max(...,0), origem padrao-2026, disclaimer, PUT Δ>0,049 persist+audit / Δ≤ 200 no-persist, aliq 40⇒422 INVALID_RATE, UF XX⇒404 UF_NOT_FOUND | `go test .../pricing/domain -run Difal` + `go test .../pricing/transport -run Difal` + `SELECT count(*) pricing_difal_rates=27` | difal_test.go; difal_seed.go; calc_handler_test.go; 0057_pricing_difal_seed.sql; 0056_pricing_difal_rates.sql (override_* columns) | F01-S2, S5, S7 |
| **C03** solver 15%⇒re-sim 15,00% exact; 60%@comissão16%⇒200 UNREACHABLE_TARGET citing ceiling; W1 shape regression | `go test .../pricing/domain -run Solve` + `go test .../pricing/transport -run Simulations` (W1 regression) | solve_golden_test.go; solve.go; calc_handler_test.go | F01-S4, S7 |
| **C04** Go contract test Decompose+DifalForUF exact IC-04 sigs; DifalForUF reflects override; zero 2nd decomposition in modules/orders | `go test .../pricing/ports -run Contract` + `grep -rn "Decompose\|margem_valor" apps/server_core/internal/modules/orders` (empty) | calc_ports_contract_test.go; calc_ports.go | F01-S3 |
| **C05** tela: decomposição per-component + MarginChip bands (18/10 CalcProfile) + market source/fetched_at/n_offers/n_sellers/match_status + destino recalc + toggle-off + sem custo⇒UnknownValue+SEM_CUSTO(never 0) + deep-link /precos?params=1 + DIFAL drawer 27+disclaimer + live comissão/frete (hardcoded=BLOCKING) | `npm run test --workspace @marketplace-central/web -- DecompositionPanel ParamsDrawer DifalDrawer MarketComparison` + P7 browser QA drive | DecompositionPanel.tsx, ParamsDrawer.tsx, DifalDrawer.tsx, MarketComparison.tsx (+tests) | F02-S2, S3, S4 |
| **C06** aplicar⇒price_update intent via /mutations, final `previewed` (UI no approve), provider dispatcher OFF, ZERO ML write requests in log window | `npm run test --workspace @marketplace-central/web -- ApplyPrice` + P7 QA security drive (inspect mutation status + provider log) | ApplyPriceAction.tsx, ApplyPrice.test.tsx | F02-S5 |
| **C07** 0055–0058 hold 3 tables (0059 reserve), fixture=real count (55), diff only in owned paths | `go test .../platform/migrate/...` + `git diff --name-only` scoped to granted paths | migrations 0055–0058, runner_test.go | F01-S5 (+ final diff audit at close) |

Every criterion maps to ≥1 named command AND carrying file(s). No unmapped criterion.

---

## SECTION E — SEAM-CLOSURE CHECKLIST

### F-01 seams
| Seam | Closure | Where |
|---|---|---|
| Composition-root wiring | ROOT-M07 grant — additive New* for calc service/repos/cost shim | F01-S8 root.go :95-100 + :689-702 |
| SDK client surface (METHODS not just types) | pricing.ts createPricingClient with getProfile/putProfile/listDifal/putDifal/scenarios CRUD/runSimulation | F01-S8 pricing.ts |
| SDK barrel export | BARREL-M07 — one additive line | F01-S8 index.ts |
| OpenAPI spec | additive /pricing/* paths+schemas, same commit as SDK | F01-S8 openapi.yaml |
| Migrations + fixture bump | 0055–0058 (0059 reserve) + runner_test 51→55 same slice | F01-S5 |
| Cost port (IC-02) consumption | pricing-side cost port + costread adapter + root shim | F01-S6 (port/adapter) + F01-S8 (shim wiring) |
| Frete port (IC-06) consumption | reuse existing connectors GetFreeShippingCost (no new seam) | F01-S6 |
| ic04-ports M-08 publication | Decompose+DifalForUF frozen, committed early, hub-reviewed | F01-S3 |

### F-02 seams
| Seam | Closure | Where |
|---|---|---|
| Router/shell test | AppRouter.test.tsx /precos case (APPTEST-M07) + route rewire | F02-S1 |
| SDK consumption (pricing) | pricing.ts client methods used by page | F02-S1..S5 |
| SDK consumption (market IC-03 HTTP) | market.ts listMarketAggregates | F02-S4 |
| ui component reuse | MarginChip/UnknownValue/Loading/Error/Empty from packages/ui (consume, not fork — packages/ui FORBIDDEN) | F02-S2, S3 |
| /mutations seam (price_update) | ApplyPrice creates intent via SDK, previewed ceiling | F02-S5 |
| Legacy package survival (IC-05) | feature-simulator retheme, keeps exports | F02-S6 |

Hub-owned post-merge steps (named, not in write-set): dev-stack rebuild/re-point to merged tree; governance registry `modules.json` entry lands via chip merge (never pre-merge on main); P7 live browser QA; provider-dispatcher-OFF log-window inspection for C06.

---

## SECTION F — GOLDEN-TEST-FIRST CARDS (complex candidates)

### GT-1 — Decimal decomposition engine (F01-S3)
- Test file (RED first): `apps/server_core/internal/modules/pricing/domain/decompose_golden_test.go`
- ≥10 vectors over the C01 matrix: {SIMPLES 4%, PRESUMIDO 9,25%} × {preço 78,90 / 79,00 / 150,00} × {UF 12%: SP · UF 7%: BA} × {custo conhecido / nil} × {classico / premium / full(tarifa_full set) / full(tarifa_full null⇒unknown)}.
- Each vector asserts EXACT decimal for comissão, taxa_fixa (6,50 <79 else 0), frete (produto ≥79 else 0), imposto (aliquota×preço), difal (efetivo_uf×preço when enabled+destino known), tarifa_full (full-only, explicit 0 elsewhere), custo, margem_valor, margem_pct; **soma fecha**: `preço = Σ componentes + margem_valor`; unknown component ⇒ margem null + `componentes_desconhecidos`. Assert types are decimal string, no float64 in the path.
- Engine written GREEN only after test is RED and reviewed.

### GT-2 — Bidirectional binary-search solver (F01-S4)
- Test file (RED first): `apps/server_core/internal/modules/pricing/domain/solve_golden_test.go`
- Vectors: (a) margem-alvo 15% ⇒ solved preço whose re-Decompose margem_pct = 15,00% EXACT; (b) target straddling preço=79 ⇒ solver crosses the taxa_fixa+frete discontinuity and still converges; (c) margem-alvo 60% with comissão 16% ⇒ UNREACHABLE_TARGET returning the attainable ceiling; (d) monotonic-bracket + tolerance invariant.
- Solver written GREEN only after RED + review. Depends on GT-1 engine (re-sim is the oracle).

---

## RETURN SUMMARY (data for orchestrator)
- Slice count: F-01 = 8 (F01-S1…S8), F-02 = 6 (F02-S1…S6). Total 14.
- DAG order: F01-S1 → F01-S2 → **F01-S3(ic04-ports)** → F01-S4 → F01-S5 → F01-S6 → F01-S7 → F01-S8 → F02-S1 → F02-S2 → {F02-S3 ∥ F02-S4 ∥ F02-S5} → F02-S6. (F01-S5 depends only on S2 → may run parallel to S3/S4; F02-S3/S4/S5 all depend only on F02-S2 → parallel.)
- ic04-ports publication slice: **F01-S3** (Decompose + DifalForUF, frozen sigs + engine + C01 golden + C04 contract), unblocks M-08.
- Contract-satisfiability flags: NONE. All /pricing/* sub-paths absent (additive-clean); POST /simulations extended additively; pricing.ts new; no sibling collision (M-05/M-08 disjoint; M-08 = read-only consumer).
- Ownership collisions: NONE. All writes inside granted exclusive paths + BARREL-M07/ROOT-M07/APPTEST-M07.
- Residual open_questions: NONE. (In-flight verify-at-execution notes, not blockers: MS interna_pct seeded 17,0 DISPUTED per R-04 — carried as code comment for execution-time confirmation, not a planning unknown; migration fixture final count = 55 = 51 current + 4 [SECTION 0 · C-A]; DIFAL override+audit = columns on pricing_difal_rates, NOT a 4th table.)
