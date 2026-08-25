# D5-B2 — Operation Admission Matrix

> **Status:** ACCEPTED / CURRENT CONSOLIDATED SEMANTIC ADMISSION AUTHORITY  
> **Parent:** `D5-B2-PRODUCT-OPERATION-SURFACE.md`  
> **Exact current machine census:** `contracts/api/product/openapi.yaml` — **106 operations**

## 1. Admission rule

Every Product operation is admitted only for a real Product 1.0 consumer and declares:

```text
owner
Q | C | P class
Organization/source-qualified subject
ordinary Permission or bounded authenticated/self condition
allowed Principal kinds H/A/S
honest read/outcome semantics
C consequence/idempotency/precondition tuple when applicable
```

No operation is admitted for symmetry, provider endpoint parity, internal debugging, legacy compatibility or frontend-screen convenience.

The exact 106 operationId/path/method/owner/class/access mapping is the canonical OAD projection and is mechanically protected. This document records **why/family admission**, avoiding a second manually maintained 106-row wire table.

## 2. Current admitted families

| Family / owner | Admitted Product jobs | Important exclusions |
| --- | --- | --- |
| D2 Identity / ordinary access | self access-context discovery; Organization member/role reads; assign/revoke AccessRole | credentials/token issuance, custom role designer, generic IAM |
| Marketplace Portfolio | Installation list/detail/create/config/deactivate; Selling Entity list | provider credentials/protocol admin as Product semantics |
| Product & Channel Readiness | source Product search; current readiness; publication requirements; resolve/clear correspondence | Product master/PIM, generic mapping/rules engine |
| Marketplace Offering Operations | Listing list/detail observations; ListingIntent list/detail/create/update/discard/submit/media; PriceIntent list/detail/create | provider-direct CRUD ontology, generic mutation owner |
| Availability Control | sellable availability list/point; InventorySource lifecycle; allocation-scope policy | native stock master, Offering writes |
| Market Intelligence | competitive position list/point; comparable offers | generic market-data query DSL, price mutation |
| Marketplace Performance Intelligence | marketplace summary; Listing performance list/point; retail-media performance list | generic analytics/metrics dimensions/aggregation API |
| Commercial Economics | expected economics list/point; scenario evaluation; sale economics; performance summary; commercial policy; economic attribution list/detail/resolve | price actuation, accounting/treasury API |
| Controlled Action Governance | actionable AuthorizationRequest list/detail; decide exact Request; AuthorizationDecision history; Delegation lifecycle | generic approval/case workflow, source action execution |
| Marketplace Sales | Sale list/detail; resolve Selling Entity attribution | materialization/fulfillment/post-sale ownership |
| Business-System Materialization | BusinessOrderIntent/InvoicingIntent list/detail; bounded party/destination resolution | ERP-native TOP/document API |
| Fulfillment Lifecycle | execution list/detail; Node lifecycle; operating targets; physical checkpoints; artifacts; Shipment list/detail | company-wide WMS/TMS |
| Post-Sale Resolution | resolution list/detail/create | generic CRM/SAC/refund provider API |
| Operational Work | Work list/detail + assign/clear/hold/resume/escalate | originating business truth/authorization |
| Personal Notifications | self Inbox list/update; routing list/set; routing-recipient candidates | source mutation, generic subscriptions/delivery channels |

## 3. Current additive/rebased admissions

### 3.1 Performance Intelligence

Four read operations were later admitted under `performance.read`, H/A/S, exact Marketplace Installation/period semantics:

```text
GetMarketplacePerformanceSummary
ListMarketplaceListingPerformance
GetMarketplaceListingPerformance
ListRetailMediaPerformance
```

### 3.2 Personal Notifications

Five operations were later admitted:

```text
ListMyNotifications                     H / authenticated-self
UpdateMyNotificationAwarenessState      H / authenticated-self
ListNotificationRoutes                  H / notifications.manage
ListNotificationRouteRecipientCandidates H / notifications.manage
SetNotificationRoute                    H / notifications.manage
```

Self surfaces require current Organization Membership and exact self-recipient semantics. Candidate discovery exposes only bounded Principal identity/presentation needed for routing, never role/Permission internals.

### 3.3 AuthorizationRequest

Two actionable reads were added, both H / `governance.decide`:

```text
ListMyActionableAuthorizationRequests
GetMyActionableAuthorizationRequest
```

`CreateAuthorizationDecision` remains the admitted decision operation but is request-scoped at:

```text
POST /organizations/{organization_id}/authorization-requests/{authorization_request_id}:decide
```

with mandatory `Idempotency-Key`, AuthorizationRequest ETag + outcome and no target in the body.

Explicitly rejected/deferred public operations include:

```text
CreateAuthorizationRequest
InvalidateAuthorizationRequest
ReauthorizeAuthorizationRequest
ListAllAuthorizationRequests
ResolveAuthorizationRequestRecipients
CreateNoApproverWork
```

AuthorizationRequest generation/invalidation/recovery remains owner-internal/domain communication semantics, not a generic public workflow surface.

## 4. C-operation safety law

Every C declares independently:

- consequence class;
- idempotency/duplicate-intake disposition;
- concurrency/precondition disposition.

Consequential create/intake uses caller idempotency unless a named structural owner anchor is explicitly sufficient. Resource update uses owner revision/ETag when stale-write prevention is material. Neither mechanism grants business validity or makes ambiguous external replay safe.

Special current case: `CreateAuthorizationDecision` uses a request-body typed AuthorizationRequest ETag rather than HTTP `If-Match`, because the custom capability decides that exact request episode and must remain idempotent under caller retry.

## 5. Collection/search admission

Collections/search exist only for accepted enumerable populations and current human/machine jobs. Pagination/filter/query vocabulary belongs W3; there is no generic list/filter/sort/search framework.

`SearchSourceProductsForMarketplace` remains the only Product Search operation; it is bounded source-qualified discovery, never an MPC Product universe/master API.

## 6. Access admission

W4 owns exact Permission/Principal laws. Current ordinary Permission vocabulary is 31; there is no implicit manage→read or wildcard hierarchy.

`authenticated` on bounded self-only operations is not a stored Permission.

## 7. Reopen trigger

A new Product operation requires a proven consumer and an honest existing/newly reopened semantic owner. If the desired operation exists only to simplify one screen, expose provider topology, mirror CRUD or create symmetry, reject it or return to the smallest upstream Product decision that actually changed.
