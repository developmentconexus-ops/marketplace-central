# P5 Sol Decomposition Audit — r02 (VERBATIM)

```yaml
id: P5-SOL-AUDIT-R02
type: planning-review
status: complete
owner: GPT-5.6 Sol (medium) via codex exec, OS-process dispatch
parent: MIS-004
created: 2026-07-17
lifecycle_scope: support
```

Provenance: codex exec --model gpt-5.6-sol -c model_reasoning_effort=medium --sandbox read-only,
prompt scratchpad/prompt-p5-sol-r02.md, log scratchpad/agent__p5-sol-r02.log,
manifest p5-input-r02.sha256 top digest 7ed7b63799ebe9185d40c7aeaf03156bbdfb3e6c0e3362eaa285c1ee4795a508.
Stdout below is UNEDITED.

---
## Checked paths

Frozen manifest:

- `planning-reviews/p5-input-r02.sha256`

Mission planning artifacts:

- `planning-reviews/p5-sol-decomposition-audit-r01.md`
- `planning-reviews/p5-reconciliation-r01.md`
- `mission.md`
- `planning-reviews/p5-decomposition-passes-r01.md`
- `research/identity-matching-interface-contract.md`
- `research/erp-xlsx-import-interface-contract.md`
- `research/market-evidence-read-interface-contract.md`
- `research/pricing-difal-interface-contract.md`
- `research/fe-shell-seams-interface-contract.md`
- `research/ml-read-ports-interface-contract.md`
- `research/difal-interna-rates-2026.md`
- `research/design-screens-2026-07-17.md`
- `research/p1-clarified-decisions-2026-07-17.md`
- `research/repo-baseline-2026-07-17.md`
- `research/w1-merge-addendum-2026-07-17.md`
- `M-01-erp-xlsx-identity/milestone.md`
- `M-01-erp-xlsx-identity/F-01-identity-semantics-fix/feature.md`
- `M-01-erp-xlsx-identity/F-02-erp-import-module/feature.md`
- `M-01-erp-xlsx-identity/F-03-reader-adapter-selection/feature.md`
- `M-02-price-intel-core/milestone.md`
- `M-02-price-intel-core/F-01-ml-adapter-read-ports/feature.md`
- `M-02-price-intel-core/F-02-market-persistence/feature.md`
- `M-02-price-intel-core/F-03-identity-resolver/feature.md`
- `M-02-price-intel-core/F-04-collect-verdict-api/feature.md`
- `M-03-shell-retheme/milestone.md`
- `M-03-shell-retheme/F-01-theme-tokens-fonts/feature.md`
- `M-03-shell-retheme/F-02-header-nav-routes/feature.md`
- `M-03-shell-retheme/F-03-shared-primitives/feature.md`
- `M-04-vinculos-import-ui/milestone.md`
- `M-04-vinculos-import-ui/F-01-product-links-api-gaps/feature.md`
- `M-04-vinculos-import-ui/F-02-vinculos-screen/feature.md`
- `M-05-anuncios-sinais/milestone.md`
- `M-05-anuncios-sinais/F-01-listings-signals-api/feature.md`
- `M-05-anuncios-sinais/F-02-anuncios-ui-sinais/feature.md`
- `M-07-simulador/milestone.md`
- `M-07-simulador/F-01-pricing-calc-difal-service/feature.md`
- `M-07-simulador/F-02-simulador-ui/feature.md`
- `M-08-pedidos/milestone.md`
- `M-08-pedidos/F-01-orders-projection-api/feature.md`
- `M-08-pedidos/F-02-pedidos-ui/feature.md`
- `M-06-produto-detalhe/milestone.md`
- `M-06-produto-detalhe/F-01-produto-detalhe-page/feature.md`
- `M-09-dashboard-demo/milestone.md`
- `M-09-dashboard-demo/F-01-dashboard-mpc/feature.md`
- `planning-reviews/p3-reconciliation-r01.md`

Repository code-fact verification:

- `packages/sdk-runtime/src/index.ts`
- `contracts/api/marketplace-central.openapi.yaml`
- `apps/server_core/internal/modules/mutations/adapters/listings/selection_resolver.go`
- `apps/server_core/internal/composition/root.go`
- `apps/server_core/internal/platform/migrate/runner_test.go`
- `apps/server_core/migrations/` — filename inventory only

