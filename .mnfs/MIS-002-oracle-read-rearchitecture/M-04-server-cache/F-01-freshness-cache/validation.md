# Feature Validation

```yaml
id: F-01
type: feature-validation
status: correction_validation_evidence_collected
owner: Correction Worker
parent: M-04
created: 2026-07-14
updated: 2026-07-14
validation_level: targeted-plus-full-suite
lifecycle_scope: feature
```

## Summary

The sixth correction round is test/evidence-only and addresses three
evidence-level defects in the prior evidence. The R1 post-invalidation singleflight regression test was
initially hollow: it could pass after the generation component was removed
from the normal-path group key. The test now deterministically observes
downstream entry for both readers before releasing the first load. The
mutation test below proves the rewritten test fails when the guarded behavior
is removed and passes after the committed production code is restored. This
document records feature evidence; it is not a milestone or formal QA verdict.

## Corrections

### Fourth correction round

1. H1 linkage-exclusion correction: `TestFreshnessCacheBypassAndLinkageExclusion`
   now repeats the identical `FindProductsInput` three times. This preserves
   the existing bypass/repopulation assertion and makes the linkage assertion
   mutation-sensitive: correct uncached behavior requires three downstream
   calls, while a cached linkage result produces one.

2. H2 failed-mutation correction: the tautological local `failed` error and
   `invalidations` counter block was deleted from `TestEvictOnMutation`. That
   block never called a mutation service or invalidated a cache and therefore
   proved nothing. Failed-write/no-invalidation evidence remains at the
   application layer, where the mutation paths actually execute.

### Fifth correction round

1. G1 nil-vs-absent correction: `TestFreshnessCacheTTLPerClass` now uses the
   comma-ok map lookup for cost key `1`, requiring the key to be present and its
   value to be explicitly nil. It separately checks that key `2` remains
   present with a value, preserves the error and one-downstream-call checks,
   and reports key absence separately from a non-nil value. The cost fake now
   returns `{1: nil, 2: &CostAsOf{}}`. The deep-copy cost, tax, and stock map
   accesses were audited; none is a nil-preservation assertion, so no other
   nil assertion required correction.

2. G2 log-hygiene correction: `TestFreshnessCacheLRUAndLogs` now parses JSON
   cache records and requires the exact attribute set `{time, level, msg,
   cache, key_class}`. It still positively observes `cache=miss`, `cache=hit`,
   and `key_class=catalog`. Any raw argument, ID, value, or digest attribute
   therefore fails the test.

#### G1 mutation proof

`cache.go` was temporarily changed to omit nil values from `cloneCostFacts`.
The required focused test then failed because the comma-ok assertion detected
the absent key:

```text
--- FAIL: TestFreshnessCacheTTLPerClass (0.01s)
    cache_test.go:215: pricecost cache hit dropped unknown fact key: calls=1 err=<nil>
FAIL
FAIL	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache	1.422s
```

After restoring `cache.go`, `git diff --
apps/server_core/internal/modules/internal_read/adapters/cache/cache.go` was
empty and the same focused test passed:

```text
ok  	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache	1.426s
```

#### G2 mutation proof

`cache.go` was temporarily changed to add a `key` attribute to each cache log.
The required focused test failed because the exact allowlist detected the
unexpected attribute:

```text
--- FAIL: TestFreshnessCacheLRUAndLogs (0.01s)
    cache_test.go:638: log record 1 has unexpected attribute "key": {"time":"2026-07-14T15:57:28.7630961-03:00","level":"INFO","msg":"internal_read.cache","cache":"miss","key_class":"catalog","key":"raw-key"}
FAIL
FAIL	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache	1.668s
```

After restoring `cache.go`, its diff was empty and the same focused test
passed:

```text
ok  	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache	1.329s
```

#### H1 linkage-exclusion mutation proof

`cache.go` was temporarily changed to route `FindProductsForLinking` through
the catalog cache. With that deliberate production mutation, the required
focused test failed:

```text
=== RUN   TestFreshnessCacheBypassAndLinkageExclusion
    cache_test.go:579: linkage was cached: calls=1
--- FAIL: TestFreshnessCacheBypassAndLinkageExclusion (0.01s)
FAIL
FAIL	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache	0.923s
```

After restoring `cache.go` exactly, `git diff --
apps/server_core/internal/modules/internal_read/adapters/cache/cache.go` was
empty and the same focused test passed:

```text
=== RUN   TestFreshnessCacheBypassAndLinkageExclusion
--- PASS: TestFreshnessCacheBypassAndLinkageExclusion (0.01s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache	1.310s
```

### Sixth correction round

