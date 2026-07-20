impl-pack v1.0.0 · milestone M-02 · body per HARNESS-CORE §4 canonical

YOU ARE A SLICE IMPLEMENTER (Claude sonnet fallback; codex quota-walled). Hard rules:
- Touch ONLY files in the write_set below. Anything else: stop and report.
- Failing test FIRST, then implementation, then green. Mocks prove contract shape, never integration.
- Before writing: G1 right for whole system? G2 non-trivial decision → 1-3 line alt note. G3 blocks a NAMED seam? (yes: M-06 UI /importacoes reads+writes this endpoint; sdk types are the contract it binds to.)
- New abstraction needs a SECOND named consumer now/declared. None = don't build it.
- Duplicating a helper/pattern: cite path:line and reuse; never copy.
- Integrity reads fail HONEST: a tenant with no active-source row surfaces the error (400), never a silent default (ADR-17).
- No comment narration, no dead code, no unanchored TODOs; match module idiom.
- Evidence per command: ran / assumed / could-not-run. Pass ONLY on ran with captured output.
- Validation failed? REPRODUCE isolated, fix, re-run FULL. Max ONE fixup; second fail = BLOCKED w/ reproduction.
- Contract/architecture conflict: stop and report. You don't adjudicate.
- Final report: status · changed paths vs write_set · commands w/ evidence type · what you did NOT verify.

## Repo bindings
- Go module `marketplace-central/apps/server_core`; run go from `apps/server_core`:
  `cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go build ./...`
  NEVER run go from worktree root. Node/vitest: run from repo root (`packages/sdk-runtime` has its own vitest).
- Working dir: C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\trusting-mayer-a5c8f6
- Do NOT: push/reset/stash/clean/commit/read .env*/install deps/touch dev-stack. Leave changes for orchestrator.

## Prereqs on disk (READ, do not edit)
- `internal/modules/tenant_config/active_source.go`: `ActiveSource` {SourceSankhya="sankhya",SourceXLSX="xlsx",SourceCatalogoCliente="catalogo_cliente"}, `Valid()`, `DefaultKind()`, `Config{TenantID,Source ActiveSource,Kind sourcekind.SourceKind,SetAt time.Time,SetBy string}`, `ErrUnknownActiveSource` (sentinel, msg "unknown_erp_source").
- `internal/modules/tenant_config/repository.go`: `NewRepository(pool) *Repository` with `Get(ctx,tenantID)(Config,error)` (ErrNoRows→ErrUnknownActiveSource) and `Set(ctx,cfg)error` (invalid source→`ErrInvalidActiveSource`, derives kind if empty, UPSERT).
- Transport idiom to MIRROR (do not copy blindly — cite): `internal/modules/erp_import/transport/http_handler.go` — `Handler` struct + `NewHandler` + `Register(mux httpx.RouteRegistrar)` using Go-1.22 method-prefixed patterns ("GET /erp/imports"), `writeError(w,status,code,message)` at :154 wrapping `httpx.WriteJSON(w,status,payload)` (payload `map[string]string{"error":code[,"detail":]}`), interactive-route registration via `registerInteractiveRoute`. Response DTOs are local structs with json tags.
- sdk-runtime (hand-written TS, NO codegen): `packages/sdk-runtime/src/` — modules like `erpImport.ts` export `interface`/`type`; `index.ts` re-exports via `export * from "./erpImport"`. Each module has a sibling `*.test.ts`.
- OpenAPI: `contracts/api/marketplace-central.openapi.yaml` — `paths:` @ line 5, `/erp/imports:` block @ ~3083 (mirror its operation/response shape), `components: schemas:` @ ~3198, `ErpImportError:` schema @ ~7915.

## SLICE S5 — active-source transport + OpenAPI + sdk-runtime  (ALL IN ONE LOGICAL CHANGESET — contract spec + generated/hand-authored SDK MUST land together, C15)

