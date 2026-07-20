# CHIP-BASELINE-FIX — integration lane compile repair

**Branch:** `chip/baseline-fix` (off `main` @078caf12) · **Scope:** EXACTLY 1 file · **Gate:** compile-green

## Problem
`dashboardapp.NewService` takes 7 params since dashboard S1 (@95a45a12):
`installation, listings, linkage, orders, sync ports.SyncSource, erp ports.ErpImportSource, now func() time.Time`.
Fixture `apps/server_core/tests/integration/aggregate_sync_read_test.go` (~line 255) still
called it with 6 args — passed `operationSvc` (sync) but was MISSING the
`erp ports.ErpImportSource` arg before `func() time.Time`. Result:
`not enough arguments in call to dashboardapp.NewService` → `tests/integration` package fails
to compile → `npm run harness:integration` exit-1s for every milestone regardless of chip
quality (pre-existing baseline-red on main).

## Fix (mirror production call root.go:614)
Inserted the missing `erp` argument in the 7-param order, built the analogous test instance
the same way root.go constructs `erpQuerySvc` (root.go:283–285:
`erpapp.NewQueryService(erppostgres.NewRepository(pool), tenantID)`), tenant-scoped to
`h.tenantA` to match the sibling summary services in the same call
(productlinks/orders both use `pool, h.tenantA`).

### Diff (1 file, 3 insertions)
```diff
@@ import block @@
 	dashboardtransport "marketplace-central/apps/server_core/internal/modules/dashboard/transport"
+	erppostgres "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres"
+	erpapp "marketplace-central/apps/server_core/internal/modules/erp_import/application"
 	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
@@ newAggregateSyncHarness NewService call @@
 		ordersapp.NewSummaryService(orderspostgres.NewOrderRepository(pool, h.tenantA)),
 		operationSvc,
+		erpapp.NewQueryService(erppostgres.NewRepository(pool), h.tenantA),
 		func() time.Time { return h.now },
 	)
```

## Compile-green proof (THE GATE)
```
cd apps/server_core
GOMODCACHE=$(pwd)/.gomodcache GOCACHE=$(pwd)/.gocache go build -tags=integration ./tests/integration/...
EXIT=0
```
Modcache warmed first per profile §3 bootstrap. Package COMPILES.

## gofmt
`gofmt -l` lists the file, but `gofmt -d` shows a whole-file line-ending-only diff (CRLF→LF,
zero token changes; the two new import lines are NOT reordered by gofmt → ordering correct).
This is the profile §2/§3 documented CRLF false-alarm signature on Windows worktrees
(autocrlf checkout artifact), NOT a formatting defect.

## Cold review (single pass — sufficient for 1-file compile-fix of a test fixture)
- Arg order matches the 7-param signature exactly (mirrors root.go:614). ✓
- `*erpapp.QueryService` satisfies `dashboard/ports.ErpImportSource`
  (`ListImports(context.Context) ([]erpimportdomain.ImportReport, error)`). ✓
- Repo ctor `erppostgres.NewRepository(pool)` — single arg, identical to root.go:283. ✓
- Tenant `h.tenantA` consistent with sibling summary services (test isolation contract). ✓
- NON-SCOPE respected: NewService signature UNCHANGED; dashboard / sync / root.go / other
  test files UNTOUCHED. Scope = 1 file, 3 insertions (`git diff --stat` confirms). ✓
- No M-01 / M-02 file overlap. ✓

## Hub re-run note
Running the test needs Postgres → HUB re-runs the full `harness:integration` lane at
acceptance. Compile-green is the chip-side gate; not booting a stack (chip rule).

**Status: READY-TO-MERGE.**
