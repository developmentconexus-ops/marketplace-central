# EVIDENCE — CHIP-FIX4-BATCH-COMMISSION

**Chip:** CHIP-FIX4-BATCH-COMMISSION · **Base SHA:** 97ef7b09e683d601fd74e02aa604d725378f44de
**Contract:** docs/architecture/domain-model-marketplace.md §7 FIX-4
**Status:** CLOSED — P6 dual gate AGREEMENT (see §5). **Commit:** ba8b25b7 (branch claude/brave-panini-30631c, not pushed).

> WORKTREE FINDING (report to hub): the dispatched worktree
> `.claude/worktrees/sleepy-wing-6d7500` is empty (0 files); its git-dir points at the
> main `.git`. `git worktree list` shows the main checkout
> `C:/Users/leandro.theodoro/Documents/marketplace-central` is itself on branch
> `claude/brave-panini-30631c` at the base SHA. All edits/tests/commit were therefore
> done in the main checkout (which physically holds this chip branch), via absolute
> paths. HEAD there == base SHA at pack time (no commit yet). Hub should ratify whether
> chip worktrees are being provisioned empty (harness upstream candidate).

---

## 1. Deliverables (write-set — 8 files, exact)

| # | File | Change |
|---|---|---|
| 1 | apps/server_core/internal/modules/pricing/application/batch_orchestrator.go | commission resolved via `ports.TariffResolver` (ML category), not `prod.CategoryID`; `WithTariffResolver`; honest-unknown; `round2` (FIX-4 parity); `commissionFraction`; `CommissionSource` field; `Modalidade` on request |
| 2 | apps/server_core/internal/modules/pricing/application/batch_orchestrator_test.go | 4 commission tests rewritten to resolver seam + unknown case + rounding-sensitive parity test (EXEMPLO-IO) |
| 3 | apps/server_core/internal/composition/root.go | drop feeAdapter; `batchOrch.WithTariffResolver(tariffResolver)` — same instance CalcService gets |
| 4 | apps/server_core/internal/modules/pricing/transport/http_handler.go | handleBatch sets `Modalidade: domain.ModalidadeClassico` |
| 5 | apps/server_core/internal/modules/pricing/ports/fee_schedule.go | DELETED (only-batch consumer) |
| 6 | apps/server_core/internal/modules/pricing/adapters/feeschedule/adapter.go | DELETED |
| 7 | contracts/api/marketplace-central.openapi.yaml | BatchSimulationItem: `commission_source` additive (required + enum resolver/override/policy/unknown) |
| 8 | packages/sdk-runtime/src/index.ts | BatchSimulationItem: `commission_source: string` |

`git diff --stat`: 8 files, +214 / −101. Nothing outside the write-set.

## 2. The bug → the fix (before/after)

**BEFORE** (batch_orchestrator.go, commission block): keyed commission by
`categoryID := prod.CategoryID` (ERP TaxonomyNodeID via adapters/catalog/reader.go:37)
→ `o.feeLookup.LookupFee(marketplaceCode, categoryID, "")`. ERP taxonomy ≠ ML listing
category → wrong commission %; violates FIX-4.

**AFTER**: commission resolved through the SAME seam the single-product /precos decompose
path uses — `o.tariffResolver.Resolve(ctx, ports.TariffRequest{Modalidade, ProductID, PriceBasis})`
→ composite resolver (degrau-3 live ML `ResolveCategory(EAN,Titulo)` over degrau-4 config).
`Comissao.Valor` (percent string) → `commissionFraction` (÷100) → `round2(price×frac)`.
Parity with decompose by construction (same resolver instance, same value, same rounding).

**Resolver seam reused (not reinvented):** batch call `ports.TariffRequest` shape ==
`CalcService.resolveTariff` (calc_service.go). In root.go the single local `tariffResolver`
is passed to BOTH `batchOrch.WithTariffResolver(tariffResolver)` and
`calcSvc…WithTariffResolver(tariffResolver)` — one instance, true parity.

## 3. Honest-unknown (ADR-17)

Resolver wired but `Comissao.Valor == nil` (unresolvable ML category) →
`CommissionSource="unknown"`, `Status="critical"`, `CommissionAmount/MarginAmount/MarginPercent`
left 0 as flagged sentinel — NEVER computed from the policy/default %. Resolver infra
error PROPAGATES (`return …, fmt.Errorf("PRICING_BATCH_RESOLVE_TARIFF: %w", rerr)`) — not
swallowed into unknown. No blanket recover/fallback on the integrity read.

**One-line honest-unknown statement:** an unresolvable ML listing category is surfaced as
`commission_source="unknown"` + `status="critical"` with a zeroed sentinel amount, never a
fabricated default commission.

## 4. Tests (GOCACHE=.gocache, from apps/server_core)

- `go build ./...` → OK (exit 0)
- `go vet ./internal/modules/pricing/... ./internal/composition/...` → OK (exit 0)
- `go test ./internal/modules/pricing/...` → all packages `ok`
- `npx tsc --noEmit -p packages/sdk-runtime/tsconfig.json` → clean (exit 0)

