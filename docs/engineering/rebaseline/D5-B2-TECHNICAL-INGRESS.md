# D5-B2 — Technical Non-Product Ingress

> **Status:** ACCEPTED / CANONICAL — Whole-Ingress operator-ratified and canonically filed  
> **Parent D5:** `D5-API.md`  
> **Operation authority:** ratified D5-B2 Operation Admission Matrix + canonical W1/W2/W3/W4  
> **External authority:** canonical D2 + D3 + D4 + D4-R1  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-19  
> **Ingress-A accepted:** 2026-08-19  
> **Ingress-B accepted:** 2026-08-19  
> **Whole-Ingress final ratification incorporated:** 2026-08-19  
> **Authored-media delivery boundary cross-reference incorporated:** 2026-08-19

## 1. Purpose, authority and stage boundary

This artifact is the single canonical D5-B2 authority for Product 1.0 inbound technical surfaces that are **not** ordinary Product Principals invoking Product API operations.

D4 remains authoritative for provider protocol, authentication, source semantics, sanctioned reads/writes and evidence coverage. This artifact crystallizes the corresponding **wire and trust boundary** at D5-B2. It does not create a second provider/source semantic authority.

Two distinct lanes are canonical:

```text
A. External Acquisition Ingress
   provider-specific discovery/notification/recovery signal
   → MPC-native recoverable marketplace acquisition
   → closed typed acquisition request
   → authoritative D4 reread/acquisition
   → consumer-owned translation
   → D1 owner commits meaning
   → D3 communication only after owner commit

B. OAuth / Authorization Ceremony
   Product-authorized human initiation
   → server-bound Authorization Attempt
   → provider authorization protocol
   → callback
   → current MPC authority revalidation
   → provider selling-account proof
   → initial bind or same-seller technical reauthorization
   → complete generation-safe credential activation
```

The lanes share an external protocol boundary but **do not share one generic event, workflow, Integration or OAuth business model**.

This artifact does not choose D7 queue/broker/database/worker/storage/secret/transaction/lock/CAS/retry/deployment realization. Provider technical routes are outside the Product API, Product OpenAPI and Product SDK. The only W4 relationship is the current-access evaluation used when a Product-authorized human explicitly begins the authorization ceremony.

Implementation remains blocked until D9.

---

## 2. Governing invariants

1. **D4 owns provider semantics.** Technical Ingress owns only D5-B2 wire/trust-boundary crystallization and bounded ingress mechanism.
2. **Provider transport is adapter-local.** HTTP webhook, SQS, EventBridge, queue, event bus, callback and provider-specific wire vocabulary remain protocol concerns.
3. **Native ingress is mechanism, not business authority.** It may centralize custody, correlation, recovery and technical provenance without owning Sale, Listing, Shipment, Payment, Claim, Market, OAuth or another D1 meaning.
4. **Provider signal != D3 domain event.** Receipt never creates owner meaning.
5. **Provider signal != current provider truth.** Current material state comes from the accepted D4 authoritative acquisition/reread contract.
6. **Organization is never payload-derived authority.** Product 1.0 marketplace ingress derives it only from an exact MPC-owned MarketplaceInstallation seller/account namespace binding.
7. **Closed typed vocabulary.** Unknown provider topics never become generic `ExternalEvent`, `WebhookEvent`, `ProviderResource` or property-bag business objects.
8. **Acquisition family != provider topic.** A family exists only when a distinct authoritative read/coverage contract is required to establish a distinct consumer-owned claim.
9. **Positive provider acknowledgement has a technical ceiling.** It rests only on recoverable attributed custody, bounded quarantine disposition or explicit terminal non-processing; it never means business acceptance/completion/convergence.
10. **Push/poll/recovery convergence is capability-qualified.** Discovery mechanisms converge only where D4 proves that the provider actually offers them for the family; the shared path never invents enumeration, polling, recovery or completeness.
11. **Delivery dedup != business idempotency.** Delivery keys/coalescing are technical optimizations; source-qualified identity and owner semantics protect correctness.
12. **Quarantine is narrow pre-attribution platform state.** It is not a feature backlog, attacker-controlled archive, Organization Work or D3 durable communication.
13. **Deactivation removes business participation, not retained attribution/history.** Same-seller technical credential restoration may remain available for evidence recovery without business reactivation or write/publication authority.
14. **Product access and provider trust are independent.** Provider credentials/signatures never authenticate Product Principals; Product bearer tokens never prove provider callbacks.
15. **OAuth state is correlation, not standing authorization.** Callback completion revalidates current MPC authority and provider selling-account identity.
16. **Initial binding requires transaction-bound anti-injection protection.** State secrecy alone is insufficient; the current Mercado Livre lane uses PKCE where supported as the selected transaction-bound control.
17. **Credential activation is generation-safe.** A stale refresh or older Authorization Attempt can never overwrite a newer active credential generation.
18. **Consuming refresh is serialized by correctness.** Current Mercado Livre refresh is single-use; one active generation therefore has one serialized refresh consumer. D7 chooses the mechanism.
19. **Historical lineage is explanation only.** Non-secret authorization/binding lineage never becomes current namespace authority; current Installation/D4 binding remains the sole current authority.
20. **No generic ingress framework by symmetry.** A later provider may reuse proven seams but keeps its actual protocol and capability evidence explicit.

