# F-11 Sankhya Linkage Read Contract

```yaml
id: F-11
type: feature
status: planned
owner: Feature Implementer
parent: M-06
```

## Brief

Add a separate MPC-owned `internal_read` contract and Oracle adapter for the
assisted linkage flow: validate an explicit configured TOP 313 header field,
read bounded exact candidates and lines, and follow exact `TGFVAR` origins to
TOP 306 descendants with attended quantity.

## Validation Expectations

- Invalid configuration, metadata, uniqueness, or Oracle availability fails
  closed before candidate or lineage SQL executes.
- Query-builder tests prove safe quoted identifiers, bound values, bounded
  candidates/lines/lineage, exact TOP 313/306 checks, and exact `TGFVAR`
  origins.
- Generic models preserve nullable evidence and every one-to-many descendant
  with explicit none/partial/complete/conflict lineage states.
- Focused fake Go tests pass without live Oracle, provider, Postgres,
  dependency, composition, or public-contract changes.

## Scope

- Add generic internal-document/linkage domain models and a separate reader
  port/service. Do not extend the broad existing `Reader` interface.
- Add an Oracle adapter configured explicitly with schema, header field name,
  configuration revision, expected origin TOP 313, destination TOP 306, and a
  uniqueness-attestation identifier. There is no default field name.
- Validate the identifier by a strict Oracle-safe allowlist and exact metadata
  lookup. Values always use binds. Missing metadata, incompatible field type or
  capacity, duplicate nonblank values, missing attestation, or unavailable DB
  fails closed.
- Candidate query matches the exact nonblank digits-only external key and TOP 313,
  is bounded, and returns generic document ID/number/TOP plus exact item lines.
  Product/date/value are returned only as consistency evidence, never proof.
- Descendant query uses exact origin document/line binds through `TGFVAR`,
  requires destination TOP 306, returns every exact destination document/line
  and nullable attended quantity, and preserves none/partial/complete/conflict
  as explicit typed states.

## Acceptance criteria

1. Invalid/missing field config or metadata/uniqueness validation returns a
   stable fail-closed error and no candidate/lineage query runs.
2. Query-builder tests prove the configured identifier is allowlisted/quoted,
   all values are binds, candidate limits are bounded, and TOP 313/306 plus
   exact `NUNOTA`/`SEQUENCIA`/`TGFVAR` predicates are present.
3. Candidate and descendant models expose only generic document/line identities,
   safe evidence refs/states, and nullable operational facts; no Oracle DTO or
   raw row/payload crosses the adapter.
4. Missing or one-to-many descendants remain explicit; no product/date match or
   zero/default tax identity is produced.
5. Focused internal-read Go tests pass. No live Oracle, provider, Postgres,
   OpenAPI/SDK, composition, dependency, or runtime-registry change occurs.
