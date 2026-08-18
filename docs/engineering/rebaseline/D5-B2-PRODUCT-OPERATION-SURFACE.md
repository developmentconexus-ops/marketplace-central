# D5-B2 — Product Operation / Resource Surface

> **Status:** OPEN / ACTIVE — B2-A Client & Authentication Admission Model accepted in-stage; operation matrix under derivation  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + `DECISION-RECONCILIATION-BASELINE.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-18

## 1. Purpose and stage boundary

D5-B2 derives the smallest coherent **Product 1.0 operation/resource surface** that real MPC clients need to interact with accepted semantic owners.

B2 is not a legacy-route cleanup exercise. It does not begin from the retired pre-R1 candidate, current OpenAPI, controller/package layout, provider endpoint inventory or frontend screen list.

B2 determines, before final path/schema spelling:

- which external Product API operations are admitted at all;
- which accepted owner owns each operation;
- whether each interaction is Q / C / P;
- which client classes may invoke it;
- ordinary Permission requirements;
- Organization and identity scope;
- read knowledge/freshness/provenance semantics;
- consequential intent/outcome/idempotency/concurrency semantics;
- provider enrichment, pagination/filter/sort and bulk only when justified by a real consumer.

B2 does **not** choose D6 screen/component topology or D7 process, worker, transaction, queue, RLS, deployment, Keycloak realm/deployment, secret-storage or token-lifetime realization.

The status/next-action wording in `D5-API.md` §21 predates formal B2 opening. The router is the sole current-program status authority; this file is the current B2 in-stage artifact.

---

## 2. Governing B2 invariant

> **A Product API operation exists only when a real Product 1.0 client/actor needs to read one accepted owner's meaning (Q), ask one accepted owner to perform/accept owner-owned work (C), or consume a justified read-only composition (P). Public API symmetry, legacy routes, provider endpoints and internal implementation convenience do not create Product operations.**

Every admitted operation must map to:

```text
real consumer/use
  → allowed client class
  → exactly one semantic owner / accepted D2 substrate authority
  → Q | C | P
  → explicit Organization scope
  → ordinary Permission
  → canonical/source-qualified subject identity
  → knowledge/outcome contract
  → idempotency/concurrency when consequential
