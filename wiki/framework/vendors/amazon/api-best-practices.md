# Amazon SP-API Best Practices

Last verified: 2026-04-27  
Scope: SP-API best-practice guidance mapped to MPC implementation priorities.

## Best-practice contract for MPC

### 1) Auth, identity, and secrets

- Keep app identifiers, LWA credentials, refresh tokens, and IAM credentials clearly separated.
- Enforce credential rotation windows (especially LWA client secret cycle).
- Never log secret material, raw tokens, or full signed headers.
- Use dedicated auth clients for:
  - standard SP-API calls,
  - restricted data token (RDT) flows.

### 2) Throttling and retries

- Treat rate limits as baseline behavior, not exceptional failures.
- Read `x-amzn-RateLimit-Limit` when present.
- On `429`/transient errors, retry with capped exponential backoff + jitter.
- Keep endpoint-specific budgets to prevent one capability starving others.

### 3) Idempotency and mutation safety

- Every write path must be retry-safe (`listings`, shipping updates, finance mutations).
- Use stable operation identifiers per tenant/account/entity change.
- Persist request/response correlation IDs for replay and audit.

### 4) Notifications before aggressive polling

- Prefer Notifications API transport (SQS/EventBridge) for near-real-time updates.
- Assume duplicates and out-of-order delivery for SQS standard queues.
- Dedupe with stable event key (`notificationId`) and keep replay windows.
- Use polling as fallback reconciliation, not only source of truth.

### 5) PII and restricted operations

- Gate restricted endpoints behind explicit capability flags and role validation.
- Acquire RDT only when required and keep short-lived token scope isolated.
- Do not reuse RDT token channels for non-restricted calls.

### 6) Listings, pricing, and inventory

- Resolve product type requirements and listing restrictions before write calls.
- Surface listing issues and suppression states as first-class connector results.
- Separate business decisioning (pricing module) from transport execution (connector).
- Prefer batch or feed patterns only when they reduce throttling pressure safely.

### 7) Orders, fulfillment, and finance

- Maintain explicit order/fulfillment status mapping tables in connector boundaries.
- Use robust pagination and NextToken handling with checkpoints.
- Keep invoice and payment data flows auditable and compliant by market.

## Recommended module ownership in MPC

- `integrations`
  - OAuth/auth session lifecycle
  - credential and rotation lifecycle
  - provider capability gating
- `connectors`
  - SigV4 + LWA/RDT execution
  - rate-limit aware request orchestration
  - notifications ingestion adapter contracts
- domain modules (`catalog`, `orders`, `pricing`, `alerts`)
  - canonical business decisions and SLA policies

## Production readiness checklist

1. OAuth + callback + refresh flow verified in runtime environment.
2. SigV4 request layer covered with deterministic tests.
3. Throttle/retry policy tested with synthetic 429 scenarios.
4. Notification ingestion includes dedupe + retry-safe handlers.
5. Restricted-data paths isolated and tested separately.
6. Structured logging present for every action with duration and result.

