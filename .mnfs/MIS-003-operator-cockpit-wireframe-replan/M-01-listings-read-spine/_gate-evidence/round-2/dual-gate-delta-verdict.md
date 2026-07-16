# Dual-gate DELTA (slices 8+9) — MERGED VERDICT: **FAIL**

```yaml
gate: P6 dual gate, DELTA scope (REVIEW-STANDARD §9)
base: e2cde36 (prior dual gate PASS on code criteria, commit b1179767)
head: c4e8ab913c132d4929c3bc60156a908379a78043
diff: _gate-evidence/round-2/delta-8-9.diff (759 lines, 12 files)
claude_side: COLD Opus subagent (model=opus, clean context, implementer≠reviewer) — VERDICT: FAIL
gpt_side:    GPT-5.6 Sol --effort medium (/codex:rescue, fresh thread) — VERDICT: FAIL
dispatch: SIMULTANEOUS at fixed SHA, both on prompt-pack §14, neither saw the other's output
merge: §8 — both sides FAIL, AGREEMENT on the blocking core. Gate FAILS. M-01 does NOT close.
```

## Reconciliation table (§8 mandatory)

| # | Finding | Opus side | Sol side | Merged | Severity |
|---|---|---|---|---|---|
| **B1** | **Degraded read path drops contract-declared fact-dependent filters** — on Oracle outage `passThrough` re-queries with the original `q` and never applies `matchesDependentFilter`. `filter.has_exception` / `filter.exception=below_margin` therefore return provably NON-matching rows as a confident 200. | **blocking** (`read_service.go:365-375,:321-322`) | **blocking** (`read_service.go:321`) | **AGREE — blocking** | blocking |
| **B2** | **All cost-reader errors collapse into "facts unavailable"** — cancellation, timeout, adapter/data defects become null facts indistinguishable from genuinely-absent data. Only source-unavailable should degrade; the rest must propagate. | important (framed as telemetry gap, `:394-396,:122`) | **blocking** (`read_service.go:394`) | **AGREE on substance; take the HIGHER severity** | blocking |
| **B2b** | Outage is swallowed with **zero telemetry** — `ceilingErr` never logged, cost error discarded; `read_service.go` has no `slog` call at all. Total Oracle outage = HTTP 200, null margins cockpit-wide, not one log line. Same delta logs a WARN for one unmapped status string → internally inconsistent about operator honesty. | **important** (`:394-396,:122`) | folded into its blocking(:394) | **AGREE** | important |
| **I1** | `unavailablePolicyService` wrapper is now **verified-unreachable** — slice 8 removed its only producer of `ReadErrorSourceUnavailable`; the `IsReadErrorCode` branch is dead. Name asserts "unavailable" while wrapping the live policy service. | **important** (`root.go:106-116`) | not raised | **Opus-only — accepted** (verified claim, not speculative) | important |
| **I2** | **Tests pin the defect as intended behavior** — `TestReadServiceDependentFilterPassesThroughWhenOptionalFactsUnavailable` asserts `len(page.Items)==1` for a row that does NOT match `exception=below_margin`. `has_exception` has NO degraded-path coverage. | **important** (`read_service_test.go:499-564`) | noted ("the new tests explicitly pin these false positives") | **AGREE** | important |
| **S1** | `NewBatchReader(oracleDB,…)` unguarded vs. the `internalReadAvailable`/`oracleDB==nil` idiom elsewhere in the file. Verified **safe** (`ensureBatchAvailable` nil-checks). Also skips the `TimingReader`/cache profitability applies → live per-page Oracle latency unproven (C09 p95 3.17ms was measured against a fake). | suggestion (`root.go:485`) | not raised | Opus-only — accepted as suggestion | suggestion |
| **N1** | `enrichGroups` returns `false` for `factsUnavailable` on the error path, discarding known state. Harmless (callers return on `err`). | nit (`:246-249`) | not raised | Opus-only — accepted as nit | nit |

**No contradictions between the two sides.** Sol raised a strict subset of Opus's blocking/important set, at equal-or-higher severity. Opus added I1/S1/N1 with verified receipts. Merge = union; severity = max.

## Both sides independently CONFIRMED correct (no finding)

- **C10's PASS is HONEST, not threshold-gaming.** `under_review` is a real documented ML status; the 7 rows carried it as their raw provider value and previously fell through `default→unknown`. The explicit case is a genuine mapping win — collapsing them into `paused` would have been the gaming move. Unchanged summary buckets (active 10 / paused 17) corroborate no bucket was distorted.
- **Slice 9 is CLEAN.** `slog.Warn` fires only on `default→unknown`, never for a mapped status. Migration 0037 additive + reversible (DROP + re-ADD same-named constraint, widened IN-list only, no data rewrite); 0037 is the ONLY migration touched, 0036 byte-unchanged; existing `unknown` rows stay valid. Domain enum / OpenAPI (3 sites) / SDK union = identical 8 values.
- `Get`'s cost-outage path correctly keeps the real ICMS matrix with null `below_margin_at_uf` rather than nulling the whole matrix.
- Installation and policy failures correctly remain hard (pinned by `errors.Is` + not-called assertions).
- No `tenant_id` scoping regression — the delta touches no query construction.

## The verdict in one line

Slice 8 traded an honest 503 for a **confidently wrong 200**. Pre-slice-8, an Oracle outage on
`filter.has_exception=false` served no data; post-slice-8 it serves rows whose `sync_state='error'` /
unresolved link make them **provable non-matches** — facts read from Postgres, fully available,
independent of Oracle — with every field non-null and no log line. That is not "unknown rendered as
null" (ADR-17 compliant); it is a **known-false fact returned as a match**, which is precisely what
ADR-17 exists to prevent. The degrade is the right *direction*; its filter semantics are wrong.

## Consequence

**M-01 does NOT close.** P8 CLOSED is withheld. A corrective **slice 10** is required, scoped to the
slice-8 degrade path only (slice 9 needs no rework). Degraded-filter semantics require a hub ruling —
see `HUB-EVENT-BLOCKED-dual-gate-delta-fail.md`.
