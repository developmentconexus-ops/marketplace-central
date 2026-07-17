# Milestone Validation Contract

```yaml
id: M-05
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: milestone
```

## Milestone ID

M-05-anuncios-sinais.

## QA Level

Dual gate frio (Opus + Sol medium) + QA live-drive browser fresh (/anuncios).

## Required Outcome

Sinais competitivos por anúncio no módulo `listings` (via port `market.EvidenceReader`, M-02) + AnunciosPage estendida exibindo sinais honestos com evidência completa.

## Criteria

## Criterion: Listings enriquecidos via port interno
ID: M-05-C01
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: GET /listings (estendido, aditivo) pós-coleta; inspecionar código do enriquecimento
- Expected: campos de sinal por anúncio (nosso preço vs sinal de mercado, posição/vencedor, alvo) vindos do port Go `market.EvidenceReader` — SEM HTTP self-call (grep por chamada HTTP interna a /market = zero); response antigo do /listings preservado (aditivo); anúncio sem vínculo ⇒ sinal ausente com motivo (sem codprod), não zero
- Actual:
- Artifact: `M-05-anuncios-sinais/validation-result.md` §port (response + citação código)
Blocking failure: self-call HTTP p/ market, campo antigo quebrado, ou sinal fabricado p/ anúncio sem vínculo
Blocking failure observed: No
Owner: QA Validator

## Criterion: Evidência IC-03 em cada sinal
ID: M-05-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: GET /listings estendido p/ conjunto com: anúncio com sinal forte, anúncio NO_PRICE_EVIDENCE, anúncio INSUFFICIENT_MARKET
- Expected: cada sinal presente carrega source, fetched_at, n_offers, n_sellers, match_status; estados negativos retornados como estados nomeados (NO_PRICE_EVIDENCE/INSUFFICIENT_MARKET), nunca null silencioso nem 0
- Actual:
- Artifact: `M-05-anuncios-sinais/validation-result.md` §evidencia (response bodies)
Blocking failure: sinal sem campos de evidência ou estado negativo colapsado
Blocking failure observed: No
Owner: QA Validator

## Criterion: AnunciosPage com sinais honestos
ID: M-05-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive em /anuncios
- Expected: colunas/área de sinal por anúncio: nosso preço, sinal de mercado com FreshnessIndicator (idade de fetched_at) e n_offers/n_sellers visíveis; anúncio sem evidência exibe o estado nomeado (copy ADR-06 — sem promessa de preço automático); refresh manual dispara coleta e atualiza idade; light+dark OK
- Actual:
- Artifact: `M-05-anuncios-sinais/validation-result.md` §tela (screenshots light+dark)
Blocking failure: preço de mercado exibido sem evidência/idade, estado negativo como R$0 ou verde enganoso
Blocking failure observed: No
Drive (UI — agent-browser):
- Fixture: stack local com wave A fechada + vínculos M-04 aplicados (≥1 anúncio vinculado com sinal, ≥1 sem evidência); sem auth
- Steps:
  - open http://localhost:5174/anuncios
  - assert text <título do anúncio vinculado>
  - assert text <idade da coleta, ex. "há">
  - click <botão atualizar/refresh sinais>
  - assert text <idade renovada>
- Expected: sinal com evidência completa visível; anúncio sem evidência mostra estado nomeado, não valor
Owner: QA Validator

## Criterion: Ownership limpo
ID: M-05-C04
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: diff do chip vs matriz; lanes L0–L2
- Expected: writes só em `modules/listings/**`, `sdk-runtime/src/listings.ts` (aditivo), páginas Anuncios*/Listing* existentes, OpenAPI `/listings*` aditivo; ZERO migração (join read); lanes verdes
- Actual:
- Artifact: `M-05-anuncios-sinais/validation-result.md` §seams
Blocking failure: migração criada, ou write fora do ownership
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- `M-05-anuncios-sinais/validation-result.md` com seções port, evidencia, tela, seams.
- Dual gate antes do live-drive.

## Blocking Failures

Depende de M-02 (port EvidenceReader publicado) e M-03 (shell) — pré-condição de fixture. Dados de sinal vêm de coleta live (M-02 seam); coleta indisponível ⇒ estados negativos honestos são o comportamento esperado, não falha.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned (P6).
- Next owner: QA Validator (pós-close F-01/F-02 do chip M-05).
- Next action: rodar critérios com fixture wave A + M-04.
- Required files/evidence: `M-05-anuncios-sinais/validation-result.md`.
- Blockers or open decisions: none.
