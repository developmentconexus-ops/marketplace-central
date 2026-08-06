# ADR-027: Absent from a partial pull is not closed

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed the MIS-007-ml-sync mission's listings backfill/sweep
from F-02 onward but no document was ever written. It is reconstructed here from the
6 live-code citations of the two-digit token `ADR-06` (a numbering collision with an
unrelated pre-existing `006-oracle-internal-read-owned-by-mpc.md`; see
`docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md`), harvested at
`docs/architecture/decisions/_citations/adr-06-twodigit-citations.md`. Every clause below
is traceable to code or tests that already assert it. Nothing here is new policy.

## Context

The listings module keeps a per-installation mirror of a marketplace's catalog by pulling
pages from the provider and upserting rows locally. Enumerating and hydrating that catalog
takes multiple ticks, and any tick can fail: a 429, a deadline, a kill signal, a crashed
worker. Before this rule, the pull's completion writer (`ApplyCompletedPull`) treated every
row it had not seen in a run as gone, and closed it — a single unconditional UPDATE over
"everything not touched this run." That is correct only if the run actually walked the
whole catalog. A run truncated partway through walks only a prefix of it, and the same
UPDATE then closes every listing the run simply never got to: audit D-120 flaw F1, risk
R-B. A rate limit or a crash was silently indistinguishable, downstream, from a seller
delisting most of their catalog.

## Decision

**A listing's absence is only ever established by a run that reached its own end. A
partial run — one that stopped before enumeration and hydration fully drained — persists
whatever rows it did see and marks nothing absent.**

### The clauses

**§1 — Upsert and completion are two different calls, not one flag.** The old
`ApplyCompletedPull(rows, runStartedAt, complete bool)` scoped its "did this call see any
rows" safety pin to a single call, which broke on a resumable multi-tick run: a
legitimate empty terminal tick (the scroll cursor exhausted, nothing left to page in)
looked identical to "the caller claims completion over zero rows" and disabled keep-absent
for the whole run. The port now splits the write in two: `UpsertPulledRows` persists one
tick's rows and is called once per page; `MarkRunComplete` runs the keep-absent step once,
at the end, and only if the run's own enumeration and hydration fully drained.
> `apps/server_core/internal/modules/listings/ports/ingestion.go:50` — "MarkRunComplete
> runs the keep-absent step (ADR-06/IC-06) for a run whose enumeration + hydration fully
> drained."

**§2 — A tick's upsert never closes anything.** `UpsertPulledRows` writes only the rows it
received, verbatim; it contains no UPDATE that marks other rows absent. A run aborted after
any number of ticks — including the very first — has upserted whatever it saw and touched
nothing else.
> `apps/server_core/internal/modules/listings/adapters/postgres/repository.go:391` — "never
> a blanket close — ADR-06, audit D-120 F1, risk R-B: a run truncated by a 429/deadline/kill
> must never wipe listings it simply never got to."