1. C03 now has composed `TestComposedCatalogHTTPNoCacheBypassesAndRepopulates`
   in `tests/unit/cache_composed_test.go`: real catalog transport handler →
   real catalog cache decorator → fake Oracle page reader. It proves two warm
   GETs make one Oracle call with identical response `as_of`, `Cache-Control:
   no-cache` makes call two with a strictly newer `as_of`, and a following
   ordinary GET makes no third call while returning the bypass-refreshed
   `as_of`. The older cache-only bypass-repopulation assertion at
   `cache_test.go:585` is insensitive; this response-body assertion supplies
   the mutation-sensitive coverage.
2. C05 now has
   `TestAssistedSankhyaConfirmInvalidatesCatalogAfterSuccessfulPersistence`,
   requiring exactly `[]string{"catalog"}`, and composed
   `TestComposedLinkageConfirmEvictsCatalogButPreservesPriceCost`, proving the
   real handler → cache → successful `Confirm` flow evicts the warm catalog
   key inside TTL while preserving a warm `pricecost` key.
3. C02 evidence is stated precisely: `TestFreshnessCacheSingleflight` is a
   weak guard because it has no barrier proving all 19 other waiters parked;
   QA observed it catches singleflight deletion only about 15/100 runs.
   `TestFreshnessCacheErrorNotCached`, with the readiness barrier, is the real
   C02 mutation guard and caught that removal 20/20.

#### Sixth-round mutation evidence

The C03 bypass mutation failed at `bypass Oracle calls=1, want 2`; after
restoration the composed test passed. Removing the success-path
`InvalidateClass("catalog")` call failed the direct C05 assertion with
`invalidations=[], want exactly [catalog]` and independently failed the
composed C05 assertion with `post-confirm catalog Oracle calls=1, want 2`.
Both production files were restored to empty diffs and the three tests passed.

### Prior correction evidence, retained and not regressed

1. R1 test correction: `stagedCatalogReader` now reports every catalog
   downstream entry through a buffered `entered` channel. The regression test
   confirms call #1 is in flight, invalidates the class, advances the fake
   clock, starts reader B, and requires call #2 to be observed before releasing
   call #1. It then verifies two downstream calls and strictly increasing
   `AsOf` values. The R2 snapshot-guard and R3 invalidation-guard tests were
   left unchanged.

#### R1 regression-test mutation proof

The normal-path group key in `cache.go` was temporarily changed from
`strconv.FormatUint(generation, 10)` to the literal `"0"`. The required
command was:

```powershell
go test ./internal/modules/internal_read/adapters/cache/ -run PostInvalidationDoesNotJoinInFlightLoad -count=1
```

It then failed at the second-entry detector:

```text
2026/07/14 14:35:37 INFO internal_read.cache cache=miss key_class=catalog
2026/07/14 14:35:37 INFO internal_read.cache cache=miss key_class=catalog
--- FAIL: TestFreshnessCachePostInvalidationDoesNotJoinInFlightLoad (2.01s)
    cache_test.go:330: post-invalidation reader joined the pre-invalidation in-flight load
FAIL
FAIL	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache	3.171s
FAIL
```

The production line was restored exactly to
`strconv.FormatUint(generation, 10)`. `git diff --
apps/server_core/internal/modules/internal_read/adapters/cache/cache.go` was
empty, and the same focused command, with the required absolute `GOCACHE`,
passed:

```text
ok  	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache	1.388s
```

### Prior round, retained and not regressed

1. P1-A generation fence: an in-flight load captured before invalidation may
   return to its original caller but cannot populate the new generation.
2. P1-B bypass namespace: MaxAge=0 bypasses L2 and repopulates the regular
   cache without joining a normal caller's singleflight double-check.
3. P1-C deep clone: catalog pages and cost, tax, and stock maps clone all
   mutable pointers, slices, nested observed timestamps, nil maps/slices, and
   nil map values on store and return; the cost nil map value is specifically
   proven by the comma-ok assertion in `TestFreshnessCacheTTLPerClass`.
4. P2-D waiter cancellation: `DoChan` plus context selection lets canceled
   waiters return without a helper goroutine while live waiters still collapse.
5. P2-E persistence evidence: catalog/inventory/orders/profitability failure
   paths retain evidence that failed persistence does not invalidate; inventory
   also retains provider-rejection coverage.

## Changed Paths

This correction round permanently changes only:

- `apps/server_core/tests/unit/cache_composed_test.go`
- `apps/server_core/internal/modules/orders/application/assisted_sankhya_linkage_service_test.go`
- `.mnfs/MIS-002-oracle-read-rearchitecture/M-04-server-cache/validation-contract.md`
- `.mnfs/MIS-002-oracle-read-rearchitecture/M-04-server-cache/F-01-freshness-cache/validation.md`

