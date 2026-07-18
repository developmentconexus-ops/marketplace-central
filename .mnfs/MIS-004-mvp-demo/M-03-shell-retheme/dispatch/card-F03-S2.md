SLICE CARD — F03-S2 · feature F-03-shared-primitives · milestone M-03

depends_on: [F01-S1] (semantic @theme tokens committed). Independent of F-02.

goal: Retheme the six existing fact/load STATE primitives from the dead slate/red/amber-100 palette to the paper+green
semantic tokens. This is a PURE VISUAL retheme: swap raw Tailwind palette classes for the semantic utilities that resolve
to index.css @theme vars (auto light/dark via data-theme). Do NOT change component copy, structure, props, roles, aria
labels, title logic, or any behavior. The honest-unknown affordances ("—" and "desconhecido") stay visually neutral —
NEVER green/healthy.

SCOPE ADJUDICATION (binding — read before you touch FreshnessIndicator): the P2 plan floated a relative-age freshness
copy ("agora/há N min/há N h…"). That is RETIRED and OUT OF SCOPE. The real copy comes from `formatAsOf` in
`@marketplace-central/web-query` (packages/web-query) and renders `dados de <HH:MM:SS>` / `dados de desconhecido`. You
MUST NOT change `formatAsOf`, MUST NOT touch packages/web-query, and MUST NOT change FreshnessIndicator's text/behavior —
only its wrapper's visual token classes. Changing the copy is a logic change in a shared util under a retheme milestone
and would break existing green tests; do not do it.

complexity: standard (mechanical retheme + real token assertions).

write_set (retheme 6 primitives + extend their 2 existing test files — NO new files):
- packages/ui/src/UnknownValue.tsx
- packages/ui/src/FreshnessIndicator.tsx
- packages/ui/src/EmptyState.tsx
- packages/ui/src/ErrorState.tsx
- packages/ui/src/LoadingState.tsx
- packages/ui/src/ConflictTag.tsx
- packages/ui/src/FactStates.test.tsx   (extend — keep every existing assertion, ADD token assertions)
- packages/ui/src/LoadStates.test.tsx   (extend — keep every existing assertion, ADD token assertions)

────────────────────────────────────────────────────────────────────────
EXACT PALETTE MAPPING (change ONLY these class tokens; leave all other JSX/props/text byte-identical):

UnknownValue.tsx — the honest-unknown em-dash. Currently `<span title={…}>—</span>` (no color). ADD `className="text-faint"`
so unknown reads as a neutral/muted absence, never default-ink and never green. Keep the `title={hint || undefined}`
logic and the "—" glyph EXACTLY.

FreshnessIndicator.tsx — currently `<span aria-label="Atualização dos dados">{formatAsOf(asOf)}</span>`. ADD
`className="text-muted text-xs font-mono"` (metadata styling; monospace timestamp). Keep `aria-label`, keep the
`formatAsOf(asOf)` call, keep the prop signature. Do NOT import anything new. Behavior FROZEN.

EmptyState.tsx — `text-slate-500` → `text-muted` (outer div); `text-slate-400` → `text-faint` (hint <p>). Keep both copy
strings, the `hint` conditional, and the two-<p> structure EXACTLY (a test counts the <p> elements).

ErrorState.tsx — `text-red-700` → `text-warn` (paper+green negative/error token). Keep `role="alert"`, the fixed
"Erro ao carregar." prefix + optional detail, and the `<Button variant="secondary" onClick={onRetry}>Tentar novamente</Button>`
EXACTLY. Do NOT touch Button.tsx (F03-S3 owns it) — ErrorState only consumes it.

LoadingState.tsx — `text-slate-500` → `text-muted`. Keep `role="status"` and the "Carregando…" copy EXACTLY.

ConflictTag.tsx — `bg-amber-100` → `bg-amber-soft`; `text-amber-800` → `text-amber`; `rounded` → `rounded-pill`
(paper+green uses pill radii for tags). Keep "divergente" copy, the `title={detail || undefined}` logic, and
`inline-flex items-center px-2 py-0.5 text-xs font-medium` EXACTLY.

FORBIDDEN palette tokens anywhere in these 6 files after your edit: `slate-`, `red-`, `amber-100`, `amber-800`,
`bg-white`, `#0F172A`, `blue-`. Use ONLY semantic utilities: text-ink, text-muted, text-faint, bg-surface, bg-surface-2,
bg-accent-soft, text-accent-ink, text-warn, bg-warn-soft, text-amber, bg-amber-soft, rounded-pill, rounded-control, font-mono.

