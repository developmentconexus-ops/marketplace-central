# D5-B2 — Product Operation / Resource Surface

> **Status:** OPEN / ACTIVE — B2-A + Operation Admission Matrix + Whole-Matrix Global Coherence **ACCEPTED IN-STAGE / OPERATOR-RATIFIED**; Wire Contract next  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + `DECISION-RECONCILIATION-BASELINE.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-18  
> **Whole-Matrix ratified:** 2026-08-18

## 1. Purpose and stage boundary

D5-B2 derives the smallest coherent **Product 1.0 operation/resource surface and wire contract** that real MPC clients need to interact with accepted semantic owners.

B2 is not a legacy-route cleanup exercise. It does not begin from retired candidates, current OpenAPI, controller/package layout, provider endpoint inventory or frontend screen list.

The operation-admission phase is now ratified in `D5-B2-OPERATION-ADMISSION-MATRIX.md`. B2 next crystallizes that accepted meaning into resource/path/schema/HTTP/OpenAPI shape without moving business authority.

B2 determines:

- which external Product API operations are admitted at all;
- which accepted owner owns each operation;
- Q / C / P classification;
- legitimate client classes;
- ordinary Permission requirements;
- Organization and source-qualified identity scope;
- read knowledge/freshness/provenance semantics;
- owner outcomes, idempotency and concurrency/precondition semantics;
- provider enrichment, pagination/filter/search and bulk only when justified;
- the one machine-readable Product API wire authority.

B2 does **not** choose D6 screen/component topology or D7 process, worker, transaction, queue, RLS, blob storage, deployment, Keycloak realm/deployment, secret-storage or token-lifetime realization.

The router remains the sole current-program status/next-action authority.

---

## 2. Governing B2 invariant

> **A Product API operation exists only when a real Product 1.0 client/actor needs to read one accepted owner's meaning (Q), ask one accepted owner to perform/accept owner-owned work (C), or consume a justified read-only composition (P). Public API symmetry, legacy routes, provider endpoints and internal implementation convenience do not create Product operations.**

Every admitted operation maps to:

```text
real consumer/use
  → allowed client class
  → exactly one semantic owner / accepted D2 substrate authority
  → Q | C | P
  → explicit Organization scope where business-owned
  → ordinary Permission
  → canonical/source-qualified subject identity
  → knowledge/outcome contract
  → complete C-operation safety tuple when C
