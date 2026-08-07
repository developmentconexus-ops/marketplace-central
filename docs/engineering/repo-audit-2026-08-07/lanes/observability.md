# Lane: observability

> I have run something really deep into MetalDocs and I am changing the way I code there to move to
> something more professional towards issues, PRs, PR review, CodeRabbit mechanical full validation
> and so much more. For that I had to identify every error in my code, my platform, to improve it and
> to create this full validation. I want to run it here as well so we move on the same path, this way
> it gets so much harder to send bad PRs.

Calibration: solid professional level, not Google-tier. Success condition is a **mechanism** that
makes a bad change hard to land, not a cleaner codebase.

Scope: `apps/server_core` (Go) sync/scheduler/integration paths, HTTP transport, health surfaces.
Method: `Documents\MetalDocs\docs\engineering\repo-audit-playbook.md`. Repo @ `1473e863`.

## Findings

| ID | class | finding | evidence | scale |
|---|---|---|---|---|
| OBS-1 | gap | `/healthz` returns 200 unconditionally — no DB, Oracle, or ML check — and is the literal Docker container healthcheck | `internal/platform/httpx/router.go:5-14`; `docker-compose.yml:46` (`curl -fsS http://localhost:8080/healthz`) | 1 endpoint, 1 healthcheck config, 100% of container-health signal |
| OBS-2 | hazard | 5 of 6 `time.NewTicker`-driven background loops have zero panic recovery anywhere in the goroutine; an unrecovered panic in any one kills the entire process, taking down the HTTP server (and every other loop) with it | `internal/modules/integrations/background/{refresh_ticker.go:70-88, state_cleanup.go:38-52, fee_sync_scheduler.go:40-55}`, `internal/modules/internal_read/observability/pool_stats.go:44-81`, `internal/modules/mutations/background/poller.go:36-76` — none contain `recover()`. Only `internal/modules/sync/application/scheduler.go` has a `recover()` (`safeInvoke`, line 201-208), and it covers only the per-job `JobFunc` call, not the surrounding `store.Read`/`RecordFailure`/`RecordSuccess` calls in `runJob` | 5/6 background-loop files, 0 defense in depth for `cmd/server/main.go`'s only process |
| OBS-3 | gap | `FeeSyncScheduler` has zero logging anywhere in its file: ticker tick, per-installation `StartSync` error, and `RunOnce` error are all discarded with `_ =` / `_, _ =` | `internal/modules/integrations/background/fee_sync_scheduler.go:52` (`_ = s.RunOnce(ctx)`), `:74` (`_, _ = s.flow.StartSync(ctx, ...)`) — file imports neither `log/slog` nor `log` | 1 scheduler, 2 discard sites, 0 log lines possible from this file ever |
| OBS-4 | gap | `sync.Scheduler.runJob` silently skips on a cursor-read error (no log, no metric) and discards the error from writing the failure/success record itself — the durable record that is the *only* observability channel for this scheduler | `internal/modules/sync/application/scheduler.go:146-149` (`if err != nil { return }`, no log call), `:154` (`_ = s.store.RecordFailure(...)`), `:160` (`_ = s.store.RecordSuccess(...)`) — `Scheduler` has no `logger` field at all | 3 sites in the one scheduler type composing products, market, and ICMS-matrix sync |
| OBS-5 | gap | 13 of 15 Mercado Livre adapter files never call `slog` — an outbound call failing propagates as a bare Go error with no record that the call was even attempted, unless a caller several layers up happens to log it (several do not, see OBS-3) | `grep -rl "slog\." internal/modules/connectors/adapters/mercado_livre --include=*.go \| grep -v _test.go` → 2 files (`pricing_reader.go`, `catalog_offers_reader.go`) out of 15 | 13/15 files, 87% |
| OBS-6 | gap | `/sync/health` rows are deliberately not keyed by `installation_id` (one row per `sync_state` PK `(tenant_id, installation_id, entity)`, collapsed to just `entity` in the DTO); the FE list renders one row per API entity using `entity.entity` as the React key — with 2+ ML installations sharing an entity (e.g. two `orders` schedulers), the health card cannot show, or even key, more than one installation's status for that entity | `internal/modules/sync/application/health_reader.go:11-14` (doc-comment concedes this is deliberate), `internal/modules/sync/adapters/postgres/health_reader.go:33-39` (SQL selects no `installation_id`), `migrations/0075_sync_sync_state.sql:29` (PK includes `installation_id`), `apps/web/src/pages/integracoes/SyncHealthCard.tsx:169` (`key={entity.entity}`) | every tenant with 2+ ML installations; `orders`/`listings` schedulers are per-installation (`internal/modules/orders/composition/scheduler.go:129-140`) |
| OBS-7 | gap | Zero metrics or tracing libraries anywhere in the module graph — no Prometheus client, no OpenTelemetry, no `/metrics` route | `grep -rn "opentelemetry\|otel\|prometheus" go.mod` → 0 hits; `find internal -iname "*metric*"` → 0 files | 0 instrumented paths / all paths |
| OBS-8 | gap | No request/correlation ID is generated or propagated anywhere in the HTTP entry path; the only middleware chain is CORS + panic-recover | `internal/composition/root.go:994` (`httpx.CORSMiddleware(apierror.Recover(mux))`); `grep -i "requestid\|correlation"` across `internal` returns only unrelated domain idempotency keys (`market/application/collection_pipeline_service.go:467`) | 0 of N HTTP handlers carry a correlation id into their logs or DB queries |
| OBS-9 | idiom | `slog` is used everywhere (36 files) but no `slog.Handler` is ever configured (`slog.SetDefault`, `NewJSONHandler`) — every call goes through `slog.Default()`, Go's built-in **text** handler. Structured fields exist in the API call but the actual stdout bytes are `key=value` text, not machine-parseable JSON | `grep -rn "slog.SetDefault\|NewJSONHandler\|NewTextHandler\|HandlerOptions" internal cmd` → 0 hits | 36 files depend on a handler that is never set |
| OBS-10 | hand-sync | Two logging mechanisms coexist: `log/slog` (36 files, the de facto standard) and stdlib `log` (`internal/platform/logging/logger.go`, used once in `cmd/server/main.go:19,29`, plus `cmd/migrate/main.go`) — the process's own boot line (`"server starting on %s"`) is the one log statement in the whole backend that is not structured | `internal/platform/logging/logger.go:8-10`; `cmd/server/main.go:19,29` | 2 files outside the slog convention |
| OBS-11 | gap (re-measured, ledger D-45) | `cmd/catalogingest` failure surfaces only as stderr text + non-zero exit code; no `sync_state`/`sync_health` row exists for the `catalog` context, so an operator not watching the terminal at the exact moment has no way to learn a catalog ingest run failed | `apps/server_core/cmd/catalogingest/main.go:34-40` — still true against current tree; ledger entry `.mnfs/HARNESS-DEBTS.md:1536-1542` (D-45), unchanged | 1 command, 0 durable failure record |
| OBS-12 | gap | `/sync/health`'s webhook block is permanently the canonical zero (`defaultWebhookStatsReader`); `WithWebhookStatsReader` is defined but never called from `composition/root.go` — there is also no webhook receiver route registered anywhere in the backend, so this isn't a wiring gap so much as a subsystem that was never built | `internal/modules/sync/application/health_webhook_stats.go:24-30`; `grep -rn "WithWebhookStatsReader" internal` → only the definition, no call site; `grep -rli webhook internal` finds no receiver handler | 1 API block permanently zero; FE (`SyncHealthCard.tsx:78-92`) does distinguish "never configured" from "quiet" honestly in its copy, so this is low severity |

