# M-02 (mirror-port-active-source) — EVIDENCE

chip: local_014342cc · branch `claude/m02-mirror-port-active-source` · REBASED onto main @97e27859 (M-01 sync-state @18df2f44 + baseline-fix aggregate_sync_read arity). Original base 078caf12; rebase resolved 1 root.go import-hunk conflict ADDITIVE-KEEP-BOTH (M-01 synccomposition import + ticker line :583 kept alongside M-02 source-wiring :449-452) and bumped migration-count fixture 63→64 (0075 set 63, 0076 makes 64). Post-rebase build+affected-test GREEN.
Codex quota-wall (reset 2026-07-25) → workers = Claude sonnet fallback (core §1), dispatched SYNC. Gate crews = 2 independent Claude cold passes (Sol-waiver, disclosed below).

> **HUB RECONCILIATION (D-120, appended at accept).** This EVIDENCE.md is the 4f1bff70-generation copy (base 97e27859), so the body text above/below that says "F2 RE-CONFIRMED red" / "rebased @97e27859" is one rebase-generation stale. Corrections of record:
> - **F2 (pricing NewBatchOrchestrator arity) is RESOLVED** on main @8de7f49e (hub baseline-hygiene commit; verified 70 unit tests + full `go test ./...` green). It is NOT red on the merged tree.
> - **Merged**: hub merged 4f1bff70 --no-ff → main **@49ab3bdd** (over 8de7f49e, which already carried the F2 fix + M-01 done/F1 docs). Merged tree is a **superset** of the chip's later 74a30383 (8de7f49e-based) rebase — M-02 code is byte-equivalent; 74a30383 was not merged (redundant).
> - **Acceptance**: hub read-only spot-check 7/7 PASS; base full ladder + `harness:integration` green pre-merge; post-merge P7 dev-stack live-drive GREEN @b70fe1b8 (vitest 75/75 after the 400-window test fix; lane count=64; 0076 `\d` honest-NULL money cols + CHECK enums + tenant_default seed; live GET→200 xlsx/upload_snapshot, PUT invalid→400 invalid_active_source, PUT sankhya→200 live_read_through).
> - Non-blocking: `set_by` empty on PUT when no actor header is sent (awareness, not a contract defect).

