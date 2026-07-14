# F-01 validation

Status: `quick_validation_passed`

All Go commands were run from `apps/server_core`. On Windows, Go rejects a
relative `GOCACHE=.gocache` with `GOCACHE is not an absolute path`; the same
registered commands were run with the absolute path
`C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache`.

## Validation expectations

| Expectation | Result | Evidence |
| --- | --- | --- |
| Build all packages | Pass | `go build ./...` exited 0. |
| Internal-read tests | Pass | `go test ./internal/modules/internal_read/... -v` exited 0; fake, Oracle adapter, application, domain, and ports packages passed. |
| Full Go test suite | Pass | `go test ./...` exited 0; all listed packages passed, including transport and unit suites. |
| Default Oracle pool max is 12 | Pass | `TestLoadConfigFromEnvUsesBoundedDefaults` passed with no pool env overrides. |
| `MPC_ORACLE_POOL_MAX_SESSIONS=0` is rejected | Pass | `TestLoadConfigFromEnvRejectsUnsafePoolAndTimeoutValues/MPC_ORACLE_POOL_MAX_SESSIONS` passed. |
| HTTP server timeout discipline | Pass | `cmd/server/main.go` constructs `http.Server` with `ReadHeaderTimeout=10s`, `ReadTimeout=30s`, and `WriteTimeout=60s`. |
| Route-class deadlines and cancellation | Pass | Focused `httpx` httptest passed: production constants are interactive 15s and batch 120s; scaled 15ms interactive stall returned exact 504 `{"error":"deadline_exceeded"}` and canceled the handler context; the same 30ms stall completed under the scaled 120ms batch budget with 200. |
| Batch routes declared at registration | Pass | Composition registers exactly the IC-01 concrete batch patterns before transport registration: profitability margin import, profitability snapshot calculation, orders import, both product-links import/generation routes, pricing batch simulations, and fee-schedule sync. |

## Command transcripts

```text
GOCACHE=.gocache go test ./internal/platform/httpx ./internal/modules/internal_read/... -v
build cache is required, but could not be located: GOCACHE is not an absolute path

$env:GOCACHE=(Join-Path (Get-Location).Path '.gocache'); go test ./internal/platform/httpx ./internal/modules/internal_read/... -v
PASS
ok  marketplace-central/apps/server_core/internal/platform/httpx
ok  marketplace-central/apps/server_core/internal/modules/internal_read/adapters/fake
ok  marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle
ok  marketplace-central/apps/server_core/internal/modules/internal_read/application
ok  marketplace-central/apps/server_core/internal/modules/internal_read/domain
ok  marketplace-central/apps/server_core/internal/modules/internal_read/ports

$env:GOCACHE=(Join-Path (Get-Location).Path '.gocache'); go build ./...
exit code: 0

$env:GOCACHE=(Join-Path (Get-Location).Path '.gocache'); go test ./internal/modules/internal_read/... -v
PASS
ok  marketplace-central/apps/server_core/internal/modules/internal_read/adapters/fake
ok  marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle
ok  marketplace-central/apps/server_core/internal/modules/internal_read/application
ok  marketplace-central/apps/server_core/internal/modules/internal_read/domain
ok  marketplace-central/apps/server_core/internal/modules/internal_read/ports

$env:GOCACHE=(Join-Path (Get-Location).Path '.gocache'); go test ./...
PASS
ok  .../internal/modules/marketplaces/transport
ok  .../internal/modules/orders/transport
ok  .../internal/modules/pricing/transport
ok  .../internal/modules/product_links/transport
ok  .../internal/modules/profitability/transport
ok  .../internal/platform/httpx
ok  .../tests/unit
exit code: 0

$env:GOCACHE=(Join-Path (Get-Location).Path '.gocache'); go test ./internal/platform/httpx ./internal/modules/internal_read/... -v
--- PASS: TestRouteClassDeadlinesAndCancellation/interactive_expires (0.02s)
--- PASS: TestRouteClassDeadlinesAndCancellation/batch_keeps_its_longer_budget (0.03s)
--- PASS: TestRouteClassMuxDeclaresBatchBeforeTransportRegistration
--- PASS: TestLoadConfigFromEnvUsesBoundedDefaults
--- PASS: TestLoadConfigFromEnvRejectsUnsafePoolAndTimeoutValues/MPC_ORACLE_POOL_MAX_SESSIONS
exit code: 0
```

## Changed paths

- `apps/server_core/cmd/server/main.go`
- `apps/server_core/go.mod`
- `apps/server_core/go.sum`
- `apps/server_core/internal/composition/root.go`
- `apps/server_core/internal/platform/httpx/route_deadline.go`
- `apps/server_core/internal/platform/httpx/route_deadline_test.go`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/{config.go,config_test.go,open_cgo.go,open_nocgo.go,reader.go,reader_live_test.go,reader_test.go,database.go}`
- transport registration signatures under `apps/server_core/internal/modules/{catalog,classifications,connectors,integrations,inventory,marketplaces,orders,pricing,product_links,profitability}/transport/`

## Handoff

The Oracle adapter hardening refactor was preserved and completed with pool
default 12, HTTP server timeouts, composition-declared route classes, context
deadline cancellation, and redacted Oracle causes. No retry, shutdown, readyz,
or query-shape changes were added. Next owner: Milestone Orchestrator for
fixed-SHA review and proportional QA.

## Scoped correction evidence (C02)

- `apps/server_core/internal/composition/root.go` classifies both
  `/admin/fee-schedules/sync` and `/admin/fee-schedules/seed` as batch before
  transport registration.
- `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test -count=1 ./internal/composition -run TestFeeScheduleRoutesUseBatchDeadline -v` — Pass; both routes resolved with approximately the 120s batch deadline.
- `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go build ./...` — Pass.
- `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./...` — Pass.
