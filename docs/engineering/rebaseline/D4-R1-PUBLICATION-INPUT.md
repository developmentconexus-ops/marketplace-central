# D4-R1 — Publication Input & Listing Authoring Contract

> **Status:** ACCEPTED / CANONICAL — targeted D4 amendment after post-D4 discovery during D5  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authority:** accepted D0 → D4 + D5-B1  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Accepted:** 2026-08-18  
> **Operator ratification:** 2026-08-18

## 1. Role of this amendment

This file is the accepted targeted amendment to `D4-EXTERNAL-INTEGRATIONS.md` for marketplace publication input/listing authoring.

The original D4 B1/B2/B3/B4 decisions remain accepted and canonical. This amendment was required because D5-B2 discovery exposed one missing seam: D4 already defined provider requirements/protocol/effects, and D1/D2 already defined Product, Readiness, Offering and Listing Intent authority, but no accepted contract said how external Product facts plus MPC-authored channel values become a provider-valid, historically explainable listing action.

The targeted review confirmed:

- **D0→D5-B1 remains one coherent architectural system**;
- no D0/D1/D2/D3/D5-B1 reopen is required;
- the existing 12 D1 business boundaries remain justified;
- no 13th publication/content/integration business domain is introduced;
- D5-B2 may resume only from this accepted seam, not from the pre-R1 review candidate by inheritance.

`docs/README.md` remains the sole current-program status/router. The original D4 §10 next-action text is a historical closure snapshot from before this targeted reopen; current next action is defined only by the router.

---

## 2. Root cause and governing invariant

### Root cause

The previous architecture began publication one step too late.

It had:

- external Product identity and evidence;
- Product↔channel Readiness requirements/correspondence;
- Offering-owned Listing Intent and convergence;
- provider-specific publication requirements/protocol in D4;

but no accepted owner-preserving seam for:

- source facts that may be incomplete;
- channel values authored temporarily inside MPC;
- images/media from source or MPC authoring;
- changing category/product-type requirements;
- later transition from manual/agent authoring to source-following;
- provider contracts that physically combine listing content, price, availability or fulfillment in one request.

Without this seam an adapter or client could become a de facto truth selector, MPC could silently become a PIM, source/manual precedence could collapse into last-write-wins, or provider DTOs could leak into the core.

### Governing invariant

> **External Product truth, Readiness-owned requirement/source evidence, Offering-owned desired listing state, other owner-issued business meanings, provider protocol and provider actual Listing state remain distinct authorities. Every material publication input is resolved by its accepted owner before provider serialization; D4/D7 mechanisms may compose owner-issued meanings into the provider operation but never own, default or recompute them.**

---

## 3. Whole-product ownership model

### 3.1 External Product remains external

Product identity remains:

```text
Organization
+ SourceInstance
+ native Product key
```

MPC does **not** create a Product master, PIM, ProductDraft or synthetic Product identity merely to publish into marketplaces.

External Product facts remain source-qualified evidence with honest knowledge, provenance and time semantics.

### 3.2 Readiness owns requirements and source-level readiness

Product & Channel Readiness owns:

- Product↔channel correspondence;
- applicable publication/channel requirements;
- provider/source evidence needed to understand those requirements;
- source candidates that may satisfy a requirement;
- missing/conflicting/unsupported source-level state;
- source-level readiness/sufficiency under Readiness semantics.

Conceptual subject:

```text
Organization
+ Marketplace Installation
+ SourceInstance-qualified Product
```

Readiness may query D4 through a **Readiness-owned publication-requirements port**. D4 translates provider/category/product-type requirement evidence only to the extent this consumer needs it.

Readiness does **not** own desired listing copy, authored overrides, selected listing media, Listing lifecycle, Price Intent or Sellable Availability.

### 3.3 Offering owns the one create/edit authoring identity

Marketplace Offering Operations owns `ListingIntent` as the single MPC authoring/draft identity for listing creation and editing.

```text
Create
  → ListingIntent(target = no existing Listing)

Edit
  → ListingIntent(target = existing source-qualified Listing)
```

A mutable pre-freeze draft is part of the ListingIntent lifecycle; it does not create another canonical identity class.

The draft may carry only desired listing representation/action meaning owned by Offering, proportionately including:

