# D5-B2 — Wire Contract / Resource-Path-HTTP Grammar

> **Status:** W1 RESOURCE / PATH / HTTP GRAMMAR — ACCEPTED / CANONICAL  
> **Parent B2:** `D5-B2-PRODUCT-OPERATION-SURFACE.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + Decision Reconciliation Baseline + ratified D5-B2 Whole-Matrix  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Whole-W2 final ratification incorporated:** 2026-08-19  
> **Final Problem/media consistency cross-reference incorporated:** 2026-08-19

## 1. Purpose

W1 defines the canonical Product API resource/path/HTTP grammar from the accepted semantic owners and operation inventory.

> **Product API paths identify stable MPC scope/resources or explicitly source-qualified external resources; HTTP methods express honest resource-state semantics or explicit owner capabilities. URI hierarchy never defines business ownership, provider topology or workflow state.**

W1 does not choose D6 screen/BFF topology, D7 runtime/storage/transaction realization, D8 live-effect proof or implementation.

---

## 2. Server and compatibility baseline

Product paths are relative to the OpenAPI server URL.

There is **no `/v1` path prefix baseline**. No entitled production compatibility consumer currently exists, so introducing a public version axis would create unsupported compatibility tax. Reopen versioning strategy only when a real stable consumer requires an incompatible-evolution contract.

Exact host/deployment base path remains D7/runtime realization.

---

## 3. Organization path law

All Organization-owned Product API operations use:

```text
/organizations/{organization_id}/...
```

The path Organization is a scope claim, never self-authorization. Current Membership, Permission and client-class rules still apply.

Organization-scoped secondary references in body/query structures must resolve inside the path Organization unless a later accepted cross-Organization relation explicitly says otherwise. There is no ambient/default tenant header.

---

## 4. Platform-scoped self discovery

The only baseline platform-scoped Product Q is:

```http
GET /access-context
```

Rules:

- Principal is derived only from trusted authentication context;
- no caller-supplied `principal_id` or Organization parameter;
- returns current MPC Principal kind, visible Organization memberships, AccessRole keys and effective ordinary Permissions;
- it is not `/me`, `/session` or a generic Principal/profile API.

---

## 5. Public paths do not expose D1 bounded-context taxonomy

D1 owner/context names are not mandatory URI prefixes. Reject domain-namespaced roots such as `/offering`, `/availability`, `/economics` or `/governance` merely to mirror code boundaries.

Public URI structure expresses stable scope/identity, not internal package/module taxonomy.

---

## 6. MPC-owned canonical resource grammar

An MPC-owned canonical resource uses Organization root + resource collection + opaque MPC ID, proportionately:

```text
/organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}
/organizations/{organization_id}/listing-intents/{listing_intent_id}
/organizations/{organization_id}/price-intents/{price_intent_id}
/organizations/{organization_id}/inventory-sources/{inventory_source_id}
/organizations/{organization_id}/fulfillment-nodes/{fulfillment_node_id}
/organizations/{organization_id}/fulfillment-executions/{fulfillment_execution_id}
/organizations/{organization_id}/post-sale-resolutions/{post_sale_resolution_id}
/organizations/{organization_id}/work/{work_id}
```

Exact remaining collection nouns are later OpenAPI spelling against accepted semantics. Resource identity is not nested under Marketplace Installation merely because the Installation is a target/reference.

Example: ListingIntent remains Organization-level MPC identity and carries `marketplace_installation_id` as context; Installation is not its lifecycle owner.

---

## 7. Source-qualified external resource grammar

Externally authoritative resources retain their namespace qualifier and native opaque key; MPC does not mint mirror IDs for REST aesthetics.

Marketplace-native point resources use Installation qualification, for example:

```text
/organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}/listings/{native_listing_key}
/organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}/sales/{native_sale_key}
/organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}/shipments/{native_shipment_key}
```

The nesting is source namespace qualification, not Portfolio ownership of Listing/Sale/Shipment meaning.

Business-system external references remain SourceInstance-qualified when operation scope does not already make SourceInstance unambiguous.

---

## 8. Keyed Q/current meaning without fake identity

Owner current meanings or views do not receive synthetic IDs merely so every read looks like `/resources/{id}`.

`ProductChannelReadiness`, `CompetitivePosition`, `ExpectedEconomics`, `EconomicPerformanceSummary` and similar meanings may use semantically keyed path/query subjects without `ReadinessId`, `CompetitivePositionId`, `ExpectedEconomicsId` or Summary ID.

ProductChannelCorrespondence remains a keyed Readiness-owned meaning inside the Product + Marketplace subject. Resolve/Clear keep owner capability semantics and **do not** gain a synthetic Correspondence resource ID merely to obtain concurrency.

---

## 9. URI nesting invariant

Nesting means only:

1. identity/lifecycle containment; or
2. external namespace qualification.

It never means workflow order.

Do not encode process chains such as Sales → Materialization → Fulfillment → Invoicing in URI hierarchy merely because one stage follows another.

Contained singleton meanings are allowed when identity truly belongs to the parent, e.g. BusinessOrderIntent PartyResolution/DestinationRealization.

---

## 10. Standard HTTP methods when honest

Use ordinary HTTP resource methods only when they match accepted owner semantics:

- `GET` — read current owner meaning/resource;
- `POST` collection — create when a new MPC canonical resource/occurrence is genuinely created and server assigns identity;
- `PATCH` — typed partial update of a mutable MPC resource/configuration;
- `PUT` — complete replacement only where that is genuinely the resource meaning;
- `DELETE` — actual resource/relation removal only where deletion semantics are true.

Do not convert capabilities to PATCH/PUT/DELETE solely to reuse generic HTTP tooling.

---

## 11. Owner-specific capability grammar

When CRUD would lie, use explicit owner capability syntax:

```text
POST {resource-or-keyed-subject-uri}:verb
```

Examples:

```text
POST /organizations/O/listing-intents/LI:submit
POST /organizations/O/listing-intents/LI:discard
POST /organizations/O/listing-intents/LI:create-media
POST /organizations/O/economic-attributions/EA:resolve
POST /organizations/O/work/W:hold
POST /organizations/O/work/W:resume
POST /organizations/O/work/W:escalate
POST /organizations/O/marketplace-installations/I:deactivate
```

Rules:

- verb remains owner vocabulary, never generic `/actions`, `/commands`, `/operations`;
- a custom-method URI is a distinct HTTP request target; semantic attachment to a base resource does **not** make it an alias representation of the base URI;
- custom methods must not be disguised writable-state enums.

---

## 12. No writable workflow/status shortcut

Never encode consequential owner capabilities as:

```text
PATCH status = submitted
PATCH status = approved
PATCH status = packed
PATCH status = closed
```

when the real meaning is submit, authorize, establish physical evidence, reconcile or coordinate owner-specific consequence.

Lifecycle/state may be present in responses; it is not automatically a writable API field.

---

## 13. One opaque validator authority

Mutable/current MPC meanings that require stale-state protection expose a **strong opaque server-issued ETag** representing the material owner revision.

The validator is not a database sequence, timestamp, provider x-version, Postgres xmin or client-authored version number.

There is one revision authority per protected meaning. Transport may differ by operation class; validator meaning does not.

---

## 14. Revision-precondition carrier grammar — canonical

### 14.1 Standard method on the exact protected resource URI

When the HTTP request target is the actual conditionally protected resource representation, use RFC conditional semantics:

```text
ETag: "opaque-validator"
If-Match: "opaque-validator"
```

Examples: typed PATCH of ListingIntent, MarketplaceInstallation, InventorySource, CommercialPolicy, AuthorizationDelegation, FulfillmentNode and FulfillmentOperatingTargets.

Failures:

- required `If-Match` absent → `428 Precondition Required` + Problem Details;
- supplied `If-Match` false/stale → `412 Precondition Failed` + Problem Details.

### 14.2 Owner custom method

A custom method such as `/resource/{id}:verb` does **not** use the base resource's ETag as `If-Match`, because the custom-method URI is a different HTTP request target unless explicit alias semantics are separately created. MPC creates no such alias merely to make the header fit.

When the ratified safety tuple requires current owner state, the custom request carries the acted-on resource's same opaque validator as typed technical request data:

```json
{
  "etag": "\"opaque-validator\""
}
```

A custom method may therefore have no **business payload** while still carrying technical revision proof.

Multipart capabilities use a typed `etag` part.

Failures:

- required typed ETag missing/invalid → `422 validation-error`;
- supplied typed ETag stale → `409 resource-revision-conflict`.

### 14.3 Exact revision of another referenced resource

If a create/capability depends on another MPC resource's exact revision, the typed reference carries that resource's ETag adjacent to the reference.

Examples:

- AuthorizationDecision targeting an exact Intent revision;
- PriceIntent explicitly superseding an exact prior PriceIntent revision.

The same `422` missing/invalid and `409 resource-revision-conflict` grammar applies.

### 14.4 ProductChannelCorrespondence

`GetProductChannelReadiness` exposes a **correspondence-scoped opaque `etag`** distinct from `requirements_revision` and readiness/source-evidence revisions.

`ResolveProductChannelCorrespondence` and `ClearProductChannelCorrespondence` remain Readiness owner capabilities over the keyed Product + Marketplace subject and carry that typed ETag.

Do not force PUT/DELETE solely for conditionals: correspondence stale-state safety must also protect unresolved/conflicting current meaning, where a standing resolved relation may not exist. No synthetic Correspondence/Readiness ID is introduced.

---

## 15. Idempotency and concurrency are independent

`Idempotency-Key` protects duplicate consequential intake where the ratified operation matrix requires a client key.

Revision proof protects stale current-state decisions.

They are never interchangeable:

```text
Idempotency-Key != ETag / If-Match / typed etag
```

For client-keyed operations, exact prior intake is resolved before re-evaluating a now-stale revision proof on a lost-response retry. A changed ETag/precondition under the same key is a materially different request.

Idempotency never authorizes blind redispatch after ambiguous possible external acceptance.

Exact intake equivalence/problem grammar is owned by canonical W2.

---

## 16. No provider-shaped Product API roots

Reject provider/business-system top-level Product roots such as:

```text
/mercado-livre
/sankhya
/providers
/integrations
```

Provider OAuth/webhook/callback/connector ingress is a separate D4 protocol surface and is not Product API business vocabulary.

---

## 17. OpenAPI version decision remains tooling-scoped

One machine-readable OpenAPI document remains the Product API wire authority.

The exact OpenAPI minor is **not selected by recency alone**. Choose it only when W3/later schema requirements and actual server/client generator compatibility prove the smallest suitable version.

No second manually authoritative SDK/schema representation is admitted.

---

## 18. W1 negative controls

Later OpenAPI/contract proof must falsify at least:

1. `/v1` compatibility axis introduced without a consumer;
2. D1 domain names becoming mandatory URI roots;
3. external Listing/Sale/Shipment receiving fake MPC mirror IDs;
4. MPC ListingIntent nested under Installation merely because it references it;
5. workflow-order nesting becoming URI ownership;
6. provider top-level Product API namespace;
7. arbitrary `PATCH status` replacing an owner capability;
8. generic `/actions`/`/commands` escape hatch;
9. `If-Match` carrying a base-resource ETag on a different `:verb` request target by private convention;
10. custom header introduced when typed revision data is sufficient;
11. typed custom/reference ETag becoming a second version authority instead of carrying the resource ETag;
12. ProductChannelCorrespondence receiving a synthetic ID only for concurrency;
13. idempotency used as concurrency authority or revision proof used as duplicate-intake identity.

---

## 19. W1 outcome

**Outcome:** `CURRENT PARENT STRUCTURE CONFIRMED` + ratified Whole-W2 correction incorporated.

> **Use identity-oriented Organization/source-qualified paths, honest standard resource methods, explicit owner capabilities and exactly one opaque revision authority. HTTP `If-Match` is reserved for a request whose literal target is the protected resource; custom/reference revision proofs travel as typed request data.**

No D0→D5-B1 parent reopen is required.

W1 is canonical and closed. W2, W3 and W4 are separately consolidated canonical schema, collection and enforcement authorities, and `D5-B2-TECHNICAL-INGRESS.md` is the canonical technical non-Product ingress authority. Current status and exact next action are owned only by the rebaseline router.

Implementation remains blocked until D9.

---

## 20. Canonical non-Product route routing

This section is the canonical routing statement for externally reachable routes that are not Product API operations, and it sharpens the routing statement in §16. It changes **no Product resource, path, HTTP method, operation or schema semantics**.

Canonical routing is:

- D4 owns provider protocol/authentication/source semantics;
- `D5-B2-TECHNICAL-INGRESS.md` owns their D5-B2 technical wire/trust-boundary crystallization;
- provider OAuth begin/callback, notifications and acquisition custody are not Product operations and are excluded from the Product OpenAPI and SDK;
- exact technical host/prefix/method remains deferred to its proper realization stage;
- any technical route must remain unambiguously outside and must not collide with `/access-context` or `/organizations/{organization_id}/...`;
- no generic Product `/providers`, `/integrations`, `/webhooks`, `/oauth` or `/external-events` root is admitted.

Authored-media byte delivery is a distinct non-Product case:

- it is a **separately justified technical presentation surface**, not a Product operation and not a Technical Ingress lane A or B;
- it is excluded from the Product OpenAPI and SDK;
- it reuses current Product AuthN, unique Principal binding, Principal access eligibility, Organization Membership and `offering.read` for the exact ListingIntent/media relationship, per canonical W2 §3.9.8;
- no durable anonymous or freely forwardable locator is baseline;
- its failures are technical-surface failures and never expand the Product Problem catalog;
- exact route spelling, proxy/storage/CDN topology and streaming/transformation mechanics remain D7, and must not collide with `/access-context` or `/organizations/{organization_id}/...`.

W1 remains **ACCEPTED / CANONICAL**. The Product inventory remains **95 operations / 29 Permissions**. Current status and exact next action are owned only by the rebaseline router.

Implementation remains blocked until D9.
