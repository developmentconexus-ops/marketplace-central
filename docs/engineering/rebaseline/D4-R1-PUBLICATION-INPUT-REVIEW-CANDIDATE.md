# D4-R1 — Publication Input & Listing Authoring Contract — REVIEW CANDIDATE

> **Status:** NON-AUTHORITATIVE TARGETED-REOPEN REVIEW CANDIDATE  
> **Stage:** targeted D4 reopen while D5-B2 remains NEXT / NOT YET OPENED  
> **Parent authority:** accepted D0 → D4 + D5-B1  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Prepared at:** `8a01920786184d517767498c61084ca7908cf577` on `docs/global-methodology-alignment`  
> **Purpose:** close a publication/input-contract gap discovered while deriving D5-B2. This file does not change D4/D5 authority, router status, ADR status or implementation permission.

---

## 1. Why D5-B2 stops here

D0 already requires a Product 1.0 actor to take an eligible source Product through channel readiness and **creation/publication → real channel observation**.

D1 already says:

- Product master remains external;
- Product & Channel Readiness owns Product↔channel correspondence, channel requirements and supported/missing/conflicting/readiness meaning;
- Marketplace Offering Operations owns marketplace offer/listing representation and lifecycle, Listing Intent, Price Intent and convergence;
- Availability owns Sellable Availability and Availability Intent;
- adapters own provider protocol, never business authority.

D2 already says:

- Product is `SourceInstance + native Product key`;
- provider Listing/Variation is external source-qualified identity;
- durable Listing/Price Intents may have owner-local MPC identity;
- no MPC Product mirror exists merely for normalization.

D4 already says:

- consumer owns meaning; adapter owns protocol;
- provider category/catalog/User Product/etc. topology stays provider-local;
- Offering owns Listing Intent/convergence;
- provider-specific evidence may survive only as source-qualified evidence serving a named consumer/correctness property.

What is **not** yet decided is the contract between an external Product source and a publishable Listing Intent:

- which publication inputs come from an external source such as Sankhya;
- which values may be authored/overridden inside MPC;
- how images/media enter when sourced externally or authored/uploaded in MPC;
- how category/product-type and changing provider requirements are represented;
- how source values, explicit overrides and unknowns coexist without one silently replacing another;
- how external systems push source observations into MPC;
- how a built-in Sankhya adapter and an external push integration share one semantic ingestion path without self-HTTP duplication;
- what exact resolved snapshot a Listing Intent freezes before provider publication;
- how an adapter receives enough provider-specific requirement resolution to publish without becoming the authority that selects truth;
- how later source enrichment can replace temporary manual authoring deliberately rather than by accidental last-write-wins.

D5-B2 cannot safely freeze `ListingRequirements`, `ListingIntent` request shapes or external input APIs until this seam is decided.

---

# 2. Decision question

> **What is the smallest sustainable contract that lets MPC create/edit marketplace publications when Product data may come from Sankhya or another external source, may temporarily be completed by humans/automation through MPC, and must ultimately be translated into provider-specific publication requirements without creating an MPC Product Master, generic PIM, raw provider-field bag or duplicate authority?**

The contract must support all of these present/future facts without changing product identity:

1. today, some source Product information is still incomplete or under human review;
2. an operator or automation Principal may need to complete publication data directly in MPC so publication can proceed;
3. Claude Code/Codex/Fable-style automation may legitimately call MPC APIs to author those MPC-owned publication inputs;
4. Sankhya is expected to become the dominant Product source over time, including e-commerce/product presentation data and images where available;
5. MPC may later consume source Product observations through a built-in adapter, an external connector, or a push integration;
6. each marketplace retains materially different and changing publication requirements;
7. provider protocol/resources must remain adapter-local.

---

# 3. Evidence classification

## 3.1 KNOWN from repository authority

1. Product 1.0 requires actual marketplace creation/publication, not read-only listing observation.
2. MPC is not an ERP replacement or generic integration hub.
3. Product master remains externally authoritative and source-qualified; no MPC Product Master/PIM domain exists.
4. Readiness owns Product↔channel correspondence, provider/channel requirements and readiness conclusion.
5. Offering owns listing representation/lifecycle, Listing Intent and convergence.
6. Availability owns Sellable Availability; Offering may not absorb stock authority merely because a provider publication payload happens to carry quantity.
7. Price Intent is Offering-owned; Economics informs but never writes marketplace price.
8. Provider DTO/category/catalog/User Product/warehouse topology never becomes MPC ontology by normalization.
9. Externally authoritative facts preserve source qualification/provenance/knowledge state.
10. A valid consequential intent preserves decision-time evidence/provenance sufficiently for later explanation.
11. D5-B1 fixes a semantic Product API, separate provider/business-system protocol ingress, source-qualified external identity, fail-closed consequential idempotency and one machine-readable Product API wire authority.
12. D5-B1 refused a speculative third technical/admin surface but explicitly left room for one later when a real concrete consumer proves it.
13. Current D5-B2 candidate is non-authoritative and may be revised without compatibility cost.