```

An operation that cannot fit this model without inventing authority stops and returns to the implicated parent decision.

---

# 3. B2-A — Client & Authentication Admission Model — ACCEPTED IN-STAGE

## 3.1 Root cause

D5-B1 defined authorization and Product API semantics but intentionally did not yet decide the concrete client classes and authentication contract by which a real human application or non-human automation proves who/what is calling the API.

Without that decision, B2 could enumerate write-capable operations while leaving reachable failure classes such as:

- a global/shared API key with excessive blast radius;
- browser-embedded secrets;
- client-supplied Principal identity;
- one credential implicitly authorizing every Organization;
- IdP roles becoming a second business-Permission authority;
- automation being attributed to the human who provisioned it;
- tokens for another application being accepted by MPC;
- human and machine clients requiring unrelated authentication architectures;
- building a second MPC credential/token authority even though the chosen IdP stack already supplies standards-based machine credentials.

This is essential API security complexity, not a D7 runtime detail: B2 must know which client classes can legitimately invoke each Product operation and how authenticated identity reaches MPC Principal/access semantics.

## 3.2 Known / inferred / unknown / deferred

### Known

Accepted D2 already establishes:

- interactive human authentication is external through an OIDC boundary;
- stable human binding is based on OIDC `(issuer, subject)`, not email/username;
- Principal is MPC-owned accountable actor identity;
- automation/system Principals are distinct from humans and do not impersonate them;
- Membership / AccessRole / Permission / RoleAssignment are MPC-owned ordinary-access semantics;
- Permission does not prove business disposition or consequential Governance authorization;
- Keycloak is the preferred first self-hosted candidate, while exact provider/deployment/realm topology was intentionally deferred to later technical realization.

Current standards/provider evidence establishes:

- OAuth 2.0 Authorization Code clients are protected with PKCE; current OAuth Security BCP requires PKCE for public clients and recommends it for confidential clients, with `S256` as the non-verifier-exposing method;
- access tokens should be privilege-restricted and audience-restricted to the intended resource server, which must validate that audience;
- OAuth Client Credentials is for confidential machine clients;
- current Keycloak supports OIDC clients, Authorization Code, service accounts / Client Credentials, confidential client credentials, audience configuration, client-secret rotation and stronger client-authentication mechanisms including signed JWT;
- mature integration platforms such as GitHub Apps use fine-grained permissions plus short-lived tokens to reduce credential blast radius.

### Inferred

For MPC's current first-party/internal client set, one standards-based authentication architecture can serve both humans and machine clients without MPC becoming an Authorization Server.

Sharing an IdP deployment with another internal product such as MetalDocs may reduce infrastructure duplication, but this is **not** a B2 topology decision and creates no shared Product authority, database, Organization or Permission model.

### Unknown / deferred

D7 still owns:

- exact IdP deployment and final production binding;
- Keycloak realm topology and whether MetalDocs/MPC use one workforce realm or another bounded arrangement;
- HA/database/backups/upgrades;
- token lifetime and refresh-token realization;
- client-secret/private-key storage and rotation mechanics;
- exact JWT/JWK validation library and caches;
- whether/when sender-constrained tokens such as DPoP/mTLS are justified;
- exact claim used to bind a machine service-account/client identity to an MPC non-human Principal.

Unknown/deferred does not authorize a second credential architecture by convenience.

## 3.3 Credible alternatives

### A — MPC-owned global/static API keys

Rejected.

A global or Organization-wide shared secret creates excessive blast radius, weak attribution and a second credential/token lifecycle inside MPC. It also makes browser/client misuse easy and pushes revocation/rotation/security administration into Product code without product value.

### B — MPC builds its own OAuth/authorization server

Rejected.

Password/session/MFA/client/token/signing/key-rotation machinery is commodity identity infrastructure, not Marketplace Operations authority. Reimplementing it adds large accidental/security complexity and duplicates the accepted external-IdP boundary.

### C — Standards-based external IdP for both humans and machine clients

**Selected Global Maximum.**

Use one authentication boundary based on OIDC/OAuth:

- human clients authenticate interactively through Authorization Code + PKCE;
- confidential machine clients obtain short-lived access tokens through Client Credentials / service-account semantics;
- MPC validates the token as a resource server, resolves an MPC Principal, then applies MPC Organization Membership and Permission;
- business owner disposition and Governance remain downstream and independent.

This removes duplicate credential authority while staying small enough for MPC's real scale.

## 3.4 Client classes

B2 recognizes only the client classes needed to classify Product operations.

### Human interactive client

Examples: MPC React application or another explicitly supported user-facing client.

Contract:

```text
human
  ↓ Authorization Code + PKCE
external OIDC/OAuth authorization server
  ↓ access token for MPC API
MPC resource server
  ↓ issuer + subject binding
MPC human Principal
```

Binding rules:

- no browser-stored confidential client secret;
- human identity comes from authenticated token context, never request body/header Principal fields;
- email/username does not become canonical binding;
- resource-owner password/direct-grant flow is not baseline;
- the Product request remains Organization-path scoped and requires current MPC Membership/Permission.

### Machine / automation / system client

Examples: Claude Code/Codex-style automation, an MPC-owned background client when an HTTP client is actually needed, or a future physically external application acting under an explicitly provisioned non-human identity.

Contract:

```text
confidential client
  ↓ client authentication
external OAuth authorization server
  ↓ Client Credentials
short-lived MPC-API access token
  ↓
MPC resource server
  ↓ fail-closed machine identity binding
MPC automation/system Principal
```

Binding rules:

- the machine is attributed to its own non-human Principal, never to the human who provisioned it;
- provisioning/delegation history may record the responsible human/authority context separately;
- one client credential proves authentication only; it does not prove MPC business Permission, action disposition, Governance approval or execution validity;
- a machine token is accepted only when it resolves fail-closed to the intended non-human Principal and the requested Organization is currently accessible to that Principal;
- no long-lived MPC Product API key is baseline merely for CLI/agent convenience.

A future third-party delegated-user application may create a materially different OAuth/consent requirement. It is not invented before a real consumer exists.

## 3.5 Authentication authority versus MPC authority

The external IdP/Authorization Server owns only identity/authentication protocol concerns such as:

- login/authentication ceremony;
- user credentials/MFA/passkeys/session authority;
- OAuth/OIDC client registration/credentials;
- token issuance/signing/key rotation;
- service-account/client credential authentication.

MPC remains authority for:

- Principal identity and historical attribution;
- Organization Membership;
- AccessRole / Permission / RoleAssignment;
- every D1 business authority;
- business disposition;
- Controlled Action Governance;
- execution-time validity and reconciliation.

Therefore:

```text
valid token
  != MPC Permission
  != business permitted
  != Governance authorized
  != executable now
  != externally converged
