# CHIP-M07-FE-DEGRAU — validation

**Branch:** `chip/m07-fe-degrau` · **base** `9e84fe0f` · **head** `72190e6f` (1 commit)
**Owned seam:** `apps/web/src/pages/precos/**` (respected — 4 files, all in seam).

## Defect fixed
FE hardcoded `comissao_pct` per modalidade and always sent it on decompose + solve →
backend forced MANUAL override (degrau 0), gating off the live tariff resolver chain
(degrau 3 COTACAO / degrau 4 PADRAO). Fix: omit `comissao_pct` (optional in contract);
only `modalidade` is sent. Carimbos (fonte/degrau/data/ESTIMATIVA) already render and now
show COTACAO/PADRAO honestly.

## Changed paths (vs write-set = apps/web/src/pages/precos/**)
- PricingPage.tsx — declared, in seam
- PricingPage.test.tsx — declared, in seam
- SolverPanel.tsx — declared, in seam
- SolverPanel.test.tsx — declared, in seam
No changed-undeclared paths. No backend/OpenAPI/SDK touched.

## Verification (evidence type: ran)
- **vitest** apps/web — 45 files / 359 tests GREEN (chip-local config w/ fs.allow +
  absolute setupFiles for the node_modules junction; **deleted pre-commit**, not in diff).
  PricingPage 7/7, SolverPanel 16/16.
- **vite build** — GREEN, built in 4.55s.
- **tsc --noEmit** — raw 450 in the junctioned worktree, **entirely environmental**
  (436 jest-dom matcher augmentation + 4 ImportMeta.env + 10 @mc cross-branch client-type;
  base 9e84fe0f predates tsconfig fix @4a9518a). ZERO errors on the changed surface —
  diff is tsc-neutral (grep for comissao/SolverPanelProps/MODALIDADES/modalidade → none).
  Canonical baseline=2 lives in the hub proper-env.

## Dual gate (Claude-only; codex quota dead til 2026-07-25) — AGREEMENT
- Gate A cold Opus subagent (model=opus): **PASS**, 0 findings.
- Gate B adversarial sonnet subagent: **PASS**, 0 findings.
- Reconciliation: agreement, no disagreement. Both confirmed no lingering `comissao_pct`
  sender and that the test mocks are real (useClient wiring, not disconnected stubs).

## Known post-fix behavior (NOT defects)
- solve for ERP-xlsx products → degrau 4 PADRAO (13.00): solve has no price basis, ERP
  import carries no sale price. Honest by design.
- decompose shows degrau 3 only for products whose EAN/title matches an ML category (e.g. 90001).

## Out of scope (untouched)
Backend/OpenAPI/SDK; manual-override input UI (post-demo); other pages.

## Handoff to hub
Browser P7 live-drive (decompose 90001 shows COTACAO degrau 3; ERP product shows PADRAO
degrau 4) is the hub's post-merge step. No push performed.
