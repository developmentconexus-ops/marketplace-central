# F-10 Postgres Harness Migration Inventory

```yaml
id: F-10
type: feature
status: planned
owner: Feature Implementer
parent: M-06
```

## Outcome

Make the registered ephemeral-Postgres lifecycle derive its expected first-run
migration count from the canonical `apps/server_core/migrations/*.sql`
inventory, then rerun the F-09 integration proof.

## Brief

Replace stale migration-count constants in the owned Postgres harness and fake
fixtures with one fail-closed count derived from the canonical SQL files, then
rerun the registered integration lifecycle without widening application or
migration scope.

## Acceptance criteria

1. The lifecycle no longer hardcodes 32 or 33; it fails closed when canonical
   migration inventory is absent/invalid and requires first-run count to equal
   the derived inventory and second-run count to equal zero.
2. Fake lifecycle fixtures and assertions derive the same canonical count.
3. Focused harness contract/lifecycle tests pass.
4. `./scripts/harness.ps1 -Command integration` reaches and passes the Go
   integration tests, including the F-09 orders linkage repository tests, or
   returns one exact unrelated blocker without bypass.
5. No application, migration, Oracle/provider, or production-data path changes.

## Expected Output

- Run specs carry a positive canonical SQL-file count and fail before side
  effects when the inventory is missing or empty.
- First-run migration output must equal the derived count; second-run output
  must equal zero.
- Fake and Docker assertions contain no 32/33 inventory constant.
- Focused fake proof and the registered F-09 Postgres integration rerun are
  recorded in this feature directory.
