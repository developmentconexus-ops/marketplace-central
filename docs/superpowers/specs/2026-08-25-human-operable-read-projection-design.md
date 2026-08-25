# Human-Operable Read Projection & Wire Conformance — Design

> **Status:** DESIGN CANDIDATE — architecture/authority direction operator-approved on 2026-08-25; written-spec review still required before contract implementation
> **Trigger:** D6-R2 B20 planning after B10 integration
> **Method:** DevelopmentConexus Engineering Method v1.0.0 + Frontend Product Experience Planning Method v2.3
> **Baseline:** `main@bdbbef43ed3a5e9d912e67ddac5173024352eaa3`
> **Implementation:** BLOCKED until this written design is operator-reviewed, an implementation plan is approved, and the repository roadmap authorizes execution

## 1. Decision

**Outcome: `RESTRUCTURE NOW` — targeted Readiness/Offering human-operable read projection + W2/OAD conformance repair.**

The accepted Product architecture remains intact. The repair does **not** create a presentation business domain, Product/PIM master, generic entity registry, generic key→label service, provider DTO passthrough, new Product operation, new ordinary Permission or new Principal kind.

The structural defect is narrower and deeper:

> Canonical references/keys that are correct for identity, correlation and writes have been reused as if they were sufficient human read projections. At the same time, some human-presentation distinctions already accepted in W2 or already implemented in another owner were not propagated coherently to the canonical Product OAD.

Target invariant:

> **Canonical identity and decision carriers remain key/ID based. A Product read whose admitted human job requires recognizing, comparing, selecting or explaining a subject must expose the smallest owner-correct human presentation needed for that job, with honest knowledge/freshness semantics. Presentation metadata never authenticates, authorizes, scopes, correlates, matches, converges, or substitutes for canonical identity. Historical/purpose presentation snapshots remain distinct from current presentation.**

## 2. Evidence and root cause

### 2.1 Current wire inconsistency

The canonical OAD currently demonstrates both patterns:

- owner-local human presentation exists for `MarketplaceInstallation`, `SellingEntity`, `SourceProductSearchHit`, `InventorySource` and `FulfillmentNode`;
- `MarketplaceListing`, `ListingIntentListItem`, `PriceIntentListItem`, Availability listing targets, Market/Economics listing subjects and other reads often expose only canonical IDs/keys.

Performance independently added a non-authoritative Listing `display_name`, proving the same human-recognition need was already solved locally in another owner. Governance independently uses `subject_display_label` as an immutable authorization-purpose presentation snapshot, proving that current presentation and historical/purpose presentation are materially different meanings.

### 2.2 Provider evidence

Current Mercado Livre documentation preserves technical identity separately from human presentation:

- category/domain IDs and names;
- attribute IDs and names;
- allowed-value IDs and names;
- Item IDs and titles;
- number/unit attributes with provider-defined allowed units.

D4 therefore has legitimate source evidence from which bounded human presentation can be translated without turning provider vocabulary into MPC identity or authority.

### 2.3 B10 falsifiers

B10 exposed two distinct defects:

1. `PublicationRequirement`, option keys, unit keys and publication-context keys are not human-operable without provider/source presentation metadata;
2. `ResolveProductChannelCorrespondence` requires `candidate_key`, while current correspondence reads do not always return a selectable, human-recognizable candidate population. `unresolved` returns no candidates and `conflicting` exposes only opaque keys.

The second defect is functional, not cosmetic: a human may be asked to resolve a correspondence without a contract-supplied choice set sufficient to perform the admitted capability.

### 2.4 W2 → OAD drift

Canonical W2 already requires more than the current OAD for directly implicated Offering reads, including proportionately:

- MarketplaceListing publication context;
- bounded observed representation/media;
- current observed price;
- observation/freshness/provenance;
- ListingIntent resolved requirement/provenance read axes;
- stable authored-media descriptor vs volatile presentation descriptor;
- authenticated technical authored-media delivery rather than a durable anonymous locator.

The prerequisite must distinguish:

```text
semantic insufficiency
  → a parent/W2 rule was not yet complete enough for the proven human job

projection drift
  → W2 already decided the meaning, but the canonical OAD did not realize it
```

Do not use one label patch to hide either class.

## 3. Authority disposition

