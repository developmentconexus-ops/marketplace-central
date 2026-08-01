# Interface Contract — orders persistence (order_shipments + extensão de orders + decomposição)

```yaml
id: IC-03
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-007
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

Lado persistido de pedidos entre: M-03 (popula `order_shipments` + campos de shipment/fiscal,
mata sítios A/B), M-06 (backfill 12m + decomposição + camada 3), M-08 (webhook chama o mesmo
ingester), transport/FE de /pedidos (colunas ratificadas: SLA, bucket, frete_real, rastreio,
comprador_fiscal). DDL nasce no M-02.

## Why This Contract Exists

Dois milestones escrevem o lado orders em fases diferentes (M-03 shipment-enrich; M-06
backfill+decomposição); sem shape pinado, o M-06 re-inventa o que o M-03 criou. E a regra
A-05: as colunas de /pedidos NÃO podem piscar — o que a tela mostra hoje via read-vivo tem
que sair IGUAL do Postgres.

## Resources Or Entities

- Tabela `order_shipments` (nova, M-02).
- Colunas aditivas em `orders` e itens (M-02).
- `IngestOrder(ctx, tenant, installation, providerOrderID)` (ADR-04/IC-06) — writer ÚNICO.

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| IngestOrder | import existente, backfill, incremental 5min, webhook worker, refresh | provider_order_id | upsert orders + itens + order_shipments + camada 3 + decomposição + bucket | caminho ÚNICO; per-order GETs shipment/sla/costs/billing (multiget não existe — T5) com budget de goroutines |
| ListOrders (existente) | GET /orders | filtros atuais | DTO atual + campos novos | ZERO chamada ML (ADR-05); ordenação atual preservada |

## Fields

### Required Inputs — `order_shipments`

- PK `(tenant_id, provider, provider_shipment_id)`.
- Linkage: `provider_order_id text` NOT NULL (+ índice).
- Shipment: `status text`, `substatus text` NULL, `logistic_type text` NULL,
  `tracking_number text` NULL, `tracking_method text` NULL.
- SLA (fonte `/shipments/{id}/sla`): `sla_status text` NULL, `sla_limit_at timestamptz` NULL.
- Custos (fonte `/costs`, header `X-Costs-New`): `cost_gross numeric(14,2)` NULL,
  `cost_seller numeric(14,2)` NULL, `currency char(3)` NULL.
- Destino (o que a tela já mostra): `receiver_name text` NULL, `dest_city text` NULL,
  `dest_state text` NULL, `dest_zip text` NULL.
- SEM coluna `raw` (P7 r01 B-7, ADR-03 emendado: payload de shipment carrega PII de entrega
  — receiver name/endereço/CEP, classe PII de `cmd/mlprobe/main.go:41-43`; `raw jsonb` só em
  `listings`). Os campos tipados acima são a persistência COMPLETA de shipment.
- `source_time timestamptz` NULL; `fetched_at timestamptz` NOT NULL; created/updated.

### Required Inputs — extensão `orders` (aditivas, NULL = honesto-desconhecido)

- `pack_id text` NULL; `provider_shipment_id text` NULL (link).
- `bucket text` NULL + índice `(tenant_id, bucket)`. Valores = enum produzido pela função
  EXISTENTE `domain.DeriveOrderBucket` (`orders/domain/order_bucket.go:48` — assinatura e
  truth table INALTERADAS; já vive no núcleo). M-03 move só o CALL SITE: do read path p/ o
  INGEST (derivação nunca no read). Verificado por prerequisite-existence (§1 de
  `research/p5-prerequisites.md`; a prosa anterior "transport de market" / "função MOVE
  p/ o núcleo" era FALSA contra o repo — deletada, auditoria P5 r04 F-r04-3).
- `date_last_updated_ml timestamptz` NULL (watermark IC-06).
- Comprador fiscal (fonte billing_info, TIPADO — raw de billing_info PROIBIDO, ADR-03/R-6):
  `buyer_doc_type text` NULL, `buyer_doc_number text` NULL + exatamente os campos que o
  drawer comprador_fiscal renderiza hoje (P5 enumera do DTO de `buyer_fiscal_reader.go` —
  nem mais nem menos).
- Decomposição: `decomposition jsonb` NULL; extraídos p/ sort/filtro:
  `net_amount numeric(14,2)` NULL, `margin_pct numeric(7,4)` NULL.

### Required Outputs

DTO de /pedidos: campos atuais (SLA, bucket, frete_real, rastreio, destinatário,
comprador_fiscal) byte-compatíveis com o que o read-vivo entrega hoje, agora vindos do banco;
+ margem real quando decomposição presente. Termos de fee expostos carregam proveniência
(camada=3, origem, coletado_em) vinda da decomposição (ADR-09, P7 r01 B-3).

## Enums And Statuses

- `order_shipments.status/substatus`: verdade do provider verbatim (sem enum nosso; CHECK
  não aplica — vocabulário do ML muda).
- `bucket`: enum existente de `DeriveOrderBucket` (não inventar valor novo nesta missão).

## Error Cases

Rotas existentes mantêm envelope apierror vigente; nenhum code novo no lado orders.

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| — | — | — | sem rota nova; codes existentes de /orders inalterados |

## Persistence Expectations

- Upsert idioma canônico (writer.go:74-95 `upsertSQL`; keep-absent `:104-112`
  `keepAbsentSQL` — F-r06-5): ON CONFLICT na PK, keep-absent, nunca DELETE.
- Itens: conflict target existente `(…, mpc_line_id)`; `mpc_line_id` trigger-imutável
  (0033) — preservar.
- **Decomposição** (JSONB, `"versao":1`):
  - Chaves: `receita_bruta`, `comissao_total` (= Σ sale_fee_unit×qty das linhas, IC-01),
    `comissao_origem` (= `api_order`, vocabulário IC-01), `comissao_coletado_em`,
    `frete_seller` (= cost_seller do shipment), `frete_origem` (= `api_shipment`,
    vocabulário IC-01), `frete_coletado_em`, `custo_produto`, `custo_fonte`
    (ex.: `mirror:sankhya`), `custo_congelado_em`, `liquido`, `margem_pct`, `computado_em`,
    `incompleto` (array de chaves ausentes).
  - Proveniência de fee (ADR-09, P7 r01 B-3): todo termo de fee presente carrega origem +
    coletado_em (camada = 3, constante da decomposição); termo sem os DOIS companheiros =
    defeito de ingest. É ESTA proveniência que o PedidoDrawer renderiza (M-06/F-03) e que o
    drive do M-06 assere.
  - Campo desconhecido = AUSENTE do JSON (nunca 0) e listado em `incompleto`;
    `liquido`/`margem_pct` só computados com TODOS os insumos presentes.
  - **Custo congelado**: lido do mirror via vínculo NA PRIMEIRA computação; reingests
    posteriores atualizam fees/frete mas NUNCA `custo_produto`/`custo_congelado_em`.
  - `net_amount`/`margin_pct` colunas = espelho de `liquido`/`margem_pct` (mesma transação).

## Canonical Examples

```json
{"versao":1,"receita_bruta":239.70,
 "comissao_total":48.66,"comissao_origem":"api_order",
 "comissao_coletado_em":"2026-07-31T12:00:03Z",
 "frete_seller":22.90,"frete_origem":"api_shipment",
 "frete_coletado_em":"2026-07-31T12:00:04Z",
 "custo_produto":95.10,"custo_fonte":"mirror:sankhya",
 "custo_congelado_em":"2026-07-31T12:00:05Z","liquido":73.04,"margem_pct":30.47,
 "computado_em":"2026-07-31T12:00:05Z","incompleto":[]}
