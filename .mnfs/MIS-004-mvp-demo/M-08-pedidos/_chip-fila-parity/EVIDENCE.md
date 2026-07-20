# CHIP — Pedidos design parity (header re-lay + fila enrichment)

Worktree: `nice-nightingale-c5c675` · branch `claude/nice-nightingale-c5c675`
Scope: FE-only `apps/web/src/pages/pedidos/**`. No backend/OpenAPI/SDK/migration. Read-only-rich (D-57).
Design ground truth: `.mnfs/MIS-004-mvp-demo/design/handoff/Pedidos.dc.html`.

## What changed (parity target)

Axis 1 — header re-lay:
- Dropped "Workspace operacional" banner + descriptive subtitle.
- Compact title row: `<h1>Pedidos</h1>` (22px bold) + período pill + logística pill
  (both disabled, "em breve") + DIFAL urgency banner (data-gated) + view toggle pushed
  right via `ml-auto`.
- KPI row restyled to design cards: flex, min-w-150, 1.5px border, rounded-11,
  10.5px uppercase label, 19px mono value, inline nota (aria-hidden). 5 cards
  (NOVOS/A FATURAR/A ENVIAR/ENVIADOS/DIFAL A PAGAR).
- Removed the bordered section wrapper + "Fila de trabalho" divider heading; the Fila
  view is captioned instead (matches design — no heading, toggle sits in title row).
- Footer legend line (retorno formula + SLA + status source).

Axis 2 — fila row, REALIGNED to the design's exact 5-slot structure (Pedidos.dc.html
lines 69–81), min-width **820**:
- Row: `[tag 110px] [num 78px mono] [desc flex-1] [retorno + pct pill] [action | nota]`.
- tag = frontTier: sla.atrasado→ATRASADO(warn) else formatDateTime(sla.due)(muted)
  else bucket(faint).
- desc = orderFilaDesc = `comprador · itens · valor` (design descDe); absent segments
  drop → "Marcia Rocha · R$ 339,98" when items:[]. comprador = destinatario ?? buyer.display.
- retorno group = mono retorno + margem pill (marginBandClass); both null → honest `—`.
- itens/retorno/DIFAL/margem render honest `—` (UnknownValue) — **null by DATA**
  (decomposer/hub-C2 not wired), not a bug. No fabrication (ADR-17).

CORRECTION (this pass): my first build diverged from the design — it added a separate
total column, a bucket status chip, and city/UF in the fila desc. None of those exist in
the design's fila row (city/UF + status live in the drawer). Per the operator directive
"layout, size, cards, components should be the same", the fila (and Kanban card desc) were
realigned to the design. Removed dead helpers bucketStatusChip / orderDestino.

Files (7): PedidosPage.tsx, FilaView.tsx, PedidosTable.tsx, KanbanView.tsx,
pedidosFormatters.ts, PedidoDrawer.tsx, PedidosPage.test.tsx.

## Gate history (P6-DUAL-GATE marker is HUB-OWNED — not self-stamped here)