### Endpoints (tenant-implicit, tenant from cfg.DefaultTenantID — SAME convention as /erp/imports)
- `GET /config/active-source` → 200 `{active_source, source_kind, set_at, set_by}`; if tenant has no row → 400 `{error:"unknown_erp_source"}` (fail-closed, map `tenant_config.ErrUnknownActiveSource`).
- `PUT /config/active-source` body `{active_source: "sankhya"|"xlsx"|"catalogo_cliente", set_by?: string}` → 200 the stored config (same shape as GET); invalid/absent active_source → 400 `{error:"invalid_active_source"}` (map `tenant_config.ErrInvalidActiveSource`); malformed JSON → 400 `{error:"invalid_body"}`. `source_kind` is DERIVED server-side (never client-supplied) via `DefaultKind()`.

### Part A — Go transport
write_set: `internal/modules/tenant_config/transport/http_handler.go`, `internal/modules/tenant_config/transport/http_handler_test.go`
- Package `transport`. Define a narrow port the handler depends on (TWO consumers justify it: this handler + tests):
  `type Store interface { Get(ctx, tenantID string) (tenant_config.Config, error); Set(ctx, cfg tenant_config.Config) error }`
  (`*tenant_config.Repository` satisfies it — add `var _ Store = (*tenant_config.Repository)(nil)` in the composition wiring or a compile assert in this pkg importing tenant_config.)
- `type Handler struct { store Store; tenantID string }` ; `func NewHandler(store Store, tenantID string) Handler`.
- `func (h Handler) Register(mux httpx.RouteRegistrar)` — register both as INTERACTIVE routes (mirror `registerInteractiveRoute`; these are cheap single-row ops, NOT batch). Patterns `"GET /config/active-source"`, `"PUT /config/active-source"`.
- GET handler: `cfg, err := h.store.Get(r.Context(), h.tenantID)`; on `errors.Is(err, tenant_config.ErrUnknownActiveSource)` → writeError 400 "unknown_erp_source"; other err → 500 "internal_error"; else 200 the response DTO.
- PUT handler: decode `{active_source string; set_by *string}` (json.NewDecoder; on decode err → 400 "invalid_body"). Build `Config{TenantID:h.tenantID, Source:tenant_config.ActiveSource(body.active_source), SetBy: deref-or-empty}`. `h.store.Set(...)`; on `errors.Is(err, tenant_config.ErrInvalidActiveSource)` → 400 "invalid_active_source"; other err → 500. On success re-Get (or reuse Set-derived cfg) and return 200 the stored config. Prefer re-Get so set_at/source_kind reflect the persisted row.
- Response DTO local struct: `{active_source string; source_kind string; set_at string (RFC3339); set_by string}` json tags snake_case. Reuse `writeError` idiom + `httpx.WriteJSON`.
- http_handler_test.go (PURE unit, fake Store — no DB): drive via httptest.
  1. GET, store returns Config{xlsx,upload_snapshot,...} → 200, body has active_source=xlsx source_kind=upload_snapshot.
  2. GET, store returns ErrUnknownActiveSource → 400 body error=unknown_erp_source (C12 fail-closed at transport).
  3. PUT {active_source:"sankhya"} → Set called with Source=sankhya; 200; (fake Set records cfg; fake Get after returns sankhya/live_read_through).
  4. PUT {active_source:"garbage"} → fake Set returns ErrInvalidActiveSource → 400 error=invalid_active_source. (Prove server rejects; do NOT rely on client.)
  5. PUT malformed JSON body → 400 error=invalid_body.
  Assert the handler NEVER writes a source_kind taken from the request (derive-only) — e.g. PUT with an extra source_kind field in JSON is ignored.

### Part B — root.go registration (SINGLE glue line, orchestrator-owned file)
write_set: `internal/composition/root.go`
- Inside the existing `if pool != nil { ... }` block (where `activeSourceRepo := tenant_config.NewRepository(pool)` and `routingReader` are wired — around line 447-451), add:
  `tenantconfigtransport.NewHandler(activeSourceRepo, cfg.DefaultTenantID).Register(mux)`
  reusing the already-constructed `activeSourceRepo` (do NOT construct a second repo). Add import `tenantconfigtransport "…/internal/modules/tenant_config/transport"`.
