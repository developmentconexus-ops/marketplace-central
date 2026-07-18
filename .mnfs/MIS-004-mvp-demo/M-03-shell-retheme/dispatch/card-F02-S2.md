SLICE CARD — F02-S2 · feature F-02-header-nav-routes · milestone M-03

depends_on: [F02-S1] (route modules + AppRouter indirection already committed 25b02c8), [F01-S2] (theme hook committed).

goal: Kill the dark sidebar and install the paper+green HORIZONTAL header. Extract a new `Header` component and make
`Layout` render it above the main content. Reuse `InstallationContext` (installation selector) and the F-01 `useTheme`
hook (theme control) VERBATIM — do NOT add a second theme state, provider, or storage path. This is the demo-critical
shell chrome. No route/URL changes (F02-S1 owns routing); no shared-primitive work (F-03 owns that).

complexity: standard.

write_set (create Header + its test, gut Layout):
- apps/web/src/app/Header.tsx        (new — export `Header`)
- apps/web/src/app/Header.test.tsx   (new — render tests, jsdom)
- apps/web/src/app/Layout.tsx        (edit — remove <aside> sidebar + old <header>; render <Header/> then <main>)

────────────────────────────────────────────────────────────────────────
STYLING — semantic Tailwind v4 utilities ONLY (they resolve to the index.css @theme CSS vars; light/dark auto-swap via
`data-theme`). NEVER hardcode a hex, NEVER reference #0F172A / slate-* / blue-* / white / bg-slate-900 — those are the
DEAD sidebar palette. Available utilities (already emitted by @theme): bg-bg, bg-surface, bg-surface-2, text-ink,
text-muted, text-faint, border-border, border-border-2, bg-accent, bg-accent-soft, text-accent-ink, rounded-control,
rounded-card, rounded-pill, font-mono. Use `border` + `border-border`, etc.

Element palette (compose from the above — do not invent classes beyond spacing/flex/text-size layout utilities):
- Header bar: `bg-surface border-b border-border`, horizontal flex, `h-14`, `px-4 lg:px-6`, items centered, `gap-4`.
  MUST be horizontally scrollable at narrow width: wrap the pill row (and if needed the whole bar) so every pill/control
  stays reachable — e.g. `overflow-x-auto` on the nav row; nothing may become unreachable when the viewport narrows.
- Brand wordmark (left): `text-ink font-semibold text-sm tracking-wide`, text "Marketplace Central". Single brand only
  (the old sidebar title + header h1 duplication dies).
- Enabled pill (inactive): `rounded-pill px-3 py-1.5 text-sm font-medium text-muted hover:bg-surface-2 hover:text-ink transition-colors`.
- Enabled pill (active): `bg-accent-soft text-accent-ink` (same shape). Active state MUST derive from router — use
  `NavLink` `className={({isActive}) => …}`. Enabled pills are `NavLink` to `{ pathname, search: location.search }`
  so the installation query param is preserved (same behavior the old sidebar had).
- Disabled "em breve" pill (Mercado, Repasses): NOT a link and NOT navigable — render a `<span>` (or `<button type="button" disabled>`),
  `text-faint cursor-not-allowed`, with a visible inline badge reading `em breve`
  (`text-[10px] uppercase tracking-wide rounded-pill bg-surface-2 text-faint px-1.5 py-0.5`). Clicking it MUST NOT change
  the URL and MUST NOT be a router link. This is the honest "not built yet" affordance — do not fake it as active/clickable.

────────────────────────────────────────────────────────────────────────
NAV PILLS — EXACTLY these six, in THIS order (canonical IC-05 nav; NO Vínculos, NO Catálogo/Estoque/Integrações in the
pill row — those live only in the gear menu):
1. `Visão geral`  → NavLink to "/"          (end)      [enabled]
2. `Anúncios`     → NavLink to "/anuncios"             [enabled]
3. `Mercado`      → DISABLED, em breve, non-nav         [disabled]
4. `Simulador`    → NavLink to "/precos"               [enabled]   ← deep-open /precos ⇒ this pill active
5. `Pedidos`      → NavLink to "/pedidos"              [enabled]
6. `Repasses`     → DISABLED, em breve, non-nav         [disabled]

