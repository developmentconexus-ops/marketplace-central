# D5-B2 — Technical Non-Product Ingress

> **Status:** OPEN / ACTIVE — A External Acquisition Ingress + B OAuth / Authorization Ceremony **ACCEPTED IN-STAGE / OPERATOR-RATIFIED**; Whole-Ingress adversarial coherence = NEXT  
> **Parent D5:** `D5-API.md`  
> **Operation authority:** ratified D5-B2 Operation Admission Matrix + canonical W1/W2/W3/W4  
> **External authority:** canonical D2 + D3 + D4 + D4-R1  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-19  
> **Ingress-A accepted:** 2026-08-19  
> **Ingress-B accepted:** 2026-08-19

## 1. Purpose and stage boundary

This artifact classifies Product 1.0 inbound technical surfaces that are **not** ordinary Product Principals invoking Product API operations.

It closes the wire seam between external provider protocol and accepted MPC business authority without creating a generic Integration/Ingress business domain.

Two distinct lanes are accepted:

```text
A. External Acquisition Ingress
   external change/recovery signal
   → recoverable native acquisition
   → authoritative D4 reread

B. OAuth / Authorization Ceremony
   Product-authorized human initiation
   → provider authorization protocol
   → callback
   → provider credential + seller/account namespace binding
```

These lanes share an external protocol boundary but **do not share one generic event model**.

This artifact does not choose D7 queue/broker/database/worker/storage/secret/locking/retry/deployment realization. Provider technical routes are not Product API business operations, are not included in the Product SDK, and are not governed by W4 Product Permission checks except where a Product-authorized human explicitly initiates the OAuth ceremony.

Implementation remains blocked until D9.

---

# 2. Governing invariants

1. **Provider transport is adapter-local.** HTTP webhook, queue, event bus, OAuth callback and provider-specific wire vocabulary remain protocol concerns.
2. **Native ingress is mechanism, not business authority.** It may centralize custody, correlation, replay/recovery and technical provenance without owning Sale, Listing, Shipment, Payment, Claim, Market, OAuth business policy or other D1 meaning.
3. **Provider signal != D3 domain event.** A notification is an acquisition trigger unless D4 explicitly proves that the external occurrence itself must be preserved as material evidence.
4. **Provider signal != current provider truth.** Current material state comes from the accepted D4 authoritative acquisition/reread contract.
5. **Organization is never payload-derived authority.** It is obtained only from an exact current MPC-owned MarketplaceInstallation/SourceInstance binding.
6. **Closed typed vocabulary.** Unknown provider topics do not become generic `ExternalEvent` objects.
7. **Positive provider acknowledgement has a technical ceiling.** It means recoverable technical custody/disposition only, never business acceptance/completion/convergence.
8. **Push/poll/recovery converge.** Different discovery mechanisms feed the same typed authoritative acquisition path.
9. **Delivery dedup != business idempotency.** Provider delivery IDs/coalescing are technical optimizations; source-qualified identity and owner semantics protect correctness.
10. **PII/raw payload minimization remains binding.** Ingress does not become a raw-provider archive.
11. **Product access and provider protocol trust are independent.** Provider credentials/signatures never authenticate Product Principals; Product bearer tokens never prove provider callbacks.
12. **OAuth state is correlation, not durable authorization.** Callback completion revalidates current MPC authority and provider seller/account identity.
13. **Credential activation is generation-safe.** A stale refresh or older authorization attempt can never overwrite a newer active credential generation.
14. **No generic ingress framework by symmetry.** A second provider may reuse proven technical seams, but provider-specific capabilities remain explicit.

---

# 3. Global Maximum

Rejected structures:

```text
Webhook business domain
Generic Integration API
Generic ExternalEvent { type, payload, metadata }
provider topic = MPC domain event
one webhook endpoint per D1 owner
provider-specific business processing pipeline
universal subscription/provider-resource framework
OAuth workflow engine / generic OAuth domain
Credential CRUD Product API
```

Selected structure:

```text
EXTERNAL ACQUISITION
provider-specific inbound transport adapter
        ↓
MPC-native recoverable acquisition ingress
        ↓
closed typed acquisition request
        ↓
accepted D4 authoritative acquisition/reread
        ↓
consumer-owned semantic translation
        ↓
D1 owner commits meaning
        ↓
D3 communication only after owner commit

AUTHORIZATION CEREMONY
Product-authorized human
        ↓
server-bound authorization attempt
        ↓
provider Authorization Code ceremony
        ↓
callback verification + current MPC authority revalidation
        ↓
provider-authoritative seller/account proof
        ↓
complete credential-generation activation
```