```
=== RUN   TestBatchOrchestrator_CommissionOverrideTakesPriority
--- PASS: TestBatchOrchestrator_CommissionOverrideTakesPriority (0.00s)
=== RUN   TestBatchOrchestrator_ResolverUsedWhenNoOverride
--- PASS: TestBatchOrchestrator_ResolverUsedWhenNoOverride (0.00s)
=== RUN   TestBatchOrchestrator_PolicyRateUsedWhenNoResolver
--- PASS: TestBatchOrchestrator_PolicyRateUsedWhenNoResolver (0.00s)
=== RUN   TestBatchOrchestrator_UnresolvedCategoryIsHonestUnknown
--- PASS: TestBatchOrchestrator_UnresolvedCategoryIsHonestUnknown (0.00s)
=== RUN   TestBatchOrchestrator_ParityWithDecompose
--- PASS: TestBatchOrchestrator_ParityWithDecompose (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/pricing/application	1.875s
```

**EXEMPLO-IO parity test:** 2-case table, both rounding-sensitive, each asserting
batch `CommissionAmount` == `pricingdomain.Decompose(...).Comissao` == expected:
- 78.90 @ "12.00" → 9.47 (non-boundary)
- 6.10 @ "15" → 0.92 (true .xx5 half-cent boundary; the exact defect the R2 fix closed —
  float `round2` gave 0.91, exact-decimal gives 0.92)
Fails against both the pre-fix unrounded code AND the intermediate float-`round2` code — real
regression guard. Parity-test form = pure-domain `Decompose` (full `CalcRepository` stub too
heavy; proves commission-math parity off the shared resolver value). Resolver commission is
bit-for-bit equal to decompose (both route through `domain.ParseRat`/`FormatRatHalfUp`).

## 5. P6 DUAL GATE (contingency lane: cold Opus gate-reviewer + independent adversarial sonnet)

**Round 1:**
- Cold Opus gate-reviewer → **PASS-WITH-NOTES**. All Required criteria (seam reuse,
  ML-category keying, unit conversion, ADR-17 honest-unknown, contract atomicity) PASS with
  file:line evidence. Notes: stale `CategoryID` comment/reader residue outside write-set
  (follow-up); single-modalidade (classico) scope narrowing; parity test pure-domain form.
- Adversarial sonnet → **DEFECTIVE (minor)**. Landed: (a) rounding parity gap — batch
  commission unrounded float vs decompose big.Rat half-up (78.90@12% → 9.468 vs 9.47);
  (b) parity test used rounding-clean input, can't catch (a). Core keying / ADR-17 / seam /
  contract all confirmed SOUND.

→ **DISAGREEMENT** (no AGREEMENT). Fix R1: `round2` added; parity test → 78.90@12% rounding-sensitive.

**Round 2 (re-gate):** adversarial → **STILL-DEFECTIVE (minor)**. Defect #2 (test) resolved.
Residual #1 LANDS: `round2` is float64 money math; decompose uses exact `big.Rat`
(`FormatRatHalfUp`). Named input price 6.10 @ pct "15" → batch 0.91 vs decompose 0.92
(true 0.915 half-cent boundary); 6077 divergent pairs on a 0.01–2000 scan. Float drift
reintroduced exactly what `domain/decimal.go` exists to prevent.

→ Fix R2: resolver (parity-critical) commission routed through the SAME exact-decimal
primitives decompose uses — new `commissionDecimal(price, pctStr)` =
`domain.ParseRat` × (pct ÷ big.NewRat(100,1)) → `domain.FormatRatHalfUp(_, 2)` (decimal.go:30,46).
`commissionFraction` removed. `round2` retained ONLY on override/policy branches (float-sourced
batch-only config, never compared to decompose). Parity test → 2-case table: 78.90@"12.00"→9.47
(non-boundary) + 6.10@"15"→0.92 (boundary), each asserting batch==decompose==expected.

**Round 3 (re-gate):** adversarial → **SOUND**. Brute-force re-scan of the NEW path: 300,000
prices (0.01–3000) × 26 rates (incl. every .xx5-producing rate) → **0 divergences**. 6.10@15%
now 0.92==0.92. All prior SOUND findings (ADR-17, seam reuse, contract atomicity, clean removal)
re-confirmed intact.

**Gate verdicts:**
- Cold Opus gate-reviewer: **PASS** (all Required criteria pass; notes = non-blocking follow-ups §7).
- Adversarial sonnet refuter (R3): **SOUND** (parity exact; 0 divergences; no regression).

P6-DUAL-GATE: AGREEMENT

## 6. Not verified

- No server boot / live HTTP `/pricing/simulations/batch` exercise (out of chip scope).
- Live ML-category resolution not re-exercised (this change re-wires an already-live-verified
  resolver instance; no new provider-touching code — all new tests use stubs). No LIVE-VERIFIED
  marker claimed for new provider I/O because none was added.

## 7. Follow-ups for hub (findings, not blockers)

1. Empty chip worktree provisioning (see top).
2. Stale dead-code residue: `ports/batch_ports.go` `CategoryID` field + comment and
   `adapters/catalog/reader.go:37` assignment now describe the deleted fee-schedule seam;
   `CategoryID` is populated but no longer read for commission. Outside write-set → sweep in
   a follow-up chip.
3. Batch parity verified for modalidade=classico only (batch DTO carries no modalidade).
   Contract's "(product, modalidade)" parity holds for classico; multi-modalidade batch =
   future scope.
