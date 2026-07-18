SLICE CARD — F03-S6 · feature F-03-shared-primitives · milestone M-03

depends_on: [F03-S1 (MarginChip committed), F03-S5 (DataTable committed)]. Uses existing DetailPanel (base + F03-S4 retheme).
This is the LAST F-03 slice and the SOLE owner of the packages/ui barrel.

goal: (1) Create the NEW design-canonical `DetailDrawer` shell — the right-hand detail drawer the design screens use
(Anúncios/Pedidos/Vínculos/Simulador), 300–380px wide, controlled open/close, header + scrollable body + actions footer.
(2) Wire the packages/ui BARREL (`index.ts`) to export the three primitives that the F-03 contract requires but that are
not yet exported: `MarginChip` (F03-S1), `DataTable` (F03-S5), and `DetailDrawer` (this slice). This closes milestone
criterion M-03-C05 ("primitiva do contrato F-03 ausente do export" = blocking).

complexity: standard.

DO-NOT-DUPLICATE DIRECTIVE (binding): DetailPanel ALREADY implements the drawer chrome (fixed right-side positioning,
Escape-to-close, `role="complementary"`, aria-label, close button, header/body/footer regions, width prop). DetailDrawer
MUST COMPOSE DetailPanel — delegate ALL chrome to it — and add ONLY the design-canonical defaults. Do NOT re-implement the
Escape handler, positioning, or aria (that would duplicate tested logic = anti-slop REJECT).

write_set (1 new component + its test + the barrel):
- packages/ui/src/DetailDrawer.tsx        (NEW — composes DetailPanel)
- packages/ui/src/DetailDrawer.test.tsx   (NEW)
- packages/ui/src/index.ts                (EDIT — add 3 export groups; keep every existing export intact)

────────────────────────────────────────────────────────────────────────
DetailDrawer PUBLIC API (minimal; composes DetailPanel):

```ts
import type { ReactNode } from "react";

export interface DetailDrawerProps {
  open: boolean;
  onClose: () => void;
  title: string;
  subtitle?: string;
  children: ReactNode;
  actions?: ReactNode;          // design "ações dinâmicas" — rendered in the drawer footer
  closeLabel?: string;          // default "Fechar" (pt-BR — design copy)
  width?: number;               // default 360 (design band 300–380); pass-through override
}

export function DetailDrawer(props: DetailDrawerProps): JSX.Element | null;
```

IMPLEMENTATION (thin composition):
- Render `<DetailPanel open={open} onClose={onClose} title={title} subtitle={subtitle} closeLabel={closeLabel ?? "Fechar"}
  width={width ?? 360} footer={actions}>{children}</DetailPanel>`.
- That is essentially the whole component. `actions` maps to DetailPanel's `footer`; the default width (360) sits in the
  design's 300–380 band; the default closeLabel is pt-BR "Fechar". Do NOT add positioning/Escape/aria here — DetailPanel
  owns them.
- No internal state, no fetch, no context (DetailPanel already handles the open/Escape lifecycle). Presentational.

────────────────────────────────────────────────────────────────────────
BARREL EDIT — packages/ui/src/index.ts. Current file exports Button/SurfaceCard/Badge/StatCard/LoadingState/ErrorState/
EmptyState/UnknownValue/ConflictTag/FreshnessIndicator (via `export *`), plus named ProductPicker/PaginatedTable/
DetailPanel and their types. KEEP ALL OF THOSE UNCHANGED. ADD exactly these lines (place them logically; the ProductPicker
type line stays exactly as-is):

```ts
export * from "./MarginChip";
export { DataTable } from "./DataTable";
export type { DataTableProps, DataTableColumn } from "./DataTable";
export { DetailDrawer } from "./DetailDrawer";
export type { DetailDrawerProps } from "./DetailDrawer";
```

- Do NOT remove, reorder-destructively, or alter any existing export. Additive only. No other file in the write_set edits
  the barrel; this slice is its sole writer.
- After the edit, `MarginChip`, `DataTable`, and `DetailDrawer` (plus the existing 13 primitives) are all importable from
  `@marketplace-central/ui`.

