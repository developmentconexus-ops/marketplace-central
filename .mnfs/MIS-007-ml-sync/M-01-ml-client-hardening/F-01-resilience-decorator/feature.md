# F-01-resilience-decorator

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-01
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-01 ml-client-hardening.

## Brief

Decorator de resiliência no choke point único do adapter ML (`doRawWithHeaders`,
`capability_adapter.go:712-731`): backoff exponencial + jitter, `Retry-After` honrado no
429, token-bucket por `ProviderAccountRef` (installation) compartilhado entre goroutines.
429 deixa de ser erro seco (`:654-655`, `:462-465`, `:578-583`): vira
`ErrCodeProviderRateLimited` SÓ após budget esgotado, com tentativas + último Retry-After
no erro. Flag de opt-out no-retry aplicada ao write existente (`PUT /items`, `:444`).

EARS:
- While budget disponível, when resposta 429 com Retry-After, the decorator shall aguardar
  ≥Retry-After (com jitter) e re-tentar, sem propagar erro.
- While budget esgotado, when novo 429, the decorator shall retornar
  ErrCodeProviderRateLimited nomeando tentativas e último Retry-After.
- While N goroutines da mesma installation em voo, when limite do bucket atingido, the
  decorator shall enfileirar (bloquear) até haver token — nunca estourar o limite somado.
- While a chamada é write marcado no-retry, when qualquer erro, the decorator shall
  propagar imediatamente sem retry.

## Inputs

- ADR-02 (spine mission.md) — semântica binding.
- `research/codebase-ingest-side.md` — mapa do adapter (choke point, sítios de 429,
  providerDiag clipping :681-688, AccessTokenResolver :26 / root.go:370-378).
- Limite configurável: default conservador (ex.: 900 req/min por installation), fonte
  fato #11 (`research/external-ml-api-facts.md`) `assumed`.

## Expected Output

- Arquivo novo no package `mercadolivre` com o decorator + config; `capability_adapter.go`
  editado SÓ para: envolver `doRawWithHeaders`, remover os returns secos de 429, marcar o
  write no-retry.
- Erro exaurido carrega `attempts` e `last_retry_after` (campos inspecionáveis).

## Constraints

- `AccessTokenResolver` e backoff de refresh OAuth (`refresh_policy.go:18-27`) = mecanismo
  SEPARADO — proibido fundir.
- Mapeamento de erros não-429 e clipping de `providerDiag` (512 runes) inalterados.
- Limiter configurável via config existente do adapter — NUNCA constante compilada.
- Sem persistência, sem schema, sem mudança de assinatura pública do capability set
  (`ProviderCapabilitySet` :79-92 permanece único).

## Negative Scenarios

- 429 sem header Retry-After → backoff exponencial puro com jitter (base documentada no
  código de teste), não crash, não retry imediato.
- Retry-After não-numérico → tratado como ausente (backoff puro); nunca parse panic.
- Contexto cancelado durante espera → retorna ctx.Err() imediatamente (não segura o token).

## Ownership

- Owned paths: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/`
  (decorator novo + `capability_adapter.go`).
- Forbidden paths: qualquer coisa fora do package `mercadolivre`; `refresh_policy.go`.
- Parallel-safe with: F-02 (eixo files: F-02 só cria arquivos novos de multiget/DTO).

## Validation Expectations

- Teste com transport fake 429+`Retry-After: 2`: asserção nomeia elapsed ≥2s.
- Teste de bucket: N chamadas concorrentes, asserção sobre TIMESTAMPS observados das
  requests (não sobre valor de config).
- Must-fail: remover o wait → teste falha nomeando "elapsed < retry-after".
- Teste no-retry: write com erro → exatamente 1 tentativa.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` = execução, não planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (via hub dispatch).
- Next action: criar `spec.md`.
- Required files/evidence: `F-01-resilience-decorator/validation.md` na execução.
- Blockers or open decisions: none.
