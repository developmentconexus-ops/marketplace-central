SLICE CARD — F03-S5 · feature F-03-shared-primitives · milestone M-03

depends_on: [F01-S1] (semantic @theme tokens committed). Independent of F-02, F03-S2/S3/S4 (new file, disjoint).

goal: Create the NEW design-canonical `DataTable<T>` shell primitive in packages/ui — a purely PRESENTATIONAL, fully
CONTROLLED table (columns / rows / selection / sort all driven by props + callbacks; ZERO internal data state, ZERO fetch,
ZERO TanStack). This is the shared table chrome the design screens (Anúncios, Estoque, Simulador, Pedidos, Vínculos) all
reuse; the actual data/columns/formatting live in the consuming PAGES (later milestones), NOT here. Paper+green semantic
tokens only. This is a SHELL, not a data component — do NOT bake in domain columns, grouping, pagination, or formatting.

complexity: standard. IMPORTANT: because selection and sort are CONTROLLED (props in, callback out), the component holds
NO useState/useReducer for data — it is a pure function of its props. Do NOT build an internal selection/sort state
machine; that is the page's job. (An internal state machine here = anti-slop REJECT.)

write_set (NEW files only — do NOT touch the barrel; F03-S6 exports it):
- packages/ui/src/DataTable.tsx        (NEW)
- packages/ui/src/DataTable.test.tsx   (NEW)

────────────────────────────────────────────────────────────────────────
PUBLIC API (exact — this becomes a contract consumed by waves B/C pages; keep it MINIMAL, no speculative extras):

```ts
import type { ReactNode } from "react";

export interface DataTableColumn<T> {
  key: string;                         // stable column id (also the sort key)
  header: ReactNode;                   // header cell content
  render: (row: T) => ReactNode;       // cell content — the PAGE owns all formatting + honest-value handling
  align?: "left" | "right";            // default "left"; use "right" for numeric columns
  sortable?: boolean;                  // when true, header is a sort control that calls onSortChange(key)
}

export interface DataTableProps<T> {
  columns: DataTableColumn<T>[];
  rows: T[];
  rowKey: (row: T) => string;          // stable unique key per row
  onRowClick?: (row: T) => void;       // row click (design: clique na linha abre o drawer)
  // controlled selection (optional — a leading checkbox column renders ONLY when onSelectionChange is provided):
  selectedKeys?: ReadonlySet<string>;
  onSelectionChange?: (next: Set<string>) => void;
  // controlled sort (optional):
  sortKey?: string | null;
  sortDir?: "asc" | "desc";
  onSortChange?: (key: string) => void;
  loading?: boolean;
  emptyState?: ReactNode;              // default: "Nenhum item."
  stickyHeader?: boolean;              // thead sticks to top of the scroll container
}

export function DataTable<T>(props: DataTableProps<T>): JSX.Element;
```

BEHAVIOR (all derived from props; no internal data state):
- Rendering: `<div overflow-x-auto>` → `<table>` → `<thead>` (one header row from `columns`) → `<tbody>` (one row per
  `rows` item via `rowKey`, each cell = `column.render(row)`).
- SELECTION (only when `onSelectionChange` is set): a leading checkbox column.
  · Header checkbox: checked when EVERY row key ∈ `selectedKeys`; set `indeterminate` (via a ref) when SOME but not all
    are selected; unchecked when none. Toggling it calls `onSelectionChange(new Set(all row keys))` when going to all-
    selected, or `onSelectionChange(new Set())` when clearing.
  · Row checkbox: `checked = selectedKeys.has(rowKey(row))`; toggling calls `onSelectionChange(next)` where `next` is a
    NEW Set derived from the current `selectedKeys` with that key added/removed (never mutate the incoming set).
  · A click on ANY checkbox must `stopPropagation` so it does NOT also trigger `onRowClick`.
- SORT: a `sortable` column header renders as a clickable control that calls `onSortChange(column.key)`. When
  `sortKey === column.key`, show a direction indicator (▲ for "asc", ▼ for "desc"); otherwise a neutral affordance. A
  non-sortable header is plain text (no button, no indicator).
- ROW CLICK: when `onRowClick` is set, each row is clickable (cursor-pointer + hover) and calls `onRowClick(row)`.
- SELECTED ROW: when `selectedKeys?.has(rowKey(row))`, the row carries `bg-accent-soft`.
- LOADING: when `loading`, render a centered spinner (NO rows). EMPTY: when `!loading && rows.length === 0`, render
  `emptyState ?? "Nenhum item."` centered. Otherwise render the table.

HONEST-DATA CONSTRAINT (ADR-17, binding): DataTable renders EXACTLY what `column.render(row)` returns — it must NEVER
inject a `0`, `R$ 0`, a placeholder number, or any fabricated/coerced value, and must never colour a cell green/healthy
on its own. Unknown/absent values are the page's responsibility (it passes `<UnknownValue/>` / `<MarginChip/>` into
`render`). The shell is value-agnostic.

────────────────────────────────────────────────────────────────────────
STYLING (semantic tokens ONLY — same token vocabulary as the rethemed tables in F03-S4):
- wrapper: `overflow-x-auto border border-border rounded-card`
- table: `w-full text-sm text-left`
- thead: `bg-surface-2 border-b border-border` (+ `sticky top-0` when `stickyHeader`)
- th: `px-3 py-2 text-xs font-medium text-muted uppercase tracking-wide` (+ `text-right` when the column is right-aligned)
- sortable header control: `inline-flex items-center gap-1 cursor-pointer` (a `<button type="button">`), indicator glyph
  in `text-faint` (or `text-ink` for the active direction)