---

## 3. Global Maximum

Rejected structures:

```text
Webhook business domain
Generic Integration API/domain
Generic ExternalEvent { type, payload, metadata }
provider topic = MPC domain event
one provider endpoint per D1 owner
provider-specific business-processing pipeline
universal subscription/provider-resource framework
OAuth workflow engine / generic OAuth domain
Credential CRUD Product API
Product Sync/Refresh operations for technical recovery
SourceInstance callback seam without a real admitted callback consumer
```

Selected structure:

```text
provider-specific adapter
        ↓
MPC-native MarketplaceInstallation-scoped custody/correlation seam
        ↓
closed typed acquisition request
        ↓
accepted D4 authoritative read/coverage contract
        ↓
consumer-owned semantic translation
        ↓
owner commit
        ↓
D3 communication when applicable
```

and separately:

```text
current Product-authorized human
        ↓
server-bound Authorization Attempt
        ↓
provider Authorization Code ceremony
        ↓
transaction-bound anti-injection + callback verification
        ↓
current MPC authority revalidation
        ↓
provider-authoritative selling-account proof
        ↓
complete credential-generation activation
```

**Method disposition:** `RESTRUCTURE NOW — D5-LOCAL TECHNICAL INGRESS` relative to direct provider-route→domain processing and browser/session-driven OAuth binding.

Structural Inversion remains **PASS**. No D0→W4/D3/D4 semantic parent reopen is required.

---

# Lane A — External Acquisition Ingress

## 4. Provider-specific inbound adapter

The adapter owns only protocol concerns needed to establish an admissible acquisition-signal candidate:

- HTTP/SQS/EventBridge/provider delivery mechanics actually supported by that provider;
- origin/authentication/verification evidence supported by the selected provider contract;
- provider DTO parsing;
- closed provider topic/resource grammar;
- provider application/account/seller/source evidence;
- provider delivery discriminator where actually supplied;
- translation from provider resource syntax to a bounded native acquisition key/reference.

For the current Mercado Livre HTTP notification lane, accepted official provider evidence establishes a real origin-verification basis. Concrete provider IPs or equivalent mutable network facts remain adapter-local and are revalidated; they are never frozen as MPC business constants.

The adapter does **not** choose D1 business meaning from topic names, infer Organization, commit owner state, expose raw DTOs to owners, fabricate event identity or fetch an arbitrary callback-supplied URL.

A resource URL/path is parsed against closed adapter-local grammar, converted to a native key/reference and reread only through the accepted D4 sanctioned operation.

---

## 5. MPC-native marketplace acquisition seam

Product 1.0 native acquisition is **MarketplaceInstallation-scoped**.

The seam may own only the technical mechanism/state required for:

- exact MarketplaceInstallation seller/account namespace correlation;
- Organization derivation from that MPC-owned binding;
- recoverable custody before positive external acknowledgement;
- replay/recovery bookkeeping;
- bounded pre-attribution quarantine;
- dispatch of a closed typed acquisition request;
- minimum PII-minimized technical provenance/observability needed to prove disposition, recovery and attribution.

It owns none of the business meaning acquired afterward.

Sankhya remains on its accepted embedded/outbound D4 acquisition paths. Product 1.0 does **not** fabricate a `SourceInstance` inbound callback path. A future business-system callback requires a real consumer, explicit D4 admission and its own source-qualified family/coverage contract.

### 5.1 Minimal technical context

A proportional realization may carry:

```text
organization_id                  server-derived after exact attribution
marketplace_installation_id      exactly one after attribution
mpc_received_at
provider_sent_at?                source evidence only
provider_delivery_key?           optional technical dedup evidence
native acquisition family
bounded technical provenance     only what recovery/audit requires
```

Exact schema, ID encoding and persistence remain D7. Do not create a universal payload/property bag or raw-provider archive.

### 5.2 Family admission criterion

A native acquisition family exists only when both are true:

1. a distinct authoritative D4 read/coverage contract is required; and
2. that contract establishes a distinct consumer-owned claim.

Consequences:

- one provider topic may awaken several families;
- one authoritative read may satisfy several families when its authority/coverage is sufficient;
- two provider topics may map to one family;
- another provider is never forced to imitate Mercado Livre topic/read decomposition;
- topic-name symmetry alone never admits a family.

