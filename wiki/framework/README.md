# Marketplace Framework

This folder documents the marketplace integration framework used by Marketplace Central.

Use these pages when adding a new marketplace, changing provider metadata, or deciding where a marketplace capability belongs.

## Index

- [Marketplace Integration Framework](marketplace-integration-framework.md)
- [Adding a Marketplace Provider](adding-a-marketplace-provider.md)
- [Auth Strategies](auth-strategies.md)
- [Provider Metadata Contract](provider-metadata-contract.md)
- [Capability Model](capability-model.md)
- [Shopee Fit Analysis](shopee-fit-analysis.md)
- [Vendors Hub](vendors/README.md)
- [Shopee Vendor Playbook](vendors/shopee/README.md)

## Fast Rule

New marketplace work must enter through the framework:

1. `marketplaces` owns business configuration.
2. `integrations` owns provider lifecycle, auth, credentials, installations, and operation state.
3. `connectors` owns external API adapters and capability execution.
4. `packages/feature-marketplaces` renders provider state from `sdk-runtime` only.
