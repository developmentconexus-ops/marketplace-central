# ADR-035: Architecture Rebaseline governs target design during D0–D9

**Date:** 2026-08-13  
**Status:** accepted

## Context

Marketplace Central contains several architectural eras at once:

- 21 legacy `internal/modules` directories;
- the newer `internal/contexts` shape already used by `catalog` and `listings`;
- provider adapters moving from `connectors` to `internal/adapters`;
- a large historical migration chain;
- a hand-authored API/SDK/routing model;
- several dated plans/specs/handoffs/wikis describing incompatible target structures.

The platform has no production users whose compatibility requirements force preservation of those historical shapes.

The operator has explicitly chosen to perform a full technical architecture deep dive before continuing implementation. Merely choosing folder names/context candidates is insufficient: identity, data ownership, internal communication, events, external integrations, API, frontend, runtime and golden flows must be designed deeply enough that implementation does not need to invent architectural decisions locally.

Keeping every previous `accepted` structural ADR as equal target authority would defeat that program: some decisions were correct for an earlier module architecture but directly answer questions D0–D9 is now required to re-adjudicate.

Deleting ADR history outright would also be wrong: code/citations need provenance and old decisions remain useful evidence of why the current system has its shape.

## Decision

### 1. D0–D9 is the governing target-design program

The required sequence is:

1. **D0** — Current State & Authority Baseline
2. **D1** — Context Adjudication
3. **D2** — Identity & Data Ownership
4. **D3** — Internal Communication & Event Model
5. **D4** — External Integration Contracts
6. **D5** — HTTP/API Contract
7. **D6** — Frontend Contract
8. **D7** — Runtime / Scheduler / Transactions / Outbox
9. **D8** — Golden Flow Simulation
10. **D9** — Adversarial Global-Maximum Review

Only after D9 acceptance may the repository create the implementation DAG and implementation plan for the rebaseline.

Product implementation is not authorized by approval of D0–D8 alone.

### 2. Git is the archive; active documentation has one path

Historical plans, handoffs, evidence snapshots, wikis and old design documents are removed from the active tree after their still-valid principles are absorbed.

There is no `old/`, `archive/` or parallel legacy roadmap.

Current program status lives only at `docs/engineering/rebaseline/README.md`.

### 3. Hard cutover is permitted

Because no production user depends on current compatibility, target architecture may intentionally break/delete current routes, schemas, IDs, package APIs, module boundaries and frontend redirects.

Compatibility requires a measured consumer/reason.

Hard cutover does not authorize an indefinitely red/ambiguous `main`; each landing must converge to one authority with proof.

### 4. Prior ADR records remain historical but their target authority is classified below

A prior ADR marked **reopened** by this decision remains useful evidence. It does **not** constrain target design until the named D-stage re-adjudicates it. If the old decision survives, the stage records why and may restore it as current. If it does not, a later ADR amends/supersedes it explicitly.

This later ADR takes precedence over the older status for target-design work.

## Still-binding constraints during the rebaseline

These remain current unless a new material finding explicitly reopens them:

- **ADR-005** — Mercado Livre is the first operational control plane.
- **ADR-006** — Oracle/Sankhya reads are MPC-owned behind application adapter boundaries.
- **ADR-007** — godror/OCI is the current canonical Oracle runtime.
- **ADR-009** — fee values carry provenance.
- **ADR-013** — webhook payload is a pointer/trigger, not trusted domain truth.
- **ADR-021** — TanStack Query is the frontend server-state mechanism; D6 may redesign package/route topology without duplicating server state authority.
- **ADR-025** — raw provider PII is not retained merely for convenience.
- **ADR-027** — absence from a partial pull is not closure/deletion.
- **ADR-029** — provider writes are not blindly retried after failure/ambiguous outcome.
- **ADR-033** — external marketplaces enter through vendor adapter trees implementing consumer-owned ports.
- **ADR-034** — `internal/kernel/fact` is the accepted primitive replacing the old prose-only unknown-is-never-zero mechanism; **D2 decides where uncertainty semantics actually require it**, not whether every value must use it.

These are constraints, not a complete target architecture.

## Reopened / non-authoritative for target design

The following are reopened because they encode structural choices D0–D9 must evaluate as a coherent system:

| ADR | Previous decision | Re-adjudicated in |
|---|---|---|
| 003 | integration spec split/sequencing | D4 / D9 |
| 004 | integration catalog plugin framework | D1 / D4 |
| 008 | production deploy topology | D7 |
| 010 | ML polling / visible refresh policy | D4 / D7 |
| 011 | divergence table semantics | D1 / D2 / D3 |
| 012 | DIFAL single source inside old `pricing` | D1 / D2 |
| 014 | market collection on-demand | D1 / D4 |
| 015 | listings module read-only boundary | D1 / D4 |
| 016 | hand-written SDK + same-commit atomicity | D5 |
| 018 | mutation envelope + in-process poller | D1 / D3 / D7 |
| 019 | listings ingestion feeding old snapshot observer | D1 / D3 |
| 020 | market data only through old `CollectorPort` | D1 / D4 |
| 022 | provider write requires `SELLER_SKU == CODPROD` | D1 / D2 / D4 |
| 023 | old module protocol | D1 |
| 024 | `IngestOrder` exact single write path | D1 / D3 |
| 026 | scheduler phase vocabulary | D3 / D7 |
| 028 | auto-link only on old concordant anchors | D1 / D2 |
| 030 | second scheduler instance per installation | D7 |
| 031 | products-mirror keep-absent merge | D1 / D2 |
| 032 | ML catalog-offers flag/default policy | D4 |

Reopened means **do not cite the old decision as proof that the target must keep it**.

## Already superseded historical records

- ADR-001 and ADR-002 were superseded in July 2026.
- ADR-017 was superseded by ADR-034.

Their files may remain for numbering/provenance but are not current authority.

## Consequences

1. `AGENTS.md` routes every new session through the rebaseline status before any old architectural artifact.
2. `ARCHITECTURE.md` contains only stable product-level constraints during the program; old module layout is removed from it.
3. `docs/architecture/decisions/README.md` is the only ADR status registry; its status column reflects this decision.
4. Legacy planning/documentation surfaces are deleted from the active tree rather than moved to another archive directory.
5. The current OpenAPI, schema, modules and frontend are D0 evidence. Their existence does not automatically grant them target authority.
6. D1–D9 can restore an old decision if it survives the global-maximum review; that survival must be justified by domain/failure-mode evidence, not inertia.
7. No implementation-plan generation or Codex execution of the rebaseline occurs before D9 acceptance.

## Alternatives considered

### Keep all old docs/ADRs active and add a new roadmap

Rejected. It preserves several mutually valid-looking routes through the repository and forces every fresh session to infer precedence.

### Delete all ADRs and start numbering over

Rejected. It destroys provenance, breaks citations and confuses historical evidence. The problem is authority status, not the existence of a record.

### Move legacy docs to `docs/archive/` or `old/`

Rejected. It creates another searchable in-repo knowledge surface and invites future agents to resurrect superseded designs. Git already provides immutable history.

### Begin implementation and decide technical details per context PR

Rejected. That is the root failure mode the rebaseline exists to remove: local PRs would decide identity, data ownership, events, retries, API and frontend semantics independently and recreate architectural drift under cleaner folder names.

## Proof / review

This decision is correctly applied when a fresh session can read `AGENTS.md`, `docs/engineering/rebaseline/README.md`, the current D-stage and the ADR registry and determine unambiguously:

- what is accepted;
- what is reopened;
- what work is prohibited;
- current stage/progress;
- exact next action;
- why historical documents cannot silently regain authority.