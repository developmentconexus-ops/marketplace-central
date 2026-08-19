# D5-B2 — Technical Non-Product Ingress

> **Status:** OPEN / ACTIVE — A External Acquisition Ingress **ACCEPTED IN-STAGE / OPERATOR-RATIFIED**; B OAuth / Authorization Ceremony = NEXT  
> **Parent D5:** `D5-API.md`  
> **Operation authority:** ratified D5-B2 Operation Admission Matrix + canonical W1/W2/W3/W4  
> **External authority:** canonical D3 + D4 + D4-R1  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-19  
> **Ingress-A accepted:** 2026-08-19

## 1. Purpose and stage boundary

This artifact classifies Product 1.0 inbound technical surfaces that are **not** a Product Principal invoking a Product API operation.

It exists to close a gap between accepted D4 provider acquisition semantics and the D5 wire boundary:

> **provider-specific inbound transport must not connect directly to D1 business processing by protocol shape. MPC needs a native technical acquisition seam that assumes recoverable custody, establishes namespace/Organization correlation and dispatches a closed typed acquisition request without becoming business authority.**

This artifact does **not** create a generic Integration/Ingress business domain, public event API, provider-resource business model or generic connector platform.

It does **not** choose D7 queue/broker/database/worker/storage/secret/retry/deployment realization.

It does **not** place provider callbacks/webhooks in the Product OpenAPI/SDK or W4 Product Permission surface.

Implementation remains blocked until D9.

---

# 2. Governing invariants

1. **External signal transport is adapter-local.** HTTP webhook, SQS, EventBridge or another provider delivery mechanism remains provider protocol.
2. **Native ingress is technical mechanism, not semantic authority.** It may centralize custody, correlation, replay/recovery and observability without owning Sale, Listing, Shipment, Payment, Claim, Market or other business meaning.
3. **Provider signal != D3 domain event.** A provider notification becomes only an acquisition trigger unless D4 explicitly proves that the external occurrence itself is non-reconstructable material evidence.
4. **Provider signal != current provider truth.** Current material meaning comes from the accepted D4 authoritative acquisition/reread contract.
5. **Organization is derived only from an exact current MPC-owned namespace binding.** Provider `user_id`, application ID, seller/account ID, SourceInstance marker, topic or payload field never becomes tenant authority.
6. **Closed typed acquisition vocabulary.** Unknown provider topics do not become generic `ExternalEvent` objects.
7. **Positive provider acknowledgement means technical recoverable custody only.** It never means owner/business acceptance, application, completion or convergence.
8. **Push/poll/recovery converge.** Webhook, missed-feed, scheduled reconciliation, cold-start scan and manual technical recovery may discover work differently but dispatch the same typed acquisition path for the same semantic acquisition.
9. **Delivery dedup != business idempotency.** Provider delivery IDs and technical coalescing may reduce work; correctness remains protected by source-qualified identity and owner semantics.
10. **PII/raw payload minimization remains binding.** Provider payload is not archived wholesale merely because the ingress received it.
11. **Product API remains separate.** Provider credentials/signatures never authenticate Product Principals; Product bearer tokens/Permissions never authenticate provider ingress.
12. **OAuth/authorization ceremony is a separate technical lane.** It shares the protocol boundary but does not enter the acquisition-signal inbox.

---

# 3. Global Maximum selected

Rejected structures:

```text
Webhook business domain
Generic Integration API
Generic ExternalEvent { type, payload, metadata }
provider topic = MPC domain event
one webhook endpoint per D1 owner
one provider-specific business processing pipeline per marketplace
universal subscription/provider-resource framework
```

Selected structure:

```text
external provider signal transport
        ↓
provider-specific inbound adapter
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
D3 communication only after owner commit when applicable
```

**Method disposition:** `RESTRUCTURE NOW — D5-LOCAL TECHNICAL INGRESS` relative to a direct webhook→domain interpretation.

