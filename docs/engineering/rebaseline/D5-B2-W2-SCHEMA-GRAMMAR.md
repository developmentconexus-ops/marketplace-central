# D5-B2 — W2 Request / Response Schema Grammar

> **Status:** ACCEPTED / CANONICAL — Whole-W2 operator-ratified  
> **Parent Wire Contract:** `D5-B2-WIRE-CONTRACT.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + Decision Reconciliation Baseline + ratified D5-B2 Whole-Matrix + canonical Wire W1  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Whole-W2 final ratification incorporated:** 2026-08-19  
> **Whole-W3 cursor Problem Details amendment incorporated:** 2026-08-19  
> **Final Problem/media consistency incorporated:** 2026-08-19

## 1. Purpose and authority

W2 is the **single canonical Product API request/response schema authority** for the admitted D5-B2 operation surface.

It consolidates the formerly staged W2-A/B/C/D/E work and the operator-ratified Whole-W2 adversarial review. Git history preserves those staging snapshots; they are not parallel active schema authorities.

> **A Product API schema preserves every material distinction the accepted owner can make, while adding no independent business meaning merely for uniform serialization. Reusable wire primitives normalize representation; business meaning, knowledge, lifecycle, evidence and outcomes remain owner-specific.**

W2 does not choose D6 screen/BFF topology, D7 server/storage/transaction/job realization, D8 live-effect proof or implementation.

---

# 2. Core schema grammar

## 2.1 Opaque MPC identifiers

MPC-owned canonical IDs serialize as opaque non-empty strings.

Named semantic schemas such as `OrganizationId`, `ListingIntentId`, `PriceIntentId`, `FulfillmentNodeId`, `FulfillmentExecutionId`, `PostSaleResolutionId`, `AuthorizationDelegationId` and `WorkId` may exist for documentation/tooling, but wire values do not promise UUID, ULID, database sequence or time encoding.

Do not expose implementation-derived regex/version meaning in IDs.

## 2.2 External/native identifiers

Provider/business-system native keys also serialize as opaque strings even when a current provider uses numeric identifiers.

Examples include `native_product_key`, `native_listing_key`, `native_sale_key`, `native_shipment_key`, native business-document keys and provider movement keys.

Identifiers are never quantities and are not narrowed to JavaScript-safe integers.

## 2.3 Typed source-qualified references

There is no universal `ExternalRef`, `{entity_type, entity_id}`, provider graph or relationship graph.

Use semantic typed references, for example:

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

A request already scoped by `/organizations/{organization_id}` does not repeat effective Organization inside every reference. Every Organization-owned qualifier/reference resolves fail-closed inside the path Organization.

## 2.4 Exact decimals and Money

Exact decision/economic numbers use `ExactDecimalString`:

- JSON string, not JSON number;
- ordinary base-10 notation;
- optional leading minus only when the owning semantic permits negative values;
- one or more integer digits;
- optional fractional part;
- no exponent notation;
- no authoritative binary floating-point conversion.

`Money` is:

```json
{
  "amount": "189.90",
  "currency": "BRL"
}
```

Currency is explicit alpha-3 currency code. There is no universal minor-unit representation or global `round(2)` rule.

Rates and exact material quantities may reuse `ExactDecimalString`; owner-specific scale/unit/rounding constraints remain owner meaning.

## 2.5 Temporal grammar

Material instants use unambiguous date-time strings with names preserving actual meaning, e.g. `observed_at`, `recorded_at`, `decided_at`, `authored_at`, `effective_at`, `created_at`, `updated_at`.

A generic `timestamp`/response-generation time never substitutes for source/effective/observation/decision time.

## 2.6 Request schemas are not read schemas

Create/update request bodies and resource/read responses are distinct whenever their authorities differ.

Server-owned identity, actor attribution, decision/history metadata, provider outcome evidence, current knowledge, convergence and audit fields are absent from client-write schemas rather than merely marked read-only by convention.

## 2.7 Closed semantic objects

Product semantic request objects are closed by default. Undeclared properties are contract errors.

Responses are likewise closed where this materially protects the semantic boundary. Provider enrichment uses explicit bounded variants, never `provider_fields`, arbitrary maps or raw DTO passthrough.

Exact JSON Schema closure keywords depend on the later selected OpenAPI/schema version; the semantic closure requirement is binding.

## 2.8 `null` is not a knowledge-state carrier

`null` never means unknown, unavailable, partial, absent or not-applicable by convention.

A nullable value is allowed only when semantic null itself is a legitimate explicit value for that field.

## 2.9 Knowledge-state grammar

There is no universal Product API `Fact<T>` envelope.

Where knowledge state is material, use the smallest owner/field-specific discriminated union required by that meaning. Potential variants include only those genuinely reachable, such as:

- known value;
- known empty/absent;
- unknown/insufficiently known;
- unavailable;
- partial/incomplete coverage;
- not-applicable/unsupported.

Known zero/false/empty remains distinct from unknown/absent.

## 2.10 Union validation

Material unions are mechanically exclusive using `oneOf` plus a fixed per-variant discriminant.

Use `kind` for semantic type/target unions and `state` for knowledge/lifecycle meanings where practical.

An OpenAPI `discriminator` may later be added for tooling but is never correctness authority.

## 2.11 Provenance/freshness locality

Freshness/provenance attaches at the smallest semantic unit for which it is uniformly true.

Do not create universal top-level `metadata`, `Evidence`, `Observation` or generated-at envelopes.

## 2.12 No universal result/operation envelope

Do not wrap all Product operations in `Result<T>`, `Operation`, `CommandResult`, `CapabilityResult` or a universal async tracking identity.

Resource creation returns the created resource; mutation/capability returns the current relevant resource or operation-specific semantic result; stateless evaluation returns its evaluation schema.

## 2.13 Business outcomes versus HTTP problems

Valid owner outcomes such as pending, rejected, ambiguous, insufficient-evidence, approval-required, external-required or unsupported remain Product semantics where admitted. They are not automatically HTTP 4xx/5xx.

`202 Accepted` is not a generic synonym for external effect not converged. Durable owner Intents already provide tracking identity.

## 2.14 Typed PATCH/update

Baseline updates use operation-specific typed JSON bodies, not JSON Patch and not generic Merge Patch.

For a typed update:

- omitted field = unchanged;
- present field carries the complete intended semantic replacement for that field/section;
- `null` is not generic clear;
- clear/inherit/reset is explicit semantic meaning where needed;
- arrays, when present, represent the complete intended semantic collection unless the operation explicitly defines another owner-specific meaning.

## 2.15 Problem Details base

RFC 9457 Problem Details is the API-level failure representation.

`type` is the stable primary problem identifier. Baseline fields are proportionately `type`, `title`, `status`, `detail`, `instance`.

No duplicate top-level global `code` taxonomy is introduced by default. `title` and `detail` are human-facing and never stable programmatic identifiers. A bounded problem-specific extension exists only for structured data that problem genuinely requires, such as validation locations or a stale-validator pointer.

### 2.15.1 `about:blank` for status-only failures

A failure carrying stable MPC-specific semantics uses the canonical §19 custom problem type.

A failure with no additional Product semantics beyond its standard HTTP status uses:

```text
type   = about:blank
status = applicable standard status
title  = recommended status phrase
```

Do not mint an MPC problem type merely to repeat standard HTTP meaning. Product clients branch on custom `type`, and on HTTP `status` when `type = about:blank`.

### 2.15.2 `about:blank` never waives HTTP obligations

`about:blank` governs only the Problem Details body. Status-specific HTTP obligations remain binding:

- `405 Method Not Allowed` carries the current `Allow` header value required by RFC 9110;
- `415 Unsupported Media Type` may carry an honest `Accept`/`Accept-Encoding` value when the supported request representation/coding can be expressed without inventing part-level semantics;
- `413 Content Too Large` carries `Retry-After` only when the refusal is genuinely temporary; a standing Product size bound normally is not.

## 2.16 Server-authority fence

Client requests never author effective Principal/Organization, `created_by`, `approved_by`, `authorized=true`, `converged=true`, provider-success claims or server history/evidence.

Principal/Organization IDs appear only when they are the legitimate subject/target of that operation, never as the effective caller/scope authority.

---

# 3. Offering — ListingIntent, MarketplaceListing, PriceIntent

## 3.1 One ListingIntent create/edit identity

`ListingIntent` is the one MPC-owned authoring identity for both create and edit.

Target union:

```text
new_listing
  → Marketplace Installation target; no provider Listing yet

