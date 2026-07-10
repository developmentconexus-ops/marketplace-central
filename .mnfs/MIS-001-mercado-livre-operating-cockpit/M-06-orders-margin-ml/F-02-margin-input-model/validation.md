# F-02 Validation

```yaml
id: M-06-F-02
type: feature-validation
status: quick_validation_passed
owner: Feature Implementer
parent: M-06
created: 2026-07-09
updated: 2026-07-09
```

## Scope

Validate margin input modeling and manual adjustment audit only.

## Local Contract Validation

- Command:
  - `cd apps/server_core; $env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./internal/modules/profitability/... ./internal/composition -count=1`
- Result:
  - Passed.
  - profitability application and transport tests passed.
- Command:
  - `npm run test --workspace @marketplace-central/sdk-runtime`
- Result:
  - Passed with `33` tests.
  - SDK now covers profitability import/list/manual-adjustment paths.

## Runtime Validation

- Environment:
  - Docker backend on `http://localhost:8080`
  - restart command: `docker compose restart backend`
  - backend boot evidence: `applied 1 migration(s)` for `0028_profitability_margin_inputs.sql`
- Import evidence:
  - `POST /profitability/margin-inputs/import`
  - body: `{"installation_id":"inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98","limit":1}`
  - result: `imported_count=9`
- Persisted input evidence:
  - `GET /profitability/margin-inputs?installation_id=inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98&limit=20`
  - confirmed:
    - `revenue=19.9` with `quality=complete`
    - `sale_fee=9.63` with `quality=complete`
    - `cost` and all tax components with `quality=missing`
    - order-level `freight` and `commission` with `quality=missing`
- Manual adjustment evidence:
  - `POST /profitability/manual-adjustments`
  - body includes order scope, category `freight`, amount `12.5`, operator metadata, and note `manual freight validation`
  - persisted readback via `GET /profitability/manual-adjustments?...`

## Live Provider Boundary

- Live validation covers:
  - profitability input import from a real Mercado Livre order snapshot already persisted in F-01
  - persistence of profitability input rows in local Postgres
  - append-only manual adjustment audit with actor and reason
- Live validation does not cover:
  - final margin calculation
  - profitability UI
  - a fully resolved-link order with live Oracle-enriched cost/tax in this session

## Open Blockers

- None for F-02 base model.

## Postgres Evidence

- Command:
  - `select count(*) as input_count from profitability_margin_inputs where installation_id = 'inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98';`
  - `select count(*) as adjustment_count from profitability_manual_adjustments where installation_id = 'inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98';`
- Result:
  - `input_count=9`
  - `adjustment_count=1`

## Quality Honesty Evidence

- The imported order item had no resolved product link.
- Cost and tax inputs were still assembled as explicit rows with `quality=missing`.
- No missing input was converted to `0`.
