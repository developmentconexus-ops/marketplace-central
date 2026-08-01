# F-02-mass-closure-replacement

```yaml
id: F-02
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

Morte do MASS-CLOSURE (ADR-06): `ApplyCompletedPull`
(`listings/adapters/postgres/repository.go:383-465`) deixa de emitir o UPDATE incondicional
`status='closed'` (`:390-394`). Nova semântica: upsert por row (idioma writer.go:74-95 `upsertSQL`; keep-absent em
`:104-112` `keepAbsentSQL` — F-r06-5) +
`last_seen_at` no upsert; marcação `absent_since = now() WHERE last_seen_at <
run_started_at` executa SÓ quando o run é declarado COMPLETO (IC-06: enumerador exaurido +
hidratação drenada); run truncado por 429/deadline/cancel NUNCA marca. `status` = verdade
do provider verbatim, nunca inferido de ausência.

EARS:
- While run completo, when marcação roda, the repository shall setar absent_since só nas
  rows não vistas pelo run — e NUNCA mudar status.
- While run incompleto (qualquer causa), when pull termina, the repository shall preservar
  status e absent_since de todas as rows.
- While item reporta status closed no payload, when upsert, the row shall gravar closed
  (verdade do provider).

## Inputs

ADR-06/IC-06/IC-07 (binding); `research/codebase-ingest-side.md` (ApplyCompletedPull +
statement atual); 0090 do F-01 (last_seen_at/absent_since).

## Expected Output

`ApplyCompletedPull` re-semantizado (continua o writer ÚNICO — não criar segundo writer);
assinatura pode ganhar run metadata (run_started_at, completo:bool) — mudança contida no
package listings.

## Constraints

- ANTES do cursor (F-03) — edge R-B obrigatório; este feature NÃO implementa cursor.
- Nada de DELETE físico; `absent` ≠ `closed` (leitor distingue).

## Inputs/Outputs

- Input: batch de rows do run + run metadata (`run_started_at`, `completo:bool`) — regra de
  run-completo BINDING em IC-06 (enumerador exaurido + hidratação drenada; truncado NUNCA
  marca); colunas alvo BINDING em IC-07 E3 (`last_seen_at`, `absent_since` — 0090 do F-01).
- Output: rows upsertadas com `last_seen_at` avançado; `absent_since` setado SÓ em run
  completo nas rows não vistas; `status` = verdade do provider verbatim, nunca derivado de
  ausência. (Seção adicionada — auditoria P5 r06 F-r06-6.)

## Negative Scenarios

- Run completo com 0 itens (conta vazia / scan quebrado) → marcação NÃO roda em massa cega:
  pino de segurança — run com enumeração 0 ids é tratado como INCOMPLETO (nunca "tudo
  ausente"); teste nomeia.

## Ownership

- Owned paths: `apps/server_core/internal/modules/listings/adapters/postgres/repository.go`
  + testes do package.
- Forbidden paths: connectors (F-03), scheduler, migrações.
- Parallel-safe with: none — depends on F-01 (colunas).

## Validation Expectations

- Must-fail central (R-B): seed N listings > 1 página, abort pós-página-1 → asserção ZERO
  rows `status='closed'` flipped, com nome do teste declarando a falha protegida.
- Run completo com item ausente → absent_since setado, status intocado.
- Enumeração 0 ids → nenhuma marcação (pino de segurança).

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md` após F-01.
- Required files/evidence: `validation.md` com o must-fail vermelho nomeado.
- Blockers or open decisions: none.
