# M-08 Governance Runtime Inventory

```yaml
id: R-002
type: research
status: verified
owner: Codebase Investigator
parent: MIS-001
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: support
```

## Topic

Current machine-verifiable module, runtime configuration, invariant, seam, and
PowerShell capabilities needed by M-08/F-07.

## Sources Checked

- Active Go/web/scripts/Docker source via targeted environment-read and import
  scans; excludes `.env`, secrets, historical plans, and generated output.
- `AGENTS.md`, `ARCHITECTURE.md`, OpenAPI route prefixes, composition root,
  migrations, and active module directories.
- Installed `pwsh 7.6.2` plus Microsoft PowerShell documentation for
  `Test-Json -SchemaFile` and `ConvertFrom-Json -AsHashtable`.

## Findings

### Module inventory

Exactly 11 active roots exist: `catalog`, `classifications`, `connectors`,
`integrations`, `internal_read`, `inventory`, `marketplaces`, `orders`,
`pricing`, `product_links`, and `profitability`. Each is imported by
`apps/server_core/internal/composition/root.go`.

Observed module edges must be declared, but forbidden target layers remain
violations. Exact current forbidden-layer edges become temporary exceptions
owned by M-10; new edges fail. Known independent exceptions:

- production `panic()` at
  `apps/server_core/internal/modules/product_links/application/resolution_service.go`
  is owned for removal by M-09;
- duplicate migration prefix `0021` is owned for correction before M-08/F-03;
- PostgreSQL `database/sql` prohibition applies only under
  `adapters/postgres`; Oracle driver use remains valid inside Oracle adapters.

### Runtime configuration inventory

| Canonical key(s) | Owner | Sensitivity | Lifecycle / removal |
| --- | --- | --- | --- |
| `SERVER_ADDR`, `API_PORT` | `platform/config` | public | active; `API_PORT` precedence |
| `MC_DATABASE_URL`, `MC_DEFAULT_TENANT_ID`, `MPC_ENCRYPTION_KEY` | `platform/pgdb` | secret/internal | active |
| `MC_MIGRATIONS_DIR` | `cmd/migrate` | public path | active direct-reader exception; F-03 |
| `RUN_MIGRATIONS` | `docker/dev` | public flag | active edge reader |
| `MS_DATABASE_URL`, `MS_TENANT_ID` | `platform/msdb` | secret/internal | legacy-current; remove M-09 |
| `MPC_ORACLE_USERNAME`, `MPC_ORACLE_PASSWORD`, `MPC_ORACLE_CONNECT_STRING`, `MPC_ORACLE_LIB_DIR` | `internal_read/adapters/oracle` | internal/secret | active typed owner |
| `MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID`, `MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET` | ML integration adapter | internal/secret | active direct-reader exception; M-10 |
| `MPC_PROVIDER_MAGALU_CLIENT_ID`, `MPC_PROVIDER_MAGALU_CLIENT_SECRET` | Magalu integration adapter | internal/secret | active direct-reader exception; M-10 |
| `MPC_PROVIDER_AMAZON_APPLICATION_ID`, `MPC_PROVIDER_AMAZON_CLIENT_ID`, `MPC_PROVIDER_AMAZON_CLIENT_SECRET`, `MPC_PROVIDER_AMAZON_AUTH_VERSION` | Amazon integration adapter | internal/secret/public | active direct-reader exception; M-10 |
| `MPC_PROVIDER_SHOPEE_PARTNER_ID`, `MPC_PROVIDER_SHOPEE_PARTNER_KEY`, `MPC_PROVIDER_SHOPEE_BASE_URL` | Shopee integration adapter | internal/secret/public | active direct-reader exception; M-10 |
| `MPC_OAUTH_REDIRECT_URI`, `MPC_WEB_ORIGIN` | integrations transport/web edge | internal URL | active direct-reader exception; M-10 |
| `ME_CLIENT_ID`, `ME_CLIENT_SECRET`, `ME_REDIRECT_URI` | Melhor Envio connector adapter | internal/secret | active direct-reader exception; M-10 |
| `MPC_WEB_PROXY_TARGET`, `VITE_API_BASE_URL` | web/Vite edge | internal/public URL | active |
| `NGROK_AUTHTOKEN` | Docker ngrok edge | secret | active dev/browser edge |
| `MPC_ORACLE_LIVE_TEST` | Oracle live test | public flag | lane-scoped test gate |
| `MPC_PROVIDER_MERCADOLIVRE_ACCESS_TOKEN`, `MPC_PROVIDER_MERCADOLIVRE_LIVE_TEST` | harness provider live lane | secret/public | provisional harness-only |
| `MPC_PRODUCT_LINKS_LIVE_TEST`, `MPC_PRODUCT_LINKS_POSTGRES_URL`, `MPC_PRODUCT_LINKS_INSTALLATION_ID` | product-links live test | public/secret/internal | temporary test contract; F-03/F-08 |
| `MPC_TEST_DATABASE_URL` | harness integration input | secret | reserved-not-ambient; F-03 |

