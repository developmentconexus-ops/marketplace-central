# ADR-017: Unknown is never zero — honest absence end to end

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed the codebase from its first citation and was
enforced in review, but no document was ever written. It is reconstructed here from the
1.378 live citations of `ADR-17`/`ADR-017`, harvested at
`docs/architecture/decisions/_citations/adr-017-citations.md`. Every clause below is
traceable to code, tests, or the published contract that already assert it. Nothing here
is new policy.

## Context

MPC's job is to report operational facts it does not own: stock from Sankhya, fees and
addresses from Mercado Livre, tax rates from an ERP matrix, costs from an xlsx export.
Every one of those facts can be **absent** — not yet synced, masked by the provider until
payment, structurally inapplicable, or simply missing from the source row.

A number type has no absent. `0` is a perfectly good `float64`, a perfectly good stock
count, and a perfectly good tax rate. So the cheapest thing for any layer to do with a
missing fact is to emit `0` and move on — and every layer that does so produces a screen
that looks answered, a margin that computes, and an aggregate that sums, all of them
wrong in a way that nothing downstream can detect. A fabricated zero is indistinguishable
from a measured zero once it is one function call away from where it was invented.

The failure is not hypothetical here. It is the single most-cited rule in the repository
precisely because it was violated repeatedly and caught in review each time.

## Decision

**An operational fact that could not be established is represented as absent — `nil`,
`NULL`, an omitted JSON key — at every layer it crosses. It is never substituted with
`0`, an empty string, a default, a neighbouring source's value, or a fabricated
timestamp.**

The rule is **two-sided**. A fact that was established and happens to be zero is written
as literal `0`. Suppressing a real zero into "unknown" is the same defect mirrored, and
is equally rejected.

### The clauses

Each clause is the rule applied to one situation. They are numbered so review and
checkers can cite them individually.

**§1 — No fabricated value.** A missing source fact becomes absent, never `0`, never a
default, never an empty struct.
> `pricing/domain/decimal.go:10` — "nil means unknown/unset, never a zero Money (ADR-17)"

**§2 — No suppressed zero.** A known zero is persisted and rendered as `0`. `0` is a
legitimate, distinct value, not a synonym for absent.
> `internal_read/adapters/oracle/sync.go:404` — "0 is itself a valid, distinct
> `grupo_icms` value in Sankhya."

**§3 — Every unknown is named.** When a composite value cannot be resolved, each
unresolved component is added to the unknown set by name. Absent-and-unnamed is a third
state this rule does not admit: it produces a nil that nobody can explain.
> `pricing/domain/icms.go:140-151`, `pricing/domain/decompose.go:87` — every missing
> source fact enters `ComponentesDesconhecidos` alongside `icms_saida`/`difal`/`pis_cofins`.

**§4 — "Does not apply" is not "unknown".** A component that is structurally
inapplicable on a given path is absent for a different reason: it does not enter the
unknown set and does not block a derived value.
> `pricing/domain/decompose.go:99-110` — "not an ADR-17 unknown — it never enters
> `ComponentesDesconhecidos` and never blocks `MargemValor`."

**§5 — No silent cross-source fallback.** When the configured source cannot answer, the
answer is absent. Another source's data is never served in its place, because the caller
cannot tell which source it received.
> `internal_read/adapters/routing/reader.go:20-21` — "there is no fallback between
> sources (ADR-17, fail honest)."

**§6 — No fabricated time.** A value without a real source timestamp leaves the
timestamp field nil. `time.Now()` at read time is a fabricated fact about when something
was true.
> `pricing/domain/tariff.go:28-30` — "config has no per-resolve timestamp and ADR-17
> forbids fabricating one."

**§7 — Unknown does not participate in comparison.** An absent value is excluded from
derived booleans and aggregates. It is never coerced into a number so a comparison can
run.
> `listings/domain/signal.go:78-100` — "custo desconhecido excluded, never counted below
> cost."

**§8 — No placeholder rows.** When a record cannot be resolved, no row is written.
A zero-valued or placeholder row is a fabricated fact with a primary key.
> `internal_read/domain/icms_matrix.go:42-43` — "NO row at all for this cell — never a
> zero-value or placeholder row."

