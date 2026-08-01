# F-01-sync-health-endpoint

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-09
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-09 sync-observability.

## Brief

`GET /sync/health` no módulo sync (transport NOVO — `sync/transport/health_handler.go`;
não confundir com `/sync/runs` do módulo integrations, `run_read_handler.go:34`, que fica
INTOCADO): lê sync_state (repo existente `sync_state_repo.go:35` Read — ou query de lista
nova no repo, read-only) e monta payload IC-05 verbatim:
`{entities:[{entity,last_success_at,last_incremental_at,consecutive_failures,phase,
last_error}], webhook:{last_notification_at,pending,dropped_24h}}`. Bloco webhook via porta
`WebhookStatsReader` (assinatura IC-05); implementação DEFAULT retorna o estado canônico
inicial `{"last_notification_at":null,"pending":0,"dropped_24h":0}` (IC-05 — timestamp
null, contadores ZERO: contagem vazia é fato conhecido, não desconhecido). O health service
expõe `WithWebhookStatsReader(...)` (setter/builder) p/ o M-08 trocar a impl a partir da
PRÓPRIA região ancorada dele no root.go — M-08 NUNCA edita a linha de construção do M-09;
porta opcional ganha compile-time assert (lição catalog-503). SEM param `installation_id` —
tenant vem do ctx (IC-05 pina); a leitura varre TODAS as rows de `sync_state` do tenant,
independente de `installation_id` — INCLUI o sentinela de escopo ERP
`installation_id = "erp"` (`sync/composition/scheduler.go:11`), onde vive a row de
products; scan restrito a instalações ML devolveria ZERO products e tornaria o Done Means
inatingível (auditoria P5 r05 F-r05-1). Classe
interativa, fora de registerBatchRoutes. Campos de entities sem observação = null JSON
(nunca 0, nunca string vazia). OpenAPI + SDK `getSyncHealth` (idioma dos métodos em
`index.ts:2144`) MESMO commit. Registro: `httpx.InteractiveRouteClass` + 1 linha ancorada
root.go.

EARS:
- While sync_state tem row de products, when GET roda, the entity products shall carregar
  timestamps reais e as ausentes shall vir com nulls.
- While inbox não existe (pré-M-08), when GET roda, the bloco webhook shall ser o estado
  canônico inicial `{"last_notification_at":null,"pending":0,"dropped_24h":0}` via impl
  default — endpoint nunca 500 por ausência de dado.
- While entity com falhas consecutivas, when GET roda, the consecutive_failures + last_error
  shall refletir sync_state.

## Inputs

IC-05 (payload + porta binding); fatos de `research/p5-prerequisites.md`: #10 (sync_state
repo assinaturas verbatim), #12 (/sync/runs idioma de envelope/params/registro —
REFERÊNCIA, não reuso), #7 (root.go anchors); 0075 schema VERIFICADO (`0075_sync_sync_state.sql`: cursor jsonb,
last_full_sync_at, last_incremental_at, last_error jsonb, consecutive_failures int NOT
NULL DEFAULT 0 — todo campo de entities[] tem coluna real ou derivação pinada; `phase`
deriva do cursor jsonb; `last_success_at` = `GREATEST(last_full_sync_at,
last_incremental_at)`, NULL só quando ambos NULL — NUNCA só o full, que congela em
entidade incremental-only (IC-05; auditoria P5 r04 F-r04-1)).

## Expected Output

Handler + read service fino + porta WebhookStatsReader + impl default (estado canônico
inicial) + setter `WithWebhookStatsReader` com compile-time assert + par OpenAPI+SDK +
linha root.go.

## Constraints

- sync_state READ-ONLY; scheduler.go INTOCADO.
- Sem cache — leitura direta (payload minúsculo, interativo).
- `phase` vem do cursor persistido (JSON IC-06 tem `phase`) — parse tolerante: cursor de
  formato desconhecido → phase null (nunca erro).

## Inputs/Outputs

Payload canonical IC-05 (binding — spec não re-decide shape). Erros: envelope apierror
vigente. Input: NENHUM param — tenant do ctx (IC-05).

## Negative Scenarios

- Tenant sem NENHUMA row em `sync_state` → `entities: []` (array VAZIO — IC-05 pina "uma
  row por entidade REGISTRADA", zero rows = zero entries, M09-C1; reescrito P7 r02 A-13) +
  webhook canônico inicial (não 404 — o caso vazio é ausência de ROW, não ausência de
  instalação ML — F-r05-1).
- Cursor jsonb corrompido → phase null + last_error preservado.

## Ownership

- Owned paths: `sync/transport/` (novo), `sync/application/health_*` (arquivos NOVOS com
  prefixo `health_` — grant verbatim da matriz, mission.md; read service + porta; P7 r02
  A-6), par OpenAPI+SDK, linha ancorada root.go.
- Forbidden paths: integrations transport (/sync/runs); scheduler.go; sync_state_repo
  writes; FE (F-02).
- Parallel-safe with: none — primeira do M-09.

## Validation Expectations

- Fixture: 1 entity com sucesso + 1 com falha + 1 ausente → JSON golden com nulls exatos.
- Fixture negativa F-r04-1: entidade com `last_full_sync_at` velho + `last_incremental_at`
  recente → `last_success_at` = o incremental no JSON; implementação que iguala ao full
  REPROVA.
- Teste da porta: impl default → bloco webhook byte-igual ao canônico inicial IC-05; fake
  injetado via WithWebhookStatsReader → valores do fake na resposta DA ROTA REGISTRADA
  (asserção via handler montado, não só no service isolado — injeção por
  referência/ponteiro, IC-05 §seam; P5 r04 F-r04-2. Seam provado p/ M-08 trocar sem tocar
  construção).
- Registro interativo provado (controle negativo FORA de registerBatchRoutes — lição
  deadline-class).

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md`.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none.