existing_listing
  → source-qualified Marketplace Listing
```

A draft may be contract-valid while business-incomplete. Missing publication requirements are explicit blockers, not schema errors merely because the draft exists.

## 3.2 Sparse declarative desired meaning

ListingIntent stores/serializes only Offering meaning the intent seeks to establish/maintain. It is not a full provider Listing mirror.

For edit, an unmentioned provider/listing aspect is outside that intent's desired-change scope. Mentioned aspects carry explicit declarative desired meaning.

No add/remove/replace provider-field mini-DSL is admitted.

## 3.3 Desired lifecycle

Where applicable, ListingIntent desired lifecycle is declarative owner state such as `active`, `paused`, `closed`, not provider-shaped commands.

Selected initial publication remains active creation. A paused/zero-quantity creation lane is not claimed without D4/D8 proof.

## 3.4 Publication context and requirements revision

ListingIntent may preserve:

- source-qualified Product;
- Marketplace Installation / existing Listing target;
- publication context such as provider-qualified category/product-type selection when part of desired representation;
- opaque `requirements_revision` identifying the Readiness/provider requirement basis the draft resolved against;
- one resolution per material `requirement_key`;
- listing-specific media selection/order/role.

`requirements_revision` is not the ListingIntent ETag. Requirement evidence may change without a ListingIntent write; such drift changes dispatchability, not HTTP resource concurrency.

## 3.5 Requirement resolution union

Exactly two baseline resolution authorities exist:

```text
FOLLOW_SOURCE(source_candidate_key)
EXPLICIT_OVERRIDE(PublicationValue)
```

No fallback, derived/formula, transform/mapping engine or last-write-wins mode.

FOLLOW_SOURCE candidate keys are opaque/context-bound Readiness references. At freeze/dispatch, the value is re-resolved from current source evidence. Unknown/unavailable remains unknown/unavailable; no hidden manual fallback.

EXPLICIT_OVERRIDE preserves server-attributed author Principal/time and never rewrites source Product truth.

## 3.6 PublicationValue — canonical bounded union

Baseline kinds:

- `text`;
- `exact_decimal`;
- `boolean`;
- `option`;
- `text_list`;
- `option_list`;
- `number_unit`;
- `not_applicable`.

`number_unit` is:

```json
{
  "kind": "number_unit",
  "value": "12.50",
  "unit_key": "..."
}
```

`unit_key` is requirement-scoped and opaque; Readiness publishes allowed/default unit keys where the requirement needs them; D4 maps to provider-native unit representation. No generic unit conversion/UoM engine is introduced.

`not_applicable` is permitted only when the current Readiness requirement explicitly allows N/A. It is an **EXPLICIT_OVERRIDE PublicationValue meaning**, never a third resolution mode and never inferred from source absence/unknown.

Nested arbitrary objects/maps/expressions are not baseline. A future real provider requirement that cannot be represented without information loss reopens only the smallest value-family extension.

## 3.7 Typed ListingIntent update

`UpdateListingIntentRequest` is typed PATCH. Omitted top-level properties remain unchanged.

When `desired` is present it is the complete desired section for that revision. Requirement-resolution/media arrays are complete intended values when present.

PATCH is a standard same-resource mutation and therefore uses W1 `If-Match`.

## 3.8 Listing media selection

Selection origin union:

```text
source_media(source_media_candidate_key)
authored_media(listing_intent_media_id)
```

Array order is publication order; no parallel numeric position authority.

`listing_intent_media_id` is ListingIntent-scoped identity, not ProductAsset/Asset/media-library authority.

Source media remains external evidence. Arbitrary external URLs are not trusted authored media.

## 3.9 CreateListingIntentMedia wire

Authored media intake is a ListingIntent capability:

```http
POST /organizations/{organization_id}/listing-intents/{listing_intent_id}:create-media
Idempotency-Key: ...
Content-Type: multipart/form-data
```

Multipart contains:

- one binary file part;
- required typed `etag` part carrying current ListingIntent validator;
- only bounded additional semantic metadata if later proven.

Do not use base64 JSON, arbitrary external URL, UploadSession, generic Asset resource or presigned-upload Product resource by default.

Success is `200` + operation-specific typed result (§3.9.1). There is no standalone media URI/GET, so no `201 + Location` resource fiction.

Exact idempotent replay returns the same `listing_intent_media_id`. Binary content identity + material multipart metadata + revision proof participate in semantic idempotency equivalence. D7 chooses digest/storage/CDN/resizing mechanics.

Exact accepted media families and size bounds must be explicit before implementation; their values and enforcement mechanism remain later realization.

### 3.9.1 Success result and parent validator advancement

Authored-media descriptors are one ListingIntent read axis (§3.10). Successful creation therefore changes the ListingIntent representation and advances its strong opaque owner validator.

The typed result carries:

```text
media                 → stable ListingIntentMediaDescriptor
listing_intent_etag   → current parent ListingIntent validator after acceptance
```

Do not place the parent validator in an HTTP `ETag` header on the distinct `:create-media` request URI; that URI is not the protected resource representation (§17.1). It travels as typed result data carrying the one parent-resource revision authority.

Consequences:

- two concurrent creates against the same ListingIntent validator serialize;
- the first accepted intake advances the validator;
- a materially different concurrent loser receives `409 resource-revision-conflict` and re-reads/rebases;
- a caller chains sequential successful uploads using the returned `listing_intent_etag`, without a mandatory GET after every success;
- an exact lost-response retry resolves the already accepted intake before stale-revision re-evaluation (§18) and returns the same intake result rather than creating another media ID.

No multi-file/bulk intake is admitted. Reopen only if D6 proves real multi-image authoring is unusable under this serialized model.

### 3.9.2 Exact failure grammar

Transport and representation:

| Failure | Product HTTP disposition |
|---|---|
| Request cannot be parsed as a multipart representation | `400 malformed-request` |
| Top-level request representation is not the selected multipart media type | `415` + `about:blank` |
| Binary file format is unsupported, undecodable or materially contradicts its declared content type after inspection | `415` + `about:blank` |
| Request/file exceeds the enforced bound | `413` + `about:blank` |

The operation contains exactly one binary file part, so an unsupported file representation remains a `415` representation-format failure rather than being relabeled as semantic validation merely to obtain a part pointer.

Contract, revision, idempotency and lifecycle:

| Failure | Product HTTP disposition |
|---|---|
| Required binary part or typed `etag` missing, duplicated or contract-invalid; semantic metadata invalid | `422 validation-error` with bounded field/part diagnostics |
| Supplied ListingIntent validator is stale | `409 resource-revision-conflict` |
| Current ListingIntent state does not admit media creation, after exact-repeat handling | `409 resource-state-conflict` |
| Required `Idempotency-Key` absent/invalid | `400 idempotency-key-required` |
| Same key reused with materially different bytes, metadata, target or revision proof | `422 idempotency-key-reused` |
| Equivalent prior intake still processing | `409 idempotency-request-in-progress` |
| Unexpected internal storage/scanning/transformation/runtime failure before a successful Product result is established | `500 internal-error` |

### 3.9.3 Transport-guard ordering

A transport-level refusal for wrong top-level representation or excess size:

- may occur before Product authentication;
- must be enforceable without full-body buffering merely to discover the violation;
- uses `about:blank` and discloses no Product resource existence, Organization, Membership, Permission or business state.

The §18 idempotency processing order governs semantic evaluation after transport admission; it does not require an unbounded request body to pass authentication first.

### 3.9.4 Leakage fence

Do not add Product problem types such as `blob-upload-failed`, `virus-scanner-error`, `cdn-error` or `provider-image-error`.

Raw binary, storage key, access locator, scanner result, provider payload, secret, stack detail and arbitrary PII never enter Product Problem Details.

### 3.9.5 Authored-media identity

For one ListingIntent:

```text
accepted authored binary + admitted metadata → one stable listing_intent_media_id
```

- the ID belongs to exactly one ListingIntent and cannot be referenced from another;
- it never rebinds to different bytes or material meaning;
- a materially different upload creates another ID through the same capability;
- current selection/order/role is ListingIntent desired-state meaning and changes only through the accepted draft update (§3.8);
- unselecting an authored ID does not by itself delete its bytes or historical reference;
- publication-attempt history (§3.11) retains only the identity/provenance needed to explain what was attempted;
- no Product `UpdateMedia`, `DeleteMedia`, media collection or ProductAsset CRUD is admitted.

### 3.9.6 Retention/erasure residual

Product 1.0 admits no client-facing authored-media erasure operation. Unselection is not erasure.

No accepted authority currently establishes a universal legal/contractual retention duration or an Organization-lifetime immortality rule. Retention/erasure remains **Unknown / deferred to D2 data ownership and D7 realization**, subject to these D5 constraints:

- content required by current selection or material historical explanation cannot be silently removed;
- a future retention/erasure rule must preserve enough historical identity/provenance to explain prior consequential attempts;
- a legal, privacy, contractual, operator or material cost obligation reopens the smallest D2/D7 scope rather than automatically creating Product delete CRUD.

### 3.9.7 Descriptor families

Identity/provenance and presentation access are separate schemas:

```text
ListingIntentMediaDescriptor
  → stable authored-media identity + bounded content/provenance facts
  → eligible for current selection and historical attempt basis

