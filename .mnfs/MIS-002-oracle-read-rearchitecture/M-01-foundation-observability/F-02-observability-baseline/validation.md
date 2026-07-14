# F-02-observability-baseline validation

```yaml
id: F-02
status: blocked
validation_level: QA-2
owner: Feature Implementer
base_sha: 909f61e6
```

## Result

Local implementation and quick validation passed. The governed Oracle lane is still required for the three live baseline numbers, so the feature remains externally blocked until that lane runs.

## Changed paths

- `apps/server_core/internal/modules/internal_read/observability/timing.go`
- `apps/server_core/internal/modules/internal_read/observability/pool_stats.go`
- `apps/server_core/internal/modules/internal_read/observability/config.go`
- `apps/server_core/internal/modules/internal_read/observability/observability_test.go`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/reader_live_test.go`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/reader_live_nocgo_test.go`
- `apps/server_core/internal/composition/root.go`
- `apps/server_core/cmd/server/main.go`
- `.mnfs/MIS-002-oracle-read-rearchitecture/M-01-foundation-observability/F-02-observability-baseline/validation.md`

## Commands and evidence

### Targeted observability tests

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test -count=1 -v ./internal/modules/internal_read/observability`
- Status: Pass
- Evidence type: ran
- Actual: three `oracle_read` lines were asserted and captured: fast call `duration_ms=0 slow_query=false`; forced 1.1-second call `duration_ms=1100 slow_query=true`; error call `duration_ms=0 oracle_code=942`. The captured transcript contained no SQL, bind values, credentials, or raw driver text. A `pool_stats` line contained `open=7 in_use=3 wait_count=11`. `Stop` was called twice successfully.

### Oracle adapter and live-gate tests

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test -count=1 ./internal/modules/internal_read/adapters/oracle`
- Status: Pass
- Evidence type: ran
- Actual: adapter tests passed.

- Command: `Remove-Item Env:MPC_ORACLE_LIVE_TEST -ErrorAction SilentlyContinue; $env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test -count=1 ./internal/modules/internal_read/adapters/oracle -run TestOracleLiveBaseline -v`
- Status: Pass
- Evidence type: ran
- Actual: `TestOracleLiveBaseline` skipped immediately with `set MPC_ORACLE_LIVE_TEST=1 to run live Oracle validation`; no live connection was attempted.

- Command: `$env:CGO_ENABLED='1'; Remove-Item Env:MPC_ORACLE_LIVE_TEST -ErrorAction SilentlyContinue; $env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test -count=1 ./internal/modules/internal_read/adapters/oracle -run TestOracleLiveBaseline -v`
- Status: Blocked
- Evidence type: could-not-run
- Actual: cgo build stopped before the test because `gcc` is not installed. No live connection was attempted.

### Build and full suite

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go build ./...`
- Status: Pass
- Evidence type: ran
- Actual: exit code 0 with an absolute Windows `GOCACHE`.

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./...`
- Status: Pass
- Evidence type: ran
- Actual: all packages passed; the cgo live file was not included by this shell’s default `CGO_ENABLED=0`, while the !cgo gate test passed through the explicit skip path above.

## Scope and safety checks

- Existing domain/application/ports/adapters/transport boundaries were preserved.
- Reader decorators emit only method, duration, slow flag, and safe numeric Oracle code on errors.
- Pool stats use `Database.Stats()` and stop through the existing `ListenAndServe` return path.
- Live baseline SQL is confined to the cgo live test and is read-only; plan failure reports the exact test error rather than fabricating a verdict.
- No unrelated dirty or untracked files were staged.

## Blockers and next

- Blocker: the machine lacks a C compiler for the cgo live test, and no Oracle credentials/live lane are available here.
- The governed lane must run `TestOracleLiveBaseline` with `MPC_ORACLE_LIVE_TEST=1` and record `MPC_BASELINE_TGFPRO_ACTIVE_COUNT`, `MPC_BASELINE_RTT_MS`, and `MPC_BASELINE_PAGE_PLAN`.
- Next owner: Milestone Orchestrator, for fixed-SHA review and governed-lane baseline capture.

## Commit evidence

- Command: `git log --oneline -4`
- Status: Fail for the requested direct-parent criterion; evidence type: ran.
- Actual: `dc5d79de feat(observability): port latency, slow-query and pool-stats logging + live baseline probe` is followed by external commit `777785e1 chore(tooling): remove Claude bridge`, then the supplied base `909f61e6`. Exactly one F-02 commit exists, but it is not directly based on `909f61e6`.
- Current commit parent: `777785e1`; requested parent: `909f61e6`.
