# M-01-erp-xlsx-identity

```yaml
id: M-01
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-004 mvp-demo — fundação de dados ERP + identidade correta.

## Outcome

Planilha .xlsx do cliente importada com protocolo e relatório de rejeição; produtos com identidade canônica correta (CODPROD/EAN/REFFORN per IC-01); Reader port servindo custo/estoque/identidade do snapshot importado; Oracle intacto como caminho alternativo.

## Why This Milestone Exists

Demo sem Sankhya conectado (P1a); defeito estrutural REFERENCIA≠EAN bloqueia matching honesto (R-01, research §4). Tudo na wave B consome estes dados.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | identity-semantics-fix | `F-01-identity-semantics-fix/feature.md` |
| F-02 | erp-import-module | `F-02-erp-import-module/feature.md` |
| F-03 | reader-adapter-selection | `F-03-reader-adapter-selection/feature.md` |

## Dependencies

Nenhuma de outros milestones (wave A). Base: main ≥ f4612be3. Contratos: IC-01, IC-02.

## Ownership & Concurrency

- Exclusive surfaces: `apps/server_core/internal/modules/erp_import/**` (novo), `modules/internal_read/**` (semântica identidade), `modules/catalog/**` (campos identity), OpenAPI seções `/erp/imports*` E schema de produto catalog (F-01, aditivo), `packages/sdk-runtime/src/erpImport.ts`, tabelas `erp_import_*`.
- Migration block: **0045–0049** (+ bump fixture `runner_test.go`).
- Additive-lock grant (F-01): tipos catalog em `packages/sdk-runtime/src/index.ts` — aditivo-only, mesmo commit do schema OpenAPI, liberado no close do milestone (matriz da missão).
- Predicted seam locks: export 1-linha em `sdk-runtime/src/index.ts` fora do grant acima (hub adjudica); registro módulo em composition root + `contracts/governance/modules.json` via merge do chip (hub).
- Runs in parallel with: M-02, M-03.
- Internal feature DAG: `F-01 ∥ F-02` → F-03.

## Risks

Planilha real diverge do IC-02 (R3 da missão — dry-run pré-demo); semântica REFERENCIA muda leitura Oracle existente (regressão lane live-oracle — testes de contrato cobrem ambas as fontes).

## Done Means

Import de planilha exemplo real completa com protocolo; linha inválida rejeitada com motivo inspecionável; produto sem EAN entra REVIEW-only; `GetCostAsOf`/`GetSellableStock` respondem do snapshot; Oracle path compila e passa lane; dual gate + QA live-drive PASS.

## Handoff

- Current status: planned.
- Next owner: hub (dispatch chip wave A).
- Next action: chip implementa F-01∥F-02 → F-03.
- Required files/evidence: `validation-result.md` deste milestone; relatório de import em evidência QA.
- Blockers or open decisions: none.

## Correction Handoff

n/a (planning).
