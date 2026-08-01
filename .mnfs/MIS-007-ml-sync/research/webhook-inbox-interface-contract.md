# Interface Contract — webhook + notifications_inbox

```yaml
id: IC-04
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-007
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

`POST /webhooks/{provider}` (chamado pelo ML, público, NÃO autenticado — 1ª superfície
assim no sistema) → tabela `notifications_inbox` → worker in-process → `IngestOrder`
(IC-03/IC-06). Tabela no M-02? NÃO — inbox nasce no M-08 (único produtor+consumidor,
seam não compartilhado); o contrato existe porque a POSTURA de segurança (gate P1) e o
shape do dedupe não podem drifar entre transport, worker e /integracoes (M-09 lê status).

## Why This Contract Exists

ADR-11 + gate P1: hint não-confiável, always-200, IP log-only. Worker e transport são
features distintas do M-08; M-09 lê o inbox p/ saúde. Sem pino, o transport devolve 4xx
(tempestade de 8 retries do ML) ou o worker confia no payload (injeção).

## Resources Or Entities

- Rota `POST /webhooks/{provider}` (`provider` = slug `mercado_livre`, mesmo vocabulário
  provider_code das 4 superfícies — CHIP-PED-FILA; P7 r02 ★2-B: provider_code real é
  `mercado_livre` — `market_adapters.go:239`, `p5-prerequisites.md:113`; `mercadolivre`
  é nome de PACKAGE Go, nunca slug de rota/coluna).
- Tabela `notifications_inbox` (nova, range do M-08).
- Worker in-process (goroutine no server_core), assunção ratificada em mission.md.

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| ReceiveNotification | POST do ML | body JSON (cap 64KB) | **200 vazio SEMPRE** — inclusive topic desconhecido, body malformado, user_id desconhecido | nunca 4xx/5xx p/ conteúdo (evita storm 8×/1h); grava inbox e retorna em ms; processamento NUNCA no request |
| DrainInbox | worker (loop curto) | rows `received` FIFO (`received_at` asc) | chama IngestOrder com resource id; marca status | user_id→installation via `integration_installations.external_account_id` lido por porta do package (F-r04-4); body nunca vira dado de domínio |
| InboxHealth | M-09 /integracoes | tenant | last_notification_at, pending, dropped_24h | ver IC-05 |

## Fields

### Required Inputs (colunas `notifications_inbox`)

- `id bigserial PK`; `provider text` NOT NULL.
- `tenant_id` NULL (resolvido via user_id→installation; NULL quando não mapeável).
- `installation_id` NULL (idem).
- `topic text` NULL; `resource text` NULL (ex.: `/orders/2000012345`). **Pin de formato
  (P7 r01 B-8):** `resource` só é usável como ponteiro quando casa `^/orders/[0-9]+$`
  (ancorado); fora do pin → status terminal `malformed` (200 mantido, row gravada). O
  `resource` é attacker-controlled — NUNCA concatenado em URL de chamada ML autenticada sem
  casar o pin; o guard mora no package webhook (dono M-08), nenhuma feature decide.
  `user_id text` NULL;
  `attempts_provider int` NULL (campo `attempts` do payload); `sent_at timestamptz` NULL;
  `notification_id text` NULL (campo `_id` do payload quando presente).
- `raw_body text` NOT NULL (cap 64KB, truncado com marcador — auditoria; NUNCA parseado
  para dado de domínio).
- `source_ip inet` NOT NULL; `ip_official boolean` NOT NULL (comparação à allowlist oficial
  54.88.218.97, 18.215.140.160, 18.213.114.129, 18.206.34.84 — LOG-ONLY, nunca reject;
  lista mutável). **Derivação pinada (P7 r01 B-9):** sob o túnel ngrok o peer do socket é o
  agente local — `source_ip` = primeiro IP do header `X-Forwarded-For` injetado pelo túnel;
  header ausente/inválido → peer do socket. Header é forjável fora do túnel ⇒ `ip_official`
  é INFORMATIVO/log-only por decisão P1 do operador, NUNCA gate de aceitação (registrado em
  mission.md §Non-Functional Scope). M08-C1 exige controle positivo (XFF oficial →
  `ip_official=true`) e negativo (XFF não-oficial → `false`) — critério não-constante.
- `status text` NOT NULL CHECK (`received`,`processing`,`done`,`malformed`,`unmapped`,
  `dropped`); `attempts int` NOT NULL DEFAULT 0 (NOSSAS tentativas de processamento);
  `error text` NULL; `received_at timestamptz` NOT NULL; `processed_at timestamptz` NULL.

### Required Outputs

Resposta HTTP: `200`, corpo vazio. Sem envelope de erro (não há erro externo).

## Enums And Statuses

- `received` → worker pega → `processing` → `done` | (falha) `attempts++`, re-`received` |
  attempts ≥ 5 → `dropped` (visível em IC-05).
- `malformed`: body não-JSON/sem resource OU `resource` fora do pin `^/orders/[0-9]+$`
  (B-8) — gravado e terminal (200 mesmo assim).
- `unmapped`: user_id sem installation nossa — gravado e terminal.

## Error Cases

Externos: nenhum (always-200). Internos: falha de INSERT no inbox = única condição que pode
gerar 5xx — aceita (rara) e coberta pela reconciliação 5min.

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| qualquer conteúdo inválido | 200 | — | gravado com status terminal; NUNCA 4xx |
| falha de persistência do inbox | 500 | envelope apierror vigente | única exceção; reconciliação cobre |

## Persistence Expectations

- Dedupe: `UNIQUE (provider, notification_id) WHERE notification_id IS NOT NULL`; duplicata
  → upsert `attempts_provider`/`received_at` (não cria row). Sem `notification_id`: sem
  dedupe de transport — idempotência real vem do ingest (ADR-04). (Chave estreitada da
  tupla original de ADR-11 — emenda registrada em `mission.md` ADR-11, P5 r02 N-4: tupla
  cheia com COALESCE bloquearia re-notificação legítima do mesmo resource p/ sempre.)
- Topic aceito p/ processamento: `orders_v2` SOMENTE (gate P1); outros topics = gravados
  `received`→`done` sem ação (auditoria) — pino: NÃO status novo, `done` com
  `error='topic_ignored'`.
- missed_feeds (`GET /missed_feeds`): NÃO consumido nesta missão — reconciliação 5min
  supera a janela; registrado como não-escopo.

## Canonical Examples

Payload ML (fato §5 design — dado NÃO vem):

```json
{"_id":"f9e1...","resource":"/orders/2000012345","user_id":"123456789",
 "topic":"orders_v2","application_id":"999","attempts":1,
 "sent":"2026-07-31T12:00:00.000Z","received":"2026-07-31T12:00:00.000Z"}