- intended listing lifecycle action;
- target category/product-type selection when part of desired listing representation;
- channel-facing content values;
- listing-specific media selection/order/role;
- provider-qualified requirement resolutions needed by this desired representation;
- typed references/correlation to other owner-issued meanings when the provider execution requires them.

The draft does not become Product master, Sellable Availability authority, Economics authority or Fulfillment authority.

### 3.4 Offering owns draft dispatchability, not Readiness

Readiness never reads ListingIntent content.

Offering consumes current Readiness-owned requirement/correspondence/source evidence by the already accepted `Readiness → Offering` Q edge and decides **draft dispatchability**:

> whether this specific ListingIntent's resolved desired values satisfy the current applicable Readiness meaning well enough to freeze/dispatch.

Therefore:

- Readiness may still report a source requirement missing because the source itself lacks a value;
- Offering may nevertheless have a dispatchable draft when an allowed MPC-authored override supplies the desired listing value;
- provider requirement revision can make a draft non-dispatchable without mutating Readiness into an Offering-content authority;
- D2 pre-dispatch correspondence consistency is revalidated against current Readiness meaning before effect.

---

## 4. Publication value resolution

Baseline listing-value resolution has exactly two modes.

### `FOLLOW_SOURCE`

The ListingIntent points to one explicit source-qualified candidate admitted through Readiness meaning.

Rules:

- source identity/provenance remains explicit;
- material value is resolved from the source at material decision points, including freeze and dispatch-time revalidation;
- unavailable/unknown remains unavailable/unknown;
- no hidden fallback to a manual value;
- source changes do not rewrite already frozen historical intents.

No baseline event is required merely because the source changed. A later proven independent reaction may add an event only under D3 rules.

### `EXPLICIT_OVERRIDE`

The ListingIntent contains an MPC-authored listing value.

Rules:

- effective Principal/automation identity and time/provenance are preserved;
- override never mutates or falsifies source Product truth;
- later source change does not silently replace the override;
- removing/replacing the override is an explicit Offering action;
- historical ListingIntent basis remains explainable after later source/provider change.

### No baseline `DERIVED`

No generic derived/transformation mode, feed-rule engine, expression DSL, mapping engine or formula language is admitted now.

If repeated real cases later prove a deterministic owner-specific transformation materially reduces total complexity, adjudicate the smallest explicit transformation then. Value-preserving provider wire formatting remains adapter serialization and is not business `DERIVED` meaning.

---

## 5. Human and automation authoring

Humans and automation use the same Semantic Product API authority.

An automation such as Claude/Codex/Fable-style tooling is a D2 **automation Principal**, not an AI-specific business authority or backdoor.

It remains subject to:

- Organization scope;
- ordinary Permission;
- Offering business validity/disposition;
- Governance where required;
- consequential idempotency/effect-safety;
- audit/history;
- provider convergence rules.

D2 §10.3 remains binding:

> automation recurrence never grants authority to silently reopen or reverse a standing human decision in the same semantic scope.

Therefore an automated run cannot silently replace a human-authored `EXPLICIT_OVERRIDE` or resolution choice. Supersession must be explicit domain semantics with preserved attribution and any required disposition/Governance.

---

## 6. Source acquisition and future integration ingress

Source acquisition is D4 mechanism/evidence, not a Product domain and not a generic `SourceProductObservation` business owner.

One source acquisition may feed multiple consumer-owned semantic ports. No component owns the source Product payload as a whole.

### Embedded adapter

A built-in Sankhya adapter may call consumer-owned application/semantic ports directly:

```text
Sankhya sanctioned API
    ↓
D4 Sankhya adapter
    ↓
consumer-owned semantic ports
```

No mandatory self-HTTP loopback is required merely for structural symmetry.

### Future external connector

A physically external connector may enter through a bounded Integration Ingress HTTP adapter **when a concrete connector exists**.

```text
external connector
    ↓ authenticated bounded ingress
D4 HTTP inbound adapter
    ↓ translation/admission mechanics
consumer-owned semantic ports
```

This seam is prepared but no generic Source Ingestion API, plugin registry, `/entities`, `/resources` or connector platform is frozen now.

