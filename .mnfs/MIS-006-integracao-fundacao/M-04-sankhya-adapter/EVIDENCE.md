# M-04 SankhyaAdapter — Evidence Pack

Branch `claude/dreamy-lamarr-0bd9de`. Commits:
- `18b278dc` F-03 cache-key partition by active source incl sankhya + must-fail test
- `0af8d28c` F-01/F-02 SankhyaAdapter Sync + own mirror writer
- `b185cc3c` tests (sync mapping, mirror guards, integration keep-absent)

Base: mapping doc `docs/design/evidence/sankhya-mapping-M04.md` @`dbe089ac` (mission main, merged).

## Criteria ledger

| ID | Verdict | Evidence |
|----|---------|----------|
| M04-C1 | ran | Mapping doc @`dbe089ac` (db-consult, 6/6 items ratified, verified live METALPRD 2026-07-22). Sync entrypoint commit `0af8d28c` timestamp POSTERIOR to mapping. AC-07 clear. |
| M04-C2 | n/a | [TESTAR-SKW] not blocked (mapping returned). |
| M04-C3 | ran | `oracle/reader.go` read-side untouched — `git show 0af8d28c --stat` adds only `sync.go` (new file) + mirror writer; FindProductsForLinking/GetSellableStock/GetCostAsOf/GetTaxInputs signatures + bodies unchanged. SankhyaAdapter embeds `*Reader` verbatim. |
| M04-C4 | **could-not-run (chip) → hub live-drive REQUESTED** | Chip never boots server (AC-04 no-stub, profile). Real-Oracle `SELECT * FROM products_mirror WHERE source='sankhya'` post-Sync = hub live-drive REQUEST (task #8). NOT stubbed. |
| M04-C5 | ran (logic, unit) + could-not-run→hub (real PG) | `mirror/writer_integration_test.go` (//go:build integration) proves 2-round keep-absent, present→absent flag + stale_since + last-known preserved, reappear-clears, tenant isolation, empty-snapshot guard. Runs in hub integration lane (REQUEST task #8). Guard logic unit-proven in `writer_test.go`. |
| M04-C6 | ran | Honest-NULL: `sync.go` sets custo/preco/estoque/ean via `*T` pointers left nil when unresolved; `TestSankhyaSyncMapsSnapshotHonestNull` asserts product 300 (absent everywhere) → custo/preco/estoque NULL, and product 200 non-EAN REFERENCIA → ean NULL, real 0 stock ≠ NULL. Grep: zero hardcoded 0/default on those columns (mirror migration NUMERIC nullable no-default; sync writes pointers). |
| M04-C7 | ran | `go build ./...` green; `Kind()` returns constant `sourcekind.LiveReadThrough`; `TestSankhyaAdapterKindIsLiveReadThrough`. |
| M04-C8 | ran | `Sync(ctx)` returns `ports.SyncResult{Processed:N, Errors:0}`; `TestSankhyaSyncMapsSnapshotHonestNull` asserts Processed==3. |
| M04-C9 | ran | Read-side signatures unchanged (same as C3); full-module `go build ./...` green = existing pricing/linking consumers compile. |
| M04-C10 | ran | `git show 18b278dc --stat`: `cache.go` (canonicalKey→activeSourceKey incl sankhya) + `cache_test.go` in ONE commit; source is 3rd key value. |
| M04-C11 | ran | `TestCatalogCachePartitionsSankhyaVsXlsx` (cache_test.go): (tenant,xlsx) populated, (tenant,sankhya) misses; proven load-bearing pre-commit by removing activeSourceKey → test failed ("got 1"), restored → green. |
| M04-C12 | ran | Single commit `18b278dc` (same as C10). |
| M04-C13 | ran | Tenant scoping. **Grep scope = the whole sync flow (oracle read + mirror write), NOT sync.go alone** — a literal grep of `sync.go` finds zero `tenant_id` BY DESIGN and that is not a fail: the Sankhya Oracle tables carry no tenant column (one Oracle connection = one tenant's ERP, `sync.go` header), so the tenant concept exists only at the mirror write. Every `products_mirror`/`_stock_locations` statement in `writer.go` is tenant-keyed: `upsertSQL`/`keepAbsentSQL`/`deleteLocationsSQL`/`insertLocationSQL` all `tenant_id=$1`. Cross-tenant isolation proven by `writer_integration_test.go` (neighbour tenant row untouched by sweep). Both gate + refuter concurred this is architecturally sound, not AC-01. |
| M04-C14 | ran | No `.env`/credential read or print in any commit; no Oracle DSN in transcript/evidence. |

## User-drive (D-120) — hub live-drive REQUESTED (task #8)

| ID | Plan |
|----|------|
| M04-U1 | Flip active_source→sankhya; /catalogo serves live (Oracle reachable) OR honest error (never other-source data / empty). Hub live-drive. |
| M04-U2 | Flip back→xlsx returns byte-consistent upload dataset. **VERIFY hub warning**: xlsx read path serves from xlsx import snapshot tables via routing.Reader, NOT products_mirror; if it reads mirror → cross-adapter collision, STOP+REQUEST. Hub live-drive. |
| M04-U3 | Pre-activation F1 honored: /precos with xlsx source shows no Oracle-origin number (no silent cross-source mix). Hub live-drive. |

## Scope guards honored
- No root.go edits (F1 pre-activation routing = HUB-OWNED, sequenced post M-03+M-04 merge).
- Own mirror writer INSIDE internal_read (hub R-2 (i)), disjoint from M-03 erp_import; no import from erp_import.
- No ImportSource enum extension; no sankhya_linkage_*.go touched; no promo/preco_final (M-06).
- No stub for real-Oracle proof (AC-04); no push (AC-06).
