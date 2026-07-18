SLICE CARD — F03-S3 · feature F-03-shared-primitives · milestone M-03

depends_on: [F01-S1] (semantic @theme tokens committed). Independent of F-02 and F03-S2 (disjoint files).

goal: Retheme the four base primitives Button / SurfaceCard / Badge / StatCard from the dead blue/emerald/red/slate
palette to the paper+green semantic tokens. PURE VISUAL retheme + real token assertions. Do NOT change component APIs,
props, variants, structure, copy, or behavior — only the class tokens (and the radius tokens noted below).

complexity: standard.

write_set (retheme 4 primitives + extend Badge.test + add BasePrimitives.test):
- packages/ui/src/Button.tsx
- packages/ui/src/SurfaceCard.tsx
- packages/ui/src/Badge.tsx
- packages/ui/src/StatCard.tsx
- packages/ui/src/Badge.test.tsx        (extend — keep every existing assertion, ADD token assertions)
- packages/ui/src/BasePrimitives.test.tsx (NEW — Button/SurfaceCard/StatCard render + token + behavior)

────────────────────────────────────────────────────────────────────────
EXACT PALETTE MAPPING (change ONLY these tokens; leave APIs/props/structure/children byte-identical):

Button.tsx — the `variantClasses` map and the base rounded token:
- primary:   `bg-blue-600 hover:bg-blue-700 text-white border-transparent` → `bg-accent hover:bg-accent-ink text-white border-transparent`
- secondary: `bg-white hover:bg-slate-50 text-slate-700 border-slate-200`   → `bg-surface hover:bg-surface-2 text-ink border-border`
- danger:    `bg-red-600 hover:bg-red-700 text-white border-transparent`    → `bg-warn hover:opacity-90 text-white border-transparent`
    (there is no darker-warn token; use `hover:opacity-90` — do NOT invent a shade.)
- base class: `rounded-lg` → `rounded-control`. Keep everything else in the base className string
  (`inline-flex items-center gap-2 px-3 py-2 text-sm font-medium border cursor-pointer transition-colors duration-150
  disabled:opacity-50 disabled:cursor-not-allowed`), the `loading` spinner SVG, the `isDisabled` logic, the
  `type`/`variant`/`loading`/`disabled`/`className` prop handling — ALL byte-identical.

SurfaceCard.tsx — `bg-white border border-slate-200 rounded-xl` → `bg-surface border border-border rounded-card`.
Keep `p-6`, the `className` merge, and the `<section>` element EXACTLY.

Badge.tsx — the `config` status map classes and the base rounded token:
- pending:     `bg-slate-100 text-slate-600`   → `bg-surface-2 text-muted`
- in_progress: `bg-blue-100 text-blue-700`     → `bg-info-soft text-info`
- succeeded:   `bg-emerald-100 text-emerald-700`→ `bg-accent-soft text-accent-ink`
- failed:      `bg-red-100 text-red-700`       → `bg-warn-soft text-warn`
- completed:   `bg-emerald-100 text-emerald-700`→ `bg-accent-soft text-accent-ink`
- base class: `rounded` → `rounded-pill`. Keep every `label` string ("Pending"/"In Progress"/"Succeeded"/"Failed"/
  "Completed"), the `inline-flex items-center px-2 py-0.5 text-xs font-medium` base, and the `className` merge EXACTLY.

StatCard.tsx —
- `bg-white border border-slate-200 rounded-xl` → `bg-surface border border-border rounded-card`
- label `text-slate-500` → `text-muted`
- value `text-slate-900` → `text-ink` (keep the `style={{ fontFamily: "var(--font-mono)" }}` OR equivalently replace it
  with the `font-mono` utility class — either is acceptable; the monospace value must survive)
- sub `text-slate-400` → `text-faint`
Keep the label/value/sub structure, the `sub &&` conditional, `uppercase tracking-wide`, spacing, and props EXACTLY.

FORBIDDEN palette tokens anywhere in these 4 files after your edit: `blue-`, `emerald-`, `red-`, `slate-`, `bg-white`,
`#0F172A`. Use ONLY semantic utilities: bg-surface, bg-surface-2, bg-accent, bg-accent-soft, text-accent-ink, bg-warn,
bg-warn-soft, text-warn, bg-info-soft, text-info, text-ink, text-muted, text-faint, border-border, rounded-control,
rounded-card, rounded-pill, font-mono.

