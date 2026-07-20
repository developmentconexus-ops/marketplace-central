# M-01 sync_state + scheduler skeleton — Evidence Pack

```yaml
milestone: M-01
mission: MIS-006-integracao-fundacao
chip_branch: chip-m01-sync-state-scheduler
worktree: .claude/worktrees/chip-m01
base_sha: 078caf12d8a508c387411c1a9734bb42c674c65b   # forked off current main (superset of milestone base 138aac3d)
milestone_base_sha: 138aac3d                          # contract anchor
base_slice_sha: ac8af597e06a56225bf370d49ee9a7ffec0768d4    # green base (dual-gate)
head_sha: 4c9a6494                                          # + corrective slice (delta re-gate)
validation_level: QA-0
status: READY-TO-MERGE (dual-gate AGREEMENT; C3/C4/C11 real-Postgres proof pending hub integration lane)
```

## Worktree isolation correction (reported to hub)

Assigned dir `cool-kalam-218e3c` was NOT a registered git worktree (no `.git` file → git
resolved to the MAIN/hub repo). Initial `git checkout -b`/`reset` briefly flipped the PRIMARY
repo onto the chip branch; detected, RESTORED primary to detached HEAD @078caf12 (untracked
ml-api/mlprobe preserved), then created a PROPER isolated worktree at `.claude/worktrees/chip-m01`
on branch `chip-m01-sync-state-scheduler`. Verified `git rev-parse --show-toplevel` = the chip
worktree and it appears in `git worktree list`. All code authored there. Hub asked to confirm
primary-repo state unaffected.

## Deliverables

| Deliverable | Path | Status |
|---|---|---|
| Migration bloco A (CREATE TABLE sync_state, E8) | `apps/server_core/migrations/0075_sync_sync_state.sql` | done |
| Migration shape test (C1, no DB) | `apps/server_core/migrations/sync_state_test.go` | done, ran green |
| Domain: SyncState + Entity enum + SyncError | `apps/server_core/internal/modules/sync/domain/sync_state.go` | done |
| Postgres repo (read/upsert/atomic-append) | `apps/server_core/internal/modules/sync/adapters/postgres/sync_state_repo.go` | done |
| Repo integration tests (C3/C4/C11/tenancy) | `.../adapters/postgres/sync_state_repo_integration_test.go` | done, compiles; run pending hub lane |
| Scheduler skeleton + RegisterJob seam | `apps/server_core/internal/modules/sync/application/scheduler.go` | done |
| Scheduler unit tests | `.../application/scheduler_test.go` | done, ran green (6 tests) |
| Composition wiring helper | `apps/server_core/internal/modules/sync/composition/scheduler.go` | done |
| root.go ticker-block additive-lock | `apps/server_core/internal/composition/root.go` | done (import + 1 ticker line) |
| modules.json `sync` entry | `contracts/governance/modules.json` | done |
| Migration-count fixture bump 62→63 (profile §3) | `apps/server_core/internal/platform/migrate/runner_test.go` | done, ran green |

### Write-set reconciliation (declared vs changed)

