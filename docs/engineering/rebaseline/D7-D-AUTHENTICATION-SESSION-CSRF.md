# D7-D — Authentication / Session / CSRF / Machine Token Realization

> **Status:** CANDIDATE / OPERATOR RATIFICATION PENDING  
> **Parent:** `D7-RUNTIME-JOBS-TRANSACTIONS.md`  
> **Accepted prerequisites:** D7-A + D7-B + D7-C — OPERATOR-RATIFIED  
> **Parent authority:** accepted D5-R1 Human Browser Authentication Correction + accepted D2 identity/access semantics  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Derived:** 2026-08-22

## 1. Purpose

D7-D realizes the already-accepted authentication carrier split without creating a second authorization authority:

```text
H browser  -> server-side OIDC login -> opaque MPC application session + CSRF
A / S      -> Client Credentials -> audience-bound machine bearer
```

D7-D selects the minimum provider/library/session mechanics needed to make that split executable. It does not add Product operations, Permissions, Principal kinds, IdP roles-as-Permissions, a user-profile domain, an MPC password store, or a browser bearer path.

## 2. Target invariant

> **The browser never owns an OIDC access/refresh token; a human request is authenticated only by a current server-side MPC session and independently trusted for unsafe methods by CSRF controls; a machine bearer is accepted only after cryptographic issuer/audience/time validation and a unique current binding to an A/S Principal; neither carrier grants Membership, Permission, business disposition or Governance authority.**

## 3. First provider profile

### 3.1 Keycloak — SELECT CANDIDATE

Keycloak is selected as the first concrete OIDC/OAuth provider while the application boundary remains standards-based and provider-neutral.

One configured Keycloak issuer/realm is sufficient for the MPC baseline. Exact realm naming, Keycloak version pin, HA topology, backup and upgrade mechanics remain D7-E/deployment concerns.

Use separate client roles:

```text
human web client
  -> confidential OIDC client
  -> Standard Flow / Authorization Code enabled
  -> PKCE S256 required
  -> Implicit Flow disabled
  -> Direct Access Grants disabled
  -> Service Accounts disabled

machine client(s)
  -> confidential OAuth client
  -> Service Account / Client Credentials enabled
  -> Standard Flow disabled
  -> Implicit Flow disabled
  -> Direct Access Grants disabled
  -> access token audience includes https://conexus.fun
```

Keycloak roles/scopes/service-account roles are credential-issuance configuration only. They never map directly to MPC ordinary Permissions or domain authorization.

Production `Full Scope Allowed` is not baseline; machine clients expose only the minimum provider-side token claims required for authentication and audience restriction.

### 3.2 Current evidence

Revalidated on 2026-08-22:

- RFC 9700 recommends PKCE even for confidential Authorization Code clients and requires S256-class protection rather than exposing the verifier;
- current Keycloak documentation supports confidential server-side clients, Standard Flow, required PKCE S256, service accounts, Client Credentials and explicit audience mappers;
- the current IETF browser-application BCP remains in the RFC Editor queue and continues to strongly favor server-side/BFF-class token isolation for business/sensitive browser applications.

Primary sources:

- <https://www.rfc-editor.org/rfc/rfc9700>
- <https://datatracker.ietf.org/doc/draft-ietf-oauth-browser-based-apps/>
- <https://www.keycloak.org/docs/latest/server_admin/>

## 4. Human OIDC login realization

### 4.1 Libraries

Select:

```text
github.com/coreos/go-oidc/v3/oidc
golang.org/x/oauth2
```

Reasons:

- provider discovery and ID-token signature/issuer/audience validation are standard OIDC concerns rather than MPC business code;
- `x/oauth2` provides per-authorization PKCE verifier generation and S256 challenge/exchange options;
- `go-oidc` requires nonce validation by the caller, keeping the transaction binding explicit.

Exact dependency versions remain implementation-manifest concerns.

### 4.2 Technical auth surface

The baseline same-origin non-Product surface is:

```text
GET  /_auth/login
GET  /_auth/callback
GET  /_auth/csrf
POST /_auth/logout
```

These routes are Technical Non-Product auth/session mechanics. They are excluded from the Product OpenAPI/operation census and do not collide with `/access-context` or `/organizations/{organization_id}/...`.

No CORS relaxation is baseline; the browser application and Product API remain same-origin at `https://conexus.fun`.

### 4.3 Login transaction

Starting login creates one short-lived technical login transaction in PostgreSQL containing only:

```text
one-time login handle digest
state digest
OIDC nonce
PKCE verifier
bounded relative return path
created_at / expires_at
```