No D0→W4 semantic parent reopen is required.

---

# 4. Provider-specific inbound adapter

The provider-specific inbound adapter owns only protocol concerns needed to establish an admissible signal candidate, including proportionately:

- HTTP/SQS/EventBridge/provider delivery mechanics;
- transport admission/origin/provider verification supported by the selected provider contract;
- provider DTO parsing;
- closed provider topic/resource grammar;
- provider application/account/seller/source evidence;
- provider delivery discriminator where the provider actually supplies one;
- translation from provider topic/resource syntax into a bounded native acquisition candidate.

The adapter does **not**:

- choose a D1 owner by provider topic naming;
- commit Sale/Listing/Shipment/Payment/Claim meaning;
- trust arbitrary provider `resource` URLs for outbound fetch;
- infer Organization from provider payload;
- expose raw provider DTOs to business owners;
- invent a provider-global `EventId` when the protocol has none.

A provider resource path/URL is parsed against a closed adapter-local grammar and translated into a source-qualified native key/reference. The outbound adapter then performs the accepted sanctioned D4 read; MPC never performs arbitrary callback-provided URL fetches.

---

# 5. MPC-native external acquisition ingress

The native ingress is a small technical seam between verified provider protocol and consumer-owned acquisition.

It may own technical state/mechanism required to satisfy:

- exact current MarketplaceInstallation/SourceInstance correlation;
- Organization derivation from that MPC-owned binding;
- receipt/custody before positive external acknowledgement;
- technical replay/recovery state;
- bounded quarantine for signals that cannot yet be attributed safely;
- dispatch of a closed typed acquisition request;
- technical provenance/observability needed to prove recovery and attribution.

It owns **none** of the business meaning acquired after reread.

### 5.1 Minimal technical context

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

Exact schema/encoding/persistence remains D7.

Do not create a universal provider payload/property bag merely to make future signals fit.

### 5.2 Typed acquisition requests

The accepted Product 1.0 acquisition families are:

```text
AcquireMarketplaceListing
AcquireMarketplacePrice
AcquireMarketplaceSale
AcquireMarketplaceShipment
AcquireMarketplacePayment
AcquireMarketplacePostSaleClaim
AcquireMarketplaceCompetitivePosition
```

These names denote **what current source evidence must be acquired**, not what the provider claims happened.

They do not become Product API operations and require no Product Permission mapping.

---

# 6. Exact namespace / Organization correlation

Provider signal attribution follows:

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
process-global/default Organization
```

Outcomes:

- exactly one current compatible binding → native acquisition may proceed;
- zero bindings → no D1 owner attribution;
- multiple bindings → no D1 owner attribution;
- authoritative namespace contradiction with an existing binding → no D1 owner attribution.

A signal without exact attribution may enter bounded **technical quarantine** only under §7; it cannot become Organization Work/Sales/Offering/Fulfillment state.

---

# 7. Technical quarantine

Technical quarantine is admitted because provider retry windows are bounded and silent discard would violate D3/D4 recovery expectations.

Quarantine is:

- technical/platform state, not a D1 business resource;
- not Product API-visible by default;
- not Organization-attributed until exact binding is established;
- PII-minimized;
- unable to dispatch any D1 owner acquisition while attribution is unresolved;
- recoverable/reprocessable after an explicit legitimate correlation correction.

Typical causes:

```text
zero current bindings
multiple current bindings
provider namespace contradiction
known provider signal whose currently selected ingress mapping is unavailable
```

Do not create Work from an unbound technical signal because Work requires legitimate source-domain/Organization meaning.

Exact quarantine storage, retention, operator tooling and alerting remain D7/D6 as later justified.

---

# 8. Positive acknowledgement / custody law

Provider transport acknowledgements are adapter-local protocol responses, but D5 freezes their semantic ceiling:

> **A positive ingress acknowledgement may be emitted only after MPC has established recoverable responsibility for the admitted signal or an explicit quarantine disposition.**

Therefore:

```text
HTTP 200 / SQS delete / equivalent provider acknowledgement
!= Sale committed
!= Listing updated
!= Payment processed
!= business accepted
!= externally converged
```

If the selected realization cannot recover after process crash immediately following positive acknowledgement, the realization fails this contract.

D7 chooses how to establish recoverable custody; this artifact does not mandate Postgres, broker, queue or another storage/runtime topology.

---

# 9. Push / poll / recovery convergence

The same acquisition family is used regardless of discovery source:

```text
live provider notification ───────┐
provider missed-feed/recovery ────┤
scheduled reconciliation/poll ────┼→ same typed acquisition request
cold-start/source scan ────────────┤
manual technical recovery ────────┘
```

A push source therefore improves latency/freshness but never becomes completeness authority.

Notification outage, replay, missed-feed exhaustion or provider retention limits preserve honest partial/unknown coverage under D4/W3; they do not authorize a stronger completeness claim.

Internal D3 events do **not** loop through this surface merely for transport symmetry. In-process/outbox/queue realization of owner-committed MPC events remains D3/D7 internal communication.

---

# 10. Duplicate, replay, ordering and coalescing

Three distinct layers remain separate:

```text
provider delivery dedup
    != technical acquisition coalescing
    != owner semantic idempotency / current truth
