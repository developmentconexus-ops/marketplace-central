# F-11 Sankhya Linkage Read Contract Validation

## Verdict

F-11 contract correction passed for the bounded fake/unit target. It preserves
MPC tenant/provider/installation ownership in the ledger boundary and explicit
ambiguity; it performs no live-system action. C03 remains deferred/failing and
is not authenticated by this correction.

## Correction

- `TGFCAB.AD_NUMPEDIDO_ECOM` metadata must resolve exactly once as nullable
  `CLOB`; no VARCHAR/NVARCHAR type or capacity is accepted.
- External order keys are nonblank digits-only.
- Candidate matching uses exact, parameterized CLOB comparison:
  `DBMS_LOB.COMPARE(c.<quoted field>, TO_CLOB(:1)) = 0`.
- The duplicate guard uses exact CLOB comparison between rows, avoiding an
  invalid CLOB group-by while keeping ambiguous duplicate keys explicit.

## Context evidence

- Accepted base SHA: `0cfa801b7f9cfe57c0cd81f7c953e81a8a706cbf`.
- `context.json` was recompiled from the updated feature contract and passed
  `context-validate -RequireCurrentBase`.
- The scoped pack retains the F-11 source/test seam and registered commands.

## go-sankhya-linkage-read

Command, from repository root with isolated cache:

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

Focused tests cover nullable-CLOB metadata acceptance/rejection, digits-only
keys, safe quoted identifiers, parameterized exact CLOB candidate comparison,
exact CLOB duplicate detection, bounds, TOP predicates, TGFVAR lineage, and
explicit ambiguity/lineage states.

## git-diff-check

- `git diff --check`: exit 0.
- Changed paths are confined to the dispatched F-11 artifact directory and
  permitted `internal_read` adapter/application test seam.
- No live Oracle, provider, Postgres, network, dependency installation, secret,
  or PII access ran.

## Acceptance mapping

- F11-AC01: nullable-CLOB metadata, attestation, uniqueness, and source checks
  remain fail-closed before data reads.
- F11-AC02: identifiers and values remain safe; CLOB equality is exact and
  parameterized.
- F11-AC03: generic domain/port boundary remains unchanged.
- F11-AC04: existing explicit one-to-many and ambiguity states remain covered.
- F11-AC05: focused proof and scoped diff inspection pass.
