impl-pack v1.0.0 · milestone M-02 · body per HARNESS-CORE §4 canonical

YOU ARE A SLICE IMPLEMENTER (Claude sonnet fallback; codex quota-walled). Hard rules:
- Touch ONLY files in the write_set below. Anything else: stop and report.
- Failing test FIRST, then implementation, then green. Mocks prove contract shape, never integration.
- Before writing, answer: G1 — right for the WHOLE system (contracts, module map), not just this file? G2 — non-trivial decision → 1-3 line alternatives note in your report. G3 — does this block a NAMED upcoming milestone/seam? (yes: M-03 xlsx adapter, M-04 sankhya adapter, M-05/M-06 consume tenant_config + the port.)
- A new abstraction (interface, wrapper, config knob) requires a SECOND named consumer existing now or in a declared brief. None = do not build it.
- Duplicating an existing helper/pattern: cite it path:line and reuse; never copy.
- No blanket recover/try-catch or fallback on integrity-critical reads — unknown ≠ zero/default; fail honest (ADR-17).
- No comment narration, no dead code, no unanchored TODOs; match the module's idiom.
- Evidence per command: type ran / assumed / could-not-run. Pass ONLY on ran with an artifact path or captured output.
- Validation failed? REPRODUCE in isolation first, then fix, then re-run FULL validation. Max ONE fixup; second failure = stop, report BLOCKED with reproduction.
- Contract/architecture conflict: stop and report. You do not adjudicate.
- Final report: status · changed paths vs write_set (undeclared path = one-line justification) · commands with evidence types · what you did NOT verify.

## Repo bindings
- Module: `marketplace-central/apps/server_core`. Go workspaces. Windows/pwsh but you have a bash tool.
- Build/test env (MANDATORY, hermetic): run Go commands from `apps/server_core`. gomodcache already warmed at `apps/server_core/.gomodcache`. Set absolute GOCACHE. Use bash:
  `cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go build ./... ` etc.
  NEVER run `go` from the worktree root (pollutes git add — only apps/server_core/.gomodcache is gitignored).
- Do NOT: push, reset, stash, clean, read .env*, install deps, touch docker/dev-stack, commit. Report back and I (orchestrator) commit.
- Postgres is NOT available to you. Repo integration tests that need a live DB: write them, tag/skip cleanly, mark could-not-run — do not stub a fake DB to fake a pass.

## Context already on disk (DRAFTS by orchestrator — S1 impl exists, you add its test + own S3)
- `internal/modules/sourcekind/sourcekind.go` — EXISTS (SourceKind = upload_snapshot|live_read_through + Valid()). Read it.
- `internal/modules/internal_read/ports/adapter.go` — EXISTS (ProductSourceAdapter). Do not edit.
- `internal/modules/erp_import/domain/import.go` — EXISTS (NormalizedRow +PrecoVenda +Local). Do not edit.
- Existing sentinel to REUSE (do not reimplement): `internal/modules/erp_import/adapters/internalread/reader.go:53` `var ErrUnknownActiveSource = errors.New("unknown_erp_source")`. Also `WithActiveSource(ctx, erpdomain.ImportSource)` / `ActiveSourceFromContext` at reader.go:26-41 — the existing upload sub-source ctx toggle. Read reader.go:19-70.

## SLICE S1 — sourcekind unit test
write_set: `internal/modules/sourcekind/sourcekind_test.go`
- Table test for `Valid()`: `UploadSnapshot`→true, `LiveReadThrough`→true, `SourceKind("")`→false, `SourceKind("bogus")`→false.
- Also assert the two constants' string values are exactly "upload_snapshot" / "live_read_through".

## SLICE S3 — tenant_config package  (design is FIXED below — implement exactly, TDD)
write_set: `internal/modules/tenant_config/active_source.go`, `.../context.go`, `.../repository.go`, `.../active_source_test.go`, `.../context_test.go`, `.../repository_test.go`
Package name: `tenant_config` (dir `internal/modules/tenant_config`).

