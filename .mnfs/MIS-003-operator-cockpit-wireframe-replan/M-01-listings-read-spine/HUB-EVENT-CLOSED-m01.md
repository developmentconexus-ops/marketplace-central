# CLOSED — M-01-listings-read-spine

```yaml
from: chip M-01-listings-read-spine
to: HUB (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
event: CLOSED
tip: 389b8b9a          # code tip 982d44e; c89fae3d + 389b8b9a are evidence-only on top
branch: mis-003/m-01-listings-read-spine
pushed: no
gate: round 5 merged PASS (both sides, no contradiction)
qa: P7 fresh browser QA @ c89fae3d — PASS
result: PASS
```

## The close ladder

**P6 dual gate — round 5, merged PASS.** Both sides independently PASS at `982d44e`. Zero blocking,
zero important, one non-gating nit. No contradiction to reconcile.

**P7 fresh browser QA — PASS** at `c89fae3d`, the exact SHA that passed the gate. Live endpoints
driven from a clean browser session against your stack.

## Gate reconciliation

| Round | SHA | Claude side | GPT side | Merged |
|---|---|---|---|---|
| 2 | `e2cde36` | Opus PASS | Sol FAIL (C07) | **PASS** — stale `.mnfs` wording superseded by ratified D-22 + binding OpenAPI |
| 3 | `c4e8ab91` | PASS on delta | **FAIL** — left the delta, read the adapter | **FAIL** — G1..G6 |
| 4 | `a6878dc6` | **PASS** — traced the wrap chain deeper | **FAIL** — enumerated the telemetry surface | **FAIL** — H1 |
| 5 | `982d44e` | PASS | PASS | **PASS** |

Round-5 merge, item by item:

| # | Was | Merged | Receipt |
|---|---|---|---|
| **G1** | blocking | **RESOLVED** per ruling (C) | `read_service.go:22-27` checks `errors.Is(ctx.Canceled/DeadlineExceeded)` before `IsReadErrorCode`. Both sides traced the **production** wrap and agree the cause survives it (`reader.go:516` → `read_error.go:28-33`). Residual → your task #2. |
| **G2** | important | **RESOLVED** | `:338` guard back to `ceilingErr != nil`. Confirmed live: 25-entry matrix served, only `below_margin_at_uf` null. |
| **G3** | important | **RESOLVED** | `read_service_test.go:320-332` pins it — revert `:338` and `:321` fatals. Anti-vacuity checked: `Price` is set, so the assertion has teeth. |
| **G4** | blocking | **RESOLVED** | Whole-file narration grep, both files → 0 hits. |
| **G5** | important | **RESOLVED** | Wording fixed in slice 11; last attribution gap (H1) in slice 12. |
| **G6** | nit | UNRESOLVED — non-gating | Standing. |
| **H1** | important | **RESOLVED** | `:430` now `"op", op`. `grep slog | grep -v '"op"'` → empty. Sol, who raised it, verified its own finding rather than assuming. |

## P7 — live, fresh browser, at the gated SHA

The QA surface is the **API**, not a page. M-01 is the read spine; there is no `/listings` route in
the nav (Dashboard, Catalog, Classifications, Marketplaces, Integrations, Product Links, Stock
Seguro, Orders, Pricing Simulator) because the cockpit that consumes this spine lands in later
MIS-003 milestones. Claiming a UI QA for a milestone with no UI would be theater, so P7's fresh-browser
requirement is met by driving the live endpoints from a clean browser session — real HTTP, real
Oracle-backed reads.

| ID | Result |
|---|---|
| **C04** | **PASS** — `limit=5` walk: 7 pages, **34 rows, 34 unique, zero dup/skip**, title ASC held across every boundary, `next_cursor` null at end |
| **C05** | **PASS** — 2 groups, `["15956", null]`, synthetic null group **last**, 33 unlinked grouped not dropped |
| **C06** | **PASS** — 400 `installation_required` · 400 `invalid_filter` · 400 `invalid_cursor` · 404 `listing_not_found`, each on status **and** `error.code` |
| **C07** | **PASS** — 33/34 rows `cost: null`, **zero** report `below_margin_worst_case: false`; summary counts them in additive `margin_unknown`=1, `below_margin_worst_case`=0 |
| **C10** | **PASS** — 34 rows, 34/34 real `MLB…`, **0.0% unknown** vs `<20%` bar |

Summary byte-stable across four drives (only `as_of` moves) — slices 11 and 12 were behaviour-neutral
on the read surface and the live data proves it.

`filter.has_exception=true` 33 · `=false` 0 · `filter.exception=below_margin` 0. 33+0=33 of 34, and
the browser named the unclaimed row directly: `MLB4735328201`, the one linked row, margin unknown, so
both buckets honestly decline it. ADR-17 working — forcing it into either would be the zero/default
substitution the ADR forbids.

## What C10 was actually worth

Round 1 recorded *"No code/composition/security defect found by any reviewer"* — against a
five-member cold crew **plus** a dual gate — and was Blocked only because C10 could not be driven.

Driving it for real found **20.6% unknown status**, a genuine adapter mapping gap, and the gate
rounds it triggered found four more defects. The mocks proved contract behaviour and nothing about
the running system. C10 was not a formality between a finished milestone and its close; it was the
criterion that made the milestone true. Worth keeping in the doctrine ledger next to *"mocks prove
contract behavior, never live integration."*

## Self-reported

- **G2 was a regression I approved.** My slice-10 review called it a *"BONUS fix nobody asked for."*
  It nulled the whole ICMS matrix on a cost-only outage — 25 known ceilings discarded because an
  unrelated fact was missing. I attached the review, so I own the misclassification.
- **Two miscounts in my evidence prose**, both caught after the fact: 16 vs **15** `slog` sites
  (round-5 cold Opus, checking a claim I invited it to distrust), and bare vs namespaced filter
  params in the round-4 re-drive table (caught at P7 by driving the live surface instead of trusting
  my own artifact — counts right, labels would have misled a reader into seeing a regression). Both
  corrected in place, left visible. **My narrative accuracy was the weakest link this milestone, not
  the code.**
- Implementation on **sonnet** across slices 8-12 per operator directive — HARNESS §1 deviation,
  logged in every commit and event.

## Yours now

- **Merge** — hub-owned seam. Branch `mis-003/m-01-listings-read-spine`, tip `389b8b9a`. **Not pushed.**
- **Task #2** — `internal_read` taxonomy refactor (G1 residual per ruling (C)): adapter/data defects
  still misclassified as source-unavailable because `safeOracleCause` flattens the cause before the
  application layer sees it. Misclassified but logged — no longer silent.
- **Task #3** — H4 (`read_service.go:414`, ceiling outage nils a healthy `Cost`) and H5
  (`cache.go:179-197`, `DoChan` can deliver caller A's cancellation to caller B). Both routed, not
  fixed — reviews verify, they do not generate scope.

`docker/dev/*.sh` remain modified in the worktree — your pre-existing dev-stack changes, excluded
from every commit of mine.

## Artifacts

- `validation-result.md` — round 2, supersedes round 1's Blocked (round 1 preserved in history and
  `_gate-evidence/round-1/`)
- `_gate-evidence/P7-browser-qa.md`
- `_gate-evidence/round-{3,4,5}/dual-gate-round*-verdict.md`
- `_gate-evidence/round-4/{slice11,slice12}-L0-report.md`, `redrive-post-slice11.md`
