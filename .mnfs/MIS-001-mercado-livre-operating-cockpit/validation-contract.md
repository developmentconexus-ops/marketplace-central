# Mission Validation Contract

```yaml
id: MIS-001
type: mission-validation-contract
status: draft
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: mission
```

## Mission ID

MIS-001

## QA Level

QA-0 planning contract; milestone execution must raise evidence level with tests/builds/browser checks.

## Required Final State

- VTEX is removed from target runtime surfaces.
- Marketplace operations are capability-based and Mercado Livre is the first adapter.
- MPC reads internal facts through MNOS/Sankhya contracts and stores only MPC-owned state/snapshots.
- Existing Mercado Livre listings can be linked to internal products.
- Stock Seguro shows safe stock divergence and supports manual audited stock actions.
- Orders + Margin creates margin-quality visibility without treating missing costs/freight/tax as zero.

## Criteria

## Criterion: Project Conventions Preserved
ID: MIS-001-C01
Level: Mission
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `GOCACHE=.gocache go test ./...` from `apps/server_core` and `npm test -- --run`
- Expected: Go tests exit 0; frontend tests exit 0; no frontend feature package calls backend outside `packages/sdk-runtime`.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/validation-result.md`
Blocking failure: Tests fail or frontend direct backend fetch exists.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Module Boundaries Preserved
ID: MIS-001-C02
Level: Mission
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: targeted `rg` checks plus Go tests
- Expected: `inventory`, `product_links`, `orders`, and `profitability` do not import provider HTTP adapter packages, `net/http`, `pgx`, or another module's internals from application/domain layers.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/validation-result.md`
Blocking failure: Business logic depends on provider adapter internals or transport/infrastructure types.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Idempotent Ingestion And Actions
ID: MIS-001-C03
Level: Mission
Type: Reliability
Required: Yes
Status: Pending
Evidence:
- Command: module tests for listing refresh, order refresh, and stock action apply
- Expected: Reprocessing the same provider listing/order/action id produces one durable business record or one duplicate-safe update, not duplicate actions/orders.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/validation-result.md`
Blocking failure: Duplicate provider resource processing creates duplicate business actions/orders.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Provider Write Audit
ID: MIS-001-C04
Level: Mission
Type: Observability
Required: Yes
Status: Pending
Evidence:
- Command: stock action tests and audit repository inspection
- Expected: Applied or failed Mercado Livre stock write stores action id, operator/manual trigger, link id, before quantity, requested quantity, policy id, source timestamps, provider response summary, `action`, `result`, and `duration_ms`.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/validation-result.md`
Blocking failure: Provider write can occur without audit record or source/policy evidence.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Secret Handling
ID: MIS-001-C05
Level: Mission
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: tests for integrations/Sankhya error paths plus `rg` over logs/UI response models
- Expected: Access tokens, refresh tokens, app secrets, Sankhya password, and Oracle DSN password do not appear in logs, API responses, UI models, or validation artifacts.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/validation-result.md`
Blocking failure: Any secret value is returned or logged.
Blocking failure observed: No
Owner: QA Validator

## Criterion: UI Operational States
ID: MIS-001-C06
Level: Mission
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: frontend tests and browser validation for product links, stock risks, and margin views
- Expected: Each data-fetching screen visibly handles loading, error, empty, blocked, conflict, stale, and success states using SDK runtime data.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/validation-result.md`
Blocking failure: Main operator workflow lacks a visible blocked/conflict/error state.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Stock Dashboard Responsiveness
ID: MIS-001-C07
Level: Mission
Type: Performance
Required: Yes
Status: Pending
Evidence:
- Command: repository/service benchmark or test fixture with 100 linked listings
- Expected: Stock risk list calculation returns 100 rows with link state, internal stock, ML stock, risk, recommended quantity, and source timestamps in p95 < 500ms in local test environment.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/validation-result.md`
Blocking failure: Stock risk view depends on synchronous live provider calls or exceeds target in seeded test.
Blocking failure observed: No
Owner: QA Validator

## Criterion: MNOS/Sankhya Semantics Preserved
ID: MIS-001-C08
Level: Mission
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: Sankhya read contract tests using seeded/fake rows
- Expected: Sellable stock uses `SUM(ESTOQUE - RESERVADO)` for `CODEMP IN (1,2)` and `CODLOCAL=10101`; showroom `CODLOCAL=10108` is excluded; missing cost yields `missing_cost`.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/validation-result.md`
Blocking failure: Showroom stock contributes to default sellable stock or missing cost is treated as zero.
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Every milestone must produce `validation-result.md`.
- Every feature must produce `validation.md`.
- API changes must include OpenAPI and `packages/sdk-runtime` evidence together.
- UI changes must include frontend tests and at least one browser validation pass before claiming user workflow complete.
- Provider behavior claims must cite official Mercado Livre docs or local adapter tests using documented shapes.
- Fake/mock/seam evidence may prove contract and business-rule readiness only; it does not prove real dependency integration.
- Any feature or milestone that claims Sankhya, Mercado Livre, OAuth, DB, or other external/runtime-dependent behavior works must include real-environment validation evidence or stay blocked/not passed for that scope.
- Validation artifacts must label whether evidence is fake/local or real-environment.

## Blocking Failures

- Any Mercado Livre write without a resolved product link, visible policy, and audit evidence.
- Any direct React provider/backend bypass outside SDK runtime.
- Any production code path writing to Sankhya/Oracle.
- Any missing value used as zero in stock or margin calculations.
- Any VTEX route/SDK/UI remaining as active target architecture after M-01.

## Retry Policy

- Milestone correction attempts: max 2 before escalating scope.
- Provider-doc ambiguity blocks the feature until official docs or operator-provided docs resolve it.

## Handoff

- Current status: Mission validation contract drafted.
- Next owner: Mission reviewer and QA Validator.
- Next action: Run readiness review, then execute M-01.
- Required files/evidence: Mission and milestone validation artifacts.
- Blockers or open decisions: Independent readiness review pending.