## 3.2 External benchmark evidence — non-authoritative

Current official/provider/reference research supports these structural observations:

### Mercado Livre

- current User Products change publication structure materially compared with the legacy Item/Variation model;
- product/family attributes and item sales conditions are no longer one stable flat listing shape;
- current documentation states that some User Product fields are inherited/generated and must not be sent in the same way as legacy publication;
- category technical-specification surfaces expose required attributes that can change by category/context;
- edits to shared User Product characteristics may propagate across related items asynchronously;
- stock topology may require User Product/location-specific endpoints rather than a universal item quantity field.

Primary sources:

- https://developers.mercadolivre.com.br/pt_br/publicacao-de-produtos/user-products
- https://developers.mercadolivre.com.br/pt_br/api-docs-pt-br/atributos
- https://developers.mercadolivre.com.br/pt_br/autenticacao-e-autorizacao/gestao-de-estoque-multiorigem-user-products

### Amazon SP-API

- Product Type Definitions returns marketplace/product-type-specific JSON Schema, including required/conditional attributes and a version;
- product facts and sales terms can be requested separately (`LISTING_PRODUCT_ONLY`, `LISTING_OFFER_ONLY`) or together;
- Listings Items uses the resulting schema for create/update and can represent issues/attributes/offers/fulfillment availability;
- listing requirements/enumerations change over time, so static provider-field assumptions are unsafe.

Primary sources:

- https://developer-docs.amazon.com/sp-api/lang-en_EN/docs/manage-product-listings-guide
- https://developer-docs.amazon.com/sp-api/lang-en_EN/docs/building-listings-management-workflows-guide

### Mirakl

- catalog configuration distinguishes standard fields from channel-specific custom offer attributes;
- channel connectors import channel catalog/taxonomy configuration;
- synchronization distinguishes Product, Offer and Price/Stock upsert flows;
- connectors report integration feedback after channel processing.

Primary sources:

- https://developer.mirakl.com/content/product/connect-channel-platform/developer-guide/catalog-configuration
- https://developer.mirakl.com/content/product/connect-channel-platform/developer-guide/catalog-flow

### Google Merchant Center

- primary and supplemental product data sources are explicitly distinct;
- supplemental sources may enrich/override missing or channel-specific information and can be supplied through API/file/Sheets;
- attribute rules make source precedence explicit rather than relying on accidental write order.

Primary sources:

- https://support.google.com/merchants/answer/14990942
- https://support.google.com/merchants/answer/14994083

### Hub/operations benchmarks

ANYMARKET, Channable and Linnworks demonstrate the real operational need for channel attribute mapping/configuration and manual/source enrichment, but they also demonstrate the risk of letting a marketplace hub become a broad Product/PIM master. MPC uses them only as competitive/operational evidence, not target authority.

### Sankhya / JET evidence

- Sankhya currently documents Product image retrieval for external applications through a Gateway/legacy image endpoint;
- JET is an official Sankhya e-commerce partner advertised with Sankhya integration and broad marketplace connectivity;
- exact JET internal synchronization semantics are not treated as authority here.

Primary sources:

- https://ajuda.sankhya.com.br/hc/pt-br/articles/36396748479383-Obten%C3%A7%C3%A3o-de-Imagens-de-Produtos-via-API
- https://www.sankhya.com.br/gestao-de-negocios/parceiros-sankhya-erp/encontre-um-parceiro/jet-e-commerce/

## 3.3 INFERRED

1. A fixed universal `MPCListingDTO` containing every marketplace field is structurally unstable.
2. A generic `map[string]any`/`extensions` field only hides provider DTO leakage rather than solving it.
3. A single API that accepts both source facts and MPC-authored overrides would blur authority/provenance.
4. Built-in pull adapters and external push connectors should converge on one semantic ingestion use case/port, but they need not use the same physical transport.
5. A built-in Sankhya adapter calling MPC's own public HTTP API merely to reach an in-process use case adds accidental self-HTTP unless deployment separation requires it.
6. Temporary manual/agent publication authoring and future source-driven publication can coexist sustainably if each resolved value preserves its origin and selected resolution mode.
7. Provider requirement schemas/descriptors must be version/provenance aware enough that a Listing Intent can explain what requirements it satisfied at decision time.
8. The current D1 boundary set may be sufficient if mutable **Channel Publication Preparation** is treated as Readiness-owned requirement-resolution/preparation state rather than live Listing representation; independent review must attack this boundary carefully.

## 3.4 UNKNOWN

