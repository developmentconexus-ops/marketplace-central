# F-01-camada2-fee-ingest

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-05 listings-fees-divergence.

## Brief

Passo aditivo no pipeline de ingest de listings (M-04): p/ cada anúncio ingerido, observar
fee de venda camada 2 e gravar via `ChannelFeeWriter` (M-02, IC-01): subject_type=`listing`,
subject_id=provider_listing_id, layer=2, fee_kind=`commission`, value_type=`amount` +
`detail` {percentage_fee, fixed_fee, financing_add_on_fee, price_used, listing_type_id}
(canonical IC-01 — camada 2 NÃO é percent seco), origem=`api_listing_prices`,
`coletado_em`=hora da chamada. Fonte: UMA só — o reader que M-01 F-02 entrega pronto
(sale_price no multiget com `?context=channel_marketplace` OU o reader dedicado
`GET /items/{id}/prices` que M-01 F-02 entrega caso o multiget não exponha — decisão FECHADA
no M-01, este feature só consome). Upsert por natural key (re-observação avança coletado_em,
nunca row nova). Frete: NÃO observado aqui (honest-unknown IC-01 — sem chute de R$79).

EARS:
- While anúncio ingerido com fee observável, when passo roda, the sistema shall upsert row
  camada 2 com value amount + detail canônico IC-01 COMPLETO (5 chaves: percentage_fee,
  fixed_fee, financing_add_on_fee, price_used, listing_type_id — F-r05-5) + coletado_em.
- While payload não expõe fee (context ausente/campo nil), when passo roda, the sistema shall
  NÃO gravar (ausência honesta — nunca 0).
- While sweep re-observa fee igual, when upsert roda, the contagem de rows shall ficar
  estável e coletado_em avançar.

## Inputs

IC-01 (integral — enums/formatos binding); M-04 F-03 (pipeline IngestListing — ponto de
extensão aditivo); M-01 F-02 (fonte de fee PRONTA — multiget com context OU reader dedicado
/items/{id}/prices, verificado e entregue lá); memória `ml-catalog-offers-pricing-api`
(sale_price precisa `?context=channel_marketplace`).

## Expected Output

Passo de fee no ingest (arquivo novo em listings/application ou adapter — additive-only
sobre M-04) + fiação do ChannelFeeWriter no pipeline. A fiação (construção do writer e
injeção no pipeline) vive DENTRO do package de composition de listings, sob o lock aditivo
registrado do M-04 — `root.go` NÃO é tocado pelo M-05 (matriz da missão: célula root.go =
`—`; auditoria P5 r03 P-5).

## Constraints

- Lock aditivo M-04: NÃO alterar assinatura/semântica existente do IngestListing — só
  compor passo novo.
- Falha do passo de fee NÃO derruba ingest do anúncio (fee é enriquecimento; log + skip
  contado).
- SÓ camada 2 — layer CHECK do 0086 rejeita o resto de qualquer forma (defesa em
  profundidade).

## Inputs/Outputs

Row canonical: IC-01 §Canonical Examples camada 2 (binding — spec não re-decide shape).

## Negative Scenarios

- Anúncio closed/ausente no sweep → sem observação nova (rows antigas ficam; resolução de
  leitura usa coletado_em).
- 429 sob rajada → token-bucket segura; passo nunca estoura o run (medido no teste com
  clock fake do M-01).

## Ownership

- Owned paths: arquivo(s) novo(s) de passo de fee em listings (additive-lock M-04
  registrado); fiação no pipeline.
- Forbidden paths: schema channel_fees (M-02); capability_adapter.go; transport /listings
  (F-03).
- Parallel-safe with: F-02 (write-sets disjuntos: channel_fees vs divergences).

## Validation Expectations

- Fixture: ingest de N anúncios → N rows camada 2 verificadas por SELECT (amount + detail
  exatos do fixture nas 5 chaves: percentage_fee, fixed_fee, financing_add_on_fee,
  price_used, listing_type_id — asserção de subconjunto deixa writer que omite
  financing_add_on_fee passar verde; F-r06-4, lição chip-import-chain).
- Ausência de fee no payload → 0 rows + log nomeado.
- Re-run → mesma contagem, coletado_em maior.

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md`.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none — a dúvida multiget-vs-reader-dedicado é RESOLVIDA e
  entregue pelo M-01 F-02 (auditoria P5 F-9); este feature consome a fonte única pronta.
