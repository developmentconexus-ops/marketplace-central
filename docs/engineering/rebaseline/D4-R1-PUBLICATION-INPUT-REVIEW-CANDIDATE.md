# D4-R1 — Publication Input & Listing Authoring Contract — REVISED REVIEW CANDIDATE

> **Status:** NON-AUTHORITATIVE TARGETED-REOPEN REVIEW CANDIDATE — revised after operator-approved Global Coherence direction  
> **Stage:** targeted D4 reopen while D5-B2 remains NEXT / NOT YET OPENED  
> **Parent authority:** accepted D0 → D4 + D5-B1  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Prepared against:** `44d5d8b116a2dd11e2529be46406563ea9615380` on `docs/global-methodology-alignment`  
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

What is not yet canonically decided is how external Product facts, provider publication requirements and MPC-authored channel values combine into one explainable Listing Intent without creating another Product master or a second desired-listing authority.

D5-B2 cannot safely freeze Listing requirements, Listing Intent authoring operations or input-integration surfaces until this seam is closed.

---

# 2. Global Coherence review correction

A whole-product review was performed before this candidate proceeds to independent review.

Result:

> **D0→D5-B1 remains globally coherent. The material correction is local to this candidate: do not create `Channel Publication Preparation` as a new durable aggregate owned by Readiness. Reuse the already accepted Offering-owned `ListingIntent` as the single authoring/draft identity for both creation and editing.**

This removes a potential authority overlap:

```text
INVALID DIRECTION

Readiness-owned mutable Publication Preparation
        +
Offering-owned desired Listing representation / ListingIntent
```

and replaces it with:

```text
Readiness
  owns requirements / correspondence / source candidates / sufficiency
        ↓ Q
Offering
  owns ListingIntent draft / desired listing representation / lifecycle
        ↓
D4 adapter
  owns provider serialization only
```

The same model handles creation and later editing; no separate publication-draft architecture is introduced.

A second coherence correction also applies:

> **Do not introduce a generic `SourceProductObservation` business owner/service. External source acquisition is D4 mechanism/evidence. One source acquisition may feed multiple consumer-owned semantic ports; no component owns the source Product payload as a whole.**

---

# 3. Decision question

> **What is the smallest sustainable contract that lets MPC create and edit marketplace listings when Product facts may come from Sankhya or another external source, missing channel values may temporarily be authored by humans/automation in MPC, and provider-specific requirements must be satisfied without creating an MPC Product Master, generic PIM, provider-field bag, second listing-draft authority or generic integration platform?**

The solution must support:

1. current incomplete source Product data;
2. operator/automation completion through MPC;
3. future Sankhya-dominant source enrichment, including presentation/media where available;
4. built-in source adapters and later external connectors;
5. provider requirement/schema churn;
6. creation and editing through one consistent model;
7. Mercado Livre first without making its resource topology MPC ontology;
8. future second marketplace/business system without redesigning Product identity or owner boundaries.

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
8. Economics owns pricing/economic interpretation; it never writes marketplace price.
9. D4 adapters own provider protocol/DTO/resource selection and capability evidence.
10. D4 final coherence already forbids one consumer or D4 from owning a provider payload/resource wholesale.
11. A material Listing Intent may have durable owner-local identity and must remain historically explainable.
12. D5-B1 permits owner-specific operations, source-qualified external identity, honest knowledge/freshness, fail-closed consequential idempotency and hard cutover.
13. D5-B2 remains non-authoritative and frozen.

## 4.2 Current external benchmark evidence — non-authoritative

### Mercado Livre

Current official documentation shows the publication model is transitioning to User Products. In current User Product flows:

- fields that were historically Item fields may be generated/inherited from User Product;
- `title`, attributes, pictures, domain/catalog and availability behavior are not one stable flat payload;
- shared characteristic changes may propagate asynchronously across related Items;
- required/conditional attributes depend on category/context;
- stock semantics differ by current resource/model and cannot be assumed to be one universal Listing field.

