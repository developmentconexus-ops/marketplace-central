# ADR-026: Sync scheduler phase vocabulary

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed the sync scheduler under the two-digit citation
`ADR-07` inside MIS-007, and was never given a document of its own. The two-digit number
collided with an unrelated MIS-004 decision and with the pre-existing
`007-godror-oci-oracle-runtime.md`, which is about the Oracle Go driver — a different
subject entirely. This document reconstructs the MIS-007 phase-vocabulary rule under its
own global number, from the 6 live-code citations harvested at
`docs/architecture/decisions/_citations/adr-07-twodigit-citations.md` (Assertion A1). It
states the rule as it stands **after** its 2026-08-01 amendment; the original, buggy form
is recorded under Amendments. Every clause below is traceable to code that already
asserts it. Nothing here is new policy.

## Context

`sync_state.cursor` is a nullable JSONB column (migration 0075) shared by every sync job
registered with the scheduler. After each run, the scheduler needs to decide whether that
run advances `last_incremental_at` or `last_full_sync_at` — a health-card-visible
timestamp (M-09's card reads `GREATEST(last_full_sync_at, last_incremental_at)`). It makes
that decision by peeking at a `"phase"` key inside whatever JSON the job's own cursor
happens to contain, without requiring every job to declare one.

## Decision

**The sync scheduler recognizes exactly three cursor phase values —
`backfill`, `incremental`, and `sweep` — and only `phase == "incremental"` causes a run to
be recorded as incremental. Every other value, including `sweep`, an absent or empty
phase, or a cursor that is not even a JSON object with a `phase` key, resolves to `false`
(a full/backfill-class run) with no error. A cursor is not required to carry a `phase`
field at all; jobs that never distinguish backfill/incremental/sweep are unaffected.**

**§1 — The vocabulary is exactly three values, and only one resolves true.**
`inferIncremental` peeks at the terminal cursor's `"phase"` field; `"incremental"`
resolves `true`, `"sweep"` and `"backfill"` both resolve `false` because a sweep is a
full/terminal walk of the catalog (M-04's backfill cursor lands on
`{"phase":"sweep",...}`), so it must advance `last_full_sync_at`, not
`last_incremental_at`.
> `apps/server_core/internal/modules/sync/application/scheduler.go:163-178` — "ADR-07
> (amended 2026-08-01) ratifies the phase vocabulary as backfill | incremental | sweep.
> Only \"incremental\" resolves true; \"sweep\" is a full/terminal sweep of listings ...
> so it resolves false like backfill, which is correct."

**§2 — Any unrecognized or absent phase is a tolerant `false`, never an error.** An empty
cursor, malformed JSON, a JSON value that is not an object, an empty-string phase, or an
unrecognized phase string all resolve to `false` with no error raised. This tolerance is
tested directly, including the amended `sweep` case.
> `apps/server_core/internal/modules/sync/application/scheduler_test.go:258-287`
> (`TestInferIncrementalTolerates`) — cases `"phase incremental"` → `true`;
> `"phase backfill"`, `"unrecognized phase sweep falls to tolerant default (ADR-07 has no
> sweep phase)"`, `"phase empty string"`, `"phase unrecognized"` → all `false`.

**§3 — `phase` is opt-in per job; legacy cursors without it are unaffected.** The
pre-existing `products` job's `ProductsCursor` carries no `phase` key at all and must keep
resolving to `false` with zero behavior change; this is the load-bearing backward-
compatibility case in the test suite, not an edge case.
> `apps/server_core/internal/modules/sync/application/scheduler_test.go:255-257,274-277` —
> "The \"legacy products shape\" case is the load-bearing one: the real products job's
> ProductsCursor has no phase key at all and must keep recording incremental=false with
> zero behavior change (F-03)."

**§4 — Jobs that use the vocabulary declare it with the shared constants, not ad hoc
strings.** The orders job's cursor uses `Phase` with values drawn from `phaseBackfill` /
`phaseIncremental`, documented explicitly as the ADR-07 vocabulary, because the scheduler
reads exactly that field to decide whether a cycle advances `last_incremental_at`.
> `apps/server_core/internal/modules/orders/application/orders_job.go:11-13,27-28` —
> "OrdersCursor é o estado persistido em sync_state.cursor para a entidade orders. O campo
> Phase usa o vocabulário do ADR-07 (backfill|incremental|sweep) porque o scheduler lê
> exatamente esse campo para decidir se o ciclo avança last_incremental_at."

**§5 — The listings backfill/sweep cursors are the two concrete phases that are not
`incremental`.** `backfillCursor.Phase` and the terminal `sweepCursor.Phase == "sweep"`
are the two shapes a listings run's cursor takes; the sweep cursor's own doc comment notes
that ADR-07's ratified vocabulary resolves only `"incremental"` to true, never `"sweep"`,
because ML has no reliable incremental listing search — a full daily sweep is the correct
steady state for that job, not a bug to fix.
> `apps/server_core/internal/modules/listings/application/backfill.go:80-98` — "ADR-07's
> ratified vocabulary only resolves \"incremental\" to true, never \"sweep\"."

## Amendments

**Original rule (pre-2026-08-01, buggy): `case "incremental", "sweep": return true`.**
`inferIncremental` treated both `"incremental"` and `"sweep"` as incremental. This was
wrong: a sweep is a full/terminal catalog walk, and recording it against
`last_incremental_at` instead of `last_full_sync_at` would freeze `last_full_sync_at`
while a sweep ran and make the M-09 health card (which reads
`GREATEST(last_full_sync_at, last_incremental_at)`) lie about when a full sync last
completed.

**Amended rule (2026-08-01, current — stated above as the Decision).** Only
`phase == "incremental"` resolves `true`; `"sweep"` now resolves `false`, matching
`"backfill"`. The fix is recorded directly in the function's doc comment
(`scheduler.go:165`, "ADR-07 (amended 2026-08-01)") and both `scheduler_test.go` table
cases were renamed to `"...falls to tolerant default (ADR-07 has no sweep phase)"` and now
assert `false`.
> Evidence trail: `.mnfs/MIS-007-ml-sync/M-02-sync-core-seam/_chip-m02/EVIDENCE.md:82-102`.

An earlier planning draft (`p5-claude-decomposition-audit-r03.md:105`) had additionally
proposed making `phase` mandatory on every cursor, which would have contradicted both the
narrower ratified mission scope and the legacy `ProductsCursor`'s lack of a `phase` field.
That proposal was reverted back to the narrow, opt-in form (§3) in
`p5-reconciliation-r03.md:30-34` before ratification; it never reached code.

## Rationale

A shared scheduler serving jobs with very different sync shapes (a paginated orders sync
with a real incremental watermark; a listings catalog walk with no incremental search at
all) needs one bit of information from each job's cursor — was this run incremental — 
without forcing every job to restructure its cursor shape around that one bit. Peeking at
an optional `phase` key, and resolving anything unrecognized to the safe default
(`false`, i.e. "counts toward full sync"), lets old jobs keep their cursor shape
unchanged while new jobs opt in by name.

## Consequences

- Any job wanting the health card to reflect an incremental run must set its cursor's
  `phase` to the literal string `"incremental"`; any other spelling (including case
  variants) silently falls through to `false` with no error to catch the mistake.
- Introducing a fourth phase value some jobs need requires either extending
  `inferIncremental`'s switch (visible, deliberate) or accepting that it resolves to
  `false` like every other unrecognized value — there is no way for a job to introduce a
  new phase that maps to `true` without changing the scheduler.
- The `products` job's `ProductsCursor` is permanently exempt from ever declaring `phase`
  under §3; adding a phase to it is a new decision, not something this rule already
  covers.

## Alternatives Considered

**Make `phase` mandatory on every cursor.** Considered during planning
(`p5-claude-decomposition-audit-r03.md:105`) and rejected before ratification: it
contradicted the legacy `ProductsCursor` shape already in production and would have
required a migration of every existing job's cursor for no behavioral gain, since the
tolerant default already handles absence correctly.

**Treat `sweep` as incremental (the pre-amendment behavior).** Rejected by the 2026-08-01
amendment: conflating a full catalog walk with an incremental update corrupts the
`last_full_sync_at` / `last_incremental_at` distinction that the M-09 health card depends
on.

## Unverified claims

None — every clause above matches a verified anchor in code.