────────────────────────────────────────────────────────────────────────
TEST EXTENSIONS (real assertions, no theater; all prior assertions stay green):

Badge.test.tsx — keep the 5 existing label-render tests; ADD token assertions:
- pending → `bg-surface-2` `text-muted`; in_progress → `bg-info-soft` `text-info`; succeeded → `bg-accent-soft`
  `text-accent-ink`; failed → `bg-warn-soft` `text-warn`; completed → `bg-accent-soft` `text-accent-ink`.
  (Use `screen.getByText(<label>)` then `toHaveClass(...)`.)

BasePrimitives.test.tsx (NEW) — cover:
- Button: renders children; default variant is secondary → `bg-surface`/`text-ink`/`border-border`; primary →
  `bg-accent`/`text-white`; danger → `bg-warn`/`text-white`; base has `rounded-control`. `loading` sets the button
  `disabled` and renders the spinner (query the button by role and assert `toBeDisabled()`); `disabled` prop disables.
  An `onClick` fires when enabled and does NOT fire when disabled/loading.
- SurfaceCard: renders children; the `<section>` has `bg-surface`/`border-border`/`rounded-card`; a passed `className`
  is merged (appears on the section).
- StatCard: renders `label`, `value`, and (when provided) `sub`; omitting `sub` renders no sub paragraph; the value
  element is monospace (either `font-mono` class OR inline `fontFamily: var(--font-mono)` — assert whichever the impl
  uses); value/label use `text-ink`/`text-muted`. StatCard renders `value` verbatim — assert a passed value like `"—"`
  or `"1.234"` appears exactly (it must NOT coerce/round/fake a number; honest passthrough).

────────────────────────────────────────────────────────────────────────
INTEGRITY / HARD CONSTRAINTS:
- Retheme = class tokens only (+ the three radius tokens). Do NOT change any variant name, prop, label string, or
  behavior. Structure byte-identical otherwise.
- StatCard must render `value` verbatim (no coercion of unknown to 0 — honest passthrough; the caller supplies "—" for
  unknown via UnknownValue).
- Do NOT touch the index.ts barrel (F03-S6 owns it), MarginChip, the F03-S2 state primitives, or any file outside the
  write_set. Forbidden regardless: apps/server_core/**, packages/sdk-runtime/**, contracts/**, migrations, OpenAPI,
  apps/web/src/pages/**.
- No new npm dependency. No new import beyond React types already present.
- Never read/print/commit any .env* file.

validation_kind: render test / typecheck / grep-assertion.

commands (run from worktree root; SANDBOX NOTE in bindings applies to test+build — record could-not-run(sandbox) and STOP,
do not burn the fixup; chip re-runs as verification of record):
- `npm run test -w @marketplace-central/web -- src/Badge.test.tsx src/BasePrimitives.test.tsx`
   (the web vitest config globs packages/ui tests; if a relative path fails, run the whole packages/ui glob)
- `npx tsc --noEmit -p apps/web/tsconfig.json`   (expect ONLY the known baseline errors; any NEW error in these files is yours)
- `rg -n 'blue-|emerald-|red-|slate-|bg-white|#0F172A' packages/ui/src/Button.tsx packages/ui/src/SurfaceCard.tsx packages/ui/src/Badge.tsx packages/ui/src/StatCard.tsx`  → EXPECT NO MATCHES
- `rg -n 'bg-accent|bg-surface|bg-warn|bg-info-soft|text-ink|text-muted|rounded-control|rounded-card|rounded-pill' packages/ui/src/{Button,SurfaceCard,Badge,StatCard}.tsx`  (semantic tokens landed)

expected_artifacts: 4 primitives using ONLY semantic tokens (zero dead-palette refs); Button variants primary=accent/
secondary=surface/danger=warn on rounded-control; SurfaceCard/StatCard on bg-surface/border-border/rounded-card; Badge
statuses mapped to info/accent/warn softs on rounded-pill; StatCard honest value passthrough; Badge.test extended +
BasePrimitives.test added with real token+behavior assertions, all green; no barrel/other-primitive touch; no new tsc
errors vs baseline.

open_questions: [] (dispatch-ready)

Report back: status; changed paths vs write_set; each command with evidence type (ran/assumed/could-not-run) + captured
output; confirm no API/behavior/barrel change and StatCard value passthrough is honest; anything you did NOT verify.
