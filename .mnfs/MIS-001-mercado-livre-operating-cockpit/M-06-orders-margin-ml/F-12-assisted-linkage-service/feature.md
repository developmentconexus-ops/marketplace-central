# F-12 Assisted Sankhya Linkage Service

```yaml
id: F-12
type: feature
status: planned
owner: Feature Implementer
parent: M-06
```

## Outcome

Implement one orders-owned assisted linkage service that lists exact TOP 313
candidates for an imported Mercado Livre order, validates an operator's
explicit all-line selection, appends the F-09 mapping idempotently, and follows
each confirmed origin through F-11 to exact TOP 306 descendants.

## Brief

Add the bounded orders-owned application seam that converts an operator's
explicit exact TOP 313 selection into one F-09 append-only mapping and then
reads exact per-line TOP 306 descendants through F-11 without leaking Oracle
types or claiming production authentication.

## Expected Output

- Exact-key candidate listing for one persisted Mercado Livre order with no
  automatic selection or persistence.
- Fail-closed confirmation that validates one selected candidate and a full
  stable-MPC-line-to-candidate-line bijection before an exact-evidence append.
- Idempotent mapping return plus explicit per-line descendant states, including
  exact conflict versus unavailable, without undoing or manufacturing the
  confirmed mapping.
- Focused fake/unit evidence and one intentional Feature commit limited to the
  dispatched orders service, ports, bridge, order lookup, and Feature artifacts.

## Scope

- Add generic orders-domain candidate, confirmation, and lineage result types.
- Add an orders-owned reader port and an adapter to the separate internal-read
  Sankhya linkage service. Oracle DTOs/types must not cross the adapter.
- Add an exact order lookup port/method on the existing Postgres order
  repository; tenant and installation remain explicit/scoped.
- `ListCandidates` derives the exact external key
  `ml:v1:<installation>:<provider-order>`, requires a persisted Mercado Livre
  order, and returns bounded header/line candidates. It never auto-selects.
- `Confirm` requires configuration validation, exact candidate document ID,
  explicit mapping for every persisted stable MPC line exactly once, exact
  candidate line existence, actor type/ID, reason, source time, idempotency
  key, runtime-derived config revision, and a server-generated event ID.
  Ambiguous/legacy lines, missing/extra mappings, non-candidates, empty runtime
  revision, event-ID generation failure, and conflicts fail closed.
- Append the F-09 aggregate with `evidence_state=exact`. Actor provenance is
  explicitly `operator_supplied_unverified`; this Feature does not claim or fix
  production authentication/manual-adjustment authorization.
- After a confirmed/idempotent mapping, read every exact descendant set. None,
  partial, complete, conflict, and unavailable remain explicit; no product/date
  fallback and no tax default is created.

## Acceptance criteria

1. Candidate listing uses only the derived exact external key and does not
   persist or promote a candidate to proof.
2. Confirmation accepts only an explicitly selected candidate and a bijection
   covering every stable MPC order line and exact candidate line; all invalid,
   ambiguous, legacy, missing, extra, duplicate, or conflicting cases fail
   before ledger append.
3. Semantically identical retry returns the same append-only mapping; a reused
   key or origin mismatch returns conflict with no overwrite.
4. Exact TOP 306 descendants are returned per MPC line after confirmation;
   unknown quantity/lineage remains partial/unknown and Oracle unavailability
   never changes the confirmed mapping or manufactures identity.
5. Focused orders unit/adapter tests pass. No HTTP/OpenAPI/SDK, composition,
   runtime registry, migration, live Oracle/provider/Postgres, or UI change.
