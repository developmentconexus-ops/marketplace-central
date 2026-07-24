# M-03 CHIP EVIDENCE

## F-03 precondition — internalread reader caller audit (BLOCKING-2 hub constraint)

Ran (grep + read, base @78f02ac9). Hub required: prove no caller bypassing routing depends on the silent xlsx default (reader.go:47) before removing it.

**Callers that pin `erpinternalread.WithActiveSource` OUTSIDE routing:**
- `catalog/transport/http_handler.go:52-54` — `if present { WithActiveSource }`; absent `erp_source` → unpinned. Doc:45 "Absent = reader default (xlsx)".
- `market/transport/collection_handler.go:36-38` — same `if present` shape; doc:28 "Absent = default (xlsx)".
- (`cache/cache.go:288` READS `ActiveSourceFromContext` for cache-key only, serves no data.)

**Determination — both flow THROUGH routing, which re-pins → default is dead code:**
- catalog page/canonical/compat readers all = `internalReadSvc` (`root.go:465,458,462`), and `internalReadSvc = internalreadapp.NewService(routingReader)` (`root.go:451`).
- market `Collect` → `s.identity` = `marketIdentityReader = newMarketProductIdentityReaderAdapter(internalReadSvc, …)` (`root.go:588`) → also `internalReadSvc` (routing).
- `routing.Reader.resolve` (`routing/reader.go:45-63`) UNCONDITIONALLY `lookup.Get` + re-pins `tenant_config.WithActiveSource` (:50) + `erpinternalread.WithActiveSource(cfg.Source)` (:53), overriding any handler pin; no-row → `ErrUnknownActiveSource` (fail-closed). So the erp reader is NEVER reached with unpinned ctx in the wired composition; the handlers' manual pin is overridden.

**Conclusion:** removing the `return erpdomain.SourceXLSX` silent default (reader.go:47) is SAFE — no production path depends on it. Reader always observes the routing-pinned source (or the read fails closed upstream). Test files that assert the old default (`reader_test.go`, `source_contract_test.go`) are updated in the F-02/F-03 slice (reader is rewritten there anyway). Doc-comment `reader.go:25-26` to be updated to fail-closed wording (hub-mandated).

**Side FINDING (out of M-03 scope, → hub):** the `erp_source` query-param handling in catalog+market handlers is now DEAD post-M-02 — routing overrides it with the DB active_source. The client toggle moved to `PUT /config/active-source` (M-02). Not M-03's write-set; flagging for M-06/hub backlog.

---

# M-03 CLOSURE EVIDENCE PACK

```yaml
milestone: M-03-xlsx-adapter
mission: MIS-006-integracao-fundacao
branch: claude/m03-xlsx-adapter
base_sha: 78f02ac9        # M-02 merged main (DAG root); milestone diff = 78f02ac9..HEAD
head_sha: 71b34c5         # S9 corrective (D4 fixture + D5 per-source-kind data-time); code tip. Docs tip d85038b = ledger only.
validation_contract: M-03-xlsx-adapter/validation-contract.md (M03-C1..C17, AC-01..08, U1..U4)
```

## Deliverables (7 slices, serialized single-worktree DAG S1→S2→S4→S5→S3→S6→S7)

| Slice | Commit | Deliverable |
|-------|--------|-------------|
| S1 | 9fee2264 | Atomic mirror materialization: shared `mergeSnapshotTx` (staging temp table `ON COMMIT DROP` + CopyFrom → upsert touched + source-scoped set-diff `absent_in_last_snapshot=true`, NO parent delete, child stock_locations delete-touched-then-upsert, honest-NULL casts) + `SyncLatestCompletedSnapshot` + `ImportRepository` iface extension. Merge call runs in the snapshot-persist tx (`import_repository.go`). |
| S2 | c4d23b74 | Post-import hook in `import_service.go`: structural ifaces `SyncEnqueuer`+`LinkCandidateGenerator` (no product_links/sync/integrations import), `SetPostImportHooks` late setter. COMPLETED-only after `PersistSnapshotAtomically`: collect accepted Codprod → EnqueueMarketProducts (fail-loud) → GenerateLinkCandidates per installationID (fail-fast). NEVER Collect. |
| S4 | 89023c0d | Mirror query repo (F-02 reads): `MirrorRows`/`MirrorProductByCode`/`MirrorCatalogPage`(numeric keyset)/`MirrorEANCollisionCounts`(global), all tenant+source+`absent=false` scoped; `domain/mirror.go` MirrorProduct pre-computed nullable numerics (honest-NULL). |
| S5 | 9a650cf9 | Mirror-backed reader + F-03 fail-closed: 5 reads off `LatestCompletedSnapshot` rescan onto mirror; `activeSourceFromContext→(src,err)` unpinned→`ErrUnknownActiveSource` wrapped `ReadErrorSourceUnavailable`; `snapshot()` removed; both compile-asserts survive; consumer zero-diff; catalogPage keyset+limit+1+NextCursor. |
| S3 | 5fe2abc1 | Real hook adapters making S2 hook load-bearing: `adapters/sync/enqueuer.go` (MarketEnqueuer → InstallationService.List × codes → AppendPendingCodigo(EntityMarket) fail-loud → distinct installationIDs) + `adapters/productlinks/generator.go` (wraps GenerationService, Limit:0→20). root.go additive-lease: 3 imports + ONE SetPostImportHooks call, zero reorder. |
| S6 | c727a164 | Concrete `XlsxAdapter` (F-04): embeds readports.Reader (read side verbatim) + Sync()→SyncLatestCompletedSnapshot(SourceXLSX) same merge path + Kind()=upload_snapshot; `var _ readports.ProductSourceAdapter`. parser.go ZERO edit. |
| S7 | (chip static gate) | Boundary/static sweep — chip-proven evidence below. |

