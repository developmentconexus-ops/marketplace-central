# ADR-035: Architecture Rebaseline governs target design during D0–D9

**Date:** 2026-08-13  
**Amended:** 2026-08-14, 2026-08-18  
**Status:** accepted

## Context

Marketplace Central contains several architectural eras at once: legacy modules, newer contexts/adapters, multiple persistence/API/frontend shapes and historical plans/specs/handoffs that described incompatible targets.

The operator has chosen to design the technical system deeply enough that implementation does not invent architectural semantics locally. At the same time, the operator explicitly rejected turning that design program into a multi-week exhaustive audit of the existing repository before architecture discussion begins.

Those concerns are compatible only if two things stay separate:

1. **documentary/governance authority cleanup**, which removes competing historical authority now; and
2. **target-system design**, which consults current code/schema/runtime as evidence on demand for the decision being made.

Current source shape is evidence about the present system. It does not receive target-design authority merely because it exists.

## Decision

### 1. Documentary authority cleanup precedes D0

Before target design begins, PR #41 (or an accepted successor) removes/retargets retired documentary authority and its active consumers.

This cleanup does **not** decide the target disposition of legacy product/runtime code.

It ends when current governance is self-contained, active consumers no longer route to/recreate retired documentary authority, verification is green without weakening controls, and a fresh session can identify one authority path and one exact next action.

After that point, stop cleaning and begin D0.

### 2. D0–D9 is the governing target-design program

The required sequence is:

1. **D0 — Product / System Definition**
2. **D1 — Domains / Boundaries**
3. **D2 — Identity / Tenant / Data Ownership**
4. **D3 — Communication / Events**
5. **D4 — External Integrations**
6. **D5 — API**
7. **D6 — Frontend**
8. **D7 — Runtime / Jobs / Transactions**
9. **D8 — Golden Flows**
10. **D9 — Adversarial Architecture Review**

Only after D9 acceptance may the repository create the implementation DAG/plan and begin product implementation.

### 3. Evidence is pulled by the decision; exhaustive legacy census is not a prerequisite

Each D-stage begins from its architectural/product question. It identifies the evidence needed to answer that question and then inspects only the relevant code, schema, contracts, runtime, external behavior and historical evidence.

Already-collected current-state facts live as supporting evidence (for example in `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`). They do not create a gate requiring every fact to be freshly reproduced before D0 or before design discussion.

When a stage needs a package graph, writer/reader census, route topology, adapter capability map, runtime reachability measurement or other repository analysis, that measurement is performed for that decision.

### 4. Git is the archive; active documentation has one current-program path

Historical plans, handoffs, wikis and old design documents are removed from the active tree after still-valid principles/evidence are absorbed where necessary.

There is no `old/`, `archive/` or parallel legacy roadmap.

Current program status and exact next action live only at `docs/engineering/rebaseline/README.md`.

### 5. Hard cutover is permitted, but only after the relevant design decision

Because no production user requires preservation of current application compatibility, accepted target architecture may intentionally replace/delete current routes, schemas, IDs, package APIs, module boundaries and frontend redirects.

Compatibility requires a measured consumer/reason.

This permission is **not** authority to delete legacy source during documentary cleanup. A source unit is adjudicated only when the D-stage responsible for its concern has established the target owner/contract and cutover.

Hard cutover also does not authorize an indefinitely red/ambiguous `main`; each landing must converge to one authority with proof.

### 6. Prior ADR records remain historical but target authority is classified

A prior ADR marked **reopened** remains useful evidence. It does not constrain target design until the relevant D-stage re-adjudicates it. If the old decision survives, the stage records why. If it does not, later accepted authority supersedes it explicitly.

## Still-binding constraints during the rebaseline — 2026-08-14 snapshot

The following table/list was the active transition snapshot on 2026-08-14. It is retained here only as historical evidence of the transition state at that date:

