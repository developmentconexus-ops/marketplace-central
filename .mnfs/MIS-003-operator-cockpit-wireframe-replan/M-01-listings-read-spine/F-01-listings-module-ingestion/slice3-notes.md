# Slice 3 integration verification note

The prior registration blocker is resolved: `scripts/harness/Postgres.psm1` includes
`./internal/modules/listings/adapters/postgres` in the integration test arguments.

The registered command was run twice:

```text
pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command integration
```

Both runs started an ephemeral PostgreSQL container and allocated a host port, but stopped
before tagged tests during the first migration invocation. The harness-reported result was:

```text
target=ephemeral-postgres
key=MPC_TEST_DATABASE_URL
migrations=embedded
migrations_first=-1
migrations_second=-1
resource_count=0
status=blocked
postgres lifecycle failed reasons=HPG_MIGRATION_FAILED exit_code=1
```

The tagged listings package compiles independently with:

```text
go test -tags=integration ./internal/modules/listings/adapters/postgres -run '^$' -count=1
ok marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres [no tests to run]
```

Classification: verification-lane migration failure. No harness or migration files were
changed in Slice 3.

An isolated ephemeral PostgreSQL run using the same migration runner and generated-database
naming contract succeeded, and then executed the tagged listings repository test:

```text
applied 36 migration(s)
ok marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres 2.935s
```

This proves the migration set and Slice 3 integration test pass against real PostgreSQL, while
the required registered lane itself remains honestly reported as failed before test execution.