**Method disposition:** `RESTRUCTURE NOW — D5-LOCAL TECHNICAL INGRESS` relative to direct provider-route→domain processing and browser/session-driven OAuth binding.

No D0→W4 semantic parent reopen is currently required.

---

# 4. Lane A — Provider-specific inbound adapter

The inbound adapter owns only protocol concerns needed to establish an admissible acquisition signal candidate:

- HTTP/SQS/EventBridge/provider delivery mechanics;
- origin/authentication/verification supported by the selected provider contract;
- provider DTO parsing;
- closed provider topic/resource grammar;
- provider application/account/seller/source evidence;
- provider delivery discriminator where actually supplied;
- translation from provider resource syntax to a bounded source-qualified acquisition candidate.

It does **not** choose D1 business meaning from provider topic names, infer Organization, commit owner state, expose raw DTOs to owners, fabricate an event identity, or fetch arbitrary callback-supplied URLs.

A provider resource URL/path is parsed against closed adapter-local grammar, converted to a native key/reference, and reread through the accepted D4 sanctioned operation.

---

# 5. Lane A — MPC-native external acquisition ingress

The native acquisition seam may own technical mechanism/state required for:

- exact current MarketplaceInstallation/SourceInstance correlation;
- Organization derivation from that binding;
- recoverable custody before positive external acknowledgement;
- replay/recovery bookkeeping;
- bounded technical quarantine;
- dispatch of a closed typed acquisition request;
- minimum technical provenance/observability needed to prove recovery and attribution.

It owns none of the business meaning acquired afterward.

## 5.1 Minimal technical context

A valid realization may carry proportionately:

```text
organization_id                  server-derived after exact correlation
marketplace_installation_id?     exactly one when marketplace-scoped
source_instance_id?              exactly one when source-scoped
mpc_received_at
provider_sent_at?                source evidence only
provider_delivery_key?           optional technical dedup evidence
bounded technical provenance     only what recovery/audit requires
```

Exact schema, ID encoding and persistence remain D7. Do not create a universal payload/property bag.

## 5.2 Accepted typed acquisition families

```text
AcquireMarketplaceListing
AcquireMarketplacePrice
AcquireMarketplaceSale
AcquireMarketplaceShipment
AcquireMarketplacePayment
AcquireMarketplacePostSaleClaim
AcquireMarketplaceCompetitivePosition
```

These names mean **what source evidence must be acquired**, not what the provider claims happened. They are internal technical contracts, not Product API operations.

---

# 6. Lane A — Exact namespace / Organization correlation

```text
provider application/account/seller/source evidence
        ↓
current MPC external-namespace binding
        ↓
exactly one MarketplaceInstallation or SourceInstance
        ↓
Organization from that MPC-owned binding
```

Forbidden:

```text
provider user_id = OrganizationId
application_id = tenant
payload organization_id = authority
first matching Installation wins
ambient/default Organization
```

- exactly one compatible binding → acquisition may proceed;
- zero bindings → no D1 attribution;
- multiple bindings → no D1 attribution;
- authoritative provider namespace contradiction → no D1 attribution.

Unattributable but otherwise admissible signals may enter technical quarantine only under §7.

---

# 7. Lane A — Technical quarantine

Technical quarantine is admitted because provider retry/recovery windows can be bounded and silent discard would violate D3/D4 recovery expectations.

It is:

- technical/platform state, not a D1 resource;
- not Product API-visible by default;
- not Organization-attributed until exact binding exists;
- PII-minimized;
- unable to dispatch owner acquisition while attribution is unresolved;
- recoverable/reprocessable after a legitimate binding correction.

Typical causes:

```text
zero current bindings
multiple current bindings
provider namespace contradiction
known provider signal whose selected mapping is temporarily unavailable
```

Do not create Work from an unbound technical signal. Storage, retention, alerting and operator tooling remain D7/D6 as justified later.

---

# 8. Lane A — Positive acknowledgement / custody law

> **A positive external-delivery acknowledgement may be emitted only after MPC has established recoverable responsibility for the admitted signal or an explicit quarantine disposition.**

Thus:

```text
HTTP 200 / queue delete / equivalent acknowledgement
!= Sale committed
!= Listing updated
!= Payment processed
!= business accepted
!= externally converged
```

If a process crash immediately after positive acknowledgement can permanently lose required acquisition, the realization violates this contract.

D7 chooses the custody mechanism.