ListingIntentMediaPresentationDescriptor
  → ListingIntentMediaDescriptor + volatile access reference
  → response-only presentation aid
  → never persisted into history, idempotency fingerprint, logs or Problem Details
```

The `CreateListingIntentMedia` result returns the stable descriptor and the new parent validator, not a durable presentation locator. The client already possesses the submitted bytes; later server-side retrieval is a read concern through `GetListingIntent` under `offering.read`.

### 3.9.8 Authored-media byte delivery authority

Authored-media byte delivery is a **separately justified technical presentation surface** under D5-B1 route classification. It is:

- not a 96th Product operation;
- not a Technical Ingress lane A or B;
- not a Product SDK/OpenAPI business operation;
- not a generic Media/Asset service;
- not authorized by stable media ID alone.

The baseline authority law is:

> **A caller that cannot currently obtain the corresponding authored-media presentation descriptor through `GetListingIntent` under the exact path Organization and `offering.read` cannot obtain the bytes.**

The delivery realization therefore reuses current Product authentication, unique Principal binding, Principal access eligibility, Organization Membership and exact `offering.read` authorization for the referenced ListingIntent/media relationship. Exact route, proxy/storage/CDN topology, streaming and transformation mechanics remain D7.

A durable, anonymous or freely forwardable object-store/CDN locator is not baseline. A delegated bearer capability may be reconsidered only by an explicit smallest-scope D5/D7 reopen proving that authenticated delivery is materially unsuitable and preserving tenant/scope/expiry/non-enumerability and credential-handling constraints.

Technical media-delivery failures remain technical-surface failures; they do not expand the §19 Product Problem catalog.

If D6 proves that an embedded presentation reference plus the bounded technical surface cannot satisfy the real consumer, reopen only the smallest B2 operation/W4 surface before implementation. D7 may not invent a Product media GET privately.

### 3.9.9 Source and authored media remain distinct

- source media remains source-qualified external evidence owned by the Readiness/D4 seam;
- authored media remains ListingIntent-scoped MPC state;
- ListingIntent selection uses a bounded discriminated union (§3.8) because it must choose/order both origins;
- dimension/content-type primitives may be reused only where their meaning is genuinely identical;
- source-media locators and authored-media access references are distinct types even if their JSON shape happens to match, because issuer, trust, lifetime and governing authority differ;
- source candidate key and authored media ID never substitute for one another;
- arbitrary client URL remains rejected;
- provider image-error feeds remain deferred until a named consumer proves need.

No generic Media/Asset owner or cross-ListingIntent library is introduced.

## 3.10 ListingIntent read axes

A ListingIntent read preserves distinct axes proportionately:

- canonical identity and create/edit target;
- source Product;
- desired Offering meaning;
- current resolved requirement values/provenance;
- authored-media descriptors/selection;
- Intent lifecycle (`draft`, `submitted`, `discarded` proportionately);
- current draft dispatchability/blockers;
- required execution-input correlations;
- current external-effect evidence/state;
- Listing representation convergence;
- actor/time attribution;
- immutable historical dispatch/effect basis (§3.11).

No giant provider/workflow `status` replaces these dimensions.

## 3.11 Historical ListingIntent dispatch basis — append-only

Every consequential publication attempt that is materially required for explanation is preserved through the existing ListingIntent read/effect-history axis; no standalone `PublicationAttempt` Product resource/CRUD is created.

Each material attempt preserves proportionately:

- ListingIntent identity + exact submitted/material revision;
- source-qualified Product / Installation / target Listing;
- exact resolved Offering PublicationValues actually used;
- FOLLOW_SOURCE knowledge/value/provenance **as established at that attempt**, never by current re-resolution;
- EXPLICIT_OVERRIDE author provenance;
- provider requirement/schema revision materially used;
- media selection/order/role + material source/authored provenance;
- PriceIntent identity + exact material revision used, never ListingIntent-owned price;
- typed Availability-issued historical input/correlation sufficient to explain the attempt, never ListingIntent-owned quantity/current Availability authority;
- decision/disposition/AuthorizationDecision references materially used;
- intended/authorized/attempted scope where material;
- provider member/aspect outcomes and authoritative convergence evidence.

Multi-step/partial publication can accumulate append-only attempt/member/aspect occurrences. A mutable “latest attempt” blob is insufficient.

Historical snapshots may duplicate past values as historical owner-attributed context; they never become current cross-owner authority or payload archive/PIM mirror.

## 3.12 Submit / discard

`SubmitListingIntent` means submit the current ListingIntent revision. Its custom-method request carries no business payload but does carry the required typed current-resource validator:

```json
{
  "etag": "\"opaque-validator\""
}
```

Price, quantity, approval, provider payload or `execute=true` cannot appear.

Submission may be accepted while provider dispatch is blocked/pending on PriceIntent, Availability, authorization or provider-effective requirements. Dispatch revalidates current prerequisites; submission is never eternal authorization.

`DiscardListingIntentDraft` is likewise an owner capability using typed `etag` when current-state protection is required.

## 3.13 MarketplaceListing actual-state Q

MarketplaceListing remains a source-qualified external resource; no synthetic MPC Listing ID.

The bounded read shape preserves proportionately:

- qualified Listing ref;
- normalized observed lifecycle;
- publication context;
- bounded observed listing representation values/media relevant to current Offering semantics using PublicationValue/knowledge grammar rather than raw provider attributes;
- current observed marketplace price evidence where materially needed;
- source observation/freshness/provenance.

MarketplaceListing does **not** own ListingIntent convergence or PriceIntent convergence. Observed marketplace price is evidence only; price-convergence verdict remains PriceIntent-owned.

Within Offering, one bounded marketplace-price observation grammar is reused between MarketplaceListing and PriceIntent convergence evidence rather than creating two spellings. This is same-owner schema reuse, not universal Evidence ontology.

## 3.14 PriceIntent target duality

One PriceIntent identity supports:

```text
existing_listing
  → source-qualified Listing

pre_creation_listing_intent
  → ListingIntentId
