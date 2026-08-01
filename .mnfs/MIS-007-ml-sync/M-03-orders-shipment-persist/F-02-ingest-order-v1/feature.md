# F-02-ingest-order-v1

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-03 orders-shipment-persist.

## Brief

`IngestOrder(ctx, tenant, installation, providerOrderID)` (assinatura IC-06 verbatim) em
`orders/application/` — o writer único de pedido: chama F-01
(OrderDetail → se ShippingID≠nil, ShipmentDetail → se billing precisa, fluxo fiscal
existente `ports/buyer_fiscal_reader.go:19-21`), persiste em UMA transação: header orders +
colunas 0089 (pack_id, provider_shipment_id, bucket, date_last_updated_ml, fiscal tipado) +
linhas/payments (reusa upsert existente `order_repo.go:467-514` estendido) + row
order_shipments (0088, upsert por PK `(tenant_id, provider, provider_shipment_id)`). Bucket
persistido = `domain.DeriveOrderBucket(providerStatus, shipmentDetail.Status, tags,
faturadoAt != nil)` — MESMA função (`order_bucket.go:48`), agora com shipment status REAL no
caminho de escrita. Duplicate protection herdada: header ganha só com
`provider_updated_at` bumped (comportamento atual `:603` preservado); shipment upsert
idempotente. Import path atual re-fiado: `ImportService.Import` (`import_service.go:46`)
passa a enumerar ids e chamar IngestOrder por id — `normalizeOrders` morre no mesmo diff
(um caminho só, ADR-04).

EARS:
- While pedido tem shipping_id, when IngestOrder roda, the sistema shall persistir orders +
  order_shipments + bucket na MESMA tx (rollback conjunto).
- While fiscal indisponível (404 honest-absence), when ingest persiste, the colunas fiscais
  shall ficar NULL e o pedido shall persistir normal.
- While replay do mesmo id sem mudança, when IngestOrder roda, the contagem de rows shall
  ficar estável (idempotência medida).

## Inputs

F-01 (readers); IC-06 (assinatura + regra nil-cursor N/A aqui — IngestOrder é por recurso);
IC-03 (colunas verbatim); `order_repo.go:467-514` (upsert baseline);
`DeriveOrderBucket` truth table (`order_bucket_test.go:8` — INTOCADA).

## Expected Output

Service IngestOrder + repo estendido + port `OrderIngestor` (interface que M-06 estende e
M-08 worker consome) + re-fiação do Import.

## Constraints

- Uma tx por pedido — nunca orders sem shipment quando shipment veio (atomicidade IC-03).
- 403 de pedido terceiro → skip contado, run NÃO falha (fato live).
- `FaturadoAt` é local (Sankhya) — ingest NUNCA sobrescreve (coluna fora do upsert set).
- Decomposição/camada 3: FORA (M-06). net_amount/margin_pct ficam NULL.

## Inputs/Outputs

Assinatura port: `IngestOrder(ctx context.Context, installationID, providerOrderID string)
error` (tenant vem do scoping padrão). Canonical: IC-06 §orders.

## Negative Scenarios

- Shipment 404 mas order ok → orders persiste, order_shipments ausente, provider_shipment_id
  preenchido (retry futuro acha).
- Order 403/404 → erro tipado, NENHUMA escrita.

## Ownership

- Owned paths: `orders/application/` (ingest service novo + import_service.go re-fiação),
  `orders/ports/` (port novo), `orders/adapters/postgres/order_repo.go` (extensão upsert +
  shipment repo novo), região orders root.go `:576-601` (fiação do ingest nos constructors).
- Forbidden paths: transport (F-03), connectors (F-01 fechou), scheduler (M-06).
- Parallel-safe with: none — depends on F-01.

## Validation Expectations

- Fixture end-to-end hermética: id → rows orders + order_shipments + bucket correto (caso com
  tag delivered E caso paid+faturado — 2 braços da truth table via shipment real).
- Replay → zero delta.
- Must-fail: quebrar atomicidade (commit parcial injetado) → teste nomeia.

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md` após F-01.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none.
