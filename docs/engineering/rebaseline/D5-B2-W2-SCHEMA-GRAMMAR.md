# D5-B2 — W2 Request / Response Schema Grammar

> **Status:** ACCEPTED / CURRENT CONSOLIDATED AUTHORITY  
> **Parent:** `D5-B2-WIRE-CONTRACT.md`  
> **Exact machine schemas:** `contracts/api/product/openapi.yaml` + local Product components

## 1. Governing invariant

> **A Product schema preserves every material distinction the accepted owner can make, while adding no independent business meaning merely for uniform serialization. Reusable wire primitives normalize representation; business meaning, knowledge, lifecycle, evidence and outcomes remain owner-specific.**

W2 owns semantic request/response shape laws. Exact current OpenAPI schema spelling is the machine-readable projection of these laws.

## 2. Core wire grammar

### 2.1 Identity

- MPC-owned IDs are opaque non-empty strings; no UUID/ULID/database encoding is promised.
- Provider/business-system native keys are also opaque strings even when numerically shaped.
- Identifiers are never quantities.
- Use typed source-qualified refs, not universal `ExternalRef` or `{entity_type,entity_id}` graphs.
- Organization is supplied by path scope and is not redundantly embedded in every Organization-owned ref unless the semantic object itself requires it.

### 2.2 Exact numeric values

`ExactDecimalString` carries authoritative decimal numbers. `Money` is exact amount + explicit currency. Binary floating point, exponent notation and global `round(2)` are not authoritative business rules.

### 2.3 Time

Names preserve semantic time: source/effective, observed/acquired, recorded, created/updated, decided/authored, deadline etc. Generic response generation time never substitutes for owner/source time.

### 2.4 Request != read

Client-authored requests and server reads are distinct when authority differs.

Client writes never author:

- server identity;
- actor/history attribution;
- provider outcome/convergence;
- source observation provenance;
- current source/read projection;
- Governance result/history;
- presentation snapshots/locators that are not decision authority.

This rule is semantic, not merely an OpenAPI `readOnly` convention.

#### Human-operable read projection grammar

A canonical Ref/request carrier remains minimal. A current read may carry an adjacent owner-correct presentation projection when a proven human job requires recognition, selection or explanation. A purpose/historical presentation snapshot is a third meaning and is not refreshed as current source truth.

Presentation is never accepted as identity or write authority. Equal labels never collapse distinct keys. Unknown/unavailable presentation never erases the known canonical subject.

The current implementation vocabulary is intentionally bounded to these owner-specific schema families:

```text
SourceProductPresentation
PublicationContextRef / PublicationContextView
PublicationOptionDescriptor
PublicationUnitDescriptor
PublicationValue / PublicationValueView
PublicationSourceCandidateView
CorrespondenceCandidatePopulation
MarketplaceListingPresentation
ListingIntent requirement-resolution read views
source/authored/observed media presentation families
```

The MarketplaceListing/ListingIntent OAD drift directly implicated by the approved D6-R2 prerequisite is repaired under W2 without changing W1 operation/resource authority or W4 authorization semantics.

### 2.5 Closed semantic objects / bounded unions

Request objects are closed by default; responses are closed where needed to protect meaning. Material unions use exclusive typed/discriminated variants. Provider enrichment is typed and bounded; arbitrary `metadata`, raw provider DTOs or dynamic field bags are forbidden.

### 2.6 `null` and knowledge

`null` is not shorthand for unknown/unavailable/partial/absent/not-applicable. Where knowledge is material, use the smallest owner-specific state union needed by reachable states.

```text
known value
known empty/absent
unknown / insufficiently known
unavailable
partial / incomplete coverage
unsupported / not-applicable
```

Only reachable variants belong to each owner schema. Known zero/false/empty remains distinct from unknown.

### 2.7 Provenance/freshness locality

Attach observation/provenance/freshness at the smallest semantic unit for which it is uniformly true. Do not create universal top-level Evidence/metadata/generated-at envelopes.

### 2.8 No universal Result/Operation envelope

Do not wrap all operations in generic `Result<T>`, `Operation`, `CommandResult` or async tracking objects. Durable owner Intents/Resources are their own tracking identities.

---

## 3. Problem / mutation grammar

### 3.1 Product outcome != HTTP Problem

Pending, approval-required, insufficient evidence, ambiguous external effect, divergent convergence, unsupported/external-required and similar owner meanings remain valid Product semantics where admitted. They are not automatically transport errors.

Authentication/access/privacy, validation, invalid cursor, conflict, stale precondition, unsupported media/content size and internal failure use stable Product Problem semantics. Provider error DTOs remain D4-local evidence.

### 3.2 Typed updates

Use operation-specific typed JSON update bodies, not generic JSON Patch/Merge Patch. Mutation schemas include only fields the caller is actually authorized to author.

### 3.3 Revision/ETag