```

`desired_price` is exact Money. Price never becomes ListingIntent content, including initial publication.

## 3.15 Explicit PriceIntent supersession

A changed pending/pre-dispatch price is a new PriceIntent with explicit lineage, not mutable PriceDraft/latest-time-wins.

When superseding, the typed reference carries:

```text
superseded PriceIntent ID
+ exact referenced etag
```

Automation recurrence cannot silently supersede standing human price intent. Attribution/history remains server-owned.

## 3.16 Cross-intent correlation is server-established

Clients do not PATCH `price_intent_id` or Availability correlations into ListingIntent.

Pre-creation PriceIntent already points at ListingIntent. Availability uses accepted owner edges. Server/D4/D7 identifies/revalidates current owner-issued inputs for dispatch and preserves material historical correlation.

---

# 4. Availability

## 4.1 SellableAvailability target

Availability supports before/after publication targets:

```text
pre_creation_listing_intent → ListingIntentId
existing_listing            → source-qualified Listing
```

This is an Availability reference to Offering target, not Availability ownership of ListingIntent lifecycle.

## 4.2 Four separate semantic axes

SellableAvailability separates:

1. **control/effective capability** — e.g. `mpc_managed`, `external_required`, `unsupported`;
2. **desired owner meaning** — current Sellable Availability quantity/basis;
3. **provider observation** — current provider observed quantity/state/freshness;
4. **convergence** — relation between desired and observed meaning.

Never collapse these into one `available` or workflow state.

## 4.3 Desired quantity knowledge

Known desired quantity uses exact decimal string where fractional quantity is legitimate and preserves evaluation basis such as InventorySource IDs/policy revision where materially explanatory.

Known zero remains known zero.

Unknown/unavailable are explicit variants and never become zero.

## 4.4 Provider observation and convergence

Provider observation is separate source-qualified evidence. It may be known, unknown, unavailable, partial or not-applicable as the selected lane requires.

Provider observed quantity never becomes Availability desired truth by serialization convenience.

Convergence uses the smallest owner-specific variants such as pending/converged/divergent/unknown/not-applicable.

Provider observation unavailable does not erase known desired value.

## 4.5 InventorySource seam

InventorySource is MPC identity distinct from native stock/location identity and from FulfillmentNode.

Baseline shape may include display name plus one or more source-qualified native inventory scopes:

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

Native company/location/warehouse codes remain D4 evidence/realization.

## 4.6 Allocation/scope policy

Availability owns its typed allocation/scope configuration. Cross-owner inherit/override mechanics are in §10; no generic rules DSL or arbitrary scope graph.

---

# 5. Product & Channel Readiness

## 5.1 Source Product search is evidence

`SearchSourceProductsForMarketplace` returns source-qualified search hits, not MPC Product resources.

Each hit is anchored by SourceInstance + native Product key and exposes only bounded facts needed by current readiness/correspondence consumers, e.g. name/SKU/GTIN evidence with honest knowledge/provenance.

No `product_id`, Product lifecycle or generic `attributes` bag.

## 5.2 ProductChannelReadiness keyed Q

ProductChannelReadiness has no synthetic ID. Subject:

```text
source-qualified Product
+ Marketplace Installation
```

Response preserves proportionately:

- subject;
- current correspondence meaning;
- **correspondence-scoped opaque `etag`**;
- publication readiness conclusion;
- blockers/insufficiency reasons;
- `requirements_revision` where material;
- evaluation/provenance time.

Readiness states distinguish at least `ready`, `blocked`, `unknown`, `unavailable` where reachable.

`blocked` means sufficient knowledge establishes an unmet condition; `unknown` means insufficient evidence to conclude.

## 5.3 ProductChannelCorrespondence state and capabilities

Correspondence is an owner-specific union such as resolved/unresolved/conflicting/unknown/unavailable.

Candidate keys are opaque/context-bound Readiness references, not provider Product IDs or canonical entities.

`ResolveProductChannelCorrespondence` and `ClearProductChannelCorrespondence` remain explicit Readiness capabilities over the keyed Product+Marketplace subject. They do not gain a synthetic Correspondence/Readiness ID and are **not forced into PUT/DELETE** merely for HTTP conditional syntax.

Both carry the `correspondence_etag` from the current ProductChannelReadiness representation as typed request revision proof.

Why typed proof is required: stale-state safety must protect unresolved/conflicting meaning as well as a standing resolved association. Modeling a resource that exists only when resolved leaves an absent-target precondition gap; modeling it as always-existing makes DELETE semantically false.

`correspondence_etag` is distinct from `requirements_revision` and broader Readiness/source-evidence revision, so unrelated requirement churn does not spuriously conflict with correspondence decisions.

Automation cannot silently supersede a standing human decision; that remains owner business validity, not transport authority.

## 5.4 PublicationRequirements keyed Q

PublicationRequirements is keyed by source Product + Marketplace Installation + publication context.

It preserves:

- subject;
- publication context;
- opaque requirements revision;
- requirements[];
- source media candidates[];
- evaluation/provenance.

Provider condition/expression language never becomes Product rules DSL.

## 5.5 Requirement value specs

Requirement specs align with PublicationValue and expose only constraints required for safe authoring, e.g. text lengths, option sets, allowed/default unit keys and explicit N/A permission where applicable.

Option/unit/candidate keys are opaque and requirement/context-scoped.

Source candidates expose current resolved value through requirement-specific knowledge semantics without source table/column/provider JSON-path leakage.

## 5.6 Applicability classes and draft-dependent requirements

Readiness owns requirement definition/key/value-spec/source candidates and may classify applicability as:

- current/unconditional/evaluable from present readiness context; or
- `draft_dependent` where provider evidence requires concrete owner-issued listing inputs.

W2 does **not** expose provider conditional expressions.

For draft-dependent provider validation:

```text
current ListingIntent revision
+ current PriceIntent identity/revision
+ current Availability basis
+ requirements revision
        ↓
D4 technical composition / provider validation
        ↓
provider-effective required-attribute evidence
        ↓