Primary references:

- https://developers.mercadolivre.com.br/pt-br/publicacao-de-produtos
- https://developers.mercadolivre.com.br/pt_br/publicacao-de-produtos/user-products
- https://developers.mercadolivre.com.br/pt_br/tutorial-tipos-de-publicacao-y-atualizacao-de-artigos/preco-variacao

### Amazon SP-API

Current official documentation shows:

- Product Type Definitions returns product-type/marketplace-specific JSON Schema and constraints;
- Listings Items supports requirement modes including `LISTING`, `LISTING_PRODUCT_ONLY`, and `LISTING_OFFER_ONLY`;
- product-only and offer-only workflows can differ;
- attributes/enumerations change over time, so requirement/schema identity and historical context are materially useful.

Primary references:

- https://developer-docs.amazon.com/sp-api/lang-en_EN/docs/building-listings-management-workflows-guide
- https://developer-docs.amazon.com/sp-api/lang-en_EN/reference/putlistingsitem
- https://developer-docs.amazon.com/sp-api/changelog/update-january-updates-to-listing-attribute-usage-and-enumeration-values-1

### Mirakl

Current official documentation distinguishes Product, Offer and Price/Stock flows and lets channel catalog configuration expose offer/channel-specific fields. This supports keeping provider topology in the adapter rather than forcing Product, offer, price and stock into one MPC aggregate.

Primary references:

- https://developer.mirakl.com/content/product/connect-channel-platform/developer-guide/catalog-flow
- https://developer.mirakl.com/content/product/connect-channel-platform/developer-guide/catalog-configuration

### Google Merchant Center

Primary and supplemental product data sources plus attribute rules demonstrate a real need for explicit source precedence/enrichment. MPC uses this only as evidence that precedence must be explicit; it does not copy Google's rule engine.

Primary references:

- https://support.google.com/merchants/answer/14990942
- https://support.google.com/merchants/answer/14994083

### PIM/hub anti-model evidence

Akeneo/ANymarket-class systems legitimately centralize Product/catalog attributes because Product/PIM is part of their product authority. MPC has explicitly rejected that authority. Their existence validates the operational need for channel mapping without justifying a duplicate Product master inside MPC.

## 4.3 INFERRED

1. A fixed universal `MPCListingDTO` is structurally unstable.
2. A generic provider `attributes`/`extensions`/`map[string]any` bag only renames DTO leakage.
3. Source truth and MPC-authored listing values require different provenance/authority.
4. ListingIntent is the smallest existing identity that can carry mutable authoring before submission and an immutable decision-time snapshot afterward.
5. Readiness can supply requirements/source candidates without owning desired Listing representation.
6. Source acquisition should reuse adapter/translation mechanics, not create a business owner for the external Product payload.
7. Embedded adapters do not need loopback HTTP; external connectors can use an ingress transport later while reaching the same D4/consumer boundaries.
8. Baseline source resolution needs only explicit source following and explicit override; a transformation/rules engine is not yet required.

## 4.4 UNKNOWN

- exact first selected Mercado Livre create/edit User Product flow;
- whether the selected first ML publication can establish representation without simultaneously consuming Availability-owned quantity;
- exact request/response fields of ListingIntent authoring operations — D5;
- exact media storage/processing — D7;
- exact source connector HTTP contract when a real external connector is admitted — D5/D7;
- whether repeated mappings later justify a reusable transformation primitive;
- exact provider-specific requirement fields retained by each Readiness-owned contract;
- exact ListingIntent lifecycle/status names and persistence schema — later D5/D7.

Unknown remains Unknown.

---

# 5. Root cause

D4 currently begins listing publication too late.

It defines provider protocol/effect semantics, while D1/D2 define Product identity, Readiness and ListingIntent ownership, but the authority path does not yet say how source facts and authored channel values become a provider-valid ListingIntent.