## Fold verification

| r01 finding id | Status |
| --- | --- |
| P5-F-01 | CLOSED |
| P5-F-02 | CLOSED |
| P5-F-03 | CLOSED |
| P5-F-04 | CLOSED |
| P5-F-05 | NOT-CLOSED — `M-01-erp-xlsx-identity/F-01-identity-semantics-fix/feature.md`, `## Ownership` |
| P5-F-06 | NOT-CLOSED — `M-01-erp-xlsx-identity/F-01-identity-semantics-fix/feature.md`, `## Expected Output` and `## Negative Scenarios` |
| P5-F-07 | CLOSED |
| P5-F-08 | CLOSED |
| P5-F-09 | CLOSED |
| P5-F-10 | NOT-CLOSED — `M-02-price-intel-core/F-04-collect-verdict-api/feature.md`, `## Expected Output` |
| P5-F-11 | NOT-CLOSED — `research/market-evidence-read-interface-contract.md`, `## Operations`, versus M-02 F-04 |
| P5-F-12 | NOT-CLOSED — `mission.md`, `## Outcome`, `## Domain Scope`, `### Runtime Contract`, ADR-08, versus the M-04 ownership row |
| P5-F-13 | CLOSED |
| P5-F-14 | CLOSED |
| P5-F-15 | CLOSED |
| P5-F-16 | NOT-CLOSED — `M-07-simulador/F-02-simulador-ui/feature.md`, `## Expected Output` |
| P5-F-17 | NOT-CLOSED — `M-07-simulador/F-01-pricing-calc-difal-service/feature.md` and IC-04 decomposition |
| P5-F-18 | CLOSED |
| P5-F-19 | CLOSED |

## Findings

### P5-R2-01

- Severity: `BLOCKING`
- Check: A/P5-F-05; B/six-axis ownership and contract satisfiability.
- Cited excerpt: “incluindo a seção OpenAPI do schema de produto e os tipos catalog do SDK” but ownership says “Owned paths: `...internal_read/**`, `...catalog/**`” and “Forbidden paths: ... `sdk-runtime/**` (F-02).”
- Exact defect locus: `M-01-erp-xlsx-identity/F-01-identity-semantics-fix/feature.md`, `## Brief` and `## Ownership`.
- Offending token/value: forbidden `sdk-runtime/**`, with no owned catalog OpenAPI section.
- Defect: the prose fold assigns catalog OpenAPI/SDK parity to F-01, but its executable write-set still forbids the SDK and does not own the OpenAPI catalog schema. F-02 owns only `/erp/imports*` and `erpImport.ts`. Repository inspection confirms the existing catalog types remain in `packages/sdk-runtime/src/index.ts`.
- Yes-if: F-01’s ownership explicitly includes the catalog-product OpenAPI schema and the additive catalog types in `packages/sdk-runtime/src/index.ts`, removes the contradictory SDK prohibition, and retains ADR-12 same-commit parity.

### P5-R2-02

- Severity: `BLOCKING`
- Check: A/P5-F-06; C/removed identity decision.
- Cited excerpts: “REFERENCIA inválida ⇒ `ean: null` + warning ... nunca migra p/ outro campo” is contradicted by “Reader Oracle mapeia ... inválida → `refforn`,” “REFERENCIA falha checksum ... REFERENCIA em `refforn`,” and “Checksum inválido ... ⇒ vai a `refforn`.”
- Exact defect locus: `M-01-erp-xlsx-identity/F-01-identity-semantics-fix/feature.md`, `## Expected Output` and `## Negative Scenarios`.
- Offending token/value: invalid `REFERENCIA → refforn`.
- Defect: three operative requirements still preserve the exact reversed mapping rejected in r01 and contradict both the amended Brief and IC-01.
- Yes-if: every invalid REFERENCIA case yields `ean: null` plus the contracted warning, and `refforn` is populated exclusively from TGFPRO.REFFORN.

### P5-R2-03

