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
| 003 | reopened — D9 only; D4-B1 rehomed credential/identity-binding prerequisite and superseded the old OAuth→fee→frontend implementation sequence |
| 004 | superseded by rebaseline — D4-B1 rejects generic plugin/self-registration/auth-factory/fee-sync registry as target structure |
| 005 | carried constraint — Mercado Livre first; active home: D0 / `ARCHITECTURE.md` |
| 006 | reopened — D4-B3 as exception evidence only; D4-B1 superseded direct-Oracle as the normative/default Sankhya target transport |
| 007 | reopened — D4-B3/D7 only if a supported direct-DB exception survives; D4-B1 superseded godror/OCI as an inherited canonical default |
| 008 | reopened — D7 |
| 009 | carried constraint — fee/value provenance; active home: D2 provenance semantics |
| 010 | reopened — D7 only; D4-B1 superseded polling-only/no-webhook target meaning and rehomed notification→reread + coverage semantics |
| 011 | superseded by rebaseline — D1/D2 own divergence/work semantics; no generic divergence authority |
| 012 | superseded by rebaseline — D1/D2 own economics/provenance; legacy `pricing` DIFAL table is not target authority |
| 013 | carried constraint — webhook payload is not domain truth; active home: D0 / D3 / D4-B1 external-evidence semantics |
| 014 | reopened — D4-B4 (D1 portion adjudicated; on-demand/local-runtime shape has no authority by inheritance) |
| 015 | reopened — D4-B2 (D1 portion adjudicated; B1 rehomed source-qualified identity/coverage/reread, final ML listing contract remains open) |
| 016 | reopened — D5 |
| 017 | historical predicate retained as evidence until target Fact ADR rehomes still-valid domain-judgment clauses |
| 018 | reopened — D7 (D1/D3 semantic portions adjudicated; generic Mutation owner/table/poller not target authority) |
| 019 | historical — D1/D3 rehomed accepted-consumer, duplicate/idempotency and recoverable-propagation semantics |
| 020 | reopened — D4-B4 (D1 portion adjudicated; generic `CollectorPort` target shape not inherited) |
| 021 | carried constraint — TanStack Query owns frontend server state; active home: `ARCHITECTURE.md` |
| 022 | superseded by rebaseline as identity law — D2 preserves pre-dispatch correspondence safety; D4-B2 re-verifies concrete Mercado Livre mapping/identifier evidence |
| 023 | superseded by rebaseline — D1/`ARCHITECTURE.md` own semantic boundaries and private-implementation prohibition |
| 024 | historical — D1/D3 rehomed single semantic owner, trigger convergence and anti-regression semantics |
| 025 | carried constraint — provider PII minimization; active home: D0 / `ARCHITECTURE.md` |
| 026 | reopened — D7 (D3 semantic portion adjudicated; no global phase vocabulary carried forward) |
| 027 | carried constraint — partial-pull absence is not closure; active home: D0 / `ARCHITECTURE.md` / D4-B1 |
| 028 | superseded by rebaseline — D1/D2 Readiness owns correspondence/corroboration policy; D4-B2 supplies current provider identifier evidence |
| 029 | carried constraint — no blind retry of ambiguous external writes; active home: D0 / D3 / `ARCHITECTURE.md` |
| 030 | reopened — D7 |
| 031 | superseded by rebaseline — no target Product mirror; honest absence remains current authority |
| 032 | reopened — D4-B4; current catalog-offers env flag/default-off behavior is current-state evidence, not target authority |
| 033 | carried constraint — vendor adapters implement consumer-owned ports; active home: D1 / D4-B1 / `ARCHITECTURE.md` |
| 034 | carried implementation/evidence anchor — D2 decides target `Fact<T>` scope; replace with a new target Fact ADR before legacy cleanup |
| 035 | carried transition constraint — rebaseline governs target design; retain until D0–D9 program closes |

## D4-B1 transport transition note

D4-B1 explicitly re-adjudicated only the **normative Sankhya direct-transport meaning** behind ADR-006/007.

Current target authority is:

- Sankhya/business-system facts remain external to MPC;
- consumer-owned ports and adapter boundaries remain binding;
- provider-sanctioned Sankhya API Gateway is the default target transport for new MPC↔Sankhya contracts;
- a direct-database exception is not assumed or silently available;
- client-specific direct-DB contractual entitlement/compliance remains **Unknown** until proven;
- D4-B3 may admit a direct-DB exception only with explicit current provider/customer support/entitlement evidence plus proof that sanctioned APIs cannot meet required correctness, coverage and operational viability;
- if no exception survives B3, ADR-006/007 have no remaining target role and become historical evidence; if one survives, D7 adjudicates its runtime rather than inheriting ADR-007 automatically.

This transition does not reopen D0–D3.

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
