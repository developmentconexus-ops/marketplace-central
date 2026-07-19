# F-01 listings-signals-api — SLICE PLAN (P2 batch, cold Opus)

```yaml
id: F-01
type: feature-plan
parent: M-05-anuncios-sinais
mission: MIS-004-mvp-demo
base_sha: 8b6c4b3093f9465cd3b91209b054af4fa702171a
planner: P2 BATCH PLANNER (cold Opus)
lane: contingency §12 (implementers = sonnet)
slices: 5
internal_dag: S1 (contract+SDK, independent) ; S2 -> S3 -> S4 -> S5
```

## Scope recap (from feature.md + IC-03 + validation-contract)

Additive competitive-signal enrichment on `/listings*`: per-item `market_signal` object + `signal_status`
(`OK|SEM_VINCULO|NO_PRICE_EVIDENCE|STALE`), summary `exceptions` counters
(`sem_vinculo|abaixo_custo|sem_evidencia`), and `?exception=` filter values — all strictly additive, no
existing `/listings*` consumer edited. Cross-module read ONLY via `market.EvidenceReader` Go port
(`internal/modules/market/ports/evidence_reader.go:12`), never HTTP self-call, never cross-schema SQL.
Zero migration (computed projection over already-fetched `ListingReadModel.Link.ProductID`).

## Command reference (HARNESS-PROFILE §3; split build/vet/test — false-alarm signature)

- GO_BUILD: `cd apps/server_core && GOCACHE="$(pwd)/.gocache" go build ./...`
- GO_VET:   `cd apps/server_core && GOCACHE="$(pwd)/.gocache" go vet ./...`
- GO_TEST_LISTINGS: `cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/modules/listings/...`
- SDK_TYPECHECK: `cd packages/sdk-runtime && npx tsc --noEmit`
- SDK_TEST: `cd packages/sdk-runtime && npx vitest run`
- GOVERNANCE (L0, clean worktree, 40-hex BaseSha): `npm run harness:governance` (run from a CLEAN
  worktree checkout — false-fails if run against the hub checkout because it scans `.claude/worktrees`).

---

## SLICE F01-S1 — Contract + SDK surface (additive OpenAPI delta + new listings.ts)

- **id:** F01-S1
- **complexity:** standard
- **goal:** Land the additive `/listings*` OpenAPI delta and the NEW hand-written `packages/sdk-runtime/src/listings.ts`
  in the SAME commit (invariant: spec + SDK land together). Give F-02 a typed surface for `market_signal` /
  `signal_status` / extended `exceptions` WITHOUT editing the inline `index.ts` client (only BARREL-M05 line).
- **failing_test_first:** NEW `packages/sdk-runtime/src/listings.test.ts` (vitest): builds a stub client whose
  `listListings` returns a fixture with `market_signal` (OK-shape) and one `SEM_VINCULO` item (`market_signal: null`,
  `signal_status: "SEM_VINCULO"`); passes it through `withListingSignals(client)` and asserts the wrapper re-types
  items so `.market_signal` / `.signal_status` are accessible and the extended `summary.exceptions` carries
  `sem_vinculo/abaixo_custo/sem_evidencia`. Fails to compile/run before listings.ts exists.
- **validation_kind:** L0 (SDK_TYPECHECK + SDK_TEST) + GOVERNANCE (spec additive-only diff).
- **commands:** SDK_TYPECHECK ; SDK_TEST ; GOVERNANCE
- **expected_artifacts:** transcript of SDK_TEST green; governance additive-diff pass; `git diff` on openapi.yaml
  showing only NEW optional properties / enum-value additions (no removals/renames).