**§9 — The wire says absent.** In DTOs, the OpenAPI contract and the SDK, an
unestablished value is an omitted key or an explicit `null`. Absence is part of the
published contract, documented per field, not an accident of serialization.
> 11 clauses of `contracts/api/marketplace-central.openapi.yaml` already state this
> (L5721, 5727, 5791, 5797, 5802, 5810, 5832, 5843, 5856, 5913, 5920), e.g. L5810:
> "Every amount is optional/nullable: absent (nil) means the cost could not be honestly
> sourced, never a fabricated zero (ADR-17)."

**§10 — Opaque stays opaque.** A provider-supplied string whose value space we do not
own is rendered exactly as received. Mapping it onto a local enum invents a fact about
the provider.
> OpenAPI L5843 — "Document type string EXACTLY as the provider returns it. Opaque:
> never assume/map to a CPF/CNPJ enum — render as-is."

**§11 — Provenance is a fact too.** A decision, matched anchor, or rule name recorded in
an audit trail must reflect what was actually established. A plausible-looking
attribution nobody verified is a fabricated fact in the one place built to be trusted.
> `product_links/application/resolution_service.go:210-212` — "an anchor matched
> nothing — a fact nobody established."

**§12 — Stale presented as fresh is fabrication.** A cached answer served as current
asserts a time-of-truth that is not true.
> `product_links/adapters/postgres/link_candidate_repo.go:53-54` — "a stale answer
> presented as fresh is worse than none."

**§13 — Ingestion is lenient, not inventive.** On the lenient import path, an optional
field absent from the source file is accepted as honest-unknown with a warning. It is
neither coerced to `0` nor used to reject the row. The strict validation path is
unchanged and still rejects the same absence.
> `erp_import/adapters/xlsx/parser.go:28-42`, `erp_import/application/import_service.go:94-96`.

## Rationale

Absence is information. It says "ask the source again", "this sync has not run", "the
provider will not tell us yet" — all of which are actionable. `0` says "we measured, and
the answer is nothing", which is a different claim and, when invented, a false one. The
two are not interchangeable, and only one of them is recoverable downstream: you can
render an unknown honestly at any layer, but you can never recover an unknown from a
zero that has already been written.

The rule is deliberately end-to-end rather than per-layer. A domain that models absence
correctly is undone by a DTO with a non-pointer `float64`, and a correct DTO is undone by
a frontend that renders `null` as `0`. Every crossing is a place the fact can be
laundered, so every crossing is covered.

## Consequences

- Money, rates, stock and timestamps on unresolved paths are pointer/nullable types
  across domain, DTO, OpenAPI, SDK and FE. This is the cost of the rule and it is
  accepted.
- Screens must have a rendering for unknown (`—` plus a freshness indicator) and must
  never fall back to `0`. The FE carries `UnknownValue`/`FreshnessIndicator` for this.
- Aggregates are partial by design: a total over a set containing an unknown is itself
  unknown, or explicitly reports what it excluded.
- **Known limit:** honest absence is indistinguishable, at the point of rendering, from
  a correctly-working pipeline whose upstream table is empty. §1 working is not evidence
  that the sync ran. Verifying that upstream produced rows is a separate obligation.
  (`internal_read/application/icms_matrix_job.go:47-49`,
  `docs/engineering/defect-class-catalog.md:354-357`.)
- `float64` on the money path violates §1 structurally: a non-pointer float cannot
  express absence. The module-protocol mission carries the remaining migration.

## Alternatives Considered

**Sentinel values (`-1`, `NaN`, `9999`).** Rejected: a sentinel is a fabricated value
that every consumer must be told about out of band, and any consumer that is not told
computes with it. It moves the defect from "silently wrong" to "wrong unless everyone
remembers".

**Zero plus a separate `known` flag.** Rejected: two fields that must be read together
will be read apart. The flag is droppable at every boundary the value crosses, and
dropping it degrades to the exact defect this rule exists to prevent.

**Per-layer discretion (domain strict, wire lenient).** Rejected: the wire is where the
value is consumed. Laundering the unknown at the last boundary makes the strictness of
every layer before it decorative.
