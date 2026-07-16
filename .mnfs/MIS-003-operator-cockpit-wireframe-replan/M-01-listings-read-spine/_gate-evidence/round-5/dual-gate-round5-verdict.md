# Dual gate round 5 — merged verdict (REVIEW-STANDARD §8)

```yaml
sha_under_review: 982d44e       # fixed; both reviewers read the same tree
base: a6878dc6                  # §9 delta-only re-review = slice 12
delta: 7 hunks, 1 file (apps/server_core/internal/modules/listings/application/read_service.go)
claude_side: COLD Opus subagent (Agent tool, model=opus, clean context) — implementer != reviewer
gpt_side:    gpt-5.6-sol --effort medium --fresh
independence: neither reviewer saw the other's output; dispatched at a fixed SHA
standard: docs/REVIEW-STANDARD.md §14 (MAIN checkout; worktree docs/ stale — hub-ruled, not fixed mid-flight)
merged_verdict: PASS
```

## The checklist, and it is one item long

Round 4 merged to FAIL on a single `important` — **H1**, raised by Sol, missed by Opus. Slice 12
exists to close it and nothing else.

| # | Round-4 severity | Opus | Sol | **Merged** | Receipt (milestone owner re-verified in-tree) |
|---|---|---|---|---|---|
| **H1** cost-fact propagate ERROR unattributable to a read path | important | **RESOLVED** | **RESOLVED** | **RESOLVED** | `read_service.go:430` now reads `slog.Error("listings: cost fact read failed", "err", costErr, "op", op)`. `grep -n "slog\." read_service.go \| grep -v '"op"'` → **empty**. Both reviewers independently confirmed the closure is total, not partial. Sol — who raised it — verified its own finding rather than assuming it closed. |

**Both sides PASS. No contradiction to reconcile this round.** Sol returned zero findings. Opus
returned one `nit` (`:412`: `enrich` now carries 7 parameters, two of which — `unavailable error`
and `op string` — are state and telemetry rather than data, so the slice extends the already-noted
`(error, error)` shape instead of collapsing it). Opus itself judged no fix owed by this slice and
folded it into the standing adjudicated nit. A nit never gates (§14).

## What both sides independently verified

- **Behavior-neutral, checked from the diff rather than the claim.** 7 hunks: 5 call-site argument
  additions, 2 signature additions, 1 doc-comment sentence, 1 `slog` field. The propagate-vs-degrade
  branch (`:428-436`), the `unavailable == nil` gate (`:415`), the `Get` matrix guard (`:338`), and
  `isSourceUnavailable` (`:22-27`) appear only as unchanged context. `:430` sits inside the propagate
  branch, but the change is the field alone — the `return unavailable, costErr` beneath it is
  untouched. No input produces a different service return.
- **Once per request, never per row.** `op` is consumed only at `:430`, inside `if len(ids) > 0`
  under the unchanged collection gate. The per-item loop (`:438-464`) never reads it.
- **`op` is TRUE at every call site**, traced through real callers rather than names — the check that
  matters, since a wrong `op` would be worse than none: it would lie to the operator instead of
  leaving him to infer. `scan:357` has exactly one caller (`List:167`) and passes `"List"`.
  `scanGroups:213` has exactly one caller (`ByProduct:192`) and passes `"ByProduct"`. Direct sites:
  `List:157`, `ByProduct:198`, `Get:330`. `enrichGroups:273` forwards its own `op` rather than
  hardcoding one. Opus grepped the whole listings module for callers and found none outside the file.
- **AI-slop checklist (HARNESS.md:203-209), whole file, every item — clean.** No speculative
  abstraction (`op` is a parameter, and both consumers exist today). No narration: the only whole-file
  grep hits are the known false positives — `gate` inside `propagate`, and `sort.Slice` matching a
  `slice` pattern. The added doc sentence says what `op` is, not what a review asked for.
- **The absent test — both sides judged the call correct.** The read's return value is byte-identical
  for every input, so there is no behavior a unit test at this layer could pin. A test would have to
  install a `slog` handler and assert a log field, which asserts the telemetry rather than the read.
  Opus put it most sharply: *a test would have been the defect.* The one honest check of `op`
  truthfulness is the caller-graph trace, which is static — and both reviewers performed it.
- **Layering, tenancy, fail-honest.** Application layer only; `ports`/`domain` imports unchanged; no
  query surface touched, so tenant scoping is unaffected; the integrity-critical read still fails
  honest rather than nulling.

## A correction to my own evidence

`slice12-L0-report.md` states the file has **16** `slog` sites. It has **15**. The miscount is mine
and Opus caught it while checking a claim it was invited to distrust. The substance is unaffected —
what H1 turns on is that the set of `slog` sites *without* `op` is now empty, and it is — but the
number was wrong and is corrected here rather than left standing in the ledger. Sol was handed the
correction in its prompt-pack and independently re-derived 15.

The wider point is worth recording. Round 4 failed on a finding Opus never examined, and this round
Opus caught an error in the milestone owner's own report. Neither reviewer is reliably the deeper
one; the dual gate earns its cost precisely because which side goes deeper is not predictable in
advance.

## Consequence

**MERGED VERDICT: PASS.** Zero blocking, zero important, one non-gating nit. Every round-3 and
round-4 finding is now resolved or adjudicated:

- G1 — resolved in-module (cancel/timeout); residual deferred-with-name to **hub board task #2**
  (`internal_read` taxonomy refactor). Per hub ruling (C).
- G2, G3, G4 — resolved (slice 11), confirmed by both sides in round 4.
- G5 — resolved: the wording defect in slice 11, the last attribution gap (H1) in slice 12.
- G6 — nit, non-gating, standing.
- H2, H3 — suggestions, declined with reasons recorded in the round-4 verdict.
- H4, H5 — questions routed to the hub board; outside this delta. Reviews verify, they do not
  generate scope.

**Code criteria for M-01 are met.** This is not a close: **only browser/live QA passes a milestone**
(HARNESS). P7 fresh browser QA at `982d44e` is the remaining gate before P8 CLOSED. Not pushed.
