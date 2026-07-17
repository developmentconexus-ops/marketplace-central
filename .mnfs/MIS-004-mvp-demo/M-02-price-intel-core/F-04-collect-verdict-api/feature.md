# F-04-collect-verdict-api

```yaml
id: F-04
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-02
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-02-price-intel-core.

## Brief

Orquestração de coleta on-demand + agregação + motor de veredicto + API HTTP do market. Fecha o M-02: `POST /market/collections` executa coleta SÍNCRONA delimitada para um produto (ou anúncio), grava evidência via F-02, resolve identidade via F-03, computa `MarketAggregate` e `Verdict`; GETs servem sinais/agregados/veredictos para M-06/M-07; M-05 consome via port Go `market.EvidenceReader` publicado por este F.

## Inputs

- IC-03 — shapes/enums/endpoints (fonte única): MarketPriceSnapshot, CompetitiveSignal, MarketAggregate, Verdict; labels `saudavel|viavel_preco_mercado|apertado|nao_vale`.
- F-01 ports (coleta), F-02 repos, F-03 resolver.
- Custo p/ veredicto: Reader `GetCostAsOf` (IC-02 — via port interno; sem custo ⇒ veredicto degrada honesto).

## Expected Output

- `POST /market/collections` `{codprod}` ⇒ execução SÍNCRONA delimitada (timeout curto, sem tabela de job, sem polling), **200** com sumário do resultado `{status COMPLETED|PARTIAL, decisões, contagens, causas}`; consumidor refaz GET após a resposta.
- Pipeline por produto: identidade catalog → `SearchCatalogByEAN` (sem EAN ⇒ pula p/ NO_CANDIDATE honesto) → F-03 Resolve → se ACCEPT: `GetCatalogProduct` (+ `ListCatalogOffers` se flag ON) → snapshots F-02; para anúncios vinculados: `GetOwnItemPricing` + `GetPriceToWin` → CompetitiveSignal.
- Agregação (IC-03): BRL, condição new, só matches ACCEPT, dedupe por seller, `<5` sellers ⇒ `INSUFFICIENT_MARKET`; flag OFF/sem ofertas ⇒ `NO_PRICE_EVIDENCE`.
- Veredicto: margem simulada nos quartis do agregado vs custo ERP ⇒ label + faixa de preço + evidência citada (fontes + fetched_at).
- `GET /market/signals?listing_ids=`, `GET /market/aggregates?codprod=`, `GET /market/verdicts?codprod=` — batch, shapes e chaves EXATOS IC-03 (param `codprod`).
- Port Go `market.EvidenceReader` (batch: sinais por listing_ids + agregados/veredictos por codprod) exportado pelo módulo — assinatura congelada no close deste F; M-05 F-01 consome read-only (nunca HTTP self-call).
- Seção OpenAPI `/market/*` (aditiva — observations/references intocados) + `sdk-runtime/src/market.ts`.
- EARS: While produto sem EAN, when coleta roda, the sistema shall registrar decisão NO_CANDIDATE e veredicto NO_PRICE_EVIDENCE sem chamar rotas de catálogo. While flag catalog_offers OFF, when agregação roda, the sistema shall responder NO_PRICE_EVIDENCE explícito (nunca agregado de amostra vazia). While custo ERP indisponível, when veredicto computa, the sistema shall responder faixa de preço de mercado com `verdict_label: null` + `blocking_state: SEM_CUSTO` (IC-03).

## Inputs/Outputs

Shapes exatos: IC-03 §Operations (referência). Status codes: POST **200 síncrono** (sumário no corpo); GET 200 com arrays possivelmente vazios; codprod inexistente ⇒ 404; coleta concorrente mesmo produto ⇒ 409 `COLLECTION_IN_PROGRESS`.

## Negative Scenarios

- ML fora do ar meio-coleta ⇒ snapshots FAILED gravados, latest-valid preservado (ADR-17), resposta 200 com `status: PARTIAL` + causas no sumário.
- `ErrRateLimited` ⇒ coleta aborta com estado inspecionável; sem retry-storm.
- Veredicto pedido p/ produto nunca coletado ⇒ 200 com `status: NO_PRICE_EVIDENCE` + `collected_at: null` (não 404 — produto existe).

## Constraints

- Zero writes ML. Coleta é on-demand APENAS (sem scheduler — MIS-005).
- Telemetria da rota flag inspecionável (IC-06).
- Enums EXATOS IC-01/IC-03 — divergência de string é defeito.

## Ownership

- Owned paths: `modules/market/**` (orquestração/handlers/veredicto), seção `/market/*` do OpenAPI (aditiva), `packages/sdk-runtime/src/market.ts`, `apps/server_core/migrations/0054*` se coluna extra necessária.
- Forbidden paths: `modules/connectors/**` (mudança de port = ESCALATION), `sdk-runtime/src/index.ts` (barrel = hub), demais módulos.
- Parallel-safe with: none — depends on F-01+F-02+F-03.

## Validation Expectations

- Transcript E2E (installation real, lane live-provider-read): POST collections p/ produto com EAN real ⇒ 200 síncrono com sumário; GET verdicts ⇒ label + faixa + ≥1 fonte com fetched_at; GET signals de anúncio próprio ⇒ posição price_to_win real.
- Teste de contrato do port `market.EvidenceReader`: batch com listing conhecido + codprod conhecido ⇒ shapes idênticos aos GETs correspondentes.
- Transcript flag OFF ⇒ verdicts com NO_PRICE_EVIDENCE (JSON exato).
- Transcript produto sem EAN ⇒ NO_CANDIDATE/NO_PRICE_EVIDENCE sem hit de catálogo (telemetria zero p/ rota).
- OpenAPI lint/SDK build verdes (contrato e client sincronizados).

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-02, após F-03).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` + logs lane.
- Blockers or open decisions: none.
