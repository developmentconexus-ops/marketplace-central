# Dual-gate round 3 (slice 10 delta) — MERGED VERDICT: **FAIL**

```yaml
gate: P6 dual gate, DELTA scope (REVIEW-STANDARD §9) + explicit round-2 resolution check
base: c4e8ab91
head: 7f5a1b8c
diff: _gate-evidence/round-3/delta-round3.diff (715 lines, 3 files)
claude_side: COLD Opus subagent (model=opus, clean context, implementer≠reviewer) — VERDICT: FAIL
gpt_side:    GPT-5.6 Sol --effort medium (/codex:rescue, fresh thread) — VERDICT: FAIL
dispatch: SIMULTANEOUS at fixed SHA 7f5a1b8c, both on the §14 prompt-pack, neither saw the other's output
merge: §8 — both FAIL. One CONTRADICTION (B2), resolved in favour of the side with the deeper receipt.
verdict: FAIL. M-01 does NOT close. Slice 11 required + one hub ruling.
```

## Part 2 — round-2 resolution check

| # | Opus | Sol | Merged | Receipt |
|---|---|---|---|---|
| **B1** degraded path drops fact-dependent filters | RESOLVED | RESOLVED | **RESOLVED** | `read_service.go:361-363` (ceiling pre-check), `:377-379` (cost mid-scan), `:215-217`,`:231-233` (scanGroups); `passThrough`/`passThroughGroups` deleted; `needsBelowMarginScan:485-487` uniform. |
| **B2** all cost-reader errors collapse to "facts unavailable" | RESOLVED | **UNRESOLVED** | **UNRESOLVED — Sol wins** | See G1 below. Opus verified only the application layer, where the gate IS strict (`:65-67`, `:76-79`, `:99-102`, `:112-115`, `:144-147`, `:324-327`, `:429-432`). Sol read one layer down into the adapter and found the taxonomy itself is not discriminating. Sol's receipt is deeper and independently re-verified by the milestone owner. |
| **B2b** outage swallowed with zero telemetry | RESOLVED | RESOLVED | **RESOLVED** | `slog.Error` `:77,:100,:113,:145,:325,:430`; `slog.Warn` `:86,:103,:117,:189,:329,:433`; once per request (cost fetch gated by `if unavailable == nil` `:415`; scan/scanGroups early-return `:377`,`:231`) — no per-row flood. (But see G5: the WARN's wording can misreport the outcome.) |
| **I1** unreachable policy wrapper | RESOLVED | RESOLVED | **RESOLVED** | `root.go:472-473` wires `NewPolicyReader(marketSvc)` directly; type + `listingsports` import gone; repo-wide grep = 0 hits. |
| **I2** tests pin the defect | RESOLVED | RESOLVED | **RESOLVED** | `read_service_test.go:550` replaces the defective test; `t.Fatalf` on `err == nil` across 3 filter shapes × 2 outage sources × List/ByProduct; over-correction guard `:619`. (But G3: a *different* branch is now unpinned.) |
| **S1** `NewBatchReader` unguarded / skips `TimingReader` | UNRESOLVED — defensible | UNRESOLVED — defensible | **UNRESOLVED — defensible, non-gating** | `root.go:472` unchanged. Accepted as a *suggestion* in round 2, verified safe (`ensureBatchAvailable` nil-checks), explicitly outside slice 10's R2 scope. Both sides agree it does not gate. The unproven live per-page Oracle latency stays separately tracked. |
| **N1** `enrichGroups` discards known state | RESOLVED | RESOLVED | **RESOLVED** | `read_service.go:272-274` returns `unavailable, err`. |

**6 of 7 resolved. B2 is not.**

## Part 1 — merged delta findings

| # | Finding | Opus | Sol | Merged |
|---|---|---|---|---|
| **G1** | **Adapter taxonomy defeats the strict classification the ruling mandated.** `wrapOracleError` (`internal_read/adapters/oracle/reader.go:508-513`) labels EVERYTHING that is not `sql.ErrNoRows` as `ReadErrorSourceUnavailable` — `context.Canceled`, `DeadlineExceeded`, and query/scan/iteration defects included. So `isSourceUnavailable` returns true in production for errors that must propagate. Unit tests pass only because they inject raw fakes carrying the correct code. Net effect: on a no-filter read, a genuine data/driver defect still degrades to null facts — the exact substance B2 was raised to stop. | not raised | **blocking** | **blocking** |
| **G2** | **`Get` regressed a round-2-certified behavior.** `read_service.go:337` changed its matrix guard from `ceilingErr != nil` to `unavailable != nil`, so a **cost-only** outage now nulls the ENTIRE `ICMSWorstCaseByUF`. But `icmsWorstCaseByUF:563-567` populates `WorstCaseICMSPct` and `PriceNetBasis` (`priceNetBasis:574-576` — price + ceiling only) with **no cost dependency**, and `belowMargin:520-522` already returns nil for a nil cost. Pre-slice-10 returned the full matrix with only `below_margin_at_uf` null — which round 2 **explicitly certified as correct** ("keeps the real ICMS matrix… rather than nulling the whole matrix"). This discards KNOWN-TRUE facts because an unrelated fact is unavailable. ADR-17 licenses null for *unknown* facts, not for known ones. | **important** | not raised | **important** |
| **G3** | **The G2 regression is invisible to the suite.** `read_service_test.go:317` guards the matrix assertion behind `outage.ceilingErr &&`, so the cost-outage case asserts nothing about `ICMSWorstCaseByUF`. This is round-2's own I2 criterion turned on slice 10: a test that does not pin is not coverage. | **important** | not raised | **important** |
| **G4** | **Comment narration / PR-voice comments in production code.** `read_service.go:172-177`, `:271-278`, `:318-324` and the test doc comments narrate adjudication history — "Per the hub ruling (B1, PURE uniform fail-honest)", "pins B2". `HARNESS.md:203-204` is unambiguous: *"AI-slop checklist — any hit = REJECT the slice: … comment narration / PR-voice comments"*. A binding REJECT rule, hit directly. | not raised | **blocking** | **blocking** |
| **G5** | **The degrade WARN can state the opposite of what happened.** The ceiling WARN (`:117`, `:189`, `:329`) reads "degrading fact to null" but fires BEFORE the `needsBelowMarginScan` branch (`:153`, `:192`), so a fact-dependent-filter request logs "degrading to null" and then returns 503. The log contradicts the outcome — undercutting exactly the operator honesty B2b was raised to restore. | suggestion | **important** | **important** (higher severity wins per §8) |
| **G6** | `(error, error)` uses `error` as a state carrier in the non-terminal position; transposable at call sites, against Go convention. Correctly threaded at all six call sites today. | nit | not raised | nit |

## The one contradiction, and how it merged

Opus marked **B2 RESOLVED**; Sol marked it **UNRESOLVED**. Not a disagreement about facts — a
difference in depth. Opus reviewed the delta's three files, where the classification IS strict and
correct. Sol left the delta to read the adapter that produces the codes, and found the labels
themselves are indiscriminate. Both descriptions are true of what each read; Sol's is true of the
running system. **The milestone owner re-verified `reader.go:508-513` in-tree before merging** — the
`errors.Is(err, sql.ErrNoRows)` check is the only discrimination that exists. §8 merge = union of
findings, max severity, deeper receipt wins the contradiction.

This is also the case for reviewing at fixed SHA with *two* readers rather than one: a single
reviewer bounded strictly to the delta would have passed a defect that lives outside it.

## Milestone-owner note on my own slice review

`slice10-review.md:29-31` recorded G2 as a **"BONUS fix nobody asked for… closing a latent gap."**
That was wrong: it is a regression of behavior a prior gate had certified. The §14 slice review is
mine to own — I attached it, so I own the misclassification. Recorded here rather than quietly
corrected. It is precisely why the cold gate exists downstream of the slice review.

## Scope split — G1 is NOT the chip's to fix

`wrapOracleError` is a **cross-module shared seam**. Its `ReadErrorSourceUnavailable` output is
consumed by `catalog/transport/http_handler.go:289`, `product_links/application/generation_service.go:319`
and `profitability/adapters/internalread/fact_reader.go`, besides listings. The taxonomy
(`internal_read/domain/read_error.go:8-9`) has only TWO general codes — `source_unavailable` and
`unsupported_query` — so there is **no existing code** for cancellation or for an adapter/data
defect; fixing G1 means ADDING to the taxonomy and re-classifying for every consumer.

Per HARNESS, one writer owns a shared seam and the hub owns cross-module seams. The M-01 chip must
NOT unilaterally edit it. **G1 escalates to the hub** — see
`HUB-EVENT-BLOCKED-gate-round3-adapter-taxonomy.md`.

G2/G3/G4/G5 are inside the listings module = chip's, slice 11.

## Governance drift observed (not a code finding)

This worktree's `docs/` is behind the main repo's governance commits: it has no
`docs/REVIEW-STANDARD.md`, and its `AGENTS.md` still points at `docs/superpowers/HARNESS.md`. The
binding copies (`docs/HARNESS.md`, `docs/REVIEW-STANDARD.md`) exist only in the main checkout. The
prompt-pack handed both reviewers a worktree path for `REVIEW-STANDARD.md` that does not resolve
there; both nonetheless returned §14-conformant output, so they sourced the standard elsewhere
(codex's cwd is the main root). Flagged for the hub: reviewers should be pointed at the main-repo
governance paths until the worktree is caught up.

## Consequence

**M-01 does NOT close. P8 CLOSED withheld.** Slice 11 required (G2/G3/G4/G5), plus a hub ruling on
G1's scope before B2 can be honestly called resolved.