When the first real external connector exists, D5 defines the smallest wire contract for that source class while preserving Organization + SourceInstance binding, provenance, partial/absent semantics, replay safety and inability to impersonate MPC-authored ListingIntent values.

---

## 7. Provider publication requirements

Provider publication requirements remain provider-authoritative D4 evidence translated through consumer-owned ports.

### Human-operable publication decision seam

```text
canonical decision identity
  requirement_key / option_key / unit_key / source_candidate_key / correspondence candidate_key

current human read projection
  source/provider display presentation needed to recognize the choice

write/effect
  canonical key only + current owner revalidation
```

When a human choice is admitted, Readiness preserves the smallest current human presentation needed to recognize the SourceProduct/SourceInstance, requirements, options, units, `FOLLOW_SOURCE` candidates and correspondence candidates. This read evidence remains adjacent to the canonical key and never becomes decision identity or write authority.

Provider expressions, raw paths, adapter DTO topology, arbitrary maps and generic metadata bags remain rejected at the Product seam. Unknown/unavailable presentation blocks a choice that cannot be performed honestly rather than fabricating a label.

For a Readiness consumer the semantic contract may preserve only what is material, such as:

- provider/context-qualified requirement identity;
- required/recommended/optional/conditional applicability;
- data kind/cardinality/options/constraints;
- editability/immutability where material;
- provider schema/version/checksum/revision evidence;
- acquisition/provenance/time.

This is not a universal `ProductAttribute` catalog or generic ProviderRequirement business framework.

A provider-specific concept is promoted to a shared MPC semantic concept only when its meaning is genuinely provider-independent and an accepted owner needs that meaning independently of protocol.

---

## 8. Media / images

Media participates in listing authoring without creating Product-media authority.

Legitimate origins include:

- source-qualified Product media, such as Sankhya Product images;
- MPC-authored/uploaded media for the listing context.

Rules:

- source media remains external evidence;
- MPC-uploaded media never rewrites source Product media;
- ListingIntent owns only listing-specific selection/order/role;
- provider image IDs/CDN/resource topology remains provider-local;
- frozen intent history preserves enough media/provenance to explain what was attempted;
- exact blob storage, hashing, resizing, caching/CDN and upload realization belongs later;
- arbitrary external URL strings are not automatically trusted publication media.

No reusable ProductAsset/AssetFamily/media master is admitted. A later reusable cross-listing media master is explicit Product/PIM pressure and must be re-adjudicated.

---

## 9. R1-G1 — Mercado Livre initial publication × Availability — CLOSED / PASS-B

Current official Mercado Livre User Product evidence establishes that the selected publication direction may physically require `available_quantity` in the same new-item provider request that creates the listing/User Product representation.

D4-R1 therefore does **not** assume a normal active representation-first path without Availability input. The paused/zero-quantity path remains unclaimed for the selected Installation until real proof establishes it.

This physical provider contract does **not** create `Availability → Offering` and does not reopen D1/D3.

### Joint technical realization law

When provider protocol requires multiple owner-owned meanings in one request:

```text
Offering-owned ListingIntent
        │ owner-issued listing meaning
        │
Availability
        │ owner-issued availability meaning / Availability Intent
        │
        ▼
D4/D7 external-effect execution mechanism
  validate + correlate owner-issued inputs
  serialize one provider request when required
        │
        ▼
marketplace
        │ authoritative reread
        ├─ Offering evaluates Listing convergence
        └─ Availability evaluates Availability convergence
```

Binding rules:

1. Offering never authors, stores or owns the Availability quantity as ListingIntent content.
2. ListingIntent may keep only a typed correlation/reference to the separate Availability-owned input/intent when historical explanation requires it.
3. Availability consumes the intended marketplace target through the already accepted `Offering → Availability` meaning and owns Sellable Availability/Availability Intent.
4. Availability never authors desired Listing representation.
5. The execution mechanism may require both owner-issued inputs/proofs before dispatch; mechanism does not own either answer.
6. Adapter cannot default, calculate or infer either business value by protocol convenience.
7. Active creation without the required Availability-issued value fails closed before provider dispatch.
8. No cross-owner atomicity is claimed. One provider request may yield partial/asynchronous convergence and each owner evaluates its own result.
9. The first controlled real creation remains a D8 proof and must include authoritative reread and shared-User-Product blast-radius verification.