The browser receives a separate one-time `__Host-mpc_login` cookie containing the raw random login handle:

```text
Secure
HttpOnly
SameSite=Lax
Path=/
no Domain
short-lived
```

The transaction is one-time consumed on callback. `state` alone is not treated as the browser binding.

The return target is a validated same-origin relative path, never an arbitrary redirect URI.

### 4.4 Authorization request

The server uses provider discovery and sends an Authorization Code request with:

```text
scope=openid
state=<transaction-specific random value>
nonce=<transaction-specific random value>
code_challenge=<S256(PKCE verifier)>
code_challenge_method=S256
```

Do not request `offline_access` in the human baseline. MPC does not need provider API delegated access merely to authenticate a person.

### 4.5 Callback

Callback processing must, in order:

1. validate and one-time consume the login-handle correlation;
2. validate exact `state`;
3. exchange the code server-side using the confidential client and exact registered redirect URI;
4. verify the returned ID token signature, exact configured issuer, audience/client, expiry and nonce;
5. resolve stable `(issuer, subject)` to exactly one current MPC human Principal;
6. require current Principal access eligibility;
7. create a fresh MPC application session;
8. clear the temporary login cookie/transaction;
9. redirect only to the validated relative return path.

Email, username and mutable profile claims never auto-link or identify the MPC Principal.

### 4.6 Human OIDC token minimization

After successful callback validation, the human token response has no continuing MPC consumer.

Therefore baseline behavior is:

- do not expose access, refresh or ID tokens to browser JavaScript;
- do not persist the human access token;
- do not persist the human refresh token;
- do not use a refresh token to keep the MPC session alive;
- do not call UserInfo merely to duplicate MPC presentation identity;
- discard the token set after the verified identity binding has established the MPC session.

If a future accepted feature requires delegated provider access on behalf of the human, reopen the smallest D5/D7 auth scope rather than pre-storing tokens now.

## 5. MPC application session

### 5.1 Session handle

Generate at least 256 bits of CSPRNG entropy for each authenticated session handle.

Browser cookie:

```text
__Host-mpc_session=<opaque random handle>
Secure
HttpOnly
SameSite=Lax
Path=/
no Domain
no persistent Max-Age/Expires baseline
```

The cookie contains no Principal, Organization, Permission, role, expiry or other decodable state.

PostgreSQL stores only a cryptographic digest of the raw session handle, never the replayable handle itself.

### 5.2 Session row

The minimal server-side session record contains proportionately:

```text
session_handle_digest
principal_id
originating_identity_binding reference
created_at
last_seen_at
absolute_expires_at
revoked_at/revocation reason when applicable
csrf_token_digest
```

No Organization is embedded as ambient tenant context. Organization Membership and Permission remain current request-time checks.

Session state is platform/identity-access technical state, not an Organization-owned business entity.

### 5.3 Timeouts

Baseline security envelope:

```text
idle timeout      30 minutes
absolute timeout   8 hours
```

Both are enforced server-side. `last_seen_at` persistence may be conservatively write-coalesced, but the implementation may not silently extend the effective idle bound beyond its documented coalescing tolerance.

The session is not indefinitely sliding. After absolute expiry the human must perform a new OIDC authentication.

These values may become deployment configuration only inside the accepted maximum envelope; loosening the envelope requires an explicit D7 security review rather than an environment-variable accident.

### 5.4 Rotation and revocation

- a new, unrelated MPC session handle is always minted after successful authentication/reauthentication;
- temporary pre-auth login state is never upgraded into the authenticated session ID;
- logout invalidates the server row first and clears the browser cookie;
- current Principal access-eligibility revocation invalidates use of every existing session on the next request;
- removal/revocation of the originating `(issuer, subject)` binding invalidates the session on the next request;
- current Organization Membership/RoleAssignment/Permission changes take effect through normal current authorization checks and do not require cookie claims to expire;
- automatic periodic session-ID renewal is not baseline because the bounded idle/absolute lifetime and mandatory login-time rotation already satisfy the current consumer without introducing renewal races.

Multiple concurrent human sessions are allowed unless a later concrete security/product requirement proves single-session semantics necessary.

## 6. CSRF request-trust realization

### 6.1 Synchronizer token

Each MPC application session owns one high-entropy synchronizer CSRF token. PostgreSQL stores its digest; the browser obtains the raw value only from authenticated:

```text
GET /_auth/csrf
```

The response is `Cache-Control: no-store` and the frontend keeps the token only in memory. It is not placed in URL, localStorage, sessionStorage, persisted application state or a readable cookie.

