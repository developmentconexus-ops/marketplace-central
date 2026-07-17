# M-09-dashboard-demo

```yaml
id: M-09
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

MIS-004 mvp-demo — visão geral que abre a demo. CORTÁVEL (P1a decisão: se tempo apertar, cai primeiro).

## Outcome

`/` (Visão geral) reconstruída nos tokens: KPIs reais (vendas hoje/semana via orders, anúncios ativos/com exceção via listings, produtos sem vínculo via product_links), Fila de atenção (top exceções cross-módulo), pedidos recentes, atalhos. `/dashboard/summary` estendido aditivamente com as agregações que faltarem.

## Why This Milestone Exists

Primeira tela que o cliente vê; hoje DashboardPage é dark antigo com dados de outra era. Wave C.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | dashboard-mpc | `F-01-dashboard-mpc/feature.md` |

## Dependencies

- Produtores reais (DAG da missão): M-01 (import/custo), M-03 (tokens/seam), M-04 (vínculos), M-05 (sinais listings), M-08 (orders projection). M-07 NÃO produz artefato consumido pelo M-09.
- Contratos: IC-05.

## Ownership & Concurrency

- Exclusive surfaces: `modules/dashboard/**`, OpenAPI seção `/dashboard*` (aditivo), `sdk-runtime/src/dashboard.ts`, `apps/web/src/pages/dashboard/**` (rebuild), `apps/web/src/routes/dashboard.tsx`.
- Migration block: none (agrega leituras; precisa de tabela ⇒ REQUEST reserva 0070+).
- Predicted seam locks: export barrel SDK (hub); leituras cross-módulo APENAS via ports/APIs públicas dos módulos donos.
- Runs in parallel with: M-06 (wave C).
- Internal feature DAG: F-01 único.

## Risks

Cortável — se cortado, `/` mantém DashboardPage atual (funcional, feia) e demo abre em /anuncios; KPI com dado ausente mostra estado honesto (nunca zero fabricado, ADR-17).

## Done Means

KPIs batem com as telas de origem (mesmo número em /pedidos e no card); Fila de atenção deep-linka para tela certa com filtro aplicado; dado ausente ⇒ honesto; light/dark OK; dual gate + QA PASS. Se cortado: registro no mission close + demo abre em /anuncios.

## Handoff

- Current status: planned.
- Next owner: hub (dispatch chip wave C, pós-merge wave B; ou corte).
- Next action: chip implementa F-01.
- Required files/evidence: `validation-result.md`; screenshot dashboard light/dark.
- Blockers or open decisions: corte é decisão do hub/operator na entrada da wave C.

## Correction Handoff

n/a (planning).