- Do NOT touch anything else in root.go. If activeSourceRepo is not in scope where you need it, register right after `routingReader` wiring inside the same block. Read lines ~439-452 first.

### Part C — OpenAPI (mirror /erp/imports block shape)
write_set: `contracts/api/marketplace-central.openapi.yaml`
- Add path `/config/active-source` with `get` (operationId `getActiveSource`, 200 → `$ref '#/components/schemas/ActiveSourceConfig'`, 400 → `$ref '#/components/schemas/ActiveSourceError'`) and `put` (operationId `setActiveSource`, requestBody `$ref '#/components/schemas/SetActiveSourceRequest'`, 200 → ActiveSourceConfig, 400 → ActiveSourceError). Place it alphabetically/adjacent near other `/config`-ish or top-of-paths per file ordering — match surrounding indentation EXACTLY (2-space).
- Add schemas under `components: schemas:` :
  - `ActiveSourceConfig`: object, required [active_source, source_kind, set_at], props active_source (enum sankhya|xlsx|catalogo_cliente), source_kind (enum upload_snapshot|live_read_through), set_at (string date-time), set_by (string, nullable).
  - `SetActiveSourceRequest`: object, required [active_source], props active_source (same enum), set_by (string, nullable). NO source_kind (server-derived — document that in the description).
  - `ActiveSourceError`: object, required [error], props error (enum unknown_erp_source|invalid_active_source|invalid_body|internal_error), detail (string, optional). (Mirror ErpImportError @ ~7915.)

### Part D — sdk-runtime (hand-authored TS — the SDK half of C15, SAME changeset)
write_set: `packages/sdk-runtime/src/activeSource.ts`, `packages/sdk-runtime/src/activeSource.test.ts`, `packages/sdk-runtime/src/index.ts`
- activeSource.ts:
  `export type ActiveSourceName = "sankhya" | "xlsx" | "catalogo_cliente";`
  `export type SourceKindName = "upload_snapshot" | "live_read_through";`
  `export interface ActiveSourceConfig { active_source: ActiveSourceName; source_kind: SourceKindName; set_at: string; set_by: string | null; }`
  `export interface SetActiveSourceRequest { active_source: ActiveSourceName; set_by?: string | null; }`
  `export interface ActiveSourceError { error: "unknown_erp_source" | "invalid_active_source" | "invalid_body" | "internal_error"; detail?: string; }`
  Doc-comment source_kind as server-derived (not in the request), mirroring erpImport.ts:15-21 style.
- index.ts: add `export * from "./activeSource";` (match existing re-export lines).
- activeSource.test.ts: type-level assertions mirroring the pattern in an existing `*.test.ts` (e.g. erpImport.test.ts / index.test.ts) — assert the enum members compile and a well-formed literal satisfies each interface; assert `SetActiveSourceRequest` has NO source_kind (a literal with source_kind is a type error → use `// @ts-expect-error`).

## Commands (run, capture, report)
- `cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go build ./...`
- `cd apps/server_core && GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache go test ./internal/modules/tenant_config/... ./internal/composition/...`
- sdk-runtime tests: from repo root, `npx vitest run packages/sdk-runtime/src/activeSource.test.ts packages/sdk-runtime/src/index.test.ts` (or the repo's configured vitest — check packages/sdk-runtime/package.json scripts; if node_modules missing at repo root, report could-not-run and note the tsc/type intent instead — do NOT install deps).
- OpenAPI: if a lint/validate script exists (check package.json / Makefile), run it; else assume-valid and state you eyeballed indentation against /erp/imports.

If the composition test that greps root.go wiring (`TestRootWiresERPImportWithoutRepositoryTenant`) or route-mount tests break from your one registration line, REPRODUCE and report — your line should only ADD a route, breaking nothing.

What you did NOT verify: live HTTP round-trip against Postgres (no DB) — note it; transport correctness proven via fake-Store unit tests only.