## Deliverables (write-set)
| File | Kind | Slice |
|---|---|---|
| `apps/server_core/internal/modules/sourcekind/sourcekind.go` (+`_test.go`) | NEW dedicated SourceKind type | S1 |
| `apps/server_core/internal/modules/internal_read/ports/adapter.go` | NEW ProductSourceAdapter port (Reader + Sync + Kind) | S1 |
| `apps/server_core/internal/modules/erp_import/domain/import.go` | MOD NormalizedRow +PrecoVenda +Local | S1 |
| `apps/server_core/migrations/0076_products_mirror_active_source.sql` | NEW 3 tables + seed | S2 |
| `apps/server_core/internal/modules/tenant_config/active_source.go, context.go, repository.go` (+ tests) | NEW active-source domain/ctx/repo | S3 |
| `apps/server_core/internal/modules/internal_read/adapters/routing/reader.go` (+`_test.go`) | NEW per-tenant routing reader | S4 |
| `apps/server_core/internal/composition/root.go` | MOD source-wiring refactor (erpSource/MC_ERP_SOURCE removed, both readers + routing) | S4 |
| `apps/server_core/internal/composition/root_test.go` | MOD obsolete env-branch tests removed (co-change of deleted fn) | S4 |
| `apps/server_core/internal/platform/migrate/runner_test.go` | MOD migration count →64 (migration-owner co-change; post-rebase 63→64 atop M-01's 0075) | S2 |
| `apps/server_core/internal/modules/tenant_config/transport/http_handler.go` (+`_test.go`) | NEW GET/PUT /config/active-source | S5 |
| `contracts/api/marketplace-central.openapi.yaml` | MOD path + 3 schemas | S5 |
| `packages/sdk-runtime/src/activeSource.ts` (+`.test.ts`), `index.ts` | NEW/MOD SDK types + re-export | S5 |

## Criterion verdicts (M02-C1..C16)
| Crit | Verdict | Evidence (file:line) | Type |
|---|---|---|---|
| C1 products_mirror shape | PASS | `migrations/0076...sql:18-39` PK(tenant_id,codigo_produto); source CHECK; absent NOT NULL DEFAULT false; stale_since nullable | ran(inspect); `\d`→P7 |
| C2 stock_locations shape | PASS | `0076...sql:44-51` PK(tenant_id,codigo_produto,local_codigo) | ran(inspect); `\d`→P7 |
| C3 honest-NULL 3 cols | PASS | `0076...sql:29-31` custo/preco_venda/estoque_total bare NUMERIC, no DEFAULT/NOT NULL | ran(grep) |
| C4 active_source shape | PASS | `0076...sql:57-67` tenant_id PK; both CHECKs = ratified enums | ran(inspect); `\d`→P7 |
| C5 NormalizedRow | PASS | `erp_import/domain/import.go` +PrecoVenda +Local, prior fields intact | ran(build+diff) |
| C6 port compiles, Reader sig unchanged | PASS | `internal_read/ports/reader.go:48-55` 6 methods unchanged; `routing/reader.go:113` compile-assert | ran(build) |
| C7 SourceKind dedicated 2-val | PASS | `sourcekind/sourcekind.go:18-26` type + 2 consts; doc "NOT an extension of ImportSource" | ran(grep+test) |
| C8 ImportSource untouched | PASS | `erp_import/domain/import.go` ImportSource still {xlsx, catalogo_cliente} | ran(grep) |
| C9 seed tenant_default | PASS | `0076...sql:73-75` xlsx/upload_snapshot ON CONFLICT DO NOTHING, set_by='migration_0076' | ran(inspect) |
| C10 MC_ERP_SOURCE removed | PASS | `grep -rn MC_ERP_SOURCE apps/server_core` = 0 code (2 migration-comment mentions only); erpSource() deleted `root.go` diff | ran(grep) |
| C11 both adapters + dispatch | PASS | `routing/reader.go:45-63` resolve(); `root.go:439-451` builds uploadReader+liveReader+routingReader | ran(build+unit) |
| C12 fail-closed | PASS | `routing/reader.go:56-61` sankhya-nil→ErrActiveSourceUnavailable, default→ErrUnknownActiveSource, no upload fallback; `transport/http_handler.go:52-54` GET no-row→400; proven `routing/reader_test.go` cases 4/5 + `http_handler_test.go` | ran(unit) |
| C13 FromContext exported | PASS | `tenant_config/context.go:17-20` | ran |
| C14 additive migration | PASS | 0 ALTER/DROP/DELETE/TRUNCATE/UPDATE; only CREATE IF NOT EXISTS + INSERT ON CONFLICT | ran(grep) |
| C15 OpenAPI+SDK same changeset | PASS | `openapi.yaml:3190-3230`+`7983-8024`; `sdk-runtime/src/activeSource.ts`; PUT schema+type have NO source_kind (server-derived); error enum parity | ran(tsc+inspect) |
| C16 no stub adapter (AC-04) | PASS | `grep ProductSourceAdapter` = 2 hits, both the interface decl; no concrete impl in changeset | ran(grep) |

## P5 ladder (ran)
- `go build ./...` → GREEN (whole module).
- `go test` affected: sourcekind, tenant_config, tenant_config/transport, internal_read/adapters/routing, composition, platform/migrate → all GREEN.
- `tsc --noEmit` (sdk-runtime) → 0 errors (proves `@ts-expect-error` on source_kind is a valid rejection = server-derive contract holds).
- migration static guards: C3 no-DEFAULT/NOT-NULL ✓; C14 additive-only (0 forbidden) ✓.
- import-cycle check: tenant_config→sourcekind+erp internalread; internalread ∌ tenant_config → no cycle ✓.

## Baseline red (NOT M-02 — cited, not fixed)
- `tests/unit/pricing_{handler,service}_test.go` compile-fail: `application.NewBatchOrchestrator` called 5-arg vs 4-arg signature (`pricing_handler_test.go:122` — want (ProductProvider, PolicyProvider, FreightQuoter, string)). Pre-exists on main; RE-CONFIRMED still red after rebase onto @97e27859 — the baseline-fix merge fixed `aggregate_sync_read` NewService arity, NOT this pricing orchestrator. Zero pricing files in this write-set. **FINDING → hub** (independent baseline breakage; needs its own fix).

## Could-not-run (honest, named — not gaps)
- Migration apply proofs (C1/C2/C4/C9 `\d` output) + tenant_config Repository Get/Set fail-closed integration → need Postgres/dev-stack → **P7 hub REQUEST**. Repository integration test present, skips cleanly (reused `testsupport/postgres`).
- sdk-runtime vitest runner → env-blocked (jest-dom setupFiles resolves through worktree node_modules junction; documented false-alarm [fe-chip-vitest-node-modules-junction]). Substituted by `tsc --noEmit` PASS + manual OpenAPI-parity inspection.
- MC-01 full (real rows through both concrete adapters) = joint M-03/M-04 scope.

## P6 DUAL-GATE
Cold reviewer (harness:gate-reviewer, read-only, agentId afaaf8333e47e3aab): **PASS-WITH-NITS**, zero blocking defects; C1–C16 + non-negotiables all PASS with file:line; nits informational (EXEMPLO-IO / LIVE markers N/A for a config/routing plumbing milestone with no live-provider calls).

Adversarial refuter (independent Claude, agentId a5c92623d39bf08ae): **REFUTED round-1** — root.go now opens oracleDB independent of active_source; four oracle-batch consumers (listingCostReader, inventoryStockReader, profitability FactReader.batch, assisted Sankhya linkage) read oracle directly, bypassing routing.Reader → cross-source mix when a tenant is on xlsx but oracle is configured.

Orchestrator adjudication (with reconciliation round-2, per [chip-import-fix] no-self-verified-round-2 lesson):
- Finding is FACTUALLY CORRECT on the diff (verified `git diff main -- root.go`: old code was xlsx⊕oracle mutual-exclusion; new code opens oracle whenever its env is present).
- Finding is INERT in the current deployment: oracle env unconfigured → `LoadConfigFromEnv` err → `oracleDB == nil` → all four gates false → byte-identical runtime to old xlsx-mode. ZERO live regression today.
- Finding is UNFIXABLE in M-02: the four consumers use oracle-ONLY batch interfaces with no xlsx/mirror-backed equivalent to route to; those are M-03/M-04 deliverables. M-02's charter is the internal_read.Reader 6-method interface only, which IS routed fail-closed.
- **Verdict: DEFERRED CROSS-MILESTONE SEAM owned by M-04, not an M-02 blocker.**

Refuter round-2 (reconciliation, agentId a5c92623d39bf08ae): **withdrew the finding as M-02-blocking → verdict (B) deferred M-04 seam**. Independently confirmed `internalreadoracle.NewBatchReader(nil, ...)` (root.go:602) stores db in a field and does NO I/O at construction (nil queryer deref only at request time) → with oracleDB==nil today, listingCostReader errors honestly at request, inventoryStockReader stays at UnavailableStockBatchReader default, profitability/linkage gates false → byte-identical to old xlsx-mode, honest-unavailable not silent-wrong. Reconciliation was genuine (orchestrator supplied the `git diff main` facts; refuter re-verified NewBatchReader independently — not a self-verified round-2, per [chip-import-fix] lesson).

**P6-DUAL-GATE: AGREEMENT** (cold PASS-WITH-NITS + refuter NO-REFUTATION after round-2). Sol-waiver: GPT-5.6 Sol side quota-walled until 2026-07-25 → both gate passes run by independent Claude cold crews (core §1 sanctioned fallback), disclosed per protocol.

### FINDING → hub (M-04 pre-activation constraint, verbatim from refuter)
Before M-04 wires the live (sankhya/oracle) batch readers into `root.go`, it MUST route `listingCostReader` (root.go:602), `inventoryStockReader` (root.go:410/435), `profitabilityCfg.Internal`'s batch cost/tax facts (root.go:560), and the assisted-linkage gate (root.go:521) by the tenant's `active_source` — or select mirror-backed equivalents — because these consumers read `oracleDB` directly and bypass `routing.Reader`; once oracle is reachable while a tenant's `active_source` is xlsx/catalogo_cliente, they will silently serve live Oracle data for that tenant's listing cost, inventory risk, and profitability facts while the product catalog correctly serves upload data — a cross-source mix with no error. INERT in the M-02 foundation deployment (oracle unconfigured).

## Dispatch ledger
See `LEDGER.md` (D1 S1+S3, D2 S4, D3 S5, P5 orchestrator ladder, P6 dual-gate).
