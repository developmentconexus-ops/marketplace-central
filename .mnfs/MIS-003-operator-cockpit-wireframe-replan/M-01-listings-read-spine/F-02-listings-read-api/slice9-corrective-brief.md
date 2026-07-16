# F-02 Slice 9 — corrective brief: grow canonical listing-status enum (C10 mapping gap)

```yaml
slice: 9 (corrective, M-01, same chip/worktree)
origin: C10 re-drive found status unknown 7/34 = 20.6% > `<20%` (validation-contract.md:175 blocking "adapter mapping gap")
hub_ruling: ESCALATION C10 → R2 GROW ENUM FULL + R3 M-01 corrective slice (HUB-EVENT-ESCALATION-c10-status-mapping-gap.md; hub adjudication 2026-07-16, operator-ratified)
protocol: slice-8 flow — plan (Sol medium) → impl (Luna high, test-first) → review §14 (sonnet cold) → L0 (build+vet) / L1 (listings+connectors+composition unit + integration lane, migration touched) → commit. F-01 NOT reopened.
seams_granted_by_hub: migration number 0037 (pre-allocated; main@0035, worktree@0036, single track); OpenAPI listings-status section (enum + filter) locked to this slice.
no_push. GOCACHE absolute. no self-boot / no .env in session.
```

## Ruling (R2) — verbatim intent

Canonical listing status GROWS to add **`under_review`, `inactive`, `payment_required`, `not_yet_active`**
(the remaining documented Mercado Livre `item.status` values). **`unknown` REMAINS** as the honest ADR-17
fallback for genuinely-unrecognized provider statuses (never delete it; `under_review` ≠ `paused` — an
operator must not treat an in-review item as reactivatable). The mapper must:
- map every documented ML status **explicitly** (no silent collapse), and
- on the `default → unknown` branch, emit a **`slog` WARN carrying the raw provider status string** so a
  future mapping gap is never blind again (was: silent 20.6% flood).

FE/cockpit rendering of the new statuses = a downstream UI milestone, **out of M-01**.

## Documented ML `item.status` set (firm datum, hub R1)

`active | paused | closed | under_review | inactive | payment_required | not_yet_active`
(chip cannot read the raw values of the 7 unknown rows — no DB/.env/self-boot; DB stores canonical only,
provider payload dies at the adapter by doctrine. Exact per-status counts become R1-closing evidence
**after** this slice via a mapped re-sync.)

## Ripple loci (ALL worktree paths — one commit, doctrine §7 GOV_API_SDK_SPLIT)

| # | File | Locus | Change |
|---|------|-------|--------|
| 1 | `apps/server_core/internal/modules/listings/domain/listing.go` | `:13-18` consts | add 4 `ListingStatus` consts (`under_review`/`inactive`/`payment_required`/`not_yet_active`) |
| 1b | same | `:162-168` `IsValid()` | add the 4 to the valid `case` |
| 2 | `apps/server_core/internal/modules/listings/adapters/connectors/mapper.go` | `:102-113` `canonicalListingStatus` | add 4 explicit `case` arms; in `default` branch `slog.Warn` the **raw** `providerStatus` before returning `unknown` |
| 3 | `apps/server_core/migrations/0037_<name>.sql` | NEW | `ALTER TABLE listings DROP CONSTRAINT listings_status_check`, re-add `CHECK (status IN ('active','paused','closed','unknown','under_review','inactive','payment_required','not_yet_active'))` (match the exact constraint name/shape from `0036_listings.sql:23`) |
| 3b | `apps/server_core/migrations/listings_test.go` | `:36` | update the asserted constraint string to the grown set |
| 4 | `apps/server_core/internal/modules/listings/transport/query.go` | `:54-60` | NO code change — validation delegates to `domain.IsValid()`, auto-covered. (Confirm, don't edit.) |
| 4b | `apps/server_core/internal/modules/listings/transport/query_test.go` | `:24` valid-status set | add the 4 new values so the parse test proves they pass; `banana`-style reject case stays |
| 5 | `contracts/api/marketplace-central.openapi.yaml` | `:206`, `:304`, `:1955` | add the 4 to each `enum: [active, paused, closed, unknown]` (2 filter params + `ListingReadModel.status` response) |
| 6 | `packages/sdk-runtime/src/index.ts` | `:216` `ListingStatus` union | add the 4 to the string union (consumers `:224/:286/…` auto-covered); update `index.test.ts` if it asserts the union set |
| 7 | tests | `domain/listing_test.go` status table (`:13-16`), `connectors/mapper_test.go` (`:81-84`) | add mapped-case assertions for the 4 new statuses **and** keep one genuinely-unknown case (e.g. `something-new` → `unknown`, asserting the WARN path is reached) |

## Constraints / guard-rails

- **Do not** delete or repurpose `unknown` — honest fallback stays (ADR-17).
- **Do not** touch F-01 module beyond `mapper.go` + `domain/listing.go` (both already M-01-scoped seams). No new module, no root.go change (composition unaffected — enum grow is transparent to wiring).
- **slog nuance:** `canonicalListingStatus` is a pure func with no logger/ctx param. Prefer `slog.Warn` on the package default logger (no ctx) with structured `raw=providerStatus`; WARN fires **only** on the `default` arm, never for the newly-mapped statuses. Planner: confirm no import cycle (`log/slog` stdlib, safe).
- **Migration is additive + reversible** (constraint widening only; no data rewrite — existing `unknown` rows stay valid and will re-map on next refresh). L1 MUST run the integration lane (migration touched): ephemeral Postgres → `testdb migrate` (now applies 0037) → listings integration tests green.
- **OpenAPI+SDK same commit** (GOV_API_SDK_SPLIT, D-15 contract-first hand-maintained parity — no codegen claim).
- Slice size will exceed slice-8's ~355 lines only modestly (mostly enum lists + test rows); acceptable per hub (mechanical fan-out, single logical change). Note it in review, no artificial split.

## Definition of done (pre-commit)

- L0: `go build ./...` 0, `go vet ./...` 0 (whole repo, absolute GOCACHE).
- L1: `go test -count=1 ./internal/modules/listings/... ./internal/composition/...` 0; integration lane (ephemeral PG over 0037) listings tests 0; SDK `tsc --noEmit` 0 + `vitest` green.
- Test-first evidence: mapper + domain + migration + transport parse tests fail BEFORE impl (unknown for the 4), pass AFTER.
- §14 sonnet cold review APPROVE (no blocking/unresolved-important).

## After commit (own task, not this slice)

REQUEST hub re-sync the installation → C10 re-drive: expect unknown ≪ 20% (likely 0%); record per-status
counts = R1-closing evidence. Then dual-gate DELTA from `e2cde36` covering **slice 8 + 9 together** (COLD
OPUS SUBAGENT `model=opus` + Sol `--effort medium`, simultaneous, merge §8, reconciliation table in CLOSED)
→ P8 CLOSED.
