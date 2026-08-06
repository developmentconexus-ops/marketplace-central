# ADR-018: Mutation envelope is a protocol table plus an in-process poller

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision was ratified inside MIS-003
(`operator-cockpit-wireframe-replan`) as its "ADR-13", but the mission never produced a
standalone document — the rule lived only as a cited row in `mission.md` and in the
IC-03 interface contract it governs. It is reconstructed here from the 10 live citations
of that decision harvested at
`docs/architecture/decisions/_citations/adr-013-citations.md`, Assertion A1. The number
013 collided with an unrelated MIS-007 decision (reassigned separately to ADR-019); this
document keeps only the MIS-003 write-path sense and receives the new global number 018
per `docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md`.

## Context

Every provider write in the MIS-003 cockpit — price update, stock correction, link
resolution, pause, resync, attribute edit — was drawn in the wireframe as "preview e
protocolo": the operator sees what will change before it changes, and every change
leaves an auditable record. At the time this was decided, exactly one write path existed
in the codebase (`StockWriter`, synchronous, never run live), and the wireframe promised
five more. Built independently, five write surfaces become five ad-hoc write paths, each
re-inventing preview, idempotency, and audit — the redundancy the operator had already
forbidden.

The obvious alternative — a message bus or outbox with a background worker framework —
was rejected before it was built. The mission was single-tenant, dev-local, with one
write volume low enough that a database table and a goroutine could serve every
correctness property a queue would provide (ordering, at-least-once retry, durability)
without standing up infrastructure nobody would operate.

## Decision

**All provider writes route through one durable protocol: two Postgres tables
(`mutation_protocols`, `mutation_items`) and a single in-process poller that claims and
applies approved work. No outbox, no message bus, no background job framework.**

**§1 — One write path for every provider mutation.** Every write type — price, stock,
link, pause, resync, attribute edit — is a `type` value inside the same `mutation_protocols`
/ `mutation_items` shape, going through the same lifecycle and the same seven
provider-write gates, never a direct capability call from transport.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/mutation-envelope-interface-contract.md:17` —
> "Single durable provider-write path... `mutations` capability inside a new `mutations`
> module... final module home decided in ADR-13, not per feature."

**§2 — No external queue infrastructure.** Persistence is two tenant-scoped Postgres
tables; the applying step is a goroutine that claims work with `FOR UPDATE SKIP LOCKED`,
not a message broker or workflow engine.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/mutation-envelope-interface-contract.md:103` —
> "The applying poller is in-process (goroutine + DB claim with `FOR UPDATE SKIP LOCKED`);
> no external queue infrastructure this mission (ADR-13)."

**§3 — The table is the contract, not a placeholder for one.** The lifecycle, the seven
gates, and the item audit shape are structural (schema plus state machine), not
convention. Real queue infrastructure may replace the poller later behind the same
tables and API — the migration path is designed in, not deferred as debt.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-03-mutation-envelope-writes/F-01-protocolo-core/feature.md:48` —
> "No outbox, no message bus, no background framework — plain goroutine + ticker per
> ADR-13."

**§4 — The poller must survive a restart without re-doing or losing work.** A protocol
in `approved` or `applying` state at process death is picked up again by the next
poller instance; per-item idempotency keys make re-application safe rather than
duplicating provider calls.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:209` — "Q3 Reliability |
> poller resumes approved/applying after restart; idempotent re-apply proven; failed
> preserved never auto-mutated | ADR-13 |".

**§5 — Module home is a mission-level decision, not a per-feature one.** Which Go module
owns the `mutations` capability (`mutations` itself, or `listings`-adjacent) was fixed
once at the mission level so that no two features could each pick a different home for
the same write path.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:95` — "mutation envelope
> (module home fixed in ADR-13)".

**§6 — The route namespace follows the same designation.** The `/mutations` HTTP prefix
is mounted by whichever module ADR-018 (then ADR-13) designates as the owner, not
wherever a feature branch happens to add it.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/research/mutation-envelope-interface-contract.md:197` —
> "Server: `/mutations` prefix, mounted by the module that ADR-13 designates (M-03
> owns)."

**§7 — The rule is named in the owning milestone's own contracts, not just the mission
row.** M-03's milestone brief and its `mutation_protocols`/`mutation_items` feature both
cite the decision directly as their binding shape, not by inference from the mission
table.
> `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-03-mutation-envelope-writes/milestone.md:17` —
> "contracts: IC-03..., ADR-13/16; governance `contracts/governance/execution-lanes.json`
> 7 provider-write gates."

**§8 — The decision was still in force, unimplemented, at the next mission's baseline
audit.** MIS-004's repo baseline records the mutation envelope as designed but not yet
merged, citing the same table-plus-poller shape as the standing design to build against.
> `.mnfs/MIS-004-mvp-demo/research/repo-baseline-2026-07-17.md:66` — "NÃO está em main.
> Design em `.mnfs/MIS-003-.../mission.md`: ADR-13 protocolo table + in-process poller;
> lifecycle 8 estados (IC-03)."

## Rationale

A queue or bus buys ordering, retry, and durability that a single-tenant, dev-local
deployment does not need yet, at the cost of infrastructure someone has to run and
operate. A Postgres table already gives durability for free; `FOR UPDATE SKIP LOCKED`
gives single-claim semantics without a broker; a goroutine with a ticker gives a worker
without a framework. The five promised write surfaces sharing one envelope means the
seven provider-write gates (actor, idempotency, execute, resolved-link, policy,
source-timestamp, before/after-audit) are enforced once, structurally, instead of five
times by convention — the failure mode this decision exists to prevent is a sixth
surface that "forgot" one gate.

## Consequences

- The poller does not scale horizontally; a second instance racing the same table relies
  entirely on `FOR UPDATE SKIP LOCKED` for correctness. Accepted as a single-tenant
  trade-off (mission.md:118 accepted trade-offs row for ADR-13).
- Every new write type extends the `type` enum and one intent schema; the lifecycle,
  gates, and item shape are not allowed to vary per type (interface contract's
  Compatibility Rules).
- A future migration to real queue infrastructure is designed to slot in behind the same
  tables and API, but that migration has not happened and is not scheduled by this
  document.
- Because the module ships as one seam, any future write type surface that bypasses
  `/mutations` for a direct capability call is a regression against §1, not a new
  feature.

## Alternatives Considered

**Message bus / outbox pattern.** Rejected for this mission: it solves horizontal
scaling and cross-process delivery guarantees the single-tenant, single-poller
deployment does not need, at the cost of running and operating broker infrastructure.

**Per-write-surface synchronous calls (the pre-existing `StockWriter` pattern).**
Rejected: this is exactly the "N ad-hoc write paths" defect the envelope exists to
prevent — each surface would re-implement preview, idempotency and audit independently,
and five independent implementations of seven gates is five chances to miss one.

**Background job framework (e.g., a generic task-queue library).** Not evaluated as a
real alternative in the citations — the decision states directly that "no background
framework" was used, without recording a considered-and-rejected discussion beyond the
queue/bus framing above.
