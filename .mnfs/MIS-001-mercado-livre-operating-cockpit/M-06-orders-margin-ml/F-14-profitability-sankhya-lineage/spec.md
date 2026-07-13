# F-14 Profitability Sankhya Lineage Specification

## Outcome and ownership

Profitability resolves order tax only from a stable MPC line through the
orders-owned confirmed linkage and the exact current TOP 306 descendants of
that confirmed TOP 313 origin line. Orders owns confirmation and lineage
resolution; profitability owns its order-fact projection, tax-source
selection, component aggregation, provenance, and incomplete margin state.
The internal-read tax port remains the only boundary that reads tax for an
exact Oracle sale-line identity.

No public API, SDK, migration, Oracle query, or provider contract changes.
TOP 313 is confirmation evidence only and is never a tax source.

## Stable line propagation

The profitability-owned order fact carries the canonical opaque `MPCLineID`
for each persisted order item. Only a nonblank ID whose reconciliation state
is `stable` may request current lineage. Ambiguous, legacy-unresolved, blank,
or duplicate stable identities are not tax-resolvable and produce explicit
missing tax without an Oracle tax call.

## Read-only current lineage

The assisted-linkage application boundary exposes a read-only operation scoped
by exact tenant, installation, and provider order identity plus one stable MPC
line ID. It:

- loads the exact persisted Mercado Livre order;
- loads the exact current confirmed F-09 aggregate for the derived external
  order key;
- finds exactly one confirmed origin mapping for the requested MPC line;
- validates the persisted positive TOP 313 origin document/line identity;
- validates the F-11 reader configuration and re-reads descendants for that
  exact origin, including its nullable expected quantity; and
- returns orders-owned generic descendants and one explicit `none`, `partial`,
  `complete`, `conflict`, or `unavailable` state.

The operation never appends or rewrites a linkage event. Missing order,
confirmation, or line mapping remains `none`; malformed, invalid, duplicated,
or otherwise contradictory lineage remains `conflict`; configuration or
source failures remain `unavailable`. Partial descendants remain partial.
No state is promoted to complete by filtering, defaulting, or inference.

## Profitability tax resolution

Profitability depends on a narrow orders-lineage reader port using
profitability-owned request/result values. For each stable MPC line, the
service requests current lineage once and accepts only positive, unique exact
TOP 306 descendant `(NUNOTA, SEQUENCIA)` identities. It never passes the TOP
313 origin to the tax reader and never uses product/date matching.

`none`, `conflict`, `unavailable`, a missing mapping, or any invalid/duplicate
descendant yields missing tax and no Oracle tax call. `partial` may read every
already-known valid exact descendant, but the resulting tax/margin input stays
incomplete even if every returned tax component is known.

For each accepted descendant, profitability invokes the existing exact-line
tax read independently. Provenance is deterministic in ascending document/line
order and retains every exact identity. A component (ICMS, IPI, PIS, or
COFINS) is aggregated only when that component is non-nil for every accepted
descendant. If any descendant lacks a component, that aggregate component is
nil; known values never manufacture unknown siblings as zero. Tax quality can
be complete only when lineage is `complete`, descendants are nonempty and
valid, and every required component is known for every descendant.

## Unknowns and failure behavior

- Missing or non-stable MPC line: missing tax, no lineage/tax read as
  applicable.
- No confirmation or no exact line mapping: `none`, missing tax, no tax read.
- Partial lineage: exact known descendants may be read, but margin stays
  incomplete.
- Conflict, unavailable, invalid, or duplicate descendants: missing tax, no
  tax read.
- Per-descendant missing/partial tax: aggregate only universally known
  components; retain nil for every non-universal component and keep margin
  incomplete.
- Operational errors never become empty success, zero amounts, or inferred
  completion.

## Acceptance criteria

### F14-AC01 Read-only exact current lineage

Focused orders application tests prove exact tenant/installation/order/line
scope, exact persisted mapping use, F-11 descendant re-read, preservation of
all five states, and no ledger append.

### F14-AC02 Stable MPC line is the only resolver

Focused profitability tests prove canonical stable `MPCLineID` propagation and
that blank, legacy, ambiguous, missing, or duplicate identities cannot trigger
an Oracle tax read.

### F14-AC03 Exact TOP 306 one-to-many tax provenance

Focused profitability tests prove one call per unique positive TOP 306
descendant, deterministic multi-line provenance, and no TOP 313 tax source or
product/date fallback.

### F14-AC04 Honest aggregation and partial state

Focused profitability tests prove per-component all-descendant aggregation,
nil rather than zero for any unknown component, partial-lineage exposure of
known exact amounts, and incomplete margin quality until lineage and tax facts
are complete.

### F14-AC05 Runtime wiring remains conditional and bounded

Focused composition/compile proof establishes that profitability receives the
lineage boundary only when assisted-linkage runtime dependencies exist, with
no HTTP/OpenAPI/SDK or runtime configuration change.