---

# 9. Lane A — Push / poll / recovery convergence

```text
live provider notification ───────┐
provider missed-feed/recovery ────┤
scheduled reconciliation/poll ────┼→ same typed acquisition request
cold-start/source scan ────────────┤
manual technical recovery ────────┘
```

Push improves latency/freshness but never becomes completeness authority. Notification gaps, replay, recovery-window exhaustion and provider retention limits remain honest partial/unknown coverage under D4/W3.

Internal D3 owner events do not loop through this ingress for transport symmetry.

---

# 10. Lane A — Duplicate, ordering and coalescing

```text
provider delivery dedup
!= technical acquisition coalescing
!= owner semantic idempotency/current truth
```

Provider delivery keys may reduce duplicate work but never become Sale/Listing/Shipment/Payment/Claim identity.

Current-state invalidation signals may be coalesced only when the D4/owner claim loses no material information. Material occurrence evidence may be coalesced only when authoritative acquisition still recovers every occurrence required by D3 evidence-edge correctness.

If latest reread cannot recover a materially required external occurrence, D4 must explicitly preserve/classify that occurrence rather than inventing universal webhook event sourcing.

---

# 11. Lane A — Mercado Livre Ingress Admission Matrix

Provider topic spelling remains adapter-local evidence. Admission below freezes MPC treatment, not provider ontology.

## 11.1 ADMIT

| Provider signal/topic | Native acquisition | Authoritative D4 acquisition | Consumer meaning |
|---|---|---|---|
| `orders_v2` | `AcquireMarketplaceSale` | source-qualified Order point reread | Marketplace Sales |
| `shipments` | `AcquireMarketplaceShipment` | source-qualified Shipment reread | Fulfillment |
| `items` | `AcquireMarketplaceListing` | source-qualified Item/Listing reread | Offering; admitted evidence may feed Readiness/Availability |
| `items_prices` | `AcquireMarketplacePrice` | selected current price read | Offering/PriceIntent observation; qualified Economics input |
| `payments` | `AcquireMarketplacePayment` | accepted D4 Payment read | Commercial Economics + Post-Sale evidence |
| `post_purchase:claims` | `AcquireMarketplacePostSaleClaim` | source-qualified Claim reread | Post-Sale Resolution |
| `post_purchase:claims_actions` | same Claim acquisition | same Claim reread | Post-Sale Resolution |
| `catalog_item_competition_status` | `AcquireMarketplaceCompetitivePosition` | selected competition/price-to-win read | Market Intelligence |

The two post-purchase topics intentionally map to one acquisition family. `items` may feed multiple consumer-owned semantic ports without creating a provider-resource owner.

## 11.2 DEFER SAFELY

| Provider signal/topic | Reopen trigger |
|---|---|
| `stock-location` | real multi-origin/User-Product stock lane accepted |
| `fbm_stock_operations` | Full/FBM becomes selected Product flow |
| `flex-handshakes` | Flex becomes selected Fulfillment lane |
| `invoices` | provider-native invoice lane becomes selected |
| `user-products-families` | proactive family drift becomes a real correctness consumer |
| `catalog_suggestions` | catalog/brand suggestion gets named consumer |
| `price_suggestion` | Market/Economics adopts that evidence |
| promotion candidates/offers / `best_price_eligible` | promotion participation enters Product scope |
| provider image-error/picture feeds | named media/readiness consumer proves need |

A DEFER row is not pre-authorization to subscribe or implement.

## 11.3 REJECT BASELINE

```text
legacy orders when orders_v2 is selected
created_orders / unconfirmed checkout feed
order feedback/reputation
buyer/seller messages/chat
questions/Q&A
leads / vertical-specific lead signals
unknown/unclassified provider topics
```

The initial Mercado Livre subscription direction is therefore, subject to live capability revalidation:

```text
orders_v2
shipments
items
items_prices
payments
post_purchase:claims
post_purchase:claims_actions
catalog_item_competition_status
```

Provider application configuration follows this matrix; it never defines admission authority.

---

# 12. Product API / OpenAPI separation

Technical ingress:

- is outside `/organizations/{organization_id}/...` Product roots;
- is not `/access-context`;
- is not included in the Product SDK;
- is not itself authorized with W4 Product Permissions;
- never uses provider credentials to impersonate Product Principals;
- never uses Product bearer tokens as provider callback proof;
- may use provider-specific route/transport vocabulary inside the protocol boundary.

