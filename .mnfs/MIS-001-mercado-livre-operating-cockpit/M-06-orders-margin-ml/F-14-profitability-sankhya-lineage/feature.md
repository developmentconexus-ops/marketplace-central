# F-14 Profitability Sankhya Lineage

```yaml
id: F-14
type: feature
status: briefed
owner: Feature Implementer
parent: M-06
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-1
lifecycle_scope: feature
```

## Brief

Make profitability consume one confirmed stable MPC order-line linkage and its
exact existing one-to-many TOP 306 descendants before reading tax, preserving
partial and unknown facts without ever using TOP 313 or a zero default.

## Objective

Make profitability consume the confirmed MPC-owned order-line linkage and the
exact existing one-to-many TOP 306 descendants before reading tax. Remove the
known empty-identity gap without introducing product/date heuristics, Oracle
writes, or zero defaults.

## Required behavior

- Carry stable `MPCLineID` from canonical orders into profitability-owned order
  facts; ambiguous or legacy-reconciliation lines are never tax-resolvable.
- Extend the canonical assisted-linkage application boundary with a read-only
  current-lineage operation. It must load the exact tenant/installation/order
  confirmation, re-read descendants through the validated F-11 reader, and
  preserve `none`, `partial`, `complete`, `conflict`, and `unavailable`.
- Profitability resolves tax sources only by exact stable MPC line. A missing
  confirmation, missing line mapping, conflict, unavailable lineage, or no
  descendants yields missing tax and no Oracle tax call.
- Every tax read targets one exact TOP 306 descendant `(NUNOTA, SEQUENCIA)`.
  The TOP 313 origin is never used as a tax source.
- For one-to-many descendants, aggregate a component only when that component
  is known for every exact descendant. Preserve deterministic multi-line
  provenance. Duplicate or invalid descendants fail closed.
- Partial lineage may expose already-known exact tax amounts, but it must keep
  the resulting margin input/snapshot incomplete. Unknown operational values
  remain nil, never zero.
- Wire the boundary in the composition root only when the assisted-linkage
  runtime is available. Existing HTTP/OpenAPI/SDK contracts do not change.

## Expected Output

- A read-only exact current-lineage operation preserves `none`, `partial`,
  `complete`, `conflict`, and `unavailable` without a ledger write.
- Profitability carries stable `MPCLineID`, reads tax only for exact unique TOP
  306 descendants, and retains deterministic one-to-many provenance.
- Tax components aggregate only when known for every exact descendant;
  partial or unknown lineage and tax remain incomplete and nil rather than
  zero.
- Conditional composition wiring, focused fake tests, repository-wide Go
  proof, scoped diff evidence, and one intentional Feature commit establish
  the bounded behavior.

## Acceptance evidence

- Focused orders application tests prove read-only current-lineage resolution,
  exact scope, state preservation, and no ledger write.
- Focused profitability tests prove exact TOP 306 calls, one-to-many component
  aggregation, partial-lineage incompleteness, and fail-closed missing states.
- Focused adapter/root tests or compile proof establish boundary wiring.
- Repository-wide Go tests and `git diff --check` pass with `GOCACHE=.gocache`.

## Non-goals and stop conditions

- Do not change migrations, public APIs/SDK, UI, Docker, research, or runtime
  configuration.
- Do not invent Oracle linkage, TOP behavior, quantities, tax amounts, or
  authentication claims.
- Do not read or write live Oracle/Postgres/provider data.
- Stop for an architecture conflict, a need to weaken exact-line selection, or
  any requirement to treat a partial/missing fact as complete.

## Handoff

Write `spec.md`, `plan.md`, scoped `context.json`, and `validation.md`; create
one intentional commit and return the compact feature checkpoint to Milestone.
