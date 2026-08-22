# D7-D — Authentication / Session / CSRF / Machine Token Realization

> **Status:** OPERATOR-RATIFIED  
> **Parent:** `D7-RUNTIME-JOBS-TRANSACTIONS.md`  
> **Accepted prerequisites:** D7-A + D7-B + D7-C — OPERATOR-RATIFIED  
> **Parent authority:** accepted D5-R1 Human Browser Authentication Correction + accepted D2 identity/access semantics  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Derived:** 2026-08-22  
> **Ratified:** 2026-08-22

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

### 3.1 Keycloak — ACCEPTED

Keycloak is selected as the first concrete OIDC/OAuth provider while the application boundary remains standards-based and provider-neutral.

Keycloak provides only external authentication/token issuance. MPC keeps canonical Principal, current Principal access eligibility, Organization Membership, RoleAssignment, AccessRole/Permission and business/Governance authorization authority.

No Keycloak realm role, client role, group or scope becomes an MPC ordinary Permission or business authorization shortcut.

Baseline deployment uses one MPC realm/client set, not one realm per Organization. Organization is MPC tenancy, not IdP tenancy.

## 4. Human login mechanics — accepted

### 4.1 Authorization Code + PKCE S256

The browser begins login through a same-origin MPC technical auth route. MPC acts as a confidential OIDC client and performs the Authorization Code exchange server-side.

Each login transaction creates fresh high-entropy:

- `state`;
- OIDC `nonce`;
- PKCE `code_verifier` and derived `S256` challenge;
- bounded post-login return target.

The login transaction is one-time and short-lived. It may be persisted in PostgreSQL because callbacks may land on another application replica. It stores no end-user password and no reusable browser credential.

The callback validates exact one-time `state`, exchanges the code with the stored verifier, verifies the ID token signature/issuer/audience/expiry/nonce through OIDC discovery, and resolves the stable `(issuer, sub)` identity binding to exactly one current MPC human Principal.

Email, preferred username and display claims never auto-bind or merge Principals.

### 4.2 Human OIDC token lifetime inside MPC

MPC does **not** retain the human access token or refresh token after successful login in the Product 1.0 baseline.

After callback validation and Principal binding:

1. validate the OIDC token response and ID token;
2. establish the MPC ApplicationSession;
3. discard the human OIDC access/refresh/ID token material from normal application state.

Reason: the Product browser does not need delegated calls to another OAuth resource server after authentication. Persisting or refreshing IdP tokens would increase secret-bearing state without a current consumer.

If a later accepted human flow must call an external resource on that human's delegated authority, reopen only this D7-D token-retention boundary.

## 5. Human ApplicationSession — accepted

### 5.1 Representation

The browser receives only:

```text
__Host-mpc_session=<opaque random handle>
Secure
HttpOnly
SameSite=Lax
Path=/
no Domain
```

The session handle contains no Principal, Organization, Permission or business state. It is generated from cryptographically secure randomness with at least 256 bits of entropy.

PostgreSQL stores only a one-way digest of the presented handle together with bounded mechanism state such as:

```text
session_id / digest
principal_id
created_at
last_seen_at
absolute_expires_at
idle_expires_at
rotation lineage when needed
revoked_at when revoked
```

A database disclosure therefore does not directly disclose reusable browser session handles.

Session state is platform authentication mechanism state, not a Product entity.

### 5.2 Session lifetime

Baseline:

```text
idle timeout      30 minutes
absolute lifetime 8 hours
```

A valid authenticated request may advance the idle deadline without extending the absolute deadline. Implementations should avoid a write on every request by using a bounded touch interval while never accepting a session whose effective idle or absolute deadline has passed.

The concrete numbers are security/operability defaults, not Product semantics; they may be tightened without Product reopen. Materially lengthening them requires a D7 security review.

### 5.3 Rotation and invalidation

- successful OIDC login creates a fresh handle; pre-login IDs never become authenticated session IDs;
- logout revokes the current session server-side and expires the cookie;
- Principal access eligibility is rechecked from current MPC authority on every Product request, so disabling a Principal blocks future access even if the session row remains;
- session lookup resolves one Principal only; it never caches current Organization Membership/Permission as authority;
- sensitive identity-administration changes may revoke all sessions for the Principal when the change requires it;
- no process-memory-only session store is allowed because D7-A permits multiple replicas;
- Redis is not baseline: PostgreSQL already provides the required shared durable store at current scale.

