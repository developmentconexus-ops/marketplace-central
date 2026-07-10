# M-06 F-03 cold-gate correction report

Date: 2026-07-10

## Scope and safety

Corrected the three cold-gate findings in profitability without resetting, reverting, staging, or committing the shared dirty checkout. No provider write was issued. PostgreSQL validation used the local Docker service; the application database write was isolated to random test tenant and installation identifiers and cleaned up by the integration test.

## Root cause confirmation

1. `applyInput` accumulated each tax component with `addOptional`; a nil component could leave a non-nil aggregate, so `finalizeSnapshot` only flagged tax missing when the whole total was nil.
2. `applyAdjustment` routed every `cost` category to `CostAmount`, without considering adjustment scope.
3. The create path generated an adjustment ID then performed a bare append insert, with neither a request key nor a tenant/installation uniqueness constraint.

## TDD evidence

Focused tests were added before the corresponding production changes:

- `TestBuildProfitSnapshotsKeepsPartialTaxIncompleteAtItemAndOrder`
- `TestBuildProfitSnapshotsRoutesCostAdjustmentByScope`
- `TestCreateManualAdjustmentRejectsMissingIdempotencyKey`
- `TestCreateManualAdjustmentReturnsOriginalForDuplicateSubmission`

Initial RED command:

```powershell
$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'
go test ./internal/modules/profitability/application -run 'TestBuildProfitSnapshotsKeepsPartialTaxIncompleteAtItemAndOrder|TestBuildProfitSnapshotsRoutesCostAdjustmentByScope|TestCreateManualAdjustmentRejectsMissingIdempotencyKey|TestCreateManualAdjustmentReturnsOriginalForDuplicateSubmission' -count=1 -v
```

Initial RED result: compile failure at `service_test.go:588`: `CreateManualAdjustmentInput` had no `IdempotencyKey`. This was the expected absent contract. The batch did not reach the runtime assertions until the domain/application contract was added; this is recorded explicitly rather than presenting a compile failure as runtime evidence.

GREEN command and result:

```powershell
$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'
go test ./internal/modules/profitability/application -count=1 -v
```

Result: PASS, including all four focused cases. The full profitability suite was then green.

## Implementation

### Partial tax honesty

- A nil tax-component amount marks `missing_tax` while preserving known component amounts in `TaxAmount`.
- `finalizeSnapshot` treats the accumulated `missing_tax` flag as required-missing at item and order scope, so contribution and margin remain nil and quality is `incomplete`.
- Item flags are propagated during order aggregation; the order keeps known tax totals but cannot become realized while one item’s tax set is incomplete.

### Cost scope

- Item-scope `cost` adjustment: adds to `CostAmount`.
- Order-scope `cost` adjustment: adds to `AdjustmentAmount`.
- Freight, commission, and generic-adjustment mapping remains unchanged.

### Idempotent adjustments

- Added `idempotency_key` to create input, domain record, HTTP transport, OpenAPI, SDK, and Orders UI client contract.
- Empty/whitespace keys return `PROFITABILITY_ADJUSTMENT_IDEMPOTENCY_KEY_REQUIRED` before any ID generation or write.
- Replaced append-only store API with atomic `CreateOrGetAdjustment`.
- Migration `0032_profitability_manual_adjustment_idempotency.sql` adds the column with a default empty string for historical readability and a partial unique index on `(tenant_id, installation_id, idempotency_key)` for non-empty keys. Existing rows remain readable and are outside the new uniqueness index.
- PostgreSQL uses `INSERT ... ON CONFLICT ... DO NOTHING`, then reads the immutable original record for a duplicate key. Returned timestamp is normalized to UTC, matching the rest of the domain read path.
- Orders UI creates one stable request key for the active submission and retains it through a failed retry; it clears only after the full submit/calculate/refresh sequence succeeds.

## Changed paths

