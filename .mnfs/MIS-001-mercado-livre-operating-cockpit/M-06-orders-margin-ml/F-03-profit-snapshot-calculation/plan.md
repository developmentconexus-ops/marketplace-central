# F-03 Plan - Profit Snapshot Calculation

```yaml
id: M-06-F-03
type: feature-plan
status: drafted
owner: Feature Implementer
parent: M-06
created: 2026-07-09
updated: 2026-07-09
```

## Slice Decision

Build server-side profitability calculation on top of the persisted F-02 input model, with no UI dependency and no hidden estimations.

## Work Breakdown

1. Add `ProfitSnapshot` domain types and calculation quality flags.
2. Extend profitability persistence with profit snapshots.
3. Implement calculation service over margin inputs + manual adjustments.
4. Expose calculate/list endpoints.
5. Update OpenAPI + SDK together.
6. Run local tests and real runtime recalculation against the same Mercado Livre order used in F-01/F-02.

## Verification Commands

Backend:

```powershell
cd apps/server_core
$env:GOCACHE="$PWD\\.gocache"
go test ./internal/modules/profitability/... ./internal/composition -count=1
```

SDK:

```powershell
npm run test --workspace @marketplace-central/sdk-runtime
```

Runtime:

```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/profitability/profit-snapshots/calculate" -ContentType "application/json" -Body '{"installation_id":"inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98","limit":1}'
```