- exact current Mercado Livre create/edit payload for every User Product branch the first real lane may encounter;
- whether the first selected Mercado Livre creation path can always create a non-sellable/product-only representation before Availability convergence;
- exact persisted representation/identity, if any, for Channel Publication Preparation;
- exact media/blob storage/caching realization;
- exact generic semantic fields, if any, that deserve promotion across multiple providers;
- whether reusable cross-product Source→Requirement mappings are needed now or only after repeated manual resolution proves the abstraction;
- exact Source Ingestion HTTP schema/auth/batch shape — D5/D7 after semantic contract;
- exact API operation names/paths — D5;
- exact UI/editor — D6;
- runtime scheduling/polling/checkpoint/media cache — D7.

Unknown stays Unknown.

---

# 4. Root cause

D4 currently begins publication too late.

It knows provider requirements/protocol and external effect semantics, while D1/D2 know Product identity, readiness and Listing Intent ownership, but no accepted contract explicitly says how source Product observations and temporary MPC-authored publication values become a provider-valid, historically explainable publication input.

That gap allows several defect classes:

- adapter becomes de facto truth selector (`if Sankhya field empty, invent/default/use UI value`);
- manual authoring silently overwrites source truth or vice versa;
- source outage/unknown becomes empty publication data;
- provider requirement IDs/constraints leak into a generic core field bag;
- MPC grows a duplicate Product Master/PIM by convenience;
- source-to-provider mapping becomes scattered across adapters/UI/scripts;
- Listing Intent cannot explain which requirement/version/value provenance it used;
- later Sankhya enrichment cannot safely supersede temporary overrides;
- agents/scripts bypass semantic authority by writing database/provider payloads directly;
- built-in and external connectors implement different business rules.

---

# 5. Target invariant

> **Source Product truth, MPC publication authoring, provider requirement evidence, Listing Intent and external Listing state remain distinct authorities. Every material provider publication input is resolved under an accepted MPC owner before protocol serialization; the adapter translates the resolved meaning but never chooses which source/override is true.**

Corollaries:

1. Source observation never becomes MPC-owned merely because it is persisted.
2. An MPC-authored override never masquerades as a Sankhya/source observation.
3. Provider requirement metadata never becomes universal Product ontology by existence alone.
4. A missing source value remains missing/unknown unless an explicit MPC-owned authored value legitimately supplies that publication requirement.
5. No implicit last-write-wins across source and authored values.
6. Provider-specific requirement resolution may remain source-qualified without raw DTO mirroring.
7. Listing Intent freezes enough resolved input + requirement/provenance context to explain the attempted publication later.
8. Price, Availability, Fulfillment and other accepted owners do not lose authority merely because a provider combines their data in one wire request.
9. Built-in adapters, external connectors and agent/API authoring use the same accepted semantic use cases even when physical transports differ.
10. Current implementation convenience never creates a second source of business truth.

---

# 6. Credible alternatives

## A — Sankhya-only Product source; no MPC authoring

Wait until every e-commerce/product field is complete in Sankhya and only then publish.

**Rejected.** It blocks a present Product 1.0 operating need and makes an external data-completion process a hidden launch dependency. It also provides no legitimate path for temporary operator/agent completion.

## B — MPC becomes Product/PIM master

Import Product records from Sankhya, own title/description/images/attributes centrally, then let every marketplace adapter select fields from that master.

**Rejected.** Contradicts D1/D2 Product authority and duplicates the external master. It is a Local Maximum resembling marketplace hubs/PIMs rather than MPC's accepted control-plane boundary.

## C — Generic provider field bag on Product or ListingIntent

Store arbitrary provider field names/JSON under `attributes`, `extensions` or `provider_payload` and let adapters interpret them.

**Rejected.** Provider protocol becomes domain state by renaming; no clear authority, schema evolution, provenance or validation boundary exists.

## D — Every built-in adapter must call MPC's public HTTP Source API

Sankhya pull adapter fetches Sankhya, then loops back through the external HTTP API even when running inside the same application boundary.

**Rejected as baseline.** It preserves one wire path at the cost of accidental network coupling, auth duplication and self-HTTP. If a connector is physically external, HTTP is appropriate; an embedded adapter should call the same semantic application port/use case directly.

## E — Separate source ingestion + MPC authoring over one semantic core

- external source facts enter through a source-qualified ingestion contract;
- operator/agent authored values enter through Product API capabilities owned by the appropriate MPC domain;
- built-in pull adapters and external push connectors converge on the same source-observation application use case;
- Readiness maintains publication requirement/preparation semantics without becoming Product Master or live Listing owner;
- Offering freezes a resolved publication snapshot into Listing Intent and owns external Listing convergence;
- D4 adapter serializes provider protocol only.

**Recommended Global Maximum.** Smallest structure that supports today + the stated Sankhya-dominant future without duplicate authority or generic PIM machinery.

---

# 7. Four semantic planes

## 7.1 Plane A — Source Product Observation

Authority: external SourceInstance.

Conceptual subject:

```text
Organization
+ SourceInstance
+ native Product key
```

Observation may include only facts a concrete consumer actually needs, for example names/descriptions, brand/model, GTIN, physical dimensions/weight, source media references and existing e-commerce custom fields where the source legitimately owns/exposes them.