### 5.3 Canonical Product 1.0 families

```text
AcquireMarketplaceListing
AcquireMarketplacePrice
AcquireMarketplaceSale
AcquireMarketplaceShipment
AcquireMarketplacePayment
AcquireMarketplacePostSaleClaim
AcquireMarketplaceCompetitivePosition
```

The names describe **what source evidence must be acquired**, not what the provider claims happened. They are internal technical contracts, not Product API operations, D3 events or business resources.

---

## 6. Namespace attribution, posture, credential capability and lineage

Canonical attribution is:

```text
provider application + selling-account namespace evidence
        ↓
current/retained MPC MarketplaceInstallation binding
        ↓
exactly one MarketplaceInstallation
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
historical lineage = current binding authority
```

Disposition:

- exactly one compatible retained/current binding → attribution may proceed;
- zero compatible bindings → no D1 attribution;
- multiple bindings → no D1 attribution;
- authoritative provider namespace contradiction → no D1 attribution.

Three meanings remain separate:

1. **namespace attribution/history** — which MPC Installation is bound to the provider selling-account namespace;
2. **business participation posture** — whether the Installation may currently publish/write/participate;
3. **credential usability** — whether authoritative evidence can currently be reread.

Deactivation removes business participation authority. It does not erase the retained namespace binding, historical attribution or the technical ability to restore credentials for evidence already attributable to that Installation.

A deactivated Installation with an unambiguous retained selling-account binding may undergo **same-seller technical reauthorization for evidence recovery only**, subject to all current human/access/security checks. Credential restoration:

- is not Installation reactivation;
- creates no Product operation;
- permits no publication, price, availability or other business write;
- cannot change the selling-account binding;
- may resume only attributed acquisition/reconciliation that D4 requires.

Non-secret authorization/binding lineage may preserve the smallest security provenance needed to explain who initiated which ceremony, which Installation/provider application it concerned, which selling-account proof was established and which generation resulted. It is append-only historical explanation only. Current Installation/D4 binding remains the sole current namespace authority.

Authorization code, client secret, access token, refresh token and PKCE verifier never enter that lineage.

---

## 7. Signal disposition and bounded quarantine

Every inbound signal reaches exactly one technical disposition.

### 7.1 Class A — protocol-invalid / unverified

Examples include malformed syntax, failed provider-origin verification or a resource outside closed adapter grammar.

```text
reject
no MPC custody claim
no quarantine
no owner dispatch
```

### 7.2 Class B — explicit terminal non-admission/non-processing

Applies to:

- a protocol-admissible but non-admitted topic/resource;
- an admitted-looking signal with no exact binding and no plausible legitimate binding currently pending;
- arbitrary unknown seller identities;
- multiple contradictory bindings;
- unsupported future feature signals.

```text
record only the bounded terminal disposition/provenance required for protocol honesty
no recoverable feature backlog
no indefinite payload retention
no owner dispatch
```

### 7.3 Class C — admitted and exactly attributed

```text
recoverable Organization-scoped custody
→ typed acquisition
→ authoritative D4 reread when credential capability permits
```

If credentials are unavailable/auth-invalid, the acquisition remains attributed, detectable and recoverable. It may wait for same-seller technical credential restoration; it never silently stalls forever. D6/D7 later choose bounded observability/operator treatment without transferring source truth to Work.

### 7.4 Class D — bounded pre-attribution quarantine

Quarantine is allowed **only** when all are true:

- the topic/resource family is admitted;
- exact attribution is temporarily unresolved;
- a legitimate binding is plausibly pending for a known current MarketplaceInstallation + provider application;
- retention/capacity/reprocessing bounds are enforceable.

Quarantine is:

- platform-scoped pre-attribution protocol state;
- outside D3's Organization-scoped durable communication/recovery class;
- unable to dispatch owner acquisition;
- PII-minimized;
- bounded in capacity and time;
- observable when capacity/retention is exhausted;
- refused/fail-honestly on overflow, never silently discarded.

Once exact Installation attribution exists, the signal leaves quarantine and becomes Organization-scoped recoverable acquisition. It does not remain a second queue/authority.

Do not create Work from unbound quarantine. Storage, retention, alerting and operator tooling remain D7/D6 as justified later.

---

## 8. Positive acknowledgement and custody law

Positive provider acknowledgement is allowed only after one explicit technical decision:

1. **recoverable attributed custody** has been established;
2. a **bounded quarantine disposition** has been established; or
3. an **explicit terminal non-admission/non-processing disposition** has been established and the provider protocol should be acknowledged.

There is no silent-discard fourth basis.

```text
HTTP 200 / queue delete / equivalent acknowledgement
!= Sale committed
!= Listing updated
!= Payment processed
!= Claim resolved
!= business accepted
!= externally converged
```

