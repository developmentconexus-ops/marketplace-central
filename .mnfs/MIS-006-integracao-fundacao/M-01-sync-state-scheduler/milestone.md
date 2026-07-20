# M-01 — sync_state + scheduler skeleton

```yaml
id: M-01
type: milestone
mission: MIS-006
status: draft
depends_on: []   # Fase 0; coordena additive-lock em composition/root.go com M-02 (dono)
base_sha: 138aac3d
validation_level: QA-0
```

## Objective

Criar a tabela `sync_state` (E8) e um scheduler skeleton — loop + cursor read/write + `last_error`
+ seam de registro, cadence-agnostic — reusando o padrão de ticker já provado em produção
(`internal/composition/root.go` bloco de tickers, `integrationsbg.NewFeeSyncScheduler`
~:462-464 nesta base; mission cita `root.go:577` na base_sha original — número dá drift entre
worktrees, o TEMPLATE é o mesmo). Fase 0: fundação de rastreabilidade, sem disparar nenhum job
ML pesado.

## Scope

- Migração bloco A: `CREATE TABLE sync_state` conforme E8 (`tenant_id, installation_id, entity,
  cursor JSONB, schedule JSONB, last_full_sync_at, last_incremental_at, last_error JSONB NULL,
  consecutive_failures INT`).
- Novo módulo `internal/modules/sync/` (domain + application + adapters/postgres) com:
  - repositório de leitura/escrita de `sync_state` por `(tenant_id, installation_id, entity)`.
  - scheduler skeleton reusando o formato `NewXxxScheduler(...).Start(ctx)` (padrão
    `integrationsbg.NewFeeSyncScheduler`, `research/refactor-inventory-backend.md` §7)
    — loop com intervalo configurável, lê cursor, escreve `last_full_sync_at`/`last_incremental_at`,
    grava `last_error` em falha, incrementa `consecutive_failures`.
  - seam de registro: forma pela qual um `entity` futuro (products/listings/orders/market/tariffs)
    se registra no scheduler SEM editar o loop central (interface/callback, não switch hardcoded).
- Entry em `contracts/governance/modules.json` para o módulo `sync`.
- Append additive-lock no bloco de tickers de `composition/root.go` (só o entity `products`
  registrado nesta milestone — corpo do job é NO-OP/placeholder; M-03/M-04 conectam o Sync real).