Rules:

- source observation is qualified by SourceInstance + Product key;
- acquisition/source time and provenance are preserved where material;
- partial observation does not erase omitted facts;
- explicit source-null/absence is distinct from field-not-observed;
- unavailable acquisition never becomes empty Product data;
- no synthetic MPC Product identity is created;
- source-specific field names/protocol remain in the source adapter unless a consumer-owned semantic fact is actually needed.

This plane is external evidence, not MPC Product Master state.

## 7.2 Plane B — Channel Publication Preparation

Leading owner: **Product & Channel Readiness**, subject to independent boundary challenge.

Conceptual key:

```text
Organization
+ Marketplace Installation
+ SourceInstance-qualified Product
```

Purpose:

> hold the current preparation needed to answer whether this Product can be published/edited correctly in this channel context, including provider requirement resolution, missing/conflicting values and provenance.

This is not the live Listing and not external Product master.

It may contain proportionately:

- selected provider category/domain/product type reference;
- provider requirement descriptor/version/checksum/provenance;
- current resolution per admitted publication requirement;
- current source candidate where available;
- explicit MPC-authored override where legitimate;
- unresolved/missing/conflicting requirement state;
- media selection/order references needed for publication preparation;
- readiness conclusion/reasons.

It MUST NOT absorb:

- Sellable Availability authority;
- Price Intent/economic authority;
- Fulfillment responsibility;
- provider actual Listing state;
- generic source Product lifecycle;
- arbitrary provider payloads;
- universal PIM/category/attribute ontology.

If review proves that this state is actually Offering-owned desired Listing representation rather than Readiness preparation, D1 must be targeted-reopened rather than silently redefining the boundary here.

## 7.3 Plane C — Listing Intent publication snapshot

Authority: Marketplace Offering Operations.

A material create/edit/close publication action creates a Listing Intent containing or referencing an immutable decision-time snapshot sufficient to establish:

- target Organization + Marketplace Installation;
- source-qualified Product/correspondence;
- selected provider category/product-type context when material;
- resolved publication inputs that Offering is allowed to use;
- provider requirement descriptor/version/provenance needed to explain validity;
- Offering-owned desired lifecycle/commercial listing semantics;
- applicable Price Intent/reference where the action requires it;
- intended target scope/blast radius;
- business disposition + authorization context lineage;
- idempotency/correlation anchor;
- resulting provider acceptance/convergence evidence later.

The snapshot does not turn every source fact into MPC authority. It is historical decision context.

## 7.4 Plane D — Provider execution / external Listing

Authority split:

- Offering owns Listing Intent, desired semantic action and convergence conclusion;
- D4 adapter owns provider protocol/DTO/resource selection;
- provider owns actual Listing/User Product/Catalog/etc. current state.

Adapter receives only already-resolved semantic/provider-qualified inputs required by its publication contract.

Adapter MUST NOT:

- pick between competing source/manual values by convenience;
- invent missing required values;
- query arbitrary extra systems to fill semantic gaps secretly;
- turn 2xx into convergence;
- rewrite owner meaning to fit provider topology.

---

# 8. Two input authorities, not one generic write API

## 8.1 Source Ingestion contract — machine/source authority

This is a D4 source-integration boundary, not ordinary Product authoring.

A source push or built-in pull path submits **Source Product Observations**.

Minimum semantic obligations:

- authenticated/configured Organization + SourceInstance binding;
- source cannot choose another Organization merely in payload;
- explicit native Product key;
- observation/revision/occurrence discriminator sufficient for duplicate safety where needed;
- source/observation time where available/material;
- partial/full observation semantics explicit;
- known/unknown/absent distinctions preserved;
- media may be represented by source-qualified reference/fetch contract rather than copied blindly;
- replay of the same observation does not create duplicate Product meaning;
- no Business Intent/Listing creation is implied merely by source ingest.

### Physical entry modes

The same semantic source-observation use case may be reached by:

1. **built-in pull adapter** — e.g. MPC Sankhya adapter reads sanctioned Sankhya APIs and calls the application ingestion port directly;
2. **external push connector** — ERP/JET/other integration calls an MPC Source Ingestion HTTP contract;
3. **bounded file/import adapter** — only if a real operational source requires it; file format is mechanism, not Product API ontology.

The first two are present evidence; a generic plugin framework is not introduced.

## 8.2 Product authoring contract — MPC authority

A human or automation Principal edits MPC-owned publication preparation/owner configuration through the Semantic Product API.

Examples:

- supply a missing marketplace-specific description;
- select/confirm category/product-type evidence;
- resolve a required provider attribute not yet supplied by Sankhya;
- select/reorder publication media;
- explicitly override a source-derived publication value where business rules allow it;
- clear an override later so the preparation follows the source again.

Rules:

