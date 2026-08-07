# Lane: persistence

> I have run something really deep into MetalDocs and I am changing the way I code there to move to
> something more professional towards issues, PRs, PR review, CodeRabbit mechanical full validation
> and so much more. For that I had to identify every error in my code, my platform, to improve it and
> to create this full validation. I want to run it here as well so we move on the same path, this way
> it gets so much harder to send bad PRs.

Calibration: solid professional level. The operative sentence: **"this way it gets so much harder to
send bad PRs"** — the program's success condition is a mechanism, not a cleaner codebase.

Method: `C:\Users\leandro.theodoro\Documents\MetalDocs\docs\engineering\repo-audit-playbook.md` (cited
by path, not vendored). Established facts 1–14 from `docs/engineering/repo-audit-2026-08-07/PHASE-0.md`
are treated as given, not re-derived; D-44 and D-53 are re-measured, not re-narrated, per instruction.

## Findings

| ID | class | finding | evidence | scale |
|---|---|---|---|---|
| P-1 | gap | RLS exists on only 4 of 64 live tables, and even there it is inert against the connecting role (re-confirms D-44) | `apps/server_core/migrations/0097_catalog_context.sql:80-110` (`ENABLE`+`FORCE`+`CREATE POLICY tenant_isolation ... USING (tenant_id = current_setting('app.tenant_id', true))`); `grep -n "ENABLE ROW LEVEL SECURITY\|CREATE POLICY" migrations/*.sql` → all 12 hits in one file | 4/64 tables (6%); 0 non-catalog tables |
| P-2 | hazard | Tenant scoping is enforced by hand convention, never by machinery — no query builder, no codegen, no linter, no test gate checks for a `tenant_id` predicate | `go.mod` lists only `github.com/jackc/pgx/v5`, no ORM/builder/codegen; `grep -rln tenant scripts/lib scripts/arch-gate.sh scripts/tests` → 0 hits | 246 raw `.Query(ctx\|.QueryRow(ctx\|.Exec(ctx` call sites in non-test Go, 0 automated tenant checks |
| P-3 | hazard/drift | The one place that built real per-tenant machinery (session `set_config` + `FORCE RLS` policy) is itself neutralized by the same root cause as D-44, and it is unique — no other repo/query in the codebase reuses this pattern | `internal/contexts/catalog/internal/postgres/repository.go:49` (`SELECT set_config('app.tenant_id', $1, true)`); `grep -rn "set_config" internal` → 1 hit total | 1 of 61 tenant-scoped tables uses live per-request tenant propagation; the mechanism protects 0 rows in practice because the DSN connects as the bypassing owner (D-44) |
| P-4 | gap | `pgdb.LoadConfig()` still silently defaults an unset tenant to `"tenant_default"` for `cmd/server` — the one guard added (D-39) covers only `cmd/catalogingest`, leaving the entire HTTP server's tenant wiring on the silent path this AGENTS.md rule forbids | `apps/server_core/internal/platform/pgdb/config.go:23-24` (`if cfg.DefaultTenantID == "" { cfg.DefaultTenantID = "tenant_default" }`); consumed at `internal/composition/root.go` (70 occurrences of `cfg.DefaultTenantID`); rule at `AGENTS.md:24` ("unknown operational facts never become zero/default") | 70 repo/service constructions in `root.go` inherit one process-wide tenant value with no fail-closed check; only `cmd/catalogingest` (D-39) is guarded |
| P-5 | idiom/hazard | Money scanned from `numeric(14,2)` columns lands in Go `float64`/`pgtype.Float8` in most adapters; one reader (`pricing/matrix_reader.go`) explicitly avoids this by casting `::text`, so the codebase has two inconsistent, undocumented conventions for the same hazard | `apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go:180,1124,1154` (`scanFloat8`/`nullableFloat8` on `unit_price numeric(14,2)`, migrations/0027:31); `internal/modules/profitability/adapters/postgres/store.go:89,246-254` (10 `pgtype.Float8` fields: revenue/saleFee/cost/tax/freight/commission/adjustment/contribution/margin); contrast `internal/modules/pricing/adapters/postgres/matrix_reader.go:34-39,64-68` (`aliquota::text`) | 25 `pgtype.Float8`/`pgtype.Numeric` sites across 2 files vs 21 files using `::text`-cast money/decimal reads elsewhere — no shared convention |
| P-6 | hand-sync | Enum-like vocabularies are hand-typed into `CHECK (... IN (...))` per table, with zero shared type — the same vocabulary is independently retyped in multiple migrations | `migrations/0050_market_price_snapshots.sql:16` and `migrations/0053_market_signals_aggregates.sql:40` both declare `CHECK (source IN ('ml_sale_price','ml_price_to_win','ml_catalog_offers'))` verbatim; `migrations/0076_products_mirror_active_source.sql:38,64` and `migrations/0078_products_mirror_key_by_source.sql:61` independently declare `CHECK (source/active_source IN ('sankhya','xlsx','catalogo_cliente'))` 3 times | 65 `CHECK (...IN (...))` constraints total, 0 `CREATE TYPE`, 0 shared `DOMAIN`; at least 2 vocabularies confirmed duplicated verbatim across 2-3 files |
| P-7 | gap | The Oracle→Postgres mirror stamps staleness (`stale_since`, `absent_in_last_snapshot`) correctly and atomically, but `stale_since` is write-only — no Go code path ever reads it back | `apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_repository.go:135,142` (writes); `grep -rn stale_since internal/modules/*/adapters/postgres/*.go` → only these 2 write sites, 0 `SELECT`/scan sites anywhere in the tree | 1 column, 2 write call sites, 0 read call sites |
| P-8 | gap | Migration numbering has 2 duplicate numbers and multiple gaps; nothing enforces sequential numbering — collisions are caught only because the runner tracks full filenames, not numbers | `apps/server_core/migrations/0021_integration_operation_run_evidence.sql` + `0021_integrations_provider_auth_strategy_shopee_partner.sql`; `0093_orders_status_details_nullable.sql` + `0093_sync_state_market_queue_entity_split.sql`; runner at `internal/platform/migrate/runner.go:25` (`sort.Strings(filenames)`, tracked by filename in `schema_migrations`, not by number) | 85 `.sql` files, 83 unique leading numbers, 2 collisions, gaps at 0002, 0040-42, 0048-49, 0054, 0059-64, 0077, 0095 |
| P-9 | idiom | Dynamic WHERE-clause construction via `fmt.Sprintf`/`strings.Join` is a hand-maintained safe pattern (values always via `$N`, never string-interpolated) but nothing prevents a future call site from interpolating a request value directly into the `%s` bucket | `internal/modules/listings/adapters/postgres/repository.go:93,122,131,135,167,228,268,297`; `internal/modules/mutations/adapters/postgres/read_repository.go:74,78,82,85,115,118` | 14 call sites across 2 files build SQL text with `Sprintf`; sampled all 14, all currently safe |

