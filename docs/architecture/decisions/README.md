# Architecture Decision Records

## Current authority

**Read ADR-035 before treating any earlier structural ADR as target authority.**

Marketplace Central is in the D0–D9 Architecture Rebaseline. Earlier ADR files remain as historical records, but ADR-035 explicitly reopens decisions that the current program must re-adjudicate.

Current status vocabulary:

- **binding** — current constraint during the rebaseline;
- **reopened** — historical evidence only until the named D-stage decides it;
- **superseded** — historical only.

A direct old ADR file may still show the status it had when written. This registry plus later ADRs determine current status.

## Registry

| ADR | Current status |
|---|---|
| 001 | superseded |
| 002 | superseded |
| 003 | reopened — D4/D9 |
| 004 | reopened — D1/D4 |
| 005 | binding — Mercado Livre first |
| 006 | binding — MPC-owned Oracle reads |
| 007 | binding — godror/OCI Oracle runtime |
| 008 | reopened — D7 |
| 009 | binding — fee provenance |
| 010 | reopened — D4/D7 |
| 011 | reopened — D1/D2/D3 |
| 012 | reopened — D1/D2 |
| 013 | binding — webhook payload is not domain truth |
| 014 | reopened — D1/D4 |
| 015 | reopened — D1/D4 |
| 016 | reopened — D5 |
| 017 | superseded by ADR-034 |
| 018 | reopened — D1/D3/D7 |
| 019 | reopened — D1/D3 |
| 020 | reopened — D1/D4 |
| 021 | binding — TanStack Query owns frontend server state |
| 022 | reopened — D1/D2/D4 |
| 023 | reopened — D1 |
| 024 | reopened — D1/D3 |
| 025 | binding — provider PII retention minimization |
| 026 | reopened — D3/D7 |
| 027 | binding — partial-pull absence is not closure |
| 028 | reopened — D1/D2 |
| 029 | binding — no blind retry of provider writes |
| 030 | reopened — D7 |
| 031 | reopened — D1/D2 |
| 032 | reopened — D4 |
| 033 | binding — vendor adapters implement consumer-owned ports |
| 034 | binding primitive — D2 decides application scope |
| 035 | binding — Architecture Rebaseline governs target design |

## Numbering and provenance

ADR numbers use three digits and are never reused. The existing ADR files and `_citations/` remain for provenance; they are not a roadmap.

A reopened ADR must not be cited as proof that the target architecture should keep its old structure. The responsible D-stage may restore, amend or supersede it after current analysis.

Current program stage/status lives only in `docs/engineering/rebaseline/README.md`.