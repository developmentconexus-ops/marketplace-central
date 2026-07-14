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

The third correction round is test-only and addresses a validation defect in
the prior evidence. The R1 post-invalidation singleflight regression test was
initially hollow: it could pass after the generation component was removed
from the normal-path group key. The test now deterministically observes
downstream entry for both readers before releasing the first load. The
mutation test below proves the rewritten test fails when the guarded behavior
is removed and passes after the committed production code is restored. This
document records feature evidence; it is not a milestone or formal QA verdict.

## Corrections

### This round

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
   nil map values on store and return.
4. P2-D waiter cancellation: `DoChan` plus context selection lets canceled
   waiters return without a helper goroutine while live waiters still collapse.
5. P2-E persistence evidence: catalog/inventory/orders/profitability failure
   paths retain evidence that failed persistence does not invalidate; inventory
   also retains provider-rejection coverage.

## Changed Paths

This correction round permanently changes only:

- `apps/server_core/internal/modules/internal_read/adapters/cache/cache_test.go`
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
    --- PASS: TestFreshnessCachePostInvalidationDoesNotJoinInFlightLoad
    --- PASS: TestFreshnessCacheOlderLoadCannotOverwriteBypassRefresh
    --- PASS: TestFreshnessCacheSingleflight
    --- PASS: TestFreshnessCacheWaiterContextCancellation
    PASS
    ok  	marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache	1.613s
    ```

- `go test ./... -count=1`
  - Pass. All packages reported `ok` or `[no test files]`, including cache,
    inventory, orders, profitability, integration, and unit packages; command
    runtime was 59.1s.

- `go vet ./...`
  - Pass. No output or diagnostics.

- Race validation: not run. This environment has `CGO_ENABLED=0` and no C
  compiler, so `-race` is unavailable. No race coverage is claimed.

## Evidence Mapped to M-04 Criteria

| Criterion | Evidence actually run | Result |
| --- | --- | --- |
| M-04-C01 | `TestFreshnessCacheTTLPerClass` in cache package and FreshnessCache suite | Pass |
| M-04-C02 | `TestFreshnessCacheSingleflight` (20 live waiters, one downstream call), `TestFreshnessCacheErrorNotCached`, and `TestFreshnessCacheWaiterContextCancellation` | Pass |
| M-04-C03 | `TestFreshnessCacheOlderLoadCannotOverwriteBypassRefresh`, existing bypass namespace/repopulation tests, and `TestFreshnessCacheTTLPerClass` | Pass |
| M-04-C04 | `TestFreshnessCacheLRUAndLogs`, `TestFreshnessCacheDeepCopiesCachedFacts`, and full cache suite | Pass |
| M-04-C05 / R1 | `TestFreshnessCacheFencesInFlightLoadAfterInvalidation`, the deterministic `TestFreshnessCachePostInvalidationDoesNotJoinInFlightLoad`, `TestEvictOnMutation`, inventory final-Save failure, inventory provider rejection, and full suite | Pass for runnable validation; the R1 mutation failed as expected and the restored test passed; race evidence unavailable |

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
