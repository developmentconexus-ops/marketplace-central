# M-05-anuncios-sinais

```yaml
id: M-05
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

MIS-004 mvp-demo — Anúncios com sinais competitivos honestos.

## Outcome

AnunciosPage (existente, W1) estendida: coluna VS MERCADO (posição vs price_to_win + estado de evidência), chips de exceção no header (sem vínculo / abaixo do custo / sem evidência), toggle agrupar-por-produto, drawer com evidência de mercado (fonte + timestamp + freshness). Backend: endpoint de sinais por listing juntando `listings` × `market`.

## Why This Milestone Exists

Anúncios é a tela-âncora da demo (única tela rica que já existe); sinais competitivos são o payoff visível do M-02.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | listings-signals-api | `F-01-listings-signals-api/feature.md` |
| F-02 | anuncios-ui-sinais | `F-02-anuncios-ui-sinais/feature.md` |

## Dependencies

- M-02 (CompetitiveSignal/MarketAggregate persistidos + GET /market/signals) — F-01 consome via read port/API interna, contrato IC-03.
- M-03 (tokens/primitivas) para F-02.
- M-04 fornece vínculos resolvidos (dado, não código) — sinal sem vínculo = estado honesto, não bloqueio.

## Ownership & Concurrency

- Exclusive surfaces: `modules/listings/**`, OpenAPI seção `/listings*` (ADITIVO — campos/endpoint novos, nada removido), `sdk-runtime/src/listings.ts` (novo; client atual em index.ts intocado), `apps/web/src/pages/anuncios/**` (AnunciosPage e filhos), `apps/web/src/routes/anuncios.tsx`.
- Migration block: none (lê market/product_links; projeção computada, sem tabela própria nova — se precisar de cache/tabela, REQUEST ao hub por bloco da reserva 0070+).
- Predicted seam locks: export barrel SDK (hub); leitura cross-módulo de `market` via port público do M-02 (IC-03) — nunca query direta nas tabelas do market.
- Runs in parallel with: M-04, M-07, M-08 (wave B).
- Internal feature DAG: F-01 → F-02.

## Risks

AnunciosPage é código vivo (regressão) — QA C08/C09 regressão deep-link obrigatória; evidência escassa na conta demo ⇒ colunas dominadas por NO_PRICE_EVIDENCE (aceitável; runbook da demo garante ≥1 anúncio com evidência VÁLIDA).

## Done Means

Anúncio vinculado com evidência mostra posição vs price_to_win + freshness; anúncio sem vínculo mostra chip "sem vínculo" (nunca zero); chips de exceção filtram a tabela; agrupar-por-produto agrupa; drawer cita fonte+timestamp; regressão AnunciosPage passa; dual gate + QA PASS.

## Handoff

- Current status: planned.
- Next owner: hub (dispatch chip wave B).
- Next action: chip implementa F-01 → F-02.
- Required files/evidence: `validation-result.md`; screenshot coluna VS MERCADO + drawer evidência.
- Blockers or open decisions: none.

## Correction Handoff

n/a (planning).
