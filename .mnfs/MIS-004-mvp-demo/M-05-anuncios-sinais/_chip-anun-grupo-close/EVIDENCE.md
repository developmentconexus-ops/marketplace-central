# CHIP-ANUN-GRUPO — Hub Close Evidence Pack

```yaml
id: CHIP-ANUN-GRUPO
type: hub-close-evidence
author: Dispatch Hub (Wave C)
mission: MIS-004-mvp-demo
milestone-owner: M-05-anuncios-sinais
contract: MIS-004-C10 (design-fidelity sweep) + SCREEN-INVENTORY gap #3
created: 2026-07-19
```

- **Branch:** chip/anun-grupo
- **Tip:** bd5589a0 (LIVE-VERIFIED fill) · feat @ 1450c832
- **Base:** 783cbc0d
- **Full chip pack (brought in by this merge):** `.mnfs/MIS-004-mvp-demo/_chip-anun-grupo/EVIDENCE.md`
- **Scope:** `apps/web/src/pages/AnunciosTable.tsx` (+ `.test.tsx`) ONLY — FE render-only, zero backend/OpenAPI/SDK/migration.

## Merge-gate markers (harness 0.4.0)

P6-DUAL-GATE: AGREEMENT

Dual gate ran during the chip (cold Opus gate-reviewer PASS + adversarial sonnet no-blocking-defect, 0🔴) → agreement, recorded in the chip's EVIDENCE.md §"P6 dual gate". Hub verified the marker + diff scope (2 code files, in-bounds `pages/AnunciosTable.tsx`) at close.

LIVE-VERIFIED: 2026-07-19 hub P7 live-drive on the clean docker dev-stack (frontend :5174 + backend :8080 on the ANUN worktree mount; mount confirmed by AnunciosTable.tsx carrying "ERP est"). `/anuncios` → "Agrupar por produto" ON rendered **32 group-header rows**, each = ▾ chevron + product title + "ERP est. —" + "N anúncio(s)" + "✓ ok" pill. ADR-17 honesty proven live: "ERP est. —" on all 32 (fabricated "ERP est. 0" ABSENT from the DOM); pluralization correct ("1 anúncio"); error pill "✓ ok" accent because the fixture is all `sincronizado` (errorCount=0 → never a fabricated red). Chevron collapse exercised: click ▾ → child row removed (67→66), `aria-expanded` true→false, ▾→▸, header survives. Zero console errors. Theme paper `rgb(251,250,247)` + Instrument Sans; header `bg-surface-2` = `rgb(244,242,234)`.

## Pre-merge pill-integrity concern — RESOLVED by code

Adversarial ❓ (chip EVIDENCE.md): "`listing_count` label vs `listings.length` errorCount could drift if a group's listings were a partial subset. Flagged for hub confirm." Hub verified: `apps/server_core/internal/modules/listings/application/read_service.go:365 finalizeGroups` sets `groups[i].ListingCount = len(groups[i].Listings)`; the fact-filtered `scanGroups` path sets `group.ListingCount = len(survivors)` at :325. Both the label (M) and the error count derive from the same reconciled `Listings` slice — the by-product endpoint ships full listings per group, no intra-group pagination. **Drift is impossible by construction.** Concern closed, not merely accepted.

## Design FINDING (advisory, non-blocking) — carried forward

Design's literal "ERP est. **N**" wants a product-ERP-stock datum absent from `ListingGroup` (SDK `index.ts:368-374` carries only product_id/title/listing_count/group_state/listings[]). Seam FROZEN T-1 → chip shipped honest "ERP est. —". To show a real number post-demo, operator/hub decide: unfreeze + separate backend chip adding ERP-stock to the by-product rollup. Non-blocking for demo.

## Verdict

**CHIP-ANUN-GRUPO: PASS.** P6 dual-gate AGREEMENT + P7 live-drive PASS + pill-integrity concern resolved by code. Cleared for `--no-ff` merge into main.
