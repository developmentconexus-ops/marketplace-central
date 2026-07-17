# F-03-identity-resolver

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-02
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-02-price-intel-core.

## Brief

Resolver determinístico de identidade produto↔catálogo-ML no `market`: aplica o gate de matching IC-01 (2 âncoras independentes ⇒ auto-ACCEPT; **negativa dura — kit/combo, cor, medida, voltagem — ⇒ REJECT mesmo com EAN igual**; contradição não-dura vence EAN ⇒ teto REVIEW; sem EAN ⇒ teto REVIEW; nunca auto-ACCEPT só por título) e persiste `market_match_decisions` com âncoras citadas. Alimenta agregados (F-04) — só evidência de match ACCEPT entra em agregado.

## Inputs

- IC-01 — tabela-verdade completa de âncoras/confiança/enums (fonte única).
- F-02 — tabela `market_match_decisions` + repos.
- F-01 — `SearchCatalogByEAN`, `GetCatalogProduct` (candidatos + atributos p/ âncoras).
- Identidade canônica do produto via API catalog (M-01 F-01): `ean`, `refforn`, `marca`, descrição.

## Expected Output

- Serviço `Resolve(productIdentity, catalogCandidates) → MatchDecision` puro/determinístico: mesmas entradas ⇒ mesma decisão (sem LLM, sem rede).
- Âncoras avaliadas: EAN exato, marca normalizada, modelo/refforn, atributos-chave; título é evidência FRACA (nunca âncora sozinha p/ ACCEPT).
- Persistência da decisão com: status (`ACCEPT|REVIEW|REJECT|NO_CANDIDATE`), confiança (≥85 / 50–84 / <50), âncoras a favor, contradições.
- EARS: While produto tem EAN e candidato bate EAN + marca, when Resolve roda, the sistema shall decidir ACCEPT com ambas as âncoras citadas. While EAN bate mas candidato exibe negativa DURA (kit/combo, cor, medida, voltagem divergentes — IC-01), when Resolve roda, the sistema shall decidir REJECT citando a negativa (negativa dura vence EAN). While EAN bate mas marca CONTRADIZ sem negativa dura, when Resolve roda, the sistema shall decidir no máximo REVIEW citando a contradição. While produto sem EAN, when Resolve roda, the sistema shall decidir no máximo REVIEW independente das demais âncoras. While zero candidatos, when Resolve roda, the sistema shall decidir NO_CANDIDATE (não REJECT).

## Negative Scenarios

- Candidato com atributos vazios ⇒ âncora "indisponível", não "a favor" nem "contra".
- Dois candidatos empatados em ACCEPT-elegível ⇒ REVIEW com ambos listados (ambiguidade nunca auto-resolve).
- EAN do produto em colisão (flag M-01 F-01) ⇒ EAN vira âncora inválida; teto REVIEW.

## Constraints

- Determinismo total (tabela-verdade testável).
- Fixtures de colisão IC-01 são o contrato de teste — MESMA tabela-verdade que M-04 F-01 usa p/ ranking de candidatos (implementações independentes, verdade única).
- Sem UI, sem endpoint (F-04 expõe).

## Ownership

- Owned paths: `apps/server_core/internal/modules/market/**` (subpacote resolver + testes).
- Forbidden paths: `modules/connectors/**`, `modules/product_links/**` (M-04), handlers HTTP (F-04), OpenAPI, SDK, migrations novas.
- Parallel-safe with: none — depends on F-02 (tabela decisions) e F-01 (shapes de candidato).

## Validation Expectations

- Suite de fixtures IC-01: ≥12 casos da tabela-verdade (2-âncoras ACCEPT, negativa dura com EAN igual ⇒ REJECT, contradição não-dura ⇒ REVIEW, sem-EAN, título-só, empate, colisão, NO_CANDIDATE...) — cada um com decisão + confiança + âncoras EXATAS esperadas.
- Transcript: mesma entrada 3× ⇒ decisão idêntica (determinismo).
- Linha persistida em `market_match_decisions` com âncoras serializadas inspecionáveis.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-02, após F-01+F-02).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` com suite de fixtures.
- Blockers or open decisions: none.
