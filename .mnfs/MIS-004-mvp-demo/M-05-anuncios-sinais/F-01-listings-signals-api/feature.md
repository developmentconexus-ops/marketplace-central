# F-01-listings-signals-api

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-05-anuncios-sinais.

## Brief

Enriquecer a API de listings com sinais competitivos: campo aditivo `market_signal` por anúncio (via product_links→market, contrato IC-03), contadores de exceção no summary, sem quebrar nenhum consumidor atual de `/listings*`.

## Inputs

- IC-03 — `CompetitiveSignal` shape + estados de evidência (fonte única).
- OpenAPI atual `/listings*` (list, by-product, summary, refresh) — tudo preservado.
- Port Go `market.EvidenceReader` publicado por M-02 F-04 (batch: sinais por listing_ids + agregados/veredictos por `codprod`) — transporte FIXADO: consumo via port interno, NUNCA HTTP self-call; contrato de dados é IC-03.
- Vínculo listing→CODPROD via `product_links` (leitura por port/API pública do módulo — sem query cross-schema).

## Expected Output

- `GET /listings` (e by-product): cada item ganha `market_signal: {status, position, price_to_win, delta_pct, match_status, n_offers, n_sellers, evidence: {source, fetched_at, freshness}} | null` + `signal_status: OK|SEM_VINCULO|NO_PRICE_EVIDENCE|STALE` — aditivo, nullable. `match_status` (IC-01) e contagens de amostra (`n_offers`, `n_sellers` — IC-03) propagados do market, nunca omitidos quando existem.
- `GET /listings/summary` estendido: `exceptions: {sem_vinculo: n, abaixo_custo: n, sem_evidencia: n}` (abaixo_custo usa custo Reader IC-02; custo desconhecido NÃO conta como abaixo_custo).
- Filtros aditivos: `?exception=sem_vinculo|abaixo_custo|sem_evidencia`.
- EARS: While anúncio sem vínculo resolvido, when lista responde, the sistema shall enviar `market_signal: null` + `signal_status: SEM_VINCULO`. While evidência mais velha que TTL IC-03, when lista responde, the sistema shall enviar sinal com `signal_status: STALE` + idade (evidência velha visível, nunca escondida). While custo desconhecido, when exceção abaixo_custo computa, the sistema shall excluir o anúncio do contador (desconhecido ≠ violação).

## Inputs/Outputs

Shapes: IC-03 §CompetitiveSignal (referência). Codes: 200 sempre p/ lista (estados por item); filtro exception inválido ⇒ 422.

## Negative Scenarios

- Market indisponível (módulo erro) ⇒ lista responde 200 com `signal_status: NO_PRICE_EVIDENCE` + telemetria de erro — lista NUNCA cai por causa de sinal.
- Vínculo existe mas produto sem agregado/sinal ⇒ NO_PRICE_EVIDENCE (não SEM_VINCULO — estados distintos).

## Constraints

- ADITIVO estrito: nenhum campo/param existente removido/renomeado; consumidores atuais (AnunciosPage W1) continuam funcionando SEM edição.
- Leitura cross-módulo APENAS via ports/APIs públicas (market, product_links, Reader) — zero SQL nas tabelas alheias.
- Sem migration própria (projeção computada); precisar de cache ⇒ REQUEST reserva 0070+.

## Ownership

- Owned paths: `modules/listings/**`, seção `/listings*` OpenAPI (aditiva), `sdk-runtime/src/listings.ts` (novo).
- Forbidden paths: `modules/market/**`, `modules/product_links/**`, `modules/pricing/**`, barrel SDK, `apps/web/**` (F-02).
- Parallel-safe with: none — F-02 depende da seção OpenAPI deste F.

## Validation Expectations

- Transcript GET /listings com 4 anúncios seedados cobrindo os 4 `signal_status` ⇒ JSON exato por item (incl. null shape).
- Transcript summary: contadores exatos batendo com os itens da lista (mesma seed).
- Teste de regressão: response atual do W1 (snapshot) continua subset válido do novo (prova aditiva).
- Filtro `?exception=abaixo_custo` retorna só os corretos; item de custo desconhecido fora.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-05).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` com transcripts acima.
- Blockers or open decisions: none.
