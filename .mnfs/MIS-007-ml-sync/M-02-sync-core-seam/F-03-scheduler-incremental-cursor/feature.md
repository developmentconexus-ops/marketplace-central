# F-03-scheduler-incremental-cursor

```yaml
id: F-03
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

Duas mudanças pontuais na infra de sync (IC-06):
1. `incremental` deixa de ser hardcoded false (`scheduler.go:160`): passa a refletir o tipo
   do run DERIVADO do `cursor.phase` do cursor que o job retorna (backfill=false;
   incremental/sweep=true) → `last_incremental_at` vira real (pré-condição Q4/IC-05).
   `phase` é obrigatória SÓ nos jobs NOVOS (M-04/M-06 — escopo ratificado de ADR-07); o
   parse é TOLERANTE: phase ausente/desconhecida ⇒ `incremental=false` (comportamento de
   hoje) — o cursor legado de products (`ProductsCursor`, sem phase) fica intocado e nunca
   vira erro (mesmo pino de M-09 F-01; auditoria P5 r03 P-1).
2. Contrato de cursor executável: helper/teste de contrato que REJEITA job retornando cursor
   nil (IC-06: nil apaga — scheduler.go:42-45) e valida `phase` obrigatório no JSON —
   domínio do contrato = jobs NOVOS de M-04/M-06.
   Vocabulário: `backfill`/`incremental`/`sweep`.

EARS:
- While job reporta run incremental, when RecordSuccess, the sync_state shall gravar
  last_incremental_at (e não gravar em backfill).
- While job retorna cursor nil, when contrato roda, the teste shall falhar nomeando
  "terminal cursor must be non-nil".
- While job products existente roda, when fix aplicado, the comportamento de products shall
  permanecer idêntico (regressão).

## Inputs

IC-06 (binding: JobFunc signature intocada, RegisterJob fail-closed intocado, cadence-
agnostic Start intocado, `schedule jsonb` continua não-lido); `scheduler.go:42-46,105-161`;
`sync_state_repo.go:95-102`.

## Expected Output

Diff mínimo em `sync/application/scheduler.go`: o scheduler DERIVA o tipo do run do
`cursor.phase` quando presente (obrigatório SÓ nos jobs novos M-04/M-06 — ADR-07;
IC-06 §Compatibility Rules; mesma fonte que IC-05 usa); phase ausente ⇒ `incremental=false`
(parse tolerante — nunca erro, mesmo pino de M-09 F-01). `JobFunc` fica BYTE-IDÊNTICA (`scheduler.go:46` — Must Preserve
IC-06); NENHUM tipo de retorno muda; NENHUM job concreto editado (`products_job.go`
intocado — auditoria P5 r02 N-1). + teste de contrato de cursor reutilizável pelos jobs de
M-04/M-06.

## Constraints

- NENHUMA migração (enum 0075 já tem listings/orders — IC-06; milestone propondo migração
  p/ registrar job está ERRADO).
- Round-trip JSONB nunca byte-exact.
- Isolamento de falha por entidade preservado.

## Inputs/Outputs

In: cursor JSONB atual. Out: cursores com phase (exemplos canônicos IC-06 — binding).

## Negative Scenarios

- Cursor read error → pula o ciclo, NUNCA fabrica nil (comportamento existente preservado,
  teste pina).
- Cursor de job NOVO sem `phase` → contrato falha nomeando o campo (domínio do contrato =
  jobs de M-04/M-06).
- Cursor legado sem `phase` (products) → parse tolerante do scheduler: `incremental=false`,
  SEM erro (regressão pina — P5 r03 P-1).

## Ownership

- Owned paths: `apps/server_core/internal/modules/sync/application/scheduler.go` + testes;
  helper de contrato de cursor (arquivo novo no package sync).
- Forbidden paths: `sync_state_repo.go` (só leitura), migrações, jobs concretos.
- Parallel-safe with: F-01, F-04 (eixo files).

## Validation Expectations

- Teste: job fake reporta incremental → `last_incremental_at` NOT NULL; backfill → NULL
  mantido.
- Must-fail nil-cursor com a mensagem nomeada.
- Regressão products contra o job REAL: cursor `ProductsCursor` sem phase →
  `incremental=false` gravado, mesmo fluxo de antes, asserções nos MESMOS campos
  (falsificável no único job vivo — P5 r03 P-1).

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md`.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none.
