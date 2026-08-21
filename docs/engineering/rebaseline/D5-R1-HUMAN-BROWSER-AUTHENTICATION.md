# D5-R1 — Human Browser Authentication Correction

> **Status:** OPERATOR-APPROVED BOUNDED D5 CORRECTION — executable proof required in the current D6 candidate  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Discovered during:** D6-B2 cross-repository review with MetalDocs  
> **Parent authorities:** accepted D2 identity/access + accepted D5-B2 client/auth, W4 and OpenAPI wire authority  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Approved:** 2026-08-21

## 1. Role and supersession

This document is the smallest targeted correction to the accepted D5 authentication/wire profile after current browser-application security evidence falsified one earlier mechanism choice.

It **supersedes only** the following prior D5 statements where they require the first-party interactive human browser to receive an OAuth access token and call the Product API with the same bearer scheme used by machine clients:

- `D5-B2-PRODUCT-OPERATION-SURFACE.md` §3 B2-A human authentication carrier;
- `D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` only to the extent its first authentication gate assumed one undifferentiated external bearer carrier;
- `D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md` OA-C6 one-`MpcBearerAuth` projection.

All other accepted D0→D5 meaning remains unchanged unless executable proof exposes a material contradiction.

This correction does **not** add a Product operation, Permission, Principal kind, business authority, BFF business API, session business entity, OAuth server, runtime, database, deployment topology or D7 implementation.

## 2. Root cause

The previous D5-B2-A profile selected one external OAuth architecture for both humans and machines:

```text
human browser -> Authorization Code + PKCE -> browser bearer -> MPC Product API
machine       -> Client Credentials          -> bearer         -> MPC Product API
```

That decision correctly avoided browser secrets and MPC-owned credential infrastructure, but it still made the first-party browser JavaScript execution context a bearer-token holder.

Current IETF browser-application guidance materially changes the evidence for this failure class:

- `draft-ietf-oauth-browser-based-apps-27`, dated 2026-07-06 and currently in the RFC Editor queue with intended status **Best Current Practice**, states that a BFF architecture is strongly recommended for business/sensitive/personal-data applications;
- it states that a browser-only OAuth client significantly increases attack surface and is not recommended for those application classes;
- it separately observes that a frontend + API on a common domain can use OpenID Connect for federated login and retain authentication state with a protected cookie session instead of using OAuth access tokens between that frontend and backend.

Primary evidence:

- <https://datatracker.ietf.org/doc/draft-ietf-oauth-browser-based-apps/>
- <https://www.ietf.org/archive/id/draft-ietf-oauth-browser-based-apps-27.html>

Marketplace Central is a first-party internal business application, handles personal/business data and performs consequential writes. There is no current requirement for its human SPA to act as an independent third-party OAuth client to a separately hosted Product resource server.

Therefore the previous **browser bearer** choice is no longer the smallest secure baseline for the current human consumer.

## 3. Corrected authentication profile

Authentication carrier and Product Principal kind remain separate concepts.

### 3.1 Human interactive browser — H

```text
browser
  -> external OIDC login
  -> Authorization Code exchange terminates server-side
  -> external OIDC/OAuth access/refresh tokens remain server-side
  -> MPC issues/maintains one opaque application-session handle
  -> browser sends same-origin Secure HttpOnly session cookie
  -> unsafe Product requests also carry CSRF request-trust proof
  -> server resolves exactly one current MPC human Principal
```

Binding laws:

- architecture remains OIDC/provider-neutral;
- Keycloak remains the preferred first concrete provider candidate;
- stable external binding remains trusted `(issuer, subject)` → one MPC Principal;
- no OIDC access token or refresh token is exposed to ordinary browser JavaScript as the Product API credential;
- the application-session handle is authentication mechanism state, never MPC Permission/business authority;
- baseline cookie profile is `__Host-mpc_session`, `Secure`, `HttpOnly`, `SameSite=Lax`, `Path=/`, no `Domain`;
- unsafe browser methods `POST | PUT | PATCH | DELETE` require an `X-CSRF-Token` request-trust carrier;
- CSRF proof is not authentication, Principal identity, Permission, business disposition or Governance authorization;
- exact session persistence, expiry, rotation, cryptographic representation, CSRF issuance/rotation and OIDC token storage/caching are D7 realization questions.

The server-side security mediation may serve the role commonly called a BFF in OAuth threat-model literature, but MPC does **not** introduce a separate screen-shaped/business BFF service or API. Product operations remain the canonical Product API.

### 3.2 Automation / system — A / S

Machine clients keep the already-accepted external OAuth model:

```text
confidential machine client
  -> Client Credentials
  -> audience-bound bearer access token
  -> MPC Product API
  -> fail-closed resolution to exactly one current non-human Principal
```

