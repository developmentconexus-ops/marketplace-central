# Dual gate round 4 — merged verdict (REVIEW-STANDARD §8)

```yaml
sha_under_review: a6878dc6      # fixed; both reviewers read the same tree
base: 7f5a1b8c                  # §9 delta-only = slice 11
diff: _gate-evidence/round-4/delta-round4-code.diff   # 410 lines, 2 files
claude_side: COLD Opus subagent (Agent tool, model=opus, clean context) — implementer != reviewer
gpt_side:    gpt-5.6-sol --effort medium --fresh
independence: neither reviewer saw the other's output
standard: docs/REVIEW-STANDARD.md §14 (MAIN checkout; worktree docs/ stale — hub-ruled, not fixed mid-flight)
merged_verdict: FAIL
```

## Part 1 — round-3 resolution check (merged)

| # | Round-3 severity | Opus | Sol | **Merged** | Receipt (milestone owner re-verified in-tree) |
|---|---|---|---|---|---|
| **G1** adapter taxonomy defeats strict classification | blocking | RESOLVED (in-module half) | PARTIAL — accepted, non-gating | **RESOLVED per ruling (C)** | `read_service.go:23-25` checks `errors.Is(ctx.Canceled/DeadlineExceeded)` **before** `IsReadErrorCode`. Both sides independently traced the **production** wrap and agree the cause survives it: `oracle/reader.go:516` `safeOracleCause` returns context sentinels unflattened → `read_error.go:28-33` `Unwrap()` returns `Cause` → `errors.Is` walks it. Opus went further and proved the chain is lossless end-to-end (`batch_reader.go:57`, `cache/cache.go:355-357,379-381`, `internalread/cost_reader.go:16-17,29-30` all return `err` verbatim). Residual (adapter/data defect still flattened at `reader.go:519`) = **deferred-with-name to hub task #2**, adjudicated. Neither reviewer re-blocked on it. |
| **G2** Get nulls whole ICMS matrix on cost-only outage | important | RESOLVED | RESOLVED | **RESOLVED** | `read_service.go:338` guard back to `ceilingErr != nil`. Cost-outage path traced: `enrich:441` nils `Cost`, `Get:341` still calls `icmsWorstCaseByUF`, which fills `WorstCaseICMSPct`/`PriceNetBasis` from price+ceiling alone (`:571-573`) and `belowMargin:528` nils only `BelowMarginAtUF`. Independently confirmed on live data — `redrive-post-slice11.md` (25-entry matrix served, only `below_margin_at_uf` null). |
| **G3** the G2 branch is unpinned | important | RESOLVED | RESOLVED | **RESOLVED** | `read_service_test.go:320-332` asserts the matrix non-nil, each row's ICMS/basis non-nil, `BelowMarginAtUF` nil. Both sides independently confirmed it **pins**: revert `:338` → `:321` `t.Fatalf`. Opus added the anti-vacuity check — `resolved():108` sets `Price: money("90")`, so `PriceNetBasis` is computable and its assertion has teeth. |
| **G4** comment narration / PR-voice | blocking | RESOLVED | RESOLVED | **RESOLVED** | Both ran **whole-file** greps (not diff hunks) for hub rulings / finding IDs / round numbers / superseded-test narration / unanchored TODOs → **0 hits** on both files. Surviving comments state invariants only. |
| **G5** WARN reports an outcome that didn't happen | important | RESOLVED | **PARTIAL** | **PARTIAL — see below** | The wording defect **is** fixed: every degrade WARN now fires after the outcome is known (`:162` inside the `!needsBelowMarginScan` block after `enrich`; `:203` after the scan return at `:191-193`; `:335`), and the propagate arms ERROR + return immediately (`:215`, `:232`, `:359`, `:376`). Once per request, never per row (cost fetch gated by `enrich:414`; scans early-return). **But `:429` remains.** |
| **G6** `(error, error)` as state carrier | nit | UNRESOLVED | UNRESOLVED | **UNRESOLVED — nit, non-gating** | `read_service.go:411`, `:268` unchanged. Correctly threaded at all call sites. Never gates (§14). |

**Test-theater check — both sides cleared it independently.** `wrappedCause:18-21` builds
`NewReadError(ReadErrorSourceUnavailable, …, cause)` — the exact shape
`wrapOracleError`+`safeOracleCause` emits for a context cause, not a fake carrying whatever the
assertion wants. The slice's self-reported caveat is **genuinely structural, not a gap explained
away**: 24 subtests = 2 causes × 2 facts × 6 paths; the 8 `dependent_filter_(group_)scan` subtests
cannot pin because `scan`/`scanGroups` (`:214-217`, `:231-234`, `:358-361`, `:375-378`) return
`unavailable` as an error regardless of cause — outcome identical with the fix reverted. 24−8=16
carry the pinning force, matching the slice's claim. (The slice review's "2 subtests" is loose
wording for 2 subtest *shapes*; the arithmetic is right.)

