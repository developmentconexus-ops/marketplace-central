# CHIP-ANCHORS-3 — impl R4 + R5

Slice: R4 (`queued_products` join canonicalization + regression test) and R5 (one sentence on
what `buildConcordantCandidate`'s degradation produces).

Files touched (exactly three, all in the allowed set):

- `apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go`
- `apps/server_core/internal/modules/erp_import/adapters/postgres/chain_query_repository_integration_test.go`
- `apps/server_core/internal/modules/product_links/application/generation_service.go`

Zero `apps/web/`, zero migrations, zero `platform/httpx`. No commit, no push, no
reset/revert/stash/clean/checkout.

---

## Item 1 — `queued_products` identity rule

### SQL before (verbatim, `query_repository.go` lines 102–120 pre-edit)

```sql
		queued_products AS (
			-- COALESCE only defends against NULL. A cursor whose 'pending' is an
			-- object or a scalar makes jsonb_array_elements_text raise at query
			-- time and takes the whole endpoint down for the tenant, so the type
			-- is checked before expansion. CASE (not AND) because a join/filter
			-- predicate gives no evaluation-order guarantee.
			SELECT DISTINCT products.codprod AS codprod
			FROM sync_state AS state
			CROSS JOIN LATERAL jsonb_array_elements_text(
				CASE
					WHEN jsonb_typeof(state.cursor -> 'pending') = 'array'
						THEN state.cursor -> 'pending'
					ELSE '[]'::jsonb
				END
			) AS pending(codprod)
			JOIN import_products AS products ON products.codprod = pending.codprod
			WHERE state.tenant_id = $1
			  AND state.entity = 'market'
		)
```

### SQL after (verbatim)

```sql
		queued_products AS (
			-- importados / vinculados / enfileirados are read on one screen as a
			-- decomposition of a single population, so the two joined counters have
			-- to agree on what makes two CODPRODs the same product. resolved_products
			-- already answers that canonically; leaving this side raw would let one
			-- padded CODPROD be linked but not queued, and the operator would read
			-- the gap between two numbers as a stalled queue that never existed. The
			-- counted key stays the raw codprod, which is what importados counts —
			-- only the identity test is canonicalized. Text-to-text for the same
			-- reason as above: codprod is not guaranteed numeric.
			--
			-- COALESCE only defends against NULL. A cursor whose 'pending' is an
			-- object or a scalar makes jsonb_array_elements_text raise at query
			-- time and takes the whole endpoint down for the tenant, so the type
			-- is checked before expansion. CASE (not AND) because a join/filter
			-- predicate gives no evaluation-order guarantee.
			SELECT DISTINCT products.codprod AS codprod
			FROM sync_state AS state
			CROSS JOIN LATERAL jsonb_array_elements_text(
				CASE
					WHEN jsonb_typeof(state.cursor -> 'pending') = 'array'
						THEN state.cursor -> 'pending'
					ELSE '[]'::jsonb
				END
			) AS pending(codprod)
			JOIN import_products AS products
			  ON ltrim(products.codprod, '0') = ltrim(pending.codprod, '0')
			WHERE state.tenant_id = $1
			  AND state.entity = 'market'
		)
```

Constraints honoured as adjudicated:

- `SELECT DISTINCT products.codprod` unchanged in BOTH CTEs — only the JOIN predicate is
  canonicalized. I agree with the adjudication: `importados` is `count(*)` over raw
  `import_products` rows, so the raw codprod is the only counted key consistent with it.
  I have no dissent to record.
- The `CASE` / `jsonb_typeof` guard is untouched — its paragraph is preserved verbatim and its
  own regression test (`TestGetImportChainNonArrayPendingCursorDoesNotFailQuery`) still passes.
- The new comment paragraph says WHY (the two counters must agree on identity, and what the
  operator misreads when they do not) and does not narrate the SQL line below it.

### Fixture direction chosen, and why

Chosen: **`erp_import_products.codprod = '00101'`, cursor `pending` containing the unpadded
`"101"`.**

Reason, from the code that writes each side:

- `apps/server_core/internal/modules/erp_import/application/import_service.go:153-158` builds
  the enqueue slice as `row.Codprod` for every accepted row and hands it to
  `EnqueueMarketProducts`; `apps/server_core/internal/modules/erp_import/adapters/sync/enqueuer.go:44`
  passes it straight to `AppendPendingCodigos`, which appends it as `to_jsonb($4::text[])`
  (`apps/server_core/internal/modules/sync/adapters/postgres/sync_state_repo.go:138-158`).
  That producer writes the SAME raw string that `erp_import_products` stores, so it can never
  disagree with itself — a fixture in that direction would not exercise the bug.
- The padding is therefore lost only by the OTHER producer(s) of the cursor —
  `AppendPendingCodigo` singular exists for exactly that (its own doc comment names "two
  concurrent producers (e.g. an M-03 import hook and an M-07 enqueue)"), and anything sourcing
  codes from the integer-keyed product side loses the leading zero. That is the identical loss
  `resolved_products` already had to canonicalize (`product_links.internal_product_id` is an
  integer column, so `'00101'` linked as `101`).

I added a **second** fixture row in the symmetric direction — raw `'102'` in
`erp_import_products` against `"00102"` in the cursor — so a half-applied fix (`ltrim` on only
one side of the comparison) also fails. That is one extra row in the same test, not a second
test.

### New test code (verbatim, appended to `chain_query_repository_integration_test.go`)

```go
// TestGetImportChainCountsLeadingZeroCodprodInQueue pins the queued_products
// join fix. importados, vinculados and enfileirados are read on one screen as
// a decomposition of the same population, so the two joined counters have to
// agree on what makes two CODPRODs the same product. resolved_products already
// compares canonically; while queued_products compared raw, a padded '00101'
// was counted as vinculados and NOT as enfileirados, and the operator read the
// gap between the two numbers as a queue that had stalled.
//
// The fixture drives the mismatch from the cursor side — '00101' imported, the
// cursor carrying the unpadded "101" — because that is how it enters: the
// post-import hook writes the accepted rows' raw CODPROD straight back into
// the cursor (erp_import/adapters/sync.MarketEnqueuer), so that producer can
// never disagree with itself; the padding is lost by any other producer that
// sources codes from the integer-keyed product side, which is the same loss
// resolved_products already had to canonicalize. The second row ("00102"
// pending against a raw '102') pins the symmetric direction, so the fix cannot
// be half-applied to one side of the comparison.
func TestGetImportChainCountsLeadingZeroCodprodInQueue(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	importID := domain.ImportID("66000000-0000-0000-0000-000000000001")

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_links WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_products WHERE tenant_id=$1`, tenant)
	})

	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_protocols
			(id, tenant_id, file_sha256, protocol, source, imported_at, status, accepted_count, rejected_count, warning_count)
		VALUES ($1, $2, 'chain-queue-leading-zero-hash', '#660-E', 'xlsx', now(), 'COMPLETED', 3, 0, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_products (tenant_id, protocol_id, codprod, descrprod, custo, stock_physical)
		VALUES
			($2, $1, '00101', 'Product 00101', 1, 0),
			($2, $1, '102', 'Product 102', 1, 0),
			($2, $1, '103', 'Product 103', 1, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	// '00101' is linked, so vinculados counts it. enfileirados has to count the
	// same row from the same cursor entry, or the two numbers stop describing
	// the same product.
	if _, err := pool.Exec(ctx, `
		INSERT INTO product_links
			(tenant_id, installation_id, provider_code, provider_item_id, state, internal_product_id)
		VALUES ($1, 'installation-a', 'mercadolivre', 'item-00101', 'resolved', 101);
	`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, cursor)
		VALUES ($1, 'installation-a', 'market', '{"pending":["101","00102","OUTSIDE"]}'::jsonb);
	`, tenant); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetImportChain(ctx, tenant, importID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Protocol != "#660-E" || got.Importados != 3 || got.Vinculados != 1 || got.Enfileirados != 2 {
		t.Fatalf("chain=%#v", got)
	}
}
```

Shape follows the existing `TestGetImportChainCountsLeadingZeroCodprod`: same
`integrationRepo(t)` helper (which calls `testpostgres.SkipWithoutTarget`), same
`testpostgres.OpenPool`, same `t.Cleanup` deletes, same `//go:build integration` file.
`"OUTSIDE"` is the negative control (queued but not imported → not counted).

### Real-Postgres run (direct invocation, NOT the lane artifact)

Target: the same ephemeral harness Postgres the lane uses — this checkout's session container
`mpc-pg-session-1d7cca5b` (`scripts/.runs/pg-session.json`, the same container
`Invoke-Integration` reuses when it prints `container=session-reuse`). I created a
harness-shaped database (`^mpc_test_[0-9a-f]{32}$`, enforced by
`internal/testsupport/postgres/target.go:34`), migrated it with `go run ./cmd/testdb migrate`
(`applied 69 migration(s)`), ran the tests, and dropped the database afterwards. The session
container was already running before this slice and is left running, untouched.

I did NOT close on `npm run harness:integration` — as the dispatch says, its artifact cannot
distinguish a fully-skipped run from a fully-green one.

```
$ cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache \
  go test -tags=integration -count=1 -v -run 'TestGetImportChain' ./internal/modules/erp_import/adapters/postgres/

=== RUN   TestGetImportChainCountsCurrentQueueAcrossInstallations
--- PASS: TestGetImportChainCountsCurrentQueueAcrossInstallations (0.04s)
=== RUN   TestGetImportChainHTTPIntegration
--- PASS: TestGetImportChainHTTPIntegration (0.03s)
=== RUN   TestGetImportChainMissingAndProtocolWithoutSyncState
--- PASS: TestGetImportChainMissingAndProtocolWithoutSyncState (0.04s)
=== RUN   TestGetImportChainCountsLeadingZeroCodprod
--- PASS: TestGetImportChainCountsLeadingZeroCodprod (0.03s)
=== RUN   TestGetImportChainCountsLeadingZeroCodprodInQueue
--- PASS: TestGetImportChainCountsLeadingZeroCodprodInQueue (0.03s)
=== RUN   TestGetImportChainNonArrayPendingCursorDoesNotFailQuery
=== RUN   TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_object
=== RUN   TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_scalar
--- PASS: TestGetImportChainNonArrayPendingCursorDoesNotFailQuery (0.03s)
    --- PASS: TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_object (0.02s)
    --- PASS: TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_scalar (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres	3.222s
EXIT=0
```

The load-bearing line: `--- PASS: TestGetImportChainCountsLeadingZeroCodprodInQueue (0.03s)`.
It is a real PASS, not a skip — a skip would print `--- SKIP` with
`MPC_TEST_DATABASE_URL is unset`.

### Must-fail (mandatory)

Mutation: the join predicate reverted to the raw form, everything else untouched.

```sql
			JOIN import_products AS products ON products.codprod = pending.codprod
```

Same command, same database. VERBATIM output:

```
=== RUN   TestGetImportChainCountsCurrentQueueAcrossInstallations
--- PASS: TestGetImportChainCountsCurrentQueueAcrossInstallations (0.06s)
=== RUN   TestGetImportChainHTTPIntegration
--- PASS: TestGetImportChainHTTPIntegration (0.04s)
=== RUN   TestGetImportChainMissingAndProtocolWithoutSyncState
--- PASS: TestGetImportChainMissingAndProtocolWithoutSyncState (0.03s)
=== RUN   TestGetImportChainCountsLeadingZeroCodprod
--- PASS: TestGetImportChainCountsLeadingZeroCodprod (0.05s)
=== RUN   TestGetImportChainCountsLeadingZeroCodprodInQueue
    chain_query_repository_integration_test.go:354: chain=domain.ImportChain{Protocol:"#660-E", Importados:3, Vinculados:1, Enfileirados:0, QueueReadAt:time.Date(2026, time.July, 28, 14, 18, 13, 715864000, time.Local)}
--- FAIL: TestGetImportChainCountsLeadingZeroCodprodInQueue (0.09s)
=== RUN   TestGetImportChainNonArrayPendingCursorDoesNotFailQuery
=== RUN   TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_object
=== RUN   TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_scalar
--- PASS: TestGetImportChainNonArrayPendingCursorDoesNotFailQuery (0.04s)
    --- PASS: TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_object (0.03s)
    --- PASS: TestGetImportChainNonArrayPendingCursorDoesNotFailQuery/pending_as_scalar (0.01s)
FAIL
FAIL	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres	2.765s
FAIL
EXIT=1
```

**The actual `Enfileirados:` number under the raw join is `Enfileirados:0`**, on the same read
where `Vinculados:1` and `Importados:3`. That is the defect stated in one line: the same
CODPROD is linked and simultaneously not queued.

Note the OTHER four chain tests still PASS under the mutation — the mutation is targeted, and
only the new test detects it.

### Restore

Restored by editing FORWARD (Edit tool, re-applying the two-line canonicalized predicate). No
`git checkout`, `reset`, `revert`, `stash` or `clean` was run at any point in this slice.

Post-restore rerun, same command, same database:

```
=== RUN   TestGetImportChainCountsLeadingZeroCodprodInQueue
--- PASS: TestGetImportChainCountsLeadingZeroCodprodInQueue (0.16s)
...
PASS
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres	3.443s
EXIT=0
```

`git status --short` after restore — only the three intended files, plus pre-existing untracked
`.mnfs` chip artifacts that were already there when this slice started:

```
 M apps/server_core/internal/modules/erp_import/adapters/postgres/chain_query_repository_integration_test.go
 M apps/server_core/internal/modules/erp_import/adapters/postgres/query_repository.go
 M apps/server_core/internal/modules/product_links/application/generation_service.go
?? .mnfs/MIS-006-integracao-fundacao/_chip-anchors-3/EVIDENCE.md
?? .mnfs/MIS-006-integracao-fundacao/_chip-anchors-3/dispatches/
?? .mnfs/MIS-006-integracao-fundacao/_chip-anchors-3/p6-input-r1.patch
?? .mnfs/MIS-006-integracao-fundacao/_chip-anchors-3/p6-input-r2.patch
```

`git diff --stat`:

```
 .../chain_query_repository_integration_test.go     | 71 ++++++++++++++++++++++
 .../adapters/postgres/query_repository.go          | 13 +++-
 .../application/generation_service.go              |  7 ++-
 3 files changed, 89 insertions(+), 2 deletions(-)
```

The `query_repository.go` diff contains exactly the comment paragraph plus the predicate line —
no residue of the mutation.

---

## Item 2 — what the degradation produces

### Verifying BOTH halves of the existing claim before touching it

Half A — "this site really did deref unconditionally". CONFIRMED from history, not from
memory. `git log -L 489,496:...generation_service.go`:

- `9c030154 feat(links): add incomparable anchor reasons` introduced the function body as
  `product := *comparison.product` — no guard.
- `9555ad30 chip(CHIP-ANCHORS-3): CORR-1/CORR-4/CORR-6b` replaced that line with the zero-value
  + nil-check and wrote the comment.

Half B — "both sibling scorers nil-check this pointer". CONFIRMED with a caveat on wording.
`grep -n 'comparison\.product'` returns exactly four hits in the file: lines 493–494 (this
site), 519–520 (`applySingleAnchorScore`: `if comparison.product != nil { product = *comparison.product }`),
and 714/718 (`classifyProviderIdentityAnchor`: passes the pointer to `identityAnchorValues`,
then `if comparison.product == nil`). So **every other consumer of the pointer nil-checks it**,
and the count "both" (two other consumers) is right. The caveat: only one of those two,
`applySingleAnchorScore`, is literally a scorer; `classifyProviderIdentityAnchor` is a reason
classifier. `applyUnresolvedScore` is nil-safe by a different route — it passes a literal `nil`
to `missingMatchedAnchorReason`, which handles `product == nil`. I did NOT edit that wording —
it is loose, not false, and rewording it is outside this slice.

### Verifying what this site actually produces on the degraded path

Read and traced, not assumed. With `comparison.product == nil`:

- `product` is the zero `internalreaddomain.ProductCandidate{}`.
- `newCandidate` (line 396) → `canonicalProductID` returns `(0, false)` because
  `InternalProductID == nil` (line 465), so `InternalProductID` on the candidate stays `nil`;
  `InternalProductName` and `InternalReferenceCode` are empty.
- `reasons` is unconditionally seeded with `seller_sku` FOR and `ean` FOR (lines 498–501).
- `detectHardNegative(snapshot.Title, product.Name)` returns `(false, "")` immediately because
  `in == ""` (line 813), so the REJECT branch cannot be reached on the degraded path.
- Therefore `Confidence = 95`, `ConfidenceBand = ALTA`, `MatchStatus = ACCEPT` (lines 508–510).
- `autoApprovals` (line 236) selects on `MatchStatus == ACCEPT` and emits an auto-approval.

So the degradation is not parity with the siblings' absence reasons: it is two FOR reasons and
an ALTA/ACCEPT for a CODPROD that is not there, on a row whose `internal_product_id` is null.

### Comment — exact OLD text

```go
	// Both sibling scorers nil-check this pointer; an unconditional deref here
	// made this the one site that panics instead of degrading.
```

### Comment — exact NEW text

```go
	// Both sibling scorers nil-check this pointer; an unconditional deref here
	// made this the one site that panics instead of degrading. What it degrades
	// INTO is not the siblings' absence reasons: on the zeroed candidate this
	// function still emits seller_sku FOR and ean FOR at 95 / ALTA / ACCEPT,
	// asserting corroboration for a CODPROD that is not there. The row carries a
	// null internal_product_id and autoApprovals still reads that ACCEPT as
	// auto-approvable.
```

Two sentences added. The existing text is preserved verbatim and is not annotated as wrong —
it is true, only narrower than it reads. No code change, no test change in this file.

---

## Verification — all four, post-restore, with the mandated cd/GOCACHE form

Command form used for every one:
`cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go ...`
(no `-mod=mod` anywhere).

```
### 1 go build ./...
exit=0
### 2 go vet ./...
exit=0
### 3 go vet -tags=integration ./internal/modules/erp_import/...
exit=0
### 4 go test -count=1 ./internal/modules/erp_import/... ./internal/modules/product_links/...
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread	4.730s
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres	3.519s
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/productlinks	3.231s
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/sync	1.583s
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/adapters/xlsx	4.595s
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/application	3.199s
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/domain	1.836s
?   	marketplace-central/apps/server_core/internal/modules/erp_import/ports	[no test files]
ok  	marketplace-central/apps/server_core/internal/modules/erp_import/transport	3.680s
ok  	marketplace-central/apps/server_core/internal/modules/product_links/adapters/connectors	4.050s
?   	marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres	[no test files]
ok  	marketplace-central/apps/server_core/internal/modules/product_links/application	4.047s
ok  	marketplace-central/apps/server_core/internal/modules/product_links/composition	3.029s
ok  	marketplace-central/apps/server_core/internal/modules/product_links/domain	3.003s
?   	marketplace-central/apps/server_core/internal/modules/product_links/ports	[no test files]
ok  	marketplace-central/apps/server_core/internal/modules/product_links/transport	3.778s
exit=0
```

Build 1–3 printed no diagnostics; the `exit=0` lines are the whole output.

---

## What I could not verify, and why

1. **The degraded path of `buildConcordantCandidate` is currently UNREACHABLE, and I did not
   prove otherwise.** Its only caller (`buildExactCandidates`,
   `generation_service.go:303`) always passes `&product` where `product :=
   skuMatches.Products[0]`, so `comparison.product` is never nil today. Everything I wrote in
   the new sentences is derived by reading the function and its callees, NOT by executing the
   nil path — there is no test that constructs a nil-product comparison and observes 95 / ALTA
   / ACCEPT. The claim is a code-reading claim about what the code WOULD produce. The slice
   forbade a test change, so I did not add one; if the gate wants it executed rather than read,
   that is a follow-up.
2. **I did not verify that any real producer writes an unpadded CODPROD into
   `sync_state.cursor->'pending'`.** I proved the negative (the import hook writes the raw
   string and therefore cannot mismatch) and I identified `AppendPendingCodigo` singular as the
   second producer named in its own doc comment — but `grep` finds NO non-test caller of
   `AppendPendingCodigo` in the tree today. So the padded-vs-unpadded cursor mismatch is, right
   now, a latent inconsistency between two counters rather than an observed production
   incident. What IS proven is the invariant: the two counters disagreed on identity, and now
   they do not.
3. **`ltrim(x, '0')` collides all-zero codprods.** `'0'`, `'00'` and `'000'` all canonicalize
   to the empty string, so they would be treated as the same product by this predicate. That
   hazard already existed in `resolved_products` (shipped at `bba08b41`); I extended it to
   `queued_products` for consistency rather than introducing a third identity rule. I did not
   test it and I did not fix it — a codprod of all zeros is presumably not valid ERP data, but
   I did not confirm that against `IsValidCodprod` or against real Sankhya data.
4. **I did not run `npm run harness:integration`.** By design — the dispatch rules out closing
   on its artifact, and running it would have added no evidence the direct `-v` run does not
   already carry. The full integration suite beyond `-run 'TestGetImportChain'` was therefore
   NOT executed against Postgres by me; only the six chain tests were. `go vet -tags=integration`
   proves the rest of the module's integration files still compile, not that they pass.
5. **The FE / `apps/web` reading of these three counters was not inspected.** Out of scope by
   the slice's hard rules (zero `apps/web/`), so I cannot say whether the UI derives anything
   from the difference between `vinculados` and `enfileirados`.
6. **Housekeeping note, not a verification gap:** to get a direct real-Postgres run I created
   and then dropped a database (`mpc_test_0123456789abcdef0123456789abcdef`) inside this
   checkout's pre-existing session container `mpc-pg-session-1d7cca5b`. The container was
   already running before this slice and is still running; I did not start it, stop it, or
   touch `scripts/.runs/pg-session.json`. The container password was printed to the session
   transcript once while I was locating the target; it is a throwaway tmpfs-backed local test
   container, but the hub should know it is in the transcript.