## The five heaviest, with detail

**P-1/P-3 — RLS protects almost nothing, and the one place it was done right is neutralized anyway.**
`migrations/0097_catalog_context.sql` is the only migration that touches Row Level Security at all
(`grep -n "ENABLE ROW LEVEL SECURITY\|FORCE ROW LEVEL SECURITY\|CREATE POLICY" migrations/*.sql` → 12
lines, all in this one file, covering 4 tables: `catalog.products`, `catalog.product_identifiers`,
`catalog.source_product_keys`, `catalog.source_observations`). The other 60 tenant-scoped tables have
no RLS at all — not even a decorative policy. Migration `0098_catalog_app_role.sql` then goes further
than any other part of the schema: it creates a dedicated `mpc_app` role, explicitly `NOSUPERUSER
NOBYPASSRLS`, scoped only to the `catalog` schema (`GRANT ... ON catalog.* TO mpc_app`, no grant on any
other schema) — the comment at `0098_catalog_app_role.sql:1-7` names the exact defect it fixes
("Every policy was evaluated against a role that skips policies"). This is genuinely correct database
engineering. But `grep -rn "mpc_app\|SET ROLE\|SET SESSION AUTHORIZATION" apps/server_core --include=*.go`
returns zero hits outside the migration file itself — no Go code ever connects as, or switches into,
`mpc_app`. `internal/platform/pgdb/config.go` and `pool.go` read one `MC_DATABASE_URL` and open one pool;
there is no code path that would ever use the role the migration built. Established fact 12/D-44 says the
production DSN connects as the table owner (a superuser), which always bypasses RLS regardless of
`FORCE`. So the 4-table mechanism — `set_config('app.tenant_id', ...)` per transaction
(`internal/contexts/catalog/internal/postgres/repository.go:49`) plus the `FORCE`d policy — is complete,
correct, and unreachable in the running application. It is not "half-built"; it is finished and disconnected.

