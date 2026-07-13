# F-12 Assisted Sankhya Linkage Service Plan

## Fixed decisions

- Owner: orders domain/application/ports plus the bounded orders adapters named
  in dispatch; `internal_read` remains owner of its F-11 service and types.
- Consumer: a future transport/composition feature; F-12 adds no current
  runtime consumer, HTTP route, public contract, or activation claim.
- Ports: `OrderLookup` loads one exact tenant/installation/provider-order; the
  orders-owned `SankhyaLinkageReader` exposes configuration, candidates, and
  descendants using only orders domain values; the existing F-09 repository
  port remains unchanged.
- Bridge: `orders/adapters/internalread` translates every F-11 value/state into
  orders types. No Oracle/internal-read DTO crosses into orders domain,
  application inputs/results, or the ledger.
- Proof rule: confirmation re-reads and matches one explicitly selected TOP
  313 document, then proves a bijection over every stable persisted MPC line and
  every exact selected-document line before append.
- Audit: actor type is fixed to `operator_supplied_unverified`; actor ID and
  intent facts are required, but this Feature makes no authentication or
  authorization claim. Configuration revision comes only from the validated
  reader bridge and event ID only from a server-side crypto-random generator;
  neither is caller-controlled. Evidence state is `exact` only after all
  validation and successful event generation.
- Idempotency: use the aggregate returned by F-09 append, including on an
  identical retry. Repository conflicts propagate without overwrite.
- Post-append failure: read descendants for every persisted mapping line;
  exact lineage errors and malformed responses become per-line `conflict`,
  only configuration/source errors become `unavailable`, and neither state
  fabricates descendants or undoes/manufactures the mapping.
- Unknowns: nullable candidate/line/descendant facts remain nil. No zero,
  product/date fallback, tax identity, or inferred operational fact is created.
- Postgres lookup: explicit tenant equality plus exact installation/order and
  `mercado_livre` provider predicates; focused proof does not connect to a DB.

## Implementation steps

1. Compile and validate `context.json` from this plan and the named routes. If
   the known L1 route estimate defect recurs, record the exact failed code and
   compile the reduced required-source pack without falsifying risk, paths,
   side effects, stops, or the already-read packet sources.
2. Add orders-owned candidate, confirmation-selection, lineage, unavailable
   state, and bijection validation types with focused domain tests.
3. Add the narrow orders reader and exact order lookup ports.
4. Add the application service with exact order/key checks, explicit candidate
   revalidation, runtime revision, server event generation, pre-append
   bijection/audit construction, F-09 append, and persisted-line descendant
   reads with exact failure classification.
5. Add the internal-read bridge with complete generic DTO/state conversion and
   fake tests proving nullable facts and error propagation.
6. Implement the exact tenant/installation/provider-order method on the
   existing Postgres order repository without wiring or live DB proof.
7. Run `go-assisted-linkage-service` with `GOCACHE=.gocache`, run
   `git-diff-check`, write `validation.md`, and create one intentional commit.

## Registered proof

- `go-assisted-linkage-service`: from `apps/server_core`, set repository-local
  `GOCACHE=.gocache` and run focused orders domain/application/internal-read
  adapter tests. Unit lane, no database or network.
- `git-diff-check`: run whitespace validation and inspect exact changed/staged
  paths against the dispatch allowlist.

Stop on proof promotion without explicit selection, an incomplete/non-bijective
mapping, DTO leakage, an authentication claim, an out-of-scope required edit,
or any live Oracle/provider/Postgres/dependency/secret/PII operation.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-12",
  "required_sources": [],
  "allowed_paths": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-12-assisted-linkage-service/**",
    "apps/server_core/internal/modules/orders/domain/assisted_sankhya_linkage.go",
    "apps/server_core/internal/modules/orders/domain/assisted_sankhya_linkage_test.go",
    "apps/server_core/internal/modules/orders/ports/sankhya_linkage_reader.go",
    "apps/server_core/internal/modules/orders/ports/order_lookup.go",
    "apps/server_core/internal/modules/orders/application/assisted_sankhya_linkage_service.go",
    "apps/server_core/internal/modules/orders/application/assisted_sankhya_linkage_service_test.go",
    "apps/server_core/internal/modules/orders/adapters/internalread/sankhya_linkage_reader.go",
    "apps/server_core/internal/modules/orders/adapters/internalread/sankhya_linkage_reader_test.go",
    "apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go",
    "apps/server_core/internal/modules/orders/adapters/postgres/order_repo_test.go"
  ],
  "forbidden_paths": [
    "apps/server_core/internal/modules/internal_read/**",
    "apps/server_core/internal/modules/profitability/**",
    "apps/server_core/internal/modules/orders/transport/**",
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
    {"id": "go-assisted-linkage-service", "command_id": "go-assisted-linkage-service", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "git-diff-check", "command_id": "git-diff-check", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F12-AC01", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-assisted-linkage-service"]},
    {"id": "F12-AC02", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-assisted-linkage-service"]},
    {"id": "F12-AC03", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-assisted-linkage-service"]},
    {"id": "F12-AC04", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-assisted-linkage-service"]},
    {"id": "F12-AC05", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-assisted-linkage-service", "git-diff-check"]}
  ],
  "stop_conditions": [
    {"code": "proof-not-explicit", "condition": "A candidate or filter would become proof without one explicitly selected exact TOP 313 document."},
    {"code": "bijection-incomplete", "condition": "Every stable persisted MPC line and every exact selected candidate line cannot be covered exactly once."},
    {"code": "boundary-leak", "condition": "An Oracle or internal-read type would cross the orders adapter boundary."},
    {"code": "authentication-claim", "condition": "The implementation would claim or infer production authentication or manual-adjustment authorization."},
    {"code": "path-outside-scope", "condition": "A required edit falls outside the dispatched paths."},
    {"code": "forbidden-side-effect", "condition": "Validation would access a live database/provider, install dependencies, or expose secrets/PII."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "commit", "changed-paths", "evidence", "blockers", "next"]
}
```