- authenticated Principal/automation attribution is explicit;
- ordinary Permission does not imply publication authorization;
- authored values are MPC-owned preparation decisions, not source observations;
- an automation Principal (Claude/Codex/Fable-style agent) uses the same Product API authority and audit path as another MPC client;
- direct database writes/provider calls are not a legitimate authoring API;
- consequential publish action still occurs only through Listing Intent/Governance/effect-safety semantics.

---

# 9. Source-following versus explicit override

A global `last write wins` rule is rejected.

For a publication requirement that Readiness legitimately owns resolving, the smallest baseline modes are:

### `FOLLOW_SOURCE`

The preparation resolves from one explicitly selected source-qualified fact/reference.

- current value/provenance follows that source observation;
- source change causes readiness/revalidation, not a hidden new manual value;
- source unavailable/unknown remains unavailable/unknown;
- no silent fallback to a stale/manual candidate unless an explicit rule admits it.

### `EXPLICIT_OVERRIDE`

The preparation uses an MPC-authored value for this channel context.

- records effective Principal/automation + time/provenance;
- does not mutate or falsify source Product truth;
- survives later source updates until explicitly cleared/replaced;
- source drift may be surfaced if materially relevant;
- clearing the override may return the requirement to FOLLOW_SOURCE when a valid source binding exists.

### `DERIVED` — admitted only when a named deterministic owner rule exists

A derived value may be produced by a small accepted transformation whose owner and inputs are explicit.

This does **not** authorize a generic expression/rules DSL.

If no resolution exists, the requirement remains unresolved/unknown. There is no `DEFAULT_PLAUSIBLE` mode.

### Future Sankhya takeover path

Temporary authored values therefore do not create migration debt:

```text
Today:
Sankhya fact missing
  + EXPLICIT_OVERRIDE authored in MPC
  -> publication can become ready

Later:
Sankhya begins supplying sufficient fact
  -> source candidate becomes available
  -> operator/automation deliberately clears override / selects FOLLOW_SOURCE
  -> future preparation follows Sankhya
```

No mass migration from an MPC Product Master is required because none exists.

---

# 10. Provider requirements are evidence, not universal fields

D4 adapter acquires provider publication requirements for the concrete context.

A **Provider Requirement Descriptor** may preserve as much as materially needed:

- provider / Marketplace Installation / site or marketplace context;
- category/domain/product type identity;
- provider requirement identifier;
- required/recommended/optional applicability;
- data type/cardinality/allowed values/constraints/conditionality;
- provider editability/immutability where material;
- requirement/version/checksum/schema identity when exposed;
- acquisition time/provenance;
- provider ownership.

Examples:

- Mercado Livre category technical specs / User Product domain attributes;
- Amazon Product Type Definition JSON Schema + productTypeVersion;
- Mirakl category/offer custom-attribute configuration.

The descriptor is source-qualified provider evidence. It is not a new universal `ProductAttribute` master.

Readiness decides whether the requirement is sufficiently resolved for the Product/channel use. The adapter does not.

### Semantic-core promotion test

A provider field becomes a shared MPC semantic concept only after evidence proves:

1. the meaning is genuinely provider-independent;
2. at least one accepted owner needs that meaning independent of protocol;
3. normalization reduces total complexity rather than hiding differentiated semantics.

Otherwise it stays a provider-qualified requirement resolution.

---

# 11. Media / images

Images are explicitly part of this seam because current Sankhya and marketplace publication flows both expose them.

Publication preparation may reference media from at least two legitimate origins:

1. **source media** — e.g. a Sankhya Product image exposed by the SourceInstance adapter;
2. **MPC-authored/uploaded publication media** — supplied by an authorized Principal/automation for the channel context.

Rules:

- source media provenance remains source-qualified;
- MPC upload does not rewrite source Product media;
- preparation owns selection/order/role only where that meaning is publication-specific;
- a provider's image ID/thumbnail/CDN representation remains provider-local external state;
- Listing Intent freezes the selected media set/provenance required to explain the attempted publication;
- exact object/blob storage, hashing, resizing/cache/CDN and upload transport are D7/implementation decisions;
- raw arbitrary external URLs are not automatically trusted publication media merely because the field is a string.

This preserves a future in which most media follows Sankhya while allowing present publication before every source record is fully enriched.

---

# 12. Price, availability and fulfillment do not disappear into publication preparation

Provider APIs may physically combine product content, price, quantity, shipping/fulfillment and offer conditions in one request. That does not merge MPC authorities.

### Price

Offering already owns Price Intent. Publication execution may correlate to a current Price Intent/Offering-owned price decision when the provider requires initial price, without moving pricing analysis to Readiness or adapter.

### Availability

Availability Control owns Sellable Availability and Availability Intent.

Baseline preference:

> use a provider publication path that can establish the product/listing representation without inventing sellable quantity, then let Availability converge through its own accepted path.

Current Amazon evidence explicitly supports product-only listing submission distinct from offer terms. Current Mercado Livre User Products evidence separates/inherits stock differently from the legacy item model in some flows.