**P-2 — tenant scoping is 100% hand-discipline, verified correct on every sample checked, enforced by nothing.**
Of 64 live tables, 61 carry a `tenant_id` column (`scratch_analyze_tables.py`, parses every
`CREATE TABLE` block for `tenant_id`, cross-checked against `DROP TABLE` in `0081_drop_legacy_...sql`).
Every non-test Go file that issues a `.Query`/`.QueryRow`/`.Exec` against Postgres textually mentions
`tenant_id`, with exactly one legitimate exception (`marketplaces/adapters/postgres/fee_schedule_repo.go`,
which only touches the two genuinely global tables `marketplace_definitions`/`marketplace_fee_schedules`).
A statement-level scan (`scratch_sql_scan.py`, extracts every backtick Go string containing
SELECT/INSERT/UPDATE/DELETE) found 285 of 285 statements against tenant tables carry `tenant_id` in the
same statement or a concatenated sibling segment — sampled 6 of the ones the naive scanner flagged as
missing it (`erp_import/mirror_query_repository.go`, `listings/repository.go`,
`mutations/read_repository.go`, `product_links/link_candidate_repo.go`) and all 6 turned out to build the
predicate correctly, just via string concatenation the regex couldn't see across segments. The convention
is real and consistently followed — `internal/modules/listings/adapters/postgres/repository.go:263-265`
(`args := []any{tenantID, ...}; where := []string{"l.tenant_id = $1", ...}`) is the canonical shape,
repeated by hand in every adapter. Nothing compiles, lints, or tests for its absence: `arch-gate.sh` and
everything under `scripts/` never mention "tenant" (`grep -rln tenant scripts/lib scripts/arch-gate.sh
scripts/tests` → 0 hits). The correctness of every one of those 285 statements today rests entirely on a
human having typed the right thing, with zero mechanical backstop if a new query omits it.

**P-4 — `tenant_default` silent fallback is only half-closed; `cmd/server` still inherits it.**
D-39 records this exact defect for `cmd/catalogingest` and shows it CLOSED with a real guard
(`requireTenantConfigured`, `cmd/catalogingest/main.go:48,98-123`, RED/GREEN proof included in the ledger).
But the ledger entry itself says the fix was scoped narrowly: *"`pgdb.LoadConfig()`/`config.go` NÃO foi
tocado — o fallback silencioso para `cmd/server` e os demais consumidores partilhados continua exatamente
como estava."* Verified unchanged at HEAD: `pgdb/config.go:23-24` still does
`if cfg.DefaultTenantID == "" { cfg.DefaultTenantID = "tenant_default" }`, and `cmd/server/main.go` calls
`pgdb.LoadConfig()` with no equivalent guard. `internal/composition/root.go` then threads
`cfg.DefaultTenantID` into 70 repository/service constructors (`grep -c "cfg.DefaultTenantID"
internal/composition/root.go` → 70) — every one of them, including `erpRepo`, `classRepo`,
`installationRepo`, `credentialRepo`, `mutationRepo`, `productLinkSnapshotRepo`, inherits whatever value
`LoadConfig()` produced, silently, for the lifetime of the process. `AGENTS.md:24` states "unknown
operational facts never become zero/default" as a binding rule; this is the rule's namesake violation,
alive today, for the code path that matters more than the batch importer — the running HTTP server.