| Authority | Disposition | Reason |
| --- | --- | --- |
| D0 Product/System | CONFIRMED | Product purpose/actors unchanged |
| D1 domains/boundaries | CONFIRMED | no new business owner; existing owners still own meaning |
| D2 identity/data ownership | CONFIRMED | existing laws already separate canonical identity, projection and historical snapshot authority |
| D3 communication/events | CONFIRMED unless later falsified | no new semantic dependency/event is required by presentation itself |
| D4 External Integrations | BOUNDED REOPEN | external requirement/listing/source evidence must preserve materially required presentation and selectable candidate evidence |
| D4-R1 Publication Input | BOUNDED REOPEN | Readiness publication vocabulary/candidate evidence was semantically incomplete for human authoring |
| D5 W1 | CONFIRMED | paths, resource identity and custom capability grammar remain sound |
| D5 W2 | REOPEN | read/write distinction, dynamic vocabulary, correspondence candidate population and direct Offering drift require schema repair |
| D5 W3 | BOUNDED REOPEN | affected collection items must satisfy scan/select/navigate without N+1 detail fan-out |
| D5 W4 | CONFIRMED | operations/Permissions/Principal kinds unchanged |
| Product OAD | REPAIR | realize the repaired W2/W3 meanings and already-accepted W2 semantics |
| D6 frontend | BOUNDED REBASELINE | B10 correspondence region + P9; B20 stays paused; unrelated locks remain protected |

Expected Product inventory after this prerequisite remains:

```text
106 Product operations
31 ordinary Permissions
Principal kinds H / A / S
runtime baseline NONE
```

Any need to change those counts is a new material finding and must stop/reopen rather than enter this repair by convenience.

## 4. Rejected alternatives

### 4.1 Per-screen `display_name` patches — REJECTED / Local Maximum

Adding a label whenever a new screen complains leaves the same defect reachable in B23, B24, Availability, Market, Economics and future Listing consumers. Performance has already demonstrated this fragmentation.

### 4.2 Generic PresentationService / EntityDescriptor / key→label registry — REJECTED / overengineering

No independent presentation business lifecycle exists. A generic service would create a new authority or dependency graph merely because multiple owners have human reads.

### 4.3 Frontend dictionaries or direct provider lookups — REJECTED

Frontend hardcodes/provider calls would become parallel Product truth, bypass Product access/knowledge semantics, and couple UX to provider DTOs.

### 4.4 N+1 detail fan-out — REJECTED as production baseline

A collection required for human scanning must expose its necessary item-level presentation directly when the owner can supply it. Point GET fan-out does not repair a deficient collection contract.

## 5. Shared grammar: Ref ≠ Current Read Projection ≠ Snapshot

This is a reasoning/wire grammar, **not** a universal Product object hierarchy.

### 5.1 Canonical Ref

Canonical refs stay minimal and stable.

Example:

```text
MarketplaceListingRef
  marketplace_installation_id
  native_listing_key
```

They remain appropriate for:

- request targets;
- intent targets;
- correlation;
- identity qualification;
- authorization/business checks that already own the referenced meaning.

Presentation fields MUST NOT enter canonical refs merely for reuse convenience.

### 5.2 Current read presentation

When a current human read requires recognition, the read schema adds an **owner-specific typed presentation projection** adjacent to the canonical ref. Current presentation is source/owner-attributed read evidence, mutable and non-unique.

A dynamic external presentation that can be unavailable independently must preserve honest knowledge rather than silently falling back to a fabricated label.

For Marketplace Listing, target shape is conceptually:

```text
MarketplaceListingPresentation
  known(display_name)
  unknown
  unavailable
```

The UI may fall back to the canonical native key as a visibly technical identifier when presentation is unknown/unavailable, but the key never becomes a fabricated name.

### 5.3 Purpose/historical presentation snapshot

A label preserved specifically to explain an authorization/decision/history occurrence is a snapshot, not current source truth.

Governance's accepted `subject_display_label` is the model: immutable for its purpose, non-authoritative for identity and not silently refreshed from current source state.

No shared `DisplaySnapshot` business resource is created.

## 6. Publication vocabulary — read/write split

Dynamic provider vocabulary is where the current key-only model is insufficient. The target preserves keys for decision authority while returning names for human reads.

### 6.1 Publication context

Request/intent form remains key-based:

```text
PublicationContextRef
  category_key?
  product_type_key?
```

Human read form is distinct:

```text
PublicationContextView
  category?
    category_key
    display_name
  product_type?
    product_type_key
    display_name
```

`GetPublicationRequirements` query parameters remain keys. Its response returns the human read form. ListingIntent stores/serializes the key-based context when the selected context is part of desired Offering meaning, as canonical W2 already permits/requires proportionately.

### 6.2 Requirement identity

`PublicationRequirement` remains keyed by opaque `requirement_key` but gains non-empty current `display_name` in the human read schema.

The display name is provider/source presentation evidence translated by Readiness/D4. It is not a stable API identifier and cannot be submitted instead of `requirement_key`.