A targeted D1/D3 reopen is required only if real execution proves one owner cannot issue its meaning for the intended target without a genuinely new semantic dependency, or if correctness requires one owner to persist/own another owner's business meaning.

The same mechanism/authority distinction applies if a provider create/edit operation physically requires a Fulfillment-owned input: technical composition is legal; hidden business-edge creation is not.

---

## 10. Multi-step / partial publication convergence

Listing creation/edit is not assumed to be one atomic provider request or one boolean outcome.

Current Mercado Livre evidence includes multi-step and asynchronous behavior such as:

- item creation plus description operations;
- User Product shared-field propagation;
- family-level asynchronous tasks/member outcomes;
- shared-field blast radius across related Items;
- provider-specific per-UP and edit restrictions.

Therefore D4 publication effects preserve where material:

- accepted / rejected / pending / ambiguous;
- step/member/aspect-level confirmed / rejected / pending / ambiguous / not-executed;
- authoritative reread per material resource/aspect;
- no overall convergence merely because an early provider request returned `2xx/201/202`;
- no blind retry of already-confirmed or ambiguous steps;
- intended/authorized/attempted blast radius when provider shared resources widen the physical effect.

Offering concludes its listing-representation convergence. Availability and other owners conclude their own semantics independently.

D5-B2 must not freeze a `createListing = success` contract that collapses these distinctions.

---

## 11. Historical ListingIntent snapshot

Before consequential provider dispatch, the material ListingIntent becomes immutable enough to explain the attempted action.

Preserve proportionately:

- intent identity/action/target scope;
- source-qualified Product + Installation + target Listing where applicable;
- exact desired Offering-owned listing values;
- `FOLLOW_SOURCE` references or `EXPLICIT_OVERRIDE` author provenance;
- provider requirement/schema revision materially used;
- listing-specific media selection/provenance;
- decision-time Readiness/disposition/authorization references;
- typed correlation to other owner-issued inputs participating in joint realization;
- provider attempt/result/convergence evidence.

This is historical decision/action context, never current source/provider authority. Snapshot only what materially supports explanation/correctness; do not create a payload archive/PIM copy by history.

---

## 12. Complexity / Global Maximum

Accepted Global Maximum:

> **Readiness requirements/source evidence + Offering-owned ListingIntent as the single create/edit authoring identity + owner-preserving D4/D7 provider realization.**

Explicitly rejected:

- Sankhya-only publication with no MPC authoring;
- MPC Product/PIM master;
- provider JSON/field bag as canonical state;
- Readiness-owned `PublicationPreparation` aggregate;
- generic `SourceProductObservation` business service;
- generic transformation/rule/mapping engine;
- mandatory self-HTTP for embedded adapters;
- generic connector/plugin platform;
- AI-specific authoring API/backdoor;
- separate creation and editing architectures;
- giant Publication aggregate absorbing Price/Availability/Fulfillment.

The seam survives:

- future Sankhya-complete source data without identity migration;
- a second business system through SourceInstance + adapter;
- a second marketplace through provider-specific requirements/serialization;
- provider schema churn through current requirement evidence + historical intent context;
- heavier agent automation through the existing Principal/Product API/Governance model.

No new business boundary or infrastructure framework is required.

---

## 13. Proof obligations

Before/through later implementation, the architecture must be falsifiable by at least:

1. **Source-to-override-to-source transition:** missing source field → explicit ListingIntent override → frozen historical intent → later source value appears → new intent deliberately follows source without rewriting history.
2. **Authority spoofing:** Product API cannot claim an override came from Sankhya; integration ingress cannot author an Offering override by payload convention.
3. **Partial source observation:** omitted field in partial acquisition does not erase previously known source fact unless source semantics prove removal.
4. **Provider schema change:** old historical intent remains explainable under old requirement revision; current draft dispatchability re-evaluates under new requirement meaning.
5. **Second provider:** Mercado Livre/Amazon/Mirakl-class protocol differences do not change Product identity or D1 ownership.
6. **Media provenance:** source media and MPC-authored listing media coexist without creating Product-media master.
7. **Embedded vs external connector:** both reach the same consumer semantics without HTTP-only business rules.
8. **No adapter truth selection:** missing owner-resolved input blocks before provider dispatch; adapter never invents/defaults.
9. **R1-G1 joint realization:** active listing dispatch without required Availability-issued input fails closed; quantity cannot be sourced from ListingIntent content; each owner evaluates its own convergence.
10. **Human/automation override:** recurring automation cannot silently replace a standing human override; explicit supersession preserves lineage.
11. **Multi-step publication:** an early provider 2xx cannot falsely close later description/shared-UP/asynchronous steps.