Global Keycloak Single Logout/back-channel logout is deferred because current MPC sessions do not retain IdP tokens and no accepted requirement says remote IdP logout must instantly revoke every MPC session. A real requirement may reopen the smallest D7-D seam.

## 6. CSRF / browser request-trust — accepted

Unsafe human requests (`POST | PUT | PATCH | DELETE`) require both the current HttpOnly session and a current CSRF proof.

### 6.1 Synchronizer-token baseline

Each ApplicationSession has a cryptographically random server-held CSRF secret/token. A same-origin technical endpoint such as:

```text
GET /_auth/csrf
```

returns the current token only to an already authenticated session. The frontend keeps it in memory and sends:

```text
X-CSRF-Token: <token>
```

on unsafe requests.

The token is not placed in a browser-readable cookie, URL, localStorage or sessionStorage.

Baseline comparison is constant-time and bound to the current ApplicationSession. Rotate it when the session handle is rotated or security-sensitive reauthentication occurs.

### 6.2 Cross-origin defense in depth

Use Go `net/http.CrossOriginProtection` on state-changing browser routes in addition to the synchronizer token. Same-origin requests are admitted; unsafe cross-origin browser requests fail before owner business effects.

Public technical endpoints whose protocols legitimately require cross-origin/provider callbacks are not silently exempted from trust validation: they live behind their own D4/D5 Technical Ingress validation and explicit bypass route registration rather than a global CORS relaxation.

Baseline application exposes no wildcard CORS. The first-party frontend and Product API share `https://conexus.fun`.

## 7. OIDC/OAuth Go libraries — accepted

Use:

```text
github.com/coreos/go-oidc/v3/oidc
golang.org/x/oauth2
```

for OIDC discovery/ID-token verification and Authorization Code/PKCE exchange respectively.

These are protocol/mechanism dependencies only. A thin MPC auth package owns the accepted constraints: issuer/client/redirect configuration, state/nonce/PKCE login transaction, stable Principal binding and session establishment.

Do not build a generic IdP abstraction/framework. Provider-neutral OIDC boundaries plus one Keycloak implementation are sufficient.

## 8. Machine A/S bearer realization — accepted

### 8.1 Keycloak client credentials

Each admitted machine integration/automation credential is a confidential Keycloak client/service account using Client Credentials. Shared human credentials are never reused for machine execution.

The token must be short-lived and carry an audience including the stable MPC Product audience:

```text
https://conexus.fun
```

Keycloak token/client roles/scopes may support IdP protocol configuration but are not MPC ordinary Permissions.

### 8.2 Verification

Use `github.com/lestrrat-go/jwx/v3` for machine JWT/JWKS verification against the configured trusted Keycloak issuer/JWKS boundary.

Validate at least:

- accepted signing algorithm/key from trusted JWKS;
- exact trusted issuer;
- MPC audience;
- expiration and not-before/time validity;
- expected token/client identity claim used by the configured binding;
- token is not accepted through the H session carrier.

Unknown `kid` triggers bounded JWKS refresh/retry; it does not fall back to unverified decode or an arbitrary JWK URL from the token.

JWKS may be cached in process memory because it is public verification material and can be reacquired. Cache expiry/refresh must not accept a token under a key outside the configured issuer trust boundary.

### 8.3 Machine Principal binding

MPC stores an explicit technical binding from the trusted issuer + machine client identity to exactly one current Principal whose kind is `A` or `S`.

No token claim creates a Principal or selects Organization/Permission dynamically.

After bearer verification:

```text
trusted machine binding
  -> exactly one current A/S Principal
  -> current Principal eligibility
  -> Organization Membership when required
  -> allowed Principal kind
  -> ordinary Permission
  -> remaining W4/business/Governance gates
```

A bearer resolving to no binding, multiple bindings or an H Principal fails closed.

## 9. Secrets and credential handling carried to D7-E

D7-D establishes what secrets exist; D7-E decides production injection/rotation mechanics.

Secret-bearing runtime inputs include proportionately:

- Keycloak confidential-client secret for human code exchange;
- machine client secrets only in the external machine clients that own them, not in the MPC server unless MPC itself needs a distinct client;
- database runtime credential;
- provider/business-system adapter credentials already owned by D4 realization;
- any cryptographic application key later proven necessary.

