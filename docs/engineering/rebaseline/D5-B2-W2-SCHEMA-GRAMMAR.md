# D5-B2 — W2 Request / Response Schema Grammar

> **Status:** OPEN / ACTIVE — W2-A Core Schema Grammar + W2-B ListingIntent/PriceIntent/Availability **ACCEPTED IN-STAGE**; W2-C Readiness + Market Intelligence + Economics next  
> **Parent Wire Contract:** `D5-B2-WIRE-CONTRACT.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + ratified D5-B2 Whole-Matrix + Wire W1  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **W2-A accepted:** 2026-08-18  
> **W2-B accepted:** 2026-08-18

## 1. Purpose

W2 converts the ratified operation inventory and W1 path/HTTP grammar into precise request/response schema laws without introducing a second cross-domain ontology for serialization convenience.

> **Reusable wire primitives may normalize representation; business meaning, knowledge, lifecycle, evidence and outcomes remain owner-specific.**

W2 is not derived from legacy OpenAPI, provider DTOs, database rows or frontend forms. It does not choose D6 UI composition, D7 server/generator/blob/persistence/runtime realization, D8 live-effect proof or implementation.

The complete W2 package will receive one coherent independent Fable review after its internal sections converge. Do not fragment W2 into micro-reviews unless a material contradiction requires it.

---

# 2. W2-A — Core Schema Grammar — ACCEPTED IN-STAGE

## 2.1 Governing invariant

> **A Product API schema preserves every material distinction the accepted owner can make, while adding no independent business meaning merely for uniform serialization.**

Corollaries:

1. reusable primitives are value/mechanism primitives, not horizontal business wrappers;
2. request schemas never allow clients to author server-owned authority/history fields;
3. unknown/unavailable/partial/not-applicable semantics are explicit where material and never hidden in `null`, zero, false or empty collections;
4. provider/source qualification remains explicit and mechanically distinguishable;
5. schema reuse must not create a universal entity/evidence/result/operation graph.

## 2.2 MPC-owned identifier grammar

MPC-owned canonical identifiers serialize as opaque non-empty strings.

Named semantic schemas such as `OrganizationId`, `ListingIntentId`, `PriceIntentId`, `FulfillmentNodeId` may exist for documentation/tooling, but their wire value does not promise UUID, ULID, database sequence or any other physical encoding.

Do not expose implementation-derived regex/version/timestamp meaning in the ID contract.

## 2.3 Native external identifier grammar

Provider/business-system native keys serialize as opaque strings even when a current provider happens to use numeric identifiers.

Examples include `native_listing_key`, `native_sale_key`, `native_shipment_key`, `native_product_key` and native business-document keys when materially exposed.

The Product API does not invite arithmetic over identifiers, narrow them to JavaScript-safe integers or assume a provider can never evolve an identifier to alphanumeric form.

## 2.4 Typed source-qualified references

No universal `ExternalRef`, `{entity_type, entity_id}`, provider graph or relationship graph is admitted.

Use semantic typed references, for example conceptually:

```json
{
  "marketplace_installation_id": "...",
  "native_listing_key": "MLB123"
}
```

```json
{
  "source_instance_id": "...",
  "native_product_key": "42664"
}
```

A request already scoped by `/organizations/{organization_id}` does not repeat an effective Organization field inside each reference. The server resolves every Organization-owned qualifier/reference inside the path Organization and fails closed on cross-Organization resolution.

## 2.5 Exact decimal primitive and Money

D2 exactness is preserved on the wire with an `ExactDecimalString` primitive.

Baseline lexical properties:

- JSON string, not JSON number;
- ordinary base-10 decimal notation;
- optional leading minus when the owning semantic allows negative values;
- one or more integer digits;
- optional fractional part;
- no exponent notation;
- no binary-floating-point conversion as the authoritative contract representation.

`Money` is:

```json
{
  "amount": "189.90",
  "currency": "BRL"
}
```

where `amount` is `ExactDecimalString` and currency is an explicit ISO-4217-style alpha-3 code in the Product API contract.

No universal minor-unit representation or global `round(2)` rule is introduced. Domain/provider rounding remains owned by the applicable accepted rule.

Rates and exact material quantities may reuse `ExactDecimalString` when exactness matters, while their semantic constraints remain owner-specific. W2 does not create a generic Unit-of-Measure/conversion framework.

## 2.6 Temporal grammar

Material instants use unambiguous date-time strings. Field names preserve the accepted temporal meaning, for example `observed_at`, `recorded_at`, `decided_at`, `effective_at` or an externally authoritative deadline/window field.

A generic `timestamp` or response-generation time never substitutes for source/effective/observation/decision time.

Timezone/local-calendar structures are introduced only where an accepted business/external rule actually depends on them.

## 2.7 Request and response schemas are not one CRUD blob

Create/update request schemas and resource/read response schemas are separate contracts when their authorities differ.

Examples of the pattern:

```text
CreateListingIntentRequest
UpdateListingIntentRequest
ListingIntent
```

Server-owned fields such as canonical ID, actor attribution, decision/history metadata, convergence evidence and current owner state are absent from client-write schemas rather than merely relying on clients to ignore them.

## 2.8 Closed semantic objects

Product semantic request objects are closed by default: undeclared properties are contract errors.

This is a wire-level guard against provider DTO/property-bag leakage and accidental client-authored authority fields.

Response objects are likewise closed where this materially protects the semantic contract; provider enrichment uses explicit bounded variants rather than `additionalProperties` property bags.

The exact JSON Schema/OpenAPI keywords used for closure depend on the later selected OAS/schema version, but the semantic closure requirement is accepted here.

## 2.9 `null` is not a knowledge-state transport

`null` never means unknown, unavailable, partial, not applicable or not observed by convention.

A nullable value is admitted only when semantic null itself is an accepted value for that specific field. Omission and null are not silently overloaded to carry multiple knowledge meanings.

## 2.10 Knowledge-state grammar

No universal Product API `Fact<T>` envelope is admitted.

Where knowledge state is material, use the smallest owner/field-specific discriminated union needed by that meaning. Potential states are drawn only from those materially applicable, such as known value, known empty/absent, unknown/insufficiently known, unavailable, partial/incomplete coverage and not-applicable/unsupported where semantically distinct.

Different meanings need not expose identical state sets. A known zero/false/empty value remains distinguishable from unknown/absent when material.

## 2.11 Union validation grammar

Material unions are mechanically exclusive.

Baseline schema technique is conceptually:

```text
oneOf
+ per-variant const discriminator value
```

Use `state` for knowledge/lifecycle-state discriminants and `kind` for semantic type/target unions where practical.

An OpenAPI `discriminator` may later be added as a tooling hint if useful, but correctness does not depend on it; the underlying union schemas must validate independently and exclusively.

## 2.12 Freshness and provenance locality

Freshness/provenance metadata attaches at the smallest semantic unit for which it is uniformly true.

Do not create a universal top-level `metadata`, `Evidence`, `Observation` or generated-at envelope.

Where one owner resource shares one observation context, resource-level provenance may be sufficient. Where components have different sources/times, each component preserves the materially correct provenance/time instead of inheriting a fabricated global timestamp.

## 2.13 No universal result/operation envelope

Do not wrap all Product operations in a generic shape such as:

```json
{ "status": "accepted", "data": {}, "error": null }
```

Resource creation returns the created owner resource; resource mutation/capability returns the current relevant owner resource when that honestly represents the outcome; stateless evaluation returns an operation-specific evaluation schema.

Create `oneOf` outcome variants only when one operation has materially distinct valid outcomes that cannot be represented honestly in the returned owner resource.

No generic `Operation`, `CapabilityResult`, `CommandResult` or universal async tracking identity is introduced.

## 2.14 Business outcomes versus HTTP problems

Valid domain `pending`, `rejected`, `ambiguous`, approval-required or unsupported/external-required outcomes remain semantic response meaning where the owner contract admits them.

They are not automatically HTTP 4xx/5xx and do not become `403` merely because the caller dislikes the outcome.

`202 Accepted` is not a baseline synonym for external effect not converged. Owner-local durable Intents already provide tracking identity. Use `202` only if a later concrete HTTP interaction genuinely has HTTP-level deferred processing that cannot be represented better by the current owner resource/result.

## 2.15 Typed PATCH/update bodies

Baseline Product updates use operation-specific typed JSON request bodies, not JSON Patch and not a generic mutation DSL.

For a typed update request:

- omitted property = unchanged;
- present property carries the complete semantic replacement for that field/relationship according to owner rules;
- `null` is not a generic clear instruction;
- clear/inherit/reset operations are explicit semantic variants where the domain needs them;
- array/list fields, when present, represent the intended complete semantic value unless the operation explicitly defines another owner-specific meaning.

Generic JSON Merge Patch is not adopted by default because `null`, arrays and discriminated knowledge/value semantics can carry material domain meaning.

## 2.16 Problem Details baseline

RFC 9457 Problem Details remains the API-level failure shape.

The problem `type` is the primary stable machine identifier for the problem class. Baseline fields are proportionately `type`, `title`, `status`, `detail`, `instance`.

Do not introduce an independent duplicate top-level `code` taxonomy by default. Add a stable MPC extension only when a real programmatic consumer needs information not already represented by the problem type/extensions.

Validation-detail extensions may expose bounded machine-readable locations/details without provider payloads, secrets or PII.

## 2.17 Provider-enriched evidence schema

Provider-enriched evidence is an owner-local closed variant, not a provider property bag.

Use bounded discriminated variants only for named Product 1.0 consumer/correctness needs. Another provider may honestly expose unsupported/not-applicable/unavailable semantics rather than fabricated equivalence.

Raw provider DTOs, arbitrary provider error text and free-form provider-field maps do not cross the Product API boundary.

## 2.18 Server-authority field fence

Client request schemas never contain effective-authority claims such as effective Principal/`created_by`/`approved_by` as caller identity, effective Organization duplicating path scope, `authorized=true`, `approved=true`, `converged=true` or provider-success assertions, or server-owned history/convergence/evidence fields.

A Principal/Organization ID may appear only when it is the legitimate business/admin subject of that operation, e.g. assigning an AccessRole to another Principal, never as the effective actor/scope authority.

## 2.19 W2-A proof / negative controls

Later executable contract proof must be able to falsify at least:

1. an MPC ID schema constraining clients to a physical UUID/ULID representation without accepted need;
2. a native external key becoming numeric/arithmetic wire state;
3. a bare external key accepted where its Marketplace Installation/SourceInstance qualifier is not otherwise unambiguous;
4. a universal `{entity_type, entity_id}` / `ExternalRef` replacing typed references;
5. Money represented as authoritative JSON floating-point number;
6. a global minor-unit/`round(2)` rule silently imposed on domain/provider economics;
7. unknown/unavailable/partial collapsing into `null`, zero, false or empty collection;
8. an invalid knowledge/type union matching multiple or no legal variants while passing validation;
9. client request containing an undeclared provider field/property bag;
10. one request/response CRUD schema making server-owned authority/history writable;
11. a generic Result/Operation wrapper becoming required Product ontology;
12. generic JSON Patch/merge semantics using `null` or array behavior to bypass owner meaning;
13. response generation time impersonating source/owner observation time;
14. arbitrary provider enrichment/property bag escaping its bounded owner-local variant;
15. request attempting to author effective Principal, Organization scope, authorization or convergence state.

## 2.20 W2-A outcome

**Parent architecture/W1:** `CURRENT STRUCTURE CONFIRMED`.

> **Use strong reusable scalar and union mechanics with owner-specific semantic schemas. Preserve exact decimals, source qualification, knowledge state, provenance and server authority without exporting a universal Fact/Evidence/Result/Operation/resource wrapper ontology.**

No parent-stage reopen is required.

---

# 3. W2-B — ListingIntent + PriceIntent + Availability — ACCEPTED IN-STAGE

## 3.1 Governing invariant

> **Dynamic marketplace publication is represented through Readiness-qualified requirements resolved by Offering-owned ListingIntent, while PriceIntent and Availability remain independent owner meanings before and after listing creation. Clients author each owner meaning through its own contract; D4/D7 may correlate and physically serialize them together without merging wire authority.**

This section stress-tests W2-A against authoring/control semantics. It does not create a PIM, provider field bag, mutable PriceDraft or public AvailabilityIntent authoring API.

## 3.2 ListingIntent is one create/edit authoring identity

`ListingIntent` remains one MPC-owned identity with a target union:

```text
new_listing
  → Marketplace Installation target, no provider Listing yet

