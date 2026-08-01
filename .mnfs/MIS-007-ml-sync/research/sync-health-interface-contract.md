# Interface Contract — saúde de sync (/integracoes)

```yaml
id: IC-05
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

Endpoint novo `GET /sync/health` (server_core) + seção nova na IntegracoesPage (FE) +
consumo de `sync_state` (0075) e `notifications_inbox` (IC-04). Dono: M-09. As entidades
`listings`/`orders` acendem quando M-04/M-06 registram jobs — o endpoint nasce correto com
qualquer subconjunto.

## Why This Contract Exists

Q4 é critério de missão dirigido em browser. `GET /sync/runs` existe SEM consumidor FE —
decisão de planning: NÃO reusar (shape é lista de runs, não saúde agregada por entidade);
`/sync/health` é agregado novo. Sem pino, M-09 inventa shape que não bate com o que
M-04/M-06 gravam em `sync_state`.

## Resources Or Entities

- Rota `GET /sync/health` (server_core, classe interativa).
- Fontes: `sync_state` (por entidade), `notifications_inbox` (agregados).
- SDK: método novo tipado; FE: seção "Saúde do sync" na IntegracoesPage existente
  (nenhuma tela nova — design §8).

## Seam M-09 ↔ M-08 — porta `WebhookStatsReader` (binding)

- Porta definida no application do módulo sync (M-09, dono):
  `WebhookStatsReader interface { WebhookStats(ctx, tenantID) (WebhookStats, error) }`
  com `WebhookStats {LastNotificationAt *time.Time; Pending int; Dropped24h int}` (nomes
  finais Go = spec; SEMÂNTICA e forma binding).
- Impl DEFAULT (M-09): retorna o estado canônico inicial
  `{last_notification_at: null, pending: 0, dropped_24h: 0}` — nunca erro.
- Troca: health service expõe `WithWebhookStatsReader(...)` (setter/builder). M-08 chama o
  setter DA PRÓPRIA região ancorada no root.go com a impl real (lê o inbox). Código e linha
  de construção do M-09 NUNCA editados pelo M-08.
- SEMÂNTICA DE MUTAÇÃO pinada (auditoria P5 r04 F-r04-2): a injeção é por
  REFERÊNCIA/ponteiro — o reader injetado tem que ser observável através do handler JÁ
  REGISTRADO pelo M-09. Builder por VALOR que retorna cópia (idioma `root.go:853-855`
  `pricingHandler = pricingHandler.WithCalc(calcSvc)`, `Register(mux)` em `:858` —
  linhas medidas, p5-prerequisites §7; precisão corrigida P7 r02 A-5,
  F-r06-5
  `pricingHandler = pricingHandler.WithCalc(...)`) é PROIBIDO aqui: M-08 mutaria uma cópia
  que ninguém serve e o endpoint devolveria o default p/ sempre. Prova obrigatória na
  ROTA, não na porta: após a injeção, `GET /sync/health` (handler registrado) retorna
  stats derivadas do inbox.
- Porta opcional ⇒ compile-time assert obrigatório (lição catalog-503: decorator apagou
  porta opcional em silêncio).

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| GetSyncHealth | render de /integracoes | tenant (ctx) | payload abaixo | ordenação: entities por nome asc (determinística) |

## Fields

### Required Inputs

Nenhum parâmetro além do tenant do contexto.

### Required Outputs

```json
{
  "entities": [
    {
      "entity": "orders",
      "last_success_at": "2026-07-31T12:05:00Z",
      "last_incremental_at": "2026-07-31T12:05:00Z",
      "consecutive_failures": 0,
      "phase": "incremental",
      "last_error": null
    }
  ],
  "webhook": {
    "last_notification_at": "2026-07-31T12:04:58Z",
    "pending": 0,
    "dropped_24h": 0
  }
}
```

- `entities[]`: uma row por entidade REGISTRADA em `sync_state` (products incluído — o
  endpoint é entidade-agnóstico); campos NULL = honesto-desconhecido (nunca timestamp
  fabricado).
- `last_success_at`: sucesso mais recente de QUALQUER tipo —
  `GREATEST(last_full_sync_at, last_incremental_at)`, NULL só quando ambos NULL
  (auditoria P5 r04 F-r04-1: `last_full_sync_at` CONGELA quando a entidade roda
  incremental — `sync_state_repo.go:62,74-79` grava por COALESCE só a coluna do tipo do
  run; igualar `last_success_at` ao full faria a tela dizer "há N dias" p/ entidade
  saudável sincronizando a cada 5min. Fixture negativa obrigatória: entidade com full
  velho + incremental recente ⇒ `last_success_at` = o incremental, badge fresco).
- `phase`: extraído do cursor JSONB (`cursor.phase`, IC-06); NULL quando cursor ausente.
- `webhook`: estado canônico inicial EXATO quando inbox ainda não existe:
  `{"last_notification_at":null,"pending":0,"dropped_24h":0}` (timestamp null = nunca
  observado; contadores ZERO = contagem vazia é fato conhecido). M-09 pode fechar antes do
  M-08 — FE discrimina pelo `last_notification_at === null` e renderiza o FATO observado:
  "nenhuma notificação recebida". Rótulo de VEREDITO de configuração ("webhook não
  configurado") é PROIBIDO: pós-M-08, instalação CONFIGURADA com inbox vazio (registro
  feito e nenhum evento ainda; janela quieta; worker travado) produz o MESMO
  `{null,0,0}` — o payload não distingue configuração de silêncio, e a tela não pode
  afirmar o que o payload não carrega (honest-unknown; auditoria P5 r07 F-r07-3).
- `last_incremental_at` REAL exige o fix `incremental` (ADR-08, M-02 F-03) — pré-condição
  nomeada do critério de MISSÃO/live: a reprovação acende só depois que M-04/M-06
  registram jobs com `phase` e o campo continua NULL. Pré-fix, NULL uniforme = honesto
  (§NULLs acima) e o gate do M-09 PASSA com ele — a dependência M-02 F-03 permanece SOFT
  no M-09 (auditoria P5 r05 F-r05-4).

## Enums And Statuses

`phase`: vocabulário de IC-06 (`backfill`, `incremental`, `sweep`) — verbatim do cursor,
sem re-mapear.

## Error Cases

Envelope apierror vigente; nenhum code novo (500 genérico existente cobre falha de leitura).

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| falha de leitura interna | 500 | code interno vigente | envelope universal CHIP-ERROR-UNIFY |

## Persistence Expectations

Read-only. Zero escrita; zero chamada ML (ADR-05 vale aqui também).

## Canonical Examples

Sucesso acima. Estado inicial honesto (só products registrado, sem webhook):

```json
{"entities":[{"entity":"products","last_success_at":"2026-07-31T11:00:00Z",
 "last_incremental_at":null,"consecutive_failures":0,"phase":null,"last_error":null}],
 "webhook":{"last_notification_at":null,"pending":0,"dropped_24h":0}}
