# Architecture Decision Records

> **Role:** active ADR status/disposition registry during the D0–D9 rebaseline.  
> **Current program status:** `docs/roadmap.md`  
> **Target architecture:** accepted D-stage owners + `ARCHITECTURE.md`

## Current posture

Pre-rebaseline ADRs are historical evidence, not target architecture by inheritance. Still-valid target meaning must live in accepted D-stage authority or `ARCHITECTURE.md`; Git history is the archive and ADR numbers are never reused.

An old ADR remains in the active tree only while it has a named unresolved current consumer/retirement condition.

## Retained legacy residues

| ADR | Current status | Why it remains | Retirement condition |
| --- | --- | --- | --- |
| [017](017-unknown-is-never-zero.md) | **KEEP CURRENT EVIDENCE** | remaining Fact/domain-judgment evidence still protected by the D2 rehome gate | retire together with 034 after the remaining Fact clauses are explicitly rehomed |
| [034](034-fact-substitui-adr-017.md) | **KEEP CURRENT EVIDENCE** | retained `Fact<T>` replacement/evidence anchor; old implementation/package owns no target structure | retire with 017 after the remaining Fact clauses are rehomed |
| [035](035-architecture-rebaseline-governs-target-design.md) | **KEEP TRANSITION AUTHORITY** | governs the D0–D9 target-design transition and implementation block | retire only after D0–D9 closes |

ADR-035's embedded old status/reopen tables are historical snapshots. This registry plus current accepted D-stage owners supersede those status snapshots where later adjudication differs.

## Retired legacy ADRs

Previously retired on 2026-08-18:

`001, 002, 003, 004, 005, 006, 007, 009, 011, 012, 013, 014, 015, 016, 019, 020, 021, 022, 023, 024, 025, 027, 028, 029, 031, 032, 033`.

Repository-health retirement after accepted D7 consolidation:

`008, 010, 018, 026, 030`.

Their surviving target meaning is now owned by current D4/D7/D3 authority as applicable; the original ADR text remains available through Git history.

Notable current dispositions:

- deployment/host/backup/runtime topology → D7/D7-E;
- Mercado Livre acquisition cadence/recovery → D4 + D7;
- generic Mutation envelope/poller business shape → rejected; owner intents + D3/D7 external-effect safety are current;
- global scheduler phase vocabulary → not current Product/runtime authority;
- scheduler/worker-per-Installation topology → D7 runtime mechanism, not legacy ADR authority.

## Citation archaeology

`docs/architecture/decisions/_citations/` retains only evidence still consumed by retained ADRs.

Current retained files:

- [`RENUMBERING-REGISTRY.md`](_citations/RENUMBERING-REGISTRY.md) — provenance for retained reconstructed ADRs;
- [`adr-017-citations.md`](_citations/adr-017-citations.md) — consumed by ADR-017/034.

Citation files whose last active ADR consumer was retired are removed from the active tree and remain in Git history.

## New target ADRs

New target ADRs continue at **ADR-036+**. Historical numbers are never reused.

Create a new ADR only when an accepted target decision materially benefits from durable ADR treatment. Do not mechanically explode accepted D-stage authority into one ADR per decision bullet.

## Authority rule

- `docs/roadmap.md` owns mutable current status/allowed work/next action.
- `docs/index.md` owns selective task routing.
- `ARCHITECTURE.md` owns stable cross-stage constraints.
- accepted D-stage artifacts own current detailed target semantics.
- **this registry owns only ADR file status/disposition.**

If another active document gives a conflicting legacy ADR status, resolve it against this registry and the responsible current D-stage owner rather than preserving two status catalogs.
