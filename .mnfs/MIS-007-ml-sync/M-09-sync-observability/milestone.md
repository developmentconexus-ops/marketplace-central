# M-09-sync-observability

```yaml
id: M-09
type: milestone
status: passed
owner: Mission Strategist
parent: MIS-007
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-007 ml-sync — [mission.md](../mission.md); IC-05 (sync health, binding integral).

## Outcome

Operador vê saúde do sync sem SQL: `GET /sync/health` NOVO (decisão IC-05: NÃO reusa
/sync/runs — runs é telemetria de operação, health é estado consolidado por entidade),
payload IC-05 (entities[] de sync_state + bloco webhook), método SDK, e seção "Saúde do
sync" na IntegracoesPage. NULLs honestos: entidade nunca rodada = null; webhook antes do
M-08 = estado canônico inicial IC-05 `{"last_notification_at":null,"pending":0,
"dropped_24h":0}` via impl default da porta (M-08 injeta a real via setter, sem editar
código do M-09).

## Why This Milestone Exists

Design §8: sync invisível = sync que apodrece em silêncio (MIS-006 provou — scheduler
NO-OP por semanas sem ninguém ver). Lane A: só lê sync_state (0075, existente) — zero
dependência dos milestones de dado. `consecutive_failures`/`phase` leem o que o schema
0075 tem HOJE; campos que dependem do fix incremental (M-02 F-03) ficam honestos até lá
(IC-05 §NULLs).

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | sync-health-endpoint | [F-01-sync-health-endpoint/feature.md](F-01-sync-health-endpoint/feature.md) |
| F-02 | integracoes-health-section | [F-02-integracoes-health-section/feature.md](F-02-integracoes-health-section/feature.md) |

## Dependencies

- Nenhuma dura (sync_state 0075 existe). SOFT: M-02 F-03 (last_incremental_at vira
  significativo), M-04/M-06 (entities acendem), M-08 (bloco webhook real).

## Ownership & Concurrency

- Exclusive surfaces: handler novo no módulo sync (transport novo), porta
  `WebhookStatsReader` + impl default (estado canônico inicial IC-05) + setter
  `WithWebhookStatsReader` com compile-time assert (M-08 injeta a real DA REGIÃO DELE no
  root.go — nunca edita construção do M-09), par OpenAPI+SDK `/sync/health`
  (`getSyncHealth`), 1 linha ancorada root.go, seção nova na IntegracoesPage.
- Migration block: nenhum.
- Predicted seam locks: sync_state READ-ONLY; `scheduler.go` INTOCADO (M-02 é o dono do
  fix); porta WebhookStatsReader = seam publicado p/ M-08 (assinatura em IC-05, binding).
- Runs in parallel with: M-01, M-02 (lane A — superfícies disjuntas; único toque comum =
  root.go linha ancorada própria, ADR-14). Contrato FE: M-09 é o ÚNICO milestone FE da
  lane A — sem colisão de contrato.
- Internal feature DAG: F-01 → F-02.

## Risks

- Tela de saúde que mente é pior que nenhuma: TODO campo sem observação renderiza
  `—`/`nunca` — nunca zero, nunca "ok" default (ADR-17 / AGENTS unknown≠zero).
- IntegracoesPage é superfície compartilhada FE — M-09 adiciona seção NOVA (arquivo novo de
  card, 1 linha de mount em `IntegracoesPage.tsx:558-574`); nenhum outro milestone da
  missão toca essa página.

## Done Means

- `GET /sync/health` responde payload IC-05 com products (entidade viva de MIS-006) real e
  listings/orders nulls-honestos ANTES dos M-04/M-06 (prova do NULL honesto no live).
- Após M-04/M-06 fecharem (re-drive do hub): entities acendem SEM mudança no M-09.
- Bloco webhook: estado canônico inicial IC-05 via impl default; contrato do seam provado
  na ROTA registrada (default byte-igual ao canônico + fake injetado via setter observável
  no `GET /sync/health` montado — injeção por referência, IC-05 §seam; P5 r04 F-r04-2).
- Seção na IntegracoesPage renderiza os 3 estados (verde/falha/nunca) — fixture por estado;
  tsc verde; par OpenAPI+SDK mesmo commit.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator (hub) — lane A.
- Next action: F-01 spec.
- Required files/evidence: `validation-contract.md` (P6), `validation-result.md`.
- Blockers or open decisions: none.

## Correction Handoff

N/A.
