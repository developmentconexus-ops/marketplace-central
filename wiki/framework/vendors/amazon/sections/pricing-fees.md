# Amazon SP-API: Pricing and Fees

Last verified: 2026-04-27

## What this section covers

- Product Fees and pricing APIs
- Feed pathways for bulk updates
- Rate-limit-aware pricing operations

## Representative docs

- [Product Fees API](https://developer-docs.amazon.com/sp-api/docs/product-fees-api)
- [Get Product Fee Estimates (ASIN)](https://developer-docs.amazon.com/sp-api/docs/get-product-fee-estimates-asin)
- [Get Product Fees (Batch)](https://developer-docs.amazon.com/sp-api/docs/get-product-fees-batch)
- [Product Pricing API](https://developer-docs.amazon.com/sp-api/docs/product-pricing-api)
- [Manage Automated Pricing Rules](https://developer-docs.amazon.com/sp-api/docs/manage-automated-pricing-rules)
- [Price Adjustment Automation Workflows Guide](https://developer-docs.amazon.com/sp-api/docs/price-adjustment-automation-workflows-guide)
- [Feeds API Best Practices](https://developer-docs.amazon.com/sp-api/docs/feeds-api-best-practices)
- [Feeds API Rate Limits](https://developer-docs.amazon.com/sp-api/docs/feeds-api-rate-limits)

## MPC notes

- Keep pricing decisions in `pricing`, execution in `connectors`.
- Preserve fee estimation responses for audit and simulation explainability.

