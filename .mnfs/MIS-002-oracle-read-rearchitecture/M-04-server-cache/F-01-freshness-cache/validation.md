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

The second and final correction round is implemented on top of `0e1c2469`.
The three returned defects have focused regression tests and pass the targeted
and full non-race validation below. This document records feature evidence;
it is not a milestone or formal QA verdict.

## Corrections

### This round

1. M-04-C05 stale delivery after invalidation: `Cache.load` captures the class
   generation once before lookup and includes `(class, generation, namespace,
   key)` in the singleflight key. The same captured generation is passed to
   `storeIfCurrent`, so post-invalidation readers start a fresh downstream
   load while the existing generation fence still drops stale stores.
   `TestFreshnessCachePostInvalidationDoesNotJoinInFlightLoad` verifies two
   downstream calls and a strictly newer `AsOf` for the post-mutation reader.
2. M-04-C03 older-load repopulation: `storeIfCurrent` now drops an incoming
   snapshot only when the existing snapshot is strictly newer, while holding
   the same mutex as the generation check. Dropped writes do not clone, renew
   `created`, or move the LRU entry; equal snapshots still replace normally.
   `TestFreshnessCacheOlderLoadCannotOverwriteBypassRefresh` verifies that a
   slower normal load cannot overwrite a newer MaxAge=0 refresh.
3. M-04-C05 inventory evidence: the `memoryActionStore` test fake now fails
   only Save #2. The persistence subtest reaches the approved Save, performs
   an Applied provider write, reaches the final Applied Save, and fails there;
   it asserts one provider call, two saves, and no invalidation. The sibling
   provider-rejection subtest remains genuine and unchanged in intent.

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

- `apps/server_core/internal/modules/internal_read/adapters/cache/cache.go`
- `apps/server_core/internal/modules/internal_read/adapters/cache/cache_test.go`
- `apps/server_core/internal/modules/inventory/application/stock_action_service_test.go`
- `.mnfs/MIS-002-oracle-read-rearchitecture/M-04-server-cache/F-01-freshness-cache/validation.md`

## Commands Run

All Go commands used this absolute cache path:

```powershell
$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\focused-borg-9c9811\.gocache'
```

- `go test ./internal/modules/internal_read/adapters/cache/ -count=1 -v`
  - Pass. Excerpts:

    ```text
    --- PASS: TestFreshnessCachePostInvalidationDoesNotJoinInFlightLoad
    --- PASS: TestFreshnessCacheOlderLoadCannotOverwriteBypassRefresh
    --- PASS: TestFreshnessCacheSingleflight
    --- PASS: TestFreshnessCacheWaiterContextCancellation
    PASS
    ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/cache 0.908s
    ```

- `go test ./... -run FreshnessCache -count=1 -v`
  - Pass. The cache package passed the TTL, singleflight, invalidation,
    older-load, bypass, cancellation, deep-copy, error, linkage, and LRU/log
    FreshnessCache tests.

- `go test ./... -run EvictOnMutation -count=1 -v`
  - Pass. Excerpts:

    ```text
    --- PASS: TestEvictOnMutationFailedWriteDoesNotInvalidate
    --- PASS: TestEvictOnMutation
    PASS
    ```

- `go test ./internal/modules/inventory/... -count=1 -v`
  - Pass. Excerpts:

    ```text
    --- PASS: TestApplyManualStockActionDoesNotInvalidateFailedOrRejectedWrites/persistence_failure
    --- PASS: TestApplyManualStockActionDoesNotInvalidateFailedOrRejectedWrites/provider_rejection
    PASS
    ok   marketplace-central/apps/server_core/internal/modules/inventory/application 9.483s
    ```

- `go test ./... -count=1`
  - Pass. All packages reported `ok` or `[no test files]`, including cache,
    inventory, orders, profitability, integration, and unit packages.

- `go vet ./...`
  - Pass. No output or diagnostics.

- `go test ./internal/modules/internal_read/adapters/cache/ -race -count=1`
  - Not runnable in this worktree environment. First attempt reported
    `-race requires cgo; enable cgo by setting CGO_ENABLED=1`; the retry with
    `CGO_ENABLED=1` reported `C compiler "gcc" not found`. No C compiler was
    available as `gcc`, `clang`, or `cl`, and no dependency/toolchain install
    was performed.

## Evidence Mapped to M-04 Criteria

| Criterion | Evidence actually run | Result |
| --- | --- | --- |
| M-04-C01 | `TestFreshnessCacheTTLPerClass` in cache package and FreshnessCache suite | Pass |
| M-04-C02 | `TestFreshnessCacheSingleflight` (20 live waiters, one downstream call), `TestFreshnessCacheErrorNotCached`, and `TestFreshnessCacheWaiterContextCancellation` | Pass |
| M-04-C03 | `TestFreshnessCacheOlderLoadCannotOverwriteBypassRefresh`, existing bypass namespace/repopulation tests, and `TestFreshnessCacheTTLPerClass` | Pass |
| M-04-C04 | `TestFreshnessCacheLRUAndLogs`, `TestFreshnessCacheDeepCopiesCachedFacts`, and full cache suite | Pass |
| M-04-C05 | `TestFreshnessCacheFencesInFlightLoadAfterInvalidation`, `TestFreshnessCachePostInvalidationDoesNotJoinInFlightLoad`, `TestEvictOnMutation`, inventory final-Save failure, inventory provider rejection, and full suite | Pass for runnable validation; race evidence unavailable due host toolchain |

## Risks and Remaining Blockers

- The required cache race command could not execute because the environment
  has no C compiler. This is an external validation blocker, not a reported
  test failure or a code change.
- Cache invalidation remains intentionally per-process; no cross-process bus
  was added in this correction.
- Fixed-SHA re-review and milestone QA remain the Milestone Orchestrator's
  responsibility.

## Handoff

- Status: `correction_validation_evidence_collected`
- Three correction tests and all runnable required validation: passed
- Remaining blocker: cache `-race` requires a C compiler unavailable in this environment
- Next owner: Milestone Orchestrator / fixed-SHA review and QA