Binding laws:

- Client Credentials does not create a Principal kind;
- bearer authentication admits only A/S through this baseline carrier;
- current Principal eligibility, Organization Membership, allowed Principal kind, exact ordinary Permission, physical qualification where required, owner disposition and Governance remain independent current gates;
- no MPC Product API key is introduced;
- token lifetime, client authentication mechanism, sender constraint and Keycloak client/realm/deployment details remain D7 unless a later accepted requirement makes them Product-contract relevant.

## 4. Unchanged authorization sequence

After one admitted authentication carrier succeeds, existing W4 authority remains binding:

```text
1. authenticate admitted carrier
   H -> current application session + required CSRF on unsafe browser request
   A/S -> current audience-bound machine bearer
2. resolve exactly one MPC Principal
3. require current Principal access eligibility
4. identify Product operation
5. resolve current Organization Membership when required
6. enforce allowed Principal kind
7. enforce exact ordinary Permission
8. enforce physical qualification when required
9. resolve resources/secondary references fail-closed in Organization scope
10. apply wire/idempotency/revision/precondition grammar
11. evaluate owner business disposition
12. evaluate Governance when required
13. establish owner intake/effect
```

Consequently:

```text
valid session or bearer
!= current eligibility
!= Membership
!= Permission
!= allowed Principal kind
!= business permitted
!= Governance authorized
!= executable/converged
```

No IdP role, scope, Keycloak role/group or session field becomes MPC ordinary-access authority.

## 5. Canonical OAD projection

The canonical Product OAD represents exactly two alternative authentication schemes:

```yaml
security:
  - MpcHumanSessionAuth: []
  - MpcMachineBearerAuth: []
```

`MpcHumanSessionAuth` is an opaque cookie carrier. `MpcMachineBearerAuth` is the external OAuth bearer carrier.

A bounded top-level `x-mpc-authentication-profile` records the cross-operation client-class projection that OpenAPI security objects alone cannot express:

- H ↔ human session;
- server-side Authorization Code exchange;
- no browser OIDC token exposure;
- required cookie security profile;
- CSRF carrier/method set;
- A/S ↔ Client Credentials bearer.

This extension is projection authority only. It does not create a new business model or allow an operation to override W4 `x-mpc-principal-kinds`.

No Product operation carries an operation-local `security` override in the baseline. The root admits either authentication carrier; current W4 Principal-kind enforcement decides whether the resolved Principal class is legal for the exact operation.

OIDC login/callback/logout protocol endpoints remain Technical Non-Product ingress/runtime surfaces and are not added to the Product OAD operation census.

## 6. Product invariants preserved

The correction must mechanically preserve:

```text
Product operations                  99
ordinary Permissions                30
Principal kinds                     H / A / S only
Performance operations              exact 4 Qs
Performance Permission              performance.read
Product stable origin               https://conexus.fun
Organization scoping                unchanged
Idempotency / ETag grammar           unchanged
Problem Details                     unchanged
Technical Ingress separation        unchanged
Product implementation              BLOCKED UNTIL D9
D7 runtime/deployment                NOT OPEN
```

The accepted 95/29 D5 baseline and the pre-auth 99/30 D6-R1 surface remain executable historical non-regression fixtures. Current generated TypeScript/Go projections must also remain deterministic/compilable under the corrected security profile.

## 7. Negative controls

The corrected proof must fail if any of these become reachable:

1. one universal `MpcBearerAuth` silently returns;
2. the human session becomes browser-readable instead of HttpOnly;
3. an unsafe human method escapes the required CSRF method set;
4. H is admitted through the machine-bearer profile;
5. OIDC access/refresh tokens become the normal browser Product credential again.

These controls supplement, rather than replace, the accepted D5 and D6-R1 semantic/OAD negative controls.

## 8. Deferred realization — D7 only

D7, when opened, owns the mechanism details required to realize this already-decided security profile, including:

- exact server-side ApplicationSession persistence/cryptographic representation;
- session expiry/rotation/revocation propagation;
- CSRF secret generation, carriage/bootstrap and rotation;
- OIDC library and code-exchange mechanics;
- Keycloak realm/client/deployment/HA/backup topology;
- browser/static/API serving topology that preserves the accepted same-origin security property;
- machine token validation/cache/client-auth mechanics.

D7 may not reopen the browser-token isolation property merely for implementation convenience. A materially new human client requiring delegated independent-resource access is a legitimate bounded reopen trigger.

## 9. Continuation law

D6-B1 remains operator-ratified. D6-B2 may resume only after this bounded D5 correction is executable and the current 99/30 Product candidate is green again.

D7–D9 remain blocked. Product implementation remains blocked until D9.