```

Forja (rejeição semântica, não HTTP): body plausível com resource alheio → 200, inbox row,
worker faz refetch AUTENTICADO → ML responde com o dado REAL ou 404/403 → zero escrita de
domínio com valor do body. Este é o critério Q2.

## Database Shape

- `notifications_inbox` no range de migração do M-08; índices: partial unique acima +
  `(status, received_at)` p/ drain.

## Seed Data

Nenhum. Fixture: notificação válida + forjada + malformada + user_id desconhecido.

## Timestamp And ID Semantics

- `sent_at` = do payload (UTC, provider); `received_at` = clock nosso.
- `notification_id` = `_id` verbatim ou NULL.

## Compatibility Rules

- Topic novo (`items`) = row aditiva no CHECK de processamento, missão futura; shape do
  inbox NÃO muda.
- Worker separável p/ processo próprio sem mudança de schema (assunção ratificada).

## Route Namespace

- `POST /webhooks/{provider}` — mount method-aware (`"POST /webhooks/{provider}"`), classe
  INTERATIVA, explicitamente FORA de `registerBatchRoutes` (root.go:259-272).
- OpenAPI: path documentado. SDK: **NENHUM método novo** (endpoint não tem consumidor FE) —
  registrado aqui p/ o gate same-commit não acusar par incompleto (decisão de planning,
  ADR-14).
- Registro do callback no app ML: domínio ngrok fixo do operador; troca = 1 edição no
  cadastro, zero código.

## Transport And Integration

- Sem cookie/CORS/auth — endpoint público por design; mitigação = conteúdo nunca confiável
  + refetch autenticado.
- Cap de body 64KB ANTES do parse (LimitReader).

## Must Preserve

- Always-200; ponteiro-nunca-dado; IP log-only; dedupe; classe interativa; scheduler 5min
  continua obrigatório (nunca um sem o outro).

## Must Not Decide In Feature Execution

- Postura de segurança; enum de status; regra de topic ignorado; cap de body; não-consumo
  de missed_feeds.

## Validation Impact

- Q2: forja → zero escrita de domínio + row com IP off-allowlist logado.
- Live-drive: evento real (mexer anúncio/pedido) → notificação → inbox → ingest → tela em
  segundos (design §9).
- Replay da MESMA notificação → zero duplicata (Q3, descarregado contra ADR-04).
