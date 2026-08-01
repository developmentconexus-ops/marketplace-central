# Milestone Validation Contract — M-09-sync-observability

```yaml
id: M-09-VC
type: milestone-validation-contract
status: passed
owner: Mission Strategist
parent: MIS-007
created: 2026-08-01
updated: 2026-08-01
validation_level: QA-0
lifecycle_scope: milestone
base_sha: dd89d4b3
```

Verdicts binários. Evidência = caminho inspecionável concreto. Seams: sync_state (0075)
existe e tem dado real de products (MIS-006) — o live deste milestone é barato e imediato;
bloco webhook usa impl DEFAULT da porta até M-08 (estado canônico inicial IC-05 — não é
stub de seam, é o contrato).

## Milestone ID

M-09

## QA Level

QA-0

## Required Outcome

Operador vê saúde do sync sem SQL: `GET /sync/health` NOVO (não reusa /sync/runs), payload
IC-05 §Required Outputs (entities[] de sync_state + bloco webhook), método SDK
`getSyncHealth`, porta `WebhookStatsReader` + default + setter, seção "Saúde do sync" na
IntegracoesPage. NULLs honestos em toda parte.

## Criteria

## Criterion: Payload IC-05 com products real + nulls honestos
ID: M09-C1
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: `GET /sync/health` no dev stack (pré M-04/M-06) + SELECT sync_state
- Expected: products com timestamps REAIS (== SELECT); TODAS as rows de sync_state
  retornadas incluindo sentinel ERP `installation_id="erp"` (F-r05-1 — scan all-rows);
  entidade nunca-rodada = campos null (NUNCA timestamp fabricado);
  `last_success_at = GREATEST(last_full_sync_at, last_incremental_at)` (F-r04-1);
  `phase` do cursor jsonb, null quando ausente; tenant sem NENHUMA row → `entities: []`
  (lista vazia honesta)
- Actual:
- Artifact:
Blocking failure: row de sync_state filtrada, null virando 0/"", ou last_success_at
ignorando incremental
Blocking failure observed: No
Owner: QA Validator

## Criterion: Fixture negativa incremental-only (F-r04-1)
ID: M09-C2
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: fixture: entidade com full VELHO + incremental RECENTE
- Expected: `last_success_at` == o incremental (fresco); FE renderiza badge VERDE +
  timestamp relativo do incremental — NUNCA "há N dias"/cinza. ("Verde" é decidido por
  ESTADO, não corte temporal: `last_success_at` não-nulo ∧ `consecutive_failures == 0`
  — P7 r02 A-8)
- Actual:
- Artifact:
Blocking failure: entidade saudável 5min mostrada como velha (o defeito que F-r04-1
nomeia)
Blocking failure observed: No
Owner: QA Validator

## Criterion: Seam WebhookStatsReader provado na ROTA
ID: M09-C3
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: teste na rota REGISTRADA (não no handler isolado): default + fake injetado
- Expected: default → bloco webhook BYTE-IGUAL ao canônico IC-05
  `{"last_notification_at":null,"pending":0,"dropped_24h":0}`; fake injetado via
  `WithWebhookStatsReader` (injeção por referência/ponteiro — F-r04-2) → valores do fake
  observáveis no `GET /sync/health` montado; compile-time assert da porta presente
  (lição catalog-503)
- Actual:
- Artifact:
Blocking failure: injeção por cópia (fake invisível na rota), ou default divergente do
canônico
Blocking failure observed: No
Owner: QA Validator

## Criterion: FE — 4 estados renderizados honestos
ID: M09-C4
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: fixtures por estado no teste de componente + browser QA
- Expected: verde (`last_success_at` não-nulo ∧ `consecutive_failures == 0` — discriminador
  por estado, P7 r02 A-8) / vermelho (falhas consecutivas + last_error tooltip) /
  cinza "nunca" (nulls — NUNCA "0 min atrás") / webhook inicial "nenhuma notificação
  recebida" (o FATO; rótulo de veredito "não configurado" PROIBIDO — F-r07-3); erro de
  fetch → card com erro nomeado (apierror/hasCode), resto da página INTACTO
- Actual:
- Artifact:
Blocking failure: qualquer estado fabricando atividade/configuração que o payload não
carrega
Blocking failure observed: No
Drive (UI — agent-browser; UI criteria only):
- Fixture: dev stack real (products vivo; listings/orders ainda nulls; webhook default)
- Steps:
  - open /integracoes
  - assert text "Saúde do sync"
  - assert text "nenhuma notificação recebida"
  - assert text "nunca"
- Expected: seção nova após ProviderConnectCard; products com timestamp relativo real;
  entidades não-rodadas "nunca" em cinza; webhook com o rótulo de FATO
Owner: QA Validator

## Criterion: Par SDK/OpenAPI + página intacta
ID: M09-C5
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `git show --stat` do commit de contrato; tsc; `IntegracoesPage.test.tsx`
- Expected: `/sync/health` no OpenAPI + `getSyncHealth` no sdk-runtime MESMO commit; tsc
  verde; teste de página existente VERDE sem edição; diff da página = 1 linha de mount
  (card é arquivo novo)
- Actual:
- Artifact:
Blocking failure: par quebrado, ou IntegracoesPage editada além do mount
Blocking failure observed: No
Owner: QA Validator

## Criterion: Re-drive pós-lanes — entities acendem SEM mudança
ID: M09-C6
Level: Milestone
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: re-drive do hub após M-04/M-06 fecharem (critério DIFERIDO — não bloqueia close
  do M-09; bloqueia close da MISSÃO via MIS07-C8)
- Expected: listings/orders aparecem na seção sem NENHUM commit no código do M-09
  (endpoint entidade-agnóstico provado no mundo real)
- Actual:
- Artifact:
Blocking failure: entidade nova exigindo mudança no M-09 (lista hardcoded escondida)
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Payload JSON capturado (não descrição); SELECT sync_state de conferência salvo.
- Screenshots light+dark da seção em `docs/design/evidence/`.
- Byte-igualdade do default provada por comparação de JSON serializado, com a ressalva
  round-trip JSONB (nunca string-eq byte-exact em round-trip — lição M-01 MIS-006; comparar
  estrutura).

## Blocking Failures

- Tela afirmando o que o payload não carrega (veredito de configuração, zero fabricado,
  "0 min atrás") = blocking (M09-C1/C2/C4).
- Fake invisível na rota (injeção por cópia) = blocking (M09-C3).
- Regressão na IntegracoesPage = blocking (M09-C5).

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: n/a

## Handoff

- Current status: planned.
- Next owner: hub (lane A, ∥ M-01 ∥ M-02).
- Next action: F-01 → F-02.
- Required files/evidence: este arquivo; `M-09/validation-result.md`.
- Blockers or open decisions: none.

## Critérios de user-drive (mandato do operador — obrigatório)

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M09-U1 | Seção "Saúde do sync" dirigida com dado REAL: products com timestamp relativo que confere com SELECT de sync_state (hover mostra absoluto) | browser drive + SELECT lado a lado |
| M09-U2 | Estados honestos visíveis no mesmo drive: entidade nunca-rodada = "nunca" cinza; webhook = "nenhuma notificação recebida" (fato, não veredito) | browser drive + screenshot dos 2 estados |
| M09-U3 | Health 500 simulado → card mostra erro nomeado, resto de /integracoes segue operável (cards existentes intactos) | browser drive com fetch quebrado |
