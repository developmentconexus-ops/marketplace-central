# Slice 10 — §14 independent review — APPROVE

```yaml
reviewer: cold Claude sonnet subagent (independent context, implementer≠reviewer)
standard: REVIEW-STANDARD.md §14
candidate: slice10-candidate.diff (716 lines, 3 files, base 1f6b72d8)
ran: concurrently with L1 (§15)
verdict: APPROVE — zero blocking, zero important
```

## Verdict

**APPROVE.** No blocking, no important. Remaining findings are suggestion / nit / question only, none of
which gate a commit under §14.

## Reviewer's in-tree verifications (explicitly "confirmed in-tree, not merely claimed")

- **Fail-honest gate COMPLETE for BOTH Oracle facts** — ceiling pre-check (`scan:361-363`,
  `scanGroups:215-217`) AND cost discovered mid-scan via `enrich` (`scan:373-379`). Uniform for
  `has_exception` (both directions) and `exception=below_margin` via `needsBelowMarginScan:485-487`. New
  tests cover all three filter shapes × both outage sources × both `List`/`ByProduct`.
- **`matchesDependentFilter` genuinely still applied** on the happy path (`scan:383`, `scanGroups:240`) —
  the fix removed the passThrough escape hatch, not the filtering.
- **No-filter degrade-to-null preserved** and pinned by `TestReadServiceNoDependentFilterStillDegradesOnSourceUnavailable` (over-correction guard).
- **Classification strict, existing taxonomy only** — `isSourceUnavailable` → `IsReadErrorCode(err, ReadErrorSourceUnavailable)`; no string-matching. `context.Canceled`/`DeadlineExceeded`/plain errors propagate, pinned by three dedicated tests.
- **Tests PIN, not merely execute** — every dependent-filter-outage test `t.Fatalf`s if `err == nil`, so
  reverting to the old `passThrough` (which returned 200) would FAIL them. This is the exact criterion the
  round-2 gate used to catch the defective test.
- **`(error, error)` threading correct at every call site.** Reviewer found a BONUS fix nobody asked for:
  `Get` now nulls `ICMSWorstCaseByUF` on the combined `unavailable` (ceiling OR cost), closing a latent gap
  in the pre-slice-10 code that checked only `ceilingErr`.
- **Dead wrapper removal clean** — `unavailableListingPolicyReader` + its now-unused `listingsports` import
  gone from `root.go`; live wiring behaviorally unchanged.
- **Option (A) absent** — never nulls-margin-while-keeping-SQL-filter, exactly as the hub ruled.
- **Telemetry once per request, never per row** (guarded by early-return-on-unavailable in `scan`/`scanGroups` and the `if unavailable == nil` gate before the cost re-fetch).
- No `tenant_id` / layering regressions; no query construction touched.

## Non-gating findings (recorded, not fixed in this slice)

| Severity | Locus | Finding | Disposition |
|---|---|---|---|
| suggestion | `read_service.go:430` | `enrich` cost-path logs omit the `"op"` field the top-level call sites carry; an operator tracing a cost-outage line loses the calling operation. | **Accepted as follow-up.** Requires plumbing `op` through `enrich` (signature change) — out of the hub's R2 scope for slice 10. Logged here rather than silently dropped. |
| question | `transport/http_handler.go:68,113,214,274` | The catch-all maps ANY non-nil service error to `503 source_unavailable`, so a propagated `context.Canceled`/adapter defect surfaces identically to a genuine outage. | **Pre-existing, unchanged by this diff**, and the hub ruling explicitly blessed "no transport change". Mirrors the already-named deferred option (B) `service_degraded` signal. Tracked as the same cockpit-UI-milestone follow-up. |
| nit | `read_service.go:267,412` | `(error, error)` returns are transposable at a call site; named returns or a small struct would self-document. | Not a bug today (every call site handles `err` before `unavailable`). Readability only. |
