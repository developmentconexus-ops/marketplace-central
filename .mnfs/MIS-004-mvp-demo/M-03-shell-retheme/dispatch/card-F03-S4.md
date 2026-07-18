SLICE CARD — F03-S4 · feature F-03-shared-primitives · milestone M-03

depends_on: [F01-S1] (semantic @theme tokens committed). Independent of F-02, F03-S2, F03-S3 (disjoint files).

goal: Retheme the three EXISTING complex primitives — ProductPicker / PaginatedTable / DetailPanel — from the dead
slate/blue palette to the paper+green semantic tokens. PURE VISUAL retheme + real token assertions. Do NOT change any
component API, prop, variant, structure, copy string, aria-label, keyboard behavior, pagination math, selection logic, or
the currency formatter — ONLY the class tokens (and the two radius tokens noted below). These are the design's table +
drawer surfaces; existing consumers must keep compiling and behaving byte-for-byte.

complexity: standard (mechanical retheme across 3 files + additive token assertions).

write_set (retheme 3 primitives + extend their 3 existing test files — NO new files):
- packages/ui/src/ProductPicker.tsx
- packages/ui/src/PaginatedTable.tsx
- packages/ui/src/DetailPanel.tsx
- packages/ui/src/ProductPicker.test.tsx   (extend — keep EVERY existing assertion, ADD token assertions)
- packages/ui/src/PaginatedTable.test.tsx  (extend — keep EVERY existing assertion, ADD token assertions)
- packages/ui/src/DetailPanel.test.tsx     (extend — keep EVERY existing assertion, ADD token assertions)

────────────────────────────────────────────────────────────────────────
EXACT PALETTE MAPPING — apply this table to EVERY occurrence across the 3 primitives; change ONLY these class tokens and
leave all other JSX / props / text / handlers / imports byte-identical. (Occurrence counts are given so you can confirm
you caught them all; they are not a scope limit — retheme every match.)

| raw token (FORBIDDEN after edit) | semantic replacement |
|---|---|
| text-slate-500, text-slate-600            | text-muted        |
| text-slate-700, text-slate-900            | text-ink          |
| text-slate-400                            | text-faint        |
| hover:text-slate-600                      | hover:text-ink    |
| border-slate-100, border-slate-200, border-slate-300 | border-border |
| divide-slate-100, divide-slate-200        | divide-border     |
| bg-white                                  | bg-surface        |
| bg-slate-50                               | bg-surface-2      |
| hover:bg-slate-50                         | hover:bg-surface-2|
| bg-blue-50                                | bg-accent-soft    |
| border-t-blue-600                         | border-t-accent   |
| focus:ring-blue-500                       | focus:ring-accent |
| focus:border-blue-500                     | focus:border-accent|
| rounded-lg                                | rounded-control   |
| rounded-xl                                | rounded-card      |

Per-file occurrence inventory to reconcile against (from the base source):
- ProductPicker.tsx: bg-blue-50 ×1, bg-slate-50 ×1, border-slate-200 ×1, border-slate-300 ×6, border-t-blue-600 ×1,
  divide-slate-100 ×1, divide-slate-200 ×1, focus:border-blue-500 ×3, focus:ring-blue-500 ×3, hover:bg-slate-50 ×1,
  rounded-lg ×4, text-slate-400 ×3, text-slate-500 ×1, text-slate-600 ×6, text-slate-700 ×6, text-slate-900 ×1.
- PaginatedTable.tsx: bg-slate-50 ×1, border-slate-100 ×1, border-slate-200 ×3, border-slate-300 ×1, border-t-blue-600 ×1,
  hover:bg-slate-50 ×2, rounded-lg ×2, rounded-xl ×1, text-slate-400 ×1, text-slate-500 ×2, text-slate-600 ×3.
- DetailPanel.tsx: bg-white ×1, border-slate-100 ×2, border-slate-200 ×1, hover:text-slate-600 ×1, text-slate-400 ×1,
  text-slate-500 ×1, text-slate-900 ×1.

Notes:
- The bare `rounded` on the checkbox inputs in ProductPicker (e.g. `rounded border-slate-300`) is NOT a forbidden token —
  keep `rounded`, only swap its `border-slate-300` → `border-border`.
- Non-color utilities MUST survive untouched: `shadow-xl`, `z-40`, `animate-spin`, `overflow-x-auto`, `overflow-y-auto`,
  every flex/grid/spacing/`min-w-*`/`truncate`/`shrink-0`/`sticky` class, and the lucide `Search`/`X`/`ChevronLeft`/
  `ChevronRight` imports.

────────────────────────────────────────────────────────────────────────
BEHAVIOR / API FREEZE (retheme = class tokens only; everything below is byte-identical after your edit):
- ProductPicker.tsx: `formatCurrency` stays EXACTLY as-is — pt-BR `toLocaleString` BRL. It must NOT coerce, round, or
  fabricate a value; an unknown/absent price is NEVER rendered as R$ 0,00 or a green/healthy state (ADR-17 fail-honest).
  Search input, filter, checkbox selection state, `onSelectionChange`/callbacks, the selected-row highlight condition,
  and every prop signature are unchanged. Only the class strings change.
