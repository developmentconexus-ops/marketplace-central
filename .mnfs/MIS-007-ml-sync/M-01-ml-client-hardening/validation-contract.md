# Milestone Validation Contract — M-01-ml-client-hardening

```yaml
id: M-01-VC
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

Verdicts binários. Evidência = caminho inspecionável concreto. Tipos: `ran`, `assumed`,
`could-not-run` (nomear bloqueio). Seam real (API ML viva) NÃO é exigido pelos critérios
hermético-testáveis deste milestone; o comportamento live é certificado transitivamente
pelos live-drives de M-04/M-06 (backfills reais atravessam o decorator). Deferral
registrado: 2026-08-01, dono = live-drives das lanes B/C.

## Milestone ID

M-01

## QA Level

QA-0

## Required Outcome

Client ML resiliente no choke point único `doRawWithHeaders` (`capability_adapter.go:712`):
backoff exponencial + jitter + `Retry-After` honrado + token-bucket POR INSTALLATION entre
goroutines; multiget `/items?ids=` 20/batch; regra DTO `Raw json.RawMessage`. Zero schema,
zero UI, zero mudança para caminhos sem erro; `capability_adapter.go` CONGELA pós-close
(endpoint novo = arquivo novo).

## Criteria

## Criterion: Retry-After honrado com tempo NOMEADO
ID: M01-C1
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: `cd apps/server_core && GOCACHE=.gocache go test ./internal/modules/connectors/adapters/mercado_livre/... -run <TestRetryAfter>` (nome exato no spec)
- Expected: teste com servidor fake respondendo 429 `Retry-After: 2` PASSA asserindo tempo
  decorrido ≥2s entre tentativas (asserção sobre o RELÓGIO observado; "eventually succeeds"
  sem asserção de tempo = passe vácuo, REPROVA o critério)
- Actual:
- Artifact:
Blocking failure: retry imediato ignorando Retry-After, ou teste que passa sem nomear o
tempo decorrido
Blocking failure observed: No
Owner: QA Validator

## Criterion: Token-bucket por installation sob concorrência
ID: M01-C2
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: teste de concorrência (N goroutines contra o decorator, mesma installation)
- Expected: timestamps OBSERVADOS das requests respeitam o rate configurado (asserção sobre
  os timestamps coletados, não sobre a config); limite CONFIGURÁVEL, nunca constante
  compilada (fato #11 é `assumed`)
- Actual:
- Artifact:
Blocking failure: burst acima do bucket nos timestamps, ou limite hardcoded
Blocking failure observed: No
Owner: QA Validator

## Criterion: Budget esgotado → erro tipado nomeando o contexto
ID: M01-C3
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: teste com 429 persistente até esgotar tentativas
- Expected: `ErrCodeProviderRateLimited` retornado com nº de tentativas e último
  `Retry-After` na mensagem/campos — nunca erro seco genérico
- Actual:
- Artifact:
Blocking failure: erro genérico sem tentativas/Retry-After, ou retry infinito
Blocking failure observed: No
Owner: QA Validator

## Criterion: Multiget 20/batch com Raw populado
ID: M01-C4
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: teste do reader multiget com 45 ids
- Expected: exatamente 3 chamadas (20+20+5) observadas no fake; DTOs retornados com
  `Raw json.RawMessage` não-vazio por item
- Actual:
- Artifact:
Blocking failure: batching errado (1 chamada de 45, ou 45 de 1) ou Raw vazio
Blocking failure observed: No
Owner: QA Validator

## Criterion: Write marcado no-retry (opt-out provado)
ID: M01-C5
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: teste do caminho de write (`PUT /items`) sob 429
- Expected: ZERO retry automático em write (opt-out explícito no decorator) — write
  idempotência não é assumida; erro sobe tipado na 1ª falha
- Actual:
- Artifact:
Blocking failure: write re-tentado automaticamente
Blocking failure observed: No
Owner: QA Validator

## Criterion: Lanes verdes + freeze do adapter
ID: M01-C6
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `cd apps/server_core && GOCACHE=.gocache go test ./...` + lane de integração
  hermética; `git diff --stat <base>..<tip>` do chip
- Expected: lanes verdes; diff toca SÓ `connectors/adapters/mercado_livre/` (arquivos
  existentes: só capability_adapter.go; resto = arquivos novos), zero migração, zero UI
- Actual:
- Artifact:
Blocking failure: lane vermelha, ou escrita fora do write-set (Ownership M-01)
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Output de teste salvo (não afirmação); teste de tempo/timestamps cita os valores medidos.
- `git diff --stat` do chip salvo.
- Baseline PRÉ-MERGE capturada e salva ANTES do merge (P7 r01 B-6): resposta da operação
  real que atravessa o adapter ML (payload/headers relevantes) — irrecuperável depois;
  M01-U2 = diff das DUAS capturas salvas (pré vs pós), não afirmação de identidade.
- Scrub das capturas de baseline ANTES de qualquer commit (P7 r02 A-11, mesma regra do
  M-03): PII de comprador/entrega removida; headers de auth (`Authorization`/Bearer,
  tokens) NUNCA capturados — redigir no momento da captura, não depois.

## Blocking Failures

- Passe vácuo em teste de tempo (asserção sem relógio) = blocking (M01-C1/C2).
- Retry em write = blocking (M01-C5).
- Comportamento mudado em caminho sem erro (regressão nos testes existentes do adapter) =
  blocking.

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: n/a

## Handoff

- Current status: planned.
- Next owner: hub (lane A, ∥ M-02 ∥ M-09).
- Next action: dispatch F-01 ∥ F-02 pós P7 Ready.
- Required files/evidence: este arquivo; `M-01/validation-result.md`; outputs de teste.
- Blockers or open decisions: none.

## Critérios de user-drive (mandato do operador — obrigatório)

M-01 é infra invisível por design ("zero mudança de comportamento para caminhos sem erro") —
o user-drive prova exatamente ISSO.

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M01-U1 | Zero regressão visível com o decorator no fio: /anuncios, /pedidos, /precos, /integracoes carregam e operam sem erro novo no dev stack re-apontado pós-merge | browser drive 4 telas + console limpo (zero erro novo) |
| M01-U2 | Refresh/import existente que atravessa o adapter ML continua funcional (caminho sem erro inalterado) | browser drive de 1 operação real que toca o adapter + comparação das DUAS capturas salvas (baseline pré-merge de Evidence Requirements vs captura pós-merge) — B-6 |