### 6.3 Options

Write value stays canonical:

```text
PublicationOptionValue
  kind = option
  option_key
```

Allowed-option read specification uses typed descriptors:

```text
PublicationOptionDescriptor
  option_key
  display_name
```

`PublicationOptionRequirementSpec` and option-list specs expose descriptor arrays rather than key-only arrays. Default/selected/correlated decisions continue to use keys.

### 6.4 Units

Write value stays canonical:

```text
PublicationNumberUnitValue
  value
  unit_key
```

Allowed-unit read specification uses:

```text
PublicationUnitDescriptor
  unit_key
  display_name
```

`default_unit_key` remains a key and must reference one current allowed descriptor. No generic unit-conversion/UoM engine is admitted.

### 6.5 PublicationValue vs PublicationValueView

Client-authored `PublicationValue` remains label-free.

Human-facing source evidence and observed Listing values use a distinct read union, conceptually `PublicationValueView`:

- text/exact-decimal/boolean/text-list preserve their canonical value;
- option includes `option_key + display_name`;
- option-list includes ordered keyed/name descriptors;
- number-unit includes exact value + `unit_key + unit display_name`;
- not-applicable stays explicit.

This prevents a client from authoring provider labels while allowing the Product to render values without hidden dictionaries.

## 7. Product-channel correspondence — operable candidate projection

Correspondence state and candidate population are distinct axes.

### 7.1 Current correspondence

Retain the owner-specific current meaning:

```text
resolved(candidate_key)
unresolved
conflicting(candidate_keys)
unknown
unavailable
```

The correspondence-scoped ETag remains unchanged and distinct from requirements/source-evidence revisions.

### 7.2 Candidate population

Add a current Readiness-owned candidate-population axis sufficient for human resolution:

```text
CorrespondenceCandidatePopulation
  known(candidates[])
  unknown
  unavailable
```

A known empty array is distinct from unknown/unavailable.

Each candidate is:

```text
ProductChannelCorrespondenceCandidate
  candidate_key
  display_label
  display_context[]?   # human-only bounded disambiguation strings
```

Rules:

- `candidate_key` is the only decision carrier;
- `display_label` / `display_context` never become matching or identity authority;
- D4/provider raw fields do not leak as an arbitrary map;
- Readiness may compose bounded human disambiguation text from legitimate provider/source evidence;
- equal presentation never collapses candidates with different keys;
- candidate population may contain more candidates than the currently conflicting subset; `conflicting.candidate_keys`, when present, must reference keys in the current known population;
- unknown/unavailable candidate population blocks human resolution rather than fabricating a choice.

### 7.3 Resolve capability

`ResolveProductChannelCorrespondence` keeps the same Product operation, Permission, subject and typed correspondence ETag. The request continues to submit only `candidate_key`.

At effect time the server revalidates:

- current Organization/subject;
- current correspondence revision;
- candidate still belongs to the current admissible candidate population;
- current business/automation safety rules, including no silent automation override of a standing human decision.

No candidate-list Product operation is added; candidate population belongs in the current Readiness Q because it is part of the already-admitted decision meaning.

## 8. MarketplaceListing read projection

### 8.1 Canonical identity

`MarketplaceListingRef` remains unchanged.

### 8.2 Collection item

`MarketplaceListingListItem` must carry, proportionately:

```text
listing: MarketplaceListingRef
presentation: MarketplaceListingPresentation
lifecycle
observed_at
```

This lets R20 scan/select/navigate without per-row point GETs.

`ListMarketplaceListings` coverage/cursor semantics remain unchanged.

### 8.3 Point actual-state read

Repair `MarketplaceListing` to satisfy canonical W2 proportionately:

```text
listing
presentation
lifecycle
publication_context: PublicationContextView
observed_fields[]
observed_media[] when materially available
observed_price? when materially available
observed_at
source observation/provenance sufficient to distinguish current/preserved evidence where the owner can serve both
```

`ListingObservedField` carries `requirement_key`, human requirement presentation, explicit knowledge state and `PublicationValueView` when known.

The actual-state read never owns ListingIntent convergence or PriceIntent convergence.

### 8.4 Cross-owner Listing presentation reuse

Performance currently has a local `display_name` solution for Listing. Normalize it to the same **MarketplaceListing presentation meaning** rather than preserving a parallel spelling.

Availability, Market and Economics reads that enumerate existing Listings may reuse this typed presentation projection only where their admitted P5 human job needs it. Reuse does not transfer Offering business authority; it is source-attributed presentation of the same external Listing identity.

