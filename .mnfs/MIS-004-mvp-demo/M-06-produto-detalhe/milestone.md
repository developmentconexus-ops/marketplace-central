# M-06-produto-detalhe

```yaml
id: M-06
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

MIS-004 mvp-demo — página do produto: o veredicto "vale a pena vender?".

## Outcome

`/catalogo/produtos/:productId` real: header de identidade (CODPROD/EAN/REFFORN/custo/completude), box VEREDICTO (estado honesto IC-03 + faixa de preço + evidência citada), tab Anúncios vinculados (posição vs mercado por anúncio), tab Estoque (KPIs físico/reservado/disponível + banner cobertura). Composição 100% client-side de APIs existentes (M-01/M-02/M-04/M-05) — ZERO endpoint novo.

## Why This Milestone Exists

É a cena central do pitch da demo (estoque → veredicto). Wave C: consome tudo da wave B.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | produto-detalhe-page | `F-01-produto-detalhe-page/feature.md` |

## Dependencies

- M-01 (identidade/custo/estoque via APIs catalog/erp), M-02 (GET /market/verdicts|aggregates), M-04 (vínculos), M-05 (sinais por listing), M-03 (shell/tokens).
- Contratos: IC-01, IC-03, IC-05.

## Ownership & Concurrency

- Exclusive surfaces: `apps/web/src/pages/produto/**` (novo), `apps/web/src/routes/produto.tsx` (conteúdo pós-seam). SEM superfícies backend, SEM OpenAPI, SEM SDK novo (consome clients existentes).
- Migration block: none.
- Predicted seam locks: none.
- Runs in parallel with: M-09 (wave C).
- Internal feature DAG: F-01 único.

## Deliberate Exclusions (não é esquecimento)

- Tabela de movimentações de estoque (design Estoque tab): SEM fonte de dados no MVP (nenhuma API de movimentos existe) → MIS-005 M-08.
- KPI Full ML (inventories): MIS-005 M-08. Card omitido, não zerado.
- Tabs Concorrência/Pedidos/Histórico/Dados: MIS-005 M-05.

## Risks

Estados compostos (produto sem EAN + sem vínculo + sem evidência simultâneos) — matriz de estados no brief cobre; página depende de 4 milestones ⇒ chip só despacha pós-merge wave B (hub confere).

## Done Means

Produto com evidência mostra veredicto com faixa + fontes; produto sem EAN mostra REVIEW-only + call-to-action vínculo; produto sem mercado mostra INSUFFICIENT_MARKET/NO_PRICE_EVIDENCE (nunca zero); KPIs de estoque batem com snapshot importado; deep-link + F5 OK; dual gate + QA PASS.

## Handoff

- Current status: planned.
- Next owner: hub (dispatch chip wave C, pós-merge wave B).
- Next action: chip implementa F-01.
- Required files/evidence: `validation-result.md`; screenshots dos 4 estados de veredicto.
- Blockers or open decisions: none.

## Correction Handoff

n/a (planning).
