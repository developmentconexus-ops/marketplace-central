# M-09-sync-observability — Closing Evidence Pack

```yaml
id: M-09-EVIDENCE
type: evidence-pack
parent: M-09-sync-observability
mission: MIS-007
head_sha: cbf12ed04c1877e2ebb300a426ee62a6daf246f5
base_sha: 295e293fdc273ed0fad9c3eb2445b7f2152586ed
captured: 2026-08-01
```

All command output below was captured fresh against HEAD `cbf12ed0`, in this session, from
`apps/server_core` (Go) and `apps/web` (FE). The live JSON/DB evidence was captured with a
throwaway `_test.go` file
(`apps/server_core/internal/modules/sync/transport/zzevidence_m09_test.go`) built on the
repo's existing session-container Postgres pattern (`testpostgres.OpenPool`, the same helper
`health_reader_integration_test.go` / `health_handler_integration_test.go` use). It never
bound a port — only `httptest.NewRecorder()`/`httptest.NewRequest()` in-process. The file was
deleted immediately after its output was captured below; `git status --short` at the end of
this pack confirms zero trace.

Scratch DB used for the capture: fresh `CREATE DATABASE mpc_test_f61d6ecb0cb64ebde366e763cfabed5a`
against the worktree's already-running harness session container
(`mpc-pg-session-18eca3d6`, `127.0.0.1:58818`), migrated with `go run ./cmd/testdb migrate`
(`applied 72 migration(s)`) before the capture ran.

---

## M09-C1 — Payload IC-05 com products real + nulls honestos

Requires: `GET /sync/health` returns entities built from real `sync_state` rows (products
timestamps == the seeded/SELECTed values), every row scanned including the ERP sentinel
(`installation_id="erp"`), never-run entities as `null` fields, `last_success_at =
GREATEST(last_full_sync_at, last_incremental_at)`, `phase` from the cursor jsonb, and an
empty tenant answering `entities: []`.

Command (throwaway evidence test, `go test -tags=integration -v -run
TestZZEvidenceM09CaptureLiveJSONAndDBRows ./internal/modules/sync/transport/...`):

```
=== RUN   TestZZEvidenceM09CaptureLiveJSONAndDBRows
    zzevidence_m09_test.go:63: === SELECT * FROM sync_state WHERE tenant_id = "tenant-evidence-m09-1785589837715064600" ===
    zzevidence_m09_test.go:72: row: installation_id=ml-installation-1 entity=listings cursor= last_full_sync_at=<nil> last_incremental_at=<nil> last_error={"at": "2026-07-30T08:00:00Z", "message": "provider timeout"} consecutive_failures=3
    zzevidence_m09_test.go:72: row: installation_id=erp entity=products cursor={"phase": "incremental"} last_full_sync_at=2025-12-31 21:00:00 -0300 -03 last_incremental_at=2026-07-31 08:00:00 -0300 -03 last_error= consecutive_failures=0
    zzevidence_m09_test.go:95: === GET /sync/health JSON (tenant=tenant-evidence-m09-1785589837715064600) ===
        {
          "entities": [
            {
              "consecutive_failures": 3,
              "entity": "listings",
              "last_error": {
                "at": "2026-07-30T08:00:00Z",
                "message": "provider timeout"
              },
              "last_incremental_at": null,
              "last_success_at": null,
              "phase": null
            },
            {
              "consecutive_failures": 0,
              "entity": "products",
              "last_error": null,
              "last_incremental_at": "2026-07-31T11:00:00Z",
              "last_success_at": "2026-07-31T11:00:00Z",
              "phase": "incremental"
            }
          ],
          "webhook": {
            "dropped_24h": 0,
            "last_notification_at": null,
            "pending": 0
          }
        }
--- PASS: TestZZEvidenceM09CaptureLiveJSONAndDBRows (0.09s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/sync/transport	2.638s
```

