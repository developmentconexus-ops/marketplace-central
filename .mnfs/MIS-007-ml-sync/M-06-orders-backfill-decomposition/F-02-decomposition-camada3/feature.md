# F-02-decomposition-camada3

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-06
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-06 orders-backfill-decomposition.

## Brief

Estende IngestOrder (M-03 — MESMO writer, ADR-04): pós-persistência do pedido, computa e
grava (a) rows camada 3 do fee ledger — uma por linha de pedido, subject_type=`order_line`,
subject_id=`<provider_order_id>:<provider_item_id>` (formato IC-01), value = TOTAL da linha
(sale_fee_unit × quantity — sale_fee é POR UNIDADE, fato T2), detail obrigatório
`{"sale_fee_unit":..,"quantity":..}`, origem=`api_order`; frete seller do shipment (M-03)
→ row fee_kind=`freight`, subject_type=`order`, subject_id=`<provider_order_id>`
(frete é do pedido, não da linha — IC-01), origem=`api_shipment`; (b) decomposição canônica
IC-03 no jsonb
de orders (versao:1; receita_bruta, comissao_total, frete_seller, custo_produto,
custo_fonte, custo_congelado_em, liquido, margem_pct, computado_em, incompleto[]) +
colunas net_amount/margin_pct. Custo do produto: lido do mirror ERP via vínculo no momento
da 1ª computação e CONGELADO (custo_congelado_em; recomputações futuras NÃO re-leem custo).
Campo indisponível → AUSENTE do JSON + nomeado em `incompleto[]`, margem NULL — nunca 0
(ADR-17 lado honest-unknown).

EARS:
- While pedido com N linhas + shipment + custo, when ingest completa, the sistema shall
  gravar N rows camada 3 commission + 1 freight + decomposição com liquido = receita_bruta
  − comissao_total − frete_seller − custo_produto (exemplo canônico IC-03:
  239.70−48.66−22.90−95.10=73.04).
- While custo ERP indisponível na 1ª computação, when decompõe, the JSON shall omitir
  custo_produto, listar em incompleto[], e margem_pct/margin_pct shall ser NULL.
- While pedido re-ingerido (incremental), when decomposição recomputa, the custo_produto
  shall permanecer o CONGELADO e camada 3 shall fazer upsert (natural key), nunca duplicar.

## Inputs

IC-01 (camada 3 + detail obrigatório, binding); IC-03 (JSON canônico verbatim — chaves e
regra ausente≠zero); M-03 (tx do ingest — decomposição entra na mesma tx ou tx subsequente
idempotente: spec decide e justifica); mirror/vínculo (leitura por porta, como M-05 F-02);
fato T2 (sale_fee POR UNIDADE).

## Expected Output

Passo de decomposição no ingest + escritas camada 3 via ChannelFeeWriter + repo estendido
(jsonb + colunas).

## Constraints

- comissao_total SEMPRE derivada de sale_fee_unit × quantity — nunca do valor agregado do
  payload (fonte única, evita dupla contagem).
- Pedido cancelado: decomposição computada mesmo assim (margem de cancelado = informação);
  bucket já discrimina na tela.
- Sem vínculo → custo ausente (incompleto[]), NUNCA bloqueia ingest.
- Soma-das-partes é invariante testada: liquido = receita_bruta − comissao_total −
  frete_seller − custo_produto (IC-03 canonical; tolerância centavo por arredondamento —
  regra pinada na spec, half-even). Qualquer parcela ausente → liquido/margem AUSENTES +
  incompleto[] — nunca soma parcial.

## Inputs/Outputs

JSON canonical IC-03 §decomposition (binding). Camada 3 canonical IC-01.

## Negative Scenarios

- Linha sem sale_fee no payload → comissão da linha ausente + incompleto[] nomeia
  `comissao:<item_id>` — total NÃO soma parcial silenciosa (ausência propaga).
- Shipment ausente (M-03 honest) → frete ausente + incompleto[].

## Ownership

- Owned paths: `orders/application/` (passo de decomposição), `orders/adapters/postgres/`
  (extensão), porta de custo (nova, read-only mirror).
- Forbidden paths: schema; scheduler; transport (F-03); pricing.
- Parallel-safe with: none — depends on F-01 (na prática estende o mesmo ingest; serial).

## Validation Expectations

- Golden do JSON canônico p/ pedido fixture completo E p/ pedido incompleto (2 goldens).
- Invariante soma-das-partes por property test simples (N fixtures).
- Congelamento: mudar custo no mirror + re-ingerir → decomposição INALTERADA.
- Camada 3: SELECT confere value TOTAL + detail contra fixture.

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md` após F-01.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none.