If a crash immediately after positive acknowledgement can permanently lose a required acquisition or its explicit terminal disposition, the realization violates this contract.

D7 chooses the custody mechanism.

---

## 9. Push, poll and recovery: one path, explicit family coverage

The convergence shape is conditional on real provider capability:

```text
supported live notification ─────────┐
supported missed-feed/recovery ──────┤
supported reconciliation/enumeration ├→ same family acquisition request
cold-start/source scan when proven ──┤
manual technical recovery ───────────┘
```

A shared path never proves that every branch exists. Push improves latency; it never becomes completeness authority. Provider retention, recovery-window exhaustion and unsupported enumeration remain honest partial/unknown coverage under D4/W3.

Current Mercado Livre family obligations are:

| Family | Current trigger/read anchor | Coverage/recovery statement |
|---|---|---|
| `AcquireMarketplaceListing` | `items` → source-qualified Item/Listing reread | current Listing evidence follows the accepted D4 read contract; notification completeness is not claimed |
| `AcquireMarketplacePrice` | `items_prices` → selected current price read | establishes current qualified price evidence only; no historical occurrence stream is invented |
| `AcquireMarketplaceSale` | `orders_v2` → source-qualified Order point reread | a cancellation-inclusive recovery universe is **not proven**; seller Order enumeration does not currently discharge it; D8 proof obligation |
| `AcquireMarketplaceShipment` | `shipments` → source-qualified Shipment reread | coverage follows the selected D4 Shipment read contract; no universal provider completeness claim |
| `AcquireMarketplacePayment` | `payments` → accepted Payment read anchored from accepted D4 identity relations | material refund/reversal recovery is only what D4's authoritative occurrence/anchor evidence supports; delivery receipt is not payment truth |
| `AcquireMarketplacePostSaleClaim` | `post_purchase:claims` and `post_purchase:claims_actions` → source-qualified Claim reread | both topics awaken one family; whether every material Claim-action occurrence is reconstructable remains **Unknown**, a D4/D8 proof obligation |
| `AcquireMarketplaceCompetitivePosition` | `catalog_item_competition_status` → selected competition/price-to-win read | current provider-qualified competitive evidence only; provider population is not universal market completeness |

Internal D3 owner events do not loop through this ingress for transport symmetry.

---

## 10. Duplicate, replay, ordering and coalescing

```text
provider delivery dedup
!= technical acquisition coalescing
!= owner semantic idempotency/current truth
```

Provider delivery keys may reduce duplicate work but never become Sale/Listing/Shipment/Payment/Claim identity.

Current-state invalidation signals may be coalesced only when the authoritative D4 reread plus owner claim loses no material information. Material occurrence evidence may be coalesced only when authoritative acquisition still recovers every occurrence required by D3 evidence-edge correctness.

If latest reread cannot recover a materially required external occurrence, D4 must explicitly preserve/classify that occurrence rather than inventing universal webhook event sourcing.

Out-of-order provider payload never directly regresses owner state. Authoritative reread and owner semantics decide current/historical meaning.

---

## 11. Mercado Livre ingress admission matrix

Provider topic spelling remains adapter-local evidence. Admission freezes MPC treatment, not provider ontology.

### 11.1 ADMIT

| Provider signal/topic | Native acquisition | Authoritative D4 acquisition | Consumer meaning |
|---|---|---|---|
| `orders_v2` | `AcquireMarketplaceSale` | source-qualified Order point reread | Marketplace Sales |
| `shipments` | `AcquireMarketplaceShipment` | source-qualified Shipment reread | Fulfillment |
| `items` | `AcquireMarketplaceListing` | source-qualified Item/Listing reread | Offering; admitted evidence may feed Readiness/Availability |
| `items_prices` | `AcquireMarketplacePrice` | selected current price read | Offering/PriceIntent observation; qualified Economics input |
| `payments` | `AcquireMarketplacePayment` | accepted D4 Payment read | Commercial Economics + Post-Sale evidence |
| `post_purchase:claims` | `AcquireMarketplacePostSaleClaim` | source-qualified Claim reread | Post-Sale Resolution |
| `post_purchase:claims_actions` | same Claim acquisition | same Claim reread | Post-Sale Resolution; action-occurrence completeness remains Unknown |
| `catalog_item_competition_status` | `AcquireMarketplaceCompetitivePosition` | selected competition/price-to-win read | Market Intelligence |

The two post-purchase topics intentionally map to one acquisition family. `items` may awaken multiple consumer-owned family/semantic paths only when each required D4 read/coverage contract is satisfied; it never creates a provider-resource owner.

### 11.2 DEFER SAFELY

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

A DEFER row is not pre-authorization to subscribe, retain or implement.

### 11.3 REJECT BASELINE

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

## 12. Product API, OpenAPI, SDK and route separation