- **write_set:**
  - `contracts/api/marketplace-central.openapi.yaml` — `/listings` + `/listings/by-product`: extend `filter.exception`
    enum (lines 556 & 654) additively `[sync_error, stale, unlinked, below_margin, sem_vinculo, abaixo_custo, sem_evidencia]`;
    `ListingReadModel` schema (2895): add optional `market_signal` object + `signal_status` enum property (NOT added to
    `required` at 2897); `ListingSummaryExceptions` schema (3036): add optional `sem_vinculo`/`abaixo_custo`/`sem_evidencia`
    integer counters (NOT added to `required` at 3038).
  - `packages/sdk-runtime/src/listings.ts` — NEW (D-30 standalone pattern). Imports base `ListingReadModel`,
    `ListingListOptions`, `ListingPage`, `ListingSummary`, `Client` from `./index`; declares `SignalStatus`,
    `ListingSignalEvidence` (`source`, `fetched_at`, `freshness`), `ListingMarketSignal` (`status`, `position`,
    `price_to_win`, `delta_pct`, `match_status`, `n_offers`, `n_sellers`, `evidence`), `ListingWithSignal extends
    ListingReadModel` (`+ market_signal: ListingMarketSignal | null; signal_status: SignalStatus`),
    `ListingPageWithSignals`, `ListingSummaryExceptionsWithSignals`, and the client wrapper
    `withListingSignals(client: Client)` re-typing `listListings`/`listListingsByProduct`/`getListingsSummary`
    responses (no duplicate fetch, base client untouched).
  - `packages/sdk-runtime/src/listings.test.ts` — NEW.
  - `packages/sdk-runtime/src/index.ts` — BARREL-M05: exactly ONE additive line `export * from "./listings";`.
- **done_criteria:** SDK_TEST green; `tsc --noEmit` clean; openapi diff strictly additive (governance green);
  index.ts touched only on the single barrel line; no fetch logic duplicated.
- **open_questions:** none.

---

## SLICE F01-S2 — Listings domain: signal types + status/exception derivation (pure, no port)

- **id:** F01-S2
- **complexity:** standard
- **goal:** Add listings-owned domain types and PURE derivation functions for `signal_status`, `market_signal`
  mapping, and the `abaixo_custo` rule — no port dependency yet, all additive on `ListingReadModel`.
- **failing_test_first:** NEW `internal/modules/listings/domain/signal_test.go` — table test over `DeriveSignalStatus`:
  (a) `Link.ProductID == nil` ⇒ `SEM_VINCULO` + nil signal; (b) linked but no aggregate/signal ⇒ `NO_PRICE_EVIDENCE`;
  (c) signal `FetchedAt` older than `signalStaleTTL` ⇒ `STALE` (value + age retained, never hidden); (d) fresh signal
  present ⇒ `OK`. Plus `belowCost` unit: Cost known AND target(=price_to_win, fallback winner) < Cost ⇒ true;
  Cost nil ⇒ false (custo desconhecido excluded, ADR-17).
- **validation_kind:** L0 (GO_BUILD + GO_VET) + L1 (GO_TEST_LISTINGS).
- **commands:** GO_BUILD ; GO_VET ; GO_TEST_LISTINGS
- **expected_artifacts:** GO_TEST_LISTINGS transcript green; existing listings unit tests still green (additive proof —
  new `ListingReadModel` fields default to zero value so struct-equality tests unaffected).
- **write_set:**
  - `internal/modules/listings/domain/signal.go` — NEW: `SignalStatus` enum (`OK|SEM_VINCULO|NO_PRICE_EVIDENCE|STALE`);
    `MarketSignal` struct (Status, Position *{Rank,Total}, PriceToWin *Money, DeltaPct *string, MatchStatus, NOffers,
    NSellers, Evidence{Source, FetchedAt, Freshness}); `const signalStaleTTL = time.Hour` with a comment: mirrors market
    `snapshotTTL` (`market/application/collection_pipeline_service.go:19`, `ExpiresAt = fetchedAt + snapshotTTL`) — the
    IC-03 EXPIRED horizon; listings owns the composite per IC-03 line 29 (reader port does NOT expose `expires_at`, so age
    is derived from `FetchedAt`); pure `DeriveSignalStatus(link, signal, aggregate, now)` and `belowCost(cost, target)`.
  - `internal/modules/listings/domain/read_model.go` — additive fields on `ListingReadModel` (:119):
    `MarketSignal *MarketSignal`, `SignalStatus SignalStatus` (zero-value = empty until enrichment sets it).
  - `internal/modules/listings/domain/filter.go` — extend `FilterKeys` + `SetFilterValue` exception validation (:46)
    additively to accept `sem_vinculo|abaixo_custo|sem_evidencia` alongside the existing 4; add matching
    `ListingException*` consts.