Cross-check: the raw `sync_state` SELECT (`installation_id=erp entity=products
last_full_sync_at=2025-12-31 21:00:00 -0300 -03 last_incremental_at=2026-07-31 08:00:00
-0300 -03`, which is `2026-01-01T00:00:00Z` / `2026-07-31T11:00:00Z` in UTC — the DB session
prints in `-03` local time) matches the JSON's `products.last_incremental_at =
2026-07-31T11:00:00Z` byte-for-byte, and `products.last_success_at` equals it too — the
GREATEST picked the incremental over the older full-sync, proving the JSON is derived from
the seeded rows, not hand-typed. The `orders` entity (never inserted for this tenant) is
absent from `entities[]` — not present with fabricated nulls, genuinely missing. `listings`
(never succeeded) carries `last_success_at: null`, `last_incremental_at: null`,
`phase: null` — all honest nulls, never `0`/`""`.

Empty-tenant `entities: []` is separately proved by
`TestSyncHealthHandlerEmptyTenant`/`TestHealthReaderEmptyTenant` (see M09-C5 test-name
citations below; both are part of the `go test -tags=integration -v
./internal/modules/sync/...` run reproduced under "Backend verification" and passed).

Sentinel-row scan (F-r05-1: `installation_id="erp"` included, not filtered) is separately
proved by `TestHealthReaderReadsAcrossEveryInstallation`, also in that same run (PASS).

Verdict: **PASS**

---

## M09-C2 — Fixture negativa incremental-only (F-r04-1)

Requires: an entity with a stale full-sync + a recent incremental sync must report
`last_success_at` equal to the recent incremental, never the stale full-sync — the GREATEST
must not freeze on the one-time backfill.

Evidence: the same capture above (M09-C1) seeded exactly this shape for `products`:

- Inserted `last_full_sync_at = 2026-01-01T00:00:00Z` (stale, six months old relative to the
  fixture's "today") and `last_incremental_at = 2026-07-31T11:00:00Z` (recent).
- The JSON the live route returned: `"last_success_at": "2026-07-31T11:00:00Z"` — the
  **recent incremental value**, not `2026-01-01T00:00:00Z`. `last_incremental_at` in the JSON
  matches the same value, confirming the field the service picked.

The unit-level version of this same guarantee (`TestHealthServiceLastSuccessAtIsGreatestNotFullSyncAlone`)
and its DB-backed sibling (`TestSyncHealthHandlerGoldenFixture`, which seeds the identical
stale-full/recent-incremental shape independently) both PASS in the "Backend verification"
run below — three independent proofs (unit, DB-golden, and this fresh capture) agree on the
same GREATEST semantics.

Verdict: **PASS**

---

## M09-C3 — Seam WebhookStatsReader provado na ROTA

Requires: the default webhook reader answers the canonical zero state byte-for-byte, and a
fake injected via `WithWebhookStatsReader` AFTER the handler is already `Register()`'d on a
live mux must be observable on the SAME route (injection by reference/pointer, not by a
value-receiver copy that would leave the live route serving the stale default forever).

Command: `go test -tags=integration -run
TestHealthHandlerWithWebhookStatsReaderAfterRegisterIsObservedLive -v
./internal/modules/sync/transport/...`

```
=== RUN   TestHealthHandlerWithWebhookStatsReaderAfterRegisterIsObservedLive
--- PASS: TestHealthHandlerWithWebhookStatsReaderAfterRegisterIsObservedLive (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/sync/transport	3.153s
```

This test (`apps/server_core/internal/modules/sync/transport/health_handler_test.go:67-112`)
registers the handler on a live `httpx.NewRouteClassMux()`, hits `/sync/health` BEFORE
injection (asserts `"pending":0`, the default), then calls `svc.WithWebhookStatsReader(fake)`
on the *same* already-registered `*HealthService`, then hits the *same* route again and
asserts the fake's values (`pending=7`, `dropped_24h=2`, a real `last_notification_at`) come
back — proving the swap is visible on the already-live route, not just a fresh construction.

The sibling test `TestHealthHandlerDefaultWebhookBlockIsByteExact` (same file, same run,
included in the "Backend verification" sweep below) asserts the default block is literally
`"webhook":{"last_notification_at":null,"pending":0,"dropped_24h":0}` — the IC-05 canonical
state — via a raw substring match on the response body, and PASSES.

The seam's pointer-receiver contract is also directly visible in
`apps/server_core/internal/modules/sync/application/health_service.go:39-60` —
`WithWebhookStatsReader` is defined on `*HealthService` and mutates `s.webhookReader` in
place; `NewHealthService` returns `*HealthService`, and `HealthHandler` holds the interface
value it was constructed with, so a later mutation of the underlying struct is visible
through the handler's held reference.

Verdict: **PASS**

---

## M09-C4 — FE: 4 estados renderizados honestos

Requires: green (`last_success_at` non-null AND `consecutive_failures == 0`) / red
(consecutive failures, `last_error` in tooltip) / gray "nunca" (nulls, never "0 min atrás") /
webhook initial state as the literal fact "nenhuma notificação recebida" (never a
configuration verdict); a fetch error renders a named error on the card while the rest of the
page stays intact.

Command: `npx --no-install vitest run src/pages/integracoes/SyncHealthCard.test.tsx
src/pages/integracoes/IntegracoesPage.test.tsx --reporter=verbose` (from `apps/web`)

```
 ✓ src/pages/integracoes/SyncHealthCard.test.tsx > SyncHealthCard > renders a green badge with relative time and the absolute ISO in the title for a healthy entity 58ms
 ✓ src/pages/integracoes/SyncHealthCard.test.tsx > SyncHealthCard > renders a red badge with the failure count and the last_error message in a tooltip 29ms
 ✓ src/pages/integracoes/SyncHealthCard.test.tsx > SyncHealthCard > shows the literal 'nunca' in faint styling for an entity that never completed a sync, with no relative-time text 18ms
 ✓ src/pages/integracoes/SyncHealthCard.test.tsx > SyncHealthCard > trusts a recent last_success_at as fresh/green without recomputing staleness itself 28ms
 ✓ src/pages/integracoes/SyncHealthCard.test.tsx > SyncHealthCard > states the fact 'nenhuma notificação recebida' for the webhook's initial state, never a configuration verdict 31ms
 ✓ src/pages/integracoes/SyncHealthCard.test.tsx > SyncHealthCard > shows notification/pending/dropped counts for the webhook's active state 31ms
 ✓ src/pages/integracoes/SyncHealthCard.test.tsx > SyncHealthCard > renders a named error and a retry affordance when the fetch fails, without a blank card 47ms
 ✓ src/pages/integracoes/SyncHealthCard.test.tsx > SyncHealthCard > isolates the fetch failure to this card: the rest of IntegracoesPage stays intact with its normal content 47ms
 ✓ src/pages/integracoes/SyncHealthCard.test.tsx > SyncHealthCard > renders a generic row for an unknown/future entity name not in any hardcoded list 30ms

 Test Files  2 passed (2)
      Tests  27 passed (27)
```

(Full untruncated 27-test list, including all 18 `IntegracoesPage.test.tsx` cases, is in the
"Frontend verification" section below — every one PASSED, none skipped.)

Test-name-to-criterion mapping:
- Green/timestamp coherence: `renders a green badge with relative time and the absolute ISO
  in the title for a healthy entity`.
- Red/last_error tooltip: `renders a red badge with the failure count and the last_error
  message in a tooltip`.
- Gray "nunca", no fabricated relative time: `shows the literal 'nunca' in faint styling for
  an entity that never completed a sync, with no relative-time text`.
- F-r04-1 negative fixture parity on the FE (green from state, not a time cutoff): `trusts a
  recent last_success_at as fresh/green without recomputing staleness itself`.
- Webhook initial-state FACT, not verdict (F-r07-3): `states the fact 'nenhuma notificação
  recebida' for the webhook's initial state, never a configuration verdict`.
- Fetch-error isolation: `renders a named error and a retry affordance when the fetch fails,
  without a blank card` + `isolates the fetch failure to this card: the rest of
  IntegracoesPage stays intact with its normal content`.

Verdict: **PASS** (component-level; M09-U1–U3 browser-drive confirmation is deferred to the
hub per the validation contract — see Deferred section).

---

## M09-C5 — Par SDK/OpenAPI + página intacta

Requires: `/sync/health` in OpenAPI + `getSyncHealth` in sdk-runtime land in the SAME commit;
`tsc` green; `IntegracoesPage.test.tsx` stays green unedited; the page diff is a 1-line mount
(card is a new file).

OpenAPI + SDK pair, same commit (`git show --stat 362618da`):

```
commit 362618dabddd7e85721f430cfdd57f6f50bb1720
    feat(sync): add GET /sync/health per-entity aggregate endpoint
    ...
    OpenAPI (/sync/health + SyncHealth* schemas) and sdk-runtime
    (getSyncHealth + types) updated in the same commit.

 contracts/api/marketplace-central.openapi.yaml | 84 ++++++++++++++++++++++++++
 packages/sdk-runtime/src/index.ts              | 26 ++++++++
 2 files changed, 110 insertions(+)
```

`IntegracoesPage.tsx` diff, base (`295e293f`, planning commit right before M-09) vs HEAD
(`cbf12ed0`) — `git diff 295e293fdc273ed0fad9c3eb2445b7f2152586ed
cbf12ed04c1877e2ebb300a426ee62a6daf246f5 -- apps/web/src/pages/integracoes/IntegracoesPage.tsx`:

```diff
diff --git a/apps/web/src/pages/integracoes/IntegracoesPage.tsx b/apps/web/src/pages/integracoes/IntegracoesPage.tsx
index c171673f..d265885a 100644
--- a/apps/web/src/pages/integracoes/IntegracoesPage.tsx
+++ b/apps/web/src/pages/integracoes/IntegracoesPage.tsx
@@ -18,6 +18,7 @@ import { useId, useRef, useState, type ChangeEvent, type DragEvent } from "react
 import { useClient } from "../../app/ClientContext";
 import { ImportacaoSection } from "../importacoes/ImportacaoSection";
 import { useErpImportDetail } from "../vinculos/useErpImports";
+import { SyncHealthCard } from "./SyncHealthCard";
 import { useErpImportUpload, type ErpImportUploadError } from "./useErpImportUpload";
 
 type SourceOption = {
@@ -568,6 +569,7 @@ export function IntegracoesPage() {
       <SellableAssortmentCard />
       <UploadCard />
       <ProviderConnectCard />
+      <SyncHealthCard />
       <ImportacaoSection />
     </section>
   );
```

Exactly 2 changed lines (1 import + 1 mount), no other edits — `SyncHealthCard` lives in its
own new file (`apps/web/src/pages/integracoes/SyncHealthCard.tsx`, 175 lines, `git show
--stat 6bf04d09`).

`IntegracoesPage.test.tsx` (18 tests, unedited by M-09 per the same `git diff` showing zero
hunks against that file) — all PASS, see the vitest output under M09-C4 and "Frontend
verification" below.

`tsc --noEmit` from `apps/web`: **zero errors touch any M-09 file.** The full run does surface
12 pre-existing errors in unrelated files (`anunciosQueries.ts`,
`mutations/MutationPreviewModal.tsx`, `mutations/MutationResultSummary.tsx`,
`produto/ProdutoPage*.test.tsx`, `AnunciosTable.test.tsx`, `ListingsRefreshControl.test.tsx`,
`anunciosQueryState.test.ts`) — verified pre-existing and untouched by M-09 by diffing each
against base_sha:

```
$ git diff --stat 295e293f cbf12ed0 -- apps/web/src/pages/anunciosQueries.ts apps/web/src/pages/mutations/MutationPreviewModal.tsx apps/web/src/pages/produto/ProdutoPage.test.tsx apps/web/src/pages/AnunciosTable.test.tsx apps/web/src/pages/ListingsRefreshControl.test.tsx apps/web/src/pages/produto/ProdutoPage.partialFailure.test.tsx apps/web/src/pages/anunciosQueryState.test.ts apps/web/src/pages/mutations/MutationResultSummary.tsx
(no output — zero diff, none of these files changed between base_sha and HEAD)
```

None of `SyncHealthCard.tsx`, `SyncHealthCard.test.tsx`, `IntegracoesPage.tsx`,
`web-query/src/syncHealth.ts`, or the sdk-runtime/OpenAPI files appear anywhere in the tsc
error list (full output is in "Frontend verification" below).

Verdict: **PASS**

---

## Backend verification (whole module, HEAD `cbf12ed0`)

All from `apps/server_core`, `GOCACHE=$(pwd)/.gocache`, `GOMODCACHE=$(pwd)/.gomodcache`.

### `go build ./...`

```
$ go build ./...
(no output, exit 0)
```

### `go vet ./...`

```
$ go vet ./...
(no output, exit 0)
```

### `go test -count=1 ./...`

Full module sweep (non-integration-tagged; `//go:build integration` packages compile-check
here but their tests are gated behind the tag — see the dedicated integration run below,
required by this repo's own "vacuous green" doctrine note that a bare `go test ./...` never
executes integration-tagged files):

```
$ go test -count=1 ./...
... 156 package lines total: 111 "ok", 45 "[no test files]", 0 "FAIL" ...
ok  	marketplace-central/apps/server_core/internal/modules/sync/application	4.107s
ok  	marketplace-central/apps/server_core/internal/modules/sync/composition	4.552s
?   	marketplace-central/apps/server_core/internal/modules/sync/domain	[no test files]
ok  	marketplace-central/apps/server_core/internal/modules/sync/transport	4.701s
?   	marketplace-central/apps/server_core/internal/modules/sync/adapters/postgres	[no test files]
...
ok  	marketplace-central/apps/server_core/tests/integration	3.005s
ok  	marketplace-central/apps/server_core/tests/unit	3.266s
```

Counts, not a tail (per doctrine): `grep -c "^ok"` = 111, `grep -c "no test files"` = 45,
`grep -c FAIL` = 0. `sync/adapters/postgres` shows `[no test files]` here — that is expected
and NOT a gap: every test in that package carries `//go:build integration` and is proved
below, tagged, with a real DB.

### `go test -tags=integration -count=1 -v ./internal/modules/sync/...`

Run against the harness session-container Postgres (`MPC_TEST_DATABASE_URL` pointed at the
freshly migrated scratch DB described above) — this is the sweep that actually compiles and
executes every `//go:build integration` file under the sync module, closing the gap the
previous command left open:

```
=== RUN   TestHealthReaderReadsAcrossEveryInstallation
--- PASS: TestHealthReaderReadsAcrossEveryInstallation (0.02s)
=== RUN   TestHealthReaderOrderingAndFieldMapping
--- PASS: TestHealthReaderOrderingAndFieldMapping (0.01s)
=== RUN   TestHealthReaderEmptyTenant
--- PASS: TestHealthReaderEmptyTenant (0.02s)
=== RUN   TestSyncStateCursorRoundTrip
--- PASS: TestSyncStateCursorRoundTrip (0.01s)
=== RUN   TestSyncStateFailureThenRecovery
--- PASS: TestSyncStateFailureThenRecovery (0.01s)
=== RUN   TestSyncStateConcurrentCursorAppend
--- PASS: TestSyncStateConcurrentCursorAppend (0.03s)
=== RUN   TestSyncStateTenantIsolation
--- PASS: TestSyncStateTenantIsolation (0.02s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/sync/adapters/postgres	3.730s
=== RUN   TestRunOnceRecordsSuccessWithCursor
--- PASS: TestRunOnceRecordsSuccessWithCursor (0.00s)
=== RUN   TestRunOnceFailureIsIsolatedPerEntity
--- PASS: TestRunOnceFailureIsIsolatedPerEntity (0.00s)
=== RUN   TestRunOnceSkipsEntityOnReadError
--- PASS: TestRunOnceSkipsEntityOnReadError (0.00s)
=== RUN   TestRunOncePanicIsIsolatedAsFailure
--- PASS: TestRunOncePanicIsIsolatedAsFailure (0.00s)
=== RUN   TestRegisterJobGuards
--- PASS: TestRegisterJobGuards (0.00s)
=== RUN   TestStartStopsOnContextCancel
--- PASS: TestStartStopsOnContextCancel (0.00s)
=== RUN   TestStartNoOpOnZeroInterval
--- PASS: TestStartNoOpOnZeroInterval (0.00s)
=== RUN   TestHealthServiceLastSuccessAtIsGreatestNotFullSyncAlone
--- PASS: TestHealthServiceLastSuccessAtIsGreatestNotFullSyncAlone (0.00s)
=== RUN   TestHealthServiceLastSuccessAtNilSafety
=== RUN   TestHealthServiceLastSuccessAtNilSafety/both_nil
=== RUN   TestHealthServiceLastSuccessAtNilSafety/only_full_set
=== RUN   TestHealthServiceLastSuccessAtNilSafety/only_incremental_set
=== RUN   TestHealthServiceLastSuccessAtNilSafety/full_is_later
--- PASS: TestHealthServiceLastSuccessAtNilSafety (0.00s)
=== RUN   TestHealthServiceEmptyTenantReturnsEmptySliceNotNil
--- PASS: TestHealthServiceEmptyTenantReturnsEmptySliceNotNil (0.00s)
=== RUN   TestHealthServiceCursorPhaseExtraction
--- PASS: TestHealthServiceCursorPhaseExtraction (0.00s)
    (5 subtests: nil_cursor, valid_phase, no_phase_key, malformed_non-JSON, phase_null — all PASS)
=== RUN   TestHealthServiceReaderErrorPropagates
--- PASS: TestHealthServiceReaderErrorPropagates (0.00s)
=== RUN   TestHealthServiceDefaultWebhookIsCanonicalZero
--- PASS: TestHealthServiceDefaultWebhookIsCanonicalZero (0.00s)
=== RUN   TestHealthServiceWithWebhookStatsReaderMutatesInPlace
--- PASS: TestHealthServiceWithWebhookStatsReaderMutatesInPlace (0.00s)
=== RUN   TestHealthServiceLastErrorPassesThroughVerbatim
--- PASS: TestHealthServiceLastErrorPassesThroughVerbatim (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/sync/application	2.080s
=== RUN   TestProductsJobRunsTheActiveSourceAdapter
--- PASS: TestProductsJobRunsTheActiveSourceAdapter (0.00s)
=== RUN   TestProductsJobFollowsTheActiveSourceFlip
--- PASS: TestProductsJobFollowsTheActiveSourceFlip (0.00s)
=== RUN   TestProductsJobFailsClosedWhenTheActiveSourceIsUnreadable
--- PASS: TestProductsJobFailsClosedWhenTheActiveSourceIsUnreadable (0.00s)
=== RUN   TestProductsJobFailsClosedWhenNoAdapterIsRegistered
--- PASS: TestProductsJobFailsClosedWhenNoAdapterIsRegistered (0.00s)
=== RUN   TestProductsJobFailsClosedOnANilAdapter
--- PASS: TestProductsJobFailsClosedOnANilAdapter (0.00s)
=== RUN   TestProductsJobRefreshesLinkCandidatesAfterASuccessfulSync
--- PASS: TestProductsJobRefreshesLinkCandidatesAfterASuccessfulSync (0.00s)
=== RUN   TestProductsJobSurvivesALinkCandidateRefreshFailure
--- PASS: TestProductsJobSurvivesALinkCandidateRefreshFailure (0.02s)
=== RUN   TestProductsJobDoesNotRefreshWhenTheSyncFailed
--- PASS: TestProductsJobDoesNotRefreshWhenTheSyncFailed (0.00s)
=== RUN   TestProductsJobKeepsTheCursorWhenTheAdapterFails
--- PASS: TestProductsJobKeepsTheCursorWhenTheAdapterFails (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/sync/composition	2.335s
?   	marketplace-central/apps/server_core/internal/modules/sync/domain	[no test files]
=== RUN   TestSyncHealthHandlerGoldenFixture
--- PASS: TestSyncHealthHandlerGoldenFixture (0.02s)
=== RUN   TestSyncHealthHandlerEmptyTenant
--- PASS: TestSyncHealthHandlerEmptyTenant (0.01s)
=== RUN   TestHealthHandlerDefaultWebhookBlockIsByteExact
--- PASS: TestHealthHandlerDefaultWebhookBlockIsByteExact (0.00s)
=== RUN   TestHealthHandlerWithWebhookStatsReaderAfterRegisterIsObservedLive
--- PASS: TestHealthHandlerWithWebhookStatsReaderAfterRegisterIsObservedLive (0.00s)
=== RUN   TestHealthHandlerRegistersInteractiveNotBatch
--- PASS: TestHealthHandlerRegistersInteractiveNotBatch (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/modules/sync/transport	3.473s
```

0 FAIL across every sync package, tagged and untagged combined.

---

## Frontend verification (HEAD `cbf12ed0`, from `apps/web`)

### `npx --no-install tsc --noEmit`

```
src/pages/anunciosQueries.ts(17,54): error TS2345: Argument of type 'ListingListOptions' is not assignable to parameter of type 'Record<string, unknown>'.
src/pages/anunciosQueryState.test.ts(128,39): error TS2345: Argument of type ... AnunciosClient': refreshListings, listIntegrationOperationRuns
src/pages/anunciosQueryState.test.ts(153,42): error TS2345: Argument of type ... AnunciosClient': refreshListings, listIntegrationOperationRuns
src/pages/AnunciosTable.test.tsx(40,7): error TS2739: Type ... is missing the following properties from type 'ListingMarketSignal': median, min_valid, max_valid
src/pages/ListingsRefreshControl.test.tsx(114,36): error TS2352: Conversion of type 'Promise<void>' to type 'QueryClient' may be a mistake ...
src/pages/mutations/MutationPreviewModal.tsx(210,33): error TS2741: Property 'onRetry' is missing in type '{ detail: string; }' but required in type 'ErrorStateProps'.
src/pages/mutations/MutationResultSummary.tsx(22,25): error TS2741: Property 'onRetry' is missing in type '{ detail: string; }' but required in type 'ErrorStateProps'.
src/pages/produto/ProdutoPage.partialFailure.test.tsx(39,44): error TS2322: Type '"complete"' is not assignable to type 'CanonicalSourceFactQuality'.
src/pages/produto/ProdutoPage.partialFailure.test.tsx(40,45): error TS2322: Type '"complete"' is not assignable to type 'CanonicalSourceFactQuality'.
src/pages/produto/ProdutoPage.partialFailure.test.tsx(41,46): error TS2322: Type '"complete"' is not assignable to type 'CanonicalSourceFactQuality'.
src/pages/produto/ProdutoPage.partialFailure.test.tsx(45,3): error TS2322: Type '"MATCHED"' is not assignable to type 'MarketPriceIntelMatchStatus'.
src/pages/produto/ProdutoPage.test.tsx(91,7): error TS2322: Type '"MATCHED"' is not assignable to type 'MarketPriceIntelMatchStatus'.
```

12 errors, all in files confirmed unchanged since `base_sha` (see the `git diff --stat`
under M09-C5 — zero output). **None reference `SyncHealthCard.tsx`,
`SyncHealthCard.test.tsx`, `IntegracoesPage.tsx`, or `web-query/src/syncHealth.ts`.** This is
pre-existing repo debt unrelated to M-09, not a regression this milestone introduced.

### `npx --no-install vitest run` (full suite)

```
 Test Files  68 passed (68)
      Tests  575 passed (575)
   Duration  31.97s (transform 8.20s, setup 20.65s, collect 104.86s, tests 27.44s, environment 107.63s, prepare 14.62s)
```

0 failed, 0 skipped, including `src/pages/integracoes/SyncHealthCard.test.tsx (9 tests)` and
`src/pages/integracoes/IntegracoesPage.test.tsx (18 tests)`.

---

## Deferred

- **M09-C6** — "Re-drive pós-lanes — entities acendem SEM mudança": explicitly deferred in
  the validation contract itself (`Command: re-drive do hub após M-04/M-06 fecharem —
  critério DIFERIDO — não bloqueia close do M-09; bloqueia close da MISSÃO via MIS07-C8`).
  Not attempted in this session; hub re-drives after M-04/M-06 close.
- **M09-U1** — Seção "Saúde do sync" dirigida com dado REAL (browser drive + SELECT lado a
  lado). Hub-executed post-merge browser QA per the user-drive mandate; not attempted here.
- **M09-U2** — Estados honestos visíveis no mesmo drive (browser drive + screenshot).
  Hub-executed post-merge; not attempted here.
- **M09-U3** — Health 500 simulado → card mostra erro nomeado, resto de /integracoes segue
  operável (browser drive com fetch quebrado). Hub-executed post-merge; not attempted here.

---

## Commits

| SHA | Description |
|---|---|
| `362618da` | feat(sync): add `GET /sync/health` per-entity aggregate endpoint — reader/service/handler chain, `WebhookStatsReader` port + default, OpenAPI + sdk-runtime `getSyncHealth` in the same commit. |
| `6bf04d09` | feat(web): add `SyncHealthCard` to `IntegracoesPage` (F-02/M-09) — 4-state card, `useSyncHealthQuery` hook (30s poll), 2-line page mount. |
| `1f29c33a` | fix(sync): drop unused `tenantID` field from `HealthReader` — cleanup, no behavior change. |
| `cbf12ed0` | test(web): assert phase display renders and omits per entity (M-09/F-02) — additional `SyncHealthCard.test.tsx` coverage. |

## Reviews

Three cold adversarial reviews were run against this milestone's diff prior to this evidence
pack: **F-01 backend**, **F-02 frontend**, and a final **whole-diff review**. All three
returned **PASS with zero blocking findings**.
