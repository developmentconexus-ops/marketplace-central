# Milestone Validation Contract — M-04-listings-backfill-ingest

```yaml
id: M-04-VC
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

Verdicts binários. Evidência = caminho inspecionável concreto. Seams: backfill contra a
conta ML REAL é live-driven (hub); R-3 obriga fixtures multi-página nos testes (conta real
pequena esconde defeito de paginação — lição CHIP-MERCADO).

## Milestone ID

M-04

## QA Level

QA-0

## Required Outcome

Catálogo ML completo e fresco no Postgres: backfill scan→multiget retomável por cursor
(IC-06), MASS-CLOSURE substituído por absent≠closed, E3 + `listing_variations` +
`available_quantity` (IC-07), scheduler diário registrado, refresh manual re-apontado pro
mesmo `IngestListing`. Writer único preservado.

## Criteria

## Criterion: Backfill completo + retomada sem duplicata
ID: M04-C1
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive do hub: backfill na conta real; kill do processo no meio; re-boot
- Expected: retomada do cursor persistido (zero re-início, zero duplicata —
  SELECT count por (tenant, provider_listing_id) distinto == count total); cursor terminal
  `{"phase":"sweep",...}` persistido não-nil em sync_state
- Actual:
- Artifact:
Blocking failure: re-início do zero, duplicata, ou cursor terminal nil
Blocking failure observed: No
Owner: QA Validator

## Criterion: Abort pós-página-1 → ZERO rows flipped closed (R-B)
ID: M04-C2
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: must-fail nomeado — fixture multi-página, abort (429 simulado) após página 1
- Expected: nenhuma row viva flipa p/ closed em run INCOMPLETO (run-complete rule IC-06);
  o teste VERMELHO (com a semântica antiga MASS-CLOSURE simulada) nomeia o teste
  (`failure_token=test=`), depois verde com a nova
- Actual:
- Artifact:
Blocking failure: run parcial fechando listings (catalog-wiper) — o defeito capital da
missão (MASS-CLOSURE em pull parcial, audit D-120 F1)
Blocking failure observed: No
Owner: QA Validator

## Criterion: E3 + variations no grão certo
ID: M04-C3
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: ingest de fixture com anúncio COM variações e SEM variações + SELECTs
- Expected: colunas E3 (IC-07) populadas; `listing_variations` com rows por variação;
  `available_quantity` no grão VARIAÇÃO quando existe, no grão listing quando não (IC-07)
- Actual:
- Artifact:
Blocking failure: estoque agregado no grão errado, ou variações perdidas
Blocking failure observed: No
Owner: QA Validator

## Criterion: Âncoras de snapshots não-regressivas (ADR-13)
ID: M04-C4
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: comparação `AbsorbProviderSnapshots` antes/depois da hidratação nova
  (contagem + conteúdo de snapshot por listing)
- Expected: observer de snapshots recebe o MESMO material (contagem e campos) — re-vínculo
  não degrada silencioso
- Actual:
- Artifact:
Blocking failure: snapshots esfomeados pela hidratação multiget
Blocking failure observed: No
Owner: QA Validator

## Criterion: Scheduler diário + refresh manual no mesmo caminho
ID: M04-C5
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: SELECT sync_state (entidade listings) pós-registro + drive do refresh em lote
- Expected: run real de listings registrado em sync_state (scheduler diário); refresh
  manual responde 202 async e executa pelo MESMO `IngestListing` (log/trace do caminho
  único — writer único ADR-04)
- Actual:
- Artifact:
Blocking failure: refresh por caminho paralelo de escrita, ou scheduler NO-OP (lição
MIS-006)
Blocking failure observed: No
Owner: QA Validator

## Criterion: Paginação provada com fixture >1 página
ID: M04-C6
Level: Milestone
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: teste de backfill com fixture ≥2 páginas de scan + multiget (R-3)
- Expected: todas as páginas atravessadas (count ingerido == count da fixture, > page size);
  scroll_id do scan usado até exaustão
- Actual:
- Artifact:
Blocking failure: fixture de 1 página (passe vácuo de paginação — lição CHIP-MERCADO
truncamento silencioso)
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Kill/resume com timestamps do live-drive salvos; SELECTs de contagem salvos.
- Must-fail M04-C2 exige o par vermelho-nomeado → verde.
- Medição do backfill completo (duração, nº requests) salva — insumo do R-1.

## Blocking Failures

- Catalog-wipe em run parcial = blocking (M04-C2).
- Duplicata ou re-início na retomada = blocking (M04-C1).
- Scheduler registrado mas sem run real = blocking (M04-C5 — NO-OP silencioso).

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: n/a

## Handoff

- Current status: planned.
- Next owner: hub (lane B, ∥ M-03; após M-01+M-02 merged).
- Next action: F-01 → F-02 → F-03 → F-04 (cadeia serial R-B).
- Required files/evidence: este arquivo; `M-04/validation-result.md`.
- Blockers or open decisions: none.

## Critérios de user-drive (mandato do operador — obrigatório)

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M04-U1 | /anuncios mostra o catálogo REAL completo pós-backfill: count na tela bate com SELECT count de listings não-closed | browser drive + SELECT de conferência |
| M04-U2 | Refresh manual dirigido na tela: dispara, responde async, e dado atualizado aparece após conclusão (sem F5 infinito nem erro) | browser drive do fluxo de refresh + sync_state atualizado |
| M04-U3 | Anúncio pausado no painel ML mostra o status pausado (lifecycle da row) na tela após sweep; item ABSENT do scan NÃO flipa p/ closed até o sweep confirmar (absent≠closed na prática) | browser drive + row lifecycle no DB (coluna de status antes/depois do sweep) |
