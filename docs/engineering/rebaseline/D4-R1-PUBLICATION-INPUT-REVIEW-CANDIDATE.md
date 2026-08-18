# D4-R1 — Publication Input & Listing Authoring Contract — CONVERGED REVIEW CANDIDATE

> **Status:** NON-AUTHORITATIVE — Fable reviewed / GPT adjudicated / awaiting explicit operator ratification  
> **Stage:** targeted D4 reopen while D5-B2 remains NEXT / NOT YET OPENED  
> **Parent authority:** accepted D0 → D4 + D5-B1  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Review evidence:** `AI-DIALOG.md`, Fable D4-R1 + Pre-D5-B2 Global Coherence Independent Adversarial Review, 2026-08-18  
> **Purpose:** close the publication/input seam discovered before D5-B2 without creating a second product architecture, duplicate authority, hidden PIM, generic integration framework or client-visible workflow choreography. This file does not change D4/D5 authority, router status, ADR status or implementation permission.

---

# 1. Why D5-B2 stops here

D0 already requires Product 1.0 to take an eligible source Product through readiness and **creation/publication → real channel observation**.

Accepted authority already establishes:

- Product master remains external and source-qualified;
- Product & Channel Readiness owns Product↔channel correspondence, channel requirements and supported/missing/conflicting/readiness meaning;
- Marketplace Offering Operations owns offer/listing representation and lifecycle, Listing Intent, Price Intent and convergence;
- Availability Control owns Sellable Availability and Availability Intent;
- provider/business-system adapters own protocol, never business authority;
- durable Listing Intent already has an owner-local MPC identity;
- external Listing remains source-qualified provider identity;
- D5-B1 requires a semantic Product API and forbids provider DTO/resource ontology as the client contract.

What was not canonically decided is how external Product facts, provider publication requirements and MPC-authored channel values combine into one explainable Listing Intent without creating another Product master, a second desired-listing authority or a generic source/integration platform.

D5-B2 cannot safely freeze Listing requirements, Listing Intent authoring operations or input-integration surfaces until this seam is accepted.

---

# 2. Global Coherence result

Whole-product review and independent Fable challenge converge on the following result:

> **D0→D5-B1 remains one coherent architectural system. The missing publication seam closes inside existing owners and identities. No 13th business domain, Product/PIM master, PublicationPreparation aggregate, SourceProductObservation owner, rule engine, workflow engine or generic connector platform is required.**

Two earlier parasitic directions are explicitly rejected:

```text
REJECTED

Readiness-owned mutable PublicationPreparation
        +
Offering-owned desired Listing representation / ListingIntent
```

and:

```text
REJECTED

generic SourceProductObservation business owner/service
```

The accepted candidate direction is:

```text
external Product/source evidence
        ↓ D4 acquisition / consumer-owned ports
Readiness
  requirements
  correspondence
  source candidates
  source-level readiness
        ↓ Q
Offering
  ListingIntent draft
  desired listing representation
  draft dispatchability
        ↓ freeze / authorization / execution safety
D4 + D7 execution mechanism
  provider protocol / joint realization where required
        ↓
marketplace
        ↓ authoritative reread
owner-specific convergence / divergence
```

The same model handles creation and editing.

---

# 3. Decision question

> **What is the smallest sustainable contract that lets MPC create and edit marketplace listings when Product facts may come from Sankhya or another external source, missing channel values may temporarily be authored by humans/automation in MPC, and provider-specific requirements must be satisfied without creating an MPC Product Master, generic PIM, provider-field bag, second listing-draft authority, generic source-observation authority or generic integration platform?**

The solution must support:

1. current incomplete source Product data;
2. operator/automation completion through MPC;
3. future Sankhya-dominant source enrichment, including presentation/media where available;
4. built-in source adapters and later external connectors;
5. provider requirement/schema churn;
6. creation and editing through one consistent model;
7. Mercado Livre first without making its resource topology MPC ontology;
8. future second marketplace/business system without redesigning Product identity or owner boundaries;
9. providers that physically combine listing, price, quantity or fulfillment inputs in one request without merging MPC authorities.

---

# 4. Evidence classification

## 4.1 KNOWN from repository authority

