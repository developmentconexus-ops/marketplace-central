impl-pack v1.0.0 · milestone M-02 · body per HARNESS-CORE §4 canonical

YOU ARE A SLICE IMPLEMENTER (Claude sonnet fallback; codex quota-walled). Hard rules:
- Touch ONLY files in the write_set below. Anything else: stop and report.
- Failing test FIRST, then implementation, then green. Mocks prove contract shape, never integration.
- Before writing, answer: G1 — right for the WHOLE system? G2 — non-trivial decision → 1-3 line alternatives note. G3 — blocks a NAMED seam? (yes: M-04 extends the cache key from the ctx active-source you inject here; M-03/M-04 implement the adapters this wiring builds.)
- A new abstraction requires a SECOND named consumer now or in a declared brief. None = do not build it.
- Duplicating an existing helper/pattern: cite it path:line and reuse; never copy.
- No blanket recover/try-catch or fallback on integrity-critical reads — unknown ≠ zero/default; fail honest (ADR-17). A tenant with no active-source row MUST surface the error, never silently read xlsx.
- No comment narration, no dead code, no unanchored TODOs; match the module's idiom.
- Evidence per command: ran / assumed / could-not-run. Pass ONLY on ran with captured output.
- Validation failed? REPRODUCE in isolation first, then fix, re-run FULL validation. Max ONE fixup; second failure = BLOCKED with reproduction.
- Contract/architecture conflict: stop and report. You do not adjudicate.
- Final report: status · changed paths vs write_set (undeclared = one-line justification) · commands with evidence types · what you did NOT verify.

## Repo bindings
- Module `marketplace-central/apps/server_core`. Run Go from `apps/server_core`; caches warmed.
  `cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go build ./...`
  NEVER run go from worktree root. Do NOT push/reset/stash/clean/commit/read .env/install deps/touch dev-stack. Leave changes for the orchestrator to commit.
- Working dir: C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\trusting-mayer-a5c8f6

## Prereqs already on disk (from prior slice — READ, do not edit)
- `internal/modules/tenant_config/active_source.go`: `ActiveSource` {SourceSankhya,SourceXLSX,SourceCatalogoCliente}, `Config{TenantID,Source,Kind,SetAt,SetBy}`, `ActiveSourceLookup interface { Get(ctx, tenantID string) (Config, error) }` (returns `tenant_config.ErrUnknownActiveSource` on no row), `WithActiveSource(ctx,Config)`, `FromContext(ctx)`.
- `internal/modules/tenant_config/repository.go`: `NewRepository(pool *pgxpool.Pool) *Repository` implements ActiveSourceLookup.
- `internal/modules/erp_import/adapters/internalread/reader.go:26-48`: `WithActiveSource(ctx, erpdomain.ImportSource) context.Context` — the EXISTING upload sub-source ctx toggle. `erpdomain.SourceXLSX="xlsx"`, `erpdomain.SourceCatalogoCliente="catalogo_cliente"` (erp_import/domain/import.go:64-67).
- `internal/modules/internal_read/ports/reader.go:48-55`: `Reader` interface — 6 methods (FindProductsForLinking, GetSellableStock, GetCurrentPrice, GetCostAsOf, GetSalesHistory, GetTaxInputs). Signatures MUST NOT change.

## SLICE S4 — routing reader + root.go refactor

