# F-10 Specification — Postgres Harness Migration Inventory

## Outcome

The ephemeral-Postgres lifecycle derives its expected first-run migration count
from the repository's canonical `apps/server_core/migrations/*.sql` files. It
accepts exactly that count on the first run and zero on the second run, without
embedding the current inventory size in harness source or fixtures.

## Contract

- Owner: this Feature Implementer owns the dispatched `harness-control` paths.
- Port: `New-HarnessPostgresRunSpec` resolves and carries the canonical count;
  `Invoke-HarnessPostgresLifecycle` consumes it for the idempotency gate.
- Consumers: fake lifecycle tests, the fake Go probe, and the Docker cleanup
  integration assertion.
- Canonical source: direct SQL files under
  `<repository-root>/apps/server_core/migrations`.
- Invalid/unknown state: a missing directory, unreadable inventory, or zero SQL
  files raises `HPG_MIGRATION_INVENTORY_INVALID` before Docker/database work.
- Legacy decision: remove every harness/fixture assertion of 32 or 33; do not
  preserve a fallback count.

## Acceptance proof

### F10-AC01 — Canonical lifecycle count

The lifecycle rejects an absent/empty canonical inventory and otherwise
requires the first migration run to equal the resolved SQL-file count and the
second run to equal zero.

### F10-AC02 — Derived fake inventory

The fake probe and lifecycle assertions derive the inventory from repository
truth and contain no current-count constant.

### F10-AC03 — Focused harness proof

The fake lifecycle suite passes, including missing/empty inventory cases.

### F10-AC04 — F-09 integration rerun

The registered integration lifecycle reaches the Go repository suite, passes
the F-09 orders-linkage tests, and proves owned database/container cleanup.

- `postgres-harness-fake`: focused fake lifecycle, including invalid inventory.
- `go-orders-linkage-postgres`: registered `integration` harness command reaches
  and passes the F-09 repository tests and proves owned cleanup.
- `git-diff-check`: only dispatched harness/test and Feature evidence paths
  change; no application or migration file changes.

## Boundaries

No application, migration, Oracle, provider, production-data, dependency, API,
SDK, or governance changes. Unknown migration inventory never becomes zero or
a historical default.
