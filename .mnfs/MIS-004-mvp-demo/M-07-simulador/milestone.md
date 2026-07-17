# M-07-simulador

```yaml
id: M-07
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

MIS-004 mvp-demo — simulador de margem real (comissão + frete + imposto + DIFAL + custo ERP).

## Outcome

Backend: CalcProfile persistido, tabela DIFAL 27 UFs seedada (IC-04), motor de decomposição único, simulação bidirecional (preço→margem, margem-alvo→preço), cenários salvos, port `DifalForUF` para M-08. Frontend: `/precos` reconstruída per design (picker de anúncio, tabela com veredicto, painel bidirecional com decomposição completa, drawer de parâmetros com labels de origem, toggle DIFAL, cenários, "aplicar" enfileira mutação local SEM execução).

## Why This Milestone Exists

"Simulador de margem real" é metade do pitch. Simulador atual ignora DIFAL/frete real e usa custo chumbado. DIFAL é dor declarada do cliente.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | pricing-calc-difal-service | `F-01-pricing-calc-difal-service/feature.md` |
| F-02 | simulador-ui | `F-02-simulador-ui/feature.md` |

## Dependencies

- M-01 (custo via Reader — edge de dado; F-01 programa contra port IC-02 §Reader).
- M-02 (preço de mercado como referência — via IC-03 GET; ausência ⇒ estado honesto).
- M-03 (tokens/primitivas/MarginChip) para F-02.
- Contratos: IC-04 (formula/decomposição/DIFAL/endpoints), IC-06 (GetFreeShippingCost/GetShipmentInfo via ports M-02).

## Ownership & Concurrency

- Exclusive surfaces: `modules/pricing/**`, OpenAPI seção `/pricing/*` (aditivo — simulations existentes preservadas), `sdk-runtime/src/pricing.ts`, `apps/web/src/pages/precos/**`, `packages/feature-simulator/**` (PricingSimulatorPage legado — retematizado/absorvido por F-02, package sobrevive per IC-05), `apps/web/src/routes/precos.tsx`, tabelas `pricing_calc_profiles`, `pricing_difal_rates`, `pricing_scenarios`.
- Migration block: **0055–0059** (+ bump fixture).
- Predicted seam locks: export barrel SDK (hub); M-08 consome `DifalForUF` (port público — lock aditivo pré-autorizado: assinatura congelada no IC-04, M-07 não a altera pós-publicação).
- Runs in parallel with: M-04, M-05, M-08 (wave B).
- Internal feature DAG: F-01 → F-02.

## Risks

Precisão fiscal (mitigação: label obrigatória "seed padrão 2026 — não é orientação fiscal" + override editável); decomposição divergir entre simulador e pedidos (mitigação: motor único no pricing, M-08 consome via port — nunca reimplementa).

## Done Means

Simulação com todos os componentes conhecidos bate com cálculo manual (golden tests IC-04); custo desconhecido ⇒ margem desconhecida (nunca zero); margem-alvo→preço converge; DIFAL muda com UF; toggle remove componente com aviso; cenário salvo recarrega; "aplicar" cria mutação PENDING local e execução NÃO ocorre; dual gate + QA PASS.

## Handoff

- Current status: planned.
- Next owner: hub (dispatch chip wave B).
- Next action: chip implementa F-01 → F-02.
- Required files/evidence: `validation-result.md`; screenshot decomposição + golden test log.
- Blockers or open decisions: none.

## Correction Handoff

n/a (planning).