- Severity: `BLOCKING`
- Check: A/P5-F-10; B/contract propagation of exact parameter names.
- Cited excerpts: F-04 specifies “`GET /market/aggregates?codprods=`” and “`GET /market/verdicts?codprods=`”; IC-03 specifies “`GET /market/aggregates?codprod=`” and “`GET /market/verdicts?codprod=`.”
- Exact defect locus: `M-02-price-intel-core/F-04-collect-verdict-api/feature.md`, `## Expected Output`; `research/market-evidence-read-interface-contract.md`, `## Operations`.
- Offending token/value: `codprods` instead of `codprod`.
- Defect: the fold changed away from `product_ids` but did not adopt IC-03’s exact query key.
- Yes-if: producer, OpenAPI/SDK plan, and all HTTP consumers use IC-03’s exact `codprod` query parameter and identifier semantics.

### P5-R2-04

- Severity: `BLOCKING`
- Check: A/P5-F-11; B/contract satisfiability and prerequisite existence.
- Cited excerpts: F-04 fixes “processamento SÍNCRONO ... **200** ... sem tabela de job, sem polling”; IC-03 still fixes “`202 run id`.”
- Exact defect locus: `M-02-price-intel-core/F-04-collect-verdict-api/feature.md`, `## Expected Output` and `## Inputs/Outputs`; `research/market-evidence-read-interface-contract.md`, `## Operations`.
- Offending token/value: `200 summary/no job` versus `202 run id`.
- Defect: the reconciliation claims the execution model was frozen, but the binding interface contract was not folded. A `202 run id` requires a result/status resource that the brief explicitly excludes. M-06 also consumes the new synchronous `200` behavior.
- Yes-if: IC-03, M-02 F-04, and M-06 agree on one frozen synchronous response/status model, with every required resource named at planning time.

### P5-R2-05

- Severity: `BLOCKING`
- Check: A/P5-F-12; C/consistency with removed M-04 mutation-envelope decision.
- Cited excerpts: the mission says “Toda mutação entra na fila M-03,” “fila M-03 preview+protocolo p/ vínculos em massa,” “Único write path de mutação = envelope M-03,” and ADR-08 says “toda mutação via fila M-03.” Its M-04 ownership row instead says “batch preview/apply local DENTRO de product_links ... mutations fora do MVP p/ vínculos.”
- Exact defect locus: `mission.md`, `## Outcome`, `## Domain Scope`, `### Runtime Contract`, ADR-08, and `### Ownership matrix`.
- Offending decision: universal M-03 mutation routing versus the reconciled local M-04 batch mechanism.
- Defect: the M-04 milestone and briefs now describe a satisfiable local mechanism, but the mission-level architecture still requires bulk linkage through the envelope. The fold therefore leaves two mutually exclusive write paths as binding instructions.
- Yes-if: the mission-level mutation rule and M-04’s reconciled local mechanism state one consistent architecture, explicitly carrying any approved M-04 exception instead of simultaneously requiring envelope routing.

### P5-R2-06

- Severity: `BLOCKING`
- Check: A/P5-F-16; B/exact evidence-field propagation.
- Cited excerpts: IC-03 requires “Toda exibição de preço de mercado” to carry “`source`, `fetched_at`, `n_offers`/`n_sellers` ... `match_status`.” M-07 promises only “`n_offers`/`n_sellers`, fetched_at, freshness.”
- Exact defect locus: `M-07-simulador/F-02-simulador-ui/feature.md`, `## Expected Output`; `research/market-evidence-read-interface-contract.md`, `## Evidence Fields`.
- Offending omissions: `source` and `match_status`.
- Defect: the M-05 displays were folded, but the Simulador market comparison still lacks two mandatory IC-03 fields.
- Yes-if: the Simulador market-price comparison receives and visibly presents `source`, `fetched_at`/freshness, aggregate counts, and `match_status`, using honest unknown states.

### P5-R2-07

- Severity: `BLOCKING`
- Check: A/P5-F-17; B/contract propagation and full decomposition.
- Cited excerpts: F-01 declares `Decompose(input) → {preco, comissao, taxa_fixa, frete, imposto, difal, custo, margem_valor, margem_pct, ...}` while separately saying full “usa `tarifa_full`.” IC-04’s unique formula likewise omits `tarifa_full` despite CalcProfile defining it.
- Exact defect locus: `M-07-simulador/F-01-pricing-calc-difal-service/feature.md`, `## Expected Output`; `research/pricing-difal-interface-contract.md`, `## Resources Or Entities`.
- Offending omission: `tarifa_full` from the decomposition output and unique formula.
- Defect: the reconciliation claims this component was added, but neither the port output nor the binding formula accounts for it. Full-mode results therefore cannot expose or balance the complete approved breakdown.
- Yes-if: IC-04 and F-01 include `tarifa_full` as an explicit nullable decomposition component and debit it in the unique formula for `full`, with unknown propagation when absent.

