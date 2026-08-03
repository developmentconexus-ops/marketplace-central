# M-06-orders-backfill-decomposition

```yaml
id: M-06
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

MIS-007 ml-sync — [mission.md](../mission.md); IC-03 (decomposição, binding), IC-06
(cursor orders canonical), IC-01 (camada 3 + auditoria 3→2), IC-02 (divergência tarifa).

## Outcome

Todo pedido dos últimos 12 meses no banco, atualizando sozinho a cada 5min, com margem
calculada: enumerador de backfill (/orders/search, cursor canonical IC-06:
backfill/window_until/offset → incremental/watermark−10min) + scheduler de orders 5min
(segunda instância, ADR-08) chamando o IngestOrder do M-03; decomposição financeira por
pedido (JSON canônico IC-03: receita_bruta, comissao_total = sale_fee POR UNIDADE × qty,
frete_seller, custo_produto CONGELADO na 1ª computação, liquido, margem_pct, incompleto[])
persistida + camada 3 do fee ledger (order_line, valor TOTAL + detail
{sale_fee_unit, quantity}); auditoria 3→2 (tolerância R$0,01) gera divergência kind=tarifa;
FE /pedidos ganha coluna MARGEM + drawer de decomposição.

## Why This Milestone Exists

Design §6: margem por pedido é o outcome nº1 da missão. Estende o writer do M-03 (ADR-04 —
NUNCA segundo writer) e consome camada 2 do M-05 p/ auditoria. M-08 (webhook) chama o
IngestOrder estendido daqui.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | backfill-incremental | [F-01-backfill-incremental/feature.md](F-01-backfill-incremental/feature.md) |
| F-02 | decomposition-camada3 | [F-02-decomposition-camada3/feature.md](F-02-decomposition-camada3/feature.md) |
| F-03 | audit-fe-pedidos | [F-03-audit-fe-pedidos/feature.md](F-03-audit-fe-pedidos/feature.md) |

## Dependencies

- M-03 (IngestOrder v1 — F-01 enumera e chama; F-02 estende persistência).
- M-02 (F-03 scheduler fix: incremental flag + cursor não-nil; ports fee/divergence).
- M-05 SOFT (auditoria 3→2 precisa de camada 2 populada — sem ela auditoria fica muda,
  não quebra; lane C permite qualquer ordem de close).

## Ownership & Concurrency

- Exclusive surfaces: `orders/**` (herda posse da lane B — M-03 FECHADO quando lane C abre),
  composição do scheduler de orders (arquivo novo + 1 linha ancorada root.go), par
  OpenAPI+SDK /orders (margem), FE /pedidos (PedidosPage, PedidoDrawer.tsx, fila).
- Migration block: **0094-0095 (reserva — só se a execução descobrir coluna faltante; gap
  numérico ok, já há precedente no repo)**. Esperado: ZERO migração (0088/0089 cobrem).
- Predicted seam locks: channel_fees SÓ camada 3; divergences SÓ kind=tarifa; scheduler.go
  INTOCADO (fix foi M-02); contrato FE serializado na lane C pelo hub (ADR-14).
- Runs in parallel with: M-05, M-07 (código ∥; commits de contrato serializados).
- Internal feature DAG: F-01 → F-02 → F-03.

## Risks

- Backfill 12m ≈ milhares de pedidos × 3-4 GETs — token-bucket do M-01 é o que impede 429
  storm; run particionado por janela (cursor retomável IC-06; nunca run monolítico).
- MASS-CLOSURE não existe em orders (upsert puro) — mas run INCOMPLETO não pode avançar
  watermark (IC-06 run-complete rule; senão pedidos somem do incremental p/ sempre).
- Custo congelado: pedido antigo sem custo ERP na 1ª computação → `incompleto[]` nomeia
  `custo_produto`, margem NULL — nunca 0 (ADR/AGENTS honest-unknown).
- 403 de pedidos de terceiro no search (fato live) → skip contado, run segue.

## Done Means

- Backfill dos 12m completo na conta real (34 anúncios, poucos pedidos — live-drive barato);
  sync_state de orders com watermark e `incremental=true` visível.
- Kill do scheduler no meio do backfill → re-boot retoma do cursor (zero re-início).
- Pedido real: decomposição com margem consistente (liquido = receita_bruta −
  comissao_total − frete_seller − custo_produto, IC-03; sale_fee × qty confere com camada 3
  TOTAL + detail).
- Divergência tarifa plantada (camada 2 ≠ observado no pedido além de R$0,01) → row aberta;
  convergência → auto-resolve.
- /pedidos: coluna MARGEM (% formatado, NULL → `—`), drawer com decomposição linha a linha
  + `incompleto[]` visível; tsc verde; par OpenAPI+SDK mesmo commit.

## Handoff

- Current status: planned. Prep CLOSED: P2 "fios cortados em /pedidos" —
  [_evidence/p2-fios-cortados-CLOSED.md](_evidence/p2-fios-cortados-CLOSED.md), branch
  `worktree-p2-dinheiro-real-pedidos` (7 commits, base `70523e92`, tip `e338b279`), merge
  pendente. Corrige currency/fulfillment/logistic_type/tracking_method/cancellation_detail
  no read model de `/pedidos` — NÃO toca M06-C1..C6 (backfill/decomposição/margem seguem
  Pending; margem já media 33/38 antes deste prep, fato re-medido não construído).
- Next owner: Milestone Orchestrator (hub) — decidir merge do prep, depois lane C após M-03.
- Next action: hub aprova/mergeia o prep; então F-01 spec.
- Required files/evidence: `validation-contract.md` (P6), `validation-result.md`.
- Blockers or open decisions: none.

## Correction Handoff

N/A.