Technical ingress:

- is outside `/organizations/{organization_id}/...` Product roots;
- is not `/access-context`;
- is not one of the 95 Product operations;
- is not included in the Product SDK;
- is not itself authorized by W4 Product Permission mapping;
- never uses provider credentials to impersonate Product Principals;
- never uses Product bearer tokens as provider callback proof;
- may use provider-specific route/transport vocabulary inside the protocol boundary.

The single Product OpenAPI remains the Semantic Product API wire authority. Provider protocol conformance and internal typed ingress contracts are separate technical contracts, not parallel Product business wire authorities.

Exact provider-facing host/prefix/method/redirect spelling remains `DEFER SAFELY` until final wire/runtime realization. Any future spelling must be unambiguously separate and must not collide semantically or syntactically with `/access-context` or `/organizations/{organization_id}/...`.

Do not create generic Product `/providers`, `/integrations`, `/webhooks`, `/oauth` or `/external-events` roots.

### 12.1 Authored-media byte delivery is not ingress

Authored-media byte delivery is a **separately justified technical presentation surface** owned by canonical W2 §3.9.8. It is:

- not Lane A external acquisition ingress;
- not Lane B authorization ceremony;
- not a Product operation, and not part of the Product OpenAPI/SDK;
- authorized by current Product AuthN + unique Principal binding + Principal access eligibility + Organization Membership + `offering.read` for the exact ListingIntent/media relationship — never by provider credentials, callback origin or a durable anonymous locator.

Its failures are technical-surface failures and never enter the W2 Product Problem catalog. This artifact does not own that surface; it is named here only so no future reader files a presentation surface into an ingress lane.

---

# Lane B — OAuth / Authorization Ceremony

## 13. Boundary and initiation authority

OAuth begin/callback/reauthorization is **not** an acquisition signal and is **not** a Product business operation.

Its purpose is to establish or renew the provider credential plus selling-account namespace proof of an existing MarketplaceInstallation without allowing browser session, provider operator identity or OAuth protocol state to become Product authority.

Only an authenticated **human** may initiate the current Mercado Livre ceremony.

At begin time all must be current:

```text
valid Product AuthN
exactly one MPC Principal binding
Principal.kind = human
current Principal access eligibility
current Organization Membership
portfolio.manage
exact MarketplaceInstallation belongs to that Organization
Installation is eligible for the requested ceremony purpose
```

This is a **W4-equivalent current-access evaluation** reused by a technical route. It does not add a 96th Product operation, a 30th Permission, a new Principal class or Product OpenAPI/SDK surface.

For an active Installation, the ceremony may establish initial binding or renew same-seller credentials. For a deactivated Installation, eligibility is limited to same-seller technical evidence-recovery restoration under §6; business participation remains disabled.

The begin is an explicit state-creating technical action. Ambient/cross-site GET navigation must not silently create an Authorization Attempt. Exact route/method/redirect status remain deferred technical wire spelling.

---

## 14. Server-bound Authorization Attempt

One bounded technical Authorization Attempt correlates the ceremony. It is not a Product resource, durable business workflow or current authorization authority.

It binds proportionately:

```text
provider application/client
Organization
MarketplaceInstallation
initiating MPC Principal
ceremony purpose: initial/same-seller active renewal/evidence-recovery renewal
opaque state
transaction-bound anti-injection binding
PKCE verifier/challenge when selected
created/expiry boundary
server-controlled safe post-callback destination
technical attempt generation
```

Exact storage/ID/TTL mechanics remain D7.

### 14.1 State law

`state` is:

```text
high entropy
opaque
single-use
finite-lived
bound to exactly one Authorization Attempt
```

It is correlation/anti-forgery evidence only. It never proves current Membership, `portfolio.manage`, Installation eligibility, seller compatibility or authority to complete after those facts change.

The provider redirect URI remains static/provider-registered. Organization, Installation, ceremony purpose and dynamic return targets are not trusted from callback query parameters.

### 14.2 One current attempt

For one MarketplaceInstallation + provider application, at most one unfinished Authorization Attempt is current.

A new authorized begin supersedes the prior unfinished attempt. A callback for superseded state cannot exchange code or activate credentials.

### 14.3 Transaction-bound anti-CSRF / code-injection protection

Initial binding may not rely on state secrecy alone.

Every provider lane requires one explicit transaction-bound control that ties the authorization response/code to the initiating transaction or user-agent.

For current Mercado Livre, PKCE is selected where supported and carries this proof together with opaque state and the other load-bearing controls. A future provider without usable PKCE must bind state securely to the initiating user-agent session or require equivalent explicit human confirmation before initial bind activation. MPC does not invent a universal OAuth capability to hide provider differences.

---

## 15. Current Mercado Livre authorization protocol

The current lane uses:

