# F-12 Assisted Sankhya Linkage Service Specification

## Outcome and ownership

The orders module owns a bounded assisted-linkage application service. It
requires one exact persisted Mercado Livre order, derives only
`ml:v1:<installation_id>:<provider_order_id>`, lists generic exact TOP 313
candidates through an orders port, and never selects or persists a candidate.
Confirmation re-reads the exact candidate set, validates the operator's
explicit document and all-line selection, appends the F-09 aggregate, and then
reads exact TOP 306 descendants for every persisted mapping line.

Orders domain/application/ports remain generic. A narrow `internalread` adapter
is the only code that sees F-11 types. The existing Postgres order repository
implements an exact tenant/installation/order lookup. There is no transport,
composition, runtime registry, migration, public contract, or live-system work.

## Exact order and candidate listing

The service requires nonblank tenant, installation, and provider order IDs.
The order lookup is scoped by all three values and must return an already
persisted order whose provider code is exactly `mercado_livre` and whose stored
installation/order identity equals the request. Absence, scope mismatch, or a
different provider fails closed.

Before listing, the service validates the dedicated reader configuration. It
derives the external order key from the validated scope and passes only that
exact key to candidate lookup. Candidate results carry generic document IDs,
the proved header operation code, nullable header/line consistency evidence,
and exact generic line identities. Listing has no repository write and never
promotes filtering, ordering, or a single result into identity proof.

## Confirmation contract

Confirmation requires the exact scoped order, a positive explicitly selected
candidate document ID, one supplied mapping per order line, actor ID, reason,
source time, and idempotency key. Configuration revision is runtime-derived
from the validated reader bridge and is never accepted from the caller. Event
ID is generated inside the service through an injectable generator whose
default uses cryptographic randomness; it is never caller-controlled. Empty
runtime revision or event-ID generation failure returns before append. Actor
provenance is fixed to `operator_supplied_unverified`; it records supplied
intent but does not assert authentication or manual-adjustment authorization.

The service validates reader configuration and re-runs candidate lookup using
the exact derived key. Exactly one returned candidate must match the selected
document ID, it must prove header operation code TOP 313, and every candidate
line must be a positive unique identity on that selected document. Every
persisted order line must have a valid opaque MPC line ID and reconciliation
state `stable`; legacy-unresolved or ambiguous rows fail closed.

The supplied mapping must be a bijection:

- every persisted stable MPC line appears exactly once;
- no unknown, missing, extra, or duplicate MPC line appears;
- every exact line of the selected candidate appears exactly once;
- no unknown, missing, extra, duplicate, or different-document candidate line
  appears; and
- the order-line and candidate-line cardinalities are equal.

All checks happen before the ledger port is called. The aggregate uses the
derived key, selected header, every validated line mapping, event type
`confirmed`, evidence state `exact`, the fixed actor type, and supplied audit
facts. Configuration revision and event identity come only from the runtime
bridge and service respectively; recorded time is an application clock fact.
F-09 owns transactional append/idempotency/conflict behavior: a semantically
identical retry returns the persisted aggregate; reused keys or origin
conflicts fail without overwrite.

## Descendant behavior

Only after `AppendConfirmation` succeeds does the service read descendants.
It iterates every exact line in the persisted/idempotently returned aggregate,
not an unpersisted request copy. Each request binds that exact internal origin
document/line and its nullable expected quantity when candidate evidence can be
matched safely. Results are converted to orders-owned generic descendant types
with explicit `none`, `partial`, `complete`, or `conflict` state.

An exact F-11 lineage conflict, unexpected error, or malformed returned lineage
is preserved as per-line `conflict`. Only configuration/source unavailability
is returned as `unavailable`; neither state carries fabricated descendants.
The service continues attempting the remaining persisted lines. Such a failure
never rolls back, deletes, rewrites, or fabricates the already confirmed
mapping, and no product/date fallback or tax identity is created.

## Errors and unknowns

Invalid request/scope/order/provider/candidate/bijection/audit data returns a
stable orders-domain invalid or not-found error before append. F-09 conflict
and not-found errors remain explicit. Nullable document number, date, product,
quantity, amount, and attended quantity remain nil. Unknown operational values
never become zero/default facts.

## Acceptance criteria

### F12-AC01 Exact listing without proof promotion

Focused service tests prove exact scoped persisted-order lookup, exact derived
key use, configuration validation, generic candidate conversion, and no append
during listing.

### F12-AC02 Full exact confirmation bijection

Domain/service tests reject non-candidates, non-TOP-313 candidates,
different-document lines, missing/extra/duplicate mappings, unknown MPC lines,
missing/extra candidate lines, invalid identities, and every ambiguous or
legacy order line before ledger append.

### F12-AC03 Exact audited idempotent append

Tests prove valid confirmation appends all lines with `evidence_state=exact`,
fixed actor type `operator_supplied_unverified`, runtime-derived configuration
revision, server-generated event ID, supplied actor/intent facts, and the exact
derived key. Empty revision and event-ID generation failure stop before append;
an identical retry returns the repository aggregate and conflicting append
errors are not converted into success or overwrite.

### F12-AC04 Per-line descendants preserve mapping and unknowns

Tests prove every persisted origin is read after append, exact per-line
none/partial/complete/conflict states and nullable facts are preserved,
lineage-conflict/malformed responses remain `conflict`, only configuration or
source failures become `unavailable`, and the confirmed mapping plus remaining
exact reads remain intact.

### F12-AC05 Bounded seam and fake proof

Focused Go fake/unit tests and scoped diff inspection prove only the dispatched
orders domain/ports/application/internal-read bridge/order lookup and Feature
artifacts changed, with no HTTP/OpenAPI/SDK, composition, runtime, migration,
live Oracle/provider/Postgres, dependency, secret, or PII operation.
