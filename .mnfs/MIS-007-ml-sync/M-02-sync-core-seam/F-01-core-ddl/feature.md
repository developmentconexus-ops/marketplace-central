# F-01-core-ddl

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-02
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-02 sync-core-seam.

## Brief

Migrações 0086-0089, TODAS aditivas, shapes verbatim de IC-01/IC-02/IC-03:
- 0086 `channel_fees` (IC-01 Database Shape: natural unique + CHECKs layer/fee_kind/
  subject_type/value_type/origem + currency-when-amount).
- 0087 `divergences` (IC-02: partial unique WHERE resolved_at IS NULL + índice de contagem;
  timestamps dos 2 lados NOT NULL).
- 0088 `order_shipments` (IC-03: PK (tenant_id, provider, provider_shipment_id) + índice
  provider_order_id).
- 0089 colunas aditivas em `orders`/itens (IC-03: pack_id, provider_shipment_id, bucket +
  índice (tenant_id,bucket), date_last_updated_ml, buyer fiscal tipado, decomposition
  jsonb, net_amount, margin_pct) — NULL default (honest-unknown).

EARS: While banco na 0085, when `cmd/migrate` roda, the schema shall conter as 4 novas
superfícies com constraints exatas dos ICs; when re-rodado, shall ser no-op.

## Inputs

IC-01/02/03 (shapes binding); estilo de teste regex `migrations/listings_test.go:25`,
`product_link_decisions_test.go:23`; runner `runner.go:33-40` (filename = PK — NUNCA
renomear aplicada).

## Expected Output

4 arquivos de migração + 4 testes regex no estilo existente.

## Constraints

- Nada de ALTER destrutivo; nada em products_mirror/listings; 0021 duplicado NÃO se toca.
- Campos fiscais do 0089: colunas tipadas SÓ (raw de billing_info PROIBIDO — R-6);
  lista exata dos campos do drawer vem de `research/p5-prerequisites.md` §2 (buyer fiscal
  DTO) — o spec copia dali, não inventa.
- `cmd/server` não migra no boot — não mexer.

## Inputs/Outputs

In: schema atual 0085. Out: DDL SQL; sem superfície de runtime.

## Negative Scenarios

- Migração re-aplicada → no-op (IF NOT EXISTS / guards padrão do repo).
- INSERT violando CHECK (ex.: layer=4, amount sem currency) → erro de constraint (teste).

## Ownership

- Owned paths: `apps/server_core/migrations/0086_*..0089_*` + testes.
- Forbidden paths: migrações existentes; qualquer código Go de aplicação.
- Parallel-safe with: F-03, F-04 (eixo files/migração).

## Validation Expectations

- `cmd/migrate` verde em banco hermético; segunda rodada no-op.
- Teste de constraint: row camada 3 sem currency REJEITADA (SQLSTATE de CHECK).
- Partial unique provado: 2ª row aberta mesma (entidade,kind) REJEITADA; com a 1ª resolvida,
  aceita.

## Execution Artifact Rules

`spec.md`/`plan.md`/`validation.md` = execução.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md`.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none — p5-prerequisites §2 COMPLETO no repo
  (`buyer_fiscal_reader.go` + DTO enumerados; blocker descarregado, observação r07).
