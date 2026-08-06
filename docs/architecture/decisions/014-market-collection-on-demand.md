# ADR-014: Market-reference collection is on-demand; runtime is a local docker dev stack

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-004's market-intelligence scope and was
enforced in planning review, but no document was ever written under this number. It is
reconstructed here from the 5 live citations of `ADR-11` that assert this meaning,
harvested at `docs/architecture/decisions/_citations/adr-011-citations.md` (Assertion
A2). The same `ADR-11` label was also used, unrelated, for two other decisions in two
other missions — see the registry at
`docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md`. During planning the
assertion itself moved between candidate numbers (`ADR-11`/`ADR-12`) across reviewer
rounds before ratification; this document fixes the final, ratified meaning.

## Context

MIS-004 needed to bring market-reference data (competitor pricing signals) into the
system without committing to a scraping pipeline, a scheduled collector, or a claim of
historical coverage the mission had no way to back. The infrastructure to run any
collection at all was new and had to land within the mission's timeframe, on a local
runtime.

## Decision

**Market-reference collection runs on demand only, against a local docker dev stack,
and makes no claim of retroactive or historical data.**

### The clauses

**§1 — Collection is on-demand, not scheduled.** There is no background job or cron
pulling market-reference data; collection happens when triggered.
> `.mnfs/MIS-004-mvp-demo/mission.md:95` — "ADR-11 Coleta on-demand + runtime docker
> local | ratificada | infra nova no prazo | sem retroativo: \"sem histórico ainda\"
> honesto | rehearsal completo no stack local"

**§2 — Runtime is a local docker dev stack.** The collection infrastructure for this
mission is new and runs locally; it is not a hosted service the mission stood up.
> `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:81` — "### ADR-11:
> Coleta on-demand; runtime demo local"

**§3 — No historical/retroactive backfill is claimed.** Because collection is on-demand
against new infrastructure, the mission does not claim any market history predating the
first on-demand run. The absence of history is stated honestly rather than backfilled or
approximated.
> `.mnfs/MIS-004-mvp-demo/mission.md:95` — "sem retroativo: \"sem histórico ainda\"
> honesto"

## Rationale

Building a scheduled or continuously-polling collector was out of scope for the
infrastructure timeline available; on-demand collection against a local stack was the
smallest thing that could produce a real, honestly-labeled market-reference signal within
the mission window. Declaring "no history yet" rather than inventing or backfilling a
history keeps the honest-absence discipline (ADR-017) intact at the one place a shortcut
would have been easiest to take — synthesizing a plausible-looking historical trend where
none was measured.

## Consequences

- Any screen or report that shows market-reference data must be able to represent "not
  yet collected" rather than assume a value exists.
- There is no mechanism in this mission to answer "what was the competitor price a month
  ago" — only "what is it as of the last on-demand collection."
- Because the runtime is a local docker dev stack, this collection path is not, by this
  decision alone, deployed as a standing service; operating it outside local development
  is a separate decision.

## Alternatives Considered

**Scheduled/continuous polling collector.** Rejected for this mission: building and
operating a scheduler within the available timeline was not feasible alongside the rest
of MIS-004's scope; on-demand collection was the achievable subset.

**Backfilling historical data from an external source.** Rejected: no source of verified
historical market-reference data was available, and fabricating a plausible-looking
history would have violated the honest-absence rule this decision was written under.