D8 owns the first controlled real Mercado Livre creation/write proof. Implementation remains blocked until D9.

---

## 14. Reopen triggers

Reopen only the implicated decision when material evidence proves:

1. ListingIntent cannot legitimately carry editable desired listing representation before consequential submission → D1/D2 targeted review.
2. Readiness cannot expose requirement/source-level meaning without owning authored listing state → D1.
3. a selected provider requires a cross-owner semantic dependency not representable as joint technical realization → D1/D3.
4. publication correctness requires atomic cross-owner mutation → D1/D3.
5. a new durable identity/lineage class is required beyond external Product/Listing + ListingIntent → D2.
6. provider requirement/effect cannot fit consumer-owned semantic ports without raw DTO leakage → D4.
7. an external source connector cannot preserve Organization/SourceInstance/provenance safely → D2/D4/D5 boundary review.
8. repeated real source→listing transforms prove a small reusable transformation is essential → targeted owner decision, never generic PIM by default.
9. a real external connector creates a concrete wire/security/compatibility requirement → D5 defines the smallest ingress contract.
10. reusable cross-listing media/content authority becomes a real product requirement → explicit D0/D1 Product/PIM review.

Framework preference, provider symmetry, current implementation convenience and hypothetical future systems are not reopen evidence.

---

## 15. Final disposition / handoff to D5

**Outcome:** `CURRENT D0→D5-B1 STRUCTURE CONFIRMED` + `RESTRUCTURE NOW` for the missing D4 publication-input seam.

Independent Fable Global Coherence review returned `REVISE`; GPT adjudicated F-R1-1..F-R1-5 as local corrections, all incorporated before operator ratification. No second review round was required because no contradiction survived adjudication.

**R1-G1 = CLOSED / PASS-B.**

No D0/D1/D2/D3/D5-B1 reopen is required.

D4 remains **CLOSED / ACCEPTED AS A WHOLE**, now including this accepted targeted amendment. The next architecture action is **D5-B2 — Product Operation / Resource Surface**, derived from D0–D4 including this R1 amendment plus D5-B1.

The existing pre-R1 D5-B2 review candidate remains non-authoritative evidence only and must be revised/re-derived against current authority before any B2 ratification.

---

## 16. Listing variations amendment — 2026-08-26

A ratified D6-R2/B23 upstream finding (Mercado Livre variations with no MPC model) reopened this owner boundedly. The operator adjudicated the [variations Global Maximum design](../../superpowers/specs/2026-08-26-listing-variations-global-maximum-design.md); it is now authority:

- provider variation vocabulary (axes/options per publication context) is census evidence served through the existing Readiness publication-requirements port: `PublicationRequirements.variation_axes` (`VariationAxisSpec`, option or text kind) and `PublicationRequirement.scope: listing | per_variation`;
- Offering's `ListingIntentDesired` gains optional `variations` (chosen axes + authored options); each option carries coordinate identity (`VariationCoordinate`: axis/option canonical keys or axis text value), an optional SKU-level `source_product` ref, and its own per-variation `requirement_resolutions`/`media_selection` using the unchanged resolution kinds;
- per-option quantity and price remain owned by Sellable Availability and PriceIntent; D4 composes them at dispatch under the §9 R1-G1 composition law — they never enter ListingIntent authoring;
- the Offering `MarketplaceListing` observation gains optional `observed_variations` (coordinates + typed presentation + honest observed fields/price);
- no new operation, path, Permission or Principal kind; the shapes are protected by `verify-human-operable-read-projection.mjs` (16 negative controls, including: no price/quantity/label in variation authoring, coordinate-typed identity, honest observed presentation).

Kits/bundles remain an explicitly separate future finding. Variation-aware Performance/Economics splits remain DEFERRED with a recorded reopen trigger.
