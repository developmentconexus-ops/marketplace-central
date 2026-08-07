# Lane: security

Method: `C:\Users\leandro.theodoro\Documents\MetalDocs\docs\engineering\repo-audit-playbook.md`.
Repo: `marketplace-central` @ `11b9e4943b24717ebecddbcac9042b92ee99d8f9` (main, working tree
carried two pre-existing modified files and two untracked scratch files not authored by this
lane — left untouched, see Commands run).

Operator's own words, calibration: **solid professional level**, mechanism over discipline,
success = "harder to send bad PRs," not just cleaner code. No paid tooling. No second human
reviewer. This lane is discovery only — findings, not remediation architecture.

## Findings

| ID | class | finding | evidence | scale |
|---|---|---|---|---|
| SEC-1 | hazard | Zero authentication on the entire HTTP API. Final handler chain is CORS + panic-recovery only; no route requires a credential. | `internal/composition/root.go:994` `httpx.CORSMiddleware(apierror.Recover(mux))`; `internal/platform/httpx/router.go` (no auth); `grep -rln "RequireAuth\|checkAuth\|apiKeyMiddleware" internal/modules/*/transport` → 0 files; FE has no login page, no `Authorization` header sender (`grep -rn "Authorization" apps/web/src` → 0 hits) | 94 route registrations (`grep -rn 'mux\.HandleFunc(\|mux\.Handle(' internal/modules/*/transport/*.go` → 94), 0 authenticated |
| SEC-2 | hazard | Buyer PII — full name, government tax ID (CPF/CNPJ equivalent), full street address — served by the unauthenticated orders API. | `internal/modules/orders/transport/http_handler.go:608-618` `mapCompradorFiscal` puts `info.DocNumber` straight into `compradorFiscalDTO.DocNumero`; source columns `buyer_name, buyer_doc_type, buyer_doc_number, buyer_address_*` in `orders_marketplace_orders` (migration `0089_orders_marketplace_orders_sync_fields.sql:35,39`) | 1 endpoint (orders enrich/detail), 9 buyer_* PII columns per order row |
| SEC-3 | hazard/drift | `pgdb.LoadConfig()` — the config loader used by `cmd/server` (the real running application, not a batch tool) — silently defaults the tenant identity to `"tenant_default"` when `MC_DEFAULT_TENANT_ID` is unset, wiring that value into 30+ tenant-scoped repositories. This is the exact defect class just fixed in `cmd/catalogingest` (established fact: "an absent tenant variable silently became tenant_default"), unfixed in its much larger-blast-radius sibling. | `apps/server_core/internal/platform/pgdb/config.go:23-24`; used by `cmd/server/main.go:22` → `composition.NewRootRuntime(pool, dbCfg)`; `grep -c "cfg.DefaultTenantID" internal/composition/root.go` → 70 call sites across ~30 repository constructors | 1 config loader, ~30 tenant-scoped repos inherit it, entire server process |
| SEC-4 | gap | RLS is decorative today, and even the one code path that satisfies its precondition (`SET app.tenant_id`) is not on the app's live connection pool. Re-measures and extends established fact 12 / D-44. | Policies use `current_setting('app.tenant_id', true)` (`migrations/0097_catalog_context.sql:92-108`); only 1 file in the whole tree ever calls `set_config('app.tenant_id', ...)`: `internal/contexts/catalog/internal/postgres/repository.go:49`, inside a `WithTx` helper on the new-context write path, not the legacy `internal/modules/catalog` read path that actually serves `/catalog/*`; separately, `internal/composition/root.go` builds the pool from `MC_DATABASE_URL`, which is the `marketplace` superuser (`rolsuper=t, rolbypassrls=t`, D-44), so RLS is bypassed at the connection layer regardless | 4 tables have RLS+FORCE (`catalog.products/product_identifiers/source_product_keys/source_observations`); 0 of them are ever reached with the GUC set from the live app pool |
| SEC-5 | gap | RLS coverage is narrow even counting the decorative case: 1 of 85 migrations enables RLS at all; 52 migration files add a `tenant_id` column with zero DB-level backstop beyond the hand-written predicate. | `grep -rl "ROW LEVEL SECURITY" migrations/*.sql` → 1 file (`0097_catalog_context.sql`); `grep -rl "tenant_id" migrations/*.sql \| wc -l` → 52; `ls migrations/*.sql \| wc -l` → 85 | 51 of 52 tenant-bearing migrations, 0 RLS |
| SEC-6 | idiom | Permissive CORS (`Access-Control-Allow-Origin: *`) explicitly allows `Authorization` in `Access-Control-Allow-Headers`, on every route. Currently moot (SEC-1: nothing checks Authorization), but is a live footgun the moment auth is added — a wildcard origin plus an allowed Authorization header is the classic CSRF-adjacent misconfiguration once credentials exist. | `internal/platform/httpx/router.go:18-30`, own comment admits "in production this should be scoped to known domains" | 1 middleware, applies to all 94 routes |
| SEC-7 | idiom | Panic value passed verbatim to the log sink (not the client). If a panic ever carries a secret in its message (e.g. a formatted error including a token), it lands in logs unredacted, unlike the deliberate token/secret redaction built for provider response bodies elsewhere. | `internal/platform/apierror/recover.go:35` `slog.Error("apierror.Recover", "result", "500", "panic", rec, "stack", ...)` — `rec` is the raw `recover()` value | 1 sink, log-only exposure, not response body |
| SEC-8 | gap | `pgdb.DefaultTenantID(ctx, fallback string)` is a second, dead, fail-open tenant helper (`fallback == "" → "tenant_default"`) with zero callers — not currently exploitable, but the same fail-open shape as SEC-3 sitting in the same package, one call site away from being wired in by a future change without anyone noticing the precedent. | `apps/server_core/internal/platform/pgdb/tenant.go:5-10`; `grep -rn "pgdb\.DefaultTenantID" apps/server_core` → 0 call sites outside its own file | 1 function, 0 callers today |

