# F-02 Slice 2 — BLOCKED

## Classification

Contract/runtime conflict: exact cost-dependent exception filtering cannot be combined with the required PostgreSQL `limit+1` keyset query under the ratified cross-database read design.

## Evidence

- `plan.md:169-177` requires every exception/`has_exception` predicate to be exact, requires known-formula filtering, and requires filtering before truncation with a PostgreSQL `limit+1` keyset query.
- The Slice 2 prompt likewise requires `exception`/`has_exception` predicates to be exact in the list repository SQL and requires the joined projection/filtering before truncation.
- D-17 defines `below_margin` as an active exception participating in both `exception` and `has_exception`.
- D-20 ratifies latest cost as one bounded Oracle `GetCostFactsByIDs` read and expressly rejects a local persisted cost projection.
- D-22 computes `below_margin_worst_case` from listing price plus Oracle cost and Oracle ICMS ceiling plus PostgreSQL marketplace policy. It forbids persisting or defaulting those inputs.
- The `listings` and product-link PostgreSQL tables contain no cost, ceiling, policy-derived below-margin value, or queryable equivalent.

Therefore PostgreSQL cannot evaluate `filter.exception=below_margin`, nor can it evaluate `filter.has_exception=true|false` exactly because below-margin participates in the D-17 active-issue set. Applying these predicates after fetching `limit+1` rows can return short/empty pages despite later matches and can emit a cursor that skips qualifying rows. Fetching an unbounded PostgreSQL candidate set and filtering in application code would violate the required `limit+1` query/bounded-read design and the 2,000-row keyset criterion.

## Smallest unblock decision

Ratify one of these contract changes before implementation:

1. Defer `exception=below_margin` and `has_exception` list filtering until an architecture-owned, tenant-scoped PostgreSQL projection makes all inputs queryable; or
2. Amend list pagination to a bounded iterative keyset scan (multiple PostgreSQL pages and bounded Oracle batches until `limit+1` qualifying results), explicitly replacing the one-query/one-batch requirement; or
3. Authorize a local tenant-scoped projection for cost/rate/policy evaluation, reversing the D-20 no-projection ruling.

No source files remain changed by this attempted slice. No dependency, composition, OpenAPI, SDK, migration, or cross-module internal file was modified.
