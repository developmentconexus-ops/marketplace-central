# Architecture Decision Records

## Current rebaseline posture

Marketplace Central is rebuilding its target architecture through D0–D9. **Pre-rebaseline ADR files are legacy architecture history/evidence, not the new system's ADR baseline.**

During the rebaseline, target authority comes from the active authority path (`ARCHITECTURE.md` + accepted/current D-stage artifacts) and any explicitly created **new target ADRs**. A legacy ADR may remain in the active tree because:

- a later D-stage still needs it as evidence before adjudication; or
- a currently carried constraint still needs a safe rehome before the legacy file can be removed.

Do not infer target architecture from a legacy ADR merely because that file once said `accepted`/`binding`.

Status vocabulary in this registry:

- **carried constraint** — the meaning is actively carried by current rebaseline authority; the legacy ADR is evidence/provenance, not a second target authority;
- **reopened — D<n>** — legacy evidence awaiting the named stage's adjudication;
- **superseded by rebaseline** — the responsible stage has adjudicated the old shape; use current D-stage/`ARCHITECTURE.md` authority instead;
- **historical** — no current target role.

## Legacy ADR registry / transition state

| ADR | Current rebaseline status |
|---|---|
| 001 | historical |
| 002 | historical |
| 003 | reopened — D4/D9 |
| 004 | reopened — D4 (D1 portion adjudicated) |
| 005 | carried constraint — Mercado Livre first; active home: D0 / `ARCHITECTURE.md` |
| 006 | carried constraint — MPC-owned Oracle reads; active home: `ARCHITECTURE.md` |
| 007 | carried constraint — godror/OCI current Oracle runtime; active home: `ARCHITECTURE.md` |
| 008 | reopened — D7 |
| 009 | carried constraint — fee/value provenance; active home: D2 provenance semantics |
| 010 | reopened — D4/D7 |
| 011 | superseded by rebaseline — D1/D2 own divergence/work semantics; no generic divergence authority |
| 012 | superseded by rebaseline — D1/D2 own economics/provenance; legacy `pricing` DIFAL table is not target authority |
| 013 | carried constraint — webhook payload is not domain truth; active home: D0 external-evidence semantics |
| 014 | reopened — D4 (D1 portion adjudicated) |
| 015 | reopened — D4 (D1 portion adjudicated) |
| 016 | reopened — D5 |
| 017 | historical predicate retained as evidence until target Fact ADR rehomes still-valid domain-judgment clauses |
| 018 | reopened — D7 (D1/D3 semantic portions adjudicated; generic Mutation owner/table/poller not target authority) |
| 019 | historical — D1/D3 rehomed accepted-consumer, duplicate/idempotency and recoverable-propagation semantics |
| 020 | reopened — D4 (D1 portion adjudicated) |
| 021 | carried constraint — TanStack Query owns frontend server state; active home: `ARCHITECTURE.md` |
| 022 | superseded by rebaseline as identity law — D2 preserves pre-dispatch correspondence safety; D4 re-verifies concrete provider mapping |
| 023 | superseded by rebaseline — D1/`ARCHITECTURE.md` own semantic boundaries and private-implementation prohibition |
| 024 | historical — D1/D3 rehomed single semantic owner, trigger convergence and anti-regression semantics |
| 025 | carried constraint — provider PII minimization; active home: D0 / `ARCHITECTURE.md` |
| 026 | reopened — D7 (D3 semantic portion adjudicated; no global phase vocabulary carried forward) |
| 027 | carried constraint — partial-pull absence is not closure; active home: D0 / `ARCHITECTURE.md` |
| 028 | superseded by rebaseline — D1/D2 Readiness owns correspondence; D2 preserves corroboration/no-silent-human-override safety |
| 029 | carried constraint — no blind retry of ambiguous external writes; active home: D0 / `ARCHITECTURE.md` |
| 030 | reopened — D7 |
| 031 | superseded by rebaseline — no target Product mirror; honest absence remains current authority |
| 032 | reopened — D4 |
| 033 | carried constraint — vendor adapters implement consumer-owned ports; active home: `ARCHITECTURE.md` |
| 034 | carried implementation/evidence anchor — D2 decides target `Fact<T>` scope; replace with a new target Fact ADR before legacy cleanup |
| 035 | carried transition constraint — rebaseline governs target design; retain until D0–D9 program closes |

## Legacy-retirement gates

Legacy ADR deletion follows D2 §12:

1. a reopened ADR remains available until its responsible D-stage adjudicates the relevant meaning;
2. a carried constraint must have an active home in `ARCHITECTURE.md`, an accepted D-stage artifact or a new target ADR before its legacy file is removed;
3. ADR-017/034 still-valid domain-judgment clauses must be rehomed before those files are removed;
4. ADR-035 remains through rebaseline closure.

Git history is the archive; no active legacy ADR tree is retained merely for compatibility after these gates are satisfied.

## New target ADR numbering

New target ADRs continue from the next unused number, **ADR-036+**. Historical numbers are never reused, even after legacy files are deleted, so citations in Git history remain unambiguous.

Only target decisions that materially benefit from ADR treatment become ADRs; accepted D-stage artifacts are not mechanically split into one ADR per decision.

Current program stage/status lives only in `docs/engineering/rebaseline/README.md`.
