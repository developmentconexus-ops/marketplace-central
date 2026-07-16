# HUB EVENT — BLOCKED (gate round 3 FAIL; G1 needs a ruling on a shared seam)

```yaml
event: BLOCKED
from: chip M-01-listings-read-spine
to: HUB (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
gate: P6 dual gate round 3, delta c4e8ab91..7f5a1b8c
result: FAIL — both sides, independently, simultaneous at fixed SHA
verdict: _gate-evidence/round-3/dual-gate-round3-verdict.md (full §8 reconciliation table)
blocking_on: ONE ruling (G1). G2/G3/G4/G5 are chip-owned and already moving as slice 11.
```

## Round-2 resolution check: 6 of 7 resolved

B1, B2b, I1, I2, N1 **RESOLVED** with receipts, both sides agreeing. S1 **UNRESOLVED but
defensible** (suggestion, out of R2 scope, verified safe) — both sides agree it does not gate.

**B2 is UNRESOLVED**, and that is G1.

## G1 — the ruling I need

`wrapOracleError` (`internal_read/adapters/oracle/reader.go:508-513`) labels **everything** that is
not `sql.ErrNoRows` as `ReadErrorSourceUnavailable` — `context.Canceled`, `DeadlineExceeded`, and
query/scan/iteration defects included. Slice 10's classification is strict at the application layer,
but the labels it classifies on are indiscriminate at the source. The unit tests pass only because
they inject raw fakes carrying the correct code; **in production a data/driver defect on a
no-filter read still degrades to null facts** — the exact substance B2 was raised to stop.

Sol found it; Opus (bounded to the delta's 3 files) marked B2 resolved. Not a contradiction of fact
— a difference in depth. **I re-verified `reader.go:508-513` in-tree myself** before merging: the
`errors.Is(err, sql.ErrNoRows)` check is the only discrimination that exists. Sol's receipt wins.

### Why this is NOT mine to fix

`wrapOracleError` is a **cross-module shared seam**. Its `source_unavailable` output is consumed by
`catalog/transport/http_handler.go:289`, `product_links/application/generation_service.go:319`, and
`profitability/adapters/internalread/fact_reader.go` — not just listings. And the taxonomy
(`internal_read/domain/read_error.go:8-9`) has only TWO general codes, `source_unavailable` and
`unsupported_query`: there is **no code for cancellation and none for an adapter/data defect**.
Fixing it means ADDING to the taxonomy and re-classifying for every consumer. One writer owns a
shared seam; the hub owns cross-module ones. I am not touching it unilaterally.

### What I CAN do without the seam — verified, not assumed

`ReadError.Unwrap` returns `Cause` (`read_error.go:28-33`) and `safeOracleCause:515-518`
**preserves** `context.Canceled` / `DeadlineExceeded` as that cause. So `errors.Is` reaches through
the wrap **today**, and listings can de-classify cancellation locally:

```go
func isSourceUnavailable(err error) bool {
    if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
        return false
    }
    return internalreaddomain.IsReadErrorCode(err, internalreaddomain.ReadErrorSourceUnavailable)
}
```

That closes the **cancellation/timeout half** of B2 inside my module. The **defect half does not
close this way**: `safeOracleCause:519-523` flattens a query/scan defect into
`errors.New("oracle driver error")` / `oracle error code=N`, erasing the distinction from a genuine
connectivity outage before listings ever sees it. No application-layer code can recover it.

## Options

- **(A) Fix the taxonomy at the adapter.** Add codes (cancellation, adapter/data defect), classify
  in `wrapOracleError`, update the three other consumers. The correct fix. But it is a cross-module
  seam change that widens M-01 from a listings read spine into an `internal_read` refactor, and it
  puts a chip inside a seam three other modules read.
- **(B) Chip-local only.** Ship the `errors.Is` de-classification above in slice 11. Cancellation
  and timeout stop degrading. Adapter/data defects still degrade to null — **but now with a WARN**
  (B2b landed), so it is misclassified, not silent. Defer the residue with a name.
- **(C) (B) now + (A) as its own hub-owned work** on the `internal_read` seam, scheduled separately.

**My recommendation: (C).** With the honest caveat: **(B) does not fully close B2** — it closes the
cancellation half and leaves the defect half misclassified-but-logged. If you consider that residue
unacceptable for M-01's bar, then M-01 cannot close without (A), and (A) is yours to place, not
mine to take.

I am not going to overstate this one. Last round my recommendation (A+C) was partly wrong and you
corrected it — you rejected my (A) because it produced a silently under-inclusive result. So: the
above is my read, the residue is stated plainly, and the scope call is yours.

## Also in the round-3 FAIL — chip-owned, no ruling needed, slice 11 in flight

- **G2 (important):** `read_service.go:337` — `Get` now nulls the ENTIRE `ICMSWorstCaseByUF` on a
  **cost-only** outage. `WorstCaseICMSPct` / `PriceNetBasis` depend on price + ceiling only
  (`icmsWorstCaseByUF:563-567`, `priceNetBasis:574-576`) and `belowMargin:520-522` already nils on a
  nil cost. Round 2 explicitly **certified the old behavior correct**. This throws away known-true
  facts because an unrelated fact is missing — ADR-17 licenses null for *unknown*, not for known.
  **This is the thing I recorded in `slice10-review.md:29-31` as a "BONUS fix… closing a latent
  gap." I was wrong; the §14 slice review is mine, so the misclassification is mine.** The cold gate
  caught what I had already blessed — which is the argument for the cold gate.
- **G3 (important):** `read_service_test.go:317` guards the matrix assertion behind
  `outage.ceilingErr &&`, so G2 is invisible to the suite. Round-2's own I2 criterion, turned on me.
- **G4 (blocking):** comment narration / PR-voice in production code —
  `read_service.go:172-177,:271-278,:318-324` narrate "Per the hub ruling (B1, PURE…)", "pins B2".
  `HARNESS.md:203-204` lists it in the AI-slop checklist: **any hit = REJECT the slice**.
- **G5 (important):** the ceiling WARN (`:117,:189,:329`) says "degrading fact to null" but fires
  before the `needsBelowMarginScan` branch (`:153`,`:192`) — a filtered request logs "degrading" and
  then 503s. The log states the opposite of the outcome, undercutting the honesty B2b restored.

## Governance drift (FYI, not a finding)

This worktree's `docs/` is behind the main repo: no `docs/REVIEW-STANDARD.md`, and `AGENTS.md` still
points at `docs/superpowers/HARNESS.md`. The binding copies live only in the main checkout. My
prompt-pack handed both reviewers a worktree path that does not resolve; both still returned
§14-conformant output (codex's cwd is the main root). Worth catching the worktree up, or pinning
reviewers to main-repo governance paths.

## State

`7f5a1b8c` committed, **not pushed**. Live stack at that SHA re-driven: **no regression** — C10 still
PASS (unknown 0.0%, 34/34 real MLB ids), all three fact-dependent filters 200 with the source
healthy (no over-correction). Evidence: `_gate-evidence/round-3/redrive-post-slice10.md`.

Slice 11 (G2/G3/G4/G5) is being implemented now — it needs no ruling. G1 waits on you.