1. Product 1.0 requires real creation/publication, not read-only observation.
2. Product master remains externally authoritative; MPC has no Product/PIM/MDM domain.
3. Product identity is `SourceInstance + native Product key`.
4. Provider Listing/Variation remains external source-qualified identity.
5. Readiness owns correspondence, channel requirements and readiness conclusion.
6. Offering owns listing representation/lifecycle, Listing Intent, Price Intent and convergence.
7. Availability owns Sellable Availability and Availability Intent.
8. Economics owns economic interpretation; it never writes marketplace price.
9. D4 adapters own provider protocol/DTO/resource selection and provider capability/requirement evidence.
10. D4 final coherence forbids one consumer or D4 from owning a provider payload/resource wholesale.
11. Material Listing Intent may have durable owner-local identity and must remain historically explainable.
12. D2 allows immutable decision-time snapshots but they never become current producer authority.
13. D2 forbids automation from silently reopening/reversing a standing human decision in the same semantic scope.
14. D3 requires no cross-owner atomic mutation; cross-owner workflows are correlated/convergent.
15. D3 allows shared external-effect safety mechanics to verify required owner-issued proofs without owning the answers.
16. D5-B1 fixes semantic Product API, source-qualified external identity, honest knowledge/freshness, fail-closed consequential idempotency and one machine-readable Product API wire authority.
17. D5-B2 remains non-authoritative and frozen.

## 4.2 Current Mercado Livre evidence — non-authoritative external evidence

Current official Mercado Livre documentation reviewed on 2026-08-18 establishes materially relevant facts for the selected User Product direction:

- current User Product publication changes the legacy Item/Variation shape;
- `family_name` is seller-supplied while `title` is generated in the current documented UP flow;
- current documented new-item examples carry `available_quantity` with the same creation request;
- without multi-origin stock, stock update uses Item `available_quantity` and Mercado Livre synchronizes associated Items of the same User Product;
- UP-level shared characteristics may replicate asynchronously;
- family-level edits can use asynchronous task semantics with per-member status;
- provider documentation also contains specific zero-quantity/paused and existing-UP condition-of-sale flows, but those do not prove a representation-first PASS-A path for the selected first active from-scratch lane.

Primary references used by Fable/current review:

- https://developers.mercadolibre.com.ar/es_ar/api-docs-es/precio-variacion
- https://developers.mercadolibre.com.ar/es_ar/preguntas-y-respuestas/stock-distribuido
- https://developers.mercadolibre.com.ar/es_ar/atributos-y-variaciones/user-products
- https://developers.mercadolibre.com.ar/es_ar/es_ar/publica-productos

The accepted D4-B2 Installation evidence remains the local context: current seller is User Product model, measured 34/34 Items in that model, no multi-origin/Full observed, selected availability candidate is non-multi-origin Item-path `available_quantity`, and shared-UP blast radius must remain revalidated.

## 4.3 External structural benchmarks — non-authoritative

- Amazon Product Type Definitions demonstrates marketplace/product-type-specific requirement schemas and schema churn; product-only and offer-only modes show provider topology need not equal MPC ownership.
- Mirakl distinguishes Product, Offer and Price/Stock flows; provider workflow decomposition therefore must remain adapter-local.
- Google Merchant Center demonstrates that primary/supplemental source precedence must be explicit rather than accidental last-write-wins; MPC does not copy its general rules engine.
- Akeneo/marketplace-hub systems legitimately centralize Product/catalog attributes because Product/PIM is part of their product authority; that is anti-model evidence for MPC, whose D1/D2 authority rejects a Product master.

## 4.4 INFERRED

1. A fixed universal `MPCListingDTO` is structurally unstable.
2. A generic provider `attributes`/`extensions`/`map[string]any` bag only renames DTO leakage.
3. Source truth and MPC-authored listing values require different provenance/authority.
4. ListingIntent is the smallest existing identity that can carry mutable authoring before freeze without creating a second publication aggregate.
5. Readiness must not read ListingIntent drafts; otherwise the accepted `Readiness → Offering` edge silently becomes bidirectional.
6. Offering can evaluate draft dispatchability by consuming current Readiness meaning through the existing Q edge.
7. Providers may physically require values from multiple owners in one external call; the correct solution is joint technical realization of owner-issued meanings, not transfer of one owner's meaning into another owner.
8. Embedded pull adapters and future external connectors may reuse translation/admission mechanics without sharing a mandatory HTTP hop or a generic source service.
9. FOLLOW_SOURCE + EXPLICIT_OVERRIDE is sufficient baseline; a generic DERIVED/rule/mapping engine lacks a present consumer.

## 4.5 UNKNOWN / DEFERRED

- exact final OpenAPI paths/request/response schemas — D5;
- exact reusable source→requirement mapping primitive, if repeated real cases later prove one necessary;
- exact media/blob/CDN/storage mechanics — D7;
- exact connector ingress wire contract — only when a real external connector consumer appears;
- exact provider task/poll/retry mechanism — D7;
- exact first controlled Mercado Livre create/edit effect — D8;
- exact D6 editor/screens/workflow UX.

Unknown stays Unknown.

---

# 5. Root cause

D4 previously began publication too late.

