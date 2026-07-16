# F-02 Slice 9 — implementation plan (grow canonical listing-status enum)

```yaml
author: GPT-5.6 Sol --effort medium (feature planner, codex thread 019f6961-fd63-7a30-a3a8-a2a2a8bf68c0)
persisted_by: milestone-owner (planner ran read-only sandbox; content = planner findings + slice9-corrective-brief.md loci, transcribed + anchor-verified in-tree)
slice: 9 corrective (M-01, worktree m01-listings, branch mis-003/m-01-listings-read-spine)
protocol: slice-8 + hub flow-audit correction (§7 + NEW §15, REVIEW-STANDARD 3d15888). Sequence: test-first impl (Luna high) → **L0 (build+vet) green FIRST** → **[§14 sonnet cold review ∥ L1 unit + integration lane over 0037]** run CONCURRENTLY on the same candidate diff → both green → commit → emit `COMMITTED <sha>` event (hub restart pre-armed, no formal REQUEST). Single dual-gate DELTA covers slice 8+9. No push. GOCACHE absolute. F-01 not reopened.
```

## Planner-resolved findings (drift + omissions vs brief)

1. **Migration is a NEW file, 0036 stays frozen.** `0036_listings.sql:23` keeps `CHECK (status IN ('active','paused','closed','unknown'))` (history). Add **`migrations/0037_listings_status.sql`** widening the constraint at runtime.
2. **`migrations/listings_test.go:36` — DO NOT alter.** It asserts a substring of the **0036** DDL (`normalizedSQL` = 0036 content); 0036 is unchanged so the assertion stays valid. Prove the grown set via the integration lane (applied schema) / a dedicated 0037 check, not by editing the 0036 assertion.
3. **OMITTED LOCUS (brief miss): `apps/server_core/internal/platform/migrate/runner_test.go`** has **two** hard-coded canonical-migration counts = `36` (`:25-26` and `:64-65`). Adding 0037 makes 37 → both `36`→`37` or L1 reds. MANDATORY.
4. **`testdb migrate` auto-discovers 0037** via `//go:embed *.sql` — no registration needed. Integration harness verifies first-apply count 37, second-apply 0.
5. **OpenAPI: exactly 3 status-enum sites** (`:206` filter GET /listings, `:304` filter by-product, `:1955` `ListingReadModel.status` response). `ListingDetail` + grouped responses reuse `ListingReadModel`; summary is counts-only. No other enum.
6. **SDK tests don't assert an exhaustive status set** — growing the `ListingStatus` union (index.ts:216) needs no test-set edit; `tsc --noEmit` must stay 0, `vitest` green.
7. **Integration status fixtures are scenarios, not exhaustive sets** — no fixture expansion; none assert rejection of the 4 new statuses.
8. No requested source anchor drifted numerically (domain :13-18/:162-168, mapper :102-113, query.go :54-60, query_test.go :24 all accurate).

## New canonical values

`under_review`, `inactive`, `payment_required`, `not_yet_active` — added to `active|paused|closed`. **`unknown` STAYS** (ADR-17 honest fallback for genuinely-unrecognized ML statuses).

## Test-first sequence (Luna high)

**Step 1 — RED (behavioral).** Add assertions that fail against current code:
- `connectors/mapper_test.go` (around the `:81-84` table): add rows `{provider:"under_review", want: domain.ListingStatusUnderReview}`, `inactive`, `payment_required`, `not_yet_active`; KEEP a genuinely-unknown row (`{provider:"something-new", want: domain.ListingStatusUnknown}`). Also assert an unrecognized status still yields `unknown` (the WARN path).
- `domain/listing_test.go` (`:13-16` status table): add the 4 `{raw, ListingStatus...}` rows + `IsValid()==true` for each.
- `transport/query_test.go:24`: add the 4 to the valid `status` set (they must parse, not 400).
- (These reference the new consts → package won't compile until Step 2; that compile-gate IS the first red. After Step 2 compiles, mapper_test stays red until Step 3.)

**Step 2 — domain consts.** `domain/listing.go:13-18` add 4 `ListingStatus` consts; `:162-168 IsValid()` add them to the `case`. Package compiles; `query_test.go` + `domain` tests go green; `mapper_test` new rows still RED (mapper returns `unknown`).

**Step 3 — mapper.** `connectors/mapper.go:102-113 canonicalListingStatus`: add 4 explicit `case` arms returning the matching const. In the `default` arm, before `return ListingStatusUnknown`, add `slog.Warn("unmapped provider listing status", "raw", providerStatus, "provider", "mercado_livre")` (package default logger — pure func, no ctx; `import "log/slog"` — stdlib, no cycle). WARN fires ONLY on `default`, never for the mapped 7. `mapper_test` → GREEN.

**Step 4 — migration + counts.**
- NEW `migrations/0037_listings_status.sql`:
  ```sql
  ALTER TABLE listings DROP CONSTRAINT listings_status_check;
  ALTER TABLE listings ADD CONSTRAINT listings_status_check
    CHECK (status IN ('active','paused','closed','unknown','under_review','inactive','payment_required','not_yet_active'));
  ```
  Additive/reversible; existing rows (incl. `unknown`) stay valid; no data rewrite.
- `internal/platform/migrate/runner_test.go` `:26` + `:65`: `36` → `37` (both the `len(want) != 37` guard and the foreign-CWD `len(got) != 37`), and the two `want 36` messages → `want 37`.

**Step 5 — contract (same commit, GOV_API_SDK_SPLIT / D-15 hand-maintained parity).**
- `contracts/api/marketplace-central.openapi.yaml` `:206`, `:304`, `:1955`: extend each `enum: [active, paused, closed, unknown]` → `[active, paused, closed, unknown, under_review, inactive, payment_required, not_yet_active]`.
- `packages/sdk-runtime/src/index.ts:216`: `ListingStatus` union → add `| "under_review" | "inactive" | "payment_required" | "not_yet_active"`. Consumers (`:224`/`:286`/…) auto-covered.

**Step 6 — GREEN gate (L0/L1).** See DoD.

## Definition of Done (pre-commit)

- **L0** (absolute GOCACHE, from `apps/server_core`): `go build ./...` = 0; `go vet ./...` = 0 (whole repo).
- **L1 unit:** `go test -count=1 ./internal/modules/listings/... ./internal/composition/... ./internal/platform/migrate/... ./migrations/...` = 0.
- **L1 integration lane** (migration touched): ephemeral `postgres:16` → `CREATE DATABASE` (retry loop, first-boot 3D000 race) → `go run ./cmd/testdb migrate` applies **37** (second run 0) → `go test -tags=integration -run TestListingsRead -count=1 ./tests/integration` = 0.
- **SDK:** `packages/sdk-runtime` `tsc --noEmit` = 0; `vitest` green.
- **Test-first proof:** mapper/domain/query parse tests RED before Steps 2–3, GREEN after; migration count test RED before Step 4, GREEN after.
- Scope = exactly: `domain/listing.go`, `connectors/mapper.go`, `connectors/mapper_test.go`, `domain/listing_test.go`, `transport/query_test.go`, `migrations/0037_listings_status.sql` (new), `internal/platform/migrate/runner_test.go`, `contracts/api/marketplace-central.openapi.yaml`, `packages/sdk-runtime/src/index.ts` (+ `.mnfs` evidence). NO root.go, NO `query.go` (IsValid delegation covers it), NO 0036 edit, NO F-01 module beyond mapper/domain.
- `docker/dev/*.sh` + worktree `.env` remain uncommitted.