### P5-R2-08

- Severity: `BLOCKING`
- Check: B/contract propagation of exact endpoint paths.
- Cited excerpts: F-01 plans “`GET/PUT /pricing/difal`”; IC-04 contracts “`PUT /pricing/difal/{uf}`.”
- Exact defect locus: `M-07-simulador/F-01-pricing-calc-difal-service/feature.md`, `## Expected Output`; `research/pricing-difal-interface-contract.md`, `## Operations`.
- Offending token/value: `PUT /pricing/difal` instead of `PUT /pricing/difal/{uf}`.
- Defect: the producing brief does not propagate the contracted per-UF override route.
- Yes-if: the feature brief, OpenAPI/SDK plan, and UI consumer use IC-04’s exact `PUT /pricing/difal/{uf}` operation.

### P5-R2-09

- Severity: `BLOCKING`
- Check: B/DAG edges and prerequisite existence.
- Cited excerpts: M-09 says “Wave B mergeada (`M-04/M-05/M-07/M-08` fornecem os dados agregáveis).” The mission DAG says M-09’s real producers are `{M-01, M-03, M-04, M-05, M-08}` and explicitly states “M-07 NÃO produz artefato consumido pelo M-09.”
- Exact defect locus: `M-09-dashboard-demo/milestone.md`, `## Dependencies`; `mission.md`, `## Parallel Execution Plan`.
- Offending token/value: M-07 as an M-09 data producer/prerequisite.
- Defect: the fold correctly removed false aggregate edges in the mission DAG, but the M-09 milestone reintroduces a false M-07 dependency.
- Yes-if: M-09’s dependency statement names only the actual producing milestones already identified by the mission DAG and its own feature inputs.

### P5-R2-10

- Severity: `ADVISORY`
- Check: B/no new product scope; C/contract consistency.
- Cited excerpt: R-04 says product-category overrides “are expected to be entered per-SKU by the operator.”
- Exact defect locus: `research/difal-interna-rates-2026.md`, `## Caveats`, item 4.
- Offending scope: operator-entered per-SKU tax-rate overrides.
- Defect: IC-04 and M-07 define only a tenant profile plus per-UF DIFAL overrides; no per-SKU fiscal override surface is approved or assigned.
- Yes-if: R-04 does not assert an MIS-004 per-SKU override capability absent from IC-04, and treats product-specific taxation as outside this seed’s contracted behavior.

## Non-findings worth noting

- The manifest contains 45 entries. All 45 file hashes match, and the normalized entry aggregate is exactly `7ed7b63799ebe9185d40c7aeaf03156bbdfb3e6c0e3362eaa285c1ee4795a508`.
- P5-F-01’s intra-wave gate, P5-F-02’s two missing edges, and P5-F-03’s removal of the false M-06→M-09 edge are present in `mission.md`.
- All active migration ownership paths now use `apps/server_core/migrations/`; the numeric blocks remain disjoint and the shared `runner_test.go` seam is named.
- The ERP import lifecycle, unknown reserved-stock behavior, hard-negative `REJECT`, frozen M-05 internal port, Configurações→DIFAL entry, M-03 F-03 interaction model, R-04 dataset reference, and disabled Cancelados stub are folded.
- Removed tokens such as `db/migrations`, `NO_SOLUTION`, and `reservado_desconhecido` remain only in historical audit/reconciliation/context artifacts, not active briefs. The active `202` in IC-03 is covered by P5-R2-04.
- Required UI state/interaction sections are present in the changed UI briefs.
- M-04’s local batch design is internally satisfiable against the inspected mutation-envelope limitation; the remaining defect is the contradictory mission-level architecture, not the local batch brief itself.

## Verdict

**NEEDS-FOLD** — 12 of 19 r01 findings are closed; P5-F-05, P5-F-06, P5-F-10, P5-F-11, P5-F-12, P5-F-16, and P5-F-17 remain open, with additional blocking regressions P5-R2-08 and P5-R2-09.