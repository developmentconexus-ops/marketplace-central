# F-11 Sankhya Linkage Read Contract Specification

## Outcome and ownership

`internal_read` owns a separate typed Sankhya-linkage read boundary. It validates
one deployment-supplied `TGFCAB` header-field identifier, reads bounded exact TOP
313 candidates and their lines, and follows exact `TGFVAR` origins to every TOP
306 descendant. The existing broad `ports.Reader` and Oracle `Reader` remain
unchanged. Future orders code consumes this feature only through the new
`SankhyaLinkageReader` port/application service; this feature adds no composition
or runtime activation.

## Configuration and activation

The adapter receives an explicit schema, header field name, configuration
revision, expected origin TOP, expected destination TOP, uniqueness-attestation
identifier, candidate limit, candidate-line limit, and lineage limit. Every
limit is positive with a fixed maximum. There is no default field name, schema,
TOP, attestation, or limit. Identifiers must match
`^[A-Z][A-Z0-9_$#]{0,29}$` after requiring
the configured value already be uppercase and unchanged. SQL quotes only an
identifier that passed this allowlist; all request and configuration values are
binds.

Before any candidate or lineage read, validation must prove:

- the database is available;
- the configuration is complete, uses origin TOP 313 and destination TOP 306,
  carries a nonblank revision/attestation, and has positive bounded header,
  line, and lineage limits;
- an exact `ALL_TAB_COLUMNS` lookup for the configured owner, `TGFCAB`, and field
  returns exactly one compatible text column with at least 160 character
  capacity; and
- a bounded aggregate probe finds no duplicate nonblank field value.

Missing or incompatible metadata, absent attestation, duplicate values, invalid
identifier/configuration, or unavailable Oracle returns a stable fail-closed
typed error. No candidate or descendant query runs after such a failure. The
attestation is an explicit administrative prerequisite, not a claim that MPC can
prove the external uniqueness mechanism from catalog facts.

## Separate port and service

The new port exposes configuration validation, exact bounded candidate lookup,
and exact descendant lookup. Inputs carry only the external key or a generic
origin document/line identity. Outputs carry generic document IDs/numbers, the
proved header operation code, exact item-line identities, nullable consistency
evidence, descendants, and explicit lineage state. Oracle rows, driver values,
table names, and SQL never cross the adapter.

The application service exposes validation and delegates only through the
dedicated reader. The Oracle implementation validates before each data read,
keeping activation fail-closed even if a consumer skips a separate startup
validation call.

## Candidate semantics

Candidate lookup requires a nonblank exact external key, binds that value and
expected TOP 313, orders deterministically, and fetches at most the configured
limit plus one. More than the configured limit fails as ambiguous instead of
silently truncating. For every retained header, a second bound query reads exact
`TGFITE` lines by `NUNOTA`, fetches the candidate-line limit plus one, and
returns positive `SEQUENCIA` identities. Line overflow fails as candidate
ambiguous.

The deployed Sankhya effective-TOP history join is not proved by repository
truth. This contract therefore returns and checks the proved `TGFCAB.CODTIPOPER`
header fact. It does not invent a `DHTIPOPER`/effective-TOP join or claim temporal
TOP semantics. Product, date, quantity, and value may be nullable consistency
evidence only and never identity proof.

## Descendant and state semantics

Descendant lookup binds exact origin `NUNOTA` and `SEQUENCIA`, joins `TGFVAR` to
destination `TGFCAB`/`TGFITE`, requires destination header operation code TOP
306, and returns every destination `NUNOTA`/`SEQUENCIA` plus nullable
`QTDATENDIDA`. It fetches the lineage limit plus one and fails as a lineage
conflict on overflow. Rows are not collapsed, so bounded one-to-many lineage is
generic and preserved.

Lineage states are explicit:

- `none`: no exact descendants;
- `partial`: descendants exist and expected origin quantity is unknown, at
  least one attended quantity is missing, or their known attended total is
  below the expected origin quantity;
- `complete`: descendants exist, all attended quantities are known, and no
  optional expected-origin-quantity mismatch is present;
- `conflict`: any destination identity is duplicated, a destination identity is
  invalid, the lineage bound is exceeded, or known attended quantity exceeds
  an expected origin quantity.

Unknown expected quantity remains nil; it never becomes zero. No descendant is
converted into tax identity by this feature.

## Stable errors and verification

The feature defines stable error classes for configuration invalid, metadata
mismatch, uniqueness unproved, candidate ambiguous, lineage conflict, and source
unavailable. Focused fake tests prove identifier allowlisting/quoting, bound
values, bounded limits, exact TOP and `NUNOTA`/`SEQUENCIA`/`TGFVAR` predicates,
validation ordering, nullable facts, no descendants, and one-to-many/partial/
conflict state. Tests perform no live Oracle, provider, Postgres, dependency, or
network access.

## Acceptance Criteria

### F11-AC01 Configuration fails closed

Invalid configuration, missing/incompatible metadata, absent uniqueness
attestation, duplicate nonblank values, and unavailable Oracle prevent all data
queries and return stable typed errors.

### F11-AC02 SQL is identifier-safe, bound, and bounded

Tests prove the configured identifier is strictly allowlisted and quoted, every
value is a bind, all read overflows are detected, and exact TOP 313/306 plus
`NUNOTA`/`SEQUENCIA`/`TGFVAR` predicates are present.

### F11-AC03 Adapter rows do not leak

Only generic document/line identities, safe evidence, and nullable operational
facts cross the dedicated adapter port.

### F11-AC04 Lineage states stay honest

No descendants and one-to-many descendants remain explicit; missing attended
quantity, conflicting duplicates, and optional quantity mismatch never become a
zero/default identity or tax fact.

### F11-AC05 Scope and side effects remain bounded

Focused fake Go tests and scoped diff inspection prove the separate
`internal_read` seam without changing the broad reader, runtime wiring,
OpenAPI/SDK, or accessing live systems.
