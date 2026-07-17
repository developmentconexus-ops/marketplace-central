# Milestone Validation Contract

```yaml
id: M-04
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

M-04-vinculos-import-ui.

## QA Level

Dual gate frio (Opus + Sol medium) + QA live-drive browser fresh (tela Vínculos).

## Required Outcome

Gaps de `product_links` fechados (batch preview/apply LOCAL com auditoria — exceção P5-F-12) + tela Vínculos & Importação funcional consumindo produtos importados (M-01) e candidatos do resolver (M-02).

## Criteria

## Criterion: Batch preview dry-run local
ID: M-04-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: POST batch-preview de vínculos (endpoint local do product_links) com seleção contendo: candidato ACCEPT, candidato REVIEW, candidato já vinculado
- Expected: response dry-run SEM persistir nada (SELECT antes/depois idêntico): lista por item com ação prevista (aplicar/pular) + motivo; item já vinculado ⇒ pulado com motivo de duplicata; NENHUM protocolo em /mutations criado (batch é estado LOCAL, fora do envelope)
- Actual:
- Artifact: `M-04-vinculos-import-ui/validation-result.md` §preview (response + SELECTs)
Blocking failure: preview persistindo estado, ou intent aparecendo na fila /mutations
Blocking failure observed: No
Owner: QA Validator

## Criterion: Batch apply com auditoria e undo
ID: M-04-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: POST batch apply da mesma seleção; SELECT das tabelas product_links + auditoria; exercitar undo do lote
- Expected: vínculos aplicados com confiança/motivo/lote registrados; entrada de auditoria por item (quem/quando/lote; "quem" = ator fixo `operator` no MVP sem auth); undo do lote reverte TODOS os itens do lote e registra auditoria da reversão; writes provider = zero (log adapter limpo na janela)
- Actual:
- Artifact: `M-04-vinculos-import-ui/validation-result.md` §apply (SELECTs + log)
Blocking failure: apply sem trilha de auditoria, undo parcial, ou qualquer write provider
Blocking failure observed: No
Owner: QA Validator

## Criterion: Tela Vínculos — fluxo preview→apply→Resolvidos
ID: M-04-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive em /vinculos
- Expected: pendentes listados com candidato+confiança+motivo do resolver (IC-01); seleção múltipla → botão de preview abre resumo dry-run; confirmar aplica; itens aplicados saem de Pendentes e aparecem em Resolvidos SEM F5 manual (refetch pós-apply); estados honestos p/ item sem candidato (NO_CANDIDATE visível, não linha vazia)
- Actual:
- Artifact: `M-04-vinculos-import-ui/validation-result.md` §tela (screenshots light+dark + transcript)
Blocking failure: apply sem etapa de preview, Resolvidos não refletindo sem reload, ou NO_CANDIDATE renderizado como vazio
Blocking failure observed: No
Drive (UI — agent-browser):
- Fixture: stack local com import M-01 concluído + coleta M-02 rodada (candidatos existem); sem auth
- Steps:
  - open http://localhost:5174/vinculos
  - assert text "Pendentes"
  - click <checkbox primeiro candidato ACCEPT>
  - click <botão pré-visualizar/aplicar seleção>
  - assert text <resumo dry-run com contagem>
  - click <confirmar>
  - assert text "Resolvidos"
- Expected: item confirmado listado em Resolvidos com lote/motivo; Pendentes decrementado
Owner: QA Validator

## Criterion: Import na tela — protocolo visível
ID: M-04-C04
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: seção Importação da tela exibe último import (M-01)
- Expected: protocolo `#NNN-E`, status COMPLETED|REJECTED, contagens importadas/rejeitadas e link/expansão do relatório de rejeição por linha — dados do GET /erp/imports (sem UI de upload, conforme Non-Scope: import via endpoint/runbook)
- Actual:
- Artifact: `M-04-vinculos-import-ui/validation-result.md` §import (screenshot)
Blocking failure: rejeições invisíveis na UI, ou upload UI criada fora do escopo
Blocking failure observed: No
Owner: QA Validator

## Criterion: Migrações e seams
ID: M-04-C05
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `ls apps/server_core/migrations/ | grep -E '^006[5-9]'` + `runner_test.go`; diff vs ownership
- Expected: ALTERs `product_links_*` (confiança/motivo/lote/undo/auditoria) SÓ no bloco 0065–0069; fixture = contagem real; diff só em `modules/product_links/**`, `sdk-runtime/src/productLinks.ts`, OpenAPI `/product-links/*`, `apps/web/src/pages/vinculos/**`, `routes/vinculos.tsx`
- Actual:
- Artifact: `M-04-vinculos-import-ui/validation-result.md` §seams
Blocking failure: migração fora do bloco ou write fora do ownership
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- `M-04-vinculos-import-ui/validation-result.md` com seções preview, apply, tela, import, seams.
- Log adapter na janela do C02 (prova zero provider writes).
- Dual gate antes do live-drive.

## Blocking Failures

Depende de M-01 (produtos importados) e M-02 (candidatos) fechados — pré-condição de fixture, não seam externo próprio. Batch é LOCAL por contrato (exceção P5-F-12): item aparecendo em /mutations = blocking.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned (P6).
- Next owner: QA Validator (pós-close F-01/F-02 do chip M-04).
- Next action: rodar critérios com fixture wave A fechada.
- Required files/evidence: `M-04-vinculos-import-ui/validation-result.md`.
- Blockers or open decisions: none.
