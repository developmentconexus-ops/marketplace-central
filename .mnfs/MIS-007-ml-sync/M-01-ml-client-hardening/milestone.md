# M-01-ml-client-hardening

```yaml
id: M-01
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

MIS-007 ml-sync — [mission.md](../mission.md); spine ADR-01..14; contratos IC-01..07.

## Outcome

O client ML sobrevive a 429 e a concorrência: backoff exponencial + jitter + `Retry-After`
honrado + token-bucket POR INSTALLATION compartilhado entre goroutines, tudo UMA vez no
choke point `doRawWithHeaders`; multiget `/items?ids=` (20/batch) disponível; regra DTO
`Raw json.RawMessage` estabelecida. Zero schema, zero UI, zero mudança de comportamento
para caminhos sem erro.

## Why This Milestone Exists

Design §6: pré-requisito das duas ondas. Todo backfill multiplica volume contra um adapter
onde 429 é erro seco no primeiro contato (`capability_adapter.go:654-655`); backfill escrito
antes = escrito contra semântica de falha que muda embaixo dele. Depois deste milestone,
`capability_adapter.go` CONGELA (ADR-02 do spine/A-02): endpoint novo = arquivo novo.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | resilience-decorator | [F-01-resilience-decorator/feature.md](F-01-resilience-decorator/feature.md) |
| F-02 | items-multiget-raw-dto | [F-02-items-multiget-raw-dto/feature.md](F-02-items-multiget-raw-dto/feature.md) |

## Dependencies

Nenhuma (raiz do DAG, lane A — ∥ M-02 ∥ M-09).

## Ownership & Concurrency

- Exclusive surfaces: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/`
  (dir `mercado_livre`, package `mercadolivre`; NÃO confundir com
  `internal/modules/integrations/adapters/mercadolivre/` — OAuth/metadata de catálogo,
  NINGUÉM edita nesta missão — P5 r02 N-2). Exclusividade = ARQUIVOS EXISTENTES do package
  (ÚNICO a editar
  `capability_adapter.go` nesta missão); arquivos NOVOS no package são permitidos a
  milestones downstream (M-03 F-01 readers de ingest; M-01 F-02 pode entregar reader de
  prices — ver F-02).
- Migration block: none.
- Predicted seam locks: none (nenhum port novo; decorator é interno ao adapter).
- Runs in parallel with: M-02, M-09.
- Internal feature DAG: F-01 ∥ F-02 (arquivos disjuntos: F-01 = capability_adapter.go +
  arquivo novo do decorator; F-02 = arquivo novo de multiget + DTOs).

## Risks

- Limite ~1500 req/min é `assumed` (fato #11, `research/external-ml-api-facts.md`) →
  limiter CONFIGURÁVEL, nunca constante
  compilada; default conservador.
- Fundir com o backoff de refresh OAuth (`refresh_policy.go:18-27`) — mecanismo SEPARADO,
  não tocar.

## Done Means

- 429 com `Retry-After: 2` → retry após ≥2s, teste NOMEIA o tempo decorrido (asserção
  "eventually succeeds" = passe vácuo, reprova).
- N goroutines concorrentes → timestamps observados das requests respeitam o bucket
  (asserção sobre timestamps, não sobre config).
- Budget esgotado → `ErrCodeProviderRateLimited` nomeando tentativas + último Retry-After.
- Multiget: 45 ids → 3 chamadas (20+20+5), DTOs com `Raw` populado.
- Write (`PUT /items`) marcado no-retry (opt-out provado por teste).
- Lanes verdes (unit + integration hermética `GOCACHE=.gocache`).

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator (hub).
- Next action: despachar F-01 ∥ F-02.
- Required files/evidence: `validation-contract.md` (P6), `M-01/validation-result.md`.
- Blockers or open decisions: none.

## Correction Handoff

N/A (planning inicial).
