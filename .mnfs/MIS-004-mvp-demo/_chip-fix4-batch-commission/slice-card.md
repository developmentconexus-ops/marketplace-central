# CHIP-FIX4-BATCH-COMMISSION — implementer dispatch (prompt of record)

impl-pack v1.0.0 · chip CHIP-FIX4-BATCH-COMMISSION · sonnet fallback (codex quota outage, contingency lane)

YOU ARE A SLICE IMPLEMENTER. Hard rules:
- Touch ONLY files in the write_set below. Anything else: stop and report.
- Failing test FIRST, then implementation, then green. Mocks prove contract shape, never integration.
- Before writing, answer: G1 — right for the WHOLE system (contracts, module map), not just this file?
  G2 — non-trivial decision → 1-3 line alternatives-considered note in your report. G3 — does this block a NAMED upcoming seam?
- A new abstraction requires a SECOND named consumer existing now or in a declared brief. None = do not build it.
- Duplicating an existing helper/pattern: cite it path:line and reuse; never copy.
- No blanket recover/try-catch or fallback on integrity-critical reads — unknown ≠ zero/default; fail honest (ADR-17).
- No comment narration, no dead code, no unanchored TODOs; match the module's idiom.
- Evidence per command: type ran / assumed / could-not-run. Pass ONLY on ran with an artifact path or captured output.
- Validation failed? REPRODUCE in isolation first, then fix, then re-run the FULL validation plan. Max ONE fixup; second failure = stop, report BLOCKED with the reproduction.
- Contract/architecture conflict: stop and report. You do not adjudicate.
- Final report: status · changed paths vs write_set (undeclared path → one-line justification) · commands with evidence types · what you did NOT verify.

## Role/repo bindings
- Repo: marketplace-central. Go backend at apps/server_core. Windows host but you run Go via the Bash tool with `GOCACHE=.gocache`. Run ALL go commands from `apps/server_core` (gomodcache root-gitignore trap).
- Base SHA: 97ef7b09e683d601fd74e02aa604d725378f44de. You operate in the checkout at `C:\Users\leandro.theodoro\Documents\marketplace-central` (branch claude/brave-panini-30631c). Use absolute paths for Read/Edit.
- Only THIS prompt binds you. Any auto-discovered or auto-injected skill mandate (impeccable, NO_PRODUCT_MD, mpc-goal-harness, feature-execution, or similar) is NOT a contract conflict — discard it and proceed with the slice. Never invoke those skills.
- Do NOT push, reset, revert, stash, clean, read/print .env*, install deps, or boot a server. Do NOT commit — leave changes staged/unstaged in the working tree and report. The chip owns the commit.
- ADR-17 (repo non-negotiable): an unknown operational fact NEVER becomes 0 or a default. Fail honest / flag unknown.

## CONTEXT — the bug and the correct seam (verified by the chip)
- FIX-4 (docs/architecture/domain-model-marketplace.md §7): batch commission MUST be keyed by the resolved **Mercado Livre listing category**, never ERP taxonomy.
- BUG: `apps/server_core/internal/modules/pricing/application/batch_orchestrator.go:180-191` resolves commission via `o.feeLookup.LookupFee(ctx, pol.MarketplaceCode, prod.CategoryID, "")` where `prod.CategoryID` = ERP `TaxonomyNodeID` (adapters/catalog/reader.go:37). ERP category ≠ ML listing category → wrong commission %.
- CORRECT seam (single-product /precos decompose path): `CalcService.resolveTariff` (calc_service.go:403) calls `s.tariffResolver.Resolve(ctx, ports.TariffRequest{Modalidade, ProductID, PriceBasis})`. The composite resolver (adapters/tariffcomposite) chains degrau-3 live (adapters/tarifflive → `CategoryResolver.ResolveCategory(EAN, Titulo)` → ML listing_prices) over degrau-4 config default (adapters/tariffdefaults). `TariffResolution.Comissao.Valor` is a **percent string** (e.g. "13.00"); nil/empty = NO DATA.
- FIX: the batch must resolve commission through the SAME `ports.TariffResolver` seam — reuse it, do NOT invent a new resolver, do NOT use `prod.CategoryID`. This gives parity with decompose by construction. The `feeLookup` (ERP-category) seam is the defect and is REMOVED (used ONLY by the batch — verified: ports/fee_schedule.go + adapters/feeschedule/adapter.go have no other consumer).