```text
Authorization Code
+ statically registered redirect URI
+ opaque single-use finite-lived server-bound state
+ PKCE where supported/selected
+ server-to-server code exchange using the confidential client
+ callback-time current-authority revalidation
+ authoritative selling-account proof
```

PKCE is the selected provider-lane hardening/transaction-binding control; the whole security claim also depends on the static redirect, server-side attempt, confidential exchange, current-authority checks and selling-account proof.

Provider/client secret, authorization code, PKCE verifier, access token and refresh token are sensitive technical credentials and must not enter Product schemas, browser storage, normal logs, Problem Details, business history or non-secret lineage.

Current Mercado Livre authorization guidance requires an **administrator** of the selling account. Collaborator/operator authorization is not accepted as selling-account binding proof in the current lane.

A future provider with delegated operators requires explicit D4 proof that the authenticated operator can be mapped authoritatively to the selling-account namespace. It is not inferred from nickname, email, display name or provider role text.

---

## 16. Callback current-authority revalidation

The callback is not authorized by whatever Product browser session happens to exist when the redirect arrives.

Before credential activation, MPC must prove:

```text
state exists, is current, unexpired and unconsumed
attempt is still current for the Installation/application
transaction-bound anti-injection proof succeeds
initiating Principal still resolves uniquely and is access-eligible
initiating Principal is still human
initiating Principal still has current Membership in the Organization
initiating Principal still has portfolio.manage
MarketplaceInstallation still exists, remains in that Organization and allows the recorded ceremony purpose
provider application/configuration remains compatible
provider code exchange yields one complete credential candidate
provider selling-account namespace is proven authoritatively
selling-account binding is initial or exactly compatible with the retained/current Installation binding
```

If authority was revoked after begin, callback completion fails closed. State never becomes eternal authorization.

---

## 17. Selling-account binding proof

The Installation binding subject is the provider **selling-account namespace**, not an arbitrary authorizing operator Principal.

For current Mercado Livre, token-returned identity is corroborated through the accepted provider-authoritative authenticated identity surface under the selected D4 contract. Material contradiction fails closed.

Do not use nickname, email or display name as seller identity.

### 17.1 Initial authorization

If the Installation has no standing seller binding:

```text
valid ceremony
+ administrator authorization
+ selling account S proven
→ establish Installation ↔ S binding
→ activate one complete credential generation
```

### 17.2 Same-seller reauthorization

If the Installation is bound to `S1`:

- callback proves `S1` → credentials may be renewed/replaced for the admitted active or evidence-recovery purpose;
- callback proves `S2` → **fail closed**; never silently rebind the Installation.

A future deliberate seller rebind requires a separately admitted Portfolio/D4 decision or a new Installation. It is not created by OAuth convenience.

Collaborator/operator credentials are outside the current Mercado Livre lane even when they could act on behalf of the seller; administrator proof remains required.

---

## 18. Credential generation and consuming refresh

A new credential set becomes active only after all callback, current-authority, anti-injection and selling-account checks are complete. No partial/mixed active generation is permitted.

On failed reauthorization, the prior valid active generation remains unchanged where provider semantics permit.

A stale refresh result or older Authorization Attempt can never overwrite a newer generation:

```text
current generation G1
reauthorization activates G2
late refresh derived from G1 returns
→ cannot replace G2
```

Current official Mercado Livre evidence establishes:

- only the latest refresh token is accepted; and
- each refresh token is single-use/consuming.

Therefore refresh for one active credential generation must be **serialized/single-consumer by correctness**, not merely for performance. A concurrent loser must not destroy the only usable generation, and no stale result may replace a newer OAuth or refresh generation.

D7 chooses transaction/CAS/lock/lease/storage/secret mechanisms. Refresh remains outbound D4/D7 credential lifecycle, not Product operation and not inbound acquisition.

---

## 19. Cross-lane recovery after credential activation

Successful initial/same-seller credential activation may awaken only bounded technical work that the new credential capability makes possible, such as:

- initial provider capability/bootstrap reads;
- pending attributed acquisition recovery;
- bounded reprocessing of signals that leave legitimate pre-attribution quarantine after binding becomes exact;
- source reconciliation already required by D4.

It never creates:

```text
Product Sync/Refresh operation
D3 business event merely for OAuth success
business reactivation
publication/write authorization
owner semantic success
```

Owner meaning still begins only after authoritative D4 acquisition and owner commit.

---

## 20. Replay, ambiguity, navigation and failure honesty

### 20.1 Callback replay

A consumed state/attempt cannot cause a second code exchange or credential activation. Browser refresh/replay is harmless at the semantic boundary.

### 20.2 Ambiguous code exchange

If MPC sends a code exchange but receives no complete trusted token response, it does not declare completion or activate credentials. The safe baseline is a fresh ceremony rather than heuristic recovery of an externally possibly-issued orphan credential.