Declared (chip prompt): migration sql · domain · repo · scheduler+test · root.go · modules.json.
Changed-undeclared (each additive, within M-01 ownership axis — justified):
- `migrations/sync_state_test.go` — static shape test of MY table; proves C1 `ran` without DB.
- `sync/composition/scheduler.go` — wiring helper that shrinks the root.go additive-lock to
  1 import + 1 line (reduces collision surface with M-02's root.go ownership); intra-module.
- `sync/adapters/postgres/sync_state_repo_integration_test.go` — real-Postgres proof of MY table.
- `internal/platform/migrate/runner_test.go` — count-fixture bump **mandated by profile §3**
  ("every migration grant implies bumping that fixture in the same slice"). ⚠️ M-02 also adds a
  migration (0076) → M-02 must rebase this count 63→64; flagged to hub for the M-01→M-02 serialize.
Declared-but-unchanged: none.

## Verification ladder

| Lane | Command | Result | Type |
|---|---|---|---|
| L0 build | `GOCACHE=… GOMODCACHE=… go build ./...` | exit 0 | ran |
| L0 vet | `go vet ./internal/modules/sync/... ./internal/composition/...` | exit 0 | ran |
| L1 unit | `go test ./internal/modules/sync/...` | ok (application 6 tests) | ran |
| L1 migration shape | `go test ./migrations/ -run TestSyncStateMigrationShapeE8` | PASS | ran |
| L1 migrate runner | `go test ./internal/platform/migrate/` | ok (count 63) | ran |
| L0 governance | `npm run harness:governance` | HUB-OWNED at acceptance (profile §2: hub re-runs on integrated branch; my `sync` modules.json entry satisfies GOV_MODULE_COVERAGE by inspection) | could-not-run (deferred to hub) |
| L1 integration (C3/C4/C11/tenancy) | `npm run harness:integration` (sync/adapters/postgres) | REQUEST sent to hub (real Postgres; profile §4 hub-serialized for divergent migrations) | could-not-run (blocker: hub-serialized integration lane) |

Fresh-worktree bootstrap ran: `cd apps/server_core && GOMODCACHE=$(pwd)/.gomodcache go mod download all` (exit 0). gomodcache pollution check on `git status`: 0.

## Per-criterion verdicts (M01-C1..C11)

| ID | Verdict | Evidence (file:line) | Type |
|----|---------|----------------------|------|
| M01-C1 | PASS | `migrations/0075_sync_sync_state.sql:19-30` full E8 shape + PK (tenant_id,installation_id,entity); asserted by `migrations/sync_state_test.go` TestSyncStateMigrationShapeE8 | ran |
| M01-C2 | PASS | `internal/composition/root.go:579` new `go synccomposition…Start(ctx)` in ticker block; :575-577 NewRefreshTicker/NewStateCleanup/NewFeeSyncScheduler intact | ran |
| M01-C3 | PENDING-HUB | `sync_state_repo.go:36-73` Read+RecordSuccess (last_full_sync_at set); proof = `sync_state_repo_integration_test.go` TestSyncStateCursorRoundTrip | could-not-run (hub lane) |
| M01-C4 | PENDING-HUB (logic proven) | `sync_state_repo.go:78-95` single-statement `consecutive_failures+1`+last_error; scheduler isolation proven `scheduler_test.go` TestRunOnceFailureIsIsolatedPerEntity (ran); DB proof = TestSyncStateFailureThenRecovery | mixed: unit ran / DB could-not-run |
| M01-C5 | PASS | no import of connectors/mercado_livre/oracle in `internal/modules/sync/*`; job body `composition/scheduler.go` returns cursor unchanged (no external call) | ran (grep diff) |
| M01-C6 | PASS | `scheduler.go` interval is a `time.Duration` ctor param; no "daily"/"diário"/cron literal (asserted negative in both migration + [gate]) | ran (grep) |
| M01-C7 | PASS | `scheduler.go` ctor + `Start(ctx)` loop + `RunOnce` mirror `integrations/background/fee_sync_scheduler.go:31-81` | ran |
| M01-C8 | PASS | `0075_sync_sync_state.sql` pure `CREATE TABLE IF NOT EXISTS`, no ALTER; 0075 < M-02 0076 | ran (diff) |
| M01-C9 | PASS | `contracts/governance/modules.json` sync entry (id/root/code_owner_path/composition_required) | ran |
| M01-C10 | PASS | every statement in `sync_state_repo.go` keys on `tenant_id=$1` (Read:44, RecordSuccess:60, RecordFailure:86, AppendPendingCodigo:110); tenancy DB proof = TestSyncStateTenantIsolation | ran (grep) / DB pending |
| M01-C11 | PENDING-HUB (single-statement proven) | `sync_state_repo.go:100-118` single `INSERT…ON CONFLICT DO UPDATE SET cursor = jsonb_set(…|| to_jsonb($4))` — no app-side read-modify-write; concurrency proof = TestSyncStateConcurrentCursorAppend (50 goroutines) | inspection ran / DB could-not-run |

## Anti-criteria (all must be ABSENT)

| AC | Present? | Evidence |
|----|----------|----------|
| AC-M01-1 (query w/o tenant_id) | ABSENT | all 4 statements tenant-scoped |
| AC-M01-2 (fabricated 0/{} vs NULL) | ABSENT | `jsonbArg` maps empty→NULL; cursor/schedule/last_error columns nullable, no default; only consecutive_failures DEFAULT 0 |
| AC-M01-3 (real ML/Oracle job) | ABSENT | no-op placeholder, no external client |
| AC-M01-4 (root.go outside ticker block) | ABSENT | diff = 1 import + 1 ticker line; source-wiring untouched |
| AC-M01-5 (hardcoded cadence) | ABSENT | interval injected; no cadence literal |

## Design notes (pre-empting gate)

- RecordSuccess treats the returned cursor as authoritative: a nil cursor persists as SQL NULL
  (documented). The M-01 no-op returns the cursor it read, so an existing cursor is never
  clobbered. M-03/M-04 pass their real cursor.
- Read-error handling in `scheduler.go` runJob: a cursor READ error is NOT treated as a first
  run — the entity is skipped this cycle (fail honest, unknown≠zero), next tick retries. Only
  found=false (no error) yields a nil first-run cursor.
- No per-job panic recover: the milestone scopes failure isolation to error-RETURNS; jobs
  return errors, not panic. A blanket recover would trip the anti-slop checklist. Documented.

## P6 Dual-Gate

P6-DUAL-GATE: AGREEMENT

- **Cold reviewer (Claude, gate-reviewer)** — PASS on all C1–C11. Only AC-M01-4 raised as a
  CONCERN, solely because the read-only reviewer had no git-diff access; the chip verified
  independently that the `root.go` change is exactly 2 hunks (1 import + 1 ticker line,
  source-wiring untouched) → AC-M01-4 is a hard ABSENT/PASS.
- **Adversarial refuter (Claude, general-purpose)** — NOT-REFUTED, no BLOCKING. Raised
  IMPORTANT residuals #3 (nil-cursor clobber doc), #7 (no panic isolation), #8 (unvalidated
  entity on the direct repo write path), #9 (C11 needs a real DB run). #3/#7/#8 hardened in a
  corrective slice (below); #9 is the outstanding hub integration-lane REQUEST.
