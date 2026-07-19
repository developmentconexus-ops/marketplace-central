# M-08-pedidos — Validation Result (P7 QA COLD STATIC)

```yaml
id: M-08
type: milestone-validation-result
validator: QA Validator (fresh/cold, static only)
worktree: chip-m08-pedidos
head_sha: 5051ae53
delta: 8b6c4b30..5051ae5
date: 2026-07-18
mode: READ-ONLY, no server/dev-stack/DB, no browser live-drive (stack hub-owned)
```

## Acceptance bar honored
- D-57 READ-ONLY RICH: all mutation controls render disabled ("em breve"/"disponível em breve") — CORRECT, not a defect.
- D-62 HYBRID slice C: decomposição/DIFAL/retorno are C1 = additive nullable, honest-nil today; real decomposer is C2, hub-owned, wired post-merge. This report verifies the SEAM and honesty, not live values.
- Stack is hub-owned: browser live-drive is could-not-run here and routed to hub.

---

## §projecao (C01)

VERDICT: **PASS-STATIC-hub-drive-pending**

Evidence (static, ran):
- `apps/server_core/internal/modules/orders/transport/http_handler.go:370-392` — `enrichedOrderDTO` is additive: embeds `domain.OrderReadModel` (old shape preserved) and adds `items` (overridden with cost), `vinculo_status`, `buyer`, `sla`, `destino_uf`, `rastreio`, `bucket`, `retorno_liquido`, `margem_pct`, `decomposicao`, `difal`. No old field removed/renamed.
- ML/provider raw DTO does not leak: `internal/composition/orders_adapters.go` — `ordersShipmentReaderAdapter.GetShipment` returns `connectorsdomain.ShipmentInfo` (module-owned normalized domain type), the raw ML shipment payload dies inside the pre-existing connector `CapabilityAdapter`, never reaching orders/transport. `grep -rn "modules/pricing" orders/` → no hits (see §decomposicao).
- Cost via `GetCostAsOf` with source time visible: `application/enrich_service.go:160-202` `resolveItemCosts`/`orderEffectiveAt` — the "as of" date is the first non-nil of `ProviderClosedAt`/`ProviderCreatedAt`/`CreatedAt` (never a zero time), passed to `CostReader.GetCostAsOf`.
- Unavailable ⇒ null never 0 (ADR-17): `enrich_service.go:168-183` — any resolution miss (`cost==nil`, no `InternalProductID`, no effective date, reader error, nil `Amount`) appends the identifier to `unknown` and leaves `ItemCost.UnitCost` nil; DTO field `custo_unitario` is `*float64,omitempty` (`http_handler.go:335`) — omitted, never `0`.
- Lanes: `go build`/`go vet`/`go test ./apps/server_core/internal/modules/orders/... ./apps/server_core/internal/composition/...` → build/vet clean, all 9 packages `ok`.

Evidence-type: static/ran. LIVE response bodies against a real ML installation = **could-not-run (stack hub-owned)** → ROUTED TO HUB.

---

## §decomposicao (C02)

VERDICT: **PASS (architecture half, chip-verifiable)** / value-parity = **DEFERRED-C2**

Evidence (ran):
- `grep -rn "modules/pricing" apps/server_core/internal/modules/orders/` → **zero hits**. Orders does not import the pricing module.
- `grep -rniE "difal.*=.*\*|margem.*calc|decompose\(" .../orders/domain .../orders/application` (excluding tests) → only `enrich_service.go:148: if p, ok := s.decomposer.Decompose(ctx, order); ok {` — a **port call**, not a formula. No arithmetic formula reimplementation found.
- `apps/server_core/internal/modules/orders/ports/decompose_reader.go:9-19` — `Decomposer` is a module-local, nil-able consumer port (`Decompose(ctx, order) (domain.OrderProfitability, bool)`), explicitly documented as mirroring the `CostReader`/`ShipmentReader` nil-port idiom; the real IC-04 pricing engine adapter is C2/hub-owned/post-merge.
- `application/enrich_service.go:144-151` `resolveProfitability` — nil decomposer or not-ok both degrade to `domain.UnknownOrderProfitability()` (honest-empty), never a panic, never a fabricated value.
- DIFAL is per-order, honest-nil until UF route is resolvable: `enrich_service.go:141-143` comment — C1 deliberately leaves `Difal.UFRoute` nil because a route needs both origin (seller UF, engine-owned fact) and destino; only C2 (which has tenant origin) can fill it.