```

A provider delivery key may be used only when supplied by the provider and may not become a business identity.

For current-state invalidation-style signals, D7 may coalesce pending acquisition requests only when the selected D4/owner correctness claim cannot lose material information.

For material occurrence evidence such as financial reversals/refunds or post-sale occurrences, coalescing is permitted only when authoritative acquisition still recovers every occurrence material to D3 evidence-edge correctness.

If latest reread cannot recover a materially required external occurrence, D4 must explicitly classify/preserve the occurrence itself. Do not silently invent universal webhook event sourcing.

---

# 11. Mercado Livre Ingress Admission Matrix — Product 1.0

Current Mercado Livre topic/provider spelling remains adapter-local protocol evidence. Admission below freezes MPC treatment, not a permanent provider taxonomy.

## 11.1 ADMIT

| Provider signal/topic | Native MPC acquisition | Authoritative D4 acquisition | Product consumer / owner meaning | Classification |
|---|---|---|---|---|
| `orders_v2` | `AcquireMarketplaceSale` | source-qualified Order point reread | Marketplace Sales | current-state acquisition hint |
| `shipments` | `AcquireMarketplaceShipment` | source-qualified Shipment point reread | Fulfillment | current-state acquisition hint |
| `items` | `AcquireMarketplaceListing` | source-qualified Item/Listing reread | Offering; translated evidence may feed accepted Readiness/Availability ports | provider-resource invalidation |
| `items_prices` | `AcquireMarketplacePrice` | selected current price observation/read contract | Offering/PriceIntent observation; qualified Economics input where accepted | price-state invalidation |
| `payments` | `AcquireMarketplacePayment` | accepted D4 Payment read contract | Commercial Economics + Post-Sale evidence consumers | financial evidence acquisition |
| `post_purchase:claims` | `AcquireMarketplacePostSaleClaim` | accepted source-qualified Claim reread | Post-Sale Resolution | claim-state acquisition |
| `post_purchase:claims_actions` | `AcquireMarketplacePostSaleClaim` | same Claim reread | Post-Sale Resolution | additional trigger; no new MPC meaning |
| `catalog_item_competition_status` | `AcquireMarketplaceCompetitivePosition` | selected competition/price-to-win evidence read | Market Intelligence | market-position invalidation |

The two post-purchase provider signals intentionally feed one native acquisition family; provider action vocabulary does not become Product operation vocabulary.

`items` may support one provider acquisition feeding multiple accepted consumer-owned ports. That reuse does not create a generic provider-resource owner and does not merge Offering/Readiness/Availability authority.

## 11.2 DEFER SAFELY

| Provider signal/topic | Disposition | Reopen trigger |
|---|---|---|
| `stock-location` | DEFER | real multi-origin/User-Product stock lane becomes accepted Product scope |
| `fbm_stock_operations` | DEFER | Full/FBM observation/control becomes a selected Product flow |
| `flex-handshakes` | DEFER | Mercado Envios Flex becomes a selected Fulfillment lane |
| `invoices` | DEFER | provider-native invoice generation becomes a selected Materialization/fiscal evidence lane |
| `user-products-families` | DEFER | proactive family/relation change detection becomes a real Readiness/Offering correctness consumer |
| `catalog_suggestions` | DEFER | provider catalog/brand suggestion meaning becomes an accepted consumer requirement |
| `price_suggestion` | DEFER | Market/Economics adopts that provider evidence for a named Product use |
| promotion candidates/offers / `best_price_eligible` | DEFER | promotion/campaign participation becomes accepted Product scope |
| provider image-error/picture feeds | DEFER | a named Offering/Readiness media-correctness consumer proves the need |

A DEFER row is not pre-authorization to subscribe or implement the topic.

## 11.3 REJECT FROM PRODUCT 1.0 BASELINE

| Provider signal/topic | Disposition |
|---|---|
| legacy `orders` when `orders_v2` is selected | REJECT duplicate/legacy acquisition source for same purpose |
| `created_orders` / unconfirmed checkout-only feed | REJECT as Sale business meaning baseline |
| order feedback/reputation feed | REJECT — outside current Product charter |
| buyer/seller messages/chat | REJECT — outside Product 1.0 |
| questions/Q&A | REJECT — outside Product 1.0 |
| leads / vertical-specific lead-credit signals | REJECT — outside Product 1.0 |
| unknown/unclassified provider topic | REJECT until explicitly admitted; never become generic event |

---

# 12. Initial Mercado Livre subscription direction

Subject to real provider application/Installation capability revalidation, the smallest current subscription set implied by Product 1.0 is:

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

Provider application configuration never becomes the authority for admission; it is configured to follow this accepted matrix.

A future provider topic may map to an already accepted acquisition family without creating a new D1 meaning, but still requires explicit ingress-admission review before enablement.

---

# 13. Product API / OpenAPI separation

Provider protocol ingress:

- is not under `/organizations/{organization_id}/...` Product roots;
- is not `/access-context`;
- is not included in the Product SDK;
- is not authorized with W4 Product Permissions;
- does not use provider access tokens to impersonate Product Principals;
- does not use Product bearer tokens as provider callback proof;
- may use provider-specific technical route/transport vocabulary because this is deliberately the protocol boundary.

The single Product OpenAPI wire authority remains the Semantic Product API authority. Provider protocol conformance and internal typed acquisition contracts are separate technical contracts, not parallel Product business wire authorities.

Exact public technical HTTP path spelling remains later ingress/Wire closure after the OAuth lane is decided; do not create generic `/providers`, `/integrations`, `/webhooks` Product roots.

---

# 14. OAuth / Authorization Ceremony lane — NOT DECIDED HERE

OAuth begin/callback/re-authorization belongs to the same technical protocol boundary but is **not** an external acquisition signal.

It has different semantics:

```text
Product-authorized human initiates technical authorization ceremony
        ↓