Do not automatically attach presentation to every nested `MarketplaceListingRef` in every schema. Canonical refs inside machine/historical correlations remain minimal unless the specific human read requires adjacent presentation.

## 9. ListingIntent read conformance

Canonical W2 already requires ListingIntent reads to preserve more axes than the current OAD. Directly implicated conformance repair must restore, proportionately:

- key-based publication context where selected;
- current resolved requirement values/provenance sufficient to render and explain current draft state;
- authored-media stable descriptors;
- authored-media presentation descriptors for currently authorized reads;
- desired media selection/order;
- lifecycle, dispatchability/blockers, current external-effect/convergence evidence;
- actor/time attribution and historical attempt basis already required by W2 where the current OAD omitted them.

Writes remain sparse/declarative and do not gain response-only labels, presentation locators or server attribution.

No giant workflow/status aggregate is introduced.

## 10. Media presentation — preserve trust boundaries

### 10.1 Authored media

Retain W2's already-accepted distinction:

```text
ListingIntentMediaDescriptor
  stable identity + bounded content/provenance facts

ListingIntentMediaPresentationDescriptor
  stable descriptor + volatile authorized presentation reference
```

`CreateListingIntentMedia` returns the stable descriptor and parent validator, not a durable access locator. `GetListingIntent` may return presentation descriptors under current authorized read semantics.

Authored byte delivery remains the already-accepted authenticated technical presentation surface, outside Product operation count/SDK business operations.

### 10.2 Source media

Source media remains external Readiness/D4 evidence and must not reuse authored-media access semantics merely because both display images.

When B23 requires preview, `SourceMediaCandidate` may carry a distinct response-only `SourceMediaPresentationDescriptor` with volatile access reference under the appropriate source/Organization authorization model.

No generic Asset/Media owner, Product media library, arbitrary client URL or freely forwardable durable CDN locator is introduced.

### 10.3 Observed MarketplaceListing media

Actual provider Listing media may expose a distinct Offering read-only presentation descriptor sufficient to render the current observed representation. Provider image IDs/CDN topology remain D4-local unless an exact source key is independently material to Product semantics.

## 11. Collection/read consumer audit

The invariant is applied by proven human job, not by mechanically adding labels to every ID.

### Mandatory in this prerequisite

| Consumer | Why presentation is proven now |
| --- | --- |
| B10 / R10 | requirement/context/options/units/correspondence choices are human authoring inputs |
| B20 / R20–R21 | Listing collection/detail must support scan/select/understand |
| B23 / R22–R23 | ListingIntent collection/editor must identify source/target and render dynamic publication vocabulary/media |
| B24 / R24 | PriceIntent target must be recognizable to the operator |
| B30 / R30 | Availability target must be recognizable without N+1 Offering reads |
| B40 / R41 | Performance Listing already exposes a local label; normalize duplicated meaning |
| B50 / R50/R60 | Market/Economics Listing subjects must be recognizable in strategic analysis |

The wire repair may introduce adjacent presentation fields/projections to these read schemas without changing their canonical target refs.

### Not automatically modified

Sales, Shipment, Work, Governance history, Fulfillment, Post-Sale and other opaque IDs are **not** bulk-enriched by this prerequisite merely because they contain IDs. Their later P8/P9 work applies the invariant and must prove a real recognition/selection defect before changing wire shape.

Governance actionable requests are explicitly not a gap because they already carry purpose-specific `subject_display_label` snapshots.

## 12. B10 bounded rebaseline

The integrated B10 Global Maximum remains accepted:

```text
marketplace requirements
+ available source values
+ downstream ListingIntent authoring
+ provider validation where applicable
```

Preserve:

- search-first structure;
- exact Organization/Installation/source subject;
- requirement/value table;
- missing/conflicting/unknown/unavailable honesty;
- ListingIntent handoff;
- operator wording/layout already accepted.

Reopen only:

```text
correspondence region
  → real candidate presentation/selection when resolvable
  → known-empty / unknown / unavailable candidate population
  → effect + authoritative reread

P9
  → rerun against repaired canonical OAD
```

Because the correspondence interaction itself was not contract-supplied, the affected P8 region requires a fresh operator LOCK after executable candidate repair. Unaffected B10 structure remains protected under Frontend Method §5.3 bounded rebaseline.

## 13. B20 sequencing

PR #69 / B20 remains **PAUSED / NO P8** until this prerequisite is accepted and integrated.

After prerequisite integration:

1. repair/revalidate the bounded B10 correspondence region;
2. rerun B10 P9;
3. obtain operator re-LOCK for the affected region only;
4. resume B20 R20/R21 structural design from the repaired Listing read contract.

