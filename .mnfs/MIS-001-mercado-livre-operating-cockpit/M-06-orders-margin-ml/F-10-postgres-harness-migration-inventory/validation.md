# F-10 Validation Evidence

## Status

Passed. The Postgres harness now derives a positive first-run expectation from
the canonical `apps/server_core/migrations/*.sql` inventory, rejects missing or
empty inventory before Docker work, and requires a zero-migration second run.
The registered integration reached the Go integration package set, including
the F-09 orders Postgres repository tests, and cleaned every owned resource.

## Context compile and validation

- Base SHA: `a4b54e935cb525bb8b0470ad1bb1fe53065cde59`.
- Final compile: passed with all five dispatched allowed path selectors.
- `context-validate -RequireCurrentBase`: passed.
- Artifact: `context.json`.

## postgres-harness-fake

Command:

```text
pwsh -NoProfile -File scripts/tests/postgres-lifecycle.tests.ps1
```

Result: passed (`PASS postgres lifecycle tests target=fake`). The suite proves
missing and empty canonical inventories return
`HPG_MIGRATION_INVENTORY_INVALID`, the fake probe reads direct `.sql` files
from its controlled `apps/server_core` working directory, first-run output
equals the run spec's derived count, and second-run output equals zero.

## go-orders-linkage-postgres

Registered command:

```text
pwsh -NoProfile -File scripts/harness.ps1 -Command integration
```

Result: passed. Bounded output:

```text
target=ephemeral-postgres
migrations_first=33
migrations_second=0
resource_count=0
status=passed
run_id=5215d376eecb48e39a36458c44e40bc1
```

The registered lifecycle runs Go with `-tags=integration -count=1` across its
approved package set, including
`./internal/modules/orders/adapters/postgres`; therefore the F-09 linkage
repository tests executed against the freshly migrated ephemeral database.
Cleanup reported zero owned resources.

## git-diff-check

- `git diff --check`: passed; only line-ending conversion warnings appeared.
- Scoped inspection found no stale migration-count comparison or fixture
  literal in the four dispatched harness/test files.
- Changed paths are limited to the F-10 artifact directory and the four owned
  harness/test paths.
- Pre-existing untracked `docs/research/**` and `output/` state remained
  untouched and excluded from this feature commit.

## Criteria disposition

- F10-AC01: passed by fail-closed fake cases and both registered lifecycle runs.
- F10-AC02: passed; fake migration output is derived from canonical CWD truth.
- F10-AC03: passed by the focused fake lifecycle suite.
- F10-AC04: passed by the registered integration, including F-09 repository
  execution and zero-resource cleanup.

## Blockers and next

Blockers: none. Next: Milestone Orchestrator reviews the intentional commit and
includes this evidence in the fixed-SHA milestone review/QA flow.