- `schedule` gravado como JSONB genérico (ex.: `{"kind":"interval","every_seconds":900}` ou
  `{"kind":"cron","expr":"..."}") — nenhuma cadência real hardcoded como "diário" (D6).

## Non-Scope

- Qualquer sync ML pesado (anúncios/pedidos/mercado/tarifas) — missões seguintes.
- `products_mirror`, `ProductSourceAdapter`, `active_source` config — M-02.
- Corpo real do job de sync de produtos (o que o Sync efetivamente faz) — M-03 (xlsx) / M-04
  (Sankhya) conectam a lógica; M-01 só entrega o esqueleto que os recebe.
- UI de observabilidade (`/importacoes`) — M-06.
- Qualquer alteração ao ticker `NewFeeSyncScheduler`/`NewRefreshTicker` existentes (KEEP).

## Feature Briefs

### F-01 — Tabela `sync_state` (migração bloco A)

**Inputs/Outputs**
- Input: nenhum (DDL puro).
- Output: tabela `sync_state` com colunas conforme E8; PK/unique em
  `(tenant_id, installation_id, entity)`; `cursor`/`schedule`/`last_error` NULL-able JSONB;
  `consecutive_failures` default 0.

**EARS**
- While a migração bloco A ainda não rodou, when o serviço sobe, the sistema shall falhar o boot
  se `sync_state` for referenciada por código que já assume sua existência (nenhum código deve
  assumir isso antes da migração aplicar — ordem de deploy).
- While o banco está no estado pós-migração, when uma row é inserida com `entity` fora do enum
  documentado (`products|listings|orders|market|tariffs`), the sistema shall aceitar a escrita
  (enum é semântico na aplicação, não constraint de banco — ver Non-negotiables) mas a camada de
  aplicação shall validar antes de persistir.
- While nenhuma row existe para `(tenant_id, installation_id, entity)`, when o repositório de sync
  lê o cursor, the sistema shall retornar "nunca sincronizado" honesto (não erro, não cursor
  fabricado) e permitir first-run.
- While uma row já existe, when um novo ciclo grava `last_full_sync_at`, the sistema shall fazer
  UPSERT (nunca duplicar linha por `(tenant_id, installation_id, entity)`).
- While dois produtores concorrentes (ex. M-03 hook de import E M-07 enqueue) acumulam
  `codigo_produto` no MESMO `cursor` de `(tenant, installation_id, entity=market)`, when ambos
  gravam, the repositório `sync_state` shall expor uma operação de append-ao-cursor ATÔMICA
  (single-statement UPDATE ... jsonb, não read-modify-write na aplicação) — evita race de
  lost-update. A operação atômica é propriedade de M-01 (dono de `sync_state`); consumidores
  (M-03/M-07) chamam, não reimplementam.

**Negative Scenarios**
- Insert sem `tenant_id` — rejeitado (NOT NULL; nenhuma query tenant-scoped sem tenant_id).
- Campo `cursor`/`schedule` ausente na primeira escrita — grava NULL, nunca `{}` fabricado como
  "sem cursor" disfarçado de dado real (ADR-17 aplicado a JSONB: NULL é o estado honesto).

**Write-set**
- `apps/server_core/migrations/<próximo bloco A>_sync_sync_state.sql` (número alocado pelo hub
  no merge — só a ordem bloco A < bloco B de M-02 importa, ver mission §Clarified Decisions).
- `internal/modules/sync/domain/sync_state.go` (struct + enum de `entity`).
- `internal/modules/sync/adapters/postgres/sync_state_repo.go` (upsert + read por chave composta).

### F-02 — Scheduler skeleton (loop + cursor + seam de registro)

**Inputs/Outputs**
- Input: lista de "jobs registrados" (entity → função de sync a executar) + intervalo/config lido
  de `sync_state.schedule` (ou default se ausente na primeira execução).
- Output: side-effect — chamada periódica ao job registrado; `sync_state` atualizado
  (`last_full_sync_at`/`last_incremental_at` em sucesso, `last_error`+`consecutive_failures++`
  em falha); NENHUM job de sync ML pesado registrado nesta milestone (só o seam existe, vazio
  ou com um placeholder no-op para `entity=products`).

**EARS**
- While o processo sobe, when `composition/root.go` inicializa, the scheduler shall iniciar como
  goroutine adicional no MESMO bloco dos tickers existentes (`go ....Start(ctx)`), sem remover ou
  alterar os tickers já presentes (`NewRefreshTicker`, `NewStateCleanup`, `NewFeeSyncScheduler`).
- While um job registrado executa com sucesso, when o ciclo termina, the scheduler shall persistir
  cursor + `last_full_sync_at`/`last_incremental_at` e zerar `consecutive_failures`.
- While um job registrado retorna erro, when o ciclo termina, the scheduler shall persistir
  `last_error` (JSONB com mensagem+timestamp) e incrementar `consecutive_failures`, sem derrubar
  o processo (falha isolada por entity, não propaga ao loop de outro entity).
- While nenhum job pesado está registrado (estado desta milestone), when o loop roda, the
  scheduler shall completar ciclos sem nenhuma chamada a API externa (ML/Oracle) — só leitura/
  escrita local de `sync_state` (ou no-op puro).
- While um novo `entity` precisa ser adicionado no futuro (M-03/M-04/mercado), when o registro
  acontece, the seam shall permitir isso via função/interface de registro (ex.: `RegisterJob(entity,
  fn)`), sem editar o corpo do loop central.

**Negative Scenarios**
- Job registrado trava (nunca retorna) — fora de escopo desta milestone tratar timeout/cancelamento
  do job em si (skeleton só garante que o LOOP não trava se cada ciclo respeitar ctx; timeout de
  job real é responsabilidade de quem registra em M-03+).
- Dois registros para o mesmo `entity` — comportamento indefinido documentado como tal (skeleton
  não resolve conflito; primeira versão assume 1 registro por entity, decisão revisitável).

**Write-set**
- `internal/modules/sync/application/scheduler.go` (loop + `RegisterJob` seam).
- `internal/modules/sync/application/scheduler_test.go`.
- `internal/composition/root.go` (append additive-lock no bloco de tickers, ~linha do bloco atual
  de tickers — ver Ownership abaixo).
- `contracts/governance/modules.json` (entry do módulo `sync`, `composition_required: true`,
  `dependencies: []` — skeleton não depende de nenhum outro módulo).

## Ownership & Concurrency

Eixo canônico (6 eixos, ver `architecture-map.md` "Superfícies compartilhadas"):

| Eixo | M-01 |
|---|---|
| Migração | bloco A (`sync_state`) — precede bloco B de M-02 só por convenção de número, tabelas independentes |
| DB shape (tabelas) | `sync_state` (novo, exclusivo) |
| Módulo Go | `internal/modules/sync/*` (novo, exclusivo) |
| root.go (seam) | **additive-lock**: M-02 é dono do arquivo (source-wiring); M-01 tem permissão pré-autorizada de **append** no bloco de tickers (~:575-577 na base_sha da missão; ~:462-464 nesta worktree — número dá drift, o bloco-alvo é "onde os `go ...Ticker/...Scheduler(...).Start(ctx)` já vivem", logo após os tickers de integrations). Grant liberado no início da onda 1, **release no close de M-01** (hub reconcilia qualquer conflito de merge com M-02; fallback serial se risco alto — ver mission Risks). |
| Contrato/SDK | nenhum endpoint novo nesta milestone (sync_state não é exposto via API ainda — consumo interno) |
| FE surface | nenhum |

Regra de convivência com M-02 (paralelo na mesma onda): M-01 e M-02 tocam `root.go` mas em
seções disjuntas dentro do arquivo (M-01 só o bloco de tickers; M-02 o wiring de source). Se
o merge de ambos colidir de forma não-trivial, hub decide serializar (M-02 primeiro, por ser
"fundação — tudo depende daqui" per mission.md).

## Dependencies

Nenhuma — Fase 0, primeira milestone a executar. Corre em paralelo com M-02 (onda 1), coordenado
só pelo additive-lock de `root.go` acima. M-03/M-04 dependem de M-01 (usam o seam de registro do
scheduler para conectar o Sync real dos adapters); M-06 depende de M-01 (chain-viz lê `sync_state`
para a contagem N-enqueued).

## Validation

Prova MC-07 (mission `validation-contract.md`): "`sync_state` rastreia entidades cadence-agnostic;
scheduler skeleton roda loop sem job ML pesado". Prova mínima inspecionável:
- Migração aplicada; `sync_state` populada com pelo menos 1 row (`entity=products`) após o
  scheduler completar um ciclo (mesmo que o corpo do job seja no-op nesta milestone).
- `cursor`/`schedule` são JSONB genéricos — nenhuma string "daily"/"diário" hardcoded no código
  do scheduler (D6, cadence-agnostic).
- `composition/root.go` grep mostra o novo `go ....Start(ctx)` no bloco de tickers, tickers
  existentes intactos.
- Grep/log da execução do ciclo: zero chamadas a hosts ML (`api.mercadolibre.com`) ou Oracle —
  nenhum job pesado disparado (AC alinhado a Non-Scope).
- Ver `M-01-sync-state-scheduler/validation-contract.md` para critérios binários M01-C1..C5.