RIGHT CLUSTER (in the header, right-aligned via `ml-auto`, `gap-2`/`gap-3`):
a. Theme control — a `<button type="button">` that calls `toggleTheme` from `useTheme()`. Label/aria reflects current
   theme (e.g. aria-label `Alternar tema` ; may show a sun/moon glyph or the text of the target theme). Use ONLY
   `useTheme` from "./theme/useTheme" — do NOT add another useState for theme, do NOT re-read/write localStorage here.
b. Installation selector — REUSE the existing logic verbatim from the current Layout: `useInstallation()` →
   `{ installationId, setInstallationId, installations, status }`, the same three status branches
   (ready+selected ⇒ the `ML:` <select> of installations by `installation_id`/`display_name`;
   `empty` ⇒ the "Conecte uma conta em Integrações" hint; otherwise a neutral spacer). Retheme its wrapper to semantic
   tokens (`border-border bg-surface text-ink`, `rounded-pill`, focus ring via `focus-within:border-accent`) — same
   behavior, new skin. Keep `aria-label="Selecionar instalação"`.
c. Gear menu — a ⚙ control opening a dropdown. Prefer native `<details><summary>` (accessible, no click-outside JS,
   testable) OR a controlled dropdown if you add proper click-outside handling; either is fine. The menu has EXACTLY
   four top-level entries/groups, in this order, and NOTHING else (NO Vínculos anywhere):
     1. `Configurações` — a group header exposing ONE nested item:
          `DIFAL` → a `Link` to the FIXED string "/precos?params=1". It MUST be exactly that literal —
          do NOT merge `location.search`, do NOT inherit unrelated params.
     2. `Integrações` → Link to "/integracoes"
     3. `Catálogo`    → Link to "/catalogo"
     4. `Estoque`     → Link to "/estoque"
   Menu surface: `bg-surface border border-border rounded-card shadow`, items `text-ink hover:bg-surface-2 rounded-control px-3 py-2`.

LAYOUT.tsx after edit — the WHOLE `<aside>` sidebar block and the old inline `<header>` are DELETED. Remove the now-unused
lucide icon imports and the `navItems` array (they belonged to the sidebar). New structure:
```
export function Layout() {
  return (
    <div className="flex min-h-screen flex-col bg-bg text-ink">
      <Header />
      <main className="flex-1 overflow-auto p-4 lg:p-6">
        <InstallationGate>
          <Outlet />
        </InstallationGate>
      </main>
    </div>
  );
}
```
(Keep `Outlet` from react-router-dom and `InstallationGate` from "./InstallationContext". The installation SELECT logic
moves INTO Header — Layout no longer calls useInstallation directly. `InstallationGate` still wraps `<Outlet/>` in main.)

────────────────────────────────────────────────────────────────────────
INTEGRITY / HARD CONSTRAINTS:
- Mercado and Repasses are NEVER navigable and NEVER rendered as active/enabled. They are honest "em breve" stubs. A
  disabled pill must NOT be a NavLink/Link and must NOT register or hit any route.
- Do NOT touch routing: no changes to AppRouter.tsx, no route path strings, no new routes, no catch-all. The disabled
  pills have NO route (that is correct — they are not wired).
