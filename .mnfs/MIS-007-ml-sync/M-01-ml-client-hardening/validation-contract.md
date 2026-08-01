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
Status: **Passed**
Evidence:
- Command: `cd apps/server_core && GOCACHE=.gocache go test ./internal/modules/connectors/adapters/mercado_livre/... -run TestResilienceDecoratorHonorsRetryAfterHeader -v -count=1`
- Expected: teste com servidor fake respondendo 429 `Retry-After: 2` PASSA asserindo tempo
  decorrido ≥2s entre tentativas (asserção sobre o RELÓGIO observado; "eventually succeeds"
  sem asserção de tempo = passe vácuo, REPROVA o critério)
- Actual: PASS (2.39s); log `measured total elapsed = 2.3922754s, inter-call gap = 2.3917459s`;
  reconfirmado por reviewer adversarial independente em run separado (2.17s). Não vácuo —
  asserção sobre `time.Now()` medido, não "eventually succeeds".
- Artifact: `_chip-m01/EVIDENCE.md` §M01-C1; `resilience_decorator_test.go:61-114`.
Blocking failure: retry imediato ignorando Retry-After, ou teste que passa sem nomear o
tempo decorrido
Blocking failure observed: No
Owner: QA Validator

## Criterion: Token-bucket por installation sob concorrência
ID: M01-C2
Level: Milestone
Type: Functional
Required: Yes
Status: **Passed**
Evidence:
- Command: `... -run TestResilienceDecoratorTokenBucketThrottlesConcurrentRequests -v -count=1`
- Expected: timestamps OBSERVADOS das requests respeitam o rate configurado (asserção sobre
  os timestamps coletados, não sobre a config); limite CONFIGURÁVEL, nunca constante
  compilada (fato #11 é `assumed`)
- Actual: PASS (0.40s); log `gap[1..4]≈99.6-100.2ms, total elapsed for 5 concurrent requests
  at 600/min = 400.4ms (min expected ~400ms)` — 5 goroutines reais, asserção sobre timestamps
  observados. `RateLimitPerMinute` confirmado como campo real wired em
  `newResilienceDecorator` (não constante); singleton do adapter confirmado em
  `internal/composition/root.go:371` (bucket compartilhado no processo, não por-request).
  Debt registrado (não-blocking): composition root ainda não passa valor não-default de
  env/config — mecanismo é configurável, falta knob operacional (ver EVIDENCE.md "Registered
  debts").
- Artifact: `_chip-m01/EVIDENCE.md` §M01-C2; `resilience_decorator_test.go:121-189`;
  `resilience_decorator.go:86-132`.
Blocking failure: burst acima do bucket nos timestamps, ou limite hardcoded
Blocking failure observed: No
Owner: QA Validator

## Criterion: Budget esgotado → erro tipado nomeando o contexto
ID: M01-C3
Level: Milestone
Type: Functional
Required: Yes
Status: **Passed**
Evidence:
- Command: `... -run TestResilienceDecoratorBudgetExhaustedReturnsTypedError -v -count=1`
- Expected: `ErrCodeProviderRateLimited` retornado com nº de tentativas e último
  `Retry-After` na mensagem/campos — nunca erro seco genérico
- Actual: PASS (0.13s); `MaxRetryAttempts=3`; `domain.ErrorCodeOf(err)==ErrCodeProviderRateLimited`;
  mensagem contém "3 attempts"; `errors.As` recupera `*RateLimitExhaustedError{Attempts:3}`,
  `LastRetryAfter>0`; exatamente 3 chamadas reais observadas no fake.
- Artifact: `_chip-m01/EVIDENCE.md` §M01-C3; `resilience_decorator.go:57-64,255-286`;
  `resilience_decorator_test.go:204-252`.
Blocking failure: erro genérico sem tentativas/Retry-After, ou retry infinito
Blocking failure observed: No
Owner: QA Validator

## Criterion: Multiget 20/batch com Raw populado
ID: M01-C4
Level: Milestone
Type: Functional
Required: Yes
Status: **Passed** (com amendment sobre o cap de truncamento — ver abaixo)
Evidence:
- Command: `... -run TestGetItemsMultigetPartitionsInBatchesOf20 -v -count=1`
- Expected: exatamente 3 chamadas (20+20+5) observadas no fake; DTOs retornados com
  `Raw json.RawMessage` não-vazio por item
- Actual: PASS; 45 ids → 3 chamadas, batches `[20,20,5]`, ordem global preservada; Raw
  não-vazio e byte-idêntico ao sub-objeto do próprio item (não o array inteiro) confirmado em
  `TestGetItemsMultigetSuccessfulDTOsHaveNonEmptyRaw`. Isolamento por item testado com batch
  MISTO real (`TestGetItemsMultigetPerItemErrorDoesNotFailBatch`: item 17 falha 404, outros 44
  OK), não a versão fraca all-fail/all-succeed. Achado de review adversarial: cap de Raw
  por-item implementado é 256KiB, não 1MB como `feature.md`/ADR-03 diziam — classificado e
  resolvido via amendment em `F-02-items-multiget-raw-dto/feature.md` (razão: cap externo de
  1MB é da RESPOSTA INTEIRA, compartilhado por até 20 itens; 1MB por-item nunca seria
  alcançável). Comentário e teste corrigidos, constante pinada (`itemMultigetRawCap==256*1024`
  explícito no teste).
- Artifact: `_chip-m01/EVIDENCE.md` §M01-C4; `items_multiget_reader.go:169-343`;
  `items_multiget_reader_test.go`.
Blocking failure: batching errado (1 chamada de 45, ou 45 de 1) ou Raw vazio
Blocking failure observed: No
Owner: QA Validator

## Criterion: Write marcado no-retry (opt-out provado)
ID: M01-C5
Level: Milestone
Type: Engineering
Required: Yes
Status: **Passed**
Evidence:
- Command: `... -run TestResilienceDecoratorStockWriteDoesNotRetry -v -count=1`
- Expected: ZERO retry automático em write (opt-out explícito no decorator) — write
  idempotência não é assumida; erro sobe tipado na 1ª falha
- Actual: PASS; 1 chamada PUT, `elapsed=67ms` contra `RetryBaseDelay=500ms` (escolhido para
  pegar retry acidental por tempo). Garantia ESTRUTURAL, não sorte de teste: `doRawWithHeaders`
  roteia qualquer `method != GET` para `doRawWithHeadersNoRetry` ANTES do loop de retry — os
  3 call sites de write (stock, price, listing) convergem aqui, todos PUT. Regressão
  confirmada: `price_writer_test.go`/`listing_writer_test.go` (subtestes 429, zero diff) ainda
  passam sem qualquer edição.
- Artifact: `_chip-m01/EVIDENCE.md` §M01-C5; `capability_adapter.go:745-767`.
Blocking failure: write re-tentado automaticamente
Blocking failure observed: No
Owner: QA Validator

## Criterion: Lanes verdes + freeze do adapter
ID: M01-C6
Level: Milestone
Type: Engineering
Required: Yes
Status: **Passed**
Evidence:
- Command: `cd apps/server_core && GOCACHE=.gocache go build ./... && GOCACHE=.gocache go vet ./...`;
  `GOCACHE=.gocache GOMODCACHE=.gomodcache go test ./internal/modules/connectors/... -count=1`;
  `npm run harness:integration` (lane hermética, repo root); `git diff --stat HEAD -- internal/modules/connectors/adapters/mercado_livre/`
- Expected: lanes verdes; diff toca SÓ `connectors/adapters/mercado_livre/` (arquivos
  existentes: `capability_adapter.go` + `pricing_reader_test.go` +
  `items_multiget_reader_test.go`; resto = arquivos novos), zero migração, zero UI
- Actual: build/vet limpos; `go test .../connectors/...` — `ok` para melhorenvio (3.21s),
  mercado_livre (12.10s), application (1.93s), domain (1.56s), transport (2.12s);
  `npm run harness:integration` → `status=passed`, `run_id=56b877e5001f457dadd0638014660ba2`,
  72 migrações embarcadas aplicadas. Diff stat confirmado: 4 arquivos existentes tocados
  (exatamente os 3 nomeados + validation-contract.md, doc) + 2 arquivos novos
  (resilience_decorator.go/_test.go). Zero migração, zero UI. `capability_adapter.go`
  confirmado intocado no commit próprio do F-02 (`e6dcf75e`) antes da correção do F-01.
- Artifact: `_chip-m01/EVIDENCE.md` §M01-C6 (comandos completos + saída).
Blocking failure: lane vermelha, ou escrita fora do write-set (Ownership M-01)
Blocking failure observed: No
Owner: QA Validator

**Amendment (2026-08-01, orchestrator M-01, self-classified contract conflict):**
Required Outcome names `doRawWithHeaders` as THE single choke point that must be resilient.
F-01's first pass wired retry only into `ReadFeeQuote` via a bypass (`doRawRetrying`) to avoid
touching two pre-existing tests that locked single-attempt-on-429 through the shared choke
point — leaving F-02's multiget/prices readers (the mission's stated reason this milestone
exists — see mission.md "why this milestone") completely unprotected. That is a direct
contradiction between C6's literal file-touch restriction and the Required Outcome's own
"choke point único" clause. Resolution: `doRawWithHeaders` now branches by HTTP method (GET
retries via the decorator, non-GET stays single-attempt) — the actual universal fix — which
required correcting `pricing_reader_test.go`'s `TestPriceToWinErrorMappingWithoutRetry`
"rate limited" case (extracted into `TestPriceToWinRateLimitedRetriesThenExhausts`, asserting
the NEW correct behavior: retry-then-exhaust, not single-attempt) and
`items_multiget_reader_test.go`'s `TestGetItemsMultigetBatch429Propagates` (same pattern).
Both are inside M-01's exclusively-owned directory; operator's dispatch-prompt authorization
("autorizo deletar e refatorar o que for possível para atingir a solução correta e global
maximum... não somente feature a feature") covers this. `price_writer_test.go` and
`listing_writer_test.go` (PUT call sites) needed zero changes and remain green unmodified —
verified by direct rerun, not assumed.

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

- Current status: **hermetic criteria M01-C1..C6 all Passed** (2026-08-01), independently
  adversarially reviewed per feature (1 sonnet reviewer each, lean-review operator override) —
  see `_chip-m01/EVIDENCE.md` for full command/output trail. Changes are uncommitted working
  tree on top of `e6dcf75e`; not yet committed, not yet merged.
- Next owner: hub.
- Next action: (1) commit the uncommitted working-tree changes as the milestone's close
  commit; (2) `REQUEST dev-stack-capture` to the hub for the pre-merge baseline (B-6,
  Evidence Requirements) — this session cannot boot the dev stack (hub-owned seam); (3) after
  merge + dev-stack re-point, hub-driven browser QA of M01-U1/M01-U2 per the table below.
- Required files/evidence: este arquivo; `_chip-m01/EVIDENCE.md` (written); pre-merge baseline
  capture (NOT yet written — irrecuperável depois, see Outstanding note in EVIDENCE.md).
- Blockers or open decisions: pre-merge baseline capture requires dev-stack access this
  session does not have — hub action needed before merge.

## Critérios de user-drive (mandato do operador — obrigatório)

M-01 é infra invisível por design ("zero mudança de comportamento para caminhos sem erro") —
o user-drive prova exatamente ISSO.

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M01-U1 | Zero regressão visível com o decorator no fio: /anuncios, /pedidos, /precos, /integracoes carregam e operam sem erro novo no dev stack re-apontado pós-merge | browser drive 4 telas + console limpo (zero erro novo) |
| M01-U2 | Refresh/import existente que atravessa o adapter ML continua funcional (caminho sem erro inalterado) | browser drive de 1 operação real que toca o adapter + comparação das DUAS capturas salvas (baseline pré-merge de Evidence Requirements vs captura pós-merge) — B-6 |