ApplicationSession handles/CSRF tokens and transient OIDC code/verifier/token material never enter logs/traces/metrics.

## 10. Rejected/deferred alternatives

- browser-held OAuth access/refresh tokens — **REJECT**, contradicts D5-R1;
- persistent human refresh-token/token cache — **REJECT baseline**, no delegated-resource consumer;
- stateless JWT MPC browser session — **REJECT baseline**, complicates immediate revocation/current access checks and exposes mechanism claims to the browser without need;
- in-memory-only session store — **REJECT**, unsafe across allowed replicas/restarts;
- Redis session store — **REJECT baseline**, second state service without current need;
- Keycloak roles/groups as MPC Permissions — **REJECT**, violates D2/D5 authority;
- realm-per-Organization — **REJECT baseline**, duplicates MPC tenancy into IdP mechanism;
- password/direct-access grant — **REJECT**;
- OAuth implicit flow — **REJECT**;
- refresh-token rotation infrastructure for humans — **DEFER** unless human delegated-resource access appears;
- generic IAM/policy engine — **REJECT**;
- wildcard cross-origin CORS — **REJECT**;
- universal auth middleware that treats H cookie and A/S bearer as interchangeable identity sources before carrier-specific validation — **REJECT**.

## 11. Falsifiable proof contract

D7-D cannot close without proof design capable of falsifying at least:

1. browser JavaScript can obtain an OIDC access/refresh token through the normal login flow;
2. callback with wrong/replayed state, wrong nonce, wrong issuer or wrong audience creates a session;
3. pre-login/fixated session handle survives successful authentication;
4. raw session handle stored in PostgreSQL remains directly reusable after a database leak fixture;
5. expired/revoked/idle-expired session still authenticates;
6. disabling Principal access eligibility leaves Product access usable through an old session;
7. unsafe H request without valid session-bound `X-CSRF-Token` reaches owner effect;
8. cross-origin unsafe H request bypasses request-trust defenses;
9. Product handler accepts browser bearer for an H Principal;
10. machine token with wrong issuer/audience/signature/time validity succeeds;
11. arbitrary token-supplied JWK/JWKS location alters the trust root;
12. unknown/ambiguous machine client binding succeeds;
13. A/S bearer resolves to H;
14. Keycloak role/group/scope grants an MPC Permission absent current MPC RoleAssignment;
15. server stores persistent human access/refresh tokens despite no admitted consumer;
16. logs/traces/metrics reveal session, CSRF, OAuth code/verifier/token, client secret or provider credential.

Real OIDC/Keycloak integration proof is required for callback/token/client-credentials/JWKS claims; pure JWT/session mocks cannot close the slice. CSRF/cookie behavior requires a real HTTP/browser-capable proof seam. D8 will exercise composed Product flows after D7 closes.

## 12. Current evidence basis

Revalidated 2026-08-22 against current primary/security sources:

- OAuth 2.0 Security Best Current Practice, RFC 9700;
- IETF OAuth browser-based apps BCP draft in RFC Editor queue, recommending BFF/server-side patterns for sensitive business apps;
- current Keycloak securing-apps/server administration documentation for Authorization Code, service accounts/client credentials, audience and token behavior;
- Go `net/http.CrossOriginProtection` current standard-library API;
- OWASP Session Management and CSRF guidance;
- current `coreos/go-oidc/v3`, `golang.org/x/oauth2` and `lestrrat-go/jwx/v3` APIs.

Exact dependency versions, timeout constants and Keycloak realm export are implementation-manifest/configuration concerns unless a version-specific behavior becomes architectural.

## 13. Adjudication

**OPERATOR-RATIFIED:** Keycloak is the first OIDC/OAuth provider; human login uses confidential Authorization Code + PKCE S256 with one-time state/nonce/verifier, server-side ID-token validation and stable `(issuer, sub)` Principal binding; the human OIDC token set is discarded after callback; PostgreSQL stores only digests for opaque MPC sessions with 30-minute idle / 8-hour absolute baseline; unsafe browser requests require a session-bound synchronizer CSRF token plus Go cross-origin protection; A/S use Keycloak Client Credentials with audience-bound JWT verification through trusted JWKS and explicit machine-client → A/S Principal binding. No IdP role becomes MPC Permission.

Next is **D7-E — Operability / Secrets / Migrations / Deployment & Proof Baseline**.

Do not begin D8, D9 or Product implementation.