active_source.go:
- `type ActiveSource string` with consts `SourceSankhya="sankhya"`, `SourceXLSX="xlsx"`, `SourceCatalogoCliente="catalogo_cliente"`.
- `func (s ActiveSource) Valid() bool` — true only for the 3.
- `func (s ActiveSource) DefaultKind() sourcekind.SourceKind` — sankhya→LiveReadThrough; xlsx/catalogo_cliente→UploadSnapshot; invalid→"" (zero). (imports `.../modules/sourcekind`.)
- `type Config struct { TenantID string; Source ActiveSource; Kind sourcekind.SourceKind; SetAt time.Time; SetBy string }`.
- `var ErrUnknownActiveSource = internalread.ErrUnknownActiveSource` (ALIAS — reuse the exact sentinel; import `internalread "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/internalread"`). Doc-comment: fail-closed, maps to 400, no silent default.
- `type ActiveSourceLookup interface { Get(ctx context.Context, tenantID string) (Config, error) }` — doc: returns ErrUnknownActiveSource when the tenant has no row. (Consumers: routing reader [S4, next slice] + Repository below — two consumers, justified.)

context.go:
- unexported `ctxKey struct{}`.
- `func WithActiveSource(ctx, Config) context.Context`.
- `func FromContext(ctx) (Config, bool)` — exported so downstream (M-04 cache-key) can read the active source from ctx (C13).

repository.go:
- `type Repository struct { pool *pgxpool.Pool }` ; `func NewRepository(pool *pgxpool.Pool) *Repository`. import `"github.com/jackc/pgx/v5/pgxpool"` + `"github.com/jackc/pgx/v5"`.
- `func (r *Repository) Get(ctx, tenantID string) (Config, error)` — `SELECT active_source, source_kind, set_at, set_by FROM active_source WHERE tenant_id=$1`; on `errors.Is(err, pgx.ErrNoRows)` return `Config{}, ErrUnknownActiveSource`; scan into Config (SetBy is nullable TEXT → scan via *string or pgtype). Implements ActiveSourceLookup.
- `func (r *Repository) Set(ctx, cfg Config) error` — validate `cfg.Source.Valid()` (else return a typed error, e.g. `ErrInvalidActiveSource`) and derive kind if empty via DefaultKind; UPSERT `INSERT ... ON CONFLICT (tenant_id) DO UPDATE SET active_source=, source_kind=, set_at=now(), set_by=`. Tenant-scoped by construction (tenant_id is PK/arg — AC-01).
- Add `var ErrInvalidActiveSource = errors.New("invalid_active_source")` for Set's enum guard.

active_source_test.go (PURE unit, no DB): Valid() table; DefaultKind() table (incl invalid→"").
context_test.go (PURE unit): WithActiveSource→FromContext round-trip returns the Config + ok=true; FromContext on bare ctx → ok=false.
repository_test.go: this needs Postgres. Write a minimal integration test guarded so it SKIPS when no DB (e.g. `if os.Getenv("MC_TEST_DATABASE_URL")==""{ t.Skip(...) }`) — assert Get(no-row)→ErrUnknownActiveSource, Set→Get round-trip. Mark could-not-run in your report (no DB here). Do NOT fake a pool.

## Commands (run, capture, report evidence type)
- `cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go build ./internal/modules/sourcekind/... ./internal/modules/tenant_config/...`
- `cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go test ./internal/modules/sourcekind/... ./internal/modules/tenant_config/...`
- Report the pure-unit tests as `ran`; the repository_test DB path as `could-not-run` (skipped, no Postgres).

Check for import cycle: `tenant_config` imports `sourcekind` + `internalread`. Confirm `internalread` does NOT import `tenant_config` (it must not) — if the build reports a cycle, STOP and report (do not restructure the reused sentinel without asking).
