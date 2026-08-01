# F-01-ml-ingest-readers

```yaml
id: F-01
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

Readers de INGEST no adapter ML, arquivos novos em
`connectors/adapters/mercado_livre/` (ex.: `order_ingest_reader.go`,
`shipment_ingest_reader.go`): `GetOrderDetail(ctx, accountRef, providerOrderID)` (GET
/orders/{id} — payload COMPLETO: status, status_detail, tags, date_last_updated, pack_id,
shipping.id, buyer, payments, order_items com sale_fee POR UNIDADE — fato T2) e
`GetShipmentDetail(ctx, accountRef, shipmentID)` (GET /shipments/{id}: status, substatus,
logistic_type, tracking_number, sla/status.limit, costs gross/seller, campos de destino
TIPADOS — receiver name/cidade/UF/CEP conforme colunas IC-03; nome exato da chave no
payload verificado no spec contra fixture real, não assumido).
Ambos via `doRawWithHeaders` decorado (M-01) — token-bucket + retry herdam de graça. Decode
structs de provider FICAM no adapter (AGENTS); retorno = DTOs de `connectors/domain` novos
(`OrderDetail`, `ShipmentDetail`) com raw jsonb seletivo carregado junto (ADR-03:
billing_info raw NUNCA — reusar o fluxo de `buyer_fiscal_reader.go:59` p/ fiscal, que já
tipifica e descarta raw).

EARS:
- While installation autenticada, when IngestOrder pede detail, the reader shall retornar DTO
  tipado + raw seletivo em 1 GET.
- While ML responde 404/410 p/ shipment, when reader consulta, the retorno shall ser
  honest-absence (DTO vazio + found=false), nunca erro fatal.
- While billing_info presente, when raw é montado, the raw shall EXCLUIR billing_info
  (assert de ausência por teste).

## Inputs

M-01 (client decorado); `buyer_fiscal_reader.go:18-52` (idioma de decode structs);
`capability_adapter.go` INTOCADO; IC-03 §colunas (o que o DTO precisa carregar).

## Expected Output

2 arquivos novos de reader + DTOs em connectors/domain + testes com fixtures HTTP mock
(contrato de decode, não integração viva).

## Constraints

- Nenhum campo desconhecido vira zero — ausente = ponteiro nil (AGENTS).
- sale_fee é POR UNIDADE (fato live T2) — DTO nomeia `SaleFeeUnit`, nunca `SaleFeeTotal`.
- raw seletivo (campo DTO em memória, ADR-03): ordem = payload de /orders/{id} MENOS
  buyer.billing_info. PERSISTÊNCIA de raw no lado orders/shipments: NENHUMA (P7 r01 B-7 —
  ADR-03 emendado: `raw jsonb` só em `listings`; `order_shipments` SEM coluna raw, só
  campos tipados IC-03).

## Inputs/Outputs

DTO OrderDetail: campos p/ TODAS colunas 0089 + linhas/payments existentes
(`normalizeOrders` `import_service.go:108-162` enumera o baseline). ShipmentDetail: campos
p/ TODAS colunas 0088 (IC-03 verbatim).

## Negative Scenarios

- 403 (pedido de terceiro — fato live) → erro tipado não-retryable, ingest marca e segue.
- Payload sem shipping.id → OrderDetail.ShippingID nil; ingest NÃO chama shipment reader.

## Ownership

- Owned paths: arquivos NOVOS em `connectors/adapters/mercado_livre/`, DTOs novos em
  `connectors/domain/`.
- Forbidden paths: `capability_adapter.go` (M-01 fechou), orders module (F-02).
- Parallel-safe with: none — primeira do M-03.

## Validation Expectations

- Fixture de /orders/{id} real (scrubbed) decodifica TODAS colunas 0089 alvo.
- Assert: raw não contém substring `billing_info`.
- Assert: caminho de persistência de shipment não grava payload bruto (nenhuma coluna
  raw em `order_shipments` — B-7).
- 404 shipment → found=false sem erro.

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md`.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none.
