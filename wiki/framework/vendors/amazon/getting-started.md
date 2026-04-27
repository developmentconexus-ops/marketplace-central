# Amazon SP-API Getting Started

Last verified: 2026-04-27  
Scope: Amazon Selling Partner API onboarding and authorization flows required by MPC.

## Practical onboarding sequence for MPC

1. Register in Solution Provider Portal and choose app type (private/public).
2. Create SP-API app credentials and Login With Amazon credentials.
3. Configure website authorization redirect URI for MPC callback.
4. Implement OAuth state validation and auth code exchange.
5. Persist refresh token securely and implement token refresh automation.
6. Implement AWS SigV4 request signing for every SP-API call.
7. Validate first production-safe calls (catalog/orders read paths).
8. Add restricted-data (RDT) flow only for PII endpoints.
9. Configure notifications (SQS or EventBridge) with idempotent consumers.
10. Run sandbox and pre-production evidence checks before seller rollout.

## Core technical rules

- SP-API requests require both:
  - LWA access token (`x-amz-access-token`), and
  - AWS SigV4 signature.
- Restricted operations use Tokens API (RDT) instead of plain LWA token.
- OAuth callback must validate `state` and reject mismatches.
- Treat `amzn1.sellerapps.app.*` as consent app identifier, not `amzn1.sp.solution.*`.
- Token lifecycle:
  - short-lived access token,
  - long-lived refresh token,
  - secret rotation cadence enforced.

## MPC implementation notes

- `integrations`:
  - installation draft -> connected transitions
  - OAuth initiation/callback handling
  - credential encryption and rotation metadata
- `connectors`:
  - SigV4 signing wrapper
  - RDT-aware client paths
  - endpoint throttling/backoff execution
- Keep write operations idempotent and retry-safe with deterministic operation keys.

## First-call validation checklist

1. Authorization URL generated correctly (includes application id + callback + state).
2. Callback exchanges code for refresh token successfully.
3. Refresh path can mint new access token.
4. Signed API request succeeds on a low-risk endpoint (`getOrders` or catalog read).
5. Errors map into MPC structured codes and logs include action/result/duration.