────────────────────────────────────────────────────────────────────────
TESTS:

DetailDrawer.test.tsx (NEW):
- `open` → renders the `role="complementary"` panel (from DetailPanel) with the `title`, the `subtitle` when provided, the
  `children`, and the `actions` node in the footer region.
- `open=false` → renders nothing (`container` empty / query returns null).
- Close: clicking the close button (aria-label "Fechar") calls `onClose`; pressing Escape calls `onClose` (chrome
  inherited from DetailPanel — assert it still works through the composition).
- Default width: with no `width` prop, the panel's inline style width is 360 (assert the rendered `role="complementary"`
  element's `style.width`); passing `width={320}` overrides it.
- Token sanity: the panel carries the rethemed semantic tokens (`bg-surface`, `border-border`) inherited from DetailPanel
  (F03-S4) — assert at least `bg-surface`.

barrel export-contract test — ADD to DetailDrawer.test.tsx (or a tiny separate `Barrel.test.tsx`, your call, but if new it
must be in the write_set — prefer folding into DetailDrawer.test.tsx to avoid touching the write_set):
- `import { MarginChip, DataTable, DetailDrawer } from "./index";` and assert each is a function (`typeof … === "function"`).
  This is the real guard for M-03-C05: it FAILS if any contract primitive is missing from the barrel.

────────────────────────────────────────────────────────────────────────
INTEGRITY / HARD CONSTRAINTS:
- DetailDrawer composes DetailPanel — no duplicated chrome/Escape/positioning/aria. Presentational, no internal state.
- Barrel edit is ADDITIVE only — every existing export preserved byte-for-byte; the 5 new lines added.
- Do NOT modify DetailPanel, DataTable, MarginChip, or any other primitive (only compose/export them). Do NOT touch any
  file outside the write_set. Forbidden regardless: apps/server_core/**, packages/sdk-runtime/**, contracts/**,
  migrations, OpenAPI, apps/web/src/app/**, apps/web/src/routes/**, apps/web/src/pages/**.
- No new npm dependency. Import only from "react" (types), "./DetailPanel", and (in tests) @testing-library/react + "./index".
- FORBIDDEN palette tokens in DetailDrawer.tsx: slate-/blue-/emerald-/red-/bg-white/#0F172A (there should be little/no raw
  styling here at all since chrome is delegated — any styling you add uses semantic tokens only).
- Never read/print/commit any .env* file.

validation_kind: render test / typecheck / grep-assertion / barrel-export contract.

commands (run from worktree root; SANDBOX NOTE in bindings applies to test+build — record could-not-run(sandbox) and STOP;
chip re-runs as verification of record):
- `npm run test -w @marketplace-central/web -- src/DetailDrawer.test.tsx`
- `npx tsc --noEmit -p apps/web/tsconfig.json`  (expect ONLY the known baseline errors; any NEW error here is yours)
- `rg -n 'MarginChip|DataTable|DetailDrawer' packages/ui/src/index.ts`  (all three exported)
- `rg -n 'ProductPicker|PaginatedTable|DetailPanel|LoadingState|UnknownValue' packages/ui/src/index.ts`  (existing exports intact)
- `rg -n 'slate-|blue-|emerald-|red-|bg-white|#0F172A' packages/ui/src/DetailDrawer.tsx`  → EXPECT NO MATCHES

expected_artifacts: a NEW `DetailDrawer` that composes DetailPanel with design defaults (width 360, pt-BR "Fechar",
actions→footer); the barrel additively exporting MarginChip + DataTable + DetailDrawer (plus the existing 13, untouched);
DetailDrawer.test.tsx with render/close/width/token assertions AND the barrel export-contract test (all three primitives
resolve as functions), all green; no DetailPanel/other-primitive change; no new tsc errors vs baseline; no new dependency.

open_questions: [] (dispatch-ready)

Report back: status; changed paths vs write_set; each command with evidence type (ran/assumed/could-not-run) + captured
output; confirm DetailDrawer composes (not duplicates) DetailPanel, the barrel edit is additive with all prior exports
intact, and MarginChip/DataTable/DetailDrawer all resolve from the barrel; anything you did NOT verify.
