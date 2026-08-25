# D5-B2 — W1 Resource / Path / HTTP Grammar

> **Status:** ACCEPTED / CURRENT CONSOLIDATED AUTHORITY  
> **Parent:** `D5-B2-PRODUCT-OPERATION-SURFACE.md`  
> **Machine wire:** `contracts/api/product/openapi.yaml`

## 1. Governing invariant

> **Product paths identify stable MPC scope/resources or explicitly source-qualified external resources; HTTP methods express honest resource-state semantics or explicit owner capabilities. URI hierarchy never defines business ownership, provider topology or workflow stage.**

W1 owns path/resource/custom-capability grammar. Exact current spelling is the canonical OAD.

## 2. Server / compatibility baseline

No `/v1` path-prefix baseline exists while no production compatibility consumer requires one. Exact host/deployment base remains D7/runtime. Hard cutover is allowed before real compatibility obligations exist.

## 3. Organization scope

Every Organization-owned Product operation is rooted at:

```text
/organizations/{organization_id}/...
```

Path Organization is a scope claim, never self-authorization. Secondary Organization-owned references resolve within that same Organization unless an explicitly accepted cross-Organization relationship exists. No ambient/default tenant header exists.

The one bounded platform-scoped self-discovery operation remains:

```http
GET /access-context
```

It accepts no Principal/Organization selector; Principal comes from trusted authentication context.

## 4. Resource identity grammar

### 4.1 MPC-owned durable resources

Externally addressable MPC-owned identities use Organization root + resource collection + opaque MPC ID, proportionately, e.g.:

```text
marketplace-installations/{marketplace_installation_id}
listing-intents/{listing_intent_id}
price-intents/{price_intent_id}
inventory-sources/{inventory_source_id}
authorization-requests/{authorization_request_id}
authorization-decisions/{authorization_decision_id}
authorization-delegations/{authorization_delegation_id}
business-order-intents/{business_order_intent_id}
invoicing-intents/{invoicing_intent_id}
fulfillment-executions/{fulfillment_execution_id}
fulfillment-nodes/{fulfillment_node_id}
post-sale-resolutions/{post_sale_resolution_id}
work/{work_id}
notifications/{notification_id}
```

Resource nesting does not transfer lifecycle ownership. Example: ListingIntent remains Organization-owned Offering identity even though it targets a Marketplace Installation.

### 4.2 Source-qualified external resources

External resources keep their namespace qualifier and native key. Marketplace-native examples:

```text
marketplace-installations/{marketplace_installation_id}/listings/{native_listing_key}
marketplace-installations/{marketplace_installation_id}/sales/{native_sale_key}
marketplace-installations/{marketplace_installation_id}/shipments/{native_shipment_key}
```

No synthetic MPC mirror ID is minted merely for REST aesthetics.

## 5. Collection / point grammar

A collection path names a real enumerable owner population. Point resource paths exist only where stable point identity/current-owner meaning is admitted. Collections do not exist merely because a point GET exists and vice versa.

Query/filter/pagination semantics belong W3; URI nouns do not encode sort, workflow status or screen grouping.

## 6. Custom owner capabilities

When ordinary CRUD semantics would lie, use explicit colon-suffixed owner capability operations, e.g.:

```text
.../{resource_id}:deactivate
.../{listing_intent_id}:submit
.../{listing_intent_id}:discard
.../product-channel-readiness:resolve-correspondence
.../{work_id}:assign
.../{work_id}:hold
.../{work_id}:resume
.../{work_id}:escalate
.../{authorization_request_id}:decide
```

A custom operation still belongs to one semantic owner; colon syntax does not create a generic command/action model.

## 7. AuthorizationRequest current route

Current exact Governance decision path is:

```http
GET  /organizations/{organization_id}/authorization-requests
GET  /organizations/{organization_id}/authorization-requests/{authorization_request_id}
POST /organizations/{organization_id}/authorization-requests/{authorization_request_id}:decide
```

The first two are exact-human actionable views, not general Governance history browse. `CreateAuthorizationDecision` is the POST operationId at `:decide`; the target is obtained from the Request, not resubmitted by the caller.

AuthorizationDecision history remains separately addressable under `/authorization-decisions` for principals with the admitted Governance read access.

No public generic create/invalidate/reauthorize/list-all/recipient-resolution AuthorizationRequest routes are admitted.

## 8. Personal Notifications current routes

Current Personal Notifications Product routes are:

```http
GET   /organizations/{organization_id}/notifications
PATCH /organizations/{organization_id}/notifications/{notification_id}

GET   /organizations/{organization_id}/notification-routes
GET   /organizations/{organization_id}/notification-route-recipient-candidates
PATCH /organizations/{organization_id}/notification-routes/{notification_kind}
```

The Inbox is exact-self awareness; there is no caller-supplied recipient Principal. Notification routing is Organization configuration, not provider notification/webhook ingress.

No generic `/subscriptions`, `/preferences`, `/delivery-channels`, `/notification-templates`, mark-all-read or source-action route is admitted.

## 9. HTTP method law

- `GET` reads current owner meaning/read projection and has no business side effect.
- `POST` creates durable owner state or invokes an owner capability when creation/custom semantics require it.
- `PATCH` changes an existing owner resource/configuration partially where current contract admits that meaning.
- `DELETE` is not used to imply business deactivation/terminal meaning when explicit capability/resource state is more honest.
- `PUT` is not introduced merely for replacement symmetry.

Exact admitted method per operation is frozen in canonical OAD.

## 10. Idempotency / precondition carriers

Consequential create/intake operations use required `Idempotency-Key` where admitted. Stable idempotency keys identify the same semantic request, not authorization or provider replay permission.

Owner resource updates use typed ETag/revision preconditions where stale-write prevention is material, usually `If-Match` plus 412/428 semantics.

Current deliberate custom exception:

```text
CreateAuthorizationDecision
→ Idempotency-Key header
→ body: AuthorizationRequest StrongETag + outcome
→ no If-Match header
→ no target/target ETag in body
```

The Request ETag is the concurrent terminal-decision carrier; material business validity remains separately revalidated.

## 11. Responses / privacy

Product responses preserve authentication/access/privacy, validation, conflict, stale-write and internal-failure distinctions where applicable. Organization/resource privacy may use 404 rather than disclosing inaccessible resource existence.

A transport/provider result never silently strengthens Product effect/convergence truth.

Exact Problem/status/header spelling remains canonical OAD/W2.

## 12. Product vs technical ingress

Provider/business-system callbacks, webhooks, import/acquisition endpoints and protocol auth are technical D4 ingress, not Product operation symmetry. They do not enter the 106-operation Product census unless a separate real Product client need is admitted.

## 13. Forbidden route patterns

Reject by default:

- bounded-context/package names as URI roots merely to mirror code;
- generic `/actions`, `/commands`, `/mutations`, `/workflow`;
- provider ontology roots exposed as MPC Product semantics;
- screen/dashboard/BFF-shaped paths;
- bare external IDs without their accepted namespace qualifier;
- hidden/default Organization/Installation/SourceInstance;
- paths encoding mutable workflow stage as resource identity.

## 14. Reopen trigger

Reopen W1 only when an admitted Product meaning cannot be represented honestly by existing resource/scope/capability grammar. Route-style preference or implementation convenience is not evidence.