## The five heaviest, with detail

**1. SEC-1 — no authentication anywhere (heaviest by construction).** `cmd/server/main.go` builds
`runtime.Handler = httpx.CORSMiddleware(apierror.Recover(mux))` and hands that straight to
`http.Server`. `apierror.Recover` only catches panics; `CORSMiddleware` only sets headers. Neither
inspects identity. Every one of the 94 `mux.HandleFunc`/`mux.Handle` registrations across
`internal/modules/*/transport/*.go` is reachable by anyone who can open a TCP connection to the
port — `docker-compose.yml:40` publishes it as `8080:8080` (host-exposed by default in the dev
stack; unknown in production, see Unverified). The frontend confirms this is not "auth happens
elsewhere and just isn't in this repo": there is no login page under `apps/web/src`, and no
component ever sets an `Authorization` header. `ARCHITECTURE.md:183`'s "authentication management"
line is scoped to the integrations/connectors module — that is marketplace **OAuth token**
management for talking *to* Mercado Livre/Melhor Envio, a completely different thing from
authenticating a caller *of* this API. There is no confusion in the code about this; the two
concepts simply never meet. Blast radius: every read (orders, pricing, margin, buyer PII,
inventory) and every write this API exposes (installation create/disconnect, mutation dispatch,
pricing config) is open to whatever can reach the socket.