provider authorization protocol
        ↓
callback / code / state validation
        ↓
credential + provider seller/account namespace binding
```

Do not force OAuth callbacks into the acquisition inbox or generic `ExternalIngressEvent` model.

**Exact next sub-batch:** derive the OAuth/Authorization Ceremony lane, including initiation authority, state/correlation, current-authority revalidation at callback time, provider seller identity proof, rebinding/mismatch rules, replay/expiry and post-callback acknowledgement/redirect semantics.

---

# 15. Negative controls

Later proof must make at least these defects invalid/unreachable:

1. unknown provider topic becomes a generic event;
2. provider DTO reaches a D1 owner as its public semantic contract;
3. provider user/application/account identity becomes Organization directly;
4. ambiguous Installation correlation selects the first match;
5. arbitrary callback resource URL causes arbitrary outbound HTTP;
6. positive callback acknowledgement is recorded as owner/business success;
7. positive acknowledgement occurs before recoverable responsibility is established;
8. process crash immediately after acknowledgement permanently loses required acquisition;
9. duplicate `orders_v2` delivery creates duplicate Sale meaning;
10. out-of-order provider payload directly regresses owner state without authoritative reread;
11. provider topic name becomes D3 domain-event name by normalization;
12. polling/recovery uses a distinct business processing path from push;
13. missed-feed/recovery window becomes a completeness claim;
14. raw provider payload/PII is retained indefinitely by convenience;
15. unbound quarantine enters Organization-scoped Sales/Offering/Work/Fulfillment state;
16. provider credential authenticates Product API;
17. Product bearer token authenticates provider protocol ingress;
18. deferred stock/full/flex/invoice/promotion topic is activated by provider availability alone;
19. another marketplace must emulate an HTTP webhook to enter the same MPC acquisition path;
20. provider delivery ID becomes Sale/Listing/Shipment/Payment/Claim identity.

---

# 16. Proof strategy

Architecture/D7/D8 proof must be able to falsify at least:

1. duplicate Sale notification → one Sale meaning;
2. positive acknowledgement followed by process crash → acquisition remains recoverable;
3. webhook + missed-feed/recovery for same source resource → same native acquisition path;
4. zero or multiple namespace bindings → no owner commit;
5. arbitrary provider resource URL → no arbitrary fetch;
6. provider payload state disagrees with authoritative reread → reread/owner meaning wins;
7. Payment refund/reversal after earlier release → occurrence/history remains recoverable where D3/D4 requires it;
8. claim + claim-action notifications → same Claim acquisition without generic claims workflow;
9. competition signal → Market evidence may change, no PriceIntent is created by ingress;
10. a second provider using queue/event-bus transport can produce the same accepted native acquisition family without introducing webhook assumptions.

Exact executable mechanism belongs to D7/D8.

---

# 17. Reopen triggers

Reopen only the smallest implicated ingress/D4/owner decision when material evidence shows:

- a provider signal itself contains a material occurrence that cannot be recovered from the accepted authoritative read and must therefore be preserved as external occurrence evidence;
- exact provider namespace/Organization correlation cannot be established with current D2/D4 bindings;
- a real Product consumer requires a deferred/rejected provider signal;
- a second provider proves the typed acquisition family cannot remain provider-independent without semantic loss;
- recoverable custody cannot satisfy a provider acknowledgement contract without changing D3/D4 recovery assumptions;
- Product 1.0 requires a business-system inbound callback not currently admitted by D4;
- a technical ingress state becomes user/business manipulable in a way that would require a real Product operation/owner.

Do not reopen for queue preference, current handler layout, provider topic availability or desire for a generic integration framework.

---

# 18. Ingress-A outcome

**Outcome:** `RESTRUCTURE NOW — D5-LOCAL TECHNICAL INGRESS` + `DEFER SAFELY` unsupported provider modes/topics.

> **Use provider-specific inbound adapters feeding one MPC-native recoverable acquisition seam and closed typed acquisition requests. Push, polling and recovery converge on the same authoritative D4 acquisition path. External delivery never becomes Product authentication, Product operation, domain event by receipt or generic Integration authority.**

Ingress-A is accepted in-stage by explicit operator ratification on 2026-08-19.

Technical non-Product ingress remains open because the **OAuth / Authorization Ceremony lane is NEXT**. Whole-Ingress coherence and final canonical acceptance occur only after that lane is resolved and the combined package is adversarially reviewed.

Implementation remains blocked until D9.
