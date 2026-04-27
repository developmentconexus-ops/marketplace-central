# Shopee Vendor Playbook

Last verified: 2026-04-27

## Why this exists

This is the canonical Shopee operational guide for Marketplace Central implementation work.
It complements (does not replace) `wiki/framework/shopee-fit-analysis.md`.

## Current fit for MPC

- Auth model: signed partner authorization flow (`/api/v2/shop/auth_partner`, token exchange and refresh)
- Integration shape: `integrations` owns install/auth/credential lifecycle, `connectors` owns API execution
- Initial capability priority:
  - auth + installation lifecycle
  - product + stock/price flows
  - orders + logistics + returns
  - push subscriptions for reactive updates

## Reading order

1. [Getting Started](getting-started.md)
2. [API Best Practices](api-best-practices.md)
3. [Capability Matrix](capability-matrix.md)
4. [Sources and Coverage](sources.md)
5. [Shopee Fit Analysis (framework mapping)](../../shopee-fit-analysis.md)

## LLM retrieval helpers

- Full extracted document index:
  - `docs/marketplaces/shopee-openplatform-docs-index.json`
- Contains 42 Shopee docs (Getting Started + API Best Practices) with:
  - `id`, `title`, `url`, `headings`, detected endpoints, and timestamps.
