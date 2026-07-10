# Feature Validation

```yaml
id: F-02
type: feature-validation
status: in_progress
owner: Feature Implementer
parent: F-02
created: 2026-07-08
updated: 2026-07-08
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-02-link-candidate-engine

## Summary

`product_links` now generates and persists exact-first link candidates from imported listing snapshots, exposes manual generation/listing APIs, and degrades honestly when the Oracle read boundary is unavailable.

## Current Validation State

- Result: Passed for local contract behavior, degraded-runtime behavior, and positive live Oracle/Postgres generation
- Result owner: Feature Implementer
- Decision date: 2026-07-08
- Final feature state for handoff: ready_for_resolution_workflow

## Evidence

- Command: `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/product_links/... ./apps/server_core/internal/composition -count=1`
  - Result: Pass
- Command: `npm run test --workspace @marketplace-central/sdk-runtime -- src/index.test.ts`
  - Result: Pass
- Command: `Invoke-WebRequest -Uri 'http://localhost:8080/product-links/link-candidates/generations' -Method Post -ContentType 'application/json' -Body '{"installation_id":"inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98","limit":5}' -SkipHttpErrorCheck`
  - Result: Pass for degraded-runtime proof
  - Observed: `503` with `PRODUCT_LINKS_INTERNAL_READ_UNAVAILABLE`
- Command: `docker compose exec -T postgres psql -U marketplace -d marketplace_central -c "SELECT COUNT(*) FROM product_link_candidates WHERE tenant_id = 'tenant_default' AND installation_id = 'inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98';"`
  - Result: Pass
  - Observed: `0` rows after blocked generation, consistent with fail-fast behavior
- Command: `Test-NetConnection -ComputerName 10.55.10.101 -Port 1521`
  - Result: Pass
  - Observed: `TcpTestSucceeded : True`
- Command: `$env:PATH='C:\Users\leandro.theodoro\AppData\Local\codex-tools\winlibs-mcf-ucrt\bin;C:\oracle\instantclient_23_0;' + $env:PATH; $env:CC='C:\Users\leandro.theodoro\AppData\Local\codex-tools\winlibs-mcf-ucrt\bin\gcc.exe'; $env:CXX='C:\Users\leandro.theodoro\AppData\Local\codex-tools\winlibs-mcf-ucrt\bin\g++.exe'; $env:CGO_ENABLED='1'; $env:SANKHYA_ORACLE_HOST='10.55.10.101'; $env:SANKHYA_ORACLE_PORT='1521'; $env:SANKHYA_ORACLE_SERVICE_NAME='ORCL'; $env:SANKHYA_ORACLE_USER='leandro'; $env:SANKHYA_ORACLE_PASSWORD='troca#123'; $env:MPC_ORACLE_USERNAME=$env:SANKHYA_ORACLE_USER; $env:MPC_ORACLE_PASSWORD=$env:SANKHYA_ORACLE_PASSWORD; $env:MPC_ORACLE_CONNECT_STRING="$($env:SANKHYA_ORACLE_HOST):$($env:SANKHYA_ORACLE_PORT)/$($env:SANKHYA_ORACLE_SERVICE_NAME)"; $env:MPC_ORACLE_LIB_DIR='C:/oracle/instantclient_23_0'; $env:MPC_ORACLE_LIVE_TEST='1'; $env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/internal_read/adapters/oracle -run TestOracleLiveSmoke -v -count=1`
  - Result: Pass
  - Observed: live smoke passed for product lookup, stock, current price, cost, sales history, and tax inputs
- Command: `$env:PATH='C:\Users\leandro.theodoro\AppData\Local\codex-tools\winlibs-mcf-ucrt\bin;C:\oracle\instantclient_23_0;' + $env:PATH; $env:CC='C:\Users\leandro.theodoro\AppData\Local\codex-tools\winlibs-mcf-ucrt\bin\gcc.exe'; $env:CXX='C:\Users\leandro.theodoro\AppData\Local\codex-tools\winlibs-mcf-ucrt\bin\g++.exe'; $env:CGO_ENABLED='1'; $env:SANKHYA_ORACLE_HOST='10.55.10.101'; $env:SANKHYA_ORACLE_PORT='1521'; $env:SANKHYA_ORACLE_SERVICE_NAME='ORCL'; $env:SANKHYA_ORACLE_USER='leandro'; $env:SANKHYA_ORACLE_PASSWORD='troca#123'; $env:MPC_ORACLE_USERNAME=$env:SANKHYA_ORACLE_USER; $env:MPC_ORACLE_PASSWORD=$env:SANKHYA_ORACLE_PASSWORD; $env:MPC_ORACLE_CONNECT_STRING="$($env:SANKHYA_ORACLE_HOST):$($env:SANKHYA_ORACLE_PORT)/$($env:SANKHYA_ORACLE_SERVICE_NAME)"; $env:MPC_ORACLE_LIB_DIR='C:/oracle/instantclient_23_0'; $env:MPC_PRODUCT_LINKS_LIVE_TEST='1'; $env:MPC_PRODUCT_LINKS_INSTALLATION_ID='inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98'; $env:MPC_PRODUCT_LINKS_POSTGRES_URL='postgres://marketplace:marketplace@127.0.0.1:5435/marketplace_central?sslmode=disable'; $env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/product_links/application -run TestGenerateLinkCandidatesLive -v -count=1`
  - Result: Pass
  - Observed: `generated_count=5 persisted_count=5 states=map[exact_ean:5]`
- Command: `docker exec marketplace-central-postgres-1 psql -U marketplace -d marketplace_central -c "SELECT state, COUNT(*) FROM product_link_candidates WHERE tenant_id = 'tenant_default' AND installation_id = 'inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98' GROUP BY state ORDER BY state;"`
  - Result: Pass
  - Observed: `exact_ean | 5`

## Observed

- Exact SKU, exact EAN, exact conflict, title fallback, and unresolved flows are covered in focused unit tests.
- Manual generation/listing contracts exist in Go transport, OpenAPI, and `sdk-runtime`.
- The runtime now returns a product-links-specific `503` when Oracle is unavailable instead of pretending generation succeeded.
- Positive live generation now completes against real Oracle data and persists candidates in Postgres with exact-first classification.
- Live validation still covers reads and candidate persistence only; it does not imply any stock write path is approved or available.

## Scope Declaration

- contract_validated: Yes
- local_business_validation: Yes
- degraded_runtime_validation: Yes
- positive_live_oracle_validation: Yes
- blocked_for_real_validation: No

## Handoff

- Current status: `passed`
- Next owner: Feature Implementer
- Next action: build F-03 resolution workflow on top of persisted candidate evidence and exact-first states
- Required files/evidence: API/SDK/UI resolution flow, audit persistence, and operator-visible conflict/unresolved handling
- Blockers or open decisions: none for the candidate-generation slice
