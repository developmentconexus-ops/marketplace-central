# M-06 F-02 Correction Brief

## Assigned Failure

The M-06 cold audit rejected F-02 because manual adjustments could be created without operator identity, manual commission was not proven, and append-only history had no real PostgreSQL evidence.

## Required Behavior

1. Require non-empty `actor_type` and `actor_id`; trim all actor fields before persistence.
2. Accept only `order` or `item` scope. Omitted scope may default to `order`.
3. Accept only `freight`, `commission`, `cost`, or `generic_adjustment` category. Omitted category may default to `generic_adjustment`.
4. Item scope requires `provider_item_id`; order scope must not carry item or variation identifiers.
5. Preserve signed amount semantics and never default unknown margin inputs to zero.
6. Generate adjustment IDs independently from the clock so two events at the same timestamp remain distinct.
7. Keep manual adjustments append-only: creating freight and commission events for the same order persists two immutable history rows with distinct IDs and complete actor/reason/time audit.
8. Enforce durable database checks for actor, reason, currency, scope, and category through a new forward migration; do not rewrite applied migration history.
9. Make `actor` required in the OpenAPI create request and in the SDK request type; update SDK tests with the contract.

## TDD Contract

- Add focused failing application/contract tests first and record the RED output.
- Implement the smallest change that makes the tests pass.
- Add PostgreSQL integration coverage gated by `MC_DATABASE_URL`; fake stores do not prove append-only persistence.
- Controller will apply migrations and run the real PostgreSQL gate after handoff.

## Allowed Paths

- `apps/server_core/internal/modules/profitability/**`
- `apps/server_core/migrations/0030_profitability_manual_adjustment_invariants.sql`
- `contracts/api/marketplace-central.openapi.yaml`
- `packages/sdk-runtime/src/index.ts`
- `packages/sdk-runtime/src/index.test.ts`
- `.superpowers/sdd/m06-f02-correction-report.md`

Do not change F-03 calculation semantics, UI, composition, roadmap, or unrelated dirty files. Do not commit in the shared worktree.

## Validation

- Focused Go application tests.
- Profitability module tests.
- SDK runtime tests.
- Real PostgreSQL test after `0028` and `0030` migrations.

