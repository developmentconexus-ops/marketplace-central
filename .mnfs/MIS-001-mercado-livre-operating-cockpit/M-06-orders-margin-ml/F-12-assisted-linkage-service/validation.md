# F-12 Assisted Sankhya Linkage Service Validation

## Verdict

Feature validation passed for the bounded fake/unit target. The orders service
requires one exact persisted Mercado Livre order, re-reads one explicitly
selected TOP 313 candidate, proves the complete stable-MPC-line/candidate-line
bijection before append, records exact evidence with actor provenance
`operator_supplied_unverified`, and reads descendants from every line of the
persisted/idempotently returned aggregate. Runtime configuration revision and
server-generated event identity are not caller-controlled. Exact lineage
conflicts remain `conflict`; only configuration/source failures become
`unavailable`, and neither undoes or manufactures the confirmed mapping.

## Context evidence

- Accepted base SHA: `50ac8baee6c658bc987ab13941e47137a37c17a5`.
- The initial truthful L1 compile with `portfolio-core` and `orders-margin`
  failed `CTX_TOKEN_BUDGET_EXCEEDED` at `estimated-input-tokens`.
- A second truthful L1 compile without those route expansions but with the full
  packet-required F-09/F-11 specs failed the same estimate check.
- Per the dispatch's explicit known-defect instruction, `context.json` was
  reduced to full Feature/spec/plan/AGENTS contracts after all packet-required
  and named-route sources had already been read. No risk, path, side effect,
  command, criterion, stop, or original `dispatch.json` fact was changed or
  falsified.
- The reduced L1 pack compiled and passed `context-validate
  -RequireCurrentBase` at the accepted base.

## Implementation evidence

- Orders domain owns generic candidate/line/selection/lineage/read-error types
  and the fail-closed bijection validator. Mutable product, quantity, date, and
  amount evidence never participates in identity proof.
- The application derives only `ml:v1:<installation>:<provider-order>` after
  exact tenant/installation/order lookup and exact `mercado_livre` validation.
  Listing performs no append.
- Confirmation re-validates configuration and candidates, matches exactly one
  explicit document with TOP 313, rejects legacy/ambiguous/invalid/missing/
  extra/duplicate mappings, and validates all audit fields before F-09 append.
- The F-09 aggregate uses the exact derived key, all selected lines,
  `evidence_state=exact`, and fixed actor type
  `operator_supplied_unverified`. Its configuration revision comes from the
  validated bridge and its event ID from an injectable service generator with
  a cryptographic default; caller input exposes neither field. Empty revision
  and generation failure stop before append. No authentication claim is made.
- The post-append loop uses only mappings returned by `AppendConfirmation`.
  It attempts every exact origin, preserves lineage errors and malformed
  responses as empty per-line `conflict`, and uses `unavailable` only for typed
  configuration/source failures while continuing.
- The `internalread` adapter translates all F-11 DTOs, states, and errors into
  orders-owned values. Tests prove internal error types do not leak.
- The Postgres lookup requires exact tenant equality and exact scoped
  installation/order plus provider code. It was compile-checked only; no
  Postgres query or integration test executed.

## go-assisted-linkage-service

Command from `apps/server_core` with repository-local `GOCACHE=.gocache`:

```text
go test -count=1 \
  ./internal/modules/orders/domain \
  ./internal/modules/orders/application \
  ./internal/modules/orders/adapters/internalread
```

Result: exit 0.

```text
ok marketplace-central/apps/server_core/internal/modules/orders/domain
ok marketplace-central/apps/server_core/internal/modules/orders/application
ok marketplace-central/apps/server_core/internal/modules/orders/adapters/internalread
```

The focused tests cover exact-key/no-write listing, exact order/provider scope,
wrong/non-candidate TOP, legacy and ambiguous lines, missing/extra/duplicate/
different-document mappings, complete valid append audit, idempotent persisted
aggregate return, repository conflict propagation, exact per-line descendant
reads, runtime revision and event generation failures before append, exact
conflict/malformed/unavailable classification, nullable conversion, and
internal error translation.

Additional compile-only evidence:

```text
go test -tags=integration -run '^$' \
  ./internal/modules/orders/adapters/postgres
```

Result: exit 0, `[no tests to run]`; this compiled the exact lookup and its
scoping test without opening or reading a database.

## git-diff-check

- `git diff --check`: exit 0 before staging.
- Production scan of the new domain/application/bridge files found no `panic`,
  HTTP, OpenAPI, or SDK runtime reference.
- Changed paths are limited to the dispatched F-12 directory and the ten exact
  orders domain/ports/application/internal-read bridge/Postgres lookup files.
- Existing unrelated untracked `docs/research/**` and `output/**` paths were
  not modified, staged, or included.
- No live Oracle, provider, Postgres, network, dependency installation, secret,
  PII, transport, composition, runtime registry, migration, OpenAPI, SDK, or UI
  operation ran.

## Acceptance mapping

- F12-AC01: exact persisted ML order and derived-key listing tests pass with no
  append.
- F12-AC02: domain/service adversarial tests prove the full exact bijection and
  fail before append.
- F12-AC03: exact audit, fixed unverified actor provenance, runtime revision,
  server event generation, failure-before-append, persisted retry result, and
  conflict propagation tests pass.
- F12-AC04: every persisted origin is attempted; exact conflict, malformed,
  unavailable, and complete states preserve the confirmed mapping.
- F12-AC05: focused fake/unit proof, compile-only lookup proof, and scoped diff
  inspection pass without forbidden seams or side effects.
