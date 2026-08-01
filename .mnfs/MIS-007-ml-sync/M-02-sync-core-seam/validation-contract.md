# Milestone Validation Contract — M-02-sync-core-seam

```yaml
id: M-02-VC
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

Verdicts binários. Evidência = caminho inspecionável concreto. Tipos: `ran`, `assumed`,
`could-not-run`. Seams: banco hermético (lane de integração `GOCACHE=.gocache`, retry do
CREATE DATABASE) é a dependência real deste milestone — DDL/ports provados contra Postgres
REAL da lane, nunca contra mock de SQL.

## Milestone ID

M-02

## QA Level

QA-0

## Required Outcome

Shapes e ports compartilhados existem ANTES de qualquer produtor: migrações 0086-0089
(channel_fees, divergences, order_shipments, colunas aditivas orders) verbatim IC-01/02/03;
ports de fee/divergência com resolução e tolerâncias pinadas; fix `incremental`; guard
allowlist encolhente dos 4 sítios read-time ML. Zero código de provider, zero UI.

## Criteria

## Criterion: 4 migrações aplicam limpo, re-run no-op, shapes exatos
ID: M02-C1
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: `cmd/migrate` em banco hermético (2 rodadas) + testes regex das 4 migrações
  (estilo `migrations/listings_test.go:25` substrings + `:101` createTableBody regex)
- Expected: rodada 1 aplica 0086-0089; rodada 2 no-op; PKs/uniques/CHECKs byte-exatos aos
  Database Shapes de IC-01/IC-02/IC-03 (natural unique de channel_fees; partial unique
  WHERE resolved_at IS NULL; PK (tenant_id, provider, provider_shipment_id))
- Actual:
- Artifact:
Blocking failure: constraint divergente do IC, ou re-run que altera schema
Blocking failure observed: No
Owner: QA Validator

## Criterion: CHECKs rejeitam o inválido (provado por INSERT)
ID: M02-C2
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: teste de integração — INSERTs violando CHECKs
- Expected: row camada 3 amount sem currency REJEITADA (SQLSTATE de CHECK nomeado);
  2ª row aberta mesma (entidade,kind) REJEITADA pelo partial unique; com a 1ª resolvida,
  aceita (os DOIS lados)
- Actual:
- Artifact:
Blocking failure: INSERT inválido aceito, ou partial unique bloqueando após resolve
Blocking failure observed: No
Owner: QA Validator

## Criterion: Ports com contrato provado (resolução, tolerância, recusa)
ID: M02-C3
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: testes de contrato dos ports (channelfees/divergences)
- Expected: `ResolveListingFees` ledger-only camada 2→1→ausente-TIPADO com a tupla completa
  `{value, value_type, currency, layer, detail, origem, coletado_em}` (detail = jsonb
  VERBATIM da row; degrau config NÃO existe no port — é composição do M-07); tolerância
  R$0.01 aplicada no comparador; one-open-row upsert + auto-resolve; escrita camada 3
  fee_kind=commission RECUSA detail sem sale_fee_unit/quantity (freight aceita NULL —
  IC-01)
- Actual:
- Artifact:
Blocking failure: port devolvendo config, tupla incompleta (sem detail/origem), ou recusa
ausente
Blocking failure observed: No
Owner: QA Validator

## Criterion: Cursor terminal não-nil (must-fail nomeado)
ID: M02-C4
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: teste de contrato do cursor (IC-06)
- Expected: job retornando cursor nil FALHA com mensagem contendo "terminal cursor must be
  non-nil"
- Actual:
- Artifact:
Blocking failure: nil aceito como cursor terminal
Blocking failure observed: No
Owner: QA Validator

## Criterion: Fix incremental sem regressão em products
ID: M02-C5
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: teste do scheduler (`sync/application/scheduler.go`) + regressão contra o job
  products REAL
- Expected: `incremental` reflete o tipo do run (cursor com `phase` incremental → true);
  phase AUSENTE ⇒ false (parse tolerante ADR-07); job products com cursor sem phase →
  `incremental=false` e fluxo byte-idêntico ao atual
- Actual:
- Artifact:
Blocking failure: mudança de comportamento do job products vivo
Blocking failure observed: No
Owner: QA Validator

## Criterion: Guard allowlist — passa hoje, reprova sítio novo
ID: M02-C6
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: guard arquitetural na main atual + must-fail com sítio simulado
- Expected: allowlist com as 4 entradas (A/B orders enrich, C/D pricing) PASSA na main;
  sítio read-time ML novo simulado FALHA com o guard NOMEANDO o sítio (arquivo/símbolo) —
  não contagem genérica
- Actual:
- Artifact:
Blocking failure: guard vácuo (passa com sítio novo) ou que não nomeia o sítio
Blocking failure observed: No
Owner: QA Validator

## Criterion: Q2/R-6 — zero raw de PII no schema novo
ID: M02-C7
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: inspeção das 4 migrações — colunas de buyer fiscal (0089)
- Expected: colunas TIPADAS só (lista exata de p5-prerequisites §2); nenhuma coluna
  `raw`/jsonb de billing_info em tabela com dado fiscal/comprador
- Actual:
- Artifact:
Blocking failure: qualquer raw jsonb de payload de comprador no schema
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Outputs da lane hermética salvos (aquecer `.gomodcache` — alarme falso
  HPG_MIGRATION_FAILED conhecido; retry do CREATE DATABASE no 1º boot).
- Must-fail (M02-C4/C6) exige o output VERMELHO nomeando o alvo, depois o verde.

## Blocking Failures

- Shape divergente de IC = blocking (M02-C1/C2) — propaga p/ 6 milestones.
- Regressão no job products = blocking (M02-C5).
- Guard vácuo = blocking (M02-C6).

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: n/a

## Handoff

- Current status: planned.
- Next owner: hub (lane A).
- Next action: dispatch F-01 ∥ F-03 ∥ F-04; F-02 após F-01.
- Required files/evidence: este arquivo; `M-02/validation-result.md`.
- Blockers or open decisions: none.

## Critérios de user-drive (mandato do operador — obrigatório)

M-02 é seam puro (zero UI) — o user-drive prova que o schema/fix não quebrou o mundo vivo.

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M02-U1 | Pós-migração 0086-0089 no dev stack: /catalogo, /anuncios, /pedidos, /integracoes carregam sem erro novo (colunas aditivas não quebram reads existentes) | browser drive 4 telas + console limpo |
| M02-U2 | Job products vivo continua sincronizando pós-fix do scheduler: /integracoes (ou SELECT sync_state) mostra run de products RECENTE após o merge | browser drive /integracoes + SELECT sync_state products (last run > deploy time) |
