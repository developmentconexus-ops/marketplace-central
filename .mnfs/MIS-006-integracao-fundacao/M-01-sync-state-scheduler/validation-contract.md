# M-01 Validation Contract — sync_state + scheduler skeleton

```yaml
id: M-01-VC
type: milestone-validation-contract
mission: MIS-006
milestone: M-01
validation_level: QA-0
```

Verdicts binários. Evidência = caminho inspecionável concreto (mission profile §5, ladder L0-L4).
Tipos: `ran` (executado, output salvo), `assumed` (design, não executado), `could-not-run`
(bloqueado — nomear bloqueio).

| ID | Critério | Prova mínima inspecionável | Ladder | Evidência |
|----|----------|----------------------------|--------|-----------|
| M01-C1 | `sync_state` existe com shape E8 completo | `\d sync_state` (ou migração lida) mostra `tenant_id, installation_id, entity, cursor JSONB, schedule JSONB, last_full_sync_at, last_incremental_at, last_error JSONB NULL, consecutive_failures INT`; PK/unique em `(tenant_id, installation_id, entity)` | L1 | ran |
| M01-C2 | Scheduler loop registrado em `root.go` sem remover tickers existentes | grep `root.go` pós-diff: `NewRefreshTicker`, `NewStateCleanup`, `NewFeeSyncScheduler` ainda presentes E novo `go ....Start(ctx)` do módulo `sync` presente no mesmo bloco | L1 | ran |
| M01-C3 | Cursor read/write funciona ponta-a-ponta | Rodar 1 ciclo do scheduler (dev stack ou teste de integração local) → `SELECT * FROM sync_state WHERE entity='products'` retorna row com `last_full_sync_at` não-nulo após sucesso | L2 | ran |
| M01-C4 | Falha de job grava `last_error` + incrementa `consecutive_failures`, não derruba o processo | Teste força job a retornar erro → row atualizada com `last_error` JSONB preenchido + `consecutive_failures=1`; processo continua rodando (próximo ciclo de outro entity não afetado) | L1 | ran |
| M01-C5 | NENHUM job ML pesado disparado nesta milestone | Grep do diff: nenhuma chamada a `api.mercadolibre.com`/Oracle/`internal_read`/`connectors/adapters/mercado_livre` a partir do módulo `sync`; corpo do job registrado (se algum) é no-op/placeholder | L0 | ran (grep diff) |
| M01-C6 | Cadence-agnostic: `schedule` é JSONB genérico, sem cadência hardcoded | Grep do código do scheduler: nenhuma string literal `"daily"`/`"diário"`/cron fixo embutido; intervalo/cron vem de `sync_state.schedule` ou de config passada ao construtor, nunca de constante de negócio | L0 | ran (grep diff) |
| M01-C7 | Padrão de ticker reusado, não reinventado | `internal/modules/sync/application/scheduler.go` segue a MESMA forma de `integrationsbg.NewFeeSyncScheduler` (construtor + `.Start(ctx)` + goroutine em `root.go`) — comparação estrutural lado a lado no PR/evidência | L1 | ran |
| M01-C8 | Migração é aditiva (bloco A, sem colidir com bloco B de M-02) | Diff da migração = só `CREATE TABLE` novo; nenhuma alteração a tabela existente; número de arquivo < número alocado ao bloco B (M-02) no merge | L0 | ran (diff) |
| M01-C9 | `modules.json` tem entry do módulo `sync` | `contracts/governance/modules.json` contém `{"id":"sync", "root":"apps/server_core/internal/modules/sync", ...}` | L0 | ran |
| M01-C10 | Toda query a `sync_state` é tenant-scoped | Grep do repositório postgres: toda query SELECT/UPSERT em `sync_state` inclui `tenant_id` na cláusula WHERE/chave (AC-01 da mission) | L0 | ran (grep) |
| M01-C11 | Repo expõe append-ao-cursor atômico (single-statement), não read-modify-write | diff do repo `sync_state`: a operação de acúmulo no `cursor` é um único `UPDATE ... SET cursor = ...jsonb...` (não SELECT-then-UPDATE na aplicação); teste concorrente: 2 appends simultâneos ao mesmo `(tenant,installation_id,entity)` → ambos os `codigo_produto` presentes no cursor (nenhum lost-update) | L1 | ran |

## Anti-critérios (falha se presente)

- AC-M01-1: qualquer query a `sync_state` sem `tenant_id` (herda AC-01 da mission).
- AC-M01-2: cursor/schedule ausente virando `0`/`{}` fabricado em vez de NULL honesto (ADR-17).
- AC-M01-3: job ML/Oracle real disparado pelo scheduler nesta milestone (viola Non-Scope +
  Boundary MC-11 da mission — "MIS-006 ENFILEIRA, não executa").
- AC-M01-4: `root.go` editado fora do bloco de tickers (invasão do wiring de source, propriedade
  de M-02 — quebra o additive-lock).
- AC-M01-5: cadência hardcoded ("diário"/cron fixo) no código do scheduler em vez de vir de
  `schedule` JSONB (viola D6 cadence-agnostic).

## Notas

- M01-C3/C4 exigem dev stack (docker) para rodar um ciclo real contra Postgres — se bloqueado
  no momento da execução, registrar `could-not-run` nomeando o bloqueio (ex.: "REQUEST dev-stack
  ao hub pendente") em vez de pular a prova.
- `last_error` NÃO deve conter payload de provider (AC-02 da mission não se aplica diretamente
  aqui pois M-01 não fala com provider nenhum — mas se um job placeholder futuro escrever erro,
  a mensagem deve ser genérica, nunca vazar credencial/token).