Every human `POST | PUT | PATCH | DELETE` Product request must carry:

```text
X-CSRF-Token: <session-bound token>
```

The server compares the presented value to the current session token in constant time before owner effect.

The token rotates with the MPC session. It is request trust only: never Principal identity, Membership, Permission, disposition or Governance authorization.

### 6.2 Standard-library cross-origin defense

In addition to the required synchronizer token, use Go `net/http.CrossOriginProtection` as defense in depth for browser-exposed unsafe requests.

The accepted Go floor is already >= 1.25.1, which excludes the initial Go 1.25.0 bypass-pattern vulnerability. No `AddInsecureBypassPattern` is required by the MPC baseline.

Current standard-library behavior uses `Sec-Fetch-Site` and/or `Origin`/Host checks for unsafe cross-origin browser requests. The synchronizer token remains mandatory because the standard control intentionally allows requests lacking browser metadata headers and D5-R1 already requires `X-CSRF-Token`.

Primary evidence:

- <https://pkg.go.dev/net/http#CrossOriginProtection>
- <https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html>

### 6.3 Safe methods

`GET`, `HEAD` and `OPTIONS` do not perform Product business mutations. Technical login initiation may create bounded ephemeral anti-forgery transaction state but never owner/business state or external effect.

## 7. Logout

`POST /_auth/logout` requires the current H session and valid CSRF trust, invalidates the MPC session server-side and clears `__Host-mpc_session`.

Local MPC logout is the security baseline and does not require retaining an ID/refresh token.

For the selected Keycloak profile the browser may additionally be redirected to the discovered RP-Initiated Logout endpoint using the configured `client_id` and pre-registered `post_logout_redirect_uri`. Absence/failure of provider-wide SSO logout does not resurrect the already-invalid MPC session.

Back-channel logout, global session management UI and persisted ID-token hints are deferred absent a concrete consumer/security requirement.

## 8. Machine bearer realization

### 8.1 Credential issuance

Each independently accountable A/S credential boundary uses its own confidential Keycloak client/service account. A client credential is not a Principal; it binds to exactly one current MPC non-human Principal through explicit technical identity binding.

First-profile client authentication to Keycloak uses standard confidential-client credentials. Exact secret injection/rotation is D7-E. Asymmetric `private_key_jwt`, mTLS and DPoP remain bounded security upgrades, not baseline Product requirements.

RFC 9700 recommends sender-constrained access tokens where feasible, but D5-R1 currently admits an audience-bound bearer carrier. D7-D does not silently change that accepted wire contract; a future sender-constraint requirement is a bounded D5/D7 reopen trigger.

### 8.2 Token validation

Select `github.com/lestrrat-go/jwx/v3` as the bounded JWT/JWK verifier for machine access tokens.

Use only the trusted configured/discovered issuer metadata and JWKS URL. Never follow a JWT-controlled `jku`/key URL.

For every A/S bearer request require, before Principal resolution:

```text
valid cryptographic signature from trusted JWKS
allowed configured signing algorithm
exact configured issuer
not expired
not before satisfied when present
required audience includes https://conexus.fun
verified stable Keycloak client identity claim
```

The first Keycloak profile normalizes its verified authorized-client identity (`azp`/equivalent configured provider claim) behind the auth adapter. Product/business code never depends on provider JWT claim spelling.

The verified `(issuer, machine client identity)` binding resolves to exactly one current A or S Principal. It must never resolve H.

Current Principal access eligibility, Organization Membership, allowed Principal kind, ordinary Permission, owner disposition and Governance remain later independent gates.

Keycloak token roles/scopes never become MPC Permissions.

### 8.3 JWKS caching and key rotation

JWKS is cached from the trusted provider metadata endpoint. Unknown `kid`, cache expiry or rotation triggers bounded refresh from that trusted endpoint; inability to establish a valid signing key fails closed.

Existing validated cached keys may be used according to bounded cache lifetime, but a token is never accepted solely because it parsed successfully.

Machine access tokens remain short-lived and client-credentials responses do not establish a refresh-token path. Exact token lifetime is provider/deployment configuration; D7 proof must establish a bounded revocation window and no indefinite bearer lifetime.

## 9. Authentication libraries and authority fence

Accepted candidate dependency set for D7-D:

```text
github.com/coreos/go-oidc/v3/oidc   human OIDC discovery + ID-token validation
golang.org/x/oauth2                 Authorization Code exchange + PKCE
net/http.CrossOriginProtection      cross-origin browser defense in depth
github.com/lestrrat-go/jwx/v3       machine JWT/JWK verification
```

