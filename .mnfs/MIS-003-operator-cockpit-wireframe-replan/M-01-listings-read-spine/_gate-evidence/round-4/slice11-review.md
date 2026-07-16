# Slice 11 — §14 independent review — APPROVE

```yaml
reviewer: cold Claude sonnet subagent (independent context, implementer≠reviewer)
standard: docs/REVIEW-STANDARD.md §14 (MAIN checkout path, per hub governance ruling)
candidate: slice11-candidate.diff (410 lines, 2 files, base 7f5a1b8c)
ran: concurrently with the integration lane (§15)
verdict: APPROVE — zero blocking, zero important
```

## Verdict

**APPROVE.** G1–G5 correctly implemented and verified against the real wrap/unwrap chain; no
must-not-break invariant regressed; no narration/PR-voice comments remain anywhere in either file
(checked whole-file, not just the diff hunks). Only non-gating findings remain.

## Reviewer's in-tree verifications (confirmed against source, not accepted as claims)

- **G1 works through the production wrap.** `wrapOracleError`/`safeOracleCause`
  (`reader.go:508-519`) preserve `context.Canceled`/`DeadlineExceeded` as `ReadError.Cause`, and
  `ReadError.Unwrap` (`read_error.go:28-33`) returns `Cause` — so `errors.Is` genuinely walks the
  production shape. The new `wrappedCause` helper (`read_service_test.go:206-210`) builds exactly
  that shape (correct code + real context cause), **not a fake-code shortcut**.
- **G2 restored.** `Get`'s guard is `ceilingErr != nil` (`read_service.go:338`). Reviewer traced the
  cost-outage path through `enrich` and `icmsWorstCaseByUF` and confirmed the full matrix survives
  with only `BelowMarginAtUF` nulled.
- **G3 pins.** New assertions (`read_service_test.go:237-249`) genuinely fail if G2 is reverted
  (traced by hand, independently of the worker's executed proof).
- **G4 clean.** Whole-file grep for narration/PR-voice: **no hits** in either file.
- **G5 truthful at every site.** WARN/ERROR placement checked at Summary, List, ByProduct, Get,
  `scan`, `scanGroups`, `enrich` — all fire only once the outcome (degrade vs fail) is known, none
  inside a per-row loop.

## The worker's self-reported caveat — independently CONFIRMED, and it holds

The worker disclosed that its 2 `dependent_filter_scan` subtests pass even with the `errors.Is`
check removed. The reviewer verified the reasoning against the code and **confirms it**: with G1
reverted, `isSourceUnavailable` would misclassify the wrapped context cause and `enrich` would
degrade instead of propagate — but `scan`/`scanGroups` return `unavailable` as an error whenever it
is non-nil (`read_service.go:358-361`, `:375-378`, `:214-217`, `:231-234`) **regardless of why**. So
those two subtests cannot pin G1 there; the behavior is identical either way. The pinning force
comes entirely from the **16 flat-path/Summary/List/ByProduct/Get subtests**, where the top-level
guard or `enrich`'s propagate-vs-degrade branch is the only thing between a 200 and an error.

**Not a coverage gap — correctly characterized by the worker.** Recorded because a self-reported
caveat that survives independent verification is evidence, not noise.

Worker's executed pin proof on the 16 that do pin: removing the `errors.Is` check produced
`read_service_test.go:543: expected context canceled to propagate, got summary={...MarginUnknown:<nil>...} err=<nil>`
— a 200 with null facts on a cancelled context. Restored → green.

## Non-gating findings

| Severity | Locus | Finding | Disposition |
|---|---|---|---|
| suggestion | `read_service.go:467-469` | `unavailableFact`'s doc comment reasoned backwards — it said "the ceiling fetch never happened", but the ceiling fetch *did* happen and failed (that is how `ceilingErr` got set); it is the **cost** fetch that `enrich` skips when `unavailable` is already non-nil (`:414`). The returned label was correct; only the comment's reasoning was wrong. | **FIXED by the milestone owner** before commit. Reworded to state that a non-nil `ceilingErr` means the ceiling fetch already failed and `enrich` skipped the cost fetch. Comment-only change; build + vet + `listings/application` tests re-run green after it. |
| question | `read_service_test.go:381-400` | The "ceiling" half of the new dependent-filter scan/group-scan subtests never reaches `scan`/`scanGroups` — `List`/`ByProduct`'s top-level `ceilingErr != nil && !isSourceUnavailable(ceilingErr)` guard returns before the `needsBelowMarginScan` branch — so they duplicate the plain "list"/"by product" subtests above them. Not a correctness bug. | **Accepted as-is.** Redundant coverage, not wrong coverage; the subtests still assert the required invariant. Removing them buys nothing and risks the pinning set. Recorded, not actioned. |