D1/D2 already defined Product identity, Readiness and Offering/ListingIntent ownership. D4 already knew provider requirements/protocol/effect safety. The missing contract was the seam that turns external Product facts + current channel requirements + explicit MPC-authored values into one owner-valid ListingIntent and then into provider protocol without moving truth selection into the adapter.

Without an explicit seam, the system could drift into:

- adapter-side truth selection/fallback;
- manual values masquerading as source facts;
- source changes silently overwriting operator decisions;
- Readiness and Offering each owning adjacent desired-publication state;
- provider requirement IDs/payloads becoming a hidden Product model;
- a generic source-observation service acquiring Product authority;
- agents/scripts bypassing owner APIs;
- initial provider create merging Offering and Availability authority because both fields happen to share a request;
- client/API design freezing single-shot success where provider creation/edit is partially asynchronous.

---

# 6. Target invariant

> **Source Product truth, Readiness requirement/correspondence meaning, Offering-owned ListingIntent authoring, Availability/Price/Fulfillment owner meanings, provider requirement evidence and external Listing state remain distinct authorities. Every material provider publication input is owner-issued/resolved before protocol serialization; technical execution may jointly realize multiple owner-issued meanings in one provider request without any owner storing or acquiring another owner's meaning.**

Corollaries:

1. Source observation never becomes MPC Product authority merely because it is persisted or reused.
2. MPC-authored override never masquerades as Sankhya/source observation.
3. Provider requirement metadata never becomes universal Product ontology.
4. Missing source value remains missing unless an explicit MPC-owned ListingIntent value legitimately supplies the listing requirement.
5. No implicit last-write-wins between source and authored values.
6. Readiness never owns desired Listing representation.
7. Offering never owns Sellable Availability merely because provider creation carries quantity.
8. Availability never owns listing representation.
9. ListingIntent is the single create/edit authoring identity and only carries Offering-owned desired listing meaning plus typed references/correlations to other owner meanings where needed.
10. Adapter/execution mechanism may compose already-issued owner effect inputs; it never computes, defaults or reassigns their authority.
11. Each owner concludes its own convergence from authoritative external evidence.
12. Historical ListingIntent snapshots are decision/action context, never a second current authority.
13. Automation does not gain override authority by recurrence.
14. Source changes are re-observed/revalidated at material decision points; no baseline acquisition→Offering push edge is invented.
15. Provider multi-step/async behavior remains explicit; accepted submission never implies complete convergence.

---

# 7. Selected Global Maximum

## Rejected A — Sankhya-only authoring

Wait until every source field is complete.

Rejected: makes external data-completion a hidden Product 1.0 gate and provides no legitimate temporary operator/agent path.

## Rejected B — MPC Product/PIM master

Import/own complete Product content, then map it to channels.

Rejected: duplicates external Product authority and moves the product boundary toward PIM/hub semantics not accepted by D0/D1.

## Rejected C — generic provider field bag

Store arbitrary provider fields in Product/ListingIntent JSON.

Rejected: provider protocol becomes domain state by renaming and loses authoritative typing/provenance.

## Rejected D — Readiness PublicationPreparation + Offering ListingIntent

Rejected by Global Coherence: adjacent desired-publication authorities and separate create/edit architectures become reachable.

## Rejected E — generic SourceProductObservation business service/domain

Rejected: source acquisition mechanism would acquire authority over a payload that can legitimately feed multiple consumer-owned semantics.

## Rejected F — mandatory self-HTTP for embedded adapters

Rejected: transport consistency buys accidental network/auth/deployment coupling without improving semantic authority.

## Selected G — Readiness requirements + Offering ListingIntent + owner-preserving joint provider realization

- Readiness owns requirements, correspondence, source candidates and source-level readiness;
- Offering owns ListingIntent draft, desired representation and draft dispatchability;
- source/manual resolution inside ListingIntent is FOLLOW_SOURCE or EXPLICIT_OVERRIDE;
- Availability/Price/Fulfillment retain their own meanings;
- D4/D7 may jointly serialize already-issued owner inputs when the provider requires one physical request;
- external provider topology remains adapter-local;
- source acquisition remains D4 evidence/mechanism feeding consumer-owned ports;
- no connector platform or extra business domain.

**Global Maximum / recommended.**

---

# 8. Readiness contract — source/channel readiness, not intent dispatchability

Readiness remains authority for:

- Product↔channel correspondence;
- applicable channel/publication requirements;
- provider/category/product-type requirement meaning translated through its consumer-owned D4 port;
- source candidates/evidence relevant to requirements;
- missing/conflicting/unsupported source-level state;
- source-level Product↔channel readiness/sufficiency.

Conceptual subject:

```text
Organization
+ Marketplace Installation
+ SourceInstance-qualified Product
```

Readiness may query D4 through a **Readiness-owned publication-requirements port**.

The D4 adapter translates current provider requirements into only the semantics Readiness needs, proportionately including:

