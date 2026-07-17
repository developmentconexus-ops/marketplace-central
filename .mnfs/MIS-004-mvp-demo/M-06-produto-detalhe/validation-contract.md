# Milestone Validation Contract

```yaml
id: M-06
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

M-06-produto-detalhe.

## QA Level

Dual gate frio (Opus + Sol medium) + QA live-drive browser fresh (/catalogo/produtos/:productId).

## Required Outcome

Página Produto Detalhe: header de identidade/custo/estoque, box de veredicto "vale a pena vender?" com coleta síncrona + refetch, abas Anúncios vinculados e Estoque — leitura composta client-side via SDK (sem backend novo).

## Criteria

## Criterion: Header de identidade honesto
ID: M-06-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive em produto importado com EAN válido e outro com `ean: null`
- Expected: header exibe CODPROD, EAN, REFFORN, marca, NCM, custo (com source time do import) e estoque; campo null renderizado como UnknownValue/ausente (nunca string vazia ou 0); disponível DESCONHECIDO quando reservado ausente (físico segue exibido como físico)
- Actual:
- Artifact: `M-06-produto-detalhe/validation-result.md` §header (screenshots dos 2 produtos)
Blocking failure: null renderizado como vazio/0, ou disponível fabricado
Blocking failure observed: No
Owner: QA Validator

## Criterion: Box veredicto com coleta síncrona + refetch
ID: M-06-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: acionar coleta no box de veredicto
- Expected: botão dispara POST /market/collections (síncrono, IC-03); durante execução UI mostra estado de espera delimitado; ao 200, refetch de GET /market/verdicts?codprod= atualiza o box SEM polling e sem F5; veredicto exibido ∈ estados ADR-06 com evidência (source, fetched_at, n_offers/n_sellers, match_status) e blocking_state quando aplicável (incl. SEM_CUSTO)
- Actual:
- Artifact: `M-06-produto-detalhe/validation-result.md` §veredicto (screenshot + network transcript)
Blocking failure: polling/job implementado, box atualizando só com F5, ou veredicto sem evidência
Blocking failure observed: No
Drive (UI — agent-browser):
- Fixture: stack local com waves A+B fechadas; produto com vínculo e outro sem evidência; sem auth
- Steps:
  - open http://localhost:5174/catalogo/produtos/<codprod-vinculado>
  - assert text "vale a pena"
  - click <botão coletar/atualizar veredicto>
  - assert text <fetched_at/idade renovada>
  - open http://localhost:5174/catalogo/produtos/<codprod-sem-evidencia>
  - assert text "NO_PRICE_EVIDENCE"
- Expected: veredicto novo sem reload manual; produto sem evidência exibe estado nomeado, nunca R$0/verde
Owner: QA Validator

## Criterion: Aba Anúncios vinculados
ID: M-06-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: aba Anúncios vinculados de produto com ≥1 vínculo aplicado (M-04)
- Expected: lista os anúncios vinculados via links resolvidos + campos de sinal por anúncio (M-05 `listings.ts`); anúncio com sinal exibe FreshnessIndicator; produto sem vínculo ⇒ aba com estado vazio explicativo (call-to-action p/ /vinculos), não tabela em branco
- Actual:
- Artifact: `M-06-produto-detalhe/validation-result.md` §anuncios (screenshots)
Blocking failure: vínculo aplicado não aparecendo, ou estado vazio mudo
Blocking failure observed: No
Owner: QA Validator

## Criterion: Aba Estoque
ID: M-06-C04
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: aba Estoque
- Expected: físico, reservado e disponível do snapshot (Reader M-01) com source time do import; reservado ausente ⇒ disponível DESCONHECIDO via UnknownValue; abas fora do escopo (Concorrência/Pedidos/Histórico) AUSENTES (Non-Scope MIS-005 M-05)
- Actual:
- Artifact: `M-06-produto-detalhe/validation-result.md` §estoque (screenshot)
Blocking failure: disponível numérico fabricado, ou aba fora de escopo criada
Blocking failure observed: No
Owner: QA Validator

## Criterion: Ownership limpo — sem backend novo
ID: M-06-C05
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: diff do chip vs matriz; lanes L0–L2
- Expected: writes SÓ em `apps/web/src/pages/produto/**` + `routes/produto.tsx` (swap de placeholder); ZERO migração, ZERO OpenAPI, ZERO módulo Go; leituras compostas client-side via SDK existente (catalog/inventory/listings/market); light+dark OK
- Actual:
- Artifact: `M-06-produto-detalhe/validation-result.md` §seams
Blocking failure: endpoint/módulo backend novo criado, ou write fora do ownership
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- `M-06-produto-detalhe/validation-result.md` com seções header, veredicto, anuncios, estoque, seams.
- Screenshots light+dark da página completa.
- Dual gate antes do live-drive.

## Blocking Failures

Wave C: depende de M-01/M-02/M-04/M-05 fechados (fixture). Coleta ao vivo via seam M-02 — indisponibilidade ⇒ estados negativos honestos esperados.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned (P6).
- Next owner: QA Validator (pós-close F-01 do chip M-06).
- Next action: rodar critérios com fixture waves A+B.
- Required files/evidence: `M-06-produto-detalhe/validation-result.md`.
- Blockers or open decisions: none.
