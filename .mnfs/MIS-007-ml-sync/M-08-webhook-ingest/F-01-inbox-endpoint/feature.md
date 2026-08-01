# F-01-inbox-endpoint

```yaml
id: F-01
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

Migração 0093 `notifications_inbox` (colunas/CHECK/partial-unique verbatim IC-04) +
`POST /webhooks/{provider}` transport fino: LimitReader 64KB ANTES do parse, extrai
(topic, resource, user_id, attempts, sent, _id) quando parseável, grava row com source_ip +
ip_official (allowlist IC-04, LOG-ONLY), responde **200 vazio SEMPRE** — inclusive
malformado, topic desconhecido, user_id desconhecido. Dedupe
UNIQUE(provider,notification_id) parcial → duplicata faz upsert de attempts_provider, não
row nova. Mount method-aware `"POST /webhooks/{provider}"`, classe INTERATIVA, fora de
`registerBatchRoutes`. OpenAPI: path documentado; SDK: NENHUM método (IC-04 — registrar a
decisão no diff p/ o gate same-commit).

EARS: verbatim IC-04 Operations/Enums (binding — spec não re-decide).

## Inputs

IC-04 (integral); slug provider_code das 4 superfícies (`mercado_livre` — P7 r02 ★2-B,
não confundir com package Go `mercadolivre`); route class idiom
(`route_deadline.go:23-28`); envelope apierror só p/ o 500 de falha de INSERT.

## Expected Output

Package novo (transport + repo do inbox) + 0093 + entrada OpenAPI + wiring ancorado.

## Constraints

- Processamento NUNCA no request (worker é F-02); handler só valida-grava-200.
- `raw_body` NUNCA parseado p/ dado de domínio; cap com marcador de truncamento.
- Nenhum 4xx por conteúdo (storm 8×/1h — fato #7, `research/external-ml-api-facts.md`).

## Inputs/Outputs

IC-04 Canonical Examples (binding). Response: 200 corpo vazio; única exceção 500 envelope
vigente (falha de INSERT).

## Negative Scenarios

IC-04: malformado→`malformed`; user_id sem installation→`unmapped`; topic≠orders_v2→
`done`+`error='topic_ignored'`; body >64KB→truncado com marcador; duplicata→upsert.

## Ownership

- Owned paths: package novo webhook/inbox, `migrations/0093_*`, entrada OpenAPI do path,
  linha ancorada root.go.
- Forbidden paths: SDK client literal (sem método — só se o gate exigir nota); orders
  ingest (F-02 consome).
- Parallel-safe with: none — primeira do M-08.

## Validation Expectations

- 4 fixtures (válida/forjada/malformada/unmapped) → todas 200; statuses exatos por row;
  ip_official=false logado na forjada.
- Duplicata → contagem de rows estável, attempts_provider incrementado.
- Rota responde <15s deadline interativo (classe provada por teste de registro — lição
  deadline-class-measured-by-trickle: controle negativo FORA de registerBatchRoutes).

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md`.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none.
