# ADR-014 — what the citations assert
**Harvested:** 2026-08-05 · **Total citations (excl. scripts/.runs):** ~40
**Spellings found:** ADR-14 only (no `ADR-014`, no `ADR 14`, no lowercase hits found).

**Same cross-mission collision pattern as ADR-13.** MIS-003 and MIS-007 each mint an unrelated "ADR-14."

## Assertion A1 — Market/competitor data must come only from a `CollectorPort` behind a contract-only module; no scraping, no fabricated facts, honest-empty when a signal is unavailable (MIS-003 usage)
- Citations: 3
- Verbatim: "Market contract-only behind CollectorPort | decided | fabricated market facts; scraping drift | 6-signal separation; honest-empty; no production adapter; G1–G7 sequence for successor"
- Anchors:
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/mission.md:113`
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-06-corrigir-atributo-market-contracts/milestone.md:17` ("R-04 (G1–G7, scraping forbidden), ADR-14 (market contract-only) / ADR-12 (`listings` module hosts the category-attributes read)")
  - `.mnfs/MIS-003-operator-cockpit-wireframe-replan/M-01-listings-read-spine/F-01-listings-module-ingestion/evidence/planner-sol-medium.log:1455` (log restatement of the same milestone.md line)

## Assertion A2 — `root.go` wiring and the OpenAPI+SDK contract pair are hub-arbitrated seams: each milestone enters via its own composition constructor in an anchored region, and at most one FE-contract commit may be "in flight" per lane at a time (code parallelizes, the contract commit does not) (MIS-007 usage)
- Citations: ~35 (mission.md, every MIS-007 milestone.md, all `research/*-interface-contract.md` files, several `planning-reviews/*.md`)
- Verbatim: "commits de contrato FE serializados pelo hub — ≤1 COMMIT de contrato em voo por vez; código paraleliza" (as amended; see Contradictions)
- Anchors spanning subsystems:
  - `.mnfs/MIS-007-ml-sync/mission.md:248,286,337,340,403` (root.go ownership + amended commit-serialization rule + R-7 risk row)
  - `.mnfs/MIS-007-ml-sync/research/webhook-inbox-interface-contract.md:152`
  - `.mnfs/MIS-007-ml-sync/research/sync-ingest-ports-interface-contract.md:145`
  - `.mnfs/MIS-007-ml-sync/research/sync-health-interface-contract.md:177,179`
  - `.mnfs/MIS-007-ml-sync/research/orders-persistence-interface-contract.md:174`
  - `.mnfs/MIS-007-ml-sync/research/channel-fees-interface-contract.md:175`
  - `.mnfs/MIS-007-ml-sync/M-04-listings-backfill-ingest/milestone.md:31,52` and `F-04-scheduler-refresh-wiring/feature.md:39`
  - `.mnfs/MIS-007-ml-sync/M-05-listings-fees-divergence/milestone.md:64` and `validation-contract.md:124`
  - `.mnfs/MIS-007-ml-sync/M-07-pricing-fee-read/milestone.md:61`
  - `.mnfs/MIS-007-ml-sync/M-09-sync-observability/milestone.md:60`
  - `.mnfs/MIS-007-ml-sync/planning-reviews/p5-claude-decomposition-audit-r04.md:55,57` and `p5-reconciliation-r03.md:47-49,70,84`

## Contradictions
- **Cross-mission number collision** (same class as ADR-13): MIS-003's ADR-14 (market data honesty) and MIS-007's ADR-14 (root.go/contract-commit serialization) are unrelated decisions sharing a number.
- **Self-detected and repaired drift inside Assertion A2.** The original ADR-14 wording ratified "≤1 **milestone** with an FE contract in flight," which directly contradicted the mission's own ratified Lane C schedule running three FE-contract milestones (M-05/M-06/M-07) in parallel — flagged as **blocking finding P-2** in `.mnfs/MIS-007-ml-sync/planning-reviews/p5-claude-decomposition-audit-r03.md:115,155` ("ADR-14 ratifies ≤1 FE-contract milestone in flight *overriding code disjointness*, and the Parallel Execution Plan schedules three, tripping the mission's own R-7 trigger by construction"). It was fixed, not by changing the lane plan, but by amending ADR-14 itself to the weaker rule "≤1 **commit** of contract in flight, hub arbitrates order, code parallelizes" — recorded in `planning-reviews/p5-reconciliation-r03.md:47-49` and `p5-claude-decomposition-audit-r04.md:55-58` (P-2, "CLOSED, verified exhaustively"). Every citation of ADR-14 in MIS-007 post-dating that fold cites the amended (weaker, per-commit) form; pre-fold artifacts are the ones documenting the contradiction itself, not surviving instances of it.

## Exceptions / carve-outs
- `.mnfs/MIS-007-ml-sync/M-04-listings-backfill-ingest/milestone.md:31`: "ADR-14: lane B não carrega contrato FE" — Lane B (M-03/M-04) is explicitly exempted from the FE-contract serialization rule because it carries no FE contract at all.
- `.mnfs/MIS-007-ml-sync/M-04-listings-backfill-ingest/_chip-m04/EVIDENCE.md:333`: a second `root.go` edit locus flagged against the "one anchored region per milestone" reading of ADR-14 was **downgraded** on review — ruled a necessary, minimal companion edit (a forced struct-field rename), not a violation of the anchored-region rule.
