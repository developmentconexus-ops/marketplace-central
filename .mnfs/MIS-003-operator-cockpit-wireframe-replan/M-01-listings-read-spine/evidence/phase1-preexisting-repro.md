# Evidence — TestPhase1SmokeFlow is a PRE-EXISTING failure, independent of M-01 changes

Per hub ruling B (2026-07-15): M-01 acceptance does NOT require full integration-lane green;
it requires proof that `tests/integration/TestPhase1SmokeFlow` fails WITHOUT M-01's changes.
Hub owns the fix (post-M-01 corrective chip). Chip must not fix it.

## Environment
Fresh worktree `.claude/worktrees/m01-listings`, hermetic per HARNESS §5:
`cd apps/server_core && GOMODCACHE=$(pwd)/.gomodcache go mod download all` (cache warmed, ~129M).
Ephemeral postgres:16-alpine, DB name `mpc_test_0123456789abcdef0123456789abcdef` (LoadConfig
`^mpc_test_[0-9a-f]{32}$` rule), `go run ./cmd/testdb migrate` → `applied 36 migration(s)`.
Env: `GOCACHE=.gocache GOMODCACHE=.gomodcache GOPROXY=off GOSUMDB=off MPC_TEST_DATABASE_URL=<ephemeral>`.

## Repro 1 — ORIGINAL 4-package set (NO listings), clean migrated DB
Command:
```
go test -tags=integration ./tests/integration ./internal/modules/orders/adapters/postgres \
  ./internal/modules/profitability/adapters/postgres ./internal/modules/product_links/application -count=1
```
Output:
```
--- FAIL: TestPhase1SmokeFlow (0.44s)
    phase1_smoke_test.go:75: run simulation error: PRICING_INVALID_PRODUCT_ID
FAIL
FAIL	marketplace-central/apps/server_core/tests/integration	39.298s
ok  	marketplace-central/apps/server_core/internal/modules/orders/adapters/postgres	39.994s
ok  	marketplace-central/apps/server_core/internal/modules/profitability/adapters/postgres	8.167s
ok  	marketplace-central/apps/server_core/internal/modules/product_links/application	7.584s
FAIL
```
Listings NOT in this set → failure reproduces with zero M-01 changes.

## Repro 2 — isolation, single test, clean migrated DB
Command:
```
go test -tags=integration ./tests/integration -run TestPhase1SmokeFlow -count=1
```
Output:
```
--- FAIL: TestPhase1SmokeFlow (0.34s)
    phase1_smoke_test.go:75: run simulation error: PRICING_INVALID_PRODUCT_ID
FAIL	marketplace-central/apps/server_core/tests/integration	3.389s
```
Deterministic, fast, isolated → not cross-test contamination.

## Non-linkage to D-24
`grep -n "internal_read|ICMSCeiling|GetICMSCeilingByOrigin" tests/integration/phase1_smoke_test.go`
→ no matches. The pricing smoke path does not touch internal_read; D-24's additive ceiling
method is not implicated.

## Contrast — M-01-owned packages GREEN in the same hermetic lane
```
ok  marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres  7.802s
ok  marketplace-central/apps/server_core/internal/modules/internal_read/... (build+test)
```