Strong owner revisions/ETags are typed. `If-Match` is used when ordinary resource stale-write prevention is material, with explicit stale/missing-precondition semantics.

Current deliberate Governance exception:

```text
CreateAuthorizationDecisionRequest
  etag: StrongETag        // exact AuthorizationRequest revision
  outcome: authorize | reject
```

The custom request-scoped decision operation carries this ETag in its body, not target state and not a target ETag. It also uses caller Idempotency-Key at W1/operation level.

---

## 4. Product references / target unions

There is no platform-wide generic entity graph. Use owner-specific typed refs/unions, including proportionately:

```text
SourceProductRef
MarketplaceListingRef
MarketplaceSaleRef
MarketplaceShipmentRef
ListingIntentRef
PriceIntentRef
BusinessOrderIntentRef
InvoicingIntentRef
FulfillmentExecutionRef
PostSaleResolutionRef
WorkRef
AuthorizationRequestRef
AuthorizationDecisionRef
AuthorizationTargetRef
AvailabilityTargetRef
EconomicAttributionRef
```

`AuthorizationTargetRef` remains a closed union over:

```text
listing_intent
price_intent
business_order_intent
invoicing_intent
```

Current Request/actionable-read wire must not reintroduce generic target/entity payloads.

---

## 5. Publication / Readiness schema grammar

The accepted publication model remains key/identity based and source/provider qualified. Current human reads add only the owner-correct presentation required by the grammar above; requests and effects remain canonical-key based.

### 5.1 Publication context / requirements

Readiness current publication-requirements meaning preserves proportionately:

```text
exact source-qualified subject
Marketplace Installation
category_key? / product_type_key? context
requirements_revision
requirements[]
source_media_candidates[]
observation/evaluation provenance
```

`PublicationRequirement` preserves:

- opaque `requirement_key`;
- class/applicability;
- value specification/constraints;
- whether not-applicable is allowed;
- source evidence.

Allowed option/unit/candidate keys are opaque requirement/context-scoped decision keys. W2 does **not** promote them into global attribute/unit identity or a generic provider-rules DSL.

### 5.2 Source evidence

Source evidence preserves distinct states such as known, conflicting, missing, unknown, unavailable and unsupported, with candidate/value evidence only where semantically known. `FOLLOW_SOURCE` and `EXPLICIT_OVERRIDE` remain the two accepted ListingIntent resolution modes.

### 5.3 PublicationValue

`PublicationValue` is a closed authored/provider-semantic union containing the currently accepted value kinds, including text, exact-decimal, boolean, option, text-list, option-list, number+unit and explicit not-applicable. Client-authored value identity remains key/value based; provider presentation enrichment proposed by PR #70 is not silently added here.

---

## 6. Offering / MarketplaceListing grammar

### 6.1 Canonical Listing identity

`MarketplaceListingRef = MarketplaceInstallation + native_listing_key`; do not add a synthetic MPC listing identity merely as a mirror alias.

### 6.2 Current actual-state meaning

Canonical W2 MarketplaceListing meaning includes proportionately:

- qualified Listing ref;
- normalized lifecycle;
- publication context;
- bounded observed representation values/media;
- current observed marketplace price when material;
- observation/freshness/provenance.

The D6-R2 prerequisite repairs the directly implicated OAD underprojection without adding another Listing identity, owner or operation.

### 6.3 ListingIntent

ListingIntent is Offering-owned desired state with stable identity/lifecycle/revision and exact source/target context. Its write schemas carry only desired authoring decisions, including requirement resolutions/media selections allowed by D4-R1.

Read schemas preserve current authored state, owner lifecycle/dispatchability/convergence/provenance sufficiently for accepted consumers. They do not make Readiness/provider evidence Offering-owned.

### 6.4 Media trust separation

Accepted W2 distinguishes:

```text
ListingIntentMediaDescriptor
  stable authored-media identity + bounded content/provenance facts

ListingIntentMediaPresentationDescriptor
  stable descriptor + volatile authorized access reference
  response/presentation only
```

The access reference is never persisted into durable history/idempotency/logs/Problems as authority. Source-media locators, authored-media access references and observed-provider media are distinct trust types; there is no generic Media/Asset business service.

The D6-R2 prerequisite restores this accepted response-only presentation descriptor while preserving the trust separation above.

---

## 7. Availability / Market / Economics / Performance read grammar

These owners preserve their own knowledge/evidence rather than one universal analytics schema.

- Availability distinguishes pre-creation ListingIntent target vs existing source-qualified Listing target and honest known/unknown/unavailable quantity/control meaning.
- Market uses typed MarketSubject/comparable-offer evidence, comparability/provenance and explicit coverage where material.
- Economics keeps Expected, Sale/Realized and Attribution meanings distinct; price scenario evaluation is stateless read-like evaluation even when encoded as POST.
- Performance uses exact Installation + explicit period/comparison semantics and provider measurement/coverage/comparability evidence. No arbitrary metric/dimension/group-by query/payload grammar exists.