```

## Database Shape

Nenhuma tabela nova. Consultas: `sync_state` scan por tenant; inbox agregados
(`status='received'` count; `dropped` nas últimas 24h; max `received_at`).

Semântica de tenant dos agregados de inbox (pinada P7 r02 A-10): os agregados são
GLOBAIS — `notifications_inbox.tenant_id` é NULL p/ notificações não-mapeadas (IC-04),
logo NÃO há predicado de tenant nessas 3 consultas hoje. A porta `WebhookStats(ctx,
tenantID)` já carrega o parâmetro; sob a assunção registrada de tenant único
(mission.md Accepted assumptions) global == tenant. Missão futura multi-tenant DEVE
adicionar o predicado antes de expor cross-tenant — sem isso os counts vazam.

## Seed Data

Nenhum. Fixture: sync_state multi-entidade + inbox com pending/dropped.

## Timestamp And ID Semantics

timestamptz UTC ISO-8601 no JSON; NULL viaja como null literal, nunca omitido (shape
estável p/ SDK tipado).

## Compatibility Rules

Campos novos por entidade = aditivos no objeto; FE ignora desconhecidos.

## Route Namespace

- `GET /sync/health` sob prefixo `/sync/` existente (onde vive `/sync/runs`); mount pelo
  M-09 via constructor próprio (ADR-14).
- OpenAPI + SDK MESMO commit; o commit de contrato do M-09 entra na serialização do hub
  (ADR-14 emendado P5 r03 P-2: ≤1 COMMIT FE em voo; código paraleliza).
- FE: seção dentro de `IntegracoesPage` — sem rota de página nova.

## Transport And Integration

Auth/session existentes da API; nada cruza origem nova.

## Must Preserve

- `GET /sync/runs` permanece intocado (sem consumidor FE, sem quebra).
- Entidade-agnóstico: registrar job novo NÃO muda o endpoint.

## Must Not Decide In Feature Execution

- Shape do payload; decisão de não-reusar /sync/runs; fonte do `phase`; semântica NULL.

## Validation Impact

- Q4 dirigido em browser: /integracoes mostra saúde por entidade + última notificação.
- Critério de estado honesto: entidade sem run ainda mostra NULLs, não zeros/datas fake
  (M0X-U*).
