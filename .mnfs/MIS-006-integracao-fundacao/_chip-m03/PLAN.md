# M-03 — XlsxAdapter · P2 Plan (orchestrator inventory + design + slice cards)

base_sha: 78f02ac9 (chip prompt BASE; contract yaml says 138aac3d = stale pre-M01/M02 anchor) · branch: claude/m03-xlsx-adapter
Codex quota returned 2026-07-22 → workers = **REAL codex** (GPT-5.6 Luna high std / Sol low complex); dual gate = cold Opus + **GPT-5.6 Sol medium REAL** (onda-1 Sol-waiver DEAD).

## Inventory of record (file:line, read pre-write)

- **reader.go (erp_import/adapters/internalread)** satisfies TWO ports — asserted `reader.go:401` `var _ readports.Reader` + `reader.go:402` `var _ readports.CatalogPageReader`. Both MUST survive F-02 rewrite (D-120 /catalogo-503 regression = exactly a dropped CatalogPageReader). Compile-time asserts stay.
- **Two-layer source model** (do NOT conflate):
  - erp_import `ImportSource` = {xlsx, catalogo_cliente} (`domain/import.go:87`) — WHICH upload dataset.
  - `sourcekind.SourceKind` = {upload_snapshot, live_read_through} (`internal/modules/sourcekind`) — WHICH adapter shape.
- **active_source lookup already owned by M-02 routing** — `routing/reader.go:45-63` `resolve()` calls `lookup.Get(ctx,tenantID)` (tenant_config.Repository, DB-backed `repository.go:31`), fail-closed `ErrUnknownActiveSource` on no-row, then pins `erpinternalread.WithActiveSource(ctx, ImportSource(cfg.Source))` (`routing/reader.go:53`) for the upload sub-toggle. The erp reader's `activeSourceFromContext` (`reader.go:43-48`) just READS that pinned value, silent-defaulting to `SourceXLSX` when unpinned.
- **IMPORT CYCLE FACT**: `tenant_config` imports `erp_import/adapters/internalread` (`active_source.go:11`, for the `ErrUnknownActiveSource` alias). ⇒ `internalread` CANNOT import `tenant_config`. F-03/C13 literal "reader consulta o lookup de active_source (M-02)" is **not implementable inside internalread** without a cycle. See BLOCKING-2.
- **Mirror I/O can flow through existing injected `erpRepo`** (`ImportRepository`, `erp_import/ports/repository.go:9`, I OWN it). Add mirror upsert + mirror read methods to interface + `erppostgres` impl. `NewReader(erpRepo,...)` (`root.go:442`) and `NewImportService(parser,erpRepo,...)` (`root.go:284`) signatures UNCHANGED ⇒ no root.go delta for mirror read/write.
- **F-01 enqueue seam pre-built by M-01**: `SyncStateRepository.AppendPendingCodigo(ctx, installationID, entity, codigo)` (`sync/adapters/postgres/sync_state_repo.go:112`) = atomic cursor-append (M01-C11). Hook calls it per touched codigo, `entity=market`.
- **F-01 generation trigger target**: `GenerationService.GenerateLinkCandidates(ctx, {InstallationID,Limit})` (`product_links/application/generation_service.go:60`). Currently wired at `root.go:469`, AFTER `NewImportService` (`root.go:284`).
- **Mirror schema (0076)**: PK `(tenant_id, codigo_produto)`, source CHECK ∈ {sankhya,xlsx,catalogo_cliente}, custo/preco_venda/estoque_total NUMERIC NULLable no-default (honest-NULL), `absent_in_last_snapshot BOOL default false`, `stale_since TIMESTAMPTZ`, `protocol_id UUID`, `updated_at`. Plus `products_mirror_stock_locations` PK `(tenant_id,codigo_produto,local_codigo)`.
- **NormalizedRow E2 fields** (`domain/import.go:23`): Codprod, Descrprod, Custo, PrecoVenda, StockPhysical, StockReserved, EAN, Refforn, Marca, NCM, Grupo, DescrGrupo, Local. estoque_total = physical (stock breakdown → stock_locations via Local).

## Blocking items → hub (sent as REQUEST, awaiting ruling before dispatch)

