# M-08-pedidos

```yaml
id: M-08
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-004 mvp-demo — Pedidos funcional com margem real por pedido.

## Outcome

Backend: projeção de lista de pedidos enriquecida (SLA/atraso via shipment IC-06, retorno líquido + decomposição via port pricing IC-04, chip DIFAL por UF destino, custo ERP via Reader), KPIs, filtros/tabs por estado. Frontend: `/pedidos` per design — KPIs clicáveis, Fila de atenção rankeada, Lista com tabs, Kanban read-only, drawer com decomposição de margem + timeline + rastreio, comprador mascarado (LGPD), banner agregado DIFAL.

## Why This Milestone Exists

"Pedidos funcional" é pedido explícito do cliente (brief §2). Margem real por pedido conecta ERP→venda→lucro na demo.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | orders-projection-api | `F-01-orders-projection-api/feature.md` |
| F-02 | pedidos-ui | `F-02-pedidos-ui/feature.md` |

## Dependencies

- M-01 (custo — edge de dado via Reader IC-02).
- M-02 (ports `GetShipmentInfo`/`GetFreeShippingCost` IC-06 — edge de CÓDIGO: F-01 consome ports publicados pelo M-02; hub sincroniza merge M-02 antes do trecho shipment, resto do F-01 independe).
- M-07 (`DifalForUF` + motor decomposição via port IC-04 — edge de código; lock aditivo pré-autorizado no IC-04).
- M-03 (tokens) para F-02. M-04 (vínculo pedido-item→CODPROD — dado; sem vínculo ⇒ margem desconhecida honesta).

## Ownership & Concurrency

- Exclusive surfaces: `modules/orders/**`, OpenAPI seção `/orders*` (ADITIVO — sankhya-linkage/import existentes preservados), `sdk-runtime/src/orders.ts`, `apps/web/src/pages/pedidos/**`, `apps/web/src/routes/pedidos.tsx`, tabelas novas de projeção `orders_*` se necessárias.
- Migration block: **0060–0064** (+ bump fixture).
- Predicted seam locks: export barrel SDK (hub); consumo aditivo dos ports M-02/M-07 (read-only, sem edição de `modules/connectors|pricing/**`).
- Runs in parallel with: M-04, M-05, M-07 (wave B) — trechos dependentes de ports gated pelo hub (edge M-02→M-08 e M-07→M-08 no DAG da missão).
- Internal feature DAG: F-01 → F-02.

## Risks

Pedidos reais escassos na conta demo (runbook: pedido de teste próprio pré-demo); decomposição por pedido com dados parciais (payment sem shipment etc.) — matriz unknown-propagation no brief; PII comprador (mascaramento obrigatório, nunca nome/endereço completos na lista).

## Done Means

Lista mostra retorno líquido com decomposição inspecionável no drawer; pedido sem custo ERP mostra margem desconhecida (nunca zero); chip DIFAL aparece p/ UF destino ≠ SC quando dado existe; SLA/atraso honesto; Kanban read-only espelha estados; comprador mascarado; KPIs batem com lista; dual gate + QA PASS.

## Handoff

- Current status: planned.
- Next owner: hub (dispatch chip wave B; gate nos ports M-02/M-07).
- Next action: chip implementa F-01 → F-02.
- Required files/evidence: `validation-result.md`; screenshot drawer decomposição + lista.
- Blockers or open decisions: none.

## Correction Handoff

n/a (planning).
