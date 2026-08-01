# M-03-orders-shipment-persist

```yaml
id: M-03
type: milestone
status: passed
owner: Mission Strategist
parent: MIS-007
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-007 ml-sync — [mission.md](../mission.md); IC-03 (orders persistence, binding), IC-06
(IngestOrder seam), IC-01 (NÃO escreve camada 3 aqui — é M-06).

## Outcome

Pedido lido da tela vem 100% do Postgres: `IngestOrder` v1 (ordem completa + shipment +
comprador fiscal persistidos em uma passada), tabela `order_shipments` populada (DDL veio no
M-02/0088), colunas aditivas de `orders` preenchidas (0089: pack_id, provider_shipment_id,
bucket persistido + índice, date_last_updated_ml, fiscal TIPADO), e os 2 sites de leitura
viva ML do caminho de pedidos (A: shipment enrich, B: buyer fiscal enrich) MORTOS — enrich lê
banco. Allowlist de leitura viva (M-02 F-04) encolhe em 2 entradas no MESMO commit.

## Why This Milestone Exists

Hoje `GET /orders/{id}` faz 3-4 GETs vivos no ML por pedido
(`enrich_service.go:194-198`, readers ML em root.go:591-592 — `:590` é o cost reader
INTERNO e sobrevive; F-r08-3) — lento, caro de rate-limit, e mente
quando ML degrada. ADR-04: caminho único de ingest resource-addressed. M-06 (backfill 12m +
decomposição) ESTENDE o IngestOrder daqui — sem M-03 não há writer p/ estender.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | ml-ingest-readers | [F-01-ml-ingest-readers/feature.md](F-01-ml-ingest-readers/feature.md) |
| F-02 | ingest-order-v1 | [F-02-ingest-order-v1/feature.md](F-02-ingest-order-v1/feature.md) |
| F-03 | read-path-switch | [F-03-read-path-switch/feature.md](F-03-read-path-switch/feature.md) |

## Dependencies

- M-01 (client resiliente — ingest de 1 pedido = 3-4 GETs; 429 sem token-bucket mata backfill
  futuro; F-01 consome `doRawWithHeaders` decorado).
- M-02 (0088 order_shipments + 0089 orders aditivas aplicadas; guard/allowlist F-04 existe p/
  encolher).

## Ownership & Concurrency

- Exclusive surfaces: `orders/**` (application + transport + adapters do módulo), readers ML
  novos em `connectors/adapters/mercado_livre/` (ARQUIVOS NOVOS — não edita
  capability_adapter.go, que foi do M-01), par OpenAPI+SDK dos DTOs de /orders, região
  orders existente do root.go `:576-601` editada in-place (troca readers A/B `:591-592`
  por readers de banco, deleções inclusas — edita SÓ dentro da região; uma das DUAS
  exceções à regra linha-nova, a outra é M-07 região pricing `:828-858`; hub arbitra —
  F-r08-2), entradas A/B da allowlist (remoção).
- Migration block: nenhum (DDL veio no M-02).
- Predicted seam locks: `IngestOrder` vira interface estável — M-06 herda a posse da lane
  (nunca simultâneos); DeriveOrderBucket (`order_bucket.go:48`) REUSADO verbatim, não
  re-derivado.
- Runs in parallel with: M-04 (lane B — superfícies disjuntas: orders vs listings).
- Internal feature DAG: F-01 → F-02 → F-03 (readers antes do ingest; switch só depois que o
  banco tem o dado).

## Risks

- Enrich atual degrada honesto em 404 (`buyer_fiscal_reader.go:101-103`) — ingest TEM que
  preservar honest-absence (fiscal ausente = colunas NULL, nunca erro que bloqueia o pedido).
- `order_repo.go:378` chama `DeriveOrderBucket(providerStatus, "", tags, ...)` sem shipment —
  bucket persistido corrige esse braço; contagens de summary mudam p/ MELHOR (before/after
  obrigatório, não regressão silenciosa).
- FE /pedidos NÃO muda contrato aqui (margem é M-06) — DTO só ganha campos aditivos.

## Done Means

- `GET /orders/{id}` responde sem NENHUM GET vivo ML (log/contador de adapter no teste).
- Row de order_shipments existe p/ pedido com shipping_id; fiscal tipado persistido;
  bucket coluna == DeriveOrderBucket com shipment status REAL (não mais "").
- Allowlist: 2 entradas a menos; guard must-fail prova que ressuscitar site A quebra build.
- Import path existente (`/orders/import`, 202 batch) passa pelo IngestOrder — um writer só.
- Truth table `TestDeriveOrderBucket` intocada e verde.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator (hub) — despacho lane B junto com M-04.
- Next action: F-01 spec.
- Required files/evidence: `validation-contract.md` (P6), `validation-result.md`.
- Blockers or open decisions: none.

## Correction Handoff

N/A.
