# F-11 Sankhya Linkage Read Contract Validation

## Verdict

Feature validation passed for the bounded fake/unit target. The implementation
adds only the separate `internal_read` Sankhya-linkage domain, port, service,
Oracle adapter, and focused tests. It does not extend the existing broad
`Reader`, wire composition, change OpenAPI/SDK, or access a live system.

## Context evidence

- Accepted base SHA: `011cb2d3f53135e72911eac773483c954436cd04`.
- `context.json` compiled and passed `Test-HarnessContextPack
  -RequireCurrentBase` at L1 with truthful `repository-write` and
  `isolated-cache-write` side effects.
- Exact compiled selectors: full `feature.md`, `spec.md`, `plan.md`, `AGENTS.md`,
  F-08 `implementation-contract.md`, and F-08 `sankhya-admin-spec.md`.
- After the focused review correction, `context.json` was recompiled to refresh
  changed contract hashes and passed canonical validation against the original
  accepted base (without claiming that the post-feature HEAD equals that base).
- The original `portfolio-core` plus `orders-margin` route expansion was omitted
  under Milestone authorization because its registered selectors alone estimate
  2,600 tokens, while the canonical compiler rejects every truthful single-module
  L1 pack above 2,000. The packet-required sources had already been read. Risk,
  side effects, paths, and the original `dispatch.json` were not widened or
  falsified.

## Contract and schema-fact decisions

- Configuration has no default schema, header field, revision, TOP, uniqueness
  attestation, candidate, candidate-line, or lineage limit. Every limit is
  explicit, positive, and capped by a fixed maximum.
- Schema and field identifiers must already match the strict uppercase
  Oracle-safe allowlist and are quoted before use. Metadata lookup and every
  request/configuration value use binds.
- Metadata must resolve exactly once to `VARCHAR2`/`NVARCHAR2` with character
  semantics and capacity at least 160. A duplicate-nonblank aggregate probe and
  explicit attestation are required before every candidate or lineage read.
- Candidate reads bind the exact account-scoped external key and TOP 313, fetch
  configured header and line limits plus one, reject either overflow as
  candidate ambiguous, and read lines by exact `NUNOTA`.
- Descendant reads bind exact origin `NUNOTA`/`SEQUENCIA`, traverse `TGFVAR`,
  require destination header TOP 306, fetch the lineage limit plus one, reject
  overflow as conflict, and preserve every bounded destination line plus
  nullable `QTDATENDIDA`.
- Repository truth does not prove an effective-TOP history column/join. The
  adapter therefore returns and checks the proved `TGFCAB.CODTIPOPER` header
  operation-code fact and makes no effective-time TOP claim. Activation remains
  fail-closed.

## go-sankhya-linkage-read

Command, from `apps/server_core`, with repository-local cache:

```text
GOCACHE=<repository>/.gocache go test -count=1 \
  ./internal/modules/internal_read/domain \
  ./internal/modules/internal_read/application \
  ./internal/modules/internal_read/adapters/oracle
```

Result: exit 0.

```text
ok marketplace-central/apps/server_core/internal/modules/internal_read/domain
ok marketplace-central/apps/server_core/internal/modules/internal_read/application
ok marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle
```

Evidence covers invalid identifier and missing attestation before DB access,
missing/incompatible metadata and duplicate values before candidate SQL, quoted
identifier/bind/limit construction, header/line/lineage overflow, exact TOP and
document-line predicates, nil DB typing, generic candidate/line models, unknown
expected quantity remaining partial, and none/partial/complete/conflict lineage
including one-to-many nullable descendants and duplicate-identity rejection.

## Focused review correction

- Removed the production `mustQuoteSankhyaIdentifier` panic path. Every dynamic
  query builder now returns a typed `configuration_invalid` error.
- An AST-based governance test parses all changed production Go files and fails
  on any `panic` call.
- Added bound-plus-one SQL and overflow handling for candidate lines and
  descendants, with all limit values bound.
- Corrected lineage derivation so an unknown expected origin quantity remains
  `partial` and any repeated destination identity is `conflict`, even when its
  observed quantity is identical.

## git-diff-check

- `git diff --check`: exit 0 before staging.
- Changed paths are exactly the dispatched F-11 directory and seven new bounded
  `internal_read` source/test files.
- Existing unrelated untracked `docs/research/**` and `output/**` paths were not
  read, modified, staged, or included.
- No live Oracle, provider, Postgres, network, dependency installation, secret,
  or PII operation ran. All database behavior is deterministic fake-driver
  contract evidence only.

## Acceptance mapping

- F11-AC01: typed fail-closed configuration, metadata, uniqueness, and source
  tests pass before data-query execution.
- F11-AC02: strict quoted identifiers, typed no-panic builder errors, bound
  values, bounded candidates/lines/lineage, and exact TOP/identity/TGFVAR SQL
  tests pass.
- F11-AC03: port and service expose generic domain values only.
- F11-AC04: domain and adapter tests preserve nullable and bounded one-to-many
  lineage, unknown expected quantity, and duplicate conflicts with explicit
  states and no zero/default tax identity.
- F11-AC05: focused Go proof and scoped path inspection pass; forbidden runtime
  and public seams remain unchanged.