- **done_criteria:** new table test green; all pre-existing listings unit tests green unchanged; no existing symbol
  renamed/removed; `go vet` clean.
- **open_questions:** none. (TTL value resolved: `time.Hour`, mirroring market `snapshotTTL`.)

---

## SLICE F01-S3 — ReadService signal enrichment via EvidenceReader (batch + failure isolation)

- **id:** F01-S3
- **complexity:** complex
- **goal:** Enrich List/ByProduct/Get items with real competitive signals through `market.EvidenceReader`, via an
  ADDITIVE constructor so none of the ~40 existing `NewReadService` callers (many out-of-ownership:
  `modules/mutations/integration/lifecycle_test.go`, `tests/integration/*`, repo integration tests) need editing.
- **failing_test_first:** extend `internal/modules/listings/application/read_service_test.go` with a fake
  `EvidenceReader`: (a) strong signal ⇒ `OK` + full evidence (`source/fetched_at/n_offers/n_sellers/match_status`);
  (b) linked, aggregate `NO_PRICE_EVIDENCE` ⇒ item `NO_PRICE_EVIDENCE`, not `SEM_VINCULO`; (c) unlinked ⇒ `SEM_VINCULO`,
  `market_signal` nil, reader NOT queried for that listing; (d) reader returns error ⇒ items degrade to
  `NO_PRICE_EVIDENCE`, telemetry counter incremented, List still returns 200 (error surfaced via telemetry/log, NOT
  swallowed silently, NOT 500); (e) age > TTL ⇒ `STALE`.
- **validation_kind:** L0 (GO_BUILD + GO_VET) + L1 (GO_TEST_LISTINGS).
- **commands:** GO_BUILD ; GO_VET ; GO_TEST_LISTINGS
- **expected_artifacts:** read_service_test transcript covering all 5 cases; grep proof of zero HTTP self-call to
  `/market` in listings (C01 blocking-failure guard).
- **write_set:**
  - `internal/modules/listings/application/read_service.go` — add `NewReadServiceWithEvidence(repo, facts, policies,
    installations, now, evidence market.EvidenceReader)` (existing `NewReadService` retained, delegates with a no-op/absent
    evidence path used ONLY by legacy callers/tests — the LIVE path is always wired with the real reader at ROOT-M05, so no
    permanent nil on the production path); add an `enrichSignals(ctx, items, op)` step invoked from List/ByProduct/Get:
    collect `Link.ProductID` (codprods) + listing IDs from already-fetched rows, batch `evidence.Signals(listingIDs)`,
    `evidence.Aggregates(codprods)`, `evidence.Verdicts(codprods)` (input-order), map to `MarketSignal` +
    `DeriveSignalStatus`; per-listing failure isolation: on reader error set `NO_PRICE_EVIDENCE` + telemetry, List stays 200.
  - `internal/modules/listings/ports/` — NEW minimal consumer port interface (listings-side) mirroring the three
    `EvidenceReader` methods it consumes, so listings depends on its own port (hexagonal), wired to market via the D-21
    composition adapter in S5 (named second consumer = listings ReadService itself — not speculative).
- **done_criteria:** all 5 fake-reader cases green; existing List/ByProduct/Get/Summary tests green; no market package
  imported into listings application (only the listings-side port); List never 500s on reader error.
- **open_questions:** none.

---

## SLICE F01-S4 — Summary exception counters + abaixo_custo + exception filter application

- **id:** F01-S4
- **complexity:** standard
- **goal:** Roll up `sem_vinculo/abaixo_custo/sem_evidencia` into the summary and apply the 3 new `?exception=` filter
  values, reusing the existing exception-filter seams (SQL predicate for link-state; computed scan-and-filter for the
  cost/evidence-derived ones — same machinery as `below_margin`).
- **failing_test_first:** extend `read_service_test.go` + transport test: summary counts exact per a 4-item seed covering
  all 4 statuses (custo-desconhecido item EXCLUDED from `abaixo_custo`); `?exception=abaixo_custo` returns only items with
  Cost known AND target < Cost; `?exception=sem_evidencia` returns only `NO_PRICE_EVIDENCE`; `?exception=sem_vinculo`
  returns only unlinked; invalid `?exception=` ⇒ 422 (existing behavior preserved).