Gate ordering check (contract: "trecho decomposição só commitado APÓS ports M-07 publicados"):
- Since the real M-07 decompose engine (`d288cd12`/merge `42e0b204`) is intentionally **not consumed** by this branch (D-62: orders defines its own local nil-able port, not an import of pricing's engine — see above), the literal "consume M-07's published ports" gate does not apply to C1. What was verified instead: orders never imports `modules/pricing` and never reimplements the formula, which is the substantive protection the gate is checking for. No violation found.

Blocking-failure check: **No** — no divergent component (nothing to diverge; all nil/honest-unknown) and no duplicated formula in orders.

Evidence-type: static/ran for architecture (no-reimplementation, honest-nil seam). Live VALUE-parity vs `POST /pricing/simulations` = **DEFERRED-C2 (hub-owned post-merge)** → ROUTED TO HUB. Do not fail the milestone for the absence of live parity.

---

## §tela (C03)

VERDICT: **PASS-STATIC-hub-drive-pending**

Evidence (ran):
- KPI/Lista single-source: `adapters/postgres/order_repo.go:209-239` `GetOrderBucketCounts` (summary `by_status` path) and `transport/http_handler.go:419` `mapEnrichedOrder` (per-order path, commit `5051ae53` "P6 — bucket hasShipment from shipping_id (parity w/ summary)") both call `domain.DeriveOrderBucket(status, hasShipment)` fed by the **same** `shipping_id` presence signal (repo SQL selects `o.shipping_id` at lines 59/107/221; handler uses `e.Order.ShippingID != ""`). This is the fix that makes KPI and Lista counts provably consistent by construction — cannot diverge.
- Kanban strictly read-only: `grep -rn "useMutation|fetch(.*method.*POST|PUT|PATCH|DELETE" apps/web/src/pages/pedidos/*.tsx` → zero hits. `grep -n "draggable|onDrag|drop" KanbanView.tsx` → zero hits (no drag handlers exist at all).
- Cancelados disabled stub: `pedidosTabs.ts:16-21` — `concluido`/`cancelado`/`devolucao` all carry `placeholder: true`; `isLiveBucketTab`/`filterOrdersByTab` only recognize `novo|faturar|enviar|enviado` as live buckets — cancelado renders as a non-functional placeholder tab, matches "Non-Scope stub" requirement.
- Full dataset load (KPI/Lista agreement precondition): `PedidosPage.tsx:45-59` `fetchAllOrders` follows `next_cursor` in a loop (capped at `MAX_ORDER_PAGES=20`, with an honest console.warn if the cap is hit rather than a silent "complete" claim) — confirms the "fetch loads the full dataset" requirement.
- Lanes: vitest `src/pages/pedidos` → 22/22 tests pass (`PedidosPage.test.tsx`).

Evidence-type: static/ran for the counts-consistency mechanism, Kanban read-only guarantee, and Cancelados stub. Light+dark browser render, pixel/interaction QA = **could-not-run (stack hub-owned)** → ROUTED TO HUB.

---

## §drawer (C04)

VERDICT: **PASS-STATIC-hub-drive-pending** (honesty verified); margin color-chip + DIFAL disclaimer label = **DEFERRED-C2**

Evidence (ran):
- Custo ausente ⇒ UnknownValue, never fabricated: `PedidoDrawer.tsx:105-171` `DecompRow`/`DecomposicaoSection` — every decomposition field routes through `formatMoney`/`formatPercent` (`pedidosFormatters.ts:16-41`, both null-guarded: `value === null || undefined → null`), and `DecompRow` renders `value ?? <UnknownValue hint={...}/>`. `componentes_desconhecidos` visible: lines 129, 137-141 render a "Pendente: …" banner listing the unresolved components whenever `decomposicao.componentes_desconhecidos.length > 0`.
- Retorno never shown as complete: `retorno_liquido` (`DecompRow` line 153) goes through the same null-guarded `formatMoney` — with the decomposer nil (C1 today) this is always `null` → always renders `UnknownValue`, never a fabricated partial number.
- Two-case behavior (known-cost vs unknown-cost item) is exercised in `ItemRow` (`PedidoDrawer.tsx:62-80`): `custo` is `null` when `item.custo_unitario === undefined`, rendering `UnknownValue hint="sem custo ERP no vínculo"`; when present it renders the formatted amount — both code paths exist and are exercised by `PedidosPage.test.tsx` (22 passing tests cover drawer states).
- Mutation buttons disabled: `PedidoDrawer.tsx:282-307` `DrawerActions` — every button `disabled`, `title="disponível em breve"` (D-57 compliant).
- Margin color-chip (verde/âmbar/vermelho) and DIFAL "seed padrão 2026 — não é orientação fiscal" label: `grep -rn "seed padrão|não é orientação fiscal" apps/web apps/server_core` → **zero hits**. Neither exists in the codebase yet. This is consistent with the contract's own framing (dispatch prompt: "depend on real C2 values — note which exist now vs C2-deferred") — since `margem_pct`/`difal` values are always null under C1 (no decomposer wired), no chip is rendered at all (no fabricated color, no fabricated label) — this is the honest behavior, not a defect. The chip/label component work is C2-deferred.

Evidence-type: static/ran for honesty guarantees (no fabrication under absent cost). Two-case live screenshots (light/dark), and the margin-chip/DIFAL-label UI itself = **could-not-run (stack hub-owned)** / **DEFERRED-C2** → ROUTED TO HUB.

---

## §seams (C05)

VERDICT: **PASS**

Evidence (ran):
- `ls apps/server_core/migrations | grep -E '^006[0-4]'` → **no output** (absence confirmed) — matches plan expectation: M-08 ships NO migration (C1 seam is Go-only, additive; no `orders_*` schema change).
- `git diff 8b6c4b30..5051ae5 --stat` (35 files, +4276/-11) — all writes fall within ownership:
  - `apps/server_core/internal/composition/orders_adapters*.go` (composition wiring — expected per D-62 note on where an orders-side adapter over connectors must live to avoid an import cycle)
  - `apps/server_core/internal/modules/orders/**` (adapters/postgres, application, domain, ports, transport)
  - `apps/web/src/app/AppRouter.test.tsx` (+4/-? — the one explicitly-allowed touch)
  - `apps/web/src/pages/pedidos/**` (new page module)
  - `apps/web/src/routes/pedidos.tsx`
  - `contracts/api/marketplace-central.openapi.yaml` (additive `/orders*`)
  - `packages/sdk-runtime/src/index.ts` (additive, orders region)
  - No file outside this list. No write into `modules/pricing`, `modules/market`, or any other module's ownership.
- Decomposição seam gate: `9c3a6812 feat(orders): F01-C1 decomposição/DIFAL/retorno seam — honest-empty (additive)` is the seam commit. As established in §decomposicao, C1 never consumes M-07's actual decompose engine (`d288cd12`, merged via `42e0b204`) — it defines a module-local nil-able port instead — so there is no ordering violation to check against those specific SHAs; the substantive protection (no formula duplication, no cross-module import) is independently verified. Ledger `b6b6beb6 docs(hub): ledger D-62 — CHIP-M08 IC-04 reachability ruled (C1-now + C2-hub-post-merge)` documents this ruling.
- Lanes: `go build`/`go vet` clean on `orders/...` + `composition/...`.

Blocking-failure check: **No** — no migration out of block (none exists, none expected), no write outside ownership, no formula-duplication/gate violation found.

---

## §pii (C06)

VERDICT: **PASS**

Evidence (ran):
- Buyer DTO on transport: `http_handler.go:310-317` `enrichedBuyerDTO{Display, City *string, UF *string}` — the ONLY buyer fields the transport can emit are Display (masked) + honest-possibly-nil City/UF.
- Masking logic: `domain/order_enrichment.go:15-40` `MaskBuyer` — `Display` is derived as `tokens[0]` (first token) plus, if a second token exists, `" " + firstRuneUpper(tokens[1]) + "."` (initial only). Never more than first name + initial. Empty/blank nickname ⇒ empty `MaskedBuyer{}` (no fabrication).
- Raw nickname field: `enrichedOrderDTO` embeds `domain.OrderReadModel` which carries `BuyerNickname *string json:"buyer_nickname"` (no `omitempty`) — confirmed this raw field is **never populated** anywhere in the orders module: `grep -rn "BuyerNickname" .../orders/` (excluding tests) → only 2 references, both READS in `enrich_service.go:82-83` (used solely as input to `MaskBuyer`); `order_repo.go` has zero `BuyerNickname` references (repo scan never populates it). So `buyer_nickname` serializes as `null` in every response — an inert always-nil key, not a live PII leak. (Note for hub: this relies on no future repo change silently populating `BuyerNickname` without updating the DTO — worth a follow-up guard, not a blocker.)
- No CPF/CNPJ/email/phone/full-address/full-name anywhere in the read path: `domain/order.go:63-82` `MarketplaceOrderItem`/`MarketplaceOrderPayment` structs carry no buyer-identifying fields (only product/payment amounts, provider IDs, SKU/title). `PedidoDrawer.tsx:260-266` explicitly renders "Documento" as `<UnknownValue hint="documento do comprador não disponível no backend"/>` — comment states no CPF/CNPJ field exists on `OrderRead` at all (never fabricated).
- Logs: `enrich_service.go:114-121` shipment-lookup warn logs `installation_id`, `provider_order_id`, `shipping_id`, `error` only — no buyer field logged.

Blocking-failure check: **No** unmasked buyer field found in payload, log, or UI.

---

## Design 1:1 (structural, read-only)

Source: `git show main:.mnfs/MIS-004-mvp-demo/design/handoff/Pedidos.dc.html`, `git show main:.mnfs/MIS-004-mvp-demo/design/DESIGN-REFERENCE.md` (line 67 Pedidos row).

| Design element | Implemented | Match |
|---|---|---|
| KPIs: Novos/A faturar/A enviar/Enviados/DIFAL a pagar | `PedidosPage.tsx:157-184` — exactly these 5 cards, in this order | MATCH |
| Views: Fila (default)/Lista (abas+contagem)/Kanban (read-only) | `PedidosPage.tsx:16-22,84-85` — `view` defaults to `"fila"`; `viewOptions` = Fila/Lista/Kanban | MATCH |
| Kanban: "arrastar NÃO muda status" (design line 55) | No drag handlers exist at all in `KanbanView.tsx` (grepped, zero) — stricter than "drag exists but doesn't mutate," there is no drag UI at all | MATCH (superset-safe) |
| Drawer section order: itens+vínculo, decomposição+DIFAL, timeline ML→interno, NF, rastreio, comprador | `PedidoDrawer.tsx:271-280` `DrawerBody`: `ItemsSection → DecomposicaoSection → TimelineSection → FactsSection` (Facts contains NF, Rastreio, Destino, Código rastreio, Comprador, Documento in that order) | MATCH |
| DIFAL column in Lista | `PedidosTable.tsx:113-114` `key: "difal", header: "DIFAL"` | MATCH |
| DIFAL coluna-chip (pago✓/vence hoje/ag dd-mm) | Not present as a colored chip yet — DIFAL column renders `formatMoney(item.difal?.amount)` or `UnknownValue`; the chip states depend on live paid/due_date values that are C2-deferred (always null under C1) | STRUCTURAL SEAM PRESENT, chip styling C2-deferred (not a defect at this bar) |
| Bloco DIFAL in drawer (agendar/lembrete/marcar pago) | `DrawerActions` renders "DIFAL agendar" button, disabled (D-57) | MATCH (read-only rich) |

No structural mismatches found. Per ratified design memory, the M-05 Anúncios prototype's vs-mercado-column omission is a different screen and does not apply to Pedidos.dc.html.

Pixel-exact 1:1, light/dark rendering = **could-not-run (stack hub-owned)** → ROUTED TO HUB.

---

## §lanes (evidence)

| Lane | Command | Result |
|---|---|---|
| Go build | `go build ./apps/server_core/internal/modules/orders/... ./apps/server_core/internal/composition/...` (GOCACHE absolute) | clean, no output |
| Go vet | `go vet` same packages | clean, no output |
| Go test | `go test` same packages | **9/9 packages ok** (adapters/integrations, adapters/internalread, adapters/postgres, adapters/productlinks, application, domain, ports, transport, composition) |
| vitest | `npx vitest run src/pages/pedidos` | **1 file, 22/22 tests passed** |
| tsc --noEmit | `npx tsc --noEmit` (apps/web) | Errors present, but **all** are: (a) `TS2339 ... toBeInTheDocument/toBeDisabled/toHaveTextContent` jest-dom matcher errors — confirmed pre-existing project-wide baseline (identical errors also fire in unrelated pre-existing files `vinculos/QueueTab.test.tsx`, `vinculos/VinculosPage.test.tsx`, `vinculos/BatchPreviewModal.test.tsx`, `vinculos/ImportacaoSection.test.tsx`, none touched by this diff), or (b) `src/routes/precos.tsx TS2322` — a pre-existing M-07 file, not part of this diff's file-scope (§seams). **No new non-baseline tsc error introduced by M-08.** |

---

## ROUTED TO HUB (not chip-satisfiable, browser stack owned / C2 post-merge)

1. C01 — live `GET /orders` response bodies against a real ML installation (post-sync).
2. C02 — live value-parity of decomposição components between Pedidos drawer and `POST /pricing/simulations` (requires C2 decomposer wiring).
3. C03 — browser live-drive of `/pedidos` (light+dark render, KPI-vs-Lista live agreement, Kanban interaction, network transcript).
4. C04 — two-case drawer screenshots (known-cost vs unknown-cost order); margin color-chip (verde/âmbar/vermelho) and DIFAL "seed padrão 2026 — não é orientação fiscal" label rendering (both depend on C2 values/UI not yet built).
5. Design — pixel-exact 1:1 comparison, light/dark theme rendering.
6. The C05 "gate respected" check as literally worded (ordering vs M-07 ports SHA) does not strictly apply since C1 never consumes those ports directly (see §decomposicao) — hub should confirm this reading is accepted, or re-verify once C2 wires the real decomposer.

## OVERALL

**Safe to CLOSE on the READ-ONLY-RICH + C1 bar (D-57/D-62).**

All 6 criteria pass their chip-verifiable half:
- C01: PASS-STATIC (DTO additive, no provider leak, ADR-17 honest-nil, cost sourced with visible time) — live values routed to hub.
- C02: PASS (architecture — no pricing import, no formula duplication, honest-nil decomposer port) — value parity DEFERRED-C2.
- C03: PASS-STATIC (single-source bucket derivation proven by shared `shipping_id` signal, Kanban has zero write/drag surface, Cancelados is an inert placeholder, full-dataset load confirmed) — browser render routed to hub.
- C04: PASS-STATIC (no fabricated margin/retorno under absent cost, both cost-present/cost-absent code paths exist and are unit-tested) — chip/label UI + screenshots DEFERRED-C2/routed to hub.
- C05: PASS (no migration, all writes in-ownership, no formula duplication).
- C06: PASS (buyer DTO is display+city+uf only; raw `buyer_nickname` field is structurally always-nil; no CPF/CNPJ/email/phone/full-name/full-address anywhere in payload, log, or UI).

No blocking failure was found for anything chip-verifiable: no write leak, no fabricated value, no PII exposure, no ownership breach, no migration out of block. All go/vitest lanes pass; the only tsc errors are the pre-existing jest-dom baseline plus one unrelated pre-existing file outside this diff's scope.

Design structure matches Pedidos.dc.html/DESIGN-REFERENCE.md with no mismatches; the DIFAL coluna-chip styling and drawer margin-chip/disclaimer label are honestly absent (C2-deferred), not defects.

---

## HUB P7 BROWSER LIVE-DRIVE — CLOSE STAMP (2026-07-19, D-72)

```yaml
gate: P7 browser QA (hub live-drive)
verdict: PASS → M-08 CLOSED
stack: integrated main @72b083e, MC_ERP_SOURCE=xlsx (#003-E demo-client snapshot), backend :8080 / frontend :5174
driver: hub (fresh live-drive on the integrated stack)
```

Discharges every ROUTED-TO-HUB item from the cold-static verdict above:

- **C01 live GET /orders** — DISCHARGED. 24 real mercado_livre orders (bucket A ENVIAR), real ML order ids + lifecycle timestamps (09/07/2026); list, summary?by=status, and detail /orders/{id} enrich all returned 200. No 500 on the read/enrich path (distinct from the M-07 decompose blocker D-71).
- **C02 value-parity decomposição** — remains **DEFERRED-C2 → post-demo backlog** per D-70 (order carries real SaleFeeAmount but no modalidade; feeding IC-04 simulator engine would fabricate modalidade/fee = ADR-17 violation). NOT a close-blocker: operator ruled pedidos ships the honest "—" it renders; the real-numbers story is the simulador (M-07).
- **C03 browser live-drive /pedidos** — DISCHARGED. KPI cards + all 3 views (Fila default / Lista abas+contagem / Kanban read-only) render light + dark. KPI↔Lista agree (A ENVIAR 24 in both). Kanban has no drag surface ("sem arrastar · ações nos cards em breve").
- **C04 two-case drawer + honesty** — DISCHARGED. Cost-present path: item vinculado + CODPROD 15956 with "custo incompleto" honest badge (custo unit "—"). Decomposição/DIFAL/margem/retorno all render UnknownValue "—" with the "Pendente: comissao, taxa_fixa, frete, imposto, difal, tarifa_full, custo — … mostram '—' até a decomposição ser calculada" ComponentesDesconhecidos banner. No fabricated value anywhere. Margin color-chip + DIFAL "seed padrão 2026" disclaimer label remain C2-deferred (honestly absent under C1, not defects).
- **C06 PII (live)** — CONFIRMED in-drive. Drawer Comprador "—" / Documento "—"; Lista COMPRADOR column "—". No raw buyer nickname/CPF/CNPJ/email/phone rendered. LGPD-safe honest degrade.
- **Design pixel/theme 1:1 vs Pedidos.dc.html** — DISCHARGED (structural). Tokens hold light + dark (paper/dark bg, green accent, IBM Plex Mono numerals); KPI order + view set + drawer section order match the cold-static structural table above.
- **C05 gate-wording reading** — ACCEPTED by hub (D-62/D-70): C1 defines a module-local nil-able Decomposer port and never imports modules/pricing, so the literal "consume M-07 ports SHA" ordering gate does not apply; the substantive protection (no formula duplication, no cross-module import) was independently verified. Re-verify when C2 wires the real decomposer (post-demo).

**Zero Mercado Livre writes**: every mutation control disabled ("disponível em breve" fila+drawer; Kanban "sem arrastar"); MPC_PROVIDER_WRITES_ENABLED unset holds.

**OVERALL: M-08 pedidos CLOSED on the read-only-rich + C1 honest-"—" bar (D-57/D-62/D-70). Only QA passes a milestone — this P7 is that QA and it PASSED.** C2 order-actuals decomposition = post-demo backlog.

---

## ⚠ CLOSE REVERSED — M-08 REOPENED (2026-07-19, D-73)

The D-72 P7 CLOSE STAMP above is **REVERSED**. The hub P7 "PASS" was wrong: the pedidos KPI "A ENVIAR 24 / ENVIADOS 0" is a DEFECT, not real data. Operator (ground truth on the real ML account) flagged that most of those orders are already shipped/delivered; the hub had wrongly accepted the blank Rastreio/"—" as honest-degrade when it was a live shipment-read failure.

Root cause (Context7 ML-docs-grounded): the connector omits the required `x-format-new: true` header on GET /shipments/{id}/costs → ML returns the legacy shape → decode fails (CONNECTORS_PROVIDER_PAYLOAD_INVALID); getShipmentInfo treats that decode failure as fatal (fallback catches only 404) → the whole shipment read dies → shipment status is lost → DeriveOrderBucket falls back to order.provider_status ("paid" for life on ML) → every paid order with a label buckets A ENVIAR, ENVIADO unreachable. KPI compound: GetOrderBucketCounts is pure SQL with no shipment-status column.

Correction dispatched: **CHIP-M08-SHIPFIX** (task_0ed1ee6f) — connector x-format-new + non-fatal costs-decode degrade + substatus; DeriveOrderBucket consumes shipment status; FE-derived KPI counts. M-08 is NOT closed; re-close requires that chip landing + a fresh hub P7 live-drive. See HUB-LEDGER D-73.
