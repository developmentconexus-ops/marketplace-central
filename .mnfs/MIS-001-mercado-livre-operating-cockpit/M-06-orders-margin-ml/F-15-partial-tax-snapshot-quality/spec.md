# F-15 Partial Tax Snapshot Quality Specification

## Outcome and ownership

Profitability snapshot aggregation retains every known non-nil exact-line tax
amount while treating any tax component whose input quality is not `complete`
as missing for margin-completeness purposes. Profitability application owns
this correction. Existing orders lineage, internal-read tax, repository,
transport, API, SDK, migration, UI, and runtime boundaries do not change.

## Snapshot quality contract

For each ICMS, IPI, PIS, and COFINS input, snapshot input application records
the amount when it is known. Independently, it records missing tax whenever
the amount is nil or the input quality is not `complete`.

Consequently:

- partial lineage with all four component amounts known retains those amounts
  but produces incomplete item and order snapshots;
- complete lineage with any exact-line source component quality that is not
  complete has the same incomplete result while retaining known amounts;
- any incomplete item propagates missing tax to the order rollup; and
- contribution and margin remain nil whenever the tax completeness condition
  is not satisfied.

Fully complete lineage and fully complete exact-line tax inputs preserve the
existing complete contribution and margin behavior. Unknown amounts remain
nil and are never replaced with zero.

## Interfaces, consumers, and legacy decision

No port or interface changes are required. `ImportOrder` remains the consumer
that resolves profitability inputs and builds item/order snapshots. Existing
snapshot DTOs, persistence contracts, and callers retain their shapes.

The legacy amount-only missing-tax test in snapshot aggregation is corrected
in place: amount presence and input quality are separate facts. No fallback,
estimation, compatibility path, or manual adjustment behavior is introduced.

## Acceptance criteria

### F15-AC01 Partial-lineage known amounts stay incomplete

An import-to-snapshot regression with partial lineage and all four known tax
components proves known amounts are retained, item and order quality include
missing tax, and both item/order contribution and margin remain nil.

### F15-AC02 Complete lineage does not promote incomplete source quality

An import-to-snapshot regression with complete lineage and all four known tax
amounts but at least one incomplete exact-line source quality proves the same
incomplete/missing-tax and unrealized-margin behavior.

### F15-AC03 Fully complete control remains complete

An import-to-snapshot regression with complete lineage and complete qualities
for all known tax inputs proves the existing complete contribution and margin
behavior remains intact.

### F15-AC04 Evidence is corrected append-only

F-14 validation is amended to identify its overclaim and point to the F-15
regression proof. The M-06 correction ledger appends attempt 2 without
rewriting prior rows and continues to state that C03 is deferred/failing.
