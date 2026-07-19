# CHIP-T1 pricing-solver-tier4 — dispatch ledger

Contingency lane §12 (codex quota dead til 2026-07-25): implementers = sonnet Agent-tool workers; chip authors each slice card, P5-verifies independently (re-runs gates, reads crux diffs), commits per green slice.

| # | slice | worker | model | outcome | chip P5 re-verify | commit |
|---|-------|--------|-------|---------|-------------------|--------|
| 1 | A tariff-defaults stack | general-purpose | sonnet | green (integration repo test correctly skipped — no live PG; covered by stub transport tests) | re-ran build/vet/test; scanned for 13/16 Go literals (none, comments only); migration ADR-17 shape confirmed | `a7a5bc9` |
| 2 | B TariffResolver + adapter | general-purpose | sonnet | green (8 subtests) | re-ran build/vet/test; confirmed frete Valor nil default (no zero substitution) | `7bfcdc0` |
| 3 | C domain segment-conditional frete | general-purpose | sonnet | green (8 solve tests) | re-ran full pricing suite + isolated frozen ports contract (3 ShapeFrozen PASS); read solve.go+decompose.go diff line-by-line | `065e06e` |
| 4 | D service resolver injection | general-purpose | sonnet | green (5 new resolver tests); worker self-corrected golden target for high-segment frete case | re-ran build/vet/application tests; read resolveTariff impl (back-compat + FonteManual + sem_dados nil) | `b8bf522` |
| 5 | F root.go composition | chip inline (≤10-line additive glue) | opus | build/vet OK | direct build+vet of composition | `3f0c309` |
| 6 | E transport + OpenAPI + SDK | general-purpose | sonnet | green (tsc TSC_OK; PyYAML lint OK; transport tests pass) | contract-lock ACK'd by hub; re-ran build/vet/test pricing; read handler solveCode branch + toTarifaDTO; verified OpenAPI↔SDK field parity by hand | `c57de66` |

## Events emitted to hub (local_7c41fdd8-dd6f-487c-a7ea-1b511e214fa9)
- REQUEST contract-lock (pre-compaction, original) — non-contract slices proceeded per its statement.
- REQUEST contract-lock reminder + A–F status (this session) — awaiting ACK.

## P6 gate (chip-side, contingency lane §12: cold Opus + adversarial sonnet REFUTE)
| round | cold Opus | adversarial REFUTE | outcome |
|-------|-----------|--------------------|---------|
| 1 | ad6fc0f3 SHIP-WITH-NITS (inv 1–7 ✓) | a47645ed **REFUTED — FINDING-P6-SOLVER blocker** | disagree ⇒ NOT passed; chip reproduced blocker, fixed @865863e; nits @0f5537a |
| 2 (re-gate) | a2c5bc56 SHIP-WITH-NITS (cap=N1) | ae3f42d1 **REFUTED — cap=blocker** | disagree on SAME cap issue's severity ⇒ NOT passed; chip eliminated it (FINDING-P6-SOLVER-2) via exact bracket + bounded ±20k window @40be57a; also cleared N2/N3/N4 |
| 3 (re-gate) | SHIP-WITH-NITS (IC-04/ADR-17/parity ✓; 4 nits, N1 soundness) | **COULD NOT REFUTE** (~240 near-ceiling + 200-iter gap=0.005 fuzz, worst disp ~11.7k<20k) | **both ship, no BLOCKER ⇒ PASSED**; N1–N4 resolved @81d61df (derived-window fix, regression-covered) |

## Findings (for hub ratification)
- **FINDING-P6-SOLVER (BLOCKER, self-found & fixed @865863e):** solver binary search assumed margem_pct monotone; Decompose's 2dp component rounding makes it non-monotone (600–1000+ dips on ordinary tariffs) ⇒ reachable low-segment targets mis-reported as SEM_FRETE/UNREACHABLE. Fix = exhaustive low scan + exact-monotone high bracket. Amendment candidate: solver goldens must brute-force cheapest-exact-match, not sample a monotone assumption (the old test was theater). See EVIDENCE.md.
- **FINDING-P6-SOLVER-2 (residual, self-found & fixed @40be57a):** the round-1 high-segment fix capped its post-bracket linear scan at `highScanCapCents=4e6` ⇒ (a) reachable near-ceiling targets whose cheapest price sits beyond the cap → false UNREACHABLE (CONFIRMED: com10/ali0/frete0/custo5000, target 89.40 → R$826.445,64), (b) 20–55s solves. Round-2 reviewers split on severity (Opus N1 nit / REFUTE blocker). Fix = exact bracket to `firstCentExactAtLeast(target−0.005)` + fixed ±20 000-cent window (provably covers the ≤150/(ceiling−target)≤15 000-cent perturbation); full R$1M reach, O(log)+bounded, <1s. Amendment candidate: a "bounded scan cap" as a safety net silently converts reachable→UNREACHABLE — bound the search by the PROVABLE perturbation window, not an arbitrary cent cap.
- None blocking. Slice A worker confirmed calc_repository_test.go is `//go:build integration` (needs live PG) — SQL for tariff_defaults is exercised only by build/vet + migrate-count; handler contract covered by stub-backed transport tests. Consistent with existing repo pattern, not a defect.
