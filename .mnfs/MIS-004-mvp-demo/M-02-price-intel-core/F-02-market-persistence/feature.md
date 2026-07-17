# F-02-market-persistence

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-02
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-02-price-intel-core.

## Brief

Camada de persistência de evidência de mercado no módulo `market`: tabelas + repositórios para snapshots de preço, ofertas validadas, decisões de match, sinais competitivos e agregados (shapes IC-03). Regra ADR-17 de retenção: uma coleta FAILED NUNCA anula a última evidência VALID.

## Inputs

- IC-03 (`research/market-evidence-read-interface-contract.md`) — entidades, enums, regra de retenção.
- IC-01 — enums `match_status`, `price_evidence_status` (compartilhados).
- Bloco de migrations 0050–0054; fixture `runner_test.go`.
- Padrões de repositório existentes no repo (`modules/market/**` atual — observations/references SAT ficam intocados).

## Expected Output

- Migrations: `market_price_snapshots` (source `ml_sale_price|ml_price_to_win|ml_catalog_offers`, status `VALID|FAILED|EXPIRED`, amounts nullable, fetched_at, tenant_id, installation_id), `market_validated_offers`, `market_match_decisions` (âncoras citadas + resultado IC-01), `market_competitive_signals`, `market_aggregates` — todas tenant-scoped.
- Repositórios Go com escrita idempotente por (tenant, alvo, source, fetched_at) — recoleta no mesmo instante não duplica.
- Leituras "latest VALID por alvo+source" e "série p/ freshness".
- EARS: While existe snapshot VALID para um alvo, when nova coleta FALHA, the sistema shall gravar snapshot FAILED novo e manter o VALID anterior como latest-valid (consulta latest-valid retorna o antigo + idade). While TTL do IC-03 expira, when leitura latest-valid ocorre, the sistema shall marcar a evidência como EXPIRED na resposta (estado, não deleção).

## Negative Scenarios

- Insert com amount 0 vindo de "desconhecido" ⇒ REJEITADO na camada de domínio (unknown chega como null ou não chega — ADR-17).
- Snapshot FAILED sem causa registrada ⇒ inválido (causa obrigatória).
- Query sem tenant_id ⇒ impossível por construção (scope obrigatório no repo).

## Constraints

- Nenhum endpoint HTTP neste feature (F-04). Nenhuma chamada ML (F-01).
- `/market/observations|references` (SAT existente) intocados — tabelas novas, sem ALTER nas antigas.
- Migrations só no bloco 0050–0054.

## Ownership

- Owned paths: `apps/server_core/internal/modules/market/**` (persistência/domínio; NÃO handlers HTTP), `apps/server_core/migrations/0050*–0053*`, fixture `apps/server_core/internal/platform/migrate/runner_test.go` (bump).
- Forbidden paths: `modules/connectors/**` (F-01), handlers/rotas do market (F-04), OpenAPI, SDK.
- Parallel-safe with: F-01 (disjoint: módulos e migrations exclusivos).

## Validation Expectations

- Teste ADR-17 negativo: seed VALID → coleta FAILED → query latest-valid retorna o VALID antigo com `fetched_at` original + FAILED visível na série (transcript SQL/JSON exato).
- Teste idempotência: mesma (alvo, source, fetched_at) 2× ⇒ 1 linha.
- `go test` do módulo verde com `GOCACHE=.gocache`; migrations sobem em lane hermética (log).

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-02).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` com transcripts acima.
- Blockers or open decisions: none.
