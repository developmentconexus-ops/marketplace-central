# F-04-scheduler-refresh-wiring

```yaml
id: F-04
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

Composição: instância de Scheduler DIÁRIA p/ listings (padrão
`synccomposition.NewProductsScheduler` — segunda instância, ADR-08), job do F-03 registrado
na entity `listings` (SEM migração — enum 0075 já tem); refresh manual em lote (rota batch
existente, 202 async) re-apontado pro caminho IngestListing — pull antigo morre de vez
(um caminho só, ADR-04).

EARS:
- While server boota, when composição roda, the scheduler diário shall registrar o job de
  listings e o sync_state shall refletir runs reais.
- While operador dispara refresh em lote na tela, when 202 aceito, the processamento shall
  correr pelo MESMO IngestListing do scheduler (batch class, nunca no request).

## Inputs

ADR-08/ADR-14 (constructor local, região ancorada no root.go); F-03 (job); rota de pull
manual existente (202 async — `research/codebase-ingest-side.md`).

## Expected Output

Constructor de composição local do módulo listings + 1 chamada no root.go; rota de refresh
existente redirecionada; caminho antigo de pull removido (não coexistem dois writers).

## Constraints

- root.go: SÓ a chamada do constructor na região ancorada (hub arbitra conflito — ADR-14).
- Cadência diário fixa; `schedule jsonb` não-lido.
- Sem contrato FE novo (rota de refresh já existe).

## Negative Scenarios

- Job já registrado p/ entity (RegisterJob fail-closed) → boot falha alto — teste cobre
  double-registration.
- Refresh durante backfill em curso → segundo run não sobrepõe: a rota EXISTENTE responde
  `409 refresh_in_progress` (202 quando livre) — comportamento PRESERVADO, medido em
  `research/codebase-read-side.md:87`; row correspondente na Error Matrix do IC-07 por
  referência (de-hedge P7 r02 A-12; sem code novo).

## Ownership

- Owned paths: composição do módulo listings (arquivo novo), linha ancorada no root.go,
  rota de refresh existente (re-point interno).
- Forbidden paths: scheduler.go; outros módulos; OpenAPI/SDK (sem mudança de contrato).
- Parallel-safe with: none — depends on F-03.

## Validation Expectations

- Boot com scheduler diário registrado (sync_state row de listings após run).
- Refresh em lote → 202 + rows atualizadas pelo caminho novo (log/efeito do IngestListing).
- Caminho antigo de pull: ausente (grep no módulo — writer único).

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md` após F-03.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none.