- **validation_kind:** L0 (GO_BUILD + GO_VET) + L1 (GO_TEST_LISTINGS).
- **commands:** GO_BUILD ; GO_VET ; GO_TEST_LISTINGS
- **expected_artifacts:** summary + filter test transcripts; JSON body showing `exceptions` counters; regression proof
  W1 summary response is a subset of the extended one (additive).
- **write_set:**
  - `internal/modules/listings/application/read_service.go` — `Summary` (:56) rolls up the 3 new counters over enriched
    items; `abaixo_custo = belowCost(item.Cost, target)` (target = signal `price_to_win`, fallback `winner_price`);
    generalize `needsBelowMarginScan` (:493) to `needsExceptionScan` so `abaixo_custo`/`sem_evidencia` reuse the capped
    scan-and-filter loop (`maxBelowMarginScanPages`); `sem_vinculo` reuses the link-state SQL predicate.
  - `internal/modules/listings/adapters/postgres/repository.go` — additive `case` arms in the two `switch q.Filter.Exception`
    blocks (:99, :274) mapping `sem_vinculo` to the existing unlinked link-state predicate (SQL-filterable); `abaixo_custo`
    and `sem_evidencia` fall through to the service-layer computed scan (no SQL arm, same as `below_margin`).
  - `internal/modules/listings/transport/http_handler.go` — extend `listingSummaryExceptions` envelope (add 3 counters)
    and item serialization to carry `market_signal`/`signal_status` (additive; existing fields untouched).
- **done_criteria:** counts + filters exact per seed; custo-desconhecido excluded; invalid exception still 422; W1 JSON
  is a valid subset of extended JSON; L1 green.
- **open_questions:** none.

---

## SLICE F01-S5 — ROOT-M05 composition wiring (adapter + live path)

- **id:** F01-S5
- **complexity:** standard
- **goal:** Wire the real `market.EvidenceReader` into the listings `ReadService` at the composition root via the D-21
  adapter pattern, so the live `/listings*` path always enriches (no permanent nil/stub on production).
- **failing_test_first:** `internal/composition/*` wiring/smoke test (or the existing composition test) asserting the
  listings read handler is constructed with a non-nil evidence dependency and that a booted-server smoke List response
  carries the `signal_status` field; if no composition-level test exists, the failing test is the L1 listings
  end-to-end contract test (`TestListingsReadContractEndToEnd`, allowlisted-flaky per profile) extended to assert
  `signal_status` present.
- **validation_kind:** L0 (GO_BUILD + GO_VET) + L1 (GO_TEST_LISTINGS) + GOVERNANCE.
- **commands:** GO_BUILD ; GO_VET ; GO_TEST_LISTINGS ; GOVERNANCE
- **expected_artifacts:** `git diff` root.go/market_adapters.go showing only additive imports + the LISTINGS wiring
  region changed; booted List response JSON with `signal_status`.
- **write_set:**
  - `internal/composition/market_adapters.go` — NEW composition-owned adapter `listingsEvidenceAdapter` implementing the
    listings-side evidence port over the market `EvidenceReader`, with `var _ listingsports.EvidenceReader =
    (*listingsEvidenceAdapter)(nil)` and factory `newListingsEvidenceAdapter(marketEvidenceSvc)` (D-21).
  - `internal/composition/root.go` — ROOT-M05: additive imports + in the LISTINGS wiring region (:586-597) swap
    `NewReadService(...)` for `NewReadServiceWithEvidence(..., newListingsEvidenceAdapter(marketEvidenceSvc))`
    (`marketEvidenceSvc` already exists at :580). Zero edits to any other module's wiring.
- **done_criteria:** live List path enriches; governance clean; diff confined to the two composition files, additive
  only; no nil evidence on the production path.
- **open_questions:** none.

## Seam note (F-01)

Passing the already-constructed `marketEvidenceSvc` (`root.go:580`) through the D-21 adapter satisfies the enrichment
dependency with zero new module construction. `modules.json` needs NO new entry (listings module already registered —
no `GOV_MODULE_COVERAGE` hit). See P2-plan-outputs.md §Additive Locks + §Seam-Closure.
