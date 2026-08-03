# CLOSED — P2 fios cortados em /pedidos

```yaml
plan: docs/superpowers/plans/2026-08-02-p2-dinheiro-real-pedidos.md
milestone: M-06-orders-backfill-decomposition
parent: MIS-007-ml-sync
branch: worktree-p2-dinheiro-real-pedidos
base_sha: 70523e92
tip_sha: e338b279
commits:
  - 1bfa74d5  evidence(M-06/P2): medicao de premissas + divida D-19 CODEMP
  - 3ba82cd2  fix(orders): status_detail e cancellation_detail viram *string; NULL volta a ser alcancavel
  - f63aa804  fix(mercado_livre): DTO de envio declara logistic.type e tracking_method
  - 29c9d6b8  fix(mercado_livre): DTO de pedido declara cancel_detail e currency_id
  - d6c54fd9  feat(orders): currency e fulfillment ganham produtor real
  - 21aebbb7  feat(orders): cancellation_detail sai no read model; provider_status_detail vira nullable
  - e338b279  contract(orders): currency, fulfillment e motivo de cancelamento no read model
migration: 0093 (NULL alcançável em provider_status_detail/cancellation_detail), 0096 (orders_currency)
review_model: self-check por task (mecânico, bem-especificado) — subagent-driven-development, sem review formal de spec/quality (dispensado pelo operador por escopo estreito)
status: CLOSED. Merge em main ainda PENDENTE (branch não integrada).
```

## Escopo — o que este plano NÃO é

Este NÃO fecha M06-C1..C6 (backfill 12m, incremental 5min, decomposição persistida,
camada 3, auditoria 3→2). Margem já media 33/38 não-nulo ANTES deste plano
(`enrich_service.go:391` → `BuildProfitability`) — fato re-medido, não construído aqui.
P2 é a metade "fio cortado" do M-06: seis campos que o provider manda e o DTO/tipo
descartava — mesma classe de defeito que o P1 fechou em `/anuncios`.

## Medições finais (checklist do plano, seção final)

| # | Critério | Medida |
|---|---|---|
| 1 | currency/fulfillment não-nulos | 0/38 → **38/38 cada** |
| 2 | logistic_type/tracking_method no banco | 0 → **38/38/38** |
| 3 | motivo de cancelamento | **7/7** cancelados com valor real; **31/31** não-cancelados `NULL` (SQL `IS NULL`, 0 `''`) |
| 4 | margem intacta | **33/38**, sem regressão |
| 5 | imposto/DIFAL intactos (fora de escopo) | imposto ≈4%, DIFAL honesto-vazio `{}`; `pricing_calc_profiles` 0 rows — confirmado intocado |
| 6 | `/pedidos` browser drive | pedido `2000017258505630`: "Motivo do cancelamento → vendedor · produto indisponível", "Modalidade de envio → Cross-docking (drop-off)"; console limpo; zero chamadas ao Mercado Livre no `/pedidos` (só `localhost:8080`/`5174`) |

## Dívida registrada

D-19 (`.mnfs/HARNESS-DEBTS.md`, worktree branch): `CODEMP` fixo em 1 no leitor de custo
de pedidos — não bloqueia P2, escopo de M-06/M-07.

## Handoff

- Current status: CLOSED (branch), aguardando merge/repoint definitivo pelo hub.
- Next owner: hub — decidir merge para `main` (branch `worktree-p2-dinheiro-real-pedidos`,
  7 commits, base `70523e92`, sem conflito conhecido — nenhum outro chip tocou
  `orders/**`/`PedidoDrawer.tsx` desde o base_sha).
- Dev stack: backend+frontend repontados para `main` ao final desta sessão (verificação
  usou repoint temporário para o worktree; já revertido).
- Next action: se hub aprovar merge, integrar e então despachar F-01 (backfill-incremental)
  da lane C — M06-C1..C6 continuam Pending, nenhum tocado por este plano.
