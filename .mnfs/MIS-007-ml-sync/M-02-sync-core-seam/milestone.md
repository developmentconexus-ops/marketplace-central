# M-02-sync-core-seam

```yaml
id: M-02
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-007
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-007 ml-sync — [mission.md](../mission.md); IC-01 (channel_fees), IC-02 (divergences),
IC-03 (order_shipments/orders ext), IC-06 (ports/cursor/scheduler).

## Outcome

Os shapes e ports compartilhados existem ANTES de qualquer produtor: migrações 0086-0089
(channel_fees, divergences, order_shipments, colunas aditivas de orders), ports de
fee/divergência com resolução e tolerâncias pinadas, fix do `incremental=false`, e o guard
allowlist encolhente dos 4 sítios read-time ML. Um dono, zero código de provider, zero UI.

## Why This Milestone Exists

Movimento estrutural do plano (P3 r01): extrair o que M-03/M-04/M-05/M-06/M-07 inventariam
incompativelmente (ADR-09/10, IC-01/02/03). É o que torna M-03 ∥ M-04 e a lane C seguras.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | core-ddl | [F-01-core-ddl/feature.md](F-01-core-ddl/feature.md) |
| F-02 | fee-divergence-ports | [F-02-fee-divergence-ports/feature.md](F-02-fee-divergence-ports/feature.md) |
| F-03 | scheduler-incremental-cursor | [F-03-scheduler-incremental-cursor/feature.md](F-03-scheduler-incremental-cursor/feature.md) |
| F-04 | read-guard-allowlist | [F-04-read-guard-allowlist/feature.md](F-04-read-guard-allowlist/feature.md) |

## Dependencies

Nenhuma (raiz, lane A — ∥ M-01 ∥ M-09).

## Ownership & Concurrency

- Exclusive surfaces: `apps/server_core/migrations/0086-0089_*`; packages novos de núcleo
  p/ fees/divergences sob `apps/server_core/internal/modules/` (ex.:
  `internal/modules/channelfees/`, `internal/modules/divergences/` — nomes finais no spec,
  layering AGENTS obriga `modules/`, DONO é este milestone);
  `apps/server_core/internal/modules/sync/application/scheduler.go` (fix pontual);
  arquivo novo do guard allowlist em teste arquitetural.
- Migration block: **0086-0089** (0086 channel_fees, 0087 divergences, 0088 order_shipments,
  0089 colunas aditivas orders+índice bucket).
- Predicted seam locks: none a conceder; CONCEDE depois — ports deste milestone são
  consumidos por M-03..M-08 (interfaces estáveis pós-close).
- Runs in parallel with: M-01, M-09.
- Internal feature DAG: F-01 → F-02; F-03 ∥ F-04 ∥ F-01.

## Risks

- Shape errado aqui propaga p/ 6 milestones — mitigação: shapes vêm PRONTOS de IC-01/02/03
  (o spec não re-decide, implementa).
- `scheduler.go` é compartilhado com o job `products` vivo — fix do incremental não pode
  mudar comportamento de products (teste de regressão).

## Done Means

- 4 migrações aplicam limpo + testes regex; PKs/uniques/CHECKs exatamente como IC-01/02/03.
- Ports com testes de contrato: resolução ledger-only camada 2→1→ausente-tipado com
  proveniência (degrau config = composição do M-07, fora do port); tolerância
  R$0.01; one-open-row upsert + auto-resolve; camada 3 recusa detail sem
  sale_fee_unit/quantity.
- Cursor de contrato: job retornando nil FALHA teste nomeando "terminal cursor must be
  non-nil" (IC-06).
- `incremental` reflete tipo do run (phase ausente ⇒ false — parse tolerante, ADR-07/P5
  r03 P-1); regressão verde contra o job products REAL (cursor sem phase →
  `incremental=false`, fluxo idêntico).
- Guard allowlist com 4 entradas passa na main atual; sítio novo simulado FALHA.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator (hub).
- Next action: despachar F-01 ∥ F-03 ∥ F-04; F-02 após F-01.
- Required files/evidence: `validation-contract.md` (P6), `validation-result.md`.
- Blockers or open decisions: none.

## Correction Handoff

N/A.
