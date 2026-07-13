# F-09 Stable Order-Line Identity and Linkage Ledger Plan

## Fixed decisions

- Owner: orders domain/ports/Postgres adapter.
- Consumer now: the existing Postgres order refresh path for stable line
  reconciliation. Consumer later: the orders application linkage service via
  `SankhyaLinkageRepository`.
- Legacy decision: migrated rows get opaque IDs but remain
  `legacy_unresolved`; neither `line_no` nor an ambiguous attribute tuple is
  proof.
- Reconciliation key: provider item ID plus variation and seller-SKU
  presence/value only. Quantity, unit price, and `line_no` are mutable snapshot
  evidence, not identity proof.
- Duplicate decision: preserve the prior opaque ID multiset deterministically,
  allocate IDs only for new excess members, and keep every member ambiguous and
  non-linkable without claiming an individual pairing.
- Unknown states: `legacy_unresolved`, `ambiguous`, and linkage evidence
  `unknown` remain explicit. They never become zero/default identities.
- Transaction boundary: one repository append owns the header event and every
  line event. The scoped idempotency key and all uniqueness checks are resolved
  inside that transaction.
- External key: confirmation requires the exact
  `ml:v1:<installation_id>:<provider_order_id>` value; Postgres stores and
  uniquely constrains it within tenant/installation scope for later TOP 313
  header comparison.
- Refresh persistence: retain rows in place by `mpc_line_id`, insert only new
  identities, delete only removed unlinked rows, and fail the transaction with
  an explicit identity conflict if a removed row is linkage-protected.
- Port shape: load current mapping plus append confirmation; domain values use
  generic document and document-line identities.
- No application, transport, composition, OpenAPI/SDK, Oracle, provider, or
  production changes are in scope.

## Implementation steps

1. Compile and validate `context.json` from this plan and the dispatched
   knowledge routes. Read only its selectors before implementation.
2. Add migration 0033:
   - add immutable `mpc_line_id` and reconciliation state to order items;
   - backfill existing rows as `legacy_unresolved` without positional proof;
   - add scoped uniqueness needed for identity ownership;
   - create append-only header/line event storage, current projections if
     required, scoped idempotency and origin-line uniqueness, and audit fields;
   - add triggers preventing ledger update/delete.
3. Add domain types and tests for opaque line IDs, reconciliation state,
   generic internal identities, audit metadata, append validation, and
   semantic equality.
4. Extend order refresh persistence and tests:
   - assign fresh opaque IDs on first import;
   - preserve IDs for a unique provider item/variation/seller-SKU identity even
     when quantity, unit price, order, or `line_no` changes;
   - preserve the existing ID multiset for duplicate identity groups and create
     IDs only for new excess rows;
   - mark all indistinguishable duplicates ambiguous without using `line_no` as
     identity proof;
   - update retained identities in place across reorder, preserving linkage
     foreign keys and row creation audit;
   - reject removal of a confirmed line with an explicit conflict before any
     partial refresh commits;
   - retain stale-refresh behavior.
5. Add the orders-owned repository port and Postgres adapter. Insert the
   mapping plus exact external order key atomically, return an identical scoped
   retry, and fail closed on semantic or uniqueness conflicts without overwrite.
6. Add focused repository tests. Use only the approved ephemeral Postgres
   target when configured; otherwise record the integration proof as not run,
   not passed.
7. Update the migration embed-count test, run registered commands, write
   `validation.md`, inspect the scoped diff, and create one intentional commit.

## Allowed paths

- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-09-order-line-identity-linkage-ledger/**`
- `apps/server_core/migrations/0033_orders_sankhya_linkage.sql`
- `apps/server_core/internal/platform/migrate/runner_test.go`
- `apps/server_core/internal/modules/orders/domain/order.go`
- `apps/server_core/internal/modules/orders/domain/sankhya_linkage.go`
- `apps/server_core/internal/modules/orders/domain/sankhya_linkage_test.go`
- `apps/server_core/internal/modules/orders/ports/sankhya_linkage_repository.go`
- `apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go`
- `apps/server_core/internal/modules/orders/adapters/postgres/order_repo_test.go`
- `apps/server_core/internal/modules/orders/adapters/postgres/sankhya_linkage_repo.go`
- `apps/server_core/internal/modules/orders/adapters/postgres/sankhya_linkage_repo_test.go`

## Registered proof

- `go-orders-linkage-unit`: focused orders domain and order repository tests
  with `GOCACHE=.gocache`.
- `go-orders-linkage-postgres`: focused linkage adapter tests against the
  approved ephemeral Postgres target, when available.
- `go-migration-embed`: migration embedding/count test with
  `GOCACHE=.gocache`.
- `git-diff-check`: scoped status/diff review proving allowed-path ownership
  and no forbidden side effects.

Evidence is recorded in `validation.md`; no raw secret, PII, provider payload,
or external-system output is retained.

## Stop conditions

Stop and classify the conflict if stable identity requires positional or
ambiguous proof, atomic scoped append cannot be enforced, an out-of-scope path
is required, repository truth conflicts, or validation would access Oracle,
provider, production, secrets, or PII.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-09",
  "required_sources": [
    "ARCHITECTURE.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/mission.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/milestone.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-contract.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/implementation-contract.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/sankhya-admin-spec.md",
    "contracts/governance/execution-lanes.json",
    "contracts/governance/shared-seams.json"
  ],
  "knowledge_route_ids": ["portfolio-core", "orders-margin"],
  "allowed_paths": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-09-order-line-identity-linkage-ledger/**",
    "apps/server_core/migrations/0033_orders_sankhya_linkage.sql",
    "apps/server_core/internal/platform/migrate/runner_test.go",
    "apps/server_core/internal/modules/orders/domain/order.go",
    "apps/server_core/internal/modules/orders/domain/sankhya_linkage.go",
    "apps/server_core/internal/modules/orders/domain/sankhya_linkage_test.go",
    "apps/server_core/internal/modules/orders/ports/sankhya_linkage_repository.go",
    "apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go",
    "apps/server_core/internal/modules/orders/adapters/postgres/order_repo_test.go",
    "apps/server_core/internal/modules/orders/adapters/postgres/sankhya_linkage_repo.go",
    "apps/server_core/internal/modules/orders/adapters/postgres/sankhya_linkage_repo_test.go"
  ],
  "forbidden_paths": [
    "apps/server_core/internal/modules/internal_read/**",
    "apps/server_core/internal/modules/orders/application/**",
    "apps/server_core/internal/modules/orders/transport/**",
    "apps/server_core/internal/composition/**",
    "contracts/api/**",
    "packages/sdk-runtime/**",
    "docs/research/**",
    "contracts/governance/**"
  ],
  "side_effects": {
    "allowed": ["repository-write", "database-write", "isolated-cache-write"],
    "forbidden": ["external-network", "provider-write"]
  },
  "commands": [
    {"id": "go-orders-linkage-unit", "command_id": "go-orders-linkage-unit", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "go-orders-linkage-postgres", "command_id": "go-orders-linkage-postgres", "lane_id": "integration", "expected_exit_code": 0},
    {"id": "go-migration-embed", "command_id": "go-migration-embed", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "git-diff-check", "command_id": "git-diff-check", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F09-AC01", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-orders-linkage-unit", "go-orders-linkage-postgres"]},
    {"id": "F09-AC02", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-orders-linkage-unit"]},
    {"id": "F09-AC03", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-migration-embed", "go-orders-linkage-postgres"]},
    {"id": "F09-AC04", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-orders-linkage-unit", "go-orders-linkage-postgres"]},
    {"id": "F09-AC05", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-orders-linkage-postgres", "go-migration-embed"]},
    {"id": "F09-AC06", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-orders-linkage-unit", "go-migration-embed", "git-diff-check"]}
  ],
  "stop_conditions": [
    {"code": "identity-proof-ambiguous", "condition": "Stable identity would require treating mutable line_no or an ambiguous attribute tuple as proof."},
    {"code": "atomic-scope-unproved", "condition": "The mapping cannot be inserted atomically and idempotently under tenant and installation scope."},
    {"code": "path-outside-scope", "condition": "A change outside the allowed persistence and domain paths is required."},
    {"code": "truth-conflict", "condition": "An architecture, contract, ownership, runtime, or verification conflict is found."},
    {"code": "forbidden-side-effect", "condition": "Any command would access Oracle, a provider, production data, secrets, or PII."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "commit", "changed-paths", "evidence", "blockers", "next"]
}
```