## Chip-proven evidence (ladder L0/L1 — `ran`, outputs saved this session)

- **M03-C1 / AC-04 (parser.go zero-edit)** — `ran` L0: `git diff 78f02ac9..HEAD -- .../adapters/xlsx/parser.go` = EMPTY. S6 added `adapter.go`+`adapter_test.go` only.
- **M03-C9 / AC-06 boundary** — `ran` L0 grep: forbidden-token sweep on ADDED lines `78f02ac9..HEAD` → 0 production hits of `Collect(`/`market_aggregates`/`competitor_offers`. Production NEVER `DELETE FROM products_mirror` (parent); only sanctioned child-table `DELETE FROM products_mirror_stock_locations` @mirror_repository.go:132. All bare parent deletes are `_test.go` teardown.
- **M03-C11 (mirror read, no rescan)** — `ran` L1: reader.go 5 reads consume the Mirror* methods; `snapshot()` + all `LatestCompletedSnapshot` calls removed from reader.go.
- **M03-C11b (deterministic pagination)** — `ran` L1: `MirrorCatalogPage` `ORDER BY codigo_produto::bigint ASC` + `^[0-9]{1,18}$` guard + `LIMIT CASE WHEN $5>0 THEN $5+1 END`; reader hasMore(limit+1)+truncate+NextCursor=last id. Unit fixtures 1,2,10 → page1=[1,2] page2=[10], no dup/skip.
- **M03-C12 (consumer zero-diff)** — `ran` L1: `go build ./...` GREEN; `git diff --stat 78f02ac9..HEAD -- ':!.../erp_import'` = only root.go 7-line additive lease. Both compile-asserts survive reader.go:393-394.
- **M03-C13/C14 (fail-closed)** — `ran` L1: `activeSourceFromContext(ctx)(ImportSource,error)`; unpinned→`ErrUnknownActiveSource` wrapped `ReadErrorSourceUnavailable`, NO repo read on miss. Test `TestUnpinnedReadsFailClosed`.
- **M03-C15/C16 (adapter + Kind)** — `ran` L1: `var _ readports.ProductSourceAdapter` compiles; `Kind()==sourcekind.UploadSnapshot`.
- **M03-C17 (idempotent Sync, adapter level)** — `ran` L1: Sync delegates shared `SyncLatestCompletedSnapshot(SourceXLSX)` PK-upsert; stable delegation+mapping test. End-to-end = S1 PG integration + hub P7.
- **Build/suites** — `ran`: build OK; `erp_import/...` all ok; `composition/...` ok (4.086s); `vet -tags integration erp_import/...` compiles.

## Deferred to hub P7 (`could-not-run` — named block: no Postgres / no browser in chip sandbox)

- **M03-C2/C3/C4/C5/C6/C7 (L2 integration)** — DB state post-import: keep-absent flag+stale_since, product_links intact, flag clear/reappear, honest-NULL rows, parser-rejected orphan absent. Integration tests WRITTEN (`mirror_repository_integration_test.go`, `reader_integration_test.go`, `//go:build integration`) — compile (`go vet -tags integration` OK) but need `MPC_TEST_DATABASE_URL`. Hub P7 runs against dev-stack Postgres.
- **M03-C8 (hook fires generation, DB-observable)** — code-side function-call proven; DB candidate-after-import = hub P7.
- **M03-C10 (upsert-failure propagation, DB)** — tx-boundary design proven (merge in snapshot-persist tx; preco_venda mirror-only → invalid text passes protocol insert, fails `::numeric` cast → whole tx rolls back); DB-forced-error run = hub P7.
- **M03-U1/U2/U3/U4 (user-drive, D-120)** — browser drive /catalogo (2-source flip), /integracoes (upload+409), vínculos post-import, 4-screen no-regression. Hub P7 live-drive; fixtures/flows ready.

## Boundary/scope notes carried into the gate
- MC-01 full ("fed by BOTH adapters") closes only with M-04 merged; M-03 proves the xlsx side (C2). `WHERE source IN ('xlsx','sankhya')` = `could-not-run` named until M-04.
- M03-C8 proves the hook FIRES candidate generation; the EAN-exact-unique auto-approve RULE is M-05, not M-03 (contract lines 57-63).

## P6 DUAL-GATE

**P6-DUAL-GATE: AGREEMENT** (corrective HEAD `fff4128`, both reviewers independent, no waiver)