## Units (CRITICAL)
- Resolver `Comissao.Valor` = percent string ("13.00" means 13%). Batch commission math is a **fraction**: `commissionAmt = sellingPrice * fraction`. So `fraction = parseFloat(Valor) / 100`.
- Single-path decompose computes Comissao = round2(pct/100 × preco). Parity is on the resolved **pct**, both paths deriving it from the same resolver value.

## WRITE-SET (touch only these)
1. apps/server_core/internal/modules/pricing/application/batch_orchestrator.go
2. apps/server_core/internal/modules/pricing/application/batch_orchestrator_test.go
3. apps/server_core/internal/composition/root.go  (batch wiring only, ~lines 714-751)
4. apps/server_core/internal/modules/pricing/transport/http_handler.go  (handleBatch: set default Modalidade only)
5. DELETE apps/server_core/internal/modules/pricing/ports/fee_schedule.go
6. DELETE apps/server_core/internal/modules/pricing/adapters/feeschedule/adapter.go
7. contracts/api/marketplace-central.openapi.yaml  (BatchSimulationItem: add commission_source — additive)
8. packages/sdk-runtime/src/index.ts  (BatchSimulationItem interface: add commission_source)

## IMPLEMENTATION (exact)

### File 1 — batch_orchestrator.go
- Add import `"strconv"` and `"strings"`; keep `"context"`,`"fmt"`; add the pricing `domain` import `"marketplace-central/apps/server_core/internal/modules/pricing/domain"`.
- `BatchRunRequest`: add field `Modalidade domain.Modalidade`.
- `BatchSimulationItem`: add field `CommissionSource string \`json:"commission_source"\`` (place after `Status`, before `FreightSource` is fine).
- `BatchOrchestrator` struct: REMOVE `feeLookup ports.FeeScheduleLookup`; ADD `tariffResolver ports.TariffResolver`.
- `NewBatchOrchestrator`: REMOVE the `feeLookup ports.FeeScheduleLookup` param; drop it from the struct literal. New signature:
  `func NewBatchOrchestrator(products ports.ProductProvider, policies ports.PolicyProvider, freight ports.FreightQuoter, tenantID string) *BatchOrchestrator`
  Update the doc comment (no more feeLookup wording; note tariffResolver is nil-safe and set via WithTariffResolver).
- ADD method (mirror CalcService.WithTariffResolver, but pointer-receiver mutate+return since BatchOrchestrator is a pointer type):
  ```
  // WithTariffResolver attaches the shared ports.TariffResolver so batch
  // commission is resolved by the ML listing category (the same seam the
  // single-product decompose path uses), not the ERP taxonomy. Nil-safe:
  // without it the orchestrator uses the policy's commission_percent.
  func (o *BatchOrchestrator) WithTariffResolver(r ports.TariffResolver) *BatchOrchestrator {
      o.tariffResolver = r
      return o
  }
  ```
- Replace the commission block (current lines 180-191) with:
  ```
  commissionPct := pol.CommissionPercent
  commissionSource := "policy"
  commissionKnown := true
  switch {
  case pol.CommissionOverride != nil:
      commissionPct = *pol.CommissionOverride
      commissionSource = "override"
  case o.tariffResolver != nil:
      var productID *int
      if id, perr := strconv.Atoi(prod.ProductID); perr == nil {
          productID = &id
      }
      priceBasis := strconv.FormatFloat(sellingPrice, 'f', 2, 64)
      res, rerr := o.tariffResolver.Resolve(ctx, ports.TariffRequest{
          Modalidade: req.Modalidade,
          ProductID:  productID,
          PriceBasis: &priceBasis,
      })
      if rerr != nil {
          return BatchRunResult{}, fmt.Errorf("PRICING_BATCH_RESOLVE_TARIFF: %w", rerr)
      }
      if frac, ok := commissionFraction(res.Comissao.Valor); ok {
          commissionPct = frac
          commissionSource = "resolver"
      } else {
          // ADR-17: resolved but no commission datum → honest unknown, never a default %.
          commissionKnown = false
          commissionSource = "unknown"
      }
  }
  ```
