# Milestone Validation Result

```yaml
id: M-02
type: milestone-validation-result
status: passed
owner: QA Validator
parent: MIS-001
created: 2026-07-06
updated: 2026-07-08
validation_level: QA-0
lifecycle_scope: milestone
```

## Milestone

M-02-marketplace-capability-framework

## Verdict

- Result: `passed`
- Blocking failures: none
- Summary: Marketplace Central exposes provider-agnostic marketplace capability contracts, aligned provider/runtime vocabulary, and a live-wired Mercado Livre read/probe surface with real account, listing, order, fee-quote, and stock-read evidence through the platform installation endpoints.

## Validation Scope Declaration

- contract_validated: Yes
- integration_validated: Yes, for read/probe flows only
- blocked_for_real_validation: provider writes remain intentionally unvalidated and out of runtime scope for this milestone

This pass now covers both capability-contract readiness and live Mercado Livre runtime validation for read/probe flows only. It does not claim write validation, shipment validation, messaging validation, or any unsafe provider-side mutation.

## Feature Evidence

- F-01 capability port contract: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/F-01-capability-port-contract/validation.md`
- F-02 provider capability registration: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/F-02-provider-capability-registration/validation.md`
- F-03 Mercado Livre adapter spine: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-02-marketplace-capability-framework/F-03-mercado-livre-adapter-spine/validation.md`

## Criterion Review

### M-02-C01 — Business Modules Can Depend On Capability Ports

- Status: `passed`
- Commands:
  - `go test ./...` from `apps/server_core/internal/modules/connectors`
  - `rg -n "modules/connectors/adapters/mercado_livre|modules/integrations/adapters/mercadolivre" apps/server_core/internal/modules/connectors/application apps/server_core/internal/modules/connectors/ports apps/server_core/internal/modules/connectors/domain`
- Expected:
  - normalized capability contracts compile and no Mercado Livre adapter imports leak into the business-facing connector surface.
- Actual:
  - connector module tests passed during F-01 validation
  - import-boundary search returned no matches
- Blocking failure observed:
  - `No`

### M-02-C02 — Provider And Public Capability Vocabulary Is Aligned

- Status: `passed`
- Commands:
  - `go test ./internal/modules/marketplaces/... ./internal/modules/integrations/...`
  - `npm run test --workspace @marketplace-central/sdk-runtime -- --run src/index.test.ts`
- Expected:
  - provider declarations, marketplace profiles, OpenAPI, and SDK all use `listing_read`, `stock_read`, `stock_write`, and `order_read` with conservative future-provider support.
- Actual:
  - backend marketplace/integrations tests passed
  - SDK runtime test passed (`16` tests)
  - Mercado Livre is the only provider advertising the new operational business capabilities
- Blocking failure observed:
  - `No`

### M-02-C03 — Mercado Livre Maps Official Listing, Stock, And Order Shapes

- Status: `passed`
- Commands:
  - `go test ./internal/modules/connectors/adapters/mercado_livre`
  - `go test ./internal/modules/connectors/application ./internal/modules/connectors/adapters/mercado_livre`
- Expected:
  - Mercado Livre adapter maps item-without-variation, item-with-variation, and order-with-sale-fee shapes, while explicit provider rejection and `429` behavior are covered by tests.
- Actual:
  - adapter package tests passed
  - shared capability service plus adapter package passed together
  - tests cover item mapping, variation stock mapping, order fee/payment mapping, provider rejection, rate-limit mapping, and unsupported variation shape
- Blocking failure observed:
  - `No`

### M-02-C04 — Installation Runtime Surface Is Honest And Live-Validated

- Status: `passed`
- Commands:
  - `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/integrations/... -count=1`
  - `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/connectors/... -count=1`
  - `npm test -- --run packages/sdk-runtime/src/index.test.ts packages/feature-integrations/src/IntegrationsHubPage.test.tsx`
  - `npm run build --if-present`
  - `Invoke-WebRequest -UseBasicParsing 'http://localhost:8080/integrations/installations'`
  - `Invoke-WebRequest -UseBasicParsing 'http://localhost:8080/integrations/installations/inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98/auth/status'`
  - `Invoke-WebRequest -UseBasicParsing 'http://localhost:8080/integrations/installations/inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98/probes/account' -Method Post`
  - `Invoke-WebRequest -UseBasicParsing 'http://localhost:8080/integrations/installations/inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98/probes/listings?limit=1'`
  - `Invoke-WebRequest -UseBasicParsing 'http://localhost:8080/integrations/installations/inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98/probes/orders?limit=1'`
  - `Invoke-WebRequest -UseBasicParsing 'http://localhost:8080/integrations/installations/inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98/probes/fee-quote?listing_type_id=gold_special&price=199.9&currency_id=BRL&site_id=MLB&category_id=MLB1000'`
  - `Invoke-WebRequest -UseBasicParsing 'http://localhost:8080/integrations/installations/inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98/probes/stock?provider_item_id=MLB4807275656'`
- Expected:
  - backend, SDK, and integrations UI remain aligned around connection snapshot, runtime capabilities, operation evidence, and read/probe endpoints
  - Mercado Livre installation returns connected/healthy connection truth
  - runtime capabilities become live-validated only after real provider execution
  - no write capability is exposed as runtime-available proof
- Actual:
  - Go integrations suite passed
  - Go connectors suite passed
  - SDK + integrations UI tests passed (`44` tests across the two files)
  - frontend production build passed
  - installation list now returns `connection`, `runtime_capabilities`, and `stock_read.live_validated = true`
  - auth status mirrors the same connected/healthy Mercado Livre account snapshot for installation `inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98`
  - live probes succeeded for `account_probe`, `listing_read`, `order_read`, `fee_quote_read`, and `stock_read`
  - operation timeline persists translated evidence such as `available_quantity: 11 | provider_item_id: MLB4807275656 | provider_status: paused | scope: item`
- Blocking failure observed:
  - `No`

## Validation Notes

- Official Mercado Livre docs were re-checked through Context7 during F-03 and used as the seam source for `GET /items/{ITEM_ID}`, `PUT /items/{ITEM_ID}`, `GET /orders/{ORDER_ID}`, `GET /orders/search?seller=...`, and `GET /users/{USER_ID}/items/search`.
- A broad `go test ./internal/modules/connectors/...` run stalled in this environment after partial output, so milestone evidence uses narrower load-bearing package runs that completed successfully instead of overstating a full-module pass.
- Live validation on 2026-07-08 used the local Docker stack (`backend`, `frontend`, `ngrok`, `postgres`) with the active Mercado Livre installation and captured successful platform-level reads through `/integrations/installations/{id}/probes/*`.
- In-app browser QA on 2026-07-08 confirmed the live integrations page renders `METALNOBREACABAMENTOS`, `connected`, `healthy`, `next action = none`, runtime capabilities (`account_probe`, `listing_read`, `order_read`, `fee_quote_read`, `stock_read`), and the persisted operations timeline without a misleading `Not connected` card state.
- Live provider-write validation remains intentionally outside M-02 because the milestone brief forbids uncontrolled real provider writes.

## Live Mercado Livre Read Validation

Date: 2026-07-08
Environment: local Docker backend/frontend with active development credentials
Installation: `inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98`

Evidence:
- Installation list returned `connection.state = connected`, `connection.health = healthy`, `connection.external_account_id = 691607102`, `connection.external_account_name = METALNOBREACABAMENTOS`, and runtime capabilities with `stock_read.live_validated = true`.
- Auth status mirrored the same connection snapshot with `next_action = none`.
- `POST /probes/account` returned the live Mercado Livre account identity for seller `691607102`.
- `GET /probes/listings?limit=1` succeeded and persisted `listing_read`.
- `GET /probes/orders?limit=1` succeeded and persisted `order_read`.
- `GET /probes/fee-quote?...` succeeded and persisted `fee_quote_read` without defaulting unknown values to zero.
- `GET /probes/stock?provider_item_id=MLB4807275656` returned `available_quantity = 11`, `provider_status = paused`, `ean = 7898016503119`, `title = Placa De Sinalização Aviso Aos Passageiros 18x15cm-sinalize`, `scope = item`, and persisted `stock_read`.
- `GET /operations` now shows successful `stock_read`, `order_read`, `listing_read`, `fee_quote_read`, and `account_probe` runs with structured `provider_evidence`.

Boundary:
- No provider writes were executed.
- Mock/fake tests were used only for local contract and deterministic behavior.
- Live validation covers Mercado Livre read/probe behavior only.

## Handoff

- Milestone status: `ready for mission continuation`
- Next recommended action: move into the next milestone that consumes the now live-validated read/probe foundation for product-link, stock safety, and order-margin workflows
- Open blockers: none for M-02