**P-5 — money type handling is inconsistent between two adapters that both know better.**
`pricing/matrix_reader.go:64-68` reads `icms_matrix_mirror.codtrib`/`icms_aliquota_interna.aliquota` via
explicit `::text` casts, with a doc comment explaining exactly why: *"Numeric columns come back via
`::text` so values never round-trip through float64"* (`matrix_reader.go:17-18`). That discipline is not
shared. `orders/adapters/postgres/order_repo.go:180,1029` scans `unit_price numeric(14,2)`
(`migrations/0027_orders_marketplace_orders.sql:31`) through `scanFloat8`/`pgtype.Float8`
(`order_repo.go:1124,1154`) straight into `*float64`. `profitability/adapters/postgres/store.go:89,
246-254` does the same for 10 distinct money fields in one function (`revenue`, `saleFee`, `cost`, `tax`,
`freight`, `commission`, `adjustment`, `contribution`, `margin`, plus `amount`). Existing memory
(`orders-float64-safe-by-accident.md`) already establishes that for `orders` specifically the
`numeric(14,2)` domain keeps this safe today by coincidence, not by contract, and that no line of code
asserts the premise. This lane's contribution is sizing: the `::text` discipline that would remove the
risk exists, is documented, and is used in exactly 1 module (`pricing`) out of at least 3 that scan money
(`orders`, `profitability`, `pricing`).

**P-6 — 65 hand-typed CHECK vocabularies, 0 shared enum type, confirmed duplication in at least 2 cases.**
`grep -n "CHECK.*IN (" migrations/*.sql | wc -l` → 65. `grep -n "CREATE TYPE" migrations/*.sql` → 0 — no
native Postgres enum, no domain type, anywhere in 98 migrations. Two vocabularies are independently
retyped, verbatim, in separate files with separate constraint names: the market "source" enum
(`ml_sale_price`/`ml_price_to_win`/`ml_catalog_offers`) at `0050_market_price_snapshots.sql:16`
(`market_price_snapshots_source_check`) and again at `0053_market_signals_aggregates.sql:40`
(`market_aggregates_source_check`); and the mirror "source" enum (`sankhya`/`xlsx`/`catalogo_cliente`) at
`0076_products_mirror_active_source.sql:38` (`products_mirror.source`), again at
`0076_products_mirror_active_source.sql:64` (`products_mirror.active_source`, same file, same migration,
second column), and again at `0078_products_mirror_key_by_source.sql:61`. Nothing keeps these three
declarations in agreement if a fourth ERP source is ever added — each `ALTER TABLE ... ADD CONSTRAINT`
would need to be found and edited by hand (0072 already demonstrates the pattern of widening one CHECK at
a time: `catalogo_cliente` was added to the `erp_import_protocols.source` vocabulary by a dedicated
migration).

## Machinery vs hand

| correctness property | enforced by | count enforced (machinery) | count by hand |
|---|---|---|---|
| Tenant isolation on read/write | Postgres RLS policy (only where present) | 4 tables (`catalog.*`, `0097_catalog_context.sql:80-110`) | 57 tenant-scoped tables outside `catalog` schema (61 total minus 4) rely solely on the hand-typed `WHERE tenant_id=$1` convention |
| Tenant isolation, effective (connecting role) | RLS + non-bypassing role | 0 tables (D-44: DSN connects as owner; confirmed no `SET ROLE`/`mpc_app` usage anywhere in Go) | all 64 live tables |
| SQL well-formedness / no injection | query builder or codegen | 0 — no ORM, no sqlc, no squirrel/goqu/jet in `go.mod` | 246 raw `.Query`/`.QueryRow`/`.Exec` call sites |
| Enum/vocabulary consistency | `CREATE TYPE` / shared domain | 0 | 65 `CHECK (...IN(...))` constraints, confirmed duplicated in ≥2 vocabularies |
| Money precision (numeric → Go) | explicit `::text` cast convention | 1 module (`pricing/matrix_reader.go`), 21 files use `::text` somewhere | 25 `pgtype.Float8`/`pgtype.Numeric` scan sites in `orders`, `profitability` |
| Migration application idempotency (same file twice) | `schema_migrations` filename-tracked runner, per-migration transaction | all 85 migration files (`internal/platform/migrate/runner.go:33-105`) | — this one IS machinery; see "what is fine" |
| Migration sequencing (no two migrations claim the same number/intent slot) | none — numbers are cosmetic, filename is the real key | 0 | 85 files, author must pick a free-looking number by eye (2 collisions found) |
| Startup tenant identity, fail-closed on missing config | `requireTenantConfigured` guard | 1 command (`cmd/catalogingest`, D-39) | `cmd/server` (70 wiring sites in `root.go`) still silently defaults |
| Oracle mirror staleness detection | staging table + transactional upsert + `stale_since`/`absent_in_last_snapshot` columns | write side: 100% (`mirror_repository.go:135,142`) | read side: the columns exist but 0 Go call sites select `stale_since` — surfacing it is by-hand-not-yet-done, i.e. not done |