Declared Oracle legacy aliases:

- `MPC_ORACLE_USERNAME` <- `SANKHYA_ORACLE_USER`;
- `MPC_ORACLE_PASSWORD` <- `SANKHYA_ORACLE_PASSWORD`;
- `MPC_ORACLE_CONNECT_STRING` <- `SANKHYA_ORACLE_CONNECT_STRING`;
- or composite `SANKHYA_ORACLE_HOST`, `SANKHYA_ORACLE_PORT`,
  `SANKHYA_ORACLE_SERVICE_NAME`.

Do not register hostile fixture keys, generic DB/proxy patterns, host/tool keys,
or Vite built-in `import.meta.env.DEV` as application runtime keys.

Approved typed readers are platform config/pgdb, temporary platform/msdb, and
Oracle config. Vite, Docker, and harness are runtime-edge readers. Provider,
OAuth, Melhor Envio, migration, integration-test, and live-test direct readers
must be exact temporary exceptions; no wildcard reader is allowed.

### Enforceable initial invariants

- registry covers module directories exactly;
- observed module edges are declared and forbidden layers require exact
  exceptions;
- composition-required modules appear in root composition;
- application packages do not import `net/http`, `pgx`, provider SDKs, or
  another module adapter/transport/registry;
- PostgreSQL adapters do not use `database/sql`;
- production Go contains no `panic()` without exact exception;
- OpenAPI and `packages/sdk-runtime` change atomically relative to base SHA;
- frontend source contains no direct `fetch()` outside SDK runtime;
- migration numeric prefixes are unique without exact exception.

### Shared seams

Initial exclusive seams: `api-sdk`, `migration-sequence`, `composition-root`,
`dependency-graph`, `architecture-decisions`, and the narrow provider
capability contract paths. Broad provider/module globs would serialize
unrelated work and are rejected by design.

### PowerShell capability

Installed `pwsh 7.6.2` provides `Test-Json -LiteralPath -SchemaFile` and
`ConvertFrom-Json -AsHashtable`. Draft 2020-12 strict schemas require no third
party dependency. Token count cannot be exact without a tokenizer; F-07 uses
the documented deterministic estimate `ceil(serialized characters / 4)`.

## Recommendation

Build strict schemas first, then semantic drift, then context compilation.
Keep current violations as exact expiring exceptions and fail any growth.
Later owning milestones remove exceptions; F-07 must not refactor product
runtime to make its own checker green.

## Impact On Mission

- F-07 scope is limited to facts with deterministic checks.
- Evidence/run-state/eval schemas remain with F-04/F-05/F-09.
- F-08 replaces environment inheritance; F-03 removes database-test shortcuts;
  M-09/M-10 remove legacy config/import exceptions.

## Handoff

- Current status: Verified and ready for F-07 planning/execution.
- Next owner: F-07 Feature Implementer.
- Next action: Create registries from this inventory and verify every record
  against active source before committing.
- Required files/evidence: F-07 feature/spec/plan, active source scans, schema
  and drift fixtures.
- Blockers or open decisions: None.