### Part A — routing reader (NEW pkg)
write_set: `internal/modules/internal_read/adapters/routing/reader.go`, `internal/modules/internal_read/adapters/routing/reader_test.go`
Package `routing`. Type `Reader` implementing `internalreadports.Reader`:
```
type Reader struct {
    upload   internalreadports.Reader   // xlsx/catalogo_cliente path (erp internalread, cached)
    live     internalreadports.Reader   // sankhya/oracle path, may be nil if oracle unconfigured
    lookup   tenant_config.ActiveSourceLookup
    tenantID string                     // single-tenant resolution key (cfg.DefaultTenantID)
}
func NewReader(upload, live internalreadports.Reader, lookup tenant_config.ActiveSourceLookup, tenantID string) *Reader
```
Core resolution (private helper), used by EVERY method:
```
func (r *Reader) resolve(ctx) (internalreadports.Reader, context.Context, error) {
    cfg, err := r.lookup.Get(ctx, r.tenantID)   // ErrUnknownActiveSource propagates as-is (fail-closed)
    if err != nil { return nil, ctx, err }
    ctx = tenant_config.WithActiveSource(ctx, cfg)   // downstream/cache-key prereq (M-04)
    switch cfg.Source {
    case tenant_config.SourceXLSX, tenant_config.SourceCatalogoCliente:
        ctx = internalread.WithActiveSource(ctx, erpdomain.ImportSource(cfg.Source))  // upload sub-toggle
        return r.upload, ctx, nil
    case tenant_config.SourceSankhya:
        if r.live == nil {
            return nil, ctx, <typed unavailable error>   // NOT a fallback to upload
        }
        return r.live, ctx, nil
    default:
        return nil, ctx, tenant_config.ErrUnknownActiveSource   // defensive; lookup shouldn't return invalid
    }
}
```
For the sankhya-but-live-nil case define a package error, e.g. `var ErrActiveSourceUnavailable = errors.New("active_source_unavailable")` (doc: the configured live source has no reader wired — fail honest, never serve the other source's data).
Each of the 6 Reader methods: `rd, ctx, err := r.resolve(ctx); if err != nil { return <zero of return type>, err }; return rd.MethodName(ctx, input)`. Return-zero types: FindProductsForLinking→`nil`; the Get* ones return domain structs (`internalreaddomain.SellableStock{}` etc.) — check reader.go for exact types.

reader_test.go (PURE unit, no DB) — fakes only:
- fake `internalreadports.Reader` recording the ctx it was called with (two instances: upload, live).
- fake `tenant_config.ActiveSourceLookup` returning a scripted (Config,error).
Cases (drive via FindProductsForLinking, one representative method is enough + assert others compile):
  1. lookup→{Source:xlsx} : upload reader invoked; its ctx has erp `ActiveSourceFromContext`==xlsx AND `tenant_config.FromContext` ok==true Source==xlsx. live NOT invoked. (proves BOTH-adapters wired + correct dispatch, C11)
  2. lookup→{Source:catalogo_cliente} : upload reader invoked, erp ctx==catalogo_cliente.
  3. lookup→{Source:sankhya}, live present : live reader invoked; tenant_config.FromContext Source==sankhya; erp active-source NOT set.
  4. lookup→{Source:sankhya}, live==nil : returns ErrActiveSourceUnavailable; NEITHER reader invoked (no silent upload fallback).
  5. lookup→ErrUnknownActiveSource : method returns that error; neither reader invoked. (C12 fail-closed)

### Part B — root.go source-wiring refactor
write_set: `internal/composition/root.go`
This is the orchestrator-OWNED source-wiring section. Do NOT touch any ticker/scheduler block (that is another track's additive-lock). Edits:
1. DELETE the `erpSource` function (currently ~line 772-782) and its call `source, err := erpSource(os.Getenv)` (~line 271-274). Remove now-unused imports if any (`strings` may become unused — check).
2. DELETE the `MC_ERP_SOURCE` reference entirely. After this, `grep -rn MC_ERP_SOURCE apps/server_core` MUST return 0.
3. Replace the single-source branch (currently ~lines 415-453 `var internalReadSvc ...; if source == "xlsx" {...} else if oracleCfg ... {...}`) with: build BOTH readers unconditionally-where-possible, then a routing reader.
   - Build `uploadReader` ALWAYS: `erpinternalread.NewReader(erpRepo, cfg.DefaultTenantID)` → wrap in `internalreadcache.NewReader(_, freshnessCache)` → wrap in `internalreadobservability.NewTimingReader(_, slog.Default(), observabilityCfg.SlowQueryThreshold)`. (Same wrapping the old xlsx branch used.)
   - Build `liveReader` (may stay nil): keep the existing oracle path (`internalreadoracle.LoadConfigFromEnv` / `OpenDB` — same warn-on-failure graceful degrade). On success set `liveReader` = the cached+timed oracle reader, keep `oracleDB`, `poolStats`, and the `inventoryStockReader` oracle wiring EXACTLY as today (guarded by oracle availability). On failure liveReader stays nil (app still serves the upload source).
   - `activeSourceRepo := tenant_config.NewRepository(pool)`.
   - `routingReader := routing.NewReader(uploadReader, liveReader, activeSourceRepo, cfg.DefaultTenantID)`.
   - `internalReadSvc = internalreadapp.NewService(routingReader)`; `internalReadAvailable = true` (xlsx always available now); `productMatcher = internalReadSvc`.
   - Everything downstream (canonicalCatalogReader, catalogPageReader, etc.) is UNCHANGED — it consumes internalReadSvc as before.
   NOTE: when `pool == nil` (some test/degraded boots), NewRepository(nil) would panic on use. Preserve existing nil-pool tolerance: if `pool == nil`, keep the old unavailable-reader behavior (internalReadAvailable=false, productMatcher=unavailable). Guard the routing construction behind `pool != nil`. Read the surrounding code (lines ~405-465) and preserve its nil-safety exactly.
4. Add imports: `routing "…/internal_read/adapters/routing"`, `tenant_config "…/internal/modules/tenant_config"`. `erpinternalread`, `internalreadcache`, `internalreadobservability`, `internalreadoracle`, `internalreadapp` already imported.

## Commands (run, capture, report)
- `cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go build ./...`  (C6 full build; expect green)
- `cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go vet ./internal/composition/... ./internal/modules/internal_read/adapters/routing/...`
- `cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go test ./internal/modules/internal_read/adapters/routing/...`  (routing unit GREEN)
- `grep -rn MC_ERP_SOURCE apps/server_core` → MUST be 0 lines (report the count).
- If any existing test in `internal/composition` or `internal/modules/internal_read` breaks from the wiring change, REPRODUCE, then report — do not broaden scope beyond making the refactor build+pass. If a broken test asserts the OLD env-branch behavior, STOP and report it as a finding (it may need orchestrator decision), do not silently rewrite it beyond the mechanical wiring change.

What you did NOT verify: live product reads against Postgres/oracle (no DB here) — note it.
