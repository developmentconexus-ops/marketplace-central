# CHIP-VINC-DESIGN — Hub Close Evidence Pack

```yaml
id: CHIP-VINC-DESIGN
type: hub-close-evidence
author: Dispatch Hub (Wave C)
mission: MIS-004-mvp-demo
milestone-owner: M-04-vinculos-import-ui
contract: MIS-004-C10 (design-fidelity sweep) + SCREEN-INVENTORY gap #2
created: 2026-07-19
```

- **Branch:** chip/vinc-design
- **Tip:** 868a9b01 (LIVE-VERIFIED fill) · feat @ f70bd54f
- **Base:** 783cbc0d
- **Full chip pack (brought in by this merge):** `.mnfs/MIS-004-mvp-demo/M-04-vinculos-import-ui/_chip-vinc-design/EVIDENCE.md`
- **Scope:** 13 files under `apps/web/src/pages/vinculos/**` ONLY (+593/−207) — FE render-only, zero backend/OpenAPI/SDK/migration.

## Merge-gate markers (harness 0.4.0)

P6-DUAL-GATE: AGREEMENT

Dual gate ran during the chip: cold Opus gate-reviewer PASS (8/8 contract items, file:line) + adversarial sonnet PASS (could not refute theme / GTIN-gating / band-thresholds / undo-audit-wiring / test-theater / scope) → AGREEMENT. Hub verified the marker + scope (all writes under `pages/vinculos/**`) at close.

LIVE-VERIFIED: 2026-07-19 hub P7 live-drive on the clean docker dev-stack (frontend :5174 + backend :8080 on the VINC worktree mount). `/vinculos` verified live:
- **P0 THEME** — off-theme scan over `main *` for literal tailwind numbered-palette classes = **zero offenders**; body paper `rgb(251,250,247)`; Confiança bands via `@theme` tokens (MEDIA = amber `rgb(138,109,31)`), never a literal palette color.
- **P1 REORIENT** — Fila is the 9-col anúncio-cêntrica table exactly (Selecionar · Anúncio ML · SKU ML · Produto sugerido · SKU HUB · GTIN · Confiança · Motivo · Ação), 20 rows. GTIN "✓ igual" only on ean-corroborated rows, "—" on title-match — honest, never fabricated.
- **P2 RESOLVIDOS** — real (not stub): 20 rows from `listProductLinkWorkflows` (SKU HUB 90001/90002/90003/90006… = D-95 EAN linkages), Estado "Vinculado ✓", real resolution timestamps, per-row "Desfazer" enabled (audit_id resolved). Desfazer NOT clicked — live `undoProductLinkResolution` mutation; would corrupt the fixture's 20 resolved links. Presence + enabled + audit-wiring verified structurally.
- **M-09 fold** — default tab "Fila" (`aria-selected=true`); the Anúncios "sem vínculo" deep-link lands on the queue.
- Zero console errors. NO_CANDIDATE honest row data-gated (fixture has all-candidate rows), golden-tested (VinculosDesign.golden.test.tsx).

## Verdict

**CHIP-VINC-DESIGN: PASS.** P6 dual-gate AGREEMENT + P7 design-fidelity live PASS (theme, GTIN honesty, real Resolvidos, Fila-default fold). Cleared for `--no-ff` merge into main.
