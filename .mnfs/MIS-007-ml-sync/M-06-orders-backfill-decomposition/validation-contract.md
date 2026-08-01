# Milestone Validation Contract — M-06-orders-backfill-decomposition

```yaml
id: M-06-VC
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-007
created: 2026-08-01
updated: 2026-08-01
validation_level: QA-0
lifecycle_scope: milestone
base_sha: dd89d4b3
```

Verdicts binários. Evidência = caminho inspecionável concreto. Seams: backfill 12m e
decomposição de pedido real são live-driven (conta do operador; R-4 sale_fee POR UNIDADE é
fato SÓ por medição live — a decomposição TEM que ser conferida contra pedido real).

## Milestone ID

M-06

## QA Level

QA-0

## Required Outcome

Todo pedido dos últimos 12 meses no banco, atualizando a cada 5min, com margem: enumerador
de backfill + scheduler orders 5min chamando IngestOrder do M-03; decomposição canônica
IC-03 persistida (custo CONGELADO na 1ª computação) + camada 3 do ledger; auditoria 3→2
(R$0,01) gera divergência kind=tarifa; /pedidos com coluna MARGEM + drawer.

## Criteria

## Criterion: Backfill 12m + incremental visível
ID: M06-C1
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: backfill 12m na conta real; SELECT sync_state entidade orders
- Expected: pedidos dos 12m no banco (count vs painel ML); watermark persistido;
  `incremental=true` no run incremental (fix M-02 F-03 em uso); 403 de pedido de terceiro
  → skip CONTADO, run segue (fato live)
- Actual:
- Artifact:
Blocking failure: run monolítico não-particionado, 403 abortando run, ou watermark
avançando em run INCOMPLETO (IC-06 run-complete — pedidos sumiriam do incremental p/
sempre)
Blocking failure observed: No
Owner: QA Validator

## Criterion: Kill do scheduler → retomada do cursor
ID: M06-C2
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: kill do processo no meio do backfill; re-boot; logs + SELECT cursor
- Expected: retomada EXATA do cursor persistido (janela/offset — IC-06 canonical), zero
  re-início, zero duplicata de domínio (upsert por resource id)
- Actual:
- Artifact:
Blocking failure: re-início do zero ou duplicata
Blocking failure observed: No
Owner: QA Validator

## Criterion: Decomposição de pedido real consistente
ID: M06-C3
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: pedido real decomposto; SELECT decomposition jsonb + channel_fees
  camada 3
- Expected: liquido == receita_bruta − comissao_total − frete_seller − custo_produto
  (identidade IC-03 exata nos valores persistidos); comissao_total == sale_fee POR
  UNIDADE × qty (R-4); camada 3: valor TOTAL da linha + detail {sale_fee_unit, quantity}
  consistentes entre si; frete camada 3 = custo seller do shipment, detail NULL permitido
  (F-r06-1); custo_produto CONGELADO (re-ingest do pedido NÃO re-computa custo)
- Actual:
- Artifact:
Blocking failure: identidade quebrada, sale_fee tratado como total (regressão R-4), ou
custo re-computado em re-ingest
Blocking failure observed: No
Owner: QA Validator

## Criterion: Honest-unknown na margem
ID: M06-C4
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: fixture: pedido antigo sem custo ERP na 1ª computação
- Expected: `incompleto[]` nomeia `custo_produto`; margem NULL persistida (NUNCA 0);
  /pedidos renderiza `—` na coluna MARGEM
- Actual:
- Artifact:
Blocking failure: margem 0 fabricada onde custo é desconhecido
Blocking failure observed: No
Owner: QA Validator

## Criterion: Auditoria 3→2 nas 2 direções
ID: M06-C5
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: plantar camada 2 ≠ observado no pedido além de R$0,01 → auditoria; convergir →
  re-auditoria
- Expected: row divergence kind=tarifa aberta (timestamps 2 lados NOT NULL); convergência →
  auto-resolve; delta ≤ R$0,01 → NENHUMA divergência (tolerância provada dos 2 lados);
  fórmula esperada = percentage_fee × unit_price/100 + fixed_fee + financing_add_on_fee
  (3 termos — F-r05-2); sem camada 2 → auditoria MUDA (não quebra, edge soft M-05)
- Actual:
- Artifact:
Blocking failure: falso positivo dentro da tolerância, fórmula de 2 termos, ou auditoria
quebrando sem camada 2
Blocking failure observed: No
Owner: QA Validator

## Criterion: /pedidos MARGEM + drawer; par mesmo commit
ID: M06-C6
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser QA /pedidos + `git show --stat` do commit de contrato + tsc
- Expected: coluna MARGEM (% formatado; NULL → `—`); drawer com decomposição linha a linha
  + `incompleto[]` visível quando presente; termos de fee no drawer exibem proveniência
  (origem + coletado_em, camada 3 — ADR-09, P7 r01 B-3); OpenAPI+SDK mesmo commit; tsc verde
- Actual:
- Artifact:
Blocking failure: margem sem drawer explicativo, ou par de contrato quebrado
Blocking failure observed: No
Drive (UI — agent-browser; UI criteria only):
- Fixture: tenant real pós-backfill com ≥1 pedido completo e ≥1 com incompleto[]
- Steps:
  - open /pedidos
  - assert text "Margem"
  - click <linha de pedido decomposto>
  - assert text "Comissão"
  - assert text "api_order"
- Expected: coluna MARGEM na lista; drawer decompõe receita/comissão/frete/custo/líquido;
  proveniência de fee visível no drawer (origem `api_order`/`api_shipment` vocabulário
  IC-01 + coletado_em — ADR-09, B-3)
Owner: QA Validator

## Evidence Requirements

- Conferência do pedido real contra painel ML (screenshot lado a lado) — R-4.
- SELECTs de decomposition + camada 3 salvos com valores.
- Kill/resume com logs timestamped.

## Blocking Failures

- Identidade financeira quebrada ou sale_fee×qty errado = blocking (M06-C3).
- Margem fabricada = blocking (M06-C4).
- Watermark avançando em run incompleto = blocking (M06-C1).

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: n/a

## Handoff

- Current status: planned.
- Next owner: hub (lane C, após M-03).
- Next action: F-01 → F-02 → F-03.
- Required files/evidence: este arquivo; `M-06/validation-result.md`.
- Blockers or open decisions: none.

## Critérios de user-drive (mandato do operador — obrigatório)

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M06-U1 | Margem de pedido REAL na tela confere com conta manual feita sobre o painel ML do mesmo pedido (sale_fee unitário × qty) | browser drive + screenshot painel ML + conta registrada |
| M06-U2 | Drawer de decomposição linha a linha dirigido: cada termo visível e a soma fecha na tela | browser drive do drawer |
| M06-U3 | Pedido sem custo ERP: tela mostra `—` na margem e o drawer nomeia o que falta (incompleto[]) — nunca 0% | browser drive do caso incompleto |
| M06-U4 | Pedido novo criado no ML aparece na tela em ≤5min sem F5 ritual (scheduler incremental vivo) | drive cronometrado + sync_state |
