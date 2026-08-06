# ADR-031: Products-mirror upsert keep-absent merge

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-006's `products_mirror` write path from
migration 0076 onward but was only ever cited by the two-digit token `ADR-04`, a number
that collides with three other unrelated mission decisions and with the pre-existing
three-digit document `004-integration-catalog-plugin-framework.md` (provider-plugin
self-registration — not this rule). It is reconstructed here from Assertion A2 of
`docs/architecture/decisions/_citations/adr-04-twodigit-citations.md` (4 live-code
citations). Every clause below is traceable to code or the migration that already
asserts it.

## Context

`products_mirror` is a current-state table: each source (`sankhya`, `xlsx`,
`catalogo_cliente`) upserts its own rows keyed by `(tenant_id, source, codigo_produto)`
every time it syncs. A sync run reads a full snapshot of the source, not a diff, so a
product that existed in the previous snapshot but is missing from the new one has to be
handled somehow — the sync cannot tell, from the snapshot alone, whether the product was
discontinued, temporarily filtered out upstream, or the read itself came back short.

## Decision

**A sync run never deletes a `products_mirror` row for a product that is missing from
its current snapshot. It flags the row `absent_in_last_snapshot = true`, stamps
`stale_since` once on the present-to-absent transition, and preserves the row's
last-known field values. If the product reappears in a later snapshot, the upsert clears
`absent_in_last_snapshot` back to `false` and clears `stale_since` to `NULL`.**

### The clauses

**§0 — The rule is named at both the package that writes it and the adapter that
triggers it.** The mirror writer package doc and the Sankhya sync entry point both cite
the rule by name at the point they implement it.
> `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go:4` — "It
> applies a full snapshot (ADR: full-snapshot v1) to products_mirror + the child
> products_mirror_stock_locations table with keep-absent merge semantics (ADR-04) and
> honest-NULL fields (ADR-17 — a missing value stays NULL, never a fabricated 0)."
> `apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go:56-60` —
> "Sync reads the current Sankhya snapshot (as-of today) and applies it to
> products_mirror via keep-absent merge (ADR-04) with honest-NULL fields (ADR-17).
> Full-snapshot v1: keep-absent needs the complete set and TGFEST has no change
> timestamp, so every run reads the whole active catalogue."

**§1 — Presence and absence are recorded on the row, never expressed by deleting it.**
The `products_mirror` schema carries `absent_in_last_snapshot BOOLEAN NOT NULL DEFAULT
false` and `stale_since TIMESTAMPTZ`, added specifically so a rebuild has somewhere to
record "not in this snapshot" without touching the row's existence.
> `apps/server_core/migrations/0076_products_mirror_active_source.sql:13` — "Keep-absent
> (ADR-04): a rebuild never physically deletes. When a product from the previous
> snapshot is missing from a new one, the merge logic (M-03/M-04) sets
> `absent_in_last_snapshot=true` + `stale_since=now()`; when it reappears the flag
> clears back to `false`."

**§2 — The sweep is a targeted UPDATE, not a delete, and it is scoped to the source that
ran.** `keepAbsentSQL` sets the two flag columns for rows of the given `tenant_id` and
`source` whose `codigo_produto` is not in the current snapshot's code list, and only if
they are not already flagged absent. It contains no `DELETE`.
> `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go:114-129` —
> "`keepAbsentSQL` implements ADR-04: rows present in the mirror for this source but
> absent from the current snapshot are flagged (never physically deleted, and their
> last-known values are preserved). ... Scope is THIS source only: since 0078 each
> source owns its own rows, so a Sankhya sync must never mark an xlsx-sourced product as
> absent."

**§3 — `stale_since` is stamped once, on the transition, not on every absent run.**
`COALESCE(stale_since, now())` leaves an already-flagged row's original absence
timestamp untouched across repeated syncs that keep finding it missing; only the
first transition into absence sets it.
> `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go:121-125`
> — `SET absent_in_last_snapshot = true, stale_since = COALESCE(stale_since, now())`.