Without a contract, these defects remain reachable:

- adapter selects truth by convenience;
- UI/script chooses fallback precedence implicitly;
- manual data masquerades as Sankhya truth;
- source change silently overwrites deliberate channel authoring;
- source outage/unknown becomes empty/default listing data;
- provider requirement fields leak into a universal core DTO;
- Readiness grows into a hidden desired-listing owner;
- MPC grows a hidden PIM;
- external connector and embedded adapter implement different business rules;
- ListingIntent cannot explain the source/override/requirements used by the attempted publication;
- creation and editing acquire different architectures.

---

# 6. Target invariant

> **External Product truth, Readiness requirement/sufficiency meaning, Offering-owned ListingIntent authoring, provider requirement evidence and external Listing state remain distinct. Every material publication value is selected under an accepted MPC owner before provider serialization; the adapter translates but never chooses business truth.**

Corollaries:

1. Source observation stays external evidence.
2. Readiness owns requirement applicability/sufficiency, not desired Listing representation.
3. ListingIntent is the single MPC-owned listing-authoring/draft identity for create/edit actions.
4. An authored override never masquerades as source truth.
5. No implicit last-write-wins across source and authored values.
6. Missing/unknown source remains missing/unknown unless a legitimate explicit override satisfies the listing requirement.
7. Provider requirement metadata is source-qualified evidence, not Product ontology.
8. ListingIntent freezes enough resolved values + requirement/provenance context at consequential submission to explain the action later.
9. Price, Availability and Fulfillment retain their accepted authorities even if the provider combines their values in one wire request.
10. Source acquisition mechanism never becomes a Product/domain owner.
11. Creation and editing use the same Offering-owned intent model.
12. Agent automation uses the same Product API authority path as other MPC clients; no AI backdoor exists.

---

# 7. Global architecture shape

```text
EXTERNAL PRODUCT SOURCES
Sankhya / later source
        ↓
D4 source adapter / acquisition
        ↓
consumer-owned semantic facts / queries
        │
        ├───────────────┐
        ▼               ▼
READINESS            OTHER CONSUMERS
requirements         Availability/Economics/etc.
correspondence       only for their own meanings
source candidates
sufficiency
        │
        │ Q
        ▼
OFFERING — ListingIntent
  draft authoring for create/edit
  ├─ FOLLOW_SOURCE references
  ├─ EXPLICIT_OVERRIDE values
  ├─ selected media for listing representation
  ├─ desired listing change
  └─ requirement/provenance context
        │
        │ freeze / consequential intake
        ▼
Governance if required
        │
        ▼
D4 MARKETPLACE ADAPTER
        │ provider serialization/protocol
        ▼
MARKETPLACE
        │ authoritative reread
        ▼
Offering convergence/divergence
```

No additional Product, PublicationPreparation, ProviderResource or generic workflow authority is introduced.

---

# 8. Readiness contract — requirements and sufficiency only

Readiness remains authority for:

- Product↔channel correspondence;
- applicable channel/publication requirements;
- source candidates/evidence relevant to those requirements;
- missing/conflicting/unsupported state;
- readiness/sufficiency conclusion.

Conceptual subject remains:

```text
Organization
+ Marketplace Installation
+ SourceInstance-qualified Product
```

Readiness may query D4 through a **Readiness-owned publication-requirements port**.

D4 adapter translates current provider/category/product-type requirements into only the semantics Readiness needs, proportionately including:

- requirement identity inside provider context;
- required/recommended/optional/conditional applicability;
- data kind/cardinality/options/constraints where materially needed;
- editability/immutability when material;
- provider schema/version/checksum/revision evidence when available;
- acquisition/provenance/time.

This is not a generic Provider Requirement framework, registry or new business owner.

Readiness MUST NOT own:

- desired listing copy/text/media selection as a separate aggregate;
- Listing lifecycle;
- Price Intent;
- Sellable Availability;
- provider actual state;
- arbitrary provider payloads;
- Product master lifecycle.

---

# 9. Offering contract — ListingIntent is the one authoring aggregate

D2 already admits durable ListingIntent identity. B1 already allows owner-specific intent operations. This candidate reuses that existing concept instead of adding PublicationPreparation.

A ListingIntent may begin as an editable **draft** under Offering for one desired create/edit/close listing action.

Conceptual subject/correlation includes:

```text
Organization
+ Marketplace Installation
+ SourceInstance-qualified Product
+ existing source-qualified Listing where editing/closing an existing listing
```

Exact path/addressing remains D5.

The ListingIntent draft may carry only data whose meaning belongs to the desired listing representation/action, for example:

- target category/product-type selection where Offering must choose the intended listing representation context;
- channel-facing content values needed by the provider;
- media selection/order/role for the listing;
- provider-qualified requirement resolutions used for this desired representation;
- intended listing lifecycle action;
- references to accepted owner-specific price/availability/fulfillment inputs only where the concrete provider publication contract requires them.

It MUST NOT become Product master, Economics authority, Sellable Availability authority or Fulfillment authority.

### Creation and editing are one architecture

```text
Create listing
  → ListingIntent(target=no existing Listing)

Edit listing
  → ListingIntent(target=existing source-qualified Listing)
```

Provider differences determine adapter protocol, not a second MPC authoring model.

---

# 10. Publication value resolution — baseline only

No generic rule/mapping DSL is admitted.

For a publication value legitimately carried by ListingIntent, baseline resolution has only two modes.

## `FOLLOW_SOURCE`

The ListingIntent references one explicit source-qualified fact/candidate admitted by Readiness.

Rules:

- source identity/provenance is explicit;
- current value follows the selected source fact until the intent is frozen/submitted;
- source unavailable/unknown stays unavailable/unknown;
- a source change triggers readiness/current-intent revalidation where material;
- no silent fallback to a manual value.

## `EXPLICIT_OVERRIDE`

The ListingIntent contains an MPC-authored channel value.

Rules:

- effective Principal/automation and time are recorded;
- override never mutates or falsifies source Product truth;
- source change does not silently replace the override;
- source drift may make the draft stale/conflicting if Readiness/Offering semantics judge it material;
- removing the override may return the draft to FOLLOW_SOURCE only through an explicit edit.

### No baseline `DERIVED`

A generic derived/transformation mode is rejected for Product 1.0 baseline.

If repeated concrete cases later prove a deterministic owner-specific transformation materially reduces total complexity, admit the smallest explicit transformation with a named owner and proof. Do not prebuild feed rules, expression DSL, mapping engine or formula language.

---

# 11. Source acquisition / integration ingress

Source acquisition remains D4 mechanism/evidence, not a new Product domain.

## Embedded source adapter

Example: Sankhya adapter inside MPC.

```text
Sankhya sanctioned API
    ↓
D4 Sankhya adapter
    ↓
consumer-owned semantic ports / D4 acquisition translation
```

No mandatory self-HTTP loopback.

## Future external connector

A separately deployed connector may enter through a bounded Integration Ingress HTTP adapter when a concrete consumer requires it.

```text
external connector
    ↓ authenticated source-specific/bounded ingress
D4 HTTP inbound adapter
    ↓ same translation/admission logic
consumer-owned semantic ports
```

This candidate prepares the seam but **does not freeze or require a generic Source Ingestion HTTP API now**.

When the first real external connector is admitted, D5 may define the smallest contract necessary for that source class, preserving:

- Organization + SourceInstance binding;
- native Product identity;
- observation provenance/time/revision where material;
- partial/full/absent semantics;
- replay/duplicate safety;
- no ability to impersonate MPC-authored ListingIntent values.

No generic `/entities`, `/resources`, plugin registry or connector platform is introduced.