**2. SEC-2 — real buyer PII served by that unauthenticated surface.** This is not a theoretical
join of two unrelated facts; the code deliberately carries the data end to end. Migration `0089`
adds 9 `buyer_*` columns to `orders_marketplace_orders`. `orders/adapters/postgres/buyer_fiscal_reader.go:48-57`
reads all 9, scoped only by a hand-written `WHERE tenant_id = $1 AND installation_id = $2 AND
provider_order_id = $3` (no RLS backstop per SEC-4/SEC-5). `orders/transport/http_handler.go:608-618`
maps the result straight onto the wire: `Nome: info.Name`, `DocTipo: info.DocType`, `DocNumero:
info.DocNumber` (Brazil's CPF/CNPJ — the direct equivalent of a US SSN/EIN for tax purposes), plus
a full billing address. The code is otherwise careful about PII — `redactBillingInfo` strips
`buyer.billing_info` from the *raw* ML payload before any raw capture is persisted
(`connectors/adapters/mercado_livre/order_ingest_reader.go:164-184`), `order_shipments_test.go:59-68`
enforces by test that shipments carry *zero* raw/PII columns, and the DTO comment at
`http_handler.go:393-395` explicitly names "C06/LGPD" (Brazil's data-protection law) as the reason
buyer identity is masked everywhere *except* this one fiscal block, which is deliberately
unmasked because a Brazilian nota fiscal legally requires the payer's real name and document
number. The gap is not carelessness about PII — it's that the one block designed to legitimately
carry unmasked PII is standing behind a `writeError`/`WriteJSON` handler with the exact same
"anyone who can reach the port" exposure as everything else in SEC-1.

**3. SEC-3 — the fail-open sibling the brief asked this lane to find, and it is bigger than the
one already fixed.** `cmd/catalogingest/main.go` (fixed per commit `47a76837`, `.mnfs/HARNESS-DEBTS.md`
lines ~1291-1429) now reads `os.Getenv("MC_DEFAULT_TENANT_ID")` directly and refuses to run if it
is empty. That guard is local to one manually-invoked batch command. `pgdb.LoadConfig()` —
`apps/server_core/internal/platform/pgdb/config.go:14-30` — is the loader `cmd/server/main.go`
calls to boot the actual always-on HTTP service, and it does the opposite: `MC_DATABASE_URL` and
`MPC_ENCRYPTION_KEY` both fail closed (`return Config{}, errors.New(...)`) when absent, but
`DefaultTenantID` silently becomes `"tenant_default"` (lines 23-24) and is returned as `nil` error
— the caller cannot distinguish "operator configured tenant_default on purpose" from "operator
forgot the variable." That single value then flows into `cfg.DefaultTenantID`, which
`internal/composition/root.go` threads into ~30 repository constructors (orders, catalog
enrichments, classifications, integrations installations, credentials, pricing, market, listings,
inventory — `grep -c "cfg.DefaultTenantID" internal/composition/root.go` → 70 occurrences). A
deployment that forgets this one variable does not crash; it runs correctly-looking and writes/reads
every tenant-scoped table under a fabricated identity — precisely the failure mode the operator
named as the reason this defect class matters, at ~30x the surface of the instance already fixed.

**4. SEC-4 — RLS's stated precondition is unmet by the only two things that could meet it.**
`migrations/0097_catalog_context.sql:92-108` writes a genuinely correct policy —
`USING (tenant_id = current_setting('app.tenant_id', true))`, `FORCE ROW LEVEL SECURITY` on all
four `catalog.*` tables (the migration's own comment names exactly the trap it's avoiding: "without
[FORCE] the table owner silently bypasses every policy"). But that policy has two independent
preconditions, and both fail in this deployment as measured: (a) the connecting role must not
bypass RLS — established fact 12 / D-44 says `marketplace` is superuser + bypassrls, and this
lane's own read of `docker-compose.yml:6-7,27` (`POSTGRES_USER: marketplace`,
`MC_DATABASE_URL: postgres://marketplace:marketplace@postgres...`) confirms the app pool has no
other choice today — `mpc_app` is not used anywhere in `internal/composition/root.go`; (b) the
session must `SET app.tenant_id` before the query runs, and the *only* place in the entire Go tree
that does this is `internal/contexts/catalog/internal/postgres/repository.go:49`, inside a
transaction helper on the new-context write path used by the manually-invoked
`cmd/catalogingest` — not by `internal/modules/catalog/transport/http_handler.go`, which is what
actually answers `/catalog/*` HTTP requests today. Even in the counterfactual where `mpc_app` were
wired in tomorrow, the live read path would set no GUC, and `tenant_id = NULL` matches nothing —
RLS would fail *closed* (breaking the feature, not leaking data) rather than leaking, which is the
better of two bad outcomes but still means the control does not do what its own migration comment
describes for the code that actually runs.

**5. SEC-5 — the RLS gap is a rounding error against the total surface.** Even granting SEC-4's
inertness a pass, RLS was only ever built for 4 tables. `grep -rl "ROW LEVEL SECURITY"
migrations/*.sql` returns exactly one file among 85. `grep -rl "tenant_id" migrations/*.sql | wc -l`
returns 52 — meaning 51 of the 52 migrations that introduce a `tenant_id` column ship with no
database-level enforcement mechanism at all, ever, by design (not a gap in an otherwise-covered
system — the system was never meant to have DB-level tenant enforcement outside the one context
that's mid-migration). This matches — and should be read as full agreement with, not a
contradiction of — established fact 12 and constraint 10's own framing: tenant isolation for the
other 51 tables is, today, entirely the hand-written `WHERE tenant_id = $1` in each query. This
lane did not find a query that omits it in a sample of `orders`, `listings`, and `mutations`
repositories (see Commands run), but did not attempt an exhaustive audit of all ~30 repository
files; that remains open (see Unverified).

## Control inventory

| control | exists? | fires? | contingent on what config | effect if absent |
|---|---|---|---|---|
| HTTP request authentication | no | — | — | every route open to network access (SEC-1) |
| HTTP request authorization (per-tenant/per-role) | no | — | — | n/a — no identity to authorize |
| Postgres RLS on `catalog.*` (4 tables) | yes (migration) | no | connecting role must lack BYPASSRLS **and** session must `SET app.tenant_id`; neither holds on the live app pool | tenant isolation on those 4 tables reduces to the hand-written query predicate, i.e. same as everywhere else (SEC-4) |
| Postgres RLS on the other ~51 tenant tables | no | — | — | same as above, by design, not a regression |
| Hand-written `WHERE tenant_id = $1` predicate | yes | yes (sampled) | correct call-site wiring of `cfg.DefaultTenantID`/`tenantID` param per repo | a missing predicate would be a full cross-tenant leak; not found in samples, not exhaustively audited |
| Tenant identity resolution | yes | yes | `MC_DEFAULT_TENANT_ID` env var | falls open to `"tenant_default"`, silently (SEC-3) |
| ML/Melhor Envio OAuth credential encryption at rest | yes | yes | `MPC_ENCRYPTION_KEY` (32-byte, AES-256-GCM) | `NewLocalKeyService` fails closed — `len(key)!=32` → error; `pgdb.LoadConfig` fails closed if unset |
| Oracle/Sankhya credential loading | yes | yes | `MPC_ORACLE_USERNAME/PASSWORD/CONNECT_STRING` | fails closed — `LoadConfigFromEnv` returns an error for any missing required field |
| Live provider write gate (`MPC_PROVIDER_WRITES_ENABLED`) | yes | yes | env var, must be literal `"true"` | fails closed — defaults to disabled, matches AGENTS.md "Live ML writes require explicit operator authorization" |
| Assisted Sankhya linkage feature gate | yes | yes | `MPC_ASSISTED_SANKHYA_LINKAGE_ENABLED` via `requiredBool` | fails closed — absent/unparseable → `ErrSankhyaLinkageRuntimeUnavailable` |
| Provider secret/token redaction in logs (ML responses) | yes | yes | none — always runs | `market_adapters.go:559-580`, `price_writer.go:128-129` regex-redact bearer/credential-shaped fields before logging; a benign body is asserted not over-redacted (`market_adapters_test.go:255-259`) |
| Panic-value redaction in the generic recover handler | no | — | — | raw `recover()` value goes to the log sink verbatim (SEC-7); never to the client |
| CORS origin restriction | no | — | — | `Access-Control-Allow-Origin: *` on every route, with `Authorization` in the allowed-headers list (SEC-6) |
| SQL parameterization | yes (sampled) | yes | none | dynamic `WHERE` clauses build placeholder *indices* via `fmt.Sprintf("...$%d...", n)`, values always travel through `args...`; no raw value interpolation found in `listings`/`mutations`/`orders` repositories sampled |

## Fail-open defaults

| env var | security property it changes | behaviour when absent |
|---|---|---|
| `MC_DEFAULT_TENANT_ID` | tenant identity for the entire server process (~30 repos) | silently becomes `"tenant_default"` — `apps/server_core/internal/platform/pgdb/config.go:23-24` (SEC-3) |
| `MC_DATABASE_URL` | which database/credentials the app uses | fails closed — `pgdb.LoadConfig` returns an error, server does not start |
| `MPC_ENCRYPTION_KEY` | whether OAuth credentials can be encrypted at rest | fails closed — `pgdb.LoadConfig` returns an error, server does not start |
| `MPC_ORACLE_USERNAME` / `MPC_ORACLE_PASSWORD` / `MPC_ORACLE_CONNECT_STRING` | Oracle/Sankhya read-path credentials | fails closed — `LoadConfigFromEnv` returns an error per missing field |
| `MPC_PROVIDER_WRITES_ENABLED` | whether the app is allowed to write to a live marketplace provider | fails closed — defaults to `false` (`strings.EqualFold(..., "true")`, anything else is falsy) |
| `MPC_ML_CATALOG_OFFERS_ENABLED` | whether the ML catalog-offers read route is reachable | fails closed — same `EqualFold(...,"true")` pattern, defaults off (ADR-032) |
| `MPC_ASSISTED_SANKHYA_LINKAGE_ENABLED` | whether the assisted-linkage Oracle read runtime activates | fails closed — `requiredBool` errors on empty/unparseable |

## What is actually fine

- **Credential encryption at rest for marketplace OAuth tokens is real and fails closed.**
  `LocalKeyService` (AES-256-GCM, `internal/modules/integrations/adapters/crypto/local_key.go`)
  refuses a key of the wrong length; `pgdb.LoadConfig` refuses to boot without one. Do not touch.
- **Oracle/Sankhya credential loading is fully fail-closed**, including numeric pool/timeout
  bounds (`oracle/config.go:99-137` validates every duration is positive and every pool bound is
  internally consistent). This is a well-built config loader; use it as the reference shape for
  fixing SEC-3.
- **PII redaction on the raw-capture path is deliberate and tested**, not accidental: `buyer.billing_info`
  is stripped from the raw ML order payload before persistence (`order_ingest_reader.go:164-184`),
  a dedicated test (`order_ingest_reader_test.go:18-22,190-200`) asserts a synthetic PII marker
  never survives, and shipment rows are asserted by test to carry zero raw/PII columns
  (`migrations/order_shipments_test.go:59-68`, citing "P7 r01 B-7, ADR-025"). SEC-2 is about the
  one deliberately-unmasked fiscal block being unauthenticated, not about the PII-handling
  discipline itself, which is above the bar seen elsewhere in this audit.
- **Provider secret redaction in logs/error surfaces for ML responses is real, tested both
  directions** (redacts credential-shaped substrings, asserted not to over-redact benign bodies —
  `market_adapters_test.go:232-259`).
- **The live-write gate defaults closed and matches the operator's own stated constraint**
  (AGENTS.md: "Live ML writes require explicit operator authorization").
- **SQL parameterization, sampled across `listings`, `mutations`, and `orders` postgres
  repositories, is correct** — dynamic `WHERE` fragments interpolate only positional parameter
  *indices*, never raw values; not exhaustively verified across all ~30 repositories (see
  Unverified).
- **`go.sum` and root `package-lock.json` exist** — dependency resolution is content-hash pinned
  for Go and lockfile-pinned for npm regardless of the 16 caret ranges in `apps/web/package.json`.

## Unverified / needs judgment

- **Whether `MC_DEFAULT_TENANT_ID` is actually set in the live dev `.env`.** This lane's sandbox
  denies reading `.env` directly (`Permission to read ... .env has been denied` — a correct guard,
  not a workaround target). SEC-3 is a code-level finding independent of the current `.env`
  content: the defect is that the loader has no way to fail loudly if it *were* unset, in
  production or any future deployment.
- **Whether the backend port is reachable from outside the operator's own machine/network in any
  real deployment.** `docker-compose.yml` publishes `8080:8080` for the dev stack; nothing in this
  repo describes a production topology, reverse proxy, or VPN boundary. PHASE-0.md marks repo
  visibility itself as `unverified`. SEC-1's severity is stated as "the code enforces nothing" —
  a network-layer compensating control may or may not exist outside this repository, and this lane
  has no way to check that from here.
- **Whether any repository beyond the sampled `listings`/`mutations`/`orders` ever omits the
  `tenant_id` predicate.** Not exhaustively audited across all ~30 postgres adapter files; the
  established-fact framing (constraint 10, D-44) already treats this predicate as the sole
  enforcement mechanism, which makes a single miss a full cross-tenant leak on that table. A
  dedicated instrument (grep every `postgres` adapter's `SELECT`/`UPDATE`/`DELETE` for a
  `tenant_id =` clause, flag any without one) would close this; this lane did not build it within
  budget.
- **Whether `internal/contexts/catalog`'s `Reader()` (fact 13, D-53: zero non-test callers) is
  meant to eventually replace the legacy `internal/modules/catalog` read path**, at which point
  the `set_config('app.tenant_id', ...)` call this lane found would start mattering for reads, not
  just the current writer-only use. This is a forward-looking observation, not a present-tense
  finding — recorded here so a future lane doesn't re-derive it.
- **Whether the `enrichedBuyerDTO`/masked-`Display` buyer name shown on other order endpoints is
  actually masked in practice** (this lane read the DTO's doc comment claim and the struct shape
  but did not trace the `Display` field's construction end to end the way SEC-2's `DocNumero` path
  was traced). Lower confidence than SEC-2; not asserted as a finding.

## Commands run

```
git rev-parse HEAD
git status --porcelain --untracked-files=all

grep -rl "ROW LEVEL SECURITY" apps/server_core/migrations/*.sql
grep -rn "FORCE ROW LEVEL SECURITY" apps/server_core/migrations/*.sql
grep -rln "CREATE POLICY" apps/server_core/migrations/*.sql
grep -rl "tenant_id" apps/server_core/migrations/*.sql | wc -l
ls apps/server_core/migrations/*.sql | wc -l

grep -rn "func.*Middleware|middleware\.|AuthN|AuthZ|Authenticate|Authorize|RequireAuth" apps/server_core/internal --include=*.go
cat apps/server_core/internal/platform/httpx/router.go
cat apps/server_core/cmd/server/main.go

grep -rn "TenantID|tenant_id|X-Tenant-Id|tenantFromContext|CtxTenant" apps/server_core/internal/platform/httpx
grep -n "tenant" -i apps/server_core/internal/composition/root.go
grep -rn "Authorization|Bearer |jwt\.|APIKey|X-Api-Key" apps/server_core/internal --include=*.go

cat apps/server_core/internal/platform/pgdb/config.go
cat apps/server_core/internal/platform/pgdb/tenant.go
grep -rn "pgdb\.DefaultTenantID" apps/server_core

sed -n '1,60p' docker-compose.yml
grep -n "MC_DEFAULT_TENANT_ID\|MC_DATABASE_URL\|MPC_ENCRYPTION_KEY" -r docker-compose.yml docker/

cat apps/server_core/internal/modules/integrations/adapters/crypto/local_key.go
cat apps/server_core/internal/modules/internal_read/adapters/oracle/config.go

grep -rn "buyer|cpf|email|phone|address" -i apps/server_core/migrations/0027_orders_marketplace_orders.sql
cat apps/server_core/migrations/0079_orders_buyer_nickname.sql
grep -rn "raw_provider_ref|RawProviderRef" apps/server_core/internal/modules/orders --include=*.go
cat apps/server_core/internal/modules/orders/adapters/postgres/buyer_fiscal_reader.go
sed -n '580,630p' apps/server_core/internal/modules/orders/transport/http_handler.go

grep -rn "scrub|redact|mask.*pii|PII" -i apps/server_core --include=*.go
grep -n "buyer_name|buyer_doc_number" apps/server_core/migrations/0089_orders_marketplace_orders_sync_fields.sql

grep -n "app\.tenant_id|set_config|SET LOCAL|SetTenant" -r apps/server_core/internal --include=*.go
sed -n '60,110p' apps/server_core/migrations/0097_catalog_context.sql

grep -n "runtime.Handler|CORSMiddleware|apierror.Recover" apps/server_core/internal/composition/root.go
grep -rn "mux\.HandleFunc(|mux\.Handle(" apps/server_core/internal/modules/*/transport/*.go | wc -l
grep -rln "func.*Middleware|RequireAuth|checkAuth|verifyToken|apiKeyMiddleware" apps/server_core/internal/modules/*/transport/*.go apps/server_core/internal/platform/httpx/*.go

grep -rn "headers:|Authorization|X-Api-Key" apps/web/src --include=*.ts*

head -5 apps/server_core/go.mod
ls apps/server_core/go.sum
grep -c '"\^' apps/web/package.json
ls package-lock.json

grep -rn 'fmt\.Sprintf\(`|fmt\.Sprintf\("SELECT|"UPDATE|"INSERT|"DELETE' apps/server_core/internal --include=*.go
sed -n '90,140p' apps/server_core/internal/modules/listings/adapters/postgres/repository.go

cat apps/server_core/internal/platform/config/*.go
sed -n '990,1010p' apps/server_core/internal/composition/root.go
cat apps/server_core/internal/modules/internal_read/adapters/oracle/sankhya_linkage_config.go
```

Read in full: `docs/engineering/repo-audit-2026-08-07/PHASE-0.md`, `.mnfs/HARNESS-DEBTS.md`
(sections around D-2, D-9, D-28, D-31/32/35, D-40/41/42, D-43 through D-50, grepped for
tenant/secret/auth/PII keywords rather than read linearly — file is 1846 lines).
