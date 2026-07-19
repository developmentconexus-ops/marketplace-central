# CHIP-T1 pricing-solver-tier4 — P2 PLAN

`base_sha: 18fbd91a5706f8a58f41ddf27f13c8f58a5dd6e4` (main) · branch `chip/t1-solver-tier4`
Scope = DESIGN-TARIFAS-ML §8.1 DEMO-CRÍTICO + §4 wiring + §2.5 degrau 4. Contingency lane §12: Claude-only (sonnet workers, plan by me).

## Key repo-truth constraints discovered (drive the design)

- **IC-04 FROZEN CONTRACT** (`ports/calc_ports_contract_test.go`): `assertFields` asserts EXACT field count on
  `domain.DecomposeInput`, `domain.Decomposition`, `domain.DifalForUFResult`, + `var _ func(DecomposeInput) Decomposition = Decompose`.
  → **Cannot add fields to DecomposeInput/Decomposition, cannot change Decompose signature** (shared w/ M-08).
  `SolveInput`/`SolveResult` are M-07-only (NOT frozen) → safe to extend.
- Threshold param therefore threads via an **unexported `decomposeWithLimiar(in, limiarCents)`**; `Decompose(in)` wraps it with the named default. Literal `7900` survives ONLY as `defaultTaxaFixaLimiarCents` (a named default, §8.1.3 partial removal, honest).
- Migration top on disk = 0067. **0068 free** (M-04 used 0065-0067 of its 0065-0069 block). I take 0068 ONLY.
- SDK is a single `packages/sdk-runtime/src/index.ts` (ownership-matrix `pricing.ts` was aspirational). Pricing block @1344-1442.
- OpenAPI `contracts/api/marketplace-central.openapi.yaml`: paths @2496 (/decompose) /2525 (/solve); schemas @4100 (PricingCalcInput) /4126 (PricingDecomposition) /4165 /4175 (Solve/Decompose responses).
- `cfg pgdb.Config` has `DefaultTenantID` only, NO installation. Config table scoped (tenant_id, installation_id) with installation_id DEFAULT '' ; resolver uses a default-installation sentinel (`""`) for demo (multi-installation-ready schema, cenário 12; single installation wired). Deferral documented.
- `13/16` must be **DB seed, never Go constant** → `GetDefaults` materializes via `INSERT (tenant_id,installation_id) ON CONFLICT DO NOTHING` (DB column DEFAULTs fill 13.00/16.00/sem_dados), then SELECT. No Go literal.
- `NewCalcService` call site = root.go:716 only (+ tests). Resolver added via additive `WithTariffResolver` chained method (nil-safe) → existing calc_service tests untouched.
- Only `calc_service.go:270` constructs SolveInput externally. Goldens build SolveInput w/o threshold → solver treats `TaxaFixaLimiarCents==0` as default.

## Slice cards (failing-test-first each; sonnet workers; I re-verify P5)

### Slice A — migration 0068 + domain TariffDefaults + repo + GET/PUT endpoints
- write_set: `migrations/0068_pricing_tariff_defaults.sql`, `pricing/domain/tariff_defaults.go`, `pricing/ports/tariff.go` (TariffDefaultsStore part), `pricing/adapters/postgres/calc_repository.go` (+GetTariffDefaults/UpsertTariffDefaults), `pricing/application/calc_service.go` (+GetTariffDefaults/PutTariffDefaults), `pricing/transport/calc_handler.go` (+2 handlers +routes), `pricing/adapters/postgres/calc_repository_test.go`, `platform/migrate/runner_test.go` (bump count fixture)
- table: tenant_id,installation_id (PK), comissao_classico_pct DEFAULT 13.00, comissao_premium_pct DEFAULT 16.00 (CHECK 0..100), frete_estimativa_amount numeric NULL (ADR-17 nunca 0), frete_policy text DEFAULT 'sem_dados' CHECK IN ('estimativa','sem_dados'), updated_at
- validation: `go test ./internal/modules/pricing/...`, migration count fixture green
- verification map: deliverable 1 → GET/PUT handlers + repo test; ADR-17 → frete_estimativa_amount NULL not 0