- tbody row: `border-b border-border last:border-0`; when clickable add `cursor-pointer hover:bg-surface-2`; selected row
  adds `bg-accent-soft`
- td: `px-3 py-2 text-sm text-ink` (+ `text-right` when the column is right-aligned)
- checkbox: `rounded border-border` (bare `rounded` is fine — not a forbidden token)
- spinner (loading): `animate-spin rounded-full h-5 w-5 border-2 border-border border-t-accent` inside a
  `flex items-center justify-center py-12 text-sm text-muted` wrapper
- empty: `flex items-center justify-center py-16 text-sm text-faint`

FORBIDDEN tokens anywhere in DataTable.tsx: `slate-`, `blue-`, `emerald-`, `red-`, `amber-100`, `amber-800`, `bg-white`,
`#0F172A`. Use ONLY: bg-surface, bg-surface-2, bg-accent-soft, border-border, text-ink, text-muted, text-faint,
border-t-accent, rounded-card, rounded-control, rounded-pill (as needed).

────────────────────────────────────────────────────────────────────────
TESTS — DataTable.test.tsx (NEW; real render assertions, no theater):
- Renders each column header and, for a small `rows` array, each cell via `render` (assert real cell text appears).
- `onRowClick`: clicking a row calls it with THAT row; clicking a row's checkbox does NOT call `onRowClick`
  (stopPropagation) but DOES call `onSelectionChange`.
- Selection: with `onSelectionChange` provided, the leading checkbox column renders; the header select-all calls
  `onSelectionChange` with a Set of all row keys; toggling one row checkbox calls `onSelectionChange` with the correct
  next Set (added when unselected, removed when selected); a row whose key is in `selectedKeys` carries `bg-accent-soft`.
  Without `onSelectionChange`, NO checkbox column renders.
- Sort: a `sortable` column header is a control that calls `onSortChange(key)`; when `sortKey` matches, the direction
  indicator reflects `sortDir`; a non-sortable header renders no sort control.
- `loading` → spinner shown, no data rows. `rows: []` (not loading) → `emptyState` (or the default "Nenhum item.") shown.
- Token assertions: wrapper has `border-border` + `rounded-card`; thead has `bg-surface-2`; a selected row has
  `bg-accent-soft`; a right-aligned column's cell has `text-right`.
- HONEST: pass a column whose `render` returns the string "—" and assert "—" appears verbatim; assert the shell injects
  no "0"/"R$ 0" of its own (e.g. `screen.queryByText("0")` is null when no cell rendered it).

────────────────────────────────────────────────────────────────────────
INTEGRITY / HARD CONSTRAINTS:
- Presentational + fully controlled: NO useState/useReducer holding selection/sort/data; NO fetch, NO TanStack, NO
  context. The one allowed `useRef` is solely to set the header checkbox `indeterminate` DOM property (React has no
  prop for it). Nothing else stateful.
- Value-agnostic + honest (ADR-17): never inject/coerce/fake a value; never colour a cell healthy on its own.
- Minimal API: implement EXACTLY the props above — no extra props, no speculative flexibility (density variants, pinned
  columns, virtual scroll, built-in pagination, etc. are OUT — YAGNI / anti-speculative-abstraction).
- Do NOT touch the index.ts barrel (F03-S6 owns it), PaginatedTable/DetailPanel/ProductPicker (F03-S4), or any file
  outside the write_set. Forbidden regardless: apps/server_core/**, packages/sdk-runtime/**, contracts/**, migrations,
  OpenAPI, apps/web/src/app/**, apps/web/src/routes/**, apps/web/src/pages/**.
- No new npm dependency. Import only from "react" (types) and @testing-library/react in the test.
- Never read/print/commit any .env* file.

validation_kind: render test / typecheck / grep-assertion.

commands (run from worktree root; SANDBOX NOTE in bindings applies to test+build — record could-not-run(sandbox) and STOP;
chip re-runs as verification of record):
- `npm run test -w @marketplace-central/web -- src/DataTable.test.tsx`
- `npx tsc --noEmit -p apps/web/tsconfig.json`  (expect ONLY the known baseline errors; any NEW error here is yours)
- `rg -n 'slate-|blue-|emerald-|red-|amber-100|amber-800|bg-white|#0F172A' packages/ui/src/DataTable.tsx`  → EXPECT NO MATCHES
- `rg -n 'useState|useReducer|useQuery|fetch\(' packages/ui/src/DataTable.tsx`  → EXPECT NO MATCHES (controlled/presentational)
- `rg -n 'bg-surface-2|border-border|bg-accent-soft|border-t-accent|rounded-card|text-muted|text-faint' packages/ui/src/DataTable.tsx`  (semantic tokens landed)

expected_artifacts: a NEW generic `DataTable<T>` shell — controlled columns/rows/selection/sort, presentational, semantic
tokens only, honest/value-agnostic, minimal API exactly as specified; DataTable.test.tsx with real render + interaction +
token + honest-value assertions, all green; barrel untouched; no new tsc errors vs baseline; no new dependency.

open_questions: [] (dispatch-ready)

Report back: status; changed paths vs write_set; each command with evidence type (ran/assumed/could-not-run) + captured
output; confirm the component is fully controlled (no internal data state), value-agnostic/honest, minimal-API, and does
NOT touch the barrel; anything you did NOT verify.
