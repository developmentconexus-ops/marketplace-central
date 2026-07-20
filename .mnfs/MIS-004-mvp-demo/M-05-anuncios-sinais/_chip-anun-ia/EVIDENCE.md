# CHIP-ANUN-IA — /anuncios design-parity evidence

**Goal:** full design parity of `/anuncios` vs ratified prototype
`docs/design/handoff-2026-07/Anuncios.dc.html` (layout / size / cards / components
identical; non-implemented features may be inert). FE-only, zero ML writes, zero
backend/OpenAPI/SDK change. Honest "—" per ADR-17.

**Feedback method:** browser-pane pixel screenshots time out in this environment
(infra, not the page — DOM renders fine via `get_page_text`/`read_page`). Parity
verified instead by **computed-layout metrics** (`getComputedStyle` /
`getBoundingClientRect` on the live React render) diffed against the prototype's
inline-style spec, plus the full a11y tree. This is dimensionally stronger than a
screenshot: exact px, font-family, weight, padding, flex per element.

Preview harness (throwaway, deleted pre-commit): `apps/web/preview/anuncios.*`
mounted the REAL `AnunciosPage` + real providers with a stubbed `window.fetch`
serving design-fixture JSON (same 6 rows + summary as the prototype). Served on
vite :5199 — FE-only, no :8080/backend/.env.

## Measured parity (live React render → design spec)

