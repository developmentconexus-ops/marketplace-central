# Interface Contract

```yaml
id: IC-03
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

market (M-02, produtor) ↔ pricing (M-07), listings (M-05), FE Produto Detalhe (M-06), Dashboard (M-09).

## Why This Contract Exists

Evidência de preço tem significados distintos (nosso preço, vencedor, alvo competitivo, agregado de catálogo) que consumidores confundiriam; e o target P1c exige campos de evidência em TODA UI de preço.

## Resources Or Entities

- **MarketPriceSnapshot**: scope (`listing`|`catalog_product`), ref_id, source (`ml_sale_price`|`ml_price_to_win`|`ml_catalog_offers`), price (decimal nullable), currency `BRL`, observed_at, fetched_at, expires_at, status (`VALID`|`FAILED`|`EXPIRED`), failure_reason nullable, request_id.
- **CompetitiveSignal** (por anúncio próprio): listing_id, our_price, winner_price nullable, target_price nullable (= price_to_win), position nullable (`{rank, total}`), fetched_at.
- **MarketAggregate** (por codprod): median, min_valid, n_offers, n_sellers, computed_at, status (`OK`|`INSUFFICIENT_MARKET`|`NO_PRICE_EVIDENCE`). Regras: só BRL, condição new, identidade ACCEPT, dedupe por seller (menor oferta válida), n_sellers<5 ⇒ `INSUFFICIENT_MARKET`.
- **Verdict** (por codprod): match_status (IC-01), price_evidence_status, verdict_label nullable (`saudavel`|`viavel_preco_mercado`|`apertado`|`nao_vale`), blocking_state nullable (`NO_CANDIDATE`|`NO_PRICE_EVIDENCE`|`INSUFFICIENT_MARKET`|`SEM_CUSTO`), inputs_used. `SEM_CUSTO`: evidência de mercado existe mas custo ERP indisponível ⇒ faixa de preço de mercado é servida SEM verdict_label (P5 ruling 2026-07-17; consumido por M-02 F-04 e M-06).
- Nota de propriedade: `signal_status` (`OK|SEM_VINCULO|NO_PRICE_EVIDENCE|STALE`, M-05 F-01) é enum COMPOSTO derivado no módulo listings (dono: M-05) sobre estados deste contrato — NÃO é enum de source/status do IC-03 e não entra no "Must Not Decide" deste IC.

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| `POST /market/collections` | on-demand (runbook/UI) | `{codprod}` | **200 síncrono** com sumário `{status COMPLETED\|PARTIAL, decisões, contagens, causas}` | execução SÍNCRONA delimitada (timeout curto); sem tabela de job, sem polling — consumidor refaz GET após a resposta (P5 ruling r02, P5-F-11) |
| `GET /market/signals?listing_ids=` | M-05/M-06 | ids | CompetitiveSignal[] | ordem = input |
| `GET /market/aggregates?codprod=` | M-06/M-07 | codprod[] | MarketAggregate[] | ordem = input |
| `GET /market/verdicts?codprod=` | M-06/M-05/Simulador picker | codprod[] | Verdict[] | ordem = input |

## Evidence Fields (obrigação de UI — target P1c)

Toda exibição de preço de mercado carrega: `source`, `fetched_at` (⇒ idade via FreshnessIndicator), `n_offers`/`n_sellers` quando agregado, `match_status`. Valor ausente ⇒ `UnknownValue` ("—"), NUNCA 0.

## Persistence Expectations (ADR-17)

Coleta falha/null grava tentativa `FAILED` NOVA; último snapshot `VALID` permanece intacto e visivelmente envelhecido. `buy_box_winner` null ⇒ NÃO produz snapshot de preço (fica `NO_PRICE_EVIDENCE`). Expiração marca `EXPIRED`, não deleta.

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| coleta com flag catalog_offers OFF pedindo offers | 200 | — | responde com `NO_PRICE_EVIDENCE` + reason `FLAG_DISABLED` (fallback explícito, nunca 500) |
| provider 4xx/5xx durante coleta | 200 (coleta) | — | tentativa FAILED gravada + telemetria; leitura segue servindo último VALID |
| codprod sem identidade ACCEPT em aggregate | 200 | — | aggregate status segue match_status (não inventa) |
| coleta p/ codprod inexistente | 404 | — | produto não existe no catálogo (P7 r01, F-04 alinhado) |
| coleta concorrente p/ mesmo codprod | 409 | `COLLECTION_IN_PROGRESS` | segunda coleta síncrona sobre o mesmo produto enquanto a primeira executa |

## Must Not Decide In Feature Execution

Enums de source/status, regras de agregação, semântica price_to_win (alvo competitivo ML, NÃO "menor preço de mercado"), copy de veredicto (labels fixos acima).

## Validation Impact

Teste negativo ADR-17 (VC M-02); QA visual: FreshnessIndicator + fonte presentes em Anúncios, Produto Detalhe, Simulador, Dashboard.

## Amendment A1 — 2026-07-17 (hub-ratified, D-12; drafted by CHIP-M02 per D-04 flow)

F-04 public JSON freeze. Money = decimal string; null stays null (ADR-17).

**POST /market/collections 200 envelope:**
```json
{
  "status": "COMPLETED|PARTIAL",
  "decisões": [{"codprod": ..., "match_status": ..., "price_evidence_status": ..., "blocking_state": "...|null"}],
  "contagens": {"ok": 0, "no_price_evidence": 0, "insufficient_market": 0, "no_candidate": 0, "sem_custo": 0},
  "causas": [{"codprod": ..., "reason": "FLAG_DISABLED|PROVIDER_4XX|PROVIDER_5XX|NO_IDENTITY|TIMEOUT", "detail": "...|null"}]
}
```
- `decisões` in input order; `causas` = non-OK rows only; `contagens` rolls up `decisões`.
- RULING: `contagens` keys = lower snake_case (hub-approved) — wire keys decoupled from
  UPPER Go enum constants per D-04 "enum-keyed snake_case count maps".

**Verdict (per codprod):**
```json
{
  "match_status": ...,
  "price_evidence_status": ...,
  "verdict_label": null,
  "blocking_state": "NO_CANDIDATE|NO_PRICE_EVIDENCE|INSUFFICIENT_MARKET|SEM_CUSTO|null",
  "inputs_used": {"<input>": {"source": "ml_sale_price|ml_price_to_win|ml_catalog_offers|erp_cost", "as_of": "RFC3339"}},
  "market_range": {"min_valid": "dec|null", "median": "dec|null", "currency": "BRL", "n_offers": 0, "n_sellers": 0}
}
```
- `verdict_label` ALWAYS null from M-02 in MIS-004 (Q3/D-04: no margin label crosses the
  M-02 boundary; verde/âmbar = M-07-owned via IC-04 CalcProfile).
- `market_range` mirrors frozen MarketAggregate EXACTLY — NO max field; present whenever
  price evidence exists, including SEM_CUSTO.
- SEM_CUSTO cost input = consumer-side minimal port mirroring IC-02 GetCostAsOf,
  hub-wired at composition root post-merge (Q3).
- Aggregates product-vs-source-keyed: F-04-S1 internal decision; GET /market/aggregates
  stays product-keyed; exposed JSON unaffected (excluded from this amendment).
