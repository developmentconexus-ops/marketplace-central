# Slice 10 — L0/L1 deterministic pre-pass report (§7)

```yaml
base_sha: 1f6b72d8 (c4e8ab91 slice 9 + gate-round-2 evidence)
candidate: working-tree diff (uncommitted) — 3 files, 716 diff lines
gocache: absolute (worktree .gocache)
buildvcs: false (sandbox Git safe.directory VCS-stamp; no repo config changed)
run_by: milestone-owner (chip), INDEPENDENT of the worker (not a relayed claim)
```

## L0 (precedes review dispatch)

- `go build ./...` → **exit 0**
- `go vet ./...` → **exit 0** (whole repo)

## L1 unit (runs ∥ review per §15)

`go test -count=1 ./internal/modules/listings/... ./internal/composition/...` → **exit 0**

connectors, integrations, application, domain, ports, transport, composition all `ok`.

## L1 integration lane — GREEN (regression proof)

No migration touched by slice 10, but `read_service.go` was substantially reworked, so the lane ran anyway
to prove the read spine did not regress end-to-end. Ephemeral `postgres:16-bookworm`:

- migrate #1 → **applied 37**, #2 → **applied 0** (idempotent); `listings_status_check` still the 8-value
  set (slice 9 intact).
- `go test -tags=integration -run TestListingsRead -count=1 ./tests/integration` → **PASS, exit 0** (26.3s):
  - `TestListingsReadContractEndToEnd` PASS — all 8 subtests, including **`null_cost_honesty_known_margin_and_summary`** (C07 null-honesty NOT regressed by the fail-honest change) and `tenant_isolation_all_read_paths_and_cursors`.
  - `TestListingsReadPerformance2000` PASS — nearest-rank p95 3.93ms, index-only scan, summary aggregate query count=1.

Lane infra note (not a code signal): the first two lane attempts reported `CREATE_DB_OK=False` / a 3D000
test failure. Root cause was the postgres first-boot `initdb` exceeding the 90×700ms retry cap under host
load, plus a second attempt racing a `docker stop` from the first script. Re-run on a clean container with
a 240 retry cap: `CREATE_DB_OK=True after 61 tries`, everything green. No code involvement.

## Scope

Exactly 3 files: `internal/composition/root.go`, `internal/modules/listings/application/read_service.go`,
`internal/modules/listings/application/read_service_test.go`. No migration, no OpenAPI, no SDK, no
transport, no slice-9 file, no `docker/dev/*.sh`, no `.env`. Matches hub R2 scope exactly.

## Milestone-owner independent verification (not relayed from the worker)

The worker's load-bearing claim — *"transport's catch-all already maps any error to 503
source_unavailable, no transport change needed"* — **VERIFIED IN-TREE**:
`transport/http_handler.go:68, :113, :214, :274` all `writeListError(… 503, "source_unavailable" …)`,
and `http_handler_test.go:262` pins a generic `errors.New("down")` → `503 source_unavailable`. The hub's
"no OpenAPI/SDK change" condition therefore holds as ruled.

B1 gate read directly in-tree (both Oracle facts covered, not just the pre-check):
- `read_service.go:361-363` — `scan`: ceiling source-unavailable + fact-dependent filter → return error.
- `read_service.go:377-379` — `scan`: cost source-unavailable discovered DURING the scan (threaded out of
  `enrich`) → return error. This is the path a pre-check alone would have missed.
- `read_service.go:383` — `matchesDependentFilter` is genuinely APPLIED when facts are available.
- `scanGroups:215-217` — same semantics for the by-product path.
- `passThrough` / `passThroughGroups` — **0 references** (deleted; unreachable once fixed).
- `needsBelowMarginScan:485` is the uniform trigger (true for `has_exception` AND `below_margin`) — no
  special case, per the hub's rejection of option (A).

Classification mechanism: `internalreaddomain.IsReadErrorCode(err, ReadErrorSourceUnavailable)` — the
EXISTING taxonomy, not a new one. Only that code degrades; `context.Canceled`/`DeadlineExceeded`/adapter
defects propagate.

Telemetry: `slog.Error` on the propagate arm, `slog.Warn` on the degrade arm, both carrying `err` + `op`
+ `fact`, at the per-request fact-fetch sites (Summary/List/ByProduct/Get) — not per row.