The single Product OpenAPI remains the Semantic Product API wire authority. Provider protocol conformance and internal typed ingress contracts are separate technical contracts, not parallel Product business wire authorities.

Exact provider-facing route spelling remains later final ingress/Wire closure; do not create generic Product `/providers`, `/integrations` or `/webhooks` roots.

---

# 13. Lane B — Authorization Ceremony boundary

OAuth begin/callback/reauthorization is **not** an acquisition signal.

Its purpose is to establish or renew the provider authorization/credential binding of an existing MarketplaceInstallation without allowing browser/session/provider identity to become business authority.

```text
Product-authorized human
        ↓
explicit technical begin
        ↓
server-bound authorization attempt
        ↓
provider authorization protocol
        ↓
callback verification
        ↓
current MPC authority revalidation
        ↓
provider seller/account identity proof
        ↓
initial bind or same-seller reauthorization
        ↓
complete credential-generation activation
```

OAuth callback never creates Organization or MarketplaceInstallation.

---

# 14. Lane B — Initiation authority

Only an authenticated **human** may initiate the current Mercado Livre authorization ceremony.

At begin time all must be current:

```text
valid Product AuthN
Principal.kind = human
current Principal access eligibility
current Organization Membership
portfolio.manage
MarketplaceInstallation belongs to that Organization
Installation is eligible for authorization/reauthorization
```

`portfolio.manage` means only that the human may initiate the technical connection ceremony for that Installation. It does not authenticate the provider callback or grant provider seller authority.

The begin is an explicit state-creating technical action; ambient/cross-site GET navigation must not silently create an authorization attempt. Exact route/method/redirect status remain final ingress wire spelling.

---

# 15. Lane B — Technical Authorization Attempt

One bounded technical Authorization Attempt correlates the ceremony. It is not a Product resource or business workflow.

It binds proportionately:

```text
provider application/client
Organization
MarketplaceInstallation
initiating MPC Principal
opaque state
PKCE verifier/challenge binding when supported/selected
created/expiry boundary
server-controlled safe post-callback destination
technical attempt generation
```

Exact storage/ID/TTL mechanics remain D7.

## 15.1 `state` law

For the selected Mercado Livre lane, `state` is:

```text
high entropy
opaque
single-use
finite-lived
bound to exactly one authorization attempt
```

`state` is correlation/anti-forgery proof only. It never proves current Membership, `portfolio.manage`, Installation validity, seller compatibility or authorization to complete after those facts change.

The provider redirect URI remains static. Organization, Installation and dynamic return targets are not trusted from callback query parameters.

## 15.2 One live attempt per Installation/application

For one MarketplaceInstallation + provider application, at most one unfinished authorization attempt is current.

A new authorized begin supersedes the prior unfinished attempt. A callback for the superseded state cannot activate credentials.

This prevents callback races without introducing a generic workflow engine.

---

# 16. Lane B — Provider Authorization Code + PKCE

For the current Mercado Livre lane, selected ceremony protection is:

```text
Authorization Code
+ server-issued state
+ PKCE where supported by the current provider contract
```

Provider/client secret, authorization code, PKCE verifier, access token and refresh token are sensitive technical credentials and must not enter Product schemas, browser storage, normal logs, Problem Details or business history.

The authorization-code exchange is server-to-server.

A provider that later lacks PKCE uses its actual strongest sanctioned protocol; MPC does not fabricate a universal provider capability.

---

# 17. Lane B — Callback revalidation

The callback is not authorized by whatever Product browser session happens to exist when the redirect arrives.

Before credential activation, MPC must prove:

```text
state exists, is current, unexpired and unconsumed
attempt is still the current attempt for the Installation/application
initiating Principal is still access-eligible
initiating Principal is still human
initiating Principal still has current Membership in the Organization
initiating Principal still has portfolio.manage
MarketplaceInstallation still exists/currently allows the ceremony and remains in that Organization
provider application/configuration remains compatible
provider code exchange completes sufficiently to produce a complete credential candidate
provider seller/account identity is proven authoritatively
seller binding is initial or compatible with the standing Installation binding
```

If authority was revoked after begin, the callback cannot complete. `state` never becomes eternal authorization.

---

# 18. Lane B — Provider seller/account identity proof

Provider seller/account identity is external namespace evidence, not Organization or Principal identity.

For current Mercado Livre authorization, token-returned seller/user identity is corroborated through the accepted provider-authoritative authenticated identity surface (for example `/users/me` under the selected D4 contract).

Material contradiction fails closed.

Do not use nickname, email or display name as seller identity.

## 18.1 Initial authorization