```

Every admitted C declares:

```text
consequence class
idempotency disposition
concurrency / precondition disposition
```

Silence is non-conformant. Idempotency never substitutes for concurrency and never authorizes blind retry of ambiguous external acceptance.

An operation that cannot fit this model without inventing authority stops and returns to the implicated parent decision.

---

# 3. B2-A — Client & Authentication Admission Model — ACCEPTED IN-STAGE

## 3.1 Root cause

D5-B1 defined authorization and Product API semantics but intentionally did not decide the concrete client classes and authentication contract by which a real human application or non-human automation proves who/what is calling the API.

Without that decision, B2 could freeze write-capable operations while leaving reachable failures such as:

- global/shared API keys with excessive blast radius;
- browser-embedded secrets;
- client-supplied Principal identity;
- one credential implicitly authorizing every Organization;
- IdP roles becoming a second business-Permission authority;
- automation attributed to the provisioning human;
- tokens for sibling applications being replayed against MPC;
- unrelated human/machine authentication architectures;
- a second MPC credential/token authority duplicating external IdP capability.

This is essential API security complexity, not merely D7 runtime detail.

## 3.2 Known / inferred / unknown / deferred

### Known

Accepted D2 already establishes:

- interactive human authentication is external through OIDC;
- stable human binding uses OIDC `(issuer, subject)`, not email/username;
- Principal is MPC-owned accountable actor identity;
- automation/system Principals are distinct from humans and never impersonate them;
- Membership / AccessRole / Permission / RoleAssignment are MPC-owned ordinary-access semantics;
- Permission does not prove business disposition or Governance authorization;
- Keycloak is the preferred first self-hosted candidate, while exact provider/deployment/realm topology remains later realization.

Current standards/provider evidence establishes:

- Authorization Code clients use PKCE; current OAuth Security BCP requires PKCE for public clients and recommends it more broadly, with `S256` as the non-verifier-exposing method;
- access tokens should be privilege- and audience-restricted to the intended resource server;
- Client Credentials is for confidential machine clients;
- current Keycloak supports OIDC clients, Authorization Code, service accounts / Client Credentials, audience configuration, client-secret rotation and stronger client authentication such as signed JWT;
- mature machine-integration systems benefit from fine-grained permissions and short-lived tokens instead of broad long-lived credentials.

### Inferred

For MPC's current first-party/internal clients, one standards-based authentication architecture can serve humans and machines without MPC becoming an Authorization Server.

Sharing IdP infrastructure with another internal product such as MetalDocs may reduce operational duplication, but creates no shared Product authority, database, MPC Organization or Permission model.

### Unknown / deferred

D7 still owns:

- exact IdP deployment and final production binding;
- Keycloak realm topology, including any bounded sharing with MetalDocs;
- HA/database/backups/upgrades;
- token lifetime and refresh-token realization;
- client-secret/private-key storage and rotation mechanics;
- exact JWT/JWK validation library/caches;
- whether sender-constrained tokens such as DPoP/mTLS become justified;
- exact machine-token claim/binding mechanics to MPC non-human Principal.

Unknown/deferred never authorizes a second credential architecture by convenience.

## 3.3 Credible alternatives

### A — MPC-owned global/static API keys

Rejected. They create excess blast radius, poor attribution and duplicate credential lifecycle inside Product code.

### B — MPC builds its own OAuth/authorization server

Rejected. Password/session/MFA/client/token/signing/key-rotation machinery is commodity identity infrastructure, not Marketplace Operations authority.

### C — Standards-based external IdP for humans and machines

**Selected Global Maximum.**

```text
human → Authorization Code + PKCE ┐
                                  ├→ external OIDC/OAuth authority
machine → Client Credentials      ┘
                                     ↓ audience-bound access token
                                  MPC resource server
                                     ↓ resolved MPC Principal
                                  Membership / Permission
                                     ↓ business owner disposition
                                  Governance when required
```

## 3.4 Client classes

### Human interactive client

Examples: MPC React application or another explicitly supported user-facing client.

Rules:

- no browser-stored confidential client secret;
- human identity comes from authenticated token context, never Principal request fields;
- email/username is not canonical identity;
- resource-owner password/direct-grant is not baseline;
- Organization-owned business calls remain explicit Organization-path scoped and require current Membership/Permission;
- **exception:** `GetCurrentAccessContext` is a bounded platform-scoped **self-only** discovery Q because the client cannot know its Organization memberships before discovery. It accepts no Principal parameter and returns only memberships of the authenticated Principal.

### Machine / automation / system client

Examples: Claude Code/Codex-style automation, bounded MPC-owned machine clients, or future physically external applications acting under explicitly provisioned non-human identity.

Rules:

- machine is attributed to its own non-human Principal, never the provisioning human;
- provisioning/delegation history may preserve responsible human/authority context separately;
- client credential proves authentication only, not MPC Permission, business disposition, Governance approval or execution validity;
- machine token resolves fail-closed to the intended non-human Principal and requested Organization access;
- no long-lived MPC Product API key exists merely for CLI/agent convenience;
- holding Permission does not prove epistemic/physical ability to establish a fact. Physical Fulfillment facts require human evidence or an explicitly proven system Principal/source.

A future third-party delegated-user application may create materially different OAuth/consent requirements; it is not invented before a real consumer exists.

## 3.5 Authentication authority versus MPC authority

The external IdP/Authorization Server owns only identity/authentication protocol concerns such as credentials, login/MFA/session, OAuth/OIDC clients and token issuance/signing.

MPC remains authority for Principal identity/history, Organization Membership, AccessRole/Permission/RoleAssignment, every D1 business meaning, business disposition, Governance and execution-time validity/reconciliation.

```text
valid token
  != MPC Permission
  != business permitted
  != Governance authorized
  != executable now
  != externally converged