```

Parcial honesto (sem custo — produto não vinculado):

```json
{"versao":1,"receita_bruta":239.70,
 "comissao_total":48.66,"comissao_origem":"api_order",
 "comissao_coletado_em":"2026-07-31T12:00:03Z",
 "frete_seller":22.90,"frete_origem":"api_shipment",
 "frete_coletado_em":"2026-07-31T12:00:04Z",
 "computado_em":"2026-07-31T12:00:05Z","incompleto":["custo_produto"]}
```

(sem `liquido`, sem `margem_pct` — ausentes, não zero.)

## Database Shape

- `order_shipments` PK acima; índice por `provider_order_id`; migrações range M-02;
  teste regex estilo `listings_test.go:25`.
- `orders`/itens: SÓ colunas aditivas; zero ALTER destrutivo.

## Seed Data

Nenhum. Fixtures: pedido multi-linha qty>1 (R-4); pedido sem vínculo (decomposição parcial);
fixture >1 página p/ backfill (R-3).

## Timestamp And ID Semantics

- `fetched_at` = clock nosso; `source_time` = provider; `date_last_updated_ml` = verbatim do
  ML (base do watermark).
- ids provider verbatim; separadores compostos só nos formatos IC-01/IC-02.

## Compatibility Rules

- `decomposition.versao` incrementa em mudança de shape; leitor tolera versões anteriores.
- `order_shipments` extensível aditivamente (ex.: eventos de tracking futuro).

## Route Namespace

Nenhuma rota nova. Mudanças de DTO em /orders (campos novos) = par OpenAPI/SDK no milestone
que as expõe (ADR-14 emendado P5 r03 P-2: commits de contrato serializados — ≤1 COMMIT
FE em voo; código paraleliza).

## Transport And Integration

N/A.

## Must Preserve

- Colunas ratificadas de /pedidos NUNCA piscam (A-05): morte dos sítios A/B no MESMO merge
  que popula os substitutos, com before/after na mesma tela.
- `MissingSaleFee` continua honesto (IC-01).
- Classe de rota: /orders continua interativa (15s); ingest roda em batch/worker.

## Must Not Decide In Feature Execution

- Shape da decomposição; regra de congelamento; enum de bucket; colunas fiscais tipadas
  (nunca raw); formato de linkage shipment↔order.

## Validation Impact

- Live-drive: pedido real decomposto, margem confere com sale_fee POR UNIDADE (R-4 must-fail
  com qty>1).
- Before/after mesma tela: SLA/frete/rastreio/comprador_fiscal idênticos pré/pós morte dos
  sítios A/B; /pedidos <2s medido com dados reais.
- Kill-and-resume do backfill: zero duplicata (PK + upsert), retomada do watermark.
