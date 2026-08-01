# Milestone Validation Contract — M-08-webhook-ingest

```yaml
id: M-08-VC
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

Verdicts binários. Evidência = caminho inspecionável concreto. Seams: registro do callback
no app ML = WRITE de configuração na conta ML → **autorização explícita do operador ANTES,
registrada aqui** (AGENTS live-ML-writes). O drive de evento real depende desse registro;
todo o resto é hermético-testável (endpoint + inbox + worker com fixtures IC-04).

## Milestone ID

M-08

## QA Level

QA-0

## Required Outcome

Pedido novo em segundos: `POST /webhooks/{provider}` público (always-200,
ponteiro-nunca-dado, IP log-only, cap 64KB, classe interativa 15s), `notifications_inbox`
(0093), worker in-process chamando IngestOrder do M-06, callback registrado (ngrok fixo),
topic `orders_v2` somente. Scheduler 5min permanece como reconciliação.

## Criteria

## Criterion: Q2 — forja não injeta dado
ID: M08-C1
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: POST forjado plausível (user_id válido, resource inventado) + POST com
  `resource` fora do pin (ex.: `/users/123` e `/orders/../users/me` — traversal, P7 r01
  B-8) + par de POSTs p/ derivação de IP (header `X-Forwarded-For` com IP oficial e com IP
  não-oficial — controles positivo/negativo, P7 r01 B-9)
- Expected: 200 + row no inbox + ZERO escrita em tabela de domínio (o worker REFETCH na
  API autenticada e o recurso não existe → status terminal, sem dado); resource fora do
  pin `^/orders/[0-9]+$` → row `malformed`, NENHUMA URL ML construída com ele (IC-04 B-8);
  XFF oficial → `ip_official=true`, XFF não-oficial → `false` (derivação pinada IC-04 —
  critério não-constante), ambos LOGADOS e processados (log-only, decisão P1); corpo do
  webhook NUNCA vira dado de domínio (ADR-11)
- Actual:
- Artifact:
Blocking failure: forja alcançando tabela de domínio, ou corpo de webhook persistido como
dado
Blocking failure observed: No
Owner: QA Validator

## Criterion: Envelope do endpoint público (always-200, cap, deadline)
ID: M08-C2
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: testes do endpoint: payload válido / malformado / >64KB / topic estranho /
  user_id desconhecido; controle negativo de deadline FORA de `registerBatchRoutes`
- Expected: 200 corpo vazio em TODOS os casos processáveis (única exceção: falha de
  INSERT no inbox → 500 envelope apierror); >64KB → corpo TRUNCADO em 64KB com marcador
  (IC-04 — alinhado P7 r01), leitura bounded sem consumir o resto; classe
  interativa 15s (`route_deadline.go:23-28`) provada por trickle; malformado→`malformed`,
  topic estranho→`done`+topic_ignored, user_id desconhecido→`unmapped` (status terminais
  IC-04)
- Actual:
- Artifact:
Blocking failure: endpoint devolvendo erro a caller externo em caso processável, ou
payload gigante lido inteiro
Blocking failure observed: No
Owner: QA Validator

## Criterion: Dedupe de transporte e idempotência de DOMÍNIO
ID: M08-C3
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: replay com `_id` presente ×3; replay SEM `_id` ×3; SELECTs inbox + domínio
- Expected: com `_id` → 1 row de inbox (upsert attempts_provider++); sem `_id` → rows
  extras de inbox INOFENSIVAS, e a prova real é de DOMÍNIO: IngestOrder idempotente, zero
  efeito duplicado no pedido (ADR-11 emendado, P5 r02 N-4)
- Actual:
- Artifact:
Blocking failure: duplicata de DOMÍNIO em qualquer replay
Blocking failure observed: No
Owner: QA Validator

## Criterion: Retry → dropped visível
ID: M08-C4
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: fixture com ingest falhando persistente; observar máquina de status
- Expected: `received`→`processing`→re-`received` com attempts++ até ≥5 → `dropped`;
  dropped conta em `dropped_24h` do `GET /sync/health` (IC-05); worker em shutdown termina
  a row em curso (nunca meio-processada sem status)
- Actual:
- Artifact:
Blocking failure: retry infinito, ou dropped invisível na saúde
Blocking failure observed: No
Owner: QA Validator

## Criterion: Live — evento real na tela em segundos
ID: M08-C5
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive (PÓS autorização do registro): evento real na conta ML
  (ex. pergunta/alteração que dispare orders_v2 — na prática: pedido/mudança de pedido) →
  cronometrar até a tela
- Expected: notificação → inbox → worker → IngestOrder → /pedidos atualizado em SEGUNDOS
  (< cadência do scheduler; design §9); `WebhookStatsReader` real injetado — bloco webhook
  do health com `last_notification_at` real
- Actual:
- Artifact:
Blocking failure: evento real que só entra pela reconciliação 5min (webhook surdo)
Blocking failure observed: No
Owner: QA Validator

## Criterion: Autorização do operador registrada (gate de write ML)
ID: M08-C6
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: leitura do registro de autorização no artefato de execução do chip
- Expected: autorização EXPLÍCITA do operador, com data, ANTES do registro do callback no
  app ML; sem ela, M08-C5 fica `could-not-run` e o milestone NÃO passa por stub
- Actual:
- Artifact:
Blocking failure: callback registrado sem autorização registrada
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Fixture de forja com documento sintético (R-6: teste com PII sintética prova que raw não
  persiste).
- Cronometragem do live com timestamps (notificação recebida vs tela).
- Controle negativo do deadline FORA de registerBatchRoutes (classe medida por trickle).

## Blocking Failures

- Escrita de domínio a partir de corpo de webhook = blocking (M08-C1 — Q2 inteira).
- Duplicata de domínio = blocking (M08-C3 — Q3).
- Write ML sem autorização = blocking (M08-C6 — AGENTS).

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: n/a

## Handoff

- Current status: planned.
- Next owner: hub (lane D, após M-06; M-09 fecha na lane A antes do despacho da lane D —
  P7 r02 A-1: dependência de construção de lane, não fato consumado).
- Next action: F-01 → F-02; PEDIR autorização do operador p/ callback antes do live.
- Required files/evidence: este arquivo; `M-08/validation-result.md`.
- Blockers or open decisions: autorização de registro do callback (operador, na execução).

## Critérios de user-drive (mandato do operador — obrigatório)

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M08-U1 | Evento real no ML → pedido aparece/atualiza em /pedidos em segundos, cronometrado, sem F5 ritual além do refetch da tela | browser drive cronometrado + row do inbox correspondente |
| M08-U2 | /integracoes bloco webhook mostra "última notificação" com o timestamp do evento real (não mais estado inicial) | browser drive + SELECT inbox |
| M08-U3 | Derrubar o túnel (ngrok off) + evento real → pedido AINDA chega via reconciliação 5min; saúde mostra última notificação envelhecendo (R-2 na prática) | drive com túnel morto + sync_state/inbox comparados |
