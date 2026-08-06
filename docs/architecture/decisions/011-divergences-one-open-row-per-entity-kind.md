# ADR-011: Divergences are one open row per (entity, kind), detected at ingest

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed MIS-007's shared `divergences` ledger from its
first citation but no document was ever written under a stable number. It is
reconstructed here from Assertion A1 of the 12 live citations of `ADR-10` gathered in
`docs/architecture/decisions/_citations/adr-010-citations.md`. Every clause below is
traceable to code, migration, or mission text that already asserts it.

## Context

Two independent producers detect mismatches against the same set of entities on
different cadences: a daily stock producer (listings ingest, M-05) and a five-minute
order producer (order ingest, M-06). Both need to record "expected vs. observed"
disagreements — stock count drift, tariff (fee) drift — for two frontend consumers
(`/anuncios` badge and filter, `/precos` warning) to read. Without a single decided
shape, the daily producer would naturally append one event per detection while the
five-minute producer would naturally want a single row it updates in place; badges and
filters built against one shape would be meaningless against the other, and the
five-minute producer's append-only table would grow without bound.

## Decision

**A divergence is tracked as at most one open row per `(tenant_id, provider, entity_type,
entity_id, kind)`. Detection happens only at ingest time, never on read. Opening or
refreshing a row requires both sides' observation timestamps to be present; a later
convergent observation resolves the existing open row in place rather than deleting it or
opening a new one.**

**§1 — At most one open row per (entity, kind); "open" is `resolved_at IS NULL`, not a
separate status column.** A partial unique index enforces this at the database, and the
column comment states the state is carried by the NULL itself.
> `apps/server_core/migrations/0087_divergences.sql:9-13` — "Row \"open\" = resolved_at
> IS NULL; there is no separate status column, the NULL is the state (IC-02 Enums And
> Statuses). Convergence sets resolved_at and the row stays as history — never a physical
> delete. A later flap opens a brand new row (new detected_at); the partial unique index
> below permits exactly that, one open row per (entity, kind) at a time."

**§2 — Detection happens at ingest, never on read.** Both producers call `Evaluate()`
into the same table at their own ingest cadence; nothing computes a divergence lazily
when a screen is opened.
> `apps/server_core/migrations/0087_divergences.sql:4-7` — "Detection always happens at
> ingest (ADR-10): a daily stock producer (M-05) and a 5-minute order producer (M-06)
> each Evaluate() into the SAME table instead of inventing incompatible shapes
> (append-events vs. one-open-row). No producer writes into this table yet."
> `.mnfs/MIS-007-ml-sync/research/divergences-interface-contract.md:37` — "Evaluate | TODO
> ingest da entidade (ADR-10: detecção no INGEST, nunca no read) | esperado + observado +
> fontes + timestamps de observação | upsert de row aberta OU resolve".

**§3 — Both sides' observation timestamps are mandatory, never defaulted.** The recorder
rejects an evaluation that is missing either timestamp rather than substituting a
fabricated one — the same non-fabrication discipline `ADR-017` states generally, applied
here to mitigate false positives from comparing observations taken at different times.
> `.mnfs/MIS-007-ml-sync/M-05-listings-fees-divergence/validation-contract.md:64` —
> "Expected: row aberta com timestamps dos 2 lados NOT NULL (ADR-10); após convergência,
> `resolved_at` preenchido na MESMA row (one-open-row)."
> `apps/server_core/internal/modules/divergences/adapters/postgres/recorder.go:83-88` —
> "Order of checks: (1) no linkage → skip ... (2) both observation timestamps are
> mandatory, named-error rejected, never defaulted; (3) compare expected vs observed
> within Kind's tolerance and either auto-resolve an open row (convergent) or open/refresh
> one (divergent)."

**§4 — `detected_at` is immutable once a row is opened; only `last_evaluated_at` moves on
a later divergent re-evaluation.** The upsert statement's `DO UPDATE SET` list
deliberately omits `detected_at`.
> `apps/server_core/internal/modules/divergences/adapters/postgres/recorder.go:54-56` —
> "detected_at is only ever written by the INSERT arm — the DO UPDATE SET list
> deliberately omits it, so a later divergent re-evaluation of the same open row can never
> move it (IC-02 Persistence Expectations)."

**§5 — Convergence resolves the existing row; it is never a delete and never a new row.**
> `.mnfs/MIS-007-ml-sync/mission.md:218-221` — "ADR-10 Divergences =
> one-open-row-per-(entity,kind), upsert. Chave natural `(tenant_id, provider,
> entity_type, entity_id, kind)` com no máximo 1 row aberta; `expected_*`/`observed_*` +
> timestamps de observação dos DOIS lados NOT NULL (mitiga falso-positivo R-5); detecção =
> upsert; convergência grava `resolved_at`."

## Contradictions

**Number collision across three unrelated rules.** The token `ADR-10` (and its
zero-padded spelling `ADR-010`) also names a MIS-004 DIFAL-single-source rule
(reconstructed separately as `012-difal-single-source-in-pricing.md`) and a thinner
MIS-001/MIS-003 "mocks never claim live integration" rule, out of scope for this
reconstruction. Nothing in the citing text disambiguates the number by itself; only the
mission the citing file belongs to does.

## Exceptions / carve-outs

No carve-out narrows this rule beyond the timing constraint that is itself part of the
decision: detection must happen at ingest, explicitly "nunca no read" (never on read).
There is no partial-guard variant recorded for divergences under this number.

## Rationale

A shared ledger between two producers on different cadences only stays coherent if the
shape is decided once, centrally, before either producer writes to it — otherwise each
producer optimizes for its own cadence (append vs. upsert) and the two readers
(`/anuncios`, `/precos`) inherit incompatible semantics. Detecting only at ingest, rather
than lazily at read time, keeps the comparison anchored to the exact pair of observations
that were actually made together, rather than re-comparing against whatever the current
state happens to be when someone opens a screen — which would silently change what
"divergent" means depending on when it is checked. Requiring both timestamps and
rejecting their absence is the direct application of `ADR-017`'s no-fabricated-value rule
to comparison inputs: a comparison built on a guessed timestamp is a comparison with a
fabricated premise.

## Consequences

- The `divergences` table never grows without bound from repeated re-detection of the
  same unresolved mismatch — a re-evaluation updates the existing open row rather than
  appending.
- History is preserved: a resolved row is never deleted, so a later flap on the same
  entity opens a genuinely new row with its own `detected_at`, distinguishable from the
  earlier resolved one.
- Producers cannot evaluate without a link between their entity and the compared side; an
  unlinked entity is skipped, not evaluated against a guessed counterpart.

## Alternatives Considered

**Append-only event log (one row per detection).** Rejected as the daily producer's
natural shape: it would make "how many divergences are open right now" a query over the
whole history rather than a direct read, and gives the five-minute producer an unbounded
table.

**Detect divergences lazily on read (compare current mirror state to current provider
state when a screen opens).** Rejected: this decouples the divergence from the specific
observation pair that triggered it, makes "divergent" a moving target that depends on
when it's checked, and cannot honor the two-observation-timestamp discipline in §3.

## Unverified claims

None. All anchors cited for Assertion A1 were read and confirmed to match verbatim at the
cited line (or within 1-2 lines of surrounding context for multi-line comments).
