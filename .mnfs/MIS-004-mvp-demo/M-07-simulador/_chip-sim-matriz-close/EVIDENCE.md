# CHIP-SIM-MATRIZ — Hub Close Evidence Pack

```yaml
id: CHIP-SIM-MATRIZ
type: hub-close-evidence
author: Dispatch Hub (Wave C)
mission: MIS-004-mvp-demo
milestone-owner: M-07-simulador
contract: MIS-004-C10 (design-fidelity sweep) + SCREEN-INVENTORY gap #1
created: 2026-07-19
```

- **Branch:** chip/sim-matriz
- **Tip:** 411e00b0 (LIVE-VERIFIED fill) · hardened head aeee4b53
- **Base:** 783cbc0d
- **Full chip pack (brought in by this merge):** `.mnfs/MIS-004-mvp-demo/_chip-sim-matriz/EVIDENCE.md`
- **Scope:** `apps/web/src/pages/precos/**` ONLY — FE render-only, OpenAPI/SDK FROZEN, zero backend/migration.

## Merge-gate markers (harness 0.4.0)

P6-DUAL-GATE: AGREEMENT

Dual gate ran during the chip (2 rounds). Round 1: cold Opus PASS, adversarial sonnet FAIL (ADR-17 loading/error states collapsed into confirmed SEM_EVIDENCIA; INSUFFICIENT_MARKET untested) → no agreement, fixed at 02ed4a14 + aeee4b53. Round 2 (hardened head aeee4b53): cold Opus PASS + adversarial sonnet PASS → **AGREEMENT**. Hub verified the marker + diff scope (confined to `pages/precos/**`) at close.

LIVE-VERIFIED: 2026-07-19 hub P7 live-drive on the clean docker dev-stack (frontend :5174 + backend :8080 on the SIM worktree mount). `/precos` renders the product matrix as the main surface with the exact 7 design columns (`SKU · DESCRIÇÃO · CUSTO · NOSSO PREÇO · PREÇO MERCADO · MARGEM · VEREDICTO`), 50 rows. Honest ADR-17 live: real CUSTO off the fact (R$ 42.1 / id 412), unknown market/preço/margem = "—", VEREDICTO SEM_EVIDENCIA in neutral ink `rgb(37,41,31)` (never a fabricated verde/red), "novo" tag on zero-listing products, zero fabricated "R$ 0" margin. Row selection opens the 380px aside reusing the existing DecompositionPanel + SolverPanel + MarketComparison (not rebuilt). Theme paper `rgb(251,250,247)` + Instrument Sans. Zero console errors.

## Demo-data FINDING (hub-verified; NOT a chip defect) — action for operator

Every ERP catalog product returns `NO_PRICE_EVIDENCE` from `GET /market/aggregates?codprod=<internal_product_id>` (hub verified: ids 412–511 all `NO_PRICE_EVIDENCE`). Live market evidence exists only under the **ML listing codprod space** (e.g. 90008 Papeleira, 16 offers) — disjoint from the ERP `internal_product_id` space the matrix keys on. Consequence: on `/precos` the PREÇO MERCADO / MARGEM / VEREDICTO columns render all-`—`/SEM_EVIDENCIA — honest and correct for this fixture, but visually empty for a *pricing* demo. Same data-gating M-06 product 412 hit. The matrix faithfully renders what the endpoint returns; the populated OK/margin path is golden-tested. **Operator decision for T-0:** either provision market aggregates keyed to a few ERP `internal_product_id`s (ties to `demo-data-provisioned` + D-95 EAN linkage), or demo market intel via `/produto/:id` and `/anuncios` (which already surface it). Non-blocking for the SIM merge.

## Verdict

**CHIP-SIM-MATRIZ: PASS.** P6 dual-gate AGREEMENT + P7 design-fidelity live PASS. Cleared for `--no-ff` merge into main. Demo-data market-keyspace finding carried to the mission-close sweep for operator ruling.