- PaginatedTable.tsx: pagination math (`totalPages`/`safePage`/`start`/`end`/`pageItems`), `handlePrev`/`handleNext`, the
  `useEffect(() => setPage(1), [items])` reset, the `loading`/empty branches, and all copy ("Loading...", "No items
  found.", "Showing {start+1}–{end} of {items.length}", "Page {safePage} of {totalPages}", "Prev", "Next") + aria-labels
  ("Prev page"/"Next page") are unchanged. `renderHeader`/`renderRow`/`emptyState`/`pageSize` props unchanged.
- DetailPanel.tsx: the Escape-key `useEffect` handler, `open` early-return, `role="complementary"`, `aria-label={title}`,
  `closeLabel` default "Close panel", the `width` prop + inline `style={{ width }}`, `onClose`, and the title/subtitle/
  children/footer structure are unchanged. Only the wrapper/header/close-button/footer class tokens change.

FORBIDDEN palette tokens anywhere in the 3 primitives after your edit: `slate-`, `blue-`, `emerald-`, `red-`, `bg-white`,
`#0F172A`. Use ONLY the semantic replacements in the mapping table above.

────────────────────────────────────────────────────────────────────────
TEST EXTENSIONS (extend the three existing files — REAL toHaveClass assertions on rendered DOM, no theater; EVERY current
assertion stays and stays green — additive only, delete/weaken nothing):

ProductPicker.test.tsx — ADD:
- Render with a small item list and a selected item; assert the selected row carries `bg-accent-soft` (the retheme of the
  old bg-blue-50 highlight). Assert a container/surface element carries `bg-surface` or `border-border`.
- HONEST CURRENCY: render an item with a known numeric price and assert its formatted BRL string appears verbatim (e.g.
  the value the existing formatCurrency produces). Do NOT add any assertion that expects "R$ 0" for an unknown/absent
  price — honest passthrough only.
- (Env note below: run this file IN ISOLATION; the pre-existing "renders product rows" case is a known parallel-run flake,
  not your concern — do not modify or "fix" it.)

PaginatedTable.test.tsx — ADD:
- Assert the table wrapper `<div>` carries `border-border` AND `rounded-card`; the `<thead>` carries `bg-surface-2`.
- Assert the "Showing …" range `<span>` (or the page counter) carries `text-muted`. Keep the existing pagination-behavior
  and copy assertions unchanged.

DetailPanel.test.tsx — ADD:
- Render `open`; assert the `role="complementary"` panel carries `bg-surface` AND `border-border`. Optionally assert the
  title carries `text-ink` and the close button `text-faint`. Keep the existing open/close + Escape-key + aria assertions
  unchanged.

(Use `toHaveClass(...)` on the real rendered node — never assert against a string literal or a tautology.)

────────────────────────────────────────────────────────────────────────
INTEGRITY / HARD CONSTRAINTS:
- Retheme = class tokens only (+ the two radius tokens rounded-lg→rounded-control, rounded-xl→rounded-card). Do NOT change
  any prop, variant, label/copy string, aria-label, keyboard handler, pagination math, selection logic, or formatCurrency.
  Structure byte-identical otherwise. Existing consumers must compile with ZERO edits (public API preserved).
- ProductPicker currency stays honest: unknown ≠ 0, never green (ADR-17).
- Do NOT touch the index.ts barrel (F03-S6 owns it), MarginChip, the F03-S2 state primitives, the F03-S3 base primitives
  (Button/SurfaceCard/Badge/StatCard), or any file outside the write_set. Forbidden regardless: apps/server_core/**,
  packages/sdk-runtime/**, contracts/**, migrations, OpenAPI, apps/web/src/app/**, apps/web/src/routes/**,
  apps/web/src/pages/**.
- No new npm dependency. No new import beyond what already exists in each file.
- Never read/print/commit any .env* file.

validation_kind: render test / typecheck / grep-assertion.

commands (run from worktree root; SANDBOX NOTE in bindings applies to test+build — record could-not-run(sandbox) and STOP,
do not burn the fixup; chip re-runs as verification of record):
- `npm run test -w @marketplace-central/web -- src/ProductPicker.test.tsx src/PaginatedTable.test.tsx src/DetailPanel.test.tsx`
   (run these three files IN ISOLATION per F-ENV-5 — the ProductPicker case flakes only under the full parallel run)
- `npx tsc --noEmit -p apps/web/tsconfig.json`  (expect ONLY the known baseline errors; any NEW error in these files is yours)
- `rg -n 'slate-|blue-|emerald-|red-|bg-white|#0F172A' packages/ui/src/ProductPicker.tsx packages/ui/src/PaginatedTable.tsx packages/ui/src/DetailPanel.tsx`  → EXPECT NO MATCHES
- `rg -n 'bg-surface|bg-surface-2|border-border|divide-border|bg-accent-soft|border-t-accent|focus:ring-accent|text-ink|text-muted|text-faint|rounded-control|rounded-card' packages/ui/src/{ProductPicker,PaginatedTable,DetailPanel}.tsx`  (semantic tokens landed)
- `rg -n 'toLocaleString|formatCurrency' packages/ui/src/ProductPicker.tsx`  (currency formatter preserved — no coercion)

expected_artifacts: 3 primitives using ONLY semantic tokens (zero dead-palette refs); selected-row highlight = bg-accent-soft;
tables/drawer on bg-surface/border-border with rounded-card/rounded-control; ProductPicker currency honest + unchanged;
public APIs/behavior byte-identical; 3 test files extended with real token assertions and ALL prior assertions still green;
no barrel/MarginChip/other-primitive touch; no new tsc errors vs baseline.

open_questions: [] (dispatch-ready)

Report back: status; changed paths vs write_set; each command with evidence type (ran/assumed/could-not-run) + captured
output; confirm no API/behavior/copy/aria/formatCurrency/barrel change and currency stays honest; anything you did NOT verify.