- Do NOT add a theme provider/context or a second theme state — ONLY `useTheme()`. One theme source of truth (F01-S2).
- Reuse InstallationContext as-is; do NOT change InstallationContext.tsx or its API.
- Forbidden paths (regardless): apps/web/src/pages/**, apps/server_core/**, contracts/**, migrations, OpenAPI, any
  packages/**. FE app-shell only. No new npm dependency (if you think you need one, STOP and report — do not install).
- Never read/print/commit any .env* file.

────────────────────────────────────────────────────────────────────────
Header.test.tsx (jsdom render tests — REAL assertions, no test theater). Render `<Header/>` inside a
`<MemoryRouter>` (and the `InstallationProvider` if the selector needs context — if `useInstallation` throws outside a
provider, wrap in `InstallationProvider`; installations may be empty, that is fine, assert the pills regardless). Cover:
1. The six pill labels render in exact order: Visão geral, Anúncios, Mercado, Simulador, Pedidos, Repasses.
2. Mercado and Repasses show visible `em breve` text AND are NOT links — assert they have no `href`/`role="link"`
   (e.g. `screen.getByText("Mercado").closest("a")` is null), i.e. they cannot navigate.
3. Enabled pills (Visão geral, Anúncios, Simulador, Pedidos) ARE links (anchor present).
4. Active pill derives from router: render inside `<MemoryRouter initialEntries={["/precos"]}>` and assert the
   `Simulador` pill has the active classes (`bg-accent-soft`/`text-accent-ink`) while e.g. `Anúncios` does not. Also a
   second case: initialEntries `["/"]` ⇒ `Visão geral` active.
5. Gear menu (open it if it uses <details>/state) exposes exactly the four entries Configurações, Integrações, Catálogo,
   Estoque and the nested DIFAL; assert DIFAL's anchor href is exactly "/precos?params=1"; assert NO "Vínculos" text
   appears anywhere in the header.
6. Theme button present (query by its aria-label/role) and clicking it does not throw (it calls toggleTheme). You do NOT
   need to assert the data-theme flip (that is F01-S2's tested surface) — just that the control exists and is wired to
   useTheme, not a local state.

validation_kind: render test / typecheck / grep-assertion / manual-QA-drive (chip drives the browser QA at close).

commands (run from worktree root; SANDBOX NOTE in bindings applies to test+build — record could-not-run(sandbox) and
STOP, do not burn the fixup; chip re-runs as verification of record):
- `npm run test -w @marketplace-central/web -- src/app/Header.test.tsx`
- `npx tsc --noEmit -p apps/web/tsconfig.json`   (expect ONLY the known pre-existing baseline errors — 162 TS2339
   jsdom-matcher + AppRouter/precos Client mismatches; ANY new error in Header.tsx/Layout.tsx is yours to fix)
- `rg -n '#0F172A|bg-slate-9|border-slate-7|text-slate-|bg-blue-|<aside|Sidebar|navItems' apps/web/src/app/Layout.tsx apps/web/src/app/Header.tsx`  → EXPECT NO MATCHES (dead sidebar palette fully gone)
- `rg -n 'Visão geral|Anúncios|Mercado|Simulador|Pedidos|Repasses|Configurações|DIFAL|Integrações|Catálogo|Estoque' apps/web/src/app/Header.tsx`  (all present)
- `rg -n 'Vínculos' apps/web/src/app/Header.tsx apps/web/src/app/Layout.tsx`  → EXPECT NO MATCHES
- `rg -n '/precos\?params=1' apps/web/src/app/Header.tsx`  (DIFAL fixed target present)
- `rg -n 'useTheme' apps/web/src/app/Header.tsx`  (theme control wired to the F-01 hook)

expected_artifacts: horizontal paper+green Header with the exact six pills in order (Mercado+Repasses honest non-nav
"em breve"); active pill from router incl. deep /precos ⇒ Simulador; theme button via useTheme (no 2nd state);
installation selector behavior preserved (installationId/setInstallationId/installations/status); gear menu with exactly
Configurações(▸DIFAL→/precos?params=1)/Integrações/Catálogo/Estoque and no Vínculos; Layout sidebar+old header removed,
renders <Header/> then <main><InstallationGate><Outlet/></InstallationGate></main>; header horizontally usable at narrow
width; only semantic tokens (zero dead-palette refs); green Header.test.tsx; no new tsc errors vs baseline.

open_questions: [] (dispatch-ready)

Report back: status; changed paths vs write_set; each command with evidence type (ran/assumed/could-not-run) + captured
output; confirm no routing/pages/**/InstallationContext/packages touched; confirm Mercado+Repasses are non-navigable and
no Vínculos anywhere; anything you did NOT verify.