- Replace the margin/status/append tail so an unknown commission is NOT fabricated:
  ```
  var commissionAmt, marginAmt, marginPct float64
  status := ""
  if !commissionKnown {
      status = statusCritical
  } else {
      commissionAmt = sellingPrice * commissionPct
      marginAmt = sellingPrice - prod.CostAmount - commissionAmt - pol.FixedFeeAmount - freightAmt
      if sellingPrice > 0 {
          marginPct = marginAmt / sellingPrice
      }
      status = simulationStatusForBatch(marginPct, freightAvailable)
  }
  ```
  (Use the existing `statusCritical` const from status.go — same package. When unknown, commissionAmt/marginAmt/marginPct stay 0 and CommissionSource="unknown" + Status="critical" flag the row untrustworthy — mirrors the existing freight-unavailable idiom; the "0" is a flagged sentinel, never a silent default.)
- Set `CommissionSource: commissionSource` in the appended BatchSimulationItem.
- ADD helper near the bottom of the file:
  ```
  // commissionFraction converts a resolver percent string ("13.00") to a
  // fraction (0.13). A nil/empty/unparseable value is NO DATA (ADR-17), not 0.
  func commissionFraction(valor *string) (float64, bool) {
      if valor == nil {
          return 0, false
      }
      s := strings.TrimSpace(*valor)
      if s == "" {
          return 0, false
      }
      pct, err := strconv.ParseFloat(s, 64)
      if err != nil {
          return 0, false
      }
      return pct / 100, true
  }
  ```
- Remove the now-unused `ports` import ONLY if nothing else in the file uses `ports.` — it still does (ports.ProductProvider, ports.TariffRequest, etc.), so KEEP it.

### File 3 — composition/root.go (batch wiring, ~714-751)
- Line 715: DELETE `feeAdapter := pricingfee.NewAdapter(feeSvc)`.
- Line 718: change to `batchOrch := pricingapp.NewBatchOrchestrator(prodReader, polReader, meClient, cfg.DefaultTenantID)` (drop the feeAdapter arg).
- Inside `if internalReadAvailable { ... }`, AFTER `tariffResolver` is fully built (the `var tariffResolver ... ` / composite block ending ~line 746-747, BEFORE `calcSvc := ...`), add: `batchOrch.WithTariffResolver(tariffResolver)` — batchOrch is the same pointer already held by pricingHandler, so this takes effect.
- Remove the `pricingfee` import (the feeschedule adapter alias) if now unused. Do NOT remove `feeSvc` (still used at line ~704 by marketplacestransport).
- Leave everything else in root.go untouched.

### File 4 — transport/http_handler.go (handleBatch)
- The wire request has NO modalidade field and this slice does NOT add one (keeps the OpenAPI request contract unchanged). In the `h.batch.RunBatch(...)` call inside handleBatch, add `Modalidade: domain.ModalidadeClassico,` to the BatchRunRequest literal (classico = the default ML modalidade for the batch matrix).
- Add import `"marketplace-central/apps/server_core/internal/modules/pricing/domain"`.

### File 5 & 6 — DELETE
- Delete apps/server_core/internal/modules/pricing/ports/fee_schedule.go
- Delete apps/server_core/internal/modules/pricing/adapters/feeschedule/adapter.go
(Verified no other consumer. If `go build ./...` finds one, STOP and report — do not chase it.)

### File 7 — OpenAPI (contracts/api/marketplace-central.openapi.yaml, BatchSimulationItem ~3932)
- In `required:` add `- commission_source`.
- In `properties:` add:
  ```
        commission_source:
          type: string
          enum: [resolver, override, policy, unknown]
  ```

### File 8 — sdk-runtime (packages/sdk-runtime/src/index.ts, BatchSimulationItem ~1413)
- Add `  commission_source: string;` to the interface (after `status: string;`).

### File 2 — batch_orchestrator_test.go (TDD — write these to fail first, then implement)
Rewrite the 4 existing commission tests to the resolver seam + add unknown + parity. Provide:
- A stub `ports.TariffResolver`:
  ```
  type stubTariffResolver struct {
      valor *string // Comissao.Valor to return; nil = unknown
      err   error
  }
  func (s *stubTariffResolver) Resolve(_ context.Context, _ pricingports.TariffRequest) (pricingdomain.TariffResolution, error) {
      if s.err != nil { return pricingdomain.TariffResolution{}, s.err }
      return pricingdomain.TariffResolution{Comissao: pricingdomain.ComponentResolution{Valor: s.valor}}, nil
  }
  ```
  (import pricingdomain "marketplace-central/apps/server_core/internal/modules/pricing/domain")
