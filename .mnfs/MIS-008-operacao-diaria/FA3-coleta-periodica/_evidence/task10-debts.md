# Task 10 — Dívidas registradas no fechamento de F-A3

**D-48 — RESOLVIDA** nesta fatia (`formatRelativeAge`).
**D-49 — RESOLVIDA** nesta fatia (um só `FreshnessIndicator`).
**D-50 — RESOLVIDA** nesta fatia (idade nas abas de `/mercado`).
**D-51 — RESOLVIDA** nesta fatia (scheduler de coleta).
**D-52 — asserção de presença em rede de teste de tela.** `AnunciosTable.test.tsx:190` assertava `toBeInTheDocument()` sobre o marcador de frescor e por isso não pegou D-48 por toda a vida do defeito. Corrigida aqui pontualmente. **A classe não foi varrida:** ninguém mediu quantas outras asserções `toBeInTheDocument()` sobre células de valor existem no repo. Varredura é fatia própria.
**D-16 — continua ABERTA e esta fatia não a fecha.** `syncapp.Scheduler.Start` (`sync/application/scheduler.go:105-118`) segue sem execução no boot e sem vencimento persistido. O job de mercado registrado aqui herda o defeito, mas com intervalo de 30min o pior caso é um ciclo perdido por reinício — inofensivo pela letra da própria dívida (`.mnfs/HARNESS-DEBTS.md:301-303`). **A vítima real é `listings` a 24h**, que em dev nunca roda. Fatia própria.
**D-53 — o job renova, não descobre.** Produto sem agregado nunca entra no ciclo periódico; só o clique do operador cria a primeira evidência. Consequência aceita e medida (orçamento de 6 chamadas ML × 2.923 produtos contra bucket de 900/min), não esquecida.

## Achado extra desta fatia (fora do texto do plano, resolvido no mesmo commit)

`sync_state.entity="market"` era compartilhado por dois produtores sem relação: a fila
`MarketEnqueuer` do `erp_import` (nunca chama `RecordSuccess`) e o job periódico desta fatia
(tenant-wide). `/sync/health` devolvia duas linhas com `entity="market"` e sem
`installation_id` no DTO — `SyncHealthCard.tsx` usa `entity` como `key` do React, então a
identidade/ordem de renderização entre as duas linhas era indeterminada (achado do peer
session `local_002d928a`, verificado por mim contra `health_reader.go` e
`SyncHealthCard.tsx` antes de aceitar).

**RESOLVIDO, não registrado como dívida nova**: migração `0093_sync_state_market_queue_entity_split.sql`
renomeia as linhas da fila para `entity="market_queue"`; `sync/domain/sync_state.go` ganha
`EntityMarketQueue`; `sync/composition/scheduler.go` nomeia o sentinela do job
(`InstallationScopeMarket = "market"`, valor inalterado — só o Go-level constant ganhou nome).
Commit `afb6b54a`, mesmo commit da migração — sem janela onde uma tela "verde" mentiria sobre
dado ainda não migrado.

Lane hermética (`-tags=integration` contra `mpc_test_*`) **não foi rodada** para este fix —
fora do escopo do live-drive contra o dev stack compartilhado (banco `marketplace_central`,
não um `mpc_test_*` efêmero). Cobertura obtida: `go build ./...` limpo, testes unitários dos
dois módulos tocados verdes, grep exaustivo confirmando zero literal `'market'` (não-fila)
remanescente nos sítios produtor/consumidor/fixture, e prova ao vivo via `psql` (ver
`task9-live-drive-fatia-b.md`) mostrando as duas entidades já separadas no banco real.
