# Amazon SP-API: Onboarding and Auth

Last verified: 2026-04-27

## What this section covers

- Developer/service provider onboarding
- App registration and authorization models
- OAuth callback flow and token lifecycle
- Signed request requirements

## Representative docs

- [Welcome](https://developer-docs.amazon.com/sp-api/docs/welcome)
- [Onboarding Overview](https://developer-docs.amazon.com/sp-api/docs/onboarding-overview)
- [Service Provider Onboarding Overview](https://developer-docs.amazon.com/sp-api/docs/service-provider-onboarding-overview)
- [Authorizing SP-API Applications](https://developer-docs.amazon.com/sp-api/docs/authorizing-selling-partner-api-applications)
- [Website Authorization Workflow](https://developer-docs.amazon.com/sp-api/docs/website-authorization-workflow)
- [Selling Partner Appstore Authorization Workflow](https://developer-docs.amazon.com/sp-api/docs/selling-partner-appstore-authorization-workflow)
- [Connecting to the SP-API](https://developer-docs.amazon.com/sp-api/docs/connecting-to-the-selling-partner-api)
- [Configuration Details](https://developer-docs.amazon.com/sp-api/docs/configuration-details)

## MPC notes

- Keep auth transitions in `integrations`.
- Keep request signing in `connectors`.
- Do not mix app identifiers between consent and LWA token exchange contexts.