## Agent/manual authoring

Claude/Codex/Fable-style automation is a D2 automation Principal and uses the **Semantic Product API** to edit Offering-owned ListingIntent drafts. It does not use Integration Ingress and cannot claim a value came from Sankhya.

---

# 12. Media/images

Media participates in listing authoring without creating Product media authority.

Possible origins:

- source-qualified Product media from Sankhya/another external source;
- MPC-authored/uploaded media for this listing context.

Rules:

- source media remains source evidence;
- MPC-uploaded media never rewrites the source Product;
- ListingIntent owns only the listing-specific selection/order/role desired for publication;
- provider image IDs/CDN/resource topology remains provider-local;
- the consequential ListingIntent snapshot preserves enough media/provenance to explain what was attempted;
- blob/object storage, hashing, resizing, caching/CDN and upload mechanics remain D7/implementation;
- arbitrary external URL strings are not automatically trusted media.

No `ProductAsset`/`AssetFamily`/PIM media hierarchy is introduced without new evidence.

---

# 13. Consequential snapshot / history

Before provider dispatch, the material ListingIntent becomes immutable enough to explain the attempted action.

The historical snapshot preserves, proportionately:

- intent identity/action/target scope;
- source-qualified Product + Installation + target Listing when applicable;
- exact desired listing values;
- per-value resolution source (`FOLLOW_SOURCE` reference or `EXPLICIT_OVERRIDE` + Principal provenance);
- provider requirement/schema revision/provenance materially used by readiness/validation;
- media selection/provenance;
- decision-time readiness/disposition/authorization references;
- correlation to provider attempt/result/convergence.

This does not convert source facts into MPC current authority. It is decision/action history.

Exact snapshot/reference persistence strategy remains D2/D7 realization as long as historical explanation survives later source/provider changes.

---

# 14. Price, Availability and Fulfillment authority fence

A provider may physically combine content, price, quantity and fulfillment in one request. That does not merge MPC owners.

## Price

Offering already owns Price Intent. Listing creation may correlate to a valid current Price Intent when the provider requires initial price. Economics informs; adapter does not calculate/choose price.

## Availability

Availability owns Sellable Availability and Availability Intent.

Preferred shape:

```text
Listing representation can be established without inventing Sellable Availability
        ↓
Availability observes Offering target
        ↓
Availability converges through its own owner path
```

But this is **not yet assumed true for the selected first Mercado Livre creation lane**.

### OPEN GATE — R1-G1: Mercado Livre initial publication × Availability

Before D4-R1 can be ratified, establish from current official/real non-destructive evidence whether the selected current Mercado Livre User Product publication lane can create the required representation without an Offering/Readiness-owned fabricated quantity.

Acceptable closure:

**PASS-A — separate progression works**

- selected creation path establishes representation/product/offer without requiring Offering to own Sellable Availability;
- Availability can converge afterward through its accepted path.

**PASS-B — provider physically requires quantity but existing authority can supply it without semantic distortion**

- the exact accepted owner dependency is already legal under D1/D3 and adapter only composes an already-authoritative Availability meaning.

**FAIL / REOPEN**

- selected initial publication semantically requires a new `Availability → Offering` or other cross-owner dependency absent from D1/D3;
- correctness requires Offering/Readiness to own quantity;
- provider contract forces atomic cross-owner semantics not representable by current boundaries.

In FAIL, reopen only implicated D1/D3 decisions before D5-B2 resumes.

No live marketplace write is authorized by this gate. Use current official docs, read-only evidence, sandbox/protocol validation or another non-destructive proof unless the operator separately authorizes a live effect.

## Fulfillment

Provider shipping/fulfillment requirements remain D4 evidence and Fulfillment retains business responsibility. ListingIntent may reference a Fulfillment-owned semantic input if the concrete provider contract requires it; it does not absorb Fulfillment policy/state.

---

# 15. Provider execution / convergence