**§4 — The upsert path is the reappearance path.** `upsertSQL`'s `ON CONFLICT ... DO
UPDATE` unconditionally sets `absent_in_last_snapshot = false` and `stale_since = NULL`
on every row present in the current snapshot, so a product that comes back in a later
sync clears its own absence in the same statement that re-writes its fields — there is
no separate "un-flag" step.
> `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go:91,110-111`
> — the `ON CONFLICT` clause sets `absent_in_last_snapshot = false, stale_since = NULL`
> alongside every other column.

**§5 — Consumers distinguish present-and-absent from never-seen by whether a row exists
at all.** A product the source never returned in any snapshot has no `products_mirror`
row (§1); a product the source used to return but no longer does has a row with
`absent_in_last_snapshot = true`. "Never seen" is the absence of the row itself;
"present, then absent" is a row that exists and says so. The two are not the same state
and the schema does not conflate them.

**§6 — An empty snapshot is refused, not applied.** `ApplySnapshot` returns
`ErrEmptySnapshot` and does not run the keep-absent sweep at all when the snapshot
carries zero rows, because a same-run empty read is more likely an upstream failure
(Oracle unreachable, query error) than "the source now has no products," and applying
the sweep over an empty code list would flag every row in the mirror absent.
> `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go:24-29` —
> "Applying it would mark the entire mirror absent (keep-absent sweep), which for an
> empty read is almost always an upstream failure ... rather than \"the ERP has no
> products\". We refuse it fail-closed instead of wiping the mirror's presence flags."

## Rationale

A full-snapshot sync has no diff to consume: it only knows what it read this time. Row
deletion would destroy the last-known field values (cost, price, stock) that a caller
may still need to show alongside an honest "not currently available" signal, and it
would make the transition itself unobservable — nothing would record when or that a
product left the snapshot. Flagging the row instead keeps the history in place and adds
exactly the two facts (that it's currently absent, and since when) that a consumer
needs.

## Consequences

- `products_mirror` grows monotonically per source; nothing is ever deleted from it by
  the sync path. Any operational cleanup of long-absent rows is a separate, undesigned
  concern.
- A consumer reading `products_mirror` must check `absent_in_last_snapshot` explicitly;
  reading the other fields without checking it risks presenting stale-since-absent data
  as current.
- The per-source scoping (§2) means one source's absence sweep can never mark another
  source's rows absent, which is what makes multiple sources sharing the same table
  (`sankhya`, `xlsx`, `catalogo_cliente`) safe.

## Relationship to ADR-027

ADR-027 ("Absent from a partial pull is not closed") governs a different mechanism over
a different table: it gates *when* the `listings` sweep is allowed to mark a row absent
at all, specifically to stop a run that aborted partway through a multi-tick pull from
treating "I didn't get to it" as "it's gone." This ADR governs a single-tenant,
single-call, full-snapshot write (`ApplySnapshot`) against `products_mirror`, which has
no multi-tick resumability to protect — the closest analogue is §6's refusal of an
*empty* snapshot, not a *partial* one. Both ADRs share the same underlying value (an
absence claim must be earned by the read that makes it, not assumed by default) but
protect against different failure shapes: ADR-027 against a truncated run over many
ticks, this ADR against an empty run over one call. Neither clause substitutes for the
other, and a `products_mirror` writer does not automatically inherit ADR-027's
run-completion gating just because both rules are about "absence."

## Alternatives Considered

**Physically delete a row on absence.** Rejected: it destroys the last-known values a
consumer may still want to show, and makes the presence-to-absence transition
unobservable — there would be nothing to query for "which products recently
disappeared."

**A separate `absent_products` table instead of flag columns on the same row.**
Rejected in favor of the simpler shape actually built: two nullable/defaulted columns on
the row itself keep the read a single-table query instead of a join, at the cost of
those two columns being meaningless before a sync has ever run for that row (mitigated
by `DEFAULT false` / `NULL`, which is itself honest-unknown per ADR-017).

## Unverified claims

None. All four live-code anchors cited above were read in the current source files and
match the clauses drawn from them.
