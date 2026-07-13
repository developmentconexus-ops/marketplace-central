# F-09 Validation Evidence

## Status

Implementation and focused fake-target proof pass. The approved
ephemeral-Postgres lane applied all 33 migrations and proved a zero-migration
second pass, but its lifecycle stopped before Go integration tests because the
separately owned harness still hardcodes an expected inventory of 32. Postgres
integration execution is therefore **blocked, not passed**.

The single focused review fixup replaces destructive order-item refresh with
identity-based in-place upsert, adds explicit linked-line removal conflict, and
persists the validated account-scoped external order key. No harness path was
changed.

## Evidence

### Context compile and validation

- Result: passed.
- Base SHA: `02538fe05307cec91bd445afac92ae4474d68338`.
- Artifact: `context.json`.
- The schema-backed feature work contract resolves to only the dispatched
  paths, selectors, migration seam, side effects, and registered command IDs.

### go-orders-linkage-unit

Command from `apps/server_core` with repository-local `GOCACHE`:

```text
go test ./internal/modules/orders/... -count=1
```

Result: passed. This covers the domain reconciliation and linkage validation
tests and compiles every orders package. Focused domain proof includes unique
provider-identity retention across mutable quantity/price, preserved duplicate
ID multisets with excess-only allocation, explicit legacy/ambiguous state, and
semantic linkage retry validation.

### go-migration-embed

Command from `apps/server_core` with repository-local `GOCACHE`:

```text
go test ./internal/platform/migrate -count=1
```

Result: passed. The canonical embedded inventory contains migration 0033 and
the foreign-working-directory invariant remains green.

### go-orders-linkage-postgres compile proof

Command from `apps/server_core` with repository-local `GOCACHE`:

```text
go test -tags integration ./internal/modules/orders/adapters/postgres -run '^$' -count=1
```

Result: passed compilation; no integration tests were executed by this command.
The compiled tests cover refresh reconciliation, atomic append, semantic and
concurrent retry, conflicting-key/origin rejection, line ownership,
tenant/installation scoping, and append-only database enforcement.
They additionally cover refreshing a confirmed linked order without recreating
its rows, atomic rollback with `ErrOrderLineIdentityConflict` when a linked line
is removed, exact external-order-key load/retry semantics, and conflicting
external-key/idempotency intent.

### go-orders-linkage-postgres approved target

Registered command:

```text
./scripts/harness.ps1 -Command integration
```

Observed bounded output:

```text
migrations_first=33
migrations_second=0
resource_count=0
status=blocked
postgres lifecycle failed reasons=HPG_MIGRATION_NOT_IDEMPOTENT exit_code=1
```

Classification: verification ownership conflict. The migration was accepted by
Postgres and the second runner pass applied nothing. The lifecycle then compared
the first count to a hardcoded `32` in `scripts/harness/Postgres.psm1` and
stopped before Go integration execution. That harness path is forbidden and
owned by a separate writer, so F-09 did not patch or bypass it. A distinct
harness-count compatibility slice must update the registered lane and rerun
this proof.

### git-diff-check

- `git diff --check`: passed; only line-ending conversion warnings were emitted.
- Scoped status inspection: all F-09 changes are within the dispatch allowlist.
- Pre-existing untracked `docs/research/**` files remain untouched and excluded
  from this feature commit.

## Criteria disposition

- F09-AC01: focused domain proof passed; repository execution awaits the
  harness-count compatibility correction. Review coverage now includes
  in-place linked refresh and explicit removed-linked-line conflict.
- F09-AC02: passed by domain reconciliation tests; line number, quantity, and
  unit price are not identity-key fields.
- F09-AC03: migration applied and embed proof passed; Postgres assertion awaits
  the registered integration rerun.
- F09-AC04: domain semantic proof and integration compilation passed;
  transactional/concurrent Postgres execution awaits rerun. The aggregate now
  validates and compares the exact account-scoped external order key.
- F09-AC05: migration applied; constraint/append-only runtime assertions await
  rerun.
- F09-AC06: fake-target, embed, side-effect, and scoped-diff proof passed; no
  Oracle, provider, production, secret, or PII access occurred.

## Blocker and next

Blocker: registered ephemeral integration is incompatible with the new
canonical migration count. Next: the Milestone Orchestrator dispatches the
separately owned harness-count compatibility slice, then reruns
`./scripts/harness.ps1 -Command integration` at the F-09 commit SHA before QA.
