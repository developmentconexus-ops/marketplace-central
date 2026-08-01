# Milestone Validation Contract — M-07-pricing-fee-read

```yaml
id: M-07-VC
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

Verdicts binários. Evidência = caminho inspecionável concreto. Seams: resolver é
ledger-only (channel_fees READ + pricing_tariff_defaults) — hermético-testável integral;
before/after de simulações REAIS é live-driven (comportamento MUDA onde degrau-3 vivo
resolvia).

## Milestone ID

M-07

## QA Level

QA-0

## Required Outcome

Pricing lê tarifa do LEDGER: cascata camada 2 → 1 → `config` (pricing_tariff_defaults como
fallback ratificado-honesto, ADR-09); `tarifflive/resolver.go` deletado + composite
simplificado + região root.go:845-851 refiada (`:844 var tariffResolver` SOBREVIVE); DTOs
/pricing com proveniência (origem + coletado_em); /precos mostra origem com ⚠ quando
`config`; allowlist -2 (C/D).

## Criteria

## Criterion: Grep zero tarifflive + imports mortos
ID: M07-C1
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: `grep -rn "tarifflive" apps/server_core/` + build verde
- Expected: 0 hits (arquivo deletado, imports root.go:99,101 removidos, composite
  simplificado); `var tariffResolver` (`root.go:844`) e `batchOrch.WithTariffResolver`
  (`:852`) vivos com o resolver NOVO
- Actual:
- Artifact:
Blocking failure: referência sobrevivente a tarifflive, ou dangling brace da região
845-851 (F-r07-5)
Blocking failure observed: No
Owner: QA Validator

## Criterion: Cascata com proveniência (truth table)
ID: M07-C2
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: truth table do resolver (fixtures nomeadas): camada 2 presente / só camada 1 /
  ledger vazio / store de defaults FALHA
- Expected: camada 2 → valor do detail (percentage_fee/fixed_fee) + origem
  `api_listing_prices` + coletado_em da row; camada 1 → origem da PRÓPRIA row de fixture
  (camada 1 não tem produtor nesta missão — IC-01; o teste PLANTA a row, logo a origem
  esperada é fixture-defined dentro do CHECK, P7 r02 A-4); vazio → defaults com origem
  `config` (13.00/16.00 materialize-on-read
  `calc_repository.go:239` INTOCADO); store FALHA → erro nomeado, NUNCA valor silencioso
- Actual:
- Artifact:
Blocking failure: valor sem origem, origem errada por degrau, ou fallback silencioso em
falha de store
Blocking failure observed: No
Owner: QA Validator

## Criterion: Os DOIS caminhos (single + batch) no resolver novo
ID: M07-C3
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: teste dos 2 orquestradores (simulação single + batch via
  `WithTariffResolver`)
- Expected: ambos consomem o MESMO resolver; resultados idênticos p/ mesmo input; zero
  chamada ML em qualquer um (contador de adapter = 0)
- Actual:
- Artifact:
Blocking failure: batch ainda no caminho velho, ou divergência single×batch
Blocking failure observed: No
Owner: QA Validator

## Criterion: Allowlist -2 (C/D) com must-fail
ID: M07-C4
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: guard M-02 F-04 pós-merge + must-fail reintroduzindo chamada ML read-time em
  pricing
- Expected: entradas C/D removidas; reintrodução simulada REPROVA nomeando o sítio (output
  vermelho salvo)
- Actual:
- Artifact:
Blocking failure: guard aceitando sítio de pricing ressuscitado
Blocking failure observed: No
Owner: QA Validator

## Criterion: Before/after de simulações reais explicado
ID: M07-C5
Level: Milestone
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: N≥3 simulações reais pré-merge (baseline salvo) e pós-merge —
  incluindo ≥1 anúncio COM camada 2 e ≥1 sem ledger (degrau config), P7 r02 A-7
- Expected: CADA delta explicado e classificado (live→camada2 ou live→config); nenhum
  delta órfão
- Actual:
- Artifact:
Blocking failure: delta sem explicação nomeada
Blocking failure observed: No
Owner: QA Validator

## Criterion: /precos com proveniência visível
ID: M07-C6
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser QA /precos; tsc; `git show --stat` do commit de contrato
- Expected: origem + coletado_em na tela por valor de tarifa; ⚠ quando origem=config;
  NUNCA valor sem origem; par OpenAPI+SDK mesmo commit; tsc verde
- Actual:
- Artifact:
Blocking failure: número de tarifa sem origem na tela (a doença que ADR-09 mata)
Blocking failure observed: No
Drive (UI — agent-browser; UI criteria only):
- Fixture: tenant com ≥1 anúncio com camada 2 e ≥1 produto sem ledger (resolve config)
- Steps:
  - open /precos
  - click <simular no produto com camada 2>
  - assert text "api_listing_prices"
  - click <simular no produto sem ledger>
  - assert text "⚠"
- Expected: origem visível nos 2 casos; ⚠ SÓ no caso config
Owner: QA Validator

## Evidence Requirements

- Truth table com fixtures NOMEADAS e outputs salvos.
- Baseline before/after salvo ANTES do merge (irrecuperável depois).
- Screenshots light+dark do /precos.

## Blocking Failures

- Valor sem origem = blocking (M07-C2/C6).
- Chamada ML em simulação = blocking (M07-C3 — allowlist).
- Delta órfão no before/after = blocking (M07-C5).

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: n/a

## Handoff

- Current status: planned.
- Next owner: hub (lane C — pode ser o primeiro da lane).
- Next action: F-01 → F-02; baseline before ANTES do merge.
- Required files/evidence: este arquivo; `M-07/validation-result.md`.
- Blockers or open decisions: none.

## Critérios de user-drive (mandato do operador — obrigatório)

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M07-U1 | Simulação dirigida de produto COM camada 2: tarifa vem do ledger com origem+coletado_em na tela, e o valor bate com a coluna TARIFA de /anuncios do mesmo anúncio | browser drive 2 telas + SELECT camada 2 |
| M07-U2 | Simulação de produto SEM ledger: ⚠ config visível, valor = defaults do tenant (13/16), sem erro | browser drive do caso fallback |
| M07-U3 | Network do browser durante simulações: ZERO request ML (sites C/D mortos de fato no fio) | network capture do drive |