- requirement identity inside provider context;
- required/recommended/optional/conditional applicability;
- data kind/cardinality/options/constraints where materially needed;
- editability/immutability when material;
- provider schema/version/checksum/revision evidence when exposed;
- acquisition/provenance/time.

This is not a generic ProviderRequirement business framework.

Readiness MUST NOT:

- own desired listing copy/text/media selection;
- own a PublicationPreparation aggregate;
- read ListingIntent content in order to produce readiness;
- own Listing lifecycle;
- own Price Intent;
- own Sellable Availability;
- own provider actual state;
- own arbitrary provider payloads;
- own Product master lifecycle.

### Source-level readiness

Readiness answers what the current Product/channel/source evidence can establish and what requirements are missing/conflicting at that level.

An EXPLICIT_OVERRIDE inside a ListingIntent does **not** rewrite this source-level conclusion. A source requirement may remain source-missing while a particular ListingIntent is still dispatchable because Offering has an explicit channel value.

This distinction removes the residual Readiness↔Offering ownership cycle found by Fable.

---

# 9. Offering contract — ListingIntent is the single create/edit authoring aggregate

D2 already admits durable ListingIntent identity. D5-B1 allows owner-specific intent operations.

A ListingIntent may begin as an editable **draft** under Offering for one desired listing action.

Conceptual correlation includes:

```text
Organization
+ Marketplace Installation
+ SourceInstance-qualified Product
+ existing source-qualified Listing where editing/closing an existing listing
```

Exact API path/addressing remains D5.

A ListingIntent draft may carry only Offering-owned desired Listing meaning, for example:

- intended create/edit/close action;
- target category/product-type context where this expresses desired Listing representation;
- channel-facing content values;
- media selection/order/role for this Listing;
- provider-qualified requirement resolutions used for this desired representation;
- FOLLOW_SOURCE references;
- EXPLICIT_OVERRIDE values;
- typed references/correlations to other owner intents/inputs required for one provider effect, without copying their owned values into Offering state.

It MUST NOT become Product master, Economics authority, Sellable Availability authority, Fulfillment authority or provider DTO mirror.

### Creation and editing are one architecture

```text
Create listing
  → ListingIntent(target=no existing Listing)

Edit listing
  → ListingIntent(target=existing source-qualified Listing)
```

Provider differences determine protocol realization, not a second MPC authoring model.

### Draft dispatchability belongs to Offering

Offering owns the validity question:

> **Can this particular ListingIntent be frozen/dispatched now under the current Readiness-owned requirements/correspondence and current owner-specific prerequisites?**

Offering answers this by Q against Readiness and other already legal producer meanings. It does not mutate or recompute Readiness authority.

A requirement revision may therefore make a draft non-dispatchable while leaving Readiness' source-level conclusion independently unchanged.

---

# 10. Publication value resolution — FOLLOW_SOURCE | EXPLICIT_OVERRIDE only

No generic rule/mapping DSL is admitted.

## FOLLOW_SOURCE

The ListingIntent references one explicit source-qualified fact/candidate admitted through current Readiness/source evidence.

Rules:

- source identity/provenance is explicit;
- value re-resolves from current source meaning at material decision points before freeze/dispatch;
- source unavailable/unknown stays unavailable/unknown;
- no silent fallback to a manual value;
- a prior draft read never proves current source validity at dispatch.

A source change does **not** imply a baseline push event into Offering. Current value is revalidated through existing Q/current-source semantics at material decision points such as draft read where needed, freeze and dispatch-time safety checks.

## EXPLICIT_OVERRIDE

The ListingIntent contains an MPC-authored channel value.

Rules:

- effective Principal/automation and time are preserved;
- override never mutates or falsifies source Product truth;
- source change does not silently replace the override;
- source/provider requirement drift may make the draft stale/non-dispatchable if material;
- removing/replacing the override is an explicit Offering edit;
- future source takeover is deliberate: a new edit can clear the override and use FOLLOW_SOURCE without rewriting historical intents.

## No baseline DERIVED

A generic derived/transformation mode is rejected for Product 1.0 baseline.

Automation Principals can author explicit values with attributable history. If repeated real cases later prove a deterministic owner-specific transformation materially reduces total complexity, admit the smallest explicit transformation then; do not prebuild feed rules, expression DSL, mapping engine or formula language.

---

# 11. Human / automation authority protection

Claude/Codex/Fable-style automation is a D2 automation Principal and uses the same Semantic Product API as other Product API clients.

It does not use Integration Ingress and cannot claim a value came from Sankhya.

D2 §10.3 remains binding:

> automation recurrence never grants authority to silently reopen or reverse a standing human decision in the same semantic scope.

Therefore:

