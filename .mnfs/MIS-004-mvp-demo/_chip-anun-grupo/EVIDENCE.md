# CHIP-ANUN-GRUPO — validation-result

- Chip: CHIP-ANUN-GRUPO (design-fidelity CORRECTION, milestone owner M-05 Anúncios)
- Contract: MIS-004-C10 (design-fidelity sweep) + SCREEN-INVENTORY gap #3
- Branch: `chip/anun-grupo` @ commit `1450c83` (base `783cbc0d`)
- Scope: `apps/web/src/pages/AnunciosTable.tsx` (+ `AnunciosTable.test.tsx`) ONLY

## Markers

P6-DUAL-GATE: AGREEMENT
LIVE-VERIFIED: pending

(cold Opus gate PASS + adversarial sonnet no-blocking-defect → agreement.
LIVE-VERIFIED pending = hub drives P7 fresh browser QA; chip does not self-drive.)

## What shipped

Enriched the "Agrupar por produto" group-header row (was: title + `(count)` at
`AnunciosTable.tsx:198-208`) to the design spec (design-screens §Anúncios:
"linhas grupo (chevron, meta 'ERP est. N · M anúncios', pill ok/N erro)"):

1. **▸/▾ chevron** — per-group expand/collapse. Collapsed → child listing rows
   hidden (tbody flatMap drops them), header row survives. Groups start expanded
   (preserves flat-list demo behaviour). `aria-expanded` + `aria-label`.
2. **meta "ERP est. — · M anúncios"** — M = `group.listing_count`, pluralized.
3. **pill "✓ ok" / "N erro(s)"** — errorCount = `group.listings.filter(sync_state === "error")`,
   warn pill when > 0 else accent "✓ ok". Pluralized.

## ADR-17 (unknown ≠ zero)

ERP stock per group is **genuinely unavailable FE-only**: `ListingGroup`
(`packages/sdk-runtime/src/index.ts:368-374`) carries only
`product_id / product_title / listing_count / group_state / listings[]` — **no
ERP-stock field**. Seam FROZEN T-1 → chip did NOT add a backend field. Row
renders honest **"ERP est. —"**, never a fabricated `0`. Golden test guards that
`ERP est. 0` is absent. Error pill counts only real `error` sync_states (not
stale/queued). `published_quantity` is ML-published qty, not ERP stock — not
substituted.

## FINDING / ESCALATION (advisory, non-blocking — row ships either way)

Design's literal "ERP est. **N**" requires a product-ERP-stock datum not present
on `ListingGroup`. To show a real number live, backend must add an ERP-stock
field to the by-product group rollup — a frozen-seam change the chip is barred
from. Operator/hub decide: unfreeze + separate backend chip, OR accept
"ERP est. —" for demo. Chip proceeds with the honest fallback now.

## Tests (TDD, golden per dispatch EXEMPLO-IO)

`apps/web/src/pages/AnunciosTable.test.tsx` — new describe "enriched group header":
- GOLDEN: 3 listings / 1 sync error → "ERP est. —", "3 anúncios", "1 erro"; not "ERP est. 0".
- GOLDEN: 3 listings all synced → "ERP est. —", "3 anúncios", "✓ ok", no error text.
- pluralization: 1 → "anúncio"; 2 errors → "erros".
- chevron collapse: child rows removed on click, `aria-expanded` flips, header stays.

Run (chip-local vitest config w/ fs.allow + absolute setupFiles over node_modules
junctions to main; config deleted pre-commit):

```
✓ src/pages/AnunciosTable.test.tsx (20 tests) 236ms
Test Files  1 passed (1)  |  Tests  20 passed (20)
```

Typecheck (`apps/web tsc --noEmit`): 10 `error TS` = main baseline (memory
web-tsc-lane), **zero net new**, none at chip lines.

## P6 dual gate

| Gate | Model | Framing | Verdict |
|------|-------|---------|---------|
| Cold | Opus (harness:gate-reviewer, read-only) | contract compliance | **PASS** — all Required criteria proven file:line; agrees "ERP est. —" is correct honest fallback |
| Adversarial | Sonnet (refutation-framed) | find fault, default reject | **no blocking defect** (0🔴 1🟡 2❓ 1🔵) |

**AGREEMENT: PASS.**

Adversarial non-blocking notes (dispositions):
- 🟡 null-product-id groups share collapse key `"sem-produto"` → both toggle
  together. By design the synthetic no-product bucket is a **single** group per
  page (pre-existing React-key assumption, unchanged). Not triggerable.
- ❓ error pill re-derives instead of using `group_state`. **Contract-mandated**:
  dispatch requires the pill be a *count* derived from per-listing `sync_state`;
  `group_state` is tri-state, gives no count.
- ❓ `listing_count` (label) vs `listings.length` (errorCount source) could drift
  if a group's `listings` were ever a partial subset. By-product endpoint ships
  full listings per group; not paginated within a group. Flagged for hub confirm.
- 🔵 chevron lacks `aria-controls`. Acceptable; header/rows are sibling `<tr>`s.

## Zero-writes / hard rules

FE render-only. Zero ML writes. No server booted. No push. Branch `chip/anun-grupo`.
node_modules junctions to main + chip-local vitest config used for the test lane
then removed; `git status` pre-commit = exactly the 2 intended files.