Offering owns ListingIntent and convergence conclusion.

D4 marketplace adapter receives already resolved/authorized inputs and owns:

- exact provider endpoint/resource choice;
- DTO serialization;
- provider auth/protocol;
- provider-specific prerequisite translation;
- response interpretation no stronger than evidence;
- authoritative reread/reconciliation surfaces.

Adapter MUST NOT:

- choose between competing source/override values;
- invent missing required values;
- query hidden systems to fabricate semantic readiness;
- copy provider DTOs into a generic Listing model;
- turn transport 2xx into convergence;
- widen intended/authorized scope because provider resources are shared;
- recalculate Price/Sellable Availability/Fulfillment meaning.

Provider actual Listing/User Product/Catalog/etc. remains external authority.

---

# 16. Credible alternatives

## A — Sankhya-only, no MPC authoring

Rejected. Current product can need publication before all source e-commerce fields are complete; it also creates a hidden source-data-completion launch gate.

## B — MPC Product/PIM master

Rejected. Duplicates external Product authority and drifts toward Akeneo/commerce-hub semantics outside D0/D1.

## C — generic provider field bag

Rejected. Hides provider DTO leakage and destroys authority/provenance/type evolution.

## D — Readiness-owned PublicationPreparation aggregate + Offering ListingIntent

Rejected by Global Coherence. Two adjacent desired-publication authorities and creation/edit drift are foreseeable.

## E — generic SourceProductObservation service/domain

Rejected. D4/source acquisition mechanism would acquire ownership of a payload that may feed multiple consumer-owned semantics; conflicts with D4 final coherence.

## F — every embedded adapter calls public HTTP ingress

Rejected baseline. Self-HTTP adds network/auth/deployment coupling without changing authority. External connector HTTP remains legitimate when physically external.

## G — Readiness requirements + Offering ListingIntent draft + D4 translation

**Recommended Global Maximum.** Reuses existing identities/owners, supports current manual/agent authoring and future Sankhya dominance, keeps source/provider distinctions honest and minimizes implementation objects.

---

# 17. Complexity / YAGNI

## Essential complexity retained

- source-qualified Product facts;
- provider/category/product-type-specific requirements;
- requirement/schema churn/version evidence;
- Readiness sufficiency;
- ListingIntent authoring/history;
- source-follow versus explicit channel override;
- media origin/selection;
- provider serialization/convergence;
- price/availability/fulfillment authority separation;
- agent attribution;
- source acquisition provenance;
- initial publication/Availability proof gate.

## Accidental complexity removed/refused

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
- AI-specific API/backdoor;
- creation architecture separate from edit architecture;
- one giant Publication aggregate owning price/stock/fulfillment.

---

# 18. Whole-product coherence checks

This seam must remain coherent with the existing system rather than create an alternate architecture.

## Identity

- Product remains external source-qualified identity.
- Listing remains external source-qualified identity.
- ListingIntent is the existing MPC-owned action identity.
- no new ProductDraft/PublicationDraft identity is presumed.

## Communication

- Readiness→Offering remains Q.
- source acquisition remains D4 boundary/mechanism.
- provider execution remains owner intent → D4 effect → authoritative reread.
- no new workflow bus/command taxonomy.

## Work/Governance

- missing/conflicting listing requirements may become explicit Work under existing semantics when materially actionable;
- assignment never becomes authorization;
- Governance authorizes consequential ListingIntent where policy requires it;
- neither Work nor Governance owns listing content/business validity.

## Historical explainability

- mutable source facts do not rewrite past ListingIntent basis;
- current provider requirements do not rewrite the schema/requirements used by a prior attempted publication.

## Implementation posture

- compatible with modular monolith + one PostgreSQL target;
- does not require microservices, broker, schema registry, PIM engine or workflow engine;
- D7 can choose minimal storage/runtime mechanisms later.

---

# 19. Proof strategy before implementation

