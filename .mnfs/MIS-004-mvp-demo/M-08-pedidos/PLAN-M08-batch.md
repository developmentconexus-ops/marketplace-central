# M-08-pedidos — P2 Batch Plan (chip CHIP-M08)

Base SHA: 8b6c4b3093f9465cd3b91209b054af4fa702171a · Worktree: .claude/worktrees/chip-m08-pedidos
Planner: cold Opus subagent a9daa09eb1a70de71 (contingency lane, plan-only), 2026-07-18.

## Key architecture decision (G2)
NO projection cache tables → NO migration. Enrichment (cost/shipment/decomposition) computed
LIVE at read time. Migration block 0060–0064 stays RESERVED-UNUSED; runner_test.go (51) NOT bumped.
C05 satisfied by ABSENCE of any migration outside the block. If a hard persistence need is proven
mid-impl, use 0060_* + bump fixture in the SAME slice = scope escalation to hub, not a silent add.

## Gating truth
DECOMPOSITION-GATED on M-07 (commit real wiring only after hub relays `ic04-ports:<sha>`):
retorno_liquido, margem_pct, decomposicao, difal, exception=margem_negativa, summary retorno/difal
aggregates. NON-GATED (ship first): buyer mask, vinculo_status, ERP cost, shipment SLA/UF/rastreio,
summary counts, status_tab/exception=atrasado|sem_vinculo filters, SLA ordering, timeline.

## Additive contract-locks
- SDK-ORDERS-M08 (index.ts orders region, additive) — HUB GRANT PENDING (REQUEST sent; supersedes
  orders.ts-NEW + BARREL-M08). All F-01 contract-bearing slices touch it.
- OpenAPI /orders* additive (my exclusive write).
- ROOT-M08 (root.go ORDERS region additive wiring). APPTEST-M08 (/pedidos test case only).
- IC-04 gate = sequencing gate for S6b real wiring.

## F-01 slices (single sequential writer)
- S1 buyer PII mask + vinculo_status + enrichment scaffold (OrderEnricher seam). NON-GATED. ready.
- S2 ERP cost per item via GetCostAsOf (nil⇒null+componentes_desconhecidos, never 0). NON-GATED. ready.
- S3 shipment SLA/atraso/rastreio/destino_uf via IC-06 (errors⇒nulls+telemetry, list never fails).
  COMPLEX. NOT dispatch-ready → needs ShipmentReader/ProviderAccountRef wiring path (investigator).
- S4 GET /orders/summary counts + REGISTER unwired handler (ordersSummarySvc:500). NON-GATED. ready.
- S5 additive filters (status_tab, exception=atrasado|sem_vinculo, SLA order) + timeline in /orders/{id}. ready.
- S6a decomposition/retorno/margem/difal via LOCAL frozen IC-04 interface + fakes. COMPLEX. ready (fakes).
- S6b real pricing-port adapter + ROOT wiring + golden live test. GATED on ic04-ports:<sha>.
- S7 gated summary aggregates + exception=margem_negativa. logic ready; live verify gated.

## F-02 slices (single sequential writer, pages/pedidos/** + route)
- S1 route swap + PedidosPage + Lista tabs + states (Cancelados disabled no fetch) + APPTEST-M08.
- S2 KPI row clickable + DIFAL banner + Fila de atenção.
- S3 drawer (decomposition + timeline + rastreio + itens + masked buyer) + ?order= deep-link.
- S4 Kanban read-only (same query, zero fetch/write) + view switch + light/dark.

## Open gates
1. F01-S3 investigator (ShipmentReader bound to installation ProviderAccountRef).
2. IC-04 gate: S6b + live verify wait for hub ic04-ports:<sha>.
3. SDK-ORDERS-M08 hub grant (all F-01 contract-bearing slices; index.ts additive).
4. Non-blocking note: orders storage lacks buyer full-name/city → MVP mask = nickname + shipment UF.

Full plan detail: planner transcript (agent a9daa09eb1a70de71). Verification map C01–C06 + seam-closure
checklist recorded there; C05 = no migration (absence-verified).