## The five heaviest, with detail

**1. OBS-2 — one panic anywhere in five background loops takes the whole server down.**
`cmd/server/main.go` runs a single OS process. `composition/root.go` starts, via bare `go` statements,
at least these ticker loops with no panic isolation: `RefreshTicker.Start` (`refresh_ticker.go:70`),
`StateCleanup.Start` (`state_cleanup.go:38`), `FeeSyncScheduler.Start` (`fee_sync_scheduler.go:40`),
`PoolStatsLoop.Start` (`pool_stats.go:44`, itself the DB-observability primitive), and
`mutations/background.Poller.Run` (`poller.go:36`). Go's runtime terminates the entire process on an
unrecovered panic in *any* goroutine — there is no per-goroutine isolation by default. A nil map/slice
access inside any one of these five (e.g. a `nil` `*pgxpool.Pool` result row, a type assertion on an
unexpected provider payload) ends order ingestion, listing sync, fee sync, pool stats, *and* the HTTP
API simultaneously, with the only trace being whatever `slog.Error` fires from `apierror.Recover`
(`internal/platform/apierror/recover.go:35`) for the *next* HTTP request that happens to hit the dead
process before the container restarts it. Only `sync.Scheduler.safeInvoke` (`scheduler.go:201-208`)
demonstrates the fix is already known in this codebase — it is just not applied to the other five
loops.