**§3 — Completion itself is gated on the run having produced at least one row.** Before
`MarkRunComplete` marks anything absent, it checks whether this run's own timestamp
(`runStartedAt`, the same value stamped as `seenAt` on every tick's upsert) appears on any
row at all. If it does not — an enumerator that came back empty and terminal on its very
first tick — `MarkRunComplete` is a no-op. A run that produced zero rows is read as "found
nothing to say," never as "the whole catalog vanished."
> `apps/server_core/internal/modules/listings/adapters/postgres/repository.go:503` —
> "MarkRunComplete runs the keep-absent step (ADR-06/IC-06) for a run whose enumeration +
> hydration fully drained... an enumerator that produced nothing across every tick of a run
> must never be read as 'the whole catalog vanished'."

**§4 — Only when both gates hold does the keep-absent UPDATE run, and only over the exact
unseen set.** With the run confirmed to have touched at least one row, the UPDATE marks
`absent_since` only for rows of that installation whose `last_seen_at` is strictly older
than `runStartedAt` and that are not already marked absent — rows this run should have
seen (because they are in the same installation being fully walked) but did not.
> `apps/server_core/internal/modules/listings/adapters/postgres/repository.go:544` —
> `UPDATE listings SET absent_since = $3 WHERE tenant_id = $1 AND installation_id = $2 AND
> last_seen_at < $3 AND absent_since IS NULL`.

**§5 — Status is never inferred from absence.** A row's `status` column is written
verbatim from the provider payload on the tick that upserted it. Absence changes only
`absent_since`; it never synthesizes a status value such as `closed`.
> `apps/server_core/internal/modules/listings/adapters/postgres/repository.go:391` —
> "status is always written verbatim from row.Status; it is never inferred here."

**§6 — The CHECK that used to enumerate a fixed status vocabulary is gone, and the
lifecycle columns are additive and nullable.** `last_seen_at`/`absent_since` were added
specifically so a row's presence history is tracked independently of its provider status,
honest-unknown (ADR-017) until a run actually populates them.
> `apps/server_core/migrations/0090_listings_e3_fields_status_relax.sql:48` — "Lifecycle
> support (ADR-06): a row's absence is only ever learned at the end of a COMPLETE run
> (IC-06) — status is never inferred from absence."

The tests that pin this behavior assert the run-scoped mechanism directly rather than
re-deriving it: `fakeRunStore` in the application-layer test suite records every
`UpsertPulledRows`/`MarkRunComplete` call verbatim so a test can assert the pin survives
across ticks and resumes.
> `apps/server_core/internal/modules/listings/application/backfill_test.go:59` — "records
> every UpsertPulledRows/MarkRunComplete call verbatim... so tests can assert the run-scoped
> ADR-06 pin survives across ticks and resumes."

The integration suite that used to assert MASS-CLOSURE directly (an empty pull closes the
whole installation) was deleted along with the defect it existed to catch, and replaced by
tests of the behavior that took its place.
> `apps/server_core/internal/modules/listings/adapters/postgres/repository_integration_test.go:24`
> — "The MASS-CLOSURE assertions that used to live here... asserted the exact defect F-02
> kills (ADR-06, audit D-120 F1) and are gone — see
> TestApplyCompletedPull_AbortAfterPage1_NeverClosesUnseenRows and friends below."

## Rationale

A pull's failure mode is asymmetric: failing to write a row you did see is recoverable on
the next run, but wrongly closing a row you never got to is not recoverable without a
human noticing the seller's catalog "shrank" and reopening it by hand. Treating "not
present in this run's rows" as equivalent to "not present in the catalog" collapses two
different facts — *the pull didn't reach it* and *the provider removed it* — into one, and
the pull is not in a position to tell them apart mid-run. Only a run that reached its own
end, having enumerated the whole catalog, is in that position. Gating the keep-absent step
on run completion, and gating "any absence at all" on the run having touched at least one
row, are the two guards that keep a partial pull from asserting a fact it did not observe.

## Consequences

- A pull that aborts at any tick — for any reason — never marks a single listing absent.
  The cost is that a genuinely delisted item is not reflected as such until a run
  completes cleanly; this is accepted in exchange for never fabricating a closure.
- A run whose first tick returns empty-and-terminal is indistinguishable, at the database,
  from a run that never started. This mirrors ADR-017's known limit for honest absence:
  the guard proves the write path did not fabricate a closure, not that the sync ran
  successfully end to end.
- The keep-absent UPDATE is scoped per installation, per run's own timestamp — a resumable
  multi-tick run must pass the identical `runStartedAt` to every `UpsertPulledRows` call
  and to the final `MarkRunComplete`, or the touched-check and the keep-absent bound both
  break silently.

## Alternatives Considered

**Keep the single `ApplyCompletedPull(rows, runStartedAt, complete bool)` call and require
callers to pass `complete=false` on any non-terminal tick.** Rejected: this is exactly the
shape that produced the defect. A legitimate empty terminal tick and an aborted run both
look like "zero rows, nothing more coming" to a boolean-flag caller, and the two must be
told apart to decide whether keep-absent may run at all.

**Retry the whole pull from scratch on any tick failure, so a completed run is always
either fully fresh or entirely discarded.** Rejected: not evaluated against the mission's
rate-limit and cost constraints; discarding partial progress on every transient failure
would multiply provider calls without changing the correctness question this ADR answers.