### 20.3 Identity read failure

Obtaining a token is insufficient if the selling-account namespace cannot be established. No activation occurs until the D4 binding proof is satisfied.

### 20.4 Provider denial/error

Provider denial, canceled consent, invalid/expired state/code and other ceremony failures change no business configuration and do not become Product business rejection vocabulary.

### 20.5 Post-callback navigation

Any later browser redirect to Product UI is server-controlled and derived from the Authorization Attempt, never arbitrary callback `return_url` input.

The redirect does not assert Product success. The UI performs a fresh Product API read such as `GetMarketplaceInstallation`; current W4 rules still govern visibility.

### 20.6 Product Problem separation

Do not add provider-protocol errors such as:

```text
oauth-state-expired
provider-code-invalid
seller-mismatch
provider-origin-invalid
```

to the W2 Product Problem Details catalog merely because technical routes need failure handling.

Provider callback/notification response vocabulary remains protocol-local. Product clients learn current Installation/authorization posture through accepted Product resources/queries.

---

## 21. Whole Technical Ingress negative controls

Later D7/D8/OpenAPI/runtime proof must make at least these defects invalid or unreachable.

### Acquisition

1. unknown provider topic becomes generic event/resource;
2. provider DTO reaches a D1 owner as semantic contract;
3. provider user/application/account identity becomes Organization directly;
4. Product 1.0 fabricates a SourceInstance/Sankhya inbound path;
5. ambiguous Installation correlation selects first match;
6. arbitrary callback resource URL causes arbitrary outbound HTTP;
7. protocol-invalid input gains custody/quarantine;
8. arbitrary unknown seller identity gains indefinite quarantine;
9. quarantine becomes feature backlog, attacker payload sink, Organization Work or D3 communication;
10. positive acknowledgement occurs without one of the three admitted bases;
11. positive acknowledgement is recorded as business success;
12. crash after acknowledgement permanently loses required acquisition/disposition;
13. duplicate Sale notification creates duplicate Sale meaning;
14. out-of-order payload directly regresses owner state without authoritative reread;
15. shared ingress path invents provider polling/enumeration/recovery/completeness;
16. `orders_v2` is claimed cancellation-complete without D8 proof;
17. Claim-action occurrence completeness is claimed without D4/D8 proof;
18. raw provider PII is archived by convenience;
19. deferred topic activates merely because provider supports it;
20. a second provider must fake an HTTP webhook or Mercado Livre topic decomposition.

### OAuth / credentials

21. begin creates a 96th Product operation or new Permission;
22. begin succeeds without current authenticated H + eligibility + Membership + `portfolio.manage` + exact Installation;
23. ambient/cross-site GET silently creates an Authorization Attempt;
24. state secrecy alone authorizes initial binding without transaction-bound anti-injection proof;
25. missing/forged/replayed/expired/superseded state activates credentials;
26. callback query changes Organization, Installation, purpose or return authority;
27. current browser Product session substitutes for the initiator's current authority;
28. initiator loses access after begin but callback still activates;
29. deactivated Installation gains business participation/write authority through evidence-recovery reauthorization;
30. different-seller reauthorization silently rebinds Installation;
31. collaborator/operator identity is accepted as current Mercado Livre selling-account proof;
32. token seller identity contradicts authoritative provider identity but activation proceeds;
33. partial/mixed credential generation becomes active;
34. two concurrent consuming refreshes can invalidate the active generation;
35. stale G1 refresh overwrites G2;
36. secrets appear in Product schema/log/problem/browser/history/lineage;
37. historical lineage becomes a second current binding authority;
38. callback creates Organization or MarketplaceInstallation;
39. provider seller identity becomes MPC Principal;
40. OAuth failure becomes Product business rejection or D3 business event;
41. successful OAuth creates Product Sync/Refresh or business reactivation;
42. Product bearer token authenticates provider callback or provider credential authenticates Product API;
43. technical route collides with `/access-context` or `/organizations/{organization_id}/...`;
44. technical ingress appears in Product OpenAPI/SDK.

---

## 22. Proof strategy and explicit residuals

D7/D8 proof must be able to falsify at least:

