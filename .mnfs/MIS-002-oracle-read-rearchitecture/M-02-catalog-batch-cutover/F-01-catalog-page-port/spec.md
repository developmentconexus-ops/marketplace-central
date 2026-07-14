# F-01 catalog page port — execution spec

## Scope

Implement the IC-01 catalog page read seam in `internal_read`, backed by one
Oracle query per list/search call, then route the catalog adapter's listing
projection through that seam. The M-01-C04 `FULL_SCAN` verdict requires the
query to join base TGF* tables directly.

## Contract

The port exposes list and bounded search operations returning `CatalogFactPage`
with nullable operational facts, quality flags, and `as_of`. `Cursor` is the
opaque base64 encoding of the positive last `CODPROD` keyset value. Malformed,
non-numeric, or non-positive cursor values return the typed invalid-cursor
error before any Oracle call.

## Acceptance

- Interface/types are stable in `internal_read/ports` and were committed first.
- List uses one set-based base-table JOIN with stock aggregation, price, cost,
  keyset predicate, and `limit+1` peek.
- Search uses one bounded query, sorts by internal product id, and has no next
  cursor.
- Missing stock/price/cost remain nil and carry their quality flags; duplicate
  active price rows carry `ambiguous_price` with a nil amount.
- Catalog listing no longer composes one product lookup plus three fact reads.
- Required build, targeted tests, and full test evidence are recorded in
  `validation.md` without claiming unrun commands.