existing_listing
  → source-qualified Listing target
```

Creation and editing do not gain separate aggregate/resource types.

Conceptual create input:

```json
{
  "source_product": {
    "source_instance_id": "...",
    "native_product_key": "42664"
  },
  "target": {
    "kind": "new_listing",
    "marketplace_installation_id": "..."
  },
  "desired": {
    "lifecycle": "active",
    "publication_context": {
      "native_category_key": "...",
      "native_product_type_key": "..."
    },
    "requirements_revision": "opaque-revision",
    "requirement_resolutions": []
  }
}
```

The draft may be created before all publication requirements are resolved. A contract-valid but incomplete draft is an MPC business draft with explicit blockers, not an API-schema error.

For edit, `target.kind = existing_listing` carries a source-qualified Listing reference. The same ListingIntent lifecycle is used.

## 3.3 Sparse declarative desired meaning; no provider mirror

A ListingIntent stores/serializes only Offering meaning that the intent seeks to establish/maintain. It does not copy the provider Listing wholesale merely to construct a full desired-state mirror.

For an edit intent, an unmentioned provider/listing aspect is outside that intent's desired change scope; a mentioned aspect carries explicit declarative desired meaning.

No procedural mini-DSL such as add/remove/replace provider field is admitted.

## 3.4 Desired lifecycle is state, not provider command vocabulary

Where the selected Listing lifecycle supports it, desired lifecycle is declarative owner state such as `active`, `paused` or `closed`, not separate provider-shaped `pause`, `reactivate`, `close` request types.

The selected initial publication lane remains active creation. A hypothetical paused/zero-quantity creation path is not admitted without D4/D8 proof.

## 3.5 Requirements and publication context

Dynamic provider/category/product-type publication shape enters through Readiness-owned requirement meaning, not a generic ProductAttribute/provider property bag.

A ListingIntent may preserve proportionately:

- publication context needed for desired representation, such as provider-qualified category/product-type selection;
- an opaque `requirements_revision` identifying the requirement meaning the draft was resolved against;
- one resolution per material `requirement_key`.

`requirement_key` is Readiness-owned within the relevant Product + marketplace/publication context. It is not required to equal a raw provider attribute ID, database column or JSON pointer.

The provider/native requirement identifiers may appear only as bounded source-qualified evidence/enrichment where needed for explanation/correctness.

`requirements_revision` is distinct from the ListingIntent `ETag`: requirement evidence may change without a concurrent write to the draft. Such drift changes dispatchability/business state, not HTTP representation concurrency.

## 3.6 Requirement resolution union

Each requirement resolution has exactly one of two baseline meanings:

```text
FOLLOW_SOURCE(source_candidate_key)
EXPLICIT_OVERRIDE(PublicationValue)
```

Conceptually:

```json
{
  "requirement_key": "title",
  "resolution": {
    "kind": "follow_source",
    "source_candidate_key": "..."
  }
}
```

or:

```json
{
  "requirement_key": "title",
  "resolution": {
    "kind": "explicit_override",
    "value": {
      "kind": "text",
      "value": "Torneira Monocomando"
    }
  }
}
```

No `fallback`, `derived`, expression, transformation, mapping or last-write-wins mode is admitted.

`source_candidate_key` is opaque and context-bound to Readiness meaning. It is not a new canonical entity identity and never exposes a business-system column as the Product API's source vocabulary.

At freeze/dispatch points, FOLLOW_SOURCE is re-resolved against current Readiness/source evidence; missing/unknown remains missing/unknown rather than falling back to an override.

## 3.7 Bounded PublicationValue

W2-B admits a small closed value union sufficient for the accepted authoring seam rather than `any`/arbitrary JSON or a speculative universal attribute system.

Baseline kinds:

- `text`;
- `exact_decimal`;
- `boolean`;
- `option`;
- `text_list`;
- `option_list`.

The applicable Readiness requirement determines which kind/cardinality/options are valid.

An option uses an opaque requirement-scoped `option_key`; it does not make raw provider option codes the Product API's canonical ontology.

Nested arbitrary objects/maps/expressions are not baseline. A future provider requirement that cannot be represented without material information loss is a concrete W2/D4-R1 reopen trigger for the smallest additional value form, not permission for a generic property bag now.

## 3.8 Typed ListingIntent update

`UpdateListingIntentRequest` is a typed PATCH body. Omitted top-level properties remain unchanged.

When the `desired` section is present, it represents the complete desired section for that draft revision rather than procedural incremental mutations. Arrays such as requirement resolutions or media selection are complete intended semantic values when present.

No JSON Patch/merge-patch null semantics, add/remove commands or duplicate ordering axes are introduced.

Current-state protection remains `If-Match` from W1.

## 3.9 Media reference/selection grammar

Listing media selection is a closed union of legitimate origins:

```text
source_media(source_media_candidate_key)
authored_media(listing_intent_media_id)
```

Example:

```json
[
  {
    "kind": "authored_media",
    "listing_intent_media_id": "...",
    "role": "primary"
  },
  {
    "kind": "source_media",
    "source_media_candidate_key": "...",
    "role": "gallery"
  }
]
```

Array order is publication order; there is no second numeric `position` authority.

`CreateListingIntentMedia` mints only a ListingIntent-scoped authored-media identity required for selection/history. It does not create ProductAsset/Asset/media-library authority.

Arbitrary external URLs are not trusted authored media. Source media remains a Readiness/D4 acquisition/evidence path. Exact binary upload/storage/hash/resizing/CDN/content-delivery mechanics remain D7.

`GetListingIntent` carries the media references/descriptors needed to explain/select the draft; no standalone Product media GET is admitted by symmetry.

## 3.10 ListingIntent response separates semantic axes

A ListingIntent response does not collapse all progress into one `status`.

It preserves proportionately distinct axes such as:

- canonical identity and create/edit target;
- source Product reference;
- desired Offering meaning;
- current resolved requirement values/provenance;
- authored-media descriptors/selection;
- Intent lifecycle (`draft | submitted | discarded` proportionately);
- draft dispatchability/blockers;
- current required execution-input readiness/correlation references;
- external-effect evidence/state where material;
- Listing representation convergence;
- actor/time attribution.

The exact field names/variant inventories are finalized with the remaining W2 owner schemas, but these meanings may not be collapsed into one provider/workflow state.

## 3.11 FOLLOW_SOURCE read-back and override provenance

For a FOLLOW_SOURCE resolution, the response may expose the current resolved value using a ListingIntent-specific knowledge union, for example known value, unknown or unavailable, with material source observation/provenance time.

An EXPLICIT_OVERRIDE response preserves server-attributed author Principal and authored time. Client requests do not author those fields.

The wire does not export a universal `Fact<PublicationValue>` type to achieve this.

## 3.12 PriceIntent target duality

`PriceIntent` remains one MPC-owned identity with two target variants:

```text
existing_listing
  → source-qualified Listing