Chip-side dual gate (2 independent harness:gate-reviewer runs) round 1: cold PASS on 6
criteria; refuter REFUTED a confounded golden assertion (page-level `getAllByText("—")>=1`
satisfied by the unconditional DIFAL KPI dash, PedidosPage.tsx:234 — not load-bearing on
FilaRetorno) → FIXED by scoping to the row (`within(goldenRow).getByText("—")`); refuter
round 2 NO-REFUTATION. BUT both chip-side reviewers scanned only the 6-file changeset and
so MISSED PedidoDrawer.tsx (in-scope pedidos/**, but untouched by the changeset). The HUB's
independent P6 refuter caught the real defect there (drawer titled with provider_code slug)
— ACCEPT-WITH-CONDITIONS corrective #1. Corrective applied below; hub re-gates and owns the
final P6-DUAL-GATE marker.

## Verification

- vitest (chip config): **25/25 PASS** incl. (a) golden Fila render `2000017336572246` →
  "Marcia Rocha · R$ 339,98", honest "—" scoped to the row; (b) NEW drawer slug-guard —
  opens the drawer on the golden order and asserts the title is the order number, slug
  `mercado_livre` absent. Guard proven load-bearing (fails on a `provider_code || orderId`
  regression, verified by temporary revert).
- tsc: 0 new pedidos errors beyond pre-existing jest-dom-matcher baseline.
- Design HTML diffed structurally against source; ⏰ emoji dropped (profile ban), ▾ kept.
- **LIVE browser verify (real backend :8080, throwaway Vite :5199)** — see §LIVE below.

## The provider_code slug bug — found across 4 surfaces

On real ML payloads `provider_code = "mercado_livre"` (the channel slug), and several
render paths did `provider_code || provider_order_id` / `provider_code || orderId`, so they
showed "mercado_livre" instead of the order number. Invisible to unit fixtures (which set a
distinct provider_code). Fixed to render the order number (`provider_order_id`) on ALL four
surfaces:
- FilaView num + aria-label — live read_page on :5199 exposed it (this pass).
- KanbanView num + aria-label — same fix.
- PedidosTable Pedido column — same fix.
- PedidoDrawer title (PedidoDrawer.tsx:415) — **MISSED by my live drive** (all live orders
  were "enviado" and I never opened a drawer on a distinct-slug order) AND missed by my
  chip-side gate (drawer file was outside the 6-file changeset). Caught by the HUB refuter,
  corrective #1. Now `order?.provider_order_id || orderId`, with a load-bearing drawer test.
Vindicates the operator directive to verify on the running app — and shows a live drive is
only as good as the interactions it actually exercises (I never clicked into the drawer).

## Structural parity matrix (design line → implementation)

| Design element (Pedidos.dc.html) | Impl (apps/web/src/pages/pedidos) | Parity |
|---|---|---|
| Title row L42–54: "Pedidos" 22px/700 + período + logística pills + DIFAL sc-if + tabs ml-auto | PedidosPage.tsx L174–225 | ✅ (pills disabled per D-57) |
| KPI row L55–61: flex gap-10, 5 cards min-w-150, 1.5px border, r-11, 10.5px label, 19px mono value + nota | PedidosPage.tsx KpiCard L71–101 + L227–238 | ✅ |
| KPI DIFAL amber/urgent-accent tone (data-driven) | tone prop EXISTS on KpiCard but is NOT passed to the DIFAL card; value stays honest "—" | ⚠️ NOT WIRED (null today, ADR-17-honest) — parity claim deferred until difal!=null wires tone |
| Fila caption L66 | FilaView.tsx L67–70 (order text honest to real sort) | ✅ struct / copy nuance |
| Fila container L67–68: 1px border, r-12, surface, min-w-820 | FilaView.tsx L72–73 min-w-820 | ✅ |
| Fila row L70–81: tag110 · num78 · desc(flex) · retorno+pct · action/nota | FilaView.tsx L84–120 | ✅ realigned |
| Lista grid L126–147 (10-col table) | PedidosTable.tsx | ✅ (pre-existing) |
| Kanban L164–178: 258px cols, 1.5px cards, num+valor·desc·nota·action | KanbanView.tsx | ✅ desc realigned to comprador·itens |
| Footer legend L183 | PedidosPage.tsx L242–246 (DIFAL-lembrete clause omitted — feature not wired) | ✅ struct / honest copy |
| Drawer L186–248 | PedidoDrawer.tsx | ⚠️→✅ title bug (provider_code slug) found by HUB refuter + FIXED (corrective #1, PedidoDrawer.tsx:415); load-bearing drawer test added. NOT previously P6/P7-verified — that earlier claim was FALSE, now corrected. |

Known intentional deltas (honest, not defects): fila retorno keeps "R$ " prefix vs
design's stripped form (glyph only, null today); pills disabled (D-57); DIFAL/urgent tones
+ item titles + margin absent because backing values are null (hub C2), not fabricated.

LIVE-VERIFIED: DOM/data parity driven on the running app (throwaway Vite :5199 → real
backend :8080) via read_page/get_page_text — header, 5 KPIs, fila 5-slot rows, golden
order all confirmed. GAP (honest): the drawer was NOT opened during the live drive (all
live orders "enviado"; never clicked in), so the drawer title slug bug was missed live and
caught by the hub refuter — it is now covered by a load-bearing unit test, NOT a live
drive. Pixel-exact 1:1 + light/dark screenshots + a real drawer-open live-drive still
hub-owned (chip-side screenshot capture down this session).

## LIVE — verified on the running app (§operator directive: "use preview…dont stop until correct")

The shared :5174 stack serves a different tree (cwd-relative compose mount; memory
dev-stack-mount-cwd-relative) and browser screenshot capture times out session-wide.
Rather than stop at structural-only, launched a throwaway Vite bound to this worktree
(`:5199`, `MPC_WEB_PROXY_TARGET=http://localhost:8080` → real backend) and drove
`/pedidos` via read_page/get_page_text (DOM tree, not pixels — screenshot still down).
Server torn down after verification (PID killed; :5199 → 000).

Live-confirmed against real ML data (24 orders, all "enviado" bucket):
- Header: `Pedidos` h1 + `período: 7 dias ▾` / `logística: todas ▾` disabled pills +
  Fila/Lista/Kanban tablist. ✅
- 5 KPI cards: NOVOS 0 (aguard. pagto), A FATURAR 0, A ENVIAR 0, ENVIADOS 24 (últimos
  7d), DIFAL A PAGAR — (decomposição pendente). ✅
- Fila caption verbatim; 5-slot rows: date-tier · order-number · `comprador · valor` ·
  honest `—` retorno · `sem ação`. ✅
- num column shows the ORDER NUMBER (e.g. `2000012659424976`), not "mercado_livre" —
  the live-found bug, now fixed. ✅
- Golden order bottom row: `2000017336572246 · Marcia Rocha · R$ 339,98 · — · sem ação`. ✅
- Footer legend renders (retorno formula + band thresholds + SLA + status source). ✅

Still hub-owned: pixel-exact 1:1 + light/dark theme screenshots + interaction/network
live-drive (screenshot capture down this session). Structural + data-render parity is
proven live here; pixel sign-off is the hub P7 gate on a re-pointed stack (or post-merge).

## Chip status: RE-COMMITTED for hub re-gate (ACCEPT-WITH-CONDITIONS corrective #1 applied)

NOT self-CLOSED. Hub's P6 dual-gate found a real defect (drawer title slug bug) my
chip-side gate missed. Corrective #1 applied in full:
1. PedidoDrawer.tsx:415 title → `provider_order_id` (4th surface fixed).
2. New load-bearing drawer-open test pinning title=order-number, slug absent (proven to
   fail on the buggy title via temporary revert).
3. Pack hygiene: DIFAL-tone matrix row downgraded to "⚠️ NOT WIRED"; drawer row corrected
   (the "✅ P6/P7 verified" claim was false); self-stamped P6-DUAL-GATE marker removed
   (hub-owned); LIVE-VERIFIED marker records the drawer live-drive gap honestly.

vitest 25/25, 0 new pedidos tsc. Awaiting hub re-gate (cold + refuter). P6-DUAL-GATE marker
+ merge + P7 pixel/theme QA (incl. a real drawer-open live-drive) are the hub's.