- a recurring automatic run cannot silently replace a human-authored EXPLICIT_OVERRIDE;
- it cannot silently flip a human-selected resolution mode;
- replacement/supersession must be explicit Offering/domain semantics with preserved historical attribution;
- where current business disposition or Governance requires additional authority, ordinary Product API access does not bypass it;
- a one-off explicit human-directed agent operation is still attributed to the automation Principal and must follow the same explicit supersession semantics rather than mutating history invisibly.

No AI-specific authority/backdoor exists.

---

# 12. Source acquisition / future integration ingress

Source acquisition remains D4 mechanism/evidence, not a Product domain or generic source owner.

One source acquisition may legitimately feed multiple consumer-owned semantic ports. No component owns the source Product payload as a whole.

## Embedded source adapter

Example:

```text
Sankhya sanctioned API
    ↓
D4 Sankhya adapter
    ↓
consumer-owned semantic ports
```

No mandatory self-HTTP loopback.

## Future external connector

A separately deployed connector may enter through a bounded Integration Ingress HTTP adapter **when a concrete external connector exists**.

```text
external connector
    ↓ authenticated bounded ingress
D4 HTTP inbound adapter
    ↓ translation/admission mechanics
consumer-owned semantic ports
```

This candidate prepares that seam but does **not** freeze or require a generic Source Ingestion HTTP API now.

When the first real external connector is admitted, D5 defines the smallest necessary contract while preserving Organization + SourceInstance binding, source provenance, partial/absent semantics, replay safety and inability to impersonate MPC-authored values.

No generic `/entities`, `/resources`, plugin registry or connector platform is introduced.

---

# 13. Media / images

Media participates in Listing authoring without creating Product-media authority.

Possible origins:

- source-qualified Product media from Sankhya/another external source;
- MPC-authored/uploaded media for this Listing context.

Rules:

- source media remains source evidence;
- MPC upload never rewrites the source Product;
- ListingIntent owns only listing-specific selection/order/role;
- provider image IDs/CDN/resource topology remains provider-local;
- ListingIntent historical snapshot preserves enough media/provenance to explain what was attempted;
- blob storage, hashing, resizing, cache/CDN and upload mechanics remain D7/implementation;
- arbitrary external URL strings are not automatically trusted media.

No reusable ProductAsset/AssetFamily/media library is admitted. If a reusable cross-listing media master becomes a real requirement, that is explicit Product/PIM pressure and must be re-adjudicated rather than grown silently.

---

# 14. Price / Availability / Fulfillment authority fence

A provider may physically combine content, price, quantity and fulfillment in one request. This does not merge MPC owners.

## Price

Offering already owns Price Intent. Listing creation may correlate to a current valid Price Intent where the provider requires initial price. Economics informs; adapter does not calculate/choose price.

Price remains Offering-owned and may therefore be correlated naturally with ListingIntent without a cross-owner transfer.

## Availability

Availability owns Sellable Availability and Availability Intent.

### R1-G1 — CLOSED / PASS-B

Current official Mercado Livre User Product evidence shows the selected publication direction may physically require `available_quantity` in the same new-item request that creates the Listing/UP representation. PASS-A — normal active representation-first creation without availability input — is **not claimed** for the selected current lane.

This does **not** require `Availability → Offering` and does not reopen D1/D3.

The correct composition is:

```text
Offering-owned ListingIntent
  desired representation / PriceIntent correlation
        │
        │ owner-issued effect input/proof
        ▼

Availability
  consumes intended marketplace target through existing Offering → Availability meaning
  computes/owns Sellable Availability
  issues Availability-owned effect input / Availability Intent
        │
        │ owner-issued effect input/proof
        ▼

D4/D7 external-effect execution mechanism
  validates/correlates both owner-issued inputs
  serializes one provider request when provider protocol requires one
        │
        ▼
Mercado Livre
        │ authoritative reread
        ├─ Offering concludes representation/Listing convergence
        └─ Availability concludes availability convergence
```

Binding laws:

1. Offering never authors, stores or owns the Availability quantity as ListingIntent content.
2. ListingIntent may preserve a typed correlation/reference to the jointly realized Availability-owned input/intent where needed for historical explanation.
3. Availability never authors desired Listing representation.
4. The technical execution mechanism may require both owner-issued inputs/proofs before one provider call; mechanism does not own either answer.
5. Adapter cannot read quantity from an Offering field, default it, calculate it or infer it from provider data.
6. Active creation without the required Availability-issued value fails closed before provider dispatch.
7. No cross-owner atomicity is claimed; one provider request may yield partial/asynchronous semantic convergence and each owner evaluates its own result.
8. The paused/zero-quantity representation-first path remains unclaimed for this Installation until a real proof establishes it.
9. The first controlled real creation write remains D8 proof, including authoritative reread and shared-UP blast-radius verification.

