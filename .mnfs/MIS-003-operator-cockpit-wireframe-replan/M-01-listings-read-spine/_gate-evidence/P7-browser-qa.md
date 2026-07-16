# P7 — fresh browser QA @ c89fae3d

```yaml
tip: c89fae3d          # code tip 982d44e; c89fae3d is evidence-only on top
gate: round 5 merged PASS (both sides, no contradiction)
env: hub-owned docker compose; hub sent ENV-READY after restart
backend: marketplace-central-backend-1, up, "server starting on :8080" 13:39:56, applied 0 migration(s)
surface: http://127.0.0.1:8080 driven from a FRESH browser session (no curl, no reused client, no cache)
installation_id: inst-mercado_livre-d373dc64-577f-4950-b11b-b90244f30cb2
as_of: 2026-07-16T13:41:42.000513923Z
verdict: PASS
```

## Why the QA surface is the API, not a page

M-01 is `validation_level: QA-2` — the listings **read spine**. Its Required Outcome is
*"five endpoints serving canonical shapes"*; there is no `/listings` route in the frontend nav
(Dashboard, Catalog, Classifications, Marketplaces, Integrations, Product Links, Stock Seguro,
Orders, Pricing Simulator), because the operator cockpit that consumes this spine lands in later
MIS-003 milestones. So P7's *fresh browser* requirement is met by driving the live endpoints from a
clean browser session against the running stack — real HTTP, real Oracle-backed reads, at the exact
SHA that passed the gate. Claiming a UI QA for a milestone with no UI would be theater.

## Criteria exercised live

| ID | Criterion | Result |
|---|---|---|
| **C04** | List endpoint contract — cursor walk, sort, JSON null | **PASS** — `limit=5` walk: 7 pages, **34 rows, 34 unique, zero duplicates, zero skips**; title ASC order held across every page boundary; `next_cursor` null at end. |
| **C05** | By-product grouping, null group last | **PASS** — HTTP 200, 2 groups, `product_id` = `["15956", null]`; the synthetic null-product group is at index 1 of 2 — **last**, and the 33 unlinked listings are grouped, not dropped. |
| **C06** | Error matrix | **PASS** — missing installation → **400 `installation_required`**; `filter.bogus=x` → **400 `invalid_filter`**; garbage cursor → **400 `invalid_cursor`**; unknown composite id → **404 `listing_not_found`**. Each asserted on status **and** `error.code`. |
| **C07** | below_margin unknown honesty | **PASS** — 33 of 34 rows carry `cost: null`; **zero** of them report `below_margin_worst_case: false`; **zero** rows anywhere in the page report `false`. Summary counts them in the additive `margin_unknown` counter, not in `below_margin_worst_case`. |
| **C10** | Live provider read | **PASS** — 34 rows, all `MLB…` real provider ids, **0.0% unknown status** against a `<20%` bar (paused 17 / active 10 / under_review 7). |

`C01`, `C02`, `C03`, `C08`, `C09` are integration-lane and engineering criteria, already evidenced
at their own SHAs (`slice9-L0-report.md`, `slice11-L0-report.md`, feature `validation.md`). C02's
409 path is a **write** (`POST /listings/refresh`); it is pinned in the integration lane and was not
re-driven live here, because firing a refresh at the real installation to re-prove an already-proven
guard is a live provider write this QA does not need.

## Summary served live, byte-stable across four drives

```json
{"total":34,"active":10,"paused":17,"exceptions":{"sync_error":0,"stale":0,"unlinked":33,"below_margin_worst_case":0,"margin_unknown":1},"as_of":"2026-07-16T13:41:42.000513923Z"}
```

Identical to the slice-9, slice-10 and slice-11 drives apart from `as_of`. Slices 11 and 12 were
behaviour-neutral on the read surface and the live data proves it.

## The fact-dependent filters, and the row neither claims

```
filter.has_exception=true      -> 33 rows
filter.has_exception=false     ->  0 rows
filter.exception=below_margin  ->  0 rows
filter.exception=sync_error    ->  0 rows
filter.exception=unlinked      -> 33 rows
all=34  true=33  false=0  claimed=33
unclaimed by either bucket: ["MLB4735328201"]
```

33 + 0 = 33 against 34, and the browser named the missing row directly rather than leaving it to
inference. `MLB4735328201` is the single linked row; its margin is unknown, so `has_exception=true`
(which needs a *provable* exception) and `has_exception=false` (which needs
`BelowMarginWorstCase != nil && !active`) both honestly decline it. ADR-17 working as designed: a
row whose exception status is unknown belongs to neither answer, and forcing it into one would be
the zero/default substitution the ADR forbids. Its `Get` confirms the shape: `cost` known
(`91.57 BRL`), `below_margin_worst_case: null`, and the ICMS matrix present with 25 entries —
the G2 revert holding on live data.

## Something QA found that the contract does not specify

The contract's C06 row names `filter.bogus=x`, which passes. Driving the surface more widely
surfaced the shape rule behind it, which is worth recording because it nearly produced a false
finding of my own:

- **Namespaced filters validate values and fail honest.** `filter.exception=nonsense`,
  `filter.has_exception=maybe`, `filter.status=nonsense`, `limit=-5` → **400 `invalid_filter`**,
  every one. No silent filter-drop, no unknown-becomes-default. This is exactly the fail-honest
  behaviour slice 10 established, and it holds at the transport edge too.
- **A bare key is not a lenient alias — it is an unrecognized key, ignored.** `status=paused`
  returns all 34 rows; `filter.status=paused` returns 17. Ignoring unknown query keys is ordinary
  HTTP behaviour and not a defect.

I first read the bare-key 200s as a silent filter-drop — an operator typing `?exception=below_margins`
getting the full list back and believing it filtered. That would have been an honesty defect of the
same family as ADR-17. It is not: the bare key was never a filter. I checked before reporting, and
the check is what turned a would-be blocker into a footnote. Recorded because the near-miss is the
useful part.

## A correction to my own evidence

`_gate-evidence/round-4/redrive-post-slice11.md` wrote the filter keys **bare**
(`has_exception=true`). The counts in it are right — re-verified here through the browser under the
namespaced shape, 33/0/0 with `MLB4735328201` unclaimed, reproducing exactly — but the labels were
wrong, and a reader who copied them would have got 34 rows back and concluded a regression. Fixed in
place, with the correction left visible rather than silently edited.

That is the second miscount of mine caught after the fact this milestone (the first: 16 vs 15 `slog`
sites, caught by the round-5 cold Opus reviewer). Both were in evidence prose rather than in code,
and both were caught — one by a reviewer I invited to distrust me, one by driving the live surface
instead of trusting my own artifact. Worth the ledger: my narrative accuracy has been the weakest
link in this milestone, not the code.

## Verdict

**P7 PASS.** Every live-exercisable criterion green from a fresh browser at the gated SHA. No
regression, no over-correction into a blanket 503, no honesty violation on any read path.
