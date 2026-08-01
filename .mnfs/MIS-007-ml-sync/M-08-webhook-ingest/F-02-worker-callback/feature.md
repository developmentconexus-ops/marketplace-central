# F-02-worker-callback

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-08
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-08 webhook-ingest.

## Brief

Worker in-process (goroutine no server_core — assunção ratificada): drena inbox FIFO
(`received_at` asc, status `received`), resolve user_id→installation via
`external_account_id` das installations (porta própria do package — F-r04-4), extrai o
order id do `resource` (`/orders/<id>`), chama `IngestOrder` (M-06 —
refetch AUTENTICADO; corpo da notificação NUNCA vira dado), transiciona statuses
(processing→done | attempts++→re-received | ≥5→dropped). + Runbook de registro do callback
no app ML (domínio ngrok fixo, topic orders_v2) — EXECUTADO SÓ com autorização explícita
do operador (write de config na conta ML).

EARS:
- While row received com user_id mapeável, when worker processa, the sistema shall chamar
  IngestOrder(resource id) e marcar done com processed_at.
- While IngestOrder falha, when attempts <5, the row shall voltar a received com attempts+1;
  when ≥5, shall virar dropped (visível IC-05).
- While notificação é forja com resource inexistente, when refetch autenticado retorna
  404/403, the worker shall marcar done com error nomeado e ZERO escrita de domínio.

## Inputs

IC-04 (statuses/fluxo binding); IC-06 (IngestOrder seam); mapa user_id→installation:
fonte = `integration_installations.external_account_id` (ML user_id persistido —
`auth_adapter.go:192,261` → `auth_flow_service.go:691`), lido pelo repo/serviço de
installations EXISTENTE (`installation_repo.go:81` ListInstallations já seleciona
`external_account_id`) atrás de PORTA do package webhook — NUNCA SQL cross-module direto;
sem store novo (auditoria P5 r04 F-r04-4). `AccessTokenResolver` (wiring root.go:370-378)
fica SÓ como fonte de token do refetch autenticado — não é invertível p/ lookup.

## Expected Output

Worker no package do M-08 + composição (goroutine no boot, shutdown limpo) + runbook
`callback-registration.md` no milestone root + implementação REAL de `WebhookStatsReader`
(porta IC-05 publicada pelo M-09): last_notification_at/pending/dropped_24h lidos do
inbox; injeção via `WithWebhookStatsReader(...)` chamado DA REGIÃO ANCORADA DO M-08 no
root.go — código/construção do M-09 NÃO é editado (impl default do M-09 fica como
fallback, só deixa de ser usada). Prova na ROTA: `GET /sync/health` pós-injeção retorna
stats do inbox — injeção por referência/ponteiro, nunca builder-cópia (IC-05 §seam; P5
r04 F-r04-2).

## Constraints

- Worker chama SÓ IngestOrder — zero lógica de pedido própria (ADR-04).
- Polling do inbox com intervalo curto (segundos) — spec pina valor; sem
  LISTEN/NOTIFY nesta missão (aditivo futuro).
- Shutdown: worker termina row em curso, nunca meio-processada sem status.

## Inputs/Outputs

- Input: rows `notifications_inbox` status=`received` FIFO (`received_at` asc) — máquina de
  status BINDING em IC-04 §Enums (`received`→`processing`→`done` | re-`received` c/
  attempts++ | ≥5 → `dropped`; terminais `malformed`/`unmapped`); shape de colunas IC-04
  §Fields verbatim.
- Output 1: chamada a `IngestOrder` (IC-03/IC-06) com resource id — corpo do webhook NUNCA
  vira dado de domínio (refetch autenticado).
- Output 2: impl REAL de `WebhookStatsReader` — shape de saída BINDING em IC-05
  §Required Outputs, bloco `webhook` (`last_notification_at`, `pending`, `dropped_24h`;
  F-r08-4 — "InboxHealth" é o nome da operação em IC-04 §Operations, não um heading do
  IC-05), consumido por `GET /sync/health`;
  injeção por referência/ponteiro via `WithWebhookStatsReader` (IC-05 §seam).
  (Seção adicionada — auditoria P5 r06 F-r06-6.)

## Negative Scenarios

- Resource não-order (`/items/...` etc. — só se topic passasse) → done+topic_ignored já
  cobre no F-01; worker ignora defensivamente com error nomeado.
- Installation revogada entre gravação e processamento → unmapped tardio (attempts caem
  nesse terminal, não loop).

## Ownership

- Owned paths: worker no package webhook/inbox; impl real de WebhookStatsReader (mesmo
  package); composição (mesma região ancorada do M-08).
- Forbidden paths: orders application (IngestOrder é interface importada); transport F-01.
- Parallel-safe with: none — depends on F-01.

## Validation Expectations

- Fixture end-to-end hermética: row received → done + efeito do IngestOrder (row de pedido
  atualizada) — provado contra o MESMO seam do scheduler (Q3).
- Replay com `_id` → zero row nova (contagem); replay SEM `_id` → rows extras processadas
  com ZERO efeito duplicado de domínio (IngestOrder idempotente — asserção no EFEITO, não
  na contagem de inbox; IC-04/N-4).
- Retry/dropped: falha injetada 5× → dropped.
- Live-drive (hub, pós-autorização): evento real → tela em segundos.

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md` após F-01.
- Required files/evidence: `validation.md`; runbook; autorização do operador p/ registro.
- Blockers or open decisions: autorização do callback (operador).
