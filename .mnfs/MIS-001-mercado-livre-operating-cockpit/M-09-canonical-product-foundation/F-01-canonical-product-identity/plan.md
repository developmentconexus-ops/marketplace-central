# F-01 Canonical Product Identity — Plan

## Ownership

- Allowed: `apps/server_core/internal/modules/catalog/**`,
  `apps/server_core/internal/modules/internal_read/domain/**`,
  `apps/server_core/internal/modules/internal_read/ports/**`,
  `contracts/api/marketplace-central.openapi.yaml`, `packages/sdk-runtime/**`,
  and this feature directory.
- Forbidden: adapter Oracle/MSDB work, composition, migrations, web, M-06, and
  provider/auth/write changes.
- Shared seams: catalog/internal_read identity types and atomic OpenAPI/SDK.

## Steps

1. Add canonical catalog identity/source-fact domain types with constructors
   that reject non-positive identities and invalid unknown/stale facts.
2. Replace the unqualified internal-read candidate ID with the canonical typed
   identity and adjust deterministic fixtures/tests.
3. Define the matching nullable canonical product contract in OpenAPI and the
   SDK, including quality/source metadata and separate identifiers.
4. Add serialization/parity tests: CODPROD 1001 remains separate from EAN,
   reference and seller SKU; missing numeric values remain null; known zero is
   preserved.
5. Compile and validate the bounded context; run registered targeted Go, SDK,
   and OpenAPI/SDK parity proof; record results in `validation.md` and commit.

## Risks And Stops

No legacy string identity is converted in this feature. Any required conversion,
zero/default unknown fact, OpenAPI/SDK divergence, or need to edit an excluded
path stops the work for Milestone direction.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-01",
  "required_sources": [],
  "knowledge_route_ids": ["portfolio-core"],
  "allowed_paths": [
    "apps/server_core/internal/modules/catalog/**",
    "apps/server_core/internal/modules/internal_read/domain/**",
    "apps/server_core/internal/modules/internal_read/ports/**",
    "contracts/api/marketplace-central.openapi.yaml",
    "packages/sdk-runtime/**",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-09-canonical-product-foundation/F-01-canonical-product-identity/**"
  ],
  "forbidden_paths": [
    "apps/web/**",
    "apps/server_core/internal/platform/msdb/**",
    "apps/server_core/cmd/**",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/**"
  ],
  "side_effects": {
    "allowed": ["repository-write", "isolated-cache-write"],
    "forbidden": ["database-mutation", "external-network", "provider-write"]
  },
  "commands": [
    {"id": "targeted-go-tests", "command_id": "targeted-go-tests", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "sdk-runtime-tests", "command_id": "sdk-runtime-tests", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "openapi-sdk-parity", "command_id": "openapi-sdk-parity", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F01-AC01", "milestone_criterion_id": "M-09-C01", "command_ids": ["targeted-go-tests"]},
    {"id": "F01-AC02", "milestone_criterion_id": "M-09-C02", "command_ids": ["targeted-go-tests", "sdk-runtime-tests"]},
    {"id": "F01-AC03", "milestone_criterion_id": "M-09-C01", "command_ids": ["sdk-runtime-tests", "openapi-sdk-parity"]}
  ],
  "stop_conditions": [
    {"code": "legacy-identity-unproved", "condition": "A legacy string identity would be mapped to CODPROD without deterministic proof."},
    {"code": "zero-default", "condition": "An unknown numeric source fact would be represented as zero or another default."},
    {"code": "boundary-leak", "condition": "An Oracle or provider type would cross into catalog or internal-read domain/application."},
    {"code": "api-sdk-drift", "condition": "OpenAPI and sdk-runtime cannot be updated atomically."},
    {"code": "path-outside-scope", "condition": "A required change falls outside the dispatched paths."}
  ],
  "retry_budget": {"max_correction_attempts": 2},
  "handoff_fields": ["status", "commit", "changed-paths", "evidence", "blockers", "next"]
}
```