**2. OBS-1 — the only automated health signal that exists is structurally incapable of ever failing.**
`httpx.NewRouter()`'s `/healthz` handler (`router.go:7-12`) writes `{"service":..., "status":"ok"}`
unconditionally — it takes no dependency, checks no pool, calls no `Ping`. `docker-compose.yml:46`
wires this exact endpoint as the container's `HEALTHCHECK`. This means Docker's own judgment of
"is this container healthy" is, today, `net/http`'s ability to accept a TCP connection and nothing
else. A Postgres pool at `MaxConns` for an hour, an Oracle connection wedged, or every Mercado Livre
credential in `requires_reauth` — none of these move `/healthz` off 200. (`deploy/docker-compose.prod.yml`
has no healthcheck on the backend service at all — only Postgres does, `deploy/docker-compose.prod.yml:16-19`.)
There *is* a real, DB-backed health surface — `/sync/health` (`internal/modules/sync/transport/health_handler.go`)
— but it is not wired as any healthcheck, liveness, or readiness probe; it is FE-only.

**3. OBS-3/OBS-4 — the mechanism that is supposed to make failure durable can itself fail silently, twice over.**
The comment at `auth_flow_service.go:494-498` documents that the ML-token-refresh class of invisible
failure was *fixed* by writing `RefreshFailureCode`/`ConsecutiveFailures`/`NextRetryAt` to
`integration_auth_sessions` (see "What is actually fine" below) — but that fix pattern was not applied
uniformly. `FeeSyncScheduler` (OBS-3) has no logging mechanism in the file at all: a failed
`StartSync` for every connected, fee-sync-capable installation, every 15 minutes, vanishes with
`_, _ = s.flow.StartSync(...)`. `sync.Scheduler.runJob` (OBS-4) *does* attempt to persist failure via
`RecordFailure`, but that write's own error is discarded (`_ = s.store.RecordFailure(...)`,
`scheduler.go:154`) — if Postgres is unreachable at the exact moment a job fails, both the original
failure *and* the fact that it could not be recorded disappear with no log line anywhere. The read
side is worse: `s.store.Read` failing (`scheduler.go:146-149`) returns silently with no log call at
all — a job entity can go dark forever (every tick, `Read` fails, cycle is skipped) with zero signal
distinguishable from "interval hasn't elapsed yet."

