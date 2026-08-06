# ADR-01 (two-digit citations) — what they assert
**Harvested:** 2026-08-05 · **Two-digit citations:** 29 (live code: 0, .mnfs only: 25; docs/other: 4)
**Three-digit `ADR-001` citations:** 3 (1 is the document's own title line `001-metalshopping-direct-read.md:1`; 1 real external citation in `docs/superpowers/plans/2026-08-05-arquitetura-protocolo-de-modulo-plan.md`)
**Existing document `001-metalshopping-direct-read.md`:** "MPC reads products directly from MetalShopping Postgres" (superseded 2026-07-07) — DOES NOT MATCH: none of the four live two-digit assertions (dispatch precondition, core/adapter boundary, mirror-as-current-state, paginated Oracle ports) has anything to do with the retired MetalShopping-Postgres direct-read topology.

## Assertion A1 — Wave execution must start from main post-merge-W1, with a 40-hex base SHA recorded per chip (mission: MIS-004)
This is a dispatch/process precondition, not a system-architecture rule.
- Citations: 9 (live code: 0)
- Verbatim: "ADR-01 Base = main pós-merge W1 (≥ f4612be3); base SHA 40-hex registrado por chip"
- Anchors:
  - `.mnfs/MIS-004-mvp-demo/mission.md:85`
  - `.mnfs/MIS-004-mvp-demo/research/w1-merge-addendum-2026-07-17.md:25`
  - `.mnfs/MIS-004-mvp-demo/research/w1-merge-addendum-2026-07-17.md:47`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:9`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-reconciliation-r01.md:65`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-claude-candidate-r01.md:12`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-sol-decomposition-audit-r03.md:24`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-sol-decomposition-audit-r04.md:25`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-reconciliation-r03.md:27`

## Assertion A2 — Core module is provider-agnostic; the Mercado Livre adapter is the tested boundary ("núcleo nativo × adapter") (mission: MIS-007)
- Citations: 8 (live code: 0; 3 of these are spine range-enumerations `ADR-01..ADR-14` rather than restatements of the content)
- Verbatim: "ADR-01 Núcleo nativo × adapter (design §3). Núcleo agnóstico de provider; adapter ML é..."
- Anchors:
  - `.mnfs/MIS-007-ml-sync/mission.md:144`
  - `.mnfs/MIS-007-ml-sync/M-02-sync-core-seam/F-02-fee-divergence-ports/feature.md:25`
  - `.mnfs/MIS-007-ml-sync/M-01-ml-client-hardening/milestone.md:17`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p3-claude-candidate-r01.md:9`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p3-reconciliation-r01.md:14`

## Assertion A3 — `products_mirror` is the materialized current-state; snapshots become history, never the live read model (mission: MIS-006)
- Citations: 6 (live code: 0)
- Verbatim: "### ADR-01: products_mirror = estado corrente materializado; snapshots viram history"
- Anchors:
  - `.mnfs/MIS-006-integracao-fundacao/mission.md:104`
  - `.mnfs/MIS-006-integracao-fundacao/M-02-mirror-port-active-source/milestone.md:22`
  - `.mnfs/MIS-006-integracao-fundacao/M-03-xlsx-adapter/milestone.md:21`
  - `.mnfs/MIS-006-integracao-fundacao/M-03-xlsx-adapter/milestone.md:24`
  - `.mnfs/MIS-006-integracao-fundacao/M-06-telas-sdk/milestone.md:113`
  - `.mnfs/MIS-006-integracao-fundacao/_chip-anchors/chip.md:111`

## Assertion A4 — Use-case-shaped paginated/batch ports in `internal_read`: one Oracle query per page/batch, keyset cursor (mission: MIS-002)
- Citations: 2 (live code: 0 — this ADR is cited only in the mission record; no `ADR-01` comment found in the M-01/foundation code that implements it)
- Verbatim: "ADR-01 Use-case-shaped paginated/batch ports in `internal_read` (1 Oracle query per page/batch, keyset cursor)"
- Anchors:
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:111`
  - `.mnfs/MIS-002-oracle-read-rearchitecture/mission.md:142`

## Contradictions
- **Four unrelated decisions share one number, one per mission.** MIS-002, MIS-004, MIS-006, and MIS-007 each independently renumbered their own ratified-decisions table starting at ADR-01. None of the four subjects (dispatch precondition, core/adapter boundary, mirror data model, Oracle port shape) overlaps with another.
- **Intra-mission candidate churn (MIS-004).** Before reconciliation, Sol's counterproposal used `ADR-01` for a completely different subject: `.mnfs/MIS-004-mvp-demo/planning-reviews/p3-sol-counterproposal-r01.md:17` — "### ADR-01: Dual ERP source behind one canonical import boundary". This was superseded during reconciliation; the ratified `mission.md:85` kept the number for the "base = main" precondition instead, and the dual-ERP-source content was ratified as MIS-004's ADR-02 (see `adr-02-twodigit-citations.md`).
- **Intra-mission amendment (MIS-004).** Early planning briefly let an *unqualified* "M-03 envelope = único write path" rule ride under the `ADR-01` label (`p5-sol-decomposition-audit-r03.md:24`, `r04.md:25`, flagged as finding P5-R2-05). It was folded/qualified into the Must-preserve column of `mission.md:85` as "ADR-08"-scoped instead, per `p5-reconciliation-r03.md:27`: "FOLDED — ADR-01 coluna Must-preserve qualificada".

## Amendments
- MIS-004 mission.md:85's Must-preserve column was rewritten from an unqualified universal rule to "único write path com alvo PROVIDER (estado LOCAL = módulo dono, ADR-08)" — see Contradictions above.