### Why no D1/D3 reopen

D1 already establishes `Offering → Availability`: Availability may consume the marketplace representation/target needed to synchronize availability. D1 also explicitly says a provider API combining multiple fields/actions does not merge business authorities.

D3 already says:

- cross-owner workflows are correlated/convergent rather than atomically co-owned;
- no cross-owner atomicity is invented;
- shared external-effect safety mechanics may verify required owner-issued validity/authorization/correlation without owning the answers.

Therefore one provider request can be a **joint technical realization of two owner-owned effects** without either business owner consuming the other's value.

### Reopen trigger

Reopen only if real provider execution proves Availability cannot issue a valid value for the intended not-yet-existing target without a new business dependency, or if correctness requires Offering to persist/own Availability meaning rather than merely correlate to a separate owner input.

## Fulfillment

Provider shipping/fulfillment requirements remain D4 evidence and Fulfillment retains business responsibility.

If a concrete provider create/edit request physically requires a Fulfillment-owned semantic input, the same law applies: a technical execution mechanism may compose owner-issued meanings, but ListingIntent does not absorb Fulfillment policy/state and no new edge is hidden by convenience.

A genuinely new semantic dependency must reopen D1 rather than be buried in the adapter.

---

# 15. Provider execution / multi-step convergence

Offering owns ListingIntent and its listing-convergence conclusion. D4 owns provider protocol/effect translation.

Provider execution is not assumed to be one atomic request or one boolean outcome.

Current Mercado Livre evidence includes:

- item creation and description publication as separate provider operations in relevant flows;
- asynchronous User Product shared-field propagation;
- family-level asynchronous tasks with member statuses;
- shared-field blast radius across Items of the same User Product;
- provider-specific caps/constraints such as per-UP sales-condition limits and family-edit restrictions.

Therefore the D4 external-effect contract preserves, where material:

- accepted/rejected/pending/ambiguous;
- step/member/aspect-level confirmed/rejected/pending/ambiguous/not-executed outcomes;
- authoritative reread per material resource/aspect;
- no whole-operation convergence merely because the first provider request returned 2xx/201/202;
- no blind retry of already-confirmed or ambiguous steps;
- shared-UP intended/authorized/attempted blast radius.

B2 must not freeze a single-shot `createListing = success` contract.

Adapter MUST NOT:

- choose between competing source/override values;
- invent missing required values;
- query hidden systems to fabricate semantic validity;
- copy provider DTOs into a generic Listing model;
- turn transport 2xx/201/202 into convergence;
- widen intended/authorized scope silently;
- recalculate Price/Sellable Availability/Fulfillment meaning.

---

# 16. Whole-product coherence

## One architectural grammar

The publication seam uses the same grammar already accepted elsewhere:

```text
external evidence
  → semantic owner
  → owner intent/decision
  → controlled external effect
  → authoritative reread
  → owner convergence/divergence
```

It does not introduce a second product architecture.

## Identity

- Product remains external source-qualified identity.
- Listing remains external source-qualified identity.
- ListingIntent is the existing Offering-owned action identity.
- Availability Intent remains Availability-owned.
- no ProductDraft/PublicationDraft/ProviderResource/SourceObservation identity is added.

## Communication

- Readiness → Offering remains Q.
- Offering → Availability remains the existing target/representation dependency.
- source acquisition remains D4 evidence/mechanism.
- no Availability → Offering edge is created by joint provider dispatch.
- no generic workflow bus/command taxonomy.

## Work / Governance

- missing/conflicting requirements may become Work only when materially actionable;
- assignment never becomes authorization;
- Governance authorizes consequential ListingIntent where required;
- neither Work nor Governance owns Listing content, source truth or Sellable Availability.

## Materialization / Fulfillment consistency

The same pattern already used by Materialization/Fulfillment survives: each owner owns its own meaning, cross-owner progression is explicit/correlated, and no shared mutable workflow entity is introduced.

## Implementation posture

The architecture remains sustainably implementable as a **modular monolith first**:

```text
one Go backend
  contexts/      semantic owners
  adapters/      provider/source translation
  platform/      execution safety/runtime mechanics
  composition/   assembly only
  views/         read projections

one PostgreSQL target
```

No microservice split, broker, generic event platform, workflow engine, schema registry, PIM engine or connector framework is required by this decision.

---

# 17. Complexity / YAGNI

## Essential complexity retained

- source-qualified Product facts;
- provider/category/product-type-specific requirements;
- requirement/schema revision/provenance;
- source-level readiness versus draft dispatchability;
- ListingIntent draft/history;
- FOLLOW_SOURCE versus EXPLICIT_OVERRIDE;
- human/automation supersession safety;
- media origin/selection;
- provider-specific serialization;
- multi-step/partial convergence;
- Price/Availability/Fulfillment authority separation;
- joint realization mechanism where provider physically combines fields;
- agent attribution;
- source acquisition provenance;
- shared-UP blast radius.

