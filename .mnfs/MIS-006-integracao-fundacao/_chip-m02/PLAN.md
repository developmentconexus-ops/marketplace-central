# M-02 — P2 Plan (slice cards + write-set + verification map + seam-closure)

base_sha: 078caf12 (main superset of contract base 138aac3d) · branch: claude/m02-mirror-port-active-source
Codex quota-wall (reset 2026-07-25) → all workers = **Claude sonnet fallback** subagents (core §1), dispatched SYNC (§3).

## Package/design decisions (orchestrator, contract-fit)
- `SourceKind` in dedicated dep-free pkg `internal/modules/sourcekind` → both `ports` + `tenant_config` import, no cycle (ADR-02). `{upload_snapshot, live_read_through}` (M02-C8).
- `ProductSourceAdapter` in `internal_read/ports/adapter.go`: embeds existing `Reader` verbatim + `Sync(ctx)(SyncResult,error)` + `Kind() sourcekind.SourceKind`. No signature change to Reader (M02-C6/C7). Named consumers = M-03/M-04 briefs (not speculative).
- `tenant_config` pkg owns: `ActiveSource` enum {sankhya,xlsx,catalogo_cliente}, `Config`, ctx `WithActiveSource/FromContext` (C13), `ActiveSourceLookup` port (consumers: routing reader + postgres repo — 2 consumers, justified), postgres `Repository` (Get→ErrUnknownActiveSource on no-row, Set upsert). **Reuse** sentinel `internalread.ErrUnknownActiveSource` via alias (C12; contract says reuse not reimplement; no cycle: internalread ∌ tenant_config).
- Routing: new `internal_read/adapters/routing/reader.go` — `routingReader` implements `internalreadports.Reader`, holds upload-reader + live-reader(nullable) + lookup. Per call: resolve active_source → inject into ctx (tenant_config.WithActiveSource + erp internalread.WithActiveSource for xlsx/catalogo sub-toggle) → delegate. sankhya w/ nil oracle = typed unavailable error; unknown tenant = ErrUnknownActiveSource. Fail-closed, no silent default (C11/C12).
- Endpoint tenant-IMPLICIT (repo convention `/erp/imports`): `GET /config/active-source` + `PUT /config/active-source`, tenant from cfg.DefaultTenantID. PUT invalid enum → 400 (C16).
- Migration `0076_products_mirror_active_source.sql` (bloco B, hub-ratified number). Seeds tenant_default=xlsx/upload_snapshot (explicit config, not code fallback — preserves boot after MC_ERP_SOURCE removal).

## Slice cards

### S1 — leaf domain types (sourcekind + port + NormalizedRow)  [DRAFTED by orch, worker owns+tests]
- validation_kind: unit + build
- write_set: `internal/modules/sourcekind/sourcekind.go`(+`_test.go`), `internal/modules/internal_read/ports/adapter.go`, `internal/modules/erp_import/domain/import.go`
- failing-test-first: `sourcekind_test.go` table test for `Valid()` (upload_snapshot✓ live_read_through✓ ""✗ "bogus"✗)
- commands: `GOCACHE=.gocache go build ./internal/modules/sourcekind/... ./internal/modules/internal_read/ports/... ./internal/modules/erp_import/domain/...` ; `GOCACHE=.gocache go test ./internal/modules/sourcekind/...`
- done: build green, Valid test green. Crit: C5, C6(compile), C7, C8
- open_questions: none

### S2 — migration 0076  [DRAFTED by orch, worker verifies static guards]
- validation_kind: static-grep (L0) + apply-proof deferred to P7/Postgres
- write_set: `apps/server_core/migrations/0076_products_mirror_active_source.sql`
- commands (static): grep guard C3 `custo|preco_venda|estoque_total` have NO `DEFAULT 0`/`NOT NULL`; C14 additive-only (only CREATE/INSERT, zero ALTER..TYPE/DROP)
- done: guards pass. Crit shape C1/C2/C4/C9 = could-not-run til Postgres apply (P7 hub REQUEST)
- open_questions: none

### S3 — tenant_config pkg (domain + ctx + lookup + postgres repo)  [worker]
- validation_kind: unit (pure) + integration (Postgres, deferred)
- write_set: `internal/modules/tenant_config/{active_source.go, context.go, repository.go, active_source_test.go, context_test.go}`
- failing-test-first: `active_source_test.go` — `ActiveSource.Valid()` table + `DefaultKind()` (sankhya→live_read_through, xlsx/catalogo→upload_snapshot); `context_test.go` — WithActiveSource/FromContext round-trip + absent→false
- commands: `GOCACHE=.gocache go build ./internal/modules/tenant_config/...`; `GOCACHE=.gocache go test ./internal/modules/tenant_config/...`
- done: build+unit green. Repo Get/Set integration (fail-closed on no-row → ErrUnknownActiveSource) = integration test present, marked could-not-run til Postgres. Crit: C9(shape via DDL), C12(logic proven in S4 via fake lookup), C13
- open_questions: none

