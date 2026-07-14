# F-02-cost-tax-batch

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-2
lifecycle_scope: feature
```

## Mission

MIS-002. IC-01 batch rules.

## Milestone

M-03. Depends on F-01 (shared `oraclebatch` helper).

## Brief

Batch ports `GetCostFactsByIDs` and `GetTaxFactsByIDs` (same map-shape as F-01, chunk 500, semaphore) sourced from LEANDRO.VW_PRECO_TABELA / VW_IMPOSTO_ITEM mappings already in the adapter; convert `profitability/application/service.go:239-287` to batch reads; cap the sales-item query at 5000 rows with explicit `truncated` marker (IC-01); enforce `ImportMarginInputs.Limit` ceiling 200 → 422 `limit_exceeded` naming the cap (IC-01 error matrix, ADR-05).

## Inputs

- N+1 site: `apps/server_core/internal/modules/profitability/application/service.go:239-287`.
- Existing cost/tax SQL + column mappings in `internal_read/adapters/oracle/reader.go`.
- F-01 helper + port conventions.
- IC-01: sales cap 5000, batch route class (profitability report route = batch, 120s).

## Expected Output

- Two batch ports + adapters; profitability service uses ≤ ceil(N/500) queries per fact type per run.
- Sales query `FETCH FIRST 5000 ROWS ONLY` + peek → result carries `Truncated bool`; service surfaces it (report marked partial, never silently truncated).
- One intentional commit.

## Constraints

- Unknown cost/tax → nil + `missing_cost`/`missing_tax` flag; NEVER zero — a zero cost silently corrupts margin math (mission root rule).
- No caching here (M-04).
- SQL in adapter only; `GOCACHE=.gocache`.

## Inputs/Outputs

Map-keyed results like F-01. `Truncated` propagates to the profitability report struct and (existing) response as boolean field documented in OpenAPI IF the report route already exposes a shape change — if so, OpenAPI + sdk-runtime same commit.

## Negative Scenarios

- While `ImportMarginInputs.Limit=201`, when the import request runs, the system shall return 422 `{"error":"limit_exceeded"}` naming the 200 cap, with zero Oracle calls.
- While a product has price but no tax rows, when facts load, the system shall return nil tax + `missing_tax` and profitability marks that product incomputable rather than assuming 0 tax.
- While sales rows exceed 5000, when the report builds, the system shall set `truncated=true` and log a `slow_query`-style structured warning.

## Validation Expectations

- Fake-queryer counts (per fact type, N=1200 → 3 chunks).
- Truncation test at 5001 fixture rows.
- Nil-fact → incomputable-product test.
- `GOCACHE=.gocache go test ./...` green; sdk build green if OpenAPI touched.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` at execution time.

## Handoff

- Current status: briefed
- Next owner: Feature Implementer (mpc-implementer, gpt-5.6-luna high)
- Next action: spec → plan → implement → evidence
- Required files/evidence: `validation.md`
- Blockers or open decisions: None