`cache.go` was changed only temporarily for the mutation proof and was
restored with an empty diff.

## Commands Run

All Go commands used this absolute cache path:

```powershell
$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\.gocache'
```

- `go test ./internal/modules/internal_read/adapters/cache/ -count=1 -v`
  - Pass. The complete cache package passed, including:

    ```text
    --- PASS: TestFreshnessCacheBypassAndLinkageExclusion
    --- PASS: TestFreshnessCacheLRUAndLogs
    --- PASS: TestEvictOnMutation
    --- PASS: ExampleCache
    PASS
    ok  	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache	1.438s
    ```

- `go test ./apps/server_core/tests/unit ./apps/server_core/internal/modules/orders/application -run 'TestComposedCatalogHTTPNoCacheBypassesAndRepopulates|TestComposedLinkageConfirmEvictsCatalogButPreservesPriceCost|TestAssistedSankhyaConfirmInvalidatesCatalogAfterSuccessfulPersistence' -count=1 -v`
  - Pass. Both composed flows and the direct success-path invalidation assertion passed.

- `go test ./... -count=1`
  - Pass. All packages reported `ok` or `[no test files]`, including cache,
    inventory, orders, profitability, integration, and unit packages; command
    runtime was 59s.

- `go vet ./...`
  - Pass. No output or diagnostics.

- Race validation: not run. This environment has `CGO_ENABLED=0` and no C
  compiler, so `-race` is unavailable. No race coverage is claimed.

## Evidence Mapped to M-04 Criteria

| Criterion | Evidence actually run | Result |
| --- | --- | --- |
| M-04-C01 | `TestFreshnessCacheTTLPerClass` in cache package and FreshnessCache suite; its comma-ok lookup proves the cached unknown cost key is present with an explicit nil value | Pass |
| M-04-C02 | `TestFreshnessCacheSingleflight` is retained as a weak guard because it lacks an all-waiters barrier; `TestFreshnessCacheErrorNotCached` is the real mutation guard for singleflight removal (20/20), plus `TestFreshnessCacheWaiterContextCancellation` | Pass |
| M-04-C03 | `TestComposedCatalogHTTPNoCacheBypassesAndRepopulates` proves the real handler→cache→Oracle flow, including response-body `as_of` repopulation; the older cache-only bypass-repopulation assertion is insensitive. Linkage exclusion remains mutation-sensitive in `TestFreshnessCacheBypassAndLinkageExclusion`, with `TestFreshnessCacheOlderLoadCannotOverwriteBypassRefresh` and `TestFreshnessCacheTTLPerClass` retained | Pass |
| M-04-C04 | `TestFreshnessCacheLRUAndLogs` proves structured fields are present and every `internal_read.cache` record has exactly `{cache, key_class}` plus standard `time`, `level`, and `msg`; its temporary raw-`key` mutation failed. `TestFreshnessCacheDeepCopiesCachedFacts` and the full cache suite also passed. The C01 comma-ok assertion proves nil map-value preservation. | Pass |
| M-04-C05 / R1 | `TestComposedLinkageConfirmEvictsCatalogButPreservesPriceCost` proves the real handler→cache→successful `Confirm` flow, catalog refresh, newer response `as_of`, and warm pricecost survival. `TestAssistedSankhyaConfirmInvalidatesCatalogAfterSuccessfulPersistence` requires exactly `[catalog]`. Existing generation-fence tests and `TestEvictOnMutation` remain. Failed-write/no-invalidation coverage is provided by `TestEvictOnMutationFailedWriteDoesNotInvalidate`, `TestApplyManualStockActionDoesNotInvalidateFailedOrRejectedWrites`, `TestAssistedSankhyaConfirmDoesNotInvalidateFailedPersistence`, and `TestImportMarginInputsDoesNotInvalidateFailedPersistence`. | Pass for runnable validation; the C03 bypass and both C05 success-path mutations failed as expected and restored tests passed; race evidence unavailable |

## Risks and Remaining Blockers

- Race validation remains unavailable because the environment has no C
  compiler. This is an external validation limitation, not a reported test
  failure or a production code change in this round.
- Cache invalidation remains intentionally per-process; no cross-process bus
  was added in this correction.
- Fixed-SHA re-review and milestone QA remain the Milestone Orchestrator's
  responsibility.

## Handoff

- Status: `correction_validation_evidence_collected`
- R1 correction test and all runnable required validation: passed
- Remaining blocker: cache `-race` requires a C compiler unavailable in this environment
- Next owner: Milestone Orchestrator / fixed-SHA review and QA
