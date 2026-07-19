# CHIP-T1 pricing-solver-tier4 — EVIDENCE

`branch chip/t1-solver-tier4` · base `18fbd91a` · scope DESIGN-TARIFAS-ML §8.1 DEMO-CRÍTICO + §4 wiring + §2.5 degrau 4.
Contingency lane §12 (codex quota dead til 2026-07-25): Claude-only — sonnet workers implement, chip plans + P5-verifies each slice, cold-Opus + adversarial-sonnet dual gate at close.

## Commits (one per green slice)

| slice | sha | summary |
|-------|-----|---------|
| A | `a7a5bc9` | migration 0068 pricing_tariff_defaults + domain TariffDefaults + ports.TariffDefaultsStore + repo Get/Upsert + CalcService Get/Put + GET/PUT /pricing/tariff-defaults (Go) + runner_test 55→56 |
| B | `7bfcdc0` | domain TariffResolution/ComponentResolution/Fonte + ports.TariffResolver + tariffdefaults degrau-4 adapter |
| C | `065e06e` | decomposeWithLimiar (frozen Decompose sig intact) + SolveInput.TaxaFixaLimiarCents + SolveResult.FreteDesconhecido; removed R$150 probe + const thresholdCents; segment-conditional frete; 3 new goldens |
| D | `b8bf522` | CalcService.WithTariffResolver + resolveTariff (comissao_pct optional, manual override=FonteManual, sem_dados frete stays nil); DecomposeResult/SolveOutput.Tarifa |
| F | `3f0c309` | root.go additive: WithTariffStore(calcRepo) + WithTariffResolver(tariffdefaults.NewResolver(calcRepo, tenant, "")) inside internalReadAvailable branch |
| E | `c57de66` | transport `code` branch (solveCode) + `tarifa` block + `frete_desconhecido` + OpenAPI pricing delta + sdk-runtime — OpenAPI & SDK in the SAME commit (contract-lock ACK'd by hub) |

## ADR-17 (unknown ≠ zero) conformance
- `frete_estimativa_amount` NULLABLE numeric (migration 0068); NULL never coerced to 0.
- `13.00`/`16.00` commissions are DB column DEFAULTs only — verified NO `13`/`16` Go literal (only doc comments). Materialize-on-read via `INSERT ON CONFLICT DO NOTHING` then `::text` SELECT (mirrors 0055).
- Resolver frete: `sem_dados` (or estimativa with nil amount) → `ComponentResolution.Valor == nil` (NO-DATA), never `"0"`.
- Solver: `FreteProduto == nil` in the ≥limiar segment → `SolveResult.FreteDesconhecido = true` (distinct from structural `Desconhecidos` and from `CeilingPct`); the <limiar segment solves regardless.

## Frozen IC-04 contract (shared with M-08) — INTACT
`ports/calc_ports_contract_test.go`: `TestDecompositionShapeFrozen`, `TestDecomposeInputShapeFrozen`, `TestDifalForUFResultShapeFrozen` all PASS. Threshold threaded via unexported `decomposeWithLimiar` + M-07-only `SolveInput`/`SolveResult` extensions — no field added to `DecomposeInput`/`Decomposition`, `Decompose(DecomposeInput) Decomposition` signature unchanged.

## Gate ladder (P5 self-verify, GOCACHE/GOMODCACHE absolute, no GOFLAGS)
- L0 `go build ./...` → BUILD_OK. `go vet ./internal/modules/pricing/... ./internal/composition/...` → VET_OK.
- L1 `go test ./internal/modules/pricing/... ./internal/platform/migrate/...` → all `ok` (incl. new tariffdefaults, tariff_defaults domain, 3 new solve goldens, 5 new calc_service resolver tests, migrate 56-count, frozen ports contract).
- Each slice independently re-verified by the chip (not only the worker's claim); domain crux (Slice C) diff read line-by-line.

## Deferred / notes
- `installation_id` uses `""` default-installation sentinel (cfg has no installation field); schema is multi-installation-ready (PK tenant_id+installation_id), single installation wired for the demo.
- Per-solve `TaxaFixaLimiarCents` defaults to policy (7900) for the demo — threshold is now a domain parameter, not a literal (§8.1.3).
- **Latent design note (round-3 REFUTE, NOT a shipped defect):** the solver's ground truth is `decomposeWithLimiar(…, limiar)` at whatever limiar it solved under; the exported frozen `Decompose()` hardcodes 7900. These agree only at the default limiar. If a future caller sets `TaxaFixaLimiarCents` to a non-default value, a solved price would re-`Decompose()` (at 7900) to a different margem_pct than it solved for. Not reachable today — every live caller leaves `TaxaFixaLimiarCents=0` → default 7900 (`calc_service.go:352`), which equals the exported `Decompose` threshold, so the re-decompose invariant holds on the demo path. Whoever wires a real per-tenant/per-policy limiar override must thread that limiar through the ground-truth re-decompose too (the frozen `Decompose` signature can't take a limiar param — IC-04). Amendment candidate for the profile.
- Repo-level SQL for tariff_defaults exercised only by build/vet + migrate-count (calc_repository_test.go is `//go:build integration`, needs live PG); handler contract covered by stub-backed transport tests.

## Released contract shape (Slice E — for hub COMMITTED)
`tarifa` (nullable object, present on decompose + solve):
```
tarifa: {
  comissao: { valor: string|null, fonte: "PADRAO"|"MANUAL", degrau: int, data: string|null, estimativa: bool } | null,
  frete:    { valor: string|null, fonte: "PADRAO"|"MANUAL", degrau: int, data: string|null, estimativa: bool, sem_dados: bool } | null
}
```
Solve response adds `frete_desconhecido: bool` (required). `code` enum = `UNREACHABLE_TARGET | SEM_FRETE | DADOS_INCOMPLETOS` (empty string = solved). `PricingCalcInput.required` narrowed to `[modalidade]` (comissao_pct now optional — hub-approved widening).
- **DEVIATION FLAGGED:** `tarifa.comissao.data` / `tarifa.frete.data` are ALWAYS `null` at degrau 4. Design §4.4 stamps a `data` timestamp; at the padrão degrau there is no policy/source time to stamp, and fabricating `now()` would violate ADR-17 (unknown ≠ invented). Field is PRESENT (nullable) so FE tolerant parsing binds; value stays honest-null until a manual/dated source lands (degrau ≤3). Reported in COMMITTED.

## Close gates (all green)
- L0/L1 re-run after Slice E: `go build ./...` BUILD_OK · `go vet` VET_OK · `go test ./internal/modules/pricing/... ./internal/platform/migrate/... ./internal/composition/...` all `ok`.
- Governance: `pwsh scripts/harness.ps1 -Command governance -BaseSha 18fbd91a5706f8a58f41ddf27f13c8f58a5dd6e4` from clean worktree → **status=passed** (baseline_exceptions all pre-existing).
- SDK: `tsc` → TSC_OK. OpenAPI: PyYAML parse OK.
- Branch diffstat: 6 commits (a7a5bc9·7bfcdc0·065e06e·b8bf522·3f0c309·c57de66), ~1370 insertions, 20 files.

## P6 dual gate — round 1 (found a blocker)
- Cold Opus (ad6fc0f3): SHIP-WITH-NITS, invariants 1–7 confirmed, 3 nits.
- Adversarial sonnet REFUTE (a47645ed): **REFUTED — BLOCKER FINDING-P6-SOLVER.** Disagreement ⇒ gate not passed.

### FINDING-P6-SOLVER (BLOCKER, fixed @865863e) — for hub ratification
`Decompose` rounds comissão/imposto/difal to 2dp **before** subtracting, so its 2dp `margem_pct` is only *piecewise*-increasing in preço — it has sub-cent downward dips (601 dips at com12/ali0, 897 at com13/ali4, 1071 at com16/ali12 — ordinary tariffs, not edge values). The old per-segment **binary search assumed strict monotonicity**, so it skipped valid prices and mis-reported reachable targets. Reproduced against production code (independently re-run by the chip): com=1%/ali=0%/custo=10, target 38.98 has an exact solution at R$27,50, but `SolveTargetPrice` returned `FreteDesconhecido=true` (→ HTTP 200 `SEM_FRETE`) with frete nil, and `Ceiling=99.00` (→ `UNREACHABLE_TARGET`) with frete=15,00. Both flatly wrong; reachable via public `/pricing/solve` with ordinary inputs. **Demo-critical.**

**Fix (@865863e):** low segment scanned EXHAUSTIVELY (bounded by limiar, cmpPct match); high segment bracketed by the EXACT unrounded margem_pct `100·k − 100·F/preço` (strictly increasing) via `firstCentExactAtLeast`, then linear-scanned within a `2/preço+0.02` pp wiggle window (capped `highScanCapCents=4e6`) against the real `Decompose`. Replaced the false monotonic-invariant test (passed only by coarse sampling that never landed on a dip — test theater) with `TestSolveMatchesBruteForceLowSegment` (brute cheapest-exact-match equivalence over dip grids) + `TestSolveHighSegmentExactAndCheapest`. Post-fix probe: target 38.98 now `Reached=true Preco=27.50` for both frete nil and frete=15.

### Cold-Opus nits addressed (@0f5537a)
- nit 1: `PutTariffDefaults` now guards the numeric fields → 422 `INVALID_PRICE` (was 500) + regression test.
- nit 2 (latent, both reviewers): `structuralUnknowns` floors the probe at 1 cent → no divide-by-preço-0 if `TaxaFixaLimiarCents=1`.
- nit 3: OpenAPI `tarifa` nullable refs wrapped in `allOf` (repo idiom; strict-3.0 keeps nullability). SDK already `PricingTarifa | null` — unchanged.

## Re-verify after fix (all green)
- `go build ./...` BUILD_OK · `go vet` VET_OK · `go test ./internal/modules/pricing/... ./internal/platform/migrate/... ./internal/composition/...` all `ok` (domain suite ~17s incl. brute-force test).
- Governance (clean chip worktree, BaseSha 18fbd91a) → **status=passed**.
- `tsc` TSC_OK · OpenAPI YAML_OK.
- Range: 18fbd91a..HEAD = 8 commits (…c57de66, **865863e**, **0f5537a**), 20 files, ~1558 insertions.

## P6 dual gate — round 2 (re-gate on the @865863e/@0f5537a branch)
- Cold Opus (a2c5bc56): SHIP-WITH-NITS — classified the high-segment scan cap as NIT N1 (outside realistic envelope, degrades honestly) + N2 doc/behavior mismatch (R$1M reach vs cap), N3 test-sampling nondeterminism, N4 OpenAPI 422 description narrow.
- Adversarial REFUTE (ae3f42d1): **REFUTED — cap defect classified BLOCKER.** Reviewers disagreed on severity of the SAME issue ⇒ gate not passed; chip eliminated the issue entirely rather than adjudicate severity.

### FINDING-P6-SOLVER-2 (residual high-segment defect, fixed @40be57a) — for hub ratification
The round-1 fix bracketed the high segment with a `2/preço+0.02` pp wiggle window then **capped the linear scan at `highScanCapCents=4e6`**. Two failure modes, both independently reproduced by the chip:
- **False UNREACHABLE:** a reachable near-ceiling target whose cheapest exact price sits beyond the cap. CONFIRMED: com=10/ali=0/frete=0/custo=5000, target **89.40** has an exact match at **R$826.445,64** (`Decompose` = 89.40 there), but the capped solver returned UNREACHABLE (Ceiling 90.00). Reachable via public `/pricing/solve`.
- **Slow solve:** mid-range near-ceiling targets scanned 20–55s (non-cancellable).

**Fix (@40be57a):** replace wiggle-window+cap with an EXACT bracket + BOUNDED window. `round2(exact.pct) ≥ target ⇔ exact.pct ≥ target−0.005`, so bisect `firstCentExactAtLeast(target−0.005)` to the crossing `cStar`, then scan a fixed **±highSegmentWindowCents (20 000)** window against the real `Decompose`. The real match differs from `cStar` only by the 2dp component-rounding perturbation: `|Δpreço| ≤ 150/(ceiling−target)` cents, and `ceiling−target ≥ 0.01` (targets at/above the analytic ceiling are rejected by `SolveTargetPrice` before any search) ⇒ ≤ 15 000 cents < 20 000. Search now reaches the full `solveMaxCents` (R$1M) in O(log)+bounded scan — no false UNREACHABLE, <1s. Post-fix probe: target 89.40 → **R$826.445,64** (0.89s); 89.00 → R$497.512,04 (cheaper than the R$500k witness — cheapest honored); 89.50 → R$990.098,24. Probe deleted (no cruft committed).

**Cold-Opus round-2 nits also cleared @40be57a:** N2 resolved by the same fix (reach is genuinely full R$1M now); N3 → deterministic cent-stride sampling in `TestSolveMatchesBruteForceLowSegment`; N4 → PUT `/pricing/tariff-defaults` 422 description widened to cover non-numeric commission/frete (`INVALID_PRICE`). SDK types unchanged (description-only).

## Re-verify after cap fix (all green)
- `go build ./...` BUILD_OK · `go vet ./internal/modules/pricing/...` VET_OK.
- `go test ./internal/modules/pricing/... ./internal/platform/migrate/... ./internal/composition/...` all `ok` (domain ~32s incl. brute-force map build).
- Governance (clean chip worktree, BaseSha 18fbd91a5706f8a58f41ddf27f13c8f58a5dd6e4) → **status=passed** (only pre-existing baseline exceptions; no new `production-panic` site — `mustRat` was already used in solve.go).
- `tsc` TSC_OK · OpenAPI YAML_OK.
- Range: `18fbd91a..HEAD` = 9 commits (…c57de66, 865863e, 0f5537a, **40be57a**).

## P6 dual gate — round 3 (re-gate on @40be57a) — PASSED (both ship)
- Adversarial REFUTE: **COULD NOT REFUTE.** Brute-oracle vs solver over ~240 near-ceiling worst-case fuzz cases (gap 0.005–0.05) + a 200-iter gap=0.005 fuzz + low-segment brute — all exact cheapest, no wrong/non-cheapest price. Worst observed displacement below the crossing ~11 692 cents (< the 20 000 window). Latent non-shipped note: custom `TaxaFixaLimiarCents` would diverge internal `margemDecompose` from the frozen exported `Decompose` — unreachable today (see Deferred).
- Cold Opus: **SHIP-WITH-NITS.** IC-04 frozen contract, segment split, ADR-17, OpenAPI/SDK parity all confirmed correct. 4 nits, all non-blocking; N1 the substantive one.

Both reviewers ship, no BLOCKER ⇒ **gate passed at round 3.** All 4 cold-Opus nits then resolved @81d61df (correctness-hardening, algorithm shape unchanged — the same shape both reviewers validated — regression-covered; no round-4 re-gate warranted for non-blocking-nit resolution):
- **N1 (soundness, eliminated @81d61df):** the ±20 000 fixed window's justification assumed `ceiling−target ≥ 0.01`, but the reachability gate compares against the 2dp-**rounded** ceiling, so a >2dp comissão (arbitrary precision at `calc_service.go:410`) yields a 3-decimal exact ceiling and a sub-0.01 gap. Replaced the constant with a window **derived per-input from the exact gap**: `cStar = firstCentExactAtLeast(target−0.005)` is the `round2(exact)=target` crossing; the cheapest Decompose match lies ≤ `150/gStar` cents from it, `gStar = exactCeiling − (target−0.005)`. The −0.005 (final-round2 half-band) is folded into `cStar`, so `gStar ≥ 0.005` for any input and ≥ 0.01 for every **reachable** (2dp) target — provably finite/sufficient for ALL inputs, never a magic cap. Added `exactCeilingRat()` + `ceilRatToCents()`. (Note: this also shows N1's stated `150/g=30000` was conservative — measured from `cStar` the bound is `150/gStar ≤ 15000` for reachable targets, so the old constant was in fact safe; the fix makes the argument sound rather than empirical.)
- **N2 (custom-limiar low-segment invariant) — dissolved:** the high-segment bracket path is now sound regardless of segment width, so a `TaxaFixaLimiarCents` override no longer leans on the low-segment exhaustive-scan invariant.
- **N3 (weak high-segment test) — addressed:** `TestSolveHighSegmentNearCeilingCheapest` exercises the sub-0.01 ceiling-gap regime (comissão 12.005 → exactCeiling 87.995) with a brute cheapest oracle below the solved price, proving the window is not truncated.
- **N4 (manual comissão >2dp) — sound by design:** rather than reject >2dp commission (a legitimate manual override), the solver is now provably correct for ANY comissão precision via the derived window.

## Re-verify after N1 fix (all green)
- `go build ./...` BUILD_OK · `go vet` VET_OK.
- `go test ./internal/modules/pricing/... ./internal/platform/migrate/... ./internal/composition/...` all `ok` (domain ~39s incl. brute-force + near-ceiling oracles).
- Governance (clean chip worktree, BaseSha 18fbd91a5706f8a58f41ddf27f13c8f58a5dd6e4) → **status=passed**.
- `tsc` TSC_OK · OpenAPI YAML_OK.
- Range: `18fbd91a..HEAD` = 10 commits (…c57de66, 865863e, 0f5537a, 40be57a, **81d61df**). Post-fix probe (deleted): target 89.40 → R$826.445,64; near-ceiling comissão 12.005 targets 87.90/87.98/87.99 all exact-cheapest.