**4. OBS-5 — the marketplace integration itself is nearly unobserved at the adapter layer.**
13 of the 15 files under `internal/modules/connectors/adapters/mercado_livre/` — `listing_writer.go`,
`price_writer.go`, `items_multiget_reader.go`, `shipping_reader.go`, `shipment_ingest_reader.go`,
`order_ingest_reader.go`, `capability_adapter.go`, `buyer_fiscal_reader.go`,
`catalog_identity_reader.go`, `catalog_match_reader.go`, `items_scan_ids_reader.go`,
`resilience_decorator.go`, `oauth`-adjacent files — never call `slog`. There is no access-log
equivalent for calls to Mercado Livre: no line recording that a request to `/items/{id}` or a
listing write was even attempted, its status code, or its latency, independent of whether the caller
several layers up chooses to log the returned error. Combined with OBS-3/OBS-4 (several of those
callers *don't* log), a live-API failure in price or stock writes can be completely untraceable after
the fact.

**5. OBS-7/OBS-8 — no metrics, no tracing, no correlation id, anywhere.**
`go.mod` has no OpenTelemetry or Prometheus dependency; no `/metrics` route exists. The HTTP entry
path (`CORSMiddleware(apierror.Recover(mux))`, `root.go:994`) never mints or reads a request id, so a
log line from a database timeout inside a handler cannot be joined to the HTTP request that triggered
it, nor to the downstream Mercado Livre call it may have made, without hand-matching timestamps. For
a "solid professional level" bar this is the baseline gap underneath every other finding in this
report: every other observability signal in this codebase is a log line or a DB row keyed by
installation/entity, never a trace.

## Silent-failure census

| file:line | failure | surfaces where? | operator-visible? |
|---|---|---|---|
| `integrations/background/fee_sync_scheduler.go:52` | `RunOnce` error each tick | nowhere | no |
| `integrations/background/fee_sync_scheduler.go:74` | `StartSync` error per installation | nowhere | no |
| `sync/application/scheduler.go:146-149` | `store.Read` (cursor read) error | nowhere (no log call exists) | no |
| `sync/application/scheduler.go:154` | `RecordFailure` write itself failing | nowhere | no |
| `sync/application/scheduler.go:160` | `RecordSuccess` write itself failing | nowhere | no |
| `integrations/application/fee_sync_service.go:182` | `persistFailedOperationRun` failing after `ExecuteSync` already failed | nowhere (original `execErr` and the persist error both vanish) | no |
| `orders/adapters/postgres/order_repo.go:297,413,1109,1112` | `json.Unmarshal` of stored `tags`/`raw_provider_ref` JSONB failing | nowhere — field silently reads as zero value | no (contradicts AGENTS.md: unknown → never a default) |
| `apps/server_core/cmd/catalogingest/main.go:34-40` | ingest run failing | stderr + exit code only, no durable row | only if someone is watching the terminal live (ledger D-45, re-measured, still true) |
| `sync/composition/products_job.go:52`, `orders/composition/scheduler.go:145`, `listings/composition/scheduler.go:164` | `RegisterJob` error at boot discarded (`_ =`) | nowhere | no — but comment claims (and code shows) it structurally cannot fail at these call sites, so low severity |
| `integrations/background/state_cleanup.go:49` | expired-OAuth-state cleanup `RunOnce` error | nowhere (`_ = s.RunOnce(ctx)`, no logger field in the struct at all) | no — low severity (cleanup only, not a data-correctness path) |

Cross-cutting: none of the above is at debug level — the codebase has **zero** `slog.Debug` calls
(`grep -rn "slog.Debug" internal` → 0 hits), so "logged but invisible in prod" is not a failure mode
here; the failure mode is "not logged at all."

**What is *not* silent, verified against the specific prior claim this brief asked to re-check:**
Mercado Livre token-refresh failure is **no longer** invisible. `AuthFlowService.RefreshCredential`
(`auth_flow_service.go:454-532`) calls `recordRefreshFailure` (`:537-574`) on any refresh error, which
(a) writes `RefreshFailureCode`/`ConsecutiveFailures`/`NextRetryAt` to `integration_auth_sessions`
(durable), and (b) calls `degradeAfterRefreshFailure` (`:582-622`), which flips the installation to
`requires_reauth`/`degraded` + `HealthStatusCritical`/`Warning` when the failure is terminal or
persistent. `RefreshTicker` (`integrations/background/refresh_ticker.go:58-65`) also logs every
per-item failure via `slog.Error`. The frontend reads this through `ConnectionSnapshot.state`
(`ConnectionHealthCard.tsx:14-21,55`), which renders `needs_reauth`/`degraded` as a red/amber badge
with the raw provider error shown verbatim (`ConnectionHealthCard.tsx:111-118`) — the screen does
**not** stay green. This is proven against real Postgres by
`apps/server_core/tests/integration/integrations_refresh_failure_test.go`. The comment at
`auth_flow_service.go:494-498` explicitly narrates the old bug this fix replaced ("o ticker
descarta o erro, a sessão continua marcada como válida e /integracoes segue verde com o token
morto") — that description matches the prior measurement but no longer matches the code.

## "A bad sync happened last night" — what the operator can actually look at

1. **`/sync/health`** (`GET`, DB-backed, real): per-entity `consecutive_failures`, `last_error`
   (verbatim provider/job error text), `last_success_at` (GREATEST of full/incremental), `phase`.
   Rendered on `IntegracoesPage` via `SyncHealthCard.tsx` — red badge with failure count, "nunca"
   for never-synced. **Caveat (OBS-6):** with 2+ installations sharing an entity, only one row can be
   shown/keyed per entity name — the operator cannot tell *which* installation is failing from this
   card alone.
2. **`ConnectionHealthCard`** on the same page: per-installation OAuth/connection state
   (`connected`/`degraded`/`needs_reauth`/`disconnected`), with the raw reauth reason string and a
   one-click reauth button when applicable.
3. **Container logs** (`docker logs`, or wherever stdout is captured): `slog`-formatted `key=value`
   text lines for whichever code paths actually call `slog` (36 of the ~90+ Go files with error-return
   paths in scope) — **not JSON** (OBS-9), so grepping by field name works but piping into a log
   aggregator that expects JSON does not without a parser shim. Silent paths (the census above)
   produce **no line at all**, so their absence from the logs is indistinguishable from "nothing went
   wrong there."
4. **`sync_state` table directly** (`psql`): the ground truth `/sync/health` reads from — same
   caveats (per-`(tenant, installation, entity)` row, no aggregation across installations).
5. **What does *not* exist**: no dashboard, no metrics, no trace, no alert — an operator only learns
   of a bad sync by opening `/integracoes` (or `psql`) and looking; nothing pages them. `/healthz`
   (OBS-1) will not have flagged anything regardless of how bad the night was, since it checks
   nothing. If the failure happened in `FeeSyncScheduler` (OBS-3) or was a cursor-read failure inside
   `sync.Scheduler.runJob` (OBS-4), there is no `last_error` written at all — the operator sees a
   merely-stale `last_success_at` with no explanation, indistinguishable from "the interval hasn't
   elapsed yet." If it was `cmd/catalogingest` (OBS-11), there is nothing on any screen — only the
   terminal that ran it, if anyone was watching.

## What is actually fine

- **The ML-token-refresh failure path is fixed and durable**, contrary to the prior measurement this
  brief asked to re-verify — see the census section above. Do not re-open this as if it were still
  broken; the fix is tested against real Postgres.
- **`/sync/health` + `SyncHealthCard.tsx` is genuinely well built.** It reads a state field, not a
  time cutoff, for its red/green discriminant (`entityTone`, `SyncHealthCard.tsx:13-17`); it
  distinguishes `pending`/`fetching`/`paused-idle`/`error`/`success` react-query states explicitly
  rather than collapsing them (`SyncHealthCardBody`, `:156-194`, with an explicit regression-guard
  comment citing a prior live-drive defect); it never renders blank on an unreadable state (ADR-17
  compliance, `:150-155,182-193`); and its webhook block honestly states "no notification received"
  rather than fabricating a "not configured" vs "quiet" distinction it cannot know
  (`SyncHealthCard.tsx:78-92`, `health_webhook_stats.go:24-30`).
- **`sync.Scheduler.safeInvoke`** (`scheduler.go:201-208`) is a real, working panic-isolation pattern
  — proof the fix for OBS-2 is already known inside this codebase, just not spread to the other five
  background loops.
- **Oracle connectivity is checked at boot and degrades honestly.** `composition/root.go:451-457`
  logs at `Warn` (with the unwrapped cause) and falls every downstream reader to an explicit
  `UnavailableStockBatchReader`/`unavailableProductMatcher` rather than a fabricated zero — consistent
  with AGENTS.md's "unknown operational facts never become zero/default."
  `internal_read/observability/pool_stats.go` and `NewTimingReader`
  (`internal_read/observability/…`, wired at `root.go:464-469`) give real, if log-line-only, DB pool
  and slow-query visibility.
- **`apierror.Recover`** (`internal/platform/apierror/recover.go`) correctly isolates HTTP-handler
  panics, logs the recovered value + stack trace, and answers a proper 500 instead of a dropped
  connection.
- **`cmd/catalogingest`'s tenant-fallback guard** (`main.go:98-123`) is a good example of the
  "unknown never becomes a default" rule applied correctly to a destructive write path — cited here
  because it is adjacent to OBS-11, which is a genuinely separate, still-open gap in the same command.
- **Zero `slog.Debug` usage anywhere** means the "errors logged at a level nobody reads" failure mode
  the brief asked about does not occur in this codebase — the actual failure mode is "not logged at
  all," which is the census above, not misleveled logging.

## Unverified / needs judgment

- Whether `docker-compose.prod.yml`'s absence of a backend healthcheck is intentional (external LB
  probe instead) or an oversight — `unverified`, no external LB config found in-repo.
- Runtime behavior of the duplicate React key in OBS-6 (silently drops one row vs. a dev-console
  warning with both rendered, order-dependent) was not exercised live — the structural collision
  (same `key` value for 2+ rows) is verified by code/schema; the exact React reconciliation outcome
  for this specific case is `unverified`.
- Whether any of the 5 unrecovered-panic background loops (OBS-2) has ever actually panicked in
  production — `unverified`, no crash logs or incident record reviewed as part of this lane.
- Log retention/shipping outside the container (whether `docker logs` output is captured anywhere
  durable) is outside this repo and `unverified`.

## Commands run

```
cd apps/server_core
grep -rn "RefreshToken|refresh_token|TokenRefresh" internal --include=*.go   # locate ML token-refresh code
grep -rln "\"log/slog\"" internal --include=*.go | grep -v _test.go | wc -l  # 36
grep -rl "^\s*\"log\"$" internal --include=*.go cmd --include=*.go | grep -v _test.go
grep -rl "time.NewTicker" internal --include=*.go | grep -v _test.go        # 6 background loop files
grep -rl "recover()" internal --include=*.go | grep -v _test.go             # 3 of the 6 have any recover, 1 covers job-body only
grep -rn "opentelemetry|otel|prometheus" go.mod                             # 0 hits
find internal -iname "*metric*"                                             # 0 files
grep -rli "requestid|correlation" internal --include=*.go | grep -v _test.go
grep -rln "slog\." internal/modules/connectors/adapters/mercado_livre --include=*.go | grep -v _test.go   # 2 of 15
ls internal/modules/connectors/adapters/mercado_livre/*.go | grep -v _test.go | wc -l                     # 15
grep -n "healthz" docker-compose.yml deploy/docker-compose.prod.yml
grep -n "WithWebhookStatsReader" internal -r --include=*.go
grep -rn "^\s*_\s*=\s*[a-zA-Z]" internal/modules/{sync,integrations,mutations,orders,listings,market} --include=*.go | grep -v _test.go
grep -n "PRIMARY KEY" apps/server_core/migrations/0075_sync_sync_state.sql
grep -n "^\*\*D-[0-9]*\." .mnfs/HARNESS-DEBTS.md   # cross-checked against existing ledger (D-16, D-43, D-44, D-45 relevant)
```
