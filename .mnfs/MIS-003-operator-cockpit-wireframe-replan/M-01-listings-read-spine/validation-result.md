# Milestone Validation Result — M-01-listings-read-spine

```yaml
milestone: M-01-listings-read-spine
mission: MIS-003-operator-cockpit-wireframe-replan
round: 2
review_sha: c89fae3d          # code tip 982d44e; c89fae3d is evidence-only on top
diff_base: e2cde36            # round 1's review SHA
validation_level: QA-2
status: Pass
verdict: Pass
pushed: no
supersedes: round 1 (Blocked pending C10) — preserved in git history and _gate-evidence/round-1/
```

## Summary

M-01 (IC-02 listings read-spine) **passes**. Round 1 was **Blocked** for one reason: C10, the
live-provider-read criterion, could not be driven — no connected Mercado Livre installation, so it
was honestly recorded `could-not-drive` rather than fabricated. That environment now exists.

Driving C10 for real did not merely tick the box. It **exposed defects no lane had caught**, and
closing them took four more slices and three more gate rounds:

- The live read flooded **20.6% unknown status** — over the contract's `<20%` bar. Not a threshold
  quibble: a real adapter mapping gap. Slice 9 grew the canonical enum
  (`under_review`/`inactive`/`payment_required`/`not_yet_active`, migration 0037). Unknown → **0.0%**.
- Slice 10 made fact-dependent filters **fail honest** rather than serve a page with the filter
  silently dropped.
- Slice 11 fixed cancellation/timeout degrading facts to null, reverted a **regression I had
  personally approved** in the slice-10 review, and fixed telemetry that reported an outcome that
  had not happened.
- Slice 12 attributed the last unattributable error log to its read path.

Round 1's own words were *"No code/composition/security defect found by any reviewer"* against a
five-member cold crew plus a dual gate. Every defect above was real, and **only the live lane and
the gate rounds it triggered found them.** That is the finding of record: for this milestone, the
mocks proved contract behavior and nothing about the running system. C10 was not a formality
standing between a finished milestone and its close — it was the criterion that made the milestone
true.

## Per-criterion result

| ID | Criterion | Verdict | Evidence |
|----|-----------|---------|----------|
| M01-C01 | Refresh ingestion upserts and closes | **PASS** | `TestListingsRefreshSeedsIC02RowsAndClosesMissing`; composite tenant-scoped PK, close-only-removed. `F-01.../validation.md` |
| M01-C02 | Concurrent refresh guarded | **PASS** | `TestListingsRefreshRejectsConcurrentRunWithActiveID` + `TestOperationRunRepositoryBeginExclusiveIsAtomic`; atomic advisory-lock, 409 with active run id. Not re-driven live at P7 — firing a refresh at the real installation to re-prove an already-pinned guard is a live provider write this QA does not need. |
| M01-C03 | Unmappable → NULL honesty | **PASS** | `sales_30d IS NULL`, `listing_type` NULL, no 0/default (`listings_refresh_test.go:177-183`, `mapper.go:110-119`) |
| M01-C04 | List endpoint contract | **PASS** | **Live @ c89fae3d** (`_gate-evidence/P7-browser-qa.md`): `limit=5` walk → 7 pages, **34 rows, 34 unique, zero dup/skip**, title ASC held across every boundary, `next_cursor` null at end. Lane: `small_page_cursor_walk_and_JSON_null_contract`, `all_filter_keys`, `q_*`. |
| M01-C05 | By-product grouping (null-last) | **PASS** | **Live**: 2 groups, `product_id` = `["15956", null]`, synthetic null group **last**, 33 unlinked grouped not dropped. Lane: `by_product_cursor_walk_tie_order_and_null_last`. |
| M01-C06 | Error matrix (status + code) | **PASS** | **Live**, each asserted on status AND `error.code`: 400 `installation_required`, 400 `invalid_filter`, 400 `invalid_cursor`, 404 `listing_not_found`. 409 `refresh_in_progress` via lane. |
| M01-C07 | below_margin unknown honesty | **PASS** | **Live**: 33/34 rows `cost: null`, **zero** report `below_margin_worst_case: false`, zero anywhere in page; summary counts them in additive `margin_unknown` (=1), not `below_margin_worst_case` (=0). Lane: `null_cost_honesty_known_margin_and_summary`. Field name per ratified D-22 + binding OpenAPI; round 1's contract-wording reconciliation is applied (`validation-contract.md:127`). |
| M01-C08 | OpenAPI/SDK same-commit | **PASS** | `77845a59`, `1f0bbc66` each pair openapi.yaml + sdk-runtime same commit. Slice 9 pairs its 3-enum-site change likewise. Slices 10-12 touch no endpoint path. |
| M01-C09 | List perf (Q1) | **PASS** | `TestListingsReadPerformance2000` — nearest-rank **p95 3.56 ms** (ceiling 500), Index Only Scan `idx_listings_f02_title_key`, no Seq Scan, summary aggregate query count=1 (`slice11-L0-report.md`). |
| M01-C10 | Live read ingestion (live-provider lane) | **PASS** | **34 rows, 34/34 real `MLB…` ids, 0.0% unknown status** vs `<20%` bar (paused 17 / active 10 / under_review 7). Was **20.6%** before slice 9's enum growth — a real mapping gap, not threshold-gaming. `P7-browser-qa.md`, `_gate-evidence/round-4/redrive-post-slice11.md`. |