When a read list item needs only a semantic subset of the point owner meaning, W3 governs collection projection; it must not invent list-only business conclusions.

---

## 8. Personal Notifications schema grammar

### 8.1 Notification

Current `Notification` is a closed kind-constrained union consistent with D2's fourteen kinds. Common owner meaning includes exact Organization/historical recipient, typed source ref, source occurrence lineage/times, retained `subject_display_label`, personal awareness state/revision and only kind-specific typed result/replacement atoms that are actually admitted.

No generic fields such as `payload`, `metadata`, `data`, `summary`, `reason`, arbitrary attributes or template variables are admitted.

`subject_display_label` is an immutable purpose snapshot, not current source identity/truth or navigation authority.

### 8.2 Source refs

Kind/source compatibility remains exact. Important current corrections:

- F13 `AUTHORIZATION_ACTION_REQUIRED` source = `AuthorizationRequestRef`;
- F14 `AUTHORIZATION_DECISION_RESULT` source = `AuthorizationTargetRef`, not AuthorizationDecisionRef.

This lets actionable humans continue to the Request and lets requesters continue to the governed target without broad Governance history access.

### 8.3 Typed result atoms

Only current admitted result atoms exist:

```text
F02 offering_async_result_outcome
  converged | rejected | ambiguous | divergent

F14 authorization_decision_outcome
  authorize | reject
```

No generic `status/result/reason` property is added.

### 8.4 Recipient candidate

Routing-recipient candidate projection is deliberately minimal:

```text
principal_id
display_name
```

No role keys, Permission internals, email/claims or generic person profile are exposed merely for routing selection.

### 8.5 Notification routing state

Routing schemas preserve typed configured/unconfigured state/revision semantics from D2. Current Product API sets one explicit route state for an ORG_ROUTED kind; it never treats empty recipient list as hidden default/all-admin routing.

---

## 9. AuthorizationRequest / Governance schema grammar

### 9.1 AuthorizationReviewBasis

`AuthorizationReviewBasis` is a closed four-way typed immutable union aligned to `AuthorizationTargetRef`:

```text
listing_intent
price_intent
business_order_intent
invoicing_intent
```

Each variant contains only the bounded decision-time owner evidence required to understand what is being authorized. Generic `payload`, `metadata`, `data`, `attributes`, `entity_type/id`, raw provider fields or current-screen DTOs are forbidden.

### 9.2 ActionableAuthorizationRequestView

Current exact-human point view is a closed union whose branch requires:

```text
authorization_request_id
target
subject_display_label
review_basis
created_at
```

Target kind and review-basis kind must match. This is a **currently actionable decision view**, not general Request history.

### 9.3 ActionableAuthorizationRequestListItem

Collection item carries only the scan/select subset:

```text
authorization_request_id
target
subject_display_label
created_at
```

No generic search/status/assignee/role payload is admitted.

### 9.4 AuthorizationDecision

Decision remains historical Governance state with:

```text
authorization_decision_id
authorization_request_id
target
review_basis
outcome
decided_by_principal_id
decided_at
```

Decision target is no-ETag `AuthorizationTargetRef`; current request revision belongs to the decision command carrier, not the historical target reference.

### 9.5 Work origin

`WorkOrigin` admits bounded `authorization_request` origin for zero-current-decider/recovery cases while preserving other accepted owner-specific origins. Work never becomes generic approval authority.

---

## 10. Collection envelopes / cursor hooks

W2 defines item/schema meaning; W3 defines collection population/query/order/cursor semantics. Named collection responses use owner-specific array property + optional `next_cursor` and any legitimate owner coverage/provenance—never a universal Page/metadata wrapper.

---

## 11. Media upload technical contract

ListingIntent authored media creation uses one bounded multipart request with exact file + current parent ETag/revision carrier and returns stable media descriptor/current parent validator. MIME/size/content validation is Product/technical contract, not provider media authority.

Actual byte delivery is a separately justified authenticated technical presentation surface under the already-authorized ListingIntent read context; it is not a new Product business operation.

---

## 12. Explicit non-goals

W2 does not create:

- universal entity/reference/evidence/result/page wrappers;
- provider DTO passthrough or dynamic provider-field maps;
- generic attribute/PIM/rules/unit engine;
- generic Media/Asset business service;
- generic authorization/workflow payload;
- schema fields solely because one frontend component would be easier to code;
- presentation fields on canonical refs or writes, or owner-wide enrichment without a proven human job.

## 13. Reopen trigger

Reopen W2 only when an accepted owner/consumer requires a material distinction the current schema cannot express honestly, or when downstream frontend/runtime evidence falsifies an existing schema law. Reopen the smallest owner/schema family; do not patch with generic metadata or a second DTO authority.