However this is a **proof obligation, not a universal assumption**.

If the selected first Mercado Livre publication effect materially cannot be performed without a simultaneous Availability-owned value and no existing accepted edge can supply it without authority distortion, STOP and targeted-reopen D1/D3. The adapter may not secretly query/decide Availability just because the provider payload requires it.

### Fulfillment / shipping

Provider-specific shipping/fulfillment requirements remain D4 evidence and the accepted owner retains semantic responsibility. Publication preparation may surface a requirement/dependency without absorbing Fulfillment authority.

---

# 13. Built-in adapter versus external API — decided direction

## Embedded/built-in adapter

Example: target Sankhya adapter inside MPC.

```text
Sankhya API Gateway
    ↓ D4 Sankhya adapter
Source Product Observation use case / port
    ↓
consumer-owned source observations
```

No mandatory loopback HTTP.

## External connector

Example: a separately deployed JET/ERP connector or future partner integration.

```text
External source/connector
    ↓ authenticated MPC Source Ingestion HTTP contract
HTTP inbound adapter
    ↓
SAME Source Product Observation use case / port
```

## Agent/manual authoring

```text
Claude/Codex/operator
    ↓ Semantic Product API
Principal-authenticated Product API adapter
    ↓
Readiness/Offering owner use case
```

The three paths share semantic core/use cases but preserve different authority/provenance.

> **Reuse the use case/port, not necessarily the network hop.**

This is ports/adapters applied to input as well as marketplace output.

---

# 14. Why one universal inbound API is rejected

A tempting API shape would be:

```text
POST /products/{id}/fields
{
  "title": "...",
  "images": [...],
  "attributes": {...},
  "source": "sankhya|manual|agent"
}
```

Rejected because:

- creates a synthetic MPC Product authority;
- lets callers self-declare source authority;
- mixes observation and authoring semantics;
- invites arbitrary provider field bags;
- makes source precedence a payload convention;
- obscures Organization/SourceInstance/channel scope;
- cannot preserve D1 ownership of price/availability/etc.;
- makes later Sankhya migration a conflict-resolution project.

The target separates **source observation** from **publication authoring** structurally.

---

# 15. API implications for D5 — not final routes yet

If this candidate survives review/ratification, D5 must recognize at least these external surface classes:

1. **Semantic Product API** — ordinary MPC clients/agents; Organization path scope; Principal-authenticated; owner-specific authoring/operations.
2. **Source Ingestion API** — machine-to-machine integration surface for SourceInstance-qualified observations; not normal Product SDK semantics and not Principal-authored business overrides.
3. **Provider/business-system protocol ingress** — OAuth callbacks/webhooks/notifications/provider protocol as already accepted by D5-B1/D4.

A separate source-ingestion surface is now backed by a real named consumer: Sankhya/external source synchronization and future external connectors. Therefore it is no longer speculative technical taxonomy.

D5 later decides exact HTTP paths, auth mechanism, schemas, batch/pagination and generation packaging.

The built-in Sankhya adapter need not consume the HTTP Source Ingestion surface merely because external connectors do; both consume the same semantic application contract.

---

# 16. Legacy/current API implications

The following current shapes are evidence, not target publication architecture:

- `/catalog/products*` — cannot become MPC Product Master;
- ERP/XLSX import routes — source-ingestion mechanism only if still needed;
- provider category/attribute routes — requirement evidence, not universal Product API taxonomy;
- `/product-links/*/generations` — candidate-generation mechanism, not user-owned Product master;
- generic `/mutations` — still rejected;
- direct provider payload creation/update routes — rejected as Product API ontology;
- manual database/scripts that fill publication/provider fields — replaced by owner API/Source Ingestion semantics.

Useful requirements from current code may be rehomed; route/module shape is not preserved.

---

# 17. Authority / reopen analysis

## D0

**No reopen currently required.** Creation/publication and operational source flexibility are already inside the accepted Product 1.0 mission.

## D1

**No reopen presumed. Boundary challenge required.**

Leading interpretation:

- Readiness owns mutable Channel Publication Preparation because it is requirement-resolution/readiness state for `Product + Installation`, not live Listing representation;
- Offering owns Listing Intent, actual listing representation/lifecycle and convergence.

Reopen D1 only if independent review/concrete provider proof shows one of:

- publication authoring itself has an independent lifecycle/authority not safely owned by Readiness or Offering;
- Readiness would become de facto Listing representation authority;
- Offering needs a new semantic dependency absent from the accepted edge set;
- Availability/Fulfillment must become part of initial publication in a way no accepted edge can supply.

## D2

**No reopen presumed.**

- external Product identity already exists;
- Listing Intent identity already exists;
- Channel Publication Preparation can initially be relationship-scoped without inventing a new canonical Product/Draft identity;
- Listing Intent can freeze historical preparation values/provenance.

Reopen D2 only if correctness proves a new durable canonical identity/lineage class is actually required.