Offering-owned draft dispatchability/blockers
```

This is mechanism, not owner transfer:

- Readiness never owns mutable ListingIntent content;
- Offering never owns requirement definitions;
- Price/Availability remain separate authorities;
- client never submits raw provider validation payload/result as truth.

Evaluation evidence is anchored to the exact input basis. When any anchor is no longer current, the conclusion becomes stale/unknown and cannot masquerade as current dispatchability.

Provider validation unavailable remains unknown/unavailable, never “not required”.

The concrete selected User-Product lane validation endpoint remains D4/D8 proof; current classic Items documentation proves the failure class and seam, not the exact selected-lane endpoint.

---

# 6. Market Intelligence

## 6.1 Market subject

Bounded owner-local union:

```text
existing_listing
source_product_marketplace_context
```

Pre-listing market reasoning therefore does not require ListingIntent identity.

No universal CommerceSubject is introduced.

## 6.2 CompetitivePosition

Contextual Q with no synthetic ID. Preserves proportionately:

- subject;
- coverage;
- evidence sufficiency;
- own-offer evidence where available;
- best/relevant comparable;
- factual competitive relation/gap;
- change when a comparable prior basis exists;
- bounded provider enrichment;
- evaluation/provenance.

No mutable MarketObservation resource.

## 6.3 Coverage != sufficiency

Coverage describes the actually traversed/available source/provider universe, e.g. complete-for-declared-source-scope, partial, unknown, unavailable.

Evidence sufficiency independently describes whether the specific conclusion is supportable.

Provider enumeration completeness never means universal market completeness.

## 6.4 Comparable offers are evidence values

ComparableOffer receives no synthetic MPC ID. If provider exposes no stable native offer identity, MPC does not fabricate one.

Expose only material price/shipping/delivered-price/provenance evidence; competitor PII only when a named correctness/operations consumer requires it.

## 6.5 Provider-rich evidence

Provider-specific enrichment is a closed owner-local variant when materially useful, not a `provider_fields` map and not forced cross-provider equivalence.

---

# 7. Commercial Economics

## 7.1 Economics subject

Expected/scenario economics uses a bounded Economics-local subject union such as existing Listing or source Product + marketplace context. No universal SubjectRef.

## 7.2 ExpectedEconomics is components-first

Response:

```text
subject
components
conclusion
evaluated_at / material provenance
```

Components remain distinct, proportionately including sale/candidate price, selected Cost Basis, expected tax, expected marketplace fee, seller-borne shipping and material promotion/discount effects.

Each component uses Economics-local knowledge/provenance; equal numeric values with observed/modeled/configured/derived provenance remain distinguishable.

## 7.3 Fail-honest conclusion

Known contribution/margin/profitability exists only when required components are sufficiently known under current Economics rules.

Conclusion variants include proportionately `known`, `insufficient_evidence`, `unavailable`.

Unknown tax/fee/shipping/cost never defaults to zero to produce a plausible margin.

## 7.4 EvaluatePriceScenario

Side-effect-free capability accepts legitimate hypothetical variables only: subject, candidate Money, quantity or accepted basis-selection inputs when actually needed.

Caller cannot submit authoritative cost/tax/fee/competitor evidence replacements.

No Simulation/Recommendation ID, persistent simulation resource or price actuation authority.

## 7.5 SaleEconomics lineage

SaleEconomics is keyed by source-qualified Sale and preserves:

```text
expected_basis      L0 historical basis
order_economics     L1
realized_economics  L2
reconciliation      R1/R2
coverage/evaluated_at
```

No synthetic SaleEconomics/Reconciliation IDs.

Historical L0 is distinct from current ExpectedEconomics. L1 keeps order-time evidence. L2 is occurrence-based; refund/reversal appends material occurrences instead of overwriting prior release/settlement history.

R1/R2 are assessments, not resources, with states such as reconciled/divergent/pending/insufficient-evidence where applicable.

## 7.6 Economic Attribution

Persistent MPC Economics state with bounded local subject union only for Product 1.0 scopes, e.g. marketplace Sale, PostSaleResolution, Installation scope where material.

States include exact/partial/ambiguous/unresolved.

Partial monetary allocation uses Money and explicit unresolved remainder; no percentage-rounding shortcut or universal entity graph.

Resolve capability selects/interprets current attribution and carries typed current-resource `etag`; caller cannot rewrite authoritative movement facts.

## 7.7 CommercialPolicy

Closed typed Economics configuration containing only admitted policy meanings. Update is typed PATCH + true same-resource `If-Match`.

No generic `rules[]`, condition/action expression DSL or Governance-owned commercial threshold model.

## 7.8 EconomicPerformanceSummary

`GetEconomicPerformanceSummary` is a period/scope keyed Q with no synthetic Summary ID.

It preserves:

- explicit period and accepted Economics scope;
- coverage/partiality derived from underlying SaleEconomics universe/evidence;
- bounded exact monetary aggregates, rates and counts admitted by the concrete summary contract;
- evaluation/provenance.

The summary can never claim completeness/finality stronger than its underlying SaleEconomics coverage and never becomes a ledger/general analytics query language.

Baseline schema is closed; adding a metric requires a named Economics consumer/meaning, not a `map<string,number>`.

---

# 8. Governance and Sales

## 8.1 AuthorizationDecision

Immutable Governance-owned occurrence containing proportionately:

- Decision ID;
- typed concrete target Intent reference;
- exact target material revision/context;
- decision outcome;
- intended/authorized scope snapshots;
- authority/delegation context;
- server-attributed Principal;
- decision time.

No PATCH of historical Decision.

Target union includes only Product 1.0 Intents genuinely subject to Governance, e.g. ListingIntent, PriceIntent, BusinessOrderIntent, InvoicingIntent.

Create request carries exact target referenced `etag` adjacent to the typed reference. Governance authorization never executes or mutates the target Intent.

Baseline authorize scope equals the exact current intended scope; no generic partial-approval DSL without a real consumer.

## 8.2 AuthorizationDelegation

Stable opaque ID is justified by real update/revoke consumers while delegate/action/scope may change.

No generic Grant/IAM policy engine.

Update is typed PATCH + `If-Match`. Revoke is explicit monotonic capability; stale snapshots must not keep targeted standing authority alive and therefore no stale-snapshot block is required. Later re-grant is new explicit authority.

## 8.3 MarketplaceSale

Source-qualified external identity only: Marketplace Installation + native Sale key.

Read exposes bounded Sales interpretation/context, Selling Entity attribution, line/quantity facts required downstream, material post-sale facts and source provenance.

No Pack/Shipment/payment DTO topology absorbed into Sales.

## 8.4 `sale_line_key` lifetime

Each interpreted Sale line exposes an opaque Sale-scoped `sale_line_key` for durable Fulfillment/Post-Sale/Work selection/correlation without provider array position or global SaleLine entity.

> **Once minted within a Sale, a `sale_line_key` permanently denotes that line meaning. Reinterpretation may retire the key and mint a new one; it never rebinds an existing key.**

This closes read-then-act races while preserving local/non-global identity.

## 8.5 Selling Entity attribution

Owner-specific resolved/ambiguous/unresolved/unavailable meaning where reachable.

Resolve request selects one eligible same-Organization SellingEntity and carries current Sale-interpretation `etag` as typed custom-method request data.

Caller never supplies company/provider facts as replacement truth.

---

# 9. Business-System Materialization

## 9.1 BusinessOrderIntent

Owner-triggered MPC read/tracking resource; Product clients do not create it.

Carries proportionately:

- Intent ID;
- source-qualified Sale;
- target SourceInstance;
- PartyResolution;
- DestinationRealization;
- prerequisites;
- external-effect state/evidence;
- source-qualified native business-order result when established;
- convergence;
- Work/history references.

No TOP/NUNOTA/provider status as MPC ontology.

## 9.2 PartyResolution singleton

Contained meaning of one BusinessOrderIntent; no independent global PartyResolution ID is required.

States proportionately include unresolved/ambiguous/resolved-existing/resolved-new/unavailable.

Ambiguous candidates use opaque Materialization-local key + minimum PII-safe disambiguation evidence; raw `CODPARC` is not Product vocabulary.

Resolve request union:

```text
use_existing(candidate_key)
establish_new_from_current_sale_identity
```

Create-new variant accepts no arbitrary Customer master fields. D4 uses current authoritative Sale/source identity facts.

Operation requires mandatory `Idempotency-Key` and typed PartyResolution `etag`; candidate set is part of that representation, so one validator protects current resolution + candidates.

## 9.3 DestinationRealization

BusinessOrderIntent-contained singleton with honest states such as realized/external-required/unsupported/ambiguous/unknown/unavailable.

No Product Customer/Address/Contact CRUD. No write capability until D8 proves a safe selected business-system lane.

## 9.4 InvoicingIntent

Owner-triggered MPC read/tracking resource; Product clients do not create it.

Carries BusinessOrderIntent reference, current Fulfillment physical-readiness reference/evidence, prerequisites, external-effect state, source-qualified fiscal/native result(s), convergence/history.

No TOP 306/provider status as canonical Product identity/state.

## 9.5 Materialization axes

Keep prerequisites, external-effect state, native-result knowledge and convergence distinct. No global OperationState/giant status.

---

# 10. Fulfillment

## 10.1 FulfillmentExecution — one durable lifecycle identity

`FulfillmentExecutionId` is the **one** concrete durable Fulfillment lifecycle identity for the accepted physical checkpoint/history/artifact/Work/Materialization correlation meaning.

It is justified now by current admitted operations/history, not by hypothetical split fulfillment.

No parallel `FulfillmentIntentId`/Workflow ID may denote the same lifecycle. A future genuinely different authorizable routing/dispatch intent must prove distinct D2 meaning rather than duplicate identity by symmetry.

`FulfillmentState` is only a semantic pre-wire label from the admitted operation inventory; the wire home for List/Get Fulfillment state is `FulfillmentExecution`, not a second resource.

## 10.2 FulfillmentExecution scope

Explicit Sale-relative scope:

```text
whole_sale
sale_lines[{sale_line_key, quantity}]
```

Baseline implements no split-routing policy. The scope seam is sufficient for later real multiplicity without pre-building that feature.

## 10.3 FulfillmentExecution representation

Carries proportionately:

- ID;
- source-qualified Sale;
- scope;
- selected FulfillmentNode where established;
- separation checkpoint;
- physical-conference checkpoint;
- physical-readiness conclusion;
- provider-readiness/requirement closure;
- packing checkpoint;
- dispatch-handoff checkpoint;
- correlated Shipment refs;
- applicable provider deadline evidence;
- MPC internal operating-target evidence;
- Work/history refs.

No single writable status.

## 10.4 Physical checkpoints

Separation, physical conference, packing and dispatch handoff are explicit Fulfillment capabilities/occurrences, not status patches.

Custom request carries current FulfillmentExecution typed `etag`. Operations whose ratified safety tuple requires client idempotency also require `Idempotency-Key`.

Server attributes effective Principal/source/time. Client cannot submit `trusted_physical_evidence=true`.

PhysicalConference distinguishes checkpoint occurrence from readiness; discrepancy can be recorded without fabricating ready state.

## 10.5 Physical evidence client-class fence

Physical facts may be established only by admitted human Principal or an explicitly proven/provisioned system Principal/source capable of establishing that physical fact.

An ordinary automation Principal does not gain epistemic authority merely from `fulfillment.execute`.

## 10.6 FulfillmentNode

Closed MPC-owned resource distinct from InventorySource and external warehouse/company/location identity.

Baseline read/create/update shape is intentionally small:

- opaque `fulfillment_node_id`;
- organization-facing `display_name`;
- lifecycle (`active`/`inactive` proportionately);
- server-owned create/update provenance.

Create baseline requires only the bounded owner data needed to establish that Node (initially display identity). Typed PATCH updates only admitted fields. Deactivation is an explicit capability with typed current Node `etag`.

Do not add a generic capability graph, WMS topology, arbitrary warehouse attributes or provider/native location mapping without a concrete Fulfillment correctness need.

## 10.7 Fulfillment artifacts

Artifact keys are FulfillmentExecution-scoped opaque selectors, not global Asset IDs.

Descriptors preserve only needed kind/content type/source/provenance/sensitivity. Blob delivery/storage/signed URL mechanics remain D7.

## 10.8 Shipment

Source-qualified external resource; no synthetic MPC Shipment ID. Read exposes bounded Fulfillment interpretation, deadlines/outcomes/exceptions/freshness and provider enrichment needed by Product 1.0, not raw status/substatus/Pack ontology.

## 10.9 Fulfillment operating targets — bounded baseline

MPC internal targets remain distinct from provider deadlines.

Whole-W2 found only one current evidence-backed Product 1.0 target meaning:

```text
dispatch_handoff_lead_time_before_provider_deadline
```

Meaning: organization-owned internal safety lead time that targets a dispatch handoff **before** the externally authoritative provider dispatch deadline.

This field uses explicit inherit/override/effective provenance (§12). Provider deadline remains external evidence and cannot be edited/relaxed by this setting.

No generic `map<string,duration>`, SLA engine or speculative targets for separation/conference/packing are admitted. Add another target only when a concrete owner consumer/obligation proves it.

---

# 11. Post-Sale Resolution

## 11.1 Scoped coordination identity

PostSaleResolution is MPC-owned canonical identity with explicit Sale-relative scope:

```text
whole_sale
sale_lines[{sale_line_key, quantity}]
```

Product create establishes a coordination obligation, not immediate provider/ERP/refund effects.

## 11.2 Concerns are a set

Initiating concerns may coexist, e.g. cancellation/return/refund. Do not collapse them into one mutually exclusive type.

Concern set explains why Resolution exists; it does not command all consequence owners.

## 11.3 Independent consequence tracks

Resolution may track Sales cancellation, physical return/reverse logistics, refund/payment, business-system/fiscal and economic consequence evidence independently.

Each track carries typed owner/result references/current evidence, not provider `available_actions` vocabulary.

## 11.4 Closure

Resolution may expose open/closed lifecycle, but no direct Product close/status write exists. Post-Sale closes only when applicable consequence owners provide sufficient committed evidence.

---

# 12. Operational Work

## 12.1 Work is canonical obligation, not Task platform

Work carries proportionately:

- Work ID;
- typed originating condition;
- responsibility role;
- optional assignment;
- hold/resume lifecycle;
- escalation target/state;
- source-condition current/reconciliation state;
- closure provenance/history/time.

No direct arbitrary Work create baseline.

## 12.2 Work origin union

Work-local closed union with concrete Product 1.0 origins only, e.g. Readiness, Listing/Price, Sale, Materialization, Fulfillment, Shipment, Post-Sale, Economic Attribution, Authorization conditions.

No universal entity graph.

## 12.3 Responsibility role != AccessRole

Work `responsibility_role_key` is operational responsibility meaning, not AccessRole/Permission. Assignment never alters access authority.

## 12.4 Assignment / hold / resume / escalation

These remain explicit owner capabilities and carry current Work typed `etag`.

Escalation is declarative to an explicit target responsibility/state; increment/occurrence escalation is not baseline. If later changed to occurrence semantics, its safety tuple must be re-adjudicated for client idempotency.

## 12.5 No generic Work resolution/close

No `SubmitWorkResolution` or direct close. Source owner retains concrete closure truth; Work reconciles committed source resolution/coverage/supersession evidence.

## 12.6 Work closure-path audit — PASS

Current Product 1.0 Work-producing classes have legitimate owner closure paths through source auto-resolution/reconciliation or admitted owner-specific capabilities, including Readiness correspondence, Listing authoring, Offering/Availability ambiguity/convergence, Sales attribution, Party resolution, Governance Decision, Fulfillment checkpoints, Shipment evolution/Post-Sale, Economics attribution and Post-Sale consequence-owner results.

Reopen only if a concrete human-only evidence class cannot be carried by any owner capability and cannot auto-resolve; then add the smallest Work→source evidence-submission capability without transferring source truth to Work.

---

# 13. Identity/access and Marketplace Portfolio simple schemas

## 13.1 AccessContext

`GET /access-context` returns authenticated Principal's MPC context only:

- Principal ID/kind (`human`, `automation`, `system` as accepted semantics);
- visible Organization memberships;
- AccessRole keys per Organization;
- effective ordinary Permissions.

No raw OIDC/IdP claims, realm/client roles, MFA/session state or credentials as Product authority.

## 13.2 AccessRole

Product-defined AccessRoles use stable semantic `role_key`; no opaque ID solely for REST aesthetics. No custom-role CRUD baseline.

## 13.3 Membership and role assignment

Membership need not receive synthetic ID. Organization + Principal identifies membership relation; Principal + AccessRole is the structural role-assignment anchor.

Repeated assignment is structurally idempotent. Revocation remains monotonic/fail-safe and is not blocked by stale whole-role-set snapshots.

## 13.4 MarketplaceInstallation

MPC-owned canonical resource carrying bounded marketplace participation/lifecycle, Organization-facing configuration, eligible SellingEntity participation and external account/posture evidence only when materially required.

Marketplace kind is bounded supported Product vocabulary, not provider catalog API.

Credentials/OAuth codes/tokens/client secrets never enter Product schemas; D4 owns protocol auth.

## 13.5 SellingEntity

Organization-scoped MPC identity. May expose bounded display/legal/source-qualified references when operationally needed. CNPJ/native company/provider account identifiers never become SellingEntity identity by convenience.

Arbitrary SellingEntity create/edit remains deferred.

---

# 14. Cross-owner policy/configuration mechanics

Availability allocation/scope policy, CommercialPolicy and Fulfillment operating targets remain their respective owner authorities. There is no global `/policies`, `/rules`, `/settings` or generic configuration owner.

Where an accepted owner supports inheritance/override, write uses explicit variants:

```text
inherit
override(value)
```

Never `null` to mean clear/inherit.

Read exposes:

- configured owner value/mode;
- effective value;
- effective-source provenance such as default/inherited/explicit override.

Each owner uses only a closed scope union proven by current semantics. No generic `{entity_type, entity_id}`, tag selector, arbitrary hierarchy or expression language.

Externally governed provider deadlines/requirements remain external evidence and never become editable MPC policy.

---

# 15. Outcome / effect / convergence grammar

Do not collapse Product progress into one generic status.

Where applicable keep three independent axes:

1. owner business/capability outcome — e.g. accepted/rejected/pending/ambiguous;
2. external-effect meaning — e.g. not-attempted/accepted/rejected/pending/ambiguous under owner contract;
3. convergence — e.g. pending/converged/divergent/unknown/not-applicable.

Shared words may be reused only when meaning is actually the same. No universal OperationState/Result resource.

Generic `failed` is not baseline business meaning because it collapses HTTP failure, business rejection, provider rejection, ambiguous acceptance, divergence and unavailable evidence.

---

# 16. Response/status grammar

- successful Q/resource read → `200` + current representation;
- creation of a real canonical MPC resource/occurrence with standalone URI → `201 Created` + `Location` + created/current representation;
- typed PATCH → `200` + current updated representation and current ETag where applicable;
- custom `POST ...:verb` → `200` + current affected owner resource or operation-specific semantic result;
- stateless side-effect-free evaluation → `200` + evaluation.

Baseline avoids `204` where returning current representation/validator materially keeps clients synchronized.

`CreateListingIntentMedia` returns `200` + its operation-specific typed result (stable media descriptor + current parent ListingIntent validator, §3.9.1) because authored media deliberately has no standalone Product URI.

Failure statuses carry §19 Problem Details. A failure with no additional Product semantics beyond its standard status uses `about:blank` (§2.15.1) without waiving status-specific HTTP obligations (§2.15.2).

`202` is not used merely because an external effect has not converged.

---

# 17. Canonical revision-precondition grammar

W1 owns path/HTTP law; W2 owns request/result schemas and full admitted-operation application.

## 17.1 Carrier class A — true HTTP conditional

Use `If-Match` only when a standard method's literal request URI is the protected resource representation.

Missing required header → `428 precondition-required`; stale/false → `412 precondition-failed`.

Current admitted examples:

- PATCH ListingIntent;
- PATCH MarketplaceInstallation configuration;
- PATCH InventorySource;
- PATCH owner Availability allocation/scope policy resource;
- PATCH CommercialPolicy;
- PATCH AuthorizationDelegation;
- PATCH FulfillmentNode;
- PATCH FulfillmentOperatingTargets.

## 17.2 Carrier class B — owner custom method subject ETag

Custom owner method carries required typed `etag` request field/part when current state is required.

Missing/invalid → `422 validation-error`; stale → `409 resource-revision-conflict`.

Current admitted examples:

- Submit/Discard ListingIntent;
- CreateListingIntentMedia (multipart `etag` part);
- Deactivate MarketplaceInstallation;
- Deactivate InventorySource;
- Deactivate FulfillmentNode;
- ResolveEconomicAttribution;
- ResolveSaleSellingEntityAttribution;
- ResolveBusinessSystemPartyResolution;
- RecordSeparation / RecordPhysicalConference / RecordPacking / RecordDispatchHandoff;
- AssignWork / ClearWorkAssignment / HoldWork / ResumeWork / EscalateWork;
- Resolve/Clear ProductChannelCorrespondence using `correspondence_etag` from ProductChannelReadiness.

The request ETag is the same server-issued opaque validator; transporting it in body/part does not create a second version authority.

## 17.3 Carrier class C — exact referenced resource

A create/capability depending on another resource's exact revision carries that resource's ETag inside the typed reference.

Current admitted examples:

- CreateAuthorizationDecision target Intent;
- CreatePriceIntent superseding a PriceIntent.

Same `422` missing/invalid and `409 resource-revision-conflict` grammar applies.

## 17.4 No revision carrier by ratified design

No optimistic-revision carrier is added where the operation matrix explicitly says it is not material, including structural AccessRole assign, monotonic/fail-safe revocations, creates with no prior-resource revision axis and side-effect-free EvaluatePriceScenario.

Business/source validity is still evaluated; absence of revision proof does not waive semantic prerequisites.

---

# 18. Idempotency grammar

For operations requiring client `Idempotency-Key`, dedupe is resolved **before** re-evaluating stale revision proof on an exact retry of an already durably accepted intake.

Processing order proportionately:

1. AuthN;
2. decode/basic contract validity;
3. Organization Membership + Permission + client-class check;
4. require/validate Idempotency-Key — absent or invalid → `400 idempotency-key-required`;
5. derive semantic request fingerprint;
6. if key exists:
   - materially different request → `422 idempotency-key-reused`;
   - same request still processing → `409 idempotency-request-in-progress`;
   - same request already durably accepted/completed at intake boundary → resolve existing intake/result without rerunning stale revision proof or external effect;
7. only for new intake: evaluate revision proofs, owner/Governance prerequisites, then establish durable intake/effect.

Transport-level representation/size admission may precede step 1 (§3.9.3). It is not part of semantic evaluation and discloses no Product resource, tenant or business state.

Semantic fingerprint includes material:

- Organization;
- semantic operation/path target;
- material query parameters;
- semantic request body;
- all material subject/reference ETag proofs;
- binary content identity for multipart media;
- material multipart metadata.

It excludes bearer credentials, wall-clock request time, JSON whitespace/property order and Idempotency-Key itself.

Same key + changed ETag/precondition/file bytes = materially different request.

API failures before durable business intake do not themselves create consequential intake. Once durable intake exists, pending/ambiguous/rejected durable owner state remains tied to that intake.

Idempotency never authorizes blind provider/business-system redispatch after ambiguous possible acceptance.

---

# 19. Problem Details catalog

Canonical small problem families:

- `malformed-request`;
- `validation-error`;
- `authentication-required`;
- `access-denied`;
- `resource-not-found`;
- `idempotency-key-required`;
- `idempotency-key-reused`;
- `idempotency-request-in-progress`;
- `precondition-required`;
- `precondition-failed`;
- `resource-revision-conflict`;
- `resource-state-conflict`;
- `invalid-cursor`;
- `cursor-expired`;
- `internal-error`.

Problem `type` is an absolute stable URI; exact host is later OpenAPI/serving topology.

Validation errors may include bounded `errors[]` with machine-readable pointer/location + detail.

`resource-revision-conflict` may identify which supplied typed validator was stale via bounded pointer/location extension; do not create separate subject-vs-reference conflict types with identical semantics.

`resource-state-conflict` means the server evaluated current authoritative resource state and that state does not admit the requested capability:

```text
resource-revision-conflict → the supplied typed validator is stale
resource-state-conflict    → the supplied/current view may be current, but the capability is no longer admissible
```

Apply `resource-state-conflict` only after any operation-specific exact-repeat or structural-idempotency rule has resolved a harmless repeat. It closes the admitted custom-capability population without introducing one state-specific type per owner.

Canonical status mapping for these families:

```text
malformed-request                  400
idempotency-key-required           400
invalid-cursor                     400
cursor-expired                     400
authentication-required            401
access-denied                      403
resource-not-found                 404
resource-revision-conflict         409
resource-state-conflict            409
idempotency-request-in-progress    409
precondition-failed                412
validation-error                   422
idempotency-key-reused             422
precondition-required              428
internal-error                     500
```

A failure with no additional Product semantics beyond its standard status uses `about:blank` (§2.15.1) instead of a new type.

Product and technical-protocol problem registries stay separate:

```text
Product route
→ Product OpenAPI + this W2 Product Problem contract