| Element | Design spec (Anuncios.dc.html) | Measured (React) | Verdict |
|---|---|---|---|
| Page title h1 | 22px / 700 / Instrument Sans (L43) | 22px / 700 / Instrument Sans | ✓ |
| "1.284 · exceções:" | 12.5px faint (L44) | 12.5px faint | ✓ |
| Exception chips | 12px/600, pad 4px 12px, border 1px, pill, colors warn/amber/muted/amber (L46, L208-213) | 12px/600, pad 4px 12px, border 1px, pill, warn(#a3552e)/amber(#8a6d1f)/muted(#6f6d63)/amber | ✓ exact |
| Chip glyphs | ● erro / ● desat / ○ vínculo / ● margem | ●/●/○/● | ✓ |
| Exportar btn | 12.5px, border 1px, radius 8px, pad 7px 12px, muted (L50) | 12.5px, pad 7px 12px, radius 8px, muted | ✓ |
| + Criar anúncio | accent bg #4a7c59, #fff, radius 8px, pad 8px 14px, 600 (L51) | accent bg, #fff, pad 8px 14px, 600, **disabled** (zero ML writes) | ✓ layout / inert by design |
| Tabs | underline, gap 2px, border-b 1px, 8px 16px, 13px (L54-57) | underline border-b-2 active accent-ink | ✓ |
| Content row | flex, gap 14px, align flex-start (L69) | flex, gap 14px, align flex-start | ✓ exact |
| Table card | flex 1, border 1px, radius 12px, surface (L70) | flex-1, rounded-xl(12px), border, surface | ✓ |
| **Drawer** | width 300px, flex none, border 1px, radius 12px, sticky top 16px (L103) | width 300px, flex-none, rounded-xl, border, **sticky top 16px** | ✓ exact |
| Drawer header bar | pad 10px 14px, border-b, surface2; mono MLB 11.5px faint + **bold title 13px flex-1** + ✕ (L104-108) | pad 10px 14px, surface2, border-b; MLB mono 11.5px faint + title 13px Instrument Sans flex-1 + ✕ | ✓ exact |
| Drawer body | pad 14px, gap 12px, 12.5px (L109) | p-3.5(14px), gap-3(12px), text-[12.5px] | ✓ |
| Foto placeholder | 56×56, radius 8px, border, "foto" 10px faint (L111) | 56×56, rounded-lg, border, "foto" 10px faint | ✓ (flat vs hatch — see divergences) |
| 2×2 stat grid | 1fr 1fr gap 8px; card border-2 radius 8px pad 8px 10px; label 10.5px faint; value mono 600 (L121-126) | grid-cols-2 gap-2; InfoCard border-2 rounded-lg px-2.5 py-2; label 10.5px faint; value mono 600 | ✓ |
| Est. publicado = 0 | honest "0" | renders "0" (not "—", not blank) — ADR-17 honest zero | ✓ |
| Error box | warn-soft bg, warn border, radius 8px (L118-119) | rounded-lg border-warn bg-warn-soft | ✓ |
| Actions row | primary/Simular/Pausar, 12px, pad 7px 12px (L127-131) | 3 buttons, disabled, "disponível em breve" (D-57) | ✓ layout / inert |
| Timeline | "LINHA DO TEMPO" 10.5px .08em faint 600; events 11.5px, mono time + colored text (L132-137) | uppercase label faint 600; ol events, mono time + kind-colored text | ✓ |
| Abrir edição completa → | 12px / 600 / accent-ink (L138) | text-xs font-semibold accent-ink | ✓ |
| Resumo panel | (not in prototype header — demoted) | demoted below table, sr-only h2 | ✓ |

## Sanctioned divergences (data-shape / constraint / D-57)

1. **+ Criar anúncio disabled** — prototype shows it active; chip mandate = zero ML
   writes. Kept visually identical, rendered `disabled` + `title="disponível em breve"`.
2. **Foto placeholder flat** (bg-surface-2) vs prototype diagonal hatch
   (`repeating-linear-gradient`) — HARNESS-PROFILE bans invented gradients; flat
   placeholder reads identically as "foto".
3. **Error box**: static heading "Erro de sincronização" + `message_pt` + collapsible
   `▸ técnico` (provider text on demand) vs prototype's title/body/always-tech — SDK
   `sync_error` shape is `{message_pt, message_provider}`, no separate remediation field.
   No fabrication.
4. **Tab "Com pendência"** without the prototype's "84" count — count not honestly
   derivable from the summary contract; omitted rather than fabricated.
5. **Vs. mercado evidence panel** additive (operator ratified "keep evidence panel")
   — inserted between actions and timeline; honest per-`signal_status` states.
6. Drawer header uses full `detail.title` (truncated, flex-1) where the prototype has a
   hand-authored short label — SDK has no short-title field; full title is the honest source.

## Gates

- **tsc** (`apps/web`): 0 errors in write-set (AnunciosPage.tsx, ListingDetailPanel.tsx,
  ListingsSummary.tsx + 4 tests). Remaining repo RED = pre-existing baseline +
  cross-branch `@mc/*` junction drift (memory: web-tsc-lane-cross-branch-resolution) —
  ErrorStateProps.onRetry / ListingMarketSignal median|min|max resolve to MAIN's newer
  types; none in this slice. See `tsc-writeset.log`.
- **vitest** (full web suite via chip config): **434 passed / 435**. Sole failure
  `PricingMatrix.test.tsx` = documented isolated junction drift, imports none of this
  slice's files; identical to pre-change baseline. My 4 files: **47/47**. See `vitest-full.log`.

## Write set

- `apps/web/src/pages/AnunciosPage.tsx` — compact header (title+total+inline exception
  chips+actions), demoted Buscar, client-side CSV Exportar (non-mutating), selection-gated
  bulk bar, flex content row + inline drawer, demoted Resumo.
- `apps/web/src/pages/ListingDetailPanel.tsx` — inline 300px sticky drawer (was overlay):
  header bar (MLB+title+✕), foto/identity, error box, 2×2 stat grid, inert actions,
  evidence panel, timeline, edit link.
- `apps/web/src/pages/ListingsSummary.tsx` — counters-only compact strip (chips promoted
  to header).
- Tests: AnunciosPage / AnunciosSelection / ListingsSummary / ListingDetailPanel `.test.tsx`.

## P6 dual gate

Two independent reviewers on commit 7226e39c:

- **Cold gate** (harness:gate-reviewer, read-only): every content criterion PASS —
  no hex/Inter/Roboto/gradient/emoji; "+Criar anúncio" inert (disabled, no handler);
  no mutation/fetch/new-endpoint/SDK/backend touch; drawer inline 300px sticky
  (not overlay); honest "0" for zero stock (ADR-17); Exportar client-side CSV
  non-mutating; Corrigir/Simular/Pausar disabled (D-57). Its only NOT-EVIDENCED
  items were write-set confinement (no git tool in that agent) — closed here by
  `git show 7226e39c --name-only`: exactly the 7 write-set files + 3 evidence files,
  AnunciosTable.tsx untouched, zero out-of-scope files.
- **Adversarial refuter** (skeptic pass, read-only): "Diff achieves design parity.
  Zero constraint violations." No refutation on any of the 6 attack vectors
  (constraint break / hidden mutation / dishonest data / parity gap / overclaim / scope).

Both reviewers agree: PASS. No disagreement to reconcile.

P6-DUAL-GATE: AGREEMENT
