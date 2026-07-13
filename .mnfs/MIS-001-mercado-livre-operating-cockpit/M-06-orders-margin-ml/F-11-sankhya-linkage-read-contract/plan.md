# F-11 Sankhya Linkage Read Contract Plan

## Fixed decisions

- Owner: `internal_read` domain, dedicated port/application service, and Oracle
  adapter files named `sankhya_linkage*`.
- Consumer: future orders linkage code through the dedicated port; there is no
  composition or current runtime consumer in this feature.
- Port: `SankhyaLinkageReader` is separate from and does not extend `Reader`.
- Legacy/effective-TOP decision: use and return the proved `TGFCAB.CODTIPOPER`
  header operation-code fact. Do not invent an unproved effective-TOP history
  column or join.
- Unknowns: nullable product/date/quantity/value evidence, nullable attended
  quantity, and optional expected origin quantity remain nil. No zero/default
  identity or tax fact is derived.
- Activation: validate complete explicit configuration, Oracle availability,
  exact nullable-CLOB metadata, nonblank uniqueness attestation,
  and duplicate-free nonblank values before every data read.
- Identifier/value safety: only already-uppercase Oracle-safe schema/field
  identifiers are quoted; every lookup value, TOP, document/line identity,
  digits-only external key, and limit is bound; CLOB equality is exact via
  `DBMS_LOB.COMPARE(..., TO_CLOB(:bind))`.
- Candidate bounds: fetch header and per-header line limits plus one and fail
  ambiguous on either overflow; read exact lines only for retained TOP 313
  headers.
- Descendant bound: fetch the explicit lineage limit plus one, fail conflict on
  overflow, preserve bounded exact `TGFVAR` one-to-many destinations, and derive
  none/partial/complete/conflict without collapsing rows.
- Panic policy: every query builder returns a typed configuration error for an
  invalid identifier; changed production files contain no `panic` call.

## Implementation steps

1. Compile and validate `context.json` from this plan, dispatch, and routed
   sources; then read only its selectors.
2. Add domain models, state derivation, stable errors, and focused tests for
   generic candidates/lines/descendants and nullable operational facts.
3. Add the separate `ports.SankhyaLinkageReader` inputs and interface.
4. Add an application service that validates configuration before delegating
   candidate and descendant reads.
5. Add an independent Oracle adapter with explicit linkage configuration,
   strict identifier validation/quoting, exact metadata and duplicate probes,
   bound/bounded candidate and line queries, and exact `TGFVAR` descendant SQL.
6. Add fake database/query-builder tests for validation ordering, SQL safety,
   exact predicates, limits, row mapping, and lineage states.
7. Run `go-sankhya-linkage-read` with `GOCACHE=.gocache`, inspect the scoped
   diff with `git-diff-check`, write `validation.md`, and create one intentional
   commit.

## Allowed paths

- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-11-sankhya-linkage-read-contract/**`
- `apps/server_core/internal/modules/internal_read/domain/sankhya_linkage.go`
- `apps/server_core/internal/modules/internal_read/domain/sankhya_linkage_test.go`
- `apps/server_core/internal/modules/internal_read/ports/sankhya_linkage.go`
- `apps/server_core/internal/modules/internal_read/application/sankhya_linkage_service.go`
- `apps/server_core/internal/modules/internal_read/application/sankhya_linkage_service_test.go`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/sankhya_linkage_reader.go`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/sankhya_linkage_reader_test.go`

## Registered proof and stops

- `go-sankhya-linkage-read`: focused internal-read domain/application/Oracle
  fake tests, network and databases disabled, isolated `GOCACHE=.gocache`.
- `git-diff-check`: status/diff inspection for exact allowed-path ownership and
  forbidden-side-effect absence.

Stop on a repository-truth conflict, a need for an unproved Oracle schema fact,
unsafe interpolation, adapter DTO leakage, an out-of-scope edit, or any live
Oracle/provider/Postgres/dependency/secret/PII access.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-11",
  "required_sources": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/implementation-contract.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/sankhya-admin-spec.md"
  ],
  "allowed_paths": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-11-sankhya-linkage-read-contract/**",
    "apps/server_core/internal/modules/internal_read/domain/sankhya_linkage.go",
    "apps/server_core/internal/modules/internal_read/domain/sankhya_linkage_test.go",
    "apps/server_core/internal/modules/internal_read/ports/sankhya_linkage.go",
    "apps/server_core/internal/modules/internal_read/application/sankhya_linkage_service.go",
    "apps/server_core/internal/modules/internal_read/application/sankhya_linkage_service_test.go",
    "apps/server_core/internal/modules/internal_read/adapters/oracle/sankhya_linkage_reader.go",
    "apps/server_core/internal/modules/internal_read/adapters/oracle/sankhya_linkage_reader_test.go"
  ],
  "forbidden_paths": [
    "apps/server_core/internal/modules/orders/**",
    "apps/server_core/internal/modules/profitability/**",
    "apps/server_core/internal/composition/**",
    "apps/server_core/migrations/**",
    "contracts/api/**",
    "packages/**",
    "contracts/governance/**",
    "scripts/**",
    "docs/research/**"
  ],
  "side_effects": {
    "allowed": ["repository-write", "isolated-cache-write"],
    "forbidden": ["database-mutation", "external-network", "provider-write"]
  },
  "commands": [
    {"id": "go-sankhya-linkage-read", "command_id": "go-sankhya-linkage-read", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "git-diff-check", "command_id": "git-diff-check", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F11-AC01", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-sankhya-linkage-read"]},
    {"id": "F11-AC02", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-sankhya-linkage-read"]},
    {"id": "F11-AC03", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-sankhya-linkage-read"]},
    {"id": "F11-AC04", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-sankhya-linkage-read"]},
    {"id": "F11-AC05", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-sankhya-linkage-read", "git-diff-check"]}
  ],
  "stop_conditions": [
    {"code": "oracle-fact-unproved", "condition": "Implementation would require inventing an unproved Oracle column, table, or effective-TOP join."},
    {"code": "unsafe-sql", "condition": "An identifier is not strictly validated/quoted or any value would be interpolated."},
    {"code": "boundary-leak", "condition": "An Oracle DTO or row would cross the adapter boundary."},
    {"code": "path-outside-scope", "condition": "A required edit falls outside the dispatched paths."},
    {"code": "forbidden-side-effect", "condition": "Validation would access a live database/provider, install dependencies, or expose secrets/PII."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "commit", "changed-paths", "evidence", "blockers", "next"]
}
```
