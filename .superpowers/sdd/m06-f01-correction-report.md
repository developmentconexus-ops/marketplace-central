# M-06 F-01 Correction Report

## Status

`DONE`

## Root Cause

The repository freshness decision was not atomic with the conflicting write, so a stale snapshot could win between a read/check and the upsert and then replace its children. Freshness is now enforced by the `ON CONFLICT ... DO UPDATE ... WHERE` statement itself; item and payment replacement runs only when that statement reports a winning row. Provider payload minimization is enforced at the application boundary by deriving only the safe Mercado Livre order path and ignoring the connector's raw payload.

## RED Evidence

- Command:
  - `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./internal/modules/orders/... -run TestImportServiceDiscardsBuyerPIIFromRawProviderReference -count=1`
- Expected failure excerpt:
  - `import_service_test.go:155: raw provider reference = "/orders/2003", want safe reference only`
  - `FAIL marketplace-central/apps/server_core/internal/modules/orders/application`

## GREEN Evidence

- Command:
  - `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; $env:MC_DATABASE_URL='postgres://mpc:mpc-test@localhost:5436/marketplace_central?sslmode=disable'; go test -json -count=1 -timeout=5m ./internal/modules/orders/...`
- Result:
  - Exit code `0` on 2026-07-09.
  - All seven packages under `internal/modules/orders/...` passed or reported `[no test files]`.
  - Both PostgreSQL repository tests passed, covering absent, newer, equal, older, known/unknown timestamps, child replacement, row counts, and safe provider-reference readback.
  - Application tests passed for paid normalization, cancelled-order metadata, missing product link, payments, and buyer-PII minimization.

## Changed Paths

- `apps/server_core/internal/modules/orders/application/import_service_test.go`
- `apps/server_core/internal/modules/orders/application/import_service.go`
- `apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go`
- `apps/server_core/internal/modules/orders/adapters/postgres/order_repo_test.go`
- `.superpowers/sdd/m06-f01-correction-report.md`

## Real-DB Verification

Executed against a disposable PostgreSQL 16 instance with the real `0027_orders_marketplace_orders.sql` migration. The repository tests perform real inserts/upserts and repository readback. They verify stable order/item/payment row counts and compare the persisted `raw_provider_ref` with the safe `/orders/{id}` reference. This is real database integration evidence; it is not a mock and does not claim live Mercado Livre validation.

## Self-Review Concerns

- The shared worktree contains unrelated changes; they were not reverted or included in this correction.
- The RED excerpt retained above documents the PII expectation correction. The earlier atomic-upsert RED output was not preserved by the worker, so this report does not fabricate it; the real PostgreSQL regression test and GREEN output are the available durable proof.
- No commit was created.

## Boundary Fix Cycle

### RED

- Command:
  - `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./internal/modules/orders/application ./internal/modules/orders/adapters/integrations -count=1`
- Result:
  - Exit code `1` on 2026-07-09.
  - Application and integration-adapter tests failed to compile because `orders/domain.OrderIngestionSnapshot`, `OrderIngestionItem`, and `OrderIngestionPayment` did not yet exist.
  - This is the expected RED state: application tests require an orders-owned ingestion contract and adapter tests require translation from connector-owned snapshots.

### GREEN

- Command:
  - `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./internal/modules/orders/application ./internal/modules/orders/adapters/integrations -count=1`
- Result:
  - Exit code `0` on 2026-07-09.
  - `internal/modules/orders/application` passed.
  - `internal/modules/orders/adapters/integrations` passed.
  - The adapter test verifies every non-raw connector snapshot field is translated and `OrderIngestionSnapshot` does not expose `RawProviderRef`.