────────────────────────────────────────────────────────────────────────
TEST EXTENSIONS (extend the two existing files — REAL assertions, no theater; every current assertion stays and stays green):

FactStates.test.tsx — ADD:
- UnknownValue: assert the "—" span has class `text-faint` (neutral-unknown token present; not green/accent). Keep the
  existing title/hint assertions unchanged.
- ConflictTag: assert the "divergente" tag has class `bg-amber-soft` AND `text-amber` (the retheme landed, not raw
  amber-100). The existing `className.toContain("amber")` assertion still holds — keep it.
- FreshnessIndicator: keep the existing `dados de <time>` and `dados de desconhecido` copy assertions UNCHANGED (they
  prove the copy/behavior was NOT altered). Optionally assert the wrapper has `text-muted`.

LoadStates.test.tsx — ADD:
- EmptyState: assert the outer container carries `text-muted` (and, when hint provided, the hint <p> carries `text-faint`).
  Keep the <p>-count and copy assertions unchanged.
- ErrorState: assert the `role="alert"` element (or its container) carries `text-warn`. Keep the retry-behavior and
  copy assertions unchanged.
- LoadingState: assert the `role="status"` element carries `text-muted`. Keep the copy assertion unchanged.

(Use `toHaveClass(...)` for the token assertions. These lock the retheme so a future regression to raw palette fails CI.)

────────────────────────────────────────────────────────────────────────
INTEGRITY / HARD CONSTRAINTS:
- Honest-unknown (ADR-17): UnknownValue "—" and FreshnessIndicator "desconhecido" must stay neutral/muted — never an
  accent/green/healthy token, never rendered as 0 or a real value.
- Do NOT change any copy string, prop, role, aria-label, or title logic. Retheme = class tokens only (+ the ConflictTag
  radius). Structure byte-identical otherwise.
- Do NOT touch packages/web-query (formatAsOf frozen), Button.tsx, index.ts barrel (F03-S6 owns it), or any file outside
  the write_set. Forbidden regardless: apps/server_core/**, packages/sdk-runtime/**, contracts/**, migrations, OpenAPI,
  apps/web/src/pages/**.
- No new npm dependency. No new import in any of the 6 primitives.
- Never read/print/commit any .env* file.

validation_kind: render test / typecheck / grep-assertion.

commands (run from worktree root; SANDBOX NOTE in bindings applies to test+build — record could-not-run(sandbox) and STOP,
do not burn the fixup; chip re-runs as verification of record):
- `npm run test -w @marketplace-central/web -- src/FactStates.test.tsx src/LoadStates.test.tsx`
   (note: the web vitest config globs packages/ui tests; if a relative path fails, run the whole `packages/ui` glob)
- `npx tsc --noEmit -p apps/web/tsconfig.json`   (expect ONLY the known baseline errors; any NEW error in these 6/2 files is yours)
- `rg -n 'slate-|red-|amber-100|amber-800|bg-white|#0F172A|blue-' packages/ui/src/UnknownValue.tsx packages/ui/src/FreshnessIndicator.tsx packages/ui/src/EmptyState.tsx packages/ui/src/ErrorState.tsx packages/ui/src/LoadingState.tsx packages/ui/src/ConflictTag.tsx`  → EXPECT NO MATCHES
- `rg -n 'text-faint|text-muted|text-warn|bg-amber-soft|text-amber|rounded-pill' packages/ui/src/{UnknownValue,FreshnessIndicator,EmptyState,ErrorState,LoadingState,ConflictTag}.tsx`  (semantic tokens landed)
- `rg -n 'formatAsOf|dados de|desconhecido' packages/ui/src/FreshnessIndicator.tsx`  (behavior preserved — formatAsOf call + no copy change)

expected_artifacts: 6 primitives using ONLY semantic tokens (zero dead-palette refs); honest-unknown stays neutral;
ConflictTag on bg-amber-soft/text-amber + rounded-pill; FreshnessIndicator copy/behavior unchanged (formatAsOf frozen);
both test files extended with real token assertions and ALL prior assertions still green; no barrel/Button/web-query
touch; no new tsc errors vs baseline.

open_questions: [] (dispatch-ready)

Report back: status; changed paths vs write_set; each command with evidence type (ran/assumed/could-not-run) + captured
output; confirm no copy/behavior/formatAsOf/barrel/Button change and honest-unknown stays neutral; anything you did NOT verify.
