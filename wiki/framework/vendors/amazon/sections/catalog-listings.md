# Amazon SP-API: Catalog and Listings

Last verified: 2026-04-27

## What this section covers

- Catalog lookup and ASIN discovery
- Listing create/update/delete lifecycle
- Product type definitions and listing restrictions
- Listing issue detection and remediation loops

## Representative docs

- [Catalog Items API](https://developer-docs.amazon.com/sp-api/docs/catalog-items-api)
- [Search Catalog Items](https://developer-docs.amazon.com/sp-api/docs/search-catalog-items)
- [Retrieve Catalog Item Details](https://developer-docs.amazon.com/sp-api/docs/retrieve-catalog-item-details)
- [Manage Product Listings Guide](https://developer-docs.amazon.com/sp-api/docs/manage-product-listings-guide)
- [Create a Listing](https://developer-docs.amazon.com/sp-api/docs/create-a-listing)
- [Preview Listing Errors](https://developer-docs.amazon.com/sp-api/docs/preview-errors-before-creating-a-listing)
- [Product Type Definitions API](https://developer-docs.amazon.com/sp-api/docs/product-type-definitions-api)
- [Listings Restrictions API](https://developer-docs.amazon.com/sp-api/docs/listings-restrictions-api)

## MPC notes

- Resolve product type + restrictions before write calls.
- Parse and propagate listing `issues` even when HTTP status is 2xx.

