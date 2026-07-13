# F-13 Assisted Linkage API and Runtime Wiring Specification

```yaml
id: F-13
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-13
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-13-assisted-linkage-api-runtime

## Problem

The assisted Mercado Livre to Sankhya linkage domain and application behavior
exist, but they have no public HTTP/OpenAPI/SDK contract or runtime wiring.
Configuration must not invent customer Oracle identifiers, audit provenance,
limits, or revisions, and an invalid deployment must leave the workflow
predictably unavailable without taking down unrelated server behavior.

## Requirements

- Load the enabled flag, Oracle schema, TOP 313 header field, configuration
  revision, uniqueness attestation ID, and candidate/header-line/lineage limits
  explicitly from runtime environment. No required value has a default.
- Register every setting in `contracts/governance/runtime-config.json`.
- Compose the F-11 Oracle reader, F-12 bridge, and F-12 service only from a
  completely valid enabled configuration. Missing, disabled, malformed, or
  out-of-range settings yield a stable unavailable service/HTTP response while
  server startup and unrelated Oracle reads continue.
- Add current-linkage retrieval to the F-12 service when required by transport,
  scoped by server-selected default tenant plus exact installation/order.
- Derive the configuration revision and safe evidence reference once on the
  server and reuse them through reader, bridge, confirmation audit, and public
  response. Requests cannot supply or override either value.
- Expose GET current, GET candidates, and POST confirm endpoints with explicit
  `installation_id`. Confirm accepts only selected document/line mappings,
  unverified operator actor ID, reason, source time, and idempotency key.
- Reject caller-controlled tenant, event ID, configuration revision, evidence
  reference, actor type, external order key, and unknown JSON fields.
- Map invalid input, not found, conflict, and unavailable conditions to stable,
  safe HTTP responses. Preserve nullable unknown facts and explicit lineage
  states without exposing raw Oracle/provider payloads, SQL, configured field
  names, secrets, credentials, or PII.
- Update OpenAPI and `sdk-runtime` atomically and prove request/response parity.

## Acceptance Evidence

- `go-assisted-linkage-http`: focused Go config, composition, service, adapter,
  and transport tests/builds with `GOCACHE=.gocache`.
- `sdk-assisted-linkage`: SDK runtime tests plus API/SDK parity/build checks.
- `governance-contracts`: governance registry/schema checks.
- `git-diff-check`: whitespace and exact allowed-path inspection.

## Non-Goals

- No migration, UI, profitability calculation, Docker change, dependency
  install, live Oracle/Postgres/provider read or write, production
  authentication claim, Oracle mutation, or provider payload exposure.
- No configurable TOP codes; product decisions keep header TOP 313 and
  descendant TOP 306 fixed.
- No inferred candidate selection or caller-generated audit control facts.

## Design

Runtime parsing produces either a complete validated assisted-linkage config or
an unavailable reason; it never substitutes a value. Composition always
registers the routes and injects either the fully wired service or a fail-closed
unavailable implementation. One server-derived evidence reference is a safe,
opaque derivative of the validated deployment attestation/revision and is
reused with the same runtime revision.

Transport uses strict JSON decoding, path/query validation, the composition
root's default tenant, and narrow DTO conversion. Application/domain types stay
generic; only the Oracle adapter sees F-11 configuration. Current linkage reads
the F-09 aggregate through the F-12 service and does not query Oracle.
Candidates and confirm retain F-12 exact-scope and exact-bijection behavior.
OpenAPI schemas are the source contract mirrored by typed SDK runtime methods.

## Edge Cases

- Feature disabled, missing variables, malformed booleans/integers/identifiers,
  blank revision/attestation, and zero/negative/excessive limits are unavailable,
  not startup failures and not defaults.
- Blank or malformed installation/order IDs, invalid JSON, extra fields,
  invalid RFC3339 source time, empty reason/actor/idempotency, invalid document
  IDs, or duplicate/missing/extra line mappings are safe invalid responses.
- Missing scoped order/current linkage is not found; duplicate/idempotency or
  origin conflicts remain conflict; reader/source/config failures are
  unavailable. Internal errors reveal no operational details.
- Nullable NUMNOTA, quantities, amounts, timestamps, and lineage facts remain
  null; `none`, `partial`, `complete`, `conflict`, and unavailable states are
  not collapsed into zero or success.

## Acceptance Criteria

### F13-AC01 Explicit fail-closed runtime

- Criterion: Every runtime value is explicit and governed; invalid or disabled
  configuration keeps only this public workflow unavailable without startup
  failure or implicit defaults.
- Traces to milestone criterion ID: M-06-C02
- Proven by: `go-assisted-linkage-http`, `governance-contracts`.

### F13-AC02 Server-owned audit controls

- Criterion: One server-derived configuration revision/evidence reference is
  used through reader, service, audit, and response; requests cannot supply
  tenant/event/config/evidence/actor-type/external-key controls.
- Traces to milestone criterion ID: M-06-C01
- Proven by: `go-assisted-linkage-http`, `sdk-assisted-linkage`.

### F13-AC03 Safe exact-scope HTTP behavior

- Criterion: Current/candidates/confirm enforce default tenant and exact
  installation/order scope, strict input validation, and stable safe mappings
  for invalid/not-found/conflict/unavailable outcomes.
- Traces to milestone criterion ID: M-06-C01
- Proven by: `go-assisted-linkage-http`.

### F13-AC04 Nullable API/SDK parity

- Criterion: OpenAPI and SDK expose matching narrow DTOs, response-only audit
  controls, nullable unknown facts, and explicit lineage states.
- Traces to milestone criterion ID: M-06-C02
- Proven by: `sdk-assisted-linkage`, `git-diff-check`.

### F13-AC05 Bounded fake proof

- Criterion: Focused fake/unit tests and builds pass with only dispatched paths
  changed and no live side effect.
- Traces to milestone criterion ID: M-06-C01
- Proven by: all four registered commands.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Follow `plan.md`, compile/validate scoped context, implement, and
  record quick-validation evidence.
- Required files/evidence: feature brief, spec, milestone contract, dispatch,
  compiled context, focused command results
- Blockers or open decisions: None.