- **ADR-005** — Mercado Livre is the first operational control plane.
- **ADR-006** — Oracle/Sankhya reads are MPC-owned behind application adapter boundaries.
- **ADR-007** — godror/OCI is the current canonical Oracle runtime.
- **ADR-009** — fee values carry provenance.
- **ADR-013** — webhook payload is a pointer/trigger, not trusted domain truth.
- **ADR-021** — TanStack Query is the frontend server-state mechanism; D6 may redesign package/route topology without duplicating server-state authority.
- **ADR-025** — raw provider PII is not retained merely for convenience.
- **ADR-027** — absence from a partial pull is not closure/deletion.
- **ADR-029** — provider writes are not blindly retried after failure/ambiguous outcome.
- **ADR-033** — external marketplaces enter through vendor adapter boundaries implementing consumer-owned ports.
- **ADR-034** — `internal/kernel/fact` is the current uncertainty primitive; D2 decides its target scope.

These lines are **not current ADR disposition authority after 2026-08-18**. See the amendment below and the live ADR registry.

## Reopened / non-authoritative for target design — 2026-08-14 snapshot

| ADR | Previous decision | Re-adjudicated in (snapshot) |
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

This table is a **2026-08-14 snapshot**, not the current disposition map.

## Already superseded historical records — 2026-08-14 snapshot

- ADR-001 and ADR-002 were superseded in July 2026.
- ADR-017 was superseded by ADR-034.

## Consequences

1. `AGENTS.md` routes every fresh session through the current rebaseline status before historical material.
2. `docs/engineering/rebaseline/README.md` is the sole current program/status/next-action authority.
3. `ARCHITECTURE.md` contains stable constraints, not an inventory that blesses current code shape.
4. `docs/architecture/decisions/README.md` is the ADR status registry.
5. Existing OpenAPI, schema, modules, contexts, packages, runtime and frontend are supporting evidence consulted by the D-stage that needs them, not an exhaustive D0 prerequisite.
6. Legacy product/runtime disposition is decided stage by stage; documentary cleanup cannot smuggle in KEEP/REPLACE/DELETE decisions.
7. No product implementation-plan generation or implementation occurs before D9 acceptance.

## Alternatives considered

### Exhaustively census the repository before architecture discussion

Rejected. It turns current implementation history into the agenda for target design, delays product/system definition and creates a false prerequisite to inspect every legacy surface before asking what should exist.

### Delete/refactor obviously old code during documentary cleanup

Rejected. Age or directory shape is not a target-design decision. It can erase useful evidence and prematurely select one historical architecture over another.

### Keep all old docs/ADRs active and add a new roadmap

Rejected. It preserves several mutually valid-looking routes and forces fresh sessions to infer precedence.

### Delete all ADRs and start numbering over

Rejected. It destroys provenance and breaks historical citations. The problem is authority status, not existence of a record.

### Move legacy docs/source to `docs/archive/` or `old/`

Rejected. It creates another searchable in-repo authority surface. Git already provides immutable history.

### Begin implementation and decide architecture per feature PR

Rejected. Local implementation would again decide identity, data ownership, events, retries, API and frontend semantics independently.

## Proof / review

This decision remains correctly applied when a fresh session can determine unambiguously that:

- the D0–D9 program governs target design;
- code/schema/runtime are evidence pulled on demand by the relevant design decision;
- historical documents/ADRs cannot silently regain target authority;
- current status/next action is obtained from the router, not from this ADR's historical snapshots;
- product implementation remains blocked until D9 is accepted.

## Amendment 2026-08-18 — Decision Reconciliation Baseline

The accepted Decision Reconciliation Baseline completed the authority cleanup needed before D5-B2 becomes code-facing.

Binding clarification:

1. `docs/architecture/decisions/README.md` is the **sole current ADR file status/disposition authority**.
2. Accepted D-stage artifacts and `ARCHITECTURE.md` own current semantic meaning.
3. `docs/engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md` owns only decision-generation routing.
4. The “still-binding”, “reopened” and historical tables above are frozen **2026-08-14 transition snapshots** and do not override later adjudication.
5. D4-B1 superseded Direct Oracle/godror target transport; ADR-006/007 are not current target constraints and are retired from the active tree.
6. ADR-003's former D9 routing was adjudicated as vestigial; its old OAuth → fee-sync → UX sequencing is historical and retired from the active tree.
7. Git history is the archive for all retired ADR generations; numbers remain permanently reserved.

ADR-035's governing D0–D9 role is unchanged. This ADR remains active until the rebaseline program closes, at which point its transition role is itself retired according to the live ADR registry.
