# M-02-price-intel-core

```yaml
id: M-02
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

MIS-004 mvp-demo — fundação de evidência de preço/mercado.

## Outcome

Adapter ML estendido (ports IC-06, read-only, flag catalog_offers OFF); persistência de evidência no `market` (IC-03, ADR-17); resolver determinístico (IC-01); coleta on-demand; API de sinais/agregados/veredictos servindo M-05/M-06/M-07.

## Why This Milestone Exists

Núcleo da história "preço vs mercado honesto". Research de pricing (binding) fecha o COMO; este milestone materializa.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | ml-adapter-read-ports | `F-01-ml-adapter-read-ports/feature.md` |
| F-02 | market-persistence | `F-02-market-persistence/feature.md` |
| F-03 | identity-resolver | `F-03-identity-resolver/feature.md` |
| F-04 | collect-verdict-api | `F-04-collect-verdict-api/feature.md` |

## Dependencies

Nenhuma de código de outros milestones (wave A). IC-01 ratificado (planning) alimenta F-03. Installation ML conectada (preflight R4/R6).

## Ownership & Concurrency

- Exclusive surfaces: `modules/market/**`, `modules/connectors/**` (capability_adapter ML + ports novos), OpenAPI seção `/market/*` (novos endpoints; observations/references SAT intocados), `sdk-runtime/src/market.ts`, tabelas `market_price_snapshots`, `market_validated_offers`, `market_aggregates`, `market_competitive_signals`, `market_match_decisions`.
- Migration block: **0050–0054** (+ bump fixture).
- Predicted seam locks: export 1-linha barrel SDK (hub); registro composition root via merge (hub).
- Runs in parallel with: M-01, M-03.
- Internal feature DAG: `F-01 ∥ F-02` → F-03 → F-04.

## Risks

R4 (permissões conta ML — preflight na 1ª hora do chip); R5 (rota flag instável — fallback tipado); R2 da missão (evidência escassa — não é falha do milestone: estados honestos SÃO o comportamento correto).

## Done Means

Lane live-provider-read verde com installation real (sale_price + price_to_win de anúncio próprio ativo); teste negativo ADR-17 passa; fixtures de colisão IC-01 passam no resolver; telemetria da rota flag inspecionável; flag OFF ⇒ NO_PRICE_EVIDENCE explícito; dual gate + QA PASS.

## Handoff

- Current status: planned.
- Next owner: hub (dispatch chip wave A).
- Next action: chip implementa F-01∥F-02 → F-03 → F-04; preflight conta ML primeiro.
- Required files/evidence: `validation-result.md`; log lane live-provider-read.
- Blockers or open decisions: none.

## Correction Handoff

n/a (planning).
