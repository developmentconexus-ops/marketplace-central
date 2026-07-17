# Milestone Validation Contract

```yaml
id: M-09
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: milestone
```

## Milestone ID

M-09-dashboard-demo (CORTÁVEL — se cortado, registrar no mission close; demo abre em /anuncios; este contrato não roda).

## QA Level

Dual gate frio (Opus + Sol medium) + QA live-drive browser fresh (/).

## Required Outcome

Visão geral reconstruída nos tokens: KPIs reais das fontes dos milestones produtores (M-01/M-03/M-04/M-05/M-08), Fila de atenção cross-módulo com deep-links, pedidos recentes, atalhos.

## Criteria

## Criterion: KPIs batem com as telas de origem
ID: M-09-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: anotar cada KPI do dashboard e comparar com a tela de origem no mesmo instante
- Expected: vendas hoje/semana = mesmos números de /pedidos (KPIs M-08); anúncios ativos/com exceção = /anuncios (summary M-05); produtos sem vínculo = pendentes de /vinculos (M-04); último import ERP = protocolo de M-01; divergência zero entre card e origem
- Actual:
- Artifact: `M-09-dashboard-demo/validation-result.md` §kpis (screenshots pareados card↔origem)
Blocking failure: qualquer KPI divergente da tela de origem
Blocking failure observed: No
Drive (UI — agent-browser):
- Fixture: stack local com waves A+B fechadas e dados reais (import + vínculos + coleta + orders sync); sem auth
- Steps:
  - open http://localhost:5174/
  - assert text <valor KPI vendas hoje>
  - open http://localhost:5174/pedidos
  - assert text <mesmo valor no KPI de /pedidos>
- Expected: mesmo número nas duas telas
Owner: QA Validator

## Criterion: Fila de atenção deep-linka com filtro
ID: M-09-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: clicar item da Fila de atenção
- Expected: navegação leva à tela dona com o filtro da exceção JÁ aplicado (ex.: pendentes de vínculo → /vinculos filtrado em Pendentes); item de exceção sem tela dona não existe
- Actual:
- Artifact: `M-09-dashboard-demo/validation-result.md` §fila (transcript navegação)
Blocking failure: deep-link caindo na tela sem filtro, ou em rota errada
Blocking failure observed: No
Owner: QA Validator

## Criterion: Dado ausente honesto
ID: M-09-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive com uma fonte vazia (ex.: nenhum import ainda no DB de teste)
- Expected: card exibe estado honesto via UnknownValue/vazio explicativo (ADR-17) — NUNCA zero fabricado que se pareça com medição real
- Actual:
- Artifact: `M-09-dashboard-demo/validation-result.md` §ausente (screenshot)
Blocking failure: fonte vazia renderizada como 0 indistinguível de medição
Blocking failure observed: No
Owner: QA Validator

## Criterion: Extensão aditiva /dashboard/summary
ID: M-09-C04
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: GET /dashboard/summary antes/depois (contrato antigo); inspecionar leituras cross-módulo
- Expected: campos antigos preservados (aditivo); agregações novas via ports/APIs públicas dos módulos donos — grep por query SQL cross-schema em modules/dashboard = zero
- Actual:
- Artifact: `M-09-dashboard-demo/validation-result.md` §api (responses + grep)
Blocking failure: campo antigo quebrado, ou query direta em tabela de outro módulo
Blocking failure observed: No
Owner: QA Validator

## Criterion: Ownership limpo
ID: M-09-C05
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: diff do chip vs matriz; lanes L0–L2
- Expected: writes só em `modules/dashboard/**`, `sdk-runtime/src/dashboard.ts`, OpenAPI `/dashboard/*` aditivo, `apps/web/src/pages/dashboard/**` (rebuild), `routes/dashboard.tsx`; ZERO migração (tabela nova exigiria REQUEST reserva 0070+); light+dark OK
- Actual:
- Artifact: `M-09-dashboard-demo/validation-result.md` §seams
Blocking failure: migração auto-atribuída ou write fora do ownership
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- `M-09-dashboard-demo/validation-result.md` com seções kpis, fila, ausente, api, seams; screenshots light+dark.
- Se CORTADO: nada aqui roda; registro do corte no `validation-result.md` da missão (MIS-004-C03).

## Blocking Failures

Depende dos produtores reais {M-01, M-03, M-04, M-05, M-08} fechados (fixture). M-07 não produz artefato consumido aqui.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned (P6).
- Next owner: QA Validator (pós-close F-01 do chip M-09) — ou hub registra corte.
- Next action: decisão corte/execução na entrada da wave C.
- Required files/evidence: `M-09-dashboard-demo/validation-result.md` (ou registro de corte).
- Blockers or open decisions: corte = decisão hub/operador.
