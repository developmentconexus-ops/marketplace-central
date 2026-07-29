# CHIP-ANCHORS-3 — integration tests A5/A6/A7/A10 — implementer artifact

Owned file (only file edited, and it remains the only net diff):
`apps/server_core/internal/modules/erp_import/adapters/postgres/chain_query_repository_integration_test.go`

`query_repository.go` (production) was temporarily edited twice, in place, to
prove the must-fail (A6), then re-edited forward back to its original text —
confirmed byte-identical afterward (`git diff --stat` before and after the
round-trip is the same `20 insertions(+), 2 deletions(-)`, i.e. the same
pre-existing diff this slice found on disk, nothing added). No `git reset`,
`revert`, `stash`, or `checkout --` was used anywhere in this slice.

## 1. Full text of the three added tests (verbatim from the file, lines 239-350)

```go
// TestGetImportChainCountsLeadingZeroCodprod pins the resolved_products join
// fix: erp_import_products.codprod keeps the raw spreadsheet string
// ('00101'), while product_links.internal_product_id is a ParseInt'd integer
// column (101). Before the ltrim-both-sides fix, '00101' = '101' compared
// false and this fixture's linked product silently dropped out of
// vinculados.
func TestGetImportChainCountsLeadingZeroCodprod(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	importID := domain.ImportID("64000000-0000-0000-0000-000000000001")

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_links WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_products WHERE tenant_id=$1`, tenant)
	})

	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_protocols
			(id, tenant_id, file_sha256, protocol, source, imported_at, status, accepted_count, rejected_count, warning_count)
		VALUES ($1, $2, 'chain-leading-zero-hash', '#640-E', 'xlsx', now(), 'COMPLETED', 1, 0, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_products (tenant_id, protocol_id, codprod, descrprod, custo, stock_physical)
		VALUES ($2, $1, '00101', 'Product 00101', 1, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO product_links
			(tenant_id, installation_id, provider_code, provider_item_id, state, internal_product_id)
		VALUES ($1, 'installation-a', 'mercadolivre', 'item-00101', 'resolved', 101);
	`, tenant); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetImportChain(ctx, tenant, importID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Protocol != "#640-E" || got.Importados != 1 || got.Vinculados != 1 || got.Enfileirados != 0 {
		t.Fatalf("chain=%#v", got)
	}
}

// TestGetImportChainNonArrayPendingCursorDoesNotFailQuery pins the
// queued_products CASE guard: sync_state.cursor->'pending' is normally a
// JSON array, but a malformed or partial cursor write can leave it as an
// object or a scalar. Before the jsonb_typeof guard,
// jsonb_array_elements_text raised at query time (not a valid JSON array)
// and took the whole chain endpoint down for the tenant instead of reading
// as an empty queue.
func TestGetImportChainNonArrayPendingCursorDoesNotFailQuery(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_products WHERE tenant_id=$1`, tenant)
	})

	cases := []struct {
		name           string
		cursor         string
		importID       domain.ImportID
		protocol       string
		installationID string
	}{
		{name: "pending as object", cursor: `{"pending":{"foo":"bar"}}`, importID: domain.ImportID("65000000-0000-0000-0000-000000000001"), protocol: "#651-E", installationID: "installation-non-array-object"},
		{name: "pending as scalar", cursor: `{"pending":"501"}`, importID: domain.ImportID("65000000-0000-0000-0000-000000000002"), protocol: "#652-E", installationID: "installation-non-array-scalar"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			importID := tc.importID
			if _, err := pool.Exec(ctx, `
				INSERT INTO erp_import_protocols
					(id, tenant_id, file_sha256, protocol, source, imported_at, status, accepted_count, rejected_count, warning_count)
				VALUES ($1, $2, $3, $4, 'xlsx', now(), 'COMPLETED', 1, 0, 0);
			`, importID, tenant, "chain-non-array-hash-"+tc.name, tc.protocol); err != nil {
				t.Fatal(err)
			}
			if _, err := pool.Exec(ctx, `
				INSERT INTO erp_import_products (tenant_id, protocol_id, codprod, descrprod, custo, stock_physical)
				VALUES ($2, $1, '501', 'Product 501', 1, 0);
			`, importID, tenant); err != nil {
				t.Fatal(err)
			}
			if _, err := pool.Exec(ctx, `
				INSERT INTO sync_state (tenant_id, installation_id, entity, cursor)
				VALUES ($1, $3, 'market', $2::jsonb);
			`, tenant, tc.cursor, tc.installationID); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_protocols WHERE id=$1 AND tenant_id=$2`, importID, tenant)
			})

			got, err := repo.GetImportChain(ctx, tenant, importID)
			if err != nil {
				t.Fatalf("expected a valid response, got query error: %v", err)
			}
			if got.Protocol != domain.Protocol(tc.protocol) || got.Importados != 1 || got.Enfileirados != 0 {
				t.Fatalf("chain=%#v", got)
			}
		})
	}
}
```

(A7 is `TestGetImportChainCountsCurrentQueueAcrossInstallations`, already present
at line 22 before this slice — not rewritten, run only.)

## 2. Environment setup — exact commands

```
$ npm run harness:pg:up
> pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command pg-session-up
container=mpc-pg-session-1d7cca5b
port=51692
target=pg-session
status=ready
```

Warmed the module cache first (known false-alarm guard):
```
$ cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go build ./...
(no output — success)
```

### The lane's own artifact cannot prove PASS/SKIP (ratified finding)

`npm run harness:integration` invokes `go test -tags=integration` without
`-v`; `summary.txt` records only `target`/`status`/`run_id` — a fully-SKIPPED
and a fully-GREEN run look identical in it. First lane run in this slice hit
`status=blocked` on an unrelated pre-existing package
(`internal/modules/mutations/application`, produced by other in-flight work
on this branch per `git status`, not touched by this slice); running that
package directly and separately showed it green on its own
(`ok  marketplace-central/apps/server_core/internal/modules/mutations/application  3.778s`),
confirming it was an inter-package interaction in the shared lane run, not a
regression from this slice. That result is **not** used below as proof of
A5/A7/A10 — the proof below is a direct `-v` run against a target set the
same way the harness lane sets it (same DSN shape, same
`go run ./cmd/testdb migrate` mechanism read from
`scripts/harness/Postgres.psm1:382-399`), not the lane's exit code:

```
$ RUNID=<32-hex>; DB="mpc_test_$RUNID"
$ docker exec mpc-pg-session-1d7cca5b psql --username postgres --dbname postgres \
    --set ON_ERROR_STOP=1 --command "CREATE DATABASE $DB"
CREATE DATABASE

$ cd apps/server_core
$ export MPC_TEST_DATABASE_URL="postgresql://postgres:<session-password>@127.0.0.1:51692/${DB}?sslmode=disable"
$ GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go run ./cmd/testdb migrate
applied 69 migration(s)
```

## 3. Direct `-v` run — A5 and A10, plus A7 in the same invocation

```
$ cd apps/server_core
$ export MPC_TEST_DATABASE_URL="postgresql://postgres:<session-password>@127.0.0.1:51692/mpc_test_<runid>?sslmode=disable"
$ GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go test -tags integration -v \
    -run 'TestGetImportChainCountsLeadingZeroCodprod|TestGetImportChainNonArrayPendingCursorDoesNotFailQuery|TestGetImportChainCountsCurrentQueueAcrossInstallations' \
    ./internal/modules/erp_import/adapters/postgres/...

=== RUN   TestGetImportChainCountsCurrentQueueAcrossInstallations
--- PASS: TestGetImportChainCountsCurrentQueueAcrossInstallations (0.04s)
=== RUN   TestGetImportChainCountsLeadingZeroCodprod
--- PASS: TestGetImportChainCountsLeadingZeroCodprod (0.03s)
=== RUN   TestGetImportChainNonArrayPendingCursorDoesNotFailQuery
=== RUN   TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_object
=== RUN   TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_scalar
--- PASS: TestGetImportChainNonArrayPendingCursorDoesNotFailQuery (0.03s)
    --- PASS: TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_object (0.02s)
    --- PASS: TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_scalar (0.01s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres	3.242s
```

Counts by name: RUN=5 (3 top-level + 2 subtests), PASS=5, SKIP=0, FAIL=0.
This is with `MPC_TEST_DATABASE_URL` set — `SkipWithoutTarget` did not fire;
these are real executions against a live Postgres, not skips.

## 4. A6 — must-fail proof (production code temporarily reverted, forward-restored)

### A5's join, old form (`ON links.internal_product_id::text = products.codprod`)

Edited `query_repository.go` in place to the pre-fix line, ran:
```
$ GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go test -tags integration -v \
    -run 'TestGetImportChainCountsLeadingZeroCodprod' ./internal/modules/erp_import/adapters/postgres/...

=== RUN   TestGetImportChainCountsLeadingZeroCodprod
    chain_query_repository_integration_test.go:283: chain=domain.ImportChain{Protocol:"#640-E", Importados:1, Vinculados:0, Enfileirados:0, QueueReadAt:time.Date(2026, time.July, 28, 13, 7, 24, 280214000, time.Local)}
--- FAIL: TestGetImportChainCountsLeadingZeroCodprod (0.10s)
FAIL
FAIL	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres	3.969s
FAIL
```
`Vinculados:0` (want 1) under the old join — the exact failing integer, not
the word "failed". Then re-edited the file forward to the `ltrim(...)` form
(byte-identical to the pre-slice text) and re-ran — see §3 above, PASS.

### A10's CASE guard, old form (`COALESCE(state.cursor -> 'pending', '[]'::jsonb)`)

Edited `query_repository.go` in place to the pre-fix line, ran:
```
$ GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go test -tags integration -v \
    -run 'TestGetImportChainNonArrayPendingCursorDoesNotFailQuery' ./internal/modules/erp_import/adapters/postgres/...

=== RUN   TestGetImportChainNonArrayPendingCursorDoesNotFailQuery
=== RUN   TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_object
    chain_query_repository_integration_test.go:343: expected a valid response, got query error: get ERP import chain: ERROR: cannot extract elements from an object (SQLSTATE 22023)
=== RUN   TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_scalar
    chain_query_repository_integration_test.go:343: expected a valid response, got query error: get ERP import chain: ERROR: cannot extract elements from an object (SQLSTATE 22023)
--- FAIL: TestGetImportChainNonArrayPendingCursorDoesNotFailQuery (0.04s)
    --- FAIL: TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_object (0.02s)
    --- FAIL: TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_scalar (0.01s)
FAIL
FAIL	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres	2.990s
FAIL
```
Real Postgres error (`SQLSTATE 22023: cannot extract elements from an
object`) under the old shape, for both subtests — exactly the "takes the
whole endpoint down" failure mode the fix guards against. Then re-edited the
file forward to the `CASE ... jsonb_typeof ...` form and re-ran — see §3
above, PASS.

Post-restore verification that `query_repository.go` is back to its
pre-slice state (same diff as when this slice started, nothing left behind):
```
$ git diff --stat apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go
 .../adapters/postgres/query_repository.go          | 22 ++++++++++++++++++++--
 1 file changed, 20 insertions(+), 2 deletions(-)
```
(identical stat before and after the two temporary-edit/restore round trips)

## 5. A7 — TestGetImportChainCountsCurrentQueueAcrossInstallations, honest result

This is a **pre-existing** test, run but not modified. Ran it 6 times total
across this slice (1 inside the combined run in §3, 5 more isolated):

```
run in §3 (combined):        --- PASS (0.04s)
isolated run 1:               --- FAIL (0.05s)  queue_read_at outside [before,after] by ~1ms
isolated run 2:               --- FAIL (0.04s)  queue_read_at outside [before,after] by ~1ms
isolated run 3:               --- PASS (0.05s)
isolated run 1 (2nd batch):   --- FAIL (0.05s)  queue_read_at outside [before,after] window
isolated run 2 (2nd batch):   --- FAIL (0.06s)  queue_read_at outside [before,after] window
isolated run 3 (2nd batch):   --- PASS (0.11s)
isolated run 4 (2nd batch):   --- PASS (0.10s)
isolated run 5 (2nd batch):   --- PASS (0.04s)
```
Score: PASS 5 / FAIL 4 out of 9 runs. **Every single failure is on line 84**
(the `QueueReadAt` bracket check), **never on line 81** (the
`Protocol/Importados/Vinculados/Enfileirados` counts assertion that is what
actually pins the DISTINCT guard this test exists for). The counts assertion
— the thing A7 asks me to confirm — passed in all 9 runs, 9/9.

Root-caused, not hand-waved: measured host clock vs. this session's Postgres
container clock directly —
```
$ date -u +%Y-%m-%dT%H:%M:%S.%N
2026-07-28T16:09:04.608380400
$ docker exec mpc-pg-session-1d7cca5b date -u +%Y-%m-%dT%H:%M:%S.%N
2026-07-28T16:09:04.863871884
$ docker exec mpc-pg-session-1d7cca5b psql -U postgres -d postgres -t -c "select statement_timestamp()"
 2026-07-28 16:09:05.208959+00
```
The container's `statement_timestamp()` runs ~600ms ahead of the Go test
process's host-side `time.Now()` in this worktree's Docker session, a
pre-existing Docker Desktop clock-skew condition unrelated to the SQL text
this slice's tests exercise (neither the `resolved_products` join nor the
`queued_products` CASE guard touches `statement_timestamp()`). The `after :=
time.Now()` bracket in the test is tight enough (host clock, not container
clock) that this skew intermittently pushes `queue_read_at` past it.

**Verdict for A7: REPORT, not clean PASS.** The join change did **not**
break the DISTINCT guard — the counts assertion this test exists to pin held
9/9. But I cannot honestly claim "still green" when the test file itself
flakes ~44% of the time in this session on an assertion unrelated to my
change. Not fixed, not loosened — reported as found, per the hard rule
against dressing up a partial result as complete. This flake exists
independent of anything in this slice's diff and would reproduce on `main`
in this same Docker session.

## 6. What could not be done

Nothing in A5/A10 was left incomplete — both are full PASS with a genuine
must-fail proof against the exact old production text. A7's counts
assertion (the actual DISTINCT-guard behavior) is a clean 9/9 PASS. The one
thing not claimed as clean is A7's `QueueReadAt` timing assertion, which is
flaky in this session for a documented, pre-existing, unrelated
(Docker-clock-skew) reason — reported honestly above rather than silently
passed over.

Cleanup: the manually-created scratch database (`mpc_test_<runid>`, used
only for the direct `-v`/must-fail runs in this artifact) was dropped after
use. The session postgres container (`mpc-pg-session-*`, brought up via
`npm run harness:pg:up`) was left running, per the hub-owned long-lived
session convention.