## P1 — incomplete source → authored listing → future source takeover

1. Sankhya Product lacks one required publication value.
2. Readiness reports requirement missing and exposes source candidates honestly.
3. automation Principal creates/edits ListingIntent with EXPLICIT_OVERRIDE.
4. readiness becomes sufficient for that intent/context.
5. consequential intent snapshot preserves override provenance.
6. later Sankhya supplies the value.
7. old ListingIntent remains unchanged/history-valid.
8. a new edit intent may explicitly FOLLOW_SOURCE.

## P2 — source spoofing

- Product API client cannot claim EXPLICIT_OVERRIDE came from Sankhya.
- source connector/acquisition path cannot author ListingIntent override by choosing `source=manual`.

## P3 — partial source observation

Omitted field in partial acquisition does not clear known source fact unless source semantics establish removal/absence.

## P4 — provider requirement churn

Requirement revision N supports one historical intent. Revision N+1 changes requiredness/enumeration. Historical intent remains explainable under N; current draft/readiness evaluates against N+1.

## P5 — one architecture for create/edit

Create and edit differ only by target Listing existence/action semantics; both use Offering ListingIntent. No second publication object appears.

## P6 — second marketplace structural inversion

Use one source Product against:

- Mercado Livre User Product/Item;
- Amazon Product Type Definition + Listings Item product/offer modes;
- Mirakl Product/Offer/PriceStock split.

Readiness/Offering/Product identity must remain stable; only D4 requirements/serialization differ.

## P7 — media

Source media and MPC-uploaded media can both be selected without changing Product authority; intent history preserves attempted selection.

## P8 — embedded versus external source connector

Embedded Sankhya adapter and a hypothetical external connector reach equivalent consumer-owned semantic results without requiring one shared network hop or one new business owner.

## P9 — adapter truth-selection negative

Remove required resolution from ListingIntent/Readiness context. Adapter must refuse pre-dispatch/return owner insufficiency; it may not select fallback/default.

## P10 — agent path

Automation Principal uses normal Product API authoring and remains subject to Permission, owner validity, Governance, idempotency, audit and convergence.

## P11 — R1-G1 Availability gate

Prove selected first ML publication lane satisfies PASS-A or PASS-B above; otherwise name the exact D1/D3 contradiction.

---

# 20. Independent Global Coherence challenge package

Fable must reconstruct repository authority independently and attack the whole architecture, not merely this text.

Material attacks:

1. Does ListingIntent-as-draft genuinely fit accepted D2/D1 meaning or silently expand Intent into a Product/PIM document?
2. Does moving authoring to Offering fully remove Readiness/Offering overlap?
3. Is Readiness-owned publication-requirements Q sufficient, or is an independent publication-content authority actually missing?
4. Does `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` preserve enough future evolution without a mapping engine?
5. Is removing `DERIVED` correct YAGNI or does a present Sankhya/JET consumer already require deterministic transforms?
6. Does D4 acquisition feeding multiple consumer-owned ports remain implementable without a generic SourceProductObservation authority?
7. Is a future external Integration Ingress seam enough, or does a real present consumer require HTTP Source Ingestion now?
8. Can an external connector be authenticated/bound to Organization+SourceInstance without creating a generic integration registry?
9. Does agent automation remain an ordinary Principal rather than an AI-specific authority path?
10. Can current Mercado Livre User Product creation/edit flows fit the seam without provider DTOs leaking into ListingIntent?
11. Close or refute R1-G1: initial ML publication × Availability.
12. Does current ML shared User Product blast radius require different intent scope semantics?
13. Does provider requirement/schema versioning add only essential complexity?
14. Does media need a new D2 identity, or are bounded content/provenance references enough?
15. Does ListingIntent snapshot copy too much source/provider state, creating a second authority rather than history?
16. Does creation/edit lifecycle need one intent type or separate owner-specific intents?
17. Can Amazon product-only/offer-only modes fit without moving Product authority into MPC?
18. Can Mirakl Product/Offer/PriceStock split fit without one giant Publication aggregate?
19. Are ANYMARKET/Akeneo-like PIM capabilities actually a missing Product requirement, or correctly out of scope?
20. Does Google primary/supplemental precedent justify only explicit precedence or a real rule engine now?
21. Does any D1 boundary duplicate another after this change?
22. Are Work and Governance still distinct and non-platform business authorities?
23. Does Materialization↔Fulfillment remain coherent under the same intent/observation pattern used by Offering?
24. Does the whole system still fit a modular monolith without requiring microservice/distributed-workflow machinery?
25. Does any candidate concept exist only because another abstraction exists?
26. Run Structural Inversion: if current implementation/OpenAPI were opposite, do these conclusions still follow?
27. Hardest future change: Sankhya becomes complete source; second ERP; second marketplace; provider schema churn; agent-heavy operations — does any require authority migration?
28. Is any D0/D1/D2/D3/D5-B1 reopen genuinely required by evidence rather than preference?

