# F-13 Assisted Linkage API and Runtime Wiring

```yaml
id: F-13
type: feature
status: planned
owner: Feature Implementer
parent: M-06
```

## Outcome

Expose the assisted Mercado Livre→Sankhya linkage workflow through stable
orders HTTP/OpenAPI/SDK contracts and wire it fail-closed to the configured
Oracle read adapter and MPC ledger.

## Brief

Expose the assisted Mercado Livre to Sankhya linkage workflow through stable
orders HTTP, OpenAPI, and SDK contracts, wired fail-closed to the configured
Oracle read adapter and MPC ledger with server-owned audit facts.

## Scope

- Add explicit runtime settings with no field/schema/revision/attestation/limit
  defaults: enabled flag, Oracle schema, TOP 313 header field, configuration
  revision, uniqueness attestation ID, candidate/header-line/lineage limits.
  TOP codes remain fixed 313/306 by product decision.
- Register every new setting in governance. Missing, disabled, malformed, or
  out-of-range configuration leaves the public workflow registered but
  unavailable/fail-closed; server startup and unrelated Oracle reads continue.
- Wire one F-11 reader/service through the F-12 bridge and service to the F-09
  repository. Runtime revision and safe evidence reference are server-owned.
- Add `GET /orders/{provider_order_id}/sankhya-linkage`,
  `GET .../sankhya-linkage/candidates`, and
  `POST .../sankhya-linkage/confirm` with explicit `installation_id`.
- Candidate/current/confirm responses expose generic document/line IDs,
  NUMNOTA/TOP display facts, MPC line IDs, audit provenance/state, and exact
  descendants only. No Oracle/provider payload, buyer PII, credential, or raw
  SQL/field identifier is exposed.
- Confirm accepts selected document/lines, unverified operator actor ID,
  reason, source time, and idempotency key. It never accepts tenant, event ID,
  configuration revision, actor type, evidence reference, or external key.
- Update OpenAPI and `sdk-runtime` together with parity/transport tests. No UI,
  profitability calculation, migration, Oracle/provider write, or live read.

## Acceptance criteria

1. Runtime loader requires enabled=true and every explicit setting; invalid or
   absent config creates a stable unavailable reader/HTTP response without
   startup failure or implicit defaults.
2. Composition uses the same server-derived revision/evidence reference in the
   F-11 reader, F-12 bridge, and F-09 audit. Caller cannot override them.
3. GET current/candidates and POST confirm enforce default tenant plus exact
   installation/order scope, validate JSON/path/time/selection fields, and map
   configuration/unavailable/conflict/not-found errors to stable safe responses.
4. API/SDK schemas preserve nullable unknowns and explicit lineage states;
   event/config/evidence control fields are response-only where safe.
5. Focused Go transport/composition/config tests, OpenAPI/SDK parity tests, SDK
   runtime tests, governance checks, and builds pass. No live side effect runs.

## Expected Output

- Explicit governed runtime configuration with no field, schema, revision,
  attestation, evidence, or limit defaults and stable unavailable behavior.
- Exact-scope GET current/candidates and POST confirm transport with strict safe
  DTO/error mapping and caller audit-control rejection.
- Atomic OpenAPI and `sdk-runtime` contracts with focused fake/unit/build proof.
