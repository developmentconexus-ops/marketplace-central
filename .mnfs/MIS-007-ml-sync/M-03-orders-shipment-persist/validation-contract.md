# Milestone Validation Contract — M-03-orders-shipment-persist

```yaml
id: M-03-VC
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-007
created: 2026-08-01
updated: 2026-08-01
validation_level: QA-0
lifecycle_scope: milestone
base_sha: dd89d4b3
```

Verdicts binários. Evidência = caminho inspecionável concreto. Seams: ingest contra API ML
REAL é live-driven (conta do operador, poucos pedidos — barato); testes herméticos provam
contrato com fixtures de payload real (dumps `docs/design/evidence/ml-api/`, PII scrub
antes de qualquer commit).

## Milestone ID

M-03

## QA Level

QA-0

## Required Outcome

Pedido lido da tela vem 100% do Postgres: `IngestOrder` v1 (ordem + shipment + comprador
fiscal em uma passada), `order_shipments` populada, colunas 0089 preenchidas, sites A/B de
leitura viva mortos (root.go:591-592 trocados por readers de banco), allowlist -2 no mesmo
commit.

## Criteria

## Criterion: GET /orders/{id} sem NENHUM GET vivo ML
ID: M03-C1
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: teste/log com contador de calls do adapter ML durante `GET /orders/{id}` de
  pedido ingerido
- Expected: contador ML = 0 no read; enrich lê `order_shipments` + colunas fiscais do banco
- Actual:
- Artifact:
Blocking failure: qualquer chamada ML disparada pelo read de detalhe
Blocking failure observed: No
Owner: QA Validator

## Criterion: Persistência completa em uma passada
ID: M03-C2
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: ingest de pedido real (live-drive) + SELECTs
- Expected: row em `order_shipments` p/ pedido com shipping_id (PK tenant/provider/
  provider_shipment_id); colunas fiscais tipadas preenchidas; `bucket` persistido ==
  `DeriveOrderBucket` com shipment status REAL (não mais "" — `order_repo.go:378`);
  pedido SEM fiscal (404 honesto) → colunas NULL, ingest NÃO falha
- Actual:
- Artifact:
Blocking failure: fiscal ausente bloqueando pedido, ou bucket derivado sem shipment status
Blocking failure observed: No
Owner: QA Validator

## Criterion: Allowlist encolhe -2 com must-fail
ID: M03-C3
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: guard M-02 F-04 pós-merge + must-fail ressuscitando site A
- Expected: entradas A/B removidas NO MESMO commit da troca dos readers; reintroduzir a
  chamada do site A quebra o guard NOMEANDO o sítio (output vermelho salvo)
- Actual:
- Artifact:
Blocking failure: allowlist com entrada morta sobrando, ou guard que não reprova o retorno
Blocking failure observed: No
Owner: QA Validator

## Criterion: Writer único — import path passa pelo IngestOrder
ID: M03-C4
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: leitura do caminho `/orders/import` (202 batch) pós-merge + teste
- Expected: import existente chama o MESMO `IngestOrder` (ADR-04 — um writer só); zero
  segundo caminho de escrita de orders
- Actual:
- Artifact:
Blocking failure: segundo writer de orders sobrevivendo
Blocking failure observed: No
Owner: QA Validator

## Criterion: Truth table intocada
ID: M03-C5
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `git diff <base>..<tip> -- '*order_bucket*'` + `GOCACHE=.gocache go test -run TestDeriveOrderBucket`
- Expected: `TestDeriveOrderBucket` byte-intocado e verde; `DeriveOrderBucket`
  (`order_bucket.go:48`) REUSADO, não re-derivado
- Actual:
- Artifact:
Blocking failure: truth table editada ou duplicada
Blocking failure observed: No
Owner: QA Validator

## Criterion: Q1 — detalhe de pedido rápido
ID: M03-C6
Level: Milestone
Type: Performance
Required: Yes
Status: Pending
Evidence:
- Command: browser QA — abrir drawer de pedido ingerido, medir (DevTools) 3 amostras
- Expected: resposta de `GET /orders/{id}` <2s em todas (era 3-4 GETs vivos; banco puro)
- Actual:
- Artifact:
Blocking failure: detalhe >2s ou request ML no waterfall
Blocking failure observed: No
Drive (UI — agent-browser; UI criteria only):
- Fixture: tenant real com ≥1 pedido ingerido pelo IngestOrder v1 (live-drive do hub)
- Steps:
  - open /pedidos
  - click <primeira linha da lista>
  - assert text "Comprador"
- Expected: drawer abre <2s com dados fiscais/shipment do banco
Owner: QA Validator

## Evidence Requirements

- Contador/log de adapter salvo (M03-C1); SELECTs salvos (M03-C2).
- Must-fail com output vermelho NOMEANDO teste (`failure_token=test=`) antes do verde.
- Before/after de contagens de summary (bucket muda p/ MELHOR — não regressão silenciosa).

## Blocking Failures

- GET vivo ML no read = blocking (M03-C1).
- Segundo writer = blocking (M03-C4).
- Fiscal ausente virando erro (honest-absence violado) = blocking (M03-C2).

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: n/a

## Handoff

- Current status: planned.
- Next owner: hub (lane B, ∥ M-04).
- Next action: F-01 → F-02 → F-03.
- Required files/evidence: este arquivo; `M-03/validation-result.md`.
- Blockers or open decisions: none.

## Critérios de user-drive (mandato do operador — obrigatório)

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M03-U1 | Drawer de pedido real mostra comprador fiscal + shipment vindos DO BANCO (matar rede ML no stack → drawer continua completo) | browser drive + SELECT de conferência (order_shipments + colunas fiscais) |
| M03-U2 | /pedidos lista e summary sem regressão visível pós-switch; contagens de bucket explicadas no before/after | browser drive lista + summary; before/after salvo |
| M03-U3 | Console/network do browser sem chamada ML e sem erro novo nas 2 telas dirigidas | network capture do drive (zero request p/ api.mercadolibre.com) |