- `apps/server_core/internal/modules/profitability/application/service.go`
- `apps/server_core/internal/modules/profitability/application/service_test.go`
- `apps/server_core/internal/modules/profitability/domain/input.go`
- `apps/server_core/internal/modules/profitability/ports/store.go`
- `apps/server_core/internal/modules/profitability/adapters/postgres/store.go`
- `apps/server_core/internal/modules/profitability/adapters/postgres/manual_adjustment_integration_test.go`
- `apps/server_core/internal/modules/profitability/transport/http_handler.go`
- `apps/server_core/migrations/0032_profitability_manual_adjustment_idempotency.sql`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`
- `packages/sdk-runtime/src/index.test.ts`
- `packages/feature-orders/src/OrdersPage.tsx`
- `packages/feature-orders/src/OrdersPage.test.tsx`

## Validation evidence

| Validation | Result | Target |
| --- | --- | --- |
| `go test ./internal/modules/profitability/... -count=1` | PASS | Host Go; adapters/orders, adapters/postgres, application, transport |
| `go test ./internal/modules/profitability/adapters/postgres -run TestManualAdjustmentsAppendOnlyReadbackAndConstraints -count=1 -v` | PASS | Real Docker PostgreSQL (`127.0.0.1:5435`); validates duplicate key returns original and only two distinct-key rows remain |
| `npm run test --workspace @marketplace-central/sdk-runtime -- --run` | PASS, 35/35 | SDK contract behavior; bounded to 60 seconds |
| `npm run test --workspace @marketplace-central/feature-orders -- --run` | PASS, 12/12 | UI caller and identity flow; bounded to 60 seconds |
| `python -c "import yaml; yaml.safe_load(...)"` | PASS | OpenAPI YAML parse |
| `gofmt -w ...profitability...` | PASS | Changed Go source and tests formatted |
| boundary `rg` | PASS | No profitability imports of orders internals outside `profitability/adapters/orders` |
| Docker backend restart | PASS | `0032` applied; backend healthy; `/healthz` HTTP 200 |

The real database test first exposed a timezone-location mismatch on duplicate readback: the timestamps represented the same instant but Go structural equality differed. The store now normalizes duplicate readback to UTC; the same real PostgreSQL test passed after that correction.

## Concerns / remaining gates

- The focused initial RED command stopped at the expected absent `IdempotencyKey` compile error, so there is no separate pre-implementation runtime failure transcript for partial-tax and cost assertions. Their pre-fix causes were directly confirmed in `service.go`, and the focused tests were present before the production change.
- This report does not declare M-06 passed. Cold independent review, required validation artifacts, and any broader milestone gate remain owned by the orchestrator.

## Follow-up: UI idempotency lifecycle correction

An independent re-review found that a single pending UI key survived material form/order changes after a failed request. That could cause a new intended adjustment to receive the original key and be returned as the earlier immutable audit record.

### TDD evidence

Added `OrdersPage > reuses an idempotency key only for an exact failed-submission retry` before the follow-up implementation.

RED command:

```powershell
npm run test --workspace @marketplace-central/feature-orders -- --run -t "reuses an idempotency key only for an exact failed-submission retry"
```

RED result: the third request, after a failed exact retry and a selected-order change, reused the first key. Vitest assertion: `expected '<key>' not to be '<key>'`.

Implementation: `pendingAdjustmentSubmission` now stores both the key and an immutable JSON fingerprint of installation, order/item/variation, scope, category, amount, currency, resolved reason, and actor. The exact same fingerprint reuses its key; any material context change creates a new key. A successful create/calculate/refresh clears the pending submission.

GREEN and impacted validation:

| Validation | Result |
| --- | --- |
| focused lifecycle regression | PASS |
| feature-orders suite | PASS, 13/13 |
| SDK suite | PASS, 35/35 |

The host and Docker Vite production-build attempts initially remained stuck at
`transforming...` for over the bounded validation interval. Only the `npm/vite
build` processes launched by this validation were terminated; the shared
frontend dev process was preserved. A subsequent isolated verification avoided
the long-running frontend development process and passed:

```powershell
docker compose run --rm --no-deps frontend npm run build --workspace @marketplace-central/web
```

It transformed 1783 modules and completed successfully in 4.31 seconds.

### Follow-up import cleanup

Removed one duplicate identical React-hooks import from `packages/feature-orders/src/OrdersPage.tsx`. Focused lifecycle regression test passed (1/1; remaining tests skipped by name filter).
