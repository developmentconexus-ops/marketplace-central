# D5-B2 — Product Operation / Resource Surface

> **Status:** ACCEPTED / CURRENT CONSOLIDATED AUTHORITY  
> **Parent:** `D5-API.md`  
> **Exact machine wire:** `contracts/api/product/openapi.yaml`  
> **Current census:** **106 Product operations / 31 ordinary Permissions / H-A-S**

## 1. Governing admission law

A Product operation exists only for a real Product 1.0 client job and maps to:

```text
real consumer/use
→ one semantic owner or bounded D2 substrate authority
→ Q | C | P
→ explicit Organization scope where business-owned
→ canonical/source-qualified subject identity
→ exact ordinary Permission or bounded authenticated/self condition
→ admitted Principal class(es)
→ honest read/outcome contract
→ complete C safety semantics when applicable
```

Provider endpoint symmetry, old code/routes, screen convenience, debug access or hypothetical future capability do not admit Product operations.

No Product operation is provider protocol ingress. Technical provider/business-system callbacks/import mechanics remain outside the Product census.

## 2. Client classes

Only D2 Principal kinds exist:

```text
H — human
A — automation
S — bounded system
```

Human Product access uses the accepted server-side-session browser profile; non-human Product access uses the accepted machine bearer profile. Authentication grant type is not Principal kind.

Operation-specific Principal admission is exact and is not inferred from Permission naming.

## 3. Operation classes

- **Q:** current owner meaning / legitimate owner read projection.
- **C:** owner accepts/performs owner-owned work or changes owner-owned state.
- **P:** read-only multi-owner composition when independently justified.

A `POST` is not automatically C by HTTP mechanics; semantic class is owned here. Conversely, a C does not imply an external provider write or completion.

## 4. Complete C safety tuple

Every admitted C has explicit:

```text
consequence class
idempotency disposition
concurrency/precondition disposition
```

Silence is non-conformant.

- consequential creation/intake defaults to caller idempotency unless a named structural owner anchor is safer/sufficient;
- idempotency never substitutes for concurrency;
- concurrency/precondition never substitutes for idempotency;
- caller retry never waives no-blind-retry after ambiguous external acceptance.

## 5. Current operation census authority

The exact 106 operation IDs, paths, methods, semantic owners, classes, required Permissions and Principal kinds are embedded in the canonical OAD and protected mechanically by the Product proof. This file intentionally does **not** duplicate a second 106-row machine catalog.

`D5-B2-OPERATION-ADMISSION-MATRIX.md` records semantic admission by Product family and the later additive/rebased surfaces. W1–W4 freeze path/schema/collection/access laws.

Any operation count change requires explicit Product operation admission. Schema/frontend work cannot grow the operation surface incidentally.

## 6. Accepted post-baseline additive/rebased surface

The original ratified Product surface was later extended/rebased by accepted planning without changing the admission law.

### 6.1 Performance Intelligence — +4 Q operations

One accepted owner, `MarketplacePerformanceIntelligence`, exact Installation/period semantics, `performance.read`, H/A/S:

```text
GetMarketplacePerformanceSummary
ListMarketplaceListingPerformance
GetMarketplaceListingPerformance
ListRetailMediaPerformance
```

No generic analytics/metrics/strategy query language was admitted.

### 6.2 Personal Notifications — +5 operations

```text
ListMyNotifications
UpdateMyNotificationAwarenessState
ListNotificationRoutes
ListNotificationRouteRecipientCandidates
SetNotificationRoute
```

Self Inbox read/update is H-only authenticated/current-Membership personal awareness, not an ordinary `notifications.read` Permission. Organization routing administration/candidate discovery uses `notifications.manage` under the accepted disclosure fences. No source-mutation action is admitted through Notification.

### 6.3 AuthorizationRequest — +2 net operations and one route rebase

New actionable reads:

```text
ListMyActionableAuthorizationRequests
GetMyActionableAuthorizationRequest
```

Both are H-only `governance.decide` and expose only the exact Principal's currently actionable Governance work.

The already-admitted `CreateAuthorizationDecision` remains one operation but its correct current resource/capability home is:

```text
POST /organizations/{organization_id}/authorization-requests/{authorization_request_id}:decide
```

Its request uses the AuthorizationRequest ETag + outcome and mandatory Idempotency-Key. It does not accept a target/target ETag in the decision body.

No public `CreateAuthorizationRequest`, invalidate/reauthorize request, all-request browse, recipient-resolution or generic approval-workflow operation is admitted.

## 7. Stable family dispositions

Current Product 1.0 continues to expose only bounded operations under these meanings:

- D2 identity/access self discovery and Organization member/role administration;
- Marketplace Portfolio Installation/SellingEntity configuration;
- Readiness search/current readiness/publication requirements/correspondence resolution;
- Offering Listing observations + ListingIntent/PriceIntent authoring;
- Availability current sellable availability + InventorySource/allocation policy;
- Market competitive/comparable evidence;
- Performance bounded provider measurement evidence;
- Economics expected/sale economics, scenario evaluation, policy and attribution;
- Governance actionable AuthorizationRequest, Decisions and Delegations;
- Marketplace Sales observations/attribution;
- Business-System Materialization intents/resolution;
- Fulfillment executions/nodes/checkpoints/artifacts + Shipment observations;
- Post-Sale resolutions;
- Operational Work lifecycle;
- Personal Notifications self Inbox and bounded Organization routing.

No generic CRUD or symmetry requirement exists across families.

## 8. Forbidden additions by convenience

Without explicit re-admission, reject:

- generic `/actions`, `/commands`, `/mutations`, `/workflow`, `/analytics`, `/metrics` operations;
- Product-level provider auth/callback/admin plumbing;
- source/provider raw browsing for debugging;
- generic entity/relationship/evidence APIs;
- screen-shaped aggregate endpoints;
- hidden default Organization/Installation/SourceInstance;
- general Notification subscriptions/delivery-channel APIs;
- generic AuthorizationRequest workflow/case APIs;
- bulk APIs that erase member-level authority/outcomes.

## 9. Reopen rule

A new operation is justified only by a proven Product consumer and one honest owner meaning that cannot be expressed through an already-admitted operation without distorting semantics. Reopen the smallest D5/parent owner; never add the operation merely because implementation or frontend composition would be easier.