```

## 3.6 No duplicate permission authority

B2 does not create a second business-permission vocabulary in OAuth scopes or IdP roles.

The baseline rule is:

> **MPC Permission remains the ordinary Product-access authority. IdP/client roles/scopes may support protocol/client administration, audience and defense-in-depth, but they never independently grant a Product business operation that the resolved MPC Principal cannot invoke.**

For machine clients, the smallest baseline is a securely provisioned client/service-account identity mapped to one MPC non-human Principal whose current Membership/RoleAssignment yields the Product Permissions. A future proven need for an additional client-specific permission ceiling may be added only if it reduces a concrete blast-radius failure without creating duplicate business authority.

## 3.7 Token acceptance boundary

For Product API calls, B2 freezes these protected properties independent of implementation library:

- token is carried as a bearer access credential over TLS, never URL/query-string credential;
- signature/trust chain is validated against the configured issuer/authorization-server trust;
- issuer is accepted only from configured trusted identity authority;
- access token is not expired/not-before-invalid;
- token is **audience-bound to the MPC API** (or an explicitly equivalent resource indicator) and MPC rejects a token intended only for another resource/application;
- human or machine token context resolves fail-closed to one MPC Principal;
- Organization path scope remains explicit and current Membership/Permission is checked inside MPC;
- raw token/secret is treated as sensitive security material and does not enter business audit payloads/logs.

Exact JWT format is not elevated to business semantics; the standards boundary must remain replaceable by another conformant OIDC/OAuth provider if a later material reason justifies it.

## 3.8 Keycloak disposition

Keycloak is the **preferred first implementation/proof candidate** already carried from D2 and is fully capable, based on current official documentation, of satisfying the selected human + machine authentication contract.

B2 does **not** take ownership of the provider/deployment decision that D2 intentionally deferred. Therefore:

- B2 freezes the OIDC/OAuth client/authentication contract;
- Keycloak is the first implementation/proof target and expected D7 realization unless later evidence falsifies it;
- D7 binds concrete Keycloak deployment/realm/client configuration and may reopen only the technical provider binding if a material requirement cannot be satisfied;
- choosing a shared Keycloak deployment/realm with MetalDocs remains a D7 operational decision and must not create shared MPC/MetalDocs business authority, Permission, Organization or database semantics.

A Keycloak `Organization`, realm role or client role does not become MPC `Organization`, Membership or business Permission merely because the concepts have similar names.

## 3.9 Product API versus provider OAuth/protocol ingress

The authentication model above is for **clients consuming the MPC Product API**.

Marketplace/provider OAuth callbacks, webhook verification and business-system protocol ingress remain D4 integration-boundary concerns under D5-B1 §4.2. A Mercado Livre credential/session does not authenticate a Product API Principal.

No generic `/integrations/auth` Product business API is created from these protocol needs.

## 3.10 Proof / negative controls

Later executable proof must demonstrate at least:

1. valid token for another audience/resource server is rejected by MPC;
2. untrusted issuer/signature/expired token is rejected before Product semantics;
3. authenticated human resolves to stable MPC Principal by trusted external identity binding, not email/name;
4. machine client resolves to a distinct automation/system Principal and never inherits the provisioning human as effective actor;
5. valid authentication with no current Organization Membership/Permission is denied as ordinary access, without pretending business rejection;
6. client/request cannot supply another Principal identity or `approved=true` to bypass access/Governance;
7. token for MetalDocs or another sibling resource cannot be replayed against MPC merely because the same IdP issued it;
8. revocation/disablement of the relevant identity/client/access path prevents future invocation without rewriting historical attribution;
9. raw bearer/client credentials never appear in Product audit/history or provider diagnostic output;
10. business `approval-required/rejected/pending` remains reachable even after successful authentication and ordinary Permission.

Exact test harness/configuration belongs to later stages, but the protected properties belong to B2.

## 3.11 Evidence anchors

Current external evidence used for this in-stage decision:

- IETF RFC 9700 — OAuth 2.0 Security Best Current Practice: PKCE, privilege restriction, audience restriction, token replay guidance;
- IETF RFC 6749 — OAuth 2.0 Client Credentials restricted to confidential clients;
- Keycloak current Server Administration Guide — OIDC clients, Authorization Code, service accounts / Client Credentials, confidential client credentials, role-scope intersection, audience support and stronger client-authentication options;
- GitHub Apps official architecture — benchmark evidence that mature machine integrations benefit from fine-grained permissions and short-lived access tokens rather than broad long-lived application credentials.

These sources inform mechanism/security posture; they do not override MPC D0–D5 semantic authority.

---

# 4. Operation Admission Predicate — ACCEPTED DIRECTION / MATRIX NEXT

The operator approved the B2 direction that the Product API is **semantic-owner driven rather than CRUD-complete, screen-shaped or provider-shaped**.

A candidate operation is admitted only when all applicable questions can be answered:

1. **Consumer:** which real Product 1.0 actor/client needs it and for what outcome?
2. **Client class:** human, machine/automation/system, or both?
3. **Owner:** exactly one D1 semantic owner or accepted D2 substrate authority?
4. **Interaction:** Q, C or P?
5. **Organization:** why is the path scope legitimate and are all secondary references same-Organization?
6. **Permission:** what ordinary product capability allows invocation without conflating Governance/business disposition?
7. **Identity:** what canonical or source-qualified subject does the contract address?
8. **Read semantics:** how are known/empty/unknown/unavailable/partial/freshness/provenance represented when material?
9. **Consequence semantics:** does it create/advance a durable owner-local Intent/work item? What accepted/rejected/pending/ambiguous states are reachable?
10. **Idempotency:** is a client key mandatory or is a structural exemption actually proven?
11. **Concurrency/preconditions:** can stale client state cause unsafe overwrite/action?
12. **Provider enrichment:** which named consumer/correctness property justifies any provider-specific field?
13. **Collection semantics:** are pagination/filter/sort/cursor genuinely required?
14. **Bulk:** is there a real workflow requiring member-level bulk semantics rather than looped individual operations?

Reject the operation if its only justification is symmetry, current code, provider capability, debug convenience or hypothetical future use.

## 4.1 Explicit baseline exclusions

Do not admit by default:

- Product/PIM CRUD;
- generic `/integrations`, `/mutations`, `/commands`, `/actions`, `/resources` or `/reconcile` business surfaces;
- generic source-ingestion API before a real external connector consumer;
- generic stock-set/sync endpoint merely because AvailabilityIntent exists internally;
- generic provider refresh/sync-everything endpoint;
- ProductAsset/media library;
- mapping/rule DSL;
- AI-specific API;
- raw provider endpoint mirrors;
- public event stream without a real external consumer;
- bulk by symmetry;
- BFF/cockpit write operations that merge semantic owners.

---

# 5. Exact next B2 work

Build the **Operation Admission Matrix** owner by owner from Product 1.0 actors and accepted authority.

For every candidate operation record:

```text
operation family / candidate operation
consumer + use
allowed client class
semantic owner
Q | C | P
Organization scope
ordinary Permission
subject identity
knowledge/freshness contract
intent/outcome semantics
idempotency
concurrency/preconditions
provider enrichment
pagination/filter/sort
bulk
ADMIT | REJECT | DEFER
reason / parent authority
```

Derive the matrix in semantic order, not legacy package order:

1. identity/access substrate only where a Product client truly needs it;
2. Marketplace Portfolio;
3. Product & Channel Readiness;
4. Marketplace Offering Operations;
5. Availability Control;
6. Market Intelligence;
7. Commercial Economics;
8. Controlled Action Governance;
9. Marketplace Sales;
10. Business-System Materialization;
11. Fulfillment Lifecycle;
12. Post-Sale Resolution;
13. Operational Work;
14. justified read-only P compositions.

Do not spell final paths/schemas until the admission inventory is coherent enough that naming cannot hide duplicate/missing authority.

---

# 6. Reopen / stop triggers

Targeted reopen is required if material evidence shows:

- a required Product operation cannot fit one accepted owner;
- a client interaction requires a new D1 semantic edge;
- a machine/human authentication requirement cannot be represented by the accepted OIDC/OAuth boundary without MPC acquiring identity-provider authority;
- a future third-party delegated-user application creates materially different consent/delegation semantics;
- a real public/external client creates compatibility/security obligations absent from Product 1.0;
- Keycloak cannot satisfy a material D7 realization requirement and another conformant provider is materially better;
- operation enumeration exposes a Product 1.0 responsibility absent from D0/D1.

Framework preference, desire for API symmetry, sibling-product convenience and current-route compatibility are not reopen evidence.

Implementation remains blocked until D9.