If the Installation has no standing seller binding:

```text
valid ceremony
+ provider seller S proven
→ establish Installation ↔ S binding
→ activate credential generation
```

## 18.2 Reauthorization

If the Installation is already bound to seller `S1`:

- callback proves `S1` → credentials may be renewed/replaced;
- callback proves `S2` → **fail closed**; never silently rebind the Installation.

A future deliberate seller rebind requires a separately admitted Portfolio/D4 decision or a new Installation. It is not created by OAuth convenience.

---

# 19. Lane B — Credential generation / concurrency law

Credential state changes are technical, but correctness requires a generation boundary.

> **A new credential set becomes active only after all callback, current-authority and seller-binding checks are complete. No mixed active generation is permitted.**

On failed reauthorization, the prior valid active generation remains unchanged where provider semantics permit.

A stale token-refresh result or older authorization attempt must never overwrite a newer active credential generation.

Conceptually:

```text
current credential generation G1
reauthorization establishes G2
late refresh derived from G1 returns
→ cannot replace G2
```

Exact transaction/CAS/lock/storage mechanics remain D7.

Refresh itself is outbound D4/D7 credential lifecycle, not Product operation and not inbound ingress.

---

# 20. Lane B — Replay, ambiguity and failure honesty

## 20.1 Callback replay

A successfully consumed state/attempt cannot cause a second code exchange or second credential activation. Browser refresh/replay is harmless at the semantic boundary.

## 20.2 Ambiguous code exchange

If MPC sends a code exchange but receives no complete trusted token response, it does not declare the ceremony completed and does not activate credentials.

The safe baseline is a fresh authorization ceremony rather than heuristic recovery of an externally possibly-issued orphan credential.

## 20.3 Identity read failure

Obtaining a token is insufficient if required seller/account identity cannot be established. No credential activation occurs until the D4 namespace binding proof is satisfied.

## 20.4 Provider denial/error

Provider denial, canceled consent, invalid/expired state/code and other ceremony failures change no business configuration and do not become Product business rejection vocabulary.

OAuth failures remain technical protocol outcomes.

---

# 21. Lane B — Post-callback navigation

Provider redirect URI remains fixed/provider-compliant. Any later browser redirect to the Product UI is server-controlled and derived from the Authorization Attempt, never arbitrary callback `return_url` input.

The callback redirect does not assert Product success. The UI performs a fresh Product API read such as `GetMarketplaceInstallation`, and W4 current access rules still govern the human's ability to view it.

A callback cannot bypass current Product access simply because the provider redirect succeeded.

---

# 22. Product/problem separation for OAuth

Do not add provider-protocol errors such as:

```text
oauth-state-expired
provider-code-invalid
seller-mismatch
```

to the W2 Product Problem Details catalog merely because the technical route needs failure handling.

Provider callback response/error vocabulary remains protocol-local. Product clients learn current installation/authorization posture through accepted Product resources/queries.

---

# 23. Whole Technical Ingress negative controls

Later proof must make at least these defects invalid/unreachable.

## Acquisition

1. unknown provider topic becomes generic event;
2. provider DTO reaches a D1 owner as its semantic contract;
3. provider user/application/account identity becomes Organization directly;
4. ambiguous Installation correlation selects first match;
5. arbitrary callback resource URL causes arbitrary outbound HTTP;
6. positive delivery acknowledgement is recorded as business success;
7. acknowledgement occurs before recoverable responsibility exists;
8. process crash immediately after acknowledgement permanently loses required acquisition;
9. duplicate Sale notification creates duplicate Sale meaning;
10. out-of-order payload directly regresses owner state without authoritative reread;
11. polling/recovery uses a second business processing path;
12. recovery window becomes completeness authority;
13. raw provider PII is archived by convenience;
14. unbound quarantine enters Organization Work/Sales/Offering/Fulfillment state;
15. deferred provider topics activate merely because provider supports them;
16. a second provider must fake an HTTP webhook to use the native acquisition seam.

## OAuth