### S4 — root.go routing refactor + routing reader  [worker, complex]
- validation_kind: unit + full build
- write_set: `internal/modules/internal_read/adapters/routing/{reader.go, reader_test.go}`, `internal/composition/root.go`
- failing-test-first: `reader_test.go` — fake upload Reader + fake live Reader + fake ActiveSourceLookup:
  - lookup xlsx → delegates to upload reader, ctx carries erp ImportSource=xlsx + tenant_config active_source (C11 both-built dispatch)
  - lookup sankhya + live present → delegates to live reader
  - lookup sankhya + live nil → typed unavailable error (fail-closed, not upload fallback)
  - lookup returns ErrUnknownActiveSource → propagated, no default (C12)
- root.go edits (OWNED source-wiring, disjoint from M-01 ticker block): remove `erpSource(getenv)` func (~:772) + call (~:271) + MC_ERP_SOURCE (C10); build xlsxReader ALWAYS + oracleReader if configured; construct routingReader(upload, live, lookup); `internalReadSvc = NewService(routingReader)`; delete `if source=="xlsx"` branch (~:419-453)
- commands: `GOCACHE=.gocache go build ./...`; `GOCACHE=.gocache go test ./internal/modules/internal_read/adapters/routing/...`; `grep -rn MC_ERP_SOURCE apps/server_core` == 0 (C10)
- done: full build green, routing unit green, grep 0. Crit: C10, C11, C12
- open_questions: oracle graceful-degrade preserved (nil live reader when oracle unconfigured) — YES, keep existing LoadConfigFromEnv/OpenDB warn path

### S5 — active-source transport + OpenAPI + sdk-runtime (SAME commit)  [worker]
- validation_kind: unit + build + tsc
- write_set: `internal/modules/tenant_config/transport/{http_handler.go, http_handler_test.go}`, `contracts/api/marketplace-central.openapi.yaml`, `packages/sdk-runtime/src/activeSource.ts`, `packages/sdk-runtime/src/index.ts`, `packages/sdk-runtime/src/activeSource.test.ts`
- failing-test-first: `http_handler_test.go` w/ fake repo — PUT `{"active_source":"foo"}`→400 nothing written (C16); PUT `{"active_source":"xlsx"}`→200 persists w/ derived kind; GET→200 returns config; GET no-config→ maps ErrUnknownActiveSource
- commands: `GOCACHE=.gocache go build ./...`; `GOCACHE=.gocache go test ./internal/modules/tenant_config/...`; sdk-runtime `tsc --noEmit` + vitest
- done: build+unit+tsc green. Crit: C15(spec+SDK same commit), C16
- open_questions: sdk-runtime is hand-written types (no codegen) — activeSource.ts = types mirroring spec, index re-export

## Verification map (criterion → command/proof → carrier file)
| Crit | Proof | Carrier | Type |
|---|---|---|---|
| C1 products_mirror shape | `\d products_mirror` after apply | migration + Postgres | could-not-run→P7 |
| C2 stock_locations shape | `\d products_mirror_stock_locations` | migration | could-not-run→P7 |
| C3 honest-NULL 3 cols | grep migration no DEFAULT 0/NOT NULL | migration | ran(grep) |
| C4 absent/stale defaults | DDL text | migration | ran(grep)+P7 |
| C5 NormalizedRow 10 | diff import.go | import.go | ran |
| C6 port compiles no sig change | go build ./... + diff | adapter.go | ran |
| C7 Sync/Kind declared | grep adapter.go | adapter.go | ran |
| C8 SourceKind dedicated 2-val | grep + ImportSource untouched | sourcekind.go | ran |
| C9 active_source shape | `\d active_source` | migration | could-not-run→P7 |
| C10 MC_ERP_SOURCE=0 | `grep -rn MC_ERP_SOURCE apps/server_core` | root.go | ran(grep) |
| C11 both adapters, no boot branch | routing unit + root.go diff | routing+root | ran |
| C12 fail-closed | routing unit (fake lookup no-row) | reader_test.go | ran |
| C13 active_source from ctx exported | grep FromContext usable | context.go | ran |
| C14 additive migration | `git diff` only CREATE/INSERT | migration | ran(diff) |
| C15 OpenAPI+SDK same commit | `git show --stat` | commit | ran |
| C16 enum reject 400 | handler unit | http_handler_test.go | ran |

## Seam-closure checklist
- Composition root wiring: routingReader replaces single-source branch (S4) — inside owned write-set ✓
- Route registration: tenant_config transport Register(mux) called in root.go (S5) — owned ✓
- SDK surface: activeSource.ts types + index export (S5) — owned ✓
- API spec: OpenAPI path/schema (S5) — contract-lock, disjoint from M-06 chain-read ✓
- Governance/registry: modules.json entry for tenant_config? → landed via merge, NOT pre-merge (memory governance-registry-entry-timing). Check if new module needs registry row → hub at merge.
- ErrUnknownActiveSource reuse (not new sentinel) ✓ · erp ImportSource ctx sub-toggle preserved ✓

## Named could-not-run (honest, not gaps)
- MC-01 full (real rows via both adapters) = joint M-03/M-04.
- Migration apply proofs (C1/C2/C4/C9 `\d`) + repo integration (Get/Set/fail-closed against real Postgres) = need dev stack → hub REQUEST at P7.