## D3

**No reopen presumed.**

Source ingestion is acquisition/observation; Offering consuming current Readiness meaning fits existing Readiness→Offering Q direction.

Reopen D3 only if concrete publication requires a new cross-owner semantic dependency or atomicity assumption not already accepted.

## D4

**Targeted reopen required.** This candidate closes the missing source/publication/provider-requirement contract and must be independently reviewed before D5-B2 resumes.

## D5-B1

**No reopen presumed.** B1 already allows a separately justified technical surface once a concrete consumer exists. Source Ingestion now supplies that evidence. B2/later D5 must classify it explicitly.

---

# 18. Complexity law / YAGNI

## Essential complexity preserved

- source-qualified Product facts;
- source vs MPC-authored provenance;
- changing provider requirement schemas;
- category/product-type context;
- missing/conflicting requirement state;
- media origin/selection;
- explicit override/follow-source behavior;
- Listing Intent decision-time snapshot;
- provider-specific serialization/convergence;
- no hidden price/availability authority transfer;
- multiple physical source-ingress modes over one semantic use case.

## Accidental complexity rejected

- MPC Product Master/PIM;
- universal product attribute catalog;
- generic `map[string]any` provider bag;
- generic transformation/rules DSL;
- plugin/factory framework for hypothetical sources;
- built-in self-HTTP requirement;
- implicit last-write-wins;
- source values copied into manual editable fields without provenance;
- provider DTO stored as canonical draft;
- one giant universal Publication object containing price/availability/fulfillment semantics from different owners.

### Prepared seams, not full future capability

Admit now:

- Source Observation semantic port;
- external Source Ingestion as a real later D5 surface;
- per-publication source-follow/override resolution;
- provider requirement descriptors;
- Listing Intent snapshot.

Do **not** build now:

- global mapping/rule designer;
- auto-mapping ML→Amazon→Shopee attributes;
- reusable organization-wide attribute formulas;
- PIM catalog governance;
- every provider's schema engine;
- generalized connector marketplace.

If repeated per-product resolutions prove a reusable Source→Requirement mapping materially reduces total complexity, add the smallest mapping primitive later with evidence.

---

# 19. Proof strategy before implementation

## P1 — Today/future source transition

For one Product with a missing source field:

1. source observation lacks a required publication value;
2. automation Principal authors an explicit override;
3. readiness becomes sufficient and Listing Intent can snapshot the override;
4. later source begins supplying a valid value;
5. source truth updates without mutating the override;
6. override is deliberately cleared;
7. preparation now follows source;
8. no Product identity or historical Listing Intent is rewritten.

## P2 — Authority spoofing

A Product API client cannot mark its authored value as a Sankhya observation.

A Source Ingestion client cannot author an MPC business override merely by choosing `source=manual` or changing Organization/SourceInstance in payload.

## P3 — Partial source update

A partial source payload omitting `description` does not clear a previously known description unless source semantics explicitly establish absence/removal.

## P4 — Provider-schema change

Provider requirement version N allowed publication; version N+1 adds/changes a required attribute. Existing historical Listing Intent remains explainable under N; current preparation re-evaluates against N+1 and becomes missing/conflicting if necessary.

## P5 — Second-provider structural inversion

Map one conceptual Product into:

- Mercado Livre User Product/Item current model;
- Amazon Product Type Definition + Listings Item;
- Mirakl Product/Offer/PriceStock split.

The source Product identity and Readiness/Offering authorities must remain unchanged. Only provider requirement descriptors and adapter serialization differ.

## P6 — Media source transition

Use Sankhya source image for one publication and MPC-uploaded media for another without rewriting source Product media. Listing Intent preserves which media set it attempted.

## P7 — Embedded versus external connector

Run the same semantic Source Product Observation through:

- embedded Sankhya adapter → application port;
- external HTTP Source Ingestion adapter → same application port.

The committed semantic result is equivalent; no HTTP-only business rule exists.

## P8 — No adapter truth selection

Remove one required preparation resolution. Provider adapter must fail before dispatch/return owner-level insufficient readiness; it may not choose another source/default silently.

## P9 — Availability authority attack

Find the concrete selected Mercado Livre creation operation and prove either:

- publication representation can be established without inventing/owning Sellable Availability in Offering/Readiness; or
- surface the exact D1/D3 contradiction and reopen before D5 resumes.

## P10 — Agent path

Automation Principal can author allowed publication preparation through the Product API with normal access/audit/owner semantics; no special database/provider bypass is required.

---

# 20. Adversarial review package

Fable should independently reconstruct authority and attack at least:

