# F-03 Ephemeral PostgreSQL and Canonical Migrations — Specification

```yaml
id: F-03
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Problem

The repository's PostgreSQL tests read ambient `MC_DATABASE_URL`, skip when it
is absent, and resolve migrations from caller-dependent filesystem paths. A
database named `mpc_test_*` can still be remote or persistent; a prefix alone
does not prove ownership. The dev Compose PostgreSQL uses a fixed port and a
persistent volume, so it cannot be the integration target.

F-03 must own a disposable PostgreSQL 16 container and every mutable resource
for one run. It must run the real migrations and selected repository tests,
clean up through failure, and prove the persistent dev database stayed
unchanged without letting the ordinary integration lane contact it.

## Requirements

### Owned Docker lifecycle

- Normal `harness:integration` accepts no caller database URL. Retain
  `-DatabaseUrl` for command compatibility, but reject any nonblank value with
  `HPG_EXTERNAL_TARGET_FORBIDDEN` before Docker, run-directory creation, or
  process start.
- Generate one 32-lowercase-hex run ID, database
  `mpc_test_<run-id>`, container `mpc-pg-<run-id>`, and run label. Validate the
  exact relation before the first Docker call.
- Use local `postgres:16-bookworm` with `--pull=never`, `--rm`, no Compose
  project, no `container_name`, no named/host volume, a bounded `tmpfs` at
  `/var/lib/postgresql/data`, and dynamic `127.0.0.1::5432` publication.
- Pass the generated password through the Docker process environment with
  key-only `--env POSTGRES_PASSWORD`; never place it or a database URL in argv,
  output, committed evidence, or summary.
- Wait for `pg_isready`, inspect the assigned loopback port, create the exact
  test database from maintenance database `postgres`, and validate the
  resulting URL before Go receives it.
- Cleanup always connects to `postgres`, forces target connections closed,
  executes `DROP DATABASE ... WITH (FORCE)`, verifies absence, and force-removes
  the run container. Preserve the primary failure and cleanup failures; any
  cleanup/resource leak keeps the result nonzero.
- Post-cleanup Docker inventory by exact label/name must be empty. Never stop,
  inspect through credentials, or mutate the dev Compose project.

### Canonical migrations

- `apps/server_core/migrations/*.sql` is canonical. There are currently 32
  unique filenames; two intentionally share numeric prefix `0021`. Never rename
  an already-applied file or infer identity from numeric prefix.
- Embed the SQL files with `go:embed`; the migration runner consumes an `fs.FS`
  source and lexicographic full filenames. Production `cmd/migrate` uses the
  same embedded source.
- Remove `MC_MIGRATIONS_DIR` and its caller-CWD fallback from the command,
  Docker entrypoint/Compose configuration, registry, exceptions, and tests.
- Apply migrations in per-file transactions without accumulating deferred
  rollbacks in a loop. On a fresh database, first run must apply the exact
  canonical filename set; second run must apply zero; `schema_migrations` must
  equal that set.

### Test target and fixtures

- A central test-only Go boundary reads only `MPC_TEST_DATABASE_URL`, validates
  scheme, loopback host, explicit port, and exact generated database name, then
  builds explicit `pgdb.Config`/`pgxpool` values. It fails, never skips, when an
  integration-tagged test lacks or receives an unsafe target.
- Real database tests use `//go:build integration`; ordinary `go test ./...`
  cannot contact or silently skip a database.
- The real lane freezes and reports this package inventory:
  `./tests/integration`, `./internal/modules/orders/adapters/postgres`,
  `./internal/modules/profitability/adapters/postgres`, and
  `./internal/modules/product_links/application`. Missing package execution is
  a failed run, not partial success.
- Refactor F-03-owned direct readers in integration, orders, profitability, and
  product-links tests to the central boundary. Do not set global
  `MC_DATABASE_URL` inside tests.
- Use typed, FK-complete fixtures. Seed only required marketplace/provider
  definitions before installations/credentials. Set required marketplace codes
  explicitly. Missing-parent negative fixtures must observe the real FK error;
  never weaken, defer, or bypass production constraints.
- Replace the mixed product-links Oracle/PostgreSQL live test with deterministic
  ephemeral-PostgreSQL coverage. F-04 owns later Oracle live proof; F-03 cannot
  promote deterministic matcher evidence to live Oracle.

### Hermetic execution and evidence

- `Environment.psm1` treats `reserved_not_ambient` keys as explicit-only. The
  integration child receives generated `MPC_TEST_DATABASE_URL`, repo-local
  Go caches, and network-off Go settings; ambient DB/provider/Oracle values do
  not enter.
- Every Docker/Go subprocess uses the accepted typed executor: absolute tool,
  structured arguments, explicit CWD/environment, bounded timeout, and automatic
  redaction.
- Fake lifecycle fixtures use target `fake` and prove ordering/failures only.
  Real local Docker runs use target `ephemeral-postgres`. Neither is live
  provider/Oracle evidence.
- A validation-only observer may read table counts from the known healthy dev
  Compose PostgreSQL immediately before/after real F-03 runs. It emits only a
  digest of sorted table counts and must be equal. If dev observation is
  unavailable or changes, `M-08-C03` remains unpassed.
- The versioned `dev-invariance` wrapper performs digest-before, invokes at
  least one real `ephemeral-postgres` integration run, then performs
  digest-after/compare. It reports the read-only dev observer and nested
  ephemeral target separately; neither is a live target.
- Existing canonical `.sql` migration files are read-only F-03 inputs. This
  feature may add only the embed source and runner tests; modifying an applied
  SQL file is out of scope and blocks dispatch.

## Stable Reason Codes

`HPG_RUN_ID_INVALID`, `HPG_EXTERNAL_TARGET_FORBIDDEN`, `HPG_DOCKER_MISSING`,
`HPG_DOCKER_UNAVAILABLE`, `HPG_IMAGE_MISSING`, `HPG_RESOURCE_CONFLICT`,
`HPG_CONTAINER_START_FAILED`, `HPG_PORT_UNAVAILABLE`, `HPG_READY_TIMEOUT`,
`HPG_DATABASE_CREATE_FAILED`, `HPG_TARGET_INVALID`, `HPG_MIGRATION_FAILED`,
`HPG_MIGRATION_NOT_IDEMPOTENT`, `HPG_TEST_FAILED`,
`HPG_DATABASE_DROP_FAILED`, `HPG_CONTAINER_REMOVE_FAILED`, and
`HPG_RESOURCE_LEAK`.

Errors expose safe reason, run/container/database identifiers, counts, and
relative artifact paths only.

## Acceptance Criteria

### F03-AC01 — Harness owns an isolated disposable database

- Traces to `M-08-C03`.
- Caller URLs reject before contact; real run uses exact generated identifiers,
  loopback random port, `tmpfs`, local image/no pull, and no dev volume/project.

### F03-AC02 — Canonical migrations are complete and idempotent

- Traces to `M-08-C03`.
- Fresh run applies the exact embedded filename set and second run applies zero.

### F03-AC03 — Real database tests use safe typed fixtures

- Traces to `M-08-C03`.
- Integration-tagged suites pass on the generated target; unsafe/missing target
  fails before pool creation; FK-negative fixtures remain enforced.

### F03-AC04 — Cleanup survives failures and active connections

- Traces to `M-08-C03`.
- Success, child exit `17`, readiness/create/drop failures, and held-connection
  scenarios always attempt force-drop/removal and leave no labelled resource.

### F03-AC05 — Dev state and evidence classification remain honest

- Traces to `M-08-C03`.
- Dev table-count digest is unchanged; secrets/PII are absent; fake and real
  ephemeral evidence remain distinct.

### F03-AC06 — Database secrets and target labels remain safe

- Traces to `M-08-C04`.
- URLs/passwords never reach output/artifacts; fake, ephemeral, and dev-observer
  evidence cannot be promoted to live provider/Oracle evidence.

### F03-AC07 — Governance and context remain current

- Traces to `M-08-C09`.
- Repaired readers/exceptions disappear atomically; governance and the final
  source-hashed F-03 context pack validate from current HEAD.

## Non-Goals

- No provider, Oracle, OAuth, browser, Docker dev-stack, production schema,
  OpenAPI, SDK, or UI behavior change.
- F-04 owns rich manifests, cold image/module provisioning, complete cold gate,
  and Oracle/provider read evidence.
- F-05 owns cross-worktree leases and generalized runtime namespace recovery.
- No migration checksum/concurrency framework without a demonstrated need.

## Stop Conditions

- Docker needs inherited credentials/home/config or an image pull/network
  access: stop and replan the explicit capability boundary.
- Canonical migration inventory differs, first/second counts are inconsistent,
  or a migration requires schema weakening: stop.
- Dev digest changes, any run resource survives, cleanup masks primary failure,
  or a URL/value leaks: block acceptance.
- Unrelated dirty paths or a competing shared-seam writer appears: stop without
  reset, revert, stash, clean, checkout, or restore.

## Handoff

- Current status: `spec_ready`.
- Next owner: Fresh Phase 1 Feature Implementer.
- Next action: Commit adversarial RED fixtures only.
- Evidence required: fake contract RED/GREEN plus real Docker PostgreSQL and
  read-only dev observer evidence.
- Blockers: None; Docker engine and local PostgreSQL image are currently
  available, but implementation must recheck without pulling.
