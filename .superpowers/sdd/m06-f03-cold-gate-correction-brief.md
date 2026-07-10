# M-06 F-03 cold-gate correction brief

## Context

An independent cold review found three blocking correctness/safety defects in profitability. Work directly on the user-authorized `main` checkout; it is heavily dirty. Do not reset, revert, stage, or commit unrelated changes.

## Binding constraints

- Unknown cost/freight/fee/tax must never become zero or realized margin.
- Order-scope cost adjustment is not item cost; it is an order-level adjustment total.
- Every manual-adjustment write must be idempotent or explicitly duplicate-safe.
- Preserve append-only audit behavior and truthful actor/reason/timestamp.
- Any API change updates OpenAPI and `packages/sdk-runtime` together.
- Use `GOCACHE=.gocache` for host Go commands; Docker is available for real PostgreSQL verification.
- No provider writes, no WSL, no secrets.

## Root causes confirmed

1. `application/service.go` adds `nil` tax components as zero with `addOptional`; only an entirely nil aggregate gets `missing_tax`.
2. `applyAdjustment` routes every `cost` category to `CostAmount`, including order scope.
3. Create-adjustment assigns a random ID and bare INSERT; retrying a request produces a second append-only row.

## Required behavior

### A. Partial tax honesty

- An item snapshot with one or more unknown tax components must be `incomplete`, include `missing_tax`, and keep contribution/margin nil even if another tax component has an amount.
- The aggregate order snapshot must propagate this incompleteness from its item(s).
- Known tax components may still remain visible in `TaxAmount`; they must not authorize realized math while the tax set is incomplete.

### B. Cost-adjustment scope

- An item-scope `cost` adjustment changes `CostAmount`.
- An order-scope `cost` adjustment changes `AdjustmentAmount`, not `CostAmount`.
- Keep existing freight/commission behavior unchanged.

### C. Manual-adjustment idempotency

- Add an explicit `idempotency_key` to the create-manual-adjustment contract/domain/storage path.
- Require a non-empty key for new requests; use a stable structured error code for absence.
- Scope uniqueness to tenant + installation + idempotency key in Postgres using a forward-safe migration; existing audit rows remain readable.
- Retrying the same key must return the original adjustment and create exactly one persisted row. It must not mutate the original audit record.
- Update transport, OpenAPI, SDK, application/store interfaces, Postgres store, tests and any necessary UI caller/identity flow to supply a stable key for a single user submission.

## TDD evidence required

Create focused tests before production changes and record RED output, then GREEN output for:

1. Partial known/unknown tax at item and aggregate order scopes.
2. Order-scope versus item-scope cost adjustment mapping.
3. Duplicate manual-adjustment key returning one immutable persisted adjustment.
4. Missing idempotency key rejection.

Run focused Go tests, relevant transport/Postgres integration tests (against Docker PostgreSQL if practical), SDK tests, feature-orders tests, gofmt, and a profitability-only boundary `rg` proving `orders` imports remain limited to `profitability/adapters/orders`.

## Report

Write results, RED/GREEN commands, changed paths, tests and concerns to:
`.superpowers/sdd/m06-f03-cold-gate-correction-report.md`
