# Auth Strategies

## Purpose

Auth strategy describes the technical way a tenant connects a provider. It is exposed in provider catalog metadata and drives frontend labels and operational expectations.

Auth strategy does not decide business pricing behavior. Pricing setup remains in `marketplaces`.

## Current Strategies

| Strategy | Install Mode | Meaning | Current Providers |
|---|---|---|---|
| `oauth2` | `interactive` | Standard authorization-code flow with provider token endpoint | Mercado Livre, Magalu |
| `lwa` | `interactive` | Amazon Login With Amazon / SP-API consent flow | Amazon |
| `api_key` | `manual` | User submits API key or equivalent secret manually | Leroy Merlin, current Shopee placeholder |
| `token` | `manual` | User submits provider token manually | Madeira Madeira |
| `none` | any | No provider auth required | none currently |
| `unknown` | any | Auth not yet confirmed | marketplace definitions only, avoid for available providers |

## Interactive Flow Contract

Interactive providers use:

```text
POST /integrations/installations/:id/auth/authorize
```

The service:

1. Loads installation and auth adapter.
2. Creates a signed state and PKCE verifier.
3. Saves `integration_oauth_states`.
4. Calls adapter `StartAuthorize`.
5. Marks installation `pending_connection`.
6. Returns `auth_url`.

Callback handling:

```text
GET/POST callback route, then AuthFlowService.HandleCallback
```

The service:

1. Verifies and consumes signed state.
2. Calls adapter `ExchangeCallback`.
3. Rotates encrypted credential.
4. Upserts auth session.
5. Marks installation `connected/healthy`.

## Manual Flow Contract

Manual providers use:

```text
POST /integrations/installations/:id/auth/credentials
```

The service:

1. Loads installation and auth adapter.
2. Calls adapter `VerifyAPIKey`.
3. Rotates encrypted credential.
4. Marks installation `connected/healthy`.

Manual providers still implement `MarketplaceAuthAdapter`. Unsupported interactive methods return `domain.ErrNotSupported`.

## Refresh Contract

Refresh uses:

```text
AuthFlowService.RefreshCredential
```

The service:

1. Loads auth session and active credential.
2. Decrypts credential payload.
3. Extracts `refresh_token`.
4. Calls adapter `Refresh`.
5. Rotates credential.
6. Clears refresh failure state.
7. Marks installation `connected/healthy`.

The refresh ticker runs from `composition/root.go`.

## Provider-Specific Strategies

Add a provider-specific strategy when existing values would lie about the external API contract.

Example: Shopee Open Platform v2 is not plain OAuth2. It uses partner credentials to sign:

```text
/api/v2/shop/auth_partner
/api/v2/auth/token/get
/api/v2/auth/access_token/get
```

The user experience is interactive, but request signing and token exchange are Shopee-specific. A future implementation should use a first-class strategy such as `shopee_partner` or a more generic `signed_partner` if another provider shares the same pattern.

## Env Variable Pattern

Provider adapters read env variables inside their auth factory:

```go
integrationsproviders.RegisterAuthFactory(func() application.MarketplaceAuthAdapter {
    return NewAdapter(Config{
        ClientID: strings.TrimSpace(os.Getenv("MPC_PROVIDER_<PROVIDER>_CLIENT_ID")),
    })
})
```

Current examples:

```text
MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID
MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET
MPC_PROVIDER_MAGALU_CLIENT_ID
MPC_PROVIDER_MAGALU_CLIENT_SECRET
MPC_PROVIDER_AMAZON_APPLICATION_ID
MPC_PROVIDER_AMAZON_CLIENT_ID
MPC_PROVIDER_AMAZON_CLIENT_SECRET
MPC_PROVIDER_AMAZON_AUTH_VERSION
```

## Guardrails

- Do not label a provider `oauth2` unless the adapter follows standard authorization-code semantics closely enough for maintainers to understand it as OAuth2.
- Do not expose an interactive provider as available until `StartAuthorize` returns a real provider URL.
- Do not store provider secrets outside integration credentials.
- Do not allow reauth to switch external accounts.
- Do not make frontend aware of provider secrets.
- Do not add an auth strategy without updating OpenAPI, SDK, frontend labels, tests, and this document.