provider / OAuth / technical presentation route
→ protocol-local executable contract
→ never W2 Product type authority
→ never Product SDK operation/error surface
```

A technical route may reuse `application/problem+json` as a representation format. Format reuse does not merge type registries, operation identity, SDK exposure or business meaning. Technical-protocol problem types use a namespace disjoint from the Product problem-type namespace; examples such as `oauth-state-expired`, `provider-code-invalid`, `seller-mismatch` and `provider-origin-invalid` remain forbidden from the Product catalog. Product UI learns durable current posture through accepted Product reads, not callback navigation/error vocabulary.

Cursor-specific applicability is owned by canonical W3:

- `invalid-cursor` is HTTP `400` for a supplied cursor that cannot validly continue a **well-formed** collection request, including integrity/unknown/operation/Organization/semantic-query mismatch after normal access/privacy checks;
- `cursor-expired` is HTTP `400` for a legitimately issued continuation that can no longer be resumed honestly;
- a missing/invalid operation-required subject/search/filter parameter remains ordinary `422 validation-error` before cursor evaluation;
- no `cursor-stale`, `cursor-gone`, `cursor-conflict` or provider-specific cursor taxonomy is introduced.

Raw provider/business-system errors, payloads, secrets and arbitrary PII never become Product problem truth.

---

# 20. Cross-Organization privacy

Organization-scoped body/query reference invalid in the path Organization fails closed without disclosing another tenant's ownership/existence.

A path resource not visible/resolvable in requested Organization may return ordinary 404 semantics rather than prove cross-tenant existence.

---

# 21. Business unavailability is not server failure

Examples that remain successful Product semantic representations rather than generic 5xx/4xx solely by meaning:

- provider capability unsupported/external-required;
- unknown/unavailable/partial source/market evidence;
- approval-required/Governance pending or business rejection;
- ambiguous possible provider acceptance;
- Economics insufficient evidence;
- divergent/unknown convergence.

Use HTTP failure only when the Product API request itself cannot be satisfied under its transport/access/contract/server semantics.

---

# 22. D4/D8 proof obligations carried forward

These are proof notes, not W2 owner changes:

1. authoritative reread for N/A-bearing provider listing representation must prove the selected provider lane can distinguish explicit N/A from absent/unobserved; current classic-item behavior requiring internal attributes is evidence, not a universal selected-lane command;
2. selected User-Product publication lane must prove its concrete conditional-requirement validation surface before D8 claims dispatch/convergence based on the draft-dependent seam;
3. no provider write is considered converged without the D4-required authoritative reread/effect reconciliation.

No parent semantic reopen is required now. Reopen only the implicated D1/D3/D4-R1 decision if real selected-lane proof shows the owner-preserving technical composition cannot represent the provider contract.

---

# 23. Canonical Whole-W2 negative controls

Later OpenAPI/contract proof must make at least these defects invalid/unreachable:

1. physical UUID/ULID promise for opaque MPC ID without accepted need;
2. native external ID serialized as arithmetic number;
3. universal ExternalRef/entity graph;
4. Money as authoritative JSON float or global minor-unit/rounding law;
5. unknown/unavailable/partial collapsed into null/zero/false/empty;
6. raw provider property bag/DTO in Product request/response;
7. one request/read schema making server authority/history writable;
8. generic Result/Operation/Evidence/Fact/Subject/Scope/Policy/Workflow wrappers becoming business ontology;
9. price or Availability quantity entering ListingIntent desired content;
10. FOLLOW_SOURCE hidden fallback;
11. `not_applicable` inferred from missing source evidence;
12. number+unit represented as unvalidated provider string or generic UoM engine;
13. client PATCHing Price/Availability correlation into ListingIntent;
14. current source/provider state rewriting historical dispatch basis;
15. `PublicationAttempt` standalone Product CRUD introduced merely for history;
16. source Product search receiving MPC Product identity;
17. completed provider market enumeration presented as universal market completeness;
18. Economics known profitability from missing required components;
19. client scenario request spoofing authoritative cost/tax/fee/market evidence;
20. L2 refund/reversal overwriting prior occurrence history;
21. generic Reconciliation/finance ledger resource;
22. AuthorizationDecision mutation or approval executing the target;
23. Party resolution accepting arbitrary Customer master fields;
24. TOP/NUNOTA/CODPARC/provider status becoming canonical MPC semantics;
25. ordinary automation fabricating physical evidence;
26. FulfillmentState and FulfillmentExecution becoming duplicate resources;
27. parallel FulfillmentIntentId for the same lifecycle;
28. FulfillmentNode collapsing into InventorySource/native warehouse;
29. generic SLA/target map beyond the one proven dispatch lead-time field;
30. `sale_line_key` rebound to a different line meaning;
31. Post-Sale direct close/provider action vocabulary fabricating consequence closure;
32. Work responsibility encoded as AccessRole or Work close fabricating source truth;
33. base-resource `If-Match` used on a distinct custom `:verb` URI by private convention;
34. custom conditional header added when typed ETag request data suffices;
35. same resource revision protected by two independent version authorities;
36. correspondence given synthetic ID/forced PUT/DELETE solely for concurrency;
37. exact idempotent retry failing only because first call advanced revision;
38. same idempotency key with changed ETag/file bytes treated as same request;
39. idempotency authorizing blind redispatch after ambiguous effect;
40. business unknown/rejection/unavailable turned into generic Problem Details by convenience;
41. cross-Organization validation leaking another tenant's existence;
42. cursor errors maintained in a second collection-only problem catalog instead of this canonical catalog;
43. a `405` response omits or lies in `Allow`;
44. wrong top-level multipart representation classified as `422` rather than `415`;
45. unsupported, undecodable or inspected-mismatching binary format accepted, or classified as semantic metadata validation;
46. excess-size body fully buffered before the transport guard can refuse it;
47. malformed multipart reported as a provider/storage failure;
48. missing/duplicate file part or invalid typed ETag accepted;
49. absent `Idempotency-Key` creating intake or receiving an undeclared status;
50. a current but lifecycle-inadmissible capability confused with stale revision, or silently succeeding;
51. exact lost-response replay creating a second media ID;
52. successful media creation failing to advance and return the current parent ListingIntent validator;
53. two concurrent creates carrying one validator both succeeding as independent current-state writes;
54. `listing_intent_media_id` from ListingIntent A accepted in ListingIntent B;
55. accepted bytes replaced under an existing media ID;
56. unselection destroying material historical evidence;
57. a presentation access reference reaching durable history, idempotency fingerprint, logs or Problem Details;
58. stable media ID alone retrieving bytes;
59. a caller without current Organization Membership or `offering.read` retrieving authored bytes;
60. anonymous, durable or freely forwardable object-store locator becoming baseline access;
61. source locator and authored access reference collapsed into one authority/type;
62. raw storage/scanner/provider error appearing in Product Problem Details;
63. technical OAuth/provider/delivery problem types using the Product problem namespace or appearing in Product OpenAPI/SDK;
64. a status-only failure minting an MPC problem type that merely repeats standard HTTP meaning;
65. Product clients receiving duplicate global `type` and `code` discriminators;
66. D7 privately creating a media Product GET, Delete or Update operation;
67. stale active “next action” text misrouting a fresh session away from the router.

A green artifact that did not execute the protected subject is no proof.

---

# 24. W2 family coverage

Canonical W2 now covers every admitted Product schema family:

- Identity/Access — §13;
- Marketplace Portfolio — §13;
- Readiness — §5;
- Offering / MarketplaceListing / Price — §3;
- Availability / InventorySource — §4;
- Market Intelligence — §6;
- Commercial Economics / Performance / Attribution / Policy — §7;
- Governance / Sales — §8;
- Business-System Materialization — §9;
- Fulfillment / Node / Execution / Artifacts / Shipment / operating target — §10;
- Post-Sale — §11;
- Operational Work — §12;
- cross-owner configuration/outcome/revision/idempotency/problems — §§14–21.

No admitted Product family is left to a provider DTO, legacy OpenAPI, frontend form or deleted staging artifact.

---

# 25. Whole-W2 method disposition

```text
D0→D5-B1 / ratified B2 semantic authority     CURRENT STRUCTURE CONFIRMED
Wire W1                                        ACCEPTED / CANONICAL
W2 schema grammar                              ACCEPTED / CANONICAL
Whole-W2 lead + Fable R1/R2 + GPT adjudication CONVERGED / OPERATOR-RATIFIED
Whole-W3 cursor Problem Details amendment       INCORPORATED / OPERATOR-RATIFIED
Final Problem/media consistency package        INCORPORATED / OPERATOR-RATIFIED
Product operations / ordinary Permissions      95 / 29 unchanged
Parent-stage reopen                            NONE
```

**Global Maximum:** strong reusable wire/value/safety mechanics + owner-specific semantic schemas; complete Product 1.0 reachability without Product/PIM, provider mirror, Rules/UoM/SLA/Workflow/Finance/Task platform or generic operation ontology.

W2 remains closed/canonical. Collection/query applicability for cursor problems is owned by canonical W3; exact OpenAPI type URI host/tooling remains later Wire work. Current status and exact next action are owned only by the rebaseline router.

Implementation remains blocked until D9.
