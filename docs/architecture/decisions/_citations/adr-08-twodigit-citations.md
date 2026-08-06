# ADR-08 (two-digit citations) — what they assert
**Harvested:** 2026-08-05 · **Two-digit citations:** 34 (live code: 2, .mnfs only: 32, of which 2 are this repo's own prior ADR-registry notes) · **Three-digit `ADR-008` citations:** 17
**Existing document `008-production-deploy-topology.md`:** "Image-based delivery with a VPS-primary production topology" — DOES NOT MATCH: the document is about deploy infrastructure (VPS vs. serverless/PaaS, image-based delivery); neither two-digit assertion below (scheduler pattern, zero-writes-ML gate) is about deploy topology.

## Assertion A1 — "Second Scheduler instance": listings backfill/sweep and orders backfill/incremental each get their own fixed-cadence scheduler instance (orders 5min, listings daily), both built from the same `synccomposition.NewProductsScheduler` constructor pattern, wired with one anchored line in `root.go` (mission: MIS-007-ml-sync) — this is the "segunda instância do Scheduler" the task brief asked to verify, and it is CONFIRMED live in code
- Citations: 11 (live code: 2)
- Verbatim: "the ADR-08 \"segunda instância\" of the sync scheduler pattern sync/composition.NewProductsScheduler established"
- Anchors:
  - `apps/server_core/internal/composition/root.go:837`
  - `apps/server_core/internal/modules/listings/composition/scheduler.go:2`
  - `.mnfs/MIS-007-ml-sync/mission.md:201`
  - `.mnfs/MIS-007-ml-sync/M-06-orders-backfill-decomposition/milestone.md:25`
  - `.mnfs/MIS-007-ml-sync/M-06-orders-backfill-decomposition/F-01-backfill-incremental/feature.md:32`
  - `.mnfs/MIS-007-ml-sync/M-04-listings-backfill-ingest/F-04-scheduler-refresh-wiring/feature.md:26`
  - `.mnfs/MIS-007-ml-sync/M-04-listings-backfill-ingest/F-04-scheduler-refresh-wiring/feature.md:39`
  - `.mnfs/MIS-007-ml-sync/research/sync-health-interface-contract.md:114`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p5-claude-decomposition-audit-r07.md:119`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p5-claude-decomposition-audit-r05.md:195`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p5-claude-decomposition-audit-r04.md:88`

## Assertion A2 — Zero writes to Mercado Livre in the demo: every PROVIDER-target mutation goes through the M-03 envelope queue (preview + protocol), with the provider dispatcher held OFF by config; LOCAL-state writes (e.g. the M-04 product-links batch) are an explicitly named exception (P5-F-12) applied inside their own owning module with preview+audit, outside the envelope (mission: MIS-004-mvp-demo)
- Citations: 18 (live code: 0)
- Verbatim: "Zero writes ML: toda mutação com alvo PROVIDER via fila M-03 preview+protocolo (dispatcher provider OFF por config na demo); estado local ... fora do envelope — exceção P5-F-12"
- Anchors:
  - `.mnfs/MIS-004-mvp-demo/mission.md:53`
  - `.mnfs/MIS-004-mvp-demo/mission.md:85`
  - `.mnfs/MIS-004-mvp-demo/mission.md:92`
  - `.mnfs/MIS-004-mvp-demo/mission.md:98`
  - `.mnfs/MIS-004-mvp-demo/mission.md:185`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-sol-decomposition-audit-r02.md:97`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-sol-decomposition-audit-r02.md:152`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-sol-decomposition-audit-r03.md:24`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-sol-decomposition-audit-r04.md:25`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-reconciliation-r02.md:26`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p5-reconciliation-r03.md:27`
  - `.mnfs/MIS-004-mvp-demo/planning-reviews/p7-claude-readiness-r01.md:46`

## Contradictions
- Two unrelated decisions share the token `ADR-08` (MIS-007 scheduler-instance pattern vs. MIS-004 zero-writes-ML gate). Neither matches the three-digit `008-production-deploy-topology.md` document.
- Internal MIS-007 planning churn: the P3 candidate `ADR-08` meant "MASS-CLOSURE dies, absence ≠ closed, marked only after a complete run" (`.mnfs/MIS-007-ml-sync/planning-reviews/p3-claude-candidate-r01.md:82`, reconciled toward shape `A-06` in `p3-reconciliation-r01.md:21`) — that rule survived but was **renumbered to `ADR-06`** in the ratified `mission.md` (see `adr-06-twodigit-citations.md` Assertion A1), not `ADR-08`.
- Internal MIS-004 planning churn: `p3-sol-counterproposal-r01.md:73` used candidate `ADR-08` for "Retheme-first frontend convergence" (the FE shell decision), while `p3-claude-candidate-r01.md` called the same idea `ADR-07`; `p3-reconciliation-r01.md:15` shows the two candidates' numbers colliding before reconciliation. Separately, the "zero writes ML" assertion itself entered reconciliation as candidate `ADR-08` and was reassigned to ratified `ADR-09` in one candidate's numbering before settling back to `ADR-08` in the final `mission.md` — this specific churn is already recorded in `docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md:80` and `adr-009-citations.md:40`.

## Amendments
None found that amend the ratified text of either A1 or A2 after mission ratification; MIS-004's `mission.md:98` lists ADR-08 among "Accepted trade-offs" ("aplicar não conclui de verdade — protocolo na fila"), which is a scoping caveat rather than a textual amendment.