1. duplicate `orders_v2` → one Sale meaning;
2. positive acknowledgement + immediate crash → acquisition/disposition resumes;
3. notification + any proven recovery feed → same family acquisition path;
4. zero/multiple namespace bindings → zero owner commits and no unbounded quarantine;
5. legitimate pending binding → bounded quarantine → exact attribution → one recoverable acquisition;
6. arbitrary unknown seller → explicit terminal/refusal, never attacker-controlled storage;
7. arbitrary resource URL → no arbitrary fetch;
8. notification payload disagrees with authoritative reread → reread/owner wins;
9. Payment reversal after release remains recoverable exactly where D3/D4 requires;
10. claim + claim-action signals → one Claim family without generic workflow;
11. material Claim-action occurrence that cannot be reconstructed remains detectably Unknown rather than silently complete;
12. Sales recovery test exposes the current cancellation-inclusive gap until proved;
13. queue/event-bus provider can use the same family criterion without webhook assumptions;
14. user loses `portfolio.manage` between begin/callback → no activation;
15. begin B supersedes A → callback A cannot activate;
16. initial binding without valid transaction-bound anti-injection proof → no activation;
17. same-seller administrator reauthorization activates one complete newer generation;
18. different-seller or collaborator authorization fails closed without changing binding;
19. deactivated same-seller evidence-recovery reauthorization restores read capability but not business participation/write authority;
20. two concurrent current-Mercado-Livre refresh attempts cannot consume/destroy the generation incorrectly;
21. late G1 result cannot overwrite active G2;
22. callback replay cannot produce second activation;
23. incomplete exchange or failed selling-account proof cannot activate;
24. historical lineage cannot determine current binding;
25. successful authorization may wake bounded technical recovery but no Product operation/D3 event.

Explicit open proof obligations:

- `AcquireMarketplaceSale` / `orders_v2`: cancellation-inclusive recovery universe remains unproven;
- `post_purchase:claims_actions`: material occurrence reconstructability remains Unknown;
- exact provider-facing host/prefix/method and D7 custody/quarantine/credential/refresh mechanisms remain deferred to their proper stages.

---

## 23. Reopen triggers

Reopen only the smallest implicated ingress/D4/owner decision when material evidence shows:

- a provider signal contains a material occurrence not recoverable through accepted authoritative reads;
- exact provider selling-account/Organization correlation cannot be established with current D2/D4 bindings;
- a real Product consumer needs a deferred/rejected provider signal;
- a second provider proves a family cannot remain provider-independent without semantic loss;
- recoverable custody cannot satisfy a real acknowledgement contract without changing D3/D4 assumptions;
- Product 1.0 requires a real business-system inbound callback;
- technical ingress/quarantine/authorization state becomes user-business manipulable and requires a real Product owner/operation;
- a selected authorization protocol cannot satisfy correlation/replay/current-authority/anti-injection properties;
- deliberate selling-account rebind becomes a real Product workflow;
- delegated-operator authorization becomes required and can be proven to a selling-account namespace;
- multi-actor/multi-step authorization cannot be expressed by one bounded Authorization Attempt;
- official provider evidence changes the consuming-refresh or administrator requirements.

Do not reopen for queue preference, route naming, current handler layout, provider topic availability, framework symmetry or desire for a generic integration/OAuth platform.

---

## 24. Canonical Whole-Ingress outcome

### Lane A

**Outcome:** `RESTRUCTURE NOW — D5-LOCAL TECHNICAL INGRESS` + `DEFER SAFELY` unsupported provider modes/topics.

> **Use provider-specific adapters feeding one MarketplaceInstallation-scoped MPC-native recoverable acquisition seam and seven closed typed acquisition families. Admit custody, quarantine or terminal non-processing explicitly; then reread through D4 and let consumer owners commit meaning. Receipt, acknowledgement and provider vocabulary never become Product authentication, Product operation, D3 event by receipt or generic Integration authority.**

### Lane B

**Outcome:** `CURRENT PARENT STRUCTURE CONFIRMED` + `RESTRUCTURE NOW — D5-LOCAL AUTHORIZATION CEREMONY CRYSTALLIZATION`.

> **Use an explicit current Product-authorized human begin, one current server-bound Authorization Attempt, transaction-bound anti-injection protection, provider Authorization Code protocol, callback-time current MPC authority revalidation, authoritative selling-account proof, initial or same-seller reauthorization only, and generation-safe serialized consuming refresh. OAuth state/browser/provider operator identity never becomes Product authority.**

### Final Whole-Ingress closure

```text
Ingress-A External Acquisition                         ACCEPTED / CANONICAL
Ingress-B OAuth / Authorization Ceremony               ACCEPTED / CANONICAL
Whole-Ingress lead review                              COMPLETE
Fable independent review                              COMPLETE
GPT final adjudication                                CONVERGED
TI-C1…TI-C13                                           INCORPORATED
Round 2                                                NOT REQUIRED
operator final ratification                            COMPLETE / FILED
W1 / Operation Matrix / W4 semantic reopen             NONE
D0→D4 / D4-R1 / D5-B1 parent semantic reopen           NONE
Product operations                                     95 unchanged
ordinary Permissions                                   29 unchanged
Technical Ingress                                      ACCEPTED / CANONICAL
```

The Whole-Ingress review candidate is not an active authority and is removed from the active tree; Git history remains the archive.

**Current status and exact next work are owned only by the router.**

Implementation remains blocked until D9.
