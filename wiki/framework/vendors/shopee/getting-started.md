# Shopee Getting Started

Last verified: 2026-04-27  
Scope: Shopee Open Platform "Getting Started" docs shared for MPC onboarding.

## Coverage

The following documents were reviewed and mapped:

- `12` Developer account registration
- `14` App management
- `16` API calls
- `18` Push Mechanism notifications
- `20` Authorization and Authentication
- `644` Sandbox Testing V2
- `24` Service Partner Program
- `27` V2.0 API Call Flow
- `28` CNSC API Integration Guide
- `29` KRSC API Integration Guide
- `31` V2.0 Data Definition
- `718` Requesting Access to Sensitive Data

## Practical onboarding sequence for MPC

1. Register and approve a Shopee developer account (`12`).
2. Create app + move app to live as needed (`14`).
3. Implement auth handshake (`20`) and signed requests (`16`).
4. Validate in sandbox before production (`644`).
5. Enable push notifications only after idempotent processors are ready (`18`).
6. Align integration flows with V2 call/data definitions (`27`, `31`).
7. For cross-border scenarios, follow CNSC/KRSC specifics (`28`, `29`).
8. If sensitive scopes are needed, complete eligibility + security controls (`718`).

## Core technical rules

- All API requests are signed and time-sensitive (`16`, `20`).
- Access token lifecycle must include refresh and rotation (`20`).
- Integration design must support both shop-level and merchant/main-account authorization paths (`20`).
- Push channels are category-specific (order, product, marketing, chat) and require subscription lifecycle management (`18`).
- Data structures and status enums must track Shopee V2 definitions to avoid mapping drift (`31`).
- Sensitive-data scopes require additional governance requirements (security evidence, allowlist posture, eligibility) (`718`).

## MPC implementation notes

- Keep partner credentials server-side only.
- Current backend configuration keys:
  - `MPC_PROVIDER_SHOPEE_PARTNER_ID`
  - `MPC_PROVIDER_SHOPEE_PARTNER_KEY`
  - `MPC_PROVIDER_SHOPEE_BASE_URL`
- Normalize Shopee statuses into MPC internal enums at the connector boundary.
- Treat push events as eventually consistent: reconcile with pull APIs for critical workflows.
- Store authorization context (`shop_id`/merchant context) with token material for deterministic refresh.
- Keep signature generation deterministic and testable with fixed timestamps.
- Callback mapping currently accepts `shop_id`, `merchant_id`, or `selling_partner_id` as provider account identity input.
