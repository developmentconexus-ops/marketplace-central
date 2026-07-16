# Slice 11 — L0/L1 deterministic pre-pass report (§7)

```yaml
base_sha: 7f5a1b8c (slice 10)
candidate: working-tree diff — 2 files, 410 diff lines (+ a 3-line comment fix by the milestone owner)
scope: apps/server_core/internal/modules/listings/application/ ONLY
gocache: absolute (worktree .gocache)
buildvcs: false (sandbox Git safe.directory VCS-stamp; no repo config changed)
run_by: milestone-owner (chip), INDEPENDENT of the worker — not a relayed claim
```

## L0 (precedes review dispatch)

- `go build -buildvcs=false ./...` → **exit 0**
- `go vet ./...` → **exit 0** (whole repo)

Re-run green after the milestone owner's post-review comment fix.

## L1 unit (ran ∥ review per §15)

`go test -count=1 ./internal/modules/listings/... ./internal/composition/...` → **exit 0**
(connectors, integrations, application, domain, ports, transport, composition all `ok`)

## L1 integration lane — GREEN

Ephemeral `postgres:16-bookworm`, container `mpc-s11c-…`, port 55531:

- migrate #1 → **applied 37**, #2 → **applied 0** (idempotent)
- `listings_status_check` = the 8-value set (slice 9 intact):
  `active, paused, closed, unknown, under_review, inactive, payment_required, not_yet_active`
- `go test -tags=integration -run TestListingsRead -count=1 ./tests/integration` → **PASS, exit 0** (21.4s)
  - `TestListingsReadContractEndToEnd` **8/8 subtests**, including
    **`null_cost_honesty_known_margin_and_summary`** (C07 null-honesty NOT regressed by G2's revert)
    and `tenant_isolation_all_read_paths_and_cursors`.
  - `TestListingsReadPerformance2000` PASS — nearest-rank **p95 3.56ms**, index-only scan, summary
    aggregate query count=1.

### Lane infra note (NOT a code signal — read this before believing a red lane)

Two earlier lane attempts went red and **both were pure infra, caused by the milestone owner's
scripting, not by slice 11**:

1. **Attempt 1** — `sed`-derived script inherited `docker run --rm`; the container died and `--rm`
   destroyed the body before autopsy. Output read `CREATE_DB_OK=False after 240 tries` then
   `FAIL TestListingsReadContractEndToEnd` / `unexpected EOF`. Port 55511 showed `TimeWait`
   connections, proving the container had accepted connections and then died.
2. **Attempt 2** — script rewritten from memory invented the env-var contract: used
   `TEST_DATABASE_URL` + `go run ./cmd/migrate up` instead of the real `MPC_TEST_DATABASE_URL` +
   `go run ./cmd/testdb migrate`. Output read `MC_DATABASE_URL is required` and
   `HPG_TARGET_INVALID` → `FAIL TestListingsReadPerformance2000`.

**Every one of those failure modes prints output that looks like a code failure.** Neither was.
The lesson is recorded: start from the lane script that worked and change only name/port; never
`sed` it blindly (inherits `--rm`) and never rewrite it from memory (invents plausible-but-wrong
env vars). Attempt 3 (`s11-lane3.ps1`) added: no `--rm`, echoed `docker run` exit code, 300-retry
cap, and an autopsy branch (`docker ps -a` + `docker logs`, exit 9) that aborts loudly as INFRA
instead of running tests against a database that does not exist. `CREATE_DB_OK=True after 192 tries`
— initdb duration varies wildly with host load (12 / 61 / 192 / 229 tries observed on this box).

## Scope

Exactly 2 files: `internal/modules/listings/application/read_service.go` and `read_service_test.go`.
No `internal_read/` (the hub-owned cross-module seam), no composition, no transport, no migration,
no OpenAPI, no SDK, no `docker/dev/*.sh`, no `.env`. Matches the slice-11 brief exactly.

## Milestone-owner independent verification (not relayed from the worker)

- `read_service.go:338` — `Get`'s matrix guard reads `if ceilingErr != nil` (**G2 reverted**).
- `read_service.go:23-25` — `errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)`
  returns false BEFORE `IsReadErrorCode` (**G1 in-module half**, per hub ruling (C)).
- Whole-file grep for narration/PR-voice (`hub ruling`, `pins B1/B2`, `supersed`, `task #2`,
  `round-2`, `ADR-17`) across both files → **zero hits** (**G4** genuinely clean, not just in the
  diff hunks).
- `internal_read/` untouched — confirmed by `git status`: only the 2 listings files are modified.