Do not render B20 against the currently deficient Listing projection.

## 14. Proof strategy

Before contract implementation, the implementation plan must define proof capable of falsifying at least:

1. **Identity fence:** equal labels never collapse distinct IDs/keys; requests still use canonical keys.
2. **No write-label authority:** generated write/request schemas cannot author `display_name`, option names, unit names, presentation locators or source attribution.
3. **Requirement operability:** every dynamic requirement/option/unit displayed by B10/B23 has contract-supplied human presentation.
4. **Correspondence operability:** a known candidate population can drive Resolve by key; unknown/unavailable cannot be presented as selectable; known empty remains distinct.
5. **Stale candidate safety:** candidate chosen from a stale correspondence revision cannot silently resolve after current evidence changes.
6. **Listing collection operability:** R20 can identify/select a Listing without N+1 point GETs and without fabricating a label.
7. **Presentation unavailability:** missing presentation does not erase a known Listing or become known empty.
8. **Performance de-duplication:** no parallel second Listing-label meaning survives between Offering and Performance.
9. **W2 conformance:** repaired OAD covers the directly implicated W2 MarketplaceListing and ListingIntent read axes; no already-accepted axis is silently dropped.
10. **Media trust separation:** source, authored and observed Listing media presentation references remain distinct trust types; no arbitrary URL becomes write authority.
11. **OAD tooling:** source lint/bundle/generation/type compilation remain green; operation/Permission/Principal counts remain 106/31/H-A-S unless a new explicit finding stops the work.
12. **Frontend bounded rebaseline:** unchanged B10 regions remain unchanged; only the correspondence interaction is reopened before re-LOCK.

CI remains proportional: prefer semantic/schema/type negative proof and existing aggregate gate; do not add prose-string ratification tests.

## 15. Adversarial challenge

### Objection: “Why not just use native_listing_key as the label?”

Because the Product already proves a human-recognition need through Performance and provider APIs expose separate ID/title semantics. A technical key may be shown as fallback/correlation evidence but cannot be promoted to a human name by convention.

### Objection: “This is a frontend convenience endpoint disguised as architecture.”

No endpoint is added. Existing reads are incomplete for already-admitted human jobs. W3 explicitly allows collection items to carry owner-semantic subsets needed to scan/select/navigate.

### Objection: “One shared presentation schema will become a platform ontology.”

The design forbids a universal Entity/Presentation envelope. Reuse is limited to the same semantic external subject (`MarketplaceListing`) or same publication vocabulary meaning. Other owners require their own proven need.

### Objection: “Provider labels change, so they cannot be Product data.”

Mutability is exactly why labels are presentation, not identity. Current reads may change; historical/purpose snapshots remain explicitly separate.

### Objection: “Second provider portability requires generic arbitrary metadata now.”

Rejected. The seam is typed key + presentation for meanings already provider-independent at the Product boundary. A future provider with a value family not representable without information loss reopens the smallest affected family.

### Objection: “Fix every opaque ID now.”

Rejected. Opaque identity is correct. Only a proven human recognition/selection/explanation job justifies an adjacent presentation projection.

## 16. Reopen triggers

Reopen this design only when material evidence shows one of:

- a real provider requirement/value cannot be represented without information loss by the typed publication vocabulary;
- human disambiguation of correspondence candidates requires structured semantics stronger than bounded presentation strings;
- current vs historical presentation cannot remain correctly separated under a proven workflow;
- a Product read requires cross-owner presentation whose authorization cannot be preserved without a new semantic dependency;
- the source/authored/observed media presentation trust separation cannot support a real B23 consumer;
- a second provider proves a shared presentation/value abstraction materially reduces total complexity rather than merely increasing generality;
- implementation evidence shows a new operation/Permission/owner is actually required.

## 17. Explicit non-goals

This prerequisite does not create:

- Product/PIM master;
- generic `EntityRef`, `EntityPresentation`, metadata map or relationship graph;
- PresentationService or label registry;
- provider field bag;
- generic option/category/unit catalog business owner;
- mapping/transformation/rules engine;
- generic media/asset service;
- direct frontend provider calls;
- new Product search/list operations by symmetry;
- caller-controlled sorting/field projection/expand DSL;
- implementation/runtime topology.

## 18. Written-spec exit gate

Before implementation planning:

- operator reviews this written spec;
- placeholders/ambiguities/contradictions are resolved;
- prerequisite scope is accepted as sufficient Global Maximum rather than Local Maximum or speculative platform expansion.

Only then create an implementation plan and touch D4/D4-R1/W2/W3/OAD/frontend evidence.