17. OAuth begin without valid Product AuthN/current H + Membership + `portfolio.manage` creates an attempt;
18. ambient/cross-site GET silently creates an authorization attempt;
19. callback missing/forging/replaying/using expired state activates credentials;
20. superseded attempt activates after a newer begin;
21. callback query changes Organization/Installation authority;
22. current browser Product session substitutes for the initiating Principal's current authority;
23. initiating Principal loses access after begin but callback still activates;
24. deactivated/ineligible Installation completes authorization;
25. reauthorization to a different seller silently rebinds Installation;
26. token seller identity contradicts authoritative provider identity but credentials activate;
27. arbitrary callback return URL becomes redirect target;
28. partial/mixed credential generation becomes active;
29. late stale refresh overwrites newer reauthorization generation;
30. sensitive OAuth secrets appear in Product schema/log/problem/browser state;
31. callback creates Organization or MarketplaceInstallation;
32. provider seller identity becomes MPC Principal;
33. OAuth failure becomes Product business rejection;
34. Product bearer token authenticates provider callback or provider token authenticates Product API.

---

# 24. Whole Technical Ingress proof strategy

Architecture/D7/D8 proof must be able to falsify at least:

1. duplicate `orders_v2` → one Sale meaning;
2. positive acknowledgement + immediate crash → acquisition resumes;
3. webhook + recovery feed → same typed acquisition path;
4. zero/multiple namespace bindings → zero owner commits;
5. arbitrary resource URL → no arbitrary fetch;
6. notification payload disagrees with authoritative reread → reread/owner wins;
7. Payment reversal after release remains historically recoverable where D3/D4 requires;
8. claim + claim-action signals → same Claim acquisition without generic workflow;
9. competition signal cannot create PriceIntent;
10. queue/event-bus provider can use same native acquisition family without webhook assumptions;
11. Product user loses `portfolio.manage` between OAuth begin/callback → no activation;
12. begin B supersedes begin A → callback A cannot activate;
13. same-seller reauthorization activates a complete newer credential generation;
14. different-seller reauthorization fails closed without changing active binding;
15. late G1 refresh cannot overwrite active G2 credentials;
16. callback replay cannot produce a second activation;
17. incomplete token exchange or failed seller proof cannot activate credentials;
18. browser session identity different from initiator cannot transfer authority.

Exact executable mechanisms belong to D7/D8.

---

# 25. Reopen triggers

Reopen only the smallest implicated ingress/D4/owner decision when material evidence shows:

- a provider signal itself contains a material occurrence not recoverable through accepted authoritative reads;
- exact provider namespace/Organization correlation cannot be established with current D2/D4 bindings;
- a real Product consumer needs a deferred/rejected provider signal;
- a second provider proves an accepted typed acquisition family cannot remain provider-independent without semantic loss;
- recoverable custody cannot satisfy a real provider acknowledgement contract without changing D3/D4 recovery assumptions;
- Product 1.0 requires a business-system inbound callback not currently admitted by D4;
- a technical ingress/quarantine/authorization state becomes user-business manipulable and requires a real Product operation/owner;
- a selected provider OAuth protocol materially cannot support the accepted correlation/replay/current-authority properties;
- deliberate seller rebind becomes a real Product workflow;
- multi-actor/multi-step authorization requirements cannot be expressed by one current Authorization Attempt without moving authority.

Do not reopen for queue preference, route naming, current handler layout, provider topic availability or desire for a generic integration/OAuth framework.

---

# 26. Ingress-A outcome

**Outcome:** `RESTRUCTURE NOW — D5-LOCAL TECHNICAL INGRESS` + `DEFER SAFELY` unsupported provider modes/topics.

> **Use provider-specific inbound adapters feeding one MPC-native recoverable acquisition seam and closed typed acquisition requests. Push, polling and recovery converge on the same authoritative D4 acquisition path. External delivery never becomes Product authentication, Product operation, domain event by receipt or generic Integration authority.**

Ingress-A was explicitly operator-ratified on 2026-08-19.

---

# 27. Ingress-B outcome

**Outcome:** `CURRENT PARENT STRUCTURE CONFIRMED` + `RESTRUCTURE NOW — D5-LOCAL AUTHORIZATION CEREMONY CRYSTALLIZATION`.

> **Use an explicit Product-authorized human begin, one current server-bound Authorization Attempt, provider Authorization Code + state + selected PKCE protection, callback-time revalidation of current MPC authority, authoritative seller/account proof, initial-bind or same-seller reauthorization only, and generation-safe credential activation. OAuth state/browser/provider identity never becomes Product authority.**

Ingress-B was explicitly operator-ratified on 2026-08-19.

Technical non-Product ingress remains **OPEN / ACTIVE** until the combined A+B package passes Whole-Ingress adversarial coherence and final operator ratification.

**Exact next work:** Whole-Ingress adversarial coherence over Acquisition + Authorization as one trust-boundary system, followed by independent review if materially warranted and final canonical consolidation.

Implementation remains blocked until D9.
