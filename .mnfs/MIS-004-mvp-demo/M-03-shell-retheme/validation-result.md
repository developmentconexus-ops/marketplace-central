# M-03 shell-retheme — P7 QA validation-result

- Milestone: MIS-004 M-03 (shell-retheme), FE-only.
- Branch: `chip/m03-shell-retheme` — close SHA **`df236aff`** (range `59d0e62..df236aff`, 47 files).
- QA persona: fresh cold live-drive against `validation-contract.md` C01–C06.
- Surface under test: `http://localhost:5174` served from this worktree (hub re-point, ledger D-17).
- Method: browser live-drive (`javascript_tool` DOM/computed-style + `navigate`). Screenshot
  rasterizer broken in this browser pane (finding **F-ENV-10** — 30s timeouts, console clean);
  visual evidence supplied as computed-style token captures per that limitation.

## Verdict: **PASS** — C01–C06 all satisfied. Zero blockers.

---

### C01 — tema (tokens light/dark, persist, fonts self-hosted) — PASS
Computed-style capture at `/`:

| token | light | dark |
|---|---|---|
| `data-theme` | `light` | `dark` |
| body bg | `rgb(251,250,247)` (paper) | `rgb(22,24,20)` |
| body text | `rgb(37,41,31)` (ink) | `rgb(233,232,226)` |
| header bg | `rgb(255,255,255)` | `rgb(31,34,28)` |

- Toggle "Alternar tema" flips `data-theme` light↔dark; tokens repaint accordingly.
- Persist: `localStorage["marketplace-central-theme"]="dark"` survives F5 (verified earlier drive);
  `apps/web/index.html` inline no-flash script re-applies `data-theme` pre-paint on reload.
- Fonts self-hosted (NO CDN): `document.fonts.check('16px "Instrument Sans"')=true`;
  `document.fonts.load('16px "IBM Plex Mono"')→check=true` (lazy on first glyph, no 404).
  Google Fonts `<link>` (Inter/JetBrains via googleapis) **removed** from index.html; replaced by
  `@fontsource/instrument-sans` + `@fontsource/ibm-plex-mono` bundled deps.

### C02 — pills (6 exact, disabled non-navigable, active via deep-link+F5) — PASS
`nav[aria-label="Navegação principal"]` children, in order:
1. Visão geral — `A` → `/` 2. Anúncios — `A` → `/anuncios`
3. **Mercado — `SPAN`, no href, "em breve"** 4. Simulador — `A` → `/precos`
5. Pedidos — `A` → `/pedidos` 6. **Repasses — `SPAN`, no href, "em breve"**
- Disabled Mercado: real `click` dispatched → `location.href` **unchanged** (`urlChanged:false`).
- Active state: at `/precos` only Simulador carries `aria-current="page"`; survives fresh nav
  (F5 equivalent) to `/precos?installation=…`.
- Enabled pills carry `location.search` (query preserved on navigation — Layout.test.tsx guards).

### C03 — gear menu (4 items, DIFAL→/precos?params=1, /vinculos registered) — PASS
Trigger = `<summary>` "⚙" (aria-label "Abrir menu"); disclosure "Menu de configurações" opens
**exactly 4** items:
| text | href |
|---|---|
| DIFAL | `/precos?params=1` |
| Integrações | `/integracoes` |
| Catálogo | `/catalogo` |
| Estoque | `/estoque` |
- No "Vínculos" in gear. Route `/vinculos` stays **registered**: direct nav renders
  "Em construção" placeholder with header intact (not 404).

### C04 — route indirection (no URL change) — PASS
`AppRouter.tsx` diff: every `path=` string byte-identical base↔close; only `element=` swapped to
per-area indirection components (`AnunciosRoute`/`ProdutoRoute`/`VinculosRoute`/`PrecosRoute`/
`PedidosRoute`). Legacy English paths still redirect via `LegacyRedirect`
(`/products→/catalogo`, `/orders→/pedidos`, `/simulator→/precos`, …). Zero canonical URL changed.

### C05 — primitivas (unknown ≠ zero) — PASS
- `UnknownValue` → `<span class="text-faint">—</span>` — dash, never `0`/`R$0`/green (ADR-17).
- `MarginChip` null → band=`unknown` → renders "—" (`bg-surface-2 text-muted`), no fabricated %.
- Barrel `packages/ui/src/index.ts` exports (additive, prior byte-intact): MarginChip (+type),
  DataTable (+DataTableProps/DataTableColumn), DetailDrawer (+DetailDrawerProps), UnknownValue.

### C06 — ownership (FE-only, scope guard) — PASS
- 47 changed files. **Zero** `.go`/`.sql`/`openapi`/`migration`/`.proto`.
- Out of `apps/web/src` + `packages/ui/src`: only `apps/web/index.html` (font/theme seam) and
  manifests. Dep deltas = ratified grants ONLY: `@types/node` (root),
  `@fontsource/instrument-sans` + `@fontsource/ibm-plex-mono` (apps/web). No rogue deps.

---

## Seams / defers / notes
- Dev surface served from worktree via hub re-point (ledger D-17); hub re-points frontend back to
  `main` before merge ladder (post-CLOSED).
- Backend `GET /catalog/products/{id}` 500 in xlsx mode = **F-QA-M01-1** (owned by M-06) — NOT an
  M-03 defect; out of this milestone's FE scope.
- F-ENV-10 (browser-pane rasterizer broken) — visual evidence delivered as computed-style token
  captures in lieu of screenshots.

## Ladder verification of record (post-reconciliation, SHA df236aff)
- vitest **215/215** (`--no-file-parallelism`, sequential clean per F-ENV-5).
- tsc **181 baseline errors, 0 new** (jest-dom TS2339 env-gap = F-ENV-4).
- P6 dual gate CLOSED: GATE-A (cold Opus) PASS-with-nits; GATE-B (Sol medium) reconciled on merit
  (0 surviving blockers); reconciliation commit df236aff independently reviewed PASS.
