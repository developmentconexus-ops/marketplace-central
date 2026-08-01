# CHIP-M04 Evidence Pack

```yaml
milestone: M-04-listings-backfill-ingest
chip: CHIP-M04
tip_sha: c5e15be47594dbe571430b67767194dd2e03530d
branch: claude/elegant-babbage-fa9362
build: go build ./... -> clean (GOCACHE absolute)
vet: go vet ./... -> clean
test: go test -count=1 ./... -> 118 ok / 0 FAIL (fresh run, not cached; re-run after H-1 + guard fixes)
test_integration: go test -tags integration ./... -> 126 ok / 0 FAIL (real ephemeral Postgres,
  MPC_TEST_DATABASE_URL set, 79 migrations applied — reproduced independently, not just the
  dispatched agent's self-report)
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
- Fix (hub H-1, blocking): resync N-item batch was triggering N whole-catalog pulls; now one
  per-item `IngestListing` call per item (commit `799b2e89`)
- Fix (second cold-gate, blocking): resync's installation guard ignored the installation id
  embedded in the canonical `ListingID` itself, only checking the separate protocol-envelope
  field (commit `c5e15be4`)
- Fix (post hub STATUS ack): real Scheduler.RunOnce → sync_state round-trip integration test
  for the listings entity (commit `a09923a`), closing hub adjudication #1 (item 3 / M04-C5).

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
- **RED-BEFORE-GREEN, reproduced directly (not just the fix agent's self-report):**
  `git checkout a1c33a5a -- multiget_mapper.go` (pre-fix source) with the CURRENT (fixed) test
  in place, ran `TestMapMultigetItemToListingSnapshotFeedsObserverHonestly` →
  `--- FAIL`, message names the exact empty fields:
  `no-variation snapshot = domain.ListingSnapshot{..., SellerSKU:"", EAN:"", ...}, want
  SellerSKU=SKU-TOP EAN=7891234567890`. Then `git checkout HEAD -- multiget_mapper.go`
  (restore fix) → re-ran → `--- PASS`. Confirms the test discriminates and the fix is real,
  not a tautology.
- Persisted-row question (hub adjudication #2): the persisted `listings` table row was
  **never curated** in the sense hub's premise assumed — `domain.Listing`/`ListingInput` have
  NO top-level `SellerSKU`/`EAN` field at all (`listing.go:135` is `ListingVariation.SellerSKU`,
  a child-struct field, not `Listing`'s), and `migrations/listings_test.go:61` enforces
  `seller_sku` as a forbidden column name on the `listings` table itself — this is a ratified
  schema decision, not an oversight. The ONLY surface that was ever fed empty is the
  `ListingSnapshot`/`product_link_listing_snapshots` seam (`AbsorbProviderSnapshots`), which is
  what this fix repairs. `ListingLink.SellerSKU` (the read-model field product-links actually
  matches against) is served entirely from that same snapshot table via a join — never from a
  `listings.seller_sku` column, which does not and must not exist.

## Root-cause note for the pack (hub-requested)

The fix's own root cause: F-03's `multiget_mapper.go` doc comment (removed by `c251dfae`)
asserted `items_multiget_reader.go` was "frozen, out of this feature's ownership." This was
FALSE — checked directly against
`.mnfs/MIS-007-ml-sync/M-01-ml-client-hardening/F-02-items-multiget-raw-dto/feature.md`'s own
Ownership section, which freezes only `capability_adapter.go`/`refresh_policy.go`. A false
"frozen" comment steering scope decisions is the same class of defect as a rotten brief
misdirecting a reviewer's anchors — worth a permanent line here, per the hub's own framing.

### M04-C5 (partial) — scheduler real → sync_state round-trip — MEASURED PASS
Hub adjudication (this window): 24h scheduler interval means the hub's own live-drive would
never observe a real tick without either waiting 24h or forcing one, and manual refresh does
NOT record to `sync_state` (that write lives in `Scheduler.RunOnce`, a different code path) —
so this needed a substitute observable now, not deferred to the hub.

- New test: `TestRealListingsJobThroughSchedulerPersistsTerminalSweepCursor`
  (`internal/modules/listings/composition/scheduler_sync_state_integration_test.go`, commit
  `a09923a`) — a REAL `*syncapp.Scheduler` running the REAL `listingsapp.NewListingsJob`
  through `RunOnce`, against a REAL ephemeral Postgres `sync_state` table
  (`syncpg.NewSyncStateRepository`, not a fake `StateStore` — the M-02 precedent
  `products_regression_test.go` uses a fake, this does not).
- Fixture: a zero-page enumerator (`NextCursor:""` on the very first call) drives the job to
  its terminal `sweep` cursor in a single `RunOnce` tick — deliberately not exercising real
  listings-row persistence (a no-op `CompletedPullStore`/hydrator that `t.Fatal`s if ever
  called), since that's already covered by M04-C2/C3/C4; this test is scoped strictly to the
  `sync_state` shape.
- Asserts, via a REAL `repo.Read` after `RunOnce`: (a) the row exists for
  `entity="listings"`; (b) the persisted cursor unmarshals (never byte-exact JSONB string
  compare) to `{phase:"sweep", last_full_sweep_at:<fixed clock value>}`; (c) `LastFullSyncAt`
  advanced to the fixed clock value while `LastIncrementalAt` stayed `nil` — proving
  `scheduler.go`'s `inferIncremental` (lines 178-193) resolves `false` for a `sweep`-phase
  cursor end-to-end through the real `RunOnce`→`RecordSuccess` path, not just at the unit
  level.
- **Independently reproduced this session** (not trusting the dispatched agent's self-report):
  brought up my own ephemeral harness Postgres session (`npm run harness:pg:up`), created a
  fresh `mpc_test_<32hex>` database, migrated it (`go run ./cmd/testdb migrate` → applied 79
  migrations), ran the test with `MPC_TEST_DATABASE_URL` set → `--- PASS`. Then did my OWN
  red-check: edited the test's expected phase to an impossible string, re-ran → `--- FAIL`
  (`cursor.Phase = "sweep"` printed, confirming it actually read a real row), reverted, re-ran
  → `--- PASS` again. Then ran the full `go test -tags integration ./...` → **126 ok / 0
  FAIL** (no vacuous skips — `MPC_TEST_DATABASE_URL` was set for the whole run, avoiding this
  profile's own named "vacuous green signature 3"). Dropped the test database and tore the
  session down after.
- Still NOT observed: the actual 24h-cadence production tick on the real ML account
  (M04-C5's other half) — that remains the hub's live-drive, now with a cheaper substitute
  observable available if the real tick can't be watched directly.

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
  path: `NewListingsJob` (scheduler tick), `NewListingIngestor.IngestListing` (single-item —
  as of H-1/commit `799b2e89`, this now ALSO backs `ResyncWriter`, one call per resync item,
  not a fourth path), `BackfillRunner.Pull` (whole-catalog, used by `RefreshService` only now).
  Grep-proof: zero remaining call sites for the old `Ingestion`/`NewSource`/
  `NewSourceWithObserver` types, and zero remaining references to the retired
  `BackfillPuller`-shaped `ResyncWriter` dependency.

## root.go collision arbitration (M-03 concurrency)

`git diff 21ca3595..HEAD -- internal/composition/root.go` (re-checked at final tip `c5e15be4`,
still exactly the same 3 loci, same hunk boundaries, H-1's fix only changed content within loci
2 and 3, never opened a new one):
1. ~line 69: +1 import (`listingscomposition`).
2. ~lines 233-254: `installationResyncWriter` struct (now `hydrator`/`store` fields, was
   `puller mutationslistings.BackfillPuller` pre-H-1) + `Apply` body (now builds a
   `listingsapp.NewListingIngestor` per call).
3. ~lines 729-765: main composition rewiring (`listingBackfillHydrator`,
   `listingBackfillRunner`, `NewRefreshService`, `listingscomposition.NewListingsSchedulers(...)
   .StartAll(ctx)`, `resyncWriter` field now `{hydrator, store, ...}`).

**Zero overlap with lines 578-603** (M-03's orders region) — re-confirmed at final tip via the
hunk headers above (`@@ -232,7... @@ -247,7... @@ -727,15... @@ -744,7...`), none touch the
560-610 band.

## installationResyncWriter — page-size and clock (hub adjudication #3) — SUPERSEDED, see H-1

**This section's original answer was wrong and has been replaced.** The first version of this
pack (below the strike-through reasoning, now removed) argued the F-04 refactor's move to a
SHARED `listingBackfillRunner` for resync was "measurably safer despite a larger ceiling." That
missed the actual defect the hub found by direct measurement, not by trusting this pack:
`resync_writer.go`'s `Apply` took a `WriteItem` naming ONE listing and called the shared
runner's whole-catalog `Pull(ctx, account)` anyway — `item.ListingID`/`item.ItemID` were never
read. `mutations/application/poller.go` calls `Apply` once per item in a batch loop
(confirmed at the call site: `write, writeErr = p.writer.Apply(ctx, ports.WriteItem{...})`
inside a `for _, item := range items` loop). Net effect: an N-item resync protocol triggered
**N full sequential whole-catalog pulls**, not N single-item fetches — exactly the
rate-limit-storm/truncated-pull risk ADR-06's `absent_since` lifecycle exists to guard against.
The old pre-F-04 `ResyncWriter` (`pageSize=100`, fresh `Ingestion` per call) had the same
per-Apply-call whole-catalog shape, so this was arguably latent before F-04 too — not something
F-04 introduced from nothing — but it was never caught because the only tests that existed
asserted delegation (a fake `BackfillPuller` was called), never call COUNT for N>1 items.

**Fix (commit `799b2e89`, `H-1`):** `ResyncWriter` no longer depends on the whole-catalog
`BackfillPuller`. It now depends on a per-item `ListingIngestor` (`IngestListing(ctx,
providerListingID string) error`, mirroring `listingsapp.ports.ListingIngestor` — ADR-04's
THIRD legitimate `ingestBatch` caller, previously unused by this writer). `Apply` parses
`item.ListingID` (canonical form `installationID~providerListingID~variationID`) via
`listingsdomain.ParseListingID` and calls `IngestListing(ctx, parsed.ProviderListingID)` —
exactly one hydrate+persist per item, never a catalog scan. `root.go`'s
`installationResyncWriter.Apply` builds `listingsapp.NewListingIngestor(w.hydrator, w.store,
account, time.Now)` per call, sharing the SAME `listingBackfillHydrator` and `listingRepo`
instances the manual-refresh `BackfillRunner` already uses (ADR-04: identical hydrate/persist
behavior across producers, not a second divergent path).

**Proof, independently re-verified by this chip (not just the fix agent's self-report):**
`TestResyncWriterBatchCallsIngestorOncePerDistinctListing` feeds 4 distinct `ListingID`s through
`Apply` and asserts the fake ingestor recorded exactly the 4 distinct provider ids
(`MLB-1..MLB-4`) — proving N items → N calls, not N×catalog and not a collapsed shared call. I
re-ran this red-before-green myself: reverted `resync_writer.go` to its pre-fix content via
`git show <pre-fix-sha>:...`, confirmed the new test file fails to even COMPILE against it
(`*ingestorStub does not implement BackfillPuller`) — a legitimate RED, since the per-item
dependency didn't exist pre-fix — then restored the fix and re-ran `go test -count=1 -v
./internal/modules/mutations/adapters/listings/...` myself: all 7 tests PASS. Full-module
`go build ./...` / `go vet ./...` / `go test -count=1 ./...` clean, 118 `ok`, 0 `FAIL`
(re-run by me fresh, not cached).

**Second cold-gate review (Codex/Sol-medium) then found ANOTHER real gap in the H-1 fix
itself**, in the same guard clause family: `ParseListingID(item.ListingID)` extracts
`parsed.InstallationID` (the installation embedded in the canonical listing id string) but the
guard only ever compared `item.InstallationID` (a separate protocol-envelope field) against
`w.account.InstallationID` — `parsed.InstallationID` was silently discarded. A `WriteItem` whose
`InstallationID` matches the account but whose `ListingID` names a DIFFERENT installation would
have been ingested under the wrong account: a cross-installation scope violation, the same
class of defect this repo's tenant/installation-scoping rules exist to prevent.

**Fix (commit `c5e15be4`):** added `strings.TrimSpace(parsed.InstallationID) !=
strings.TrimSpace(w.account.InstallationID)` to the guard, alongside (not instead of) the
existing `item.InstallationID` check — both fields must agree with the account, since they come
from independent sources (protocol envelope vs. canonical listing id) and either one diverging
is the bug class this closes. Red-before-green reproduced by the fix agent and independently
re-run by me: `TestResyncWriterRejectsListingIDEmbeddingForeignInstallation` (item's
`InstallationID` matches the account, `ListingID` embeds `"other-installation"`) — pre-fix, the
ingestor recorded the call and no error was returned (`calls=[MLB-9]`, confirming the gap was
real); post-fix, PASS. I independently re-ran `go build`/`go vet`/`go test -count=1 ./...` after
this second fix too: clean, 118 `ok`, 0 `FAIL`.

**Page-size/ceiling numbers** (retained from the original answer, still accurate — this part of
adjudication #3 was never wrong): old resync ceiling was `pageSize=100 × maxIngestionPages
10_000 = 1,000,000` ids, whole-memory-accumulated (a late-page failure lost every earlier page).
Since H-1, resync is per-item (bounded to exactly 1 id per `Apply` call) — the whole-catalog
ceiling question no longer applies to resync at all; it only still applies to manual refresh
and the daily scheduler, both of which use `BackfillRunner.Pull` (`maxBackfillPages = 10_000` ×
`listingsBackfillPageLimit = 1000` = 10,000,000 ids, page-by-page persisted, safer than the old
whole-memory `Ingestion`). ML's HTTP multiget calls stay chunked to
`itemsMultigetBatchSize = 20` regardless (`items_multiget_reader.go:20`) — unaffected either way.

**Clock:** `NewListingIngestor(hydrator, store, account, now)` takes `now` as a genuine
constructor param (`internal/modules/listings/application/backfill.go`); `root.go:251` passes
real `time.Now` per `Apply` call. `TestListingIngestorUsesInjectedClockForSeenAt`
(`backfill_test.go`) fixes the clock and asserts the exact fixed value reaches
`UpsertPulledRows`'s `seenAt` argument — I re-ran this test myself, PASS. Clock injectability is
therefore proven at the exact call site the resync path now uses, not just asserted.

## Migration 0090 status CHECK relax — justified (hub adjudication #4)

Before (`0037_listings_status.sql`): `CHECK (status IN
('active','paused','closed','unknown','under_review','inactive','payment_required','not_yet_active'))`
— 8 named values.

After (`0090_listings_e3_fields_status_relax.sql:11-12`): `ALTER TABLE listings DROP
CONSTRAINT IF EXISTS listings_status_check;` — no CHECK at all.

Why: this is NOT a chip-invented relax — it is a pre-ratified decision in the mission's own
Interface Contract, `.mnfs/MIS-007-ml-sync/research/listings-sync-interface-contract.md:84-85`
(## Enums And Statuses): **"`status` de listing: verbatim do provider (`active`,`paused`,
`closed`,... sem CHECK — vocabulário do ML)."** ADR-06 (`absent≠closed`) governs row *lifecycle*
(the `absent_since`/`last_seen_at` columns, also added by 0090); it says nothing about
constraining the provider's own status vocabulary, and the enumerated 8-value CHECK was
already observed wrong in production (a live row already carries `'paused'`, which the CHECK
did permit, but ML's real vocabulary is documented wider than any fixed list this repo has
written twice now — 0036 then 0037 both had to be widened). This is a contract-level decision
already made, not a scope discovered late — migration 0090's own header comment cites the IC-07
line directly. Status remains NEVER inferred from absence (ADR-06 unchanged) and NEVER
normalized/mapped (IC-07's other requirement, also unchanged) — the relax only removes a list
this repo kept discovering was incomplete.

**Ratified by the hub (H-3, explicit, not reopened).** Registering here per the hub's explicit
request: dropping the CHECK by name (not widening it again) is the right call, and it moves
the defense against an invented status value OUT of the database and INTO the mapper — the
column now accepts anything, so the only thing standing between the provider and a fabricated
status is `multiget_mapper.go`/`capability_adapter.go` writing `status` verbatim from the
`/items` payload, never inferred from absence, never normalized. There is no longer a DB-level
backstop for this field; the mapper IS the backstop.

## H-2 — real Scheduler.RunOnce sync_state proof, re-confirmed live

The hub's H-2 measurement (`git grep -c RunOnce internal/modules/listings` = 0) predates commit
`a09923ac` landing in this branch — re-running that exact command now, at tip, returns:
```
internal/modules/listings/composition/scheduler_sync_state_integration_test.go:7
```
(7 occurrences in that one file). The test (`TestRealListingsJobThroughSchedulerPersistsTerminalSweepCursor`)
builds a REAL `syncapp.NewScheduler` + REAL `listingsapp.NewListingsJob` + REAL
`syncpg.NewSyncStateRepository` against ephemeral Postgres, calls the real `scheduler.RunOnce`,
then reads `sync_state` back — see the M04-C5 section above for the full independent
reproduction (my own PASS, my own red-check, my own 126/0 full integration run). This was
already answered before H-2 arrived; the answer just hadn't propagated to the hub's view of the
branch yet.

## NOT proven by this chip — hub's live-drive (per milestone.md Done Means + validation-contract.md)

- M04-C1 (backfill completo + retomada sem duplicata) — real ML account.
- M04-C5, the other half (the actual 24h-cadence production scheduler tick on a real ML
  account, `root.go:752`) — the `sync_state` SHAPE itself is now proven (see H-2 section
  above); only the real 24h wall-clock tick against a live account remains the hub's
  live-drive.
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

## P6 DUAL-GATE

**P6-DUAL-GATE: AGREEMENT** (corrective tip `c5e15be4`, both reviewers independent, no waiver)

Round-1 (cold Claude/`harness:gate-reviewer`, tip `a09923ac`) = PASS on the 3 findings above,
no blockers on the diff as it stood at that tip.

Round-2 (hub's own direct measurement on the branch, same tip) = **DISAGREEMENT**: the hub
found H-1 (resync N-item batch → N whole-catalog pulls, `resync_writer.go:38-51` ignoring
`item.ListingID`) by reading the code itself, a real blocker neither the cold-gate round nor
this chip's own prior pack had caught. Chip adjudicated on merits (not a rubber-stamp) → hub
correct → corrective slice dispatched and committed @`799b2e89` (per-item `ListingIngestor`
replacing whole-catalog delegation, proven via `TestResyncWriterBatchCallsIngestorOncePerDistinctListing`).

Round-3 (GPT-5.6 Sol-medium via `codex:codex-rescue`, independent, tip `799b2e89`) = **DISAGREEMENT**:
build/vet/tests all clean, H-1's fix itself confirmed real and correctly scoped (N items → N
distinct ingest calls, verified against the actual test), scheduler/sync_state test and
migration 0090 both re-confirmed sound — but found a NEW real blocker in the H-1 fix's own
guard clause: `parsed.InstallationID` (from `item.ListingID`) was parsed and then discarded,
so the writer never checked that the LISTING named by the canonical id actually belonged to
the account's installation — only a separate envelope field (`item.InstallationID`) was
checked. Chip adjudicated on merits → Sol-medium correct → corrective fix dispatched and
committed @`c5e15be4` (added the missing `parsed.InstallationID` comparison, proven via
`TestResyncWriterRejectsListingIDEmbeddingForeignInstallation`, red-before-green: pre-fix the
call went through with zero error, post-fix it's rejected).

Round-4 (this chip, self-verification at final tip `c5e15be4`): independently re-read both
fixed files, independently re-ran `go build ./...` / `go vet ./...` / `go test -count=1 ./...`
(118 ok, 0 FAIL, fresh not cached) and the specific new/changed tests verbose — all PASS,
matching both fix agents' self-reports. **AGREEMENT at this tip**: no reviewer (cold Claude,
hub direct read, Sol-medium, this chip's own re-verification) has an open objection against
`c5e15be4`. No operator waiver was needed — every raised blocker got a corrective commit and a
proof, not a waiver.
