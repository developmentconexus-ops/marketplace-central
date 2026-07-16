# F-01 validation (Slices 1–5)

## Scope and composition

F-01 builds the listings **write/ingestion spine**: the connectors→listings ingestion
pipeline (Slices 1–3), the async `POST /listings/refresh` write endpoint with an atomic
concurrent-run guard and honest operation-run lifecycle (Slice 4), and the end-to-end
contract proof of that endpoint (Slice 5).

Slice 5's `listings_refresh_test.go` drives real HTTP through `httpx.NewRouteClassMux`,
the production `RefreshHandler.Register`, `NewRefreshService`, the real listings-owned
`adapters/integrations` gateway over the published integrations application boundary, the
real integrations Postgres installation + operation-run repositories, and the real
listings Postgres repository over harness migrations (0036). The composition mirrors
`composition/root.go` (ingestion + gateway + refresh service + route registration). The
sole external-data substitution is a deterministic in-process `connectors/ports.ListingReader`
stub at the provider boundary — explicit fake evidence, not live-Mercado-Livre evidence.
The async runner is the production `go task()`, so the endpoint returns `202` before the
run reaches a terminal state; the tests use a bounded 2 s / 10 ms poll of the persisted
operation-run terminal status.

## Contract evidence (Slice 5 — 5 tests)

- **`TestListingsRefreshSeedsIC02RowsAndClosesMissing`** — a single-page pull of six
  snapshots (`MLBTEST0001..0006`) ingests exactly 6 tenant-scoped rows on the composite
  PK `(tenant_id, installation_id, provider_listing_id, variation_id)`, defaults
  `variation_id='-'` where no variation (5 of 6), keeps nullable facts SQL NULL
  (`sales_30d IS NULL` on all 6, ADR-17), and stores an unrecognized provider status as
  `'unknown'` (1 of 6). A second refresh dropping one provider listing closes **only** the
  removed row (`status='closed'`) while the other five retain their mapped facts/statuses.
- **`TestListingsRefreshRejectsConcurrentRunWithActiveID`** — the stub blocks its first
  `ListListings` on a channel gate; the first `POST` returns `202` with operation id A; a
  concurrent second `POST` for the same installation returns `409 refresh_in_progress`
  carrying id A at `error.details.operation_run_id`; exactly one `queued`/`running` run
  exists. Releasing the gate drives run A to `succeeded`. (Proves the Slice-4
  advisory-xact-lock exclusive-start over the real endpoint.)
- **`TestListingsRefreshCapabilityErrorMidPullLeavesRowsUnchanged`** — page size 3; page 1
  returns 3 valid snapshots (forcing a page-2 request), page 2 returns a typed
  `ErrCodeProviderAuth` capability error. The run ends `failed`; the listings rows are
  **byte-identical** before/after (`jsonb_agg` snapshot equality — a mid-pull error commits
  nothing, atomic upsert-close only on a complete pull); `translated_error_code` persists as
  `CONNECTORS_PROVIDER_AUTH`; `provider_evidence_json` contains **no** raw provider body; the
  report callback is **not** invoked (it fires only on persistence failure, not ordinary
  pull failure).
- **`TestListingsRefreshIsTenantIsolated`** — `integration_installations` PK is
  `(installation_id)` (globally unique), so tenant B takes its own distinct installation;
  tenant A refreshing and closing a listing leaves tenant B's listings (6, one closed by its
  own refresh) and operation runs untouched — every assertion carries a `tenant_id`
  predicate.
- **`TestListingsRefreshUnknownInstallation`** — a `POST` for a nonexistent installation
  returns `404 installation_not_found` and creates **no** operation-run or listing rows.

## Slice 4 evidence (endpoint + atomicity)

`TestOperationRunRepositoryBeginExclusiveIsAtomic` (integrations adapters/postgres): two
simultaneous `BeginExclusive` calls (barrier-gated goroutines, same tenant+installation,
different run ids) create exactly one `queued` run; the loser receives the active run; a
different installation and a different tenant each start independently; DB count == 1.
Exclusive-start is one pgx txn —
`pg_advisory_xact_lock(hashtextextended('tenant|installation|listings_refresh',0))` →
SELECT queued/running → return active on hit / INSERT queued on miss → COMMIT (lock
auto-released). No new migration (reuses `integration_operation_runs`, 0016/0021). The
listings lane consumes only the published integrations boundary (grep-confirmed zero
`integrations/adapters/postgres` import — condition c). OpenAPI `refreshListings`
(202/400/404/409) + `sdk-runtime` method shipped in the same commit.

## Lane result (2026-07-15, milestone-owner run)

Reviewer (cavecrew-reviewer, sonnet): **0 findings** — cursor keying correct, tenant-scoped
SQL throughout, gate-based concurrency deterministic (no timing race). Owner verification:
`go build ./...` 0; `go vet -tags=integration ./tests/integration/...` 0.

Hermetic ephemeral-Postgres lane — ephemeral `postgres:16-bookworm` with a
**Docker-assigned host port discovered via `docker port`** (race-proof: avoids the fixed-port
TIME_WAIT collision and the WinNAT reserved-range bind failure that produced two false-alarm
runs), retry-`CREATE DATABASE mpc_test_<32hex>` loop, `go run ./cmd/testdb migrate` (applied
36) → `go test -tags=integration -run TestListingsRefresh -v -count=1 ./tests/integration`:

**PASS** — `ok marketplace-central/apps/server_core/tests/integration 10.643s`:
- `TestListingsRefreshSeedsIC02RowsAndClosesMissing` PASS (1.69s)
- `TestListingsRefreshRejectsConcurrentRunWithActiveID` PASS (1.05s)
- `TestListingsRefreshCapabilityErrorMidPullLeavesRowsUnchanged` PASS (1.51s)
- `TestListingsRefreshIsTenantIsolated` PASS (3.03s)
- `TestListingsRefreshUnknownInstallation` PASS (0.25s)

Real ephemeral Postgres, real migrations (0036), real registrar/handler/gateway/repositories;
the only substitution is the deterministic `connectors/ports.ListingReader` stub at the
provider boundary. Slice-4 atomicity proven separately by
`TestOperationRunRepositoryBeginExclusiveIsAtomic` (PASS 4.51s). Environment: go1.26.x,
windows/amd64, worktree SHA at Slice-5 commit `748a9ff`. Completes **M01-C02**, the
operation-run evidence for **M01-C01**, and advances **M01-C08**.