## Accidental complexity explicitly rejected

- MPC Product Master/PIM;
- PublicationPreparation aggregate;
- universal ProductAttribute master;
- generic ProviderRequirement business framework;
- generic SourceProductObservation business service;
- arbitrary provider JSON bag;
- generic mapping/rules DSL;
- baseline DERIVED mode;
- connector/plugin framework;
- mandatory self-HTTP;
- mandatory Source Ingestion API before a real connector;
- AI-specific API/backdoor;
- creation architecture separate from edit architecture;
- one giant Publication aggregate owning price/stock/fulfillment;
- microservices/event platform/workflow engine by architectural fashion.

---

# 18. Proof strategy before implementation

## P1 — incomplete source → authored intent → future source takeover

1. Sankhya Product lacks one required source value.
2. Readiness reports the source-level requirement missing and source candidates honestly.
3. automation/human Principal creates/edits ListingIntent with EXPLICIT_OVERRIDE.
4. Readiness source-level state does not falsely become source-complete.
5. Offering determines that this particular draft is dispatchable under current requirements.
6. frozen ListingIntent preserves override provenance.
7. later Sankhya supplies the value.
8. historical ListingIntent is unchanged.
9. new edit intent may explicitly switch to FOLLOW_SOURCE.

## P2 — source spoofing

- Product API client cannot claim an override was a Sankhya observation.
- source acquisition path cannot author an Offering override by choosing a payload flag.

## P3 — partial source observation

Omitted field in partial acquisition does not clear a prior known source fact unless source semantics establish removal/absence.

## P4 — provider requirement/schema change

Historical ListingIntent remains explainable under requirement revision N. Current draft re-evaluated against N+1 becomes non-dispatchable when a new required value is unsatisfied.

## P5 — second-provider structural inversion

Mercado Livre User Product/Item, Amazon product/offer modes and Mirakl Product/Offer/PriceStock can differ in D4 serialization without changing Product identity, Readiness/Offering owners or core intent semantics.

## P6 — media source transition

One intent follows source media and another uses MPC-authored media without rewriting Product media authority.

## P7 — no adapter truth selection

Remove one required ListingIntent resolution. Adapter must fail before dispatch; it cannot choose fallback/default.

## P8 — R1-G1 joint realization

Negative fixtures:

- active ListingIntent dispatch requiring quantity but lacking Availability-issued input → fail-closed, zero provider calls;
- adapter contract cannot read quantity from ListingIntent content;
- ListingIntent snapshot may reference Availability intent/effect correlation but contains no Offering-owned quantity field;
- Offering convergence cannot mutate Availability convergence and vice versa;
- one partial/async provider outcome remains separately reconcilable.

D8 first controlled create proves the real provider lane, authoritative reread and blast radius.

## P9 — Readiness/Offering direction

- Readiness public contract has no operation accepting ListingIntent content;
- Offering draft with unresolved required value is non-dispatchable;
- EXPLICIT_OVERRIDE can make the draft dispatchable without changing Readiness' source-level missing state;
- requirement revision can flip dispatchability without requiring Offering → Readiness mutation.

## P10 — human/automation decision protection

Recurring automation attempts to replace a human-authored override in the same scope → rejected or explicitly superseded under domain semantics; both authorships remain historical.

## P11 — decision-point revalidation

FOLLOW_SOURCE value changed in source after draft creation but before freeze/dispatch → current Q/revalidation detects it; correctness does not depend on an acquisition→Offering event.

## P12 — modular-monolith realization

The design must be implementable with in-process owner ports and shared technical execution safety without requiring network service boundaries or distributed transactions.

---

# 19. GPT adjudication of Fable findings

| Finding | GPT adjudication | Resolution |
|---|---|---|
| **F-R1-1** joint ML create physically carries Availability but candidate risked implying Availability→Offering | **AGREE WITH TIGHTENING** | R1-G1 = PASS-B. Use joint technical realization of separately owner-issued inputs. ListingIntent never stores/owns quantity. D3 external-effect mechanics can verify/correlate owner proofs without owning answers. No D1/D3 reopen. |
| **F-R1-2** Readiness sufficiency collided with Offering overrides | **AGREE** | Split Readiness source-level readiness from Offering draft dispatchability. Readiness never reads ListingIntent. Existing Readiness→Offering Q remains sufficient. |
| **F-R1-3** automation may silently overwrite human override | **AGREE** | Restate D2 §10.3; automation replacement/supersession must be explicit and historically attributable. |
| **F-R1-4** source change wording implied push edge | **AGREE** | Revalidate at material decision points through current Q/source semantics; no baseline acquisition→Offering E. |
| **F-R1-5** creation/edit may be multi-request/async/partial | **AGREE / PROMOTE TO EXPLICIT CONTRACT LAW** | D4/B2 must preserve aspect/member outcomes and convergence; never freeze single-shot atomic create semantics. |