- **Focused delta re-gate (Claude, gate-reviewer)** — on the corrective diff ONLY (per lesson
  chip-import-fix: "do a real 2nd refuter pass on correctives, not self-verified round-2").
  Verdict **NOT-REFUTED**: `Valid()` guards are top-of-method with no preceding side effect
  and no hot-path change; the panic recover is scoped to the job body only (never the
  `store.Read` integrity seam) and the failure branch structurally cannot reach
  `RecordSuccess`, so a panic cannot clobber a stored cursor. It flagged one non-defect
  coverage gap (no test seeding a pre-existing cursor across a panic) → CLOSED by seeding
  `{"page":7}` in `TestRunOncePanicIsIsolatedAsFailure` and asserting products never appears
  in `successes` (`scheduler_test.go`).

**Sol-side waiver disclosure (profile §12):** the adversarial half of the gate is normally a
Sol/GPT-side crew. Sol is on the quota-wall (reset 2026-07-25). Per the pre-authorized §12
contingency, both the refuter and the delta re-gate ran as an independent second Claude crew
(distinct dispatch from the cold reviewer). Disclosed as a waiver, not a true cross-model gate.

### Corrective slice (post base @ac8af597, applied to owned seam only)

| Residual | Fix | File |
|---|---|---|
| #7 panic isolation | `safeInvoke` defer/recover → per-entity `RecordFailure`, loop survives | `application/scheduler.go` |
| #8 unvalidated entity at direct repo seam | `if !entity.Valid() { return ErrUnknownEntity }` guard on Read/RecordSuccess/RecordFailure/AppendPendingCodigo | `adapters/postgres/sync_state_repo.go` |
| #3 nil-cursor clobber | expanded JobFunc doc (nil return ERASES cursor); `TestRunOncePanicIsIsolatedAsFailure` proves panic path never records success (no clobber) | `application/scheduler.go`, `application/scheduler_test.go` |

Corrective delta = 3 files, +79/-1. L0 build/vet exit 0, L1 unit green (7 tests incl. panic
isolation). #9 remains `could-not-run` (hub integration lane).

## Dispatch ledger

| Dispatch | Role | Model | Outcome |
|---|---|---|---|
| gate-reviewer (cold) | independent P6 cold reviewer | Claude | PASS (AC-M01-4 concern = no-diff-access artifact; resolved) |
| general-purpose (refuter) | adversarial refuter (Sol-side rebind, §12 quota-wall waiver) | Claude | NOT-REFUTED; residuals #3/#7/#8 hardened, #9 → hub lane |
| gate-reviewer (delta) | focused re-gate of corrective diff (§12 waiver) | Claude | NOT-REFUTED; flagged coverage gap → CLOSED |

## Live/QA marker

L2/L3 browser QA: N/A — M-01 has NO FE surface and NO HTTP endpoint (sync_state is internal;
modules.json openapi_prefixes=[]). The runnable proof is the L1 integration cycle (C3/C4/C11)
against real Postgres → **could-not-run chip-side; REQUEST sent to hub** (blocker: profile §4
hub-serialized integration lane for divergent migrations M-01/M-02).
