# Milestone Validation Contract — M-05-listings-fees-divergence

```yaml
id: M-05-VC
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

Verdicts binários. Evidência = caminho inspecionável concreto. Seams: sale_price REAL exige
`?context=channel_marketplace` (fato live, `ml-catalog-offers-pricing-api`) — camada 2 real
é live-driven; divergência é testável hermético (mirror + fixtures).

## Milestone ID

M-05

## QA Level

QA-0

## Required Outcome

/anuncios mostra tarifa REAL por anúncio e divergência de estoque: camada 2 em
`channel_fees` (origem `api_listing_prices`, detail canônico 5 chaves), avaliador de
divergência de estoque (mirror ERP vs `available_quantity`, grão variação), DTO /listings
com `tarifa` + `divergences[]` + `filter.divergentes`, coluna TARIFA + badge na
AnunciosTable.

## Criteria

## Criterion: Camada 2 real com detail canônico COMPLETO
ID: M05-C1
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: ingest de fee de anúncio real + SELECT channel_fees layer=2
- Expected: row com amount + detail com AS 5 CHAVES (percentage_fee, fixed_fee,
  financing_add_on_fee, price_used, listing_type_id — F-r05-5), origem
  `api_listing_prices`, `coletado_em` real; re-sweep re-observa (upsert natural key,
  coletado_em avança)
- Actual:
- Artifact:
Blocking failure: detail parcial (3 chaves), origem errada, ou histórico fantasma
Blocking failure observed: No
Owner: QA Validator

## Criterion: Divergência de estoque nas 2 direções
ID: M05-C2
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: plantar estoque mirror ≠ ML (tolerância 0) → ingest; depois corrigir → ingest
- Expected: row aberta com timestamps dos 2 lados NOT NULL (ADR-10); após convergência,
  `resolved_at` preenchido na MESMA row (one-open-row); anúncio COM variação comparado no
  grão variação (cenário negativo: grão errado NÃO gera divergência falsa); SEM vínculo
  produto↔anúncio → NENHUMA avaliação (IC-02)
- Actual:
- Artifact:
Blocking failure: falso positivo por grão, avaliação sem vínculo, ou resolve ausente
Blocking failure observed: No
Owner: QA Validator

## Criterion: filter.divergentes com fixture >1 página
ID: M05-C3
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: teste de API com fixture >1 página contendo divergentes e não-divergentes
- Expected: `filter.divergentes=true` retorna SÓ divergentes ATRAVESSANDO páginas (count
  == divergentes da fixture, não == página 1)
- Actual:
- Artifact:
Blocking failure: truncamento de página-1 (lição CHIP-MERCADO) ou filtro vazando
não-divergentes
Blocking failure observed: No
Owner: QA Validator

## Criterion: AnunciosTable — TARIFA honesta + badge
ID: M05-C4
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser QA /anuncios com fixture mista (com fee, sem fee, divergente)
- Expected: coluna TARIFA com amount formatado onde há camada 2; anúncio SEM fee observada
  mostra `—` (nunca 0, nunca vazio ambíguo); badge ⚠ SÓ nas linhas divergentes
- Actual:
- Artifact:
Blocking failure: 0 fabricado, ou badge em linha não-divergente
Blocking failure observed: No
Drive (UI — agent-browser; UI criteria only):
- Fixture: tenant com ≥1 anúncio com camada 2, ≥1 sem fee, ≥1 divergência aberta
- Steps:
  - open /anuncios
  - assert text "Tarifa"
  - assert text "—"
- Expected: 3 estados visíveis na mesma tabela (valor, —, badge ⚠), light+dark
Owner: QA Validator

## Criterion: Par OpenAPI+SDK mesmo commit + tsc
ID: M05-C5
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `git show --stat <commit>` do contrato + `npx tsc` no worktree (npm ci na raiz —
  profile §3)
- Expected: OpenAPI (DTO tarifa + divergences[] + param) e `sdk-runtime` no MESMO commit;
  tsc verde; commit de contrato serializado pelo hub (ADR-14 — ≤1 em voo na lane C)
- Actual:
- Artifact:
Blocking failure: par quebrado em commits separados
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- SELECT da row camada 2 salvo com o detail integral (jsonb verbatim).
- Divergência: as DUAS transições com timestamps salvos.
- Screenshots light+dark do drive em `docs/design/evidence/`.

## Blocking Failures

- Detail incompleto = blocking (M05-C1 — M-07 lê essas chaves).
- Falso positivo de divergência por grão = blocking (M05-C2 — R-5).
- Truncamento página-1 = blocking (M05-C3).

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: n/a

## Handoff

- Current status: planned.
- Next owner: hub (lane C, após M-04 close; lock aditivo sobre listings).
- Next action: F-01 ∥ F-02 → F-03.
- Required files/evidence: este arquivo; `M-05/validation-result.md`.
- Blockers or open decisions: none.

## Critérios de user-drive (mandato do operador — obrigatório)

| ID | Critério | Prova mínima inspecionável |
|----|----------|----------------------------|
| M05-U1 | Tarifa de anúncio real na tela CONFERE com o painel ML do mesmo anúncio (spot-check ≥2 anúncios de listing_type distintos) | browser drive + screenshot painel ML lado a lado |
| M05-U2 | Divergência dirigida ida-e-volta: plantar → badge ⚠ aparece; corrigir + re-sync → badge some (sem F5 hack) | browser drive das 2 transições |
| M05-U3 | Filtro "divergentes" na tela retorna só os divergentes (fixture >1 página no stack) | browser drive do filtro + count de conferência |