**Artifact fidelity:** Opus verified `git diff 7f5a1b8c a6878dc6 -- apps/` is **byte-identical** to
the reviewed diff, exactly 2 files outside `.mnfs/`, `internal_read/` untouched. The hub-owned seam
boundary held.

## Part 2 — merged delta findings

| # | Finding | Opus | Sol | **Merged** |
|---|---|---|---|---|
| **H1** | **`read_service.go:429`: the cost-fact propagate ERROR is the only log site in the file with no `op`.** All 15 other `slog` sites carry it — including `:100`, the *same message on the same condition* in `Summary`. `enrich` is shared by List/ByProduct/Get and doesn't know its caller, so a propagated cost-reader failure reaches the operator unattributable to a read path. The slice-11 brief listed `err` + `op` as a **surviving invariant** of the G5 fix (`slice11-brief.md:75`). | not raised | **important** | **important** |
| **H2** | `read_service.go:470`: `unavailableFact` keys off `ceilingErr` while the adjacent log emits `unavailable` — two variables. Correct only via the implicit invariant that `enrich:414` skips the cost fetch when entry-`unavailable` is non-nil. Verified correct at all three call sites (`:162`, `:203`, `:335`). | suggestion | not raised | **suggestion** |
| **H3** | `read_service.go:215,232,359,376`: identical `slog.Error` literal 4×, differing only in `op` (§4 rule of three). Distinct early-returns; a helper may be worse than the duplication. | suggestion | not raised | **suggestion** |
| **H4** | `read_service.go:414` (**pre-base, cannot gate**): on a *ceiling* outage the cost fetch is skipped and `:441` nils `Cost`, so `cost: null` is served while the cost source is healthy — arguably ADR-17 direction (b) one level up. Pinned as intended (`read_service_test.go:264,287`), certified round 2/3, unchanged by this delta, outside the brief. | question — hub queue | not raised | **question — hub board, not a defect of this delta** |
| **H5** | `internal_read/adapters/cache/cache.go:179-197` (**outside delta, hub-owned**): `group.DoChan` shares one loader's result, so caller A's `context.Canceled` can be delivered to uncancelled caller B. Post-slice-11, B propagates instead of degrading. Honest either way (no false fact served); the delta did not create it but does change B's outcome at a distance. | question — hub board | not raised | **question — hub board** |

## The one contradiction, and how it merged

Opus returned **PASS**; Sol returned **FAIL** on H1. Opus did not raise `:429` at all.

This is not a disagreement about facts — it is, once again, a difference in depth, and this time the
depths are *inverted from round 3*. Opus read the wrap chain further than Sol (proving the context
cause survives four hops of collaborators Sol did not open) and cleared G1 more thoroughly. Sol read
the *telemetry surface* further, enumerated every log site, and found the one that breaks the
pattern. Neither reading is wrong; each is incomplete where the other is deep. **The milestone owner
re-verified H1 in-tree before merging** — `grep "slog\.\(Warn\|Error\)"` over the whole file returns
16 sites, and `:429` is the sole one without `op`, while `:100` proves the author knew the invariant
and applied it in the sibling function. Sol's receipt is correct and specific.

§8 merge = **union of findings, max severity, deeper receipt wins the contradiction.** H1 stands as
**important**, and §14 says an unresolved important is a FAIL. Opus's PASS does not survive contact
with a finding it did not examine; a PASS is only as strong as its coverage.

**Honest note on scope.** `:430` at base 7f5a1b8c already lacked `op`, so H1 is not a slice-11
regression, and a strict §9 delta-only reading could push it out of round 4. It stays in, for two
reasons: G5 was explicitly assigned to slice 11 as *"telemetry must not state the opposite of the
outcome… invariants that survive: ERROR on propagate, WARN on degrade, `err` + `op`"*, and slice 11
rewrote every other log site in this file to satisfy exactly that invariant. Leaving one site behind
is an incomplete G5, not a separate pre-existing issue. G5 merges **PARTIAL**, not RESOLVED.

## Consequence

**MERGED VERDICT: FAIL. M-01 does NOT close. P8 CLOSED withheld.**

One `important` (H1), zero `blocking`. Slice 12 is a micro-slice: thread the operation into `enrich`
(both callers of `enrichGroups` too) so `:429` attributes to its read path, keeping exactly one log
per request. H2/H3 are suggestions the author may take or decline. H4/H5 go to the hub board as
questions, not to slice 12 — reviews verify, they do not generate scope.
