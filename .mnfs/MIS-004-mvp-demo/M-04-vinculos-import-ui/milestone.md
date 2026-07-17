# M-04-vinculos-import-ui

```yaml
id: M-04
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

MIS-004 mvp-demo — vínculo produto↔anúncio operável (backend gaps + tela).

## Outcome

`/vinculos` funcional: fila de resolução com candidatos rankeados (confiança + motivo por âncoras IC-01), decisão individual via drawer, bulk aprovação via batch-preview (dry-run)→batch apply LOCAL no product_links com auditoria própria (sem envelope mutations — P5-F-12; sem write ML), estados resolvidos com undo. Placeholder W1 substituído via `routes/vinculos.tsx`.

## Why This Milestone Exists

Sem vínculo resolvido não há sinal de mercado confiável por anúncio (M-05) nem margem por pedido (M-08). Backend product_links existe (candidates/generations/resolutions/workflows) mas sem campos de confiança/motivo nem bulk.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | product-links-api-gaps | `F-01-product-links-api-gaps/feature.md` |
| F-02 | vinculos-screen | `F-02-vinculos-screen/feature.md` |

## Dependencies

- M-01 (identidade EAN correta alimenta candidatos) — edge de DADO, não de código: F-01 pode iniciar contra IC-01/IC-02.
- M-03 (shell/tokens/`routes/vinculos.tsx` seam) para F-02.
- Contratos: IC-01 (âncoras/confiança/enums), IC-05 (seams FE).

## Ownership & Concurrency

- Exclusive surfaces: `modules/product_links/**`, OpenAPI seção `/product-links/*` (aditivo), `sdk-runtime/src/productLinks.ts`, `apps/web/src/routes/vinculos.tsx` (conteúdo pós-seam), `apps/web/src/pages/vinculos/**` (novo), tabelas `product_link_*` (ALTER aditivo).
- Migration block: **0065–0069** (+ bump fixture).
- Predicted seam locks: export barrel SDK (hub). Envelope `mutations` NÃO é usado (batch local — P5-F-12); `modules/mutations/**` segue forbidden.
- Runs in parallel with: M-05, M-07, M-08 (wave B).
- Internal feature DAG: F-01 → F-02 (F-02 inicia contra OpenAPI de F-01 assim que a seção estiver comitada).

## Risks

Candidatos reais escassos (installation com poucos anúncios) — QA precisa de fixture de candidatos determinística; ranking de confiança divergir do IC-01 (mitigação: fixtures de colisão compartilhadas com M-02 F-03 via IC-01, implementações independentes mas mesma tabela-verdade).

## Done Means

Fila mostra candidato com % + motivo citando âncoras; aprovar individual resolve e some da fila; bulk N itens gera preview dry-run e aplica localmente com applied/failed itemizados (auditoria do batch inspecionável); undo reverte; produto sem candidato mostra NO_CANDIDATE honesto; deep-link `/vinculos` + F5 funciona; dual gate + QA PASS.

## Handoff

- Current status: planned.
- Next owner: hub (dispatch chip wave B).
- Next action: chip implementa F-01 → F-02.
- Required files/evidence: `validation-result.md`; screenshot fila + bulk preview.
- Blockers or open decisions: none.

## Correction Handoff

n/a (planning).