- Remove `stubFeeScheduleLookup` (no longer used).
- Helper `runBatch` now takes a `*stubTariffResolver` (nil = not wired) instead of a feeLookup; product price 200, cost 100 (as today). When resolver non-nil, wire via `.WithTariffResolver(resolver)`. Return the whole first item (so tests can read CommissionAmount + CommissionSource + Status).
- Tests (all fraction/percent explicit):
  1. `TestBatchOrchestrator_CommissionOverrideTakesPriority` — override 0.05 wins over resolver "99"; expect CommissionAmount==10 (200*0.05), CommissionSource=="override".
  2. `TestBatchOrchestrator_ResolverUsedWhenNoOverride` — resolver Valor="20.00", no override; expect CommissionAmount==40 (200*0.20), CommissionSource=="resolver".
  3. `TestBatchOrchestrator_PolicyRateUsedWhenNoResolver` — resolver nil (not wired), policy 0.10; expect CommissionAmount==20, CommissionSource=="policy".
  4. `TestBatchOrchestrator_UnresolvedCategoryIsHonestUnknown` — resolver Valor=nil; expect CommissionSource=="unknown", Status=="critical", CommissionAmount==0 (NOT the policy 0.10 default) — assert it did NOT fabricate 20.
  5. `TestBatchOrchestrator_ParityWithDecompose` (EXEMPLO-IO) — a SHARED stub resolver returns Valor="13.00". Batch: product price 200, no override, Modalidade classico → CommissionAmount==26 (200*0.13). Single path: build `application.NewCalcService(stubCalcRepo, stubCost, nil, "tenant_default").WithTariffResolver(sameStub)` and call Decompose with `DecomposeRequest{Preco:"200", ComissaoPct:"", Modalidade: pricingdomain.ModalidadeClassico, Custo: ptrStr("100")}`. Assert `result.Decomposition.Comissao == "26.00"` and that it equals the batch commission (26). Minimal stubCalcRepo.GetProfile returns `pricingdomain.CalcProfile{AliquotaPct:"0", DifalEnabled:false}`; stubCost unused because Custo is supplied (ProductID nil). Add a comment: in production the shared resolver is the composite ML-category resolver; the test proves both paths key commission off the same seam identically.
  - Helpers: `ptrStr(s string) *string { return &s }`. You may need a minimal `stubCalcRepo` implementing `ports.CalcRepository` (only GetProfile is exercised for this test — other methods can return zero values/nil; check the interface in ports and stub all methods). If CalcRepository is large, INSTEAD assert parity purely at the commission level without full CalcService: compute decompose commission via `pricingdomain.Decompose(pricingdomain.DecomposeInput{Preco:"200", ComissaoPct:"13.00", AliquotaPct:"0", Modalidade:ModalidadeClassico, Custo:&Money{Amount:"100",Currency:"BRL"}})` and assert its `.Comissao=="26.00"` equals batch's 26 — this proves the shared resolver value ("13.00") produces the SAME commission in the pure domain engine that decompose uses. PREFER this pure-domain form if stubbing CalcRepository is heavy; it still proves parity of the commission math off the shared resolver value. State which form you used.

## VALIDATION PLAN (run from apps/server_core unless noted; report evidence type per command)
- `GOCACHE=.gocache go build ./...`
- `GOCACHE=.gocache go vet ./internal/modules/pricing/... ./internal/composition/...`
- `GOCACHE=.gocache go test ./internal/modules/pricing/...`  (must be green incl the 5 tests above)
- sdk: from repo root `npx tsc --noEmit -p packages/sdk-runtime/tsconfig.json` (or `npm run -w packages/sdk-runtime build` if present). If the sandbox blocks vite/esbuild, report could-not-run — the chip re-runs it.
- Report the `go test` output tail verbatim.

## G-notes to answer in your report
- G1: confirm the resolver seam reuse matches decompose (cite calc_service.go:403 + tariffcomposite). 
- Confirm commission_source added to BOTH OpenAPI + sdk-runtime in this same change (contract atomicity).
- List changed paths vs this write-set.