```

## 3.6 No duplicate permission authority

MPC Permission remains the ordinary Product-access authority. IdP roles/scopes may support protocol/client administration, audience and defense-in-depth but never independently grant a Product business operation.

No OAuth-scope vocabulary duplicates MPC business Permissions by default.

## 3.7 Token acceptance boundary

Product API access preserves at least:

- bearer credential only over TLS, never URL/query-string credential;
- trusted issuer/signature chain validation;
- token time validity;
- audience bound to MPC API, rejecting sibling-resource tokens;
- fail-closed resolution to one MPC Principal;
- current MPC Membership/Permission for Organization-owned requests;
- no raw token/client secret in Product audit/history/diagnostic payloads.

Exact JWT library/format is mechanism, not business semantics.

## 3.8 Keycloak disposition

Keycloak remains the **preferred first implementation/proof candidate** and expected D7 realization unless later evidence falsifies it.

B2 freezes OIDC/OAuth behavior, not deployment/realm topology. A shared Keycloak deployment/realm with MetalDocs, if selected in D7, cannot create shared MPC/MetalDocs business authority, Permission, Organization or database semantics.

A Keycloak `Organization`, realm role or client role does not become MPC `Organization`, Membership or business Permission by name similarity.

## 3.9 Product API versus provider protocol ingress

The model above authenticates **Product API clients**.

Marketplace OAuth callbacks, webhook verification and business-system/provider protocol ingress remain D4/D5-B1 technical surfaces. A marketplace credential never authenticates a Product Principal.

No generic `/integrations/auth` business API is created.

## 3.10 Proof / negative controls

Later executable proof must demonstrate at least:

1. token for another audience/resource is rejected;
2. untrusted issuer/signature/expired token is rejected before Product semantics;
3. human maps by trusted external binding, not email/name;
4. machine maps to a distinct non-human Principal;
5. valid AuthN with no current Membership/Permission is denied as ordinary access without pretending business rejection;
6. client cannot supply another Principal or `approved=true` to bypass boundaries;
7. sibling-product token cannot be replayed against MPC merely because the same IdP issued it;
8. revocation/disablement stops future access without rewriting history;
9. raw credentials never enter Product history/diagnostics;
10. business `approval-required/rejected/pending` remains reachable after successful AuthN/ordinary access.

## 3.11 Evidence anchors

External evidence used for this decision includes IETF OAuth Security BCP, OAuth Client Credentials, current Keycloak documentation and mature short-lived/fine-grained machine integration patterns. These sources inform mechanism/security posture; they never override MPC semantic authority.

---

# 4. Operation Admission Model — ACCEPTED / MATRIX RATIFIED

The operator ratified the semantic-owner driven Product API direction and the complete operation matrix in `D5-B2-OPERATION-ADMISSION-MATRIX.md` after Whole-Matrix adversarial review.

Every admitted operation was challenged for:

1. real consumer;
2. human/machine client class;
3. exactly one owner;
4. Q/C/P honesty;
5. Organization scope;
6. ordinary Permission;
7. canonical/source-qualified identity;
8. honest read semantics;
9. consequence/outcome semantics;
10. idempotency;
11. concurrency/preconditions;
12. provider enrichment;
13. collection semantics;
14. bulk necessity.

Whole-Matrix corrections ratified:

- listing-context authored-media intake added without media master;
- Fulfillment internal operating-target Q/C added without SLA/rules platform;
- generic Work resolution deferred in favor of source-owner resolution paths;
- cross-owner Sale operational P deferred until D6 proves need; current baseline has **zero P operations**;
- every admitted C now has a complete safety tuple;
- Party Resolution requires idempotency + current candidate-set/resolution precondition;
- current access discovery is platform-scoped self-only;
- standing authority revocation is fail-safe/monotonic;
- initial publication price always remains a distinct PriceIntent correlated to the pre-creation ListingIntent context; ListingIntent never owns price;
- B2-A OIDC/OAuth + MPC-owned access/business authority remains confirmed.

No D0/D1/D2/D3/D4/D4-R1/D5-B1 reopen was required.

## 4.1 Explicit baseline exclusions

Still not admitted by default:

- Product/PIM CRUD;
- generic `/integrations`, `/mutations`, `/commands`, `/actions`, `/resources` or `/reconcile` business surfaces;
- generic source-ingestion platform before a real external connector;
- generic stock-set/sync;
- generic provider refresh/sync-everything;
- ProductAsset/media library;
- mapping/rule DSL;
- AI-specific API;
- provider endpoint mirrors;
- public event stream without real external consumer;
- generic bulk;
- generic Task/Case/Workflow;
- generic Work resolution command bus;
- cross-owner P/BFF surface before real D6 consumer evidence.

---

# 5. Exact next B2 work — Wire Contract / Resource-Path-Schema Grammar

Derive concrete wire shape **only** from the ratified matrix.

The next sub-batch must determine:

1. resource/path hierarchy, keeping Organization-owned business operations under `/organizations/{organization_id}/...` and giving current access context only its bounded self-only platform discovery shape;
2. standard HTTP methods versus owner-specific methods where CRUD would lie;
3. exact request/response schema families and source-qualified identities;
4. honest known/empty/unknown/unavailable/partial/freshness/provenance representation where material;
5. owner outcomes such as accepted/rejected/pending/ambiguous and later applied/converged distinctions;
6. RFC 9457 Problem Details for API/transport/access/precondition/idempotency/server problems without turning valid business outcomes into access errors;
7. exact `Idempotency-Key` placement/validation and opaque MPC concurrency/precondition mechanism for every admitted C;
8. pagination/filter/search/cursor grammar only for admitted real collections;
9. exact Permission→wire-operation mapping and client-class restrictions;
10. listing-media wire seam without D7 storage/blob/CDN design;
11. provider OAuth/webhook/future external-connector ingress classification outside Product API;
12. the Work closure-path audit required by the matrix before Wire Contract closure;
13. OpenAPI operation naming/spelling and the route to one machine-readable Product API wire authority.

Do not introduce D6 screen/BFF topology, D7 queues/workers/storage/transactions/Keycloak deployment, D8 live-effect proofs or implementation.

---

# 6. Reopen / stop triggers

Targeted reopen is required if material evidence shows:

- a required Product operation cannot fit one accepted owner;
- a client interaction requires a genuinely new D1 semantic edge;
- human/machine authentication cannot fit the accepted OIDC/OAuth boundary without MPC acquiring IdP authority;
- a future third-party delegated-user client introduces materially different consent/delegation semantics;
- a real public/external client creates compatibility/security obligations absent from Product 1.0;
- Keycloak cannot satisfy a material D7 realization requirement and another conformant provider is materially better;
- Wire Contract proves an admitted operation cannot preserve the ratified safety/identity/knowledge/outcome laws;
- Work closure-path audit finds a Product 1.0 actionable condition with no legitimate owner-side closure path;
- D6 proves repeated cross-owner composition pain sufficient to justify a bounded P;
- real Product need proves a mutable PriceDraft is cheaper globally than explicit PriceIntent supersession without weakening lineage/least privilege;
- operation/wire derivation exposes a Product 1.0 responsibility absent from D0/D1.

Framework preference, API symmetry, sibling-product convenience, legacy-route compatibility and provider payload adjacency are not reopen evidence.

Implementation remains blocked until D9.
