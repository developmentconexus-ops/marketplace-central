# D5-B2 — W2-E Transversal / Final Schema Consistency

> **Status:** ACCEPTED IN-STAGE  
> **Parent W2:** `D5-B2-W2-SCHEMA-GRAMMAR.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + ratified D5-B2 Whole-Matrix + Wire W1 + W2-A/B/C/D  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Accepted:** 2026-08-18

## 1. Governing invariant

> **Cross-cutting wire mechanics may be shared only when the protected meaning is genuinely the same. Configuration, safety and HTTP problem mechanics use one consistent grammar; business policy, lifecycle, evidence and outcome meaning remain owner-local.**

W2-E closes cross-owner consistency for configuration, simple Access/Portfolio resources, media intake transport, outcome/status relationships, direct versus referenced-resource preconditions, idempotency replay behavior and Problem Details.

It does not create GenericPolicy, GenericResult, GenericOperation, GenericEvidence, GenericSubject, GenericWorkflow or a platform-wide resource graph.

## 2. Owner-local policy/configuration with shared inheritance mechanics

Availability allocation/scope policy, Commercial Economics policy and Fulfillment operating targets remain owned by their respective D1 authorities.

No global `/policies`, `/rules`, `/settings` or generic configuration authority is introduced.

Where an accepted owner supports inheritance/override, the wire uses explicit semantic variants rather than `null` clearing:

```text
inherit
override(value)
```

Reads expose the owner-specific configured value plus effective value and material provenance such as default/inherited/explicit override. External provider deadlines/requirements remain source-qualified evidence and never become editable MPC policy or inherit/override state.

Each owner uses a closed scope union only for scopes actually proven by Product 1.0. No generic `{entity_type, entity_id}`, expression, tag selector or arbitrary policy-scope query is introduced.

## 3. Identity/access simple resource grammar

### 3.1 AccessContext

`GET /access-context` returns only the authenticated Principal's MPC access context, proportionately:

- MPC Principal ID and kind (`human | automation | system` as accepted semantics);
- current Organization memberships visible to that Principal;
- product-defined AccessRole keys currently assigned in each Organization;
- effective ordinary MPC Permissions.

OIDC token claims, IdP realm/client roles, MFA/session state and credentials are not Product authority fields.

### 3.2 AccessRole

Product-defined AccessRoles use stable semantic `role_key` values rather than minting opaque canonical IDs solely for API aesthetics. Product 1.0 admits no custom-role CRUD.

### 3.3 Membership and role assignment

Membership need not receive a synthetic canonical ID merely to normalize REST. Organization + Principal identifies the relevant membership relation; one Principal + AccessRole relation is the structural assignment anchor.

Repeated identical assignment is structurally idempotent. Revocation remains monotonic/fail-safe and is not blocked merely by a stale whole-role-set snapshot.

## 4. Marketplace Portfolio simple resource grammar

`MarketplaceInstallation` remains an MPC-owned canonical resource with opaque ID and proportionately carries marketplace participation/lifecycle, organization-facing configuration, eligible Selling Entity participation and bounded external account/posture evidence where materially needed.

Marketplace kind is a bounded Product vocabulary for supported marketplaces; it is not a generic provider catalog API.

Marketplace credentials, OAuth code/tokens, refresh tokens, provider client secrets and business-system credentials never enter Product Installation schemas.

Selling Entity remains an Organization-scoped MPC identity with source-qualified legal/business-system references only when operationally required. Native CNPJ/company/provider account identifiers never become Selling Entity identity by convenience. Arbitrary Selling Entity create/edit remains deferred.

## 5. ListingIntent media intake transport

The Product wire for `CreateListingIntentMedia` is closed in D5 independently of D7 storage realization.

Baseline transport:

```http
POST /organizations/{organization_id}/listing-intents/{listing_intent_id}/media
Idempotency-Key: <required>
If-Match: <current ListingIntent ETag>
Content-Type: multipart/form-data
```

The request contains one binary file part plus only bounded semantic metadata if later required by the exact media contract. Baseline does not use base64-in-JSON, arbitrary external URL, UploadSession, presigned-upload resource, ProductAsset or generic Asset API.

The response is a ListingIntent-scoped media descriptor with opaque `listing_intent_media_id`, content/media type and server-owned actor/time metadata needed for selection/history. Storage key, bucket, public URL, CDN URL, hash/resizing pipeline and object-store topology remain D7 mechanism.

ListingIntent remains authority for media selection/order/role. Array order is the one publication-order authority.

## 6. Outcome axes remain distinct

Product resources/capabilities do not collapse progress into one generic status.

Where materially reachable, keep distinct:

1. **owner capability/business outcome** such as accepted/rejected/pending/ambiguous;
2. **external-effect meaning** such as not-attempted/accepted/rejected/pending/ambiguous as defined by the owner/effect contract;
3. **convergence meaning** such as pending/converged/divergent/unknown/not-applicable as defined by that owner.

Shared words may be reused when their semantics are actually identical; no universal `OperationState`, `Result<T>` or business outcome resource is introduced.

Generic `failed` is not a baseline business state because it would collapse API failure, business rejection, external rejection, ambiguous acceptance, divergence and unavailable evidence.

## 7. Response/status baseline

Use HTTP status to report HTTP/API success/failure, never to substitute for owner outcome.

- successful owner Q/resource read → `200` + current representation;
- first successful creation of a canonical MPC resource/occurrence → `201 Created` + `Location` + created/current owner resource;
- typed PATCH of an existing resource → `200` + current updated resource and new `ETag` when concurrency applies;
- `POST {resource}:verb` → `200` + current affected owner resource or the operation-specific semantic result;
- stateless side-effect-free capability → `200` + operation-specific evaluation.

Baseline does not use `204` for Product mutations where returning the current representation/ETag materially reduces client re-read and keeps React/agent clients synchronized.

`202 Accepted` is not a synonym for an external effect that has not converged. Durable owner Intents already provide tracking identity; the resource may return with semantic pending/ambiguous/convergence state under ordinary successful HTTP status.

## 8. Direct-resource concurrency versus referenced-resource preconditions

### 8.1 Direct selected request resource

When the resource named by the request URI is the stale-state authority being protected, use W1 strong opaque `ETag` + `If-Match`.

- required `If-Match` missing → `428 Precondition Required` + Problem Details;
- supplied `If-Match` stale → `412 Precondition Failed` + Problem Details.

### 8.2 Referenced resource whose exact revision is material

A create/capability may depend on the exact current revision of another referenced MPC resource that is not the HTTP request target, for example an AuthorizationDecision targeting an Intent or a PriceIntent explicitly superseding another PriceIntent.

Use a small technical `ReferencedResourcePrecondition` carrying that referenced resource's opaque strong `ETag` inside the typed reference that must be frozen.

This is not a second version authority; it carries the same opaque validator emitted by the referenced resource.

Missing required referenced precondition is a schema/validation problem. A referenced validator that is no longer current is a semantic request conflict requiring re-read/re-decision, not an HTTP `If-Match` failure against the selected request URI.

Baseline disposition:

- missing required referenced precondition → `422` validation problem;
- stale referenced resource precondition → `409 Conflict` Problem Details with stable referenced-resource-conflict type.

Do not misuse `If-Match` to condition a different resource than the selected request target merely for stylistic uniformity.

## 9. Idempotency replay semantics and evaluation order

For operations whose ratified safety tuple requires a client `Idempotency-Key`, idempotency deduplication must be resolved before re-evaluating stale direct-resource preconditions on an exact retry of an already accepted intake.

Otherwise a lost-response retry could incorrectly turn into `412` solely because the first successful request already changed the resource ETag.

Semantic processing order, proportionately:

1. AuthN;
2. request decoding/basic contract validity;
3. Organization Membership + Permission + admitted client-class check;
4. require/validate `Idempotency-Key`;
5. derive semantic request fingerprint;
6. if key already exists:
   - materially different request → explicit idempotency reuse problem;
   - same request still processing → explicit in-progress problem;
   - same request already durably accepted/completed at intake boundary → resolve the existing intake/result without re-running stale `If-Match` or redispatching external effect;
7. only for a new key/intake: evaluate direct `If-Match`, referenced-resource preconditions, owner business prerequisites/Governance and then create/accept durable owner intake/effect.

Idempotency remains request-intake safety only. It never authorizes blind replay of an ambiguous provider/business-system effect.

## 10. Semantic request fingerprint

Idempotency equivalence is semantic, not raw-body-byte equality.

The fingerprint materially includes the Organization, semantic operation, selected path/resource target, material query parameters, semantic request body, direct `If-Match` value when applicable and referenced-resource preconditions.

It excludes bearer/client credentials, request wall-clock time, JSON whitespace/property order and the Idempotency-Key itself.

Changing a material precondition while reusing the same key is a different semantic request and fails explicitly rather than mutating the original intake.

D7 chooses hash/persistence/retention/locking mechanics; W2 owns the semantic equivalence property.

## 11. Idempotency Problem Details dispositions

MPC defines the header semantics explicitly; it does not depend on an unpublished/expired external draft becoming authority.

Baseline API problem cases:

- mandatory key absent → `400` idempotency-key-required;
- same key reused for a materially different semantic request → `422` idempotency-key-reused;
- same key + same request while first intake is still being established/processed → `409` idempotency-request-in-progress;
- exact replay after durable intake → resolve the existing intake/resource/result; never create another semantic effect.

API-level failures before durable business intake (for example AuthN/access/schema/precondition failure) do not themselves create the consequential business intake. Once durable intake exists, later owner pending/ambiguous/rejected-as-durable-decision state remains associated with that intake/key and exact retry resolves it rather than minting a second owner object/effect.

For a created resource, exact replay resolves the same canonical identity. Baseline may return `201 Created` + the same `Location`/current canonical resource rather than introducing a separate replay-specific business status. Byte-for-byte historical response replay is not required unless later evidence proves it material.

## 12. RFC 9457 Problem Details baseline catalog

Problem `type` remains the primary stable machine identifier. Baseline Product API problem families are intentionally small:

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
- `referenced-resource-conflict`;
- `internal-error`.

Problem type is an absolute stable URI; the exact durable documentation/problem host is finalized with OpenAPI/tooling/serving topology rather than invented in W2.

Validation errors may include a bounded `errors[]` extension with machine-readable JSON Pointer/location and detail. Do not add an independent duplicate global `code` taxonomy by default.

Raw provider/business-system errors, payloads, secrets and arbitrary PII never become Product problem truth.

## 13. Cross-Organization reference privacy

An Organization-scoped body/query reference that does not resolve validly inside the path Organization fails closed as a validation/reference problem without disclosing another Organization's ownership/existence details.

A path resource that is not visible/resolvable in the requested Organization may return ordinary `404` semantics rather than proving cross-tenant existence.

## 14. Business outcomes are not Problem Details

The following remain owner/product semantics where the owner contract admits them:

- approval-required / Governance pending/rejected;
- provider capability unsupported / external-required;
- unknown/unavailable/partial source/market evidence;
- ambiguous possible external acceptance;
- economics insufficient evidence;
- divergent convergence;
- business rejection under valid current owner state.

Semantic unavailability of an external evidence source can therefore appear in a successful Product Q representation; it does not automatically become `503` unless MPC itself cannot satisfy the Product API contract because of an API/server failure.

## 15. Complete W2 family coverage

W2-A→E collectively covers the admitted Product schema families:

- Identity/Access — W2-E;
- Marketplace Portfolio — W2-E;
- Readiness — W2-C;
- Offering/Price — W2-B;
- Availability — W2-B + W2-E owner configuration mechanics;
- Market Intelligence — W2-C;
- Commercial Economics — W2-C + W2-E owner configuration mechanics;
- Governance — W2-D;
- Marketplace Sales — W2-D;
- Business-System Materialization — W2-D;
- Fulfillment — W2-D + W2-E owner configuration/media/transport consistency where applicable;
- Post-Sale — W2-D;
- Operational Work — W2-D.

No admitted Product family is intentionally left to a provider DTO, legacy OpenAPI or D6 screen shape.

## 16. Work closure-path audit — closed for current Product 1.0 evidence

The W2-D audit found a legitimate owner path for the currently proven Work-producing condition classes through owner automatic resolution/reconciliation or an admitted owner-specific capability.

Generic `SubmitWorkResolution` therefore remains **DEFER**.

Reopen only on concrete evidence of a human-only resolution/evidence class where the source owner cannot resolve automatically and no legitimate owner-specific Product capability exists. If that occurs, admit the smallest bounded Work→source evidence-submission capability without transferring source truth to Work.

## 17. W2-E proof / negative controls

Later contract/conformance proof must be able to falsify at least:

1. generic Policy/Rules/Settings API becoming cross-owner policy authority;
2. provider deadline represented as editable MPC policy;
3. AccessContext leaking raw IdP/token/session authority;
4. synthetic Membership/AccessRole IDs introduced solely for REST symmetry where semantic relation/key is sufficient;
5. MarketplaceInstallation Product schema accepting credentials/secrets;
6. arbitrary SellingEntity/provider account identifiers becoming MPC identity;
7. ListingIntent media intake still depending on unspecified Product wire or generic Asset/upload-session API;
8. base64/external URL accepted as baseline trusted authored-media intake;
9. one generic status collapsing owner outcome, effect and convergence;
10. `202` used merely because an external effect has not converged;
11. `If-Match` used to condition a referenced resource other than the selected request URI;
12. stale referenced-resource validator incorrectly reported as `412` against another request target;
13. exact idempotent replay of an already accepted intake failing only because the first call advanced the ETag;
14. same idempotency key with changed material precondition being treated as the same request;
15. idempotency key authorizing provider/business-system redispatch after ambiguous possible acceptance;
16. raw-byte JSON differences causing duplicate-intake mismatch for semantically identical requests;
17. business rejection/unknown/unavailable/partial mapped to generic API Problem Details by convenience;
18. Problem Details exposing provider secrets/raw PII/errors;
19. cross-Organization reference validation disclosing another Organization's ownership;
20. an admitted Product operation family having no W2 schema home;
21. a currently proven Work condition requiring generic Work resolution despite an owner closure path.

## 18. Method outcome

**Parent D0→D5-B1 / B2 / Wire W1 / W2-A/B/C/D:** `CURRENT STRUCTURE CONFIRMED`.

**W2-E:** bounded cross-cutting mechanisms added where they eliminate proven failure classes; no new business authority introduced.

> **Owner-specific semantics remain owner-specific; the Product API shares only the wire mechanics whose correctness meaning is genuinely common: exact values, typed references, conditional requests, referenced-resource validators, idempotency intake semantics, Problem Details, bounded inherit/override mechanics and concrete media intake transport.**

No parent-stage reopen is required by W2-E itself.

## 19. Exact next action

Run a **Whole-W2 Global Coherence Review** across W2-A/B/C/D/E before any independent reviewer is invoked.

The lead review must attack at least:

- duplicate/missing identity and authority;
- same meaning with two wire homes or two IDs;
- universal wrapper/subject/scope/policy/result abstractions emerging from repeated mechanics;
- Product 1.0 operation reachability and family coverage;
- cross-owner/client orchestration accidentally required by the schemas;
- source/provider ontology leakage;
- knowledge/freshness/coverage/economic exactness honesty;
- outcome/effect/convergence separation;
- direct versus referenced-resource precondition correctness;
- idempotency retry/concurrency interaction;
- Problem Details versus valid business outcome separation;
- physical-fact client-class enforcement;
- Work closure paths;
- media semantics versus D7 mechanics;
- YAGNI/future retrofit/Structural Inversion against legacy API shape.

If the lead review finds a material contradiction, record the smallest W2/B2-local correction or targeted parent reopen; do not silently rewrite accepted sections.

Only after the Whole-W2 package converges should the lead prepare one bounded **NON-AUTHORITATIVE W2 review candidate** and invoke Fable through the canonical Standard Fable review workflow.

Implementation remains blocked until D9.
