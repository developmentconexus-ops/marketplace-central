# F-01-listings-ddl

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-04
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-04 listings-backfill-ingest.

## Brief

Migrações 0090-0092, aditivas, shapes verbatim IC-07:
- 0090: colunas E3 em `listings` (sold_quantity, category_id, condition, permalink,
  thumbnail, date_created_ml, tags text[], catalog_product_id, shipping_mode,
  free_shipping, logistic_type, available_quantity) + lifecycle (last_seen_at,
  absent_since) + raw jsonb + raw_truncated boolean. SEM colunas de fee (decisão IC-07 —
  commission_amount/pct/free_shipping_cost EXCLUÍDAS; fee mora em channel_fees).
- 0091: `listing_variations` PK (tenant_id, installation_id, provider, provider_listing_id,
  variation_id) — tuple verbatim ADR-13/IC-07 (P7 r02 ★2-A);
  price, available_quantity, sold_quantity, seller_sku, attributes jsonb, last_seen_at,
  absent_since.
- 0092: índice `(tenant_id, provider_listing_id)` em variations + o que IC-07 pinar.

EARS: While banco pós-0089, when migrate, the schema shall ganhar as superfícies acima;
re-run = no-op; PK de `listings` INALTERADA (sentinela `'-'` — ADR-13).

## Inputs

IC-07 (binding); 0036/0037 (shape atual de listings); estilo regex test.

## Expected Output

3 migrações + testes regex.

## Constraints

- PK de listings NÃO muda; 0022/0025 intocadas; NULL default em tudo (honest-unknown).
- Range 0090-0092 EXATO (0086-0089 são do M-02 — colisão = R-7).

## Inputs/Outputs

In: schema pós-0089 (0036/0037 shape atual de listings). Out: DDL verbatim IC-07 §listings
E3 + §listing_variations (colunas, tipos, PKs, defaults NULL) — spec não re-decide shape;
divergência de coluna vs IC-07 = defeito de plano, não escolha de spec.

## Negative Scenarios

- Re-run no-op; INSERT variation sem listing pai → sem FK dura? IC-07 não pina FK — spec
  decide índice sem FK (padrão do repo) e DOCUMENTA; não inventar cascade.

## Ownership

- Owned paths: `apps/server_core/migrations/0090_*..0092_*` + testes.
- Forbidden paths: migrações < 0090; código de aplicação.
- Parallel-safe with: none — primeira da cadeia do M-04.

## Validation Expectations

- Migrate verde + no-op na 2ª rodada; regex tests; colunas de fee AUSENTES assertadas
  (teste nega a presença — guarda a decisão IC-07).

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md`.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none.
