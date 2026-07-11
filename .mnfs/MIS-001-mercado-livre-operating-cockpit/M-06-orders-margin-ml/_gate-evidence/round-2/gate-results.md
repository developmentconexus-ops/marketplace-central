# M-06 round-2 execution QA evidence

Evidence type: ran

- Reviewed implementation/evidence freeze: `5548ae406cb26d0703c111236d703281bb227d3e`.
- Tested HEAD: `2e8c9250303c9e9300055dc8030e9ae7fb62093c`.
- Freeze check: `git diff --name-status 5548ae406cb26d0703c111236d703281bb227d3e..HEAD` returned only `M milestone-review.md`; no implementation, contract, OpenAPI, SDK, or runtime-source path changed.
- No approval, provider write, import, input import, recalculation, manual adjustment, or runtime-data mutation was issued by QA.

## Re-run corroboration

| Criterion / seam | Exact command | Observed | Outcome |
| --- | --- | --- | --- |
| M-06-C01 | From `apps/server_core`: `$env:GOCACHE=(Join-Path $PWD ".gocache"); go test ./internal/modules/orders/... ./internal/modules/connectors/... ./internal/modules/integrations/transport ./internal/composition -count=1` | Exit 0; orders, connector, integration-transport, and composition packages passed. | reproduced |
| M-06-C01 targeted | From `apps/server_core`: `$env:GOCACHE=(Join-Path $PWD ".gocache"); go test ./internal/modules/orders/application ./internal/modules/orders/transport -run "TestImportService|TestHandleImport" -count=1 -v` | All selected tests passed. | reproduced |
| M-06-C02/F-03 | From `apps/server_core`: `$env:GOCACHE=(Join-Path $PWD ".gocache"); go test ./internal/modules/profitability/... ./internal/composition -count=1` | Exit 0; all non-integration profitability packages passed. | reproduced |
| M-06-C02 targeted | From `apps/server_core`: `$env:GOCACHE=(Join-Path $PWD ".gocache"); go test ./internal/modules/profitability/application -run "TestImportMarginInputsCompleteAndMissingLink|TestBuildProfitSnapshotsKeepsPartialTaxIncompleteAtItemAndOrder|TestCalculateSnapshotsHandlesCompleteIncompleteAndNegative" -count=1 -v` | All three selected tests passed. | reproduced |
| M-06-C03 targeted | From `apps/server_core`: `$env:GOCACHE=(Join-Path $PWD ".gocache"); go test ./internal/modules/profitability/application ./internal/modules/profitability/transport -run "TestCreateManualAdjustment|TestHandleCreateAdjustment" -count=1 -v` | Selected audit/idempotency/validation tests passed; they reject empty actor fields but do not establish trusted actor provenance. | reproduced with security defect |
| SDK | `npm run test --workspace @marketplace-central/sdk-runtime` | 35/35 passed. | reproduced |
| Orders UI | `npm run test --workspace @marketplace-central/feature-orders -- OrdersPage.test.tsx` | 13/13 passed. | reproduced |
| Web route | `npm run test --workspace @marketplace-central/web -- AppRouter.test.tsx` | 4/4 passed. | reproduced |
| Web context | `npm run test --workspace @marketplace-central/web -- ClientContext.test.tsx` | 3/3 passed. | reproduced |
| Web proxy | `npm run test --workspace @marketplace-central/web -- viteProxy.test.ts` | 1/1 passed. | reproduced |
| Web build | `npm run build --workspace @marketplace-central/web` | Exit 0; 1,783 modules transformed and production bundle emitted. | reproduced |

The npm/Vite commands first hit a sandbox-only workspace-config access denial and then passed unchanged when rerun with the required filesystem access. The integration-tag PostgreSQL adjustment test was not rerun because the QA boundary prohibited runtime-data mutation and `MC_DATABASE_URL` was unset.

## Read-only live API flows

All requests were GET-only against `http://localhost:8080`:

- `/healthz`: HTTP 200, service status `ok`.
- `/orders?...&limit=50`: HTTP 200; 30 orders, 30 unique order IDs, 30 item identities, 30 unique item identities, 30 payment identities, and 30 unique payment identities; 24 paid and 6 cancelled.
- `/profitability/margin-inputs?...&limit=500`: HTTP 200; 270 inputs. Of 82 `missing` inputs, zero carried a numeric value; of 68 `complete` inputs, zero had a null amount. One complete zero is the known PIS zero already recorded in the reconciliation.
- `/profitability/profit-snapshots?...&limit=500`: HTTP 200; 60 snapshots (30 item, 30 order), 48 realized/incomplete and 12 not-realized. Zero `missing_tax` snapshots carried contribution or margin math.
- `/profitability/manual-adjustments?...&limit=500`: HTTP 200; zero rows. No adjustment was created for QA.

## Blocking seam adjudication

### Quantity greater than one

The live Candidate A readback reproduces the cold-review defect:

| Quantity | Revenue | Persisted CUSSEMICM cost | Result |
| ---: | ---: | ---: | --- |
| 1 | 169.99 | 91.57 | realized, incomplete |
| 2 | 339.98 | 91.57 | realized, incomplete |
| 7 | 1189.93 | 91.57 | realized, incomplete |

Revenue is extended at `apps/server_core/internal/modules/profitability/application/service.go:429`, while cost is read and passed without quantity at `:229-230`, persisted unchanged at `:508`, and subtracted once at `:818-820`. The existing quantity-two test at `service_test.go:92-135` asserts extended revenue but not cost. Because `F-02-margin-input-model/spec.md:51-57` does not define unit-versus-line CUSSEMICM semantics, the fixed SHA cannot establish correct quantity composition. This is a blocking contract and integration defect; current incomplete snapshots avoid emitting incorrect contribution math but do not cure the cost seam.

### Manual-adjustment actor boundary

`apps/server_core/internal/modules/profitability/transport/http_handler.go:97-116` accepts and forwards the request-body actor. `application/service.go:265-280` only trims and rejects empty type/ID, then persists that actor at `:309-324`. The route is registered at `internal/composition/root.go:345`, and the composed mux is wrapped only by CORS at `:404`. OpenAPI requires the caller actor at `contracts/api/marketplace-central.openapi.yaml:2779-2806`, and the SDK sends it at `packages/sdk-runtime/src/index.ts:1005-1017`.

The targeted tests corroborate audit-field validation but do not authenticate or authorize the claimed actor. A live acceptance probe was intentionally not sent because it would create a forbidden adjustment. The caller-forgeable actor remains a blocking Security defect.

## Fold

The persisted cold review is never downgraded: ★1 Pass, ★2 Fail, ★3 Fail, ★4 Fail, ★5 Pass, ★6 Fail, ★7 Fail; `2/7`, therefore **Fail**. The unavailable visual drive adds a required-evidence blocker but cannot upgrade or replace the existing Fail verdict.