pre_creation_listing_intent
  → ListingIntentId
```

Conceptually:

```json
{
  "target": {
    "kind": "pre_creation_listing_intent",
    "listing_intent_id": "..."
  },
  "desired_price": {
    "amount": "189.90",
    "currency": "BRL"
  }
}
```

or:

```json
{
  "target": {
    "kind": "existing_listing",
    "listing": {
      "marketplace_installation_id": "...",
      "native_listing_key": "MLB123"
    }
  },
  "desired_price": {
    "amount": "199.90",
    "currency": "BRL"
  }
}
```

Price is never ListingIntent content. The same PriceIntent architecture applies from initial publication onward.

## 3.13 Explicit PriceIntent supersession

A changed pending/pre-dispatch price is represented by a newer PriceIntent with explicit lineage, not by mutating a PriceDraft or taking latest timestamp wins.

Conceptually:

```json
{
  "target": { "kind": "pre_creation_listing_intent", "listing_intent_id": "..." },
  "desired_price": { "amount": "179.90", "currency": "BRL" },
  "supersedes_price_intent_id": "..."
}
```

The first PriceIntent has no supersedes reference. If a current standing PriceIntent exists for the same semantic target, a new conflicting PriceIntent cannot silently replace it without valid explicit supersession semantics.

Automation recurrence cannot silently supersede standing human-authored price intent. Attribution/history remains server-owned.

Withdrawal/cancellation remains deferred; supersession does not accidentally create a general PriceIntent mutation lifecycle.

## 3.14 Correlation is server-established, not client orchestration

A client does **not** PATCH `price_intent_id` or `availability_intent_id` into ListingIntent after creating them.

For pre-creation price, the correlation is already carried by `PriceIntent.target = pre_creation_listing_intent`.

Availability obtains the target through the accepted Offering→Availability semantic edge. The server identifies/revalidates the current owner-issued inputs required for dispatch and preserves the material correlation in historical effect/freeze context.

`GetListingIntent` may expose a typed current correlated PriceIntent reference for explanation; the price value remains PriceIntent-owned.

This removes unnecessary client orchestration/race windows without merging owner semantics.

## 3.15 SubmitListingIntent request body is empty of business meaning

The submit capability means submit the current ListingIntent revision:

```http
POST /organizations/{organization_id}/listing-intents/{listing_intent_id}:submit
If-Match: "opaque-validator"
```

No request body may carry price, quantity, approval, provider payload, `execute=true` or client-selected effective authority.

The owner freezes/accepts the current intent revision and revalidates current dependencies through accepted semantic boundaries.

Submission can be accepted while later external dispatch is pending/blocked on current required PriceIntent, Availability, authorization or other prerequisites. Submission is never eternal authorization; dispatch-time validity is revalidated.

## 3.16 SellableAvailability target and semantic axes

Sellable Availability can be queried for a target before or after provider listing creation without creating two APIs:

```text
pre_creation_listing_intent → ListingIntentId
existing_listing            → source-qualified Listing
```

This target duality is Availability's reference to the Offering target, not Availability ownership of ListingIntent/Listing lifecycle.

The read shape separates at least four different questions:

1. **control/effective capability** — does MPC currently control this availability target (`mpc_managed`, `external_required`, `unsupported`, or another accepted owner state)?
2. **desired owner meaning** — what Sellable Availability does Availability currently own/derive for the target?
3. **provider observation** — what quantity/state is currently observed at the marketplace, with honest source knowledge/freshness?
4. **convergence** — do current desired and observed meanings sufficiently converge?

These are not one `available` field or one workflow status.

## 3.17 Desired and provider quantity knowledge

Desired availability may be:

```json
{
  "state": "known",
  "quantity": "8",
  "evaluated_at": "...",
  "basis": {
    "inventory_source_ids": ["..."],
    "policy_revision": "..."
  }
}
```

Known zero is a legitimate known quantity and remains distinct from unknown/unavailable.

Provider observation is a separate owner-consumed source-qualified knowledge shape, for example:

```json
{
  "state": "known",
  "quantity": "8",
  "observed_at": "..."
}
```

or honest unknown/unavailable/partial/not-applicable variants when materially reachable.

Provider observed quantity never becomes the Availability-owned desired value by serialization convenience.

## 3.18 Availability convergence

Convergence is represented separately, with the smallest owner-specific state set needed to distinguish meanings such as `pending`, `converged`, `divergent`, `unknown` or `not_applicable` where material.

Example invariant:

```text
desired = known
provider observation = unavailable
→ convergence may be unknown
```

Unavailability of provider observation never erases the known Availability-owned desired value.

When provider topology creates shared/blast-radius effects, provider-enriched convergence/effect evidence may preserve the actual affected source-qualified Listing refs when correctness requires it. This does not introduce provider User Product/Pack graph authority into Product API ontology.

## 3.19 Inventory Source minimal schema seam

W2-B freezes only enough Inventory Source shape to preserve D2/D1 identity without importing Sankhya/WMS ontology:

```json
{
  "display_name": "CD Uberlandia",
  "source_inventory_scopes": [
    {
      "source_instance_id": "...",
      "native_inventory_scope_key": "..."
    }
  ]
}
```

An Inventory Source may map to one or several external inventory scopes when accepted business meaning requires composition.

Native `CODEMP`, `CODLOC`, warehouse IDs or equivalent remain D4/provider/source references, not MPC Inventory Source identity.

Exact allocation-policy default/inheritance/override schema is intentionally left for the policy/configuration W2 section so the same pattern can be challenged with Commercial Economics and Fulfillment configuration rather than duplicated here.

## 3.20 W2-B negative controls

Later contract proof must make at least these invalid/unreachable:

1. ListingIntent request containing price or available quantity;
2. ListingIntent request containing raw provider payload/property bag;
3. FOLLOW_SOURCE carrying a hidden manual fallback;
4. EXPLICIT_OVERRIDE client-authored actor/provenance fields;
5. duplicate material `requirement_key` resolutions in one desired section where only one resolution is legal;
6. arbitrary external URL treated as trusted authored media;
7. source-media/provider-media topology becoming ProductAsset/media-master authority;
8. media array order conflicting with a second `position` field;
9. PriceIntent desired price sent as JSON number;
10. PriceIntent target missing its required qualified Listing/ListingIntent reference or matching multiple union variants;
11. a new PriceIntent silently replacing a current standing PriceIntent without explicit supersession semantics;
12. client PATCHing PriceIntent/Availability correlation fields into ListingIntent;
13. submit request carrying price/quantity/approval/provider-action payload;
14. provider observed availability silently becoming desired owner availability;
15. known zero availability collapsing to unknown/empty;
16. provider observation unavailable causing desired value to disappear;
17. raw source inventory codes becoming Inventory Source canonical identity.

## 3.21 W2-B outcome

**Parent architecture / W1 / W2-A:** `CURRENT STRUCTURE CONFIRMED`.

> **Use Readiness-qualified requirement resolutions as the dynamic listing-authoring seam; keep ListingIntent sparse/declarative rather than provider-mirror state; keep PriceIntent and Availability independently targetable before/after listing creation; and derive required cross-intent correlations server-side instead of making clients orchestrate owner state.**

No parent-stage reopen is required.

---

# 4. Exact next W2 work

**W2-C — Product & Channel Readiness + Market Intelligence + Commercial Economics schema grammar.**

W2-C must stress-test the W2-A knowledge/provenance/value grammar against source evidence, dynamic requirement sets, market comparability and economic lineage. It must decide proportionately:

1. source Product search/reference result shapes without Product mirror identity;
2. `ProductChannelReadiness` keyed-Q shape, readiness conclusion and blockers without one mutable readiness row/ID;
3. `PublicationRequirements`/requirement/candidate/option schemas used by W2-B, including requirement revision and source/media candidate references;
4. competitive-position and comparable-offer schemas with source-qualified provider-rich evidence, explicit evidence sufficiency/coverage/freshness and no universal MarketObservation resource;
5. current Expected Economics schema with component-level knowledge/provenance and no plausible complete profitability when tax/fee/shipping/cost evidence is incomplete;
6. `EvaluatePriceScenario` request/result: caller supplies legitimate scenario variables only, never authoritative evidence overrides;
7. `SaleEconomics` L0/L1/L2 + R1/R2 schema grammar preserving historical rungs/occurrences rather than a mutable profitability blob;
8. Economic Attribution exact/partial/ambiguous/unresolved subject/reference grammar without universal entity graph;
9. bounded Commercial Policy schema only as needed to validate W2-A policy update semantics; generic Rules DSL remains rejected;
10. negative controls proving source/provider evidence cannot be client-fabricated as known authoritative economics or competitive truth.

W2-C does not yet choose collection pagination/filter grammar or full policy default/inheritance/override schema if a later cross-owner configuration section can adjudicate that pattern more coherently.

After W2-C, continue the remaining W2 owner/schema families necessary for global coherence. **Run one Fable review only after W2 as a coherent package converges**, following the canonical Standard Fable review workflow.

Implementation remains blocked until D9.
