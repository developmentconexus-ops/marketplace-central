# Architecture Decision Records

> **Role:** sole active authority for ADR file status/disposition during the D0–D9 rebaseline.  
> **Decision generation routing:** `docs/engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md`  
> **Current program status:** `docs/roadmap.md`

## Current posture

Pre-rebaseline ADRs are architecture history/evidence, not the new target architecture by inheritance. The 2026-08-18 Decision Reconciliation removed every fully rehomed legacy ADR from the active tree.

Detailed target meaning lives in accepted D-stage artifacts and `ARCHITECTURE.md`. A retained legacy ADR remains only because a later stage still has a concrete unresolved mechanism/proof/transition question for which that file is useful evidence.

Git history is the archive. ADR numbers are never reused.

## Retained legacy residues

| ADR | Active status | Why it remains | Retirement condition |
|---|---|---|---|
| 008 | reopened — D7 | deployment/publisher/host evidence for production topology; its old Oracle/Tailscale-to-Sankhya path is not current target meaning | retire when D7 adjudicates deployment topology |
| 010 | reopened — D7 | acquisition cadence/polling/runtime evidence; semantic freshness obligations already live in D0/D3 | retire when D7 adjudicates acquisition runtime |
| 017 | historical evidence retained by explicit D2 Fact rehoming gate | reconstructed domain-judgment clauses still needed until the replacement target Fact ADR rehomes what remains valid | retire together with 034 when target Fact ADR lands |
| 018 | reopened — D7 | execution-safety/runtime mechanics residue; generic Mutation business owner/table/poller is already superseded by D1/D3 | retire when D7 adjudicates execution-safety runtime |
| 026 | reopened — D7 | scheduler/cursor/phase runtime evidence; no global phase business vocabulary survives D3 | retire when D7 adjudicates scheduler/runtime mechanics |
| 030 | reopened — D7 | worker/scheduler/installation-topology evidence | retire when D7 adjudicates process/scheduler topology |
| 034 | carried historical evidence anchor | historical `Fact<T>` proof and domain-judgment clauses remain useful to D2; the former implementation/package is retired and owns no target structure | retire with 017 when target Fact ADR rehomes the remaining clauses |
| 035 | carried transition authority | governs the D0–D9 target-design/authority transition and implementation block | retire only after D0–D9 closes |

### ADR-035 snapshot fence

ADR-035's embedded “still-binding constraints” and “reopened” tables are a **2026-08-14 historical snapshot**. They do not own current ADR disposition.

This registry and accepted D-stage artifacts supersede those embedded tables wherever later adjudication differs. In particular, D4-B1 superseded Direct Oracle/godror target transport; ADR-006/007 are retired from the active tree and are not current target constraints.

### Repository Standard routing amendment — 2026-08-20

ADR-035 remains the D0–D9 transition authority. Pre-standard references inside accepted/historical artifacts to former current-status/read-order routers are frozen historical routing prose and are superseded for navigation only by:

```text
AGENTS.md
→ docs/index.md
→ docs/roadmap.md
```

`docs/roadmap.md` alone owns mutable current-program status/allowed work/next action; `docs/index.md` alone owns task routing. This changes routing only and does not reopen ADR-035's implementation block or any accepted D-stage semantics.

## Retired pre-rebaseline ADRs

The following pre-rebaseline ADRs were fully adjudicated/rehomed and retired from the active tree on 2026-08-18:

`001, 002, 003, 004, 005, 006, 007, 009, 011, 012, 013, 014, 015, 016, 019, 020, 021, 022, 023, 024, 025, 027, 028, 029, 031, 032, 033`.

Their still-valid semantic meaning, if any, now lives in accepted D-stage authority and/or `ARCHITECTURE.md`; their original files remain available only through Git history.

Notable reconciliations include:

- ADR-003's OAuth → fee-sync → UX sequence has no remaining D9 consumer and is historical;
- ADR-005 Mercado Livre-first meaning is carried by D0/`ARCHITECTURE.md`;
- ADR-006/007 Direct Oracle/godror target meaning is superseded by D4-B1 Gateway-only integration;
- ADR-009 provenance meaning is carried by D2/D4-B4;
- ADR-013 webhook-pointer meaning is carried by D3/D4;
- ADR-016 manual OpenAPI+SDK authority is superseded by D5-B1's single OpenAPI wire authority;
- ADR-021 TanStack server-state constraint is carried by `ARCHITECTURE.md` and remains subject to D6 only within that fence;
- ADR-025 PII minimization, ADR-027 honest partial absence, ADR-029 no-blind-retry and ADR-033 consumer-owned ports are carried by active architecture/D-stage authority.

## Citation archaeology

`docs/architecture/decisions/_citations/` contains only citation evidence still directly referenced by retained legacy residues. It is evidence, never target authority.

The retained citation files are:

- `RENUMBERING-REGISTRY.md` — provenance needed by retained reconstructed ADRs;
- `adr-009-citations.md` — referenced by retained ADR-010;
- `adr-013-citations.md` — referenced by retained ADR-018;
- `adr-07-twodigit-citations.md` — referenced by retained ADR-026;
- `adr-08-twodigit-citations.md` — referenced by retained ADR-030;
- `adr-017-citations.md` — referenced by retained ADR-017/034.

Everything else was retired with its citing ADRs and remains in Git history.

## New target ADRs

New target ADRs continue from **ADR-036+**. Historical numbers are never reused.

Only accepted target decisions that materially benefit from durable ADR treatment become new ADRs; accepted D-stage artifacts are not mechanically exploded into one ADR per bullet.

## Authority rule

- `docs/roadmap.md` owns mutable current program status, allowed/blocked work and exact next action.
- `docs/index.md` owns selective task routing.
- `ARCHITECTURE.md` owns stable cross-stage constraints.
- Decision Reconciliation Baseline owns current decision-generation routing.
- **This registry alone owns ADR file status/disposition.**
- Accepted D-stage artifacts own detailed semantic meaning.

If another active document appears to give an ADR a conflicting status, this registry is the status authority; resolve the semantic conflict against the responsible accepted D-stage/architecture home rather than keeping two status catalogs.
