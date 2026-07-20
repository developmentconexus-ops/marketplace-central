# Mission Validation Contract — MIS-006-integracao-fundacao

```yaml
id: MIS-006-VC
type: mission-validation-contract
mission: MIS-006
validation_level: QA-0
base_sha: 138aac3d
```

Verdicts binários por critério. Evidência = caminho inspecionável concreto (core §5, ladder L0-L4).
Tipos de evidência: `ran` (executado, output salvo), `assumed` (design, não executado), `could-not-run`
(bloqueado — nomear bloqueio). Honestidade de integração: seam contra dependência REAL nunca é
provado por stub.

## Critérios de missão

| ID | Critério | Prova mínima inspecionável | Ladder | Milestone dono |
|----|----------|----------------------------|--------|----------------|
| MC-01 | `products_mirror` existe e é alimentado por AMBOS adapters via um port | Migração aplicada; após import xlsx E sync Sankhya (dev stack), `SELECT ... FROM products_mirror WHERE source IN ('xlsx','sankhya')` retorna rows dos dois; consumidor lê mirror não rescan | L2 ran | M-02,M-03,M-04 |
| MC-02 | Upsert-merge keep-absent: produto ausente de snapshot novo mantém row flagada, nunca delete físico | IO Fase 5: importa snapshot sem produto X → `absent_in_last_snapshot=true`+`stale_since` set, row presente; `product_links` de X intacto | L2 ran | M-03 |
| MC-03 | NULL honesto: prospect #004-E custo/estoque ausente grava NULL, nunca 0 | `SELECT custo,estoque_total FROM products_mirror WHERE source='xlsx' AND custo IS NULL` retorna rows do prospect | L2 ran | M-02,M-03 |
| MC-04 | Fonte ativa por tenant em banco; `MC_ERP_SOURCE` removido; sem poluição de cache cross-source | grep `MC_ERP_SOURCE`=0 hits no código; toggle por tenant nas 2 direções retorna snapshot certo; must-fail test prova cache-key inclui source | L2 ran + L1 | M-02 |
| MC-05 | Cadeia de vínculo automática: import dispara geração de candidatos; EAN-exato-único auto-aprova com audit | IO A (90008 EAN único): pós-import, candidato existe E `product_links` aprovado + `E10 audit(actor=system,rule=exact_ean_unique)`; EAN duplicado → REVIEW | L2 ran | M-05 |
| MC-06 | Idempotência: re-run não duplica vínculo; override manual vence | 2× sync mesmo produto → 1 link; operador sobrescreve → nova audit row + `superseded_by` | L1 ran | M-05 |
| MC-07 | sync_state rastreia entidades cadence-agnostic; scheduler skeleton roda loop sem job ML pesado | `sync_state` populado pós-import (entity=products, cursor, last_full_sync_at); scheduler loop registrado em root.go; nenhum job ML disparado | L1 ran | M-01 |
| MC-08 | Tela /importacoes mostra cadeia real N-imported→N-linked→N-enqueued | P7 browser QA: /importacoes renderiza contagens vindas de join sync_state+links, não hardcode; light+dark | L3 ran | M-06 |
| MC-09 | SankhyaAdapter validado contra Oracle REAL (não stub) | `[TESTAR-SKW]` mapping confirmado via db-consult; M-04 smoke contra Oracle dev stack; se bloqueado → `could-not-run` nomeado | L2 ran/could-not-run | M-04 |
| MC-10 | F3.7: prova live T13-T16 decide construir-ou-remover; sem execução de coleta | evidência `docs/design/evidence/ml-api/`; se viável M-07 enfileira (não executa); se não, decisão REMOVE registrada honest-unknown | L2 ran/could-not-run/assumed | M-07 |
| MC-11 | Boundary respeitado: MIS-006 ENFILEIRA mercado, não executa | nenhum write a `market_aggregates`/`competitor_offers` no diff da missão; só enqueue em sync_state | L0 ran (grep diff) | todos |
| MC-12 | Contratos aditivos: nenhuma coluna existente muda tipo/semântica; SDK+OpenAPI no mesmo commit | diff de migração = só ADD; `sdk-runtime` regenerado junto do endpoint | L1 ran | M-02,M-06 |

## Anti-critérios (falha se presente)

- AC-01: qualquer query tenant-scoped sem `tenant_id`. 
- AC-02: payload de provider fora do adapter.
- AC-03: unknown vira 0/default em vez de NULL (viola ADR-17).
- AC-04: stub no lugar de prova de integração viva sem autorização explícita.
- AC-05: `.env*` lido/impresso; credencial ML/Oracle exposta.
- AC-06: push para remote (profile §9).

## Gate de prontidão (P7)

Missão vira `planned` só com verdicts model-side Claude=Ready E Sol=Ready sobre o MESMO digest
de manifesto. Sol rebind→Claude cold crew enquanto quota-wall codex (até 2026-07-25); se rebind
não-autorizado pelo hub → `blocked`, nomear, escalar (nunca skip para planned).