Do not introduce a generic IAM/policy engine, OAuth proxy, separate BFF service, Redis session store, JWT application-session cookie, browser OAuth SDK or Keycloak adapter that maps provider roles into Product authority.

## 10. Keycloak configuration fence

For the first provider profile:

- exact redirect URI and post-logout redirect URI are pre-registered; wildcard production redirects are not baseline;
- human client is confidential and PKCE S256 is required;
- machine clients have service accounts enabled and human Standard Flow disabled;
- machine tokens receive the explicit MPC audience `https://conexus.fun`;
- provider roles/scopes are kept minimal and are not authorization inputs to MPC;
- client secrets/private credentials live outside repository/database/logs and are injected under D7-E;
- admin/break-glass Keycloak credentials are not application runtime credentials.

No realm-per-Organization or client-per-Organization topology is introduced. Organization is MPC authority, not an IdP tenant mapping.

## 11. Falsifiable proof contract

D7-D cannot close D7 without executable proof capable of falsifying at least:

1. browser JavaScript can obtain an OIDC access/refresh token;
2. human access/refresh tokens remain persisted after successful login without a named consumer;
3. callback with wrong/replayed login handle, `state`, nonce or PKCE verifier establishes a session;
4. callback accepts wrong issuer, audience, signature or expired ID token;
5. email/username collision auto-links a Principal;
6. pre-authentication handle becomes the authenticated session handle;
7. leaked session-table contents directly reveal a replayable session token;
8. expired, revoked, disabled-Principal or revoked-binding session remains usable;
9. unsafe H request without/mismatching `X-CSRF-Token` reaches owner effect;
10. cross-origin unsafe browser request bypasses both synchronizer and standard cross-origin protection;
11. session/CSRF token appears in URL, persistent browser storage or logs;
12. session carries ambient Organization/Permission claims that remain valid after current authority changes;
13. machine token with wrong issuer/audience/signature/time window is admitted;
14. JWT-controlled remote key URL can influence signing-key retrieval;
15. unknown Keycloak machine client resolves by token role/scope instead of explicit MPC binding;
16. machine bearer resolves to H;
17. Keycloak role/scope grants an MPC Permission absent current Membership/RoleAssignment authority;
18. revoked MPC Principal eligibility fails to stop an otherwise valid current session/bearer;
19. wildcard CORS/redirect configuration expands the same-origin browser baseline;
20. logout clears only the browser cookie while leaving the server session usable.

Human/browser proof requires a real browser + real OIDC provider lane; machine validation proof requires real signed provider tokens and key rotation/negative fixtures. Mocks alone do not prove the carrier boundary.

## 12. Rejected/deferred alternatives

- browser Authorization Code + PKCE bearer ownership — **REJECT**, already superseded by D5-R1;
- JWT/self-contained MPC application session cookie — **REJECT baseline**, revocation/current-access semantics are simpler and safer with opaque server-side state;
- Redis session store — **REJECT**, PostgreSQL is already the replicated canonical technical store and no latency/scale consumer requires another system;
- persistent human OIDC refresh-token cache — **REJECT**, no delegated downstream token consumer exists;
- UserInfo/profile synchronization — **REJECT baseline**, stable `(issuer, sub)` is identity and MPC owns its presentation/access state;
- Keycloak roles as Product Permissions — **REJECT**;
- realm/client per Organization — **REJECT**;
- Direct Access Grant/password grant — **REJECT**;
- Implicit Flow — **REJECT**;
- global back-channel logout machinery — **DEFER**, local session invalidation + bounded lifetime satisfies current MPC security need;
- DPoP/mTLS/private-key JWT for all machine clients — **DEFER bounded security upgrade**, not silently added to the accepted bearer contract.

## 13. Adjudication

**Candidate:** Keycloak as the first OIDC/OAuth provider; server-side Authorization Code + required PKCE S256 using `go-oidc` + `x/oauth2`; one-time PostgreSQL login transaction; human token set discarded after verified `(issuer, sub)` binding; 256-bit opaque hashed server-side MPC session with 30-minute idle / 8-hour absolute limits; session-bound synchronizer `X-CSRF-Token` plus Go `CrossOriginProtection`; local CSRF-protected logout; Keycloak service-account Client Credentials for A/S with exact MPC audience; strict trusted-JWKS JWT validation through `jwx/v3`; explicit machine-client→A/S Principal binding; no IdP role/scope authorization.

No Product operation, Permission, Principal kind or semantic owner changes.

If ratified, next is **D7-E — Operability / Secrets / Migrations / Deployment & Proof Baseline**.

Do not begin D8, D9 or Product implementation.