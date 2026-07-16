# Slice 12 — L0/L1 deterministic pre-pass report (§7)

```yaml
base_sha: a6878dc6 (slice 11)
candidate: working-tree diff — 1 file, 7 hunks, ~15 changed lines
scope: apps/server_core/internal/modules/listings/application/read_service.go ONLY
gocache: absolute (worktree .gocache)
buildvcs: false (sandbox Git safe.directory VCS-stamp; no repo config changed)
run_by: milestone-owner (chip), INDEPENDENT of the worker — not a relayed claim
charter: close H1, the single `important` from the merged round-4 dual gate
```

## L0 (precedes review dispatch)

- `go build -buildvcs=false ./...` → **exit 0**
- `go vet ./...` → **exit 0** (whole repo)

## L1 unit (ran ∥ review per §15)

`go test -count=1 ./internal/modules/listings/... ./internal/composition/...` → **exit 0**
(connectors, integrations, application, domain, ports, transport, composition all `ok`)

No new test. A `slog` field is not behavior a unit test should pin at this layer, and writing one
would be test theater — asserting the log rather than the read. The slice is behavior-neutral by
construction; the existing suite passing unchanged is the signal that it stayed that way.

No integration lane re-run. The lane exercises read behavior over a real database; this slice
changes no read behavior, only which fields an error log carries. Re-running it would prove nothing
slice 11's green lane did not already prove, and the lane is expensive. Recorded here as a
deliberate call, not an omission.

## Milestone-owner independent verification (not relayed from the worker)

- **Every log site carries `op`.** `grep -n "slog\.\(Warn\|Error\)" read_service.go` → **15** sites;
  the same grep piped through `grep -v '"op"'` → **0**. (CORRECTION: this report originally said 16.
  The true count is 15 — the miscount was mine, caught by the round-5 cold Opus reviewer while
  checking a claim it was invited to distrust. The substance is unaffected: what H1 turns on is that
  the set of sites *without* `op` is empty, and it is. Left visible rather than silently edited.)
  The site that H1 named, now `:430`, reads
  `slog.Error("listings: cost fact read failed", "err", costErr, "op", op)`. `:100` — the same
  message on the same condition inside `Summary` — still reads `"op", "Summary"`. The invariant is
  now uniform across the file.
- **Behavior-neutral.** The full diff is 7 hunks: six call sites gain a trailing literal, two
  signatures gain `op string`, one doc comment gains one sentence, one log gains one field. The
  propagate-vs-degrade branch (`:428-436`), the `unavailable == nil` gate that keeps the cost fetch
  once-per-request (`:415`), the `Get` matrix guard (`:338`), and `isSourceUnavailable` (`:22-27`)
  are untouched.
- **`op` is true at every call site**, verified against actual callers rather than the worker's
  claim: `scan` is called only from `List` (`:167`) and passes `"List"`; `scanGroups` is called only
  from `ByProduct` (`:192`) and passes `"ByProduct"`. A wrong `op` would be worse than none — it
  would lie to the operator. These are right.
- **No narration.** Whole-file grep for hub rulings / finding IDs / round numbers / gate / reviewer
  / "used to" / "previously" / TODO → 3 hits, **all false positives**: the pattern `gate` matching
  inside `propagate`/`propagates`. The added comment sentence states what `op` is, not why a
  reviewer asked for it.
- **Scope.** `git status` shows the only Go file modified is `read_service.go`. No test file needed
  changing (no test calls `enrich`/`enrichGroups` directly). No `internal_read/` (the hub-owned
  cross-module seam), no composition, no transport, no migration, no OpenAPI, no SDK.
  `docker/dev/*.sh` show as modified in the worktree — these are **pre-existing hub dev-stack
  changes, not this chip's**, and are deliberately excluded from the commit.

## Independent slice review (§14)

Cold sonnet subagent, clean context, implementer ≠ reviewer. **APPROVE** — zero blocking, zero
important, one suggestion (`op` as a bare string has no compile-time typo safety; the reviewer
itself judged a typed enum to be speculative abstraction for a single-file need with four call
sites, and declined to press it). The reviewer independently re-derived the two facts that matter:
that the diff is behavior-neutral, and that `op` is correct at every site — tracing `scan`→`List`
and `scanGroups`→`ByProduct` through their real callers rather than accepting the labels.
