# M-05-listings-fees-divergence

```yaml
id: M-05
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-007
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-007 ml-sync — [mission.md](../mission.md); IC-01 (camada 2, binding), IC-02
(divergências), IC-07 (fee FORA de listings — DTO compõe de channel_fees).

## Outcome

/anuncios mostra tarifa REAL por anúncio e divergência de estoque: ingest de listings (M-04)
estendido p/ gravar `channel_fees` camada 2 (origem `api_listing_prices`, sale_price com
`?context=channel_marketplace` — fato live), avaliador de divergência de estoque
(mirror ERP vs `available_quantity` ML, tolerância 0, grão variação quando existe), DTO de
/listings ganha `tarifa` composta + `divergences[]` + filtro `filter.divergentes=true`, coluna
TARIFA + badge de divergência na AnunciosTable.

## Why This Milestone Exists

Tarifa por anúncio é o insumo do M-07 (resolver lê camada 2 antes de config) e da auditoria
3→2 do M-06. Divergência de estoque é o produto visível nº1 do espelho (design §5). Depende
de M-04: camada 2 precisa de category_id/price ingeridos; divergência precisa de
available_quantity — colunas que M-04 cria (IC-07 E3).

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | camada2-fee-ingest | [F-01-camada2-fee-ingest/feature.md](F-01-camada2-fee-ingest/feature.md) |
| F-02 | stock-divergence | [F-02-stock-divergence/feature.md](F-02-stock-divergence/feature.md) |
| F-03 | anuncios-fe-contract | [F-03-anuncios-fe-contract/feature.md](F-03-anuncios-fe-contract/feature.md) |

## Dependencies

- M-04 (IngestListing + colunas E3 + listing_variations; lock aditivo pós-close — ver
  Ownership).
- M-02 (channel_fees 0086 + ports IC-01; divergences 0087 + ports IC-02).
- (M-01 transitiva via M-04.)

## Ownership & Concurrency

- Exclusive surfaces: extensões de ingest no application de listings (LOCK ADITIVO sobre
  superfície do M-04, registrado na matriz da missão — enumeração idêntica à matriz:
  `listings/application/**`, `listings/transport/**`, `listings/composition/**`,
  `listings/ports/**` e `listings/adapters/postgres/repository.go`, additive-only, M-04 JÁ
  FECHADO quando lane C abre; ampliado P7 r01 B-4), avaliador de divergência
  (arquivos novos), par OpenAPI+SDK de /listings (DTO + param), FE /anuncios
  (AnunciosPage.tsx, AnunciosTable.tsx, anunciosQueries.ts, anunciosQueryState.ts,
  ListingDetailPanel.tsx).
- Migration block: nenhum.
- Predicted seam locks: escreve channel_fees SÓ camada 2 (camada 3 = M-06; config = M-07
  leitura); divergences SÓ kind=estoque (tarifa = M-06). Contrato FE: lane C tem 3 milestones
  com par FE — hub serializa COMMITS de OpenAPI+SDK dentro da lane (ADR-14).
- Runs in parallel with: M-06, M-07 (código ∥; commits de contrato serializados).
- Internal feature DAG: F-01 ∥ F-02 → F-03 (FE compõe os dois).

## Risks

- Fee de listing muda com preço — camada 2 re-observada a cada sweep de refresh do M-04
  (upsert por natural key IC-01, `coletado_em` avança; sem histórico nesta missão).
- Grão variação: estoque ML mora na variação quando existe (IC-07) — comparar no grão errado
  gera divergência falsa em TODO anúncio com variação (cenário negativo obrigatório).
- Sem vínculo produto↔anúncio → NÃO avaliar divergência (IC-02 — ausência de vínculo não é
  divergência).

## Done Means

- Anúncio real: row camada 2 com amount + detail canônico IC-01 COMPLETO (5 chaves:
  percentage_fee, fixed_fee, financing_add_on_fee, price_used, listing_type_id —
  F-r05-5) + origem `api_listing_prices` e `coletado_em`.
- Divergência plantada (estoque mirror ≠ ML) → row aberta; corrigida → auto-resolve
  (`resolved_at`) — 2 direções (IC-02).
- `filter.divergentes=true` retorna SÓ divergentes (fixture >1 página — lição CHIP-MERCADO).
- AnunciosTable: coluna TARIFA (amount formatado) + badge divergência; anúncio sem fee
  observada mostra honest-unknown (`—`), nunca 0.
- tsc + par OpenAPI+SDK no mesmo commit.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator (hub) — lane C após M-04 close.
- Next action: F-01 ∥ F-02.
- Required files/evidence: `validation-contract.md` (P6), `validation-result.md`.
- Blockers or open decisions: none.

## Correction Handoff

N/A.
