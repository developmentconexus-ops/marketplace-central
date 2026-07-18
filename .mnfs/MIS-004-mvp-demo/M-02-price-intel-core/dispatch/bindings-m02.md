## Repo & workspace bindings (M-02 price-intel-core)

- Repo: marketplace-central (Go backend `apps/server_core`, hexagonal modules under `internal/modules/`). Windows + PowerShell.
- Worktree root (your CWD): `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core`. Work ONLY here.
- Branch `chip/m02-price-intel-core`, base SHA `59d0e62fdbf15db068542432ef5d5731b6fa9f83`.

## Build / test commands (run from worktree root unless noted; GOCACHE MUST be absolute)

- Build:  `$env:GOCACHE="<worktree>\apps\server_core\.gocache"; cd apps\server_core; go build -buildvcs=false ./...`
- Vet:    `$env:GOCACHE="<worktree>\apps\server_core\.gocache"; cd apps\server_core; go vet ./...`
- **`-buildvcs=false` is REQUIRED in the codex workspace-write sandbox** — the sandbox blocks git, so plain `go build ./...` fails with git `exit status 128` (VCS stamping). This is a FALSE ALARM, not a code defect; the flag disables revision embedding for local L0 only (governance/close lanes run in a clean git-enabled worktree). Field-verified F-02-S1 2026-07-17.
- Test (touched pkg): `$env:GOCACHE="<worktree>\apps\server_core\.gocache"; $env:GOMODCACHE="<worktree>\apps\server_core\.gomodcache"; cd apps\server_core; go test ./internal/modules/<yourpkg>/...`
- `.gomodcache` is already warmed. If a build dies with `HPG_MIGRATION_FAILED migrations_first=-1`, the modcache is cold — re-run `GOMODCACHE=<abs> go mod download all` in apps/server_core; it is NOT a SQL defect.
- Your slice card names the EXACT commands + validation_kind. Run them; capture output. Evidence type per command: ran / assumed / could-not-run. Pass ONLY on `ran` with captured output.

## Reuse these existing seams (do NOT reinvent; cite path:line when you build on them)

- ML adapter HTTP+auth: `internal/modules/connectors/adapters/mercado_livre/capability_adapter.go` — GET via `doJSON(ctx, accountRef, token, http.MethodGet, path, nil, &out)` (:522-542 sets Authorization Bearer + X-Tenant-ID + X-Installation-ID); token via `accessToken(ctx, accountRef)` (:452-464); `readItem` (:59-62) is the GET-/items/{id} template. New read ports are ADDITIVE methods on `CapabilityAdapter` reusing these.
- Credential flow (already wired): `internal/modules/integrations/application/credential_resolver.go` `ResolveAccessToken(ctx, CredentialResolutionRef{TenantID,InstallationID,ProviderCode,ProviderAccountID})`. Do NOT add a new token path.
- Connectors ports idiom: interface in `internal/modules/connectors/ports/*.go` (e.g. marketplace_capability.go:22-29), normalized DTO struct with `FetchedAt time.Time` in `internal/modules/connectors/domain/capability.go` (e.g. AccountSnapshot :145-153).
- market module template: `internal/modules/market/adapters/postgres/repository.go` (`Repository{pool,tenantID}` :5-12); query pattern `WHERE tenant_id=$1 AND ... = ANY($2)` + `SELECT DISTINCT ON (...) ORDER BY ... captured_at DESC` (observation_repository.go:57-72); `domain/market.go` `Money{Amount string, Currency string}` (decimal STRING, validated by `decimalPattern`; NOT cents, NOT shopspring/decimal).
- Migrations: forward-only `*.sql` via `embed.FS` (`migrations/source.go`); applied lexicographically, no down path. 0044 PK `(tenant_id, product_id, captured_at)`. Fixture: `internal/platform/migrate/runner_test.go:25` asserts `if len(want) != 41` — bump this integer by exactly the number of new migration files you add.
- Composition root: `internal/composition/root.go` — market wired at :531-533, connectors capability at :330-347. Add ONLY your own registration lines near these; never edit another module's wiring.
- Telemetry: repo has NO prometheus/metrics lib. Structured `slog` only — template `internal/modules/internal_read/observability/timing.go:124-136` (logger.Info/Error with key/value attrs). Flag-route "counter + status" telemetry = slog fields (route, status, page, duration_ms, count). Do NOT add a metrics dependency.
- Time round-trip in tests: `time.Now().UTC().Truncate(time.Microsecond)` (postgres µs precision).

## Non-negotiables (reject-on-violation)

- tenant_id predicate on EVERY tenant query (`WHERE tenant_id=$1 ...`).
- Provider (ML) DTOs die at the connectors adapter. Consumers see ONLY normalized IC-06 shapes. No raw ML JSON in logs.
- Unknown NEVER becomes zero/default (ADR-17). Nullable stays null/pointer-nil. A FAILED collection NEVER overwrites the last VALID snapshot — write a new FAILED row, keep the VALID as latest-valid. No blanket recover/try-catch or fallback on integrity reads — fail honest with a typed error.
- Read-only ML: ZERO PUT/POST to Mercado Livre. FORBIDDEN to edit `listing_writer.go` / `price_writer.go` (the existing PUT /items/{id} write path). FORBIDDEN route `/sites/MLB/search` (403).
- Flag `MC_ML_CATALOG_OFFERS_ENABLED` default OFF; OFF or route-fail ⇒ typed `ErrCatalogOffersUnavailable`. `ListCatalogOffers` pagination MUST be complete (partial page = typed error, NEVER silent truncation).
- NO new dependencies. A dep change is out of your scope — stop and report.
- Contract/architecture conflict, or a needed change outside your write_set ⇒ STOP and report. You do not adjudicate or widen scope.

## Global forbidden paths (any slice)

`.env*` (never read/print/commit) · `contracts/governance/**` · `scripts/harness*` · `docs/HARNESS*` · `packages/sdk-runtime/src/index.ts` (barrel — hub) · other modules not in your write_set · migration numbers outside 0050–0054.