**Blocking failures observed: none.** Tenant scoping (Q2) holds — slices 10-12 touch no query
surface; lane subtest `tenant_isolation_all_read_paths_and_cursors` green.

## Gate record

| Round | SHA | Claude side | GPT side | Merged |
|---|---|---|---|---|
| 2 | `e2cde36` | Opus PASS | Sol FAIL (C07 only) | **PASS** — C07 divergence was stale `.mnfs` wording superseded by ratified D-22 + binding OpenAPI |
| 3 | `c4e8ab91` | PASS on delta | **FAIL** — left the delta, read the adapter | **FAIL** — G1..G6 |
| 4 | `a6878dc6` | **PASS** — traced the wrap chain deeper | **FAIL** — enumerated the telemetry surface | **FAIL** — H1 |
| 5 | `982d44e` | PASS | PASS | **PASS** — no contradiction |

Full merge tables with receipts: `_gate-evidence/round-3/dual-gate-round3-verdict.md`,
`_gate-evidence/round-4/dual-gate-round4-verdict.md`,
`_gate-evidence/round-5/dual-gate-round5-verdict.md`.

**The empirical case for the dual gate is that table.** Round 3's blocker came from Sol going
*outside* the delta to read the adapter that produces the error codes — a reviewer bounded strictly
to the delta would have passed a defect that lives outside it. Round 4's came from Sol reading a
surface Opus skipped entirely. Round 5 had Opus catch a miscount in *my own report*. Neither
reviewer was reliably the deeper one, and which side went deeper was not predictable in advance.
A single reviewer — either one — ships a defect here.

## P7 live runtime validation

- Surface: **API** (no `/listings` route in the frontend nav; the cockpit that consumes this spine
  lands in later MIS-003 milestones). P7's *fresh browser* requirement is met by driving the live
  endpoints from a clean browser session against the running hub stack — real HTTP, real
  Oracle-backed reads, at the exact SHA that passed the gate. Claiming a UI QA for a milestone with
  no UI would be theater.
- Outcome: **PASS** — C04, C05, C06, C07, C10 all driven live and green. `_gate-evidence/P7-browser-qa.md`.

## Deferred, with names — nothing closed by being forgotten

- **Hub board task #2** — `internal_read` taxonomy refactor (G1 residual, hub ruling (C)): an
  adapter or data defect is still misclassified as source-unavailable, because `safeOracleCause`
  flattens the cause before the application layer sees it. Cross-module seam (4+ consumers), not
  this chip's to widen. **Misclassified but logged — no longer silent.**
- **Hub board task #3** — H4 and H5 (raised by round-4 Opus, routed not fixed; reviews verify, they
  do not generate scope):
  - **H4** `read_service.go:414` — a *ceiling* outage skips the cost fetch and nils `Cost`, serving
    `cost: null` while the cost source is healthy. Arguably ADR-17 direction (b) one level up.
    Pre-base, pinned as intended, certified rounds 2-4.
  - **H5** `internal_read/adapters/cache/cache.go:179-197` — `group.DoChan` shares one loader's
    result, so caller A's `context.Canceled` can reach an uncancelled caller B. Honest either way;
    slice 11 changed B's outcome at a distance from degrade to propagate.
- **G6** (nit) — `(error, error)` uses `error` as a state carrier in the non-terminal position.
  Correctly threaded at every call site; slice 12 extends it rather than collapsing it.
- **S1** (suggestion, from round 2) — `NewBatchReader` unguarded / skips `TimingReader`. Verified
  safe (`ensureBatchAvailable` nil-checks). Unproven live per-page Oracle latency tracked separately.

## Self-reported

- **G2 was a regression I approved.** In the slice-10 review I recorded Get's guard change as a
  *"BONUS fix nobody asked for… closing a latent gap."* It was the opposite: it nulled the entire
  ICMS matrix on a cost-only outage, discarding 25 known ceilings because an unrelated fact was
  missing. ADR-17 licenses null for *unknown*, never for *known*. The §14 slice review was mine — I
  attached it, so I own the misclassification.
- **Two miscounts in my evidence prose**, both caught after I wrote them: 16 vs **15** `slog` sites
  (caught by the round-5 cold Opus reviewer, checking a claim I had invited it to distrust), and
  bare vs namespaced filter params in the round-4 re-drive table (caught at P7 by driving the live
  surface instead of trusting my own artifact — counts right, labels misleading). Both corrected in
  place, corrections left visible. My narrative accuracy has been the weakest link in this
  milestone, not the code.

Implementation ran on **sonnet** across slices 8-12 per explicit operator directive — a HARNESS §1
deviation (which binds implement to GPT-5.6 Luna high / Sol low). Operator authority wins; the
deviation is recorded in every commit and every event.

## Next handoff

- **Pass.** Code criteria met (gate round 5, both sides), QA passed (P7 live, fresh browser).
- **Not pushed.** Merge is a hub-owned seam.
- Next owner: **HUB** — merge, plus board tasks #2 and #3.