### Slice B — TariffResolver port + tier-4 adapter
- write_set: `pricing/ports/tariff.go` (TariffResolver, TariffRequest), `pricing/domain/tariff.go` (TariffResolution, ComponentResolution, Fonte), `pricing/adapters/tariffdefaults/resolver.go` + `_test.go`
- resolver resolves ONLY degrau 4: comissao by modalidade (classico→classico_pct, premium|full→premium_pct), stamp {Fonte:PADRAO,Degrau:4,Estimativa:true,Pct}; frete: policy estimativa+amount→Money+Estimativa, else sem_dados→Money=nil (NO-DATA). Structured so a future composite (1→2→3→4) implements same port unchanged.
- verification map: deliverable 2 → resolver_test asserts stamps + NO-DATA nil; ADR-17 → frete nil never 0

### Slice C — Layer 1 domain (threshold param + FreteDesconhecido segment-conditional + remove R$150 probe)
- write_set: `pricing/domain/decompose.go` (extract decomposeWithLimiar + defaultTaxaFixaLimiarCents), `pricing/domain/solve.go` (SolveInput+TaxaFixaLimiarCents, SolveResult+FreteDesconhecido, segment-conditional frete, remove const thresholdCents + R$150 probe), `pricing/domain/solve_golden_test.go` (+new cases)
- new behavior: structuralUnknowns evaluated at `limiar-1` (frete not consulted → catches custo/tarifa_full/difal only). frete nil blocks ONLY ≥limiar segment → FreteDesconhecido=true; <limiar segment solves normally. Custo nil stays Desconhecidos (hard block). Goldens stay green (frete present; threshold 0→default 7900).
- verification map: deliverable 3 → new golden cases (frete-nil low-segment solves; frete-nil high-segment→FreteDesconhecido; custo-nil→Desconhecidos)

### Slice D — CalcService wires resolver + injects comissao/frete before domain (+ tarifa block plumbing)
- write_set: `pricing/application/calc_service.go` (WithTariffResolver, resolve-before-domain, DecomposeResult/SolveOutput +Tarifa, threshold into SolveInput), `pricing/application/calc_service_test.go`
- comissao: req.comissao_pct empty + resolver → resolved pct (else 422 if no resolver, back-compat); frete: req.frete_produto nil → resolved frete (may be nil NO-DATA). Effective values + stamps → Tarifa on result.
- verification map: deliverable 3/4 → service test: resolver fills comissao+frete, tarifa stamped; manual override honored

### Slice E — Layer 2 transport + OpenAPI + SDK (SAME COMMIT) [needs contract-lock]
- write_set: `pricing/transport/calc_handler.go` (branch code cause; tarifa block; frete_desconhecido), `pricing/transport/calc_handler_test.go`, `contracts/api/marketplace-central.openapi.yaml`, `packages/sdk-runtime/src/index.ts`
- code branch: desconhecidos≠∅→DADOS_INCOMPLETOS · FreteDesconhecido→SEM_FRETE · !reached&&ceiling≠""→UNREACHABLE_TARGET · reached→none. tarifa block {comissao{valor,fonte,degrau,data,estimativa}, frete{valor,fonte,degrau,data,estimativa,sem_dados}}. comissao_pct→optional. code enum +SEM_FRETE/DADOS_INCOMPLETOS.
- verification map: deliverable 4 → handler test: blocking≠UNREACHABLE, tarifa present; OpenAPI+SDK atomic (§7 non-negotiable)

### Slice F — root.go additive registration
- write_set: `pricing/adapters/postgres/calc_repository.go` (store already there), `composition/root.go` (construct tariffdefaults resolver, `calcSvc.WithTariffResolver`)
- additive only; released at CLOSED, diff called out in CLOSED payload.

## Gates before CLOSED
go build/vet/test ./... (GOCACHE abs, no GOFLAGS) · governance clean worktree 40-hex BaseSha · sdk `tsc` 0 · P6 dual gate (cold Opus + adversarial sonnet REFUTE, agreement) · evidence complete.

## Contract-lock REQUEST (sent to hub before Slice E)
OpenAPI sections: /pricing/solve + /pricing/decompose responses (tarifa block, frete_desconhecido, code enum), PricingCalcInput (comissao_pct optional), + new /pricing/tariff-defaults path & PricingTariffDefaults schema, + PricingTarifa schema. SDK index.ts pricing block. All additive except comissao_pct required→optional (widening, back-compat).
