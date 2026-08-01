# CHIP-M04 Evidence Pack

```yaml
milestone: M-04-listings-backfill-ingest
chip: CHIP-M04
tip_sha: c251dfae80d1738088ebf7115cf276f4cb2bc0eb
branch: claude/elegant-babbage-fa9362
build: go build ./... -> clean (GOCACHE absolute)
vet: go vet ./... -> clean
test: go test ./... -> 118 ok / 0 FAIL (main baseline 117/0; +1 = new listings/composition package)
pushed: false
merged: false
```

## Features closed

- F-01 listings-ddl (migrations 0090-0092: E3 columns, listing_variations, indices)
- F-02 mass-closure-replacement (absent≠closed resemantization, repository.go)
- F-03 backfill-cursor-ingest (resumable enumerator/hydrator, ingestBatch single writer)
- F-04 scheduler-refresh-wiring (per-installation scheduler, BackfillRunner, root.go rewiring)
- Fix (post cold-gate): ADR-13 SellerSKU/EAN content regression for no-variation multiget items
  (commits `8c9a27eb`, `c251dfae`)

## Criterion-by-criterion (validation-contract.md)

### M04-C2 — Abort pós-página-1 → ZERO rows flipped closed (R-B) — MEASURED PASS
- Test: `TestApplyCompletedPull_AbortAfterPage1_NeverClosesUnseenRows`
  (`internal/modules/listings/adapters/postgres/repository_integration_test.go:124`)
- Fixture: `total=30`, `page1=10` (ids [0,10) — the only page an aborted run ever saw),
  `page2=20` (resume adds ids [10,20); ids [20,30) never seen in any tick) — genuinely
  multi-page, per R-3 (34-listing real account would hide this).
- Result: phase A (page-1-only, no MarkRunComplete) — 0 of 30 rows touched (status or
  absent_since). Phase B (resume + MarkRunComplete) — absent_since set ONLY on the 10 rows
  genuinely never seen ([20,30)); status column untouched. Ran with `-tags integration`
  against ephemeral migrated Postgres this session: PASS.

### M04-C3 — E3 + variations no grão certo — MEASURED PASS
- `TestApplyCompletedPull_VariationsUpsertIdempotently`
  (`repository_integration_test.go:421`) — variation rows upserted idempotently, one row
  per variation.
- `TestApplyCompletedPull_StatusPersistsVerbatimEvenForNovelProviderValues`
  (`repository_integration_test.go:368`) — status stored verbatim, no CHECK/normalization
  (IC-07).
- `AvailableQuantity` grain: variation-level when `len(Variations)>0`, listing-level
  otherwise — enforced in `multiget_mapper.go`'s `MapMultigetItemToListing` and mirrored by
  the oracle (`capability_adapter.go:842,858`).

### M04-C4 — Âncoras de snapshots não-regressivas (ADR-13) — MEASURED PASS (post-fix)
- Regression found by this chip's own cold-gate review: no-variation multiget items fed
  `AbsorbProviderSnapshots` (`product_links/application/import_service.go:84`, sole caller
  `listings/adapters/connectors/source.go:44`) an honest-empty top-level SellerSKU/EAN,
  because `items_multiget_reader.go`'s `mlMultigetItemBody` never typed the wire's
  `seller_sku`/`seller_custom_field`/`attributes` fields (present on the wire, silently
  dropped by `json.Unmarshal` — confirmed against `mlItemResponse`, which types the same
  three fields for the single-item GET). Root cause of the gap: a wrongly-assumed "frozen
  file" comment — verified against `.mnfs/MIS-007-ml-sync/M-01-ml-client-hardening/F-02-items-multiget-raw-dto/feature.md`'s
  own Ownership section, which freezes only `capability_adapter.go`/`refresh_policy.go`.
- Fix: `8c9a27eb` types the 3 fields on `ItemMultigetDTO`; `c251dfae` makes
  `MapMultigetItemToListingSnapshot` derive top-level SellerSKU/EAN unconditionally from the
  item's own fields, matching `mapListing`'s oracle behavior exactly
  (`capability_adapter.go:839-840` vs the new `multiget_mapper.go` code) — content parity,
  not just count parity.