## What is actually fine

- **The migration runner's idempotency is real machinery, not discipline.** `internal/platform/migrate/
  runner.go:33-105`: each migration runs inside its own transaction, is only recorded in
  `schema_migrations` on commit, and a partial failure rolls back cleanly and is retried whole on the next
  run. `Filenames()` sorts lexicographically and applies only unapplied filenames — re-running the binary
  is safe. No `CREATE INDEX CONCURRENTLY` anywhere (`grep -ln CONCURRENTLY migrations/*.sql` → 0), so
  nothing breaks the "run inside a transaction" assumption.
- **All 116 timestamp columns are `timestamptz`.** `grep -c timestamptz migrations/*.sql` sums to 116;
  `grep -n "TIMESTAMP\b" migrations/*.sql | grep -vi timestamptz` → 0 hits. No naive-datetime/timezone
  hazard found anywhere in the schema.
- **The ERP→Postgres mirror write path is genuinely well-built.** `mirror_repository.go:71-172`
  (`mergeSnapshotTx`): stage into a `TEMP TABLE ... ON COMMIT DROP`, bulk-load via `CopyFrom`, then a
  single transaction does upsert + explicit stale-marking (`absent_in_last_snapshot=true,
  stale_since=COALESCE(stale_since,now())` for rows missing from the new snapshot) + child-table
  replacement, all-or-nothing. It does not silently produce a stale-but-plausible row: absence is marked,
  not overwritten with defaults, consistent with ADR-17 as cited in the code's own comments
  (`mirror_repository.go:111-116`).
- **Two non-tenant tables (`marketplace_definitions`, `marketplace_fee_schedules`) correctly have no
  `tenant_id` and their one adapter (`fee_schedule_repo.go`) correctly never scopes by tenant** — verified
  by reading the migrations (`0010_marketplace_definitions.sql`, `0011_marketplace_fee_schedules.sql`):
  both are genuinely global platform metadata (marketplace plugin catalog, seeded fee schedules), not
  tenant-owned data. This is not a gap; it's the one case in the codebase where "no tenant_id" is correct.
- **The 14 `fmt.Sprintf`-built WHERE clauses sampled were all safe.** Every dynamic-predicate site in
  `listings/repository.go` and `mutations/read_repository.go` puts only statically-chosen SQL fragments
  into the `%s`/`%d` slots; every user-supplied value goes through a `$N` placeholder. No string
  concatenation of a request value into SQL text was found anywhere in the 246 sampled call sites.
- **`icms_aliquota_interna` is correctly documented and coded as the one global (non-tenant) exception
  inside an otherwise tenant-scoped reader** (`matrix_reader.go:49-61`) — the comment explains why, and
  the code matches the comment.

## Unverified / needs judgment

- Whether the dev Postgres instance currently holds more than one distinct `tenant_id` value — I did not
  query the running database (`.env`/`MC_DATABASE_URL` read was denied by the permission system before I
  could inspect it), so I cannot say whether the tenant-isolation gap (P-1/P-2) is a live risk today or
  latent for a single-tenant deployment. Treat as `unverified`.
- Whether any HTTP middleware anywhere derives a per-request tenant from an authenticated principal —
  `grep -rn "X-Tenant\|tenant.*header"` across `internal` found nothing, and `cfg.DefaultTenantID` is
  wired once at process boot in `root.go`, which reads as single-tenant-per-deployment by design rather
  than a defect. I report the fact (no per-request derivation found); whether that is intentional is a
  product/architecture question outside this lane's mandate.
- The exact number of `float64` money **fields** (as opposed to scan call sites) across the whole domain
  layer — the earlier grep sample (`internal/modules/*/domain/*.go`) returned 20+ struct fields but I did
  not exhaustively enumerate every module; `P-5`'s scale claim is about scan call sites, which I did count
  precisely (25 `pgtype.Float8`/`Numeric` sites in 2 files, vs 21 files using `::text` for *something*,
  not necessarily money in every case — the 21-file figure is an upper bound on the disciplined side, not
  an exact count of money-specific `::text` casts).
- Whether the 2 duplicate migration numbers (0021, 0093) ever caused an actual ordering problem — the
  runner sorts by full filename, and in both observed pairs the alphabetically-earlier file has no
  documented dependency on the later one, but I did not diff column-level dependencies exhaustively.

## Commands run

```
git ls-files ... (see PHASE-0.md, not rerun)
ls apps/server_core/migrations | wc -l                                           → 101 (85 .sql + 16 _test.go/source.go)
cd apps/server_core/migrations && ls *.sql | sed -E 's/^([0-9]+)_.*/\1/' | sort -u | wc -l   → 83 unique numbers / 85 files
grep -n "ENABLE ROW LEVEL SECURITY|FORCE ROW LEVEL SECURITY|CREATE POLICY|DISABLE ROW LEVEL SECURITY" migrations/*.sql
grep -rn "SET ROLE|SET SESSION AUTHORIZATION|mpc_app" --include=*.go apps/server_core   (excl. _test.go)
python scratch_analyze_tables.py   (custom: parses every CREATE TABLE block for tenant_id, cross-checks DROP TABLE)  → 64 live tables, 61 with tenant_id
python scratch_sql_scan.py          (custom: extracts Go backtick SQL literals, classifies SELECT/INSERT/UPDATE/DELETE by tenant_id presence)
grep -rlE "\.Query\(ctx|\.QueryRow\(ctx|\.Exec\(ctx" --include=*.go apps/server_core/internal | grep -v _test.go | wc -l   → 246
for f in <query files>; do grep -q tenant_id "$f" || echo "$f"; done   → 2 exceptions, both legitimate
grep -rn "mpc_app|DATABASE_URL|PGUSER|pgxpool.New|connString" --include=*.go internal/platform cmd
grep -c "cfg.DefaultTenantID" apps/server_core/internal/composition/root.go   → 70
grep -n "DefaultTenantID|MC_DEFAULT_TENANT_ID|tenant_default" apps/server_core/internal/platform/pgdb/config.go
grep -rln "pgtype.Float8|pgtype.Numeric" --include=*.go apps/server_core/internal | grep -v _test.go   → 2 files, 25 sites
grep -rln "::text" apps/server_core/internal/modules/*/adapters/postgres/*.go   → 21 files
grep -n "CHECK.*IN (" apps/server_core/migrations/*.sql | wc -l   → 65
grep -n "CREATE TYPE" apps/server_core/migrations/*.sql | wc -l   → 0
grep -c timestamptz apps/server_core/migrations/*.sql | awk -F: '{s+=$2}END{print s}'   → 116
grep -n "TIMESTAMP\b" apps/server_core/migrations/*.sql | grep -vi timestamptz   → 0 hits
grep -ln CONCURRENTLY apps/server_core/migrations/*.sql   → 0 hits
grep -rln tenant scripts/lib scripts/arch-gate.sh scripts/tests   → 0 hits
grep -i "pgx|database|sql|orm" apps/server_core/go.mod   → only github.com/jackc/pgx/v5
```