### BLOCKING-1 — root.go additive-wiring lease (F-01 hook deps)
M03-C8 requires the generation-candidate call **inside import_service.go**; F-01 also enqueues sync_state. Both are NEW deps on `NewImportService` not currently injected (`root.go:284`). root.go is M-02/hub-owned (profile §6; chip prompt "root.go ZERO toque"). Proposed minimal delta (precedent: M-01 & M-02 additive-lock leases):
- Inject an optional `LinkCandidateGenerator` iface (defined IN erp_import, structural — no product_links import) + a `SyncEnqueuer` iface (structural, no sync import) into ImportService via functional options.
- root.go: reorder so `productLinkGenerationSvc` builds before `erpImportSvc`, OR use setters; add 2 `With…` option calls + adapters. ~4-8 additive lines, zero behavior change to existing wiring.
- Ask: hub wires it, OR grants me a scoped additive-lease on `root.go:284` region.

### BLOCKING-2 — F-03/C13 reconciliation (import cycle vs landed M-02)
C13 "prova: diff mostra a função consultando o repo/lookup de active_source (M-02)". Landed M-02 already does this lookup in `routing/reader.go:46`; internalread can't re-do it (cycle). Proposed reading (needs hub ratify, gate authority = contract):
- F-03 satisfied by: (a) value observed by `activeSourceFromContext` is fed from active_source via routing pinning (ALREADY true post-M-02), + (b) **remove the silent `SourceXLSX` default** at `reader.go:47` → fail-closed when unpinned (no silent xlsx fallback — the actual F-03 EARS intent + C14).
- C13 "prova mínima" reinterpreted: diff shows the silent-default removed + no `ParseActiveSource(defaultParam)`-style hardcode remains; the lookup-of-record is M-02's routing (cite `routing/reader.go:46`), not a duplicated internalread lookup.

### OPEN (P2 resolves, not blocking) — installationID sourcing
Generation + sync enqueue are installation-scoped; xlsx import is tenant-scoped. Hook must resolve which installation(s). Candidate: iterate tenant installations (`InstallationService.List`, `installation_service.go:104`) and trigger per-installation. P2 planner (Sol medium) to decide + justify.

## Draft slice cards (provisional — finalized post-hub-ruling by Sol-medium P2)

- **S1 — mirror repo (write+read) on erpRepo**: extend `ImportRepository` iface + `erppostgres` impl with `UpsertMirrorMerge(ctx,tenant,rows,protocolID)` (keep-absent: upsert touched + flag absent set-diff in one tx) + mirror read methods for F-02. Integration (Postgres) tests C2/C3/C5/C6/C7/C17. write_set: `erp_import/ports/repository.go`, `erp_import/adapters/postgres/*`.
- **S2 — F-01 hook in import_service.go**: post-`PersistSnapshotAtomically` (`import_service.go:120`), same local tx: mirror upsert-merge → trigger generator (injected iface) → enqueue sync (injected iface, AppendPendingCodigo per codigo). Propagate upsert error (C10). grep-guard zero `Collect(`/market_aggregates/competitor_offers (C9/AC-06). write_set: `import_service.go`.
- **S3 — F-02 mirror-backed read in reader.go**: FindProductsForLinking/GetSellableStock/GetCostAsOf/catalogPage read products_mirror via erpRepo mirror-read; `catalogPage` `ORDER BY codigo_produto ASC` (C11b); KEEP both compile-asserts (401/402). write_set: `internalread/reader.go`.
- **S4 — F-03 fail-closed activeSourceFromContext**: remove silent xlsx default (per BLOCKING-2 ruling). write_set: `internalread/reader.go` (same commit as S3 possible).
- **S5 — F-04 XlsxAdapter**: concrete type in `erp_import/adapters/xlsx/` (NEW file, NOT parser.go) satisfying `ProductSourceAdapter`; `Sync()` = same path as F-01 hook (not a 2nd impl); `Kind()=upload_snapshot`; idempotency test C17. write_set: `erp_import/adapters/xlsx/adapter.go`(+test).

KEEP absolute: `parser.go` zero edit (C1/AC-04). No new migration. No root.go beyond BLOCKING-1 lease.