Reviewer severity does not create authority.

---

# 21. Reopen / stop triggers

Reopen only the implicated parent decision when material evidence proves:

1. ListingIntent cannot legitimately carry editable desired listing representation before consequential submission → D1/D2 targeted review.
2. Readiness cannot provide requirements/sufficiency without owning desired listing state → D1.
3. R1-G1 requires a new Availability/Fulfillment semantic edge absent from D1/D3 → D1/D3.
4. publication correctness requires atomic cross-owner mutation → D1/D3.
5. a new durable identity/lineage class is required beyond Product external reference + ListingIntent → D2.
6. provider requirement/effect cannot be represented through consumer-owned ports without raw DTO leakage → D4 redesign.
7. external source connector cannot preserve Organization/SourceInstance/provenance safely → D2/D4/D5 boundary review.
8. repeated real source→listing transformations prove a reusable mapping primitive essential → targeted Offering/Readiness/D4 decision, not generic PIM.
9. a second real Product master changes source-authority assumptions → D1/D2.
10. a real external connector creates a materially different API compatibility/security requirement → D5.

Framework preference, provider symmetry, current-code convenience and hypothetical providers are not reopen evidence.

---

# 22. Candidate outcome

**Proposed outcome:** `CURRENT D0→D5-B1 STRUCTURE CONFIRMED` + `RESTRUCTURE NOW` for the missing D4 publication-input seam.

If R1-G1 closes, independent Fable challenge converges and the operator explicitly ratifies the resulting package, canonical D4 should add only the following durable meaning:

- Readiness owns publication requirements, correspondence, source candidates and sufficiency;
- Offering-owned ListingIntent is the single create/edit authoring/draft identity;
- baseline listing-value resolution is `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` only;
- agent/manual authoring uses normal Product API + D2 Principal semantics;
- external source acquisition remains D4 evidence/mechanism feeding consumer-owned semantic ports, not a Product/source-observation domain;
- embedded adapters do not require self-HTTP;
- external connector ingress remains a prepared D4/D5 seam until a concrete connector requires a wire contract;
- provider requirement/schema evidence remains source-qualified and consumer-bounded;
- media selection belongs to ListingIntent only as desired listing representation, without creating Product media master;
- ListingIntent consequential snapshot preserves decision-time source/override/requirement provenance;
- adapter translates resolved owner meaning and never selects truth;
- Price/Availability/Fulfillment authorities remain intact;
- no MPC Product Master/PIM, PublicationPreparation aggregate, generic ProviderRequirement framework, generic SourceProductObservation service, generic transformation DSL, connector framework or alternate AI architecture is introduced;
- D5-B2 resumes only after this contract and any genuine parent-stage contradiction are closed.

**D5-B2 review candidate remains non-authoritative and frozen. Implementation remains blocked until D9.**
