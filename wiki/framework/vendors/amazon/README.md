# Amazon SP-API Vendor Playbook

Last verified: 2026-04-27

## Why this exists

This is the canonical Amazon SP-API guide for Marketplace Central implementation work.
It complements runtime notes in `docs/marketplaces/amazon.md` and converts the 454-page SP-API map into LLM-oriented sections.

## Current fit for MPC

- Auth model: Login With Amazon (LWA) OAuth + AWS SigV4 signed API requests.
- Integration shape: `integrations` owns install/auth/credential lifecycle, `connectors` owns SP-API execution.
- Current provider runtime status:
  - `auth_strategy`: `oauth2`
  - `install_mode`: `interactive`
  - `rollout_stage`: `v1`
  - `execution_mode`: `live`
- Primary capability groups:
  - onboarding/auth
  - catalog/listings
  - price/fees + inventory
  - orders/fulfillment
  - notifications
  - reports/analytics
  - finance/payments
  - appstore/compliance
  - vendor-specific APIs

## Reading order

1. [Getting Started](getting-started.md)
2. [API Best Practices](api-best-practices.md)
3. [Capability Matrix](capability-matrix.md)
4. [Section Index](section-index.md)
5. [Sources and Coverage](sources.md)
6. [Implementation Sync Runbook](implementation-sync.md)

## LLM retrieval helpers

- Full extracted site map:
  - `docs/marketplaces/amazon-spapi-docs-map.json`
- Includes 454 SP-API pages with `url`, `title`, and `description`.
