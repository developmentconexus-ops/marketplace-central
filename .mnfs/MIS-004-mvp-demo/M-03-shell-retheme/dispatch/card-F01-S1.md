SLICE CARD — F01-S1 · feature F-01-theme-tokens-fonts · milestone M-03

SCOPE NOTE (read first): this slice delivers the token contract + Tailwind v4 `@theme` semantic
utilities + radii + font-family variables ONLY. It does NOT add `@font-face`, does NOT add font
binary assets, and does NOT edit apps/web/index.html. The self-hosted-font wiring and CDN removal
are a SEPARATE later slice (font-binary provenance is pending a hub decision). Do not touch fonts
beyond declaring the two font-family CSS variables below.

goal: Establish the complete papel+verde light/dark token contract as CSS variables on the root,
expose semantic Tailwind v4 utilities via an `@theme` block in index.css, define the design radii,
and point the font-family variables at the design families (with system fallbacks). No layout, no
React state, no component edits.

write_set (edit ONLY this file):
- apps/web/src/index.css

Current index.css (24 lines) defines `@import "tailwindcss";`, five `@source` directives, a box-sizing
reset, `:root { --font-sans: 'Inter'...; --font-mono: 'JetBrains Mono'... }`, and a `body` rule with
`background-color:#F8FAFC; color:#0F172A`. PRESERVE the `@import "tailwindcss";` line and ALL five
`@source` directives exactly. Replace the color/font specifics as below.

REQUIRED CONTENT:

1) Light tokens on `:root` — EXACT values (do not alter a single hex):
```
--bg:#fbfaf7; --surface:#ffffff; --surface2:#f4f2ea; --border:#e6e3da; --border2:#f0eee6;
--ink:#25291f; --muted:#6f6d63; --faint:#8b887c;
--accent:#4a7c59; --accent-soft:#edf3ec; --accent-ink:#33573f;
--warn:#a3552e; --warn-soft:#f7ebe2; --info:#2f5bb7; --info-soft:#e8eef9;
--amber:#8a6d1f; --amber-soft:#fdf3d7; --donut-track:#e9e6dc;
```

2) Dark overrides on `[data-theme="dark"]` — EXACT values:
```
--bg:#161814; --surface:#1f221c; --surface2:#262a23; --border:#2e312a; --border2:#2a2d26;
--ink:#e9e8e2; --muted:#a5a399; --faint:#8b887c;
--accent:#7fb08c; --accent-soft:#243328; --accent-ink:#9ecfab;
--warn:#d08a63; --warn-soft:#3a2a20; --info:#7fa3e0; --info-soft:#222b3c;
--amber:#d4b45a; --amber-soft:#37311d; --donut-track:#2e312a;
```
(These are the exact values in docs/design/handoff-2026-07/Dashboard.dc.html:14-15 — you may open that
file to self-verify before writing; report the comparison as your token receipt.)

3) Font-family CSS variables (replace the current Inter/JetBrains Mono ones), on `:root`:
```
--font-sans: 'Instrument Sans', ui-sans-serif, system-ui, sans-serif;
--font-mono: 'IBM Plex Mono', ui-monospace, SFMono-Regular, monospace;
```
(No `@font-face` yet — the families fall back to system fonts until the font slice lands. That is
expected and acceptable for this slice.)

4) Radii variables on `:root`: control/button/input `--radius-control: 8px`; card `--radius-card: 12px`;
pill/chip `--radius-pill: 999px`.

5) A Tailwind v4 `@theme { ... }` block mapping semantic utilities to the tokens, so components can use
classes that resolve to the vars. Expose at minimum (Tailwind v4 `--color-*` / `--font-*` / `--radius-*`
namespaces generate the corresponding `bg-*`/`text-*`/`border-*`/`font-*`/`rounded-*` utilities):
```
--color-bg, --color-surface, --color-surface-2, --color-border, --color-border-2,
--color-ink, --color-muted, --color-faint,
--color-accent, --color-accent-soft, --color-accent-ink,
--color-warn, --color-warn-soft, --color-info, --color-info-soft,
--color-amber, --color-amber-soft
--font-sans (→ var(--font-sans)), --font-mono (→ var(--font-mono))
--radius-control, --radius-card, --radius-pill
```
Each `@theme` entry references the corresponding `var(--token)` (e.g. `--color-surface: var(--surface);`)
so that dark-mode overrides on `[data-theme="dark"]` flow through automatically with no JS. Verify in the
build that `bg-surface`, `text-ink`, `border-border`, `text-accent-ink`, `bg-accent-soft` are generated.

6) `body`: `font-family: var(--font-sans); background-color: var(--bg); color: var(--ink);` keep the
`-webkit-font-smoothing: antialiased;`. Keep the box-sizing reset.

validation_kind: build / typecheck / grep-assertion (NO unit test applies to a pure declarative CSS
slice — do NOT fabricate a test; the build + grep ARE the verification).

commands (run from worktree root; capture output):
- `npx tsc --noEmit -p apps/web/tsconfig.json`  (must be clean — this slice adds no TS)
- `npm run build -w @marketplace-central/web`  (must succeed; proves @theme is valid and utilities generate)
- `rg -n -- '--bg:#fbfaf7|--accent:#4a7c59|--ink:#25291f' apps/web/src/index.css`  (light tokens present)
- `rg -n -- 'data-theme="dark"' apps/web/src/index.css` and `rg -n -- '--bg:#161814|--accent:#7fb08c' apps/web/src/index.css`  (dark tokens present)
- `rg -n -- '@theme' apps/web/src/index.css`  (theme block present)
- `rg -n -- 'Instrument Sans|IBM Plex Mono' apps/web/src/index.css`  (font families set)
- `rg -n -- '@font-face|fonts.googleapis' apps/web/src/index.css`  → EXPECT NO MATCHES in index.css (fonts deferred)

expected_artifacts: index.css with :root light tokens (exact), [data-theme="dark"] dark tokens (exact),
font-family + radii vars, and an `@theme` block mapping semantic utilities; clean tsc; green web build with
generated semantic utilities; token-comparison receipt vs the .dc.html.

complexity: standard

open_questions: [] (dispatch-ready)

Report back: status; changed paths vs write_set; each command with evidence type (ran/assumed/could-not-run)
and captured output; the token-comparison receipt; anything you did NOT verify.