- Test: `TestMapMultigetItemToListingSnapshotFeedsObserverHonestly`
  (`internal/modules/listings/adapters/connectors/multiget_mapper_test.go`) — asserts a
  populated SellerSKU/EAN for a no-variation item, seller_custom_field fallback, and no
  variation-to-parent leakage (per-variation entries keep their own fields only, matching
  the oracle's per-variation loop at `capability_adapter.go:852-860`). Re-run independently
  this session: `--- PASS`.
- `TestGetItemsMultigetMapsTopLevelSellerSKUAndAttributes`
  (`internal/modules/connectors/adapters/mercado_livre/items_multiget_reader_test.go`) —
  re-run independently this session: `--- PASS`.
- Did NOT add a SellerSKU field to `domain.Listing`/`ListingInput` — neither type has one
  today; `migrations/listings_test.go:61` enforces `seller_sku` as a forbidden column on the
  `listings` table; the real surface (`ListingLink.SellerSKU`) is already served entirely by
  `product_link_listing_snapshots` via this same `ListingSnapshot` seam.

### M04-C6 — Paginação provada com fixture >1 página — MEASURED PASS
- `TestBackfillRunnerWalksEveryPageAndMarksRunComplete`
  (`internal/modules/listings/application/backfill_runner_test.go:44`) — 2-page fixture
  (`item-1`→cursor `page-2`, `item-2`→cursor `""`); asserts both pages enumerated
  (`cursors == ["", "page-2"]`), both upserted (`len(upserts)==2`), `MarkRunComplete` called
  exactly once at exhaustion.
- Combined with M04-C2's 30/10/20 fixture at the repository layer — pagination proven at
  both the runner and DB-write layers, never on a single-page fixture (R-3).

### ADR-07 — phase vocabulary untouched
- `scheduler.go` (`internal/modules/listings/composition` — new package, not the pre-existing
  `sync` scheduler's `inferIncremental`) not touched by this milestone; direct read confirms
  `inferIncremental` unchanged — only `"incremental"` resolves `true`, `"sweep"` resolves
  `false` by default.

### Single writer (ADR-04)
- `CompletedPullStore.UpsertPulledRows`/`MarkRunComplete` called ONLY from `ingestBatch`
  (`internal/modules/listings/application/backfill.go:28`). Three callers, one persistence
  path: `NewListingsJob` (scheduler tick), `NewListingIngestor.IngestListing` (single-item),
  `BackfillRunner.Pull` (whole-catalog, shared by `RefreshService` and `ResyncWriter` via the
  `BackfillPuller` interface — `mutations/adapters/listings/resync_writer.go:21`). Grep-proof:
  zero remaining call sites for the old `Ingestion`/`NewSource`/`NewSourceWithObserver` types.

## root.go collision arbitration (M-03 concurrency)

`git diff 21ca3595..HEAD -- internal/composition/root.go` — exactly 3 loci touched:
1. ~line 69: +1 import (`listingscomposition`).
2. ~lines 233-249: `installationResyncWriter` struct field rename
   (`source,store` → `puller mutationslistings.BackfillPuller`) + 1-line `Apply` body change.
3. ~lines 727-762: main composition rewiring (`listingBackfillRunner`, `NewRefreshService`,
   `listingscomposition.NewListingsSchedulers(...).StartAll(ctx)`, `resyncWriter` field).

**Zero overlap with lines 578-603** (M-03's orders region).

## NOT proven by this chip — hub's live-drive (per milestone.md Done Means + validation-contract.md)

- M04-C1 (backfill completo + retomada sem duplicata) — real ML account.
- M04-C5 (scheduler diário real + refresh manual mesmo caminho) — needs either waiting 24h
  (scheduler interval is `24*time.Hour`, `root.go:752`) or forcing a tick/using manual
  refresh to observe a real `sync_state` row.
- M04-U1/U2/U3 (browser drive criteria) — explicitly operator-mandated live-drive, not chip
  scope.

## Cold gate review

Dispatched to `harness:gate-reviewer` (read-only). Findings independently re-verified:
- ADR-13 SellerSKU/EAN regression: CONFIRMED, fixed (this pack, M04-C4).
- Provider-slug mismatch (`"mercadolivre"` vs `"mercado_livre"`): DOWNGRADED — confined to
  an unrelated integration test's self-consistent fixture (`tests/integration/listings_refresh_test.go`),
  not the system convention `scheduler.go`/`scheduler_test.go` actually use.
- ADR-14 "second locus" at root.go:233-249: DOWNGRADED — necessary, minimal companion edit
  (struct field rename forced by `BackfillPuller` replacing `PageSource`/`CompletedPullStore`),
  not an anchored-region violation.
