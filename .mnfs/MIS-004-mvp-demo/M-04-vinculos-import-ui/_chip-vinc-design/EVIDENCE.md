# CHIP-VINC-DESIGN — Validation Result

**Chip:** CHIP-VINC-DESIGN (design-fidelity correction, `pages/vinculos/**`)
**Criterion:** MIS-004-C10 (design fidelity), SCREEN-INVENTORY gap #2
**Branch:** `chip/vinc-design`  **Base:** `783cbc0d55ad91fe9cb73a97e8329d589587a0be` (main tip)
**Date:** 2026-07-19

## Markers

P6-DUAL-GATE: AGREEMENT
LIVE-VERIFIED: pending

`LIVE-VERIFIED: pending` — this is a FE design-only change (no provider write, no live
ML path). Live/browser QA (P7) is hub-driven; the chip does not self-drive it.

## Scope

13 files, all under `apps/web/src/pages/vinculos/**` (+593 / −207). No files outside
scope (AppRouter, precos, Anuncios*, produto, dashboard all untouched — confirmed by both
gate reviewers via `git diff --name-only`).

Modified: BatchPreviewModal.tsx, BatchResultFeedback.tsx, BulkBar.tsx, ImportacaoSection.tsx,
QueueRow.tsx, QueueTab.tsx, VinculoDrawer.tsx, VinculosPage.tsx, QueueTab.test.tsx,
VinculosPage.test.tsx.
New: ResolvidosTab.tsx, useVinculosResolved.ts, VinculosDesign.golden.test.tsx.

## Deliverables (fallback order T-1)

- **P0 THEME (required, shipped):** all literal numbered Tailwind color utilities
  (slate/blue/emerald/red/amber/…-N) removed; only paper+green `@theme` tokens remain.
  Grep over the dir → zero offenders. Golden test `OFF_THEME` regex scans rendered DOM
  and asserts `offenders === []`. Confidence bands keep semantic mapping via tokens:
  ALTA→accent-soft/accent-ink (verde), MEDIA→amber-soft/amber, BAIXA→warn-soft/warn.
- **P1 REORIENT (data supported, shipped):** table is anúncio-cêntrica 9-col —
  (checkbox)·Anúncio ML·SKU ML·Produto sugerido·SKU HUB·GTIN·Confiança·Motivo·Ação.
  All columns sourced from `ProductLinkCandidateItem` (FE-only, no backend change).
  GTIN "✓ igual" gated on `match_input === "ean" && Boolean(match_value)`, else "—"
  (never fabricated). EXEMPLO-IO row asserted by golden test.
- **P2 RESOLVIDOS (real, shipped):** ResolvidosTab.tsx queries `listProductLinkWorkflows`,
  filters `current_link.state === "resolved"`, renders real rows; Desfazer wired to
  `undoProductLinkResolution(audit_id)` via `useVinculosResolved.resolutionAuditId`
  (latest `next_state === "resolved"` audit entry); disabled when no audit id. No stub,
  no fabricated resolved rows. Old TODO stub removed from VinculosPage.
- **ADR-17 honest states (hard, held):** NO_CANDIDATE row → "Sem candidato", conf not a
  verde chip, GTIN "—", actions = disabled "Criar produto" (no create seam) + "Ignorar";
  no "Vincular". Nothing fabricated.
- **M-09 fold:** VinculosPage default tab = "fila" — /vinculos "sem vínculo" deep-link
  from Anúncios lands correctly.

## Verification ladder

- **tsc (apps/web):** 10 errors total = main baseline exactly; 0 new; 0 in vinculos
  (`npx tsc --noEmit | grep vinculos` → none).
- **vitest (vinculos):** 5 files / 17 tests PASS (golden 3 + QueueTab 5 + VinculosPage 4
  + ImportacaoSection 3 + BatchPreviewModal 2).
- **No lint lane** in repo (FE ladder = tsc + vitest).
- Chip-local `vitest.config.chip.ts` deleted before commit (hard rule).

## P6 Dual Gate

- **Cold gate (Opus, read-only `harness:gate-reviewer`):** VERDICT: PASS — all 8 contract
  items PASS with file:line evidence.
- **Adversarial gate (sonnet, refutation-framed):** VERDICT: PASS (could not refute) —
  attacked theme/GTIN-gating/band-thresholds/undo-audit-wiring/test-theater/scope; no
  surviving defect.
- **Result: AGREEMENT.**

## Next

Hub drives P7 browser QA (light + dark, 1:1 vs DESIGN-REFERENCE). Chip does not self-drive.

## FINDING — stop-gate.sh unscoped find (harness 0.4.0)

`hooks/stop-gate.sh:39` runs `find "$CWD/.mnfs" -name EVIDENCE.md -path '*_chip*' | head -1`
then greps THAT file for `P6-DUAL-GATE: AGREEMENT`. The `find | head -1` is NOT scoped to
the closing chip — in a fresh worktree, every historical chip pack committed in base is
present, and directory-walk order returns the first one (`M-02-price-intel-core/_chip-m02-fix/
EVIDENCE.md`), which predates the marker convention and lacks the line. Result: the gate
checks the wrong pack and blocks EVERY new chip's CLOSED, regardless of that chip's own
(correct) pack. This chip's pack at `_chip-vinc-design/EVIDENCE.md` DOES carry the marker.

Proposed fix (hub/operator to ratify): scope the find to the current chip — derive the chip
slug from the branch (`chip/<name>`) or `$CWD` and check `_chip-<name>/EVIDENCE.md`
specifically; or require the marker in the NEWEST pack; not `head -1` of an unordered walk.
