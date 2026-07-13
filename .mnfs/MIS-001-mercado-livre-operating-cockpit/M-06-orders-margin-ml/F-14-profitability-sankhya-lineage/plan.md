# F-14 Profitability Sankhya Lineage Plan

## Fixed decisions

- Owner: orders application owns read-only current confirmation/descendant
  resolution; profitability owns stable line projection, exact tax-source
  consumption, aggregation, provenance, and margin quality.
- Ports: extend the canonical assisted-linkage application boundary with one
  exact scoped current-lineage read. Profitability consumes it through a
  profitability-owned narrow port; the existing exact-line internal tax port
  remains unchanged unless a compatible provenance collection must be added
  inside the allowed profitability package.
- Consumer: profitability calculation/service is the only business consumer;
  composition supplies the adapter only when the assisted-linkage runtime is
  already available.
- Legacy decision: TOP 313, product/date matching, resolved product linkage
  alone, ambiguous/legacy reconciliation, and empty source identities are not
  valid tax provenance. They remain fail-closed.
- Identity: only exact positive, unique TOP 306 descendant document/line pairs
  produced by the confirmed stable MPC line are tax sources.
- Aggregation: sort exact descendants deterministically; aggregate a component
  only when every descendant has that component. Missing values remain nil.
- Unknown states: preserve `none`, `partial`, `complete`, `conflict`, and
  `unavailable`; partial may expose known exact tax but never complete margin.
- Side effects: current-lineage resolution is read-only; tests use fakes and
  repository-local `GOCACHE=.gocache`, never live Oracle/Postgres/provider.

## Implementation steps

1. Compile and validate `context.json` at the accepted base SHA from this
   plan, dispatch, and routed selectors; read only the compiled selectors.
2. Inspect the existing orders assisted-linkage service, ledger/reader ports,
   profitability order-fact and internal-read boundaries, calculation flow,
   tests, and composition construction named by the compiled context.
3. Add the exact read-only current-lineage operation and focused orders tests
   proving scope, confirmation lookup, mapping selection, state preservation,
   exact descendant re-read, and no append.
4. Carry stable `MPCLineID` into profitability order facts and add the narrow
   current-lineage adapter/port without leaking orders or internal-read DTOs
   into profitability domain values.
5. Resolve tax only from valid unique TOP 306 descendants, preserve sorted
   provenance, aggregate each component only across all descendants, and keep
   partial/unknown results incomplete. Add focused profitability tests for
   exact calls, one-to-many, partial values, conflicts, invalid/duplicates,
   missing mapping, and no TOP 313 calls.
6. Wire the boundary in the composition root only alongside the existing
   assisted-linkage runtime; add focused root tests if present or use compile
   proof when construction is already covered.
7. Run `go-orders-profitability-lineage`, `go-repository`, and
   `git-diff-check` with repository-local `GOCACHE=.gocache`; write
   `validation.md`, stage only allowed paths, and create one intentional
   commit.

## Registered proof

- `go-orders-profitability-lineage`: focused orders application,
  profitability, and composition tests proving exact read-only lineage and
  deterministic honest tax aggregation. Fake/unit lane only.
- `go-repository`: repository-wide Go tests from `apps/server_core` using
  `GOCACHE=.gocache`, with no live integrations.
- `git-diff-check`: whitespace validation plus changed/staged path inspection
  against the dispatch allowlist.

Stop if implementation would use TOP 313 or heuristics for tax, represent
partial/missing/conflict/unavailable as complete, aggregate an unknown
component as zero, require an out-of-scope contract/migration/adapter edit, or
touch a live system.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-14",
  "required_sources": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-contract.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-07-order-specific-tax-provenance/spec.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-09-order-line-identity-linkage-ledger/spec.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-11-sankhya-linkage-read-contract/spec.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-12-assisted-linkage-service/spec.md",
    "contracts/governance/shared-seams.json"
  ],
  "knowledge_route_ids": ["portfolio-core", "orders-margin"],
  "allowed_paths": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-14-profitability-sankhya-lineage/**",
    "apps/server_core/internal/modules/orders/application/assisted_sankhya_linkage_service.go",
    "apps/server_core/internal/modules/orders/application/assisted_sankhya_linkage_service_test.go",
    "apps/server_core/internal/modules/profitability/**",
    "apps/server_core/internal/composition/root.go"
  ],
  "forbidden_paths": [
    "apps/server_core/migrations/**",
    "apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go",
    "contracts/api/**",
    "packages/sdk-runtime/**",
    "apps/web/**",
    "scripts/**",
    "docs/research/**",
    "docker/**"
  ],
  "side_effects": {
    "allowed": ["repository-write", "isolated-cache-write"],
    "forbidden": ["database-mutation", "external-network", "provider-write"]
  },
  "commands": [
    {"id": "go-orders-profitability-lineage", "command_id": "go-orders-profitability-lineage", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "go-repository", "command_id": "go-repository", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "git-diff-check", "command_id": "git-diff-check", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F14-AC01", "milestone_criterion_id": "M-06-C02", "command_ids": ["go-orders-profitability-lineage"]},
    {"id": "F14-AC02", "milestone_criterion_id": "M-06-C02", "command_ids": ["go-orders-profitability-lineage"]},
    {"id": "F14-AC03", "milestone_criterion_id": "M-06-C02", "command_ids": ["go-orders-profitability-lineage"]},
    {"id": "F14-AC04", "milestone_criterion_id": "M-06-C02", "command_ids": ["go-orders-profitability-lineage", "go-repository"]},
    {"id": "F14-AC05", "milestone_criterion_id": "M-06-C02", "command_ids": ["go-orders-profitability-lineage", "go-repository", "git-diff-check"]}
  ],
  "stop_conditions": [
    {"code": "invalid-tax-source", "condition": "TOP 313, a product/date heuristic, or an invalid/duplicate descendant would be used for tax."},
    {"code": "unknown-promoted", "condition": "Partial, missing, conflict, or unavailable lineage would be represented as complete."},
    {"code": "zero-default", "condition": "An unknown tax component would be aggregated as zero."},
    {"code": "path-outside-scope", "condition": "A required edit falls outside the dispatched paths."},
    {"code": "forbidden-side-effect", "condition": "Validation would access a live system, install dependencies, or expose secrets/PII."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "commit", "changed-paths", "evidence", "blockers", "next"]
}
```
