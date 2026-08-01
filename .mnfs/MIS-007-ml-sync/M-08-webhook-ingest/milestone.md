# M-08-webhook-ingest

```yaml
id: M-08
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

MIS-007 ml-sync — [mission.md](../mission.md); IC-04 (webhook+inbox, binding integral),
IC-06 (worker chama IngestOrder).

## Outcome

Pedido novo aparece em segundos: `POST /webhooks/{provider}` público (always-200,
ponteiro-nunca-dado, IP log-only), `notifications_inbox` (0093), worker in-process que
chama o `IngestOrder` do M-06, callback registrado no app ML (ngrok domínio fixo), topic
`orders_v2` somente. Scheduler 5min permanece (reconciliação — nunca um sem o outro).

## Why This Milestone Exists

Última da cadeia de orders: o worker só tem o que chamar depois do M-06 (ADR-04 — construir
antes = worker sem caminho). 1ª superfície pública não-autenticada de escrita do sistema —
postura Q2 inteira pinada em IC-04.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | inbox-endpoint | [F-01-inbox-endpoint/feature.md](F-01-inbox-endpoint/feature.md) |
| F-02 | worker-callback | [F-02-worker-callback/feature.md](F-02-worker-callback/feature.md) |

## Dependencies

- M-06 (IngestOrder estendido — o worker chama exatamente esse seam).
- M-09 (porta `WebhookStatsReader` IC-05 + setter `WithWebhookStatsReader` — F-02 injeta a
  impl real; edge dura, na prática sempre satisfeita: M-09 é lane A).
- (transitivas: M-03, M-02, M-01.)

## Ownership & Concurrency

- Exclusive surfaces: package novo de webhook/inbox (transport + worker), migração 0093,
  wiring próprio no root.go (região ancorada), path `/webhooks/{provider}` no OpenAPI (SEM
  método SDK — IC-04).
- Migration block: **0093** (notifications_inbox).
- Predicted seam locks: consome IngestOrder (interface estável do M-06); IC-05 lê inbox p/
  saúde — shape de status é o de IC-04, não renegociável aqui; M-08 fornece a implementação
  REAL da porta `WebhookStatsReader` (publicada pelo M-09 com impl default = estado
  canônico inicial) e a injeta via `WithWebhookStatsReader` chamado na PRÓPRIA região
  ancorada (F-02) — código do M-09 NUNCA editado.
- Runs in parallel with: none (lane D; M-09 pode já ter fechado — seção webhook do
  /integracoes acende sozinha, IC-05).
- Internal feature DAG: F-01 → F-02.

## Risks

- Endpoint público: forja/flood — mitigação IC-04 (cap 64KB, dedupe, ponteiro-nunca-dado,
  classe interativa com deadline 15s).
- ngrok cai → webhook surdo (R-2): reconciliação 5min cobre; última notificação visível
  via IC-05.
- Registro do callback no app ML = mudança de config na conta ML → **autorização explícita
  do operador ANTES de registrar** (AGENTS: live ML writes gated; registro é write de
  configuração).

## Done Means

- Q2: forja plausível → 200 + row no inbox + ZERO escrita de domínio + IP off-allowlist
  logado.
- Live-drive: evento real → notificação → inbox → ingest → tela em segundos (design §9).
- Replay com `_id` presente → zero row nova (upsert attempts_provider). Sem `_id` NÃO há
  dedupe de transporte (IC-04/ADR-11 emendado — P5 r02 N-4): replay vira rows extras
  inofensivas; a prova de zero-duplicata é de DOMÍNIO — IngestOrder idempotente, zero
  efeito duplicado.
- Malformada/topic estranho/user_id desconhecido → 200 + status terminal correto
  (malformed/done+topic_ignored/unmapped).
- attempts ≥5 → dropped, visível na saúde (IC-05).

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator (hub) — despacho após M-06.
- Next action: F-01 → F-02; pedir autorização do operador p/ registro do callback.
- Required files/evidence: `validation-contract.md` (P6), `validation-result.md`.
- Blockers or open decisions: autorização de registro do callback (operador, na execução).

## Correction Handoff

N/A.