### R1-G1 adjudication

**PASS-B.** Current official ML evidence is sufficient to establish the architectural composition case: quantity may be physically required in initial UP publication; current owners can legally supply separate meanings to one external-effect mechanism without moving Availability authority into Offering.

The first controlled real effect remains D8 proof. PASS-A remains unclaimed.

### Parent-stage reopen adjudication

- D0 — **NO REOPEN**
- D1 — **NO REOPEN**
- D2 — **NO REOPEN**
- D3 — **NO REOPEN**
- D5-B1 — **NO REOPEN**

No surviving reviewer contradiction requires a second Fable round.

---

# 20. Reopen / stop triggers

Reopen only the implicated parent decision when material evidence proves:

1. ListingIntent cannot legitimately carry editable desired Listing representation before consequential freeze → D1/D2 targeted review.
2. Readiness must evaluate authored ListingIntent content to own the required Product/channel conclusion → D1.
3. A selected provider requires Availability/Fulfillment semantics that cannot be owner-issued/jointly realized without a new business dependency → D1/D3.
4. correctness requires atomic cross-owner mutation rather than correlated owner effects → D1/D3.
5. a new durable identity/lineage class is required beyond external Product/Listing + domain intents → D2.
6. provider requirements/effects cannot be represented through consumer-owned ports without raw DTO leakage → D4 redesign.
7. external source connector cannot preserve Organization/SourceInstance/provenance safely → D2/D4/D5 boundary review.
8. repeated real source→listing transforms prove an explicit reusable mapping primitive materially reduces total complexity → targeted Offering/Readiness/D4 decision, not generic PIM.
9. reusable cross-listing content/media/catalog governance becomes a real Product requirement → D0/D1 Product/PIM boundary review rather than silent growth.
10. a second real Product master changes source-authority assumptions → D1/D2.
11. a real external connector creates a materially different API/security/compatibility obligation → D5.
12. current provider changes remove the joint-realization path or materially change UP shared-field/blast-radius semantics → targeted D4 review.

Framework preference, provider symmetry, current-code convenience, screen preference and hypothetical providers are not reopen evidence.

---

# 21. Candidate outcome

**Proposed outcome:** `CURRENT D0→D5-B1 STRUCTURE CONFIRMED` + `RESTRUCTURE NOW` for the missing D4 publication-input seam.

If the operator explicitly ratifies this converged package, canonical D4 should add only the durable meaning below:

1. **Readiness** owns publication requirements, correspondence, source candidates and source-level Product↔channel readiness.
2. **Offering** owns `ListingIntent` as the single create/edit authoring draft, desired Listing representation, draft dispatchability and Listing convergence.
3. Baseline authoring resolution is **FOLLOW_SOURCE | EXPLICIT_OVERRIDE** only.
4. Human/automation authorship and supersession preserve D2 Principal/decision-history rules; recurring automation cannot silently reverse standing human decisions.
5. Source acquisition remains D4 evidence/mechanism feeding consumer-owned semantic ports; no SourceProductObservation owner/service.
6. Embedded adapters do not require self-HTTP; future external connector ingress remains a prepared seam until a real connector requires a wire contract.
7. Provider publication requirements remain source-qualified evidence translated through consumer-owned ports; no generic ProviderRequirement business framework.
8. Media selection belongs to ListingIntent only as desired Listing representation; no Product media master.
9. ListingIntent consequential history preserves materially used source/override/requirement/media/authority context without becoming current source/provider authority.
10. **R1-G1 = PASS-B:** when provider creation physically combines Listing and Availability, D4/D7 jointly realizes separately owner-issued meanings; Offering never owns quantity and Availability never owns Listing representation.
11. Price/Availability/Fulfillment authorities remain intact even when one provider request contains several fields.
12. Provider create/edit may be multi-request/asynchronous/partial; accepted submission never equals complete convergence.
13. Adapter translates resolved/owner-issued meaning and never selects truth, invents fallback or widens scope silently.
14. No MPC Product Master/PIM, PublicationPreparation aggregate, SourceProductObservation service, generic transformation DSL, connector framework, AI-specific architecture, workflow engine or microservice requirement is introduced.
15. The whole design remains compatible with modular-monolith-first realization.
16. D5-B2 may resume only after this R1 package is explicitly ratified and canonically consolidated.

**This file remains NON-AUTHORITATIVE until explicit operator ratification. Implementation remains blocked until D9.**
