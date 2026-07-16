# F-01-aggregate-sync-endpoints — validation

```yaml
id: F-01
type: feature-validation
status: complete
owner: CHIP-SAT
parent: M-05
updated: 2026-07-16
branch: chip/sat-m05f01-m06f02
governance_anchor: a49168e641ffd6f61932ca57c29b1d1bdcde2fb0
```

## Lock-exception D-02 (hub ruling 2026-07-16)

Legacy `listMarketplaceOrders` 200 response retyped to `ListOrdersResponse` (contract follows
committed runtime — after R1 the transport serializes `orderPageEnvelope`; the legacy
`MarketplaceOrder` required fields `installation_id`/`fetched_at` were no longer emitted
anywhere). `MarketplaceOrder` schema retained with `deprecated: true`; legacy SDK method
retained `@deprecated` (live consumer found by grep: `packages/feature-orders/src/OrdersPage.tsx`
— hub's original scan covered only `apps/web/src`); `importMarketplaceOrders`/POST untouched;
zero `apps/web` edits. Citation also present in the OpenAPI operation description (ce12fd7f)
and dispatch-ledger rows 42b–43b.

## Ladder

### L0 — build + static + contract (COMPLETE)

| Check | Result | Evidence |
|---|---|---|
| `go build ./...` (apps/server_core) | exit 0 | run 2026-07-16 |
| Governance (`harness.ps1 governance -BaseSha a49168e6…`) | status=passed | clean temp worktree `mpc-gov-check` (hub-checkout scan of untracked `.gomodcache/` produces RCFG_DYNAMIC_READER_UNBOUNDED false alarms — profile-known); only baseline `temporary_exceptions` reported; dashboard module + `/sync` prefix registered (17cce1a9) |
| SDK `tsc --noEmit` | exit 0 | sdk-runtime |
| SDK vitest | 46/46 PASS | includes null-fixture test for nullable `last_sync_at` |
| apps/web `tsc --noEmit` | TS2688 'node' — PRE-EXISTING baseline | recorded MIS-002 mission.md:195 + M-05 validation-result.md:66; `vite build` is the binding web gate |
| apps/web vitest | 11/11 PASS | provisional run 2026-07-16 |

### L1 — unit + integration lanes

Canonical lanes were initially blocked by a pre-existing pwsh 7.6.2 tooling bug:
`scripts/harness/Postgres.psm1:4` nested `Import-Module Execution.psm1 -Force` strips
`New-HarnessProcessRequest` from harness.ps1 scope → `harness:unit`/`harness:integration`
failed with "term not recognized" on BOTH chip worktree and hub checkout (probe:
scratchpad/probe-import.ps1). Hub RULING granted the 1-line fix on this branch
(commit 380c7faa; identical fix landing on main — merge auto-resolves).

| Check | Result | Evidence |
|---|---|---|
| `npm run harness:unit` (canonical, post-fix) | **status=passed**, exit 0 | run_id 81effc7521cd4bfe918f4e3b497183a7 (`scripts/.runs/…`); go unit ok + web vitest 11/11 |
| `go test ./...` (no tags, full module, provisional) | exit 0, all pkgs ok | tasks baz7h28hi.output |
| Migrations vs live pg16 | 37 applied, second pass 0 (idempotent) | session pg 51700, per-run DB `mpc_test_<32hex>` |
| `npm run harness:integration` (canonical, post-corrections) | sole failure = allowlisted `TestPhase1SmokeFlow` (`PRICING_INVALID_PRODUCT_ID`); ALL F-01 integration tests PASS live | tasks bklnq3o6p.output — failure_token set contains ONLY package=tests/integration + test=TestPhase1SmokeFlow; migrations 37/0 idempotent, resource_count=0, session pg 51700 |

### Corrections from first live integration run (ledger rows 44–46)

1. `aggregate_sync_read_test.go` cleanup DELETEd from append-only `orders_sankhya_linkage_events`
   (trigger 23514) then FK-blocked on `orders_marketplace_orders` — all subtests had PASSED;
   only cleanup failed. Fixed: cleanup removed (tenants unique-per-run; ledger is append-only
   by design).
2. `operation_run_read_integration_test.go` + `operation_run_repo_integration_test.go` harnesses
   inserted the same `installation_id` for two tenants — impossible: `integration_installations`
   PK is `(installation_id)` alone (23505). Fixed: one installation per tenant + strengthened
   cross-tenant probe (repo scoped to tenant A querying tenant B's installation returns zero rows,
   proving the tenant predicate through the FK-consistent schema).
3. `operation_run_repo_integration_test.go` `TestLatestRuns*` asserted `.Equal(now)` on
   round-tripped timestamps: Go `time.Now()` on Windows carries 100ns resolution, postgres
   timestamptz stores µs → deterministic fail when `now` not µs-aligned (reachable only after
   fix 2). Fixed: `Truncate(time.Microsecond)` on the two fixture clocks — same storage
   convention production uses (`orders/domain/sankhya_linkage.go:178` truncates both sides);
   5/5 live rounds pass post-fix.

All three corrections delta-reviewed ACCEPT (ledger rows 47, 49b/49c); commits 8d8fa9a7 +
1c4e464c.

### FINDINGS for hub (beyond my seams)

- **F-1 (tooling — RESOLVED via hub grant)**: pwsh 7.6.2 nested `Import-Module -Force`
  removes the re-imported module's functions from the caller scope; fix = drop `-Force` in
  `Postgres.psm1:4` (commit 380c7faa on this branch; hub landing identical fix on main).
- **F-2 (pre-existing flake)**: `TestOrderRepositoryDuplicateIdentityGroupPreservesIDSetAndRemainsAmbiguous`
  (`orders/adapters/postgres/order_repo_test.go:172`, untouched since M-06) asserts the two
  original sorted `mpc_line_id`s stay in positions 0–1 after a third is added — but
  `NewMPCLineID` is crypto/rand, so the new ID sorts last with P=1/3; the test fails ~2/3 of
  live runs. Failed in provisional run (`mpl_cbfb…` sorted mid). Fix suggestion: assert
  set-containment (both original IDs present, length 3) instead of positions. Outside my
  dispatch scope — not edited.
- **F-3 (allowlist confirmed)**: `TestPhase1SmokeFlow` failed with `PRICING_INVALID_PRODUCT_ID`
  — matches the known-failure allowlist; cited, not re-proved. Sole failure in the final
  canonical lane run.
- **F-4 (profile candidate)**: Windows Go `time.Now()` has 100ns resolution; postgres
  timestamptz truncates to µs. Any integration fixture asserting `.Equal()` on a round-tripped
  timestamp must `Truncate(time.Microsecond)` (production convention:
  `orders/domain/sankhya_linkage.go:178`). Failure signature: got/want differ only in the
  7th fractional digit.
- **F-5 (tooling, ratified with hub)**: pwsh 7.6.2 nested `Import-Module -Force` scope removal
  — fix on main @ 89892b5e + this branch @ 380c7faa.

### L2 — live QA

Not chip-owned (per-feature grain closes at L1 + evidence; milestone close = hub dual gate + QA).