1. Is Channel Publication Preparation genuinely Readiness-owned, or does it steal Offering representation authority?
2. Does per-requirement FOLLOW_SOURCE / EXPLICIT_OVERRIDE introduce a hidden PIM/rules system?
3. Is a separate Source Ingestion API really necessary, or can D4 provider/business-system ingress cover it without ambiguity?
4. Does the Source Ingestion API create a generic integration platform contrary to D0/D1?
5. Can a machine/source binding prevent caller authority spoofing without inventing another tenant/source registry?
6. Should a built-in adapter ever call the HTTP Source API for structural consistency, or is shared application port correctly smaller?
7. Does Listing Intent need to copy/freeze all resolved fields, or can immutable references/provenance be sufficient without stale-history risk?
8. Are provider Requirement Descriptors too generic/abstract compared with provider-specific contracts?
9. Can Amazon/Mercado Livre/Mirakl all fit without a universal PIM or provider bag?
10. Does current Mercado Livre User Product publication force Availability/Offering coupling that reopens D1/D3?
11. Are images/media a business concept requiring D2 identity, or can content/provenance references remain bounded?
12. Could an explicit override silently become stale/wrong after source/provider changes, and what owner detects that?
13. Does manual/agent authoring require separate authorization beyond ordinary access + Listing Intent Governance?
14. Is `DERIVED` already overengineering and should baseline keep only source/override?
15. Does Source→Requirement mapping need a reusable rule now because Sankhya custom e-commerce fields already make it a present consumer?
16. Is Google primary/supplemental feed precedent useful, or would copying precedence/rule mechanics be accidental complexity for MPC?
17. Are ANYMARKET/Channable/Linnworks showing a missing PIM requirement rather than merely mapping complexity?
18. Does editing an existing Listing differ materially enough from creation to require separate preparation/intent semantics?
19. Can provider shared-field blast radius cause one preparation change to affect multiple external Items/User Products and is that represented correctly?
20. Does a source observation need full-population/coverage semantics for publication, or only point Product facts?
21. Is Source Ingestion allowed to push arbitrary custom fields, and if not, how does a new source fact get admitted without code changes becoming excessive?
22. What is the hardest future change: Sankhya becomes fully authoritative; second ERP; second marketplace; provider schema churn; multi-origin stock; and does the seam survive each?
23. Does any proposed construct exist only because another abstraction exists?
24. What is the smallest proof that D4 can close and D5-B2 safely resume?
25. Is any D0/D1/D2/D3/D5-B1 reopen genuinely required by evidence rather than reviewer preference?

Reviewer severity does not create authority.

---

# 21. Reopen / stop triggers

Reopen only the implicated parent decision when material evidence proves:

1. Publication Preparation cannot be owned by Readiness without stealing Offering semantics → D1.
2. Initial publication requires Availability/Fulfillment-owned values through a dependency absent from D1 → D1/D3.
3. A new durable identity/lineage object is required beyond source Product + Listing Intent → D2.
4. Publication requires atomic cross-owner mutation for correctness → D1/D3.
5. Provider requirement/effect cannot be represented through consumer-owned resolution without raw DTO leakage → D4 redesign.
6. Source Ingestion cannot preserve Organization/SourceInstance/provenance safely → D2/D4/D5 boundary review.
7. A real external source/API consumer requires a compatibility/versioning obligation materially different from D5-B1 assumptions → D5.
8. Repeated mappings prove a reusable Source→Requirement mapping primitive is essential → targeted Readiness/D4 decision, not generic PIM.
9. A second real Product master changes source authority assumptions → D1/D2 targeted review.
10. Current Mercado Livre User Product selected lane makes the proposed separate publication/availability progression impossible → targeted D1/D3/D4 reopen before D5.

Framework preference, provider symmetry, current code convenience and hypothetical providers are not reopen evidence.

---

# 22. Candidate outcome

**Proposed outcome:** `RESTRUCTURE NOW` for the missing publication/input seam, with no presumed D0/D1/D2/D3 or D5-B1 reopen.

If independent review and operator ratification converge, D4 should canonically add:

- explicit Source Product Observation ingress semantics;
- one semantic source-observation application port shared by embedded pull adapters and external push ingress;
- no mandatory self-HTTP for embedded adapters;
- structurally separate source ingestion and MPC Product authoring authorities;
- automation Principal as a legitimate Product API authoring client, never source-authority impersonation;
- Readiness-owned Channel Publication Preparation as the leading boundary, subject to final challenge;
- explicit FOLLOW_SOURCE / EXPLICIT_OVERRIDE provenance semantics, with DERIVED only if review keeps it;
- provider Requirement Descriptors as source-qualified D4 evidence;
- media/source-image versus MPC-authored publication-media provenance;
- Listing Intent decision-time resolved publication snapshot under Offering;
- adapter as protocol translator, never truth selector;
- explicit authority fence for price/availability/fulfillment inputs;
- Source Ingestion API as a now-justified later D5 external technical surface;
- no MPC Product Master/PIM, universal field bag, generic mapping DSL or connector framework;
- D5-B2 resume only after this contract is accepted and any genuine parent-stage contradiction is closed.

**D5-B2 review candidate remains non-authoritative and frozen pending this decision. Implementation remains blocked until D9.**