Round-1 (HEAD `c727a164`) = **DISAGREEMENT**: cold Claude PASS vs GPT-5.6 Sol-medium FAIL (3 defects). Chip adjudicated on merits (no rubber-stamp) → Sol correct on all 3 → corrective slice **S8** dispatched + committed @`fff4128` (DISTINCT ON int64-identity keyset; C17 real 2×-call integration invariant + adapter test de-clawed to delegation-only; C4 stale product's OWN link).

Round-2 (HEAD `fff4128`, REAL 2nd refuter pass per CHIP-IMPORT-FIX lesson — NOT self-verified) = **BOTH PASS**:

| Reviewer | Path | Verdict |
|----------|------|---------|
| Cold Claude gate-reviewer | Agent, physically read-only | **GATE: PASS** — D1/D2/D3 all FIXED & probative (tie case would fail against pre-fix SQL; C17 2×-call exercises real merge; C4 same-identity throughout); no S8 regression; AC-01..08 clean (AC-07 could-not-run/no-git, no push evidence). |
| GPT-5.6 Sol-medium (round-1 FAIL reviewer) | codex OS-process read-only | **GATE: PASS** — same 3 fixes verified; full rubric re-run; S8 changed exactly 3 declared files, parser.go/root.go/reader.go/catalog_page.go empty vs `c727a164`; defects: none. |

Both agree the DB/browser-dependent criteria (C2–C10, C11b runtime, C17 runtime, U1–U4) are honest `could-not-run` → hub P7. Code-side + static + boundary criteria PASS with file:line evidence (C1, C11, C12–C16 PASS; AC-01..08 PASS).

### Non-blocking process note (cold Claude, → hub/QA mission-gate, NOT an M-03 code-gate blocker)
Mission `interface-contracts-mis006.md:46-49` defines an EXEMPLO-IO golden fixture (Sankhya 90008 / xlsx 74606). No M-03 test asserts that named golden case (tests use synthetic codprods `SYNC`/`42`/`X`/`1`/`01`/`2`). Not an M-03-VC criterion; hub/QA should confirm coverage at a later mission-level gate.

Gate artifacts: `scratchpad/agent__gate2-sol-m03.last.md` (GPT); cold Claude verdict in chip transcript. Round-1: `scratchpad/agent__gate-sol-m03.last.md`.

### P6 DUAL-GATE DELTA (S9 corrective @`71b34c5`, delta `fff4128..71b34c5`) = **AGREEMENT / GATE: PASS**

S9 fixes the hub's real-Postgres integration-lane RED (@fff41288): **D4** protocol-fixture CHECK violation (all 4 integration tests died at setup, never ran) + **D5** real semantic regression (reader gated cost-as-of / reported ObservedAt on `row.UpdatedAt`, which the merge stamps `now()` every sync → wrong as-of failures + fabricated sync-time ObservedAt). Fix = data-time carried per-source-kind (upload→ImportedAt honest-fail on unknown; live→UpdatedAt honest live observation), per hub ruling. Real 2nd refuter pass (both reviewers independent, NO self-substitution; Sol was the round-1 FAIL reviewer):

| Reviewer | Path | Verdict |
|----------|------|---------|
| GPT-5.6 Sol-medium (round-1 FAIL reviewer) | codex OS-process, read-only | **GATE: PASS** — D4 FIXED (helper explicit protocol; `#611-E`..`#619-E` regex-valid + unique; distinct FileSHA256; no assertion weakened), D5 FIXED (ImportedAt pointer + 3-read LEFT JOIN `p.imported_at` via sql.NullTime; per-source-kind at reader.go:25-51; all 6 sites via `dataObservationTime`, no bare `row.UpdatedAt`; honest fail-closed; live-branch honestly unit-covered, F1-labeled), MirrorCatalogPage KEEP-EXACT verbatim, KEEP-ABSOLUTE clean, NEW defects: none. |
| Cold Claude gate-reviewer | Agent, physically read-only | **GATE: PASS** — same fixes verified w/ file:line; column-order integrity of added JOIN hand-checked (12+1 cols ↔ 13 scan dests, no shift); nil-deref guarded reader.go:47-50; only `row.UpdatedAt` use is INSIDE `dataObservationTime` (reader.go:45), all 6 consumers route through it; regression tests updated honestly (ObservedAt.Equal(importedAt), not deleted); compile-asserts survive reader.go:445-446; NEW defects: none. |

Both agree DB integration (C2–C10, C17 runtime) + U1–U4 browser = honest `could-not-run` → hub P7 (S9 does NOT newly claim they ran; inherited from round-2 AGREEMENT). D4's own fix means those integration tests can now ACTUALLY execute against real PG — hub re-runs the lane before merge. EXEMPLO-IO golden still unasserted at M-03 level (carried non-blocking → mission/hub QA, not an M-03-VC criterion, not S9-scoped). Delta-gate artifacts: `scratchpad/agent__gate3-sol-m03.last.md` (GPT); cold Claude verdict in chip transcript.

**M-03 chip-side DELTA gate GREEN. Emitting new tip `71b34c5` to hub for real-PG integration re-run + P7.**
